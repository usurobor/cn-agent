<!-- sections: [Preamble, Kernel, Cell Outcomes, Recursion Modes, Scope-Lift, Two-Layer Separation, Non-goals, Closure] -->
<!-- completed: [Preamble, Kernel, Cell Outcomes, Recursion Modes, Scope-Lift, Two-Layer Separation, Non-goals, Closure] -->

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

---

## Scope-Lift

The cross-scope recursion mode (§Recursion Modes) advances the scope index from `n` to `n+1`. What carries across that scope boundary is named here as **three projections under scope-lift**. The framing is load-bearing: these are projections of scope-`n` objects onto scope-`n+1` roles, not a flat renaming of roles inside a single cell. The same actor wearing a δ hat at scope `n` and a β hat at scope `n+1` is not "δ becomes β at the next scope" — it is a δ-decision at scope `n` projecting onto β-discrimination at scope `n+1`. The object that projects is the decision (and the closed cell, and the receipt-stream observation); the role at scope `n+1` is the role that receives the projection.

### The three projections

```text
Projection 1:  closed (αₙ, βₙ, γₙ) cell        →   αₙ₊₁ matter
Projection 2:  δₙ boundary decision             →   βₙ₊₁-like discrimination
Projection 3:  εₙ receipt-stream observation    →   γₙ₊₁-like coordination / evolution
```

Each projection is a typed mapping from a scope-`n` object to a scope-`n+1` role function.

**Projection 1 — the closed cell projects as α-matter.**

The closed cell at scope `n` — the typed object `{contractₙ, matterₙ, reviewₙ, receiptₙ, verdictₙ, decisionₙ}` (and, under degraded outcomes, the override block) — is what scope-`n+1` α reads as one of its α-input objects. The cell is no longer "a project that happened"; it is "a typed input to the next-scope α's production." The receipt is the carrier (§Recursion Modes); the closed cell is what the parent scope reasons over.

**Projection 2 — δ's boundary decision projects as β-like discrimination at the parent scope.**

δₙ's recorded decision — the boundary action that determined whether the cell crossed and how — projects onto scope `n+1` as discrimination input. The scope-`n+1` β does not re-execute δₙ's decision; the scope-`n+1` β reads δₙ's decision as part of the matter it discriminates. δₙ-as-boundary-actor at scope `n` and βₙ₊₁-as-discriminator at scope `n+1` are not the same role; they perform analogous functions at different scopes, and δₙ's decision feeds βₙ₊₁'s work as one of its inputs. The qualifier "β-like" matters: scope-`n+1` β is β at scope `n+1` (a new role, not δₙ renamed), and δₙ's decision is one of the typed objects that β-at-scope-`n+1` discriminates.

**Projection 3 — ε's receipt-stream observation projects as γ-like coordination / evolution at the parent scope.**

εₙ observes the receipt stream across cells at scope `n` and detects protocol gaps — patterns of receipt-stream incoherence that name a missing predicate, a missing artifact, or a missing rule the kernel does not yet enforce. εₙ's observations project onto scope `n+1` as coordination and evolution input: scope-`n+1` γ reads the receipt-stream observation as one of the matters its scope-`n+1` cells will be contracted against, and the protocol-evolution patches εₙ emits become matter for scope-`n+1` α to produce. εₙ-as-protocol-observer at scope `n` and γₙ₊₁-as-coordinator-and-evolver at scope `n+1` are analogous in role function but not identical actors; the projection carries the observation, not the role.

### Projection-not-renaming framing

The framing is explicit:

> **This is a projection under scope-lift. It is not a flat role-renaming inside a single cell. δₙ is not literally βₙ₊₁; εₙ is not literally γₙ₊₁. The scope-`n+1` β and γ are full roles at scope `n+1` performing β-discrimination and γ-coordination over scope-`n+1` matter; what projects from scope `n` is the decision (from δₙ) and the receipt-stream observation (from εₙ) as typed objects scope-`n+1` β and γ consume as part of their work.**

The reason the framing matters: a reader who parses "δₙ → βₙ₊₁" as flat renaming will conclude that δ at scope `n` is the same role as β at scope `n+1`, which would collapse the kernel into a single-cell algorithm operating at one scope with different role labels at different times. That reading erases the scope index, breaks the recursion, and reintroduces the trust-by-seniority mode the predecessor doctrine rejects. The scope-`n+1` cell has its own α, β, γ, δ — and the projection feeds them rather than replacing them.

### β and γ have no upward projection

