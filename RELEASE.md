# 3.81.0

## Outcome

Coherence delta: C_Σ **B−** (α B+, β A, γ C) · **Level:** `L6`

The activation soul is now lean: KERNEL + cap + clp (3 surfaces, down from 7). Four CA skills that were activation-loaded — mca, mci, coherent, agent-ops — are reclassified as on-demand reference. Operational rules from all four are absorbed into cap §4–§6 (MCA, MCI, Coherent output) with explicit non-absorbed boundary notes. The deprecated `cnos.core/AGENTS.md` is deleted. The Go renderer and the R5-activate kata fixture reflect the 2-skill path.

## Why it matters

Each activation loaded 6 CA skills. Four of the six (mca, mci, coherent, agent-ops) duplicated rules already present in cap or described a deprecated runtime (the OCaml-era `cn agent` daemon). Loading all six doubled activation context and propagated a stale surface. After this release, a cnos body session loads exactly: KERNEL → cap → clp → (hub-side PERSONA + OPERATOR). Six fewer SKILL.md files parsed per activation; one fewer root file in the package build.

## What closed

- **#385** — `agent/activate/SKILL.md §2.1 step 2` now lists 2 CA skills (cap, clp) instead of 6. `calls:` frontmatter drops 4 entries. `cap/SKILL.md` gains three new sections (§4 MCA, §5 MCI, §6 Coherent output) absorbing the core rules from the standalone files. mca/mci/coherent/agent-ops each receive `activation_status: on-demand — not activation-loaded since cycle/385` in their frontmatter. `cnos.core/AGENTS.md` is deleted; pkgbuild updated. `activate.go:506` returns the 2-skill path. R5-activate kata P10 fixture updated. All 9 ACs met; 0 β findings at R1; 1 review round.

## γ note (CI)

Build workflow is red on this commit. Pre-existing failure: `Repo link validation (I4)` has been failing since cycle #369 (stale `file://` links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md`); `notify` is downstream of I4. This cycle does not introduce or scope to fix the failure. Local Go tests pass; frontmatter validator passes (67 SKILL.md, no findings); R5-activate kata passes 27/27. Remediation: fix I4 in a near-term cycle.
