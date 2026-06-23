---
title: "Dumb Models, Smart Cells"
subtitle: "Receipts, Trust Levels, and Named Authority for Agent Systems"
version: v0.7.0
status: DRAFT
author: usurobor (aka Axiom) (human & AI)
date: 2026-06-23
---

# Dumb Models, Smart Cells
## Receipts, Trust Levels, and Named Authority for Agent Systems

**Status:** v0.7.0 (DRAFT — position paper)
**Author(s):** usurobor (aka Axiom) (human & AI)
**Date:** 2026-06-23

> **Scope:** This paper explains why cnos treats language models as bounded executors, not as the home of workflow, memory, identity, evidence, trust, or release authority. It complements the CN protocol whitepaper, which explains Git as the communication substrate, and `THESIS.md`, which explains cnos as a recurrent coherence system.

---

## 0. The claim

The hard question is not whether models are smart.
They are.
The hard question is where authority lands when model output becomes work.

A model can write code, inspect files, draft plans, summarize evidence, and propose changes. A frontier model can hold more context, plan across longer horizons, and coordinate more steps than a small local model.

cnos does not deny that capability.
It refuses to confuse capability with authority.

Capability is rented.
Authority is owned.
Authority has to have a name.
And the receipt has to say what kind of name it is.

A plain label, a cryptographic signature, and a third-party attestation are not the same thing. A serious agent system should not smuggle one under the language of another.

That is the paper's claim: agent systems should rent model capability while keeping workflow, memory, evidence, validation, permissions, receipts, trust level, and release authority in a durable system the operator controls.

The model may move the work forward.
It should not own the work.

---

## 1. Vendor runtimes are tempting for a real reason

It is too easy to describe vendor runtimes only as lock-in.
They are also useful.

A hosted model runtime can hold context, remember files, call tools, run workflows, plan across a session, and hide a lot of friction. When the model and the workflow live together, the experience can be smoother. Less glue code. Fewer seams. More apparent continuity.

That is not fake value.
It is capability.

The product gravity is still real. A prompt becomes a thread. A thread gets memory. Memory gets files. Files get tools. Tools get workflows. Workflows get hosted agents. Hosted agents get logs, permissions, background tasks, connectors, deployment surfaces, registries, identities, gateways, observability, and governance surfaces.

Each feature is useful by itself.
Together they turn the vendor runtime into the place where work lives.

This is no longer a prediction. It is the roadmap.

Google's Gemini Enterprise Agent Platform brings agent development, runtime, registry, identity, gateway, observability, and governance into one platform. Google Memory Bank gives agents long-term memories across sessions. OpenAI exposes file search over vector stores. Anthropic is building managed agent infrastructure around Claude. The shape is clear: model capability is becoming a full work environment.

That is not a conspiracy.
It is what a capable product wants to become.

The risk is not that vendors are evil. The risk is that the easiest place to put the next feature is also the easiest place to lose authority.

At first the model answers.
Then it remembers.
Then it plans.
Then it acts.
Then it explains what it did.
Then the explanation, memory, workflow, logs, and decision boundary all live inside someone else's product.

That may be fine for some work.
It is not fine for all work.

---

## 2. The boundary costs something

A bounded model is not a free lunch.

If you keep workflow outside the vendor runtime, you give up some of the magic that comes from letting the runtime hold everything.

You may lose long-session continuity.
You may lose vendor-native planning features.
You may lose hidden context the model product was quietly carrying.
You may lose speed.
You may add ceremony: cells, receipts, validators, routing policy, release gates, package doctrine, and trust-level bookkeeping.

That cost is real.

A skeptic can state the objection plainly:

> You are trading capability for auditability.

Sometimes the skeptic is right.

If the work is exploratory, cheap, private, and easy to throw away, a hosted agent runtime may be enough. Let the model run. Let the thread hold the context. Let the product own the workflow.

cnos starts to earn its keep when the work matters enough that chat history is not a safe system of record.

The threshold is not philosophical.
It is practical.

Use a stronger boundary when the work:

- changes durable repo state,
- affects other people,
- crosses a release boundary,
- needs review,
- spans more than one session,
- spans more than one agent,
- needs audit later,
- can fail in a way people care about.

