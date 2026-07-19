---
name: measure/recursive-cell
description: >-
  A recursive coherence methodology for cnos cell architectures and cell
  executions. It measures the three-cell system, each cell's CCNF kernel,
  typed contract/evidence surfaces, cell and wave FSMs, and an actual receipt
  stream without allowing a strong local layer to average away a broken
  cross-level projection.
governing_question: >-
  Do the declared cell properties, their relations, their evolution, their
  internal kernels, their state machines, and the observed execution receipts
  still describe one executable coherence system?
artifact_class: measurement
scope: global
kata_surface: external
triggers:
  - cnos cell architecture measurement
  - WC PC CC ratification
  - recursive cell coherence
  - cell FSM coherence
inputs:
  - "Immutable target revision and ordered L0-L4 bundle"
  - "Pinned TSC Core/Operational and scoring instruction"
  - "Cell/wave contracts, schemas, FSM tables, implementation seams, and receipts"
  - "Control-plane evidence snapshotted or content-hashed"
outputs:
  - "Per-level TSC measurement reports for L0-L4"
  - "Cross-level projection and hard-invariant report"
  - "Overall measurement with exactly one bottleneck and provenance"
visibility: public
methodology:
  registry: "src/packages/cnos.cdd/skills/cdd/measure/recursive-cell/calibration/662/registry.tsc"
  calibration_preflight: "src/packages/cnos.cdd/skills/cdd/measure/recursive-cell/calibration/662/verify-target.sh"
  targets: ["cc662-system", "cc662-l0", "cc662-l1", "cc662-l2", "cc662-l3", "cc662-l4"]
  cross_target: true
  instruction: "src/packages/cnos.cdd/skills/cdd/measure/recursive-cell/INSTRUCTION.md"
  output_root: ".tsc/cnos-recursive-cell"
  default_mode: hybrid
  consistency:
    mechanical: identical
    llm_repeats: 3
    llm_spread: >-
      For each target, compute max absolute pairwise spread over alpha, beta,
      gamma, the three delta fields, and confidence; map max spread through
      phi(d)=d/(1-d), Coh_consistency=exp(-phi(d)). Report both max and
      mean-pairwise forms. A single semantic sample has no standing.
  standing:
    scope: house-authored-public-commons
    admissibility: public-only
    heldout_status: none
    external_anchor_count: 0
    llm_consistency_gate: reported-not-gating
    llm_consistency_floor: 0.90
  mechanical:
    backend: "usurobor/tsc@26aab5023f03dc7d0abf82e5fdba20134fc6adad:engine/ocaml/lib/mechanical_scoring.ml"
    determinism: >-
      Identical ordered file bundle, target manifest, TSC revision, and
      parameters must produce identical numeric scores and axis-detail
      signals over at least three runs.
    signals:
      alpha: ["alpha.terminology_consistency", "alpha.repeated_structure", "alpha.duplicate_definition_tension", "alpha.naming_drift"]
      beta: ["beta.cross_reference_consistency", "beta.authority_alignment", "beta.source_of_truth_alignment", "beta.target_file_fit"]
      gamma: ["gamma.canonical_generated_distinction", "gamma.version_surface_consistency", "gamma.traceability_presence", "gamma.authority_evolution_consistency"]
  llm:
    estimates: ["target", "alpha", "beta", "gamma", "delta_alpha_beta", "delta_beta_gamma", "delta_gamma_alpha", "bottleneck_axis", "confidence", "summary", "axis_evidence", "defect_cards", "unresolved_ambiguity", "next_fixes"]
    must_not:
      - "compute Coh or C_sigma; the engine owns barrier and aggregate computation"
      - "infer files outside the frozen target bundle"
      - "treat specified or illustrative behavior as shipped"
      - "average away a failed hard invariant or cross-level projection"
      - "use the CM being created or modified as its own independent warrant"
    validation: >-
      Validate the standard TSC v3.2.4 witness contract, the fixed checklist
      and defect-card correspondence, L0-L4 completeness, hard-invariant
      coverage, FSM totality coverage, and target/prompt digest equality.
    providers:
      local: "TSC external-response route using a frozen emitted prompt"
      ci: "Any pinned witness CLI restricted to the emitted prompt and one JSON response"
    ci_prompt: >-
      Read only the frozen prompt for the named L0-L4 target. Apply
      INSTRUCTION.md, emit exactly one standard TSC v3.2.4 JSON witness, and
      write no other artifact. Do not compute Coh or C-sigma.
---

# cnos Recursive Cell Coherence Methodology v0.1.0

## 0. Status and epistemic boundary

This is a **candidate CM derived from cnos#662**. That cycle is its calibration
source, not a held-out anchor. The CM may describe and reproduce the #662
measurement, but that result gives this CM no independent standing. Standing
requires at minimum:

