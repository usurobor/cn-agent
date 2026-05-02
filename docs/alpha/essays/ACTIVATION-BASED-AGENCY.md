# Activation-Based Agency

A cnos agent is a Git-native continuity surface.

It becomes active when a capable runner reads the hub, acts under its identity
and authority surfaces, and writes evidence back.

The runner is temporary.
The hub persists.

## The agent is the hub

A cnos hub is a Git repository that carries an agent's continuity.

It contains identity, doctrine, skills, memory, threads, branches, commits,
close-outs, receipts, and history. Those surfaces let a new runner understand
what agent it is activating and what work it is allowed to do.

The hub is not just storage. It is the durable body of the agent.

A Claude session may disappear. A GPT chat may end. A Gemini window may be
closed. A VPS may go offline. A Telegram gateway may be replaced.

The hub remains.

## The runner is an episode

A runner is the temporary reasoning window that activates the hub.

Examples:

- Claude Code for engineering work;
- Claude UI or Claude CLI for bounded reasoning and repo work;
- GPT for writing, critique, or synthesis;
- Gemini for research or multimodal work;
- OpenClaw for Telegram or other chat ingress;
- a future cnos runtime for structured dispatch.

None of these is the agent.

They are activation windows.

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

## The activation bridge

`cn activate` is the bridge from hub to runner.

It does not run the agent.
It does not start a daemon.
It does not replace the future runtime.

It prepares the runner.

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

The prompt is not the whole hub. It is an orientation guide that tells the
runner how to enter the hub correctly.

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

That is why the activation surface must be explicit.

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

The receipt is not the agent.

It is the witness of one episode.

The hub accumulates those witnesses. That is how continuity survives across
temporary windows.

## OpenClaw is an ingress

An always-on gateway can still be useful.

Telegram, Slack, Discord, or another channel may wake a runner. OpenClaw or a
similar gateway may provide channel routing and message ingress.

But in cnos, the gateway is not where the agent lives.

It is one activation path.

Sigma is not the Telegram bot.
Sigma is not the OpenClaw process.
Sigma is not the Claude session.

Sigma is the hub those surfaces return to.

## Design implications

Activation-based agency changes what cnos should build first.

The first priority is not a full always-on runtime. The first priority is making
the hub easy to activate and safe to re-enter.

That means cnos needs:

- **`cn activate`** — generate the activation prompt from hub state
- **activation receipts** — record what one runner did
- **runner profiles** — describe what a runner can and cannot do
- **close-out discipline** — require evidence, limits, and next steps
- **conflict handling** — detect stale or concurrent activations
- **operator boundaries** — prevent temporary runners from exceeding authority

A future runtime can still exist.

But the runtime is not the agent. It is another runner, or another way to
coordinate runners.

## Risks

### Identity drift

Different runners may interpret the same hub differently.

Mitigation:

- clear Kernel / Persona / Operator surfaces
- stable activation prompt
- runner receipts

### Memory fragmentation

A runner may work without returning evidence to the hub.

Mitigation:

- every substantial activation closes with an artifact, receipt, or debt record

### Capability mismatch

Different runners have different abilities.

A Claude Code session can edit files and run shell commands. A plain chat
session may only reason. A gateway runner may have restricted tools.

Mitigation:

- runner profiles
- explicit authority boundaries
- clear escalation path

### Concurrent activations

Two runners may activate the same hub at once.

Mitigation:

- branch isolation
- stale-state checks
- activation receipts
- merge/review protocol

### Secret exposure

Activation must not leak secrets.

Mitigation:

- never include secret files in activation prompts
- state authority boundaries explicitly
- route sensitive operations through the operator gate

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
