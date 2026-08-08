# Cell Runner — the case ladder

**Status:** Design note for β review (Pi). Companion to the kernel doctrine
`src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (CCNF) and the
ratified architecture `cnos#711` (recursive-with-predicates). Grounded in the
reference implementation `src/go/internal/cellkernel`.

## Purpose

One kernel runs every cell. This note pins **exactly how the runner behaves
across the ladder of cases** — from the empty cell to full recursion — so there
is one shared picture and no per-case confusion. Every case is the *same*
five-step closure; only **how the seats are filled** changes.

## The kernel (recap)

At one scope, a cell is a five-step closure over a `contract`:

```
1. matter   := α.produce(contract)
2. review   := β.review(contract, matter)
3. receipt  := γ.close(contract, matter, review, evidence)
4. verdict  := V(contract, receipt)
5. decision := δ.decide(receipt, verdict)
outcome     := f(verdict, decision)   ∈ {accepted, degraded, blocked, invalid}
```

Reference types: `Alpha`, `Beta`, `Gamma`, `Validator`(V), `Delta` interfaces;
`Run(ctx, Spec) (ClosedCell, error)`. `α` and `β` are the caller's
customization; `γ/V/δ` default to the mechanical kernel.

### Invariants that hold in **every** case

- **I1 — Trust moves by the typed receipt.** The parent trusts a cell only when
  **V validates the receipt** and **δ records a decision** — never because γ
  "closed." (CCNF: trust by typed receipt, not role seniority.)
- **I2 — Cognition is rentable only at α and β.** `γ` (close, pure), `V`
  (validate, mechanical/CUE), `δ` (decide, policy) are **always compiled**. V
  being mechanical is precisely what keeps the trust surface deterministic and
  reproducible.
- **I3 — Escalation is a decorator on the α/β seat**, keyed on a
  runtime-constructed, **hashed bundle**; **deterministic** (same bundle hash →
  same compiled-vs-rented decision) and **logged as evidence**. The kernel's
  `Run` is unchanged by it.
- **I4 — Two failure kinds are distinct.** A seat returning an **error** is a
  *malfunction* (the cell does not close; `Run` errors). A review returning
  `Pass:false` is *contract-unmet* (the cell closes `blocked`).
- **I5 — Loops are bounded.** The repair loop (`repair_dispatch`) runs under a
  **preregistered attempt budget**; on exhaustion it emits a terminal `held`
  (contract-unmet) or `failed` (malfunction) receipt, never an infinite spin.
- **I6 — Evidence-binding.** Evidence accrues during α/β work → **γ binds it
  into the receipt** → **V dereferences it** → β never consumes evidence, δ
  never re-reads it. Escalation decisions are part of that evidence.

---

## The case ladder

Each case names the **contract**, which **seats are compiled vs rented**, and
the **walk** through the five steps.

### Case 0 — Empty cell (all compiled)

- **Contract:** empty.
- **Seats:** α = noop (`Matter{}`), β = accept, γ/V/δ = mechanical default.
- **Walk:** `∅ → review.pass → receipt(cell-0) → verdict.pass → accept →
  accepted`.
- **Proves:** the loop turns; all five seats fire. Zero cognition.
- **Status:** implemented in `cellkernel` (`EmptySpec`).

### Case 1 — Bool cell (all compiled, real repair loop)

- **Contract:** "produce a bool that is `true`."
- **Seats:** α produces a bool (compiled — e.g. a deterministic/pseudo source),
  β checks `matter == true` (compiled — a CUE constraint `matter: true`), γ/V/δ
  mechanical.
- **Walk:** on `false`, V FAIL → δ `repair_dispatch` → α re-produces with the
  why → re-test; bounded by the attempt budget (I5).
- **New vs Case 0:** the **repair loop** with a mechanical β. Still **zero
  cognition** — everything is CUE-decidable.
- **Status:** seats/outcome/malfunction implemented; the `repair_dispatch` drive
  loop is a named increment (see §Recursion).

### Case 2 — α rents cognition (mini-CDS)

- **Contract:** "write a function that passes test `T`."
- **Seats:** **α rents cognition** — there is *no compiled path* that
  synthesizes arbitrary code, so the escalation predicate fires at α (I3); α
  calls a cognition backend. **β/V are compiled** — they *run the test*;
  PASS/FAIL is deterministic. **δ** loops on FAIL back to α.
- **Why V stays mechanical:** the test **is** the oracle. Cognition at α costs
  nothing on the trust surface (I1/I2).
- **New:** first cell that needs a backend; first genuine escalation. The
  compiled path for a seat is *empty*, so the trigger *must* fire.
- **Status:** increment — needs the escalating-α decorator + `dispatch.Backend`
  wiring + a contract carrying acceptance.

### Case 3 — β also rents cognition (full CDS)

- **Contract:** "write feature X and review it to the review skill."
- **Seats:** **α rents cognition** (writes). **β rents cognition for the
  judgment residue** — a thorough review following the cnos review skill
  (`cdd/beta`), while **mechanizing everything it can** (run tests, grep
  patterns, AC-presence stay compiled; process-economics). **V stays
  mechanical**: it validates the **receipt** — that β produced the *required
  evidence* (checks ran, ACs covered, independence/firebreak recorded) and the
  receipt is structurally consistent. **V never re-judges the code.**
