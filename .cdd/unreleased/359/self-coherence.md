<!-- sections: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness] -->
<!-- completed: [gap, skills, acs] -->

# Self-Coherence — Issue #359

## Gap

`operator/SKILL.md §5.2` titled "Single-session δ-as-γ via Agent tool" did not state which role-pair the collapse covers. The intended reading is "δ acts as γ; α and β remain isolated sub-agents." The observed misread (tsc #49 wave-1) is "all roles fused into one sub-agent," which violates `CDD.md §1.4` Triadic rule. §5.2 needed an explicit scope statement and an explicit violation shape.

Mode: docs-only. Branch: `cycle/359`. Base: `23e28e45`.

## Skills

- Tier 1 / role: `CDD.md`, `alpha/SKILL.md`. Lifecycle sub-skills not loaded — single-paragraph skill patch with three ACs needed no design or plan artifact.
- Tier 2: none — change is not engineering code.
- Tier 3: `src/packages/cnos.core/skills/write/SKILL.md` (governing question, brevity-is-earned, state-what-it-is, front-load-the-point).

## ACs

Implementation SHA: `22e9e7eb` (`α-359: clarify §5.2 collapses δ↔γ only — γ↔α↔β stays separate`). Diff: `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md +4 / -0`.

**AC1 — `operator/SKILL.md §5.2` contains explicit statement that §5.2 collapses δ↔γ only and γ↔α↔β remain structurally separate per §1.4.** Evidence: the new paragraph in `operator/SKILL.md §5.2` reads "§5.2 collapses **δ↔γ only**. γ↔α↔β remain structurally separate per `CDD.md §1.4` Triadic rule: γ scaffolds and coordinates in the parent session, α implements in its own sub-agent, β reviews and merges in its own sub-agent." `git show 22e9e7eb -- src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` confirms the line lands at the head of the three-consequences list. Met.

**AC2 — names what violation looks like (single subagent doing γ+α+β work).** Evidence: the second new paragraph reads "A single sub-agent that performs γ-selection plus α-implementation plus β-review is not §5.2 — it is a §1.4 violation. §5.2 requires three execution contexts: the parent session (γ, also δ), a separate α sub-agent, and a separate β sub-agent. Lumping γ+α+β into one sub-agent breaks role-isolation (α gains access to β's reasoning and vice versa) and is rejected." Met.

**AC3 — `operator/SKILL.md §5.2` consistent with the patch.** Evidence: re-read §5.2 mechanism paragraph ("α and β are dispatched as sub-agents") and consequence (1) ("δ=γ collapse"). The new "Scope of the collapse" block does not contradict these — it sharpens the mechanism paragraph (which already implied α/β separation) and explains the *scope* that consequence (1) is internal to. Downstream references to §5.2 in `activation/SKILL.md §8` (dispatch declaration), `release/SKILL.md §3.8` (configuration-floor clause), and `operator/SKILL.md §5.3` (escalation criteria) all describe "γ/δ separation absent" — none asserts α/β collapse, so the patch leaves them coherent. `rg "§5\.2|δ-as-γ|δ=γ|δ↔γ" src/packages/cnos.cdd/skills/cdd/` enumerated and inspected: 4 files match, all consistent. Met.
