// Package cdsadmit is the door: it decides whether a run input is executable
// under the CDS profile, and it is the only thing that decides it. Reached
// through cellfill.Registry.Door, wired by internal/cellfills. Envelope and
// payloads are decided TOGETHER, because a wrong `kind` is this profile's rule
// as much as an issue's shape is, and must refuse through the same receipt as a
// design with no approach rather than a second path with its own exit code.
//
// Everything here is structural, the cognitive arm is DEFERRED in admit, and
// nothing cognitive runs before admission passes — by shape, not by discipline:
// Decide is a pure function of the document's bytes, takes no provider, builds
// no seat, and runs before the registry's constructors (witness door_test.go).
// No IO, no paths (eng/go §2.17); its one link, cellwork.ParseSubject, is one
// decoder per fact, not §2.17's parallel parser.
package cdsadmit

import (
	"encoding/json"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cdsdesign"
	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellinput"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// ReceiptKind is the admission receipt's own tag: a refusal is NOT an episode
// closure, since no episode exists before a contract is admitted (§4.7).
const ReceiptKind = "cnos.cds.admission-receipt.v0"

// Outcome is the complete admission vocabulary (CELL-SYSTEM-DESIGN §10.3).
type Outcome string

const (
	OutcomeAdmitted Outcome = "admitted"
	OutcomeRejected Outcome = "rejected"
	// Incomplete is distinct from rejected because the repairs differ: supply the
	// missing document, versus fix the one you supplied. One case in 0.1.
	OutcomeIncomplete Outcome = "incomplete"
)

// Receipt is the typed record of the decision. InputDigest is the SHA-256 of the
// document decided on: a refusal freezes nothing and closes no episode, so the
// digest alone says WHICH one it refused — not a second proof surface for the
// admitted payloads (CDS-CELL-MIGRATION.md).
type Receipt struct {
	Kind        string  `json:"kind"`
	Outcome     Outcome `json:"outcome"`
	InputDigest string  `json:"input_digest"`
	Reason      string  `json:"reason,omitempty"`
	// SemanticAdequacy names who decided this issue and design are EXECUTABLE, not
	// merely well-formed (Pi #81 C2): no rule here can tell whether the criteria
	// cover the problem, and 0.1 rents no cognition to answer that.
	SemanticAdequacy string `json:"semantic_adequacy"`
}

// SemanticAdequacyOperatorAttested is the only value this profile emits — a
// constant, not a bool, so it reads as a sentence to someone weighing a receipt.
const SemanticAdequacyOperatorAttested = "operator-attested; this cell validated structure only"

// Door is this profile's cellfill.Door. The receipt is marshaled HERE and the
// runner emits those bytes undecoded, so its vocabulary stays this package's.
func Door(raw []byte) (cellfill.Admitted, json.RawMessage, error) {
	admitted, receipt, decision := Decide(raw)
	out, err := json.Marshal(receipt)
	if err != nil {
		// Failing to marshal this package's own struct of strings is a broken
		// runtime, not a refusal, so it must not be reported as one.
		return cellfill.Admitted{}, nil, fmt.Errorf("cds admission: encode receipt: %w", err)
	}
	return admitted, out, decision
}

// Decide returns the exact authored bytes of each payload. The SUBJECT is still
// the authored one and may name `HEAD` or a branch: pinning reads a repository
// and this function is pure, so it happens once afterwards in cellwork.Pin.
func Decide(raw []byte) (cellfill.Admitted, Receipt, error) {
	digest := cellinput.Digest(raw)
	in, err := cellinput.Decode(raw)
	if err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}
	return admit(digest, in)
}

// admit decides one already-decoded run input. The receipt is what an operator
// reads and the error what a Go caller handles; there is no third "fault"
// channel, because every rule below is a pure predicate that cannot malfunction.
func admit(digest string, in cellinput.RunInput) (cellfill.Admitted, Receipt, error) {
	// Absent and malformed are different outcomes: an absent payload is a run input
	// never finished, a malformed one a document whose author believed it was.
	for _, p := range []struct {
		name string
		raw  json.RawMessage
	}{
		{"issue", in.Issue},
		{"design", in.Design},
		{"subject", in.Subject},
	} {
		if len(p.raw) == 0 {
			return refuse(digest, OutcomeIncomplete, fmt.Sprintf("run input carries no %s", p.name))
		}
		// THE DOOR REFUSES WHAT THE KERNEL WILL REFUSE: without this check both seats
		// were built before the kernel rejected the oversize slot as `episode
		// malfunction` (CDS-CELL-MIGRATION.md). The bound is REFERENCED, never
		// restated: a second constant drifts silently into admitting what it rejects.
		if len(p.raw) > cellkernel.MaxOpaqueSlotBytes {
			return refuse(digest, OutcomeRejected, fmt.Sprintf(
				"run input %s is %d bytes, over the %d-byte contract slot bound",
				p.name, len(p.raw), cellkernel.MaxOpaqueSlotBytes))
		}
	}

	iss, err := cdsissue.Admit(in.Issue)
	if err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}
	des, err := cdsdesign.Admit(in.Design)
	if err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}
	// ParseSubject, not AdmitSubject: an AUTHORED subject may name a moving
	// revision, which is what pinning resolves. Requiring 40 hex here would make
	// `HEAD` inadmissible and move pinning into the author's hands.
	if _, err := cellwork.ParseSubject(in.Subject); err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}
	if err := relate(iss, des); err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}

	// --- cognitive arm: DEFERRED ------------------------------------------
	// "Is the problem real, are issue and design consistent, is the scope
	// executable" (CELL-SYSTEM-DESIGN §10.2 step 6) is declared and EMPTY: deferred
	// for order, not value, at the cost in CDS-CELL-MIGRATION.md. A branch and not
	// an unfilled interface field, since a nil attestor would read as a gate that ran.

	return cellfill.Admitted{Issue: in.Issue, Design: in.Design, Subject: in.Subject},
		Receipt{Kind: ReceiptKind, Outcome: OutcomeAdmitted, InputDigest: digest,
			SemanticAdequacy: SemanticAdequacyOperatorAttested}, nil
}

// relate is the cross-facet arm: rules about issue and design TOGETHER. Both
// conjuncts are ALSO enforced in cdsissue.Admit and cdsdesign.Admit, so nothing on
// the full path can reach here and fail it; written anyway because the relation is
// the door's to own if a facet relaxes its rule, and tested by direct call only.
func relate(iss cdsissue.Issue, des cdsdesign.Design) error {
	seen := make(map[string]bool, len(iss.Acceptance))
	for _, c := range iss.Acceptance {
		if seen[c.ID] {
			return fmt.Errorf("cds admission: acceptance id %q is not unique", c.ID)
		}
		seen[c.ID] = true
	}
	for i, s := range des.Impact {
		if cdsissue.Blank(s.Surface) {
			return fmt.Errorf("cds admission: design impact[%d] names no surface", i)
		}
	}
	return nil
}

// refuse states one refusal in all three channels: nothing frozen, a typed receipt,
// and an ErrRefused wrap so a caller ignoring the receipt still cannot proceed.
func refuse(digest string, o Outcome, reason string) (cellfill.Admitted, Receipt, error) {
	return cellfill.Admitted{},
		Receipt{Kind: ReceiptKind, Outcome: o, InputDigest: digest, Reason: reason,
			SemanticAdequacy: SemanticAdequacyOperatorAttested},
		fmt.Errorf("%w (%s): %s", cellfill.ErrRefused, o, reason)
}
