# α close-out — cnos#633

## Scope note

Written at the δ §9.5 converge boundary — β returned `verdict: converge` at R0 (`beta-review.md`), and this file is one of the three closeout artifacts that boundary requires. δ's next step is to open the cycle-PR and request the `status:review` transition.

## Summary — what was built

`transitions.json`'s `in-progress` rule `[6]` (`propose_status_todo_with_matter`, cnos#630) relied purely on table ordering — being positioned after the `review_request_present`-gated rules `[0]`/`[1]` — to avoid reverting a review-ready cell back to `todo`. Added `"review_request_present"` to the rule's own `all_false` array, making it self-guarding regardless of future reordering. Added three regression tests: an order-independence proof (the rule extracted alone into a single-rule table, with no other rule present, still refuses to match a `review_request_present: true` fixture), a full-real-table proof, and a cnos#630 wedge-behavior regression check.

## What worked well

- The fix was a one-line, additive change to an existing `all_false` array — no restructuring, no new guard needed (`review_request_present` was already in `table.go`'s guard registry from cnos#569).
- A pre-existing fixture (`testdata/in-progress-review-request-with-matter.json`) already carried exactly the shape the new positive-case tests needed — no new testdata file required.
- The isolated-rule extraction technique (pulling rule `[6]` out by action name into a standalone single-rule `Table`) gave a clean, mechanical way to prove "self-guarding regardless of ordering" as a real assertion rather than an inference from reading the rule in context.

## What was harder than expected

Nothing of note — this was a small, fully-scoped fix with no open design calls and no scope ambiguity.

## Debt carried forward

None. No PR opened yet for `cycle/633` — that is δ's next action per `delta/SKILL.md` §9.6, not α's or β's.

## Calibrated success claim

All three ACs (AC1–AC3) pass, independently re-derived by β from the code/tests/diff (β re-ran every test cited in `beta-review.md` from a clean shell, and additionally reverted the guard via a scoped `sed` to independently confirm the isolated-rule test fails red without the fix). The full pre-existing `issues-fsm` test suite passes unchanged alongside the three new tests. `go build`/`go vet` clean; both dispatch-integrity guard scripts (`#516`, `#524`) exit 0. Not yet verified: PR/merge (δ's next step), and live post-merge CI (no merge exists yet).
