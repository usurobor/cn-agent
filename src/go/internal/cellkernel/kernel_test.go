package cellkernel

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// seqIDs is a deterministic id source for tests.
type seqIDs struct{ ep, a, b string }

func (s seqIDs) Mint() (Identity, error) {
	return Identity{Episode: s.ep, Alpha: s.a, Beta: s.b}, nil
}

func testMeta(mode ExecutionMode) RunMeta {
	// Profile must be mode-coherent: stub ⇔ profile "stub" (round-5 D1).
	profile := "bool"
	if mode == ModeStub {
		profile = "stub"
	}
	return RunMeta{ExecutionMode: mode, ResolvedSpec: ResolvedSpec{
		Version: "cnos.cellspec.v0", DeclaredProtocol: "p", Profile: profile,
		AlphaSkills: []string{}, BetaSkills: []string{},
	}}
}

func mechClosure(t *testing.T, s Spec) Closure {
	t.Helper()
	cl, err := RunEpisode(context.Background(), s, testMeta(ModeMechanical),
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return cl
}

func roundTrip(t *testing.T, cl Closure) Closure {
	t.Helper()
	b, err := json.Marshal(cl)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Closure
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestEmptyCellRunsToAccepted(t *testing.T) {
	cl := mechClosure(t, EmptySpec())
	if cl.Status != Accepted {
		t.Fatalf("status: want accepted, got %q", cl.Status)
	}
	if err := VerifyClosure(EmptySpec().Contract, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("closure must self-verify: %v", err)
	}
}

func TestBoolCellAcceptedAndRepair(t *testing.T) {
	if cl := mechClosure(t, BoolSpec(true)); cl.Status != Accepted {
		t.Fatalf("bool true: status=%q", cl.Status)
	}
	rep := mechClosure(t, BoolSpec(false))
	if rep.Status != NeedsRepair || rep.Repair == nil {
		t.Fatalf("bool false: status=%q repair=%v", rep.Status, rep.Repair)
	}
}

// Stub runs are non-authoritative `simulated`.
func TestStubIsSimulated(t *testing.T) {
	cl, err := RunEpisode(context.Background(), EmptySpec(), testMeta(ModeStub),
		WithIDSource(seqIDs{"ep-s", "alpha-s", "beta-s"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != Simulated {
		t.Fatalf("stub status: want simulated, got %q", cl.Status)
	}
	if err := VerifyClosure(EmptySpec().Contract, testMeta(ModeStub), cl); err != nil {
		t.Fatalf("stub closure must verify: %v", err)
	}
}

// Identity is per-invocation; resolved input is covered by the one digest.
func TestEpisodeIdentityIsPerInvocation(t *testing.T) {
	a, _ := RunEpisode(context.Background(), BoolSpec(true), testMeta(ModeMechanical))
	b, _ := RunEpisode(context.Background(), BoolSpec(true), testMeta(ModeMechanical))
	if a.Receipt.Record.EpisodeID == b.Receipt.Record.EpisodeID {
		t.Fatal("identities are not per-invocation")
	}
}

func TestResolvedInputChangesDigest(t *testing.T) {
	mk := func(skill string) Closure {
		m := testMeta(ModeMechanical)
		m.ResolvedSpec.AlphaSkills = []string{skill}
		cl, _ := RunEpisode(context.Background(), BoolSpec(true), m,
			WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
		return cl
	}
	if mk("go").Receipt.ScopeLiftDigest == mk("rust").Receipt.ScopeLiftDigest {
		t.Fatal("runs differing in resolved input share a scope-lift digest")
	}
}

// --- The one verification boundary ---------------------------------------

func TestClosureReVerifiesAfterSerialization(t *testing.T) {
	cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("round-tripped closure must verify: %v", err)
	}
}

func TestTamperedClosureFails(t *testing.T) {
	tamper := map[string]func(*Closure){
		"flip status":           func(c *Closure) { c.Status = Rejected },
		"flip decision":         func(c *Closure) { c.Decision = Reject },
		"flip verdict pass":     func(c *Closure) { c.Verdict.Pass = false },
		"flip protocol claim":   func(c *Closure) { c.ProtocolValidated = true },
		"rewrite record matter": func(c *Closure) { c.Receipt.Record.Matter.Data = "changed" },
		"rewrite record review": func(c *Closure) { c.Receipt.Record.Review.Pass = false },
		"forge alpha artifact":  func(c *Closure) { c.Receipt.Record.Alpha.Artifacts[0].Text = "x" },
		"rewrite resolved spec": func(c *Closure) { c.Receipt.Record.ResolvedSpec.Profile = "x" },
		"substitute digest":     func(c *Closure) { c.Receipt.ScopeLiftDigest = strings.Repeat("0", 64) },
		"spurious repair":       func(c *Closure) { c.Repair = &RepairRequest{Reason: "spurious"} },
		"move artifact to beta": func(c *Closure) {
			c.Receipt.Record.Beta.Artifacts = c.Receipt.Record.Alpha.Artifacts
			c.Receipt.Record.Alpha.Artifacts = nil
		},
	}
	for name, mut := range tamper {
		t.Run(name, func(t *testing.T) {
			cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
			mut(&cl)
			if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err == nil {
				t.Fatalf("tamper %q passed verification", name)
			}
		})
	}
}

// --- Seat isolation (Pi #44 required action 7) ----------------------------

// mutatingAlpha tries to relax the contract it was handed.
type mutatingAlpha struct{}

func (mutatingAlpha) Produce(_ context.Context, in AlphaInput) (AlphaOutput, error) {
	for i := range in.Contract.RequiredEvidence {
		in.Contract.RequiredEvidence[i].ID = "neutralized"
	}
	in.Contract.RequiredEvidence = nil
	return AlphaOutput{Matter: Matter{Data: "m"}}, nil
}

func TestSeatCannotMutateItsInputContract(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g",
			RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
		Alpha: mutatingAlpha{},
		Beta:  AcceptBeta{},
	}
	cl := mechClosure(t, s)
	if cl.Verdict.Pass {
		t.Fatal("frozen contract was relaxed by a hostile alpha")
	}
	if len(cl.Receipt.Record.Contract.RequiredEvidence) != 1 || cl.Receipt.Record.Contract.RequiredEvidence[0].ID != "diff" {
		t.Fatalf("frozen contract mutated: %+v", cl.Receipt.Record.Contract.RequiredEvidence)
	}
}

// projectionMutatingBeta mutates everything it receives (matter projection,
// frozen contract copy); the sealed originals must be unaffected.
type projectionMutatingBeta struct{}

func (projectionMutatingBeta) Review(_ context.Context, in BetaInput) (BetaOutput, error) {
	in.Matter.Data = "corrupted-by-beta"
	for i := range in.Contract.RequiredEvidence {
		in.Contract.RequiredEvidence[i].ID = "neutralized"
	}
	return BetaOutput{Review: Review{Pass: true, Notes: "beta tried to corrupt its inputs"}}, nil
}

func TestBetaCannotMutateSealedAlpha(t *testing.T) {
	s := BoolSpec(true)
	s.Beta = projectionMutatingBeta{}
	cl := mechClosure(t, s)
	if cl.Receipt.Record.Matter.Data != "true" {
		t.Fatalf("sealed alpha matter mutated via beta projection: %q", cl.Receipt.Record.Matter.Data)
	}
	if got := cl.Receipt.Record.Alpha.Artifacts[0].Text; got != "true" {
		t.Fatalf("sealed alpha artifact mutated via beta projection: %q", got)
	}
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("closure must still verify: %v", err)
	}
}

// A seat has no authority surface to forge: candidates carry only {id, kind,
// text}; a "beta_review"-labeled α artifact lands positionally under Alpha and
// cannot satisfy a β-side requirement.
type impersonatingAlpha struct{}

func (impersonatingAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	return AlphaOutput{Matter: Matter{Data: "m"}, Artifacts: []ArtifactCandidate{
		{ID: "diff", Kind: "diff", Text: "d"},
		{ID: "beta_signoff", Kind: "review", Text: "fake"},
	}}, nil
}

func TestAlphaArtifactCannotSatisfyBetaRequirement(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g", RequiredEvidence: []RequiredRef{
			{ID: "diff", Kind: "diff", Producer: RoleAlpha},
			{ID: "beta_signoff", Kind: "review", Producer: RoleBeta},
		}},
		Alpha: impersonatingAlpha{},
		Beta:  AcceptBeta{}, // produces no artifacts
	}
	cl := mechClosure(t, s)
	if cl.Verdict.Pass {
		t.Fatal("an alpha-side artifact satisfied a beta-side requirement")
	}
}

