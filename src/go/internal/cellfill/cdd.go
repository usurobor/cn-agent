package cellfill

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
)

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }

// The generic cdd fills. These are protocol-neutral seats the cdd package
// itself owns: smoke stubs and the two honest mechanical reviewers.
const (
	FillStubAlpha       = "cdd.stub"
	FillStubBeta        = "cdd.stub"
	FillBoolAlpha       = "cdd.bool"
	FillBoolCheckBeta   = "cdd.bool-check"
	FillMechanicalUnmet = "cdd.mechanical-unmet"
)

// CddFills returns the statically assembled generic fills.
func CddFills() Registry {
	return Registry{
		// No generic seat reads the contract's subject: a stub fabricates its
		// own side's artifacts, a bool answers from its declaration and the
		// mechanical reviewers judge what they are handed, so all of them run
		// with an empty binding and none declares a requirement.
		Alpha: map[string]AlphaFill{
			FillStubAlpha: {Construct: stubAlphaFactory},
			FillBoolAlpha: {Construct: boolAlphaFactory},
		},
		Beta: map[string]BetaFill{
			FillStubBeta:        {Construct: stubBetaFactory},
			FillBoolCheckBeta:   {Construct: boolCheckBetaFactory},
			FillMechanicalUnmet: {Construct: mechanicalUnmetFactory},
		},
	}
}

// --- cdd.stub -------------------------------------------------------------

type stubDecl struct {
	Fill string `json:"fill"`
}

// decodeTagOnly is the shape shared by the fills whose entire declaration is
// the tag: each states its own exact key set rather than deriving one.
func decodeTagOnly(decl json.RawMessage, d *stubDecl) error {
	if err := OnlyKeys(decl, "declaration", "fill"); err != nil {
		return err
	}
	return StrictDecode(decl, d)
}

// The generic seats take their projection and IGNORE it: a stub fabricates its
// own side's artifacts, a bool answers from its declaration, and the
// mechanical-unmet reviewer refuses on principle rather than on obligations —
// so none of them is held to any methodology. Named `_` rather than accepted
// and dropped silently, so the disinterest is stated.
func stubAlphaFactory(_ context.Context, decl json.RawMessage, _ cellmethod.View) (ConstructedAlpha, error) {
	var d stubDecl
	if err := decodeTagOnly(decl, &d); err != nil {
		return ConstructedAlpha{}, fmt.Errorf("fill %q: %w", FillStubAlpha, err)
	}
	canon, _ := json.Marshal(d)
	return ConstructedAlpha{
		Constructed: Constructed{Decl: canon, Mode: cellkernel.ModeStub},
		Seat:        stubAlpha{},
	}, nil
}

func stubBetaFactory(_ context.Context, decl json.RawMessage, _ cellmethod.View) (ConstructedBeta, error) {
	var d stubDecl
	if err := decodeTagOnly(decl, &d); err != nil {
		return ConstructedBeta{}, fmt.Errorf("fill %q: %w", FillStubBeta, err)
	}
	canon, _ := json.Marshal(d)
	return ConstructedBeta{
		Constructed: Constructed{Decl: canon, Mode: cellkernel.ModeStub},
		Seat:        stubBeta{},
	}, nil
}

// stubAlpha fabricates its own side's required artifacts — smoke only, and
// the stub mode makes the episode honestly `simulated`.
type stubAlpha struct{}

func (stubAlpha) Produce(_ context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	var cands []cellkernel.ArtifactCandidate
	for _, req := range in.Contract.RequiredEvidence {
		if req.Producer != cellkernel.RoleAlpha {
			continue
		}
		cands = append(cands, cellkernel.ArtifactCandidate{ID: req.ID, Kind: req.Kind, Text: "stub-alpha:" + req.ID})
	}
	return cellkernel.AlphaOutput{
		Matter:    cellkernel.Matter{Data: "stub-alpha produced for " + strconv.Quote(in.Contract.Goal)},
		Artifacts: cands,
	}, nil
}

type stubBeta struct{}

func (stubBeta) Review(_ context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	var cands []cellkernel.ArtifactCandidate
	for _, req := range in.Contract.RequiredEvidence {
		if req.Producer != cellkernel.RoleBeta {
			continue
		}
		cands = append(cands, cellkernel.ArtifactCandidate{ID: req.ID, Kind: req.Kind, Text: "stub-beta:" + req.ID})
	}
	return cellkernel.BetaOutput{
		Review:    cellkernel.Review{Pass: true, Notes: "stub-beta accepted (smoke)"},
		Artifacts: cands,
	}, nil
}

// --- cdd.bool / cdd.bool-check --------------------------------------------

type boolDecl struct {
	Fill  string `json:"fill"`
	Value string `json:"value"`
}

func boolAlphaFactory(_ context.Context, decl json.RawMessage, _ cellmethod.View) (ConstructedAlpha, error) {
	var d boolDecl
	if err := OnlyKeys(decl, "cdd.bool", "fill", "value"); err != nil {
		return ConstructedAlpha{}, fmt.Errorf("fill %q: %w", FillBoolAlpha, err)
	}
	if err := StrictDecode(decl, &d); err != nil {
		return ConstructedAlpha{}, fmt.Errorf("fill %q: %w", FillBoolAlpha, err)
	}
	v, err := strconv.ParseBool(d.Value)
	if err != nil {
		return ConstructedAlpha{}, fmt.Errorf("fill %q: value %q is not a bool", FillBoolAlpha, d.Value)
	}
	canon, _ := json.Marshal(d)
	return ConstructedAlpha{
		Constructed: Constructed{Decl: canon, Mode: cellkernel.ModeMechanical},
		Seat:        cellkernel.BoolAlpha{Value: v},
	}, nil
}

func boolCheckBetaFactory(_ context.Context, decl json.RawMessage, _ cellmethod.View) (ConstructedBeta, error) {
	var d stubDecl
	if err := decodeTagOnly(decl, &d); err != nil {
		return ConstructedBeta{}, fmt.Errorf("fill %q: %w", FillBoolCheckBeta, err)
	}
	canon, _ := json.Marshal(d)
	return ConstructedBeta{
		Constructed: Constructed{Decl: canon, Mode: cellkernel.ModeMechanical},
		Seat:        cellkernel.BoolBeta{},
	}, nil
}

// --- cdd.mechanical-unmet -------------------------------------------------

func mechanicalUnmetFactory(_ context.Context, decl json.RawMessage, _ cellmethod.View) (ConstructedBeta, error) {
	var d stubDecl
	if err := decodeTagOnly(decl, &d); err != nil {
		return ConstructedBeta{}, fmt.Errorf("fill %q: %w", FillMechanicalUnmet, err)
	}
	canon, _ := json.Marshal(d)
	return ConstructedBeta{
		Constructed: Constructed{Decl: canon, Mode: cellkernel.ModeMechanical},
		Seat:        MechanicalUnmet{},
	}, nil
}

// MechanicalUnmet is the honest Case-2 reviewer: a mechanical seat that
// cannot judge whether matter meets a goal, and therefore never passes it. A
// non-empty diff is not a met contract — it is work awaiting review. The
// episode closes needs_repair with the measured matter preserved, until an
// independent rented beta (Case 3) supplies real judgement.
type MechanicalUnmet struct{}

func (MechanicalUnmet) Review(_ context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	return cellkernel.BetaOutput{Review: cellkernel.Review{
		Pass: false,
		Notes: fmt.Sprintf("mechanical seat cannot judge the goal; matter (%d bytes) preserved for independent review",
			len(in.Matter.Data)),
	}}, nil
}
