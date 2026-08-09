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
		{"unknown profile", func(m *RunMeta) { m.ResolvedSpec.Profile = "wat" }},
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
			if _, err := RunEpisode(context.Background(), s,
				WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}), WithMeta(meta)); err == nil {
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
	cl, err := RunEpisode(context.Background(), BoolSpec(true),
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}), WithMeta(testMeta(ModeStub)))
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
	if err := VerifyClosure(cl); err == nil {
		t.Fatal("promoted stub record with recomputed digest verified")
	}
	v := validate(cl.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("profile/mode incoherence must be an integrity failure, got %+v", v)
	}
}

// D1: a record rewritten to a profile outside the closed set (which CUE
// #EpisodeClosure would reject) must not self-verify either.
func TestUnknownProfileFailsAtScopeLift(t *testing.T) {
	cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
	cl.Receipt.Record.ResolvedSpec.Profile = "wat"
	cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
	if err := VerifyClosure(cl); err == nil {
		t.Fatal("unknown profile with recomputed digest verified")
	}
	v := validate(cl.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("unknown profile must be an integrity failure, got %+v", v)
	}
}
