package cellkernel

import "context"

// The empty cell: α produces nothing, β accepts unconditionally. It exercises
// the full closure and terminates `accepted` — Case 0, the runner's smoke
// reference. It carries no required evidence, so V's evidence checks are vacuous.

// NoopAlpha produces no matter and no evidence.
type NoopAlpha struct{}

func (NoopAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{}, nil
}

// AcceptBeta accepts any matter. It reviews over the runtime-owned BetaInput but
// (being the empty smoke seat) checks nothing.
type AcceptBeta struct{}

func (AcceptBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{Review: Review{Pass: true, Notes: "empty cell: accept"}}, nil
}

// EmptySpec is the empty cell: noop α, accept β, kernel-owned γ/V/δ.
func EmptySpec() Spec {
	return Spec{
		Contract: Contract{ID: "cell-0", Goal: "empty"},
		Alpha:    NoopAlpha{},
		Beta:     AcceptBeta{},
	}
}
