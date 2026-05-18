# Cross-Repo Lineage: cph → cnos (bootstrap-cdr)

## Source

- **Repo:** `usurobor/cph` (renamed 2026-05-18 from `usurobor/gait-support-paths`; bundle authored under the prior name)
- **Branch:** `claude/acquire-test-set-model-60jZ2` (intentionally on-branch-only per source-side LINEAGE — never merged to cph main)
- **Authoring commit:** `774427b` (bundle initial filing, under prior repo name)
- **Latest mirror-state commit:** `333dbeb` (cph-side LINEAGE updated to reflect rename + name the cnos-side stale-path consequence)
- **Path:** `.cdd/iterations/cross-repo/cnos/bootstrap-cdr/`

## Target

- **Repo:** `usurobor/cnos`
- **Issue:** [#376](https://github.com/usurobor/cnos/issues/376) — Bootstrap cnos.cdr — research protocol package
- **Mode:** master/tracking issue (`docs-only`); sub-issues to be filed and dispatched separately by δ + γ
- **Disposition:** accepted (body filed verbatim from source `ISSUE.md` @ `774427b`; only the `## Source Proposal` block was inserted at the canonical position per `cdd/issue/SKILL.md` minimal output pattern)
- **First filing branch:** `claude/file-cnos-cdr-issue-fi9Ld` (cnos commit `892a429`; superseded by canonical-path landing on main per this lineage)
- **Canonical path on main:** `cnos:.cdd/iterations/cross-repo/cph/bootstrap-cdr/` (this file)

## Repository rename event

On 2026-05-18, the source repo was renamed `usurobor/gait-support-paths` → `usurobor/cph` (pre-`cph#11` chore aligning the repo name with its identity as a reference cnos.cdr / Coherence Path Hypothesis project, matching the cnos / tsc compact-acronym style). GitHub auto-redirect resolves historical URLs under both names.

Consequences for this cnos-side bundle:

- **First-filed cnos-side path was stale** at `cnos:.cdd/iterations/cross-repo/gait-support-paths/bootstrap-cdr/` on branch `claude/file-cnos-cdr-issue-fi9Ld` (commit `892a429`). That branch is now superseded: the canonical path on cnos main is `cph/bootstrap-cdr/` (this directory). The branch may be deleted as a follow-up; the filing commit is preserved in git for audit.
- **`cnos#376` body is unchanged** — body was filed verbatim from source `ISSUE.md` @ `774427b` under the prior name and remains historical. GitHub redirects cover the embedded URLs.
- **Source-side LINEAGE @ `333dbeb`** already names this stale-path consequence and predicts the canonical `cph/` path; this commit fulfills that prediction. A follow-up cph-side update to mark the prediction realised is out of scope for this cnos-side cycle (cph wave is closed at `42466ad`; any cph-side update is a future cph cycle's scope).

## Bilateral trace

The source-side lineage at `usurobor/cph:.cdd/iterations/cross-repo/cnos/bootstrap-cdr/LINEAGE.md` (commit `333dbeb` on branch `claude/acquire-test-set-model-60jZ2`) names the target as `cnos#376` and the cnos-side mirror as `cnos:.cdd/iterations/cross-repo/cph/bootstrap-cdr/`. This file is the cnos-side counterpart confirming both the filing and the canonical-path landing — the bilateral trace is bilateral.

cnos γ emitted a `FEEDBACK.patch` (per `cdd/gamma/SKILL.md §"Cross-repo proposal close-out"` option (b)) at filing time because the cnos session was scoped to `usurobor/cnos` only. The patch appended the `accepted gamma@cnos cnos#376` event to the source-side `STATUS` ledger; the cph-side wave-closeout (`/.cdd/waves/cdr-refactor-2026-05-18/wave-closeout.md`) confirms the patch was applied. The `FEEDBACK.patch` file is preserved here as a historical artifact of the cross-repo handshake.

Per `cdd/post-release/SKILL.md` Step 5.6b: once this cnos-side mirror exists on cnos main + the cnos.cdr v0.1 wave closes, the source-side bundle may be archived.

## Per-sub confirmation

The master is a scope envelope; sub-issues will land independently. Status filled in as subs close:

| Sub | Proposed shape | cnos issue # | cnos commit | Status |
|---|---|---|---|---|
| Sub 1 | Draft + ratify cdr's six-field instantiation contract (`design-and-build`) | TBD | TBD | not yet filed |
| Sub 2 | Bootstrap `cnos.cdr` package skeleton + `CDR.md` + `SKILL.md` (`MCA`, cites Sub 1) | TBD | TBD | not yet filed |
| Sub 3 | Role overlays in `skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md` (`design-and-build`) | TBD | TBD | not yet filed |
| Sub 4 | Empirical-anchor doc mapping the cph zeroth-pilot wave + segmentation-fix cycle to cdr v0.1 surface (`docs-only`); proposed filename `cnos.cdr/docs/empirical-anchor-cph.md` (pre-rename proposal was `…-gait-support-paths.md` — δ chooses final filename when filing Sub 4) | TBD | TBD | not yet filed |

Proposed sub shapes are non-binding; cnos δ has authoring freedom to re-shape, split, or merge.

## Related cross-repo wave context

The cph-side `cdr-refactor-2026-05-18` wave closed at cph commit `42466ad` (cph main). That wave's close-out (`cph:.cdd/waves/cdr-refactor-2026-05-18/wave-closeout.md`) lists, under §"Out of scope", item 1: **"Cross-repo bundle ingestion on cnos side — renaming `gait-support-paths/bootstrap-cdr/` → `cph/bootstrap-cdr/` and landing on `cnos:main`."** This commit closes that out-of-scope item.

A separate cross-repo proposal bundle for the wave's `F1` next-MCA (filed as `cph#11` → expected `cnos#378`) is **not** part of this bundle and is out of scope for this commit.

Filed by γ on 2026-05-18 (initial branch filing) and ingested by γ-as-δ on 2026-05-18 (canonical-path landing on cnos main).
