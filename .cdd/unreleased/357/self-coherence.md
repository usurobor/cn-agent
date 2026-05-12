<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills] -->

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