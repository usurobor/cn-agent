package cellkernel

import (
	"context"
	"encoding/json"
	"testing"
)

// Regression pairs for Pi β msg-cn-pi-cnos-pr718-fido-beta-45 (D1–D5).

// D1: episode_id aliasing a station id is an integrity failure — even when the
// attacker recomputes the digest, the closure cannot verify with its original
// accepted tail; and a fully re-derived tail is Rejected, never accepted.
func TestEpisodeIDAliasFailsAtScopeLift(t *testing.T) {
	cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
	cl.Receipt.Record.EpisodeID = cl.Receipt.Record.Alpha.ExecutionID
	cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err == nil {
		t.Fatal("aliased episode id with recomputed digest verified")
	}
	// Even the honest re-derivation of the tampered record is not accepted.
	v := validate(BoolSpec(true).Contract, cl.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("aliased identity must be an integrity failure, got %+v", v)
	}
}

// D1: a required-evidence producer outside alpha|beta fails closed — it is an
// integrity failure and never resolves to the alpha side.
func TestInvalidProducerFailsClosed(t *testing.T) {
	cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
	cl.Receipt.Record.Contract.RequiredEvidence = []RequiredRef{{ID: "bool", Kind: "value", Producer: "gamma"}}
	cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err == nil {
		t.Fatal("invalid producer with recomputed digest verified")
	}
	v := validate(BoolSpec(true).Contract, cl.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("invalid producer must fail closed as integrity, got %+v", v)
	}
}

// D1: a structurally invalid artifact (bad encoding) is caught at scope lift
// even with a recomputed digest.
func TestBadEncodingFailsAtScopeLift(t *testing.T) {
	cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
	cl.Receipt.Record.Alpha.Artifacts[0].Encoding = "base64"
	cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
	v := validate(BoolSpec(true).Contract, cl.Receipt)
	if v.Pass || !v.hasIntegrityFailure() {
		t.Fatalf("unknown encoding must be an integrity failure, got %+v", v)
	}
}

// D2: the repair surface derives from the verdict; changing only repair fields
// fails verification.
func TestRepairTamperFails(t *testing.T) {
	cl := roundTrip(t, mechClosure(t, BoolSpec(false)))
	if cl.Status != NeedsRepair || cl.Repair == nil {
		t.Fatalf("precondition: want needs_repair with repair, got %q", cl.Status)
	}
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("untouched repair closure must verify: %v", err)
	}
	mut := cl
	mut.Repair = &RepairRequest{Reason: "rewritten", Failed: cl.Repair.Failed}
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), mut); err == nil {
		t.Fatal("rewritten repair.reason verified")
	}
	mut2 := cl
	mut2.Repair = &RepairRequest{Reason: cl.Repair.Reason, Failed: []string{"dropped"}}
	if err := VerifyClosure(BoolSpec(true).Contract, testMeta(ModeMechanical), mut2); err == nil {
		t.Fatal("rewritten repair.failed verified")
	}
}

// D3: a hostile alpha that captured the caller's declaration bytes cannot
// mutate runtime-owned invocation truth — the record binds the frozen ingress
// copy.
type declMutatingAlpha struct{ captured []byte }

func (a declMutatingAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	copy(a.captured, []byte(`{"fill":"corrupted"`))
	return AlphaOutput{Matter: Matter{Data: "true"}, Artifacts: []ArtifactCandidate{{ID: "bool", Kind: "value", Text: "true"}}}, nil
}

func TestHostileAlphaCannotMutateResolvedSpec(t *testing.T) {
	decl := []byte(`{"fill":"t.alpha","language":"go"}`)
	meta := testMeta(ModeMechanical)
	meta.ResolvedSpec.Alpha = decl // caller's bytes, also captured by alpha

	s := BoolSpec(true)
	s.Alpha = declMutatingAlpha{captured: decl}
	cl, err := RunEpisode(context.Background(), s, meta,
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if got := string(cl.Receipt.Record.ResolvedSpec.Alpha); got != `{"fill":"t.alpha","language":"go"}` {
		t.Fatalf("frozen resolved spec mutated through captured alias: %q", got)
	}
	// The parent's trusted meta is what it originally supplied — fresh bytes,
	// not the buffer alpha corrupted through the captured alias.
	trusted := testMeta(ModeMechanical)
	trusted.ResolvedSpec.Alpha = json.RawMessage(`{"fill":"t.alpha","language":"go"}`)
	if err := VerifyClosure(BoolSpec(true).Contract, trusted, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}

// D4: stub mode never masks failures as simulated — integrity fails closed and
// semantic failures keep their disposition.
type dupArtifactAlpha struct{}

func (dupArtifactAlpha) Produce(context.Context, AlphaInput) (AlphaOutput, error) {
	return AlphaOutput{Matter: Matter{Data: "m"}, Artifacts: []ArtifactCandidate{
		{ID: "x", Kind: "k", Text: "a"},
		{ID: "x", Kind: "k", Text: "b"},
	}}, nil
}

func TestStubDoesNotMaskFailures(t *testing.T) {
	t.Run("integrity -> rejected", func(t *testing.T) {
		s := Spec{Contract: Contract{ID: "c", Goal: "g"}, Alpha: dupArtifactAlpha{}, Beta: AcceptBeta{}}
		cl, err := RunEpisode(context.Background(), s, testMeta(ModeStub),
			WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		if cl.Status != Rejected {
			t.Fatalf("stub integrity failure: want rejected, got %q", cl.Status)
		}
	})
	t.Run("contract-unmet -> needs_repair", func(t *testing.T) {
		s := Spec{
			Contract: Contract{ID: "c", Goal: "g",
				RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
			Alpha: NoopAlpha{}, // produces nothing — requirement unmet
			Beta:  AcceptBeta{},
		}
		cl, err := RunEpisode(context.Background(), s, testMeta(ModeStub),
			WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		if cl.Status != NeedsRepair {
			t.Fatalf("stub contract-unmet: want needs_repair, got %q", cl.Status)
		}
	})
}

// D5: a nil or typed-nil identity source fails before alpha runs — no panic.
type typedNilIDs struct{}

func (*typedNilIDs) Mint() (Identity, error) { return Identity{}, nil }

func TestNilIDSourceFailsClosed(t *testing.T) {
	if _, err := RunEpisode(context.Background(), BoolSpec(true), testMeta(ModeMechanical), WithIDSource(nil)); err == nil {
		t.Fatal("nil id source must error, not panic")
	}
	var tn *typedNilIDs
	if _, err := RunEpisode(context.Background(), BoolSpec(true), testMeta(ModeMechanical), WithIDSource(tn)); err == nil {
		t.Fatal("typed-nil id source must error, not panic")
	}
}
