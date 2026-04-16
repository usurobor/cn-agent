# RELEASE.md

## Outcome

Coherence delta: C_Σ **A** (`α A`, `β A-`, `γ A-`) · **Level:** **L6**

Git-style subcommand dispatch (`cn <noun> <verb>`) is now the primary command surface, with flat-hyphenated form preserved as backward-compat fallback. The command surface is now extensible by convention rather than by registry mutation.

## Why it matters

cnos commands were accumulating as flat hyphenated forms (`cn kata-run`, `cn kata-list`, `cn deps-restore`). Each new command was a top-level name competing in a flat namespace. Git-style grouping (`cn kata run`, `cn kata list`, `cn deps restore`) creates natural command families discoverable via `cn help <noun>`. This is infrastructure for the package command system — packages can register verbs under their noun without namespace collision.

## Added

- **Git-style subcommand dispatch** (#254, PR #257): `cn <noun> <verb>` resolves to `<noun>-<verb>` internally. `cn help <noun>` lists all verbs in the group. Flat forms preserved as fallback. 14 dispatch tests + 2 help tests.
- **DESIGN-CONSTRAINTS.md §3.1** (1a7c362): Convention documented before implementation — design converged prior to α cycle.

## Changed

- **main.go dispatch** uses `ResolveCommand` from new `internal/cli/dispatch.go`. Lookup order: noun-verb first, flat fallback second.
- **`cn help <noun>`** now shows grouped verb listing instead of falling through to general help.

## Validation

- CI: 7/7 green (go, kata-tier1, kata-tier2, I1, I2, 2× notify)
- `cn kata run` resolves to `kata-run` command
- `cn help kata` lists all kata verbs
- `cn deps restore` (flat form) still works — backward compat preserved
- Deployment deferred (no binary-breaking change; existing flat commands unaffected)

## Known Issues

- #253 — ContentClasses divergence (8 vs 5) still open
- Squash-merge destroys `.cdd/` artifacts committed on PR branches — α close-out had to be recovered to main manually (process gap, not code)
