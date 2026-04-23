# Core Refactor for a Package-Driven cnos Runtime

**Issue:** #182  
**Version:** 0.1.0  
**Mode:** MCA  
**Active Skills:** cdd/design, eng/evolve, eng/process-economics  
**Level:** L7

---

## Authority Relationships

After this refactor:

- Package-system docs own content classes and install semantics
- Runtime contract owns what is visible at wake
- Orchestrator spec owns workflow IR
- Runtime-extensions own capability-provider rules
- Built-in CLI owns only bootstrap commands
- Package/command/orchestrator registries are derived runtime structures, not hand-maintained truth

---

## Proposal

Refactor cnos around one package-driven runtime substrate with four distinct runtime surfaces:

1. **Skills** — judgment surfaces
2. **Commands** — exact-dispatch operator actions
3. **Orchestrators** — mechanical workflows
4. **Extensions** — runtime capability providers

And one architectural split:

- **core** — pure schemas, validators, registries, protocol types
- **host** — IO/effects adapters
- **runtime** — composed services
- **builtins** — tiny bootstrap command set
- **agent** — package-authored assets for nearly everything else

---

## Problem

cnos has arrived at a much more modular architecture in theory than in code. The current documentation now describes:

- A runtime contract with:
  - `activation_index` under cognition
  - `commands` and `orchestrators` under body
- A package system with explicit content classes, including:
  - skills, commands, orchestrators, extensions, templates, doctrine, mindsets
- A runtime-extensions system that remains distinct from package-distributed logic

But the internal code structure still reflects an older, more monolithic system:

- `src/cmd/` mixes: manifests, package logic, runtime contract rendering, activation logic, command dispatch, A2A mail transport, maintenance loop, workflow/orchestration
- Many system-model concerns live in command handlers instead of a reusable core
- Built-in commands still own much of the user-facing surface that should increasingly come from packages
- Commands are not yet authored through the same source → build → install pipeline as the other content classes

The result is a structural incoherence:

- The docs describe a package-driven modular runtime
- The code still behaves like a command-centric monolith with package add-ons

---

## Constraints

- Keep one package substrate. Do not create separate package systems for "agent" vs "cnos."
- Preserve a minimal built-in bootstrap surface: help, init, setup, deps, build, doctor, status, update
- Do not collapse runtime extensions into commands or orchestrators.
- Do not move low-level A2A transport into packages.
- Preserve current cnos semantics during migration.
- Prefer staged refactor over flag-day rewrite.

## Challenged Assumption

This design challenges the assumption that the command layer is the natural place to hold system truth.

That assumption no longer serves cnos. The code should no longer treat package install, activation, command discovery, orchestrator execution, runtime contract rendering, and protocol envelopes as command-level details. They are now core/runtime concerns that command handlers should consume, not own.

---

## 1. One Package Substrate, Not Two Package Systems

cnos should use one package system for both:

- "Agent intelligence" (doctrine, mindsets, skills, templates)
- "cnos behaviors" (commands, orchestrators, extensions)

Do not create one package type for agents and another for cnos itself. Instead, introduce package roles on top of one substrate.

### Proposed Package Roles

| Role | Description | Example |
|------|-------------|---------|
| base | Bootstrap/runtime-critical | `cnos.core` |
| cognition | Doctrine, mindsets, skills, templates | `cnos.cognition` |
| command | Operator command packages | `cnos.ops` |
| workflow | Orchestrator-heavy packages | `cnos.workflows` |
| capability | Runtime extension providers | `cnos.net.http` |
| bundle | Packages that aggregate several of the above | `cnos.full` |

The artifact format, index, lockfile, restore flow, and install path stay the same. Only policy/default behavior differs by role.

---

## 2. Commands Become a First-Class Package/Runtime Surface

Commands should share the package substrate but remain distinct from skills.

### Rule

Commands are:

- Exact-dispatch operator actions
- Not activated by keyword
- Not judgment surfaces
- Not runtime extensions

### Discovery Precedence

1. Built-in commands
2. Repo-local commands
3. Vendored package commands

No PATH-wide discovery in v1.

### Package Manifest Command Metadata

