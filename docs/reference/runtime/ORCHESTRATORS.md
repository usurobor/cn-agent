# Orchestrators and Command Providers

**Version:** 0.1.0
**Status:** Draft
**Issue:** #170
**Doc-Class:** canonical-spec
**Canonical-Path:** docs/reference/runtime/ORCHESTRATORS.md
**Owns:** orchestrator model, command-provider model, activation index, compiled orchestrator IR, runtime contract additions
**Does-Not-Own:** low-level capability provider execution, package transport details, skill corpus standards

---

## 0. Purpose

This document defines a second extension axis for cnos:
- not runtime capability providers,
- not skills,
- but package-distributed commands and orchestrators.

The goal is to make more of cnos mechanical:
- skills teach judgment
- commands expose exact operator actions
- orchestrators execute declared workflows
- runtime extensions provide low-level capabilities

This document keeps those surfaces separate and composable.

---

## 1. Problem

cnos currently has:
- skills for judgment
- runtime extensions for capability growth
- built-in commands for operator control

But it lacks a coherent model for package-distributed:
- operator-facing commands
- mechanical workflows that sequence LLM calls and typed ops

That creates three pressures:
1. More commands drift into core than necessary.
2. Reusable workflows have no first-class runtime form.
3. "Skill activation" risks becoming vague keyword matching instead of explicit runtime policy.

---

## 2. Core Decisions

### 2.1 An orchestrator is not a skill

A skill teaches repeatable judgment in a domain. An orchestrator is a mechanical workflow:
- sequence
- bind
- branch
- retry
- fan-out / fan-in
- call the LLM
- call typed runtime ops
- emit artifacts / receipts

A skill may choose an orchestrator. An orchestrator does not replace the skill.

### 2.2 A computation expression is an authoring surface, not the runtime artifact

Humans may author orchestrators in a host DSL. F# computation expressions are a preferred authoring surface for non-trivial orchestrators. The runtime does not execute "F#." The runtime executes a compiled orchestrator IR.

### 2.3 Commands and orchestrators share package infrastructure

Package artifacts may distribute:
- skills
- commands
- orchestrators
- templates
- doctrine / mindsets
- runtime extensions

Commands and orchestrators are package content classes. They use the same package install/restore substrate underneath.

### 2.4 Runtime extensions remain separate

Runtime capability providers remain the existing extensions system. They are:
- typed capability providers
- runtime/body-plane components
- not skills
- not commands
- not orchestrators

### 2.5 Activation must be declarative, not just keyword-based

Skill activation uses an activation index built from:
- events
- work shapes
- path hints
- level/risk hints
- token hints

Keywords are allowed as weak signals. They are not the whole activation model.

---

## 3. New Package Content Classes

Add two explicit package content classes:
- commands
- orchestrators

The package-system content class set becomes:
- doctrine
- mindsets
- skills
- commands
- orchestrators
- extensions
- templates
- metadata

---

## 4. Runtime Contract Changes

### 4.1 Cognition

Add:
- activation_index

Purpose:
- declares which skills may activate under which conditions
- makes activation visible at wake
- allows the agent to reason from explicit activation metadata rather than hidden prompt conventions

Suggested shape:

```json
"cognition": {
  "installed_packages": [...],
  "active_overrides": {
    "doctrine": [],
    "mindsets": [],
    "skills": []
  },
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
}
```

### 4.2 Body

Add:
- commands
- orchestrators

Purpose:
- these are executable surfaces the runtime can dispatch
- they belong to the body plane, not the cognition plane

Suggested shape:

```json
"body": {
  "capabilities": { ... },
  "peers": ["sigma"],
  "commands": [
    {
      "name": "daily",
      "source": "package",
      "package": "cnos.pm",
      "summary": "Create or show the daily reflection thread"
    }
  ],
  "orchestrators": [
    {
      "name": "daily-review",
      "source": "package",
      "package": "cnos.pm",
      "trigger_kinds": ["command", "schedule"]
    }
  ]
}
```

---

## 5. Command Provider Model

### 5.1 What a command is

A command is an operator-facing executable entrypoint resolved by exact name. Commands are not activated by keyword. They are dispatched directly.

### 5.2 Discovery precedence

v1 precedence:
1. built-in commands
2. repo-local commands
3. package commands

No PATH discovery in v1.

### 5.3 Repo-local commands

Repo-local commands live under:

```
.cn/commands/
  cn-daily
  cn-weekly
```

These are suitable for project-local workflows.

### 5.4 Package commands

Package commands are declared in cn.package.json and installed into the vendored package tree.

Suggested metadata:

```json
"commands": {
  "daily": {
    "entrypoint": "commands/cn-daily",
    "summary": "Create or show the daily reflection thread",
    "needs_hub": true,
    "dangerous": false
  },
  "weekly": {
    "entrypoint": "commands/cn-weekly",
    "summary": "Create or show the weekly review thread",
    "needs_hub": true,
    "dangerous": false
  }
}
```

### 5.5 Help and doctor

cn help lists:
- built-ins
- repo-local commands
- package commands

cn doctor validates:
- duplicate command names
- missing entrypoints
- non-executable files
- malformed metadata

