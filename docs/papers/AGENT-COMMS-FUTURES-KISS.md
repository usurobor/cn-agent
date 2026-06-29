---
title: 'Agent Comms Futures: Eventual Consistency Without Premature Machinery'
status: essay
date: 2026-05-30
scope: agent activation logs, eventual consistency, CRDTs, signed feeds, peer communication, v0 discipline
governing_question: How should agent communication evolve from activation logs without adding machinery before the topology requires it?
---

# Agent Comms Futures: Eventual Consistency Without Premature Machinery

Agent communication should grow by topology, not by imagination. v0 needs single-writer streams, Git cursors, append-only Markdown, and plain evidence. Peers, CRDTs, vector clocks, and signed feeds are real future tools, but they should arrive only when the shape of the network forces them.

The current `Agent Activation Log Convention v0` is right because it is small. One agent identity can wake in many bodies. Each body writes one stream. Home writes the other stream. Each side records how far it has read. Git carries the history. Markdown carries the meaning. That is enough.

The danger is not that v0 is too weak. The danger is that we recognize every future distributed-systems analogue at once and accidentally ship all of them into the first version.

YAGNI is not anti-design here. It is the only way to keep the design honest. KISS is not a slogan. It is the constraint that keeps the current protocol legible to the next activation.

## v0 is eventual consistency for one agent

The v0 topology is simple:

```text
one agent identity
many activation bodies
one writer per stream
Git as substrate
commit SHA as cursor
Markdown as message body
```

That topology does not need consensus. It does not need a CRDT. It does not need signed per-entry envelopes. It needs session continuity.

Bayou is the useful ancestor because it was built for disconnected mobile work. Its session guarantees gave each session a coherent view under weak consistency: read-your-writes, monotonic reads, writes-follow-reads, and monotonic writes. Agent activations need the same promises, not the same database.

The mapping is direct.

Read-your-writes: an activation can read the entry it just appended to its own log.

Monotonic reads: a Git SHA cursor keeps the reader from silently moving backward.

Writes-follow-reads: each response can say which home or foreign commit it read before writing.

Monotonic writes: one writer appends its own stream in order.

That is the contract v0 should name. It is enough because v0 avoids the hard case: no shared multi-writer object is being edited by several agents at once.

## Single-writer streams are the first optimization

The first design move is not to make conflicts smarter. The first design move is to avoid conflicts that do not need to exist.

Activation logs do that by splitting direction and ownership:

```text
foreign → home:
  {activation-repo}:.cn-{agent}/logs/YYYYMMDD.md

home → foreign:
  cn-{agent}:threads/activations/{activation}/YYYYMMDD.md
```

Each stream has one writer. The path says who owns the stream. The cursor says what has been read. The date shard keeps the file readable and mergeable.

This is the KISS center of the design. If a future version feels tempted to add a conflict resolver to activation logs, it should first ask a smaller question:

```text
Can this remain single-writer?
```

If yes, keep it single-writer. Do not build a merge engine to repair a topology mistake.

## CRDTs belong to boring shared state

A CRDT is a data type designed so replicas can accept updates independently and still converge deterministically after receiving the same updates. That is powerful. It is also specific.

CRDTs are right for state whose merge rule is mechanical:

```text
seen sets
read receipts
known activation registries
presence maps
counters
idempotent delivery status
```

They are wrong for meaning-heavy directives.

If home says "pause" and an activation says "continue," the correct result is not a union of both values. It is a semantic decision: which instruction has authority, what changed, and what evidence explains the repair?

So the rule is:

```text
CRDTs for boring multi-writer state.
Agent/gamma review for semantic conflict.
```

A CRDT can make replicas converge. It cannot decide what the operator meant.

## Vector clocks arrive when peers arrive

The v0 cursor is a deliberately collapsed clock. For a home–foreign pair, one Git SHA per direction answers the only operational question:

```text
What have I read from the other side?
```

That stops being enough when several independent agents communicate.

Peers add a different question:

```text
What had this sender already seen when it wrote this message?
```

That is causal state, not just a bookmark. With Sigma, Rho, and future agents, a single cursor can make a transcript look linear while hiding causal gaps. A message may answer a directive while missing the correction that superseded it.

