---
cell_kind: implementation
---

# γ scaffold — cnos#656: repo lifecycle Phase 1 (ledger + `cn repo status`)

## Source of truth

- Issue: cnos#656 (sub of cnos#655, the `cn repo` lifecycle wave).
- Design surface: `docs/development/design/cn-repo-lifecycle.md` (§Phase-contract
  precision P1–P3 is the normative AC surface for this cell; the earlier
  "core decision" section's illustrative JSON predates the tightening and is
  superseded where they conflict — see Decision 1 below).
- Mock-first target (new, authored this cycle before implementation, mirroring
  `cn-repo-install-MOCKS.md`'s discipline per A6):
  `docs/development/design/cn-repo-status-MOCKS.md`.

## run_class

`first_pass`. No prior rejection, no scanner "MECHANICAL reversion" comment,
no operator continuation-authorization comment, no pre-existing `cycle/656`
branch/PR/artifacts at claim time. Steps B–E of the repair re-entry preflight
do not apply.

## Scope for this cell (Phase 1 only)

In scope: `cn.repo.state.v1` schema + CI validation gate (P1); ledger write
wired into `cn repo install` (idempotent, deterministic) with workflow render-
contract records (P2); `cn repo status` (+`--json`, `--check`), read-only (P3);
A3 best-effort reconstruction when the ledger is absent.

Out of scope (later phases per cnos#655): `cn repo update`, `cn repo repair`,
`cn repo prune`, `cn repo uninstall`, the A7 two-step install-semantics flip,
`cn repo doctor` (A5 only requires *delineating* status from the existing
global `cn status`/`cn doctor` in help text — it does not require building
`cn repo doctor` in this phase; the design doc lists `doctor` under Phase 1's
command surface table but the issue's own Acceptance section only requires
`status`, so `cn repo doctor` is deferred to a later phase and this scaffold
records that as an explicit scope call, not an oversight).

## Decision 1 — `managed_files` requires `id`/`sha256` uniformly

The design doc's illustrative ledger JSON (§"The core decision") omits
`id`/`sha256` for the `manifest`/`lockfile` `managed_files` entries. But the
issue's own P1 Acceptance text is more specific and is the authoritative AC
for this cell: "`managed_files`/workflow record missing required fields" and
"managed_files records require path/kind/id/sha256". Read literally, this
requires `id`+`sha256` on **every** `managed_files` entry, not just workflow
records. This cell follows the issue text over the design doc's earlier
illustration: every `managed_files` record gets a content `sha256` (computed
at ledger-write time) and a stable `id`. This is a strict improvement over the
illustration, not a regression — a `sha256` on `deps.json`/`deps.lock.json`
too is exactly what Phase 4's `repair`/local-edit detection (A4) will need for
non-workflow managed files as well, so paying that cost now avoids a schema
bump later.

## Decision 2 — schema location: `schemas/repo_state.cue`, not `schemas/cdd/`

`schemas/cdd/` is explicitly scoped ("Machine-checkable CUE schemas for the
**generic coherence-cell kernel**" — `schemas/cdd/README.md`) to the CDD
contract/receipt/boundary-decision triad. `.cn/repo.state.json` is an
installer-lifecycle artifact, not a coherence-cell receipt — putting it under
`schemas/cdd/` would conflate two unrelated schema families sharing a
directory only by accident of the design doc's "e.g." phrasing (which itself
says "or the repo's chosen location"). The repo's existing precedent for a
standalone, non-cell schema is `schemas/skill.cue` (top-level, own fixtures
dir, own validation script) — this cell follows that precedent:
`schemas/repo_state.cue` + `schemas/fixtures/repo-state/{valid,invalid}/` +
`scripts/ci/validate-repo-state.sh`.

## Decision 3 — closed schema (not open like `#Receipt`)

`schemas/cdd/receipt.cue`'s `#Receipt` and `schemas/skill.cue`'s `#Skill` are
deliberately **open** (`...`) because they type externally-authored,
forward-compatible formats (a receipt closure record; SKILL.md frontmatter a
human writes, where LANGUAGE-SPEC mandates unknown-key tolerance). `.cn/repo.
state.json` is different in kind: it is a **machine-written ledger**, entirely
owned by `cn repo install`/`update`/`repair`, never hand-authored. Making
`#RepoState` (and its nested `#ManagedFile`/`#ManagedDir`/`#Source`
definitions) **closed** (the CUE default for a `#Definition` with no trailing
`...`) is what makes the P1 invariants structural rather than script-enforced:

- **timestamp-free**: no timestamp-shaped field is declared anywhere in the
  schema, so a fixture that adds `timestamp`/`created_at`/`installed_at` fails
  `cue vet` as an unknown-field error under closedness — no separate
  denylist/blocklist logic needed.
- **no lock duplication (A1)**: there is no top-level `packages` field
  declared, so a fixture that adds one (with `version`/`sha256`) fails the
  same way. Package identity is referenced by `name` string only inside
  `managed_dirs[].package`, joined against `.cn/deps.lock.json` (`cn.lock.v2`)
  at read time by the Go code — never copied into the ledger.

This mirrors the same "keep the proof in CUE, not in runtime logic" doctrine
`schemas/cdd/README.md` §"Architectural choice" reasons 1–3 (structural proof
over required-keys-in-open-map) already establishes for this repo — applied
here to closedness instead of required-key unification, same underlying
principle.

## Decision 4 — Go package layout: `internal/repostate`

Pure types + (de)serialization + validation-adjacent helpers (no `os`/`net`
per `internal/pkg`'s own discipline, mirrored here) live in a new
`internal/repostate` package: `RepoState`, `Source`, `ManagedFile`,
`ManagedDir`, `ExternalExpectations` Go structs with `json` tags mirroring the
CUE field names exactly. `internal/repoinstall` gets a new
`writeRepoState`-shaped step wired into `applyInstall` (after the existing
`writeManifest`/lock/restore/gitignore steps) that constructs a `RepoState`
from the just-computed `manifest`/lockfile/dispatch outcome and writes
`.cn/repo.state.json` via the same `json.MarshalIndent(v, "", "  ")` + trailing
newline convention `writeManifest` already uses (`repoinstall.go:580-590`) —
deterministic byte-for-byte output on a same-input rerun (no timestamps, no
map-iteration-order-dependent fields — `managed_files`/`managed_dirs` are
sorted by `path` before serialization).

`cn repo status` gets its own package `internal/repostatus` (mirrors
`internal/hubstatus`'s "no `cli/` import" discipline) plus a thin
`cli.RepoStatusCmd` wrapper registered as `Spec().Name = "repo-status"`
(dispatches from `cn repo status` via the existing noun-verb joiner, same
mechanism as `cmd_repo_install.go`'s `"repo-install"`). `NeedsHub: false`;
resolves `RepoRoot` via the same `gitRepoRoot(ctx)` helper `RepoInstallCmd`
uses (not `inv.HubPath` — A5's repo-scope-vs-hub-scope distinction).

## Decision 5 — drift classification (A2) scope for Phase 1

Full "renderer-moved" classification (re-invoking `cn-install-wake` fresh and
diffing) is implementable now because the renderer script + its exact flags
are already wired in `runDispatchCds`. `cn repo status` re-renders the
recorded workflow tier in a temp dir and compares the fresh sha256 against
both the live file and the ledger sha to produce the three-way classification
(`matches_ledger` / `user_edit` / `renderer_moved`). When the ledger is
absent (A3, pre-3.82.x installs), there is no ledger sha to anchor the
three-way split, so status degrades to a two-way `matches_fresh_render` /
`unknown` classification for that case only (recorded explicitly in Mock D2
of the status mocks doc) — this is a genuine information loss inherent to a
missing ledger, not a cut corner.

## Non-goals this cell explicitly does not attempt

- `cn repo install`'s existing re-resolve-latest-on-rerun behavior (A7) is
  UNCHANGED by this cell — Phase 1 only adds the ledger write as an
  additional side effect of the existing install flow; it does not touch
  install's version-resolution semantics (that is Phase 2's job).
- Backfill triggers only on `cn repo install` (the only lifecycle-write
  command that exists yet). The design doc's "backfill on next
  install/update/repair" is forward-looking; `update`/`repair` do not exist
  in this repo yet, so this cell cannot wire backfill into commands that
  don't exist — install-triggered backfill satisfies A3 for what exists today,
  and later phases inherit the same `internal/repostate` reconstruction helper
  this cell ships.

## Acceptance mapping (this cell's own checklist)

- P1 → `schemas/repo_state.cue` + `schemas/fixtures/repo-state/{valid,invalid}/`
  + `scripts/ci/validate-repo-state.sh` + CI wiring + Go-side
  `internal/repostate` tests.
- P2 → workflow `managed_files` record round-trip fields (Decision 4);
  `cn repo status`'s drift classification (Decision 5).
- P3 → `cn repo status` never writes; verified in tests via
  `git status --porcelain`-equivalent (temp-dir mtime/content diff check).
- A6 → `docs/development/design/cn-repo-status-MOCKS.md` authored before
  implementation; `mock_parity` block required in `gamma-closeout.md`.
