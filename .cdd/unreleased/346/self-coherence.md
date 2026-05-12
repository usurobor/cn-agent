# Self-Coherence — Cycle #346

## Gap

`epsilon/SKILL.md §1` conditioned `cdd-iteration.md` production on the
existence of `cdd-*-gap` findings, creating an ambiguity: a missing file
could indicate either a clean cycle or a skipped close-out.
`activation/SKILL.md §22` prescribes unconditional production ("every
cycle"). The two authoritative sources were in conflict; this cycle
harmonizes epsilon §1 with activation §22.

## ACs

| AC | Oracle | Result |
|---|---|---|
| AC1 | `rg 'Empty cycles produce no file' epsilon/SKILL.md` → 0 hits | PASS |
| AC2 | `rg 'every cycle' epsilon/SKILL.md` → ≥1 hit in §1 | PASS |
| AC3 | epsilon §1 and activation §22 both prescribe every-cycle production | PASS |

## CDD Trace

- **Source of truth:** `activation/SKILL.md §22` (prescribes every cycle)
- **File changed:** `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` §1, second paragraph
- **Change type:** docs-only, no behavior change to other artifacts

## Review-Readiness

- AC1 oracle: `rg 'Empty cycles produce no file' src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` → exit 1 (0 hits) ✓
- AC2 oracle: `rg 'every cycle' src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` → 1 hit in §1 ✓
- AC3: epsilon §1 and activation §22 are now consistent — both prescribe unconditional every-cycle production ✓
- Author email: `alpha@cdd.cnos` ✓
- No other files modified (only `epsilon/SKILL.md` and this artifact) ✓

Alpha signals **review-ready**.
