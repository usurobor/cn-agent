// Schema ID: cnos.cdd.cds.spec.v1
//
// #CDSCellSpec — the CDS input-side overlay on the generic #CellSpec. It pins
// the declared protocol_id, the canonical diff-first evidence rule, and the
// CLOSED shapes of the CDS seats: the fill-owned construction boundary means
// the generic schema only sees a tagged envelope, while this overlay — like
// the cds.patch Go decoder it mirrors — strictly defines every constructor
// argument. Unknown, mixed-case, and null fields still fail.
//
// AUTHORED vs RESOLVED are separate definitions on purpose. An authored spec
// may carry `$param` holes; a resolved declaration — the one the closure
// records and the digest covers — may not, and its base_sha must be a real
// commit rather than a moving ref. One schema serving both would let a hole
// survive into a receipt unnoticed, so the corpus vets the emitted closure's
// declaration against the RESOLVED shape.
//
// NOTE (Pi #32 D1): declaring protocol_id "cnos.cdd.cds.receipt.v1" is a
// declaration of intent, not proof. The v0 runner carries it as provenance and
// sets protocol_validated=false.
package cds

import "cnos.dev/cnos/schemas/cdd"

// #Hole is an unresolved `$param` reference.
#Hole: =~"^\\$[A-Za-z_][A-Za-z0-9_]*$"

// #Cognition is the inline provider declaration. A provider that really rents
// cognition must name an EXACT model; only the deterministic fake may omit
// it. Written as a disjunction rather than two independent fields so this
// schema rejects exactly what cellcog.New rejects — the two authorities must
// not disagree. There is no argv/flags escape: a cell cannot smuggle
// arguments into an adapter.
//
// codex-cli is deliberately absent (Pi #55 D1). Its available suppression
// flags reach $CODEX_HOME/config.toml and execpolicy .rules only, while
// global/project AGENTS.md and discovered skills still load — a second,
// unreceipted component definition beside the fill's digested skills. It is
// held here exactly as it is held in cellcog.New.
#Cognition: {provider: "fake", model: ""} |
	{provider: "claude-cli", model: string & !=""}

// #CDSPatchAlphaResolved is what a closure records: no holes, and a base_sha
// pinned to a commit at construction.
#CDSPatchAlphaResolved: {
	fill:      "cds.patch"
	cognition: #Cognition
	workspace: {
		kind:     "git-worktree"
		repo:     string & !=""
		base_sha: =~"^[0-9a-f]{40}$"
	}
	// Ordered canonical refs with the content digest of the body that was
	// actually injected — naming a skill is not loading it.
	skills: [{ref: string & !="", sha256: =~"^[0-9a-f]{64}$"}, ...{ref: string & !="", sha256: =~"^[0-9a-f]{64}$"}]
}

// #CDSPatchAlphaAuthored is what a cell spec may carry: the same shape with
// holes admitted in the positions resolution fills.
#CDSPatchAlphaAuthored: {
	fill: "cds.patch"
	cognition: #Cognition | {provider: #Hole, model: string}
	workspace: {
		kind:     "git-worktree"
		repo:     (string & !="") | #Hole
		base_sha: (string & !="") | #Hole
	}
	skills: [(string & !="") | #Hole, ...(string & !="") | #Hole]
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

	alpha: #CDSPatchAlphaAuthored
	beta:  #CDSMechanicalUnmetBeta
}
