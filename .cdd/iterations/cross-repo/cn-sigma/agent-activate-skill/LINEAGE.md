# Cross-Repo Lineage: cn-sigma → cnos (agent-activate-skill)

## Source

- **Repo:** `usurobor/cn-sigma`
- **Branch:** `claude/review-context-EbgR2`
- **Commit:** `8f153c15e634032c13283fca7b5c5a64eaa5cf4b` (short: `8f153c1`)
- **Path:** `.cdd/iterations/cross-repo/cnos/agent-activate-skill/`
- **Source posture:** cn-sigma functions as an agent hub rather than a fully CDD-activated tenant repo. The bundle adopts the cross-repo trace convention without otherwise activating CDD in cn-sigma. cnos#377 tracks codification of the hub→repo case.

## Target

- **Repo:** `usurobor/cnos`
- **Issue:** [#379](https://github.com/usurobor/cnos/issues/379) — *agent/activate skill — single source of truth for agent self-activation*
- **Filed on branch:** `claude/file-cnos-cdr-issue-fi9Ld`
- **Filing commit:** TBD (this commit)
- **Disposition:** accepted (body filed verbatim from source `ISSUE.md`; only the placeholder fields in `## Source Proposal` were filled in with concrete values)
- **Mode:** design-and-build (per source `ISSUE.md` mode declaration; 7 ACs, two modules tightly coupled)
- **Labels applied:** `design`, `P2`, `core` (matches source proposed labels; `core` was auto-created by MCP — did not previously exist in cnos)

## Bilateral trace

This file is the cnos-side mirror following the precedent at `cnos:.cdd/iterations/cross-repo/tsc/cdd-supercycle/LINEAGE.md` and `cnos:.cdd/iterations/cross-repo/gait-support-paths/bootstrap-cdr/LINEAGE.md`.

The source-side `LINEAGE.md` at `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/LINEAGE.md` (commit `8f153c1`) predicts the cnos-side mirror at exactly this path; that prediction is now realised.

cnos γ cannot write to `usurobor/cn-sigma` from this session (MCP scoped to `usurobor/cnos` only). Per `cdd/gamma/SKILL.md §"Cross-repo proposal close-out"` option (b), a feedback patch is emitted in this same bundle at `FEEDBACK.patch` for the cn-sigma-side γ to apply on receipt; the patch lands the `accepted gamma@cnos cnos#379` event on the source `STATUS` ledger.

## Protocol observations recorded for cnos#377

The intake of this proposal surfaced three additional cross-repo protocol gaps (the same `cdd-protocol-gap` finding tracked by cnos#377):

1. **STATUS event vocabulary drift.** Source ledger uses `drafted` as the initial event, not `submitted` (the documented vocabulary is `submitted | accepted | modified | rejected | landed`). `drafted` was treated as filing-ready for intake purposes; cnos#377 should reconcile or expand the vocabulary.
2. **Bundle file set is partial.** The cn-sigma bundle ships `ISSUE.md`, `STATUS`, `LINEAGE.md` — no `README.md`. The bootstrap-cdr bundle shipped all four. cnos#377 should codify whether `README.md` is required or optional.
3. **Source `## Source Proposal` block carries placeholders.** Source `ISSUE.md` shipped with `Source commit: filled at filing` and `Disposition: pending` as deliberate template fields for cnos γ to fill in. This is a cleaner author-side convention than the bootstrap-cdr pattern (insert a brand-new block); cnos#377 should consider adopting it.

## Per-deliverable confirmation

The source body deliverable maps 1:1 to cnos#379. This is *not* a master/sub-shaped proposal — it is a single-cycle, 7-AC, design-and-build proposal.

| Source artifact | cnos artifact | Status |
|---|---|---|
| `ISSUE.md` (body) | cnos#379 body | filed verbatim |
| Proposed labels (`design`, `P2`, `core`) | cnos#379 labels | applied; `core` auto-created by MCP |
| `STATUS` event ledger | source `STATUS` + `FEEDBACK.patch` (this bundle) | feedback patch emitted for source-side close-out |
| Source `LINEAGE.md` | this file | mirror created |

Filed by γ on 2026-05-19.
