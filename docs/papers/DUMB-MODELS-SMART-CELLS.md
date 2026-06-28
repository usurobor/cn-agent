---
title: "Dumb Models, Smart Cells"
subtitle: "Receipts, Trust Claims, and Coherence Witnesses for Agent Systems"
version: v0.11.0
status: DRAFT
author: usurobor (aka Axiom) (human & AI)
date: 2026-06-26
---

# Dumb Models, Smart Cells
## Receipts, Trust Claims, and Coherence Witnesses for Agent Systems

**Status:** v0.11.0 (DRAFT — position paper)
**Author(s):** usurobor (aka Axiom) (human & AI)
**Date:** 2026-06-26

> **Scope:** This paper explains how cnos keeps releasable agent work outside model runtimes. The cnos answer is a *work body* — cells, packages, receipts, coherence reports, and repo state — carrying explicit trust claims, coherence witnesses, and named boundary authority.

---

## 0. The claim

The hard question for agent systems starts after the model answers.
Who owns the resulting work?

Models supply capability. They write code, inspect files, draft plans, summarize evidence, and propose changes. Frontier models hold more context, plan across longer horizons, and coordinate more steps than smaller local models.

cnos rents that capability.
The durable work remains outside the model, in a work body the system owns (§5).

Within that body:
the cell owns workflow,
the package owns local cognition,
the receipt owns evidence,
the coherence report owns measurement,
the repo owns continuity,
and the boundary authority decides what crosses into release.

Two claims ride on every receipt that matters. The trust claim names the strength of attribution. The coherence witness asks whether the artifacts still describe one unit of work.

Capability is rented.
Authority is owned.
The receipt must not claim more trust than its artifact carries.
The system must not claim more coherence than it measured.

---

## 1. Capability attracts workflow

Vendor runtimes are useful.

A hosted model runtime can hold context, remember files, call tools, run workflows, plan across a session, and hide a lot of friction. When the model and workflow live together, the user sees less glue code, fewer seams, and more apparent continuity.

That value is real.
The same value creates product gravity.

A prompt becomes a thread. A thread gets memory. Memory gets files. Files get tools. Tools get workflows. Workflows get hosted agents. Hosted agents get logs, permissions, background tasks, connectors, deployment surfaces, registries, identities, gateways, observability, evaluation, and governance.

Each feature helps.
Together they make the vendor runtime the place where work lives.

That shape has already arrived. Google's Gemini Enterprise Agent Platform puts agent development, runtime, registry, identity, gateway, observability, governance, and long-running state in one platform. Memory Bank gives agents long-term memories across sessions. OpenAI exposes file search over vector stores. Anthropic is building managed agent infrastructure around Claude.

The risk is structural.
The easiest place to put the next feature is also the easiest place to lose authority.

At first the model answers.
Then it remembers.
Then it plans.
Then it acts.
Then it explains what it did.
Then the explanation, memory, workflow, logs, and decision boundary all live inside someone else's product.

That may be the right trade for disposable work.
It is a poor system of record for work that must be reviewed, released, resumed, audited, or defended later.

---

## 2. The boundary earns its cost

A bounded model costs something.

Keeping workflow outside the vendor runtime gives up some convenience. Long sessions may feel less smooth. Vendor-native planning may be harder to use. Hidden context carried by the product must become explicit. The system gains ceremony: cells, receipts, validators, routing policy, release gates, package doctrine, trust claims, and coherence reports.

The trade is capability for inspectability.

Some work does not need that trade.

For exploratory, private, cheap, throwaway work, a hosted agent runtime may be enough. Let the model run. Let the thread hold the context. Let the product own the workflow.

cnos starts to earn its keep when chat history is too weak to be the system of record.

Use a stronger boundary when work:

- changes durable repo state,
- affects other people,
- crosses a release boundary,
- needs review,
- spans more than one session,
- spans more than one agent,
- needs audit later,
- can fail in a way people care about.

