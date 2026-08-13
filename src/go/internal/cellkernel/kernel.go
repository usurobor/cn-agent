// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) single-episode kernel, implemented under the
// operator-ratified FIDO/functional doctrine
// (msg-cn-pi-cnos-cell-runner-fido-functional-44):
//
//	alphaIn  := AlphaInput{frozen contract}
//	aOut     := α(alphaIn)                        Result<AlphaOutput, error>
//	sealedA  := sealAlpha(aOut)                   runtime-owned, immutable
//	betaIn   := BetaInput{contract, matter}
//	bOut     := β(betaIn)                         Result<BetaOutput, error>
//	sealedB  := sealBeta(bOut)
//	record   := compose(start, sealedA, sealedB)  one immutable EpisodeRecord
//	receipt  := γ(record)                         canonical bytes + ONE digest
//	verdict  := V(contract, receipt)              typed failures
//	decision := δ(receipt, verdict)
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
// package), β receives a fresh projection of the sealed α MATTER only (never
// α's artifacts — those are γ/V's evidence channel), and the
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
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"
)

// Schema / canon identifiers.
const (
	ClosureSchema = "cnos.cellkernel.episode-closure.v0"
	RecordCanon   = "cnos.cellkernel.episode-record-canon.v0"
)

// Bounds enforced at the kernel boundary.
const (
	maxRequiredEvidence  = 64
	maxMatterBytes       = 1 << 20 // 1 MiB
	maxReviewNotesBytes  = 64 << 10
	maxArtifacts         = 64
	maxArtifactBytes     = 1 << 20 // per artifact
	maxAggregateArtifact = 4 << 20 // sum over both seats
	maxAssessmentUnits   = 256
	// maxAssessmentBytes bounds the assessment's text — unit ids, reasons and
	// citations summed. One aggregate rather than a bound per field: what must
	// not grow without limit is the value that travels, and three separate
	// numbers would be three things to justify for one property.
	maxAssessmentBytes = 64 << 10
)

// MaxOpaqueSlotBytes is the one bound for every opaque contract slot. A
// per-slot bound would be the kernel deciding that one opaque payload deserves
// more room than another, which is a judgement about what they MEAN — exactly
// what this boundary refuses to make.
//
// It is EXPORTED because a profile's admission gate must refuse what this
// boundary will refuse, and the only way for two gates to agree on a number is
// for there to be one number. A gate that restated it would drift silently:
// nothing fails when two constants disagree, a document simply passes the door
// and then breaks the episode.
const MaxOpaqueSlotBytes = 64 << 10

// Role names which station a required artifact must come from. The check is
// positional (which side of the record the artifact sits on), never a stamp.
type Role string

const (
	RoleAlpha Role = "alpha"
	RoleBeta  Role = "beta"
)

// ExecutionMode is how the episode's work was produced. It is invocation
// truth: the parent supplies it and VerifyClosure binds it, so a closure
// cannot misreport how it ran.
//
//	stub       — nothing real was produced; non-authoritative (`simulated`)
//	mechanical — deterministic seats; reproducible from the record
//	cognitive  — a provider (a model) held at least one seat; authoritative
//	             work, but NOT reproducible by re-running
type ExecutionMode string

const (
	ModeStub       ExecutionMode = "stub"
	ModeMechanical ExecutionMode = "mechanical"
	ModeCognitive  ExecutionMode = "cognitive"
)

func knownMode(m ExecutionMode) bool {
	return m == ModeStub || m == ModeMechanical || m == ModeCognitive
}

// --- Immutable values ----------------------------------------------------

type RequiredRef struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Producer Role   `json:"producer"`
}

