// schemas/skill.cue — CUE schema for SKILL.md frontmatter (CTB v0.1).
//
// Validated by tools/validate-skill-frontmatter.sh as the I5 coherence-CI
// job. Field shape, type, and enum constraints live here; file discovery,
// frontmatter extraction, exception handling, and filesystem path checks
// are owned by the script (#301 AC2 surface boundary).
//
// LANGUAGE-SPEC §11 says loaders MUST ignore unknown keys, so this schema
// is open (`...` at the end) and does NOT reject package-local extension
// keys (e.g. `parent:`).
//
// Hard-gate fields fail the build if missing (no exception path). Fields
// declared optional here as `?` — `artifact_class`, `kata_surface`,
// `inputs`, `outputs` — are spec-required-but-exception-backed in §301:
// the script enforces "present-or-excepted-in-skill-exceptions.json" for
// those fields after `cue vet` succeeds on type/shape.

package skill

#Skill: {
	// Hard gate (#301 AC1): must be present, no exception.
	// Schema requires these unconditionally; `cue vet` fails if missing.
	//
	// `name` accepts a slash to allow compound sub-skill names like
	// `review/contract` and `review/issue-contract` that the cdd package
	// uses for orchestrator/sub-skill nesting.
	name:               string & =~"^[a-z][a-z0-9_/-]*$"
	description:        !=""
	governing_question: !=""
	triggers: [...string]
	scope: "global" | "role-local" | "task-local"

	// Spec-required-but-exception-backed (#301 AC1, AC4).
	// Schema marks optional; the script enforces presence-or-exception.
	artifact_class?: "skill" | "runbook" | "reference" | "deprecated"
	kata_surface?:   "embedded" | "external" | "none"
	inputs?: [...string]
	outputs?: [...string]

	// Optional with default per LANGUAGE-SPEC §4.1 (loader defaults to
	// `internal` when omitted). Schema validates the enum if present.
	visibility?: "internal" | "public"

	// Reserved per LANGUAGE-SPEC §11 — type-checked only when present.
	requires?: [...string]

	// `calls`: each entry is either a bare path string or a mapping with
	// at least a `path` field (LANGUAGE-SPEC §2.4.1). Filesystem-existence
	// of each target is checked by the script (#301 AC6), not here.
	calls?: [...(string | {
		path: string
		...
	})]

	// `calls_dynamic`: declares the *source* of dynamic targets, not the
	// targets themselves (LANGUAGE-SPEC §2.4.2). Target existence is NOT
	// validated by I5 (§301 non-goal). Only source/constraint shape.
	calls_dynamic?: [...{
		source:      string
		constraint?: string
		...
	}]

	runs_after?: [...string]
	runs_before?: [...string]
	excludes?: [...string]

	// LANGUAGE-SPEC §11: loaders MUST ignore unknown keys. The schema is
	// open: package-local extension keys (e.g. `parent:`) pass through.
	...
}