Those cases need a place where authority lands.
They also need a way to ask whether the artifacts still describe the same work.

---

## 3. Dumb means bounded

"Dumb model" means bounded model.

The model may be strong inside the call. It may write code, inspect evidence, propose a repair, draft release notes, or explain a tradeoff.

It stays non-authoritative.

The model produces matter for a cell.
The runtime chooses the provider route.
The validator checks the receipt.
The coherence witness measures the work body.
The boundary authority decides what crosses.

The model performs a move inside a system-owned boundary.

That boundary is the difference between a capable model and a governed work system.

---

## 4. Providers are engines, not sovereigns

cnos draws the model boundary in the provider contract.

A provider is an executable capability endpoint. It may be `cnos.llm.local`, `cnos.llm.anthropic`, `cnos.net.http`, or `cnos.transport.git`.

The provider declares what it can do.
The kernel decides how it may act.

The provider contract states the rule:

> Package distributes. Extension declares. Provider executes. Kernel decides.

An LLM backend sits in the provider slot. It receives a typed operation, effective permissions, limits, and context. It returns a typed result, artifacts, timing, or an error.

The provider may be powerful.
The provider is governed.

Model choice belongs in runtime policy. A skill may shape the task. A command may declare a hint. The runtime chooses the provider and records the route.

Every routed model call should leave a route receipt:

- route decision,
- reason set,
- provider,
- model name,
- fallback status,
- latency,
- token estimate,
- validation signal when available.

Models remain unequal. A small local model, a code-specialized model, a frontier model, and a long-context model have different capability. Swapping them can change whether the work succeeds.

cnos does not promise model fungibility.
It gives the system a governed seam where model choice is explicit, recorded, and replaceable without moving workflow, evidence, trust, coherence, memory, or release authority into the vendor runtime.

---

## 5. The work body

A model output becomes work only when it enters a body.

In cnos, that body is made of five parts: cells, packages, receipts, coherence reports, and repo state.

### Cells carry workflow

A cell is a bounded unit of work. The CDD kernel names the loop:

`contract -> matter -> review -> receipt -> verdict -> decision`

The kernel names roles, artifacts, validator `V`, evidence, verdicts, and decisions. It does not depend on GitHub, Claude, prompts, CI, branches, or any particular invocation surface.

A model may produce matter for a cell.
The cell gives the work its shape.

### Packages carry local cognition

Skills, doctrine, runtime adapters, and domain rules arrive as packages.

An agent wakes into local cognition instead of depending on a vendor's hidden system prompt to know how a repo works.

Packages can be versioned, reviewed, replaced, and composed.

### Receipts carry evidence

A receipt is the closeout artifact.

It records what happened, what evidence was bound, what validator ran, what trust claim was carried, what coherence witness exists, what decision was made, and what can cross the boundary.

A receipt gives another human, agent, validator, or release step something to inspect.

### Coherence reports carry measurement

A coherence report asks whether the work body still hangs together.

Do the contract, matter, review, evidence, validator verdict, delta decision, release note, and handoff still describe one unit of work?

That is the TSC question. The receipt references the report; the report is its own artifact.

### Repo state carries continuity

The repo is the durable software body.

Git gives history, signatures, diffs, branches, tags, merges, clones, and offline operation. Git is not the source of coherence. It is the lowest durable surface where coherence can be inspected.

The chat can end.
The repo remembers.

---

## 6. A trace

Here is the shape in miniature.

An operator asks Sigma to turn a new repo into a Hello World app. Sigma creates an issue contract first.

```yaml
issue: 1
title: Create minimal Hello World app
contract:
  requirements:
    - create the smallest runnable app appropriate for this repo
    - add run instructions
    - add one validation step
    - prepare release notes
  acceptance:
    - app runs
    - validation evidence is recorded
    - release boundary is explicit
    - closeout receipt exists
    - coherence measurement exists or absence is recorded
status: ready
```

