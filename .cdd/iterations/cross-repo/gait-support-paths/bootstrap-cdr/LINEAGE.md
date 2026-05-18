# Cross-Repo Lineage: gait-support-paths → cnos (bootstrap-cdr)

## Source

- **Repo:** `usurobor/gait-support-paths`
- **Branch:** `claude/acquire-test-set-model-60jZ2`
- **Commit:** `774427b`
- **Path:** `.cdd/iterations/cross-repo/cnos/bootstrap-cdr/`

## Target

- **Repo:** `usurobor/cnos`
- **Issue:** [#376](https://github.com/usurobor/cnos/issues/376) — Bootstrap cnos.cdr — research protocol package
- **Filed on branch:** `claude/file-cnos-cdr-issue-fi9Ld`
- **Disposition:** accepted (body filed verbatim from source `ISSUE.md`; only the `## Source Proposal` block was inserted at the canonical position per `cdd/issue/SKILL.md`)
- **Mode:** master/tracking issue (`docs-only`); sub-issues to be filed and dispatched separately by δ + γ

## Bilateral trace

The source-side bundle at `usurobor/gait-support-paths:.cdd/iterations/cross-repo/cnos/bootstrap-cdr/` (commit `774427b`) carries the master issue body that this lineage file confirms landed as cnos#376. This file is the cnos-side mirror following the precedent at `cnos:.cdd/iterations/cross-repo/tsc/cdd-supercycle/LINEAGE.md`.

cnos γ could not write to `usurobor/gait-support-paths` from this session (cross-repo write not authorised in this environment). Per `cdd/gamma/SKILL.md §"Cross-repo proposal close-out"` option (b), a feedback patch was emitted in this same bundle at `FEEDBACK.patch`. The gait-side γ applies the patch on receipt; that lands the `accepted gamma@cnos cnos#376` event on the source `STATUS` ledger.

Per `cdd/post-release/SKILL.md` Step 5.6b: once the source-side bundle archives in response to the cnos-side mirror existing, the bilateral trace is bilateral and the source-side copy may be deleted.

## Per-sub confirmation

The master is a scope envelope; sub-issues will land independently. Status filled in as subs close:

| Sub | Proposed shape | cnos issue # | Status |
|---|---|---|---|
| Sub 1 | Draft + ratify cdr's six-field instantiation contract (`design-and-build`) | TBD | not yet filed |
| Sub 2 | Bootstrap `cnos.cdr` package skeleton + `CDR.md` + `SKILL.md` (`MCA`, cites Sub 1) | TBD | not yet filed |
| Sub 3 | Role overlays in `skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md` (`design-and-build`) | TBD | not yet filed |
| Sub 4 | Empirical-anchor doc mapping gait waves to cdr v0.1 surface (`docs-only`) | TBD | not yet filed |

Filed by γ on 2026-05-18.
