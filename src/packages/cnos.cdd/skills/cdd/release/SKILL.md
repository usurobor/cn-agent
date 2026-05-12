---
name: release
description: Ship a version as a measured coherence delta with matching code, changelog, tag, release notes, deployment, and validation.
artifact_class: skill
kata_surface: embedded
governing_question: How does β turn a converged branch into a released version with a complete and auditable coherence delta?
visibility: internal
parent: cdd
triggers:
  - release
  - ship
  - tag
  - version
  - deploy
  - changelog
scope: task-local
inputs:
  - approved branch
  - version decision
  - release context
outputs:
  - release artifact set
  - beta close-out
requires:
  - approval and release readiness
calls: []
---

# Release

## Core Principle

**Coherent release: every version bump is a measured coherence delta with a complete audit trail.**

β owns: review approval outcome, `git merge` into main, and β close-out. **δ owns tag/release/deploy** (the release boundary — `scripts/release.sh`, binary publication, deployment). γ owns `RELEASE.md` authoring, cycle-directory movement, and the post-release assessment. β does not tag, push tags, bump versions for release, or cut the disconnect release. See `CDD.md` §1.4 β algorithm, γ algorithm Phase 5a, and `operator/SKILL.md` §3.4.

Canonical artifact locations (β close-out path, RELEASE.md, snapshot dirs, tag policy) are defined in `CDD.md` §5.3a (Artifact Location Matrix). All tags are bare `X.Y.Z`; `v`-prefixed tags are legacy and warn-only.

A release has parts: readiness check, version decision, changelog, release notes, tag, binaries, deployment, validation. Coherence = each part completed and each artifact matches the others (version in code = tag = changelog = binary = deployed agent). The released system is validated, not merely published.

Failure mode: version drift — tag says X, binary says Y, agent reports Z. Or: releasing without validation, so breakage ships silently. Or: incomplete deploy where artifacts ship but assessment never lands.

## 1. Define

1.1. **Identify the parts**
  - Readiness: CI green, branches merged or returned, main clean
  - Version: semver bump decided from commit content
  - Artifacts: cn.json, VERSION, package manifests, CHANGELOG, RELEASE.md, tag, GitHub release, binaries
  - Deployment: binary on target hosts, `cn deps restore`, validation
  - ❌ Push tag, hope it works
  - ✅ Every artifact checked, deployed, validated end-to-end

1.2. **Articulate how they fit**
  - Version flows: VERSION → `stamp-versions.sh` → cn.json + package manifests → CHANGELOG → commit → tag → CI → binaries → deploy → validate
  - Each step depends on the prior; skip one and the chain breaks
  - ❌ Tag pushed before version bumped in code (binary reports old version)
  - ✅ Version string updated everywhere before commit, tag matches code

1.3. **Name the failure mode**
  - Version drift: any two artifacts disagree on the version
  - Silent breakage: release ships but CI failed or validation skipped
  - Incomplete deploy: binary updated but packages stale (`cn deps restore` not run)
  - ❌ "cn update worked" (but agent still reports old version)
  - ✅ `cn --version` matches tag, `cn deps restore` run, agent can execute ops

## 2. Unfold

