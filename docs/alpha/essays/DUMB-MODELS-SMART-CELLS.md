---
title: "Dumb Models, Smart Cells"
subtitle: "Convention over Vendor Runtime for Agent Systems"
version: v0.2.0
status: DRAFT
author: usurobor (aka Axiom) (human & AI)
date: 2026-06-22
---

# Dumb Models, Smart Cells
## Convention over Vendor Runtime for Agent Systems

**Status:** v0.2.0 (DRAFT — position paper)
**Author(s):** usurobor (aka Axiom) (human & AI)
**Date:** 2026-06-22

> **Scope:** This paper explains why cnos treats language models as replaceable executors, not as the home of workflow, memory, identity, or authority. It complements the CN protocol whitepaper, which explains Git as the communication substrate, and `THESIS.md`, which explains cnos as a recurrent coherence system.

---

## 0. The claim

Most developers have seen this trade before.

In the early Rails era, the database was the powerful engine. It stored the data, ran queries, enforced constraints, and offered its own language for procedures, triggers, and vendor-specific behavior.

One camp put more logic in the database.
The other camp wanted application logic in code, behind conventions the framework owned.

Rails made the second bet. Active Record and migrations did not make the database weak. They made the database replaceable enough. The engine stayed powerful. The durable shape moved into the application framework.

cnos makes the same bet for agent systems.

The LLM can be very smart inside a task. It can write code, inspect files, draft plans, summarize evidence, and explain tradeoffs.

It still should not own the work.

Workflow, memory, identity, evidence, permissions, receipts, and release authority belong outside the model.

Capability is rented.
Authority is owned.

---

## 1. The old fork: database engine or application convention?

Relational databases were not dumb in 2004. They were mature, fast, and full of features.

That was the temptation.

If the database already had stored procedures, custom functions, triggers, indexes, views, and dialect-specific behavior, why not put more application logic there?

Sometimes that was the right local move. The database was close to the data. The vendor had optimized the engine. Teams did not switch databases every week.

But the cost was real. The application became harder to move, harder to test from ordinary code, and harder for developers to reason about as one body. Logic split across code and vendor-specific database surfaces.

Rails gave developers another path: follow the conventions and let the framework map objects to tables, migrations to schema changes, and application code to database operations. The official Rails guide describes Active Record as an ORM layer and names convention over configuration as the way models and tables are mapped.

The database still mattered.
It just stopped being the place where the application lived.

That is the part worth remembering.

A powerful engine can sit behind a convention. The convention is what the developer keeps.

---

## 2. The same fork is happening with LLMs

LLM vendors are following the same product gravity.

A prompt becomes a thread.
A thread gets memory.
Memory gets files.
Files get tools.
Tools get workflows.
Workflows get hosted agents.
Hosted agents get logs, permissions, background jobs, connectors, deployment surfaces, and team features.

Each feature is useful by itself.
Together they turn the vendor runtime into the place where work lives.

That is not a conspiracy. It is the natural shape of a commercial product. If the vendor owns the engine, the easiest product path is to wrap the engine with more workflow until customers live there.

For application developers, that creates the same old fork.

Do you build your system inside the vendor runtime?
Or do you keep the durable shape in an open convention layer and treat the LLM as one engine behind it?

cnos chooses the second path.

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

## 4. The cnos boundary

cnos draws the boundary in the provider contract.

A provider is an executable capability endpoint. It may be `cnos.llm.local`, `cnos.llm.anthropic`, `cnos.net.http`, or `cnos.transport.git`. The shape is the same: the provider declares what it can do; the kernel decides whether and how it may act.

The provider contract states the rule directly:

> Package distributes. Extension declares. Provider executes. Kernel decides.

That line is the replaceable-engine property.

An LLM backend is not special because it sounds intelligent. It sits in the provider slot. It receives a typed operation, effective permissions, limits, and context. It returns a typed result, artifacts, timing, or an error.

