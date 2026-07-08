# self-coherence — cnos#630

## Gap

**Issue:** [usurobor/cnos#630](https://github.com/usurobor/cnos/issues/630) — "partial-matter in-progress wedge: give the reconciler a mechanical path out."

**Mode:** design-and-build. **Cell kind:** `implementation` (reconciler gap fix).

**The gap.** A `status:in-progress` dispatch cell whose run dies *after* producing matter (a `cycle/{N}` branch + a checkpointed draft PR, per the cnos#591 finalizer) but *before* `REVIEW-REQUEST.yml` is written had **no mechanical path out**: it could not requeue (the cnos#575 release-back-to-queue guard blocks over existing matter — the cnos#368 blind-requeue protection), could not advance (no `REVIEW-REQUEST.yml`, so `in-progress -> review` is blocked), and was never re-claimed (the dispatch selector only claims `status:todo`). Concrete evidence: cnos#614 sat wedged in `status:in-progress` for over a day across many scheduled sweeps, recoverable only by manual κ (operator) intervention. The literal wedge was `scan.go`'s `else` branch of the `propose_delta_recovery` case (the pre-fix `scan.go:252-256`), which re-evaluated to the identical no-op decision on every scan tick once matter was checkpointed.

**Round:** R0 (first implementation pass on `cycle/630`, `run_class: first_pass` per the γ scaffold — clean first dispatch, no prior claim history).

**Base:** `cycle/630`, created from `main@6143b53c9098e415c4e7e83791843bcad2313314` per `.cdd/unreleased/630/gamma-scaffold.md`.

## Skills

**Tier 1 (always loaded):** `CDD.md`; `cnos.cdd/skills/cdd/alpha/SKILL.md` (this role's contract — §2.1 dispatch intake, §2.5 self-coherence authoring, §2.6 pre-review gate, §2.7 request-review signal, all followed).

**Tier 2 (always-applicable engineering skills, per `cnos.eng/skills/eng/README.md`):** Go authoring conventions (table-driven evaluator discipline — `table.go`'s own doc comment: "never switched on a CDS-specific state name... entirely table-driven" — honored by adding the new rule as declarative JSON data, not a Go `switch` branch); test-first discipline for behavior changes (every fixture/action-name change is backed by a table-driven or `RunScan`-level test).

**Tier 3 (issue-specific):**
- `cnos.core/skills/skill` — not invoked as a distinct skill-authoring pass; the diff amends existing `SKILL.md` prose (`cds-dispatch/SKILL.md`, `delta/SKILL.md`) rather than authoring a new skill file, so the existing files' own section/heading conventions were followed by direct read-and-match rather than a separate skill-authoring skill load.
- The γ scaffold's "Recommended design" and "Open design call" sections (`.cdd/unreleased/630/gamma-scaffold.md`) functioned as the binding Tier-3 design constraint for this cycle — read in full before any code change, per α's dispatch-intake discipline (`alpha/SKILL.md` §2.1 step 4).
- `cnos.cds/skills/cds/fsm/transitions.json`'s own `_doc`/`_doc_phase2`/`_doc_phase3` blocks were read as the binding authoring convention for how new rules must be documented (each `_doc*` block names the issue that added the rule and the rationale) — the new rule in this diff follows that convention (a `reason` string naming cnos#630, the mechanism, and the guardrail it preserves).

**Implementation contract (pinned by δ; not relaxed):** Go; `cn issues fsm evaluate`/`scan` subcommand family, no new binary; package scoping `src/packages/cnos.issues/commands/issues-fsm/` + `src/packages/cnos.cds/skills/cds/fsm/transitions.json` + `cds-dispatch/SKILL.md` + `delta/SKILL.md` §9; existing-binary disposition preserved; no new runtime dependencies; JSON/wire contract additive-only; backward-compat preserved except the one AC2-redirected test (named in §ACs below). Every diff hunk in this cycle maps to one or more of these pinned rows — no improvisation on language, package location, or binary shape.

## ACs

Implementation SHA for all evidence below: `a2e9e7dcb0078947a99fe8891c1d68ab1cfb17b3` (the implementation commit; this self-coherence write continues on top of it).

### Open design calls (resolved first, since AC evidence below depends on both)

**Design call 1 — extend the existing rule in-place, or add a distinct new action/rule?** Resolved: **(ii) distinct new action/rule**, per the γ scaffold's own recommendation. Implemented as `transitions.json`'s new `propose_status_todo_with_matter` rule (`in-progress` state, positioned between the `run_active` valid rule and the pre-existing `propose_delta_recovery` dead-with-matter rule), gated on `all_false: [run_active]` + `all_true: [pr_exists]`. Rationale, restated concretely against the actual guard vocabulary (not just the scaffold's abstract argument): `pr_exists` is the exact boolean that flips from false to true the moment `cn cell finalize` successfully checkpoints matter into a draft PR — before that flip, `branch_has_commits` may be true with `pr_exists` still false (matter exists but is not yet checkpointed; finalize still needs to run), and collapsing the two states into one rule would either (a) skip the finalize step entirely for freshly-dead cells with only commits, or (b) require scan.go to grow a conditional inside the single `propose_delta_recovery` case to decide "call finalize" vs. "just relabel," reintroducing exactly the kind of state-name-shaped branching `table.go`'s doc comment forbids ("never switched on a CDS-specific state name... entirely table-driven"). Two rules, ordered by `pr_exists`, keep the decision entirely in the declarative table: rule N ("checkpoint already done") matches first when `pr_exists` is true; the untouched rule N+1 ("checkpoint not yet done") still matches when only `branch_has_commits`/`pr_has_commits` are true. `scan.go`'s existing generic reconcile branch (`dec.Outcome == "proposed" && dec.TargetState != ""`) already handles the new action with zero new Go code — it was written generically enough (Cases 2/5) to absorb a third case for free. New action-id vocabulary added: `propose_status_todo_with_matter` (target_state `"todo"`), documented in `table.go`'s `Action` field doc comment alongside the pre-existing action-id examples.

**Design call 2 — does AC3's claim-time resume need new Go surface, or is it doctrine-only?** Resolved: **doctrine-only, confirmed** (the scaffold's recommendation). Evidence for "no Go-level data gap": `snapshot.go`'s `FactSnapshot` already carries `BranchExists`, `PRExists`, and `CDDArtifacts` (`[]string` of filenames under `.cdd/unreleased/{N}/`) as first-class JSON fields, all populated by the live path (`fetch.go`'s `assembleLive`) and already exercised by `cn issues fsm evaluate --issue {N} --json`. `transitions.json`'s `todo` state rules (read in full at scaffold time, confirmed again here) gate the `todo -> in-progress` claim rule on `claim_request_present` + `all_false: [run_active]` only — no guard in the `todo` rule list references `branch_exists`/`pr_exists`/`cdd_artifacts` at all, so there is nothing for the claim path to "miss": the FSM already proposes the claim transition regardless of pre-existing matter (confirmed as a regression, not new behavior, by `TestAC630_TodoWithExistingBranchAndPRResumesNotDeferred` below). The actual gap was never a data gap; it was that `cds-dispatch/SKILL.md`'s claim-mechanism prose and `delta/SKILL.md` §9's routing prose never *named* "claimed cell already has branch/PR matter" as an expected, resume-worthy shape — a claiming wake or δ session reading the doctrine cold had no textual instruction to treat it as normal rather than an anomaly to defer. This cycle closes that doctrine gap directly: `cds-dispatch/SKILL.md` gets a new paragraph after claim step 9 naming the shape and pointing to the new `delta/SKILL.md` §9.11; `cds-dispatch/SKILL.md`'s Repair re-entry preflight Step A gets a `resumed_from_matter` `run_class` value, checked *before* the repair-re-entry classification (a genuine risk otherwise: a cnos#630-resumed cell's pre-existing branch-with-commits would literally satisfy the old Step A's repair-re-entry bullet list, sending the wake down the wrong path — loading a nonexistent rejection contract). No Go code changed for AC3 itself; `TestAC630_TodoWithExistingBranchAndPRResumesNotDeferred` is a regression lock proving the FSM layer, not new coverage of new code.

### AC1 — mechanically driven to a valid next state, no human intervention

**Oracle:** a fixture matching "in-progress, matter present (branch+PR), no REVIEW-REQUEST.yml, no live run" produces a concrete non-`""` action on first reconciliation tick.

**Evidence:** `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` (`issuesfsm_test.go`) runs `testdata/scan-died-after-pr-before-review-request.json` (the exact #614-shaped fixture) through `Evaluate` against the real, current `transitions.json` and asserts `Outcome == "proposed"`, `Action == "propose_status_todo_with_matter"`, `TargetState == "todo"`, `EnabledTransition == "in-progress -> todo"`. At the `RunScan` level, `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` (`scan_test.go`) runs the same fixture through the full reconciliation path with a fake GitHub server and asserts `res.Reconciled == true`, `res.TargetState == "todo"`, exactly one label remove+add request pair recorded by the fake server, and zero finalize calls (matter is already checkpointed — nothing further to checkpoint). Both pass (`go test ./...` output below). No human-authored label edit or comment appears anywhere in the fixture/test harness — the label move and the audit-note comment are both produced by the existing `applyStatusLabel`/`PostComment` machinery driven purely by the table's decision.

### AC2 — the #614 wedge fixture: reproduces the strand before the fix, resolves after

**Oracle:** the fixture must demonstrably fail (reproduce the strand) against pre-fix logic and pass (reach a valid next state) against the fix, with zero human intervention on the pass side.

**Evidence — the "before" half.** A single-branch cycle cannot literally run "the old binary" against a fixture that now lives on the same branch as the fix (Friction note 1 in the γ scaffold names this constraint explicitly and leaves the resolution mechanism to α). Resolution taken: `TestAC630_WedgePreFixRuleReproducesStrand` (`issuesfsm_test.go`) constructs a literal, inline `Table` carrying only the pre-cnos#630 `in-progress` dead-with-matter rule (`all_false: [run_active]`, `any_true: [branch_has_commits, pr_exists, pr_has_commits]` → `propose_delta_recovery`, `target_state: ""` — copied verbatim from the rule's pre-diff shape, still visible in this PR's diff as the "old" side of the `transitions.json` hunk) and runs it against `testdata/scan-died-after-pr-before-review-request.json`. Asserts `Outcome == "proposed"`, `Action == "propose_delta_recovery"`, `TargetState == ""` — the permanent-no-op decision that IS the wedge (a `""` target state means `scan.go`'s `propose_delta_recovery` case never applies a status-label move, so the same fixture would re-evaluate to this identical decision on every future scan tick, forever, exactly matching #614's observed >24h stranding across many scheduled sweeps). This test is not a one-time manual verification pasted into this document — it is a permanent, committed regression lock that will keep failing loudly if anyone ever reintroduces the old rule shape.

**Evidence — the "after" half.** `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` (same fixture, real current table) — see AC1 evidence above; reaches `propose_status_todo_with_matter` / `TargetState: "todo"`, a valid next state. `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` confirms the full `RunScan` path reaches this state with zero human-authored intervention (label move and comment both driven by the existing FSM/reconcile machinery).

**Test redirect named explicitly (per AC5's "name exactly which, and why"):** `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` (pre-existing, `scan_test.go`) is REDIRECTED to `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter`. Its old assertion (`Action == "propose_delta_recovery"`, finalize invoked, zero comments) was a literal name-level assertion of the wedge's own steady-state behavior against this exact fixture — the fixture this issue exists to fix. Keeping the old assertion unchanged would be asserting the bug is still present; redirecting it to assert the fix is the correct AC2 closure, not an incidental behavior change. `TestScan_Idempotent` (pre-existing) is also MODIFIED (not renamed) in place: its `checkpointed` fixture sub-case (same `scan-died-after-pr-before-review-request.json`) previously asserted "finalize re-invoked every tick, zero new comments" (the old idempotence-of-the-wedge shape); it now asserts "reconciles on tick 1 with one comment, zero finalize calls, and is absent from tick 2's active-issue listing" — the same test, same fixture, updated expectations to match the new (non-wedged) terminal behavior, following exactly the same "drop from listing once reconciled" convention the test already used for its `deadNoMatter` sub-case.

**A third, non-behavior-changing test was added to preserve coverage**, not to satisfy AC2 directly but disclosed here for completeness: `TestScan_DeltaRecoveryObserveApplyRace_FinalizeNoOpNoComment` keeps `scan.go`'s `propose_delta_recovery` "PR already existed" else-branch (now an observe-vs-apply race fallback, not the steady-state path) under test coverage, using `in-progress-dead-with-commits.json` (a fixture with `pr_exists: false` at observation time) plus a fake finalizer stdout simulating a PR that appeared between observation and the finalize call.

### AC3 — claim-time resume, not defer

**Oracle:** `status:todo` + pre-existing `branch_exists`/`pr_exists` + `run_active: false` still proposes `todo -> in-progress` (FSM-level regression, per Key finding 2), AND the claim-doctrine layer demonstrably describes resuming from existing matter, not merely allowing the label transition.

**FSM-level evidence:** `testdata/todo-resume-existing-matter.json` (new fixture, mirrors `todo-competing-run.json`'s shape but `run_state: ""` and `branch_exists`/`pr_exists: true`) run through `TestAC630_TodoWithExistingBranchAndPRResumesNotDeferred` (`issuesfsm_test.go`) asserts `Outcome == "proposed"`, `TargetState == "in-progress"`, `Action == "propose_status_in_progress"` — confirming Key finding 2's claim (no `todo`-state rule change was needed or made; `git diff` on `transitions.json`'s `todo` block is empty, confirmed by inspection of the diff hunk, which touches only the `in-progress` block).

**Doctrine-level evidence (the actual AC3 gap, per Design call 2 above):** `cds-dispatch/SKILL.md`'s claim-mechanism section gets a new paragraph (after step 9) naming the pre-existing-matter shape as normal and pointing to `delta/SKILL.md` §9.11. `delta/SKILL.md` §9.11 ("resumed-from-mechanical-reversion shape") is new: it names the trigger (a scanner-driven `status:in-progress -> status:todo` reconciliation, not an operator rejection), gives a comparison table against the structurally similar but causally distinct §9.10 "resumed-from-changes shape," names the detection mechanism (read the cycle-branch state at spawn time — no new Go surface, `FactSnapshot` already exposes what's needed per Design call 2), and gives the routing sequence: δ reads existing `.cdd/unreleased/{N}/` artifacts first and does NOT re-dispatch γ for a fresh scaffold when one already exists. This satisfies the scaffold's own oracle text verbatim: "the label transition alone does NOT satisfy AC3 — the resume behavior is what AC3 requires" — the resume *behavior* (reading existing state, not re-scaffolding) is what the new §9.11 names, distinct from the (unchanged, already-passing) label transition.

### AC4 — machine-readable audit note; not misread as an unexplained orphan

**Oracle:** every new mechanical reversion posts an audit-note comment; a subsequent claim/scan pass does not classify it as an unexplained orphan.

**Comment-posting evidence:** the new rule's `reason` text is posted verbatim as the recovery-scan comment body via `scan.go`'s pre-existing, unmodified reconcile branch (`body := fmt.Sprintf("cds recovery scan: reconciled #%d -- %s.\n\n%s", issue, dec.EnabledTransition, dec.Reason)`), which was already the mechanism Cases 2/5 used — no new `PostComment` call site. `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` confirms exactly one comment is posted (`com.count() == 1`).

**"Names it a mechanical reversion" evidence:** `TestAC630_AuditNoteReasonNamesMechanicalReversion` (`issuesfsm_test.go`) asserts the rule's `Reason` string contains the literal substrings `"MECHANICAL reversion"`, `"PRESERVED"`, and `"cnos#368"` — a programmatic lock on the audit-note's content, not just its existence.

**"Not misread as an orphan" evidence:** this is the doctrine half again — `cds-dispatch/SKILL.md`'s Repair re-entry preflight Step A now explicitly checks for the scanner's "MECHANICAL reversion" comment phrase (and the absence of operator-rejection evidence) BEFORE the repair-re-entry classification, and routes that case to `run_class: resumed_from_matter` → `delta/SKILL.md` §9.11, rather than into Steps B–E's rejected-findings machinery (which would misfire: there is no rejection to load a repair contract for). This is a textual/doctrine oracle, not a Go test, per the scaffold's own allowance ("If no Go-testable surface exists for the doctrine-only half, the oracle is the doctrine text itself").

### AC5 — no regression

**Package-level:** `go test ./...` from `src/packages/cnos.issues/commands/issues-fsm` — full output pasted below (§Self-check has the exact command + output; summary: 77 tests run, 77 pass, 0 fail, `ok` in 0.042s; `go test ./... -race` also green, 1.196s).

**Tests deliberately redirected (named exactly, per AC5's own requirement):**
1. `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` → `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` (rename + behavior-assertion flip). Why: this test's old assertions were a literal specification of the wedge bug against the exact fixture this issue fixes; asserting the fix's behavior against the same fixture is AC2's own requirement, not an incidental behavior change.
2. `TestScan_Idempotent` (kept its name; assertions updated in place for the `checkpointed` sub-case only — the `deadNoMatter` and `blocked` sub-cases are unchanged). Why: the `checkpointed` fixture's old expected idempotence shape ("finalize re-invoked forever, no new comment") was the wedge's own steady-state; the new expected shape ("reconciles once, absent from the next scan's active-issue listing") is the correct idempotence property of the fix, using the exact same "drop from listing once reconciled" convention this test already applied to its other sub-case.

No other pre-existing test's assertions changed. Confirmed by re-running the full suite (§Self-check) and by the fixture-level audit above (§ACs, Design call 1) showing only one fixture (`scan-died-after-pr-before-review-request.json`, issue #9601) in the entire `testdata/` directory has `pr_exists: true` + dead-run + `review_request_present: false` + no release/block markers — i.e. the new rule's guard combination is provably reachable by exactly one pre-existing fixture, confirmed by a table scan of every fixture's `pr_exists`/`run_state`/`review_request_present`/`release_request_present`/`block_request_present` fields (all 22 pre-existing fixtures under `testdata/` enumerated; only issue #9601's fixture matches).

`TestSeam_CellKindNotEnforced` — located via `grep -rn TestSeam_CellKindNotEnforced src/` → `issuesfsm_test.go:810`; included in the full-suite run above; passes.

`check-dispatch-repair-preflight.sh` (#516) and `check-dispatch-closeout-integrity.sh` (#524) — both run directly (`bash scripts/ci/check-dispatch-repair-preflight.sh` / `bash scripts/ci/check-dispatch-closeout-integrity.sh`), both exit 0.

**Broader gate:** `cd src/go && go build ./... && go vet ./...` — clean, no output, exit 0. All four `go.work` modules (`src/go`, `src/packages/cnos.cdd/commands/cdd-verify`, `src/packages/cnos.issues/commands/issues-map`, `src/packages/cnos.issues/commands/issues-fsm`) individually build + vet clean. `go test ./...` from `src/go` — all 15 internal packages `ok`, zero failures (no packages under `src/go` import `issues-fsm`, confirmed no cross-module regression).

**install-wake-golden gate:** the rendered `.github/workflows/cnos-cds-dispatch.yml` and the per-package golden `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` were both re-rendered via `cn-install-wake cds-dispatch` (the shell-script renderer at `src/packages/cnos.core/commands/install-wake/cn-install-wake`, built from `/tmp/cn` per the dispatch prompt — actually invoked directly as the shell script itself, since `/tmp/cn`'s compiled kernel binary does not expose `install-wake` as a subcommand at all; the renderer is a package command, not a kernel command). `sha256sum` of the live workflow and the golden are byte-identical (`78711c5c...`). A second consecutive render of the golden produces the identical hash (idempotence, matching the CI gate's own "Verify idempotence" step). Both goldens parse as valid YAML (`yaml.safe_load`). The substrate structural-shape needles the CI gate checks (`^on:`, `^permissions:`, `^concurrency:`, `^jobs:`, `anthropics/claude-code-action@v1`, `types: \[labeled\]`, `^  group: cds-dispatch-sigma$`) are all present, replicating `install-wake-golden.yml`'s own grep checks locally.

### AC6 — idempotent

**Oracle:** running the new reconciliation path twice against the same fixture state produces zero new mutations/comments/PRs on the second pass.

**Evidence:** `TestScan_Idempotent` (`scan_test.go`, updated per AC5 above) now covers this directly: first `RunScan` pass against `{deadNoMatter, checkpointed, blocked}` produces exactly 2 comments (one per requeue path, including the new with-matter requeue) and 0 finalize calls; the second pass — simulating the post-reconciliation world by dropping both reconciled issues from the active-issue listing, exactly mirroring the convention the test already used for `deadNoMatter` — produces zero new comments, zero new finalize calls, and zero new label requests (`secondLabelCount == 0`, asserted). This is the same "state moved out of the active-scan set" idempotence proof the pre-existing `deadNoMatter` case already used, extended to the new with-matter requeue path per the scaffold's explicit instruction not to duplicate the finalizer's own idempotence coverage (already proven by `TestScan_DeltaRecoveryObserveApplyRace_FinalizeNoOpNoComment` and the pre-existing `TestScan_DeadWithMatterNoPR_FinalizesAndComments`) and to test only the new requeue/resume path's idempotence.

## Self-check

**Role self-check.** Did this round push ambiguity onto β? Both open design calls the scaffold flagged are resolved and documented above with concrete rationale tied to the actual guard vocabulary (not restated scaffold text) — β should not need to re-derive either decision from scratch, only verify it. Every AC claim above cites a specific test name and file; β can re-run each one independently rather than trusting narrative. The one area of judgment β should specifically re-examine: whether the `resumed_from_matter` / §9.11 doctrine addition is *sufficiently* load-bearing to satisfy AC3, since (per the scaffold's own Friction note 2) this is a doctrine-quality question, not a code-coverage one, and reasonable readers could disagree on how explicit is explicit enough.

**Exact commands run and their output** (all from `/home/runner/work/cnos/cnos` unless noted; SHA at time of this run: `a2e9e7dcb0078947a99fe8891c1d68ab1cfb17b3` plus this document's own commits on top):

```
$ cd src/packages/cnos.issues/commands/issues-fsm && go build ./... && go vet ./...
(no output; exit 0)

$ go test ./... -v 2>&1 | grep -c '^--- PASS'
77
$ go test ./... -v 2>&1 | grep -c '^--- FAIL'
0
$ go test ./...
ok  	github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm	0.042s -0.187s (varies by cache state)

$ go test ./... -race
ok  	github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm	1.196s

$ go test ./... -v 2>&1 | grep -E 'TestAC630|TestScan_Dead(WithCheckpointedPR|NoMatter)|TestScan_DeltaRecoveryObserveApplyRace|TestScan_Idempotent|TestSeam_CellKindNotEnforced'
--- PASS: TestScan_DeadNoMatter_RequeuesToTodo (0.00s)
--- PASS: TestScan_DeltaRecoveryObserveApplyRace_FinalizeNoOpNoComment (0.00s)
--- PASS: TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter (0.00s)
--- PASS: TestScan_Idempotent (0.00s)
--- PASS: TestAC630_WedgePreFixRuleReproducesStrand (0.00s)
--- PASS: TestAC630_WedgeFixResolvesToTodoWithMatterPreserved (0.00s)
--- PASS: TestAC630_AuditNoteReasonNamesMechanicalReversion (0.00s)
--- PASS: TestAC630_TodoWithExistingBranchAndPRResumesNotDeferred (0.00s)
--- PASS: TestSeam_CellKindNotEnforced (0.00s)
```

```
$ bash scripts/ci/check-dispatch-repair-preflight.sh
cnos#516 repair-preflight guard: dispatch surface carries the repair re-entry contract (protocol + prompt + golden + live workflow).
(exit 0)

$ bash scripts/ci/check-dispatch-closeout-integrity.sh
cnos#524 closeout-integrity guard: dispatch surface carries the deliverable-proof contract (protocol + prompt + SKILL + golden + live). The empty-review detector itself is proven live by the Go FSM test suite (see header note) rather than by a bash self-test.
(exit 0)
```

```
$ for d in src/go src/packages/cnos.cdd/commands/cdd-verify src/packages/cnos.issues/commands/issues-map src/packages/cnos.issues/commands/issues-fsm; do (cd "$d" && go build ./... && go vet ./...); done
(all four modules: no output, exit 0)

$ cd src/go && go test ./...
ok  	github.com/usurobor/cnos/src/go/internal/activate	0.038s
ok  	github.com/usurobor/cnos/src/go/internal/activation	0.069s
ok  	github.com/usurobor/cnos/src/go/internal/binupdate	0.018s
ok  	github.com/usurobor/cnos/src/go/internal/cell	0.017s
ok  	github.com/usurobor/cnos/src/go/internal/cli	2.631s
ok  	github.com/usurobor/cnos/src/go/internal/discover	0.014s
ok  	github.com/usurobor/cnos/src/go/internal/dispatch	0.006s
ok  	github.com/usurobor/cnos/src/go/internal/doctor	0.308s
ok  	github.com/usurobor/cnos/src/go/internal/hubinit	0.039s
ok  	github.com/usurobor/cnos/src/go/internal/hubsetup	0.008s
ok  	github.com/usurobor/cnos/src/go/internal/hubstatus	0.016s
ok  	github.com/usurobor/cnos/src/go/internal/pkg	0.004s
ok  	github.com/usurobor/cnos/src/go/internal/pkgbuild	0.110s
ok  	github.com/usurobor/cnos/src/go/internal/repoinstall	0.974s
ok  	github.com/usurobor/cnos/src/go/internal/restore	0.036s
(cmd/cn: [no test files])
```

```
$ sha256sum .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
78711c5c5aaaef3270d532dfd1daeda750b24c22124ecb6bf40f620d9465c817  .github/workflows/cnos-cds-dispatch.yml
78711c5c5aaaef3270d532dfd1daeda750b24c22124ecb6bf40f620d9465c817  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
(byte-identical, confirms the render committed matches the live workflow)
```

**Artifact-enumeration-matches-diff check (pre-review gate row 11).** `git diff --stat origin/main..HEAD` (run just before this section was written) against this section's + §ACs' coverage:

| File | Named above? |
|---|---|
| `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | Yes — §ACs Design call 1, AC1-AC6 |
| `src/packages/cnos.issues/commands/issues-fsm/table.go` | Yes — §ACs Design call 1 (Action doc comment) |
| `src/packages/cnos.issues/commands/issues-fsm/scan.go` | Yes — §ACs AC1/AC4 (comment updates only, no new branch) |
| `src/packages/cnos.issues/commands/issues-fsm/scan_test.go` | Yes — §ACs AC2/AC5/AC6 (redirected + new tests) |
| `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` | Yes — §ACs AC1/AC2/AC3/AC4 (new `TestAC630_*` tests) |
| `src/packages/cnos.issues/commands/issues-fsm/testdata/todo-resume-existing-matter.json` | Yes — §ACs AC3 |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | Yes — §ACs Design call 2 / AC3 / AC4 |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | Yes — §ACs AC5 (install-wake-golden gate) |
| `.github/workflows/cnos-cds-dispatch.yml` | Yes — §ACs AC5 (install-wake-golden gate) |
| `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | Yes — §ACs Design call 2 / AC3 |

**Caller-path trace for new modules (pre-review gate row 12).** No new Go module or package was added this cycle — `propose_status_todo_with_matter` is a new *action-id string* (data, not code), consumed by `scan.go`'s pre-existing generic reconcile branch (`scanOne`, the `case dec.Outcome == "proposed" && dec.TargetState != ""` arm, unchanged code, already the caller for Cases 2/5). No new Go function was added without a non-test caller; every new test function's caller is the Go test runner itself (the standard, expected caller for `Test*` functions).
