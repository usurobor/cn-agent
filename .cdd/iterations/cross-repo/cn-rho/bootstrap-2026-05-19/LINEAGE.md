# Bootstrap drafts — `usurobor/cn-rho` (research persona hub)

## Purpose

Stage `spec/PERSONA.md` and `spec/OPERATOR.md` content for the `cn-rho` research persona hub, ahead of the operator scaffolding the `usurobor/cn-rho` repository. The drafts here are the **content**; the **repo** is operator-scaffolded (high-blast-radius action — separate operator-side step).

## Source

- **Trigger:** cnos session conversation 2026-05-19; operator delegated drafting per the persona/protocol/project boundary (`ROLES.md §4a`)
- **Reference doctrine:** `ROLES.md §4a` (persona / operator / protocol / project / receipt enforcement chain), §4a.1 (discipline-profile required section), §4a.2 (loss-function distinction)
- **Sibling exemplar:** `usurobor/cn-sigma` carries the same two-file persona-hub shape (`spec/PERSONA.md`, `spec/OPERATOR.md`)
- **Empirical anchor:** `usurobor/cph#32` (the data-mounted-gate failure pattern that motivates Rho's refusal conditions)

## Bundle contents

- `PERSONA.md` — Rho's identity, voice, default station, memory discipline, `## Discipline` section per `ROLES.md §4a.1`
- `OPERATOR.md` — Rho's dispatch conditions, refusal conditions, boundary authority, and authority exclusions
- `LINEAGE.md` — this file

## Operator-side scaffolding steps

When the operator is ready to bring `usurobor/cn-rho` online:

1. Create `usurobor/cn-rho` repository (operator action; cnos session cannot create cross-repo via MCP).
2. Copy `PERSONA.md` to `cn-rho:spec/PERSONA.md`.
3. Copy `OPERATOR.md` to `cn-rho:spec/OPERATOR.md`.
4. Add a minimal `cn-rho:README.md` (router pattern from cnos#379 when that lands, or the cn-sigma README shape in the interim).
5. Add `cn-rho:threads/` and `cn-rho:state/` directories (mirrors cn-sigma's layout).
6. First Rho activation: load `cnos.core` skills (KERNEL, CA skills), this `PERSONA.md`, this `OPERATOR.md`. Once `cnos.cdr` v0.1 ships (cnos#376), Rho additionally loads `cnos.cdr/skills/cdr/SKILL.md` for protocol overlay.

## What is NOT in scope here

- Creating the cn-rho repository (operator action).
- Authoring `cnos.cdr` protocol skills (cnos#376 wave).
- Tightening `cph#32` (separate feedback patch at `cnos:.cdd/iterations/cross-repo/cph/issue-32-tightening-2026-05-19/`).
- Adding a `## Discipline` section to `cn-sigma/spec/PERSONA.md` (separate feedback patch at `cnos:.cdd/iterations/cross-repo/cn-sigma/discipline-section-2026-05-19/`).

## Disposition

Status: **drafted** (operator-pending — repo creation is the gate).

When the cn-rho repo lands and these drafts are applied, this bundle may be archived. The persona identity is preserved in the live cn-rho hub; the lineage is preserved in this LINEAGE.md.

Filed by γ on 2026-05-19.
