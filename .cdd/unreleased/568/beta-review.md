# β independent review — cnos#568 (cds/issues: FSM Phase 1 — read-only issue-state reconciler)

**Reviewer:** β · **Cycle:** cycle/568 @ `7096250c` · reviewed against `origin/main` @ `cd61f3d2` (merge-base) and the γ scaffold's operationalized AC→oracle→surface table.

## §R0

**verdict: converge**

I independently re-derived every AC by reading the actual source (not α's prose), running `go test ./... -v` myself, building the `cn` binary and driving the CLI against fixtures myself, and grepping/diffing the protected files myself. I found no AC violations and no scope creep. One non-blocking observation is noted under AC5/AC6 below for α/δ's awareness but it does not gate convergence — it is a latent-precision note, not a violation of any AC oracle as written.

---

### AC1 — FSM transition table exists (single declarative source)

**Independently verified.** Read `src/packages/cnos.cds/skills/cds/fsm/transitions.json` directly: it is pure JSON (`states`, `guards`, `transitions[].rules[]` with `all_true`/`any_true`/`all_false`/`outcome`/`action`/`target_state`/`reason`/`evidence_guards`) — no Go literals. Read the engine (`src/packages/cnos.issues/commands/issues-fsm/table.go:157-210`, `Evaluate`): it loops `t.Transitions`, matches `tr.State == state` (a runtime string from `FactSnapshot.Labels`, `snapshot.go:73-80`), then matches `Rule`s in file order. There is no `switch`/`if`-chain anywhere in `table.go` or `decision.go` naming a CDS state (`"review"`, `"changes"`, `"in-progress"`, ...). Confirmed by my own grep:

```
$ grep -n '"review"\|"changes"\|"in-progress"' src/packages/cnos.issues/commands/issues-fsm/*.go | grep -v _test.go
(no output)
```

`guardFuncs` (`table.go:90-100`) is a fixed vocabulary of generic evidence predicates (`run_active`, `branch_has_commits`, `pr_exists`, ...) over `FactSnapshot` fields — not a CDS-state switch. Package split matches the scaffold: `cnos.issues` owns the engine, `cnos.cds` owns the data file. `go test ./... -v` in the new module: `TestAC1_TableLoadsAsData` PASS (ran myself, see below).

### AC2 — `cn issues fsm evaluate --issue {N}` exists and explains

**Independently verified.** Built the binary myself (`cd src/go && go build -o cn ./cmd/cn`, exit 0) and ran it directly against three different fixtures. Sample (issue #604, `review-empty.json`):

```
$ ./cn issues fsm evaluate --issue 604 --fixture .../testdata/review-empty.json
Current state: review
Observed facts: labels, run_state, run_conclusion, branch_exists, commits_beyond_base,
  pr_exists, pr_commit_count, cdd_artifacts, review_request_present,
  repair_contract_present, checks_state  [all present]
Decision:
  outcome: blocked
  enabled_transition: (none)
  blocked_reason: empty review: status:review with no PR, no commits beyond base, and no REVIEW-REQUEST.yml present.
  missing_evidence: pr_exists, branch_has_commits, review_request_present
  proposed_action: block
(read-only: no label was written; this is a proposal, not an action)
```

`./cn issues fsm apply --issue 1` → rejected, exit 1, stderr explains apply is Phase 2 (ran myself). `./cn help | grep issues` lists both `issues-map` and `issues-fsm`. `cmd_issues_fsm.go` (`src/go/internal/cli/cmd_issues_fsm.go`) is a one-line delegation (`return issuesfsm.Run(...)`), matching the AC2 CLI surface.

### AC3 — empty-review is detected

**Independently verified.** `TestAC3_EmptyReviewBlocked` and `TestAC3_ReviewWithEvidenceIsValid` both PASS under my own `go test -v` run. I also drove both fixtures through the built CLI directly (`review-empty.json` → `blocked` with the 3-item missing-evidence list shown above; `review-with-pr.json` → `outcome: valid`, confirmed by reading `testdata/review-with-pr.json`, which sets `pr_exists: true`). The guard logic (`transitions.json` `review` state, rule 2: `all_false: [pr_exists, branch_has_commits, review_request_present]`) is what fires — read directly, not inferred from the test name.

### AC4 — dead in-progress / no matter → proposes requeue

**Independently verified — ran the CLI myself.**

```
$ ./cn issues fsm evaluate --issue 602 --fixture .../testdata/in-progress-dead-no-matter.json
Decision:
  outcome: proposed
  enabled_transition: in-progress -> todo
  proposed_action: propose_status_todo
```

Fixture (`in-progress-dead-no-matter.json`, read directly): `run_state: "completed"`, `branch_exists: false`, `pr_exists: false`. This is the table's third (fallback) rule for `in-progress` — no `all_true`/`all_false`/`any_true` guard, i.e. "if nothing else matched" — which is correct only because I confirmed (next AC) that the with-commits case is caught by an earlier rule, so the fallback genuinely means "no matter."

### AC5 — dead in-progress / branch-with-commits → proposes δ-recovery (the cnos#368 guard)

**Independently verified — read the guard logic directly, not just the test.** `transitions.json`'s `in-progress` state has 3 rules in this file order:

1. `all_true: [run_active]` → valid/none
2. `all_false: [run_active]`, `any_true: [branch_has_commits, pr_exists, pr_has_commits]` → `proposed` / `propose_delta_recovery` / `target_state: ""`
3. (fallback, no conditions) → `proposed` / `propose_status_todo` / `target_state: "todo"`

`Evaluate` (`table.go:169`) iterates rules **in file order** and returns on **first match** (`table.go:174-198`) — so rule 2 (delta-recovery) is checked and matched *before* the fallback rule 3 (todo) is ever reached, for any dead run with matter. I ran the actual CLI against `in-progress-dead-with-commits.json` (`branch_exists: true`, `commits_beyond_base: 4`, `pr_exists: false`) myself:

```
$ ./cn issues fsm evaluate --issue 603 --fixture .../testdata/in-progress-dead-with-commits.json
Decision:
  outcome: proposed
  enabled_transition: (none)
  proposed_action: propose_delta_recovery
  reason: dead run with matter: ... Never re-dispatch blind over existing work (cnos#368 failure mode) -- propose delta-recovery ... instead.
```

`target_state` is empty (no direct status flip) and `enabled_transition` is `(none)` — it genuinely never proposes bare `status:todo` for this case. `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery` additionally asserts `TargetState != "todo"` directly as a regression guard; I ran it and it PASSes. This is the strongest AC and I gave it the most scrutiny — confirmed by rule ordering in the data file, by the pure-function engine code, and by live CLI output, not merely by trusting the test's name.

**Non-blocking observation (not an AC violation):** `RunConclusion` is captured in `FactSnapshot` and printed in the facts block, but no guard function in `table.go`'s `guardFuncs` ever reads it — the `in-progress` rules key off `run_active` (queued/in_progress) only, treating any non-active run (success or failure conclusion) with matter identically as "dead run with matter → delta-recovery." AC5's oracle as literally written ("no active run + commits beyond base + no PR → δ-recovery") does not require a success/failure distinction, so this is not a violation. I considered whether this could misfire on the *current* multi-agent cycle (α→β→δ running inside one dispatch job), but confirmed via `.github/workflows/cnos-cds-dispatch.yml` that one wake = one workflow run covering the whole cell (not one run per role-turn), so `run_active` would correctly still read `in_progress` for the duration of an in-flight cycle like this one. Flagging only so α/δ are aware `checks_passing` and `RunConclusion`/`cdd_artifacts_present` are collected but currently unused by any rule — possibly intentional headroom for Phase 2, not a Phase-1 defect.

### AC6 — `changes → todo` without repair context is blocked

**Independently verified.** Read `testdata/changes-no-repair.json` (`repair_contract_present: false`) and `changes-with-repair.json` (`repair_contract_present: true`) directly. `transitions.json`'s `changes` state: rule 1 `all_true: [repair_contract_present]` → `proposed`/`propose_repair_dispatch`/`target_state: "todo"`/`repair_pass: true`; rule 2 `all_false: [repair_contract_present]` → `blocked`/`target_state: ""`. Ran `go test -run TestAC6` myself: both `TestAC6_ChangesWithoutRepairContextBlocked` and `TestAC6_ChangesWithRepairContextEnablesRepairPass` PASS.

### AC7 — idempotence

**Independently verified, two ways.** (1) `TestAC7_Idempotent` PASS under my own run, across all 5 non-trivial fixtures, asserting both the pure `Evaluate()` struct fields and the byte-identical rendered `Run()` CLI output across two calls. (2) I additionally ran the built CLI twice myself outside of the test suite and diffed the raw stdout:

```
$ ./cn issues fsm evaluate --issue 603 --fixture .../in-progress-dead-with-commits.json > run1.txt
$ ./cn issues fsm evaluate --issue 603 --fixture .../in-progress-dead-with-commits.json > run2.txt
$ diff run1.txt run2.txt && echo IDENTICAL
IDENTICAL
```

`Evaluate` and `Decision.Render` are pure functions of `(Table, FactSnapshot)` with no global state, confirmed by reading both functions end to end.

### AC8 — no mutation, no authority flip (Phase-1 boundary)

**Independently verified, three ways, all run myself (not just re-read from self-coherence.md).**

1. `./cn issues fsm apply --issue 1` → exit 1, "unknown subcommand" + explains Phase 2/cnos#569.
2. Grepped the entire new source tree myself for label-write signatures:
   ```
   $ grep -rniE 'gh issue edit|gh label|-X PATCH|-X POST|/labels|AddLabel|RemoveLabel|--apply' \
       src/packages/cnos.issues/commands/issues-fsm/*.go src/go/internal/cli/cmd_issues_fsm.go
   ```
   Every hit is inside `issuesfsm_test.go` (string literals used by `TestAC8_NoLabelMutationCodeInSource` to check *for* their absence) or comment prose in `issuesfsm.go` documenting the absence — none in executable non-test code. Separately grepped `fetch.go` for `http.Method`: only `http.MethodGet` appears, no `Post`/`Patch`/`Put`/`Delete`.
3. Diffed the three protected files against `origin/main` myself:
   ```
   $ git diff origin/main -- .github/workflows/cnos-cds-dispatch.yml \
       src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md \
       src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md | wc -l
   0
   ```
   Zero diff, confirmed directly (not copy-pasted from α's claim).

### AC9 — CI guard for evaluator fixtures

**Independently verified.** Read `.github/workflows/build.yml` directly (lines 45-47): a new step `Test issues-fsm module (cnos#568 AC9)` with `working-directory: src/packages/cnos.issues/commands/issues-fsm` and `run: go test ./... -v`. I ran that exact `cd` + command myself from a fresh checkout and got 20/20 tests PASS (see AC1-AC8 above). I also independently confirmed the γ scaffold's claim that this is *not* redundant with the existing `go` job's `Test` step: `go build ./...`/`go test ./...` run from `src/go` do not reach this module (it's a separate `go.work`-listed module with its own `go.mod`, not a subpackage of `src/go`). I did not actually break a fixture and watch the CI step go red in an actual CI run (no CI runner available to me either), but the working-directory and command are correct and I exercised the identical command locally with a real failure signal available (the assertions are concrete value/field checks, not existence-only checks), so a broken rule genuinely fails the step.

### AC10 — current gates remain green

**Independently verified for everything I could run locally; explicitly not verified for gates that require infrastructure unavailable in this environment.**

Ran myself:
```
$ cd src/go && go build ./...        # exit 0
$ go vet ./...                       # exit 0
$ go test ./... -v                   # all 14 existing internal/* packages PASS, no regressions
$ cd src/packages/cnos.issues/commands/issues-fsm && go vet ./...   # exit 0
$ gofmt -l src/packages/cnos.issues/commands/issues-fsm/*.go src/go/internal/cli/cmd_issues_fsm.go
(no output — clean)
$ ./cn build --check     # I1 package-source-drift
✓ cnos.issues: valid
✓ cnos.cds: valid
✓ All packages valid.
$ ./cn help | grep issues
  issues-map   ...
  issues-fsm   Evaluate CDS issue state from observed facts (read-only; Phase 1, cnos#568)
```

**Not independently re-run (infrastructure unavailable to me, same as α's own honesty note):** I5 skill-frontmatter-check requires `cue`, which is not installed in this environment (`which cue` → exit 1) — I instead read the `src/packages/cnos.issues/SKILL.md` diff directly and confirmed the frontmatter's YAML shape is unchanged (same keys: `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, `visibility`, `triggers`, `scope`, `inputs`, `outputs`, `calls`; only prose/array-item additions, no new/removed/renamed fields), which is the failure mode I5 would most plausibly catch. I2 (protocol-contract-check), I4 (link-check), install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity were not re-run by me in an actual CI runner — none of the changed files fall inside those gates' known scope (no protocol-contract files, no linked docs, no wake-render golden, no repair/closeout logic touched), but I am stating explicitly that this is inference from the diff's file list, not a runtime pass I observed.

**Scope check:** confirmed via `git diff origin/main...origin/cycle/568 --stat` (triple-dot, merge-base-relative) that exactly 24 files changed, all inside: two `.cdd/unreleased/568/` artifacts, the new `issues-fsm` package tree + its 9 testdata fixtures, `transitions.json`, the 4 wiring points (`go.work`, `src/go/go.mod`, `src/go/cmd/cn/main.go`, `cmd_issues_fsm.go`), the AC9 CI step, and the `cnos.issues/SKILL.md` doc update. I initially saw `.cn-sigma/logs/20260703.md` and `docs/development/board/*` in a plain two-dot `git diff origin/main origin/cycle/568`, but confirmed via `git merge-base` + triple-dot diff that these are an artifact of `origin/main` having advanced past the branch point (unrelated commits `45476d24`/`18ec7cb2`), not files α touched — no scope creep.

---

## Summary

All 10 ACs verified independently against the actual code, actual test runs, and actual CLI output — not against α's self-coherence prose. AC5 (the cnos#368 regression guard, the highest-stakes AC in this cycle) holds up under direct inspection of rule ordering, the pure-function engine, and live CLI output on the exact fixture. AC8's zero-mutation boundary and AC1's genuine declarative-data/generic-engine split both hold under grep and direct code reading. No scope creep. No dispatch-boundary violation. `go build`/`go vet`/`go test`/`gofmt -l` all clean, run by me from a fresh worktree.

**verdict: converge**
