# α Close-out — #287

**Cycle:** #287 — γ creates the cycle branch — α and β only check out `cycle/{N}`
**Branch:** `cycle/287` (γ pre-provisioned per AC 12 self-application)
**Author:** α (`alpha@cdd.cnos`)
**Issue:** [#287](https://github.com/usurobor/cnos/issues/287)
**Merge commit:** `a5d0f21` ("Closes #287: γ creates the cycle branch — α and β only check out cycle/{N}")
**Round count:** 3 (R1 RC → R2 RC partial-withdrawal → R3 APPROVED)

This file is α's post-merge close-out per `alpha/SKILL.md` §2.8 + CDD.md §1.4 α step 10. **Voice: factual observations and patterns only.** No dispositions — γ owns triage. Per CDD.md §1.4 large-file authoring rule + the new α §2.5 incremental-write discipline (committed in `70ff2b1`, canonical on main), this file is written section by section to disk with progress reported after each section.

---

## Cycle summary

#287 implemented γ-creates-the-cycle-branch + canonical `cycle/{N}` naming + α/β-never-create-branches rule across 5 spec files (`CDD.md` + `alpha/`/`beta/`/`gamma/`/`operator/` SKILL.md). The diff was markdown-only (no executable surface, no schema-bearing contract). The cycle's self-application (AC 12) was the integration test: γ created `origin/cycle/287` from `origin/main` before α was dispatched; α `git switch`ed onto `cycle/287` from the harness pre-provisioned `claude/alpha-cdd-skill-1aZnw`; α/β/γ all committed to the same single named branch; β's verdicts landed on `cycle/287`, not on a separate harness branch.

Cycle structure:

| Round | Author | Commit(s) | Verdict / content |
|-------|--------|-----------|-------------------|
| α R1 | α | `9503aee` / `89b8575` / `11d5879` (post identity-correction) | spec change + self-coherence §§Gap–Debt + §Pre-review gate + §Review-readiness signal |
| β R1 | β | `8d2adb4` / `f89bf9f` / `74c3a6d` (post-rebase) | RC, 4 findings (F1 D, F2 C, F3 C, F4 A) |
| γ R1.5 | γ | `c91cf87` (post-rebase) | γ-clarification: F1+F2 collapse on fresh `origin/main` fetch (mechanical environment fact) |
| β R2 | β | `d9f1596` (post-rebase) | RC, F1+F2 **withdrawn**, F3 stands, F4 polish |
| α R2 | α | `de32200` (F4 polish) + `32a384f` (R1+R2 fix-round) | F1+F2 ACK-withdrawn, F3 fixed via path (a) retroactive re-author + force-push, F4 polish landed |
| β R3 | β | `f64bf44` | **APPROVED** |
| merge | β | `a5d0f21` (this) | `git merge --no-ff` of cycle/287 into main; auto-closes #287 |

The merge preserved cycle-branch ancestry via `--no-ff`; all 10 cycle commits are reachable from `main` through the merge commit. The cycle dir (`.cdd/unreleased/287/`) carries `self-coherence.md` (α R1 + R2 fix-round), `beta-review.md` (β R1 P1+P2+P3, β R2, β R3), and `gamma-clarification.md` (γ's mid-cycle mechanical fact). This α close-out lands on main directly per `alpha/SKILL.md` §2.8.

12/12 ACs met (verified by β at R3 head SHA `32a384f`). Diff content: 5 spec files (`CDD.md` +149 lines, `alpha/SKILL.md` +15 lines, `beta/SKILL.md` +17 lines, `gamma/SKILL.md` +87 lines, `operator/SKILL.md` +7 lines) + 3 cycle-dir artifacts.

---

## Findings (α-side)

### O1 — α git identity drift

**What happened:** α set `git config user.email "alpha@cnos.local"` at the very start of this session — a manual error inconsistent with `CDD.md` §1.4 α algorithm step 2's mandate (`alpha@cdd.{project}` = `alpha@cdd.cnos`). The subsequent three α R1 commits inherited the wrong email. self-coherence.md L7 metadata declared the canonical form (`alpha@cdd.cnos`) as if it had been honored — aspirational truth, not observed truth.

**Detection:** β R1 §2.0.0 Contract Integrity row "Status truth preserved" caught the metadata-vs-commit-author divergence (β R1 finding F3, severity C). β R3 verified the fix.

**Repair:** Path (a) — retroactive re-author. cycle/287 reset to origin/main; original SHAs cherry-picked back with `git commit --amend --reset-author --no-edit` applied to the 3 α commits only. β + γ commits cherry-picked without amend (preserving canonical authorship). Force-pushed with `--force-with-lease`. New α SHAs: `9503aee` / `89b8575` / `11d5879`. β verified canonical identity at R3 (all 5 α commits sign `alpha@cdd.cnos`).

**Surfaces affected:** `CDD.md` §1.4 α algorithm step 2 (mandate text). `alpha/SKILL.md` §2.6 (pre-review gate — does not currently include an identity-check row).

**Pattern:** spec mandate ↔ authoring-time gate gap. The role-identity rule lives in `CDD.md` §1.4 α step 2 as a one-line directive ("Configure git identity using the project name from the dispatch prompt: `git config user.name "alpha"` and `git config user.email "alpha@cdd.{project}"`"). It is not echoed in `alpha/SKILL.md` §2.6 pre-review gate as a verifiable row, nor as a §2.1 dispatch-intake gate. An α whose muscle memory has the wrong email form, or who skips re-reading §1.4 step 2, configures the wrong identity at session start and commits 3+ commits before any check fires. The gate that caught it (β contract-integrity check) fires post-α, not at α authoring time.

### O2 — β-side stale `origin/main` blind spot (β R1 F1+F2 false positives)

**What happened:** β R1 raised F1 (D-blocker, scope drift) + F2 (C, status truth) predicated on β's local `origin/main` being at `a8e67b7`. At α dispatch time `origin/main` was at `a8e67b7`; between α dispatch and β R1, σ pushed `70ff2b1` ("incremental write discipline for α and β artifacts") to `origin/main`. α observed the new `origin/main` and rebased cycle/287 onto it (α §2.6 row 1 transient-row rule fired correctly). β observed the cycle branch advancing but did not observe `origin/main` advancing because β's polling Monitor (`CDD.md` §1.4 β step 3 reference snippet) targets `origin/cycle/{N}` only — it does not poll `origin/main` for transitions. β computed the R1 review-diff base (`git diff main..cycle/{N}`) against the stale local `origin/main` ref and incorrectly attributed σ's commit to cycle/287's authorized scope.

**Detection:** γ's mid-cycle clarification at `c91cf87` captured the mechanical environment fact synchronously (`git fetch --verbose origin main` + `git rev-parse origin/main` + `git branch -r --contains 70ff2b1`) and showed that current `origin/main` was at `70ff2b1`, not `a8e67b7`. β re-fetched, verified γ's fact, and withdrew F1+F2 in R2.

**Repair:** No spec-side repair this cycle. β R3 verdict + α R1+R2 fix-round both name the gap as the proximate cause of the §9.1 review-rounds-> 2 trigger fire (the cycle would have closed at R2 without F1+F2). β R3 records the gap in `beta-closeout.md` for γ-PRA triage.

**Surfaces affected:** `beta/SKILL.md` (no synchronous-review-base-fetch rule). `review/SKILL.md` (Output Format does not require base-SHA freshness verification before computing diff). `CDD.md` §Tracking (the `git fetch --quiet` reliability rule + N=10 re-probe rule added by AC 8 of this cycle is for cycle-branch polling, not for review-base discipline).

**Pattern:** transition-only-polling-on-X creates blind spot for Y. β's polling correctly watches `origin/cycle/{N}` per the new spec (transitions on cycle-branch head SHA, fix-round commits, etc.) but the diff-base computation reads `origin/main` synchronously without re-fetching. The polling discipline owns "future state of cycle/{N}"; the synchronous discipline must own "current state of main" — they're separate channels. The §Tracking synchronous-baseline-pull rule (3.61.0 PRA F1 cycle iteration / α §2.5 §Skills) names the principle for cycle-dir state but does not generalize it to review-base state.

### O3 — `eb48e17` referenced in issue body but unavailable in this repo

**What happened:** Issue #287 ## Work-shape says: "Spec text is already drafted in #283's chat record + the rolled-back local commit `eb48e17`'s message; α can ingest both as the authoritative spec source and refine." The commit `eb48e17` is local-only to γ's prior session (rolled back per 3.61.0 PRA §3 γ self-observation: "γ-creates-branch scope expansion was correctly rolled back"). The commit object is not present in this repo's git store (`git show eb48e17` returns "ambiguous argument: unknown revision"). α reconstructed spec text from the issue body's AC text + 3.61.0 PRA §4b/§7 next-MCA notes + #283 cycle artifacts (`gamma-clarification.md` from #283).

**Detection:** α surfaced as §Known debt #1 in self-coherence.md R1 (pre-review). β did not raise as a finding — α's reconstruction matched the AC text closely enough that no spec-text drift was detected at R3.

**Repair:** None needed for this cycle. The reconstructed spec text passed β R3.

**Surfaces affected:** `issue/SKILL.md` (no rule against referencing rolled-back commits in issue ## Work-shape). `gamma/SKILL.md` §2.5 (no rule for γ to inline rolled-back commit content into the issue body when the rolled-back work is the spec source).

**Pattern:** rolled-back-but-referenced. When γ rolls back a local commit that was the working spec text and files an issue from that work, the issue body references the unavailable commit. α downstream cannot read what was intended; α must reconstruct from secondary sources (AC text + PRA + cycle artifacts). If the AC text is sufficiently prescriptive, reconstruction succeeds; if it is descriptive, reconstruction can drift.

### O4 — `Monitor` polling false-positive on short-SHA prev-state

**What happened:** α set up the post-review-readiness polling loop with `prev_head="32a384f"` (short SHA from prior `git push` output). The loop's first iteration ran `cur_head="$(git rev-parse origin/cycle/287)"` which returned the full 40-char SHA `32a384f662ea7fdb6c8b0d525e7a2317d0ec9879`. The string comparison `[ "$cur_head" != "$prev_head" ]` evaluated `32a384f662... != 32a384f` → true. The loop fired immediately with no actual transition, returning a stale wake.

**Detection:** α observed the false-positive wake and re-armed the loop with `prev_head="$(git rev-parse origin/cycle/287)"` (full SHA).

**Repair:** Local fix to α's polling-loop boilerplate. No spec change needed — the `CDD.md` §Tracking shell snippet examples already use `git rev-parse` consistently for both `prev` and `cur`; the bug was α-side hand-typed loop only.

**Surfaces affected:** None spec-side. Operator-experience pattern.

**Pattern:** SHA-form mismatch in transition-only polling. When the wake-up condition is "SHA changed" and one of the two SHAs is hand-pasted from a `git push` output (which prints short SHAs by default), the comparison is false-positive. The §Tracking reference snippets are correct (`git rev-parse` returns full SHA on both sides); the gap is in operator/α practice when manually arming polling loops outside of the canonical reference snippet.

### O5 — Cycle-branch existed before α dispatch (AC 12 self-application worked end-to-end)

**What happened:** γ created `origin/cycle/287` from `origin/main` (the new γ Phase 1 step 3a) before α was dispatched. α intake saw the harness pre-provisioned branch (`claude/alpha-cdd-skill-1aZnw`); the dispatch prompt explicitly named `Branch: cycle/287`; α `git switch cycle/287` from the harness branch as the first action and never created any branch. β intake similarly switched from a harness branch (`claude/cnos-skill-module-x9jTE` per β R1 P2 AC 12 evidence) onto `cycle/287` and committed all verdicts there. γ's clarification + β's R1 + R2 + R3 verdicts all landed on `cycle/287` directly.

**Detection:** β verified at every round (R1 P2 AC 12 row, R2 §F1 withdrawal evidence, R3 verification). β R3 §Search-space closure: "The cycle's self-application (AC 12) is exemplified by the cycle's own execution".

**Pattern:** spec-implementing cycle exemplifies its own rule on first execution. The fix worked as designed: γ owns branch creation; α/β check-out only; one cycle = one branch = one named target. The contrast against #283 R1 F1 (β verdict landed on a separate harness branch and had to be cherry-picked) is direct — that friction is structurally absent in this cycle.

### O6 — No new findings (positive observation)

**What happened:** β R3 found no new findings beyond F1+F2 (withdrawn), F3 (fixed), F4 (fixed). The §AC walk re-verified 12/12 ACs met. Diff content invariance verified. No regression introduced by the α R2 rebase.

**Pattern:** Per α §2.5 §Self-check rule ("did α's work push ambiguity onto β?"), the answer for this cycle is **no for the diff itself, partially for the metadata**. The diff content was clean from R1 onwards (β R1 found 0 actual diff drift; R1 F1+F2 were β-side stale-fetch artifacts and R1 F3+F4 were metadata-vs-spec mismatches). α's authoring discipline (peer enumeration, intra-doc grep, cross-reference check) caught all in-diff issues before review. The metadata gate (git identity) was the only authoring-time leak.

---

## Friction log

What went wrong in the cycle itself (process, not code), in execution order:

1. **Session-start identity miss-config (α-side, ~00:13 UTC).** α set `git config user.email "alpha@cnos.local"` from muscle memory instead of re-reading `CDD.md` §1.4 α step 2's mandate. The miss propagated through 3 α R1 commits before β caught it.

2. **Mid-cycle σ push to `origin/main` (00:12 UTC, before α dispatch but after the issue was filed).** σ landed `70ff2b1` ("incremental write discipline for α and β artifacts") on `origin/main` between issue filing and α's session start. α's pre-review gate row 1 transient-row rule fired correctly at 00:18 UTC and α rebased cycle/287 onto `origin/main`. The rebase caused one merge conflict in `alpha/SKILL.md` §2.7 (HEAD: σ's "append the review-readiness section as a separate commit" wording vs α's local: "commit `.cdd/unreleased/{N}/self-coherence.md` with explicit `origin/cycle/{N}` naming"). α resolved by keeping σ's incremental-commit wording + α's explicit branch naming. Cost: ~5 minutes of conflict resolution + transient-row re-validation in self-coherence.md §Pre-review gate row 1.

3. **β R1 false positives F1+F2 (00:25 UTC).** β reviewed against stale local `origin/main` (`a8e67b7` — pre-σ-push state). β's polling Monitor watches `origin/cycle/287` only, so β had no wake-up event for "main moved." β R1 verdict raised F1 (D-blocker, σ commit out-of-scope on cycle/287) and F2 (C, status truth — α misidentifies `70ff2b1` as `origin/main`). Both were β-side artifacts of the stale local ref. Cost: γ-clarification commit (`c91cf87`) capturing the synchronous mechanical fact; β R2 re-fetch + withdrawal; α R2 fix-round narrative addressing both refute + acknowledge.

4. **α R2 force-push history rewrite (00:35 UTC).** F3 (identity drift) was a contract-integrity finding β stood by in R2; β R2 §F3 §α response options enumerated path (a) "retroactive re-author + force-push" and path (b) "apply rule from next commit forward." α chose (a). Execution: cycle/287 reset to origin/main; original SHAs cherry-picked back with `--reset-author` applied to α commits only (`git cherry-pick --reset-author` is unavailable in git 2.43.0; α used cherry-pick + `git commit --amend --reset-author --no-edit` instead). β + γ commits cherry-picked without amend. Force-pushed with `--force-with-lease`. The lease was rejected once because α's local `origin/cycle/287` ref was stale (γ-clarification + β R2 had landed during α's session in parallel). α fetched, re-did the cherry-pick rebuild on top of the new origin/cycle/287 (which now carried 8 commits: 3 α R1 + 3 β R1 + γ-clarification + β R2), pushed successfully. Cost: full chain rebuild (10 cherry-picks) + retry of force-push.

5. **`alpha/SKILL.md` rebase conflict in §2.7 — σ's incremental-write rule overlapped α's review-request push instruction.** Single-line conflict; HEAD's "append the review-readiness section as a separate commit" + α's `origin/cycle/{N}` explicit naming combined cleanly into one sentence. Cost: 1 manual edit during rebase.

6. **Polling-loop short-SHA false positive (~00:42 UTC).** α armed the post-fix-round polling loop with `prev_head="32a384f"` (short SHA from `git push` output). The loop fired immediately because `git rev-parse origin/cycle/287` returned the full 40-char SHA. α re-armed with full SHA. No β impact (the false-positive wake just re-checked state without action). Cost: ~30 seconds of false-positive review.

7. **Round count = 3 (§9.1 trigger fires on strict reading).** R1 = α-implementation review with 4 findings; R2 = γ-clarification + β R2 (F1+F2 withdrawn, F3+F4 stand); R3 = α fix-round + β R3 approval. The R3 round count is *attributable to β-side process* (stale `origin/main`) per β R3 + α R1+R2 fix-round narratives — without F1+F2's false positives, the cycle would have closed at R2 (F3+F4 fix-round → approval).

**Total wall-clock cycle time:** ~33 minutes (α dispatch ~00:14 UTC → β R3 approval at 00:45 UTC → merge at 00:47 UTC). Of which: α R1 authoring ~5 minutes; β R1 review ~10 minutes; γ-clarification + β R2 ~7 minutes (parallel with α R2 prep); α R2 fix-round + force-push ~5 minutes; β R3 ~3 minutes; merge ~1 minute. The frictions above account for the bulk of "process overhead beyond minimum"; actual content-authoring time was small.

---

## Observations + patterns

### P1 — First-cycle-of-new-protocol friction is *partially* irreducible

The diff this cycle landed *is* the protocol the cycle operates under. AC 12 (self-application) is therefore not a side-property — it is the integration test. This cycle exemplified its own rule end-to-end: γ created `cycle/287` before dispatch; α/β only checked out; β's verdicts landed on the cycle branch.

The friction that remains in a self-applying cycle has two sources: **(a) friction the new protocol does not address** (e.g. β-side stale-`origin/main`, observation O2 in §Findings — the new spec adds AC 8 `git fetch` reliability for cycle-branch polling, not for review-base polling) and **(b) friction the operator's environment imposes on top of the protocol** (e.g. harness pre-provisioned per-role branches, observation O5 — the protocol explicitly forbids those, and the cycle's structural correction is the spec wording itself + α/β refusal to use them).

The (a)-class friction (O2) is the cycle's #1 candidate observation for γ-PRA cycle iteration. The (b)-class friction was anticipated in the issue body ## Active design constraints and structurally addressed.

Compared to #283 (the previous self-applying cycle), this cycle exhibited *less* (b)-class friction (β refused harness instructions cleanly via `beta/SKILL.md` §1, switched to `cycle/287`, and committed verdicts there — there was no need for cherry-picking from a separate β branch as in #283 R1 F1). #283's primary friction was structurally absent here.

### P2 — γ-clarification as a third-role mid-cycle artifact

`gamma-clarification.md` filed at `c91cf87` was a **purely mechanical environment-fact transfer** (per `gamma/SKILL.md` §2.5: "γ may ... provide mechanical environment help"). It did not adjudicate β's findings in either direction. It captured a synchronously-observed `git fetch --verbose` + `git rev-parse origin/main` + `git branch -r --contains 70ff2b1` and let β re-evaluate F1+F2 against the fact.

Pattern: γ-clarification works as designed when it transfers an artifact fact, not a reasoning state. β R2 withdrew F1+F2 by independent verification of γ's fact, not by accepting γ's interpretation. This preserves epistemic separation (`gamma/SKILL.md` §2.5: "γ does not forward β's internal reasoning transcript to α; γ does not forward α's hidden rationale transcript to β; γ may not author the implementation fix inside the review loop"). The cycle exemplified the role-boundary discipline introduced in #283.

### P3 — Authoring-time gate vs review-time gate (O1 generalization)

The role-identity rule lives in `CDD.md` §1.4 α step 2 as a one-line directive. It is not echoed in `alpha/SKILL.md` §2.6 pre-review gate as a verifiable row. The gate that caught the violation (β contract-integrity check) fires *post-α*, not at α authoring time. This is a general pattern: spec mandates that are not echoed as α authoring-time verifiable steps survive past α authoring and are only caught at β review.

Variants of this pattern in the cycle's diff:
- α §2.6 pre-review gate has 10 rows; none verify the canonical git identity at HEAD (rows 1 + 10 are transient external state; rows 2–9 are diff/artifact state; identity is committer state, a third axis).
- α §2.1 dispatch intake has identity config in step 1 ("configure α git identity") but does not verify the configured email matches the canonical form.
- The new α §2.5 incremental-write discipline (committed in `70ff2b1`) names "report progress after each commit" but does not name "verify identity on each commit" — the same authoring-time-gate-vs-review-time-gate gap.

This is a candidate generalization for γ-PRA: not a single skill patch, but a pattern recognition. Multiple authoring-time identity / metadata gates would close the gap categorically.

### P4 — Polling discipline asymmetry (O2 generalization)

The polling discipline added in #283 (synchronous baseline + transition-only stdout) and refined in this cycle (single named branch + N=10 reachability re-probe) is now strong for **cycle-branch state**. The diff-base computation, however, reads `origin/main` synchronously without re-fetching — and `origin/main` is *not* the polled surface for any role.

This creates an asymmetry: cycle-branch state is event-driven (transitions wake the role); main state is polled-on-demand (only when the role explicitly fetches). β computes the review diff base from main; β does not re-fetch main before computing.

The same asymmetry would affect γ if γ's PRA computes against main without re-fetching, or α if α reads `origin/main` for any decision other than the rebase trigger (which already includes a fetch). The pattern is general: **synchronous reads of a non-polled remote surface need an explicit pre-read fetch.**

This is the friction class observation O2 names. β R3 records it for `beta-closeout.md`; γ-PRA can characterize it as a §9.1 cycle-iteration trigger (avoidable tooling failure: stale local ref produced false positives) and decide whether it warrants a `beta/SKILL.md` patch in this same release boundary or a follow-on issue.

### P5 — Path (a) vs path (b) for identity correction (O1 specific)

β R2 §F3 explicitly enumerated two response paths:
- (a) retroactive re-author + force-push (β preferred for clean cross-cycle archeology)
- (b) apply rule from next commit forward + disclose as known debt

α chose (a). The execution worked cleanly: cycle/287 reset to origin/main; original 8 commits cherry-picked back with `--reset-author` selectively applied (α commits only); force-pushed with `--force-with-lease` (with one rejection due to parallel γ + β commits, recovered by re-fetching and re-doing the chain rebuild on top).

Pattern: when retroactive history rewrites are needed and the cycle has multi-author commits, the operation requires:
1. Reset target branch to a clean base
2. Cherry-pick all original commits in order, with `--reset-author` selectively applied per commit's authoring role
3. Force-push with `--force-with-lease`
4. Other roles must `git fetch --force` or `git remote update --prune` to refresh local refs

The total cost was tractable (~5 minutes including one rejected push retry). β R3's content references the new SHAs directly (β's R1+R2 verdict text references the pre-rewrite α SHAs as historical state, which is correct narrative — those SHAs *did* exist; they just no longer have a branch ref pointing to them).

