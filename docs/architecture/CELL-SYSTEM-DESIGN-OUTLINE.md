# Cell System Design — Outline

**Status:** Design outline · pre-ratification · not an implementation contract
**Purpose:** Establish the questions, types, boundaries, and proof obligations the
canonical cell design must settle before Case 3 or later cell work continues.
**Review mode:** Pi/Sigma design convergence. Review this artifact as architecture;
do not implement around an unresolved question.
**Current baseline:** `usurobor/cnos` main after the Case-2 merge.
**Companions:** CCNF owns the substrate-independent closure kernel; TSC owns
coherence-methodology semantics; this document will own their end-to-end CNOS cell
realization once promoted from outline to design.

> **This is deliberately an outline of the design, not the design itself.**
> Working positions are marked as such. Questions remain questions until Pi, Sigma,
> and the operator converge on them. JSON below is a candidate used to test whether
> the conceptual model is complete; it is not yet a frozen schema.

---

## 0. Why this outline exists

### Named incoherence

CNOS has a kernel equation, several strong conceptual papers, a generic episode
runner, and a CDS Case-2 realization, but it has no single current-state design
that names all cell components, their construction, their inputs and outputs,
their authority, and their failure paths.

Implementation is therefore discovering architecture locally. The same missing
design has appeared as several apparently separate defects:

- alpha was offered file editing without the engineering capability needed to run
  the tests that justify its work;
- beta was asked for a verdict from a diff that did not contain enough information
  to verify its own claims, and invented a concrete false defect;
- the issue was referred to in a one-line goal but was absent from the episode;
- admission, production, review, evidence, telemetry, and boundary policy have
  been added one at a time without a common type model;
- acceptance criteria were treated as if they enumerated every affected surface,
  leaving no rule for newly discovered impact outside those criteria;
- implementation and documentation repeatedly made stronger claims than the
  executed boundary established.

Sigma's current Case-3 r0 records the first two failures directly. The active Case-3
branch then began adding a typed issue and admission logic. Those changes are useful
evidence, but they are not a substitute for deciding the cell design first.

### Executed Case-3 evidence for the design-first pause

Sigma's design-first event 68 adds a stronger empirical claim: eight material defects
were all discovered after implementation even though each was decidable from a design.
The outline accepts that evidence and requires the final design to make each class
falsifiable before code resumes.

| Observed defect | Missing design proof |
|---|---|
| Diff measured from mutable `HEAD`, not the pinned base | complete value flow and consumer catalogue |
| Alpha withheld Bash and called that containment | enforce/declare/claim-nothing authority table |
| Shell was available but ordinary commands were not approved | separate capability availability, approval, and containment semantics |
| Beta emitted a confident false finding from insufficient input | per-seat input-sufficiency argument and an `unverified` outcome |
| New task matter was not traced into every closure/schema surface | cross-boundary value and scope-lift catalogue |
| An admission barrier leaked into unrelated specs and fixtures | barrier location and type-boundary specification |
| CUE and Go implemented different definitions of nonblank | rule owner, mirror, and executable parity witness |
| Local corpus ran a stale binary | oracle artifact provenance and proof-limit statement |

These are not eight independent patch requests. They are evidence that the final design
must expose the system's dataflow and authority before another Case-3 implementation
round.

### Challenged assumptions

This design challenges the following assumptions:

1. An issue or prose goal alone is enough to contract implementation.
2. A reviewer can reconstruct hidden context from candidate matter.
3. Withholding a capability constitutes containment.
4. Skills are merely prompt fragments rather than part of a methodology.
5. Acceptance criteria are a complete impact boundary.
6. Authored configuration, normalized IR, compiled plan, and run request can be one
   undifferentiated JSON object.
7. Another local implementation round can settle these questions more cheaply than
   a design.

### Design discipline

The final design must satisfy both current CNOS design disciplines:

- `cdd/design`: named incoherence, constraints, challenged assumption, impact
  graph, typed proposal, leverage, negative leverage, alternatives, and falsifiable
  acceptance criteria;
- `cnos.core/design`: one reason to change per boundary, policy above detail,
  truthful interfaces, explicit runtime surfaces, no duplicate source of truth,
  and no service-locator or registry drift.

