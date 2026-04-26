# β Close-Out — 3.60.0 (#235, PR #276)

## Cycle context

- **Cycle:** #235 (B2 — `cn build --check` validates command entrypoints + skill SKILL.md presence; single-issue bundle, no closures).
- **PR:** #276 (https://github.com/usurobor/cnos/pull/276).
- **Branch:** `claude/alpha-tier-3-skills-M8Vce` (non-canonical per CDD §4.2 — same harness pattern as #230's `…-IZOsO`; see §Findings).
- **Author identity:** `α` (squash-merge committer recorded as `usurobor` on origin).
- **Reviewer identity:** `β <beta@cdd.cnos>` — shared GitHub identity with α (per CDD §1.4 review note); review verdict posted as comment per `cdd/review/SKILL.md` §7.1 with explicit "shared identity" annotation.
- **Tier 3 skills loaded by β:** `eng/go` (incl. §2.17 "one parser per fact" + diff-local generalization, §2.18 dispatch boundary), `eng/test` (proof-depth ladder + negative-space mandatory), plus `docs/alpha/DESIGN-CONSTRAINTS.md` (§1, §2, §3.2, §6 named/touched by the issue) and `docs/alpha/architecture/INVARIANTS.md` (T-002, T-004 named in the issue).
- **Architecture check (§2.2.14):** not strongly triggered — additive validator inside one module, no boundary moves; spot-checked anyway, all axes `yes` or `n/a`, no `no`.

## Review context

### Round 1 — APPROVED at SHA `f7d27b4`

Walked `cdd/review/SKILL.md` §2.0 (issue contract) → §2.1 (diff internal coherence + mechanical scans) → §2.2 (sibling/context audit incl. §2.2.1a input enumeration, §2.2.5 cross-ref, §2.2.8 authority surface, §2.2.9 module-truth, §2.2.10 contract-implementation confinement, §2.2.13 design constraints, §2.2.14 architecture check spot-check) → §2.3 (verdict).

**Findings table:** empty. The diff's contract surface (`pkgbuild.checkCommandEntrypoints`, `pkgbuild.checkSkillDirectories`, `pkgbuild.containsSkillMd`) plus its surrounding `pkgbuild`/`activation`/`cli` modules cohere on first pass.

**Notes (deferred to γ; not findings against the diff):**

| # | Surface | Pattern |
|---|---------|---------|
| N1 | `pkgbuild.PackageManifest` + `ParseManifestData` (minimal, local) | Pre-existing parallel parser to `pkg.PackageManifest` + `pkg.ParseInstalledManifestData`. α's PR body §Known debt names this; eliminating it changes `DiscoveredPackage.Manifest`'s static type and ripples through `BuildOne`/`UpdateIndex`/`UpdateChecksums`. The diff itself uses the canonical `pkg.ParseFullManifestData` for the new code path — `eng/go §2.17` "one parser per fact" + diff-local generalization is honored on the new surface. |
| N2 | `pkgbuild.BuildOne` does not gate on `CheckOne` | `cli/cmd_build.go::Run` dispatches `runBuild` and `runCheck` as separate paths. `cn build` (without `--check`) tarballs unvalidated. CI runs `--check` separately so the gate exists in CI; locally the safety is opt-in. Pre-existing structure — the prior `no content class` check had the same shape — not introduced by this PR. |
| N3 | `checkCommandEntrypoints` symlink semantics | `os.Stat` follows symlinks; an entrypoint that is a symlink to a regular file outside the package root passes `IsRegular()`. The packlist + tarball preserve symlinks, so install-host `/etc/passwd` could become the resolved target. Pre-existing concern at the runtime/dispatch layer too. AC1 wording ("regular file under the package root") could be tightened later via `os.Lstat` + symlink-target containment. Out of scope for this issue. |
| N4 | Issue body cites stale path `docs/alpha/INVARIANTS.md` | Canonical is `docs/alpha/architecture/INVARIANTS.md`. PR body cites the constraints by name (T-002, T-004) without a path, so the implementation is unaffected. γ-side issue-quality observation. |

CI state at the time of review: 7/7 success completed at `f7d27b4` (`go`, `Package/source drift (I1)`, `Protocol contract schema sync (I2)`, `kata-tier1`, `kata-tier2`, `notify`×2). Approval was unconditional on CI grounds.

### β positive checks (the parts that earned the approval)

- **AC1 contract-implementation confinement (review §2.2.10):** `checkCommandEntrypoints` rejects empty entrypoints, missing files, non-regular targets, and path-traversal (`../`). Three of the four extend beyond strict AC1 wording but are necessary for the "regular file under the package root" claim to be earned. Verified the path-traversal logic against three edge cases mentally: `..` alone, `../foo`, `./../foo` — all caught by `rel == ".." || strings.HasPrefix(rel, "../"+sep)`. Absolute-path entrypoint (`/etc/passwd`) is collapsed by `filepath.Join` into `<pkgDir>/etc/passwd` and flagged as non-existent — safe.
- **AC2 authority-surface alignment (review §2.2.8):** `checkSkillDirectories` reads `pkgtypes.ClassSkills` (the canonical content-class constant) and walks `skillsRoot` for `SKILL.md` — exactly the surface `internal/activation/index.go::discoverPackageSkills` walks at runtime. Build-time check and runtime activation now agree on what counts as a valid skill. The "subtree contains SKILL.md" rule (rather than "this directory has SKILL.md") correctly exempts namespace containers like `cnos.eng/skills/eng/`.
- **AC3 negative-space coverage (eng/test §2.7 / §3.6):** all four boundaries the validator names have explicit "must not happen" tests — `TestCheckOneEntrypointMissing`, `TestCheckOneEntrypointIsDirectory`, `TestCheckOneEntrypointEscapesPackageRoot`, `TestCheckOneSkillDirMissingSkillMd`. The two false-positive guards (`TestCheckOneSkillContainerNamespaceExempt`, `TestCheckOneSkillResourceSubdirExempt`) test what the validator must **not** flag.
- **AC4 evidence depth (review §2.3.6):** α's verification path is local deliberate-break repro on a copy of the repo + the existing I1 job's exit-code semantics. The issue's "in a CI dry-run" wording is tighter than what landed (no on-CI deliberate-break test added), but the substance — both new error classes flow through the unchanged exit-code → I1-job-fail chain — is verified end-to-end. Acceptable for the AC.
- **eng/go §2.17 "one parser per fact" + diff-local generalization:** the new code reads `cn.package.json` once via `pkg.ParseFullManifestData`, the canonical pure parser also used by `internal/discover/discover.go` for runtime command discovery. Build-time validator and runtime command dispatch now share one full-manifest parser. No parallel parser introduced by the diff. (The pre-existing `pkgbuild.ParseManifestData` minimal-shape parser is left alone with explicit known-debt note — see N1.)
- **CI green on head commit (alpha §2.6 / β §3.7):** all 7 required checks completed `success` at the reviewed SHA before approval. No checks pending, queued, or re-running.

## Release evidence

- **Merge:** PR #276 squash-merged via GitHub MCP at `d814e1665228237fa52922f4855f8cd502f6c09f`. Merge commit message preserves the diff intent + AC mapping and carries `Closes #235`.
- **Issue auto-close:** #235 transitioned to `state: closed` automatically via the `Closes #235` reference in the squash-merge commit message; no manual close needed.
- **Version decision:** 3.59.3 → 3.60.0 (minor bump). New behavior shipped (`cn build --check` rejects new failure classes); existing valid packages remain valid (CI's I1 job confirmed `./cn build --check` passes on all 5 live packages at `f7d27b4`); backwards-compatible additive validation; no API/ABI breakage.
- **Version stamping:** `VERSION` updated to `3.60.0`; `scripts/stamp-versions.sh` derived `cn.json` (3.59.3 → 3.60.0) and `src/packages/cnos.cdd/cn.package.json`, `src/packages/cnos.core/cn.package.json`, `src/packages/cnos.eng/cn.package.json` (each: `version` and `engines.cnos`); `scripts/check-version-consistency.sh` PASSED — 8/8 surfaces agree on `3.60.0`.
- **Release artifacts on main:** `VERSION`, `cn.json`, three package manifests, `CHANGELOG.md` ledger row, `RELEASE.md`. Released as commit `cce905e` on `origin/main`. `.cdd/unreleased/` was empty at release time (no close-outs to migrate per release/SKILL.md §2.5a).
- **Tag:** `3.60.0` (bare per CDD §5.3a; `v`-prefix is legacy). **Local only — origin push deferred to δ; see §Tag-push deferral.**
- **Local sanity at release time:** `go test ./internal/pkgbuild/...` PASSED (`ok`, 0.034s); `go vet ./...` clean.
- **Release CI (`release.yml`) expected post-tag:** triggers on tag push — builds the 4-platform matrix (linux-x64, linux-arm64, macos-x64, macos-arm64) with `-ldflags "-X main.version=${{ github.ref_name }}"` from #230, runs `kata-tier1` per platform on the shipped binary, generates `checksums.txt`, attaches binaries + `packages/*.tar.gz` + `index.json` to the GitHub release. `RELEASE.md` (committed in the release commit) becomes the release body via `softprops/action-gh-release@v1` + `body_path: RELEASE.md`.
- **Release smoke (`release-smoke.yml`) expected post-tag:** fires on `release: published` event after the release job uploads assets — exercises the full bootstrap chain on each of 4 platforms. rc=0 → `::notice::`; rc=2 (offline) → `::warning::`; rc=1 → `::error::`.
- **Validation expected post-tag (operator-side once δ pushes the tag):**
  - `release.yml` matrix completes green on all four platforms.
  - GitHub release `3.60.0` exists with `cn-{linux-x64,linux-arm64,macos-x64,macos-arm64}` binaries, `checksums.txt`, `packages/*.tar.gz`, `packages/index.json`.
  - `release-smoke.yml` matrix completes per-platform — chain integrity confirmed end-to-end.
  - On any operator-side `cn build --check` invocation against `cnos.cdd` with a deliberately-broken `commands.fake-cmd.entrypoint` and a `skills/orphan-skill/` containing only a stray `notes.md`: both new error classes are surfaced with exit code 1 (mirrors α's local AC4 dry-run).

## Cycle findings

Voice per CDD §1.4 β step 9: factual observations and patterns only. Triage is γ's at PRA time.

### CDD §9.1 trigger fired — "avoidable tooling/environmental failure" + "loaded skill failed to prevent a finding"

**What happened.** On the new-branch transition event for `origin/claude/alpha-tier-3-skills-M8Vce`, β's first synchronous PR-existence check used `mcp__github__list_pull_requests state=open head=usurobor:claude/alpha-tier-3-skills-M8Vce` and got back `[]`. β reported "no PR open yet" to the operator. PR #276 was already open at that exact head ref (created at the same wall-clock moment as the branch push). The operator caught the gap within one polling tick and prompted β to re-check; an unfiltered `mcp__github__list_pull_requests state=open` returned the PR immediately. β then read the diff and the rest of the cycle proceeded normally.

**Pattern.** Three filter forms were attempted; all three failed silently:

1. **MCP `head=owner:branch` filter**: returned `[]` even though the PR's `head.ref` exactly matches the branch name. Likely cause: the slash inside `claude/alpha-tier-3-skills-M8Vce` breaks the filter's URL-encoding or query semantics; the filter is unreliable for in-repo branches with slashes. Behavior is the bug regardless of cause.
2. **`gh pr list --search 'closes:#235 OR refs:#235 OR fixes:#235'` (in the polling Monitor)**: `closes:`/`refs:`/`fixes:` are not valid GitHub search qualifiers. The query returns empty for any input.
3. **`gh pr list --state open --json number,headRefName,body,title | jq … test('#235|/235\\b')` (fallback in the polling Monitor)**: would have matched, but the Monitor only polls every 60s, so on the new-branch transition there's a ≤60s window before the next tick.

The unifying failure mode: β trusted a narrowly-scoped query result as authoritative without falling back to a broad sync.

**Skill that was loaded but did not prevent the finding.** `CDD.md §Tracking` was loaded by β at intake. The relevant rule there is "synchronous-baseline-pull is a precondition of transition-only polling" — currently named **only for session-start**, not for per-event handling after a transition fires. The reference table of query forms also implies that filtered queries are reliable.

**Surfaces affected.** `CDD.md §Tracking` (reference query forms + baseline rule scope). Same general surface as 3.59.0's "first-iteration absorption" finding, narrower instance.

**Recommended `§Tracking` patch (β observation; γ owns triage):**

> The MCP `head=owner:branch` filter is unreliable for in-repo branches whose name contains a slash; do not use it as a primary query. The `gh pr list --search` qualifier set does not include `closes:`/`refs:`/`fixes:`; those query strings always return empty. On any new-branch transition event from the polling Monitor, do a synchronous broad open-PR list (`mcp__github__list_pull_requests state=open` with no head/search filter) and scan client-side for `(head.ref == new-branch) OR (body matches (?i)\b(closes\|fixes\|resolves\|refs)\s*#N\b)`. This is the same shape as the existing baseline rule, extended from session-start to per-event handling.

### Harness branch-naming drift vs CDD §4.2 (recurring from 3.59.0)

α's branch was `claude/alpha-tier-3-skills-M8Vce`. CDD §4.2 canonical format is `{agent}/{issue}-{scope}` (or `{agent}/{version}-{issue}-{scope}` when version is known) — the issue number is part of the canonical branch identity. The harness assigned the name; α did not choose. Identical pattern shape to #230's `claude/alpha-tier-3-skills-IZOsO`. The 3.59.0 close-out already surfaced this; appears to remain unresolved upstream (harness ↔ CDD §4.2 contract). Surfaces affected: CDD §4.2 ↔ harness/dispatch convention.

### Issue body's stale doc-path reference

Issue #235 body cites "`docs/alpha/INVARIANTS.md` T-002, T-004" in the Active design constraints section. Canonical path is `docs/alpha/architecture/INVARIANTS.md`. PR body cites the constraints by name without a path, so the implementation is unaffected; this is a γ-side issue-quality observation about the dispatch artifact, not an α responsibility. Surfaces affected: issue authoring ↔ canonical doc paths. Note: 3.59.0's β close-out also surfaced stale doc-path references on `main` (in `eng/go/SKILL.md` and `pkg.go:134`); whether those have been corrected or remain pre-existing on this release's `main` was not re-audited (out of scope for #235 review).

### Mechanical ratio

Findings: 0 total. Mechanical ratio undefined (no findings to ratio). Recorded for completeness — the §9.1 ratio trigger requires ≥10 findings to fire on ratio alone.

### Round economics

1 review round — meets the §9.1 default target of ≤2 with margin. The diff was authoring-time clean: every AC mapped to code + test evidence in the PR body §Self-coherence, the canonical parser was reused at the call site that needed full-manifest commands, the negative-space tests (path-traversal, directory-as-entrypoint, container-namespace, resource-subdir) were already named in the test file. β's review walked the surface to confirm rather than to discover.

## Invariant audit

Project design constraints checked per `cdd/review/SKILL.md` §2.2.13 against `docs/alpha/DESIGN-CONSTRAINTS.md`:

| Constraint | Pre-cycle state | Post-cycle state | Direction |
|---|---|---|---|
| §1 Source of truth (each fact in one place) | Build-time validation only checked content-class directory presence; the manifest's commands list and the filesystem skill layout were not cross-checked. Activation/discover already trusted filesystem-as-authority. | Build-time validator and runtime activation now read from the same authority surface (`pkg.ClassSkills`, `pkg.ParseFullManifestData`, `pkg.CommandEntries()`); a single `pkg`-package fact has a single canonical parser path on the new code. | tightened |
| §2.1 One package substrate (`cn.package.v1`) | All packages on `cn.package.v1`; one manifest format, one build pipeline. | Same. No new manifest types or schema fields introduced. | preserved |
| §2.2 Source / artifact / installed clarity | `src/packages/<name>/` (authored) → `dist/packages/*.tar.gz` (built; not in git per 3.58.0) → `.cn/vendor/packages/<name>/` (installed). Build-time check did not enforce manifest-claim ↔ filesystem-fact link. | Build-time check now refuses to produce a tarball whose declared commands point at missing files or whose top-level skill dirs lack `SKILL.md`. The source-side authority surface is now mechanically self-consistent. | tightened |
| §3.1 Git-style subcommands | Not touched. `cn build`, `cn build --check`, `cn build clean` retain their shape. | Same. | preserved |
| §3.2 Dispatch boundary (`cli/` owns dispatch only) | `internal/cli/cmd_build.go::runCheck` is a thin wrapper that invokes `pkgbuild.CheckOne` per package and renders results. | Same — all new validation logic is in `internal/pkgbuild/build.go` (`checkCommandEntrypoints`, `checkSkillDirectories`, `containsSkillMd`); `cli/cmd_build.go` is unchanged. The `cli/` thin-wrapper invariant is preserved. | preserved |
| §4.1 Surface separation (skills / commands / orchestrators / providers) | Not touched. | Same — `commands` declared in the manifest are still distinct from `skills` discovered on the filesystem; the validator reinforces both surfaces independently. | preserved |
| §4.2 Registry normalization (one runtime descriptor) | Not touched. The build-check produces a separate issue list, not a runtime registry. | Same. | preserved |
| §5.0 OCaml deprecated, no edits | No OCaml files modified by this cycle. | Same. | preserved |
| §6.1 Reason to change (one per boundary) | `pkgbuild` owns build-time validation; `activation` owns runtime skill discovery; `discover` owns runtime command discovery. | Same — boundaries reinforced (the new validator stays in `pkgbuild`; it does not bleed into `activation` or `cli`). | preserved |
| §6.2 Policy above detail (kernel owns policy) | Not touched. | Same. | preserved |
| §6.3 Degraded-path visibility | Pre-cycle: missing entrypoints / orphan skill dirs were silent at build time, producing broken tarballs that failed only at install or runtime. | Post-cycle: both classes surface as named issues at `cn build --check` time, exit-coded, and routed through I1's CI gate. The previously-silent failure class is now testable and reviewable on every PR. | tightened |

Cross-reference to `docs/alpha/architecture/INVARIANTS.md`:

| Invariant | Status |
|---|---|
| INV-001 (one package substrate) | preserved — no new manifest format introduced |
| INV-003 (commands/providers/orchestrators/skills distinct) | preserved — validator treats commands and skills as independent surfaces, mirrors the runtime distinction |
| T-002 (kernel remains minimal and trusted; cli/ owns dispatch only) | preserved — verified via diff: zero `cmd_*.go` changes; all new logic in `internal/pkgbuild/` |
| T-004 (source/artifact/installed explicit) | preserved-tightened — manifest-claim ↔ filesystem-fact link is now enforced at build time, eliminating one class of source-side incoherence that previously leaked into the artifact |

Architecture check (`cdd/review/SKILL.md` §2.2.14) result on the merged diff: 6 of 7 axes `yes`, 1 `n/a` (Registry normalization — not a registry-touching change). No `no` answers; the §2.3.5 architecture-gate was not blocked.

## Tag-push deferral (CDD §1.4 β-step 8)

`git push origin 3.60.0` returned `HTTP 403` from the sandbox-side git proxy:

```
To http://127.0.0.1:46731/git/usurobor/cnos
remote: ...
error: RPC failed; HTTP 403 curl 22 The requested URL returned error: 403
send-pack: unexpected disconnect while reading sideband packet
fatal: the remote end hung up unexpectedly
Everything up-to-date
```

`git push -u origin main` from the same session succeeded (release commit `cce905e` is on `origin/main`) — the 403 is specific to tag refs from this β session, not a general remote-write block. Identical environment behavior to 3.59.0's tag-push deferral. Per CDD §1.4 β-step 8:

> If tag push fails due to env constraints (e.g. sandbox HTTP 403), commit all release artifacts to main and defer tag push to δ (operator) — do not block closure on it.

Status:

- All release artifacts are on `origin/main` at `cce905e`: `VERSION=3.60.0`, three `cn.package.json` files (3.60.0), `cn.json` (3.60.0), `CHANGELOG.md` ledger row, `RELEASE.md`.
- The tag `3.60.0` exists locally at `cce905e` but is NOT on origin (`git ls-remote --tags origin 3.60.0` would return empty).
- **Deferred to δ:** `git push origin 3.60.0` from a session with non-sandboxed write access. Once the tag lands on origin, `release.yml` triggers automatically (4-platform matrix build → GitHub release with `RELEASE.md` body → `release-smoke.yml` fires on `release: published`).
- β does not block closure on this step. β close-out (this file) lands directly to `origin/main` per CDD §5.3a; γ's PRA and α's close-out can proceed in parallel.

## β verdict trail

- Round 1 (and only round): APPROVED — https://github.com/usurobor/cnos/pull/276 (review submitted via `mcp__github__pull_request_review_write event=COMMENT` at SHA `f7d27b4` per shared-identity rule, `cdd/review/SKILL.md` §7.1).
- Merge commit: `d814e1665228237fa52922f4855f8cd502f6c09f` (squash of α PR #276; merged via `mcp__github__merge_pull_request method=squash`).
- Release commit: `cce905e…` (`release: 3.60.0 — cn build --check validates entrypoints + skill paths (#235, PR #276)`).
- Tag (local only, pending δ push): `3.60.0` → `cce905e`.
- Issue closed: #235 `state: closed` (auto-close via `Closes #235` in squash-merge commit).
- β process finding (this close-out, §Cycle findings): recommended `CDD.md §Tracking` patch deferred to γ for triage at PRA time.

