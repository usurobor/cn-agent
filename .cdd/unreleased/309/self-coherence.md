---
cycle: 309
role: alpha
---

# Self-Coherence — Cycle 309

## §Gap

**Issue:** #309 — skill(eng): create eng/troubleshoot — live diagnosis skill for environmental and runtime failures

**Version / mode:** cnos · CDD triadic cycle · α implementation

**Gap:** No live troubleshooting skill existed. When an agent hits a live failure mid-task, diagnosis was improvised. The 2026-04-30 identity-rotation dispatch test produced five β dispatch failures across three root-cause classes (OOM kill, `gh` GraphQL error, background-process lifecycle kill), each diagnosed ad-hoc with a wrong first hypothesis. A structured skill for live diagnosis would have forced process state and kernel log checks before model/token speculation.

**Target artifact:** `src/packages/cnos.eng/skills/eng/troubleshoot/SKILL.md`

**Not in scope:** Automated diagnostic tooling, platform-specific runbooks, changes to `eng/rca`, new commands, dispatch behavior changes.

## §Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract
- `CDD.md` — canonical lifecycle

**Tier 2:**
- `src/packages/cnos.core/skills/write/SKILL.md` — writing standard (top-down, point first, no clutter)

**Tier 3 (issue-specified):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — classification, formula, kata surface requirements
- `src/packages/cnos.core/skills/write/SKILL.md` — governing question, levels of structure, direct rules
- `src/packages/cnos.core/skills/design/SKILL.md` — boundary discipline; skill vs runbook vs reference; surface separation
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — invariant, oracle, positive/negative proof model
- `src/packages/cnos.eng/skills/eng/rca/SKILL.md` — adjacent sibling boundary; RCA handoff definition

All five Tier 3 skills were loaded and read before drafting began.
