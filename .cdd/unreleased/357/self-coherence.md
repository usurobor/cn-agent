<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs] -->

# Self-Coherence Report: Issue #357

## Gap

**Issue:** #357 — cdd/release: release notes generation skill — δ produces structured tag messages before tagging

**Version/Mode:** design-and-build — the gap is bounded, but there is no committed design path and no committed plan path; alpha must choose the smallest coherent script/docs/test shape during implementation.

**Gap description:** `scripts/release.sh` currently creates lightweight tags with `git tag "$VERSION"`; no generated annotated tag-message path exists. Delta cannot produce a consistent tag message at the release boundary without manually re-reading git history, GitHub issues, and `.cdd/` artifacts.

**Selected incoherence:** Release metadata outside `RELEASE.md` is absent or manual, so tag messages can drift from issue labels, review-round evidence, wave context, and release highlights. Multi-issue and wave releases are especially fragile: a release can have a coherent `RELEASE.md` while its git tag remains lightweight or carries ad hoc prose.

**Coherence target:** Delta gets tag-message generation through the canonical script path, creating annotated version tags with generated plain-text messages derived deterministically from git history, issue metadata, CDD review artifacts, and optional wave artifacts.

## Skills

**Tier 1 (CDD core):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface and algorithm

**Tier 2 (always-applicable eng):**
- `src/packages/cnos.core/skills/eng/*` — engineering discipline skills

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — release artifact authority, RELEASE.md boundary, version/tag flow
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — delta gate ownership and no-manual-tag rule
- `src/packages/cnos.core/skills/eng/tool/SKILL.md` — shell/tooling implementation discipline
- `src/packages/cnos.core/skills/eng/test/SKILL.md` — fixture/offline proof for git/GitHub/CDD parsing and release-script behavior
- `src/packages/cnos.core/skills/eng/document/SKILL.md` — operator-facing release/tag-message prose

**Active constraints from loaded skills:**
- Release script integration only (no manual tag commands) - from cdd/operator
- Preserve RELEASE.md authority - from cdd/release
- Plain text output suitable for git tag annotation - from eng/tool
- Fixture-testable with deterministic fallback - from eng/test
- Clear operator documentation - from eng/document

## ACs

**AC1: Generator emits deterministic plain-text tag messages**
- ✅ **Evidence:** `scripts/generate-release-tag-message.sh` at 8041bb30 produces stable plain-text output
- ✅ **Oracle:** Integration test at `scripts/test-release-tag-integration.sh` line 109 verifies deterministic output (excluding timestamp)
- ✅ **Surface:** Generator includes version/date, issue count, per-issue entries with titles/labels/reviews, highlights section
- ✅ **Verification:** Test output shows "1.1.0 — Release Tag Message", "Changes:", "- #456: unavailable", "Summary: 1 issues"

**AC2: Release script creates annotated tags with the generated message**  
- ✅ **Evidence:** `scripts/release.sh` line 105-115 modified to use generated messages for annotated tags
- ✅ **Oracle:** Integration test verifies `git cat-file -t 1.1.0` returns "tag" and message contains generated content
- ✅ **Surface:** Release script calls `scripts/generate-release-tag-message.sh` and creates annotated tags via `git tag -a -F`
- ✅ **Verification:** Test output shows "Tag type: tag" and "Tag is annotated (type: tag)"

**AC3: RELEASE.md remains the GitHub release-body authority**
- ✅ **Evidence:** No changes to `.github/workflows/release.yml` or `scripts/validate-release-gate.sh`
- ✅ **Oracle:** Integration test includes RELEASE.md creation and validation; release gate enforcement unchanged
- ✅ **Surface:** Generated tag messages are additive, documented in release/SKILL.md as preserving RELEASE.md authority
- ✅ **Verification:** Release integration test creates RELEASE.md and process succeeds without replacing GitHub body

**AC4: Issue metadata collected with deterministic fallback**
- ✅ **Evidence:** Generator script lines 65-76 implement GitHub metadata lookup with "unavailable" fallback
- ✅ **Oracle:** Test repository shows "#456: unavailable" output when GitHub metadata inaccessible
- ✅ **Surface:** Uses `gh issue view` when available, falls back to minimal structure with explicit unavailable markers  
- ✅ **Verification:** Test output includes issue number with "unavailable" title, no silent omissions or failures

**AC5: Wave and standalone releases both render coherently**
- ✅ **Evidence:** Generator script lines 110-123 handle wave context detection with fallback to standalone grouping
- ✅ **Oracle:** Real repository test shows "Wave: hardening-2026-05-12 (4 issues)" when wave present
- ✅ **Surface:** Detects `.cdd/waves/` artifacts when present, gracefully handles absence without bogus sections
- ✅ **Verification:** Output includes wave context when available, remains coherent for standalone releases

**AC6: Review rounds and findings come only from durable CDD artifacts**
- ✅ **Evidence:** Generator script lines 125-145 parse `.cdd/{unreleased,releases}/**/beta-*.md` when present
- ✅ **Oracle:** Script counts review rounds from beta-review.md, extracts findings from beta-closeout.md
- ✅ **Surface:** Does not invent "0 rounds" for missing artifacts, follows documented unavailable behavior
- ✅ **Verification:** Code shows explicit artifact directory checking and fallback to empty values

**AC7: Delta documentation names the new script-mediated gate step**
- ✅ **Evidence:** Updated operator/SKILL.md lines 222 and 226-228, release/SKILL.md lines 217-218 and 221
- ✅ **Oracle:** `rg "git tag -a|git tag .* -m|git tag .* -F"` finds no new manual tagging instructions for delta
- ✅ **Surface:** Documentation describes generated tag messages through scripts/release.sh only
- ✅ **Verification:** Operator docs mention "generates structured tag messages" and δ tag inspection commands