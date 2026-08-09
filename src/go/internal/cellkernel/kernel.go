// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) single-episode kernel, implemented under the
// operator-ratified FIDO/functional doctrine
// (msg-cn-pi-cnos-cell-runner-fido-functional-44):
//
//	alphaIn  := AlphaInput{frozen contract}
//	aOut     := α(alphaIn)                        Result<AlphaOutput, error>
//	sealedA  := sealAlpha(aOut)                   runtime-owned, immutable
//	betaIn   := BetaInput{contract, projection(sealedA), policy}
//	bOut     := β(betaIn)                         Result<BetaOutput, error>
//	sealedB  := sealBeta(bOut)
//	record   := compose(start, sealedA, sealedB)  one immutable EpisodeRecord
//	receipt  := γ(record)                         canonical bytes + ONE digest
//	verdict  := V(receipt)                        typed failures
//	decision := δ(verdict)
//	status   := lift(decision, mode)
//
// Governing rule: NO MUTABLE SHARED EPISODE STATE. Each seat is a pure-shaped
// function invoked with exactly the immutable data it needs, returning one
// typed value. Ownership is positional and structural — the runtime knows a
// value came from α because it invoked α and received the return. Seats never
// declare producer roles, execution ids, hashes, verdicts, receipts, status,
// or decisions; there are no seat-visible authority fields to forge.
//
// The primary safety mechanism is structural isolation of the untrusted
// cognitive seats, not the trusted runtime proving its own internal steps to
// itself: sealed results carry unexported state (unforgeable outside this
// package), β receives fresh copies (projections) of sealed α output, and the
// single scope-lift digest over the canonical EpisodeRecord is the one proof a
// downstream verifier recomputes (VerifyClosure).
//
// Retained mechanical gates (Pi β #31–#33, PR-#718 β): typed
// semantic-vs-integrity failure routing; fail-closed identity minting;
// explicit non-authoritative stub (`simulated`); size bounds + cancellation
// between stations; explicit UTF-8 artifact text contract with an aggregate
// bound; γ/V/δ kernel-owned and mechanical; only α and β are open seats.
package cellkernel

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"unicode/utf8"
)

// Schema / canon identifiers.
const (
	ClosureSchema     = "cnos.cellkernel.episode-closure.v0"
	RecordCanon       = "cnos.cellkernel.episode-record-canon.v0"
	BetaInputPolicyID = "cnos.cellkernel.beta-input.v0"
)

// Bounds enforced at the kernel boundary.
const (
	maxRequiredEvidence  = 64
	maxMatterBytes       = 1 << 20 // 1 MiB
	maxReviewNotesBytes  = 64 << 10
	maxArtifacts         = 64
	maxArtifactBytes     = 1 << 20 // per artifact
	maxAggregateArtifact = 4 << 20 // sum over both seats
)

// Role names which station a required artifact must come from. The check is
// positional (which side of the record the artifact sits on), never a stamp.
type Role string

const (
	RoleAlpha Role = "alpha"
	RoleBeta  Role = "beta"
)

// ExecutionMode: a stub run is a non-authoritative smoke test.
type ExecutionMode string

const (
	ModeStub       ExecutionMode = "stub"
	ModeMechanical ExecutionMode = "mechanical"
)

func knownMode(m ExecutionMode) bool { return m == ModeStub || m == ModeMechanical }

// --- Immutable values ----------------------------------------------------

type RequiredRef struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Producer Role   `json:"producer"`
}

type Contract struct {
	ID               string        `json:"id"`
	Goal             string        `json:"goal"`
	RequiredEvidence []RequiredRef `json:"required_evidence,omitempty"`
}

func (c Contract) clone() Contract {
	cp := c
	if c.RequiredEvidence != nil {
		cp.RequiredEvidence = append([]RequiredRef(nil), c.RequiredEvidence...)
	}
	return cp
}

type Matter struct {
	Data string `json:"data"`
}

