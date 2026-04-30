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