The human reviews the issue and authorizes dispatch.

```yaml
dispatch:
  issue: 1
  protocol: cds
  status: todo
  authority:
    subject: peter@operator
    trust_claim:
      governance_scope: G1_repo_commons
      evidence_strength: E0_unsigned
      process_integrity: P1_committed_history
      admissibility: A0_none
    signature: null
    attestation: none
  authorized_at: 2026-06-22T14:31:00Z
```

`protocol: cds` selects the concrete software-development protocol. The cell still runs the generic CDD kernel internally.

The wake claims the issue and routes work.

```yaml
route_receipt:
  cell: issue-1
  provider: cnos.llm.anthropic
  model: claude-code
  reason_set:
    - implementation touches repo structure
    - release notes required
    - validation evidence required
    - coherence report required before non-degraded release
  fallback: none
  latency_ms: 1842
  token_estimate: 18231
```

The cell produces matter on dev, binds evidence, and closes with a receipt.

This trace is illustrative. It records missing coherence measurement instead of inventing a score.

```yaml
cell_receipt:
  cell: issue-1
  contract: Create minimal Hello World app
  matter:
    - app source added
    - README run instructions added
    - validation command recorded
    - release notes prepared
  evidence:
    - git diff
    - validation output
    - release notes path
    - target commit
  trust_claim:
    governance_scope: G1_repo_commons
    evidence_strength: E0_unsigned
    process_integrity: P1_committed_history
    admissibility: A0_none
  coherence:
    required: true
    measured: false
    status: NOT_MEASURED
    expected_report: .tsc/cells/issue-1/report.json
    reason:
      - TSC coherence measurement was not run in this illustrative trace
  validator:
    name: cdd.v.basic-release
    verdict: FAIL
    witnesses:
      - validation_command_recorded
      - release_notes_present
      - target_commit_present
    reason:
      - validation command recorded but not executed in clean CI
      - coherence report missing
  delta:
    authority:
      subject: peter@operator
      trust_claim:
        governance_scope: G1_repo_commons
        evidence_strength: E0_unsigned
        process_integrity: P1_committed_history
        admissibility: A0_none
      signature: null
      attestation: none
    decision: override
    reason:
      - demo repo
      - no production users
      - validation weakness recorded
      - coherence measurement missing and recorded
  outcome: degraded
  output:
    branch: deploy
    tag: v0.1.0
    tag_signature: null
```

This receipt proves very little by itself.

It gives the decision a body.

A later agent can see what ran, which model ran, what evidence was bound, what trust claim was carried, what validation failed, what coherence measurement was missing, who accepted the degradation, and what tag crossed the boundary.

The live Hello World demo should replace this illustrative trace with a captured receipt and a real TSC report.

---

## 7. The same-work problem

TSC solves the same-work problem.

After agent work moves through an issue, a cell, a commit, a validator, a release note, and a handoff, the system needs one more question:

Do these artifacts still describe one unit of work?

A receipt can be named, signed, attested, retained, and incoherent.

A test can pass while the release note describes the wrong change.
A commit can be signed while the evidence points at the wrong commit.
A validator can emit FAIL while the handoff says "ready for production."
A delta override can record a demo exception while the next cell treats the output as clean.

Every artifact exists.
The body has drifted.

TSC measures that drift.

Trust says who stands behind a decision.
Validation says which checks passed or failed.
TSC asks whether the artifacts still describe the same work.

---

## 8. TSC as witness

A TSC report is a witness to `V`.

The CDD loop already has the right places:

`contract -> matter -> review -> receipt -> verdict -> decision`

The receipt binds evidence.
`V` checks the receipt against the contract.
`delta` decides what crosses the boundary.

TSC adds a coherence witness to the validator's evidence set.

For a cnos cell, the witness asks three questions:

