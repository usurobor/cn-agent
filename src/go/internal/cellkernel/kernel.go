// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) kernel: the substrate-independent five-step closure a
// coherence cell runs at one scope.
//
//  1. matter   := α.produce(contract)
//  2. review   := β.review(contract, matter)
//  3. receipt  := γ.close(contract, matter, review, evidence)
//  4. verdict  := V(contract, receipt)
//  5. decision := δ.decide(receipt, verdict)
//
// See src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md and
// cnos#711. This is the walking skeleton: no recursion (repair_dispatch), no
// evidence dereferencing, no CUE binding yet — those are increments into seams
// that already exist here. α and β are supplied by the caller (the cell's
// customization); γ/V/δ default to the mechanical kernel implementations.
//
// Purity boundary (eng/go §2.17): α, β, and V may do IO or rent cognition, so
// they take a context and return an error. γ and δ are pure — γ closes a
// receipt value; δ decides from (receipt, verdict). A seat that returns an
// error is a malfunction (the cell cannot close); a review that returns
// Pass=false is contract-unmet (the cell closes blocked).
package cellkernel

import (
	"context"
	"fmt"
)

// --- Artifacts (typed; minimal shapes for now) --------------------------

// Contract is the cell's input: what to produce and what "done" means.
type Contract struct {
	ID   string
	Goal string
	// acceptance criteria, scope, non-goals: added as the kernel grows.
}

// Matter is α's product.
type Matter struct{ Data string }

// Review is β's discrimination of the matter against the contract. β consumes
// matter only — never evidence (that is V's job at step 4).
type Review struct {
	Pass  bool
	Notes string
}

// Evidence accumulates during α/β work and is bound into the receipt by γ.
type Evidence struct{ Refs []string }

// Receipt is the parent-facing artifact γ emits: the single typed surface that
// crosses the scope boundary. V dereferences everything it needs from here.
type Receipt struct {
	Contract Contract
	Matter   Matter
	Review   Review
	Evidence Evidence
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

// --- The five seats -----------------------------------------------------

// Alpha produces matter against the contract. It sees only the contract, and
// may rent cognition or do IO, so it takes a context and may fail.
type Alpha interface {
	Produce(ctx context.Context, c Contract) (Matter, error)
}

// Beta discriminates the matter against the contract. Matter only, no evidence.
// A returned error is a malfunction; Review{Pass:false} is a contract-unmet
// review, not an error.
type Beta interface {
	Review(ctx context.Context, c Contract, m Matter) (Review, error)
}

// Gamma closes the cell: binds contract, matter, review, and evidence into a
// typed receipt. γ is pure — it produces a receipt value and does not validate,
// decide, or perform IO.
type Gamma interface {
	Close(c Contract, m Matter, r Review, e Evidence) Receipt
}

// Validator (V) is a predicate: it reads the contract and receipt, dereferences
// evidence from the receipt (IO), and emits PASS/FAIL.
type Validator interface {
	Validate(ctx context.Context, c Contract, rc Receipt) (Verdict, error)
}

// Delta decides at the boundary on receipt and verdict only (never evidence).
// δ is pure.
type Delta interface {
	Decide(rc Receipt, v Verdict) Decision
}

// --- Spec + result ------------------------------------------------------

// Spec is a cell: a contract plus the seats. α and β are the caller's
// customization; γ/V/δ default to the mechanical kernel when nil.
type Spec struct {
	Contract Contract
	Alpha    Alpha
	Beta     Beta
	Gamma    Gamma     // optional; defaults to DefaultGamma
	V        Validator // optional; defaults to DefaultV
	Delta    Delta     // optional; defaults to DefaultDelta
}

// Outcome is the cell's terminal state, determined by (verdict, decision).
type Outcome string

const (
	Accepted Outcome = "accepted"
	Degraded Outcome = "degraded"
	Blocked  Outcome = "blocked"
	Invalid  Outcome = "invalid" // non-terminal in full CCNF; δ re-decides
)

// ClosedCell is the kernel's terminal object at one scope.
type ClosedCell struct {
	Contract Contract
	Matter   Matter
	Review   Review
	Receipt  Receipt
	Verdict  Verdict
	Decision Decision
	Outcome  Outcome
}

// Run executes the CCNF five-step closure at one scope. A returned error means
// a seat malfunctioned and the cell did not close; otherwise the ClosedCell's
// Outcome reports how it closed. No recursion yet: a repair_dispatch decision
// surfaces as Blocked rather than opening a child cell.
func Run(ctx context.Context, s Spec) (ClosedCell, error) {
	g := gammaOrDefault(s.Gamma)
	v := validatorOrDefault(s.V)
	d := deltaOrDefault(s.Delta)

	matter, err := s.Alpha.Produce(ctx, s.Contract) // 1
	if err != nil {
		return ClosedCell{}, fmt.Errorf("alpha produce: %w", err)
	}
	review, err := s.Beta.Review(ctx, s.Contract, matter) // 2
	if err != nil {
		return ClosedCell{}, fmt.Errorf("beta review: %w", err)
	}
	receipt := g.Close(s.Contract, matter, review, Evidence{}) // 3
	verdict, err := v.Validate(ctx, s.Contract, receipt)       // 4
	if err != nil {
		return ClosedCell{}, fmt.Errorf("validate receipt: %w", err)
	}
	decision := d.Decide(receipt, verdict) // 5

	return ClosedCell{
		Contract: s.Contract,
		Matter:   matter,
		Review:   review,
		Receipt:  receipt,
		Verdict:  verdict,
		Decision: decision,
		Outcome:  outcomeOf(verdict, decision),
	}, nil
}

func gammaOrDefault(g Gamma) Gamma {
	if g == nil {
		return DefaultGamma{}
	}
	return g
}

func validatorOrDefault(v Validator) Validator {
	if v == nil {
		return DefaultV{}
	}
	return v
}

func deltaOrDefault(d Delta) Delta {
	if d == nil {
		return DefaultDelta{}
	}
	return d
}

// outcomeOf maps (verdict, decision) to one of the four CCNF outcomes.
func outcomeOf(v Verdict, d Decision) Outcome {
	switch {
	case v.Pass && (d == Accept || d == Release):
		return Accepted
	case !v.Pass && d == Override:
		return Degraded
	case d == Reject || d == RepairDispatch:
		return Blocked
	default:
		// PASS+override, or !PASS+accept/release: inconsistent.
		return Invalid
	}
}

// --- Default mechanical γ / V / δ ---------------------------------------

// DefaultGamma mechanically closes: binds inputs into a receipt, no judgment.
type DefaultGamma struct{}

func (DefaultGamma) Close(c Contract, m Matter, r Review, e Evidence) Receipt {
	return Receipt{Contract: c, Matter: m, Review: r, Evidence: e}
}

// DefaultV is the minimal predicate: PASS iff β passed. A richer V will cue-vet
// the receipt and dereference evidence against the contract.
type DefaultV struct{}

func (DefaultV) Validate(_ context.Context, _ Contract, rc Receipt) (Verdict, error) {
	if !rc.Review.Pass {
		return Verdict{Pass: false, Failed: []string{"review.pass=false"}}, nil
	}
	return Verdict{Pass: true}, nil
}

// DefaultDelta is the minimal boundary policy: PASS→accept, FAIL→repair_dispatch.
type DefaultDelta struct{}

func (DefaultDelta) Decide(_ Receipt, v Verdict) Decision {
	if v.Pass {
		return Accept
	}
	return RepairDispatch
}
