# α design notes — cycle/392

## 1. cn CLI dispatch & how new commands register

Read against `src/go/cmd/cn/main.go` and `src/go/internal/cli/`.

**Command resolution flow** (`internal/cli/dispatch.go` lines 21–82):
1. Noun-verb form: `args[0] + "-" + args[1]` looked up in registry (e.g. `cn cdd verify` → `cdd-verify`).
2. Flat form: `args[0]` looked up directly (e.g. `cn doctor`).
3. Group prefix: if `args[0]` matches `<noun>-*` with no matching verb, list the group.

**Registry tier-based precedence** (`internal/cli/registry.go`):
- `TierKernel = 0` (highest), `TierRepoLocal = 1`, `TierPackage = 2`.
- Lower tier wins. Registering a kernel-tier `cdd-verify` supersedes any package-tier `cdd-verify`.

**Three existing surfaces a new command can use**:
- (a) Kernel command registered in `cmd/cn/main.go` (e.g. `cli.DoctorCmd`, `cli.BuildCmd`).
- (b) Package-vendored command via `cn.package.json` declaration → `discover.ScanPackageCommands` → `ExecCommand` execs the entrypoint.
- (c) Repo-local command at `.cn/commands/cn-<name>` → `discover.ScanRepoLocalCommands`.

**Cycle #392 path** (per δ pin): kernel command (a). Rationale: the pin says "integrated through `cn` CLI dispatch, NOT a separate binary." A package-vendored command (b) is dispatched via `ExecCommand` exec'ing a separate binary, which is exactly what the pin forbids. A kernel command compiles into the cn binary itself.

This requires a Go workspace because the V source must live at `src/packages/cnos.cdd/commands/cdd-verify/` (package-scoping pin) but the cn binary lives at `src/go/cmd/cn/` (single Go module today). Workspaces solve exactly this: separate modules in one build graph.

## 2. Architecture & package layout

```
src/packages/cnos.cdd/commands/cdd-verify/
├── go.mod                      # module github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify
├── README.md                   # updated for Go
├── exceptions.example.yml      # unchanged
├── cddverify.go                # public API: ValidateReceipt, ValidationVerdict types, JSON marshal
├── dispatch.go                 # protocol_id → schema-package dispatch table
├── counterfeit.go              # C1-C5 counterfeit rules
├── verdict.go                  # ValidationVerdict types + JSON marshaling
├── parse.go                    # Pure: ParseReceiptData([]byte) → Receipt (eng/go §2.17)
├── read.go                     # IO: ReadReceipt(path) → calls ParseReceiptData
├── cuevet.go                   # adapter: invokes `cue vet` (eng/go §2.6 side-effect boundaries)
├── git.go                      # adapter: git log -1 --format=%ae / %cI
├── ledger.go                   # legacy artifact-presence checks (--unreleased/--all/--version/--pr/--cycle/--triadic)
├── run.go                      # top-level Run(ctx, args, stdout, stderr) → exit code
├── cddverify_test.go           # table-driven tests (eng/go §2.10)
└── test/                       # existing test fixtures (kept; bash test harness)

src/go/internal/cli/cmd_cdd_verify.go    # thin kernel-command wrapper (eng/go §2.18)
src/go/cmd/cn/main.go                    # register cli.CddVerifyCmd

go.work                                  # repo-root workspace file
```

### Go workspace
```go
// go.work
go 1.24

use (
    ./src/go
    ./src/packages/cnos.cdd/commands/cdd-verify
)
```

This is the cleanest way to honor the package-scoping pin without inflating `src/go/internal/`. Reference: eng/go §2.11 "Keep the build boring: standard module layout, strict module hygiene (`go mod tidy`)" — workspaces are a stdlib-blessed feature since go 1.18 and are boring/standard.

