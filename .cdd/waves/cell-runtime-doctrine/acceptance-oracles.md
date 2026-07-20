<!-- wave-revision: R10 -->
# Acceptance oracles — honestly classified, registry-projected (cnos#671 R10)

Every load-bearing acceptance predicate (wave + WC-1..WC-5) is classified as **exactly one** of five
honest kinds. Nothing cognitive is dressed up as mechanical; a semantic-absence claim is **never**
implemented as grep-and-called-mechanical.

**This document is a PROJECTION.** The authoritative machine-readable source is the **TOTAL**
[`oracle-registry.yaml`](./oracle-registry.yaml), whose shape is validated now by
`cue vet ./schema/ ./oracle-registry.yaml -d '#AssuranceRegistry'` (classification enum + per-category
required fields). Each child contract consumes its slice via a canonical `inputs.required[].external`
control_plane ref (no non-canonical contract key). The two tables below are generated from the registry.

## R8 re-architecture (operator directive: Go + CUE repo, NO Python)

The prior rounds enforced the wave-level structural predicates with a hand-rolled Python validator
(`validate.py`) — **removed entirely** in R8. Structural pre-authorization validation is now `cue vet`
against the plan-local closed schemas under [`schema/`](./schema/); procedural/semantic checks are
**named-and-deferred to Go**. There is **no `.py`** anywhere: every `mechanically-verifiable` oracle is
a **Go checker or a CUE schema** the child WC emits.

## Classification kinds

- **structural-cue** — enforced **NOW** by `cue vet` against a plan-local `#Def`
  (`#CellContract`/`#Wave`/`#AssuranceRegistry`/`#Intent`). Credential-free; exits non-zero on any
  closed-struct / typed-field / enum / cardinality / required-field / pinned-const violation. **Not** a
  blanket soundness certificate — it warrants exactly the declared constraints, nothing beyond.
- **deferred-go** — a **named Go validator** with **exactly one** in-wave `deferred_owner` (procedural/
  semantic: DAG acyclicity, edge parity, ref/hash resolution, classification/oracle-ownership bijection,
  write-surface disjointness, ledger consistency, completion-evidence derivation). #627 is never the
  owner (S2–S3 are downstream consumers/canonicalizers). **Not** run at this tree; run when the owning
  WC executes. Named, not implemented (D9, option B).
- **mechanically-verifiable** — a per-child acceptance oracle the child WC **emits** as a **Go checker**
  (`go run .cdd/unreleased/wc-N/fixtures/CHECKER.go FIXTURE.json`) **or a CUE schema**
  (`cue vet FIXTURE.json .cdd/unreleased/wc-N/fixtures/NAME.schema.cue`), with named positive/negative
  fixtures at concrete repo-root-relative paths, bound (bytes + results) into the receipt. **No `.py`.**
  The checkers/schemas/fixtures do not exist yet — they are the mechanical acceptance contract each
  Working Cell inherits and must emit under `.cdd/unreleased/wc-N/fixtures/`.
- **evidenced** — verified by **bound evidence in the receipt**, not self-report.
- **cognitive-review** — an **independent β / CC doctrine judgment**, honestly **NOT mechanical**.

**Summary counts (78 child acceptance predicates, one classification each):** structural-cue **1** ·
deferred-go **19** · mechanically-verifiable **30** · evidenced **7** · cognitive-review **21**. The
single-owner deferred-validator **gating predicates** (9) are each a `deferred-go` entry with exactly one
`deferred_owner`. **R10 — forward-only:** every child deferred acceptance entry is verified by the child
**itself** or a construction-**predecessor** (`deferred_owner ∈ {owner} ∪ predecessors(owner)`, CUE-enforced),
so the combined assurance graph is acyclic; the whole-wave oracle-ownership bijection is a **wave-boundary**
predicate (owned by no child) held in `wave_predicates:`.

## Wave-level structural / deferred predicates

