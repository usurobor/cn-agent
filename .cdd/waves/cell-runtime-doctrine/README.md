# Wave: cell-runtime-doctrine (cnos#671 — R2)

**Planning Cell output.** This directory is the matter of the Planning Cell #671 (child of parent
wave #627): a mature, executable `cn.wave.v1` plan that decomposes the cell-runtime doctrine into
single-purpose Working-Cell contracts, grounded in an immutable coherence measurement. It is the
**R2 repair** of the wave against the external-β ITERATE (PR/issue #672): six exact-contract +
assurance repairs. The accepted decomposition graph shape is **unchanged**
(WC-2 root → WC-1 → {WC-3a, WC-3b, WC-4} → WC-5); R2 repairs the contracts and adds assurance.

**This is a plan, not doctrine.** The `docs/` artifacts named as `requested_output` paths are the
**future output of the Working Cells** — they are *not* authored here. This wave authors the
*contracts* that will produce them. Nothing under `docs/` or `schemas/` is edited by this cell.

## Goal

Produce a coherent, mechanically-checkable specification of the cnos cell-runtime doctrine — the
typed cell contract, the coherence-measurement (CM) backbone, the cell/wave FSMs, and the
shipped→specified migration — decomposed so each artifact is single-purpose and individually
coherent. **Done means:** a Working Cell can implement the runtime (Go/CUE) from these artifacts
**without inventing** dependency, measurement, authority, or gate policy.

## Grounding CM (content-bound to the TRUE source)

The whole wave is **grounded in** an immutable measurement, not "producing CM as a telos."

- **Immutable source snapshot (exact bytes):** [`grounding-source-5015460988.md`](./grounding-source-5015460988.md)
  — byte-for-byte capture of PR #667 comment `5015460988` (the external-CC recursive TSC measurement,
  frozen R7).
- **True source SHA-256:** `9d1ab3a50d00e642fdeb87728dd71f7c7499c60878afe7001f2ddb832b161dbb` (150-line body).
- **α's derivative:** [`grounding-cm.md`](./grounding-cm.md) — clearly marked as α's **derived**
  FAIL-classification mapping that *references* the snapshot by path + source hash. It is **not** the
  exact bytes (an earlier revision wrongly inlined an abridged block and called it verbatim; repaired).
- `wave.cn-wave-v1.yaml` grounding points to **both**: the source snapshot (by source hash) and α's
  derivative (by content hash). Every node's grounding `external` ref content-binds the **source hash**.

**CM ownership.** CM is a **runner/provider-produced measurement object** (producer = a CM
provider/runtime op, e.g. tsc — **not** CC). A **CC's output is a JUDGMENT that consumes immutable CM
refs**; CM and CC-judgment are **distinct tagged objects**. **V gates on CM.** A **PC's wave is
grounded in / references an immutable CM.** WC-2 encodes exactly this.

## Intent (durable `cn.intent.v1`, materialized)

[`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) is the durable, first-class planning object
(id, source: operator, captured_by: kappa, statement, scope, constraints, desired_outcome) — SHA-256
`f3f484132c819c2a184df06ba99317bd032678f2cc6e47be74446f222005cb5c`. Per the D9 dependency order
(`cn.intent.v1 → cn.cell.contract.v1 → …`) it exists **before** any cell. Every child `intent_ref`
resolves to it by an **immutable, content-bound** locator (`sha256:<hash>@<path>`) through a typed
**bootstrap adapter** (`carrier.kind: repo_artifact_bootstrap`), **not** the mutable issue carrier
`cnos#671`. The §2 `intent_ref: {schema, id, carrier{kind, ref}}` key-path shape is preserved.

## Corrected dependency graph (WC-2 keystone root) — unchanged shape

```
external roots: grounding-CM(@sha 9d1ab3a5) + #628/S1 + shipped schemas/CCNF/transitions.json
  → WC-2  (CM measurement object + receipt_core→CM→V→δ→final_receipt type path)   [keystone root]
  → WC-1  (cell classes + typed cell contract; imports cm_ref from WC-2)
      → WC-3a (cell FSM) · WC-3b (wave FSM) · WC-4 (shipped→specified migration)   [each dep WC-1 + WC-2]
          → WC-5  (integration / seal: whole-wave proof)   [dep WC-1,2,3a,3b,4]
```

Edges are the mechanical projection of the child contracts' `sibling_output` refs (see
`wave.cn-wave-v1.yaml` `edge_derivation`); `validate.py` check (c) proves parity. Critical path:
**WC-2 → WC-1 → WC-3b → WC-5**.

