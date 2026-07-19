// cnos#671 R8 — PLAN-LOCAL, TRANSITIONAL structural schema for oracle-registry.yaml.
//
// Design input; NOT the canonical deliverable. Closed #AssuranceRegistry so `cue vet` validates the
// registry's shape natively: the top-level key set, a per-entry `classification` ENUM over the five
// honest categories, and the required fields per category (a structural-cue entry MUST name its
// enforcing `cue vet`; a deferred-go entry MUST name the Go validator that ships it; a
// mechanically-verifiable entry MUST carry ownership + fixtures + command + receipt binding; an
// evidenced entry MUST bind receipt evidence; a cognitive-review entry MUST name the independent
// review). The bijection registry ⇄ union(child acceptance.predicates) is a DEFERRED-GO check.
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

#AssuranceEntry: {
	owner!:          string
	predicate!:      string
	classification!: #Classification

	// superset of category fields (optional at the top; required per-category below)
	enforced_by?:                 string
	deferred_to?:                 string
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
	// deferred-go: a named Go validator WC-1/WC-2/#627 S2–S3 ship; not run at this tree.
	if classification == "deferred-go" {
		deferred_to!: string
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
	if classification == "structural-cue" {
		enforced_by!: string
	}
	if classification == "deferred-go" {
		deferred_to!: string
	}
}
