---
name: handoff
description: Inter-agent / inter-activation / inter-repo handoff. Use when work, state, or messages cross a boundary between two cells, two roles, two activations of the same role, two repos, or two coherence-cell-shaped protocol packages. Owns the wire format; consumers (cnos.cdd, cnos.cds, cnos.cdr, future c-d-X) own the binding.
artifact_class: skill
kata_surface: embedded
governing_question: How do two parts of the system that are not co-located in time, role, repo, or activation transfer state coherently, without inventing the protocol per cycle?
visibility: public
triggers:
  - handoff
  - dispatch
  - cross-repo
  - bundle
  - proposal
  - clarification
  - artifact-channel
  - receipt-stream
  - inter-activation
  - inter-agent
scope: global
inputs:
  - handoff event (intra-cycle / inter-cycle / inter-activation / inter-repo)
  - active role (γ for cross-repo; δ for dispatch enrichment; ε for receipt-stream)
  - consumer protocol package (cnos.cdd / cnos.cds / cnos.cdr / future c-d-X)
outputs:
  - canonical handoff contract loaded (HANDOFF.md)
  - active sub-skill loaded (cross-repo / dispatch / mid-flight / artifact-channel / receipt-stream)
  - required Tier 3 skills selected
requires:
  - handoff applies (the matter crosses a cell / role / activation / repo boundary)
calls:
  - HANDOFF.md
  - cross-repo/SKILL.md
  - dispatch/SKILL.md
  - mid-flight/SKILL.md
  - artifact-channel/SKILL.md
  - receipt-stream/SKILL.md
---

# Handoff

## Load order

This skill is the package-visible loader entrypoint for `cnos.handoff`. External dispatch enters through `handoff` only — internal sub-skill triggers (`cross-repo`, `dispatch`, `mid-flight`, `artifact-channel`, `receipt-stream`) are advisory and are used only after `handoff` and the active sub-skill have been loaded.

**v0.1 status.** `HANDOFF.md` and the five per-surface sub-skill files named in this loader's `calls:` frontmatter are **landed and canonical**. v0.1 is complete; the cnos#404 wave closed on 2026-05-22 via cnos#420. See [`HANDOFF.md`](HANDOFF.md) for the package contract and [`docs/extraction-map.md`](../../docs/extraction-map.md) for the per-surface source pin and migration provenance.

When handoff applies (the matter crosses a cell / role / activation / repo boundary; not "I am writing matter for my own scope at my own time on my own branch"):

1. Load [`HANDOFF.md`](HANDOFF.md) in this directory as the canonical wire-format contract. HANDOFF.md declares the package contract (what cnos.handoff owns; what it does not own) and indexes the five Landed sub-surfaces.
2. Load the sub-skill for the active handoff event:
   - cross-repo (proposal intake, outbound iteration, bilateral iteration, operator-pending): [`cross-repo/SKILL.md`](cross-repo/SKILL.md)
   - dispatch (operator → α/β/γ/δ prompt + implementation-contract enrichment): [`dispatch/SKILL.md`](dispatch/SKILL.md)
   - mid-flight (γ → in-flight α re-pinning; `gamma-clarification.md` mechanism): [`mid-flight/SKILL.md`](mid-flight/SKILL.md)
   - artifact-channel (intra-cycle α→β→γ sequential handoff over `.cdd/unreleased/{N}/`): [`artifact-channel/SKILL.md`](artifact-channel/SKILL.md)
   - receipt-stream (cross-cycle `cdd-iteration.md` + INDEX.md aggregator; ε reader feed): [`receipt-stream/SKILL.md`](receipt-stream/SKILL.md)

   All five v0.1 sub-skills are canonical at their `cnos.handoff` homes — landed via cnos#416 (cross-repo), cnos#417 (dispatch), cnos#418 (mid-flight + artifact-channel), cnos#419 (receipt-stream), with the cross-reference sweep closed via cnos#420. See [`docs/extraction-map.md`](../../docs/extraction-map.md) for the per-surface source pin.

3. Load Tier 2 (process / role-local) and Tier 3 (lifecycle / kata) skills as directed by the sub-skill and the event shape.

4. Read the consumer protocol package's loader for any binding-specific overrides (e.g., CDS's `.cdd/unreleased/{N}/` rooting; CDR's wave-shaped `.cdr/` rooting; CDD's generic instantiation rules).

## Rule

[`HANDOFF.md`](HANDOFF.md) is the only normative source for:

