---
cycle: 362
issue: "https://github.com/usurobor/cnos/issues/362"
role: alpha
manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace]
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

## Self-check

**Did α push ambiguity onto β?** No. The patch is one new subsection in one file. Every AC maps to a quoted line in §1.5; β does not need to derive evidence from elsewhere.

**Is every claim backed by diff evidence?** Yes. AC1/AC2/AC3 each quote the exact §1.5 text they rely on. `git diff origin/main..HEAD -- src/packages/cnos.core/skills/agent/cap/SKILL.md` shows the 16-line addition; no other file is touched.

**Peer enumeration.** Skill-class peers were considered. `cap/SKILL.md` is a global agent skill, not a role skill — the contract is a behavioral rule for any agent reading the skill, not a lifecycle hand-off contract. Adjacent agent skills (`reflect`, `coherent`, `CLP` — listed in `cap/SKILL.md` §4) describe distinct concerns (post-hoc interpretation, output verification, pre-publish refinement) and do not duplicate the input-classification rule; none of them encode a question-vs-instruction gate that could drift from this one. Lifecycle skills (`alpha/`, `beta/`, `gamma/`, `review/`, `release/`, `post-release/`, `design/`, `plan/`, `issue/`) do not consume the UIE contract operationally — they constrain role behavior, not agent input handling. Conclusion: no peer surfaces require sibling updates.

**Harness audit.** No schema or runtime contract changed; no harness audit required.

**Post-patch re-audit (polyglot).** The diff is Markdown only. Re-audit: read the inserted §1.5 against the existing §1 structure — numbering continues from §1.4, the `---` separator before §2 remains in place, the ❌/✅ pair format matches surrounding sections (compare §1.2a, §1.4 in the unchanged file). Cross-reference check: the new section is referenced from nowhere else in the file, so no internal link drift. No additional toolchains apply.

**Artifact enumeration matches diff.** `git diff --stat origin/main..HEAD` lists two files: `src/packages/cnos.core/skills/agent/cap/SKILL.md` (the implementation, mapped in §ACs) and `.cdd/unreleased/362/self-coherence.md` (this file, the cycle artifact). Both declared.

**Caller-path trace.** No new module or function — skill-patch only. Row 12 of pre-review gate not applicable.

**Test assertion count.** No tests added — this patch modifies a prose skill that constrains agent authorship, not executable code. Row 13 of pre-review gate not applicable; no test runner output to paste.

## Debt

None. The patch closes the named gap; no follow-up cycles are implied by this change.

Adjacent observations (not debt, for γ's awareness): §1.5 placement is at the end of §1 rather than the top. The content correctly states the ordering ("the U step is the answer. Deliver the answer before moving to I or E"), and the existing §1 intro paragraph ("Read state before acting") frames the whole section consistently, so the trailing position does not undermine the rule. Renumbering §1.1-§1.4 to position the gate first was considered and rejected as larger diff for marginal semantic gain; the γ scope said "small diff to one file."

Provisional close-out fallback: this cycle uses bounded dispatch (§5.2 single-session δ-as-γ per γ scaffold). α will exit after signaling review-readiness; close-out follows via re-dispatch or γ-provisional path per `alpha/SKILL.md` §2.8.

## CDD-Trace

1. **Define** — issue #362 names the gap: `cap/SKILL.md` §UIE has no rule distinguishing questions from instructions; agent skips to action; "what is X?" becomes "fix X."
2. **Design** — not required. Single-file additive subsection to one skill, no impact graph beyond the file itself; γ scope explicitly said "small diff to one file."
3. **Plan** — not required. Implementation is one subsection insertion; sequencing is trivial.
4. **Tests** — not applicable. The artifact is a prose skill that constrains future agent authorship; correctness is verified by reading the new §1.5 against the three ACs, which §ACs above does explicitly.
5. **Code** — `src/packages/cnos.core/skills/agent/cap/SKILL.md` §1.5 added (commit `925095cd`). One file, +16 lines.
6. **Docs** — the patched skill *is* the doc. No other doc references §UIE by structural number, so no sibling doc updates needed. Files in diff: `src/packages/cnos.core/skills/agent/cap/SKILL.md` (implementation, mapped in §ACs); `.cdd/unreleased/362/self-coherence.md` (this artifact). Both enumerated.
7. **Self-coherence** — this file. §Gap → §Skills → §ACs → §Self-check → §Debt → §CDD-Trace → §Review-readiness, each committed incrementally per `alpha/SKILL.md` §2.5.
