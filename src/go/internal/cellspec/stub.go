package cellspec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// Builtin seat profiles. A profile selects a seat pair; the loader — never the
// kernel — owns this whitelist, so the kernel stays domain-neutral.
const (
	ProfileStub      = "stub"      // smoke: fabricates required artifacts; non-authoritative `simulated`
	ProfileBool      = "bool"      // real: alpha produces a bool, beta INDEPENDENTLY verifies it
	ProfileCognitive = "cognitive" // rented alpha behind a provider; beta still mechanical (Case 2)
	ProfileCode      = "code"      // rented alpha changes real code; the runtime measures the diff
)

func isKnownProfile(p string) bool {
	return p == ProfileStub || p == ProfileBool || p == ProfileCognitive || p == ProfileCode
}

// Providers the cognitive profile may rent. Closed set: a typo must fail
// resolution rather than silently pick a backend.
const (
	ProviderClaude = "claude" // the Claude Code CLI — real cognition
	ProviderFake   = "fake"   // deterministic, rents nothing (CI)
)

// buildProfile constructs the seat pair for a resolved spec's profile and the
// execution mode that honestly describes how the work was produced: `stub`
// fabricates (non-authoritative `simulated`), `mechanical` is deterministic
// and reproducible from the record, `cognitive` means a provider held a seat.
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
	case ProfileCognitive:
		provider, mode, err := buildProvider(r.Params["provider"])
		if err != nil {
			return nil, nil, "", err
		}
		return cellcog.Alpha{Provider: provider, Skills: r.AlphaSkills}, cellcog.MatterBeta{}, mode, nil
	case ProfileCode:
		coder, mode, err := buildCoder(r.Params["provider"])
		if err != nil {
			return nil, nil, "", err
		}
		base := r.Params["base_sha"]
		if base == "" {
			return nil, nil, "", fmt.Errorf("profile %q requires a non-empty %q parameter", ProfileCode, "base_sha")
		}
		repo := r.Params["repo"]
		if repo == "" {
			repo = "."
		}
		return cellcog.CodeAlpha{Coder: coder, Repo: repo, BaseRef: base, Skills: r.AlphaSkills},
			cellcog.MatterBeta{}, mode, nil
	default:
		return nil, nil, "", fmt.Errorf("unknown profile %q", r.Spec.Profile)
	}
}

// buildProvider pairs a provider with the mode that tells the truth about it:
// only a provider that actually rents cognition may run `cognitive`.
func buildProvider(name string) (cellcog.Provider, cellkernel.ExecutionMode, error) {
	switch name {
	case ProviderClaude:
		return cellcog.ClaudeCLI{}, cellkernel.ModeCognitive, nil
	case ProviderFake:
		return cellcog.Fake{}, cellkernel.ModeMechanical, nil
	default:
		return nil, "", fmt.Errorf("profile %q: unknown provider %q (want %q or %q)",
			ProfileCognitive, name, ProviderClaude, ProviderFake)
	}
}

// buildCoder is buildProvider's file-capable twin; the mode tells the same
// truth about which one actually rents cognition.
func buildCoder(name string) (cellcog.Coder, cellkernel.ExecutionMode, error) {
	switch name {
	case ProviderClaude:
		return cellcog.ClaudeCLI{}, cellkernel.ModeCognitive, nil
	case ProviderFake:
		return cellcog.FakeCoder{}, cellkernel.ModeMechanical, nil
	default:
		return nil, "", fmt.Errorf("profile %q: unknown provider %q (want %q or %q)",
			ProfileCode, name, ProviderClaude, ProviderFake)
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