---

## 1. Working thesis

### 1.1. Conceptual definition

**Working position:**

> A cell is a bounded, cognition-bearing transformation that applies a declared
> methodology to a contract and a subject, produces candidate matter, assesses it
> independently, and emits a validated receipt before any result crosses its
> boundary.

The cell is not defined by Claude, Codex, or another model. Its identity is its
contract type, subject type, output-matter type, methodology, and authority
boundary. A model is a replaceable execution engine inside a declared component.

### 1.2. Cognitive-only v0

**Working position:** every deployed CNOS cell in v0 contains at least one cognitive
station. A purely mechanical transformation should normally be an ordinary
function, workflow, oracle, or CM property provider; it does not need cell ceremony.

Do not permanently encode this product-scope decision as an ontological prohibition
in the kernel. A future mechanical-only plan may earn cell treatment if uniform
composition and receipts provide real leverage. That is not a v0 requirement.

### 1.3. Four distinct artifacts before a run

The design must distinguish:

```text
Issue             why/what must change; scope and acceptance
Design            how it will change; invariants and impact graph
SubjectSnapshot   exact state being operated on
CellDefinition    reusable methodology and component construction
```

For CDS, implementation must not start without both an issue and a design. A small
change may carry a compact inline design; the logical artifact remains distinct so
"what" cannot silently choose "how".

### 1.4. Three output planes plus the receipt

```text
matter       the proposed semantic result
evidence     runtime-observed witnesses supporting or refuting claims
telemetry    non-authoritative progress/debug/usage events
receipt      immutable closure envelope binding inputs, matter, assessment,
             evidence references, provenance, verdict, and decision
```

Telemetry is never evidence merely because it was logged. A named collector or
oracle must deliberately promote an observation into a digested evidence artifact.

---

## 2. Does the CDS cell now flatten?

### 2.1. Working answer

Yes, substantially—but not into one model call.

The cell-specific, configurable spine appears to be:

```text
admit contract + subject
  → produce candidate matter
  → assess candidate matter
```

The remaining tail is fixed mechanics:

```text
γ closes → V validates → δ decides
```

For the bootstrap:

- **admission** is a CDS-owned component combining structural validation with an
  optional cognitive issue/design review;
- **production** is a full engineering cognition component over a disposable,
  pinned subject;
- **assessment** is a fixed cognitive review component using the admitted contract,
  runtime-derived evaluation view, and declared skills;
- **gamma** is mechanical receipt construction;
- **V** is mechanical receipt validation;
- **delta** is boundary policy/authority, not cognition.

In the future, admission and assessment both become executions of compiled Coh
methodologies. Their position stays; their bootstrap implementation changes.

### 2.2. Important correction

Bootstrap beta is **not mechanical** merely because its orchestration and output
shape are fixed. It is rented cognition executing a fixed falsification procedure.
The eventual Coh evaluator is mechanically orchestrated, but individual property
providers may still be cognitive.

### 2.3. Where input checking belongs

The generic runner cannot know what makes a CDS issue, writing brief, research
question, or dataset admissible. Admission must therefore be a first-class
component of the cell definition.

The runner only performs the generic operation:

```text
construct the declared admission component
→ give it the run's contract + subject envelope
→ continue only on an admitted result
```

For CDS v0, a `cds.work-contract-admission` constructor owns:

1. exact structural validation of issue + design;
2. the structural relation between issue, design, and subject;
3. a cognitive semantic check, using issue/design skills, when structure alone
   cannot establish executability;
4. an explicit `admitted | rejected | incomplete | failed` result;
5. a digest binding the admitted issue/design bytes.

Malformed or unreviewable input emits an admission receipt and stops before alpha.
It is never silently discarded.

---

## 3. JSON question — source, IR, and run request

### 3.1. Working answer

We should now be able to express a CDS cell fully in JSON, but we need to stop
calling every JSON stage "the IR."

```text
CellSourceJSON       human/compiler-authored declaration; may contain holes
NormalizedCellIR     holes filled; exact keys/types; no unresolved symbols
CompiledCellPlan     skills loaded/digested; providers validated; capabilities and
                     execution plan bound by the composition root
RunRequest           admitted issue + design + pinned subject for one episode
```

