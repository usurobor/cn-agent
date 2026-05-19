<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills, acs]
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

Per-AC oracles run against branch HEAD `dd86e283` (last implementation commit before this self-coherence file). AC numbering follows the issue body.

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