1. three validated same-route semantic samples over a frozen prompt;
2. `Coh_consistency >= 0.90`;
3. a public calibration commons containing positive and adversarial cases;
4. at least one held-out cnos cell architecture not used to design this CM;
5. CM-of-CMs/admissibility measurement of this artifact.

The current TSC CLI can run the bundled #662 target manifests with its built-in
CM after `calibration/662/verify-target.sh` accepts the exact frozen receipt
checkout, but it cannot yet consume
this arbitrary CM declaration as a first-class instrument.
Until that executor exists, `INSTRUCTION.md` is the semantic authority and the
target manifests are the deterministic corpus projection. A new application
must materialize its own immutable registry by specializing those manifests;
the bundled registry is calibration evidence, not a wildcard production target.

## 1. Measurement object and recursion

The object is not one Markdown file. It is a recursive system observed at five
levels. Every level is articulated symmetrically as α/pattern, β/relation, and
γ/process. “Dominant telos” never exempts a level from either other axis.

| Level | System under measurement | α — properties/pattern | β — relations | γ — process/evolution |
|---|---|---|---|---|
| L0 | Three-cell agent | WC/PC/CC definitions and output types | handoffs, authority, next-MCA, ε scope lift | intent→plan→work→judgment loop and system evolution |
| L1 | One WC, PC, or CC | contract, α/β/γ/V/δ components, role boundaries | signatures, evidence flow, firebreaks, scope projection | production→review→closure→validation→decision and repair recursion |
| L2 | Typed surfaces | schemas, tags, locators, receipt forms, CM result | contract↔matter↔review↔measurement↔receipt↔V↔δ | versioning, binding, provisional/final boundaries, migrations |
| L3 | Runtime/FSM | state/event/guard/action vocabularies | transition authority, command/table parity, class/domain routing | complete cell lifecycle, wave lifecycle, recovery, replan, termination |
| L4 | Concrete execution | actual matter, receipts, review artifacts, claimed identities | exact-SHA and role/evidence links | round history, monotonic repair, mutability, exit behavior |

The recursion closes only if both conditions hold:

- **within-level coherence:** each Lk is coherent under α/β/γ;
- **cross-level coherence:** every authoritative Lk claim has a lossless,
  non-contradictory realization or explicitly owned migration at Lk+1.

Overall coherence is the geometric mean of per-level results only after every
hard invariant passes. A failed hard invariant forces a non-accepting
disposition regardless of a high mean.

## 2. L0 — three-cell system

### 2.1 Classifier

Class is determined by the **canonical requested output and its operational
consumer**, never by filename, prose topic, or matter domain:

- WC: canonical output is a repository/product artifact.
- PC: canonical output is executable relational structure—a D0 relation
  contract or wave graph whose consumers can derive work/dependencies.
- CC: canonical output is an evidence-backed state/process judgment with one
  disposition and, when required, a next-MCA.

A Markdown file may project any class for humans. If the Markdown itself is the
normative merged deliverable, it is an artifact and therefore WC. Mixed outputs
must split into cells or declare one canonical typed output and mark every other
surface as a derived projection.

### 2.2 Relations

Measure all edges, including:

```text
intent → PC-D0? → WC/PC-Wave → CC judgment → V → δ/FSM → next cell/human gate
closed cells/receipts → ε observation → CC at n+1
```

Every edge must name source type, target type, immutable identity, authority,
and failure behavior. A narrative arrow is not an executable relation.

### 2.3 Process

Check that system evolution cannot silently change class, domain, authority, or
measurement instrument. A CM must be selected and pinned before its result is
used as warrant. A CC may request creation or repair of a CM; it may not invent
the CM that retroactively justifies the same judgment.

## 3. L1 — internal cell kernel

For every WC/PC/CC instance, measure:

```text
matter := α(contract)
review := β(contract, matter)
measurement := CM(pinned bundle)
receipt := γ.bind(contract, matter, review, measurement, evidence)
verdict := V(contract, receipt)
decision := δ(receipt, verdict)
```

This is the intended non-circular form. If source doctrine instead types a
receipt before measurement/V/δ while requiring those outputs inside that same
receipt, record a blocking type cycle. A conforming design may use a typed
`receipt_core` followed by a final receipt, but both forms and their hash edge
must be explicit.

Hard firebreaks:

- α and β are distinct for the matter under review.
- κ is outside the cell and is not equal to α/β/γ/δ.
- shared account/model hosting is disclosed separately from functional-role
  identity; authorization does not turn unequal role types into equality.
- CC does not implement its recommendation.

## 4. L2 — contracts, CM, evidence, and receipts

The cell contract must pin:

```yaml
measurement_policy:
  cm_ref: { immutable_revision, path_or_id, content_digest }
  target_projection_ref: { immutable_revision, content_digest }
  mode: mechanical | llm | hybrid | auto
  thresholds_ref: { immutable_revision, content_digest }
  standing_required: <typed policy>
```

The measurement result must bind:

