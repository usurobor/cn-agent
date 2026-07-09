# β review — cnos#633

## §R0

**Reviewer:** β (independent review, cycle/633)

### Summary

Independently re-derived every AC in gamma-scaffold.md's oracle table and every guardrail without trusting self-coherence.md's claims. All claims check out. No blocking findings.

### AC1 — self-guarding rule + isolation proof

- Confirmed `transitions.json`'s `in-progress` state rule `[6]` (`propose_status_todo_with_matter`) now has `"all_false": ["run_active", "review_request_present"]` (previously `["run_active"]`). Its `reason` text was updated to document the added guard and its ordering-independence rationale.
- Ran `go test ./... -run 'TestAC633' -v` from `src/packages/cnos.issues/commands/issues-fsm`: all three new tests pass.
  - `TestAC633_TodoWithMatterRuleIsSelfGuardedAgainstReviewRequest` extracts rule `[6]` alone (via `findRuleByAction`, by action name not index) into a single-rule `Table` with no other rule present and asserts no match when `review_request_present: true`. This is the order-independence proof the issue specifically asks for — distinct from a real-table-ordering test.
  - `TestAC633_ReviewReadyCellNeverResolvesToTodoWithMatter` runs the real table against `testdata/in-progress-review-request-with-matter.json` (`review_request_present: true`, `pr_exists: true`, `pr_commit_count: 5`) and asserts the outcome is `propose_status_review`/`review`, explicitly asserting it is never `propose_status_todo_with_matter`.
- **Independent regression-lock verification (not just trusting α's claim):** I reverted rule `[6]`'s `all_false` back to `["run_active"]` via a scoped `sed` (not `git stash`), re-ran `TestAC633_*`, and confirmed `TestAC633_TodoWithMatterRuleIsSelfGuardedAgainstReviewRequest` FAILS red (`propose_status_todo_with_matter matched in isolation against a review_request_present fixture`) while `TestAC633_ReviewReadyCellNeverResolvesToTodoWithMatter` still PASSES (because rules `[0]`/`[1]`'s ordering still covers the full-table case) — this is exactly the ordering-vs-self-guarding distinction the issue asks to close. Restored the fix via `git checkout --`, re-ran all `TestAC633`/`TestAC630` tests, all green again, working tree clean.

**Verdict: AC1 satisfied.**

### AC2 — cnos#630 wedge fixture unaffected

- Ran `testdata/scan-died-after-pr-before-review-request.json` (confirmed on disk: `review_request_present: false`, `pr_exists: true`, `pr_commit_count: 6`, `run_state: completed`) against the fixed real table via the test suite.
- `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved`, `TestAC630_AuditNoteReasonNamesMechanicalReversion`, and the new `TestAC633_WedgeFixStillRequeuesWithoutReviewRequest` all PASS — the fixture still resolves to `propose_status_todo_with_matter`/`todo`, unchanged.

**Verdict: AC2 satisfied.**

### AC3 — full CI gate

- `go test ./../packages/cnos.issues/commands/issues-fsm/...` from `src/go`: `ok` (green, exit 0).
- `go test ./...` from `src/packages/cnos.issues/commands/issues-fsm` (full package, all ~40 tests, not just the new ones): all PASS.
- `go build ./...` from `src/go`: clean, exit 0.
- `go vet ./...` from `src/go`: clean, exit 0.
- `bash scripts/ci/check-dispatch-repair-preflight.sh`: exit 0.
- `bash scripts/ci/check-dispatch-closeout-integrity.sh`: exit 0.

**Verdict: AC3 satisfied.**

### Guardrail verification

- `git diff main...cycle/633 --stat` shows changes to: `.cdd/unreleased/633/CLAIM-REQUEST.yml`, `.cdd/unreleased/633/gamma-scaffold.md`, `.cdd/unreleased/633/self-coherence.md`, `docs/development/board/board-data.json`, `docs/development/board/index.html`, `src/packages/cnos.cds/skills/cds/fsm/transitions.json`, `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go`.
  - **Non-blocking observation:** `docs/development/board/board-data.json` and `docs/development/board/index.html` show up in the three-dot diff against `main`, but these come from two intervening automated "board-map: regenerate docs/development/board from live open issues" commits (`9a2b39f`, `32decae`) that landed on `cycle/633` between the branch point and α's own commit — they are not part of α's authored work. α's own commit (`075773f`) touches exactly 5 files: the three `.cdd/unreleased/633/` artifacts, `transitions.json`, and `issuesfsm_test.go`. Verified via `git show --stat 075773f`. This matches the scaffold's scope exactly; the board-map diffs are pre-existing branch-history noise unrelated to this cell, not a scope violation by α.
- Confirmed via diff that exactly one rule in transitions.json's `in-progress` state changed (rule `[6]`'s `all_false` array + its `reason` text) — no other rule in `in-progress` or any other state touched.
- Confirmed no new `status:*` label anywhere in the diff. The only `status:` occurrences are pre-existing prose references inside `.cdd/unreleased/633/*.md` (scope-guardrail text quoting "no new status:* label") and an unrelated line in the auto-regenerated `board-data.json` describing issue #640's title — none of these are new label additions produced by this cell's code/table changes.
- `issuesfsm_test.go` diff is purely additive (119 new lines at EOF: doc comment, `findRuleByAction` helper, three new `TestAC633_*` tests) — no pre-existing test body modified.

### Findings

None — no blocking or non-blocking defects found in the implementation itself. One non-blocking observational note recorded above regarding the unrelated automated board-map commits present in the branch's ancestry (informational only, not attributable to α, does not affect correctness or scope of the actual fix).

### Verdict

verdict: converge
