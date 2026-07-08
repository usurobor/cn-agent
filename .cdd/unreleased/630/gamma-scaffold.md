# γ R0 scaffold — cnos#630

## Issue reference

- **Issue:** [usurobor/cnos#630](https://github.com/usurobor/cnos/issues/630) — "partial-matter in-progress wedge: give the reconciler a mechanical path out"
- **Mode:** design-and-build
- **Cell kind:** `implementation` (reconciler gap fix, per `CELL-KINDS.md` §2)
- **Family:** #583 (mechanical-dispatch doctrine) / #593 (recovery reconciler, `cn issues fsm scan`) / #575 (FSM lifecycle transitions, Phase 3) / #591 (checkpoint+PR-open finalizer) / #614 (concrete evidence — the stranded cell) / #615 (terminal-hygiene reconciler, sibling precedent for this scaffold's shape) / #368 (blind-requeue-over-matter failure the fix must avoid) / #626 (adjacent, explicitly out of scope)
- **Protocol:** cds
- **Dispatch mode:** bootstrap-δ, single-session δ-as-γ via Agent tool (Claude Code) — per `.cdd/DISPATCH` §5.2. Not wake-invoked-δ (`delta/SKILL.md` §9.1 table: wake-invoked-δ "not yet observed in production" as of this cycle).
- **Base SHA:** `6143b53c9098e415c4e7e83791843bcad2313314` — δ's pinned re-base (STOP decision (a): re-pinned from the original claim-time SHA to current `origin/main` HEAD after confirming the two trailing commits, `8267b5d3`/heartbeat and `6143b53c`/board-map regen, touch only `.cn-sigma/logs/` and `docs/development/board/**` — no `src/`, no skill files, no FSM/reconciler code). Verified: `git rev-parse origin/main` == `6143b53c9098e415c4e7e83791843bcad2313314` exactly (zero further drift at scaffold time) and `.cdd/unreleased/630/CLAIM-REQUEST.yml`'s `head_commit` (`313cb1f3a4e5d809f09e5c51660351ec1ed189c1`, one commit older — the claim itself, `5231430a`, was the commit immediately after that head, and is included in the pinned SHA's ancestry).
- **Cycle branch:** `cycle/630`, created from `main@6143b53c`.
- **run_class:** `first_pass` — per `CLAIM-REQUEST.yml`'s own note: no prior `status:changes` history (label timeline `status:ready -> status:todo` only), no prior `cycle/630` branch (confirmed: `git ls-remote origin` had no `refs/heads/cycle/630` before this scaffold), no prior `.cdd/unreleased/630/` artifacts beyond `CLAIM-REQUEST.yml` itself, zero issue comments beyond the operator's dispatch-authorization comment.

## Governing rule (restated from the issue)

Fix the **"partial-matter in-progress wedge"**: a `status:in-progress` dispatch cell whose run dies *after* producing matter (a `cycle/{N}` branch + a checkpointed draft PR, per the #591 finalizer) but *before* `REVIEW-REQUEST.yml` is written has **no mechanical path out** today. It cannot requeue (matter exists, so the FSM's release-back-to-queue guard blocks it), cannot advance (no `REVIEW-REQUEST.yml`, so `in-progress -> review` is blocked), and is never re-claimed (the dispatch selector only claims `status:todo`). AC1/AC2 require this to self-heal without a human. AC3 requires the claim path to *resume* rather than defer when it finds a `status:todo` cell with pre-existing `cycle/{N}` matter. AC4 requires mechanical state reversions to leave a machine-readable audit note.

## Key finding 1 — the wedge is the literal, already-shipped no-op branch in `scan.go`

Per `gamma/SKILL.md` §2.2a peer-enumeration: I read `src/packages/cnos.issues/commands/issues-fsm/scan.go` in full (197 lines of orchestration logic) rather than asserting the gap from the issue text alone. The wedge is not hypothetical — it is the `else` branch of `scanOne`'s `propose_delta_recovery` case (`scan.go:252-256`):

```go
} else {
    // A PR (or the branch/commits it depends on) already existed;
    // cn cell finalize's own idempotent no-op path ran. Nothing
    // new happened this tick -- do not re-comment.
    res.Note = "dead run with matter: PR (or commits) already checkpointed; cn cell finalize confirmed idempotent no-op. Status remains in-progress -- awaiting REVIEW-REQUEST.yml or operator/δ follow-up."
}
```

Once a dead in-progress cell's matter is checkpointed (draft PR open), every subsequent scan tick re-evaluates to the identical `propose_delta_recovery` decision and hits this exact `else` branch again — forever. There is no rule anywhere in `transitions.json`'s `in-progress` state (verified: read the full table, all 8 rules) that proposes anything other than `propose_delta_recovery` (target_state `""`, i.e. no label move) for "dead run + matter present." This is confirmed by the **existing fixture and test** `testdata/scan-died-after-pr-before-review-request.json` / `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` (`scan_test.go:199`) — that test's own name asserts the current (buggy, per this issue) terminal behavior: finalize no-ops, no comment, no further action. **AC2's fixture is not new ground — it is this exact scenario, re-asserted to reach a *different*, non-wedged terminal outcome.** α should treat `scan-died-after-pr-before-review-request.json` as the base fixture to extend/redirect, not duplicate.

## Key finding 2 — the FSM already half-solves AC3, but the claim-prose layer doesn't know it

I read `transitions.json`'s `todo` state (3 rules) and the `todo-competing-run.json` fixture. The `todo -> in-progress` claim rule gates *only* on `claim_request_present && !run_active` — it does **not** check `branch_exists`/`pr_exists`. Mechanically, the FSM already permits claiming a `status:todo` cell that happens to have a pre-existing `cycle/{N}` branch + draft PR (no competing run). **The actual gap AC3 names lives one layer up**, in the prose-authored claim doctrine (`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`, rendered to `.github/workflows/cnos-cds-dispatch.yml`) and in δ's wake-invoked routing (`delta/SKILL.md` §9) — neither currently *names* "detect pre-existing `cycle/{N}` + open draft PR at claim time and resume from it" as a step. The closest existing precedent is §9.10 "resumed-from-changes shape" (cnos#500) — a different transition (`status:changes -> status:todo -> status:in-progress`, re-claim after an operator verdict) that already establishes the pattern this cell needs for a *different* trigger (reconciler-driven safe-requeue, not operator verdict): δ detects prior-round artifacts on the cycle branch (`CDDArtifacts` in `FactSnapshot`, already observed and JSON-exposed — see `snapshot.go`) and routes accordingly instead of re-scaffolding from scratch. #614's own recovery (see `.cdd/unreleased/614/gamma-closeout.md` — "resuming from a prior firing's infra-interrupted `gamma-scaffold.md`") is empirical evidence this resume pattern already works ad hoc under bootstrap-δ; this cell's job is to name it as a **generalized, mechanical, wake-invoked-safe** shape (§9.10's sibling) rather than something only bootstrap-δ can improvise.

**Grep confirms no prior art collides:** `grep -rln "partial-matter|unwedge|strand" src/` returns zero Go/skill hits outside `.cdd/` (i.e., outside prior cycle artifacts) — this is a genuinely new mechanism, not a duplicate.

## Recommended design (α may diverge with documented rationale — see "Open design call")

The issue itself names two candidate mechanisms and says "the design sub picks the mechanism." Given Key finding 2, I recommend **(b) safe requeue-with-matter** as the single coherent mechanism, because it makes AC1/AC2/AC3/AC4 fall out of *one* new FSM rule + *one* claim-doctrine addition, rather than two independent mechanisms:

1. **New `transitions.json` `in-progress` rule** (positioned after the existing dead-with-matter `propose_delta_recovery` rule, or replacing its "already checkpointed" half): when a dead run's matter is *already checkpointed* (branch + PR exist, no live run, no `REVIEW-REQUEST.yml`) — i.e., exactly the state `scan.go`'s current `else` branch no-ops on — propose `in-progress -> todo` (a new/extended action, e.g. `propose_status_todo_with_matter`), **preserving** the branch/PR (never deleting them — this is the #368 protection restated for the *reconciler's own* action, not just the wake's `RELEASE-REQUEST.yml` path). This is what actually unwedges the cell: it re-enters the claim queue instead of being invisible to it forever.
2. **`scan.go`'s reconciliation of that new proposed action**: apply the label move via the same `applyStatusLabel` primitive `evaluate --apply` already uses (no new label-writer), and post the **AC4 audit-note comment** naming this as a *mechanical reversion* (not an operator/κ action) so the next firing's claim-time observation doesn't read it as an unexplained orphan. `scan.go` already has a working audit-comment precedent for FSM-driven label moves (`case dec.Outcome == "proposed" && dec.TargetState != ""` at `scan.go:259-283` already posts a reconciliation comment) — reuse that shape rather than inventing a new one.
3. **Claim-doctrine addition** (`cds-dispatch/SKILL.md` §"claim" steps + `delta/SKILL.md` §9, new subsection sibling to §9.10): when the claim sequence's well-formedness check (`cds-dispatch/SKILL.md` step 3-4) observes a `status:todo` cell that *also* carries a pre-existing `cycle/{N}` branch + open PR (this is now a **normal**, audit-noted outcome of step 1-2 above, not an anomaly), δ detects this as a resume shape — mirroring §9.10's `CDDArtifacts`-presence detection — and routes γ/α/β to continue from the existing branch/PR (γ reads what's already there before re-scaffolding; if a scaffold already exists, γ does not overwrite it) rather than deferring as "orphaned prior claim, unexplained" (AC3, AC1's "no human intervention").

This design keeps `transitions.json` as the single source of truth for the *decision* (per `table.go`'s own doc comment: "never switched on a CDS-specific state name... entirely table-driven") and keeps the *doctrine* (what "resume" means for a claimed cell) in the existing prose layer that already owns it (`cds-dispatch/SKILL.md`, `delta/SKILL.md` §9), rather than inventing a third location.

## Open design call α must resolve and document

Whether to (i) extend the existing `propose_delta_recovery` rule's `else`-branch behavior in-place, or (ii) add a wholly new rule/action distinguishing "matter not yet checkpointed" (still `propose_delta_recovery`, finalize runs) from "matter already checkpointed, no review-request, no live run" (new `propose_status_todo_with_matter` or equivalent). I recommend (ii) — the two states are observably different (`scan.go`'s own `strings.Contains(stdout, "opened draft PR for")` branch already distinguishes them) and collapsing them risks re-triggering `cn cell finalize` redundantly on every tick before also proposing the todo move. Pick one, implement it, name the exact new action string(s) added to `table.go`'s action-id vocabulary, and write one paragraph in `self-coherence.md` explaining the choice — mirroring the documented-design-call precedent in cnos#574 AC4 / cnos#615's own scaffold.

A second open call: whether AC3's "resumed, not deferred" claim-time behavior is best implemented as new prose-only doctrine (`cds-dispatch/SKILL.md` + `delta/SKILL.md` §9, no new Go) or requires a Go-level assist (e.g. a `cn issues fsm evaluate`/`scan` flag that surfaces "this candidate has pre-existing matter" for the wake/δ to consume mechanically, rather than the wake re-deriving it from raw `git`/`gh` calls each time). Given `FactSnapshot.CDDArtifacts` + `BranchExists`/`PRExists` are already computed and JSON-exposed by the existing `evaluate`/`scan` machinery, I recommend **no new Go surface** for this — the wake/δ already has everything it needs from `cn issues fsm evaluate --issue {N} --json` at claim time; this is a doctrine-naming gap, not a data gap. Confirm or refute this in `self-coherence.md` with evidence (what data, if any, the wake's claim sequence is actually missing).

## Source of truth

| Claim / surface | Canonical source | Status |
|---|---|---|
| Reconciler this cell extends | `src/packages/cnos.issues/commands/issues-fsm/scan.go` (`RunScan`/`scanOne`, the `propose_delta_recovery` case at lines 222-257 — the exact wedge) | On `main`; extend, do not fork |
| FSM transition table (`in-progress`, `todo` states) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | On `main`; `in-progress` has 8 rules, `todo` has 3 — read in full at scaffold time |
| Guard/action vocabulary + evaluator engine | `src/packages/cnos.issues/commands/issues-fsm/table.go` (`guardFuncs`, `Rule`, `Evaluate`) | On `main`; guard vocabulary is generic/table-driven — do not switch on CDS-specific state names in Go |
| Fact model (branch/PR/artifact observation) | `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` (`FactSnapshot`) + `fetch.go` (live assembly) | On `main`; `CDDArtifacts`, `BranchExists`, `PRExists`, `ReviewRequestPresent` already observed — reuse, do not re-derive |
| Existing near-duplicate fixture (base for AC2) | `testdata/scan-died-after-pr-before-review-request.json` + `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` (`scan_test.go:199`) | On `main`; this test's current assertion is the wedge behavior AC2 requires to change |
| Claim doctrine (prose, rendered to workflow) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` §"claim" (steps 1-8) + `.github/workflows/cnos-cds-dispatch.yml` (rendered artifact — re-render via `cn install-wake cds-dispatch`, never hand-edit) | On `main` |
| δ wake-invoked routing + resume precedent | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 (esp. §9.10 "resumed-from-changes shape," the closest existing precedent for this cell's AC3/AC1 resume behavior) | On `main` |
| Finalizer this cell's reconciliation calls (unchanged) | `src/go/internal/cell/cell.go` (`cn cell finalize`, invoked via `scan.go`'s `liveRunFinalize` self-re-exec) | On `main`; out of scope to modify — only its *caller*'s post-checkpoint next-step changes |
| #368 blind-requeue protection (must not regress) | `transitions.json`'s existing `release_request_present` + `all_false: [branch_has_commits, pr_exists, pr_has_commits]` rule, and the dead-with-matter `propose_delta_recovery` rule's own reasoning text | On `main`; the new rule this cell adds must preserve branch/PR on every requeue, exactly like these do |
| CI gates named in AC5 | `.github/workflows/build.yml` (`go`, `package-verify`/Package, `binary-verify`/Binary), dispatch guard scripts (`check-dispatch-repair-preflight.sh` / #516, `check-dispatch-closeout-integrity.sh` / #524) | On `main` |
| Concrete incident evidence (#614) | `.cdd/unreleased/614/gamma-closeout.md`, `alpha-closeout.md` (resume-from-infra-interruption narrative — the empirical precedent this cell generalizes) | On `main` (614 not yet release-moved) |

## Per-AC oracle list

| AC | Oracle (mechanical pass/fail) |
|---|---|
| **AC1** | A fixture/fake-server scenario matching "in-progress, matter present (branch+PR), no `REVIEW-REQUEST.yml`, no live run" run through the new reconciliation path (`cn issues fsm scan --apply` or equivalent) produces a **non-`""` `TargetState`/action** on first tick after checkpointing (not the current permanent no-op) — assert the report shows a concrete next action taken (label move and/or resume-drive), not `res.Note` reading "awaiting... follow-up" indefinitely. |
| **AC2** | The #614-shaped fixture (extend/redirect `testdata/scan-died-after-pr-before-review-request.json`): run it against `main`-at-base (pre-fix) — assert the existing `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment`-style behavior (permanent no-op) reproduces the strand (test **fails** / demonstrates the bug against pre-fix logic, per the issue's explicit "must fail before the fix" requirement — this may need a new test file/branch that intentionally asserts against old behavior first, or a documented before/after pair in `self-coherence.md`). Then assert against the fixed code: the cell reaches a valid next state (resumed to review, or safely requeued+resumed) without any human-authored label edit or comment in the fixture/test harness. |
| **AC3** | A fixture/scenario: `status:todo` + `claim_request_present` + pre-existing `branch_exists`/`pr_exists` (mirroring `todo-competing-run.json`'s shape but with `run_active: false`) run through the claim path is **resumed** — i.e., the FSM still proposes `todo -> in-progress` (already true today per Key finding 2 — confirm no regression), **and** the claim-doctrine/δ-routing layer's documented behavior (self-coherence.md narrative + any new doctrine text) demonstrably treats this as "continue from existing matter," not "ambiguous, defer." If no Go-testable surface exists for the doctrine-only half, the oracle is the doctrine text itself plus a self-coherence.md walk-through citing the exact new/amended prose section. |
| **AC4** | Every new mechanical state-reversion this cell introduces (the new todo-with-matter requeue, and any other label move the reconciler makes) posts a comment via the existing `PostComment`/audit-note shape (`scan.go`'s existing reconciliation-comment pattern at lines 274-282) that a subsequent fixture run's "claim sees a documented reversion" scenario does **not** classify as an unexplained orphan (i.e., a follow-on claim/scan pass against the post-reversion state produces zero "ambiguous"/"orphaned" notes). |
| **AC5** | No regression: `go test ./...` from `src/packages/cnos.issues/commands/issues-fsm` — every **pre-existing** test (`TestScan_*`, `TestAC1Sub593_ScanSubcommandWired`, all `issuesfsm_test.go`/`terminal_test.go` cases) still passes unchanged (except the one(s) AC2 deliberately redirects — name exactly which, and why, in `self-coherence.md`). `TestSeam_CellKindNotEnforced` (locate via `grep -rn TestSeam_CellKindNotEnforced src/`) passes. Dispatch guard scripts `check-dispatch-repair-preflight.sh` (#516) and `check-dispatch-closeout-integrity.sh` (#524) pass unchanged. Full CI gate set green (see Source-of-truth row). |
| **AC6** | Idempotence, mirroring `scan_test.go`'s existing `TestScan_Idempotent` shape: run the new reconciliation path twice against the same fixture state. First pass reconciles/requeues; second pass (against the post-reconciliation state) makes zero further label mutations, zero duplicate comments, zero duplicate PRs (finalizer's own idempotence, already covered — do not re-test that path, only the new requeue/resume path's idempotence). |

## Scope (restated from issue)

**In scope:** a mechanical path out of the partial-matter in-progress wedge (resume or safe-requeue-with-matter — design sub's choice, see above); wiring "resume from pre-existing `cycle/{N}` branch + draft PR" into the claim doctrine so it isn't deferred as an orphan; a machine-readable audit note on mechanical state reversions; the #614-wedge fixture (AC2) and the todo-with-existing-PR fixture (AC3); no regression to #575/#591/#593/#516/#524.

**Out of scope (non-goals, restated):** no new `status:*` labels; no change to what `status:review` requires (deliverable-proof bar stays); no Demo 0; no `cell_kind` enforcement; no #626 work (cell⇄substrate identity isolation — separate architectural wave, explicitly named as "next up" after this cell in the issue's operator-authorization comment, not part of it).

## The α prompt

```
Branch: cycle/630

You are α, implementing cnos#630: "partial-matter in-progress wedge: give
the reconciler a mechanical path out." Read the full issue body (including
the operator's "Dispatch authorized (operator, 2026-07-08)" comment — it
sets Authority A2/CAP: you may fix in-scope CI, fixtures, prompt/prose, PR
body, and receipt honesty WITHOUT asking; escalate only if the fix would
require a new label, weaken status:review's deliverable-proof bar, or touch
worker-identity/substrate boundaries (that is #626, explicitly out of
scope)). Read .cdd/unreleased/630/gamma-scaffold.md (this scaffold) in full
before starting — it is your scaffold contract for this round.

## What you are building

A mechanical path out of the "partial-matter in-progress wedge": a
status:in-progress cell whose run died after producing matter (cycle/{N}
branch + checkpointed draft PR, per #591's finalizer) but before
REVIEW-REQUEST.yml currently stalls forever (see scan.go:252-256, the exact
no-op branch this cell must change) because it can neither requeue (matter
exists, #575's blind-requeue guard blocks it) nor advance (no
REVIEW-REQUEST.yml) nor get re-claimed (dispatch only claims status:todo).

## Recommended design (you may diverge — see "Open design call" below)

1. Add a new transitions.json rule (or extend the existing dead-with-matter
   in-progress rule) so that "matter already checkpointed (branch+PR exist),
   no live run, no REVIEW-REQUEST.yml" proposes in-progress -> todo,
   PRESERVING the branch/PR (never deleting them -- this is the #368
   blind-requeue protection restated for the reconciler's own action, not
   just the wake's RELEASE-REQUEST.yml path). This is the mechanical exit:
   it puts the cell back in the claim queue instead of leaving it invisible
   to it forever.
2. In scan.go, reconcile that new proposed action via the SAME
   applyStatusLabel primitive evaluate --apply already uses (no new
   label-writer), and post an audit-note comment (reuse the existing
   reconciliation-comment shape at scan.go:274-282) naming this as a
   MECHANICAL reversion -- this is AC4: a subsequent claim/scan pass must be
   able to tell this was not an unexplained orphan.
3. Update the claim doctrine (src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
   "claim" steps, and delta/SKILL.md §9 -- add a subsection sibling to
   §9.10 "resumed-from-changes shape") so that when the claim sequence finds
   a status:todo cell with a pre-existing cycle/{N} branch + open PR (now a
   NORMAL, audit-noted outcome of steps 1-2, not an anomaly), δ resumes from
   the existing matter (reads what's already on the branch, e.g. does not
   overwrite an existing gamma-scaffold.md) rather than deferring as
   "orphaned prior claim, unexplained." This is AC3 / AC1's "without human
   intervention" requirement. FactSnapshot already exposes BranchExists,
   PRExists, and CDDArtifacts (snapshot.go) -- reuse these, do not invent a
   new data surface unless you find a concrete gap and document it.

## Open design call you must resolve and document

(i) Extend the existing propose_delta_recovery rule's behavior in place, vs
(ii) add a distinct new action/rule distinguishing "matter not yet
checkpointed" from "matter already checkpointed, no review-request, no live
run." I recommend (ii) -- scan.go's own strings.Contains(stdout, "opened
draft PR for") check already distinguishes these two states, and collapsing
them risks re-running cn cell finalize redundantly every tick before also
proposing the todo-with-matter move. Pick one, implement it, name the exact
new action string(s), and write one paragraph in self-coherence.md.

Second open call: whether AC3's claim-time resume behavior needs any new Go
surface, or is doctrine-only (cds-dispatch/SKILL.md + delta/SKILL.md §9
prose, no new code). I recommend doctrine-only -- `cn issues fsm evaluate
--issue {N} --json` already exposes everything (BranchExists, PRExists,
CDDArtifacts) the claim sequence needs. Confirm or refute with evidence in
self-coherence.md.

## Fixtures (AC2, AC3, AC6)

- AC2: extend or add alongside testdata/scan-died-after-pr-before-review-request.json
  (the existing near-duplicate of the #614 scenario) a fixture/test pair
  that demonstrates (a) the pre-fix behavior would strand permanently
  (document this explicitly -- e.g. a comment or a before/after assertion
  pair in self-coherence.md, since you cannot literally run "the old code"
  once you've changed it) and (b) the post-fix behavior reaches a valid next
  state with zero human intervention.
- AC3: a fixture mirroring testdata/todo-competing-run.json's shape but with
  run_active: false and branch_exists/pr_exists: true -- confirm todo ->
  in-progress still proposes cleanly (should already pass -- this is a
  regression check per Key finding 2, not new FSM work) and that your
  doctrine-layer changes describe the resume path explicitly.
- AC6: run your new reconciliation path twice against the same fixture;
  assert zero new mutations/comments/PRs on the second pass. Do not
  re-test the finalizer's own idempotence (#591 already covers that) --
  only your new requeue/resume path.

## Guardrails (binding -- violating any of these is a scope violation, not a
judgment call)

- Never delete or discard a cycle/{N} branch or its PR when requeuing --
  every requeue in this cell preserves matter (the #368 protection).
- No new status:* label. No change to what status:review requires. No
  cell_kind enforcement. No #626 (worker-identity/substrate isolation) work.
- Zero behavior change to existing evaluate/scan/terminal tests except the
  ones you deliberately redirect for AC2 -- name exactly which and why.
- Reuse applyStatusLabel, PostComment, and the existing finalizer
  (cn cell finalize) verbatim -- no parallel label-writer, no parallel
  comment-poster, no parallel checkpoint mechanism.
- Full existing #516 (dispatch-repair-preflight) and #524
  (dispatch-closeout-integrity) guard scripts must still pass.
- Write self-coherence.md documenting: your two open-design-call decisions
  (design shape + Go-vs-doctrine-only for AC3) with rationale, exactly
  which pre-existing test(s) you redirected for AC2 and why, and
  confirmation every AC1-AC6 oracle in the scaffold passes.

## When done

Commit, push to cycle/630, append your review-readiness signal to
self-coherence.md, and stop -- you do not dispatch β yourself.
```

## The β prompt

```
Branch: cycle/630

You are β, independently reviewing α's implementation of cnos#630 on
cycle/630. Read .cdd/unreleased/630/gamma-scaffold.md (this scaffold) and
.cdd/unreleased/630/self-coherence.md (α's round record) in full before
forming any verdict. Do not take α's self-coherence claims as verified until
you have independently walked each oracle below yourself.

## Independent AC walk (do not skip any -- re-derive each verdict yourself)

- AC1: find the new reconciliation path. Confirm a fixture with
  in-progress + matter present + no REVIEW-REQUEST.yml + no live run
  produces a concrete next action (not the old permanent no-op) on first
  reconciliation tick.
- AC2: confirm a fixture/test demonstrably tied to the #614 scenario shape
  exists, that the pre-fix behavior is documented as reproducing the strand
  (read self-coherence.md's before/after account critically -- if it just
  asserts this without a mechanical demonstration, that is a finding), and
  that the post-fix run reaches a valid next state with zero
  human-authored label edits or comments anywhere in the test/fixture path.
- AC3: confirm todo -> in-progress claiming still works over a fixture with
  pre-existing branch/PR (regression check against Key finding 2's claim
  that this already worked at the FSM layer -- if α's diff changed the todo
  rules, that itself is a finding requiring explicit justification). Then
  independently read whatever doctrine/code change α made for "resume, not
  defer" and confirm it actually describes resuming from existing matter,
  not merely allowing the label transition (the label transition alone does
  NOT satisfy AC3 -- the resume behavior is what AC3 requires).
- AC4: confirm every new mechanical label reversion posts an audit-note
  comment, and confirm (via fixture or documented reasoning) that a
  subsequent claim/scan pass does not misread that comment's state as an
  unexplained orphan.
- AC5: run go test ./... from src/packages/cnos.issues/commands/issues-fsm
  yourself. Confirm every pre-existing test passes except the ones
  self-coherence.md names as deliberately redirected -- verify those
  redirects are justified, not silent behavior changes. Confirm
  TestSeam_CellKindNotEnforced passes. Confirm #516/#524 guard scripts pass.
  Confirm CI green (Go, Package, Binary + any dispatch-integrity gates).
- AC6: run (or read) the idempotence test yourself -- second pass against
  post-reconciliation state produces zero new mutations.

## Guardrail verification (independent of the AC table)

- Confirm no cycle/{N} branch or PR is ever deleted/discarded by the new
  code path (grep for any git branch -D / PR-close call in the diff).
- Confirm applyStatusLabel/PostComment/cn cell finalize are reused, not
  forked -- grep the diff for any duplicate HTTP label-mutation or
  comment-posting call.
- Confirm no new status:* label was introduced (grep the diff and
  src/packages/cnos.core/labels.json for any new status: entry).
- Confirm the two open-design-call decisions are documented in
  self-coherence.md with stated rationale.

## Verdict

Write .cdd/unreleased/630/beta-review.md with a full R[N] section: outcome
per AC, any findings (with severity), and `verdict: converge` or
`verdict: iterate`. Commit, push, and stop -- you do not re-dispatch α
yourself.
```

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go |
| CLI integration target | `cn` subcommand — extends the existing `cn issues fsm evaluate`/`scan` verb family; no new binary. Any claim-doctrine change is prose (`cds-dispatch/SKILL.md`) rendered via `cn install-wake cds-dispatch` into `.github/workflows/cnos-cds-dispatch.yml` — never hand-edit the rendered workflow file directly. |
| Package scoping | `src/packages/cnos.issues/commands/issues-fsm/` (Go: `table.go`, `scan.go`, new/extended tests + `testdata/*.json`) + `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (declarative rule table) + `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` and `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 (doctrine prose, if the design needs it — see "Open design call") |
| Existing-binary disposition | Preserve — extends `cn issues fsm evaluate`/`scan`; no existing binary replaced, deprecated, or forked |
| Runtime dependencies | None new — reuses the existing GitHub REST client helpers (`ghRequest`/`ghAddLabel`/`ghRemoveLabel`/`ghGetJSON` in `fetch.go`) and the existing `cn cell finalize` self-re-exec pattern |
| JSON/wire contract preservation | Preserve as-is — `FactSnapshot`, `Decision`, `ScanReport` JSON shapes are additive-only (new fields/action strings permitted; no existing field renamed or removed; every pre-existing fixture must still decode and evaluate identically except the ones AC2 deliberately redirects) |
| Backward-compat invariant | Existing `evaluate`/`scan`/`terminal` behavior and all pre-existing tests/fixtures pass unchanged except the AC2-redirected ones (named explicitly in `self-coherence.md`); `transitions.json`'s `todo` state and every non-`in-progress` state's rules are untouched unless α's diff explicitly justifies otherwise |

No TBD rows — every axis above is pinned from direct repo-convention evidence gathered during peer enumeration (Key findings 1–2 and the Source-of-truth table); none required escalation to δ beyond this single-session bootstrap-δ scaffold itself.

## Scope guardrails (binding on α)

- Every requeue/resume this cell introduces MUST preserve the `cycle/{N}` branch and its PR — never discard matter (the #368 failure this issue's own STOP condition names).
- No new `status:*` label; no change to `status:review`'s deliverable-proof bar; no Demo 0; no `cell_kind` enforcement; no #626 (worker-identity/substrate isolation) work.
- All label mutation routes through the existing `applyStatusLabel` primitive; all comments through the existing `PostComment` shape; all checkpointing through the existing `cn cell finalize` — no parallel/duplicate mechanism.
- `#516`/`#524` dispatch guard scripts (`check-dispatch-repair-preflight.sh`, `check-dispatch-closeout-integrity.sh`) must stay green.
- Any rendered-workflow change (`.github/workflows/cnos-cds-dispatch.yml`) must go through `cn install-wake cds-dispatch`, never a hand-edit (the `install-wake-golden` CI gate enforces this).

## Friction notes

1. **AC2's "must fail before the fix" requirement is mechanically awkward for a single-branch cycle.** Standard TDD red/green requires running the *old* code against the new fixture and observing failure, but α is implementing the fix in the same cycle branch — there is no separate "pre-fix" binary to run the new fixture against without a throwaway commit/stash. I recommend α either (a) commit the new fixture + a test asserting the *old* wedge behavior first, observe it pass (proving the fixture faithfully reproduces the strand against pre-fix logic), then implement the fix and update the test's assertion to the new behavior in a follow-on commit (both commits visible in the cycle branch history as the red/green pair) — or (b) if a single-commit implementation is preferred, explicitly document in `self-coherence.md` a manual verification step (e.g. `git stash` the fix, run the new test, observe failure, `git stash pop`) with the failure output pasted in. Left as α's call since the scaffold should not mandate git-choreography beyond what `alpha/SKILL.md` already governs.

2. **The claim-doctrine (AC3) half of this cell may end up being pure prose with no new Go test surface**, per the second open design call. If α confirms (as I recommend) that `cn issues fsm evaluate --json` already exposes everything the claim sequence needs, then AC3's oracle is partly a documentation-quality check, not a code-coverage check — β should not force a fabricated Go test for a doctrine-only change; the AC oracle table already accounts for this ("If no Go-testable surface exists for the doctrine-only half, the oracle is the doctrine text itself").

3. **`origin/main` was re-pinned mid-claim** (δ's STOP decision, restated at the top of this scaffold) from the claim-time SHA (`313cb1f3`) to the current SHA (`6143b53c`) two commits later. Both trailing commits (`8267b5d3` admin-wake heartbeat, `6143b53c` board-map regen) are confirmed docs/log-only and touch nothing this cell's diff graph intersects — if α's diff touches `docs/development/board/**` or `.cn-sigma/logs/**`, treat that as an unrelated pre-existing/concurrent-automation file, not part of this cell's scope, exactly per the cnos#615 scaffold's Friction note 4 precedent for the same class of drift.

4. **This is a bootstrap-δ (single-session δ-as-γ), not wake-invoked-δ, cycle.** Per `delta/SKILL.md` §9.1's own table, wake-invoked-δ (the production mode this issue's fix is ultimately *for*) is "not yet observed in production" — meaning this very cycle is being driven the same ad hoc way #614's resume was (a human/δ session improvising continuity), which is itself a small piece of evidence for why AC3's mechanical, wake-invoked-safe resume path is worth landing: it turns today's improvisation into a named, repeatable mechanism.
