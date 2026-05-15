<!-- sections: [Preamble, Q1, Q2, Q3, Q4, Q5, Validation Interface, Non-goals, Closure] -->
<!-- completed: [Preamble, Q1, Q2, Q3, Q4, Q5, Validation Interface, Non-goals, Closure] -->

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

## Q4 — How is an override receipted?

**Chosen position.** An override is a **structured `override:` block inside the receipt's `boundary` block**, populated by δ when δ records a `BoundaryDecision` of `override` against a non-PASS `ValidationVerdict`. The `validation` block of the receipt continues to carry `V`'s original verdict unchanged — override does not rewrite, mask, or replace the `ValidationVerdict`. The presence of a non-null `boundary.override` is the structural signal to every downstream consumer that the receipt closed in degraded state.

### Rationale

The doctrine constraint from `COHERENCE-CELL.md` §Trust Boundary and §δ Boundary Complex is that an override is **a degraded boundary action, not a form of validity**. δ may override a non-PASS verdict; δ may not pronounce a non-PASS verdict valid. The override exists to let δ effect a boundary action when the cell's matter is unavoidably degraded (production deadline, legacy debt, explicitly-tolerated divergence) without lying about the verdict.

The three candidate shapes the predecessor doctrine names trade off three properties:

| Candidate | Inline with receipt? | Single artifact? | Detectable by downstream V? |
|---|---|---|---|
| Structured `override:` block inside receipt | Yes | Yes | Yes — field presence |
| Sibling artifact `override.md` | No | No | Indirectly — requires consumer to check for a sibling file |
| Degraded-state flag with override metadata block | Yes | Yes | Yes — but the flag and the metadata are two surfaces |

The structured `override:` block wins on all three counts simultaneously. It keeps the override **inline with the receipt** (a receipt that closed under override is one artifact, not a pair to be carried together), **single-artifact** (the receipt is the parent-facing trust surface; degradation is part of it, not a separate document), and **detectable by downstream V** (the parent scope's `V` reads `boundary.override` from the receipt's evidence-graph input directly, with no need to crawl sibling files).

The sibling-artifact framing reintroduces the failure mode `design/SKILL.md` §3.2 names — two places claiming authority over the same fact ("did the cell close under override?"). One downstream consumer might read `override.md` and miss the `validation` block; another might read the `validation` block and miss `override.md`. The structured-block-in-receipt framing keeps the override grammar inside the parent-facing trust surface.

The degraded-state flag with separate metadata block reintroduces a similar split: the flag is one signal, the metadata is another, and downstream consumers can drift between checking the flag and checking the metadata. The structured `override:` block is its own signal — if the field is non-null, the override exists; if the field is null, no override; field presence is the detection rule.

### Required fields

The override block, when populated, carries the following fields. Names are illustrative; Phase 2 pins schema syntax.

| Field | Type | Purpose |
|---|---|---|
| `actor` | role-identity ref (e.g., `delta@cdd.cnos`) | δ actor who authored the override |
| `justification` | free text | Narrative reason for the override — what compelled δ to proceed despite a non-PASS verdict |
| `original_validation_verdict` | `ValidationVerdict` ref or inline copy | The `V`-emitted verdict the override is overriding. This is `validation.result == FAIL` at minimum; the field carries the structured copy or a ref so the override block stands self-contained |
| `failed_predicates_overridden` | list of predicate refs | The specific failed predicates from the original verdict the override covers. Naming each is required: a blanket "override all failures" is forbidden because it makes downstream consumers unable to reason about which specific failures were tolerated |
| `degraded_state` | boolean | Always `true` when override is non-null. Redundant with the field-presence detection rule, but explicit so a downstream consumer reading the field directly observes the degraded property without inferring it |
| `downstream_consumer_detection_rule` | rule statement | The detection rule itself, included for human readers: "any receipt with `boundary.override != null` is a degraded receipt; PASS-equivalence requires `validation.result == PASS` AND `boundary.override == null`." Phase 3 may elect to fold this into doctrine and omit the per-receipt copy; the doctrine commitment is that the rule exists and is documented somewhere downstream consumers can read |

