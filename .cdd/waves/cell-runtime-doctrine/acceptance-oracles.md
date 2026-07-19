# Acceptance oracles — decidable checks for every load-bearing predicate (cnos#671 R2)

Repairs external-β finding 4: every load-bearing `acceptance.predicate` (per contract) and every
wave-level predicate names a **decidable oracle** — a fixture/schema path, a command, the expected
rejection cases, and the evidence output it produces. An oracle is *decidable* when a credential-free
mechanical run yields PASS/FAIL deterministically. The Planning-Cell **pre-authorization** subset
(shape, refs, parity, DAG, surfaces, gates, completion-acyclicity) is fully automated **now** by
`validate.py` (checks a–g); it runs at the wave tree and exits non-zero on any violation. The
**per-artifact** oracles below are the checks each Working Cell's own coherence loop must pass when it
authors its `docs/` artifact (they run against artifacts that do not yet exist — they are the
mechanical acceptance contract the WC inherits).

## Legend
- **Oracle**: `validate.py:<check>` (automated now) · `schema-fixture` (CUE/JSON-schema validate) ·
  `text-assertion` (grep/parse over the produced artifact) · `round-trip` (encode∘decode identity).
- **Evidence**: the artifact a passing run emits into the cell's receipt.

## Wave-level (checked now by `validate.py` at the wave tree)

| Predicate (wave) | Oracle | Command | Expected rejection | Evidence |
|---|---|---|---|---|
| §2 key-path conformance, no extra top-level keys | `validate.py:a` | `python3 validate.py` | add `consumers`/`completion_signal` or drop a §2 key → exit 1 | check-(a) PASS line |
| immutable refs resolve (intent, grounding **source 9d1ab3a5…**, pins) | `validate.py:b` | same | corrupt snapshot / flip a `contract_sha256` / bare comment-id revision → exit 1 | check-(b) PASS line |
| edge parity (authored == sibling-output-derived) | `validate.py:c` | same | add/remove an edge or change an `output_id` → exit 1 | derived-vs-authored edge sets equal |
| dependency graph is a DAG | `validate.py:d` | same | introduce a cycle in `sibling_output` refs → exit 1 | topo-sort covers all nodes |
| parallel nodes share no write surface | `validate.py:e` | same | two nodes share `requested_output.path`/`allowed_paths` → exit 1 | per-pair disjointness |
| gate invariants | `validate.py:f` | same | `doctrine_affecting:true` + `operator_acceptance_required:false`, or missing `reason` → exit 1 | check-(f) PASS line |
| completion-predicate graph acyclic (non-recursive) | `validate.py:g` | same | put `wc-5` in `construction_set_N`, or make `wc5_completion` read `wave_complete` → exit 1 | predicate-DAG topo + text anti-recursion |
| completion fixtures behave | `text-assertion` on `wave.completion.fixtures` | evaluate the 4 named fixtures | `all-N-complete/WC-5-absent` yielding `wave_complete=true`, or `WC-5-FAIL` yielding `true` → reject | fixture truth-table receipt |

## WC-2 — Coherence-Measurement contract (`COHERENCE-MEASUREMENT.md`)

| Predicate | Oracle | Command / check | Expected rejection | Evidence |
|---|---|---|---|---|
| CM is a runner/provider-produced object (producer = tsc/runtime op, **not** CC) | `text-assertion` | assert `producer` field names a CM provider; assert no clause equates CM with a CC judgment | a CM whose `producer` is a Cohering Cell | producer-typed CM fixture |
| CC result is a separate tagged object consuming immutable CM refs | `schema-fixture` | validate a CM object and a CC-judgment object as **distinct** tagged types | a single object tagged both CM and judgment | two-fixture validate log |
| `receipt_core → CM → V → δ → final_receipt` type path declared | `schema-fixture` | a `receipt_core` (no `validation`/`boundary_decision`) validates; a `final_receipt` satisfies shipped `receipt.cue` | a receipt claiming final status without V/δ outputs | round-trip receipt fixtures |
| CM→V edge pinned (V consumes an immutable `cm_ref`) | `schema-fixture` | V input carries a named immutable `cm_ref` | a V invocation with no `cm_ref` | V-signature fixture |
| D9 four-schema boundary settled; CM realized **within** four schemas | `text-assertion` + `validate.py`(reject-fifth note) | assert no fifth canonical `cn.cm.v1`; `cm_ref` resolves within the four schemas | **negative fixture:** a proposed fifth canonical schema → reject | four-schema conformance receipt |
| reproducible (same frozen inputs + provider version → same CM) | `round-trip` | re-run `coh` on the frozen bundle digest | a differing CΣ under identical frozen inputs | reproduction diff = ∅ |

## WC-1 — Cell classes & typed contract (`CELL-RUNTIME-CLASSES.md`)

