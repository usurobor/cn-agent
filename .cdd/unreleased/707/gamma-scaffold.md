## Issue / Mode

**Issue:** [cnos#707](https://github.com/usurobor/cnos/issues/707) — "cli: add `cn repo uninstall` — symmetric teardown of everything install creates (manifest-driven)".

**Mode:** substantial. New Go domain package + new `cn` subcommand + tests + an
additive write-path change in an existing installed package
(`internal/repoinstall`). No wave-mode (`§5.2`) applies — this is a single
cell, not a multi-cell wave.

Full issue body was fetched fresh via `gh issue view 707 --json body -q .body`
at scaffold time (2026-08-06) — this scaffold does not paraphrase from any
stale summary. AC1–AC6 numbers below are the issue's own numbering, verbatim.

---

## IMPORTANT — pre-existing overlapping wave discovered (§2.2a peer enumeration)

Before framing surfaces, γ ran the required peer enumeration (`gamma/SKILL.md
§2.2a`): `rg -il "uninstall|install-manifest|installmanifest" src/go` (no
hits — confirms the issue's own claim that no uninstall code exists yet is
empirically true) — but a broader `gh issue list` sweep for "repo lifecycle"
surfaced a **directly overlapping, currently-open wave** that #707's issue
body does not reference:

- **[cnos#655](https://github.com/usurobor/cnos/issues/655)** (OPEN, master
  tracker) — "wave: cn repo lifecycle commands (ledger → sync-install →
  update → repair/prune → uninstall)". Design doc:
  `docs/development/design/cn-repo-lifecycle.md`.
