# alpha close-out — cycle #608

**Issue:** [cnos#608](https://github.com/usurobor/cnos/issues/608) — `cn repo install` base installer (Sub 1 / PR 1 of the #607 wave). **Verdict at merge:** converge (β R1), one iteration.

---

## What α built

A new kernel command, `cn repo install`, registered through the existing noun-verb resolver (`RepoInstallCmd`, `Spec().Name == "repo-install"`, `Source: SourceKernel`, `NeedsHub: false`) so it can run in a plain `git init`-only repo before any `.cn/` hub exists. The domain logic lives in a new `internal/repoinstall/` package; the CLI surface is a thin `internal/cli/cmd_repo_install.go` wrapper. Both reuse the existing substrate directly — `restore.GenerateLockFromIndex`, `restore.Restore`, the `pkg.Manifest`/`pkg.Lockfile`/`pkg.PackageIndex` types, `binupdate`'s release-resolution shape — with zero reimplementation of lock generation, SHA-256 verification, or tar extraction. `docs/development/design/cn-repo-install-MOCKS.md` was landed verbatim from the operator-reviewed design branch (not re-authored), and `docs/guides/INSTALL-CDS.md` was rewritten so `cn repo install` is the canonical Layer-1 onboarding path.

All eleven ACs closed with PASS-backed evidence, and the `mock_parity` block converged at 10/10 matched (A1–A5, B1–B4, E1), 0 missed. Final test count: 307 passing, 0 failing (`go test ./...`), plus clean `go vet ./...` and `go test -race ./...` on the touched packages — up from 304 at R0's signal, the +3 being the R1 fix-round's regression tests (two real-dispatch-path tests, one flag-combination test closing a narrow R0 debt item).

## The one iteration, and what I learned from it

β's R0 finding was real, not a nitpick: `RepoInstallCmd.Run` trusted `inv.HubPath` (populated by `main.go`'s unbounded ancestor `.cn/` walk) as the install root whenever it was non-empty, skipping the actual git-repo-root check entirely. Against a real binary, standing in a directory that is not a git repo at all but sits under some unrelated ancestor `.cn/`, the command would silently proceed — exactly the negative case AC1's own text calls out ("does not silently walk up") and exactly the failure mode `cn repo install`'s whole premise (installing where no `.cn/` yet exists) makes newly reachable, unlike the existing `NeedsHub: true` commands where the walk-up target and the git root coincide by construction.

I wrote a negative test for this AC (`TestRepoInstall_NotAGitRepo_FailsClearly`) at R0 and it passed — because it hand-constructed `Invocation{HubPath: ""}` rather than driving the real dispatch path through `main.go`'s `discoverHub()`. The test satisfied the AC's literal text while stubbing around the exact boundary the AC concerned. That is the lesson worth carrying forward: a negative test built by constructing the harness input directly, instead of building the binary and driving it as a subprocess, can pass while the real invariant is false. β's Role Rule 6 (code-first, build-the-binary) is precisely the check that catches this, and it did.

The R1 fix was a net deletion, not new complexity: `gitRepoRoot(ctx)` is now called unconditionally, first, before any `inv.HubPath` read; the alternative (trust `inv.HubPath` after confirming it equals the git root) was considered and rejected as an unnecessary second code path. Two new tests build the real `cn` binary and drive it as a subprocess against β's exact repro shape (ancestor `.cn/` outside any git repo) plus a stronger control (ancestor `.cn/` alongside a correctly-nested inner git repo), and I verified both genuinely fail against the pre-fix code before confirming they pass against the fix — the same "revert and re-run" discipline β used to verify my fix, applied to my own regression tests before handing them back.

## Debt disclosed

Seven items were named explicitly in `self-coherence.md` §Debt, none blocking any AC: (1) `cn --version` exits non-zero — pre-existing, now cross-referenced against #612 since the new install doc references it; (2) `cn cdd verify`'s small-change classifier hard-fails an in-progress `self-coherence.md` rather than warning — tooling friction forcing a mid-cycle header-format change, cross-referenced against #577; (3) `--index <URL>` + explicit `--release <tag>` untested — closed within the cycle at R1; (4) `docs/guides/README.md` nav link not added — correctly out of scope per the issue's own text (deferred to #611); (5) no lockfile-level concurrency guard — inherited from `restore.Restore`, not new; (6) no Windows support — pre-existing kernel-wide scope boundary; (7) transient CI red on `cdd-artifact-check` during incremental self-coherence authoring — expected, resolved once all sections landed, same root cause as (2).

One accounting note on my own artifact, not β's or γ's: γ's scaffold text said "9 rows total" for the mock_parity set while separately naming a 10-item set (A1–A5, B1–B4, E1). I used the correct 10-row count and flagged the scaffold's arithmetic slip in `self-coherence.md` rather than silently dropping a row to match "9" or silently correcting γ's prose without comment — β independently confirmed the 10-row accounting was right.

## Artifact list

`src/go/internal/repoinstall/repoinstall.go` + `repoinstall_test.go` (new, 22+1 tests), `src/go/internal/cli/cmd_repo_install.go` + `cmd_repo_install_test.go` (new, 9+2 tests), `src/go/cmd/cn/main.go` (+1 registration line), `src/go/internal/hubsetup/hubsetup.go` (export rename only, `ensureGitignoreEntry` → `EnsureGitignoreEntry`, no behavior change), `docs/development/design/cn-repo-install-MOCKS.md` (new, landed verbatim), `docs/guides/INSTALL-CDS.md` (new), `.cdd/unreleased/608/self-coherence.md` (this cycle's primary artifact, R0 + R1 sections).

## Final AC status

AC1–AC11: all PASS, all confirmed independently by β against code and a built binary, not merely read off this artifact. AC1's negative half was the one AC that needed a round-2 fix; AC2–AC11 were unaffected by that fix and re-confirmed unchanged at R1.
