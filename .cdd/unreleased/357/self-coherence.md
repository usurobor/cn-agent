<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap] -->

# Self-Coherence Report: Issue #357

## Gap

**Issue:** #357 — cdd/release: release notes generation skill — δ produces structured tag messages before tagging

**Version/Mode:** design-and-build — the gap is bounded, but there is no committed design path and no committed plan path; alpha must choose the smallest coherent script/docs/test shape during implementation.

**Gap description:** `scripts/release.sh` currently creates lightweight tags with `git tag "$VERSION"`; no generated annotated tag-message path exists. Delta cannot produce a consistent tag message at the release boundary without manually re-reading git history, GitHub issues, and `.cdd/` artifacts.

**Selected incoherence:** Release metadata outside `RELEASE.md` is absent or manual, so tag messages can drift from issue labels, review-round evidence, wave context, and release highlights. Multi-issue and wave releases are especially fragile: a release can have a coherent `RELEASE.md` while its git tag remains lightweight or carries ad hoc prose.

**Coherence target:** Delta gets tag-message generation through the canonical script path, creating annotated version tags with generated plain-text messages derived deterministically from git history, issue metadata, CDD review artifacts, and optional wave artifacts.