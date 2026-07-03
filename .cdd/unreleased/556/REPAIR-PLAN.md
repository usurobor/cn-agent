# REPAIR-PLAN — cnos#556 (R2)

## Context

The operator reopened #556 with `status:changes` and posted a correction
(issue comment, 2026-07-03) stating that the R1 repair round was itself
wrong: it reverted the physical relocation of the Go implementation into
`src/packages/cnos.issues/commands/issues-map/` on the strength of an
earlier operator comment ("the Go implementation may remain under
`src/go/internal/issuesmap/`") that the operator has now explicitly
withdrawn ("That was a misunderstanding of intent and is now withdrawn.").
The operator's final instruction: move the Go implementation into the
package and **keep it there**; do not revert the relocation again.

R0 (commits `c5c7db79`..`bdad537f`) already implemented exactly this,
mirroring the `cnos#392` `cdd-verify` precedent, and β converged on it at
R0. R1 (commits `f9707a9d`, `97a4b7d3`, `aa635c7a`) reverted that work in
response to a δ override that, per the operator's correction, misweighted
a withdrawn comment over the issue's actual intent.

## Rejected findings → planned repairs

| # | Rejected finding | Source | Planned repair | Evidence target |
|---|---|---|---|---|
| 1 | Go implementation must physically live under `src/packages/cnos.issues/commands/issues-map/`, not `src/go/internal/issuesmap/` | Operator `status:changes` comment (2026-07-03), "Required implementation" items 1-2 | `git revert` the two R1 commits that moved the code back (`f9707a9d`, `97a4b7d3`), reinstating R0's relocation, `go.work` entry, and `src/go/go.mod` require+replace wiring | `find src/go/internal/issuesmap src/packages/cnos.issues -type f`; `git diff` vs R0 state |
| 2 | `cmd_issues_map.go` must stay a thin shim over the package-owned implementation | Operator item 4; precedent `cmd_cdd_verify.go` | Confirmed by repair #1 (the revert restores R0's thin-shim import); no additional change needed beyond re-verifying the file is still parse+delegate only | `cmd_issues_map.go` diff/contents |
| 3 | `cnos.issues` doctrine (`SKILL.md`, `skills/map/SKILL.md`) must state the implementation lives in the package now, not narrate the withdrawn rule as current | Operator item 5 ("Teach current truth — no R0/R1 narration in active skill prose") | Do not simply re-revert `aa635c7a` (which would leave R0/R1 narration in prose per its own original wording); rewrite `SKILL.md` + `skills/map/SKILL.md` fresh to state current truth only: implementation at `src/packages/cnos.issues/commands/issues-map/`, kernel-dispatched thin shim (Option B, `#216` debt), no R0/R1 history in the active doctrine text (history belongs in this repair record / CDD artifacts, not skill prose) | Read diff of both `SKILL.md` files; grep for "R0"/"R1" in active skill prose (should be absent) |
| 4 | Board Action must still call the public `cn issues map` command, no behavior change | Operator items 6-9 | No planned change — confirm `.github/workflows/board-map.yml` is untouched by the repair (it never referenced internal paths) | `git diff` shows no change to `board-map.yml` |
| 5 | `cn build --check`, `go test ./...`, and all required CI jobs (I1/I2/I4/I5/I6/Go/Package/Binary + cnos#516/#524 guards) must pass at the new HEAD | Operator "Acceptance" list | Re-run local build/vet/test + `cn build --check` + `cn issues map` (fixture and live) after the revert and doctrine fix; push and confirm remote CI green before any closeout | local command output; `gh run view` on the pushed head |
| 6 | Receipt must explicitly state package-command dispatch disposition (Option A vs B) | Operator "Acceptance" item + AC10 | State Option B explicitly in the rewritten `SKILL.md`: implementation is package-owned and physically relocated, but dispatch remains the compiled-in kernel command (thin shim) — true package-command exec-dispatch is `#216` debt, matching the `cdd-verify` precedent's own disposition | `SKILL.md` "Command-dispatch disposition" section |

## Non-goals for this repair round

- Do not reopen AC1-AC9 as previously verified at R0/R1 (board output, taxonomy, CLI dispatch shape, GitHub Action are unaffected by this round).
- Do not implement generic `#216` package-command exec-dispatch.
- Do not add a `commands` entry to `cnos.issues/cn.package.json` (precedent: `cnos.cdd/cn.package.json` has no `cdd-verify` entry either — physical co-location is not the same as declared package-command dispatch).

## Sequencing

1. This file (write-before-any-other-file gate).
2. `git revert` the two R1 revert commits on `cycle-556-local`.
3. Rewrite `SKILL.md` + `skills/map/SKILL.md` doctrine text (fresh, not a re-revert of `aa635c7a`).
4. Local verification: build/vet/test, `cn build --check`, `cn issues map` (fixture + live), dangling-reference grep.
5. `self-coherence.md` §R2 + `beta-review.md` §R2 (independent re-verification) + closeouts, each carrying a `repair_evidence` block.
6. Push to `cycle/556`, confirm PR #557 reflects new head, wait for remote CI, confirm all required jobs green.
7. Only then: `status:in-progress → status:review` with `deliverable_evidence`.