### Package boundaries (eng/go §2.1, §3.1)
- Single domain package `cddverify` owning V's responsibility.
- Adapters (cue, git) are internal helpers within the same package since the package is small (~700 LOC); injected via concrete types (function values) per eng/go §2.6 / §2.3, NOT producer-owned interfaces.
- The kernel-command wrapper (`internal/cli/cmd_cdd_verify.go`) is ≤30 lines — pure delegation to `cddverify.Run` (eng/go §2.18).

### Concrete types (eng/go §2.2)
```go
type Receipt map[string]any  // YAML-parsed receipt; at the JSON boundary, map[string]any is permitted per §2.3.
type FailedPredicate struct {
    Predicate   string `json:"predicate"`
    EvidenceRef string `json:"evidence_ref,omitempty"`
    Diagnostic  string `json:"diagnostic"`
}
type Provenance struct {
    ValidatorIdentity string    `json:"validator_identity"`
    ValidatorVersion  string    `json:"validator_version"`
    CheckedAt         string    `json:"checked_at"`
    InputRefs         InputRefs `json:"input_refs"`
}
type InputRefs struct {
    ContractRef     string `json:"contract_ref,omitempty"`
    ReceiptRef      string `json:"receipt_ref"`
    EvidenceRootRef string `json:"evidence_root_ref,omitempty"`
}
type ValidationVerdict struct {
    Result           string            `json:"result"`           // PASS|FAIL
    FailedPredicates []FailedPredicate `json:"failed_predicates"`
    Warnings         []string          `json:"warnings"`
    Provenance       Provenance        `json:"provenance"`
}
type DispatchEntry struct {
    Definition  string   // e.g. "#Receipt"
    SchemaFiles []string // repo-relative paths
}
```
Note: `Receipt` is `map[string]any` only at the YAML/JSON boundary (eng/go §2.3 explicitly permits this at true JSON boundaries). All counterfeit rules access fields through typed accessor helpers.

### Errors (eng/go §2.5, §3.4)
- Functions return `error`; wrap with `%w` and lower-case prefixes.
- No `log.Fatal`/`os.Exit` outside `main` (cmd/cn already owns that).
- `errors.Is` used where checking for missing files vs other failures.
- V's exit codes preserved from Python design: 0=PASS, 1=FAIL, 2=internal error. Encoded in `Run`'s return value.

### Context (eng/go §2.4, §3.3)
- All exported entrypoints take `ctx context.Context` first.
- `cue vet` and `git` subprocess invocations use `exec.CommandContext(ctx, ...)`.

### Resources (eng/go §2.7, §3.6)
- `os.ReadFile` is the standard mode for receipt loading — no manual file handles.
- Subprocess `cmd.Run()` doesn't acquire a long-lived resource.
- No explicit `defer Close()` needed since we avoid streaming IO.

### Side-effect boundaries (eng/go §2.6, §3.7)
- `parse.go` is pure (no `os` imports, only `[]byte` in, typed values out) per eng/go §2.17 "Purity boundary: Parse vs Read."
- `read.go` is the IO wrapper.
- `cuevet.go` and `git.go` shell out via `exec.CommandContext` with argv (no string concatenation) per eng/go §3.10.

### Observability (eng/go §2.8)
- The V command's stdout/stderr is operator-facing prose or JSON, not structured logs. Per eng/go §2.8 the runtime logs use `log/slog`. The V command is a one-shot CLI that emits its verdict — slog is appropriate for fallback/degraded-path notes that don't belong in the verdict itself.

### Schema/compatibility (eng/go §2.12, §3.7)
- ValidationVerdict JSON schema is the public contract; Go types map 1:1.
- Receipt YAML parsing: tolerant of unknown fields (`map[string]any`); strict on required fields enforced by CUE.

