# CDR — Coherence-Driven Research

**CDR (Coherence-Driven Research)** is the research-protocol overlay in the c-d-X family. It instantiates the role-cell grammar (α/β/γ/δ/ε) and the six-field instantiation contract (per [`ROLES.md §3`](../../../ROLES.md)) for research work — investigation, synthesis, citation, dataset stewardship, claim transmission under uncertainty.

CDR is **not** a persona and **not** a project. It is the protocol layer — layer 3 of the five-layer enforcement chain — that sits between the persona (Rho lives in `cn-rho`) and the project binding (cph lives in `usurobor/cph`). See [`skills/cdr/CDR.md` §"Persona, Protocol, Project"](skills/cdr/CDR.md) for the boundary.

## What CDR Does

CDR is the research discipline used to transmit claims about the world (or about a system under study) coherently. Its loss function — per [`ROLES.md §4a.2`](../../../ROLES.md) — is **truth-preserving claim transmission under uncertainty**. The primary failure mode CDR must structurally resist is **overclaim**: a claim becoming stronger than its evidence; false knowledge propagating.

CDR shares the role-cell kernel (α/β/γ/δ/ε) with [`cnos.cdd`](../cnos.cdd/) (the generic kernel) and is a sibling of [`cnos.cds`](../cnos.cds/) (engineering realization). The discipline profile diverges: engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*. The role names are shared; the loss functions diverge; the receipts and oracles diverge to enforce the divergent disciplines.

