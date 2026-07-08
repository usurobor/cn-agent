# α close-out — cnos#630

## Scope note

This close-out is written at the **δ §9.5 converge boundary** — β
returned `verdict: converge` at R0 (`beta-review.md`, head
`94dc605c`), and this file is one of the three artifacts that boundary
requires (`alpha-closeout.md` + `beta-closeout.md` + `gamma-closeout.md`),
per `delta/SKILL.md` §9.5's per-R[N] artifact contract for the
`converge` row. It is **not** the post-merge/post-release close-out
`alpha/SKILL.md` §2.8 otherwise describes (that flow assumes β has
already merged to `main`; here, merge/PR/release are still ahead —
δ's next step, per §9.6, is to open the cycle-PR and request the
`status:review` transition). The retrospective content below is
written against the full R0 cycle record (`gamma-scaffold.md`,
`self-coherence.md`, `beta-review.md`) as it stands at this boundary,
all on `cycle/630`, none of it yet on `main`.

## Summary — what was built

A mechanical exit from the "partial-matter in-progress wedge"
(cnos#630, concrete evidence cnos#614): a `status:in-progress` dispatch
cell whose run died after producing matter (a `cycle/{N}` branch +
checkpointed draft PR, per the cnos#591 finalizer) but before
`REVIEW-REQUEST.yml` was written previously stalled forever — it could
neither requeue (the cnos#575/#368 blind-requeue guard blocks a
requeue over existing matter), advance (no `REVIEW-REQUEST.yml`), nor
be re-claimed (the dispatch selector only claims `status:todo`). The
literal wedge was `scan.go`'s `else` branch of the
`propose_delta_recovery` case (pre-diff `scan.go:252-256`), which
re-evaluated to the identical permanent no-op on every scan tick once
matter was checkpointed.

The fix is a single new declarative rule in `transitions.json`
(`in-progress` state, `all_false: [run_active]` + `all_true:
[pr_exists]` → action `propose_status_todo_with_matter`, target state
`todo`), positioned ahead of the existing dead-with-matter
`propose_delta_recovery` rule and gated on the exact boolean
(`pr_exists`) that distinguishes "matter already checkpointed" from
"matter not yet checkpointed." The new action flows through `scan.go`'s
pre-existing generic reconcile branch (the same one Cases 2/5 already
used) with zero new Go control-flow — no new function was added to
`scan.go` at all. Two doctrine-prose additions close the claim-time
half of the gap: a new paragraph in `cds-dispatch/SKILL.md`'s claim
mechanism, and a new `delta/SKILL.md` §9.11 ("resumed-from-mechanical-
reversion shape," sibling to §9.10) naming the trigger, detection
mechanism, and routing sequence for a wake/δ session that claims a
`status:todo` cell carrying pre-existing `cycle/{N}` matter.

## The two open-design-call resolutions, restated with outcome

**Design call 1 — extend the existing rule in-place, or add a distinct
new action/rule?** The γ scaffold recommended (ii), a distinct rule, on
an abstract argument (the two states are "observably different").
Resolved as (ii), `propose_status_todo_with_matter`, but the rationale
actually written into `self-coherence.md` is concrete, not a restatement
of the scaffold: `pr_exists` is the literal boolean that flips the
instant `cn cell finalize` checkpoints matter into a draft PR; before
that flip, `branch_has_commits` may already be true with `pr_exists`
still false (matter exists, checkpointing still needed). Collapsing the
two states into one rule would force `scan.go` to grow a conditional
inside a single case deciding "call finalize" vs. "just relabel" —
reintroducing exactly the state-name-shaped Go branching `table.go`'s
own doc comment forbids ("never switched on a CDS-specific state
name... entirely table-driven"). Two rules, ordered by `pr_exists`,
keep the decision entirely declarative. **Outcome:** β independently
re-derived this same argument from the diff rather than accepting the
narrative (`beta-review.md` "Guardrail verification": "substantive
engineering rationale, not restated scaffold text") — the resolution
held without amendment.

**Design call 2 — does AC3's claim-time resume need new Go surface, or
is it doctrine-only?** The scaffold recommended doctrine-only,
predicting `cn issues fsm evaluate --json` already exposes everything
the claim sequence needs. Confirmed by a concrete negative check, not
just agreement: `transitions.json`'s `todo` state rules (read in full)
gate the claim rule on `claim_request_present` + `all_false:
[run_active]` only — no `todo`-state guard references
`branch_exists`/`pr_exists`/`cdd_artifacts` at all, so there was never
anything for the claim path to be missing at the data layer. The actual
gap was that the claim-doctrine prose never *named* "claimed cell
already has branch/PR matter" as an expected, resume-worthy shape —
a claiming wake or δ session reading the doctrine cold had no textual
instruction to treat it as normal rather than an anomaly. This
resolved as pure doctrine (`cds-dispatch/SKILL.md` + `delta/SKILL.md`
§9.11), zero new Go. **Outcome:** β independently confirmed the `todo`
block of the diff hunk is byte-for-byte unchanged and independently
read §9.11's routing sequence in full (not from the self-coherence.md
summary), confirming it names a trigger, a detection method, and a
concrete 3-step routing sequence rather than vague prose — the "label
transition alone does NOT satisfy AC3" bar the scaffold set was cleared
on independent re-derivation, not narrative trust.

## What worked well

- **The scaffold's Key finding 2** (peer-enumeration-derived: the FSM's
  `todo` state already permitted claiming a cell with pre-existing
  matter, so AC3's real gap was one layer up in prose) meant α did not
  have to discover this from scratch — it was verified, not
  re-derived, saving a full investigation pass. This is the
  `gamma/SKILL.md` §2.2a peer-enumeration discipline paying off
  concretely: γ read `scan.go` and `transitions.json` in full at
  scaffold time rather than asserting the gap from the issue text
  alone.
- **The existing near-duplicate fixture and test**
  (`testdata/scan-died-after-pr-before-review-request.json` /
  `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment`) being the
  literal #614 shape, already committed on `main`, meant AC2's fixture
  work was a redirect of existing coverage rather than authoring a new
  scenario from zero — the scaffold correctly identified this and told
  α to treat it as the base to extend, not duplicate.
- **`scan.go`'s existing generic reconcile branch being written broadly
  enough to absorb a third case for free** (`case dec.Outcome ==
  "proposed" && dec.TargetState != ""`) meant the entire fix required
  zero new Go functions in `scan.go` — only a new data row in
  `transitions.json` and a doc-comment update in `table.go`. This is a
  direct payoff of the table-driven design discipline the codebase
  already enforces.

## What was harder than expected

**The mid-round rebase.** `origin/main` advanced during α's own round —
by three merges/commits (PR #631/cycle/612, PR #624/cycle/611, and an
unrelated "cn CLI ergonomics" commit `8db8ff4b`) beyond the SHA δ had
already re-pinned once at scaffold time. This meant re-running the
overlap check α had already planned to run once, a second time, against
a wider commit range, and force-pushing after a clean rebase — not
difficult mechanically (the overlap check returned zero lines against
every path this cell's diff touches), but it meant the "base SHA"
narrative in `self-coherence.md` §Gap had to carry two separate re-pin
events (δ's original claim-time-to-scaffold-time re-pin, then this
round's scaffold-time-to-review-readiness-time rebase) rather than one,
and every downstream evidence claim (test counts, `cn cdd verify`
output, the golden-file sha256 comparison) had to be re-verified
post-rebase rather than assumed stable. This added real time to the
round without changing the shape of the implementation itself.

**The AC2 red/green choreography.** The scaffold's own Friction note 1
named this honestly as mechanically awkward on a single cycle branch:
standard TDD red-then-green requires running the *old* code against the
*new* fixture and observing failure, but there is no separate "pre-fix"
binary to run the new fixture against without a throwaway
commit/stash, since the fix lands in the same branch as the fixture
that tests it. The scaffold offered two named options and left the
choice to α. Neither named option was actually used — the resolution
that shipped is a third path: an inline `Table` literal in
`issuesfsm_test.go` (`TestAC630_WedgePreFixRuleReproducesStrand`)
carrying a verbatim copy of the pre-diff rule shape, run against the
real fixture, asserting the wedge (`TargetState == ""`) reproduces.
This took more thought than either scaffold-offered option to arrive
at, but produces a strictly better artifact than either: a permanent,
committed regression lock rather than a one-time manual verification
(option (b)) or a two-commit history someone has to go read to find the
red/green pair (option (a)). β checked this construction with explicit
stated skepticism ("checked with real skepticism, per the scaffold's
own friction-note flag") and independently diffed the inline literal
against the diff's own removed lines before accepting it as a faithful,
non-strawman reproduction.

## Debt carried forward — status at this boundary

`self-coherence.md` §Debt named five items. Status now:

1. **"§9.11's doctrine addition has no live wake-invoked-δ empirical
   witness yet"** — **not resolved, not expected to be at this
   boundary**, by design. This cycle itself ran under bootstrap-δ (the
   same production mode gap this fix targets), mirroring §9.10's own
   history (its only witness, cycle/497, was also a bootstrap
   exception). A genuine wake-invoked-δ firing claiming a
   cnos#630-resumed cell is the empirical anchor §9.11 currently lacks;
   γ's closeout addresses whether this rises to a tracking issue.
2. **"comment-text-based audit-note detection, not structured data"**
   — **not resolved, and not expected to be**, by explicit scope
   decision (Design call 2: doctrine-only, no new mechanism). Named as
   a candidate follow-up only if it proves fragile in practice.
3. **"resumed-then-died-again mid-cycle not independently fixture-
   tested"** — **not resolved**; argued safe via `Evaluate`'s purity
   (same rule, same guards, same outcome on a second dead-run event)
   but not independently exercised by a fixture chain. Disclosed, not
   silently assumed. β raised no objection to this disclosure.
4. **"only one pre-existing fixture matches the new rule's guard
   combination — point-in-time audit"** — **expected to become stale
   as new fixtures are added**, by design; named so a future reader
   auditing `testdata/` does not misread a changed count as drift or a
   defect.
5. **"no PR opened for `cycle/630` yet"** — **still true at this
   boundary**, and correctly so: PR-opening and the
   `status:in-progress → status:review` transition request are δ's
   next actions per `delta/SKILL.md` §9.6, not something either α or β
   does.

No new debt surfaced during β's review beyond what `self-coherence.md`
already disclosed — `beta-review.md` §Findings reports zero blocking
findings and explicitly treats the §9.11-witness gap as "a legitimate
scope boundary, not a gap this cycle should have closed."

## Calibrated success claim

What is verified, as of this boundary:

- All six scaffold ACs (AC1–AC6) pass, independently re-derived by β
  from the code/tests/diff — every test cited in `beta-review.md` was
  re-run by β from a clean shell, not read off α's pasted output.
- The full pre-existing `issues-fsm` test suite passes unchanged (77
  tests) alongside 4 new `TestAC630_*` functions and 1 new coverage-
  preservation test, with exactly two pre-existing tests deliberately
  redirected/updated (named explicitly in both `self-coherence.md` and
  independently re-confirmed in `beta-review.md`).
- `transitions.json`'s `todo` state and `src/packages/cnos.core/
  labels.json` are both byte-identical to `origin/main`.
- Live GitHub Actions CI is green on the actual cycle-branch head
  commit at β's review time (`a63bc9d2`), independently confirmed by β
  via `gh api` — 11/11 checks `success`.
- The `.github/workflows/cnos-cds-dispatch.yml` rendered-workflow
  change was independently reproduced by β from source (rebuilding the
  renderer and re-running it against the checked-out tree), not merely
  sha256-compared against an already-committed golden — a stronger
  check than α's own self-coherence.md evidence performed.

What is **not yet verified**, because it is not this boundary's job:

- No PR has been opened and no merge to `main` has happened — that is
  δ's next step (§9.6 `status:review` transition), not something
  either α or β has done or claims to have done.
- CI has not been (and cannot yet be) observed on a post-merge `main`
  commit, since no merge exists yet.
- §9.11's doctrine has not been exercised by a real wake-invoked-δ
  firing (see Debt item 1 above).

**Claim:** the implementation is correct and complete against the
scaffold's six ACs and all stated guardrails, as verified by two
independent passes (α's own self-coherence walk and β's from-scratch
re-derivation including live CI). It is ready for δ to proceed to the
`status:review` transition. It is not yet released, merged, or tagged,
and no claim is made about post-merge CI or a live wake-invoked-δ
firing beyond what is stated above.
