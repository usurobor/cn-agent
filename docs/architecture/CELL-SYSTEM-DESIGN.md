# Cell System Design

**Version:** 0.1.0-draft
**Status:** Proposed canonical design · pre-ratification · not yet an implementation contract
**Doc class:** Architecture design
**Evidence baseline:** `claude/cds-case3-rented-beta@5b0274c215839a193009ff3484e3f2ed6bcb87df`
**Owns:** The end-to-end CNOS realization of a cell: definition, construction, admission,
subject handling, cognition, methodology, assessment, evidence, closure, failure, and
composition.
**Does not own:** The substrate-independent CCNF kernel, TSC/Coh methodology semantics,
provider-specific model behavior, or application-specific issue and design content.

This document is the proposed design authority for the CNOS cell system. It is intentionally
ahead of the runtime. Statements marked **current** describe code that exists on the stated
baseline. Statements marked **target** define the design to implement after ratification.
Nothing in this draft authorizes Case 3 implementation, merge, release, or boundary effects.

The governing companions are:

- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` — the
  substrate-independent `alpha -> beta -> gamma -> V -> delta` kernel;
- `docs/papers/CELL-OF-CELLS.md` — the recursive system model and downward-contract /
  upward-receipt relation;
- TSC Coh/CM — the property and methodology language this design intends to consume;
- `docs/architecture/CELL-RUNNER-CASES.md` — executable-case history, not architecture
  authority;
- `docs/architecture/CDS-CELL-MIGRATION.md` — migration history, not architecture
  authority.

---

## 1. Purpose

CNOS has a functional CCNF kernel and a generic one-episode runner, but it lacks one design
that says what a complete cell is. Implementation has consequently discovered system
boundaries one defect at a time:

- the producing seat was unable to run the engineering checks required by its task;
- a reviewing seat was asked to decide claims its input could not support;
- an issue was mentioned but not actually carried into the episode;
- candidate state was measured against mutable `HEAD` rather than the pinned base;
- admission, evidence, telemetry, and failure were added without a common type model;
- skills were treated as independent prompt lists rather than one assessable methodology;
- code, CUE, fixtures, and prose repeatedly described different languages or guarantees.

These are not local defects. They are consequences of an incomplete system definition.
This document closes that gap before more runtime work proceeds.

### 1.1 Goals

The design must:

1. define a cell conceptually and as a set of typed artifacts;
2. name every component, its one responsibility, and its forbidden knowledge;
3. distinguish authored definition, normalized IR, compiled plan, run input, matter,
   evidence, telemetry, and receipt;
4. define construction without reflection, service location, or a second binding plane;
5. make admission and assessment honest about what they can establish;
6. give production enough capability to perform the work while keeping authority outside
   cognition;
7. preserve the CCNF signatures and the `alpha != beta` independence boundary;
8. give every good and bad path a typed, observable outcome;
9. let bootstrap skills be replaced by Coh without changing the cell boundary;
10. support CDS, writing, planning, and later cell classes through one runner.

### 1.2 Non-goals

This design does not:

- ratify CCNF or TSC doctrine;
- implement Coh, Case 3, repair recursion, or the main cell loop;
- define a generic model router or arbitrary command runner;
- make GitHub, Git, Claude, Codex, or CUE part of the generic kernel;
- make telemetry authoritative;
- allow a cell to amend the contract that authorizes it;
- guarantee host containment merely because a process is fresh or a workspace is disposable;
- require mechanical-only cells in v0;
- freeze the final source-language syntax.

### 1.3 Constraints

- The CCNF kernel remains `matter := alpha(contract)`, `review := beta(contract, matter)`,
  `receipt := gamma(contract, matter, review, evidence)`, `verdict := V(contract, receipt)`,
  `decision := delta(receipt, verdict)`.
- The generic runner remains blind to domain semantics, providers, Git, issues, and skills.
- Cognition rents capability but never owns identity, evidence, validation, or boundary
  authority.
- A component is declared once. Sibling fields in its declaration are constructor arguments;
  no separate `bindings.alpha` or component-reference plane exists.
- Exact installed skill bytes, requested model selectors, component policies, contracts, and
  subject snapshots are content-bound before use.
- Current behavior and target behavior remain visibly distinct until implementation catches up.

### 1.4 Challenged assumptions

This design rejects seven assumptions that produced the current thrash:

1. a prose goal or issue alone is an executable implementation contract;
2. a reviewer can recover information that was never supplied;
3. withholding a tool is equivalent to containing the process;
4. separate alpha and beta skill lists can remain coherent by convention;
5. acceptance criteria enumerate every legitimate impact of a change;
6. source configuration, normalized IR, compiled plan, and run request are one object;
7. another implementation round is cheaper than settling these boundaries first.

---

## 2. Conceptual model

### 2.1 Definition

> A cell is a bounded, receipt-producing transformation that applies a declared methodology
> to an admitted contract and subject, produces candidate matter, assesses that matter
> independently, and crosses its boundary only through mechanical closure, validation, and
> decision.

A cell is not a model, agent persona, task, issue, workflow, or provider. Those may appear
inside or around a cell. Cell identity is determined by its contract kind, subject kind,
matter kind, methodology, component declarations, protocol, and authority boundary.

The semantic phases are:

```text
admit -> produce -> assess
```

The fixed closure tail is:

```text
gamma.close -> V.validate -> delta.decide
```

`admit`, `produce`, and `assess` are phases, not necessarily three cognitive calls.
Admission may be mechanical, cognitive, or composed. Production and assessment are cognitive
in the v0 CDS realization. The type system does not prohibit a future mechanically realized
cell when uniform composition and receipts pay for the boundary.

### 2.2 What requires a cell

A purely mechanical function normally remains a function, workflow step, or property
provider. A cell earns its cost when at least one of the following is true:

- cognition makes the transformation nondeterministic;
- independent discrimination is load-bearing;
- evidence must be bound across an authority boundary;
- the result will become matter at a parent scope;
- typed refusal, custody, and replay history are needed.

V0 assumes every deployed cell contains at least one cognitive component. This is a product
constraint, not an ontological claim that mechanical cells are impossible.

### 2.3 Invariants

Every conforming cell preserves these invariants:

1. **One contract:** production and assessment receive the same frozen admitted contract.
2. **One methodology:** the WORK methodology is one normative bundle from which the
   constructive and adversarial views derive. Admission judges the input contract rather than
   the work, so it declares its own separate bundle; the invariant forbids two authorities
   over one obligation set, not two obligation sets over different questions.
3. **One matter channel:** assessment sees production only through sealed matter.
4. **Independent assessment:** no private producer state, session, or self-report crosses to
   assessment.
5. **Runtime-observed evidence:** seats may propose claims; runtime-owned components create
   authoritative evidence.
6. **Mechanical closure:** gamma, V, and delta are re-derivable mechanics in v0.
7. **No self-authorization:** cognition cannot validate its own receipt or effect its result.
8. **Explicit uncertainty:** a component that lacks sufficient information returns
   `incomplete` or `unverified`; it does not guess.
9. **Immutable scope lift:** only a validated receipt and permitted decision cross upward.
10. **Visible refusal:** every semantically handled invocation constructs one typed outcome,
    including pre-episode refusal; process/transport failure may still prevent delivery.
11. **Content-bound execution:** contract, subject, methodology, skills, and execution policy
    are digest-bound before cognition runs.
12. **Stateless episodes:** immutable compiled definitions may be shared, but provider calls,
    handles, and all mutable/live component state are fresh per attempt unless a future
    contract explicitly introduces state.

### 2.4 Distinct channels

Four channels must never be collapsed:

| Channel | Meaning | Authority |
|---|---|---|
| Matter | The candidate semantic result | proposed, not yet transmissible |
| Evidence | Runtime-observed witnesses for or against contract claims | authoritative only through named producers and digests |
| Telemetry | Progress, tool, timing, usage, and diagnostic events | non-authoritative observation |
| Receipt | Gamma's immutable envelope binding contract, matter, assessment, evidence references, provenance, record, and digest | authoritative input to V; it does not contain verdict or decision |

Telemetry becomes evidence only when an explicit collector or property provider emits a typed
artifact candidate, the runtime seals it with positional provenance, and gamma binds its
digested bytes. Logging does not confer authority.

---

## 3. Status truth

This draft describes both an executed baseline and a target realization.

| Surface | Current baseline | Target in this design |
|---|---|---|
| Kernel | Functional single-episode `cellkernel` with frozen inputs, sealed station values, digest, V, delta, and self-verification | Retained |
| Generic runner | Parses, resolves, builds fills, runs one episode, emits a closure | Also emits typed pre-episode refusal outcomes and consumes a compiled plan |
| CDS production | `cds.patch` builds cognition, installed skills, and a disposable Git worktree | Full engineering capability inside an explicit execution substrate; subject handling becomes a declared component |
| CDS assessment | Branch implementation has a cognitive `cds.review` fill that reconstructs a bounded evaluation view from `(contract.subject, matter)` before renting, and is offered no tools | Fixed bootstrap falsifier derived from the same methodology, later replaced by Coh execution |
| Admission | Branch code structurally admits the typed CDS issue and the pinned subject inside seats; failure emits no closure | Peer admission component before the episode; success creates an admitted contract, refusal creates a receipt |
| Methodology | Independent seat skill lists | One compiled bundle with deterministic constructive/adversarial projections |
| Oracles | CI/corpus scripts outside the cell | Methodology-declared property providers invoked inside assessment and bound as evidence |
| Subject | `contract.subject` carries a `git.snapshot/0.1` pinned once at construction; both stations receive those bytes, and the adapter also reconstructs assessment's view from them. The adapter is wired at the composition root rather than declared in source | Declared subject/matter adapter shared by production and assessment through a pinned subject contract |
| Containment | Disposable worktree and bounded argv; no general host containment proof | Explicit substrate policy; full capability is permitted only inside the named boundary |
| Coh | Design direction only; no general runtime | Embedded methodology/property graph; runner ownership remains in CNOS |

The target columns are not claims about `main` or the active Case-3 branch.

---

## 4. Seven inventories

The inventories are the checkable spine of the design. Narrative sections may explain them
but may not contradict them.

### 4.1 Component inventory

| Component | One reason to change | Owns | Must not know |
|---|---|---|---|
| Invocation runner | Operator invocation surface changes | argument reading, input decoding, outcome encoding, exit mapping | domain meaning, providers, Git, issue semantics |
| Source normalizer | Cell source language changes | exact keys, hole grammar, parameter domains, canonical normalized IR | installed skills, providers, filesystem |
| Plan compiler | Construction/linking rules change | component lookup, strict component decode, skill loading/digests, policy resolution, immutable plan | episode matter or seat reasoning |
| Composition root | Available runtime implementations change | one explicit registry and constructor wiring | component internals |
| Admission component | Input contract policy changes | structural/cognitive admission and admission receipt | how work is produced or judged |
| Subject/matter adapter | Subject or matter representation changes | pin, materialize, measure, reconstruct, structural validation, disposal | contract semantics, prompts, verdicts |
| Methodology compiler | Property/skill composition changes | one bundle, its digest, static methodology catalog, projections, checker declarations, episode-catalog merge policy | provider credentials or episode effects |
| Producing component | Production strategy changes | constructive prompt/context and producer artifacts | assessment internals or boundary decision |
| Assessing component | Falsification strategy changes | reconstruction/checker orchestration, findings, uncertainty | producer private state or authority to accept |
| Cognition adapter | Provider process protocol changes | fixed argv recipe, model selector, bounded I/O, cancellation, execution mode | cell meaning, evidence policy, worktree semantics |
| Skill loader | Installed-skill format changes | exact canonical refs, ordered bodies, content digests | why a skill was selected |
| Protocol kernel | CCNF changes | frozen contract, sealed matter/review, gamma, V, delta, status lift, verification | providers, Git, task schema, prompts |
| Evidence store | Custody backend changes | immutable bytes, typed refs, content digests, retrieval | whether evidence proves a claim |
| Telemetry sink | Observability backend changes | bounded event transport, redaction, retention | verdict or decision |
| Effector | Boundary integration changes | applying only an already-permitted decision | validity, methodology, cognition |
| Epsilon/learning consumer | Cross-episode learning policy changes | observing accepted receipt streams and proposing later work | changing the closed episode |

No component may gain a second reason to change merely for convenience. In particular,
subject handling does not belong privately to the producing fill, and provider selection does
not belong in the generic runner.

### 4.2 Value catalogue

| Value | Produced by | Consumed by | Custody / digest rule |
|---|---|---|---|
| `CellSource` | author or future source compiler | normalizer | source bytes retained and hashed |
| `NormalizedCellIR` | normalizer | plan compiler, CUE oracle | canonical bytes and schema version hashed |
| `CompiledCellPlan` | plan compiler | admission and episode runtime | immutable; binds component policies, skills, projections, adapter support, and requested models |
| `UntrustedRunInput` | invoker/adapter | subject resolver and admission | raw envelope digest retained even on rejection |
| `SubjectSnapshot` | subject/matter adapter | admission, admitted contract, materialization, assessment reconstruction | exact immutable identity plus digest |
| `AdmissionReceipt` | admission component | outcome custody; episode record on success | always immutable; successful receipt binds admitted-contract digest |
| `AdmittedContract` | admission component | production, assessment, V | normalized domain contract, subject, obligations, scope, non-goals, and required evidence bound once |
| `EpisodeRequest` | runtime | kernel | plan digest + admitted contract + attempt identity |
| workspace handle | subject/matter adapter | producing component only | runtime-local; never serialized |
| `Matter` | subject/matter adapter's measurement of candidate state | assessment, gamma, receipt | bounded immutable value or content-addressed `MatterRef`; no evidence authority implied |
| evaluation view | subject/matter adapter reconstruction | assessment internals | derived from `(contract.subject, matter)`; no independent authority |
| `PropertyCheckObservation` | property provider inside assessment | falsifier, then runtime sealing | typed result + provider/provenance + measured artifact digest; not yet CCNF evidence |
| `AssessmentResult` | assessing composite | runtime sealer | canonical `Assessment` plus unsealed positional artifact candidates |
| `Assessment` | runtime sealer from beta result | gamma, V | bounded findings, citations, complete coverage, and unverified properties |
| `EvidenceRef` | runtime producers | gamma, V, parent | content-addressed; producer authority explicit |
| `EpisodeRecord` | gamma | receipt digest and V | canonical immutable record nested in the receipt |
| `Receipt` | gamma | V and delta | binds record + one digest; contains no verdict or decision |
| `EpisodeClosure` | receipt + V + delta | invoker, parent cell, effector | binds receipt, verdict, decision, status, and repair metadata |
| `AdmissionRefusalReceipt` | admission path | invoker, custody, learning | terminal pre-episode outcome; never masquerades as a closure |
| `RunFailureReceipt` | invocation/runtime fault path | invoker and operational custody | structured fault when result channel remains available; never masquerades as semantic refusal or closure |
| `TelemetryEvent` | runtime/provider/tool observers | telemetry sink | non-authoritative; loss never changes verdict |
| credentials | operator substrate | provider child process | never stored in source, plan, receipt, evidence, or telemetry |

### 4.3 Authority table

| Surface | Enforces | Declares | Claims nothing about |
|---|---|---|---|
| Normalizer | exact source grammar, parameter presence/domain, no unresolved holes | canonical normalized form | constructability or semantic sense |
| CUE validation | independent structural constraints over source/IR/receipts | schema conformance | runtime effects or semantic correctness |
| Plan compiler | known fills, strict component keys, installed skills, supported subject/matter pair, allowed provider policy | complete constructed plan | outcome of cognition |
| Admission structural gate | contract shape, required sections, unique criterion ids, verification-route presence, referenced snapshot shape | well-formedness | semantic executability |
| Admission cognitive gate | bounded attestation procedure and typed output | `attested` judgement with findings | proof of correctness or completeness |
| Cognition adapter | exact fixed process recipe, bounded I/O/time, requested model selector, no cell-supplied argv | offered capabilities and baseline configuration | host containment or model correctness |
| Subject/matter adapter | exact pinned snapshot, measured matter, contract-bound reconstruction semantics, structural matter validity | workspace disposability | effects outside the containment boundary; bounded reconstruction I/O may still be required |
| Assessment | every declared property receives `pass`, `finding`, or `unverified` | reasoned findings and citations | unobserved facts |
| Property-check provider | typed check against a named artifact and provenance | result and proof limit | properties outside its declared predicate |
| Gamma | canonical record construction and digest | closure provenance | semantic validity |
| V | receipt integrity, contract equality, evidence authority/availability, assessment policy | `PASS` or `FAIL` with predicates | permission to effect |
| Delta | CCNF boundary action for receipt/verdict pair | five-value decision type; normal v0 policy emits `accept | reject | repair_dispatch` | independent re-judgement of evidence |
| Effector | action is in the permitted decision and operator policy | effect result | validity of the underlying work |

Capability availability, tool approval, operating-system containment, evidence authority, and
boundary permission are different facts. No sentence may use one as proof of another.

### 4.4 Phase specifications

#### Admission

| Field | Contract |
|---|---|
| Question | Is the proposed run input executable under this cell definition? |
| Input | `UntrustedRunInput`, optional pinned `SubjectSnapshot`, compiled admission policy |
| Output | `AdmissionReceipt` with `admitted | rejected | incomplete`; mechanism faults take `RunFailureReceipt` |
| Structural sufficiency | Complete for declared structural predicates |
| Cognitive sufficiency | Attested, not proven; must expose uncertainty and findings |
| Independence | Sees no candidate matter and cannot perform production |

#### Production

| Field | Contract |
|---|---|
| Question | None; production proposes matter rather than authorizing it |
| Input | frozen `AdmittedContract` and constructive methodology projection; its composite realization materializes the subject |
| Output | seat result plus runtime-measured `Matter` and producer evidence |
| Capability | Full practical domain capability inside the declared execution boundary |
| Failure discipline | Provider/tool/substrate malfunction is distinct from valid empty or partial matter |
| Independence | Cannot see assessment internals, verdict, or decision |

#### Assessment

| Field | Contract |
|---|---|
| Question | For every declared obligation, does the candidate matter satisfy it? |
| Kernel input | Exactly `(AdmittedContract, Matter)` |
| Internal derivations | evaluation view from `contract.subject + matter`; property-check observations from declared providers |
| Output | property coverage plus `pass | finding | unverified` for every obligation; mechanism faults take `RunFailureReceipt` |
| Sufficiency | A property may pass only when its declared route resolves against matter/view or a valid property-check observation |
| Independence | No producer session, workspace handle, private transcript, or self-reported evidence |

Assessment remains the CCNF function `beta(contract, matter)`. Reconstruction and property-check
execution are internal to that function's realization and derive only from the contract and
matter. They do not add a third kernel input.

### 4.5 Rule ownership

| Rule | Normative owner | Runtime mirror | Agreement mechanism |
|---|---|---|---|
| Cell source structure | versioned CUE `#CellSource` | normalizer | shared positive/negative corpus |
| Hole/name grammar | CUE definition | parser/resolver | same corpus plus direct diagnostic tests |
| Fill-owned keys | each fill's schema | strict fill decoder | case/unknown/trailing-data negatives |
| Provider/model pairing | cognition component schema | adapter constructor | shared constructor/schema fixtures |
| Subject/matter compatibility | subject/matter adapter descriptor | plan compiler | positive pair plus unsupported-pair negative |
| Admission contract | domain schema | admission implementation | single-reason negative corpus and semantic fixtures |
| Methodology obligations | one bundle/CM + static methodology catalog + deterministic episode merge | constructive/adversarial projections and beta coverage | bootstrap ordered-body identity/catalog closure; Coh property-id coverage and projection-digest tests |
| Property-check contract | methodology catalog/checker declaration | linked provider implementation | provider fixtures, unavailable/fault cases, receipt validation |
| Closure shape | CCNF/CUE receipt schema | kernel | live output vet plus mutation negatives |
| Boundary action algebra | CCNF/boundary schema | delta/lift | full verdict x action table |

