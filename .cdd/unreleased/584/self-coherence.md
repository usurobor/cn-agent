<!--
section-manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
  completed: [Gap, Skills]
-->

# Self-coherence — cnos#584

## §Gap

**Issue:** [usurobor/cnos#584](https://github.com/usurobor/cnos/issues/584) — "arch(cds/cdd): codify mechanical cell runtime vs cognitive skills (mechanism/cognition boundary)". Sub 1 of 5 under parent #583 (master wave — mechanical dispatch-cell architecture). Lands first.

**Mode:** design-and-build, doctrine-only — no runtime change.

**Governing rule to codify (operator-stated):** cells are mechanical and defer cognition to skills; skills don't control anything. The cell is the execution engine — it owns control flow and all side effects, and calls a skill when it needs a judgment. A skill returns cognition; it never opens a PR, writes a label, or drives the lifecycle.

**Branch:** `cycle/584`, branched from `origin/main@9309de97d7e6d90637012839163e8d0511b56ca6` (per γ's scaffold; the wake-invoked-δ input's pinned SHA `9277c7e3...` had already advanced by one unrelated commit — board-map regen — by branch time; γ verified and branched from the current HEAD per its branch pre-flight rule, not the stale pinned value).

**Scope guardrail (restated):** in scope — codify the governing rule in cell-framework doctrine; classify the 9 named dispatch lifecycle steps mechanical vs. cognitive; audit 5 named files (`cdd/gamma`, `cdd/alpha`, `cdd/beta`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`) for control-implying prose and correct it. Out of scope — any runtime/workflow/transition-table change (Subs 2–4), implementing checkpoint/PR-open/recovery, `cell_kind` enforcement, new status labels/taxonomy, Demo 0.

## §Skills

**Tier 1 (loaded, canonical):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle + kernel doctrine; the chosen AC1/AC2 doctrine home.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (this file's governing load order).

**Tier 1 sub-skill:**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — cell-shaped issue contract; used to read the issue pack correctly (mode = design-and-build, source-of-truth table, non-goals discipline).

**Tier 3 (issue-specific, named by γ's dispatch prompt):**
- `src/packages/cnos.core/skills/write/SKILL.md` — prose-quality discipline (one governing question per doc, front-loaded point, stable-fact-once-then-point, contrastive examples). Applied throughout the new CDD.md section and this artifact: the section leads with the rule, states the worked example before generalizing, and the lifecycle table states each row's fact once rather than repeating it in prose.

**Read for context, not loaded as an authoring skill (per scope guardrails):**
- `docs/papers/CELL-OF-CELLS.md` — read in full as the other candidate doctrine home; not chosen (see §ACs → AC1 for the placement justification).
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`, `.../beta/SKILL.md` (already Tier-1-adjacent for α context), `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — the 5 AC3 audit files, read in full for the audit.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 — read for citation (existing positive instance of the mechanism/cognition split); not an audit target, not restructured.

No Tier 2 `eng/*` skills were loaded — this cycle produces no source code in any language (per the pinned implementation contract's Language axis: N/A).
