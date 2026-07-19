// schemas/coherence_methodology.cue — cnos projection of the generic TSC
// coherence-methodology essence contract.
//
// Compatibility source:
// usurobor/tsc@26aab5023f03dc7d0abf82e5fdba20134fc6adad
// schemas/skill.cue:#CoherenceMethodology.
//
// TSC owns the generic contract and executor. cnos carries this pinned,
// deliberately essence-only projection so a cnos-owned domain methodology can
// be validated without pretending TSC is a cnos package or requiring network
// access in CI. Deployment bindings remain TSC/runtime concerns.

package skill

#CoherenceMethodology: {
	registry:     !=""
	targets:      [...string]
	cross_target: bool
	instruction:  !=""
	// cnos currently validates and executes methodology authority paths from a
	// repository checkout. Installed-package activation is refused until TSC
	// can separate an external authority base from the measured target root.
	path_base:  "repository-root"
	activation: "source-checkout-only"
	target_root_contract: "explicit-coh-root"
	target_revision?: string & =~"^[0-9a-f]{40}$"
	output_root:  string & =~"^\\.tsc(/|$)"
	default_mode: "mechanical" | "llm" | "hybrid" | "auto"

	consistency: {
		mechanical: "identical"
		llm_repeats: int & >=2
		llm_spread:  !=""
		script?:     !=""
		...
	}

	admissibility?: !=""

	standing?: {
		scope:                 "house-authored-public-commons" | "house-authored-blind-heldout" | "external-blind-heldout"
		admissibility:         "public-only" | "public-plus-house-heldout" | "public-plus-heldout"
		heldout_status:        "none" | "registered-and-revealed"
		external_anchor_count: int & >=0
		llm_consistency_gate:  "reported-not-gating" | "passed" | "failed"
		llm_consistency_floor: float & >=0 & <=1
		...
	}

	mechanical: {
		backend:     !=""
		determinism: !=""
		signals: {
			alpha: [...string]
			beta:  [...string]
			gamma: [...string]
		}
	}

	llm: {
		estimates:  [...string]
		must_not:   [...string]
		validation: !=""
		providers: {
			local: !=""
			ci:    !=""
			...
		}
		ci_prompt: !=""
		...
	}

	...
}

#Measurement: #Skill & {
	artifact_class: "measurement"
	methodology:     #CoherenceMethodology
}
