// schema/receipt.cue — the runner receipt, shaped after schemas/cdd/receipt.cue.
//
// A receipt is COMPUTED at γ close-out from the cell's transition evidence
// (and finalized by δ) — it is never free-text authored. Every cell writes
// exactly one. The verdict × action → transmissibility table is enforced
// STRUCTURALLY below (mirroring schemas/cdd/receipt.cue): an author-asserted
// transmissibility inconsistent with the (validation.verdict × action) pair
// unifies to _|_ and fails `cue vet`. This is the "fail closed" surface.
//
// This is a minimal spike receipt (see ../DESIGN.md), not the canonical
// cnos.cdd.receipt.v1 — it types only the fields the walking skeleton needs.
package schema

// #ValidationVerdict is V's mechanical output, recorded into the receipt. δ
// never rewrites it (mirrors schemas/cdd/boundary_decision.cue §Q4).
#ValidationVerdict: {
	verdict: "PASS" | "FAIL"
	failed_predicates: [...string]
	warnings: [...string]
}

// #BoundaryDecision is δ's parent-facing action at the cell boundary.
#BoundaryDecision: {
	actor:  string
	action: "accept" | "release" | "reject" | "repair_dispatch" | "override"
}

// #Transmissibility is the derived trust property — what crosses the boundary.
#Transmissibility: "accepted" | "not_transmissible" | "degraded"

// #TransitionRecord is one edge the kernel actually traversed. The receipt's
// transitions[] is the evidence the verdict is computed from.
#TransitionRecord: {
	from:  string
	to:    string
	step:  string
	guard: string
}

// #Receipt is the parent-facing trust surface of a closed (or held) cell.
#Receipt: {
	id:    string
	cell:  "CC" | "PC" | "WC"
	telos: #Telos
	transitions: [...#TransitionRecord]
	evidence_refs: [string]: string
	validation:        #ValidationVerdict
	boundary_decision: #BoundaryDecision
	transmissibility:  #Transmissibility
	// verdict is the one-word terminal outcome of the cell: accept (reached
	// `decided`) or hold (reached `held`, fail-closed).
	verdict:          "accept" | "hold"
	produced_at_ref:  string

	// Structural transmissibility derivation (the fail-closed gate). The four
	// rows the spike can produce are pinned; any other (verdict × action)
	// pair leaves transmissibility unifying to _|_ and fails vet.
	if validation.verdict == "PASS" {
		if boundary_decision.action == "accept" || boundary_decision.action == "release" {
			transmissibility: "accepted"
		}
		if boundary_decision.action == "reject" || boundary_decision.action == "repair_dispatch" {
			transmissibility: "not_transmissible"
		}
	}
	if validation.verdict == "FAIL" {
		if boundary_decision.action == "reject" || boundary_decision.action == "repair_dispatch" {
			transmissibility: "not_transmissible"
		}
		if boundary_decision.action == "override" {
			transmissibility: "degraded"
		}
	}
}
