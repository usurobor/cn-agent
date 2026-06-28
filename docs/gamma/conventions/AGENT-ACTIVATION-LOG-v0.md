---
title: Agent Activation Log Convention v0
status: field convention
version: v0
date: 2026-05-30
scope: cross-activation continuity for one agent identity operating across multiple hubs/bodies
related:
  - cn-sigma:.cn-sigma/spec/OPERATOR.md § Activation logs
  - cn-sigma:.cn-sigma/threads/adhoc/20260530-sigma-activation-log-v0.md
  - cn-sigma:.cn-sigma/state/activations.md
  - docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md
  - docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md (cnos#150)
---

# Agent Activation Log Convention v0

A minimal two-artifact, single-writer, append-only convention for cross-activation continuity when one agent identity operates across multiple hubs/bodies.

This is a **field convention** for the topology we have: one agent identity across multiple bodies, one operator owning all push permissions. The whitepaper v1 elaborations (signed + entry-IDed + union-merged) and cnos#150 (ref-based packets) solve problems for a different topology — adversarial routing, distrusted operators, cross-organization peer comms. Evolve toward those elaborations only when the actual topology forces it; do not pre-build for hypothetical adversaries. See [GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md](../../reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md) and [MESSAGE-PACKET-TRANSPORT.md](../../reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md) for what those elaborations look like when needed.

## §0 Writer Locality invariant

**Writer locality (hard).** Every body writes only to its own repo. Cross-repo communication is exclusively read-direction. No body ever pushes, comments, or dispatches with payload to another body's repo.

Consequences:
- Foreign hubs write to `{hub}:.cn-{agent}/logs/` only.
- Home writes to `cn-{agent}:.cn-{agent}/threads/activations/{name}/` and `cn-{agent}:.cn-{agent}/threads/reflections/` only.
- No cross-repo PR comments. No cross-repo `repository_dispatch` with payload. No PATs at foreign hubs for write.
- Notification flows back via pull: home reads, home decides what to surface.

**Single writer per surface** (§6 Trust boundary) is a derived consequence of this invariant: writer locality is the topology; single-writer-per-surface is the file-level enforcement of that topology. A body that would violate writer locality necessarily violates single-writer-per-surface — but the invariant named here is architectural, not file-level.

This invariant prevents a class of design pressure that recurs when notification or relay flows are considered: the instinct is "push notification cross-repo," but writer locality forecloses that direction entirely. Notification flows via pull — home reads the foreign log, home decides what to surface.

**Same-repo terminal-report carve-out.** A body writing to an issue, PR, workflow report, or other notify surface inside its own repo is not a Writer Locality violation. Same-repo terminal reports are the preferred operator-visible notify surface. The forbidden case is writing into another repo's surfaces: cross-repo comments, cross-repo dispatch payloads, or cross-repo pushes.

## §0.1 Wake-class writer ownership (same-repo)

**Same-repo cross-wake exclusion (hard).** When multiple wakes operate in the same repo as different bodies of the same agent identity (e.g., `cnos-agent-admin` and `cnos-cds-dispatch` at `cnos`), the writer-locality invariant of §0 is necessary but NOT sufficient — same-repo writer-locality has historically been "the body that owns this repo writes here," but the wake-orchestration architecture (cnos#467) introduces multiple wakes-as-bodies in the same repo with different ownership. This section carves out the §0 invariant for that case.

**Writer ownership of `.cn-{agent}/logs/` (the channel log surface):**

| Wake class | Writer of `.cn-{agent}/logs/`? | Rationale |
|---|---|---|
| **admin** (e.g. `cnos-agent-admin`) | **Yes — sole writer.** | The admin wake's responsibility is channel sync + status reporting per `cnos.core/orchestrators/agent-admin/wake-provider.json`. Its activation loop (per §3) ends with appending a channel entry. Writer-locality is satisfied: the admin wake at this hub writes to this hub's log. |
| **package-dispatch** (e.g. `cnos-cds-dispatch`, future `cnos-cdr-dispatch`, `cnos-cdw-dispatch`) | **No — non-writer.** | The dispatch wake's responsibility is cell execution per `cnos.cds/orchestrators/cds-dispatch/wake-provider.json`. Its inbound is the open-issue queue, not the home thread; its outbound is cell artifacts (`.cdd/unreleased/{N}/`), cell PRs, and lifecycle label transitions on the claimed cell. **A dispatch wake firing — real claim OR no-op — MUST NOT commit to `.cn-{agent}/logs/`.** |

**Why this partition.** Dispatch firings are package-runtime telemetry (selector scans, no-op exits, claim sequences, drift diagnostics). Channel logs are agent activation memory (durable narrative across activations of the agent identity). Mixing them pollutes activation memory with runtime telemetry at every dispatch firing. The empirical proof-point (2026-06-24, `cn-sigma:.cn-sigma/logs/20260624.md`): the live `cds-dispatch` wake's provider prompt forbade the write in prose; the model wrote channel entries four times anyway (selector-scan no-op entries at ~05:55Z, 08:45Z, 12:38Z, and 13:58Z). Prompt-only prohibition has been empirically falsified.

**Mechanical enforcement (the rule is known; use mechanics).** Cycle/496 (cnos#496 Sub 1 of cnos#495 umbrella) installs the mechanical guard:

1. **Manifest field (`activation_log_writer: bool`).** Every `cn.wake-provider.v1` manifest declares its writer status. Default-when-absent is `true` (backward-compat with existing admin behavior). Package-dispatch wakes (`role: dispatch + admin_only: false`) MUST declare `activation_log_writer: false`.
2. **Render-time refusal (`cn install-wake` exit code 4).** The renderer refuses to materialize a package-dispatch wake whose manifest declares (or defaults to) `activation_log_writer: true` — a mis-declaration. Mis-declared providers do not reach the substrate at all.
3. **Run-time write fence (workflow-level step).** When `activation_log_writer: false`, the renderer appends a final write-fence step to the rendered workflow. The fence inspects this wake's local working-tree state + this run's commit graph for paths under `.cn-{agent}/logs/`; non-empty → fail the workflow with `dispatch_activation_log_write_violation`. **The fence is local-scoped, not remote-state-delta** (false-positive resistance: concurrent admin-wake writes on `main` do not cause the dispatch fence to fail).

This is belt and suspenders: render-time refusal is primary defense (mis-declared providers don't ship); run-time fence is secondary defense (catches drift if a hand-edited workflow bypasses the renderer; catches model writes that bypass the provider prompt).

**Cross-references.**
- **Long-arc partition framing** (cnos#495 umbrella, [comment #4792969173](https://github.com/usurobor/cnos/issues/495#issuecomment-4792969173)): *Use intelligence where meaning is unresolved. Use mechanics where the rule is known. Identity can be Sigma. Ownership is package-local. Memory is summarized upward, not dumped sideways.* This wake-class partition is the first concrete enforcement of that direction.
- **Empirical motivator** (2026-06-24 mixed log): `cn-sigma:.cn-sigma/logs/20260624.md` shows both admin-wake entries (legitimate) and dispatch-wake selector-scan entries (the falsification). The historical entries are preserved per cnos#496 AC6 as the witness for why mechanical enforcement is necessary.
- **Implementation cycle** (cnos#496 Sub 1): renders the field declaration into `cn.wake-provider.v1`; extends `cn install-wake` with exit code 4 + the workflow-level fence; updates `cds-dispatch`'s manifest + prompt; re-renders the production substrate.

**No-op evidence destination.** Dispatch firings that produce no claim (selector scan finds nothing) surface in the GitHub Actions job summary (workflow log + run UI) — NOT in the channel log. The job summary is structured + UI-visible + ephemeral (does not pollute durable activation memory). Cycle/495 Sub 2 (HELD) will land an admin-wake `class: dispatch-summary` channel entry that *reads* the job summaries periodically and writes a single rolled-up summary into the channel; until then, the job summary is sufficient operator-visible surface.

## §1 Agents, activations, and peers

A clear conceptual distinction:

- **Agent:** the coherent identity rooted in a home hub, e.g. `cn-sigma` (engineer/cdd profile) or future `cn-rho` (researcher profile, planned). In path templates below, `{agent}` is the hub slug, e.g. `sigma` or `rho`. This is not the local model process; it is the named continuity that can be activated in different bodies.
- **Activation:** the *same agent identity* taking up residence in *another body*. Sigma-at-cnos, Sigma-at-bumpt, Sigma-at-cph are all the same Sigma — same identity, same operator, same continuity — operating through different repos as bodies. Activations are registered in `cn-{agent}:.cn-{agent}/state/activations.md`.
- **Peer:** a *different agent identity* — its own hub, its own identity, its own keys. E.g., cn-rho would be a peer to cn-sigma. Peers are registered in `cn-{agent}:.cn-{agent}/state/peers.md`. Currently empty; no peer agents exist yet.

This convention covers **activations of one agent identity**. Peer↔peer comms is a different design problem (different identity, separate trust boundary, different lifecycle); deferred until the first peer registers.

## §2 The two artifacts

| Surface | Single writer | Direction |
|---|---|---|
| `{activation-repo}:.cn-{agent}/logs/YYYYMMDD.md` | the agent at that activation | foreign → home |
| `cn-{agent}:.cn-{agent}/threads/activations/{activation}/YYYYMMDD.md` | the agent at home | home → foreign |

Both surfaces are single-writer, append-only, plain markdown. Routing is implicit in the path.

**Both directions are date-sharded. The stream owner differs, but the stream shape is symmetric.** Each side has a per-activation directory of per-day files; the cursor (a Git commit SHA, see §4) spans each directory naturally — no per-file cursor is needed. Sharding is the v0 default not because of volume but because each activation is a long-lived narrative channel, and durable channels are smaller, more mergeable, and more readable when split by day.

The home-to-foreign side stays under `.cn-{agent}/threads/` because these *are* threads — append-only narrative channels — not registry state. The registry is `.cn-{agent}/state/activations.md` (routing, cursors, durable machine-ish state); the conversation per activation is a thread under `.cn-{agent}/threads/activations/{activation}/`. Mixing them would conflate state with conversation; separating them keeps the ontology clean. A top-level `activations/` directory would be justified later only if an activation becomes a multi-surface object (`activations/cnos/{log,state,receipts,directives}/…`); at v0, one date-sharded thread per activation is simpler.

Sigma is the first adopter. Its current instantiation is:

- `cnos:.cn-sigma/logs/YYYYMMDD.md` ← Sigma-at-cnos writes
- `bumpt:.cn-sigma/logs/YYYYMMDD.md` ← Sigma-at-bumpt (bump-sigma) writes
- `cn-sigma:.cn-sigma/threads/activations/cnos/YYYYMMDD.md` ← Sigma-at-home writes
- `cn-sigma:.cn-sigma/threads/activations/bumpt/YYYYMMDD.md` ← Sigma-at-home writes

Both sides read each other's `main` HEAD. Channel artifacts (logs, registry, spec) live on `main` by convention; work-in-progress on feature branches is invisible to the channel until merged. See §6 Branch discipline for the full statement.

## §3 The activation loop

On activation, in order:

1. Pull `cn-{agent}` and the current activation repo.
2. Read the activation's `.cn-{agent}/logs/` from `last_read_foreign_log` (in `cn-{agent}:.cn-{agent}/state/activations.md`) to activation-repo HEAD — walk the directory; any file changed since the cursor is in range.
3. Read `cn-{agent}:.cn-{agent}/threads/activations/{activation}/` for directives from home — walk the directory from the activation-side cursor (the activation scans its own logs for the last `Read home directives through cn-{agent}@{sha}` entry to find that cursor) forward to `cn-{agent}` HEAD.
4. Do the work.
5. Append to the local writer-owned surface (create today's file if it doesn't exist yet):
   - at home, append to `.cn-{agent}/threads/activations/{activation}/YYYYMMDD.md`;
   - at the activation, append to `.cn-{agent}/logs/YYYYMMDD.md`.
6. Commit and push.
7. The home agent updates `last_read_foreign_log` in `.cn-{agent}/state/activations.md`.
8. The activation records the home-read cursor in the entry's frontmatter (`cursor_out: cn-{agent}@<sha>`); the H2 header carries the wake's full UTC timestamp.

## §4 Cursor model

Cursors answer "what's new since I last looked?" Each side keeps one number — a Git commit SHA — marking how far it has read the other side's stream.

| Direction | Cursor location | Written by |
|---|---|---|
| home reads foreign | `cn-{agent}:.cn-{agent}/state/activations.md` → `last_read_foreign_log: <sha>` per activation | the agent at home |
| foreign reads home | `<activation>:.cn-{agent}/logs/YYYYMMDD.md` → entry frontmatter `cursor_out: <agent>@<sha>` (most recent entry's value) | the agent at that activation |

Each side stores its cursor in its **own** surface. Each cursor is a Git commit SHA — Git already addresses "the state of the tree at a point in time," which is why the cursor still works when either side is a directory of per-day files. The same SHA-as-cursor mechanism applies symmetrically on both directions.

**Read direction is asymmetric by purpose.** Cursor extraction from own writer-surface is **bottom-up** (most recent entry first) — the latest entry's `cursor_out` is the active cursor; scanning from the bottom of the most recent file finds it in O(1) without parsing every entry. Processing inbound directives from the reader-surface during sync is **forward-chronological** (earliest first) — directives accumulate causally; the first directive's response sets context for later directives. Same file, two read patterns, both correct for their purpose.

On next activation, each side reads from its cursor forward to the other side's HEAD — only what's new. The cursor advances when the reader records "I've consumed up through here." Idempotent re-activation: if a wake crashes before commit, the cursor hasn't advanced, and the next wake re-reads the same range.

The cursor also serves as an implicit receipt. When home sees the activation-side cursor match the commit SHA where it wrote a directive, that's the ACK — no separate receipt-message needed.

## §5 Entry format

```
## YYYY-MM-DDTHH:MM:SSZ — short subject

---
at: <hub-name>
mode: home | foreign-activation | ephemeral
cursor_in: <agent>@<sha>
cursor_out: <agent>@<sha>
class: substantive | heartbeat | inaugural | directive-out
---

Body. Free-form markdown. Blank line at end.
```

The entry header is a full UTC timestamp (ISO 8601 extended, second precision) plus a short subject. Files are date-sharded (`YYYYMMDD.md`); entry H2 timestamps make multiple wakes on the same day uniquely identifiable and chronologically ordered without inspecting Git history.

Frontmatter is YAML, five required fields:

- **at:** the hub this entry was written from (e.g., `cnos`, `bumpt`, `cn-sigma`). Disambiguates which body authored when reading across hubs.
- **mode:** the attach mode this wake operated in. `home` = home reading foreign logs; `foreign-activation` = body at a registered hub reading home's thread; `ephemeral` = no persistent channel.
- **cursor_in:** the other-side cursor SHA at start of this wake (where the read began).
- **cursor_out:** the other-side cursor SHA at end of this wake. Equal to `cursor_in` when no advance (heartbeat, directive-out, etc.).
- **class:** entry type — `substantive` (real walk with work done), `heartbeat` (no-op cursor advance), `inaugural` (first attach), `directive-out` (outbound directive, no walk).

Body author and timestamp come from Git metadata (`git log <file>`); no need to duplicate in frontmatter.

Multiple entries within a single `YYYYMMDD.md` are bottom-appended in chronological order. Cursor extraction reads `cursor_out` from the most recent (last) entry's frontmatter.

## §6 Trust boundary

- **Single writer per file** — home writes only the home thread; each activation writes only its own local log. This is the file-level enforcement of the §0 Writer Locality invariant.
- **Repo push permission** — GitHub does the authentication via deployed token.
- **Git history** — `git log <file>` answers "did I send this?" with author + commit SHA.
- **Good enough until** volume forces real signing.

### Branch discipline

The channel reads the activation repo's default branch (`main`). Hub-state surfaces (`.cn-{agent}/logs/`, `.cn-{agent}/state/`, `.cn-{agent}/spec/`) are main-only by convention. Project work on feature branches is invisible to home until merged; merging is the act that makes hub state visible to the channel.

This is a clarification of v0 — not an addition to its trust mechanics. It makes explicit what was already implicit: every cursor in the convention is a Git commit SHA on `main`; every log file the next activation walks is the state at `main` HEAD.

The collision case bump-sigma surfaced 2026-06-01: at a fresh project hub where dev work lives on feature branches by default, hub state can be invisible to home until the feature branch merges. The discipline says: keep hub state on `main`; merge to make it channel-visible.

## §7 What is deferred at v0

| Deferred | Why deferred at v0 |
|---|---|
| Signed commits + `cn.json` | Repo push permission is the trust anchor; no impersonation surface at v0 volume. |
| Entry IDs (ULID) | File path + commit SHA + line range serves; only matters under union-merge concurrency. |
| Envelope / frontmatter per entry | Single H2 header is enough; routing is implicit in the file. |
| `merge=union` driver | Single-writer per file → no union needed. |
| CN mail dirs (`inbox/outbox/sent`) | Already explored and dropped — overkill for single-writer logs. |
| Ref-based packets (cnos#150) | Future evolution at higher volume and adversarial routing. |
| Peer↔peer convention | Different agent identities mean a separate trust boundary; design when the first peer lands. |

## §8 Single-writer caveat

Single-writer is **logical**, not physical. Per-day sharding (both directions: foreign `.cn-{agent}/logs/YYYYMMDD.md`, home `.cn-{agent}/threads/activations/{activation}/YYYYMMDD.md`) is the v0 default and handles the common case (one activation per day, or several appended sequentially in the same day's file). If concurrent activations on the same day still race, the next shard is per-activation-hour or per-activation-session. Signatures belong to a different topology — adversarial routing, distrusted operators — not to a more-volume version of this one (`docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md` v1).

## §9 Origin and evolution

- **Origin:** `cn-sigma:spec/OPERATOR.md` § Activation logs (commit `7d8edc0`, 2026-05-30). Sigma is the first adopter; predecessor "Sigma peer log v0" (`89404dd`) was renamed to reflect the activations/peers conceptual split.
- **Generalization:** The Sigma-specific `SIGMA-ACTIVATION-LOG-v0.md` name is superseded by this agent-level convention. The mechanics are unchanged; the scope is now any agent identity with a home hub and foreign activations.
- **Field writeup:** `cn-sigma:threads/adhoc/20260530-sigma-activation-log-v0.md` (rationale + open questions).
- **Different-topology elaborations:** [GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md](../../reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md) v1 (signed envelopes, entry IDs, union merge) and [MESSAGE-PACKET-TRANSPORT.md](../../reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md) (cnos#150, ref-based packets) solve problems for topologies v0 does not have (adversarial routing, distrusted operators, cross-organization peer comms). They are not the next version of v0; they are the right tools when the topology changes. Cross-reference: `cn-sigma:threads/adhoc/20260530-sigma-activation-log-v0.md § Reframe 2026-06-01` for the full topology argument.

## §10 Naming

The canonical name is **Agent Activation Log Convention v0**.

The convention is not Sigma-specific. Sigma is the first adopter and remains the running example until another agent identity activates elsewhere. Future peer-to-peer work should not overload this file; peer comms will need its own convention or promotion into the canonical CN mail protocol.
