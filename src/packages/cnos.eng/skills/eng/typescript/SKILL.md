---
name: typescript
description: Write TypeScript that keeps runtime behavior, type surfaces, and module boundaries aligned. Use when building CLIs, services, tooling, registries, package flows, or browser/node apps in TypeScript.
artifact_class: skill
kata_surface: embedded
governing_question: How do we write TypeScript that is explicit at runtime, strict at compile time, and hard to lie to?
triggers:
  - implementing cnos-adjacent tools or services in TypeScript
  - building package/index/registry tooling in TS
  - reviewing TypeScript code for runtime/type drift
  - debugging async, module, or schema issues in TS
  - migrating JavaScript into TypeScript
---

# TypeScript

## Core Principle

**Write TypeScript so the runtime truth and the type story stay aligned.**

Use schema-backed validation at external boundaries, discriminated unions for state, and branded primitives where accidental interchange would be costly. Prefer strict compiler settings and explicit runtime checks over `any`, unchecked assertions, and wishful typing.

## Algorithm

1. **Define** — name the runtime model, the module boundary, the async/error policy, and the validation boundary.
2. **Unfold** — design types, modules, runtime validation, branded primitives, async flow, and tests together.
3. **Rules** — keep the code strict, schema-backed at boundaries, exhaustively handled, and hostile to type/runtime drift.
4. **Verify** — compile, lint, test, and check the same boundary cases across the whole surface.

---

## 1. Define

### 1.0. Classify the artifact first

Use this skill for TypeScript code that belongs to:

- CLI/tools
- services
- registries/manifests
- package/index handling
- command dispatch
- orchestrator tooling
- browser/node apps where runtime/type drift matters

Do **not** use this skill for:

- prose-only documentation
- framework style wars with no runtime consequence
- "add types later" JavaScript scaffolding

### 1.1. Identify the parts

A coherent TypeScript subsystem has:

- **Runtime model** — what actually exists at runtime
- **Type model** — how that reality is represented statically
- **Boundary validation** — where `unknown` input becomes trusted data
- **Branded primitives** — distinct identities for easily confused scalar values
- **Module surface** — what a module exports and what it hides
- **Async policy** — how promises, cancellation, and retries are handled
- **Error policy** — how failure is represented and propagated
- **Build/lint/test contract** — what proves the code is safe to ship
  - ❌ Start with types disconnected from runtime parsing
  - ✅ Define the runtime payloads and the type model together

### 1.2. Articulate how they fit

External values cross boundaries as `unknown`. Validation makes them trusted. Types model trusted shapes. Branded primitives keep distinct runtime meanings from collapsing into one wide scalar type. Modules expose only what callers need. Async flow is explicit. Errors are visible and policy-driven.

- ❌ Parse JSON straight into `any`
- ✅ Parse as `unknown`, validate, then narrow to the domain type

### 1.3. Name the failure mode

TypeScript fails through **type/runtime drift**. Typical signs:

- `any` leakage
- unchecked `as X` at trust boundaries
- promises started and ignored
- unions never exhaustively checked
- optional fields used as if always present
- raw values standing in for many incompatible identities
- type-only imports mixed carelessly with runtime imports
- module surface confusion between ESM/CJS/runtime shape
- invisible degraded paths
- mutation of shared state that changes behavior silently
  - ❌ `const cfg = JSON.parse(text) as Config`
  - ✅ `const raw: unknown = JSON.parse(text); const cfg = parseConfig(raw)`

---

## 2. Unfold

### 2.1. Compiler posture

Start strict. Baseline:

- `strict: true`

Recommended where it fits:

- `noUncheckedIndexedAccess`
- `exactOptionalPropertyTypes`
- `noImplicitOverride`
- `useUnknownInCatchVariables`
  - ❌ Start loose and promise to tighten later
  - ✅ Start strict and carve out only justified exceptions

### 2.2. Runtime boundary first

Every external boundary begins as `unknown`. Typical boundaries:

- JSON parse
- HTTP request/response
- env/config
- local storage / database reads
- message bus / queue payloads
- plugin/extension manifests

Prefer a runtime schema library at external trust boundaries. Good examples:

- Zod
- Valibot
- ArkType
- similar runtime-first schema systems

The rule is not "use Zod." The rule is: use explicit runtime schemas where untrusted data enters the system.

- ❌ `function load(x: any): Settings`
- ❌ `const manifest = raw as PackageManifest`
- ✅ `function parseSettings(x: unknown): Settings`

### 2.3. Type design

Prefer:

- concrete object types
- discriminated unions
- `readonly` where mutation is not part of the design
- `keyof`, indexed access, and generics when they reduce duplication without hiding the model

Use discriminated unions for state machines and command/result surfaces.

```typescript
type RestoreResult =
  | { kind: "ok"; packageName: string; version: string }
  | { kind: "missing"; packageName: string }
  | { kind: "checksum-mismatch"; packageName: string; expected: string; actual: string };
```

- ❌ One wide object with many optional fields and no discriminant
- ✅ One discriminated union with one obvious runtime case at a time

### 2.4. Branded primitives

Use branded types for domain primitives that are:

- easily confused
- semantically distinct
- costly to interchange accidentally

Good candidates:

- paths
- checksums
- package names
- versions
- command names
- IDs

```typescript
type Brand<T, B extends string> = T & { readonly __brand: B };

type AbsolutePath = Brand<string, "AbsolutePath">;
type Checksum = Brand<string, "Checksum">;
type PackageName = Brand<string, "PackageName">;
```

- ❌ `function verify(path: string, checksum: string)`
- ✅ `function verify(path: AbsolutePath, checksum: Checksum)`

Do not brand everything. Brand the primitives where accidental substitution would create real bugs.

### 2.5. `unknown`, guards, and assertions

Use:

- `unknown` at boundaries
- schema-backed parsing for external input
- local type predicates / assertion functions where they genuinely clarify internal control flow
- `satisfies` to check object shapes without widening away useful literal information

Do not use handwritten type guards as the primary trust boundary for external data.

- ❌ `function isConfig(x: unknown): x is Config` as the only protection on arbitrary JSON
- ✅ `ConfigSchema.parse(x)` at the boundary, then local narrowing helpers inside the trusted domain

Avoid:

- `any`
- unchecked `as`
- double assertions (`as unknown as X`)

### 2.6. Errors and fallback

Make the error policy explicit. Use:

- thrown `Error` objects for exceptional failure
- discriminated result types for expected domain failure
- `unknown` in catch blocks
- explicit wrapping/normalization at boundaries

```typescript
try {
  await execute();
} catch (err: unknown) {
  throw new Error(`execute failed: ${err instanceof Error ? err.message : String(err)}`);
}
```

- ❌ `throw "bad config"`
- ❌ `catch (err: any) { ... }`
- ✅ Throw `Error` instances and inspect caught values explicitly

Every fallback must be:

- intentional
- visible
- testable

### 2.7. Modules and imports

Keep module boundaries explicit. Prefer:

- small modules
- named exports by default
- `import type` for type-only imports
- one clear ESM/CJS strategy per package
  - ❌ Mix runtime and type-only imports casually
  - ✅ `import type { PackageLock } from "./types.js"`

### 2.8. Async and promises

Every promise must be:

- awaited
- returned
- or intentionally handled

No floating promises. No async callbacks smuggled into sync-shaped APIs by accident.

- ❌ `doRestore()` and ignore the returned promise
- ✅ `await doRestore()`

If concurrency exists, own it explicitly:

- `Promise.all`
- bounded task groups
- retries with clear policy
- `AbortSignal` where cancellation matters

### 2.9. Exhaustiveness

If a union matters, handle it exhaustively.

```typescript
function assertNever(x: never): never {
  throw new Error(`unreachable case: ${String(x)}`);
}

switch (result.kind) {
  case "ok": return ...;
  case "missing": return ...;
  case "checksum-mismatch": return ...;
  default: return assertNever(result);
}
```

- ❌ Default branch that silently absorbs future cases
- ✅ Exhaustive switch with `assertNever`

### 2.10. Schema and compatibility

Treat manifests, contracts, and registries as public contracts.

When handling:

- package manifests
- lockfiles
- indexes
- runtime contracts
- orchestrator IR
- extension manifests

prefer:

- explicit version fields
- additive evolution
- strict validation of required fields
- clear behavior for unknown fields
  - ❌ Silently coerce malformed manifest data into defaults
  - ✅ Reject malformed input with repairable errors

### 2.11. Determinism and reproducibility

For outputs that are:

- hashed
- compared in tests
- written to reports
- stored for later diffing

ensure:

- stable ordering
- canonical serialization
- no dependence on object key enumeration accidents where that matters
- no hidden timestamps unless intentionally included
  - ❌ Render objects in arbitrary order where stability matters
  - ✅ Sort before rendering

### 2.12. Idempotence, retry, and receipts

Assume package restore, update, command execution, and workflow steps may be retried.

For state-changing operations, define:

- what happens on interruption
- whether the operation is idempotent
- what partial state is left behind
- how cleanup works
- what visible receipt/log/status tells the caller what happened
  - ❌ Write directly to the final target and hope the process completes
  - ✅ Write to temp, validate, then atomically move into place
  - ❌ Degraded path only visible in transient stderr
  - ✅ Degraded path visible in logs, receipts, or structured state

### 2.13. Side-effect boundaries

Keep parsing/validation/planning separate from:

- filesystem
- HTTP
- subprocesses
- archive extraction
- env access
  - ❌ Parse package manifests inside the fetcher
  - ✅ Fetch bytes in one module; parse and validate in the domain module

### 2.14. Node/browser reality

Types do not exist at runtime. Module systems, globals, file APIs, and network APIs do. Be explicit about:

- environment assumptions
- polyfills/shims if any
- Node vs browser boundaries
- runtime feature detection where necessary
  - ❌ Write code that type-checks but assumes a runtime that isn't there
  - ✅ Make the runtime target obvious in the package and module surface

### 2.15. Testing

Use:

- `tsc --noEmit`
- linter
- unit tests
- boundary tests
- unhappy-path tests

Test:

- malformed input
- missing fields
- unknown union members
- async failure
- retry/degraded paths
- deterministic output
  - ❌ Test only the happy path because "types prove the rest"
  - ✅ Test the runtime boundary and degraded behavior explicitly

---

## 3. Rules

### 3.1. Start strict

Use `strict` mode by default.

### 3.2. External input starts as `unknown`

Do not trust JSON, HTTP, env, storage, or plugin data without validation.

### 3.3. Prefer schema-backed trust boundaries

At external I/O boundaries, prefer runtime schema parsing over handwritten guards.

### 3.4. Avoid `any`

Use `unknown`, unions, generics, and validation instead. If `any` appears, treat it as debt.

### 3.5. Prefer discriminated unions for state and results

They keep runtime cases explicit and exhaustiveness possible.

### 3.6. Use branded primitives where confusion is costly

Brand paths, checksums, IDs, package names, and similar domain primitives when accidental interchange would create real bugs.

### 3.7. Use `satisfies` for object-shape validation

Use `satisfies` when you want to validate shape without erasing more specific inferred structure.

### 3.8. Keep imports honest

Use `import type` for type-only imports and keep runtime imports explicit.

### 3.9. Make error policy explicit

Choose clearly between:

- thrown `Error`
- discriminated result values
- or a deliberate combination by layer

Do not let failure handling emerge accidentally.

### 3.10. No floating or misused promises

Await them, return them, or handle them deliberately.

### 3.11. Make fallback explicit

Fallback is policy, not accident. If a degraded path exists, it must be visible, testable, and reviewable.

### 3.12. Prefer immutable updates by default

Do not silently mutate shared objects when returning a new state/value is clearer.

- ❌ Mutate shared config or registry objects in-place across module boundaries
- ✅ Return new values unless local mutation is clearly contained and justified

### 3.13. Preserve cnos boundaries

Keep these concepts distinct:

- **skills** choose
- **commands** dispatch
- **orchestrators** execute
- **extensions** provide capability

Prefer explicit metadata and registries over hidden inference.

### 3.14. Shell, archive, and path safety

When touching files, subprocesses, or archives:

- avoid shell string interpolation
- validate archive extraction paths
- reject path traversal
- do not log secrets
- make destructive operations explicit

### 3.15. Override precedence must be explicit

If behavior may come from:

- flags
- environment
- config
- runtime discovery
- defaults

define the precedence clearly and test it.

### 3.16. Build and verify before push

Run the equivalent of:

- `tsc --noEmit`
- linter
- test suite

before claiming the code is ready.

---

## 4. Verify

### 4.1. File-level check

- Is the runtime model explicit?
- Are boundaries typed honestly?
- Is validation present where trust begins?

### 4.2. Module check

- Are modules small and specific?
- Are type-only imports separated cleanly?
- Is ESM/CJS/runtime behavior explicit?

### 4.3. Async/error check

- Any floating promises?
- Any swallowed errors?
- Any degraded path that is invisible?

### 4.4. Schema/contract check

- Are manifests/versioned contracts validated?
- Are `unknown` and malformed cases handled explicitly?
- Is output deterministic where it matters?

### 4.5. Build/test check

- compile
- lint
- test
- boundary/degraded-path tests present

---

## 5. Kata

### 5.1. Kata A — package index client

**Scenario:** Implement a TypeScript client that loads a package index, validates package entries, and resolves name + version to a downloadable artifact.

**Task:** Write code that:

- reads JSON as `unknown`
- validates it into a domain type
- resolves a package version
- returns discriminated result types for missing or malformed cases

**Governing skills:** typescript, design, test

**Inputs:**

- valid package index fixture
- malformed index fixture
- missing package fixture
- missing version fixture

**Expected artifacts:**

- domain types using discriminated unions
- explicit schema-backed parse layer
- branded types where domain primitives are easily confused
- no `any`
- tests for valid and invalid cases

**Verification:**

- no unchecked assertions at the trust boundary
- all union results handled exhaustively
- compile, lint, and tests pass

**Common failures:**

- `JSON.parse(...) as Index`
- one wide result type with many optional fields
- missing tests for malformed input
- swallowed errors on missing packages

**Reflection:**

- Did the type story match the runtime truth?
- Did invalid input fail explicitly enough to repair?

### 5.2. Kata B — async command runner

**Scenario:** Implement a command runner that downloads an artifact, verifies it, and writes it atomically.

**Task:** Write code that:

- uses explicit async flow
- never drops promises
- models success/failure as explicit result cases
- keeps I/O separate from parsing/validation

**Governing skills:** typescript, design, test

**Inputs:**

- successful fetch path
- checksum mismatch path
- network failure path
- atomic replace path

**Expected artifacts:**

- no floating promises
- clear error/result policy
- deterministic verification flow
- tests for unhappy paths

**Verification:**

- `no-floating-promises` and `no-misused-promises` would pass
- atomic path is explicit
- checksum mismatch is visible and test-covered

**Common failures:**

- fire-and-forget promise
- mixing validation inside transport
- broad catch that erases failure detail
- optional fields used as if required

**Reflection:**

- Did the async flow stay visible and auditable?
- Did retries/failures remain policy-driven rather than accidental?