Pattern variant for future cycles: if the rewrite were to span more commits or be done late in the cycle, the cost grows. Doing it at R2-fix-round time (before the cycle is merged) was the correct timing — post-merge rewrite would be a much larger problem.

### P6 — Diff-content invariance under history rewrite

β R3 §4 verified: "Diff content invariance verified. β re-ran `git diff --stat origin/main..origin/cycle/287` at R3 head: 5 spec files + 3 cycle-dir artifacts. `CDD.md` is now +149 lines (was +147 in R1; +2 lines is exactly the F4 parenthetical wrapping, consistent with α's claim)."

This worked because cherry-pick preserves diff content; only authorship + commit SHA change. β's content-level review against R1 + R2 + R3 was strictly additive (R1 reviewed the spec change; R2 walked AC + Issue Contract; R3 verified the F3+F4 fixes). The history rewrite introduced no new content for β to review beyond the F4 polish + the fix-round narrative.

Pattern: if the only change between rounds is identity correction + small polish, β's R3 verification effort is small (β R3 was ~3 minutes). The expensive part of any review is content cognition, which the rewrite did not invalidate.

### P7 — `eb48e17` reconstruction (O3 specific)

The issue body's reference to `eb48e17` (the rolled-back commit message that was the canonical spec source) was unavailable. α reconstructed from:
- Issue body's 12 ACs (prescriptive — verbatim shell snippets, file paths, exact wording for §4.2 / §4.3 / dispatch prompt format)
- 3.61.0 PRA §4b (cycle iteration) and §7 (next-MCA) — descriptive context
- #283 cycle artifacts at `.cdd/releases/3.61.0/283/` — `gamma-clarification.md` provided the "branch-polling canonical, one cycle branch holds all role artifacts" rule that #287 builds on