---

## 6. Skill Activation Model

### 6.1 What activation is

Activation is the process of turning installed skills into candidate skills for the current context.

### 6.2 Activation inputs

A skill may activate from:
- event kind
- work shape
- path hints
- level hints
- dominant risk
- token hints
- requires / excludes

### 6.3 Example activation declaration

```json
"activation": {
  "events": ["cdd:step=5", "command:design"],
  "tokens": ["impact graph", "acceptance criteria", "proposal"],
  "paths": ["docs/**", "src/**"],
  "work_shapes": ["design", "docs/process", "runtime/platform"],
  "levels": ["L6", "L7"],
  "requires": ["testing"],
  "excludes": ["release"]
}
```

### 6.4 Rule

Keywords are weak signals. Events, work shape, path, and explicit constraints are stronger.

---

## 7. Orchestrator Model

### 7.1 What an orchestrator is

An orchestrator is a mechanical workflow that can:
- sequence steps
- bind outputs to later inputs
- branch on results
- call the LLM
- call typed runtime ops
- emit artifacts
- leave receipts / trace events

An orchestrator is not:
- a skill
- a runtime extension
- arbitrary executable code produced by the agent

### 7.2 Authoring stack

Three authoring surfaces, one canonical IR:

| Surface | Author | Strength |
|---------|--------|----------|
| **CTB (Coherent Triadic Binding)** | LLMs + Humans | Deterministic, expression-only, effects-as-data. Small surface area, canonical formatting, LLM-friendly. See [CTB-v4.0.0-VISION.md](../ctb/CTB-v4.0.0-VISION.md) |
| **F# computation expressions** | Humans | Sequencing, binding, branching, effect boundaries — compact and type-safe |
| **YAML** | Humans | Simple workflows, low ceremony |

All three compile to the same canonical IR (`cn.orchestrator.v1` JSON).

```
CTB / F# CE / YAML  →  compiler  →  cn.orchestrator.v1 JSON  →  runtime
  (authoring)            (build)         (canonical IR)          (execution)
```

**CTB** is the canonical authoring language for orchestrators. A CTB program evaluates to an **Effect Plan** — a pure data term describing what should happen — which maps directly to the orchestrator IR step graph. This is the connection between CTB (§4: separation of logic and effects) and the orchestrator runtime: CTB programs produce effect plans, the runtime executes them with explicit capabilities.

CTB is the preferred authoring surface because it is:
- **Deterministic** — same inputs, same outputs (CTB-v4.0.0-VISION §3.1)
- **Misuse-resistant** — no ambiguous clause overlap, no hidden state (§3.1.3)
- **LLM-friendly** — small surface area, predictable syntax, witness-based errors (§3.1.5)
- **Effects-as-data** — programs return effect plans, not side effects (§3.1.4)

F# computation expressions and YAML remain as alternative authoring surfaces for humans who prefer them. All compile to the same IR.

