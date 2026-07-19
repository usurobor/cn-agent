<!-- wave-revision: R7 -->
# Acceptance oracles — honestly classified, registry-projected (cnos#671 R7)

Every load-bearing acceptance predicate (wave + WC-1..WC-5) is classified as **exactly one** of four
honest kinds. Nothing cognitive is dressed up as mechanical; a semantic-absence claim is **never**
implemented as grep-and-called-mechanical. **Every category value is a single enum member** — no
compound `enforced + evidenced` rows (`validate.py` check (i) fails closed on a compound cell).

**This document is a PROJECTION.** The authoritative machine-readable source is now the **TOTAL**
[`oracle-registry.yaml`](./oracle-registry.yaml) — a separate content-bound registry that each child
contract consumes via a canonical `inputs.required[].external` control_plane ref (no non-canonical
contract key). Its `assurance:` list carries an entry for **every** scalar in **every** child contract's
`acceptance.predicates`, classified into **exactly one** of the four categories with that category's
required fields; wave-only enforced predicates live in a **separate** `wave_predicates:` section.

- `validate.py` check **(h)** proves a total + singular bijection over the **mechanically-verifiable
  subset** (registry ⇄ each child's `acceptance.predicates` ⇄ the *Registry projection* table below).
- `validate.py` check **(j)** proves a **total + singular bijection over the COMPLETE child acceptance
  set**: `union(acceptance.predicates)` ⇄ the child-owned `assurance:` entries, each classified **exactly
  once** across all four categories — rejecting a predicate present in a contract but absent from the
  registry (the coordinated-omission the pre-R6 checks missed), a double-classified predicate, a phantom
  entry, a bad category, or a break against the *Complete assurance classification* table below.

The R4 regression — carrying oracle ownership as an extra top-level `oracles` key that broke the exact §2
contract shape — remains repaired: the contracts are exact §2, and an added top-level key FAILS check (a).

## Classification kinds

- **enforced** — checked by `validate.py` **now**, at this wave tree (credential-free; exits non-zero
  on violation). The wave-level pre-authorization predicates are all here.
- **mechanically-verifiable** — a **real command** over a **concretely-located checker/schema + named
  structured fixtures the child WC must produce**, owned by **exactly one** child contract in
  `oracle-registry.yaml` at concrete repo-root-relative paths (NO placeholders, NO implicit CWD), in
  that contract's `constraints.allowed_paths` (emitted) or a pinned immutable input (existing), required
  by the contract's `acceptance.predicates`, and bound (bytes + results) into the receipt via the child's
  `emits_machine_readable_acceptance_fixtures_into_receipt__…` predicate. The child WC emits the
  checkers/schemas/fixtures under `.cdd/unreleased/wc-N/fixtures/` (its allowed-path); they do not exist
  yet — they are the mechanical acceptance contract each Working Cell inherits and must emit.
- **evidenced** — verified by **bound evidence in the receipt**, not self-report (e.g. a reproduction
  diff, a receipt binding the emitted fixtures). Not a static credential-free command.
- **cognitive-review** — an **independent β / CC doctrine judgment**, honestly **NOT mechanical**. A
  human/independent reviewer decides it (e.g. "no clause equates CM with a CC judgment", "exactly one
  κ≠α logic", "exactly one authority per FSM edge"). Grep can *locate* text; it cannot *decide* these.

**Summary counts (child acceptance predicates):** enforced **11** · mechanically-verifiable **30** ·
evidenced **7** · cognitive-review **21** — **69** total, one classification each. (These are the
per-category totals over the **69** child acceptance predicates, derived from the TOTAL
`oracle-registry.yaml` `assurance:` list; the mechanically-verifiable count equals the 30 ownership
entries. `validate.py` check (i) recomputes these from the registry and fails closed on any drift, and
check (j) proves the registry bijects with `union(acceptance.predicates)` exactly once each.) The **12**
wave-only enforced predicates (checks (a)–(k) + the mutation harness) are held **separately** in
`wave_predicates:` so they can neither inflate nor mask child coverage; the wave-level table just below
documents them for the reader.

## How the mechanically-verifiable rows become real commands — owned in the registry

Every checker, schema, and fixture is **owned by exactly one child contract** in `oracle-registry.yaml`
and lives at a **concrete path** under that WC's fixtures dir `.cdd/unreleased/wc-N/fixtures/`
(repo-root-relative; no placeholder operands; no implicit CWD). Two ownership modes:

- **emitted** — the WC emits the checker/schema itself (`python3 .cdd/unreleased/wc-N/fixtures/CHECKER.py`,
  or `cue vet FIXTURE.json .cdd/unreleased/wc-N/fixtures/NAME.schema.cue` for a fixture schema the WC
  writes — never editing shipped `schemas/**`). The checker/schema and both fixtures are in the
  contract's `allowed_paths`.
- **existing** — the command vets against an **existing shipped schema at a pinned commit** (e.g.
  `schemas/cdd/receipt.cue`) that is a pinned immutable input of the owning contract; only the
  positive/negative fixtures are emitted.

The commands are credential-free: `cue vet FIXTURE.json SCHEMA.cue` or `python3 CHECKER.py FIXTURE.json`.
The registry is authoritative; this file is its human-readable projection, parity-checked against it by
check (h) (see the *Registry projection* table at the bottom).

---

## Wave-level — ENFORCED now by `validate.py` / `validate_test.py`

| Predicate | Kind | Oracle (command) | Expected positive / negative | Evidence |
|---|---|---|---|---|
| §2 constraint model (key paths + enums + types + cardinalities); **no extra top-level key** | enforced | `python3 validate.py` (a) | clean tree exit 0 / a `cell.class: nonsense`, a non-bool gate, an added top-level `oracles` key → exit 1 | check-(a) PASS |
| immutable refs resolve (intent id/schema; git `commit:path`; content hashes incl. the oracle-registry ref; source `9d1ab3a5…`) | enforced | `validate.py` (b) | 0 / wrong `intent_ref.id`, `docs/.../DOES-NOT-EXIST.md`, a stale registry content hash → exit 1 | check-(b) PASS |
| edge parity (authored == sibling-output-derived) | enforced | `validate.py` (c) | 0 / add/remove an edge or change an `output_id` → exit 1 | derived==authored |
| dependency graph is a DAG | enforced | `validate.py` (d) | 0 / a cycle in `sibling_output` refs → exit 1 | topo covers all |
| parallel nodes share no write surface | enforced | `validate.py` (e) | 0 / two concurrent nodes overlap `allowed_paths` → exit 1 | per-pair disjoint |
| gate invariants | enforced | `validate.py` (f) | 0 / `doctrine_affecting:true` + `operator_acceptance_required:false`, or missing `reason` → exit 1 | check-(f) PASS |
| completion-predicate graph acyclic (built by AST walk) | enforced | `validate.py` (g) | 0 / tautology `wave_complete := wave_complete` → cycle → exit 1 | predicate-DAG topo |
| completion **fixtures evaluate to their authored `expected`** | enforced | `validate.py` (g) | 0 / a flipped `expected` (computed≠authored) → exit 1 | per-fixture eval |
| **oracle-ownership bijection** — total + singular map (registry ⇄ each child's `acceptance` ⇄ this projection); every mechanically-verifiable predicate owned exactly once | enforced | `validate.py` (h) | 0 / a removed entry, a duplicate owner, an unowned mv predicate, an extra owner absent from acceptance, a reclassified entry, a placeholder, or a projection parity break → exit 1 | check-(h) PASS |
| **ledger consistency** — revision labels agree; reported counts == total-registry child totals == mv ownership size; every category a single enum member | enforced | `validate.py` (i) | 0 / a disagreeing revision label, a miscount, or a compound category → exit 1 | check-(i) PASS |
| **classification totality + singularity** — `union(acceptance.predicates)` ⇄ the total `assurance:` registry, every child predicate classified exactly once across all four categories | enforced | `validate.py` (j) | 0 / a predicate in a contract but absent from the registry (a coordinated omission), a double-classified/phantom entry, a bad category, or a projection parity break → exit 1 | check-(j) PASS |
| **full `cn.wave.v1` envelope** — exact top-level keys + `schema` const + types; node shape + output_path/output_id parity; edge key-set + `kind` enum; sole-root + critical-path (root→terminal); typed gates; non-empty typed STOP; completion shape; **duplicate mapping keys rejected** | enforced | `validate.py` (k) | 0 / `schema: nonsense`, `stop_conditions: []`, `gates: {}`, a truncated critical path, a node `output_path` diverging from its contract, an `edge.kind: nonsense`, or a duplicate `schema:` key → exit 1 | check-(k) PASS |
| the twenty-nine adversarial mutations each fail for their own predicate; clean passes | enforced | `python3 validate_test.py` | harness exit 0 (29 fail, clean passes) | harness receipt |

---

## WC-2 — Coherence-Measurement (`COHERENCE-MEASUREMENT.md`) · fixtures `.cdd/unreleased/wc-2/fixtures/`

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| CM producer is a provider/runtime op (e.g. tsc), **not** CC | mechanically-verifiable | `python3 check_cm_producer.py cm-object.positive.json` | valid CM accepted / producer = cohering-cell → reject | cm-object fixture |
| CM object fields typed (`s_α,s_β,s_γ,C_Σ,bottleneck,witnesses,defects,provenance,…`) | mechanically-verifiable | `cue vet cm-object.positive.json cm-object.schema.cue` (emitted schema) | full CM accepted (exit 0) / `cm-object.missing-bottleneck.negative.json` → reject | cm-object.schema.cue + pos/neg |
| CM and CC-judgment are **distinct tagged objects** (CC consumes immutable `cm_ref`) | mechanically-verifiable | `python3 check_cm_cc_distinct.py cm-cc-distinct.positive.json` | two distinct tags accepted / one object tagged both → reject | two-fixture log |
| CM→V edge pinned (V consumes a named immutable `cm_ref`) | mechanically-verifiable | `python3 check_v_signature.py v-signature.positive.json` | V input carries `cm_ref` / V invocation with no `cm_ref` → reject | V-signature fixture |
| `receipt_core → CM → V → δ → final_receipt` type path | mechanically-verifiable | `cue vet final-receipt.positive.json schemas/cdd/receipt.cue` (existing shipped schema) | final validates / a "final" receipt with no V/δ outputs → reject | round-trip receipts |
| CM realized **within** the four schemas (typed CM field/edge + `cm_ref` resolving within them) | mechanically-verifiable | `cue vet cm-ref-embedded.positive.json receipt-with-cmref.schema.cue` (emitted schema) | `cm_ref` resolves within the four schemas / `cm-ref-needs-fifth.negative.json` → reject | receipt-with-cmref.schema.cue |
| negative fixture: reject a **fifth** canonical `cn.cm.v1` schema | mechanically-verifiable | `python3 reject_fifth.py four-schema.positive.json` | four-schema proposal accepted / a fifth-canonical-schema proposal → reject | reject-fifth fixture |
| `cm_ref` field shape is the single interface all nodes import (no second shape) | mechanically-verifiable | `python3 cm_ref_parity.py cm-ref-parity.positive.json` | identical `cm_ref` shape across nodes / a divergent downstream `cm_ref` → reject | cm_ref-parity log |
| CM is a runner/provider object and **not** the CC's judgment (no clause equates them) | cognitive-review | independent β/CC read | β confirms no clause collapses CM into a CC judgment | β judgment |
| CM producers **and** consumers pinned by role (V gates; CC judgment consumes refs; PC wave grounded in it) | cognitive-review | independent β/CC read | β confirms the role assignment is coherent doctrine | β judgment |
| D9 four-schema boundary **settled** (operator decision; no reopening) | cognitive-review | operator/CC read | operator-settled; not a mechanical property of the artifact | operator matter |
| reproducible (same frozen inputs + provider version → same CM) | evidenced | re-run `coh` on the frozen bundle digest (needs the provider) | reproduction diff = ∅ bound in receipt / a differing CΣ under identical inputs | reproduction evidence |
| emits machine-readable fixtures into receipt (cm-object, cc-judgment, receipt-core/final, reject-fifth) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / a missing fixture → V FAIL | receipt evidence |

---

## WC-1 — Cell classes & typed contract (`CELL-RUNTIME-CLASSES.md`) · fixtures `.cdd/unreleased/wc-1/fixtures/`

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| `cn.cell.contract.v1` envelope = §2 key paths (template ≡ worked instance) | mechanically-verifiable | `python3 keypath_diff.py contract-template.yaml contract-worked-instance.yaml` | normalized key-path sets identical (diff = ∅) / an extra/missing top-level key → non-empty diff → reject | key-path diff |
| input reference model is the provenance-tagged union | mechanically-verifiable | `python3 union_validate.py input-union.positive.yaml` | each `inputs.required[]` is `external{locator}` or `sibling_output{producer,output_id}` / an untyped/opaque optional scalar → reject | union-validate log |
| `cm_ref` imported verbatim from WC-2 (within four schemas; no fifth; no re-derivation) | mechanically-verifiable | `python3 cm_ref_parity.py cm-ref-parity.positive.json` | WC-1 `cm_ref` == WC-2 output shape / a re-derived or fifth-schema `cm_ref` → reject | cm_ref-parity |
| gate invariants hold on the worked instance | enforced | `validate.py` (f) | as wave-level (f) | check-(f) |
| classes defined by canonical output telos (WC=artifact, PC=relation_graph, CC=judgment); closes L0-B1 | cognitive-review | independent β/CC read | β confirms the class↔telos rule is correct doctrine | β judgment |
| class-specific V predicates reference an imported `cm_ref`; WC-1 inlines no CM internals | cognitive-review | independent β read | β confirms no CM internals leaked into WC-1 | β judgment |
| κ≠α role logic unconditional; State-A = hosting-identity collapse, not role-equality (carries CC-1) | cognitive-review | independent β/CC read | β confirms exactly one κ≠α logic, State-A framing honest | β judgment |
| review evidence content-bound not identity-bound (carries CC-2) | cognitive-review | independent β read | β confirms the §11.6 content-bound-review rule is carried | β judgment |
| emits machine-readable fixtures into receipt (template, worked instance, input-union) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / a missing fixture → V FAIL | receipt evidence |

---

## WC-3a — Cell FSM (`CELL-FSM.md`) · fixtures `.cdd/unreleased/wc-3a/fixtures/`

The cell FSM is emitted as a machine-readable state table `cell-fsm-state-table.json`
(`{states, transitions:[{from,to,event,guard,authority}], initial, terminals}`).

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| state set closed & declared (`blocked` is a declared member, not only a transition target) | mechanically-verifiable | `python3 fsm_closure.py cell-fsm-state-table.positive.json` | every transition target ∈ declared states / `blocked` reachable but undeclared → reject | closure report |
| total event/guard/action relation (every non-terminal state has owned exits) | mechanically-verifiable | `python3 fsm_totality.py cell-fsm-state-table.positive.json` | every non-terminal has ≥1 exit / a non-terminal with no exit → reject | totality report |
| deterministic / normatively-prioritized transition selection (same bundle hash → same decision) | mechanically-verifiable | `python3 fsm_determinism.py cell-fsm-state-table.positive.json` | no unresolved (state,event) ambiguity / two equal-priority exits → reject | determinism report |
| CC-disposition→transition-request adapter (CC never mutates state) | mechanically-verifiable | `cue vet cc-transition-request.positive.json cell-transition-request.schema.cue` (emitted schema) | CC emits a request the FSM validates / `cc-mutates-state.negative.json` → reject | schema + pos/neg |
| command-table parity (`cn cell` / `cn issues` map onto declared edges) | mechanically-verifiable | `python3 fsm_command_parity.py cell-fsm-state-table.positive.json command-map.positive.json` | each command maps to a declared edge / a command with no edge → reject | parity table |
| reachability proved (BFS from initial covers all states) | mechanically-verifiable | `python3 fsm_reach.py cell-fsm-state-table.positive.json` | reachable set == declared states / an unreachable declared state → reject | reachability set |
| exactly one authority per edge (reconciles `transitions.json` / dispatch / resume / spec) | cognitive-review | independent β read | β confirms each edge has a single owning authority | β judgment |
| invalid-transition semantics defined | cognitive-review | independent β read | β confirms the invalid-transition rule is coherent | β judgment |
| review-return conflict resolved (§11.1 `→todo` vs §11.2 `→in-progress` settled to one edge) | cognitive-review | independent β read | β confirms one settled edge | β judgment |
| request-marker pattern pinned (typed marker + evidentiary guards) | cognitive-review | independent β read | β confirms the marker/guard pattern is carried | β judgment |
| emits machine-readable fixtures into receipt (cell-fsm-state-table, command-map) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## WC-3b — Wave FSM (`WAVE-FSM.md`) · fixtures `.cdd/unreleased/wc-3b/fixtures/`

Emitted machine-readable: `wave-fsm-state-table.json` and `disposition-to-transition-map.json`.

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| wave state set closed (enumerated states = declared set) | mechanically-verifiable | `python3 fsm_closure.py wave-fsm-state-table.positive.json` | every target ∈ declared set / a target outside the set → reject | closure report |
| total event/guard/action relation (distinct state space from the cell FSM) | mechanically-verifiable | `python3 fsm_totality.py wave-fsm-state-table.positive.json` | every non-terminal has owned exits / a non-terminal with no exit → reject | totality report |
| total mapping of all **8** CC dispositions → wave transition requests | mechanically-verifiable | `python3 disposition_totality.py disposition-to-transition-map.positive.json` | all 8 dispositions mapped / a disposition with no transition → reject | mapping table |
| child-receipt aggregation rule (child completion predicates aggregate; CC does not mutate state) | mechanically-verifiable | `cue vet wave-aggregation-request.positive.json wave-transition-request.schema.cue` (emitted schema) | aggregation request validates / `cc-mutates-wave-state.negative.json` → reject | schema + pos/neg |
| wave-transition-request shape pinned (mirrors the cell request-marker pattern + guards) | mechanically-verifiable | `cue vet wave-transition-request.positive.json wave-transition-request.schema.cue` (emitted schema) | request with `{from,to,guard}` validates / `wave-request-missing-guard.negative.json` → reject | schema + pos/neg |
| child & whole-wave completion predicates **derived** (not a contract field); acyclic | enforced | `validate.py` (g) | as wave-level (g) / a per-contract `completion_signal` field or recursive completion → exit 1 | check-(g) |
| holding/replanning exits & terminal semantics defined | cognitive-review | independent β read | β confirms exits and terminals coherent | β judgment |
| invalid-transition semantics defined | cognitive-review | independent β read | β confirms the invalid-transition rule | β judgment |
| dependency & gate dispatch rules (wave-boundary auth once; children scheduled by deps) | cognitive-review | independent β read | β confirms §9 wave-boundary policy carried | β judgment |
| separate authority from the cell FSM (two independently total machines) | cognitive-review | independent β read | β confirms two distinct total machines, not one conflated | β judgment |
| emits machine-readable fixtures into receipt (wave-fsm-state-table, disposition map) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## WC-4 — Shipped→specified migration (`CONTRACT-MIGRATION.md`) · fixtures `.cdd/unreleased/wc-4/fixtures/`

Emitted machine-readable: `field-map.json`, `round-trip/*.json`, and five `negative/*.json`.

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| field-by-field preservation (every shipped field mapped/added/deprecated) | mechanically-verifiable | `python3 field_coverage.py field-map.positive.json schemas/cdd/contract.cue` | every shipped field present with a status / an unmapped shipped field → reject | coverage report |
| round-trip oracle (lossless paths proven lossless; lossy points named) | mechanically-verifiable | `python3 round_trip.py round-trip.positive.json` | declared-lossless path round-trips identically / a "lossless" path that loses a field → reject | round-trip log |
| negative fixtures (reject unknown field / dropped required / mutable-root input / identity-only revision / missing adapter case) | mechanically-verifiable | `python3 run_negatives.py negatives-manifest.positive.json` | all 5 rejected / any of the 5 accepted → the oracle fails | 5 rejection receipts |
| CM edge representation consumed from WC-2 (target `cm_ref` == WC-2 shape; within four schemas) | mechanically-verifiable | `python3 cm_ref_parity.py cm-ref-parity.positive.json` | WC-4 target `cm_ref` == WC-2 shape / a re-derived or fifth-schema `cm_ref` → reject | cm_ref-parity |
| versioned projection specified (`cnos.cdd.contract.v1 → cn.cell.contract.v1`, named version + direction) | cognitive-review | independent β read | β confirms the projection is named + directed | β judgment |
| State-A/State-B field enumeration (shipped commands/fields vs target fields) | cognitive-review | independent β read | β confirms the enumeration is complete/honest | β judgment |
| migration is doctrine/schema-boundary only (no CUE/Go adapter implemented here) | cognitive-review | independent β read | β confirms no adapter code authored | β judgment |
| emits machine-readable fixtures into receipt (field-map, round-trip, 5 negatives) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## WC-5 — Integration seal (`CELL-RUNTIME-INTEGRATION.md`) · fixtures `.cdd/unreleased/wc-5/fixtures/`

Emitted machine-readable: `reconcile-627-audit.json`, `fail-disposition-audit.json`,
`vocab-consistency.json`, `schema-fixture-log.json`.

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| graph parity (derived == authored edges) | enforced | `validate.py` (c) | as wave-level (c) / a spurious/missing edge → exit 1 | check-(c) |
| schema fixtures pass (every worked instance validates against §2) | mechanically-verifiable | `python3 schema_fixture_run.py schema-fixture-log.positive.json` | every worked instance validates / a non-conformant instance → reject | fixture pass log |
| reconcile-627 mapping complete (every S-node classified exactly once, no double owner) | mechanically-verifiable | `python3 reconcile_audit.py reconcile-627-audit.positive.json` | each S-node classified once / a doubly-owned or omitted S-node → reject | #627 map audit |
| residual debt disposed exactly once (each grounding FAIL closed-by / retired / tracked once) | mechanically-verifiable | `python3 fail_disposition_audit.py fail-disposition-audit.positive.json` | each FAIL disposed once / a FAIL disposed twice or not at all → reject | FAIL-disposition audit |
| migration round-trips (WC-4 round-trip oracle green on named fixtures) | mechanically-verifiable | `python3 round_trip.py migration-round-trip.positive.json` | WC-4 round-trips green / a failing round-trip → reject | round-trip log |
| wc5 readiness over non-recursive N (receipt exists for each of wc-1,wc-2,wc-3a,wc-3b,wc-4) | enforced | `validate.py` (g) N-set | N derived correctly (g) / a recursive N including wc-5 → exit 1 | check-(g) |
| wc5 own completion non-recursive (own output+acceptance+V+receipt; no self/wave quantification) | enforced | `validate.py` (g) | (g) proves non-recursion / a wc-5 completion reading `wave_complete` → (g) exit 1 | check-(g) |
| cross-artifact vocabulary consistent (one `cm_ref` vocab, one CM object, one κ≠α logic across WC-1..WC-4) | cognitive-review | independent β/CC read | β confirms one vocabulary, one κ≠α logic (semantic consistency, not grep) | β judgment |
| final integration V fails-closed on a conflicting `cm_ref`/FSM-state/schema-field/missing-adapter/unmapped-627-node | cognitive-review | independent β/CC read | β confirms the seal fails closed on each conflict class | β judgment |
| emits machine-readable fixtures into receipt (reconcile audit, FAIL audit, vocab, schema-fixture log) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## Complete assurance classification — projection of the total registry (check (j))

This table is the human-readable PROJECTION of the **complete** `assurance:` list: **one row per every**
child acceptance predicate (all six contracts, 69 rows), each classified into **exactly one** category.
`validate.py` check (j) parses it and proves it is set-identical to the registry (same `(owner, predicate)`
pairs, same classification) AND that this set equals `union(acceptance.predicates)` exactly once each.
A predicate left in a contract but dropped from the registry + this table (a coordinated omission) FAILS (j).

| owner | predicate | classification |
|---|---|---|
| wc-1 | `wc_pc_cc_classes_defined_by_canonical_output_telos__wc_artifact_pc_relation_graph_cc_judgment__closes_L0_B1` | cognitive-review |
| wc-1 | `cn_cell_contract_v1_envelope_specified__cell_id_class_mode_protocol_matter_domain__inputs_required_provenance_tagged_union__requested_output_id_kind_path__acceptance__constraints__gates__doctrine_affecting__stop_conditions` | cognitive-review |
| wc-1 | `input_reference_model_is_the_provenance_tagged_union__external_locator_classed_or_sibling_output_producer_output_id` | mechanically-verifiable |
| wc-1 | `class_specific_v_predicates_reference_an_imported_cm_ref__wc1_does_not_inline_cm_internals` | cognitive-review |
| wc-1 | `cm_ref_field_shape_is_imported_verbatim_from_wc2_output_realized_within_the_four_schemas__no_re_derivation_no_fifth_schema` | mechanically-verifiable |
| wc-1 | `kappa_neq_alpha_role_logic_unconditional__state_a_is_hosting_identity_collapse_not_role_equality__carries_repaired_CC_1_as_settled_input` | cognitive-review |
| wc-1 | `review_evidence_content_bound_not_identity_bound__carries_repaired_CC_2_as_settled_input` | cognitive-review |
| wc-1 | `gate_invariants_hold__doctrine_affecting_true_implies_operator_acceptance_required_true__reason_present_iff_a_gate_bool_true` | enforced |
| wc-1 | `key_path_parity__template_and_worked_instance_normalize_to_identical_key_path_sets` | mechanically-verifiable |
| wc-1 | `scope_guard__no_cm_internals_no_cell_or_wave_fsm_no_migration_authored_here` | cognitive-review |
| wc-1 | `emits_machine_readable_acceptance_fixtures_into_receipt__contract_template_yaml__worked_instance_yaml__input_union_fixtures__each_validated_by_the_named_acceptance_oracle_command` | evidenced |
| wc-1 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | enforced |
| wc-2 | `cm_defined_as_runner_or_provider_produced_object__producer_is_cm_provider_or_runtime_op_eg_tsc_not_cc` | mechanically-verifiable |
| wc-2 | `cm_object_fields_typed__s_alpha_s_beta_s_gamma_c_sigma_geom_mean_degeneracy_bottleneck_witnesses_defects_mode_target_bundle_digest_thresholds_standing_provenance` | mechanically-verifiable |
| wc-2 | `cc_class_result_is_a_separate_tagged_object_that_consumes_immutable_cm_refs__cm_and_cc_judgment_validate_as_distinct_tagged_objects` | mechanically-verifiable |
| wc-2 | `cm_to_v_edge_pinned__v_consumes_an_immutable_cm_ref_as_a_named_input_the_frozen_contract_times_receipt_signature_is_extended_to_carry` | mechanically-verifiable |
| wc-2 | `receipt_core_to_cm_to_v_to_delta_to_final_receipt_type_path_declared__provisional_receipt_core_without_validation_or_boundary_decision_then_measurement_then_final_receipt_satisfying_shipped_receipt_cue` | mechanically-verifiable |
| wc-2 | `cm_producers_and_consumers_pinned__producer_provider_runtime_op__consumers_v_gates_on_it_cc_judgment_consumes_refs_pc_wave_grounded_in_it` | cognitive-review |
| wc-2 | `d9_four_schema_boundary_settled__cm_is_realized_within_the_authorized_four_schemas_not_as_a_fifth_canonical_cn_cm_v1__operator_settled_no_authorized_reopening` | cognitive-review |
| wc-2 | `cm_realized_within_four_schemas__a_typed_cm_field_or_edge_in_the_receipt_plus_a_cm_ref_resolving_within_the_existing_four_schemas__wc1_and_wc4_consume_this_same_cm_ref_shape` | mechanically-verifiable |
| wc-2 | `negative_fixture_reject_fifth_canonical_schema__a_proposed_fifth_cn_cm_v1_canonical_schema_is_rejected_while_the_four_schema_boundary_remains_settled` | mechanically-verifiable |
| wc-2 | `reproducible__same_frozen_inputs_plus_provider_version_reproduce_the_mechanically_claimed_cm_result` | evidenced |
| wc-2 | `cm_ref_field_shape_is_the_single_interface_all_other_nodes_import__no_second_cm_shape_introduced_downstream` | mechanically-verifiable |
| wc-2 | `emits_machine_readable_acceptance_fixtures_into_receipt__cm_object_json__cc_judgment_json__receipt_core_json__final_receipt_json__reject_fifth_schema_json__each_validated_by_the_named_acceptance_oracle_command` | evidenced |
| wc-2 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | enforced |
| wc-3a | `cell_fsm_state_set_closed_and_declared__blocked_is_a_declared_member_not_only_a_transition_target` | mechanically-verifiable |
| wc-3a | `exactly_one_authority_per_edge__reconciles_transitions_json_and_cn_issues_dispatch_and_cn_cell_resume_and_spec` | cognitive-review |
| wc-3a | `total_event_guard_action_relation__every_nonterminal_state_has_owned_exits` | mechanically-verifiable |
| wc-3a | `invalid_transition_semantics_defined` | cognitive-review |
| wc-3a | `deterministic_or_normatively_prioritized_transition_selection__same_bundle_hash_same_decision` | mechanically-verifiable |
| wc-3a | `review_return_conflict_resolved__the_11_1_changes_to_todo_vs_11_2_changes_to_in_progress_ambiguity_settled_to_one_edge` | cognitive-review |
| wc-3a | `cc_disposition_to_cell_transition_request_adapter__cc_selects_a_judgment_the_fsm_validates_and_effects_the_transition__cc_never_mutates_state_directly` | mechanically-verifiable |
| wc-3a | `command_table_parity__cn_cell_and_cn_issues_surfaces_map_onto_declared_table_edges` | mechanically-verifiable |
| wc-3a | `request_marker_pattern_pinned__typed_marker_plus_evidentiary_guards` | cognitive-review |
| wc-3a | `reachability_proved__every_declared_state_reachable_from_the_initial_state` | mechanically-verifiable |
| wc-3a | `emits_machine_readable_acceptance_fixtures_into_receipt__cell_fsm_state_table_json__for_state_set_closure_totality_and_reachability_oracles` | evidenced |
| wc-3a | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | enforced |
| wc-3b | `wave_fsm_state_set_closed__intent_planning_required_d0_review_wave_planning_wave_review_dispatchable_executing_holding_replanning_completing_complete_complete_with_residuals` | mechanically-verifiable |
| wc-3b | `total_event_guard_action_relation__distinct_from_the_cell_fsm_state_space_and_events` | mechanically-verifiable |
| wc-3b | `invalid_transition_semantics_defined` | cognitive-review |
| wc-3b | `child_receipt_aggregation_rule__child_completion_predicates_aggregate_into_wave_state_without_cc_mutating_state` | mechanically-verifiable |
| wc-3b | `total_mapping_all_eight_cc_dispositions_to_wave_transition_requests__request_planning_request_working_hold_request_human_continue_wave_complete_complete_with_residuals_block` | mechanically-verifiable |
| wc-3b | `dependency_and_gate_dispatch_rules__wave_boundary_authorization_once_then_children_scheduled_by_dependencies_no_per_child_operator_gate` | cognitive-review |
| wc-3b | `holding_and_replanning_exits_and_terminal_semantics_defined` | cognitive-review |
| wc-3b | `wave_transition_request_shape_pinned__mirrors_the_cell_fsm_request_marker_pattern_with_guards` | mechanically-verifiable |
| wc-3b | `separate_authority_from_the_cell_fsm__two_independently_total_machines_one_shared_transition_algebra_or_two_tagged_instances` | cognitive-review |
| wc-3b | `child_and_whole_wave_completion_predicates_defined__derived_from_child_acceptance_v_receipt_surfaces_not_a_new_contract_field` | enforced |
| wc-3b | `emits_machine_readable_acceptance_fixtures_into_receipt__wave_fsm_state_table_json__and_eight_disposition_to_transition_map_json__for_totality_and_mapping_oracles` | evidenced |
| wc-3b | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | enforced |
| wc-4 | `versioned_projection_specified__cnos_cdd_contract_v1_to_cn_cell_contract_v1_with_a_named_version_and_direction` | cognitive-review |
| wc-4 | `field_by_field_preservation_rules__every_shipped_field_mapped_or_explicitly_marked_added_or_deprecated` | mechanically-verifiable |
| wc-4 | `cm_edge_representation_consumed_from_wc2__the_target_contract_cm_ref_shape_is_wc2_output_realized_within_the_four_schemas_not_a_fifth_schema_not_re_derived` | mechanically-verifiable |
| wc-4 | `state_a_state_b_field_enumeration__shipped_commands_and_schema_fields_vs_target_fields_enumerated_per_field` | cognitive-review |
| wc-4 | `round_trip_oracle__lossless_paths_proven_lossless_and_lossy_points_named` | mechanically-verifiable |
| wc-4 | `negative_fixtures__reject_unknown_field_dropped_required_field_mutable_root_input_identity_only_revision_missing_adapter_case` | mechanically-verifiable |
| wc-4 | `migration_is_doctrine_schema_boundary_only__no_cue_or_go_adapter_implemented_here` | cognitive-review |
| wc-4 | `emits_machine_readable_acceptance_fixtures_into_receipt__field_map_json__round_trip_fixtures__five_negative_fixtures__for_field_preservation_and_round_trip_oracles` | evidenced |
| wc-4 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | enforced |
| wc-5 | `cross_artifact_vocabulary_consistent__one_canonical_input_ref_and_requested_output_vocabulary_one_cm_object_one_kappa_neq_alpha_logic_across_wc1_through_wc4` | cognitive-review |
| wc-5 | `schema_fixtures_pass__every_worked_instance_validates_against_the_selected_cn_cell_contract_v1_shape` | mechanically-verifiable |
| wc-5 | `graph_parity__wave_edges_equal_the_sibling_output_relations_exactly_no_missing_or_spurious_edge` | enforced |
| wc-5 | `migration_round_trips__wc4_round_trip_oracle_is_green_on_the_named_fixtures` | mechanically-verifiable |
| wc-5 | `reconcile_627_mapping_complete__every_S_node_classified_exactly_once_no_duplicate_owner_no_omitted_executable_child` | mechanically-verifiable |
| wc-5 | `residual_debt_disposed__every_grounding_fail_is_closed_by_or_retired_instance_defect_or_tracked_follow_up_exactly_once` | mechanically-verifiable |
| wc-5 | `final_integration_v_fails_closed_on__a_conflicting_cm_ref_a_conflicting_fsm_state_a_conflicting_schema_field_a_missing_adapter_case_or_an_unmapped_627_node` | cognitive-review |
| wc-5 | `wc5_readiness_over_non_recursive_construction_set_N__a_completion_receipt_exists_for_every_node_in_N_where_N_equals_wc1_wc2_wc3a_wc3b_wc4_and_excludes_wc5` | enforced |
| wc-5 | `wc5_own_completion_is_non_recursive__wc5_requested_output_produced_plus_wc5_acceptance_pass_plus_integration_v_pass_plus_bound_receipt__does_not_quantify_over_wc5_itself_nor_over_the_whole_wave_predicate` | enforced |
| wc-5 | `emits_machine_readable_acceptance_fixtures_into_receipt__reconcile_627_audit_json__fail_disposition_audit_json__vocab_consistency_json__schema_fixture_log_json__for_the_integration_seal_oracles` | evidenced |
| wc-5 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | enforced |

---

## Registry projection — parity-checked against `oracle-registry.yaml` (check (h))

This table is the human-readable PROJECTION of the 30 authoritative `oracle-registry.yaml` entries.
`validate.py` check (h) parses it and proves it is byte-for-field identical to the registry (same set of
predicates; same owner / classification / ownership / checker-or-schema / positive / negative per row).
Removing a registry entry, adding a projection row with no registry entry, reclassifying, or diverging any
field FAILS check (h).

| owner | predicate | classification | ownership | checker_or_schema | positive_fixture | negative_fixture |
|---|---|---|---|---|---|---|
| wc-1 | `key_path_parity__template_and_worked_instance_normalize_to_identical_key_path_sets` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-1/fixtures/keypath_diff.py` | `.cdd/unreleased/wc-1/fixtures/contract-worked-instance.yaml` | `.cdd/unreleased/wc-1/fixtures/contract-extra-key.negative.yaml` |
| wc-1 | `input_reference_model_is_the_provenance_tagged_union__external_locator_classed_or_sibling_output_producer_output_id` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-1/fixtures/union_validate.py` | `.cdd/unreleased/wc-1/fixtures/input-union.positive.yaml` | `.cdd/unreleased/wc-1/fixtures/input-union-opaque-scalar.negative.yaml` |
| wc-1 | `cm_ref_field_shape_is_imported_verbatim_from_wc2_output_realized_within_the_four_schemas__no_re_derivation_no_fifth_schema` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-1/fixtures/cm_ref_parity.py` | `.cdd/unreleased/wc-1/fixtures/cm-ref-parity.positive.json` | `.cdd/unreleased/wc-1/fixtures/cm-ref-rederived.negative.json` |
| wc-2 | `cm_defined_as_runner_or_provider_produced_object__producer_is_cm_provider_or_runtime_op_eg_tsc_not_cc` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/check_cm_producer.py` | `.cdd/unreleased/wc-2/fixtures/cm-object.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-producer-cc.negative.json` |
| wc-2 | `cm_object_fields_typed__s_alpha_s_beta_s_gamma_c_sigma_geom_mean_degeneracy_bottleneck_witnesses_defects_mode_target_bundle_digest_thresholds_standing_provenance` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/cm-object.schema.cue` | `.cdd/unreleased/wc-2/fixtures/cm-object.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-object.missing-bottleneck.negative.json` |
| wc-2 | `cc_class_result_is_a_separate_tagged_object_that_consumes_immutable_cm_refs__cm_and_cc_judgment_validate_as_distinct_tagged_objects` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/check_cm_cc_distinct.py` | `.cdd/unreleased/wc-2/fixtures/cm-cc-distinct.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-cc-conflated.negative.json` |
| wc-2 | `cm_to_v_edge_pinned__v_consumes_an_immutable_cm_ref_as_a_named_input_the_frozen_contract_times_receipt_signature_is_extended_to_carry` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/check_v_signature.py` | `.cdd/unreleased/wc-2/fixtures/v-signature.positive.json` | `.cdd/unreleased/wc-2/fixtures/v-signature-no-cmref.negative.json` |
| wc-2 | `receipt_core_to_cm_to_v_to_delta_to_final_receipt_type_path_declared__provisional_receipt_core_without_validation_or_boundary_decision_then_measurement_then_final_receipt_satisfying_shipped_receipt_cue` | mechanically-verifiable | existing | `schemas/cdd/receipt.cue` | `.cdd/unreleased/wc-2/fixtures/final-receipt.positive.json` | `.cdd/unreleased/wc-2/fixtures/final-receipt-no-v-delta.negative.json` |
| wc-2 | `cm_realized_within_four_schemas__a_typed_cm_field_or_edge_in_the_receipt_plus_a_cm_ref_resolving_within_the_existing_four_schemas__wc1_and_wc4_consume_this_same_cm_ref_shape` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/receipt-with-cmref.schema.cue` | `.cdd/unreleased/wc-2/fixtures/cm-ref-embedded.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-ref-needs-fifth.negative.json` |
| wc-2 | `negative_fixture_reject_fifth_canonical_schema__a_proposed_fifth_cn_cm_v1_canonical_schema_is_rejected_while_the_four_schema_boundary_remains_settled` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/reject_fifth.py` | `.cdd/unreleased/wc-2/fixtures/four-schema.positive.json` | `.cdd/unreleased/wc-2/fixtures/fifth-schema.negative.json` |
| wc-2 | `cm_ref_field_shape_is_the_single_interface_all_other_nodes_import__no_second_cm_shape_introduced_downstream` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/cm_ref_parity.py` | `.cdd/unreleased/wc-2/fixtures/cm-ref-parity.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-ref-divergent.negative.json` |
| wc-3a | `cell_fsm_state_set_closed_and_declared__blocked_is_a_declared_member_not_only_a_transition_target` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_closure.py` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-blocked-undeclared.negative.json` |
| wc-3a | `total_event_guard_action_relation__every_nonterminal_state_has_owned_exits` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_totality.py` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-nonterminal-noexit.negative.json` |
| wc-3a | `deterministic_or_normatively_prioritized_transition_selection__same_bundle_hash_same_decision` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_determinism.py` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-ambiguous.negative.json` |
| wc-3a | `cc_disposition_to_cell_transition_request_adapter__cc_selects_a_judgment_the_fsm_validates_and_effects_the_transition__cc_never_mutates_state_directly` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/cell-transition-request.schema.cue` | `.cdd/unreleased/wc-3a/fixtures/cc-transition-request.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cc-mutates-state.negative.json` |
| wc-3a | `command_table_parity__cn_cell_and_cn_issues_surfaces_map_onto_declared_table_edges` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_command_parity.py` | `.cdd/unreleased/wc-3a/fixtures/command-map.positive.json` | `.cdd/unreleased/wc-3a/fixtures/command-map-orphan.negative.json` |
| wc-3a | `reachability_proved__every_declared_state_reachable_from_the_initial_state` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_reach.py` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-unreachable.negative.json` |
| wc-3b | `wave_fsm_state_set_closed__intent_planning_required_d0_review_wave_planning_wave_review_dispatchable_executing_holding_replanning_completing_complete_complete_with_residuals` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/fsm_closure.py` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-target-outside.negative.json` |
| wc-3b | `total_event_guard_action_relation__distinct_from_the_cell_fsm_state_space_and_events` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/fsm_totality.py` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-nonterminal-noexit.negative.json` |
| wc-3b | `total_mapping_all_eight_cc_dispositions_to_wave_transition_requests__request_planning_request_working_hold_request_human_continue_wave_complete_complete_with_residuals_block` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/disposition_totality.py` | `.cdd/unreleased/wc-3b/fixtures/disposition-to-transition-map.positive.json` | `.cdd/unreleased/wc-3b/fixtures/disposition-missing.negative.json` |
| wc-3b | `child_receipt_aggregation_rule__child_completion_predicates_aggregate_into_wave_state_without_cc_mutating_state` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/wave-transition-request.schema.cue` | `.cdd/unreleased/wc-3b/fixtures/wave-aggregation-request.positive.json` | `.cdd/unreleased/wc-3b/fixtures/cc-mutates-wave-state.negative.json` |
| wc-3b | `wave_transition_request_shape_pinned__mirrors_the_cell_fsm_request_marker_pattern_with_guards` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/wave-transition-request.schema.cue` | `.cdd/unreleased/wc-3b/fixtures/wave-transition-request.positive.json` | `.cdd/unreleased/wc-3b/fixtures/wave-request-missing-guard.negative.json` |
| wc-4 | `field_by_field_preservation_rules__every_shipped_field_mapped_or_explicitly_marked_added_or_deprecated` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/field_coverage.py` | `.cdd/unreleased/wc-4/fixtures/field-map.positive.json` | `.cdd/unreleased/wc-4/fixtures/field-map-unmapped.negative.json` |
| wc-4 | `round_trip_oracle__lossless_paths_proven_lossless_and_lossy_points_named` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/round_trip.py` | `.cdd/unreleased/wc-4/fixtures/round-trip.positive.json` | `.cdd/unreleased/wc-4/fixtures/round-trip-lossy.negative.json` |
| wc-4 | `negative_fixtures__reject_unknown_field_dropped_required_field_mutable_root_input_identity_only_revision_missing_adapter_case` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/run_negatives.py` | `.cdd/unreleased/wc-4/fixtures/negatives-manifest.positive.json` | `.cdd/unreleased/wc-4/fixtures/negatives-regression.negative.json` |
| wc-4 | `cm_edge_representation_consumed_from_wc2__the_target_contract_cm_ref_shape_is_wc2_output_realized_within_the_four_schemas_not_a_fifth_schema_not_re_derived` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/cm_ref_parity.py` | `.cdd/unreleased/wc-4/fixtures/cm-ref-parity.positive.json` | `.cdd/unreleased/wc-4/fixtures/cm-ref-rederived.negative.json` |
| wc-5 | `schema_fixtures_pass__every_worked_instance_validates_against_the_selected_cn_cell_contract_v1_shape` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/schema_fixture_run.py` | `.cdd/unreleased/wc-5/fixtures/schema-fixture-log.positive.json` | `.cdd/unreleased/wc-5/fixtures/schema-fixture-nonconformant.negative.json` |
| wc-5 | `reconcile_627_mapping_complete__every_S_node_classified_exactly_once_no_duplicate_owner_no_omitted_executable_child` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/reconcile_audit.py` | `.cdd/unreleased/wc-5/fixtures/reconcile-627-audit.positive.json` | `.cdd/unreleased/wc-5/fixtures/reconcile-627-double-owner.negative.json` |
| wc-5 | `residual_debt_disposed__every_grounding_fail_is_closed_by_or_retired_instance_defect_or_tracked_follow_up_exactly_once` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/fail_disposition_audit.py` | `.cdd/unreleased/wc-5/fixtures/fail-disposition-audit.positive.json` | `.cdd/unreleased/wc-5/fixtures/fail-disposed-twice.negative.json` |
| wc-5 | `migration_round_trips__wc4_round_trip_oracle_is_green_on_the_named_fixtures` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/round_trip.py` | `.cdd/unreleased/wc-5/fixtures/migration-round-trip.positive.json` | `.cdd/unreleased/wc-5/fixtures/migration-round-trip-lossy.negative.json` |

---

*Rows marked **cognitive-review** are honestly not mechanical: an independent β (outside the Sigma
lineage) or a CC doctrine judgment decides them. They are surfaced here so the reviewer knows exactly
which predicates carry warrant by judgment vs. by machine.*