In those cases, the system needs more than a capable model.
It needs a place where authority lands.

---

## 3. "Dumb" means non-authoritative

"Dumb model" does not mean weak model.
It means bounded model.

A bounded model can still do hard work. cnos can route low-context work to a local model and escalate architecture, cross-package refactors, release work, and policy reconciliation to a frontier model when the task needs macro-coherence.

The boundary is authority.

The model may produce matter.
The model may propose changes.
The model may summarize evidence.
The model may recommend a next move.

But the model does not decide its own permissions. It does not choose its own provider route. It does not validate its own receipt. It does not own identity, memory, or the definition of done.

That is the difference between a smart product and a smart system.

In a smart product, the model runtime absorbs authority.
In a smart system, the model performs a move inside a boundary the system owns.

---

## 4. Replaceable does not mean interchangeable

The model slot is replaceable.
The models are not equivalent.

This is not a small distinction.

A local model, a frontier model, a code-specialized model, and a long-context model do not have the same capability. Swapping one for another can change whether the work succeeds at all.

So cnos should not promise magic portability.

The promise is narrower and more useful: model choice becomes an explicit system decision instead of a hidden habit.

A small task may run locally.
A release boundary may require a frontier model.
A local attempt may fail validation and escalate.
A remote provider may be unavailable, and the system may fail closed instead of silently pretending the local model is good enough.

That is not fungibility.
It is governed substitution.

The model can change without moving workflow, memory, evidence, identity, or release authority into the vendor runtime.

That is the property cnos is after.

---

## 5. The cnos boundary

cnos draws the boundary in the provider contract.

A provider is an executable capability endpoint. It may be `cnos.llm.local`, `cnos.llm.anthropic`, `cnos.net.http`, or `cnos.transport.git`. The shape is the same: the provider declares what it can do; the kernel decides whether and how it may act.

The provider contract states the rule directly:

> Package distributes. Extension declares. Provider executes. Kernel decides.

An LLM backend is not special because it sounds intelligent. It sits in the provider slot. It receives a typed operation, effective permissions, limits, and context. It returns a typed result, artifacts, timing, or an error.

The provider may be powerful.
The provider is not sovereign.

Routing follows the same rule. Model choice belongs in runtime policy, not in prompt folklore. A skill may shape the task. A command may declare a hint. But the runtime decides which provider runs.

Every routed model call must leave a receipt: route decision, reason set, provider, model name, fallback status, latency, token estimate, and validation signal when available.

Do not assume two engines are equal.
Do keep the workflow outside either engine.

Change the provider when policy allows it. Escalate when capability requires it. Fail closed when the available engine is not good enough.

The application shape stays in cells, packages, receipts, and repo state.

---

## 6. Where the work lives

If the model does not own the work, something else must.

In cnos, durable work lives in four places: cells, packages, receipts, and repo state.

### Cells own the workflow

A cell is a bounded unit of work. The CDD kernel names the generic loop:

`contract -> matter -> review -> receipt -> verdict -> decision`

The kernel is substrate-independent. It names roles, artifacts, validator `V`, evidence, verdicts, and decisions. It does not name GitHub, Claude, prompts, CI, branches, or any other invocation surface.

A model can produce matter for a cell. But the cell gives the work its shape. The model can change. The cell loop remains.

### Packages own local cognition

Skills, doctrine, runtime adapters, and domain rules arrive as packages.

That means an agent does not depend on a vendor's hidden system prompt to know how this repo works. It wakes into a local body of installed cognition.

The package can be versioned.
The package can be reviewed.
The package can be replaced.

### Receipts own evidence

A receipt is the durable close-out artifact.

It records what happened, what evidence was bound, what decision was made, and what can cross the boundary.

The receipt is not a chat summary.
It is the artifact another human, agent, validator, or release step can inspect later.

### Repo state owns continuity

The repo is the durable software body.

CN uses Git because Git already gives history, signatures, diffs, branches, tags, merges, clones, and offline operation.

Git is not the source of coherence. It is the lowest durable surface where coherence becomes inspectable.

The agent's memory should not disappear when a chat ends.
The repo remembers.

---

## 7. A trace, not just a principle

Here is the shape in miniature.