**Natural language → CTB** is a future planner step, not a v1 requirement. The safe sequence:
1. Define canonical JSON IR (#174)
2. Build CTB → IR compiler (#175)
3. Let LLMs generate CTB directly
4. Later: natural language → CTB as a planner step

Do not make natural language compile directly to runtime IR. That bundles planning, normalization, capability selection, control-flow construction, and validation into one implicit step.

### 7.3 Three-layer execution model

CTB is the source language, not the runtime IR. Three layers, each with a distinct role:

| Layer | What | Form | Deterministic? |
|-------|------|------|---------------|
| **CTB program** | Source / authoring language | `.ctb` source | Yes — pure, equation-only |
| **Effect Plan** | Semantic IR / runtime contract | JSON (`cn.effect-plan.v1`) | Yes — pure data, hashable, verifiable |
| **Runtime bridge** | Executor with capabilities | `cn` binary | No — impure at effect boundaries |

```
CTB source  →  compile/evaluate  →  Effect Plan (JSON)  →  runtime bridge execution
  (pure)          (pure)              (pure data)           (capability-gated)
```

The Effect Plan is the actual IR — pure data, capability-declarable, runtime-interpretable, reproducible. CTB programs evaluate to effect plans. The runtime bridge interprets effect plans under explicit capability policy. Peers can re-run the same CTB program and get the same plan.

This matches CTB-v4.0.0-VISION.md exactly:
- CTB programs are equations over triadic terms (§3.1)
- A CTB "skill" is a pure function `State → EffectPlan` (§4)
- The runtime bridge interprets the plan with explicit capabilities (§4)
- The runtime becomes a driver, not the source of truth

### 7.4 Why JSON as the wire/storage encoding

The Effect Plan is the semantic IR. JSON is its serialized wire/storage form:
- Package manifests are already JSON
- Typed ops manifests are already JSON
- Runtime validation trusts the serialized plan, not the source language
- Machine artifacts benefit from rigid structure
- Plans are hashable and diffable in JSON

JSON is not the authoring surface — CTB is. JSON is one layer down: the compiled, validated, transportable form of the effect plan.

### 7.5 What CTB is for cnos

CTB is the programmable skill language for cnos. The CTB vision doc says it explicitly:
- "Skills are programs, not prose" (§1)
- "A shared skills layer written in Markdown or prose cannot be executed consistently" (§1)
- "We need a skills substrate where peers can re-run the same program and get the same plan" (§1)

This means:
- **Markdown skills** remain as human-facing judgment docs (CAP, review, design)
- **CTB programs** become the canonical source for executable workflow logic / orchestrators
- **Commands** dispatch into CTB where appropriate — `cn daily` → CTB program → Effect Plan → runtime execution

The execution path:
1. Command selected explicitly by operator
2. Command may invoke a CTB orchestrator
3. CTB produces effect plan (pure data)
4. Runtime bridge executes under capability policy

---

## 8. Orchestrator IR

### 8.1 Minimal shape

```json
{
  "kind": "cn.orchestrator.v1",
  "name": "daily-review",
  "trigger": {
    "kind": "command",
    "name": "daily"
  },
  "inputs": {
    "requires_hub": true
  },
  "permissions": {
    "llm": true,
    "ops": ["fs_read", "fs_write", "git_status"],
    "external_effects": false
  },
  "steps": [
    {
      "id": "load-context",
      "kind": "op",
      "op": "fs_read",
      "args": { "path": "threads/reflections/daily/today.md" },
      "bind": "today"
    },
    {
      "id": "summarize",
      "kind": "llm",
      "prompt": "daily-review-v1",
      "inputs": ["today"],
      "bind": "summary"
    },
    {
      "id": "write-output",
      "kind": "op",
      "op": "fs_write",
      "args": { "path": "threads/adhoc/daily-summary.md" },
      "inputs": ["summary"]
    },
    {
      "id": "return",
      "kind": "return",
      "value": { "artifact": "threads/adhoc/daily-summary.md" }
    }
  ]
}
```

### 8.2 Step kinds

v1 step kinds:
- `op`
- `llm`
- `if`
- `match`
- `parallel`
- `return`
- `fail`

### 8.3 Determinism rule

The orchestrator control structure must be deterministic. Non-determinism is allowed only at explicit effect boundaries:
- `llm`
- runtime ops with external state

Execution must leave receipts / events sufficient to replay the control flow.

---

## 9. Package Manifest Extensions

Extend cn.package.json with:
- commands
- orchestrators
- optional activation_index

Suggested full example:

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.pm",
  "version": "3.40.0",
  "kind": "package",
  "engines": { "cnos": "3.40.0" },
  "sources": {
    "skills": ["pm/daily", "pm/weekly"],
    "commands": ["daily", "weekly"],
    "orchestrators": ["daily-review"],
    "templates": []
  },
  "commands": {
    "daily": {
      "entrypoint": "commands/cn-daily",
      "summary": "Create or show the daily reflection thread"
    }
  },
  "orchestrators": {
    "daily-review": {
      "entrypoint": "orchestrators/daily-review/orchestrator.json",
      "summary": "Run the daily review workflow",
      "trigger_kinds": ["command", "schedule"]
    }
  },
  "activation_index": {
    "skills": [
      {
        "id": "pm/daily",
        "activation": {
          "events": ["command:daily"],
          "tokens": ["daily", "reflection"],
          "paths": ["threads/reflections/**"]
        }
      }
    ]
  }
}
```

---

## 10. Build / Install Layout

### 10.1 Source layout

```
src/agent/
  skills/
  commands/
  orchestrators/
  extensions/
  templates/
  doctrine/
  mindsets/
```

### 10.2 Installed package layout

```
packages/<name>/
  skills/
  commands/
  orchestrators/
  extensions/
  templates/
  cn.package.json
```

### 10.3 Vendored layout

```
.cn/vendor/packages/<name>@<version>/
  skills/
  commands/
  orchestrators/
  extensions/
  templates/
  cn.package.json
```

---

## 11. Execution Rules

### 11.1 Skills
- activated by context through the activation index
- remain judgment surfaces
- may recommend commands or orchestrators

### 11.2 Commands
- invoked by exact operator name
- may call an orchestrator
- may also remain small standalone entrypoints

### 11.3 Orchestrators
- selected by trigger contract or called explicitly
- execute mechanically
- may call LLM and typed ops
- may not bypass runtime policy

### 11.4 Runtime extensions
- provide low-level typed capabilities
- remain separate and policy-governed

---

## 12. Invariants

1. Orchestrators are not skills.
2. Computation expressions are source forms, not runtime contracts.
3. Commands are exact-dispatch surfaces, not keyword activations.
4. Skills activate from declarative context signals, not trigger words alone.
5. Commands and orchestrators are package content classes.
6. Runtime extensions remain separate from commands and orchestrators.
7. The runtime validates compiled manifests, not authoring languages.

---

## 13. Non-goals

This design does not:
- add PATH-based command discovery in v1
- collapse runtime extensions into package commands
- let the agent emit arbitrary executable code
- require JSON as a hand-authored source language
- define one mandatory host DSL

---

## 14. Short form

- Skills decide.
- Commands dispatch.
- Orchestrators execute.
- Extensions provide capability.
- Packages distribute all of them.