The minimum guarantee: every override carries `actor`, `justification`, `original_validation_verdict`, `failed_predicates_overridden`, and `degraded_state`. Phase 2 may add fields (expiry, follow-up issue ref, risk note); the design freezes the floor, not the ceiling.

### What override never does

Three things are explicit failure modes if the design or implementation drifts:

1. **Override never rewrites the `ValidationVerdict`.** The `validation` block of the receipt continues to carry `V`'s emitted verdict. `validation.result` stays at FAIL; `validation.failures` stays populated. δ does not edit, redact, or replace `V`'s verdict. The verdict is what `V` saw; the override is what δ did with it. The two are separate fields with separate authors.
2. **Override never substitutes for PASS in downstream consumers.** A receipt where `validation.result == FAIL` and `boundary.override` is populated is **not equivalent** to a receipt where `validation.result == PASS` and `boundary.override == null`. Downstream `V` at scope n+1 — reading the n-receipt as evidence — must distinguish these two cases. Treating a degraded n-receipt as PASS-equivalent would let degradation propagate silently into the parent scope, which is exactly the failure mode the deep invariant rejects ("a closed cell is trusted because its receipt validates against its contract").
3. **Override never emits `OVERRIDE-PASS` as a `ValidationVerdict`.** `V` emits PASS, FAIL, and possibly WARN. There is no `OVERRIDE-PASS` value. Override lives in the `BoundaryDecision` shape (an enumerant of what δ decided), not in the `ValidationVerdict` shape (an enumerant of what `V` observed). Fusing them would let δ-side decisions reach back and rewrite `V`-side observations, collapsing the two surfaces the design works to keep separate.

### Downstream-consumer detection rule

The rule, stated precisely so the design and the implementation both inherit it:

```text
A receipt is degraded ⇔ boundary.override != null.
A receipt is PASS-equivalent ⇔ validation.result == PASS AND boundary.override == null.

Every downstream consumer of a receipt — including:
  - the parent scope's V (V at scope n+1)
  - the parent scope's δ
  - operators auditing the cycle
  - cn-cdd-verify when reading a closed cycle
— MUST check both validation.result and boundary.override before
treating the receipt as a clean closure. Reading only one of the
two fields is incorrect.
```

The biconditional form is load-bearing. A consumer that checks only `validation.result == PASS` and assumes clean closure is incorrect (it misses degraded cells that V did not approve). A consumer that checks only `boundary.override == null` and assumes clean closure is also incorrect (it misses cells where V emitted FAIL but no override exists — which is itself an incomplete-closure state per `COHERENCE-CELL.md` §Receipt Validity: "A missing override block plus a non-PASS V is not a receipt — it is an incomplete closure"). Both fields must be checked together.

The detection rule has a structural consequence: the design commits to **two independent verdict surfaces** (`ValidationVerdict` from `V`, `BoundaryDecision` from δ) rather than collapsing them. The override exists in the `BoundaryDecision` surface; the failed predicates exist in the `ValidationVerdict` surface; the rule above joins them. The next section (§Validation Interface) names both surfaces explicitly.

### Consequence

Phase 2 (schemas) inherits the override block's required fields and the detection rule. Phase 3 (validator implementation) inherits the rule that `V` never emits `OVERRIDE-PASS` and never rewrites its own verdict in response to an override decision. Phase 4 (δ split) inherits the structural prediction that δ's role doctrine names the override path as a degraded boundary action with mandatory receipted fields, not as a discretionary skip-validation lever.

---

## Q5 — Are role closeouts evidence-graph inputs or parent-facing artifacts?

