// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) single-episode kernel: the substrate-independent
// five-step closure a coherence cell runs at one scope.
//
//  1. matter   := α.produce(contract)
//  2. review   := β.review(betaInput)
//  3. receipt  := γ.close(contract, matter, review, evidence)   [kernel-owned]
//  4. verdict  := V(contract, receipt)                          [kernel-owned]
//  5. decision := δ.decide(receipt, verdict)                    [kernel-owned]
//
// See docs/architecture/CDS-CELL-MIGRATION.md and cnos#711. This is the
// single-episode engine: no repair loop (a future Drive wrapper), no
// composition (α-proposes / runtime-executes), no protocol dispatch (the
// runner emits a generic episode receipt; protocol-specific validation is a
// later phase).
//
// Authority model (Pi β msg-cn-pi-cnos-cell-prototype-beta-32, D1–D4):
//
//   - The runtime FREEZES the contract at episode start: it deep-copies and
//     hashes it, hands each seat an isolated copy, and V/γ bind the frozen
//     snapshot + hash. A seat cannot mutate the terms it is judged against (D3).
//   - Evidence is RUNTIME-AUTHENTICATED, not seat-asserted: the runtime stamps
//     each evidence ref's producer role and execution id from the seat that
//     actually returned it and hashes its content. α cannot mint β's evidence,
//     because the producer role is assigned by the runtime, not claimed by the
//     seat (D2).
//   - β receives a runtime-owned BetaInput (the review surface + the frozen
//     contract + authenticated α evidence + a bundle hash), never α's private
//     state (D4).
//   - Only α and β are open seats. γ binds, V validates bindings, δ applies
//     boundary policy — all kernel-owned and mechanical (no injectable seats).
//
// Purity boundary (eng/go §2.17): α and β may do IO or rent cognition, so they
// take a context and may fail. A seat error is a malfunction (the episode
// cannot close); a review with Pass=false is contract-unmet (closes needing
// repair). An inconsistent (verdict, decision) pair is a typed error, never a
// returned closed cell.
package cellkernel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Role names which seat produced an artifact. The runtime assigns it; a seat
// cannot claim a role it does not hold.
type Role string

const (
	RoleAlpha Role = "alpha"
	RoleBeta  Role = "beta"
)

// BetaInputPolicyID identifies the review-surface policy: β sees the frozen
// contract, the matter, and α's authenticated evidence — never α's private
// reasoning or session state.
const BetaInputPolicyID = "cnos.cellkernel.beta-input.v0"

// --- Artifacts (typed; minimal shapes for now) --------------------------

// RequiredRef names an evidence ref the contract requires, and which producer
// role is authorized to mint it. V checks presence *and* producer authority.
type RequiredRef struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Producer Role   `json:"producer"`
}

// Contract is the cell's input: what to produce, what "done" means, and the
// evidence the closed receipt must carry (with producer authority).
type Contract struct {
	ID               string        `json:"id"`
	Goal             string        `json:"goal"`
	RequiredEvidence []RequiredRef `json:"required_evidence,omitempty"`
}

// clone deep-copies the contract so a seat cannot mutate the frozen snapshot
// through a shared slice backing array (D3).
func (c Contract) clone() Contract {
	cp := c
	if c.RequiredEvidence != nil {
		cp.RequiredEvidence = append([]RequiredRef(nil), c.RequiredEvidence...)
	}
	return cp
}

// canonicalHash is the content address of the frozen contract.
func (c Contract) canonicalHash() string {
	b, _ := json.Marshal(c) // struct field order is deterministic
	return sha256hex(b)
}

// Matter is α's product. It is a public result surface, never a channel for a
// model's private reasoning.
type Matter struct {
	Data string `json:"data"`
}

// EvidenceRef is one typed evidence record. Producer, ProducerExecutionID, and
// SHA256 are RUNTIME-stamped (a seat cannot forge them); Content is the produced
// bytes the runtime hashes (kept out of the serialized receipt).
type EvidenceRef struct {
	ID                  string `json:"id"`
	Kind                string `json:"kind"`
	Producer            Role   `json:"producer"`
	ProducerExecutionID string `json:"producer_execution_id"`
	Ref                 string `json:"ref"`
	SHA256              string `json:"sha256"`
	Content             string `json:"-"` // hashed by the runtime; not serialized
}

// AlphaResult is α's output: the matter plus candidate evidence (the runtime
// authenticates the candidates).
type AlphaResult struct {
	Matter   Matter
	Evidence []EvidenceRef
}

// Review is β's discrimination of the matter against the contract.
type Review struct {
	Pass  bool   `json:"pass"`
	Notes string `json:"notes"`
}

// BetaResult is β's output: the review plus candidate evidence.
type BetaResult struct {
	Review   Review
	Evidence []EvidenceRef
}

