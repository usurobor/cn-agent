package cellkernel

import (
	"context"
	"strconv"
)

// Case 1 on the ladder: a one-shot mechanical cell over a bool. α produces a
// bool as its matter; β passes iff the bool is true. No repair loop — a false
// bool closes the episode as needs_repair (contract-unmet), which a Drive wrapper
// (future) would re-attempt. This is the smallest cell that exercises both a
// PASS→accepted path and a FAIL→needs_repair path through the real kernel.

// BoolAlpha produces a fixed bool as its matter, plus one evidence ref recording
// what it produced. Deterministic: the same BoolAlpha always produces the same
// matter (no cognition rented here).
type BoolAlpha struct {
	Value bool
}

func (a BoolAlpha) Produce(_ context.Context, _ Contract) (AlphaResult, error) {
	return AlphaResult{
		Matter: Matter{Data: strconv.FormatBool(a.Value)},
		EvidenceRefs: []EvidenceRef{{
			ID:                  "bool",
			Kind:                "value",
			Ref:                 strconv.FormatBool(a.Value),
			ProducerExecutionID: "boolalpha",
		}},
	}, nil
}

// BoolBeta discriminates: it passes iff the matter parses as bool true. A matter
// that does not parse is contract-unmet (not a malfunction) — β ran, it said no.
type BoolBeta struct{}

func (BoolBeta) Review(_ context.Context, _ Contract, m Matter) (BetaResult, error) {
	v, err := strconv.ParseBool(m.Data)
	if err != nil || !v {
		return BetaResult{Review: Review{Pass: false, Notes: "matter is not bool true"}}, nil
	}
	return BetaResult{Review: Review{Pass: true, Notes: "matter is bool true"}}, nil
}

// BoolSpec is the one-shot bool cell with the given α value and a contract that
// requires the α "bool" evidence ref to be present.
func BoolSpec(value bool) Spec {
	return Spec{
		Contract: Contract{
			ID:               "cell-bool",
			Goal:             "produce bool true",
			RequiredEvidence: []RequiredRef{{ID: "bool", Kind: "value"}},
		},
		Alpha: BoolAlpha{Value: value},
		Beta:  BoolBeta{},
	}
}
