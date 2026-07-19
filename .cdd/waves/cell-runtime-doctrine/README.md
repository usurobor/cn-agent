<!-- wave-revision: R8 -->
# Wave: cell-runtime-doctrine (cnos#671 — R8)

**Planning Cell output.** This directory is the matter of the Planning Cell #671 (child of parent
wave #627): a mature, executable `cn.wave.v1` plan that decomposes the cell-runtime doctrine into
single-purpose Working-Cell contracts, grounded in an immutable coherence measurement.

**R8 re-architecture (operator directive).** This is a **Go + CUE** repository. **No Python.**
Schemas are **CUE**; procedural code is **Go**. The prior rounds shipped a hand-rolled Python
validator (`validate.py` / `validate_test.py`) — a tool-choice error. R8 **removes it entirely** and
re-architects pre-authorization validation the right way:

- **Structural** pre-authorization validation = **`cue vet`** against plan-local **closed CUE**
  schemas under [`schema/`](./schema/) (`#CellContract`, `#Wave`, `#AssuranceRegistry`, `#Intent`).
  CUE closes the R7 wave-envelope + completion blockers **natively** — closed structs reject unknown
  fields; typed fields reject `"false"` for a bool; enums/cardinality/required-field constraints;
  duplicate mapping keys and duplicate node/edge ids conflict; the completion formula is a pinned
  required const; `evidence_bound` is a mandatory typed bool. No procedural validator is involved.
- **Procedural / semantic** validation = **deferred Go validators** the Working Cells (WC-1/WC-2) and
  **#627 S2–S3** produce (option B, D9: *named*, not implemented here): graph acyclicity (DAG),
  sibling-output → edge **parity**, git ref/content-hash resolution, the classification **bijection**,
  and **completion-evidence derivation** over resolved records. CUE fixes the shape now; Go computes
  the procedural facts when each WC executes under its own α→β→γ→CC.

This is a **plan, not doctrine**. The `docs/` artifacts named as `requested_output` paths are the
**future output of the Working Cells** — not authored here. Nothing under `docs/` or `schemas/cdd/`
is edited by this cell. The `schema/` CUE is **plan-local and transitional**: the design input WC-1/
WC-2 formalize into the canonical `schemas/cdd/*.cue`, **not** the canonical deliverable.

The accepted decomposition (WC-2 → WC-1 → {WC-3a, WC-3b, WC-4} → WC-5), the six-node graph, D9-four,
intent provenance, grounding, and the §2 contracts are **unchanged** across R2–R8. R8 is a
validation-tooling re-architecture + honest scoping.

## Goal

Produce a coherent, mechanically-checkable specification of the cnos cell-runtime doctrine — the
typed cell contract, the coherence-measurement (CM) backbone, the cell/wave FSMs, and the
shipped→specified migration — decomposed so each artifact is single-purpose and individually
coherent. **Done means:** a Working Cell can implement the runtime (Go/CUE) from these artifacts
**without inventing** dependency, measurement, authority, or gate policy.

## Validation architecture (R8)

### Structural — enforced NOW by `cue vet`

Run `make -C schema clean` (or the per-artifact invocations in [`schema/README.md`](./schema/README.md)):

```
cue vet ./schema/ ./wave.cn-wave-v1.yaml                      -d '#Wave'
cue vet ./schema/ ./contracts/wc-{1,2,3a,3b,4,5}.cn-cell-contract-v1.yaml -d '#CellContract'
cue vet ./schema/ ./oracle-registry.yaml                      -d '#AssuranceRegistry'
cue vet ./schema/ ./intent.cn-intent-v1.yaml                  -d '#Intent'
```

The R7 blockers close natively (full table in [`schema/README.md`](./schema/README.md)):
`receipt_bound: "false"` (string) → type conflict; a record missing `evidence_bound` → required-field
fail; a duplicate node id / duplicate edge / duplicate mapping key → conflicting values;
`external_roots: []` / `stop_conditions: []` → nonempty-list fail; `contract_ref_resolution: {}` and a
missing `gates.wave_authorization.revision_bound` → required-field fail; a deleted or flipped
`edge_derivation.parity_required`/`acyclic` → required/const fail; a deleted or literal-mutated
completion `wave_complete` formula → required/const fail; an added top-level key (e.g. the R4-regression
`oracles`) → "field not allowed"; a bad `cell.class`/`requested_output.kind`/STOP `effect` → empty
disjunction; a doctrine-affecting cell that does not require operator acceptance → gate-invariant
conflict.

