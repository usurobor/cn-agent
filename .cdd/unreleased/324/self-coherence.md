# Self-Coherence — Issue #324

## §Gap

**Issue:** [#324 — skill(cdd/issue): split issue skill into focused subskills and add label taxonomy](https://github.com/usurobor/cnos/issues/324)

**Version/mode:** CDD 3.15.0 — MCA (skill refactor + new label taxonomy subskill)

**Named incoherence:**

`cdd/issue/SKILL.md` has grown into a large all-in-one skill carrying label taxonomy, AC/proof rules, constraint strata, exception discipline, path resolution, cross-surface updates, handoff checks, and katas inline. The root skill loads too much context when only one concern is needed and has multiple reasons to change. New label rules would make it larger if added directly.

**Target state:** Root `issue/SKILL.md` becomes a small orchestrator that delegates to focused subskills: `labels/`, `contract/`, `proof/`, `constraints/`. Label taxonomy lives in `issue/labels/SKILL.md` and requires exactly one kind label and one priority label per issue.

**Design:** Not required — subskill boundaries are specified in the issue. No novel architecture decisions.

**Plan:** Not required — five files in sequence, no complex dependencies, no multi-step execution ordering.

---

## §Skills

**Tier 1a — CDD authority:**
- `CDD.md` v3.15.0
- `src/packages/cnos.cdd/skills/cdd/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`

**Tier 1b — CDD lifecycle:**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (interpreting issue ACs and quality)

**Tier 2 — Writing bundle:**
- Applied: write/SKILL.md constraints — one governing question per artifact, no repetition, front-load the point

**Tier 3 — Issue-specific:**
- `src/packages/cnos.core/skills/write/SKILL.md`
- `src/packages/cnos.core/skills/design/SKILL.md`
- `src/packages/cnos.core/skills/skill/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` (loaded per issue body §Skills to load)