**βₙ and γₙ have no projection to scope `n+1`.** Their work is intra-scope: β discriminates matter against the contract at scope `n`; γ closes the cell and emits the receipt at scope `n`. Neither βₙ nor γₙ produces an object that crosses to scope `n+1` as a separately typed projection.

What does cross from scope `n` to scope `n+1` carrying βₙ's and γₙ's work is the **closed cell itself** (projection 1) — and the closed cell carries `reviewₙ` (β's record) and `receiptₙ` (γ's record) as components. Scope-`n+1` α reads the closed cell, which means scope-`n+1` α reads the embedded β review record and γ receipt record as part of the closed-cell object. But βₙ and γₙ do not have their own upward projections distinct from the closed cell. Only δ (boundary decision) and ε (receipt-stream observation) have separately typed projections that scope-`n+1` β and γ receive directly.

The asymmetry is structural. β and γ are **intra-cell roles** — they operate inside the cell at scope `n` and their work is fully captured by the closed cell. δ and ε are **boundary-or-cross-cell roles** — δ is the cell's boundary actor (it operates at the membrane between the cell and the parent scope), ε is a cross-cell observer (it operates across the receipt stream of many cells at scope `n`). The boundary-or-cross-cell roles project upward because their work is already at or beyond the cell's boundary; the intra-cell roles do not project upward because their work is *the cell*.

### Alignment with the receptor design

The three-projection structure aligns with `#369` AC3 — the scope-lift projection that the receptor's schema work is typing. The kernel pins the *what* (three projections; the projection-not-renaming framing; β/γ-no-upward-projection); the schema work pins the *type* (the field shapes that carry the projected objects across the boundary). When the schema work lands, the three projections become typed handoffs the parent scope reads as fields on the received closed cell, the received decision, and the received protocol-observation patch.

Where this document and the receptor's schema work disagree on the projection structure, the kernel statement here is the load-bearing claim and the schema work aligns to it — not the reverse.

### What scope-lift does not do

Three non-projections, stated to head off conflation:

- **β has no upward projection.** βₙ's review is part of the closed cell; it does not have a separately typed projection at scope `n+1`. Treating β as having an upward projection would let scope-`n+1` reason over the review independent of the closed cell, which would break the receipt's role as the single typed handoff surface.
- **γ has no upward projection.** γₙ's receipt is the parent-facing artifact of the closed cell; it is *what carries* projection 1, not a separate projection. Treating γ as having an upward projection distinct from the closed cell would double-count the receipt's role.
- **Scope-lift is not a single-cell loop.** The scope-`n+1` cell is a new cell with new α, β, γ, δ, ε at scope `n+1`. It is not "the same cell at the next scope"; the scope-`n+1` cell has its own contract, its own matter, its own review, its own receipt, its own verdict, its own decision. The scope-`n` cell projects *into* the scope-`n+1` cell's inputs; it does not *become* the scope-`n+1` cell.

---

## Two-Layer Separation

This document is the **kernel layer**. The kernel layer names the *what* — the substrate-independent recursion algorithm that any instantiation of the role-scope ladder pattern can carry. The kernel does not specify a particular tooling, a particular schema language, a particular invocation surface, or a particular platform. Substrate-specific facts live in the **realization layer**, which expands the kernel on a chosen substrate.

> **Kernel is the *what*. Realization is the *how-on-this-substrate*.**

The split protects two properties simultaneously. The kernel stays portable: a future re-instantiation of the role-scope ladder on a different substrate inherits the kernel without re-derivation. The realization stays expandable: substrate-specific concerns (schema syntax, command-wrapper shape, dispatch mechanism, validator implementation) can evolve without dragging the kernel along, because the kernel does not name them.

### Realization peers

Four realization-layer surfaces expand this kernel on the current substrate. Each is cited here as a realization peer; none of them is edited by this cycle. The kernel is the doctrine; the realization peers are the substrate expansions.

**`RECEIPT-VALIDATION.md` — parent-facing typed-interface design.**
This is the frozen validation interface for `V` — input contract (what `V` reads as references), output contract (`ValidationVerdict` shape), invocation contract (when and how δ invokes `V`), and the structural distinction between `ValidationVerdict` (V-emitted) and `BoundaryDecision` (δ-recorded). The kernel positions `V` at step 4 and δ at step 5; the receptor design types the interface and resolves the five Open Questions seeded by the predecessor doctrine. Where the kernel commits to "V dereferences evidence from the receipt," the receptor design types the dereferencing as `EvidenceRootRef` and pins the override-detection biconditional.

