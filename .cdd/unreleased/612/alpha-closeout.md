# α close-out — cycle #612

**Verdict at hand-off:** β converged in R0 after one blocking finding (fixed same round — see `beta-review.md` Finding 1 / `self-coherence.md`).

**Diff summary:**
- `src/go/cmd/cn/main.go` — top-level `--version`/`-V`/`--help`/`-h`; per-command `--help`/`-h` before the hub gate.
- `src/go/internal/cli/dispatch.go` — `InvocationName`.
- `src/go/internal/cli/help.go` (new) — `HelpProvider` interface + `PrintCommandHelp`.
- `src/go/internal/cli/cmd_help.go`, `cmd_status.go` — use `InvocationName` for display.
- `src/go/internal/cli/cmd_init.go` — refuse unrecognized flags.
- `src/go/internal/cli/cmd_activate.go` — `Help()` method (parity fix, prevents regression from the new generic interception).
- `.github/workflows/build.yml` — F1-F4 end-to-end smoke test against the real binary.
- Tests: `dispatch_test.go` (+1 test), `help_test.go` (new, 3 tests), `cmd_help_test.go` (new, 1 test), `cmd_init_test.go` (new, 3 tests).

**No production behavior change outside the four named commands' CLI-argv handling.** No FSM change, no installer change (issue's own non-goals), no schema change, no new external dependency.

**Debt disclosed (none blocking):**
1. CI smoke test's F4 leak-check is a hand-enumerated grep list (7 names) rather than a fully generic assertion over the registry — acceptable given the existing style of the surrounding smoke-test step, flagged by β as a non-blocking observation, not filed as a follow-up issue (low value: a new command with this exact collision shape is rare, and the Go unit tests already cover the mechanism generically via `TestInvocationName`).

**mock_parity:** see `self-coherence.md` — F1/F3 match, F2/F4 exceed (both justified, both safe, both are "fix the same class of bug in an adjacent location β found," not new capability). 0 missed.

**Build/test status at hand-off:** `go build ./...`, `go vet ./...`, `go test ./...` (src/go) — all green. `issues-fsm`, `issues-map`, `cdd-verify` co-located modules — all green. Manual end-to-end verification against the built binary for every AC1-AC4 scenario, run twice (once pre-status.go-fix, once post-fix).
