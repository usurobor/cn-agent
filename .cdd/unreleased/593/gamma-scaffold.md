# γ R0 scaffold — cnos#593

## Issue reference

- **Issue:** [usurobor/cnos#593](https://github.com/usurobor/cnos/issues/593) — "cds/runtime: mechanically recover stranded dispatch cells (Sub C of #583)"
- **Mode:** design-and-build
- **Parent:** #583 (master wave — make the dispatch cell fully mechanical). **Sub C.**
- **Preconditions (all merged):** #584 (Sub A — doctrine), #575 (FSM owns all four lifecycle transitions), #591 (Sub B — mechanical checkpoint + PR-open finalizer, `cn cell finalize`).
- **Protocol:** cds
- **run_class:** `first_pass`. A prior `cds-dispatch` firing had written `.cdd/unreleased/593/CLAIM-REQUEST.yml` and created `cycle/593`, but the FSM claim transition never applied (issue remained `status:todo`; no CDD implementation artifact existed on the branch — only the claim marker itself). No `status:changes` history, no beta/delta findings, no closeouts. Re-claimed cleanly; see the issue-593 claim comment.
- **Cycle branch:** `cycle/593`, rebased onto `main`@`2b735066` (1 commit ahead: the pre-existing `CLAIM-REQUEST.yml`).

## Governing rule (restated from the issue)

No κ/operator liveness guessing. Facts observed (live run? cycle branch? branch commits? draft/open PR? `REVIEW-REQUEST.yml`? CDD artifacts? repair context? current status labels?) map to exactly one reconcile action: **no-op · finalize/checkpoint · requeue · move to review via FSM · block with evidence.**

## Key finding: the six behavior cases are ALREADY implemented by the existing FSM table

Reading `src/packages/cnos.cds/skills/cds/fsm/transitions.json`'s `in-progress` and `review` rule sets (landed across Phase 1/#568, Phase 2/#569, cnos#574, Phase 3/#575) shows every one of the issue's six state cases already has a matching rule:

| Issue case | Existing rule (transitions.json, state `in-progress` unless noted) |
|---|---|
| 1. live run active → no-op | `all_true: [run_active]` → `outcome: valid, action: none` |
| 2. dead, no matter → requeue | catch-all rule → `outcome: proposed, action: propose_status_todo, target_state: todo` |
| 3. dead, matter, no PR → finalize + PR | `all_false: [run_active], any_true: [branch_has_commits, pr_exists, pr_has_commits]` → `outcome: proposed, action: propose_delta_recovery, target_state: ""` |
| 4. dead, PR exists, no REVIEW-REQUEST.yml → keep as checkpointed matter | same `propose_delta_recovery` rule (PR-exists alone still satisfies `any_true`) — target_state stays `""`, so no blind status move; matter is preserved for the next worker |
| 5. PR + REVIEW-REQUEST.yml → review via FSM | `all_true: [review_request_present, pr_exists, pr_has_commits]` → `outcome: proposed, action: propose_status_review, target_state: review` |
| 6. stale review, missing proof → blocked | state `review`'s catch-all → `outcome: blocked`, `missing_evidence` populated |

**Consequence for scope:** `transitions.json` and the FSM engine (`table.go`/`decision.go`) need **zero changes**. Sub C is the **iteration + finalizer-invocation glue** around the existing per-issue `evaluate` logic, not a new rule set. This matches the issue's own "Implementation preference": *"Add a mechanical command/workflow step ... that scans active CDS dispatch issues and reconciles only those with stale/inconsistent state."*

## Design

New sub-verb: **`cn issues fsm scan --protocol P [--apply]`**, in the same package as `evaluate` (`src/packages/cnos.issues/commands/issues-fsm/`, `scan.go`):

1. **List** — GitHub REST `issues?labels=dispatch:cell,protocol:{P}&state=open`, client-filtered to `status:in-progress` / `status:review` (the only two states this scanner reconciles — `todo`/`ready`/`changes`/`blocked` await a claim event or an operator, they don't "strand").
2. **Per issue** — `assembleLive` (the exact function `evaluate` uses) → `Evaluate(table, snap)` (same engine, same table).
3. **Reconcile per decision**, using ONLY primitives that already exist:
   - `outcome: valid` → no-op (case 1).
   - `outcome: blocked` → report only; **never comment** (see idempotence note below) (case 6, and any other blocked state).
   - `outcome: proposed, action: propose_delta_recovery` → shell out to `cn cell finalize --issue N` (cnos#591's finalizer, unchanged) — checkpoints matter and idempotently ensures a draft PR; **no label write** (cases 3/4).
   - `outcome: proposed, target_state != ""` → `--apply` writes the label via `applyStatusLabel` (the exact function `evaluate --apply` calls) (cases 2/5).
4. **Comment discipline (AC9 idempotence):** only post an issue comment on an actual state-changing event — a label reconciliation, or a finalize call that just opened a **new** draft PR (detected via the finalizer's own stdout marker `"opened draft PR for"`, distinguishing it from `cn cell finalize`'s pre-existing idempotent-no-op wording). A blocked or already-checkpointed state produces no repeat comment on every ~15-minute sweep tick.

### Why the finalizer is invoked as a subprocess, not a Go import

`cn cell finalize`'s implementation (`cell.Finalizer`) lives in `src/go/internal/cell`, part of the `src/go` Go module. `issues-fsm` is a **separate module** (its own `go.mod`, per `go.work`). Go's `internal/` import-visibility rule scopes `.../src/go/internal/cell` to importers rooted under `.../src/go/...` — `issues-fsm` cannot import it directly. `scan.go`'s `RunFinalize` therefore self-re-execs the currently-running `cn` binary (`os.Executable()`) as `cn cell finalize --issue N`, reusing cnos#591's finalizer verbatim (same binary, same code path a human or the workflow's own finalizer step would run) rather than duplicating its checkpoint/PR logic in a second implementation.

## Discovered blocker + minimal fix: `cn cell finalize` required a `.cn/` hub that does not exist in this repo

Testing `cn cell finalize` directly (to verify AC4's finalizer-invocation path) surfaced a **pre-existing, independent bug**: `CellFinalizeCmd.Spec()` declared `NeedsHub: true`, and `main.go`'s `discoverHub()` gate only recognizes a `.cn/` **package-vendor hub directory** — which this repository (the `cnos` package source tree itself, not a downstream product hub) does not have. Any invocation of `cn cell finalize` against this exact checkout — including the **already-merged** `.github/workflows/cnos-cds-dispatch.yml`'s "Mechanical checkpoint + PR finalizer" step — hits `✗ Command 'cell-finalize' requires a hub, but no hub found.` before `Finalize` ever runs, unconditionally. `cell.Finalizer.RepoRoot` only ever needs a plain git-repository root (to locate `.cdd/unreleased/{N}/` and run `git`/`gh`), not `.cn/vendor/packages/` discovery — `CellReturnCmd` already recognized this distinction (`NeedsHub: false`, "operates via gh CLI; no hub required").

**Minimal fix (src/go/internal/cli/cmd_cell.go):** `CellFinalizeCmd.Spec()` now declares `NeedsHub: false`; `Run()` resolves the repo root itself via a new `gitRepoRoot(ctx)` helper (`git rev-parse --show-toplevel`) when `inv.HubPath` is empty, falling back to the discovered hub path when one exists (so behavior in an actual `.cn/`-carrying downstream hub is unchanged). Verified in a sandboxed synthetic git repo (no `.cn/`): the command now reaches real matter-detection/checkpoint logic instead of exiting at the hub gate. `CellResumeCmd`/`CellReturnCmd` are untouched — this fix is scoped to the one command Sub C's scanner depends on.

This is necessary, not optional, for AC4/AC10 to hold at all in this repository — the recovery scanner's finalize-invocation path (and arguably #591's own already-merged step) was otherwise dead code in this exact environment.

## Surfaces touched

| Surface | Path | What changed |
|---|---|---|
| New scanner engine + CLI wiring | `src/packages/cnos.issues/commands/issues-fsm/scan.go` (new), `issuesfsm.go` (add `scan` sub-verb) | `RunScan`, `ScanOptions` (injectable `ListActiveIssues`/`AssembleFacts`/`RunFinalize`/`PostComment`), live defaults, `runScan` CLI wrapper |
| Hub-gate fix | `src/go/internal/cli/cmd_cell.go` | `CellFinalizeCmd.Spec().NeedsHub` → `false`; `gitRepoRoot` fallback helper |
| Fixtures + tests | `src/packages/cnos.issues/commands/issues-fsm/testdata/scan-*.json`, `scan_test.go` (new) | AC10's five strand fixtures + AC2 (live-run no-op) + AC9 (idempotence) |
| Renderer (mechanical wiring) | `src/packages/cnos.core/commands/install-wake/cn-install-wake` | New `if: always()` sibling step to the existing finalizer step, gated on `role == "dispatch"`, running `cn issues fsm scan --protocol "$protocol" --apply` (protocol read from the manifest via `jq`, not hardcoded — AC7/AC8 leak-audit constraint) |
| Regenerated goldens | `cnos-cds-dispatch.golden.yml`, `.github/workflows/cnos-cds-dispatch.yml` | Re-rendered via `cn-install-wake cds-dispatch [--out ...]` |

## Per-AC oracle list (abbreviated — full detail in self-coherence.md)

- **AC1** — `grep -n "scan" src/packages/cnos.issues/commands/issues-fsm/issuesfsm.go`; `cn issues fsm scan --help` exits 0.
- **AC2** — fixture with `run_active: true` at `in-progress` → `RunScan` reports `outcome: valid`, no finalize/apply/comment call.
- **AC3** — fixture with no matter at `in-progress` → `--apply` applies `todo`; a second identical run is a no-op (already `todo`, out of scan's active-state set).
- **AC4** — fixture with `branch_has_commits: true`, `pr_exists: false` → injected `RunFinalize` invoked; injected fake returns "opened draft PR for" → `Commented: true`.
- **AC5** — fixture with `pr_exists: true` already (no REVIEW-REQUEST.yml) → `RunFinalize` invoked (idempotent per #591), injected fake returns the "idempotent no-op" wording → `Commented: false` (no duplicate).
- **AC6** — fixture with `review_request_present`, `pr_exists`, `pr_has_commits` all true at `in-progress` → `--apply` applies `review`.
- **AC7** — fixture at `review` with missing evidence → `outcome: blocked`, no mutation, no comment.
- **AC8** — every mutation path in `scan.go` routes through `applyStatusLabel` (label) or the finalizer subprocess (matter/PR only, no label code path) — grep confirms zero direct `gh issue edit --add-label/--remove-label` calls in `scan.go`.
- **AC9** — `TestScan_Idempotent`: run `RunScan` twice against unchanged injected facts; assert the second run's `Commented`/`Reconciled`/`Finalized` counts are zero (or identical no-op notes), and the injected `PostComment`/`RunFinalize` call counts do not grow on the second pass for an unchanged blocked/checkpointed state.
- **AC10** — fixtures named per the issue's own strand-class list; `go test ./...`, `install-wake-golden` re-render+diff, `check-dispatch-repair-preflight.sh`, `check-dispatch-closeout-integrity.sh`, Go/Package/Binary CI all green.

## Scope guardrails honored

- No new status labels — scan only ever proposes/applies transitions the existing table already encodes.
- No daemon; scan runs as one `cn issues fsm scan` invocation per mechanical workflow step.
- No claim/resume-guard changes: the delta-recovery path never returns a checkpointed cell to `status:todo` (target_state stays `""`), so the "todo issue with an existing draft PR" resume-vs-conflict question this issue flagged as conditionally in-scope never arises — matter always keeps the cell at `status:in-progress` until an operator/δ adds `REVIEW-REQUEST.yml` or otherwise intervenes.
- No FSM guard weakening: `applyStatusLabel` is called with exactly the `(CurrentState, TargetState)` pair `Evaluate` already proposed; scan does not invent transitions.
- Zero diff to `transitions.json` / `table.go` / `decision.go`.
