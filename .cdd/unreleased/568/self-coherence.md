# α self-coherence — cnos#568 (cds/issues: FSM Phase 1 — read-only issue-state reconciler)

**Cycle:** cycle/568 · base `origin/cycle/568@ce478656` (γ scaffold commit) · run_class `first_pass`.

## §R0 — AC1–AC10 evidence

### AC1 — FSM transition table exists (single declarative source)

**Met.** The transition table exists as pure JSON data at
`src/packages/cnos.cds/skills/cds/fsm/transitions.json` (states, guards,
and per-state `rules[]` with `all_true`/`any_true`/`all_false`
guard-condition lists, `outcome`/`action`/`target_state`/`reason` — no Go
literals). It is loaded by `Table`/`LoadTable`
(`src/packages/cnos.issues/commands/issues-fsm/table.go:64-82`) and
consumed by `Evaluate` (`table.go:157-210`), which looks up
`t.Transitions` by the state string parsed at runtime from
`FactSnapshot.Labels` (`CurrentState`, `snapshot.go:64-72`) — there is no
`switch`/`if`-chain on a CDS state name (`"review"`, `"changes"`,
`"in-progress"`, ...) anywhere in `table.go` or `decision.go`; the only
place a Go string literal names a CDS state is in test fixtures
(`testdata/*.json`) and test assertions. Verified by grep:

```
$ grep -n '"review"\|"changes"\|"in-progress"' src/packages/cnos.issues/commands/issues-fsm/*.go | grep -v _test.go
(no output)
```

Guard *names* referenced by the table (`run_active`, `branch_has_commits`,
`pr_exists`, ...) are the fixed, generic predicate vocabulary in
`guardFuncs` (`table.go:90-100`) — these are evidence primitives over the
generic `FactSnapshot`, not CDS-specific state names, satisfying the
package-ownership split: `cnos.cds` owns the table, `cnos.issues` owns the
engine. Test: `TestAC1_TableLoadsAsData`
(`issuesfsm_test.go`) — confirms the table parses and contains all five
states.

### AC2 — `cn issues fsm evaluate --issue {N}` exists and explains

**Met.** Registered as a compiled-in kernel command: `src/go/cmd/cn/main.go:48`
(`reg.Register(&cli.IssuesFsmCmd{})`), `Name: "issues-fsm"`
(`src/go/internal/cli/cmd_issues_fsm.go:28`), which noun-verb resolves to
`cn issues fsm ...` per `ResolveCommand`. `Run` dispatches on the first
remaining arg (`issuesfsm.go:44-58`); `evaluate` is the only recognized
sub-verb (`issuesfsm.go:50`). `runEvaluate` (`issuesfsm.go:62-127`) supports
`--issue`, `--fixture` (offline), `--table`, `--repo`, `--token`.
`Decision.Render` (`decision.go:44-88`) prints, in order: current state,
observed facts (every FactSnapshot field), and the decision block
(`outcome`, `enabled_transition`, `blocked_reason`, `missing_evidence`,
`proposed_action`, `reason`).

Live-run proof (built binary, real CLI path, auto-discovered table — no
`--table` flag needed):

```
$ ./cn issues fsm evaluate --issue 601 \
    --fixture src/packages/cnos.issues/commands/issues-fsm/testdata/review-empty.json
cn issues fsm evaluate — issue #601

Current state: review

Observed facts:
  labels: dispatch:cell, protocol:cds, status:review
  run_state: completed
  run_conclusion: success
  branch_exists: true
  commits_beyond_base: 0
  pr_exists: false
  pr_commit_count: 0
  cdd_artifacts: gamma-scaffold.md
  review_request_present: false
  repair_contract_present: false
  checks_state: (none)

Decision:
  outcome: blocked
  enabled_transition: (none)
  blocked_reason: empty review: status:review with no PR, no commits beyond base, and no REVIEW-REQUEST.yml present.
  missing_evidence: pr_exists, branch_has_commits, review_request_present
  proposed_action: block
  reason: empty review: status:review with no PR, no commits beyond base, and no REVIEW-REQUEST.yml present.

(read-only: no label was written; this is a proposal, not an action)
```

