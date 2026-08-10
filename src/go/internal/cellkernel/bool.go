package cellkernel

import (
	"context"
	"strconv"
)

// Case 1: a one-shot mechanical cell where β INDEPENDENTLY reviews the matter
// from its own immutable projection (non-tautological). α produces a bool as
// its matter plus the required "bool" artifact; β re-reads only the matter
// from its BetaInput and passes iff it is true — V, not β, checks the
// required α artifact. Provenance is positional: the "bool" artifact
// satisfies the contract only because it sits under record.Alpha.

// BoolAlpha produces a fixed bool as its matter plus the "bool" artifact.
type BoolAlpha struct {
	Value bool
}

func (a BoolAlpha) Produce(_ context.Context, _ AlphaInput) (AlphaOutput, error) {
	s := strconv.FormatBool(a.Value)
	return AlphaOutput{
		Matter:    Matter{Data: s},
		Artifacts: []ArtifactCandidate{{ID: "bool", Kind: "value", Text: s}},
	}, nil
}

// BoolBeta independently checks the matter from its projection: pass iff it
// parses as bool true. It does not trust α's claim — it re-reads and decides.
type BoolBeta struct{}

func (BoolBeta) Review(_ context.Context, in BetaInput) (BetaOutput, error) {
	v, err := strconv.ParseBool(in.Matter.Data)
	if err != nil || !v {
		return BetaOutput{Review: Review{Pass: false, Notes: "matter is not bool true"}}, nil
	}
	return BetaOutput{Review: Review{Pass: true, Notes: "beta independently verified matter is bool true"}}, nil
}

// BoolSpec is the one-shot bool cell requiring the α-side "bool" artifact.
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
