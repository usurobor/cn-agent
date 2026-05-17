<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap] -->

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
