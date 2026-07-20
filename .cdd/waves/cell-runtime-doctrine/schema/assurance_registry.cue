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
// Invocation:
//   cue vet ./schema/ ./oracle-registry.yaml -d '#AssuranceRegistry'

package crd

#AssuranceRegistry: {
	schema!:   "cn.assurance-registry.v1"
	wave!:     string
	revision!: string & =~"^R[0-9]+$"
	assurance!: [...#AssuranceEntry]
	wave_predicates!: [...#WavePredicate]
}

// The five honest categories.
#Classification: "structural-cue" | "deferred-go" | "mechanically-verifiable" | "evidenced" | "cognitive-review"

// A single in-wave producer WC — the SOLE owner of a deferred-go Go validator. A single string
// (never a list / slash-alternative) enforces "exactly one owner"; #627 is not a member.
#DeferredOwner: "wc-1" | "wc-2" | "wc-3a" | "wc-3b" | "wc-4" | "wc-5"

#AssuranceEntry: {
	owner!:          string
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
	// deferred-go: a named Go validator ships it, owned by EXACTLY ONE in-wave WC.
	if classification == "deferred-go" {
		deferred_to!:    string
		deferred_owner!: #DeferredOwner
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

#WavePredicate: {
	predicate!:      string
	classification!: #Classification
	enforced_by?:    string
	deferred_to?:    string
	deferred_owner?: #DeferredOwner
	// canonical deferred-go validator spec (Required 2 pinning)
	go_artifact_id?:       string
	go_artifact_path?:     string
	inputs?:               string
	result_shape?:         string
	positive_fixture?:     string
	negative_fixture?:     string
	downstream_consumers?: [...string]
	gating_predicate?:     string
	if classification == "structural-cue" {
		enforced_by!: string
	}
	// deferred-go wave validator: single owner + the full pinned spec (id/path/inputs/result/
	// fixtures/consumers/gating). WC-5 consumes each via a sibling_output edge, so it cannot seal
	// until every validator exists and PASSES its gating predicate.
	if classification == "deferred-go" {
		deferred_to!:          string
		deferred_owner!:       #DeferredOwner
		go_artifact_id!:       string
		go_artifact_path!:     string
		inputs!:               string
		result_shape!:         string
		positive_fixture!:     string
		negative_fixture!:     string
		downstream_consumers!: [string, ...string]
		gating_predicate!:     string
	}
}
