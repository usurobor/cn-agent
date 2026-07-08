# α self-coherence — cycle #612

## §R1

**Branch:** `cycle/612`, from `main` at the claimed base SHA.

### What I did

Reproduced all four symptoms against the built binary at base SHA (see `gamma-scaffold.md` §"Bug reproduction") before writing any fix, to confirm the issue text against reality rather than trusting it.

**AC1/AC2 (`src/go/cmd/cn/main.go`):** Added a `switch os.Args[1]` block for `--version`/`-V`/`--help`/`-h` immediately after the existing "no args → help" check, before `ResolveCommand` is called (so these two flags never hit the "Unknown command" path). Added a second interception after command resolution (`args := res.Remaining`), checking `args[0] == "--help" || args[0] == "-h"` and calling the new `cli.PrintCommandHelp`, placed *before* the `NeedsHub` gate — this is the actual fix for `doctor`/`status`/`deps`/`dispatch`, whose `--help` previously errored on the hub-requirement check before ever reaching their `Run`.

**Discovered while implementing AC2:** `DispatchCmd`, `CellReturnCmd`, `CellResumeCmd`, `CellFinalizeCmd` each already define a `Help() string` method with detailed usage text — but grep confirms `Run()` in each never calls it, and nothing else in the codebase calls it either. It was dead code. Rather than write new help text for these four, I added a `HelpProvider` interface (`help.go`) and wired `main.go`'s new interception to check for it, so these four commands' existing rich help becomes reachable for free. `RepoInstallCmd` already had *both* a `Help()` method and its own inline `--help` scan inside `Run` — the interface pickup gives it a second, earlier path to the same content (verified byte-identical output), and its own inline check remains as a harmless no-op for the common case (still useful for `--help` appearing anywhere in its arg list, not just position 0).

**Gap found:** `ActivateCmd` had inline `--help` handling inside `Run` but no `Help()` method — meaning the new generic interception in `main.go` would have *regressed* it (short-circuiting before `Run`, falling through to the generic fallback instead of its existing rich `activateHelp` text). Fixed by adding `Help() string` to `ActivateCmd` returning the existing `activateHelp` const — one-line addition, no content duplication.

