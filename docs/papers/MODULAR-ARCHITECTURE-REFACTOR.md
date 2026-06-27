# Modular Architecture Refactor

_Decision record: why and how the cnos codebase catches up to the modular architecture the docs already describe._

**Status:** Proposed  
**Date:** 2026-04-08  
**Context:** v3.36.0 shipped orchestrator IR runtime (#174). Package system has 8 content classes. Runtime contract exposes activation index + command/orchestrator registries. Docs describe a highly modular system; code layout is still monolithic.

---

## 1. The Gap

The docs already describe a modular system:

- Runtime contract includes cognition-side `activation_index` and body-side `commands` and `orchestrators`
- Package system treats commands and orchestrators as content classes alongside skills, extensions, and templates
- Runtime extensions are explicitly a separate system for capability providers

But the codebase still looks monolithic:

- `src/cmd/` contains a broad mix of concerns: activation, agent loop, build, deps, commands, context, mail, maintenance, orchestrator, runtime contract, output, system
- `src/lib/` is tiny (cn_json, cn_lib, cn_sha256) — domain logic lives in the command layer instead of a reusable core
- `cn_lib.ml` hardcodes the visible command/help surface in one block — the opposite of package-driven discovery

## 2. Architectural Decisions

### 2.1 Super-lean core

The core shrinks to pure contracts:

- Manifests and schemas
- Registries
- Lockfile/index resolution
- Activation rules
- Orchestrator IR + validator
- Runtime contract types/rendering
- Protocol envelope/receipt types

Everything else — network, git, fs, LLM calls, package fetch, command execution — sits above as adapters/services.

### 2.2 Package-driven intelligence with built-in bootstrap

Almost all intelligence comes from installed packages (skills, activation, commands, orchestrators, templates). But a minimal bootstrap kernel stays built-in:

- Package restore/install
- Package index / lockfile validation
- Runtime contract emission
- help / doctor / status / init / setup / update / deps / build
- Workflow engine / extension host

### 2.3 Commands follow the standard content-class pipeline

Current asymmetry: most content classes are copied from `src/agent/…` by `cn build`, but commands are authored directly under `packages/<name>/commands/`. This breaks the "one source → build → install" shape.

Fix: `src/agent/commands/<id>/…` → `cn build` → `packages/<name>/commands/<id>/…` → `cn deps restore` → `.cn/vendor/…`

This makes commands follow the same build/install model as skills, orchestrators, extensions, and templates.

### 2.4 Most commands become packages

Built-in command kernel (makes the package system possible):

- help, init, setup, deps, build, doctor, status, update

Everything else moves to packages over time:

- daily, weekly, save, commit, push, gtd, adhoc, future workflow commands

### 2.5 A2A: transport kernel in core, workflows in packages

Low-level A2A substrate stays in core:

- Peer config, sync/fetch/push, message envelope validation, routing invariants, receipts/rejection/retry

Higher-level A2A behaviours become packages/orchestrators:

- Review request flows, issue handoff, release coordination, daily/weekly peer rituals, inbox triage styles

### 2.6 One package system with roles, not two systems

One artifact format, one index, one restore/install pipeline, one manifest schema. But add a role/kind notion for policy and defaults:

- **base** — bootstrap/runtime-critical (cnos.core)
- **cognition** — doctrine/mindsets/skills/templates
- **command** — operator command packages
- **workflow** — orchestrator-heavy packages
- **capability** — runtime extension providers
- **bundle** — packages that aggregate several of the above

### 2.7 Four distinct content concerns

- **Skills** → cognition / judgment
- **Commands** → operator dispatch
- **Orchestrators** → mechanical workflows
- **Extensions** → capability providers

Do not merge these. Each has a different runtime role.

## 3. Phased Refactoring Plan

### Phase 1 — Define the real core

Create a proper core layer. Move pure logic into it:

- Package manifest schema + validation
- Package index + lockfile schema
- Activation-index schema + evaluator
- Command registry schema
- Orchestrator IR schema + validator
- Runtime-extension manifest schema
- Runtime contract types + renderers
- Protocol envelope / receipt / trace types

Key moves from `src/cmd/`:

- `cn_deps.ml` → manifest/index/lockfile types into core
- `cn_runtime_contract.ml` → contract types/rendering into core
- `cn_activation.ml` → activation rule evaluation into core
- `cn_workflow.ml` → workflow IR + validator into core

Goal: `src/cmd/` stops owning the system model; it only orchestrates host operations.

### Phase 2 — Split host adapters from runtime services

Host/adapters layer for impure operations: git, filesystem, HTTP/download, env/secrets, process exec, archive extract, LLM provider client.

Runtime/services layer that composes them: package manager, command registry builder, activation loader, workflow runner, runtime contract builder, maintenance loop, A2A transport coordinator.

### Phase 3 — Make packages the primary source of intelligence

At install/wake time, build two registries from installed packages:

**activation_index** — declarative activation rules per skill (event / work-shape / path / level / risk / token hints, requires / excludes).

**command_registry** — exact name, entrypoint, source, summary, trust/policy flags.

### Phase 4 — Make orchestrators first-class runtime artifacts

CTB / F# CE / YAML as source forms → compiled orchestrator IR as runtime artifact → runtime executes IR under policy. The orchestrator runtime belongs in core/runtime services.

### Phase 5 — Keep runtime extensions separate

Extensions are capability providers, not packages or skills. Final separation: skills (cognition), commands (dispatch), orchestrators (workflows), extensions (capabilities).

## 4. Target Structure

```
src/
  core/                    # pure contracts, schemas, validators
    package/
    activation/
    commands/
    orchestrators/
    extensions/
    contract/
    protocol/
  host/                    # impure adapters
    fs/
    git/
    http/
    env/
    process/
    archive/
    llm/
  runtime/                 # composed services
    package_manager/
    activation_loader/
    command_dispatch/
    workflow_runner/
    contract_builder/
    maintenance/
    transport/
  builtins/                # tiny bootstrap command set
    help/
    init/
    setup/
    deps/
    build/
    doctor/
    status/
    update/
  agent/                   # package-authored cognition assets
    doctrine/
    mindsets/
    skills/
    commands/
    orchestrators/
    templates/
```

Not a one-shot rename. Move logic gradually toward this structure.

## 5. Execution Order

### Immediate execution slice (KISS/YAGNI)

Full design is the north star. Execute three moves first:

1. **Command pipeline symmetry** — commands through `src/agent/` → `cn build` → `packages/`. First migrations: daily, weekly, save. Built-in shrinks to 8.
2. **Pure-model gravity into `src/lib/`** — no new directories. Pure types sink down, impure stays in `src/cmd/`.
3. **Retire beta package doc (#180)** — redirect stub to alpha spec.

Skip for now: package roles, richer activation schema, `src/host/`/`src/runtime/` split, gh-pages (#181), full directory tree.

### Full plan (north star)

1. Command pipeline symmetry
2. `src/lib/` extraction
3. Doc coherence (#180)
4. Package source (#181) — when scale justifies
5. `src/core/` / `src/host/` / `src/runtime/` — when `src/lib/` gets crowded
6. Full design phases incrementally

## 6. Related

- #182 — Core refactor umbrella issue (L7)
- #180 — Doc incoherence: beta claims Git-native transport
- #181 — gh-pages package source
- #174 — Orchestrator IR runtime (shipped v3.36.0)
- #167 — Package artifact distribution (shipped v3.34.0)
- #162 — Modular CLI commands
- #96 — Docs taxonomy alignment
