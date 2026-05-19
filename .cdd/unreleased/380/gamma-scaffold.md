---
cycle: 380
issue: "https://github.com/usurobor/cnos/issues/380"
date: "2026-05-19"
dispatch_configuration: "§5.1 canonical multi-session (one `claude -p` per role) — continuity with cycle/379 same-day infra; §5.3 criteria do not strictly require §5.1 (5 ACs, single Go package), but δ continues §5.1 for tooling-loop coherence with the just-shipped 3.78.0 cycle"
base_sha: "319893a4ea9de1b89989ef6e6dc44cd3e1cad147"
mode: "design-and-build"
mode_rationale: "Design converges in the issue body during γ (flag names --claude / --codex, exec mechanism syscall.Exec, mutual-exclusion ordering, pre-render LookPath, render-capture seam Option A vs B). α implements + tests in the same cycle. Not MCA — no separately-committed DESIGN.md path; the issue body is the design surface, and the only ambiguity (render-capture seam) is a small bounded design call γ delegates to α with both options pre-validated as acceptable to β."
---

# γ Scaffold — #380

## Gap

cnos 3.78.0 shipped the source-of-truth half of activation: `cnos.core/skills/agent/activate/SKILL.md` is the body-agnostic skill, and `src/go/internal/activate/activate.go` renders the prompt from it. The hub-side flow stops short: `cn activate HUB_DIR` writes the rendered prompt to stdout. The intended pattern was `cn activate cn-sigma | claude`, but `claude` (and `codex`) in their non-interactive pipe form are one-shot — they print the model's response and exit. The operator does not land in an interactive session. To recover interactive mode, the operator copy-pastes the prompt into a fresh `claude` / `codex` invocation. The bridge from "prompt rendered" to "body operating interactively at this hub" is missing.

This cycle adds the spawn-into-CLI surface on top of the 3.78.0 renderer. Two flags — `--claude` and `--codex` — replace the cn process with the chosen AI CLI invoked bare-positionally (`claude "$prompt"` / `codex "$prompt"`), preserving the TTY so the operator lands in a live REPL with the body already in the post-activation state. Default no-flag behavior is preserved bytes-equal — the pipe-and-redirect path for advanced consumers must not regress.

The 3.78.0 skill (`agent/activate/SKILL.md`) is body-agnostic by design and stays untouched. The 3.78.0 renderer (`activate.go`) stays untouched (Option A) or grows one render-only entrypoint (Option B); γ delegates the seam choice to α with both pre-validated as β-acceptable. The full surface this cycle adds is one new helper file (`spawn.go`), its test file (`spawn_test.go`), and three new flag-handling lines plus the spawn wiring in `cmd_activate.go`.

## Peer enumeration (gamma/SKILL.md §2.2a)

Before authoring §Gap, γ enumerated the directories named by the issue's impact graph and grepped for the proposed surfaces:

```bash
# (a) The target Go package family — confirms what exists today
ls src/go/internal/activate/
# → activate.go (3.78.0 renderer), activate_test.go (existing tests)
#   spawn.go and spawn_test.go do NOT exist — this cycle is additive at canonical paths

ls src/go/internal/cli/cmd_activate*
# → src/go/internal/cli/cmd_activate.go (current CLI entry, ~64 lines, single ActivateCmd struct)

# (b) Search for any pre-existing --claude / --codex / spawn surface in src/go
rg -l '\-\-claude|\-\-codex|syscall\.Exec|spawn' src/go
# → src/go/internal/dispatch/types.go (cn dispatch backend identifiers, not cn activate)
#   src/go/internal/cli/cmd_dispatch.go (cn dispatch CLI, not cn activate)
#   NO matches in src/go/internal/activate/ or src/go/internal/cli/cmd_activate.go
# → confirms: this cycle is additive at canonical paths; no consolidation with cn dispatch

# (c) Pattern reference for exec.LookPath in this repo (issue's AC5 mechanism)
rg -n 'LookPath' src/go
# → src/go/internal/dispatch/backend_claude.go:18 — `if _, err := exec.LookPath("claude"); err != nil`
#   pattern this cycle reuses for AC5 (but the spawn path itself uses syscall.Exec, not os/exec.Cmd)

# (d) The 3.78.0 surfaces this cycle builds on (binding non-modify targets)
ls src/packages/cnos.core/skills/agent/activate/SKILL.md
# → present (3.78.0 ship); body-agnostic; non-goal for this cycle to modify

grep -n 'writePrompt\|loadActivateSkillOrdering' src/go/internal/activate/activate.go
# → 70: writePrompt(opts.Stdout, …) — Option A reuses unchanged
#   65: readFirst, source := loadActivateSkillOrdering(cnDir) — Option B's Render entrypoint extracts above this

# (e) CLI claim verification (issue says: claude "PROMPT" + codex "PROMPT" are interactive launches)
claude --help | head -5
# → "Usage: claude [options] [command] [prompt]
#    Claude Code - starts an interactive session by default, use -p/--print for non-interactive output"
# → confirms: bare-positional `claude "PROMPT"` lands in interactive REPL

codex --help | head -5
# → "Codex CLI
#    If no subcommand is specified, options will be forwarded to the interactive CLI."
# → confirms: bare-positional `codex "PROMPT"` is the canonical interactive launch
#   ALSO confirms: `codex exec "PROMPT"` is the explicit non-interactive variant — α must NOT use this
```

All five checks confirm: gap is real, additive at canonical paths, no consolidation required. The issue's CLI claims are empirically verified by γ. The pattern-reference for `exec.LookPath` is already in the repo (`backend_claude.go`); the only net-new mechanism is the `syscall.Exec` call site, which has no precedent in this repo's Go code.

## Mode

**design-and-build.** Design converges in the issue body during γ:
- flag names (`--claude`, `--codex`) — issue body §Scope
- exec mechanism (`syscall.Exec`, not `os/exec.Cmd`) — issue body §Active design constraints
- bare-positional invocation (`[binary, prompt]`, no `-p` or `exec` subcommand) — issue body §Active design constraints + γ-verified via `claude --help` / `codex --help`
- mutual-exclusion ordering (pre-render) — issue body AC4 + γ tightening in α prompt
- missing-binary detection (pre-render via `exec.LookPath`) — issue body AC5 + γ tightening in α prompt
- render-capture seam — one bounded design call γ delegates to α with both options pre-validated (see α prompt §Render-capture seam)

α executes both the design step (pick render-capture seam option, declare in self-coherence.md) and the build step (write spawn.go, wire cmd_activate.go, write tests) inside this cycle.

Mode is **not MCA**: no separate committed `DESIGN.md` / `PLAN.md` at a stable path. The issue body is the design surface; the source-of-truth rows cite the artifact paths the cycle will produce/evolve, not pre-existing design docs.

## ACs reference

All 5 ACs are defined in the issue body at https://github.com/usurobor/cnos/issues/380 — they are the binding contract for α/β. α maps each AC to evidence in `.cdd/unreleased/380/self-coherence.md` §ACs.

## Acceptance posture summary (γ pre-flagged to α)

The issue body is authoritative. The notes below are γ's compressed reading of each AC for handoff; they do not override the issue.

- **AC1 — `--claude` spawns interactively.** Flag is added to the Cobra-style command spec in `cmd_activate.go`; documented in `cn activate --help`. When set, after rendering, the cn process is replaced by `claude "$prompt"` via `syscall.Exec`. argv shape (after argv[0]) is exactly `[claude, prompt]` — no `-p`, no `--print`, no other flags. TTY/stdin/stdout inherited from the parent terminal (this is what makes the REPL handoff work). Testability: the syscall.Exec call goes through a test-injectable hook so spawn_test.go can record the planned argv without actually replacing the process.

