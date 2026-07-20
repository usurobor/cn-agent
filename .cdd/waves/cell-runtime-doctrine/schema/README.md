<!-- wave-revision: R9 -->
# Plan-local structural schemas (cnos#671 R9) — `cue vet`, no Python

This directory is the pre-authorization validation for the cell-runtime wave.

**Operator directive:** this is a **Go + CUE** repository. **No Python.** Schemas are **CUE**; procedural
code is **Go**. A prior round shipped a hand-rolled Python validator (`validate.py`) — a tool-choice
error, **removed entirely** (R8). Structural pre-authorization validation is **`cue vet`** against the
closed CUE definitions here; the procedural/semantic checks CUE cannot express are **named-and-deferred
to single-owner Go validators**.

**R9 changes.** (1) `#CellContract` is now the **faithful** canonical `cn.cell.contract.v1` §2 shape —
nullable scope refs, 1+ cardinalities, the `prior_receipt` locator class, all three output kinds, and
the gate **truth table** (see below). (2) A named **`#WorkingCellContract`** refinement (class=working,
output kind=artifact) validates the six real Working-Cell contracts, while `#CellContract` stays the
canonical shape (exercised by the canonical-variant fixtures). (3) Every **deferred-go** procedural check
is owned by **exactly one** in-wave WC. (4) The completion constituents are made consistent
(`evidence_bound` now listed in `required_constituents`) and a **typed resolver-input contract**
(`#CompletionEvidenceInput`) types the completion-validator's input.

