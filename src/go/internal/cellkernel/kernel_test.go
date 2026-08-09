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

func (s seqIDs) Mint() (Identity, error) { return Identity{Episode: s.ep, Alpha: s.a, Beta: s.b}, nil }

func testMeta(mode ExecutionMode) RunMeta {
	return RunMeta{ExecutionMode: mode, ResolvedSpec: ResolvedSpec{Version: "cnos.cellspec.v0", DeclaredProtocol: "p", Profile: "bool", AlphaSkills: []string{}, BetaSkills: []string{}}}
}

func mechEnvelope(t *testing.T, s Spec) Envelope {
	t.Helper()
	env, err := RunEpisode(context.Background(), s,
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}), WithMeta(testMeta(ModeMechanical)))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return env
}

func TestEmptyCellRunsToAccepted(t *testing.T) {
	env := mechEnvelope(t, EmptySpec())
	if env.Status != Accepted {
		t.Fatalf("status: want accepted, got %q", env.Status)
	}
	if err := VerifyEnvelope(env); err != nil {
		t.Fatalf("envelope must self-verify: %v", err)
	}
}

func TestBoolCellAcceptedAndRepair(t *testing.T) {
	acc := mechEnvelope(t, BoolSpec(true))
	if acc.Status != Accepted {
		t.Fatalf("bool true: status=%q", acc.Status)
	}
	rep := mechEnvelope(t, BoolSpec(false))
	if rep.Status != NeedsRepair || rep.Repair == nil {
		t.Fatalf("bool false: status=%q repair=%v", rep.Status, rep.Repair)
	}
}

// D5: a stub run is non-authoritative `simulated`, never accepted.
func TestStubIsSimulated(t *testing.T) {
	env, err := RunEpisode(context.Background(), EmptySpec(),
		WithIDSource(seqIDs{"ep-s", "alpha-s", "beta-s"}), WithMeta(testMeta(ModeStub)))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if env.Status != Simulated {
		t.Fatalf("stub status: want simulated, got %q", env.Status)
	}
	if err := VerifyEnvelope(env); err != nil {
		t.Fatalf("stub envelope must verify: %v", err)
	}
}

// D1: identity is per-invocation.
func TestEpisodeIdentityIsPerInvocation(t *testing.T) {
	a, _ := RunEpisode(context.Background(), BoolSpec(true), WithMeta(testMeta(ModeMechanical)))
	b, _ := RunEpisode(context.Background(), BoolSpec(true), WithMeta(testMeta(ModeMechanical)))
	if a.Receipt.EpisodeID == b.Receipt.EpisodeID || a.Receipt.AlphaExecutionID == b.Receipt.AlphaExecutionID {
		t.Fatal("identities are not per-invocation")
	}
}

func TestResolvedInputBoundInEnvelope(t *testing.T) {
	mk := func(v string) Envelope {
		m := testMeta(ModeMechanical)
		m.ResolvedSpec.AlphaSkills = []string{v}
		env, _ := RunEpisode(context.Background(), BoolSpec(true), WithMeta(m))
		return env
	}
	if mk("go").ResolvedSpecHash == mk("rust").ResolvedSpecHash {
		t.Fatal("runs differing in resolved input share a resolved_spec_hash")
	}
}

// --- D1/D2: whole-envelope verification; every field re-derives ----------

