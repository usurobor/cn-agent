# γ close-out — cycle/574

Written per `gamma/SKILL.md` §2.7 (close-out / process-gap audit / triage), adapted to this
cycle's actual state: δ has not yet opened the PR or requested the `status:review` transition,
so the post-merge steps of §2.7 (CI-on-merge-commit verification, cycle-directory move to
`.cdd/releases/`, POST-RELEASE-ASSESSMENT.md) do not apply yet and are out of scope for this
artifact. This close-out covers the process-gap audit and the F1/F2 triage that are available
now, and declares the cycle complete pending δ's PR + transition request.

## Close-outs collected

- `.cdd/unreleased/574/alpha-closeout.md` — present, on branch.
- `.cdd/unreleased/574/beta-closeout.md` — present, on branch.
- `.cdd/unreleased/574/beta-review.md` — present, verdict **APPROVED**, round 1 (R0).
- `.cdd/unreleased/574/REPAIR-PLAN.md` — δ's repair-reentry note; read in full; informs the
  process-gap audit below.

No re-dispatch of α or β was required to obtain either close-out — both were already on the
branch at the start of this firing.

## Process-gap audit

**What went right.**

- The γ scaffold (`gamma-scaffold.md`) was substantive enough that α needed **zero mid-flight
  clarification**: self-coherence.md §Gap records that the issue's only comment was an operator
  note unrelated to AC scope, and α's own §Self-check section shows every ambiguity the scaffold
  anticipated (AC2's internal tension, AC4's design call, AC5's edit-vs-append question) was
  resolved by α using the scaffold's own reasoning, not left open for β to discover.
- The scaffold's Friction notes explicitly named the exact trap AC2/AC3 contained: that tightening
  `any_true → all_true` is not a mechanical rewrite of the paired blocked/negative rule, because the
  literal complement of `all_true:[g1,g2,g3]` is "NOT all_true" (any one guard false), not
  `all_false` (all three false) — and that a naive `all_false` rewrite would let a partial-evidence
  state fall through to the evaluator's generic, non-diagnostic "no rule matched" fallback instead of
  an explicit `blocked` outcome with a populated `evidence_guards` list. α applied this correctly in
  both AC2 and AC3 (the "unconditional second rule" design), and β independently traced the fix
  through `table.go`'s `ruleMatches`/`Evaluate` rather than accepting the claim as asserted, confirming
  it actually behaves as documented. This is the scaffold doing its job: naming the trap before
  anyone could fall into it, in exactly the place the issue's own guard-tightening work required it.

**What went wrong.**

- The cell **stranded once, before β ever ran**. Per `REPAIR-PLAN.md`'s timeline reconstruction
  (`gh api repos/usurobor/cnos/issues/574/timeline`): α's implementation completed and pushed at
  `897faad1` (02:16:24Z, with `self-coherence.md` already carrying full `MET` writeups for AC1–AC7
  and both AC5 external artifacts already live), but the label history shows
  `in-progress → todo (03:38:42Z, release-back-to-queue) → in-progress (03:41:33Z, next wake's claim)`
  with **no `status:review`/`status:changes` event ever occurring in between**. In other words: the
  substrate/workflow firing that ran α completed its own work and pushed a fully review-ready diff,
  but that same firing ended (or was cut off) before β was ever dispatched against it, and the
  dispatch wake's own release-back-to-queue mechanism reclassified the cell back to `todo` rather
  than leaving it silently claimed-and-abandoned. A second wake firing then had to run the full
  dispatch-wake **claim → release-back-to-queue-recognition → repair-reentry preflight → reclaim**
  cycle (`REPAIR-PLAN.md`) just to discover that no repair was actually needed — the branch already
  held a complete, unreviewed diff — and resume straight to dispatching β.

  Naming this precisely as the process gap: **α→β dispatch continuity is not guaranteed within a
  single firing.** There is no mechanism that makes "α finished and pushed" atomically imply "β gets
  dispatched next" — the boundary between those two steps is a substrate/workflow lifetime edge, and
  if that edge is hit, the cell strands until some later wake polls, notices the stranding, and
  reclaims it. This is exactly the class of failure this cell's own **AC4/#368 concern** is about
  (a dead run whose matter exists but whose completion isn't correctly observed must resolve to
  recovery, never a blind requeue-from-scratch) and exactly what the **dispatch-protocol's
  release-back-to-queue mechanism** exists to catch at the orchestration layer — applied here
  reflexively to this very cell's own lifecycle, not to the FSM guards the cell was hardening.
  The good news: the mechanism worked as designed on both sides of that reflexive application —
  release-back-to-queue put the cell back in a discoverable state instead of losing it silently,
  and the reclaiming wake (per `REPAIR-PLAN.md`) read the branch's actual content (commit history,
  `self-coherence.md`'s AC-by-AC evidence, the two live GitHub artifacts for AC5) before acting,
  correctly recognizing "matter already exists, resume/continue" rather than re-scaffolding or
  re-implementing from `todo`. That is precisely the delta-recovery outcome AC4 hardens `fetch.go`
  to reach for a dead run with remote matter, rather than the `propose_status_todo` blind-requeue
  #368 named as the failure mode — observed here operating correctly one layer up, at the
  dispatch-wake level, using the same principle (check matter before concluding "nothing happened").
  No new tracking issue is filed for this gap: it is not a defect introduced by this cycle, the
  existing release-back-to-queue + repair-reentry preflight mechanism handled it correctly end to
  end, and it stands as a second empirical anchor (alongside `gamma/SKILL.md` §2.5's cycle #369
  anchor for the scaffold-presence gate) that this recovery path works under a real, unplanned
  stranding — not just in the abstract.

## Triage: F1 and F2 (from `beta-review.md` §Findings)

| Finding | Source | Type | Disposition | Artifact / commit |
|---|---|---|---|---|
| **F1** — `TestAssembleLive_RemoteOnlyBranchResolvesToDeltaRecovery`'s name implies an `assembleLive`-level integration test, but the body calls `observeRemoteBranch` directly and hand-constructs a `FactSnapshot`, never invoking `assembleLive` itself. Disclosed by the function's own doc comment; β independently read the actual 5-line call site in `fetch.go` and confirmed it correct by direct inspection. | β (`beta-review.md` §Findings; echoed in `alpha-closeout.md` §Non-blocking findings) | One-off — cosmetic test-name imprecision, no hidden behavior claim | **Drop explicitly.** No tracking issue: the drift is name-only, the adjacent doc comment already discloses it honestly, and the actual call site has been independently verified correct by direct reading (not by this test). A future touch of `issuesfsm_test.go` may rename it to something like `TestObserveRemoteBranch_ComposesIntoDeltaRecovery` for precision, but this is not worth a standalone issue for a single test name. | n/a (recorded here + in `beta-review.md`/`alpha-closeout.md`) |
| **F2** — `cn-cdd-verify`'s `classifyCycleType` (`ledger.go:425`) classifies a cycle as "small-change" whenever `beta-review.md` doesn't yet exist, without checking `gamma-scaffold.md`'s presence as an alternative triadic signal; the small-change path's stricter, literal `"## Gap"`-prefix section check then hard-fails on this repo's own `§`-prefixed `self-coherence.md` convention, producing a spurious CI-red on I6 for every commit between α's push and β's first push. Reproduced live (CI run `28693885842` + 7 preceding), mechanism confirmed via a merge-tree stub-file test, self-resolving once β's own artifacts land. | β (`beta-review.md` §Findings F2; echoed in `beta-closeout.md` §Process Observations item 2, `alpha-closeout.md` §Non-blocking findings F2) | **Project MCI** — real, recurring process gap in `cn-cdd-verify` itself, not in this cycle's own diff (`ledger.go` is untouched by cycle/574) | **Tracking issue filed**, since this will recur on every future triadic cycle in the α-push→β-push window, not just this one. Names both suggested fixes (recognize `gamma-scaffold.md` as an independent triadic signal; and/or teach `sectionPresent`/`validateSections` to tolerate the `§`-prefixed header convention). | [cnos#577](https://github.com/usurobor/cnos/issues/577) |

Both findings received an explicit disposition; neither was silently dropped (per §2.7's
"Silence is not triage" rule and CAP's four-way disposition set).

## Cycle closure declaration

All three CDD artifacts required before closure are present and internally consistent:
`gamma-scaffold.md`, `self-coherence.md` (AC1–AC7 all MET, independently re-derived by β rather
than accepted on α's prose), `beta-review.md` (verdict APPROVED, R0, both findings non-blocking
and now triaged above), `alpha-closeout.md`, `beta-closeout.md`, and `REPAIR-PLAN.md` (confirming
this cycle's one stranding was a substrate/workflow recovery, not an operator rejection, and that
no re-scaffolding or re-implementation was warranted or performed).

No AC is reopened by this close-out. No implementation was touched in writing it.

**cycle/574 is complete from γ's side**, pending δ opening the PR against `cycle/574` (with
`Closes #574` in the merge-commit subject, per `beta-review.md`'s merge instruction) and
requesting the `in-progress → status:review` transition via
`cn issues fsm evaluate --issue 574 --apply`, per `REPAIR-PLAN.md`'s "What actually remains" items
4–5. That step, the post-merge CI verification, the `.cdd/unreleased/574/` → `.cdd/releases/{X.Y.Z}/574/`
move, and the POST-RELEASE-ASSESSMENT.md are δ's/a subsequent γ-post-merge pass's actions, not
this artifact's — γ does not perform the PR-open or transition-request step itself.
