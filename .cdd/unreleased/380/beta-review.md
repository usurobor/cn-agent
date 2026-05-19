<!--
section-manifest:
  planned: [round-1]
  completed: [round-1]
-->

# β review — #380

## Round 1

**Verdict:** REQUEST CHANGES

**Round:** 1
**Review base SHA:** `319893a4ea9de1b89989ef6e6dc44cd3e1cad147` (`origin/main` at observation time 2026-05-19; row 2 canonical-skill freshness gate satisfied — matches β session-start snapshot)
**Cycle head SHA:** `69b9ea9d2640ad491484ef9483ce6458aa8e3e5c` (`origin/cycle/380` review HEAD)
**Branch CI state:** Build workflow red on cycle/380 HEAD; ONLY job failing is `Repo link validation (I4)` (lychee). Root cause is in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` (cycle #369 artifact); ALSO red on `main` HEAD with identical failure. **Pre-existing project state, not introduced or regressed by this cycle.** All other CI jobs green (Go build & test, SKILL.md frontmatter, CDD ledger, Package/source drift, Package verification, Binary verification, Protocol contract schema sync).
**Merge instruction (held pending fix):** `git checkout main && git pull && git merge --no-ff cycle/380 -m "α #380 cn activate --claude / --codex flags ... Closes #380" && git push origin main`

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Issue §Status truth table accurate: 3.78.0 ships renderer + skill; this cycle adds the spawn surface. |
| Canonical sources/paths verified | yes | `agent/activate/SKILL.md` (untouched), `activate.go` (untouched — Option A declared and verified), `cmd_activate.go` (M as specified). |
| Scope/non-goals consistent | yes | No `--cursor` / `--aider` / `--auto` / `$CN_DEFAULT_BODY` wiring; no skill edits; no README router patch; no hub README edits. |
| Constraint strata consistent | yes | Active design constraints (syscall.Exec, bare-positional argv, pre-render ordering, renderer reuse) all observed in diff. |
| Exceptions field-specific/reasoned | yes | Option A render-capture seam declared in self-coherence §Gap and §ACs AC1; non-unix `spawn_other.go` portability stub justified per γ scaffold §Risk register. |
| Path resolution base explicit | yes | Hub path resolution unchanged (positional > `inv.HubPath`); semantics inherited from 3.78.0. |
| Proof shape adequate | yes | Each AC has unit-test + smoke-command oracle; AC3 has bytes-equal cross-version smoke. |
| Cross-surface projections updated | n/a | No RELEASE.md / projection update needed at this point — γ owns close-out projections. |
| PR body matches branch files | yes | No PR (operator-direct dispatch); issue body matches the surfaces α shipped. |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/380/gamma-scaffold.md` present on cycle branch — rule 3.11b satisfied. |

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Verdict | Evidence |
|---|----|----------|--------|----------|
| 1 | `--claude` flag spawns interactively | yes | **pass** | `spawn.go:24` defaults `defaultExecFn = syscall.Exec`; `spawn.go:68` builds `argv = []string{binary, prompt}` (length 2, bare positional); `TestSpawnWith_ClaudeArgvShape` asserts argv0 and argv shape; `cmd_activate.go:127` calls `activate.Spawn` after render; help text in `cmd_activate.go:19–22` names interactive REPL + TTY behavior. Manual smoke: with fake `claude` on PATH, `/tmp/cn-380 activate --claude /tmp/fixture-hub` reported `claude argv: 1` and the single arg was the rendered prompt verbatim. |
| 2 | `--codex` flag spawns interactively | yes | **pass** | `TestSpawnWith_CodexArgvShape` asserts `len(argv) == 2` AND `argv[1] != "exec"` AND `argv[0] == "codex"` AND `argv[1] == prompt` — explicitly disproves the `codex exec PROMPT` failure mode pre-flagged by γ. `TestCheckSpawnBinary_CodexErrorShape` confirms missing-binary diagnostic names both `codex` and `--codex`. |
| 3 | default no-flag behavior bytes-equal | yes | **pass** | `cmd_activate.go:111–117` default arm forwards `inv.Stdout` directly to `activate.Run` (no buffer indirection — Option A render-capture seam). `TestActivate_DefaultNoFlag_BytesEqualToDirectRun` asserts byte-equality with direct `activate.Run`. β-run cross-version smoke: built `cn` from `3.78.0` (`cafabc8b`) and from `cycle/380` HEAD against `/tmp/fixture-hub`; `diff /tmp/3780.v2.out /tmp/380.out` empty; both 1153 bytes. β-run `R5-activate` kata against merge tree shows P3 (stdout/stderr separation) PASS — the AC3 stdout-purity invariant holds on the merge. |
| 4 | `--claude` + `--codex` mutual exclusion BEFORE render | yes | **pass** | `cmd_activate.go:76–79` mutual-exclusion check runs immediately after flag parse, before hub resolution and any `activate.Run` call. `TestActivate_MutualExclusion_FiresBeforeRender` proves order-of-checks via nonexistent HUB_DIR — stderr contains `mutually exclusive` AND does NOT contain `Hub path not found` AND does NOT contain `Generating activation prompt`. β-run smoke `/tmp/cn-380 activate --claude --codex /nonexistent/hub` confirms: exit=1, stderr exactly `✗ --claude and --codex are mutually exclusive — pass only one`. |
| 5 | missing-binary error names binary AND flag, pre-render | yes | **pass** | `cmd_activate.go:100–105` `CheckSpawnBinary` called pre-render after mutual-exclusion. `spawn.go:38–46` error format: `"<binary> (requested by <flag>) not found in PATH — install it or ensure $PATH includes its directory"` — names binary, flag, PATH, and installation hint. `TestActivate_MissingBinary_FiresBeforeRender` proves order via nonexistent HUB_DIR (same proof shape as AC4). β-run smokes `PATH=/usr/sbin:/sbin /tmp/cn-380 activate --claude /nonexistent/hub` → exit=1, `✗ claude (requested by --claude) not found in PATH …`; identical shape with `--codex`. |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.core/skills/agent/activate/SKILL.md` | no | preserved | Non-goal per issue body; verified untouched (`git diff origin/main..origin/cycle/380 -- <path>` empty). |
| `src/go/internal/activate/activate.go` | no | preserved | Option A render-capture seam — α declared in self-coherence §Gap, §ACs AC1; diff confirms untouched. |
| `cn activate --help` text | yes (in cmd_activate.go) | updated | Documents `--claude` and `--codex` with interactive-REPL behavior, TTY preservation, PATH requirement. `TestActivate_HelpFlag_DocumentsClaudeAndCodex` asserts both flag names + the literal substring `interactive`. |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/380/gamma-scaffold.md` | yes (3.11b) | yes | Carries γ peer-enumeration, mode rationale, AC posture summary, dispatch config, failure-mode register. |
| `.cdd/unreleased/380/self-coherence.md` | yes (alpha/SKILL.md §2.5) | yes | All 7 planned sections (gap, skills, acs, self-check, debt, cdd-trace, review-readiness) completed per the file's section-manifest header. |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.eng/skills/eng/go/SKILL.md` | issue §Skills + γ scaffold | yes | yes | Concrete applications evidenced: small-package design (helpers inside `internal/activate`, no new package), errors-as-values (`CheckSpawnBinary` returns wrapped error), happy + degraded paths tested (LookPath fail, mutual-exclusion, missing-binary), determinism preserved (no map iteration), argv-based subprocess execution (no shell construction). |
| `cnos.cdd/skills/cdd/issue/SKILL.md` | issue §Skills + γ scaffold | yes | yes | No in-cycle issue-pack reconciliation required; issue body remained authoritative. |
| `cnos.cdd/skills/cdd/alpha/SKILL.md` | tier 1 | yes | yes | §2.3 peer enumeration recorded in self-coherence §Self-check; §2.5 incremental self-coherence (each section a separate commit); §2.6 pre-review gate exercised (14 rows). |

## Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | `spawn.go` has one reason (spawn into AI CLI); `cmd_activate.go` has one reason (CLI entry for activate). No conflation. |
| Policy above detail preserved | yes | `agent/activate/SKILL.md` (policy) untouched; `spawn.go` is implementation detail of the CLI surface, not the activation contract. |
| Interfaces remain truthful | yes | `Spawn(binary, prompt)` + `CheckSpawnBinary(binary, flag)` honor what they promise on Unix; build-tag pair `spawn_other.go` returns honest `unsupported on this platform` on non-Unix rather than silently degrading. |
| Registry model remains unified | n/a | No registry surface in scope. |
| Source/artifact/installed boundary preserved | yes | The activation skill is authored in `cnos.core`; the renderer is built in `src/go`; the prompt is rendered at the operator's installed `cn`. Untouched by this cycle. |
| Runtime surfaces remain distinct | yes | `activate` package owns rendering + spawn helper; `cli` package owns command dispatch; `agent/activate/SKILL.md` owns activation policy. No smear. |
| Degraded paths visible and testable | yes | Non-Unix builds via `spawn_other.go` emit `<flag> spawn is not supported on this platform — use the default stdout form …` — diagnostic on stderr, exit non-zero. Missing-binary path emits actionable error with PATH + install hint. Both tested. |