func roundTrip(t *testing.T, env Envelope) Envelope {
	t.Helper()
	b, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Envelope
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestEnvelopeReVerifiesAfterSerialization(t *testing.T) {
	env := roundTrip(t, mechEnvelope(t, BoolSpec(true)))
	if err := VerifyEnvelope(env); err != nil {
		t.Fatalf("round-tripped envelope must verify: %v", err)
	}
}

func TestTamperedEnvelopeFails(t *testing.T) {
	tamper := map[string]func(*Envelope){
		"flip status":          func(e *Envelope) { e.Status = Rejected },
		"flip decision":        func(e *Envelope) { e.Decision = Reject },
		"flip verdict pass":    func(e *Envelope) { e.Verdict.Pass = false },
		"flip protocol valid":  func(e *Envelope) { e.ProtocolValidated = true },
		"flip execution mode":  func(e *Envelope) { e.ExecutionMode = ModeStub },
		"rewrite spec version": func(e *Envelope) { e.ResolvedSpec.Version = "x" },
		"rewrite inner matter": func(e *Envelope) { e.Receipt.Matter.Data = "changed" },
		"forge inner evidence": func(e *Envelope) { e.Receipt.Evidence[0].Content = "x" },
		"drop repair":          func(e *Envelope) { e.Repair = &RepairRequest{Reason: "spurious"} },
	}
	for name, mut := range tamper {
		t.Run(name, func(t *testing.T) {
			env := roundTrip(t, mechEnvelope(t, BoolSpec(true)))
			mut(&env)
			if err := VerifyEnvelope(env); err == nil {
				t.Fatalf("tamper %q passed verification", name)
			}
		})
	}
}

// --- D3: typed failure routing ------------------------------------------

func TestIntegrityFailureFailsClosed(t *testing.T) {
	// A verdict carrying an integrity failure must reject (fail closed), not
	// route to the ordinary alpha repair path.
	integrity := Verdict{Pass: false, Failures: []Failure{{InvalidEvidence, "x"}}}
	if d := decide(integrity); d != Reject {
		t.Fatalf("integrity failure -> %q, want reject", d)
	}
	unmet := Verdict{Pass: false, Failures: []Failure{{ContractUnmet, "x"}}}
	if d := decide(unmet); d != RepairDispatch {
		t.Fatalf("contract-unmet -> %q, want repair_dispatch", d)
	}
}

// --- D4: fail-closed identity minting ------------------------------------

type errIDs struct{}

func (errIDs) Mint() (Identity, error) { return Identity{}, errors.New("rand failure") }

type dupIDs struct{}

func (dupIDs) Mint() (Identity, error) { return Identity{"x", "x", "x"}, nil }

func TestIdentityFailsClosed(t *testing.T) {
	for name, src := range map[string]IDSource{"mint error": errIDs{}, "not distinct": dupIDs{}} {
		t.Run(name, func(t *testing.T) {
			if _, err := RunEpisode(context.Background(), BoolSpec(true), WithIDSource(src)); err == nil {
				t.Fatalf("want error for %s", name)
			}
		})
	}
}

// --- C1: evidence bytes must be valid UTF-8 ------------------------------

type badBytesAlpha struct{}

func (badBytesAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: "m"}, Evidence: []EvidenceCandidate{{ID: "x", Kind: "k", Bytes: []byte{0xff, 0xfe}}}}, nil
}

func TestNonUTF8EvidenceRejected(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c", Goal: "g"}, Alpha: badBytesAlpha{}, Beta: AcceptBeta{}}
	if _, err := RunEpisode(context.Background(), s, WithMeta(testMeta(ModeMechanical))); err == nil {
		t.Fatal("want error for non-UTF-8 evidence")
	}
}

// --- Malfunction, self-cert, producer authority --------------------------

type brokenBeta struct{}

func (brokenBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{}, errors.New("backend unavailable")
}

func TestBetaMalfunctionReturnsError(t *testing.T) {
	s := EmptySpec()
	s.Beta = brokenBeta{}
	if _, err := RunEpisode(context.Background(), s, WithMeta(testMeta(ModeMechanical))); err == nil {
		t.Fatal("want error from beta malfunction")
	}
}

func TestNilSeatsFailClosed(t *testing.T) {
	s := EmptySpec()
	s.Alpha = nil
	if _, err := RunEpisode(context.Background(), s, WithMeta(testMeta(ModeMechanical))); err == nil {
		t.Fatal("want error for nil alpha")
	}
}

func TestCancelledContextFailsClosed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RunEpisode(ctx, BoolSpec(true), WithMeta(testMeta(ModeMechanical))); err == nil {
		t.Fatal("want error for cancelled context")
	}
}

type forgingAlpha struct{}

