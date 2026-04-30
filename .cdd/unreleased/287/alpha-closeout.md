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
