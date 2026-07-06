# beta-review — cnos#608 — R0

## Verdict: **iterate**

Base SHA at review time: `origin/main` = `80778d688e04e61c66d38f2bd5962fafb0729e95` (re-fetched and confirmed current; matches α's rebased base — no drift since α's last rebase). Cycle-branch head reviewed: `226181f4c5a70fd6ab632e23f37ddc9cb5bf7925`.

## Independent verification performed (not a re-read of α's claims)

- Fresh clone (`git clone` into a scratch dir, not reusing any prior worktree), `git fetch origin cycle/608 && git switch -c cycle/608-review origin/cycle/608` — clean checkout, no leftover build artifacts.
- `cd src/go && go build ./...` — clean, exit 0.
- `go vet ./...` — clean, exit 0.
- `go test ./...` — all packages pass, including `internal/repoinstall` and `internal/cli` (0 failures). Matches α's claimed 304-test count in shape (all green).
- Read `internal/repoinstall/repoinstall.go`, `internal/cli/cmd_repo_install.go`, `src/go/cmd/cn/main.go`, `internal/hubsetup/hubsetup.go` diff, and the landed `docs/development/design/cn-repo-install-MOCKS.md` / `docs/guides/INSTALL-CDS.md` in full — not just the self-coherence report's excerpts.
- Confirmed no new `go.mod`/`go.sum` entries (runtime-dependency contract axis holds).
- Confirmed `.github/workflows/*` is untouched in the diff (`git diff --stat origin/main...HEAD` — no workflow file appears).
- Confirmed no SKILL.md files are touched (frontmatter validator not applicable to this cycle; pre-merge gate row 3 collapse for that specific validator is legitimate).
- Confirmed CI is green at the reviewed HEAD via `gh run list --repo usurobor/cnos --branch cycle/608` (`Build` → `success` at `226181f4`).
- Built the actual `cn` binary (`go build -o cn ./cmd/cn`) and drove `cn repo install` directly (not through the test harness) to independently re-derive AC1's behavior — this is where a real, reproducible bug surfaced (see Findings).

## AC-by-AC walk (code-first, per Rule 6)