An operator asks Sigma to turn a new repo into a Hello World app.

Sigma does not implement immediately. It creates an issue contract first.

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
status: ready
```

The human reviews the issue and authorizes dispatch.

```yaml
dispatch:
  issue: 1
  protocol: cdd
  status: todo
  authority:
    subject: peter@operator
    trust_level: L1_named_within_repo
    signature: null
    attestation: none
  authorized_at: 2026-06-22T14:31:00Z
```

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
  fallback: none
  latency_ms: 1842
  token_estimate: 18231
```

The cell produces matter on dev, binds evidence, and closes with a receipt.

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
  validator:
    name: cdd.v.basic-release
    verdict: FAIL
    reason:
      - validation command recorded but not executed in clean CI
  delta:
    authority:
      subject: peter@operator
      trust_level: L1_named_within_repo
      signature: null
      attestation: none
    decision: override
    reason:
      - demo repo
      - no production users
      - degradation recorded
  outcome: degraded
  output:
    branch: deploy
    tag: v0.1.0
    tag_signature: null
```

This example is intentionally `L1_named_within_repo`.

That means the decision is attributable inside the repo's trust commons. It does not claim cryptographic non-repudiation, third-party timestamping, or legal admissibility.

A stronger receipt can add a cryptographic signature over the receipt object, commit, tag, or release artifact. A stronger still receipt can add third-party attestation. Those are higher trust levels.

The base claim here is narrower: the decision is no longer ambient chat. It lands at a named seam with evidence attached.

The point is not that this receipt proves the app is correct.

It does not.

The point is that the decision has a body.

A later agent can see what ran, which model ran, what evidence was bound, which validator was weak, who accepted the degradation, and what tag crossed the boundary.

That is the difference between "the model said it was done" and "the system can inspect what done meant."

---

## 8. Receipts do not validate themselves

A receipt is not proof.
A receipt is a body the system can inspect.

For some work, validation is cheap. The validator can check that files exist, schemas parse, tests pass, commands ran, tags point to the expected commit, release notes exist, signatures verify, and policy transitions happened in the right order.

For other work, validation is hard. A large refactor may compile and still be wrong. A plan may sound good and still miss the real constraint. A research claim may cite evidence and still overstate it.

cnos does not remove that problem.
It makes the problem visible.

In CDD, `V` is the validator predicate. It checks the receipt against the contract. `delta` is the boundary authority.

A weak validator should not produce a strong claim.
A degraded release should say it is degraded.
A repairable failure should create repair work.

The goal is not to invent a cheap oracle for correctness.
The goal is to stop hiding the oracle inside chat.

---

## 9. Named is not signed

A name is not a signature.
A signature is not an attestation.
An attestation is not automatically admissible evidence.

Those are different claims.
cnos should not blur them.

A useful receipt needs to say what trust level it carries. One possible ladder:

```yaml
trust_levels:
  L0_label:
    claim: a string names an actor
    strength: useful for local logs
    limitation: anyone with write access can forge the label
  L1_named_within_repo:
    claim: the repo's trust commons attributes the decision to a subject
    strength: useful for team accountability
    limitation: not cryptographic non-repudiation
  L2_signed_artifact:
    claim: a key signed the receipt, commit, tag, or release artifact
    strength: tamper-evident within the key infrastructure
    limitation: proves key control, not necessarily human intent
  L3_external_attestation:
    claim: an external witness, timestamping service, CI authority, or identity provider attested to the event
    strength: stronger chain of evidence
    limitation: depends on the attester and retention model
  L4_compliance_record:
    claim: the record satisfies a specific compliance or legal evidentiary regime
    strength: jurisdiction-specific admissibility target
    limitation: requires policy, retention, chain of custody, and legal review
