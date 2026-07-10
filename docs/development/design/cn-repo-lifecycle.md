# Architecture note: `cn repo` lifecycle commands

**Status:** Proposed (design surface — pre-implementation).

CNOS has a working repo installer (`cn repo install`, cnos#607/#608, shipped
in 3.82.x). It needs a coherent **lifecycle model** before install, update,
repair, prune, and uninstall grow ad hoc.

## The core decision

- **`cn repo install` is convergence** — a desired-state *sync*, not an
  implicit updater forever. Boring, rerunnable, safe.
- **`cn repo update` is movement** — the explicit command that advances the
  desired state.
- **`cn repo uninstall` is conservative removal** — driven by an ownership
  **ledger**, removing only CNOS-managed surfaces, never user matter.

> Do not make `cn repo install` the command that silently does everything.
> That is the path to surprising users and unsafe deletion.

## Five first-class states

| State | Artifact |
|---|---|
| Desired | `.cn/deps.json` |
| Exact resolved | `.cn/deps.lock.json` (schema `cn.lock.v2`) |
| Materialized | `.cn/vendor/packages/<name>/` |
| **Managed-surface (new)** | **`.cn/repo.state.json`** — the lifecycle ledger |
| External side-effect | GitHub labels, workflows, secrets, App installation, issue-FSM labels |

The new addition — **`.cn/repo.state.json`** — is what makes uninstall safe.
It is **committed, deterministic, timestamp-free**, and records the
managed-surface ownership CNOS needs to remove only what it installed:

```json
{
  "schema": "cn.repo.state.v1",
  "profile": "cds",
  "source": { "channel": "stable", "release": "3.82.6",
              "index": "github:usurobor/cnos/releases/3.82.6/index.json" },
  "managed_files": [
    { "path": ".cn/deps.json", "kind": "manifest" },
    { "path": ".cn/deps.lock.json", "kind": "lockfile" },
    { "path": ".github/workflows/cnos-cds-dispatch.yml", "kind": "workflow",
      "id": "cnos-cds-dispatch", "tier": "engine", "sha256": "<content-sha>" }
  ],
  "managed_dirs": [
    { "path": ".cn/vendor/packages/cnos.core", "package": "cnos.core" }
  ],
  "external_expectations": {
    "labels": { "mode": "ensure", "source": "cnos.core/labels.json",
                "delete_on_uninstall": false }
  }
}
```

## Command surface (keep it small)

```
cn repo install     # bootstrap or sync (convergence)
cn repo install --locked   # deterministic CI/frozen sync
cn repo install --clean    # clean exact reinstall from lock
cn repo update      # move to latest selected release/channel (movement)
cn repo status      # installed state + drift (truth)
cn repo doctor      # diagnose lifecycle problems (+ --fix → safe subcommands)
cn repo repair      # rehydrate damaged materialized state (no version change)
cn repo prune       # remove extraneous materialized state
cn repo uninstall   # remove CNOS-managed surfaces, conservatively
```

- **Install** is convergence · **Update** is movement · **Repair** is
  restoration · **Prune** is cleanup · **Uninstall** is conservative managed
  removal · **Status/doctor** is truth.

### Semantics (key points)

- **`cn repo install`**: if no lock → bootstrap (resolve latest stable, write
  deps/lock, restore, write ledger). If lock exists → **sync the vendor tree
  to the lock; do not float to latest** unless an explicit `--release` /
  `--index` / `update` is given. Flags: `--release`, `--index`, `--locked`
  (never mutate deps/lock; fail if changes would be needed), `--clean` (remove
  + recreate vendor from lock, spirit of `npm ci`), `--dispatch cds`,
  `--engine`.
- **`cn repo update`**: resolve selected channel/release/index → update deps +
  lock + changed packages + ledger. Does **not** prune unless `--prune`. Does
  not remove dispatch unless asked. `--dry-run` shows exact package changes.
- **`cn repo status`**: source, deps vs lock vs materialized versions, dispatch
  present/tier, canonical labels (ok/missing/drifted/unknown), orphan vendored
  packages, local edits to managed files, available update. `--json`,
  `--check` (nonzero exit on drift).
- **`cn repo doctor`**: diagnoses (reads ledger + deps/lock/vendor + workflow
  markers + labels if a token is available); `--fix` calls safe subcommands.