// BetaInput is the runtime-owned review surface handed to β (D4). It carries the
// frozen contract + hash, the matter, α's authenticated evidence, and the
// bundle/policy identity — and nothing of α's private state.
type BetaInput struct {
	Contract      Contract
	ContractHash  string
	Matter        Matter
	AlphaEvidence []EvidenceRef
	PolicyID      string
	BundleHash    string
}

// Receipt is the parent-facing artifact γ emits. It binds the frozen contract
// and its hash, the episode and seat execution ids, and the authenticated
// evidence.
type Receipt struct {
	EpisodeID        string        `json:"episode_id"`
	Contract         Contract      `json:"contract"`
	ContractHash     string        `json:"contract_hash"`
	Matter           Matter        `json:"matter"`
	Review           Review        `json:"review"`
	Evidence         []EvidenceRef `json:"evidence_refs"`
	AlphaExecutionID string        `json:"alpha_execution_id"`
	BetaExecutionID  string        `json:"beta_execution_id"`
}

// Verdict is V's typed output.
type Verdict struct {
	Pass     bool     `json:"pass"`
	Failed   []string `json:"failed,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// Decision is δ's boundary decision.
type Decision string

const (
	Accept         Decision = "accept"
	Release        Decision = "release"
	Override       Decision = "override"
	Reject         Decision = "reject"
	RepairDispatch Decision = "repair_dispatch"
)

// --- The two open seats -------------------------------------------------

// Alpha produces matter against the contract (an isolated frozen copy).
type Alpha interface {
	Produce(ctx context.Context, c Contract) (AlphaResult, error)
}

// Beta discriminates over the runtime-owned review surface. A returned error is
// a malfunction; Review{Pass:false} is contract-unmet.
type Beta interface {
	Review(ctx context.Context, in BetaInput) (BetaResult, error)
}

// Spec is a single-episode cell: a contract plus the two open seats.
type Spec struct {
	Contract Contract
	Alpha    Alpha
	Beta     Beta
}

// Status is how a single episode closed.
type Status string

const (
	Accepted    Status = "accepted"
	Degraded    Status = "degraded"
	Rejected    Status = "rejected"
	NeedsRepair Status = "needs_repair" // nonterminal; the parent stays open
)

// RepairRequest is the typed reason an episode needs repair.
type RepairRequest struct {
	Reason string   `json:"reason"`
	Failed []string `json:"failed,omitempty"`
}

// EpisodeResult is the kernel's object for one closed-or-held episode.
type EpisodeResult struct {
	EpisodeID    string
	ContractHash string
	Contract     Contract
	Matter       Matter
	Review       Review
	Receipt      Receipt
	Verdict      Verdict
	Decision     Decision
	Status       Status
	Repair       *RepairRequest // set iff Status == NeedsRepair
}

// ErrInvalidClosure is returned for an inconsistent (verdict, decision) pair.
var ErrInvalidClosure = errors.New("cellkernel: inconsistent (verdict, decision) pair")

// RunEpisode executes the CCNF five-step closure for one episode.
func RunEpisode(ctx context.Context, s Spec) (EpisodeResult, error) {
	if seatIsNil(s.Alpha) { // D4/hardening: reject nil and typed-nil seats.
		return EpisodeResult{}, errors.New("cellkernel: spec has nil alpha")
	}
	if seatIsNil(s.Beta) {
		return EpisodeResult{}, errors.New("cellkernel: spec has nil beta")
	}
	if err := ctx.Err(); err != nil { // honor a cancelled context.
		return EpisodeResult{}, fmt.Errorf("cellkernel: context: %w", err)
	}

	// Freeze the contract: deep-copy + hash. Everything downstream binds the
	// frozen snapshot; seats receive isolated copies (D3).
	frozen := s.Contract.clone()
	frozenHash := frozen.canonicalHash()
	episodeID := "ep-" + frozenHash[:12]
	alphaExec := "alpha-" + episodeID
	betaExec := "beta-" + episodeID

	aRes, err := s.Alpha.Produce(ctx, frozen.clone()) // 1 — α gets its own copy
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("alpha produce: %w", err)
	}
	alphaEv := authenticate(aRes.Evidence, RoleAlpha, alphaExec)

	bin := BetaInput{
		Contract:      frozen.clone(),
		ContractHash:  frozenHash,
		Matter:        aRes.Matter,
		AlphaEvidence: cloneEvidence(alphaEv),
		PolicyID:      BetaInputPolicyID,
		BundleHash:    bundleHash(frozenHash, aRes.Matter, alphaEv),
	}
	bRes, err := s.Beta.Review(ctx, bin) // 2 — β gets the runtime-owned surface
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("beta review: %w", err)
	}
	betaEv := authenticate(bRes.Evidence, RoleBeta, betaExec)

	evidence := append(cloneEvidence(alphaEv), betaEv...)
	receipt := closeReceipt(frozen, frozenHash, episodeID, aRes.Matter, bRes.Review, evidence, alphaExec, betaExec) // 3
	verdict := validate(frozen, frozenHash, receipt)                                                                // 4
	decision := decide(verdict)                                                                                     // 5

	status, err := statusOf(verdict, decision)
	if err != nil {
		return EpisodeResult{}, err
	}

	res := EpisodeResult{
		EpisodeID:    episodeID,
		ContractHash: frozenHash,
		Contract:     frozen,
		Matter:       aRes.Matter,
		Review:       bRes.Review,
		Receipt:      receipt,
		Verdict:      verdict,
		Decision:     decision,
		Status:       status,
	}
	if status == NeedsRepair {
		res.Repair = &RepairRequest{Reason: "contract unmet", Failed: verdict.Failed}
	}
	return res, nil
}

// statusOf maps a consistent (verdict, decision) to a Status, or returns
// ErrInvalidClosure for an inconsistent pair (D1).
func statusOf(v Verdict, d Decision) (Status, error) {
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

// --- Runtime evidence authentication ------------------------------------

// authenticate stamps each candidate evidence ref with the producing role, the
// seat execution id, and the content hash — overwriting anything the seat
// claimed. This is what stops α from minting β's evidence (D2).
func authenticate(candidates []EvidenceRef, role Role, execID string) []EvidenceRef {
	out := make([]EvidenceRef, 0, len(candidates))
	for _, e := range candidates {
		e.Producer = role
		e.ProducerExecutionID = execID
		e.SHA256 = sha256hex([]byte(e.Content))
		out = append(out, e)
	}
	return out
}

func cloneEvidence(in []EvidenceRef) []EvidenceRef {
	return append([]EvidenceRef(nil), in...)
}

func bundleHash(contractHash string, m Matter, alphaEv []EvidenceRef) string {
	h := sha256.New()
	h.Write([]byte(contractHash))
	h.Write([]byte(m.Data))
	for _, e := range alphaEv {
		h.Write([]byte(e.ID))
		h.Write([]byte(e.SHA256))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// --- Kernel-owned mechanical γ / V / δ (D2; not injectable) --------------

// closeReceipt (γ) mechanically binds the frozen snapshot and authenticated
// evidence. Pure: no judgment, no IO.
func closeReceipt(c Contract, hash, episodeID string, m Matter, r Review, evidence []EvidenceRef, alphaExec, betaExec string) Receipt {
	return Receipt{
		EpisodeID:        episodeID,
		Contract:         c,
		ContractHash:     hash,
		Matter:           m,
		Review:           r,
		Evidence:         evidence,
		AlphaExecutionID: alphaExec,
		BetaExecutionID:  betaExec,
	}
}

// validate (V) verifies the receipt's bindings against the frozen contract: the
// contract hash, then — for each required ref — a bound evidence entry with the
// matching id, kind, and authorized producer, unique, with an intact content
// hash. It does not re-judge β's prose (Pi D2/D4).
func validate(frozen Contract, frozenHash string, rc Receipt) Verdict {
	var failed []string

	if rc.ContractHash != frozenHash {
		failed = append(failed, "contract binding mismatch")
	}

	// Uniqueness: no two bound evidence refs may share an id.
	seen := make(map[string]int)
	for _, e := range rc.Evidence {
		seen[e.ID]++
	}
	for id, n := range seen {
		if n > 1 {
			failed = append(failed, "duplicate evidence id: "+id)
		}
	}

	// Integrity: each bound ref's content must hash to its stamped SHA256.
	for _, e := range rc.Evidence {
		if sha256hex([]byte(e.Content)) != e.SHA256 {
			failed = append(failed, "evidence hash mismatch: "+e.ID)
		}
	}

	// Required presence + producer authority.
	for _, req := range frozen.RequiredEvidence {
		if !hasAuthorizedEvidence(rc.Evidence, req) {
			failed = append(failed, fmt.Sprintf("missing/unauthorized required evidence: %s (want producer %s)", req.ID, req.Producer))
		}
	}

	if !rc.Review.Pass {
		failed = append(failed, "review.pass=false")
	}

	if len(failed) > 0 {
		return Verdict{Pass: false, Failed: failed}
	}
	return Verdict{Pass: true}
}

func hasAuthorizedEvidence(refs []EvidenceRef, req RequiredRef) bool {
	for _, e := range refs {
		if e.ID == req.ID && e.Kind == req.Kind && e.Producer == req.Producer {
			return true
		}
	}
	return false
}

// decide (δ) is the minimal mechanical boundary policy: PASS→accept,
// FAIL→repair_dispatch.
func decide(v Verdict) Decision {
	if v.Pass {
		return Accept
	}
	return RepairDispatch
}

// --- helpers ------------------------------------------------------------

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// seatIsNil reports whether a seat is nil or a typed-nil (interface holding a
// nil pointer), which a bare `== nil` misses.
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
