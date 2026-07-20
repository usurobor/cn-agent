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

## R1 — review of α repair `ccfa1723`

**Verdict: REQUEST CHANGES.**

Review base: `8cf73eaec53c8b42f22e4ae9286125fc8005a210`. Reviewed α head: `ccfa1723ebb2cdda1b1082f20d1701cca1c56dfa` (tree `65658964fc26efdf8a192cbbb042ae86e2245ba4`). The repair is directly parented by the R0 β findings commit and has canonical `alpha <alpha@cdd.cnos>` author and committer identity. Original matter remains `db9ce95b1beb025d24f43fba2df58692665f80be`; final accepted review head remains `4ff33d1f8aff835aa24a8b11237724cf191d2a6d`; merge commit remains `46ee622aa9bd942ec4d483332c74438f037c776c`; standing remains `none`.

The repair correctly makes the activation artifact trigger cycle-push-only, makes golden and live dispatch byte-identical, gives `cdd-status` and both post-merge/release gates the exact `CDD-Post-Merge-Closeout: complete` marker, and removes the pre-tag cycle-directory move from `scripts/release.sh`. Three binding contradictions remain.

### Findings

1. **D — mechanical / honest claim: an unknown cycle shape still fails open as `small-change`.** `scripts/validate-release-gate.sh` says a missing or unrecognized declaration fails closed, but its condition enters the triadic path only when `beta-review.md` exists or a `gamma-scaffold.md` exists with a non-collapsed parsed mode. A selected directory with `self-coherence.md` but no scaffold, and even an entirely empty selected directory, therefore reaches the unconditional `small-change` success branch. Fresh adversarial fixtures for cycles 706 and 707 both returned exit 0 from `--mode pre-merge --cycle N` with `cycle N (small-change)`; explicit `small-change` and `immediate-output` returned 0, while explicit `substantial`, an unrecognized mode, and a scaffold missing its mode line correctly returned 1. This contradicts the R1 claim that a missing declaration fails closed and permits a malformed exact target to pass review-readiness. Require an explicit recognized collapsed-mode declaration for the collapsed path; absence of the scaffold/declaration must fail closed, with regressions for self-coherence-only and empty selected directories.

2. **D — contract: review-readiness peers still retain closeout-bound convergence requirements.** `delta/SKILL.md §9.5` says every R[N] boundary MUST carry the table's artifacts, and its `converge` row still enumerates `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, and optional PRA. The adjacent sentence and §9.6 now say the same converge/status transition needs only the three review-time artifacts, so the canonical δ contract contradicts itself. In `cds-dispatch`, golden, and live workflow, the pre-`status:review` preflight still directs the actor to record `deliverable_evidence` “on the cycle's closeout,” despite the repaired list and responsibilities saying closeouts are post-merge. This also makes the R1 alpha-closeout claim that `cds-dispatch` and δ bind readiness only to pre-merge artifacts factually too broad. Replace the δ converge row with the review-time set and give `deliverable_evidence` an explicitly pre-merge owner/path in the source, golden, and live projections; no closeout artifact or ambiguous “cycle's closeout” may be a review-readiness prerequisite.

3. **D — contract / ordering: canonical lifecycle authorities still put terminal closure before the disconnect or the marker-producing closeout after release.** The repaired CDS state machine has the executable order `S9 marked post-merge closeout → S10 release prep → S11 δ disconnect/green CI → S12 γ terminal archive`, but the same file declares its 0–13 step table the canonical ordering and still orders `Gate`/`Release` at steps 9–10 before PRA/`gamma-closeout.md`/closure at steps 12–13. The duplicated ordered-artifact flow and manifest retain that sequence. Other unsuperseded authoritative prose still calls `--mode post-merge` an “explicit terminal γ declaration,” says the gate moves the cycle to “terminally closed” before the later release gate, and says γ declares final closure in `gamma-closeout.md`. `gamma/SKILL.md`'s frontmatter and step map still place a “closure declaration” before step 17's δ disconnect; the operator kata and release-effector governing question still say γ declares closure and δ then cuts the release. The alpha artifact explicitly supersedes its own R0 terminal-marker proposal, so its retained R0 prose is historical; these live canonical peers are not similarly qualified. Reconcile the canonical step table, ordered flow/manifest, gamma step map, operator kata, release-effector framing, and terminal wording to the one executable order: marked nonterminal post-merge closeout while under `unreleased/`, release validation, δ tag/disconnect plus green CI, then γ archive and only then terminal closure.

### Verification

