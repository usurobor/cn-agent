# β review — cycle #612

## §R1 (R0 — first pass)

**Method:** Did not trust α's report. Built the binary independently (`go build -o /tmp/cn-review ./cmd/cn`) from the staged diff and drove it directly for every AC1-AC4 scenario named in the issue plus adjacent edge cases (hub-gated commands' `--help`, previously-dead `Help()` methods, `repo install`'s pre-existing rich help, regressions on valid `cn init <name>`). Ran the full build/vet/test matrix across `src/go` and the three co-located modules (`cnos.issues/commands/issues-fsm`, `cnos.issues/commands/issues-map`, `cnos.cdd/commands/cdd-verify`). Checked the dispatch-boundary CI guard against every changed `cmd_*.go` file. Read the new/changed Go tests for whether they're meaningful or tautological. Read the new CI smoke-test shell block for scripting bugs.

### Finding 1 (blocking → fixed same round)

`cn status`, run inside any hub, prints the raw hyphenated registry keys (`cdd-verify`, `cell-return`, `cell-finalize`, `issues-fsm`, `issues-map`, `repo-install`) — the identical AC4 bug class fixed in `cn help`, missed because the issue's own AC4 text names only `cn help`. Root cause: `cmd_status.go`'s `registryToCommandInfo()` forwards `spec.Name` unchanged into `hubstatus.CommandInfo.Name`, and `hubstatus.go` prints it verbatim.

Reproduction: `cn init tenant && cd cn-tenant && cn status` (pre-fix) shows `cdd-verify`, `cell-finalize`, etc. in the "Commands:" section.

**Disposition:** Same root cause, same fix (`InvocationName`), same architectural boundary (`cmd_status.go` already builds `CommandInfo` from a `Registry` it has in scope; `hubstatus` package deliberately stays `cli`-import-free per its own doc comment, so the fix correctly lives at the call site, not inside `hubstatus`). Fixed in this round — see `self-coherence.md` §"Not in the original AC text". Re-verified independently post-fix (`cn init verifyhub && cd cn-verifyhub && cn status` — clean output, all seven previously-hyphenated names now show their space form).

### Non-blocking observations

1. `cn init --help`'s exact behavior (usage text, exit 0) diverges from the illustrative transcript in `docs/development/design/cn-repo-install-MOCKS.md` §Mock F, which shows `--help` itself being refused with an "unknown flag" error. This is not a defect: AC2's own text requires `--help`/`-h` to work globally, which necessarily takes priority over the mock's illustrative (and, given AC2, self-contradictory) transcript — F3's own wording separately confirms `--help` is only "any unrecognized `--flag`" as an *example* of the underlying bug class (a flag being silently treated as a positional hub name), not a literal requirement that `--help` itself must error. Both interpretations are satisfied here: no stray directory is ever created either way (checked in both the `cmd_init_test.go` unit tests and the CI smoke test's F3 block).
2. The CI smoke test's F4 leak-check grep list is a fixed set of 7 command names — it will not catch a *new* future command with the same collision shape unless someone remembers to add it to the grep. Acceptable for this cycle (matches the existing style of hand-enumerated command lists already present in `build.yml`'s other steps, e.g. the `for c in ...` loop two lines above it); a fully generic regression guard would need a small dedicated Go test asserting `InvocationName` is called for every registered command's display path, which is out of scope for this cycle's AC text.
3. Pre-existing cosmetic misalignment in `cn help`'s `%-12s` column padding for names ≥13 chars (e.g. `cell finalize`) predates this diff and is unrelated to it — not a finding, noted only to distinguish it from actual review findings.

### Test quality check

`cmd_help_test.go`, `cmd_init_test.go`, `help_test.go`, and the new `TestInvocationName` case in `dispatch_test.go` were read in full: each exercises a real collision/no-collision case or a real negative-path assertion (no directory created, error names the flag), not a tautology. The CI smoke-test block was read line by line for shell correctness (quoting, `set -euo pipefail` behavior, `cd -` usage, cleanup on both success and failure paths) — no bugs found after Finding 1's status.go gap prompted a closer second pass of the whole block.

### Verdict

**converge.** One blocking finding (status.go AC4-class gap), fixed within this same review round and independently re-verified against the rebuilt binary. No other blocking findings. All of AC1-AC4 hold against the actual binary, not just the test suite's framing of them.