- **AC2 — `--codex` spawns interactively.** Identical shape to AC1, exec-ing `codex "$prompt"`. **Critical:** the argv after argv[0] is exactly `[codex, prompt]`, NOT `[codex, exec, prompt]` — `codex exec` is the non-interactive subcommand and would defeat the cycle's purpose. γ verified via `codex --help`: "If no subcommand is specified, options will be forwarded to the interactive CLI."

- **AC3 — default `cn activate HUB_DIR` (no flag) preserved bytes-equal.** The hardest invariant to keep cheaply: any diagnostic written to stdout (a header, a deprecation notice, a "use --claude" hint) breaks pipe consumers. The render-capture seam (Option A or B) must NOT redirect the no-flag path through a buffer; the default path keeps `inv.Stdout` direct. β diffs against a captured 3.78.0 fixture or a worktree-checkout smoke compare.

- **AC4 — `--claude` + `--codex` mutual exclusion before render.** Setting both flags exits non-zero with a stderr message naming both flags and the conflict, BEFORE any rendering. The order-of-checks invariant is verifiable by setting HUB_DIR to a path that does not exist; the mutual-exclusion error must fire, not the "hub path not found" error from `activate.Run`. γ added this verification trick to the α/β prompts because the issue body alone does not nail the order of checks.

- **AC5 — missing-binary error is clear and pre-render.** With `--claude` set and `claude` not on `$PATH`, exit non-zero. Error names both the missing binary (`claude`) AND the flag that requested it (`--claude`). Detection via `exec.LookPath`, called BEFORE `activate.Run`. Same order-of-checks proof as AC4: with a non-existent HUB_DIR + missing binary, the LookPath error fires first, not the hub-path-not-found error. (Issue's §Active design constraints says "pre-exec"; γ tightens to "pre-render" to avoid the stderr `→ Generating activation prompt…` diagnostic firing before the error confuses the operator.)

## Skills to load

**Tier 1 (always-on for α/β CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α) / `beta/SKILL.md` (β)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown + Go authoring
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific, per issue body §Skills to load):**
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go authoring discipline for the new flag wiring + exec path + tests (load-bearing for AC1–AC5)
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any in-cycle issue-pack reconciliation

The issue intentionally lists only two Tier 3 skills (single-focus Go-only cycle, no doctrine touch). `skill/SKILL.md` is not loaded — no SKILL.md is authored this cycle. `design/SKILL.md` is not loaded — the design surface is small and converges in the issue body, not a separate design step.

## Dispatch configuration

**§5.1 canonical multi-session.** γ produces α and β prompts at `/root/dispatch-380/{alpha,beta}-prompt.md` (NOT `/tmp/dispatch-380/` — see Dispatch-dir hygiene note below); δ (operator) routes each role to a fresh `claude -p` session sequentially. No Agent-tool sub-agent dispatch from a parent session. Each role loads its skills fresh from disk on dispatch.

**§5.3 escalation rationale.** Per `operator/SKILL.md` §5.3, the strict criteria for §5.1 over §5.2 are: ≥3 ACs with new contract surface AND cross-repo proposal AND active doctrine evolution. cycle/380 has 5 ACs and a new contract surface (new helper file + new flag wiring) but NO cross-repo proposal and NO doctrine evolution. §5.2 (γ-in-hub) would technically suffice. δ continues §5.1 for **continuity with cycle/379** (just shipped at 3.78.0 the same day, via §5.1) — keeping the dispatch infrastructure stable across the 3.78.0 → 3.79.0 disconnect avoids re-exercising §5.2 wiring under a hot tooling loop. The operator's same-day re-dispatch loop is the load-bearing variable; ACs + surface are insufficient signal alone.