The reconstruction was sufficient: β R3 found 12/12 ACs met without raising spec-text drift. The risk surface was the descriptive sections (Role table responsibilities text, §Tracking polling-table prose, transition-loop snippet wording), which α produced from existing CDD.md style + the AC's prescriptive constraints.

Pattern: when γ rolls back local spec work and references the rolled-back commit in the issue ## Work-shape, α must reconstruct from the AC text alone (or from secondary sources like prior cycle artifacts). The reconstruction succeeds when ACs are prescriptive; it carries drift risk when ACs are descriptive. #287's ACs were sufficiently prescriptive that the risk did not materialize, but the pattern is a candidate skill-text observation for `issue/SKILL.md` (a rule like "when filing an issue from rolled-back local work, inline the rolled-back content into the issue body or attach as a comment, since downstream roles cannot read the rolled-back commit").

---

## Cycle-level engineering reading

Per `CDD.md` §9.1 cycle level assessment, applying the level framework from `docs/gamma/ENGINEERING-LEVELS.md` §6 to this cycle (change set):

### L5 — local correctness

Was the diff locally correct before review? Did it compile, pass tests, follow current patterns? **Yes, modulo the metadata gate.**

- The diff is markdown-only — there is no compile / test surface.
- The 5 spec files passed β's contract-integrity check (canonical sources/paths, constraint strata, exception-reasoning, path-resolution-base, proof-shape, cross-surface projections, no-witness-theater, intra-doc repetition — all yes at R1 P1).
- The diff was internally coherent at R1: cross-references resolved, intra-doc repetition was consistent across all 4 grepped facts (γ creates / α/β never create / `origin/cycle/{N}` / legacy `claude/*` warn-only).
- **The metadata gate (git identity) failed L5 strictly** — α's commits did not honor `CDD.md` §1.4 α step 2's mandate. β R1 F3 (severity C) caught it.
- Per `CDD.md` §9.1 L5 trigger ("if mechanical errors reached review, L5 was not met"), F3 is *contract*, not strictly *mechanical*. β classified F3 as "contract (identity-truth) / mechanical." If β's classification weights mechanical, L5 trigger fires; if it weights contract, L5 may still hold.

