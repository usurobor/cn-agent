# Self-Coherence — Issue #324

## §Gap

**Issue:** [#324 — skill(cdd/issue): split issue skill into focused subskills and add label taxonomy](https://github.com/usurobor/cnos/issues/324)

**Version/mode:** CDD 3.15.0 — MCA (skill refactor + new label taxonomy subskill)

**Named incoherence:**

`cdd/issue/SKILL.md` has grown into a large all-in-one skill carrying label taxonomy, AC/proof rules, constraint strata, exception discipline, path resolution, cross-surface updates, handoff checks, and katas inline. The root skill loads too much context when only one concern is needed and has multiple reasons to change. New label rules would make it larger if added directly.

**Target state:** Root `issue/SKILL.md` becomes a small orchestrator that delegates to focused subskills: `labels/`, `contract/`, `proof/`, `constraints/`. Label taxonomy lives in `issue/labels/SKILL.md` and requires exactly one kind label and one priority label per issue.

**Design:** Not required — subskill boundaries are specified in the issue. No novel architecture decisions.

**Plan:** Not required — five files in sequence, no complex dependencies, no multi-step execution ordering.
