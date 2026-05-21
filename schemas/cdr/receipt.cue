// Schema ID: cnos.cdd.cdr.receipt.v1
//
// Research protocol overlay receipt. Imports the generic recursive-cell
// kernel from schemas/cdd/ and adds the CDR-specific required
// evidence-ref keys via CUE unification.
//
// The CDR receipt is anchored in `ROLES.md §4a.3` — research evidence
// discipline requires claim refs, data refs, method refs, result refs, a
// claim-status enum, optional limitations list, and an optional
// reproduction record (β re-ran the producing commands, output matched).
// The empirical anchor is `usurobor/cph` (the research repo): cph's
// `.cdr/waves/<wave>/receipt.md` exemplifies the field set.
//
// **Skeleton, not frozen.** Cycle/388 ships the minimum CDR field set per
// §4a.3 sketch. cnos#376 Sub 1/Sub 2 may refine field types or
// cardinalities; that refinement is a known follow-on (named in
// cycle/388's gamma-closeout as Phase 2.5.1 / Sub 1 precondition).
//
// CUE inheritance idiom: `#CDRReceipt: cdd.#Receipt & { ... }`. The
// generic Receipt's structural invariants are inherited; #CDRReceipt only
// constrains `protocol_id`, the keys required in `evidence_refs`, and
// the CDR-specific top-level fields (claim_status, limitations,
// reproduction).

package cdr

import "cnos.dev/cnos/schemas/cdd"

// Claim status enumerates the calibration of the receipt's asserted
// claims. Distinguishes empirical observation (data-backed) from computed
// derivation (script-output) from inferred reasoning (model/extrapolation)
// from hypothesized (proposed; not yet measured) from indeterminate (data
// inconclusive). The taxonomy is from ROLES.md §4a.3.
#ClaimStatus: "observed" | "computed" | "inferred" | "hypothesized" | "indeterminate"

// Reproduction record. β (the research reviewer) re-ran the producing
// commands; the output matches what the receipt's `result` refs point at.
// Required when claim_status is "observed" or "computed" (i.e. anything
// data-backed); optional for purely-inferred or hypothesized claims where
// no reproduction step exists.
#Reproduction: {
	reviewer_rerun: bool   // β re-ran the producing commands
	command:        string // the canonical command run (with args/SHA)
	output_match:   bool   // observed output matches receipt-claimed output
	notes?:         string // optional reviewer notes
}

#CDRReceipt: cdd.#Receipt & {
	// Pinned protocol_id. V (Phase 3) dispatches on this string.
	protocol_id: "cnos.cdd.cdr.receipt.v1"

	// CDR-required evidence refs. Per ROLES.md §4a.3 sketch: claim refs
	// (which claims this receipt asserts), data refs (dataset / mount /
	// manifest / checksum), method refs (script paths + commit SHA),
	// result refs (output file paths). Each is a list with at least one
	// entry because a research receipt that asserts nothing, references
	// no data, runs no method, or produces no result is not a research
	// receipt — it is an incomplete artifact. The `& [_, ...]` constraint
	// pins minimum-length-one so `cue vet` rejects empty lists structurally
	// (ROLES.md §4a.3: "no `data_refs` → V rejects research receipt").
	evidence_refs: {
		claim:  [_, ...cdd.#EvidenceRef] & [...cdd.#EvidenceRef] // ≥1 claim
		data:   [_, ...cdd.#EvidenceRef] & [...cdd.#EvidenceRef] // ≥1 data source
		method: [_, ...cdd.#EvidenceRef] & [...cdd.#EvidenceRef] // ≥1 method
		result: [_, ...cdd.#EvidenceRef] & [...cdd.#EvidenceRef] // ≥1 result
	}

	// CDR-specific top-level fields. claim_status is required; limitations
	// and reproduction are optional (a hypothesized-status receipt has no
	// reproduction record; an observed/computed receipt should have one).
	claim_status: #ClaimStatus
	limitations?: [...string]
	reproduction?: #Reproduction
}