These schemas are **PLAN-LOCAL and TRANSITIONAL**: design inputs that the Go/CUE deliverables (and
#627 S2–S3, as downstream consumers/canonicalizers) formalize into the **canonical** `schemas/cdd/*.cue`.
They are **not** the canonical deliverable and they do **not** edit the shipped `schemas/cdd/`.

## Files

- [`cell_contract.cue`](./cell_contract.cue) — `#CellContract`: the **exact** `cn.cell.contract.v1` §2
  envelope (closed struct; nullable `scope.wave`/`scope.parent_cell`; 1+ `inputs.required`; external-only
  `inputs.optional`; `repo_artifact | control_plane | prior_receipt` locator union; `artifact |
  relation_graph | judgment` output kinds; required `forbidden_paths`/`non_goals`; the gate truth-table).
  Plus **`#WorkingCellContract`** — the named WC refinement (class=working, kind=artifact).
- [`wave.cue`](./wave.cue) — `#Wave`: the full `cn.wave.v1` envelope; plus `#ChildCompleteDef` (5 pinned
  constituents incl. `evidence_bound`), `#ChildRecord`, `#CompletionEvidenceInput` (the typed resolver
  input), and `#EvidenceResolver` (owned by WC-5).
- [`assurance_registry.cue`](./assurance_registry.cue) — `#AssuranceRegistry`: registry shape,
  classification enum, per-category required fields; deferred-go entries require a single `deferred_owner`.
- [`intent.cue`](./intent.cue) — `#Intent`: the transitional bootstrap intent projection shape.
- [`Makefile`](./Makefile) — `make clean` · `make regressions` · `make all`.
- [`regressions/`](./regressions/) — clean/canonical bases + one-mutation bad fixtures (see its README).

## Validation architecture

### STRUCTURAL — enforced NOW by `cue vet` (this tree)

```
cue vet ./schema/ ./wave.cn-wave-v1.yaml                        -d '#Wave'
cue vet ./schema/ ./contracts/wc-{1,2,3a,3b,4,5}.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
cue vet ./schema/ ./oracle-registry.yaml                        -d '#AssuranceRegistry'
cue vet ./schema/ ./intent.cn-intent-v1.yaml                    -d '#Intent'
# canonical §2 variants against the FAITHFUL #CellContract (not the WC refinement):
cue vet ./schema/ ./regressions/contract.canonical-nullable-scope.yaml        -d '#CellContract'
cue vet ./schema/ ./regressions/contract.canonical-optional-prior-receipt.yaml -d '#CellContract'
cue vet ./schema/ ./regressions/contract.canonical-relation-graph.yaml        -d '#CellContract'
cue vet ./schema/ ./regressions/contract.canonical-judgment.yaml              -d '#CellContract'
```

The faithful `#CellContract` rejects every §2 drift the external β found, natively:

| §2-drift mutation | Rejected by |
|---|---|
| `scope.wave`/`parent_cell` forced to `string` (no null) | `string \| null` — a null is now valid, and the canonical variants prove it |
| empty `inputs.required` | `[#InputRef, ...#InputRef]` (1+) |
| a `sibling_output` in `inputs.optional` | `optional!: [...#ExternalRef]` — external-only |
| missing the `prior_receipt` locator class | `#Locator = repo_artifact \| control_plane \| prior_receipt` |
| `requested_output.kind` limited to `artifact` | `artifact \| relation_graph \| judgment` |
| empty `acceptance.predicates` / empty `allowed_paths` | `[string, ...string]` (1+) |
| missing `forbidden_paths` / `non_goals` key | required (`!`) keys, each 0+ |
| both gates false but `reason` non-null (or a gate true but `reason` null/absent) | truth-table conditionals: nonempty string iff a gate bool true, null iff both false |
| `doctrine_affecting: true` with `operator_acceptance_required: false` | `if doctrine_affecting { gates: operator_acceptance_required: true }` |
| unknown top-level key | closed `#CellContract` struct |
| `required_constituents` drops `evidence_bound` / definition ignores a constituent | `#ChildCompleteDef` pins all 5 |

### PROCEDURAL / SEMANTIC — named-and-deferred to SINGLE-OWNER Go validators (D9)

Each deferred-go check has **exactly one** in-wave owner (`deferred_owner` in `oracle-registry.yaml`;
never #627 — #627 S2–S3 stay downstream consumers/canonicalizers). WC-5 depends (sibling_output) on
wc-1/wc-2/wc-3b and each owner's gating acceptance predicate is a constituent of its `child_complete`,
so **WC-5 cannot seal until every deferred validator exists and passes**.

| Deferred-Go check | Single owner | Go artifact |
|---|---|---|
| git **ref / content-hash resolution** | **WC-2** | `.cdd/unreleased/wc-2/validators/ref_resolve.go` |
| graph **acyclicity** (DAG) | **WC-3b** | `.cdd/unreleased/wc-3b/validators/wave_dag.go` |
| sibling-output → wave-edge **parity** | **WC-3b** | `.cdd/unreleased/wc-3b/validators/edge_parity.go` |
| parallel **write-surface disjointness** | **WC-3b** | `.cdd/unreleased/wc-3b/validators/write_surface.go` |
| **oracle-ownership bijection** (each mech-verifiable predicate → one checker) | **WC-1** | `.cdd/unreleased/wc-1/validators/oracle_ownership_bijection.go` |
| **completion-evidence derivation** (typed resolver input → 5 derived constituents) | **WC-5** | `.cdd/unreleased/wc-5/validators/completion_evidence.go` |
| **ledger consistency** (revision markers agree; per-category counts) | **WC-5** | `.cdd/unreleased/wc-5/validators/ledger_consistency.go` |
| **classification-totality bijection** (union(predicates) ⇄ registry, each once, all 5 categories) | **WC-5** | `.cdd/unreleased/wc-5/validators/classification_bijection.go` |

Each owner contract pins the Go artifact path in `allowed_paths` and a gating acceptance predicate whose
PASS gates its completion; the full spec (typed inputs, result/evidence shape, positive + named negative
fixtures, downstream consumers, gating predicate) is pinned in `oracle-registry.yaml` `wave_predicates`.
CUE fixes the **shape** now; Go **computes** the procedural facts when the owning WC executes.

## Reproduce

```
go install cuelang.org/go/cmd/cue@latest   # cue v0.17+
make -C .cdd/waves/cell-runtime-doctrine/schema all
```