```yaml
measurement:
  cm_ref: <same immutable CM>
  target_bundle_hash: <exact ordered bytes>
  prompt_or_config_hash: <exact instrument invocation>
  scores: { alpha, beta, gamma }
  pair_deltas: { alpha_beta, beta_gamma, gamma_alpha }
  bottleneck: <one axis>
  witnesses: <TSC operational witnesses>
  consistency: <sample count and spread>
  provenance: <engine/model/version/time>
```

The receipt binds the measurement artifact by immutable digest. V reads the
measurement only through the receipt. CM measures; V gates; δ effects.

External issue/comment identities never bind content by themselves. Every
control-plane input needs a content hash plus retrievable snapshot or an
immutable repository artifact.

## 5. L3 — FSM measurement

FSM coherence is triadic at two scopes.

### 5.1 Cell FSM

- α: one state/event/guard/action vocabulary; terminal and exceptional states
  are typed; no state is simultaneously “not a state” and a reachable state
  without a distinct terminal type.
- β: exactly one authority owns each transition. Commands submit typed requests
  to the table/controller; they do not implement competing transition paths.
  Guards match contracts, V, δ, labels, and receipts.
- γ: every nonterminal state has an owned exit or explicit wait condition;
  repair, recovery, rejection, override, and resumption preserve history.

Required mechanical checks:

1. state closure: every target state is declared or typed terminal;
2. totality: every declared state/event has an explicit result;
3. determinism: at most one rule wins for a fact snapshot, or priority is
   normative and overlap-tested;
4. reachability and no accidental dead states;
5. authority uniqueness per transition;
6. positive/negative guard fixtures;
7. command↔table parity;
8. class/domain orthogonality;
9. State-A→State-B migration preserving live cycles.

### 5.2 Wave FSM

Apply the same checks plus:

1. total mapping from CC dispositions to wave transition requests;
2. child receipt/completion aggregation rule;
3. dependency/gate prerequisites for dispatchability;
4. holding and replanning exits;
5. completion and complete-with-residuals predicates;
6. operator authority only at declared irreversible/scope-expanding gates;
7. no CC mutation of wave state—the FSM owns transition.

A sequence of state names is not an FSM. It must define events, guards,
actions, authority, invalid transitions, and terminal semantics.

## 6. L4 — execution and receipt stream

Bind the exact matter bytes, complete receipt directory, review-content
snapshots, PR/issue history snapshots, checks, and state timeline. Check:

- the instance's output type matches its class;
- role claims match actual authorship/review;
- every review names the exact matter SHA;
- γ binds both matter and immutable review bytes;
- repairs narrow or explain any widening;
- historical artifacts are explicitly superseded;
- State-A limitations are true and do not redefine target invariants;
- the emitted disposition is within CC authority.

## 7. Cross-level hard invariants

These are gating, not averaged:

1. one CCNF kernel for all classes;
2. class determined by canonical output telos;
3. class and matter domain orthogonal;
4. α≠β;
5. κ outside and κ≠α/β/γ/δ;
6. ε is an observation/projection, not a cell role or CC itself;
7. immutable matter, review, CM, target-bundle, measurement, and receipt links;
8. CM/V/δ separation with a type-correct data path;
9. single dependency authority;
10. single transition authority per FSM edge;
11. complete wave-FSM semantics before “executable” claims;
12. honest shipped/specified/illustrative partition;
13. CC judgment distinct from FSM transition and operator acceptance.

## 8. Aggregation and disposition

For each level, TSC computes `C_sigma`. Cross-level aggregate is the geometric
mean of L0-L4. Report the lowest level and axis as the dominant bottleneck.

Disposition rules:

- `complete`: every hard invariant passes and formal TSC Operational ACCEPT
  has standing.
- `complete_with_residuals`: invariants pass; only bounded non-gating debt
  remains and standing policy permits it.
- `hold`: target or assurance is frozen but one or more repairable hard
  invariants fail.
- `request_planning`: relations/FSM/migration require redesign before work.
- `request_working`: design is coherent but an implementation artifact is
  absent or defective.
- `request_human`: methodology/authority choice is genuinely undecidable.
- `block`: continuing would destroy evidence, violate authority, or act on a
  moved target.

For ratification, a high score cannot override a failed invariant. Name the
invariant, evidence, level, and minimum next MCA.

## 9. Calibration commons required before standing

At minimum:

- positive: coherent WC artifact cycle;
- positive: coherent PC wave with complete typed child contracts;
- positive: coherent CC judgment with pinned CM and immutable receipt chain;
- negative: PC emits canonical artifact instead of relation graph;
- negative: κ/α or α/β role equality;
- negative: receipt requires its own downstream outputs;
- negative: mutable comment ID used as content revision;
- negative: wave state list without transitions/guards;
- negative: command and table own the same transition differently;
- comparative: each positive must outrank its paired negative.

Until these are independently authored or held out, the CM remains
`house-authored-public-commons` with no off-diagonal standing.
