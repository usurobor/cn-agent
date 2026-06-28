# Design: Prove the Distribution Pipeline End-to-End

**Issue:** #227
**Mode:** MCA
**Active Skills:** eng/go, cdd/design
**Engineering Level:** L5 — local correctness, no boundary change

## Problem

The artifact-first distribution pipeline (`src/packages/ → cn build → dist/packages/ → cn deps restore → .cn/vendor/packages/`) has never been proven end-to-end. Both the build and restore modules were written to the same spec (BUILD-AND-DIST.md) but have never exchanged real data. Three specific incoherences exist:

1. **URL format mismatch:** `pkgbuild.UpdateIndex` writes relative URLs (`cnos.core-1.0.0.tar.gz`). `restore.downloadFile` constructs an HTTP request — a relative URL is not a valid HTTP URL. No local file path support exists.
2. **No lockfile generation:** `restore.Restore` reads `.cn/deps.lock.json` but nothing in the system creates it from build output. The pipeline has no entry point for restore.
3. **VendorPath design divergence:** `pkg.VendorPath` produces `<name>@<version>/` but BUILD-AND-DIST.md §2 says the installed path is `<name>/` (no version). Quote: "the installed path does not (`.cn/vendor/packages/cnos.core/`). This keeps the active runtime simple."

## Constraints

- Go sole implementation language (T-003)
- One package substrate (INV-001) — no second restore path
- Source → artifact → installed explicit (T-004) — this cycle proves the chain
- Parse/Read purity boundary (eng/go §2.17, INVARIANTS T-004): pure logic in `internal/pkg/`, IO wrappers in `internal/restore/` and `internal/pkgbuild/`
- cli/ dispatch boundary (INVARIANTS T-002): domain logic stays in packages, not in CLI wiring

## Impact Graph

### Downstream consumers
- `restore.Restore` reads the index URL field → must handle relative paths
- `restore.restoreOne` calls `downloadFile` → needs local file copy path
- `restore.restoreOne` calls `pkg.VendorPath` → path format must match design
- All existing restore tests use `@version` vendor paths → must update
- `internal/doctor/` (hub doctor) reads vendor paths → must stay consistent

### Upstream producers
- `pkgbuild.BuildOne` produces tarballs
- `pkgbuild.UpdateIndex` writes index with relative URLs
- Source manifests (`src/packages/*/cn.package.json`) define name + version

### Copies and embeddings
- `pkg.VendorPath` is the single source of truth for vendor path format
- `pkg.TmpTarballPath` uses a separate format (no `@`, uses `-`)
- BUILD-AND-DIST.md defines the canonical path format

## Proposal

### 1. Add local file restore to `restore.go`

When the index URL has no scheme (`://`), resolve it relative to the index file directory and use file copy instead of HTTP download. This is the preferred dev-workflow approach per the issue.

Pass `indexDir` (directory containing index.json) through the restore call chain. In `restoreOne`, resolve relative URLs against `indexDir`.

### 2. Add lockfile generation to `pkgbuild`

Add `GenerateLockfile(repoRoot string, results []BuildResult) error` that writes `.cn/deps.lock.json` from build results. This is build-adjacent — the build knows the name, version, and SHA-256 of each package.

Wire into `cn build` as an optional `--lockfile <path>` flag, with default output to `.cn/deps.lock.json` when a `.cn/` directory exists in the hub.

### 3. Fix VendorPath to match BUILD-AND-DIST.md

Change `pkg.VendorPath` from `<name>@<version>` to `<name>` — matching the design doc's explicit decision that "the installed path does not" carry the version.

### 4. Add `cn deps lock` subcommand

Add a `lock` subcommand to `cn deps` that reads `dist/packages/index.json` and generates `.cn/deps.lock.json` for all packages in the index. This decouples lockfile generation from the build step and follows the Unix philosophy of composable commands.

## Non-goals

