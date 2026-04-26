# β Close-Out — 3.59.0 (#230, PR #274)

## Cycle context

- **Cycle:** #230 (B1 distribution chain honesty — version-aware reinstall + release-bootstrap smoke; subsumes #238).
- **PR:** #274 (https://github.com/usurobor/cnos/pull/274).
- **Branch:** `claude/alpha-tier-3-skills-IZOsO` (non-canonical per CDD §4.2 — see §Findings).
- **Author identity:** `α (claude-opus-4-7) <alpha@cnos.local>`.
- **Reviewer identity:** `β <beta@cdd.cnos>` — shared GitHub identity with α (per CDD §1.4 review note); review verdicts posted as comments per `cdd/review/SKILL.md` §7.1 with explicit "shared identity" annotation.
- **Tier 3 skills loaded by β:** `eng/go`, `eng/test`, `eng/tool`, `cnos.core/skills/design` (architecture check §2.2.14 active because the diff touches source/artifact/installed flow), plus `docs/alpha/DESIGN-CONSTRAINTS.md` (§1, §2.2, §3.2, §6.3 named in the issue).

## Review context

### Round 1 — REQUEST CHANGES at SHA `6cbd225`

Walked `cdd/review/SKILL.md` §2.0 (issue contract) → §2.1 (diff internal coherence + mechanical scans) → §2.2 (sibling/context audit incl. §2.2.1a input enumeration, §2.2.5 cross-ref, §2.2.8 authority surface, §2.2.9 module-truth, §2.2.13 design constraints, §2.2.14 architecture check) → §2.3 (verdict).

