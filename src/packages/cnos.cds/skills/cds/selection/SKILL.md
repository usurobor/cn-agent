---
name: selection
description: CDS selection function — v0.1 thin overlay. Canonical rules at CDS.md; mechanics delegate to cnos.cdd gamma/SKILL.md.
artifact_class: skill
kata_surface: pointer
governing_question: How does γ select the next CDS cycle's gap by applying the canonical selection function in rule order?
visibility: public
parent: cds
triggers:
  - selection
scope: lifecycle-phase
status: v0.1-thin-overlay
---

# CDS Selection — v0.1 thin overlay

This v0.1 operational overlay delegates mechanics to existing cnos.cdd
gamma/alpha skills until the v1 role rewrite. Canonical rules live at
[`../CDS.md §"Selection function"`](../CDS.md).

## Canonical rules

The 10-rule selection function (P0 override, operational-infrastructure
override, assessment-commitment default, stale-backlog re-evaluation,
MCI freeze check, weakest-axis rule, maximum-leverage rule, dependency
order, effort-adjusted tie-break, no-gap case) plus the four
observation-surface inputs are stated canonically at
[`../CDS.md §"Selection function"`](../CDS.md). Cite that path.

## v0.1 operational realization

Mechanics γ runs (read inputs, build candidate table, apply rules in
order, name decisive clause, size intervention) — v0.1 overlay:

- [`../../../../cnos.cdd/skills/cdd/gamma/SKILL.md §2.1`](../../../../cnos.cdd/skills/cdd/gamma/SKILL.md) — observe + candidate table.
- [`../../../../cnos.cdd/skills/cdd/gamma/SKILL.md §2.2`](../../../../cnos.cdd/skills/cdd/gamma/SKILL.md) — select by rule order; size.
- [`../../../../cnos.cdd/skills/cdd/gamma/SKILL.md §2.2a`](../../../../cnos.cdd/skills/cdd/gamma/SKILL.md) — peer enumeration at scaffold.

When v1 lands, this overlay retires or expands into a full CDS-side skill.