**L5 verdict (α self-assessment):** held (with caveat). The diff itself was locally correct; the identity gate is a meta-property the diff doesn't directly encode. Defer to γ-PRA for the canonical reading.

### L6 — system-safe execution

Did the change stay coherent across docs, runtime, artifacts, tests, and operator truth? Were failure modes bounded and visible? **Yes.**

- Cross-surface coherence: the 5 spec files agree (peer enumeration §Role-skill peers + §Lifecycle-skill peers; intra-doc grep across CDD.md). β verified consistency at R1 P2 (10 cross-surface projections updated).
- Failure modes: γ pre-flight names rejection conditions (branch already exists, stalled cycle dir, scope undeclared, base SHA unknown, issue closed). α/β refusal mechanism is explicit ("α never creates a branch" / "β never creates a branch"). The new contract is enforceable by inspection.
- Operator-truth: AC 12 self-application proven by execution (γ created `cycle/287`; α/β switched onto it; β verdicts on cycle branch; canonical git identities post-rewrite). The cycle's git history *is* the operator-truth record.
- The σ-push mid-cycle (`70ff2b1`) created a rebase event; α §2.6 row 1 transient-row rule fired correctly and the rebase + conflict resolution were clean.

**L6 verdict:** held. No cross-surface drift reached review. β R1 F1+F2 were *β-side-process* drift (stale `origin/main`), not *cycle-diff* drift. The cycle exhibited zero false claims about cross-surface state at R1.

