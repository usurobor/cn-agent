package cellkernel

import (
	"context"
	"encoding/json"
	"testing"
)

// Regression pairs for Pi β msg-cn-pi-cnos-pr718-fido-round7-beta-48.

// D1: a COHERENT rewrite — execution_mode stub→mechanical with digest and
// the entire tail honestly recomputed — must still fail against the
// parent-trusted RunMeta. The record's internal consistency proves nothing
// about invocation authority; only the metadata the parent supplied does.
func TestCoherentModeRewriteFailsAgainstTrustedMeta(t *testing.T) {
	trusted := testMeta(ModeStub)
	cl, err := RunEpisode(context.Background(), BoolSpec(true), trusted,
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != Simulated {
		t.Fatalf("precondition: want simulated, got %q", cl.Status)
	}
	if err := VerifyClosure(BoolSpec(true).Contract, trusted, cl); err != nil {
		t.Fatalf("honest stub closure must verify against its own meta: %v", err)
	}

	// Paired rewrite: mode AND profile move together, so the record stays
	// internally coherent; the attacker re-derives the full tail honestly.
	laundered := roundTrip(t, cl)
	laundered.Receipt.Record.Mode = ModeMechanical
	laundered.Receipt.ScopeLiftDigest = sha256hex(laundered.Receipt.Record.canonicalBytes())
	laundered.Verdict = validate(BoolSpec(true).Contract, laundered.Receipt)
	laundered.Decision = decide(laundered.Receipt, laundered.Verdict)
	st, err := lift(laundered.Verdict, laundered.Decision, laundered.Receipt.Record.Mode)
	if err != nil {
		t.Fatalf("lift: %v", err)
	}
	laundered.Status = st
	laundered.Repair = repairFrom(laundered.Verdict, st)
	if laundered.Status != Accepted {
		t.Fatalf("attack precondition: laundered tail should read accepted, got %q", laundered.Status)
	}

	if err := VerifyClosure(BoolSpec(true).Contract, trusted, laundered); err == nil {
		t.Fatal("dual mode/profile rewrite verified against the original trusted meta")
	}
}

// C2 (re-anchored by fill construction): seat declarations are opaque at the
// generic boundaries — a mechanical episode under an arbitrary fill tag runs
// and verifies; fill whitelists are the registry's business, not the kernel's.
func TestOpaqueDeclarationVerifies(t *testing.T) {
	meta := testMeta(ModeMechanical)
	meta.ResolvedSpec.Alpha = json.RawMessage(`{"fill":"custom.fill","anything":"goes"}`)
	cl, err := RunEpisode(context.Background(), BoolSpec(true), meta,
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != Accepted {
		t.Fatalf("status: want accepted, got %q", cl.Status)
	}
	if err := VerifyClosure(BoolSpec(true).Contract, meta, roundTrip(t, cl)); err != nil {
		t.Fatalf("opaque-declaration closure must verify: %v", err)
	}
}
