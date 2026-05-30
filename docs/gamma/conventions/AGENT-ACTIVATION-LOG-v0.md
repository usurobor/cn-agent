---
title: Agent Activation Log Convention v0
status: field convention
version: v0
date: 2026-05-30
scope: cross-activation continuity for one agent identity operating across multiple hubs/bodies
related:
  - cn-sigma:spec/OPERATOR.md § Activation logs
  - cn-sigma:threads/adhoc/20260530-sigma-activation-log-v0.md
  - cn-sigma:state/activations.md
  - docs/alpha/protocol/WHITEPAPER.md
  - docs/alpha/protocol/MESSAGE-PACKET-TRANSPORT.md (cnos#150)
---

# Agent Activation Log Convention v0

A minimal two-artifact, single-writer, append-only convention for cross-activation continuity when one agent identity operates across multiple hubs/bodies.

This is a **field convention**, not the canonical CN mail protocol. The canonical protocol lives at [WHITEPAPER.md](../../alpha/protocol/WHITEPAPER.md) (v1: signed + entry-IDed + union-merged) with the ref-based evolution at [MESSAGE-PACKET-TRANSPORT.md](../../alpha/protocol/MESSAGE-PACKET-TRANSPORT.md) (cnos#150). v0 is the bridge to today's volume.

## §0 Agents, activations, and peers

A clear conceptual distinction:

- **Agent:** the coherent identity rooted in a home hub, e.g. `cn-sigma` (engineer/cdd profile) or future `cn-rho` (researcher profile, planned). In path templates below, `{agent}` is the hub slug, e.g. `sigma` or `rho`. This is not the local model process; it is the named continuity that can be activated in different bodies.
- **Activation:** the *same agent identity* taking up residence in *another body*. Sigma-at-cnos, Sigma-at-bumpt, Sigma-at-cph are all the same Sigma — same identity, same operator, same continuity — operating through different repos as bodies. Activations are registered in `cn-{agent}:state/activations.md`.
- **Peer:** a *different agent identity* — its own hub, its own identity, its own keys. E.g., cn-rho would be a peer to cn-sigma. Peers are registered in `cn-{agent}:state/peers.md`. Currently empty; no peer agents exist yet.

This convention covers **activations of one agent identity**. Peer↔peer comms is a different design problem (different identity, separate trust boundary, different lifecycle); deferred until the first peer registers.

## §1 The two artifacts

| Surface | Single writer | Direction |
|---|---|---|
| `{activation-repo}:.cn-{agent}/logs/YYYYMMDD.md` | the agent at that activation | foreign → home |
| `cn-{agent}:threads/activations/{activation}/YYYYMMDD.md` | the agent at home | home → foreign |

Both surfaces are single-writer, append-only, plain markdown. Routing is implicit in the path.

**Both directions are date-sharded. The stream owner differs, but the stream shape is symmetric.** Each side has a per-activation directory of per-day files; the cursor (a Git commit SHA, see §3) spans each directory naturally — no per-file cursor is needed. Sharding is the v0 default not because of volume but because each activation is a long-lived narrative channel, and durable channels are smaller, more mergeable, and more readable when split by day.

The home-to-foreign side stays under `threads/` because these *are* threads — append-only narrative channels — not registry state. The registry is `state/activations.md` (routing, cursors, durable machine-ish state); the conversation per activation is a thread under `threads/activations/{activation}/`. Mixing them would conflate state with conversation; separating them keeps the ontology clean. A top-level `activations/` directory would be justified later only if an activation becomes a multi-surface object (`activations/cnos/{log,state,receipts,directives}/…`); at v0, one date-sharded thread per activation is simpler.

Sigma is the first adopter. Its current instantiation is:

- `cnos:.cn-sigma/logs/YYYYMMDD.md` ← Sigma-at-cnos writes
- `bumpt:.cn-sigma/logs/YYYYMMDD.md` ← Sigma-at-bumpt (bump-sigma) writes
- `cn-sigma:threads/activations/cnos/YYYYMMDD.md` ← Sigma-at-home writes
- `cn-sigma:threads/activations/bumpt/YYYYMMDD.md` ← Sigma-at-home writes

## §2 The activation loop

On activation, in order:

1. Pull `cn-{agent}` and the current activation repo.
2. Read the activation's `.cn-{agent}/logs/` from `last_read_foreign_log` (in `cn-{agent}:state/activations.md`) to activation-repo HEAD — walk the directory; any file changed since the cursor is in range.
3. Read `cn-{agent}:threads/activations/{activation}/` for directives from home — walk the directory from the activation-side cursor (the activation scans its own logs for the last `Read home directives through cn-{agent}@{sha}` entry to find that cursor) forward to `cn-{agent}` HEAD.
4. Do the work.
5. Append to the local writer-owned surface (create today's file if it doesn't exist yet):
   - at home, append to `threads/activations/{activation}/YYYYMMDD.md`;
   - at the activation, append to `.cn-{agent}/logs/YYYYMMDD.md`.
6. Commit and push.
7. The home agent updates `last_read_foreign_log` in `state/activations.md`.
8. The activation records the home-read cursor inline: `## YYYY-MM-DD — Read home directives through cn-{agent}@{sha}`.

## §3 Cursor model

Cursors answer "what's new since I last looked?" Each side keeps one number — a Git commit SHA — marking how far it has read the other side's stream.

| Direction | Cursor location | Written by |
|---|---|---|
| home reads foreign | `cn-{agent}:state/activations.md` → `last_read_foreign_log: <sha>` per activation | the agent at home |
| foreign reads home | `<activation>:.cn-{agent}/logs/YYYYMMDD.md` → dated entry `Read home directives through cn-{agent}@<sha>` | the agent at that activation |

Each side stores its cursor in its **own** surface. Each cursor is a Git commit SHA — Git already addresses "the state of the tree at a point in time," which is why the cursor still works when either side is a directory of per-day files. The same SHA-as-cursor mechanism applies symmetrically on both directions.

On next activation, each side reads from its cursor forward to the other side's HEAD — only what's new. The cursor advances when the reader records "I've consumed up through here." Idempotent re-activation: if a wake crashes before commit, the cursor hasn't advanced, and the next wake re-reads the same range.

The cursor also serves as an implicit receipt. When home sees the activation-side cursor match the commit SHA where it wrote a directive, that's the ACK — no separate receipt-message needed.

## §4 Entry format

```
## YYYY-MM-DD — short subject

Body. Free-form markdown. Blank line at end.
```

No YAML frontmatter per entry. No `entry_id`. No envelope. The H2 header is the entry boundary; the file's Git history is the trace. Multiple entries land in the same `YYYYMMDD.md` if multiple activations happen on the same day.

## §5 Trust boundary

- **Single writer per file** — home writes only the home thread; each activation writes only its own local log.
- **Repo push permission** — GitHub does the authentication via deployed token.
- **Git history** — `git log <file>` answers "did I send this?" with author + commit SHA.
- **Good enough until** volume forces real signing.

## §6 What is deferred at v0

| Deferred | Why deferred at v0 |
|---|---|
| Signed commits + `cn.json` | Repo push permission is the trust anchor; no impersonation surface at v0 volume. |
| Entry IDs (ULID) | File path + commit SHA + line range serves; only matters under union-merge concurrency. |
| Envelope / frontmatter per entry | Single H2 header is enough; routing is implicit in the file. |
| `merge=union` driver | Single-writer per file → no union needed. |
| CN mail dirs (`inbox/outbox/sent`) | Already explored and dropped — overkill for single-writer logs. |
| Ref-based packets (cnos#150) | Future evolution at higher volume and adversarial routing. |
| Peer↔peer convention | Different agent identities mean a separate trust boundary; design when the first peer lands. |

## §7 Single-writer caveat

Single-writer is **logical**, not physical. Per-day sharding (both directions: foreign `.cn-{agent}/logs/YYYYMMDD.md`, home `threads/activations/{activation}/YYYYMMDD.md`) is the v0 default and handles the common case (one activation per day, or several appended sequentially in the same day's file). If concurrent activations on the same day still race, the next shard is per-activation-hour or per-activation-session. Signatures come later under volume/adversarial pressure (`docs/alpha/protocol/WHITEPAPER.md` v1).

## §8 Origin and evolution

- **Origin:** `cn-sigma:spec/OPERATOR.md` § Activation logs (commit `7d8edc0`, 2026-05-30). Sigma is the first adopter; predecessor "Sigma peer log v0" (`89404dd`) was renamed to reflect the activations/peers conceptual split.
- **Generalization:** The Sigma-specific `SIGMA-ACTIVATION-LOG-v0.md` name is superseded by this agent-level convention. The mechanics are unchanged; the scope is now any agent identity with a home hub and foreign activations.
- **Field writeup:** `cn-sigma:threads/adhoc/20260530-sigma-activation-log-v0.md` (rationale + open questions).
- **Future evolution:** [WHITEPAPER.md](../../alpha/protocol/WHITEPAPER.md) v1 (signed envelopes, entry IDs, union merge) → [MESSAGE-PACKET-TRANSPORT.md](../../alpha/protocol/MESSAGE-PACKET-TRANSPORT.md) (cnos#150, ref-based packets).

## §9 Naming

The canonical name is **Agent Activation Log Convention v0**.

The convention is not Sigma-specific. Sigma is the first adopter and remains the running example until another agent identity activates elsewhere. Future peer-to-peer work should not overload this file; peer comms will need its own convention or promotion into the canonical CN mail protocol.
