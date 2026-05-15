---
cycle: 365
type: alpha-closeout
date: "2026-05-15"
author: alpha@cdd.cnos
---

# α Close-out — #365

## Summary

Bumped `is_legacy_version` in `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` from `minor < 74` to `minor < 77`, rewrote the helper's comment block to name the new cutoff + covered artifacts + rationale, updated both downstream warn strings (`(pre-v3.65.0)` → `(pre-v3.77.0)`), added three boundary tests (29 → 39 assertions) plus two fixture writers, relocated the `incomplete-triadic` test fixture from `3.75.0/200/` to `3.78.0/200/` so the strict-path proof sits post-cutoff, and added a `CHANGELOG.md ## Unreleased` bullet. R1 surfaced one B-tier finding (F1, pre-existing I4 lychee stale link in cycle-360 surface, unrelated to α's #365 implementation surface); γ authored the fix as `c3a48741`. R2 APPROVED. Merge at `cbb8848b`.

## Implementation observations

- **Two-option design call collapsed to KISS.** γ scaffolded two paths: single cutoff bump vs. per-artifact era policy table. The failure shape (α-closeout / γ-closeout / CDD-Trace section all absent at the same v3.75/v3.76 boundary) collapses to one threshold, so three knobs all set to v3.77.0 would have been complexity without payoff. The per-artifact path is the only shape that would preserve "v3.74 strict" once v3.74.0 has zero observable misses; today the semantic shift from "v3.74 strict" → "v3.74 legacy" is invisible (cycle 327 ships clean). β observation N2 records this trade-off explicitly.
- **Comment / code paired in one commit.** Commit `1b7e77a1` carries the cutoff change *and* the helper's comment-block rewrite together. The prior comment named `v3.65.0 (epoch threshold)` and `3.74.0`; both literals would have been stale within minutes if the comment had been a follow-up commit. AC4 ("Comment / code coherence") was authored as a gate against exactly the "I forgot to update the comment" failure mode, and pairing kept that gate cheap to satisfy.
- **Boundary tests as positive/negative pair across the cutoff.** Tests 11 (v3.75.0 missing α-closeout → warn), 12 (v3.76.0 missing CDD-Trace → warn), and 13 (v3.77.0 missing α/γ close-outs → fail with exit-1 and the literal `❌ alpha-closeout.md` / `❌ gamma-closeout.md` wire strings) form the regression pin. The negative oracle (test 13) is the single highest-precision AC — any future cutoff drift past v3.77.0 fires it on real wire format, not a synthetic substring.
- **Fixture relocated, not invented.** The `incomplete-triadic` fixture was the prior strict-path proof at `3.75.0/200/`. Renamed (git rename score R100) to `3.78.0/200/` rather than creating a new post-cutoff fixture. This keeps `test-fixtures.sh` test 2 invariant (it names the fixture path's existence, not the version literal) while moving the proof's era to a still-strict era. Two-file pure rename, zero content delta.
- **Test runner output line count, not manual enumeration.** Per α/SKILL.md §2.6 row 13, the assertion count quoted in §AC5 is pasted from `bash test-cn-cdd-verify.sh | tail` (`ran 39 assertions: 39 passed, 0 failed`), not counted by reading the script.

## Friction log

- **In-flight triadic cycles classify as small-change during α-side authoring.** While α was writing self-coherence and before β had committed `beta-review.md`, the validator's `classify_cycle_type` routed `.cdd/unreleased/365/` to the small-change path (because `beta-review.md` was absent) and ran the strict section validator against α's incrementally-committed self-coherence. The failure was transient (cleared once §Gap..§CDD-Trace were all on the branch), but every α-side incremental commit drove a transient `cdd-artifact-check` red on the cycle branch. Out of #365's scope per γ non-goal #2 ("do not change which artifacts are required for the current era"); the classifier predates and is orthogonal to the era cutoff. Pattern: classifier inputs are which post-merge artifacts are present, but the triadic-vs-small-change signal is actually `gamma-scaffold.md` presence — a future cycle could route classification on that signal so in-flight triadic cycles are recognized before β has written.
- **Section-name slug drift between `alpha/SKILL.md` and the validator.** `alpha/SKILL.md` §2.5 names the self-coherence step-7 section `**§CDD-Trace**` (hyphen). The validator's grep (`cn-cdd-verify` line 491) requires `## CDD Trace` (space). v3.76.0 cycle 362 used the hyphen form and was flagged. Resolved within #365 scope by using the validator-compatible space form for the rendered heading while keeping the manifest's section slug `CDD-Trace` (matching α/SKILL.md). The dual form works but encodes a drift between two surfaces that name the same thing — out of #365's scope per γ non-goal #2.
- **F1 authorship landed at γ, not α.** β R1 named F1 (pre-existing I4 lychee stale link in `.cdd/releases/3.75.0/360/beta-closeout.md:54` after the cycle directory moved from `.cdd/unreleased/360/` to `.cdd/releases/3.75.0/360/` changed what `../../../` resolved to) and offered two recovery paths: (a) one-character path-depth fix on cycle/365, or (b) defer-by-design-scope via a separate γ-filed tracking issue. The fix that actually shipped was path (a) but authored by γ (cherry-pick of σ's `76fb7780` as `c3a48741`) rather than by α. The content boundary held — the broken file is from cycle 360, unrelated to α's #365 surface, so α had no authoring claim on it; γ as orchestrator was the right authorship. β/SKILL.md Role Rule 1 ("β does not author the fix it judges") was honored even though the fix-author was γ rather than α. β R2 verdict and closeout both record this explicitly.

## Patterns observed

- **Era cutoffs concentrate at one literal.** γ's peer enumeration found exactly one place encoding the cutoff (`is_legacy_version`'s `minor < 74` literal) and exactly one place where comment-vs-code had drifted. Pattern: a single-knob era policy stays single-knob through cutoff bumps as long as the warn-string sites read the helper's verdict rather than re-test the literal. Both `check_triadic_artifacts` and `validate_artifact_sections` already followed that shape; only the warn-string text (`(pre-v3.65.0)` vs `(pre-v3.77.0)`) needed an update at each call site. The literal `3.74` / `3.65` did not exist outside the helper + its comment, so the diff stayed narrow.
- **Cycle scope can absorb an incidental pre-existing red CI surface.** F1 was pre-existing on `origin/main@16aaef89` (Build red at base SHA, not introduced by #365). Rule 3.10 (CI-green-on-review-SHA) treated it as binding regardless, and the one-character fix landed inside the cycle rather than as a separate cycle. Pattern: when a cycle is the next push past a pre-existing red, the cycle is the right place to clear it if the fix is mechanical and orthogonal in content; β can name the disposition path (γ-files-issue vs in-cycle fix) without forcing the bigger refactor. The two-path framing in R1's F1 evidence column made the cheap path discoverable.
- **In-flight transient CI noise on α-side incremental commits.** The validator's small-change vs triadic classifier reads post-merge artifact presence, so in-flight triadic cycles trip a transient strict-section failure during α-side authoring. The signal is not actionable mid-cycle (β has not yet written `beta-review.md`), and it clears the moment α's §CDD-Trace section commits. Pattern: classifier inputs that depend on post-merge artifacts produce ephemeral red during the authoring window between scaffold and review-readiness signal.

## CDD Trace evidence

- §Gap: `.cdd/unreleased/365/self-coherence.md` §Gap.
- Implementation commits: `1b7e77a1` (validator cutoff + comment + warn strings + fixture rename), `722fe414` (boundary tests + fixture writers), `1536a843` (CHANGELOG bullet).
- Self-coherence commits: `42e608e2` (§Gap + §Skills), `757776ef` (§ACs), `d5561963` (§Self-check + §Debt + §CDD-Trace), then review-readiness section appended.
- Pre-review gate evidence: `self-coherence.md` §Review-readiness.
- β R1 verdict: REQUEST CHANGES on F1 (pre-existing I4 lychee stale link, cycle-360 surface unrelated to #365); β R2 verdict at `8810e4a5` — APPROVED.
- F1 fix: `c3a48741` (γ-authored cherry-pick of `76fb7780`, one extra `..` in `.cdd/releases/3.75.0/360/beta-closeout.md:54`).
- Merge: `cbb8848b` (`Closes #365`).

No undeclared debt. The two §Debt items in `self-coherence.md` (validator's `.cdd/unreleased/$issue_num/` text in released-cycle warnings, and `## CDD Trace` vs `## CDD-Trace` slug drift) carry forward as β observations N1 / N3 and γ PRA candidates; both are pre-existing, both are out of #365's scope per γ non-goal #2.
