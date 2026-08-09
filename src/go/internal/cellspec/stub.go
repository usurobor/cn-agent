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
	ProfileStub = "stub" // smoke: fabricates required evidence, beta accepts (tautological)
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

// stubAlpha produces a summary matter and candidate bytes for exactly the
// α-owned required evidence. The runtime stamps producer/exec/digest, so a stub
// can only satisfy its own role.
type stubAlpha struct{ skills []string }

func (a stubAlpha) Produce(_ context.Context, c cellkernel.Contract) (cellkernel.AlphaResult, error) {
	matter := cellkernel.Matter{Data: fmt.Sprintf("stub-alpha produced for %q with skills [%s]", c.Goal, strings.Join(a.skills, ", "))}
	var cands []cellkernel.EvidenceCandidate
	for _, req := range c.RequiredEvidence {
		if req.Producer != cellkernel.RoleAlpha {
			continue
		}
		cands = append(cands, cellkernel.EvidenceCandidate{ID: req.ID, Kind: req.Kind, Bytes: []byte("stub-alpha:" + req.ID)})
	}
	return cellkernel.AlphaResult{Matter: matter, Evidence: cands}, nil
}

// stubBeta accepts and produces candidate bytes for β-owned required evidence
// other than beta_review (the runtime mints beta_review from the actual review).
type stubBeta struct{ skills []string }

func (b stubBeta) Review(_ context.Context, in cellkernel.BetaInput) (cellkernel.BetaResult, error) {
	var cands []cellkernel.EvidenceCandidate
	for _, req := range in.Contract.RequiredEvidence {
		if req.Producer != cellkernel.RoleBeta || req.ID == "beta_review" {
			continue
		}
		cands = append(cands, cellkernel.EvidenceCandidate{ID: req.ID, Kind: req.Kind, Bytes: []byte("stub-beta:" + req.ID)})
	}
	return cellkernel.BetaResult{
		Review:   cellkernel.Review{Pass: true, Notes: fmt.Sprintf("stub-beta accepted with skills [%s]", strings.Join(b.skills, ", "))},
		Evidence: cands,
	}, nil
}
