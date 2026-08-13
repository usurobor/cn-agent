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
// cognition must name a REQUESTED MODEL SELECTOR; only the deterministic fake
// may omit it. Written as a disjunction rather than two independent fields so this
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

// #CognitionAuthored is the AUTHORING-time shape: what is STRUCTURALLY
// POSSIBLE before resolution. It is deliberately wider than what will
// succeed — a provider hole with the model omitted may resolve to claude-cli
// and then fail construction, so this shape admits it and `Resolve` plus the
// fill's constructor judge the selected combination (Pi #59 B1). Three
// resolution stages exist, so there are three arms:
//
//   - a literal fake ignores the model, so it may omit the key, write it
//     empty, or carry a hole that resolution fills with the empty value;
//   - a literal real provider needs a model selector, present either as a
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

// _patchAlphaSkillsResolved is the ordered canonical refs with the content
// digest of the body that was actually injected — naming a skill is not
// loading it. Named once because two definitions below carry it.
_patchAlphaSkillsResolved: [{ref: string & !="", sha256: =~"^[0-9a-f]{64}$"}, ...{ref: string & !="", sha256: =~"^[0-9a-f]{64}$"}]

// #CDSPatchAlphaResolved is what a closure records: the canonical structural
// shape and digested skills. Hole-freedom is proven by resolution, not by this
// shape — see the header.
//
// NO WORKSPACE. A cds.patch seat does not name a repository: it reads the
// repository and the base commit from the run's pinned contract subject, which
// is the single source. The key is absent from a CLOSED definition, so a
// declaration carrying one is rejected here exactly as the fill's exact-key set
// rejects it in Go — the two authorities delete the field together, or the
// deletion is only a Go convention.
#CDSPatchAlphaResolved: {
	fill:      "cds.patch"
	cognition: #Cognition
	skills:    _patchAlphaSkillsResolved
}

// #CDSPatchAlphaAuthored is what a cell spec may carry: the same shape with
// holes admitted in the positions resolution fills.
#CDSPatchAlphaAuthored: {
	fill: "cds.patch"
	cognition: #CognitionAuthored
	skills: [#Concrete | #Hole, ...#Concrete | #Hole]
}

// #CDSPatchAlphaResolvedPreWorkspaceDeletion is a FROZEN HISTORICAL shape, and
// exists for exactly one artifact: docs/architecture/evidence/
// cds-case2-claude-closure.json, the one committed record of a rented-cognition
// run. That episode ran while the fill still declared its own workspace, and
// its record is covered by a scope-lift digest that recomputes — editing the
// declaration out of it would break the digest and turn evidence into a claim.
// So the artifact stays byte-for-byte, and this definition is what keeps it
// accountable to a CLOSED shape rather than to nothing.
//
// It is not a second way to name a repository. Nothing authorable references
// it: #CDSCellSpec.alpha is #CDSPatchAlphaAuthored, the runtime has no decoder
// for a workspace key, and no new document can be produced in this shape. It
// describes a record that already exists and cannot be regenerated.
#CDSPatchAlphaResolvedPreWorkspaceDeletion: {
	fill:      "cds.patch"
	cognition: #Cognition
	skills:    _patchAlphaSkillsResolved
	workspace: {
		kind:     "git-worktree"
		repo:     string & !=""
		base_sha: =~"^[0-9a-f]{40}$"
	}
}

// #CDSMechanicalUnmetBeta is Case 2's honest reviewer: a mechanical seat that
// cannot judge the goal and therefore never passes it. Case 3 replaces beta
// alone — the fill boundary makes that a one-field change.
#CDSMechanicalUnmetBeta: {
	fill: "cdd.mechanical-unmet"
}

// #GitSnapshotPinned is a subject as a RECORD carries it: one repository at
// one exact commit. It mirrors cellwork.Subject field for field, and the two
// are vetted against ONE corpus (schemas/cds/fixtures/subject/), so a subject
// admitted by one authority and rejected by the other is a gate failure rather
// than a discovery in production.
//
// It lives in this package because the CDS run input is where a subject is
// required and CUE has no separate adapter package; the LANGUAGE is the git
// subject adapter's, not CDS's — cellwork owns it on the Go side.
//
// `!` throughout, for the reason #CDSIssue carries it: without the required-
// field marker an ABSENT field unifies with its declared value and vets clean.
#GitSnapshotPinned: {
	kind!:     "git.snapshot/0.1"
	repo!:     string & !=""
	base_sha!: =~"^[0-9a-f]{40}$"
}

// #GitSnapshotAuthored is what a run input may carry: the same shape with a
// base that may still be a moving revision. Pinning happens once, before
// either seat is constructed; the record then carries the pinned form above.
// This pair is what shows the two definitions are not the same definition.
#GitSnapshotAuthored: {
	kind!:     "git.snapshot/0.1"
	repo!:     string & !=""
	base_sha!: string & !=""
}

