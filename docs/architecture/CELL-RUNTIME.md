# Cell Runtime — Classes, Domains, and the Generic Runner

**Status:** Proposed architecture note (realization layer). Not ratified. Ratification is the job of a CDD `doctrine` cell (see #628 / the #627 wave).
**Owns:** How the one CCNF kernel cell *deploys* at CNOS operating scale — three output-telos cell **classes** (WC/PC/CC), the matter/contract **domains** they carry, and the single generic **runner** that executes them.
**Does-not-own:** The kernel algorithm (`COHERENCE-CELL-NORMAL-FORM.md`), the receipt rule and four surfaces (`COHERENCE-CELL.md`), the validator interface (`RECEIPT-VALIDATION.md`), the CUE schemas (`schemas/cdd/`), or the #366 phase roadmap. This note is a **realization peer** that cites them; it does not restate or supersede them.

---

## Thesis

CNOS has **one** recursive coherence cell — the CCNF five-step kernel: `α.produce → β.review → γ.close → V.validate → δ.decide`, recursing across scopes. It does not have a closed roster of cell "kinds," and it does not need a runner per kind. At operating scale the *one* kernel is **deployed** in three output-telos shapes, executed by **one generic runner**:

```
WC — Working Cell    — matter targets TSC-α (pattern)     — produces an artifact
PC — Planning Cell   — matter targets TSC-β (relational)  — produces a relation graph
CC — Cohering Cell   — matter targets TSC-γ (process)     — produces a process judgment
```

This corrects two framings at once. It corrects the "canonical cell kinds" reading (the `CELL-KINDS` names are *domains*, not kinds). And it corrects #627's "WC/PC/CC as a lens over the kinds" reading (WC/PC/CC are not a lens over the domains — they are **deployment shapes of the kernel**, orthogonal to the domains).

## Four orthogonal classifications (the core correction)

The confusion this note removes is a single conflated axis where there are **four**:

| # | Classification | Values | What it answers |
|---|---|---|---|
| 1 | **Kernel cell** (canonical unit) | one CCNF cell | *What is the unit of work?* — always the five-step closure |
| 2 | **Cell roles** | α, β, γ, δ (+ ε, cross-cell) | *Who acts inside the cell?* — α/β/γ/δ run the five-step closure; ε is the cross-cell receipt-stream observer, **not** a step in it |
| 3 | **Output-telos class** | WC, PC, CC | *Which TSC axis does the cell's matter target?* |
| 4 | **Matter / contract domain** | implementation, doctrine, … | *What kind of matter, in what mode?* |

These are independent. A single cell is simultaneously: the one kernel (1), running α→β→γ→V→δ (2), of some class (3), over some domain (4). Example: a doctrine change is *one kernel cell*, running the internal closure, of class *WC* (its matter is a pattern artifact — a doc), in domain *doctrine*. Collapsing any two of these axes is the surface-smearing failure `COHERENCE-CELL.md` §"Structural Prediction" names.

Axis 2's five-step closure is `α→β→γ→V→δ`. **ε is not a step in that closure** — it observes receipt streams *across* cells and belongs to the protocol-evolution layer (CCNF §Scope-Lift). Listing it beside α/β/γ/δ is a convenience, not a claim that ε runs inside every cell.

## WC/PC/CC are deployment shapes, not role-relabelings

**Load-bearing distinction.** WC/PC/CC name axis 3 (output telos). They are **not** axis 2 (α/β/γ roles) re-labeled into three sibling cells. Reading "WC=α, PC=β, CC=γ as three cells" is exactly the *flat-renaming collapse* CCNF §Scope-Lift explicitly forbids — it "erases the scope index, breaks the recursion," and it would break the α≠β **contagion firebreak** (`COHERENCE-CELL.md` §"β Independence"): β must review α's matter from outside α's frame *within the same cell*. Every WC, PC, and CC runs its own internal, independent α→β→γ→V→δ.

**Corollary — the β word.** Do not call PC a "discriminator." In CDD, *discriminate/review* is the intra-cell **β role** (axis 2). PC's dominant telos is **relational structure** (axis 3, TSC-β "do the parts fit") — it *constructs* a relation graph; it does not *review* matter. Reserve "discriminator/review" for β-the-role; use "relator / relational-structurer" for PC-the-class. (This is precisely the β-fracture #627 flagged; here it is resolved by keeping axis 2 and axis 3 apart.)

**S₃-invariance.** `tsc:spec/tsc-core.md` §0 makes axis symmetry normative — *"No axis privilege."* So "α/β/γ-dominant" is a claim about the cell's **telos**, never about the measurement privileging an axis. The CM stays S₃-symmetric; a WC simply reads its telos-axis score (`s_α`) as the primary signal for its own gate.

## The generic runner (realization substrate)

One runner executes all three classes. It is the *runtime-substrate* surface `COHERENCE-CELL.md` §"Two-Layer Separation" names — below δ. It owns the mechanical loop, and **only** the mechanical loop:

```
read contract → assemble input bundle → hash evidence → run class adapter
→ collect matter + receipt → run CM → run V → apply permitted bounded action
→ emit durable receipt → route next pulse
```

It does **not** own planning intelligence, worker implementation, validation policy (`V`), boundary policy (δ), or release effection — the four-surface rule holds. It extends the shipped `cn cell` family (`return`/`resume`/`finalize`) with `run`, `measure`, `bundle`, `act`, `pulse`, selected by `--class working|planning|cohering` (not `--kind implementation` — that repeats the taxonomy mistake).

**CM measures; V gates; δ effects — three separate surfaces (pinned to the shipped schema).**

```
CM → { s_α, s_β, s_γ, C_Σ (geom-mean, strict Degeneracy Axiom), bottleneck, witnesses, defects, provenance }   # tsc Core
V  → #ValidationVerdict.verdict ∈ { PASS, FAIL }                                                                # schemas/cdd/receipt.cue
δ  → #BoundaryDecision.action  ∈ { accept, release, reject, repair_dispatch, override }                         # schemas/cdd/boundary_decision.cue
runner routes → { continue, hold, escalate, human_gate, no_op_with_reason }                                     # runtime state, not a verdict
transmissibility (derived) → { accepted, not_transmissible, degraded }   from (verdict × action)
```

CM computes V's mechanical score; V is the same `V : Contract × Receipt → Verdict` CCNF pins at step 4 and #366 Phase 3 implements — **not a second validator, and this note does not change the verdict enum.** `ESCALATE` is runner routing, not a verdict; `override` (carrying `#Override`, a degraded boundary action) is a δ decision, not a verdict. tsc supplies CM's *measurement* semantics — including the Θ-threshold/witness/provenance logic of tsc Operational — which CNOS adapts into CDD's typed `CM → V → δ` boundary rather than copying tsc's operational state-names into `ValidationVerdict`. (Θ default 0.75; any richer verdict state, e.g. `WARN`, is a future schema/validator migration, a #366 Phase 3 concern.)

**Mechanical-first, no silent no-op.** The next pulse does not call an LLM by default. Run CM mechanically; escalate only on a **deterministic, logged** predicate (evidence missing / axis score in ambiguity band / >1 valid next-move / sources contradict / Θ unevaluable). Same bundle hash ⟹ same escalation. Every pulse emits a receipt; a no-op states *why* (`hold | blocked | waiting_human_gate | already_satisfied | no_unblocked_cell | evidence_missing`) — never an empty green. (Motivated by this session's silent-wake incident: green ~30× over ~22h, zero commits.)

## Matter/contract domains (the ex-"cell kinds")

The `CELL-KINDS.md` names are demoted to axis 4 — but they are **not homogeneous**, and flattening them all to "matter_domain" loses information. Classify them:

| Sub-type | Members | Note |
|---|---|---|
| **Matter domain** (what α emits) | implementation, doctrine, audit, experiment, cleanup, planning | true domains; the common case (`planning` added by #614 — its matter is the plan/wave contract) |
| **Contract mode** (how the cell is contracted) | repair, recovery | modifiers on an existing cell's contract, not fresh matter |
| **Coordination / boundary** | wave (scope structure), release (δ boundary action) | not matter at all — a scope, and a δ effect |

`CELL-KINDS.md`'s content is **demoted, not deleted** — it becomes the domain vocabulary. Its own rule ("cell kind does not fork the kernel") is what licenses this: a domain was always a typed refinement, never an ontological type. The re-head is **in place** (filename kept for link stability); a physical rename to `CELL-DOMAINS.md` + inbound-link sweep is a deferred follow-up. Compatibility: legacy `cell_kind: X` → `cell_class + matter_domain`.

## Human gates and safety (unchanged from #627/S1)

Mechanical dispatch may read state, assemble bundles, run CM/V, apply launch labels, open/update issues+PRs, emit receipts, hold with reason, request a gate. It may **not** by default merge product PRs, override a failed V, declare done-done with unresolved blockers, change release authority, invent scope, write product code from CC, or bypass independent β review. CC v0 is **receipts-only**. Human approval gates: operator intent, material merges, V-override, done-done shipment, security/permission escalation, release authority.

## Reconciliation map (what this note extends, and how)

| Existing surface | Status | This note's relation |
|---|---|---|
| `COHERENCE-CELL-NORMAL-FORM.md` (#370) | Landed kernel doctrine | **Parent.** WC/PC/CC deploy this exact five-role kernel; this note is its realization peer. |
| `COHERENCE-CELL.md` (#364) | Landed doctrine | Four surfaces + firebreak inherited unchanged. |
| `CELL-KINDS.md` (#570) | Landed doctrine | Re-headed **in place** as domain vocabulary (axis 4), heterogeneity classified; physical rename to `CELL-DOMAINS.md` deferred. Constitutive change → migration note in the file. |
| `schemas/cdd/{contract,receipt,boundary_decision}.cue` (#369) | Shipped | **Extend** with `cell_class` + `matter_domain` + CM fields; do not recreate. A `cell.cue` descriptor is the new artifact. |
| `#366` roadmap Phases 3–7 (V validator, δ-split, CDD.md rewrite) | Open | The runner/CM/V work belongs **inside** this roadmap, not parallel to it. Phase 3's `V` **consumes** this note's CM output as one validation input; CM is the measurement sub-surface — **not** a second validator, and **not** the same surface as `V` (which stays the typed `PASS\|FAIL` gate). |
| `cn cell` Go family (#500/#593) | Shipped | `run/measure/bundle/act/pulse` extend it. |
| `usurobor/tsc` Core/Operational + modes | Shipped | Supplies CM semantics (α/β/γ, geometric-mean+degeneracy, Θ verdict, `mechanical|llm|hybrid|auto`). Measurement provider, not a dispatch protocol. |

## Non-goals

A runner-per-kind or a wake-per-kind; LLM-first dispatch; auto-merge or gate-crossing authority; a global cross-repo scheduler; deleting `CELL-KINDS.md` content; rewriting `CDD.md` (Phase 7); new agent-personality doctrine. First cut makes the four axes explicit, the runner generic, and the measurement mechanical.

## Open questions

1. **Class ⟂ domain typing:** is `cell_class` a required top-level field alongside `matter_domain`, and does the FSM's `CellKind` observation seam migrate to both? (Recommended: yes; add both, keep `cell_kind` as a compat alias.)
2. **Does this fold into #366 as new Phases, or a sibling wave?** (Recommended: new Phases 8–10 under #366, since it extends the same kernel.)
3. **CM↔V vs #366 Phase 3:** are they one surface? (**No** — Phase 3's `V` *consumes* CM output; CM is a measurement sub-surface, not a second validator; `V` stays the `PASS|FAIL` gate.)
4. **`CELL-KINDS.md` rename vs augment-in-place** — this PR chose augment-in-place (link stability); the physical rename to `CELL-DOMAINS.md` is a deferred follow-up.

---

*Authoring note (κ).* Drafted by Sigma-as-κ from the operator's "triadic cell runtime" architecture note (2026-07-07), reconciled against landed cnos doctrine the operator's draft predated — chiefly `COHERENCE-CELL-NORMAL-FORM.md` (#370), the #366 roadmap, and the shipped `schemas/cdd/`. The operator's core thesis (one recursive cell; the `CELL-KINDS` names are not canonical types; one generic runner; CM measures / V gates / δ effects; mechanical-first; no silent no-op) is preserved. κ's reconciliation: (1) the four-orthogonal-axes frame; (2) WC/PC/CC as output-telos deployment shapes, not α/β/γ re-labeled (preserving the α≠β firebreak per CCNF); (3) "PC = relator," not "discriminator"; (4) the domains are heterogeneous (matter / contract-mode / coordination); (5) position the work inside #366, not parallel to it; (6) ε is cross-cell, not an internal-closure role; CM feeds V (not the same surface). Proposed, not ratified — a `doctrine` cell ratifies.