- **[cnos#656](https://github.com/usurobor/cnos/issues/656)** (OPEN,
  `status:review`, PR #663 open/unmerged) — "repo lifecycle Phase 1 — state
  ledger (`.cn/repo.state.json`, schema `cn.repo.state.v1`) + `cn repo
  status`". This is install writing a managed-surface ledger — the *same
  shape of obligation* as #707's AC2 ("install writes a manifest recording
  its side effects").
- **[cnos#660](https://github.com/usurobor/cnos/issues/660)** (OPEN) — "repo
  lifecycle Phase 5 — `cn repo uninstall` (conservative, ledger-driven)" —
  **the same command #707 asks for**, designed against the #656 ledger.

**Disposition (γ, not a silent decision — recorded for δ/operator
visibility):** #707 is independently claimed (`CLAIM-REQUEST.yml`,
`cds-dispatch`, `todo → in-progress`) and dispatch-ready with its own
numbered ACs naming a specific, narrower artifact
(`.cn/install-manifest.json`, not `.cn/repo.state.json`). γ is **not**
reopening the selection decision mid-scaffold. α proceeds against #707's
literal design (§ below) — but this scaffold pins the wire schema name
distinctly (`cn.install-manifest.v1`, not `cn.repo.state.v1`) specifically
*so it does not collide with or get confused for* #656/#660's eventual
ledger. γ's close-out triage for #707 MUST carry a disposition row
reconciling #655/#656/#660 against what #707 actually ships (e.g. "#656/#660
superseded by #707", "#656/#660 to be redesigned to consume #707's
manifest", or "#707's manifest to be renamed/merged into the #656 ledger
schema in a follow-on cycle") — this is flagged here so it is not lost.
**α/β do not need to resolve this; γ does, at close-out.**

---

## Surfaces γ expects α to touch

All paths below were read in full or grepped concretely — not guessed.

1. **`src/go/internal/repoinstall/repoinstall.go`** (922 lines). AC2 requires
   `cn repo install` to write an install manifest. The natural insertion
   point is `applyInstall` (lines 537–574), which already writes
   `.cn/deps.json` via `writeManifest` (lines 580–590) after
   `restore.GenerateLockFromIndex` / `restore.Restore` / label-doctor
   (`ensureCanonicalDispatchLabels`, lines 516–531) have run — the manifest
   write must happen **last**, after every side effect it records has
   actually occurred (files, gitignore line via
   `hubsetup.EnsureGitignoreEntry`, vendor path, dispatch workflow path
   `dispatchWorkflowPath` line 392–394, and — when `--dispatch cds` ran —
   the labels `ensureCanonicalDispatchLabels` applied). `Result` (lines
   234–240) and `Run` (lines 243–338) also need to carry/report the new
   manifest.
2. **`src/go/internal/cli/cmd_repo_install.go`** (172 lines) — the existing
   thin `cli/`-boundary wrapper (eng/go §2.18) around `repoinstall`:
   flag parsing → `gitRepoRoot(ctx)` (reused from `cmd_cell.go:17`, **not**
   `inv.HubPath` — see the existing comment at lines 133–146 explaining why)
   → `repoinstall.Options` → `repoinstall.Run`. A new
   `cmd_repo_uninstall.go` is expected to follow this exact shape 1:1:
   parse flags, resolve `gitRepoRoot`, build `Options`, call
   `repouninstall.Run`. `RepoInstallCmd.Spec()` (lines 104–112) is the
   template for `RepoUninstallCmd.Spec()`: `Name: "repo-uninstall"`,
   `Source: SourceKernel`, `Tier: TierKernel`, `NeedsHub: false` (uninstall
   must run even if `.cn/` is in a partially-torn-down state — same
   reasoning install uses for running before `.cn/` exists).
3. **`src/go/internal/repoinstall/repoinstall_test.go`** (1596 lines) and
   **`src/go/internal/cli/cmd_repo_install_test.go`** (662 lines) — existing
   test conventions: `httptest.Server` fixtures for GitHub API/release
   endpoints, `makeTarGz`/`writeLocalIndex`/`writeFixtureIndex` helpers,
   `t.TempDir()` + `git init` for repo-root fixtures. `repouninstall`'s test
   file should mirror the `_test.go` naming/fixture-builder pattern (test
   fixture duplication across packages is accepted precedent here — see
   `repoinstall_test.go`'s own `makeTarGz` doc comment citing eng/go §2.17
   as governing parser dedup, not test fixtures).
4. **`src/packages/cnos.core/commands/label-doctor`** (`doctor.go`,
   `github.go`, `manifest.go`, `resolve.go`, own `go.mod`, workspace-linked
   via root `go.work` line 9 — **no** `require`/`replace` pair needed in
   `src/go/go.mod`, confirmed: `grep -n label src/go/go.mod` returns
   nothing, yet `repoinstall.go` already imports it directly and `go build
   ./...` passes clean). `repoinstall.go`'s `ensureCanonicalDispatchLabels`
   (lines 498–531) is the in-process-call precedent AC4 should follow. **Gap
   α must close:** `github.go`'s `ghListLabels` (GET, paginated) exists, but
   there is **no delete primitive** in this package today — `doctor.go`'s
   `Audit`/repair path only creates/PATCHes. α either (a) adds an exported
   delete function to `labeldoctor` (extending the existing package,
   mirroring `ghListLabels`'s dependency-free `net/http` idiom exactly — no
   `gh` CLI shellout), or (b) implements a minimal GitHub label-DELETE call
   directly inside `repouninstall` reusing only `ghListLabels`-shaped
   listing via an exported `labeldoctor` list function if one exists (it
   does not today — `ghListLabels` is unexported). **This is a real
   sub-decision or the α prompt would under-specify it** — pinned in the
   Implementation contract below as "extend `labeldoctor` with an exported
   list + delete primitive" (mirrors how the package already grew a GET
   primitive over the `issues-fsm` precedent per `github.go`'s own doc
   comment).
5. **`src/go/cmd/cn/main.go`** — command registration. Line 51:
   `reg.Register(&cli.RepoInstallCmd{})`. `RepoUninstallCmd` gets registered
   immediately alongside it (same pattern as line 52's
   `reg.Register(&cli.LabelDoctorCmd{})` sitting next to `RepoInstallCmd`).
   **No other help-text change is required for AC1** — confirmed by reading
   `src/go/internal/cli/dispatch.go`'s `GroupMembers`/`PrintGroup` (lines
   112–155): `cn repo --help` / `cn repo` (no verb) is **not** a static help
   string; it is generated by `PrintGroup` iterating every registered
   command whose `Spec().Name` has the `"repo-"` prefix, in registry order.
   Registering `RepoUninstallCmd` with `Name: "repo-uninstall"` is
   sufficient for AC1 by construction — α should verify this with a test
   asserting the group listing, not by hand-editing a help string.
