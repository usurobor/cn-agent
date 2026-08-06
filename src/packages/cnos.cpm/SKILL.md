---
name: cpm
description: cnos.cpm package entrypoint — Coherent Product Management. Use when structuring a body of work into owned workstreams and tracks, deciding a single source of truth per initiative, or reasoning about planning nomenclature and ownership.
artifact_class: skill
kata_surface: none
governing_question: What owns product-management doctrine — planning hierarchy, canonical sources of truth, single-threaded ownership — and how does it realize as issue labels and the board?
visibility: public
triggers:
  - workstream
  - track
  - planning hierarchy
  - single-threaded owner
  - canonical document
  - extreme clarity
  - roadmap
  - prioritization
scope: global
calls:
  - skills/cpm/planning-hierarchy/SKILL.md
---

# cnos.cpm

## Core Principle

`cnos.cpm` is the package home for **coherent product management**: how a body of work is structured, owned, and made legible, distinct from how it is dispatched (`cnos.cds`), developed (`cnos.cdd`), or labeled and rendered (`cnos.issues`).

`cnos.cpm` owns:

- **planning hierarchy** — the `workstream ⊃ track ⊃ issue` structure, single-threaded ownership per level, and one canonical source of truth per initiative (`skills/cpm/planning-hierarchy/SKILL.md`).

## Boundary

`cnos.cpm` owns the **doctrine**; it does not own the label mechanics or the board. The `workstream/*` and `track/*` labels, their definitions, and the treemap render are owned by `cnos.issues` (canonical [`docs/development/issues/TAXONOMY.md`](../../../docs/development/issues/TAXONOMY.md), `cn issues map`). The planning-hierarchy skill **cites** that realization; it does not fork it. If doctrine and taxonomy ever disagree, extend `TAXONOMY.md` to match the doctrine, then let the board consume it.

## Source

The planning hierarchy adapts **Naomi Gleit's "Canonical Everything" / "Extreme Clarity"** program-management practice (workstreams with single-threaded owners, each with its own canonical doc, recursively). This package names the source; it does not restate her essays.