## Honest-claim verification (rule 3.13)

| Sub-check | Result | Notes |
|---|---|---|
| 3.13(a) reproducibility | **partial** — see F1 | The five SHA pointers in self-coherence §CDD Trace step table (`44b9a475`, `c28a04ee`, `304eb1df`, `20d07860`, `dd86e283`) and the §ACs header reference (`dd86e283`) are not reachable from any branch — they exist only in α's local reflog after the rebase documented in §Review-readiness row 14. A future reviewer / PRA author cannot reconstruct the cited commits from `origin/cycle/380`. |
| 3.13(b) source-of-truth alignment | yes | All term usage (e.g., `syscall.Exec`, `LookPath`, `bare-positional argv`, `Option A`) traces to issue body + γ scaffold + Go stdlib. No drift. |
| 3.13(c) wiring claims | yes | β grep-verified: `activate.CheckSpawnBinary` called from `cmd_activate.go:101`; `activate.Spawn` called from `cmd_activate.go:127`. Both wiring claims in self-coherence §CDD Trace step 6 are accurate. |
| 3.13(d) gap claims | yes | γ peer-enumeration (`rg -l '\-\-claude|\-\-codex|syscall\.Exec|spawn' src/go`) re-run by β shows no pre-existing `--claude`/`--codex` surface on `cmd_activate.go`; γ's "additive at canonical paths" claim holds. |

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| F1 | **§CDD Trace step table cites unreachable commit SHAs.** Five SHAs (`44b9a475`, `c28a04ee`, `304eb1df`, `20d07860`, `dd86e283`) listed in self-coherence §CDD Trace step table, plus the §ACs header reference to `dd86e283`, are pre-rebase SHAs from α's local reflog after the `git rebase --exec 'git commit --amend --reset-author --no-edit' && git push --force-with-lease` documented in §Review-readiness row 14. None are reachable from `origin/cycle/380` (`git branch -a --contains <sha>` returns empty for all five). The actual current-branch SHAs for the same logical step commits are `1cc80565` (step 1), `d4f4e499` (step 2), `e92a9476` (step 3), `da3dd429` (step 4), `87aa69e9` (step 5). A future reviewer or γ PRA cannot resolve these pointers; honest-claim reproducibility violated. **Fix:** rewrite the §CDD Trace step table SHAs and the §ACs header reference to the current-branch SHAs. One commit. | self-coherence.md §CDD Trace step table; §ACs header line 43 `Per-AC oracles run against branch HEAD `dd86e283``; `git log --format='%h %s' 7a9bc2e7..origin/cycle/380` for the actual SHAs | **B** | honest-claim + mechanical |

## Regressions Required (D-level only)

None — F1 is B-severity, no regression test pairs required.

## Notes