**Chosen position.** The typed **receipt is the parent-facing artifact** of the closed cell. The five existing per-cycle markdown files (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`) **become evidence-graph inputs to `V`** and **receipt-derivation sources** — they continue to exist as the in-cell record, structurally unchanged in their current authoring shape, but the parent scope reads the receipt rather than crawling the five files. The receipt is **computed at γ close-out** from the five files plus the diff plus the issue body.

### Rationale

The doctrine constraint from `COHERENCE-CELL.md` §cn-cdd-verify Target Position is that "the parent-facing trust surface becomes the typed receipt: a single artifact that summarises the cell at the validator's interface" while "the existing per-cycle markdown artifacts continue to exist as evidence-graph inputs the validator reads when scoring." That doctrine fragment fixes the direction. This section names it precisely.

Two reasons make the typed receipt the right parent-facing surface and the markdown files the right evidence-graph inputs.

**1. Parent scopes read receipts, not transcripts.** The recursion equation in `COHERENCE-CELL.md` §Recursion Equation has the accepted receipt at cell n becoming α-level matter at cell n+1. The parent scope reasons over the cell as a typed object — fields, structured failures, explicit override blocks — not over the cell's internal narrative. The five markdown files are narrative artifacts: `beta-review.md` is a round-by-round review transcript; `gamma-closeout.md` is γ's triage narrative; `alpha-closeout.md` is α's cycle observations. Narrative artifacts are the right medium for cell-internal coordination (which they have been since CDD began) and the wrong medium for cross-scope reasoning (which is a typed-object job). Keeping the typed receipt above the markdown layer is what makes the recursion work; collapsing the receipt into "go read these five markdown files" would force every parent scope to re-derive structure from narrative each time.

**2. The five files are exactly what makes a coherent receipt computable.** The receipt's contract block needs the issue body and the AC oracles — `self-coherence.md` already maps these to evidence. The receipt's production block needs α actor and artifact list — `alpha-closeout.md` (and the cycle's git history) carries the α actor identity, and the diff carries the artifact list. The review block needs β verdict, findings, and merge SHA — `beta-review.md` carries the verdict and findings; `beta-closeout.md` carries the merge SHA and final acceptance state. The closure block needs γ triage, unresolved debt, and protocol gap signals — `gamma-closeout.md` carries all of these. The five files are not a coincidence; they correspond to the receipt blocks one-to-one (with α-closeout and β-closeout overlapping production / review respectively because the closeouts narrate what production and review did). Deprecating any of the five would remove a derivation source for one of the receipt's blocks. They stay.

### Per-file roles

The five files, with their named role under this design:

| File | Role under V | Receipt block(s) derived from it | Continues to exist? |
|---|---|---|---|
| `self-coherence.md` | Evidence-graph input + receipt-derivation source | `contract` (AC oracle results), `production` (artifact enumeration via CDD Trace step 6), `closure` (known debt rows) | Yes — unchanged in current shape |
| `beta-review.md` | Receipt-derivation source | `review.verdict`, `review.rounds`, `review.findings`, `review.findings_dispositions` | Yes — unchanged in current shape |
| `alpha-closeout.md` | Evidence-graph input | `production.alpha_actor` (already in git), α-side cycle observations contributing to `closure.protocol_gap_refs` when patterns surface | Yes — unchanged in current shape |
| `beta-closeout.md` | Receipt-derivation source | `review.merge_sha`, `review.beta_actor`, `boundary.accepted` (when β's acceptance is the closure point under current CDD; future cycles may relocate this to δ explicitly) | Yes — unchanged in current shape |
| `gamma-closeout.md` | Receipt-derivation source | `closure.gamma_actor`, `closure.triage_ref`, `closure.unresolved_debt`, `closure.next_commitments`, `closure.protocol_gap_count`, `closure.protocol_gap_refs` | Yes — unchanged in current shape |

None of the five files is **deprecated** by this design. The role-skill instantiation of CDD continues to require the same artifacts in the same shapes. What changes is what crosses the parent boundary: the receipt does; the five files do not (they remain available to V and to operators as the evidence graph, but they are not parent-scope inputs).

### Derivation rule

The receipt is computed at γ close-out, after all five files have landed on `main`:

```text
γ writes gamma-closeout.md (last of the five cycle artifacts to land)
γ invokes receipt-derivation:
    receipt := derive_receipt(
        contract       := <issue body + AC oracles + active design constraints>,
        diff           := <git diff origin/main..merge_sha>,
        self_coherence := <self-coherence.md>,
        beta_review    := <beta-review.md>,
        alpha_closeout := <alpha-closeout.md>,
        beta_closeout  := <beta-closeout.md>,
        gamma_closeout := <gamma-closeout.md>,
    )
γ emits receipt as the typed parent-facing artifact
γ optionally invokes V on the receipt as non-authoritative preflight (§Q1)
δ invokes V on the receipt; δ records BoundaryDecision (§Q1)
```

The derivation function `derive_receipt` is Phase 3 work; the doctrine commits to the existence of this function and to the input set (the five files plus the diff plus the contract). Phase 2 designs the receipt schema; Phase 3 implements the derivation. The contract this cycle freezes is: receipt fields map to evidence in specific files in a documented way; no receipt field is invented out of nothing; every receipt is reproducible from the cycle's `.cdd/` artifacts plus the diff.

The doctrine consequence is that the **receipt is computable, not authored**. `derive_receipt` is a function over evidence, not a free-text artifact γ writes by hand. γ's free-text artifact is `gamma-closeout.md`, which is one of the function's inputs. This separation is what keeps the receipt typed: a typed artifact computed from inputs has a checkable shape; a free-text artifact a role writes has only the shape the role chose at writing time.

### Why not the inverse reading

The Open Question's alternative reading was "closeouts remain parent-facing and the receipt is a thinner summary." That reading is rejected for two reasons.

First, it leaves the parent scope with five markdown files to read as inputs, which is the failure mode the receipt exists to fix. The recursion equation needs an α-input at cell n+1 — a typed object the next-scope α produces matter against. Five markdown files are not a typed input.

Second, it leaves the validator without a single artifact to validate. `V` is `Contract × Receipt × EvidenceGraph → ValidationVerdict`. If the receipt is a thin summary that does not carry the structured fields, `V` collapses into "read all five markdown files and infer structure." That defeats the purpose of having `V` as a typed predicate. The receipt is the parent-facing artifact precisely because `V` consumes it.

### Consequence

Phase 2 (schemas) inherits the receipt's structural commitments — which blocks exist, which fields they carry, which derivation sources feed them. Phase 3 (validator + `derive_receipt`) inherits the function shape and the input set. Phase 5 (γ shrink) inherits the constraint that `gamma-closeout.md` remains one of `derive_receipt`'s required inputs even after γ's role-skill shrinks; γ's coordination-function continues to produce the closeout, even if γ's runtime-supervision-mechanics responsibilities move to the harness.

---

## Validation Interface

This section freezes the parent-facing validation interface at doctrine level. **Schema syntax is Phase 2** — the shapes below are prose plus minimal illustrative examples, not `.cue` definitions. Phase 2 chooses the schema language and pins exact field types; this section commits to which contracts exist, what each carries, and how they compose.

### Input contract — what V reads

`V` reads three references:

```text
V : ContractRef × ReceiptRef × EvidenceRootRef → ValidationVerdict
```

- **`ContractRef`** points to the contract the cell was accountable to. The contract is the issue body (acceptance criteria, scope, non-goals, related artifacts, active design constraints) at the SHA committed to in the issue's `## Source of truth` table. The issue body is the authoring substrate; the AC oracles are the executable predicates. The contract is what `V` checks the receipt against; the cycle does not edit the contract.

- **`ReceiptRef`** points to the typed receipt artifact emitted by γ at close-out (§Q5). The receipt carries the `contract` / `production` / `review` / `closure` / `validation` / `boundary` blocks per the illustrative shape in `COHERENCE-CELL.md` §Illustrative Receipt Shape. The receipt is the artifact `V` operates on as primary input.

- **`EvidenceRootRef`** points to the evidence graph — the cycle's `.cdd/` artifacts (the five markdown files per §Q5), the git diff (`origin/main..merge_sha`), the CI run results, and any project-specific evidence sources the contract names. The evidence graph is what `V` consults to verify the receipt's claims. Receipt fields are claims; evidence graph fields are checkable facts.

The input contract is **reference-shaped**, not value-shaped: `V` reads refs that resolve to artifacts at specific SHAs, not inline copies. This matters for two reasons. First, evidence-graph artifacts can be large (review transcripts, full diffs); inlining them into `V`'s call would force a copy of the cycle's content into the call shape. Second, refs preserve reproducibility: the receipt + the refs let any operator re-run `V` against the same inputs the original δ invocation saw.

### Output contract — what V emits

`V` emits a `ValidationVerdict`:

```text
ValidationVerdict {
    verdict          : PASS | FAIL | (optionally) WARN
    failed_predicates: list of predicate references that did not hold
    warnings         : list of advisory observations
    provenance       : { validator_identity, validator_version, checked_at, input_refs }
}
```

- **`verdict`** is the headline value. `PASS` means every predicate in the contract held against the receipt+evidence. `FAIL` means at least one predicate did not hold. `WARN` (if Phase 2 elects to include it) names contract-defined advisory observations that do not block PASS — Phase 2 may collapse this into `warnings` and emit only `PASS | FAIL` as the verdict enum.

- **`failed_predicates`** is the structured list of predicates that did not hold, with refs back to the contract (which AC, which oracle) and refs into the evidence graph (which artifact/field the predicate read). The list is empty when `verdict == PASS`. The list is non-empty when `verdict == FAIL`. This is what makes `V` discriminative rather than ceremonial — a FAIL verdict carries the structured story of which predicates failed where.

- **`warnings`** carries advisory observations the contract permits without failing PASS — for example, a deferred sub-AC explicitly carried as known debt that the receipt acknowledges in `closure.unresolved_debt`. Phase 2 names the warning catalog; the doctrine commits to the field existing.

- **`provenance`** records who emitted the verdict, against which inputs, at what time. The validator identity is the capability's canonical name; the validator version is the deployed implementation version; `checked_at` is the timestamp of the invocation; `input_refs` echo the three refs above. Provenance lets a future operator re-run `V` deterministically against the same inputs and observe the same verdict — or, if the verdict differs, locate the validator-version drift that caused the divergence.

The output contract is **structured, not free text**. `ValidationVerdict` is a typed object Phase 2 pins; consumers (δ at scope n, `V` at scope n+1, operators, CI) read fields, not prose.

### Invocation contract — when and how V is invoked

The invocation contract is fixed by §Q1 (when) and §Q2 (how), reproduced here as the joint contract for completeness:

**When.** `V` fires authoritatively at the δ boundary, after γ emits the receipt and before δ records the `BoundaryDecision`. γ may invoke `V` as a non-authoritative preflight while drafting the receipt; γ preflight does not populate the receipt's `validation` block and does not authorize δ to skip its own invocation. The receipt's `validation` block is populated only by δ's invocation.

**How.** `V` is a typed predicate exposed as a capability by the `cnos.cdd` package, invoked by δ from within the boundary-action sequence. The operator-facing `cn-cdd-verify` command wraps the same predicate as a separate runtime surface; δ does not invoke `V` through the command shape.

**Composition with BoundaryDecision.** δ reads the `ValidationVerdict` returned by `V` and records a `BoundaryDecision`. The composition rule:

```text
ValidationVerdict.verdict == PASS
    → δ may record BoundaryDecision in {accept, release}
    → δ may also record BoundaryDecision in {reject, repair-dispatch} if boundary-policy
      reasons orthogonal to V demand it (release freeze, downstream readiness check)

ValidationVerdict.verdict == FAIL
    → δ may record BoundaryDecision in {reject, repair-dispatch}
    → δ may record BoundaryDecision == override only as a degraded boundary action,
      with the override block populated per §Q4. boundary.override != null is the
      structural signal that downstream consumers MUST detect.
```

`V` makes acceptance **available** to δ on PASS; `V` does not make acceptance **automatic**. δ holds final gate authority over what crosses the boundary, and the `BoundaryDecision` is what δ records after consulting the verdict.

### ValidationVerdict vs BoundaryDecision — the structural distinction

The two surfaces are independent and must not be fused. The distinction is the load-bearing claim this section commits to:

| Surface | Emitted by | Describes | Lives in receipt block |
|---|---|---|---|
| `ValidationVerdict` | `V` (the predicate) | Whether the receipt satisfies the contract+evidence predicates | `validation` |
| `BoundaryDecision` | δ (the boundary actor) | What δ decided to do with the receipt at the boundary | `boundary` |

The composition rule binds them but does not collapse them. Three explicit constraints follow:

1. **`V` does not record boundary decisions.** `V` emits a verdict — PASS or FAIL — and a structured failure list. `V` does not record "accepted" or "rejected" or "released." Those are δ's. The `validation` block of the receipt is `V`-owned; the `boundary` block is δ-owned.

2. **δ does not rewrite `ValidationVerdict`.** δ may decide to override a FAIL verdict (§Q4), but δ does not rewrite the verdict to PASS. The verdict stays FAIL; the override block carries the degraded boundary action. The two blocks are independent.

3. **PASS-equivalence is the conjunction.** A receipt is PASS-equivalent for downstream consumers iff `validation.verdict == PASS` AND `boundary.override == null` (§Q4). Either condition failing makes the receipt non-PASS-equivalent. Downstream consumers (including the parent scope's `V`) must check both.

This distinction is what keeps the trust grammar honest. If `V` could emit `OVERRIDE-PASS`, δ-side decisions would be reaching back into the validity surface. If δ's `BoundaryDecision` could be `validated`, δ would be conferring validity rather than observing it. Neither happens in the design this section freezes.

### Illustrative example (prose form)

The example below is **illustrative**, not a schema definition. Phase 2 pins the schema syntax (the shape this example sketches is `.cue`-compatible-ish but is not `.cue`). Reading the example: an invocation of `V` produces a `ValidationVerdict` that δ then composes with a `BoundaryDecision`.

Suppose δ invokes `V` for cycle #N at merge SHA `<merge-SHA>` with the issue body at the merge SHA's view of the contract:

```text
INPUT:
    contract_ref       = issue:#N @ contract_sha
    receipt_ref        = .cdd/releases/X.Y.Z/N/receipt.<format> @ merge_sha
    evidence_root_ref  = .cdd/releases/X.Y.Z/N/  @ merge_sha
                         + git_diff(origin/main..merge_sha)

V's WORK (illustrative — actual predicate set is contract-derived, not fixed here):
    - For each AC in the contract, check the oracle against the evidence graph
    - Check β-actor != α-actor (structural review-independence per
      COHERENCE-CELL.md §Review-Independence Evidence)
    - Check verdict-precedes-merge ancestry
    - Check artifact enumeration matches diff
    - ... etc, per Phase 3's predicate set

OUTPUT (case PASS):
    ValidationVerdict {
        verdict           : PASS
        failed_predicates : []
        warnings          : []
        provenance        : {
            validator_identity : cdd.validate_receipt
            validator_version  : <Phase 3 version>
            checked_at         : <timestamp>
            input_refs         : { contract_ref, receipt_ref, evidence_root_ref }
        }
    }

OUTPUT (case FAIL — illustrative shape, not a real failure):
    ValidationVerdict {
        verdict           : FAIL
        failed_predicates : [
            {
                predicate_ref : contract.AC4.oracle
                evidence_ref  : receipt.review.findings
                detail        : "AC4 oracle: review carries no finding-disposition rows"
            },
        ]
        warnings          : []
        provenance        : { ... as above ... }
    }

δ's COMPOSITION (PASS case):
    BoundaryDecision { kind: accept, accepted_at: <ts>, release_tag: <future> }
    Receipt.validation = ValidationVerdict (above)
    Receipt.boundary   = { delta_actor, accepted: true, override: null, ... }

δ's COMPOSITION (FAIL case, override path):
    BoundaryDecision { kind: override, ... }
    Receipt.validation = ValidationVerdict (FAIL, above) — UNCHANGED
    Receipt.boundary   = {
        delta_actor,
        accepted: true,
        override: {
            actor                          : delta@cdd.cnos
            justification                  : <text>
            original_validation_verdict    : ValidationVerdict (FAIL, above)
            failed_predicates_overridden   : [ contract.AC4.oracle ]
            degraded_state                 : true
            downstream_consumer_detection_rule : "boundary.override != null
                ⇒ receipt is degraded; PASS-equivalence requires
                validation.verdict == PASS AND boundary.override == null"
        }
    }
```

The example is **illustrative** in two senses. First, the shape (curly-brace block, indented fields) is prose-tabular, not a schema. Phase 2 chooses whether the schema language is CUE, JSON Schema, TypeScript-types, or something else. Second, the field names are illustrative — `validator_identity`, `failed_predicates_overridden`, etc. — and Phase 2 may rename any of them for consistency with its chosen schema's conventions. The doctrine commitment is the **shape**: three contracts (input, output, invocation), two independent verdict surfaces (`ValidationVerdict`, `BoundaryDecision`), and the composition rule that joins them.

### What Phase 2 chooses

Phase 2 (`receipt.cue` + `contract.cue` design) chooses:

- Schema language (CUE, per `COHERENCE-CELL.md` §Practical Landing Order item 2)
- Exact field names and field types
- Cardinality and optionality of fields
- Validation predicate language (how predicates in `failed_predicates` are referenced and serialized)
- Whether `WARN` is a third verdict enumerant or a `warnings`-list affair
- Whether `provenance.input_refs` includes the issue body SHA or a content hash

Phase 2 does **not** change:

- The three-contract shape (input, output, invocation)
- The `ValidationVerdict`-vs-`BoundaryDecision` distinction
- The PASS-equivalence biconditional rule (§Q4)
- The δ-authoritative firing point (§Q1)
- The capability-not-command framing (§Q2)
- The receipt-is-parent-facing rule (§Q5)

This is the line between design (this document) and schema (Phase 2). The design freezes the shape; Phase 2 pins the syntax.

### What Phase 3 inherits

Phase 3 (`cn-cdd-verify` refactor) inherits the function signature, the input refs, the output verdict shape, the invocation composition rule, and the override-detection rule. Phase 3 implements `V` and `derive_receipt`; it does not redesign the interface. If Phase 3 finds a missing field in the design — a real consumer pressure the doctrine surface did not anticipate — the right response is to amend this document in a follow-up cycle, not to silently extend the interface in implementation.

---

## Non-goals

This document is a design surface. Its non-goals are the surfaces it must not edit and the choices it must not make. The AC9 contract in [#367](https://github.com/usurobor/cnos/issues/367) is binding; this section restates the non-goals so they are internal to the document and observable by any reader.

**Files this document does not edit.** The diff for this cycle, before α signals review-readiness, is constrained to exactly:

- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` (new file, status `A`)
- `.cdd/unreleased/367/*.md` (cycle evidence — `gamma-scaffold.md`, `alpha-codex-prompt.md`, `beta-codex-prompt.md`, `self-coherence.md`, and any β/γ artifacts that land later in the cycle)

The cycle's surface-containment contract specifically forbids:

- Any `.cue` file (Phase 2 work)
- Any edit under `src/packages/cnos.cdd/commands/cdd-verify/` (Phase 3 work)
- Any edit to `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (Phase 4 work)
- Any edit to `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (Phase 5 work)
- Any edit to `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` (Phase 6 work)
- Any edit to `ROLES.md` (Phase 6 work, named target only)
- Any edit to `src/packages/cnos.cdd/skills/cdd/CDD.md` (Phase 7 work; explicitly out of scope for this cycle by the issue's `## Scope` section)
- Any edit to `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` (the predecessor doctrine surface — frozen by this cycle's contract)

**Choices this document does not make.** The design names the contracts; subsequent phases pin the syntax and the implementation:

- Schema language and exact field types for the receipt and contract (Phase 2)
- The capability identifier and package-manifest binding for `V` (Phase 3)
- The transport mechanism δ-the-skill uses to invoke the capability (Phase 3)
- The argv/stdout shape of `cn-cdd-verify` post-refactor (Phase 3)
- The exact `ROLES.md` text that lifts ε's generic doctrine (Phase 6)
- Any change to `CDD.md`'s canonical algorithm description (Phase 7)

**Author-discipline non-goals.** The design is intentionally **decisive over exhaustive**. The five Open Questions get single chosen answers; the rationale exists to justify the choice and to anchor future reading, not to enumerate every alternative the design could have made. The reader who wants the full alternative space can read `COHERENCE-CELL.md` §Open Questions; this document commits to one position per question and explains why.

---

## Closure

The design surface is complete when:

1. The five Open Questions seeded by `COHERENCE-CELL.md` §Open Questions have single chosen positions with rationale — `## Q1` through `## Q5` above.
2. The validation interface is frozen at doctrine level — `## Validation Interface` above, naming input contract, output contract, invocation contract, and the `ValidationVerdict`-vs-`BoundaryDecision` distinction.
3. The non-goals are restated and observable in the document — `## Non-goals` above, listing the surfaces this cycle does not touch.

Subsequent phases of [#366](https://github.com/usurobor/cnos/issues/366) inherit this document as their input contract:

| Phase | Inherits from this document |
|---|---|
| Phase 2 — `receipt.cue` + `contract.cue` | §Validation Interface (output shape, input refs), §Q4 (override fields), §Q5 (receipt blocks and derivation sources) |
| Phase 3 — `cn-cdd-verify` refactor | §Q1 (firing point), §Q2 (capability framing), §Validation Interface (function signature) |
| Phase 4 — δ split (`operator/` → `delta/` + harness) | §Q1 (δ owns the authoritative firing), §Q2 (δ binds to capability not command), §Q4 (override is a δ-owned receipt block) |
| Phase 5 — γ shrink | §Q1 (γ preflight is non-authoritative), §Q5 (`gamma-closeout.md` remains a derivation source even after γ shrinks) |
| Phase 6 — ε relocation | §Q3 (target is `ROLES.md`; `epsilon/SKILL.md` shrinks to CDD-instantiation pointer) |
| Phase 7 — `CDD.md` rewrite | §Q1 + §Validation Interface (the rewrite can cite the δ-boundary firing point and the typed receipt as the parent-facing trust surface) |

The doctrine surface (`COHERENCE-CELL.md`) remains the source of truth for *what the cell is*; this design surface fixes *how the cell's receipt is validated and where the validator's interface lives*. When Phase 3 lands and `V` becomes an executable predicate, this document moves from draft design to binding doctrine — at that point, the document's `Binding-status` line in the preamble is updated by the Phase 3 cycle's authoring (not by this cycle), and the design becomes the ratified contract Phases 4–7 satisfy.
