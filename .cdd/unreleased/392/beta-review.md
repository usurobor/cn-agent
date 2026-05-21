# β review — cycle/392

**Issue:** cnos#392 (Phase 3 remediation v2; supersedes #391).
**Mode:** β-α-collapsed (γ+α+β-collapsed-on-δ). Acknowledged.
**Special rigor:** implementation-contract conformance + eng/go SKILL section-by-section, per dispatch (the wrong-language and wrong-location failure classes from #389 and #391 motivated this rigor).

## Round 1

### Implementation-contract conformance audit (7 pinned axes)

| Axis | Pinned | Observed | Verdict |
|---|---|---|---|
| Language | Go | All V source under `src/packages/cnos.cdd/commands/cdd-verify/` is `.go`. `cn-cdd-validate-receipt` (Python) `git rm`'d. | ✅ |
| CLI integration target | `cn cdd-verify` subcommand through `cn` CLI dispatch — NOT a separate binary | `cli.CddVerifyCmd` registered at `src/go/cmd/cn/main.go:43`. Tier `TierKernel` (`cli/cmd_cdd_verify.go:21`). cddverify package compiles INTO the cn binary via go.work + replace directive in `src/go/go.mod`. No standalone binary; `which cn-cdd-verify` returns nothing post-uninstall. | ✅ |
| Package scoping | `src/packages/cnos.cdd/commands/cdd-verify/` — NOT `src/go/internal/` | All V source under the cnos.cdd package. Only the 47-line thin dispatch wrapper `cmd_cdd_verify.go` lives in `src/go/internal/cli/`, satisfying eng/go §2.18. | ✅ |
| Existing-binary disposition | Remove Python; remove bash (or thin shim) | Both fully removed (operator preference); `git rm` recorded. | ✅ |
| Runtime deps | Go + cue only; no new deps | `src/go/go.mod` adds only the local `replace` directive (1 require for the local module). `src/packages/cnos.cdd/commands/cdd-verify/go.mod` is empty after `go mod tidy`. YAML parsing routes through `cue export` subprocess (cue already pinned). | ✅ |
| JSON schema contract | Unchanged | `git diff origin/main..HEAD -- schemas/cdd/validation_verdict.schema.json` is empty. | ✅ |
| Backward compat | All existing `cn cdd-verify` invocations work | `--unreleased`, `--all`, `--version`, `--pr`, `--cycle`, `--triadic`, `--receipt`, `--json`, `--contract`, `--structural-only`, `--repo-root`, `--bundle`, `--exceptions` all reproduced. BC.a-b PASS in main oracle (37/37). 39/39 PASS in `test-cn-cdd-verify.sh`. | ✅ |

**7/7 axes honored. No deviation.**

### eng/go SKILL compliance audit (§1.0 – §4.9)

Walked the self-coherence compliance table line-by-line against the actual source. Spot checks:

- §2.4 Context: `Validate(ctx, opts ValidateOptions)` at `cddverify.go:33` — ctx first parameter ✓.
- §2.5 Errors: `parse.go:25` `return nil, fmt.Errorf("decode receipt json: %w", err)` — wrap with `%w` ✓; lower-case prefix ✓.
- §2.6 Side-effect boundaries: `cuevet.go::cueVet:18` shells out via `exec.CommandContext`; `git.go::gitLogAuthor:14` similarly. Pure parsing in `parse.go::ParseReceiptJSON:18`. ✓
- §2.10 Testing: 27 tests in `cddverify_test.go`, table-driven with `t.Run` subtests. ✓
- §2.13 Determinism: `sortedKeys` (`counterfeit.go:288-300`) ensures C4 rule fires in sorted-key order across the `evidence_refs` map. `KnownProtocols()` returns sorted slice. ✓
- §2.17 Purity boundary: `ParseReceiptJSON([]byte)` pure ✓; `ReadReceipt(ctx, path)` IO wrapper ✓. Confirmed no `os.` imports inside `ParseReceiptJSON`.
- §2.18 Dispatch boundary: `cli/cmd_cdd_verify.go` is 47 lines total; `Run` method is 9 lines, pure delegation. ✓
- §3.10 Shell safety: every subprocess uses argv-style `exec.CommandContext` (`cuevet.go:23`, `git.go:14,24,34`, `ledger.go:486,497`). No `sh -c`. ✓
- §3.11 Override precedence: `--repo-root` > `git rev-parse` > `os.Getwd()` documented in `run.go:188-194`. ✓

**Walked every section. No non-conformance detected.**

### AC verification

| AC | Oracle | Verdict |
|---|---|---|
| AC1 | 9 Go files under `src/packages/cnos.cdd/commands/cdd-verify/`; Python+bash removed. | ✅ |
| AC2 | `cn cdd verify --help` works; routes through Go (validator_identity emitted). No separate binary. | ✅ |
| AC3 | 37/37 oracle sub-checks pass. | ✅ |
| AC4 | 5/5 counterfeit fixtures rejected with named diagnostics. | ✅ |
| AC5 | JSON schema diff empty; emitted JSON validates structurally. | ✅ |
| AC6 | eng/go compliance: all sections walked, no non-conformance. | ✅ |
| AC7 | All flag forms preserved; ledger output format identical. | ✅ |
| AC8 | `go build`, `go test`, `go test -race`, `go vet`, `go mod tidy` all green. `cn cdd verify --help` works post-build. | ✅ |

**8/8 AC oracles pass.**

### Diff sanity

- Diff scope matches the gamma-scaffold (new go.work; new module; new kernel wrapper; main.go registration; Python+bash removed; cn.package.json edit; README + test scripts updated; CHANGELOG line; cycle artifacts). No scope creep.
- No edits to `src/go/internal/cli/dispatch.go`, `registry.go`, `command.go` (the kernel dispatch contract is untouched — only a new command was registered).
- No edits to `schemas/`.
- No edits to other packages.

### Findings

**None at R1.** No follow-up rounds required.

### Special-rigor notes

The two architectural failure classes from #389 (wrong language) and #391 (wrong location, wrong CLI) are both eliminated:

- **Wrong language (#389 class):** All V code is Go. No Python files remain.
- **Wrong location (#391 class):** V code lives at the cnos.cdd package's command space, not in `src/go/internal/`. Verified by file paths.
- **Wrong CLI (#391 class):** V is exposed as `cn cdd-verify` through the cn dispatcher (kernel-tier registration). No separate `cn-cdd-verify` top-level binary.

The implementation contract pinned at dispatch by δ-as-architect is honored axis-by-axis.

### Verdict

**APPROVE (R1).** No findings. Cycle ready for close-out.

---

β-α-collapse acknowledged in the receipt mode/structure. The collapse-mode warning surface that V itself emits (C1/C2 warnings) is the structural-trust surface the operator can inspect post-merge to confirm the collapse was honest.
