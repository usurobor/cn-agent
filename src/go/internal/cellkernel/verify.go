package cellkernel

import (
	"encoding/json"
	"errors"
	"strings"
)

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

// VerifyClosure is the ONE scope-lift verification boundary (FIDO doctrine,
// msg-cn-pi-cnos-cell-runner-fido-functional-44): a parent or another runtime
// re-checks a serialized closure against the contract AND invocation metadata
// IT trusts (Pi round-6 D1, round-7 D1). Neither ever comes out of the
// closure being judged, so a substituted embedded contract — or a coherent
// mode+profile rewrite (stub/simulated laundered to mechanical/accepted) —
// cannot verify even with an honestly recomputed digest and tail. The parent
// invoked the episode, so it owns both values; this is ordinary argument
// passing, not a second proof surface.
//
//  1. the scope-lift digest recomputes over the record's canonical bytes;
//  2. the record's frozen contract, execution mode, and full resolved spec
//     equal the trusted expected values;
//  3. verdict ← V(expected, receipt), decision ← δ(receipt, verdict),
//     status ← lift(...) all re-derive to exactly the closure's values;
//  4. schema/protocol pins hold and repair is present iff status needs it.
//
// Returns nil iff the closure is fully self-consistent against the expected
// contract and metadata.
func VerifyClosure(expected Contract, meta RunMeta, cl Closure) error {
	var f []string
	add := func(cond bool, msg string) {
		if cond {
			f = append(f, msg)
		}
	}

	add(cl.Schema != ClosureSchema, "wrong closure schema")
	add(cl.ProtocolValidated, "protocol_validated must be false in v0")

	// (1) the single reproducible proof.
	add(sha256hex(cl.Receipt.Record.canonicalBytes()) != cl.Receipt.ScopeLiftDigest,
		"scope-lift digest does not recompute")

	// (2) invocation authority: mode and the full resolved spec bind to the
	// parent-trusted metadata, canonicalized exactly as the runtime froze it
	// at ingress.
	rec := cl.Receipt.Record
	add(rec.Mode != meta.ExecutionMode, "record execution mode does not match the expected metadata")
	add(string(mustJSON(rec.ResolvedSpec)) != string(mustJSON(meta.ResolvedSpec.clone())),
		"record resolved spec does not match the expected metadata")

	// (3) pure re-derivation of the mechanical tail — including the repair
	// surface, so repair is never a second unauthenticated authority (D2).
	wantVerdict := validate(expected, cl.Receipt)
	add(string(mustJSON(wantVerdict)) != string(mustJSON(cl.Verdict)), "verdict does not derive from the receipt")
	wantDecision := decide(cl.Receipt, wantVerdict)
	add(wantDecision != cl.Decision, "decision does not derive from the verdict")
	if wantStatus, err := lift(wantVerdict, wantDecision, rec.Mode); err != nil {
		f = append(f, "status inconsistent: "+err.Error())
	} else {
		add(wantStatus != cl.Status, "status does not derive from (verdict, decision, mode)")
		wantRepair := repairFrom(wantVerdict, wantStatus)
		add(string(mustJSON(wantRepair)) != string(mustJSON(cl.Repair)), "repair does not derive from the verdict")
	}

	if len(f) > 0 {
		return errors.New("closure verification failed: " + strings.Join(f, "; "))
	}
	return nil
}
