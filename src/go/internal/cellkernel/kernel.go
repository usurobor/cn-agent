// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) single-episode kernel: the substrate-independent
// five-step closure a coherence cell runs at one scope.
//
//  1. matter   := α.produce(contract)
//  2. review   := β.review(betaInput)
//  3. receipt  := γ.close(record)                               [kernel-owned]
//  4. verdict  := V(record, receipt)                            [kernel-owned]
//  5. decision := δ.decide(receipt, verdict)                    [kernel-owned]
//
// See docs/architecture/CDS-CELL-MIGRATION.md and cnos#711/#717. Single-episode
// engine: no repair loop, no composition, no protocol dispatch.
//
// Authority model (Pi β #31–#33 + PR-#718 β). The terminal object is an
// Envelope, and VerifyEnvelope re-derives EVERY field of it from content — a
// parent trusts nothing that is merely copied in:
//
//   - WHOLE-ENVELOPE PROOF (D1): the emitted Envelope (schema, protocol_validated,
//     execution_mode, verdict, decision, status, repair, resolved_spec, receipt)
//     is one verified object. verdict←V(receipt), decision←δ(verdict),
//     status←(decision, execution_mode); protocol_validated is pinned false.
//   - REPRODUCIBLE BINDINGS (D2): the Envelope carries the normalized
//     resolved_spec (version/protocol/profile/params/skills/contract) so
//     resolved_spec_hash recomputes; beta_input_hash recomputes from the
//     receipt; execution ids are distinct and every evidence ref is bound to its
//     producer's station id.
//   - TYPED FAILURE ROUTING (D3): V classifies failures. Only contract_unmet may
//     become needs_repair; integrity failures (invalid_receipt/_evidence/
//     _identity/_independence) fail closed to rejected — never the α repair path.
//   - FAIL-CLOSED IDENTITY (D4): identities are minted through one error-returning
//     op and must be non-empty and pairwise distinct before α runs.
//   - NON-AUTHORITATIVE SMOKE (D5): a stub run yields status `simulated`, never
//     ordinary accepted authority.
//   - EVIDENCE BYTES (C1): evidence bytes must be valid UTF-8 and bounded per-item
//     and in aggregate.
//
// γ, V, δ are kernel-owned and mechanical; only α and β are open seats.
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

// Schema ids.
const (
	EnvelopeSchema        = "cnos.cellkernel.episode-envelope.v0"
	betaInputCanonVersion = "cnos.cellkernel.beta-input-canon.v0"
	resolvedSpecCanon     = "cnos.cellkernel.resolved-spec-canon.v0"
	BetaInputPolicyID     = "cnos.cellkernel.beta-input.v0"
)

// Output bounds enforced at the kernel boundary (D6/C1).
const (
	maxRequiredEvidence  = 64
	maxMatterBytes       = 1 << 20 // 1 MiB
	maxReviewNotesBytes  = 64 << 10
	maxEvidenceItems     = 64
	maxEvidenceBytes     = 1 << 20 // per item
	maxAggregateEvidence = 4 << 20 // sum of all evidence bytes
)

// Role names which seat produced an artifact. The runtime assigns it.
type Role string

const (
	RoleAlpha Role = "alpha"
	RoleBeta  Role = "beta"
)

// ExecutionMode reports how authoritative a run is. A stub run is a smoke test
// (non-authoritative); a mechanical run has β independently verify.
type ExecutionMode string

const (
	ModeStub       ExecutionMode = "stub"
	ModeMechanical ExecutionMode = "mechanical"
)

func knownMode(m ExecutionMode) bool { return m == ModeStub || m == ModeMechanical }

// --- Artifacts ----------------------------------------------------------

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

// EvidenceCandidate is what a seat returns: semantic identity + bytes only.
type EvidenceCandidate struct {
	ID    string
	Kind  string
	Bytes []byte
}

// EvidenceRef is an authenticated, self-verifying evidence record. Content is
// the inlined UTF-8 bytes the runtime hashed; Ref is the content address.
type EvidenceRef struct {
	ID                  string `json:"id"`
	Kind                string `json:"kind"`
	Producer            Role   `json:"producer"`
	ProducerExecutionID string `json:"producer_execution_id"`
	Ref                 string `json:"ref"`
	SHA256              string `json:"sha256"`
	Content             string `json:"content"`
}

