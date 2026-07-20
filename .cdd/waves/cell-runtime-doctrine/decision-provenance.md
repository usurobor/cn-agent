# Decision provenance — α/β planning conclusions for the cell-runtime wave (cnos#671)

This file holds the **planning conclusions** of this Planning Cell — the matter that is **α/β work**,
**not operator intent**. Finding-2 (R3) moved these out of `intent.cn-intent-v1.yaml` because they
had been mislabeled there as operator/κ pre-cell intent. Each conclusion below cites its **actual**
source (α round, external-β ITERATE, or settled doctrine). The durable operator intent lives in
[`intent.cn-intent-v1.yaml`](./intent.cn-intent-v1.yaml) (statement + carriers) and, authoritatively,
in the operator's wave master **#627**.

## What is operator intent vs. what is α/β planning matter

| Matter | Where it lives | Authority |
|---|---|---|
| The objective + the operator's verbatim final doctrine line | `intent.cn-intent-v1.yaml` `statement`; **#627** | Operator |
| Settled S1 doctrine (one kernel, WC/PC/CC telos classes; CM measures / V gates / δ effects) | `docs/architecture/CELL-RUNTIME.md` (#628, landed) | Settled doctrine |
| The D9 four-schema boundary being **settled** (no reopening) | operator decision, carried in `intent` `out_of_scope` | Operator |
| **The WC decomposition graph** (node set, edges, roots, critical path) | **this file** + `wave.cn-wave-v1.yaml` | **α (R1), shape accepted by external-β** |
| **The D9 realization detail** (CM realized *within* the four schemas; the `cm_ref` shape) | **this file** + WC-2/WC-1/WC-4 acceptance | **α**, consistent with settled S1 |
| **The β-repair outcomes** (R2/R3 dispositions) | **this file** | **external-β ITERATE #672** + this α (R2/R3) |

## 1. The WC decomposition graph — α conclusion (shape accepted by external-β)

**Source:** authored by α in R1 (`cnos#671` R1 commit); the **graph shape was accepted by the
external-β ITERATE (#672)** and is FIXED for R2/R3 (no shape change). This is α's decomposition of the
operator objective into single-purpose Working-Cell contracts — it is **planning matter, not intent**.

Accepted, fixed shape:

```
external roots (immutable): grounding-CM(@sha 9d1ab3a5) · #628/S1 · shipped schemas/CCNF/transitions.json
  → WC-2  (CM measurement object + receipt_core→CM→V→δ→final_receipt type path)   [keystone root]
  → WC-1  (cell classes + typed cell contract; imports cm_ref from WC-2)
      → WC-3a (cell FSM) · WC-3b (wave FSM) · WC-4 (shipped→specified migration)   [each dep WC-1 + WC-2]
          → WC-5  (integration / seal: whole-wave proof)                          [dep WC-1,2,3a,3b,4]
```

- **Node set** `N ∪ {wc-5}` where `N = {wc-1, wc-2, wc-3a, wc-3b, wc-4}`.
- **Edges** are the mechanical projection of the child contracts' `sibling_output` refs (an edge
  `A→B` exists iff B carries a `sibling_output` ref resolving to A's `requested_output.id`);
  `external` refs create no edge. The deferred-Go edge-parity validator (WC-3b/WC-5) proves authored == derived parity.
- **Roots:** `wc-2` (sole keystone; every `cm_ref` imports from its output).
- **Critical path:** `wc-2 → wc-1 → wc-3b → wc-5`.
- **WC-5 (integration seal)** was **required by the external-β** (four/six child closures do not by
  themselves prove the objective); it has no #627 S-counterpart.
- **#662 → WC-1 reclassification** (a converged PC-D0 artifact re-landing as a WC-doctrine node
  because its canonical output is a repo artifact) is α's classification per the settled telos rule.

The node/edge data itself is authored, mechanically, in `wave.cn-wave-v1.yaml`; this section records
that it is **α planning matter** and names its acceptance source.

## 2. The D9 realization detail — α conclusion (consistent with settled S1)

**Source:** α (R1/R2). The operator settled that D9's four-schema boundary stays **fixed** (that is the
operator constraint, in `intent` `out_of_scope`). **How** CM is realized *within* those four schemas is
α's design conclusion:

- CM is realized as a **typed CM field/edge in the receipt** plus a **`cm_ref`** that resolves within
  the existing four schemas — **not** a fifth canonical `cn.cm.v1` schema.
- WC-2 owns the `cm_ref` interface shape; **WC-1 and WC-4 consume the same shape** (no second CM shape,
  no re-derivation). WC-2 carries a reject-a-fifth-canonical-schema negative fixture.

This realization is α's, aligned with the landed S1 doctrine ("CM measures; V gates; δ effects",
`CELL-RUNTIME.md` #628) and the operator's settled four-schema boundary. It is **not** an operator
decision beyond "four schemas stay settled."

