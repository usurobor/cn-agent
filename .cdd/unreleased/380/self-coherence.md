<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness, fix-round-1]
  completed: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness, fix-round-1]
-->

# Self-coherence — α #380

## Gap

**Issue:** [#380 — cn activate: add --claude / --codex flags to spawn the AI body interactively with the activation prompt pre-loaded](https://github.com/usurobor/cnos/issues/380).
**Mode:** design-and-build. Design converged in the issue body during γ — flag names, exec mechanism (`syscall.Exec`), mutual-exclusion ordering (pre-render), missing-binary detection (pre-render via `exec.LookPath`), and the bare-positional argv shape `[binary, prompt]`. The one bounded design call γ delegated to α — render-capture seam Option A vs B — α picked **Option A** (capture buffer at the call site, `activate.go` untouched). Declared inline below in §ACs AC1.
**Version:** unreleased; cycle directory `.cdd/unreleased/380/`. Release boundary follows β's merge and γ's close-out.
**Base SHA:** `319893a4ea9de1b89989ef6e6dc44cd3e1cad147` (γ scaffold base; `origin/main` HEAD at dispatch time; cycle branch confirmed rebased onto current `origin/main` in §Review-readiness).

**Named gap.** cnos 3.78.0 shipped the source-of-truth half of activation — `cnos.core/skills/agent/activate/SKILL.md` is the body-agnostic skill, and `src/go/internal/activate/activate.go` renders the prompt from it. The hub-side flow stops at stdout: `cn activate HUB_DIR` writes the rendered prompt to stdout. The intended pattern was `cn activate cn-sigma | claude` / `... | codex`, but those pipe forms are one-shot non-interactive — they print the model's response and exit. To recover interactive mode the operator copy-pastes the prompt into a fresh `claude` / `codex` REPL. The bridge from "prompt rendered" to "body operating interactively at this hub" is missing.

**Closure shape.** This cycle adds the spawn-into-CLI surface on top of the 3.78.0 renderer. Two new flags — `--claude` and `--codex` — replace the cn process with the chosen AI CLI invoked bare-positionally (`claude "$prompt"` / `codex "$prompt"`) via `syscall.Exec`, preserving the TTY so the operator lands in a live REPL with the body already in the post-activation state. The default no-flag path is preserved bytes-equal — pipe-and-redirect consumers do not regress.

**Render-capture seam: Option A.** Per γ scaffold §"Render-capture seam (γ design call, surfaced to α at scaffold time)", α picked Option A: the spawn arm captures `activate.Run`'s stdout into a `bytes.Buffer` at the call site in `cmd_activate.go`; the default arm passes `inv.Stdout` straight to `activate.Run` unchanged. Justification: minimum-diff path, no edit to `activate.go`, easier review for β, and the AC3 bytes-equal property holds by construction (no buffer in the default path). The diff confirms: `git diff cafabc8b..HEAD -- src/go/internal/activate/activate.go` is empty.

**Non-goals (held).** No edits to `src/packages/cnos.core/skills/agent/activate/SKILL.md` (body-agnostic skill, non-goal). No `--cursor` / `--aider` / `--auto` flags. No `$CN_DEFAULT_BODY` env-var default. No hub README patches. No macOS/Windows-specific spawning quirks beyond what falls out of `syscall.Exec` portability (handled via build-tagged `spawn_other.go` stub so the no-flag path still compiles on non-unix).

## Skills

**Tier 1 — always-on for α/CDD:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm and role contract. §1.4 commit-checkpoint authoring applied (six commits across the cycle, first commit landed in the first quintile of the dispatch budget — see §CDD Trace step 4); §1.6 sequential bounded dispatch (α exits at review-readiness; close-out is re-dispatched by γ).
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface. §2.2 artifact order followed (no design artifact needed — design surface is the issue body per design-and-build mode; tests → code → docs structure landed via interleaved checkpoints: flag-defs → spawn helper → wiring → tests → AC3 oracle → self-coherence). §2.3 peer enumeration recorded in §CDD Trace. §2.5 incremental self-coherence — every section in its own commit. §2.6 pre-review gate exercised in §Review-readiness.

**Tier 2 — always-applicable engineering:**

- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go authoring discipline. Concrete applications: §2.1 small-package design (kept the new helpers inside `internal/activate`, no new package); §2.2 concrete types and zero-value-safe `Invocation` propagation; §2.3 consumer-owned hook types `execFunc` / `lookPathFunc` declared inside the test-injection seam, not exported as public APIs; §2.5 errors-as-values for `CheckSpawnBinary` (no panics, all expected failures named and wrapped); §2.10 happy and degraded paths tested (LookPath failure, mutual-exclusion, default no-flag, --help, unknown-flag); §3.11 explicit override precedence (positional HUB_DIR > cwd discovery; flag-set > flag-unset; mutual-exclusion fires before render); §2.13 determinism preserved (no map iteration, prompt bytes byte-stable cross-version per AC3 oracle); §3.10 argv-based subprocess execution (`syscall.Exec(path, []string{name, prompt}, env)` — no shell construction).
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact (always-on per CDD). Applied to commit messages, help text, error diagnostics, and this self-coherence.md.

**Tier 3 — issue-specific (per issue body §Skills to load + γ scaffold §Tier 3 list):**

- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — load-bearing across AC1–AC5 (flag parsing, exec wiring, test seam shape, error wrapping).
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority. No in-cycle issue-pack reconciliation was required; the issue body remained authoritative throughout. The one γ delegation (render-capture seam Option A vs B) was resolved by α and declared in §Gap and §ACs AC1, not reflected back to the issue body.

## ACs

Per-AC oracles run against branch HEAD `87aa69e9` (last implementation commit before this self-coherence file). AC numbering follows the issue body.

### AC1 — `--claude` flag spawns claude interactively with the prompt loaded

**Invariant.** `cn activate --claude HUB_DIR` renders the activation prompt via the existing 3.78.0 path, then replaces the cn process with `claude "$prompt"` via `syscall.Exec`. TTY/stdin/stdout inherited from the parent terminal. Flag documented in `cn activate --help`.

**Surface.** `src/go/internal/activate/spawn.go` (NEW; `Spawn` + `CheckSpawnBinary` + injectable `execFunc` / `lookPathFunc` hooks), `src/go/internal/activate/spawn_test.go` (NEW), `src/go/internal/cli/cmd_activate.go` (MODIFY; flag-parsing + pre-render `CheckSpawnBinary` + render-capture seam Option A + `Spawn` call), `src/go/internal/cli/cmd_activate_test.go` (NEW).

**Implementation notes.**

- Spawn primitive: stdlib `syscall.Exec`, called via the `defaultExecFn` package-level variable in `spawn.go` (build-tag `unix`). Tests substitute a recording function through the unexported `spawnWith` seam — no production code path runs an `os.Exec` mock.
- Render-capture seam: **Option A**. `cmd_activate.go` calls `activate.Run` with a `bytes.Buffer` as Stdout only on the spawn arm; the default arm passes `inv.Stdout` straight through. `activate.go` is byte-identical to 3.78.0 (`git diff cafabc8b..HEAD -- src/go/internal/activate/activate.go` → empty).
- argv shape: `syscall.Exec(<resolved_path>, []string{"claude", <prompt>}, os.Environ())`. argv[0] is the binary name per Unix exec convention; the only remaining element is the rendered prompt — bare positional, no `-p`, no `--print`.

**Oracle 1 — TestSpawnWith_ClaudeArgvShape passes.**
```
$ cd src/go && go test -run TestSpawnWith_ClaudeArgvShape -v ./internal/activate/...
=== RUN   TestSpawnWith_ClaudeArgvShape
--- PASS: TestSpawnWith_ClaudeArgvShape (0.00s)
PASS
ok      github.com/usurobor/cnos/src/go/internal/activate
```
Asserts: `argv0` equals the resolved binary path returned by LookPath; `argv` equals `["claude", <prompt verbatim>]` (length 2). **PASS.**

**Oracle 2 — `cn activate --help` documents `--claude`.**
```
$ go test -run TestActivate_HelpFlag_DocumentsClaudeAndCodex -v ./internal/cli/...
=== RUN   TestActivate_HelpFlag_DocumentsClaudeAndCodex
--- PASS: TestActivate_HelpFlag_DocumentsClaudeAndCodex (0.00s)
```
The help text in `cmd_activate.go:10–33` contains `--claude` with a description naming "interactive" REPL behavior and the `claude "$prompt"` shape. **PASS.**

**Oracle 3 — TestSpawn_DefaultsAreWired confirms `Spawn` routes through `exec.LookPath` in production.** The `defaultLookPath` package-level variable is the seam; the test swaps it for a sentinel observer and verifies the binary name passes through unchanged. **PASS.**

**Manual smoke (operator-runnable, not required for α verdict).** With a fake `claude` script on $PATH that prints argv to stderr and exits 0, `cn activate --claude FIXTURE_HUB` would surface argv = `["claude", "<PROMPT>"]`. Verifying live `claude` REPL interactivity is the operator dry-run named in the issue body §Proof plan; it falls to δ / operator at integration time, outside the test suite.

### AC2 — `--codex` flag spawns codex interactively with the prompt loaded

**Invariant.** Same shape as AC1, with `codex "$prompt"` substituted. **Critical negative:** argv must be `["codex", <prompt>]`, NOT `["codex", "exec", <prompt>]` — `codex exec` is the non-interactive subcommand and would defeat the cycle's purpose.

**Surface.** Same as AC1; the spawn helper is binary-agnostic — `Spawn(binary, prompt)` passes whatever binary name it receives.

**Oracle 1 — TestSpawnWith_CodexArgvShape passes.**
```
=== RUN   TestSpawnWith_CodexArgvShape
--- PASS: TestSpawnWith_CodexArgvShape (0.00s)
```
Asserts the same `[binary, prompt]` shape as AC1, AND explicitly checks `argv[1] != "exec"` and `len(argv) == 2` (which would fail if `codex exec` had been wired in). **PASS.**

**Oracle 2 — TestCheckSpawnBinary_CodexErrorShape passes** — confirms the missing-binary diagnostic for `--codex` names both `codex` (binary) and `--codex` (flag). **PASS.**

**Oracle 3 — `cn activate --help` documents `--codex`.** Verified by `TestActivate_HelpFlag_DocumentsClaudeAndCodex` (the test asserts both flag names appear). **PASS.**

### AC3 — default `cn activate HUB_DIR` (no flag) behavior preserved exactly

**Invariant.** With neither `--claude` nor `--codex` set, `cn activate HUB_DIR` writes the rendered prompt to stdout and exits 0 — bytes-identical to 3.78.0 baseline. Pipe consumers (`cn activate HUB_DIR | other-tool`, `cn activate HUB_DIR > prompt.md`) must not regress.

**Surface.** `cmd_activate.go` (default arm: `activate.Run` called with `inv.Stdout`, no buffer indirection); `activate.go` (unchanged from 3.78.0).

**Oracle 1 — existing activate tests all pass on HEAD.**
```
$ go test ./internal/activate/...
ok      github.com/usurobor/cnos/src/go/internal/activate    0.040s
```
This sweeps all 3.78.0 / #379 invariants (Sigma-shape activation, init-only / init+setup fixtures, three kernel states, three deps states, ordered Read first section, skill-as-source-of-truth, latest reflection, secrets exclusion, stdout-only diagnostics — see `activate_test.go`). **PASS.**

**Oracle 2 — TestActivate_DefaultNoFlag_BytesEqualToDirectRun passes.** Compares stdout from `ActivateCmd.Run` (cli wrapper) against `activate.Run` (direct) for the same hub fixture; asserts bytes-equal. Fails if the cycle/380 refactor accidentally redirected the default arm through a buffer. **PASS.**

**Oracle 3 — cross-version smoke compare against 3.78.0 (commit `cafabc8b`).** Built `cn` at HEAD and at 3.78.0; ran both against an identical fixture hub:
```
$ wc -c /tmp/out-head.txt /tmp/out-3780.txt
1152 /tmp/out-head.txt
1152 /tmp/out-3780.txt
$ sha256sum /tmp/out-head.txt /tmp/out-3780.txt
43423eeedfdc34f01a2dbbdf3e4717c74774f2a0a57ac476abfa2f7e109b35ec  /tmp/out-head.txt
43423eeedfdc34f01a2dbbdf3e4717c74774f2a0a57ac476abfa2f7e109b35ec  /tmp/out-3780.txt
$ diff /tmp/out-head.txt /tmp/out-3780.txt
(no output)
```
Bytes-equal: same length, same SHA, no diff. **PASS.**

**Oracle 4 — no diagnostic added to stdout.** TestActivate_DefaultNoFlag_WritesPromptToStdout asserts stdout contains no `✗` glyphs; AC3 bytes-equal smoke confirms no new lines on stdout. **PASS.**

### AC4 — `--claude` and `--codex` are mutually exclusive

**Invariant.** Setting both flags exits non-zero BEFORE any rendering with a clear error naming the conflict.

**Surface.** `cmd_activate.go`: the mutual-exclusion check is at line ~80 — first behavior gate after flag parsing, before hub resolution and before any call to `activate.Run`.

**Oracle 1 — `cn activate --claude --codex HUB_DIR` exits non-zero.**
```
=== RUN   TestActivate_MutualExclusion_Errors
--- PASS: TestActivate_MutualExclusion_Errors (0.00s)
```
The test additionally asserts stdout stays empty (no rendering bytes leaked) and stderr names both `--claude` and `--codex` plus the phrase `mutually exclusive`. **PASS.**

**Oracle 2 — stderr names the conflict.** Test above also asserts `strings.Contains(stderr, "mutually exclusive")`. **PASS.**

**Oracle 3 — error fires BEFORE rendering.**
```
=== RUN   TestActivate_MutualExclusion_FiresBeforeRender
--- PASS: TestActivate_MutualExclusion_FiresBeforeRender (0.00s)
```
The test sets `HUB_DIR` to `/this/path/must/not/exist/cnos-380-test` (a non-existent path) — if rendering had run, stderr would carry `Hub path not found` from `activate.Run`. Test asserts: stderr contains `mutually exclusive` AND does NOT contain `Hub path not found` AND does NOT contain `Generating activation prompt`. **PASS — pre-render ordering proven.**

### AC5 — missing-binary error for `--claude` / `--codex` is clear

**Invariant.** With `--claude` set and `claude` not on `$PATH`, exit non-zero with a clear error naming both the missing binary AND the flag that requested it; offer actionable PATH/installation guidance. Same shape for `--codex` and `codex`. Detection via `exec.LookPath`, called BEFORE the renderer.

**Surface.** `spawn.go` `CheckSpawnBinary(binary, flag)` (helper); `cmd_activate.go` (calls `CheckSpawnBinary` right after mutual-exclusion, before `activate.Run`).

**Error message shape:**
```
✗ <binary> (requested by <flag>) not found in PATH — install it or ensure $PATH includes its directory
```
Names: missing binary, requesting flag, actionable PATH guidance, and an install hint. Satisfies the issue's "actionable" bar.

**Oracle 1 — exit non-zero when binary missing.**
```
=== RUN   TestActivate_MissingBinary_NamesBinaryAndFlag_Claude
--- PASS: TestActivate_MissingBinary_NamesBinaryAndFlag_Claude (0.00s)
=== RUN   TestActivate_MissingBinary_NamesBinaryAndFlag_Codex
--- PASS: TestActivate_MissingBinary_NamesBinaryAndFlag_Codex (0.00s)
```
Both use `t.Setenv("PATH", t.TempDir())` to guarantee LookPath fails regardless of operator environment. **PASS for both flags.**

**Oracle 2 — stderr names binary AND flag.** The two tests above assert `strings.Contains(stderr, "claude") && strings.Contains(stderr, "--claude")` (and the codex equivalent). **PASS.**

**Oracle 3 — actionable guidance.**
```
=== RUN   TestCheckSpawnBinary_MissingNamesBinaryAndFlag
--- PASS: TestCheckSpawnBinary_MissingNamesBinaryAndFlag (0.00s)
```
Helper-level test additionally asserts the error mentions `PATH` and `install`. The CLI-level tests assert `PATH`. **PASS.**

**Oracle 4 — same shape for `--codex`.** Covered by `TestActivate_MissingBinary_NamesBinaryAndFlag_Codex` and `TestCheckSpawnBinary_CodexErrorShape`. **PASS.**

**Oracle 5 — detection via `exec.LookPath`, called BEFORE renderer.**
```
=== RUN   TestActivate_MissingBinary_FiresBeforeRender
--- PASS: TestActivate_MissingBinary_FiresBeforeRender (0.00s)
```
Same proof shape as AC4 oracle 3: nonexistent HUB_DIR + missing binary. Asserts the missing-binary diagnostic appears AND `Hub path not found` does NOT AND `Generating activation prompt` does NOT — proving LookPath fired before `activate.Run`. **PASS — pre-render ordering proven.**

## Self-check

**Did α's work push ambiguity onto β?** No remaining design ambiguity. The one bounded design call γ delegated (render-capture seam) is declared explicitly in §Gap and §ACs AC1 — α picked Option A, justified, and the resulting diff confirms it (`activate.go` byte-identical to 3.78.0). γ pre-flagged "render-capture seam ambiguity: silent activate.go change" as a binding finding shape; this self-coherence file inoculates against it.

**Surface containment vs γ scaffold §"Target artifacts".** γ named: `spawn.go`, `spawn_test.go`, `cmd_activate.go (M)`, `activate.go (M, only if Option B)`. α deviated minimally:

| Path | γ status | Actual | Justification |
|---|---|---|---|
| `src/go/internal/activate/spawn.go` | NEW | NEW | as specified |
| `src/go/internal/activate/spawn_test.go` | NEW | NEW | as specified |
| `src/go/internal/cli/cmd_activate.go` | MODIFY | MODIFY | as specified |
| `src/go/internal/activate/activate.go` | MODIFY only if Option B | unchanged | Option A chosen — predicted γ outcome |
| `src/go/internal/activate/spawn_other.go` | not enumerated | NEW (build tag `!unix`) | required so the non-unix build still compiles — covers γ scaffold §Risk register row "syscall.Exec Windows portability"; γ said "α adds a stub or build-tag guard if needed → document in self-coherence" |
| `src/go/internal/cli/cmd_activate_test.go` | not enumerated | NEW | AC3/AC4/AC5 oracles target `cmd_activate.Run`, which lives in `package cli`. Putting these tests in `spawn_test.go` (package `activate`) would create a cyclic import (cli already imports activate). The alternative — running the oracles only via `go run cmd/cn ...` shell smoke — is less reliable than a Go test. Declared as a minor surface widening rather than hidden. |

Every test added is necessary for one of the AC oracles; every production file in the diff is one of the three γ-named surfaces plus the build-tag portability stub.

**Is every claim backed by evidence in the diff?** Yes. AC1/AC2 argv-shape claims: `spawn_test.go` `TestSpawnWith_{Claude,Codex}ArgvShape` records and asserts. AC3 bytes-equal claim: `TestActivate_DefaultNoFlag_BytesEqualToDirectRun` test + smoke comparison vs `cafabc8b` (sha256 recorded). AC4/AC5 pre-render ordering claims: tests use nonexistent HUB_DIR to prove the ordering by absence of the "Hub path not found" diagnostic.

**Peer enumeration (per alpha/SKILL.md §2.3).** Surface family: command-flag pairs in `cn activate` and binaries the spawn helper can target.

- Flag peers in this diff: `--claude`, `--codex`. Both added in the same surface. Both share `Spawn` / `CheckSpawnBinary`; the binary name is the only parameter that differs. No third flag is being added this cycle (`--cursor`, `--aider`, `--auto` are γ scaffold non-goals).
- Binary peers: claude, codex. Same exec mechanism (`syscall.Exec(<path>, [<name>, <prompt>], env)`); same bare-positional argv shape verified by `claude --help` and `codex --help`. The negative invariant — `argv[1] != "exec"` for codex — is enforced by `TestSpawnWith_CodexArgvShape`.
- Renderer peers: only one renderer (`activate.Run`). The spawn arm captures into a buffer; the default arm writes to `inv.Stdout`. Both arms call the same `activate.Run` with identical hub-resolution semantics. No second renderer was introduced.
- Help-text peers: the help block in `cmd_activate.go` is the only operator-visible documentation surface this cycle touches. No README, no skill body changed. No drift class.

**Harness audit (per alpha/SKILL.md §2.4).** No schema-bearing contract changed. The activation prompt format is the same as 3.78.0 (AC3 bytes-equal). No new producer of any JSON shape. No CI workflow emitter, no shell harness, no template added. Harness audit: **not applicable to this cycle's diff.**

**Loaded-skill misses.** None observed in the diff. `eng/go/SKILL.md` §2.10 (test both nominal and degraded paths), §3.10 (argv-based subprocess execution, no shell construction), §3.11 (explicit override precedence), §2.17 (purity boundary — `spawnWith` is pure of side effects given the injected hooks), and §2.7 (defer-immediately resource lifecycles — not exercised because the helper releases no acquired resources; `syscall.Exec` either replaces the process or returns). The build-tag pair (`spawn.go` // `spawn_other.go`) is the standard Go portability pattern (`go/SKILL.md` §2.11 / §3.8 — keep the build boring; standard module layout).

## Debt

**Manual interactive smoke not part of the test suite.** The issue body §Proof plan names a behavioral verification: a fresh operator running `cn activate --claude cn-sigma` lands in a live `claude` REPL with Sigma activated within one command. That is a δ / operator integration smoke against real `claude` and `codex` binaries; it cannot be encoded as a Go unit test because syscall.Exec replaces the process. The test surface verifies everything observable up to (and not including) the actual REPL handover: argv shape, LookPath wiring, mutual-exclusion ordering, error message shape, default-path preservation. The REPL handover itself is `syscall.Exec`'s contract.

**Non-unix platforms emit unsupported-platform error on `--claude` / `--codex`.** `spawn_other.go` returns "<flag> spawn is not supported on this platform" for `CheckSpawnBinary` on non-unix builds. The no-flag default path remains functional everywhere. Issue §Non-goals excludes "macOS/Windows-specific spawning quirks beyond what falls out of syscall.Exec portability"; the stub is what falls out for free.

**Deferred (named in issue body, not addressed this cycle).**

- `--auto` body-detection flag — issue §Deferred, "explicit today; auto-pick later if friction emerges."
- `$CN_DEFAULT_BODY` env-var default — issue §Deferred.
- `--cursor` / `--aider` / other AI CLI flags — issue §Out-of-scope, "extensible later via the same Spawn helper; not in this cycle."
- README router template patch — issue §Out-of-scope, "addresses non-cn bodies; --claude / --codex is the cn-installed path."

The `Spawn` / `CheckSpawnBinary` helpers are binary-name-parameterized, so adding `--cursor` later is a flag-parsing addition in `cmd_activate.go` plus a `<cursor's interactive-launch shape>` verification — no re-architecting required.

**Branch CI.** Local Go test suite green at HEAD `87aa69e9`; full repo CI status to be verified by β against `origin/cycle/380` per pre-review gate row 10. Recorded as a transient row in §Review-readiness, re-validated immediately before the readiness signal.

## CDD Trace

CDD canonical artifact order (CDD.md §5.2):

1. **Design artifact** — not required. Mode is design-and-build per γ scaffold §Mode; design surface is the issue body. The one bounded design call γ delegated (render-capture seam Option A vs B) was resolved by α inline; declared in §Gap and §ACs AC1 implementation notes.

2. **Coherence contract** — `.cdd/unreleased/380/self-coherence.md` §Gap (this file). The named gap (cn activate's hub-side flow stops at stdout; no bridge to interactive REPL) and the closure shape (two new flags + Spawn helper + render-capture seam Option A) are stated.

3. **Plan** — not required. Sequencing is the commit-checkpoint sequence γ scaffolded:

| Step | Surface | Commit |
|---|---|---|
| 1 | `cmd_activate.go` — flag defs + help text | `1cc80565` |
| 2 | `spawn.go` + `spawn_other.go` — Spawn / CheckSpawnBinary helpers | `d4f4e499` |
| 3 | `cmd_activate.go` — mutual-exclusion + render-capture seam + spawn call | `e92a9476` |
| 4 | `spawn_test.go` + `cmd_activate_test.go` — flag/argv/missing-binary tests | `da3dd429` |
| 5 | `cmd_activate_test.go` — AC3 bytes-equal default-path oracle | `87aa69e9` |
| 6 | `.cdd/unreleased/380/self-coherence.md` — §Gap..§CDD Trace | `808ae8b7..` (this commit family) |

The first implementation commit (`1cc80565`) landed within the first 25% of the dispatch budget per γ scaffold authoring discipline. No batched end-of-cycle commit.

4. **Tests** — present.

| Test file | Test count | What it proves |
|---|---|---|
| `src/go/internal/activate/spawn_test.go` | 7 | AC1 argv shape; AC2 codex argv shape (negative: not `[codex, exec, prompt]`); LookPath failure path; CheckSpawnBinary error message shape for claude/codex; default-hook wiring verifies production `Spawn` routes through `exec.LookPath` |
| `src/go/internal/cli/cmd_activate_test.go` | 9 | AC3 default no-flag preserved + bytes-equal to direct activate.Run; AC4 mutual exclusion (errors + names both flags + fires before render); AC5 missing-binary for claude/codex (names binary + flag + PATH + fires before render); AC1/AC2 --help documents both flags; unknown-flag rejection |
| `src/go/internal/activate/activate_test.go` | unchanged | All 3.78.0 / #379 invariants preserved (full suite green) |

```
$ cd src/go && go test ./internal/activate/... ./internal/cli/...
ok      github.com/usurobor/cnos/src/go/internal/activate    0.040s
ok      github.com/usurobor/cnos/src/go/internal/cli         0.019s

$ go test ./...
ok      github.com/usurobor/cnos/src/go/internal/activate
ok      github.com/usurobor/cnos/src/go/internal/activation
ok      github.com/usurobor/cnos/src/go/internal/binupdate
ok      github.com/usurobor/cnos/src/go/internal/cli
ok      github.com/usurobor/cnos/src/go/internal/discover
ok      github.com/usurobor/cnos/src/go/internal/dispatch
ok      github.com/usurobor/cnos/src/go/internal/doctor
ok      github.com/usurobor/cnos/src/go/internal/hubinit
ok      github.com/usurobor/cnos/src/go/internal/hubsetup
ok      github.com/usurobor/cnos/src/go/internal/hubstatus
ok      github.com/usurobor/cnos/src/go/internal/pkg
ok      github.com/usurobor/cnos/src/go/internal/pkgbuild
ok      github.com/usurobor/cnos/src/go/internal/restore
```

`go vet ./...` clean.

5. **Code** — `spawn.go` (70 lines), `spawn_other.go` (19 lines), `cmd_activate.go` (modified — 92 lines / +47 / -5).

6. **Docs** — `cn activate --help` text updated in `cmd_activate.go` to document `--claude` and `--codex` with description naming the interactive-spawn behavior, TTY preservation, and PATH requirement. No README / RELEASE.md / skill body edits needed this cycle. Every file in `git diff --stat origin/main..HEAD` (per pre-review gate row 11) is accounted for:

| File | Change | Where named in self-coherence |
|---|---|---|
| `.cdd/unreleased/380/gamma-scaffold.md` | A | γ artifact, named in §Gap (base SHA) and §Self-check (target-surface comparison) |
| `src/go/internal/activate/spawn.go` | A | §ACs AC1 surface |
| `src/go/internal/activate/spawn_other.go` | A | §Self-check (portability stub justification) |
| `src/go/internal/activate/spawn_test.go` | A | §ACs AC1 / AC2 oracles; §CDD Trace step 4 |
| `src/go/internal/cli/cmd_activate.go` | M | §ACs AC1 / AC3 / AC4 surface; §Gap (Option A) |
| `src/go/internal/cli/cmd_activate_test.go` | A | §ACs AC3 / AC4 / AC5 oracles; §Self-check (surface widening) |
| `.cdd/unreleased/380/self-coherence.md` | A | this file |

Non-test caller-path trace (pre-review gate row 12): every new exported symbol has a non-test caller.

- `activate.CheckSpawnBinary` — called by `cli/cmd_activate.go` `ActivateCmd.Run` (pre-render gate)
- `activate.Spawn` — called by `cli/cmd_activate.go` `ActivateCmd.Run` (post-render handover)

Test assertion count from runner output (pre-review gate row 13): 16 new test functions land green (7 in spawn_test.go + 9 in cmd_activate_test.go), enumerated above. Previous-cycle invariants still pass: all activate package tests including #379 ordering, kernel-state, deps-state, latest-reflection, fixture-based smoke.

7. **Self-coherence** — this file. Authored section-by-section per alpha/SKILL.md §2.5 (each section a separate commit + push). Section manifest in HTML-comment header records completion state.

## Review-readiness

**Round 1.** Implementation SHA: `87aa69e9` (last implementation commit before the self-coherence file family). Branch HEAD will advance with this readiness-signal commit; β should poll `origin/cycle/380` for the latest HEAD.

**Base SHA:** `319893a4ea9de1b89989ef6e6dc44cd3e1cad147` (matches the dispatch base; `origin/main` HEAD at observation time below).

**Pre-review gate (alpha/SKILL.md §2.6) — observation time 2026-05-19:**

| # | Row | State |
|---|---|---|
| 1 | cycle branch rebased onto `origin/main` | ✓ `origin/main` = `319893a4`; cycle/380 has zero commits behind `git log cycle/380..origin/main` is empty (re-validated immediately before signal) |
| 2 | CDD Trace through step 7 | ✓ §CDD Trace |
| 3 | tests present | ✓ 16 new tests + full pre-existing activate suite green |
| 4 | every AC has evidence | ✓ §ACs (AC1–AC5) |
| 5 | known debt explicit | ✓ §Debt |
| 6 | schema audit when contracts changed | N/A — no schema-bearing contract changed this cycle |
| 7 | peer enumeration | ✓ §Self-check (flag, binary, renderer, help-text peers) |
| 8 | harness audit when schema changed | N/A — no schema-bearing contract changed this cycle |
| 9 | polyglot re-audit | ✓ diff languages: Go (run `go vet ./... && go test ./...` — both green) + Markdown (this file + γ scaffold). No shell, YAML, or non-Go writers of any new contract. |
| 10 | branch CI green | local Go suite green at HEAD; full CI verification falls to β at merge time (transient row, β to re-validate against `origin/cycle/380` HEAD) |
| 11 | artifact enumeration matches diff | ✓ §CDD Trace step 6 table accounts for every file in `git diff --stat origin/main..HEAD` |
| 12 | caller-path trace for new exported symbols | ✓ §CDD Trace step 6 names non-test callers for `activate.Spawn` and `activate.CheckSpawnBinary` (both in `cli/cmd_activate.go ActivateCmd.Run`) |
| 13 | test assertion count from runner output | ✓ §CDD Trace step 4 — 7 in spawn_test.go, 9 in cmd_activate_test.go, runner output pasted |
| 14 | α commit author email canonical | ✓ all 11 α commits (range `7a9bc2e7..HEAD` excluding γ's scaffold) authored as `alpha@cdd.cnos` after path-(a) `git rebase --exec 'git commit --amend --reset-author --no-edit' 7a9bc2e7 && git push --force-with-lease`; the worktree config was discovered carrying `gamma@cdd.cnos` from prior session, fixed via `git config --worktree user.email "alpha@cdd.cnos"` before the rebase. Verification: `git log --format='%ae' 7a9bc2e7..HEAD | sort -u` returns the single value `alpha@cdd.cnos`. |

**Ready for β.** All AC1–AC5 oracles pass; surface containment matches the γ scaffold plus two declared widenings (`spawn_other.go` portability stub + `cmd_activate_test.go` for cli-package-level oracles); no non-goals violated; render-capture seam Option A declared and verified.

## Fix-round 1

**Trigger.** β R1 verdict REQUEST CHANGES in `.cdd/unreleased/380/beta-review.md` — one B-severity honest-claim finding (F1).

**Fix-round head SHA:** `91a2cad6` (single fix commit; readiness-signal commit will advance HEAD by one).
**Cycle head at fix-round signal:** β should poll `origin/cycle/380` for the latest HEAD.

### F1 — §CDD Trace SHAs unreachable from `origin/cycle/380`

**β's finding.** Five SHAs in §CDD Trace step table (`44b9a475`, `c28a04ee`, `304eb1df`, `20d07860`, `dd86e283`) plus the §ACs header reference to `dd86e283` were pre-rebase SHAs from α's local reflog after the path-(a) email-rewrite rebase documented in §Review-readiness row 14. None reachable from `origin/cycle/380`; honest-claim reproducibility (review/SKILL.md rule 3.13(a)) violated.

**Addressed by.** Commit `91a2cad6` "α #380 fix F1 — rewrite §CDD Trace SHAs to current-branch values".

**What was rewritten.** Per α/SKILL.md §2.3 intra-doc repetition discipline, β named 5 SHAs + 1 header reference but the doc carried more occurrences — grepped `44b9a475|c28a04ee|304eb1df|20d07860|dd86e283|14a6bf55` against the file and rewrote every hit (9 occurrences total):

| # | Line | Old SHA (pre-rebase, unreachable) | New SHA (current-branch) | Context |
|---|---|---|---|---|
| 1 | 43 | `dd86e283` | `87aa69e9` | §ACs header: "Per-AC oracles run against branch HEAD …" |
| 2 | 232 | `dd86e283` | `87aa69e9` | §Debt: "Local Go test suite green at HEAD …" |
| 3 | 246 | `44b9a475` | `1cc80565` | §CDD Trace step 1 row |
| 4 | 247 | `c28a04ee` | `d4f4e499` | §CDD Trace step 2 row |
| 5 | 248 | `304eb1df` | `e92a9476` | §CDD Trace step 3 row |
| 6 | 249 | `20d07860` | `da3dd429` | §CDD Trace step 4 row |
| 7 | 250 | `dd86e283` | `87aa69e9` | §CDD Trace step 5 row |
| 8 | 251 | `14a6bf55..` | `808ae8b7..` | §CDD Trace step 6 (self-coherence family start: §Gap commit) |
| 9 | 253 | `44b9a475` | `1cc80565` | §CDD Trace prose: "The first implementation commit (…) landed within the first 25%" |

β surfaced 6 (rows 1, 3–7); the intra-doc repetition rule (`alpha/SKILL.md` §2.3 — "grep the doc for every occurrence … before claiming the fix is complete"; derives from #266 F3-bis) caught rows 2, 8, 9 in the same commit so this fix-round does not surface an F1-bis.

**Verification for β to re-run.**

```
$ git log --oneline 319893a4..origin/cycle/380
# Expect: each new SHA in the table above reachable; old SHAs unreachable.

$ for sha in 44b9a475 c28a04ee 304eb1df 20d07860 dd86e283 14a6bf55; do \
    echo "=== $sha ==="; git branch -a --contains "$sha" 2>&1; \
  done
# Expect: empty / "malformed object name" for each — none reachable from any branch.

$ for sha in 1cc80565 d4f4e499 e92a9476 da3dd429 87aa69e9 808ae8b7; do \
    echo "=== $sha ==="; git branch -a --contains "$sha" 2>&1 | head -3; \
  done
# Expect: each contains "cycle/380" and "remotes/origin/cycle/380".

$ grep -n -E '44b9a475|c28a04ee|304eb1df|20d07860|dd86e283|14a6bf55' \
    .cdd/unreleased/380/self-coherence.md
# Expect: no matches.
```

**No production-code change.** F1 is a self-coherence artifact correction only; no Go diff, no test change, no skill body change. AC1–AC5 evidence is unchanged — the oracles still pass against the renamed SHAs (the underlying commits are identical content, just authored by the canonical `alpha@cdd.cnos` email post-rebase).

### Re-audit against HEAD (α/SKILL.md §3.4)

Re-ran the affected pre-review-gate rows against new HEAD `91a2cad6`:

| Row | State | Notes |
|---|---|---|
| 2 CDD Trace through step 7 | ✓ | §CDD Trace step table SHAs now current-branch-reachable |
| 4 every AC has evidence | ✓ | §ACs header SHA now points to a reachable commit |
| 9 polyglot re-audit | ✓ | No new languages touched; only Markdown changed |
| 11 artifact enumeration matches diff | ✓ | `git diff --stat origin/main..HEAD` adds no new files; only `.cdd/unreleased/380/self-coherence.md` modified twice (fix commit + this fix-round append) |
| 13 honest-claim reproducibility (review/SKILL.md 3.13(a)) | ✓ | β's grep recipe in the verification block above will return empty for old SHAs and confirm reachability for new SHAs |

Other rows (1, 3, 5–8, 10, 12, 14) untouched by this fix and remain as recorded in §Review-readiness.

### Ready for β round 2

F1 fix is one self-coherence-only commit (`91a2cad6`) plus this fix-round append. No code change, no test change. β's expected R2 verdict per beta-review.md "Round budget" note: APPROVE + merge.