```

The example trace in this paper is L1.
It is named inside the repo.
It is not signed.
It is not third-party attested.
It is not a compliance record.

That is fine, as long as the paper says so.

Git signatures, signed tags, signed receipts, W3C-style verifiable credentials, third-party timestamping, and agent audit-trail formats can all strengthen the claim. But the receipt must show that strength. Prose should not imply it.

The rule is simple:

> Do not claim more trust than the artifact carries.

---

## 10. The floor is named authority

Push the question down far enough and it bottoms out.

Is the model trustworthy? Wrong question — bound it.
Is `V` strong? Sometimes. The receipt should say when it is not.
Does `delta` act well on the verdict? That is the floor.

cnos does not remove this floor. No architecture can.
Judgment cannot be deleted. It can only be placed.

In CDD, `V` checks the receipt against the contract. `delta` decides what crosses the boundary. It may accept, reject, repair-dispatch, or release with degradation recorded.

That means cnos does not guarantee good judgment.
It guarantees a place where judgment lands: under a name, at a seam, with the receipt attached.

That is different from chat.

In chat, authority diffuses. The model suggested something. The human nodded. The thread moved on. Later, nobody can say exactly where the decision lived.

In cnos, the override is attributable. A degraded release is not a mood. It is a recorded decision.
If the receipt carries a cryptographic signature, it is also a signed decision.
If it carries third-party attestation, it may be stronger evidence outside the commons.
Those are different claims.

A degraded cell also propagates. It becomes degraded matter for the next cell. A later boundary can refuse to build on it, ask for repair, or accept the risk again under another name.

And it can become visible as a pattern. If overrides become routine, `epsilon` can read the receipt stream and surface that as system behavior, not one-off noise.

This still does not force anyone to care.
A negligent community can rubber-stamp degraded work forever. cnos will record that too.

The guarantee is not moral. It is structural: not-caring leaves a trail.

The rest is culture. A commons has to agree that receipts and overrides are binding enough to answer.

Fail-closed versus degrade-and-record follows the same rule. It is not a fact about the model. It is a policy decision made at a named seam, by an authority responsible for the cost of stopping and the risk of shipping.

---

## 11. The neighbors are real

The choice is not only cnos versus vendor runtime.
That would be too easy.

There is already a middle ecosystem trying to keep parts of agent systems outside any one model vendor.

Stateful-agent systems keep memory and context across conversations. Dedicated memory layers make memory portable across orchestrators. Tool protocols standardize how models reach external systems. Agent-to-agent protocols standardize how agents discover each other and exchange work.

Those are real seams.

Letta is a stateful-agent system. Mem0 is a memory layer. MCP is a tool and context seam. A2A is an agent-to-agent communication seam.

cnos should not pretend those neighbors do not exist.

The distinction is where the boundary closes.

A memory layer can preserve what an agent knows.
A tool protocol can expose what an agent can call.
An agent protocol can move messages between agents.

cnos is after a different unit: releaseable work.

A cnos cell does not only remember context or call tools. It binds a contract, produced matter, evidence, validator verdict, boundary decision, trust level, and receipt into one body the next human or agent can inspect.

That is the gap.

Memory says what the agent may know.
Tools say what the agent may touch.
Agent protocols say who the agent may talk to.
Receipts say what happened, what evidence was bound, what failed validation, who accepted the risk, what trust level the decision carries, and what crossed the boundary.

cnos belongs at that boundary.

---

## 12. The Rails echo

Rails is not proof that LLMs are database engines.
They are not.

Databases had SQL, shared relational ideas, and years of adapter work. LLMs do not yet have a stable capability contract. A frontier model and a small local model are not interchangeable in the way two SQL-backed adapters can sometimes be interchangeable.

So the Rails analogy should not carry more than it can.

The useful echo is narrower.

In the early Rails era, the database was the powerful engine. It stored the data, ran queries, enforced constraints, and offered its own language for procedures, triggers, and vendor-specific behavior.

One camp put more logic in the database.
The other camp wanted application logic in code, behind conventions the framework owned.

Rails made the second bet. Active Record and migrations did not make databases weak. They gave developers a stable seam. The database stayed powerful, but the application stopped living inside stored procedures, triggers, and vendor-specific behavior.

cnos makes a similar move, but at a harder boundary. There is no stable LLM capability contract yet. That makes the convention more important, not less.

The point is not "LLMs are just databases."
The point is that powerful engines should not automatically become the place where the application lives.

A powerful engine can sit behind a convention.
The convention is what the developer keeps.

---

## 13. Why this is an open-source shape

Open source cannot win by copying every vendor feature one at a time.

That game favors the vendor. The vendor can integrate faster, polish harder, subsidize usage, and collapse more of the stack into one product.

Open source wins where Unix and Rails won: small parts, strong conventions, stable interfaces, inspectable state, and community extension.

For agent systems, that means:

> the repo owns state,
> the runtime owns effects,
> the protocol owns handoff,
> the cell owns workflow,
> the receipt owns evidence,
> the package owns local cognition,
> the model owns only the bounded move it was asked to make.

That is the shape cnos is trying to normalize.

A community can add a package.
A team can swap a model.
A domain can define its own evidence.
A validator can check a receipt.
A future agent can continue from repo state instead of guessing from chat history.

The model is replaceable at the workflow boundary because the work has a body outside the model.

---

## 14. What the boundary buys

It buys substitution, not fungibility.
A team can add a new provider without moving workflow into that provider. The new provider may be better or worse. The receipt records which one ran, why it was chosen, and what validation followed.

It buys audit, not certainty.
A receipt does not prove the work was correct. It records what the system believed, what evidence was bound, what validator ran, and what decision crossed the boundary.

It buys explicit trust level.
A named decision, a signed decision, an externally attested decision, and a compliance record are not the same thing. The receipt can say which one it is. That prevents the system from smuggling weak evidence under strong words.

It buys durability.
The work survives a model call, a vendor outage, a price change, a lost chat, or a new runner. The repo keeps the body.

It buys cost control.
Simple work can stay local. High-context work can escalate. Failed local work can repair-dispatch or route up. The decision is policy, not vibes.

It buys audit readiness, not automatic compliance.
The EU AI Act's Article 12 requires high-risk AI systems to technically allow automatic event logging over the lifetime of the system. NIST and OWASP governance work point in the same direction: systems need records, monitoring, risk management, and evidence. A cnos receipt is not a legal compliance program. It is the kind of body a compliance program can inspect.

The common pattern is simple.
Keep durable concerns where the system can reach them.

---

## 15. Relationship to the CN whitepaper and THESIS

The CN whitepaper answers a substrate question:

> Where do agents coordinate?

Answer: Git, with a thin CN convention layer.

`THESIS.md` answers the system question:

> What is cnos?

Answer: a recurrent coherence system whose articulations include doctrine, documents, packages, runtime modules, repos, traces, releases, and agents.

This paper answers a placement question:

> Which parts of agent work must not live inside the vendor model?

Answer: workflow, memory, identity, evidence, permissions, receipts, trust level, and release authority.

The three papers draw the same boundary at different layers.

CN draws it under transport: forges are projections; the repo is the substrate.
`THESIS.md` draws it under the whole system: Git is the lowest durable articulation, not the source of coherence.
This paper draws it under the model: the LLM is a provider the kernel governs, not an authority the system obeys.

---

## 16. Conclusion

The strongest case for vendor runtimes is capability.
They are good because they keep context close to the model. They can plan, act, remember, and hide friction.

The strongest case for cnos is not that this capability is fake.
It is that capability should not silently become authority.

Let the model be smart.
Keep memory, workflow, identity, evidence, receipts, trust level, and release authority outside it — at a seam with a name.

The architecture does not make judgment correct.
It makes judgment land somewhere you can see, under a name that can be asked.
If the receipt is signed, the signature should be visible.
If the receipt is attested, the attestation should be visible.
The system should not claim more trust than the artifact carries.
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
- W3C, "Verifiable Credential Data Integrity 1.0." https://www.w3.org/TR/vc-data-integrity/
- Rails Guides, "Active Record Basics." https://guides.rubyonrails.org/active_record_basics.html
- [`docs/alpha/protocol/WHITEPAPER.md`](../protocol/WHITEPAPER.md)
- [`docs/THESIS.md`](../../THESIS.md)
- [`docs/alpha/agent-runtime/PROVIDER-CONTRACT-v1.md`](../agent-runtime/PROVIDER-CONTRACT-v1.md)
- [`docs/alpha/agent-runtime/HYBRID-LLM-ROUTING.md`](../agent-runtime/HYBRID-LLM-ROUTING.md)
- [`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../../src/packages/cnos.cdd/skills/cdd/CDD.md)
