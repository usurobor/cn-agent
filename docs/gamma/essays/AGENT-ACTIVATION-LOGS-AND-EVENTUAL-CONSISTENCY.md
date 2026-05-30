---
title: Agent Activation Logs and Eventual Consistency
status: essay
date: 2026-05-30
scope: agent activation logs, eventual consistency, Bayou, session guarantees, future peer transport
governing_question: How should agent activation logs borrow from eventual consistency without becoming a database protocol?
---

# Agent Activation Logs and Eventual Consistency

Agent activation logs are an eventual-consistency protocol for agent memory. They let the same agent act through many bodies without requiring consensus, locks, or a central live process. Each body writes where it has authority, records how far it has read, and trusts later synchronization to converge the narrative.

The design looks new because the actors are agents, but the distributed-systems problem is old. Bayou had mobile computers that could disappear from the network and still needed to read and write useful shared data. Agent activations have foreign bodies that may wake under different harnesses, repos, branches, permissions, and clocks. Both systems reject the same false choice: either stay online and consistent, or stop working.

The activation-log answer is weaker and more useful. Work now. Write locally. Sync later. Preserve enough causal context that the next activation can continue without pretending the world was synchronous.

This is the same kind of move cache-coherence protocols make for processors and anti-entropy makes for mobile databases: do not pretend that every replica is the source of truth at every instant. Give each local actor rules that make later convergence possible.

## Bayou is the nearest ancestor

Bayou is the closest prior art because it was built for disconnected work, not just replicated storage. Its premise was that mobile users should keep reading and writing even when connectivity is poor, intermittent, or expensive. Bayou used read-any/write-any weak replication, propagated updates later, and relied on anti-entropy sessions to move replicas toward convergence.

That is the shape of an agent activation network. Sigma-at-cnos and Sigma-at-home are not always awake together. They do not share one memory process. They share Git, logs, commits, and cursors. When one wakes, it reads what the other wrote since the last cursor. Then it writes its own next entry. Later, the other side does the same.

Bayou's important lesson is not "copy its database." The lesson is that availability can be made safe enough by giving each session a coherent view of its own history. The activation-log spec is a small, Git-shaped version of that idea.

## Session guarantees map to activations

Bayou's four session guarantees name exactly the promises an activation needs.

**Read your writes.** An activation must not lose what it just wrote. In v0, this is local and simple: the activation writes its own log file and commits it. Sigma-at-cnos can immediately read `cnos:.cn-sigma/logs/YYYYMMDD.md`; home Sigma can immediately read `cn-sigma:threads/activations/cnos/YYYYMMDD.md`.

**Monotonic reads.** A reader should not move backward in the stream without saying so. The cursor gives this property. Home records `last_read_foreign_log` for each activation. The foreign activation records "Read home directives through cn-sigma@{sha}." On the next wake, the reader scans forward from that SHA to HEAD. A rewind is allowed only as an explicit repair, not as ordinary operation.

**Writes follow reads.** A response should carry the context it depended on. The foreign activation writes after reading home directives, and its entry records the home commit it consumed. Home writes after reading the foreign log and updates the registry cursor. This makes each response causally legible: the entry says, implicitly or explicitly, "I wrote this after seeing that."

**Monotonic writes.** One writer should not reorder its own history. The v0 spec gets this by making each channel single-writer and append-only. Date-sharded files do not break the promise because the Git commit sequence and H2 entry order still tell the writer's story. If same-day concurrent activations race, the design should shard more narrowly before it reaches for a multi-writer protocol.

These four guarantees are enough for v0 because v0 has one persona and many bodies, not many independent agents arguing over one shared object.

## Git is anti-entropy here

Bayou replicas exchanged missing writes during anti-entropy sessions. Agent activations use Git pull, log scan, patch, commit, and push.

That sounds less formal, but the job is similar. A sleeping activation is a disconnected replica of the agent's working continuity. When it wakes, Git tells it what changed. The cursor tells it what it has already consumed. The log files tell it what the other body meant. The next commit carries its own response back into the shared substrate.