type Review struct {
	Pass  bool   `json:"pass"`
	Notes string `json:"notes"`
}

// ArtifactCandidate is what a seat returns: semantic identity + UTF-8 text.
// No provenance fields exist for a seat to forge.
type ArtifactCandidate struct {
	ID   string
	Kind string
	Text string
}

// Artifact is a runtime-normalized artifact inside the sealed record. Its
// provenance is its position (under alpha or beta), not a stamp.
type Artifact struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Encoding string `json:"encoding"` // "utf8" (base64 is a future extension)
	Text     string `json:"text"`
}

// --- Seat inputs/outputs (immutable scopes) ------------------------------

// AlphaInput is exactly what α may see: an isolated frozen contract copy.
type AlphaInput struct {
	Contract Contract
}

type AlphaOutput struct {
	Matter    Matter
	Artifacts []ArtifactCandidate
}

// BetaInput is the runtime-owned review surface: a fresh frozen contract copy,
// a projection of sealed α output (copies — never α's live scope), and the
// review policy. Nothing of α's private state crosses.
type BetaInput struct {
	Contract       Contract
	Matter         Matter
	AlphaArtifacts []Artifact
	PolicyID       string
}

type BetaOutput struct {
	Review    Review
	Artifacts []ArtifactCandidate
}

// --- Sealed results (runtime-owned; unforgeable outside this package) ----

// SealedAlpha holds α's normalized return. Fields are unexported: no seat or
// external caller can construct or mutate a sealed value; accessors return
// copies, so a β that mutates its projection cannot reach the sealed original.
type SealedAlpha struct {
	exec      string
	matter    Matter
	artifacts []Artifact
}

func (s SealedAlpha) projection() ([]Artifact, Matter) {
	return append([]Artifact(nil), s.artifacts...), s.matter
}

type SealedBeta struct {
	exec      string
	review    Review
	artifacts []Artifact
}

// --- The two open seats ---------------------------------------------------

type Alpha interface {
	Produce(ctx context.Context, in AlphaInput) (AlphaOutput, error)
}

type Beta interface {
	Review(ctx context.Context, in BetaInput) (BetaOutput, error)
}

// Spec is a single-episode cell: a contract plus the two open seats.
type Spec struct {
	Contract Contract
	Alpha    Alpha
	Beta     Beta
}

// --- Identity (fail-closed) ----------------------------------------------

type Identity struct {
	Episode string
	Alpha   string
	Beta    string
}

func (id Identity) valid() error {
	if id.Episode == "" || id.Alpha == "" || id.Beta == "" {
		return errors.New("identity has empty id")
	}
	if id.Episode == id.Alpha || id.Episode == id.Beta || id.Alpha == id.Beta {
		return errors.New("identity ids are not pairwise distinct")
	}
	return nil
}

// IDSource mints the whole identity tuple in one error-returning operation.
type IDSource interface {
	Mint() (Identity, error)
}

type randomIDs struct{}

func (randomIDs) Mint() (Identity, error) {
	parts := [3]string{}
	for i := range parts {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return Identity{}, fmt.Errorf("mint id: %w", err)
		}
		parts[i] = hex.EncodeToString(b)
	}
	return Identity{Episode: "ep-" + parts[0], Alpha: "alpha-" + parts[1], Beta: "beta-" + parts[2]}, nil
}

// --- Episode record, receipt, closure ------------------------------------

// ResolvedSpec is the normalized executable spec, carried whole so the scope-
// lift digest covers it and any verifier can reproduce it.
type ResolvedSpec struct {
	Version          string            `json:"version"`
	DeclaredProtocol string            `json:"declared_protocol"`
	Profile          string            `json:"profile"`
	Params           map[string]string `json:"params,omitempty"`
	AlphaSkills      []string          `json:"alpha_skills"`
	BetaSkills       []string          `json:"beta_skills"`
}

