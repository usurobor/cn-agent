# self-coherence — cycle/574

manifest: sections planned = [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
completed: [Gap, Skills, ACs, Self-check]

## §Gap

**Issue:** [cnos#574](https://github.com/usurobor/cnos/issues/574) — "fix(cds/fsm): harden review guards + remote-branch observation; correct wave-closure language (post-#567 remediation)"

**Mode:** design-and-build (per issue header). Bug-fix hardening cell, not greenfield design — AC2/AC3's target guard strings were fully specified by the issue's own "Status truth" table and γ's scaffold; AC1/AC5/AC6 are verification/correction tasks; AC4 is the one open design call (fetch-vs-API), explicitly reserved for α.

**Parent/wave:** #567 (master, CLOSED/COMPLETED) — this cell is an operator-review remediation of the #567 wave (#568 Phase 1, #570 cell-kind doctrine, #569 Phase 2), all three merged and closed, found not fully clean by operator review.

**Branch:** `cycle/574`, created by γ from `origin/main@452191fe28bae8f7fbad53fe2010d4c122645342`. Confirmed unchanged at α's dispatch time (`git rev-parse origin/main` still resolves to `452191fe...` as of this writing) — no rebase required (pre-review gate row 1).

**γ scaffold read in full** at `.cdd/unreleased/574/gamma-scaffold.md` before any implementation commit, per dispatch instructions. Followed its load order, exact guard strings, per-AC oracle list, source-of-truth table, scope guardrails, and friction notes. Did not improvise beyond it except where explicitly documented below (AC4's design call, and one internal scaffold tension in AC2 resolved and documented in §ACs).

**Underlying issue body** read in full via `gh issue view 574 --repo usurobor/cnos --json title,body,state,comments` (single comment present, an operator clarification unrelated to AC scope — no mid-flight clarification needed).

## §Skills

**Tier 1 (loaded per α's own load order):**
- `CDD.md` — canonical lifecycle/role contract (loaded as the framing document before role work).
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — this file's own role surface, followed section-by-section: §2.1 dispatch intake, §2.5 self-coherence (incremental-write discipline followed here), §2.6 pre-review gate, §2.7 request review.
- `issue/SKILL.md` (referenced) — for interpreting the AC boundary questions that arose in AC2 (see §ACs AC2's documented tension).

**Tier 2 (eng/*, always-applicable):** `eng/go` conventions implicitly applied (idiomatic Go: doc comments on every new exported/unexported function, table-driven test style matching the existing file's idiom, `go vet`/`go test -race` as the gate). `eng/ship` bug-fix TDD explicitly followed for AC2/AC3 (failing fixture first, confirmed red, then tighten until green — see §ACs).

**Tier 3 (issue-specific):**
- No new language/package-scoping skill needed — the issue is scoped entirely within already-loaded Go (`cnos.issues/commands/issues-fsm`) and CDS-owned JSON data (`cnos.cds/skills/cds/fsm/transitions.json`), per the γ scaffold's pinned implementation contract (§3.6 in alpha/SKILL.md: the 7 axes are δ's, not α's — α did not relax any of them; see §ACs AC7 / implementation-contract row below).
- GitHub CLI (`gh`) operational knowledge for AC1 (authenticated issue-state re-verification), AC5 (comment posting + issue filing + label discovery via `gh label list`).

No skill gap identified that would have prevented remaining debt (see §Debt).

## §ACs

### AC1 — authenticated wave-issue-state verification + discrepancy record

**Status: MET.**

Re-verified (2026-07-04, this session) via `gh issue view {N} --repo usurobor/cnos --json state,stateReason,labels` for all four:

| Issue | state | stateReason | labels |
|---|---|---|---|
| #567 | CLOSED | COMPLETED | P1, area/cdd, area/cds, area/wake, kind/tracking, area/issues |
| #568 | CLOSED | COMPLETED | P1, kind/tooling, area/cds, area/issues |
| #569 | CLOSED | COMPLETED | P1, area/cdd, kind/process, area/cds, area/wake |
| #570 | CLOSED | COMPLETED | P1, area/cdd, kind/process, area/cds, area/coherence, area/wake, area/issues |

Matches γ's pre-gathered values in gamma-scaffold.md §Per-AC oracle list verbatim. No `dispatch:cell`/`protocol:cds`/`status:*` label present on any of the four.

**Public-render discrepancy check:** fetched the unauthenticated public HTML for all four issue pages (`curl -sSL https://github.com/usurobor/cnos/issues/{N}`) and inspected the definitive header-state badge (`data-testid="header-state" ... data-status="issueClosed"`), which is GitHub's own rendered status indicator, not a heuristic. All four showed `data-status="issueClosed"` — consistent with the authenticated API state. **No discrepancy found.**

Method note: an initial blunt `grep -io octicon-issue-opened` over the raw HTML produced false-positive matches (these came from embedded issue-reference chips for *other*, currently-open issues linked within each issue's own markdown body — e.g. #567's body links #566/#532/#516 — not from the primary issue's own status badge). The `data-testid="header-state"` selector is the correct, unambiguous signal and was used for the final determination.

### AC2 — `review`-state validity requires deliverable matter, not review-request alone

**Status: MET.**

TDD sequence followed exactly as the issue/scaffold specify:
1. Wrote `testdata/review-request-only.json` (status:review, `review_request_present:true`, `pr_exists:false`, `pr_commit_count:0`) and `testdata/review-partial-evidence.json` (`pr_exists:true`, `pr_commit_count:0` i.e. `pr_has_commits:false`, `review_request_present:true`) plus corresponding tests `TestAC574_ReviewRequestAloneNoLongerValid` / `TestAC574_ReviewPartialEvidenceBlocked` in `issuesfsm_test.go`.
2. Ran `go test -run TestAC574 -v` against the **unmodified** `transitions.json` — both new tests **FAILED** as expected (outcome was `valid`, wanted `blocked`) — captured in this session's transcript. Positive/backward-compat tests (`TestAC574_ReviewWithPRStillValid` etc.) already passed pre-tightening, confirming the fixtures were well-formed.
3. Tightened `transitions.json`'s `review` state: rule 1 (valid) changed from `any_true: [pr_exists, branch_has_commits, review_request_present]` to `all_true: [review_request_present, pr_exists, pr_has_commits]`; rule 2 (blocked) changed from a symmetric `all_false` over the old 3-guard set to an **unconditional** rule (no guard clause at all) positioned second, so it fires on ANY combination rule 1 didn't match — not just all-three-false. This resolves the Friction-notes concern directly: a naive `all_false` rewrite would let a partial-evidence state (e.g. `pr_exists:true, pr_has_commits:false`) fall through both explicit rules to the engine's generic "no rule in the transition table matched" fallback (table.go's `Evaluate`, which produces no `MissingEvidence` list). The unconditional-second-rule design instead gives an explicit, diagnostic `blocked` outcome with a populated `evidence_guards`/`MissingEvidence` list for every partial-evidence combination. This was a deliberate choice, not an accident — documented here per the Friction notes' explicit request.
4. Re-ran the full package test suite — all green, including the two new tests and all pre-existing ones (with one pre-existing test's assertion updated, see below).

**Internal scaffold tension found and resolved (α's documented judgment call):** the γ scaffold's literal instruction ("replace `any_true: [...]` with `all_true: [review_request_present, pr_exists, pr_has_commits]`") is the exact 3-guard AND the issue body's own Status-truth table also states. But the scaffold's own β-verification bullet for AC2 separately says "confirm the existing `review-with-pr.json` fixture (PR + commits) still evaluates `valid` (no regression)" — and that pre-existing fixture had `review_request_present:false`. Applying the literal 3-guard `all_true` as-is would make that pre-existing fixture evaluate `blocked`, directly contradicting the "no regression" instruction in the same document.

Resolution: treated this the same way the Friction notes explicitly sanction for the AC3 in-progress fixtures ("existing fixtures may need their expected outcomes updated, not just re-verified... this is an intentional behavior change, not a regression to avoid" — generalized here to "existing fixtures may need their *field values* updated to remain a meaningful full-evidence positive case under the tightened invariant"). Updated `testdata/review-with-pr.json` to also carry `"review_request_present": true` (and added `REVIEW-REQUEST.yml` to its `cdd_artifacts` list) — this is realistic, not a fudge: the fixture models "an in-progress review that already has a PR + commits," and nothing in the system deletes `REVIEW-REQUEST.yml` once a review is underway, so a full-evidence review-state fixture plausibly does carry all three facts. Kept the literal δ-pinned 3-guard `all_true` exactly as specified rather than weakening it back to a 2-guard `all_true: [pr_exists, pr_has_commits]`, since the issue body's own Status-truth table is the more authoritative (δ-pinned) source for the invariant text, and weakening it would silently under-implement the stated invariant.

Consequential test update: `TestAC3_EmptyReviewBlocked`'s `wantMissing` set changed from `{pr_exists, branch_has_commits, review_request_present}` to `{pr_exists, pr_has_commits, review_request_present}` — `branch_has_commits` dropped out of the review-state guard set entirely (replaced by `pr_has_commits`), so it can no longer appear in `MissingEvidence`. `testdata/review-empty.json` (unchanged) still evaluates all three of the new guard set as false, so the test's assertion of exactly-3-missing-items still holds, just with the corrected member set. Documented inline in the test file at the point of change.

**Evidence:** `TestAC574_ReviewRequestAloneNoLongerValid`, `TestAC574_ReviewPartialEvidenceBlocked` (new, positive), `TestAC574_ReviewWithPRStillValid`, `TestAC3_ReviewWithEvidenceIsValid`, `TestAC3_EmptyReviewBlocked` (backward-compat, all green post-tightening) in `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go`.

### AC3 — `in-progress → review` requires a PR

**Status: MET.**

Same TDD discipline: wrote `testdata/in-progress-review-request-branch-only.json` (`review_request_present:true`, `branch_has_commits` true via `commits_beyond_base:4`, `pr_exists:false`) and `TestAC574_InProgressBranchOnlyNoLongerProposesReview`, confirmed it **FAILED** against the unmodified table (outcome was `proposed`/`review`), then tightened `transitions.json`'s `in-progress` rule 1 from `all_true:[review_request_present]` + `any_true:[pr_exists, pr_has_commits, branch_has_commits]` to a single `all_true:[review_request_present, pr_exists, pr_has_commits]`; rule 2 (blocked) simplified from `all_true:[review_request_present]` + `all_false:[pr_exists, pr_has_commits, branch_has_commits]` to `all_true:[review_request_present]` alone (unconditional given that gate) with `evidence_guards:[pr_exists, pr_has_commits]` — the same "unconditional second rule catches every partial-evidence case" design as AC2, and for the identical reason (rule 2's own `all_true:[review_request_present]` guard means rule 2 only ever runs when review_request_present is already true, and rule 1 already tested the full 3-guard AND and didn't match, so rule 2 firing here is logically exactly the "at least one of pr_exists/pr_has_commits is false" case — no separate re-statement of that logic is needed in JSON).

Re-checked the two existing #573-era fixtures per the Friction notes' explicit instruction to verify field values, not assume: `in-progress-review-request-with-matter.json` sets `pr_exists:true, pr_commit_count:5` (so `pr_has_commits:true`) — satisfies the tightened rule unchanged, no update needed. `in-progress-review-request-no-matter.json` sets `pr_exists:false, pr_commit_count:0, branch_exists:false` — still blocks unchanged. Neither required a fixture or expected-outcome edit.

**`run_active` non-gating preserved** exactly per constraint 5 (scope guardrails) and the `_doc_phase2` rationale: rules 1/2 (both gated only on `review_request_present` and the matter guards, never on `run_active`) are unchanged in position (still ahead of rule 3's `run_active` check) and unchanged in their non-gating on `run_active`. Added `TestAC574_InProgressRunActiveNonGatingPreserved`, which asserts the with-matter fixture (which has `run_state:"in_progress"`) still proposes `review` — proving a worker's own synchronous in-flight request still works. Rules 3–5 (the pre-existing dead-run reconciliation rules) are byte-unchanged.

**Evidence:** `TestAC574_InProgressBranchOnlyNoLongerProposesReview`, `TestAC574_InProgressWithMatterStillProposesReview`, `TestAC574_InProgressRunActiveNonGatingPreserved` (new), `TestAC569_InProgressReviewRequestWithMatterProposesReview`, `TestAC569_InProgressReviewRequestNoMatterBlocked`, `TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid` (all pre-existing, all still green — rules 3–5 unaffected).

### AC4 — matter/branch observation robust to remote-only branches

**Status: MET. Design decision documented below (α's call, per dispatch instructions).**

**Decision: option (b) — GitHub-API branch/commit observation in `fetch.go`**, following γ's non-binding recommendation, over option (a) (workflow-level `git fetch refs/heads/cycle/*`).

**Reasoning:** (1) Option (b) fixes the observation primitive (`assembleLive`) itself, so it benefits every caller — the CI workflow, a local operator run, or any future non-workflow invocation — not just the one CI-workflow invocation path option (a) would patch. (2) `fetch.go` already has the exact idiom to extend: `ghGetJSON` is used for PR/run observation a few lines below the block this touches, so adding a branch-existence GET (`/repos/{repo}/branches/{branch}`, 200 vs 404) and a compare GET (`/repos/{repo}/compare/main...{branch}`'s `.ahead_by`) is a natural, small extension of an existing pattern — no third-party GitHub client, no new runtime dependency (net/http stdlib only, per the pinned implementation contract). (3) Option (a) would also require a corresponding ref-namespace change in `fetch.go`'s `exec.Command` calls (checking `refs/remotes/origin/cycle/{N}` instead of `refs/heads/cycle/{N}`) — a two-surface change (workflow YAML + Go) for a fix that still only covers the CI invocation path. Given the issue's own cycle-scope-sizing "Decision" explicitly permits (b) even though it changes Go source, and does not name diff-size as having ballooned, no AC4 sub-issue split was triggered.

**Implementation:** `assembleLive`'s existing local-git block is unchanged (still tries local git first). A new block runs only when `repo != "" && !snap.BranchExists` (i.e., local git found nothing) and calls the new `observeRemoteBranch(ctx, repo, branch, token)` function, which performs exactly the two GETs described above and returns `(exists bool, commitsBeyondBase int)`. A GitHub API failure (404, network hiccup) is tolerated, not a hard error — mirrors the rest of `assembleLive`'s "kept honest but deliberately simple" live-path philosophy (per the file's existing package doc comment).

**Testability design note:** `observeRemoteBranch` deliberately reads the `githubAPIBase` package var (already used by the Phase 2 mutation-path tests' `withFakeGitHub` helper) rather than a hardcoded `"https://api.github.com"` literal, unlike the pre-existing labels/PR/run-observation GETs in the same file (which intentionally keep their Phase-1 literal URLs — see `githubAPIBase`'s doc comment). Since this is new code (not part of the Phase-1 byte-identical-URL guarantee), routing it through `githubAPIBase` was a safe, in-scope choice that makes it unit-testable against a fake HTTP server with zero live-network dependency, matching γ's oracle text verbatim ("a mock/fake HTTP server ... using `githubAPIBase` override"). Extracted as its own function (rather than left inline in `assembleLive`) specifically so tests could exercise it in isolation without also triggering `assembleLive`'s other real-network calls.

**Proof (fully hermetic — no live network in any test):**
- `TestObserveRemoteBranch_ExistsWithCommits` — fake server reports branch exists (200) + `ahead_by:3`; confirms `observeRemoteBranch` returns `(true, 3)`.
- `TestObserveRemoteBranch_AbsentReportsFalse` — fake server 404s; confirms `(false, 0)`, not an error.
- `TestAssembleLive_RemoteOnlyBranchResolvesToDeltaRecovery` — the AC4 oracle's exact scenario: composes `observeRemoteBranch`'s result (simulating "local git found nothing, but the branch exists remotely with 3 commits") into a `FactSnapshot` with `status:in-progress` and no active run (`RunState:""`), runs it through the **real** transition table via `Evaluate`, and asserts `outcome=="proposed"`, `action=="propose_delta_recovery"`, and explicitly asserts `TargetState != "todo"` — proving the full pipeline (observation → table) lands on delta-recovery, never blind requeue (the cnos#368 failure mode).
- `TestAssembleLive_RemoteBranchFallbackNeverRunsWithoutRepo` — confirms the fallback makes zero GitHub API calls when `repo==""` (the existing `assembleLive(ctx, "", issue, "")` calling convention used by the pre-existing cell-kind test), locking "only fills the gap, never runs without a repo to query."

**Scope guardrail honored:** no `cell_kind`-based branching added anywhere in this change (constraint 2) — `observeRemoteBranch` and its call site touch only `BranchExists`/`CommitsBeyondBase`.

### AC5 — wave-closure language corrected + Phase 3 filed (deferred)

**Status: MET.**

Re-ran `rg -i "FSM owns" src/ docs/ .cdd/ .github/` before touching any doctrine prose, per the Friction notes' explicit instruction not to assume a repo-file fix is needed. Result: **zero hits** in any shipped file (`src/`, `docs/`, `.github/`) — the only hits are inside `.cdd/unreleased/574/gamma-scaffold.md` itself (γ's own scaffold quoting the overclaim string for reference), which is not a shipped doctrine surface and is expected to contain that quote. No repo-file edit was made, matching γ's own re-grep finding — the doctrine files (`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md`) were already correctly scoped from #569's landing.

The actual overclaim — #567's own closing comment's bolded headline, `"✅ FSM wave complete — \`Workers produce matter. The FSM owns status labels.\`"` — was corrected via an **appended correction comment**, not an edit of the original: https://github.com/usurobor/cnos/issues/567#issuecomment-4880282434. The correction states the true scoped boundary ("FSM gates the `in-progress → review` transition; claim/hard-block/release-back-to-queue remain direct wake writes") and references cnos#574.

**Edit-vs-append judgment call:** `gh auth status` / `gh api user` showed the authenticated identity for this session's `gh` calls is `usurobor` — the same account that authored the original #567 comment — so a direct edit was technically possible (unlike the scaffold's assumption that α's identity would differ from the comment author's). α chose to **append rather than edit** anyway: editing a historical closure comment in place would silently rewrite the wave's own record with no visible trail that a correction occurred, which conflicts with CDD's evidence-based, explicit-correction ethos (the same principle the loaded skill applies to self-coherence itself — corrections are appended, not silently absorbed). An appended, dated, cross-referenced correction preserves the full audit trail for any future reader of the #567 thread.

Filed the Phase 3 tracking issue: https://github.com/usurobor/cnos/issues/575 ("cds/fsm Phase 3 (deferred): route claim/hard-block/release-back-to-queue through the FSM"), labeled `area/cds`, `kind/process`, `priority/deferred` (checked `gh label list --repo usurobor/cnos` first — `priority/deferred` is the repo's existing deferred-label convention, description "Deferred — not required for current implementation scope"; no separate `held`-type label exists). Linked from #567 via a second comment: https://github.com/usurobor/cnos/issues/567#issuecomment-4880283288. #575's body itself also states Phase 3 is explicitly **not implemented** here, names the three transitions in scope for a future cell (claim, hard-block, release-back-to-queue), and links back to #567 and cnos#574.

**Scope guardrail honored:** Phase 3 itself was not implemented (constraint 1) — no FSM-routed logic for claim/hard-block/requeue was added to `transitions.json` or any Go source in this cycle.

### AC6 — hidden/bidi Unicode sweep of wave-touched files

**Status: MET — sweep clean.**

Ran a byte-level sweep for bidirectional-control and hidden-format Unicode characters (`rg -nP '[\x{200B}-\x{200F}\x{2028}-\x{202E}\x{2060}-\x{2069}\x{FEFF}]'`) over the full wave-touched file set named in gamma-scaffold.md §Surfaces (the union of #571/#572/#573's changed files): all three named workflow files, `go.work`, `src/go/cmd/cn/main.go`, `src/go/go.mod`, `src/go/internal/cli/cmd_issues_fsm.go`, `transitions.json`, both `cds-dispatch` doctrine files, `delta/SKILL.md`, `CDD.md`, `CELL-KINDS.md`, `cdd/issue/SKILL.md`, `cnos.issues/SKILL.md`, all `issues-fsm/*.go` and `issues-fsm/testdata/*.json`, `GLOSSARY.md`, and all six `.cdd/unreleased/{568,569,570}/*.md` closeout artifacts.

**Result: zero matches** (`rg` exit code 1 = no matches found across the entire file set). Verified every named path actually exists before trusting the "no matches" result (a missing file would silently contribute zero without being a real "clean" signal) — all files present, none skipped.

**#572 checked specifically first** (the PR GitHub flagged with a hidden/bidi warning): fetched its file list via `gh pr view 572 --json files` — all 12 files (`.cdd/unreleased/570/*.md` ×6, `GLOSSARY.md`, `CDD.md`, `CELL-KINDS.md`, `cdd/issue/SKILL.md`, `fetch.go`, `issuesfsm_test.go`) are already members of the swept union above and returned zero matches. No unexplained bidi/hidden-Unicode content found anywhere in the wave-touched set; the GitHub UI warning did not correspond to a persisting hidden-character issue at the current state of these files (it may have flagged content in an intermediate diff hunk that was since resolved, or a false-positive on legitimate non-ASCII prose — this sweep only speaks to the current file contents, which are clean).

### AC7 — current gates remain green

**Status: MET (locally verified; CI on the pushed head commit not directly observable from within this session — see §Debt).**

`TestSeam_CellKindNotEnforced` (`issuesfsm_test.go`) confirmed **byte-identical** to `origin/main` via an automated function-body extraction diff (not just "still passes") — see evidence below. Not touched by any edit in this cycle.

Full workspace build + test, run per go.work module (root `go build ./...`/`go test ./...` do not work directly since the repo root itself is not a workspace module — each `use` entry was built/tested individually):

```
=== src/go ===                                                    all packages ok (14 test packages, 1 no-test-files)
=== src/packages/cnos.cdd/commands/cdd-verify ===                 ok
=== src/packages/cnos.issues/commands/issues-map ===               ok
=== src/packages/cnos.issues/commands/issues-fsm ===               ok (go test -race also green)
```

`issues-fsm` package test runner output (this is the package every AC2/AC3/AC4 change touches): `go test ./... -v` → 48 `=== RUN` lines (including subtests), **37 top-level `--- PASS` entries, 0 `--- FAIL`**, final `ok` line. `go test ./... -race` → green, no data races detected (relevant given the new `observeRemoteBranch`/`withFakeGitHub` tests use goroutine-backed `httptest` servers).

`TestSeam_CellKindNotEnforced` byte-identity check: extracted the function body from `origin/main`'s `issuesfsm_test.go` and from this branch's HEAD via a brace-balanced Python extraction (not a line-range diff, which could miss whitespace-only drift) — result: **IDENTICAL**.

Named CI jobs (I1/I2/I4/I5/I6, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity, go build/test, Package, Binary) were not independently re-run inside this session (no local CI harness invocation was performed beyond `go build`/`go test`/`go vet`) — this branch's actual GitHub Actions run status on the pushed head commit is deferred to β per the pre-review gate's explicit allowance ("if local CI is unavailable, the artifact's review-readiness section says so explicitly and β waits for green before merge" — alpha/SKILL.md §2.6 row 10). Declared as known debt in §Debt.

**Scope guardrails honored (full re-check across all seven):**
1. Phase 3 not implemented — only filed (issue #575).
2. `cell_kind` stays observed-only — `grep -n "CellKind" src/packages/cnos.cds/skills/cds/fsm/transitions.json` returns nothing; no rule in the diff references `FactSnapshot.CellKind`; `TestSeam_CellKindNotEnforced` byte-identical (above).
3. No new status labels / taxonomy change — the diff touches only existing guard combinators (`any_true`→`all_true`) and `reason`/`evidence_guards` strings in `transitions.json`; no `status:*`/`protocol:*` label string appears anywhere in the diff.
4. No external controller service, no wake-source-model change, no #216 package-command-discovery change, no Demo 0 — none touched.
5. `run_active` non-gating preserved on `in-progress` rules 1/2 (see §ACs AC3) — explicitly tested via `TestAC574_InProgressRunActiveNonGatingPreserved`.
6. No Phase-1/Phase-2 evaluator logic re-derived — `table.go`'s `Evaluate`/`ruleMatches`/`guardFuncs` are byte-unchanged; AC2/AC3 are data-only edits to `transitions.json`; AC4 added one new Go function (`observeRemoteBranch`) plus its call site, not a rewrite of existing evaluator logic.
7. STOP-and-ask trigger — AC4's chosen approach (option (b)) uses only the GitHub REST endpoints already implicitly authorized by the existing `GITHUB_TOKEN`/`ghGetJSON` pattern already used for PR/run observation in the same file (`/repos/{repo}/branches/{branch}` and `/repos/{repo}/compare/...` are both standard-scope read endpoints, no broader token authority than the pre-existing PR/issue/run GETs) — no broadening of GitHub API/token authority; no STOP triggered.

**Implementation-contract conformance (the 7 axes pinned in gamma-scaffold.md, alpha/SKILL.md §3.6):** Language — Go + JSON data, as pinned (no Python/other-language introduced). CLI integration target — `cn issues fsm evaluate [--apply]` unchanged; no new subcommand, no flag rename. Package scoping — all Go changes stayed within `src/packages/cnos.issues/commands/issues-fsm/`; the only data change is `src/packages/cnos.cds/skills/cds/fsm/transitions.json`; no workflow YAML was touched (AC4 picked option (b), not (a)/hybrid, so the workflow-file row of the pinned scoping table was correctly left untouched). Existing-binary disposition — preserved, single compiled-in `cn issues fsm evaluate [--apply]`. Runtime dependencies — none new; `observeRemoteBranch` reuses `net/http` stdlib + the existing `ghGetJSON` idiom. JSON/wire contract preservation — `FactSnapshot`/`Decision` struct shapes unchanged; all fixture files remain loadable with unchanged field names/types (only the *values* of two testdata JSON files changed: `review-with-pr.json`'s `review_request_present`/`cdd_artifacts`, and the three new fixture files added). Backward-compat invariant — every existing fixture whose scenario the invariant-tightening doesn't specifically probe evaluates identically (verified by the full green test suite); the two exceptions (`review-with-pr.json`'s field values, `TestAC3_EmptyReviewBlocked`'s expected `MissingEvidence` set) are both documented, invariant-tied, deliberate changes, not silent regressions.

## §Self-check

**Did α's work push ambiguity onto β?**

- AC2's internal scaffold tension (literal 3-guard `all_true` vs. the "no regression for `review-with-pr.json`" instruction) was resolved by α, not left for β to discover — the resolution (update the pre-existing fixture's field values rather than weaken the guard) is documented in §ACs AC2 with explicit reasoning, and the resulting test-assertion change is called out inline in the test file itself, not just in this artifact.
- AC4's design decision (option (b) over (a)) is made and reasoned, not deferred — β does not need to re-litigate fetch-vs-API, only verify the chosen implementation is correct and hermetically tested.
- AC5's edit-vs-append judgment (technically-possible edit, chose append anyway) is explicit with reasoning, not silently defaulted.
- AC7's CI-on-pushed-head-commit gap is named as debt (below), not silently asserted green — β is told exactly what was and wasn't independently verified in this session.

**Is every claim backed by evidence in the diff?**

- AC1: authenticated `gh issue view` output quoted in §ACs; public-render check names the exact selector used and the result.
- AC2/AC3: TDD red→green sequence narrated with the specific test names and the specific guard-string before/after; both are inspectable directly in `transitions.json` and `issuesfsm_test.go`.
- AC4: three new/extracted-function tests named, each mapped to what it proves; the design-decision reasoning is inline as a doc comment in `fetch.go` itself (not only in this artifact — so a future reader of the code, not just this cycle's artifact, sees the reasoning).
- AC5: both GitHub URLs (correction comment, new issue) are live links a reviewer can open directly.
- AC6: the exact `rg` invocation and file set are reproducible by β re-running the same command.
- AC7: exact test-runner counts (37 PASS / 0 FAIL / 48 RUN) pasted from actual output, per the pre-review gate's row 13 requirement, not manually enumerated.

**Peer enumeration:** the family of surfaces touched is small and fully enumerated in §ACs per-AC — `transitions.json` (2 states' rules), `fetch.go` (1 new function + 1 call-site edit), `issuesfsm_test.go` (new tests + one pre-existing assertion update), 4 testdata fixture files (3 new, 1 edited). No sibling command/renderer/writer of the same schema was identified beyond what's already named: `table.go`'s `guardFuncs` registry was checked and confirmed to already contain every guard AC2/AC3 need (no new guard function required, per scope guardrail 6) and was not modified.

**Harness audit:** the only non-Go harness in the wave-touched surface set is the `.github/workflows/*.yml` CI definitions (checked, unchanged — AC4 did not pick the workflow-fetch option) and the JSON transition table itself (which *is* the schema-bearing data structure this cycle's AC2/AC3 changes edit — its consumer, `table.go`'s `Evaluate`, was read and confirmed unchanged/compatible, not re-derived).
