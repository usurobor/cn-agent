---
name: cds
description: Coherence-Driven Software. Use for engineering work — code, tests, branches, diffs, CI, releases, deployments — that requires structurally-enforced discipline against artifact degradation under unmeasured churn. Software-development realization of CCNF; sibling to cnos.cdr (research).
artifact_class: skill
kata_surface: embedded
governing_question: How do we improve software artifacts under repairable feedback coherently, without the engineering loop drifting farther from the kernel with each round?
visibility: public
triggers:
  - software
  - code
  - test
  - branch
  - diff
  - ci
  - release
  - deploy
  - cycle
  - engineering
scope: global
inputs:
  - software-engineering gap or unshipped feature
  - active role
  - cycle or wave context
outputs:
  - canonical CDS contract loaded (CDS.md)
  - active role overlay loaded (when role overlays land at Subs 3–5)
  - required sub-skills selected
requires:
  - CDS applies (software-class matter, not research-class)
calls:
  - CDS.md
  - alpha/SKILL.md
  - beta/SKILL.md
  - gamma/SKILL.md
  - delta/SKILL.md
  - epsilon/SKILL.md
---

# CDS

## Load order

This skill is the package-visible loader entrypoint for CDS. External dispatch enters through `cds` only — internal sub-skill triggers (`alpha`, `beta`, `gamma`, `delta`, `epsilon`) are advisory and are used only after `cds` and the active role overlay have been loaded.