The [`schema/regressions/`](./schema/regressions/) suite proves each rejection: a clean base PASSES and
**20 one-mutation bad fixtures are each rejected** by `cue vet` (`make -C schema regressions`; verified
green with cue v0.17.1). This is **not** a blanket soundness certificate — `cue vet` warrants exactly
the closed-struct / typed-field / enum / cardinality / required-field / pinned-const constraints the
schemas declare, and **nothing beyond them**. Everything else is named-and-deferred to Go.

### Procedural / semantic — named-and-deferred to Go

| Deferred-Go check (not run at this tree; the WC/#627 Go validator runs it when the WC executes) | Owner |
|---|---|
| graph **acyclicity** (DAG over derived dependency edges) | WC-3b/WC-5 |
| sibling-output → wave-edge **parity** (authored == derived, exactly) | WC-3b/WC-5 |
| git **ref / content-hash resolution** (intent id/schema; `commit:path`; revision hashes; grounding source `9d1ab3a5…`) | WC-1 / #627 S2 |
| classification **bijection** (`union(child acceptance.predicates)` ⇄ registry, each classified once) | WC-1 / #627 S2–S3 |
| **completion-evidence derivation** (child_complete AND over records; predicate-graph acyclicity; fixture `expected`) | WC-3b/WC-5 |
| ledger consistency (revision markers agree; per-category counts) · parallel write-surface disjointness | WC-5 |

## Grounding CM (content-bound to the TRUE source)

The whole wave is **grounded in** an immutable measurement, not "producing CM as a telos."

- **Immutable source snapshot (exact bytes):** [`grounding-source-5015460988.md`](./grounding-source-5015460988.md)
  — byte-for-byte capture of PR #667 comment `5015460988`.
- **True source SHA-256:** `9d1ab3a50d00e642fdeb87728dd71f7c7499c60878afe7001f2ddb832b161dbb` (150-line body).
- **α's derivative:** [`grounding-cm.md`](./grounding-cm.md) — clearly marked as α's **derived**
  FAIL-classification that *references* the snapshot by path + source hash (not the exact bytes).
- `wave.cn-wave-v1.yaml` grounding points to **both** (source hash + derivative content hash). Every
  node's grounding `external` ref content-binds the **source hash**. (Ref/hash resolution is a
  **deferred-Go** check.)

**CM ownership.** CM is a **runner/provider-produced measurement object** (producer = a CM provider/
runtime op, e.g. tsc — **not** CC). A **CC's output is a JUDGMENT that consumes immutable CM refs**;
CM and CC-judgment are **distinct tagged objects**. **V gates on CM.** A **PC's wave is grounded in /
references an immutable CM.** WC-2 encodes exactly this.

## Intent (TRANSITIONAL bootstrap projection — honest provenance)

[`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) is an **explicitly transitional bootstrap
projection** of operator intent, authored **during** this cell — SHA-256
`1e246846774c43365cb00dc1b39b6998bab535f4833a16fcfe32ca3e2825560e`. It does **not** claim to have
existed before any cell: D9 **names** `cn.intent.v1` but does not implement it. Its `statement` carries
**only authoritative operator matter** (the objective as the operator directed it in #627, incl. the
operator's **verbatim final doctrine line**); its carriers point to the real operator-intent sources
(#627, the operator verbatim line). Every **α-derived conclusion** lives in
[`decision-provenance.md`](./decision-provenance.md). Identity (`intent_ref.id`) vs carrier
(`carrier.kind: repo_artifact_bootstrap`, a projection pointer) are kept distinct. `#Intent` (cue vet)
validates the projection's shape; identity/carrier **resolution** is a deferred-Go check.

## Corrected dependency graph (WC-2 keystone root) — unchanged shape

```
external roots: grounding-CM(@sha 9d1ab3a5) + #628/S1 + shipped schemas/CCNF/transitions.json
  → WC-2  (CM measurement object + receipt_core→CM→V→δ→final_receipt type path)   [keystone root]
  → WC-1  (cell classes + typed cell contract; imports cm_ref from WC-2)
      → WC-3a (cell FSM) · WC-3b (wave FSM) · WC-4 (shipped→specified migration)   [each dep WC-1 + WC-2]
          → WC-5  (integration / seal: whole-wave proof)   [dep WC-1,2,3a,3b,4]
```

Edges are the mechanical projection of the child contracts' `sibling_output` refs (see
`wave.cn-wave-v1.yaml` `edge_derivation`); **parity (authored == derived) is a deferred-Go check**.
Critical path: **WC-2 → WC-1 → WC-3b → WC-5**.

## §2 conformance

