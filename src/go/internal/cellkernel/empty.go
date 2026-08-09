package cellkernel

import "context"

// The empty cell: α produces nothing, β accepts unconditionally. Case 0 — the
// runner's smoke reference. It carries no required evidence.

// NoopAlpha produces no matter and no evidence.
type NoopAlpha struct{}

func (NoopAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{}, nil
}

// AcceptBeta accepts any matter, reviewing over the runtime-owned BetaInput.
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