| Predicate | Oracle | Command / check | Expected rejection | Evidence |
|---|---|---|---|---|
| classes defined by canonical output telos (closes L0-B1) | `text-assertion` | WC=artifact, PC=relation_graph, CC=judgment | a PC whose canonical output is an artifact | class↔telos table |
| `cn.cell.contract.v1` envelope = §2 key paths | `schema-fixture` | template + worked instance normalize to identical key-path sets | any instance with an extra/missing top-level key | key-path diff = ∅ |
| input reference model is the provenance-tagged union | `schema-fixture` | each `inputs.required[]` is `external{locator}` or `sibling_output{producer,output_id}` | an untyped/opaque optional scalar | union-validate log |
| `cm_ref` imported verbatim from WC-2, within four schemas, no fifth | `text-assertion` | WC-1 `cm_ref` == WC-2 output shape | a re-derived or fifth-schema `cm_ref` | cm_ref-parity check |
| gate invariants hold | `validate.py:f` | as above | — | check-(f) |
| key-path parity (template ≡ worked instance) | `schema-fixture` | parse both, diff normalized key paths | non-identical key-path sets | diff = ∅ |

## WC-3a — Cell FSM (`CELL-FSM.md`)

| Predicate | Oracle | Command / check | Expected rejection | Evidence |
|---|---|---|---|---|
| state set closed & declared (`blocked` is a member) | `text-assertion` | every transition target ∈ declared state set | `blocked` reachable but undeclared | state-set closure report |
| exactly one authority per edge | `text-assertion` | each edge has a single owner among table/dispatch/resume/spec | an edge owned by two authorities (the §11.1 vs §11.2 conflict) | edge-authority table |
| total event/guard/action relation | `text-assertion` | every non-terminal state has owned exits | a non-terminal state with no exit | totality report |
| CC-disposition→transition-request adapter (CC never mutates state) | `schema-fixture` | CC emits a request the FSM validates | a CC transition that mutates state directly | adapter fixture |
| reachability proved | `text-assertion` | BFS from initial covers all states | an unreachable declared state | reachability set |

## WC-3b — Wave FSM (`WAVE-FSM.md`)

| Predicate | Oracle | Command / check | Expected rejection | Evidence |
|---|---|---|---|---|
| wave state set closed | `text-assertion` | declared set matches the enumerated states | a target outside the set | state-set report |
| total mapping of all 8 CC dispositions → wave transition requests | `text-assertion` | all 8 dispositions mapped | a disposition with no transition | disposition→transition table |
| child-receipt aggregation rule (no CC state mutation) | `schema-fixture` | child completion predicates aggregate into wave state | CC mutating wave state | aggregation fixture |
| child & whole-wave completion predicates **derived**, not a contract field | `text-assertion` + `validate.py:g` | predicates derived from child acceptance/V/receipt; acyclic | a per-contract `completion_signal` field; recursive completion | check-(g) + derivation note |
| separate authority from the cell FSM | `text-assertion` | two independently total machines | one machine conflated across scopes | two-machine proof |

## WC-4 — Shipped→specified migration (`CONTRACT-MIGRATION.md`)

| Predicate | Oracle | Command / check | Expected rejection | Evidence |
|---|---|---|---|---|
| versioned projection specified | `text-assertion` | named version + direction `cnos.cdd.contract.v1 → cn.cell.contract.v1` | an unversioned/undirected mapping | projection header |
| field-by-field preservation | `text-assertion` | every shipped field mapped/added/deprecated | an unmapped shipped field | preservation matrix |
| round-trip oracle | `round-trip` | lossless paths proven lossless; lossy points named | a "lossless" path that loses a field | round-trip green log |
| negative fixtures | `schema-fixture` | reject unknown field / dropped required / mutable-root input / identity-only revision / missing adapter case | any of those accepted | 5 rejection receipts |

## WC-5 — Integration seal (`CELL-RUNTIME-INTEGRATION.md`)

| Predicate | Oracle | Command / check | Expected rejection | Evidence |
|---|---|---|---|---|
| cross-artifact vocabulary consistent | `text-assertion` | one `cm_ref` vocab, one CM object, one κ≠α logic across WC-1..WC-4 | two divergent `cm_ref` shapes | vocab-consistency report |
| schema fixtures pass | `schema-fixture` | every worked instance validates against §2 | a non-conformant instance | fixture pass log |
| graph parity | `validate.py:c` | derived == authored edges | a spurious/missing edge | check-(c) |
| migration round-trips | `round-trip` | WC-4 round-trip oracle green on named fixtures | a failing round-trip | round-trip log |
| reconcile-627 mapping complete | `text-assertion` | every S-node classified exactly once, no double owner | a doubly-owned or omitted S-node | #627 map audit |
| residual debt disposed exactly once | `text-assertion` | each grounding FAIL closed-by / retired / tracked exactly once | a FAIL disposed twice or not at all | FAIL-disposition audit |
| **wc5 readiness over N (non-recursive)** | `validate.py:g` + receipts | completion receipt exists for each of wc-1,wc-2,wc-3a,wc-3b,wc-4 | evaluating readiness over a set containing wc-5 | readiness receipt |
| **wc5 own completion (non-recursive)** | receipts | wc-5 output+acceptance+V+receipt; does not quantify over itself/wave predicate | a wc-5 completion that reads `wave_complete` | own-completion receipt |
