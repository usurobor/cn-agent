// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) single-episode kernel: the substrate-independent
// five-step closure a coherence cell runs at one scope.
//
//  1. matter   := α.produce(contract)
//  2. review   := β.review(contract, matter)
//  3. receipt  := γ.close(contract, matter, review, evidence)   [kernel-owned]
//  4. verdict  := V(contract, receipt)                          [kernel-owned]
//  5. decision := δ.decide(receipt, verdict)                    [kernel-owned]
//
// See src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md,
// docs/architecture/CDS-CELL-MIGRATION.md, and cnos#711. This is the
// single-episode engine: no repair loop (that is a Drive wrapper), no
// composition (that is α-proposes / runtime-executes), no CUE binding yet.
//
// Design corrections from Pi β (msg-cn-pi-cnos-cell-runner-cases-review-31):
//
//   - D1 honest closure: RunEpisode returns an EpisodeResult whose Status is a
//     terminal outcome (accepted|degraded|rejected) OR needs_repair (the parent
//     stays open). An inconsistent (verdict, decision) pair is a typed error,
//     never a returned result. Repair looping belongs to a future Drive.
//   - D2 no self-certification: the caller supplies Contract + α + β only. The
//     kernel OWNS mechanical γ/V/δ; there are no injectable seat interfaces, so
//     a rejecting β cannot be rewritten into an acceptance. V verifies bindings,
//     it does not merely mirror β.
//   - D3 evidence seam: α and β return typed EvidenceRefs that γ binds.
//   - D4 fail closed: a nil α or β is a wrapped error before any seat runs.
//
// Purity boundary (eng/go §2.17): α and β may do IO or rent cognition, so they
// take a context and may fail. γ/V/δ are mechanical and kernel-owned. A seat
// that returns an error is a malfunction (the episode cannot close); a review
// that returns Pass=false is contract-unmet (the episode closes needing repair).
package cellkernel

import (
	"context"
	"errors"
	"fmt"
)

// --- Artifacts (typed; minimal shapes for now) --------------------------

// RequiredRef names an evidence ref the contract requires γ to bind and V to
// find present. Kind is the ref's typed role (e.g. "diff", "review").
type RequiredRef struct {
	ID   string
	Kind string
}

// Contract is the cell's input: what to produce, what "done" means, and the
// evidence the closed receipt must carry.
type Contract struct {
	ID               string
	Goal             string
	RequiredEvidence []RequiredRef
	// acceptance criteria, scope, non-goals: added as the kernel grows.
}

// Matter is α's product.
type Matter struct{ Data string }

// EvidenceRef is one typed, content-addressed evidence record accrued during
// α/β work (Pi Q4 shape). ProducerExecutionID ties the ref to the seat run
// that produced it, so V can check it was not γ-authored after the fact.
type EvidenceRef struct {
	ID                  string `json:"id"`
	Kind                string `json:"kind"`
	Ref                 string `json:"ref"`
	SHA256              string `json:"sha256,omitempty"`
	ProducerExecutionID string `json:"producer_execution_id"`
}

// AlphaResult is α's output: the matter plus the evidence its run accrued.
type AlphaResult struct {
	Matter       Matter
	EvidenceRefs []EvidenceRef
}

// Review is β's discrimination of the matter against the contract. β consumes
// matter only — never evidence (that is V's job at step 4).
type Review struct {
	Pass  bool
	Notes string
}

// BetaResult is β's output: the review plus the evidence its run accrued.
type BetaResult struct {
	Review       Review
	EvidenceRefs []EvidenceRef
}

// Receipt is the parent-facing artifact γ emits: the single typed surface that
// crosses the scope boundary. V dereferences everything it needs from here.
type Receipt struct {
	Contract     Contract
	Matter       Matter
	Review       Review
	EvidenceRefs []EvidenceRef
}

