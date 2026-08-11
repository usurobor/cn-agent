// schema/contract.cue — the WC/PC work contract handed to a kernel run.
//
// A contract is the α-input of a cell: it pins the telos (output kind), the
// target the cell acts on, the instruction, and the mechanical acceptance
// predicates V asserts at validation. This is an exploratory-spike schema
// (see ../DESIGN.md), not canonical R12 doctrine.
package schema

// #Telos classifies a cell purely by its output kind — the R12 rule that
// distinguishes CC / PC / WC without any other structural difference.
#Telos: "artifact" | "relation_graph" | "judgment"

// #Contract is the pinned unit of work for one kernel run.
#Contract: {
	id:          string
	telos:       #Telos
	target:      string
	instruction: string
	// acceptance predicates V asserts mechanically at validation (V pass iff
	// every predicate holds). Free-text ids here; the actor maps each to a
	// mechanical check.
	acceptance: [...string]
}
