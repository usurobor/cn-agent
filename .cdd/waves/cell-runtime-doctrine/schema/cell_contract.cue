// cnos#671 R8 — PLAN-LOCAL, TRANSITIONAL structural schema for this wave's child contracts.
//
// This is NOT the canonical deliverable. It is a design input: the exact cn.cell.contract.v1 §2
// envelope that WC-1/WC-2 formalize into the canonical CUE under schemas/cdd/. It exists so that
// `cue vet` structurally validates the six plan-local contracts/*.yaml BEFORE authorization —
// natively, with no procedural validator:
//   * this is a CLOSED definition (#CellContract), so an unknown top-level key is rejected;
//   * every field is typed, so a bool given as the string "false" is rejected;
//   * enums (cell.class, ref_kind, locator.kind, requested_output.kind) reject any other member;
//   * required (`!`) fields reject an omitted section;
//   * the gate invariant (doctrine_affecting ⇒ operator_acceptance_required, and reason required
//     when acceptance is required) is a native CUE conditional constraint.
//
// Procedural/semantic checks CUE cannot express declaratively — git ref / content-hash resolution
// of every locator (intent_ref.id, repo_artifact commit:path, control_plane revision hashes),
// sibling-output → wave-edge parity, and the classification bijection — are NAMED-AND-DEFERRED to
// the Go validators WC-1/WC-2 and #627 S2–S3 ship (see README §Validation architecture). They are
// NOT reimplemented here, and never in Python.
//
// Invocation (one per contract):
//   cue vet ./schema/ ./contracts/wc-1.cn-cell-contract-v1.yaml -d '#CellContract'
//   … wc-2 … wc-3a … wc-3b … wc-4 … wc-5 …

package crd

// #CellContract — the exact cn.cell.contract.v1 §2 top-level key set. Closed: an added top-level
// key (e.g. the R4-regression `oracles`, or `consumers`/`completion_signal`) fails `cue vet`.
#CellContract: {
	schema!:            "cn.cell.contract.v1"
	cell!:              #Cell
	scope!:             #Scope
	intent_ref!:        #IntentRef
	inputs!:            #Inputs
	requested_output!:  #RequestedOutput
	acceptance!:        #Acceptance
	constraints!:       #Constraints
	gates!:             #Gates
	doctrine_affecting!: bool
	stop_conditions!: [...string]

	// Gate invariant 1 (native CUE, replaces the removed validate.py check (f)):
	// a doctrine-affecting cell MUST require operator acceptance.
	if doctrine_affecting {
		gates: operator_acceptance_required: true
	}
	// reason is mandatory whenever a gate bool is asserted true.
	if gates.operator_acceptance_required {
		gates: reason!: string
	}
	if gates.operator_authorization_required {
		gates: reason!: string
	}
}

#Cell: {
	id!:            string
	class!:         "working" | "planning" | "cohering" // cell.class enum
	mode!:          string
	protocol!:      string
	matter_domain!: string
}

#Scope: {
	repo!:        string
	wave!:        string
	parent_cell!: string
}

#IntentRef: {
	schema!: "cn.intent.v1"
	id!:     string
	carrier!: {
		kind!: string
		ref!:  string
	}
}

#Inputs: {
	required!: [...#InputRef]
	optional!: [...#InputRef]
}

// #InputRef — the provenance-tagged union: an external ref (with a classed locator) OR a
// sibling_output ref (producer + output_id). ref_kind is the discriminant enum.
#InputRef: #ExternalRef | #SiblingRef

#ExternalRef: {
	ref_kind!: "external"
	locator!:  #Locator
}

#SiblingRef: {
	ref_kind!:  "sibling_output"
	producer!:  string
	output_id!: string
}

#Locator: #RepoArtifactLocator | #ControlPlaneLocator

#RepoArtifactLocator: {
	kind!:   "repo_artifact"
	repo!:   string
	commit!: string
	path!:   string
}

#ControlPlaneLocator: {
	kind!:     "control_plane"
	ref!:      string
	revision!: string
}

#RequestedOutput: {
	id!:   string
	kind!: "artifact" // requested_output.kind enum
	path!: string
}

#Acceptance: {
	predicates!: [...string]
}

#Constraints: {
	allowed_paths!: [...string]
	forbidden_paths?: [...string]
	non_goals?: [...string]
}

#Gates: {
	operator_authorization_required!: bool
	operator_acceptance_required!:    bool
	reason?:                          string
}