- Remote package index / gh-pages hosting (#181)
- Package command discovery + dispatch (#226)
- Runtime contract changes
- CI release workflow integration
- Changing the tarball internal path format (already correct — no top-level wrapper)

## File Changes

### Edit
| File | Change |
|------|--------|
| `src/go/internal/pkg/pkg.go` | Fix `VendorPath` to `<name>/` (no `@version`). Add `GenerateLockfileData` pure function. |
| `src/go/internal/pkg/pkg_test.go` | Update `TestVendorPath` expected value. Add lockfile generation test. |
| `src/go/internal/restore/restore.go` | Add `indexDir` parameter threading. Add `copyFile` for local paths. Route relative URLs to local copy. |
| `src/go/internal/restore/restore_test.go` | Update vendor path expectations. Add local-file restore test. Add build→restore round-trip test. |
| `src/go/internal/pkgbuild/build.go` | Add `GenerateLockfile` IO wrapper. |
| `src/go/internal/pkgbuild/build_test.go` | Add lockfile generation test. |
| `src/go/internal/cli/cmd_deps.go` | Add `lock` subcommand. |

## Acceptance Criteria

- [x] AC1: `cn build` produces `dist/packages/cnos.{core,cdd,eng}-*.tar.gz`
- [x] AC2: `cn build` produces `dist/packages/index.json` with `cn.package-index.v1` schema
- [x] AC3: `cn build --check` validates all 3 packages
- [x] AC4: `cn deps restore` installs all 3 packages into `.cn/vendor/packages/<name>/`
- [x] AC5: Round-trip verified: tarball SHA-256 matches index, extracted `cn.package.json` parses
- [x] AC6: Tests updated/added for the end-to-end path

## Plan

### Step 1: Fix VendorPath (P0)
- Change `pkg.VendorPath` to `<hubPath>/.cn/vendor/packages/<name>/`
- Update `pkg_test.go` TestVendorPath
- Update all test expectations in `restore_test.go`
- **AC:** `VendorPath("/hub", "cnos.core", "1.0.0")` → `/hub/.cn/vendor/packages/cnos.core/`
- **Depends on:** nothing
- **Unblocks:** Steps 2, 3

### Step 2: Add local file restore (P0)
- Add `copyFile(src, dest string) error` to `restore.go`
- Modify `Restore` to accept `indexPath` (already does) and derive `indexDir`
- In `restoreOne`, check if URL contains `://` — if not, resolve relative to indexDir and copy
- **AC:** `Restore(ctx, hub, indexPath)` works with relative URLs in index
- **Depends on:** Step 1 (vendor path)
- **Unblocks:** Step 4

### Step 3: Add lockfile generation (P0)
- Add `GenerateLockfileData(results []BuildResult) ([]byte, error)` to `pkgbuild` (pure)
- Add `GenerateLockfile(hubPath string, results []BuildResult) error` (IO wrapper)
- Add `lock` subcommand to `cmd_deps.go` that reads index → writes lockfile
- **AC:** `cn deps lock` reads `dist/packages/index.json` → writes `.cn/deps.lock.json`
- **Depends on:** nothing
- **Unblocks:** Step 4

### Step 4: End-to-end test + verification (P0)
- Add `TestBuildRestoreRoundTrip` in `restore_test.go` or new `e2e_test.go`
- Build test packages → generate lockfile → local restore → verify content
- **AC:** round-trip test passes
- **Depends on:** Steps 1–3

## Known Debt

- `cn deps lock` generates lockfile for ALL packages in index — no filtering by profile/manifest yet
- No `cn deps restore --local` flag — local vs HTTP is auto-detected from URL format
- Index URL format diverges between local (relative) and remote (absolute) — release workflow handles absolutization

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 0 Observe | CHANGELOG TSC, lag table, POST-RELEASE v3.50.0, INVARIANTS | — | MCI frozen, 6 growing lag items, Phase 3 complete |
| 1 Select | Issue #227 | — | P1 distribution pipeline proof — prerequisite for #226, #186 |
| 4 Gap | this artifact | — | Build and restore modules never exchanged real data: 3 mismatches |
| 5 Mode | this artifact | eng/go, cdd/design | MCA, L5 local correctness, Go implementation |
| 6 Artifacts | this artifact | — | design + plan combined (issue provides thorough spec) |
| 7 Tests | TestBuildRestoreRoundTrip + 4 more | eng/go | round-trip pipeline proved in automated test |
| 8 Code | pkg.go, restore.go, build.go, cmd_deps.go, doctor.go | eng/go | 3 mismatches fixed, pipeline proven end-to-end |
| 9 Self-coherence | SELF-COHERENCE-227.md | cdd | all 6 ACs met, invariants preserved |