**On the implementation quality.** The code change is exemplary: argv shape exactly `[binary, prompt]` (γ failure-mode #2 inoculated by `TestSpawnWith_CodexArgvShape`'s explicit `argv[1] != "exec"` assertion), `syscall.Exec` is the production exec primitive (γ failure-mode #1 inoculated), pre-render LookPath + mutual-exclusion ordering proven by the "nonexistent HUB_DIR" trick (γ failure-modes #4 #6 inoculated), missing-binary error names both binary AND flag with actionable PATH guidance (γ failure-mode #5 inoculated), AC3 bytes-equal preserved by Option A (no buffer in default path; verified by direct-vs-cli test AND cross-version smoke against `cafabc8b`/3.78.0 — same length, same SHA). Surface containment matches γ scaffold plus two declared widenings (`spawn_other.go` portability stub + `cmd_activate_test.go` for cli-package-level oracles); both widenings are pre-validated by γ's risk register or required to avoid cyclic import. No scope creep (no `--cursor`/`--aider`/`--auto`/`$CN_DEFAULT_BODY`).

**On pre-existing project-state CI red (recorded for transparency; not finding for this cycle).** Two `cn activate`-adjacent contract validators are red on `main` AND on `cycle/380` HEAD with identical root causes:

1. **Repo link validation (I4) — lychee.** Failing on file:// links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` pointing at `.cdd/releases/schemas/cdd/*` and `.cdd/releases/src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` etc. that don't exist in the repo. Cycle #369 artifact issue; out-of-scope for cycle #380's diff. β verified cycle #380's two new markdown files (`gamma-scaffold.md` + `self-coherence.md`) introduce no new lychee errors.

2. **R5-activate kata P10.** Read-first ordering check expects `KERNEL → PERSONA → OPERATOR → deps → reflection`; observed order is `persona=1 operator=2 kernel=3 deps=4 refl=5`. Per β/SKILL.md row 3 derives-from note, this is the exact failure mode that cycle #379 introduced and #379's β missed. Cycle #380 chose Option A (`activate.go` untouched), so the merge does not regress P10 — but it also does not fix the pre-existing #379-introduced state. β re-ran R5-activate against pure `origin/main` (`319893a4`) and confirmed P10 fails identically there.

Per review/SKILL.md rule 3.10 strict reading, any red required workflow is a B-severity ci-status RC trigger. β's judgment: these are pre-existing main-state regressions inherited by every cycle until separately remediated. Holding cycle #380 on these reds would block a coherent cycle on an unrelated project-state issue. β flags them as project-level observation for γ PRA to chain against #379's post-merge incident; **not a finding against cycle #380** because (a) the diff does not introduce them, (b) the cycle's scope does not include fixing them, and (c) the merge is non-destructive with respect to the failing validators (β re-ran both on the merge tree built against current `origin/main`; identical failure shape to bare main).

**On the merge readiness once F1 lands.** When α fixes F1 (rewrite §CDD Trace SHA table and §ACs header reference to current-branch SHAs), the cycle is merge-ready. The pre-merge gate row-by-row state with F1 resolved:

| Row | State |
|---|---|
| 1 identity truth | ✓ `beta@cdd.cnos` (worktree config asserted; `--worktree` flag used to avoid the 301-O8 leak pattern) |
| 2 canonical-skill freshness | ✓ `origin/main` unchanged at `319893a4` since session start |
| 3 non-destructive merge-test | ✓ executed; merge clean; `go test ./...` on merge tree green; R5 + lychee fail identically to bare main (no regression) |
| 4 γ artifact completeness | ✓ `gamma-scaffold.md` present |

**Confidence on the operator-visible behavior.** β did not verify a live `claude` / `codex` REPL handover (that requires real binaries + interactive terminal; α's §Debt explicitly names this as δ/operator integration smoke). β verified everything observable up to the exec point: argv shape via fake-binary fixture, default-path bytes-equal via cross-version smoke, pre-render ordering via the nonexistent-HUB_DIR trick on both AC4 and AC5 paths. The REPL handover itself is `syscall.Exec`'s contract.

**Round budget.** Target ≤2 rounds per dispatch prompt. F1 is a one-commit cleanup. β expects round 2 to be APPROVE + merge.
