package cellspec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// Builtin v0 seat profiles. v0 has no cognition; a profile selects a mechanical
// seat pair. Rented cognition (Phase 3) adds a provider-backed profile.
const (
	ProfileStub = "stub" // smoke: fabricates required artifacts; non-authoritative `simulated`
	ProfileBool = "bool" // real: alpha produces a bool, beta INDEPENDENTLY verifies it
)

func isKnownProfile(p string) bool { return p == ProfileStub || p == ProfileBool }

// buildProfile constructs the seat pair for a resolved spec's profile and its
// execution mode. A stub run is non-authoritative (kernel status `simulated`).
func buildProfile(r Resolved) (cellkernel.Alpha, cellkernel.Beta, cellkernel.ExecutionMode, error) {
	switch r.Spec.Profile {
	case ProfileStub:
		return stubAlpha{skills: r.AlphaSkills}, stubBeta{skills: r.BetaSkills}, cellkernel.ModeStub, nil
	case ProfileBool:
		v, ok := r.Params["value"]
		if !ok {
			return nil, nil, "", fmt.Errorf("profile %q requires parameter %q", ProfileBool, "value")
		}
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, nil, "", fmt.Errorf("profile %q: value %q is not a bool", ProfileBool, v)
		}
		return cellkernel.BoolAlpha{Value: b}, cellkernel.BoolBeta{}, cellkernel.ModeMechanical, nil
	default:
		return nil, nil, "", fmt.Errorf("unknown profile %q", r.Spec.Profile)
	}
}

// stubAlpha produces a summary matter and candidates for the α-side required
// artifacts. Provenance is positional — a stub can only fill its own side.
type stubAlpha struct{ skills []string }

func (a stubAlpha) Produce(_ context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	matter := cellkernel.Matter{Data: fmt.Sprintf("stub-alpha produced for %q with skills [%s]", in.Contract.Goal, strings.Join(a.skills, ", "))}
	var cands []cellkernel.ArtifactCandidate
	for _, req := range in.Contract.RequiredEvidence {
		if req.Producer != cellkernel.RoleAlpha {
			continue
		}
		cands = append(cands, cellkernel.ArtifactCandidate{ID: req.ID, Kind: req.Kind, Text: "stub-alpha:" + req.ID})
	}
	return cellkernel.AlphaOutput{Matter: matter, Artifacts: cands}, nil
}

// stubBeta accepts and produces candidates for the β-side required artifacts.
type stubBeta struct{ skills []string }

func (b stubBeta) Review(_ context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	var cands []cellkernel.ArtifactCandidate
	for _, req := range in.Contract.RequiredEvidence {
		if req.Producer != cellkernel.RoleBeta {
			continue
		}
		cands = append(cands, cellkernel.ArtifactCandidate{ID: req.ID, Kind: req.Kind, Text: "stub-beta:" + req.ID})
	}
	return cellkernel.BetaOutput{
		Review:    cellkernel.Review{Pass: true, Notes: fmt.Sprintf("stub-beta accepted with skills [%s]", strings.Join(b.skills, ", "))},
		Artifacts: cands,
	}, nil
}