- the directional case taxonomy (intra-cycle / inter-cycle / inter-activation / inter-repo)
- the wire-format artifacts per case (bundle file sets; dispatch prompts; clarification files; artifact-channel rules; receipt-stream records)
- the typed schemas (if any are lifted from Markdown into `schemas/handoff/`)
- the consumer protocol (how a CDD/CDS/CDR/future-c-d-X package binds the wire format)

Per-sub-skill files may add:

- per-case operational detail (e.g., specific bundle file sets per directional case)
- per-event evidence requirements
- checklists, examples, katas
- execution detail for the steps the sub-skill owns

Per-sub-skill files may **not** redefine:

- the directional case taxonomy
- the wire-format artifact families
- the consumer-protocol binding contract
- the boundary against CCNF-O (the orchestration grammar — `cnos.handoff` transports work; CCNF-O composes cells; see "Cross-protocol relationship" below)

## Sub-skills

All five v0.1 sub-skills are Landed; the [cnos#404](https://github.com/usurobor/cnos/issues/404) handoff extraction wave closed on 2026-05-22 via [cnos#420](https://github.com/usurobor/cnos/issues/420):

- [`cross-repo/`](cross-repo/SKILL.md) — cross-repo state machine (8-event STATUS lifecycle: drafted/submitted/accepted/modified/rejected/landed/withdrawn/revised); bundle file sets per directional case (a/b/c/d); LINEAGE.md schemas per case; feedback-patch format; bundle archival rule; hat-collapse attribution. **Migrated wholesale from `cnos.cdd/skills/cdd/cross-repo/SKILL.md` by Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — landed under [cnos#416](https://github.com/usurobor/cnos/issues/416).**
- [`dispatch/`](dispatch/SKILL.md) — dispatch-prompt template (operator → α/β/γ/δ); implementation-contract 7-axes schema; δ-as-inward-membrane enrichment doctrine; γ template ↔ δ enrichment ↔ α constraint ↔ β verification four-surface mesh. **Migrated from `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` + `cnos.cdd/skills/cdd/operator/SKILL.md §3a` + `cnos.cdd/skills/cdd/delta/SKILL.md §2` by Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — landed under [cnos#417](https://github.com/usurobor/cnos/issues/417).**
- [`mid-flight/`](mid-flight/SKILL.md) — `gamma-clarification.md` mechanism for intra-cycle re-pinning of axis values; γ-to-in-flight-α channel; cache-bust semantics; spec-staleness propagation. **Authored from cnos#391 empirical anchor by Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — landed under [cnos#418](https://github.com/usurobor/cnos/issues/418).**
- [`artifact-channel/`](artifact-channel/SKILL.md) — intra-cycle α→β→γ sequential handoff rules over `.cdd/unreleased/{N}/`; per-role artifact write ownership; frozen-snapshot rule on merge; release-time directory move. **Migrated from `cnos.cds/skills/cds/CDS.md §"Artifact contract"` wire-format invariants by Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — landed under [cnos#418](https://github.com/usurobor/cnos/issues/418).** Per-artifact filename set, ownership matrix, and 13-stage Ordered flow remain canonical at cnos.cds (consumer-protocol-specific); the wire-format pattern (channel-shape, sequential rule, frozen-snapshot, release-time move) lives here.
- [`receipt-stream/`](receipt-stream/SKILL.md) — cross-cycle `cdd-iteration.md` receipt-stream; per-finding shape (F1, F2, …); `.cdd/iterations/INDEX.md` aggregator; ε reader feed; cross-repo trace bundle on `patch-landed` with `Cross-repo` disposition. **Migrated from `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` by Sub 5 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — landed under [cnos#419](https://github.com/usurobor/cnos/issues/419).**

Each sub-skill is a handoff-specific extension of the wire-format contract in `HANDOFF.md`. The kernel grammar (cell-cell composition, scope-lift projections) is inherited from `cnos.cdd` by reference; the wire format itself is the substrate-independent piece this package owns.

## Cross-protocol relationship

`cnos.handoff` owns **wire formats and handoff doctrine for inter-agent / inter-activation / inter-repo transport.** Three current consumers — [`cnos.cdd`](../../../cnos.cdd/) (which uses handoff doctrine natively for its own coordination layer), [`cnos.cdr`](../../../cnos.cdr/) (which has loaded `cross-repo/SKILL.md` via cross-package reference since v0.1 per [cnos#376](https://github.com/usurobor/cnos/issues/376)), and [`cnos.cds`](../../../cnos.cds/) (which absorbed some surfaces during the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave and continues to cite others from `cnos.cdd`) — and an open class of future c-d-X protocol packages that will consume the same primitives once the extraction lands.

**Boundary vs CCNF-O (Track A of [cnos#405](https://github.com/usurobor/cnos/issues/405)):** `cnos.handoff` transports work between cells; CCNF-O composes cells into cycles, waves, roadmaps, gates, joins. The two have **different reasons to change**: handoff changes when the wire format changes; CCNF-O changes when the composition grammar changes. Each has one reason to change. Sub 1 of [cnos#404](https://github.com/usurobor/cnos/issues/404) preserves this boundary by explicitly not preempting CCNF-O surfaces — no orchestration-composition grammar appears in this package; Track A of [cnos#405](https://github.com/usurobor/cnos/issues/405) is gated on the `cnos#404` wave closing.

**Boundary vs the domain protocols (CDD / CDS / CDR / future c-d-X):** handoff is the *wire format*; the domain protocols are the *cycle lifecycles*. CDD changes when the generic kernel changes; CDS changes when the software-engineering cycle shape changes; CDR changes when the research-cycle shape changes. None of those changes ought to ripple into handoff doctrine, and vice versa.

**`cnos.handoff` is not CDD-with-different-words.** CDD owns the cycle lifecycle and the CCNF kernel; `cnos.handoff` owns the wire format by which cycles talk to each other and the wire format by which roles within a cycle talk to each other. Loading the `handoff` skill is a class-routing act — the matter is a boundary-crossing event (intra-cycle, inter-cycle, inter-activation, or inter-repo); not "I am writing α-matter for my own scope at my own time on my own branch" (that's a domain-protocol activity, not a handoff event).

## Conflict rule

If this file and [`HANDOFF.md`](HANDOFF.md) disagree, `HANDOFF.md` governs.

If a sub-skill and `HANDOFF.md` disagree on the directional case taxonomy, the wire-format artifact families, or the consumer-protocol binding contract, `HANDOFF.md` governs.

If a sub-skill adds execution detail, evidence requirements, or operational gates for a wire-format case **without** changing that case's taxonomy or artifact family, the sub-skill governs for that local execution detail.

If `HANDOFF.md` and a consumer protocol (CDD / CDS / CDR / future c-d-X) disagree on the wire format for a handoff event, `HANDOFF.md` governs (consumers bind the wire format; the wire format does not adapt to its consumers). The consumer may add binding-specific overrides (e.g., directory rooting; filename conventions) but must do so in its own loader, not by re-stating handoff doctrine.

If `HANDOFF.md` and `cnos.cdd` (CCNF kernel) disagree on the cell-cell composition or scope-lift projections, `cnos.cdd` governs (handoff transports between cells; cell grammar lives in the kernel).

## v0.1 surface

`HANDOFF.md` and the five sub-skill files named in this loader's `calls:` frontmatter are all **landed and canonical** as of 2026-05-22 (cnos#404 wave closed via cnos#420). See [`docs/extraction-map.md`](../../docs/extraction-map.md) for the surface-by-surface migration provenance.

The canonical homes are:

- Cross-repo doctrine: [`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](cross-repo/SKILL.md) (landed under Sub 2 / [cnos#416](https://github.com/usurobor/cnos/issues/416); compatibility pointer at `cnos.cdd/skills/cdd/cross-repo/SKILL.md`)
- Dispatch + implementation-contract: [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](dispatch/SKILL.md) (landed under Sub 3 / [cnos#417](https://github.com/usurobor/cnos/issues/417); consumer-side role-skill realizations at `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5`, `delta/SKILL.md §2`, `alpha/SKILL.md §3.6`, `beta/SKILL.md` Rule 7)
- Mid-flight rescue: [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](mid-flight/SKILL.md) (landed under Sub 4 / [cnos#418](https://github.com/usurobor/cnos/issues/418); v0.1 software-cycle operational realization at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"`)
- Artifact channel: [`cnos.handoff/skills/handoff/artifact-channel/SKILL.md`](artifact-channel/SKILL.md) (landed under Sub 4 / [cnos#418](https://github.com/usurobor/cnos/issues/418); CDS-specific per-artifact contract at `cnos.cds/skills/cds/CDS.md §"Artifact contract"`)
- Receipt-stream: [`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`](receipt-stream/SKILL.md) (landed under Sub 5 / [cnos#419](https://github.com/usurobor/cnos/issues/419); cdd-side runbook pointer at `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`; live aggregator data at `.cdd/iterations/INDEX.md`)
