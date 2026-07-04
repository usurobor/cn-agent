<!--
section-manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
  completed: [Gap]
-->

# Self-coherence — cnos#584

## §Gap

**Issue:** [usurobor/cnos#584](https://github.com/usurobor/cnos/issues/584) — "arch(cds/cdd): codify mechanical cell runtime vs cognitive skills (mechanism/cognition boundary)". Sub 1 of 5 under parent #583 (master wave — mechanical dispatch-cell architecture). Lands first.

**Mode:** design-and-build, doctrine-only — no runtime change.

**Governing rule to codify (operator-stated):** cells are mechanical and defer cognition to skills; skills don't control anything. The cell is the execution engine — it owns control flow and all side effects, and calls a skill when it needs a judgment. A skill returns cognition; it never opens a PR, writes a label, or drives the lifecycle.

**Branch:** `cycle/584`, branched from `origin/main@9309de97d7e6d90637012839163e8d0511b56ca6` (per γ's scaffold; the wake-invoked-δ input's pinned SHA `9277c7e3...` had already advanced by one unrelated commit — board-map regen — by branch time; γ verified and branched from the current HEAD per its branch pre-flight rule, not the stale pinned value).

**Scope guardrail (restated):** in scope — codify the governing rule in cell-framework doctrine; classify the 9 named dispatch lifecycle steps mechanical vs. cognitive; audit 5 named files (`cdd/gamma`, `cdd/alpha`, `cdd/beta`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`) for control-implying prose and correct it. Out of scope — any runtime/workflow/transition-table change (Subs 2–4), implementing checkpoint/PR-open/recovery, `cell_kind` enforcement, new status labels/taxonomy, Demo 0.
