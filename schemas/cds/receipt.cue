// Schema ID: cnos.cdd.cds.receipt.v1
//
// Software-development protocol overlay receipt. Imports the generic
// recursive-cell kernel from schemas/cdd/ and adds the CDS-specific
// required evidence-ref keys via CUE unification.
//
// The CDS receipt requires the closure-record evidence set #369 shipped:
// the five per-cycle artifact refs (self_coherence, beta_review,
// alpha_closeout, beta_closeout, gamma_closeout), the evidence-root ref,
// the cycle's diff, and optional CI refs. These are the artifacts every
// CDS (engineering) cycle emits at γ close-out; V dereferences them when
// validating a CDS receipt.
//
// **Forward-looking debt.** ROLES.md §4a.3 names a CDS receipt sketch
// using forward-looking field names (`artifact_refs`, `test_refs`,
// `ci_refs`, `diff_ref`, `debt_refs`). Cycle/388 implements the field set
// #369 actually shipped (the closure-record shape). A future cycle may
// align §4a.3 and `schemas/cds/receipt.cue` if the field names should
// converge; the cycle/388 close-out names this as known debt.
//
// CUE inheritance idiom: `#CDSReceipt: cdd.#Receipt & { ... }`. The
// generic Receipt's structural invariants (transmissibility derivation,
// override-iff, required boundary_decision, protocol_gap consistency) are
// inherited via unification; #CDSReceipt only constrains `protocol_id` and
// the keys required in `evidence_refs`.

package cds

import "cnos.dev/cnos/schemas/cdd"

#CDSReceipt: cdd.#Receipt & {
	// Pinned protocol_id. V (Phase 3) dispatches on this string.
	protocol_id: "cnos.cdd.cds.receipt.v1"

	// CDS-required evidence refs. Each key is the closure-record artifact
	// the CDS cycle's γ emits at close-out (see CDD.md §close-out and the
	// per-cycle .cdd/releases/<version>/<cycle>/ corpus). V dereferences
	// each ref against the cycle's evidence root.
	evidence_refs: {
		// Root of the cycle evidence corpus (typically
		// ".cdd/releases/<version>/<cycle>/" or ".cdd/unreleased/<cycle>/").
		evidence_root: cdd.#EvidenceRef

		// The five per-cycle artifacts CDS emits.
		self_coherence: cdd.#EvidenceRef
		beta_review:    cdd.#EvidenceRef
		alpha_closeout: cdd.#EvidenceRef
		beta_closeout:  cdd.#EvidenceRef
		gamma_closeout: cdd.#EvidenceRef

		// The cycle's diff against base SHA.
		diff: cdd.#EvidenceRef

		// Optional CI artefact refs (list — a cycle may have many CI jobs).
		ci?: [...cdd.#EvidenceRef]
	}
}