A comment that two authorities agree is not an agreement mechanism. At least one executable
witness must fail when they diverge.

### 4.6 Gate specification

Every gate states which artifact it measured and how that artifact was obtained.

**Provenance rule.** A gate must build or fetch the artifact it measures from the revision
under review; it may not measure whatever happens to be present. This is not a general
caution. The eighth observed defect was exactly this: the shared corpus ran `./cn` from the
repository root rather than building it, so local runs reported every CLI check green against
a binary that could predate the change, and a mutation test passed after the guard it was
testing had been deleted. CI, which builds first, was never affected — which is why the gap
survived. A gate whose artifact provenance is unstated reports on an unknown object.

| Gate | Artifact | Green proves | Green does not prove |
|---|---|---|---|
| Source/IR CUE vet | exact source or normalized IR bytes | structural schema conformance | construction, runtime behavior, semantic adequacy |
| Plan compile | installed packages and normalized IR | components construct and policy/skill/adapter dependencies resolve | provider success or contract admission |
| Admission shape | exact run-input bytes and subject-reference bytes | declared contract shape and reference syntax hold | subject availability or semantic executability |
| Admission relational | normalized contract plus pinned subject | declared contract/subject relations hold | semantic executability |
| Admission cognitive | same bound input | one bounded attestor found no blocking incoherence, or named it | truth or complete coverage |
| Property checker | reconstructed candidate view identified in receipt | its one declared predicate for that artifact | any undeclared property |
| Assessment coverage | methodology + assessment output | every property has a typed disposition | that every `pass` is mechanically true unless its provider proves it |
| V | contract + receipt + dereferenced evidence | receipt integrity and validation policy | authorization to effect |
| Corpus / CLI gate | fixtures plus an executable **built from the revision under review** | both authorities agree over the corpus and the live cells close and vet | anything about a rented provider; the corpus rents only the deterministic fake |
| Exact-head CI | source at named commit | encoded tests and checks pass for that commit | absence of unencoded defects or cognitive reproducibility |

