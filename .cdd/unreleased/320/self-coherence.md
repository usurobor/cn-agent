## Gap — #320 cn activate

**Issue:** #320 — feat(cli): cn activate — bootstrap prompt generation from hub state

**Version/mode:** MCA — adds `cn activate` as a kernel CLI command that generates a bootstrap prompt from local hub state. No runtime is started; no model is invoked.

**Incoherence closed:** `cn init` and `cn setup` exist, but no CLI path helps the operator assemble the activation prompt a capable model needs to orient to the hub. Operators must manually know which hub files matter and assemble the prompt themselves, producing inconsistent results.

---

## Skills

**Tier 1 (always loaded):**
- `CDD.md` (canonical lifecycle)
- `cdd/alpha/SKILL.md` (α role surface)

**Tier 2 (general engineering):**
- `cnos.core/skills/design/SKILL.md`
- `cnos.core/skills/write/SKILL.md`
- `cnos.eng/skills/eng/tool/SKILL.md`
- `cnos.eng/skills/eng/test/SKILL.md`
- `cnos.eng/skills/eng/ux-cli/SKILL.md`
- `cnos.eng/skills/eng/go/SKILL.md`

**Active constraints applied:**
- `design`: prompt generation separated from runtime/daemon/dispatch; `cn activate` is a flat kernel command (not `cn agent activate`); `cli/` owns dispatch only (`eng/go §2.18`)
- `write`: help text and prompt prose say what the command does without overclaiming runtime behavior
- `tool`: stdout = prompt only; stderr = diagnostics; fail-fast on invalid hub
- `test`: invariants named before tests written; positive and negative cases both covered; no model invocation in tests
- `ux-cli`: error messages are causal and actionable; `✗` and `→` glyphs used per semantic table
- `go`: `eng/go §2.17` (Parse/Read split not applicable — no schema parsing introduced); `eng/go §2.18` (domain logic in `internal/activate/`, thin wrapper in `internal/cli/cmd_activate.go`)

---

## ACs

**AC1: Command exists — `cn activate` prints a bootstrap prompt**
- Evidence: `ActivateCmd` registered in `cmd/cn/main.go:43`; `Run()` calls `activate.Run()` which writes prompt to stdout.
- Test: `TestRunPositive_CwdHub` — verifies prompt header in stdout; `TestRunPositive_StdoutOnly` — verifies stdout contains prompt.
- ✅ Met.

**AC2: Flat kernel command, not agent namespace**
- Evidence: `ActivateCmd.Spec()` returns `Source: SourceKernel, Tier: TierKernel`. No `cn agent` command introduced. `git grep -r "cn agent" src/go/` → 0 hits.
- Test: `TestRunPositive_CwdHub` exercises the command directly.
- ✅ Met.

**AC3: Hub discovery works from cwd**
- Evidence: `cmd_activate.go` passes `inv.HubPath` (cwd-discovered in `main.go:discoverHub()`) when no explicit arg is given. `activate.Run()` accepts it and validates `.cn/`.
- Test: `TestRunPositive_CwdHub` passes `HubPath: hub` (simulating cwd discovery result); `TestRunNegative_EmptyHubPath` confirms failure when no hub is found.
- ✅ Met.

**AC4: Explicit local HUB_DIR works**
- Evidence: `cmd_activate.go:Run()` — if `inv.Args[0]` is non-flag it overrides `inv.HubPath`; `activate.Run()` validates the path.
- Test: `TestRunPositive_ExplicitHubDir` passes an explicit path.
- ✅ Met.

**AC5: Command does not require pre-discovered hub**
- Evidence: `CommandSpec.NeedsHub: false` — `main.go` hub gate does not reject the command when `hubPath == ""`. `cmd_activate.go` resolves hub itself.
- Test: `TestRunNegative_EmptyHubPath` confirms the command runs (and fails with a diagnostic) even without a hub.
- ✅ Met.

**AC6: Prompt uses safe hub state**
- Evidence: `activate.go:Run()` reads only `.cn/config.json` (name/version/created). It does not read `.cn/secrets.env`, `.env`, or any token-containing file. Prompt body contains path/structure references only — no file contents.
- Test: `TestRunPositive_SecretsExcluded` — writes fake secret to `.cn/secrets.env`, verifies secret content absent from stdout; verifies "secrets.env" mentioned in Notes.
- ✅ Met.

**AC7: Prompt is model-ready but does not call a model**
- Evidence: `activate.go` has zero `os/exec` or HTTP calls. `writePrompt()` only writes to `io.Writer`.
- Test: `TestRunPositive_NoModelInvocation` — verifies Run completes and prompt states "No model is invoked."
- ✅ Met.

**AC8: stdout/stderr contract preserved**
- Evidence: `writePrompt()` writes only to `opts.Stdout`; diagnostic (`→ Generating...`) writes only to `opts.Stderr`; error messages write only to `opts.Stderr`.
- Test: `TestRunPositive_StdoutOnly` — checks `→` diagnostic is in stderr not stdout; checks prompt is in stdout.
- ✅ Met.

**AC9: Claude CLI example is valid**
- Evidence: `activateHelp` constant and `writePrompt()` Notes section both use `claude -p "Activate this cnos hub using the bootstrap prompt on stdin."` — includes a query argument per documented `claude -p` shape. Neither uses bare `claude -p` without a query.
- Test: `TestRunPositive_ClaudeCliExampleInPrompt` — verifies the example string appears in stdout.
- ✅ Met.

**AC10: Docs preserve current-vs-target truth**
- Evidence: `activateHelp` says "generates a prompt only. No model is invoked." `writePrompt()` Notes say "This command generates a prompt only. No model is invoked." No claim of daemon, scheduler, or runtime is made.
- Test: `TestRunPositive_NoModelInvocation` checks the explicit statement.
- ✅ Met.

**AC11: Namespace migration not part of MVP**
- Evidence: No `cn agent` command added. No existing command renamed or replaced. `git diff main...cycle/320 -- src/go/cmd/cn/main.go` shows only one new `reg.Register(&cli.ActivateCmd{})` line added.
- ✅ Met.

**AC12: Tests cover positive and negative cases**
Positive cases covered:
- `TestRunPositive_CwdHub` — cwd hub discovery succeeds, prompt produced
- `TestRunPositive_ExplicitHubDir` — explicit HubPath succeeds
- `TestRunPositive_PackagesListed` — packages in prompt
- `TestRunPositive_SecretsExcluded` — secret content absent from stdout
- `TestRunPositive_StdoutOnly` — stdout = prompt, stderr = diagnostics
- `TestRunPositive_NoModelInvocation` — no model called, prompt states it
- `TestRunPositive_ClaudeCliExampleInPrompt` — valid Claude CLI example

Negative cases covered:
- `TestRunNegative_EmptyHubPath` — no hub fails with diagnostic
- `TestRunNegative_MissingPath` — nonexistent path fails
- `TestRunNegative_PathWithoutDotCn` — directory without `.cn/` fails
- ✅ Met.