func (forgingAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: "junk"}}, nil
}

type rejectBeta struct{}

func (rejectBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{Review: Review{Pass: false, Notes: "reject"}}, nil
}

func TestNoSelfCertification(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c", Goal: "g"}, Alpha: forgingAlpha{}, Beta: rejectBeta{}}
	env := mechEnvelope(t, s)
	if env.Status == Accepted || env.Verdict.Pass {
		t.Fatalf("self-cert leaked: status=%q", env.Status)
	}
}

type alphaMintingBetaEvidence struct{}

func (alphaMintingBetaEvidence) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: "m"}, Evidence: []EvidenceCandidate{
		{ID: "diff", Kind: "diff", Bytes: []byte("d")},
		{ID: "beta_review", Kind: "review", Bytes: []byte("fake")},
	}}, nil
}

func TestAlphaCannotMintBetaEvidence(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g", RequiredEvidence: []RequiredRef{
			{ID: "diff", Kind: "diff", Producer: RoleAlpha},
			{ID: "beta_review", Kind: "review", Producer: RoleBeta},
		}},
		Alpha: alphaMintingBetaEvidence{},
		Beta:  AcceptBeta{},
	}
	env := mechEnvelope(t, s)
	if env.Verdict.Pass {
		t.Fatal("V passed although beta_review was forged by alpha")
	}
}

func TestHostileAlphaCannotMutateFrozenContract(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g", RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
		Alpha:    mutatingAlpha{},
		Beta:     AcceptBeta{},
	}
	env := mechEnvelope(t, s)
	if env.Verdict.Pass {
		t.Fatal("frozen contract was relaxed by a hostile alpha")
	}
	if len(env.Receipt.Contract.RequiredEvidence) != 1 || env.Receipt.Contract.RequiredEvidence[0].ID != "diff" {
		t.Fatalf("frozen contract mutated: %+v", env.Receipt.Contract.RequiredEvidence)
	}
}

type mutatingAlpha struct{}

func (mutatingAlpha) Produce(_ context.Context, c Contract) (AlphaResult, error) {
	for i := range c.RequiredEvidence {
		c.RequiredEvidence[i].ID = "neutralized"
	}
	c.RequiredEvidence = nil
	return AlphaResult{Matter: Matter{Data: "m"}}, nil
}

func TestBetaInputHashStableAndSensitive(t *testing.T) {
	base := BetaInput{Contract: Contract{ID: "c"}, ContractHash: "h", Matter: Matter{Data: "m"}, PolicyID: BetaInputPolicyID}
	if hashBetaInput(base) != hashBetaInput(base) {
		t.Fatal("beta-input hash not stable")
	}
	other := base
	other.PolicyID = "different"
	if hashBetaInput(base) == hashBetaInput(other) {
		t.Fatal("beta-input hash ignored the policy id")
	}
}

func TestKernelRejectsInvalidSpec(t *testing.T) {
	cases := map[string]Spec{
		"empty contract id": {Contract: Contract{ID: ""}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
		"bad producer":      {Contract: Contract{ID: "c", RequiredEvidence: []RequiredRef{{ID: "x", Kind: "k", Producer: "gamma"}}}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
		"dup required":      {Contract: Contract{ID: "c", RequiredEvidence: []RequiredRef{{ID: "x", Kind: "k", Producer: RoleAlpha}, {ID: "x", Kind: "k", Producer: RoleBeta}}}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
	}
	for name, s := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := RunEpisode(context.Background(), s, WithMeta(testMeta(ModeMechanical))); err == nil {
				t.Fatalf("invalid spec %q ran", name)
			}
		})
	}
}

type oversizeAlpha struct{}

func (oversizeAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: strings.Repeat("x", maxMatterBytes+1)}}, nil
}

func TestBoundedOutput(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c"}, Alpha: oversizeAlpha{}, Beta: AcceptBeta{}}
	if _, err := RunEpisode(context.Background(), s, WithMeta(testMeta(ModeMechanical))); err == nil {
		t.Fatal("want error for oversized matter")
	}
}
