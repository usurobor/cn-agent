---
name: cdr
description: Coherence-Driven Research. Use for research work — investigation, synthesis, citation, dataset stewardship, claim transmission under uncertainty — that requires structurally-enforced discipline against overclaim.
artifact_class: skill
kata_surface: embedded
governing_question: How do we transmit claims about the world (or about a system under study) coherently under uncertainty, without the system believing something it did not measure?
visibility: public
triggers:
  - research
  - claim
  - hypothesis
  - dataset
  - analysis
  - field-report
  - wave
  - reproduce
  - cite
  - overclaim
scope: global
inputs:
  - research gap or unanswered question
  - active role
  - wave or project context
outputs:
  - canonical CDR contract loaded
  - active role overlay loaded (when shipped)
  - required sub-skills selected
requires:
  - CDR applies (research-class matter, not engineering-class)
  - CDR.md exists in this directory
calls:
  - CDR.md
  - alpha/SKILL.md
  - beta/SKILL.md
  - gamma/SKILL.md
  - delta/SKILL.md
  - epsilon/SKILL.md
---

# CDR

## Load order

This skill is the package-visible loader entrypoint for CDR. External dispatch enters through `cdr` only — internal sub-skill triggers (`alpha`, `beta`, `gamma`, `delta`, `epsilon`) are advisory and are used only after `cdr` and the active role overlay have been loaded.

When CDR applies (the matter is research-class — claim transmission under uncertainty — not engineering-class — artifact improvement under repairable feedback):

1. Load [`CDR.md`](CDR.md) in this directory as the canonical instantiation contract. CDR.md declares the six fields per `ROLES.md §3`: matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule.
2. Load the role overlay for the active role (forthcoming — Sub 3 of [cnos#376](https://github.com/usurobor/cnos/issues/376)):
   - α: `alpha/SKILL.md`
   - β: `beta/SKILL.md`
   - γ: `gamma/SKILL.md`
   - δ: `delta/SKILL.md`
   - ε: `epsilon/SKILL.md`
3. Load Tier 2 (process / role-local) and Tier 3 (lifecycle / kata) skills as directed by the role overlay and the wave shape.
4. Read the project binding's `.cdr/` directory (the canonical exemplar is [`usurobor/cph:.cdr/`](https://github.com/usurobor/cph)) for project-specific gates, datasets, policy.

## Rule

[`CDR.md`](CDR.md) is the only normative source for:

- the six instantiation-contract fields (`§"Six-field instantiation contract"`)
- the architectural-choice inheritance from cnos#388 (`§"Architecture choice"`)
- the persona/protocol/project boundary (`§"Persona, Protocol, Project"`)
- the empirical-anchor citation (`§"Empirical anchor"`)
- the actor-collapse floor (Field 6)

Per-role overlays (when shipped under Sub 3) may add:

- evidence requirements specific to the role's matter type
- role-local gates and procedural detail
- checklists, examples, katas
- execution detail for the steps the role owns

Per-role overlays may **not** redefine:

- the six fields or their meanings
- the architectural-choice declaration
- the persona/protocol/project boundary
- the actor-collapse floor (α=β is always prohibited for research claims)
- closeout obligations already named in `CDR.md`

## Role overlays (forthcoming — Sub 3 of cnos#376)

- `alpha/` — α role: research-claim, hypothesis, method, dataset, analysis, report production
- `beta/` — β role: review against falsifiability + diagnostic-oracle + reproduction-from-clean + citation-integrity + data-policy-compliance + claim/evidence-alignment
- `gamma/` — γ role: wave coordination, typed `#CDRReceipt` close-out
- `delta/` — δ role: gap selection, wave dispatch, gate-transition-shaped cadence
- `epsilon/` — ε role: receipt-stream review over protocol gaps; protocol patches

Until Sub 3 ships these files, CDR runs against `CDR.md` alone with role discipline supplied by the persona hub (e.g. [`cn-rho/spec/`](https://github.com/usurobor/cn-rho)) and the project binding (e.g. [`usurobor/cph:.cdr/`](https://github.com/usurobor/cph)).

## Cross-protocol relationship

CDR shares the role-cell kernel with [`cnos.cdd`](../../../cnos.cdd/) (engineering); the discipline profile diverges. The kernel — `COHERENCE-CELL`, role grammar, generic receipt-validation interface — lives in `cnos.cdd`. The CDR-specific procedures (matter type, review oracle, γ close-out, δ cadence, ε cadence, actor collapse) live here. The option (a) split was decided at [cnos#388](https://github.com/usurobor/cnos/issues/388) for schemas and extended to skills at [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376); see [`CDR.md` §"Architecture choice"](CDR.md) for the inheritance record.

**CDR is not CDS-with-different-words.** Engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*. The role names are shared; the loss functions diverge; the receipts and oracles diverge to enforce the divergent disciplines. Loading the cdr skill is a class-routing act — the matter is research-class, not engineering-class.

## Conflict rule

If this file and [`CDR.md`](CDR.md) disagree, `CDR.md` governs.

If a role overlay (when shipped) and `CDR.md` disagree on field meanings, the architectural-choice declaration, the persona/protocol/project boundary, or the actor-collapse floor, `CDR.md` governs.

If a role overlay adds execution detail, gates, or evidence requirements for a CDR field **without** changing that field's meaning, the role overlay governs for that local execution detail.

If `CDR.md` and `ROLES.md` disagree on the six-field shape or the role-cell grammar, `ROLES.md` governs (CDR instantiates; ROLES governs the grammar).
