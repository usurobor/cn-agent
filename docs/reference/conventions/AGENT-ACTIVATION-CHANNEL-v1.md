---
title: Agent Activation Channel Convention v1
status: field convention
version: v1
date: 2026-08-01
supersedes: docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
scope: cross-activation continuity for one agent identity operating across multiple hubs/bodies — the channel is an independent append-only ref, off `main` HEAD
related:
  - docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md (superseded — retained as the v0 historical record; see that file's own superseded-by note)
  - cnos#684 (this convention's design issue; source of AC1–AC5)
  - cnos#682 ("architecture: dematerialize closed CDD cells from HEAD…" — the sibling *ancestry* plane; out of scope here, referenced only as a migration blast-radius item)
  - .cdd/unreleased/684/dry-run-migration-plan.md (AC7 — dry-run-only migration plan for the existing on-`main` `.cn-sigma/logs/**` history)
  - .cdd/unreleased/684/verify-channel-reconstruction.sh (AC6 — content-preservation verification procedure; precondition gate before any strip step in the migration plan)
  - docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md
  - docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md (cnos#150)
  - src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.12 (cell/substrate identity boundary — this convention is a hub-state/substrate concern; cell-execution cognition does not read or write the channel surfaces this doc describes)
---

# Agent Activation Channel Convention v1

A convention for cross-activation continuity when one agent identity operates across multiple hubs/bodies, where the communication stream is an **independent append-only ref**, structurally separate from `main`'s HEAD tree.

This is v1 of the convention `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` established. v0's topology assumptions (one operator, one agent identity, non-adversarial routing) are unchanged and still hold — see that file for the full "why not the whitepaper v1 elaborations yet" argument, which this document does not repeat. What changes in v1 is **where the stream lives** and **what a registration entry has to declare**, not the topology or the trust model.

## §0 Rename rationale

v0 named the surface `.cn-{agent}/logs/` — a folder inside `main`'s tree. Two problems drove the rename:

1. **"Logs" is too generic.** The folder is specifically the *foreign→home activation stream* — one half of a two-direction channel — not a general-purpose log sink. Calling it "logs" invited exactly the confusion that produced the 2026-07 cleanup friction: readers of `main`'s tree could not tell, from the name alone, that this was an independent communication timeline rather than product state.
2. **The stream is not part of the product tree.** `main` HEAD is supposed to hold only *current, actionable* state (`docs/development/board`'s own regeneration discipline makes the same argument for board state). A parallel communication timeline — narrative, append-only, growing without bound — does not fit that shape. Materializing it in `main`'s tree meant every reader of `main` paid its churn cost (~1024 commits / ~10 MB, the dominant source of `main` HEAD churn per the operator's dispatch-authorization comment) whether or not they cared about the channel.

The fix folds the name into the mechanism: in the target model, **the name is a *ref* name, not a folder name.** Dropping "logs" for "channel" makes the rename observable in the artifact itself — `refs/heads/channels/{agent}/…` — rather than only in prose. There is no `.cn-sigma/logs/` in the target architecture; there is a ref.

## §1 Agents, activations, and peers

Unchanged from v0 §1 — restated here for a self-contained v1 read:

- **Agent:** the coherent identity rooted in a home hub, e.g. `cn-sigma` (engineer/cdd profile) or a future `cn-rho` (researcher profile). `{agent}` in path/ref templates below is the hub slug (`sigma`, `rho`, …). This is the named continuity, not the local model process.
- **Activation:** the *same agent identity* taking up residence in *another* body/repo. Sigma-at-cnos, Sigma-at-bumpt, Sigma-at-cph are all the same Sigma — same identity, same operator, same continuity — through different repos as bodies.
- **Peer:** a *different* agent identity, its own hub, its own identity, its own keys. Peer↔peer comms remains a deferred, separate design problem (v0 §1's scoping still applies; no peer agent exists yet as of this writing).

## §2 The two orphan refs

v0's two artifacts were files in each side's `main` tree (`.cn-{agent}/logs/YYYYMMDD.md` at the activation, `.cn-{agent}/threads/activations/{activation}/YYYYMMDD.md` at home). v1 replaces both with **orphan refs** — branches with no shared history with `main`, holding the same append-only, date-sharded content, but structurally outside the product tree.

| Direction | Ref | Single writer | Reader |
|---|---|---|---|
| foreign → home | `refs/heads/channels/{agent}/{activation}-to-home` | the agent at that activation (e.g. cnos-side Sigma writes `refs/heads/channels/sigma/cnos-to-home`) | the agent at home |
| home → foreign | `refs/heads/channels/{agent}/home-to-{activation}` | the agent at home (writes `refs/heads/channels/sigma/home-to-cnos`) | the agent at that activation |

Sigma's concrete instantiation of the two refs:

- `cnos:refs/heads/channels/sigma/cnos-to-home` ← Sigma-at-cnos is the single writer; Sigma-at-home is the reader.
- `cn-sigma:refs/heads/channels/sigma/home-to-cnos` ← Sigma-at-home is the single writer; Sigma-at-cnos is the reader.

Both refs are **orphaned from `main`** — no merge-base, no shared history — precisely so that neither ref's churn ever appears in a `main`-tree diff, log, or blame. Each ref is still date-sharded internally (one commit or one file-append per wake, per v0's original sharding rationale — durable narrative channels are smaller and more reviewable split by day); the sharding discipline itself is unchanged, only the parent ref moved off `main`.

**Own stream: write + push. Other stream: fetch + poll only.** This restates v0 §0's Writer Locality invariant at the ref level: a body may only ever push commits onto the ref it is the declared writer of. It may fetch and read the other ref freely, but a push (of any kind — normal or force) to the other side's ref is a protocol violation regardless of content.

## §3 Registration schema

Registration state per activation needs **repo + ref + cursor**, across both directions — a bare commit SHA is no longer sufficient once the stream does not share `main`'s tree (there is no longer an implicit "current repo, current branch" context a bare SHA could borrow). AC1's oracle: a registration entry carrying only a bare SHA, with no `repo:`/`ref:` context, is **invalid** and MUST be rejected by any validating tool or reviewer.

The schema is split into two parts, per the pre-dispatch cursor-ownership amendment (cnos#684, operator comment @2026-08-01T06:59:04Z), which supersedes the single combined example in the issue's Proposal §2:

### §3.1 Endpoints (centrally registered)

*Endpoints* — repo + ref for each direction — may live in one shared registry. Endpoints are not secrets and not writer-sensitive; declaring "where the channel lives" centrally does not violate Writer Locality, because Writer Locality governs *who may push commits*, not *who may know a ref's coordinates*.

```yaml
activation: cnos
channels:
  foreign_to_home:
    repo: usurobor/cnos
    ref:  refs/heads/channels/sigma/cnos-to-home
  home_to_foreign:
    repo: usurobor/cn-sigma
    ref:  refs/heads/channels/sigma/home-to-cnos
```

Note deliberately: **no `cursor:` field appears in this shared block, on either the `foreign_to_home` or `home_to_foreign` entry.** This is the amended shape — the inline `foreign:{…cursor}` / `home:{…cursor}` example in the issue's original Proposal §2 is superseded and MUST NOT be implemented as written; a shared registry entry that carries a cursor field is itself the schema-level version of the WRITER_LOCALITY_VIOLATION case named in §3.3 below.

### §3.2 Read cursors (writer-local)

Each side stores **its own** read cursor — how far it has read the *other* side's stream — on a surface **it**, not the other side, owns and writes:

| Cursor | Owner (stores it) | Location |
|---|---|---|
| home's read-progress through `foreign_to_home` | home (`cn-sigma`) | `cn-sigma:.cn-sigma/state/activations.md` → per-activation `last_read_foreign_channel: { repo: usurobor/cnos, ref: refs/heads/channels/sigma/cnos-to-home, cursor: <sha> }` |
| foreign's read-progress through `home_to_foreign` | foreign (the activation, e.g. `cnos`) | the activation's own writer-owned stream state — the most recent entry's frontmatter on `refs/heads/channels/sigma/cnos-to-home` carries `cursor_out: { repo: usurobor/cn-sigma, ref: refs/heads/channels/sigma/home-to-cnos, cursor: <sha> }` |

Every cursor record carries `repo:` + `ref:` alongside the `<sha>` — the SHA alone does not identify *which* ref it addresses once refs are no longer implicitly "the current repo's `main`." A cursor entry with a bare SHA and no `repo:`/`ref:` pair fails validation (AC1's falsifiable check).

### §3.3 WRITER_LOCALITY_VIOLATION (negative case)

**Trigger.** A side attempts to write or update the read-cursor record that lives on the *other* side's owned surface — e.g. the foreign activation attempting to write into `cn-sigma:.cn-sigma/state/activations.md`'s `last_read_foreign_channel`, or home attempting to write a `cursor_out` value into the foreign activation's own stream entries on its behalf.

**Response.** Any tool, CI check, or reviewing agent that detects this pattern MUST reject the write and surface `WRITER_LOCALITY_VIOLATION` as the named error/rejection state. This is the cursor-level instance of the same Writer Locality invariant that governs ref pushes in §2 — the invariant is uniform across "which ref may I push commits to" and "whose cursor record may I write."

## §4 Attach-contract sequence

The attach loop changes from v0's "commit + push to `main`" shape to a fetch/write-and-fast-forward shape against the orphan refs:

1. **Fetch declared reader ref.** Pull the other side's writer-owned ref (the one this side is the registered *reader* of) to local.
2. **Read cursor → reader-tip.** From this side's own stored read cursor (§3.2) forward to the fetched ref's current tip, walk the new commits/entries — this is what's new since last read.
3. **Do the work.**
4. **Append to writer stream.** Append the new entry to *this side's own* writer-owned ref (never the other side's ref).
5. **Commit + fast-forward writer ref.** Commit locally, then push as a fast-forward-only update to this side's writer ref.
6. **Record new reader cursor.** Update this side's own stored read-cursor record (§3.2) to the SHA just read through, on the ref pair just consumed.

**Own stream: write + push. Other stream: fetch + poll only.** Two invariants fall directly out of this sequence and are stated explicitly because they are the exact shape of mistake the v0→v1 migration exists to foreclose:

- **A writer never pushes the other stream.** Step 4/5 only ever target this side's own writer ref. There is no step in this sequence in which a side pushes a commit onto the ref it is registered as *reader* of.
- **A reader never writes.** Step 1/2 are fetch + read only. Consuming the other side's stream never produces a commit on that ref.

Per direction, concretely:

| Channel | Writer (executes steps 4–6) | Reader (executes steps 1–2, consumes) |
|---|---|---|
| `refs/heads/channels/sigma/cnos-to-home` | Sigma-at-cnos | Sigma-at-home |
| `refs/heads/channels/sigma/home-to-cnos` | Sigma-at-home | Sigma-at-cnos |

## §5 Orphan-ref invariants and enforcement

Six invariants govern every channel ref (Proposal §4 of cnos#684), each with a named enforcement mechanism concrete enough to configure without re-deriving the rule from prose:

| # | Invariant | Enforcement mechanism |
|---|---|---|
| 1 | **Orphaned start** — the ref has no merge-base with `main`; it starts life via `git commit --orphan` (or equivalent), never `git branch <ref> main`. | Repo-side CI check on ref creation/first-push: `git merge-base <ref> main` must fail (no common ancestor) or the push is rejected as malformed channel history. |
| 2 | **Exactly one declared writer** | GitHub branch-protection "restrict who can push" rule on `refs/heads/channels/{agent}/*`, scoped to the single writer identity's deploy token/bot account per direction; any other identity's push is rejected at the platform layer before it reaches the ref. |
| 3 | **Fast-forward only** | GitHub branch-protection "require linear history" + reject non-fast-forward updates, enabled per-ref (or per-pattern `channels/{agent}/*`) — a non-ff push (including any push whose new tip does not descend from the ref's current tip) is rejected by GitHub itself, not by convention. |
| 4 | **Force-push prohibited** | Same branch-protection surface as #3 — "allow force pushes" left disabled for `channels/{agent}/*`. A force-push attempt is rejected identically to a non-ff push; there is no separate force-push allowance to configure away. |
| 5 | **Deletion prohibited while registered** | Branch-protection "restrict deletions" enabled on `channels/{agent}/*`, combined with a registration-state check (the endpoint entry in §3.1's shared registry) — a ref MUST NOT be deleted while a `channels:` entry in the shared registry still names it. Deletion is only valid after the corresponding registry entry is first removed (de-registration precedes ref deletion, never the reverse). |
| 6 | **Cursor commit reachable from declared ref** | A CI/validation check run at read time (or periodically): for every stored cursor record (§3.2), `git merge-base --is-ancestor <cursor-sha> <ref>` must succeed. A cursor pointing to a SHA not reachable from the current ref tip indicates either a corrupted cursor record or an out-of-band history change on the ref and MUST be surfaced as a hard failure, not silently re-derived. |

Invariants 2–5 are GitHub branch-protection configuration (no custom tooling required — "normal GitHub fetch, CI, and branch protection all work" on `refs/heads/…`, per Proposal §6's substrate choice). Invariants 1 and 6 require a small validation script/CI check, since GitHub's native branch protection has no built-in "must be orphaned" or "cursor must be reachable" rule.

## §6 Promotion boundary

**A channel observation becomes project memory or project authority only through promotion.** The channel entry — on either `refs/heads/channels/{agent}/cnos-to-home` or `refs/heads/channels/{agent}/home-to-cnos` — stays the communication record: a durable, readable trace of what was observed or said, when. It does **not**, by merely existing on the channel, become project memory. Promotion is the act of carrying an observation into a main-reachable artifact: an issue, a CDD contract or design doc, an explicit decision record, or a code change landed on `main`.

**Binding statement.** Absent promotion, a channel entry MUST NOT silently govern current behavior. This is Kernel §2.1's no-silent-drops principle applied to the channel: the failure mode it forecloses is a future reader (agent or operator) treating "I said this in a channel entry once" as equivalent to "this is now project policy" without the intervening step of landing that policy somewhere `main`'s ancestry can be asked about. If a channel observation should govern behavior, the correct next action is to promote it — file the issue, write the design section, land the code — not to let the channel entry stand in for that artifact.

This section is unchanged in substance from v0's Proposal-adjacent framing (v0 did not yet have this section explicitly, since v0 predates the promotion-boundary requirement named in cnos#684); v1 makes it an explicit, separately-headed rule rather than an implicit assumption.

## §7 Entry format

Unchanged from v0 §5 in shape; cursor fields now carry the `repo:`/`ref:`-qualified form from §3.2 instead of a bare `<agent>@<sha>`:

```
## YYYY-MM-DDTHH:MM:SSZ — short subject

---
at: <hub-name>
mode: home | foreign-activation | ephemeral
cursor_in: { repo: <repo>, ref: <ref>, sha: <sha> }
cursor_out: { repo: <repo>, ref: <ref>, sha: <sha> }
class: substantive | heartbeat | inaugural | directive-out
---

Body. Free-form markdown. Blank line at end.
```

Frontmatter fields (five required, same as v0):

- **at:** the hub this entry was written from.
- **mode:** the attach mode this wake operated in (`home`, `foreign-activation`, `ephemeral`).
- **cursor_in:** the other-side cursor at start of this wake (where the read began) — the qualified `repo`/`ref`/`sha` form, per §3.2.
- **cursor_out:** the other-side cursor at end of this wake (equal to `cursor_in` when no advance).
- **class:** `substantive` | `heartbeat` | `inaugural` | `directive-out`.

Multiple entries within a shard are bottom-appended in chronological order; cursor extraction reads `cursor_out` from the most recent entry.

## §8 Rename status

**This section exists precisely because AC5's full text ("`.cn-sigma` HEAD retains only current cursors/state + the convention README, not the dated stream") cannot be fully satisfied inside this design cycle** — δ confirmed this split at dispatch (`.cdd/unreleased/684/gamma-clarification.md` §1); it is not a silent redefinition.

- **Design-landed (this cycle, cnos#684):** the rename (§0), the two orphan-ref names and their writer/reader roles (§2), the amended registration schema (§3), the attach-contract sequence (§4), the orphan-ref invariants and enforcement (§5), and the promotion boundary (§6) are all fully specified in this document. Any *new* activation-channel entry authored going forward should target the orphan refs under this convention, not `.cn-sigma/logs/YYYYMMDD.md` under `main`.
- **Physical-strip deferred (operator-executed, AC7):** the physical fact that `.cn-sigma` HEAD on `main` retains *only* current cursors/state (not the already-committed dated-stream history under `.cn-sigma/logs/**`) requires removing already-landed content from `main`'s tree via a history rewrite (`git rm` of a live path plus ancestor rewriting to actually shrink `main`'s history). The binding worker guardrails for this cycle (cnos#684 operator comment @2026-08-01T07:50:17Z) explicitly forbid the worker from performing that action: **no history rewrite, no `main` force-push, no ref deletion, no `git rm` of live paths on `main`.** That physical step is planned — not executed — in `.cdd/unreleased/684/dry-run-migration-plan.md`, and is explicitly named there as an operator-executed, separately-gated step, preceded by the content-preservation verification in `.cdd/unreleased/684/verify-channel-reconstruction.sh` (AC6).

Until the physical-strip step runs, `.cn-sigma/logs/**` remains present on `main` HEAD (as historical content only — no new entries should be appended there under this convention). This convention document is authoritative for the *target* shape; `main`'s tree catches up to it via the separately-gated migration.

## §9 Migration references (AC6/AC7)

The go-forward mechanism above (§0–§6) is independent of the migration of already-committed history. Two documents under `.cdd/unreleased/684/` cover that migration, in required execution order:

1. `.cdd/unreleased/684/verify-channel-reconstruction.sh` — the content-preservation verification procedure (AC6). Confirms every `.cn-sigma/logs/**` payload in `main`'s history is reconstructable from the orphan-ref import before any strip is proposed. **Specified, not executed, by this cell** (see that file's own header for why — `.cn-sigma/**` is outside this cell's reach per `delta/SKILL.md` §9.12 and this checkout's sparse-checkout exclusion).
2. `.cdd/unreleased/684/dry-run-migration-plan.md` — the dry-run-only `git filter-repo`-style rewrite plan (AC7), gated on (1) passing first, and gated overall on operator execution — no command in that plan has been run by this cell.

## §10 What changed from v0 (summary)

| Surface | v0 | v1 |
|---|---|---|
| Storage | `.cn-{agent}/logs/YYYYMMDD.md`, `.cn-{agent}/threads/activations/{activation}/YYYYMMDD.md` — files in each side's `main` tree | `refs/heads/channels/{agent}/{activation}-to-home`, `refs/heads/channels/{agent}/home-to-{activation}` — orphan refs, off `main` |
| Cursor | bare Git commit SHA | `{repo, ref, sha}` triple, writer-local (§3.2) |
| Registration | implicit (path convention only) | explicit shared endpoint registry (repo+ref, §3.1) + writer-local cursor records (§3.2) |
| Attach contract | commit + push to `main` | fetch reader ref / write-and-fast-forward writer ref (§4) |
| Ref invariants | none named (relied on §6 "Branch discipline: main-only by convention") | six explicit invariants with named enforcement mechanisms (§5) |
| Promotion boundary | not explicitly named | explicit binding section (§6) |
| Writer Locality (§0 of v0) | repo-level ("every body writes only to its own repo") | unchanged at the repo level; extended to ref-level (§2) and cursor-level (§3.3, WRITER_LOCALITY_VIOLATION) |
| Wake-class writer ownership (§0.1 of v0) | unchanged — still governs which *wake* at a hub may write the channel surface; orthogonal to which *ref* that surface lives on |

## §11 What remains deferred

Carried forward from v0 §7, still deferred at v1 for the same reasons (topology has not changed — one operator, one agent identity, non-adversarial routing):

| Deferred | Why deferred |
|---|---|
| Signed commits + `cn.json` | Repo push permission (now: branch-protection-scoped push permission per ref) remains the trust anchor. |
| Entry IDs (ULID) | Ref + commit SHA + line range still serves; only matters under union-merge concurrency. |
| `merge=union` driver | Single-writer per ref → no union needed. |
| Custom ref namespace beyond `refs/heads/…` | Deferred per cnos#684 Scope — a protected orphan branch under `refs/heads/…` is the first-implementation substrate; a custom ref namespace can follow once tooling proves consistent support (Proposal §6). |
| Peer↔peer convention | Different agent identities mean a separate trust boundary; design when the first peer lands. |

## §12 Origin and evolution

- **v0 origin:** `cn-sigma:spec/OPERATOR.md` § Activation logs (commit `7d8edc0`, 2026-05-30). See `AGENT-ACTIVATION-LOG-v0.md` §9 for the full v0 lineage.
- **v1 origin:** cnos#684 ("architecture: Sigma activation channel — rename `.cn-sigma/logs`, move to symmetric append-only orphan refs"), design-first/explore cell, dispatched under operator authorization 2026-08-01. Cursor-ownership amendment per operator comment @2026-08-01T06:59:04Z; bounded migration scope (AC6/AC7) per operator comment @2026-08-01T07:50:17Z.
- **Relationship to #682:** cnos#682 (dematerializing closed CDD cells from `main` HEAD — the *ancestry* plane) is a sibling design, not a dependency of this one. This convention's migration plan (§9, AC7) references #682's CDD-by-SHA custody mechanism only as a blast-radius item to account for — it does not modify #682's scope.
- **Different-topology elaborations:** unchanged from v0 — `GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md` v1 and `MESSAGE-PACKET-TRANSPORT.md` (cnos#150) remain the right tools for a different topology (adversarial routing, distrusted operators, cross-organization peer comms), not the next version of this convention.

## §13 Naming

The canonical name is **Agent Activation Channel Convention v1**. Sigma remains the first and, as of this writing, only adopter, and remains the running example. The convention itself is agent-identity-agnostic — any future agent identity with a home hub and foreign activations uses the same `refs/heads/channels/{agent}/*` shape.
