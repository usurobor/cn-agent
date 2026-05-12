# Self-Coherence — Cycle #354

## §Gap

**Issue:** #354 — δ must poll release CI after tag push — no silent release failures

**Version:** 3.15.0

**Mode:** docs-only

**Gap summary:** δ pushes a tag, the release workflow fires, and nobody watches it. Release failures are silent until a human notices. Need δ to poll release CI after `git push origin <tag>` — the same way β polls cycle CI before APPROVED (rule 3.10) and γ polls post-merge CI before close-out (§2.7). The tag push is not the end of the release gate — the release workflow passing is.

**Evidence:** v3.66.0 release smoke failed (binary version-pin mismatch). Failure notification arrived in Telegram but no automated or role-prescribed response existed.

**Origin:** v3.66.0 release smoke failure (2026-05-12). Binary version-pin mismatch between compiled Go binary and package index.

## §Skills

**Tier 1 (CDD lifecycle):**
- `cdd/alpha/SKILL.md` — α role surface and implementation algorithm
- `cdd/CDD.md` — canonical lifecycle and role contract

**Tier 2 (always-applicable eng):**
- N/A for docs-only mode

**Tier 3 (issue-specific):**
- `cnos.core/skills/write/SKILL.md` — coherent writing surface for docs modification