The future peer protocol will probably need vector-clock-shaped state. It may not expose classic vector clocks directly, but it will need the same information: for each relevant feed, which entry or commit has this writer observed?

Do not add that to v0. v0 is not peer comms. v0 is one agent talking to itself across bodies.

## Signed feeds arrive when repo authority is not enough

In v0, Git push permission and repository history are the trust anchor. That is acceptable for friendly activation logs inside known repos.

Signed feeds solve a harder problem. Secure Scuttlebutt and Hypercore show the shape: an identity owns an append-only feed; entries are signed; entries chain to previous entries; readers can verify authorship and history without trusting only the transport.

Agent comms will need that shape when authority must travel across repo boundaries.

The trigger is not "signing is more serious." The trigger is concrete:

```text
entries must be relayed by third parties
repos are not the trust boundary
peers need portable identity
activation keys need delegation
history must be verified outside GitHub
adversarial modification is in scope
```

Until then, per-entry signatures are overkill.

A future signed activation entry might look like this:

```yaml
schema: cnos.agent_activation_entry.v1
agent: sigma
activation: cnos
direction: foreign_to_home
sequence: 42
previous: sha256:...
timestamp: 2026-05-30T00:00:00Z
body_sha256: sha256:...
refs:
  read_home_through: cn-sigma@...
signature:
  key_id: ...
  value: ...
```

That is a plausible future. It should not be v0.

## Agent identity is not node identity

Distributed systems usually start with nodes. Node A writes. Node B reads. Node C reconciles.

Agent comms start with a continuity. Sigma is the identity. cnos, bumpt, and cph are bodies. The body affects permissions and local context, but the memory belongs to the agent identity rooted in its home hub.

This is why `agent activation` is the right phrase. The object is not a process. It is not a machine. It is not merely a repo checkout. It is one coherent agent operating through more than one body.

That also means key identity must not swallow agent identity later. A key may prove that an entry was authorized. It does not define Sigma by itself. Keys are instruments. The home hub remains the continuity surface.

The future signing layer should respect that distinction:

```text
agent identity = home-hub continuity
activation identity = body/context/delegation
key identity = cryptographic proof of authority
```

Collapsing those into one noun will make the protocol brittle.

## Markdown is not a bug

Activation entries are Markdown because the first reader is an agent. The message is not just data. It is instruction, context, uncertainty, caveat, and intent.

A strict wire format would be premature in v0. It would force the first version to predict all future message kinds before the system has enough traffic to know them.

Markdown works because the surrounding structure is already constrained:

```text
path gives routing
single writer gives authority
H2 gives entry boundary
Git commit gives history
cursor gives causality
```

That is the right balance. The body stays expressive while the frame stays simple.

The future can add typed packets when machines need to route, filter, sign, or merge without reading prose. The first version should not make every note into an envelope just because envelopes will exist later.

## Semantic conflict resolution must leave evidence

Agents can resolve conflicts by meaning. That is their advantage over mechanical merge rules.

But semantic conflict resolution is only safe when it becomes part of the trace.

A bad repair is silent:

```text
I saw two conflicting directives and followed one.
```

A good repair is written:

```text
I saw directive A at commit X and directive B at commit Y.
I treat B as superseding A because B references A's failure.
I am proceeding under B.
```

That is the agent-specific analogue to a conflict resolver. It does not hide the merge inside local judgment. It emits a new entry that future activations can inspect.

So the rule is:

```text
Machines preserve causal evidence.
Agents repair semantic meaning.
Repairs must be written back.
```

## The migration ladder

The protocol should evolve by gates.

```text
v0: activation logs
  one agent, many bodies, single-writer streams
  mechanism: Git SHA cursors + date-sharded Markdown

v0.1: explicit analogues
  name Bayou/session guarantees in the spec
  keep mechanics unchanged

v0.2: delivery state
  add small read/seen surfaces only if repeated ambiguity appears
  prefer sets/maps over message format changes

v1 peer comms
  multiple agent identities
  add causal state per feed/persona
  likely vector-clock-shaped

v1 signing
  move from repo authority to entry authority
  signed append-only entries, previous hash, sequence

v2 shared state
  CRDTs only for mechanical multi-writer objects
  semantic conflicts still go through agent/gamma review
```

