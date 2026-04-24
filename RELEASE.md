# RELEASE.md

## Outcome

Coherence delta: C_Σ A- (`α A-`, `β A-`, `γ A-`) · **Level:** `L7`

Packaging pipeline is now deterministic, content-model-driven, and dist-free. The rebase-race class that blocked two review rounds in the prior cycle is structurally eliminated. CDD triadic coordination is significantly tightened — agents now act on steps instead of asking permission.

## Why it matters

Every packaging PR since 3.57.0 hit the same race: α rebases, rebuilds dist, pushes; main moves; dist is stale. This release removes committed tarballs entirely, makes `cn build` derive its packlist from the content-class model, and normalizes tar/gzip output for cross-build determinism. CDD skill patches close a class of agent coordination failures where system prompts overrode skill instructions.

## Fixed

- **Deterministic package tarballs** (#264): zeroed tar header timestamps, uid/gid, uname/gname, and gzip header time/OS byte. Same source tree now produces byte-identical tarballs regardless of checkout time, machine, or user.
- **Cross-Go-toolchain gzip non-determinism**: bumped go.mod and all CI workflows from Go 1.22 to Go 1.24. `compress/gzip` flate output is not stable across minor Go versions.
- **CDD agent compliance gap**: agents treated skill instructions as suggestions when system prompts conflicted. Six patches make instructions imperative and explicitly override environment defaults (PR creation, PR subscribe, reply to RC comments, β subscribe+wait, section-by-section writing, dispatch prompt trimming).
- **CDD β polling gap**: β had "wait for PR" with no mechanism. Now has a concrete `gh pr list` polling loop.
- **CDD event subscription unreliability**: GitHub notifications are unreliable with shared identity. All roles now use periodic `gh` CLI polling instead of event subscriptions.

## Added

- **Content-class packlist derivation** (#262 partial, PR #265): `cn build` now derives its packlist from `cn.package.json` + recognized content-class directories, not the whole package root. Stray files no longer ship in tarballs. Shared `pkg.ContentClasses` constant used by both build and activation.
- **CDD package audit** (`docs/gamma/cdd/CDD-PACKAGE-AUDIT.md`): 34-finding structural review of the CDD skill program. Filed as #268 for convergence work.
- **CTB §8.5.1**: four skills-language design evidence items from CDD practice (global aspects, authority hierarchy, visibility enforcement, dependency graph).
- **CDD §1.4 large-file authoring rule**: generic cross-role constraint — any file >50 lines written section by section. Replaces 4 per-step duplicates.
- **CDD close-out voice constraint**: α/β close-outs report factual observations only; triage is γ's job.
- **eng/go §2.13**: cross-toolchain determinism axis documented.

## Changed

- **dist/ removed from git** (#266, PR #267): tarballs are no longer committed. I3 (coherence CI) rebuilds from source and compares checksums against committed `checksums.txt`. Eliminates the rebase-race class that hit twice in the #262 cycle.
- **CDD dispatch model**: γ produces both α and β prompts at dispatch time. β starts intake immediately — no need to wait for α's PR before subscribing.
- **CDD γ skill**: reads last PRA first (binding selection input), subscribes to issue before dispatching.

## Removed

- `dist/packages/*.tar.gz` — committed tarballs eliminated from git history going forward.
- `docs/gamma/essays/SKILLS-LANGUAGE-EVIDENCE.md` — consolidated into CTB §8.5.1.
- Per-step section-by-section writing instructions (4 locations) — replaced by generic §1.4 rule.

## Validation

- `go test ./...` green (all packages)
- `go test -race ./internal/activation/... ./internal/doctor/...` green
- `scripts/check-version-consistency.sh` passes
- `cn build` produces deterministic, byte-identical tarballs across runs
- I3 coherence check rebuilt and verified

## Known Issues

- #262 stays open — AC1 (cn build → cn pack rename), AC5 (--list/--dry-run), AC3 root-file narrowing, .gitignore honoring deferred
- #268 — CDD package convergence (34 audit findings) dispatched, not yet landed
- `cn publish` command (#266 design addendum) — not yet implemented; currently manual release-asset upload