type AlphaResult struct {
	Matter   Matter
	Evidence []EvidenceCandidate
}

type BetaResult struct {
	Review   Review
	Evidence []EvidenceCandidate
}

// BetaInput is the runtime-owned review surface handed to β.
type BetaInput struct {
	Contract      Contract
	ContractHash  string
	Matter        Matter
	AlphaEvidence []EvidenceRef
	PolicyID      string
	BundleHash    string
}

// Receipt is the kernel's closure record (no resolved-spec identity; that lives
// in the Envelope's ResolvedSpec).
type Receipt struct {
	EpisodeID        string        `json:"episode_id"`
	Contract         Contract      `json:"contract"`
	ContractHash     string        `json:"contract_hash"`
	Matter           Matter        `json:"matter"`
	MatterHash       string        `json:"matter_hash"`
	Review           Review        `json:"review"`
	ReviewHash       string        `json:"review_hash"`
	Evidence         []EvidenceRef `json:"evidence_refs"`
	EvidenceHash     string        `json:"evidence_hash"`
	AlphaExecutionID string        `json:"alpha_execution_id"`
	BetaExecutionID  string        `json:"beta_execution_id"`
	PolicyID         string        `json:"beta_input_policy_id"`
	BetaInputHash    string        `json:"beta_input_hash"`
}

// FailureClass separates repairable contract-unmet from fail-closed integrity
// failures (D3).
type FailureClass string

const (
	ContractUnmet       FailureClass = "contract_unmet"
	InvalidReceipt      FailureClass = "invalid_receipt"
	InvalidEvidence     FailureClass = "invalid_evidence"
	InvalidIdentity     FailureClass = "invalid_identity"
	InvalidIndependence FailureClass = "invalid_independence"
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
	Simulated   Status = "simulated" // stub profile: non-authoritative
)

type RepairRequest struct {
	Reason string   `json:"reason"`
	Failed []string `json:"failed,omitempty"`
}

// ResolvedSpec is the normalized executable spec resolved_spec_hash covers.
type ResolvedSpec struct {
	Canon            string            `json:"canon"`
	Version          string            `json:"version"`
	DeclaredProtocol string            `json:"declared_protocol"`
	Profile          string            `json:"profile"`
	Params           map[string]string `json:"params,omitempty"`
	AlphaSkills      []string          `json:"alpha_skills"`
	BetaSkills       []string          `json:"beta_skills"`
	Contract         Contract          `json:"contract"`
}

func (r ResolvedSpec) hash() string { return hashJSON(r) }

// Envelope is the terminal, whole-object-verifiable artifact `cn cell run`
// emits. VerifyEnvelope re-derives every field.
type Envelope struct {
	Schema            string         `json:"envelope_schema"`
	ProtocolValidated bool           `json:"protocol_validated"`
	ExecutionMode     ExecutionMode  `json:"execution_mode"`
	Status            Status         `json:"status"`
	Decision          Decision       `json:"decision"`
	Verdict           Verdict        `json:"verdict"`
	ResolvedSpec      ResolvedSpec   `json:"resolved_spec"`
	ResolvedSpecHash  string         `json:"resolved_spec_hash"`
	Receipt           Receipt        `json:"receipt"`
	Repair            *RepairRequest `json:"repair,omitempty"`
}

// --- The two open seats -------------------------------------------------

type Alpha interface {
	Produce(ctx context.Context, c Contract) (AlphaResult, error)
}

type Beta interface {
	Review(ctx context.Context, in BetaInput) (BetaResult, error)
}

type Spec struct {
	Contract Contract
	Alpha    Alpha
	Beta     Beta
}

// --- Identity source (D1/D4) --------------------------------------------

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

// IDSource mints the whole identity tuple; it may fail (crypto/rand).
type IDSource interface {
	Mint() (Identity, error)
}

type randomIDs struct{}

func (randomIDs) Mint() (Identity, error) {
	e, err := randHex(16)
	if err != nil {
		return Identity{}, err
	}
	a, err := randHex(16)
	if err != nil {
		return Identity{}, err
	}
	b, err := randHex(16)
	if err != nil {
		return Identity{}, err
	}
	return Identity{Episode: "ep-" + e, Alpha: "alpha-" + a, Beta: "beta-" + b}, nil
}

func randHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("mint id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// RunMeta carries the resolved-spec identity + execution mode from the runner.
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

// RunEpisode executes the CCNF five-step closure and returns the terminal
// Envelope. A returned error is a malfunction (a seat failed, the spec was
// invalid, identity minting failed, or the closure was internally inconsistent);
// otherwise the Envelope's Status reports how the episode closed.
func RunEpisode(ctx context.Context, s Spec, opts ...RunOption) (Envelope, error) {
	cfg := runConfig{ids: randomIDs{}, meta: RunMeta{ExecutionMode: ModeMechanical}}
	for _, o := range opts {
		o(&cfg)
	}
	if !knownMode(cfg.meta.ExecutionMode) {
		return Envelope{}, fmt.Errorf("cellkernel: unknown execution mode %q", cfg.meta.ExecutionMode)
	}
	if err := validateSpec(s); err != nil {
		return Envelope{}, err
	}
	if err := ctx.Err(); err != nil {
		return Envelope{}, fmt.Errorf("cellkernel: context: %w", err)
	}

	id, err := cfg.ids.Mint() // D4: fail-closed identity before α runs.
	if err != nil {
		return Envelope{}, fmt.Errorf("cellkernel: %w", err)
	}
	if err := id.valid(); err != nil {
		return Envelope{}, fmt.Errorf("cellkernel: %w", err)
	}

	frozen := s.Contract.clone()
	contractHash := hashJSON(frozen)

	aRes, err := s.Alpha.Produce(ctx, frozen.clone()) // 1
	if err != nil {
		return Envelope{}, fmt.Errorf("alpha produce: %w", err)
	}
	if err := boundMatter(aRes.Matter); err != nil {
		return Envelope{}, err
	}
	alphaEv, err := authenticate(aRes.Evidence, RoleAlpha, id.Alpha)
	if err != nil {
		return Envelope{}, fmt.Errorf("alpha evidence: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return Envelope{}, fmt.Errorf("cellkernel: context after alpha: %w", err)
	}

	bin := BetaInput{
		Contract:      frozen.clone(),
		ContractHash:  contractHash,
		Matter:        aRes.Matter,
		AlphaEvidence: cloneEvidence(alphaEv),
		PolicyID:      BetaInputPolicyID,
	}
	betaInputHash := hashBetaInput(bin)
	bin.BundleHash = betaInputHash

	bRes, err := s.Beta.Review(ctx, bin) // 2
	if err != nil {
		return Envelope{}, fmt.Errorf("beta review: %w", err)
	}
	if err := boundReview(bRes.Review); err != nil {
		return Envelope{}, err
	}
	betaEv, err := mintBetaEvidence(bRes, id.Beta)
	if err != nil {
		return Envelope{}, fmt.Errorf("beta evidence: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return Envelope{}, fmt.Errorf("cellkernel: context after beta: %w", err)
	}

	evidence := append(cloneEvidence(alphaEv), betaEv...)
	if err := boundAggregateEvidence(evidence); err != nil {
		return Envelope{}, err
	}

	receipt := Receipt{ // 3 (kernel γ)
		EpisodeID:        id.Episode,
		Contract:         frozen,
		ContractHash:     contractHash,
		Matter:           aRes.Matter,
		MatterHash:       hashJSON(aRes.Matter),
		Review:           bRes.Review,
		ReviewHash:       hashJSON(bRes.Review),
		Evidence:         evidence,
		EvidenceHash:     hashJSON(evidence),
		AlphaExecutionID: id.Alpha,
		BetaExecutionID:  id.Beta,
		PolicyID:         BetaInputPolicyID,
		BetaInputHash:    betaInputHash,
	}

	verdict := validate(receipt) // 4 (kernel V)
	decision := decide(verdict)  // 5 (kernel δ)

	rs := cfg.meta.ResolvedSpec
	rs.Canon = resolvedSpecCanon
	rs.Contract = frozen // bind the frozen contract into the resolved spec

	status, err := statusOf(verdict, decision, cfg.meta.ExecutionMode)
	if err != nil {
		return Envelope{}, err
	}

	env := Envelope{
		Schema:            EnvelopeSchema,
		ProtocolValidated: false,
		ExecutionMode:     cfg.meta.ExecutionMode,
		Status:            status,
		Decision:          decision,
		Verdict:           verdict,
		ResolvedSpec:      rs,
		ResolvedSpecHash:  rs.hash(),
		Receipt:           receipt,
	}
	if status == NeedsRepair {
		env.Repair = &RepairRequest{Reason: "contract unmet", Failed: failureDetails(verdict)}
	}
	return env, nil
}

func failureDetails(v Verdict) []string {
	out := make([]string, 0, len(v.Failures))
	for _, f := range v.Failures {
		out = append(out, string(f.Class)+": "+f.Detail)
	}
	return out
}

// statusOf maps (verdict, decision, mode) to a terminal Status. A stub run is
// always non-authoritative `simulated` (D5). An inconsistent pair is an error.
func statusOf(v Verdict, d Decision, mode ExecutionMode) (Status, error) {
	if mode == ModeStub {
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

// --- D6: kernel-boundary spec validation + output bounds ----------------

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

func boundMatter(m Matter) error {
	if len(m.Data) > maxMatterBytes {
		return fmt.Errorf("cellkernel: matter exceeds %d bytes", maxMatterBytes)
	}
	return nil
}

func boundReview(r Review) error {
	if len(r.Notes) > maxReviewNotesBytes {
		return fmt.Errorf("cellkernel: review notes exceed %d bytes", maxReviewNotesBytes)
	}
	return nil
}

func boundAggregateEvidence(refs []EvidenceRef) error {
	total := 0
	for _, e := range refs {
		total += len(e.Content)
	}
	if total > maxAggregateEvidence {
		return fmt.Errorf("cellkernel: aggregate evidence exceeds %d bytes", maxAggregateEvidence)
	}
	return nil
}

// --- D2/C1: runtime evidence authentication -----------------------------

func authenticate(candidates []EvidenceCandidate, role Role, execID string) ([]EvidenceRef, error) {
	if len(candidates) > maxEvidenceItems {
		return nil, fmt.Errorf("too many evidence items (%d > %d)", len(candidates), maxEvidenceItems)
	}
	out := make([]EvidenceRef, 0, len(candidates))
	for _, c := range candidates {
		if c.ID == "" || c.Kind == "" {
			return nil, errors.New("evidence candidate needs non-empty id and kind")
		}
		if len(c.Bytes) > maxEvidenceBytes {
			return nil, fmt.Errorf("evidence %q exceeds %d bytes", c.ID, maxEvidenceBytes)
		}
		if !utf8.Valid(c.Bytes) { // C1: bytes must round-trip through JSON.
			return nil, fmt.Errorf("evidence %q is not valid UTF-8", c.ID)
		}
		sum := sha256hex(c.Bytes)
		out = append(out, EvidenceRef{
			ID:                  c.ID,
			Kind:                c.Kind,
			Producer:            role,
			ProducerExecutionID: execID,
			Ref:                 "sha256:" + sum,
			SHA256:              sum,
			Content:             string(c.Bytes),
		})
	}
	return out, nil
}

// mintBetaEvidence authenticates β's candidates, then mints the canonical
// beta_review from the ACTUAL review (a seat cannot substitute it).
func mintBetaEvidence(b BetaResult, execID string) ([]EvidenceRef, error) {
	filtered := make([]EvidenceCandidate, 0, len(b.Evidence))
	for _, c := range b.Evidence {
		if c.ID == "beta_review" {
			continue
		}
		filtered = append(filtered, c)
	}
	reviewBytes, _ := json.Marshal(b.Review)
	filtered = append([]EvidenceCandidate{{ID: "beta_review", Kind: "review", Bytes: reviewBytes}}, filtered...)
	return authenticate(filtered, RoleBeta, execID)
}

func cloneEvidence(in []EvidenceRef) []EvidenceRef {
	return append([]EvidenceRef(nil), in...)
}

func hashBetaInput(in BetaInput) string {
	type evCanon struct{ ID, Kind, Producer, Exec, SHA string }
	type canon struct {
		Version      string
		PolicyID     string
		ContractHash string
		Matter       string
		Alpha        []evCanon
	}
	c := canon{Version: betaInputCanonVersion, PolicyID: in.PolicyID, ContractHash: in.ContractHash, Matter: in.Matter.Data}
	for _, e := range in.AlphaEvidence {
		c.Alpha = append(c.Alpha, evCanon{e.ID, e.Kind, string(e.Producer), e.ProducerExecutionID, e.SHA256})
	}
	return hashJSON(c)
}

// --- Kernel-owned γ / V / δ ----------------------------------------------

// validate (V) recomputes every binding from the receipt's own content and
// classifies each failure. It does not re-judge β's prose.
func validate(rc Receipt) Verdict {
	var fs []Failure
	add := func(cond bool, class FailureClass, msg string) {
		if cond {
			fs = append(fs, Failure{Class: class, Detail: msg})
		}
	}

	add(hashJSON(rc.Contract) != rc.ContractHash, InvalidReceipt, "contract hash mismatch")
	add(hashJSON(rc.Matter) != rc.MatterHash, InvalidReceipt, "matter hash mismatch")
	add(hashJSON(rc.Review) != rc.ReviewHash, InvalidReceipt, "review hash mismatch")
	add(hashJSON(rc.Evidence) != rc.EvidenceHash, InvalidReceipt, "evidence hash mismatch")

	// Identity: non-empty, distinct, and each evidence ref bound to its station.
	add(rc.AlphaExecutionID == "" || rc.BetaExecutionID == "", InvalidIdentity, "missing execution id")
	add(rc.AlphaExecutionID == rc.BetaExecutionID, InvalidIdentity, "alpha and beta share an execution id")

	seen := make(map[string]int)
	for _, e := range rc.Evidence {
		seen[e.ID]++
		if sha256hex([]byte(e.Content)) != e.SHA256 {
			fs = append(fs, Failure{InvalidEvidence, "evidence hash mismatch: " + e.ID})
		}
		if e.Ref != "sha256:"+e.SHA256 {
			fs = append(fs, Failure{InvalidEvidence, "evidence ref not content-addressed: " + e.ID})
		}
		wantExec := rc.AlphaExecutionID
		if e.Producer == RoleBeta {
			wantExec = rc.BetaExecutionID
		}
		if e.Producer != RoleAlpha && e.Producer != RoleBeta {
			fs = append(fs, Failure{InvalidEvidence, "evidence has no producer: " + e.ID})
		} else if e.ProducerExecutionID != wantExec {
			fs = append(fs, Failure{InvalidIndependence, "evidence producer-execution not bound to station: " + e.ID})
		}
	}
	for id, n := range seen {
		if n > 1 {
			fs = append(fs, Failure{InvalidEvidence, "duplicate evidence id: " + id})
		}
	}

	// Required presence + producer authority (a work gap, not an integrity fault).
	for _, req := range rc.Contract.RequiredEvidence {
		if !hasAuthorizedEvidence(rc.Evidence, req) {
			fs = append(fs, Failure{ContractUnmet, fmt.Sprintf("missing/unauthorized required evidence: %s (want producer %s)", req.ID, req.Producer)})
		}
	}
	if !rc.Review.Pass {
		fs = append(fs, Failure{ContractUnmet, "review.pass=false"})
	}

	// Deterministic order so the verdict hashes/compares stably (map iteration
	// over duplicate ids above is otherwise unordered).
	sort.Slice(fs, func(i, j int) bool {
		if fs[i].Class != fs[j].Class {
			return fs[i].Class < fs[j].Class
		}
		return fs[i].Detail < fs[j].Detail
	})
	return Verdict{Pass: len(fs) == 0, Failures: fs}
}

func hasAuthorizedEvidence(refs []EvidenceRef, req RequiredRef) bool {
	for _, e := range refs {
		if e.ID == req.ID && e.Kind == req.Kind && e.Producer == req.Producer {
			return true
		}
	}
	return false
}

// decide (δ) routes by failure class (D3): an integrity failure fails closed to
// reject; a pure contract-unmet goes to repair; a pass accepts.
func decide(v Verdict) Decision {
	if v.Pass {
		return Accept
	}
	if v.hasIntegrityFailure() {
		return Reject
	}
	return RepairDispatch
}

// --- helpers ------------------------------------------------------------

func hashJSON(v any) string {
	b, _ := json.Marshal(v)
	return sha256hex(b)
}

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
