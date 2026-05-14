---
cycle: 362
issue: "https://github.com/usurobor/cnos/issues/362"
role: alpha
manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, ACs]
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

## ACs

**AC1** — `cap/SKILL.md` §UIE contains an explicit rule: if the input is a question, the U step is the answer — deliver it before moving to I or E.

Evidence: `cap/SKILL.md` §1.5 first bullet — *"**Question** ("what is X?", "how does Y work?", "is this Z?") — the U step is the answer. Deliver the answer before moving to I or E. The answer makes your model of the situation visible, so the operator can correct it before any action fires."* Located in §1 Understand (the U step), so the placement and the text agree.

**AC2** — `cap/SKILL.md` contains a question-vs-instruction distinction: questions get answers, instructions get actions. "What is X?" is never "fix X."

Evidence: `cap/SKILL.md` §1.5 — section title is *"Classify the input — question or instruction"*, and the two bullets explicitly split the two cases: Question → answer is the U step; Instruction → proceed through U, I, E as normal. The third paragraph handles ambiguous input by routing it to the question side first. The final ❌/✅ pair literally instantiates the "what is X?" example: *"❌ Operator: 'what is the dispatch protocol?' → agent edits the dispatch skill file"* vs *"✅ Operator: 'what is the dispatch protocol?' → agent answers in prose, then asks whether a change is requested."*

**AC3** — the rule names the failure mode it prevents: invisible understanding that skips to action, making the operator unable to verify the agent's model before work begins.

Evidence: `cap/SKILL.md` §1.5 paragraph beginning *"Failure mode: **invisible understanding that skips to action.**"* — names the mode explicitly in bold and explains the mechanism (model forms internally, action fires before operator can verify or correct). Closes with the concrete instantiation: *"'What is X?' silently becomes 'fix X.'"*