Also verified: `./cn issues` (no verb) lists both `issues map` and
`issues fsm` as a group (registry co-listing works); `./cn issues fsm apply
--issue 1` is rejected (exit 1, explains apply is Phase 2 / cnos#569 — see
AC8). Test: `TestAC2_RunPrintsFullDecisionBlock`,
`TestAC2_UnknownSubcommandRejected` (`issuesfsm_test.go`).

### AC3 — empty-review is detected

**Met.** Fixture `testdata/review-empty.json`: `status:review`, `pr_exists:
false`, `commits_beyond_base: 0`, `review_request_present: false`. The
`review` state's second rule (`transitions.json`, `all_false:
["pr_exists","branch_has_commits","review_request_present"]`) matches,
producing `outcome: blocked` with `missing_evidence: [pr_exists,
branch_has_commits, review_request_present]` (all three, since all three
guards are false). Test: `TestAC3_EmptyReviewBlocked` — asserts
`Outcome == "blocked"`, non-empty `BlockedReason`, and the exact 3-element
missing-evidence set. Negative counter-test:
`TestAC3_ReviewWithEvidenceIsValid` (fixture `review-with-pr.json`, PR
present) asserts `Outcome == "valid"` — proves the guard is discriminating,
not always-blocked.

### AC4 — dead in-progress / no matter → proposes requeue

**Met.** Fixture `testdata/in-progress-dead-no-matter.json`: `status:
in-progress`, `run_state: completed` (not active), `branch_exists: false`,
`pr_exists: false`. The `in-progress` state's third (fallback) rule fires
(`transitions.json`), producing `outcome: proposed`, `action:
propose_status_todo`, `target_state: todo`, `enabled_transition: "in-progress
-> todo"`. Test: `TestAC4_DeadInProgressNoMatterProposesRequeue` — asserts
`Outcome == "proposed"`, `TargetState == "todo"`,
`EnabledTransition == "in-progress -> todo"`, `Action ==
"propose_status_todo"`.

### AC5 — dead in-progress / branch-with-commits → proposes δ-recovery

**Met.** Fixture `testdata/in-progress-dead-with-commits.json`: same as
AC4 but `commits_beyond_base: 4`. The `in-progress` state's *second* rule
(`all_false: ["run_active"]`, `any_true: ["branch_has_commits", "pr_exists",
"pr_has_commits"]`) matches first (rule order in `transitions.json` puts
this before the requeue fallback), producing `outcome: proposed`, `action:
propose_delta_recovery`, `target_state: ""` (no direct label move — the
proposal is "open/request a recovery PR", never a status flip). Test:
`TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery` — asserts
`Action == "propose_delta_recovery"` **and explicitly asserts
`TargetState != "todo"`**, directly guarding the cnos#368 regression named
in the AC. Positive companion:
`TestAC5_HealthyActiveInProgressIsValid` (fixture
`in-progress-active.json`, `run_state: in_progress`) asserts
`Outcome == "valid"`, `Action == "none"` — proves the guard distinguishes
live runs from dead ones.

### AC6 — `changes → todo` without repair context is blocked

**Met.** Fixture `testdata/changes-no-repair.json`:
`repair_contract_present: false` → the `changes` state's second (fallback)
rule matches, `outcome: blocked`, `target_state: ""`. Test:
`TestAC6_ChangesWithoutRepairContextBlocked` asserts `Outcome ==
"blocked"` and `TargetState != "todo"`. Positive: fixture
`testdata/changes-with-repair.json` (`repair_contract_present: true`) →
first rule matches, `outcome: proposed`, `target_state: "todo"`,
`repair_pass: true`. Test:
`TestAC6_ChangesWithRepairContextEnablesRepairPass` asserts `TargetState ==
"todo"` and `RepairPass == true`.

### AC7 — idempotence

**Met.** `TestAC7_Idempotent` runs both layers twice per fixture, across
all five non-trivial fixtures: (a) `Evaluate(tab, snap)` called twice,
asserting `Outcome`/`Action`/`TargetState`/`Reason` are identical; (b) the
full `Run()` CLI entry point invoked twice with identical args, asserting
the rendered stdout bytes are byte-identical (`out1.String() ==
out2.String()`). This holds structurally, not just empirically: `Evaluate`
and `Decision.Render` are pure functions of `(Table, FactSnapshot)` with no
mutation of any input, no I/O side effect beyond reading the two input
files, and no global/hidden state — there is nothing in the code path that
could cause the second call to diverge.

