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
// deliberately does NOT enumerate any fill's fields, providers, or skills;
// the selected fill's overlay (e.g. cds.#CDSPatchAlphaAuthored)
// owns the closed schema for the complete object, exactly as the fill's Go
// decoder owns the strict shape at build time.
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
//
// Division of labour, stated so neither authority is credited with the
// other's work (Pi #57 B2): CUE validates the DECLARATION — that a parameter's
// name, flags and domain are well shaped. Go's `Resolve` validates SUPPLIED
// VALUES — that required parameters were given and that a given value lies in
// its declared domain. A schema cannot check the second; it never sees the
// invocation.
#Param: {
	required: bool | *false
	default?: string
	domain?: [...string]
}

// #Hole is an unresolved `$param` reference. Its identifier grammar is
// #ParamName's — a hole IS a parameter name, so the two cannot diverge.
//
// #Hole and #Concrete live in the GENERIC layer because holes are a generic
// resolution concept and, since the methodology below is a generic field, the
// hole-versus-concrete distinction is now needed here. The CDS overlay aliases
// these rather than restating them: two transcriptions of the same pattern are
// two chances for one authority to admit what the other rejects.
#Hole: =~"^\\$[A-Za-z_][A-Za-z0-9_]*$"

// #Concrete is an already-resolved value: nonempty, and NOT hole-shaped.
// Written as an explicit exclusion because Go treats EVERY `$...` string as a
// hole. A plain `string & !=""` arm silently accepts `$bad-name` — the value
// looks concrete to CUE and is a malformed hole to Go, so the two authorities
// disagree about the same document (Pi #56 C1).
#Concrete: string & !="" & !~"^\\$"

// #Fillable is a position that may carry either, at authoring time.
#Fillable: #Concrete | #Hole

// #Methodology is the cell's ONE ordered methodology bundle: the obligations
// every seat of this cell is held to, declared once. It is a CELL field and
// not a seat field on purpose — while each seat carried its own skill list,
// two lists that were supposed to state the same obligations could drift with
// nothing able to notice, because each was internally consistent.
//
// `!` throughout: without the required-field marker an ABSENT field unifies
// with its declared value and vets clean, so a methodology omitting `skills`
// would pass here while the Go loader rejects it.
//
// Order is meaning — later skills refine earlier ones — so this is a list, not
// a set, and at least one entry: a bundle with no skills states no obligations.
#Methodology: {
	kind!: string & !=""
	skills!: [#Fillable, ...#Fillable]
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
	// OPTIONAL here, and a profile overlay may require it. A cell whose seats
	// are mechanical is held to nothing and declares nothing; making the
	// generic envelope demand a methodology would make one profile's rule
	// everyone's. `cellspec.Parse` mirrors this exactly: it carries the field
	// raw and leaves the shape to the one methodology parser.
	methodology?: #Methodology
	alpha:        #Seat
	beta:         #Seat
}
