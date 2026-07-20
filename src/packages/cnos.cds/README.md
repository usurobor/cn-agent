# CDS — Coherence-Driven Software

**CDS (Coherence-Driven Software)** is the software-development realization of CCNF — the generic recursive coherence-cell algorithm that lives in [`cnos.cdd`](../cnos.cdd/). CDS binds the kernel to source code, tests, branches, diffs, CI, releases, deployments, and the software-specific evidence the validator predicate `V` reads. It sits as a peer to [`cnos.cdr`](../cnos.cdr/) (research realization) and [`cnos.handoff`](../cnos.handoff/) (inter-agent / inter-activation / inter-repo wire-format doctrine); all three packages inherit the generic kernel from `cnos.cdd` without re-deriving it.

CDS is **not** a persona and **not** a project. It is the protocol layer — the substrate-specific instantiation of CCNF for software-class matter — that lives below the kernel (CCNF) and above per-project bindings (a project's `.cds/` directory, when CDS reaches v1). The kernel/realization split is the option (a) decision recorded at [cnos#388](https://github.com/usurobor/cnos/issues/388) (schemas) and extended to skills at [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376) (research realization, the structural precedent for this package).

## What CDS Does

CDS is the engineering discipline used to improve artifacts under repairable feedback. Its loss function — per [`ROLES.md §4a.2`](../../../ROLES.md) — is **artifact improvement under repairable feedback**. The primary failure mode CDS must structurally resist is **artifact degradation under unmeasured churn**: shipping code that drifts farther from coherence with each round, the failure of the engineering loop to bind on its own evidence.

CDS shares the role-cell kernel (α/β/γ/δ/ε) with [`cnos.cdr`](../cnos.cdr/) (research). The kernel grammar — `COHERENCE-CELL`, `COHERENCE-CELL-NORMAL-FORM`, role grammar, generic receipt-validation interface — lives in `cnos.cdd` and is cited by reference, never restated. CDS owns the software-specific instantiation: the selection function, the development lifecycle, the artifact contract, the review CLP, the gate / closure verification checklist, the assessment and cycle-iteration framework, the retro-packaging rule, and the software-class non-goals.

The architectural split — common kernel in `cnos.cdd`, per-protocol procedures here — inherits the option (a) decision from [cnos#388](https://github.com/usurobor/cnos/issues/388) and is the same split that produced [`cnos.cdr`](../cnos.cdr/) v0.1 at [cnos#376](https://github.com/usurobor/cnos/issues/376).

## Package Structure

This package is **v0.1 complete**. The doctrine surface (`CDS.md` — the 3,588-line canonical instantiation contract), the loader skill (`SKILL.md`), the extraction map (`docs/extraction-map.md`), the operational sub-area overlays (`skills/cds/lifecycle/SKILL.md` + `skills/cds/selection/SKILL.md`), and the empirical-anchor doc (`docs/empirical-anchor-cdd.md`) are all landed under [cnos#403](https://github.com/usurobor/cnos/issues/403). The wave bootstrapped under Sub 1 ([cnos#406](https://github.com/usurobor/cnos/issues/406)), authored CDS.md under Sub 2 ([cnos#407](https://github.com/usurobor/cnos/issues/407)), migrated content from `CDD.md`'s "pending cds extraction" section across Subs 3–5 ([cnos#408](https://github.com/usurobor/cnos/issues/408), [cnos#409](https://github.com/usurobor/cnos/issues/409), [cnos#410](https://github.com/usurobor/cnos/issues/410)), swept the `CDD.md` markers under Sub 6 ([cnos#411](https://github.com/usurobor/cnos/issues/411)), and closed with the empirical-anchor doc under Sub 7 ([cnos#412](https://github.com/usurobor/cnos/issues/412)) on 2026-05-22.

### `/skills/cds/` — The Method

- **`SKILL.md`** — Loader skill. Names `CDS.md` and the per-role overlay positions as advisory targets; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel) and `cnos.cdr` (sibling research realization). **Landed via Sub 1 of [cnos#403](https://github.com/usurobor/cnos/issues/403) — [cnos#406](https://github.com/usurobor/cnos/issues/406).**
- **`CDS.md`** — Canonical instantiation contract (3,588 lines). Declares the six fields per [`ROLES.md §3`](../../../ROLES.md): matter type (software-class), review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule. The §Selection function, §Development lifecycle, §Coordination surfaces, §Artifact contract, §Mechanical / §Review / §Gate / §Assessment / §Closure / §Retro / §Non-goals / §Large-file sections together are the canonical software-cycle doctrine. **Landed via Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403) — [cnos#407](https://github.com/usurobor/cnos/issues/407).** Migrated content from `cnos.cdd/skills/cdd/CDD.md` "pending cds extraction" section across Subs 3–5 (see "Surfaces" below).

### `/skills/cds/<area>/` — Operational sub-area overlays (v0.1 thin)

- **`lifecycle/SKILL.md`** — CDS development lifecycle thin overlay. Canonical 0–13 steps + S0–S13 state machine + branch rule + pre-flight + tier structure live at [`CDS.md §"Development lifecycle"`](skills/cds/CDS.md); this overlay delegates mechanics to existing `cnos.cdd` role + harness skills until the v1 role rewrite. **Landed via Subs 3–5 of [cnos#403](https://github.com/usurobor/cnos/issues/403)** as the cross-cutting operational binding for the lifecycle migration.
- **`selection/SKILL.md`** — CDS selection-function thin overlay. Canonical rules at [`CDS.md §"Selection function"`](skills/cds/CDS.md); mechanics delegate to `cnos.cdd/skills/cdd/gamma/SKILL.md` until the v1 role rewrite. **Landed via Sub 3 of [cnos#403](https://github.com/usurobor/cnos/issues/403) — [cnos#408](https://github.com/usurobor/cnos/issues/408).**

### `/docs/` — The Extraction Plan + Empirical Anchor

- **`extraction-map.md`** — Surface-by-surface migration plan from the pre-#402 software-lifecycle content quarantined in [`cnos.cdd/skills/cdd/CDD.md` §"Software-specific realization — pending cds extraction"](../cnos.cdd/skills/cdd/CDD.md) onto its planned CDS canonical home, with the migration sub-issue named per surface. Every row carries a "Status" field declaring the v0.1 outcome (migrated / deferred-by-design / complete). **Authored under Sub 1 of [cnos#403](https://github.com/usurobor/cnos/issues/403) — [cnos#406](https://github.com/usurobor/cnos/issues/406); finalised under Sub 6 — [cnos#411](https://github.com/usurobor/cnos/issues/411)** (every "pending cds extraction" marker resolved; cross-references in 11 cdd skill files + 2 YAML templates re-pointed at CDS canonical homes).
- **`empirical-anchor-cdd.md`** — Surface-by-surface mapping of the current cnos repository's `.cdd/` cycle artifacts to CDS's six-field structure, demonstrating that CDS v0.1 can host current practice without contradiction. The structural precedent is [`cnos.cdr/docs/empirical-anchor-cph.md`](../cnos.cdr/docs/empirical-anchor-cph.md). **Landed via Sub 7 of [cnos#403](https://github.com/usurobor/cnos/issues/403) — [cnos#412](https://github.com/usurobor/cnos/issues/412).**

### Surfaces (v0.1 Landed)

The v0.1 wave shipped, per the extraction map's Sub 3/4/5 split, as **canonical CDS.md sections plus the two operational sub-area overlays above** (lifecycle, selection). Per-role overlays at `skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` are **deferred to v1**; in v0.1 the role-cell discipline continues to be served by `cnos.cdd/skills/cdd/{alpha,beta,gamma,operator,epsilon}/SKILL.md` (cited by reference from `CDS.md`). The Sub 3/4/5 split landed as:

- **Sub 3 — [cnos#408](https://github.com/usurobor/cnos/issues/408):** §Selection function + §Development lifecycle migrated into `CDS.md` as canonical sections; the operational thin overlay `skills/cds/selection/SKILL.md` shipped alongside (B-lite thin extract).
- **Sub 4 — [cnos#409](https://github.com/usurobor/cnos/issues/409):** §Coordination surfaces + §Artifact contract migrated into `CDS.md` as canonical sections (B-lite thin extract); cross-cuts to §Inputs / §Tracking absorbed.
- **Sub 5 — [cnos#410](https://github.com/usurobor/cnos/issues/410):** §Mechanical / §Review / §Gate / §Assessment / §Closure / §Retro-packaging / §Non-goals / §Large-file migrated into `CDS.md` as canonical sections (B-lite thin extract).

### `/schemas/cds/` — Typed receipt surface

`schemas/cds/receipt.cue` already exists per [cnos#388](https://github.com/usurobor/cnos/issues/388); CDS Field 3 (γ close-out artifact) cites `#CDSReceipt` for the typed receipt shape. No schema-side migration was required by the #403 wave; the wave was doctrine-surface only.

## Quick Start

**New to CDS?** Start with:

1. [`skills/cds/CDS.md`](skills/cds/CDS.md) — The canonical instantiation contract: six fields per `ROLES.md §3`, plus the full software-cycle doctrine (selection, lifecycle, artifacts, review, gate, closure, assessment, retro, non-goals).
2. [`skills/cds/SKILL.md`](skills/cds/SKILL.md) — The loader skill; declares the load order and the cross-protocol relationship with `cnos.cdd` (kernel), `cnos.cdr` (research sibling), and `cnos.handoff` (wire-format peer).
3. [`docs/empirical-anchor-cdd.md`](docs/empirical-anchor-cdd.md) — Demonstrates that current cnos `.cdd/` cycle practice maps onto CDS v0.1's six-field structure without contradiction; useful for grounding the doctrine in observable evidence.
4. [`docs/extraction-map.md`](docs/extraction-map.md) — Per-surface migration map with v0.1-complete status; archaeological context on what moved from `CDD.md`'s "pending cds extraction" section to which `CDS.md` section under which Sub.
5. [`../cnos.cdd/skills/cdd/CDD.md`](../cnos.cdd/skills/cdd/CDD.md) — The CCNF spine (160 lines post-#402). CDS instantiates this kernel; reading the kernel first grounds the realization.
6. [`COHERENCE-CELL-NORMAL-FORM.md`](../cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md) and [`COHERENCE-CELL.md`](../cnos.cdd/skills/cdd/COHERENCE-CELL.md) — the kernel doctrine CDS cites by reference and does not restate.
7. [`schemas/cds/receipt.cue`](../../../schemas/cds/receipt.cue) — the typed γ close-out surface (`#CDSReceipt`) that CDS Field 3 references.
8. [cnos#403](https://github.com/usurobor/cnos/issues/403) — parent tracker (closed 2026-05-22 via [cnos#411](https://github.com/usurobor/cnos/issues/411)); rationale for the extraction; sub-issue wave history.

**Looking for the kernel?** See [`../cnos.cdd/`](../cnos.cdd/) — the generic recursive coherence-cell algorithm.

**Looking for the research-side sibling?** See [`../cnos.cdr/`](../cnos.cdr/) — the research realization shipped at [cnos#376](https://github.com/usurobor/cnos/issues/376). CDR is **not** CDS-with-different-words: engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*. The role names are shared; the loss functions diverge; the receipts and oracles diverge to enforce the divergent disciplines.

**Looking for the wire-format peer?** See [`../cnos.handoff/`](../cnos.handoff/) — the inter-agent / inter-activation / inter-repo handoff doctrine; CDS consumes handoff surfaces (cross-repo, dispatch, mid-flight, artifact-channel, receipt-stream) without re-deriving them.

## Cross-protocol relationship

CDS shares the role-cell kernel with [`cnos.cdd`](../cnos.cdd/) (the generic algorithm) and is a sibling of [`cnos.cdr`](../cnos.cdr/) (the research realization). The kernel — `COHERENCE-CELL`, `COHERENCE-CELL-NORMAL-FORM`, role grammar, generic receipt-validation interface — lives in `cnos.cdd`. CDS owns the **software-specific lifecycle, artifact contract, review CLP, gate, closure, assessment, retro-packaging, and non-goals.** No generic CCNF kernel doctrine is restated in CDS — only pointers from CDS to CDD/CCNF. The pointer discipline is the same one [`cnos.cdr/skills/cdr/SKILL.md`](../cnos.cdr/skills/cdr/SKILL.md) uses for its kernel citations.

**Boundary vs cnos.handoff:** handoff is the *wire format* (cross-repo bundles, dispatch envelopes, mid-flight rescue, artifact-channel write-ownership rules, cross-cycle receipt streams); CDS is the *software-cycle lifecycle*. CDS consumes handoff surfaces by reference; it does not own them and does not re-derive them. The two have different reasons to change.

## Status

**v0.1 complete — cnos#403 wave closed 2026-05-22.** All seven subs landed; the canonical `CDS.md` doctrine surface is in place; the two operational sub-area overlays (`lifecycle/`, `selection/`) are in place; the extraction map is finalised; the empirical-anchor doc is in place. Tracker closed via [cnos#411](https://github.com/usurobor/cnos/issues/411). The wave shipped, in order:

- `cn.package.json` (schema `cn.package.v1`; name `cnos.cds`; version `0.1.0`; kind `package`; engines `cnos >= 3.81.0`), `README.md`, `skills/cds/SKILL.md` (loader, advisory targets), `docs/extraction-map.md` (the migration plan Subs 3–5 dispatch against) — Sub 1 / [cnos#406](https://github.com/usurobor/cnos/issues/406);
- `skills/cds/CDS.md` (the 3,588-line six-field software instantiation contract; declares matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule) — Sub 2 / [cnos#407](https://github.com/usurobor/cnos/issues/407);
- §Selection function + §Development lifecycle migrated into `CDS.md` as canonical sections; `skills/cds/selection/SKILL.md` operational thin overlay shipped (B-lite extract) — Sub 3 / [cnos#408](https://github.com/usurobor/cnos/issues/408);
- §Coordination surfaces + §Artifact contract migrated into `CDS.md` as canonical sections (B-lite extract) — Sub 4 / [cnos#409](https://github.com/usurobor/cnos/issues/409);
- §Mechanical / §Review / §Gate / §Assessment / §Closure / §Retro-packaging / §Non-goals / §Large-file migrated into `CDS.md` as canonical sections; `skills/cds/lifecycle/SKILL.md` operational thin overlay shipped (B-lite extract) — Sub 5 / [cnos#410](https://github.com/usurobor/cnos/issues/410);
- `CDD.md` "pending cds extraction" markers replaced with CDS pointers; cross-references in 11 cdd skill files + 2 YAML templates re-pointed at CDS canonical homes — Sub 6 / [cnos#411](https://github.com/usurobor/cnos/issues/411);
- `docs/empirical-anchor-cdd.md` mapping current `.cdd/` practice onto the CDS six-field structure (precedent: [`cnos.cdr/docs/empirical-anchor-cph.md`](../cnos.cdr/docs/empirical-anchor-cph.md)) — Sub 7 / [cnos#412](https://github.com/usurobor/cnos/issues/412).

Per-role overlays at `skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` are **deferred to v1** (post-v0.1 field-evidence cycle). In v0.1 the role-cell discipline continues to be served by the existing `cnos.cdd/skills/cdd/{alpha,beta,gamma,operator,epsilon}/SKILL.md` overlays, cited by reference from `CDS.md`. No schemas were lifted beyond the pre-existing `schemas/cds/receipt.cue` (per [cnos#388](https://github.com/usurobor/cnos/issues/388)). The package is loadable and discoverable by cnos package-discovery (per `pkg.ContentClasses` filesystem-presence convention); commands are not declared in this version.

## License

Part of the cnos project.
