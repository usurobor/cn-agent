# Cross-Repo Lineage: cn-sigma → cnos (agent-activate-skill)

## Source

- **Repo:** `usurobor/cn-sigma`
- **Branch:** `claude/review-context-EbgR2` (drafted in a Claude Code on the web session; merged to cn-sigma main 2026-05-19)
- **Authoring commit:** `8f153c15e` (initial bundle filing on the source branch)
- **AC4-refinement commit:** `bdda457f5` (3-tier capability-matrix replaces single capability-gate AC4; folded into bundle 7 min after issue filing per `threads/adhoc/20260519-git-read-and-untested-limits.md`)
- **Post-merge cn-sigma main HEAD:** `1a4e25f75` (carries full bundle + AC4 refinement)
- **Path:** `.cdd/iterations/cross-repo/cnos/agent-activate-skill/` (full bundle on cn-sigma main; not minimal-anchor pattern — see "Source posture" below)
- **Source posture:** cn-sigma is an agent hub, not a fully CDD-activated tenant repo. This bundle adopts the cross-repo trace convention (cnos `CDD.md §5.3`, `gamma/SKILL.md "Cross-repo proposal close-out"`) and the pattern established by `cph/bootstrap-cdr` (cnos commit `7a7f7152`, 2026-05-18). cnos #377 (open) tracks codification of the hub→CDD-repo case; this bundle is a concrete exemplar. Unlike `cph/bootstrap-cdr` (minimal anchor on source main; full bundle on dormant branch), this bundle is carried in full on cn-sigma main per operator preference.

## Target

- **Repo:** `usurobor/cnos`
- **Issue:** [#379](https://github.com/usurobor/cnos/issues/379) — agent/activate skill — single source of truth for agent self-activation
- **Mode:** `design-and-build` (per `ISSUE.md` mode declaration); single-focus, ~7 ACs, two tightly coupled modules (new skill + Go renderer)
- **Disposition:** accepted (body filed verbatim from source `ISSUE.md` @ `8f153c15e`; `## Source Proposal` block inserted at canonical position per `cdd/issue/SKILL.md`; AC4 capability-matrix refinement folded into issue body post-filing per Delta below)
- **First filing commit (issue body):** filed at `2026-05-19T02:38:33Z` from bundle SHA `8f153c15e` (pre-AC4-refinement)
- **Canonical-path landing on main:** `cnos:.cdd/iterations/cross-repo/cn-sigma/agent-activate-skill/` (this file)

## Delta

The cnos#379 issue body was edited at `2026-05-19T10:21:01Z` to fold in the AC4 capability-matrix refinement from cn-sigma branch HEAD `bdda457f5`. Specifically:

- **AC4 title:** "Skill is body-agnostic and names the capability gate" → "Skill is body-agnostic and names the capability matrix"
- **AC4 invariant:** single fetch requirement + degraded path → 3-tier body-capability ladder enumerating (a) shell + git, (b) HTTP fetch only, (c) no fetch
- **AC4 oracle/positive/negative:** updated to match the 3-tier framing
- **`## Source Proposal` block:** Source commit advanced from `8f153c1` (filing-time branch HEAD) to `1a4e25f75` (post-merge cn-sigma main HEAD); filing-time SHA retained as a separate field; Delta field expanded with the refinement record

The refinement is anchored in cn-sigma `threads/adhoc/20260519-git-read-and-untested-limits.md` (commit `bdda457f5`, also merged to cn-sigma main at `1a4e25f75`), which documents that bodies in Claude Code on the web/phone have `git` available and can `git clone` the cnos repo — invalidating the pre-refinement assumption that WebFetch was the only file-read mechanism for non-cn bodies.

## Bilateral trace

The source-side lineage at `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/LINEAGE.md` (commit `1a4e25f75` on cn-sigma main) names the target as `cnos#379` and the cnos-side mirror as `cnos:.cdd/iterations/cross-repo/cn-sigma/agent-activate-skill/`. This file is the cnos-side counterpart confirming acceptance, the canonical-path landing, and the post-filing AC4 reconciliation — the bilateral trace is bilateral.

Per `cdd/post-release/SKILL.md` Step 5.6b: once this cnos-side mirror exists on cnos main and the agent-activate-skill cycle closes, the source-side bundle may be archived.

## Related cross-repo wave context

This is the second cross-repo proposal landing on cnos main with the asymmetric mirror pattern (target carries the canonical archive). The first was `cph/bootstrap-cdr` @ `7a7f7152` (2026-05-18). The cn-sigma instance differs in source posture: source main carries the full bundle (ISSUE.md + LINEAGE.md + STATUS) rather than a minimal anchor — chosen because cn-sigma is not a CDD-activated tenant repo and the bundle is also a public exemplar for the hub→CDD-repo case tracked by cnos#377.

Filed by usurobor on 2026-05-19 (initial filing) and reconciled by γ-as-δ on 2026-05-19 (canonical-path landing on cnos main + AC4 fold-in on issue body).
