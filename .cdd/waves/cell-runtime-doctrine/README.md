# Wave: cell-runtime-doctrine (cnos#671 — R1)

**Planning Cell output.** This directory is the matter of the Planning Cell #671 (child of parent
wave #627): a mature, executable `cn.wave.v1` plan that decomposes the cell-runtime doctrine into
single-purpose Working-Cell contracts, grounded in an immutable coherence measurement. It is the
**R1 repair** of the draft #671 plan against the external-β ITERATE (PR/issue #671, comment
`5015777442`).

**This is a plan, not doctrine.** The `docs/` artifacts named as `requested_output` paths are the
**future output of the Working Cells** — they are *not* authored here. This wave authors the
*contracts* that will produce them. Nothing under `docs/` or `schemas/` is edited by this cell.

## Goal

Produce a coherent, mechanically-checkable specification of the cnos cell-runtime doctrine — the
typed cell contract, the coherence-measurement (CM) backbone, the cell/wave FSMs, and the
shipped→specified migration — decomposed so each artifact is single-purpose and individually
coherent. **Done means:** a Working Cell can implement the runtime (Go/CUE) from these artifacts
**without inventing** dependency, measurement, authority, or gate policy.

## Grounding CM (content-bound)

The whole wave is **grounded in** an immutable measurement, not "producing CM as a telos."

- **File:** [`grounding-cm.md`](./grounding-cm.md) — captures the external CC's recursive TSC
  measurement (`coh 0.12.0`, L0–L4) verbatim, plus the FAIL-classification table.
- **SHA-256 (authoritative binding):** `585955d8cdd4aa0556aa0daa0f955a00d1fdf69921c9cadea2abf25b620aebd3`
- **Source:** PR #667 comment `5015460988` (*External CC Ratification — Corrective recursive
  measurement addendum, frozen R7 only*), disposition `request_planning`.
- Because the comment is editable under its id, every node content-binds the **file hash**, not the
  comment id.

**CM ownership (the repaired BLOCKER).** CM is a **runner/provider-produced measurement object**
(producer = a CM provider/runtime op, e.g. tsc — **not** CC). A **CC's output is a JUDGMENT that
consumes immutable CM refs**; CM and CC-judgment are **distinct tagged objects**. **V gates on CM.**
A **PC's wave is grounded in / references an immutable CM.** WC-2 encodes exactly this; nothing here
says "CC produces CM as its judgment."

## Corrected dependency graph (WC-2 keystone root)

```
external roots: grounding-CM(@sha) + #628/S1 + shipped schemas/CCNF/transitions.json
  → WC-2  (CM measurement object + receipt_core→CM→V→δ→final_receipt type path)   [keystone root]
  → WC-1  (cell classes + typed cell contract; imports cm_ref from WC-2)
      → WC-3a (cell FSM) · WC-3b (wave FSM) · WC-4 (shipped→specified migration)   [each dep WC-1 + WC-2]
          → WC-5  (integration / seal: whole-wave proof)   [dep WC-1,2,3a,3b,4]
```

Every `cm_ref` imports from **WC-2's** stable output — no hidden `cm_ref` edge, no "light
coordination" prose as an edge. Edges are the mechanical projection of the child contracts'
`sibling_output` refs (see [`wave.cn-wave-v1.yaml`](./wave.cn-wave-v1.yaml) `edge_derivation`).
Critical path: **WC-2 → WC-1 → WC-3b → WC-5**.

## Node list (derived edges)

| Node | Class / domain | Output (`docs/…`) | Depends on (sibling_output edges) |
|---|---|---|---|
| **WC-2** *(keystone root)* | working / doctrine | `COHERENCE-MEASUREMENT.md` | — (external roots only) |
| **WC-1** | working / doctrine | `CELL-RUNTIME-CLASSES.md` | WC-2 |
| **WC-3a** | working / doctrine | `CELL-FSM.md` | WC-1, WC-2 |
| **WC-3b** | working / doctrine | `WAVE-FSM.md` | WC-1, WC-2 |
| **WC-4** | working / doctrine | `CONTRACT-MIGRATION.md` | WC-1, WC-2 |
| **WC-5** *(seal)* | working / doctrine | `CELL-RUNTIME-INTEGRATION.md` | WC-1, WC-2, WC-3a, WC-3b, WC-4 |