```text
pattern:
  Do contract, matter, receipt, and release notes use one stable structure and vocabulary?
relation:
  Do evidence, validator verdict, delta decision, commit, tag, and release note point to the same work?
process:
  Do status transitions, handoff, degradation, and downstream propagation preserve the cell over time?
```

The bottleneck gives repair direction.

Low pattern coherence points to schema, vocabulary, or artifact-shape repair.
Low relation coherence points to reference repair.
Low process coherence points to handoff, state-transition, or propagation repair.

Mechanical mode gives the first useful witness. It catches schema drift, missing files, broken references, inconsistent paths, missing reports, and obvious cross-artifact mismatch. It is deterministic and CI-friendly, and it runs without a model.

Mechanical mode is a structural witness.

Hybrid mode adds LLM-backed semantic judgment under the same boundary. The model returns as a witness, not as authority. To keep the witness independent of the author, the coherence provider should be a different route than the model that produced the matter. When the same route serves both, the work is not independently witnessed, and the receipt must record that the witness and the author were not independent. The receipt records mode, provider, report path, and result in every case.

The rule is the same as trust:

> Do not claim more coherence than the measurement mode supports.

TSC does not decide release.
It gives `V` and `delta` a clearer surface.

---

## 9. Coherence and correctness

Coherence gives work a condition for judgment.
It does not certify truth.

A coherent cell can still contain a broken app.
A coherent plan can still be the wrong plan.
A coherent research claim can still overstate its evidence.

TSC does not replace tests, type checks, security review, product judgment, legal review, signature verification, gateway logs, runtime policy, or human responsibility.

It measures whether the work body is coherent enough to be judged.

The human, or a policy a human owns, still decides:

- what belongs to the cell bundle,
- which artifacts are canonical,
- which TSC mode is enough,
- which threshold applies,
- whether a threshold is a baseline or policy override,
- whether FAIL blocks release,
- whether FAIL_DEGENERATE invalidates the measurement,
- whether low pattern, relation, or process coherence requires repair,
- whether an override is acceptable,
- who holds delta,
- who answers for the risk.

The score is a witness.
The human still owns the boundary.

---

## 10. Trust has axes

The Agent Audit Trail Internet-Draft moves agent logging in the right direction. It defines JSON audit records with agent identity, action classification, outcome tracking, trust-level reporting, SHA-256 hash chaining, and optional ECDSA signatures.

That machinery is useful.

The flat trust ladder is too coarse for cnos receipts.

The draft's L0 to L4 field combines different claims into one line. The first step from a bare label to a named repo authority may change the social agreement while leaving the bytes unchanged. The first cryptographic step changes the artifact itself. A signature appears. A receipt, commit, tag, or release object now carries tamper-evident material.

Those claims need separate axes.

```yaml
trust_claim:
  governance_scope: G1_repo_commons     # who vouches
  evidence_strength: E0_unsigned        # what the artifact proves
  process_integrity: P1_committed_history
  admissibility: A0_none
```

`governance_scope` names the commons or authority that accepts the subject.
`evidence_strength` names the artifact's cryptographic strength.
`process_integrity` names the strength of ordering, retention, and chain reconstruction.
`admissibility` names any policy, compliance, or legal regime the record targets.

A repo-local named decision can be useful with `E0_unsigned`.
A signed artifact can still have weak governance.
A hash-chained log can still lack legal admissibility.
A compliance record can still carry a weak coherence witness.

The receipt should state the claim directly.

```yaml
trust_axes:
  governance_scope:
    G0_local_label: a string names an actor
    G1_repo_commons: the repo commons attributes the decision to a subject
    G2_org_identity: an organization identity system vouches for the subject
    G3_external_authority: an external witness or authority vouches for the event
  evidence_strength:
    E0_unsigned: no cryptographic proof is attached
    E1_signed_artifact: a key signed the receipt, commit, tag, or release artifact
    E2_hash_chained_log: records are linked by hashes
    E3_external_attestation: an external witness or timestamping service attested to the event
  process_integrity:
    P0_unordered_log: events are recorded with weak ordering
    P1_committed_history: events are recorded in repository history
    P2_signed_chain: event sequence is signed or hash chained
    P3_retained_chain: records are retained under explicit policy
  admissibility:
    A0_none: no compliance or legal claim is made
    A1_policy_record: the record targets an internal policy
    A2_compliance_record: the record targets a compliance regime
    A3_legal_evidence: the record targets legal admissibility
```

