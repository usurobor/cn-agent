---
title: "CCNF and Typed Trust: CDD as a Recursive Protocol Kernel"
status: DRAFT
version: v0.1.0
date: 2026-05-21
proposed-path: docs/gamma/essays/CCNF-AND-TYPED-TRUST.md
related:
  - docs/gamma/ENGINEERING-LEVELS.md
  - docs/gamma/essays/EXECUTABLE-SKILLS.md
  - ROLES.md
  - src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md
  - src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
  - src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md
  - schemas/cdd/README.md
  - schemas/cdd/*.cue
  - cnos#366
  - cnos#376
---

# CCNF and Typed Trust

## CDD as a Recursive Protocol Kernel

This document explains how cnos should move from prose-trusted CDD to CCNF-described, CUE-validated, domain-extensible trust surfaces.

The direction is L7 work: it changes the system boundary so future work becomes simpler, safer, cheaper, or unnecessary. CDD stops being a software-process encyclopedia. CDD becomes the recursive coherence-cell kernel. CDS binds that kernel to software engineering. CDR binds that kernel to research. CUE validates the typed receipts, contracts, and boundary decisions that cross trust boundaries.

## Executive summary

CDD should be written in **Coherence-Cell Normal Form**.

```text
matterₙ   := αₙ.produce(contractₙ)
reviewₙ   := βₙ.review(contractₙ, matterₙ)
receiptₙ  := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)
verdictₙ  := V(contractₙ, receiptₙ)
decisionₙ := δₙ.decide(receiptₙ, verdictₙ)
```

The parent scope does not trust a cell because γ says it closed. The parent scope trusts a cell only when the receipt validates against the contract and δ records a boundary decision.

```text
closed (αₙ, βₙ, γₙ) cell → αₙ₊₁ matter
δₙ boundary decision     → βₙ₊₁-like discrimination
εₙ receipt stream        → γₙ₊₁-like coordination/evolution
```

CUE should validate the typed trust surfaces, not every runtime fact. CUE checks the shape and structural constraints of contracts, receipts, validation verdicts, boundary decisions, transmissibility, protocol-gap signals, and domain evidence references. Runtime validators such as `cn-cdd-verify` wrap CUE and add evidence checks that CUE alone cannot prove.

The target is:

```text
CCNF defines the algorithm.
CUE validates the typed trust surfaces.
V validates boundary readiness.
δ decides after V emits a verdict.
CDS and CDR bind the generic cell to concrete domains.
```

## Problem

CDD has accumulated two jobs.

First, it is the generic coherence discipline: α produces, β reviews, γ closes, V validates, δ decides, and ε learns across receipt streams.

Second, it is the software-development procedure: GitHub issues, code diffs, tests, CI, release tags, role closeouts, and verifier checks over software artifacts.

Those two jobs now need different homes. If CDD keeps both, the generic kernel will inherit software-specific fields. CDR research cycles will then have to fake software receipts, or `cn-cdd-verify` will become a CDS validator while pretending to be generic CDD.

The same risk appears at the agent level. Sigma is an engineering persona. Rho is a research persona. CDD/CDS/CDR are protocols, not personas. If the system stores research discipline inside Sigma or software discipline inside CDD, the wrong loss function leaks into the wrong work.

Engineering can ship a bounded artifact and repair it under tests. Research must avoid transmitting unsupported claims. Those are sibling disciplines, not one mode switch.

## Current state

The following surfaces already exist.

`ROLES.md` defines the generic α/β/γ/δ/ε role-scope ladder and the five-layer enforcement chain:

```text
persona
operator contract
protocol overlay
project binding
receipt / validator
```

`COHERENCE-CELL.md` states the cell doctrine: roles are cell functions, runtime is substrate, validation is the cell wall / receptor assay, and release is boundary effection.

`COHERENCE-CELL-NORMAL-FORM.md` names the kernel algorithm. It defines CCNF as five steps, four closed-cell outcomes, two recursion modes, three scope-lift projections, and a two-layer separation between kernel and realization.

`RECEIPT-VALIDATION.md` freezes the parent-facing validator surface. `V` fires authoritatively at the δ boundary after γ emits a receipt and before parent acceptance or release. `V` emits `ValidationVerdict`; δ records `BoundaryDecision`.

`schemas/cdd/` contains draft CUE schemas for contracts, receipts, and boundary decisions. Those schemas are not yet binding runtime law. They become binding when Phase 3 implements `V` against them.

`cnos#366` tracks the roadmap from doctrine to executable protocol. Phase 2.5 exists because Phase 2 shipped schemas with software-shaped evidence refs. Phase 2.5 must separate generic CDD shape from CDS/CDR domain evidence before Phase 3 implements the validator.

`cnos#376` tracks the CDR package. It treats CDR as a protocol overlay, not as a persona hub and not as a project binding. Its empirical anchor is `usurobor/cph`, where local CDR-like research process already exists.

## Goals

1. Express CDD in CCNF.
2. Make typed receipts the parent-facing trust surface.
3. Use CUE to validate the structural invariants of contracts, receipts, verdicts, decisions, transmissibility, and protocol-gap fields.
4. Keep runtime evidence checking in `V` / `cn-cdd-verify`, not in CUE alone.
5. Split generic CDD from domain overlays:
   - CDD: recursive cell kernel;
   - CDS: software development overlay;
   - CDR: research overlay.
6. Keep persona hubs separate from protocols:
   - Sigma may operate CDS cycles;
   - Rho may operate CDR cycles;
   - neither persona is the protocol itself.
7. Preserve a migration path that does not rewrite `CDD.md` before V is executable and domain evidence has somewhere else to live.

## Non-goals

This document does not implement CUE schemas, rewrite `cn-cdd-verify`, split δ, shrink γ, relocate ε, create `cnos.cds`, create `cnos.cdr`, or rewrite `CDD.md`.

It does not claim CUE validates everything. CUE validates typed structure and constraint logic. Runtime validators still check evidence existence, command provenance, CI state, data mounts, script outputs, and domain-specific facts.

It does not collapse personas into protocols. Sigma, Rho, CDS, and CDR remain separate surfaces.

## Design principles

### 1. Trust moves by typed receipt, not role seniority

A closed cell is not trusted because a higher role approved it. A closed cell is trusted because its receipt validates against its contract and δ records a boundary decision.

### 2. Keep the kernel substrate-independent

The CCNF kernel must not name GitHub, CI, CUE, `cn-cdd-verify`, `claude -p`, data mounts, or release tags. Those are realization details. The kernel names only roles, contracts, matter, review, receipt, evidence, verdict, decision, outcome, and scope-lift.

### 3. Put domain evidence below the generic kernel

Generic CDD must not require `diff_ref`, `ci_refs`, `self_coherence_ref`, or `field_report_ref` by name. Generic CDD should require typed evidence refs. CDS and CDR should define which evidence refs their domains require.

### 4. Validation and decision are different surfaces

`V` emits `ValidationVerdict`. δ records `BoundaryDecision`. A PASS verdict enables ordinary acceptance. A non-PASS verdict may be overridden only as degraded boundary action. Override never rewrites the original verdict.

### 5. Personas carry discipline; protocols carry roles

Sigma is the engineering persona. Rho is the research persona. CDS and CDR are protocol overlays. A project such as cph binds CDR to concrete data, scripts, gates, and reports.

## Proposed design

### 1. Write CDD in CCNF

The final CDD algorithm should be short enough to quote:

```text
matterₙ   := αₙ.produce(contractₙ)
reviewₙ   := βₙ.review(contractₙ, matterₙ)
receiptₙ  := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)
verdictₙ  := V(contractₙ, receiptₙ)
decisionₙ := δₙ.decide(receiptₙ, verdictₙ)
```

The outcomes are:

```text
accepted := verdictₙ = PASS ∧ decisionₙ ∈ {accept, release}
degraded := verdictₙ ≠ PASS ∧ decisionₙ = override
blocked  := decisionₙ ∈ {reject, repair_dispatch}
invalid  := inconsistent(verdictₙ, decisionₙ)
```

The recursion modes are:

```text
repair_dispatch → same-scope recursion
accept/release  → cross-scope projection
override        → degraded cross-scope projection
reject          → no projection
invalid         → non-terminal; repair the receipt/protocol
```

The scope-lift is:

```text
closed (αₙ, βₙ, γₙ) cell → αₙ₊₁ matter
δₙ boundary decision     → βₙ₊₁-like discrimination
εₙ receipt stream        → γₙ₊₁-like coordination/evolution
```

This is not flat role renaming. δ is not β inside the child cell. ε is not γ inside the child cell. The mapping happens when a closed cell crosses scope.

### 2. Use CUE for typed trust surfaces

CUE schemas should type the parent-facing trust surfaces:

```text
#Contract
#Receipt
#ValidationVerdict
#BoundaryDecision
#Override
#Transmissibility
#ProtocolGapRef
```

CUE should enforce structural constraints such as:

```text
PASS + accept/release       → accepted
PASS + override             → invalid
non-PASS + override         → degraded
non-PASS + accept/release   → invalid
missing boundary_decision   → invalid receipt
γ-preflight-only            → invalid receipt
```

CUE should also enforce required protocol-gap signals:

```yaml
protocol_gap_count: 0
protocol_gap_refs: []
```

A no-gap receipt records absence explicitly. A gap-bearing receipt records structured refs that ε can read across the receipt stream.

### 3. Let V wrap CUE and evidence checks

CUE validates structure. `V` validates boundary readiness.

```text
V : Contract × Receipt → ValidationVerdict
```

The receipt carries evidence refs. `V` dereferences those refs and checks the evidence required by the declared protocol. δ does not inspect raw evidence directly; δ reads the receipt and verdict, then records a boundary decision.

For CDS, V checks software evidence: diffs, tests, CI, closeouts, release refs, and review-independence evidence.

For CDR, V checks research evidence: claim refs, data refs, method refs, result refs, data-mounted gates, reproduction status, claim status, limitations, and reviewer rerun evidence.

### 4. Split generic CDD from CDS and CDR

The package model should become:

```text
cnos.cdd
  Generic recursive coherence-cell kernel.
  Owns CCNF, generic contract/receipt/boundary-decision concepts,
  V semantics, scope-lift, and ε receipt-stream signal.

cnos.cds
  Software-development realization.
  Owns code/tests/docs matter, diffs, CI, releases,
  software closeouts, and software evidence rules.

cnos.cdr
  Research realization.
  Owns claims, hypotheses, data, methods, reports,
  reproducibility gates, claim/evidence alignment, and research verdicts.
```

The immediate migration should not force a large move. CDS can first extract by reference: name which current CDD surfaces are software-specific, then move them only after V works.

### 5. Separate persona, protocol, and project

The discipline stack is:

```text
persona          who this agent is
operator         what this agent refuses to do
protocol overlay what counts as valid work in this domain
project binding  which concrete gates and artifacts apply here
receipt/V        what evidence makes the work transmissible
```

Sigma and Rho should be distinct persona hubs because their loss functions differ.

```text
Sigma / CDS:
  artifact improvement under repairable feedback

Rho / CDR:
  truth-preserving claim transmission under uncertainty
```

Sigma may ship a bounded improvement when tests and receipts make the debt visible. Rho must not transmit an empirical claim without data, method, provenance, and claim calibration.

## Runtime flow

The target runtime flow is:

```text
contract created
  ↓
α produces matter
  ↓
β reviews matter
  ↓
γ emits typed receipt with evidence refs
  ↓
CUE checks receipt shape and structural constraints
  ↓
V dereferences receipt-bound evidence refs and emits ValidationVerdict
  ↓
δ records BoundaryDecision
  ↓
receipt becomes accepted, degraded, blocked, or invalid
  ↓
ε reads receipt streams for protocol evolution
```

CUE participates at the structural trust boundary. It does not replace runtime checks. The validator must still prove that referenced evidence exists and means what the receipt claims it means.

## Schema direction

Phase 2 shipped `schemas/cdd/`. Phase 2.5 should prevent those schemas from freezing software-specific evidence into generic CDD.

Preferred shape:

```text
schemas/cdd/
  generic cell contract, receipt, boundary decision
  evidence refs as protocol-extensible typed members

schemas/cds/
  software receipt extension
  requires diff refs, CI refs, software closeouts, release refs

schemas/cdr/
  research receipt extension
  requires claim refs, data refs, method refs, result refs,
  claim status, limitations, reproduction refs
```

An adapter-boundary variant is also acceptable if it preserves the same invariant:

```text
Generic CDD validates the cell shape.
Domain schemas validate domain evidence.
V dispatches by declared protocol.
```

## Alternatives considered

### Keep CDD as the software process

This preserves short-term convenience but makes CDR fake software receipts or fork the process. It also prevents final `CDD.md` from becoming compact.

### Put everything into CUE

CUE can validate structure and constraints, but not all runtime truth. It cannot by itself prove that CI ran, data was mounted, or a number came from a script. Those checks belong in V.

### One persona switches between CDS and CDR

Mode-switching saves one hub but mixes loss functions. Engineering rewards bounded shipping under repair. Research punishes overclaiming under uncertainty. The dispatch boundary should pick Sigma or Rho; one persona should not silently switch discipline.

### Create CDR before separating generic/domain schema shape

This would give CDR a package sooner, but CDR would inherit unstable or software-shaped validation. The safer order is Phase 2.5, then V, then CDR package hardening.

## Risks and mitigations

### Risk: generic CDD absorbs software-specific fields

Mitigation: Phase 2.5 splits generic schema from CDS/CDR evidence requirements before Phase 3 implements V.

### Risk: CUE gets overstated as full runtime truth

Mitigation: every design surface states the distinction: CUE validates typed trust surfaces; V checks evidence and provenance.

### Risk: CDD.md is rewritten too early

Mitigation: Phase 7 gates on V being executable and on CDS/CDR evidence having homes outside generic CDD.

### Risk: new packages add process weight

Mitigation: extract by reference first. Do not move files until the boundary proves itself.

### Risk: persona/protocol boundaries drift

Mitigation: require every protocol package to declare whether a statement belongs to persona, operator contract, protocol overlay, project binding, or receipt/validator.

## Migration plan

### Wave 0 — Design essay

Land this document under `docs/gamma/essays/` and point to it from the essays README.

Purpose: give the whole direction one citable design surface.

### Wave 1 — Phase 2.5 generic/domain schema boundary

Split or adapt `schemas/cdd/` so generic CDD does not require CDS-specific evidence fields.

### Wave 2 — Phase 3 validator

Implement V as the runtime trust gate:

```text
V : Contract × Receipt → ValidationVerdict
```

`cn-cdd-verify` becomes the operator-facing wrapper around V.

### Wave 3 — CDS extraction by reference

Create `cnos.cds` as the software overlay. First name the software-specific surfaces; move them only after the references stabilize.

### Wave 4 — CDR bootstrap

Create `cnos.cdr` as the research overlay, anchored in cph but not defined by cph. CDR owns research role overlays. Rho owns the research persona. cph owns the project binding.

### Wave 5 — δ split

Split δ into boundary policy, harness/platform driver, and release effector. Keep `delta/SKILL.md` small.

### Wave 6 — γ shrink and ε upscope

Move runtime supervision mechanics below γ. Move ε to protocol-iteration doctrine above ordinary CDD metabolism.

### Wave 7 — final CDD.md rewrite

Rewrite `CDD.md` around CCNF once V is executable and domain evidence has homes in CDS/CDR.

## Success criteria

This direction succeeds when:

1. `CDD.md` can state the algorithm in CCNF without hiding software assumptions.
2. `cn-cdd-verify` validates typed receipts through V.
3. CUE rejects structurally invalid receipts before δ accepts them.
4. CDS validates software evidence without making CDD software-specific.
5. CDR validates research evidence without copying cph-local prose into protocol doctrine.
6. Sigma and Rho operate different disciplines while sharing the same role-cell kernel.
7. ε can read protocol-gap signals across receipt streams without parsing role narratives.

## Open questions

1. Should Phase 2.5 prefer split schemas (`schemas/cdd`, `schemas/cds`, `schemas/cdr`) or a generic schema with domain adapters?
2. Should CDS exist before CDR, or can CDR bootstrap against generic CDD plus cph evidence while CDS extracts later?
3. Where should ε generic doctrine live: `ROLES.md`, `cnos.core`, or a new protocol-iteration package?
4. What is the minimal CDR fixture that proves V can validate research receipts without software-shaped evidence?
5. How much of the current five-file CDD closeout set remains as CDS evidence after typed receipts become parent-facing?

## Decision summary

Adopt CCNF as the compact language for CDD.

Use CUE to validate typed trust surfaces.

Implement V as the runtime validator over contract and receipt.

Separate generic CDD from CDS and CDR domain overlays.

Keep personas separate from protocols.

Do not finalize `CDD.md` until V works and domain evidence has somewhere else to live.

## References

- `ROLES.md`
- `docs/gamma/ENGINEERING-LEVELS.md`
- `docs/gamma/essays/EXECUTABLE-SKILLS.md`
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md`
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`
- `schemas/cdd/README.md`
- `cnos#366`
- `cnos#376`
