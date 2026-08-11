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

	// Paired rewrite: mode AND status move together, so the record stays
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
		t.Fatal("dual mode/status rewrite verified against the original trusted meta")
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

// Pi pr719 C1: the generic boundary validates the one thing it knows about a
// seat declaration — a JSON object with a non-empty string `fill`. Without
// it a direct caller could emit `{}`, `[]`, or a scalar, self-verify it, and
// still fail #EpisodeClosure: the kernel would be certifying something no
// reader can parse. Fill-specific fields stay opaque.
func TestMalformedSeatEnvelopeRejected(t *testing.T) {
	bad := map[string]string{
		"empty object":      `{}`,
		"array":             `[]`,
		"scalar":            `"cds.patch"`,
		"number":            `7`,
		"empty fill":        `{"fill":""}`,
		"fill not a string": `{"fill":{"name":"x"}}`,
		"case-aliased tag":  `{"Fill":"cds.patch"}`,
	}
	for name, decl := range bad {
		t.Run(name, func(t *testing.T) {
			meta := testMeta(ModeMechanical)
			meta.ResolvedSpec.Alpha = json.RawMessage(decl)
			if _, err := RunEpisode(context.Background(), BoolSpec(true), meta,
				WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"})); err == nil {
				t.Fatalf("%s must not run", name)
			}
			cl := roundTrip(t, mechClosure(t, BoolSpec(true)))
			cl.Receipt.Record.ResolvedSpec.Alpha = json.RawMessage(decl)
			cl.Receipt.ScopeLiftDigest = sha256hex(cl.Receipt.Record.canonicalBytes())
			if err := VerifyClosure(BoolSpec(true).Contract, meta, cl); err == nil {
				t.Fatalf("%s must not verify", name)
			}
			v := validate(BoolSpec(true).Contract, cl.Receipt)
			if v.Pass || !v.hasIntegrityFailure() {
				t.Fatalf("%s must be an integrity failure, got %+v", name, v)
			}
		})
	}
}

// The positive half: any well-formed tagged object is accepted whole, with
// its interior opaque to the kernel.
func TestOpaqueTaggedEnvelopeAccepted(t *testing.T) {
	meta := testMeta(ModeMechanical)
	meta.ResolvedSpec.Alpha = json.RawMessage(`{"fill":"anyone.at.all","nested":{"deep":[1,2,{"x":null}]}}`)
	cl, err := RunEpisode(context.Background(), BoolSpec(true), meta,
		WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
	if err != nil {
		t.Fatalf("a tagged declaration must run whatever its interior: %v", err)
	}
	if err := VerifyClosure(BoolSpec(true).Contract, meta, roundTrip(t, cl)); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}