The runner may fill holes while bootstrapping. The executable IR is the resolved,
canonical result, not the hole-bearing source.

### 3.2. Why canonical JSON remains the likely v0 IR

JSON is currently the strongest fit because it is:

- language-neutral between a future Coh/cell compiler and the Go runtime;
- directly validated by CUE as an independent oracle;
- content-addressable under a canonical serialization;
- inspectable and diffable during bootstrap;
- closed enough to reject unknown fields and arbitrary executable configuration;
- already the proven input/output shape of both CNOS cells and TSC CM work.

The final design must still compare:

| Alternative | Main cost |
|---|---|
| Execute CUE directly | conflates constraints/defaulting with operational semantics |
| YAML | duplicate keys, coercion, and weak canonical bytes |
| Protobuf | code-generation and poor review/authoring surface; useful later as transport |
| CBOR | compact canonical transport, poor human/debug surface |
| Compiler-private AST | not portable, independently validatable, or durable |

### 3.3. Candidate authored JSON

This candidate is a completeness probe, not a schema decision:

```json
{
  "format": "cnos.cell.source/0.1",
  "id": "cnos.cds.implement",
  "version": "0.1",
  "protocol": "cnos.cdd.cds.receipt.v1",
  "parameters": {
    "language": {
      "kind": "skill",
      "required": true,
      "domain": [
        "cnos.eng:eng/go",
        "cnos.eng:eng/ocaml",
        "cnos.eng:eng/typescript"
      ]
    },
    "style": {
      "kind": "skill",
      "required": false,
      "default": "cnos.eng:eng/write-functional"
    },
    "provider": {
      "kind": "value",
      "required": true,
      "domain": ["claude-cli", "fake"]
    },
    "model": {
      "kind": "value",
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
    },
    "admission": {
      "fill": "cds.work-contract-admission",
      "cognition": {
        "provider": "$provider",
        "model": "$model"
      },
      "skills": [
        "cnos.cdd:cdd/issue",
        "cnos.cdd:cdd/design"
      ]
    }
  },
  "produce": {
    "fill": "cds.patch",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    },
    "skills": [
      "cnos.eng:eng/code",
      "cnos.eng:eng/test",
      "$language",
      "$style"
    ]
  },
  "assess": {
    "fill": "cdd.review",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    },
    "skills": [
      "$language",
      "$style"
    ]
  },
  "output": {
    "matter": {
      "kind": "git.patch/0.1"
    },
    "required_evidence": [
      {
        "id": "candidate-diff",
        "kind": "git.diff",
        "producer": "runtime"
      },
      {
        "id": "assessment",
        "kind": "property-assessment",
        "producer": "assess"
      }
    ]
  }
}
```

### 3.4. What the candidate deliberately omits

- **Issue, design, repo, and base SHA are not holes.** They are run-specific
  contract/subject values carried by `RunRequest`.
- **Gamma, V, and delta are not configurable components here.** The declared
  protocol selects the fixed mechanical closure/validation/boundary contracts.
- **No arbitrary binary, argv, environment, command, or capability list appears.**
  A fill declares its required capability; runtime policy decides whether and how
  it can be provided.
- **No `bindings.alpha` or component-reference plane exists.** Each component is
  declared once with its constructor properties inline.
- **No task-specific skills hole appears yet.** The design must decide whether
  admitted issue/design contracts may declare additional canonical skills and how
  the fill merges them without creating an unbounded prompt-extension escape.

### 3.5. Which holes have earned their place?

**Working position:**

| Value | Classification | Why |
|---|---|---|
| language skill | semantic configuration hole | changes applicable engineering properties |
| style skill | project-policy hole with default | not universal to every CDS use |
| provider | execution-configuration hole | replaceable engine, selected before construction |
| model selector | execution-configuration hole | must be explicit in the resolved plan/receipt |
| issue + design | runtime contract, not holes | they define this episode's work |
| repo + base SHA | runtime subject, not holes | they define the exact state operated on |
| timeout/tool policy | runtime/fill policy, not holes yet | making them arbitrary weakens the boundary without demonstrated need |
| issue-specific skills | open decision | should be admitted contract data if allowed, not arbitrary CLI injection |

