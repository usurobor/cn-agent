<!-- wave-revision: R8 -->
# Plan-local structural schemas (cnos#671 R8) — `cue vet`, no Python

This directory is the R8 re-architecture of pre-authorization validation for the cell-runtime wave.

**Operator directive (R8):** this is a **Go + CUE** repository. **No Python.** Schemas are **CUE**;
procedural code is **Go**. The prior rounds shipped a hand-rolled Python validator
(`validate.py` / `validate_test.py`) — a tool-choice error. It is **removed entirely**. Structural
pre-authorization validation is now **`cue vet`** against the closed CUE definitions here; the
procedural/semantic checks CUE cannot express are **named-and-deferred to Go**.

These schemas are **PLAN-LOCAL and TRANSITIONAL**: they are the design inputs that WC-1/WC-2 (and
#627 S2–S3) formalize into the **canonical** `schemas/cdd/*.cue`. They are **not** the canonical
deliverable and they do **not** edit the shipped `schemas/cdd/`.

## Files

- [`cell_contract.cue`](./cell_contract.cue) — `#CellContract`: the exact `cn.cell.contract.v1` §2
  envelope (closed struct; `cell.class` / `ref_kind` / locator `kind` / `requested_output.kind`
  enums; the gate-invariant conditional).
- [`wave.cue`](./wave.cue) — `#Wave`: the full `cn.wave.v1` envelope (closed top-level; typed nodes/
  edges with uniqueness folds; nonempty typed `external_roots` and STOP set; required revision-bound
  `gates.wave_authorization`; required `contract_ref_resolution`; pinned `edge_derivation` constants;
  the pinned completion `wave_complete` const + closed 5-bool child records).
- [`assurance_registry.cue`](./assurance_registry.cue) — `#AssuranceRegistry`: registry shape,
  classification enum, per-category required fields.
- [`intent.cue`](./intent.cue) — `#Intent`: the transitional bootstrap intent projection shape.
- [`Makefile`](./Makefile) — `make clean` (vet real artifacts, PASS) · `make regressions` (each bad
  fixture rejected) · `make all`.
- [`regressions/`](./regressions/) — a clean base + one-mutation bad fixtures (see its README).

## Validation architecture

### STRUCTURAL — enforced NOW by `cue vet` (this tree)

One invocation per artifact/schema pair (run via `make -C schema clean`):

```
cue vet ./schema/ ./wave.cn-wave-v1.yaml                        -d '#Wave'
cue vet ./schema/ ./contracts/wc-1.cn-cell-contract-v1.yaml     -d '#CellContract'
cue vet ./schema/ ./contracts/wc-2.cn-cell-contract-v1.yaml     -d '#CellContract'
cue vet ./schema/ ./contracts/wc-3a.cn-cell-contract-v1.yaml    -d '#CellContract'
cue vet ./schema/ ./contracts/wc-3b.cn-cell-contract-v1.yaml    -d '#CellContract'
cue vet ./schema/ ./contracts/wc-4.cn-cell-contract-v1.yaml     -d '#CellContract'
cue vet ./schema/ ./contracts/wc-5.cn-cell-contract-v1.yaml     -d '#CellContract'
cue vet ./schema/ ./oracle-registry.yaml                        -d '#AssuranceRegistry'
cue vet ./schema/ ./intent.cn-intent-v1.yaml                    -d '#Intent'
```

CUE closes the R7 wave-envelope + completion blockers **natively** — no procedural code:

| R7 blocker (the Python validator false-passed) | Closed natively by |
|---|---|
| `bool("false")` coercion (`receipt_bound: "false"` accepted) | typed `bool` field → `"false"` (string) conflicts |
| non-mandatory `evidence_bound` | `#ChildRecord` closed 5-bool struct, every constituent required |
| duplicate node id / duplicate edge invisible | `_unique_node_ids` / `_unique_edges` folds → conflicting values |
| duplicate mapping key silently collapsed | CUE errors on conflicting values for the repeated key |
| empty `external_roots` | `[string, ...string]` nonempty |
| `contract_ref_resolution: {}` | required (`!`) fields → "field is required but not present" |
| `gates.wave_authorization` missing revision-bound fields | `revision_bound!` (and authority/boundary) required |
| deletable/literal-mutated completion formula | `wave_complete!: "all(N child_complete) AND child_complete(terminal)"` — required + pinned const |
| deleted/flipped `edge_derivation.parity_required`/`acyclic` | pinned `true` consts → delete = required-fail, flip = conflict |
| un-typed authority substructure / unknown top-level key | closed structs → "field not allowed" |
| bad `cell.class` / `requested_output.kind` / STOP `effect` | enum disjunctions → "empty disjunction" |
| gate invariant (doctrine_affecting ⇒ acceptance) | CUE conditional `if doctrine_affecting { gates: operator_acceptance_required: true }` |

### PROCEDURAL / SEMANTIC — named-and-deferred to Go (option B, D9)

The checks CUE cannot express **declaratively** are **named as deliverables of Go tooling**
(WC-1/WC-2 and #627 S2–S3), **not** implemented here (and never in Python). They run when each WC
executes under its own α→β→γ→CC:

| Deferred-Go check | Owner |
|---|---|
| graph **acyclicity** (DAG over the derived dependency edges) | WC-3b/WC-5 |
| sibling-output → wave-edge **parity** (authored == derived, exactly) | WC-3b/WC-5 |
| git **ref / content-hash resolution** (intent id, `commit:path`, revision hashes, grounding source) | WC-1 / #627 S2 |
| classification **bijection** (union(child acceptance.predicates) ⇄ registry, each once) | WC-1 / #627 S2–S3 |
| **completion-evidence derivation** over resolved records (child_complete AND, predicate-graph acyclicity, fixture `expected`) | WC-3b/WC-5 |
| ledger consistency (revision markers agree; per-category counts) · parallel write-surface disjointness | WC-5 |

CUE fixes the **shape** now; Go **computes** the procedural facts when the wave executes.

## Reproduce

```
go install cuelang.org/go/cmd/cue@latest   # cue v0.17+
make -C .cdd/waves/cell-runtime-doctrine/schema all
```
