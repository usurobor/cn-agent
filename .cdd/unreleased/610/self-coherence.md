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

## §Skills

- **Tier 1:** `CDD.md` (canonical lifecycle/role contract) + `cnos.cdd/skills/cdd/alpha/SKILL.md`
  (this file's own role surface, followed for load order, peer enumeration, self-coherence
  discipline, and the pre-review gate).
- **Tier 2:** none newly loaded beyond the always-applicable `eng/*` the α load order names
  generically (Go is the only language touched besides one markdown prose edit; no new
  Tier 2 bundle was required beyond what the existing `repoinstall`/`cli` packages already
  assume).
- **Tier 3 (issue-specific, named by γ's scaffold):**
  - `cnos.core/skills/write/SKILL.md` — applied to the two prose edits in
    `cds-dispatch/SKILL.md` (front-load the point, one fix per leak, no filler) and to the
    help-text rewrite in `cmd_repo_install.go` (state the workflow-scope PAT fact once,
    point to it rather than repeating).
  - `docs/development/design/cn-repo-install-MOCKS.md` (Mock C + Mock E) — read in full;
    Mock C's transcript/invariants (C1–C6) and Mock E's tenant-portability invariants
    (E1–E4) are the design source of truth this cell's oracle table maps to.
  - `.cdd/unreleased/610/gamma-scaffold.md` — read in full before coding; its peer
    enumeration (confirming the #609 renderer's tenant-portable finalizer/scanner paths
    already exist; confirming cnos#493's mechanism is genuinely absent via `find`/`rg`) was
    reused rather than re-derived, per the scaffold's own instruction.
