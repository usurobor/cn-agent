# γ close-out — cnos#630

## Scope note

This file is written at the **δ §9.5 converge boundary**: β returned
`verdict: converge` at R0 (`beta-review.md`, head `94dc605c`), and δ's
§9.5 per-R[N] artifact contract's `converge` row requires all three
closeout artifacts (`alpha-closeout.md`, `beta-closeout.md`,
`gamma-closeout.md`) to land on `cycle/630` before δ writes the
`status:review` return token. This is **not** the full γ §2.7/§2.10
closure declaration — that closure gate (post-release assessment, cycle-
directory move to `.cdd/releases/{X.Y.Z}/{N}/`, `RELEASE.md`, hub-memory
update, `gh issue view {N}` close-state assertion, merged-branch cleanup,
δ release-boundary preflight) presumes a merge and a release boundary
that have not happened yet — no PR exists on `cycle/630` at the time of
this writing, `main` is untouched by this cycle, and δ's next step per
§9.6 is to open the cycle-PR and request the `in-progress → review`
transition. What follows is the **process-gap audit** portion of γ's
role at this boundary: did the cycle reveal friction worth converting
into a follow-up, or not.

## Cycle summary

- **Issue:** cnos#630 — "partial-matter in-progress wedge: give the
  reconciler a mechanical path out."
- **Shape:** single R0 pass, `run_class: first_pass`. γ scaffolded
  (`gamma-scaffold.md`, commit `5b9aaa2c`), α implemented and signaled
  review-readiness once (implementation commit `79fc8565`, readiness
  signal at `a63bc9d2`), β reviewed once and converged
  (`beta-review.md §R0`, commit `94dc605c`). No `iterate` round, no RC,
  no re-dispatch of α mid-cycle.
