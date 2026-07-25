// schema/wave.cue — the PC's relation-graph output: a wave of WC contracts.
//
// v0's wave carries exactly one WC contract (multi-WC waves / parallel
// branches are explicitly deferred — see ../DESIGN.md §Explicitly deferred).
// #Contract resolves from contract.cue in this same package (no import).
package schema

// #Wave is the α-output of a PC run: a set of downstream WC contracts derived
// from a measurement, plus the ref of the #CM it was planned from.
#Wave: {
	id:          string
	from_cm_ref: string
	contracts: [...#Contract]
}