The example trace in this paper is:

```yaml
trust_claim:
  governance_scope: G1_repo_commons
  evidence_strength: E0_unsigned
  process_integrity: P1_committed_history
  admissibility: A0_none
```

The decision is named inside the repo.
It is unsigned.
It is not externally attested.
It makes no compliance claim.

That is acceptable because the receipt says so.

The rule is simple:

> Do not claim more trust than the artifact carries.

---

## 11. Named authority

Push the question down far enough and it bottoms out at judgment.

The model is bounded.
The validator may be weak.
The coherence report may be mechanical.
The receipt may be unsigned.

`delta` still has to decide what crosses the boundary.

cnos places that decision at a named seam.

In CDD, `V` checks the receipt against the contract. `delta` decides whether to accept, reject, repair-dispatch, or release with degradation recorded.

A degraded release becomes a recorded decision. If the receipt carries a signature, the decision is signed. If it carries third-party attestation, the decision has an external witness. If it carries a TSC report, the decision has a coherence witness.

These are separate claims.

A degraded cell also propagates. It becomes degraded matter for the next cell. A later boundary can refuse to build on it, ask for repair, or accept the risk again under another name.

Patterns become visible. If overrides become routine, `epsilon` — the system's pattern-reader over the receipt stream — can surface that as system behavior rather than one-off noise.

The system still depends on a commons that cares.

A negligent community can rubber-stamp degraded work forever. cnos will record that too.

The guarantee is structural: not-caring leaves a trail.

Fail-closed versus degrade-and-record follows the same rule. It is a policy decision made at a named seam, by an authority responsible for the cost of stopping and the risk of shipping.

---

## 12. Neighboring seams

The choice is larger than cnos versus vendor runtime.

There is already a middle ecosystem that keeps pieces of agent systems outside any one model vendor.

Stateful-agent systems preserve memory and context across conversations.
Dedicated memory layers make memory portable across orchestrators.
Tool protocols standardize how models reach external systems.
Agent-to-agent protocols standardize how agents discover each other and exchange work.

Those seams matter.

Letta is a stateful-agent system.
Mem0 is a memory layer.
MCP is a tool and context seam.
A2A is an agent-to-agent communication seam.

cnos closes a different boundary.

Memory says what the agent may know.
Tools say what the agent may touch.
Agent protocols say who the agent may talk to.
A cnos receipt says what happened, what evidence was bound, what failed validation, what coherence was measured, who accepted the risk, what trust claim the decision carries, and what crossed the boundary.

cnos belongs at the boundary of releasable work.

---

## 13. Trust Gates and cell receipts

The closest neighbor to cnos receipts is the Trust Gate.

The Agent Identity Framework Internet-Draft defines a Trust Gate as middleware before business logic. It verifies an agent passport, checks trust level, evaluates policy and context, returns ALLOW or DENY, and logs its reasoning.

That is the right shape for per-operation governance.

A Trust Gate answers:

- may this agent perform this operation now?
- what identity and trust level are attached?
- what policy applied?
- was the operation allowed or denied?
- what audit record should be written?

A cnos cell asks a larger question:

- what was the contract?
- what matter was produced?
- what evidence was bound?
- what review happened?
- what did the validator decide?
- what did the coherence witness say?
- who accepted the risk?
- what crossed the release boundary?
- what should the next cell inherit?

A Trust Gate records runtime decisions.
A cnos receipt closes a unit of releasable work.

They compose cleanly. Trust Gate logs can become evidence inside a cnos receipt.