### Determinism (eng/go §2.13)
- `failed_predicates` order: deterministic — same order Python emits (C5 → C3 → C4 → C1 → C2; CUE diagnostics emitted in lex order from cue's stderr).
- `checked_at` is intentionally non-deterministic (timestamp) per the schema.

### Tests (eng/go §2.10, §3.8)
- Table-driven Go tests for each counterfeit rule, JSON emission, dispatch table.
- Integration via the existing 37-check bash oracle (`tests/cdd/test_cn_cdd_validate_receipt.sh`).
- Run `go test ./...` and `go build ./...` at AC8 verification.

## 3. Python → Go function mapping

| Python (cn-cdd-validate-receipt) | Go (`cddverify` pkg) | Notes |
|---|---|---|
| `validate(receipt_path, contract_path, repo_root, structural_only)` | `Validate(ctx, ReceiptPath, ContractPath, RepoRoot, StructuralOnly) (ValidationVerdict, error)` | top-level entrypoint |
| `load_yaml(path)` | `ReadReceipt(path) → ([]byte → Receipt)`; uses `yaml.v3` via `gopkg.in/yaml.v3` — but eng/go §2.11 says "minimal external dependencies" and `runtime deps: Go + cue` from the issue. **Use `cue export` fallback** (already in Python) to avoid adding `gopkg.in/yaml.v3` |
| `cue_vet(receipt_path, dispatch_entry, repo_root)` | `cueVet(ctx, receipt, entry, repoRoot) (ok bool, diagnostics []string, err error)` | uses `exec.CommandContext` |
| `rule_c1_actor_separation` | `ruleC1ActorSeparation(receipt, repoRoot, &fail, &warn)` | unchanged logic; uses git log adapter |
| `rule_c2_verdict_precedes_merge` | `ruleC2VerdictPrecedesMerge(...)` | unchanged |
| `rule_c3_override_does_not_rewrite` | `ruleC3OverrideDoesNotRewrite(...)` | unchanged |
| `rule_c4_evidence_deref` | `ruleC4EvidenceDeref(...)` | unchanged |
| `rule_c5_protocol_id_mismatch` | `ruleC5ProtocolIDMismatch(...)` | unchanged |
| `is_collapsed_mode(receipt)` | `isCollapsedMode(receipt)` | unchanged |
| `is_filesystem_path(s)` | `isFilesystemPath(s)` | unchanged |
| `_extract_actor(file_path)` | `extractActor(filePath)` | unchanged |
| `_git_log_author(repo_root, file_path)` | `gitLogAuthor(ctx, repoRoot, filePath)` | adapter |
| `_git_commit_time(repo_root, file_path)` | `gitCommitTime(ctx, repoRoot, filePath)` | adapter |
| `_parse_iso(ts)` | `parseISO(ts)` | uses `time.Parse(time.RFC3339, …)` |
| `emit_verdict(...)` | `emitVerdict(...)` | builds the struct |
| `render_prose(verdict, receipt_path)` | `RenderProse(verdict, receiptPath) string` | text output |
| Legacy ledger (bash `cn-cdd-verify`) | `runLedger(ctx, args, stdout)` | NEW Go port of artifact-presence checks (--unreleased/--all/--version/--pr/--cycle/--triadic) |

## 4. YAML parsing decision

**Pinned runtime deps: Go + cue only. No new deps.** Therefore: do NOT add `gopkg.in/yaml.v3`.

Reuse the Python design: invoke `cue export --out json <yaml-path>` to coerce YAML → JSON, then `encoding/json` decode into `map[string]any`. This adds no new dependencies and reuses cue (already required).

Implementation: `parse.go::ParseReceiptYAML(ctx, path) (Receipt, error)` invokes `cue export` via `exec.CommandContext`. This is technically not a pure parser (it shells out), so it lives in `parse.go` as `ParseReceiptYAML` (impure but acceptable: it's the YAML adapter). To keep eng/go §2.17 split clean, also expose `ParseReceiptJSON([]byte) (Receipt, error)` as a pure function over JSON bytes — used by tests that supply JSON directly.

## 5. Test harness adaptation

The existing `tests/cdd/test_cn_cdd_validate_receipt.sh` already invokes `$V` with `--receipt …` etc. It currently points `V` to `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` (the bash wrapper). We change `V` to invoke the new path: the cn binary with the `cdd-verify` subcommand.

Build the cn binary into a known location (e.g. `dist/cn` or `tmp/`) and set `V="<built-cn-path> cdd-verify"`. Use `bash -c` invocation in the script.

Backwards-compatibility: bash arrays are messy for `$V` if it has spaces. Solution: in the test harness, define `V` as the cn binary path and use `"$V" cdd-verify` everywhere — minimal patch.

Alternative: introduce a `CN_BIN` env var the script reads, defaulting to a built path. Cleaner.

## 6. cnos.cdd `cn.package.json` change

Currently declares `cdd-verify` → entrypoint `commands/cdd-verify/cn-cdd-verify`. Since the kernel command supersedes (TierKernel < TierPackage), the package declaration could remain. But the file the entrypoint points to is being removed. The validation in pkgbuild (`build.go:531-557`) would then report "entrypoint does not exist" → build failure.

**Decision**: remove the `cdd-verify` entry from `src/packages/cnos.cdd/cn.package.json`. Keep `cdd-status`. Kernel-tier `cdd-verify` is the canonical surface.

## 7. CHANGELOG entry

Add a line under the unreleased / 3.81.0 section noting V ported to Go.

## 8. eng/go SKILL compliance walk (preview — full evidence in self-coherence.md)

| Section | How honored |
|---|---|
| §1.0 Classify the artifact | V is runtime/kernel-adjacent Go code; correct choice |
| §1.1 Identify the parts | packages, types, errors, ctx, adapters, tests all named |
| §1.3 Failure modes | covered: no package sprawl; no producer interfaces; no panic; no silent fallback; no globals |
| §2.1 Package design | single `cddverify` package + thin `cli` wrapper |
| §2.2 Type design | concrete struct types; zero-value safe; map[string]any only at JSON boundary |
| §2.3 Interfaces | no producer-owned interfaces; function-value injection for adapters |
| §2.4 Context | ctx first parameter throughout |
| §2.5 Errors | %w wrap; lowercase; no log.Fatal in library |
| §2.6 Side-effect boundaries | cuevet/git adapters separate from logic |
| §2.7 Resource lifecycles | use `os.ReadFile`; no long-held handles |
| §2.8 Observability | slog for degraded paths; verdict output is the contract |
| §2.10 Testing | table-driven tests for rules and verdict emission |
| §2.11 Build/shipping | go workspace; stdlib first; no new deps |
| §2.12 Schema/compat | JSON schema preserved; struct tags pinned |
| §2.13 Determinism | rule firing order is fixed; sorted CUE diagnostics |
| §2.14 Idempotence | V is read-only; no state changes |
| §2.16 cnos boundaries | command-only logic in cddverify; cli/ owns dispatch wrapping |
| §2.17 Purity boundary | ParseReceiptJSON pure; ReadReceipt IO wrapper |
| §2.18 Dispatch boundary | cli/cmd_cdd_verify.go ≤30 lines, pure delegation |
| §3.10 Shell safety | exec.CommandContext with separate args |

(Full citations with file:line in self-coherence.md after build.)

## Decisions log

- **D1**: Kernel command (not package-vendored) — required by CLI integration pin.
- **D2**: Go workspace — required to honor package-scoping pin while building into cn binary.
- **D3**: YAML via `cue export` subprocess (no new deps) — required by runtime-deps pin.
- **D4**: Remove the `cdd-verify` entry from cnos.cdd `cn.package.json` — required because kernel ownership.
- **D5**: Bash wrapper fully removed (operator preference); Python script fully removed.
- **D6**: Legacy ledger logic ported into Go (not deferred) — backward-compat pin requires `--unreleased`/`--all`/`--version`/`--pr` all still work.