- **AC1 (kernel command, fresh-repo, no prior state)**: dispatch shape confirmed (`RepoInstallCmd.Spec().Name == "repo-install"`, `Source: SourceKernel`, `NeedsHub: false`, registered in `main.go`). The "runs from a fresh git repo" half is confirmed. **The negative half of AC1 — "fails with a clear error if not at a git repo root (does not silently walk up or scaffold)" — is FALSE at the real dispatch boundary.** See Finding 1 below. This is CONFIRMED, not asserted from the test suite: I reproduced it against the built binary.
- **AC2** (deterministic `.cn/deps.json`): confirmed — `writeManifest` marshals a typed struct via `json.MarshalIndent`, no timestamp field, stable order. `TestRun_Idempotent_ByteIdenticalArtifacts` is a real byte-comparison test, not tautological.
- **AC3** (`cn.lock.v2`, exact versions, SHA-256 pins): confirmed — lockfile generation is 100% delegated to `restore.GenerateLockFromIndex` (reused, not reimplemented); `grep -n sha256` in `internal/repoinstall/repoinstall.go` (non-test) returns zero hits, confirming no parallel SHA computation was introduced.
- **AC4** (restores default triple, checksums verified): confirmed — `TestRun_SHAMismatchPropagates` uses a deliberately-corrupted fixture and asserts the *existing* verify path actually fires (a real negative test, not a happy-path-only check).
- **AC5** (idempotent): confirmed — `TestRepoInstall_Idempotent_NoDiffOnSecondRun` uses the stronger commit-then-rerun form (a bare pre-commit `git diff --exit-code` would pass trivially; this test correctly avoids that trap).
- **AC6** (`--dry-run` writes nothing): confirmed — code-level guarantee (`Run`'s `DryRun` branch returns before `applyInstall`, the sole owner of every write, is ever called) plus test-level `os.ReadDir` emptiness assertion.
- **AC7** (`--release latest` / `--release <tag>` / `--index <path-or-url>`): confirmed — real `httptest.Server`-backed fixtures for the GitHub API and release-asset download path, including the relative-URL rewrite case. The live-network path is deliberately not exercised in CI, which matches the scaffold's own guidance and is not a gap.
- **AC8** (no `.github/workflows/`, no secret requirement): confirmed by direct grep of the new source (zero hits for the named secret names, zero non-prose `.github/workflows` path construction) and by test assertion (`os.Stat(".github")` is `IsNotExist` post-install).
- **AC9** (`--dispatch cds` fails explicitly, no partial write): confirmed — `validateDispatch` is the first statement in `Run`, before any I/O; test asserts exact stderr string, `#609` reference, and zero files written (not just "no `.github/`" — the whole `RepoRoot` is checked empty).
- **AC10** (no agent-hub scaffold): confirmed — `internal/repoinstall` does not import `internal/hubinit` at all (checked the import block directly, not just trusting the self-coherence claim).
- **AC11** (docs updated): confirmed — `docs/guides/INSTALL-CDS.md` opens with the one-command flow as canonical; `grep -c "cn repo install"` → 17 hits; zero "7 manual steps" framing. The design-branch draft's stale `cn.lock` filename was correctly *not* carried forward (the shipped doc uses the real `.cn/deps.lock.json` / `cn.lock.v2` path) — a genuine reconciliation, not a doc/code drift.

## Implementation-contract conformance (Rule 7)

All seven pinned axes checked against the diff directly:

- **Language**: Go only. Confirmed.
- **CLI integration target**: kernel `cn` subcommand via noun-verb resolution (`ResolveCommand` builds `"repo"+"-"+"install"`), not a separate binary, not shelling out to `cn` internally. Confirmed by reading `dispatch.go` and `cmd_repo_install.go`.
- **Package scoping**: new `internal/repoinstall/` package + `cli/cmd_repo_install.go`, reusing `restore/`, `pkg/`, `binupdate/`, `hubsetup/`. No new top-level directory. Confirmed.
- **Existing-binary disposition (coexist)**: `cmd_init.go`, `cmd_setup.go`, `cmd_deps.go` are untouched in the diff (`git diff --stat` shows zero hunks in any of them). `hubsetup.go`'s only change is a pure export rename (`ensureGitignoreEntry` → `EnsureGitignoreEntry`), same body, same call site in `hubsetup.Run` — behavior-preserving. Confirmed.
- **Runtime dependencies**: no `go.mod`/`go.sum` changes. Confirmed.
- **JSON/wire contract preservation**: `pkg.Manifest`/`pkg.Lockfile`/`pkg.PackageIndex` are read/written via their existing typed shapes; `git diff` on `pkg.go` shows zero hunks. Vendor path is name-based (`.cn/vendor/packages/<name>/`), confirmed unchanged (`pkg.VendorPath` untouched). Confirmed.
- **Backward-compat invariant**: no shared code path (`restore.go`, `pkg.go`, `binupdate.go`) is modified. Confirmed.

Implementation-contract conformance: **PASS** on all seven axes. The blocking finding below is a behavioral/AC-oracle bug, not a contract-axis drift.

## Findings

### Finding 1 — AC1's "does not silently walk up" invariant is violated at the real dispatch boundary (BLOCKING)

**Severity:** major (confirmed, reproducible, silently writes into the wrong directory with exit code 0).

**What's wrong:** `RepoInstallCmd.Run` (`src/go/internal/cli/cmd_repo_install.go`, lines ~107–115) does:

```go
repoRoot := inv.HubPath
if repoRoot == "" {
    root, gerr := gitRepoRoot(ctx)
    ...
}
```

`inv.HubPath` is populated by `main.go`'s `discoverHub()`, which walks **upward from cwd looking for any `.cn/` directory**, with no bound at the current git repository's root and no verification that the found directory is even a git repository. Because `RepoInstallCmd.Spec().NeedsHub == false`, `discoverHub()` still runs unconditionally before dispatch (per `main.go`'s existing, unconditional call to `discoverHub()`).

If `inv.HubPath` comes back non-empty (an ancestor `.cn/` exists anywhere above cwd, for any reason — an unrelated project, a stale hub, a parent workspace), `cmd_repo_install.go` **skips `gitRepoRoot` entirely** and installs into that ancestor directory, without ever checking that cwd (or the ancestor) is a git repository at all.

This directly contradicts the issue's own Mock A1 / AC1 text: *"Fails with a clear error if not at a git repo root (does not silently walk up or scaffold)."* The mock_parity block in `self-coherence.md` marks A1 as `match`, citing `TestRepoInstall_NotAGitRepo_FailsClearly` as evidence — but that test constructs `Invocation{HubPath: ""}` directly (see `runRepoInstall` test helper, `cmd_repo_install_test.go` line ~91), bypassing `main.go`'s real `discoverHub()` entirely. **No test in the diff exercises the `inv.HubPath != ""` branch at all.** The AC1 "match" claim is true only for the code path the tests happen to hit, not for the actual production dispatch path.

**Reproduction (against the built binary, this branch, HEAD `226181f4`):**

```
mkdir -p /tmp/repro/outer-hub/nested/not-a-git-repo
mkdir /tmp/repro/outer-hub/.cn                    # unrelated ancestor "hub"
cd /tmp/repro/outer-hub/nested/not-a-git-repo      # NOT a git repo at all
git rev-parse --show-toplevel                      # fatal: not a git repository...
./cn repo install --dry-run
```

Output:
```
→ cn repo install (dry-run) — no files will be written
✓ Git repository root: /tmp/repro/outer-hub
✓ Resolved cnos release: 3.82.1
✗ package(s) not found in index: cnos.core@3.82.1, ...
```

Exit code 0 up through the (unrelated) index-resolution failure — i.e., the command proceeds past the git-repo check entirely, silently treating `/tmp/repro/outer-hub` (never verified to be a git repository) as the install root. With a real reachable index, this would write `.cn/deps.json`, `.cn/deps.lock.json`, `.cn/vendor/`, and mutate `.gitignore` under `/tmp/repro/outer-hub` — not under the directory the operator is actually standing in, and not gated by any git-repo-root check at all.

This is a materially different risk class from existing `NeedsHub: true` commands' use of the same `discoverHub()` walk: for those commands, `.cn/` was created by `cn init`/`cn setup` *at* the git root the operator is already working in, so the walk-up target and the git root coincide by construction. `cn repo install`'s entire premise is running in a repo that has *no* `.cn/` yet — exactly the scenario where an unrelated ancestor `.cn/` (a different project, a parent workspace, a monorepo-of-repos layout) can exist and silently hijack the install root.

**What would fix it:** `RepoInstallCmd.Run` should always resolve the repo root via `gitRepoRoot(ctx)` (git's own `rev-parse --show-toplevel`), independent of `inv.HubPath`. `inv.HubPath` may still be consulted as an idempotent-rerun optimization, but only after confirming it equals the git-root value — never as a substitute for the git-repo check. The existing negative test (`TestRepoInstall_NotAGitRepo_FailsClearly`) should also be extended (or a new test added) that goes through `main.go`'s actual `discoverHub()` path (or an equivalent harness that doesn't hand-construct `Invocation{HubPath: ""}`), so this class of drift is caught by CI rather than requiring a manual binary repro.

### Non-blocking observations (α already disclosed most of these; independently reviewed and agreed non-blocking)

1. Design doc's Mock B table header reads "Invariants B1–B3" but lists four rows (B1–B4) — a pre-existing inconsistency in the verbatim-landed source file, correctly not silently "fixed" by α and correctly disclosed. Non-blocking.
2. γ's scaffold miscounted "9 rows total" when the named set (A1–A5, B1–B4, E1) is actually 10 — α used the correct 10-row set and flagged the arithmetic error rather than silently dropping a row. Non-blocking, correctly handled.
3. `cn --version` referenced in the new install doc currently exits non-zero — pre-existing, out of scope for #608 (Mock F1, later sub), correctly disclosed rather than silently reproduced without comment. Non-blocking.
4. `--index <URL>` + explicit `--release <tag>` combination has no dedicated test (only each sub-case individually) — low risk, same code path either way, correctly disclosed as narrow test-coverage debt. Non-blocking.
5. `checkSmallChangeArtifacts`'s hard-fail-vs-warn asymmetry between small-change and triadic in-progress artifacts (CI tooling friction, not this cycle's diff) — correctly surfaced as a γ/δ tooling observation, not self-triaged. Non-blocking for this cycle.
6. No lockfile-level concurrency guard for simultaneous `cn repo install` runs — inherited from `restore.Restore`'s pre-existing lack of a lock primitive, not new. Non-blocking.

## Verdict rationale

Build/vet/test are clean from a fresh checkout; the seven implementation-contract axes all hold; ten of eleven ACs are independently confirmed against the actual code (not the self-coherence report). AC1's negative case, however, fails at the real dispatch boundary in a way the shipped test suite cannot detect (the test harness bypasses `main.go`'s hub-discovery walk entirely) — I reproduced this independently against the built binary. Per Role Rule 6 ("anchor oracle evidence on code, not doc" — extended here to "not a test that stubs around the real dispatch path") and Role Rule 4 ("do not merge with unresolved findings"), this is an iterate-triggering finding: it is squarely inside AC1's own explicit negative-case text, it is silent (exit 0, no error) and reproducible, and it can cause `cn repo install` to write into a directory other than the one the operator invoked it in.

**Iterate.** α should fix `cmd_repo_install.go`'s repo-root resolution to always verify the actual git root (not trust an unbounded ancestor `.cn/` walk) and add a test that exercises the real `main.go`-equivalent dispatch path for this case. Once that lands, re-review should be a narrow re-check of AC1 plus a regression pass on AC2–AC11 (none of which are expected to be affected by this fix).

---

# beta-review — cnos#608 — R1 (re-review)

## Verdict: **converge**

`origin/main` re-fetched synchronously for this pass: `git rev-parse origin/main` = `a8a3d8de1b69e0b9f747508d1e5894125514e82d`. This is one commit ahead of the `c08a7483` tip α rebased onto for the R1 fix (`sigma: agent-admin heartbeat wake at cnos, 2026-07-06T18:03:31Z`, `.cn-sigma/logs/20260706.md` only — disjoint from every path this cycle's diff touches; `git diff --stat c08a7483..a8a3d8de` confirms a single-file, unrelated log append). No re-rebase needed for this review pass; nothing in the intervening commit affects any AC or the diff base.

Cycle-branch head reviewed: `712f6b8743d6ebe420fc5e617e8b1d0a68d38c10` (fresh checkout via `git switch -c cycle/608-r1review origin/cycle/608`, no reused worktree).

This is an independent re-review, not a rubber-stamp of α's R1 claims — every item below was re-derived from the code and a freshly built binary, not read off `self-coherence.md`'s §R1 prose.

## 1. R0 blocking finding — independently re-verified fixed

Read `src/go/internal/cli/cmd_repo_install.go` at current HEAD directly (not diffed against α's narrative). Confirmed:

- `inv.HubPath` does not appear anywhere in the file (`grep -n HubPath src/go/internal/cli/cmd_repo_install.go` → zero hits). The R0-flagged branch (`repoRoot := inv.HubPath; if repoRoot == "" { ... }`) is gone entirely, not merely reordered.
- `RepoInstallCmd.Run` now calls `gitRepoRoot(ctx)` unconditionally as its first substantive step after arg parsing, and fails with the exact AC1 message (`✗ cn repo install must be run inside a Git repository.`) whenever that call errors — with no conditional path around it.

**Repro re-run against a freshly built binary from this exact branch HEAD** (not trusting α's transcript):

```
cd src/go && go build -o /tmp/cn_r1 ./cmd/cn
mkdir -p /tmp/repro2/outer-hub/nested/not-a-git-repo
mkdir /tmp/repro2/outer-hub/.cn
cd /tmp/repro2/outer-hub/nested/not-a-git-repo
/tmp/cn_r1 repo install --dry-run
```

Result: `✗ cn repo install must be run inside a Git repository.`, exit code 1, and `/tmp/repro2/outer-hub/.cn` remains empty (`ls -la` → only `.`/`..`). This is the exact scenario that silently proceeded to index resolution at R0. **Finding 1 is fixed**, confirmed independently.

## 2. New tests — confirmed non-vacuous, real dispatch path

Read `TestRepoInstall_RealDispatch_AncestorHubOutsideGitRepo_FailsClearly` and `TestRepoInstall_RealDispatch_AncestorHubWithNestedGitRepo_UsesInnerRoot` in `src/go/internal/cli/cmd_repo_install_test.go` in full. Both call `buildCnBinary(t)` (a real `go build -o <tmp>/cn ./cmd/cn` subprocess, not a stub) and drive the resulting binary via `exec.Command` with `cmd.Dir` set to a constructed cwd — this exercises `main.go`'s actual `discoverHub()` and dispatch path, not a hand-constructed `Invocation{HubPath: ...}` (which is exactly the bypass that let the R0 bug slip through the original `TestRepoInstall_NotAGitRepo_FailsClearly`).

**Reverted the fix locally and confirmed both new tests fail** (per this task's suggested verification, done rather than reasoned-about only): checked out `git show 5fd47490:.../cmd_repo_install.go` (the pre-fix content) into an isolated `git worktree` (not the reviewed branch itself), ran the two new tests against it:

```
--- FAIL: TestRepoInstall_RealDispatch_AncestorHubOutsideGitRepo_FailsClearly
    stderr = "✗ resolve latest cnos release: http status 403"
    (proceeded past the git-repo check; failed on live release resolution
    instead of the git-repo check — same bug class as β's R0 repro, though
    the exact downstream error differs from α's transcript because this
    sandbox's network path hit GitHub's rate limit rather than the
    package-index-lookup failure α saw; the load-bearing fact — the git-repo
    check is skipped entirely — is confirmed either way)
--- FAIL: TestRepoInstall_RealDispatch_AncestorHubWithNestedGitRepo_UsesInnerRoot
    stdout reports the ancestor as "Git repository root"; .cn/deps.json is
    written under the ancestor .cn/, not the inner git repo — silent
    wrong-directory install, exit 0
```

Both tests genuinely fail against the pre-fix code and pass against the fix (re-ran `go test ./internal/cli/... -run 'TestRepoInstall_RealDispatch_.*'` against the actual reviewed HEAD: both PASS). The tests are not tautological — they detect the exact regression class they were added for. Worktree torn down after (`git worktree remove --force`), reviewed branch untouched.

One incidental observation: `TestRepoInstall_RealDispatch_AncestorHubOutsideGitRepo_FailsClearly` passes no `--index`/`--release`, so if the git-repo-root check were ever accidentally removed or reordered again, this test would attempt a live `api.github.com` call rather than failing fast/hermetically — as demonstrated by the 403 above. Against the *current, fixed* code this never fires (the git-repo check runs before any network access), so this is not a regression risk today, only a latent hermeticity gap in the test's own defense-in-depth should the fix regress a second time. Non-blocking, noted for completeness rather than requiring a change.

## 3. No regression — fresh build/vet/test/race, this branch

Run from a clean checkout of `cycle/608-r1review` (tracking `origin/cycle/608` HEAD `712f6b87`):

- `go build ./...` — clean, exit 0.
- `go vet ./...` — clean, exit 0.
- `go test ./...` — all 15 packages `ok`, 0 failures (307 `--- PASS`, 0 `--- FAIL` via `go test ./... -v | grep -c`). Matches α's claimed count exactly (304 at R0 + 3 new this round).
- `go test -race -count=1 ./...` — clean, all 15 packages pass, no race reports.

## 4. AC1–AC11 re-walk

- **AC1**: now fully confirmed both halves — positive (fresh git repo, no prior state) unaffected by the fix, still passes via `TestRepoInstall_FreshGitRepo_EndToEnd`; negative (not-a-git-repo, including the ancestor-hub-hijack shape) now genuinely enforced at the real dispatch boundary, confirmed above. **Closed.**
- **AC2–AC11**: no code path other than repo-root resolution changed this round (`git diff origin/main...origin/cycle/608 -- src/go/internal/cli/cmd_repo_install.go` is the only production-code hunk touching dispatch logic; `internal/repoinstall/repoinstall_test.go` gained one additional test, `TestRun_IndexHTTPURL_WithExplicitReleaseTag`, closing R0's non-blocking observation 4 — verified it exercises the remote-URL + explicit-release-tag combination, not a restatement of an existing case). Spot-re-ran the R0 evidence set (idempotence, dry-run empty-write, SHA-mismatch propagation, dispatch-cds guard, no-agent-hub-scaffold) — all still pass, unchanged from R0's independent confirmation. No regression across AC2–AC11.

## 5. Rebase/force-push drift check

`git log --oneline origin/main..origin/cycle/608` shows a clean, complete line: γ scaffold → α's design-doc landing → domain/cli/registration commits → α's incremental self-coherence sections → β's R0 review commit (`0e67adf5`) → α's R1 fix commit (`aef74f3e`) → α's R1 self-coherence appendix (`712f6b87`, current HEAD). No commits are missing from the history; the SHA-remapping table α added in `self-coherence.md` §R1 correctly accounts for the rebase (pre-rebase SHAs `226181f4`/`2a3699b7`/`49d0cbb2` → post-rebase `1346f894`/`68ab7edf`/`0e67adf5`, content-identical).

Confirmed present and non-empty on `origin/cycle/608` at current HEAD:
- `.cdd/unreleased/608/CLAIM-REQUEST.yml` (12 lines) — γ's pre-flight claim record, intact.
- `.cdd/unreleased/608/gamma-scaffold.md` (189 lines) — intact.
- `.cdd/unreleased/608/self-coherence.md` (1074 lines) — both the original §Gap/§Skills/§ACs/§Self-check/§Debt/§CDD Trace/§Review-readiness sections (R0) and the new §R1 appendix are present; nothing was rewritten in place (verified the R0 prose is byte-identical to the version β read at R0, only new content was appended).
- `.cdd/unreleased/608/beta-review.md` — my own R0 section (this file, above the separator) is untouched by α; this R1 section is a pure append, per Role Rule "never restart completed sections."

No artifact loss or drift from the force-push.

## Verdict rationale

α's fix is a net deletion of the buggy conditional, not new conditional complexity, and closes exactly the bug class β's R0 finding named (verified independently against a freshly built binary from this exact branch HEAD, not merely re-reading α's transcript). The two new tests exercise the real dispatch path (build + subprocess) and were confirmed — by actually reverting the fix in an isolated worktree and re-running them — to genuinely fail against the pre-fix code, so they are not vacuous regression coverage. Build/vet/test/race are all clean from a fresh checkout. AC1's negative half is now closed at the real dispatch boundary; AC2–AC11 show no regression. `origin/main` advanced by one unrelated commit since α's rebase base; it does not touch any path this cycle's diff touches and does not require another rebase before merge. All required `.cdd/unreleased/608/` artifacts are present and intact after the force-push rebase.

**Converge.** This closes the review loop; the cycle proceeds to closeout (β merge + β close-out, per `beta/SKILL.md`).
