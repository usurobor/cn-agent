package cellkernel

import "context"

// The empty cell: α does nothing, β accepts unconditionally. It exercises the
// full five-step closure and terminates in `accepted`, proving the loop turns.
// It is the first cell spec; real cells swap NoopAlpha/AcceptBeta for α/β that
// produce and discriminate real matter.

// NoopAlpha produces no matter.
type NoopAlpha struct{}

func (NoopAlpha) Produce(context.Context, Contract) (Matter, error) {
	return Matter{}, nil
}

// AcceptBeta accepts any matter.
type AcceptBeta struct{}

func (AcceptBeta) Review(context.Context, Contract, Matter) (Review, error) {
	return Review{Pass: true, Notes: "empty cell: accept"}, nil
}

// EmptySpec is the empty cell: noop α, accept β, mechanical γ/V/δ.
func EmptySpec() Spec {
	return Spec{
		Contract: Contract{ID: "cell-0", Goal: "empty"},
		Alpha:    NoopAlpha{},
		Beta:     AcceptBeta{},
	}
}
