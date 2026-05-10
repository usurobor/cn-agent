## CDD Iterations — cycle #338

**Issue:** #338 — cdd: Add §1.6c — dispatch sizing, prompt scope, and commit checkpoints
**Date:** 2026-05-10
**Mode:** docs-only

---

### F1: initial-dispatch-no-sizing-rule

- **Source:** γ-triage / issue #338 body (Problem statement + Status truth table)
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — pattern observed in N=4 dispatch failures (cycle #335, TSC supercycle 3/5 α re-dispatches); no existing rule in `CDD.md §1.6` for scaling the initial dispatch timeout budget with work complexity
- **Description:** `CDD.md §1.6` covered the sequential bounded dispatch coordination model and §1.6b covered re-dispatch prompt complexity, but there was no rule governing the initial dispatch: no budget heuristic tied to AC count, no prompt-scope guidance, and no commit-checkpoint mandate. When agents were SIGTERMed before committing, recovery depended on operator memory rather than a documented procedure. The gap produced recoverable failures in cycle #335 (18 files recovered from worktree, 0 commits at SIGTERM) and unrecoverable failures in 3/5 TSC supercycle α close-out re-dispatches.
- **Root cause:** §1.6b was written to address re-dispatch failures (TSC supercycle observation), but the family — context-load vs task-complexity mismatch — was not addressed for initial dispatches. The initial dispatch was left with a fixed operator-chosen budget and no structural guidance. The commit-checkpoint discipline was not encoded anywhere in the skill files.
- **Disposition:** `patch-landed`
- **Patch:** commit `69de7ef8` (`CDD.md §1.6c`), commit `b2f5ee3b` (`operator/SKILL.md §7 timeout-recovery`), commit `482f1c81` (`post-release/SKILL.md §4 telemetry fields`)
- **Affects:** `cdd/CDD.md §1.6c`; `cdd/operator/SKILL.md §7`; `cdd/post-release/SKILL.md §4`
