---
name: ocaml
description: Write native OCaml for cnos. Use when building the runtime/kernel, package tooling, command dispatch, contract/rendering logic, transport support, or other low-level infrastructure in OCaml.
artifact_class: skill
kata_surface: embedded
governing_question: How do we write OCaml that keeps cnos explicit, pure where possible, and hostile to silent fallback?
triggers:
  - implementing cnos runtime/kernel code in OCaml
  - extracting pure model logic out of command handlers
  - building package/command/orchestrator infrastructure in OCaml
  - reviewing OCaml code for purity leaks or silent fallback
  - debugging type ambiguity, fallback drift, or resource discovery bugs
---

# OCaml

## Core Principle

**Write OCaml so the compiler carries the burden early and the runtime makes degraded behavior visible.**

Use types to remove ambiguity at definition time, modules to enforce purity boundaries, Result for expected failure, and explicit fallback policy for every degraded path. Prefer a small number of clear representations over cleverness that moves ambiguity downstream.

## Algorithm

1. **Define** — name the domain model, the purity boundary, the fallback policy, and the build/test contract.
2. **Unfold** — design types, modules, result/error flow, resource discovery, determinism, and tests together.
3. **Rules** — keep the code type-first, purity-aware, explicit about fallback, and hostile to fabricated empty state.
4. **Verify** — build, run expect tests, and audit the same failure class across sibling surfaces.

---

## 1. Define

### 1.0. Classify the artifact first

Use this skill for OCaml code that belongs in the cnos runtime/kernel layer:

- package management
- command dispatch
- runtime contract emission
- activation/registry logic
- transport/protocol support
- orchestrator/workflow runtime
- file/git/process/http adapters
- contract/rendering logic

Do **not** use this skill for:

- prose-only documentation
- package-authored Markdown skills
- shell one-offs that should stay shell one-offs
- workflow judgment better expressed as package skills

### 1.1. Identify the parts

A coherent OCaml subsystem has:

- **Types** — the primary safety mechanism
- **Modules** — the primary purity and ownership boundary
- **Results** — the normal channel for expected failure
- **Fallback policy** — the explicit rule for degraded behavior
- **Resource discovery** — path/binary/config/package lookup with visible outcomes
- **Determinism** — stable ordering, stable rendering, stable hashing inputs
- **Observability** — logs/receipts/structured state for degraded paths
- **Tests** — expect tests and focused pure-logic checks
- **Build** — dune build, dune runtest, native binary discipline

- ❌ Start with helpers and exceptions, then discover the model later
- ✅ Name the types and module boundaries first, then implement behavior

### 1.2. Articulate how they fit

Types define the domain. Modules isolate purity and ownership. Result carries expected failure. Fallbacks are policy decisions, not convenience. Resource discovery validates reality before work depends on it. Deterministic rendering and stable ordering keep reports/contracts comparable. Tests lock behavior to visible output.

- ❌ Put I/O in a pure parsing module
- ✅ Parse returns `Result`; caller performs I/O

### 1.3. Name the failure mode

OCaml code fails through:

- **type ambiguity** — overlapping record fields create annotation debt everywhere
- **purity leaks** — Unix/Sys/IO code drifts into model/parse modules
- **silent fallback** — missing, malformed, or unreadable becomes `[]`, `None`, or `""`
- **exception abuse** — expected failures become implicit jumps
- **resource discovery drift** — paths/binaries are constructed but not validated
- **nondeterministic output** — unstable ordering leaks into contracts, reports, hashes
- **invisible degraded paths** — fallback happened but no log/receipt/state makes it inspectable

- ❌ `with _ -> []`
- ✅ `Error "resource not found: ..."` or explicit logged fallback by policy

---

## 2. Unfold

### 2.1. Type design

Disambiguate at definition, not at use. When two record types in the same module share a field name, every access site pays the annotation cost. Avoid that by naming the fields differently at the type definition.

```ocaml
(* ❌ recurring disambiguation burden *)
type extension_op = { kind : string; op_class : string }
type backend = { kind : string; command : string list }

(* ✅ ambiguity removed at definition *)
type extension_op = { op_kind : string; op_class : string }
type backend = { backend_kind : string; command : string list }
```

- ❌ Rely on incremental type annotations every time the compiler complains
- ✅ Fix the representation once so the whole module gets simpler

Prefer:

- concrete algebraic data types
- records with semantically specific field names
- variants for explicit state
- small helper constructors only when invariants need them

### 2.2. Module structure

Use modules to enforce purity and ownership. Prefer a split like:

- `lib/` — pure types, parsing, validation, rendering
- `protocol/` — pure FSMs / protocol state
- `ffi/` — system bindings / side-effectful primitives
- `cmd/` — business logic wiring pure and impure pieces
- `transport/` — adapters for git/mail/network
- `cli/` — thin entrypoint/dispatch only

- ❌ One module owns parsing, I/O, rendering, and policy
- ✅ Pure modules decide; impure modules execute

### 2.3. Results, exceptions, and fallback

Use `Result` for expected failure. Expected failures include:

- malformed manifests
- missing config
- checksum mismatch
- unsupported target
- absent resource by policy

Reserve exceptions for:

- impossible internal invariants
- programmer errors
- truly exceptional system conditions

```ocaml
let parse s : (value, string) result = ...
```

- ❌ `raise Parse_error`
- ✅ `Error "expected package version"`

Every fallback must be:

- intentional
- visible
- justified

### 2.4. Resource discovery

Any function that resolves binaries, files, package contents, transport roots, or thread stores must validate existence or return a visible failure.

- ❌ `let resolve name = Filename.concat dir name`
- ✅ `if Sys.file_exists p then Ok p else Error (Printf.sprintf "not found: %s" p)`

Do not defer discovery of missing resources to a later opaque exec/IO failure.

### 2.5. Compatibility and fallback discipline

Every compatibility path must be explicitly classified:

- **remove** — delete it
- **convert** — log + fallback
- **keep** — justify it, make it fail-closed or operator-visible, and give it an expiry/follow-up where appropriate

- ❌ "Keep for safety"
- ✅ "Keep until all peers use packet refs; reject loudly on ambiguity"

If "empty/skip" is the safe semantic outcome, fallback may be tolerated, but it must still be:

- logged or surfaced
- bounded
- auditable

### 2.6. Determinism and reproducibility

For anything that is:

- hashed
- diffed in tests
- emitted in reports
- stored in runtime contracts
- compared across runs

ensure:

- sorted ordering
- canonical rendering
- stable traversal
- no hidden timestamps unless deliberate

- ❌ Use `Sys.readdir` order directly in a contract/report
- ✅ `Sys.readdir dir |> Array.to_list |> List.sort String.compare`

### 2.7. Idempotence, retry, and receipts

Assume package restore, update, sync, and workflow steps may be retried. For state-changing operations, define:

- interruption behavior
- cleanup behavior
- whether the operation is idempotent
- what receipt/log/state proves what happened

Patterns to prefer:

- temp + validate + atomic move
- lock or exclusive create where needed
- receipt/log/state update on degraded path

- ❌ Write directly into final state and hope the process completes
- ✅ Write to temp, verify, then atomically install

### 2.8. Side-effect boundaries

Keep parsing/validation/planning separate from:

- filesystem
- git
- HTTP
- subprocesses
- archive extraction
- env access

- ❌ Parse a package manifest inside the fetcher or subprocess wrapper
- ✅ Fetch bytes / run process in the adapter; parse/validate in the domain module

### 2.9. Patterns

Prefer:

```ocaml
(* Pipeline style *)
input |> parse |> validate |> output

(* Result-based parsing *)
let parse s : (value, string) result = ...

(* Pattern matching over partial stdlib access *)
match xs with
| x :: _ -> x
| [] -> default

(* filter_map — filter + map in one pass *)
List.filter_map extract items

(* Localized mutable state hidden behind a pure API *)
type state = { src : string; mutable pos : int }

(* RAII cleanup *)
Fun.protect ~finally:(fun () -> close_in_noerr ic) (fun () -> ...)

(* argv-based subprocess, not shell concatenation *)
let code, output = exec_args ~prog:"curl"
  ~args:["--config"; "-"] ~stdin_data:config ()
```

Avoid:

- partial functions (`List.hd`, `Option.get`)
- `with _ ->`
- exceptions for expected domain failure
- mutable state escaping module-local scope
- shell interpolation for subprocesses

### 2.10. Testing

Use:

- `ppx_expect` / inline expect tests for behavior contracts
- pure logic tests close to the module
- degraded-path tests for fallback/compatibility/resource discovery

Test both:

- positive path
- missing/malformed/degraded path

- ❌ Only test the happy path
- ✅ Test success and the explicit failure/fallback behavior

### 2.11. Build and shipping

Keep the build boring:

- `dune build`
- `dune runtest`
- native binary path
- minimal toolchain assumptions

Run before push:

```
dune build
dune runtest
```

If the toolchain is unavailable in the current environment, do a class-wide audit of the touched failure pattern instead of relying on CI to discover all sibling cases.

---

## 3. Rules

