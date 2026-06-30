# Conflict: cnos#524 W4 RCA landed an overlapping guard on `main` while `cycle/532` was open

**Status: STOP — surfaced to γ/δ/operator. Not resolved by α. Do not improvise past this point without explicit direction.**

## What happened

While implementing cnos#532 on `cycle/532` (branched from `origin/main` at `8f4f01b0`), `origin/main` advanced to `7950ab3d` via PR #531: **"cnos#524 W4 RCA (A2, manual_delta_repair): dispatch closeout-integrity guard"** (commit `368a1832`, merged 2026-06-30T19:1Xish UTC — i.e. landed *after* cycle/532 was scaffolded and *during* this α dispatch).

This is the exact same root-cause incident (#524 W4 empty-review) that issue #532 itself names as its motivating evidence ("Empty-review failure | Observed in #524 W4 attempt | Real" — #532's own Status-truth table). A separate, independent remediation (`run_class: manual_delta_repair`, not a dispatched cell — a direct operator/manual-δ repair) landed a working guard for the same failure mode, on `main`, concurrently with this cycle.

## What #524 W4's RCA shipped (already on `main`, not something I wrote)

- `dispatch-protocol/SKILL.md` §2.9 ("Closeout integrity") + §4.9 (verify check) + failure mode **D12** ("Empty review").
- `cds-dispatch/SKILL.md` + `prompt.md`: a new "Closeout integrity preflight" section, inserted at the same anchor point (immediately before `## Lifecycle transitions`) that I also touched — but as an **additive new section**, not an edit to the same prose I edited (step 7 + the lifecycle-transitions table row). My edits and #524 W4's section insertion are at different line ranges in the same file region.
- A `deliverable_evidence:` YAML block (embedded in the cycle's **closeout** document, e.g. `gamma-closeout.md`) naming: `pr`, `head_sha`, `base_sha`, `commits_beyond_base` (count, must be > 0), `closeout_artifacts` (list).
- `scripts/ci/check-dispatch-closeout-integrity.sh` — a presence-of-contract guard (same `need()` pattern as #516) **plus** an embedded `closeout_violation(status, has_pr, has_commits)` boolean self-test detector (not a structural YAML-field validator like my Guard B — it takes three flags, not a file).
- `.github/workflows/build.yml` job `dispatch-closeout-integrity`, appended at the exact same insertion point (end of file, after `dispatch-repair-preflight`) where I appended my `review-request-preflight` job — **mechanical git conflict on rebase**, trivial to resolve (keep both jobs), confirmed via a real `git rebase origin/main` attempt (see below).
- `docs/evidence/rca/2026-06-30-cnos524-w4-empty-review.md` — the RCA + an AC-style remediation receipt.

## Why this is a genuine conflict, not just a git merge nuisance

The `build.yml` job-append conflict is mechanically trivial (I confirmed via a real rebase attempt — see "What I verified" below — that it is a clean two-way add, resolvable by keeping both jobs). The **doctrinal** conflict is the real issue:

1. **Two named proof artifacts for the same precondition.** #524 W4 created `deliverable_evidence:` (a block embedded in a closeout doc). #532 (this cycle, per the issue body's own explicit "Proposed artifact" section and the operator's clarification comment) creates `REVIEW-REQUEST.yml` (a standalone file at `.cdd/unreleased/{N}/REVIEW-REQUEST.yml`). Both claim to prove the same five-ish facts (PR exists, commits beyond base, branch exists/differs, required closeout artifacts exist, evidence named) before the same `status:in-progress → status:review` transition. A wake author reading both sections of `cds-dispatch/prompt.md` after a naive merge would be told to produce **two different proof artifacts** for what is conceptually one gate.

2. **Two named guard scripts with overlapping but non-identical checks.** `check-dispatch-closeout-integrity.sh` (#524 W4) is a boolean detector (`closeout_violation(status, has_pr, has_commits)`) plus presence-of-contract grep. `check-review-request-preflight.sh` (#532, this cycle) is a richer structural YAML validator (Guard B) plus presence-of-contract grep (Guard A). They are not aliases of each other — #524 W4's guard does not check `matter.branch` regex shape, `no_op`/`no_op_approval`, or `artifacts.*` on-disk existence; #532's guard does not independently check `commits_beyond_base` as a literal count field. Running both is not wrong, exactly, but it is redundant doctrine: two independently-evolving "is status:review proven" mechanisms now exist in the same dispatch substrate, which is itself a coherence smell (and arguably the kind of thing #532's own kernel-doctrine addition — "no matter, no review; no proof, no status:review" — exists to prevent: a single invariant should have one mechanically-checkable realization, not two that can drift apart).

3. **`dispatch-protocol/SKILL.md` would carry two numbered sections (§2.9 from #524 W4, plus my new unnumbered §"Review-request proof gate") asserting overlapping but differently-worded preconditions on the same `status:review` row**, and two failure-mode catalogue entries (D12 from #524 W4, no equivalent named entry from #532 — I did not add a `D`-numbered failure mode, since γ's scaffold and δ's dispatch did not instruct me to touch the `## 5. Failure modes catalogue` section, which is exactly where D12 lives and where #532's failure mode should arguably also be named, or rather, should **be the same named failure mode** as D12, since they describe the identical incident).

## What I verified (not improvised — mechanical facts only)

- Ran `git fetch origin main` — confirmed `origin/main` advanced from `8f4f01b0` (cycle/532's base) to `7950ab3d` (13 commits ahead), including `368a1832` (the #524 W4 RCA commit) merged via PR #531.
- Ran a real `git rebase origin/main` attempt on `cycle/532` (with all 12 of my commits applied). Result: rebase proceeded cleanly through commit 5 of 13 of my history, then hit a genuine `CONFLICT (content)` in `.github/workflows/build.yml` at the exact job-append point described above. `.github/workflows/cnos-cds-dispatch.yml` and the golden auto-merged without conflict markers (git's three-way merge resolved them, since #524 W4's and #532's prompt insertions are at different line ranges) — but I have **not** verified that the auto-merged *content* is doctrinally coherent (per the redundancy concern above), only that git's line-level merge did not raise a conflict marker. I aborted the rebase (`git rebase --abort`) rather than resolve the `build.yml` conflict and accept the auto-merged prompt content, since accepting it would silently ship the doctrinal redundancy as fact without anyone deciding it should ship that way.
- Confirmed `git status` is clean post-abort and all 12 of my commits are intact at their original SHAs (branch unchanged, still based at `8f4f01b0`).

## What I did NOT do

- I did not rebase `cycle/532` onto the new `origin/main`.
- I did not resolve the `build.yml` conflict.
- I did not decide whether `REVIEW-REQUEST.yml` should subsume/supersede `deliverable_evidence:`, whether the two guards should merge into one, whether #532's doctrine should reference D12 instead of (or in addition to) defining its own un-numbered review-request-proof-gate section, or whether #524 W4's `check-dispatch-closeout-integrity.sh` should be deprecated in favor of `check-review-request-preflight.sh` (or vice versa, or both kept as independently-useful layered checks — that is itself a legitimate design option I am not positioned to choose unilaterally).
- I did not improvise a "least-surprise" merge (e.g., silently keeping both doctrines side by side) because that would ship exactly the redundancy named in §2 above without anyone having decided it is acceptable.

## Candidate resolution paths (named for δ/operator judgment, not chosen by me)

- **(a) Supersede.** #532's `REVIEW-REQUEST.yml` + Guard A/Guard B fully subsumes #524 W4's `deliverable_evidence:` + `check-dispatch-closeout-integrity.sh` (the issue's own scope is broader and the operator's clarification comment for #532 is more recent and more detailed than the RCA's remediation). Deprecate/retire #524 W4's guard and §2.9 in favor of #532's surfaces; fold D12 into a cross-reference to #532's invariant; retarget `dispatch-protocol/SKILL.md`'s numbered section so there is exactly one canonical "deliverable proof before status:review" section, not two.
- **(b) Layer.** Keep both as independently-useful layered checks (defense-in-depth — #524 W4's lighter boolean check as a fast first-pass, #532's richer structural check as the authoritative gate), but explicitly cross-reference each other in doctrine so a reader does not think they are two unrelated gates, and reconcile the `deliverable_evidence:` vs `REVIEW-REQUEST.yml` field vocabularies (e.g. make `deliverable_evidence:` a derived/computed view of `REVIEW-REQUEST.yml`'s `matter:` block, or vice versa) so they cannot silently drift apart.
- **(c) Merge.** Combine the two YAML shapes into one (e.g. `REVIEW-REQUEST.yml` gains a `deliverable_evidence:`-shaped sub-block, or `deliverable_evidence:` is renamed/reshaped to match `REVIEW-REQUEST.yml`'s `matter:`/`artifacts:` naming), and combine the two guard scripts into one (my `check-review-request-preflight.sh` could absorb `check-dispatch-closeout-integrity.sh`'s self-test detector as an additional check, or vice versa).

## Current branch state

`cycle/532` is fully implemented against issue #532's 10 ACs as scoped (see `.cdd/unreleased/532/self-coherence.md`), all local CI gates verified green (see self-coherence §Self-check), but **NOT rebased onto current `origin/main`** because the rebase surfaces this unresolved doctrinal conflict. All 12 implementation commits are pushed-ready but not yet pushed to `origin/cycle/532` pending this conflict's resolution (see next steps).

## Recommended next step

γ/δ should read this document, decide among (a)/(b)/(c) (or a different resolution), and either: (i) re-dispatch α with the resolution as a pinned instruction, or (ii) resolve the rebase/merge themselves and hand the reconciled branch back to α for final self-coherence re-verification, or (iii) treat #532 as superseded by #524 W4's already-shipped remediation and close #532 with a closure note pointing at the RCA — but that is also not α's call to make.