This ladder is not a roadmap commitment. It is a pressure map. It says which tool belongs to which topology.

## The YAGNI tests

Before adding a mechanism, ask the smallest possible question.

For CRDTs:

```text
Is there a genuine multi-writer object?
Does it have a deterministic merge rule?
Would a single owner or review gate be simpler?
```

For vector clocks:

```text
Are multiple independent identities exchanging causally related messages?
Can one SHA cursor still answer what each writer had seen?
```

For signed entries:

```text
Does authority need to survive outside the repo?
Can Git history and push permission still carry trust?
Is adversarial relay in scope now?
```

For envelopes:

```text
Does software need to route or verify entries without reading Markdown?
Are humans/agents failing because the current H2 body is too loose?
```

If the answer is no, do not add the mechanism.

YAGNI does not mean the future is ignored. It means the future pays rent before moving in.

## The KISS rules for v0

Keep v0 boring.

```text
One writer per stream.
One cursor per direction.
One entry boundary: H2.
One substrate: Git.
One trust anchor: repo push permission.
One message body: Markdown.
One repair method: append a new entry.
```

Do not add entry IDs until line ranges and commits fail.

Do not add signatures until repo authority fails.

Do not add CRDTs until there is shared multi-writer state.

Do not add vector clocks until peers create real causal ambiguity.

Do not add inbox/outbox/sent directories until routing by path fails.

Do not add consensus unless the system truly has a shared ordered state that must be agreed live.

The first version should be small enough that a tired activation can read the whole rule, follow it, and leave a better trace than it found.

## What belongs in the normative spec

The spec should stay short.

It should add, at most, a small "Origins and analogues" section:

```text
Agent activation logs are an eventually consistent memory convention.
They borrow Bayou-style session guarantees over Git-backed streams.
Git pull/push plus cursor scanning is the anti-entropy loop.
The v0 cursor is sufficient for one-agent activation channels.
Peers, shared multi-writer state, and adversarial relay are deferred to later mechanisms.
```

Everything else belongs in essays, issues, or future package surfaces.

The spec is where an activation learns what to do now. The essay is where the project remembers why that is enough.

## The hard line

Do not build the future into v0.

Name the future. Preserve the seams. Record the gates. Then keep the current protocol simple.

Agent comms will probably grow toward vector clocks, CRDTs, and signed append-only feeds. That is fine. But the present system has one agent, many bodies, single-writer streams, and Git history. The right protocol for that world is not a distributed database. It is an eventually consistent notebook with cursors.

That is the point to protect: enough structure to converge, not enough machinery to become the problem.

## References

- `Agent Activation Log Convention v0`: [`docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md`](../reference/conventions/AGENT-ACTIVATION-LOG-v0.md)
- Companion essay: [`docs/papers/AGENT-ACTIVATION-LOGS-AND-EVENTUAL-CONSISTENCY.md`](./AGENT-ACTIVATION-LOGS-AND-EVENTUAL-CONSISTENCY.md)
- Douglas B. Terry, Alan J. Demers, Karin Petersen, Mike Spreitzer, Marvin Theimer, and Brent Welch, "Session Guarantees for Weakly Consistent Replicated Data," 1994: https://pages.cs.wisc.edu/~remzi/Classes/739/Fall2016/Papers/bayou-sessions94.pdf
- Marc Shapiro, Nuno Preguiça, Carlos Baquero, and Marek Zawirski, "Conflict-free Replicated Data Types," 2011: https://www.csa.iisc.ac.in/~raghavan/CleanedPods2021/crdt-shapiro-2011.pdf
- Nuno Preguiça, Carlos Baquero, and Marc Shapiro, "Conflict-free Replicated Data Types," 2018 survey: https://arxiv.org/abs/1805.06358
- Secure Scuttlebutt documentation, "Secure Scuttlebutt": https://ssbc.github.io/ssb-db/
- Hypercore Protocol documentation, "Hypercore": https://hypercore-protocol.github.io/new-website/guides/modules/hypercore/