## 3. The β-repair outcomes — external-β ITERATE (#672) + this α (R2/R3)

**Source:** the external-β ITERATE on PR/issue #672 (findings), and this Planning Cell's α repairs
(R2 = the six exact-contract/assurance repairs; R3 = this repair). These are **review→repair
outcomes**, not intent.

### R2 dispositions (external-β ITERATE #672)

| # | Finding | R2 disposition |
|---|---|---|
| 1 | **[BLOCKER]** whole-wave completion was recursive | Non-recursive `N`; WC-5 readiness over N; WC-5 completion over its own surfaces; whole-wave = all(N) AND wc5_complete; predicate DAG. |
| 2 | **[BLOCKER]** child contracts not valid §2 instances | Removed non-canonical `consumers`/`completion_signal`; all six normalize to §2 key paths; materialized the intent projection. |
| 3 | **[REQUIRED]** grounding called an abridgment "verbatim" | Byte-exact `grounding-source-5015460988.md` (SHA `9d1ab3a5…`); `grounding-cm.md` honestly a derivative. |
| 4 | **[REQUIRED]** acceptance not mechanically decidable; no validator | `validate.py` (a–g) + `acceptance-oracles.md`. |
| 5 | **[REQUIRED]** no wave STOP conditions; `contract_ref` not immutable | Typed wave `stop_conditions`; revision-bound `contract_ref_resolution`; per-node `contract_sha256`. |
| 6 | **[REQUIRED]** WC-2 4-vs-5 schema contradiction | D9 four-schema kept settled; CM realized within (see §2 above). |

### R3 dispositions (this repair)

