// Schema ID: cnos.cdd.spec.v1
//
// #CellSpec — the input-side contract: the serialized, CUE-vettable cell that a
// compiled `main.cell` (or a hand-authored contract file) produces and that the
// runner binds to the kernel. It is the mirror of schemas/cdd/receipt.cue
// #Receipt on the *output* side: γ emits a receipt vetted against #Receipt; the
// invoker compiles a spec vetted against #CellSpec.
//
// γ, V, and δ are kernel-owned and mechanical (Pi β
// msg-cn-pi-cnos-cell-runner-cases-review-31, D2); they deliberately do NOT
// appear here. A cell spec supplies only the contract, the protocol_id (which
// selects the receipt schema V dispatches on), the typed parameter holes, and
// the α/β skill lines.
//
// See docs/architecture/CDS-CELL-MIGRATION.md (Phase 0).
package cdd

// #RequiredRef names an evidence ref the closed receipt must carry; γ binds it
// and V checks its presence.
#RequiredRef: {
	id:   string & !=""
	kind: string & !=""
}

// #Param is a Unix-shaped typed hole. `kind` is the resolution family (e.g.
// "skill"); a `required` hole with no `default` must be filled by the invoker;
// an optional `domain` closes the value set so a typo fails resolution.
#Param: {
	kind:     string & !=""
	required: bool | *false
	default?: string
	domain?: [...string]
}

// #Seat is a seat's skill line. Entries are literal skill names or `$param`
// references spliced from resolved parameters.
#Seat: {
	skills: [...string]
}

// #CellSpec is the serialized cell.
#CellSpec: {
	contract: {
		id:   string & !=""
		goal: string
		required_evidence?: [...#RequiredRef]
	}
	// Selects the receipt schema V dispatches on (e.g. cnos.cdd.cds.receipt.v1).
	protocol_id: string & !=""
	params?: {[string]: #Param}
	alpha: #Seat
	beta:  #Seat
	budget?: {
		tokens?: int & >0
	}
}
