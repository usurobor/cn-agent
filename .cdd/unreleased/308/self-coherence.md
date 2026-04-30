---
cycle: 308
role: α
issue: "skill(cdd/review): split phase 2 into review/issue-contract, review/diff-context, review/architecture siblings"
---

# Self-Coherence — Cycle #308

## §Gap

**Issue:** #308 — skill(cdd/review): split phase 2 into review/issue-contract, review/diff-context, review/architecture siblings

**Version/mode:** content-preserving split; no content rewriting; CTB v0.1 frontmatter required on each new sub-skill.

**Gap being closed:** `review/implementation/SKILL.md` bundles three cognitive modes (issue-contract walk §2.0, diff/context inspection §2.1.1–§2.1.13, architecture/design check §2.2) in one file. The orchestrator (`review/SKILL.md`) already names the deferred-split as the correct decomposition. This cycle lands that decomposition: three new siblings each owning one reason to change, the monolith deleted, orchestrator updated, M2-review kata reference rewired.

**Invariant:** One reason to change per skill (design-skill discipline). Each mode changes for a different reason — issue-contract walks change when CDD artifact rules change; diff/context heuristics change when implementation review patterns change; architecture checks change when design boundaries shift.