**Dispatch-dir hygiene.** This cycle's dispatch artifacts (logs, prompts, scratch) live at `/root/dispatch-380/`, **not** `/tmp/dispatch-380/`. Per today's δ adhoc `cn-sigma:threads/adhoc/20260519-dispatch-log-recovery-via-proc.md`, dispatched agents may rm -rf `/tmp` scratch under their own scratch-hygiene rules, deleting the JSONL log mid-write. `/root/dispatch-N/` is outside any agent's scratch hygiene scope. γ writes prompts to `/root/dispatch-380/`; δ reads them from there; α/β read no dispatch artifacts (their loaded skills come from the repo).

**Dispatch order:**
1. δ dispatches α with `/root/dispatch-380/alpha-prompt.md` → α writes `spawn.go` + `spawn_test.go` + modifies `cmd_activate.go` (and possibly `activate.go` if Option B), commits to `cycle/380`, writes `self-coherence.md` with review-readiness signal, exits
2. δ dispatches β with `/root/dispatch-380/beta-prompt.md` → β reviews on `cycle/380`; on APPROVE, merges into `main` per `beta/SKILL.md` §1, writes `beta-review.md` + `beta-closeout.md`, exits
3. γ resumes for close-out: re-dispatches α for `alpha-closeout.md` (per CDD §1.6a — α exits at review-readiness, returns for close-out only on γ re-dispatch); writes PRA + cdd-iteration (if findings tagged) + RELEASE.md; moves cycle dir to `.cdd/releases/{X.Y.Z}/380/`
4. δ executes release-boundary gates (tag, branch cleanup, release CI) per disconnect §3.4 after `gamma-closeout.md` lands on main

Per the operator/SKILL.md dispatch convention, α gets `--allowedTools "Read,Write,Bash"` and `--permission-mode acceptEdits`; β gets the same allowedTools (β needs Bash to run the AC-verification Go tests and the order-of-checks smoke commands during review).

## Failure modes to guard against (γ-side, pre-flagged for α/β)

1. **Exec mechanism drift: `os/exec.Cmd` instead of `syscall.Exec`.** `os/exec.Cmd` introduces a parent-child relationship that breaks interactive REPL handoff for many CLI tools (the parent holds a pipe; the child's TTY is the pipe, not the operator's terminal). `syscall.Exec` replaces the cn process image; the operator interacts directly with claude/codex. Tests passing on `os/exec.Cmd` mocks would produce a deceptively-green branch with broken production behavior. β verifies the non-test default exec function is `syscall.Exec`.

2. **Argv shape drift: `-p` / `--print` / `exec` between binary and prompt.** `claude -p "$prompt"` is the non-interactive print form; `codex exec "$prompt"` is the non-interactive subcommand. Either would render the cycle's spawn path useless — the operator would see a one-shot output, not a REPL. argv after argv[0] must be exactly `[prompt]`, one element.

3. **Default no-flag stdout regression (AC3).** The most subtle failure: a single new diagnostic line on stdout — a "use --claude for interactive" hint, a flag-availability header, a deprecation notice — breaks pipe consumers. Diagnostics belong on stderr; the 3.78.0 contract is that `cn activate HUB_DIR > prompt.md` produces the prompt with no surrounding noise. β diffs against a captured 3.78.0 fixture or a worktree-checkout smoke compare.

4. **Order-of-checks: mutual exclusion or LookPath after render.** Verifiable by the "non-existent HUB_DIR" trick documented in the α/β prompts. If the operator sees a "hub path not found" error before the mutual-exclusion error or the missing-binary error, the order is wrong. Both checks belong in the flag-validation step, before `activate.Run` is called.

5. **AC5 partial-fail: error names binary OR flag, but not both.** Operators set the flag deliberately and want to know which flag triggered the binary check. "claude: not found" is half-actionable; "claude (requested by --claude): not found in PATH; install it or ensure PATH includes the binary's directory" is fully actionable.