The architectural split — common kernel in `cnos.cdd`, per-protocol procedures in `cnos.cdr` — inherits the option (a) decision from [cnos#388](https://github.com/usurobor/cnos/issues/388) and is recorded at [`skills/cdr/CDR.md` §"Architecture choice"](skills/cdr/CDR.md). The same split produced [`cnos.cds`](../cnos.cds/) v0.1 at [cnos#403](https://github.com/usurobor/cnos/issues/403); CDR is the structural precedent CDS mirrored.

## Package Structure

This package is **v0.1 complete**. The doctrine surface (`CDR.md` — the 616-line canonical instantiation contract), the loader skill (`SKILL.md`), all five per-role overlays under `skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md`, and the empirical-anchor doc (`docs/empirical-anchor-cph.md`) are all landed under [cnos#376](https://github.com/usurobor/cnos/issues/376). The wave shipped Sub 1 ([cnos#390](https://github.com/usurobor/cnos/issues/390)), Sub 2 ([cnos#394](https://github.com/usurobor/cnos/issues/394)), Sub 3 ([cnos#395](https://github.com/usurobor/cnos/issues/395)), and Sub 4 ([cnos#396](https://github.com/usurobor/cnos/issues/396)) sequentially.

### `/skills/cdr/` — The Method

- **`CDR.md`** — Canonical instantiation contract (616 lines). Declares the six fields per `ROLES.md §3`: matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule. This file is the package's **primary doctrine surface**. **Landed via Sub 1 of [cnos#376](https://github.com/usurobor/cnos/issues/376) — [cnos#390](https://github.com/usurobor/cnos/issues/390).**
- **`SKILL.md`** — Loader skill. Names `CDR.md` and the per-role overlays; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel) and `cnos.cds` (engineering sibling). **Landed via Sub 2 of [cnos#376](https://github.com/usurobor/cnos/issues/376) — [cnos#394](https://github.com/usurobor/cnos/issues/394).**

### `/skills/cdr/<role>/` — Per-role overlays (Sub 3 / [cnos#395](https://github.com/usurobor/cnos/issues/395))

- **`alpha/SKILL.md`** — α role in CDR. Produces research matter — claims, hypotheses, methods, datasets, analyses, reports — under the truth-preserving claim-transmission discipline.
- **`beta/SKILL.md`** — β role in CDR. Audits α's research matter for claim/evidence alignment, falsifiability, reproduction, citation integrity, and data-policy compliance.
- **`gamma/SKILL.md`** — γ role in CDR. Coordinates the research wave, selects research gaps, dispatches α/β, and closes the wave by emitting the typed `#CDRReceipt`.
- **`operator/SKILL.md`** — δ role in CDR (directory name `operator` per cdr convention; the δ role-cell position). Enforces the data-mounted gate, accepts or rejects research receipts, and records gate verdicts on the wave-transition cadence.
- **`epsilon/SKILL.md`** — ε role in CDR. Iterates the CDR protocol itself — observes research-protocol gaps from the receipt stream, applies MCA discipline, writes the CDR-iteration artifact.

All five role overlays are CDR-specific extensions of the corresponding `cnos.cdd/skills/cdd/<role>/SKILL.md` generic doctrine. The kernel grammar (role-cell shape, algorithm structure, independence rules, resumption protocol) is inherited by reference; the discipline profile and matter type diverge for the research-protocol loss function per `ROLES.md §4a.2`.

### `/docs/` — The Empirical Anchor

- **`empirical-anchor-cph.md`** — Surface-by-surface mapping of [`usurobor/cph`](https://github.com/usurobor/cph)'s `.cdr/` artifact set to `CDR.md`'s six-field structure, demonstrating that CDR v0.1 can host the cph research project without contradiction. **Landed via Sub 4 of [cnos#376](https://github.com/usurobor/cnos/issues/376) — [cnos#396](https://github.com/usurobor/cnos/issues/396).** The structural precedent for [`cnos.cds/docs/empirical-anchor-cdd.md`](../cnos.cds/docs/empirical-anchor-cdd.md) (cds Sub 7).

## Quick Start

**New to CDR?** Start with:

1. [`skills/cdr/CDR.md`](skills/cdr/CDR.md) — The instantiation contract; the doctrinal anchor. Six fields per `ROLES.md §3`; architecture-choice declaration; persona/protocol/project boundary; empirical-anchor citation to `usurobor/cph`.
2. [`skills/cdr/SKILL.md`](skills/cdr/SKILL.md) — The loader skill; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel) and `cnos.cds` (engineering sibling).
3. [`ROLES.md §3` + `§4a`](../../../ROLES.md) — The generic role-cell grammar and the five-layer enforcement chain CDR inherits.
4. [`schemas/cdr/receipt.cue`](../../../schemas/cdr/receipt.cue) — The typed γ close-out surface (`#CDRReceipt`) referenced by CDR.md Field 3.
5. [`docs/empirical-anchor-cph.md`](docs/empirical-anchor-cph.md) — Maps `usurobor/cph`'s research-cycle `.cdr/` artifacts to CDR's six-field structure; useful for grounding the doctrine in observable evidence.
6. [cnos#376](https://github.com/usurobor/cnos/issues/376) — parent tracker; rationale for the research-protocol extraction; sub-issue wave history.

**Looking for the engineering-side sibling?** See [`../cnos.cds/`](../cnos.cds/) — the engineering realization shipped at [cnos#403](https://github.com/usurobor/cnos/issues/403). CDS is **not** CDR-with-different-words: the role names are shared; the loss functions diverge.

**Looking for the kernel?** See [`../cnos.cdd/`](../cnos.cdd/) — the common kernel and the generic CCNF algorithm.

**Looking for the wire-format peer?** See [`../cnos.handoff/`](../cnos.handoff/) — the inter-agent / inter-activation / inter-repo handoff doctrine. CDR consumes `cnos.handoff/skills/handoff/cross-repo/SKILL.md` for the bootstrap-cdr and cph cross-repo bundles; CDR was the original cross-package consumer of cross-repo doctrine and is the empirical anchor proving handoff doctrine was already a separable boundary.

## Cross-protocol relationship

CDR shares the role-cell kernel with [`cnos.cdd`](../cnos.cdd/) (the generic algorithm) and is a sibling of [`cnos.cds`](../cnos.cds/) (the engineering realization). The kernel lives in `cnos.cdd`; CDR owns the **research-specific instantiation** (matter type, review oracle, γ close-out artifact, δ cadence, ε cadence, actor collapse). No generic CCNF kernel doctrine is restated in CDR — only pointers from CDR to CDD/CCNF.

**Boundary vs cnos.handoff:** handoff is the *wire format*; CDR is the *research-cycle lifecycle*. CDR consumes handoff surfaces (cross-repo bundles, dispatch envelopes, receipt streams) by reference; it does not own them. The two have different reasons to change.

## Status

**v0.1 complete — cnos#376 wave closed 2026-05-22.** All four v0.1 subs landed; doctrine surface in place; loader in place; all five role overlays in place; empirical-anchor doc in place. The wave shipped, in order:

- `cn.package.json` (schema `cn.package.v1`; name `cnos.cdr`; version `0.1.0`; kind `package`; engines `cnos >= 3.80.0`), `skills/cdr/CDR.md` (the 616-line six-field research instantiation contract; declares matter type, review oracle, γ close-out artifact, δ cadence, ε cadence, actor collapse rule; architecture-choice declaration per [cnos#388](https://github.com/usurobor/cnos/issues/388); persona/protocol/project boundary; empirical-anchor citation to `usurobor/cph`) — Sub 1 / [cnos#390](https://github.com/usurobor/cnos/issues/390);
- `README.md`, `skills/cdr/SKILL.md` (loader; declares load order + cross-protocol relationship) — Sub 2 / [cnos#394](https://github.com/usurobor/cnos/issues/394);
- five per-role overlays at `skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md` (each a CDR-specific extension of the corresponding `cnos.cdd/skills/cdd/<role>/SKILL.md`) — Sub 3 / [cnos#395](https://github.com/usurobor/cnos/issues/395);
- `docs/empirical-anchor-cph.md` mapping `usurobor/cph`'s `.cdr/` artifact set onto the CDR six-field structure — Sub 4 / [cnos#396](https://github.com/usurobor/cnos/issues/396).

The package is loadable and discoverable by cnos package-discovery (per `pkg.ContentClasses` filesystem-presence convention); commands are not declared in this version. The empirical-anchor doc became the structural precedent for [`cnos.cds/docs/empirical-anchor-cdd.md`](../cnos.cds/docs/empirical-anchor-cdd.md) (cds Sub 7 / [cnos#412](https://github.com/usurobor/cnos/issues/412)).

## License

Part of the cnos project.