### 4.7 End-to-end flow

| # | Step | Success | Typed failure/outcome |
|---|---|---|---|
| 1 | Mint runtime invocation id; capture invocation/source | authoritative correlation + source bytes | usage/read failure; transport may still prevent a result |
| 2 | Normalize and resolve holes | `NormalizedCellIR` | structural or parameter rejection |
| 3 | Compile plan | immutable `CompiledCellPlan` | construction/dependency rejection |
| 4 | Decode and digest run input | `UntrustedRunInput` + raw/reference digests | input read failure |
| 5 | Validate cheap contract/reference shape | structurally admissible candidate | rejected refusal receipt |
| 6 | Resolve/pin subject | `SubjectSnapshot` | malformed rejected; unavailable incomplete; adapter fault |
| 7 | Run relational/semantic admission | `AdmissionReceipt(admitted)` + `AdmittedContract` | rejected/incomplete refusal or runtime fault |
| 8 | Mint attempt/station identities and freeze inputs | `EpisodeRequest` | identity/freeze failure |
| 9 | Alpha materializes, produces, and measures | sealed `Matter` and producer artifact candidates | producer/substrate/measurement fault |
| 10 | Beta reconstructs assessment view | contract-bound candidate view | reconstruction fault |
| 11 | Beta runs declared property checks | typed observations and artifact candidates | fail finding, unavailable/unverified, or checker fault |
| 12 | Beta assesses all properties | complete assessment coverage | assessor fault or explicit unverified properties |
| 13 | Gamma creates receipt | canonical record + one digest | record construction fault |
| 14 | V validates | verdict | typed predicate failures |
| 15 | Delta decides | deterministic v0 boundary decision | inconsistent implementation is a fault |
| 16 | Self-verify and encode | `EpisodeClosure` | output-channel fault |

An admission refusal is not an episode closure: no episode exists before a contract is
admitted. Both are members of the run outcome algebra defined below.

Effection is downstream of the cell. A permitted closure may be passed to an effector, whose
separate `EffectResult` does not belong to `RunOutcome` and cannot change the closure.

---

## 5. Type model

The design distinguishes authored, derived, compiled, runtime-local, and closed artifacts.
Using one JSON object for all five phases would erase the point at which authority and effects
enter the system.

### 5.1 Source and construction types

```text
CellSource
  format                 source-language version
  identity               cell id + design version
  protocol               selected CCNF/receipt protocol
  parameters             bounded semantic or execution holes
  input                   contract + subject declarations
  admit                   one inline component declaration
  subject_adapter         one inline adapter declaration
  methodology             one inline methodology declaration
  produce                 one inline producer declaration
  assess                  one inline assessor declaration
  output                  matter kind; methodology/protocol and admitted contract own evidence requirements

NormalizedCellIR
  canonical source semantics after defaults and hole substitution
  closed keys and concrete values
  no unresolved symbols, secrets, paths, process handles, or arbitrary commands

CompiledCellPlan
  normalized_ir_digest
  constructor/policy identities
  supported subject/matter adapter pair
  ordered installed skill refs + content digests
  one methodology digest
  static MethodologyCatalog + digest
  EpisodeAssessmentCatalog merge-policy digest
  deterministic constructive/adversarial projection descriptors + digests
  portable checker requirements plus linked provider identities/digests
  requested cognition adapter/model per cognitive component
  plan format/canonicalization version
```

`NormalizedCellIR` is a portable semantic artifact. `CompiledCellPlan` is an immutable,
reusable link product: exact factories, installed skill/methodology bytes, supported adapter
pairs, requested provider/model selectors, and portable policy requirements. It contains no
live subject handle, credential, sandbox grant, or process. Per-attempt substrate resolution
occurs after `UntrustedRunInput` arrives and is recorded in admission/episode provenance. The
plan does not introduce a second cell ontology; it proves that the portable declaration can be
constructed by this runtime if a run-specific substrate satisfies its requirements.

### 5.2 Run and closure types

```text
UntrustedRunInput
  client_request_id?        untrusted correlation only
  contract_source          domain contract bytes or content-bound ref
  subject_reference        untrusted reference supplied by the invoker

SubjectSnapshot
  kind
  adapter_id + adapter_version
  immutable_identity
  content_digest
  reconstruction_locator     repository/content identity needed by the adapter
  custody_ref?

AdmissionReceipt
  invocation_id             runtime-minted before admission
  client_request_id?
  untrusted_input_digest
  subject_reference_digest
  plan_digest
  subject_snapshot?         present only after successful pin
  policy/provenance
  outcome                  admitted | rejected | incomplete
  structural_findings[]
  semantic_findings[]
  attestation              verified | attested_unverified | not_run
  admitted_contract_digest?  required only on admitted

AdmittedContract
  contract_id
  contract_kind
  domain_contract           normalized domain payload + digest
  subject                   exact SubjectSnapshot
  scope / non_goals / invariants
  acceptance_obligations[]  stable ids + verification routes
  assessment_catalog        ordered coverage units + digest
  required_evidence[]       derived during admission from domain contract,
                            methodology/checkers, and selected protocol

EpisodeRequest
  invocation_id
  attempt_id
  compiled_plan_digest
  admission_receipt_digest
  admitted_contract
  execution_policy_snapshot

RunOutcome = AdmissionRefusalReceipt | EpisodeClosure

InvocationResult = Completed(RunOutcome) | Fault(RunFailureReceipt)

RunFailureReceipt
  format + version
  invocation_id?           present once the runtime minted it
  client_request_id?       correlation only
  phase
  failure_class
  diagnostic_code + bounded diagnostic
  known_digests            every already-known source/IR/plan/input/admission/
                           contract/subject/matter/receipt digest
```

`AdmissionRefusalReceipt` is the serialized rejected/incomplete form of
`AdmissionReceipt`. `EpisodeClosure` is available only after an admitted contract has entered
CCNF. This sum preserves visibility without pretending that a refused contract produced an
episode. `RunFailureReceipt` reports a runtime or custody malfunction without pretending that
the semantic run completed or crossed the CCNF boundary.

