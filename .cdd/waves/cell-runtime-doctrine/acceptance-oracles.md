# Acceptance oracles — honestly classified (cnos#671 R3)

Repairs external-β finding 4 and the R3 finding-3: **every** load-bearing acceptance predicate (wave
+ WC-1..WC-5) is classified as **exactly one** of four honest kinds. Nothing cognitive is dressed up
as mechanical; a semantic-absence claim is **never** implemented as grep-and-called-mechanical.

## Classification kinds

- **enforced** — checked by `validate.py` **now**, at this wave tree (credential-free; exits non-zero
  on violation). The wave-level pre-authorization predicates are all here.
- **mechanically-verifiable** — a **real command** over a **named structured fixture the child WC must
  produce**. Each row gives the command, the fixture path, and the expected **positive** and
  **negative** outputs. These run against the child's fixtures (which do not exist yet — they are the
  mechanical acceptance contract each Working Cell inherits and must emit, per its
  `emits_machine_readable_acceptance_fixtures_into_receipt__…` predicate and its
  `.cdd/unreleased/<wc>/fixtures/**` allowed-path).
- **evidenced** — verified by **bound evidence in the receipt**, not self-report (e.g. a reproduction
  diff, a receipt binding the emitted fixtures). Not a static credential-free command.
- **cognitive-review** — an **independent β / CC doctrine judgment**, honestly **NOT mechanical**. A
  human/independent reviewer decides it (e.g. "no clause equates CM with a CC judgment", "exactly one
  κ≠α logic", "exactly one authority per FSM edge"). Grep can *locate* text; it cannot *decide* these.

**Summary counts:** enforced **9** · mechanically-verifiable **28** · evidenced **8** ·
cognitive-review **17**. (Per-row below; every acceptance predicate in the six contracts and every
wave-level predicate is assigned exactly one kind.)

## How the mechanically-verifiable rows become real commands

The child WC emits the named fixtures into `.cdd/unreleased/<wc>/fixtures/` (allowed_path) and binds
them in its receipt (`emits_machine_readable_acceptance_fixtures_into_receipt__…` acceptance
predicate). The commands below are then credential-free: `cue vet <fixture> <schema>` for the shipped
`schemas/cdd/*.cue` surfaces, or `python3 <check>` over a JSON/YAML fixture. "Expected negative" is the
mutated fixture the command must reject; the WC's own V runs these and binds PASS/FAIL.

---

## Wave-level — ENFORCED now by `validate.py` / `validate_test.py`

| Predicate | Kind | Oracle (command) | Expected positive / negative | Evidence |
|---|---|---|---|---|
| §2 constraint model (key paths + enums + types + cardinalities) | enforced | `python3 validate.py` (a) | clean tree exit 0 / `cell.class: nonsense`, a non-bool gate, an under-full seq → exit 1 | check-(a) PASS |
| immutable refs resolve (intent id/schema; git `commit:path`; content hashes; source `9d1ab3a5…`) | enforced | `validate.py` (b) | 0 / wrong `intent_ref.id`, `docs/.../DOES-NOT-EXIST.md`, corrupted snapshot → exit 1 | check-(b) PASS |
| edge parity (authored == sibling-output-derived) | enforced | `validate.py` (c) | 0 / add/remove an edge or change an `output_id` → exit 1 | derived==authored |
| dependency graph is a DAG | enforced | `validate.py` (d) | 0 / a cycle in `sibling_output` refs → exit 1 | topo covers all |
| parallel nodes share no write surface | enforced | `validate.py` (e) | 0 / two concurrent nodes overlap `allowed_paths` → exit 1 | per-pair disjoint |
| gate invariants | enforced | `validate.py` (f) | 0 / `doctrine_affecting:true` + `operator_acceptance_required:false`, or missing `reason` → exit 1 | check-(f) PASS |
| completion-predicate graph acyclic (built by AST walk) | enforced | `validate.py` (g) | 0 / tautology `wave_complete := wave_complete` → cycle → exit 1 | predicate-DAG topo |
| completion **fixtures evaluate to their authored `expected`** | enforced | `validate.py` (g) | 0 / a flipped `expected` (computed≠authored) → exit 1 | per-fixture eval |
| the five adversarial mutations each fail for their own predicate; clean passes | enforced | `python3 validate_test.py` | harness exit 0 (5 fail, clean passes) | harness receipt |

---

## WC-2 — Coherence-Measurement (`COHERENCE-MEASUREMENT.md`) · fixtures `.cdd/unreleased/wc-2/fixtures/`

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| CM producer is a provider/runtime op (e.g. tsc), **not** CC | mechanically-verifiable | `cue vet cm-object.json schemas/cdd/*.cue` (producer.kind ∈ provider set) | valid CM accepted / producer = cohering-cell → reject | cm-object fixture |
| CM object fields typed (`s_α,s_β,s_γ,C_Σ,bottleneck,witnesses,defects,provenance,…`) | mechanically-verifiable | `cue vet cm-object.json <cm-schema>` | full CM accepted / CM missing `bottleneck`/`provenance` → reject | cm-object fixture |
| CM and CC-judgment are **distinct tagged objects** (CC consumes immutable `cm_ref`) | mechanically-verifiable | `cue vet cm-object.json + cc-judgment.json` as distinct tags | two distinct tags accepted / one object tagged both → reject | two-fixture log |
| CM→V edge pinned (V consumes a named immutable `cm_ref`) | mechanically-verifiable | `python3 check_v_signature.py cm-object.json` | V input carries `cm_ref` / V invocation with no `cm_ref` → reject | V-signature fixture |
| `receipt_core → CM → V → δ → final_receipt` type path | mechanically-verifiable | `cue vet receipt-core.json receipt.cue`; `cue vet final-receipt.json receipt.cue` | core (no validation/boundary_decision) + final validate / a "final" receipt with no V/δ outputs → reject | round-trip receipts |
| CM realized **within** the four schemas (typed CM field/edge + `cm_ref` resolving within them) | mechanically-verifiable | `cue vet cm-ref-embedded.json <extended receipt.cue>` | `cm_ref` resolves within the four schemas / a `cm_ref` needing a fifth schema → reject | four-schema fixture |
| negative fixture: reject a **fifth** canonical `cn.cm.v1` schema | mechanically-verifiable | `python3 reject_fifth.py reject-fifth-schema.json` | four-schema proposal accepted / a fifth-canonical-schema proposal → reject | reject-fifth fixture |
| `cm_ref` field shape is the single interface all nodes import (no second shape) | mechanically-verifiable | `python3 cm_ref_parity.py wc2 wc1 wc4` | identical `cm_ref` shape across nodes / a divergent downstream `cm_ref` → reject | cm_ref-parity log |
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
| input reference model is the provenance-tagged union | mechanically-verifiable | `python3 union_validate.py input-union-*.yaml` | each `inputs.required[]` is `external{locator}` or `sibling_output{producer,output_id}` / an untyped/opaque optional scalar → reject | union-validate log |
| `cm_ref` imported verbatim from WC-2 (within four schemas; no fifth; no re-derivation) | mechanically-verifiable | `python3 cm_ref_parity.py wc1 wc2` | WC-1 `cm_ref` == WC-2 output shape / a re-derived or fifth-schema `cm_ref` → reject | cm_ref-parity |
| gate invariants hold on the worked instance | enforced | `validate.py` (f) | as wave-level (f) | check-(f) |
| classes defined by canonical output telos (WC=artifact, PC=relation_graph, CC=judgment); closes L0-B1 | cognitive-review | independent β/CC read | β confirms the class↔telos rule is correct doctrine (a class-telos *table* is a locating aid, not the decision) | β judgment |
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
| state set closed & declared (`blocked` is a declared member, not only a transition target) | mechanically-verifiable | `python3 fsm_closure.py cell-fsm-state-table.json` | every transition target ∈ declared states / `blocked` reachable but undeclared → reject | closure report |
| total event/guard/action relation (every non-terminal state has owned exits) | mechanically-verifiable | `python3 fsm_totality.py cell-fsm-state-table.json` | every non-terminal has ≥1 exit / a non-terminal with no exit → reject | totality report |
| reachability proved (BFS from initial covers all states) | mechanically-verifiable | `python3 fsm_reach.py cell-fsm-state-table.json` | reachable set == declared states / an unreachable declared state → reject | reachability set |
| deterministic / normatively-prioritized transition selection (same bundle hash → same decision) | mechanically-verifiable | `python3 fsm_determinism.py cell-fsm-state-table.json` | no unresolved (state,event) ambiguity / two equal-priority exits on one (state,event) → reject | determinism report |
| CC-disposition→transition-request adapter (CC never mutates state) | mechanically-verifiable | `cue vet cc-transition-request.json <request-schema>` | CC emits a request the FSM validates / a CC transition mutating state directly → reject | adapter fixture |
| exactly one authority per edge (reconciles `transitions.json` / dispatch / resume / spec) | cognitive-review | independent β read | β confirms each edge has a single owning authority | β judgment |
| invalid-transition semantics defined | cognitive-review | independent β read | β confirms the invalid-transition rule is coherent | β judgment |
| review-return conflict resolved (§11.1 `→todo` vs §11.2 `→in-progress` settled to one edge) | cognitive-review | independent β read | β confirms one settled edge | β judgment |
| command-table parity (`cn cell` / `cn issues` map onto declared edges) | mechanically-verifiable | `python3 fsm_command_parity.py cell-fsm-state-table.json command-map.json` | each command maps to a declared edge / a command with no edge → reject | parity table |
| request-marker pattern pinned (typed marker + evidentiary guards) | cognitive-review | independent β read | β confirms the marker/guard pattern is carried | β judgment |
| emits machine-readable fixtures into receipt (cell-fsm-state-table, command-map) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## WC-3b — Wave FSM (`WAVE-FSM.md`) · fixtures `.cdd/unreleased/wc-3b/fixtures/`

Emitted machine-readable: `wave-fsm-state-table.json` and `disposition-to-transition-map.json`.

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| wave state set closed (enumerated states = declared set) | mechanically-verifiable | `python3 fsm_closure.py wave-fsm-state-table.json` | every target ∈ declared set / a target outside the set → reject | closure report |
| total event/guard/action relation (distinct state space from the cell FSM) | mechanically-verifiable | `python3 fsm_totality.py wave-fsm-state-table.json` | every non-terminal has owned exits / a non-terminal with no exit → reject | totality report |
| total mapping of all **8** CC dispositions → wave transition requests | mechanically-verifiable | `python3 disposition_totality.py disposition-to-transition-map.json` | all 8 dispositions mapped / a disposition with no transition → reject | mapping table |
| child-receipt aggregation rule (child completion predicates aggregate; CC does not mutate state) | mechanically-verifiable | `cue vet wave-transition-request.json <request-schema>` | aggregation request validates / a CC mutating wave state → reject | aggregation fixture |
| wave-transition-request shape pinned (mirrors the cell request-marker pattern + guards) | mechanically-verifiable | `cue vet wave-transition-request.json <request-schema>` | request with `{from,to,guard}` validates / a request missing a guard → reject | request fixture |
| child & whole-wave completion predicates **derived** (not a contract field); acyclic | enforced | `validate.py` (g) | as wave-level (g) / a per-contract `completion_signal` field or recursive completion → exit 1 | check-(g) |
| holding/replanning exits & terminal semantics defined | cognitive-review | independent β read | β confirms exits and terminals coherent | β judgment |
| invalid-transition semantics defined | cognitive-review | independent β read | β confirms the invalid-transition rule | β judgment |
| dependency & gate dispatch rules (wave-boundary auth once; children scheduled by deps) | cognitive-review | independent β read | β confirms §9 wave-boundary policy carried | β judgment |
| separate authority from the cell FSM (two independently total machines) | cognitive-review | independent β read | β confirms two distinct total machines, not one conflated | β judgment |
| emits machine-readable fixtures into receipt (wave-fsm-state-table, disposition map) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## WC-4 — Shipped→specified migration (`CONTRACT-MIGRATION.md`) · fixtures `.cdd/unreleased/wc-4/fixtures/`

Emitted machine-readable: `field-map.json` (`{shipped_field → {target, status ∈ mapped|added|deprecated}}`),
`round-trip/*.json`, and five `negative/*.json`.

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| field-by-field preservation (every shipped field mapped/added/deprecated) | mechanically-verifiable | `python3 field_coverage.py field-map.json schemas/cdd/contract.cue` | every shipped field present with a status / an unmapped shipped field → reject | coverage report |
| round-trip oracle (lossless paths proven lossless; lossy points named) | mechanically-verifiable | `python3 round_trip.py round-trip/*.json` | declared-lossless path round-trips identically / a "lossless" path that loses a field → reject | round-trip log |
| negative fixtures (reject unknown field / dropped required / mutable-root input / identity-only revision / missing adapter case) | mechanically-verifiable | `python3 run_negatives.py negative/*.json` | all 5 rejected / any of the 5 accepted → the oracle fails | 5 rejection receipts |
| CM edge representation consumed from WC-2 (target `cm_ref` == WC-2 shape; within four schemas) | mechanically-verifiable | `python3 cm_ref_parity.py wc4 wc2` | WC-4 target `cm_ref` == WC-2 shape / a re-derived or fifth-schema `cm_ref` → reject | cm_ref-parity |
| versioned projection specified (`cnos.cdd.contract.v1 → cn.cell.contract.v1`, named version + direction) | cognitive-review | independent β read | β confirms the projection is named + directed (a header string is a locating aid, not the decision) | β judgment |
| State-A/State-B field enumeration (shipped commands/fields vs target fields) | cognitive-review | independent β read | β confirms the enumeration is complete/honest | β judgment |
| migration is doctrine/schema-boundary only (no CUE/Go adapter implemented here) | cognitive-review | independent β read (also `validate.py` allowed/forbidden paths bound `schemas/**`, `src/**`) | β confirms no adapter code authored | β judgment |
| emits machine-readable fixtures into receipt (field-map, round-trip, 5 negatives) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

## WC-5 — Integration seal (`CELL-RUNTIME-INTEGRATION.md`) · fixtures `.cdd/unreleased/wc-5/fixtures/`

Emitted machine-readable: `reconcile-627-audit.json`, `fail-disposition-audit.json`,
`vocab-consistency.json`, `schema-fixture-log.json`.

| Predicate | Kind | Oracle | Expected positive / negative | Evidence |
|---|---|---|---|---|
| graph parity (derived == authored edges) | enforced | `validate.py` (c) | as wave-level (c) / a spurious/missing edge → exit 1 | check-(c) |
| schema fixtures pass (every worked instance validates against §2) | mechanically-verifiable | `python3 schema_fixture_run.py schema-fixture-log.json` | every worked instance validates / a non-conformant instance → reject | fixture pass log |
| reconcile-627 mapping complete (every S-node classified exactly once, no double owner) | mechanically-verifiable | `python3 reconcile_audit.py reconcile-627-audit.json` | each S-node classified once / a doubly-owned or omitted S-node → reject | #627 map audit |
| residual debt disposed exactly once (each grounding FAIL closed-by / retired / tracked once) | mechanically-verifiable | `python3 fail_disposition_audit.py fail-disposition-audit.json` | each FAIL disposed once / a FAIL disposed twice or not at all → reject | FAIL-disposition audit |
| migration round-trips (WC-4 round-trip oracle green on named fixtures) | mechanically-verifiable | `python3 round_trip.py ../wc-4/fixtures/round-trip/*.json` | WC-4 round-trips green / a failing round-trip → reject | round-trip log |
| wc5 readiness over non-recursive N (receipt exists for each of wc-1,wc-2,wc-3a,wc-3b,wc-4) | enforced + evidenced | `validate.py` (g) N-set; receipts | N derived correctly (g) / receipts for all N bound at seal time | readiness receipt |
| wc5 own completion non-recursive (own output+acceptance+V+receipt; no self/wave quantification) | enforced + evidenced | `validate.py` (g); receipt | (g) proves non-recursion / a wc-5 completion reading `wave_complete` → (g) exit 1 | own-completion receipt |
| cross-artifact vocabulary consistent (one `cm_ref` vocab, one CM object, one κ≠α logic across WC-1..WC-4) | cognitive-review | independent β/CC read | β confirms one vocabulary, one κ≠α logic (semantic consistency, not grep) | β judgment |
| final integration V fails-closed on a conflicting `cm_ref`/FSM-state/schema-field/missing-adapter/unmapped-627-node | cognitive-review | independent β/CC read (the mechanical parts route to the audits above) | β confirms the seal fails closed on each conflict class | β judgment |
| emits machine-readable fixtures into receipt (reconcile audit, FAIL audit, vocab, schema-fixture log) | evidenced | receipt binds the emitted fixtures | receipt lists + hashes each fixture / missing → V FAIL | receipt evidence |

---

*Rows marked **cognitive-review** are honestly not mechanical: an independent β (outside the Sigma
lineage) or a CC doctrine judgment decides them. They are surfaced here so the reviewer knows exactly
which predicates carry warrant by judgment vs. by machine — the R2 error of implementing a
semantic-absence claim as grep-and-calling-it-mechanical is not repeated.*
