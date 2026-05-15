<!-- sections: [Preamble, Q1, Q2, Q3, Q4, Q5, Validation Interface, Non-goals, Closure] -->
<!-- completed: [Preamble, Q1, Q2, Q3] -->

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

## Q1 — When does V fire?

**Chosen position.** `V`'s authoritative firing point is **δ-boundary validation**, invoked by δ after γ emits the receipt and before δ accepts, rejects, repair-dispatches, overrides, or releases. γ may invoke `V` as a non-authoritative preflight while drafting the receipt, but γ preflight never substitutes for δ boundary validation. The parent scope may not treat the cell as one of its inputs until the δ-boundary invocation produces a `ValidationVerdict` and δ records a `BoundaryDecision`.

### Rationale

The doctrine fixes where authority lives. From `COHERENCE-CELL.md` §Trust Boundary, validity is a property of `(contract, evidence, receipt)`. From §δ Boundary Complex, "δ invokes validation; δ does not embody validation." Together these say: the cell's transmissibility to the parent scope is a property the boundary observes, not a property a role internal to the cell confers. γ closes the cell and emits the receipt; that is γ's biological function. γ does not also issue the verdict that makes the cell transmissible. If γ both produced the receipt and pronounced it valid, the parent scope would be trusting the cell on γ's seniority, which is exactly the trust-by-approval mode the deep invariant rejects.

Putting `V` authoritatively at the δ boundary keeps two facts separate:

1. **Closure** is γ's function. γ fuses contract, matter, and review into a closed-cell record and emits the receipt. The closure record exists once γ has done this, independent of whether `V` returns PASS.
2. **Transmissibility** is δ's gate. The cell becomes parent-scope matter only when δ has invoked `V`, observed a verdict, and recorded a boundary decision (accept / reject / repair-dispatch / override / release).

A closed cell that has not been δ-validated is in an intermediate state: γ has done γ's work; the parent scope has not yet been handed an input. The receipt is emitted but not yet typed as accepted matter. This is the structural meaning of `closed_cell ≠ accepted_cell` from the cell recursion equation (`COHERENCE-CELL.md` §Recursion Equation).

### γ-preflight posture (non-authoritative)

γ is permitted to invoke `V` while drafting the receipt — running the same validator implementation against the contract, the receipt-in-progress, and the evidence graph as an informational check on whether the receipt γ is about to emit will pass δ-boundary validation. This is a useful affordance: it lets γ catch receipt-shape problems (missing fields, structural-evidence gaps, undisclosed degradation) before signaling closure, without requiring a round-trip through δ to discover them. A γ preflight that returns FAIL is information γ can act on before emitting the receipt.

The strict constraint: **γ preflight is non-authoritative**. Specifically:

- A γ-preflight PASS does not authorize δ to skip its own invocation. δ-boundary validation is not optional even when γ preflight reports PASS.
- A γ-preflight verdict is not recorded as the cell's `ValidationVerdict`. The validation block of the receipt is populated by δ's invocation, not γ's.
- A γ-preflight FAIL is information γ uses internally; it does not change the receipt's `validation` block, because that block is δ-owned.
- If γ omits the preflight entirely, the receipt is still valid for δ to consume. The preflight is permitted, not required.

The reason for asymmetric authority is the same as the reason α≠β within a cell (`COHERENCE-CELL.md` §β Independence as Contagion Firebreak). γ both authoring the receipt and authoritatively validating it would let γ's blind spots in receipt construction pass through validation unchallenged. The δ boundary is structurally outside γ's frame; that is what makes its validation discriminative rather than confirmatory.

This is the position the issue's AC3 asks for: γ preflight may exist, but γ preflight is explicitly non-authoritative, and the authoritative firing point is δ-boundary validation.

### Ordering rule

The parent scope cannot treat a closed cell as one of its inputs until the δ-boundary invocation has produced a verdict and δ has recorded a boundary decision. Operationally:

```text
γ closes cell                    →   closed_cell exists; receipt drafted
γ emits receipt                  →   receipt exists; cell is "closed", not "accepted"
δ invokes V on receipt           →   V emits ValidationVerdict
δ records BoundaryDecision       →   {accept | reject | repair-dispatch | override | release}
ACCEPTED                         →   receipt transmissible to parent scope
```

Three constraints follow from this ordering:

1. **δ invocation is mandatory before parent acceptance.** A receipt that exists on the cycle branch but has not been δ-validated is in the closure-emitted-but-not-accepted intermediate state. The parent scope does not see it as an input.
2. **γ preflight does not substitute.** No matter how many γ preflights pass, δ has not invoked `V` until δ explicitly does so, and the cell is not transmissible until then.
3. **The `BoundaryDecision` is the moment of acceptance, not the `ValidationVerdict` alone.** Even when `V` returns PASS, the cell is not transmissible until δ records the boundary decision — because δ holds gate authority over what crosses the boundary, including the authority to reject a PASS receipt for boundary-policy reasons orthogonal to `V`'s predicate (e.g., a release freeze, a transport problem, a downstream-scope readiness check). PASS makes acceptance available to δ; PASS does not make acceptance automatic.

### Why not pre-merge or "possibly twice"

`COHERENCE-CELL.md` §Open Questions framed Q1 as "Pre-merge (δ accept gates the merge) or post-merge (δ accept gates release/tag)? Possibly twice (γ preflight + δ authoritative)." The chosen position is **post-receipt-emission, pre-acceptance** — which is structurally adjacent to the merge boundary in the current CDD cycle (β merges into `main` after APPROVE; γ writes the closeout; δ's boundary actions follow). The placement is fixed relative to the receipt, not relative to the merge: the receipt is what `V` reads, so `V` must fire after the receipt exists.

The "possibly twice" framing is reconciled here as "preflight is permitted; only the second firing is authoritative." Naming both firings as equally authoritative would diffuse authority across γ and δ and reintroduce the trust-by-seniority mode the deep invariant rejects. The single authoritative firing point at the δ boundary is the position this design commits to.

### Consequence

Phases 2–7 inherit a fixed firing point. Schema design (Phase 2) can assume `V` reads a complete, γ-emitted receipt — not a partial intermediate object. The validator implementation (Phase 3) can be authored as a function invoked by δ, with γ-preflight as a separate non-authoritative call against the same implementation. The δ split (Phase 4) can locate `V` invocation explicitly inside the δ skill's boundary-action sequence. The `CDD.md` rewrite (Phase 7) can name the δ-boundary firing point as the authoritative validation moment.

---

## Q2 — Is V a capability or a command?

**Chosen position.** `V` is a **typed predicate exposed as a capability by the `cnos.cdd` package, invoked by δ from within the boundary-action sequence**. It is not modeled as a freshly-introduced operator command, and it is not modeled as a runtime provider with its own registry. The existing `cn-cdd-verify` command surface continues to exist as the operator-facing invocation path (per `COHERENCE-CELL.md` §cn-cdd-verify Target Position), but δ's role doctrine binds to the capability — not to the command — so that command and capability stay separate runtime surfaces rather than smearing into one.

### Rationale

The doctrine constraint is that `V` is `Contract × Receipt × EvidenceGraph → ValidationVerdict` — a pure predicate. The design question is which runtime surface that predicate is exposed through. `src/packages/cnos.core/skills/design/SKILL.md` §2.5 names four runtime surfaces that must not smear together: skills choose, commands dispatch, orchestrators execute, providers provide capability. The four options the predecessor doctrine names for `V` map cleanly onto three of those surfaces:

| Option (from `COHERENCE-CELL.md` §Open Question 2) | Runtime surface |
|---|---|
| Predicate invoked by δ in-skill | Capability (provided to δ-the-skill by the package) |
| Capability registered in package manifest | Capability |
| Separate command under `cnos.cdd/commands/` | Command |
| Provider registered in a registry | Provider |

The capability framing wins on three independent counts.

**1. Substitutability is wrong for command and provider framings.** Per `design/SKILL.md` §3.5 (interfaces belong to consumers), an interface is justified by real substitution pressure. `V` has exactly one canonical implementation per CDD doctrine version — the schemas in `receipt.cue`/`contract.cue` plus the validator that interprets them. There is no consumer-side variation pressure that would justify a provider registry (alternative validator implementations swapping behind a contract). Modeling `V` as a provider would invent substitution that does not exist. Modeling `V` as a registered command would build a runtime registry entry that exists only to be invoked by one caller (δ), with all the registry/help/doctor/status overhead that `design/SKILL.md` §2.7 names. A capability is the surface that matches the doctrine: one predicate, one canonical implementation, invoked by one caller in its boundary-action sequence.

**2. Policy stays above detail.** From `design/SKILL.md` §3.3, policy lives in the kernel/core; details live in packages. δ's boundary policy is "invoke validation before accepting the receipt; gate transmissibility on the verdict; treat override as degraded." That policy is δ-skill doctrine. The implementation — reading the receipt, parsing the contract, walking the evidence graph, computing failed predicates — is package detail. A capability boundary is what keeps those layered. A command boundary would let invocation mechanics (argv parsing, stdout shape, exit-code conventions) leak into δ's policy.

**3. The operator command and the capability are different consumers.** `cn-cdd-verify` exists today as an artifact-presence checker invoked by operators and by CI. Phase 3 refactors it into the implementation of `V`. After Phase 3, two callers exist: operators (and CI) running `cn cdd-verify {N}` to validate a cycle, and δ-the-skill invoking `V` during the boundary-action sequence. The operator-facing invocation is a command — exact-dispatch, argv-driven, stdout/exit-code surfaced. The δ-facing invocation is a capability — function-call shape, structured `ValidationVerdict` return, no argv. Both wrap the same underlying predicate implementation. Modeling `V` itself as a command would force δ to consume it through the operator-facing surface (argv + stdout parsing) when δ has no operator-facing reason to do so. Modeling `V` as a capability and `cn-cdd-verify` as the command keeps the two consumer paths separate per `design/SKILL.md` §3.7.

The reading: **`V` is the capability; `cn-cdd-verify` is the command that wraps the capability for operator/CI consumption.** Phase 3 ships both — refactors `cn-cdd-verify` into the wrapper, exposes the underlying predicate as a capability the `cnos.cdd` package provides, and updates δ-the-skill (Phase 4) to bind to the capability.

### Minimal invocation sketch

The invocation shape at doctrine level — illustrative, not schema-pinned:

```text
# δ-the-skill invokes V as a capability during boundary actions
verdict := cdd.validate_receipt(
    contract:       <ref to issue body + AC oracles>,
    receipt:        <ref to typed receipt artifact>,
    evidence_graph: <ref to cycle .cdd/ artifacts + diff>,
)

# verdict is a typed ValidationVerdict — see §Validation Interface
if verdict.result == PASS:
    δ records BoundaryDecision: accept
elif verdict.result == FAIL and δ chooses override:
    δ records BoundaryDecision: override (degraded) — see §Q4
else:
    δ records BoundaryDecision: reject | repair-dispatch
```

The capability surface name `cdd.validate_receipt` is illustrative; Phase 3 chooses the actual capability identifier and Phase 2 names the `ContractRef` / `ReceiptRef` / `EvidenceRootRef` types. The doctrine commits to the shape — function call from δ, structured return — not to the syntax.

For comparison, the operator path is the same predicate exposed through the command surface:

```text
$ cn cdd-verify 367
# wraps the capability; prints human-readable summary; sets exit code
```

Both paths reach the same underlying predicate. The two callers (δ-the-skill, operators+CI) consume different runtime surfaces over the same implementation.

### What the design does not commit to

The design does **not** commit to:

- The capability identifier (the Phase 3 implementation chooses)
- The package manifest format that registers the capability
- The transport mechanism δ-the-skill uses to invoke the capability (in-process call, MCP, subprocess — Phase 3 chooses)
- The argv/stdout shape of `cn-cdd-verify` post-refactor (Phase 3 chooses)

The design commits to: `V` is a capability, not a command and not a provider; δ binds to the capability, not to the command; the command continues to exist as a parallel operator-facing surface.

### Consequence

Phase 3 (`cn-cdd-verify` refactor) knows what it is building: a package-provided capability (the predicate) and an operator-facing command (the wrapper). Phase 4 (δ split) knows what δ-the-skill binds to: the capability, named in δ's role doctrine. Phase 6 (ε relocation) is not affected — ε's protocol-iteration role does not consume `V` directly.

---

## Q3 — Where does ε relocate?

**Chosen position.** ε relocates to **`ROLES.md`** — the generic role-scope ladder doctrine at the repo root. The CDD-specific instantiation (`src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`) is rewritten in Phase 6 as a thin pointer to the generic doctrine in `ROLES.md`, retaining only the CDD-instantiation details ε needs (where receipt streams live, what protocol gaps look like in CDD, what `cdd-iteration.md` is). **This cycle names the target; Phase 6 ships the move.**

### Rationale

The doctrine constraint from `COHERENCE-CELL.md` §ε as Protocol Evolution is that ε is outside ordinary cell metabolism. Inside the cell, α produces, β discriminates, γ closes, δ effects the boundary. ε observes the cell from outside, across many cells, and patches the protocol when receipt-stream patterns reveal a structural gap the protocol does not prevent. ε's matter is the protocol itself.

The doctrine consequence is that ε is generic to the role-scope ladder pattern, not specific to CDD. Any ladder instantiation that emits receipts can be observed by an ε that evolves the protocol. CDD is one such instantiation; if `ROLES.md` ever documents another instantiation (e.g., a research-cycle ladder, a product-cycle ladder), each instantiation has its own α/β/γ/δ doing in-cell metabolism and a shared ε that observes receipt streams across instantiations. ε's doctrine therefore belongs at the level where the ladder pattern itself is documented — `ROLES.md` at the repo root.

The three candidates the predecessor doctrine names map to three different authority scopes:

| Candidate | Authority scope | Reading |
|---|---|---|
| `ROLES.md` | Generic role-scope ladder | ε is a generic role function of the ladder pattern; instantiations inherit it |
| `cnos.core/doctrine/` | Core runtime substrate | ε is a core surface that crosses all packages |
| New `cnos.protocol-iteration` package | Dedicated subsystem | ε is a first-class subsystem with its own package boundary |

`ROLES.md` is chosen for three reasons.

**1. ε is a role function, not a runtime substrate.** Per `COHERENCE-CELL.md` §Structural Prediction, role function (α/β/γ/δ/ε) and runtime substrate are distinct surfaces — fusing them is the surface-smearing the structural prediction guards against. `cnos.core` is the runtime-substrate home (kernel, dispatch, harness primitives). Putting ε under `cnos.core/doctrine/` would smear role doctrine onto the runtime layer. ε is a role; `ROLES.md` is the doctrine home for the role-scope ladder.

**2. A new `cnos.protocol-iteration` package adds boundary without adding substitutability.** Per `design/SKILL.md` §3.10 (prefer package/install cohesion over topic labeling), a new package boundary is justified when there is a real install/use unit to separate. ε's matter is doctrine evolution — patches to skills, schemas, validators, harness contracts. That matter lands across `cnos.cdd`, `cnos.core`, and `cnos.eng` depending on where the protocol gap surfaced. There is no install/use unit "the protocol-iteration runtime" that a separate package would carry. Inventing one would create a package whose only purpose is to host one role's doctrine, which is the false-package-boundary failure `design/SKILL.md` §3.10 names.

**3. `ROLES.md` already owns the generic ladder.** `CDD.md` opens with the pointer to `ROLES.md`: "cdd is the reference instantiation of the generic role-scope ladder pattern. The pattern (α/β/γ/δ/ε roles, scope-escalation contract, instantiation fields) is documented at `ROLES.md` at the repo root." The α/β/γ/δ generic doctrine already lives there. ε's generic doctrine joining it preserves a single home for the ladder pattern and avoids splitting role doctrine across two surfaces.

### What Phase 6 will do

Phase 6 is not authored here; this cycle commits to the target. The expected Phase 6 work, named so subsequent cycles inherit a stable target shape:

1. **Lift ε generic doctrine into `ROLES.md`** — the protocol-evolution function, the "ε observes receipt streams across cells" framing, the "ε is not required in-cycle for cells with no protocol gaps" rule, and the `protocol_gap_count` / `protocol_gap_refs` artifact rule.
2. **Shrink `cnos.cdd/skills/cdd/epsilon/SKILL.md`** to a CDD-instantiation pointer — what receipt stream ε reads in CDD (the `.cdd/releases/{X.Y.Z}/{N}/` cycle records), what counts as a CDD protocol gap, where `cdd-iteration.md` lives when emitted.
3. **Remove ordinary-metabolism framing from CDD cycle requirements** — cycle artifacts no longer require an ε-produced `cdd-iteration.md` unconditionally; the receipt's `protocol_gap_count == 0` is the no-gap signal, and `cdd-iteration.md` is required only when `protocol_gap_count > 0` (per `COHERENCE-CELL.md` §ε Artifact Rule, already named in target doctrine).

The Phase 6 surface — `ROLES.md` plus a thinned `epsilon/SKILL.md` — is what this cycle commits to. The Phase 6 implementation is out of scope for #367 by construction.

### What this cycle does not do

This cycle does **not** edit `ROLES.md`. This cycle does **not** edit `epsilon/SKILL.md`. The AC9 surface-containment contract forbids both. Phase 6 is the cycle that ships the move; #367's job is to fix the target so Phase 6 has stable doctrine to refer back to. The relocation is named, not enacted.

### Consequence

Phase 6 has a target. The Phase 6 issue body can cite this section as its design surface, propose the `ROLES.md` patch, propose the `epsilon/SKILL.md` thin, and inherit the rationale for not choosing `cnos.core/doctrine/` or a new package. The doctrine surface (`COHERENCE-CELL.md` §ε as Protocol Evolution) remains the source of truth for what ε *is*; this design surface fixes *where ε lives*.

---
