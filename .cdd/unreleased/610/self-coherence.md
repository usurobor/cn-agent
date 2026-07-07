# α self-coherence — cnos#610

## §Gap

**Issue:** [cnos#610](https://github.com/usurobor/cnos/issues/610) — "cds-install Sub 3:
`cn repo install --dispatch cds` (dispatch layer)". Mode (issue header): `design-and-build`.
CDD cycle-sizing mode (per γ-scaffold): **substantial**.

**Gap (from the issue's Problem/Expected/Divergence):** no way to render a repo's dispatch
workflow through the installer; `--dispatch cds` unconditionally refused (Mock B4 guard,
pending #609). #609 (renderer generalization) merged via PR #619, so the guard is now
stale. Separately, the `cds-dispatch` wake's rendered prompt body still carries two
tenant-visible hardcoded-sigma leaks, confusing for a real third-party render (operator
#610 directive, folded into AC5 and the issue's Scope).

**What this cycle closes:** wires `cn repo install --dispatch cds` to the (now-merged)
#609 renderer with explicit caller identity, detects (but does not implement) the
cnos#493 canonical-label-install gap and reports it actionably, states the
`workflow`-scope PAT + never-pushes-main facts, and de-sigmas the two rendered-prompt-body
leaks while preserving `--agent sigma` operational compatibility.

**Depends on:** #608 (base installer, merged), #609 (renderer generalization, merged via
PR #619), cnos#493 (canonical label install, **still OPEN** — confirmed via
`gh issue view 493 --repo usurobor/cnos --json state` at authoring time; this cell does
**not** implement it, per Non-goals).

**Branch:** `cycle/610`. HEAD at review-readiness time: `f334807` (implementation SHA;
rebased onto `origin/main` at `f7e9aaa` — see §Review-readiness for the rebase record).
