# α close-out — cnos#593

## Cycle summary

Built `cn issues fsm scan` (the mechanical recovery scanner named by the issue), reusing the existing `evaluate`/`transitions.json` decision engine and the existing `cn cell finalize` checkpoint primitive without modifying either. Discovered and fixed an independent, pre-existing bug (`cn cell finalize` required a `.cn/` hub directory that doesn't exist in this repo) that would otherwise have made AC4 impossible to satisfy in this environment. Wired the scanner into the `cds-dispatch` renderer as a pre-claim step and regenerated both goldens + the live workflow.

## What was built

- `src/packages/cnos.issues/commands/issues-fsm/scan.go` — `RunScan`/`ScanOptions` (injectable dependencies mirroring `cell.Finalizer`'s testing pattern) + live defaults (`liveListActiveIssues`, `liveRunFinalize`, `livePostComment`) + `renderScan`.
- `issuesfsm.go` — `scan` sub-verb wired into `Run()`; `runScan` CLI wrapper (flag parsing mirrors `runEvaluate`).
- `scan_test.go` + two new fixtures — nine tests covering all six behavior cases, dry-run, and idempotence.
- `src/go/internal/cli/cmd_cell.go` — `CellFinalizeCmd.Spec().NeedsHub` → `false`; `gitRepoRoot` fallback.
- `cn-install-wake` — new "Mechanical recovery scanner" step, gated on `role == "dispatch"`, inserted before the existing write-fence/finalizer steps.
- Regenerated `cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml`.

## What was hard

Realizing the hub-detection bug existed at all. My first live (read-only) smoke test of the scanner against the real `usurobor/cnos` repo returned real stranded issues (#503, #532, plus #593 itself) — exactly what AC1/AC4 are supposed to find. But invoking `cn cell finalize` against one of those (in a disposable sandbox, never against the real issues #503/#532 themselves, to avoid taking any action outside my claimed cell) failed at a gate I hadn't written and hadn't expected: "Command 'cell-finalize' requires a hub, but no hub found." Tracing that back to `discoverHub()`'s `.cn`-directory-only check, and confirming this repo genuinely has no `.cn/` anywhere in its history, took longer than writing the scanner itself.

## What I'd do differently

I would have run a live dry-run smoke test *before* writing any fixture tests, not after — the sandboxed fixture tests all passed using injected fakes (which never touch the CLI's hub gate at all), so I would have shipped a scanner whose real-world finalize-invocation path was silently broken had I not separately smoke-tested against the real repo.

## Friction-note resolutions — confirmed holding up

The gamma-scaffold's "Key finding" (the six behavior cases already exist in `transitions.json`) held throughout: zero diff was needed to that file, `table.go`, or `decision.go`. This was verified, not assumed — every fixture test asserts against the real, unmodified table (`loadRealTable(t)`), not a synthetic one.

## Scope confirmation

No new status labels. No daemon. No claim/resume-guard changes (delta-recovery's `target_state` stays empty; a cell with a checkpointed-but-unreviewed PR never returns to `status:todo` through this scanner, so the "resume vs. conflict" question the issue flagged as conditionally in-scope never had to be answered). No FSM guard weakening. Zero diff to `transitions.json`/FSM Go files/`CellResumeCmd`/`CellReturnCmd`.
