---
title: "From Task Execution to Decreasing Incoherence"
status: DRAFT
version: v0.1.0
date: 2026-05-22
proposed-path: docs/gamma/essays/DECREASING-INCOHERENCE.md
class: design essay
axis: gamma
related:
  - docs/THESIS.md
  - docs/alpha/essays/COHERENCE-SYSTEM.md
  - docs/alpha/essays/FOUNDATIONS.md
  - docs/essays/agent-first.md
  - docs/gamma/essays/CCNF-AND-TYPED-TRUST.md
  - src/packages/cnos.cdd/skills/cdd/CDD.md
  - src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
  - src/packages/cnos.cds/skills/cds/CDS.md
  - src/packages/cnos.cdr/skills/cdr/CDR.md
  - ROLES.md
  - schemas/cdd/
  - schemas/cds/
  - schemas/cdr/
  - usurobor/tsc:docs/THESIS.md
  - usurobor/tsc:runtime/SELF-MEASURE.md
---

# From Task Execution to Decreasing Incoherence

This document explains how cnos moves from agents completing tasks to agents following a measured path of decreasing incoherence.

## Executive summary

cnos exists because agent work fails in repeatable ways:

```text
it sounds done but is not done
it claims evidence it did not produce
it reviews itself
it forgets why a decision was made
it ships local fixes that break global coherence
it generates follow-up work without grounding it
it learns nothing from repeated failure
```

The answer is not a smarter prompt. The answer is a recursive, receipted, validated, and measured work loop:

```text
contract
  → CCNF cell
  → receipt
  → V validation
  → δ boundary decision
  → TSC measurement
  → ε finding stream
  → grounded issue proposal
  → next contract
```

The goal is not task completion. The goal is a durable path of decreasing incoherence.

## Governing question

How should cnos convert valid agent work into measured coherence deltas and grounded next issues?

This document answers that question at the vision/design level. It does not implement the validator, the orchestration grammar, or the issue-proposal schema.

## Current foundations

The current docs already establish the lower layers.

`docs/THESIS.md` defines cnos as a recurrent coherence system with Git as its lowest durable substrate. Its coherence-gradient section already states that every cycle should produce a coherence delta and every release should reduce a named incoherence.

`docs/alpha/essays/COHERENCE-SYSTEM.md` names the system instruction set:

```text
CMP → CAP (MCA/MCI) → CLP → update articulations → CMP
```

That is the whole-system loop: look, move, check, repeat.

`docs/alpha/essays/FOUNDATIONS.md` gives the theory/practice stack:

```text
C≡ → TSC → CTB → cnos
axiom → measure → execute → coordinate
```

`docs/essays/agent-first.md` names the activation problem: probabilistic work needs durable continuity, explicit authority, and receipt-bearing evidence. It argues that a receipt is the legible handoff substrate for agent composition.

`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` establishes the typed-trust layer: CCNF defines the cell algorithm, CUE validates typed trust surfaces, `V` validates boundary readiness, and δ records a boundary decision.

`src/packages/cnos.cdd/skills/cdd/CDD.md` now owns the generic CCNF kernel. `src/packages/cnos.cds/skills/cds/CDS.md` and `src/packages/cnos.cdr/skills/cdr/CDR.md` bind that kernel to software and research loss functions.

`ROLES.md` records the five-layer enforcement chain:

```text
persona
operator contract
protocol overlay
project binding
receipt / validator
```

It also records the discipline split:

```text
CDS: artifact improvement under repairable feedback
CDR: truth-preserving claim transmission under uncertainty
```

This document adds the next layer: coherence-driven issue generation.

## Problem

A task-completion system can report success while increasing incoherence.

The task may finish locally while the system loses global coherence. The agent may produce a plausible artifact, but the receipt may lack produced evidence. The branch may pass tests, but the package boundary may drift. The agent may generate a follow-up issue, but the issue may be grounded only in its own preference rather than in a validated receipt and measured bottleneck.

The failure is not only that agents sometimes make mistakes. The deeper failure is that the system has no mechanical way to ask:

```text
What incoherence did this work reduce?
What incoherence remains?
Which axis is now the bottleneck?
What bounded contract should be generated next?
```

Without that question, agent autonomy becomes a queue of locally plausible tasks.

## Design thesis

cnos should treat every accepted shippable cell as a measurement opportunity.

