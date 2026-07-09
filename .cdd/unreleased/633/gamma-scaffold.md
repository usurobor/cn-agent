# γ R0 scaffold — cnos#633

## Issue reference

- **Issue:** [usurobor/cnos#633](https://github.com/usurobor/cnos/issues/633) — "cds/fsm: lock rule ordering so `propose_status_todo_with_matter` can't shadow the in-progress→review transition"
- **Mode:** design-and-build (defensive hardening)
- **Cell kind:** `implementation` (small test/guard addition)
- **Family:** #630 (the wedge fix this hardens) / #575 (FSM transitions) / #593 (reconciler) / #614 (mandatory-learning doctrine origin for this follow-up)
- **Protocol:** cds
- **Dispatch mode:** wake-invoked-δ (`cds-dispatch`, sigma), per `delta/SKILL.md` §9
- **Base SHA:** `32decae982ccd2246390fcf7a048d3085924d9d2` (`origin/main` HEAD at claim time — "board-map: regenerate docs/development/board from live open issues")
- **Cycle branch:** `cycle/633`, created from `main@32decae9`
- **run_class:** `first_pass` — no prior `cycle/633` branch, no prior PR, no prior `.cdd/unreleased/633/` artifacts, no `status:changes` history on the issue.

## Governing rule (restated from the issue)

`transitions.json`'s `in-progress` state rule `[6]` (`propose_status_todo_with_matter`, added by cnos#630) currently guards only on `all_false: [run_active]` + `all_true: [pr_exists]`. It relies purely on being positioned *after* rules `[0]`/`[1]` (which both gate on `review_request_present`) to avoid reverting a review-ready cell back to `todo`. That protection holds today only because of table ordering — a future reorder or rule insertion could silently regress it and reopen a variant of the #630 wedge. The issue asks for an explicit self-guarding fix (`all_false: [review_request_present]` added to rule `[6]`) plus a regression test proving the invariant holds independent of ordering.

## Key finding — the fix is a one-line guard addition, not a restructure

I read `transitions.json`'s full `in-progress` state (9 rules, `[0]`–`[8]`). Rule `[6]` is:

```json
{
  "all_false": ["run_active"],
  "all_true": ["pr_exists"],
  "outcome": "proposed",
  "action": "propose_status_todo_with_matter",
  "target_state": "todo",
  ...
}
```

Adding `"review_request_present"` to the existing `all_false` array is the minimal, self-guarding fix — the rule already has an `all_false` list, so this is additive, not a shape change. No other rule in the `in-progress` state needs to change: rules `[0]`/`[1]` are unaffected, rules `[2]`–`[5]`, `[7]`, `[8]` don't reference `pr_exists`+`review_request_present` in the same combination.

## Recommended design

**Do both, per the operator's directive:**

1. **Explicit guard** — add `"review_request_present"` to rule `[6]`'s existing `all_false` array in `src/packages/cnos.cds/skills/cds/fsm/transitions.json`.
2. **Regression test** — in `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go`, add:
   - An **isolated-rule test**: extract rule `[6]` alone (from the real, loaded table) into a single-rule `Table`/`StateTransition`, and assert that with `review_request_present: true` + `pr_exists: true` + `run_active: false`, **no rule matches** (`Decision.Outcome` is the zero value, not `"proposed"`) — this proves the guard holds standing alone, independent of rule `[0]`/`[1]`'s ordering, not merely as an emergent effect of table order.
   - A **full-table regression test**: run the real table against a fixture with `review_request_present: true`, `pr_exists: true`, `pr_commit_count > 0`, `run_active: false` — assert the outcome resolves to `propose_status_review`/`review` (rule `[0]`), never `propose_status_todo_with_matter`.
3. **AC2 regression check** — confirm the existing #630 wedge fixture (`testdata/scan-died-after-pr-before-review-request.json`, `review_request_present: false`) is unaffected by the added guard and still resolves to `propose_status_todo_with_matter`/`todo` (existing test `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` already covers this — must still pass unchanged).

## Source of truth

| Claim / surface | Canonical source | Status |
|---|---|---|
| FSM transition table (`in-progress` state, rule `[6]`) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | On `main`; rule `[6]` is the target of this fix |
| Guard/action vocabulary + evaluator engine | `src/packages/cnos.issues/commands/issues-fsm/table.go` (`guardFuncs`, `Rule`, `Evaluate` — first-match-wins, in order) | On `main`; `review_request_present` guard already exists (`ReviewRequestPresent` field), reuse — no new guard needed |
| Existing #630 wedge fixture/tests (must not regress) | `testdata/scan-died-after-pr-before-review-request.json` + `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` / `TestAC630_AuditNoteReasonNamesMechanicalReversion` (`issuesfsm_test.go:1335`, `:1357`) | On `main`; both fixtures set `review_request_present: false` — unaffected by the added guard |
| Fact model | `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` (`FactSnapshot.ReviewRequestPresent`) | On `main`; already observed/JSON-exposed — reuse |
| PR #632 (the #630 fix this hardens) | merged to `main` | Landed; this cell only adds a guard + test, no behavior change to #630's own fixture |

## Per-AC oracle list

| AC | Oracle (mechanical pass/fail) |
|---|---|
| **AC1** | A review-request-bearing `in-progress` cell (fixture: `review_request_present: true`, `pr_exists: true`, `pr_commit_count > 0`, `run_active: false`) never resolves to `propose_status_todo_with_matter` against the real table (resolves to `review` via rule `[0]`) — AND rule `[6]` in isolation (single-rule table) does not match the same facts, proving the guard is self-contained, not order-dependent. |
| **AC2** | Existing #630 behavior unchanged: `testdata/scan-died-after-pr-before-review-request.json` (`review_request_present: false`) still resolves to `propose_status_todo_with_matter`/`todo` — `TestAC630_WedgeFixResolvesToTodoWithMatterPreserved` and `TestAC630_AuditNoteReasonNamesMechanicalReversion` pass unchanged. |
| **AC3** | `go test ./...` from `src/packages/cnos.issues/commands/issues-fsm` green (all pre-existing tests pass unchanged); full CI gate set green (I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity, Go, Package, Binary). |

## Scope (restated from issue)

**In scope:** `transitions.json` rule `[6]`'s guard + regression test(s) in `issuesfsm_test.go` (+ a new testdata fixture if needed).

**Out of scope:** no change to #630's own recovery behavior; no new status labels; no Demo 0; no lifecycle redesign; no touching #626/#639 or the sparse-checkout/write-fence work.

## The α prompt

```
Branch: cycle/633

You are α, implementing cnos#633: "cds/fsm: lock rule ordering so
propose_status_todo_with_matter can't shadow the in-progress→review
transition." Read the full issue body (including the operator's "✅
Dispatched (operator, 2026-07-09)" comment — Authority A2/CAP: proceed
without asking for these small in-scope fixes; escalate only if the fix
would require changing lifecycle semantics or would weaken the #630
recovery behavior). Read .cdd/unreleased/633/gamma-scaffold.md (this
scaffold) in full before starting.

## What you are building

1. In src/packages/cnos.cds/skills/cds/fsm/transitions.json, in the
   "in-progress" state's rule [6] (action: propose_status_todo_with_matter),
   add "review_request_present" to the existing "all_false" array (currently
   ["run_active"] -> ["run_active", "review_request_present"]). This makes
   the rule self-guarding regardless of table ordering. Update the rule's
   "reason" text to note the added guard.
2. In src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go, add
   regression tests proving AC1 (see gamma-scaffold.md's oracle table):
   - an isolated-rule test extracting rule [6] alone into a single-rule
     Table and confirming it does NOT match when review_request_present is
     true (Decision.Outcome stays the zero value) -- proves the guard holds
     standing alone, not merely via ordering
   - a full-real-table test with a fixture carrying review_request_present:
     true + pr_exists: true + pr_commit_count > 0 + run_active: false,
     asserting the outcome resolves to propose_status_review/review, never
     propose_status_todo_with_matter
   Add a new testdata/*.json fixture if one matching this shape doesn't
   already exist (check testdata/in-progress-review-request-with-matter.json
   first -- it may already fit).
3. Confirm AC2: run the existing #630 tests
   (TestAC630_WedgeFixResolvesToTodoWithMatterPreserved,
   TestAC630_AuditNoteReasonNamesMechanicalReversion) and confirm they still
   pass unchanged against the fixed table.
4. Run go test ./... from src/packages/cnos.issues/commands/issues-fsm and
   confirm everything is green.

## Guardrails (binding)

- Do NOT change any other rule in the in-progress state.
- Do NOT weaken or alter #630's own recovery behavior (the wedge fixture
  must still resolve to propose_status_todo_with_matter/todo).
- No new status:* label. No lifecycle redesign. No Demo 0.
- Do not touch #626/#639 or sparse-checkout/write-fence work.
- Write self-coherence.md documenting: which testdata fixture you used/added
  and why, and confirmation every AC above passes.

## When done

Commit, push to cycle/633, append your review-readiness signal to
self-coherence.md, and stop -- you do not dispatch β yourself.
```

## The β prompt

```
Branch: cycle/633

You are β, independently reviewing α's implementation of cnos#633 on
cycle/633. Read .cdd/unreleased/633/gamma-scaffold.md and
.cdd/unreleased/633/self-coherence.md in full before forming any verdict.
Do not take α's self-coherence claims as verified until you have
independently walked each oracle yourself.

## Independent AC walk

- AC1: confirm rule [6] in transitions.json now has "review_request_present"
  in its all_false array. Confirm a test proves a review-request-bearing
  in-progress+matter fixture resolves to review (rule [0]), never
  propose_status_todo_with_matter, against the REAL table. Confirm a
  SEPARATE test proves rule [6] alone (isolated from the rest of the table)
  does not match when review_request_present is true -- this is the
  order-independence proof the issue specifically asks for; a test that
  only exercises the real table's rule ORDER is not sufficient by itself.
- AC2: run testdata/scan-died-after-pr-before-review-request.json against
  the fixed table yourself; confirm it still resolves to
  propose_status_todo_with_matter/todo. Confirm
  TestAC630_WedgeFixResolvesToTodoWithMatterPreserved and
  TestAC630_AuditNoteReasonNamesMechanicalReversion pass unchanged.
- AC3: run go test ./... from src/packages/cnos.issues/commands/issues-fsm
  yourself. Confirm all tests green, no pre-existing test's assertions
  changed except by the deliberate additions named above.

## Guardrail verification

- Confirm the diff touches only transitions.json (rule [6]) +
  issuesfsm_test.go (+ new testdata fixture if added) -- nothing else.
- Confirm no other in-progress rule changed.
- Confirm no new status:* label anywhere in the diff.

## Verdict

Write .cdd/unreleased/633/beta-review.md with a full R[N] section: outcome
per AC, any findings (with severity), and `verdict: converge` or
`verdict: iterate`. Commit, push, and stop -- you do not re-dispatch α
yourself.
```

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go + JSON (declarative rule table) |
| CLI integration target | `cn issues fsm evaluate`/`scan` — no new subcommand, no new binary |
| Package scoping | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (rule `[6]` guard only) + `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` (+ new `testdata/*.json` fixture if needed) |
| Existing-binary disposition | Preserve — no binary replaced or forked |
| Runtime dependencies | None new |
| JSON/wire contract preservation | Additive-only — `all_false` array gains one entry; no field renamed/removed; every pre-existing fixture must still decode and evaluate identically |
| Backward-compat invariant | All pre-existing tests pass unchanged; #630's wedge fixture behavior unchanged |

No TBD rows — this is a small, fully-pinned fix; no escalation to operator needed.

## Scope guardrails (binding on α)

- Touch only `transitions.json` rule `[6]` + `issuesfsm_test.go` (+ optionally one new `testdata/*.json` fixture).
- Do not weaken or alter #630's own recovery behavior.
- No new `status:*` label, no lifecycle redesign, no Demo 0.
- Do not touch #626/#639 or sparse-checkout/write-fence work.

## Friction notes

None — this is a small, well-scoped, fully-specified fix with no open design calls.
