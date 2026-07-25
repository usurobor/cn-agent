// Package stub provides deterministic, external-call-free actors so the loop
// runs end-to-end and proves the MACHINERY (not the intelligence). Each cell's
// telos is expressed only by which actor is bound; all three drive the same
// kernel. Agent-backed actors implement kernel.Actor next (v1) with no kernel
// change. Spike code; see ../../DESIGN.md.
package stub

import (
	"regexp"
	"strings"

	"cnos.dev/cnos/cell-runner/kernel"
	"cnos.dev/cnos/cell-runner/model"
)

// CoherenceHeading is the required section the toy target must carry.
const CoherenceHeading = "## Coherence"

// coherenceBody is the section body α appends (the WC's artifact edit).
const coherenceBody = "This artifact declares its coherence: every required section is present\n" +
	"and mutually consistent with the target's stated telos.\n"

var coherenceRe = regexp.MustCompile(`(?m)^##\s+Coherence\s*$`)

// hasCoherenceSection reports whether content contains a `## Coherence` heading.
func hasCoherenceSection(content string) bool { return coherenceRe.MatchString(content) }

// coherenceSectionWellFormed reports whether the `## Coherence` heading is
// present AND followed by a non-empty body (before the next `## ` heading or
// EOF). This is the mechanical acceptance predicate V asserts for a WC.
func coherenceSectionWellFormed(content string) bool {
	loc := coherenceRe.FindStringIndex(content)
	if loc == nil {
		return false
	}
	rest := content[loc[1]:]
	for _, line := range strings.Split(rest, "\n") {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "## ") {
			break // reached the next section with no body in between
		}
		return true // a non-blank body line
	}
	return false
}

// baseActor supplies the telos-agnostic steps every cell shares: γ close, δ
// accept, and δ hold. The three concrete actors embed it and override only the
// telos-specific steps (Alpha, Beta, Validate; CC also overrides DeltaAccept
// to emit its judgment).
type baseActor struct{}

// Gamma closes the cell: it computes the receipt skeleton from the transition
// evidence gathered so far and binds it. The verdict is filled by V; the
// boundary decision by δ. The receipt is COMPUTED here, never authored.
func (baseActor) Gamma(rc *kernel.RunContext) error {
	ensureReceipt(rc)
	rc.Facts.ReceiptBound = true
	rc.AddEvidence("gamma_close", rc.Ref+"/gamma")
	rc.Log("  γ close: receipt bound from %d transition(s) of evidence", len(rc.Transitions))
	return nil
}

// DeltaAccept records the accept boundary decision and finalizes the receipt.
func (baseActor) DeltaAccept(rc *kernel.RunContext) error { return acceptReceipt(rc) }

// DeltaHold finalizes a STOP receipt with a reject boundary decision (fail
// closed). It is reached on β-reject (from `reviewed`) or V-fail (from
// `validated`); it distinguishes the two by the edge that led here.
func (baseActor) DeltaHold(rc *kernel.RunContext) error { return holdReceipt(rc) }

// ensureReceipt lazily creates the receipt (the β-reject path skips γ, so δ
// must be able to compute the receipt itself).
func ensureReceipt(rc *kernel.RunContext) {
	if rc.Receipt == nil {
		rc.Receipt = &model.Receipt{
			ID:            rc.Cell,
			Cell:          rc.Cell,
			Telos:         rc.Telos,
			ProducedAtRef: rc.Ref,
		}
	}
}

// setValidation records V's verdict into the receipt and sets the VPass/VFail
// fact the kernel table reads next.
func setValidation(rc *kernel.RunContext, pass bool, failed []string) {
	ensureReceipt(rc)
	if pass {
		rc.Facts.VPass = true
		rc.Receipt.Validation = model.ValidationVerdict{Verdict: "PASS", FailedPredicates: []string{}, Warnings: []string{}}
	} else {
		rc.Facts.VFail = true
		rc.Receipt.Validation = model.ValidationVerdict{Verdict: "FAIL", FailedPredicates: nonNil(failed), Warnings: []string{}}
	}
}

// acceptReceipt finalizes the pass-path receipt.
func acceptReceipt(rc *kernel.RunContext) error {
	ensureReceipt(rc)
	rc.Receipt.Transitions = snapshotTransitions(rc)
	rc.Receipt.EvidenceRefs = snapshotEvidence(rc)
	rc.Receipt.BoundaryDecision = model.BoundaryDecision{Actor: "delta", Action: "accept"}
	rc.Receipt.Verdict = "accept"
	rc.Receipt.Transmissibility = model.DeriveTransmissibility(rc.Receipt.Validation.Verdict, "accept")
	rc.Log("  δ accept: boundary=accept transmissibility=%s", rc.Receipt.Transmissibility)
	return nil
}

// holdReceipt finalizes the fail-path (STOP) receipt.
func holdReceipt(rc *kernel.RunContext) error {
	ensureReceipt(rc)
	// If V never ran (β-reject skips γ/V), the failing predicate is the review
	// itself; record it so the receipt is a complete evidence surface.
	if rc.Receipt.Validation.Verdict == "" {
		rc.Receipt.Validation = model.ValidationVerdict{
			Verdict:          "FAIL",
			FailedPredicates: []string{"beta_review_approved"},
			Warnings:         []string{},
		}
	}
	rc.Receipt.Transitions = snapshotTransitions(rc)
	rc.Receipt.EvidenceRefs = snapshotEvidence(rc)
	rc.Receipt.BoundaryDecision = model.BoundaryDecision{Actor: "delta", Action: "reject"}
	rc.Receipt.Verdict = "hold"
	rc.Receipt.Transmissibility = model.DeriveTransmissibility(rc.Receipt.Validation.Verdict, "reject")
	rc.Log("  δ HOLD: boundary=reject transmissibility=%s (fail closed)", rc.Receipt.Transmissibility)
	return nil
}

func snapshotTransitions(rc *kernel.RunContext) []model.TransitionRecord {
	out := make([]model.TransitionRecord, len(rc.Transitions))
	copy(out, rc.Transitions)
	return out
}

func snapshotEvidence(rc *kernel.RunContext) map[string]string {
	out := make(map[string]string, len(rc.Evidence))
	for k, v := range rc.Evidence {
		out[k] = v
	}
	return out
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
