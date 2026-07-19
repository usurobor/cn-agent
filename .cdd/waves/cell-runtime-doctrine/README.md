<!-- wave-revision: R5 -->
# Wave: cell-runtime-doctrine (cnos#671 — R5)

**Planning Cell output.** This directory is the matter of the Planning Cell #671 (child of parent
wave #627): a mature, executable `cn.wave.v1` plan that decomposes the cell-runtime doctrine into
single-purpose Working-Cell contracts, grounded in an immutable coherence measurement. **R2** repaired
the external-β ITERATE (PR/issue #672); **R3** repaired a second external-β ITERATE (sound validator,
honest intent provenance, honestly-classified oracles); **R4** added per-WC oracle ownership. **R5**
(this revision) repairs how R4 added that ownership: R4 carried oracle ownership as a **non-canonical
top-level `oracles` key**, breaking the exact §2 contract shape, and its check (h) was not a true
bijection. R5 restores the exact §2 contracts, moves oracle ownership into a separate **content-bound
`oracle-registry.yaml`** each contract consumes via a canonical `inputs.required[].external` ref, makes
check (h) a **total + singular bijection** with adversarial coverage, and corrects the revision +
classification ledger. The accepted decomposition graph shape is **unchanged** across R2–R5
(WC-2 root → WC-1 → {WC-3a, WC-3b, WC-4} → WC-5).

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

## Intent (TRANSITIONAL bootstrap projection — honest provenance, finding-2 repair)

