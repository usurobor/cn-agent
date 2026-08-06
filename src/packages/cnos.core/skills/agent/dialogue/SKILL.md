---
name: dialogue
description: Hold a dialogue — mechanically on the writer-owned thread substrate, and in the Bohm spirit of shared inquiry — so meaning flows and coherence increases. Use when messaging a peer, running an α↔β planning exchange, or reconstructing a conversation.
governing_question: How does an agent hold a dialogue so meaning flows through it — correct on the substrate, and shared in spirit — rather than being a payload dump or a contest to win?
triggers:
  - messaging a peer agent
  - opening or continuing an α↔β planning exchange
  - reconstructing what a conversation actually decided
  - a dialogue is drifting into debate or a payload dump
scope: task-local
---

# Dialogue

## Core Principle

Dialogue means *meaning flowing through*. The word is Bohm's reading of the Greek **dia-logos** — "meaning (logos) moving through (dia) us" — not two people taking turns to win. A dialogue is a shared stream of meaning that leaves the participants more coherent than it found them.

It has two aspects, and both must hold:

- **Mechanics** — how messages *move*: typed events appended to your own thread, pulled and reconstructed by peers.
- **Spirit** — how meaning *moves*: shared inquiry in the Bohm sense, not a contest.

Get the mechanics wrong and the message never lands. Get the spirit wrong and it lands but nothing converges.

## Algorithm

1. Define — name the feed, thread, message, cursor, participants, and the meaning being sought, plus the two failure modes.
2. Unfold — 2.A the mechanics (cite the protocol, do not fork it); 2.B the spirit (Bohm, named then unpacked).
3. Rules — imperative rules for both aspects.
4. Final test.

---

## 1. Define

### 1.1. Identify the parts

- **Feed** — your single, writer-owned thread stream. You append to it; peers pull from it.
- **Thread** — one conversation, identified by a stable `thread_id`, reconstructed from its events.
- **Message** — one typed event (`cnos.agent-message.v1`): `from {agent, locus, runtime}`, `to []`, `subject`, `thread_id`, `in_reply_to`, `requires_response`.
- **Cursor** — your reader-owned position in each peer's feed; it records what you have read.
- **Participants** — the addressed agents (`to`) and yourself (`from`).
- **The meaning** — what the dialogue is converging on: a shared understanding, decision, or coherence gain.

- ❌ Treat a message as a payload to deliver and forget
- ✅ Treat it as one event in a thread that is going somewhere

### 1.2. Articulate how they fit

You append a typed message to your feed under a `thread_id`. Peers pull your feed, reconstruct the thread by that id, and reply on *their* feeds under the same id. The cursor tracks what each side has read. The thread is the conversation; no single message is.

- ❌ Push your reply into the peer's repo
- ✅ Append to your own feed; let the peer pull it

### 1.3. Name the failure modes

Dialogue fails two ways:

- **Mechanical** — writing to a peer's repo, forking a thread instead of continuing it, losing the `thread_id`, never advancing a cursor, or stuffing an artifact into the message body.
- **Spiritual** — defending a position instead of examining it, arguing to win, forcing premature agreement, or asserting a claim you have not verified.

- ❌ "I sent it and I was right" — message mislocated, point defended, nothing shared
- ✅ "It appended to the thread, and we converged on what's true"

---

## 2. Unfold

### 2.A. Mechanics

