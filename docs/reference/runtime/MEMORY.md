# Memory in cnos — ranked model

**Version:** 0.3.0
**Status:** Canonical doctrine, ratified by #690 (operator clarification, 2026-08-02). The ranked model below is the authoritative agent-memory doctrine for cnos going forward. The runtime implementation of the box topology (per-locus write-local refs, home's rollup tower, cursor tracking in `state/cursors.yaml`) does **not** exist yet — that is #690 Sub 2 (cn-sigma dry-run migration map) through Sub 5 (cmp alignment), future work. This document states the doctrine that implementation must converge on, not a description of what already runs.
**Supersedes:** the three-surface split (episodic memory, reflective memory, working continuity, bound respectively to `threads/adhoc/`, `threads/reflections/`, `state/conversation.json`) that was v0.2.0 of this document is **retired**, not relabeled. That split does not survive in any form below: there is one memory primitive — an append-only thread of typed entries — differentiated only by rank, not by three fixed named classes bound to three fixed paths.
**Also supersedes (for agent-memory purposes only):** [`docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md`](../conventions/AGENT-ACTIVATION-LOG-v0.md) is no longer the mechanism agent memory persists through. See that document's own status note for the narrowed scope it retains (writer-locality and wake-class-ownership mechanics, which remain live and are orthogonal to the rank question this document answers).
**Prior exploration, subsumed (not controlling):** #684 and its held PR #688 explored a direction-based ref scheme (`refs/heads/channels/sigma/{cnos-to-home,home-to-cnos}`) as the r0 substrate. Per #690's operator clarification (2026-08-02): "#684/#688 are prior substrate exploration, not the controlling topology; their salvage... carries into Sub 2. #684/#688 closed as subsumed." The salvaged artifacts (`dry-run-migration-plan.md`, `verify-channel-reconstruction.sh`) carry forward into #690 Sub 2 on the box topology stated below — the direction-based ref scheme itself does not.

---

## Principle

`main` holds only what the project **is now**; how-we-got-here lives off-HEAD (`DOCUMENTATION-SYSTEM.md §5`, kernel §2.1). Applied to agent memory:

> **Memory is one primitive — an append-only thread of typed entries — at different *ranks*.** `r0` = raw evidence; `rN` (N≥1) = a compaction that *reads* rank N-1 and cites it. "events vs reflections" and "adhoc vs daily" are the same entry at different ranks/cadences, not separate subsystems.

There is no separate episodic/reflective/working-continuity ontology and no fixed binding of memory classes to fixed paths. A single entry schema (below) covers raw evidence and every compaction over it; what differs between entries is rank and provenance, not kind.

And at **home**, the hub's content is **not** hidden inside a `.cn-{agent}/` dotdir. That container prefix is a *foreign-vendoring* concept only (keeps Sigma's files from colliding with a host repo). At home — the cn-sigma repo, which *is* Sigma — the hub materializes at repo **root**.

## The model (KISS / YAGNI)

1. **r0 — distributed, write-local.** Each activation appends to its **own box** (an orphan ref) **at its own repo**, where it already has push access. One box per locus. **No mirror, no cross-repo copy** — each box is the single source of truth for that locus's raw evidence.
2. **rN — centralized at home.** Home is the **only** reader-across and the **only** compactor: it fetches every box by `(repo, ref, cursor)` and writes the rollup tower — daily **r1** over r0, weekly **r2** over r1, monthly **r3** over r2. A reflection is a **new, smaller artifact**, never a copy.
3. **Asymmetry.** Activations are dumb producers (append-only, local, never read others, never compact); home is the one synthesizer. **r0 fans out; rN funnels in.** One Sigma identity → one home (cn-sigma) → one tower; every box (cnos, bumpt, tsc, cmp, *and* `sigma/home`) feeds it.
4. **Cursors = state.** Home tracks how far it has read per box in `state/cursors.yaml` — one line each `{repo, ref, last_read_sha}`, nothing else.
5. **Provenance.** Each rollup names the SHAs it read (`reads:`). The one non-negotiable field — it makes a bad reflection repairable from raw instead of a summary that drifts.
6. **Promotion ≠ rank.** Moving a stable lesson into spec/state/protocol is a different event, not `r(N+1)`. The reflective tower stays one kind of artifact (`r0 → r1 → r2 → …`); promotion moves a stable lesson out of that tower into a governing surface (spec/state/protocol) as an explicit, cited edit — never by incrementing rank. See "Promotion is not a rank" in [`docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md`](../../papers/AGENT-MEMORY-LOG-STRUCTURED.md) for the full argument.

