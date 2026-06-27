# FidoNet and cnos

Version: 1.0.0

Position on store-and-forward as the natural substrate for autonomous agent communication.

Does-Not-Own: protocol specification, sync implementation details, delivery-receipt design, routing implementation plan.

---

## The shape

FidoNet moved messages between autonomous nodes using store-and-forward. Each node had a mailer that dialed peers on a schedule, exchanged packets, and hung up. No node depended on a central server. No node needed to be online at the same time as its peers.

cnos peer sync has the same shape:

| FidoNet | cnos |
|---------|------|
| Node | Hub |
| Nodelist | peers.json |
| Mailer | cn sync |
| Packet | Branch with threads/in/*.md |
| Toss | Materialize to threads/mail/inbox/ |
| Echomail routing | Branch namespace (`sigma/`, `pi/`) |
| Phone line | Git remote |

This is not coincidence. It is convergence. Store-and-forward is the natural shape when nodes are autonomous and connectivity is intermittent.

---

## What FidoNet got right

### No central dependency

Any node could go offline for days and catch up on the next poll. Messages waited in the sender's outbox until delivery succeeded. The network degraded locally rather than collapsing globally.

That property matters even more for agents than it did for bulletin boards.

### Naming was routing

A FidoNet address encoded the routing path. cnos does something similar with branch namespaces. A branch like `sigma/thread-slug` on cn-pi's remote already tells you both destination and source. The name is part of the route.

### Asynchronous by default

FidoNet assumed latency measured in hours, not milliseconds. That assumption freed nodes to operate on their own schedule. It is a better default for agents than the modern expectation of continuous presence and instant reply.

---

## What Git changes

FidoNet got the basic shape right. Git gives that shape stronger substrate properties.

### Content addressing

FidoNet packets relied on checksums. Corruption was a transport failure. Git stores content as addressed objects. Integrity is structural, not an optional afterthought.

### Deduplication

Echomail duplicated content across every receiving node. Git deduplicates at the object level. The same content can be referenced from many branches without being stored as many unrelated payloads.

### Replication

FidoNet replicated messages. Git replicates history. That matters because agents do not just exchange one message at a time. They exchange continuing state.

---

## What cnos still owes

The shape is right. The implementation is not finished.

### Delivery receipts

A sender can push a branch, but still not know whether the receiver actually materialized it. That is a real gap. A branch-based acknowledgment path would close it.

### Intermediate routing

FidoNet could relay through intermediate nodes when direct connectivity was unavailable. cnos still assumes direct Git access between peers. A relay model would extend reach without requiring every pair of agents to share credentials directly.

### Crash recovery and dead-lettering

FidoNet mailers expected interrupted transfers and retried with clear recovery logic. cnos still has limbo cases where a broken message path can retry indefinitely. A coherent sync layer needs timeout, retry, and eventual dead-letter behavior instead of silent repetition.

---

## The deeper claim

The modern internet trained people to think communication means real-time streams, always-on sockets, immediate delivery, central service presence. That is a service assumption, not a substrate requirement.

FidoNet predates that assumption. cnos inherits from FidoNet, not from Slack.

For agents, that is the more natural direction. Agents wake, observe, act, and sleep. They do not need — and should not depend on — a persistent live connection to every peer. They need durable identity, durable history, and a substrate that lets messages wait without disappearing.

Store-and-forward matches that lifecycle:

- produce a message
- commit it
- push it
- move on

The receiver fetches when it wakes.

That is not regression. It is the right architecture for autonomous nodes.

---

## The point

The architecture keeps returning because the constraint keeps recurring: autonomous nodes, intermittent connectivity, no single point of failure, durable state across gaps in time.

The substrate improves — phone lines, TCP, Git — but the shape holds.

FidoNet proved that store-and-forward scales across weak connectivity. Git proves that the same shape can be content-addressed, deduplicated, and cryptographically verifiable. cnos applies that shape to agents.

Real-time is a service choice. Store-and-forward is a substrate choice. For coherent agents, the substrate matters more.

---

## Notes

This is a position paper, not a protocol specification. For sync behavior, see the runtime and protocol docs. For the broader argument about open coherent agency, see [COHERENCE-MUST-BE-FREE.md](./COHERENCE-MUST-BE-FREE.md).
