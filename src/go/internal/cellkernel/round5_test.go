package cellkernel

import (
	"context"
	"testing"
)

// Regression pairs for Pi β msg-cn-pi-cnos-pr718-fido-round5-beta-46 (D1).

// D1: meta is validated at the kernel boundary — empty fields or a
// profile/mode incoherence error out before alpha ever runs.
type recordingAlpha struct{ ran *bool }

func (a recordingAlpha) Produce(ctx context.Context, in AlphaInput) (AlphaOutput, error) {
	*a.ran = true
	return BoolAlpha{Value: true}.Produce(ctx, in)
}

func TestMetaValidatedAtIngress(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*RunMeta)
	}{
		{"empty version", func(m *RunMeta) { m.ResolvedSpec.Version = "" }},
		{"empty protocol", func(m *RunMeta) { m.ResolvedSpec.DeclaredProtocol = "" }},
		{"empty profile", func(m *RunMeta) { m.ResolvedSpec.Profile = "" }},
		{"stub mode, non-stub profile", func(m *RunMeta) {
			m.ExecutionMode = ModeStub
			m.ResolvedSpec.Profile = "bool"
		}},
		{"mechanical mode, stub profile", func(m *RunMeta) {
			m.ExecutionMode = ModeMechanical
			m.ResolvedSpec.Profile = "stub"
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			meta := testMeta(ModeMechanical)
			tc.mutate(&meta)
			ran := false
			s := BoolSpec(true)
			s.Alpha = recordingAlpha{ran: &ran}
			if _, err := RunEpisode(context.Background(), s, meta,
				WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"})); err == nil {
				t.Fatal("invalid meta must error at ingress")
			}
			if ran {
				t.Fatal("alpha ran despite invalid meta")
			}
		})
	}
}

// D1: promoting a stub record to mechanical (laundering simulated into
// accepted authority) fails closed even with a recomputed digest — the
// record replays profile/mode coherence at scope lift.
func TestStubPromotionFailsClosed(t *testing.T) {
	cl, err := RunEpisode(context.Background(), BoolSpec(true), testMeta(ModeStub),
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != Simulated {
		t.Fatalf("precondition: want simulated, got %q", cl.Status)
	}
	cl = roundTrip(t, cl)
	cl.Receipt.Record.Mode = ModeMechanical
	cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
	cl.Status = Accepted
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeStub), cl); err == nil {
		t.Fatal("promoted stub record with recomputed digest verified")
	}
	v := validate(BoolSpec(true).Contract, cl.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("profile/mode incoherence must be an integrity failure, got %+v", v)
	}
}

// Regression pairs for Pi β msg-cn-pi-cnos-pr718-fido-round6-beta-47 (D1).

// D1: V's trusted contract arrives as an argument, never out of the receipt.
// Weakening the embedded snapshot and honestly recomputing digest + tail must
// fail verification against the parent's contract.
func TestSubstitutedContractFailsAgainstExpected(t *testing.T) {
	original := BoolSpec(true).Contract
	cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
	if err := VerifyClosure(original, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("honest closure must verify against its own contract: %v", err)
	}

	// Attacker weakens the embedded contract (drops the required evidence),
	// then re-derives the ENTIRE tail honestly against the substituted record.
	weak := cl
	weak.Receipt.Record.Contract.RequiredEvidence = nil
	weak.Receipt.ScopeLiftDigest = sha256hex(weak.Receipt.Record.canonicalBytes())
	weak.Verdict = validate(weak.Receipt.Record.Contract, weak.Receipt)
	weak.Decision = decide(weak.Receipt, weak.Verdict)
	st, err := lift(weak.Verdict, weak.Decision, weak.Receipt.Record.Mode)
	if err != nil {
		t.Fatalf("lift: %v", err)
	}
	weak.Status = st
	weak.Repair = repairFrom(weak.Verdict, st)

	if err := VerifyClosure(original, testMeta(ModeMechanical), weak); err == nil {
		t.Fatal("substituted embedded contract verified against the original contract")
	}
	v := validate(original, weak.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("contract substitution must be an integrity failure, got %+v", v)
	}
}

// D1: the CCNF role split — beta reviews matter and can pass while V alone
// holds the evidence contract: a missing required artifact yields
// contract_unmet -> needs_repair with review.pass still true.
func TestRoleSplitBetaPassesVCatchesUnmet(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g",
			RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
		Alpha: NoopAlpha{}, // valid matter, no artifacts
		Beta:  AcceptBeta{},
	}
	cl := mechClosure(t, s)
	if !cl.Receipt.Record.Review.Pass {
		t.Fatal("precondition: beta's matter review must pass")
	}
	if cl.Status != NeedsRepair {
		t.Fatalf("V must catch the unmet evidence contract: want needs_repair, got %q", cl.Status)
	}
	for _, f := range cl.Verdict.Failures {
		if f.Class.integrity() {
			t.Fatalf("unmet evidence must be contract_unmet, got %+v", f)
		}
	}
}

// D2 (round 6): a nil required array surviving into canonical JSON (null) must
// not self-verify, even with a recomputed digest — CUE rejects it, so must Go.
func TestNullRequiredArraysFailAtScopeLift(t *testing.T) {
	muts := map[string]func(*EpisodeRecord){
		"alpha skills null":    func(r *EpisodeRecord) { r.ResolvedSpec.AlphaSkills = nil },
		"beta skills null":     func(r *EpisodeRecord) { r.ResolvedSpec.BetaSkills = nil },
		"alpha artifacts null": func(r *EpisodeRecord) { r.Alpha.Artifacts = nil },
		"beta artifacts null":  func(r *EpisodeRecord) { r.Beta.Artifacts = nil },
	}
	for name, mut := range muts {
		t.Run(name, func(t *testing.T) {
			cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
			mut(&cl.Receipt.Record)
			cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
			if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err == nil {
				t.Fatalf("%s: null array with recomputed digest verified", name)
			}
			v := validate(BoolSpec(true).Contract, cl.Receipt)
			if v.Pass || !v.hasIntegrityFailure() {
				t.Fatalf("%s: null array must be an integrity failure, got %+v", name, v)
			}
		})
	}
}
