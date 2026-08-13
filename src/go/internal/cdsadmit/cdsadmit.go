// Package cdsadmit is the door: it decides whether a run input is executable
// under the CDS profile, and it is the only thing that decides it.
//
// Everything here is structural. The question "is this issue REAL, is this
// design consistent with it, are these criteria sufficient" is the cognitive
// arm, and it is deliberately not implemented in 0.1 — see the branch marked
// DEFERRED in admit for why, and what it costs.
//
// THE WHOLE DOCUMENT IS DECIDED HERE, envelope and payloads together. Decide
// takes raw bytes: the declared `kind` and the closed key language are this
// profile's rules exactly as much as an issue's shape is, so a wrong `kind`
// refuses through the same receipt as a design with no approach. While the
// runner decoded the envelope itself, those two questions had two answers —
// one produced a receipt and its own exit code, the other produced neither.
//
// The property this package exists to hold is narrow and mechanical: nothing
// cognitive runs before admission passes. It is held by shape rather than by
// discipline — Decide is a pure function of the document's bytes, it takes no
// provider, constructs no seat, and runs before the fill registry's
// constructors are dispatched. Its witness is internal/cdsadmit/door_test.go,
// which invokes the real runner with a recording registry and fails if a seat
// was ever built.
//
// Decide performs no IO and touches no path (eng/go §2.17): it is a function of
// bytes already in memory. Stated at that precision rather than as "this
// package is pure", because it calls cellwork.ParseSubject — the pure half of
// a package that also contains the git adapter, and therefore links one. One
// decoder per fact is worth that link; a second subject decoder written here
// would be the parallel-parser violation §2.17 names.
//
// It is reached through cellfill.Registry.Door, wired by internal/cellfills.
// The runner names no CDS package.
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

// ReceiptKind is the pinned tag of an admission receipt. An admission refusal
// is NOT an episode closure — no episode exists before a contract is admitted
// (CELL-SYSTEM-DESIGN §4.7) — so it carries its own kind and is never emitted
// under the closure schema. Spelled under the repo's `cnos.<domain>.<name>.vN`
// convention, like every other `cnos.` tag a reader of this record meets.
const ReceiptKind = "cnos.cds.admission-receipt.v0"

// Outcome is the complete admission vocabulary (CELL-SYSTEM-DESIGN §10.3).
type Outcome string

const (
	// OutcomeAdmitted: every structural gate passed. Production may run.
	OutcomeAdmitted Outcome = "admitted"
	// OutcomeRejected: a decisive contract defect exists. The document is
	// present and wrong, and no amount of waiting changes that.
	OutcomeRejected Outcome = "rejected"
	// OutcomeIncomplete: a required part is UNAVAILABLE rather than wrong.
	// Structurally that is exactly one case in 0.1 — an absent payload — and
	// it is kept distinct from rejection because the repairs differ: supply
	// the missing document, versus fix the one you supplied.
	OutcomeIncomplete Outcome = "incomplete"
)

// Receipt is the typed record of the decision.
//
// InputDigest is the SHA-256 of the exact run-input document this receipt
// decided on. It is the receipt's ARTIFACT IDENTITY, and it is what a refusal
// most needs: when a document is refused, nothing is frozen into a contract and
// no episode closure exists, so the receipt is the only record of the run and
// the digest is the only thing that says WHICH document it refused.
//
// It is not a second proof surface for the admitted payloads. Those are frozen
// into the contract and covered there by the kernel's one scope-lift digest,
// which is taken over different bytes — the canonical episode record, not the
// authored envelope. This digest identifies an untrusted input; that one proves
// an episode. Carried on every outcome rather than only on refusals, because a
// receipt that sometimes says what it decided about is harder to read than one
// that always does.
type Receipt struct {
	Kind        string  `json:"kind"`
	Outcome     Outcome `json:"outcome"`
	InputDigest string  `json:"input_digest"`
	Reason      string  `json:"reason,omitempty"`
}

// Door is this profile's cellfill.Door: the whole decision, from raw bytes to a
// serialized receipt, in the form the generic runner dispatches. It is what
// internal/cellfills wires into the registry, exactly as it wires a fill.
//
// The receipt is marshaled HERE. The runner emits those bytes without decoding
// them, so the receipt's kind tag and vocabulary stay this package's and the
// runner never learns a word of CDS.
func Door(raw []byte) (cellfill.Admitted, json.RawMessage, error) {
	admitted, receipt, decision := Decide(raw)
	out, err := json.Marshal(receipt)
	if err != nil {
		// The receipt is this package's own struct of strings; failing to
		// marshal it is a broken runtime, not a refusal, so it must not be
		// reported as one.
		return cellfill.Admitted{}, nil, fmt.Errorf("cds admission: encode receipt: %w", err)
	}
	return admitted, out, decision
}

