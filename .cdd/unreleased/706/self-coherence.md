## Gap

Issue: [usurobor/cnos#706](https://github.com/usurobor/cnos/issues/706) — "CDS install preflight-first: ask for operator prerequisites before doing anything, stop inventing a bot."

Version/mode: `design-and-build` (per γ scaffold — the design is converged in the issue's "CONSOLIDATED FINAL SPEC" comment, but not filed at a stable `docs/{tier}/{bundle}/{X.Y.Z}/DESIGN.md` path, so MCA preconditions are not met; this is not a blocker, just the correct mode label).

Governing gap (restated from the consolidated spec, authoritative over the original issue body per its own "where it differs, this wins" clause): `cn repo install --dispatch cds` did agent-doable work first (render/labels/commit) and surfaced operator-only gates (secrets, push access, merge-to-default-branch) last, invented a non-existent "bot" concept, and named the workflow-PAT secret opaquely (`SIGMA_WORKFLOW_PAT`, hardcoding the agent name). Fix: ask the operator for what only they can provide, explain exactly how to get it, before doing anything — and stop inventing a bot.

This cycle (R0) implements the consolidated spec's Deliverables 1–7 against Final ACs 1–10, working from `cycle/706` (cut from `main` HEAD `7f249ddbb50f230d5d41287b6554ab17b5a1d1d5`) on top of γ's scaffold (`.cdd/unreleased/706/gamma-scaffold.md`, commit `3bf1b2d`).

## Skills

Active skills this cycle (Tier 1/2/3, per `cnos.cdd/skills/cdd/alpha/SKILL.md`):

- `cdd/alpha/SKILL.md` §2.1 (dispatch intake — branch checkout, reading the scaffold as the authoritative contract), §2.5 (this file's incremental, canonical-header-form write discipline), §2.6 (pre-review gate, applied before signaling review-readiness).
- No `cnos.handoff/skills/handoff/dispatch/SKILL.md` 7-axis contract discipline — γ's scaffold explicitly waives it for this cycle ("no axis is undecidable in a way that blocks starting").
- `eng/go` conventions (implicit, not a loaded skill file this session, but followed): cli/-boundary compliance (`cli/cmd_repo_install.go` stays a thin wrapper; all domain logic in `internal/repoinstall`), dependency-free `net/http` GitHub REST idiom (mirrored from `label-doctor/github.go`, not imported — separate go.work module).