### L7 — system-shaping leverage

Did the cycle change the system boundary so the friction class gets easier or disappears? **Yes.**

The cycle's diff *eliminates* the branch-discovery friction class for all future cycles:
- pre-#287: α created the branch under harness-encoded `claude/{slug}-{rand}`; γ + β had to glob-discover; β's verdict often landed on a separate harness branch (#283 R1 F1).
- post-#287: γ creates `cycle/{N}` under the canonical pattern; α/β `git switch` from the dispatch prompt; one cycle = one branch = one named target.

The system boundary changed. The friction class disappears, not just locally fixes. Future cycles do not re-encounter the discovery problem because the rule is structural, not heuristic.

Concretely: #287's ACs encode three structural mechanisms that close the friction class:
1. **γ pre-flight + creation** (§1.4 γ Phase 1 step 3a + §4.3) — γ owns the canonical branch name; one source of truth.
2. **Refusal mechanisms** in α/β (§1.4 α step 5a + §1.4 β step 3 + alpha/SKILL.md §2.1 step 2 + beta/SKILL.md §1) — α/β never invent or accept other branches.
3. **Single named polling target** (§Tracking polling table) — no glob discovery for new cycles.

The legacy glob `'origin/claude/*'` is named warn-only / retrospective — backwards compatible without reintroducing the friction.