// Verdict is V's typed output. WARN is not a verdict value; advisories live in
// Warnings (mirrors schemas/cdd/receipt.cue #ValidationVerdict).
type Verdict struct {
	Pass     bool
	Failed   []string
	Warnings []string
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

// Alpha produces matter against the contract. It sees only the contract, and
// may rent cognition or do IO, so it takes a context and may fail. It returns
// the evidence its run accrued (D3).
type Alpha interface {
	Produce(ctx context.Context, c Contract) (AlphaResult, error)
}

// Beta discriminates the matter against the contract. Matter only, no evidence.
// A returned error is a malfunction; Review{Pass:false} is a contract-unmet
// review, not an error.
type Beta interface {
	Review(ctx context.Context, c Contract, m Matter) (BetaResult, error)
}

// --- Spec + result ------------------------------------------------------

// Spec is a single-episode cell: a contract plus the two open seats. γ/V/δ are
// kernel-owned and mechanical, deliberately NOT part of the Spec (D2): a caller
// cannot inject a γ that certifies its own receipt.
type Spec struct {
	Contract Contract
	Alpha    Alpha
	Beta     Beta
}

// Status is how a single episode closed.
type Status string

const (
	// Terminal statuses: the episode is closed.
	Accepted Status = "accepted"
	Degraded Status = "degraded"
	Rejected Status = "rejected"
	// Non-terminal: the parent cell stays open; a Drive may re-attempt.
	NeedsRepair Status = "needs_repair"
)

// RepairRequest is the typed reason an episode needs repair. It is surfaced to
// a Drive loop (future); the kernel does not itself re-attempt.
type RepairRequest struct {
	Reason string   `json:"reason"`
	Failed []string `json:"failed,omitempty"`
}

// EpisodeResult is the kernel's object for one closed-or-held episode. When
// Status == NeedsRepair, Repair is set and the parent is still open; otherwise
// Status is terminal. An inconsistent (verdict, decision) pair is never
// represented here — RunEpisode returns an error instead (D1).
type EpisodeResult struct {
	Contract Contract
	Matter   Matter
	Review   Review
	Receipt  Receipt
	Verdict  Verdict
	Decision Decision
	Status   Status
	Repair   *RepairRequest // set iff Status == NeedsRepair
}

// ErrInvalidClosure is returned when (verdict, decision) is an inconsistent
// pair (PASS+override, or FAIL+accept/release). Such a pair is nonterminal in
// full CCNF (δ must re-decide); in this v0 kernel it is a typed malfunction
// rather than a returned result.
var ErrInvalidClosure = errors.New("cellkernel: inconsistent (verdict, decision) pair")

// RunEpisode executes the CCNF five-step closure for one episode. It returns:
//   - a terminal EpisodeResult (accepted|degraded|rejected), or
//   - a NeedsRepair EpisodeResult (the parent stays open), or
//   - an error: a seat malfunction, a fail-closed spec violation, or an
//     inconsistent closure (ErrInvalidClosure).
//
// It never returns a non-terminal state dressed as a closed cell.
func RunEpisode(ctx context.Context, s Spec) (EpisodeResult, error) {
	if s.Alpha == nil { // D4: fail closed before any seat runs.
		return EpisodeResult{}, errors.New("cellkernel: spec has nil alpha")
	}
	if s.Beta == nil {
		return EpisodeResult{}, errors.New("cellkernel: spec has nil beta")
	}

	aRes, err := s.Alpha.Produce(ctx, s.Contract) // 1
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("alpha produce: %w", err)
	}
	bRes, err := s.Beta.Review(ctx, s.Contract, aRes.Matter) // 2
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("beta review: %w", err)
	}

	evidence := append(append([]EvidenceRef{}, aRes.EvidenceRefs...), bRes.EvidenceRefs...)
	receipt := closeReceipt(s.Contract, aRes.Matter, bRes.Review, evidence) // 3 (kernel γ)
	verdict := validate(s.Contract, receipt)                                // 4 (kernel V)
	decision := decide(verdict)                                             // 5 (kernel δ)

	status, err := statusOf(verdict, decision)
	if err != nil {
		return EpisodeResult{}, err
	}

	res := EpisodeResult{
		Contract: s.Contract,
		Matter:   aRes.Matter,
		Review:   bRes.Review,
		Receipt:  receipt,
		Verdict:  verdict,
		Decision: decision,
		Status:   status,
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
		// PASS+override/reject, or FAIL+accept/release: inconsistent.
		return "", fmt.Errorf("%w: verdict.pass=%v decision=%q", ErrInvalidClosure, v.Pass, d)
	}
}

// --- Kernel-owned mechanical γ / V / δ (D2; not injectable) --------------

// closeReceipt (γ) mechanically binds the episode's inputs and accrued evidence
// into a receipt. It is pure: no judgment, no IO, no validation.
func closeReceipt(c Contract, m Matter, r Review, evidence []EvidenceRef) Receipt {
	return Receipt{Contract: c, Matter: m, Review: r, EvidenceRefs: evidence}
}

// validate (V) verifies the receipt's bindings — it does not re-judge the
// matter (Pi D2). Because γ is kernel-owned, β's review reaches V unrewritten;
// V additionally checks the contract binding and required-evidence presence,
// then derives PASS from the (trusted) review.
func validate(c Contract, rc Receipt) Verdict {
	var failed []string

	if rc.Contract.ID != c.ID {
		failed = append(failed, "contract binding mismatch")
	}
	for _, req := range c.RequiredEvidence {
		if !hasEvidence(rc.EvidenceRefs, req.ID, req.Kind) {
			failed = append(failed, "missing required evidence: "+req.ID)
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

func hasEvidence(refs []EvidenceRef, id, kind string) bool {
	for _, e := range refs {
		if e.ID == id && e.Kind == kind {
			return true
		}
	}
	return false
}

// decide (δ) is the minimal mechanical boundary policy: PASS→accept,
// FAIL→repair_dispatch. A richer δ (override, reject) arrives with policy.
func decide(v Verdict) Decision {
	if v.Pass {
		return Accept
	}
	return RepairDispatch
}
