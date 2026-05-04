# The Continuity of the Agent

Related:
- [agent doctrine](https://github.com/usurobor/cnos/blob/main/docs/alpha/doctrine/README.md) — structural ground (CFA, EFA, JFA, IFA); see Doctrine lineage below
- [KERNEL](https://github.com/usurobor/cnos/blob/main/src/packages/cnos.core/doctrine/KERNEL.md) — canonical for the Kernel/Persona/Operator distinction this essay assumes
- [AGENT-OPS](https://github.com/usurobor/cnos/blob/main/src/packages/cnos.core/doctrine/AGENT-OPS.md) — operational output discipline; receipts and coordination ops at runtime
- [ACTIVATION-NOT-DEPLOYMENT](ACTIVATION-NOT-DEPLOYMENT.md) — companion essay; how activation differs from deployment, MCP, memory, and runtimes

---

A cnos agent is accumulated continuity that can be activated by temporary runners.

It is persistent because it survives any one session.
It is portable because different runners can activate it.
It is evolving because each meaningful activation can change what it knows, how it works, or what evidence it carries forward.

The runner is temporary.

The continuity is the agent.

## Doctrine lineage

This essay is an application of the agent doctrine to activation. It is not foundational doctrine and does not claim to derive new structure.

[Coherence for Agents](https://github.com/usurobor/cnos/blob/main/docs/alpha/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md) grounds the receipt requirement: claim, evidence, and verdict must stay connected. A runner episode that leaves no evidence asks the next activation to inherit unverifiable closure.

[Inheritance for Agents](https://github.com/usurobor/cnos/blob/main/docs/alpha/doctrine/inheritance-for-agents/INHERITANCE-FOR-AGENTS.md) grounds the continuity requirement: later cycles inherit named constraints, not appearances. A runner can write something back and still weaken continuity if it softens a prior commitment into a form that no longer constrains action.

[Judgment for Agents](https://github.com/usurobor/cnos/blob/main/docs/alpha/doctrine/judgment-for-agents/JUDGMENT-FOR-AGENTS.md) grounds debt and close-out language: when no move preserves every boundary, the agent must name what it protected, what it breached, and what remains unpaid.

[Ethics for Agents](https://github.com/usurobor/cnos/blob/main/docs/alpha/doctrine/ethics-for-agents/ETHICS-FOR-AGENTS.md) grounds the future-cycle concern: undocumented drift affects later agents that must inherit the surface left behind.

## The higher-level object

Most agent systems begin with the execution surface:

- model;
- tools;
- runtime;
- process;
- platform;
- channel.

cnos begins one level higher.

It asks what persists across those surfaces.

That persistent thing is the accumulated continuity of the agent:

- identity and stance;
- doctrine and operating frame;
- operator relationship;
- skills and habits;
- memories and reflections;
- decisions and judgments;
- failures and repairs;
- artifacts, receipts, and close-outs.

A runner may contribute to continuity, but the runner is not the continuity itself.

## Persistent

A cnos agent persists beyond any one activation.

A chat session can end. A CLI process can exit. A gateway can be replaced. A model provider can change. The agent remains coherent only if the important residue returns to its continuity surface.

Persistence means a later runner can inherit:

- what happened;
- what changed;
- what failed;
- what was learned;
- what remains unfinished.

Continuity is the answer to the question:

> What survives the episode?

## Portable

A cnos agent is not bound to one model or platform.

The same agent can be activated through different runners:

- Claude Code for engineering;
- GPT for writing or critique;
- Gemini for research or multimodal work;
- OpenClaw or another gateway for chat ingress;
- a local CLI for deterministic checks;
- a future cnos runtime for structured orchestration.

The runner is selected by the work.

Portability means the agent can move across runners without becoming a different agent. That only works if each runner reads the same continuity surface and writes evidence back to it.

## Evolving

A cnos agent should change when it runs.

A meaningful activation leaves something durable:

- work product;
- reflection;
- close-out;
- receipt;
- decision;
- repair;
- skill update;
- named limitation;
- debt record.

The value of an agent is not only what it can do now. The value is what it has accumulated through use.

An agent becomes more valuable when its episodes add up.

## The hub implements continuity

In cnos today, continuity lives in a hub.

The hub is the durable body that stores, versions, and exposes the agent's continuity. It gives the accumulated agent a place to persist between runners.

Git is not the philosophical point.

Continuity is the point.

cnos uses a hub because versioned repositories are practical, portable, inspectable, and legible to both humans and tools. Branches, commits, diffs, reviews, and history make the continuity durable enough for many runners to return to it.

The hub is the current body.

The agent is the continuity that body preserves.

## Activation

Activation is the bridge from continuity to runner.

A runner activates the agent when it:

1. reads the continuity surface;
2. understands the task and authority boundary;
3. performs bounded work;
4. writes evidence back.

`cn activate` prepares that bridge.

It does not run the agent, start a daemon, or replace the future runtime.

It generates an orientation surface for the runner.

A good activation prompt tells the runner:

- what continuity surface it is entering;
- what to read first;
- what identity and operator surfaces matter;
- what skills and memories are available;
- what authority it has;
- what not to leak;
- what evidence to return.

The prompt is not the agent.

The prompt helps the runner enter the agent.

## Kernel, Persona, Operator

Activation needs three distinct surfaces.

**Kernel** is inherited doctrine and the universal operating frame.

**Persona** is the local agent identity, stance, role, and self-description.

**Operator** is the human relationship, authority boundary, and working contract.

These should not collapse into one vague system prompt.

The Kernel says what coherent agency means.
The Persona says who this agent is.
The Operator surface says how this agent relates to the human and what authority it has.

A runner becomes this agent by reading those surfaces and acting within them.

## Receipts

A receipt turns an episode into continuity.

A receipt answers:

- which runner activated the agent;
- what task was attempted;
- what branch or surface changed;
- what artifacts were produced;
- what failed;
- what was learned;
- what should happen next.

A minimal receipt might look like this:

```json
{
  "runner": "claude-code-web",
  "model": "claude-sonnet",
  "started_at": "...",
  "ended_at": "...",
  "task": "...",
  "branch": "...",
  "artifacts": ["..."],
  "closeout": "...",
  "limitations": "...",
  "next_recommended_activation": "..."
}
```

The receipt is not bureaucracy.

It is how one episode becomes part of the agent.

## Failure modes

Continuity-based agency fails when temporary work does not return to durable continuity.

| Failure | Meaning | Minimum mitigation |
| --- | --- | --- |
| Identity drift | Different runners interpret the agent differently. | Make Kernel, Persona, and Operator explicit. |
| Memory fragmentation | Work happens in a runner but is not written back. | Require a receipt, close-out, artifact, or debt record. |
| Capability mismatch | A runner is asked to do work it cannot perform. | Use runner profiles and authority boundaries. |
| Concurrent activation | Two runners act on the same continuity surface at once. | Use branch isolation, stale-state checks, and review before merge. |
| Secret exposure | Activation surfaces leak private material. | Exclude secrets and route sensitive actions through the operator boundary. |
| False growth | The runner says it learned, but nothing changes in continuity. | Require an MCI, reflection, skill update, or explicit "no durable lesson" close-out. |

These are continuity failures before they are runtime failures.

The answer is not always to build a daemon.

Sometimes the answer is to write back better.

## What cnos is protecting

cnos is not primarily protecting a session.

It is protecting the accumulated agent.

That means the important question is not only:

```text
Can this runner do the work?
```

It is also:

```text
Will this work return to the agent?
```

A runner that produces useful output but leaves no durable evidence has helped the moment but not the agent.

A runner that writes evidence back lets the agent persist, move, and evolve.

## Rule

Do not ask first:

> Where does this agent run?

Ask:

> What continuity does this runner activate?
> What authority does it receive?
> What evidence will it return?
> How will the agent grow from this episode?

The agent is persistent.
The agent is portable.
The agent is evolving.

The runner is temporary.

The continuity is the point.
