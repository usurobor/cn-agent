# RELEASE.md

## Outcome

Coherence delta: C_Σ A- (`α A-`, `β A`, `γ A-`) · **Level:** `L6` (diff: `L7`; cycle cap: `L6`)

Distribution chain is now structurally honest. The lockfile → installed `cn.package.json` → released-binary chain previously lied silently in two places: a lockfile version bump was ignored if the version-less vendor dir already existed, and there was no automated check that a freshly released binary could bootstrap a fresh hub. Both gaps are closed: `cn deps restore` now reads its own authority on every run, `cn doctor` surfaces the same drift visibly, and a release-bootstrap smoke runs against production endpoints on every published release.

## Why it matters

The `packages/index.json` migration leak (v3.x era) was found by accident, not by a named test. The silent v1→v2 skip was the same shape — a stale install masquerading as a current one. After this release, both failure classes are caught mechanically: `restoreOne` removes the stale tree before reinstalling, `doctor.checkPackages` fails with `(installed X, locked Y)` so the operator can see the drift at runtime, and `release-smoke.yml` blocks any release whose binary cannot self-bootstrap.

## Fixed

- **Silent v1→v2 skip in `cn deps restore`** (#230, AC1+AC2): `restore.restoreOne` now reads the installed `cn.package.json` version and reinstalls on drift. The reinstall path `os.RemoveAll`s the prior tree first, so v1-only files cannot leak into the v2 install. The same-version fast-path is preserved (one extra `os.ReadFile` of a small file).
- **Doctor agreed with the lie** (#230, AC6): `doctor.checkPackages` now reads the same installed manifest via the shared pure parser and reports `StatusFail` with `(installed X, locked Y)`, `(no manifest)`, or `(unparseable manifest)` so the runtime surface no longer claims "all present" when the lockfile and the vendor disagree.

## Added

- **`scripts/smoke/90-release-bootstrap.sh`** (#230, AC3): exercises the full bootstrap chain from a fresh empty directory — downloads `cn-<platform>` + `index.json` + every referenced tarball from the GitHub release URL, verifies each tarball's SHA-256 against the `index.json` entry (the same authority `cn deps restore` itself trusts), then `cn init` + `cn setup` + `cn deps lock` + `cn deps restore` + per-package version-match check. Exit codes 0 (ok) / 1 (chain broken) / 2 (skipped — offline / no release).
- **`.github/workflows/release-smoke.yml`** (#230, AC4+AC5): runs the smoke on `release: types: [published]` (4-platform matrix: linux-x64, linux-arm64, macos-x64, macos-arm64) plus `workflow_dispatch`. rc=2 → `::warning::` (graceful offline skip); rc=1 → `::error::` (release flagged).
- **`pkg.ParseInstalledManifestData`** + **`restore.ReadInstalledManifest`** (#230): pure parser + IO wrapper for installed `cn.package.json` (name + version), shared by `restore.restoreOne` and `doctor.checkPackages`. `ValidatePackageManifestData` reuses the same parser — one parser per fact (eng/go §2.17 + design §3.2). `pkg.PackageManifest` gained an additive `Version` field; JSON wire format unchanged.
- **Doctor `(unparseable manifest)` diagnostic + test** (#230, round-2 F3): completes the three-way stale-state coverage in `doctor.checkPackages`.

## Changed

- **`release.yml` build step stamps the version** (#230, supports AC3): `go build -ldflags "-X main.version=${{ github.ref_name }}"` so the released binary's compiled-in `version` equals the release tag. Without this, `cn setup` writes deps.json with `version="dev"` and `cn deps lock` cannot resolve packages in the released index — i.e. the smoke would correctly fail. (Regular CI keeps `dev`; kata 06 already overrides deps.json with src-derived versions.)

## Removed

- (no deletions in this release.)

## Validation

- `go test ./...` green; `go test -race ./internal/restore/... ./internal/doctor/... ./internal/pkg/...` green.
- `scripts/check-version-consistency.sh` passes against `VERSION=3.59.0`.
- CI 7/7 success at the merge-base SHA `93ea1d6` — `go`, `kata-tier1`, `kata-tier2`, `Package/source drift (I1)`, `Protocol contract schema sync (I2)`, `notify` ×2.
- Smoke verified locally to exit 2 (`RESULT: skipped (offline)`) when offline; full pass-path validated against the released `3.59.0` binary by the `release-smoke.yml` matrix once the release is published.

## Known Issues

- **`pkg.ContentClasses` documentation reference** (`pkg.go:134`) still points to `POLYGLOT-PACKAGES-AND-PROVIDERS.md` and `eng/go/SKILL.md` §2.17/§2.18 still point to `INVARIANTS.md` — the actual files live at `docs/alpha/agent-runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md` and `docs/alpha/architecture/INVARIANTS.md`. Pre-existing on `main` from before this cycle; out of #230 scope. Surfaced as a γ follow-up.
- **Branch-naming harness drift**: PR #274's branch was `claude/alpha-tier-3-skills-IZOsO`, non-canonical per CDD §4.2 (`{agent}/{issue}-{scope}`). Harness/dispatch convention should be reconciled with §4.2 — γ-side process iteration.
- **β polling baseline gap**: CDD §Tracking's reference Monitor scripts use transition-only emission, which silently absorbs pre-existing state on first iteration. β's #274 round 1 missed PR creation as a result. Synchronous-baseline check before transition-only polling is the missing template step — γ-side process iteration for the next CDD pass.
