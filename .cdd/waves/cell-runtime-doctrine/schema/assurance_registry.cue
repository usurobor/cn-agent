// cnos#671 R9 — PLAN-LOCAL, TRANSITIONAL structural schema for oracle-registry.yaml.
//
// Design input; NOT the canonical deliverable. Closed #AssuranceRegistry so `cue vet` validates the
// registry's shape natively: the top-level key set, a per-entry `classification` ENUM over the five
// honest categories, and the required fields per category (a structural-cue entry MUST name its
// enforcing `cue vet`; a deferred-go entry MUST name the Go validator that ships it AND its SINGLE
// in-wave owner `deferred_owner`; a mechanically-verifiable entry MUST carry ownership + fixtures +
// command + receipt binding; an evidenced entry MUST bind receipt evidence; a cognitive-review entry
// MUST name the independent review).
//
// R9 (Required 2): every deferred-go check is owned by EXACTLY ONE in-wave producer WC, pinned as a
// single-valued `deferred_owner` enum (not a slash-alternative, and never #627 — #627 S2–S3 stay
// DOWNSTREAM consumers/canonicalizers, not owners of this wave's closure evidence). The eight
// canonical wave-level validators additionally pin, in `wave_predicates`, a stable Go artifact
// id+path, typed inputs + result/evidence shape, positive + named negative fixtures, downstream
// consumers, and a gating predicate whose PASS gates the owning child's completion — so WC-5 cannot
// seal until every deferred validator exists and passes.
//
// R10 (this repair) — FORWARD-ONLY ACYCLIC ASSURANCE GRAPH. Each child deferred acceptance entry
// creates an assurance edge deferred_owner -> owner. To keep the combined (construction + assurance)
// graph acyclic, a child acceptance predicate may be verified ONLY by the child ITSELF or a
// construction-PREDECESSOR of that child — never a successor. #AllowedVerifier pins the transitive
// predecessor closure over the ACCEPTED six-node graph; every deferred-go #AssuranceEntry is
// constrained so deferred_owner ∈ {owner} ∪ predecessors(owner). A WC-2 acceptance predicate with
// deferred_owner wc-1, or a WC-3b one with deferred_owner wc-5, is a BACKWARD edge -> REJECTED
// (see schema/regressions/registry.bad-backward-edge-*.yaml). Whole-wave-property validators are
// relocated off child edges: the cross-contract oracle-ownership/classification bijection is a
// WAVE-BOUNDARY pre-authorization predicate (deferred_owner: "wave"); combined-graph acyclicity is a
// forward WC-3b self-owned check; ledger + classification-totality remain forward-owned by terminal WC-5.
//
// Invocation:
//   cue vet ./schema/ ./oracle-registry.yaml -d '#AssuranceRegistry'

package crd

import "list"

// The six construction nodes (the ACCEPTED graph). owner + deferred_owner range over these.
#Node: "wc-1" | "wc-2" | "wc-3a" | "wc-3b" | "wc-4" | "wc-5"

// {self} ∪ transitive construction-PREDECESSORS, per node, over WC-2 -> WC-1 -> {WC-3a,WC-3b,WC-4}
// -> WC-5. A child deferred acceptance predicate may be verified only by a member of this set for its
// owner — guaranteeing every assurance edge (deferred_owner -> owner) is FORWARD, so the union with
// the construction edges stays acyclic (sole root wc-2, terminal wc-5, topological sort succeeds).
#AllowedVerifier: {
	"wc-2":  ["wc-2"]
	"wc-1":  ["wc-1", "wc-2"]
	"wc-3a": ["wc-3a", "wc-1", "wc-2"]
	"wc-3b": ["wc-3b", "wc-1", "wc-2"]
	"wc-4":  ["wc-4", "wc-1", "wc-2"]
	"wc-5":  ["wc-5", "wc-1", "wc-2", "wc-3a", "wc-3b", "wc-4"]
}