```json
"commands": {
  "daily": {
    "entrypoint": "commands/daily/cn-daily",
    "summary": "Create or show the daily reflection thread",
    "needs_hub": true,
    "dangerous": false
  }
}
```

### Consequence

Most non-bootstrap commands can leave core over time.

**Likely built-in forever:** help, install, deps, build, doctor, status, update

`init` and `setup` are retired — agent creation moves to `cn agent create` (from the `cnos.agent` package). The binary is a package manager + command dispatcher; the agent lifecycle is a package.

**Likely package/repo-local:** daily, weekly, save, commit, push, gtd, adhoc, release wrappers, project-specific workflows

---

## 3. Skills Become Fully Package-Driven and Activation-Indexed

Skills should not be "loaded by trigger words" alone. The robust model is a declarative activation index.

### Skill Activation Inputs

- Event kind
- Work shape
- Path hints
- Level hints
- Dominant risk
- Token hints
- Requires / excludes

Keywords remain a weak signal only.

### Runtime Contract Addition

Under cognition, the runtime contract exposes `activation_index`:

```json
"activation_index": {
  "skills": [
    {
      "id": "cdd/review",
      "package": "cnos.eng",
      "summary": "Review a change against its declared contract",
      "activation": {
        "events": ["cdd:step=8", "inbox:review_request"],
        "tokens": ["review", "pr", "diff"],
        "paths": ["docs/**", "src/**"],
        "work_shapes": ["review", "docs/process"],
        "levels": ["L5", "L6", "L7"],
        "requires": [],
        "excludes": []
      }
    }
  ]
}
```

This makes the "agent as package-loaded intelligence" model real.

---

## 4. Orchestrators Become a Distinct Runtime Surface

An orchestrator is not a skill. A skill teaches judgment. An orchestrator executes a mechanical workflow.

### Rule

Orchestrators: sequence, bind, branch, retry, call the LLM, call typed runtime ops, emit artifacts/receipts.

They do not decide what is wise. They do what has already been selected.

### Runtime Contract Addition

Under body, expose `orchestrators`:

```json
"orchestrators": [
  {
    "name": "daily-review",
    "source": "package",
    "package": "cnos.core",
    "trigger_kinds": ["command", "schedule"]
  }
]
```

### Authoring vs Runtime Form

Humans may author orchestrators in source forms such as CTB, F# computation expressions, or YAML workflows. The runtime validates and executes a compiled orchestrator IR.

---

## 5. Runtime Extensions Remain Separate

Do not merge commands or orchestrators into the runtime-extensions system.

Extensions remain: typed capability providers, runtime/body-plane components, policy-governed. Not skills, not commands, not orchestrators.

This boundary should stay sharp.

---

## 6. A2A Communication Split

**Keep in core:** peer registry, message envelope validation, routing invariants, receipts/rejection/retry/dedupe, sync/fetch/push substrate.

**Package/orchestrator surfaces may own:** review request flows, issue handoffs, release coordination, daily/weekly peer rituals.

This keeps the transport substrate trustworthy while allowing communication behavior to evolve.

---

## 7. Commands Use the Same Source → Build → Install Pipeline

Current package docs still make commands an exception to the normal content flow. That should be removed.

### Desired Shape

```
src/agent/commands/<id>/...
  ↓ cn build
packages/<name>/commands/<id>/...
  ↓ cn deps restore
.cn/vendor/packages/<name>@<version>/commands/<id>/...
  ↓ command discovery
cn <name>
```

That makes commands behave like the other content classes: authored in `src/agent/`, built into package artifacts, restored into vendor, discovered at runtime.

---

## 8. Internal Code Structure

### Target Structure

