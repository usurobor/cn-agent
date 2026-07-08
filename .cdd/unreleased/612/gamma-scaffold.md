# γ scaffold — cycle #612

**Issue:** [cnos#612](https://github.com/usurobor/cnos/issues/612) — "cds-install Sub 5: `cn` CLI ergonomics for onboarding (--help / --version / init-flag / name-align)"
**Mode:** `design-and-build` (per the issue's own `**Mode:**` field).
**cell_kind:** `implementation` (defaulted; issue body carries no explicit `cell_kind:` field).
**Refs:** #607 (wave master tracker). **Customer source:** cnos#606 C3 (first external tenant, `tsc`).
**Design surface:** `docs/development/design/cn-repo-install-MOCKS.md` §"Mock F — CLI ergonomics" (on `main` already — confirmed present at `.../cn-repo-install-MOCKS.md:209-232`, no design-branch pull needed, unlike #608's Mock A/B/E1).
**Base SHA:** `main` HEAD at claim time (see CLAIM-REQUEST.yml / claim comment on the issue).
**Branch:** `cycle/612`, created from `main` at the base SHA, no prior commits.
**Branch pre-flight:** no prior `origin/cycle/612`; no prior `.cdd/unreleased/612/` artifacts beyond this wake's own `CLAIM-REQUEST.yml`; no `status:changes` history; issue state `OPEN`. Classified `run_class: first_pass` (see the claim comment) — repair re-entry preflight (cds-dispatch/SKILL.md §"Repair re-entry preflight") not required.

---

## Bug reproduction (verified against the built binary before any fix)

All three symptoms named in the issue reproduce exactly as described, confirmed by building `./src/go/cmd/cn` at the claimed base SHA and running it directly (not inferred from the issue text):

1. `cn --version` / `-V` → `✗ Unknown command: --version` / `-V`.
2. `cn --help` / `-h` (top level) → same "Unknown command" path. Per-command `--help` is inconsistent: `doctor`/`status`/`deps`/`dispatch` (all `NeedsHub: true`) error `"requires a hub"` before ever looking at `--help`; `build`/`update`/`setup` (all `NeedsHub: false`, no internal `--help` handling) silently no-op, exit 0, print nothing; `repo install` and `activate` already handle `--help` correctly (existing inline checks). `dispatch`, `cell return`, `cell resume`, `cell finalize` each already carry a rich `Help() string` method — but nothing in `Run()` ever calls it; the method is dead code, never reached.
3. `cn init --help` → creates a directory literally named `cn---help/`. Root cause: `hubinit.validHubName` accepts `-` as a valid hub-name character, and `cmd_init.go` treats any positional `Args[0]` as the hub name with no flag check — `"--help"` passes `validHubName` unmodified.
4. (AC4, not independently reproduced as a "failure" but confirmed as the described mismatch): `cn help` prints internal registry keys (`issues-fsm`, `cell-finalize`, `repo-install`, `cdd-verify`, `cell-return`, `cell-resume`, `issues-map`) verbatim. None of these hyphenated flat tokens actually dispatch to the command — `cli.ResolveCommand`'s flat-form path explicitly rejects any name whose pre-first-hyphen prefix collides with a sibling noun-group member (`dispatch.go` lines ~58-69), which is true for every one of the seven names above (each has at least one sibling under the same noun, or — for the singleton "cdd"/"repo" groups — collides with itself, since `GroupMembers` includes the exact-name registration). Typing the displayed name literally (e.g. `cn issues-fsm`) resolves to a noun-group listing, not the command; only the two-word form (`cn issues fsm`) actually runs it.

---

## Surfaces touched

| Path | Role |
|---|---|
| `src/go/cmd/cn/main.go` | Add top-level `--version`/`-V` and `--help`/`-h` interception (before `ResolveCommand`); add per-command `--help`/`-h` interception (before the `NeedsHub` gate, after command resolution). |
| `src/go/internal/cli/dispatch.go` | Add `InvocationName(reg, name) string` — returns the noun-verb space form for any name shadowed by a noun-group collision (reusing the existing `GroupMembers` check `ResolveCommand` already performs for its own reject logic), unchanged otherwise. |
| `src/go/internal/cli/help.go` (new) | `HelpProvider` interface (`Help() string`) + `PrintCommandHelp(w, reg, cmd)`: prints the command's own `Help()` text if it implements the interface, else a generic `"Usage: cn <InvocationName> [args...]\n\n<Summary>\n"` fallback. |
| `src/go/internal/cli/cmd_help.go` | Use `InvocationName` instead of raw `spec.Name` in all three listing sections (kernel/repo-local/package). |
| `src/go/internal/cli/cmd_status.go` | Same fix as `cmd_help.go`, applied to `registryToCommandInfo` (β found this during review — `cn status` had the identical display bug, not named in the issue's own AC text but the same bug class). |
| `src/go/internal/cli/cmd_init.go` | Refuse any `Args[0]` that `isFlag()` (existing helper from `cmd_activate.go`) before treating it as the positional hub name. |
| `src/go/internal/cli/cmd_activate.go` | Add `Help() string` (returning the existing `activateHelp` const) so the new generic `--help` interception in `main.go` doesn't regress `activate`'s already-correct rich help output. |
| `.github/workflows/build.yml` | Extend the existing "Smoke test" step with an end-to-end F1-F4 oracle block against the real built binary. |

**No change:** `RepoInstallCmd`'s own inline `--help` handling (left in place — still correct, still reachable for `--help` anywhere in its arg list, not just position 0); `hubinit.go` (the flag guard is added at the `cli` dispatch layer, not the domain layer, matching the existing `RepoInstallCmd`/`gitRepoRoot` precedent of keeping domain packages free of CLI-argv concerns); FSM/installer behavior (issue's own §Out-of-scope).

---

## AC oracle approach

| AC | Concrete check |
|---|---|
| **AC1** | `cn --version` / `-V` exit 0, stdout matches `^cn v`. |
| **AC2** | `cn --help` / `-h` exit 0. Every registered command's `--help` exits 0 (looped in the CI smoke test and exercised individually in Go unit tests for the generic-fallback and HelpProvider paths). `cn repo install --help` still documents `--release`/`--index`/`--packages`/`--dispatch`/`--dry-run` (pre-existing content, unaffected by the new dispatch path — verified byte-identical). |
| **AC3** | `cn init --help` and `cn init --anything-unrecognized` both: exit nonzero from `InitCmd.Run` when reached directly (Go unit tests), and — critically — never create a directory at all when run through the real binary (CI smoke test; the `--help` case never even reaches `InitCmd.Run` in the live binary, since `main.go`'s generic interception now short-circuits first — the negative "no stray dir" assertion covers both paths). `cn init myhub` (valid positional name) still creates `cn-myhub/` (regression check). |
| **AC4** | `cn help` and `cn status` never print `issues-fsm`, `cell-finalize`, `cell-return`, `cell-resume`, `repo-install`, `cdd-verify`, or `issues-map` — only their space-separated invocable forms. Checked in Go unit tests (`TestHelpCmd_DisplayNamesMatchInvocation`, `TestInvocationName`) and the CI smoke test against the live binary in both a hub-less and a hubbed context (`cn status` requires a hub). |

---

## mock_parity contract (pinned; α/β must fill actual rows in their closeouts)

```yaml
mock_parity:
  source: docs/development/design/cn-repo-install-MOCKS.md
  rows:
    # F1-F4 only — this cell's scope
  summary:
    matched: <n>
    exceeded: <n>
    missed: 0   # MUST be 0 to converge
```

---

## α prompt

Implement the fixes per the "Surfaces touched" table above. Build the binary and manually verify each of AC1-AC4 against it (not just unit tests) before calling this done — the bug class here is specifically "looks fixed in isolation, still broken through main.go's actual dispatch path" (e.g. a command with its own `Help()` method that nothing ever calls), so end-to-end verification against the real binary is the load-bearing check, not a substitute for one. Do not touch FSM or installer behavior (issue's own non-goal). Write the `mock_parity` block with real `observed`/`evidence` per row in `alpha-closeout.md`.

## β prompt

Independently rebuild the binary and re-run AC1-AC4 yourself — do not trust α's report. Specifically try: commands whose `NeedsHub` is true (confirm `--help` no longer requires a hub); commands that already had a `Help()` method sitting unused (confirm it's now wired, not replaced by a worse generic fallback); any other command surface that prints `spec.Name` directly to a user (this is exactly how the `cn status` instance of the AC4 bug was found — grep for other call sites, not just `cmd_help.go`). Confirm the CI smoke-test shell script is correct (quoting, exit codes, cleanup) since it runs unattended.
