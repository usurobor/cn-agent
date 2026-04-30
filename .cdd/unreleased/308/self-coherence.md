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

---

## §Skills

**Tier 1:**
- `CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2:**
- `src/packages/cnos.eng/skills/eng/` (always-applicable engineering skills) — not loaded explicitly; no language-specific implementation work in this cycle (pure skill authoring / YAML/Markdown)

**Tier 3 (issue-declared):**
- `cnos.core/skills/design/SKILL.md` — split design with one reason to change per skill; governs the boundary call for this decomposition
- `cnos.cdd/skills/cdd/review/SKILL.md` — parent skill being modified; loaded to verify the orchestrator's phase contract stays intact

**Applied constraints from Tier 3:**
- design/SKILL.md §2.1: decompose by reason to change → each new sub-skill owns one cognitive mode
- design/SKILL.md §3.2: one source of truth per fact → orchestrator owns phase order + verdict; sub-skills own their cognitive mode
- design/SKILL.md §2.4: truthful interfaces → each sub-skill's `inputs`/`outputs` declares exactly what it consumes/produces
- CTB v0.1 `LANGUAGE-SPEC.md` §2: hard-gate frontmatter fields on every new sub-skill — `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, `triggers`, `scope`
