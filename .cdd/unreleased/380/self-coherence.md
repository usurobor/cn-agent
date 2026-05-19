<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap]
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