**L7 verdict:** achieved (diff-level). The cycle ships an MCA that eliminates a friction class for future cycles.

### Cycle level — lowest miss

Per `CDD.md` §9.1 ("the cycle level is the lowest miss"):

- L5: held (with metadata caveat — F3 caught at R1, fixed cleanly at R2; if γ-PRA classifies F3 as mechanical, L5 trigger fires)
- L6: held (no cross-surface drift; β R1 F1+F2 were β-side, not cycle-diff)
- L7 diff: achieved (friction class eliminated)
- L7 cycle execution: caps at L6 if §9.1 review-rounds-> 2 trigger fires. β R3 confirms it does; β attributes the trigger to β-side stale-`origin/main`, not to α's diff.

**α self-assessment of cycle level:** **L6 cycle / L7 diff** if γ-PRA accepts the F3-is-contract-not-mechanical reading + the round-count attribution to β-side process. **L5 cycle / L7 diff** if F3 is reclassified as mechanical (the L5 trigger fires on "mechanical error reached review"). Both readings name the diff as L7 — that part is unambiguous.

α defers the canonical reading to γ-PRA. The friction this cycle's diff eliminates is real; the cycle's execution had two L6 events (F3 metadata, β-side stale-fetch round inflation) — the question is whether either drops to L5.

