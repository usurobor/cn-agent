---
cycle: 361
issue: "https://github.com/usurobor/cnos/issues/361"
role: alpha
sections:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills]
---

# α Self-coherence — #361

## §Gap

**Issue**: cdd: review rules 3.3/3.4 — verdict-shape lint (no conditional APPROVEDs)

β@S4 in tsc #53 issued "APPROVED with 3 unresolved C findings + conditional language." Rules 3.3 (APPROVED = zero unresolved findings) and 3.4 (verdict before details) both forbid this shape — but the ban is implicit. No explicit lint enumerates invalid verdict shapes, names the conditional qualifiers that smuggle in unresolved findings, or rules out split verdicts.

**Mode**: skill-patch, single-file diff to `src/packages/cnos.cdd/skills/cdd/review/SKILL.md`.

## §Skills

- **Tier 1** — `CDD.md`, `cdd/alpha/SKILL.md` (this role surface). No lifecycle sub-skill needed (design/plan marked not-required: single-file skill patch, no impact graph beyond the file itself).
- **Tier 2** — none. The diff is normative-prose-only inside one role surface; no Tier 2 engineering bundle applies.
- **Tier 3** — `cnos.core/skills/write/SKILL.md` (per dispatch). Governs prose: front-load the point (3.4), one fact one home (3.3), specifics over abstractions (3.8), state what it is (3.1).
