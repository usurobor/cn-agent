# 3.79.0

## Outcome

Coherence delta: C_Σ **B−** (`α B+`, `β A`, `γ C`) · **Level:** `L6`

`cn activate HUB_DIR` no longer stops at "prompt rendered" — it can hand the operator into a live interactive AI-CLI session with the body already activated. `cn activate --claude HUB_DIR` and `cn activate --codex HUB_DIR` replace the cn process with `claude "$prompt"` or `codex "$prompt"` via `syscall.Exec`, preserving the TTY cleanly. One command bootstraps a body session on a workstation with the AI CLI installed. Default no-flag stdout behavior is preserved bytes-equal to 3.78.0 by construction (Option A render-capture seam: the default arm has no buffer indirection, so the bytes-equal property holds for AC3 even if the unit oracle were absent).

## Why it matters

3.78.0 shipped the activation prompt as a body-agnostic skill artifact and `cn activate` as the renderer; that closed the "body cannot self-bootstrap without a fetchable procedure" gap at the protocol layer. But the canonical local-dev flow (workstation with `claude` or `codex` installed, plus `cn`) still required three steps: render the prompt, start a fresh CLI session, paste the prompt. The pipe form (`cn activate cn-sigma | claude`) was non-interactive — the prompt printed and the session exited. The bridge from "prompt rendered" to "body operating interactively at this hub" was missing, and the cn-installed local-dev path was less ergonomic than the URL-based router path the activate skill names for non-cn bodies.

This release closes that gap. The activation skill is unchanged; the renderer is unchanged. The new surface is a `cn`-side concern: two flags, a spawn helper, and a portability stub for non-unix. Pre-render validation (mutual exclusion, missing-binary detection) fires before any render call — proven by a "nonexistent HUB_DIR" oracle that asserts both the presence of the pre-render diagnostic AND the *absence* of the renderer's own error in stderr. Tests assert exec shape via injection hooks (no real `claude` / `codex` required on the test machine) and explicitly disprove the `codex exec PROMPT` non-interactive-subcommand drift.

## Added

- **#380** — `src/go/internal/activate/spawn.go` (NEW, `//go:build unix`): exports `Spawn(binary, prompt)` (uses `syscall.Exec` for TTY-clean process replacement) and `CheckSpawnBinary(binary, flag)` (LookPath-based pre-render check with an actionable missing-binary error naming the binary, the requesting flag, the PATH variable, and an install hint).
- **#380** — `src/go/internal/activate/spawn_other.go` (NEW): non-unix build-tag stub returning `<flag> spawn is not supported on this platform — use the default stdout form …`. The no-flag default path remains functional everywhere.
- **#380** — `src/go/internal/activate/spawn_test.go` (NEW, 7 tests): argv-shape assertions for both binaries (`TestSpawnWith_CodexArgvShape` explicitly asserts `argv[1] != "exec"` to disprove the non-interactive-subcommand drift); LookPath failure path; missing-binary error message shape; default-hook wiring (`TestSpawn_DefaultsAreWired` verifies the production path actually routes through `exec.LookPath`).
- **#380** — `src/go/internal/cli/cmd_activate.go`: adds `--claude` and `--codex` flags; mutual-exclusion check fires before hub resolution and before any `activate.Run` call; LookPath check fires pre-render; spawn arm captures into `bytes.Buffer` only on spawn paths; default arm passes `inv.Stdout` straight through (Option A render-capture seam).
- **#380** — `src/go/internal/cli/cmd_activate_test.go` (NEW, 9 tests): `TestActivate_MutualExclusion_FiresBeforeRender`, `TestActivate_MissingBinary_FiresBeforeRender` (both prove ordering via nonexistent HUB_DIR — stderr contains the pre-render diagnostic AND does NOT contain `Hub path not found` AND does NOT contain `Generating activation prompt`); `TestActivate_DefaultNoFlag_BytesEqualToDirectRun` (AC3 unit oracle); `TestActivate_HelpFlag_DocumentsClaudeAndCodex` (help-text contract).

## Changed

- **#380** — `cn activate --help` text documents the new `--claude` and `--codex` flags with interactive-REPL behavior, TTY preservation, and PATH requirement.
- **`src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`** §2.6 — added an auxiliary paragraph after row 14 naming SHA-citation invalidation as a downstream consequence of path (a) rebase (`git rebase --exec 'git commit --amend --reset-author --no-edit'`). The paragraph prescribes two coherent resolution paths: (i) run row 14's identity verification at session-start, before authoring any SHA-bearing artifact (preferred); (ii) re-stamp every SHA citation in `self-coherence.md` and other authored artifacts to the current-branch SHA immediately after the rebase, applying §2.3's intra-doc repetition rule. Empirical anchor `#380 R1 F1` cited inline. (γ step 13a; lands in this release.)

## Validation

- **Default-path bytes-equal (AC3, structural):** `go test -count=1 -run TestActivate_DefaultNoFlag_BytesEqualToDirectRun ./internal/cli/...` PASS on cycle HEAD and merge tree (β R1 + R2). `git diff origin/main..origin/cycle/380 -- src/go/internal/activate/activate.go` empty — Option A render-capture seam makes the property structural, not test-only.
- **Argv shape (AC1, AC2):** `go test -count=1 -run TestSpawnWith ./internal/activate/...` PASS — both `TestSpawnWith_ClaudeArgvShape` and `TestSpawnWith_CodexArgvShape` (codex test explicitly disproves `argv[1] == "exec"`).
- **Pre-render ordering (AC4, AC5):** `go test -count=1 -run TestActivate_(MutualExclusion|MissingBinary)_FiresBeforeRender ./internal/cli/...` PASS — proven by stderr-does-NOT-contain assertions against the renderer's own error path.
- **Help-text contract:** `TestActivate_HelpFlag_DocumentsClaudeAndCodex` PASS (asserts both flag names + the literal substring `interactive`).
- **End-to-end interactive bootstrap (manual dry-run):** deferred — no workstation session with `claude` installed available in this PRA. Commit on next operator session checklist.

## Known Issues

- `Repo link validation (I4)` CI workflow remains red on `main` — lychee fails on file:// links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` (cycle #369's archived artifact). Carried forward from 3.78.0; this release does not regress or fix. `notify` job is downstream of Build success and stays red for the same reason. The next cnos-internal MCA remediates I4 (see PRA §7 next-move list and gamma-closeout §Deferred Outputs).
- Hub README router template ships in the 3.78.0 activate skill (§2.3) but no hub has adopted it yet; per-hub README patches are deferred to per-hub cycles (carried-forward MCA: cn-sigma adoption first, cn-pi second).
- Renderer fallback `canonicalReadFirstOrdering` duplicates the activate skill's §4.1 ordering as in-Go constants (3.78.0 carry-forward). Durable fix is `//go:embed` of the skill into the binary; deferred to a future cycle.
- `--cursor` / `--aider` / other AI CLI flags held by explicit non-goal in #380 issue body until friction emerges. `--auto` body-detection flag + `$CN_DEFAULT_BODY` env-var default same disposition.
- Worktree-config identity-leak class (`extensions.worktreeConfig=true` + non-`--worktree` writes) confirmed across three consecutive cycles in three role surfaces (#379, #370, #380); issue #373 (Preventive `--worktree` identity write across role skills) carries the converged design at escalating lag.

## Source

Issue #380 (https://github.com/usurobor/cnos/issues/380). Operator-direct request during the cnos 3.78.0 disconnect phase. Anchor: cnos δ session of 2026-05-19 (cn-sigma daily reflection `threads/reflections/daily/20260519.md`).
