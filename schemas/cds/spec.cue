// Schema ID: cnos.cdd.cds.spec.v1
//
// #CDSCellSpec — the CDS input-side overlay on the generic #CellSpec, mirroring
// how schemas/cds/receipt.cue #CDSReceipt overlays #Receipt on the output side.
// It pins the declared protocol_id and requires the `language` parameter hole.
// The generic cell's structural invariants are inherited via CUE unification.
//
// NOTE (Pi #32 D1): declaring protocol_id "cnos.cdd.cds.receipt.v1" is a
// declaration of intent, not proof. The v0 runner carries it as provenance and
// sets protocol_validated=false; running full CDS protocol/receipt validation
// is a later phase.
//
// See docs/architecture/CDS-CELL-MIGRATION.md (Phase 0).
package cds

import "cnos.dev/cnos/schemas/cdd"

#CDSCellSpec: cdd.#CellSpec & {
	protocol_id: "cnos.cdd.cds.receipt.v1"

	params: {
		language: {kind: "skill", required: true, domain: ["go", "ocaml", "rust", "python", "typescript"]}
		...
	}

	// CANONICAL ORDER (explicit rule, not an accident): a CDS spec's first
	// required_evidence entry IS the alpha diff. Chosen over order-independent
	// membership deliberately — specs are machine-emitted, canonical form is
	// deterministic, and CUE v0.17 computed validators (list.MatchN/Contains,
	// comprehension guards) do not fire reliably through `vet -d` against a
	// closed base definition. Mechanically enforced by this structural pattern.
	contract: required_evidence: [{id: "diff", kind: "diff", producer: "alpha"}, ...cdd.#RequiredRef]
}
