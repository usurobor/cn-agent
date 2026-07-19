# Grounding CM — α's derivative FAIL-classification (cnos#671)

**This file is a DERIVATIVE, authored by α. It is NOT the exact bytes of the source measurement.**
The exact, byte-for-byte source is captured immutably in a separate snapshot:

- **Immutable source snapshot:** [`grounding-source-5015460988.md`](./grounding-source-5015460988.md)
- **True source SHA-256:** `9d1ab3a50d00e642fdeb87728dd71f7c7499c60878afe7001f2ddb832b161dbb`
- **Source object:** `usurobor/cnos` PR #667, issue-comment id `5015460988`
  (*External CC Ratification — cnos#662 (PC-D0) · Corrective recursive measurement addendum — frozen
  R7 only*), created `2026-07-19T11:01:36Z`, disposition `request_planning`.

The measurement that grounds this wave is the source snapshot above. This file only carries **α's
derived FAIL-classification table** — the trace mapping each measured FAIL to the wave node that
disposes it. Every claim below is α's reading of the snapshot; it is **not** a verbatim reproduction
of the source, and it must never be cited as "the exact bytes." Anything needing the exact source
bytes MUST resolve the snapshot by its path + source SHA-256 above.

> **Prior-revision defect (repaired here).** An earlier revision of this file inlined a fenced block
> and called it "the exact bytes of comment 5015460988." That block was **abridged** (the true
> source is 150 lines, SHA-256 `9d1ab3a5…`), so the claim was false. The exact bytes now live only in
> the snapshot; this file is honestly labeled a derivative.

## CM ownership (the settled BLOCKER-1 doctrine)

The object measured in the source is a **runner/provider-produced measurement**: the external
Cohering Cell ran the released `coh 0.12.0` engine (`usurobor/tsc`) over frozen recursive
projections and emitted α/β/γ scores, a bottleneck, defects, and provenance. The CC's own class
output was a **separate `request_planning` JUDGMENT that consumes this CM by reference**; the CC
refused to treat its judgment as the CM. **CM and CC-judgment are distinct tagged objects. V gates on
CM. A PC's wave is grounded in / references an immutable CM.** WC-2 encodes exactly this; nothing
here collapses CM into CC-judgment.

## Provenance (immutable binding) — as read from the snapshot

