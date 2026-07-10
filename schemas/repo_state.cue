// schemas/repo_state.cue — CUE schema for `.cn/repo.state.json` (cnos#656,
// Phase 1 of the cnos#655 `cn repo` lifecycle wave).
//
// Types the managed-surface ownership ledger `cn repo install` writes and
// `cn repo status`/`doctor`/`uninstall` read. Design surface:
// docs/development/design/cn-repo-lifecycle.md (§Phase-contract precision
// P1-P3 is the normative AC surface this schema satisfies).
//
// Validated by scripts/ci/validate-repo-state.sh, mirroring
// scripts/ci/validate-skill-frontmatter.sh's schema-owns-shape /
// script-owns-discovery split for schemas/skill.cue.
//
// Unlike schemas/cdd/receipt.cue's #Receipt or schemas/skill.cue's #Skill
// (both deliberately OPEN — they type externally-authored, forward-compatible
// formats), #RepoState and its nested definitions are CLOSED (the CUE
// default for a #Definition with no trailing `...`). `.cn/repo.state.json`
// is a machine-written ledger, entirely owned by `cn repo install`/`update`/
// `repair` — never hand-authored — so closedness is the right default here,
// and it is what makes two P1 invariants structural instead of
// script-enforced:
//
//   - timestamp-free: no timestamp-shaped field is declared anywhere below,
//     so a fixture adding `timestamp`/`created_at`/`installed_at`/etc. fails
//     `cue vet` as an unknown-field error under closedness.
//   - no lock duplication (design doc A1): there is no top-level `packages`
//     field declared, so a fixture adding one (with `version`/`sha256`)
//     fails the same way. Package identity is referenced by `name` string
//     only (managed_dirs[].package), joined against `.cn/deps.lock.json`
//     (schema cn.lock.v2) at read time by Go code — never copied here.
//
// managed_files entries require path/kind/id/sha256 uniformly (not just
// workflow records) — see gamma-scaffold.md Decision 1 (cnos#656) for why
// this reads the issue's P1 Acceptance text over the design doc's earlier
// illustrative JSON, which omitted id/sha256 for manifest/lockfile entries.

package repostate

#RepoState: {
	schema:  "cn.repo.state.v1"
	profile: string

	source: #Source

	managed_files: [...#ManagedFile]
	managed_dirs: [...#ManagedDir]

	external_expectations: {
		labels: {
			mode:                 "ensure" | "ignore"
			source:               string
			delete_on_uninstall:  bool
		}
	}
}

#Source: {
	channel: string
	release: string
	index:   string
}

// #ManagedFile types one entry in `managed_files`. Every entry requires
// path/kind/id/sha256 (P1). Workflow entries additionally require the full
// render contract (P2) so drift can be classified as matches-ledger /
// user-edit / renderer-moved without re-deriving render inputs blind.
#ManagedFile: {
	path:   string
	kind:   "manifest" | "lockfile" | "workflow" | "gitignore" | "other"
	id:     string
	sha256: string & =~"^[0-9a-f]{64}$"

	if kind == "workflow" {
		tier:                    "agent" | "engine"
		renderer:                string
		renderer_package:        string
		renderer_version_source: string
		agent:                   string
		workflow_pat_secret?:    string
		bot_name?:               string
		bot_id?:                 string
	}
}

// #ManagedDir types one entry in `managed_dirs` — a materialized vendor
// directory mapped to the package that owns it. No version/sha256 here
// (A1): package identity is looked up from `.cn/deps.lock.json` by `package`.
#ManagedDir: {
	path:    string
	package: string
}
