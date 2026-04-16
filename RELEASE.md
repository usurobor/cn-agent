# RELEASE.md

## Outcome

Coherence delta: C_Σ **A** (`α A-`, `β A`, `γ A`) · **Level:** **L6**

Content-class authority converged to one source of truth. Two surfaces (`cn build --check` and `cn status`) that disagreed on what counts as a package content class now derive from the same canonical list and the same filesystem predicate. Doctrine (`PACKAGE-SYSTEM.md`) matches the code.

## Why it matters

`pkgbuild.ContentClasses` (filesystem discovery, 8 entries) and `pkg.FullPackageManifest.ContentClasses()` (JSON-field heuristic, 5 entries including a non-existent `providers`) gave different answers about the same package. Operators had no single truthful answer to "what does this package contain?" and every new content class risked compounding the drift silently — the 3.55.0 cycle already shipped `katas` in `pkgbuild` without updating the manifest heuristic or the canonical doctrine. This release collapses the two surfaces into one list, one predicate, one authority (filesystem presence, per §3), and syncs `PACKAGE-SYSTEM.md §1.1–§4.3` so the doctrine matches what the code actually does. Future content-class additions now touch exactly one list.

## Fixed

- **ContentClasses divergence** (#253, PR #258): `cn build --check` (8 classes) and `cn status` (5 classes) agreed on content-class membership for the first time.

## Added

- **`pkg.ContentClasses`** (`src/go/internal/pkg/pkg.go`): canonical exported list, 8 entries in §1.1 order: doctrine, mindsets, skills, extensions, templates, commands, orchestrators, katas. Pure (no IO).
- **`pkgbuild.FindContentClasses(pkgDir)`** (`src/go/internal/pkgbuild/build.go`): shared filesystem predicate iterating `pkg.ContentClasses`. Single authority for both `cn build --check` and `cn status`.
- **4 new tests** exercising empty, full, subset, and end-to-end (`cn status` with all 8 classes) paths.

## Changed

- **`pkgbuild.CheckOne`** walks `pkgtypes.ContentClasses` instead of its former local duplicate list.
- **`hubstatus.Run`** now calls `pkgbuild.FindContentClasses(pkgDir)` instead of the JSON-field heuristic on the manifest. `cn status` output order for multi-class packages now follows the canonical §1.1 order (user-visible convergence, not regression).
- **`docs/alpha/package-system/PACKAGE-SYSTEM.md`**: §1.1 gains `katas` row as the 8th content class (completes 3.55.0 drift); §1.2 "Directory tree copy" enumeration includes katas; §3 names `pkg.ContentClasses` and `pkgbuild.FindContentClasses` as the implementation; §4.1/§4.3 class counts corrected 7→8 and 6→8; §6 history gains v3.55.0 and v3.56.1 rows.

## Removed

- **`pkgbuild.ContentClasses`** (local duplicate list) — callers now import `pkgtypes.ContentClasses`.
- **`(m *FullPackageManifest) ContentClasses() []string`** method and the `Skills`, `Orchestrators`, `Extensions`, `Providers` fields it depended on, plus the now-orphaned `SkillsJSON` type. The JSON-field heuristic was the wrong authority (PACKAGE-SYSTEM.md §3: content classes are discovered by directory presence, not manifest JSON fields). `providers` was never a content class — it's a runtime capability surface per POLYGLOT-PACKAGES-AND-PROVIDERS.md, orthogonal to the content-class model.

## Validation

- CI (merged commit `aacb817`): 7/7 green — `go`, `kata-tier1`, `kata-tier2`, `Package/source drift (I1)`, `Protocol contract schema sync (I2)`, `notify` ×2.
- `go build ./...`, `go vet ./...`, `go test ./...` — all green, including 4 new tests.
- `go run ./cmd/cn build --check` against real `src/packages/` — all 5 packages valid, including `cnos.cdd.kata` (katas-only).
- Authority-surface audit: single canonical list (`grep -rn "pkgbuild.ContentClasses" src/go` → 0 matches), no JSON-field heuristic (`grep -rn "\.ContentClasses()" src/go` → 0 matches), no orphaned readers of removed fields.
- Deployment deferred (no binary-breaking change; existing packages unaffected — the removed JSON fields were only consumed by the now-removed method).

## Known Issues

- **OCaml `cn_build.ml`** still carries a 7-field `source_decl` shape without `katas`; it is no longer the live `cn build` path (Go kernel owns it), but the OCaml build surface should either be dropped or mechanically synced in a future cycle.
- **`PACKAGE-AUTHORING.md`** uses the phrase "katas as commands" in one spot — consistent with `katas` being a content class (a package can do either — `cnos.kata` ships runner commands, `cnos.cdd.kata` ships katas as content) but not cross-linked to §1.1. Cosmetic.
- **`scripts/stamp-versions.sh`** only bumps `cnos.core`, `cnos.cdd`, `cnos.eng` — `cnos.kata` and `cnos.cdd.kata` use independent versioning (established in 3.55.0), which is intentional but not documented next to the script.