6. **cnos#706 landed-code check (explicit, per dispatch instructions):**
   `gh issue view 706 --json state,labels` → `state: OPEN`, label
   `status:review`, and its PR (whichever it is) is **not merged** —
   confirmed no manifest-writing code has landed on `main` under #706. AC2
   ("install writes a manifest") is **in scope of #707 itself**; α must not
   assume any prior-cycle manifest plumbing exists.

---

## AC oracle approach

| AC | Requirement | Oracle |
|---|---|---|
| AC1 | `cn repo uninstall` exists, appears in `cn repo --help` alongside `install` | Unit test in `cli` package asserting `dispatch.GroupMembers(reg, "repo")` contains both `repo-install` and `repo-uninstall`; a `cmd_repo_uninstall_test.go` test capturing `cn repo --help` (or bare `cn repo`) stdout and asserting it lists `uninstall` next to `install`, mirroring how `cmd_repo_install_test.go` already asserts help/flag-parse behavior |
| AC2 | `cn repo install` writes an install manifest recording side effects (files, gitignore line, vendor path, labels, workflow, dispatch config) | Extend `repoinstall_test.go`: after a non-dry-run `Run`, assert `.cn/install-manifest.json` exists, unmarshal it, and assert field-by-field JSON shape (schema `cn.install-manifest.v1`, tracked files list, `.gitignore` line recorded, vendor path, label names created, workflow path when `--dispatch cds`, dispatch config echoed). A second assertion: `--dry-run` writes **no** manifest (mirrors the existing dry-run-writes-nothing invariant already tested for `.cn/deps.json`) |
| AC3 | uninstall removes tracked files + gitignore line + on-disk vendor, driven by the manifest; idempotent; safe under prior partial removal | `repouninstall_test.go`: (a) fixture-install via `repoinstall.Run`, then `repouninstall.Run`, assert every manifest-recorded file/dir is gone and the `.gitignore` line is removed; (b) run `repouninstall.Run` **twice** in a row and assert the second run is a no-op success (idempotency); (c) manually delete one manifest-recorded file before calling uninstall and assert it does not error (partial-removal safety) |
| AC4 | uninstall lists install-created labels; deletes only with explicit opt-in + confirm; default leaves them | Table test over `--labels` flag states: (a) no flag → labels listed in stdout, zero DELETE calls hit the `httptest` GitHub fixture; (b) `--labels` without confirmation (however the confirm gate is shaped — e.g. requires `--yes` too, or an interactive prompt stubbed in tests) → still zero deletes; (c) `--labels --yes` (or equivalent explicit dual-opt-in) → DELETE calls issued for exactly the manifest-recorded label names, asserted against the fixture server's request log |
| AC5 | uninstall prints operator-only effects (secrets, PAT/bot grants) with removal instructions; never attempts to delete them | stdout-assertion test: run uninstall against a manifest that recorded `--dispatch cds` identity (agent/workflow-pat-secret/bot-name/bot-id, mirroring `runDispatchCds`'s own opts), assert the printed instructions name the secret/grant by name with manual-removal steps, and assert (via the `httptest` fixture, same technique as AC4) that **no** API call is made against any secrets/apps endpoint |
| AC6 | `--dry-run` shows exactly what would be removed; nothing removed unless shown | `repouninstall_test.go`: `--dry-run` run against a populated fixture asserts (a) stdout enumerates every file/label/workflow path that would be removed, matching AC3/AC4's real-run enumeration byte-for-byte in content (not just superset), and (b) the filesystem is byte-identical before/after (`os.Stat` on every manifest-recorded path, or a full tree diff) |

---

## Empirical anchor — precedent for CLI/package wiring conventions

Repo is a shallow local checkout (`git rev-parse
--is-shallow-repository` → `true`; local `git log -- <file>` on
long-lived files returns only the single squashed HEAD commit), so γ used
`gh api "repos/usurobor/cnos/commits?path=..."` to recover real history:

- **`cn repo install` itself**, the closest possible precedent for wiring a
  new `cn repo <verb>`: commit `03888d6b` — "cli: register cn repo install
  kernel command (cnos#608)", shipped under **PR #617** ("cds-install Sub 1
  (Cn=1 / PR 1): implement `cn repo install` — base installer"). Commit
  message states verbatim: *"RepoInstallCmd is a thin cli/ wrapper (eng/go
  §2.18) around internal/repoinstall... Registered in cmd/cn/main.go
  alongside the other kernel commands."* — i.e. exactly the pattern named
  in "Surfaces" §2/§5 above.
- **`cn label-doctor`** (the in-process label mechanism AC4 extends):
  commit `9681e85b` — "cnos#493: register `cn label-doctor` kernel command",
  described as mirroring `cmd_issues_fsm.go` exactly: *"Run delegates
  entirely to the labeldoctor package, no os/net/http/encoding/json/path/
  filepath import"* — the thin-wrapper pattern `cmd_repo_uninstall.go`
  should also follow for its own domain calls into `repouninstall`.

---

## Expected diff scope

**New files:**
- `src/go/internal/repouninstall/repouninstall.go` — domain package (`Args`,
  `Options`, `Result`, `ParseArgs`, `Run`), mirroring `repoinstall.go`'s
  shape (flag parsing → validation → manifest read → phased teardown →
  structured stdout, dry-run branch mirrors `printPlan`).
- `src/go/internal/repouninstall/repouninstall_test.go`.
- `src/go/internal/cli/cmd_repo_uninstall.go` — thin wrapper (`RepoUninstallCmd`, `Spec`, `Help`, `Run`), mirrors `cmd_repo_install.go`.
- `src/go/internal/cli/cmd_repo_uninstall_test.go`.
- Manifest type: either a new `installmanifest` type colocated in
  `internal/repoinstall` (e.g. `repoinstall/manifest.go`, schema
  `cn.install-manifest.v1`, exported so `repouninstall` can import it) or a
  shared type under `internal/pkg` alongside the existing `pkg.Manifest`
  (`cn.deps.v1`) — **α's call, not pinned here**, but whichever location is
  chosen, `repouninstall` reads the same exported type `repoinstall` writes
  (no duplicate parsing of the same JSON shape — eng/go §2.17).

**Files needing edits:**
- `src/go/internal/repoinstall/repoinstall.go` — `applyInstall` (or a new
  step called after it) writes `.cn/install-manifest.json` last, after
  every side effect it records has occurred; `Result` gains a manifest
  field if useful for tests.
- `src/go/internal/repoinstall/repoinstall_test.go` — AC2 assertions.
- `src/go/cmd/cn/main.go` — one new `reg.Register(&cli.RepoUninstallCmd{})`
  line beside line 51/52.
- `src/packages/cnos.core/commands/label-doctor/github.go` (and/or
  `doctor.go`) — new exported list/delete primitive(s) for AC4, extending
  (not duplicating) the existing dependency-free `net/http` idiom.
- `src/packages/cnos.core/commands/label-doctor/*_test.go` — tests for the
  new delete primitive.

**No edit needed to any static help string for AC1** (see Surfaces §5 —
`cn repo --help` is generated from the registry).

---

## Scope guardrails (restated from the issue body, verbatim intent)

**In scope:** the manifest written by install; the `uninstall` command;
label handling (list / opt-in delete); a printed operator-effects summary
(secrets/grants to remove manually); `--dry-run`; `cn repo --help` lists
both verbs.

**Out / deferred (explicit non-goals — α must not implement):**
- Deleting repo secrets or bot/PAT grants **programmatically** — uninstall
  prints removal instructions only; it never calls any GitHub secrets or
  Apps-installation API.
- Removing `.cdd/unreleased/*` cell artifacts — those are cell history, not
  install output; uninstall must not touch `.cdd/`.

---

## α prompt

```text
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 707 --json title,body,state,comments
Branch: cycle/707
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.core/skills/eng/go/SKILL.md
```

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go |
| CLI integration target | `cn` subcommand — `cn repo uninstall`, registered in `src/go/cmd/cn/main.go` via `reg.Register(&cli.RepoUninstallCmd{})` beside the existing `RepoInstallCmd`/`LabelDoctorCmd` registrations (line 51–52), resolved through the existing noun-verb dispatcher (`cli.ResolveCommand` builds `"repo"+"-"+"uninstall"`) — no new dispatch mechanism |
| Package scoping | New domain package `src/go/internal/repouninstall/` (mirrors `src/go/internal/repoinstall/` 1:1: `Args`/`Options`/`Result`/`ParseArgs`/`Run`); thin wrapper `src/go/internal/cli/cmd_repo_uninstall.go` (eng/go §2.18 cli/-boundary — domain logic never lives in `cli/`); the install-manifest write (AC2) is added to the **existing** `src/go/internal/repoinstall` package (extend `applyInstall`/`Run`, not a new package — it is install's own output); label list/delete primitives (AC4) extend the **existing** `src/packages/cnos.core/commands/label-doctor` package (workspace-linked via root `go.work`, already imported in-process by `repoinstall.go`'s `ensureCanonicalDispatchLabels` — same in-process pattern, not a subprocess) |
| Existing-binary disposition | Extend — this is the existing `cn` binary (`src/go/cmd/cn/main.go`); no new binary, no separate module entrypoint |
| Runtime dependencies | None new. Reuses `net/http` (already a `repoinstall`/`labeldoctor` dependency — no third-party GitHub client, mirroring `github.go`'s existing dependency-free idiom) |
| JSON/wire contract preservation | AC2 introduces a **new** wire schema α must define and pin: `.cn/install-manifest.json`, schema value `cn.install-manifest.v1` (naming convention mirrors `pkg.Manifest`'s `cn.deps.v1` in `src/go/internal/pkg/pkg.go`). This is additive-only — `.cn/deps.json` (`cn.deps.v1`) and `.cn/deps.lock.json` shapes are unmodified; `cn.install-manifest.v1` must not be confused with, or reuse, the schema name `cn.repo.state.v1` from the unrelated, unmerged cnos#656 ledger design (see "overlapping wave" note above) |
| Backward-compat invariant | `cn repo install` must remain byte-identical in its non-manifest output/behavior for every existing flag combination (deterministic `.cn/deps.json`/`.cn/deps.lock.json`, existing dry-run text) — the manifest write is purely additive. `cn repo uninstall` MUST handle the **no-manifest-present** case gracefully (any repo that ran install before this cycle, or under pre-#707 code, has zero `.cn/install-manifest.json`) — uninstall must not crash or silently corrupt state on a missing manifest; it must either refuse with a clear, actionable message, or perform a conservative best-effort/no-op consistent with AC3's "safe under prior partial removal" language. This case is not spelled out as its own AC in the issue body — pinning it here so α does not treat a missing manifest as unreachable |
```

## β prompt

```text
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 707 --json title,body,state,comments
Branch: cycle/707
```

β verifies α's diff against:
1. Every pinned axis in the Implementation contract above (Rule 7 —
   `beta/SKILL.md §"Role Rules" Rule 7`), independent of behavioral AC
   passage — in particular: manifest schema value is exactly
   `cn.install-manifest.v1`; `RepoUninstallCmd` is registered in
   `cmd/cn/main.go` the same way `RepoInstallCmd` is; no domain logic leaked
   into `cli/cmd_repo_uninstall.go` (thin-wrapper check against eng/go
   §2.18, same shape as `cmd_repo_install.go`); no new third-party
   dependency was added for label list/delete.
2. The AC oracle table above, AC1–AC6, independently re-derived from the
   diff and the issue body — not from α's `self-coherence.md` narrative.
3. The two explicit non-goals (no programmatic secret/grant deletion; no
   `.cdd/unreleased/*` cleanup) — β confirms the diff contains **zero**
   calls to any GitHub secrets/Apps-installation endpoint and **zero**
   filesystem operations under `.cdd/`.
4. The no-manifest backward-compat case (Implementation contract row 7) —
   confirm a test exists exercising `cn repo uninstall` against a repo with
   no `.cn/install-manifest.json` and that it fails safe / no-ops rather
   than panicking or deleting unrelated files.