**v0.1 status.** `CDS.md` is shipped in this directory by Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403) and is loaded as Step 1 below. The per-role overlays (`alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `delta/SKILL.md`, `epsilon/SKILL.md`) remain deliverables of Subs 3–5 of the same tracker; until those surfaces land, this loader names the forthcoming overlay files in `calls:` as **advisory targets**, mirroring the [`cnos.cdr` Sub 1 precedent](../../../cnos.cdr/skills/cdr/SKILL.md). Readers who need the surface mapping for migrating lifecycle/artifacts/review/gate content should consult [`docs/extraction-map.md`](../../docs/extraction-map.md): every software-lifecycle surface currently resident in [`cnos.cdd/skills/cdd/CDD.md` §"Software-specific realization — pending cds extraction"](../../../cnos.cdd/skills/cdd/CDD.md) is named with its planned CDS destination and the migration sub-issue.

When CDS applies (the matter is software-class — artifact improvement under repairable feedback — not research-class — claim transmission under uncertainty):

1. Load [`CDS.md`](CDS.md) in this directory as the canonical instantiation contract. CDS.md declares the six fields per [`ROLES.md §3`](../../../../../ROLES.md): matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule. It also records the architectural-choice inheritance from [cnos#388](https://github.com/usurobor/cnos/issues/388), the persona/protocol/project boundary per [`ROLES.md §4a`](../../../../../ROLES.md), and the empirical-anchor citation to `usurobor/cnos` itself.
2. Load the role overlay for the active role:
   - α: [`alpha/SKILL.md`](alpha/SKILL.md)
   - β: [`beta/SKILL.md`](beta/SKILL.md)
   - γ: [`gamma/SKILL.md`](gamma/SKILL.md)
   - δ: [`delta/SKILL.md`](delta/SKILL.md)
   - ε: [`epsilon/SKILL.md`](epsilon/SKILL.md)

   **Pending Subs 3–5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).** Pointer overlays delegating to the existing [`cnos.cdd/skills/cdd/<role>/SKILL.md`](../../../cnos.cdd/skills/cdd/) overlays are acceptable in v0.1; deep rewrites are deferred.
3. Load Tier 2 (process / role-local) and Tier 3 (lifecycle / kata) skills as directed by the role overlay and the cycle shape.
4. Read the project binding's `.cds/` directory (when CDS reaches v1 — current cnos practice still uses `.cdd/` artifacts for software cycles; the empirical-anchor doc planned as Sub 7 of [cnos#403](https://github.com/usurobor/cnos/issues/403) will map current practice to CDS's six-field shape).

## Rule

[`CDS.md`](CDS.md) is the only normative source for:

- the six instantiation-contract fields (`§"Six-field instantiation contract"`)
- the architectural-choice inheritance from [cnos#388](https://github.com/usurobor/cnos/issues/388) (`§"Architecture choice"`)
- the persona/protocol/project boundary (`§"Persona, Protocol, Project"`)
- the empirical-anchor citation to `usurobor/cnos` (`§"Empirical anchor"`)
- the actor-collapse floor (Field 6 — including the α=β prohibition for substantive software work and the β-α-collapse-on-δ permission for skill/docs-class cycles per the breadth-2026-05-12 wave manifest precedent)

Per-role overlays may add:

- evidence requirements specific to the role's matter type
- role-local gates and procedural detail
- checklists, examples, katas
- execution detail for the steps the role owns

Per-role overlays may **not** redefine:

- the six fields or their meanings
- the architectural-choice declaration
- the persona/protocol/project boundary
- the actor-collapse floor (α=β is always prohibited for software-engineering close-outs)
- closeout obligations already named in `CDS.md`

## Role overlays

Forthcoming per Subs 3–5 of [cnos#403](https://github.com/usurobor/cnos/issues/403). The five role positions and their intended responsibilities are:

- [`alpha/`](alpha/SKILL.md) — α role: code, test, doc, design-note matter production; software-engineering discipline against artifact degradation.
- [`beta/`](beta/SKILL.md) — β role: review against compilation/passing-tests, AC oracle, mechanical-vs-judgment boundary, configuration-floor, evidence-binding rule; verdict in the typed `#CDSReceipt`.
- [`gamma/`](gamma/SKILL.md) — γ role: cycle coordination, typed `#CDSReceipt` close-out, triage of `cds-*-gap` findings to ε.
- [`delta/`](delta/SKILL.md) — δ role: gap selection, cycle dispatch, gate-transition-shaped cadence, merge-mounted gate enforcement.
- [`epsilon/`](epsilon/SKILL.md) — ε role: receipt-stream review over software protocol gaps; protocol patches via the CDS-iteration artifact (CDS analogue of `cdd-iteration.md`). Generic ε doctrine inherited from [`ROLES.md §4b`](../../../../../ROLES.md) per the cycle/401 precedent.

Each role overlay is a CDS-specific extension of the corresponding `cnos.cdd/skills/cdd/<role>/SKILL.md` generic doctrine. The kernel grammar (role-cell shape, algorithm structure, independence rules, resumption protocol) is inherited by reference; the discipline profile and matter type diverge for the software-engineering loss function per `ROLES.md §4a.2`.

## Cross-protocol relationship

CDS shares the role-cell kernel with [`cnos.cdd`](../../../cnos.cdd/) (the generic algorithm) and is a sibling of [`cnos.cdr`](../../../cnos.cdr/) (the research realization). The kernel — `COHERENCE-CELL`, `COHERENCE-CELL-NORMAL-FORM`, role grammar, generic receipt-validation interface — lives in `cnos.cdd`. The CDS-specific procedures (matter type, review oracle, γ close-out, δ cadence, ε cadence, actor collapse) will live in `CDS.md` when Sub 2 lands. The option (a) split was decided at [cnos#388](https://github.com/usurobor/cnos/issues/388) for schemas and extended to skills at [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376); the v0.1 [`cnos.cdr`](../../../cnos.cdr/) bootstrap is the structural precedent this package mirrors.

**CDS is not CDR-with-different-words.** Engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*. The role names are shared; the loss functions diverge; the receipts and oracles diverge to enforce the divergent disciplines. Loading the cds skill is a class-routing act — the matter is software-class, not research-class.

CDD owns the **generic recursive kernel.** CDS owns the **software-specific lifecycle, artifact contract, review CLP, gate, closure, assessment, retro-packaging, and non-goals.** No generic CCNF kernel doctrine is restated in CDS — only pointers from CDS to CDD/CCNF. The pointer discipline is the same one [`cnos.cdr/skills/cdr/SKILL.md`](../../../cnos.cdr/skills/cdr/SKILL.md) uses for its kernel citations.

## Conflict rule

If this file and [`CDS.md`](CDS.md) disagree, `CDS.md` governs.

If a role overlay and `CDS.md` disagree on field meanings, the architectural-choice declaration, the persona/protocol/project boundary, or the actor-collapse floor, `CDS.md` governs.

If a role overlay adds execution detail, gates, or evidence requirements for a CDS field **without** changing that field's meaning, the role overlay governs for that local execution detail.

If `CDS.md` and `ROLES.md` disagree on the six-field shape or the role-cell grammar, `ROLES.md` governs (CDS instantiates; ROLES governs the grammar).

If `CDS.md` and `cnos.cdd` (CCNF kernel) disagree on the kernel grammar (role-cell shape, recursion modes, scope-lift projections), `cnos.cdd` governs (CDS realizes the kernel; the kernel does not adapt to its realizations).

## v0.1 caveat

`CDS.md` is shipped (Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403); cycle/407). The role-overlay files named in this loader's `calls:` frontmatter and Step 2 of "Load order" above remain **forthcoming** until Subs 3–5 land. The advisory nature is the cdr-Sub-1 precedent: a loader can name overlay surfaces before they exist as long as the cross-reference to the planned destination is unambiguous. See [`docs/extraction-map.md`](../../docs/extraction-map.md) for the surface-by-surface migration plan Subs 3–5 dispatch against.
