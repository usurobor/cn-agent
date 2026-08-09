// Schema ID: cnos.cdd.spec.v1
//
// #CellSpec — the input-side contract: the serialized, CUE-vettable cell that a
// compiled `main.cell` (or a hand-authored contract file) produces and that the
// runner binds to the kernel. It is the mirror of schemas/cdd/receipt.cue
// #Receipt on the *output* side.
//
// γ, V, and δ are kernel-owned and mechanical (Pi β
// msg-cn-pi-cnos-cell-runner-cases-review-31 D2 / cell-prototype-beta-32); they
// deliberately do NOT appear here. A cell spec supplies only the pinned schema
// version, the contract (with producer-attributed required evidence), the
// protocol_id (declared provenance), a builtin seat profile, the typed
// parameter holes, and the α/β skill lines.
//
// See docs/architecture/CDS-CELL-MIGRATION.md (Phase 0).
package cdd

// #Role is which seat is authorized to produce an artifact.
#Role: "alpha" | "beta"

// #RequiredRef names an evidence ref the closed receipt must carry, and the
// producer role authorized to mint it. γ binds it; V checks presence AND
// producer authority (Pi #32 D2).
#RequiredRef: {
	id:       string & !=""
	kind:     string & !=""
	producer: #Role
}

// #Param is a Unix-shaped typed hole. "skill" resolves to a skill and may be
// spliced into a seat via `$name`; "value" is a literal scalar passed to a
// builtin profile.
#Param: {
	kind:     "skill" | "value"
	required: bool | *false
	default?: string
	domain?: [...string]
}

// #Seat is a seat's skill line: literal skill names or `$param` splices.
#Seat: {
	skills: [...string]
}

// #Profile is a builtin v0 seat profile (no cognition yet). "stub" is smoke;
// "bool" is a real independently-checked mechanical episode.
#Profile: "stub" | "bool"

// #CellSpec is the serialized cell.
#CellSpec: {
	// Pinned schema version — a spec must declare it exactly.
	version: "cnos.cellspec.v0"
	contract: {
		id:   string & !=""
		goal: string
		required_evidence?: [...#RequiredRef]
	}
	// Declared protocol (provenance). The v0 runner emits a generic episode
	// receipt and sets protocol_validated=false; it does not validate this.
	protocol_id: string & !=""
	profile:     #Profile // explicit; no default (Pi PR-#718 β D5)
	params?: {[string]: #Param}
	alpha: #Seat
	beta:  #Seat
}
