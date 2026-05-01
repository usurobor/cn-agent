---
name: issue/constraints
description: Write the constraint sections of a cnos issue — constraint strata, exceptions, path resolution, and cross-surface projections.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ define and document constraints so α cannot inadvertently violate them and β can verify compliance?
visibility: internal
parent: issue
triggers:
  - constraint strata
  - exceptions
  - path resolution
  - cross-surface
scope: task-local
inputs:
  - scope
  - acceptance criteria
  - affected runtime surfaces
outputs:
  - constraint strata section (hard gate / exception-backed / optional / validated-if-present / deferred)
  - exceptions ledger
  - path resolution rule with example
  - cross-surface projection list
---

# Issue Constraints

## Core Principle

Constraints protect invariants. A constraint that cannot be violated is not a constraint — it is a fact. A constraint that can be violated but goes undetected is silent technical debt.

Failure mode: **exception theater** — exceptions hide the rule instead of documenting temporary debt.

---

## 1. Define

### 1.1 Identify the parts

The constraint surface of an issue has:

- **Constraint strata** — what is hard-gated vs. exception-backed vs. optional vs. validated-if-present vs. deferred
- **Exceptions ledger** — field-specific exceptions with reason and removal condition
- **Path resolution** — how implementers resolve relative paths
- **Cross-surface projections** — all surfaces that must update when one surface changes

### 1.2 Articulate how they fit

Strata name the enforcement model. Exceptions document permitted violations and their removal conditions. Path resolution prevents mis-implementation when multiple bases are possible. Cross-surface projections prevent incomplete updates.

### 1.3 Name the failure mode

Constraints fail through:

- **Exception blanket** — one exception covers a whole file instead of specific fields
- **Stale path base** — issue says "relative to the skill's directory" but no example is given
- **Missing projection** — CI job added but notification aggregation not updated
- **Deferred field as hard gate** — a field marked deferred appears in a hard-gate check

---

## 2. Unfold

### 2.1 Constraint strata

Stratify every rule the issue defines or enforces:

| Stratum | Meaning |
|---|---|
| **Hard gate** | Failure if missing or violated; no exception |
| **Exception-backed** | Failure unless explicitly listed with reason |
| **Optional/defaulted** | May be omitted; default or absence semantics defined |
| **Validated if present** | Not required, but checked when declared |
| **Ignored/deferred** | Explicitly not checked in this issue |

```markdown
## Constraint strata

Hard gate:
- name
- description
- governing_question

Exception-backed:
- artifact_class
- kata_surface

Optional/defaulted:
- visibility (default: internal)

Validated if present:
- calls
- calls_dynamic

Ignored/deferred:
- CTB v0.2 witnessed close-out semantics
```

Rules:

- Hard-gate fields cannot appear in exceptions.
- Exception-backed fields require field-specific (not file-level) exceptions.
- Optional/defaulted fields must name the default.
- Validated-if-present fields must name what is checked.
- Deferred fields must be explicitly listed as not enforced in this issue.

- ❌ Hard gate: `artifact_class` AND exception example shows `artifact_class` as `allowed_missing`.
- ✅ Hard gate: `name`, `description`; Exception-backed: `artifact_class` with a specific path and reason.

### 2.2 Exceptions ledger

Exceptions are debt, not permission to ignore the rule.

```json
[
  {
    "path": "src/packages/cnos.core/skills/example/SKILL.md",
    "allowed_missing": ["artifact_class", "kata_surface"],
    "reason": "Legacy skill; v0.1 migration pending",
    "spec_ref": "docs/alpha/ctb/LANGUAGE-SPEC.md §2.1",
    "owner": "cnos.core",
    "remove_when": "frontmatter migration completes"
  }
]
```

Rules:

- Exceptions are field-specific, not whole-file blanket passes.
- Every exception needs a reason.
- Include `remove_when` or a follow-up issue reference.
- Exception examples must not contradict hard-gate rules.

- ❌ `{ "path": "src/...", "ignore": true }`
- ✅ `{ "path": "src/...", "allowed_missing": ["artifact_class"], "reason": "Legacy skill", "remove_when": "migration completes" }`

### 2.3 Path resolution semantics

When an issue asks implementers to validate paths, load files, resolve calls, or traverse directories, it must define the resolution base.

Questions to answer:

- Is the path repo-root-relative? Package-root-relative? Caller-file-relative?
- CWD-relative? Branch-relative? Generated-artifact-relative?

Include at least one concrete resolution example.

- ❌ `Validate each target file exists relative to the skill's parent directory.`
- ✅ `Validate static calls targets relative to the package skill root. Example: src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md with calls: design/SKILL.md resolves to: src/packages/cnos.cdd/skills/cdd/design/SKILL.md`

### 2.4 Cross-surface projections

When an issue changes a surface that appears in more than one place, name all projections.

| Change | Surfaces to update |
|---|---|
| CI job added | workflow job; notification aggregation; local command; expected diagnostics; link/check job |
| README source map changed | root README; docs/README; link checker; status labels |
| Schema added | schema; fixtures; exceptions; docs; CI invocation |
| New command or runtime surface | help text; doctor; status; operator docs |

- ❌ `Add CI job I5.`
- ✅ `Add CI job I5 and update notify job status aggregation so failure summaries include I5.`

---

## 3. Rules

3.1. **Hard gates have no exceptions.** If a rule admits exceptions, it is not a hard gate.

3.2. **Exceptions are field-specific.** Blanket file-level passes are incoherent. Exceptions must name the specific fields allowed to be missing.

3.3. **Every exception has a reason and a removal condition.** Exceptions without `remove_when` or a follow-up issue are permanent debt.

3.4. **Path resolution must name the base and give an example.** Prose that says "relative to X" without an example will produce mis-implementation.

3.5. **Cross-surface projections must be enumerated.** If a change affects CI, docs, runtime status, commands, schemas, or source maps, name every surface that must also update.

---

## 4. Kata

### Scenario

You are filing an issue to add a CUE schema for SKILL.md frontmatter validation. The schema has three hard-gate fields (`name`, `description`, `governing_question`) and two exception-backed fields (`artifact_class`, `kata_surface`). The CI job validates files under `src/packages/`. A CI notification job aggregates job statuses.

### Task

Write:

- constraint strata (all five levels, with at least one entry each or explicit "none")
- one exception ledger entry for a legacy skill missing `artifact_class`
- path resolution rule with example
- cross-surface projection list

### Expected artifacts

- Strata that do not list a hard-gate field in the exceptions entry.
- Exception with `allowed_missing`, `reason`, and `remove_when`.
- Path resolution with a concrete before/after example.
- Cross-surface list that includes the notification aggregation job.

### Common failures

- Hard gate lists `artifact_class`, but the exception entry also lists `artifact_class` as `allowed_missing`.
- Exception is `{ "ignore": true }` without field specificity.
- Path resolution says "relative to the skill root" without an example.
- Cross-surface list names the CI job but not the notification aggregation.