Construction order is acyclic: the runtime mints `invocation_id`, digests raw input and the
subject reference, validates cheap contract shape, optionally pins a snapshot, constructs and
digests `AdmittedContract`, then places that contract digest in the successful
`AdmissionReceipt`. `EpisodeRequest` binds both the contract value and admission-receipt
digest; `EpisodeRecord` preserves both. Neither artifact contains the other's digest.

### 5.3 Assessment types

```text
PropertyDisposition
  property_id
  outcome       pass | finding | unverified
  reasons[]
  citations[]
  observation_refs[]

Assessment
  methodology_digest
  projection_digest
  catalog_digest
  coverage[]    exactly one PropertyDisposition per assessment-catalog unit
  summary       derived: pass iff every disposition is pass; otherwise nonpass + reasons

AssessmentResult
  assessment
  artifact_candidates[]   unsealed, non-authoritative candidates returned positionally by beta

PropertyCheckObservation
  property_id
  provider_id + provider_digest
  measured_artifact_kind + measured_artifact_digest
  outcome       pass | fail | unavailable
  artifact_candidates[]
  proof_limit
```

A checker `fail` is a successful measurement containing negative evidence candidates.
`unavailable` means the property remains unverified. A provider fault or malformed result does
not produce an observation; it returns `RunFailureReceipt`. Beta performs valid checks inside its composite implementation;
runtime sealing assigns positional provenance, gamma binds the resulting artifacts as
evidence, and V later dereferences them. Beta is never handed a pre-existing evidence object.
These paths must not share one generic non-zero branch.

### 5.4 Subject/matter adapter contract

The generic contract is deliberately not Git-shaped:

```text
SubjectMatterAdapter
  supports(subject_kind, matter_kind) -> bool
  pin(untrusted_reference) -> SubjectSnapshot
  materialize(SubjectSnapshot) -> SubjectHandle
  measure(SubjectHandle, ProductionResult) -> Matter
  validate_matter(Matter) -> StructuralResult
  reconstruct(SubjectSnapshot, Matter) -> EvaluationView
  dispose(SubjectHandle) -> DisposalResult
```

`SubjectHandle` and `ProductionResult` are runtime-local values opaque to the generic runner;
they cross only the checked adapter/producer port. `EvaluationView` stays inside composite
beta. Any resource allocated for materialization or reconstruction is attempt-scoped and
disposed by the owning composite without changing semantic closure.
The Git adapter realizes them as a disposable worktree, workspace mutation, and reconstructed
candidate tree. A writing adapter may realize them as an immutable source bundle, output slot,
and candidate document. A planning adapter may realize them as a program-state snapshot,
typed plan value, and validated graph view.

The adapter validates structure only. Whether a patch, document, or plan is *good* remains a
methodology question.

### 5.5 Identity

- `cell_id` identifies a reusable semantic definition and version.
- `plan_digest` identifies one exact compiled realization.
- `invocation_id` is runtime-minted at invocation start; a client request id is correlation
  only and never authoritative.
- `attempt_id` identifies one admitted episode; retries mint new attempt ids.
- station ids are pairwise distinct within the attempt.
- property and evidence ids are stable inside their versioned definition.

No model mints authoritative identity. The runtime does.

---

## 6. Architecture

```text
AUTHORING                         PURE / PORTABLE

 Cell source / future .cell ──> Normalize + CUE vet ──> NormalizedCellIR
         Domain contract ─┐                                  │
       Subject reference ─┴──> UntrustedRunInput              │ link/construct
                                                                  ▼
                                                   CompiledCellPlan
                                                           │
RUNTIME / AUTHORITY                                        │
                                                           ▼
 UntrustedRunInput ──> SubjectMatterAdapter.pin ──> Admission component
                                                  │
                         refusal ─────────────────┼──> AdmissionRefusalReceipt
                                                  │ admitted
                                                  ▼
                                      AdmittedContract + AdmissionReceipt
                                                  │
 constructive methodology view ─────────────> PRODUCE (alpha composite)
                                      SubjectMatterAdapter.materialize
                                                  │ ProductionResult
                                      SubjectMatterAdapter.measure
                                                  │ Matter
                                                  ▼
                                      ASSESS (beta composite)
                         ┌─ adapter.reconstruct -> EvaluationView
                         ├─ compiled methodology/checkers
                         └─ mechanical providers -> cognitive falsifier
                                                  │ Assessment + ArtifactCandidates
                                                  ▼
                                  gamma.close -> V -> delta
                                                  │
                                           EpisodeClosure
                                         /                \
                               custody / parent         permitted effector

OBSERVATION

 provider/tool/runtime events ──> Telemetry sink ──> epsilon/learning candidates
                                  (never a verdict input by itself)
```

### 6.1 Boundary placement

- Admission is before CCNF because it decides whether a contract may enter the episode.
- The subject/matter adapter is not a station. It is substrate used by admission, production, and
  assessment.
- Alpha produces; it does not decide whether its product is sufficient.
- Beta is the composite assessment function. It receives the CCNF pair `(contract, matter)`;
  its realization reconstructs the view and invokes declared property providers before any
  cognitive falsification. The generic runtime does not pass beta a third evidence argument.
- Gamma binds evidence references; it does not invent or judge them.
- V validates the receipt against the contract; delta alone chooses a permitted boundary
  action.
- The effector is outside the cell wall and cannot turn failure into validity.
- Epsilon observes receipt/telemetry streams across episodes and may propose later work. It
  cannot mutate a closed receipt.

### 6.2 Property-check placement and CCNF

CCNF explicitly fixes beta as `review(contract, matter)` and excludes a separate evidence
input. Therefore assessment is implemented as one composite beta function:

```text
beta.review(contract, matter):
  view         := adapter.reconstruct(contract.subject, matter)
  observations := methodology.run_mechanical_properties(contract, view)
  review       := falsifier.assess(contract, matter, view, observations)
  return AssessmentResult(review, observations.artifact_candidates)
```

The returned review remains beta's value. Runtime sealing assigns the returned artifact
candidates beta-positioned provenance; they then accumulate as episode evidence for gamma and
V. The cognitive falsifier may consume the mechanical results beta itself obtained;
no producer evidence channel is introduced. A future stricter CCNF reading may place all
checker interpretation in V; that would be a deliberate kernel change, not an accidental third
beta input.

### 6.3 Functional and stateless execution

At the semantic boundary, every phase is a function over immutable inputs. Adapter
implementations may perform bounded I/O without changing those value contracts:

```text
admit(plan, untrusted_input, snapshot) -> AdmissionResult
produce(contract) -> AlphaOutput
assess(contract, matter) -> AssessmentResult
gamma(contract, matter, assessment, evidence_refs) -> Receipt(EpisodeRecord, digest)
validate(contract, receipt) -> Verdict
decide(receipt, verdict) -> BoundaryDecision
```

The constructed alpha closes over its subject/matter adapter and constructive methodology
projection; its per-episode implementation internally materializes, invokes production, and
measures matter. Those constructor dependencies do not widen CCNF's outer `alpha(contract)`
signature.

Effects are localized to adapters: reading installed skills, pinning/materializing subjects,
starting provider processes, running tools, storing evidence, emitting telemetry, and applying
permitted boundary actions. No shared mutable episode object, global current provider, or
resumable model session is part of v0.

---

## 7. Definition, construction, and dependency injection

### 7.1 One-place construction

Each configurable component is declared once. Its `fill` or `kind` is a closed constructor
discriminator; sibling fields are the constructor arguments.

```json
"produce": {
  "fill": "cds.patch",
  "cognition": {
    "provider": "$provider",
    "model": "$model"
  }
}
```

The whole `produce` object is passed to the `cds.patch` constructor. The generic runner does
not know that the fill needs cognition, skills, or a workspace. The fill/compiler owns that
construction from its declared arguments and injected shared components.

Shared components are also declared once:

- `subject_adapter` is declared and compiled once; the plan instantiates attempt-local ports
  for production and assessment;
- `methodology` bytes and projection definitions are compiled once and may be shared
  immutably across attempts;
- the protocol kernel is selected once by the plan;
- each cognitive component gets a fresh provider process from its own inline cognition
  declaration.

This is constructor injection, not a reference graph. The fixed cell vocabulary tells the
compiler where each constructed value goes. V0 has no arbitrary component refs or circular
dependencies.

### 7.2 What CNOS borrows from Spring

The useful Spring analogy is manual, explicit dependency injection:

- one declaration per component;
- one composition root;
- closed typed constructor arguments;
- validation and linking before instantiation;
- constructor injection rather than service lookup;
- immutable compiled wiring;
- attempt scope for every live process, handle, and mutable component instance.

CNOS rejects the rest:

- reflection or class names in JSON;
- package scanning and implicit registration;
- autowiring by name or type;
- mutable application contexts or singletons;
- lifecycle callbacks and hidden proxies;
- ambient property files that silently override the cell;
- arbitrary executable, argv, environment, or capability maps;
- a component looking up another component at runtime.

The compiled plan is the auditable wiring graph that a general-purpose container would hide.

### 7.3 Composition root and registry

The runtime owns one explicit immutable registry of supported constructor descriptors. A
descriptor includes:

```text
id, version, declaration schema, constructor, capability requirements,
supported contract/subject/matter kinds, and policy digest
```