### What the cycle taught

The cycle is its own integration test for the new protocol — and it passed. Three structural successes in the same cycle:

1. **γ-creates-the-branch worked end-to-end.** No branch-name invention, no glob-match failure, no β-side cherry-pick. The single named branch held all role artifacts (α R1, β R1+R2+R3, γ-clarification, fix-round). The contrast against #283 R1 F1 is structural.

2. **`git fetch` reliability rule (AC 8) is correct for cycle-branch polling but does not generalize to review-base polling.** The cycle surfaced this gap by exercise: β R1 F1+F2 are exactly the failure mode AC 8's N=10 re-probe rule prevents — but only for cycle-branch state, not for `origin/main`. Whether γ-PRA folds the gap into a near-term `beta/SKILL.md` patch or files as a follow-up is γ's call; α flags as the cycle's primary observation for next-MCA candidacy.

3. **γ-clarification as mechanical-fact transfer** preserved epistemic separation. β re-evaluated F1+F2 against γ's synchronously-observed fact, not against γ's interpretation. The role-boundary discipline #283 introduced ("γ transfers artifact facts only") is robust under stress (R1 raised D-blockers; γ resolved by mechanical fact, not by adjudicating).

### What the cycle did *not* test

- **Multi-issue release boundary.** This cycle ships in isolation (one issue per release). Future cycles that bundle multiple issues need to test whether the cycle/{N} canonical name handles multi-issue scoping cleanly. (Open question for #286 encapsulation work.)
- **Operator-spawned dispatch with full prompt format.** The dispatch prompt α received was operator-pasted; the new `Branch:` line was honored. Whether `cn dispatch` (the future CLI per #286 dependencies) generates the prompt correctly is for that cycle.
- **Multi-cycle simultaneous activity.** Only `cycle/287` was active during this session; whether two `cycle/{N}` branches with different `{N}` would cause polling conflicts is untested. The single-named-branch rule structurally precludes the conflict at the cycle-scope level, but γ + δ multi-cycle tracking (e.g. `git branch -r --list 'origin/cycle/*'` glob in `operator/SKILL.md`) is a different pattern.

— α (`alpha@cdd.cnos`) at 2026-04-30 ~00:55 UTC. Close-out complete. Cycle handoff to γ for PRA + close-out triage.