- `git diff --check` and `bash -n` on every changed shell file: pass.
- YAML parse for the activation template, dispatch golden, and live workflow: pass.
- `bash scripts/test-validate-release-gate.sh`: 36/36 assertions pass; the additional no-scaffold and empty-directory fail-closed cases expose finding 1.
- `bash src/packages/cnos.cdd/commands/cdd-status/test-cn-cdd-status.sh`: 27/27 assertions pass.
- `bash scripts/ci/check-dispatch-closeout-integrity.sh`: pass.
- `bash scripts/test-release-tag-integration.sh`: pass; annotated tag behavior preserved.
- Dispatch golden re-render: unchanged; golden and live workflow are byte-identical at SHA-256 `b2580e1d254d5031a0a2ce9a9ea865b3d2a91a98d015aa56bc89236faf43d014`.
- Exact #669 pre-merge gate: pass. Exact post-merge and release gates fail only for missing `beta-closeout.md` and the absent exact post-merge marker, as expected before β/γ closeout writes.
- Frozen diff against the R0 β base for the recursive-cell CM, #662 inputs, and accepted R6 `self-coherence.md` / `beta-review.md`: empty. The accepted CM was not rerun.
- Exhaustive authoritative-source grep finds no surviving literal `CDD-Cycle-Closure: terminal`; it does find the unqualified terminal/ordering claims in finding 3. No check run or commit status was published for exact alpha SHA `ccfa1723ebb2cdda1b1082f20d1701cca1c56dfa`; local named gates above are the available evidence.

This R1 artifact remains findings-only. No `beta-closeout.md`, γ marker, release assignment, tag/disconnect, branch deletion, issue/FSM action, recursive-CM execution, or standing change is authorized.

## R2 — review of α repair `472b53d2`

**Verdict: REQUEST CHANGES.**

Review base: `bec4e4375315d7e54f86c1ff4ffdbcb71f19251a`. Reviewed α head: `472b53d2159acd62c467a3429a229a4716cdb90c` (tree `09d6c587fb0ff6a64e274c81be0475312c00f910`). The repair is directly parented by the R1 β findings commit and has canonical `alpha <alpha@cdd.cnos>` author and committer identity. Original matter remains `db9ce95b1beb025d24f43fba2df58692665f80be`; final accepted review head remains `4ff33d1f8aff835aa24a8b11237724cf191d2a6d`; merge commit remains `46ee622aa9bd942ec4d483332c74438f037c776c`; standing remains `none`.

R2 repairs the exact-selector cases, the δ converge table, the `deliverable_evidence` authority/path, and the core CDS S0–S13 state machine. The dispatch renderer is idempotent and source/golden/live are byte-identical. Four binding contradictions remain.

### Findings

1. **D — mechanical / fail-closed contract: the default release path treats every unscaffolded directory as legacy small-change.** Exact `--cycle` validation now correctly rejects missing, empty, self-only, missing-mode, and unknown-mode shapes, and correctly distinguishes explicit small/immediate/substantial modes. Repository-wide mode has no legacy cutoff, manifest, exception, age, or provenance test: the `-z "$CYCLE"` branch unconditionally classifies any unscaffolded directory without `beta-review.md` as `small-change`. `scripts/release.sh` invokes that repository-wide mode. A fresh adversarial release fixture containing a newly-created `.cdd/unreleased/999/self-coherence.md`, no scaffold, and `RELEASE.md` returned exit 0 with `cycle 999 (small-change)` and `Release gate passed`. The 51/51 regression suite encodes the same unconditional bypass as “legacy compatibility,” so its green result does not make the compatibility explicit or non-dangerous. Bound legacy compatibility to an observable legacy population/cutoff or exception and fail closed for unclassified new directories; add the new-malformed repository-wide release case as a refusal regression.

2. **D — contract: source/golden/live dispatch still make a future closeout a `status:review` prerequisite on repair passes.** `cds-dispatch/SKILL.md` Step E says every closeout landed by a repair re-entry MUST carry `repair_evidence` and the cell MUST NOT advance to `status:review` until it is complete. The same wording is rendered in the golden and live workflow. Later responsibility text says closeouts MUST NOT gate review and shifts the requirement to a review-time receipt, so each of the three authoritative surfaces contradicts itself. The repaired top-level `REVIEW-REQUEST.yml` `deliverable_evidence` block and δ §9.5 converge row are correct, but they do not supersede Step E. Give repair evidence one pre-merge owner/path (for example a field in `REVIEW-REQUEST.yml` or `REPAIR-PLAN.md`) and remove every closeout prerequisite from the transition; extend the integrity guard to reject `closeout ... status:review` coupling instead of merely proving the prose is present.

