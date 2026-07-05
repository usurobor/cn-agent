# self-coherence — cnos#593

## Gap

Stranded CDS dispatch cells (`status:in-progress`/`status:review` left behind by a dead workflow run) had no mechanical reconciliation path. `cn issues fsm evaluate` (Phase 1–3, #568/#569/#574/#575) already encodes the correct per-issue decision for every strand shape, and `cn cell finalize` (#591) already checkpoints matter into a draft PR — but nothing iterated across every active issue and invoked these two primitives together on a schedule. Operators (or a human noticing a stale board entry) were the only reconciliation path.

## Skills

Tier 3: `cdd`, `cdd/issue`, `cdd/issue/proof`, `eng`, `eng/ship`, `eng/evolve`, `eng/process-economics`, `eng/test`, `go`, `write`.

## AC Coverage

| AC | Status | Evidence |
|---|---|---|
| AC1 | ✅ | `cn issues fsm scan` exists (`scan.go`, `issuesfsm.go`); `--help` documents behavior |
| AC2 | ✅ | `TestScan_LiveRunActive_NoOp` |
| AC3 | ✅ | `TestScan_DeadNoMatter_RequeuesToTodo` |
| AC4 | ✅ | `TestScan_DeadWithMatterNoPR_FinalizesAndComments` |
| AC5 | ✅ | `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` (draft PR recognized as checkpointed matter; finalize's own idempotent no-op path taken; no duplicate) |
| AC6 | ✅ | `TestScan_ReviewRequestWithMatter_AppliesReview` |
| AC7 | ✅ | `TestScan_StaleReviewNoProof_BlockedNoComment` |
| AC8 | ✅ | Every mutation in `scan.go` routes through `applyStatusLabel` (label) or the `cn cell finalize` subprocess (matter/PR only); `grep -n "add-label\|remove-label\|issue edit" scan.go` → zero matches |
| AC9 | ✅ | `TestScan_Idempotent` (label transitions never repeat once the issue leaves the active-scan-state set; finalize is re-invoked but its own idempotent no-op suppresses the duplicate comment) |
| AC10 | ✅ | Fixtures named per the issue's own strand-class list (see "Fixture naming" below); full gate run in "Full verification run" |
| AC11 | ✅ (see "Full verification run") | Go/I1/I6/install-wake-golden/dispatch-repair-preflight/dispatch-closeout-integrity all green locally; I2 spot-checked (untouched surface); I4/I5 tools unavailable in this sandbox, scope of change does not touch either surface (no SKILL.md frontmatter edits, no new external/broken markdown links) |

## CDD Trace

- Parent: #583 (master wave). Preconditions: #584 (merged), #575 (merged), #591 (merged) — all verified present on `main` before this cycle began.
- This cycle adds **zero** diff to `src/packages/cnos.cds/skills/cds/fsm/transitions.json`, `table.go`, or `decision.go` — see gamma-scaffold.md's "Key finding" section for why the existing table already covers all six behavior cases.

## Self-check

Read back against the issue's own scope guardrails: no daemon, no new status labels, no claim/resume-guard changes (delta-recovery never returns a cell to `status:todo`), no FSM guard weakening (scan calls `applyStatusLabel` with exactly the `(CurrentState, TargetState)` pair `Evaluate` proposed), zero diff to `transitions.json`/FSM Go files. All held.

## §R0

### Summary of what was built

1. **`cn issues fsm scan --protocol P [--apply]`** (`src/packages/cnos.issues/commands/issues-fsm/scan.go`, wired into `issuesfsm.go`'s `Run()` dispatcher as a second sub-verb alongside `evaluate`). Lists open `dispatch:cell`+`protocol:{P}` issues currently at `status:in-progress`/`status:review`, observes each via the *same* `assembleLive` function `evaluate` uses, evaluates against the *same* transition table, and reconciles:
   - `outcome: valid` → no-op.
   - `outcome: blocked` → reported only, never mutated, never commented (idempotence: a blocked state does not accumulate a comment every ~15-minute sweep tick).
   - `outcome: proposed, action: propose_delta_recovery` → invokes `cn cell finalize --issue N` (self-re-exec of the running `cn` binary — `issues-fsm` is a separate Go module from `src/go/internal/cell` and cannot import it directly past Go's `internal/` visibility rule) to checkpoint matter into a draft PR; comments only when the finalizer's own stdout shows a **new** PR was opened (`"opened draft PR for"`), never on its idempotent no-op wording.
   - `outcome: proposed, target_state != ""` → applies the label via `applyStatusLabel`, the exact function `evaluate --apply` calls, then comments once.
2. **A pre-existing, independent bug fix**: `cn cell finalize` required a `.cn/` "hub" directory (`CellFinalizeCmd.Spec().NeedsHub: true`) that does not exist anywhere in this repository — meaning the command (and, by the same defect, the *already-merged* #591 finalizer workflow step) exited `"no hub found"` before ever reaching its own logic, unconditionally, in exactly the environment it is meant to run in. Verified in a sandboxed synthetic git repo with no `.cn/`: pre-fix, `cn cell finalize --issue N` failed at the hub gate; post-fix, it reached real matter-detection/checkpoint logic (failing only on the synthetic repo's missing `origin` remote, as expected). Fix: `NeedsHub: false` (mirroring `CellReturnCmd`'s existing precedent) + a `gitRepoRoot(ctx)` fallback (`git rev-parse --show-toplevel`) in `CellFinalizeCmd.Run()` when `inv.HubPath` is empty. `CellResumeCmd`/`CellReturnCmd` untouched.
3. **Renderer wiring**: a new `if`-gated (`role == "dispatch"`) "Mechanical recovery scanner" step in `cn-install-wake`, inserted right after `actions/checkout@v4` (a pre-claim pass, distinct from the post-run finalizer step) — builds `cn` fresh and runs `cn issues fsm scan --protocol "$scan_protocol" --apply`, with `$scan_protocol` read from the manifest's `.protocol` field via `jq` **at render time** (baked into the emitted YAML as a fixed value, never a literal token in the renderer's own source — the AC7/AC8 leak-audit constraint). Both goldens + the live `.github/workflows/cnos-cds-dispatch.yml` regenerated; sha256 of live == golden; second render is byte-identical (idempotent).

### AC-by-AC verification

See "AC Coverage" table above; test names are exact, all runnable via `cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -run TestScan -v`.

### Golden regeneration

```
./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin   # unchanged
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch # rendered (this cycle's step added)
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml
sha256sum .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
# identical hashes; re-running cds-dispatch a second time reports "(unchanged)"
```

### Leak-audit (AC7/AC8 of install-wake-golden.yml; HARD gate)

```
grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' src/packages/cnos.core/commands/install-wake/cn-install-wake   # => 0
grep -nE 'admin_only|disallowed_surfaces|defer_path|cell_execution' src/packages/cnos.core/commands/install-wake/cn-install-wake | grep -vE '^[0-9]+:[[:space:]]*#' | grep -vE '(die|manifest_error|log_writer_mis_declaration)[[:space:]]'   # => empty
```
Both clean. `$scan_protocol` is a jq-read shell variable, never a literal "cds"/"cdr"/"cdw" substring in the renderer source.

### Design choices where the scaffold left flexibility

- **Comment idempotence via stdout substring match** rather than a GitHub comment-listing round-trip: `cn cell finalize`'s own stdout already distinguishes "opened draft PR for" (new) from "PR already open for ... idempotent no-op" (unchanged). Parsing that avoids an extra GitHub API call per scan tick per issue and keeps the idempotence guarantee purely local/testable (no network mock needed for the comment-dedup logic itself).
- **Scan step placement**: inserted as a pre-claim step (right after checkout), not as a sibling to the post-run finalizer — the issue's own text ("runs mechanically from the dispatch workflow schedule or pre-claim phase") and the semantic difference (reconciling *other* stranded cells vs. checkpointing *this* run's own matter) both point at "before the cognitive work phase," not "after."
- **No re-observation after invoking finalize**: once `propose_delta_recovery` fires and finalize runs, the decision does not change (PR-exists alone still satisfies the rule's `any_true` guard, and `review_request_present` — the only thing that would unlock a *different* rule — is never set by this scanner). Re-evaluating post-finalize would be a no-op computation; skipped for simplicity, documented here so a future reader doesn't wonder why it's absent.

### Dogfood smoke (AC1/AC4, real end-to-end run against this repo)

A live (read-only, no `--apply`) run against `usurobor/cnos` during this cycle's own development correctly identified three real stranded issues in production (#503 at `status:review` with missing evidence → `blocked`; #532 and #593 itself at `status:in-progress` with matter and no PR → `propose_delta_recovery`) — surfacing the exact hub-detection bug fixed in this cycle (finalize failed with "no hub found" against these real issues before the fix). This is the dogfood evidence the hub-gate fix was necessary, not speculative.

### Friction encountered

1. Cross-module import restriction (`issues-fsm` cannot import `src/go/internal/cell` — separate Go modules, `internal/` visibility scoped to the `src/go` module tree) — resolved via self-re-exec subprocess invocation (`os.Executable()` + `cell finalize --issue N`), reusing #591's finalizer verbatim rather than duplicating its logic.
2. The `cn cell finalize` hub-requirement bug (see "Summary" #2 above) — discovered only by attempting to actually exercise AC4 against a real issue; not something a fixture-only test suite would have caught, since `Finalizer`'s unit tests never go through the CLI's `NeedsHub` gate at all.

### Full verification run (AC10/AC11, checked here at R0 as well as at β/converge)

```
cd src/packages/cnos.issues/commands/issues-fsm && go test ./...                     # ok
cd src/go && go build ./... && go build -o cn ./cmd/cn && go test ./... && go vet ./... && ./cn help   # ok
./cn build --check                                                                     # I1 ok
./cn cdd verify --unreleased --exceptions .cdd/exceptions.yml                          # I6: 0 failed
src/packages/cnos.cdd/commands/cdd-verify/test-fixtures.sh                            # ok
bash scripts/ci/check-dispatch-repair-preflight.sh                                    # ok
bash scripts/ci/check-dispatch-closeout-integrity.sh                                  # ok
diff docs/reference/schemas/protocol-contract.json tests/fixtures/protocol-contract.json  # I2 ok
scripts/kata/run-all.sh (Tier 1)                                                       # ok
```
I4 (lychee) and I5 (CUE frontmatter) tools are unavailable in this sandbox; this cycle does not touch any SKILL.md frontmatter or introduce new markdown links, so no regression is expected — flagged for β to independently confirm scope, not a self-graded pass.

### Review-ready signal

α believes this cycle is ready for β review: every AC has a named, runnable test; the two structural gates this cycle's changes could plausibly break (install-wake-golden leak-audit, I6 CDD ledger) were run and are green; the one non-obvious discovery (the hub-detection bug) is documented with before/after evidence, not asserted.