The cell receipt answers:

```text
Did this work close validly?
```

TSC answers:

```text
Did the target cohere more truthfully after the work?
```

ε answers:

```text
What repeated or newly exposed protocol gap should become the next bounded contract?
```

The future loop is:

```text
cellₙ closes
V validates receiptₙ
δ records decisionₙ
TSC measures targetₙ₊₁
ε reads receipt + measurement stream
issue_proposalₙ₊₁ is generated from the bottleneck
δ accepts, rejects, or dispatches the next contract
```

## Core model

### 1. CCNF closes work

The CCNF kernel remains the cell-level algorithm:

```text
matterₙ   := αₙ.produce(contractₙ)
reviewₙ   := βₙ.review(contractₙ, matterₙ)
receiptₙ  := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)
verdictₙ  := V(contractₙ, receiptₙ)
decisionₙ := δₙ.decide(receiptₙ, verdictₙ)
```

CCNF makes work transmissible. It does not, by itself, decide whether the system's coherence improved.

### 2. TSC measures the accepted target

After δ accepts or releases a cell, TSC measures the target bundle.

TSC axes are measurement axes, not role stations:

```text
α_TSC = pattern coherence
β_TSC = relational coherence
γ_TSC = process coherence
```

They are not the same as CDD/CDS/CDR roles:

```text
α_role = producer station
β_role = reviewer / discriminator station
γ_role = closer / coordinator station
```

The names rhyme because both are triadic. They must not be collapsed.

### 3. The bottleneck selects the investment axis

The post-cell TSC report names the weakest coherence axis or the highest leverage contributor to incoherence.

The issue generator maps that bottleneck to a next investment class:

```text
low α_TSC → pattern investment
low β_TSC → relation investment
low γ_TSC → process investment
```

Pattern investment stabilizes local form. Relation investment repairs cross-surface fit. Process investment repairs how the system changes.

### 4. ε turns measured gaps into issue proposals

ε reads receipt streams and coherence measurements. It does not invent arbitrary work.

An issue proposal must cite:

```text
accepted receipt
TSC report
bottleneck axis
finding or protocol-gap ref
bounded scope
acceptance criteria
proof plan
dispatch risk
```

A proposal without those links is not autonomous learning. It is another ungrounded task suggestion.

## Failure-to-mechanism map

| Failure | cnos mechanism |
|---|---|
| It sounds done but is not done. | typed receipt + `V` validation |
| It claims evidence it did not produce. | γ-bound evidence refs; `V` dereferences them |
| It reviews itself. | α ≠ β structural firebreak |
| It forgets why a decision was made. | `BoundaryDecision`, closeout, provenance |
| It ships local fixes that break global coherence. | TSC measurement over pattern, relation, process |
| It generates follow-up work without grounding it. | issue proposals cite receipt + TSC report + finding |
| It learns nothing from repeated failure. | ε reads receipt streams and patches protocol |

## Incoherence-decrease loop

The loop is not "do task, make next task." The loop is:

```text
1. Build the MCP for a target.
2. Dispatch a bounded CCNF cell.
3. Validate the receipt.
4. Accept, degrade, block, or repair.
5. Measure the accepted target with TSC.
6. Identify the bottleneck axis.
7. Generate a grounded issue proposal.
8. Validate the proposal.
9. Dispatch the next bounded cell or ask for human δ approval.
```

This is the practical form of:

```text
CMP → CAP → CLP → update articulations → CMP
```

The update is no longer only a code or docs patch. It includes the measured coherence delta and the next proposed contract.

## Issue proposal surface

A future issue proposal should be a typed artifact, not only prose.

Illustrative shape:

```yaml
schema: cnos.issue_proposal.v1

source:
  receipt_ref: .cdd/releases/docs/2026-05-22/407/receipt.yaml
  tsc_report_ref: .tsc/cds-407.json
  finding_ref: F2
  bottleneck_axis: beta

coherence:
  alpha: 0.86
  beta: 0.71
  gamma: 0.82
  investment_axis: beta
  rationale: >
    The accepted cell improved local CDS doctrine but left a relation gap
    between the CDS artifact contract and schemas/cds receipt fields.

proposal:
  title: Align CDS artifact contract with schemas/cds receipt fields
  domain: cds
  kind: relation-repair
  priority: P1
  risk: low

contract:
  in:
    - src/packages/cnos.cds/skills/cds/CDS.md
    - schemas/cds/
    - src/packages/cnos.cds/README.md
  out:
    - CDD.md rewrite
    - cn-cdd-verify implementation
  acceptance_criteria:
    - CDS.md names schemas/cds as the software evidence home.
    - schemas/cds README points back to CDS.md.
    - package README source map matches both.

dispatch:
  can_auto_dispatch: true
  requires_human_approval: false
```