Registration does not grant authority. The plan compiler still validates policy and type
compatibility. A fill cannot widen its tool grants, load arbitrary skills, or choose a binary
outside its descriptor.

---

## 8. Source language, JSON IR, and CUE

### 8.1 Artifact pipeline

```text
future .cell / direct JSON
        -> CellSource
        -> Normalize + resolve + default
        -> CUE structural validation
        -> NormalizedCellIR
        -> link installed implementations and policies
        -> CompiledCellPlan
        -> bind UntrustedRunInput
        -> InvocationResult
```

The future F#/OCaml-like computation-expression syntax is an authoring surface. It must lower
to the same `NormalizedCellIR` as direct JSON. Source syntax does not become runtime authority.

### 8.2 Why JSON is the interchange IR

JSON is selected for v0 because it is:

- language-neutral between a future compiler and the Go runtime;
- independently checkable by CUE;
- inspectable, diffable, and easy to preserve in receipts;
- compatible with closed schemas and exact-key rejection;
- suitable for deterministic content addressing under a versioned canonical encoder;
- already proven by the cell runner and TSC experiments.

JSON is not the preferred final authoring language, a constraint solver, a DI container, an
executable plan, or a place for host paths and secrets.

### 8.3 Canonicalization

Content addressing applies only after a typed value and encoder version are known:

```text
CanonicalJSON(format_version, typed_value) -> bytes
```

The current Go implementation has more than one concrete encoder convention: episode records
use schema/struct order plus `json.Marshal`, while fill declarations re-decode and compact
maps. The implementation plan must either unify these under one named cross-language contract
or version the conventions explicitly. This draft does not pretend ordinary source JSON bytes
are canonical or silently choose JCS.

### 8.4 Validation ownership

| Concern | Owner |
|---|---|
| syntax lowering and default materialization | source compiler/normalizer |
| exact structural conformance | CUE plus parser mirror |
| hole substitution and domain membership | normalizer/resolver |
| constructor, provider, adapter, and capability compatibility | plan compiler |
| operational behavior | constructed component |
| receipt integrity and validation | kernel/V |

CUE must not select providers, execute workflows, infer ambient resources, or authorize
effects.

---

## 9. Cell source shape

The following JSON is normative for the design vocabulary but not yet a frozen schema. It is
complete enough to describe CDS without hidden semantics in the generic runner.

```json
{
  "format": "cnos.cell.source/0.1",
  "id": "cnos.cds.implement",
  "version": "0.1",
  "protocol": "cnos.cell.ccnf/0.1-draft",
  "parameters": {
    "language": {
      "required": true,
      "domain": [
        "cnos.eng:eng/go",
        "cnos.eng:eng/ocaml",
        "cnos.eng:eng/typescript"
      ]
    },
    "style": {
      "required": false,
      "default": "cnos.eng:eng/write-functional"
    },
    "provider": {
      "required": true,
      "domain": ["claude-cli", "fake"]
    },
    "model": {
      "required": false,
      "default": ""
    }
  },
  "input": {
    "contract": {
      "kind": "cnos.cds.implementation-contract/0.1"
    },
    "subject": {
      "kind": "git.snapshot/0.1"
    }
  },
  "admit": {
    "fill": "cds.work-contract-admission",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    },
    "skills": [
      "cnos.cdd:cdd/issue",
      "cnos.cdd:cdd/design"
    ]
  },
  "subject_adapter": {
    "fill": "git.patch-workspace",
    "version": "0.1"
  },
  "methodology": {
    "kind": "skills.methodology/0.1",
    "skills": [
      "cnos.eng:eng/code",
      "cnos.eng:eng/test",
      "$language",
      "$style"
    ],
    "checkers": [
      {"unit": "checker:build", "checker": "project.build/0.1"},
      {"unit": "checker:tests", "checker": "project.test/0.1"}
    ]
  },
  "produce": {
    "fill": "cds.patch",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    }
  },
  "assess": {
    "fill": "methodology.falsify",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    }
  },
  "output": {
    "matter": {"kind": "git.patch/0.1"}
  }
}
```

`Assessment` is beta's one canonical review value in the episode record. It is not duplicated
as an independently authoritative evidence artifact. Property-check artifacts that support
the assessment are declared by the methodology, sealed positionally by the runtime, and bound
as evidence by gamma. Matter has its own record slot and is never duplicated as evidence
merely to prove that it exists.

Each `checker` string names a portable capability/evidence contract, not an installed
executable. The plan linker selects a concrete provider under host policy and records its
exact id, version, implementation digest, and grants in run provenance.

### 9.1 Holes

| Value | Classification |
|---|---|
| language | semantic skill hole |
| style | project-policy skill hole with a default |
| provider | execution selection constrained by the cell |
| model | requested model selector; explicit in the compiled plan |
| issue/brief/design payload | domain-defined run contract, not holes |
| subject reference/base | run subject, not holes |
| timeout/tool grants | runtime policy, not user-configurable holes in v0 |
| arbitrary task skills | excluded until admission and merge semantics are designed |

The generic parameter declaration owns only presence, default, and closed-domain rules. The
consuming fill or methodology field schema owns whether a substituted value is a skill ref,
model selector, or ordinary scalar. A future source language may type holes before lowering;
v0 JSON IR does not make the normalizer understand skill semantics.

The same provider/model parameters may populate several component declarations while each
component still receives a fresh invocation. Per-station provider holes are not added until a
real need appears.

### 9.2 Run input example

The reusable cell definition above is separate from the per-run contract and subject:

```json
{
  "format": "cnos.cell.run-input/0.1",
  "client_request_id": "client-req-01J...",
  "contract_source": {
    "kind": "cnos.cds.implementation-contract/0.1",
    "issue": {
      "source": "github:usurobor/cnos#720",
      "content_digest": "sha256:..."
    },
    "design": {
      "source": "repo:docs/designs/cnos-720.md",
      "content_digest": "sha256:..."
    }
  },
  "subject_reference": {
    "kind": "git.snapshot-ref/0.1",
    "repository": "usurobor/cnos",
    "revision": "<exact-commit-sha>"
  }
}
```

The generic input has one domain contract. Its CDS schema requires the distinct issue and
design facets shown above; writing and planning define different payloads without changing
the runner. The admission membrane reads and normalizes those content-bound sources, pins the
repository identity and revision into `SubjectSnapshot`, and constructs `AdmittedContract`. A
live issue lookup, mutable branch name, or ambient checkout never substitutes for that frozen
contract.

### 9.3 Skills and methodology

Skill bodies are not independent alpha and beta configuration. The cell declares one ordered
methodology bundle. Current engineering skills mix assessable obligations, constructive
procedure, examples, heuristics, and advice; v0 must not pretend that this prose is already an
executable property graph. Compilation loads exact installed bytes and records canonical refs
and digests. It derives:

- a **constructive projection**: obligations plus guidance useful for producing matter;
- an **adversarial projection**: the same obligations expressed as falsification questions,
  provider routes, refusal conditions, and output schema.

The honest bootstrap invariant is narrower: the production and assessment components receive the same
ordered, complete skill bodies with the same content digests and fixed role wrappers. Once
Coh supplies stable property ids, every property must appear in both the constructive and
adversarial projections. Projection descriptors and digests are derived provenance, not
additional sources of truth. Coh later replaces the bootstrap projection mechanism with
compiled property semantics; the surrounding cell does not change.

### 9.4 Assessment catalog

The compiled methodology emits a closed, ordered static `MethodologyCatalog`. After admission,
a pure versioned merge creates the episode's one coverage authority:

```text
MethodologyCatalog =
  canonical methodology skill refs as coarse coverage-unit ids
  + explicitly typed checker-only units

EpisodeAssessmentCatalog =
  merge_v1(MethodologyCatalog, admitted-contract obligation ids)
```

Every checker declaration must resolve to one existing methodology unit or explicitly declare
a typed checker-only unit; the candidate's `checker:build` and `checker:tests` do the latter.
The reusable plan binds only the static catalog and merge-policy digests. Admission constructs
and binds `EpisodeAssessmentCatalog` in `AdmittedContract` and `EpisodeRequest`; beta emits
exactly one disposition per episode unit and may cite finer skill passages without creating
new units. V rejects missing, duplicate, unknown, or reordered units.

This is deliberately coarse in bootstrap: a whole skill ref is reviewable, while its prose is
not falsely treated as a machine-extracted property graph. Coh later replaces skill-ref units
with stable property ids and preserves the same exact-coverage contract. Contract obligations
remain catalog inputs rather than disappearing into the engineering methodology.

---

## 10. Contract and admission

### 10.1 Issue plus design

CDS production requires two logically distinct contract artifacts:

- **Issue:** problem, impact, status truth, source of truth, scope, non-goals, constraints,
  acceptance obligations, verification routes, and closure condition.
- **Design:** proposed system change, invariants, impact graph, alternatives, trade-offs,
  authority changes, and known debt.

A small change may carry an inline design, but it remains separately typed and digested. The
issue cannot silently choose architecture, and the design cannot silently redefine the
problem.

The `cdd/issue` and `cdd/design` skills are bootstrap authoring/review methodologies. Their
stable structured fields should be extracted into a CDS implementation-contract schema. CUE
can prove shape and local relations; cognition remains responsible for semantic executability
until Coh properties replace the bootstrap review.

### 10.2 Admission algorithm

