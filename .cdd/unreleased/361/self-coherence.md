---
cycle: 361
issue: "https://github.com/usurobor/cnos/issues/361"
role: alpha
sections:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, ACs]
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

## §ACs

| # | AC | Status | Evidence |
|---|----|--------|----------|
| AC1 | review/SKILL.md §3.3/3.4 contains explicit verdict-shape lint rules: `APPROVED` + unresolved findings = invalid; `APPROVED` + "conditional/pending/modulo" qualifier = invalid | met | review/SKILL.md §3.4a bullet 1 (`APPROVED` + unresolved findings at any severity, cross-refs 3.3); §3.4a bullet 2 enumerates the six qualifier tokens `conditional`, `pending`, `modulo`, `subject to`, `assuming`, `provisional on`, `with follow-up` |
| AC2 | Recovery path documented: conditional becomes RC with conditions as required-fix findings | met | review/SKILL.md §3.4a Recovery paragraph: "any conditional or split shape becomes `REQUEST CHANGES`; each named condition is reformatted as a required-fix finding at its severity. The reviewer re-emits a single terminal verdict in the next round once findings clear." |
| AC3 | Reviewers must select one terminal verdict per round — no split verdicts | met | review/SKILL.md §3.4a bullet 3 (Split verdicts): "two terminal verdicts in one round … One round, one decision." Plus checklist line: "Verdict-shape lint passed (no `APPROVED` + unresolved findings; no conditional qualifier; no split verdict) per rule 3.4a." |
