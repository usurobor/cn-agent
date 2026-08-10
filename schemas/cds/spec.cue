// Schema ID: cnos.cdd.cds.spec.v1
//
// #CDSCellSpec — the CDS input-side overlay on the generic #CellSpec. It pins
// the declared protocol_id, the canonical diff-first evidence rule, and the
// CLOSED shapes of the CDS seats: the fill-owned construction boundary means
// the generic schema only sees a tagged envelope, while this overlay — like
// the cds.patch Go decoder it mirrors — strictly defines every constructor
// argument. Unknown, mixed-case, and null fields still fail.
//
// AUTHORED vs RESOLVED are separate definitions on purpose, and each proves a
// different thing — stated exactly, because the earlier wording overclaimed
// (Pi #57 C1):
//
//   - RUNTIME RESOLUTION proves every authored reference was filled. That is
//     the only authority that can prove it, because it is the only one that
//     saw the parameters.
//   - THE RESOLVED SCHEMA proves canonical structural shape, a base_sha
//     pinned to a real commit rather than a moving ref, and skills carrying
//     content digests.
//
// The resolved schema does NOT mechanically prove "no holes remain". A
// supplied parameter value may legitimately begin with `$` — it is then a
// resolved literal, and the record keeps no provenance letting CUE tell it
// from an unresolved token. Banning such values outright would break honest
// inputs to satisfy a sentence, so the sentence is what changed.
//
// NOTE (Pi #32 D1): declaring protocol_id "cnos.cdd.cds.receipt.v1" is a
// declaration of intent, not proof. The v0 runner carries it as provenance and
// sets protocol_validated=false.
package cds

import "cnos.dev/cnos/schemas/cdd"

// #Hole is an unresolved `$param` reference. Its identifier grammar is the
// generic one — a hole IS a parameter name, so the two cannot diverge.
#Hole: =~"^\\$[A-Za-z_][A-Za-z0-9_]*$"

// #Concrete is an already-resolved value: nonempty, and NOT hole-shaped.
// Written as an explicit exclusion because Go treats EVERY `$...` string as a
// hole. A plain `string & !=""` arm silently accepted `$bad-name` — the value
// looks concrete to CUE and is a malformed hole to Go, so the two authorities
// disagreed about the same document (Pi #56 C1). Anywhere a field may carry
// either, the union must be #Concrete | #Hole, never string | #Hole.
#Concrete: string & !="" & !~"^\\$"

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

// #CognitionAuthored is the same rule at AUTHORING time, and it has to admit
// exactly what survives resolution — no more and no less (Pi #57 C1). Three
// stages exist, so there are three arms:
//
//   - a literal fake ignores the model, so it may omit the key, write it
//     empty, or carry a hole that resolution fills with the empty value;
//   - a literal real provider needs an exact model, present either as a
//     concrete literal or as a hole;
//   - a HOLE in the provider position means which arm applies is not known
//     until resolution, so the model must admit the union of both.
//
// Every model position is `#Concrete | #Hole` rather than bare `string`. A
// bare string arm accepted `$bad-name`, which CUE read as a concrete value
// and Go read as an illegal hole — the same divergence #Concrete closed for
// workspace and skills.
//
// The RESOLVED form above still requires `model: ""` present for the fake,
// because a receipt records what held the seat rather than what the author
// chose not to type.
#CognitionAuthored: {provider: "fake", model?: "" | #Hole} |
	{provider: "claude-cli", model: #Concrete | #Hole} |
	{provider: #Hole, model?: "" | #Concrete | #Hole}

// #CDSPatchAlphaResolved is what a closure records: the canonical structural
// shape, a base_sha pinned to a commit at construction, and digested skills.
// Hole-freedom is proven by resolution, not by this shape — see the header.
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
	cognition: #CognitionAuthored
	workspace: {
		kind:     "git-worktree"
		repo:     #Concrete | #Hole
		base_sha: #Concrete | #Hole
	}
	skills: [#Concrete | #Hole, ...#Concrete | #Hole]
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
