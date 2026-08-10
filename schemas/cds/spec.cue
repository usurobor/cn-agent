// Schema ID: cnos.cdd.cds.spec.v1
//
// #CDSCellSpec — the CDS input-side overlay on the generic #CellSpec. It pins
// the declared protocol_id, the canonical diff-first evidence rule, and the
// CLOSED shapes of the CDS seats: the fill-owned construction boundary means
// the generic schema only sees a tagged envelope, while this overlay — like
// the cds.patch Go decoder it mirrors — strictly defines every constructor
// argument. Unknown, mixed-case, and null fields still fail.
//
// NOTE (Pi #32 D1): declaring protocol_id "cnos.cdd.cds.receipt.v1" is a
// declaration of intent, not proof. The v0 runner carries it as provenance and
// sets protocol_validated=false.
package cds

import "cnos.dev/cnos/schemas/cdd"

// #Hole is an unresolved `$param` reference. Authored form may carry holes in
// argument positions; resolution replaces them in place, so this closed shape
// vets both — enum-constrained fields admit a hole explicitly.
#Hole: =~"^\\$[A-Za-z_][A-Za-z0-9_]*$"

// #CDSPatchAlpha is the complete cds.patch alpha declaration.
#CDSPatchAlpha: {
	fill: "cds.patch"
	cognition: {
		provider: "claude-cli" | "codex-cli" | "fake" | #Hole
		model:    string // exact model id; may be empty only for "fake"
	}
	workspace: {
		kind:     "git-worktree"
		repo:     string & !=""
		base_sha: string & !=""
	}
	skills: [string & !="", ...string & !=""] // ordered canonical refs (or holes)
}

// #CDSMechanicalUnmetBeta is Case 2's honest reviewer: a mechanical seat that
// cannot judge the goal and therefore never passes it. Case 3 replaces beta
// alone — the fill boundary makes that a one-field change.
#CDSMechanicalUnmetBeta: {
	fill: "cdd.mechanical-unmet"
}

#CDSCellSpec: cdd.#CellSpec & {
	protocol_id: "cnos.cdd.cds.receipt.v1"

	// CANONICAL ORDER (explicit rule, not an accident): a CDS spec's first
	// required_evidence entry IS the alpha diff. Chosen over order-independent
	// membership deliberately: the structural list pattern below enforces the
	// rule through plain unification, whereas computed validators
	// (list.MatchN/Contains, comprehension guards) did not fire reliably
	// through `vet -d` against a closed base definition on the CI-pinned cue
	// (v0.11.0). A diff present but not first is a schema violation
	// (fixtures/invalid/cds-diff-not-first.json).
	contract: required_evidence: [{id: "diff", kind: "diff", producer: "alpha"}, ...cdd.#RequiredRef]

	alpha: #CDSPatchAlpha
	beta:  #CDSMechanicalUnmetBeta
}
