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
  - active role overlay loaded
  - required sub-skills selected
requires:
  - CDR applies (research-class matter, not engineering-class)
  - CDR.md exists in this directory
calls:
  - CDR.md
  - alpha/SKILL.md
  - beta/SKILL.md
  - gamma/SKILL.md
  - operator/SKILL.md
  - epsilon/SKILL.md
---

# CDR

## Load order

This skill is the package-visible loader entrypoint for CDR. External dispatch enters through `cdr` only — internal sub-skill triggers (`alpha`, `beta`, `gamma`, `operator`, `epsilon`) are advisory and are used only after `cdr` and the active role overlay have been loaded.

When CDR applies (the matter is research-class — claim transmission under uncertainty — not engineering-class — artifact improvement under repairable feedback):

1. Load [`CDR.md`](CDR.md) in this directory as the canonical instantiation contract. CDR.md declares the six fields per `ROLES.md §3`: matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule.
2. Load the role overlay for the active role:
   - α: [`alpha/SKILL.md`](alpha/SKILL.md)
   - β: [`beta/SKILL.md`](beta/SKILL.md)
   - γ: [`gamma/SKILL.md`](gamma/SKILL.md)
   - δ: [`operator/SKILL.md`](operator/SKILL.md)
   - ε: [`epsilon/SKILL.md`](epsilon/SKILL.md)
3. Load Tier 2 (process / role-local) and Tier 3 (lifecycle / kata) skills as directed by the role overlay and the wave shape.
4. Read the project binding's `.cdr/` directory (the canonical exemplar is [`usurobor/cph:.cdr/`](https://github.com/usurobor/cph)) for project-specific gates, datasets, policy.

## Rule

[`CDR.md`](CDR.md) is the only normative source for:

- the six instantiation-contract fields (`§"Six-field instantiation contract"`)
- the architectural-choice inheritance from cnos#388 (`§"Architecture choice"`)
- the persona/protocol/project boundary (`§"Persona, Protocol, Project"`)
- the empirical-anchor citation (`§"Empirical anchor"`)
- the actor-collapse floor (Field 6)

Per-role overlays may add:

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

## Role overlays

- [`alpha/`](alpha/SKILL.md) — α role: research-claim, hypothesis, method, dataset, analysis, report production; discipline against overclaim.
- [`beta/`](beta/SKILL.md) — β role: review against falsifiability + diagnostic-oracle + reproduction-from-clean + citation-integrity + data-policy-compliance + claim/evidence-alignment; verdict in the typed `#CDRReceipt`.
- [`gamma/`](gamma/SKILL.md) — γ role: wave coordination, typed `#CDRReceipt` close-out, triage of protocol-gap findings to ε.
- [`operator/`](operator/SKILL.md) — δ role: gap selection, wave dispatch, gate-transition-shaped cadence, data-mounted gate enforcement. (The directory is named `operator/` to mirror the engineering doctrine's `cnos.cdd/skills/cdd/operator/SKILL.md` exemplar; the role itself is δ per `CDR.md` Field 4.)
- [`epsilon/`](epsilon/SKILL.md) — ε role: receipt-stream review over protocol gaps; protocol patches via the CDR-iteration artifact.

Each role overlay is a CDR-specific extension of the corresponding `cnos.cdd/skills/cdd/<role>/SKILL.md` generic doctrine. The kernel grammar (role-cell shape, algorithm structure, independence rules, resumption protocol) is inherited by reference; the discipline profile and matter type diverge for the research loss function per `ROLES.md §4a.2`.

## Cross-protocol relationship

CDR shares the role-cell kernel with [`cnos.cdd`](../../../cnos.cdd/) (engineering); the discipline profile diverges. The kernel — `COHERENCE-CELL`, role grammar, generic receipt-validation interface — lives in `cnos.cdd`. The CDR-specific procedures (matter type, review oracle, γ close-out, δ cadence, ε cadence, actor collapse) live here. The option (a) split was decided at [cnos#388](https://github.com/usurobor/cnos/issues/388) for schemas and extended to skills at [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376); see [`CDR.md` §"Architecture choice"](CDR.md) for the inheritance record.

**CDR is not CDS-with-different-words.** Engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*. The role names are shared; the loss functions diverge; the receipts and oracles diverge to enforce the divergent disciplines. Loading the cdr skill is a class-routing act — the matter is research-class, not engineering-class.

## Conflict rule

If this file and [`CDR.md`](CDR.md) disagree, `CDR.md` governs.

If a role overlay and `CDR.md` disagree on field meanings, the architectural-choice declaration, the persona/protocol/project boundary, or the actor-collapse floor, `CDR.md` governs.

If a role overlay adds execution detail, gates, or evidence requirements for a CDR field **without** changing that field's meaning, the role overlay governs for that local execution detail.

If `CDR.md` and `ROLES.md` disagree on the six-field shape or the role-cell grammar, `ROLES.md` governs (CDR instantiates; ROLES governs the grammar).
