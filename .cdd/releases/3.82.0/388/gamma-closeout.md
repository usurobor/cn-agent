<!-- sections: [Cycle, Decision, Merge Plan, Post-Merge Cross-References, Debt-Forwarded, Skills-That-Helped, Outputs, Lifecycle Status] -->
<!-- completed: [Cycle, Decision, Merge Plan, Post-Merge Cross-References, Debt-Forwarded, Skills-That-Helped, Outputs, Lifecycle Status] -->

# γ Closeout — Cycle #388

**Cycle:** [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5 (#366): split generic `schemas/cdd/` from CDS/CDR domain evidence
**Date:** 2026-05-21
**Branch:** `cycle/388`
**Pre-merge cycle SHA:** `7294b473d849740f9fbae252d267b3396fe27497`
**Base SHA:** `7cdcd8d4e4d17340734d0c9df3786d499655c0b0`
**γ identity:** gamma / gamma@cdd.cnos (γ+α+β-collapsed on δ; per breadth-2026-05-12 wave manifest precedent)

---

## §Decision

γ accepts β's APPROVE verdict (R1, no fix-round). Cycle/388 is ready for merge to `origin/main`. Phase 2.5 of cnos#366 is shipped; Phase 3 (V implementation) is unblocked. cnos#376 Sub 1 inherits the architectural-choice decision (option (a) split, recorded in `schemas/cdd/README.md §"Architectural choice"`).

## §Merge Plan

```bash
git checkout main
git pull origin main
git merge --no-ff cycle/388 -m "Merge cycle/388: Phase 2.5 schema generic/domain separation (#388)

Closes #388

[merge body...]"
git push origin main
```

The merge commit message includes `Closes #388` to auto-close the issue. After merge, γ adds the additional artifacts before the close-out commits land in releases/ (per CDD.md §close-out — this cycle's evidence is currently in `.cdd/unreleased/388/`; release versioning is a separate concern handled by the release skill).

## §Post-Merge Cross-References

After the merge lands on `origin/main`, γ posts the following:

### Comment on cnos#366 — Phase 2.5 shipped, Phase 3 unblocked

**Content (template):**

> Phase 2.5 shipped via #388 (merged at SHA [merge-sha]). `schemas/cdd/`, `schemas/cds/`, and `schemas/cdr/` are now distinct CUE packages; the architectural-choice decision (option (a) split) is recorded in `schemas/cdd/README.md §"Architectural choice"`.
>
> The generic kernel schema (`schemas/cdd/`) carries no CDS-specific field by name; `evidence_refs` is a typed open primitive (`[string]: #EvidenceRefValue`) whose required keys are decided by the domain protocol overlays (`schemas/cds/`, `schemas/cdr/`). Receipt frontmatter carries `protocol_id` as a required structural field; Phase 3's V dispatches validation on it.
>
> Phase 3 (V implementation / `cn-cdd-verify` rewrite) is now unblocked. V's input contract is the per-protocol schema package; the dispatch convention is documented in `schemas/cdd/README.md §"protocol_id dispatch convention"`.
>
> Known debt named in cycle/388 close-out:
> - ROLES §4a.3 CDS sketch (`artifact_refs`/`test_refs`/`ci_refs`/`diff_ref`/`debt_refs`) vs `schemas/cds/` realized shape (closure-record names: `self_coherence`/`beta_review`/three closeouts/`diff`/`ci`) divergence — a future cycle may align if convergence is desired.
> - CDR schema is a skeleton; cnos#376 Sub 1/Sub 2 may refine.
> - Essay v0.1.0 open-question #1 resolution not folded back into `CCNF-AND-TYPED-TRUST.md` open-questions section — deferred to a follow-up doc cycle.

### Comment on cnos#376 — Architectural-choice decision inherited

**Content (template):**

> The architectural-choice decision for #376 AC7 (Sub 1) is inherited from cnos#388 (Phase 2.5 of #366, merged at SHA [merge-sha]).
>
> **Option (a) split** is the recorded reading: generic surface = common kernel; per-domain surfaces = protocol overlays in distinct packages/directories. Same decision applies to skills/role-overlays (this issue's scope) as to schemas (cnos#388's scope).
>
> Rationale 1-5 is documented in `schemas/cdd/README.md §"Architectural choice"` and is citable as-is by Sub 1's implementation. The five points:
> 1. CUE's named-field type system (analogous: explicit role-overlay file structure) is the strong-typed primitive.
> 2. Clarity at the protocol boundary (per-role-overlay files name required roles directly).
> 3. Mechanical generic-vs-domain boundary (different package paths).
> 4. Future c-d-X generalizes mechanically (add a directory).
> 5. Unifies with cnos#388's schema split (decision once, applied twice).
>
> Sub 1 does not need to re-decide the architectural-choice — the decision is recorded once and inherited.

## §Debt-Forwarded

The following debt is forwarded out of this cycle for tracking:

1. **§4a.3 vs `schemas/cds/` divergence** — documented in `schemas/cds/README.md`; a future cycle decides alignment direction.
2. **CDR schema refinement** — cnos#376 Sub 1/Sub 2 owns.
3. **Essay v0.1.0 open-question #1 resolution fold-back** — deferred to a follow-on doc cycle.
4. **γ skill `protocol_id` automation** — Phase 3 or γ-skill follow-up.

None of these blocks Phase 3; all four are explicit follow-up surfaces.

## §Skills-That-Helped

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — lifecycle and close-out conventions (the per-cycle artifact set: gamma-scaffold, design-notes, self-coherence, beta-review, alpha/beta/gamma-closeout).
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (implicit, via gamma-scaffold and gamma-closeout conventions seen in `.cdd/releases/3.81.0/385/`).
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — L7 output format informed the architectural-choice decision-record section.

## §Outputs

**Cycle artifacts** (all in `.cdd/unreleased/388/`):
- `gamma-scaffold.md` (γ pre-dispatch)
- `design-notes.md` (L7 design record)
- `self-coherence.md` (α + per-AC oracle evidence)
- `beta-review.md` (R1 APPROVE)
- `alpha-closeout.md` (α α-finished signal)
- `beta-closeout.md` (β β-finished signal)
- `gamma-closeout.md` (this file)

**Substrate changes:**
- `cue.mod/module.cue` (new repo-level CUE module)
- `schemas/cdd/` — refactored to generic-only; README rewritten with decision record
- `schemas/cds/` — new package (CDS protocol overlay)
- `schemas/cdr/` — new package (CDR protocol overlay)
- Fixture corpus re-categorised

## §Lifecycle Status

| Phase | Cycle | Status | Inherits from |
|---|---|---|---|
| Phase 1 — RECEIPT-VALIDATION.md | #367 (closed) | shipped | doctrine |
| Phase 1.5 — COHERENCE-CELL-NORMAL-FORM.md | #370 (closed) | shipped | doctrine |
| Phase 2 — initial CUE schemas | #369 (closed) | shipped | Phase 1 + 1.5 |
| **Phase 2.5 — schema generic/domain split** | **#388 (this cycle)** | **shipping** | **Phase 2; CCNF-AND-TYPED-TRUST.md v0.1.0** |
| Phase 3 — V implementation / cn-cdd-verify rewrite | next dispatch | **unblocked** | Phase 2.5 |
| Phase 4 — δ split | deferred | blocked on Phase 3 | Phase 2.5 |
| Phase 5 — γ shrink | deferred | blocked on Phase 4 | — |
| Phase 6 — ε relocation | deferred | blocked on Phase 5 | — |
| Phase 7 — CDD.md rewrite | deferred | blocked on Phase 3 + 4 + 5 + 6 | — |

cnos#376 Sub 1 (CDR role-overlay architectural-choice) inherits this cycle's decision and is unblocked.

γ proceeds with merge.