The same provider/model holes may populate admission, alpha, and beta while each
station still receives a fresh, independent invocation. We should not add
station-specific provider/model parameters until a real need appears.

### 3.6. Which CDS skills belong in alpha?

The known suitable baseline is currently in `cnos.eng`:

```text
eng/code
eng/test
eng/<language>
eng/write-functional (when project policy selects it)
```

The current root `cnos.cds:cds` skill describes lifecycle/orchestration concerns
that a workspace-edit alpha may not have authority or inputs to perform. Do not
inject it ceremonially. The design must decide whether CDS has a distinct,
capability-matched producer obligation deserving a narrow `cds/implement` skill,
or whether the fill + issue/design contract already supply all CDS-specific
semantics.

### 3.7. Bootstrap beta

The bootstrap intent is sound: beta is a fixed cognitive review procedure so CDS
can implement TSC, after which beta becomes Coh execution.

However, the current `cdd/review` skill cannot simply be assumed suitable. It
contains branch, issue/PR, `.cdd`, CI, and write/merge expectations. A no-workspace
answering seat cannot obey that full contract. The design must choose one of:

1. define a narrow bootstrap review skill over `(issue, design, subject, matter)`;
2. define a mechanically selected projection of the existing review skill;
3. give the review component the complete capability/input contract the existing
   skill requires.

The fixed review algorithm should be:

```text
for every declared property / acceptance obligation:
  attempt to falsify it against a runtime-derived evaluation view
  emit pass | finding | unverified | failed, with citations
```

It must not guess when the supplied view cannot decide a property.

### 3.8. Resolved IR obligations

After parameter filling and construction, the resolved/compiled representation
must bind at least:

- exact cell source digest and schema version;
- concrete provider adapter id and requested model selector for every cognitive
  station;
- ordered canonical skill refs plus exact content digests;
- fill implementation/version or policy digest;
- contract kind, subject kind, matter kind, and protocol;
- declared admission and assessment result vocabularies;
- capability/resource requirements as runtime-owned policy identifiers;
- no unresolved holes, symbolic component refs, arbitrary commands, or ambient
  configuration dependencies.

Whether this is one `NormalizedCellIR` or a normalized IR followed by a distinct
`CompiledCellPlan` is a design decision. TSC already argues for the two-stage split;
CNOS should not collapse it without evidence.

---

## 4. Architecture sketch the design must refine

```text
AUTHORING / CONFIGURATION
┌──────────┐  ┌──────────┐  ┌────────────────────────┐
│  Issue   │  │  Design  │  │ Cell source            │
│ what/why │  │ how/impact│ │ admission/α/assessment │
└────┬─────┘  └────┬─────┘  └───────────┬────────────┘
     └──────┬──────┘                    │
            │                    compile → CUE vet
            │                    → resolve/link
            └─────────────┬─────────────┘
                          ▼
                   Compiled Cell Plan
                          │
RUNTIME                   ▼
Subject snapshot → RunRequest → admission
                                │
                                ▼
                         α produces matter
                                │
                runtime reconstructs candidate view
                                │
                                ▼
                    β / Coh assesses properties
                                │
                     runtime-owned evidence
                                │
                                ▼
                        γ closes receipt
                                │
                        V validates receipt
                                │
                        δ decides boundary
                                │
                          effect adapter

OBSERVATION
runtime events ──────────→ telemetry sink       non-authoritative
closed receipts ─────────→ ε / learning cells   cross-episode
```

---

## 5. Proposed canonical design-document outline

The final design should answer the following questions in this order.

It must also include seven falsifiable inventories rather than leaving the answers only
in narrative prose:

1. component inventory — responsibility, ownership, and forbidden knowledge;
2. value catalogue — producer, consumers, schema surfaces, and scope-lift binding;
3. authority table — enforces, declares, and claims nothing, including barrier location;
4. seat specifications — decision, inputs, sufficiency argument, outcomes, independence;
5. rule ownership — normative owner, mirror, and executable agreement mechanism;
6. gate specification — artifact provenance, what green proves, and what it does not;
7. end-to-end flow — every gate in order, with typed failure and exit behavior.

The narrative sections below explain and justify those inventories; they do not replace
them.

