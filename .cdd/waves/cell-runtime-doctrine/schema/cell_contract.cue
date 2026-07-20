// cnos#671 R9 — FAITHFUL canonical cn.cell.contract.v1 §2 schema (Go + CUE repo; NO Python).
//
// #CellContract is the EXACT §2 `cn.cell.contract.v1` shape as specified by the authoritative
// "Normative constraint model" paragraph of CELL-RUNTIME-CLASSES.md §2 (the constraint model, not
// the illustrative YAML). Every drift the external β found is corrected here:
//   * scope.wave / scope.parent_cell are `string | null` (scalar-or-null), not string;
//   * inputs.required is a NONEMPTY (1+) list of the provenance-tagged reference union
//     ({ref_kind: external, locator} | {ref_kind: sibling_output, producer, output_id});
//   * inputs.optional is a 0+ list of EXTERNAL-locator refs ONLY (never sibling_output);
//   * the external locator union carries all three classes: repo_artifact | control_plane |
//     prior_receipt, each with its class-specific immutable revision component;
//   * requested_output.kind is the full enum artifact | relation_graph | judgment;
//   * acceptance.predicates is 1+; constraints.allowed_paths is 1+; constraints.forbidden_paths
//     and constraints.non_goals are REQUIRED present keys, each 0+;
//   * gates.reason is an ALWAYS-PRESENT `string | null` obeying the §2 truth table — a nonempty
//     string iff either gate boolean is true, null iff both are false (native CUE conditionals);
//   * doctrine_affecting: true ⟹ gates.operator_acceptance_required: true (gate invariant 1);
//   * the struct is CLOSED through stop_conditions: exactly the §2 top-level key set, no more.
//
// #WorkingCellContract is a NAMED REFINEMENT of the canonical shape (NOT a substitute for
// cn.cell.contract.v1): the six children of this wave are Working Cells producing artifacts, so the
// refinement pins cell.class == "working" and requested_output.kind == "artifact". The six real
// child contracts are validated against #WorkingCellContract; #CellContract stays the faithful
// canonical shape and is exercised by the canonical-variant fixtures (nullable scope, an optional
// prior_receipt, the relation_graph and judgment output kinds).
//
// Procedural/semantic checks CUE cannot express declaratively — git ref/content-hash resolution of
// every locator, sibling-output → wave-edge parity, DAG acyclicity, the classification bijection,
// write-surface disjointness, ledger consistency, and completion-evidence derivation — are each
// owned by exactly ONE in-wave Working Cell (see oracle-registry.yaml `deferred_owner` and
// schema/README.md §Validation architecture). They are NOT reimplemented here, and never in Python.
//
// Invocation:
//   cue vet ./schema/ ./contracts/wc-1.cn-cell-contract-v1.yaml -d '#WorkingCellContract'   (…the 6 real WCs)
//   cue vet ./schema/ ./regressions/contract.canonical-relation-graph.yaml -d '#CellContract'  (canonical variants)

package crd

// #CellContract — the exact cn.cell.contract.v1 §2 top-level key set. Closed: any added top-level
// key (e.g. the R4-regression `oracles`, or `consumers`/`completion_signal`) fails `cue vet`.
#CellContract: {
	schema!:             "cn.cell.contract.v1"
	cell!:               #Cell
	scope!:              #Scope
	intent_ref!:         #IntentRef
	inputs!:             #Inputs
	requested_output!:   #RequestedOutput
	acceptance!:         #Acceptance
	constraints!:        #Constraints
	gates!:              #Gates
	doctrine_affecting!: bool
	stop_conditions!: [...string]

	// Gate invariant 1 (§2, decidable): a doctrine-affecting matter ALWAYS requires operator
	// acceptance. doctrine_affecting is authoritative; operator_acceptance_required is derived.
	if doctrine_affecting {
		gates: operator_acceptance_required: true
	}

	// Gate invariant 2 (§2 truth table over gates.reason, always a present key of type
	// `string | null`): a NONEMPTY string iff at least one gate boolean is true; NULL iff both are
	// false. Expressed as two mutually-exclusive CUE conditionals over the concrete gate booleans.
	if gates.operator_authorization_required || gates.operator_acceptance_required {
		gates: reason!: string & !="" // some gate true ⇒ reason is a nonempty scalar
	}
	if !gates.operator_authorization_required && !gates.operator_acceptance_required {
		gates: reason!: null // both gates false ⇒ reason is null
	}
}

#Cell: {
	id!:            string
	class!:         "working" | "planning" | "cohering" // cell.class enum
	mode!:          string
	protocol!:      string
	matter_domain!: string
}

// scope.wave / scope.parent_cell are scalar-or-null (§2): a root cell carries null, not a string.
#Scope: {
	repo!:        string
	wave!:        string | null
	parent_cell!: string | null
}

#IntentRef: {
	schema!: "cn.intent.v1"
	id!:     string
	carrier!: {
		kind!: string
		ref!:  string
	}
}

// inputs.required — NONEMPTY (1+) provenance-tagged reference union.
// inputs.optional — 0+ EXTERNAL-locator refs ONLY (never sibling_output).
#Inputs: {
	required!: [#InputRef, ...#InputRef] // 1+ (a legitimate external root is a valid sole member)
	optional!: [...#ExternalRef]         // 0+, external-locator refs only; a sibling_output here fails
}

// #InputRef — the provenance-tagged union: an external ref (classed locator) OR a sibling_output
// ref (producer + output_id). ref_kind is the discriminant enum.
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

// #Locator — the locator-classed external union: all THREE §2 classes, each carrying its
// class-specific immutable revision component (a bare mutable path/identity is insufficient).
#Locator: #RepoArtifactLocator | #ControlPlaneLocator | #PriorReceiptLocator

#RepoArtifactLocator: {
	kind!:   "repo_artifact"
	repo!:   string
	commit!: string // REQUIRED immutable commit/tree/blob revision
	path!:   string // resolved AT that revision
}

#ControlPlaneLocator: {
	kind!:     "control_plane"
	ref!:      string // stable object identity (carrier)
	revision!: string // REQUIRED: binds immutable CONTENT (content hash + snapshot); identity-only rejected
}

#PriorReceiptLocator: {
	kind!:    "prior_receipt"
	receipt!: string // receipt identity
	binding!: string // REQUIRED content-addressed / revision-bound immutable binding
}

// requested_output.kind — the full §2 enum. WC/PC/CC output telos: artifact | relation_graph | judgment.
#RequestedOutput: {
	id!:   string
	kind!: "artifact" | "relation_graph" | "judgment"
	path!: string
}

#Acceptance: {
	predicates!: [string, ...string] // 1+ (§5)
}

// constraints — allowed_paths is 1+; forbidden_paths and non_goals are REQUIRED present keys, each 0+.
#Constraints: {
	allowed_paths!: [string, ...string] // 1+
	forbidden_paths!: [...string]       // required key, 0+
	non_goals!: [...string]             // required key, 0+
}

// gates — reason is an ALWAYS-PRESENT scalar-or-null; the truth table is enforced in #CellContract.
#Gates: {
	operator_authorization_required!: bool
	operator_acceptance_required!:    bool
	reason!:                          string | null
}

// #WorkingCellContract — a NAMED refinement of the canonical shape for this wave's six Working
// Cells (WC output telos = artifact). It does NOT replace cn.cell.contract.v1; it further constrains
// the two class/telos fields and inherits every canonical §2 constraint by unification.
#WorkingCellContract: #CellContract & {
	cell: class:            "working"
	requested_output: kind: "artifact"
}
