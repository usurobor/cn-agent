# Activation-Based Agency

A cnos agent is a Git-native continuity surface activated by temporary runners.

The hub persists.
The runner is temporary.
Activation is the bridge between them.

## The agent is the hub

A cnos hub is a Git repository — the durable body of an agent. It carries
identity, doctrine, skills, memory, threads, branches, commits, close-outs,
receipts, and history. Those surfaces let a new runner understand what agent it
is activating and what work it is allowed to do.

A Claude session may disappear. A GPT chat may end. A Gemini window may be
closed. A VPS may go offline. A Telegram gateway may be replaced.

The hub remains.

## The runner is an episode

A runner is the temporary reasoning window that activates the hub.

Examples:

- Claude Code for engineering work
- Claude UI or Claude CLI for bounded reasoning and repo work
- GPT for writing, critique, or synthesis
- Gemini for research or multimodal work
- OpenClaw or another always-on gateway for chat ingress
- a future cnos runtime for structured dispatch

An always-on runner is still a runner — duration does not turn a window into an
agent.

The agent is the continuity that survives them.

## Activation is not deployment

Deployment asks:

```text
How do we keep this process alive?
```

Activation asks:

```text
How does this temporary runner become this agent long enough to do coherent work?
```

A deployed agent needs process supervision, routing, restarts, live sessions, and
runtime ownership.

An activated cnos agent needs a different surface:

- identity
- doctrine
- operator contract
- skills
- current threads
- recent reflections
- working branch
- authority boundaries
- expected close-out

The runner reads those surfaces, performs bounded work, and returns evidence to
the hub.

## What `cn activate` does

`cn activate` prepares the runner. It does not run the agent or start a daemon.

The command should answer:

- What hub is this?
- What should the runner read first?
- What identity and operator surfaces matter?
- What skills and memories are available?
- What must not be leaked?
- What evidence should the runner leave behind?

The shape is:

```text
hub state → activation prompt → runner orientation
```

The prompt is an orientation guide, not the hub itself. It points the runner at
the surfaces that matter.

## Kernel, Persona, Operator

Activation needs distinct surfaces.

- **Kernel** — inherited doctrine and universal operating frame
- **Persona** — local agent identity, stance, role, and self-description
- **Operator** — human/operator contract, authority boundary, and working relationship

The runner should not have to infer these from scattered files or prior chat
history.

A new activation should be able to read the hub and answer:

- What doctrine governs this agent?
- Who is this agent locally?
- Who is the operator?
- What authority does this runner have?
- What should be done next?

## Episodes need receipts

A runner may do meaningful work and then disappear.

If that work is not written back, the agent fragments.

Every substantial activation should leave a receipt or close-out.

A useful receipt records:

```json
{
  "runner": "claude-code-web",
  "model": "claude-sonnet",
  "started_at": "...",
  "ended_at": "...",
  "hub": "...",
  "branch": "...",
  "task": "...",
  "artifacts": ["..."],
  "closeout": "...",
  "limitations": "...",
  "next_recommended_activation": "..."
}
```

The receipt is the witness of one episode, not the agent.

The hub accumulates those witnesses; continuity survives through them.

## Failure modes

Activation-based agency fails when the runner's episode does not return coherent
evidence to the hub.

| Failure | Meaning | Minimum mitigation |
| --- | --- | --- |
| Identity drift | Different runners read the same hub differently. | Make Kernel, Persona, and Operator explicit. |
| Memory fragmentation | Work happens in a runner but is not written back. | Require a receipt, close-out, artifact, or debt record. |
| Capability mismatch | The runner is asked to do work it cannot perform. | Use runner profiles and explicit authority boundaries. |
| Concurrent activation | Two runners act on the same hub state at once. | Use branch isolation, stale-state checks, and review before merge. |
| Secret exposure | Activation surfaces leak private material. | Exclude secret files and route sensitive actions through the operator boundary. |

These mitigations do not require a runtime. They require the hub to preserve
the right surfaces and the runner to close its episode honestly.

## Design implications

Activation-based agency changes what cnos should build first. The hub must be
easy to activate, safe to enter, and hard to leave without evidence. That
implies a small set of mechanisms — an activation command, durable receipts,
runner profiles, close-out discipline, conflict handling, and operator
boundaries — sized to the hub, not to a long-running runtime. Hub- and
agent-specific applications belong with operational notes.

## Rule

Do not ask first:

```text
Where does this agent run?
```

Ask:

```text
What hub does this runner activate,
what authority does it have,
and what witness will it leave behind?
```

A cnos agent is not the runner.

The runner is an episode.

The hub is continuity.
