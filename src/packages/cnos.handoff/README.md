# cnos.handoff — Inter-Agent / Inter-Activation / Inter-Repo Handoff

**`cnos.handoff`** is the package that owns the **wire formats and handoff doctrine** by which work, state, and messages move between cells, agents, roles, activations, repos, and packages. It sits as a peer to [`cnos.cdd`](../cnos.cdd/) (the generic recursive coherence-cell kernel), [`cnos.cds`](../cnos.cds/) (software realization), and [`cnos.cdr`](../cnos.cdr/) (research realization). Those three packages **consume** handoff doctrine; `cnos.handoff` is the canonical home of the doctrine itself.

The package extracts a coordination layer that has quietly accreted inside `cnos.cdd` over the past wave of cycles — cross-repo bundles, dispatch-prompt templates, implementation-contract enrichment, mid-flight clarification, intra-cycle artifact channels, cross-cycle receipt streams — and gives it one canonical home with one reason to change. The accretion was structurally mixed with the CCNF cycle-lifecycle kernel; the extraction is the Unix-discipline application of [`cnos.core/skills/design/SKILL.md`](../cnos.core/skills/design/SKILL.md) §3.1 (one reason to change per boundary) + §3.5 (interfaces belong to consumers) + §3.10 (package/install cohesion) + §3.11 (small bootstrap kernel) once the consumer count crossed the two-real-plus-one-imminent threshold per the parent tracker [cnos#404](https://github.com/usurobor/cnos/issues/404).

## What handoff Does

The handoff package's loss function is **wire-format coherence under multi-consumer extension** — every doctrine surface in this package answers a single question: *how do two parts of the system that are not co-located in time, role, repo, or activation transfer state without inventing the protocol per cycle?* Handoff doctrine is mechanically separable from the cycle lifecycle because the two have **different rates of change**: the cycle lifecycle changes when a domain protocol (CDD, CDS, CDR, future c-d-X) refines its loop; the handoff layer changes when the *agent-to-agent wire format* changes. Different rates, different reasons, different boundaries.

The primary failure mode `cnos.handoff` must structurally resist is **per-cycle doctrine re-invention**: γ inventing the cross-repo bundle shape because the protocol does not name the directional case; δ inventing the implementation-contract format because the dispatch template lives across three different role skills; α/β/γ inventing the intra-cycle artifact channel because `CDD.md` references it but no skill canonicalizes it. The result of the failure mode is doctrine that lives inside individual bundles and individual cycles rather than in the protocol — drift that compounds because no single file is the citable source.

Handoff doctrine is **substrate-independent for the message; substrate-specific for the realization.** The cross-repo state machine, the dispatch-prompt template, the implementation-contract schema, the mid-flight rescue mechanism, the intra-cycle artifact channel rules, the cross-cycle receipt-stream — these are wire-format facts that any coherence-cell-shaped protocol package consumes. Their realization in a specific package (CDD's `gamma-scaffold.md` filenames; CDS's `.cdd/unreleased/{N}/` rooting; CDR's wave-shaped variants) is the consumer's call. Handoff names the format; consumers name the binding.

## Package Structure

This package is a **v0.1 skeleton**. The doctrine surface (`HANDOFF.md`) and the per-surface sub-skills (`cross-repo/`, `dispatch/`, `mid-flight/`, `artifact-channel/`, `receipt-stream/`) are forthcoming sub-deliverables of [cnos#404](https://github.com/usurobor/cnos/issues/404), executed by Subs 2–5. This Sub 1 ([cnos#415](https://github.com/usurobor/cnos/issues/415)) ships the package boundary, the loader skill, and the extraction map; no content is migrated in this cycle.

### `/skills/handoff/` — The Method

- **`SKILL.md`** — Loader skill. Names `HANDOFF.md` and the forthcoming per-surface sub-skill files as advisory targets; declares the load order and the cross-protocol relationship. Shipped by [cnos#415](https://github.com/usurobor/cnos/issues/415) (Sub 1 of [cnos#404](https://github.com/usurobor/cnos/issues/404)) — this sub. Does not call `HANDOFF.md` substantively yet (Subs 2–5 land the substance).
- **`HANDOFF.md`** — Canonical instantiation contract for the handoff wire format. Will declare: directional cases (intra-cycle / inter-cycle / inter-activation / inter-repo); the artifact families per case; the typed schema(s) (if any are lifted from Markdown into `schemas/handoff/`); the consumer protocol (how a CDD/CDS/CDR/future-c-d-X package binds the wire format to its substrate). **Pending Subs 2–5 of [cnos#404](https://github.com/usurobor/cnos/issues/404).** This file is intentionally absent from the v0.1 skeleton; the `skills/handoff/` directory is tracked via `.gitkeep`.

### `/docs/` — The Extraction Plan

- **`extraction-map.md`** — Names, for every handoff-class surface currently resident in `cnos.cdd` (or already migrated into `cnos.cds` per the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave), the planned `cnos.handoff` destination and the sub-issue of [cnos#404](https://github.com/usurobor/cnos/issues/404) that will execute the migration. **The load-bearing artifact of Sub 1.** Subs 2–6 dispatch against this map. Six required surface families are covered; additional surfaces discovered during authoring are recorded as additional rows.

### Forthcoming surfaces

- **`skills/handoff/HANDOFF.md`** — the canonical wire-format contract. Shipped by Sub 2 or Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404) (the sub that lands the first substantive content decides the file's authoring order; the loader names it as the canonical home regardless).
- **`skills/handoff/cross-repo/SKILL.md`** — cross-repo state machine, bundle file sets per directional case, LINEAGE schemas, feedback-patch format, bundle archival rule, hat-collapse attribution. Shipped by Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404); migrates from `cnos.cdd/skills/cdd/cross-repo/SKILL.md` (644 lines; the primary surface this package absorbs).
- **`skills/handoff/dispatch/SKILL.md`** — dispatch-prompt template (operator → α/β/γ/δ); implementation-contract 7-axes schema; δ-as-inward-membrane enrichment doctrine. **Landed via Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — [cnos#417](https://github.com/usurobor/cnos/issues/417).** Synthesized from `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` + `cnos.cdd/skills/cdd/operator/SKILL.md §3a` + `cnos.cdd/skills/cdd/delta/SKILL.md §2`.
- **`skills/handoff/mid-flight/SKILL.md`** — `gamma-clarification.md` mechanism for intra-cycle re-pinning of axis values; empirically anchored at [cnos#391](https://github.com/usurobor/cnos/issues/391). Shipped by Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
- **`skills/handoff/artifact-channel/SKILL.md`** — intra-cycle α→β→γ sequential handoff rules over `.cdd/unreleased/{N}/`. Shipped by Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404); the canonical home post-#411 is `cnos.cds/skills/cds/CDS.md §"Artifact contract"`, so Sub 4 decides migration semantics (move / mirror / cite-only) per the extraction-map row note.
- **`skills/handoff/receipt-stream/SKILL.md`** — cross-cycle `cdd-iteration.md` receipt-stream + `.cdd/iterations/INDEX.md` aggregator; ε reader feed. Shipped by Sub 5 of [cnos#404](https://github.com/usurobor/cnos/issues/404); migrates from `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`.

## Quick Start

**New to cnos.handoff?** Start with:

1. [`docs/extraction-map.md`](docs/extraction-map.md) — Read the extraction map first. It names where each handoff-class surface currently in `cnos.cdd` (or in `cnos.cds` post-#411) will live once Subs 2–5 of [cnos#404](https://github.com/usurobor/cnos/issues/404) land.
2. [`skills/handoff/SKILL.md`](skills/handoff/SKILL.md) — The loader skill; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel), `cnos.cds` (software realization), and `cnos.cdr` (research realization).
3. [`../cnos.cdd/skills/cdd/cross-repo/SKILL.md`](../cnos.cdd/skills/cdd/cross-repo/SKILL.md) — The 644-line primary extraction candidate. Reading this surface today gives the clearest picture of what `cnos.handoff` will own once Sub 2 lands the move.
4. [`skills/handoff/dispatch/SKILL.md`](skills/handoff/dispatch/SKILL.md) — the dispatch-prompt template + 7-axis implementation-contract schema + δ-as-inward-membrane doctrine, landed under Sub 3 / [cnos#417](https://github.com/usurobor/cnos/issues/417). Consumer-side role-skill realizations live at `../cnos.cdd/skills/cdd/{gamma,operator,delta,alpha,beta}/SKILL.md`.
5. [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker; rationale for the extraction; consumer enumeration; design-skill citations.

**Looking for the cycle-lifecycle kernel?** See [`../cnos.cdd/`](../cnos.cdd/) — CCNF (Coherence Cell Normal Form) and the generic recursive coherence-cell algorithm.

**Looking for the software realization?** See [`../cnos.cds/`](../cnos.cds/) — the software-development binding of CCNF; current consumer of handoff doctrine via `cnos.cdd/skills/cdd/cross-repo/SKILL.md` (until Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404) moves it here).

**Looking for the research realization?** See [`../cnos.cdr/`](../cnos.cdr/) — the research binding; the original cross-package consumer of `cdd/cross-repo/SKILL.md` (the empirical anchor that proved handoff doctrine was already a separable boundary).

## Cross-protocol relationship

`cnos.handoff` owns **wire formats and handoff doctrine for inter-agent / inter-activation / inter-repo transport.** Three current consumers — [`cnos.cdd`](../cnos.cdd/) (which uses handoff doctrine natively for its own coordination), [`cnos.cdr`](../cnos.cdr/) (which loads `cross-repo/SKILL.md` via cross-package skill reference), and [`cnos.cds`](../cnos.cds/) (which absorbed some surfaces during the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave and continues to cite others from `cnos.cdd`) — and an open class of future c-d-X protocol packages that will consume the same primitives once the extraction lands.

**Boundary vs CCNF-O (Track A of [cnos#405](https://github.com/usurobor/cnos/issues/405)):** `cnos.handoff` transports work between cells; CCNF-O composes cells into cycles, waves, roadmaps, gates, joins. The two have **different reasons to change**: handoff changes when the wire format changes (a new directional case, a new bundle schema, a new clarification mechanism); CCNF-O changes when the composition grammar changes (a new gate type, a new join semantic, a new sequencing primitive). Each has one reason to change. Sub 1 of [cnos#404](https://github.com/usurobor/cnos/issues/404) preserves this boundary by explicitly not preempting CCNF-O surfaces — no orchestration-composition grammar appears in this package; Track A of [cnos#405](https://github.com/usurobor/cnos/issues/405) is gated on the `cnos#404` wave closing.

**Boundary vs the domain protocols (CDD / CDS / CDR / future c-d-X):** handoff is the *wire format*; the domain protocols are the *cycle lifecycles*. CDD changes when the generic kernel changes; CDS changes when the software-engineering cycle shape changes; CDR changes when the research-cycle shape changes. None of those changes ought to ripple into handoff doctrine, and vice versa — a new cross-repo directional case ought not require a CDS lifecycle revision. The extraction makes this independence structural rather than aspirational.

## Status

**v0.1 — package skeleton + extraction map.** This Sub 1 of [cnos#404](https://github.com/usurobor/cnos/issues/404) (executed under [cnos#415](https://github.com/usurobor/cnos/issues/415)) ships:

- `cn.package.json` (schema `cn.package.v1`; name `cnos.handoff`; version `0.1.0`; kind `package`; engines `cnos >= 3.81.0`),
- `README.md` (this file),
- `skills/handoff/SKILL.md` (loader, advisory; calls forthcoming sub-skills),
- `skills/handoff/.gitkeep` (placeholder for the empty directory; `HANDOFF.md` is intentionally absent in v0.1),
- `docs/extraction-map.md` (the migration plan Subs 2–6 dispatch against).

`HANDOFF.md` (the canonical wire-format contract), the per-surface sub-skills (`cross-repo/`, `dispatch/`, `mid-flight/`, `artifact-channel/`, `receipt-stream/`), the cross-reference sweep in cdd / cds / cdr (Sub 6), and any schemas lifted into `schemas/handoff/` are in flight under [cnos#404](https://github.com/usurobor/cnos/issues/404). The package is loadable and discoverable by cnos package-discovery (per `pkg.ContentClasses` filesystem-presence convention); commands are not declared in this version.

Until Subs 2–5 land, the handoff-class surfaces remain resident at their current canonical homes — most in `cnos.cdd`, some already in `cnos.cds` per the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave closure. The extraction map names how each surface will resolve.

## License

Part of the cnos project.
