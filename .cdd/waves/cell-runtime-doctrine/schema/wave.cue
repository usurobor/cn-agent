// cnos#671 R9 — PLAN-LOCAL, TRANSITIONAL structural schema for wave.cn-wave-v1.yaml.
//
// Design input (WC-1/WC-3b formalize the durable cn.wave.v1 CUE); NOT the canonical deliverable.
// Closes the R7 wave-envelope + completion blockers NATIVELY via `cue vet`:
//   * closed top-level struct + `schema` const     -> `schema: nonsense` / an added key is rejected;
//   * nonempty typed external_roots ([_, ...string]) -> `external_roots: []` is rejected;
//   * gates.wave_authorization with required revision-bound fields -> a missing field is rejected;
//   * contract_ref_resolution with required fields   -> `contract_ref_resolution: {}` is rejected;
//   * edge_derivation constants pinned true          -> deleting/flipping parity_required/acyclic
//     (or external_refs_create_no_edge) is rejected;
//   * a FIXED completion governing formula const     -> a deleted / literal-mutated `wave_complete`
//     formula is rejected (required + pinned value);
//   * per-fixture child records typed as CLOSED 5-bool structs -> a missing `evidence_bound`, or a
//     bool given as the string "false" (e.g. receipt_bound: "false"), is rejected;
//   * uniqueness folds over nodes/edges              -> a duplicate node id / duplicate edge conflicts;
//   * a non-empty typed STOP set with an effect enum -> `stop_conditions: []` / a bad effect rejected.
//
// Procedural/semantic checks NAMED-AND-DEFERRED to Go (WC-3b/WC-1, #627 S2–S3): graph acyclicity
// (DAG), sibling-output→edge PARITY, git ref/content-hash resolution, and the completion-evidence
// DERIVATION over resolved child receipts (that a fixture's `expected` equals the value computed from
// the records, and that child_complete ignores no constituent). CUE fixes the SHAPE; Go computes.
//
// Invocation:
//   cue vet ./schema/ ./wave.cn-wave-v1.yaml -d '#Wave'

package crd

#Wave: {
	schema!:                 "cn.wave.v1"
	wave!:                   string
	revision!:               string & =~"^R[0-9]+$"
	goal!:                   string
	intent!:                 #WaveIntent
	grounding!:              #Grounding
	external_roots!: [string, ...string] // NONEMPTY + fully typed: `external_roots: []` fails
	contract_ref_resolution!: #ContractRefResolution
	nodes!: [...#WaveNode]
	edges!: [...#Edge]
	edge_derivation!:        #EdgeDerivation
	gates!:                  #WaveGates
	stop_conditions!: [#StopCondition, ...#StopCondition] // NONEMPTY + fully typed STOP set
	completion!:             #Completion

	// Uniqueness folds: a repeated node id / edge (from,to) yields two distinct indices at the
	// same key -> conflicting values -> `cue vet` error (closes duplicate-node/edge invisibility).
	_unique_node_ids: {for i, n in nodes {"\(n.id)": i}}
	_unique_edges: {for i, e in edges {"\(e.from)=>\(e.to)": i}}
}

#WaveIntent: {
	schema!:  "cn.intent.v1"
	id!:      string
	locator!: string
}

#Grounding: {
	source_snapshot!: {
		ref_kind!: "external"
		locator!: {
			kind!:     "control_plane"
			ref!:      string
			revision!: string
			note?:     string
		}
	}
	alpha_derivative!: {
		ref_kind!: "external"
		locator!: {
			kind!:              "repo_artifact"
			repo!:              string
			resolved_at!:       string
			path!:              string
			content_hash!:      string
			references_source!: string
		}
	}
	reconcile_627_ref!: {
		ref_kind!: "external"
		locator!: {
			kind!:         "repo_artifact"
			repo!:         string
			resolved_at!:  string
			path!:         string
			content_hash!: string
		}
	}
}

#ContractRefResolution: {
	rule!:                   string
	authorized_wave_commit!: string
	per_child_content_binding!: string
	stale_on!:               string
}

#WaveNode: {
	id!:             string
	role?:           string
	contract_ref!:   string
	contract_sha256!: string & =~"^[0-9a-f]{64}$"
	output_id!:      string
	output_path!:    string
}

#Edge: {
	from!: string
	to!:   string
	kind!: "depends_on" // edge kind enum
}

#EdgeDerivation: {
	rule!:                    string
	external_refs_create_no_edge!: true // pinned const
	parity_required!:         true       // pinned const — deleting or flipping fails
	acyclic!:                 true       // pinned const
	roots!: [...string]
	critical_path!: [...string]
}

#WaveGates: {
	wave_authorization!: {
		authority!:     string
		boundary!:      string
		revision_bound!: string // revision-bound authorization field required
	}
	per_child!: {
		operator_authorization_required!: bool
		operator_acceptance_required!:    bool
		reason!:                          string
	}
}