They should not be mistaken for the receipt.

"Every tool call was logged" is a per-operation claim.
"This unit of work passed, or failed and was overridden by a named authority, with coherence measured and degradation recorded" is a per-cell claim.

Releasable work is larger than a call.

---

## 14. The Rails echo

Rails is an echo, not the proof.

The analogy carries one useful memory: a powerful engine does not have to become the place where the application lives.

Databases in the early Rails era were mature and feature-rich. They stored data, ran queries, enforced constraints, and offered stored procedures, triggers, functions, and vendor-specific behavior.

One camp put more logic in the database.
Another kept application logic in code behind framework-owned conventions.

Rails made the second bet. Active Record and migrations gave developers a stable seam. The database stayed powerful, while the application stopped living inside stored procedures and vendor-specific surfaces.

LLMs are harder.

There is no stable LLM capability contract like the database adapters eventually had. A frontier model and a small local model are not interchangeable like two SQL-backed adapters can sometimes be.

The echo still holds.

A powerful engine can sit behind a convention.
The convention is what the developer keeps.

---

## 15. Why this is an open-source shape

Open source wins by making strong seams normal.

For agent systems, that means:

> the repo owns state,
> the runtime owns effects,
> the protocol owns handoff,
> the cell owns workflow,
> the receipt owns evidence,
> the coherence report owns measurement,
> the package owns local cognition,
> the model owns only the bounded move it was asked to make.

A community can add a package.
A team can swap a model.
A domain can define its own evidence.
A validator can check a receipt.
A TSC report can measure coherence.
A Trust Gate log can become evidence.
A future agent can continue from repo state instead of guessing from chat history.

The model is replaceable at the workflow boundary because the work has a body outside the model.

---

## 16. What the boundary gives

The boundary gives substitution.
A team can add a provider without moving workflow into that provider. The new provider may be better or worse. The receipt records which one ran, why it was chosen, and what validation followed.

The boundary gives audit.
A receipt records what the system believed, what evidence was bound, what validator ran, what coherence was measured, and what decision crossed the boundary.

The boundary gives explicit trust.
A named decision, a signed decision, an externally attested decision, and a compliance record are different claims. The receipt can say which one it carries.

The boundary gives same-work detection.
A signed release can still be incoherent. A TSC report gives the validator a witness for whether the artifacts still describe one unit of work.

The boundary gives repair direction.
Low pattern coherence suggests schema or vocabulary repair. Low relation coherence suggests reference repair. Low process coherence suggests handoff, transition, or propagation repair.

The boundary gives durability.
The work survives a model call, a vendor outage, a price change, a lost chat, or a new runner. The repo keeps the body.

The boundary gives cost control.
Simple work can stay local. High-context work can escalate. Failed local work can repair-dispatch or route up. The decision is policy, not habit.

The boundary gives audit readiness.
The EU AI Act's Article 12 requires high-risk AI systems to technically allow automatic event logging over the lifetime of the system. NIST and OWASP governance work point in the same direction: systems need records, monitoring, risk management, and evidence. A cnos receipt is not a legal compliance program. It is the kind of body a compliance program can inspect.

The pattern is simple.
Keep durable concerns where the system can reach them.

---

## 17. Relationship to CN, THESIS, CDD, and TSC

The CN whitepaper answers a substrate question:

> Where do agents coordinate?

Answer: Git, with a thin CN convention layer.

`THESIS.md` answers the system question:

> What is cnos?

Answer: a recurrent coherence system whose articulations include doctrine, documents, packages, runtime modules, repos, traces, releases, and agents.

CDD answers the cell question:

> How does bounded work close?

Answer: contract, matter, review, receipt, verdict, decision.

TSC answers the measurement question:

> Do independent descriptions of the same system still describe one system?

Answer: measure pattern, relation, and process coherence.

This paper answers a placement question:

> Which parts of agent work belong outside the vendor model?