**AC3 (`src/go/internal/cli/cmd_init.go`):** Added a check at the top of `Run`: if `Args[0]` is flag-shaped (`isFlag`, the existing helper from `cmd_activate.go`), refuse with a clear error naming the flag and do not call `hubinit.Run` at all. In the live binary this code path is mostly a defense-in-depth backstop — `--help`/`-h` specifically never reach it (main.go's generic interception catches them first) — but any *other* unrecognized flag (`--bogus`, etc.) does reach it and is refused here. `hubinit.validHubName` itself is unchanged (it still accepts `-` for legitimate hyphenated hub names like `cn-sigma`) — the fix belongs at the CLI-argv boundary, not the domain layer, matching the existing separation of concerns in this package (`RepoInstallCmd` resolves its own repo root via `gitRepoRoot` rather than pushing git logic into `hubinit`/`restore`).

**AC4 (`dispatch.go`, `cmd_help.go`):** Added `InvocationName(reg, name) string`, reusing the exact collision check `ResolveCommand`'s flat-form reject branch already performs (`GroupMembers(reg, prefix-before-first-hyphen)`). Wired `cmd_help.go`'s three listing loops (kernel/repo-local/package) to use it instead of raw `spec.Name`.

**Not in the original AC text, found by β during review (see `beta-review.md`):** `cmd_status.go`'s `registryToCommandInfo` forwards `spec.Name` unchanged into `hubstatus.CommandInfo`, which `hubstatus.go` prints verbatim — the identical bug, in a second location the issue's own AC4 text didn't name. Fixed with the same `InvocationName` call. `hubstatus` itself is untouched (per its own package doc, it deliberately takes pure data and does not import `cli/` — the fix belongs at the boundary in `cmd_status.go`, which already builds the `CommandInfo` slice and has `Registry` in scope).

### Tests written

- `dispatch_test.go`: `TestInvocationName` (collision vs. no-collision cases, using the existing `newTestRegistry` kata-group fixture).
- `help_test.go` (new): `TestPrintCommandHelp_UsesHelpProvider`, `TestPrintCommandHelp_GenericFallback`, `TestPrintCommandHelp_GenericFallback_UsesInvocationName`.
- `cmd_help_test.go` (new): `TestHelpCmd_DisplayNamesMatchInvocation` — asserts both presence of the space form and absence of the hyphenated form.
- `cmd_init_test.go` (new): `TestInit_RejectsUnrecognizedFlag`, `TestInit_HelpFlagDoesNotScaffold` (the literal cnos#606 C3 repro, run against `InitCmd.Run` directly), `TestInit_AcceptsPositionalHubName` (regression check for the normal path).
- `.github/workflows/build.yml`: extended the existing "Smoke test" job step with an end-to-end F1-F4 block against the actually-built binary — `--version`/`-V`, `--help`/`-h` on `cn` and all 15 registered commands (including `cn repo install --help`'s documented flags), `cn init --bogus` and `cn init --help` (both asserted to scaffold nothing), and `cn help`/`cn status` display-name checks (both hub-less and hubbed contexts, since `status` requires a hub).

All of `go build ./...`, `go vet ./...`, `go test ./...` (from `src/go`) and the three co-located modules (`cnos.issues/commands/issues-fsm`, `cnos.issues/commands/issues-map`, `cnos.cdd/commands/cdd-verify`) pass clean, no regressions. The dispatch-boundary CI check (cmd_*.go must not import `os`/`net/http`/`encoding/json`/`archive/`/`compress/`/`crypto/`/`path/filepath` directly) was manually re-checked against every changed `cmd_*.go` file — none of the changes introduce a new import of any of those packages.

### mock_parity

```yaml
mock_parity:
  source: docs/development/design/cn-repo-install-MOCKS.md
  rows:
    - id: F1
      expectation: "cn --version / -V prints the version (exit 0), not \"Unknown command\"."
      observed: "cn v3.82.0 (exit 0) for both --version and -V, at both top level and via the ldflags-injected version string."
      evidence: "src/go/cmd/cn/main.go switch block; .github/workflows/build.yml F1 smoke assertion; manually run against /tmp/cn-review by independent β review agent."
      verdict: match
      how: "Matches the invariant text exactly — no behavior beyond what F1 asks for."
    - id: F2
      expectation: "--help / -h works on cn and every command incl. cn repo install."
      observed: "cn --help/-h: full command listing, exit 0. All 15 registered kernel commands (doctor, status, deps, build, update, setup, dispatch, cell finalize/return/resume, issues fsm/map, cdd verify, activate, repo install) return exit 0 on --help, including previously-hub-gated ones (doctor/status/deps/dispatch) and previously-silent-no-op ones (build/update/setup). cn repo install --help still documents --release/--index/--packages/--dispatch/--dry-run verbatim."
      evidence: "src/go/internal/cli/help.go; main.go per-command interception; .github/workflows/build.yml F2 loop; help_test.go; cmd_repo_install.go unchanged (verified byte-identical --help output)."
      verdict: exceed
      how: "Exceeds the literal AC text by also fixing 4 commands (dispatch, cell finalize/return/resume) whose Help() text existed but was dead code — not a new capability, just making already-authored content reachable. Safe: no new user-facing surface, purely a bugfix to existing intent."
    - id: F3
      expectation: "cn init --help (any unrecognized --flag) refuses with a clear error and scaffolds nothing. Negative: any directory created from a flag (e.g. cn---help/) fails."
      observed: "cn init --bogus: exit 1, stderr names the flag, no directory created. cn init --help: exit 0 (usage text, per F2's own requirement that --help succeed globally), no directory created — main.go's generic --help interception now catches this before InitCmd.Run is ever reached. cn init myhub: unaffected, still creates cn-myhub/."
      evidence: "src/go/internal/cli/cmd_init.go isFlag guard; cmd_init_test.go (3 tests); .github/workflows/build.yml F3 block (both --bogus and --help checked against the real binary, asserting no directory of any kind appears)."
      verdict: match
      how: "F3's own wording names --help as an example of 'any unrecognized --flag' that must not scaffold; both the --help path (via F2's global interception) and the general unrecognized-flag path (via cmd_init.go's own guard) independently guarantee no stray directory, satisfying the negative clause by two layers rather than one."
    - id: F4
      expectation: "cn help display names match real invocation (cn issues fsm, cn cell finalize), or list both forms — no name/invocation mismatch."
      observed: "cn help now prints 'issues fsm', 'issues map', 'cell return', 'cell resume', 'cell finalize', 'repo install', 'cdd verify' — the space form that actually dispatches each command. The hyphenated internal registry keys no longer appear anywhere in the output."
      evidence: "src/go/internal/cli/dispatch.go InvocationName; cmd_help.go; dispatch_test.go TestInvocationName; cmd_help_test.go; .github/workflows/build.yml F4 block."
      verdict: exceed
      how: "Exceeds the literal AC text (scoped to 'cn help') by also fixing the identical bug in 'cn status' (found during β review, not named in the issue), since it is the same root cause and the same fix (InvocationName) applies cleanly at the same architectural boundary. Safe: cmd_status.go already builds CommandInfo from the Registry it has in scope; no new surface, no behavior change beyond display text."
  summary:
    matched: 2
    exceeded: 2
    missed: 0
    exceed_justified: true
```
