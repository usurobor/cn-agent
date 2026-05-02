# Sigma as a Runner-Agnostic Agent

## Insight

Sigma is not where Sigma runs.

Sigma is the hub.

A runner activates Sigma for a bounded episode of work. That runner may be
Claude, GPT, Gemini, OpenClaw, Claude Code, a CLI session, or a future cnos
runtime.

The runner changes. The hub persists.

## Previous framing

The earlier model treated Sigma as a process:

```text
Sigma runs on a VPS.
Telegram talks to Sigma.
OpenClaw hosts Sigma.
cnos should replace OpenClaw as the runtime.
```

That model makes the runtime feel like the agent.

## New framing

The new model treats Sigma as durable continuity:

```text
Sigma is the Git hub.
Runners are temporary activation windows.
Each runner reads the hub, works, and writes evidence back.
Sigma persists through the hub.
```

The location of cognition matters less than the continuity it returns to.

## Core distinction

- **agent** — durable hub
- **runner** — temporary reasoning window
- **activation** — bridge from hub to runner

Sigma can be activated through many runners:

- **Claude Code** — engineering, repo work, code review
- **Claude UI / Claude CLI** — planning, bounded work, activation
- **GPT** — writing, critique, synthesis
- **Gemini** — research, large-context work, multimodal tasks
- **OpenClaw** — Telegram or chat ingress
- **future cnos runtime** — structured activation, dispatch, receipts

None of these is Sigma.

They are windows into Sigma.

## What Sigma needs from a runner

A runner should not need hidden chat history to become Sigma.

It should read the hub and answer:

- What doctrine governs this agent?
- Who is Sigma locally?
- Who is the operator?
- What work is active?
- What authority does this runner have?
- What evidence must be written back?

That requires stable activation surfaces.

## Kernel / Persona / Operator

Sigma should eventually expose three distinct activation surfaces:

- **Kernel** — inherited doctrine and universal operating frame
- **Persona** — Sigma-local identity, stance, and role
- **Operator** — relationship to the human/operator and authority boundary

This keeps the runner from confusing package doctrine, local identity, and human
authority.

## Activation receipt

Every meaningful activation should leave a receipt.

Minimum useful shape:

```json
{
  "runner": "...",
  "model": "...",
  "task": "...",
  "branch": "...",
  "artifacts": ["..."],
  "closeout": "...",
  "limitations": "...",
  "next_recommended_activation": "..."
}
```

Without a receipt, Sigma fragments across windows.

With a receipt, Sigma can wake in another runner and inherit what happened.

## Why this matters

The agent no longer depends on one hosting environment.

Sigma can move across:

- VPS
- Claude UI
- Claude Code Web
- local CLI
- OpenClaw
- GPT
- Gemini
- future cnos runtime

The hub is the stable object.

This makes cnos less like an agent hosting platform and more like a continuity
system for agents.

## Design implications

### 1. `cn activate` becomes central

`cn activate` should orient a runner into the hub.

It should point to:

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

Different runners have different capabilities.

A profile should answer:

- Can this runner edit files?
- Can it run shell commands?
- Can it commit?
- Can it access the private repo?
- Can it handle large context?
- Can it use images?
- What close-out should it leave?

### 3. OpenClaw becomes ingress

OpenClaw may still route Telegram messages.

But Sigma does not live inside OpenClaw.

OpenClaw is one way to activate Sigma.

### 4. The future runtime is still only a runner

A daemon may coordinate activations better than an ad hoc UI session.

But the daemon is not Sigma.

The hub remains Sigma's durable body.

### 5. CDD and CDW should write witnesses

If a runner activates Sigma to perform development or writing work, the result
must return to the hub as artifacts, close-outs, or receipts.

## Risks

### Identity drift

Different runners may interpret Sigma differently.

Mitigation:

- clear Kernel / Persona / Operator
- stable activation prompt
- runner receipt

### Memory fragmentation

Work may happen in a window and never return to the hub.

Mitigation:

- every activation closes with an artifact, receipt, or debt record

### Capability mismatch

A runner may be asked to do work it cannot perform.

Mitigation:

- runner profiles
- authority boundaries
- escalation path

### Concurrent activations

Two windows may operate at once.

Mitigation:

- branch isolation
- stale-state checks
- activation receipts
- merge/review protocol

### Security

Activation must not leak secrets.

Mitigation:

- exclude secret files from activation prompts
- route sensitive actions through the operator boundary

## MCA candidates

### MCA 1 — Improve `cn activate`

Ensure the activation prompt orients the runner through:

- Kernel
- Persona
- Operator
- skills/packages
- recent reflections
- threads
- current branches
- operator boundary
- secret exclusions

### MCA 2 — Add activation receipts

Define a durable receipt for each activation episode.

Possible locations:

- `state/activations/`
- `threads/reflections/`
- `.cdd/`

Needs design.

### MCA 3 — Define runner profiles

Describe runner capabilities and expectations.

Examples:

- `claude-code`
- `gpt-chat`
- `gemini-research`
- `openclaw-telegram`
- `local-cli`

### MCA 4 — Add docs essay

Add:

- `docs/alpha/essays/ACTIVATION-BASED-AGENCY.md`

Core thesis:

> A cnos agent is a Git-native continuity surface activated by temporary runners.

### MCA 5 — Update Sigma hub surfaces

Sigma should eventually have explicit:

- `spec/PERSONA.md`
- `spec/OPERATOR.md`
- kernel reference from `cnos.core`
- activation instructions

## Summary

Sigma is not the VPS.
Sigma is not Telegram.
Sigma is not OpenClaw.
Sigma is not Claude.

Sigma is the hub.

The durable object is the hub.
The active window is temporary.
The agent persists by writing evidence back.