[`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) is an **explicitly transitional bootstrap
projection** of operator intent, authored **during** this cell — SHA-256
`1e246846774c43365cb00dc1b39b6998bab535f4833a16fcfe32ca3e2825560e`. It does **not** claim to have
existed before any cell: D9 **names** `cn.intent.v1` but does not implement it, and **no durable
pre-cell `cn.intent.v1` existed** in this manual bootstrap (the file first appeared mid-cell). Its
`statement` carries **only authoritative operator matter** — the objective as the operator directed it
in the #627 wave master, including the operator's **verbatim final doctrine line** — and its carriers
point to the real operator-intent sources (**#627**, the operator verbatim line). `captured_by: kappa`
records the intent-capture **role**, with an honest transitional note and a stated **supersession**
relation to a future durable object (#627 S2 / #644). Every **α-derived conclusion** (the WC graph, the
D9 realization detail, the β-repair outcomes) has been **moved out** of intent into
[`decision-provenance.md`](./decision-provenance.md), each citing its actual α/β/settled-doctrine
source. **Identity vs carrier are kept distinct:** identity is `intent_ref.id`; the child contracts'
`carrier.kind: repo_artifact_bootstrap` is a **projection pointer** to this transitional snapshot, not
the identity-authoring mechanism. `validate.py` (b) now checks each contract's `intent_ref.id`/`schema`
resolves to this object's `id`/`schema`. The §2 `intent_ref: {schema, id, carrier{kind, ref}}`
key-path shape is preserved.

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
top-level key set and **no extras**. The non-canonical `consumers`, `completion_signal`, and (R5) the
R4 `oracles` keys are **removed** — oracle ownership now lives in the separate content-bound
[`oracle-registry.yaml`](./oracle-registry.yaml), which each contract consumes via a canonical
`inputs.required[].external` control_plane ref (immutable locator + content hash). An added top-level
key (incl. `oracles`) now FAILS check (a). Reverse-consumer information is **derived** (see
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

## SOUND pre-authorization validator + negative fixtures (finding 1) + honest oracles (finding 3)

- [`validate.py`](./validate.py) — credential-free (stdlib + PyYAML + local `git`), **fail-closed and
  genuinely SOUND**. It derives every fact from the authored data (no hard-coded node/edge/predicate
  list) and checks: (a) the **full §2 constraint model** — key paths + **enums** (`cell.class`,
  `ref_kind`, locator `kind`, `requested_output.kind`) + types + cardinalities; (b) **real ref
  resolution** — each contract's `intent_ref.id`/`schema` compared to the actual intent object, every
  repo-artifact locator resolved with `git cat-file -e <commit>:<path>`, every `sha256:h@path`
  content-hash verified, the grounding source hash == `9d1ab3a5…`; (c) **derive-from-authored-data** —
  nodes/outputs/edges/roots/critical-path derived and cross-checked (authored == derived); (d) DAG;
  (e) parallel write surfaces; (f) gate invariants; (g) the **authored completion relation evaluated
  as structured data** — the predicate dependency graph is built by **walking the `expr` ASTs** and
  proved acyclic (a tautology → self-cycle → FAIL), and **each fixture is computed and compared to its
  authored `expected`** (a flipped expectation → FAIL); (h) the **oracle-ownership TOTAL + SINGULAR
  bijection** (registry ⇄ each child's `acceptance` ⇄ the md projection); (i) **ledger consistency**
  (revision labels agree · reported counts == derived table counts == registry size · every category a
  single enum member). **Exits non-zero on any violation; passes at this tree (all nine green).**
- [`validate_test.py`](./validate_test.py) — executable adversarial-mutation harness. Materializes the
  **eleven** mutations (six originals: wrong intent id; nonsense `cell.class`; nonexistent artifact
  path; tautological whole-wave predicate; flipped fixture expectation; placeholder oracle operand —
  plus the five R5 bijection mutations: remove a registry entry; duplicate an owner; reclassify an
  entry; leave an mv predicate unowned; add an owner entry absent from acceptance) in temp trees —
  honestly re-pinning every changed `contract_sha256` and the registry hash — and asserts **each exits
  non-zero for its own named
  predicate** while the clean tree exits 0. **Verified: harness passes** (clean 0; the six originals
  fail on (b)/(a)/(b)/(g)/(g)/(h) and the five bijection mutations on (h)).
- [`oracle-registry.yaml`](./oracle-registry.yaml) — the authoritative registry (30 mechanically-verifiable
  entries); [`acceptance-oracles.md`](./acceptance-oracles.md) is its projection, every predicate classified
  as **exactly one** of **enforced** / **mechanically-verifiable** / **evidenced** / **cognitive-review**
  (single-kind; check (i) fails closed on a compound category). No semantic-absence claim is implemented
  as grep-and-called-mechanical.

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

Note: R2 finding-2's phrasing above ("materialized … exists before any cell") is **superseded by R3
finding 2** — the intent is honestly a transitional bootstrap projection authored during the cell; it
does not exist before the cell.

## Per-finding disposition (R5 — repair of the R4 oracle-ownership method)

| # | Finding | Repair in this R5 |
|---|---|---|
| 1 | **[BLOCKER]** R4 broke the exact §2 shape by adding a top-level `oracles` key; oracles were not content-bound | Removed the `oracles` key from all six contracts (each normalizes to the EXACT §2 top-level paths again); reverted `validate.py` `TOP_KEYS` to the §2 set so check (a) FAILS on any extra top-level key (incl. `oracles`); created the separate authoritative [`oracle-registry.yaml`](./oracle-registry.yaml) (30 entries, no placeholders); each contract consumes its slice via a canonical `inputs.required[].external` control_plane ref (content-hash pinned); `acceptance-oracles.md` is now a parity-checked **projection** of the registry. |
| 2 | **[BLOCKER]** check (h) was not a true bijection (β removed an entry and it still passed) | Rewrote check (h) as a **total + singular bijection** across registry ⇄ each child's `acceptance.predicates` ⇄ the md projection; rejects a missing entry, duplicate owner, extra owner absent from acceptance, classification mismatch, placeholder, and any registry↔projection parity break. Extended `validate_test.py` with the five bijection mutations (incl. the exact R4-breaking omission) plus the six originals — **all 11 fail for their own predicate; clean passes**. |
| 3 | **[REQUIRED]** revision + classification ledger was stale/inconsistent (said R3/R4; mv miscounted 28) | Advanced every revision label to **R5** (machine-checkable `wave-revision:` markers that check (i) proves agree); corrected classification counts derived from the tables/registry (mechanically-verifiable **30**); **atomized** the two compound `enforced + evidenced` WC-5 rows to single-kind `enforced`; added check (i) (labels agree · reported counts == derived == registry size · every category a single enum member). |

## Per-finding disposition (second external-β ITERATE — R3)

| # | Finding | Repair in this R3 |
|---|---|---|
| 1 | **[BLOCKER]** `validate.py` false-passed five adversarial mutations | Rewrote the validator to be genuinely **sound**: full §2 constraint model (enums/types/cardinalities), real ref resolution (intent id/schema vs the intent object; every repo-artifact locator resolved with `git cat-file -e`; grounding source hash), derive-from-authored-data (nodes/edges/roots/critical-path), and **evaluation** of the authored completion predicates + fixtures as structured data (AST-walked predicate-graph acyclicity + per-fixture computation vs `expected`). Added [`validate_test.py`](./validate_test.py); **all five now exit non-zero for their own predicate; clean tree exits 0** (verified). |
| 2 | **[BLOCKER]** intent masqueraded as pre-cell κ/operator intent | Rewrote `intent.cn-intent-v1.yaml` as an explicitly **transitional bootstrap projection** authored during the cell (no pre-cell existence claim); `statement` carries only operator matter (#627 + the operator verbatim doctrine line); α conclusions moved to [`decision-provenance.md`](./decision-provenance.md); identity vs carrier kept distinct; README no longer says the intent existed before the cell. |
| 3 | **[REQUIRED]** oracles mislabeled cognitive review as mechanical | Rewrote `acceptance-oracles.md`: every predicate classified **enforced / mechanically-verifiable / evidenced / cognitive-review**; mechanically-verifiable rows name the fixture path + command + expected positive/negative the child WC must emit (added `.cdd/unreleased/<wc>/fixtures/**` allowed-paths + a receipt-evidence acceptance predicate to each contract); semantic-absence claims are cognitive-review, not grep. |

## Files

- [`README.md`](./README.md) — this overview.
- [`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) — transitional bootstrap intent projection (operator matter only; honest provenance).
- [`decision-provenance.md`](./decision-provenance.md) — α/β planning conclusions moved out of intent (WC graph, D9 realization, β-repair outcomes), each citing its actual source.
- [`wave.cn-wave-v1.yaml`](./wave.cn-wave-v1.yaml) — the `cn.wave.v1` instance (nodes, derived edges, gates, STOP conditions, `contract_ref` resolution, **structured** non-recursive completion).
- [`contracts/wc-1..wc-5.cn-cell-contract-v1.yaml`](./contracts/) — the six §2-conformant child contracts.
- [`grounding-source-5015460988.md`](./grounding-source-5015460988.md) — byte-exact source snapshot (SHA `9d1ab3a5…`).
- [`grounding-cm.md`](./grounding-cm.md) — α's derivative FAIL-classification (references the snapshot).
- [`reconcile-627.md`](./reconcile-627.md) — the #627 S1–S8 refinement map + derived reverse-consumers.
- [`oracle-registry.yaml`](./oracle-registry.yaml) — AUTHORITATIVE machine-readable oracle registry (30 mechanically-verifiable entries), consumed by each contract via a content-bound `inputs.required[].external` ref.
- [`acceptance-oracles.md`](./acceptance-oracles.md) — projection of the registry; every predicate classified enforced / mechanically-verifiable / evidenced / cognitive-review (single-kind).
- [`validate.py`](./validate.py) — SOUND pre-authorization validator (checks a–i; exits non-zero on violation).
- [`validate_test.py`](./validate_test.py) — executable adversarial-mutation harness (11 mutations each fail for their own predicate; clean tree passes).

---
*Status: R5 wave, α matter — under operator review; external-β and CC review next. No child WCs
dispatched; no control-plane action taken by this cell.*
