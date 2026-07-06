# self-coherence — cycle #608

**manifest:** sections = [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
**completed:** [Gap, Skills, ACs, Self-check]

---

## Gap

**Issue:** [cnos#608](https://github.com/usurobor/cnos/issues/608) — "cds-install Sub 1 (Cn=1 / PR 1): implement `cn repo install` — base installer (dispatch none)"

**Mode:** `design-and-build` (per the issue body's own `**Mode:**` field).
**cell_kind:** `implementation` (per the issue body's own `**cell_kind:**` field).

**Parent:** #607 (wave master tracker — "first-class CDS repo installer"). This
cycle is PR1 of the wave (base installer, dispatch none). PR2 (#609, renderer
generalization), PR3 (#610, dispatch install), PR4 (#611, bootstrap delegation)
are explicitly out of scope for this cycle — only #608 is dispatched now.

**Gap being closed:** cnos ships the hub/package path (`cn init`, `cn setup`,
`cn deps restore`, `cn doctor`, `cn build`, `cn update`) but has no first-class
"install CNOS/CDS into this repo" command. Onboarding an arbitrary Git
repository into the CDS method requires ~7 manual steps (install `cn`, fetch
the release index, hand-write `.cn/deps.json`, run `cn deps lock`, run
`cn deps restore`, edit `.gitignore`, commit) and `cn init` is the wrong tool
for this — it scaffolds a full agent-hub (`spec/SOUL.md`, `agent/`, `threads/`,
`state/`) into a `cn-<name>/` subdirectory, not a package-substrate install
into the current repo (cnos#606 C4).

**What closes it:** a new kernel command, `cn repo install`, that:
- runs in a plain `git init`-only repo before any `.cn/` hub or package state
  exists;
- resolves a cnos release (default: latest) or an explicit `--index`
  path/URL;
- writes deterministic `.cn/deps.json` (schema `cn.deps.v1`) pinning the
  default package set (`cnos.core`, `cnos.cdd`, `cnos.cds`) to exact,
  index-resolved versions;
- reuses the existing `internal/restore` substrate directly
  (`restore.GenerateLockFromIndex`, `restore.Restore`) to write
  `.cn/deps.lock.json` (schema `cn.lock.v2`, SHA-256-pinned) and restore
  packages under `.cn/vendor/packages/<name>/` (name-based, not
  `<name>@<version>/`);
- ensures the target repo's `.gitignore` excludes `.cn/vendor/`;
- supports `--release latest|<tag>`, `--index <path-or-url>`,
  `--packages <csv>`, `--dry-run`;
- gates autonomous dispatch install behind an explicit, currently-failing
  `--dispatch cds` flag (default `--dispatch none`) — base install never
  writes `.github/workflows/` and never requires a workflow/agent secret;
  dispatch install itself is #609/#610 territory, out of scope here.

**Design surface:** `docs/development/design/cn-repo-install-MOCKS.md`
(Mocks A, B, E1 — the base-install half of the wave's design doc; Mocks
C/D/F/G belong to later subs). This file did not exist on `origin/main` at
cycle start; per γ's scaffold it was pulled verbatim (operator-reviewed,
"Enacted" per #607) from `origin/claude/cds-install-guide-2ka54j` rather than
re-authored — landed in commit `f22d5a78` (pre-rebase) / `40a6706c`
(pre-rebase second push), now at `d783ba38`'s ancestor
`f22d5a78` → after the rebase onto `origin/main` documented below, the
current SHA for that commit is unchanged in content (rebase was a clean,
non-conflicting fast-forward-style replay — no file in this cycle's diff
overlapped the intervening `origin/main` commits).

**Base SHA (original, at branch creation):** `3dad64285026582a161549b8fd10108dd67a369e`
**Base SHA (current, after α's pre-review-gate rebase):** `80778d688e04e61c66d38f2bd5962fafb0729e95`
(`origin/main` advanced by two unrelated automation commits — board-map
regeneration + a sigma heartbeat log entry — during this cycle; α rebased
`cycle/608` onto the new tip per alpha/SKILL.md §2.6 row 1. No conflicts;
the two commits touch `docs/development/board/*` and `.cn-sigma/logs/*`,
disjoint from every path this cycle's diff touches.)

---

## Skills

**Tier 1 (lifecycle):**
- `CDD.md` — canonical lifecycle/role contract.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` — this role's execution surface (§2.1
  dispatch intake, §2.5 self-coherence, §2.6 pre-review gate, §2.7 request
  review, §3.6 implementation-contract constraint).
- `design/SKILL.md` not separately re-loaded as a standalone step: the design
  artifact (`cn-repo-install-MOCKS.md`) was already operator-reviewed/pinned
  per γ's scaffold; this cycle's "design" work was landing it verbatim, not
  authoring new design content requiring the design skill's judgment
  algorithm.
- `plan/SKILL.md` not separately invoked as a standalone artifact: γ's
  scaffold already carried a concrete "Surfaces α is expected to touch" table
  + AC oracle mapping + expected diff scope, which functioned as the plan.
  Sequencing was straightforward (land design doc → domain package → cli
  wrapper → registration → tests → docs → self-coherence) with no
  non-trivial ordering decisions requiring the plan skill's own algorithm.

**Tier 2 (`cnos.eng`):**
- `eng/go/SKILL.md` — load-bearing throughout. §2.18 ("Dispatch boundary:
  cli/ owns dispatch only") directly shaped the repoinstall/cli split;
  §2.17 (parse/read split, parallel-parser prohibition) governed reuse of
  `pkg.ParsePackageIndex`/`pkg.ParseManifest`/`pkg.ParseLockfile` rather than
  inventing new ad hoc JSON decoding; the determinism conventions (sorted
  map-key iteration, no map-order-dependent output) shaped
  `writeManifest`/`rewriteRelativeEntriesFromBase`.

**Tier 3 (issue-specific):** none loaded beyond `eng/go` — this is a
single-language (Go) kernel-command cycle with no CLI-ergonomics-package,
security-review, or other Tier-3 skill named by the issue or scaffold.

**Implementation contract (δ-pinned, per alpha/SKILL.md §3.6 — not
improvised):** all 7 axes were populated in γ's scaffold (Language: Go;
CLI integration target: kernel `cn` subcommand via noun-verb resolution;
Package scoping: `cli/cmd_repo_install.go` + new `internal/repoinstall/`
package, reusing `restore/`, `pkg/`, `binupdate/`; Existing-binary
disposition: coexist, `cn init`/`setup`/`deps` unchanged; Runtime
dependencies: none new; JSON/wire contract preservation: `cn.deps.v1` /
`cn.lock.v2` / `cn.package-index.v1` schemas unchanged, vendor path
name-based; Backward-compat invariant: `cn init`/`setup`/`deps` behavior
unchanged). No row was empty or TBD at dispatch time, so no γ/δ escalation
was needed before coding. Every diff hunk in this cycle maps to one or more
of these pinned rows — see §CDD Trace step 6 for the file-by-file mapping.

---

## ACs

All commands below were run against branch HEAD at commit `a3f9b86e` (the
implementation SHA is fixed at `d783ba38` — the last pure-implementation
commit before the rebase and before this self-coherence artifact's own
commits; SHA convention per alpha/SKILL.md §2.6 "SHA convention for
readiness signal" — implementation SHA, not the advancing artifact-commit
HEAD).

### AC1 — kernel command, available from a fresh Git repo, no prior state

**Claim:** `cn repo install` resolves via the kernel registry
(`RepoInstallCmd`, `Spec().Name == "repo-install"`, `Source: SourceKernel`,
`NeedsHub: false`) and runs from a fresh `git init`-only repo with no `.cn/`
hub, no vendor packages, no prior state.

**Evidence:**
- `src/go/internal/cli/cmd_repo_install_test.go::TestRepoInstall_ResolvesViaNounVerb`
  — asserts `ResolveCommand(reg, []string{"repo","install",...})` resolves to
  `repo-install` with `Source: SourceKernel`, `Tier: TierKernel`,
  `NeedsHub: false`. PASS.
- `src/go/internal/cli/cmd_repo_install_test.go::TestRepoInstall_FreshGitRepo_EndToEnd`
  — `git init` in a temp dir, `t.Chdir` into it, run `RepoInstallCmd.Run`
  with a fixture `--index`; asserts `.cn/deps.json`, `.cn/deps.lock.json`,
  `.gitignore`, and `.cn/vendor/packages/cnos.core/cn.package.json` all
  exist afterward. PASS.
- Manual smoke test (built binary, real `git init`, no `.cn/` present
  beforehand): `cn help` lists `repo-install` in the kernel-commands section
  before any hub exists; `cn repo` (bare noun) lists the `repo install`
  subcommand via the noun-group printer, matching the `cell`/`issues`
  pattern exactly. Transcript captured during authoring (not persisted as a
  file — reproducible via the commands above).
- AC1's "not a git repo" negative half is AC1 territory too (issue's own
  Implementation requirement 2): see `TestRepoInstall_NotAGitRepo_FailsClearly`
  (PASS) — exact message `✗ cn repo install must be run inside a Git
  repository.`, no scaffolding.

### AC2 — deterministic `.cn/deps.json`

**Claim:** stable order, no timestamps, byte-identical across two runs.

**Evidence:**
- `src/go/internal/repoinstall/repoinstall_test.go::TestRun_Idempotent_ByteIdenticalArtifacts`
  — runs `Run` twice against the same `RepoRoot` + fixture index, asserts
  `bytes.Equal` on both `.cn/deps.json` and `.cn/deps.lock.json` reads. PASS.
- `writeManifest` (`repoinstall.go`) uses `json.MarshalIndent` over a typed
  `pkg.Manifest` struct (not a map) — field order is Go struct-tag order,
  not map-iteration order; no field carries a timestamp (`pkg.Manifest` has
  only `Schema`, `Profile`, `Packages`).
- `TestRun_DefaultPackageSet_AllThreeRestored` additionally asserts manifest
  package order matches `DefaultPackages`'s declared order (not resorted),
  confirming "manifest order is operator-controlled" (restore.go's own
  determinism doctrine) holds for the new writer too.

### AC3 — deterministic `.cn/deps.lock.json` (`cn.lock.v2`, exact versions, SHA-256 pins)

**Claim:** schema `cn.lock.v2`; every entry has an exact version + non-empty
sha256.

**Evidence:**
- `TestRun_LocalIndex_EndToEnd` parses the written lockfile via
  `pkg.ParseLockfile`, asserts `Schema == "cn.lock.v2"` and a non-empty
  `SHA256` field. PASS.
- The lockfile is written by `restore.GenerateLockFromIndex` — reused
  directly, not reimplemented (see §Self-check "reuse audit" below); its
  existing sort-by-name-then-version determinism is inherited for free.

### AC4 — restores `cnos.core`/`cnos.cdd`/`cnos.cds` under name-based vendor paths; checksums verified

**Claim:** `.cn/vendor/packages/<name>/cn.package.json` exists post-restore
for all three default packages; the existing SHA-256 verify path is
actually wired (not bypassed).

**Evidence:**
- `TestRun_DefaultPackageSet_AllThreeRestored` — installs the literal
  default triple from a fixture index, asserts each of
  `.cn/vendor/packages/{cnos.core,cnos.cdd,cnos.cds}/cn.package.json`
  exists and validates via `pkg.ValidatePackageManifestData`. PASS.
- `TestRun_SHAMismatchPropagates` — a deliberately-corrupted fixture (tarball
  SHA-256 does not match the index's declared SHA-256) surfaces
  `"sha256 mismatch"` on stderr and leaves no vendor directory behind. PASS.
  This is the proof the verify step is actually wired: a bypassed/faked
  verify would either not error, or would leave a partial vendor tree.
- No new SHA-256 computation exists in `internal/repoinstall` — the package
  imports neither `crypto/sha256` nor duplicates `fileSHA256`; verification
  is 100% `restore.Restore`'s existing code path (confirmed by reading
  `repoinstall.go`'s import list and grepping for `sha256` — zero hits
  outside test fixture helpers, which build fixture tarballs, not verify
  them).
- Vendor path is name-based: `pkg.VendorPath(hubPath, name)` is unchanged
  (`<hubPath>/.cn/vendor/packages/<name>`) — not modified by this cycle,
  confirmed via `git diff` showing zero hunks in `pkg.go`'s `VendorPath`.

### AC5 — idempotent (second run: no diff, no churn)

**Claim:** `cn repo install; git diff --exit-code; cn repo install; git diff
--exit-code` — both exit 0.

**Evidence:**
- `TestRepoInstall_Idempotent_NoDiffOnSecondRun` — runs the literal sequence
  from the AC's own oracle text (`git diff --exit-code` after each of two
  installs) via real `exec.Command("git", ...)` against a real
  `git init`-created repo, **plus** the stronger form: `git add` + `git
  commit` after the first install, then a second install, then asserts
  `git status --porcelain` is empty. PASS. (The stronger form is the
  meaningful proof — a bare `git diff --exit-code` before any `git add`
  passes trivially on any first run since nothing is tracked yet; the
  commit-then-rerun form is what actually exercises "no churn.")
- `TestRun_Idempotent_ByteIdenticalArtifacts` (repoinstall-level) — see AC2.
- No "already installed" special-casing was added; `Run` always regenerates
  `.cn/deps.json`/`.cn/deps.lock.json` deterministically from the same
  inputs, and `restore.Restore`'s existing version-match skip (restore.go
  lines ~169-192) prevents vendor-tree rewrites when the installed version
  already matches the lockfile pin. Idempotence is inherited from the reused
  substrate, not reimplemented.

### AC6 — `--dry-run` writes nothing; lists the planned diff

**Claim:** `git status --porcelain` empty after; stdout lists exactly
`.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore`.

**Evidence:**
- `TestRun_DryRun_WritesNothing` (repoinstall) — asserts `os.ReadDir(RepoRoot)`
  returns zero entries after a dry-run (i.e. not even a `.cn/` directory is
  created), and stdout contains all three planned-diff filenames plus the
  dispatch-mode statement. PASS.
- `TestRepoInstall_DryRun_GitStatusStaysClean` (cli, real git repo) —
  `git status --porcelain` is empty after `--dry-run`. PASS.
- Code-level guarantee: `Run`'s dry-run branch (`repoinstall.go`, the
  `if opts.DryRun` block) returns immediately after `printPlan`, before
  `applyInstall` (which owns every `os.MkdirAll`/`os.WriteFile` call in the
  non-test code path) is ever invoked.

### AC7 — `--release latest`, `--release <tag>`, `--index <path-or-url>` all work

**Claim:** all three resolution modes function independently.

**Evidence:**
- `TestRun_ReleaseLatest_ResolvesAndInstalls` — a fake GitHub API
  (`httptest.Server`) answers `/repos/{repo}/releases/latest`; a fake
  download server serves that tag's `index.json` (relative tarball URL,
  exactly as `cn build`'s local index shape declares it) + the tarball at
  the release-download path. Asserts `Result.ReleaseTag == "9.9.9"` and the
  package is restored — which is only possible if the relative URL was
  correctly rewritten to the release-download absolute URL (an
  un-rewritten fetch would 404 against the download server, since
  `restore.go`'s `fetchTarball` resolves local-looking URLs relative to
  the *local temp cache dir*, not the origin server). PASS.
- `TestRun_ReleaseTag_Pinned` — same shape with an explicit tag; the "latest"
  API endpoint is deliberately not wired to a distinct/reachable path,
  proving the pinned-tag path skips latest-resolution entirely (if it had
  called `FetchLatestRelease` regardless, the fetch would fail against the
  unregistered path). PASS.
- `TestRun_LocalIndex_EndToEnd` / `TestRepoInstall_FreshGitRepo_EndToEnd` —
  `--index <local path>`. PASS.
- `TestRun_IndexHTTPURL_RelativeRewrite` — `--index <http URL>`, relative
  tarball URL rewritten against the URL's own directory. PASS.
- Note: `--release latest`'s live-network path (a real call to
  `api.github.com`) is intentionally NOT exercised in CI — only the
  `httptest`-stubbed form is, matching the scaffold's own guidance
  ("`--release latest` test may need to stub/mock the GitHub API call
  rather than hitting it live in CI"). This is a deliberate test-design
  choice, not a gap: `binupdate.FetchLatestRelease` (the function this
  cycle reuses) already has its own live-shape unit test coverage in
  `internal/binupdate/binupdate_test.go` from prior cycles.

### AC8 — base mode never writes `.github/workflows/`, never requires a secret

**Claim:** no `.github/workflows/cnos-cds-dispatch.yml`; no
`SIGMA_WORKFLOW_PAT`/`CNOS_WORKFLOW_PAT`/`CLAUDE_CODE_OAUTH_TOKEN` reference
outside the dispatch-guard error string (which doesn't reference any of
these anyway).

**Evidence:**
- `grep -n "SIGMA_WORKFLOW_PAT\|CNOS_WORKFLOW_PAT\|CLAUDE_CODE_OAUTH_TOKEN" src/go/internal/repoinstall/repoinstall.go src/go/internal/cli/cmd_repo_install.go`
  → zero matches (run at HEAD `d783ba38`, confirmed unchanged through the
  rebase since neither file was touched by the intervening `origin/main`
  commits).
- `grep -n "\.github/workflows" src/go/internal/repoinstall/repoinstall.go src/go/internal/cli/cmd_repo_install.go`
  → the only three hits are prose/stdout strings stating that nothing is
  written there (`"Dispatch: none (base install only — no
  .github/workflows/ changes)"`, help text, package doc comment) — no code
  path constructs a path under `.github/workflows/` or calls `os.WriteFile`
  /`os.MkdirAll` against one.
- `TestRun_LocalIndex_EndToEnd` and `TestRepoInstall_FreshGitRepo_EndToEnd`
  both assert `os.Stat(filepath.Join(repoRoot, ".github"))` is
  `os.IsNotExist` after a full base install. PASS.

### AC9 — `--dispatch cds` fails explicitly, no partial `.github/workflows/` write

**Claim:** exits non-zero with a clear message; leaves zero
`.github/workflows/` files behind.

**Evidence:**
- `TestRun_DispatchCds_FailsWithNoPartialWrite` (repoinstall) — asserts the
  error contains `"#609"`, stderr carries the **exact** string
  `"✗ --dispatch cds requires generalized wake renderer support (#609)"`,
  `os.ReadDir(RepoRoot)` is empty (zero files/dirs written at all, not just
  no `.github/`), and `.github/workflows/cnos-cds-dispatch.yml` specifically
  does not exist. PASS.
- `TestRepoInstall_DispatchCds_CliWiring` (cli, full wiring through
  `RepoInstallCmd.Run`) — same assertions via the real command dispatch
  path. PASS.
- Code-level guarantee: `validateDispatch` is the first call in `Run`,
  before repo-root validation, before index resolution, before any
  `os.MkdirAll`/`os.WriteFile`. A `cds` value returns an error immediately;
  no subsequent code in `Run` executes.

### AC10 — no agent-hub scaffold (no `spec/SOUL.md`, `agent/`, `threads/`, `state/`)

**Claim:** regression guard distinguishing this installer from `cn init`'s
agent-hub scaffold (cnos#606 C4).

**Evidence:**
- `TestRun_LocalIndex_EndToEnd` and `TestRepoInstall_NoAgentHubScaffold` both
  assert `spec/SOUL.md`, `agent`, `threads`, `state` do not exist under
  `RepoRoot` after a full install. PASS.
- Code-level: `internal/repoinstall` does not import `internal/hubinit`
  (the package `cn init` uses for agent-hub scaffolding) at all — confirmed
  by reading `repoinstall.go`'s import block. The only cross-package reuse
  is `internal/hubsetup` (for the single `.gitignore`-entry helper),
  `internal/pkg`, `internal/restore`, `internal/binupdate`.

### AC11 — docs updated so the canonical base path is `cn repo install`

**Claim:** `docs/guides/INSTALL-CDS.md` names `cn repo install` as the
canonical base path; no lingering "7 manual steps" framing presented as the
only path.

**Evidence:**
- `docs/guides/INSTALL-CDS.md` (new file — did not exist on `origin/main`
  before this cycle) opens with the one-command flow
  (`curl ... | sh` then `cn repo install`) as "the canonical way to install
  Layer 1," under a two-layer framing (Layer 1 base / Layer 2 autonomous
  dispatch, opt-in) matching the issue's own §Docs requirement.
  `grep -c "cn repo install" docs/guides/INSTALL-CDS.md` → 17 occurrences,
  all consistent with the one-command path being canonical; zero
  occurrences of "7 manual steps" or "7 steps" framing.
- The design-branch draft's stale `cn.lock` filename (inconsistent with the
  actual shipped `.cn/deps.lock.json` / `cn.lock.v2` schema this cycle
  confirmed against `restore.go`) was corrected, not carried forward — see
  §Self-check "docs reconciliation" below.
- `docs/guides/README.md` was deliberately NOT modified — it carries no
  existing reference to `INSTALL-CDS.md` or the old manual-install path
  (`grep -n "INSTALL-CDS" docs/guides/README.md` → no hits), so there is no
  stale cross-reference to reconcile; adding a new nav link there is outside
  this issue's named §Docs scope (only `docs/guides/INSTALL-CDS.md` is
  named).

---

## Self-check

**Role self-check (alpha/SKILL.md §2.5 §Self-check): did this cycle's work
push ambiguity onto β? Is every claim backed by evidence in the diff?**

- Every AC above cites either an automated test (name + PASS, independently
  re-run at HEAD to produce this section, not copy-pasted from an earlier
  run) or a `grep`/`git diff`/manual-command result with the literal
  command shown. No AC claim rests on unverified prose.
- **Reuse audit (issue's own risk control 3: "don't duplicate `cn deps
  lock`/`restore` logic"):** `internal/repoinstall` imports
  `internal/restore` and calls `restore.GenerateLockFromIndex` and
  `restore.Restore` directly — zero reimplementation of lockfile generation,
  SHA-256 verification, or tar extraction. `grep -n "sha256\|Sum256"
  src/go/internal/repoinstall/repoinstall.go` → zero hits (only the test
  file, which builds fixture tarballs, imports `crypto/sha256`). Confirmed
  by reading `restore.go` in full before writing any new code, per the
  scaffold's own instruction.
- **Peer enumeration (alpha/SKILL.md §2.3):** the "multiple writers of the
  same schema" case applies to the `.gitignore` `.cn/vendor/` entry — `cn
  setup` (`hubsetup.go`) already wrote it. Peer set = {`cn setup`, `cn repo
  install`}. Resolved by exporting `hubsetup.EnsureGitignoreEntry` and
  calling it from `repoinstall.applyInstall` rather than writing a second,
  independent `.gitignore`-mutation function — one writer, two callers, not
  two writers of the same invariant. No other multi-writer schema surface
  was found: `.cn/deps.json`/`.cn/deps.lock.json` have exactly one writer
  each in the whole kernel (`repoinstall.writeManifest` and
  `restore.GenerateLockFromIndex` respectively — `cn setup` writes a
  *different* default `deps.json` shape for its own `engineer` profile, at
  a different call site, not a second writer of `cn repo install`'s `cds`
  profile manifest).
- **Harness audit (alpha/SKILL.md §2.4):** this cycle does not change any
  parser, schema-bearing type, or manifest shape — `pkg.Manifest`,
  `pkg.Lockfile`, `pkg.PackageIndex` are all read exactly as shipped, zero
  field additions. The only new "producer" is `repoinstall.writeManifest`,
  which serializes the existing `pkg.Manifest` type via
  `json.MarshalIndent` — the same mechanism `restore.GenerateLockFromIndex`
  already uses for `pkg.Lockfile`, so there is no new ad hoc encoder to
  audit against a schema. Grepping `*.sh` and `.github/workflows/*.yml` for
  `deps.json` literal writes turns up two **pre-existing** hand-rolled
  `cn.deps.v1` writers: `.github/workflows/build.yml`'s Tier-2 CI test-hub
  setup step (heredoc, `profile: "engineer"`, pins `cnos.core`/`cnos.kata`/
  `cnos.cdd.kata` for CI's own package-testing purposes) and
  `scripts/kata/06-install.sh` (same shape, kata-harness purposes). Neither
  is a consumer of `cn repo install`/`cn setup`'s specific default-package
  logic and neither needed a change here — the `cn.deps.v1` schema itself is
  byte-for-byte unchanged by this cycle (no field added/removed/renamed), so
  the harness-audit trigger condition (alpha/SKILL.md §2.4: "when the branch
  changes a parser, schema-bearing type, manifest shape, or runtime
  contract") does not fire. Noted here because an earlier draft of this
  section claimed "no independent writers exist," which the grep above
  disproves — corrected before this AC evidence was finalized, per the
  intra-doc-repetition/closure-overclaim discipline (alpha/SKILL.md §2.3):
  `pkg.Manifest` already has multiple legitimate writers today (`hubsetup`,
  these two harnesses, and now `repoinstall`), each producing its own
  package list for its own purpose; that plurality is pre-existing and
  outside this cycle's scope to consolidate.
- **Docs reconciliation:** the design-branch draft
  (`origin/claude/cds-install-guide-2ka54j:docs/guides/INSTALL-CDS.md`) used
  a stale `cn.lock` filename inconsistent with the actual shipped
  `.cn/deps.lock.json` (`cn.lock.v2`) schema this cycle confirmed by reading
  `restore.go` directly. Rather than land that inconsistency, the new
  `docs/guides/INSTALL-CDS.md` in this cycle's diff uses the confirmed-real
  path throughout, and drops the draft's Track-B (`workflow_dispatch`
  GitHub-UI installer) section entirely, since `templates/cnos-install.yml`
  and the installer workflow do not exist anywhere on `origin/main` (confirmed
  via `find . -iname cnos-install.yml` → no hits) — that is #611's
  (bootstrap delegation) scope, not #608's. Documenting a workflow file that
  does not exist would have been a false claim.
- **CI-schema harness (post-hoc finding, folded into this section rather
  than treated as a separate late AC):** the `cn cdd verify` I6 gate (which
  gates merges via `.github/workflows/build.yml`'s `cdd-artifact-check`
  job) enforces bare `## Gap` / `## Skills` (or `## Mode`) / `## ACs` (or
  `## AC Coverage`) / `## CDD Trace` / a `##`-line containing "self-check"
  or "debt" section headers on `self-coherence.md` — not the `§`-prefixed
  style (`## §Gap` etc.) this cycle initially used (and which a prior cycle,
  #600, also used without ever actually satisfying the checker — #600's
  mismatch was merely masked by its "triadic" classification's lenient
  in-progress warning, not actually schema-conformant). Because cycle #608
  is classified "small-change" by `classifyCycleType` (self-coherence.md
  present, no `beta-review.md` yet — true for any α-authored cycle before β
  responds) it hits `checkSmallChangeArtifacts`, which hardcodes
  `forUnreleased=false`, turning the same header mismatch into a hard CI
  failure rather than a warning. Fixed by switching this file's headers to
  the bare form (commit `9fbd8fa1`) rather than requesting a
  `.cdd/exceptions.yml` entry — the bare form is what a genuinely
  fully-green historical cycle (#593) actually uses, so this is a
  conformance fix, not a workaround. **Pattern observed, not
  triaged/fixed at the tool level** (per alpha/SKILL.md §2.8 "voice: factual
  observations only, triage is γ's job"): `checkSmallChangeArtifacts`'s
  hardcoded `forUnreleased=false` (ledger.go line 497) appears to make it
  structurally impossible for a small-change-classified, still-in-progress
  self-coherence.md (i.e., any α authoring §2.5's incremental one-section-
  per-commit sequence, on any cycle without a pre-existing beta-review.md)
  to pass I6 mid-authorship even when using the *correct* bare header
  convention throughout, since the same "missing sections" logic fires
  hard-fail (not warn) until every required section exists in one commit.
  This cycle's own commits before the final section landed are a live
  example (`eb269381`, `a3f9b86e` both show `❌ Cycle artifact verification
  FAILED` in Actions). This is surfaced as a pattern for γ/δ, not
  self-triaged here.
- No ambiguity was pushed onto β regarding the implementation-contract axes
  (all 7 were pinned by γ at dispatch; none were relaxed or reinterpreted —
  see §Skills above).