### 5.1. Document contract and authority

- Is this the canonical current-state realization design?
- Which existing documents remain normative companions?
- Which become rationale, migration notes, or superseded proposals?
- What does CCNF own, what does TSC own, and what does CNOS cell realization own?

### 5.2. Problem, goals, non-goals, and challenged assumptions

- Which repeated failures prove the current boundary is incomplete?
- What future work should become cheaper after this design?
- What does v0 explicitly refuse to solve?

### 5.3. Conceptual definition and invariants

- What distinguishes a cell from a function, workflow, CM, provider, or agent?
- Is cognition required in v0 and/or in the type?
- What determines cell identity and episode identity?
- Which invariants must hold regardless of realization?

### 5.4. Vocabulary and type model

- What are `CellSource`, `NormalizedCellIR`, `CompiledCellPlan`, `RunRequest`,
  `Contract`, `SubjectSnapshot`, `Matter`, `Assessment`, `Evidence`, `Telemetry`,
  `Receipt`, `Verdict`, `BoundaryDecision`, and `ContractFinding`?
- Which are authored, derived, immutable, content-addressed, or runtime-local?

### 5.5. System boundary and complete component diagram

- Which components are inside the episode runtime?
- Are V and delta part of the cell wall or the enclosing parent?
- Where do effectors, storage, transport, worktrees, credentials, and model
  processes live?
- What is each component's one reason to change?

### 5.6. Definition, construction, and dependency injection

- Which properties are configured once?
- Which values are holes, run inputs, or runtime policy?
- How is every component constructed exactly once?
- What lifecycle/scope does each component have?
- Which Spring ideas are adopted and which are rejected?

Working Spring answer: explicit nested component definitions, constructor injection,
one composition root, per-episode scope, and validation before instantiation. Reject
reflection, scanning, autowiring, mutable singletons, service location, and arbitrary
runtime class/command names.

### 5.7. Source language, normalized IR, compiled plan, and CUE

- Is cell authoring a Coh `cell` form or a small composition form containing a CM?
- Which obligations are compiler checks versus CUE checks versus linker checks?
- Why JSON, and which canonicalization contract applies?
- How are schema/IR versions migrated?

### 5.8. Contract, subject, issue, and design

- What belongs in the generic contract?
- What makes an implementation issue executable?
- What must design add beyond the issue?
- May small designs be inline while remaining distinct?
- How are issue and design digests bound to the run?
- How do ACs, scope, non-goals, invariants, and the impact graph differ?

### 5.9. Admission

- What can be rejected structurally before cognition?
- What requires semantic assessment?
- Is semantic admission a preflight child cell or a station in the target cell?
- What exact artifact proves alpha received the admitted bytes?
- What happens on `INCOMPLETE`?

### 5.10. Episode execution semantics

- What exactly does alpha receive and what capabilities must it have?
- How is the subject contained and how is candidate state measured?
- What exactly does beta/Coh receive?
- Which stages may run in parallel?
- When do oracles run?
- What does cancellation mean at each stage?

### 5.11. Methodology, skills, properties, and Coh

- What is the one normative source of assessable properties?
- How are constructive alpha guidance and adversarial assessment projected from it?
- Which instructions are advice rather than assessable obligations?
- How are mechanical and cognitive property providers composed?
- How are refusal and uncertainty represented?

### 5.12. Outputs, evidence, receipt, decision, and custody

- What returns directly and what is stored by reference?
- Does a receipt embed matter or bind its digest?
- Who owns evidence bytes?
- Which artifacts ascend to the parent cell?
- How do retries and superseding attempts preserve history?

### 5.13. Telemetry and learning

- Which progress/tool/usage events exist?
- What is retained or redacted?
- When can telemetry become evidence?
- What may epsilon consume?
- Does losing telemetry ever affect the verdict? Recommended: no.

### 5.14. Failure, discovery, repair, and re-contracting

- How do malformed input, provider failure, matter defect, unverified property,
  contract defect, resource exhaustion, and cancellation differ?
- What happens when alpha discovers a necessary surface outside the design's impact
  graph?
- Can a cell amend its own contract? Recommended: no.
- When is repair valid and when is a new issue/design digest required?

