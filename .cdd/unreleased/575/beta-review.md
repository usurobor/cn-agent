# β review — cnos#575

## §R0

**Reviewed head SHA:** `c6ac67f3272745d2166ca7c16d18abb6c6a27254` (branch `cycle/575`)
**Review-diff base:** `origin/main` = `4876278` (re-fetched synchronously via `git fetch --verbose origin main` immediately before computing the diff; matches α's recorded base SHA — branch is not stale against `main`).
**γ artifact present:** `.cdd/unreleased/575/gamma-scaffold.md` confirmed on `origin/cycle/575` (pre-merge gate row 4 — satisfied).

All findings below were produced by running the oracle commands myself, not by reading α's `self-coherence.md` and trusting the pasted output.

### Per-AC verdicts

**AC1 — claim routed through the FSM: PASS**

```
$ jq '.transitions[] | select(.state=="todo")' src/packages/cnos.cds/skills/cds/fsm/transitions.json
{ "state": "todo", "trigger": "evaluate", "rules": [
  { "all_true": ["claim_request_present"], "all_false": ["run_active"],
    "outcome": "proposed", "action": "propose_status_in_progress", "target_state": "in-progress" },
  { "all_true": ["claim_request_present"], "outcome": "blocked", "action": "block", "target_state": "" },
  { "outcome": "valid", "action": "none" } ]}
```
New `todo`-state rule confirmed, gated on `claim_request_present` with `all_false: [run_active]` as the competing-run guard. Verified `testdata/todo-competing-run.json` (run_state=in_progress) resolves to `outcome: blocked`, never `proposed`/`in-progress` (checked directly against the fixture content, not just the pasted decision).

```
$ cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -run TestAC575_1 -v
=== RUN   TestAC575_1_ClaimRoutedThroughFSM
--- PASS: TestAC575_1_ClaimRoutedThroughFSM (0.00s)
=== RUN   TestAC575_1_ClaimBlockedOverCompetingRun
--- PASS: TestAC575_1_ClaimBlockedOverCompetingRun (0.00s)
=== RUN   TestAC575_1_ClaimNotRequestedStaysValid
--- PASS: TestAC575_1_ClaimNotRequestedStaysValid (0.00s)
PASS
```

`rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` → exit 1, zero matches for the claim row (verified as part of the AC4 whole-file oracle below).

**AC2 — hard-block routed through the FSM: PASS**

```
$ jq '.transitions[] | select(.state=="in-progress") | .rules[] | select(.action=="propose_status_blocked")' transitions.json
{ "all_true": ["block_request_present"], "outcome": "proposed",
  "action": "propose_status_blocked", "target_state": "blocked" }
```
Gated on the new `block_request_present` guard (backed by `BLOCK-REQUEST.yml`, wired in `fetch.go`). No-evidence fixture (`in-progress-block-no-evidence.json`) falls through to the pre-existing `run_active` rule rather than being silently applied — confirmed by reading the fixture (`block_request_present` absent, `run_state: in_progress`) and re-running `Evaluate` via the test suite below.

```
$ go test ./... -run TestAC575_2 -v
--- PASS: TestAC575_2_HardBlockRoutedThroughFSM (0.00s)
--- PASS: TestAC575_2_HardBlockRefusedWithoutEvidence (0.00s)
PASS
```

**AC3 — release-back-to-queue routed through the FSM: PASS**

```
$ jq '.transitions[] | select(.state=="in-progress") | .rules[] | select(.all_true[0]=="release_request_present")' transitions.json
{ "all_true": ["release_request_present"], "all_false": ["branch_has_commits","pr_exists","pr_has_commits"],
  "outcome": "proposed", "action": "propose_status_todo", "target_state": "todo" }
{ "all_true": ["release_request_present"], "any_true": ["branch_has_commits","pr_exists","pr_has_commits"],
  "outcome": "proposed", "action": "propose_delta_recovery", "target_state": "" }
```
Confirmed distinct from the existing dead-run reconciliation rules (`all_false: [run_active]`) — the new rules fire on `release_request_present` regardless of `run_active`, positioned earlier in the rule list. Matter-present fixture (`in-progress-release-with-matter.json`, `pr_exists: true`, `commits_beyond_base: 3`) routes to `propose_delta_recovery`, never `target_state: todo` — the #368 protection holds.

Rule-ordering non-shadowing independently re-verified (not just trusting the new `TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules`): re-ran the **full** existing suite (`TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid`) against the modified `transitions.json` — all pass unmodified. Traced the rule list manually via `jq` (reproduced below in the implementation-contract section) and confirmed first-match-wins order places the three new rules before the `all_true:[run_active]` valid rule and before both dead-run rules, so none of the pre-existing fixtures can hit the new rules by accident (none of them set `release_request_present`/`block_request_present`).

```
$ go test ./... -run TestAC575_3 -v
--- PASS: TestAC575_3_ReleaseRoutedThroughFSMWhenNoMatter (0.00s)
--- PASS: TestAC575_3_ReleaseBlockedOverMatter (0.00s)
PASS
```

**AC4 — workers/wake request, don't write: PASS**

```
$ rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
(exit 1 — zero matches)
$ rg -n "wake writes the label directly" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
(exit 1 — zero matches)
$ rg -n "direct label write|wake writes the label directly" src/packages/cnos.cds/skills/cds/CDS.md
(exit 1 — zero matches; α's debt-item self-check re-verified)
```
Read through the full diff of `cds-dispatch/SKILL.md`: frontmatter description, claim-sequence steps 5/6/8, "Lifecycle transitions" table + intro prose, and Responsibilities item 5 all updated coherently to FSM-request phrasing for all four transitions. `dispatch-protocol/SKILL.md`'s diff is exactly the scoped "Protocol-specific FSM override" paragraph in §2.2 α claimed — the generic claim-sequence code block (step 4, direct `gh issue edit`) is untouched, confirmed by reading the diff directly (2-line addition only).

**Rendered-artifact fidelity (independently re-rendered, not trusted from α's claim):**
```
$ ./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch
cn-install-wake: cds-dispatch → .../cnos-cds-dispatch.golden.yml (unchanged)
$ diff <(pre-render copy of workflow) .github/workflows/cnos-cds-dispatch.yml
IDENTICAL
$ diff <(pre-render copy of golden) cnos-cds-dispatch.golden.yml
IDENTICAL
```
Both the committed golden fixture and the live workflow are byte-identical to a fresh render of the current `cds-dispatch/SKILL.md` — no drift.

**AC5 — invariants preserved; gates green: PASS**

```
$ cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -count=1 -v
... (51 tests)
PASS
ok  	github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm	0.033s
$ grep -c "^--- PASS" <full -v output>
51
$ grep -c "^--- FAIL" <full -v output>
0
```
Matches α's claimed 51 PASS / 0 FAIL exactly. All named pre-existing tests (`TestAC574_*` ×6, `TestSeam_CellKind*` ×2, `TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid`) present and passing. Read `TestSeam_CellKindNotEnforced`'s and `TestSeam_CellKindDefaultedWhenAbsent`'s bodies directly (not just names) — only the fixture-iteration list in `TestSeam_CellKindNotEnforced` was extended (6 new fixture paths appended); the assertion logic (decision must be byte-identical across `cell_kind` values) is byte-for-byte unchanged.

```
$ cd src/go && go build ./... && go test ./...
BUILD OK
ok × 13 packages, all green
```

```
$ rg -n "cell_kind" src/packages/cnos.cds/skills/cds/fsm/transitions.json
(exit 1 — zero matches)
```

**CI on branch head (`c6ac67f3`):**
```
$ gh api repos/usurobor/cnos/commits/c6ac67f3272745d2166ca7c16d18abb6c6a27254/check-runs \
    --jq '.check_runs[] | .name + ": " + .conclusion'
Package verification: success
Binary verification: success
CDD artifact ledger validation (I6): success
Protocol contract schema sync (I2): success
Repo link validation (I4): success
Go build & test: success
SKILL.md frontmatter validation (I5): success
Package/source drift (I1): success
Dispatch closeout-integrity guard (cnos#524): success
Dispatch repair-preflight guard (cnos#516): success
```
All 10 check-runs `success`.

**install-wake-golden freshness claim — independently confirmed, not merely trusted:**
- `git diff 0defb27 HEAD -- .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` → empty diff (zero changes to golden-relevant surfaces since the commit `0defb27` that the last green `install-wake golden` run (`28703902087`, success, `2026-07-04T10:53:04Z`) executed against).
- `git diff 0defb27 HEAD` in full touches only `.cdd/unreleased/575/self-coherence.md` and a 2-line `dispatch-protocol/SKILL.md` addition (a lychee-link path fix, not a golden-input path) — confirmed via `gh run list --workflow="install-wake golden"` that `0defb27` is genuinely the most recent run and it is `success`.
- Confirmed the run's headSha's tree already contains `dd3e966`'s cds-dispatch changes (i.e., `0defb27`'s tree for `cds-dispatch/SKILL.md`/golden/workflow is identical to current HEAD's, per the diff above) — so the last green golden run really did validate the current doctrine content, not a stale pre-#575 version.

### Implementation-contract conformance verdict: CONFORMS

Verified each pinned axis against the actual diff (`git diff origin/main...origin/cycle/575`), not the doc:
- **Language:** 100% Go (`decision.go`, `fetch.go`, `snapshot.go`, `table.go`, `issuesfsm_test.go`) + JSON data (`transitions.json`) + prose (two SKILL.md files) + two machine-rendered artifacts. No new runtime.
- **CLI integration target:** `git diff origin/main...origin/cycle/575 -- src/go/internal/cli/cmd_issues_fsm.go` → empty. No new sub-verb, no new binary.
- **Package scoping:** full diff-stat (18 files) confined to `cnos.cds/skills/cds/fsm/transitions.json`, `cnos.issues/commands/issues-fsm/{decision,fetch,snapshot,table,issuesfsm_test}.go` + 6 new `testdata/*.json`, `cds-dispatch/SKILL.md` + its two rendered-artifact peers (golden + live workflow — a downstream-generated-artifact consequence of the SKILL.md edit, not a scope-creep axis), `dispatch-protocol/SKILL.md`, plus the `.cdd/unreleased/575/` artifact set. Nothing outside this set.
- **Existing-binary disposition:** `cn` extended, no new binary — confirmed no `go.mod`/`go.sum` changes (`git diff origin/main...origin/cycle/575 -- '**/go.mod' '**/go.sum'` empty) and no new `cmd/` entries.
- **Runtime dependencies:** none added.
- **JSON/wire contract preservation:** `git diff` on `transitions.json` shows only new guard-doc entries, a new `_doc_phase3` block, two new `todo`-rules, and three new `in-progress`-rules inserted before the pre-existing `run_active`/dead-run rules — the `review` and `changes` state blocks have zero diff hunks touching them (confirmed the diff contains no `"state": "review"` or `"state": "changes"` hunk).
- **Backward-compat invariants:** `run_active` non-gating preserved for the new hard-block/release rules (deliberately `all_true`-gated on request markers only, no `run_active` guard added to them); #574's review-guard rules structurally untouched; `cell_kind` observed-only preserved (`rg` zero-match + `TestSeam_CellKindNotEnforced` unmodified assertion, extended fixture list, all passing).

No axis drift found. This is not a rubber-stamp: I independently re-derived the diff-to-contract mapping rather than accepting α's self-coherence table at face value.

### Scope guardrail verdict: HELD

- No Sub 3/4 work (no PR-open/checkpoint/recovery runtime code found in the diff).
- `cell_kind` untouched for enforcement (verified above).
- No new status labels: `rg`'d the full diff for `status:new`/taxonomy language — the four transitions use the existing `status:todo`/`in-progress`/`blocked`/`review` set only.
- No Demo 0, no #585/#586 touches (diff-stat's 18 files are all within the pinned surface; no file paths reference those issues' surfaces).
- `transitions.json`'s `review`/`changes` state blocks are additive-only elsewhere (`todo`/`in-progress`) — confirmed via diff hunk inspection above.
- `delta/SKILL.md` — `git diff origin/main...origin/cycle/575 -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` → empty. Confirms Friction note 3's claimed non-edit.
- Workflow execution semantics (`permissions`, `on:` triggers, `runs-on`, `concurrency`) in `.github/workflows/cnos-cds-dispatch.yml` are byte-identical to `origin/main` — verified via a structural YAML diff with `steps` stripped out (`diff <(yq/jq del(.jobs[].steps) on main) <(same on cycle/575)` → empty). Only the embedded step-body prose text (the rendered SKILL.md doctrine) changed, which is the expected downstream artifact of the SKILL.md edit, not a scope-creeping change to how/when/with-what-permissions the workflow runs.

### Friction-note resolution verdict: ALL SIX RESOLVED, ALL COHERENT

1. **Claim guard framing** — hybrid (existing `run_active`, unmodified, + new `claim_request_present` marker) confirmed in the actual `transitions.json` diff (`all_true: ["claim_request_present"], all_false: ["run_active"]`). This does not reopen the #368/#569 tension: `run_active` still isn't gating the *review*/*hard-block*/*release* requests (those three new/existing rules never reference `run_active` in `all_true`), and its use in the *claim* rule is a genuinely different case — the requesting wake's own run cannot appear "active" against a `cycle/{issue}` branch that doesn't exist yet at claim time, so `run_active` there means "some **other** run", not "my own run" — coherent and distinct from the #569 tension it might otherwise re-trigger.
2. **`dispatch-protocol/SKILL.md` scope** — diff confirms a scoped 1-paragraph "Protocol-specific FSM override" addition in §2.2, not a rewrite; the generic step-4 direct-write code block is untouched (verified directly in the diff, not from α's description).
3. **`delta/SKILL.md` §9.6** — `git diff` confirms zero changes to this file. Leaving it as a named, undecided follow-up rather than a silent edit or a silent gap is a defensible non-goal given the issue's explicit AC4 surface list names only `cds-dispatch/SKILL.md` + `dispatch-protocol/SKILL.md`.
4. **Test naming** — 14 `TestAC575_*` functions, zero collisions with existing `TestAC1`-`TestAC7`/`TestAC569_*`/`TestAC574_*` (spot-checked via `grep -o "^func Test..." | sort | uniq -c`, no duplicates).
5. **Fixture naming** — all six new fixtures follow `<state>-<condition>.json`, consistent with existing convention (spot-checked fixture filenames and contents directly).
6. **Rule ordering** — confirmed by reading the actual `jq`-extracted rule order for the `in-progress` block (reproduced above under AC3): hard-block, release-no-matter, release-with-matter are positioned before the `run_active` valid rule and both dead-run rules. Re-ran the full pre-existing test suite against the modified table — all originally-targeted fixtures (`in-progress-dead-no-matter.json`, `in-progress-dead-with-commits.json`, `in-progress-active.json`) still resolve through their original rules (unmodified pass).

### Overall verdict: **converge**

Every AC oracle passes on independently-executed commands (not paraphrased from α's self-coherence.md). The implementation-contract conformance gate (Rule 7) holds on every pinned axis, verified against the diff directly. Scope guardrails held — no drift into Sub 3/4, `cell_kind` enforcement, new taxonomy, Demo 0, or unrelated doctrine files. All six friction notes were resolved with real, checkable justification, and each resolution was independently confirmed against the actual diff rather than accepted from the narrative. CI is green on the reviewed head SHA (10/10 check-runs `success`), and the install-wake-golden freshness claim was independently re-derived (zero diff on golden-relevant paths since the last green golden run, and a fresh local re-render is byte-identical to the committed golden/workflow).

No unresolved findings. Proceeding to merge per β's role (pre-merge gate rows to be walked next, including the non-destructive merge-test in a throwaway worktree).
