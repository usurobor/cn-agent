---
name: cross-repo
description: Compatibility pointer. Cross-repo handoff doctrine has moved to cnos.handoff/skills/handoff/cross-repo/SKILL.md as of cnos#416 (Sub 2 of cnos#404 handoff extraction).
artifact_class: pointer
visibility: internal
parent: cdd
status: moved
canonical: cnos.handoff/skills/handoff/cross-repo/SKILL.md
---

# Cross-repo doctrine — moved

The canonical home for cross-repo handoff doctrine is now:

**[`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md)**

This includes:
- STATUS state machine (8 events: drafted / submitted / accepted / modified / rejected / landed / withdrawn / revised)
- Directional cases (a inbound, b outbound, c bilateral, d operator-pending)
- LINEAGE.md schemas per case
- Feedback-patch format
- Bundle archival rule
- Hat-collapse attribution
- Known protocol edge cases

cnos.cdd is a **consumer** of cross-repo doctrine, not an owner. See cnos#404 for the extraction rationale.

This pointer is preserved for backward compatibility with cross-references that have not yet been re-pointed. A future cleanup cycle may remove this file entirely once all consumers cite cnos.handoff directly.