2.1. **Readiness check**
  - CI passing on main (unit tests + integration tests + kata suite)
  - No unmerged branches that should be in this release
  - `git branch -r --no-merged origin/main` — review each, merge or defer
  - **Non-destructive merge-test on the merge tree before merge.** Build the merge tree in a throwaway worktree (`git worktree add /tmp/cnos-merge-test/wt origin/main && cd /tmp/cnos-merge-test/wt && git merge --no-ff --no-commit origin/cycle/{N}`); run any contract validators the cycle ships (e.g. `./tools/validate-skill-frontmatter.sh` if the cycle adds or modifies SKILL.md frontmatter; `cn-cdd-verify` for `.cdd/` artifacts; `scripts/check-version-consistency.sh` for version stamping); confirm zero unmerged paths and zero new validator findings. Tear down the worktree. **Set any worktree-local git config with the explicit `--worktree` flag** (`git config --worktree user.name X`) to avoid leaking identity to the shared repo `.git/config` (cycle #301 O8). This pattern is also referenced from `beta/SKILL.md` § pre-merge gate row 3 — it is the same step, named here so the release-flow doc is authoritative.
  - ❌ Release with failing CI ("it's just a flaky test")
  - ❌ Merge without running the cycle's own validators on the merge tree ("the cycle branch passes; the merge tree will too")
  - ✅ CI green, or known failure documented and accepted (e.g. Coherence workflow on #22)
  - ✅ Merge-tree validation green (especially when the cycle ships a new contract surface like #301's I5 frontmatter validator)

2.2. **Version decision**
  - Review commits since last tag: `git log --oneline $(git describe --tags --abbrev=0)..HEAD`
  - Major: breaking changes, paradigm shift
  - Minor: new features, backwards compatible (new ops, new skills, new specs)
  - Patch: bug fixes only
  - ❌ Patch for a new runtime feature (undersells the change)
  - ✅ Minor for new `Cn_shell.execute` entry point; patch for CI fix

2.3. **Version bump — VERSION-first flow**
  - `VERSION` file at repo root is the single source of truth
  - Edit `VERSION` to the new version string (bare, no `v` prefix)
  - Run `scripts/stamp-versions.sh` — derives all manifests (`cn.json`, package `cn.package.json` files) from VERSION
  - Run `scripts/check-version-consistency.sh` — validates all version-stamped files agree
  - Binary embeds version from VERSION at build time — no manual edit needed
  - Tests and katas read version dynamically — no manual edit needed
  - ❌ Edit cn.json or package manifests by hand (bypasses single source of truth)
  - ❌ Skip stamp-versions.sh ("I'll update the manifests manually")
  - ✅ `echo "X.Y.Z" > VERSION && scripts/stamp-versions.sh && scripts/check-version-consistency.sh`

2.4. **CHANGELOG**
  - Add a ledger row matching the format defined in `CHANGELOG.md` § Release Coherence Ledger. That format is canonical — do not reinvent the row shape here.
  - The ledger row includes: Version, C_Σ, α, β, γ, Level, Rounds, and a coherence note. The **Rounds** column records the review-round count for the release's cycle (e.g. `1`, `2`, `3`); for releases bundling multiple cycles, sum or list (`1+2`). This is a learning-curve indicator — a rising trend signals process drift, a falling trend signals accumulated context.
  - **Provisional vs. final scoring rule:** β writes the CHANGELOG TSC row at release commit with `provisional, pending γ PRA` in the level cell (or in a parenthetical next to the level). γ updates the cell to the final value in the same commit as the PRA. This ensures one writer per fact while preserving the release-time ledger entry.
  - Add a detailed section below the ledger: Added / Changed / Fixed, with each commit's impact named and linked to issues.
  - ❌ "Various improvements" (no detail)
  - ❌ Ledger row without engineering level (Level column is required)
  - ✅ Each commit's impact named, linked to issues
  - ✅ Honest TSC grades — not everything is A+
  - `RELEASE.md` must restate the ledger row as the Outcome section in prose. The ledger row remains canonical; the release body explains it.
  - If release notes or CHANGELOG wording are being authored, load the write skill and record it in the CDD Trace.

2.5. **Release notes — RELEASE.md**
  - Write `RELEASE.md` at repo root before tagging. This is the GitHub release body.
  - A release is a measured coherence delta. The release body must foreground what got more coherent and why, then detail the changes.
  - The release CI workflow uses `RELEASE.md` as the release body if present; otherwise it auto-generates from commit titles (sparse and unstructured — not acceptable).
  - `RELEASE.md` is committed in the release commit and consumed by CI. It stays in the repo until the next release overwrites it.
  - ❌ Let CI auto-generate release notes (just commit titles, no detail)
  - ❌ Jump straight into Fixed/Added/Changed (change-loggy, not coherence-framed)
  - ❌ Manually create the GitHub release after CI (race condition with workflow)
  - ✅ Outcome first, then changes, then proof
  - ✅ Write `RELEASE.md` → commit in release commit → CI uses it as release body

  ```markdown
  # RELEASE.md format

  ## Outcome

  Coherence delta: C_Σ {grade} (`α {grade}`, `β {grade}`, `γ {grade}`) · **Level:** `L{level}`

  One short paragraph:
  - what became more coherent
  - what surface is now more truthful / portable / durable / explicit

  ## Why it matters

  One short paragraph:
  - what incoherence this cycle targeted
  - why it mattered operationally or structurally
  - what changed in the system's behavior or authority model

  ## Fixed

  - **Short description** (#issue): what changed and why it matters.

  ## Added

  - **Short description** (#issue): what was added and what new capability or clarity it provides.

  ## Changed

  - **Short description** (#issue): what changed in behavior, authority, or shape.

  ## Removed

  - **Short description** (#issue): what was removed and what incoherence it eliminated.

  ## Validation

  - Deployed to [target], validated [specific check].
  - State what was proven coherent by the validation, not only that deployment succeeded.

  ## Known Issues

  - #N — description (if any found during validation)
  ```

2.5a. **Move cycle directories to release directory** (tagged release)
  - Move every per-cycle directory from `.cdd/unreleased/{N}/` to `.cdd/releases/{X.Y.Z}/{N}/`:
    ```bash
    mkdir -p .cdd/releases/X.Y.Z
    for dir in .cdd/unreleased/*/; do
      [ -d "$dir" ] || continue
      mv "$dir" .cdd/releases/X.Y.Z/
    done
    ```
  - Each `.cdd/releases/{X.Y.Z}/{N}/` directory carries the same role-prefixed files that lived in `.cdd/unreleased/{N}/` during the cycle (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, plus any cycle-specific extras) per CDD.md §Tracking + §5.3a.
  - Include the moves in the release commit
  - `.cdd/unreleased/` should be empty after the release commit
  - ❌ Leave cycle directories in `unreleased/` after tagging (lose the version association)
  - ❌ Move after the release commit (they should be part of the release snapshot)
  - ✅ Move before commit, include in the release commit

2.5b. **Docs-only disconnect (no tag)**

  A cycle that ships only documentation, protocol artifacts, or assessments — with no code change and no version bump — still requires a disconnect, but the disconnect is the merge commit on main, not a tag.

  Cases this covers:

  - Retroactive close-out for a previously-released version (no new code, only frozen artifacts and CHANGELOG row backfill — the cycle for `usurobor/tsc#27` is the canonical example)
  - Self-coherence reports and post-release assessments produced as their own cycle, not bundled with a code release
  - Skill patches that ship only as protocol changes
  - Doc-cleanup cycles (e.g. retiring stale references after a paradigm shift)

  For docs-only cycles:

  - `.cdd/unreleased/{N}/` moves to `.cdd/releases/docs/{ISO-date}/{N}/` where `{ISO-date}` is the merge commit's date in `YYYY-MM-DD` form (e.g. `.cdd/releases/docs/2026-05-08/27/`).
  - **No CHANGELOG ledger row.** The Release Coherence Ledger tracks tagged releases. γ records the cycle in `docs/gamma/cdd/docs/{ISO-date}/POST-RELEASE-ASSESSMENT.md` instead — same PRA structure as a tagged release, just keyed by date.
  - The merge commit hash IS the disconnect signal. No `scripts/release.sh` invocation. No version bump. VERSION is unchanged.
  - `scripts/check-version-consistency.sh` is not required to run (nothing version-stamped changed).

  ```bash
  # Docs-only disconnect at merge time (run as part of γ closure):
  ISO_DATE="$(date -u +%Y-%m-%d)"
  mkdir -p ".cdd/releases/docs/${ISO_DATE}"
  for dir in .cdd/unreleased/*/; do
    [ -d "$dir" ] || continue
    mv "$dir" ".cdd/releases/docs/${ISO_DATE}/"
  done
  ```

  - ❌ Leave docs-only cycle directories in `.cdd/unreleased/{N}/` after merge ("there's no release to move them to")
  - ❌ Force a synthetic version bump on a docs-only cycle just to fit the tagged-release flow
  - ❌ Skip the PRA because no tag was pushed
  - ✅ Move to `.cdd/releases/docs/{ISO-date}/{N}/` in the merge commit
  - ✅ γ writes `docs/gamma/cdd/docs/{ISO-date}/POST-RELEASE-ASSESSMENT.md`

  γ declares `docs-only` in the issue mode header (per `issue/SKILL.md` §Mode declaration); the disconnect path is selected at issue-creation time, not retrofitted at merge.

2.6. **Commit and signal readiness for δ tag**
  - **Tag naming convention:** use bare version numbers without `v` prefix: `3.15.1`, not `v3.15.1`. This matches VERSION file content, branch naming (`claude/3.15.0-22-...`), and snapshot directory names (`docs/gamma/cdd/3.15.0/`). Consistency across all version surfaces.
  - **Tag message generation:** δ creates annotated tags (not lightweight tags) with structured messages that include issue metadata, wave context when present, and CDD review artifacts when available. The tag message is automatically generated by `scripts/generate-release-tag-message.sh` and includes issue titles, labels, review rounds, and deterministic fallback when metadata is unavailable.
  - Commit: `git commit -m "release: X.Y.Z — summary"` (includes VERSION, manifests, CHANGELOG, RELEASE.md)
  - **Before any push that follows a rebase, run the eng/ship rebase-integrity gate** (see `eng/ship` § Rebase-Collision Integrity)
  - Push: `git push origin main`
  - **β signals "release ready for δ tag"** in `beta-closeout.md` — β does not execute `git tag` or push tags; δ creates the single annotated tag `X.Y.Z` with generated message as part of the disconnect release flow (CDD.md Phase 6 step 17). RELEASE.md remains the GitHub release body authority; generated tag messages are additive.
  - ❌ β pushes tag directly (conflicts with δ tag authority)
  - ❌ Mix `v` prefix and bare versions across tags (`v3.14.6` then `3.14.7`)
  - ❌ Commit without RELEASE.md (CI auto-generates sparse notes)
  - ✅ Commit (with RELEASE.md), push, signal readiness in close-out; δ handles tagging

2.6a. **Delete merged branches**
  - After push, delete remote branches that were merged into this release
  - `git branch -r --merged origin/main | grep -v main | grep -v HEAD | sed 's/origin\///' | xargs -I{} git push origin --delete {}`
  - This is mechanical cleanup, not judgment — if it's merged, it's dead
  - ❌ Leave merged branches accumulating ("someone might need them")
  - ✅ Clean up immediately after release push; history is in main

2.7. **Wait for release CI**
  - Release workflow builds binaries (linux-x64, macos-x64, macos-arm64)
  - Wait for completion before deploying
  - **δ owns release CI polling** per `operator/SKILL.md` §3.4 step 4 — after tag push, δ monitors release workflow completion via `gh run list --branch <tag>`. β waits for δ confirmation.
  - ❌ Deploy while CI still running (stale binary from previous release)
  - ✅ `gh run watch {id} --exit-status` then verify assets attached

2.8. **Deploy to target hosts**
  - Stop daemon, download binary, replace, start daemon
  - `cn deps restore` to sync cognitive substrate
  - Validate: `cn --version`, test an op that exercises the fix
  - ❌ Replace binary while daemon running (`Text file busy`)
  - ✅ `systemctl stop → cp → chmod → systemctl start → cn --version`

2.9. **Validate**
  - Version matches everywhere: binary, cn.json, tag
  - The specific fix or feature works end-to-end
  - Check telemetry for errors on first cycle
  - Validation must confirm the targeted coherence delta, not only binary/version correctness.
  - ❌ "It deployed" (no functional validation)
  - ✅ Trace telemetry, confirm receipts/artifacts exist, agent responds correctly

2.10. **CDD Trace update**
  - Update the primary branch artifact's CDD Trace with the release row:
    - artifact: tag / CHANGELOG / release artifact
    - skills loaded: release, plus write if used
    - decision: released version X.Y.Z
  - If the triadic protocol is active, β also writes:
    - release evidence in `.cdd/unreleased/{N}/beta-closeout.md`
    - β close-out narrative in `.cdd/unreleased/{N}/beta-closeout.md` for γ to read (separate from `beta-review.md`, which carries the round-by-round verdicts; see `CDD.md` §Tracking canonical filename table)
  - γ writes the post-release assessment after β's release and close-out are complete
  - For triadic cycles, the primary branch artifact is `.cdd/unreleased/{N}/self-coherence.md` — it carries the trace through step 7a; β records the review verdict in `.cdd/unreleased/{N}/beta-review.md` (steps 8–9) and the close-out + release evidence in `.cdd/unreleased/{N}/beta-closeout.md` (step 10).

## 3. Rules

3.1. **All version strings must agree**
  - VERSION, cn.json, package manifests, tag, binary output, CHANGELOG
  - `scripts/check-version-consistency.sh` validates this mechanically
  - ❌ Tag is 3.15.1 but VERSION says 3.15.2 (forgot to commit before tagging)
  - ✅ `scripts/check-version-consistency.sh` passes before commit

3.2. **CI must pass before tag**
  - If CI fails after tag, fix and force-push tag (amend release commit)
  - ❌ Ship with known test failure ("we'll fix it next patch")
  - ✅ Fix the test, amend, force-push tag, wait for green

3.3. **Never deploy stale binaries**
  - Release CI builds from the tag; wait for it to complete
  - Clear cached downloads (`--clobber` or `rm` old artifacts)
  - ❌ `cn update` downloads cached v3.5.1 binary from /tmp
  - ✅ Fresh download after CI completes, verify binary version after install

3.4. **Deploy includes deps restore**
  - Binary update alone leaves cognitive substrate stale
  - `cn deps restore` syncs doctrine, mindsets, skills from packages
  - ❌ Binary is v3.9.1 but skills are from v3.8.0 (package drift)
  - ✅ `cn deps restore` immediately after binary update

3.5. **Validate with the specific fix**
  - Don't just check version — exercise the feature or fix
  - ❌ "Version looks right, ship it"
  - ✅ For #46: send a message, check telemetry for receipts, confirm artifacts exist

3.6. **If CI fails post-tag, amend don't re-tag**
  - `git commit --amend --no-edit && git tag -d X.Y.Z && git tag X.Y.Z && git push --force origin main X.Y.Z`
  - Tags are **bare** (`X.Y.Z`, never `vX.Y.Z`) — see §2.6 and CDD.md §5.3a
  - Keeps a clean single tag, single release commit
  - ❌ `git tag -d vX.Y.Z && git tag -a vX.Y.Z && git push --force --tags` (mixes legacy v-prefix; force-pushes all tags)
  - ❌ `3.9.1`, `3.9.1-fix`, `3.9.1-fix2` (proliferating tags)
  - ✅ One bare tag, amended until CI green

3.7. **RELEASE.md must exist before tag**
  - The release commit must include `RELEASE.md` at repo root
  - If `RELEASE.md` is missing, the tag must not be pushed (CI will auto-generate sparse notes)
  - ❌ Tag without RELEASE.md ("I'll update the release body manually later")
  - ✅ Write RELEASE.md → include in release commit → tag → push

3.8. **TSC scoring in CHANGELOG — honest-grading rubric**

  Rate α (pattern), β (relation), γ (process) for each release. Honest grades — not everything is A+. The rubric below makes the grading reproducible across releases and reviewers.

  **Per-axis rubric:**

  | Grade | Numeric | Meaning |
  |---|---|---|
  | **A**   | 4.0 | met all ACs, no protocol skips, zero binding findings of the wiring or honest-claim class |
  | **A-**  | 3.7 | met all ACs, ≤1 binding finding, all findings non-blocking |
  | **B+**  | 3.3 | met all ACs, ≥2 binding findings or one round of RC |
  | **B**   | 3.0 | met core ACs; partial-protocol release (e.g. missing RELEASE.md at tag time) **or** 2+ rounds of RC |
  | **C+**  | 2.3 | shipped but missing one of: ledger row, PRA, provenance attachment |
  | **C**   | 2.0 | partial-protocol release with multiple drift items |
  | **< C** | < 2 | re-open and remediate; do not close |

  **Closure-gate override (fires before geometric-mean computation).** If any artifact required by `gamma/SKILL.md` §2.10 closure gate (per `CDD.md` §5.3b ownership matrix) is absent from the cycle's `.cdd/unreleased/{N}/` directory at merge time — specifically any of `alpha-closeout.md`, `beta-closeout.md`, or `gamma-closeout.md` for a triadic cycle — then `C_Σ` is forced to `<C` regardless of per-axis math. Cycle disposition is "open and remediate"; the geometric mean is not computed. This override is verified mechanically by `scripts/validate-release-gate.sh --mode pre-merge` (see AC1, issue #339). The math below applies only after the closure gate passes.

  **Letter normalization (`<C` and `C−`).** The rubric table uses `<C` as the grade label. Prior CHANGELOG, PRA, and alpha-closeout artifacts for cycles #331, #333, and #334 used `C−` for the same disposition. These are the same grade: `C−` is the operator-visible projection of `<C`. Both mean "open and remediate; do not close." Authors may use either form; the CHANGELOG column uses `<C` for rubric fidelity and `C−` is accepted as equivalent in prose artifacts.

  **C_Σ** is the geometric mean of the three numeric grades; report the closest letter grade (e.g. (3.7 · 3.3 · 4.0)^(1/3) ≈ 3.66 ⟹ A-).

  **Score the release, not the intent.** A retroactive partial-protocol close-out earns C+ honestly (e.g. `usurobor/tsc` cycle 27: v0.4.0 retroactive close-out scored α B / β C+ / γ C / C_Σ C+ — the grade reflects what shipped, not the goodwill of fixing it later).

  **Configuration-floor clause.** Cycles run under `operator/SKILL.md §5.2` (single-session δ-as-γ via Agent tool) cap the γ axis at **A−** regardless of execution quality, because γ/δ separation is structurally absent. The A− γ floor applies from the cycle's merge forward; it is not retroactive. The configuration must be recorded explicitly in `gamma-closeout.md` (see AC6 of any cycle introducing or running under §5.2). Applying the A− cap is operator-honest discipline: if `gamma-closeout.md` does not declare §5.2, the rubric treats the cycle as §5.1 and the cap does not apply.

  **CI-red cap clause.** Cycles with red CI on the merge commit cap the γ axis at **C** (one band below current floor). Cycles that proceed to close-out without verifying CI cap the γ axis at **B−**. CI is a mechanical gate; failing or skipping it reflects structural γ-axis failure in cycle coordination, not judgment calls.

  **False-gap peer-enumeration clause.** Cycles that ship on a false-gap premise (where the §Gap assertion is contradicted by existing surfaces that basic peer-enumeration would have discovered) cap the γ axis at **B** regardless of whether β caught and corrected the false claim. A cycle requiring fix-rounds to reconcile γ's unverified gap framing reflects failure of the peer-enumeration discipline per `gamma/SKILL.md` §2.2a, not implementation or review failure.

  - ❌ Every release is A+ (grade inflation, no signal)
  - ❌ Score the intent rather than what shipped
  - ❌ Round up because "the team worked hard"
  - ❌ Rate a §5.2 cycle's γ axis above A− (grade inflation, structural absence of γ/δ separation unaccounted)
  - ✅ "β: A- — DUR skills synced but README.md still references old framing"
  - ✅ "γ: C — partial-protocol release; retroactive close-out planned"
  - ✅ "γ: A− — §5.2 cycle; γ/δ separation absent; A− γ floor applied per §3.8 configuration-floor clause"

---

## 4. Kata

**Scenario:** A bugfix is ready on main. Execute the release.

1. Readiness check — CI green, no unmerged branches
2. Version decision — patch, minor, or major
3. VERSION bump + stamp-versions.sh + check-version-consistency.sh
4. CHANGELOG ledger row with honest TSC grades
5. RELEASE.md with Outcome mirroring the ledger row
6. Commit, tag (bare version), push
7. Wait for release CI, deploy, validate with the specific fix
8. β close-out (γ writes PRA separately)

**Verify:** Does `scripts/check-version-consistency.sh` pass? Does RELEASE.md start with the coherence delta? Does validation confirm the targeted incoherence is closed?
