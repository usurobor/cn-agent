// Schema ID: cnos.cdd.cds.spec.v1
//
// #CDSCellSpec — the CDS input-side overlay on the generic #CellSpec, mirroring
// how schemas/cds/receipt.cue #CDSReceipt overlays #Receipt on the output side.
// It pins the protocol_id V dispatches on and requires the `language` parameter
// hole (which the invoker or a parent cell fills). The generic cell's
// structural invariants are inherited via CUE unification; this overlay only
// constrains the CDS-specific fields.
//
// See docs/architecture/CDS-CELL-MIGRATION.md (Phase 0).
package cds

import "cnos.dev/cnos/schemas/cdd"

#CDSCellSpec: cdd.#CellSpec & {
	// Pinned protocol_id — selects schemas/cds/receipt.cue #CDSReceipt at V.
	protocol_id: "cnos.cdd.cds.receipt.v1"

	// CDS requires a language hole; style is optional with a default. The
	// value domain closes so an out-of-set language fails resolution.
	params: {
		language: {kind: "skill", required: true, domain: ["go", "ocaml", "rust", "python", "typescript"]}
		...
	}

	// A CDS coding cell must at least carry a diff as evidence.
	contract: required_evidence: [...cdd.#RequiredRef]
}
