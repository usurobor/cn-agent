# α self-coherence — #671 (Planning Cell — cell-runtime doctrine wave)

> **Nature of this artifact.** This is the constitution + role-provenance
> ledger for the #671 Planning Cell, produced to make the cell's roles
> **observable in git** rather than only as PR comments. It is authored
> under an explicit **bootstrap protocol exemption** (see #671 body
> `## Protocol exemption`, and `review/SKILL.md` §3.11b recovery path (b)).
> It records honestly what happened, including where roles were collapsed.
> It does **not** mutate the reviewed matter under
> `.cdd/waves/cell-runtime-doctrine/`, which stays frozen at the
> externally-reviewed revision.

## Issue

[#671 — PC — plan the cell-runtime doctrine wave (CM-grounded)](https://github.com/usurobor/cnos/issues/671).
Parent wave #627. This is a **Planning Cell (PC)**: its telos is a
`cn.wave.v1` relation graph (six Working-Cell contracts) that guides
future WC work, not a `docs/` artifact and not a release. The live matter
is **PR #672**, branch `wave/671-cell-runtime-doctrine`.

## Mode

**Bootstrap Planning Cell, §5.2-adjacent (wave-producing), doctrine.**
This cell *authors the cell-runtime doctrine itself* — the classes, the
contract envelope, the FSMs, the migration, the seal — which is the very
machinery that, once implemented (#627 S2–S8), will mechanically
constitute and separate cell roles. The cell therefore **predates the
runtime that would constitute it**. That is the bootstrap paradox this
exemption names rather than papers over.

## Gap

The measured incoherence this cell reduces (grounding CM, content-bound
`sha256:9d1ab3a5…` capture of `cnos#667#issuecomment-5015460988`):

> The cell-runtime doctrine existed as scattered, partially-contradictory
> surfaces (#628/S1 runtime loop, #662 cell classes, CELL-KINDS.md,
> shipped `schemas/cdd/*.cue`) with no single mechanically-checkable
> decomposition a Working Cell could execute without inventing
> dependency, measurement, authority, gate, migration, or closure policy.

The cell closes the gap by producing one revision-bound, content-addressed
graph of six executable `cn.cell.contract.v1` Working-Cell contracts
(WC-2 → WC-1 → {WC-3a · WC-3b · WC-4} → WC-5), each grounded in the
measured incoherence, with a total assurance registry and a materialized
wave-boundary pre-authorization gate.

## ACs (cell telos → evidence)

| AC | Claim | Evidence (reviewed matter @ `614829a4`) |
|---|---|---|
| 1 | Six-node `cn.wave.v1` construction graph, single keystone root WC-2 | `.cdd/waves/cell-runtime-doctrine/wave.cn-wave-v1.yaml`; external-β independently derived 12 edges = 12 authored, sole root WC-2 |
| 2 | Six complete `cn.cell.contract.v1` Working-Cell contracts | `contracts/wc-{1,2,3a,3b,4,5}.cn-cell-contract-v1.yaml`; each carries inputs/output/surface/acceptance/STOP/gate/evidence-ownership |
| 3 | Structural assurance (CUE) | `schema/*.cue` + 31 negative regressions; `make -C schema all` → exit 0, all 31 negatives reject |
| 4 | Wave-boundary pre-authorization gate (Go, CUE-normalized) | `wave-validators/oracle_ownership_bijection.go`; `make -C wave-validators all` → exit 0; 78 ⇄ 78 ownership bijection, 30 mechanically-verifiable |
| 5 | Total assurance registry, every child predicate classified once | `oracle-registry.yaml`; whole-wave properties owned at the wave boundary |
| 6 | Content-addressed grounding + honest transitional intent | `grounding-cm.md`, `grounding-source-5015460988.md`, `intent.cn-intent-v1.yaml`, `decision-provenance.md` |
| 7 | #627 reconciliation (doctrine layer feeds S2–S8, not a fork) | `reconcile-627.md` |

## CDD Trace

- **α (matter production):** performed by **Sigma-at-repo** across rounds
  R1→R16 under bootstrap role-collapse (see `## Known debt` item 1). The
  matter is the wave under `.cdd/waves/cell-runtime-doctrine/`, frozen at
  the reviewed revision.
- **β (independent review):** performed by a **genuinely external
  reviewer of different model lineage** ("Codex", posting as `usurobor`),
  content-bound to the exact reviewed SHA each round. Terminal verdict
  **CONVERGE at R16 / `614829a4`** — no BLOCKER or REQUIRED defect.
  Captured in `beta-review.md` (this cycle dir).
- **γ (closure):** performed by a **distinct non-κ actor** — see
  `gamma-closeout.md`. The earlier κ-signed γ closeout
  ([PR #672 comment 5076124652](https://github.com/usurobor/cnos/pull/672#issuecomment-5076124652))
  is **retracted as void** (κ must not perform γ; firebreak). This is the
  load-bearing repair of the CC HOLD.
- **CC (process judgment):** returned **HOLD** on the κ-signed γ + missing
  constitution. This artifact set + the non-κ γ + the exemption are the
  response; a fresh external-CC judgment runs against the new head before
  operator authorization.
- **κ (control-plane / Herald):** Sigma-at-repo coordinated the review,
  pushed revisions, and captured evidence. κ is **outside** the cell and
  authors no in-cell judgment.

## Self-check

| Row | State | Evidence |
|---|---|---|
| Matter frozen | yes | `.cdd/waves/cell-runtime-doctrine/**` byte-identical to `614829a4`; this cycle adds only `.cdd/unreleased/671/**` + the #671 body exemption |
| β evidence content-bound | yes | `beta-review.md` binds comment id + URL + reviewed SHA + verbatim body |
| γ not by κ | yes (via `gamma-closeout.md`) | authored + committed by a distinct non-κ actor; void κ-signed closeout retracted on PR |
| 3.11b discoverable | yes | `## Protocol exemption` in #671 body (recovery path (b)) |
| firebreak κ≠α stated | yes (as bootstrap debt) | `## Known debt` item 1; named, not hidden |
| No witness theater | yes | no fabricated multi-actor history; collapse recorded truthfully |

## Known debt

1. **Bootstrap role-collapse (α/κ/δ fused in one Sigma session).** Per
   `delta/SKILL.md` §9.1, bootstrap-δ legitimately collapses γ-driver with
   δ-orchestrator and spawns roles as sub-sessions; this cell went further
   and fused α authorship with the κ/δ control plane. This is an
   **empirical bootstrap case, not the target architecture** — the target
   is exactly what this wave specifies (the generic cell runner of #627,
   which mechanically separates the roles). Debt is discharged structurally
   by shipping that runtime; named here so the collapse is auditable.
2. **γ was initially performed by κ — a firebreak violation.** Corrected:
   the void closeout is retracted and γ is re-performed by a distinct
   non-κ actor (`gamma-closeout.md`). The general fix is the FSM/role
   separation the wave itself specifies (WC-3a/WC-3b).
3. **README historical tail (`README.md:408-411`) names R13→R14 as the
   prior round** — a stale historical projection. External-β left it as a
   deliberate OBSERVATION ("do not mutate the converged R16 matter"); it is
   carried as debt for a later authorized documentation pass, **not**
   repaired here (repairing it would break the freeze / seal).

## Review-readiness | bootstrap PC | reviewed matter SHA `614829a4682e148d98c70371e600ffdc3fa6386e` (R16) | matter frozen; constitution + role receipts added under `.cdd/unreleased/671/` | ready for non-κ γ + external-CC
