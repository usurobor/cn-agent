# Plan: #174 — Orchestrator IR Runtime

## Gap

cnos has no execution engine for mechanical workflows. The runtime contract (#173) exposes an orchestrator registry, but nothing can load, validate, or execute the `cn.orchestrator.v1` IR.

## Design reference

ORCHESTRATORS.md §7–8 (orchestrator model + IR schema).

## Staging

### Stage A: Package content class + IR schema + validation (AC1, AC2, AC6)

**Goal:** `orchestrators/` becomes a recognized content class, the IR schema is defined, and `cn doctor` validates orchestrator manifests.

**Files:**
- `src/cmd/cn_build.ml` — add `orchestrators` to `source_decl`, `parse_sources`, `copy_source` loop, and tarball builder
- `src/cmd/cn_orchestrator.ml` (new) — IR types, JSON parser, schema validator
- `src/cmd/cn_doctor.ml` — add orchestrator validation (missing steps, invalid refs, permission gaps)
- `packages/cnos.core/cn.package.json` — add `"orchestrators": []` to sources (empty initially)
- `docs/alpha/package-system/PACKAGE-SYSTEM.md` — add orchestrators to content-class table (per #174 scope note)

**Tests (write first):**
- IR parser: valid manifest, missing required fields, unknown step kinds, invalid step refs
- Schema validator: permission validation, step ref integrity, duplicate IDs
- Doctor: orchestrator-specific checks
- Build: orchestrators/ copied to package dir

**AC mapping:** AC1 (content class in manifest + build), AC2 (schema defined + validated), AC6 (doctor validates)

### Stage B: Step execution engine (AC3, AC4, AC5)

**Goal:** Runtime can load and execute orchestrator IR step by step.

**Files:**
- `src/cmd/cn_orchestrator.ml` — add execution engine: step dispatch, binding environment, receipt emission
- `src/cmd/cn_executor.ml` — integration point: orchestrator execution as a new dispatch path (separate from N-pass bind loop)

**Step kinds to implement (v1):**
- `op` — dispatch to typed ops via `Cn_shell.execute`
- `llm` — call agent (requires prompt + input binding)
- `if` — evaluate condition on bound value, branch to step
- `match` — pattern match on bound value
- `return` — emit final value
- `fail` — abort with error

**Deferred to v2:**
- `parallel` — concurrent step execution (requires async model cnos doesn't have)

**Key constraints:**
- Determinism rule: control flow is deterministic; non-determinism only at `llm` and `op` boundaries
- Permission manifest enforced before execution — if an op is not in `permissions.ops`, execution fails before dispatch
- Receipts: each step emits a trace event (`orchestrator.step.start`, `orchestrator.step.complete`, `orchestrator.step.error`) to the ulog

**Tests (write first):**
- Step execution: each step kind individually
- Binding environment: values flow between steps
- Permission enforcement: rejected op, rejected llm
- Receipt emission: trace events for each step
- Error propagation: `fail` step, op error, llm error

**AC mapping:** AC3 (step execution), AC4 (permission enforcement), AC5 (execution trace)

### Stage C: Ship daily-review orchestrator (AC7)

**Goal:** A real orchestrator is authored, included in a package, and executable.

**Files:**
- `src/agent/orchestrators/daily-review/orchestrator.json` — the IR manifest
- `packages/cnos.core/cn.package.json` — add `"orchestrators": ["daily-review"]` to sources
- `src/cmd/cn_command.ml` or new `commands/cn-daily` — command that triggers the orchestrator

**The daily-review orchestrator should:**
1. Read today's daily reflection thread
2. Summarize via LLM
3. Write output to an adhoc thread
4. Return the artifact path

This is the §8.1 example from ORCHESTRATORS.md made real.

**Tests:**
- Orchestrator JSON validates against schema
- Build includes orchestrator in package
- Integration: command triggers orchestrator (may be a cram test)

**AC mapping:** AC7

## Conventions

- **No bare `with _ ->` catches** — v3.32.0 convention. Log with context or use specific exception patterns.
- **Reuse existing infrastructure:**
  - `Cn_json` for JSON parsing
  - `Cn_ffi.Fs` / `Cn_ffi.Path` for file ops
  - `Cn_shell.execute` for op dispatch
  - `Cn_ulog` for trace events
  - `build_orchestrator_registry` in `cn_runtime_contract.ml` already reads `sources.orchestrators` — coordinate, don't duplicate
- **Test before code** — write expect tests for each module before implementation
- **OCaml 4.14** — no features from later versions
- **`cn_orchestrator.ml`** as the single new module — IR types, parser, validator, executor all in one file (split later if it grows)

## Active skills

- `eng/ocaml` — OCaml conventions
- `eng/test` — test-first, expect-test format

## Artifact order

1. Tests for Stage A (IR types, parser, validator, doctor)
2. Stage A implementation
3. Tests for Stage B (execution, permissions, receipts)
4. Stage B implementation
5. Stage C (orchestrator JSON + package + command)
6. Self-coherence report

## Risk

- **`parallel` step kind:** Deferred. cnos has no async model. If needed later, it's a separate issue.
- **`llm` step kind:** Requires calling the agent from within an orchestrator. The mechanism for this (prompt + context injection) needs design during Stage B. May simplify to "call cn with a prompt template" initially.
- **Permission enforcement granularity:** v1 enforces at the op-name level. Per-argument enforcement is future scope.