## §2 conformance (the exact-contract repair)

All six `contracts/*.yaml` are canonical `cn.cell.contract.v1` §2 instances with the **exact**
top-level key set and **no extras**. The non-canonical `consumers` and `completion_signal` keys are
**removed**. Reverse-consumer information is now **derived** (see
[`reconcile-627.md`](./reconcile-627.md) "Derived reverse-consumers" — from `sibling_output` refs +
the #627 map). Child completion uses the **canonical derived rule** (requested output produced ·
acceptance predicates pass · V PASS · bound receipt), not a per-contract `completion_signal`.

## Non-recursive whole-wave completion (BLOCKER-1 repair)

`wave.cn-wave-v1.yaml` `completion` defines a **non-recursive construction set**
`N = {wc-1, wc-2, wc-3a, wc-3b, wc-4}` (excludes wc-5). **WC-5 readiness** = a completion receipt
exists for every node in N. **WC-5 completion** = wc-5's own output + acceptance + V PASS + bound
receipt (it does **not** quantify over itself or the wave predicate). **Whole-wave completion** =
`all(N complete) AND wc5_complete`. The completion-predicate dependency graph is **acyclic**
(`validate.py` check (g)). Four fixtures pin the behavior: all-N-incomplete (seal not ready);
all-N-complete/WC-5-absent (seal ready, wave incomplete); WC-5 FAIL (wave incomplete); all-N + WC-5
PASS (wave complete).

## D9 — four-schema boundary SETTLED (finding 6)

Per operator decision, D9's four-schema boundary stays **settled** (no authorized reopening). CM is a
typed object/refinement realized **within** the four schemas (a typed CM field/edge in the receipt +
a `cm_ref` resolving within the existing schemas), **not** a fifth canonical `cn.cm.v1` schema. WC-2's
contradictory 4-vs-5 choice predicate is removed; WC-1 and WC-4 consume the same in-four-schema
`cm_ref` shape. Negative fixture: a proposed fifth canonical schema is rejected while four stays
settled (WC-2 predicate + wave `stop_conditions.unresolved_d9_at_wave_scope`).

## Wave-level STOP conditions + immutable `contract_ref` resolution (finding 5)

`wave.cn-wave-v1.yaml` adds typed wave-level `stop_conditions` (stale grounding, failed predecessor,
invalid graph, revoked authorization, integration conflict, unresolved-D9-at-wave-scope) — each a
single typed non-dispatch transition (hold/replan/escalate) + receipt. `contract_ref_resolution`
states the rule: every `contract_ref` resolves **relative to the immutable authorized wave commit**;
operator authorization is **revision-bound**; re-resolution at a changed tree **fails stale**. Each
node also carries `contract_sha256` for independent content-addressable resolution.

## Pre-authorization validator + decidable oracles (finding 4)

- [`validate.py`](./validate.py) — credential-free (stdlib + PyYAML) Planning-Cell pre-authorization
  validator. Checks (a) §2 key-path shape, (b) immutable refs resolve incl. intent + grounding source
  hash, (c) output/ref edge parity, (d) DAG, (e) parallel nodes share no write surface, (f) gate
  invariants, (g) completion-predicate acyclicity. **Exits non-zero on any violation.** Passes at this
  wave tree (all seven checks green). Negative-fixture notes are in its docstring and verified
  deterministic.
- [`acceptance-oracles.md`](./acceptance-oracles.md) — for each load-bearing acceptance predicate
  (wave + WC-1..WC-5), a decidable oracle: fixture/schema/command, expected rejection cases, evidence.

## Node list (derived edges)

| Node | Class / domain | Output (`docs/…`) | Depends on (sibling_output edges) |
|---|---|---|---|
| **WC-2** *(keystone root)* | working / doctrine | `COHERENCE-MEASUREMENT.md` | — (external roots only) |
| **WC-1** | working / doctrine | `CELL-RUNTIME-CLASSES.md` | WC-2 |
| **WC-3a** | working / doctrine | `CELL-FSM.md` | WC-1, WC-2 |
| **WC-3b** | working / doctrine | `WAVE-FSM.md` | WC-1, WC-2 |
| **WC-4** | working / doctrine | `CONTRACT-MIGRATION.md` | WC-1, WC-2 |
| **WC-5** *(seal)* | working / doctrine | `CELL-RUNTIME-INTEGRATION.md` | WC-1, WC-2, WC-3a, WC-3b, WC-4 |

