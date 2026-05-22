# cnos.handoff — Inter-Agent / Inter-Activation / Inter-Repo Handoff

**`cnos.handoff`** is the package that owns the **wire formats and handoff doctrine** by which work, state, and messages move between cells, agents, roles, activations, repos, and packages. It sits as a peer to [`cnos.cdd`](../cnos.cdd/) (the generic recursive coherence-cell kernel), [`cnos.cds`](../cnos.cds/) (software realization), and [`cnos.cdr`](../cnos.cdr/) (research realization). Those three packages **consume** handoff doctrine; `cnos.handoff` is the canonical home of the doctrine itself.

The package extracts a coordination layer that has quietly accreted inside `cnos.cdd` over the past wave of cycles — cross-repo bundles, dispatch-prompt templates, implementation-contract enrichment, mid-flight clarification, intra-cycle artifact channels, cross-cycle receipt streams — and gives it one canonical home with one reason to change. The accretion was structurally mixed with the CCNF cycle-lifecycle kernel; the extraction is the Unix-discipline application of [`cnos.core/skills/design/SKILL.md`](../cnos.core/skills/design/SKILL.md) §3.1 (one reason to change per boundary) + §3.5 (interfaces belong to consumers) + §3.10 (package/install cohesion) + §3.11 (small bootstrap kernel) once the consumer count crossed the two-real-plus-one-imminent threshold per the parent tracker [cnos#404](https://github.com/usurobor/cnos/issues/404).

## What handoff Does

The handoff package's loss function is **wire-format coherence under multi-consumer extension** — every doctrine surface in this package answers a single question: *how do two parts of the system that are not co-located in time, role, repo, or activation transfer state without inventing the protocol per cycle?* Handoff doctrine is mechanically separable from the cycle lifecycle because the two have **different rates of change**: the cycle lifecycle changes when a domain protocol (CDD, CDS, CDR, future c-d-X) refines its loop; the handoff layer changes when the *agent-to-agent wire format* changes. Different rates, different reasons, different boundaries.

The primary failure mode `cnos.handoff` must structurally resist is **per-cycle doctrine re-invention**: γ inventing the cross-repo bundle shape because the protocol does not name the directional case; δ inventing the implementation-contract format because the dispatch template lives across three different role skills; α/β/γ inventing the intra-cycle artifact channel because `CDD.md` references it but no skill canonicalizes it. The result of the failure mode is doctrine that lives inside individual bundles and individual cycles rather than in the protocol — drift that compounds because no single file is the citable source.

Handoff doctrine is **substrate-independent for the message; substrate-specific for the realization.** The cross-repo state machine, the dispatch-prompt template, the implementation-contract schema, the mid-flight rescue mechanism, the intra-cycle artifact channel rules, the cross-cycle receipt-stream — these are wire-format facts that any coherence-cell-shaped protocol package consumes. Their realization in a specific package (CDD's `gamma-scaffold.md` filenames; CDS's `.cdd/unreleased/{N}/` rooting; CDR's wave-shaped variants) is the consumer's call. Handoff names the format; consumers name the binding.

## Package Structure

This package is **v0.1 complete**. The doctrine surface (`HANDOFF.md`) and all five per-surface sub-skills (`cross-repo/`, `dispatch/`, `mid-flight/`, `artifact-channel/`, `receipt-stream/`) are landed under [cnos#404](https://github.com/usurobor/cnos/issues/404). The wave bootstrapped under Sub 1 ([cnos#415](https://github.com/usurobor/cnos/issues/415)), migrated content under Subs 2–5 ([cnos#416](https://github.com/usurobor/cnos/issues/416), [cnos#417](https://github.com/usurobor/cnos/issues/417), [cnos#418](https://github.com/usurobor/cnos/issues/418), [cnos#419](https://github.com/usurobor/cnos/issues/419)), and closed via the Sub 6 cross-reference sweep ([cnos#420](https://github.com/usurobor/cnos/issues/420)) on 2026-05-22.

### `/skills/handoff/` — The Method

- **`SKILL.md`** — Loader skill. Names `HANDOFF.md` and the five canonical per-surface sub-skill files; declares the load order and the cross-protocol relationship. Shipped by [cnos#415](https://github.com/usurobor/cnos/issues/415) (Sub 1); tightened to v0.1-complete state by [cnos#420](https://github.com/usurobor/cnos/issues/420) (Sub 6).
- **`HANDOFF.md`** — Canonical instantiation contract for the handoff wire format. Declares: the package contract (what cnos.handoff owns vs does not own); the five Landed sub-surfaces with their migration provenance; the related-documents and non-goals.

### `/docs/` — The Extraction Plan

- **`extraction-map.md`** — Surface-by-surface migration map for every handoff-class surface that was resident in `cnos.cdd` (or already in `cnos.cds` per the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave) at the source pin. Subs 2–6 of [cnos#404](https://github.com/usurobor/cnos/issues/404) dispatched against this map; every row carries a "Status" field declaring the v0.1 outcome (migrated / deferred-by-design / complete). Authored under Sub 1 ([cnos#415](https://github.com/usurobor/cnos/issues/415)); finalised under Sub 6 ([cnos#420](https://github.com/usurobor/cnos/issues/420)).

### Surfaces (v0.1 Landed)

- **`skills/handoff/cross-repo/SKILL.md`** — cross-repo state machine, bundle file sets per directional case, LINEAGE schemas, feedback-patch format, bundle archival rule, hat-collapse attribution. **Landed via Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — [cnos#416](https://github.com/usurobor/cnos/issues/416).** Migrated wholesale from `cnos.cdd/skills/cdd/cross-repo/SKILL.md` (644 lines); a 28-line compatibility pointer remains at the old path for backward compatibility.
- **`skills/handoff/dispatch/SKILL.md`** — dispatch-prompt template (operator → α/β/γ/δ); implementation-contract 7-axes schema; δ-as-inward-membrane enrichment doctrine. **Landed via Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — [cnos#417](https://github.com/usurobor/cnos/issues/417).** Synthesized from `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` + `cnos.cdd/skills/cdd/operator/SKILL.md §3a` + `cnos.cdd/skills/cdd/delta/SKILL.md §2` (consumer-side role-skill realizations remain at those paths).
- **`skills/handoff/mid-flight/SKILL.md`** — `gamma-clarification.md` mechanism for intra-cycle re-pinning of axis values; empirically anchored at [cnos#391](https://github.com/usurobor/cnos/issues/391). **Landed via Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — [cnos#418](https://github.com/usurobor/cnos/issues/418).** v0.1 software-cycle operational realization remains at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"`.
- **`skills/handoff/artifact-channel/SKILL.md`** — intra-cycle α→β→γ sequential handoff rules over `.cdd/unreleased/{N}/`; per-role write ownership; frozen-snapshot rule; release-time directory move. **Landed via Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — [cnos#418](https://github.com/usurobor/cnos/issues/418).** Wire-format invariants live here; the CDS-specific per-artifact contract (filenames, ownership matrix, 13-stage ordered flow) remains canonical at `cnos.cds/skills/cds/CDS.md §"Artifact contract"`.
- **`skills/handoff/receipt-stream/SKILL.md`** — cross-cycle `cdd-iteration.md` receipt-stream + `.cdd/iterations/INDEX.md` aggregator; ε reader feed. **Landed via Sub 5 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — [cnos#419](https://github.com/usurobor/cnos/issues/419).** Migrated from `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` (now a pointer to this skill).

## Quick Start

**New to cnos.handoff?** Start with:

1. [`skills/handoff/HANDOFF.md`](skills/handoff/HANDOFF.md) — The canonical package contract: what cnos.handoff owns, what it does not own, and the five Landed sub-surfaces with their migration provenance.
2. [`skills/handoff/SKILL.md`](skills/handoff/SKILL.md) — The loader skill; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel), `cnos.cds` (software realization), and `cnos.cdr` (research realization).
3. [`skills/handoff/cross-repo/SKILL.md`](skills/handoff/cross-repo/SKILL.md) — Cross-repo state machine + bundles + LINEAGE schemas + feedback-patch format + archival + hat-collapse attribution. The largest of the five sub-surfaces (644 lines); the primary surface this package absorbed.
4. [`skills/handoff/dispatch/SKILL.md`](skills/handoff/dispatch/SKILL.md) — the dispatch-prompt template + 7-axis implementation-contract schema + δ-as-inward-membrane doctrine. Consumer-side role-skill realizations live at `../cnos.cdd/skills/cdd/{gamma,operator,delta,alpha,beta}/SKILL.md`.
5. [`docs/extraction-map.md`](docs/extraction-map.md) — Per-surface migration map with v0.1-complete status; reading order useful for archaeological context on what moved from where.
6. [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker (closed 2026-05-22); rationale for the extraction; consumer enumeration; design-skill citations.

**Looking for the cycle-lifecycle kernel?** See [`../cnos.cdd/`](../cnos.cdd/) — CCNF (Coherence Cell Normal Form) and the generic recursive coherence-cell algorithm.

**Looking for the software realization?** See [`../cnos.cds/`](../cnos.cds/) — the software-development binding of CCNF; consumer of handoff doctrine via `cnos.handoff/skills/handoff/`.

**Looking for the research realization?** See [`../cnos.cdr/`](../cnos.cdr/) — the research binding; the original cross-package consumer of cross-repo doctrine (the empirical anchor that proved handoff doctrine was already a separable boundary; now cites `cnos.handoff/skills/handoff/cross-repo/SKILL.md` as canonical).

## Cross-protocol relationship

`cnos.handoff` owns **wire formats and handoff doctrine for inter-agent / inter-activation / inter-repo transport.** Three current consumers — [`cnos.cdd`](../cnos.cdd/) (which uses handoff doctrine natively for its own coordination), [`cnos.cdr`](../cnos.cdr/) (which loads `cross-repo/SKILL.md` via cross-package skill reference), and [`cnos.cds`](../cnos.cds/) (which absorbed some surfaces during the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave and continues to cite others from `cnos.cdd`) — and an open class of future c-d-X protocol packages that will consume the same primitives once the extraction lands.

**Boundary vs CCNF-O (Track A of [cnos#405](https://github.com/usurobor/cnos/issues/405)):** `cnos.handoff` transports work between cells; CCNF-O composes cells into cycles, waves, roadmaps, gates, joins. The two have **different reasons to change**: handoff changes when the wire format changes (a new directional case, a new bundle schema, a new clarification mechanism); CCNF-O changes when the composition grammar changes (a new gate type, a new join semantic, a new sequencing primitive). Each has one reason to change. Sub 1 of [cnos#404](https://github.com/usurobor/cnos/issues/404) preserves this boundary by explicitly not preempting CCNF-O surfaces — no orchestration-composition grammar appears in this package; Track A of [cnos#405](https://github.com/usurobor/cnos/issues/405) is gated on the `cnos#404` wave closing.

**Boundary vs the domain protocols (CDD / CDS / CDR / future c-d-X):** handoff is the *wire format*; the domain protocols are the *cycle lifecycles*. CDD changes when the generic kernel changes; CDS changes when the software-engineering cycle shape changes; CDR changes when the research-cycle shape changes. None of those changes ought to ripple into handoff doctrine, and vice versa — a new cross-repo directional case ought not require a CDS lifecycle revision. The extraction makes this independence structural rather than aspirational.

## Status

**v0.1 complete — cnos#404 wave closed 2026-05-22.** All five v0.1 sub-skills landed; cross-reference sweep complete; tracker closed via [cnos#420](https://github.com/usurobor/cnos/issues/420). The wave shipped, in order:

- `cn.package.json` (schema `cn.package.v1`; name `cnos.handoff`; version `0.1.0`; kind `package`; engines `cnos >= 3.81.0`), `README.md`, `skills/handoff/SKILL.md` (loader), `docs/extraction-map.md` — Sub 1 / [cnos#415](https://github.com/usurobor/cnos/issues/415);
- `skills/handoff/HANDOFF.md` (the canonical package contract) + `skills/handoff/cross-repo/SKILL.md` (the 644-line wholesale move; compatibility pointer left at `cnos.cdd/skills/cdd/cross-repo/SKILL.md`) — Sub 2 / [cnos#416](https://github.com/usurobor/cnos/issues/416);
- `skills/handoff/dispatch/SKILL.md` (synthesised from three cdd-side sources; consumer-side pointers in cdd role skills) — Sub 3 / [cnos#417](https://github.com/usurobor/cnos/issues/417);
- `skills/handoff/mid-flight/SKILL.md` + `skills/handoff/artifact-channel/SKILL.md` (mid-flight rescue mechanism + intra-cycle artifact channel; cds-side per-artifact contract preserved) — Sub 4 / [cnos#418](https://github.com/usurobor/cnos/issues/418);
- `skills/handoff/receipt-stream/SKILL.md` (cdd-iteration.md receipt-stream + INDEX.md aggregator; `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` now a pointer) — Sub 5 / [cnos#419](https://github.com/usurobor/cnos/issues/419);
- final cross-reference cleanup + v0.1-complete tightening of loader / HANDOFF.md / README.md / extraction-map.md — Sub 6 / [cnos#420](https://github.com/usurobor/cnos/issues/420) (this cycle).

No schemas were lifted into `schemas/handoff/` during the v0.1 wave (deferred per Sub 3 / Q2 ruling — markdown 7-axis schema preserved as-is). The package is loadable and discoverable by cnos package-discovery (per `pkg.ContentClasses` filesystem-presence convention); commands are not declared in this version.

With the cnos#404 wave closed, the cnos#405 Track A (CCNF-O orchestration grammar) gate is satisfied and Track B1 (TSC report attachment) may proceed.

## License

Part of the cnos project.
