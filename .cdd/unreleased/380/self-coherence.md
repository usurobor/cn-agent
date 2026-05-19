<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills]
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
