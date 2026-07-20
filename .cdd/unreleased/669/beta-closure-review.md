---
cycle: 669
type: beta-closure-review
date: "2026-07-20"
author: beta@cdd.cnos
phase: post-merge-remediation
---

# β Closure Review — #669

**Verdict: REQUEST CHANGES.**

Review base: `4394f14834c759e8378622daa8bdcd08b5224f44`. Reviewed α head: `c28c89ce48a7de407ae466b353c1b7597b0c25a5` (tree `9e34d63a8e034587abbf391748ede5af8860cd5b`). Original matter remains `db9ce95b1beb025d24f43fba2df58692665f80be`; final accepted review head remains `4ff33d1f8aff835aa24a8b11237724cf191d2a6d`; merge commit remains `46ee622aa9bd942ec4d483332c74438f037c776c`.

The α commit has canonical `alpha <alpha@cdd.cnos>` author and committer identity, is directly parented by the stated base, changes only the closure artifact, validator/tests, and lifecycle doctrine/template surfaces, and has zero diff in the frozen recursive-cell CM, #662 inputs, or accepted R6 `self-coherence.md` / `beta-review.md` history. `alpha-closeout.md` is factually post-merge, remains in α voice, preserves standing `none`, and makes no γ disposition.

## Findings

1. **D — mechanical / contract: the new classifier rejects a canonical small-change cycle.** `validate-release-gate.sh` classifies any directory containing `gamma-scaffold.md` as triadic, while `gamma/SKILL.md §2.5` and CDS bootstrap doctrine require γ's scaffold to name `substantial / small-change / immediate-output`. A fixture containing a `mode: small-change` scaffold plus `self-coherence.md` fails pre-merge for missing `beta-review.md` and release for four triadic artifacts. The existing small-change test omits the canonical scaffold, so it does not cover this contract. Required regression pair: a scaffolded explicit small-change must preserve collapsed requirements; a scaffolded substantial cycle without `beta-review.md` must fail.

2. **D — contract: canonical dispatch peers still require future closeouts before review/merge.** `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` §Closeout integrity preflight and §Responsibilities require all six artifacts before requesting `status:review`; `delta/SKILL.md` §9.5 likewise requires all closeouts at `converge`; `activation/SKILL.md §12` still says all five files must be non-empty before merge. These directly contradict the new pre-merge gate and the same activation section's new post-merge note. Update the canonical/golden projections so review/merge readiness requires only structurally available triadic artifacts and closeouts remain post-merge.

3. **D — contract / mechanical: post-merge closure ordering is not executable across the changed peers.** The new γ paragraph says to append the terminal marker and run `--mode post-merge --cycle N`, with `RELEASE.md` later; the surrounding γ closure rows still require `RELEASE.md`, release-boundary preflight, and the cycle-directory move before closure. Once moved, the exact-cycle oracle fails because it reads only `.cdd/unreleased/N`. CDS S9–S11 and its ownership rule retain the same ordering, and the ownership rule still says filename existence satisfies δ despite the new marker requirement. `cdd-status` also reports filename existence as the closure signal, and the release-effector kata verifies only existence. Required regression pair: the post-merge oracle must pass at the one doctrinally selected pre-move point; a pre-operator assurance file without the marker must remain nonterminal on every operator-visible surface.

4. **D — mechanical: the changed activation minimum workflow is invalid on its declared pull-request trigger.** `activation/SKILL.md §9` retains `pull_request: branches: [main]` but always derives `CYCLE_NUMBER` from `GITHUB_REF_NAME`; on a PR ref such as `670/merge`, the validator correctly rejects the nonnumeric selector. The vendored workflow avoids that failure only by having no pull-request trigger and skips main entirely, leaving the two projections behaviorally divergent. Required regression pair: cycle-branch push derives the exact numeric suffix and validates; pull-request/main events select an explicit safe phase/cycle source or are excluded consistently without claiming validation ran.

## Verification

- `git diff --check` and `bash -n` on changed shell: pass.
- `bash scripts/test-validate-release-gate.sh`: 29/29 assertions pass.
- `bash scripts/test-release-tag-integration.sh`: pass; annotated tag verified.
- Activation workflow YAML parses successfully.
- Exact #669 pre-merge gate: pass. Exact post-merge and release gates: fail only for missing `beta-closeout.md` and missing exact terminal γ marker, as expected before β/γ closure writes.
- Additional selector confinement, missing-value, invalid-mode, and path-with-spaces checks: pass.

The accepted recursive CM and frozen inputs were neither changed nor rerun. Standing remains `none`. This artifact is findings-only: no β closeout, γ marker, release assignment, disconnect, issue/FSM action, or standing change is authorized.
