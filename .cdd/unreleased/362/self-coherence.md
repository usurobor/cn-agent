---
cycle: 362
issue: "https://github.com/usurobor/cnos/issues/362"
role: alpha
manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap]
---

# α self-coherence — #362

## Gap

Issue #362: `cap/SKILL.md` §UIE has no communication gate distinguishing questions from instructions. The U step is internal — agent "understands" and immediately executes, making understanding invisible. "What is X?" silently becomes "fix X." Observed 5+ times in a single δ-session (tsc #49 wave + cnos #359-361 wave).

Version: skill-patch to `src/packages/cnos.core/skills/agent/cap/SKILL.md` — additive subsection, no contract break.
Mode: MCA — add the missing gate to the skill so it constrains future authorship.