- **`cn repo repair`**: rehydrate materialized state **without changing
  versions** — verify/reinstall locked packages, restore missing dirs, replace
  damaged same-version dirs, regenerate a managed workflow only if its tier is
  recorded. Does not update to latest; does not delete orphans.
- **`cn repo prune`**: remove vendored dirs absent from the lock + stale
  temp/cache; never workflows, labels, or unknown `.cn` content. `--yes`
  required when destructive noninteractively.
- **`cn repo uninstall`**: plans from `.cn/repo.state.json` (not heuristics);
  removes base + (opt) dispatch surfaces; **labels never removed by default**
  (`--labels --yes`); refuses locally-edited generated files unless `--force`;
  `--dry-run` / `--json` / `--yes`. Tiers: base, dispatch, labels.

### Idempotence contract

| Command | Rerun (same inputs) | Changes versions? | Removes things? |
|---|---|---|---|
| `cn repo install` | no diff | no, if lock exists | no |
| `cn repo install --locked` | no diff or fail | no | no |
| `cn repo install --clean` | recreates vendor | no | replaces vendor only |
| `cn repo update` | no diff if current | yes | no, unless `--prune` |
| `cn repo repair` | no diff if healthy | no | replaces damaged managed artifacts |
| `cn repo prune` | no diff if no orphans | no | yes, only extraneous vendor/cache |
| `cn repo uninstall` | fails/empty after removed | n/a | yes, only managed surfaces |

## Safety rules

1. **Never delete unknown user matter** — `.cn/commands/`, `.cn/issues/`,
   `.cn/cells/`, unknown vendored packages, user workflows, non-canonical
   labels, secrets — unless the command explicitly names that surface and the
   user confirms.
2. **Labels are external side effects** — install may *ensure*; doctor may
   *detect drift*; uninstall must **not** delete by default (`--labels --yes`).
3. **Workflows require ownership markers** — remove a generated workflow only
   if a managed marker exists **and** the ledger records it **and** the content
   hash matches; otherwise refuse unless `--force`. Marker:
   ```
   # Managed by CNOS.
   # cn-managed-file: cnos-cds-dispatch
   # cn-managed-schema: cn.repo.managed_file.v1
   ```
4. **Update is not uninstall** — a package leaving the desired set is not
   deleted until `cn repo update --prune` or `cn repo prune`.
5. **Binary update is separate** — `cn repo update` updates repo-installed
   *packages*, never the `cn` binary. Binary update stays `install.sh` (or a
   future `cn self update`). Do not conflate repo package lifecycle with CLI
   binary lifecycle.

## Implementation phases

1. **State ledger + status** — `.cn/repo.state.json`, `cn repo status`
   (+`--json`); install writes a deterministic ledger; rerun no-diff.
2. **Install-semantics hardening** — bootstrap-if-absent, sync-from-lock-if-
   present, no float unless explicit; add `--locked`, `--clean`.
3. **Update** — `cn repo update` (+`--dry-run`, no prune unless `--prune`).
4. **Repair + prune** — content-aware repair; orphan prune; both `--dry-run`;
   prune `--yes` when destructive noninteractively.
5. **Uninstall** — ledger-driven; managed-surface-only; local-edit detection;
   labels preserved unless `--labels`; unknown `.cn` preserved; `--dry-run` /
   `--json` / `--yes`.
6. **Docs** — lifecycle section in `INSTALL-CDS.md` (install / rerun / update /
   repair / prune / uninstall / what-is-not-removed / labels / recover).

---

## Coherence review — amendments to fold into implementation

*(Grounded against current code — `internal/repoinstall`, `internal/restore`,
`cmd/cn/main.go` — as of 3.82.6. These sharpen the note; they do not change its
decisions.)*

**A1 — The ledger must not duplicate the lock (DRY / anti-drift).** The lock
(`cn.lock.v2`) already records `packages[]` with `name`, `version`, `sha256`.
`repo.state.json` should **not** re-carry package versions/shas — three copies
of the same fact (deps → lock → state) is a three-way drift hazard. The ledger
owns only what the lock cannot: `source`/`channel` provenance, `managed_files`
(workflows + markers + sha), `managed_dirs` (path→package mapping for
uninstall), and `external_expectations` (labels). Package identity is
*referenced from the lock*, not copied. (The example above already leans this
way — make it normative: no `version`/`sha256` under a `packages` key in the
ledger.)

