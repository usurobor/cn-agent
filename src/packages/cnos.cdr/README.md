# CDR — Coherence-Driven Research

**CDR (Coherence-Driven Research)** is the research-protocol overlay in the c-d-X family. It instantiates the role-cell grammar (α/β/γ/δ/ε) and the six-field instantiation contract (per [`ROLES.md §3`](../../../ROLES.md)) for research work — investigation, synthesis, citation, dataset stewardship, claim transmission under uncertainty.

CDR is **not** a persona and **not** a project. It is the protocol layer — layer 3 of the five-layer enforcement chain — that sits between the persona (Rho lives in `cn-rho`) and the project binding (cph lives in `usurobor/cph`). See [`skills/cdr/CDR.md` §"Persona, Protocol, Project"](skills/cdr/CDR.md) for the boundary.

## What CDR Does

CDR is the research discipline used to transmit claims about the world (or about a system under study) coherently. Its loss function — per [`ROLES.md §4a.2`](../../../ROLES.md) — is **truth-preserving claim transmission under uncertainty**. The primary failure mode CDR must structurally resist is **overclaim**: a claim becoming stronger than its evidence; false knowledge propagating.

CDR shares the role-cell kernel (α/β/γ/δ/ε) with cnos.cdd (engineering); the discipline profile diverges. The architectural split — common kernel in `cnos.cdd`, per-protocol procedures in `cnos.cdr` — inherits the option (a) decision from [cnos#388](https://github.com/usurobor/cnos/issues/388) and is recorded at [`skills/cdr/CDR.md` §"Architecture choice"](skills/cdr/CDR.md).

## Package Structure

This package contains a single surface in this v0.1 skeleton:

### `/skills/cdr/` — The Method

- **`CDR.md`** — Canonical instantiation contract (declares the six fields per `ROLES.md §3`: matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule). Shipped by [cnos#390](https://github.com/usurobor/cnos/issues/390) (Sub 1 of [cnos#376](https://github.com/usurobor/cnos/issues/376)). This file is the package's **primary doctrine surface**.
- **`SKILL.md`** — Loader skill. Names CDR.md and the forthcoming per-role overlays as sub-skills without defining them. Shipped by [cnos#394](https://github.com/usurobor/cnos/issues/394) (Sub 2 of cnos#376).

### Forthcoming surfaces

- **Per-role overlays** — `skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` declare each role's research-discipline procedure (matter production, review oracles, close-out artifacts, dispatch cadence, protocol-iteration cadence). Shipped by Sub 3 of cnos#376.
- **Empirical-anchor doc** — A surface-by-surface mapping of [`usurobor/cph`](https://github.com/usurobor/cph)'s `.cdr/` artifact set to CDR.md's six-field structure. Shipped by Sub 4 of cnos#376.

## Quick Start

**New to CDR?** Start with:

1. [`skills/cdr/CDR.md`](skills/cdr/CDR.md) — The instantiation contract; the doctrinal anchor.
2. [`ROLES.md §3` + `§4a`](../../../ROLES.md) — The generic role-cell grammar and the five-layer enforcement chain CDR inherits.
3. [`schemas/cdr/receipt.cue`](../../../schemas/cdr/receipt.cue) — The typed γ close-out surface (`#CDRReceipt`) referenced by CDR.md Field 3.

**Looking for engineering-side CDD?** See [`../cnos.cdd/`](../cnos.cdd/) — the common kernel and engineering-protocol overlay.

## Status

**v0.1 — package skeleton.** Doctrine (CDR.md) and loader (SKILL.md) shipped. Role overlays (Sub 3) and empirical-anchor doc (Sub 4) are in flight under cnos#376. The package is loadable and discoverable by cnos package-discovery (per `pkg.ContentClasses` filesystem-presence convention); commands are not declared in this version.

## License

Part of the cnos project.