## Reclassification (#662 → WC-1)

`#662` ("PC-D0") is **retired as a PC**. Its converged artifact (`CELL-RUNTIME-CLASSES.md`) re-lands
as **WC-1**, re-scoped and correctly classified as WC-doctrine (its canonical output is a repo
artifact — WC telos, per finding L0-B1). **No work is discarded.** WC-1 imports the `cm_ref` interface
from WC-2 rather than re-deriving CM. The κ≠α role logic and content-bound-review rule repaired in
#662 R8 are carried into WC-1 as **settled input** (the two `retired-instance-defect (#662)` rows in
`grounding-cm.md`).

## Per-finding disposition (external-β ITERATE, PR #672)

| # | Finding | Repair in this R2 |
|---|---|---|
| 1 | **[BLOCKER]** whole-wave completion was recursive | Non-recursive `N = {wc-1,wc-2,wc-3a,wc-3b,wc-4}`; WC-5 readiness = receipts over N; WC-5 completion = own output/acceptance/V/receipt (no self/wave quantification); whole-wave = all(N) AND wc5_complete; predicate DAG proven acyclic (`validate.py:g`); 4 fixtures added. WC-5's recursive acceptance predicate removed. |
| 2 | **[BLOCKER]** child contracts not valid §2 instances | Removed non-canonical `consumers`/`completion_signal` and the open-schema `...` marker from all six; reverse-consumers **derived** in `reconcile-627.md`; materialized `intent.cn-intent-v1.yaml` (SHA recorded); every `intent_ref` resolves to it by an immutable content-bound locator via a typed `repo_artifact_bootstrap` adapter; all six normalize to §2 key paths (`validate.py:a`). |
| 3 | **[REQUIRED]** grounding called an abridgment "verbatim" | Fetched the exact comment body via the GitHub API, stored byte-for-byte as `grounding-source-5015460988.md` (true SHA-256 `9d1ab3a5…`); `grounding-cm.md` is now a clearly-marked derivative referencing the snapshot; wave grounding points to both. |
| 4 | **[REQUIRED]** acceptance not mechanically decidable; no validator | `acceptance-oracles.md` names a decidable oracle per load-bearing predicate; `validate.py` mechanically checks (a)–(g), exits non-zero on violation, passes here, negatives verified deterministic. |
| 5 | **[REQUIRED]** no wave STOP conditions; `contract_ref` not immutable | Added typed wave-level `stop_conditions`; `contract_ref_resolution` binds to the immutable authorized wave commit (revision-bound authorization; stale on changed tree); per-node `contract_sha256`. |
| 6 | **[REQUIRED]** WC-2 4-vs-5 schema contradiction | D9 four-schema boundary kept **settled**; CM realized **within** the four schemas (typed CM field/edge + `cm_ref`); 4-vs-5 choice predicate removed; WC-1/WC-4 consume the same shape; reject-a-fifth negative fixture. |

## Files

- [`README.md`](./README.md) — this overview.
- [`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) — durable `cn.intent.v1` planning object.
- [`wave.cn-wave-v1.yaml`](./wave.cn-wave-v1.yaml) — the `cn.wave.v1` instance (nodes, derived edges, gates, STOP conditions, `contract_ref` resolution, non-recursive completion).
- [`contracts/wc-1..wc-5.cn-cell-contract-v1.yaml`](./contracts/) — the six §2-conformant child contracts.
- [`grounding-source-5015460988.md`](./grounding-source-5015460988.md) — byte-exact source snapshot (SHA `9d1ab3a5…`).
- [`grounding-cm.md`](./grounding-cm.md) — α's derivative FAIL-classification (references the snapshot).
- [`reconcile-627.md`](./reconcile-627.md) — the #627 S1–S8 refinement map + derived reverse-consumers.
- [`acceptance-oracles.md`](./acceptance-oracles.md) — decidable oracle per load-bearing predicate.
- [`validate.py`](./validate.py) — pre-authorization validator (checks a–g; exits non-zero on violation).

---
*Status: R2 wave, α matter — under operator review; external-β and CC review next. No child WCs
dispatched; no control-plane action taken by this cell.*
