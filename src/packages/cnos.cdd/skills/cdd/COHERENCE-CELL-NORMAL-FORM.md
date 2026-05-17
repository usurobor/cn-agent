<!-- sections: [Preamble, Kernel, Cell Outcomes, Recursion Modes, Scope-Lift, Two-Layer Separation, Non-goals, Closure] -->
<!-- completed: [Preamble] -->

# Coherence-Cell Normal Form

**Status:** Draft doctrine. Phase 1.5 of [#366](https://github.com/usurobor/cnos/issues/366) (coherence-cell executability roadmap).
**Owns:** The substrate-independent kernel algorithm the coherence-cell doctrine implies — a five-step recursion at scope `n`, four closed-cell outcomes, two recursion modes, three scope-lift projections, and the two-layer separation between kernel and realization.
**Does-not-own:** Predecessor doctrine (`COHERENCE-CELL.md`), parent-facing validator design (`RECEIPT-VALIDATION.md`), CUE schemas (Phase 2, [#369](https://github.com/usurobor/cnos/issues/369)), validator implementation (Phase 3, deferred), role-skill rewrites (Phases 4–6, deferred), `CDD.md` rewrite (Phase 7, deferred).

---

## Preamble

This document is **draft doctrine**. It is a companion at the kernel layer to `COHERENCE-CELL.md` (the receipt-rule doctrine for the coherence cell) and to `RECEIPT-VALIDATION.md` (the parent-facing typed-interface design surface for the validator `V`). It names the recursive algorithm the predecessor doctrine implies but does not state.

### What this document is

**Coherence-Cell Normal Form (CCNF)** is the substrate-independent recursion algorithm of the coherence cell. The cell at scope `n` is the unit of work; the kernel is the closed loop that produces a closed cell at `n` and projects it as α-matter at scope `n+1`. The kernel is stated here as five composed steps, four outcomes, two recursion modes, and three scope-lift projections. Nothing in the kernel names a substrate; every kernel statement is reusable across instantiations of the role-scope ladder pattern.

### Companion, not replacement

`COHERENCE-CELL.md` remains the **predecessor doctrine**. The receipt rule, the four-way structural separation (role / runtime substrate / validation / boundary effection), and the role-as-cell-function framing all live there unchanged. This document does not replace, supersede, or silently override `COHERENCE-CELL.md`; it names the algorithm the doctrine implies and pins it as a citable surface.

`RECEIPT-VALIDATION.md` remains the **parent-facing receptor design**. The `V` invocation contract, the verdict-versus-decision distinction, and the δ-authoritative firing point all live there. This document positions `V` as step 4 of the kernel and δ as step 5, but it does not re-derive the interface — the receptor design owns the typed shape; the kernel owns the algorithmic position.

The split is structural: the predecessor doctrine names the *organism*; the receptor design names the *receptor*; this document names the *algorithm*. The three surfaces together describe what a coherence cell is, how its parent boundary validates it, and how its recursion closes. None of the three replaces the others.

### Why this document exists

Phase 3 (validator implementation) consumes a recursion no doctrine names — `V`'s contract is derivable from the cell loop but not citable. Phase 4 (δ split) decides what δ is with no kernel statement of δ's signature. Phase 6 (ε relocation) moves ε without a kernel statement of ε's signature as receipt-stream observation projecting as coordination at the parent scope. Phase 7 (`CDD.md` rewrite) would have to derive and state the kernel mid-rewrite — exactly the failure mode the two-layer split exists to prevent.

With this kernel landed: Phases 3, 4, 6, and 7 each cite the relevant kernel section of this document as their input contract. The kernel becomes the spine; the operational realization (Phase 7) becomes a sympathetic expansion rather than an act of derivation.

### Authoring discipline

Two disciplines govern this document.

**Substrate-independence in the kernel.** Sections `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, and `## Scope-Lift` are the kernel layer. They name only roles (α, β, γ, δ, ε), the validator predicate (`V`), the artifacts the kernel reasons over (`contract`, `matter`, `review`, `receipt`, `verdict`, `decision`, `evidence`), the scope index (`n`, `n+1`), and the verdicts and decisions the algorithm composes. They do not name any particular tooling, platform, or invocation surface — those names appear only in `## Two-Layer Separation`, where they are explicitly labelled as realization peers.

**Decisive over exhaustive.** Each section commits to one position with the smallest set of statements that lets downstream phases cite it. This is a kernel doctrine; it is not a re-derivation of the operational realization. The target length is 200–400 lines; if a section grows past what the kernel needs, it is probably re-deriving `CDD.md` mid-cycle and should shrink.

### Reading order

Read the kernel sections in declared order — they compose. `## Kernel` names the five-step recursion and the evidence-binding rule. `## Cell Outcomes` pins the four ways a cell terminates at scope `n` (the algorithm's possible exits). `## Recursion Modes` distinguishes the two ways the recursion advances — within-scope repair-dispatch and cross-scope projection. `## Scope-Lift` names the three projections that carry from scope `n` to scope `n+1`. `## Two-Layer Separation` then declares this document as the kernel layer and names the realization peers that instantiate it on the current substrate.

The non-goals and closure sections at the end restate the surface-containment contract this cycle commits to and the criteria under which the kernel is complete.
