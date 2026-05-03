# Activation, Not Deployment

Most agent systems begin with a deployment question:

> How do we run this agent?

cnos begins with a continuity question:

> What makes the next activation still the same agent?

That difference changes the architecture. Deployment centers the runtime. Activation centers the accumulated agent.

A deployed agent is kept alive by a process, session, service, or orchestration layer. An activated agent is entered by a runner for a bounded episode of work. The runner may disappear. The agent persists if the episode writes evidence back to continuity.

## Adjacent systems answer different questions

**MCP** answers a runner-facing question:

> What tools, resources, and prompts can this model session access?

The Model Context Protocol defines resources as a standardized way for servers to expose data and content to clients for use as model context, and prompts as reusable templates that clients can discover and retrieve. That is valuable, but it describes what a runner can access, not what makes the agent persist across runners.

**Memory systems** answer another nearby question:

> What should this model remember about the user or prior conversations?

OpenAI's memory documentation distinguishes saved memories from referenced chat history, both of which help personalize future conversations. That is useful, but cnos continuity is broader than remembered facts. It includes doctrine, identity, operator relationship, skills, artifacts, receipts, close-outs, decisions, failures, and reviewable history.

**Agent runtimes** answer yet another question:

> How do we keep a stateful workflow or agent execution alive and resumable?

LangGraph, for example, explicitly frames itself around long-running, stateful agents with durable execution, persistence, human-in-the-loop control, and memory across sessions. That is a runtime-centered answer. cnos may eventually use runtime machinery, but the runtime is not the agent. The accumulated continuity is.

## The cnos question

cnos asks:

> What accumulated agent is this runner activating, and what must return to it?

The answer is not the model.
The answer is not the platform.
The answer is not the daemon.
The answer is not the chat session.

The answer is the continuity surface.

A cnos agent is persistent, portable, and evolving.

It is persistent because it survives any one runner.
It is portable because different runners can activate it.
It is evolving because each meaningful activation can change what it knows, how it works, and what evidence it carries forward.

## Activation

Activation is the event where a temporary runner enters an accumulated agent.

A runner activates the agent when it:

1. reads the continuity surface;
2. understands the current task and authority boundary;
3. performs bounded work;
4. writes evidence back.

This is why `cn activate` matters.

`cn activate` does not run the agent. It prepares a runner to enter the agent. It tells the runner what continuity surface it is entering, what to read first, what authority it has, what not to leak, and what evidence to return.

The shape is:

```text
continuity surface → activation prompt → runner orientation → work → witness
```

## Worked example

Suppose a runner is asked to review a documentation PR.

The runner reads:

- Kernel: the operating doctrine;
- Persona: the local identity and stance;
- Operator: the human relationship and authority boundary;
- the issue;
- the changed document;
- the document-review skill;
- recent relevant reflections.

The runner performs the review.

It writes:

- a review artifact;
- a close-out;
- a receipt recording runner, model, branch, task, artifacts, limitations, and next step.

Then the runner can disappear.

A later runner does not need the old chat session. It can read the review artifact, close-out, and receipt from the continuity surface. That is how the agent persists across runners.

## Receipt

A receipt is how one episode becomes part of the agent.

A concrete receipt might look like this:

```json
{
  "runner": "claude-code-web",
  "model": "claude-sonnet",
  "started_at": "2026-05-02T14:10:00Z",
  "ended_at": "2026-05-02T14:42:00Z",
  "task": "review PR #297 for CTB document coherence",
  "branch": "claude/cleanup-ctb-docs",
  "artifacts": [
    "docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md",
    ".cdd/releases/3.74.0/327/beta-review.md"
  ],
  "closeout": ".cdd/releases/3.74.0/327/beta-closeout.md",
  "limitations": [
    "did not validate browser-rendered GitHub Markdown"
  ],
  "next_recommended_activation": "run document-review after CDW lands"
}
```

The receipt is not bureaucracy. It is the bridge between a temporary episode and durable continuity.

Without the receipt, the runner may have helped the moment.

With the receipt, the agent can inherit the work.

## Why activation is not always better

Activation has costs.

The runner must orient. The continuity surface must be read. Recent state may need to be summarized. The runner may need setup. The episode needs a close-out or receipt.

Deployment has different costs.

A deployed runtime must stay alive. It must supervise sessions, route channels, protect state, recover from failure, and manage live availability.

Activation wins when work is episodic, runners vary by task, and continuity matters more than low-latency presence.

Deployment wins when the agent must respond continuously, own a live channel, or maintain a high-frequency session.

cnos does not reject deployment. It refuses to confuse deployment with identity.

A future daemon can be useful. A chat gateway can be useful. A runtime can be useful.

They are activation surfaces.

They are not the accumulated agent.

## Why Kernel, Persona, and Operator are separate

Activation fails when a runner cannot tell which instruction has which authority.

Kernel, Persona, and Operator separate three different surfaces.

**Kernel** says what coherent agency means.
**Persona** says who this agent is locally.
**Operator** says how this agent relates to the human and what authority it has.

If those collapse into one vague prompt, doctrine becomes local preference, local identity becomes universal rule, and operator authority becomes hidden assumption.

For example, "be direct and decisive" may be a persona instruction. "Do not exceed operator authority" is an operator-boundary instruction. "Preserve evidence before closure" is doctrine. A runner should not treat those as the same kind of claim.

The split is not decorative. It prevents authority, identity, and doctrine from blurring at activation time.

## What MCP can become in cnos

MCP can still fit this architecture.

An MCP server could expose hub resources, package commands, prompts, or other capabilities to a runner. That may make activation easier. But MCP does not by itself define the accumulated agent or guarantee that a runner returns evidence to continuity.

In cnos terms:

- MCP can expose capabilities.
- `cn activate` orients the runner.
- The continuity surface preserves the agent.
- Receipts prove what returned.

The distinction is not hostile. It is a boundary.

## The rule

Do not ask first:

> Where does this agent run?

Ask:

> What continuity does this runner activate?
> What authority does it receive?
> What capabilities can it access?
> What witness will it return?
> How will the agent grow from this episode?

Deployment keeps a process alive.

Activation lets an accumulated agent keep becoming itself across many processes.
