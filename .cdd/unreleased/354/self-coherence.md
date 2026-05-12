# Self-Coherence — Cycle #354

## §Gap

**Issue:** #354 — δ must poll release CI after tag push — no silent release failures

**Version:** 3.15.0

**Mode:** docs-only

**Gap summary:** δ pushes a tag, the release workflow fires, and nobody watches it. Release failures are silent until a human notices. Need δ to poll release CI after `git push origin <tag>` — the same way β polls cycle CI before APPROVED (rule 3.10) and γ polls post-merge CI before close-out (§2.7). The tag push is not the end of the release gate — the release workflow passing is.

**Evidence:** v3.66.0 release smoke failed (binary version-pin mismatch). Failure notification arrived in Telegram but no automated or role-prescribed response existed.

**Origin:** v3.66.0 release smoke failure (2026-05-12). Binary version-pin mismatch between compiled Go binary and package index.