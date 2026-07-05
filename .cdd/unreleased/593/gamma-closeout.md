# γ close-out — cnos#593

## Cycle summary

Sub C of #583: a mechanical recovery scanner for stranded CDS dispatch cells. The gamma-scaffold's "Key finding" — that the six behavior cases the issue names were already fully encoded in `transitions.json` across Phases 1–3 (#568/#569/#574/#575) — held for the entire cycle: zero diff landed in the FSM table or engine. The cycle's only genuinely new logic is the multi-issue iteration + finalizer-invocation glue (`scan.go`) and the renderer wiring to run it mechanically. One independent, pre-existing bug (`cn cell finalize`'s hub-requirement) was discovered and fixed along the way — necessary for AC4 to hold at all in this repository, not scope creep: without it, the scanner's finalize-invocation path (and the already-merged #591 finalizer step) was dead code here.

## Process-gap audit (γ→α→β pipeline)

### Claim re-entry handled correctly

A prior `cds-dispatch` firing had already created `cycle/593` and written `CLAIM-REQUEST.yml`, but the FSM claim transition never applied (issue stayed `status:todo`; no implementation artifacts existed). Classified `run_class: first_pass` (no rejection history, no prior closeouts) per the repair re-entry preflight's own detection criteria — correctly distinguished from a `status:changes` repair re-entry, which this was not. The stale branch was rebased onto current `main` (only 1 commit behind) rather than discarded, reusing the existing claim marker.

### Scope discipline — held

No new status labels, no daemon, no claim/resume-guard changes, no FSM guard weakening, zero diff to `transitions.json`/`table.go`/`decision.go`/`CellResumeCmd`/`CellReturnCmd`. The one infrastructure change (`CellFinalizeCmd`'s `NeedsHub`) is the minimum fix that makes this cycle's own AC4 achievable, verified to have zero effect on any other command.

### A discovered gap surfaced honestly, not worked around

The hub-detection bug could have been silently worked around (e.g., skip live-testing the finalize-invocation path, ship only fixture-covered logic). Instead it was reproduced in a disposable sandbox, root-caused to `discoverHub()`'s `.cn`-only check, fixed with the smallest correct change, and independently re-verified by β in a *separate* sandbox before being accepted. This is the standard this cycle held itself to for its one non-obvious finding.

## Deliverable evidence

```yaml
deliverable_evidence:
  pr: "#595 (cycle/593 -> main)"
  head_sha: "2ac5549bef69d4cfc2438eb4629f1864b2bccd11"
  base_sha: "2b7350664432eabdb697a387f68f7bc0111f9907"
  commits_beyond_base: 5
  closeout_artifacts: [gamma-scaffold.md, self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md]
```

PR #595 was opened by `cn cell finalize --issue 593` — the dogfood proof that this cycle's own hub-gate fix works: the exact same code path this cycle adds to the `cds-dispatch` renderer's post-run finalizer step created this PR.

## Close-out triage table

| Item | Disposition |
|---|---|
| Scanner logic (`scan.go`) | Shipped, tested (9 new tests), converged |
| Hub-gate fix (`cmd_cell.go`) | Shipped, tested (existing `cell` package suite green), independently reproduced by β |
| Renderer wiring | Shipped, goldens regenerated, leak-audit clean, idempotent |
| Fixtures | 2 new (`scan-died-after-pr-before-review-request.json`, `scan-died-after-review-request-before-review.json`), 4 reused from the pre-existing #568/#569/#574 fixture set |
| Deferred work | None identified within this cycle's scope. Sub D (guard consolidation / remaining scaffold cleanup, per #583) is explicitly not started here, per the issue's own "After merge" instruction — not to be dispatched without operator go. |

## Scope confirmation

Re-read against every named STOP condition in the issue: no new GitHub permissions needed, no new status labels, claim semantics unchanged beyond resume/checkpoint awareness (and even that was not needed — see "Scope discipline" above), scanner never duplicates a PR or comment (AC9, independently verified), scanner inspects only data available via the existing GitHub REST API + local git (same primitives `evaluate`/`cell finalize` already use), FSM guards unchanged, worker behavior unchanged outside recovery/resume, Demo 0 not touched. None triggered.

## Closing statement

Ready for `status:review`. #583 should be updated to mark Sub C complete once this PR merges; Sub D remains unstarted pending operator authorization, per the issue's own instruction.
