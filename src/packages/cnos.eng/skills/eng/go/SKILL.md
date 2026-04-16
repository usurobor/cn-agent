---
name: go
description: Write native Go for cnos. Use when building the runtime/kernel, package tooling, command dispatch, HTTP services, or other low-level infrastructure in Go.
artifact_class: skill
kata_surface: embedded
governing_question: How do we write Go that keeps cnos simple, explicit, testable, and easy for any agent to build anywhere?
triggers:
  - implementing cnos runtime/kernel code in Go
  - creating package/command/orchestrator infrastructure in Go
  - building HTTP/file/git/process integrations in Go
  - reviewing Go code for cnos
  - migrating logic from OCaml or shell into Go
---

# Go

## Core Principle

**Write Go that is boring to run, easy to read, and explicit about failure.**

Use packages to keep boundaries clear, concrete types to keep intent obvious, and error values to make degraded behavior visible. Prefer simple code paths that any agent can build and reason about on any machine.

## Algorithm

1. **Define** — name the model, the package boundaries, the failure policy, and the build/test contract.
2. **Unfold** — design the package layout, types, interfaces, context flow, error flow, resource lifecycle, and tests.
3. **Rules** — keep the code explicit, small-package, low-magic, and hostile to silent fallback.
4. **Verify** — build, test, and review the same failure classes across sibling surfaces.

---

## 1. Define

### 1.0. Classify the artifact first

Use this skill for Go code that belongs in the cnos runtime/kernel layer:

- package management
- command dispatch
- orchestrator runtime
- runtime contract emission
- HTTP services
- file/git/process adapters
- transport and protocol support

Do **not** use this skill for:

- product prose
- package-authored Markdown skills
- shell one-offs that should stay shell one-offs
- workflow judgment rules better expressed as package skills

### 1.1. Identify the parts

A coherent Go subsystem has:

- **Packages** — the primary unit of design
- **Types** — concrete domain models
- **Interfaces** — small consumer-owned seams
- **Errors** — explicit outcome channel
- **Context** — cancellation / deadline propagation
- **Observability** — structured `log/slog` state emission
- **Adapters** — file, git, HTTP, process, archive, env
- **Tests** — table-driven and subtest-based
- **Build** — `go build` / `go test` / `go mod tidy`, no exotic toolchain
  - ❌ Start with helpers and globals, then invent packages later
  - ✅ Name the packages and types first, then fill in behavior

### 1.2. Articulate how they fit

Packages own responsibilities. Concrete types carry domain meaning. Interfaces appear where callers need substitution. Context flows from entrypoint downward. Errors make failure explicit. Adapters isolate side effects. Tests prove both success and degraded behavior.

- ❌ One package owns CLI, HTTP, registry parsing, and git side effects together
- ✅ `cmd/cn` wires `internal/runtime`, which uses `internal/deps`, `internal/contract`, and `internal/transport`

### 1.3. Name the failure mode

Go code fails through:

- **package sprawl** — blurry boundaries, cyclic dependencies, `util` buckets
- **interface pollution** — interfaces declared before there is a real consumer need
- **context leaks** — `context` stored in structs or dropped on the floor
- **resource leaks** — failing to `defer` cleanup for files, network responses, or locks
- **opaque exhaust** — unstructured `fmt.Printf` or `log.Printf` that hides system reality
- **panic-driven control flow** — expected failure treated as exceptional collapse
- **silent fallback** — error ignored, empty value fabricated, caller never told
- **global state drift** — mutable package globals controlling behavior invisibly
- **library fatalism** — `log.Fatal` / `os.Exit` deep in reusable code
  - ❌ `if err != nil { return nil }`
  - ✅ `if err != nil { return fmt.Errorf("read package index: %w", err) }`

---

## 2. Unfold

### 2.1. Package design

**Design packages around responsibilities, not layers of convenience.**

Use short, lower-case package names. Avoid names that steal good local variable names or say nothing.

- ❌ `common`, `utils`, `helpers`, `misc`
- ✅ `deps`, `contract`, `transport`, `command`, `workflow`, `activation`

For cnos, prefer a shape like:

```text
cmd/cn/                # CLI entrypoint only
internal/deps/         # package index, lockfile, restore/install logic
internal/contract/     # runtime contract types + rendering
internal/command/      # command registry + dispatch
internal/workflow/     # orchestrator IR + runner
internal/activation/   # skill activation index + evaluator
internal/transport/    # peer/mail/A2A support
internal/host/         # fs/git/http/process/archive/env adapters
```

- ❌ Import side-effecting adapters directly from every package
- ✅ Keep impure adapters at the edge and pass them inward through concrete types or narrow interfaces