6. **Render-capture seam ambiguity: silent activate.go change.** α picks Option A (capture buffer at the call site, activate.go untouched) or Option B (add `activate.Render` entrypoint, small refactor). β cannot review a silent change to the 3.78.0 renderer. Self-coherence.md §Implementation notes must declare the choice; activate.go in the diff without a corresponding declaration is a binding finding (return RC, not silently approve).

7. **Test fragility: tests requiring `claude` / `codex` on the dev machine's $PATH.** Tests must run on a fresh CI image without those binaries. Use a fake-binary fixture in `t.TempDir()` with `t.Setenv("PATH", …)`, or substitute the exec hook directly.

8. **Scope creep: `--cursor` / `--aider` / `--auto` / `$CN_DEFAULT_BODY`.** All explicitly out of scope per the issue's §Non-goals. Any such addition is a binding finding.

9. **Tests check flag parsing only, not exec shape.** A test that asserts `--claude` is recognized without asserting the planned exec argv passes the file-read oracle while failing the source-of-truth invariant for AC1/AC2 — the operator's lived experience is the spawn, not the flag parse. β's pre-flagged AC partial-fail pattern #1.

10. **CLI surface change beyond the two flags.** The issue body's §Out-of-scope and §Non-goals do not preclude minor cmd_activate.go edits (the new flags are a CLI surface change); the binding limit is that the no-flag behavior is bytes-equal AND no other CLI shape changes (no new subcommands, no positional-arg reordering, no help-text restructuring beyond adding the two flag descriptions).

## Risk register and γ posture

| Risk | Mitigation | γ posture |
|---|---|---|
| `syscall.Exec` Windows portability | Issue's §Out-of-scope excludes Windows-specific quirks; α adds a stub or build-tag guard if needed | Acceptable; if stub is needed, document in self-coherence.md §Implementation notes |
| Render-capture seam → Option B drift | γ pre-validated both options; α declares choice in self-coherence.md | β return RC if undeclared |
| 3.78.0 stdout regression | β diffs against captured / smoke-compared 3.78.0 baseline | Binding finding if any bytes differ on the no-flag path |
| `--dry-run` as operator-facing flag (issue mentions as one option among several test-hook approaches) | γ recommends test-only internal hook; non-binding preference | Either is acceptable; γ noted preference in α prompt §Authoring discipline |
| Cross-repo / proposal lifecycle | None — operator-direct request, no source proposal | n/a; no STATUS event chain |

## Cycle-iteration triggers γ is pre-watching

- **Review churn (>2 rounds)** — single-focus Go cycle should converge in ≤2 rounds. If RC fires twice, root cause analysis required in PRA.
- **Mechanical overload** — small surface (~one new file + flag edits); mechanical ratio should be low. If high mechanical-ratio appears in β findings, suggests test fragility or build hygiene gap, not implementation drift.
- **Avoidable tooling failure** — the dispatch-dir hygiene fix landed today specifically to prevent tooling-loop failures in this cycle's δ session; γ pre-flags as a watch item.
- **Loaded-skill miss** — `eng/go/SKILL.md` should cover the Go authoring and test discipline. If a skill-miss finding lands in close-out, γ patches the skill in step 13a.

## Continuity note

cycle/380 ships on the same operator session as cycle/379 (3.78.0). The tooling loop is hot: same δ session, same Claude infrastructure, same gh / git / claude / codex binaries. γ's expectation is that this cycle ships as 3.79.0 within the same disconnect window, modulo the operator's tolerance for back-to-back releases. If the operator signals a freeze before β merge, γ holds the cycle on the branch until cleared.

## Source

Operator-direct request during the cnos 3.78.0 disconnect phase. Anchor: cnos δ session of 2026-05-19 (cn-sigma daily reflection `threads/reflections/daily/20260519.md`). No cross-repo proposal; this cycle is the direct follow-on to cycle/379 closing the loop on `cn activate`'s hub-side flow.
