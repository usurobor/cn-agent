package cellkernel

import "context"

// The empty cell: α produces nothing, β accepts unconditionally. Case 0 — the
// runner's smoke reference. It carries no required evidence.

// NoopAlpha produces no matter and no artifacts.
type NoopAlpha struct{}

func (NoopAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	return AlphaOutput{}, nil
}

// AcceptBeta accepts any matter.
type AcceptBeta struct{}

func (AcceptBeta) Review(context.Context, BetaInput) (BetaOutput, error) {
	return BetaOutput{Review: Review{Pass: true, Notes: "empty cell: accept"}}, nil
}

// EmptySpec is the empty cell: noop α, accept β, kernel-owned γ/V/δ.
func EmptySpec() Spec {
	return Spec{
		Contract: Contract{ID: "cell-0", Goal: "empty"},
		Alpha:    NoopAlpha{},
		Beta:     AcceptBeta{},
	}
}