The protocol does not need a daemon to be real. It needs repeatable reads, single-writer append, durable commits, and explicit cursors. Git provides the substrate. The convention provides the memory discipline.

## Agents add semantic conflict resolution

A database conflict resolver works over types, fields, timestamps, and merge procedures. An agent can do something stranger: it can read a conflict as content.

If two directives disagree, a CRDT can pick a deterministic winner or preserve both values. An agent can ask what the directives were trying to accomplish. It can notice that one directive supersedes another, that one was written before a correction, or that both are compatible after reframing. It can resolve by intent, not only by operation type.

This is powerful and dangerous. Semantic conflict resolution is not magic consistency. It is a higher-level repair loop. The system still needs receipts, cursors, author surfaces, and evidence. The agent's interpretation should be recorded as a new entry, not hidden inside an untraceable local decision.

The rule is simple: machines preserve the evidence; agents resolve the meaning.

## Identity attaches to the persona, not the node

Most distributed systems start with nodes. Node A writes. Node B receives. Node C reconciles.

Agent activation logs start with a persona. Sigma is the identity. cnos, bumpt, and cph are bodies. The body matters for authority and local context, but it is not the owner of the memory. The home hub owns the persona. The foreign repo hosts one activation of it.

That changes the routing model. The question is not "which node owns this write?" The question is "which agent identity is speaking, and from which activation body?"

This is why the spec should say `cn-{agent}` rather than naming Sigma in the protocol. Sigma is the first instance. The primitive is agent activation: one coherent agent identity, many execution bodies, one eventual narrative.

## Messages are interpretable content

Activation log entries are Markdown, not opcodes. That is not an accident.

The first version does not need a rigid envelope because the reader is an agent that can interpret prose, headings, file paths, patches, and intent. The H2 boundary gives enough shape. Git gives history. The cursor gives causality. The path gives routing.

This buys speed. It also creates a limit. Markdown works while the trust boundary is friendly, the writers are single, and the messages are meant to be read by the same persona. When peers arrive, the protocol will need more machine-readable fields. When adversarial routing arrives, it will need signatures. When multi-writer state arrives, it will need conflict-free data structures or explicit merge rules.

The v0 format is not under-designed. It is scoped.

## No consensus is needed in v0

Consensus solves a different problem. It makes multiple parties agree on one ordered state under failure. Agent activation v0 does not have that topology.

The v0 topology is:

```text
one agent identity
many activation bodies
single-writer channel per direction
Git history as substrate
commit SHA as cursor
Markdown as meaning layer
```

That topology wants eventual consistency, not consensus. It wants session guarantees, not leader election. It wants a clear owner per stream, not a quorum.

Adding consensus here would make the simple case worse. It would force a live coordination problem onto a system whose main requirement is to survive disconnected wakes.

## Peers change the clock

The cursor model is almost a two-party vector clock, collapsed to one SHA per direction. That collapse is valid because each channel has one writer and one intended reader context.

Peers break the simplification. Once Rho, Sigma, and another agent exchange messages as distinct identities, each party needs to know not only "what did I last read?" but "which events had this sender seen when it wrote?" That is causal metadata, not just a bookmark.

The future peer protocol will likely need vector-clock-shaped state. It may not literally use classic vector clocks, but it will need the same information: per-participant knowledge of what has been observed. Otherwise an agent can answer a message without seeing the message that caused it, and the transcript can look coherent while the causality is false.

## Shared state may need CRDTs

Logs are easy because a single writer appends to each stream. Shared state is harder.

If multiple agents can edit the same registry, status board, or task set, a log convention is not enough. The system needs either one writer, an explicit review gate, or a replicated data type whose merge semantics are known in advance. CRDTs exist for this class of problem: replicas accept updates without coordinating and converge deterministically once they have the same updates.

That does not mean cnos should turn activation logs into CRDTs. It means the system should reserve CRDTs for actual multi-writer state. Do not pay the abstraction tax before the topology requires it.