```
src/
  core/                    # pure schemas, validators, registries, protocol types
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
  builtins/                # minimal bootstrap command set
    help/
    install/
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

### Key Rule

`src/cmd/` should stop being where system truth lives. It becomes: CLI entry, bootstrap dispatch, glue into runtime services. The system model moves into `core/` and `runtime/`.

---

## Impact Graph

### Downstream Consumers

- cn help, cn doctor, cn status
- cn deps restore
- Package command discovery
- Skill activation at wake
- Orchestrator execution
- Runtime contract emission
- A2A sync / inbox / outbox flows
- Package build/install pipeline
- Docs and package manifests

### Upstream Producers

- Package manifests
- Package index / lockfile
- Repo-local commands
- Skill frontmatter / activation declarations
- Orchestrator manifests
- Runtime extension manifests
- Hub placement and runtime state

### Copies and Embeddings

These concepts currently appear in multiple places and should be normalized:

- Package content classes
- Command discovery rules
- Activation metadata
- Runtime contract fields
- Protocol authority around transport vs workflow
- Built-in help text

---

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Keep current `src/cmd/` monolith, add package hooks | Lowest immediate churn | Preserves structural incoherence; docs and code keep drifting | Rejected |
| Separate package systems for "agent" and "cnos" | Clear at first glance | Two indexes, two restore paths, duplicated complexity | Rejected |
| Make commands just a type of skill | Fewer concepts | Breaks exact dispatch, muddies judgment vs execution | Rejected |
| Make A2A entirely package-driven | Highly flexible | Weakens transport substrate and protocol kernel | Rejected |
| **One package substrate + distinct runtime registries + super-lean core** | Strong separation, scalable, coherent with current docs | Requires staged refactor | **Chosen** |

---

## Process Cost / Automation Boundary

**Automate:**

- Package manifest validation
- Package role validation
- Activation-index building
- Command registry building
- Orchestrator registry building
- Runtime contract emission
- Doctor validation for duplicates / broken entrypoints / missing registries

**Keep human / design judgment:**

- What package role is appropriate
- What should remain built-in
- Whether a workflow belongs in a command or an orchestrator
- Whether an A2A behavior is substrate or package-level ritual

---

## Leverage

This refactor makes future growth cheaper:

- New skills become package additions, not core edits
- New commands become package additions, not core edits
- Orchestrator execution becomes explicit and testable
- Activation becomes inspectable instead of hidden prompt behavior
- Built-in command count can shrink over time
- Transport and capability substrate stay stable while behavior grows around them

### Negative Leverage

This adds:

- Registry/schema complexity
- More explicit runtime surfaces
- More code movement before visible new features
- Transitional duplication while `src/cmd/` logic is extracted
- Higher bar on manifest quality and runtime validation

---

## Non-Goals

This refactor does not:

- Redesign CTB itself
- Define the full orchestrator IR (already done in #174)
- Change runtime extension semantics
- Solve package ecosystem governance beyond first-party packages
- Remove all built-ins immediately
- Move low-level transport out of core

---

## File Changes

### Create

- `docs/alpha/agent-runtime/CORE-REFACTOR.md` (this design)
- `src/core/` scaffolding
- `src/runtime/` scaffolding
- `src/builtins/` scaffolding
- `src/agent/commands/`
- `src/agent/orchestrators/` (already exists as of v3.36.0)

### Edit

- `docs/alpha/agent-runtime/RUNTIME-CONTRACT-v2.md`
- `docs/alpha/package-system/PACKAGE-SYSTEM.md`
- `docs/alpha/package-system/PACKAGE-ARTIFACTS.md`
- `docs/alpha/runtime-extensions/RUNTIME-EXTENSIONS.md`
- `docs/beta/architecture/PACKAGE-SYSTEM.md` (retired — see alpha spec)
- `src/cmd/cn_runtime_contract.ml`
- Package build/install code to support commands and orchestrators through the same pipeline

### Migrate

- Move command authoring into `src/agent/commands/`
- Extract pure logic out of `src/cmd/` into `src/core/` and `src/runtime/`

---

## Acceptance Criteria

- [ ] AC1: Runtime contract exposes `activation_index`, `commands`, and `orchestrators`
- [ ] AC2: Package system defines commands and orchestrators as first-class content classes
- [ ] AC3: Commands use the same source → build → install pipeline as the other content classes
- [ ] AC4: Built-in commands are reduced to the bootstrap set only
- [ ] AC5: Package command discovery follows built-in > repo-local > vendored package precedence
- [ ] AC6: Skills activate from a declarative activation index, not trigger words alone
- [ ] AC7: Orchestrators are represented as a distinct runtime surface, not skills
- [ ] AC8: Runtime extensions remain separate from commands/orchestrators
- [ ] AC9: Low-level A2A transport remains core; package/orchestrator A2A workflows are additive
- [ ] AC10: Core/runtime code owns schemas and registries; command handlers stop owning system truth

---

## Immediate Execution Slice (KISS/YAGNI)

The full design is the north star. These three moves execute first:

**Move 1 — Command pipeline symmetry.** Remove the biggest structural exception: commands become ordinary package content, authored in `src/agent/commands/` and flowing through `cn build` like every other content class. Built-in kernel shrinks to: help, install, deps, build, doctor, status, update. Agent creation moves to `cnos.agent` package (`cn agent create`). First migrations: daily, weekly, save. Leave commit and push built-in for one more cycle until the package-command path proves clean. Runtime extensions and low-level A2A transport remain untouched in this slice.

**Move 2 — Pure-model gravity into `src/lib/`.** Do not create `src/core/`. Widen the existing `src/lib/` with pure types and validators extracted from `src/cmd/`. Candidates: package manifest types (from `cn_deps.ml`), runtime contract record types (from `cn_runtime_contract.ml`), workflow IR types/parser/validator (from `cn_workflow.ml`), activation types/evaluator (from `cn_activation.ml`). Discipline: no filesystem / git / process / HTTP / LLM code may move into `src/lib/`. Stderr diagnostic logging (`Printf.eprintf`) is permitted for discovery-time warnings on malformed data — it is unconditional, one-shot, and testable via `ppx_expect` stdout+stderr capture; it is not a process effect in the forbidden sense. Precedent: `Cn_workflow_ir.manifest_orchestrator_ids` (v3.40.0 slice 3) preserved the one `Printf.eprintf` call from its pre-extraction home because the diagnostic was part of the function's contract, not an IO side-effect.

> **Move 2 status (v3.38.0):** First slice shipped — `src/lib/cn_package.ml` now owns the canonical `manifest_dep`, `locked_dep`, `manifest`, `lockfile`, `index_entry`, `package_index` types plus their JSON serialization, `parse_package_index`, `lookup_index`, and `is_first_party`. `src/cmd/cn_deps.ml` re-exports each type via OCaml type-equality syntax (`type t = Cn_package.t = { ... }`) so all existing callers (`cn_runtime_contract`, `cn_system`, `cn_deps_test`, `cn_runtime_contract_test`) compile unchanged. `cn_deps.ml` now contains only IO-side functions: read/write, download, extract, restore, doctor, lockfile_for_manifest. Discipline verified by grep: `cn_package.ml` imports only stdlib + `Cn_json`.
>
> **Move 2 status (v3.39.0):** Second slice shipped — `src/lib/cn_contract.ml` now owns the canonical 11 runtime-contract record types (`package_info`, `override_info`, `zone`, `zone_entry`, `identity`, `extension_contract_info`, `command_entry`, `orchestrator_entry`, `cognition`, `body_contract`, `runtime_contract`) plus `zone_to_string`. The `Cn_activation.activation_entry` transitive dependency was also pulled through — `cn_contract.ml` owns `activation_entry` as canonical, and `cn_activation.ml` re-exports it via type-equality so existing callers of `Cn_activation.activation_entry` (test files + `cn_runtime_contract.ml` render/to_json closures) compile unchanged. `src/cmd/cn_runtime_contract.ml` re-exports all 11 types via the same type-equality pattern and retains every IO function untouched (`classify_zones`, `list_md_relative`, `list_skill_overrides`, `extensions_from_registry`, `build_command_registry`, `build_orchestrator_registry`, `gather`, `render_markdown`, `to_json`, `write`). Discipline verified by grep: `cn_contract.ml` imports only stdlib.
>
> **Move 2 status (v3.40.0):** Third slice shipped — `src/lib/cn_workflow_ir.ml` now owns the canonical 6 pure orchestrator IR types (`trigger`, `permissions`, `step` as a 6-variant sum, `orchestrator`, `issue_kind` as a 7-variant sum, `issue`) plus 10 pure functions (`let ( let* )`, `require_string`, `parse_string_list`, `parse_trigger`, `parse_permissions`, `parse_step`, `parse`, `step_id`, `validate`, `manifest_orchestrator_ids`). Largest slice so far (655 LOC source module vs 567 for slice 2 and 130 for slice 1). `src/cmd/cn_workflow.ml` re-exports each type via OCaml type-equality and delegates each pure function via a one-line let-binding; the three IO-transit types (`load_outcome`, `installed`, `outcome`) stay in `cn_workflow.ml` per issue #196 option (b) because they are only consumed by IO functions in this module and by `cn_runtime_contract.ml::build_orchestrator_registry` which pattern-matches on `Loaded`/`Load_error`. Every IO function is untouched (`parse_file`, `discover`, `doctor_issues`, `execute_step`, `execute`, `typed_op_of_op_step`, `trace_event`, `find_step`, `env_get`, `as_bool`, `as_string`). Zero caller migration — 21 reference sites in `test/cmd/cn_workflow_test.ml` + 3 in `cn_runtime_contract.ml` all resolve through re-exports. Discipline verified by grep: `cn_workflow_ir.ml` imports only stdlib + `Cn_json`.
>
> **Move 2 status (v3.41.0) — ✅ Move 2 complete:** Fourth and final slice shipped — `src/lib/cn_frontmatter.ml` now owns the canonical 12 pure surface items from the activation validator: `frontmatter` record + `empty_frontmatter` constant + 6 line-level YAML-subset parsers (`split_lines`, `extract_block`, `parse_key_value`, `is_list_item`, `list_item_value`, `parse_frontmatter`) + `manifest_skill_ids` + the activation-validator's 3-variant `issue_kind` + `issue` record + `issue_kind_label`. Smallest Move 2 source module (288 LOC). `src/cmd/cn_activation.ml` re-exports all 3 types via OCaml type-equality and delegates all 9 pure functions via one-line let-bindings; the 3 IO functions (`read_skill_frontmatter`, `build_index`, `validate`) are untouched, as is the `activation_entry` re-export chain from slice 2 (which still routes through `Cn_contract`). Zero caller migration — 6 reference sites in `cn_doctor.ml` + ~20 sites in `test/cmd/cn_activation_test.ml` all resolve through re-exports or unchanged retained functions. Discipline verified by grep: `cn_frontmatter.ml` imports only stdlib + `Cn_json`. **With this slice, Move 2 is structurally complete: every pure type and parser in the codebase lives in `src/lib/`, every IO function lives in `src/cmd/`, and the Go kernel rewrite (#192) has its boundary contract — reimplement `src/cmd/` in Go using `src/lib/` as the type spec.** Remaining Move 2 candidates: **none**.

**Move 3 — Retire beta package doc (#180).** Replace `docs/beta/architecture/PACKAGE-SYSTEM.md` with a short redirect stub pointing to the alpha spec. Preserves old links while killing the authority conflict.

**Explicitly skipped for now:** package roles, richer activation schema, `src/host/` / `src/runtime/` split, gh-pages package source (#181), full directory-tree refactor. Reassess after Move 2 — if `src/lib/` gets crowded enough to justify `src/core/` / `src/host/` / `src/runtime/`, do it then with evidence.

---

## Known Debt

- Current `src/cmd/` extraction will take more than one cycle
- Built-in command shrinkage will be gradual
- Package role policy may need one refinement pass after first real packages
- Beta package-system architecture doc still needs to be brought fully in line with alpha artifact-first model (#180)

---

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Current docs describe a package-driven runtime more clearly than current code structure implements |
| 1 Select | — | — | Selected gap: internal code structure still centers command handlers instead of core/runtime registries and package-driven surfaces |
| 4 Gap | this artifact | — | Named incoherence: package/runtime modularity exists in theory but not yet in code structure |
| 5 Mode | this artifact | cdd/design, eng/evolve, eng/process-economics | L7 MCA; system-shaping refactor of internal architecture |
| 6 Artifacts | this artifact | — | Design artifact drafted; follow-up implementation plan and issue decomposition still required |

---

## Related

- #180 — Doc incoherence: beta claims Git-native transport
- #181 — gh-pages package source
- #174 — Orchestrator IR runtime (shipped v3.36.0)
- #167 — Package artifact distribution (shipped v3.34.0)
- #162 — Modular CLI commands
- #96 — Docs taxonomy alignment
- `docs/gamma/essays/MODULAR-ARCHITECTURE-REFACTOR.md` — Decision record (gamma)