**A2 — Distinguish "user edited it" drift from "renderer moved" drift.** The
workflow `sha256` in the ledger is an *as-installed* fingerprint. When CNOS
itself re-renders (a new `cn` version changes `cn-install-wake` output), the
live workflow's sha will differ from the ledger — that is **not** a local edit.
`status`/`doctor` must classify workflow drift three ways: (a) matches ledger →
clean; (b) differs from ledger **and** differs from a fresh re-render → *user
edit* (uninstall refuses without `--force`); (c) differs from ledger **but
equals a fresh re-render of the recorded tier** → *renderer moved* (`repair` /
`update` re-renders and updates the ledger sha). Without this, every `cn`
version bump reports spurious workflow drift. (Same coupling class as the
`install-wake-golden.yml` byte-identity gate.)

**A3 — Backfill path for pre-ledger installs.** Repos installed by 3.82.x have
**no `.cn/repo.state.json`**. `status`/`repair`/`uninstall` must handle its
absence gracefully: reconstruct a best-effort ledger from `deps.lock.json` +
the vendored tree + a workflow-marker scan, and have `cn repo install`/`update`
**backfill** the ledger idempotently on next run. Uninstall on a repo with no
ledger and no markers should be conservative (refuse workflow removal without
`--force`), never heuristic-delete.

**A4 — `cn repo repair` closes a real, verified gap.** `restore.restoreOne`
skips a package when the vendored `cn.package.json` **version** matches the lock
(`"package already installed, skipping"`) — it does **not** hash the vendored
*content*. So a same-version local edit inside `.cn/vendor/packages/<name>/`
is invisible to `restore`. `repair` must either (a) verify each vendored dir's
content against the locked `sha256`, or (b) delete + reinstall locked packages
unconditionally. Recommend a `--verify` (hash-check, report) vs default
(delete+reinstall) split so repair is both a diagnosis and a fix.

**A5 — `cn repo status`/`doctor` vs the existing global `cn status`/`cn
doctor`.** CNOS already registers global `StatusCmd` and `DoctorCmd`
(binary/hub health — git identity, PATH, etc.). The repo-scoped
`cn repo status`/`cn repo doctor` are *different scope* (repo-install
desired-state drift). Delineate explicitly in help text and docs so they read
as complementary siblings, not collisions: global = "is my `cn`/hub healthy?",
repo = "is this repo's CNOS install in sync?". Consider having global
`cn doctor` cross-reference `cn repo doctor` when run inside an installed repo.

**A6 — Apply the wave's own mock-parity discipline.** This lifecycle wave is
exactly the kind of human-visible surface (`status` output, `uninstall
--dry-run` plan, `doctor` diagnosis) that cnos#607 mocked *before*
implementation with a receipt-parity gate
(`cn-repo-install-MOCKS.md`). Each phase's closeout should carry `mock_parity`
rows for its human-visible outputs, `missed: 0`. Mock the `status`,
`uninstall --dry-run`, and `doctor` outputs first.

**A7 — Ordering is right; make Phase 2 a two-step behavior flip.** Phase 1
(ledger + status) correctly precedes any removal command — you cannot safely
uninstall/prune without the ledger. Phase 2 changes a **shipped, proven**
command (`cn repo install` currently re-resolves latest on every rerun, verified
at `repoinstall.go:285`+`:549`). Flip it in two steps to avoid surprising
existing users: **vNext** — when a lock exists and a plain `install` *would*
float, sync-to-lock but **print a warning** recommending `cn repo update`;
**vNext+1** — make lock-preserving the silent default. Gate the flip on the
`status`/ledger work landing first.

## Phase-contract precision (pre-dispatch review — normative)

*These tighten the phase contracts so a worker cell cannot implement unsafe
lifecycle behavior by accident. They are requirements, not suggestions.*