// --- Typed routing, identity, bounds, malfunction -------------------------

func TestIntegrityFailureFailsClosed(t *testing.T) {
	integrity := Verdict{Pass: false, Failures: []Failure{{InvalidRecord, "x"}}}
	if d := decide(Receipt{}, integrity); d != Reject {
		t.Fatalf("integrity failure -> %q, want reject", d)
	}
	unmet := Verdict{Pass: false, Failures: []Failure{{ContractUnmet, "x"}}}
	if d := decide(Receipt{}, unmet); d != RepairDispatch {
		t.Fatalf("contract-unmet -> %q, want repair_dispatch", d)
	}
}

type errIDs struct{}

func (errIDs) Mint() (Identity, error) { return Identity{}, errors.New("rand failure") }

type dupIDs struct{}

func (dupIDs) Mint() (Identity, error) { return Identity{"x", "x", "x"}, nil }

func TestIdentityFailsClosed(t *testing.T) {
	for name, src := range map[string]IDSource{"mint error": errIDs{}, "not distinct": dupIDs{}} {
		t.Run(name, func(t *testing.T) {
			if _, err := RunEpisode(context.Background(), BoolSpec(true), testMeta(ModeMechanical), WithIDSource(src)); err == nil {
				t.Fatalf("want error for %s", name)
			}
		})
	}
}

