package cellkernel

import "context"

// The empty cell: α does nothing, β accepts unconditionally. It exercises the
// full five-step closure and terminates in `accepted`, proving the loop turns.
// It is Case 0 on the ladder; real cells swap NoopAlpha/AcceptBeta for α/β that
// produce and discriminate real matter.

// NoopAlpha produces no matter and no evidence.
type NoopAlpha struct{}

func (NoopAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{}, nil
}

// AcceptBeta accepts any matter.
type AcceptBeta struct{}

func (AcceptBeta) Review(context.Context, Contract, Matter) (BetaResult, error) {
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