All nodes are `class: working, matter_domain: doctrine` (each produces a doctrine artifact that
lives in `docs/`, per the settled rule **WC = an artifact that lives in the repo**).

## How we get there (process)

1. **This PC pass** produces the wave (this plan) + carries the grounding CM → **operator review**
   → **external β** reviews the plan (this matter) → **CC** ratifies the plan (scope n+1) →
   **operator authorizes the wave** at its boundary. *(wave-authorization gate — once, not per-child)*
2. Each **WC** then runs its own coherence loop on its own surface: α → external β (exact-SHA) → γ
   (content-bound review) → CC ratification → operator final-read → merge. Each artifact is small
   and single-purpose, so it should converge faster than #662 did.
3. **WC-5** seals: it proves the six artifacts compose (cross-artifact vocabulary, schema fixtures,
   graph parity, migration round-trips, the #627 map, residual debt). Four/six child closures do
   **not** by themselves prove the objective.

## Reclassification (#662 → WC-1)

`#662` ("PC-D0") is **retired as a PC**. Its converged artifact
(`CELL-RUNTIME-CLASSES.md` — eight rounds: typed contracts, κ≠α-unconditional role logic,
content-bound evidence) **re-lands as WC-1**, re-scoped and correctly classified as WC-doctrine
(its canonical output is a repo artifact — WC telos, per finding L0-B1). **No work is discarded.**
WC-1 imports the `cm_ref` interface from WC-2 rather than re-deriving CM. The κ≠α role logic and the
content-bound-review rule repaired in #662 R8 are carried into WC-1 as **settled input** (the two
`retired-instance-defect (#662)` rows in `grounding-cm.md`).

## Per-finding disposition (external-β ITERATE, comment 5015777442)

| # | Finding | Repair in this R1 |
|---|---|---|
| 1 | **[BLOCKER]** CM incorrectly owned by CC judgment | WC-2 defines CM as a runner/provider-produced object (producer = CM provider/runtime op, not CC); the CC result is a separate tagged object consuming immutable CM refs; pins producer/inputs/identity/provenance/standing, the CM→V edge, and the `receipt_core → CM → V → δ → final_receipt` type path resolving the receipt.cue circularity. |
| 2 | **[BLOCKER]** Claimed child contracts don't exist | Six complete `cn.cell.contract.v1` bodies authored (`contracts/wc-1..wc-5`), each pinning cell/inputs/output/acceptance/constraints/gates/doctrine_affecting/stop_conditions/consumers/completion; plus the `cn.wave.v1` instance with mechanically-derived edges. |
| 3 | **[REQUIRED]** hidden `cm_ref` edge / roots non-independent / WC-4 missing dep | WC-2 is the sole keystone root; WC-1 depends on WC-2; WC-3a/WC-3b/WC-4 each depend on WC-1+WC-2. Every `cm_ref` imports from WC-2's output. No prose edges. |
| 4 | **[REQUIRED]** split the FSM node | WC-3 → **WC-3a (cell FSM)** + **WC-3b (wave FSM)** — different state spaces/events/terminals, separate authorities; each pins the CC-disposition→transition-request adapter. |
| 5 | **[REQUIRED]** author the #627 refinement table | [`reconcile-627.md`](./reconcile-627.md): every S1–S8 node classified exactly once (landed-predecessor / refined-by / unchanged-downstream / superseded); no double owner; the WC wave FEEDS S2–S8, does not fork it. |
| 6 | **[REQUIRED]** content-bind grounding + add the seal | `grounding-cm.md` content-binds the measurement (SHA above) with the FAIL-classification table; **WC-5 (integration/seal)** is the whole-wave proof node. |

## Files

- [`README.md`](./README.md) — this overview.
- [`wave.cn-wave-v1.yaml`](./wave.cn-wave-v1.yaml) — the `cn.wave.v1` instance (nodes, derived edges, gates, completion predicate).
- [`contracts/wc-1..wc-5.cn-cell-contract-v1.yaml`](./contracts/) — the six complete child contracts (WC-1, WC-2, WC-3a, WC-3b, WC-4, WC-5).
- [`grounding-cm.md`](./grounding-cm.md) — captured CC measurement + SHA-256 + FAIL-classification table.
- [`reconcile-627.md`](./reconcile-627.md) — the #627 S1–S8 refinement/supersession table.

---
*Status: R1 wave, α matter — under operator review; external-β and CC review next. No child WCs
dispatched; no control-plane action taken by this cell.*
