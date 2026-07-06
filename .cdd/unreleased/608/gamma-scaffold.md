# γ scaffold — cycle #608

**Issue:** [cnos#608](https://github.com/usurobor/cnos/issues/608) — "cds-install Sub 1 (Cn=1 / PR 1): implement `cn repo install` — base installer (dispatch none)"
**Mode:** `design-and-build` (per the issue's own `**Mode:**` field, confirmed against the fetched issue body: `**Mode:** design-and-build`)
**cell_kind:** `implementation` (per the issue's own `**cell_kind:**` field)
**Parent:** #607 (wave master tracker — "first-class CDS repo installer"). Sub-issue sequence per #607's operator-enacted launch verdict: PR1 = #608 (this cycle, base installer) · PR2 = #609 (renderer generalization) · PR3 = #610 (dispatch install) · PR4 = #611 (bootstrap delegation). Only #608 is dispatched now; #609–#613 remain filing-only/operator-gated.
**Base SHA:** `3dad64285026582a161549b8fd10108dd67a369e` (verified via `git fetch origin main && git rev-parse origin/main` at branch-creation time — this is later than the wake-invocation-time SHA `e49eaa1...`; other automation had pushed board-map commits in between. Pre-flight re-run against the live SHA, not the stale one.)
**Branch:** `cycle/608`, created from `origin/main` at the base SHA above, pushed to `origin/cycle/608`.
**Branch pre-flight (all confirmed before branch creation):** `origin/cycle/608` did not exist (`git ls-remote --heads origin cycle/608` empty) · no stalled `.cdd/unreleased/608/` beyond the dispatch-wake's `CLAIM-REQUEST.yml` marker (`git ls-tree -r --name-only origin/main -- .cdd/unreleased/608/` → exactly one file, the claim marker, which is expected and not stale cycle matter) · issue scope declared in the body · issue state `OPEN` (`gh issue view 608 --json state` → `OPEN`).

---

## CRITICAL — Design surface cited by the issue is not yet on `origin/main`

The issue's own header cites: **"Design surface: `docs/development/design/cn-repo-install-MOCKS.md` (Mocks A, B, E1)."** I verified this path via `find`/`git log --all` and **it does not exist anywhere in `origin/main`'s history** (`git log --all --oneline -- '*cn-repo-install-MOCKS*'` on `origin/main` alone returns nothing).

It **does** exist, however: parent issue #607's own Source-of-truth table names it as **"On branch `claude/cds-install-guide-2ka54j`"**. I fetched that branch (`git fetch origin claude/cds-install-guide-2ka54j`) and confirmed:
- `git merge-base origin/main origin/claude/cds-install-guide-2ka54j` == `origin/main`'s current tip (`3dad6428...`) — i.e. that branch is a strict descendant of the current main tip, not a diverged fork. It is safe to pull the one file from it without any merge-conflict risk against this cycle's base.
- The file exists there at `docs/development/design/cn-repo-install-MOCKS.md` (296 lines), status header: **"Design surface (pre-implementation). Author: kappa/operator review."** — i.e. already operator-reviewed/pinned content (#607's "Launch verdict (operator) — enacted" section references and pins the same command-name/lockfile decisions this file encodes), not a draft α needs to redesign from scratch.

**α's obligation (part of this cycle's `design-and-build` design step):** land this file at the canonical path on `cycle/608`, e.g.:

```
git fetch origin claude/cds-install-guide-2ka54j
git show origin/claude/cds-install-guide-2ka54j:docs/development/design/cn-repo-install-MOCKS.md \
  > docs/development/design/cn-repo-install-MOCKS.md
git add docs/development/design/cn-repo-install-MOCKS.md
```

Do **not** re-author it from prose memory — pull the exact reviewed content, then note in `self-coherence.md` that it was landed via this cycle's design step (design-and-build: "design needed before code; both happen in this cycle" — this satisfies that, since the design was already drafted and operator-reviewed elsewhere, this cycle is the one that makes it canonical on `main`'s lineage).

The `claude/cds-install-guide-2ka54j` branch also carries `docs/guides/INSTALL-CDS.md` and `docs/guides/templates/cnos-install.yml` drafts, plus unrelated changes (deletions of already-released `.cdd/unreleased/{593,600}/` dirs, `.github/workflows/*` edits, `.cn-sigma/logs/*`) that are **not** part of #608's scope — those belong to later subs (#609–#611) or are stale housekeeping from whoever authored that branch. **α must cherry-pick only the specific design/docs files #608's own scope names** (the MOCKS.md file for the design surface; `INSTALL-CDS.md` per the issue's own §Docs section — but reconcile its content against #608's actual base-only scope, since the drafted doc may already describe dispatch/#609/#610 material that is out of scope here). Do not pull the whole branch or its unrelated files.

I am not treating this as a dispatch-blocking ambiguity because: (a) the content is locatable, reviewed, and its exact retrieval command is now named here; (b) the mode is `design-and-build`, which explicitly allows/expects the design artifact to land within this same cycle; (c) the parent issue #607 already documents the branch/file relationship, so this is a findable fact, not an invented one. If α cannot retrieve the branch (e.g. it has been deleted or force-pushed away by the time α runs), α must STOP and escalate to γ/δ rather than inventing Mock A/B/E1 content from scratch — the exact invariant tables are reproduced below as a fallback so α is not blocked even if the branch becomes unreachable.

### Fallback — Mocks A, B, E1 content (verbatim from the design branch, in case it becomes unreachable)

**Mock A — human-deliverable invariants A1–A5** (`cn repo install --dry-run`, base mode):

| ID | Invariant |
|---|---|
| A1 | Fails with a clear error if not at a git repo root (does not silently walk up or scaffold). |
| A2 | Prints the exact release tag it will use; `--release latest` resolves to a concrete tag in the output. |
| A3 | Lists the precise planned committed diff — and it is exactly `.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore`. |
| A4 | States dispatch mode explicitly and, in base mode, states no `.github/workflows/` change. |
| A5 | Writes nothing (verified: `git status --porcelain` empty after `--dry-run`). |

**Mock B — invariants B1–B4** (base-install committed diff):

| ID | Invariant |
|---|---|
| B1 | Diff is exactly three files (`.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore`); no `.github/workflows/` file present. |
| B2 | `.cn/deps.lock.json` validates against `cn.lock.v2` and pins every package in `deps.json` by SHA-256. Versions are **exact** (resolved from the index), not semver ranges. |
| B3 | **Idempotent**: a second `cn repo install` produces no further diff (`git status --porcelain` empty), with no formatting churn. |
| B4 | **Dispatch guard**: until the renderer is generalized (#609), `cn repo install --dispatch cds` fails with an explicit "not implemented until renderer generalization lands" error and writes **no** partial `.github/workflows/cnos-cds-dispatch.yml`. |

**Mock E1** (tenant-portability surface, base-install half only — E2–E4 belong to #609/#610):

| ID | Invariant |
|---|---|
| E1 | `cn repo install` (base) writes **no** agent-hub scaffold — no `spec/SOUL.md`, `agent/`, `threads/`, `state/`; only `.cn/deps.json` + `.cn/deps.lock.json` + `.gitignore` (cnos#606 C4). |

**Receipt parity contract shape** (from the design doc's own §"Receipt parity contract" — α's closeout MUST emit this block, per the issue's own "Parity requirement" section):

```yaml
mock_parity:
  source: docs/development/design/cn-repo-install-MOCKS.md
  source_commit: "<sha of this file the cell built against>"
  rows:
    - id: A1   # one row per A1-A5, B1-B4, E1 (9 rows total for this sub)
      expectation: "..."
      observed: "<what the built command actually printed/did>"
      evidence: "<test name / transcript path / commit>"
      verdict: match | exceed | miss
      how: "<why it matches; if exceed, why safe; if miss, gap + disposition>"
  summary:
    matched: <n>
    exceeded: <n>
    missed: 0          # MUST be 0 to converge
    exceed_justified: <bool>
```

---

## Surfaces α is expected to touch

Concrete, grounded by direct reading (not asserted from the issue text alone):

| Path | Current state | Expected role in this cycle |
|---|---|---|
| `src/go/internal/cli/cmd_repo_install.go` | **Does not exist.** | **New file.** Implements the `cn repo install` command. Per the CLI's own dispatch convention (see below), this should be a `RepoInstallCmd` type registered under `Spec().Name == "repo-install"` — noun-verb resolution (`cli.ResolveCommand` in `src/go/cmd/cn/main.go` builds the lookup key `args[0]+"-"+args[1]` = `"repo"+"-"+"install"` = `"repo-install"`), exactly the same pattern `cell-finalize`/`issues-fsm` already use for `cn cell finalize`/`cn issues fsm`. Confirmed by reading `src/go/internal/cli/dispatch.go` `ResolveCommand` (lines 42-88) and `src/go/internal/cli/cmd_cell.go` (`Name: "cell-finalize"`), `cmd_issues_fsm.go` (`Name: "issues-fsm"`). |
| `src/go/cmd/cn/main.go` | Registers 14 kernel commands via `reg.Register(&cli.XCmd{})` (lines 29-49). | Add `reg.Register(&cli.RepoInstallCmd{})` alongside the existing registrations. No other change to `main.go`'s dispatch/help/hub-discovery logic is expected. |
| `src/go/internal/restore/restore.go` | Existing, shipped. `GenerateLockFromIndex(hubPath, indexPath) (*LockResult, error)` (line 448) reads `.cn/deps.json`, resolves exact-match pins against a `PackageIndex`, writes `.cn/deps.lock.json` with `Schema: "cn.lock.v2"`, sorts deterministically. `Restore(ctx, hubPath, indexPath) ([]Result, error)` (line 118) does lockfile → index lookup → fetch (HTTP or local) → SHA-256 verify → tar extract → `cn.package.json` validate. **Reuse both directly — do not reimplement.** No change to this file is expected unless α finds a genuine gap (e.g. index normalization for relative-vs-absolute tarball URLs, per Implementation requirement 4) that must live here rather than in the new installer file. |
| `src/go/internal/pkg/pkg.go` | Existing, shipped. `Manifest{Schema, Profile, Packages []ManifestDep}` (schema value `"cn.deps.v1"`, confirmed against Mock B's diff sample), `Lockfile{Schema, Packages []LockedDep}` (schema `"cn.lock.v2"`), `PackageIndex{Schema, Packages map[name]map[version]IndexEntry}` (schema `"cn.package-index.v1"`, confirmed in `pkg_test.go`). **Preserve these schemas exactly** — do not add fields or change shape. |
| `src/go/internal/binupdate/binupdate.go` | Existing, shipped. `FetchLatestRelease(ctx, client, baseURL, repo) (*Release, error)` (line 221) calls `GET {baseURL}/repos/{repo}/releases/latest`. **Reuse candidate** for `--release latest` resolution (issue's Implementation requirement 4) rather than reimplementing a GitHub-releases-API client from scratch. α should assess whether to call this directly or extract a shared helper — either is acceptable, but do not hand-roll a second GitHub-API client. |
| `docs/development/design/cn-repo-install-MOCKS.md` | **Does not exist on `origin/main`.** Exists (reviewed, "Enacted" per #607) on `origin/claude/cds-install-guide-2ka54j`. | **Land verbatim** (see CRITICAL section above) as this cycle's design artifact. |
| `docs/guides/INSTALL-CDS.md` | Existing on `main` (pre-installer content) — check current content before editing; a draft also exists on the design branch. | Update per issue §Docs: canonical base path becomes `cn repo install`; keep the two-layer framing (Layer 1 package install / Layer 2 autonomous dispatch, opt-in). Reconcile against the design-branch draft but do not pull dispatch/#609/#610-scoped content prematurely. |
| `.gitignore` (repo root, target repo being installed into — i.e. the *installed* repo's `.gitignore`, not cnos's own) | N/A — this is a runtime-written file in the *target* repo, not cnos's own `.gitignore`. cnos's own root `.gitignore` currently has no `.cn` entry (`grep -n '\.cn' .gitignore` → no hits) — irrelevant to this cycle except as a sanity check that cnos itself isn't accidentally the install target in a test. | The installer command must ensure the **target repo's** `.gitignore` contains `.cn/vendor/` (Mock A/B: exactly 3 committed files — `.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore`). |
| Test fixtures (new) | N/A | Unit + integration tests per issue §Tests: git-root detection, deterministic `.cn/deps.json` write, exact-version resolution from a fixture `index.json`, missing-package/missing-index errors, relative-tarball-URL normalization, HTTP-tarball-URL preservation, SHA-mismatch propagation, `--dry-run` no-write, idempotence, plus an integration test (temp git repo + fixture index/tarballs) and the dispatch-guard test (`--dispatch cds` fails explicitly, no partial workflow file). Follow existing test conventions in `src/go/internal/restore/restore_test.go` and `src/go/internal/cli/cmd_activate_test.go` / `registry_test.go` for fixture shape and `Invocation`-based command testing. |

---

## AC oracle approach

| AC | Concrete check |
|---|---|
| **AC1** | `cn repo install` resolves via the kernel registry (`RepoInstallCmd`, `Spec().Name == "repo-install"`, `Source: SourceKernel`, `NeedsHub: false`) and runs successfully from a fresh `git init`-only repo with only `cn` on PATH — no `.cn/` hub, no vendor packages, no prior state. Oracle: integration test — temp dir, `git init`, install only the `cn` binary (build once), run `cn repo install --index <fixture>`, assert exit 0 and the three-file diff exists. |
| **AC2** | `.cn/deps.json` byte-identical across two successive runs with the same inputs. Oracle: unit/integration test running the write twice, diffing bytes; also verify stable package ordering (sorted by name) and no timestamp field. |
| **AC3** | `.cn/deps.lock.json` validates `Schema == "cn.lock.v2"`; every entry has exact `version` (not a range) and non-empty `sha256`. Oracle: parse the written lockfile with `pkg.ParseLockfile` (or equivalent) and assert schema + per-entry SHA-256 presence; reuse `GenerateLockFromIndex`'s existing determinism (sorted by name then version). |
| **AC4** | `.cn/vendor/packages/cnos.core/cn.package.json`, `.../cnos.cdd/cn.package.json`, `.../cnos.cds/cn.package.json` all exist post-restore; checksums verified via `restore.Restore`'s existing SHA-256 verify step (not reimplemented). Oracle: integration test asserting file existence + a deliberately-corrupted fixture tarball produces a SHA-mismatch error (proves the existing verify path is actually wired, not bypassed). |
| **AC5** | `cn repo install; git diff --exit-code; cn repo install; git diff --exit-code` — both exit 0. Oracle: shell/integration test exactly reproducing this sequence against a fixture index. |
| **AC6** | `--dry-run`: `git status --porcelain` empty after; stdout lists exactly `.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore` as the planned diff (Mock A3/A5). Oracle: integration test capturing stdout + porcelain status. |
| **AC7** | Three invocations tested independently: `--release latest` (resolves a concrete tag, e.g. via `binupdate.FetchLatestRelease` or equivalent), `--release <tag>` (pinned), `--index <path-or-url>` (fixture path and, if feasible, an `http://` fixture server). Oracle: one test per mode; `--release latest` test may need to stub/mock the GitHub API call rather than hitting it live in CI. |
| **AC8** | `grep -RIl 'cnos-cds-dispatch' .github/workflows/` on the *target* test repo after a base install returns nothing; grep the installer's own source for any secret-name reference (`SIGMA_WORKFLOW_PAT`, `CNOS_WORKFLOW_PAT`, `CLAUDE_CODE_OAUTH_TOKEN`) outside the explicitly-out-of-scope `--dispatch cds` error-path string. Oracle: integration test asserting no `.github/workflows/` dir is created at all in base mode (default `--dispatch none`). |
| **AC9** | `cn repo install --dispatch cds` exits non-zero with the exact error class named in the issue's §Tests ("Dispatch guard": `✗ --dispatch cds requires generalized wake renderer support (#609)` or equivalent clear message) and leaves zero `.github/workflows/` files behind. Oracle: integration test asserting exit code, stderr message content, and `git status --porcelain` shows no `.github/` changes. |
| **AC10** | No `spec/SOUL.md`, `agent/`, `threads/`, `state/` directories created by `cn repo install` (contrast with `cn init`'s existing agent-hub scaffold — cnos#606 C4). Oracle: integration test asserting these paths do not exist post-install; explicit regression guard distinguishing installer behavior from `cn init`'s. |
| **AC11** | `docs/guides/INSTALL-CDS.md` names `cn repo install` as the canonical base path. Oracle: grep/read the doc post-edit; no lingering "7 manual steps" framing presented as the only path. |

---

## Empirical anchor citations

- **cnos#389** (Python-not-Go) / **cnos#391** (wrong package scoping + separate binary) / **cnos#392** (recovery — δ pinned the 7-axis contract ad hoc) / **cnos#393** (codification — the `## Implementation contract` block + δ-inward-membrane doctrine) — directly load-bearing for *this* cycle: the issue's own Implementation requirement 1 ("Kernel command, not a package command… registered in `cmd/cn/main.go`") is exactly the class of axis (CLI integration target, package scoping) that #389/#391 show α will improvise on if left unpinned. The 7-axis table in the α prompt below pins these explicitly so this cycle does not repeat that failure mode.
- **cnos#606** (tenant dogfooding, C1–C7) — the direct customer-evidence source for this issue's C4 (agent-hub scaffold → AC10) and general onboarding-gap framing. Cited by the issue itself; not independently re-verified beyond reading #607's C1–C7 table.
- **cnos#607** — parent wave tracker; source of the PR-sequencing verdict (base-first, dispatch-gated) and the design-branch pointer this scaffold resolves above.

---

## Expected diff scope (coarse — not line-by-line)

- **New:** `src/go/internal/cli/cmd_repo_install.go` (the command implementation) + its `_test.go`.
- **New:** possibly a small helper package if α judges the index-normalization / release-resolution logic (Implementation requirement 4) too large for the command file alone — acceptable, but must live under `src/go/internal/` following existing package-per-concern conventions (cf. `restore/`, `binupdate/`, `pkg/`, `hubsetup/`), not a new top-level directory.
- **Modified:** `src/go/cmd/cn/main.go` — one new `reg.Register(&cli.RepoInstallCmd{})` line.
- **New (landed this cycle, per CRITICAL section):** `docs/development/design/cn-repo-install-MOCKS.md`.
- **Modified:** `docs/guides/INSTALL-CDS.md` per §Docs.
- **New:** test fixtures (a small fixture `index.json` + fixture tarballs, or fixture HTTP server, for the integration tests) — likely under `src/go/internal/cli/testdata/` or `src/go/internal/restore/testdata/` following whatever convention `restore_test.go` already uses; check before inventing a new fixture location.
- **No change expected:** `src/go/internal/restore/restore.go` internals (reused, not modified) unless a genuine gap surfaces in URL normalization; `src/go/internal/pkg/pkg.go` schemas (must not change shape); `cmd_init.go`, `cmd_setup.go`, `cmd_deps.go` (non-goal: must not alter `cn init`/`cn setup`/`cn deps` behavior); `.github/workflows/*` (base install must never touch this — a change here in the diff is itself a signal something has gone wrong, per AC8/AC9).

---

## Scope guardrails (restated by reference — see issue body for full text)

- **Out of scope / non-goals** (issue's own §Non-goals): hosted registry design, full semver-range resolution, direct push to `main`, automatic autonomous dispatch, sigma-only assumptions, Claude-specific activation, agent-runtime impl, package-substrate redesign.
- **Non-negotiable launch cut**: base install must not install the autonomous dispatch workflow under any code path reachable without `--dispatch cds`, and `--dispatch cds` itself must fail cleanly (not partially render) until #609 lands.
- **Risk controls** (issue's own): don't commit tarball/generated cache into the target repo (only `.cn/deps.json` + `.cn/deps.lock.json` + `.gitignore` entry are committed; tarballs live in cache/tmp); don't duplicate `cn deps lock`/`restore` logic — reuse the internal functions directly.

---

## α prompt

**Branch:** `cycle/608` (already created and pushed from `origin/main` at `3dad64285026582a161549b8fd10108dd67a369e`; switch to it, do not create a new branch.)

**Task:** Implement `cn repo install` (base installer, dispatch none) per cnos#608. Before writing code: (1) land `docs/development/design/cn-repo-install-MOCKS.md` at the canonical path per this scaffold's CRITICAL section (pull from `origin/claude/cds-install-guide-2ka54j`, do not re-author); (2) re-read `src/go/internal/restore/restore.go` (`GenerateLockFromIndex`, `Restore`) and `src/go/internal/pkg/pkg.go` in full before writing any new code — these are reuse targets, not reference points to reimplement; (3) implement `src/go/internal/cli/cmd_repo_install.go` as a `RepoInstallCmd` registered under `Spec().Name == "repo-install"` (noun-verb resolution, matching the `cell-finalize`/`issues-fsm` pattern in `dispatch.go`/`cmd_cell.go`/`cmd_issues_fsm.go`), `NeedsHub: false`; (4) satisfy AC1–AC11 exactly per the oracle table above; (5) emit the `mock_parity` block (9 rows: A1–A5, B1–B4, E1) in the closeout per the design doc's own Receipt parity contract, `missed: 0`. If the design branch (`claude/cds-install-guide-2ka54j`) is unreachable, use the verbatim Mock A/B/E1 tables reproduced in this scaffold's fallback section instead of inventing content, and note the substitution in `self-coherence.md`.

**## Implementation contract (pinned by δ; α MUST NOT improvise)**

| Axis | Pinned value |
|---|---|
| Language | Go — matches the existing `cn` kernel (`src/go/...`); no other language is introduced. |
| CLI integration target | `cn` subcommand — kernel command registered as `RepoInstallCmd` with `Spec().Name == "repo-install"`, `Source: SourceKernel`, `Tier: TierKernel`, `NeedsHub: false`. Invoked as `cn repo install` via the existing noun-verb resolver (`cli.ResolveCommand`, `args[0]+"-"+args[1]` lookup) — the same mechanism already used for `cn cell finalize` / `cn issues fsm`. Not a standalone binary, not a package-vendored command (it must run before any hub/package state exists). |
| Package scoping | New file `src/go/internal/cli/cmd_repo_install.go` (+ `_test.go`), registered in `src/go/cmd/cn/main.go`. Reuses `src/go/internal/restore/` (`GenerateLockFromIndex`, `Restore`, `ReadPackageIndex`, `FindIndexPath`), `src/go/internal/pkg/` (Manifest/Lockfile/PackageIndex types), and `src/go/internal/binupdate/` (`FetchLatestRelease`, or an extracted equivalent) for release resolution. Any new helper package must live under `src/go/internal/` following the existing package-per-concern convention (cf. `restore/`, `binupdate/`, `pkg/`, `hubsetup/`) — no new top-level directory. |
| Existing-binary disposition | Coexist — new primary command; no existing `cn repo install` to replace or deprecate. `cn init`, `cn setup`, `cn deps` remain byte-for-byte unchanged in behavior (non-goal; AC10 is the explicit regression guard distinguishing this installer from `cn init`'s agent-hub scaffold). |
| Runtime dependencies | None beyond the existing Go toolchain and the standard library (`net/http` for index/tarball/release fetches, already used by `restore.go` and `binupdate.go`). No new external dependency (module) is introduced. |
| JSON/wire contract preservation | Preserve exactly: `.cn/deps.json` schema `"cn.deps.v1"` (`pkg.Manifest`), `.cn/deps.lock.json` schema `"cn.lock.v2"` (`pkg.Lockfile`), package-index schema `"cn.package-index.v1"` (`pkg.PackageIndex`). No field additions, no shape changes. Vendored-package layout stays name-based: `.cn/vendor/packages/<name>/`, not `<name>@<version>/`. |
| Backward-compat invariant | `cn init`, `cn setup`, `cn deps restore`, `cn deps lock` behavior is unchanged — no shared code path is modified in a way that alters their existing output or exit codes. New content (the installer command, its reused-function call sites, the design doc, the docs update) is strictly additive to the existing kernel surface. |

---

## β prompt

**Branch:** `cycle/608`. **Oracle list:** AC1–AC11 exactly as stated in the issue body and expanded in this scaffold's "AC oracle approach" table above.

Additional gates:
- Standard `beta/SKILL.md` §"Pre-merge gate" (full row set — canonical-skill staleness re-check, merge-commit close-keyword check, etc.).
- **Rule 7 (implementation-contract conformance)** is binding: verify the diff matches every pinned axis in the α prompt's `## Implementation contract` table above. In particular: `cn repo install` must resolve through the kernel registry as a genuine `cn` subcommand (not a separate binary, not shelling out to `cn` internally per Implementation requirement 5), must reuse `restore.Restore`/`GenerateLockFromIndex` rather than reimplementing lock/restore logic, and must not touch `.github/workflows/*` in base mode under any code path. A behaviorally-correct implementation that nonetheless ships as a separate binary, reimplements SHA-256 verification, or changes the `cn.lock.v2`/`cn.deps.v1`/package-index schema shape fails Rule 7 even if AC1–AC11's behavioral checks otherwise pass.
- Confirm the `mock_parity` block exists in the closeout with exactly 9 rows (A1–A5, B1–B4, E1), `missed: 0`, and every `exceed` row (if any) carries a non-empty, safety-justified `how`.
- Confirm `docs/development/design/cn-repo-install-MOCKS.md` was landed verbatim (or, if the source branch was unreachable, that α explicitly noted the fallback substitution in `self-coherence.md` rather than silently inventing content).
- Confirm no `.github/workflows/*.yml` file appears in the diff at all (base install must never touch this surface — treat any appearance as a hard finding, not a judgment call).
- Confirm AC9's dispatch-guard test actually exercises `--dispatch cds` and asserts both the explicit failure message and the absence of a partial `.github/workflows/cnos-cds-dispatch.yml` write.

---

## Friction notes (open questions for α/β — not resolved here, do not guess past them)

1. **Design-branch reachability is the single largest risk to this scaffold's plan.** I verified `origin/claude/cds-install-guide-2ka54j` exists and is a clean descendant of the current `main` tip as of this scaffold's authoring time. If that branch is deleted, force-pushed, or otherwise unreachable by the time α runs, use the verbatim fallback tables above and flag the substitution — do not silently invent Mock content that looks similar but isn't the operator-reviewed original.
2. **`docs/guides/INSTALL-CDS.md` on `main` today may already have a different/older shape** than the draft on the design branch — I did not diff the two in detail (only confirmed the design-branch draft is 256 lines of new/changed content per the `git diff --stat` I ran). α should read both and reconcile, keeping only the base-install-scoped material for this PR (the two-layer framing minus any #609/#610-specific dispatch-flag detail that oversteps this sub-issue's scope).
3. **`--release latest` resolution mechanism is my inference, not a settled design decision.** I identified `binupdate.FetchLatestRelease` as a strong reuse candidate (same GitHub-releases-API shape this issue needs), but the issue does not name it explicitly. α may reuse it directly, extract a shared helper, or write a narrower purpose-built call — any of these satisfies the "Runtime dependencies: None new" contract row as long as no new external library is introduced. Do not treat "which exact function signature" as pinned; only the *behavior* (resolve latest release + its package index) is pinned by AC7.
4. **Test-fixture location convention was not verified in depth** — I read `restore_test.go`'s existence but did not open it fully to confirm its exact fixture-directory convention. α should check `src/go/internal/restore/restore_test.go` and `src/go/internal/cli/cmd_activate_test.go` / `registry_test.go` before inventing a new fixture path.
5. **The wave issue #607 lists #608 as "PR1, LAUNCHED"** and all other subs (#609–#613) as "filing-only" / operator-gated — this scaffold and the α/β prompts above are scoped strictly to #608. Do not begin any #609–#613-scoped work (renderer generalization, `--dispatch cds` actually rendering a workflow, bootstrap delegation, CLI ergonomics, PAT-free engine) even if it seems like an easy adjacent win — AC9 requires `--dispatch cds` to explicitly *fail*, not to work.
