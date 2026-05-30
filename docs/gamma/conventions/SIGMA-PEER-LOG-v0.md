---
title: Sigma Peer Log Convention v0
status: field convention
version: v0
date: 2026-05-30
scope: cross-activation continuity for one persona operating across multiple hubs
related:
  - cn-sigma:spec/OPERATOR.md § Foreign activation peer logs
  - cn-sigma:threads/adhoc/20260530-sigma-peer-log-v0.md
  - docs/alpha/protocol/WHITEPAPER.md
  - docs/alpha/protocol/MESSAGE-PACKET-TRANSPORT.md (cnos#150)
---

# Sigma Peer Log Convention v0

A minimal two-artifact, single-writer, append-only convention for cross-activation continuity when one persona (currently Sigma) operates across multiple hubs.

This is a **field convention**, not the canonical CN mail protocol. The canonical protocol lives at [WHITEPAPER.md](../../alpha/protocol/WHITEPAPER.md) (v1, signed + entry-IDed + union-merged) with the ref-based evolution at [MESSAGE-PACKET-TRANSPORT.md](../../alpha/protocol/MESSAGE-PACKET-TRANSPORT.md) (cnos#150). v0 is the bridge to today's volume.

## §1 The two artifacts

| Surface | Single writer | Direction |
|---|---|---|
| `{foreign-hub}:.cn-{persona}/logs/YYYYMMDD.md` | Sigma at that foreign hub | foreign → home |
| `cn-{persona}:threads/peers/{peer}.md` | Sigma at home | home → foreign |

Both surfaces are single-writer, append-only, plain markdown. Routing is implicit in the path.

The foreign side **shards per day** as the default form. Each day's activations append to that day's file; a new day creates a new file. The directory `.cn-{persona}/logs/` carries a `README.md` describing the convention; per-day files carry only entries. The cursor (a Git commit SHA, see §3) spans the directory naturally — no per-file cursor is needed.

The home side is a **single file per peer** (`cn-{persona}:threads/peers/{peer}.md`). Home-to-foreign volume is low; sharding has no economy.

This asymmetry is deliberate: the side that grows is sharded, the side that doesn't isn't.

For Sigma's three current activations:

- `cnos:.cn-sigma/logs/YYYYMMDD.md` ← Sigma-at-cnos writes
- `bumpt:.cn-sigma/logs/YYYYMMDD.md` ← Sigma-at-bumpt (bump-sigma) writes
- `cn-sigma:threads/peers/cnos.md` ← Sigma-at-home writes
- `cn-sigma:threads/peers/bumpt.md` ← Sigma-at-home writes

## §2 The activation loop

On activation, in order:

1. Pull `cn-sigma` and the current peer repo.
2. Read peer's `.cn-sigma/logs/` from `last_read_foreign_log` (in `cn-sigma:state/peers.md`) to peer HEAD — walk the directory; any file changed since the cursor is in range.
3. Read `cn-sigma:threads/peers/{peer}.md` for directives from home (foreign Sigma scans its own logs for the last `Read home directives through cn-sigma@{sha}` entry to find its own cursor).
4. Do the work.
5. Append to the local writer-owned surface:
   - at home, append to `threads/peers/{peer}.md`;
   - at the peer, append to `.cn-sigma/logs/YYYYMMDD.md` (today's file; create it if it doesn't exist yet).
6. Commit and push.
7. Home Sigma updates `last_read_foreign_log` in `state/peers.md`.
8. Foreign Sigma records home-read cursor inline: `## YYYY-MM-DD — Read home directives through cn-sigma@{sha}`.

## §3 Cursor model

Cursors answer "what's new since I last looked?" Each side keeps one number — a Git commit SHA — marking how far it has read the other side's stream.

| Direction | Cursor location | Written by |
|---|---|---|
| home reads foreign | `cn-sigma:state/peers.md` → `last_read_foreign_log: <sha>` per peer | home Sigma |
| foreign reads home | `<foreign>:.cn-sigma/logs/YYYYMMDD.md` → dated entry `Read home directives through cn-sigma@<sha>` | foreign Sigma |

Each side stores its cursor in its **own** surface. Each cursor is a Git commit SHA — Git already addresses "the state of the tree at a point in time," which is why the cursor still works when the foreign side is a directory of per-day files.

On next activation, each side reads from its cursor forward to the other side's HEAD — only what's new. The cursor advances when the reader records "I've consumed up through here." Idempotent re-activation: if a wake crashes before commit, the cursor hasn't advanced, and the next wake re-reads the same range.

The cursor also serves as an implicit receipt. When home sees the foreign cursor match the commit SHA where it wrote a directive, that's the ACK — no separate receipt-message needed.

## §4 Entry format

```
## YYYY-MM-DD — short subject

Body. Free-form markdown. Blank line at end.
```

No YAML frontmatter per entry. No `entry_id`. No envelope. The H2 header is the entry boundary; the file's Git history is the trace. Multiple entries land in the same `YYYYMMDD.md` if multiple activations happen on the same day.

## §5 Trust boundary

- **Single writer per file** — no impersonation surface in-hub.
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

## §7 Single-writer caveat

Single-writer is **logical**, not physical. Per-day sharding (`.cn-sigma/logs/YYYYMMDD.md`) is the v0 default and handles the common case (one activation per persona per day, or several activations on the same day appended sequentially). If concurrent activations on the same day still race, the next shard is per-activation or per-hour. Signatures come later under volume/adversarial pressure (`docs/alpha/protocol/WHITEPAPER.md` v1).

## §8 Origin and evolution

- **Origin:** `cn-sigma:spec/OPERATOR.md` § Foreign activation peer logs (commit `89404dd`, 2026-05-30). Sigma is the first adopter.
- **Field writeup:** `cn-sigma:threads/adhoc/20260530-sigma-peer-log-v0.md` (rationale + open questions).
- **Future evolution:** [WHITEPAPER.md](../../alpha/protocol/WHITEPAPER.md) v1 (signed envelopes, entry IDs, union merge) → [MESSAGE-PACKET-TRANSPORT.md](../../alpha/protocol/MESSAGE-PACKET-TRANSPORT.md) (cnos#150, ref-based packets).

## §9 Naming

"Sigma peer log" while Sigma is the only adopter. When a second persona adopts the same convention (e.g., cn-rho, cn-pi), generalize to "peer log convention v0" or whatever name reflects the broader use.