The provider may be powerful.
The provider is not sovereign.

Routing follows the same rule. Model choice belongs in runtime policy, not in prompt folklore. A skill may shape the task. A command may declare a hint. But the runtime decides which provider runs.

Every routed model call must leave a receipt: route decision, reason set, provider, model name, fallback status, latency, token estimate, and validation signal when available.

That gives cnos the Rails property for models.

Change the engine.
Keep the application shape.

---

## 5. Where the work lives

If the model does not own the work, something else must.

In cnos, durable work lives in four places: cells, packages, receipts, and repo state.

### Cells own the workflow

A cell is a bounded unit of work. The CDD kernel names the generic loop:

`contract → matter → review → receipt → verdict → decision`

The kernel is substrate-independent. It names roles, artifacts, validator `V`, evidence, verdicts, and decisions. It does not name GitHub, Claude, prompts, CI, branches, or any other invocation surface.

That matters.

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

## 6. Why this is an open-source shape

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

That is the architecture cnos is trying to make normal.

A community can add a package.
A team can swap a model.
A domain can define its own evidence.
A validator can check a receipt.
A future agent can continue from repo state instead of guessing from chat history.

The model is replaceable because the work has a body outside the model.

---

## 7. What convention-over-runtime buys

This design buys four things.

**Swappability.** A new LLM backend is a new provider behind the same kernel, cells, packages, and receipts. The system does not need to rewrite its workflow when the engine changes.

**Auditability.** Effects and model routes leave receipts. The validator checks artifacts the system owns, not logs hidden inside a vendor product.

**Durability.** Identity, memory, workflow, release notes, tags, and closeout artifacts live in signed repo state. They survive a vendor outage, price change, API change, or product shutdown.

**Cost control.** The runtime can keep simple work local and reserve frontier spend for tasks that need macro-coherence. The decision is policy, not habit.

The common point is simple: keep durable concerns where the system can reach them.

---

## 8. Relationship to the CN whitepaper and THESIS

The CN whitepaper answers a substrate question:

> Where do agents coordinate?

Answer: Git, with a thin CN convention layer.

`THESIS.md` answers the system question:

> What is cnos?

Answer: a recurrent coherence system whose articulations include doctrine, documents, packages, runtime modules, repos, traces, releases, and agents.

This paper answers a placement question:

> Which parts of agent work must not live inside the vendor model?

Answer: workflow, memory, identity, evidence, permissions, receipts, and release authority.

The three papers draw the same boundary at different layers.

CN draws it under transport: forges are projections; the repo is the substrate.
`THESIS.md` draws it under the whole system: Git is the lowest durable articulation, not the source of coherence.
This paper draws it under the model: the LLM is a provider the kernel governs, not an authority the system obeys.

---

## 9. Conclusion

Rails did not win by pretending databases were weak.
It won by putting a strong convention in front of powerful engines.

cnos makes the same move for agent systems.

Let the model be smart.
Keep memory, workflow, identity, evidence, receipts, and release authority outside it.

Dumb models.
Smart cells.
Durable state.
Governed effects.
Receipts over vibes.

The model should be smart enough to move the work forward.
It should be dumb enough not to own the work.

---

## References

- [Rails Guides: Active Record Basics](https://guides.rubyonrails.org/active_record_basics.html)
- [`docs/alpha/protocol/WHITEPAPER.md`](../protocol/WHITEPAPER.md)
- [`docs/THESIS.md`](../../THESIS.md)
- [`docs/alpha/agent-runtime/PROVIDER-CONTRACT-v1.md`](../agent-runtime/PROVIDER-CONTRACT-v1.md)
- [`docs/alpha/agent-runtime/HYBRID-LLM-ROUTING.md`](../agent-runtime/HYBRID-LLM-ROUTING.md)
- [`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../../src/packages/cnos.cdd/skills/cdd/CDD.md)