### 2.2. Type design

Prefer concrete types. Make the zero value useful when practical. Use structs and typed constants to model the domain directly.

```go
// ✅ concrete domain model
type SourceKind string

const (
    SourceReleaseHTTP SourceKind = "release-http"
    SourceLocal       SourceKind = "local"
)

type PackageLock struct {
    Name    string
    Version string
    SHA256  string
}
```

Make the zero value useful when it does not violate invariants.

- ❌ Constructor required for a plain data holder with no invariants
- ✅ Plain struct with zero-value-safe fields

Use constructors only when you must enforce:

- invariants
- normalization
- hidden/internal fields

### 2.3. Interfaces

Define interfaces where they are consumed, not where they are produced. Use interfaces to express what a caller needs from a dependency.

- ❌ Producer package exports `type Store interface { ... }` before any consumer needs it
- ✅ Consumer package defines the minimal interface it needs

Prefer:

- concrete return types
- narrow interfaces as inputs

```go
// ✅ consumer-owned interface
type ReleaseFetcher interface {
    FetchLatest(ctx context.Context) (Release, error)
}
```

Avoid `any` / `map[string]any` as a domain model except at true JSON boundaries.

### 2.4. Context

`context.Context` is a parameter, not a field.

Rules:

- first parameter
- never nil
- do not store in structs
- pass it through long-running or blocking operations
- check `ctx.Done()` in loops and retries

```go
func (r *Registry) Resolve(ctx context.Context, name, version string) (Artifact, error) {
    // ...
}
```

- ❌ `type Client struct { ctx context.Context }`
- ❌ `func Resolve(name string) (...)`
- ✅ `func Resolve(ctx context.Context, name string) (...)`

### 2.5. Errors and fallback

Errors are values. Fallback is policy.

Return errors for expected failures. Reserve panic for impossible states or invariants broken by programmer error.

Rules:

- wrap with `%w`
- error strings start lower-case
- no trailing punctuation
- never `log.Fatal` or `os.Exit` inside reusable packages
- never silently fabricate success from failure
- explicitly use `errors.Is` and `errors.As` to inspect wrapped errors for fallback logic

```go
if err != nil {
    return fmt.Errorf("read package index: %w", err)
}

// ✅ Explicit inspection for fallback
if errors.Is(err, fs.ErrNotExist) {
    // apply intentional policy ...
}
```

Fallbacks must be explicit and visible:

- ❌ "if missing, just return empty registry"
- ✅ "if missing and local override is allowed, log/annotate the degraded path and continue by policy"

If a fallback exists, document:

- why it exists
- when it applies
- what it returns
- what signal tells the caller it happened

### 2.6. Side-effect boundaries

Keep side effects at the edge. The core logic should decide:

- what to load
- what to validate
- what to write
- what to execute

Adapters should only do:

- filesystem
- git
- HTTP
- archive
- env
- subprocess

Rules for adapters:

- Inject adapters into domain structs at initialization.
- Strictly ban package-level mutable state or `init()` wiring.
  - ❌ Parse a package manifest inside the HTTP client
  - ✅ HTTP client returns bytes; parser validates in the domain package

### 2.7. Resource lifecycles

Resources must explicitly articulate their cleanup immediately after acquisition.

Rules:

- use `defer` for `Close()`, `Unlock()`, or cleanup functions immediately after checking the acquisition error
- never assume the garbage collector handles OS-level resources

```go
resp, err := client.Do(req)
if err != nil {
    return err
}
// ✅ Teardown bonded to acquisition
defer resp.Body.Close()
```

- ❌ Acquire a file handle, then close it 50 lines later in a conditional branch
- ✅ `defer f.Close()` on the line after the nil-error check

### 2.8. Observability

State must be structured and queryable, not opaque text.

Rules:

- use `log/slog` for all runtime logging
- attach key-value attributes (e.g., `slog.String("pkg", name)`) rather than concatenating strings
- log degraded paths explicitly when a fallback policy is triggered

```go
// ✅ Structured state articulation
slog.WarnContext(ctx, "local override allowed, bypassing missing index",
    slog.String("package", pkgName),
    slog.String("fallback_path", localPath),
)
```

- ❌ `fmt.Printf("warning: package %s missing, using fallback\n", name)`
- ✅ `slog.WarnContext(ctx, "package index missing, applying local override policy", slog.String("package", name))`

### 2.9. Concurrency

Do not add goroutines because Go makes them cheap. Add them when they simplify the model or materially improve throughput.

Rules:

- own cancellation with context
- keep fan-out/fan-in explicit
- prefer sequential correctness first
- use channels when communication is the model, not as decoration
  - ❌ Spawn goroutines for every restore step with no backpressure or cancellation
  - ✅ Keep restore sequential until parallelism is justified by real latency/throughput data

### 2.10. Testing

Use `go test` as the contract.

For deterministic logic:

- table-driven tests
- subtests with `t.Run`
- golden/text fixtures under `testdata/` where useful

For adapters:

- fake servers / temp dirs / test repos
- explicit degraded-path tests

Typical commands:

```sh
go test ./...
go test -race ./...
go build ./...
```

Test both:

- success
- degraded/malformed/missing/unsupported cases
  - ❌ Only test the happy path
  - ✅ Test both the nominal path and the policy path

### 2.11. Build and shipping

Any agent on any machine should be able to build the runtime with the standard Go toolchain.

Keep the build boring:

- standard module layout
- strict module hygiene (`go mod tidy`)
- no custom code generation unless justified
- minimal external dependencies
- cross-compile only when needed, not as a daily burden

Prefer stdlib first:

- `encoding/json`
- `net/http`
- `log/slog`
- `archive/tar`
- `compress/gzip`
- `crypto/sha256`
- `os/exec`
- `path/filepath`

**Do not reimplement stdlib functions.** Before writing any helper, check whether the standard library already provides it. This is a recurring review finding:

- ❌ Hand-rolled `contains(s, sub string) bool` → ✅ `strings.Contains`
- ❌ Hand-rolled `joinArgs(args []string) string` → ✅ `strings.Join`
- ❌ Hand-rolled `envMap()` copying `os.Environ()` → ✅ `os.Getenv` at point of use

If you write a utility function, verify it doesn't exist in `strings`, `slices`, `maps`, `path/filepath`, `sort`, `bytes`, `fmt`, or `os` first.

### 2.12. Schema and compatibility

**Treat manifests, contracts, and IRs as public contracts.**

When handling:

- package manifests
- runtime contract data
- command/orchestrator registries
- lockfiles
- indexes
- workflow IR

prefer:

- additive evolution
- explicit version fields
- strict validation of required fields
- tolerant handling of unknown fields when forward compatibility matters
  - ❌ Silently coerce malformed config into defaults
  - ✅ Reject malformed input with an explicit error and enough context to repair it

If a compatibility path exists, test it explicitly.

