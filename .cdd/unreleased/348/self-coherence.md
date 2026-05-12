# Self-Coherence — Cycle #348

## §Gap

**Issue:** #348 — ROLES.md §4: hats-vs-actors principle

**Gap:** ROLES.md §3 names "actor collapse rule" as a per-protocol instantiation field, and operator/SKILL.md §5.2 documents the γ/δ collapse as a dispatch configuration. But neither states the general principle that roles are behavioral contracts (hats), not entity slots (actors). The principle that independence is the collapse constraint (not ceremony or headcount) exists but is not captured in any cnos surface.

**Version/Mode:** docs-only, design-and-build

**Selected change:** Add ROLES.md §4 stating the hats-vs-actors principle with worked examples and cross-references from operator/SKILL.md §5.2 and epsilon/SKILL.md collapse notes.

## §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface and algorithm

**Tier 2 (always-applicable engineering):**
- None required for this docs-only cycle

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/write/SKILL.md` — writing skill for coherent prose

## §ACs

**AC1: ROLES.md contains a section (§4 or amendment) stating the hats-vs-actors principle**
- **Status:** Implemented
- **Evidence:** New ROLES.md §4 "Hats vs actors: roles as behavioral contracts" added after §3 and before existing §4 (renumbered to §5). Section states "Roles are behavioral contracts (hats), not entity slots (actors)."

**AC2: The principle names independence as the collapse constraint (not ceremony, not headcount)**
- **Status:** Implemented  
- **Evidence:** ROLES.md §4 includes "Independence, not headcount" subsection stating "The constraint that α≠β is not ceremonial separation or bureaucratic headcount. The constraint is independence: review that is not independent of implementation is not review."

**AC3: At least one worked example (δ=γ safe because..., α=β unsafe because...)**
- **Status:** Implemented
- **Evidence:** ROLES.md §4 "Worked examples" subsection provides three examples:
  - δ=γ collapse (safe) — removes communication hop without removing judgment gate
  - α=β collapse (never safe) — destroys review independence entirely
  - ε=δ collapse (safe in small-protocol regimes) — no independence boundary compromised

**AC4: Cross-reference from operator/SKILL.md §5.2 to the new section**
- **Status:** Pending implementation
- **Evidence:** TBD — will add reference from existing "δ=γ collapse" discussion to new ROLES.md §4

**AC5: Cross-reference from epsilon/SKILL.md collapse note to the new section**
- **Status:** Pending implementation  
- **Evidence:** TBD — will add reference from existing "collapse onto δ permitted" note to new ROLES.md §4