// Schema ID: cnos.cdd.spec.v1
//
// #CellSpec — the input-side contract: the serialized, CUE-vettable cell that a
// compiled `main.cell` (or a hand-authored contract file) produces and that the
// runner binds to the kernel. It is the mirror of schemas/cdd/receipt.cue
// #Receipt on the *output* side.
//
// Seats are FILL-OWNED (msg-cn-pi-cnos-cds-fill-construction-51,
// operator-ratified): each seat is one tagged object whose `fill` selects a
// constructor and whose remaining fields are that constructor's arguments.
// The generic schema therefore owns only the minimum tagged envelope — it
// deliberately does NOT enumerate any fill's fields, providers, skills, or
// workspaces; the selected fill's overlay (e.g. cds.#CDSPatchAlpha) owns the
// closed schema for the complete object, exactly as the fill's Go decoder
// owns the strict shape at build time.
//
// γ, V, and δ are kernel-owned and mechanical; they deliberately do NOT
// appear here. See docs/architecture/CDS-CELL-MIGRATION.md.
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

// #ParamName is the ONE identifier grammar for a parameter, and therefore for
// the `$name` hole that references it. It lives in the generic layer because
// holes are a generic resolution concept: a name legal here must be legal
// everywhere, or a spec resolves in Go and is rejected by CUE (Pi #55 C1).
// `cellspec.Parse` enforces the identical pattern.
#ParamName: =~"^[A-Za-z_][A-Za-z0-9_]*$"

// #Param is a Unix-shaped typed hole. Holes appear as `$name` string values
// inside seat declarations and are replaced in place at resolution. There is
// deliberately no "kind": what a filled value MEANS belongs to the fill that
// consumes it, not to the generic envelope.
#Param: {
	required: bool | *false
	default?: string
	domain?: [...string]
}

// #Seat is the minimum tagged envelope: a fill id, plus whatever constructor
// arguments that fill owns. `fill!` — required-field marker, so an absent
// tag is rejected rather than defaulting (Pi round-5 D2 parity).
#Seat: {
	fill!: string & !=""
	...
}

// #CellSpec is the serialized cell.
#CellSpec: {
	// Pinned schema version — a spec must declare it exactly.
	version: "cnos.cellspec.v0"
	contract: {
		id:   string & !=""
		goal: string & !=""
		required_evidence?: [...#RequiredRef]
	}
	// Declared protocol (provenance). The v0 runner emits a generic episode
	// receipt and sets protocol_validated=false; it does not validate this.
	protocol_id: string & !=""
	params?: {[#ParamName]: #Param}
	alpha: #Seat
	beta:  #Seat
}