type badBytesAlpha struct{}

func (badBytesAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	return AlphaOutput{Matter: Matter{Data: "m"},
		Artifacts: []ArtifactCandidate{{ID: "x", Kind: "k", Text: string([]byte{0xff, 0xfe})}}}, nil
}

func TestNonUTF8ArtifactRejected(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c", Goal: "g"}, Alpha: badBytesAlpha{}, Beta: AcceptBeta{}}
	if _, err := RunEpisode(context.Background(), s, testMeta(ModeMechanical)); err == nil {
		t.Fatal("want error for non-UTF-8 artifact")
	}
}

type brokenBeta struct{}

func (brokenBeta) Review(context.Context, BetaInput) (BetaOutput, error) {
	return BetaOutput{}, errors.New("backend unavailable")
}

func TestBetaMalfunctionReturnsError(t *testing.T) {
	s := EmptySpec()
	s.Beta = brokenBeta{}
	if _, err := RunEpisode(context.Background(), s, testMeta(ModeMechanical)); err == nil {
		t.Fatal("want error from beta malfunction")
	}
}

type forgingAlpha struct{}

func (forgingAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	return AlphaOutput{Matter: Matter{Data: "junk"}}, nil
}

type rejectBeta struct{}

func (rejectBeta) Review(context.Context, BetaInput) (BetaOutput, error) {
	return BetaOutput{Review: Review{Pass: false, Notes: "reject"}}, nil
}

func TestNoSelfCertification(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c", Goal: "g"}, Alpha: forgingAlpha{}, Beta: rejectBeta{}}
	cl := mechClosure(t, s)
	if cl.Status == Accepted || cl.Verdict.Pass {
		t.Fatalf("self-cert leaked: status=%q", cl.Status)
	}
}

func TestCancelledContextFailsClosed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RunEpisode(ctx, BoolSpec(true), testMeta(ModeMechanical)); err == nil {
		t.Fatal("want error for cancelled context")
	}
}

func TestKernelRejectsInvalidSpec(t *testing.T) {
	cases := map[string]Spec{
		"empty contract id": {Contract: Contract{ID: ""}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
		"nil alpha":         {Contract: Contract{ID: "c"}, Alpha: nil, Beta: AcceptBeta{}},
		"bad producer":      {Contract: Contract{ID: "c", RequiredEvidence: []RequiredRef{{ID: "x", Kind: "k", Producer: "gamma"}}}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
		"dup required":      {Contract: Contract{ID: "c", RequiredEvidence: []RequiredRef{{ID: "x", Kind: "k", Producer: RoleAlpha}, {ID: "x", Kind: "k", Producer: RoleBeta}}}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
	}
	for name, s := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := RunEpisode(context.Background(), s, testMeta(ModeMechanical)); err == nil {
				t.Fatalf("invalid spec %q ran", name)
			}
		})
	}
}

type oversizeAlpha struct{}

func (oversizeAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	return AlphaOutput{Matter: Matter{Data: strings.Repeat("x", maxMatterBytes+1)}}, nil
}

func TestBoundedOutput(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c"}, Alpha: oversizeAlpha{}, Beta: AcceptBeta{}}
	if _, err := RunEpisode(context.Background(), s, testMeta(ModeMechanical)); err == nil {
		t.Fatal("want error for oversized matter")
	}
}