#StopCondition: {
	trigger!:    string
	detect!:     string
	effect!:     "hold" | "replan" | "escalate" // effect enum
	transition!: string
	receipt!:    string
}

// The FIVE completion constituents — the exact field set of #ChildRecord. required_constituents and
// the child_complete.definition are constrained to this set so the three cannot drift (Required 3).
#Constituent: "requested_output_produced" | "acceptance_all_pass" | "class_v_pass" | "receipt_bound" | "evidence_bound"

// #Completion — the whole-wave completion block. The governing formula is a FIXED required const,
// not a mutable or absent field. Each fixture's per-node record is a CLOSED 5-bool struct, so a
// missing `evidence_bound` or a `receipt_bound: "false"` (string) is rejected natively. R9 adds a
// typed completion-evidence RESOLVER INPUT contract (evidence_resolver) owned by WC-5: the 5 booleans
// are DERIVED by resolving typed bindings, not author-supplied.
#Completion: {
	construction_set_N!: [string, ...string] // nonempty + typed
	terminal_node!:      string
	// PINNED governing formula const (BLOCKER close): wave completion is all(N) AND terminal.
	wave_complete!: "all(N child_complete) AND child_complete(terminal)"
	note!:          string
	child_complete!:    #ChildCompleteDef
	evidence_resolver!: #EvidenceResolver // R9: typed resolver-input contract (owned by WC-5)
	predicates!: [...#Predicate]
	fixtures!: [...#Fixture]
}

// child_complete DERIVATION — R9 consistency: required_constituents lists ALL FIVE #ChildRecord
// constituents (incl. evidence_bound), and the definition reads each of the five EXACTLY (no
// constituent ignored). Both are pinned to the canonical constituent set, so a required_constituents
// that drops evidence_bound, or a definition that ignores a constituent, is rejected by cue vet NOW.
#ChildCompleteDef: {
	kind!: string
	required_constituents!: [
		"requested_output_produced",
		"acceptance_all_pass",
		"class_v_pass",
		"receipt_bound",
		"evidence_bound",
	]
	definition!: {
		all: [
			{record: "requested_output_produced"},
			{record: "acceptance_all_pass"},
			{record: "class_v_pass"},
			{record: "receipt_bound"},
			{record: "evidence_bound"},
		]
	}
}

// #EvidenceResolver — the typed completion-evidence RESOLVER INPUT contract (Required 3). The Go
// completion-validator (owned SOLELY by WC-5, Required 2) resolves these typed per-node bindings into
// the 5 derived constituent booleans; the booleans are NOT author-supplied. Its PASS binds into WC-5
// closure (it is a constituent of wc-5 child_complete via wc-5's gating acceptance predicate).
#EvidenceResolver: {
	owner!:            "wc-5" // single owner (Required 2); never #627
	go_artifact_id!:   string
	go_artifact_path!: string
	derivation_rule!:  string                 // how each of the 5 booleans is DERIVED from the bindings
	input_contract!:   #CompletionEvidenceInput // the typed per-node resolver input template
	binds_into!:       "wc-5 closure"         // resolver PASS is a constituent of wc-5 child_complete
}

// #CompletionEvidenceInput — the typed per-node evidence the resolver consumes. Each field is an
// identity/content binding; the 5 constituent booleans are DERIVED from them (never hand-authored).
#CompletionEvidenceInput: {
	node!: string
	// node/output identity + content binding -> requested_output_produced
	output_binding!: {
		output_id!:     string
		output_digest!: string // digest of the produced bytes (content identity, not the mutable path)
		output_path!:   string
	}
	// acceptance result set -> acceptance_all_pass
	acceptance_result_set!: {
		rule!: string
	}
	// V verdict/receipt binding -> class_v_pass
	v_binding!: {
		verdict_locator!: string
		receipt_ref!:     string
	}
	// β/γ evidence locator (where applicable) -> evidence_bound
	evidence_locator!: {
		rule!: string
	}
	// receipt identity/hash -> receipt_bound
	receipt_identity!: {
		receipt_id!:   string
		receipt_hash!: string
	}
	// the 5 constituent booleans are DERIVED by resolving the bindings above, not author-supplied.
	derives!: [...#Constituent]
}

#Predicate: {
	name!: string
	expr!: _
}

// Closed 5-bool record — every constituent required and typed bool. The booleans a fixture carries
// are TRUTH-TABLE FIXTURES over the values the resolver would derive; the resolver-input contract
// above (#CompletionEvidenceInput) is what a live wave resolves, per Required 3.
#ChildRecord: {
	requested_output_produced!: bool
	acceptance_all_pass!:       bool
	class_v_pass!:              bool
	receipt_bound!:             bool
	evidence_bound!:            bool
}

#Fixture: {
	name!:  string
	note?:  string
	records!: [string]: #ChildRecord
	claimed_child_complete!: [string]: bool
	expected!: {
		wc5_ready!:     bool
		wc5_complete!:  bool
		wave_complete!: bool
	}
}