type Contract struct {
	ID   string `json:"id"`
	Goal string `json:"goal"`
	// Issue and Design are the two halves of the admitted contract, opaque
	// here. The kernel never learns what either MEANS — that belongs to
	// whichever protocol authored them — so the only rules at this boundary
	// are structural: valid JSON, within MaxOpaqueSlotBytes. They are two
	// slots rather than one because they are logically distinct documents:
	// merging them here would let the record lose the distinction that the
	// problem statement and the proposed change are separately authored and
	// separately admitted.
	//
	// Because EpisodeRecord carries the frozen Contract, both are inside
	// canonicalBytes() and therefore inside the one scope-lift digest; there
	// is no second digest to bind them, and adding one would be a second proof
	// surface for the same bytes.
	Issue  json.RawMessage `json:"issue,omitempty"`
	Design json.RawMessage `json:"design,omitempty"`
	// Subject is the thing the episode acts on, opaque here under exactly the
	// rules Issue and Design obey, and for the same reason: what a subject IS
	// belongs to the adapter that pins it, never to this boundary.
	//
	// It is a CONTRACT value rather than a seat argument because both stations
	// need the same one. Frozen once, it is the same bytes on both sides by
	// construction — there is no second place to state it and therefore
	// nothing to keep in step.
	Subject          json.RawMessage `json:"subject,omitempty"`
	RequiredEvidence []RequiredRef   `json:"required_evidence,omitempty"`
}

// opaqueSlots is the complete list of the contract's opaque slots, paired with
// the name a diagnostic uses. Written once so a slot added to the struct and
// forgotten here is the only way to escape the boundary rules — rather than
// three near-identical checks in each of the three places that apply them.
func (c Contract) opaqueSlots() []struct {
	name string
	raw  json.RawMessage
} {
	return []struct {
		name string
		raw  json.RawMessage
	}{
		{"issue", c.Issue},
		{"design", c.Design},
		{"subject", c.Subject},
	}
}

func (c Contract) clone() Contract {
	cp := c
	if c.RequiredEvidence != nil {
		cp.RequiredEvidence = append([]RequiredRef(nil), c.RequiredEvidence...)
	}
	// A shared slice is a mutable value inside a struct otherwise frozen by
	// copy: without this, a seat handed AlphaInput.Contract could write through
	// an opaque slot into the contract the runtime composes the record from.
	cp.Issue = cloneRaw(c.Issue)
	cp.Design = cloneRaw(c.Design)
	cp.Subject = cloneRaw(c.Subject)
	return cp
}

// cloneRaw copies a raw slot, preserving the nil/empty distinction: a nil slot
// is absent and `omitempty` keeps it out of the canonical bytes entirely, so
// turning nil into an empty non-nil slice here would change no JSON today and
// would be a trap the first time a caller tested for absence.
func cloneRaw(raw json.RawMessage) json.RawMessage {
	if raw == nil {
		return nil
	}
	return append(json.RawMessage{}, raw...)
}

type Matter struct {
	Data string `json:"data"`
}

// Disposition is what an assessing seat decided about one obligation. The set
// is closed and three-valued, and the third value is the point: `unverified`
// says the seat could not decide, which is different from deciding against.
// Collapsing it into `finding` would report a judgement nobody made, and
// collapsing it into `pass` would report evidence nobody saw.
type Disposition string

const (
	DispositionPass       Disposition = "pass"
	DispositionFinding    Disposition = "finding"
	DispositionUnverified Disposition = "unverified"
)

func knownDisposition(d Disposition) bool {
	return d == DispositionPass || d == DispositionFinding || d == DispositionUnverified
}

// UnitResult is one obligation's disposition. `Unit` names the obligation in
// whatever catalogue vocabulary the profile that built it uses — the kernel
// reads the id as an opaque label and never learns what it means, exactly as
// it does with the contract's opaque slots.
//
// WHAT THE KERNEL DOES KNOW is the shape: a disposition from the closed set,
// and a reason whenever the disposition is not `pass`. That much has to live
// here rather than with the profile, because V re-derives the verdict from
// these values and a verifier must be able to make the same derivation from
// the serialized record alone.
type UnitResult struct {
	Unit        string      `json:"unit"`
	Disposition Disposition `json:"disposition"`
	Reason      string      `json:"reason,omitempty"`
	Citations   []string    `json:"citations,omitempty"`
}

