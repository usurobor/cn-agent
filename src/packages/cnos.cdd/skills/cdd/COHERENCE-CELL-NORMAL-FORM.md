<!-- sections: [Preamble, Kernel, Cell Outcomes, Recursion Modes, Scope-Lift, Two-Layer Separation, Non-goals, Closure] -->
<!-- completed: [Preamble, Kernel] -->

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

---

## Kernel

The coherence cell at scope `n` is a five-step closed loop. The kernel reads as a recursion: the five steps compose into a closed cell, the closed cell terminates in one of four outcomes (§Cell Outcomes), and on the cross-scope outcomes the closed cell projects as α-matter at scope `n+1` (§Recursion Modes, §Scope-Lift).

### Five-step closure at scope `n`

At scope `n`, a `contractₙ` is given. The kernel composes five steps over that contract:

```text
1.  matterₙ      := αₙ.produce(contractₙ)
2.  reviewₙ      := βₙ.review(contractₙ, matterₙ)
3.  receiptₙ     := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)
4.  verdictₙ     := V(contractₙ, receiptₙ)
5.  decisionₙ    := δₙ.decide(receiptₙ, verdictₙ)
```

Each step is a typed function with a fixed signature. The signatures load-bear on the rest of the algorithm; the rest of this section names what each signature means and what it explicitly excludes.

### Step-by-step

**Step 1 — α produces matter against the contract.**

```text
matterₙ := αₙ.produce(contractₙ)
```

α reads the contract — the issue body, its acceptance criteria, its scope, its non-goals, its active design constraints — and produces matter against that contract. The matter is what the cell builds: source, tests, doctrine, design surfaces, schemas, prose. α is accountable to the contract, and only to the contract. α does not consult the receipt (the receipt does not yet exist), and α does not consult the verdict or decision (those come later).

**Step 2 — β reviews the (contract, matter) pair. β consumes matter only.**

```text
reviewₙ := βₙ.review(contractₙ, matterₙ)
```

β's signature explicitly excludes evidence as a separate input. β discriminates the matter against the contract — that is β's biological function as the cell's immune discrimination — but β does not consume an evidence graph alongside the matter. Evidence accumulates as α and β work, and γ binds it into the receipt at close-out (step 3); β never reads a typed evidence object. Concretely: β reads the matter and the contract; if some claim in the matter is checkable, β checks it; if some predicate the contract names must be verified, β verifies it. β's review record names which checks β performed, not a typed reference graph β walked.

The reason β's signature excludes evidence is structural. If β consumed evidence as a typed input, β would be performing receipt-time validation — which is `V`'s job at step 4. Collapsing β's discrimination and `V`'s validation into one step would erase the discriminate-then-validate sequence the kernel relies on.

**Step 3 — γ closes the cell. γ binds evidence into the receipt.**

```text
receiptₙ := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)
```

γ closes the cell by fusing the contract, the matter, the review, and the cell's accumulated evidence into a typed receipt. The receipt is the parent-facing artifact of the closed cell — what crosses the boundary to scope `n+1`. γ's act of closure has two parts: it records the cell's narrative (production record, review record, closure triage) *and* it binds the evidence into the receipt as typed references the downstream consumer can dereference.

The **evidence-binding rule** is load-bearing and stated once here for the kernel:

> **Evidence accumulates during α and β work. γ binds it into the receipt at close-out as typed references. V dereferences the references to validate. β never consumes evidence directly. δ never re-reads evidence.**

The rule has four enforcement points:
- α's signature excludes the receipt (the receipt does not yet exist at step 1).
- β's signature excludes evidence (step 2).
- γ's signature includes evidence as the input γ binds (step 3) — γ is where the evidence references enter the typed receipt.
- δ's signature excludes evidence (step 5) — δ trusts the verdict, not the underlying evidence graph.

The rule is what keeps the receipt as a typed handoff rather than a narrative pointer. A receipt with no evidence references is a receipt the parent scope cannot validate. A receipt where β consumed evidence directly and γ then re-bound it would carry two evidence views at once — α's view and β's view — collapsing the receipt's role as the single typed surface that crosses the boundary.

**Step 4 — V validates the receipt against the contract. V dereferences evidence from the receipt.**

```text
verdictₙ := V(contractₙ, receiptₙ)
```

`V` is a typed predicate. It reads the contract and the receipt, dereferences the evidence references the receipt carries, walks the referenced evidence graph, and emits a verdict — `PASS`, `FAIL`, or (optionally) `WARN`. The verdict is structured: it carries a headline value, a list of failed predicates with refs back to the contract, and a list of advisory warnings.

`V`'s signature shows two inputs, not three. The receipt carries the evidence references; `V` does not take an unbound `evidenceₙ` argument. This is the canonical form of evidence-binding: the receipt is the single object that carries everything `V` needs to dereference, and the cycle's accumulated evidence is reachable only through the receipt's typed references. A `V` that read evidence directly would be a `V` that could be invoked against an inconsistent (receipt, evidence) pair; binding evidence to the receipt prevents that class of inconsistency by construction.

`V` does not consult α, β, or γ as actors. `V` reads artifacts. It does not have discretionary authority — it is a predicate, not a role. The predecessor doctrine pins this position as "`V`, not ζ"; the kernel inherits it without re-derivation.

**Step 5 — δ decides at the boundary. δ decides on receipt and verdict only.**

```text
decisionₙ := δₙ.decide(receiptₙ, verdictₙ)
```

δ's signature explicitly excludes evidence. δ reads the receipt, reads the verdict, and records a decision — one of `{accept, release, override, reject, repair_dispatch}`. δ does not re-read the evidence graph; δ trusts the verdict `V` emitted. If δ doubts the verdict, δ has two paths: invoke `V` again (the verdict is reproducible from the same inputs) or repair-dispatch the cell to recompute its matter. δ does not have the path "read the evidence directly and reach an independent verdict" — that would collapse step 4 and step 5 into one actor and let δ confer validity (which the predecessor doctrine explicitly rejects).

The decision is what closes the cell. Until δ records a decision, the cell is in the closure-emitted-but-not-decided intermediate state: γ has done γ's work, `V` has done `V`'s work, but the boundary has not yet observed the result. The closed cell exists as a typed object only once δ has recorded `decisionₙ`.

### Composition

The five steps compose into a closed cell at scope `n`:

```text
closed_cellₙ := { contractₙ, matterₙ, reviewₙ, receiptₙ, verdictₙ, decisionₙ }
```

The closed cell is the kernel's terminal object at scope `n`. Its outcome — which of the four terminal states it inhabits — is determined by `verdictₙ × decisionₙ` (§Cell Outcomes). Its recursion mode — whether the cell stays open at scope `n` under repair-dispatch or projects to scope `n+1` under accept/release/override — is determined by `decisionₙ` (§Recursion Modes). Its projection — which closed-cell components carry to scope `n+1` and which do not — is governed by the scope-lift (§Scope-Lift).

### What the kernel does not name

The kernel names roles, the validator predicate, artifacts, verdicts, decisions, and scopes. It does not name any particular tooling, platform, dispatch mechanism, schema language, or invocation surface. Those names live in the realization layer (§Two-Layer Separation). The kernel is reusable across any instantiation of the role-scope ladder pattern that emits typed receipts and invokes a parent-facing predicate.
