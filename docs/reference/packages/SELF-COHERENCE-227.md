# Self-Coherence Report — #227 Distribution Pipeline

**Issue:** #227
**Branch:** `claude/alpha-cdd-cycle-227-vCB7e`
**Mode:** MCA
**Active Skills:** eng/go, cdd/design

## Artifact Integrity

- [x] Design artifact present (DESIGN-227-distribution-pipeline.md)
- [x] CDD Trace in design artifact
- [x] All 6 ACs addressed
- [x] Tests written before / alongside code (round-trip test proves the pipeline)
- [x] Self-coherence artifact present (this file)

## AC Verification

| # | AC | Status | Evidence |
|---|---|---|---|
| AC1 | `cn build` produces 3 tarballs | met | `cn build` → cnos.{core,cdd,eng}-*.tar.gz |
| AC2 | `cn build` produces valid index.json | met | `dist/packages/index.json` with `cn.package-index.v1` schema |
| AC3 | `cn build --check` validates all 3 | met | `cn build --check` → "All packages valid" |
| AC4 | `cn deps restore` installs all 3 | met | `cn deps restore` → "Restored 3 package(s)" to `.cn/vendor/packages/<name>/` |
| AC5 | Round-trip verified: SHA-256 matches, manifest parses | met | sha256sum matches index entry; `cn.package.json` parses correctly |
| AC6 | Tests updated/added | met | `TestBuildRestoreRoundTrip`, `TestRestoreLocalFile`, `TestGenerateLockfileData`, `TestGenerateLockfile`, `TestVendorPathVersionIgnored` |

## Changes Summary

### 1. VendorPath fixed (pkg.go)
- `VendorPath` changed from `<name>@<version>` to `<name>/` per BUILD-AND-DIST.md
- Doctor updated: version drift detection now uses lockfile comparison instead of `@version` parsing
- All downstream consumers (restore, doctor, tests) updated

### 2. Local file restore (restore.go)
- `fetchTarball` dispatches: remote URLs → HTTP, relative paths → local copy
- `copyFile` for local file operations
- `indexDir` threaded through restore chain for relative URL resolution
- Existing HTTP tests continue to pass unchanged

### 3. Lockfile generation (pkgbuild/build.go, cli/cmd_deps.go)
- `GenerateLockfileData` (pure) + `GenerateLockfile` (IO) in pkgbuild
- `cn deps lock` subcommand: reads index → writes `.cn/deps.lock.json`

### 4. Tests
- `TestBuildRestoreRoundTrip`: full pipeline proof (build → index → lockfile → restore → verify)
- `TestRestoreLocalFile`: relative URL resolution
- `TestGenerateLockfileData` + `TestGenerateLockfile`: lockfile generation
- `TestVendorPathVersionIgnored`: VendorPath version independence
- `TestRunAllPackageMissing`: updated doctor test for lockfile-based drift
- All 9 test packages pass (66 total tests)

## Invariants Check

| Constraint | Touched? | Status |
|---|---|---|
| INV-001 (one package substrate) | Yes — proved end-to-end | preserved |
| INV-002 (language-neutral metadata) | No | preserved |
| INV-003 (distinct surfaces) | No | preserved |
| INV-004 (kernel owns policy) | No | preserved |
| INV-005 (protocol above transport) | No | preserved |
| INV-006 (hub state is not source) | No | preserved |
| T-001 (kernel is package-compatible) | No | preserved |
| T-002 (cli/ dispatch boundary) | Yes — `cn deps lock` domain logic in pkgbuild, CLI is thin | preserved |
| T-003 (Go sole language) | Yes — all changes in Go | preserved |
| T-004 (source/artifact/installed explicit) | Yes — pipeline proved end-to-end | tightened (VendorPath now matches BUILD-AND-DIST.md) |
| T-005 (content classes explicit) | No | preserved |

## Known Debt

- `cn deps lock` generates lockfile for ALL packages in index — no profile/manifest filtering yet
- No `cn deps restore --local` flag — local vs HTTP auto-detected from URL scheme
- Hub doctor lockfile-based drift check is simpler than the old `@version` approach but loses per-package version visibility (would need to read each installed `cn.package.json` for full version info)
