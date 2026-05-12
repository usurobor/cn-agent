# Self-Coherence Report — Issue #356

## Gap

**Issue #356**: cdd/operator: δ release gate must block until CI green — not just poll and report

**Version/Mode**: CDD lifecycle gap — δ (operator) role lacks proper CI failure blocking and recovery protocol.

**Problem**: Issue #354 landed the rule that δ polls release CI after tag push. But δ currently reports the result and moves on. If CI is red, nothing happens — the release is declared incomplete but no recovery is prescribed. β blocks merge on red (rule 3.10). γ blocks close-out on red (§2.7). δ should block release on red — same pattern.

**Gap Class**: Role contract gap in operator skill — missing release gate CI failure handling protocol.