// Review is the assessing seat's product.
//
// Assessment is OPTIONAL and omitted when empty, which is deliberate rather
// than lenient: the cells that predate assessment — the bool checker, the
// stub, the mechanical-unmet reviewer — state a verdict and no per-obligation
// coverage, and forcing an empty list into their records would be a claim
// that a review with zero units happened. What the kernel refuses is an
// assessment that disagrees with itself, and V refuses one that disagrees
// with `Pass`: every non-`pass` unit is its own contract-unmet failure, so a
// seat cannot report findings and a passing verdict in the same breath.
//
// COVERAGE AGAINST A CATALOGUE IS NOT CHECKED HERE, and cannot be: the
// catalogue is derived from the contract's opaque issue by the profile that
// admitted it, and a kernel that could tell whether coverage was exact would
// have had to learn that profile's language. The assessing fill enforces
// exactness before it returns; what the record proves at this boundary is
// that the units it carries are well formed and that the verdict follows from
// them.
type Review struct {
	Pass       bool         `json:"pass"`
	Notes      string       `json:"notes"`
	Assessment []UnitResult `json:"assessment,omitempty"`
}

// clone deep-copies the assessment so a seat retaining its own slice cannot
// write through the sealed record — the same rule Contract.clone applies to
// the opaque slots, for the same reason.
func (r Review) clone() Review {
	cp := r
	if r.Assessment != nil {
		cp.Assessment = make([]UnitResult, 0, len(r.Assessment))
		for _, u := range r.Assessment {
			u.Citations = append([]string(nil), u.Citations...)
			cp.Assessment = append(cp.Assessment, u)
		}
	}
	return cp
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

// BetaInput is the runtime-owned review surface — exactly the CCNF signature
// β.review(contract, matter): a fresh frozen contract copy and the sealed
// matter. Artifacts/evidence are γ/V's channel, never β's (Pi round-6 D1); a
// fill needing richer review matter widens Matter. The input shape is
// defined by this type and the closure schema version — there is no
// self-reported policy id (Pi round-7 C1).
type BetaInput struct {
	Contract Contract
	Matter   Matter
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

func (s SealedAlpha) projection() Matter {
	return s.matter
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
// lift digest covers it and any verifier can reproduce it. Alpha and Beta are
// the COMPLETE resolved seat declarations — each a tagged object whose "fill"
// selected its constructor and whose remaining fields were that constructor's
// arguments (fill-owned construction, Pi cds-fill-construction-51). The kernel
// treats them as opaque canonical bytes: it never learns what a fill means.
type ResolvedSpec struct {
	Version          string          `json:"version"`
	DeclaredProtocol string          `json:"declared_protocol"`
	Alpha            json.RawMessage `json:"alpha"`
	Beta             json.RawMessage `json:"beta"`
}

// clone deep-copies the resolved spec so no caller- or seat-retained alias can
// mutate runtime-owned invocation truth (Pi PR-#718-fido β D3).
func (r ResolvedSpec) clone() ResolvedSpec {
	cp := r
	cp.Alpha = append(json.RawMessage(nil), r.Alpha...)
	cp.Beta = append(json.RawMessage(nil), r.Beta...)
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
	Canon        string        `json:"canon"`
	EpisodeID    string        `json:"episode_id"`
	Mode         ExecutionMode `json:"execution_mode"`
	ResolvedSpec ResolvedSpec  `json:"resolved_spec"`
	Contract     Contract      `json:"contract"`
	Alpha        StationRecord `json:"alpha"`
	Matter       Matter        `json:"matter"`
	Beta         StationRecord `json:"beta"`
	Review       Review        `json:"review"`
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
	ids IDSource
}

type RunOption func(*runConfig)

func WithIDSource(s IDSource) RunOption { return func(c *runConfig) { c.ids = s } }

var ErrInvalidClosure = errors.New("cellkernel: inconsistent (verdict, decision) pair")

// RunEpisode executes the five-step closure as a pure-shaped pipeline over
// immutable values. meta is required — an episode that cannot state its
// resolved invocation truth is a malfunction, not a closable cell (Pi round-6
// B1). A returned error is a malfunction (invalid spec/meta, identity failure,
// seat error, bound violation); otherwise the Closure's Status reports how the
// episode closed.
func RunEpisode(ctx context.Context, s Spec, meta RunMeta, opts ...RunOption) (Closure, error) {
	cfg := runConfig{ids: randomIDs{}}
	for _, o := range opts {
		o(&cfg)
	}
	if err := validateSpec(s); err != nil {
		return Closure{}, err
	}
	if err := validateMeta(meta); err != nil { // D1: invocation truth is part of the boundary.
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
	frozenMeta := RunMeta{ExecutionMode: meta.ExecutionMode, ResolvedSpec: meta.ResolvedSpec.clone()}

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

	// Station β: fresh projection of sealed α matter — never α's live scope.
	bOut, err := s.Beta.Review(ctx, BetaInput{
		Contract: frozen.clone(),
		Matter:   sealedA.projection(),
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
	verdict := validate(frozen, receipt)
	decision := decide(receipt, verdict)
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
	if err := assessmentIntegrity(o.Review.Assessment); err != nil {
		return SealedBeta{}, err
	}
	arts, err := normalizeArtifacts(o.Artifacts)
	if err != nil {
		return SealedBeta{}, err
	}
	return SealedBeta{exec: exec, review: o.Review.clone(), artifacts: arts}, nil
}

// assessmentIntegrity is the COMPLETE set of shape rules the kernel applies to
// an assessment, written once so the seal path and validateRecord cannot
// diverge: a record that self-verifies must be one the honest path would have
// produced (eng/go §2.17, one parser per fact).
//
// Every rule is about whether the value can be READ as a verdict, never about
// whether the verdict is right:
//
//   - a blank or duplicated unit id makes coverage unreadable — two rows for
//     one obligation are not a disposition, they are two;
//   - a disposition outside the closed set says nothing at all;
//   - a non-`pass` disposition with a blank reason is a judgement without a
//     reason, which is not review;
//   - a blank citation is a reference to nothing.
func assessmentIntegrity(units []UnitResult) error {
	if len(units) > maxAssessmentUnits {
		return fmt.Errorf("assessment carries %d units (> %d)", len(units), maxAssessmentUnits)
	}
	total := 0
	seen := make(map[string]bool, len(units))
	for _, u := range units {
		switch {
		case strings.TrimSpace(u.Unit) == "":
			return errors.New("assessment unit has a blank id")
		case seen[u.Unit]:
			return fmt.Errorf("assessment reports unit %q more than once", u.Unit)
		case !knownDisposition(u.Disposition):
			return fmt.Errorf("assessment unit %q has unknown disposition %q", u.Unit, u.Disposition)
		case u.Disposition != DispositionPass && strings.TrimSpace(u.Reason) == "":
			return fmt.Errorf("assessment unit %q is %q with no reason; a judgement without a reason is not review", u.Unit, u.Disposition)
		}
		seen[u.Unit] = true
		total += len(u.Unit) + len(u.Reason)
		for _, c := range u.Citations {
			if strings.TrimSpace(c) == "" {
				return fmt.Errorf("assessment unit %q carries a blank citation", u.Unit)
			}
			total += len(c)
		}
	}
	if total > maxAssessmentBytes {
		return fmt.Errorf("assessment text exceeds %d bytes", maxAssessmentBytes)
	}
	return nil
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
		Canon:        RecordCanon,
		EpisodeID:    id.Episode,
		Mode:         meta.ExecutionMode,
		ResolvedSpec: meta.ResolvedSpec,
		Contract:     frozen,
		Alpha:        StationRecord{ExecutionID: a.exec, Artifacts: a.artifacts},
		Matter:       a.matter,
		Beta:         StationRecord{ExecutionID: b.exec, Artifacts: b.artifacts},
		Review:       b.review,
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

	// Invocation authority (Pi round-5 D1, round-7 D1): the resolved spec is
	// inside the one boundary. The seat declarations are opaque here — fill
	// semantics live with the fills — but they must be present and non-null;
	// mode itself is bound to the parent-trusted metadata by VerifyClosure,
	// so a mode rewrite cannot survive even with a recomputed digest.
	add(r.ResolvedSpec.Version == "" || r.ResolvedSpec.DeclaredProtocol == "",
		InvalidRecord, "resolved spec missing version/protocol")
	add(validSeatEnvelope(r.ResolvedSpec.Alpha) != nil, InvalidRecord, "alpha declaration is not a tagged object")
	add(validSeatEnvelope(r.ResolvedSpec.Beta) != nil, InvalidRecord, "beta declaration is not a tagged object")

	// Required arrays are arrays, never null (round-6 D2): the closure schema
	// admits no null, so a nil slice surviving to canonical JSON must not
	// self-verify.
	add(r.Alpha.Artifacts == nil || r.Beta.Artifacts == nil,
		InvalidRecord, "station artifacts must be arrays, not null")

	// Identity: non-empty and pairwise distinct across the whole triple.
	add(r.EpisodeID == "" || r.Alpha.ExecutionID == "" || r.Beta.ExecutionID == "", InvalidIdentity, "missing identity")
	add(r.Alpha.ExecutionID == r.Beta.ExecutionID, InvalidIdentity, "stations share an execution id")
	add(r.EpisodeID == r.Alpha.ExecutionID || r.EpisodeID == r.Beta.ExecutionID, InvalidIdentity, "episode id aliases a station id")

	// Contract validity (same rules validateSpec enforces on the honest path).
	add(r.Contract.ID == "", InvalidRecord, "contract id is empty")
	for _, slot := range r.Contract.opaqueSlots() {
		add(opaqueSlotIntegrity(slot.raw) != nil, InvalidRecord,
			"contract "+slot.name+" is not structurally admissible")
	}
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
	if err := assessmentIntegrity(r.Review.Assessment); err != nil {
		// Integrity, not contract-unmet: a malformed assessment is not a failing
		// review, it is a value no verifier can read as a review at all.
		add(true, InvalidRecord, "assessment is not well formed: "+err.Error())
	}
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

// validate is V(contract, receipt) — the CCNF signature (Pi round-6 D1): the
// trusted expected contract arrives as an argument, never out of the receipt
// being judged. It checks the one digest, the full record boundary
// (validateRecord), that the embedded frozen snapshot equals the expected
// contract, and contract satisfaction against the EXPECTED contract — so a
// substituted embedded contract cannot re-scope what V requires. Producer
// authority is positional and fails closed: a required artifact with an
// invalid producer never resolves to a side.
func validate(expected Contract, rc Receipt) Verdict {
	r := rc.Record
	fs := validateRecord(r)
	if sha256hex(r.canonicalBytes()) != rc.ScopeLiftDigest {
		fs = append(fs, Failure{InvalidRecord, "scope-lift digest does not recompute"})
	}
	if !bytes.Equal(mustJSON(r.Contract), mustJSON(expected)) {
		fs = append(fs, Failure{InvalidRecord, "record contract does not match the expected contract"})
	}

	for _, req := range expected.RequiredEvidence {
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
	// Every non-`pass` disposition is its own contract-unmet failure, named by
	// unit. This is what makes the verdict RE-DERIVABLE: a seat cannot report
	// findings and set pass=true, because V never reads pass to decide this —
	// it reads the dispositions the record carries. A verdict a reader cannot
	// recompute from the receipt is not a verdict.
	//
	// `unverified` fails exactly as `finding` does, and deliberately: an
	// obligation nobody could decide has not been met, whatever the reason. The
	// two stay DISTINCT in the record and in the failure text, because "we
	// checked and it is wrong" and "we could not check" call for different next
	// work — but neither is an accepted episode.
	for _, u := range r.Review.Assessment {
		if u.Disposition != DispositionPass {
			fs = append(fs, Failure{ContractUnmet, fmt.Sprintf("assessment unit %s: %s: %s", u.Unit, u.Disposition, u.Reason)})
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

// decide is δ(receipt, verdict) — the canonical CCNF signature (Pi round-7
// C1). The v0 router is verdict-only (routes by failure class: integrity
// fails closed; contract-unmet repairs; pass accepts — pure); the receipt is
// where a later δ policy would look, never a new authority surface.
func decide(_ Receipt, v Verdict) Decision {
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

// validateMeta rejects an unmetadata'd or incoherent invocation before alpha
// runs (Pi round-5 D1): the resolved spec is invocation authority, so a run
// that cannot state its version and declared protocol — or whose mode and
// seat declarations disagree — is a malfunction, not a closable episode.
func validateMeta(m RunMeta) error {
	if !knownMode(m.ExecutionMode) {
		return fmt.Errorf("cellkernel: unknown execution mode %q", m.ExecutionMode)
	}
	rs := m.ResolvedSpec
	if rs.Version == "" || rs.DeclaredProtocol == "" {
		return errors.New("cellkernel: resolved spec must carry version and declared_protocol")
	}
	// Fill semantics are opaque to the kernel (round-6 C1); the declarations
	// only have to exist. Mode truth is bound to the parent-trusted metadata
	// at VerifyClosure (round-7 D1).
	if err := validSeatEnvelope(rs.Alpha); err != nil {
		return fmt.Errorf("cellkernel: alpha declaration: %w", err)
	}
	if err := validSeatEnvelope(rs.Beta); err != nil {
		return fmt.Errorf("cellkernel: beta declaration: %w", err)
	}
	return nil
}

// validSeatEnvelope checks the ONLY thing the generic boundary knows about a
// seat declaration: it is a JSON object carrying a non-empty string `fill`.
// Everything else stays opaque. Without this a direct caller could emit `{}`,
// `[]`, or a scalar, self-verify it, and still fail the closure schema — the
// kernel would be certifying something no reader could parse.
func validSeatEnvelope(raw json.RawMessage) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return errors.New("must be a JSON object")
	}
	tag, ok := obj["fill"]
	if !ok {
		return errors.New("must carry a fill")
	}
	var id string
	if err := json.Unmarshal(tag, &id); err != nil || id == "" {
		return errors.New("fill must be a non-empty string")
	}
	return nil
}

// opaqueSlotIntegrity is the COMPLETE set of rules the kernel applies to EVERY
// opaque contract slot — issue, design and subject — and it is one function
// rather than three because the kernel's interest in them is identical: they
// must be a JSON object and stay within the declared bound. Valid JSON is not
// taste — canonicalBytes() serializes the record with encoding/json, which
// cannot represent a RawMessage that is not JSON, so an unparseable slot would
// silently collapse the canonical bytes the one digest is taken over. An
// absent slot is admissible here; requiring one is a protocol's rule, not the
// kernel's.
//
// THE OBJECT RULE IS THE BOUNDARY'S, NOT ANY PROFILE'S. It states what a slot
// IS at this layer: a carrier for one tagged value, whose tag and vocabulary
// belong to whatever admitted it. A bare scalar is not a tagged anything — no
// reader downstream can ask it what it is — so it is inadmissible here whatever
// a profile would have made of it. The kernel still reads no key and knows no
// kind. Without the rule the two authorities disagreed: this function admitted
// any valid JSON while the closure schema declares the three slots `{...}`, so
// a record with a scalar slot self-verified here and then failed `cue vet` with
// a type conflict. Today's three CDS admitters all require objects, which made
// the disagreement unreachable through the shipped door and reachable by the
// first Door that does not.
func opaqueSlotIntegrity(raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	if len(raw) > MaxOpaqueSlotBytes {
		return fmt.Errorf("exceeds %d bytes", MaxOpaqueSlotBytes)
	}
	if !json.Valid(raw) {
		return errors.New("is not valid JSON")
	}
	// Valid JSON already, so the only way this fails is a non-object. `null`
	// decodes into a nil map without error and is not an object either.
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil || obj == nil {
		return errors.New("is not a JSON object")
	}
	return nil
}

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
	for _, slot := range s.Contract.opaqueSlots() {
		if err := opaqueSlotIntegrity(slot.raw); err != nil {
			return fmt.Errorf("cellkernel: contract.%s: %w", slot.name, err)
		}
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
