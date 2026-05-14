---
cycle: 362
issue: "https://github.com/usurobor/cnos/issues/362"
role: alpha
manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills]
---

# α self-coherence — #362

## Gap

Issue #362: `cap/SKILL.md` §UIE has no communication gate distinguishing questions from instructions. The U step is internal — agent "understands" and immediately executes, making understanding invisible. "What is X?" silently becomes "fix X." Observed 5+ times in a single δ-session (tsc #49 wave + cnos #359-361 wave).

Version: skill-patch to `src/packages/cnos.core/skills/agent/cap/SKILL.md` — additive subsection, no contract break.
Mode: MCA — add the missing gate to the skill so it constrains future authorship.

## Skills

- **Tier 1** — `CDD.md`, `cdd/alpha/SKILL.md`
- **Tier 2** — none required (skill-patch, not engineering bundle)
- **Tier 3** — `cnos.core/skills/write/SKILL.md` (per γ scaffold — the patch is prose; writing rules govern)

Applied: `write` §1.3 governing question (the new subsection answers one question: how does the agent classify input?), §2.4 lead with the point (each bullet leads with the classification), §3.12 small contrastive examples (four ❌/✅ pairs), §3.1 state what it is (the failure-mode paragraph names the mode positively, not as "don't skip understanding").