### 5.15. Composition and recursion

- How do child contracts descend and receipts ascend?
- How does alpha propose children without invoking the runtime itself?
- How do `let!` and `and!` lower into a dataflow graph?
- What belongs to the main loop rather than one cell?

### 5.16. Authority, capability, security, and resources

- What is declared, available, approved, and actually contained?
- Which substrate gives engineering alpha full useful capability safely?
- Who chooses provider/model?
- Which credentials may enter execution but never the receipt?
- Which boundary actions require an operator?

### 5.17. Reliability and operational semantics

- What is deterministic and what is not?
- What does reproducibility mean for cognition?
- How are timeouts, retries, IDs, cancellation, concurrency, and canonical event
  order defined?
- Does retry mint a new attempt? Recommended: yes.

### 5.18. Worked realizations

- CDS implementation: admitted issue+design + pinned repo → patch → assessment.
- Writing: admitted brief/style + source text → document.
- Planning: goal + state snapshot → plan/relation matter.
- One mechanical and one cognitive property provider inside a cognitive cell.

### 5.19. Alternatives and decision record

Compare at minimum:

- cognitive-only versus mechanically realizable cells;
- issue-only versus issue+design;
- blind beta versus runtime-reconstructed evaluation view;
- skills-as-prompts versus properties-as-methodology;
- Coh-as-whole-cell versus Coh-as-methodology component;
- JSON versus CUE/YAML/Protobuf/CBOR/private AST;
- nested constructor injection versus component refs/service locator;
- telemetry inside receipts versus separate stream;
- automatic scope widening versus explicit re-contracting.

### 5.20. Migration, impact graph, and conformance

- Which Case-2/Case-3 code survives, changes, or is held?
- What is the smallest bootstrap beta compatible with future Coh?
- When does work switch to TSC?
- Which existing docs are updated or demoted?
- Which fixtures prove every boundary can fail for its own reason?

---

## 6. Contract-discovery rule

Acceptance criteria prove selected outcomes. They do not enumerate every effect a
valid implementation may have.

The design's impact graph, invariants, scope, and non-goals establish the execution
boundary. During production:

| Discovery | Required disposition |
|---|---|
| Necessary change already named by impact graph/invariants | continue and bind evidence |
| Necessary dependency not represented in design | emit `ContractFinding`; stop `needs_recontract` |
| Issue/design contradiction | emit `ContractFinding`; stop `needs_recontract` |
| Optional improvement outside contract | do not implement; emit follow-up candidate |
| Required oracle unavailable | assessment `INCOMPLETE`; no acceptance |
| Forbidden/out-of-scope effect | reject candidate matter |

Partial matter may be retained as a non-transmissible proposal. The cell cannot
amend the issue/design that authorizes itself. An enclosing planning/design cell or
operator produces a new contract digest, then a new episode runs.

---

## 7. Impact graph for the eventual design

### Upstream sources

- CCNF kernel doctrine and receipt-validation boundary;
- `cdd/issue`, `cdd/issue/contract`, `cdd/issue/proof`, and `cdd/design`;
- `cnos.core/design` and `cnos.eng/evolve`;
- TSC `CMSource → NormalizedCMIR → CompiledCM → RunRequest → MeasurementReceipt`;
- Case-2 executable evidence and Sigma's Case-3 r0 findings.

### Downstream consumers

- `schemas/cdd/spec.cue` and `schemas/cds/spec.cue`;
- `cellfill`, `cellspec`, `cellrun`, `cellkernel`, and composition root;
- CDS admission, patch, and bootstrap review fills;
- skill loading and resolved invocation metadata;
- workspace/evidence/telemetry/provider adapters;
- future `.cell`/Coh compiler and TSC runtime;
- CLI and future GitHub invocation adapter;
- planning/cohering/main-loop work;
- receipt custody, epsilon, and memory/learning projections.

### Existing documents requiring an authority disposition

- `COHERENCE-CELL-NORMAL-FORM.md` — retain as kernel authority;
- `CELL-OF-CELLS.md` — retain as recursive system model/rationale;
- `DUMB-MODELS-SMART-CELLS.md` — retain as product/authority rationale;
- `CELL-RUNTIME.md` — reconcile or supersede its proposed realization claims;
- `CDS-CELL-MIGRATION.md` — demote to migration/history after canonical design;
- `CELL-RUNNER-CASES.md` — retain as executable-case record, not architecture source.