Findings posted (https://github.com/usurobor/cnos/pull/274#issuecomment-4321807808):

| # | Severity | Type | Surface | Pattern |
|---|----------|------|---------|---------|
| F1 | C | judgment | `doctor.checkPackages` parsed `cn.package.json` inline via anonymous-struct unmarshal | Two parsers for the same fact ("installed package version") — `pkg.ParseInstalledManifestData` was introduced in the same PR but `doctor` didn't consume it. Authority-surface conflict per `cdd/review/SKILL.md` §2.2.8; parse/read-purity violation per `eng/go/SKILL.md` §2.17; one-source-of-truth violation per `design/SKILL.md` §3.2. |
| F2 | B | judgment | `scripts/smoke/90-release-bootstrap.sh` header (line 14-15) | Header claimed "verify SHA-256 against checksums.txt"; actual code verifies against `index.json[name][version].sha256` at line 207-211. Downloaded `checksums.txt` was unused. Two-authority comment vs implementation drift. |
| F3 | B | judgment | `doctor.checkPackages` "unparseable manifest" stale-state branch | `(installed X, locked Y)` and `(no manifest)` had tests; `(unparseable manifest)` did not. Negative-space coverage gap per `eng/test/SKILL.md` §2.7 / §3.6. |
| F4 | A | mechanical | `scripts/smoke/90-release-bootstrap.sh:127-129` | `BIN_VERSION_OUT=$(... \|\| true)` captured then never read; `if !` guard dead under `\|\| true`. Comment misnames the actual check. |

CI state at the time of review: 7/7 success completed at `6cbd225`. Approval was withheld solely on coherence findings, not on red CI.

### Round 2 — APPROVED at SHA `93ea1d6`

α pushed one round-2 commit `93ea1d6 α #230 round 2: address β findings F1–F4` — 3 files, 50 insertions / 17 deletions, strictly within the F1–F4 anchors.

| # | Fix verification |
|---|------------------|
| F1 | `doctor.go` adds `import "github.com/usurobor/cnos/src/go/internal/pkg"` and replaces the anonymous-struct unmarshal with `pm, err := pkg.ParseInstalledManifestData(manifestData)`. Comment cites eng/go §2.17 + design §3.2 explicitly. Import graph remains cycle-free (`pkg` → stdlib only; `doctor` → `pkg`, `activation`; `restore` → `pkg`). |
| F2 | Header rewritten: "verify each tarball's SHA-256 against the index.json entry — the same authority `cn deps restore` itself trusts at runtime". `checksums.txt` download removed; replacement comment names `release.yml`'s binary-asset checksums file as the unrelated authority that was being conflated. Single source of truth restored. |
| F3 | `TestRunAllPackageManifestUnparseable` added — writes `{not valid json` as the manifest body, asserts `StatusFail` with `"unparseable manifest"` in the diagnostic. Mirrors `TestRunAllPackageManifestMissing` shape. |
| F4 | `BIN_VERSION_OUT=$(... \|\| true)` replaced with `"$BIN" status >/dev/null 2>&1 \|\| true`. Comment now honestly describes what the line guards against (missing loader / wrong arch / lost exec bit). |

No new findings introduced. No drive-bys. CI 7/7 success completed at `93ea1d6`. Approval posted (https://github.com/usurobor/cnos/pull/274#issuecomment-4321845284) provisionally on CI mid-flight, finalized once all 7 runs reached `completed`/`success`.

## Release evidence

- **Merge:** PR #274 squash-merged via GitHub MCP at `9980e3f5ff6e7221c8f29463443d526ba8a5cb6d`. Merge commit message preserves the round-1 + round-2 narrative and carries `Closes #230`.
- **Issue auto-close:** #230 transitioned to `state: closed`, `state_reason: completed` at the merge time; no manual close needed.
- **#238 already closed-as-duplicate** (2026-04-25) before this cycle — PR body's "Subsumes #238" wording is informational.
- **Version decision:** 3.58.0 → 3.59.0 (minor bump). New behavior shipped (version-aware reinstall, doctor stale-vendor diagnostic, release-bootstrap smoke gate); backwards-compatible additive `Version` field on `pkg.PackageManifest`; no API/ABI breakage.
- **Version stamping:** `VERSION` updated to `3.59.0`; `scripts/stamp-versions.sh` derived `cn.json` (3.58.0 → 3.59.0) and `src/packages/cnos.core/cn.package.json`, `src/packages/cnos.cdd/cn.package.json`, `src/packages/cnos.eng/cn.package.json` (each: `version` and `engines.cnos`); `scripts/check-version-consistency.sh` PASSED — 8/8 surfaces agree on `3.59.0`.
- **Release artifacts on main:** `VERSION`, `cn.json`, three package manifests, `CHANGELOG.md` ledger row, `RELEASE.md`, `.cdd/releases/3.59.0/beta/CLOSE-OUT.md`. No close-outs in `.cdd/unreleased/` to migrate (empty at release time).
- **Tag:** `3.59.0` (bare per CDD §5.3a; `v`-prefix is legacy).
- **Release CI:** `release.yml` triggers on tag push — builds the 4-platform matrix with the new `-ldflags "-X main.version=${{ github.ref_name }}"` stamp from #230, runs `kata-tier1` per platform on the shipped binary, generates `checksums.txt`, attaches binaries + `packages/*.tar.gz` + `index.json` to the GitHub release. `RELEASE.md` (committed in the release commit) becomes the release body — the "Check for release notes" step picks up `RELEASE.md` and uses `softprops/action-gh-release@v1` with `body_path: RELEASE.md`.
- **Release smoke:** `release-smoke.yml` (introduced in this cycle) fires on the `release: published` event after the release job uploads assets — exercises the full bootstrap chain on each of linux-x64, linux-arm64, macos-x64, macos-arm64. rc=0 → `::notice::`; rc=2 (offline) → `::warning::`; rc=1 → `::error::` (release flagged).
- **Validation expected post-tag:**
  - `release.yml` matrix completes green on all four platforms.
  - GitHub release `3.59.0` exists with `cn-{linux-x64,linux-arm64,macos-x64,macos-arm64}` binaries, `checksums.txt`, `packages/*.tar.gz`, `packages/index.json`.
  - `release-smoke.yml` matrix completes — each platform downloads the just-shipped `cn-<platform>`, verifies tarball SHA-256 against `index.json`, runs `cn init` + `cn setup` + `cn deps lock` + `cn deps restore`, asserts each pinned package's installed version matches the lockfile pin. Confirms the chain α just shipped is honest end-to-end.

## Cycle findings

Voice per CDD §1.4 β step 9: factual observations and patterns only. Triage is γ's at PRA time.

### CDD §9.1 trigger fired — "loaded skill failed to prevent a finding"

α's PR body declared `eng/go` (with §2.17 Parse vs Read) as an active skill. The §2.17 rule is:

> Pure functions operate on bytes/data. IO functions operate on paths/network. Keep them separate. … `Parse*` takes `[]byte`, returns typed data. `Read*` takes a path, calls `Parse*`. No mixing.

In the same diff, α introduced the pure parser `pkg.ParseInstalledManifestData` and the IO wrapper `restore.ReadInstalledManifest` (correctly applied to `restore.restoreOne`) — and also added a parallel anonymous-struct unmarshal of the same `cn.package.json` `version` field inside `doctor.checkPackages`. The new pure parser was visible to α at authoring time but `doctor` did not consume it. β caught the duplication in round 1 (F1, severity C) via review §2.2.8 (authority-surface conflict) — α fixed cleanly in round 2 by importing `pkg` into `doctor.go` and replacing the inline unmarshal with `pkg.ParseInstalledManifestData(manifestData)`.

Pattern: §2.17 covers the abstract failure mode (parse + read mixed) and is well-followed at the call-site level. The specific failure mode that surfaced here — "do not reintroduce a parallel parser for the same fact when a new pure parser is added in the same PR" — is one step removed from §2.17's literal prose. Surfaces affected: `doctor.checkPackages` (round-1 state), `pkg.ParseInstalledManifestData` (under-consumed), `restore.ReadInstalledManifest` (correct exemplar that doctor did not mirror).

### Harness branch-naming drift vs CDD §4.2

α's branch was `claude/alpha-tier-3-skills-IZOsO`. CDD §4.2 canonical format is `{agent}/{issue}-{scope}` (or `{agent}/{version}-{issue}-{scope}` when version is known) — the issue number is part of the canonical branch identity. α's PR body acknowledged "(per dispatch)" — the harness assigned the name; α did not choose. Surfaces affected: CDD §4.2 (canonical-branch contract) ↔ harness/dispatch convention.

### β polling pattern — first-iteration absorption masks pre-existing state

CDD §Tracking gives reference shell loops for transition-only emission via `Monitor`-wrapped `git fetch` + `comm -13`. The transition-only contract is correct (avoids context flood from emitting on every poll iteration), but the first iteration sets `prev` to the empty string and absorbs whatever the current state is — so anything that already existed when the Monitor started is silently invisible to the polling channel.

Concrete shape this took on #274: β's first synchronous PR-list query at session start returned `[]`. β then armed the narrow Monitor (branch glob `*230*`); α opened PR #274 between session start and the broader Monitor's start time. When the broader Monitor started, the broad PR-set transition guard absorbed `#274` as part of the first-iteration baseline, and the narrow branch glob never matched α's harness-named branch (`claude/alpha-tier-3-skills-IZOsO`, no `230` substring). β did not see PR #274 from polling at all — only from a synchronous re-query the operator prompted.

The β-side adjustment that resolved the operational gap was: at every β intake startup, do a synchronous `mcp__github__list_pull_requests state=open` BEFORE relying on transition-only polling. The polling channel owns the future; the synchronous channel owns the past. The CDD §Tracking reference scripts do not name this baseline step.

Surfaces affected: CDD §Tracking (reference shell loops + wake-up contract); β step 3 polling spec.

### Branch-glob narrowness vs harness scope-encoding

The CDD §Tracking example glob is `'origin/claude/*-<N>-*'` — assumes the harness encodes `<N>` (issue number) somewhere in the branch path. The harness used in this cycle encodes scope words and a random suffix instead (`claude/alpha-tier-3-skills-IZOsO`); the issue number is not in the branch name. β's first-cut Monitor used the narrow glob and matched no branches. β subsequently broadened to "any new origin branch" with downstream synchronous filtering, which caught the next transition. Surfaces affected: CDD §Tracking glob template ↔ harness branch-naming convention.

### Pre-existing stale doc-path references on `main`

While walking review §2.2.5 (cross-reference validation) on the diff, β grepped `INVARIANTS.md` and `POLYGLOT-PACKAGES-AND-PROVIDERS.md` references in `eng/go/SKILL.md` (T-002, T-004) and `pkg.go:134`. Both reference paths return zero matches; the actual files live at `docs/alpha/architecture/INVARIANTS.md` and `docs/alpha/agent-runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md`. The references are pre-existing on `main` (verified via `git diff origin/main..HEAD` showing zero changes to these reference strings in this PR's diff), so they were correctly held out of scope for #230 review. Surfaces affected: `eng/go/SKILL.md` §2.17/§2.18 + `src/go/internal/pkg/pkg.go:134` ↔ actual doc paths.

### Mechanical ratio

Findings: 4 total — 1 mechanical (F4: dead `BIN_VERSION_OUT` capture), 3 judgment (F1/F2/F3). Mechanical ratio = 25% by count; below the §9.1 absolute threshold of "≥ 10 findings" so the trigger does not fire on ratio alone. Recorded for cycle-economics observation.

### Round economics

2 review rounds — meets the §9.1 default target of ≤ 2. Round 2 narrowing was scoped strictly to F1–F4 anchors (3 files, 50 / 17 lines) with no scope creep; verification re-walk took ≤ one synchronous diff inspection.

## Invariant audit

Project design constraints checked per `cdd/review/SKILL.md` §2.2.13 against `docs/alpha/DESIGN-CONSTRAINTS.md`:

| Constraint | Pre-cycle state | Post-cycle state | Direction |
|---|---|---|---|
| §1 Source of truth (each fact in one place) | Lockfile was the integrity authority; installed `cn.package.json` was nominally an installed-version authority but ignored. After F1 round-2 fix, `pkg.ParseInstalledManifestData` is the single parser. | One parser per fact, two consumers (`restore.ReadInstalledManifest` and `doctor.checkPackages`). Lockfile remains integrity authority; installed manifest is now an explicit, mechanically-consulted installed-version authority. | tightened |
| §2.1 One package substrate (`cn.package.v1`) | `pkg.PackageManifest` and `pkgbuild.PackageManifest` are separate types covering different shapes; only the former gained `Version`. JSON wire format unchanged. | Same. No new manifest types introduced; `Version` is purely additive on the install-time decode shape. | preserved |
| §2.2 Source / artifact / installed clarity | `src/` (authored) → `dist/packages/` (built, no longer in git per 3.58.0) → `.cn/vendor/packages/` (installed). Installed-version was opaque to operator at runtime. | Installed-version is now visible (doctor stale list) and enforced (restore reinstall on drift). Boundary unchanged; clarity tightened on the installed face. | tightened |
| §3.1 Git-style subcommands | Not touched. `cn doctor`, `cn deps lock`, `cn deps restore` retain noun-verb shape. | Same. | preserved |
| §3.2 Dispatch boundary (cli/ owns dispatch only) | Touched: domain logic added to `internal/restore/`, `internal/doctor/`, `internal/pkg/`. No `cmd_*.go` modified. | Same — all new logic lives in domain packages; `cli/` thin-wrapper invariant preserved. | preserved |
| §4.1 Surface separation (skills / commands / orchestrators / providers) | Not touched. | Same. | preserved |
| §4.2 Registry normalization (one runtime descriptor) | Not touched. | Same. | preserved |
| §5.0 OCaml deprecated, no edits | No OCaml files modified by this cycle. | Same. | preserved |
| §6.1 Reason to change (one per boundary) | `restore` owns "decide whether to install"; `doctor` owns "report installed truth"; `pkgbuild` owns "build artifacts". | Same — boundaries reinforced (no new convenience buckets; F1 round-2 fix specifically resisted creating a parallel "doctor parses manifest" responsibility by reusing the canonical parser). | preserved |
| §6.2 Policy above detail (kernel owns policy) | Not touched. | Same. | preserved |
| §6.3 Degraded-path visibility | Pre-cycle: silent skip on stale install was the canonical degraded-path-invisibility example named in the issue. | Post-cycle: every reinstall reason emits `slog.WarnContext` with structured key/value attrs (`installed_version`, `lockfile_version`, `expected_version`, `error`); doctor surfaces the same drift to the operator with named diagnostics. The canonical inspectability surface for this failure class is now testable and reviewable. | tightened |

Cross-reference to `docs/alpha/architecture/INVARIANTS.md` (the active invariants doc):

| Invariant | Status |
|---|---|
| T-002 (kernel remains minimal and trusted; cli/ owns dispatch only) | preserved — verified via diff: zero `cmd_*.go` changes; all new logic in domain packages |
| T-004 (source/artifact/installed explicit) | preserved-tightened — installed `cn.package.json` is now an explicit authority for installed-version, mechanically consulted by both write-path (restore) and read-path (doctor) |

Architecture check (`cdd/review/SKILL.md` §2.2.14) result on the merged diff: 6 of 7 axes `yes`, 1 `n/a` (Registry normalization — not a registry-touching change). No `no` answers; the §2.3.5 architecture-gate was not blocked.

## Tag-push deferral (CDD §1.4 β-step 8)

`git push origin 3.59.0` returned `HTTP 403` from the sandbox-side git proxy (`http://127.0.0.1:43365/git/usurobor/cnos`). `git push origin main` from the same session succeeded — the 403 is specific to tag refs from this β session, not a general remote-write block. Per CDD §1.4 β-step 8:

> If tag push fails due to env constraints (e.g. sandbox HTTP 403), commit all release artifacts to main and defer tag push to γ/operator — do not block closure on it.

Status:

- All release artifacts are on `origin/main` at `9dd30d9d`: `VERSION=3.59.0`, three `cn.package.json` files (3.59.0), `cn.json` (3.59.0), `CHANGELOG.md` ledger row, `RELEASE.md`, `.cdd/releases/3.59.0/beta/CLOSE-OUT.md`.
- The tag `3.59.0` exists locally at `9dd30d9d` but is NOT on origin (`git ls-remote --tags origin 3.59.0` empty).
- Deferred to γ/operator (CDD γ phase 2 step 6): `git push origin 3.59.0` from a session with non-sandboxed write access. Once the tag lands on origin, `release.yml` triggers automatically (4-platform matrix build → GitHub release with `RELEASE.md` body → `release-smoke.yml` fires on `release: published`).
- β does not block closure on this step.

## β verdict trail

- Round 1: REQUEST CHANGES — https://github.com/usurobor/cnos/pull/274#issuecomment-4321807808
- Round 2: APPROVED — https://github.com/usurobor/cnos/pull/274#issuecomment-4321845284
- Merge commit: `9980e3f5ff6e7221c8f29463443d526ba8a5cb6d` (squash of α PR #274)
- Release commit: `9dd30d9dd10da965dd01fad72bc04aa6b5046cf8` (`release: 3.59.0 — distribution chain honesty (#230, PR #274)`)
- Tag (local only, pending γ push): `3.59.0` → `9dd30d9d`
- Issue closed: #230 `state: closed`, `state_reason: completed`
