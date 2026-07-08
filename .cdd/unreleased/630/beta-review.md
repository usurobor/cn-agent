# beta-review — cnos#630

## §R0 (round 0)

**Reviewed:** `cycle/630` at commit `a63bc9d2331975f8f9e5467b5bb54c9199b89d8d` (§Review-readiness round 1 signal) against `origin/main` at `e059014804229c43f0595af19916e12a86a3e198` (= the branch's own merge-base — confirmed via `git merge-base origin/main cycle/630` == `git rev-parse origin/main`, zero drift at review time).
**Review surface:** `implementation`-kind cell per `CELL-KINDS.md` §2 — reconciler gap fix. Package scoping per the pinned implementation contract: `src/packages/cnos.issues/commands/issues-fsm/` + `src/packages/cnos.cds/skills/cds/fsm/transitions.json` + `cds-dispatch/SKILL.md` + `delta/SKILL.md` §9.
**Diff shape:** `git diff --stat origin/main...cycle/630` — 12 files changed, 913 insertions(+), 68 deletions(-). Matches the 10-file (excl. the two `.cdd/unreleased/630/` artifacts) enumeration in self-coherence.md's §Self-check table exactly; no unaccounted file.

I did not take α's self-coherence.md claims on trust. Every AC below was re-derived from the diff and re-run from the shell, independently of the narrative.

### AC1 — mechanical next action, not permanent no-op

**Oracle:** in-progress + matter (branch/PR) + no REVIEW-REQUEST.yml + no live run produces a concrete non-`""` action on first tick.

**Independent evidence.** Read `transitions.json` lines 184–192 directly (not via self-coherence.md's paraphrase): the new rule (`all_false: ["run_active"]`, `all_true: ["pr_exists"]` → `action: "propose_status_todo_with_matter"`, `target_state: "todo"`) sits between the `run_active → none` rule (178–183) and the old `propose_delta_recovery` dead-with-matter rule (193–201). Re-ran `go test ./... -v -run TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` from `src/packages/cnos.issues/commands/issues-fsm` myself: **PASS**. Also independently re-ran `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` (full `RunScan` path against a fake GitHub server): **PASS** — asserts `Reconciled == true`, `TargetState == "todo"`, exactly 2 label HTTP requests (remove `status:in-progress` + add `status:todo`), 0 finalize calls. **Pass.**

### AC2 — #614-shaped fixture: strand reproduced pre-fix, recovered post-fix, zero human edits

**Oracle:** the fixture must demonstrably fail (reproduce the strand) pre-fix and pass (reach a valid state) post-fix, no human-authored label/comment edit anywhere in the test path.

**Read the fixture directly:** `testdata/scan-died-after-pr-before-review-request.json` (issue `9601`, `run_conclusion: failure`, `pr_exists: true`, `review_request_present: false`) — this is genuinely the #614 shape (dead run, checkpointed PR, no review-request), not a new invented case.

**The "must fail before the fix" half — checked with real skepticism, per the scaffold's own friction-note flag.** α's chosen resolution is `TestAC630_WedgePreFixRuleReproducesStrand` (`issuesfsm_test.go`), which builds an inline `Table` literal carrying only the verbatim pre-diff rule shape (copied from `git diff`'s own `-` side: `all_false:[run_active]`, `any_true:[branch_has_commits,pr_exists,pr_has_commits]` → `propose_delta_recovery`, `target_state:""`) and runs the real #614 fixture through it, asserting `TargetState == ""`. I independently diffed this inline literal against `git diff origin/main...cycle/630 -- .../transitions.json`'s removed lines and confirm it is a faithful, unmodified copy of the old rule — not a strawman. This is a legitimate mechanical demonstration of the strand (not merely an assertion in prose): it is a committed, permanent regression lock that fails loudly if the old rule shape is ever reintroduced. It is not literally "run the old binary" (impossible on a single branch, as the scaffold's friction note acknowledged), but it satisfies the oracle's substance — I re-ran it myself: **PASS** (`dec.TargetState == ""`, the exact permanent no-op that is the wedge).

**The "reaches a valid state, zero human edits" half:** `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` + `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` (same evidence as AC1) — both driven purely by `Evaluate`/`RunScan` against fixtures and a fake HTTP server; no human-authored label edit or comment appears anywhere in the test/fixture path (I read both test bodies in full — confirmed).

**Redirected test named and justified, independently checked:** `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` → renamed/repurposed to `TestScan_DeltaRecoveryObserveApplyRace_FinalizeNoOpNoComment`, now run against a *different* fixture (`in-progress-dead-with-commits.json`, `pr_exists: false`) to preserve coverage of the old `else`-branch as an observe-vs-apply race fallback rather than the steady-state path. I confirmed by reading `scan_test.go`'s diff directly that this is not silently repurposing the *original* #614 fixture's assertions to something toothless — the #614 fixture itself moved to the new `TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter` test with flipped (fixed) assertions. **Pass.**

### AC3 — todo→in-progress claim regression + genuine resume-not-defer mechanism

**Oracle:** claiming still works over pre-existing branch/PR (regression check); the resume behavior (not just the label transition) is described somewhere, mechanism-named not vague.

**Regression check, independently confirmed:** `git diff origin/main...cycle/630 -- src/packages/cnos.cds/skills/cds/fsm/transitions.json` shows the diff touches only the `in-progress` block (lines 184–200 of the hunk); the `todo` state's 3 rules are byte-for-byte unchanged in the diff. I independently ran `TestAC630_TodoWithExistingBranchAndPRResumesNotDeferred` against the new fixture `testdata/todo-resume-existing-matter.json` (`branch_exists: true`, `pr_exists: true`, `run_state: ""`, `claim_request_present: true`): **PASS** (`Outcome == "proposed"`, `TargetState == "in-progress"`, `Action == "propose_status_in_progress"`). Confirms Key finding 2's claim was correct — no new FSM work was needed or done here, consistent with the scaffold's prediction.

**"Resume, not just label transition" — the substantive half of AC3.** I read `cds-dispatch/SKILL.md`'s new paragraph (after claim step 9) and the new `delta/SKILL.md` §9.11 in full, independently of self-coherence.md's summary. §9.11 names a concrete mechanism, not vague prose:
- **Trigger** is named precisely (scanner-driven `in-progress → todo` reconciliation via the new rule, not an operator action).
- **Detection** is named precisely — no new Go surface, reuses `FactSnapshot.BranchExists`/`PRExists`/`CDDArtifacts` already exposed by `cn issues fsm evaluate --issue {N} --json`, plus a textual "MECHANICAL reversion" phrase check on the issue comment.
- **Routing sequence** is a concrete 3-step numbered procedure: (1) δ reads `.cdd/unreleased/{N}/` on the branch before dispatching any role; (2) does NOT re-dispatch γ for a fresh scaffold if `gamma-scaffold.md` already exists; (3) dispatches whichever role is next given what artifacts already exist.
- A **comparison table** against §9.10 (cause / evidence / what "resume" means / repair-machinery applicability) makes the two shapes distinguishable, not conflated.
- `cds-dispatch/SKILL.md`'s Repair re-entry preflight Step A is reordered so the cnos#630 check runs **before** the repair-re-entry classification — this is a real, load-bearing change (not just added prose): pre-diff, a resumed cell's pre-existing branch-with-commits would have satisfied the old Step A's repair-re-entry bullet list and misrouted to Steps B–E looking for a rejection that never happened. I confirmed this misrouting risk is genuine by reading the pre-diff Step A text (`git diff`'s `-` side) — the old bullet list's second bullet (`the cycle/{N} branch already exists with commits beyond its base`) is unconditional and would indeed have matched a resumed-from-matter cell.

This clears the bar the scaffold set ("the label transition alone does NOT satisfy AC3") — the mechanism is named with a trigger, a detection method, and a routing sequence, not merely asserted. **Pass.**

**Residual honestly-disclosed debt (not blocking):** self-coherence.md's own §Debt item 1 states §9.11 has zero empirical witnesses from an actual wake-invoked-δ firing (this cycle itself ran under bootstrap-δ). This is correctly disclosed as debt, not silently glossed — I do not treat unverified-in-production doctrine as an AC3 failure, since AC3's own oracle explicitly allows a doctrine-text oracle when no Go-testable surface exists, and δ's own dispatch prompt anticipated this class of gap.

### AC4 — audit-note comment on every mechanical reversion; not misread as an orphan

**Oracle:** every new mechanical label reversion posts an audit-note comment; a subsequent claim/scan pass doesn't misread it as an unexplained orphan.

**Comment-posting mechanism, independently traced:** read `scan.go`'s diff directly — the new action flows through the pre-existing generic `case dec.Outcome == "proposed" && dec.TargetState != ""` branch (line ~283 in the current file), which already builds `body := fmt.Sprintf("cds recovery scan: reconciled #%d -- %s.\n\n%s", issue, dec.EnabledTransition, dec.Reason)` and calls `opts.PostComment` — this is the exact, unmodified Cases-2/5 mechanism; zero new comment-posting code was added. Confirmed via `git diff -- scan.go | grep '^+func'` → no output, i.e. no new function was introduced anywhere in `scan.go`.

**Content lock, independently re-run:** `TestAC630_AuditNoteReasonNamesMechanicalReversion` — re-ran myself: **PASS**, confirms `dec.Reason` contains all three of `"MECHANICAL reversion"`, `"PRESERVED"`, `"cnos#368"`.

**"Not misread as orphan" — checked from both directions.** (a) Go-level: `TestScan_Idempotent`'s second pass (re-run myself: **PASS**) confirms a subsequent scan tick against the post-reconciliation state produces zero new comments/finalize-calls/labels — the reconciled cell drops out of the active-scan set the same way `deadNoMatter` already did, so `scan` itself never re-flags it. (b) Claim-path level (the actual "orphan" ambiguity AC4 cares about, since a *human/wake* reading the issue at claim time is the one who could misclassify it, not `scan` itself): `cds-dispatch/SKILL.md`'s reordered Step A explicitly checks for the "MECHANICAL reversion" phrase before repair-re-entry classification (see AC3 evidence above) — this is the mechanism that keeps a claiming wake from reading the audit-noted state as an unexplained orphan. **Pass.**

### AC5 — no regression, full CI green

I did not trust α's pasted output — I re-ran every command myself from a clean shell.

```
$ cd src/packages/cnos.issues/commands/issues-fsm && go build ./... && go vet ./...
(no output, exit 0)

$ go test ./... -v 2>&1 | grep -c '^--- PASS'
77
$ go test ./... -v 2>&1 | grep -c '^--- FAIL'
0
$ go test ./... -v 2>&1 | grep TestSeam_CellKindNotEnforced
=== RUN   TestSeam_CellKindNotEnforced
--- PASS: TestSeam_CellKindNotEnforced (0.00s)
```

77 pass / 0 fail — matches self-coherence.md's claimed count exactly, independently reproduced. `TestSeam_CellKindNotEnforced` located via `grep -rn TestSeam_CellKindNotEnforced src/` (`issuesfsm_test.go:810`) and confirmed passing.

Only two pre-existing tests changed assertions, both named explicitly by α and independently verified as justified, not silent: `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` (renamed+repurposed, see AC2) and `TestScan_Idempotent` (assertions updated in place for the `checkpointed` sub-case only — I diffed this test and confirmed the `deadNoMatter`/`blocked` sub-cases are untouched).

```
$ bash scripts/ci/check-dispatch-repair-preflight.sh; echo $?
cnos#516 repair-preflight guard: ... (exit 0)
0
$ bash scripts/ci/check-dispatch-closeout-integrity.sh; echo $?
cnos#524 closeout-integrity guard: ... (exit 0)
0

$ cd src/go && go build ./... && go vet ./... && go test ./...
ok for all 15 internal packages; cmd/cn [no test files]
```

**Live GitHub Actions CI, independently pulled (α disclosed this as "not independently run" from the sandbox — I confirmed it from outside the sandbox via `gh api`):**

```
$ gh api repos/usurobor/cnos/commits/a63bc9d2331975f8f9e5467b5bb54c9199b89d8d/check-runs \
    --jq '.check_runs[] | "\(.name)\t\(.status)\t\(.conclusion)"'
Package verification                          completed  success
Binary verification                           completed  success
CDD artifact ledger validation (I6)           completed  success
Repo link validation (I4)                     completed  success
Protocol contract schema sync (I2)            completed  success
Dispatch closeout-integrity guard (cnos#524)  completed  success
Go build & test                               completed  success
Dispatch repair-preflight guard (cnos#516)    completed  success
SKILL.md frontmatter validation (I5)          completed  success
Package/source drift (I1)                     completed  success
Re-render + diff per-package goldens          completed  success
```

All 11 checks on the review head SHA are `success`. This resolves α's own self-disclosed gap (self-coherence.md §Review-readiness: "Branch CI is not independently confirmed green on GitHub Actions for this head commit") — it is in fact green, confirmed by β from outside the sandbox. **Pass.**

### AC6 — idempotence

**Oracle:** run the new reconciliation path twice against the same fixture state; zero new mutations on the second pass.

Re-ran `TestScan_Idempotent` myself (already covered under AC5's PASS count) and read its body directly: first pass against `{deadNoMatter, checkpointed, blocked}` produces exactly 2 comments (one per requeue, including the new with-matter requeue) and 0 finalize calls; second pass (simulating the post-reconciliation world by dropping the two reconciled issues from the active-issue listing — the same convention the test already used pre-diff for `deadNoMatter`) asserts `com.count()` unchanged, `fin.callCount()` unchanged, and `secondLabelCount == 0`. This is a genuine two-pass mutation-count assertion, not a restated narrative. **Pass.**

## Guardrail verification (independent of the AC table)

- **No branch/PR deletion anywhere in the code diff.** `git diff origin/main...cycle/630 | grep -nE "branch -D|pr close|DeleteRef|deleteBranch|closePR"` → the only hits are inside `beta-review.md`'s/`self-coherence.md`'s own prose describing the guardrail check itself (i.e., documentation text, not code). Zero hits in any `.go`/`.json`/`.yml` diff hunk. **Confirmed.**
- **`applyStatusLabel`/`PostComment`/`cn cell finalize` reused, not forked.** `git diff -- scan.go | grep '^+func'` → empty (zero new functions in `scan.go`). `grep -n "func applyStatusLabel\|func.*PostComment\|ghAddLabel\|ghRemoveLabel"` across the package shows the same single definition sites as on `main` (`fetch.go`); no parallel/duplicate label-writer or comment-poster was added. **Confirmed.**
- **No new `status:*` label.** `diff <(git show origin/main:src/packages/cnos.core/labels.json) <(git show cycle/630:src/packages/cnos.core/labels.json)` → byte-identical, no diff. 7 pre-existing `status:*` entries unchanged. **Confirmed.**
- **Two open-design-call decisions documented with rationale, not just asserted.** Read self-coherence.md's "Open design calls" subsection directly: Design call 1 (distinct new action vs. in-place edit) is justified against the actual guard vocabulary (`pr_exists` as the literal boolean that distinguishes "checkpointed" from "not yet checkpointed," and the concrete alternative-cost argument against collapsing into one rule — reintroducing state-name-shaped Go branching `table.go`'s own doc comment forbids). Design call 2 (doctrine-only vs. new Go surface) is justified with a concrete negative check ("no `todo`-state rule references `branch_exists`/`pr_exists`/`cdd_artifacts` at all, so there is nothing for the claim path to miss"), independently confirmed above under AC3. Both are substantive engineering rationale, not restated scaffold text. **Confirmed.**
- **`.github/workflows/cnos-cds-dispatch.yml` was rendered, not hand-edited — independently reproduced, not just sha256-matched.** I rebuilt the kernel binary (`cd src/go && go build -o /tmp/cn ./cmd/cn`) and confirmed `/tmp/cn` does **not** expose `install-wake` as a kernel subcommand (α's self-coherence.md correctly self-corrected on this point rather than asserting a false invocation). I then ran the actual renderer directly (`bash src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch`) against the checked-out `cycle/630` tree and confirmed: renderer output is `(unchanged)`, `diff` against a pre-render copy of `.github/workflows/cnos-cds-dispatch.yml` is empty, and `git status --short` shows zero working-tree changes to either the live workflow or the golden file. This is a stronger check than the sha256-comparison self-coherence.md performed — I independently reproduced the render from source, not merely compared two already-committed artifacts. **Confirmed: rendered, reproducible, never hand-edited.**

## Implementation-contract conformance (β Rule 7)

Checked the diff against all 7 pinned axes in the gamma-scaffold's `## Implementation contract` table:
1. **Language = Go** — all new logic is Go (`table.go` doc-comment, `scan.go` comments) plus declarative JSON (`transitions.json`) plus prose (`cds-dispatch/SKILL.md`, `delta/SKILL.md`). Confirmed.
2. **CLI integration = `cn issues fsm evaluate`/`scan`, no new binary** — confirmed; zero new `cmd/` entries, zero new binaries in the diff.
3. **Package scoping** — every changed file maps onto the pinned paths (`issues-fsm/`, `fsm/transitions.json`, `cds-dispatch/SKILL.md`, `delta/SKILL.md` §9). Confirmed via the file list in `git diff --stat`.
4. **Existing-binary disposition preserved** — confirmed, no binary replaced/forked.
5. **No new runtime dependencies** — confirmed, `go.mod`/`go.sum` untouched (not in the diff stat at all).
6. **JSON/wire contract additive-only** — new action string `propose_status_todo_with_matter` added; no existing `FactSnapshot`/`Decision`/`ScanReport` field renamed or removed (confirmed by reading `table.go`'s diff, which is comment-only).
7. **Backward-compat invariant** — all pre-existing tests pass except the two explicitly-named, justified redirects (AC5 evidence above); `todo` state rules byte-identical (AC3 evidence above).

All 7 axes conform. **Pass.**

## Findings

None blocking.

One non-blocking observation (already self-disclosed by α in §Debt, independently re-confirmed by me as accurate and adequately scoped): the §9.11 doctrine addition has no live wake-invoked-δ empirical witness yet, mirroring §9.10's own first-use history (cycle/497 was also a bootstrap-exception witness). This is a legitimate scope boundary, not a gap this cycle should have closed — a live-firing witness requires an actual production wake event this cycle cannot manufacture. No action requested.

One additional observation, non-blocking: self-coherence.md's §Review-readiness section states "Branch CI is not independently confirmed green on GitHub Actions for this head commit" as a known limitation of α's sandboxed environment. I was able to independently confirm all 11 checks green via `gh api` from outside the sandbox (see AC5 above) — this closes α's self-disclosed gap; noting it here for γ's PRA record rather than as a finding requiring action.

## Verdict

**converge**

All 6 ACs pass with independently-reproduced evidence — every test I cite above was re-run by me from a clean shell, not read off α's pasted output; the rule ordering, the redirected-test rationale, the doctrine-mechanism substance (§9.11), and the workflow-render reproducibility were each re-derived from source, not accepted from self-coherence.md's narrative. All 4 guardrails hold (no branch/PR deletion, no forked label/comment mechanism, no new `status:*` label, both open-design-calls substantively justified). Implementation-contract conformance (Rule 7) confirmed across all 7 pinned axes. Live GitHub Actions CI is green on the review head SHA (11/11 checks `success`), resolving α's own self-disclosed CI-visibility gap. No blocking findings. Ready for γ closeout / merge.