- **Trust model:** correctness of β's judgment is trusted via **discipline (the
  review skill) + independence (β ≠ α) + the evidence β must produce + V's
  structural validation** (I1/I6) — not by V re-deriving the review.
- **New:** both open-ended seats rent cognition; the trust surface still
  deterministic. This is the full **CDS cell**: α writes, β reviews per skill,
  γ/V/δ hold the mechanical spine.
- **Status:** increment — escalating-β decorator carrying the review-skill
  bundle; V extended to validate required-evidence presence.

### Case 4 — Recursion (cell of cells)

Two mechanisms, **both re-entrant `Run`** — the kernel's five steps never change;
recursion is `Run` calling `Run`.

**(a) Compositional recursion — α/β are themselves cells.** (CCNF scope-lift;
#711 §2.)

- **Parent α = decomposition:** its `Produce` turns the parent contract into
  **child contracts / an execution graph**, runs each child via `Run(ctx,
  childSpec)`, and composes the child `ClosedCell`s into the parent matter.
- **Parent β = composition:** its `Review` judges whether the accepted **child
  receipts jointly satisfy the parent contract** — the composition oracle:
  *"name one way every child receipt could be locally valid while the parent
  contract remains unmet."* A β that only counts child PASSes is degenerate.
- **Cognition placement:** leaf α/β rent cognition where the escalation
  predicate fires (Cases 2/3); interior decompose/compose nodes may be
  **all-compiled**. Novelty concentrates at the leaves, not "cognition lives at
  the leaves" (I3).
- **Task ≠ episode:** a durable task survives multiple execution episodes;
  task state is *projected from* episodes, not identical to one episode's FSM
  (#711 §3) — relevant once a decomposition spans retries.

**(b) Repair recursion — `repair_dispatch` (within-scope).** (CCNF §Recursion
Modes.)

- δ decides `repair_dispatch` → a **child cell at the same scope** runs under a
  **repair contract** derived from the failure (the failed predicates + β's
  why). On child `accepted`, the parent γ **re-emits** a fresh receipt; V and δ
  **re-fire** at the parent. Scope index unchanged. Bounded by the attempt
  budget (I5).
- In the reference this is the **`Drive`** loop wrapping `Run`: run one closure;
  if the outcome is `blocked`/`repair_dispatch`, derive a repair contract and
  run again with feedback, until a terminal outcome (`accepted`/`degraded`/
  `reject`) or budget exhaustion.

**Cross-scope projection (scope-lift):** an `accepted`/`degraded` closed cell
projects **as α-matter at scope n+1** via its receipt; that is just mechanism
(a) viewed from the parent — the parent's α reads the child's `ClosedCell`.

---

## Mapping to the reference implementation

`src/go/internal/cellkernel` today:

- **Case 0** — implemented (`EmptySpec`, `Run` → `accepted`).
- **Cases 1–3** — the **seats, outcome mapping, and malfunction-vs-contract-unmet
  distinction** are implemented; what each case adds is a *seat implementation*
  (compiled bool / escalating-α / escalating-β) plus, for Case 1+, the `Drive`
  loop. `Run` itself does not change.
- **Case 4** — compositional recursion is a **seat that calls `Run`**; repair
  recursion is the **`Drive`** loop. Both are increments; the `Alpha`/`Beta`
  interfaces and re-entrant `Run` are the seams already present.

**Nothing in the ladder changes the five-step `Run`.** Every case is a different
*fill* of the seats plus (for the loop) a `Drive` wrapper. That is the whole
claim of the runner.

---

## What to build, in order

1. `Drive(ctx, spec, budget)` — loop `Run` under `repair_dispatch` to a terminal
   outcome (enables Case 1's real loop).
2. `Contract` carrying acceptance + a compiled CUE-backed `V`/`β` (Case 1 fully
   mechanical).
3. Escalating-α decorator + `dispatch.Backend` wiring + `cn cds build --issue N`
   (Case 2).
4. Escalating-β carrying the review-skill bundle; V validates required-evidence
   presence (Case 3 = full CDS).
5. Composite α/β (`Produce`/`Review` that call `Run`) + parent-β composition
   oracle (Case 4a); task≠episode projection when decompositions span retries.

---

## Open questions for β (Pi)

1. **Repair contract derivation.** Is `Drive` the right home for repair-contract
   synthesis, or does that belong in δ (which *decides* `repair_dispatch`)? Where
   is the boundary between "δ decides to repair" and "who writes the repair
   contract"?
2. **Bundle + escalation determinism.** Is "hashed bundle → same escalation
   decision" fully specifiable at the seat, or does part of it depend on the
   backend being pinned? (Same concern you raised on #715 AC1/AC3.)
3. **Composition oracle as a seat.** Should parent-β composition be a distinct
   `Beta` impl, or is it a *contract shape* the same β discriminates against?
4. **V's required-evidence check (Case 3).** What is the minimal typed shape of
   "required evidence present" so V stays a CUE-decidable predicate over the
   receipt rather than creeping into judgment?
5. **Anything missing from the ladder** — a case between these rungs, or a
   degraded/override path (Case 3+ with `override`) we have not drawn.