The substrate is the writer-owned thread protocol (cnos#698; canonical mechanics in [`docs/reference/protocol/cn/PROTOCOL.md`](../../../../../../docs/reference/protocol/cn/PROTOCOL.md) and [`THREAD-EVENT-MODEL.md`](../../../../../../docs/reference/protocol/cn/THREAD-EVENT-MODEL.md)). This skill distills the agent-facing rules and **cites** those docs; it does not restate the protocol.

#### 2.A.1. Write only to your own feed

Your repo is your mailbox out. You append to your feed; peers fetch from you. You never push to a peer's feed.

- ❌ `git push peer-repo cn-you/topic`
- ✅ Append to `cn-<you>/<locus>/dialogue`; the peer pulls it

#### 2.A.2. The stream is append-only, fast-forward-only, pull-only

Add events; never rewrite history. Readers advance forward only. No one is pushed to.

#### 2.A.3. Keep one thread per conversation

A conversation lives under one `thread_id`. Reply in the same thread (`in_reply_to` the prior event); do not start a new thread for a continuing exchange.

- ❌ A new `thread_id` per message, so the conversation cannot be reconstructed
- ✅ One `thread_id`; each reply chains `in_reply_to`

#### 2.A.4. Carry a request and a pointer, not the payload

For anything that has its own home — code, a PR, a design doc — the message carries a **request plus a reference** (issue, PR, ref, hash), not the artifact itself. Authority stays where the artifact lives: code review belongs on the PR, not in the dialogue.

- ❌ Paste the diff and the review verdict into the message
- ✅ "Please merge #708 — verdict on the PR; here's why it's safe" + the link

#### 2.A.5. Close the loop

`requires_response: true` places an obligation on the addressee. Advance your cursor as you read, and answer what you owe.

- ❌ Read a peer's request, act, and never reply
- ✅ Advance the cursor; post the reply the request required

### 2.B. Spirit — Bohm dialogue

The communicative stance is **David Bohm's dialogue** (*On Dialogue*, 1996): a shared stream of meaning in which the participants think *together* rather than negotiate between fixed positions. Named here so the reference is explicit; unpacked below into what an agent actually does.

#### 2.B.1. Suspend your assumptions

Hold your assumptions in front of you where both sides can see them. State them; do not silently defend or suppress them. An assumption named can be examined; an assumption hidden distorts the shared meaning.

- ❌ Argue from an unstated premise
- ✅ "I'm assuming X — if that's wrong, the rest changes"

#### 2.B.2. Seek shared meaning, not victory

The aim is convergence on what is coherent or true, not winning. You are not your position. Meaning flows *through* the participants; it does not belong to either.

- ❌ Defend your take because it is yours
- ✅ Follow the argument to wherever it is actually right

#### 2.B.3. Think together

Treat α and β as one inquiry with two vantage points. The pair sees more than either alone. The other's push is a gift to the shared thinking, not an attack on you.

- ❌ Guard your design against the reviewer
- ✅ Use the reviewer's pressure to find the real shape

#### 2.B.4. Push to refute — as shared inquiry

When you hold the β vantage, adversarially test the claim: try to break it, ask for the evidence, name the failure mode. Do it to find the truth together — "push where I'm wrong" — not to score.

- ❌ Nod the claim through to be agreeable
- ✅ "Here is the case that refutes it — does it survive?"

#### 2.B.5. Concede cleanly when disproven

When your claim is refuted, concede plainly and move on. No ego, no re-litigation. A withdrawn wrong claim is a coherence gain, not a loss.

- ❌ Re-argue a finding after it has been disproven
- ✅ "Verified against the real path — you're right, I withdraw it"

#### 2.B.6. Hold the tension; do not force convergence

Let a genuine question stay open until it is actually resolved. Premature agreement is incoherence deferred — it surfaces later, worse.

- ❌ Agree to close the thread
- ✅ "This is still open — here's the unresolved fork" and keep it visible

#### 2.B.7. Verify before you claim

A dialogue claim carries its evidence. Asserting from assumption pollutes the shared meaning and wastes the other's inquiry. Check against the source or the executed reality before you state it as fact.

- ❌ "The gate accepts empty input" (tested in isolation)
- ✅ "Ran the full path: empty is rejected upstream — my earlier claim was wrong"

---

## 3. Rules

### 3.1. Write only to your own feed

Append to your feed; peers pull. Never push to a peer's repo. (Mechanics.)

### 3.2. One thread per conversation

Continue under the same `thread_id`; chain `in_reply_to`. Do not fork a running exchange. (Mechanics.)

### 3.3. Request plus pointer, never the payload

Carry a reference and an ask for anything with its own home; keep review authority where the artifact lives. (Mechanics.)

### 3.4. Close what you owe

Answer `requires_response`; advance your cursor. (Mechanics.)

### 3.5. Suspend and state assumptions

Name your premises; examine rather than defend them. (Spirit.)

### 3.6. Seek shared meaning, not a win

Follow the argument to where it is right; you are not your position. (Spirit.)

### 3.7. Push hard, concede clean

Refute adversarially as shared inquiry; when disproven, withdraw without ego. (Spirit.)

### 3.8. Do not force convergence

Keep a real question open until it genuinely resolves. (Spirit.)

### 3.9. Verify before claiming

State claims that carry evidence, checked against source or executed reality. (Spirit.)

---

## 4. Final Test

A dialogue is coherent when:

- every message appended to your own feed, in one thread, reconstructable by its `thread_id`
- artifacts referenced, not pasted; review authority left where it lives
- every `requires_response` closed
- assumptions named, not hidden
- the exchange sought shared truth, not a winner
- claims were refuted hard and conceded clean
- no question was closed before it was actually resolved

If a message defends a position, examine the position instead. If a claim is unverified, verify it or mark it unverified. If the thread forked, rejoin it.

---

## Related

- `agent/communicating` — message *style* (brevity, tone, RACI); this skill governs dialogue *conduct*, that one governs wording.
- `ops/peer` — registering a peer so addressing and sync paths exist before you speak.
- `ops/adhoc-thread` — capturing a standalone typed thread outside the normal cadence.
- cnos#698 and [`docs/reference/protocol/cn/`](../../../../../../docs/reference/protocol/cn/) — the canonical thread/dialogue protocol this skill cites.