3.1. **Type safety rule**
Disambiguate record fields at the type definition. If overlapping names are unavoidable, annotate all access sites in one pass.

3.2. **Purity boundary rule**
`lib/` and `protocol/` stay pure. I/O, Unix, Sys, subprocesses, HTTP, and filesystem access stay in `ffi/`, `transport/`, or other adapters.

3.3. **Error rule**
Use `Result` for expected failures. Do not use exceptions for normal domain errors. Never `with _ ->` without explicit classification and justification.

3.4. **Fallback rule**
Every fallback is a policy decision. If a degraded path exists, it must be:
- intentional
- visible
- justified
- testable

3.5. **Resource discovery rule**
Any function that resolves external resources must validate existence or return explicit failure. Do not return speculative paths and let later I/O fail opaquely.

3.6. **Determinism rule**
If output is compared, hashed, or rendered for later use, sort and render canonically.

3.7. **Idempotence and receipt rule**
For restore/update/sync/workflow logic:
- be safe to retry or fail explicitly if not
- leave visible receipts/logs/state on degraded paths
- make partial-state cleanup explicit

3.8. **Preserve cnos boundaries**
Keep these concepts distinct:
- skills choose
- commands dispatch
- orchestrators execute
- extensions provide capability

Do not let one module smear these ideas together.

3.9. **Shell and archive safety**
Use argv-based subprocesses. Reject path traversal on extraction. Do not log secrets. Keep destructive operations explicit.

3.10. **Smell list**
Treat these as review smells:
- `with _ -> []`
- `with _ -> None`
- `with _ -> ""`
- swallowing `Sys_error`
- fabricated empty state
- compatibility aliases with no reason or expiry
- returning unvalidated external paths
- partial stdlib functions in core/runtime logic
- resource lookup deferred to opaque later failure

3.11. **Build-before-push rule**
Run `dune build` and `dune runtest` before push. If not possible locally, audit the whole failure class, not just the edited line.

---

## 4. Verify

4.1. **File-level check**
- Is the type model explicit?
- Is the purity boundary clear?
- Is fallback policy named?

4.2. **Boundary check**
- Is I/O at the edge?
- Are parsing/validation/rendering kept pure?

4.3. **Error/fallback check**
- Any swallowed errors?
- Any fabricated empty state?
- Any degraded path that is invisible?

4.4. **Determinism/receipt check**
- Are sorted/canonical outputs enforced where needed?
- Does retry/partial failure leave visible evidence?

4.5. **Build/test check**
- `dune build`
- `dune runtest`
- degraded-path tests present

---

## 5. Kata

### 5.1. Kata A — package restore path

**Scenario:** Implement or modify first-party package restore from an artifact index and tarball.

**Task:** Write OCaml code that:
- resolves name + version
- downloads artifact + checksum
- verifies SHA-256
- extracts into vendor
- validates `cn.package.json`
- fails visibly on malformed/missing cases

**Governing skills:** ocaml, design, test

**Inputs:**
- valid package index fixture
- malformed index fixture
- missing package fixture
- checksum mismatch fixture
- tarball extraction path fixture

**Expected artifacts:**
- concrete package types
- pure parse/validate modules
- explicit fallback policy
- no fabricated empty state
- expect tests for success and degraded cases

**Verification:**
- no `with _ ->` swallowing parse/IO failure
- restore path is explicit and retry-safe
- `dune build` and `dune runtest` pass

**Common failures:**
- parse logic inside fetch/process code
- fallback to `[]` / `None`
- unchecked path construction
- missing degraded-path tests

**Reflection:**
- Did the compiler reduce ambiguity early?
- Did every fallback become visible and justified?

### 5.2. Kata B — runtime contract extraction

**Scenario:** Move runtime-contract model/rendering out of a large command-side module into a purer library surface.

**Task:** Extract:
- record types
- validation helpers
- rendering helpers

into a reusable pure module, leaving only runtime wiring in the command layer.

**Governing skills:** ocaml, design, write

**Inputs:**
- one mixed command-side contract builder
- one expected JSON contract output
- one expected markdown/text output

**Expected artifacts:**
- pure contract types/rendering
- explicit wiring in the command layer
- tests for JSON and text rendering
- no I/O inside the contract model

**Verification:**
- no system-truth logic stranded in the command module
- no cyclic dependencies introduced
- output remains stable under tests

**Common failures:**
- moving side effects into the pure module
- leaving rendering behavior untested
- preserving hidden fallback behavior

**Reflection:**
- Did the extraction reduce command-layer ownership of system truth?
- Did the new module make degraded behavior easier to audit?