This proposal is not trusted because an agent wrote it. It is trusted only if its sources validate and its scope is bounded.

## Investment classes

### Pattern investment

Pattern investment reduces local structural incoherence.

Use it when TSC α is the bottleneck.

Common cases:

```text
unstable names
duplicated concepts
ambiguous boundaries
schema shape not canonical
local definitions contradict themselves
```

Typical issue types:

```text
normalize terminology
compress duplicate doctrine
define missing primitive
split a mixed concept
stabilize schema enum or field shape
```

### Relation investment

Relation investment reduces cross-surface incoherence.

Use it when TSC β is the bottleneck.

Common cases:

```text
docs and schemas disagree
README and loader disagree
package boundary and runtime behavior diverge
authority pointers drift
same fact has multiple homes
```

Typical issue types:

```text
repair cross-references
align README / SKILL / schema / command surfaces
move content to the correct owner
add interface pointers
resolve source-of-truth conflict
```

### Process investment

Process investment reduces evolutionary incoherence.

Use it when TSC γ is the bottleneck.

Common cases:

```text
migration path unclear
roadmap state stale
generated vs canonical artifact confused
change procedure missing
repeated review failure pattern
```

Typical issue types:

```text
add migration plan
define deprecation rule
create roadmap issue
add receipt-stream rule
patch harness or validator
define phase gate
```

## Persona and protocol implications

This loop preserves the persona/protocol/project split.

Sigma may generate CDS proposals when software work exposes engineering incoherence. Rho may generate CDR proposals when research work exposes claim/evidence incoherence. Neither persona owns the protocol layer. They operate under protocol overlays.

The issue generator should route by discipline:

```text
software artifact gap → CDS / Sigma
research evidence gap → CDR / Rho
generic cell gap → CDD
orchestration gap → CCNF-O / handoff package
persona-learning gap → persona hub branch
project-binding gap → current project repo
```

A persona-learning proposal should not be stored as canonical doctrine in the project repo. The preferred transport is a branch in the persona hub, with the project repo recording only lineage if needed.

## CCNF-O relationship

CCNF-O formalizes orchestration grammar: cycles, waves, roadmaps, dependencies, gates, repair, joins, and dispatch artifacts.

This document is not CCNF-O. It defines the steering pressure that tells orchestration what to do next.

The relationship is:

```text
CCNF closes cells.
TSC measures coherence deltas.
Issue proposals name the next bounded contract.
CCNF-O composes those contracts into waves and roadmaps.
Harness executes the validated orchestration.
```

CCNF-O should not run ahead of this steering model. Otherwise the system may validate a beautiful orchestration graph that still follows the wrong path through coherence space.

## CUE and runtime relationship

CUE should validate the typed structures:

```text
receipt
boundary decision
TSC report reference shape
issue proposal
orchestration artifact
handoff bundle
```

CUE does not prove that the coherence report is semantically correct or that the next issue is wise. It rejects malformed or structurally incoherent artifacts before agents spend work on them.

The runtime and validators then check evidence:

```text
does the receipt exist?
does the TSC report exist?
does the report name the bottleneck claimed by the proposal?
does the proposal scope match the domain?
is auto-dispatch allowed by δ policy?
```

## Autonomy levels

### Level 1 — Human-authored contracts

Humans write issues. Agents implement and close receipts.

### Level 2 — Receipt-aware follow-up

Agents name follow-up findings inside closeouts. Humans choose which issues to file.

### Level 3 — Proposal autonomy

Agents generate typed issue proposals from receipts and TSC reports. Humans approve dispatch.

### Level 4 — Bounded dispatch autonomy

Low-risk proposals auto-dispatch when policy allows. High-risk proposals require human δ approval.

### Level 5 — Roadmap autonomy

The system maintains a roadmap of validated cells, repair paths, gates, and joins.