| # | Finding | R3 disposition |
|---|---|---|
| 1 | **[BLOCKER]** `validate.py` false-passed five adversarial mutations | Rewrote the validator to be genuinely SOUND: full §2 constraint model (enums/types/cardinalities), real ref resolution (intent id/schema compared to the intent object; every repo-artifact locator resolved with `git cat-file -e`; grounding source hash verified), derivation from the authored wave/contracts (nodes/edges/roots/critical-path derived and checked), and **evaluation** of the authored completion predicates + truth-table fixtures as structured data (real graph acyclicity + per-fixture computation vs. `expected`). Added `validate_test.py` materializing the five adversarial mutations; each now exits non-zero for its own predicate while the clean tree exits 0. |
| 2 | **[BLOCKER]** intent masqueraded as pre-cell κ/operator intent | Rewrote `intent.cn-intent-v1.yaml` as an explicitly **transitional bootstrap projection** authored during the cell (no pre-cell existence claim); statement carries only operator matter (#627 + verbatim doctrine line); α conclusions moved to **this file**; identity vs carrier kept distinct. |
| 3 | **[REQUIRED]** acceptance oracles mislabeled cognitive review as mechanical | Rewrote `acceptance-oracles.md`: every predicate classified as exactly one of **enforced** / **mechanically-verifiable** (named fixture + command + expected outputs the child WC must emit) / **evidenced** / **cognitive-review** (honestly not mechanical). Semantic-absence claims are cognitive-review, not grep. Mechanically-verifiable predicates require the child contract to emit the named fixture into its receipt. |

### R8 disposition — operator directive: Go + CUE repo, NO Python (validation re-architecture)

**Source:** an **operator directive** delivered for R8, plus this Planning Cell's α re-architecture.
The directive is operator matter; the re-architecture is α matter. Recorded here (not in `intent`,
which carries the objective) because it is a **tooling/method** decision, not the objective itself.

**Operator directive (verbatim intent):** *this is a Go + CUE repository; no Python; CUE is for
schemas, Go is for procedural code.* The prior rounds (R3–R7) shipped a hand-rolled Python validator
(`validate.py` / `validate_test.py`) as the pre-authorization checker — a **tool-choice error**.

**R8 α re-architecture:**

| # | Directive | R8 disposition |
|---|---|---|
| 1 | Go + CUE repo; **no Python** | Deleted `validate.py` and `validate_test.py`; no `*.py` remains under the wave dir. |
| 2 | Structural validation moved to **CUE** | Authored plan-local, transitional closed CUE (`schema/#CellContract`, `#Wave`, `#AssuranceRegistry`, `#Intent`) — the design input WC-1/WC-2 formalize into canonical `schemas/cdd/*.cue`. Wired the exact `cue vet` invocations (`schema/Makefile`, `schema/README.md`). The two external-β R7 blockers (`bool("false")` coercion, duplicate node/edge invisibility, empty `external_roots`, deletable/literal completion formula, non-mandatory `evidence_bound`, un-typed authority substructure) are closed **natively** by CUE (typed fields, closed structs, cardinality, pinned consts, required fields, uniqueness folds), proven by 20 rejected regression fixtures. |
| 3 | Procedural checks deferred to **Go** | DAG acyclicity, sibling-output edge parity, git ref/content-hash resolution, the classification bijection, and completion-evidence derivation are **named** deliverables of WC-1/WC-2 and #627 S2–S3 (D9, option B) — not reimplemented here, and never in Python. |
| 4 | Honest scoping | The assurance registry's `enforced_by` fields no longer point at `validate.py`; each points at the CUE constraint (`cue vet` against a `#Def`) or the named deferred Go validator. The blanket "genuinely SOUND / fail-closed" claim is **removed** — `cue vet` warrants exactly the declared structural constraints; everything procedural is deferred. Categories reclassified: `structural-cue` (1 child) · `deferred-go` (10) · `mechanically-verifiable` (30, now Go/CUE — no `.py`) · `evidenced` (7) · `cognitive-review` (21). |

The accepted decomposition (WC-2 → WC-1 → {WC-3a, WC-3b, WC-4} → WC-5), the §2 contracts, D9-four,
grounding, and intent provenance are **unchanged** — R8 is a validation-tooling re-architecture only.

### R9 disposition — faithful canonical §2 CUE + single-owner deferred validators + consistent completion

**Source:** the external-β ITERATE on #672 (findings) and this Planning Cell's α repair (R9). Review→repair
outcomes, not intent. The accepted six-node graph, D9-four, grounding, and intent provenance are **unchanged**.

| # | Finding | R9 disposition |
|---|---|---|
| 1 | **[BLOCKER]** `#CellContract` did not faithfully encode canonical §2 (drifts β found) | Rewrote `schema/cell_contract.cue` so `#CellContract` is the **exact** `cn.cell.contract.v1` §2 shape: `scope.wave`/`scope.parent_cell` `string \| null`; `inputs.required` 1+ of the provenance-tagged union; `inputs.optional` external-locator refs **only** (a `sibling_output` in optional is rejected); the external locator union carries `repo_artifact \| control_plane \| **prior_receipt**`; `requested_output.kind` = `artifact \| relation_graph \| judgment`; `acceptance.predicates` 1+, `allowed_paths` 1+, `forbidden_paths`/`non_goals` required present keys (0+); `gates.reason` an always-present `string \| null` obeying the **truth table** (nonempty iff a gate bool true, null iff both false) via CUE conditionals; closed through `stop_conditions`. Added a named **`#WorkingCellContract`** refinement (class=working, kind=artifact) — the canonical shape is NOT narrowed; the 6 real WCs validate against the refinement, 4 canonical variants (nullable scope, optional `prior_receipt`, `relation_graph`, `judgment`) validate against `#CellContract`, and 7 new §2-drift negatives are rejected (plus all prior regressions). |
| 2 | **[REQUIRED]** deferred Go validators had multi/slash owners ("WC-3b/WC-5", "#627 S2-S3") | Every `deferred-go` check now has **exactly one** in-wave `deferred_owner` (a single-valued enum; #627 is not a member and stays a downstream consumer/canonicalizer): **wc-2** ref/content-hash resolver; **wc-3b** DAG + edge-parity + write-surface; **wc-1** oracle-ownership bijection; **wc-5** completion-evidence + ledger + classification-totality bijection. Each of the 8 canonical wave-level validators pins a Go artifact id+path, typed inputs, a result/evidence shape, positive + named negative fixtures, downstream consumers, and a `gating_predicate`; each owner contract carries that gating acceptance predicate (PASS gates the owner's completion) + the validator path in `allowed_paths`. WC-5 depends (sibling_output) on wc-1/wc-2/wc-3b, so **WC-5 cannot seal until every deferred validator exists and passes**. |
| 3 | **[REQUIRED]** completion constituents inconsistent; the 5 booleans were author-supplied | `required_constituents` now lists all **5** (added `evidence_bound`), and `#ChildCompleteDef` (CUE) pins both the list and the definition to the canonical constituent set (a dropped constituent is rejected NOW). Added a typed **resolver-input contract** `#CompletionEvidenceInput` (node/output identity + content binding, acceptance result set, V verdict/receipt binding, β/γ evidence locator, receipt identity/hash) + `#EvidenceResolver` (owned by WC-5): the 5 booleans are **derived** by resolving these bindings, not author-supplied; the fixture `records` are labelled truth-table fixtures; the resolver PASS binds into WC-5 closure. |

**Ledger:** every `wave-revision:`/`revision:` marker advanced to **R9**; the content-hash chain re-pinned
(oracle-registry → 6 contracts → wave `contract_sha256`); grounding/reconcile/intent files unchanged so
their hashes stay. Contracts remain exactly §2 (now validated by the faithful CUE). **No Python.**

## Coordination-index note (κ / control-plane, not this cell's matter)

Recording this provenance on an immutable coordination index (an update to #627 or a named index
issue) is a **control-plane action owned by κ**, not this Planning Cell. This cell authors the matter
(this file); it files, dispatches, comments, and merges nothing.