// clone deep-copies the resolved spec so no caller- or seat-retained alias can
// mutate runtime-owned invocation truth (Pi PR-#718-fido β D3).
func (r ResolvedSpec) clone() ResolvedSpec {
	cp := r
	if r.Params != nil {
		cp.Params = make(map[string]string, len(r.Params))
		for k, v := range r.Params {
			cp.Params[k] = v
		}
	}
	// Always non-nil so the canonical JSON is [] rather than null.
	cp.AlphaSkills = append(make([]string, 0, len(r.AlphaSkills)), r.AlphaSkills...)
	cp.BetaSkills = append(make([]string, 0, len(r.BetaSkills)), r.BetaSkills...)
	return cp
}

// StationRecord is one station's contribution, owned positionally.
type StationRecord struct {
	ExecutionID string     `json:"execution_id"`
	Artifacts   []Artifact `json:"artifacts"`
}

// EpisodeRecord is the ONE immutable account of the episode, composed by the
// runtime from the sealed results. Everything downstream derives from it.
type EpisodeRecord struct {
	Canon           string        `json:"canon"`
	EpisodeID       string        `json:"episode_id"`
	Mode            ExecutionMode `json:"execution_mode"`
	ResolvedSpec    ResolvedSpec  `json:"resolved_spec"`
	Contract        Contract      `json:"contract"`
	Alpha           StationRecord `json:"alpha"`
	Matter          Matter        `json:"matter"`
	Beta            StationRecord `json:"beta"`
	Review          Review        `json:"review"`
	BetaInputPolicy string        `json:"beta_input_policy"`
}

// canonicalBytes is the record's canonical serialization (schema-ordered JSON,
// versioned by Canon). The scope-lift digest is computed over exactly these.
func (r EpisodeRecord) canonicalBytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

// Receipt is γ's output: the serialized record plus its single scope-lift
// digest — the one proof any downstream verifier recomputes.
type Receipt struct {
	Record          EpisodeRecord `json:"record"`
	ScopeLiftDigest string        `json:"scope_lift_digest"`
}

// FailureClass separates repairable contract-unmet from fail-closed integrity
// failures.
type FailureClass string

const (
	ContractUnmet   FailureClass = "contract_unmet"
	InvalidRecord   FailureClass = "invalid_record"
	InvalidIdentity FailureClass = "invalid_identity"
)

func (f FailureClass) integrity() bool { return f != ContractUnmet }

type Failure struct {
	Class  FailureClass `json:"class"`
	Detail string       `json:"detail"`
}

type Verdict struct {
	Pass     bool      `json:"pass"`
	Failures []Failure `json:"failures,omitempty"`
}

func (v Verdict) hasIntegrityFailure() bool {
	for _, f := range v.Failures {
		if f.Class.integrity() {
			return true
		}
	}
	return false
}

type Decision string

const (
	Accept         Decision = "accept"
	Release        Decision = "release"
	Override       Decision = "override"
	Reject         Decision = "reject"
	RepairDispatch Decision = "repair_dispatch"
)

type Status string

const (
	Accepted    Status = "accepted"
	Degraded    Status = "degraded"
	Rejected    Status = "rejected"
	NeedsRepair Status = "needs_repair"
	Simulated   Status = "simulated" // stub: non-authoritative
)

type RepairRequest struct {
	Reason string   `json:"reason"`
	Failed []string `json:"failed,omitempty"`
}

// Closure is the terminal object that crosses the scope boundary: receipt,
// verdict, decision, status — all re-derivable from the record by
// VerifyClosure. It carries no other proof surfaces.
type Closure struct {
	Schema            string         `json:"closure_schema"`
	ProtocolValidated bool           `json:"protocol_validated"`
	Status            Status         `json:"status"`
	Decision          Decision       `json:"decision"`
	Verdict           Verdict        `json:"verdict"`
	Receipt           Receipt        `json:"receipt"`
	Repair            *RepairRequest `json:"repair,omitempty"`
}

