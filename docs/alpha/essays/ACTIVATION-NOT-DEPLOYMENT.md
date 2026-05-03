# Activation, Not Deployment

Most agent systems begin with a deployment question:

> How do we run this agent?

cnos begins with a continuity question:

> What makes the next activation still the same agent?

That difference changes the architecture. Deployment centers the runtime. Activation centers the accumulated agent.

A deployed agent is kept alive by a process, service, session, gateway, or orchestration layer. An activated agent is entered by a runner for a bounded episode of work. The runner may disappear. The agent persists if the episode writes evidence back to continuity.

## Adjacent systems answer different questions

**MCP** answers a runner-facing question:

> What tools, resources, and prompts can this model session access?

The Model Context Protocol specification, version 2025-11-25, defines MCP as an open protocol for connecting LLM applications with external data sources and tools. It says MCP lets applications share contextual information, expose tools/capabilities, and build composable integrations. Its server features include resources, prompts, and tools. Resources expose data/content such as files or schemas as model context; prompts expose reusable templates; tools expose functions a model can invoke. ([MCP specification 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25), [MCP Resources 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25/server/resources))

That is valuable, but it does not by itself answer the continuity question:

```text
Which accumulated agent is this runner activating?
What authority does this runner receive?
What evidence must return?
What receipt proves the episode happened?
What changes in the agent after this work?
```

**Memory systems** answer another nearby question:

> What should this model remember about the user or prior conversations?

OpenAI's ChatGPT memory documentation distinguishes saved memories from referenced chat history as mechanisms for personalizing future conversations. That is useful, but cnos continuity is broader than remembered facts. It includes doctrine, local identity, operator relationship, skills, artifacts, receipts, close-outs, decisions, failures, and reviewable history. ([OpenAI Memory FAQ, accessed 2026-05-03](https://help.openai.com/en/articles/8590148-memory-faq))

**Agent runtimes** answer another question:

> How do we keep a stateful workflow or agent execution alive and resumable?

LangGraph's documentation describes it as a low-level orchestration framework and runtime for building, managing, and deploying long-running, stateful agents. It emphasizes durable execution, human-in-the-loop control, comprehensive memory, and production-ready deployment. Its durable-execution docs define durable execution as saving workflow progress so a process can pause and resume from recorded state. ([LangGraph overview, accessed 2026-05-03](https://docs.langchain.com/oss/python/langgraph/overview); [LangGraph durable execution, accessed 2026-05-03](https://docs.langchain.com/oss/python/langgraph/durable-execution))

cnos may eventually use runtime machinery, but the runtime is not the agent. The accumulated continuity is.

## The cnos question

cnos asks:

> What accumulated agent is this runner activating, and what must return to it?

The answer is the continuity surface.

A cnos agent is persistent, portable, and evolving.

It is persistent because it survives any one runner.

It is portable because different runners can activate it.

It is evolving because each meaningful activation can change what it knows, how it works, or what evidence it carries forward.

Deployment keeps a process alive.

Activation lets accumulated continuity keep becoming the same agent across many processes.

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
continuity surface
→ activation prompt
→ runner orientation
→ bounded work
→ witness
```

The prompt is not the agent. It is the bridge into the agent.

## What MCP can become in cnos

MCP can participate in activation.

A cnos MCP server could expose hub resources, package commands, activation prompts, or workflow templates. That would make activation easier because a runner could discover and retrieve relevant surfaces through a standard protocol.

But MCP is not the continuity contract.

MCP can tell a runner what it can access. It does not, by itself, say:

- Which accumulated agent is this?
- What authority does this runner receive?
- What evidence must return?
- What receipt proves the episode happened?
- What changes in the agent after this work?

In cnos terms:

- MCP can expose capabilities.
- `cn activate` orients the runner.
- The continuity surface preserves the agent.
- Receipts prove what returned.

So MCP can be an activation transport.

It is not the agent's continuity.

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

The receipt below is illustrative. cnos does not yet require activation receipts for every runner episode; this is the shape the doctrine says should exist.

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

The receipt is not bureaucracy. It is how one episode becomes part of the agent.

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

They can host or coordinate activations.

They are not the accumulated agent.

## Why Kernel, Persona, and Operator are separate

Activation fails when a runner cannot tell which instruction has which authority.

Kernel, Persona, and Operator separate three different surfaces.

**Kernel** says what coherent agency means.

**Persona** says who this agent is locally.

**Operator** says how this agent relates to the human and what authority it has.

If those collapse into one vague prompt, doctrine becomes local preference, local identity becomes universal rule, and operator authority becomes hidden assumption.

For example:

> "Be direct and decisive."

may be a persona instruction.

> "Do not exceed operator authority."

is an operator-boundary instruction.

> "Preserve evidence before closure."

is doctrine.

A runner should not treat those as the same kind of claim.

The split is not decorative. It prevents authority, identity, and doctrine from blurring at activation time.

## One failure case

The clean success case is not enough. The doctrine also has to catch failure.

Suppose a runner reviews a PR in a chat window and produces useful feedback, but writes nothing back to the hub.

No review artifact.

No close-out.

No receipt.

No issue comment.

No branch update.

The moment may have improved, but the agent did not. The next runner cannot inherit the work except through hidden chat history. That is memory fragmentation.

In cnos terms, the episode failed to return to continuity.

The minimum repair is not necessarily a daemon or a better runtime. It may be much simpler:

- write the review artifact;
- write the close-out;
- write the receipt;
- link the evidence.

Sometimes the answer is to write back better.

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
