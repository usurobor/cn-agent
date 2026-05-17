<!-- sections: [Preamble, Kernel, Cell Outcomes, Recursion Modes, Scope-Lift, Two-Layer Separation, Non-goals, Closure] -->
<!-- completed: [Preamble, Kernel, Cell Outcomes, Recursion Modes] -->

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

---

## Cell Outcomes

A closed cell at scope `n` terminates in **exactly one of four outcomes**. The outcome is determined by the pair `(verdictₙ, decisionₙ)` — `V`'s emitted verdict and δ's recorded decision. Each outcome has a precondition the kernel pins explicitly; cells with a `(verdict, decision)` pair outside the four preconditions are not terminal and must be re-decided.

### The four outcomes

```text
accepted := (verdictₙ = PASS)  ∧ (decisionₙ ∈ {accept, release})
degraded := (verdictₙ ≠ PASS)  ∧ (decisionₙ = override)
blocked  := (decisionₙ ∈ {reject, repair_dispatch})        — any verdict
invalid  := (verdictₙ = PASS   ∧ decisionₙ = override)
          ∨ (verdictₙ ≠ PASS   ∧ decisionₙ ∈ {accept, release})
```

Each outcome carries a structural meaning:

- **`accepted`** — the cell is clean. `V` returned `PASS`; δ recorded a non-degraded boundary action. The closed cell is transmissible to scope `n+1` as α-matter without qualification. This is the kernel's PASS-equivalent terminal outcome.
- **`degraded`** — the cell is transmissible under explicit override. `V` returned `FAIL` (or `WARN`); δ chose to proceed anyway by recording `override` with the override block populated. The closed cell is transmissible, but the parent scope's α at scope `n+1` reads degraded matter — the override block is the structural signal that every downstream consumer must detect. Override is a degraded boundary action, not a form of validity.
- **`blocked`** — the cell does not transmit. δ recorded `reject` (the cell's matter cannot cross the boundary as-is) or `repair_dispatch` (the cell stays open at scope `n` and a child cell at the same scope addresses the failure under a repair contract — see §Recursion Modes). The cell at `n` does not produce α-matter at `n+1` in either case. `reject` closes the cell with no parent-scope projection; `repair_dispatch` keeps the cell open until the child cell's accept lets the parent γₙ re-emit `receiptₙ` and the loop re-fires.
- **`invalid`** — the cell is malformed at the boundary. δ's decision is inconsistent with `V`'s verdict: either δ recorded `override` against a `PASS` verdict (overriding nothing — there is nothing to override) or δ recorded `accept`/`release` against a non-`PASS` verdict (accepting matter `V` failed without invoking the override path). Both branches indicate δ has misread the verdict or has bypassed the override discipline.

### `invalid` is non-terminal

**`invalid` cells do not close.** When δ's decision falls into the `invalid` precondition, the kernel does not terminate at scope `n`. δ must re-decide — re-read the verdict, choose a consistent decision (`accept`/`release` for `PASS`; `override` with the override block populated for non-`PASS`; or `reject`/`repair_dispatch` for any verdict) — before the cycle can terminate.

The non-terminal status of `invalid` is what keeps the four outcomes well-formed. If `invalid` were terminal, the kernel would have two contradictory exits ("the verdict is `PASS` and δ overrode it" / "the verdict is not `PASS` and δ accepted it") that the parent scope could not reason over. By making `invalid` non-terminal, the kernel forces δ back to a consistent decision and the closed cell that eventually emerges inhabits exactly one of `{accepted, degraded, blocked}`.

The mechanism is δ-internal. δ holds gate authority over what crosses the boundary; an `invalid` pair is δ recognising that the decision it recorded does not match the verdict it observed and recording a corrected decision. No other role re-fires. `V` is not re-invoked unless δ chooses to (the verdict is reproducible — the same inputs produce the same verdict). γ does not re-close (the receipt is already emitted). The re-decision is δ adjusting its `decisionₙ` until the resulting `(verdictₙ, decisionₙ)` pair satisfies one of the three terminal preconditions.

### Alignment with the receptor design

The four-outcome table aligns with `#369` AC4 — the verdict × action × transmissibility table that the receptor's schema work is typing. The kernel pins the *what* (four outcomes; their preconditions; `invalid` non-terminal); the schema work pins the *type* (the enums and field shapes that make these outcomes machine-checkable). The two surfaces are designed to compose: when the schema work lands, the kernel's preconditions become enforceable by the typed schemas, and the kernel's `invalid` non-terminal status becomes a schema-level constraint on the boundary record.

Where this document and the receptor's schema work disagree on the outcome preconditions, the kernel statement here is the load-bearing claim and the schema work aligns to it — not the reverse. The schema is the realization layer typing the kernel; the kernel is the *what* the realization layer types.

### Composition with the recursion modes

The four outcomes interact with the two recursion modes (§Recursion Modes) in a specific way:

- `accepted` → cross-scope projection: the closed cell becomes α-matter at scope `n+1`.
- `degraded` → cross-scope projection: the closed cell becomes α-matter at scope `n+1` carrying the override block; the parent scope reads degraded matter.
- `blocked` with `decision = reject` → no projection: the cell terminates at scope `n` with no scope-`n+1` consequence beyond the recorded rejection.
- `blocked` with `decision = repair_dispatch` → within-scope recursion: the cell stays open at scope `n`; a child cell at scope `n` runs under a repair contract; on child accept, parent γₙ re-emits `receiptₙ` and the kernel re-fires steps 4–5.
- `invalid` → no projection and no recursion: δ re-decides until the outcome inhabits one of the three terminal cases.

The four outcomes are the kernel's possible exits; the recursion modes are how those exits propagate. The next section names the recursion modes explicitly.

---

## Recursion Modes

The kernel's recursion has **two distinct modes**. They share the same recursion operator — the five-step closure at some scope — but differ in how the scope index advances. Conflating them weakens the algorithm: the same closed-cell record carries different meanings depending on which mode it belongs to, and a downstream consumer that cannot tell them apart cannot reason about the cell's position in the recursion.

### Within-scope mode — repair-dispatch (same scope index)

```text
Precondition:  decisionₙ = repair_dispatch         (outcome: blocked, repair sub-case)
Behaviour:     cell at scope n stays open
               a child cell runs at scope n under a repair contract
               on child accept: parent γₙ re-emits a fresh receiptₙ
               V re-fires at step 4
               δ re-decides at step 5
Scope index:   unchanged — both parent and child operate at scope n
```

When δ records `repair_dispatch`, the cell at scope `n` does not close. Instead, a child cell at scope `n` runs under a **repair contract** — a contract derived from the failure that triggered the repair-dispatch. The child cell goes through the same five-step closure (steps 1–5) at scope `n`, producing its own child `receipt`, child `verdict`, and child `decision`. When the child cell terminates in `accepted`, the parent γₙ re-emits a fresh `receiptₙ` that incorporates the repaired matter; the kernel's steps 4 (`V`) and 5 (δ) re-fire on the fresh receipt at the parent level; and δ records a new `decisionₙ` over the re-emitted receipt.

Three properties of this mode:

1. **Scope is preserved.** The parent cell is at scope `n` throughout; the child cell is also at scope `n`. The recursion does not advance the scope index. The parent scope (`n+1` and above) sees nothing — there is no cross-scope projection until the parent cell eventually terminates in one of `{accepted, degraded, blocked-reject}`.
2. **The receipt is parent γ-owned.** The child cell's receipt is the child's parent-facing artifact (which means: the parent cell's γₙ reads it as input alongside the child's verdict and decision). It is not the parent cell's receipt — only γₙ at the parent level re-emits `receiptₙ`, and only after the child cell has accepted. The kernel's step 3 at the parent level re-fires; the child's step 3 produced the child's own receipt, not the parent's.
3. **The same recursion operator.** The child cell is itself a five-step closed loop at scope `n`. The kernel's recursion operator is not a different operator for within-scope work — it is the same five-step closure, applied at the same scope, against a derived repair contract. The within-scope mode is recursion over a contract, not recursion across scopes.

### Cross-scope mode — accept / degraded (scope index advances)

```text
Precondition:  decisionₙ ∈ {accept, release}     ∧ verdictₙ = PASS        (outcome: accepted)
            ∨  decisionₙ = override               ∧ verdictₙ ≠ PASS        (outcome: degraded)
Behaviour:     cell at scope n closes
               closed_cellₙ projects as α-matter at scope n+1
               under degraded outcome, the override block is part of the matter
Scope index:   advances — parent cell at n+1 operates on the closed_cellₙ projection
```

When δ records `accept`, `release`, or `override` against the matching verdict, the cell at scope `n` closes terminally. The closed cell — its contract, matter, review, receipt, verdict, decision, and (for degraded outcomes) its override block — projects as α-matter at scope `n+1`. The recursion advances: scope-`n+1` α reads the projected closed cell as input and produces matter against the scope-`n+1` contract. The scope-`n` cell is now one of the scope-`n+1` cell's α-inputs.

Three properties of this mode:

1. **Scope advances.** The recursion operator fires again, but at scope `n+1`, with a new contract, new actors (αₙ₊₁, βₙ₊₁, γₙ₊₁, δₙ₊₁, εₙ₊₁), and the closed cell from scope `n` as one of α-input objects. The scope index is the kernel's witness to the recursion having advanced; the closed cell at `n` is the boundary between the two levels.
2. **The receipt is the carrier.** The closed cell projects as α-matter via its receipt. The parent scope's α reads the receipt as a typed object — it does not re-read the cell's internal artifacts. This is what makes the recursion well-formed: the parent scope reasons over typed surfaces, not over narrative. Under `accepted`, the receipt is PASS-equivalent; under `degraded`, the receipt's `boundary.override` block is non-null and the parent scope reads degraded matter.
3. **`reject` is also cross-scope, but with no projection.** When δ records `reject`, the cell closes at scope `n` (it terminates) but does not project — the cell's matter is recorded as rejected at scope `n` with no scope-`n+1` consequence beyond the rejection record itself. Reject is the cross-scope mode's null projection; accepted and degraded are the cross-scope mode's full projection (clean and degraded respectively).

### Same operator, different scope behaviour

The kernel's recursion is a single operator — the five-step closure applied to a contract at some scope. The two modes are how that operator's output relates to the scope index:

- **Within-scope (repair-dispatch):** the operator re-fires at the same scope on a derived contract. The cell stays open at the parent level; the recursion does not advance.
- **Cross-scope (accept/degraded/reject):** the cell closes at the current scope; the closed cell either projects (accepted/degraded) or does not project (reject) to the next scope.

The two modes are not two different algorithms. They are two propagation rules for the same algorithm, distinguished by `decisionₙ`:

```text
decisionₙ = repair_dispatch          → within-scope: same n, re-emit, re-fire
decisionₙ ∈ {accept, release}       → cross-scope: closed_cellₙ → α-matter at n+1
decisionₙ = override                 → cross-scope: closed_cellₙ → α-matter at n+1 (degraded)
decisionₙ = reject                   → cross-scope: closed_cellₙ terminates at n, no projection
```

A downstream consumer reading a closed cell determines the recursion mode by reading `decisionₙ`. The mode is not a separate field; it is the structural consequence of the decision the kernel pins.

### What the recursion does not do

Two non-modes are explicit to head off conflation:

- **Repair-dispatch is not scope-advancing.** A child cell under a repair contract operates at the same scope `n` as its parent. Treating repair-dispatch as `n → n+1` would let the parent γ never re-emit, which would break the kernel's step-3-then-step-4 ordering.
- **Accept is not same-scope.** A cell that closes with `decision ∈ {accept, release}` projects to scope `n+1`; the scope-`n` actors do not continue to operate on the closed cell. Treating accept as same-scope would let scope-`n` actors mutate matter that has already crossed the boundary, which would break the kernel's typed-handoff property.

The recursion is well-formed only when both modes are named and the two non-modes above are excluded by construction.