```text
1. Runtime mints invocation_id; decode and digest UntrustedRunInput and subject reference.
2. Validate cheap domain-contract and reference shape.
3. Reject malformed input with AdmissionRefusalReceipt(rejected).
4. Ask the declared adapter to pin SubjectSnapshot.
5. Treat an unavailable well-formed subject as incomplete and an adapter malfunction as Fault.
6. Validate contract/subject relations and run semantic admission when required.
7. Construct and digest AdmittedContract.
8. Emit AdmissionReceipt containing that contract digest.
9. On rejected/incomplete, return AdmissionRefusalReceipt; do not mint an episode.
```

Structural examples include required fields, allowed enums, nonblank text, unique criterion
ids, resolvable source refs, and a verification route per criterion. Cognitive examples
include whether the problem is real, issue/design are mutually consistent, criteria are
sufficient, scope is executable, and the verification route could actually decide the claim.

### 10.3 Admission result semantics

| Outcome | Meaning | May production run? |
|---|---|---|
| admitted | all structural gates passed; required semantic gate attested success | yes |
| rejected | a decisive contract defect exists | no |
| incomplete | required source, design, subject, or semantic judgement is unavailable | no |

If the admission provider or validator malfunctions or emits invalid output, invocation takes
the `Fault(RunFailureReceipt)` path. Mechanism failure is not a semantic refusal.

Bootstrap cognitive admission is labelled `attested_unverified`. It must not borrow the
standing of gamma, V, or delta, which are mechanical and re-derivable.

### 10.4 Contract discovery during production

Acceptance criteria prove selected outcomes; they do not enumerate every affected surface.
The design impact graph, scope, non-goals, and invariants define the broader execution
boundary.

| Discovery | Disposition |
|---|---|
| Necessary change already represented by design/impact graph | continue and bind evidence |
| Necessary dependency absent from design | emit `ContractFinding`; produce a nontransmissible attempt closure with `contract_unmet` and `repair_dispatch(needs_recontract)` |
| Issue/design contradiction | emit `ContractFinding`; produce a nontransmissible attempt closure with `contract_unmet` and `repair_dispatch(needs_recontract)` |
| Optional improvement outside contract | do not implement; emit follow-up candidate |
| Forbidden effect | reject candidate matter |
| Required property check unavailable | mark property `unverified`; no acceptance |

A cell cannot amend its own contract. Partial matter may be retained as a non-transmissible
proposal. An enclosing planning/design cell or operator issues a new contract digest, and a
new episode begins. `needs_recontract` is a repair reason, not a new kernel status or decision
variant. The attempt has an immutable closure artifact, but `repair_dispatch` keeps the logical
cell open at the same scope until a child repair/re-contract is accepted and the parent can
emit a new closure.

---

## 11. Production, assessment, and closure

### 11.1 Production capability

Production must have the practical capabilities necessary to discharge its contract. For CDS
that normally means repository read/write, search, shell, build, test, formatter, and Git
inspection inside the disposable subject. A coding cell unable to run the project's tests is
not safely bounded; it is merely unable to check its work.

Capability and containment remain separate:

- tool availability says what the provider can request;
- approval policy says which requests execute without intervention;
- the subject/matter adapter says which state becomes measured matter;
- the operating-system substrate says what the process can actually touch;
- the receipt says what evidence and matter gained authority.

Full engineering capability is safe only inside a named substrate boundary. A worktree limits
what is measured; it does not by itself stop absolute-path or network effects.

### 11.2 Matter production

The model's narrative is never authoritative matter. The subject/matter adapter measures the
candidate result:

- CDS: canonical patch against the pinned commit;
- writing: canonical candidate document/revision against bound source material;
- planning: typed plan graph against a program-state snapshot.

Empty matter is a valid measurement, not an implicit success or a runtime malfunction. The
domain contract and methodology assess it normally; CDS will ordinarily find it contract-unmet,
while another cell may legitimately accept an empty semantic value. Malformed or unsupported
matter is a typed structural failure.

### 11.3 Assessment algorithm

Bootstrap beta is a fixed falsification procedure:

```text
for every property in the adversarial projection:
  obtain its declared evaluation view and/or property-check observation
  attempt to falsify the property
  emit pass | finding | unverified
  include citations and artifact candidates
reject assessment output with missing or duplicate property ids
```

The cognitive component is not asked to reproduce mechanical checks from prose. Mechanical
property providers run first inside the beta composite. A checker `fail` mechanically forces
`finding`; `unavailable` forces `unverified` unless another route declared for that property
resolves it. Semantic properties rent cognition. Cognitive output that contradicts a forced
disposition, duplicates or omits a property, or is malformed is a runtime fault and produces
no completed assessment. The observations are not a third kernel input and do not yet carry
CCNF evidence authority.
Beta returns their artifact candidates; runtime sealing assigns positional provenance, gamma
binds them as evidence, and V dereferences them. The result is a complete property coverage
map, not a free-form review essay.

### 11.4 Independence by reconstruction

Assessment independence is not blindness. It is the absence of an unaccounted producer
channel.

```text
evaluation_view = reconstruct(contract.subject, matter)
```

Production can affect this view only by changing `matter`, which is exactly what beta reviews.
Beta receives no producer session, hidden transcript, mutable workspace handle, or claimed
test result. That independence survives assessment tools only when beta receives either a
bounded reconstructed value with no filesystem tools or a genuinely isolated read-only
subject substrate. Merely omitting the producer's workspace handle does not contain an
unconstrained host file tool, and the current substrate does not prove such containment.

**Bootstrap therefore takes the first option, and this is a decision rather than an
observation.** Bootstrap beta receives the reconstructed view as a bounded value and is
offered no filesystem tools, so the independence argument rests only on channel closure and
on nothing this project cannot demonstrate. Giving beta real tools waits on an isolated
read-only substrate; that trigger is recorded in the deferred-decision table. Without this
sentence the design names a containment requirement it does not satisfy and leaves the
bootstrap undefined at exactly the point where implementation would have to guess.

### 11.5 Gamma, V, and delta

- Gamma canonically composes contract, matter, assessment, and evidence refs into
  `EpisodeRecord`, then returns `Receipt(record, one_digest)`.
- V checks record shape, digest, contract equality, required evidence availability and
  producer authority, assessment coverage, and validation policy. V passes only when coverage
  is exact and every disposition is `pass`; any `finding`, `unverified`, or missing required
  evidence yields V failure with typed reasons.
- Delta retains CCNF's five-value type: `accept | release | reject | repair_dispatch |
  override`. Normal v0 policy is a pure subset that emits only `accept | reject |
  repair_dispatch`. `release` and `override` remain held until an explicit bound grant/policy
  makes them reproducible inputs to delta; they are not moved to another component.

The selected protocol owns one total, versioned v0 function and the plan binds its policy
digest:

```text
delta_v0(receipt, PASS) = accept
delta_v0(receipt, FAIL with protocol-permitted typed repair reason) = repair_dispatch
delta_v0(receipt, other FAIL) = reject
```

The set of repair reasons and their eligibility is closed by the protocol; cognition cannot
invent one. `finding`, `unverified`, and missing-evidence predicates therefore always yield V
FAIL, while delta deterministically chooses repair or rejection from the bound typed reason.

An external effector may still require operator policy before applying a delta decision, but
that gate cannot rewrite the receipt, verdict, or decision. Before serialization, an invalid
verdict/decision pair follows CCNF's mechanical re-decision rule. A serialized or
self-verified mismatch is an integrity fault, not another policy choice. The executable
schema and delta implementation must agree on the full type even while v0 emits only its
normal subset.

---

## 12. Outcomes and failure semantics

### 12.1 Run outcome algebra

```text
RunOutcome
  = AdmissionRefusalReceipt
  | EpisodeClosure

InvocationResult
  = Completed(RunOutcome)
  | Fault(RunFailureReceipt)
