<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills] -->

# Self-coherence — Cycle #369

## §Gap

**Issue:** [#369 — Phase 2: Add draft CDD contract, receipt, and boundary decision schemas](https://github.com/usurobor/cnos/issues/369).

**Parent roadmap:** #366 (coherence-cell executability). This cycle is **Phase 2**.

**Predecessors:**

- #367 (Phase 1) — froze the parent-facing validator interface in `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`. Named `ValidationVerdict` vs `BoundaryDecision`, the override discipline, the δ-authoritative firing point, and the receipt-derivation rule. Schema syntax deferred to this cycle.
- #364 (closed) — `COHERENCE-CELL.md` doctrine. Frozen by both #367 and this cycle.

**Mode:** docs + schema. No code. No validator behavior change. `cn-cdd-verify` untouched (AC8).

**Version:** schemas pinned at v1 (`cnos.cdd.contract.v1`, `cnos.cdd.receipt.v1`, `cnos.cdd.boundary_decision.v1`). v1 is the **schema-artifact** version, not a CDD protocol version (AC2).

**Gap closed by this cycle:** #367 named the receptor; the types were prose, not declarative. Without `.cue`, Phase 3's validator has to invent shape mid-cycle. This cycle pins the shape so Phase 3 becomes a structural-reader implementation cycle.

**In-scope surfaces (per issue §Scope):**

- `schemas/cdd/contract.cue`
- `schemas/cdd/receipt.cue`
- `schemas/cdd/boundary_decision.cue`
- `schemas/cdd/README.md` (single-source prose; CUE files reference via header comment)
- `schemas/cdd/fixtures/valid-receipt.yaml`
- `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml`
- `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml`
- `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml`
- `.cdd/unreleased/369/*.md` (cycle evidence)

**Out of scope (per issue §Non-goals):** any change to `cn-cdd-verify/`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `delta/SKILL.md` (must not exist yet), `ROLES.md`, `CDD.md`, `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md` (sole permitted touch is an optional trailing-pointer; this cycle takes the safe path and omits it).

## §Skills

**Tier 1 (CDD lifecycle):**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm (consumed; not edited).
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface; load order, pre-review gate (14 rows), incremental self-coherence discipline, SHA-convention rules, polyglot re-audit rule.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority. Consulted to interpret AC boundaries.

**Tier 2:** none load-bearing for a docs+schema cycle with no shell/test toolchain authoring. (The repo's `cnos.eng/` Tier 2 bundles describe language-specific authoring; CUE-shape authoring is governed by the precedent of `schemas/skill.cue` per the issue's "Skills to load" section — "No new schema skill — the existing `schemas/skill.cue` is the precedent and form authority for CUE-shape patterns in this repo.")

**Tier 3 (issue-named):**

- `src/packages/cnos.core/skills/design/SKILL.md` — applied as authoring constraints (§3.2 single source of truth; §3.5 interfaces belong to consumers; §3.7 separate runtime surfaces; §3.10 prefer package/install cohesion).
- `src/packages/cnos.core/skills/write/SKILL.md` — applied to README prose: decisive over exhaustive; surface-containment statements internal to the document.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — used to interpret AC boundaries; no issue edits this cycle.

**Form authorities consulted (not edited):**

- `schemas/skill.cue` + `schemas/README.md` + `schemas/fixtures/skill-frontmatter/` — schema-layer precedent: CUE owns shape; prose lives in `schemas/{subsystem}/README.md`; fixture corpus lives under `schemas/{subsystem}/fixtures/`.
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — semantic owner of the interface this cycle types. Schemas reference it via header comment per active design constraint.
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — doctrine. Recursive scope-lift framing inherited from §Recursion Equation.
