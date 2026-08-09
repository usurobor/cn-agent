package cellkernel

import (
	"context"
	"strconv"
)

// Case 1: a one-shot mechanical cell where β INDEPENDENTLY checks α's artifact
// (Pi D6 — a non-tautological proof, unlike the stub profile). α produces a bool
// as its matter and mints the required "bool" evidence (producer α); β re-reads
// the matter from its runtime-owned bundle and passes iff it is true. A hostile
// α that produces false — or that tries to mint the acceptance itself — cannot
// make β pass, and V independently checks producer authority and integrity.

// BoolAlpha produces a fixed bool as its matter plus the required α evidence.
type BoolAlpha struct {
	Value bool
}

func (a BoolAlpha) Produce(_ context.Context, _ Contract) (AlphaResult, error) {
	s := strconv.FormatBool(a.Value)
	return AlphaResult{
		Matter: Matter{Data: s},
		Evidence: []EvidenceRef{{
			ID:      "bool",
			Kind:    "value",
			Ref:     "bool://" + s,
			Content: s,
		}},
	}, nil
}

// BoolBeta independently checks the matter from its bundle: pass iff it parses
// as bool true. It does not trust α's claim — it re-reads and decides.
type BoolBeta struct{}

func (BoolBeta) Review(_ context.Context, in BetaInput) (BetaResult, error) {
	v, err := strconv.ParseBool(in.Matter.Data)
	if err != nil || !v {
		return BetaResult{Review: Review{Pass: false, Notes: "matter is not bool true"}}, nil
	}
	return BetaResult{Review: Review{Pass: true, Notes: "beta independently verified matter is bool true"}}, nil
}

// BoolSpec is the one-shot bool cell requiring α-produced "bool" evidence.
func BoolSpec(value bool) Spec {
	return Spec{
		Contract: Contract{
			ID:               "cell-bool",
			Goal:             "produce bool true",
			RequiredEvidence: []RequiredRef{{ID: "bool", Kind: "value", Producer: RoleAlpha}},
		},
		Alpha: BoolAlpha{Value: value},
		Beta:  BoolBeta{},
	}
}
