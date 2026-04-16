# RELEASE.md

## Outcome

Coherence delta: C_Σ **A** (`α A`, `β A`, `γ A`) · **Level:** **L5**

`cn deps lock` now means what its name says. Pre-3.56.2 the lockfile was a dump of the entire package index; post-3.56.2 it is exactly what `.cn/deps.json` pins, resolved against the index, with explicit errors for missing manifests or unresolved pins. The Go runtime now matches the OCaml `lockfile_for_manifest` behavior that was already correct, and the kata test infrastructure that had been quietly riding on the bug is now gated against regression.

## Why it matters

The operator's single source of truth for "what does this hub use?" (`.cn/deps.json`) had no effect on the lockfile. A hub pinning `cnos.core@3.54.0` got a lockfile with every `(name, version)` pair the index knew about, and `cn deps restore` installed all of them — last-writer-wins on `restoreOne` resolved the resulting duplicate-name entries to whichever `cnos.core@*` happened to be last in iteration order. The pin was inert.

Two downstream harnesses had silently grown to depend on the bug. The Tier-2 CI test hub used object-syntax `{"cnos.core": "3.54.0"}` which the parser rejects as `[]ManifestDep`; the Tier-1 `06-install.sh` accepted any installed package count ≥ 1. Both now pin the real version via `pkg_version_from_source` (derived from `src/packages/*/cn.package.json`) and assert exact-match installed sets. The over-vendoring regression class is CI-gated going forward.

## Fixed

- **`cn deps lock` honors `.cn/deps.json` pins** (#250, PR #259): the lockfile contains exactly the packages declared in the manifest, resolved against the index. Unresolved pins error explicitly. Missing `deps.json` errors explicitly. Go/OCaml semantics now agree.

## Added

- **`pkg.ParseManifest`** (`src/go/internal/pkg/pkg.go`): pure deps.json parser mirroring the existing `ParseLockfile` / `ParsePackageIndex` shape.
- **`restore.ReadManifest`** (`src/go/internal/restore/restore.go`): IO wrapper preserving the `eng/go §2.17` Parse/Read purity boundary.
- **Six AC-named tests** in `restore_test.go` — the reproducer from #250, AC2 duplicate-free, AC3 restore-only-pinned, AC4 revert-behavior, plus missing-manifest and missing-pin error paths.
- **`pkg_version_from_source`** helper in `src/packages/cnos.kata/lib.sh` — reads `src/packages/<pkg>/cn.package.json` and echoes the version, replacing hardcoded `3.54.0` strings in R3/R4 katas.

## Changed

- **`restore.GenerateLockFromIndex`**: reads manifest first, iterates `m.Packages`, collects unresolved pins into one explicit error listing every missing `name@version`. Godoc updated to spell out the contract and error cases.
- **CI `kata-tier2` setup** (`.github/workflows/ci.yml`): derives real package versions at job time, emits the canonical array-schema `deps.json`.
- **Tier-1 kata `scripts/kata/06-install.sh`**: overwrites the post-`cn setup` `deps.json` with an explicit `cnos.core` pin read from source, and tightens the post-condition from "≥ 1 installed" to "exactly `cnos.core` installed".
- **`cnos.kata` `lib.sh` `write_deps_json`**: emits `packages` as an array of `{name, version}` objects (matches the parser; the pre-#250 object-syntax worked only because the lockfile dumped the full index regardless).

## Removed

- The "iterate the entire index" code path inside `GenerateLockFromIndex`. There was no policy under which producing a lockfile of the whole index was correct — it was the bug.

## Validation

- **Build**: `go build ./... && go vet ./...` — clean on `4f860f1`.
- **Unit + integration**: `go test ./...` — all packages green.
- **Race**: `go test -race ./internal/restore/` — green.
- **Targeted**: `go test ./internal/restore -run TestGenerateLockFromIndex -v` — 6/6 pass.
- **CI on PR head** (`44d431a`): 7/7 checks green (go, kata-tier1, kata-tier2, Package/source drift I1, Protocol contract schema sync I2, notify × 2).
- **Module-truth audit** (review §2.2.9): the only `idx.Packages` iteration sites are `pkgbuild/build.go` (index *construction*, correct) and the two `idx.Lookup(name, version)` calls in `restore.go` (both lookup-by-pin). No sibling "dump the index" incoherence remains.
- **Go/OCaml parity**: `src/ocaml/lib/cn_package.ml` `parse_manifest_dep` and `src/ocaml/cmd/cn_deps.ml` `lockfile_for_manifest` already expected an array of `{name, version}` and resolved per pin; the Go path now matches. `cn setup` in `hubsetup.go` writes the same array shape.

The targeted coherence delta — "`.cn/deps.json` actually controls `.cn/deps.lock.json`" — is proved by six AC-named tests (all six fail on revert), tightened CI post-conditions on Tier-1 and Tier-2 kata suites, and a green 7/7 check matrix on the PR head.

## Known Issues

- `cn setup` still pins `deps.json` to the **binary's** version. In dev/CI environments without `-ldflags`, `cn deps lock` immediately after `cn setup` will now error explicitly rather than silently succeeding with an index dump. Tier-1 kata `06-install.sh` works around this by overwriting the default with an explicit pin read from source. A follow-up issue (make `cn setup` resolve against the local index) is a candidate if this bites operators; today it affects only dev checkouts.