**P1 — Durable ledger schema + validation gate (Phase 1).** `cn.repo.state.v1`
must ship as a **checked-in schema** (e.g. `schemas/cdd/repo_state.cue` or the
repo's chosen location) with a CI validation gate + tests asserting: a valid
`repo.state.json` passes; **timestamp fields are rejected** (deterministic,
timestamp-free); **package `version`/`sha256` duplication is rejected** (A1 —
those live in the lock); `managed_files` records require `path`/`kind`/`id`/
`sha256`; workflow records carry enough render contract to reproduce drift
classification (P2 below). The ledger is the basis for uninstall, so it needs
schema enforcement from Phase 1.

**P2 — Workflow records store the render contract, not just the tier (Phase
1).** `path`/`kind`/`id`/`tier`/`sha256` is **not enough** to fresh-render the
same workflow, so A2's drift classification (matches-ledger / user-edit /
renderer-moved) is not implementable without the non-secret render inputs.
A dispatch-workflow record must carry, e.g.:
```json
{ "kind": "workflow", "id": "cnos-cds-dispatch",
  "path": ".github/workflows/cnos-cds-dispatch.yml", "tier": "engine",
  "renderer": "cnos.core/install-wake", "renderer_package": "cnos.core",
  "renderer_version_source": "lock",
  "agent": "cnos", "workflow_pat_secret": "CNOS_WORKFLOW_PAT",
  "bot_name": "cnos-bot", "bot_id": "12345678", "sha256": "..." }
```
**Store secret *names* and render identity only — never secret values.**

**P3 — `cn repo status` is read-only by default (Phase 1).** Status is
truth/reporting: it **never writes `.cn/repo.state.json`** by default. It may
emit a *reconstructed ledger candidate* under `--json` (A3 backfill source),
but only `install`/`update`/`repair` — on an explicit backfill path — persist
the ledger.

**P4 — Existing-lock install is no-network / no-latest by default (Phase 2).**
Do not make the "would float" decision depend on discovering whether latest
changed (that reintroduces the network dependency convergence is meant to
remove). Contract: if a lock exists and no explicit `--release`/`--index`/
`update` intent is present, `cn repo install` **syncs from the lock without
resolving latest** and prints an informational line —
`using existing .cn/deps.lock.json; run cn repo update to move to a newer release`.

**P5 — `--locked` = no resolution / no lock mutation, but fetch is allowed
(Phase 2).** `--locked` must not perform dependency resolution or mutate
`deps.json`/`deps.lock.json`; it **may** fetch the package tarball URLs the
existing lock references. "No network at all" is a *future* `--offline` flag,
not implied by `--locked`.

**P6 — `cn repo repair` v1 = delete+reinstall (Phase 4).** The lock `sha256`
is the **package tarball** digest, not a canonical expanded-tree digest, so
"hash the vendored dir against the lock sha" (as loosely stated) is invalid.
v1 default: **delete + reinstall** the locked package dirs. A `--verify` mode
requires either (Option A) Phase 1 storing a post-extract **tree digest** in
the ledger that `--verify` checks, or (Option B) `--verify` checks only
manifest presence/version and reports `content verification unavailable` when
no tree digest exists. Do not compare an expanded directory against a tarball
hash.

## Comparators (why this shape)

- **npm / pnpm** — lockfile-driven reproducibility; `install` honors the lock;
  `ci`/frozen-lock is the strict path; `prune` removes *extraneous* materialized
  packages (present but not in dependency state) — not desired deps.
- **Cargo** — the cleanest command taxonomy: manifest mutation (`add`/`remove`)
  vs lock update (`update`) vs reproducibility guards (`--locked`/`--frozen`/
  `--offline`).
- **Homebrew** — maintenance vocabulary: `doctor` (diagnose, nonzero on
  problems), `autoremove`, `cleanup` (dry-run-safe).

CNOS installs **repo capabilities, not just libraries** (vendored packages +
workflows + labels + eventual App state), so the correct pattern is
**lockfile discipline (npm/pnpm/Cargo) + status/doctor/cleanup (Homebrew) + a
CNOS managed-surface ledger** for workflows and labels.

## Related

- [`cn-repo-install-MOCKS.md`](cn-repo-install-MOCKS.md) — the installer
  design surface + receipt-parity contract this wave should mirror.
- [`cnos-installer-github-app.md`](cnos-installer-github-app.md) — the App
  surface whose per-tier permission model the ledger's `external_expectations`
  (labels) aligns with.
- [`../../guides/INSTALL-CDS.md`](../../guides/INSTALL-CDS.md) — the human
  install guide (Phase 6 extends it with the lifecycle section).