**`schemas/cdd/` (Phase 2, `#369`, in flight) — typed schemas for receipt, contract, boundary decision.**
This is the realization peer that types the four cell outcomes (§Cell Outcomes) and the override block. The schema work pins the verdict × action × transmissibility table that the kernel's `(verdict, decision)` preconditions describe at a doctrine level. The kernel and the schemas compose: when the schemas land, the kernel's outcome preconditions become enforceable as typed-field constraints on the receipt's `validation` and `boundary` blocks.

**`cn-cdd-verify` (Phase 3, deferred) — V's command-wrapper implementation.**
This is the realization peer that implements `V` as an operator-facing command wrapping the predicate the kernel names at step 4. The command surface is the substrate-specific invocation path operators and the harness consume; the predicate (exposed as a capability per the receptor design) is what δ-the-skill consumes. Both paths reach the same underlying predicate. The kernel commits to `V` being a predicate; the realization commits to the command-wrapper shape on this substrate.

**`CDD.md` (Phase 7 rewrite, deferred) — canonical executable algorithm.**
This is the realization peer that expands the kernel into the operational algorithm CDD runs as a working protocol on this substrate. The kernel pins the five-step closure, the four outcomes, the two recursion modes, and the three scope-lift projections; the operational rewrite expands each into role-specific steps, dispatch sequences, branch mechanics, harness contracts, and concrete artifacts. The kernel is the spine; `CDD.md` rewrite is the body. Until the Phase 7 rewrite lands, `CDD.md` remains the canonical algorithm at its current shape — this cycle does not touch it.

### What lives where

| Concern | Kernel layer (this doc) | Realization layer |
|---|---|---|
| Recursion algorithm (five steps, four outcomes, two modes, three projections) | Yes | No (cited; not re-stated) |
| Substrate-independent role signatures (αₙ.produce, βₙ.review, …) | Yes | No |
| Validator predicate signature (`V : Contract × Receipt → Verdict`) | Yes | No |
| Typed schemas (receipt, contract, boundary, override fields) | No | `schemas/cdd/` (Phase 2) |
| Command-wrapper shape (operator-facing) | No | `cn-cdd-verify` (Phase 3) |
| Capability binding, transport, dispatch | No | δ skill + harness (Phase 4) |
| Per-role operational steps (artifact order, branch mechanics, CI gates) | No | `CDD.md` (Phase 7) |
| Receptor interface freeze (input refs, output verdict shape, invocation rule) | Position-only (positions `V` at step 4, δ at step 5) | `RECEIPT-VALIDATION.md` (Phase 1) |

The table is a guide, not a contract — the realization peers may evolve in ways the kernel does not anticipate. The load-bearing rule is the **direction of dependency**: realization cites kernel; kernel does not cite realization. The kernel sections of this document name no substrate; the §Two-Layer Separation section is where realization peers are cited because this section is in the realization layer of the document itself.

### Why the split is load-bearing

Three failure modes the split prevents:

1. **Operational drift in the kernel.** If kernel sections cite substrate-specific tooling, the kernel becomes substrate-specific. A future re-instantiation (a different ladder pattern, a different platform) cannot inherit the kernel without re-deriving it. The kernel sections of this document name only roles, verdicts, decisions, receipts, evidence, scopes — they do not name any particular tooling.
2. **Kernel collapse mid-rewrite.** Phase 7 rewrites `CDD.md` into the operational layer expanding this kernel. If the kernel were not stated independently first, Phase 7 would have to derive the kernel mid-rewrite — which is exactly the failure mode the split exists to prevent. The kernel lands first, then the rewrite expands it.
3. **Schema-kernel disagreement.** Phase 2 types receipt, contract, and boundary decision. If the kernel did not pin the four outcomes and three projections at a doctrine level, the schemas would be choosing the shape rather than typing a doctrine-defined shape. The kernel statement here is the load-bearing claim; the schema work aligns to it.

---

## Non-goals

This cycle's surface-containment contract is binding. The diff for this cycle, before α signals review-readiness, is constrained to exactly:

- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (new file, status `A`)
- `.cdd/unreleased/370/*.md` (cycle evidence — scaffold, dispatch prompts, self-coherence, role close-outs as they land)

**Files this document does not edit:**

- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — predecessor doctrine; remains the receipt-rule doctrine surface unchanged
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — receptor design; cited as realization peer; unchanged
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical executable algorithm; Phase 7 rewrite target; unchanged this cycle
- `schemas/cdd/` — Phase 2 (`#369`), runs in parallel; this cycle adds no schema files
- `cn-cdd-verify` command — Phase 3, deferred; no edits
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — Phase 4 (δ split) target; unchanged
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — Phase 5 (γ shrink) target; unchanged
- `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` — Phase 6 (ε relocation) target; unchanged
- `ROLES.md` — Phase 6 target; unchanged
- CI workflows, new packages, new tooling

**Choices this document does not make:**

- Schema syntax and exact field types for the receipt and contract (Phase 2)
- The capability identifier and package-manifest binding for `V` (Phase 3)
- The transport mechanism δ uses to invoke the capability (Phase 3)
- The argv/stdout shape of `cn-cdd-verify` post-refactor (Phase 3)
- The exact δ-skill split between boundary and harness (Phase 4)
- The exact `ROLES.md` text that lifts ε's generic doctrine (Phase 6)
- The operational algorithm details — artifact order, branch mechanics, dispatch sequences (Phase 7)

**Author-discipline non-goals:**

- The kernel is **decisive over exhaustive.** Each section commits to one position; rationale exists to anchor the position, not to enumerate every alternative the design could have made. The full alternative space lives in the predecessor doctrine's Open Questions and the receptor design's narrative.
- The kernel does not preempt the Phase 7 rewrite. The operational realization is `CDD.md`'s job; this document gives Phase 7 a citable spine, not the body.
- The kernel does not grow to operational length. The target is 200–400 lines; growth past that signals re-derivation of the operational layer mid-cycle.

---

## Closure

The kernel is complete when:

1. The five-step recursion at scope `n` is stated with explicit signatures, including β's matter-only consumption and δ's receipt-and-verdict-only consumption (§Kernel).
2. The evidence-binding rule is pinned at the kernel level — evidence accumulates during α and β work; γ binds it into the receipt; `V` dereferences from the receipt; β never consumes evidence; δ never re-reads evidence (§Kernel).
3. The four closed-cell outcomes are pinned with their `(verdict, decision)` preconditions and `invalid` is declared non-terminal (§Cell Outcomes).
4. The two recursion modes are distinguished — within-scope repair-dispatch at the same scope index, cross-scope accept / degraded / reject at the advancing scope index — and the same recursion operator is named (§Recursion Modes).
5. The three scope-lift projections are named — closed cell → α-matter; δ decision → β-like discrimination; ε receipt-stream observation → γ-like coordination / evolution — with the projection-not-renaming framing and the β/γ-no-upward-projection clause (§Scope-Lift).
6. The two-layer separation is declared and the four realization peers are cited (§Two-Layer Separation).
7. The surface-containment contract is respected: the diff is exactly `{COHERENCE-CELL-NORMAL-FORM.md} ∪ .cdd/unreleased/370/*.md` (§Non-goals).

Subsequent phases of `#366` inherit this document as their input contract:

| Phase | Inherits from this document |
|---|---|
| Phase 2 — `receipt.cue` + `contract.cue` (parallel, `#369`) | §Cell Outcomes (four-outcome preconditions); §Kernel (evidence-binding rule); §Scope-Lift (the three projected objects' shapes) |
| Phase 3 — `cn-cdd-verify` refactor | §Kernel (V's signature, evidence-binding rule); §Cell Outcomes (verdict-decision composition); §Two-Layer Separation (V is capability + command-wrapper) |
| Phase 4 — δ split | §Kernel (δ's signature, no-evidence consumption); §Cell Outcomes (δ's authoritative composition with the verdict); §Recursion Modes (δ records `repair_dispatch` for within-scope; `accept`/`release`/`override`/`reject` for cross-scope) |
| Phase 5 — γ shrink | §Kernel (γ's signature; γ binds evidence into the receipt); §Scope-Lift (γₙ has no upward projection — γ is intra-cell) |
| Phase 6 — ε relocation | §Scope-Lift (εₙ's receipt-stream observation projects as γₙ₊₁-like coordination/evolution); §Kernel (ε is cross-cell, not in the five-step closure) |
| Phase 7 — `CDD.md` rewrite | §Kernel + §Cell Outcomes + §Recursion Modes + §Scope-Lift (the entire kernel as the spine the operational rewrite expands) |

When Phase 7 lands and rewrites `CDD.md` around this kernel, this document moves from "draft doctrine the operational layer cites" to "doctrine the operational layer satisfies." At that point, this document's status line in the preamble is updated by the Phase 7 cycle's authoring — not by this cycle — and the kernel becomes the ratified contract that the operational realization expands.
