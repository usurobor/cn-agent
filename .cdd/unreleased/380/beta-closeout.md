<!--
section-manifest:
  planned: [review-summary, implementation-assessment, technical-review, process-observations, release-notes]
  completed: [review-summary, implementation-assessment, technical-review, process-observations, release-notes]
-->

# β close-out — #380

## Review summary

**Verdict:** APPROVED (round 2).
**Merge commit:** `770ea1b4` on `main` (non-fast-forward merge of `cycle/380` into `main`, message `α #380 cn activate --claude / --codex flags — spawn AI body interactively`, `Closes #380`).
**Cycle branch tip merged:** `4a7cda79` (β R2 commit) → tree-identical for production code to `d838e7e9` (α fix-round 1) and `87aa69e9` (last α implementation commit).
**Round count:** 2 (target ≤2 met; γ's review-churn cycle-iteration trigger not fired).
**Findings:** 1 (R1 F1, B-severity honest-claim + mechanical — §CDD Trace SHAs unreachable from `origin/cycle/380` after α's email-rewrite rebase). Closed by α `91a2cad6` "rewrite §CDD Trace SHAs to current-branch values".

The cycle ships the spawn-into-CLI surface that closes the loop on `cn activate`: 3.78.0 shipped the renderer; 3.79.0 (next tag, δ's release boundary) ships the operator-facing handover into an interactive `claude` / `codex` REPL with the body already activated. Default no-flag stdout behavior is preserved bytes-equal to 3.78.0.

## Implementation assessment

**Surface containment matches γ scaffold.** α landed exactly the three γ-named files (`spawn.go` NEW, `spawn_test.go` NEW, `cmd_activate.go` M) plus two declared widenings: `spawn_other.go` (non-unix build-tag stub, pre-validated by γ's risk register) and `cmd_activate_test.go` (cli-package-level oracles for AC3/AC4/AC5 — required to avoid a cyclic import between `cli` and `activate`). `activate.go` is byte-identical to 3.78.0 (Option A chosen and declared in self-coherence §Gap and §ACs AC1).

**Active design constraints honored.** β verified every binding constraint in the issue body:

- `syscall.Exec` not `os/exec.Cmd` — `spawn.go:24` `defaultExecFn = syscall.Exec` (γ failure mode #1 inoculated)
- Bare-positional argv — `spawn.go:68` `argv := []string{binary, prompt}`; `TestSpawnWith_CodexArgvShape` explicitly asserts `argv[1] != "exec"` (γ failure mode #2 inoculated; cycle's purpose preserved against `codex exec` non-interactive subcommand drift)
- Pre-render LookPath + mutual exclusion — `cmd_activate.go:76,100` checks fire before any call to `activate.Run`; proven by the "nonexistent HUB_DIR" trick in both `TestActivate_MutualExclusion_FiresBeforeRender` and `TestActivate_MissingBinary_FiresBeforeRender` (γ failure modes #4 and #6 inoculated; β-run smokes confirm stderr never carries `Hub path not found` or `Generating activation prompt` when the pre-render check fires)
- Renderer reuse — `activate.go` and `agent/activate/SKILL.md` both untouched (verified by `git diff origin/main..origin/cycle/380` against both paths returning empty); γ failure mode #6 (silent renderer change) inoculated by α's explicit Option A declaration in self-coherence §Gap

**Tests are exec-shape-asserting, not just flag-parsing.** β confirmed `spawn_test.go` records the planned argv via the test-injectable `execFunc` hook and asserts on `argv0` + `argv` length + each `argv[i]` element. γ's pre-flagged partial-fail #1 (tests check flag parsing but not exec shape) does not apply: every AC1/AC2 oracle asserts the production argv shape via the recording hook, and tests do not require `claude` / `codex` on the dev machine's PATH (use `t.Setenv("PATH", t.TempDir())` plus `fakeLookPath` fixtures throughout).

**Default-path bytes-equality is structural, not merely tested.** Per Option A, the no-flag arm in `cmd_activate.go:111-117` passes `inv.Stdout` straight to `activate.Run` with no buffer indirection. The bytes-equal property holds *by construction*, not by accident: only the spawn arm captures into a `bytes.Buffer`. `TestActivate_DefaultNoFlag_BytesEqualToDirectRun` is the unit oracle; β-run cross-version smoke against a freshly-rebuilt 3.78.0 binary at `cafabc8b` confirms `sha256sum` identical (1151 bytes; `b50c6c90a2d0…`).

**Error message ergonomics.** AC5 partial-fail #7 (binary OR flag, not both) inoculated by `spawn.go:38-46`:
```
✗ <binary> (requested by <flag>) not found in PATH — install it or ensure $PATH includes its directory
```
Names binary, flag, PATH, and install hint in one line — fully actionable.

## Technical review

**Code quality.** `spawn.go` is 70 lines of straightforward Go: typed function hooks (`execFunc`, `lookPathFunc`) for test injection; production defaults via unexported package-level vars (`defaultExecFn`, `defaultLookPath`); errors returned, never panicked; no shell construction (`syscall.Exec` takes a `[]string`, not a shell command). `cmd_activate.go` cleanly separates flag parsing, validation, hub resolution, pre-render checks, render, and spawn — each step has one reason to change. Help text additions are concise and document the interactive-spawn behavior, TTY preservation, and PATH requirement.

**Test coverage.** 16 new test functions (7 in `spawn_test.go` + 9 in `cmd_activate_test.go`) cover:
- Argv shape for both binaries (positive + critical negative for codex)
- LookPath failure path
- Missing-binary error message shape for both flags (with explicit PATH + install hint assertions)
- Default-hook wiring (production path actually uses `exec.LookPath`)
- Mutual exclusion (positive + pre-render ordering)
- Missing-binary CLI path (per flag + pre-render ordering)
- Help text documentation
- Unknown-flag rejection
- AC3 bytes-equal vs direct `activate.Run`

All existing 3.78.0 / #379 invariants in `activate_test.go` remain green on `cycle/380` HEAD and on the merge tree.

**Cross-platform handling.** `spawn.go` is `//go:build unix`; the `spawn_other.go` stub returns a clear "unsupported on this platform" error on non-unix builds for the `--claude` / `--codex` paths only. The no-flag default path remains functional everywhere. Stays inside the issue's "macOS/Windows-specific quirks beyond what falls out of syscall.Exec portability" non-goal: only what falls out for free is shipped.

**Honest-claim discipline.** α's R1 self-coherence cited five SHAs in §CDD Trace step table that were pre-rebase reflog values from a path-(a) `git rebase --exec` email rewrite — none reachable from `origin/cycle/380`. β R1 caught this via 3.13(a) reproducibility (an external reader cannot resolve the cited pointers). α's fix-round applied the intra-doc repetition rule (`alpha/SKILL.md` §2.3, derived from #266 F3-bis), grepping the file for every SHA occurrence and catching three additional hits β had not enumerated by name — F1 closed in one commit without an F1-bis.

## Process observations

**γ → α → β protocol coherence.** This cycle ran clean §5.1 canonical multi-session (per `gamma-scaffold.md` §Dispatch configuration). γ's scaffold was load-bearing: the failure-mode register pre-flagged every binding pattern β had to verify (argv shape, exec primitive, pre-render ordering, render-capture seam declaration, test fragility). α's self-coherence inoculated against every γ-pre-flagged binding pattern by name — the render-capture seam Option A declaration in §Gap directly cites γ scaffold §"Render-capture seam"; AC1/AC2 oracles cite γ failure modes #1, #2 by number; AC5 error format cites γ failure mode #5. β's R1/R2 verdicts cite the same γ-pre-flagged patterns by number in the AC evidence column. The triadic protocol passed information without loss.

**Round budget.** Target ≤2; actual 2. γ's review-churn cycle-iteration trigger (>2 rounds = §9.1 trigger) not fired. R1's single B-severity finding was a self-coherence artifact issue, not an implementation defect — no production-code delta in the R1→R2 transition. Reviewer-implementer turnaround was tight.

**Pre-existing CI red disposition.** Two project-level CI failures inherited from prior cycles:
1. `Repo link validation (I4)` — lychee failing on file:// links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` (cycle #369 artifact). Identical failure on `main` HEAD (`319893a4`) and `cycle/380` HEAD (`d838e7e9`).
2. `R5-activate` kata P10 — read-first ordering inverted from spec (KERNEL → PERSONA → OPERATOR → deps → reflection expected; persona=1 operator=2 kernel=3 deps=4 refl=5 observed). Per β/SKILL.md row 3 derives-from note, this is the failure mode #379 introduced and #379's β missed. Cycle #380 chose Option A (`activate.go` untouched) → P10 fails identically on the merge tree and on bare `origin/main`.

β R1/R2 judgment per review/SKILL.md 3.10: inherited project-state CI red is not RC for a cycle whose diff neither introduces nor includes scope to fix the failure. Recorded for γ PRA chaining. Two follow-on cycles would be appropriate: one to remediate cycle #369's stale link references, one to repair the #379-introduced read-first ordering in `activate.go`.

**Identity hygiene.** β's merge-test worktree at `/tmp/cnos-merge-380` used `git config --worktree user.email "beta-merge-test@local"` per β/SKILL.md row 1 to prevent the cycle-#301 O8 leak pattern; worktree was removed post-test. All β-authored commits on `cycle/380` and `main` are `beta@cdd.cnos`.

## Release notes (handoff to γ for RELEASE.md authoring)

γ owns the RELEASE.md cycle entry per `release/SKILL.md`. β's hand-off summary:

**3.79.0 — `cn activate` interactive-spawn flags**

`cn activate HUB_DIR` (no flag) is unchanged from 3.78.0 — bytes-equal stdout, pipe-and-redirect consumers untouched.

Two new flags add the operator-facing handover:

- `cn activate --claude HUB_DIR` — after rendering, replaces the cn process with `claude "$prompt"` via `syscall.Exec`. Operator lands in a live interactive `claude` REPL with the body already in the post-activation state. TTY/stdin/stdout inherited. Requires `claude` on `$PATH`.
- `cn activate --codex HUB_DIR` — same shape with `codex "$prompt"` (bare-positional, not `codex exec` — interactive REPL, not non-interactive subcommand).

Validation:
- `--claude` and `--codex` are mutually exclusive (exits non-zero pre-render).
- If the requested CLI is not on `$PATH`, exits non-zero pre-render with an actionable error naming both the binary and the requesting flag.

Non-unix platforms: the default path (no flag) is unchanged; `--claude` / `--codex` on non-unix exits non-zero with `<flag> spawn is not supported on this platform`.

**Skill / doctrine:** none changed this cycle. `agent/activate/SKILL.md` (body-agnostic activation skill, 3.78.0) and `activate.go` (renderer, 3.78.0) are both untouched. The spawn surface is a `cn`-side concern; the activation skill stays body-agnostic.

**Operator-visible projection:** new flags appear in `cn activate --help`. RELEASE.md should call out: pipe consumers are not impacted; the new flags are opt-in additions for the local-dev interactive flow.

**Known follow-ons (already named in issue §Deferred / §Non-goals):** `--cursor` / `--aider` (other AI CLIs), `--auto` (body-detection), `$CN_DEFAULT_BODY` (env-var default), macOS/Windows-specific quirks. None blocking.

**Pre-existing CI red to chain in PRA:** `Repo link validation (I4)` from cycle #369 artifact paths; `R5-activate` P10 read-first ordering from cycle #379 `activate.go` refactor. Neither introduced by #380; neither remediated by #380. Two separate cleanup cycles recommended.

β's last action is this close-out. Release tagging, branch cleanup, and the disconnect release boundary belong to δ via γ's close-out per CDD.md §1.6.
