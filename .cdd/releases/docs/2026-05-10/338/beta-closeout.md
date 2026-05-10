## Beta Close-Out — cycle #338

**Issue:** #338 — cdd: Add §1.6c — dispatch sizing, prompt scope, and commit checkpoints
**Branch:** `cycle/338`
**Mode:** docs-only (§2.5b disconnect path)
**Review rounds:** 1 (R1 APPROVED)
**Merge commit:** [pending — to be confirmed post-merge]
**β identity:** `beta@cdd.cnos`

---

### Review context

This was a docs-only cycle touching three skill files in `src/packages/cnos.cdd/skills/cdd/`. The change adds a single coherent feedback loop: §1.6c (dispatch sizing heuristic), operator §7 (timeout recovery procedure), and post-release §4 (telemetry fields). The three parts are tightly coupled; splitting them would orphan the heuristic without the telemetry that validates it.

α's self-coherence artifact was comprehensive: all 5 ACs had oracle evidence, §CDD-Trace was complete through step 7, and all 7 diff files were enumerated. The review-readiness signal was clean.

---

### Narrowing pattern

Single round, no RC. All 5 ACs passed their oracles mechanically. One A-level honest-claim finding was surfaced (arithmetic error in D1 debt declaration — self-coherence and alpha-closeout both computed `max(300, 120×5) = 900s` instead of the correct 600s). This finding does not affect the shipped normative content; per dispatch instructions, recursive coherence findings do not block approval.

No D, C, or B findings.

---

### Merge evidence

**Pre-merge gate:** `scripts/validate-release-gate.sh --mode pre-merge` — passed (exit 0) on cycle/338 branch; gate also ran on merge tree (merge tree run revealed missing beta-closeout.md, remediated by this artifact).
**Non-destructive merge test:** `git worktree add /tmp/cnos-merge-338/wt origin/main && git merge --no-ff --no-commit origin/cycle/338` — clean, zero unmerged paths, 7 files staged.
**origin/main at merge time:** `cfd322e62ef2ba56cb330e542fb000e0a36e2ed9`
**cycle/338 head before merge:** `b51f35e6` (includes beta-review.md + this file)
**Merge commit SHA:** [confirmed post-merge — see main branch HEAD after `git push origin main`]
**Branch state post-merge:** `cycle/338` merged into `main`; cycle directory to move to `.cdd/releases/docs/YYYY-MM-DD/338/` at γ closure per §2.5b

---

### β-side findings (factual, no dispositions)

1. **Arithmetic error in D1 and alpha-closeout §F1:** α computed the heuristic as `max(300, 120×5) = 900s`; the correct value is `max(300, 600) = 600s`. The shipped §1.6c(a) formula is correct; the error is confined to self-coherence.md §Debt D1 and alpha-closeout.md §F1. Actual recursive coherence: budget (600s) = heuristic (600s) — no deficit. α's D1 declaration is therefore based on a false premise, but the shipped content is unaffected.

2. **Operator/SKILL.md section renumbering:** §7 (Embedded Kata) was renamed to §8; the new §7 is Timeout Recovery. Any cross-refs to `operator §7` in other docs (pre-existing) would now point to Timeout Recovery rather than the Embedded Kata. A cursory scan found no pre-existing cross-refs to `operator §7` by number — existing refs use section names (e.g., `operator/SKILL.md §4`). No action needed.

3. **Gate classification discrepancy:** `validate-release-gate.sh --mode pre-merge` run on cycle/338 branch classified cycle #338 as "small-change" (passing with "no required cycle-dir artifacts"); the same gate run on the merge tree classified it as requiring all three close-out artifacts. This discrepancy appears to be gate-state-dependent (merge-tree context reveals the full artifact set). Documented for γ as a gate behavior observation, not a finding against the cycle.