// --- Run options ----------------------------------------------------------

type RunMeta struct {
	ResolvedSpec  ResolvedSpec
	ExecutionMode ExecutionMode
}

type runConfig struct {
	ids  IDSource
	meta RunMeta
}

type RunOption func(*runConfig)

func WithIDSource(s IDSource) RunOption { return func(c *runConfig) { c.ids = s } }
func WithMeta(m RunMeta) RunOption      { return func(c *runConfig) { c.meta = m } }

var ErrInvalidClosure = errors.New("cellkernel: inconsistent (verdict, decision) pair")

// RunEpisode executes the five-step closure as a pure-shaped pipeline over
// immutable values. A returned error is a malfunction (invalid spec, identity
// failure, seat error, bound violation); otherwise the Closure's Status reports
// how the episode closed.
func RunEpisode(ctx context.Context, s Spec, opts ...RunOption) (Closure, error) {
	cfg := runConfig{ids: randomIDs{}, meta: RunMeta{ExecutionMode: ModeMechanical}}
	for _, o := range opts {
		o(&cfg)
	}
	if !knownMode(cfg.meta.ExecutionMode) {
		return Closure{}, fmt.Errorf("cellkernel: unknown execution mode %q", cfg.meta.ExecutionMode)
	}
	if err := validateSpec(s); err != nil {
		return Closure{}, err
	}
	if err := ctx.Err(); err != nil {
		return Closure{}, fmt.Errorf("cellkernel: context: %w", err)
	}

	if seatIsNil(cfg.ids) { // D5: the identity source fails closed like a seat.
		return Closure{}, errors.New("cellkernel: nil identity source")
	}
	id, err := cfg.ids.Mint() // fail-closed identity, before α runs
	if err != nil {
		return Closure{}, fmt.Errorf("cellkernel: %w", err)
	}
	if err := id.valid(); err != nil {
		return Closure{}, fmt.Errorf("cellkernel: %w", err)
	}

	frozen := s.Contract.clone()
	frozenMeta := RunMeta{ExecutionMode: cfg.meta.ExecutionMode, ResolvedSpec: cfg.meta.ResolvedSpec.clone()}

	// Station α: isolated input → output → sealed.
	aOut, err := s.Alpha.Produce(ctx, AlphaInput{Contract: frozen.clone()})
	if err != nil {
		return Closure{}, fmt.Errorf("alpha produce: %w", err)
	}
	sealedA, err := sealAlpha(aOut, id.Alpha)
	if err != nil {
		return Closure{}, fmt.Errorf("seal alpha: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return Closure{}, fmt.Errorf("cellkernel: context after alpha: %w", err)
	}

	// Station β: fresh projection of sealed α — never α's live scope.
	projArtifacts, projMatter := sealedA.projection()
	bOut, err := s.Beta.Review(ctx, BetaInput{
		Contract:       frozen.clone(),
		Matter:         projMatter,
		AlphaArtifacts: projArtifacts,
		PolicyID:       BetaInputPolicyID,
	})
	if err != nil {
		return Closure{}, fmt.Errorf("beta review: %w", err)
	}
	sealedB, err := sealBeta(bOut, id.Beta)
	if err != nil {
		return Closure{}, fmt.Errorf("seal beta: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return Closure{}, fmt.Errorf("cellkernel: context after beta: %w", err)
	}

	// Compose the one immutable record; then the mechanical tail is pure.
	record, err := compose(id, frozenMeta, frozen, sealedA, sealedB)
	if err != nil {
		return Closure{}, err
	}
	receipt := gamma(record)
	verdict := validate(receipt)
	decision := decide(verdict)
	status, err := lift(verdict, decision, record.Mode)
	if err != nil {
		return Closure{}, err
	}

	cl := Closure{
		Schema:            ClosureSchema,
		ProtocolValidated: false,
		Status:            status,
		Decision:          decision,
		Verdict:           verdict,
		Receipt:           receipt,
	}
	cl.Repair = repairFrom(verdict, status)
	return cl, nil
}

// --- Sealing (runtime-only) -----------------------------------------------

func normalizeArtifacts(cands []ArtifactCandidate) ([]Artifact, error) {
	if len(cands) > maxArtifacts {
		return nil, fmt.Errorf("too many artifacts (%d > %d)", len(cands), maxArtifacts)
	}
	out := make([]Artifact, 0, len(cands))
	for _, c := range cands {
		if c.ID == "" || c.Kind == "" {
			return nil, errors.New("artifact candidate needs non-empty id and kind")
		}
		if len(c.Text) > maxArtifactBytes {
			return nil, fmt.Errorf("artifact %q exceeds %d bytes", c.ID, maxArtifactBytes)
		}
		if !utf8.ValidString(c.Text) {
			return nil, fmt.Errorf("artifact %q is not valid UTF-8 (base64 encoding is a future extension)", c.ID)
		}
		out = append(out, Artifact{ID: c.ID, Kind: c.Kind, Encoding: "utf8", Text: c.Text})
	}
	return out, nil
}

func sealAlpha(o AlphaOutput, exec string) (SealedAlpha, error) {
	if len(o.Matter.Data) > maxMatterBytes {
		return SealedAlpha{}, fmt.Errorf("matter exceeds %d bytes", maxMatterBytes)
	}
	arts, err := normalizeArtifacts(o.Artifacts)
	if err != nil {
		return SealedAlpha{}, err
	}
	return SealedAlpha{exec: exec, matter: o.Matter, artifacts: arts}, nil
}

func sealBeta(o BetaOutput, exec string) (SealedBeta, error) {
	if len(o.Review.Notes) > maxReviewNotesBytes {
		return SealedBeta{}, fmt.Errorf("review notes exceed %d bytes", maxReviewNotesBytes)
	}
	arts, err := normalizeArtifacts(o.Artifacts)
	if err != nil {
		return SealedBeta{}, err
	}
	return SealedBeta{exec: exec, review: o.Review, artifacts: arts}, nil
}

// compose builds the one immutable EpisodeRecord from sealed results (pure).
func compose(id Identity, meta RunMeta, frozen Contract, a SealedAlpha, b SealedBeta) (EpisodeRecord, error) {
	total := 0
	for _, art := range a.artifacts {
		total += len(art.Text)
	}
	for _, art := range b.artifacts {
		total += len(art.Text)
	}
	if total > maxAggregateArtifact {
		return EpisodeRecord{}, fmt.Errorf("cellkernel: aggregate artifact bytes exceed %d", maxAggregateArtifact)
	}
	return EpisodeRecord{
		Canon:           RecordCanon,
		EpisodeID:       id.Episode,
		Mode:            meta.ExecutionMode,
		ResolvedSpec:    meta.ResolvedSpec,
		Contract:        frozen,
		Alpha:           StationRecord{ExecutionID: a.exec, Artifacts: a.artifacts},
		Matter:          a.matter,
		Beta:            StationRecord{ExecutionID: b.exec, Artifacts: b.artifacts},
		Review:          b.review,
		BetaInputPolicy: BetaInputPolicyID,
	}, nil
}

// gamma serializes the record and computes the single scope-lift digest (pure).
func gamma(r EpisodeRecord) Receipt {
	return Receipt{Record: r, ScopeLiftDigest: sha256hex(r.canonicalBytes())}
}

// validateRecord replays the COMPLETE record boundary as typed failures (Pi
// PR-#718-fido β D1): everything the kernel guards on the honest path is
// re-derivable by a verifier from the serialized record alone. Pure; used by V
// before closure and by VerifyClosure at scope lift — one boundary, one
// validator.
func validateRecord(r EpisodeRecord) []Failure {
	var fs []Failure
	add := func(cond bool, class FailureClass, msg string) {
		if cond {
			fs = append(fs, Failure{Class: class, Detail: msg})
		}
	}

	add(r.Canon != RecordCanon, InvalidRecord, "wrong record canon")
	add(!knownMode(r.Mode), InvalidRecord, "unknown execution mode")
	add(r.BetaInputPolicy != BetaInputPolicyID, InvalidRecord, "unknown beta-input policy")

	// Identity: non-empty and pairwise distinct across the whole triple.
	add(r.EpisodeID == "" || r.Alpha.ExecutionID == "" || r.Beta.ExecutionID == "", InvalidIdentity, "missing identity")
	add(r.Alpha.ExecutionID == r.Beta.ExecutionID, InvalidIdentity, "stations share an execution id")
	add(r.EpisodeID == r.Alpha.ExecutionID || r.EpisodeID == r.Beta.ExecutionID, InvalidIdentity, "episode id aliases a station id")

	// Contract validity (same rules validateSpec enforces on the honest path).
	add(r.Contract.ID == "", InvalidRecord, "contract id is empty")
	add(len(r.Contract.RequiredEvidence) > maxRequiredEvidence, InvalidRecord, "too many required evidence refs")
	seenReq := make(map[string]bool)
	for _, req := range r.Contract.RequiredEvidence {
		add(req.ID == "" || req.Kind == "", InvalidRecord, "required evidence with empty id/kind")
		add(req.Producer != RoleAlpha && req.Producer != RoleBeta, InvalidRecord, "required evidence with invalid producer: "+req.ID)
		add(seenReq[req.ID], InvalidRecord, "duplicate required evidence id: "+req.ID)
		seenReq[req.ID] = true
	}

	// Output bounds (same rules the seal path enforces).
	add(len(r.Matter.Data) > maxMatterBytes, InvalidRecord, "matter exceeds bound")
	add(len(r.Review.Notes) > maxReviewNotesBytes, InvalidRecord, "review notes exceed bound")
	total := 0
	for _, side := range []StationRecord{r.Alpha, r.Beta} {
		add(len(side.Artifacts) > maxArtifacts, InvalidRecord, "too many artifacts")
		seen := make(map[string]bool)
		for _, a := range side.Artifacts {
			add(a.ID == "" || a.Kind == "", InvalidRecord, "artifact with empty id/kind")
			add(a.Encoding != "utf8", InvalidRecord, "artifact with unknown encoding: "+a.ID)
			add(!utf8.ValidString(a.Text), InvalidRecord, "artifact is not valid UTF-8: "+a.ID)
			add(len(a.Text) > maxArtifactBytes, InvalidRecord, "artifact exceeds bound: "+a.ID)
			add(seen[a.ID], InvalidRecord, "duplicate artifact id: "+a.ID)
			seen[a.ID] = true
			total += len(a.Text)
		}
	}
	add(total > maxAggregateArtifact, InvalidRecord, "aggregate artifact bytes exceed bound")

	return fs
}

// validate (V) checks the receipt at the scope-lift boundary: the one digest,
// the full record boundary (validateRecord), and contract satisfaction — with
// typed failures. Producer authority is positional and fails closed: a
// required artifact with an invalid producer never resolves to a side.
func validate(rc Receipt) Verdict {
	r := rc.Record
	fs := validateRecord(r)
	if sha256hex(r.canonicalBytes()) != rc.ScopeLiftDigest {
		fs = append(fs, Failure{InvalidRecord, "scope-lift digest does not recompute"})
	}

	for _, req := range r.Contract.RequiredEvidence {
		var side []Artifact
		switch req.Producer {
		case RoleAlpha:
			side = r.Alpha.Artifacts
		case RoleBeta:
			side = r.Beta.Artifacts
		default:
			continue // invalid producer already an integrity failure above
		}
		if !hasArtifact(side, req.ID, req.Kind) {
			fs = append(fs, Failure{ContractUnmet, fmt.Sprintf("missing required %s artifact: %s", req.Producer, req.ID)})
		}
	}
	if !r.Review.Pass {
		fs = append(fs, Failure{ContractUnmet, "review.pass=false"})
	}

	sort.Slice(fs, func(i, j int) bool {
		if fs[i].Class != fs[j].Class {
			return fs[i].Class < fs[j].Class
		}
		return fs[i].Detail < fs[j].Detail
	})
	return Verdict{Pass: len(fs) == 0, Failures: fs}
}

func hasArtifact(arts []Artifact, id, kind string) bool {
	for _, a := range arts {
		if a.ID == id && a.Kind == kind {
			return true
		}
	}
	return false
}

// decide (δ) routes by failure class: integrity fails closed; contract-unmet
// repairs; pass accepts (pure).
func decide(v Verdict) Decision {
	switch {
	case v.Pass:
		return Accept
	case v.hasIntegrityFailure():
		return Reject
	default:
		return RepairDispatch
	}
}

// lift maps (verdict, decision, mode) to the terminal status. `simulated` is
// admissible ONLY for an otherwise coherent successful stub smoke run (Pi
// PR-#718-fido β D4); a stub with integrity or semantic failures routes by the
// normal table — failures are never masked. An inconsistent pair is an error,
// never a closure.
func lift(v Verdict, d Decision, mode ExecutionMode) (Status, error) {
	if mode == ModeStub && v.Pass && d == Accept {
		return Simulated, nil
	}
	switch {
	case v.Pass && (d == Accept || d == Release):
		return Accepted, nil
	case !v.Pass && d == Override:
		return Degraded, nil
	case d == Reject:
		return Rejected, nil
	case d == RepairDispatch:
		return NeedsRepair, nil
	default:
		return "", fmt.Errorf("%w: verdict.pass=%v decision=%q", ErrInvalidClosure, v.Pass, d)
	}
}

// repairFrom derives the repair surface purely from (verdict, status) — the
// verifier recomputes and compares it, so repair is never a second
// unauthenticated terminal authority (Pi PR-#718-fido β D2).
func repairFrom(v Verdict, st Status) *RepairRequest {
	if st != NeedsRepair {
		return nil
	}
	return &RepairRequest{Reason: "contract unmet", Failed: failureDetails(v)}
}

func failureDetails(v Verdict) []string {
	out := make([]string, 0, len(v.Failures))
	for _, f := range v.Failures {
		out = append(out, string(f.Class)+": "+f.Detail)
	}
	return out
}

// --- Kernel-boundary validation -------------------------------------------

func validateSpec(s Spec) error {
	if seatIsNil(s.Alpha) {
		return errors.New("cellkernel: spec has nil alpha")
	}
	if seatIsNil(s.Beta) {
		return errors.New("cellkernel: spec has nil beta")
	}
	if s.Contract.ID == "" {
		return errors.New("cellkernel: contract.id is empty")
	}
	if len(s.Contract.RequiredEvidence) > maxRequiredEvidence {
		return fmt.Errorf("cellkernel: too many required evidence refs (%d > %d)", len(s.Contract.RequiredEvidence), maxRequiredEvidence)
	}
	seen := make(map[string]bool)
	for _, r := range s.Contract.RequiredEvidence {
		if r.ID == "" || r.Kind == "" {
			return errors.New("cellkernel: required evidence needs non-empty id and kind")
		}
		if r.Producer != RoleAlpha && r.Producer != RoleBeta {
			return fmt.Errorf("cellkernel: required evidence %q has invalid producer %q", r.ID, r.Producer)
		}
		if seen[r.ID] {
			return fmt.Errorf("cellkernel: duplicate required evidence id %q", r.ID)
		}
		seen[r.ID] = true
	}
	return nil
}

// --- helpers ---------------------------------------------------------------

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func seatIsNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Func, reflect.Map, reflect.Slice, reflect.Chan:
		return rv.IsNil()
	default:
		return false
	}
}
