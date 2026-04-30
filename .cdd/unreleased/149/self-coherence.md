## §Gap

**Issue:** #149 — UIE must include skill loading as part of Understand before Execute
**Version/mode:** MCA (small-change path — single-file prose addition, §1.2)

SOUL.md §2.1 Observation had no explicit requirement to identify and load required skills before acting. An agent could complete the Observation step and proceed to Action using stale memory of a skill rather than the current skill file. Three concrete failures on 2026-04-02 demonstrated the gap: heartbeat without inbox skill (escalation leak), release without release skill (skipped CDD bookends), review without review skill (operator had to prompt mid-task).

The fix: add a **skill-loading gate** to §2.1 Observation, adjacent to the falsification gate, that makes skill loading an explicit Observation sub-requirement with visible-load convention and ❌/✅ examples.

---

## §Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 3 (issue-named):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — meta-skill discipline: imperative form, ❌/✅ pairs, governing question
- `src/packages/cnos.core/skills/write/SKILL.md` — house style for operator-facing load-bearing prose
