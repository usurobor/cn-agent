---
name: lifecycle
description: CDS development lifecycle — v0.1 thin overlay. Canonical 0–13 steps, S0–S13 state machine, branch rule, pre-flight, and tier structure at CDS.md; mechanics delegate to cnos.cdd role + harness skills.
artifact_class: skill
kata_surface: none
governing_question: How does a CDS cycle progress from observation through close-out across the 0–13 steps and S0–S13 state machine?
visibility: public
parent: cds
triggers:
  - lifecycle
scope: task-local
status: v0.1-thin-overlay
inputs:
  - a CDS cycle in progress
  - the current lifecycle step or state
outputs:
  - the canonical 0-13 step and S0-S13 state for the current cycle position (located in CDS.md)
---

# CDS Development lifecycle — v0.1 thin overlay

This v0.1 operational overlay delegates mechanics to existing cnos.cdd
role + harness skills until the v1 role rewrite. Canonical content at
[`../CDS.md §"Development lifecycle"`](../CDS.md).

## Canonical lifecycle shape

The 0–13 step table, the S0–S13 state machine, the `cycle/{N}` branch
rule, the γ-owned branch pre-flight, and the tier-1a/1b/1c/2/3 skill
loading structure are stated canonically at
[`../CDS.md §"Development lifecycle"`](../CDS.md). Cite that path.

## v0.1 operational realization

Per-role mechanics across each state — v0.1 overlay:

- [`../../../../cnos.cdd/skills/cdd/gamma/SKILL.md`](../../../../cnos.cdd/skills/cdd/gamma/SKILL.md) — γ (Steps 0–3, 9–10, 12–13).
- [`../../../../cnos.cdd/skills/cdd/alpha/SKILL.md`](../../../../cnos.cdd/skills/cdd/alpha/SKILL.md) — α (Steps 4–7, fix-round, close-out).
- [`../../../../cnos.cdd/skills/cdd/beta/SKILL.md`](../../../../cnos.cdd/skills/cdd/beta/SKILL.md) — β (Step 8 review + merge; pre-merge gate).
- [`../../../../cnos.cdd/skills/cdd/delta/SKILL.md`](../../../../cnos.cdd/skills/cdd/delta/SKILL.md) — δ (Step 11; disconnect boundary).
- [`../../../../cnos.cdd/skills/cdd/harness/SKILL.md`](../../../../cnos.cdd/skills/cdd/harness/SKILL.md) — runtime substrate (dispatch, polling, branch, identity).
- [`../../../../cnos.cdd/skills/cdd/release-effector/SKILL.md`](../../../../cnos.cdd/skills/cdd/release-effector/SKILL.md) — tag/release/deploy at S11.
- [`../../../../cnos.cdd/skills/cdd/operator/SKILL.md`](../../../../cnos.cdd/skills/cdd/operator/SKILL.md) — δ-the-operator session-routing.

When v1 lands, this overlay retires or expands into a full CDS-side skill.