// #NonBlank is a string carrying at least one non-whitespace rune, and it is
// the ONE blankness predicate the issue and design schemas have: `!=""` is not
// it, since "   " satisfies that while reading as an authored sentence to
// every human and to neither authority.
//
// The class is transcribed CHARACTER FOR CHARACTER from cdsissue's
// nonBlankPattern, which compiles the same string with Go's regexp. Two
// hand-written predicates were what diverged before: `=~"\\S"` here beside
// `strings.TrimSpace` there accepted a NO-BREAK-SPACE-only field in CUE and
// rejected it in Go, because `\s` in RE2 is only `[\t\n\f\r ]` while
// TrimSpace strips everything unicode.IsSpace covers. The enumeration below
// is that set; issue-blank-unicode-whitespace.json carries all of it in one
// field, so a transcription slip fails the corpus instead of production.
#NonBlank: string & =~#"[^\t\n\v\f\r \x{0085}\x{00A0}\x{1680}\x{2000}-\x{200A}\x{2028}\x{2029}\x{202F}\x{205F}\x{3000}]"#

// #CDSIssue is the CDS typed issue: the problem half of the run contract. It
// mirrors cdsissue.Issue field for field, and the two are vetted against ONE
// corpus (schemas/cds/fixtures/issue/) so a document admitted by one authority
// and rejected by the other is a gate failure rather than a discovery in
// production.
//
// Closed: unknown and mixed-case keys fail here, and cdsissue.Admit closes the
// same key language with cellfill.OnlyKeys at every level, because
// encoding/json would otherwise match `Kind` case-insensitively.
#CDSIssue: {
	// `!` throughout — the required-field marker, same reason `#Seat.fill!`
	// carries it: without it an ABSENT field unifies with its declared value
	// and vets clean, so an issue document omitting `kind` entirely would pass
	// CUE while cdsissue.Admit rejects it. Absence must be rejected, not
	// defaulted.
	kind!: "cnos.cds.issue.v0"
	id!:   #NonBlank

	// The incoherence in three lines.
	problem!: {
		exists!:   #NonBlank
		expected!: #NonBlank
		diverges!: #NonBlank
	}

	// One canonical path per load-bearing claim; at least one.
	sources!: [#CDSSource, ...#CDSSource]

	// The execution boundary. `out` is REQUIRED but may be empty: non-goals are
	// load-bearing, so an author must have considered them, and an empty list
	// says "considered, none" where an absent key says nothing at all. Go reads
	// the same present-vs-absent distinction off the raw document, since both
	// decode to an empty slice.
	scope!: {
		in!: [#NonBlank, ...#NonBlank]
		out!: [...#NonBlank]
	}

	// At least one criterion, each naming its verification route, with ids
	// unique. Uniqueness is expressed as a length agreement rather than a
	// membership validator: comprehension-built key sets unify through
	// `vet -d` against a closed definition on the CI-pinned cue, where computed
	// list validators did not (the same constraint that made required_evidence
	// order structural).
	acceptance!: [#CDSCriterion, ...#CDSCriterion]
	_acceptanceIDs: {for c in acceptance {(c.id): true}}
	_uniqueAcceptanceIDs: len(_acceptanceIDs) & len(acceptance)
}

#CDSSource: {
	claim!: #NonBlank
	path!:  #NonBlank
}

// A criterion without a verification route is exactly the ill-defined
// criterion this schema exists to reject: it leaves beta judging plausibility.
#CDSCriterion: {
	id!:           #NonBlank
	statement!:    #NonBlank
	verification!: #NonBlank
}

// #CDSDesign is the CDS typed design: the change half of the run contract, and
// a SEPARATE document from the issue so that neither can silently become the
// other. It mirrors cdsdesign.Design field for field against ONE corpus
// (schemas/cds/fixtures/design/).
//
// The shape is the minimum a structural gate can decide. Alternatives,
// trade-offs and known debt belong to a design as CELL-SYSTEM-DESIGN §10.1
// describes it and are deliberately absent here: every one of them is prose
// this authority could only check for non-emptiness, which would grow the
// authoring burden without growing what admission proves.
#CDSDesign: {
	kind!:     "cnos.cds.design.v0"
	approach!: #NonBlank
	// At least one thing the change must not break, and at least one surface
	// it reaches: a design that names neither has not been thought about.
	invariants!: [#NonBlank, ...#NonBlank]
	impact!: [#CDSImpact, ...#CDSImpact]
}

#CDSImpact: {
	surface!: #NonBlank
	why!:     #NonBlank
}

// #CDSRunInput is the per-run document `cn cell run --input` reads: the issue,
// the design, and the subject the episode acts on. It is separate from
// #CDSCellSpec because a cell definition is reusable while these three belong
// to one run — folding them in would make every run author a new cell.
//
// All three are REQUIRED here. The Go door distinguishes an absent payload
// (incomplete) from a malformed one (rejected); CUE has one verdict, so this
// schema states only that a complete run input carries all three.
#CDSRunInput: {
	kind!:    "cnos.cds.run-input.v0"
	issue!:   #CDSIssue
	design!:  #CDSDesign
	subject!: #GitSnapshotAuthored
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