| Predicate | Kind | Oracle / owner |
|---|---|---|
| §2 contract constraint model (closed key set + enums + types) | structural-cue | `cue vet ./schema/ ./contracts/<wc>…yaml -d '#CellContract'` |
| full `cn.wave.v1` envelope (schema const, node/edge uniqueness, nonempty external_roots/STOP, required authority+contract_ref_resolution, pinned edge_derivation consts, pinned completion formula + closed 5-bool records, duplicate-key rejection) | structural-cue | `cue vet ./schema/ ./wave.cn-wave-v1.yaml -d '#Wave'` |
| registry shape + classification enum + per-category required fields | structural-cue | `cue vet ./schema/ ./oracle-registry.yaml -d '#AssuranceRegistry'` |
| intent projection shape | structural-cue | `cue vet ./schema/ ./intent.cn-intent-v1.yaml -d '#Intent'` |
| gate invariant (doctrine_affecting ⇒ operator_acceptance_required; reason required) | structural-cue | native `#CellContract` conditional |
| §2-drift + wave-envelope + backward-edge regression coverage (29 bad fixtures each rejected) | structural-cue | `make -C schema regressions` |
| immutable ref / content-hash resolution (intent id/schema; `commit:path`; revision hashes; grounding `9d1ab3a5…`) | deferred-go | **WC-2** — `.cdd/unreleased/wc-2/validators/ref_resolve.go` |
| edge parity (authored == sibling-output-derived) | deferred-go | **WC-3b** — `.cdd/unreleased/wc-3b/validators/edge_parity.go` |
| dependency graph is a DAG | deferred-go | **WC-3b** — `.cdd/unreleased/wc-3b/validators/wave_dag.go` |
| parallel nodes share no write surface | deferred-go | **WC-3b** — `.cdd/unreleased/wc-3b/validators/write_surface.go` |
| combined-graph acyclicity (construction ∪ cross-owner assurance edges; forward-only) — R10 | deferred-go | **WC-3b** — `.cdd/unreleased/wc-3b/validators/combined_graph_acyclic.go` |
| completion-evidence derivation (typed resolver input → 5 derived constituents; predicate-graph acyclicity; fixture `expected`) | deferred-go | **WC-5** — `.cdd/unreleased/wc-5/validators/completion_evidence.go` |
| classification totality bijection (`union(acceptance.predicates)` ⇄ registry, each once) | deferred-go | **WC-5** — `.cdd/unreleased/wc-5/validators/classification_bijection.go` |
| ledger consistency (revision markers agree; per-category counts) | deferred-go | **WC-5** — `.cdd/unreleased/wc-5/validators/ledger_consistency.go` |
| oracle-ownership bijection (registry ⇄ ALL contracts' mechanically-verifiable predicates; whole-wave) — R10: **wave-boundary** | deferred-go | **WAVE** — `.cdd/waves/cell-runtime-doctrine/wave-validators/oracle_ownership_bijection.go` |

---

## Complete assurance classification — projection of the total registry

One row per **every** child acceptance predicate (all six contracts, 78 rows), each classified into
**exactly one** category. Generated from `oracle-registry.yaml`; the `union(acceptance.predicates)` ⇄
registry bijection that proves this projection total + singular is a **deferred-go** check (owned by WC-5).
The 9 R9 single-owner deferred-validator gating predicates are listed in the dedicated block at the end.

| owner | predicate | classification |
|---|---|---|
| wc-1 | `wc_pc_cc_classes_defined_by_canonical_output_telos__wc_artifact_pc_relation_graph_cc_judgment__closes_L0_B1` | cognitive-review |
| wc-1 | `cn_cell_contract_v1_envelope_specified__cell_id_class_mode_protocol_matter_domain__inputs_required_provenance_tagged_union__requested_output_id_kind_path__acceptance__constraints__gates__doctrine_affecting__stop_conditions` | cognitive-review |
| wc-1 | `input_reference_model_is_the_provenance_tagged_union__external_locator_classed_or_sibling_output_producer_output_id` | mechanically-verifiable |
| wc-1 | `class_specific_v_predicates_reference_an_imported_cm_ref__wc1_does_not_inline_cm_internals` | cognitive-review |
| wc-1 | `cm_ref_field_shape_is_imported_verbatim_from_wc2_output_realized_within_the_four_schemas__no_re_derivation_no_fifth_schema` | mechanically-verifiable |
| wc-1 | `kappa_neq_alpha_role_logic_unconditional__state_a_is_hosting_identity_collapse_not_role_equality__carries_repaired_CC_1_as_settled_input` | cognitive-review |
| wc-1 | `review_evidence_content_bound_not_identity_bound__carries_repaired_CC_2_as_settled_input` | cognitive-review |
| wc-1 | `gate_invariants_hold__doctrine_affecting_true_implies_operator_acceptance_required_true__reason_present_iff_a_gate_bool_true` | structural-cue |
| wc-1 | `key_path_parity__template_and_worked_instance_normalize_to_identical_key_path_sets` | mechanically-verifiable |
| wc-1 | `scope_guard__no_cm_internals_no_cell_or_wave_fsm_no_migration_authored_here` | cognitive-review |
| wc-1 | `emits_machine_readable_acceptance_fixtures_into_receipt__contract_template_yaml__worked_instance_yaml__input_union_fixtures__each_validated_by_the_named_acceptance_oracle_command` | evidenced |
| wc-1 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | deferred-go |
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
| wc-2 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | deferred-go |
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
| wc-3a | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | deferred-go |
| wc-3b | `wave_fsm_state_set_closed__intent_planning_required_d0_review_wave_planning_wave_review_dispatchable_executing_holding_replanning_completing_complete_complete_with_residuals` | mechanically-verifiable |
| wc-3b | `total_event_guard_action_relation__distinct_from_the_cell_fsm_state_space_and_events` | mechanically-verifiable |
| wc-3b | `invalid_transition_semantics_defined` | cognitive-review |
| wc-3b | `child_receipt_aggregation_rule__child_completion_predicates_aggregate_into_wave_state_without_cc_mutating_state` | mechanically-verifiable |
| wc-3b | `total_mapping_all_eight_cc_dispositions_to_wave_transition_requests__request_planning_request_working_hold_request_human_continue_wave_complete_complete_with_residuals_block` | mechanically-verifiable |
| wc-3b | `dependency_and_gate_dispatch_rules__wave_boundary_authorization_once_then_children_scheduled_by_dependencies_no_per_child_operator_gate` | cognitive-review |
| wc-3b | `holding_and_replanning_exits_and_terminal_semantics_defined` | cognitive-review |
| wc-3b | `wave_transition_request_shape_pinned__mirrors_the_cell_fsm_request_marker_pattern_with_guards` | mechanically-verifiable |
| wc-3b | `separate_authority_from_the_cell_fsm__two_independently_total_machines_one_shared_transition_algebra_or_two_tagged_instances` | cognitive-review |
| wc-3b | `child_and_whole_wave_completion_predicates_defined__derived_from_child_acceptance_v_receipt_surfaces_not_a_new_contract_field` | deferred-go |
| wc-3b | `emits_machine_readable_acceptance_fixtures_into_receipt__wave_fsm_state_table_json__and_eight_disposition_to_transition_map_json__for_totality_and_mapping_oracles` | evidenced |
| wc-3b | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | deferred-go |
| wc-4 | `versioned_projection_specified__cnos_cdd_contract_v1_to_cn_cell_contract_v1_with_a_named_version_and_direction` | cognitive-review |
| wc-4 | `field_by_field_preservation_rules__every_shipped_field_mapped_or_explicitly_marked_added_or_deprecated` | mechanically-verifiable |
| wc-4 | `cm_edge_representation_consumed_from_wc2__the_target_contract_cm_ref_shape_is_wc2_output_realized_within_the_four_schemas_not_a_fifth_schema_not_re_derived` | mechanically-verifiable |
| wc-4 | `state_a_state_b_field_enumeration__shipped_commands_and_schema_fields_vs_target_fields_enumerated_per_field` | cognitive-review |
| wc-4 | `round_trip_oracle__lossless_paths_proven_lossless_and_lossy_points_named` | mechanically-verifiable |
| wc-4 | `negative_fixtures__reject_unknown_field_dropped_required_field_mutable_root_input_identity_only_revision_missing_adapter_case` | mechanically-verifiable |
| wc-4 | `migration_is_doctrine_schema_boundary_only__no_cue_or_go_adapter_implemented_here` | cognitive-review |
| wc-4 | `emits_machine_readable_acceptance_fixtures_into_receipt__field_map_json__round_trip_fixtures__five_negative_fixtures__for_field_preservation_and_round_trip_oracles` | evidenced |
| wc-4 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | deferred-go |
| wc-5 | `cross_artifact_vocabulary_consistent__one_canonical_input_ref_and_requested_output_vocabulary_one_cm_object_one_kappa_neq_alpha_logic_across_wc1_through_wc4` | cognitive-review |
| wc-5 | `schema_fixtures_pass__every_worked_instance_validates_against_the_selected_cn_cell_contract_v1_shape` | mechanically-verifiable |
| wc-5 | `graph_parity__wave_edges_equal_the_sibling_output_relations_exactly_no_missing_or_spurious_edge` | deferred-go |
| wc-5 | `migration_round_trips__wc4_round_trip_oracle_is_green_on_the_named_fixtures` | mechanically-verifiable |
| wc-5 | `reconcile_627_mapping_complete__every_S_node_classified_exactly_once_no_duplicate_owner_no_omitted_executable_child` | mechanically-verifiable |
| wc-5 | `residual_debt_disposed__every_grounding_fail_is_closed_by_or_retired_instance_defect_or_tracked_follow_up_exactly_once` | mechanically-verifiable |
| wc-5 | `final_integration_v_fails_closed_on__a_conflicting_cm_ref_a_conflicting_fsm_state_a_conflicting_schema_field_a_missing_adapter_case_or_an_unmapped_627_node` | cognitive-review |
| wc-5 | `wc5_readiness_over_non_recursive_construction_set_N__a_completion_receipt_exists_for_every_node_in_N_where_N_equals_wc1_wc2_wc3a_wc3b_wc4_and_excludes_wc5` | deferred-go |
| wc-5 | `wc5_own_completion_is_non_recursive__wc5_requested_output_produced_plus_wc5_acceptance_pass_plus_integration_v_pass_plus_bound_receipt__does_not_quantify_over_wc5_itself_nor_over_the_whole_wave_predicate` | deferred-go |
| wc-5 | `emits_machine_readable_acceptance_fixtures_into_receipt__reconcile_627_audit_json__fail_disposition_audit_json__vocab_consistency_json__schema_fixture_log_json__for_the_integration_seal_oracles` | evidenced |
| wc-5 | `mechanical_oracles_owned__every_mechanically_verifiable_predicate_binds_a_concrete_checker_or_schema_path_no_placeholder_operand_positive_fixture_exit_0_and_named_negative_fixture_exit_nonzero_within_allowed_paths_or_a_pinned_immutable_input_and_bound_in_the_receipt` | deferred-go |

---

## Registry projection (mechanically-verifiable) — parity-checked against `oracle-registry.yaml`

The 30 authoritative `mechanically-verifiable` entries. Each names a **Go checker** or a **CUE schema**
the child WC emits (no `.py`), with concrete repo-root-relative positive/negative fixtures. The
registry ⇄ projection parity is a **deferred-go** check.

| owner | predicate | classification | ownership | checker_or_schema | positive_fixture | negative_fixture |
|---|---|---|---|---|---|---|
| wc-1 | `input_reference_model_is_the_provenance_tagged_union__external_locator_classed_or_sibling_output_producer_output_id` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-1/fixtures/union_validate.go` | `.cdd/unreleased/wc-1/fixtures/input-union.positive.yaml` | `.cdd/unreleased/wc-1/fixtures/input-union-opaque-scalar.negative.yaml` |
| wc-1 | `cm_ref_field_shape_is_imported_verbatim_from_wc2_output_realized_within_the_four_schemas__no_re_derivation_no_fifth_schema` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-1/fixtures/cm_ref_parity.go` | `.cdd/unreleased/wc-1/fixtures/cm-ref-parity.positive.json` | `.cdd/unreleased/wc-1/fixtures/cm-ref-rederived.negative.json` |
| wc-1 | `key_path_parity__template_and_worked_instance_normalize_to_identical_key_path_sets` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-1/fixtures/keypath_diff.go` | `.cdd/unreleased/wc-1/fixtures/contract-worked-instance.yaml` | `.cdd/unreleased/wc-1/fixtures/contract-extra-key.negative.yaml` |
| wc-2 | `cm_defined_as_runner_or_provider_produced_object__producer_is_cm_provider_or_runtime_op_eg_tsc_not_cc` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/check_cm_producer.go` | `.cdd/unreleased/wc-2/fixtures/cm-object.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-producer-cc.negative.json` |
| wc-2 | `cm_object_fields_typed__s_alpha_s_beta_s_gamma_c_sigma_geom_mean_degeneracy_bottleneck_witnesses_defects_mode_target_bundle_digest_thresholds_standing_provenance` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/cm-object.schema.cue` | `.cdd/unreleased/wc-2/fixtures/cm-object.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-object.missing-bottleneck.negative.json` |
| wc-2 | `cc_class_result_is_a_separate_tagged_object_that_consumes_immutable_cm_refs__cm_and_cc_judgment_validate_as_distinct_tagged_objects` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/check_cm_cc_distinct.go` | `.cdd/unreleased/wc-2/fixtures/cm-cc-distinct.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-cc-conflated.negative.json` |
| wc-2 | `cm_to_v_edge_pinned__v_consumes_an_immutable_cm_ref_as_a_named_input_the_frozen_contract_times_receipt_signature_is_extended_to_carry` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/check_v_signature.go` | `.cdd/unreleased/wc-2/fixtures/v-signature.positive.json` | `.cdd/unreleased/wc-2/fixtures/v-signature-no-cmref.negative.json` |
| wc-2 | `receipt_core_to_cm_to_v_to_delta_to_final_receipt_type_path_declared__provisional_receipt_core_without_validation_or_boundary_decision_then_measurement_then_final_receipt_satisfying_shipped_receipt_cue` | mechanically-verifiable | existing | `schemas/cdd/receipt.cue` | `.cdd/unreleased/wc-2/fixtures/final-receipt.positive.json` | `.cdd/unreleased/wc-2/fixtures/final-receipt-no-v-delta.negative.json` |
| wc-2 | `cm_realized_within_four_schemas__a_typed_cm_field_or_edge_in_the_receipt_plus_a_cm_ref_resolving_within_the_existing_four_schemas__wc1_and_wc4_consume_this_same_cm_ref_shape` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/receipt-with-cmref.schema.cue` | `.cdd/unreleased/wc-2/fixtures/cm-ref-embedded.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-ref-needs-fifth.negative.json` |
| wc-2 | `negative_fixture_reject_fifth_canonical_schema__a_proposed_fifth_cn_cm_v1_canonical_schema_is_rejected_while_the_four_schema_boundary_remains_settled` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/reject_fifth.go` | `.cdd/unreleased/wc-2/fixtures/four-schema.positive.json` | `.cdd/unreleased/wc-2/fixtures/fifth-schema.negative.json` |
| wc-2 | `cm_ref_field_shape_is_the_single_interface_all_other_nodes_import__no_second_cm_shape_introduced_downstream` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-2/fixtures/cm_ref_parity.go` | `.cdd/unreleased/wc-2/fixtures/cm-ref-parity.positive.json` | `.cdd/unreleased/wc-2/fixtures/cm-ref-divergent.negative.json` |
| wc-3a | `cell_fsm_state_set_closed_and_declared__blocked_is_a_declared_member_not_only_a_transition_target` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_closure.go` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-blocked-undeclared.negative.json` |
| wc-3a | `total_event_guard_action_relation__every_nonterminal_state_has_owned_exits` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_totality.go` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-nonterminal-noexit.negative.json` |
| wc-3a | `deterministic_or_normatively_prioritized_transition_selection__same_bundle_hash_same_decision` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_determinism.go` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-ambiguous.negative.json` |
| wc-3a | `cc_disposition_to_cell_transition_request_adapter__cc_selects_a_judgment_the_fsm_validates_and_effects_the_transition__cc_never_mutates_state_directly` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/cell-transition-request.schema.cue` | `.cdd/unreleased/wc-3a/fixtures/cc-transition-request.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cc-mutates-state.negative.json` |
| wc-3a | `command_table_parity__cn_cell_and_cn_issues_surfaces_map_onto_declared_table_edges` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_command_parity.go` | `.cdd/unreleased/wc-3a/fixtures/command-map.positive.json` | `.cdd/unreleased/wc-3a/fixtures/command-map-orphan.negative.json` |
| wc-3a | `reachability_proved__every_declared_state_reachable_from_the_initial_state` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3a/fixtures/fsm_reach.go` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3a/fixtures/cell-fsm-unreachable.negative.json` |
| wc-3b | `wave_fsm_state_set_closed__intent_planning_required_d0_review_wave_planning_wave_review_dispatchable_executing_holding_replanning_completing_complete_complete_with_residuals` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/fsm_closure.go` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-target-outside.negative.json` |
| wc-3b | `total_event_guard_action_relation__distinct_from_the_cell_fsm_state_space_and_events` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/fsm_totality.go` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-state-table.positive.json` | `.cdd/unreleased/wc-3b/fixtures/wave-fsm-nonterminal-noexit.negative.json` |
| wc-3b | `child_receipt_aggregation_rule__child_completion_predicates_aggregate_into_wave_state_without_cc_mutating_state` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/wave-transition-request.schema.cue` | `.cdd/unreleased/wc-3b/fixtures/wave-aggregation-request.positive.json` | `.cdd/unreleased/wc-3b/fixtures/cc-mutates-wave-state.negative.json` |
| wc-3b | `total_mapping_all_eight_cc_dispositions_to_wave_transition_requests__request_planning_request_working_hold_request_human_continue_wave_complete_complete_with_residuals_block` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/disposition_totality.go` | `.cdd/unreleased/wc-3b/fixtures/disposition-to-transition-map.positive.json` | `.cdd/unreleased/wc-3b/fixtures/disposition-missing.negative.json` |
| wc-3b | `wave_transition_request_shape_pinned__mirrors_the_cell_fsm_request_marker_pattern_with_guards` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-3b/fixtures/wave-transition-request.schema.cue` | `.cdd/unreleased/wc-3b/fixtures/wave-transition-request.positive.json` | `.cdd/unreleased/wc-3b/fixtures/wave-request-missing-guard.negative.json` |
| wc-4 | `field_by_field_preservation_rules__every_shipped_field_mapped_or_explicitly_marked_added_or_deprecated` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/field_coverage.go` | `.cdd/unreleased/wc-4/fixtures/field-map.positive.json` | `.cdd/unreleased/wc-4/fixtures/field-map-unmapped.negative.json` |
| wc-4 | `cm_edge_representation_consumed_from_wc2__the_target_contract_cm_ref_shape_is_wc2_output_realized_within_the_four_schemas_not_a_fifth_schema_not_re_derived` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/cm_ref_parity.go` | `.cdd/unreleased/wc-4/fixtures/cm-ref-parity.positive.json` | `.cdd/unreleased/wc-4/fixtures/cm-ref-rederived.negative.json` |
| wc-4 | `round_trip_oracle__lossless_paths_proven_lossless_and_lossy_points_named` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/round_trip.go` | `.cdd/unreleased/wc-4/fixtures/round-trip.positive.json` | `.cdd/unreleased/wc-4/fixtures/round-trip-lossy.negative.json` |
| wc-4 | `negative_fixtures__reject_unknown_field_dropped_required_field_mutable_root_input_identity_only_revision_missing_adapter_case` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-4/fixtures/run_negatives.go` | `.cdd/unreleased/wc-4/fixtures/negatives-manifest.positive.json` | `.cdd/unreleased/wc-4/fixtures/negatives-regression.negative.json` |
| wc-5 | `schema_fixtures_pass__every_worked_instance_validates_against_the_selected_cn_cell_contract_v1_shape` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/schema_fixture_run.go` | `.cdd/unreleased/wc-5/fixtures/schema-fixture-log.positive.json` | `.cdd/unreleased/wc-5/fixtures/schema-fixture-nonconformant.negative.json` |
| wc-5 | `migration_round_trips__wc4_round_trip_oracle_is_green_on_the_named_fixtures` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/round_trip.go` | `.cdd/unreleased/wc-5/fixtures/migration-round-trip.positive.json` | `.cdd/unreleased/wc-5/fixtures/migration-round-trip-lossy.negative.json` |
| wc-5 | `reconcile_627_mapping_complete__every_S_node_classified_exactly_once_no_duplicate_owner_no_omitted_executable_child` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/reconcile_audit.go` | `.cdd/unreleased/wc-5/fixtures/reconcile-627-audit.positive.json` | `.cdd/unreleased/wc-5/fixtures/reconcile-627-double-owner.negative.json` |
| wc-5 | `residual_debt_disposed__every_grounding_fail_is_closed_by_or_retired_instance_defect_or_tracked_follow_up_exactly_once` | mechanically-verifiable | emitted | `.cdd/unreleased/wc-5/fixtures/fail_disposition_audit.go` | `.cdd/unreleased/wc-5/fixtures/fail-disposition-audit.positive.json` | `.cdd/unreleased/wc-5/fixtures/fail-disposed-twice.negative.json` |

### R10 single-CHILD-owner deferred-validator gating predicates (9)

Each is a child acceptance predicate whose PASS gates the owner's completion; classified `deferred-go`
with exactly one `deferred_owner` — and each is a **forward** self-edge (`deferred_owner == owner`).
WC-5 depends (sibling_output) on wc-2/wc-3b, so it cannot seal until every upstream validator passes.
**R10:** the whole-wave oracle-ownership bijection is NO LONGER here (it moved to the wave boundary,
`deferred_owner: "wave"`); WC-3b gained the combined-graph acyclicity gating predicate (net count 9).

| owner | predicate | classification | deferred_owner |
|---|---|---|---|
| wc-2 | `deferred_validator_owned__go_ref_content_hash_resolver_emitted_by_wc2_positive_and_named_negative_fixtures_green_gates_wc2_completion` | deferred-go | wc-2 |
| wc-3b | `deferred_validator_owned__go_wave_dag_validator_emitted_by_wc3b_positive_and_named_negative_fixtures_green_gates_wc3b_completion` | deferred-go | wc-3b |
| wc-3b | `deferred_validator_owned__go_edge_parity_validator_emitted_by_wc3b_positive_and_named_negative_fixtures_green_gates_wc3b_completion` | deferred-go | wc-3b |
| wc-3b | `deferred_validator_owned__go_write_surface_disjointness_validator_emitted_by_wc3b_positive_and_named_negative_fixtures_green_gates_wc3b_completion` | deferred-go | wc-3b |
| wc-3b | `deferred_validator_owned__go_combined_assurance_graph_acyclicity_validator_emitted_by_wc3b_positive_and_named_negative_fixtures_green_gates_wc3b_completion` | deferred-go | wc-3b |
| wc-5 | `deferred_validator_owned__go_completion_evidence_validator_emitted_by_wc5_typed_resolver_input_derives_the_five_constituents_positive_and_named_negative_fixtures_green_gates_wc5_completion` | deferred-go | wc-5 |
| wc-5 | `deferred_validator_owned__go_ledger_consistency_validator_emitted_by_wc5_positive_and_named_negative_fixtures_green_gates_wc5_completion` | deferred-go | wc-5 |
| wc-5 | `deferred_validator_owned__go_classification_totality_bijection_validator_emitted_by_wc5_positive_and_named_negative_fixtures_green_gates_wc5_completion` | deferred-go | wc-5 |
| wc-5 | `seal_readiness_consumes_every_upstream_deferred_validator_pass__wc2_ref_content_hash_resolution__wc3b_dag_edge_parity_write_surface_combined_graph_acyclicity__with_wave_boundary_oracle_ownership_bijection_preauthorized__each_bound_pass_is_a_precondition_of_wc5_readiness_via_child_complete` | deferred-go | wc-5 |

**Wave-boundary deferred-go predicate (owned by no child):**

| owner | predicate | classification | deferred_owner |
|---|---|---|---|
| WAVE | `wave_oracle_ownership_bijection_enforced` (whole-wave cross-contract bijection; pre-authorization) | deferred-go | wave |

---

*Rows marked **cognitive-review** are honestly not mechanical: an independent β (outside the Sigma
lineage) or a CC doctrine judgment decides them. Rows marked **deferred-go** are procedural/semantic
checks a **single-owner** Go validator runs when the owning WC executes — not at this tree, and never
in Python. #627 S2–S3 are downstream consumers/canonicalizers, never the owner.*
