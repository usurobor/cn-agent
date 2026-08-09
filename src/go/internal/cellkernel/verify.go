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
// re-checks a serialized closure with a single reproducible proof plus pure
// re-derivation — no overlapping hash authorities.
//
//  1. the scope-lift digest recomputes over the record's canonical bytes;
//  2. verdict ← V(receipt), decision ← δ(verdict), status ← lift(...) all
//     re-derive to exactly the closure's values;
//  3. schema/protocol pins hold and repair is present iff status needs it.
//
// Returns nil iff the closure is fully self-consistent.
func VerifyClosure(cl Closure) error {
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

	// (2) pure re-derivation of the mechanical tail.
	wantVerdict := validate(cl.Receipt)
	add(string(mustJSON(wantVerdict)) != string(mustJSON(cl.Verdict)), "verdict does not derive from the receipt")
	wantDecision := decide(wantVerdict)
	add(wantDecision != cl.Decision, "decision does not derive from the verdict")
	if wantStatus, err := lift(wantVerdict, wantDecision, cl.Receipt.Record.Mode); err != nil {
		f = append(f, "status inconsistent: "+err.Error())
	} else {
		add(wantStatus != cl.Status, "status does not derive from (verdict, decision, mode)")
	}

	// (3) coherence of the repair surface.
	add((cl.Status == NeedsRepair) != (cl.Repair != nil), "repair present iff status is needs_repair")

	if len(f) > 0 {
		return errors.New("closure verification failed: " + strings.Join(f, "; "))
	}
	return nil
}
