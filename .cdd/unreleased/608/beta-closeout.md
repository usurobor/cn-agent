# beta close-out — cycle #608

**Issue:** [cnos#608](https://github.com/usurobor/cnos/issues/608) — `cn repo install` base installer (Sub 1 / PR 1 of the #607 wave). **Verdict:** converge after one iteration (R0 iterate → R1 fix → R1 converge).

---

## Review Summary

Two review passes, both from fresh clones (never a reused worktree), both re-fetching `origin/main` synchronously before computing the diff base rather than trusting a session-start snapshot. R0 reviewed head `226181f4` against base `80778d68`; R1 reviewed head `712f6b87` against base `a8a3d8de` (one unrelated log-append commit ahead of α's R1 rebase point — confirmed disjoint from this cycle's diff, no re-rebase needed). One blocking finding across both rounds; everything else was independent confirmation, not correction.

## Implementation Assessment

All seven implementation-contract axes held throughout: Go only, kernel `cn` subcommand via the existing noun-verb resolver (not a separate binary, not a shell-out), new package under `internal/repoinstall/` reusing `restore/`, `pkg/`, `binupdate/`, `hubsetup/` directly, no `go.mod`/`go.sum` changes, `cn.deps.v1`/`cn.lock.v2`/package-index schemas byte-for-byte unchanged, `cn init`/`setup`/`deps` behavior untouched. The blocking finding was a behavioral bug at the AC1 negative-case boundary, not a contract-axis drift — worth separating cleanly, since Rule 7 (contract conformance) and Rule 4 (no merge with unresolved findings) are different gates and this cycle only ever tripped the second one.

## Technical Review

**Why the R0 finding was real, and how I caught it.** α's mock_parity block and self-coherence report both marked AC1 as `match`, citing `TestRepoInstall_NotAGitRepo_FailsClearly` as evidence. That test passed — but it hand-constructed `Invocation{HubPath: ""}`, bypassing `main.go`'s real `discoverHub()` walk entirely. I didn't stop at "the cited test passes." I built the actual `cn` binary and drove `cn repo install --dry-run` from a directory that was not inside any git repository but sat under an unrelated ancestor `.cn/` — the exact shape `discoverHub()`'s unbounded upward walk can hand back as `inv.HubPath`. The command proceeded past the point AC1's own text says it must fail ("does not silently walk up"), resolved the unrelated ancestor as the install root, and only stopped on an unrelated downstream error (package-index resolution). Exit code 0 up to that point, no error, silent — precisely the risk class AC1 was written to prevent, and precisely the class a report-level or test-suite-level read would have missed, since the report was true about the test and false about the invariant.

This is the concrete case for Rule 6 ("anchor oracle evidence on code, not doc") extended one level further: the doc/report can be accurate about what the test suite verifies while the test suite itself stubs around the real dispatch boundary. Building the binary and driving it as a subprocess is what surfaced the gap between "the cited test passes" and "the AC's own invariant holds."

**R1 re-review discipline.** I did not take α's "fixed, verified" claim on faith. I re-read `cmd_repo_install.go` at the new HEAD directly (confirmed `inv.HubPath` no longer appears anywhere in the file, not merely reordered), rebuilt the binary fresh, and re-ran my own R0 repro commands verbatim against it. I also reverted α's fix in an isolated worktree (not the reviewed branch) and re-ran α's two new regression tests against the pre-fix code to confirm they actually fail there — both did, for the exact bug class, before confirming both pass against the fix. That is the same standard I held α's original AC1 test to; applying it to α's fix-round evidence too, rather than trusting a second-round claim more readily than a first-round one, is what the converge verdict rests on.

## Process Observations

- The R0→iterate→R1→converge loop cost one round for one confirmed bug — not evidence of a scaffold or process failure by itself. γ's own closeout traces the gap to γ's scaffold's AC-oracle table naming only AC1's positive check, and to α's negative test stubbing around the real dispatch boundary; both are true, and the review loop caught the compound gap at normal cost, exactly as designed.
- α's fix-round discipline matched the standard I set at R0: net deletion over added complexity, regression tests that build-and-drive the real binary rather than a synthetic harness, and an explicit revert-and-confirm-failure step before claiming the tests catch the regression class. I did not have to ask for any of that — it was already in the R1 self-coherence appendix.
- Confidence in the R1 converge verdict is high: the fix is a strict simplification (one conditional deleted, not added), AC2–AC11 touch no changed code path this round, and every claim in this artifact was independently re-derived against a freshly built binary or a fresh `git test`/`vet`/`race` run, not read off α's transcript.

## For a future reviewer of a similar cycle

When a negative AC clause ("fails if X", "does not silently do Y") is tested via a hand-constructed harness input rather than a real subprocess dispatch, treat that as an open question, not closure evidence — the harness construction itself may be exactly the code path the invariant depends on. Build the binary and drive it as a subprocess whenever the AC concerns behavior at the actual dispatch/entry boundary (argument parsing, hub discovery, cwd resolution) rather than pure domain logic; a fixture-level unit test can be fully honest about the fixture and still miss the boundary.

## Release Notes

This close-out is written at the β R1 converge boundary, per γ's own process-gap-audit closeout (`.cdd/unreleased/608/gamma-closeout.md`): the review loop is closed (converge), but the actual `git merge` to `origin/main`, the close-keyword pre-merge gate row, and the issue-close have not executed yet — issue #608 remains `OPEN` at the time this artifact is written. That mechanical merge step, and everything at the release boundary (tag/deploy, cycle-directory move to `.cdd/releases/{X.Y.Z}/608/`), remains outstanding and is not claimed as done here. No release-boundary action is β's to take beyond executing the merge itself under the current bootstrap-mode stand-in (`beta/SKILL.md` "merge/no-merge judgment" vs. the not-yet-built mechanical runtime) — tag/deploy/disconnect stays δ's surface regardless.