## Signing changes the authority model

In v0, Git push permission and repository history are the trust anchor. That is enough for friendly single-writer activation logs.

Signing changes the model. Once entries are signed, authority moves from "this repo accepted a commit" toward "this agent key authored this entry." That is the shape used by signed append-only feeds such as Secure Scuttlebutt and Hypercore: entries are chained, signed, and verifiable by readers who know the writer's public key.

That future matters because Git history is repository authority, not agent authority. It tells us what was pushed. It does not prove, by itself, which persona key meant the entry. When the system needs portable identity across repos, adversarial boundaries, or third-party relay, signed entries become the right next layer.

## What the spec should borrow

The activation-log spec should borrow four names from the distributed-systems literature and keep the mechanism small.

First, name the session guarantees. The v0 convention already has read-your-writes, monotonic reads, writes-follow-reads, and monotonic writes in practical form. Naming them prevents future readers from re-deriving Bayou by intuition.

Second, name anti-entropy. Git pull/push and cursor scanning are the anti-entropy loop. The term helps distinguish eventual synchronization from live messaging.

Third, name the degenerate clock. The SHA cursor is enough for one home–foreign pair, but it is not a full peer graph clock. Naming the limit makes the next migration visible.

Fourth, name semantic conflict resolution as an agent-specific layer. Agents can resolve by intent, but they must do so by writing new evidence, not by silently editing history.

## What the spec should not absorb

The normative spec should stay tight. It should not become a distributed-systems essay.

The spec needs one short "Origins and analogues" section that says:

```text
This convention is an eventually consistent activation-memory protocol.
Bayou/session guarantees are the nearest database analogue.
Git pull/push plus cursors act as the anti-entropy loop.
The v0 cursor is sufficient for single-persona activation channels.
Peers, shared multi-writer state, and adversarial routing require later
vector-clock, CRDT, or signed-feed mechanisms.
```

The essay can carry the full argument. The spec should carry the rule.

## The core rule

Do not make activation memory strongly consistent. Make it causally legible.

A foreign activation should know what home directives it has read. Home should know which foreign entries it has consumed. Each side should write only its own stream. Each entry should leave enough context for the next body to continue without guessing.

That is eventual consistency for agents: not all bodies know the same thing at the same time, but every body can eventually learn what matters, preserve the causal trace, and repair meaning in public.

## References

- Douglas B. Terry, Alan J. Demers, Karin Petersen, Mike Spreitzer, Marvin Theimer, and Brent Welch, "Session Guarantees for Weakly Consistent Replicated Data," 1994: https://www.cs.cornell.edu/courses/cs734/2000FA/cached%20papers/SessionGuaranteesPDIS_1.html
- Alan Demers, Karin Petersen, Mike Spreitzer, Douglas Terry, Marvin Theimer, and Brent Welch, "The Bayou Architecture: Support for Data Sharing among Mobile Users," Xerox PARC CSL-94-18, 1994: https://archive.computerhistory.org/resources/access/text/2024/06/102804092-05-0001-acc.pdf
- Douglas B. Terry, Marvin M. Theimer, Karin Petersen, Alan J. Demers, Mike J. Spreitzer, and Carl H. Hauser, "Managing Update Conflicts in Bayou, a Weakly Connected Replicated Storage System," SOSP 1995: https://people.eecs.berkeley.edu/~brewer/cs262b/update-conflicts.pdf
- Leslie Lamport, "Time, Clocks, and the Ordering of Events in a Distributed System," Communications of the ACM, 1978: https://lamport.azurewebsites.net/pubs/time-clocks.pdf
- Nuno Preguiça, Carlos Baquero, and Marc Shapiro, "Conflict-free Replicated Data Types," 2018 survey: https://arxiv.org/abs/1805.06358
- Secure Scuttlebutt documentation, "Secure Scuttlebutt": https://ssbc.github.io/ssb-db/
- Dat Protocol DEP-0002, "Hypercore": https://www.datprotocol.com/deps/0002-hypercore/