All six `contracts/*.yaml` are canonical `cn.cell.contract.v1` §2 instances with the **exact**
top-level key set and **no extras** — enforced now by `cue vet … -d '#CellContract'` (a closed struct;
an added top-level key such as `oracles`/`consumers`/`completion_signal` is rejected as "field not
allowed"). Oracle ownership lives in the separate content-bound [`oracle-registry.yaml`](./oracle-registry.yaml),
which each contract consumes via a canonical `inputs.required[].external` control_plane ref (immutable
locator + content hash). Reverse-consumer information is **derived** (see [`reconcile-627.md`](./reconcile-627.md)).

## Non-recursive, structured whole-wave completion

`wave.cn-wave-v1.yaml` `completion` defines a **non-recursive construction set**
`N = {wc-1, wc-2, wc-3a, wc-3b, wc-4}` (excludes wc-5), a terminal `wc-5`, and a **pinned governing
formula const** `wave_complete: "all(N child_complete) AND child_complete(terminal)"` (required + fixed
by `#Wave` — deleting or mutating it fails `cue vet`). Each fixture's per-node record is a **closed
5-bool struct** (`requested_output_produced`, `acceptance_all_pass`, `class_v_pass`, `receipt_bound`,
`evidence_bound` — every constituent mandatory and typed bool). The **derivation** — that
child_complete(n) is the AND over node n's bound record, that the predicate graph is acyclic, and that
each fixture computes its authored `expected` — is a **deferred-Go** check (WC-3b/WC-5). Four fixtures
pin the intended behavior (all-incomplete; all-N-complete/WC-5-absent; WC-5 V-FAIL; all-complete).

## D9 — four-schema boundary SETTLED

Per operator decision, D9's four-schema boundary stays **settled**. CM is realized **within** the four
schemas (a typed CM field/edge in the receipt + a `cm_ref` resolving within the existing schemas),
**not** a fifth canonical `cn.cm.v1` schema. WC-1 and WC-4 consume the same in-four-schema `cm_ref`
shape; WC-2 carries a reject-a-fifth-canonical-schema negative fixture.

## Wave-level STOP conditions + immutable `contract_ref` resolution

`wave.cn-wave-v1.yaml` carries typed wave-level `stop_conditions` (stale grounding, failed predecessor,
invalid graph, revoked authorization, integration conflict, unresolved-D9-at-wave-scope) — each a
single typed non-dispatch transition (hold/replan/escalate) + receipt (nonempty typed set enforced by
`#Wave`). `contract_ref_resolution` binds every `contract_ref` to the **immutable authorized wave
commit**; operator authorization is **revision-bound**; re-resolution at a changed tree **fails stale**.
Each node carries `contract_sha256` (the re-pinned SHA-256 of its contract blob).

## Assurance registry — honest scoping

[`oracle-registry.yaml`](./oracle-registry.yaml) is the authoritative **TOTAL** registry: one
`assurance:` entry for **every** child acceptance predicate, classified into **exactly one** of five
honest categories, plus a separate `wave_predicates:` section. [`acceptance-oracles.md`](./acceptance-oracles.md)
is its projection. The categories (R8 taxonomy):

- **structural-cue** — enforced NOW by `cue vet` against a plan-local `#Def`.
- **deferred-go** — a named Go validator WC-1/WC-2/#627 S2–S3 ship (procedural; not run at this tree).
- **mechanically-verifiable** — a per-child acceptance oracle the child WC **emits** as a **Go checker
  or a CUE schema** (**no `.py`**) + named fixtures + a credential-free command, bound in its receipt.
- **evidenced** — bound receipt evidence, not self-report.
- **cognitive-review** — an independent β / CC doctrine judgment, honestly not mechanical.

**Child acceptance predicates (69):** structural-cue **1** · deferred-go **10** · mechanically-verifiable
**30** · evidenced **7** · cognitive-review **21** — one classification each. The former "enforced"
category is split honestly: the gate invariant is `structural-cue` (a CUE conditional); edge-parity,
completion-evidence, oracle-ownership bijection, and readiness are `deferred-go`. Wave-only structural/
deferred predicates are held **separately** in `wave_predicates:` so they can neither inflate nor mask
child coverage.

## Node list

| Node | Class / domain | Output (`docs/…`) | Depends on (sibling_output edges) |
|---|---|---|---|
| **WC-2** *(keystone root)* | working / doctrine | `COHERENCE-MEASUREMENT.md` | — (external roots only) |
| **WC-1** | working / doctrine | `CELL-RUNTIME-CLASSES.md` | WC-2 |
| **WC-3a** | working / doctrine | `CELL-FSM.md` | WC-1, WC-2 |
| **WC-3b** | working / doctrine | `WAVE-FSM.md` | WC-1, WC-2 |
| **WC-4** | working / doctrine | `CONTRACT-MIGRATION.md` | WC-1, WC-2 |
| **WC-5** *(seal)* | working / doctrine | `CELL-RUNTIME-INTEGRATION.md` | WC-1, WC-2, WC-3a, WC-3b, WC-4 |

## Reclassification (#662 → WC-1)

`#662` ("PC-D0") is **retired as a PC**. Its converged artifact (`CELL-RUNTIME-CLASSES.md`) re-lands as
**WC-1**, re-scoped and classified as WC-doctrine. **No work is discarded.** WC-1 imports the `cm_ref`
interface from WC-2. The κ≠α role logic and content-bound-review rule repaired in #662 R8 are carried
into WC-1 as **settled input**.

## Per-finding disposition (R8 — remove Python; CUE + deferred Go; honest scoping)

| # | Directive | R8 disposition |
|---|---|---|
| 1 | **Go + CUE repo; NO Python** | Deleted `validate.py` and `validate_test.py` (and every `*.py` under the wave dir → none remain). |
| 2 | Structural pre-authorization validation the right way | Authored plan-local closed CUE (`schema/#CellContract`, `#Wave`, `#AssuranceRegistry`, `#Intent`); wired the exact `cue vet` invocations (schema/Makefile + README). All clean artifacts PASS. |
| 3 | Close the R7 blockers **via CUE, natively** | Typed bool (rejects `"false"`), closed structs (reject unknown keys + un-typed authority substructure), duplicate-key/id conflicts, nonempty `external_roots`, required `contract_ref_resolution`/`wave_authorization` fields, pinned `edge_derivation` consts, a pinned required completion formula const, a mandatory `evidence_bound`. Proven by 20 rejected regression fixtures. |
| 4 | Name-and-defer the procedural checks to Go | DAG acyclicity, edge parity, ref/hash resolution, classification bijection, completion-evidence derivation → **named** deliverables of WC-1/WC-2 and #627 S2–S3 (D9, option B). Not reimplemented; never in Python. |
| 5 | Rewire the registry + honest scoping | Every `enforced_by: validate.py check (x)` now points at the CUE constraint (`cue vet` #Wave/#CellContract/…) or the named deferred Go validator; `mechanically-verifiable` checkers are Go/CUE (no `.py`); removed the blanket "sound/fail-closed" claim. |
| 6 | Provenance + ledger | Operator directive recorded in `decision-provenance.md`; every `wave-revision:` marker advanced to **R8**; content hashes re-pinned for edited files; contracts stay exactly §2. |

## Prior-round ledger (unchanged accepted structure)

R2 repaired the first external-β ITERATE (#672); R3 a second (honest intent provenance, honestly
classified oracles, byte-exact grounding); R4/R5 moved oracle ownership into the content-bound
`oracle-registry.yaml`; R6 made that registry TOTAL; R7 tried to close the wave-envelope + completion
blockers **inside the Python validator** — which R8 supersedes by removing Python and closing them in
CUE. The accepted decomposition graph (WC-2 → WC-1 → {WC-3a, WC-3b, WC-4} → WC-5), the §2 contracts,
D9-four, intent provenance, and grounding are **unchanged** across R2–R8.

## Files

- [`README.md`](./README.md) — this overview.
- [`schema/`](./schema/) — plan-local transitional CUE (`#CellContract`, `#Wave`, `#AssuranceRegistry`,
  `#Intent`), the `Makefile` runner, and `regressions/` (clean bases + 20 rejected mutation fixtures).
- [`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) — transitional bootstrap intent projection.
- [`decision-provenance.md`](./decision-provenance.md) — α/β planning conclusions + the R8 operator directive.
- [`wave.cn-wave-v1.yaml`](./wave.cn-wave-v1.yaml) — the `cn.wave.v1` instance.
- [`contracts/wc-1..wc-5.cn-cell-contract-v1.yaml`](./contracts/) — the six §2-conformant child contracts.
- [`grounding-source-5015460988.md`](./grounding-source-5015460988.md) — byte-exact source snapshot (SHA `9d1ab3a5…`).
- [`grounding-cm.md`](./grounding-cm.md) — α's derivative FAIL-classification (references the snapshot).
- [`reconcile-627.md`](./reconcile-627.md) — the #627 S1–S8 refinement map + derived reverse-consumers.
- [`oracle-registry.yaml`](./oracle-registry.yaml) — AUTHORITATIVE, TOTAL assurance registry (five honest categories).
- [`acceptance-oracles.md`](./acceptance-oracles.md) — projection of the total registry (every predicate classified single-kind).

---
*Status: R8 wave, α matter — under operator review; external-β and CC review next. No child WCs
dispatched; no control-plane action taken by this cell.*