### AC8 — no mutation, no authority flip (Phase-1 boundary)

**Met, verified three ways.**

1. **No `--apply` flag.** `runEvaluate`'s `flag.FlagSet`
   (`issuesfsm.go:62-71`) registers exactly `--issue`, `--fixture`,
   `--table`, `--repo`, `--token`. `TestAC8_NoApplyFlag` asserts passing
   `--apply` is rejected as an unknown flag (exit non-nil).
2. **No label-write call in source.** `TestAC8_NoLabelMutationCodeInSource`
   scans every non-test `.go` file in the package for `fs.Bool("apply"`,
   `fs.String("apply"`, `"gh issue edit"`, `"gh label"`, `"-X PATCH"`,
   `"-X POST"`, `"/labels\""`, `"AddLabel"`, `"RemoveLabel"` — none present.
   Manually confirmed: the only `net/http` calls in the package
   (`fetch.go`'s `ghGetJSON`) are all `http.MethodGet`; there is no
   `http.MethodPatch`/`http.MethodPost` anywhere in the tree.
3. **Protocol files byte-identical to `origin/main`.** Verified directly
   (not just claimed):

   ```
   $ git diff origin/main -- .github/workflows/cnos-cds-dispatch.yml \
       src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md \
       src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md
   (no output — zero diff)
   ```

   None of these three files appear in `git status --short` for this
   branch's changes either.

### AC9 — CI guard for evaluator fixtures

**Met.** New step added to the `go` job in `.github/workflows/build.yml`
(lines 38-46): `Test issues-fsm module (cnos#568 AC9)`, `working-directory:
src/packages/cnos.issues/commands/issues-fsm`, `run: go test ./... -v`.
Confirmed empirically (not assumed) that this is *not* redundant with the
existing `Test` step: `go list ./...` from `src/go` has zero packages
under `issues-fsm` or `issues-map` (Go workspace `./...` is CWD-relative),
matching the γ scaffold's friction note. Demonstration that the new step
would actually go red on a regression: any fixture's expected `Outcome`/
`Action`/`TargetState` value is asserted by name in
`issuesfsm_test.go` (e.g. `TestAC5_...` fails loudly if
`transitions.json`'s in-progress rule order were reordered to let a
dead-with-commits case fall through to the requeue rule) — verified by
temporarily reordering the two `in-progress` rules locally and re-running
`go test ./...`, which failed as expected, then reverting.

### AC10 — current gates remain green

**Met, verified locally; see exact commands + output below.** `cd src/go
&& go build ./... && go vet ./... && go test ./... -v` — build clean, vet
clean, all pre-existing packages still pass (no regressions; the new
module isn't in this run's package set — expected, see AC9). `./cn help`
smoke-runs and lists `issues-fsm` alongside every other kernel command.
`./cn build --check` (I1 package-source-drift) reports `✓ cnos.issues:
valid` / `✓ cnos.cds: valid` / `✓ All packages valid.`. Only 5 files
touched outside the two new package trees
(`.github/workflows/build.yml`, `go.work`, `src/go/cmd/cn/main.go`,
`src/go/go.mod`, `src/packages/cnos.issues/SKILL.md`) — no unrelated files
modified (`git diff --stat` against the γ-scaffold base confirms this).

## Exact AC10 verification commands + output

```
$ cd src/go && go build ./...
(exit 0, no output)

$ go vet ./...
(exit 0, no output)

$ go test ./... -v
... [all internal/* packages] ...
PASS
ok  	github.com/usurobor/cnos/src/go/internal/activate	0.038s
ok  	github.com/usurobor/cnos/src/go/internal/activation	0.050s
ok  	github.com/usurobor/cnos/src/go/internal/binupdate	0.018s
ok  	github.com/usurobor/cnos/src/go/internal/cell	0.011s
ok  	github.com/usurobor/cnos/src/go/internal/cli	0.010s
ok  	github.com/usurobor/cnos/src/go/internal/discover	0.010s
ok  	github.com/usurobor/cnos/src/go/internal/dispatch	0.004s
ok  	github.com/usurobor/cnos/src/go/internal/doctor	0.192s
ok  	github.com/usurobor/cnos/src/go/internal/hubinit	0.031s
ok  	github.com/usurobor/cnos/src/go/internal/hubsetup	0.006s
ok  	github.com/usurobor/cnos/src/go/internal/hubstatus	0.008s
ok  	github.com/usurobor/cnos/src/go/internal/pkg	0.007s
ok  	github.com/usurobor/cnos/src/go/internal/pkgbuild	0.024s
ok  	github.com/usurobor/cnos/src/go/internal/restore	0.020s

$ go build -o cn ./cmd/cn && ./cn help | head -3
cn — cnos kernel

Usage: cn <command> [args...]
$ ./cn help | grep issues
  issues-map   Generate the interactive issue-board (Voronoi + Pivot) from open issues
  issues-fsm   Evaluate CDS issue state from observed facts (read-only; Phase 1, cnos#568)

$ cd ../packages/cnos.issues/commands/issues-fsm && go test ./... -v
... [16 tests, all PASS, see AC1-AC9 above] ...
PASS
ok  	github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm	0.006s

$ gofmt -l src/packages/cnos.issues/commands/issues-fsm/*.go src/go/internal/cli/cmd_issues_fsm.go
(no output — all formatted)
```

## Honesty notes (partial / not exhaustively verified)

- **Live fact-assembly path (`fetch.go`) has no dedicated unit test.** Per
  the γ scaffold's implementation contract ("the live path can be less
  thoroughly tested... keep it honest but it's fine if it's simpler"), the
  live path (GitHub REST for run/PR state, `os/exec` git for branch/diff)
  compiles, vets clean, and is exercised implicitly by `go build`/`go vet`,
  but is not covered by an `httptest`-backed test the way `commands/issues-
  map`'s `fetchIssues` also isn't. This mirrors the existing precedent's
  test coverage shape, not a new gap.
- **`repair_contract_present` on the live path is a heuristic**
  (`fetch.go`: presence of `delta-repair.md`, or a `beta-review.md`
  containing the literal substring `"verdict: iterate"`), not a full parse
  of β/δ findings. The fixture path (what AC6 actually gates) takes the
  field as a direct, explicit input and is not affected by this
  simplification.
- **I did not run the full `.github/workflows/build.yml` gate suite in
  CI** (no CI runner available in this environment) — I1
  (`./cn build --check`) and the dispatch-boundary grep were run locally
  and reproduce the CI step's exact commands; I5 (skill-frontmatter-check)
  requires `cue`, which was not available in this environment, so I could
  not run it directly — I instead verified the frontmatter YAML parses
  and is structurally unchanged (same field shapes, only array items and
  prose added) relative to the pre-existing, presumably-passing
  `cnos.issues/SKILL.md`. β should re-run I5 directly against the diff.
- **No `--json` flag.** Per the implementation contract, this was
  explicitly optional/non-AC and skipped to keep scope tight; the
  human-readable block is the binding oracle for AC2.

## Files added / changed

Added:
- `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (AC1 data)
- `src/packages/cnos.issues/commands/issues-fsm/go.mod`
- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm.go` (CLI entry / Run)
- `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` (FactSnapshot + fixture loader)
- `src/packages/cnos.issues/commands/issues-fsm/table.go` (Table + generic Evaluate engine)
- `src/packages/cnos.issues/commands/issues-fsm/decision.go` (Decision + Render)
- `src/packages/cnos.issues/commands/issues-fsm/fetch.go` (live GitHub REST + git assembly)
- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` (AC1–AC8 + positive-case tests)
- `src/packages/cnos.issues/commands/issues-fsm/testdata/*.json` (9 fixtures)
- `src/go/internal/cli/cmd_issues_fsm.go` (thin dispatch shim)

Changed:
- `src/go/cmd/cn/main.go` (+1 line: register `IssuesFsmCmd`)
- `src/go/go.mod` (+require/replace for the new module)
- `go.work` (+1 line: `use` the new module)
- `.github/workflows/build.yml` (+1 CI step, AC9)
- `src/packages/cnos.issues/SKILL.md` (document ownership split + command-dispatch disposition for `issues-fsm`)
