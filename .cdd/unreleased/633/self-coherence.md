# self-coherence — cnos#633

## §R0

**Author:** α (dispatch-wake-invoked, sigma)

### Summary

Implemented both fixes the operator's dispatch directive requested:

1. **Explicit guard** — added `"review_request_present"` to `transitions.json`'s `in-progress` state, rule `[6]` (`propose_status_todo_with_matter`)'s existing `all_false` array (`["run_active"]` → `["run_active", "review_request_present"]`). Updated the rule's `reason` text to document the addition and why it makes the rule self-guarding regardless of table ordering.
2. **Regression tests** — added three tests to `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go`:
   - `TestAC633_TodoWithMatterRuleIsSelfGuardedAgainstReviewRequest` — extracts rule `[6]` alone (by action name, not index, via a new `findRuleByAction` helper) into a single-rule `Table` with no other rule present, and asserts it does not match against a `review_request_present: true` fixture. This is the order-independence proof: it cannot rely on rules `[0]`/`[1]` being present or ordered ahead of it, because they are absent from the isolated table entirely.
   - `TestAC633_ReviewReadyCellNeverResolvesToTodoWithMatter` — full real-table test asserting a review-request-bearing in-progress cell with matter resolves to `propose_status_review`/`review`, and explicitly never to `propose_status_todo_with_matter`.
   - `TestAC633_WedgeFixStillRequeuesWithoutReviewRequest` — AC2 regression check: the #630 wedge fixture (`review_request_present: false`) still resolves to `propose_status_todo_with_matter`/`todo` after the guard addition.

### Fixture choice

No new testdata fixture was needed. `testdata/in-progress-review-request-with-matter.json` (pre-existing, already used by several other tests) already carries `review_request_present: true`, `pr_exists: true`, `pr_commit_count: 5` — exactly the shape needed for both new positive-case tests. For the isolated-rule test, `RunState` is forced to `"completed"` in the test body (the fixture's own `run_state` is `"in_progress"`) so that `run_active` is false and the only guard under test — `review_request_present` — is what determines the match/no-match outcome; leaving `run_active` true would trivially block a match for an unrelated reason and not exercise the guard being tested.

### AC verification

- **AC1** — confirmed. `TestAC633_TodoWithMatterRuleIsSelfGuardedAgainstReviewRequest` proves order-independence (rule `[6]` alone refuses to match); `TestAC633_ReviewReadyCellNeverResolvesToTodoWithMatter` proves the full-table behavior. Sanity-checked the isolated-rule test is a real regression lock: temporarily reverted the `all_false` addition (`sed` in place, not via `git stash`, to avoid disturbing the test-file edit) and re-ran — the isolated test correctly **failed** (`propose_status_todo_with_matter matched in isolation against a review_request_present fixture`), while the full-table test still passed (because rules `[0]`/`[1]`'s ordering still covered it) — this is exactly the "ordering vs self-guarding" distinction the issue asks to close. Restored the fix; all three new tests pass.
- **AC2** — confirmed. `TestAC633_WedgeFixStillRequeuesWithoutReviewRequest` (new) plus the pre-existing `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` and `TestAC630_AuditNoteReasonNamesMechanicalReversion` all pass unchanged against the fixed table.
- **AC3** — confirmed. `go test ./...` from `src/packages/cnos.issues/commands/issues-fsm` is green (all pre-existing tests pass unchanged, three new tests added). `go build ./...` and `go vet ./...` from `src/go` are clean. `scripts/ci/check-dispatch-repair-preflight.sh` and `scripts/ci/check-dispatch-closeout-integrity.sh` both exit 0.

### Diff scope

Touched exactly two files: `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (rule `[6]`'s `all_false` array + `reason` text) and `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` (three new tests + one new helper, `findRuleByAction`). No other rule in the `in-progress` state changed. No new `status:*` label. No lifecycle redesign. No touch to #626/#639 or sparse-checkout/write-fence surfaces.

### Review-ready signal

R0 implementation complete, all oracles green, diff scope matches the scaffold exactly. Ready for β review.
