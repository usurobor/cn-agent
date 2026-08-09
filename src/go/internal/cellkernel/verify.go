package cellkernel

import (
	"errors"
	"fmt"
	"strings"
)

// VerifyReceipt re-verifies a receipt's internal structural integrity from its
// serialized content alone (Pi #33 D2): recompute every hash, re-derive the
// content-addressed evidence, check id uniqueness, that each evidence ref is
// bound to its producer's station id, that execution ids are non-empty and
// distinct, and that beta_input_hash recomputes. It does NOT judge contract
// satisfaction (that is V's verdict) — only that the record is self-consistent.
func VerifyReceipt(rc Receipt) error {
	var f []string
	add := func(cond bool, msg string) {
		if cond {
			f = append(f, msg)
		}
	}

	add(rc.EpisodeID == "", "missing episode_id")
	add(rc.AlphaExecutionID == "" || rc.BetaExecutionID == "", "missing execution id")
	add(rc.AlphaExecutionID == rc.BetaExecutionID, "execution ids not distinct")
	add(rc.PolicyID == "", "missing beta_input_policy_id")

	add(hashJSON(rc.Contract) != rc.ContractHash, "contract hash mismatch")
	add(hashJSON(rc.Matter) != rc.MatterHash, "matter hash mismatch")
	add(hashJSON(rc.Review) != rc.ReviewHash, "review hash mismatch")
	add(hashJSON(rc.Evidence) != rc.EvidenceHash, "evidence hash mismatch")

	seen := make(map[string]int)
	var alphaEv []EvidenceRef
	for _, e := range rc.Evidence {
		seen[e.ID]++
		add(sha256hex([]byte(e.Content)) != e.SHA256, "evidence hash mismatch: "+e.ID)
		add(e.Ref != "sha256:"+e.SHA256, "evidence ref not content-addressed: "+e.ID)
		switch e.Producer {
		case RoleAlpha:
			add(e.ProducerExecutionID != rc.AlphaExecutionID, "alpha evidence not bound to alpha station: "+e.ID)
			alphaEv = append(alphaEv, e)
		case RoleBeta:
			add(e.ProducerExecutionID != rc.BetaExecutionID, "beta evidence not bound to beta station: "+e.ID)
		default:
			f = append(f, "evidence has no producer: "+e.ID)
		}
	}
	for id, n := range seen {
		if n > 1 {
			f = append(f, "duplicate evidence id: "+id)
		}
	}

	// beta_input_hash must recompute from the receipt's own alpha evidence.
	want := hashBetaInput(BetaInput{ContractHash: rc.ContractHash, Matter: rc.Matter, AlphaEvidence: alphaEv, PolicyID: rc.PolicyID})
	add(want != rc.BetaInputHash, "beta_input_hash does not recompute")

	if len(f) > 0 {
		return errors.New("receipt verification failed: " + strings.Join(f, "; "))
	}
	return nil
}

// VerifyEnvelope re-derives EVERY field of a terminal envelope from content
// (Pi PR-#718 β D1/D2): the resolved-spec hash, the inner receipt, and the
// verdict→decision→status chain. A parent that runs this trusts nothing that is
// merely copied in. Returns nil iff the envelope is fully self-consistent.
func VerifyEnvelope(env Envelope) error {
	var f []string
	add := func(cond bool, msg string) {
		if cond {
			f = append(f, msg)
		}
	}

	add(env.Schema != EnvelopeSchema, "wrong envelope schema")
	add(env.ProtocolValidated, "protocol_validated must be false in v0")
	add(!knownMode(env.ExecutionMode), "unknown execution mode")

	// Resolved-spec identity recomputes, and binds the receipt's contract.
	add(env.ResolvedSpec.Canon != resolvedSpecCanon, "wrong resolved-spec canon")
	add(env.ResolvedSpec.hash() != env.ResolvedSpecHash, "resolved_spec_hash does not recompute")
	add(hashJSON(env.ResolvedSpec.Contract) != hashJSON(env.Receipt.Contract), "resolved spec contract != receipt contract")

	// Inner receipt integrity.
	if err := VerifyReceipt(env.Receipt); err != nil {
		f = append(f, err.Error())
	}

	// verdict ← V(receipt); decision ← δ(verdict); status ← (decision, mode).
	wantVerdict := validate(env.Receipt)
	add(hashJSON(wantVerdict) != hashJSON(env.Verdict), "verdict does not derive from the receipt")
	wantDecision := decide(wantVerdict)
	add(wantDecision != env.Decision, "decision does not derive from the verdict")
	if wantStatus, err := statusOf(wantVerdict, wantDecision, env.ExecutionMode); err != nil {
		f = append(f, fmt.Sprintf("status inconsistent: %v", err))
	} else {
		add(wantStatus != env.Status, "status does not derive from (decision, mode)")
	}

	add((env.Status == NeedsRepair) != (env.Repair != nil), "repair present iff status is needs_repair")

	if len(f) > 0 {
		return errors.New("envelope verification failed: " + strings.Join(f, "; "))
	}
	return nil
}
