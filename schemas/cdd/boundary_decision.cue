// see schemas/cdd/README.md §Scope-Lift Invariant.
//
// Schema ID: cnos.cdd.boundary_decision.v1
//
// Types the validator output (#ValidationVerdict), the boundary actor's
// decision (#BoundaryDecision), the override sub-record (#Override required
// iff action == "override"), and the derived transmissibility enum
// (#Transmissibility). The verdict × action → transmissibility table is
// enforced as a structural constraint in receipt.cue, where both surfaces
// co-exist; the table itself is documented in schemas/cdd/README.md.
//
// Semantic owner: src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md
// (§Validation Interface, §Q4). v1 is the *schema-artifact* version, not a
// CDD protocol version.

package cdd

// #ValidationVerdict is V's output. Emitted by the validator; recorded into
// receipt.validation. δ never rewrites this field; an override of a non-PASS
// verdict lives in #BoundaryDecision.override, not here (RECEIPT-VALIDATION.md
// §Q4: "Override never rewrites the ValidationVerdict").
#ValidationVerdict: {
	verdict:           "PASS" | "FAIL"
	failed_predicates: [...string]
	warnings:          [...string]
	provenance: {
		validator_identity: string
		validator_version:  string
		checked_at:         string
		input_refs: {
			contract_ref:      string
			receipt_ref:       string
			evidence_root_ref: string
		}
	}
}

// #Override carries the structured degraded-boundary action δ takes when it
// chooses to accept a non-PASS receipt anyway. degraded_state is pinned to
// true: when an override is recorded, the receipt closed in degraded state
// by definition (see §Scope-Lift Invariant in the README). Required fields
// per RECEIPT-VALIDATION.md §Q4 "Required fields".
#Override: {
	actor:                        string
	justification:                string
	original_validation_verdict:  #ValidationVerdict
	failed_predicates_overridden: [...string]
	degraded_state:               true

	followup_issue?:   string
	risk?:             string
	expires_at?:       string
	review_condition?: string
}

// #BoundaryDecision is δ's parent-facing record of what δ did with the receipt
// at the boundary. action enum membership is pinned to five values; any other
// value fails cue vet. #Override is required iff action == "override" — both
// directions of the iff are enforced.
//
// The required-iff is structurally encoded so authors cannot:
//   (a) declare action: override and omit the override block (incomplete), or
//   (b) attach an override block to action: accept/release/reject/repair_dispatch
//       (incoherent — override is a boundary action, not metadata).
#BoundaryDecision: {
	actor:       string
	action:      "accept" | "release" | "reject" | "repair_dispatch" | "override"
	decided_at?: string
	override?:   #Override

	if action == "override" {
		// override required when action is override.
		override: #Override
	}
	if action != "override" {
		// override forbidden when action is anything else.
		override?: _|_
	}
}

// #Transmissibility is the derived parent-facing property — what crosses the
// trust boundary. Computed structurally in receipt.cue from
// (#ValidationVerdict.verdict × #BoundaryDecision.action). Authors who set
// a transmissibility value inconsistent with that mapping (e.g. verdict: FAIL,
// action: override, transmissibility: accepted) fail cue vet.
#Transmissibility: "accepted" | "not_transmissible" | "degraded"