## Ref topology (box model; authoritative go-forward shape, per #690)

```
usurobor/cnos      : refs/heads/sigma/<activation-id>        Sigma-at-cnos  r0
usurobor/bumpt     : refs/heads/sigma/<activation-id>        Sigma-at-bumpt r0
usurobor/tsc       : refs/heads/sigma/<activation-id>        Sigma-at-tsc   r0
usurobor/cmp       : refs/heads/sigma/box, refs/.../cloud    two activations, one repo
usurobor/cn-sigma  : refs/heads/sigma/home                   home's own r0 box (off HEAD)
```

Each box: `README.md` + `posts/YYYYMMDD.md`. Invariants: orphan (no `main` ancestry), **single writer**, **fast-forward-only**, no force-push, no-delete-while-registered.

This topology is **doctrine** — the ratified go-forward shape — not yet-implemented runtime state. Migrating cn-sigma onto it (unwrap `.cn-sigma/` → root, bind r0 boxes to these refs including `sigma/home`, collapse the current 11-folder sprawl, drain `state/activations.md` into `state/cursors.yaml`) is #690 Sub 2/3. Aligning `cmp` onto the same shape is Sub 5, dependent on cn-sigma proving the model first.

## Entry schema (r0 and rN share it)

```yaml
---
ts:    2026-08-01T11:39:23Z
from:  <activation-id>
rank:  r0                  # r0 | r1 | r2 | …
class: note                # note | decision | request | ack | handoff | review | status | rca
to:    <activation-id>     # optional; omit for a broadcast note
reads:                     # required for rank ≥ r1; omit for r0
  - {repo: usurobor/cnos, ref: refs/heads/sigma/<id>, sha: <sha>}
---
```

Rank↔cadence bound for v0 (daily=r1, weekly=r2, monthly=r3); the "rank = what-it-reads" flexibility is deferred until a missed cadence forces a higher tier to read raw. Cadence is operational scheduling, not structural ontology — see "Rank is not frequency" in `docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md`.

## Constraints (KISS / YAGNI)

Do not add, until the ranked model above proves insufficient:

- `threads/memory/INDEX.md`
- Vector stores
- Graph stores
- Retrieval / search indexes
- Per-rank wake machinery
- Forget/eviction
- Telemetry-as-memory
- A dedicated memory package

## Scope note (this document vs. implementation cells)

This document is the **canonical doctrine** ratified by #690. It is doctrine-only: it does not itself migrate `cn-sigma`, does not implement the box topology, and does not build the compactor. Those are separate, later cells in #690's wave (Sub 2 — dry-run migration map; Sub 3 — cn-sigma migration; Sub 4 — home rollup/no-op read mechanics; Sub 5 — cmp alignment). Do not read the ref paths or invariants above as already running in production.

## Related

- [`docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md`](../../papers/AGENT-MEMORY-LOG-STRUCTURED.md) — the essay this doctrine draws its rank law and coherence framing from; cross-referenced throughout.
- [`docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md`](../conventions/AGENT-ACTIVATION-LOG-v0.md) — historical convention; superseded for agent-memory purposes by this document; still accurate for the writer-locality/wake-ownership mechanics it documents.
- Issue #690 — the ratifying design doc (ranked model, box topology, operator clarification, 2026-08-02).
- Issue #684 / PR #688 — prior substrate exploration (direction-based ref scheme); subsumed, not controlling; closed.
- AGENT-NETWORK.md — agents carry memory when deployed to new workspaces
- HUB-PLACEMENT-MODELS.md — hub is memory, workspace is workbench
- #156 — Attached hubs (agent memory stays in hub, tagged by workspace)
</content>
