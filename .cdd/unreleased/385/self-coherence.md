<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap] -->

# Self-Coherence — Cycle #385

## §Gap

**Issue:** [#385](https://github.com/usurobor/cnos/issues/385) — Streamline activation soul: collapse 6 CA skills → 2 (cap + clp); deprecate cnos.core/AGENTS.md and hub pre-skill scaffolding

**Version / mode:** 3.15.0 · design-and-build. Issue body is the authoritative design (CLP r2 converged). No separate DESIGN.md or PLAN.md. α executes from the 9 ACs.

**Gap:** `agent/activate/SKILL.md §2.1 step 2` loads 6 CA skills. Four of the six (`mca`, `mci`, `coherent`, `agent-ops`) are partially or fully covered by `cap`; `agent-ops` describes the deprecated OCaml-era daemon runtime. Loading 6 skills on every activation adds context overhead for rules already present in cap and propagates a deprecated surface. The pre-skill scaffolding in `cnos.core/AGENTS.md` references `SOUL.md`/`USER.md` (pre-cycle/379 naming), the OCaml daemon runtime, and an "OPERATIONS mindset" not loaded by activate skill.

**Goal:** Lean soul — KERNEL + cap + clp (3 surfaces). Absorb operational rules from mca/mci/coherent into cap. Standalone files kept for on-demand citation. Renderer and kata fixture updated accordingly. AGENTS.md deleted.
