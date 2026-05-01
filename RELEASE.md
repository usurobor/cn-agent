## Outcome

Coherence delta: C_Σ A (`α A`, `β A`, `γ A`) · **Level:** L5

Issue skill refactored from monolith to orchestrator + focused subskills. Label taxonomy added as a required contract surface. CI unblocked by adding R5-activate `run.sh`.

## Why it matters

The issue skill had grown too large for efficient LLM loading — one skill carrying structure, status truth, source maps, constraint strata, proof plans, path resolution, and label rules. Agents loading the root skill paid for all of it even when only writing a simple issue. The split lets agents load only what they need while preserving the full contract.

## Changed

- **skill(cdd/issue):** Split into orchestrator + 4 focused subskills (#324):
  - `issue/SKILL.md` — root orchestrator (~200 lines, down from ~800+)
  - `issue/labels/SKILL.md` — kind + priority label taxonomy
  - `issue/contract/SKILL.md` — problem, impact, status truth, source map, scope, non-goals
  - `issue/proof/SKILL.md` — AC shape, proof plan, positive/negative cases
  - `issue/constraints/SKILL.md` — constraint strata, exceptions, path resolution
- Every issue now requires exactly one kind label and one priority label.

## Fixed

- **ci:** Added R5-activate `run.sh` — Tier 2 kata suite was blocked since 3.70.0.

## Validation

- β approved R1, 0 findings, merged to main.
- All 12 ACs met. All 14 original issue-contract rules preserved across subskills.
- CTB v0.1 frontmatter valid on all new subskills.
- CI green on main after R5 `run.sh` fix.