```

`RunFailureReceipt` covers runtime or custody malfunctions that prevent a semantic outcome:
unreadable source, corrupt bytes, runtime identity failure, provider crash before bounded
output, evidence-store outage, or closure-encoding failure. It carries structured diagnostics
and an attempt/request id where one exists, but it is neither an admission refusal nor a CCNF
closure.

This corrects a tempting overclaim: v0 should make every *semantically handled refusal* emit a
receipt, but it cannot guarantee delivery when result encoding, process transport, or custody
itself fails. The runtime constructs `RunFailureReceipt` whenever the result channel remains
available; a final unreceipted process/transport failure is still possible and must be visible
to the caller. Custody must preserve failures without confusing them with validated episodes.

### 12.2 Failure classes

| Class | Meaning | Typical disposition |
|---|---|---|
| usage/input failure | invocation could not be decoded | `RunFailureReceipt`, no episode |
| construction failure | plan could not link or instantiate | `RunFailureReceipt`, no episode |
| admission rejection | contract is decisively inadmissible | `AdmissionRefusalReceipt(rejected)` |
| admission incomplete | required information or judgement unavailable | `AdmissionRefusalReceipt(incomplete)` |
| admission mechanism fault | provider/validator failed or emitted malformed output | `RunFailureReceipt`; no semantic refusal |
| producer malfunction | provider/tool/substrate failed | `RunFailureReceipt`; no fabricated closure |
| valid empty matter | adapter measured an empty but well-typed value | assess normally under the domain contract; never imply success |
| matter structural defect | adapter cannot validate/reconstruct candidate | `RunFailureReceipt`; no fabricated closure |
| property finding | candidate violates a property | beta finding -> V fail -> reject/repair |
| property unverified | declared route unavailable or insufficient | V FAIL with typed `unverified` reason; delta policy decides reject or repair |
| property check unavailable | declared checker could not run for an environmental reason | property unverified, not a fabricated fail/pass |
| property-check fault | checker contract malfunction | `RunFailureReceipt`; no fabricated observation |
| contract discovery | necessary work lies outside admitted design | nontransmissible attempt closure + `repair_dispatch(needs_recontract)`; logical cell remains open |
| cancellation/timeout | external stop or budget exhausted | typed failure preserving partial diagnostics |
| inconsistent verdict/decision before closure | delta pairing invalid | mechanically reject and re-decide under CCNF; never transmit |
| serialized/self-verified decision mismatch | emitted closure contradicts deterministic delta | integrity fault; do not reinterpret as policy |
| effect failure | valid decision could not be applied | downstream `EffectResult(failed)`; closure remains unchanged |

### 12.3 Retry and repair

- A retry mints a new attempt id and preserves the prior outcome.
- Deterministic construction/admission may be cached only by exact input and policy digests.
- Cognition is not replayed under the same identity as if it were deterministic.
- Repair runs under a new child contract that cites the failed receipt and allowed repair
  scope.
- Re-contracting creates a new domain-contract digest and therefore a new episode.
- No failed evidence is deleted merely because a later attempt succeeds.

---

## 13. Evidence, custody, telemetry, and learning

### 13.1 Evidence

Evidence is content-addressed and producer-typed. Every `EvidenceRef` names:

```text
id, kind, producer authority, content digest, custody locator,
measured subject/artifact digest, and optional proof-limit metadata
```

The receipt need not embed every byte. It must bind enough information for V and the parent to
retrieve the exact bytes. Missing required evidence makes V fail with a typed reason; it is
never silently ignored or converted into a third verdict state.

### 13.2 Telemetry

Telemetry may include bounded:

- phase start/stop and duration;
- provider and tool events;
- stdout/stderr tails and truncation markers;
- token/usage/cost observations;
- cancellation and timeout diagnostics;
- workspace/evidence lifecycle events.

Telemetry excludes credentials, secret environment, unrestricted prompts, and uncontrolled
model traces by default. A telemetry sink loss or dropped event is non-blocking, cannot become
`RunFailureReceipt`, and never turns pass into fail or fail into pass.

### 13.3 Epsilon and learning

Epsilon is outside one episode. It observes receipt and selected telemetry streams across
episodes to identify repeated incoherence, calibration opportunities, methodology defects, and
new planning inputs. It may propose a memory, issue, methodology patch, or child contract. It
cannot revise the original receipt or automatically promote telemetry into evidence.

Learning outputs require their own cell when they contain cognition or cross an authority
boundary.

---

## 14. Coh and property composition

### 14.1 Ownership boundary

Coh/CM owns:

- the composable property/methodology graph;
- stable property ids and dependencies;
- mechanical and cognitive provider declarations;
- checker linking/execution;
- typed measurement observations, artifact candidates, proof limits, and measurement-result
  derivation.

CNOS owns:

- admission policy and the admitted-contract boundary;
- mapping a CM measurement to admission refusal or episode assessment;
- subject custody and mutation;
- sealing artifact candidates into authoritative `EvidenceRef` values;
- production;
- episode orchestration;
- CCNF gamma/V/delta;
- receipt custody and boundary effects.

A CM measures. It does not admit, mutate, repair, authorize, or effect. Admission may invoke a
CM and map its measurement to an admission outcome; assessment may be realized by a CM. The CM
does not thereby become the cell runtime.

### 14.2 Whole-cell syntax versus runtime ownership

A future Coh or `.cell` language may provide syntax for the whole cell, including
computation-expression forms such as sequential and parallel composition. That syntax lowers
to two distinct structures:

1. cell orchestration in `NormalizedCellIR`;
2. embedded methodology/property graphs in normalized CM IR.

Syntactic ownership does not imply runtime ownership. CNOS keeps custody, admission, CCNF, and
effects even if Coh supplies the most convenient source language.

### 14.3 Bootstrap migration

Until general Coh execution exists:

1. installed skills are the normative methodology source;
2. a deterministic compiler produces constructive/adversarial views;
3. a fixed beta composite runs named mechanical commands and rents cognition for semantic
   properties;
4. assessment emits the same property coverage and evidence shapes intended for Coh;
5. when Coh is ready, replace the projection/provider backend without changing cell inputs,
   outputs, or CCNF.

This bootstrap is future-compatible only if beta is a generic falsifier of declared
properties, not a hard-coded CDS reviewer.

---

## 15. Composition and recursion

### 15.1 Parent-child protocol

```text
parent contract descends
child cell runs
child receipt ascends
parent V validates
parent delta decides whether accepted child receipt becomes parent matter
```

A child receives a bounded contract and cannot inherit its parent's full authority. A parent
does not trust a child because it dispatched it; trust is earned by the validated receipt.

### 15.2 Child proposals

Alpha may return proposed child contracts as matter. It does not invoke the runner directly or
mint authoritative child identity. The enclosing runtime/planning cell validates and dispatches
those proposals under resource and authority policy.

### 15.3 Sequential and parallel composition

At the future source-language level:

- sequential composition (`let!`) means the later child contract may consume an accepted
  earlier receipt;
- parallel composition (`and!`) means sibling contracts share only explicitly bound parent
  inputs and join through a parent receipt;
- failure and cancellation propagate according to the compiled orchestration graph, not model
  discretion.

Composition belongs in the main/planning loop or parent cell. One cell episode remains a
bounded closure, not a mutable workflow engine.

---

## 16. Security, capability, and operational semantics

### 16.1 Four separate questions

For every cognitive or tool component, the plan must answer:

1. What capability is declared?
2. What is actually installed and available?
3. What is approved by policy?
4. What is contained by the substrate?

The receipt then answers a fifth question: what observed result gained authority?

### 16.2 Provider and model

Provider/model selection is inline component configuration resolved before construction. A
cell may expose bounded holes for provider id and requested model selector. It may not supply a
binary path, shell command, arbitrary argv, environment map, or permission override.

The runtime maps a closed provider id to a typed adapter. A small switch over the supported
adapters is preferable to a plugin framework until real provider growth earns one. Installation
and authentication are deployment concerns; execution fails clearly when the selected adapter
is unavailable.

### 16.3 Reproducibility

- Normalization, linking, reconstruction, mechanical checks, canonical encoding, gamma, V,
  and delta are deterministic only under identical pinned inputs, installed
  implementation/policy digests, and hermetic or fully resolved substrate facts. Resolution
  of a moving ref and non-hermetic tools are observations whose exact results become
  provenance, not claims of reproducibility.
- Cognitive output is not reproducible merely because provider/model are recorded.
- A cognitive receipt proves which declared inputs and policy produced this observed result;
  it does not prove the provider would repeat it.
- Same normalized inputs and linked dependencies should yield the same plan digest.

### 16.4 Budgets and cancellation

The compiled plan names portable limit and capability-policy requirements. After run input is
known, the host resolves concrete per-attempt grants for time, output bytes, telemetry, tool
calls, subject size, and containment; their exact provenance is recorded in admission and
episode records. The reusable plan contains no live grant or sandbox handle. Cancellation is
propagated to provider processes and tools. A component may return a typed partial result only
if its contract defines one; otherwise cancellation is a failure. Cleanup is attempted after
semantic completion; disposal failure emits an operational diagnostic/cleanup result and
leaves the already-derived semantic closure unchanged.

### 16.5 Effects and human gates

The cell emits decisions; an effector applies them. Operator policy remains required for
merges, releases, overrides, authority changes, credential or permission escalation, and other
material external effects. A valid receipt is necessary but not sufficient for an effect whose
policy requires a human gate.

---

## 17. Worked realizations

These examples test genericity. They must use the same normalizer, compiler, runner, kernel,
receipt, and verifier with no `cell_id` switch.

### 17.1 CDS implementation

| Axis | Realization |
|---|---|
| Contract | admitted implementation issue + distinct design |
| Subject | `git.snapshot/0.1` |
| Adapter | disposable repository/worktree -> `git.patch/0.1` |
| Producer | full engineering cognition with repository and shell capability |
| Methodology | code/test/language/style plus issue/design obligations |
| Mechanical providers | build, test, format, CUE/schema, path/static checks |
| Cognitive properties | AC satisfaction, design fidelity, impact completeness, simplicity |
| Matter | canonical patch against pinned base |
| Effect | separate merge/release path after delta and operator policy |

### 17.2 Writing

| Axis | Realization |
|---|---|
| Contract | admitted brief: purpose, audience, required claims/sources, constraints, output form |
| Subject | `document.bundle/0.1` with source material and existing document snapshot |
| Adapter | immutable source bundle + output slot -> `document.revision/0.1` |
| Producer | writing cognition; no Git or shell required by default |
| Methodology | general writing properties + selected style/voice |
| Mechanical providers | required sections, length, links/citations, forbidden terms, source refs |
| Cognitive properties | clarity, audience fit, fidelity, argument coherence, voice |
| Matter | `document.revision/0.1` |
| Effect | publishing/updating the destination remains outside the cell |

This example must not be implemented as merely another Git patch; otherwise it does not test
whether Git semantics leaked into the runner.

### 17.3 Planning

| Axis | Realization |
|---|---|
| Contract | admitted goal, scope, constraints, decision horizon, stop conditions |
| Subject | `program-state.snapshot/0.1`: issues, dependencies, receipts, resources |
| Adapter | typed state -> `plan.graph/0.1` |
| Producer | planning cognition over the state snapshot |
| Methodology | decomposition, dependency, authority, resource, risk, and proof properties |
| Mechanical providers | schema, unique ids, edge resolution, acyclicity where required, gate presence |
| Cognitive properties | completeness, sequencing, leverage, coherence, risk/scope fit |
| Matter | typed plan/relation graph |
| Effect | issue creation or dispatch occurs only after delta through effect adapters |

Together the examples require three distinct contract schemas, subject/matter pairs, and
capability surfaces while preserving the same cell mechanics.

---

## 18. Alternatives and decisions

| Question | Alternative | Decision and reason |
|---|---|---|
| Mechanical cells | prohibit them in the type | V0 deploys cognition-bearing cells but leaves the type open; uniform receipts may later earn mechanical cells |
| Contract | issue only | require issue + logically distinct design; implementation was otherwise choosing architecture locally |
| Assessment visibility | blind beta over patch text | reconstruct an independent evaluation view from contract subject + matter; independence is channel closure, not blindness |
| Skills | separate alpha/beta lists | one methodology authority with derived views; two lists drift |
| Coh ownership | Coh owns whole runtime | Coh owns methodology/measurement; CNOS retains admission, custody, CCNF, and effects |
| IR | execute CUE | canonical JSON IR; CUE validates structure but does not own operational semantics |
| IR | YAML | rejected for duplicate/coercion/canonical-byte hazards |
| IR | Protobuf/CBOR | useful possible transport later; poor bootstrap review/authoring surface |
| Construction | refs + service locator | explicit nested declarations and fixed constructor injection; one component, one place |
| Construction | Spring-like general container | borrow constructor injection only; reject reflection, scanning, mutable context, and ambient override |
| Provider execution | arbitrary Unix command in JSON | fixed typed adapter recipes over a common bounded process runner; prevents config-shaped RCE |
| Matter | provider self-report | runtime subject/matter adapter measures candidate state |
| Oracle placement | external CI only | methodology-declared providers inside beta composite, evidence bound for V |
| Admission refusal | fake episode closure | separate refusal receipt; no contract entered CCNF |
| Telemetry | embed all events in receipt | separate non-authoritative stream; promote selected observations through evidence producers |
| Scope discovery | auto-widen issue | stop and re-contract; a cell cannot authorize itself |

---

## 19. Leverage, cost, and limitations

### 19.1 Leverage

- Future CDS, writing, planning, and cohering work share one input/output/failure model.
- Coding begins from an admitted issue and design rather than discovering architecture mid-run.
- Coh can replace bootstrap projection and falsification without changing CCNF.
- Provider choice becomes replaceable execution detail rather than cell identity.
- Receipt custody, telemetry, and learning gain explicit channels.
- Review can test one authority instead of reconciling several essays and branch comments.

### 19.2 Negative leverage and process cost

- Visible Case-3 work pauses while design and schemas catch up.
- Existing `CellSpec`, closure, and branch Case-3 shapes will require migration.
- CDS adds an explicit design artifact and admission receipt.
- Exact source/IR/plan canonicalization must be versioned.
- Some skills mix authoring guidance, process procedure, and assessable properties and will need
  projection or decomposition.
- Full engineering capability requires a real containment decision, not merely CLI flags.

The cost is justified only if implementation stops re-deciding these boundaries locally.

### 19.3 Current limitations

- General Coh execution and assessed property packages do not yet exist.
- Cognitive admission is an attestation, not proof.
- No accepted cross-language canonical JSON ABI is selected yet.
- The current branch's `cds.review` and admission code are evidence/prototypes, not the final
  architecture.
- Current CLI failure paths do not all emit durable structured outcomes.
- No general subject/matter adapter registry or non-Git implementation exists.
- Provider-specific containment remains incomplete unless an outer substrate supplies it.

---

## 20. Migration and conformance

### 20.1 Authority disposition

| Existing artifact | Disposition |
|---|---|
| `COHERENCE-CELL-NORMAL-FORM.md` | retain as kernel authority |
| `COHERENCE-CELL.md` | retain as doctrine/rationale; reconcile its closure and repair terminology before this draft is ratified |
| `RECEIPT-VALIDATION.md` | retain as validation doctrine; map its predicates explicitly to the new protocol version |
| `schemas/cdd/receipt.cue`, `schemas/cdd/boundary_decision.cue`, `schemas/cds/receipt.cue` | retain as shipped-schema authority for the current runtime; the target `cnos.cell.ccnf/0.1-draft` is a distinct protocol until migration and parity witnesses prove replacement |
| `CELL-OF-CELLS.md` | retain as recursive system rationale |
| `DUMB-MODELS-SMART-CELLS.md` | retain as product/authority rationale |
| `CELL-RUNTIME.md` | reconcile its proposed runner realization against this design; it does not supersede this document |
| `CDS-CELL-MIGRATION.md` | retain as migration/history after this design is ratified |
| `CELL-RUNNER-CASES.md` | retain as executable-case record, not architecture authority |
| `CELL-SYSTEM-DESIGN-OUTLINE.md` | delete on promotion; git history preserves it |

### 20.2 Implementation impact surface

Ratification implies later changes across:

- source and normalized IR schemas;
- `cellspec`, `cellfill`, `cellfills`, and `cellrun`;
- `cellkernel` contract/receipt/outcome surfaces;
- CDS issue/design admission;
- subject/matter adapters;
- methodology loading and projection;
- bootstrap beta and property-check providers;
- evidence custody and telemetry;
- worked fixtures for CDS, writing, and planning;
- CLI/GitHub invocation adapters and operator-visible status.

This document does not prescribe the implementation order. A separate plan must trace every
changed value through all consumers before code resumes.

### 20.3 Conformance criteria

A runtime conforms to this design only when executable witnesses show:

1. one source lowers to normalized IR with no unresolved holes;
2. the same normalized inputs and installed dependencies yield byte-identical plan digests;
3. each component is declared once and constructed without runtime service lookup;
4. unsupported subject/matter pairs fail before cognition;
5. malformed or semantically rejected input emits an admission refusal receipt;
6. successful admission binds the exact normalized domain-contract and subject digests; the
   CDS domain contract in turn binds its distinct issue and design facets;
7. production and assessment receive the same contract and cannot share private mutable state;
8. matter is measured by the adapter against the pinned subject, not accepted from model prose;
9. bootstrap projections bind the same ordered complete skill-body digests, and later Coh
   projections have checkable stable property-id coverage equality. The bootstrap half is
   **vacuous by construction** and must be labelled so where it is implemented: when both
   views carry identical bodies under role wrappers, a digest-equality check cannot fail. It
   becomes a real witness only once projections can differ — that is, once property ids exist
   and a criterion can be absent from one view;
10. beta receives the CCNF pair `(contract, matter)` and returns a complete typed property map;
11. property-check pass, fail, unavailable, and mechanism fault each take a distinct tested
    path;
12. every evidence reference identifies producer authority and exact bytes;
13. gamma/V/delta are mechanically re-derived and reject mutation of contract, matter,
    assessment, evidence, or execution policy;
14. telemetry loss cannot change validation;
15. out-of-contract necessary work closes contract-unmet, dispatches repair with reason
    `needs_recontract`, and cannot self-widen;
16. CDS, writing, and planning examples use the same generic runner/kernel without id switches;
17. at least one negative fixture proves each barrier fails for its own reason;
18. current docs, schemas, runtime contract, help/status, and code describe the same shipped
    subset.

### 20.4 Deferred decisions and triggers

| Decision | Deferred until |
|---|---|
| exact cross-language canonical JSON algorithm | first non-Go compiler or stable IR ABI work |
| final `.cell`/Coh source syntax | normalized IR and compiler contract are stable |
| mechanical-only cell support | one use case shows receipts/composition outperform a function/workflow |
| station-specific provider/model holes | one cell demonstrably needs different configured engines |
| issue-declared additional skills | admission, authority, merge, and digest semantics are designed |
| beta with genuine filesystem tools over the reconstructed subject | an isolated read-only subject substrate exists and is demonstrated; until then bootstrap beta takes the bounded-value form of §11.4 |
| generalized provider registry | more than the bounded adapter switch creates real maintenance pressure |
| telemetry retention and epsilon learning schema | a learning cell and privacy/cost policy exist |
| full repair/main-loop orchestration | working, planning, and cohering cell contracts are independently implemented |

---

## 21. Design summary

The cell is a small, explicit system:

```text
one declared definition
one compiled methodology
one admitted contract and pinned subject
one independent matter channel
one complete assessment
one mechanical closure
one typed outcome
```

Its complexity is not in adding roles or framework machinery. It is in keeping meanings that
were previously smeared together distinct: source versus plan, contract versus subject, matter
versus evidence, evidence versus telemetry, capability versus containment, review versus
validation, and validity versus permission.

Once those boundaries are explicit, the implementation becomes ordinary constructor wiring,
typed values, pure transformations, bounded adapters, and executable witnesses. That is the
point of the design: cognition remains powerful inside the cell while authority stays simple,
auditable, and outside it.