---

## 8. Leverage and negative leverage

### Positive leverage

- Stops implementation from choosing architecture one local defect at a time.
- Makes issue and design explicit prerequisites to coding.
- Gives every future cell one construction, admission, evidence, and failure model.
- Lets bootstrap skills be replaced by Coh without changing episode boundaries.
- Makes planning, working, and cohering cells composable through the same artifacts.
- Gives Sigma and Pi one review authority instead of several partially overlapping
  essays and branch-local comments.

### Negative leverage / cost

- Pauses visible Case-3 implementation progress.
- Requires a migration from current CellSpec/closure shapes.
- Adds an explicit design artifact to CDS contracting.
- Requires deciding versioned source/IR/compiled/run artifact boundaries.
- May reveal that some current skills mix incompatible roles and need projection or
  decomposition.

The cost is justified only if the final design deletes repeated local reasoning and
becomes the contract future implementation reviews against.

---

## 9. Non-goals of this outline

- No Case-3 code changes.
- No final JSON/CUE schema.
- No new issue or PR.
- No Coh runtime implementation.
- No generic provider router.
- No repair loop or main-loop implementation.
- No mechanical-only cell implementation.
- No claim that the candidate JSON is accepted by current code.
- No claim that the current `cdd/review` skill is usable unchanged.

---

## 10. Questions for Sigma's design review

Sigma should review the outline and candidate JSON against the live Case-2/Case-3
implementation evidence, then answer each item with `agree`, `change`, or `open`:

1. Is the three-component configurable spine `admit → produce → assess` complete?
2. Is the fixed tail `γ → V → δ` correctly omitted from per-cell configuration?
3. Are issue + design + subject + cell definition the complete pre-run artifact set?
4. Does the Source JSON contain every CDS-specific decision without leaking runner,
   GitHub, arbitrary command, or service-locator semantics?
5. Which holes have actually earned their place? Are any missing?
6. Does `cnos.cds` need a narrow producer skill, or is CDS semantics already supplied
   by the fill + issue/design contract?
7. Can the existing review skill be projected safely, or does bootstrap require a
   new narrow review skill?
8. What must beta's runtime-derived evaluation view contain while preserving the
   CCNF `(contract, matter)` boundary?
9. Should admission be a preflight child cell or an inline station of the run?
10. Should Coh express the whole cell, or remain a measure-only component embedded
    in a cell composition form?
11. Is `NormalizedCellIR` distinct from `CompiledCellPlan`?
12. Which conclusions are contradicted by actual executed evidence?

The requested output is a design response or patch to this outline. It is explicitly
not authorization to continue Case-3 implementation.

---

## 11. Acceptance criteria for promoting outline → design

- [ ] The named incoherence and all observed failure modes are evidenced.
- [ ] Every component appears once in the diagram and responsibility table.
- [ ] Contract, subject, definition, source, normalized IR, compiled plan, run,
      matter, assessment, evidence, telemetry, receipt, verdict, decision, and effect
      are distinct typed artifacts or intentionally unified with a stated proof.
- [ ] The CDS cell is fully expressible as valid candidate JSON without hidden
      semantics in the generic runner.
- [ ] Issue and design admission has structural and semantic paths with typed refusal.
- [ ] Alpha has truthful engineering capability under a named containment boundary.
- [ ] Beta can verify every asked property or return `INCOMPLETE`; no hidden-context
      recovery is required.
- [ ] Gamma is demonstrably mechanical; V and delta authority remain separate.
- [ ] Matter, evidence, telemetry, and receipt channels cannot be confused.
- [ ] Newly discovered impact has a typed re-contract path and cannot self-authorize.
- [ ] Coh/TSC migration preserves the same positions and artifact contracts.
- [ ] Alternatives include at least three structurally different composition/IR moves.
- [ ] The impact graph covers every existing doctrine, schema, runtime, skill, and
      invocation consumer.
- [ ] Every final implementation AC names a specific artifact and an oracle, including
      a negative/mutation witness.
