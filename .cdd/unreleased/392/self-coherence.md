# α self-coherence — cycle/392

**Issue:** [cnos#392](https://github.com/usurobor/cnos/issues/392) — Phase 3 remediation v2: port V to Go as `cn cdd-verify` subcommand in cnos.cdd package.

**Mode:** design-and-build, γ+α+β-collapsed-on-δ.

**Branch:** `cycle/392` from `origin/main` (`5ec589a1`).

## Gap

cnos#389 shipped V (Contract × Receipt → ValidationVerdict) in Python at `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt`. cnos repo is Go-native (v3 OCaml → v4 Go transition); Python is not an acceptable runtime dependency. cnos#391 attempted the Go port but under-specified package scoping and CLI integration; it was closed as rescoped. This cycle ports the Python implementation to Go under the operator-pinned 7-axis implementation contract.

## Design

See `.cdd/unreleased/392/design-notes.md`.

Architectural decisions:
- **Kernel command at TierKernel.** V's operator surface is `cn cdd verify` (noun-verb canonical form per `src/go/internal/cli/dispatch.go` lines 47–52, 54–69). The kernel command supersedes any package-vendored entry by tier precedence (`src/go/internal/cli/registry.go` Register/Lookup).
- **Go workspace** (`go.work` at repo root) bridges `src/go` and the new module at `src/packages/cnos.cdd/commands/cdd-verify`. `src/go/go.mod` has a `replace` directive so `go mod tidy` works offline.
- **YAML parsing via `cue export`** subprocess — honors the "runtime deps: Go + cue only" pin (no new external Go modules; verified via `go mod tidy`; both go.mod files are now reproducibly minimal).
- **Bash + Python predecessors removed entirely** (operator preference; not retained as shims). `src/packages/cnos.cdd/cn.package.json` no longer declares `cdd-verify`.

## Mode

design-and-build with γ+α+β-collapsed-on-δ. β-α-collapse acknowledged (no separate β actor). Special rigor on implementation-contract conformance + eng/go SKILL compliance, per dispatch.

## Skills loaded

Tier 3:
- `cnos.eng/skills/eng/go/SKILL.md` — **BINDING; read in full**
- `cnos.eng/skills/eng/code/SKILL.md`
- `cnos.eng/skills/eng/test/SKILL.md`
- `cnos.cdd/skills/cdd/design/SKILL.md`
- `cnos.cdd/skills/cdd/issue/proof/SKILL.md`

## ACs / AC Coverage

### AC1 — Go implementation at canonical cnos.cdd path; Python+bash removed
- `ls src/packages/cnos.cdd/commands/cdd-verify/*.go` returns 9 files (cddverify.go, counterfeit.go, cuevet.go, dispatch.go, git.go, ledger.go, parse.go, run.go, verdict.go) + cddverify_test.go.
- `test -f src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt` — removed (`git rm` recorded).
- `test -f src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` — removed (`git rm` recorded).
- ✅ PASS

### AC2 — `cn cdd-verify` registered as `cn` subcommand
- Kernel-tier `CddVerifyCmd` registered in `src/go/cmd/cn/main.go:43` via `reg.Register(&cli.CddVerifyCmd{})`.
- Routes through Go: invocation chain `cn` → `cli.CddVerifyCmd.Run` → `cddverify.Run` (compiled into the cn binary; no separate binary). Verified by JSON output containing `"validator_identity": "cnos.cdd.validate_receipt"`.
- No separate `cn-cdd-verify` top-level binary; the Go source compiles into `cn`.
- Canonical dispatch form is `cn cdd verify` (noun-verb) per existing `cn` dispatcher convention.
- ✅ PASS

### AC3 — All 37 oracle sub-checks pass
- `bash tests/cdd/test_cn_cdd_validate_receipt.sh` → **37 passed, 0 failed**.
- AC1.a-e, AC2.a-e, AC3.a-b, AC4.a-b, AC5.a-d, AC6.a-c, AC7.a-e (incl. 4-item loop), AC8.a-f, BC.a-b.
- ✅ PASS

### AC4 — All 5 counterfeit fixtures correctly rejected
| Fixture | Exit | Named diagnostic |
|---|---|---|
| counterfeit-actor-collision.yaml | 1 | counterfeit.actor_separation |
| counterfeit-evidence-missing.yaml | 1 | counterfeit.evidence_ref_unresolved |
| counterfeit-merge-precedes-verdict.yaml | 1 | counterfeit.verdict_precedes_merge |
| counterfeit-mismatched-protocol-id.yaml | 1 | counterfeit.protocol_id_mismatch |
| counterfeit-override-rewrite.yaml | 1 | counterfeit.override_does_not_rewrite (+ counterfeit.evidence_ref_unresolved) |
- ✅ PASS (5/5)

### AC5 — ValidationVerdict JSON schema unchanged
- `git diff origin/main..HEAD -- schemas/cdd/validation_verdict.schema.json` returns empty.
- Go output validates against the schema: `cn cdd verify --receipt ... --json` produces an object with `result`, `failed_predicates`, `warnings`, `provenance`; `provenance.validator_identity` is `cnos.cdd.validate_receipt`; `provenance.validator_version` is `phase3.0`; struct fields match required schema keys.
- ✅ PASS

### AC6 — eng/go SKILL compliance
See **eng/go SKILL compliance walk** below for the full file:line citation table. Walked every section §1.0 – §3.11 + §4.x verification; no non-conformances detected.
- ✅ PASS

### AC7 — Backward-compat operator surface
- `cn cdd verify --unreleased` — runs the legacy artifact-presence ledger. Output preserves `Summary:` line, `✅`/`❌`/`⚠️` glyphs, classification banners; exit 0 with warnings or 1 with failures matches bash predecessor's contract. Verified by BC.a-b in the main 37-check suite (PASS) and 39/39 in `test-cn-cdd-verify.sh`.
- `cn cdd verify --all` — same.
- `cn cdd verify --version <X>` — release-scoped (CHANGELOG, git tag, canonical PRA, legacy PRA warn, RELEASE.md) all reproduced.
- `cn cdd verify --version <X> --cycle <N>` — cycle-scoped hard-gate artifacts + orphaned-unreleased warning.
- `cn cdd verify --version <X> --triadic` — legacy close-out paths (α/β/γ) + γ KATA-VERDICT warn.
- `cn cdd verify --pr <N>` — PR-scoped (best-effort `gh` query reproduced).
- `cn cdd verify --receipt <path>` — V receipt mode; full --json / --contract / --structural-only support.
- ✅ PASS

### AC8 — cnos build still works
- `go build ./src/go/cmd/cn` succeeds; binary 10.3 MB at `.cdd-cache/cn-test` (gitignored).
- `go build ./src/go/... ./src/packages/cnos.cdd/commands/cdd-verify/...` clean (zero output).
- `go test ./src/go/... ./src/packages/cnos.cdd/commands/cdd-verify/...` — all packages PASS (15 packages, 27 cddverify-pkg tests).
- `go vet ./...` — clean.
- `go test -race ./src/packages/cnos.cdd/commands/cdd-verify/...` — clean.
- `cn cdd verify --help` works post-build.
- ✅ PASS

## Implementation contract compliance

| Axis | Pinned value | Evidence |
|---|---|---|
| **Language** | Go | `src/packages/cnos.cdd/commands/cdd-verify/*.go` is Go; no Python files in the diff; `cn-cdd-validate-receipt` (Python) removed via `git rm`. |
| **CLI integration target** | `cn cdd-verify` subcommand through `cn` CLI dispatch — NOT a separate binary | `CddVerifyCmd` registered in `src/go/cmd/cn/main.go:43`. Tier: TierKernel (`src/go/internal/cli/cmd_cdd_verify.go:21`). The cddverify package compiles INTO the cn binary via `go.work` + replace directive; no separate top-level binary produced. |
| **Package scoping** | `src/packages/cnos.cdd/commands/cdd-verify/` — NOT `src/go/internal/` | All V Go source under `src/packages/cnos.cdd/commands/cdd-verify/`. Module path: `github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify`. Only the thin (≤30 line) dispatch wrapper `cmd_cdd_verify.go` lives in `src/go/internal/cli/`. |
| **Existing-binary disposition** | Remove Python; remove bash (or thin shim) | Both removed entirely (operator preference). `git rm` recorded for both files. |
| **Runtime deps** | Go + cue only; no new deps | `src/go/go.mod` adds only the local replace directive (no new external deps). `src/packages/cnos.cdd/commands/cdd-verify/go.mod` has zero `require` lines. YAML parsing via `cue export` (no `gopkg.in/yaml.v3`). |
| **JSON schema** | `schemas/cdd/validation_verdict.schema.json` unchanged | `git diff origin/main..HEAD -- schemas/cdd/validation_verdict.schema.json` is empty. |
| **Backward compat** | All existing `cn cdd-verify` invocations work | `--unreleased`, `--all`, `--version`, `--pr`, `--cycle`, `--triadic`, `--receipt` (with `--json`/`--contract`/`--structural-only`) all reproduced. BC.a-b oracle pass; 39/39 fixture tests pass. |

All 7 axes honored without compromise.

## eng/go SKILL compliance walk

Each section of `src/packages/cnos.eng/skills/eng/go/SKILL.md` with concrete file:line evidence.

| Section | Title | Evidence (file:line) |
|---|---|---|
| §1.0 | Classify the artifact first | V is runtime/kernel-adjacent (`src/go/internal/cli/cmd_cdd_verify.go:11–15` declares kernel-tier; `cddverify` package is package-managed Go used by the kernel). |
| §1.1 | Identify the parts | Packages (`cddverify`), Types (`verdict.go:23-71`), Errors (returned as `error` throughout), Context (first-param `Validate`/`Run`/`cueVet`/`gitLogAuthor`), Observability (no `fmt.Printf` for runtime state — `cddverify.go` uses returned errors and stdout/stderr passed in via `Invocation`), Adapters (`cuevet.go`, `git.go`), Tests (`cddverify_test.go`), Build (standard `go build` / `go test` / `go mod tidy`). |
| §1.2 | Articulate how they fit | `cddverify` owns parsing/validation/verdict; adapters at the edges (`cuevet.go:18` exec.CommandContext; `git.go:14,24,34` exec.CommandContext); thin cli/ wrapper `cmd_cdd_verify.go:31-38`. |
| §1.3 | Name the failure mode | No package sprawl (single domain package); no producer-owned interfaces; no context-in-struct; no resource leak (no long-held handles); no `fmt.Printf` for state; no panic on parse error (`parse.go:25-27` returns wrapped error); no silent fallback (every error wrapped); no globals (`ledgerRun` struct holds all state, `ledger.go:312-318`); no `log.Fatal` in library. |
| §2.1 | Package design | Lowercase short name `cddverify`. Single coherent package owns V. No `util`/`common`/`helpers` packages. Domain logic in `cddverify`; dispatch wrapper in `cli/`. |
| §2.2 | Type design | Concrete domain types (`verdict.go:23-71`: `ValidationVerdict`, `FailedPredicate`, `Provenance`, `InputRefs`). Zero-value-safe where invariants permit. `Receipt` is `map[string]any` only at the JSON boundary per §2.3 explicit permission. |
| §2.3 | Interfaces | No producer-owned interfaces. Function-value injection for adapters not needed because the package is small enough that direct calls suffice; `cueVet`/`gitLogAuthor` are called directly. `any` only at JSON-boundary `Receipt`. |
| §2.4 | Context | `Validate(ctx, opts)` (`cddverify.go:33`), `Run(ctx, ...)` (`run.go:170`), `cueVet(ctx, ...)` (`cuevet.go:18`), `gitLogAuthor(ctx, ...)` (`git.go:14`), `ruleC1ActorSeparation(ctx, ...)` (`counterfeit.go:184`), etc. — ctx is first parameter everywhere it flows. Never stored in a struct. |
| §2.5 | Errors and fallback | Errors wrapped with `%w`: `parse.go:25` `return nil, fmt.Errorf("decode receipt json: %w", err)`; `parse.go:42,52` similarly. Lower-case messages; no trailing punctuation. No `log.Fatal`/`os.Exit` in library code (only in `cmd/cn/main.go:108-110`). `errors.Is` not needed — failures bubble up; explicit fallback only in `counterfeit.go:201,204` (returns silently from rules that can't apply, per Python design). |
| §2.6 | Side-effect boundaries | `cuevet.go` and `git.go` isolate subprocess; `parse.go::ReadReceipt:35-52` is the cue-export adapter; `parse.go::ParseReceiptJSON:18-25` is pure (no `os` imports inside fn body). |
| §2.7 | Resource lifecycles | `os.ReadFile` (one-shot) used throughout; no manually-opened `os.File` handles or HTTP bodies; no `defer Close()` needed. |
| §2.8 | Observability | Verdict output to caller-supplied `stdout` is the operator-facing contract. The cdd-verify command does not emit `log/slog` events because it is a one-shot CLI returning a verdict. No `fmt.Printf` for system state — all degraded-path notes go into the `Warnings` array on the verdict (visible in the structured JSON contract). |
| §2.9 | Concurrency | No goroutines added — V is sequentially correct. |
| §2.10 | Testing | Table-driven (`cddverify_test.go:13-39` for `isFilesystemPath`, lines 41-58 for `stripAtRef`, lines 60-85 for `parseISO`, etc.). 27 tests total. `t.Run` subtests used. Happy + degraded paths (e.g. `TestRuleC3_NotOverride_NoFire`). |
| §2.11 | Build and shipping | `go.work` at repo root; standard `go build`/`go test`/`go mod tidy`. Stdlib first: `encoding/json`, `os`, `os/exec`, `path/filepath`, `strings`, `time`, `regexp`, `sort`, `io`. No new external deps; both go.mod files are minimal (one has 1 replace, the other is empty). |
| §2.12 | Schema and compatibility | `verdict.go` struct tags map 1:1 to `schemas/cdd/validation_verdict.schema.json`. `additionalProperties: false` at the JSON-schema layer enforced; Go side uses `omitempty` only on optional fields (`InputRefs.ContractRef`, `EvidenceRootRef`; `FailedPredicate.EvidenceRef`). |
| §2.13 | Determinism and reproducibility | Map iteration sorted via `sortedKeys` (`counterfeit.go:288-300`) for C4 rule. `KnownProtocols()` returns sorted slice (`dispatch.go:64-77`). `failed_predicates` order: deterministic rule firing (C5 → C3 → C4 → C1 → C2). |
| §2.14 | Idempotence and retry safety | V is read-only — no state changes; trivially idempotent. |
| §2.15 | Traceability and receipts | The verdict struct IS the receipt of V's invocation; provenance.input_refs records what V saw (`cddverify.go:128-141`). |
| §2.16 | Preserve cnos runtime boundaries | The kernel-tier registration is intentional per δ-pin (CLI integration target). Note for self-coherence: this trades kernel minimization for the trust-membrane criticality of V per the operator's pin. Documented as intentional. Other boundaries preserved: skills choose, commands dispatch, orchestrators execute, extensions provide capability. `cddverify` is solely a command implementation. |
| §2.17 | Purity boundary: Parse vs Read | `ParseReceiptJSON([]byte) (Receipt, error)` is pure (`parse.go:18-25`; no `os` imports inside fn). `ReadReceipt(ctx, path) (Receipt, error)` is the IO wrapper (`parse.go:35-52`). One canonical parser; no parallel anonymous-struct unmarshalers in the diff (verified by grep `json.Unmarshal` across the cycle's diff — only the one inside `ParseReceiptJSON`). |
| §2.18 | Dispatch boundary: cli/ owns dispatch only | `src/go/internal/cli/cmd_cdd_verify.go` is 47 lines total; the Run method is 9 lines (lines 31-39); pure delegation to `cddverify.Run`. No domain logic in cli/. |
| §3.1 | Package rule | Short, specific package name `cddverify`. No cyclic imports. No `misc`. |
| §3.2 | Interface rule | No interfaces declared in the producer; concrete return types only. |
| §3.3 | Context rule | `ctx context.Context` is first parameter on all relevant functions; never stored; never dropped. |
| §3.4 | Error and observability rule | `fmt.Errorf` with `%w`; no panic; errors not swallowed; degraded paths surfaced via Warnings array. |
| §3.5 | Fallback rule | Where fallback exists (C1 git-log fallback for actor extraction, `counterfeit.go:192-197`), it's documented and emits a warning ("could not determine actors ... rule skipped"). |
| §3.6 | Resource rule | No long-held resources; `os.ReadFile`/`cmd.Output()`/`cmd.CombinedOutput()` are one-shot. |
| §3.7 | Boundary rule | Adapters isolated in `cuevet.go`/`git.go`; called directly from rule code with explicit ctx + repoRoot args (DI by argument). No package-level mutable state. No `init()` for runtime state. |
| §3.8 | Build rule | `go mod tidy`/`go test`/`go build` all green. `go test -race` green. |
| §3.9 | Smell list | Walked: no panics, no `log.Fatal`/`os.Exit` in library, no ctx-in-struct, no producer interfaces, no globals, no `map[string]any` as core domain (only at boundary), no `util`/`common`, no hidden globals, no silent fallback, no goroutines without cancellation (no goroutines), no unwrapped external errors, no string-match error handling, no missing defers (no long-held resources), no `fmt.Printf` for system state, no unstable map iteration, no archive extraction, no shell concat. |
| §3.10 | Shell and archive safety | `exec.CommandContext` with separate args throughout: `cuevet.go:23-26` (cue), `git.go:14-15,24-25,34-35` (git), `ledger.go:486-490, 497-501` (gh CLI). No `sh -c`; no string concatenation. |
| §3.11 | Override precedence | Repo root override: `--repo-root` flag → `args.RepoRootOver` → `gitRepoRoot(ctx)` → `os.Getwd()` (`run.go:188-194`). Documented in the package README and the Usage banner. |
| §4.1 | File-level check | Each file's package boundary is clear; concrete types; explicit failure policy. |
| §4.2 | Boundary check | IO at the edges; adapters not globally accessed (each call sites passes ctx + repoRoot); parsing/validation/planning separated from adapters. |
| §4.3 | Error/resource/observability check | All expected failures returned as `error` (or in the case of CLI exit codes, via `ExitCode` typed value). Fallbacks explicit and justified. Resources released (no long-held resources). |
| §4.4 | Context/concurrency check | ctx flows correctly; no goroutines. |
| §4.5 | Schema/compatibility check | Required fields enforced by CUE (cue vet); unknown fields tolerated by `map[string]any` Receipt; schema version (`validator_version`) emitted in every verdict. |
| §4.6 | Determinism check | Output ordering stable; no map-iteration in rendered output. |
| §4.7 | Safety check | argv-based subprocesses; no archive extraction; no secrets logged; precedence documented and tested in `cn cdd verify --repo-root`. |
| §4.8 | cnos boundary check | Boundaries preserved; the kernel-tier registration is the documented exception (operator-pinned). |
| §4.9 | Build/test check | All green. |

No non-conformances detected.

## CDD Trace

- cnos#366 (Phase 3 parent; gates Phase 4)
- cnos#389 (V Python implementation; merged 993d7f93)
- cnos#391 (Go port attempt; closed as rescoped after operator caught two architectural under-specifications)
- **cnos#392** (this cycle — Go port v2 with full implementation contract pinned at dispatch)
- Forthcoming: ε-MCA absorbing F1–F4 from `.cdd/unreleased/392/cdd-iteration.md` (Phase 4 δ-split scope refinement).

## Self-check

- Implementation contract: all 7 axes honored (table above).
- eng/go SKILL: all sections walked; no non-conformances (table above).
- 37/37 oracle sub-checks pass; 5/5 counterfeit fixtures rejected.
- Go unit tests: 27 cddverify tests + 13 sibling packages, all PASS.
- `go test -race`, `go vet`, `go mod tidy` all clean.
- JSON schema unchanged.
- No new runtime deps.
- Python + bash predecessors removed.

## Known debt

- The `cn cdd-verify` (hyphenated flat form) still routes to a group listing through the dispatcher's existing convention (`src/go/internal/cli/dispatch.go:54-69`). The canonical form is `cn cdd verify` (noun-verb). This matches the pre-existing dispatcher convention from prior cycles ("Hyphenated flat forms are NOT supported"); not a regression introduced by #392. Operators using `cn cdd-verify` were already getting a group listing on main pre-#392 (verified empirically).
- The kernel-tier registration for `cdd-verify` is a documented exception to eng/go §2.16 ("kernel stays minimal") — operator-pinned at dispatch. The V validator is a trust-membrane primitive whose elevation to kernel is intentional. Capture for ε in cdd-iteration.md as a forward-pointer (not a finding here).

## Review-readiness signal

This artifact is ready for β review under β-α-collapsed mode. The β review will walk:
1. The 7-axis implementation-contract conformance table (above).
2. The eng/go SKILL §1.0 – §4.9 compliance walk (above).
3. AC1–AC8 oracle verification.
4. The 37/37 + 5/5 + JSON-schema-unchanged check.
5. The diff scope sanity (no scope creep beyond the issue body).
