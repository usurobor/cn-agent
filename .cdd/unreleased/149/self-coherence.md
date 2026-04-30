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

---

## §ACs

### AC1: Skill-loading rule in SOUL.md §2.1 Observation

**Invariant:** An agent reading SOUL.md §2 sees the skill-loading requirement before reaching §2.2 Action.
**Oracle:** §2.1 contains a sub-rule stating "Before acting, identify which skills this action requires and load them."
**Positive:** Sub-rule explicitly names that memory of a skill ≠ loading the skill.
**Negative:** A rule that only says "consider loading skills" does not satisfy the AC.

**Evidence:** SOUL.md §2.1 now contains:

> **Skill-loading gate:** Before acting, identify which skills this action requires and load them. Memory of a skill is not the same as loading it — skill files are the constraint, and their content drifts from agent memory over time. Name the loaded skills before proceeding to Action.

The rule is imperative ("identify … and load them"), appears inside §2.1 before §2.2, and explicitly states "Memory of a skill is not the same as loading it." AC1: **met**.

### AC2: ❌/✅ examples for skill loading

**Invariant:** The rule shows what does and does not satisfy it.
**Oracle:** At least one ❌/✅ pair shows acting on memory of a skill versus loading the skill file.

**Evidence:** Added immediately after the skill-loading gate:

```
- ❌ Act on memory of what the release skill says.
- ✅ Load `release/SKILL.md`, name it as loaded, then act.
```

Same situation (release skill), contrastive (memory vs file), imperative. AC2: **met**.

### AC3: Visible-load convention named

**Invariant:** When the agent loads skills, the load is observable to the operator.
**Oracle:** Sub-rule states the agent names loaded skills before proceeding to Action.
**Negative:** Silent load does not satisfy this AC.

**Evidence:** The skill-loading gate includes: "Name the loaded skills before proceeding to Action." AC3: **met**.