#AssuranceRegistry: {
	schema!:   "cn.assurance-registry.v1"
	wave!:     string
	revision!: string & =~"^R[0-9]+$"
	assurance!: [...#AssuranceEntry]
	wave_predicates!: [...#WavePredicate]
}

// The five honest categories.
#Classification: "structural-cue" | "deferred-go" | "mechanically-verifiable" | "evidenced" | "cognitive-review"

// A single WC — the SOLE owner of a deferred-go Go validator. A single string (never a list /
// slash-alternative) enforces "exactly one owner"; #627 is not a member. Same value space as #Node.
#DeferredOwner: #Node

#AssuranceEntry: {
	owner!:          #Node
	predicate!:      string
	classification!: #Classification

	// superset of category fields (optional at the top; required per-category below)
	enforced_by?:                 string
	deferred_to?:                 string
	deferred_owner?:              #DeferredOwner
	ownership?:                   "emitted" | "existing"
	checker?:                     string
	schema?:                      string
	positive_fixture?:            string
	negative_fixture?:            string
	command?:                     string
	receipt_evidence_predicate?:  string
	receipt_evidence?:            string
	independent_review?:          string

	// structural-cue: enforced NOW by cue vet against a plan-local #Def.
	if classification == "structural-cue" {
		enforced_by!: string
	}
	// deferred-go: a named Go validator ships it, owned by EXACTLY ONE WC.
	if classification == "deferred-go" {
		deferred_to!:    string
		deferred_owner!: #DeferredOwner
		// FORWARD-ONLY (R10): the verifier of a child's acceptance predicate must be the child
		// ITSELF or a construction-PREDECESSOR — never a successor. deferred_owner ∈
		// {owner} ∪ predecessors(owner). A backward edge (e.g. owner wc-2 / deferred_owner wc-1,
		// or owner wc-3b / deferred_owner wc-5) makes _forward_only false -> conflict -> REJECTED.
		_forward_only: list.Contains(#AllowedVerifier[owner], deferred_owner)
		_forward_only: true
	}
	// mechanically-verifiable: a child-emitted Go checker or CUE schema + fixtures (no Python).
	if classification == "mechanically-verifiable" {
		ownership!:                  "emitted" | "existing"
		positive_fixture!:           string
		negative_fixture!:           string
		command!:                    string
		receipt_evidence_predicate!: string
	}
	if classification == "evidenced" {
		receipt_evidence!: string
	}
	if classification == "cognitive-review" {
		independent_review!: string
	}
}

// A wave-level validator is owned EITHER by a single construction node (forward: it ships the Go
// artifact and its PASS gates that node's OWN completion) OR by the WAVE BOUNDARY itself ("wave":
// a whole-wave / cross-contract property run at pre-authorization, owned by no child).
#WaveOwner: #DeferredOwner | "wave"

#WavePredicate: {
	predicate!:      string
	classification!: #Classification
	enforced_by?:    string
	deferred_to?:    string
	deferred_owner?: #WaveOwner
	// canonical deferred-go validator spec (Required 2 pinning)
	go_artifact_id?:       string
	go_artifact_path?:     string
	inputs?:               string
	result_shape?:         string
	positive_fixture?:     string
	negative_fixture?:     string
	downstream_consumers?: [...string]
	gating_predicate?:         string
	wave_authorization_gated?: true
	if classification == "structural-cue" {
		enforced_by!: string
	}
	// deferred-go wave validator: single owner + the full pinned spec (id/path/inputs/result/
	// fixtures/consumers). A CHILD-owned validator additionally pins a `gating_predicate` (an
	// acceptance predicate on the owner whose PASS gates that owner's forward completion). A
	// WAVE-owned validator (deferred_owner: "wave") is gated at the wave AUTHORIZATION boundary
	// (wave_authorization_gated: true) — a pre-authorization whole-wave property, not a child edge.
	if classification == "deferred-go" {
		deferred_to!:          string
		deferred_owner!:       #WaveOwner
		go_artifact_id!:       string
		go_artifact_path!:     string
		inputs!:               string
		result_shape!:         string
		positive_fixture!:     string
		negative_fixture!:     string
		downstream_consumers!: [string, ...string]
		if deferred_owner == "wave" {
			wave_authorization_gated!: true
		}
		if deferred_owner != "wave" {
			gating_predicate!: string
		}
	}
}