Answer: the work body — cells, packages, receipts, coherence reports, and repo state — carrying trust claims, coherence witnesses, and release authority.

The papers draw one boundary at different layers.

CN draws it under transport: forges are projections; the repo is the substrate.
`THESIS.md` draws it under the whole system: Git is the lowest durable articulation, not the source of coherence.
CDD draws it around work: a cell closes through receipt, verdict, and decision.
TSC draws it around measurement: coherence must be observed, not assumed.
This paper draws it under the model: the LLM is a provider the kernel governs.

---

## 18. Conclusion

Vendor runtimes have a strong case.
They keep context close to the model. They can plan, act, remember, and hide friction.

cnos makes a different trade.
It keeps the work body outside the runtime.

Let the model be smart.
Keep workflow, evidence, trust claims, coherence witnesses, receipts, and release authority at a seam with a name.

The architecture does not make judgment correct.
It makes judgment land somewhere visible.

If the receipt is signed, the signature should be visible.
If the receipt is attested, the attestation should be visible.
If coherence was measured, the report should be visible.
If coherence was not measured, the absence should be visible.

The system should not claim more trust than the artifact carries.
It should not claim more coherence than it measured.
It should not confuse coherence with correctness.

The rest is a commons that agrees to look.

Dumb models.
Smart cells.
Receipts over vibes.

The model should be smart enough to move the work forward.
It should be bounded enough that it never owns the work.

---

## References

- Google Cloud Blog, "Introducing Gemini Enterprise Agent Platform." https://cloud.google.com/blog/products/ai-machine-learning/introducing-gemini-enterprise-agent-platform
- Google Cloud Docs, "Agent Platform Memory Bank." https://docs.cloud.google.com/gemini-enterprise-agent-platform/scale/memory-bank
- OpenAI API Docs, "File search." https://developers.openai.com/api/docs/guides/tools-file-search
- Anthropic Engineering, "Scaling Managed Agents: Decoupling the brain from the hands." https://www.anthropic.com/engineering/managed-agents
- Letta Docs, "Letta Code." https://docs.letta.com/letta-code/
- Mem0 Docs, "Build with Mem0." https://docs.mem0.ai/introduction
- Model Context Protocol Docs, "Introduction." https://modelcontextprotocol.io/docs/getting-started/intro
- Agent2Agent Protocol, GitHub repository. https://github.com/a2aproject/A2A
- EU AI Act Service Desk, "Article 12: Record-keeping." https://ai-act-service-desk.ec.europa.eu/en/ai-act/article-12
- NIST, "AI Risk Management Framework." https://www.nist.gov/itl/ai-risk-management-framework
- OWASP GenAI Security Project, "State of Agentic AI Security and Governance." https://genai.owasp.org/resource/state-of-agentic-ai-security-and-governance/
- IETF Internet-Draft, "Agent Audit Trail: A Standard Logging Format for Autonomous AI Systems." https://datatracker.ietf.org/doc/draft-sharif-agent-audit-trail/
- IETF Internet-Draft, "Agent Identity Framework: Trust and Identity for Autonomous AI Agents." https://datatracker.ietf.org/doc/draft-sharif-agent-identity-framework/
- W3C, "Verifiable Credential Data Integrity 1.0." https://www.w3.org/TR/vc-data-integrity/
- Rails Guides, "Active Record Basics." https://guides.rubyonrails.org/active_record_basics.html
- TSC (Triadic Self-Coherence), repository and docs. https://github.com/usurobor/tsc — docs/THESIS.md, README.md, ARCHITECTURE.md, spec/tsc-glossary.md, engine/ocaml/test/fixtures/report.schema.json
- docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md
- docs/THESIS.md
- docs/reference/runtime/PROVIDER-CONTRACT-v1.md
- docs/reference/runtime/HYBRID-LLM-ROUTING.md
- src/packages/cnos.cdd/skills/cdd/CDD.md
