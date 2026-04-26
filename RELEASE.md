# RELEASE.md

## Outcome

Coherence delta: C_Σ A (`α A`, `β A`, `γ A-`) · **Level:** `L6` (cycle cap: L5)

`cn build --check` now catches what it used to silently let through — packages can no longer ship with a manifest declaring `commands.<X>.entrypoint` pointing at a missing or non-regular file, and a `skills/orphan/` subtree with no `SKILL.md` anywhere is now a build-time failure instead of a silent runtime omission. The build-time validator and the runtime activation walk now read from the same filesystem-as-authority surface, so what `cn build --check` accepts is exactly what the runtime can later activate.

## Why it matters

Before this release the manifest-claim ↔ filesystem-fact link was unenforced at build time. A `commands.daily.entrypoint = "commands/daily/cn-daily"` declaration with the file missing produced a tarball that installed cleanly but failed at dispatch. A `skills/orphan-skill/` directory with only stray markdown shipped as dead weight that activation silently dropped. Both classes were latent install-time / run-time degradations that the source repo had no way to surface. They surface now, at the same gate where everything else structural already surfaces (the I1 `coherence-build-check` job).

This is the same authority-surface tightening pattern as #261 (`activation` reads filesystem for skill discovery) and #230 (`restore` + `doctor` read the installed manifest version). Build-time validation, runtime activation, and runtime dispatch now agree on what constitutes a valid command/skill.

## Added

- **`pkgbuild.checkCommandEntrypoints`** (#235, PR #276) — for every declared command, verifies the entrypoint resolves to an existing regular file under the package root. Rejects missing files, non-regular targets (directory at the entrypoint path), path-traversal (`../`), and empty entrypoints. Reads `cn.package.json` via the canonical `pkg.ParseFullManifestData` parser — no parallel parser introduced (eng/go §2.17 "one parser per fact").
- **`pkgbuild.checkSkillDirectories` + `pkgbuild.containsSkillMd`** (#235, PR #276) — for every top-level subdirectory of `skills/`, demands at least one `SKILL.md` somewhere in the subtree. Filesystem-as-authority surface (DESIGN-CONSTRAINTS §1, #261), same as activation's `discoverPackageSkills`. Top-level scope is intentionally narrow: namespace containers (e.g. `cnos.eng/skills/eng/`) and resource subdirectories (e.g. `cnos.core/skills/naturalize/references/`) are exempt by construction.
- **9 tests** in `build_test.go` covering pass + fail + exempt cases for both new validation rules.

## Changed

- **`pkgbuild.CheckOne`** now extends beyond structural content-class presence to enforce the two new rules above.

## Validation

- 7/7 CI checks green on `f7d27b4` (PR #276 head): `go`, `Package/source drift (I1)`, `Protocol contract schema sync (I2)`, `kata-tier1`, `kata-tier2`, `notify`×2.
- AC4 verified by α via deliberate-break dry-run on a copy of the repo: injecting a fake `commands.fake-cmd.entrypoint` and an `skills/orphan-skill/` with only a stray `notes.md` into `cnos.cdd`, then running the freshly-built `cn`, surfaced both new error classes with exit code 1. The chain "issues → CLI exit non-zero → I1 job fails" is unchanged from the existing structural-check flow.
- β independently verified the diff against its own contract surface and the surrounding `pkgbuild`/`activation`/`cli` modules at HEAD `f7d27b4`; review APPROVED with no findings (posted as comment per shared-identity rule, `review/SKILL.md` §7.1).

## Known issues

- **N1 — pre-existing parallel manifest parsers in `pkgbuild`.** `pkgbuild.PackageManifest` + `ParseManifestData` (minimal) coexist with `pkg.PackageManifest` + `pkg.ParseInstalledManifestData` (also minimal). Eliminating the duplication changes `DiscoveredPackage.Manifest`'s static type and ripples through `BuildOne`/`UpdateIndex`/`UpdateChecksums`. Cleanly scope-defended in PR #276 §Known debt; latent refactor for a follow-up.
- **N2 — `BuildOne` does not gate on `CheckOne`.** `cn build` (without `--check`) still tarballs without invoking the validator. CI runs `--check` separately, so the gate exists in CI; the local safety is opt-in. Pre-existing structure, not introduced by this release.
- **N3 — Symlink entrypoints.** `os.Stat` follows symlinks; an entrypoint that is a symlink to a regular file outside the package root passes `IsRegular()`. Could be tightened later via `os.Lstat` + symlink-target containment if needed.

## β process finding (deferred to γ for `CDD.md §Tracking` patch)

The MCP `head=` filter and `gh pr list --search 'closes:#N'` qualifier are unreliable for in-repo branches with slashes. β's first synchronous PR-existence check after the new-branch transition event used the head filter, got `[]`, and trusted it — PR #276 was already open with that exact head ref. Operator caught it within one polling tick. Recommended `§Tracking` patch: on any new-branch transition event, do a synchronous broad open-PR list (no head/search filter) and scan client-side for `(head.ref == new-branch) OR (body matches (?i)\b(closes\|fixes\|resolves\|refs)\s*#N\b)`. Same shape as the existing baseline rule, extended from session-start to per-event handling.
