// schemas/skill.cue — CUE schema for SKILL.md frontmatter (CTB v0.1).
//
// Validated by scripts/ci/validate-skill-frontmatter.sh as the I5 coherence-CI
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
	artifact_class?: "skill" | "runbook" | "reference" | "deprecated" | "wake"
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

// #WakeOutputAdmin — role-shaped output block for wake.role == "admin".
// Required fields: channel_log_convention, writer_surface, class_taxonomy,
// cursor_advance, cursor_field. Open for renderer extensions.
#WakeOutputAdmin: {
	channel_log_convention: !=""
	writer_surface:         !=""
	class_taxonomy: [...string]
	cursor_advance: bool
	cursor_field:   !=""
	...
}

// #WakeOutputDispatch — role-shaped output block for wake.role == "dispatch".
// Required fields: cycle_artifact_root, artifact_class_taxonomy, cell_runtime.
// Open for renderer extensions.
#WakeOutputDispatch: {
	cycle_artifact_root:     !=""
	artifact_class_taxonomy: [...string]
	cell_runtime:            !=""
	...
}

// #Wake — typed wake module schema (cnos#524 W1).
//
// A wake is a typed SKILL.md module whose frontmatter carries a `wake:` block
// and whose body is the prompt. The `#Wake` definition validates the `wake:`
// block shape at I5 time; the renderer enforces runtime behavior.
//
// Embeds #Skill and constrains artifact_class + scope as required by the
// wake-as-skill design (w0-design.md §B, §C decisions 2 and 3).
//
// Role enum: "admin" | "dispatch" ONLY (OB-1 — no "observer").
// agent_variable.default: string | null (FN-3 — null for admin wake).
// surfaces.allowed / .disallowed: [...string] (FN-4 — no length constraint).
#Wake: #Skill & {
	artifact_class: "wake"
	scope:          "global"

	wake: {
		// role-discriminant: drives role-shaped output disjunction.
		// "observer" is NOT a valid schema value (OB-1).
		role: "admin" | "dispatch"

		// package owning this wake (e.g. "cnos.core", "cnos.cds")
		package: !=""

		// required on both roles
		admin_only:             bool
		activation_log_writer:  bool

		// dispatch-only fields (optional on admin)
		activation_state?: "live" | "declaration-only"
		protocol?:         string
		selector?: {
			include: [...string]
			exclude: [...string]
		}

		// input contract
		input: {
			triggers:                     [...string]
			issues_opened_title_pattern?: string
		}

		// role-shaped output disjunction:
		// admin wakes carry channel_log_convention;
		// dispatch wakes carry cycle_artifact_root.
		output: #WakeOutputAdmin | #WakeOutputDispatch

		// logical permissions
		permission_intent: [...string]

		// concurrency
		concurrency: {
			serialize: bool
			group:     !=""
		}

		// agent_variable.default is string | null:
		// null → operator-required (admin wake pattern)
		// string → operator-defaultable (dispatch wake pattern)
		agent_variable: {
			name:    !=""
			default: string | null
		}

		// allowed/disallowed write surfaces
		surfaces?: {
			allowed?:    [...string]
			disallowed?: [...string]
		}

		// defer-path routing doctrine (optional)
		defer_path?: {
			cell_shaped_directive?: string
			off_role_directive?:    string
			ambiguous_directive?:   string
		}

		// Open per LANGUAGE-SPEC §11: renderer fields not yet in #Wake pass through.
		...
	}
}