// Decide is the whole door as the runner walks it: the ENVELOPE and the
// payloads, one decision, one receipt.
//
// The envelope belongs here and not to the caller. A document declaring the
// wrong `kind`, carrying an unknown key, or spelling `Kind` with a capital is
// decisively inadmissible under this profile — the same class of fact as a
// design with no approach — so it refuses through the receipt path and not as
// some other kind of error the caller invents a code for. Two refusal paths for
// one question is how `kind` came to exit differently from `approach`.
//
// What it returns on the admitting path is cellfill.Admitted: the exact
// authored bytes of each payload, in the form the contract will carry them. The
// SUBJECT there is still the authored one — it may name `HEAD`, a branch or a
// tag. Pinning it is an effect (it reads a repository) and this function is
// pure, so resolution happens once afterwards, in cellwork.Pin, before either
// seat is constructed.
func Decide(raw []byte) (cellfill.Admitted, Receipt, error) {
	digest := cellinput.Digest(raw)
	in, err := cellinput.Decode(raw)
	if err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}
	return admit(digest, in)
}

// admit decides one already-decoded run input.
//
// The receipt and the error carry the same fact through two channels on
// purpose: the receipt is the typed result an operator reads, the error is
// what a Go caller must handle. There is no third, "fault" channel in 0.1 —
// every rule below is a pure predicate over bytes already in memory, so there
// is no mechanism here that can malfunction as distinct from refusing.
func admit(digest string, in cellinput.RunInput) (cellfill.Admitted, Receipt, error) {
	// Absent versus malformed versus oversize. All three refuse, and the first
	// two are different outcomes: an absent payload is a run input that was
	// never finished, a malformed one is a document whose author believed it
	// was.
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
		// THE DOOR REFUSES WHAT THE KERNEL WILL REFUSE. Each admitted payload
		// becomes an opaque contract slot, and the kernel bounds every slot at
		// cellkernel.MaxOpaqueSlotBytes. Without this check an oversize issue
		// was admitted, its subject was resolved against a real repository, and
		// BOTH seats were constructed — a provider adapter built and skill
		// bodies loaded — before validateSpec reported `episode malfunction`.
		// Work happened on a document that was never admissible, and the
		// operator was told the runtime broke rather than that their document
		// was too large.
		//
		// The bound is REFERENCED, never restated: one number, owned by the
		// boundary that enforces it. A second constant here would be a second
		// number to keep in step, and the failure it produces is silent — the
		// door would admit exactly the documents the kernel then rejects.
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
	// revision, which is precisely what the pinning step exists to resolve.
	// Requiring 40 hex here would make `HEAD` — the thing a human writes —
	// inadmissible, and would move pinning into the author's hands.
	if _, err := cellwork.ParseSubject(in.Subject); err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}
	if err := relate(iss, des); err != nil {
		return refuse(digest, OutcomeRejected, err.Error())
	}

	// --- cognitive arm: DEFERRED ------------------------------------------
	//
	// This is where "is the problem real, are issue and design mutually
	// consistent, is the scope executable" would be attested
	// (CELL-SYSTEM-DESIGN §10.2 step 6). It is declared and EMPTY: control
	// falls straight through to admitted, nothing is consulted, and no
	// provider is reachable from this function.
	//
	// Deferred, not dropped, and the reason is order rather than value
	// (WCC-0.1-PLAN §0 C2). By the design's own authority table the arm is
	// `attested_unverified`: it enforces nothing and cannot be re-derived, so
	// a run cannot be refused on its say-so. What it would cost now is a third
	// cognitive station, a provider fault class, an attestation vocabulary and
	// a receipt shape — plus one provider round-trip on every run, including
	// the deterministic-fake corpus. Semantic admission is the right long-term
	// answer to "is this issue executable"; it is not on the path to a first
	// accepted patch, and the structural arm above already refuses any
	// deliberately bad issue before a seat exists.
	//
	// The position is a branch and not an unpopulated interface field on
	// purpose: a nil attestor threaded through the signature would read as a
	// gate that ran and found nothing, which is a stronger claim than "no gate
	// ran".

	return cellfill.Admitted{Issue: in.Issue, Design: in.Design, Subject: in.Subject},
		Receipt{Kind: ReceiptKind, Outcome: OutcomeAdmitted, InputDigest: digest}, nil
}

// relate is the cross-facet arm: the rules that are about the issue and the
// design TOGETHER rather than about either alone.
//
// Stated exactly, because the honest description matters more than the
// appearance of coverage: both conjuncts below are ALSO enforced inside
// cdsissue.Admit and cdsdesign.Admit, so on the full Admit path no document
// can reach this function and fail it. It is a redundant guard today. It is
// written anyway because the relation is the door's to own — if either facet
// ever relaxes its own rule, the door does not silently lose it — and it is
// unit-tested directly rather than through Admit, since only a direct call can
// make it fire.
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

// refuse states one refusal in all three channels at once: nothing frozen, a
// typed receipt naming the document and the reason, and an error wrapping
// cellfill.ErrRefused so a caller that ignores the receipt still cannot proceed.
func refuse(digest string, o Outcome, reason string) (cellfill.Admitted, Receipt, error) {
	return cellfill.Admitted{},
		Receipt{Kind: ReceiptKind, Outcome: o, InputDigest: digest, Reason: reason},
		fmt.Errorf("%w (%s): %s", cellfill.ErrRefused, o, reason)
}
