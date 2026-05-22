# CDS — Coherence-Driven Software

**CDS (Coherence-Driven Software)** is the software-development realization of CCNF — the generic recursive coherence-cell algorithm that lives in [`cnos.cdd`](../cnos.cdd/). CDS binds the kernel to source code, tests, branches, diffs, CI, releases, deployments, and the software-specific evidence the validator predicate `V` reads. It sits as a peer to [`cnos.cdr`](../cnos.cdr/) (research realization); both inherit the generic kernel from `cnos.cdd` without re-deriving it.

CDS is **not** a persona and **not** a project. It is the protocol layer — the substrate-specific instantiation of CCNF for software-class matter — that lives below the kernel (CCNF) and above per-project bindings (a project's `.cds/` directory, when CDS reaches v1). The kernel/realization split is the option (a) decision recorded at [cnos#388](https://github.com/usurobor/cnos/issues/388) (schemas) and extended to skills at [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376) (research realization, the structural precedent for this package).

## What CDS Does

CDS is the engineering discipline used to improve artifacts under repairable feedback. Its loss function — per [`ROLES.md §4a.2`](../../../ROLES.md) — is **artifact improvement under repairable feedback**. The primary failure mode CDS must structurally resist is **artifact degradation under unmeasured churn**: shipping code that drifts farther from coherence with each round, the failure of the engineering loop to bind on its own evidence.

CDS shares the role-cell kernel (α/β/γ/δ/ε) with [`cnos.cdr`](../cnos.cdr/) (research). The kernel grammar — `COHERENCE-CELL`, `COHERENCE-CELL-NORMAL-FORM`, role grammar, generic receipt-validation interface — lives in `cnos.cdd` and is cited by reference, never restated. CDS owns the software-specific instantiation: the selection function, the development lifecycle, the artifact contract, the review CLP, the gate / closure verification checklist, the assessment and cycle-iteration framework, the retro-packaging rule, and the software-class non-goals.

The architectural split — common kernel in `cnos.cdd`, per-protocol procedures here — inherits the option (a) decision from [cnos#388](https://github.com/usurobor/cnos/issues/388) and is the same split that produced [`cnos.cdr`](../cnos.cdr/) v0.1 at [cnos#376](https://github.com/usurobor/cnos/issues/376).

## Package Structure

This package is a v0.1 skeleton. The doctrine surface (`CDS.md`) and the per-role overlays are forthcoming sub-deliverables of [cnos#403](https://github.com/usurobor/cnos/issues/403).

### `/skills/cds/` — The Method

- **`SKILL.md`** — Loader skill. Names `CDS.md` and the forthcoming per-role overlays as advisory targets; declares the load order and the cross-protocol relationship. Shipped by [cnos#406](https://github.com/usurobor/cnos/issues/406) (Sub 1 of [cnos#403](https://github.com/usurobor/cnos/issues/403)) — this sub. Does not call `CDS.md` yet (Sub 2 lands it).
- **`CDS.md`** — Canonical instantiation contract. Will declare the six fields per [`ROLES.md §3`](../../../ROLES.md): matter type (software-class), review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule. **Pending Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403).**

### `/docs/` — The Extraction Plan

- **`extraction-map.md`** — Names, for every pre-#402 software-lifecycle surface still resident in [`cnos.cdd/skills/cdd/CDD.md` §"Software-specific realization — pending cds extraction"](../cnos.cdd/skills/cdd/CDD.md), the planned CDS destination and the sub-issue of [cnos#403](https://github.com/usurobor/cnos/issues/403) that will execute the migration. **The load-bearing artifact of Sub 1.** Subs 3–5 dispatch against this map.

### Forthcoming surfaces

- **`skills/cds/CDS.md`** — the six-field software instantiation contract. Shipped by Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
- **Per-role overlays** — `skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` declare each role's software-engineering-discipline procedure (matter production, review oracles, close-out artifacts, dispatch cadence, protocol-iteration cadence). Shipped by Subs 3–5 of [cnos#403](https://github.com/usurobor/cnos/issues/403). Pointer overlays delegating to existing `cnos.cdd` role skills are acceptable in v0.1.
- **`docs/empirical-anchor-cdd.md`** — surface-by-surface mapping of the current cnos repository's `.cdd/` cycle artifacts to CDS's six-field structure, demonstrating that CDS v0.1 can host current practice without contradiction. Shipped by Sub 7 of [cnos#403](https://github.com/usurobor/cnos/issues/403). The structural precedent is [`cnos.cdr/docs/empirical-anchor-cph.md`](../cnos.cdr/docs/empirical-anchor-cph.md).
- **Receipt schemas** — `schemas/cds/` already exists per [cnos#388](https://github.com/usurobor/cnos/issues/388); CDS's six-field contract Field 3 (γ close-out artifact) cites it.

## Quick Start

**New to CDS?** Start with:

1. [`docs/extraction-map.md`](docs/extraction-map.md) — Read the extraction map first. It names where each software-lifecycle surface currently in `CDD.md` will live once Subs 3–5 of [cnos#403](https://github.com/usurobor/cnos/issues/403) land.
2. [`skills/cds/SKILL.md`](skills/cds/SKILL.md) — The loader skill; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel) and `cnos.cdr` (sibling research realization).
3. [`cnos.cdd/skills/cdd/CDD.md`](../cnos.cdd/skills/cdd/CDD.md) — The CCNF spine. CDS instantiates this kernel; reading the kernel first grounds the realization.
4. [`COHERENCE-CELL-NORMAL-FORM.md`](../cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md) and [`COHERENCE-CELL.md`](../cnos.cdd/skills/cdd/COHERENCE-CELL.md) — the kernel doctrine CDS cites by reference and does not restate.
5. [`schemas/cds/receipt.cue`](../../../schemas/cds/receipt.cue) — the typed γ close-out surface (`#CDSReceipt`) that CDS Field 3 will reference once `CDS.md` lands.

**Looking for the kernel?** See [`../cnos.cdd/`](../cnos.cdd/) — the generic recursive coherence-cell algorithm.

**Looking for the research-side sibling?** See [`../cnos.cdr/`](../cnos.cdr/) — the research realization shipped at [cnos#376](https://github.com/usurobor/cnos/issues/376). CDR is **not** CDS-with-different-words: engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*. The role names are shared; the loss functions diverge; the receipts and oracles diverge to enforce the divergent disciplines.

## Status

**v0.1 — package skeleton + extraction map.** This Sub 1 of [cnos#403](https://github.com/usurobor/cnos/issues/403) ships:

- `cn.package.json` (schema `cn.package.v1`; name `cnos.cds`; version `0.1.0`; kind `package`; engines `cnos >= 3.81.0`),
- `README.md` (this file),
- `skills/cds/SKILL.md` (loader, advisory),
- `docs/extraction-map.md` (the migration plan Subs 3–5 dispatch against).

`CDS.md` (Sub 2), per-role overlays (Subs 3–5), the cross-reference sweep of `CDD.md` (Sub 6), and the empirical-anchor doc (Sub 7) are in flight under [cnos#403](https://github.com/usurobor/cnos/issues/403). The package is loadable and discoverable by cnos package-discovery (per `pkg.ContentClasses` filesystem-presence convention); commands are not declared in this version.

Until Sub 6 lands, software-lifecycle realization remains resident in [`cnos.cdd/skills/cdd/CDD.md` §"Software-specific realization — pending cds extraction"](../cnos.cdd/skills/cdd/CDD.md). The extraction map names how each marker will resolve.

## License

Part of the cnos project.
