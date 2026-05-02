# Sigma as a Runner-Agnostic Agent

Sigma is the hub. Runners are temporary windows that activate Sigma for bounded
work.

The runner may be Claude, GPT, Gemini, OpenClaw, Claude Code, or a future cnos
runtime. The runner changes. The hub persists.

For the doctrine, see
[`docs/alpha/essays/AGENT-CONTINUITY.md`](../../docs/alpha/essays/AGENT-CONTINUITY.md).
That essay is canonical. This note records what the doctrine means for Sigma
operationally.

## What Sigma was

The earlier model treated Sigma as a process:

```text
Sigma runs on a VPS.
Telegram talks to Sigma.
OpenClaw hosts Sigma.
cnos should replace OpenClaw as the runtime.
```

That model makes the runtime feel like the agent.

## What Sigma is

The hub is Sigma's durable body. Runners read it, work, and write evidence back.

Sigma can be activated through:

- **Claude Code** — engineering, repo work, code review
- **Claude UI / Claude CLI** — planning, bounded work
- **GPT** — writing, critique, synthesis
- **Gemini** — research, large-context work, multimodal tasks
- **OpenClaw** — Telegram or chat ingress
- **future cnos runtime** — structured dispatch and receipts

These are windows into Sigma.

## Why this matters

The agent no longer depends on one hosting environment. cnos becomes less an
agent hosting platform and more a continuity system for agents.

## Design implications

### 1. `cn activate` becomes central

`cn activate` should orient a runner into the hub. It should point to:

- Kernel
- Persona
- Operator
- skills/packages
- recent reflections
- threads
- current branches
- operator boundary
- secret exclusions

It should not claim to run the agent.

### 2. Runner profiles matter

Different runners have different capabilities. A profile should answer:

- Can this runner edit files?
- Can it run shell commands?
- Can it commit?
- Can it access the private repo?
- Can it handle large context?
- Can it use images?
- What close-out should it leave?

### 3. Hosts are runners

OpenClaw, a future cnos daemon, and any other host are runners. They route
messages or coordinate activations. The hub remains Sigma's durable body.

### 4. CDD and CDW must write witnesses

If a runner activates Sigma to perform development or writing work, the result
returns to the hub as artifacts, close-outs, or receipts.

## MCA candidates

### MCA 1 — Improve `cn activate`

Ensure the activation prompt orients the runner through Kernel, Persona,
Operator, skills/packages, recent reflections, threads, current branches,
operator boundary, and secret exclusions.

### MCA 2 — Add activation receipts

Define a durable receipt for each activation episode. Possible locations:

- `state/activations/`
- `threads/reflections/`
- `.cdd/`

Needs design.

### MCA 3 — Define runner profiles

Describe runner capabilities and expectations. Examples:

- `claude-code`
- `gpt-chat`
- `gemini-research`
- `openclaw-telegram`
- `local-cli`

### MCA 4 — Update Sigma hub surfaces

Sigma should eventually have explicit:

- `spec/PERSONA.md`
- `spec/OPERATOR.md`
- kernel reference from `cnos.core`
- activation instructions

## Summary

The hub is Sigma. Runners are episodes. Continuity persists by writing evidence
back.
