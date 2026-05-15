<!-- sections: [Preamble, Q1, Q2, Q3, Q4, Q5, Validation Interface, Non-goals, Closure] -->
<!-- completed: [Preamble] -->

# Receipt Validation — Parent-facing Validator Surface

**Status:** Draft design. Phase 1 of #366 (coherence-cell executability roadmap).
**Owns:** The parent-facing validation interface for the coherence cell — input shape, output shape, invocation contract, and resolved positions on the five Open Questions seeded by [#364](https://github.com/usurobor/cnos/issues/364).
**Does-not-own:** Doctrine (`COHERENCE-CELL.md`), the canonical executable algorithm (`CDD.md`), CUE schema syntax (Phase 2), validator implementation (Phase 3), role-skill edits (Phases 4–6), `CDD.md` rewrite (Phase 7).

---

## Preamble

This document explains how the validation predicate `V` is positioned at the parent-facing boundary of a coherence cell — when it fires, what it consumes, what it emits, and how its output composes with δ's boundary decision. It is **draft design**, not doctrine. The doctrine surface is `COHERENCE-CELL.md`, which names the recursive cell model and declares

```text
V : Contract × Receipt × EvidenceGraph → ValidationVerdict
```

as load-bearing without specifying the parent-facing interface. This document specifies that interface at doctrine level.

### Status of this document

| Property | Value |
|---|---|
| Document class | Draft design — sibling to `COHERENCE-CELL.md` |
| Doctrine surface | `COHERENCE-CELL.md` (frozen by this cycle's contract; not edited here) |
| Executable algorithm surface | `CDD.md` (frozen by this cycle's contract; not edited here) |
| Predecessor cycle | [#364](https://github.com/usurobor/cnos/issues/364) — landed `COHERENCE-CELL.md` at merge `32b126e4` |
| Parent roadmap | [#366](https://github.com/usurobor/cnos/issues/366) — coherence-cell executability roadmap; this issue is Phase 1 |
| Binding status | **This document becomes binding once Phase 3 implements `V` as a working predicate.** Until Phase 3 lands, this document is the design surface Phase 2 (schemas) and Phase 3 (validator) consume as input contract. |
| Schema syntax | Out of scope — Phase 2 |
| Validator implementation | Out of scope — Phase 3 |
| Role-skill rewrites | Out of scope — Phases 4–6 |
| `CDD.md` rewrite | Out of scope — Phase 7 |

This document is **draft design** because the predicate it specifies does not yet have an implementation. The interface freeze is the load-bearing act: Phases 2–7 inherit the positions named here so each downstream cycle inherits a stable input contract rather than re-deriving it. When `cn-cdd-verify` implements `V` in Phase 3, this document moves from "design that constrains future work" to "doctrine that current work satisfies."

### What this document does

1. Resolves the five Open Questions seeded by `COHERENCE-CELL.md` §Open Questions: when `V` fires (§Q1), whether `V` is a capability or a command (§Q2), where ε relocates (§Q3), how an override is receipted (§Q4), and what role the existing per-cycle markdown artifacts play in the evidence graph (§Q5).
2. Freezes the parent-facing validation interface at doctrine level: input contract, output contract, invocation contract — distinguishing `ValidationVerdict` (emitted by `V`) from `BoundaryDecision` (emitted by δ) (§Validation Interface).
3. Restates the non-goals so the surface-containment contract from the issue body is internal to the document (§Non-goals).

### What this document does not do

It does not pin schema syntax — no `.cue`. It does not modify `cn-cdd-verify` code. It does not split `operator/SKILL.md` or create a new `delta/SKILL.md`. It does not shrink `gamma/SKILL.md`. It does not move ε out of its current home; it names the target for Phase 6. It does not rewrite `CDD.md`. It does not edit `COHERENCE-CELL.md` doctrine.

### Reading order

The five Open Questions are not independent. Q1 fixes the firing point at the δ boundary, which determines who invokes `V` and where the receipt is read — that is what Q2's capability-vs-command choice is in service of. Q4's override discipline only makes sense once Q1's δ-authority is named — the override is a degraded δ boundary action against a non-PASS `ValidationVerdict`. Q5's evidence-graph derivation rule is the input side of the interface Q1 invokes. Q3 is structurally separate: it names ε's home, which affects how `cdd-iteration.md` sits relative to receipts (evidence-graph input vs receipt-class artifact), and Phase 6 ships the move.

Read in order: Q1 → Q5 → Q2 → Q4 → Q3, then §Validation Interface, then §Non-goals.

The document is written in narrative form: each `## Q{n}` section opens with the chosen position, then justifies it against the doctrine, then names the structural consequence. This is design prose, not a five-row answer table. The position is the load-bearing claim; the rationale is what makes the position defensible against future drift.

### Inherited doctrine fragments (verbatim from `COHERENCE-CELL.md`)

The following two fragments from the predecessor are quoted here as the load-bearing anchors this design extends. They are quoted, not re-derived; the doctrine surface remains authoritative.

> A closed cell is not trusted because a higher role approved it. A closed cell is trusted because its receipt validates against its contract.
> — `COHERENCE-CELL.md` §Thesis

> Trust is established by the relation `contract + evidence + valid receipt`. That triple is what holds. Trust is not issued by γ closing the cell, by δ tagging the release, or by ε approving the protocol. Those are role actions. Validity is what carries trust across the boundary.
> — `COHERENCE-CELL.md` §Trust Boundary

The first fragment is what makes `V` load-bearing rather than ceremonial. The second fixes the trust grammar: validity is a property of the (contract, evidence, receipt) triple, not a property a role confers. Every position named below is consistent with both fragments.

---
