<!-- wave-revision: R9 -->
# Regression fixtures — each blocker / §2-drift mutation is rejected by `cue vet`

`wave.clean.yaml` and `contract.clean.yaml` are minimal-but-complete instances that PASS; the four
`contract.canonical-*.yaml` are canonical §2 variants that PASS `#CellContract` (nullable scope,
optional `prior_receipt`, `relation_graph`, `judgment`). Each `*.bad-*.yaml` applies exactly ONE
mutation to a clean base and MUST be rejected. Run them all with `make -C .. regressions` (fails the
build if any bad fixture passes). Individually:

```
cue vet ../ ./wave.bad-<name>.yaml     -d '#Wave'          # expect non-zero
cue vet ../ ./contract.bad-<name>.yaml -d '#CellContract'  # expect non-zero
```

## Wave fixtures (`-d '#Wave'`)

| fixture | mutation | cue vet rejection |
|---|---|---|
| `wave.bad-schema-const.yaml` | `schema: nonsense` | `schema: conflicting values "cn.wave.v1" and "nonsense"` |
| `wave.bad-duplicate-schema-key.yaml` | duplicate `schema:` key | `schema: conflicting values "nonsense" and "cn.wave.v1"` |
| `wave.bad-receipt-bound-string.yaml` | `receipt_bound: "false"` (string) | `receipt_bound: conflicting values "false" and bool` |
| `wave.bad-missing-evidence-bound.yaml` | record drops `evidence_bound` | `...evidence_bound: field is required but not present` |
| `wave.bad-empty-external-roots.yaml` | `external_roots: []` | `external_roots: incompatible list lengths (0 and 1)` |
| `wave.bad-auth-missing-revision-bound.yaml` | drop `wave_authorization.revision_bound` | `...revision_bound: field is required but not present` |
| `wave.bad-empty-contract-ref-resolution.yaml` | `contract_ref_resolution: {}` | `...authorized_wave_commit: field is required but not present` |
| `wave.bad-edge-derivation-no-parity.yaml` | drop `parity_required` | `edge_derivation.parity_required: field is required but not present` |
| `wave.bad-edge-derivation-acyclic-false.yaml` | `acyclic: false` | `edge_derivation.acyclic: conflicting values true and false` |
| `wave.bad-wave-complete-absent.yaml` | delete the `wave_complete` const | `completion.wave_complete: field is required but not present` |
| `wave.bad-wave-complete-mutated.yaml` | change the `wave_complete` const | `completion.wave_complete: conflicting values ...` |
| `wave.bad-duplicate-node-id.yaml` | two nodes with id `wc-2` | `_unique_node_ids."wc-2": conflicting values 1 and 0` |
| `wave.bad-duplicate-edge.yaml` | duplicate `wc-2 -> wc-1` edge | `_unique_edges."wc-2=>wc-1": conflicting values 1 and 0` |
| `wave.bad-stop-effect-enum.yaml` | STOP `effect: nonsense` | `stop_conditions.0.effect: ... empty disjunction` |
| `wave.bad-empty-stop-conditions.yaml` | `stop_conditions: []` | `stop_conditions: incompatible list lengths (0 and 1)` |

## Contract fixtures (`-d '#CellContract'`)

| fixture | mutation | cue vet rejection |
|---|---|---|
| `contract.bad-unknown-top-key.yaml` | add top-level `oracles:` (the R4 regression) | `oracles: field not allowed` |
| `contract.bad-cell-class.yaml` | `cell.class: nonsense` | `cell.class: ... empty disjunction` |
| `contract.bad-requested-output-kind.yaml` | `requested_output.kind: nonsense` (not in `artifact\|relation_graph\|judgment`) | `requested_output.kind: 3 errors in empty disjunction` |
| `contract.bad-gate-invariant.yaml` | `doctrine_affecting: true` + `operator_acceptance_required: false` | `gates.operator_acceptance_required: conflicting values true and false` |
| `contract.bad-doctrine-affecting-string.yaml` | `doctrine_affecting: "true"` (string) | `doctrine_affecting: conflicting values "true" and bool` |
| `contract.bad-empty-required.yaml` | `inputs.required: []` | `inputs.required: incompatible list lengths (0 and 1)` |
| `contract.bad-optional-sibling.yaml` | a `sibling_output` in `inputs.optional` | `inputs.optional.0.ref_kind: conflicting values "external" and "sibling_output"` |
| `contract.bad-empty-acceptance.yaml` | `acceptance.predicates: []` | `acceptance.predicates: incompatible list lengths (0 and 1)` |
| `contract.bad-empty-allowed-paths.yaml` | `constraints.allowed_paths: []` | `constraints.allowed_paths: incompatible list lengths (0 and 1)` |
| `contract.bad-missing-forbidden-paths.yaml` | drop `constraints.forbidden_paths` | `constraints.forbidden_paths: field is required but not present` |
| `contract.bad-missing-non-goals.yaml` | drop `constraints.non_goals` | `constraints.non_goals: field is required but not present` |
| `contract.bad-gate-reason-nonnull-both-false.yaml` | both gates false + non-null `reason` | `gates.reason: conflicting values "..." and null` |

All 27 rejections (15 wave + 12 contract, incl. the 7 §2-drift negatives) are verified green by
`make regressions` (cue v0.17.1). CI / the WC Go tooling run this suite as the structural gate.