**Sibling harnesses are contract surfaces too.** When a Go type defines a JSON schema (e.g. `pkg.Manifest` for `.cn/deps.json`), every non-Go producer of that JSON — shell test harnesses, CI workflow steps, template-emitted defaults — is a contract surface on the same schema. The schema audit must extend to those producers. Mechanical check: for each `Parse*([]byte)` added or changed, grep the repo for non-Go writers of the same JSON shape and verify they still produce the declared schema. *Derives from: 3.56.2 (#250) — `cnos.kata/lib.sh` `write_deps_json` was emitting object-syntax `packages` while the Go parser expected an array; the bug masked the mismatch until the bug itself was fixed.*

  - ❌ Audit only the Go parser when changing a schema-bearing type
  - ✅ `grep -rn '"packages"' src/` and inspect every writer (Go, shell, workflow YAML)

### 2.13. Determinism and reproducibility

Deterministic output matters. For artifacts that are:

- hashed
- compared in tests
- used in reports
- stored for later diffing
- emitted as runtime contracts

ensure:

- stable ordering
- canonical serialization
- no dependence on map iteration order
- no hidden timestamps unless intentionally included
  - ❌ Range over a map and render the result directly
  - ✅ Sort keys/items before rendering

### 2.14. Idempotence and retry safety

For state-changing operations, define:

- what happens on interruption
- whether the operation is idempotent
- what partial state is left behind
- how cleanup works
  - ❌ Write directly to the final target and hope the process completes
  - ✅ Write to temp, validate, then atomically move into place

### 2.15. Traceability and receipts

For operations like:

- package restore
- update
- command registry rebuild
- orchestrator execution
- sync / transport handling

decide explicitly:

- what is returned
- what is logged
- what is written to status/doctor/runtime contract
- what receipt or trace event exists
  - ❌ Degraded path only visible in transient stderr
  - ✅ Degraded path visible in logs and/or structured runtime state

### 2.16. Preserve cnos runtime boundaries

Keep these boundaries explicit in Go code:

- **skills** choose
- **commands** dispatch
- **orchestrators** execute
- **extensions** provide capability

Do not let core/runtime packages smear these concepts together.

- ❌ A command registry that also decides skill activation logic
- ✅ Separate registries and clear ownership

The kernel stays minimal. Package-driven commands are preferred over built-ins. Core should not accrete convenience commands casually. Mechanical first, semantic second — the kernel resolves, validates, registers, dispatches, and enforces policy. It does not do hidden interpretation or "smart" inference when explicit metadata should exist.

### 2.17. Purity boundary: Parse vs Read

*Derives from: INVARIANTS.md T-004 (source/artifact/installed explicit)*

Pure functions operate on bytes/data. IO functions operate on paths/network. Keep them separate.

```go
// Pure — no os imports, no IO
func ParseManifestData(data []byte) (*PackageManifest, error)

// IO wrapper — calls pure function
func ReadManifest(path string) (*PackageManifest, error) {
    data, err := os.ReadFile(path)
    if err != nil { return nil, err }
    return ParseManifestData(data)
}
```

This mirrors the OCaml `src/lib/` (pure) vs `src/cmd/` (IO) discipline.

- ❌ `func ParseManifest(path string)` — mixes file IO with parsing
- ✅ `func ParseManifestData(data []byte)` pure + `func ReadManifest(path string)` IO wrapper
- ❌ `internal/pkg/` importing `os`
- ✅ `internal/pkg/` has zero `os` imports; all IO enters through explicit wrappers elsewhere

**Rule:** `Parse*` takes `[]byte`, returns typed data. `Read*` takes a path, calls `Parse*`. No mixing.

### 2.18. Dispatch boundary: cli/ owns dispatch only

*Derives from: INVARIANTS.md T-002 (kernel remains minimal and trusted)*

`cli/` contains: types (`CommandSpec`, `Invocation`, `Command`, `Registry`) and thin command wrappers. Domain logic lives in domain packages.

```
internal/cli/         # dispatch only — types, registry, thin cmd_*.go wrappers
internal/doctor/      # doctor check logic
internal/pkgbuild/    # build pipeline logic
internal/restore/     # deps restore logic
internal/hubinit/     # hub creation logic
internal/hubstatus/   # status display logic
```

- ❌ `cmd_doctor.go` with 238 lines of check logic inline in `cli/`
- ✅ `cmd_doctor.go` as a 20-line wrapper calling `doctor.Run()`
- ❌ Any `cmd_*.go` in `cli/` with more than ~30 lines of domain logic
- ✅ Domain logic extracted to its own package, `cmd_*.go` delegates

**Rule:** If a command's logic exceeds a thin wrapper, extract it to a domain package. `cli/` owns dispatch. Domain packages own logic.

---

## 3. Rules

### 3.1. Package rule

Use small packages with short, specific names.

- ❌ `common`, `helpers`, `shared`
- ✅ `deps`, `contract`, `workflow`, `activation`

No cyclic imports. No package whose only meaning is "miscellaneous."

### 3.2. Interface rule

Define interfaces in the consumer package. Return concrete types by default.

- ❌ "everything behind an interface"
- ✅ interfaces only where substitution/test seams are real

### 3.3. Context rule

Pass `context.Context` explicitly as the first parameter. Do not store it in structs. Do not drop it on the floor.

### 3.4. Error and observability rule

Expected failure returns `error`. Do not panic for policy-level failure. Wrap errors with context (`%w`). Inspect errors carefully with `errors.Is`/`errors.As` before falling back. Do not swallow them. Use `log/slog` to articulate fallback state structurally.

### 3.5. Fallback rule

Every fallback is a policy decision. Allowed only when:

- the degraded behavior is intentional
- the caller can still reason about the result
- the fallback is visible through logs, return values, or explicit state
  - ❌ `return nil` on malformed config and continue
  - ✅ `return fmt.Errorf("parse activation index: %w", err)`

### 3.6. Resource rule

`defer` must immediately follow the successful acquisition of files, network bodies, and locks. Never defer before the error check.

### 3.7. Boundary rule

Parsing, validation, and planning belong in domain packages. Filesystem, git, HTTP, process, and env access belong in adapters. Dependency injection is required for adapters.

### 3.8. Build rule

Run before push:

```sh
go mod tidy
go test ./...
go build ./...
```

If the change adds concurrency or shared mutable state, also run:

```sh
go test -race ./...
```

### 3.9. Smell list

Treat these as review smells:

- `panic` for expected errors
- `log.Fatal` / `os.Exit` in libraries
- `context.Context` stored in a struct
- producer-owned interfaces with one implementation
- package-level globals or `init()` used for runtime state
- `map[string]any` as core domain model
- package `util` / `common`
- hidden package globals controlling behavior
- silent empty fallback on parse or IO error
- goroutines with no cancellation path
- unwrapped external errors
- string-matching error messages instead of `errors.Is`
- missing `defer` for `Close()` or `Unlock()`
- unstructured `fmt.Printf` or `log.Printf` for system state
- unstable map iteration in rendered/hashed output
- non-atomic writes to final targets (no temp → validate → rename)
- archive extraction without path-traversal validation
- shell command built by string concatenation
- precedence logic undocumented or untested
- command/orchestrator/extension boundaries smeared in one package
- degraded path invisible outside transient stderr

### 3.10. Shell and archive safety

When using subprocesses or archives:

- do not shell-escape by string concatenation
- prefer argv-based execution (`exec.Command` with separate args)
- validate destination paths before extraction
- reject archive entries that escape the target directory
- never log secrets
  - ❌ `sh -c "curl ... $USER_INPUT ..."`
  - ✅ `exec.CommandContext(ctx, "curl", "-fsSL", url)`

### 3.11. Override precedence must be explicit

If behavior may come from:

- flags
- environment
- config
- runtime discovery
- defaults

define the precedence clearly and test it.

- ❌ Precedence inferred from code order
- ✅ Documented and table-tested precedence

---

## 4. Verify

### 4.1. File-level check

- Is the package boundary clear?
- Is the type model concrete?
- Is the failure policy explicit?

### 4.2. Boundary check

- Is IO at the edge?
- Are adapters injected rather than globally accessed?
- Are adapters separated from parsing/validation/planning?

### 4.3. Error, resource, and observability check

- Are all expected failures returned as `error`?
- Are fallbacks explicit, inspected via `errors.Is`, and justified?
- Are resources reliably released via `defer`?
- Is state emitted via `log/slog` (not `fmt.Printf`)?

### 4.4. Context/concurrency check

- Does context flow correctly?
- Are goroutines owned, bounded, and cancellable?

### 4.5. Schema/compatibility check

- Are manifests and contracts validated strictly on required fields?
- Are unknown fields tolerated where forward compatibility matters?
- Is schema version checked explicitly?

### 4.6. Determinism check

- Is output ordering stable (sorted, not map-iteration-dependent)?
- Are hashes computed over canonical content?

### 4.7. Safety check

- Are subprocesses argv-based (no string concatenation)?
- Are archive extractions path-validated?
- Are secrets excluded from logs?
- Is override precedence documented and tested?

### 4.8. cnos boundary check

- Are command/orchestrator/extension/skill boundaries preserved?
- Is the kernel minimal (no casual built-in accretion)?
- Are degraded paths visible in structured state, not just stderr?

### 4.9. Build/test check

- `go mod tidy`
- `go test ./...`
- `go build ./...`
- `go test -race ./...` where warranted

---

## 5. Kata

### 5.1. Kata A — package restore over HTTP

**Scenario:** Implement first-party package restore from an artifact index and tarball.

**Task:** Write Go code that:

- resolves name + version
- downloads artifact + checksum
- verifies SHA-256
- extracts into vendor
- returns explicit errors for all malformed/missing cases

**Governing skills:** go, design, testing

**Inputs:**

- package index fixture
- checksum fixture
- tarball fixture
- malformed checksum/index cases

**Expected artifacts:**

- concrete package types
- no producer-owned interfaces
- explicit adapter boundaries (struct injection)
- tests for happy path and malformed checksum/index cases

**Verification:**

- `go test ./...`
- no silent fallback from malformed index or checksum mismatch
- no panic on network/IO/parse failure

**Common failures:**

- hidden global HTTP client
- parse logic inside transport layer
- panic on malformed metadata
- missing degraded-path tests

**Reflection:**

- Did the restore path stay mechanical and deterministic?
- Did every fallback remain explicit and visible?

### 5.2. Kata B — runtime contract builder

**Scenario:** Extract runtime contract model/rendering out of a large command handler.

**Task:** Move:

- types
- validation
- rendering helpers

into a reusable Go package, leaving command wiring in `cmd/cn`.

**Governing skills:** go, design, writing

**Inputs:**

- one existing command handler with mixed model/render/wiring logic
- one expected JSON output
- one expected text output

**Expected artifacts:**

- clear package boundary
- concrete types
- explicit command wiring
- tests for JSON and text rendering

**Verification:**

- no system-truth logic left stranded in the CLI entry package
- no cyclic imports
- output stays stable under tests

**Common failures:**

- moving IO into the model package
- inventing interfaces where concrete types suffice
- leaving rendering behavior untested

**Reflection:**

- Did the extraction reduce command-layer ownership of system truth?
- Did the new package keep IO at the edge?