3. **D — contract / observability: the lifecycle peers still do not tell or observe one executable closeout → disconnect → archive → terminal order.** The central CDS table/FSM and the new Phase A/B/C gamma block correctly order marked nonterminal closeout under `unreleased/` → release validation → δ tag/green CI → γ PRA/archive → later terminal declaration. Unsuperseded peers conflict with that order: `beta/SKILL.md` says the same β session owns “the release and β close-out” and gives the positive sequence `merge → release → beta-closeout`; `release/SKILL.md` still makes `beta-closeout.md` carry release evidence after the tag and labels it step 10, while also retaining the internally contradictory “CI must pass before tag” / “If CI fails after tag” rule; gamma headings number release preparation as Steps 6–7 before closeout Steps 8–9; and `receipt-stream/SKILL.md` contains a second step 6 that directly commands γ to move the directory before requesting the tag, immediately after the repaired post-disconnect archive step. `cn-cdd-status` was not changed in R2: it reads closeouts and the nonterminal marker only from `origin/cycle/{N}:unreleased`, checks an archive only in the working tree, never checks the release tag/green CI or the archived gamma terminal binding, leaves nine conditions permanently `⚠`, and therefore can never reach its own 14/14 “terminally closed” branch. Its 27/27 suite tests only the nonterminal marker, not terminal closure. Reconcile every live role/release/handoff peer and make status mechanically distinguish release-pending, released-but-open, archived, and terminally sealed states.

4. **D — mechanical / warning semantics: ledger R2 hides genuinely stale unreleased closeouts.** `checkUnreleasedTriadicArtifacts` changed every present unreleased `alpha-closeout.md`, `beta-closeout.md`, and `gamma-closeout.md` from the stale-cycle warning to unconditional PASS, without checking the nonterminal marker, assigned release, tag/disconnect, archive opportunity, terminal binding, age, or an exception. The new focused test creates five one-line placeholder artifacts and requires all three closeouts to PASS; it has no released-but-never-archived negative case. The only orphan warning is in the separate `--version --cycle` path on a current branch literally named `main`; `--unreleased` and `--all`—including the source-built command cited by α—now silently accept a stale post-release directory forever. Preserve the valid release-pending state without broadening it into a permanent exception: discriminate pre-disconnect pending receipts from records whose release/archive evidence makes them stale, and add both positive and stale-negative regressions.

### Verification

- Exact target/provenance: remote `closure/669` was `472b53d2159acd62c467a3429a229a4716cdb90c`; parent `bec4e4375315d7e54f86c1ff4ffdbcb71f19251a`; tree `09d6c587fb0ff6a64e274c81be0475312c00f910`; canonical α author/committer; clean fresh clone before this review write.
- `scripts/test-validate-release-gate.sh`: 51/51 pass. The additional new-malformed repository-wide release fixture returns exit 0 and exposes finding 1.
- `src/packages/cnos.cdd/commands/cdd-status/test-cn-cdd-status.sh`: 27/27 pass; source/test semantic inspection exposes the unobservable/unreachable terminal branch in finding 3.
- Dispatch closeout-integrity and repair-preflight guards: pass. Direct renderer replay is idempotent; rendered output, golden, and live workflow are byte-identical at SHA-256 `18c8cfb8331a3fb5f243c80d67776177c1decd5cf54f3d17d03ec306b37f3b10`.
- Release-tag integration, shell syntax, YAML parse, `git diff --check`, skill-frontmatter self-test, and full `100 SKILL.md` validation: pass.
- Exact #669 pre-merge gate passes. Exact post-merge and release gates fail only for the still-missing `beta-closeout.md` and absent nonterminal `CDD-Post-Merge-Closeout: complete` marker; no closeout action follows from this review.
- The Go toolchain is absent in this fresh environment, so `go test` and a source-built ledger replay could not be independently rerun. Static source/test inspection proves finding 4: the production branch emits unconditional PASS and the new test requires that result for content-free placeholders. The installed `cn` is older and has no `cdd` command.
- Frozen diff from R1 across the recursive-cell CM/runner/tests, schema, frozen #662 calibration inputs, accepted R6 `self-coherence.md` / `beta-review.md`, and γ-owned #669 artifact is empty. The accepted CM was not rerun.
- R0 terminal-marker prose in `alpha-closeout.md` is now clearly labeled historical and superseded by R1/R2. No literal legacy `CDD-Cycle-Closure: terminal` remains in authoritative source.
- After fetch, `origin/main` is `07078341882529b68aba255e3c775411aeb9aeba` and is 3 commits ahead / 5 behind this head. Its three new commits touch only `.cn-sigma/logs/20260720.md`; `git merge-tree` reports no conflict, so a later operator-authorized rebase appears safe. No rebase was performed.

This R2 artifact remains findings-only. No `beta-closeout.md`, γ marker, release assignment, tag/disconnect, archive/terminal declaration, branch deletion, issue/FSM action, recursive-CM execution, or standing change is authorized.