| Field | Value (per source snapshot `9d1ab3a5…`) |
|---|---|
| Source control-plane comment | `usurobor/cnos` PR #667, issue-comment id `5015460988` (2026-07-19T11:01:36Z) |
| Producer | independent external Cohering Cell (TSC-γ judgment), outside the Sigma activation lineage |
| CM provider / runtime op | `coh 0.12.0 (016c511)` — released TSC engine (a CM provider, **not** the CC's judgment) |
| TSC source | `usurobor/tsc@26aab5023f03dc7d0abf82e5fdba20134fc6adad` |
| Measured matter | `docs/architecture/CELL-RUNTIME-CLASSES.md` at exact SHA `2d6b93cc4e69e5b413a80bd8e352cb0a004da460` (frozen #662 R7) |
| Measured matter bytes | SHA-256 `80e0d8c68a3d8affabdd4bd14848cbb3f0bee27b078e435f61433e42a0ee89e0` |
| Disposition (frozen R7) | `request_planning` — ratification withheld; feeds operator final-read gate |
| Standing note | mechanical arm N=3 exact-equal; semantic arm k=1 (no semantic-consistency standing); hard-invariant FAILs gate independently |

The authoritative binding is the **source snapshot SHA-256 `9d1ab3a5…`** (the exact bytes), recorded
in `wave.cn-wave-v1.yaml` and imported by every node's grounding `external` ref. The comment id is a
mutable carrier and is retained only as the human-readable source pointer.

---

## α's derived FAIL-classification table (every measured FAIL disposed exactly once)

Each measured hard-invariant FAIL and named finding in the snapshot is α-classified as one of:
**`closed-by <WC-n>`** (a wave node's doctrine artifact closes the measured gap),
**`retired-instance-defect (#662)`** (a defect of the frozen #662 *instance* already repaired in a
later round — not a doctrine gap the wave must close), or **`tracked-follow-up`** (real, out of this
wave's doctrine scope; routed downstream). This is α's mapping, derived from the snapshot — not a
reproduction of it.

| Measured FAIL / finding | Level | Classification | Rationale |
|---|---|---|---|
| Class ≠ canonical output telos | L0-B1 | **closed-by WC-1** | WC-1 fixes the class↔telos rule and re-lands #662's converged core correctly classified as WC-doctrine; the plan itself reclassifies #662 → WC-1. |
| Receipt-ordering circular (γ receipt before V/δ, but `receipt.cue` requires their outputs) | L1-A1 | **closed-by WC-2** | WC-2 declares the provisional-receipt→measurement→final-receipt split: the `receipt_core → CM → V → δ → final_receipt` type path. |
| κ asserted both `≠α` and `=α` | L1-A2 | **retired-instance-defect (#662)** | This is CC-1 (κ-logic); repaired in #662 R8 (`κ≠α` unconditional; State A = hosting-identity collapse). An instance defect, not a doctrine gap. WC-1 carries the repaired role logic as settled input. |
| CM contract missing (`matter+receipt → CM → V` not typed; no `cm_ref`) | L2-B1 | **closed-by WC-2** | WC-2 is the keystone: CM as a runner/provider-produced first-class typed object with producer/inputs/identity/provenance/standing and the exact CM→V edge, realized **within** the four-schema boundary. |
| Shipped↔specified contracts have no adapter | L2-B2 | **closed-by WC-4** | WC-4 authors the versioned `cnos.cdd.contract.v1` ↔ `cn.cell.contract.v1` projection with preservation rules + negative fixtures. |
| γ content-binds matter but not review bytes | L2-G1 | **retired-instance-defect (#662)** | This is CC-2 (review-binding); repaired in #662 R8 (§11.6 requires content-bound review artifact). The doctrine rule already exists in the #662 core WC-1 re-lands; instance defect closed. |
| Cell lifecycle authority split (`transitions.json`/`dispatch`/`resume`/spec; `blocked` off the state set) | L3-A1/B1 | **closed-by WC-3a** | WC-3a authors one total cell FSM with exactly one authority per edge and `blocked` in the declared state set. |
| Wave "FSM" is a state sketch, not executable | L3-G1 | **closed-by WC-3b** | WC-3b authors a total, independently-total wave FSM incl. child-receipt aggregation and the CC-disposition→transition-request adapter. |
| CC-1 (κ role-logic) — earlier CC blocker | (=L1-A2) | **retired-instance-defect (#662)** | Same as L1-A2; accepted into #662 R8. |
| CC-2 (review evidence identity-bound) — earlier CC blocker | (=L2-G1) | **retired-instance-defect (#662)** | Same as L2-G1; accepted into #662 R8. |
| Semantic arm k=1, no consistency standing; TSC Operational ACCEPT/REJECT not claimed | (standing) | **tracked-follow-up** | A full-standing recursive re-measure (N≥30, bootstrap CI) belongs to a later CM lifecycle; out of this doctrine wave's scope. |
| Derived `Recursive Cell Coherence Methodology v0.1.0` has no standing / TSC cannot yet consume arbitrary CM declarations first-class | (candidate) | **tracked-follow-up** | Candidate input to the #627 S3 CM-engine implementation; not warrant here. WC-2 cites it as a design reference only. |

**Four-schema boundary (settled).** The snapshot's hard-invariant check records the four-schema
boundary as **PASS**. Per operator decision, D9's four-schema boundary stays **SETTLED**: CM is
realized as a typed object/refinement **within** the four schemas (a typed CM field/edge in the
receipt + a `cm_ref` resolving within the existing schemas), **not** as a fifth canonical `cn.cm.v1`
schema. WC-2 encodes this; WC-1 and WC-4 consume the same `cm_ref` shape (see WC-2/WC-1/WC-4
acceptance predicates and the reject-a-fifth-schema negative fixture).

**Whole-wave note.** Closing the individual FAILs above does **not** by itself prove the objective:
**WC-5 (integration / seal)** is the whole-wave proof (cross-artifact vocabulary, schema fixtures,
graph parity, migration round-trips, the #627 mapping, and residual-debt disposition). The
classification here is the per-finding trace WC-5 re-checks for parity.
