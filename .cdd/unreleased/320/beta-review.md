## β Review — #320 cn activate

**Round:** 1
**origin/main SHA at review:** `8bf7bf8413a27af8d2902380c660f57b417076e7`
**Branch head SHA:** `99bd7b8fe59871ce0513c74f7e187e6c61fd8b10`
**Branch CI state:** green — `go test ./...` 10/10 pass; `go build ./...` clean; full suite clean (no regressions)
**Verdict:** APPROVED

**Merge instruction:** `git merge --no-ff cycle/320 -m "feat(cli): cn activate — bootstrap prompt generation from hub state (Closes #320)"` into main, then push.

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Implementation is prompt-generation only; no model invoked, no daemon started — matches issue §Status truth |
| Canonical sources/paths verified | yes | main.go, cmd_activate.go, activate.go, activate_test.go all cited correctly; paths resolve |
| Scope/non-goals consistent | yes | Only `cn activate` added; no `cn agent` namespace; no existing commands modified |
| Constraint strata consistent | yes | Hard gates all satisfied: `NeedsHub: false`, stdout=prompt only, stderr=diagnostics, no secrets, no model, no daemon |
| Exceptions field-specific/reasoned | yes | No exceptions claimed; bootstrap path explicitly waived (single kernel command, no version snapshot required) |
| Path resolution base explicit | yes | `NeedsHub: false` means activate resolves hub itself; documented in AC5 |
| Proof shape adequate | yes | 10 tests covering 7 positive + 3 negative cases; fixture hub with real file system |
| Cross-surface projections updated | yes/n/a | `cn activate` auto-appears in `cn help` via registry; README/OPERATOR update is optional per issue |
| No witness theater / false closure | yes | Tests exercise real Go code against real filesystem fixtures; no mock-pass |
| PR body matches branch files | n/a | CDD triadic cycle; no PR |

---

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | Command exists — `cn activate` prints bootstrap prompt | yes | met | `ActivateCmd` registered in main.go; `Run()` → `activate.Run()` → `writePrompt(stdout)` |
| AC2 | Flat kernel command, not agent namespace | yes | met | `Source: SourceKernel, Tier: TierKernel`; no `cn agent` command introduced |
| AC3 | Hub discovery from cwd | yes | met | `discoverHub()` in main.go; `cmd_activate.go` uses `inv.HubPath` when no arg |
| AC4 | Explicit local HUB_DIR works | yes | met | `cmd_activate.go`: `if len(inv.Args) > 0 && !isFlag(inv.Args[0]) { hubPath = inv.Args[0] }` |
| AC5 | Command does not require pre-discovered hub | yes | met | `NeedsHub: false`; main.go hub gate does not fire; command resolves hub itself |
| AC6 | Prompt uses safe hub state | yes | met | reads only `.cn/config.json`; no secrets.env, no .env; test verifies exclusion |
| AC7 | Prompt is model-ready but does not call a model | yes | met | no `os/exec` or HTTP imports; `writePrompt()` writes to `io.Writer` only |
| AC8 | stdout/stderr contract preserved | yes | met | `writePrompt()` → stdout; `→ Generating...` → stderr; errors → stderr; test `TestRunPositive_StdoutOnly` verifies |
| AC9 | Claude CLI example is valid | yes | met | example includes query arg: `claude -p "Activate this cnos hub using the bootstrap prompt on stdin."` |
| AC10 | Docs preserve current-vs-target truth | yes | met | help text and prompt Notes: "generates a prompt only. No model is invoked." No daemon/runtime claimed |
| AC11 | Namespace migration not in scope | yes | met | only `reg.Register(&cli.ActivateCmd{})` added to main.go; no existing commands touched |
| AC12 | Tests cover positive and negative cases | yes | met | 7 positive + 3 negative = 10 tests; all pass |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/go/internal/activate/activate.go` | yes | present | domain logic; complete |
| `src/go/internal/activate/activate_test.go` | yes | present | 10 tests; all green |
| `src/go/internal/cli/cmd_activate.go` | yes | present | CLI dispatch wrapper |
| `src/go/cmd/cn/main.go` | yes | present | registration line added |
| README.md / OPERATOR.md | no | n/a | optional per issue; existing docs cover manual activation |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/320/self-coherence.md` | yes | yes | complete; gap, mode, AC mapping, CDD Trace through 7a, review-readiness |
| `.cdd/unreleased/320/beta-review.md` | yes | yes (this file) | β authoring now |
| `.cdd/unreleased/320/alpha-closeout.md` | yes (post-merge) | pending | α writes after merge |
| `.cdd/unreleased/320/beta-closeout.md` | yes (post-merge) | pending | β writes after merge |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| design | issue Tier 3 | yes | yes | prompt generation separated from runtime; no daemon/dispatch imported |
| write | issue Tier 3 | yes | yes | help text and Notes don't overclaim runtime |
| eng/tool | issue Tier 3 | yes | yes | stdout/stderr contract; fail-fast on invalid hub |
| eng/test | issue Tier 3 | yes | yes | invariants before tests; positive+negative coverage |
| eng/ux-cli | issue Tier 3 | yes | yes | `✗` prefix on errors; `→` for diagnostic progress; causal messages |
| eng/go | issue Tier 3 | yes | yes | domain logic in `internal/activate/`; thin CLI wrapper in `internal/cli/`; `slices.Sort` on Go 1.24 ✓ |

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| F1 | Self-coherence AC1 evidence cites `main.go:43`; actual line is 42 | `main.go` Read output line 42: `reg.Register(&cli.ActivateCmd{})` | A | mechanical |

No other findings. F1 is a clerical off-by-one in an internal artifact cite; the implementation is correct and the registered command is at line 42 not 43. Per review §3.6 (approve when coherent, not when perfect) and §3.5 (no phantom blockers), F1 does not block approval — it is noted here for the record. The system is coherent.

## Notes

- `isFlag` is defined only in `cmd_activate.go` within the `cli` package. No conflict with other files. Appropriate for MVP scope.
- Context parameter `_ context.Context` is accepted but not forwarded in `activate.go`. No async operations in this command; acceptable for MVP.
- `slices` package requires Go 1.21+; `go.mod` declares `go 1.24`. ✓
- CI verified independently by β: `go test ./...` 10/10 pass; `go build ./...` clean; full suite clean.