- **Result:** a new declarative `transitions.json` rule
  (`propose_status_todo_with_matter`, `in-progress` state) that requeues
  a dead-run-with-checkpointed-matter cell to `status:todo` while
  preserving its `cycle/{N}` branch and PR (the #368 blind-requeue
  protection restated for the reconciler's own action), reusing
  `scan.go`'s existing generic reconcile branch with zero new Go
  control-flow; a new fixture-backed regression lock proving the
  pre-fix rule shape reproduced the #614 strand
  (`TestAC630_WedgePreFixRuleReproducesStrand`); a new
  `delta/SKILL.md` §9.11 doctrine section ("resumed-from-mechanical-
  reversion shape") plus a `cds-dispatch/SKILL.md` claim-mechanism
  paragraph and reordered Step A classification, so a claiming wake
  treats pre-existing branch/PR matter as a normal resume shape rather
  than an orphan or a rejected-repair candidate. Full pre-existing
  `issues-fsm` suite (77 tests) passes with only two deliberately
  redirected/updated; `transitions.json`'s `todo` state and
  `src/packages/cnos.core/labels.json` are byte-identical to
  `origin/main`. Live GitHub Actions CI green on the review head SHA
  (`a63bc9d2`) — 11/11 checks `success`, independently pulled by β via
  `gh api`.

## R0-only convergence: a positive signal, not a gap

This cycle converged at R0 with no iteration. Per the review-churn
trigger in `gamma/SKILL.md` §2.8 (fires only above 2 review rounds), a
single clean round is the target condition, not an anomaly to explain
away. Two factors worth naming for the record, not as findings
requiring action:

- The scaffold pre-resolved both open design calls with a concrete
  recommendation and rationale before α started (§"Open design call"),
  and α's `self-coherence.md` §ACs "Open design calls" subsection
  resolved both by direct appeal to the actual guard vocabulary
  (`pr_exists` as the literal boolean distinguishing "checkpointed" from
  "not yet checkpointed"), not by restating the scaffold's abstract
  argument. β independently re-derived both decisions from source
  (`beta-review.md` "Guardrail verification" — "substantive engineering
  rationale, not restated scaffold text") rather than accepting them on
  narrative trust.
- The scaffold's Key finding 2 (the FSM's `todo` state already permits
  claiming a cell with pre-existing matter — the gap lives in the
  claim-*doctrine* prose, not the FSM) was independently reconfirmed
  three separate times across the cycle: by α (`self-coherence.md`
  §ACs AC3, "no `todo`-state rule change was needed or made"), and by
  β (`beta-review.md` AC3, "the `todo` state's 3 rules are byte-for-byte
  unchanged in the diff"). Three independent confirmations of the same
  boundary in one cycle is the discipline working as intended — nobody
  had to trust the scaffold's up-front claim, everyone re-derived it.

## Friction audit — the two frictions the scaffold flagged, checked against the final state

γ's scaffold logged four friction notes at R0. The two substantive ones
(mechanical/procedural risk, not informational) are audited here in
detail per this closeout's dispatch instruction; the other two
(doctrine-only AC3 oracle allowance; the base-SHA re-pin) are covered
in the triage table below.

**1. AC2's "must fail before the fix" requirement — mechanically awkward
for a single-branch cycle, resolved cleanly.** The scaffold's Friction
note 1 named the real constraint: a single cycle branch cannot literally
run "the old binary" against a new fixture without a throwaway
commit/stash, and left the resolution mechanism to α's judgment (option
(a), a two-commit red/green pair, vs. option (b), a documented manual
stash-and-observe). α chose neither literally — instead,
`TestAC630_WedgePreFixRuleReproducesStrand` constructs an inline `Table`
literal carrying only the verbatim pre-diff rule shape (copied from the
diff's own removed lines) and asserts it reproduces the wedge
(`TargetState == ""`) against the real #614 fixture. This is a third
option the scaffold did not anticipate, and it resolved the friction
better than either named option: unlike (a) it needs no throwaway
commit pair to read the history for, and unlike (b) it is not a
one-time manual verification pasted into a document — it is a
permanent, committed regression lock that fails loudly forever if the
old rule shape is ever reintroduced. β treated this with explicit
skepticism per its own review discipline ("checked with real
skepticism, per the scaffold's own friction-note flag") and
independently diffed the inline literal against the diff's removed
lines to confirm it was not a strawman, then re-ran it (`PASS`). **This
friction resolved cleanly; no follow-up needed** — but the resolution
pattern (a literal `Table` snapshot of pre-diff rule shape as the
red-side oracle, for a single-branch TDD-red-before-fix requirement on
declarative-table-driven code) is reusable enough that a future γ
scaffold facing the same class of constraint could name it as a
candidate mechanism directly, rather than re-deriving it from first
principles. Not filing a skill patch for this now — one occurrence is
not yet a pattern per `gamma/SKILL.md` §2.9's own bar ("did this cycle
reveal a *recurring* friction"); noting it here so a second occurrence
would be recognized as recurring rather than novel.

**2. The base-SHA re-pin (δ's STOP decision, restated in the scaffold's
header and Friction note 3) — handled correctly, no drift materialized.**
The scaffold recorded δ's re-pin from the claim-time SHA (`313cb1f3`) to
`origin/main`'s SHA two commits later (`6143b53c`), with an explicit
verification that both trailing commits touched only
`.cn-sigma/logs/**` and `docs/development/board/**` — no overlap with
this cell's diff graph. `origin/main` then advanced *again* during α's
own round (to `e0590148`, three more merges/commits); α's
`self-coherence.md` §Gap independently re-ran the same overlap check
against the new range (`git log --oneline 6143b53c..origin/main -- <this
cell's five package paths>` returning zero lines) before rebasing and
force-pushing, and re-validated zero further drift immediately before
the review-readiness signal (§Review-readiness round 1: "`cycle/630`
still contains `origin/main`'s current tip; no further drift occurred").
β independently re-derived the branch's merge-base against
`origin/main` at review time and confirmed zero drift there too
(`beta-review.md` §R0 header). This is the base-SHA re-pin protocol
(scaffold pin → α re-verify at every rebase → β re-verify at review)
working exactly as the `alpha/SKILL.md` §2.6 row-1 rebase discipline
intends, across two separate drift events in one cycle, with zero
missed overlap. **No follow-up needed** — this is a positive
confirmation of an existing gate, not a gap.

## Disclosed debt — triaged, not filed as a tracking issue

`self-coherence.md` §Debt names five items; `beta-review.md` §Findings
independently re-confirms the first as accurate and adequately scoped.
The one item this closeout's dispatch instruction specifically asks γ
to weigh: **§9.11's doctrine addition (the new `delta/SKILL.md` section
and the reordered `cds-dispatch/SKILL.md` Step A) has no live
wake-invoked-δ empirical witness yet** — this cycle itself, like the
§9.10 precedent it mirrors, ran under bootstrap-δ (single-session
δ-as-γ), and per `delta/SKILL.md` §9.1's own table, wake-invoked-δ is
"not yet observed in production" as of this cycle.

**Disposition: explicit disclosed limitation, not a tracking issue.**
Reasoning:

- This is structurally identical to §9.10's own history — §9.10's first
  and, as of this cycle, still-only empirical witness (cycle/497) was
  *also* executed via a bootstrap exception, not a genuine wake firing.
  §9.10 has stood as accepted doctrine since that cycle without a
  dedicated tracking issue demanding a live-firing witness; treating
  §9.11 differently (filing an issue that §9.10 itself never got) would
  be an inconsistent bar applied only because this cycle happens to be
  the one authoring the closeout, not because the debt is qualitatively
  different.
- A tracking issue's own first AC would necessarily be "wait for an
  actual production wake-invoked-δ firing to claim a cnos#630-resumed
  cell" — that is not a piece of work a cycle can execute; it is an
  observation that will happen automatically the first time this
  scenario recurs in production, at which point the witness is free
  evidence, not commissioned work. Filing an issue whose only AC is
  "observe a future event" is process overhead without a concrete next
  action, contrary to `gamma/SKILL.md` §3.7 ("'Noted' is not a
  disposition") — the correct disposition *is* "noted, with a named
  reason," not "filed," because there is nothing to file that isn't
  already satisfied by §9.11's own text plus this closeout's record of
  the gap.
- The gap is genuinely bounded: β independently confirmed AC3's oracle
  text explicitly allows a doctrine-only resolution when no
  Go-testable surface exists ("the resume behavior... is what AC3
  requires" — clearing the bar on mechanism-naming, not on
  live-firing proof), so the missing witness is not a finding against
  this cycle's AC3/AC4 closure, only a forward-looking observation.

**If a second doctrine-only §9.x section reaches this same boundary
with the same "no live witness" debt before either #9.10 or #9.11 gets
one, that would be the recurring-friction signal `gamma/SKILL.md` §2.9
asks γ to watch for** — worth naming here so a future γ pass has the
prior occurrence on record rather than re-discovering it from scratch.

## Cycle-iteration triggers (`gamma/SKILL.md` §2.8) — none fired

- **Review churn** (>2 rounds): did not fire — 1 round.
- **Mechanical overload** (>20% mechanical findings, ≥10 total): did not
  fire — 0 findings.
- **Avoidable tooling/environment failure**: did not fire in a blocking
  sense. α's own sandbox lacked a GitHub Actions runner (disclosed,
  not concealed, in `self-coherence.md` §Review-readiness); β's
  independent `gh api` pull against the review head SHA closed that gap
  before convergence, exactly as the review process is designed to
  compensate for a disclosed, real environment limitation.
- **Loaded-skill miss**: none identified.

Per `gamma/SKILL.md` §2.9's independent process-gap check: no recurring
friction beyond the single-occurrence AC2-oracle pattern named above (not
yet a recurring pattern by its own two-occurrence bar), no gate found
too weak or vague, no role skill failed to prevent a predictable error,
and no coordination burden suggested a better mechanical path in this
cycle.

## Triage table

| Finding | Source | Type | Disposition | Artifact / commit |
|---|---|---|---|---|
| AC2 "must-fail-before-fix" red/green choreography on a single branch | γ scaffold Friction note 1 | process/mechanical | resolved in-cycle by α (inline pre-diff `Table` literal as the red-side oracle); pattern noted for reuse, not yet filed as a skill patch (single occurrence) | `self-coherence.md` §ACs AC2; `issuesfsm_test.go` `TestAC630_WedgePreFixRuleReproducesStrand`; this file §Friction audit item 1 |
| Base-SHA re-pin across two drift events (claim-time → scaffold-time, scaffold-time → review-readiness-time) | δ STOP decision (scaffold header); α self-coherence.md §Gap | coordination | closed — both drift events verified zero-overlap and re-validated at every checkpoint; no follow-up | `gamma-scaffold.md` header; `self-coherence.md` §Gap, §Review-readiness round 1 |
| Two open design calls (rule-shape; Go-vs-doctrine-only) | γ scaffold | design | resolved in-cycle by α with rationale tied to actual guard vocabulary, independently re-verified by β | `self-coherence.md` §ACs "Open design calls"; `beta-review.md` "Guardrail verification" |
| §9.11 doctrine addition has no live wake-invoked-δ empirical witness | α self-coherence.md §Debt item 1; β beta-review.md §Findings | doctrine/process | disclosed limitation, not filed — mirrors §9.10's own unresolved-witness history; witness is a future free observation, not commissionable work | `self-coherence.md` §Debt item 1; `beta-review.md` §Findings; this file §"Disclosed debt" |
| Comment-text-based ("MECHANICAL reversion" phrase) audit-note detection, not structured data | α self-coherence.md §Debt item 2 | design (deliberate, scoped) | drop — explicitly a scoped choice consistent with the "no new mechanism" guardrail; named as a candidate follow-up only if it proves fragile in practice, not built now | `self-coherence.md` §Debt item 2 |
| "Resumed, then died again mid-cycle" not independently fixture-tested (argued-safe via `Evaluate`'s purity) | α self-coherence.md §Debt item 3 | test-coverage | drop — low risk, structurally argued, disclosed rather than silently assumed; no action requested by β | `self-coherence.md` §Debt item 3 |
| "Only one pre-existing fixture matches the new rule's guard combination" — point-in-time audit, not an enforced invariant | α self-coherence.md §Debt item 4 | documentation | drop — correct and intended behavior as new fixtures are added; named so a future reader does not misread drift as a defect | `self-coherence.md` §Debt item 4 |

## What γ is not asserting at this boundary

- No claim about the parent issue's close state (`gh issue view 630
  --json state`) — the issue is not yet closed, and closure only
  follows a merge that has not happened. The `gamma/SKILL.md` §2.10
  row-15 issue-close assertion is a full-closure-gate item, not a
  converge-boundary item; it belongs to a later γ pass after δ merges.
- No post-release assessment (PRA) is written here — `gamma/SKILL.md`
  §2.7 makes the PRA a post-merge artifact measuring released work;
  nothing has released yet.
- No hub-memory update, cycle-directory move, or `RELEASE.md` — all
  are full-closure-gate items (`gamma/SKILL.md` §2.10), out of scope
  for the §9.5 converge boundary.
- No PR has been opened for `cycle/630` (confirmed via `gh pr list
  --head cycle/630` returning empty at the time of this writing) —
  opening the cycle-PR and requesting the `status:review` transition is
  δ's next action per §9.6, not this pass's.

## Next step

Per `delta/SKILL.md` §9.5, all three converge-boundary closeout
artifacts now exist on `cycle/630`. δ's next action is to open the
cycle-PR, ensure `REVIEW-REQUEST.yml` and the closeout matter exist
(this pass also writes `.cdd/unreleased/630/REVIEW-REQUEST.yml`), and
request the `status:review` transition via `cn issues fsm evaluate
--issue 630 --apply` per §9.6 — γ is not requesting that transition
itself and is not touching issue #630's labels or comments in this
pass.
