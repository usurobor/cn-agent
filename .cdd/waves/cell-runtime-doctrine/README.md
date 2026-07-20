<!-- wave-revision: R9 -->
# Wave: cell-runtime-doctrine (cnos#671 — R9)

**Planning Cell output.** This directory is the matter of the Planning Cell #671 (child of parent
wave #627): a mature, executable `cn.wave.v1` plan that decomposes the cell-runtime doctrine into
single-purpose Working-Cell contracts, grounded in an immutable coherence measurement.

**R9 repair (external-β #672).** `#CellContract` is now the **faithful** canonical `cn.cell.contract.v1`
§2 shape (nullable scope, 1+ cardinalities, `prior_receipt` locator, all output kinds, gate truth-table);
a named **`#WorkingCellContract`** refinement validates the six WCs. Every deferred-go check has **exactly
one** in-wave owner (gating WC-5's seal). Completion constituents are consistent (`evidence_bound` listed)
with a typed resolver-input contract. Ledger advanced to R9; content-hash chain re-pinned. **No Python.**

**Go + CUE repository (operator directive, R8).** Schemas are **CUE**; procedural code is **Go**. A prior
round shipped a hand-rolled Python validator (`validate.py`) — a tool-choice error, **removed entirely**.
Pre-authorization validation is:

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

## Validation architecture (R9)

### Structural — enforced NOW by `cue vet`

Run `make -C schema clean` (or the per-artifact invocations in [`schema/README.md`](./schema/README.md)).
R9 makes `#CellContract` the **faithful** canonical §2 shape and validates the six Working-Cell contracts
against the named **`#WorkingCellContract`** refinement (class=working, kind=artifact); the canonical
variants exercise `#CellContract` directly:

```
cue vet ./schema/ ./wave.cn-wave-v1.yaml                      -d '#Wave'
cue vet ./schema/ ./contracts/wc-{1,2,3a,3b,4,5}.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
cue vet ./schema/ ./oracle-registry.yaml                      -d '#AssuranceRegistry'
cue vet ./schema/ ./intent.cn-intent-v1.yaml                  -d '#Intent'
cue vet ./schema/ ./regressions/contract.canonical-{nullable-scope,optional-prior-receipt,relation-graph,judgment}.yaml -d '#CellContract'
```

The wave-envelope + completion + §2-drift blockers close natively (full tables in
[`schema/README.md`](./schema/README.md)): nullable `scope.wave`/`parent_cell`; 1+ `inputs.required`; a
`sibling_output` in `inputs.optional` → conflict; the `prior_receipt` locator class present; the full
`artifact | relation_graph | judgment` output enum; 1+ `acceptance.predicates`/`allowed_paths`; required
`forbidden_paths`/`non_goals`; the `gates.reason` truth table (nonempty iff a gate bool true, null iff
both false); `receipt_bound: "false"` (string) → type conflict; a record or `required_constituents`
missing `evidence_bound` → required-field fail; duplicate node/edge/mapping keys → conflicting values;
empty `external_roots`/`stop_conditions`; a deleted/flipped `edge_derivation` const; a mutated completion
formula; an added top-level key; a doctrine-affecting cell not requiring acceptance → gate-invariant conflict.

The [`schema/regressions/`](./schema/regressions/) suite proves each rejection: the clean/canonical bases
PASS and **27 one-mutation bad fixtures are each rejected** by `cue vet` (`make -C schema regressions`;
verified green with cue v0.17.1) — 15 wave + 12 contract, incl. the 7 new §2-drift negatives (empty
`inputs.required`; a `sibling_output` in optional; empty `acceptance.predicates`/`allowed_paths`; missing
`forbidden_paths`/`non_goals`; false/false gates with a non-null reason). This is **not** a blanket
soundness certificate — `cue vet` warrants exactly the closed-struct / typed-field / enum / cardinality /
required-field / pinned-const / truth-table constraints the schemas declare, and **nothing beyond them**.
Everything else is named-and-deferred to single-owner Go validators.

### Procedural / semantic — named-and-deferred to SINGLE-OWNER Go validators

R9: each deferred-go check has **exactly one** in-wave `deferred_owner` (never #627 — #627 S2–S3 stay
downstream consumers/canonicalizers). Each owner contract carries a gating acceptance predicate whose PASS
gates its completion; WC-5 depends (sibling_output) on wc-1/wc-2/wc-3b, so **WC-5 cannot seal until every
deferred validator exists and passes**.

| Deferred-Go check (the owner's Go validator runs it when that WC executes) | Single owner | Go artifact |
|---|---|---|
| git **ref / content-hash resolution** | **WC-2** | `.cdd/unreleased/wc-2/validators/ref_resolve.go` |
| graph **acyclicity** (DAG) | **WC-3b** | `.cdd/unreleased/wc-3b/validators/wave_dag.go` |
| sibling-output → wave-edge **parity** | **WC-3b** | `.cdd/unreleased/wc-3b/validators/edge_parity.go` |
| parallel **write-surface disjointness** | **WC-3b** | `.cdd/unreleased/wc-3b/validators/write_surface.go` |
| **oracle-ownership bijection** | **WC-1** | `.cdd/unreleased/wc-1/validators/oracle_ownership_bijection.go` |
| **completion-evidence derivation** (typed resolver input → 5 derived constituents) | **WC-5** | `.cdd/unreleased/wc-5/validators/completion_evidence.go` |
| **ledger consistency** (revision markers agree; per-category counts) | **WC-5** | `.cdd/unreleased/wc-5/validators/ledger_consistency.go` |
| **classification-totality bijection** | **WC-5** | `.cdd/unreleased/wc-5/validators/classification_bijection.go` |

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

All six `contracts/*.yaml` are canonical `cn.cell.contract.v1` §2 instances. `#CellContract` is the
**faithful** canonical §2 shape (nullable `scope` refs; 1+ `inputs.required`; external-only
`inputs.optional`; `repo_artifact | control_plane | prior_receipt` locator union; `artifact |
relation_graph | judgment` output kinds; required `forbidden_paths`/`non_goals`; gate truth-table; closed
struct). The six Working-Cell contracts are validated against the **named refinement**
`#WorkingCellContract = #CellContract & { cell.class: "working", requested_output.kind: "artifact" }` —
the canonical shape is **not** narrowed; the 4 canonical-variant fixtures validate directly against
`#CellContract`. Oracle ownership lives in the separate content-bound [`oracle-registry.yaml`](./oracle-registry.yaml),
which each contract consumes via a canonical `inputs.required[].external` control_plane ref (immutable
locator + content hash). Reverse-consumer information is **derived** (see [`reconcile-627.md`](./reconcile-627.md)).

## Non-recursive, structured whole-wave completion

`wave.cn-wave-v1.yaml` `completion` defines a **non-recursive construction set**
`N = {wc-1, wc-2, wc-3a, wc-3b, wc-4}` (excludes wc-5), a terminal `wc-5`, and a **pinned governing
formula const** `wave_complete: "all(N child_complete) AND child_complete(terminal)"` (required + fixed
by `#Wave` — deleting or mutating it fails `cue vet`). Each fixture's per-node record is a **closed
5-bool struct** (`requested_output_produced`, `acceptance_all_pass`, `class_v_pass`, `receipt_bound`,
`evidence_bound` — every constituent mandatory and typed bool). **R9 (Required 3):** `required_constituents`
now lists all **5** (matching `#ChildRecord`; `#ChildCompleteDef` pins the list and definition consistently),
and a typed **resolver-input contract** `#CompletionEvidenceInput` (node/output identity + content binding,
acceptance result set, V verdict/receipt binding, β/γ evidence locator, receipt identity/hash) is contracted
in `evidence_resolver` (owned by WC-5): the 5 booleans are **derived** by resolving those bindings, not
author-supplied; the fixture `records` are truth-table fixtures. The **derivation** — child_complete(n) AND,
predicate-graph acyclicity, and each fixture computing its authored `expected` — is the WC-5 completion-evidence
**deferred-Go** check, whose PASS binds into WC-5 closure. Four fixtures pin the intended behavior
(all-incomplete; all-N-complete/WC-5-absent; WC-5 V-FAIL; all-complete).

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
is its projection. The categories (R9 taxonomy):

- **structural-cue** — enforced NOW by `cue vet` against a plan-local `#Def`.
- **deferred-go** — a named Go validator with **exactly one** in-wave `deferred_owner` (never #627; #627
  S2–S3 stay downstream consumers/canonicalizers). Procedural; not run at this tree.
- **mechanically-verifiable** — a per-child acceptance oracle the child WC **emits** as a **Go checker
  or a CUE schema** (**no `.py`**) + named fixtures + a credential-free command, bound in its receipt.
- **evidenced** — bound receipt evidence, not self-report.
- **cognitive-review** — an independent β / CC doctrine judgment, honestly not mechanical.

**Child acceptance predicates (78, TOTAL):** structural-cue **1** · deferred-go **19** ·
mechanically-verifiable **30** · evidenced **7** · cognitive-review **21** — one classification each,
matching the 78 predicates across the six contracts (13+14+12+15+9+15). R9 added the single-owner
deferred-validator **gating predicates** (wc-1 oracle-ownership; wc-2 ref-resolution; wc-3b DAG/parity/
write-surface; wc-5 completion-evidence/ledger/classification + upstream-consumption), each a `deferred-go`
entry with a single `deferred_owner`, whose PASS gates the owner's completion. Wave-only structural/deferred
predicates are held **separately** in `wave_predicates:` (with the 8 canonical validators' full specs) so
they can neither inflate nor mask child coverage.

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

## Per-finding disposition (R9 — faithful §2 CUE; single-owner deferred validators; consistent completion)

| # | Finding | R9 disposition |
|---|---|---|
| 1 | **[BLOCKER]** `#CellContract` did not faithfully encode canonical §2 | Rewrote `schema/cell_contract.cue`: nullable `scope` refs; 1+ `inputs.required`; external-only `inputs.optional`; `repo_artifact\|control_plane\|prior_receipt` locator union; `artifact\|relation_graph\|judgment` output kinds; required `forbidden_paths`/`non_goals`; the `gates.reason` truth table (CUE conditionals); closed through `stop_conditions`. Added the named `#WorkingCellContract` refinement (the 6 WCs validate against it; 4 canonical variants against `#CellContract`; 7 new §2-drift negatives rejected). |
| 2 | **[REQUIRED]** deferred Go validators had multi/slash owners | Every deferred-go check has **exactly one** `deferred_owner`: wc-2 ref-resolution; wc-3b DAG/parity/write-surface; wc-1 oracle-ownership; wc-5 completion-evidence/ledger/classification. The 8 canonical validators pin artifact id+path, typed inputs, result/evidence shape, positive+named negative fixtures, downstream consumers, and a gating predicate; each owner contract carries the gating acceptance predicate. #627 S2–S3 stay downstream. **WC-5 cannot seal until every deferred validator passes.** |
| 3 | **[REQUIRED]** completion constituents inconsistent; 5 booleans author-supplied | `required_constituents` now lists all **5** (added `evidence_bound`), pinned consistent with the definition and `#ChildRecord` by `#ChildCompleteDef`; added the typed `#CompletionEvidenceInput` resolver-input contract + `#EvidenceResolver` (owned by WC-5); the 5 booleans are derived, the fixtures labelled truth-table fixtures; resolver PASS binds into WC-5 closure. |
| 4 | Ledger + provenance | Every `wave-revision:`/`revision:` marker advanced to **R9**; the content-hash chain re-pinned (registry → 6 contracts → wave `contract_sha256`); grounding/reconcile/intent unchanged. Contracts stay exactly §2 (validated by the faithful CUE). **No Python.** |

## Prior-round ledger (unchanged accepted structure)

R2/R3 repaired the external-β ITERATEs (non-recursive completion, honest intent provenance, byte-exact
grounding); R4/R5 moved oracle ownership into the content-bound `oracle-registry.yaml`; R6 made it TOTAL;
R7 closed the wave-envelope + completion blockers; R8 removed Python and moved structural validation to
CUE; **R9** makes `#CellContract` faithfully canonical, gives every deferred validator a single in-wave
owner, and makes completion constituents consistent with a typed resolver input. The accepted decomposition
graph (WC-2 → WC-1 → {WC-3a, WC-3b, WC-4} → WC-5), D9-four, intent provenance, and grounding are
**unchanged** across R2–R9.

## Files

- [`README.md`](./README.md) — this overview.
- [`schema/`](./schema/) — plan-local transitional CUE (`#CellContract` + `#WorkingCellContract`, `#Wave`,
  `#AssuranceRegistry`, `#Intent`), the `Makefile` runner, and `regressions/` (clean/canonical bases + 27
  rejected mutation fixtures).
- [`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) — transitional bootstrap intent projection.
- [`decision-provenance.md`](./decision-provenance.md) — α/β planning conclusions + the R8/R9 dispositions.
- [`wave.cn-wave-v1.yaml`](./wave.cn-wave-v1.yaml) — the `cn.wave.v1` instance.
- [`contracts/wc-1..wc-5.cn-cell-contract-v1.yaml`](./contracts/) — the six §2-conformant child contracts.
- [`grounding-source-5015460988.md`](./grounding-source-5015460988.md) — byte-exact source snapshot (SHA `9d1ab3a5…`).
- [`grounding-cm.md`](./grounding-cm.md) — α's derivative FAIL-classification (references the snapshot).
- [`reconcile-627.md`](./reconcile-627.md) — the #627 S1–S8 refinement map + derived reverse-consumers.
- [`oracle-registry.yaml`](./oracle-registry.yaml) — AUTHORITATIVE, TOTAL assurance registry (five honest categories).
- [`acceptance-oracles.md`](./acceptance-oracles.md) — projection of the total registry (every predicate classified single-kind).

---
*Status: R9 wave, α matter — under operator review; external-β and CC review next. No child WCs
dispatched; no control-plane action taken by this cell.*