### Level 6 — Protocol evolution autonomy

ε observes repeated incoherence patterns and proposes protocol patches: skill edits, schema rules, validator predicates, harness rules, or package-boundary changes.

## Anti-gaming rules

A coherence-driven system must not optimize a number while degrading reality.

Reject or flag proposals that:

```text
increase coherence by deleting useful complexity
avoid needed refactors because they temporarily lower γ
hide evidence to make β look cleaner
split work into trivial cells only to inflate local pass rates
turn every low score into documentation churn
claim TSC improvement without comparing target-equivalent bundles
```

Coherence measurement is an observation, not a vanity metric. The next issue must still serve the project goal, respect risk policy, and fit a bounded contract.

## Design decision

Adopt the path-of-decreasing-incoherence frame.

cnos should describe future autonomy as:

```text
accepted receipt
  → measured coherence posture
  → bottleneck axis
  → grounded issue proposal
  → validated next contract
```

The system should not describe autonomy as merely "agent creates issues." Autonomous issue generation is valid only when the issue is grounded in a receipt, a measurement, and an explicit incoherence class.

## Migration plan

### Wave 0 — Land this essay

Add this file under `docs/gamma/essays/` and add a short pointer to `docs/gamma/essays/README.md`.

### Wave 1 — Define TSC report attachment for CCNF receipts

Name where post-cell TSC reports live, which targets get measured, and how reports cite receipts.

### Wave 2 — Define issue proposal schema

Add a CUE schema for grounded issue proposals.

The schema should require receipt ref, TSC report ref, bottleneck axis, proposal kind, bounded scope, acceptance criteria, and dispatch risk.

### Wave 3 — Add issue proposal validation

Teach the validator or a sibling command to reject ungrounded proposals.

### Wave 4 — Add ε proposal generation

Allow ε to produce issue proposals from receipt streams and TSC reports.

### Wave 5 — Add bounded auto-dispatch policy

Allow only low-risk, validated proposals to auto-dispatch.

### Wave 6 — Feed CCNF-O

Let CCNF-O compose validated issue proposals into waves and roadmaps.

## Success criteria

This direction succeeds when:

1. Every accepted shippable cell can be linked to a post-cell coherence measurement.
2. Every generated issue proposal cites the receipt and TSC report that caused it.
3. The proposal schema rejects ungrounded follow-up work.
4. ε can distinguish one-off defects from repeated protocol gaps.
5. δ can auto-dispatch low-risk proposals and route high-risk proposals to a human.
6. CCNF-O can compose validated proposals without hiding their source evidence.
7. The project can show a history of accepted cells whose remaining incoherence decreases over time.

## Open questions

1. Which targets should be measured after every cell: changed files only, package bundle, roadmap scope, or whole repo?
2. Should TSC reports live under `.tsc/`, the protocol evidence directory (`.cdd`, `.cds`, `.cdr`), or both?
3. What is the minimal mechanical TSC proxy that can run in CI before semantic scoring is available?
4. Which issue proposal kinds may auto-dispatch without human approval?
5. How should contradictory TSC reports be handled?
6. How should project-goal pressure combine with incoherence-reduction pressure?
7. Which layer owns the issue proposal schema: `cnos.core`, `cnos.cdd`, CCNF-O, or a future `cnos.coherence` package?

## Decision summary

cnos should not aim for agents that merely complete tasks. cnos should aim for agents that close valid cells, measure the resulting coherence posture, and generate the next bounded contract from the remaining incoherence.

The operating slogan:

```text
From task execution to decreasing incoherence.
```

Or in the full form:

```text
cnos turns fallible agent work into a measured path of decreasing incoherence.
```

## References

- `docs/THESIS.md`
- `docs/alpha/essays/COHERENCE-SYSTEM.md`
- `docs/alpha/essays/FOUNDATIONS.md`
- `docs/essays/agent-first.md`
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`
- `src/packages/cnos.cdd/skills/cdd/CDD.md`
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`
- `src/packages/cnos.cds/skills/cds/CDS.md`
- `src/packages/cnos.cdr/skills/cdr/CDR.md`
- `ROLES.md`
- `schemas/cdd/`
- `schemas/cds/`
- `schemas/cdr/`
- `usurobor/tsc:docs/THESIS.md`
- `usurobor/tsc:runtime/SELF-MEASURE.md`
