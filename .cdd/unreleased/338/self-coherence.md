## Self-Coherence — cycle #338

### §Gap

**Issue:** #338 — cdd: Add §1.6c — dispatch sizing, prompt scope, and commit checkpoints
**Version / mode:** docs-only (`design-and-build`; §2.5b disconnect path)
**Branch:** `cycle/338`
**Gap being closed:** `CDD.md §1.6` had no rule for scaling the initial dispatch timeout budget with work complexity, no prompt-scope guidance matching load to task, and no commit-checkpoint mandate. This created a recoverable-only-by-luck failure class: agents SIGTERMed without committing lose all in-progress work. Pattern observed in N=4 instances (cnos #335, TSC supercycle 3/5 re-dispatches).

This cycle also validates its own heuristic: 5 ACs × 120s + 300s floor = 900s. Budget set to 600s per the dispatch prompt — this is below the heuristic. β will apply recursive coherence check per the proof plan. Known debt: budget discrepancy declared in §Debt.

### §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle; §1.6 coordination model is the change target
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface and load order (this file)

**Tier 2 (always-applicable eng):**
- `src/packages/cnos.eng/skills/eng/writing` — prose patches to skill files (docs-only cycle)

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-program/frontmatter coherence on the three modified files (constraint: no frontmatter changes)

**Active design constraints loaded:**
- No frontmatter changes to any of the three modified SKILL.md files
- No new sub-skills — §timeout-recovery lives inline in operator/SKILL.md
- Heuristic constants marked "initial; refine with telemetry"
- Cross-refs are normative: §1.6c ↔ §1.6b ↔ operator §timeout-recovery ↔ post-release §4
- Empirical citations reproducible (close-out paths)
- No retroactive changes to §1.6a or §1.6b


