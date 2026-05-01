---
name: issue/contract
description: Write the contract sections of a cnos issue — problem, impact, status truth, source of truth, scope, non-goals, and closure condition.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ write the contract sections of an issue so α can act, β can verify, and γ can close without hidden context?
visibility: internal
parent: issue
triggers:
  - issue contract
  - problem statement
  - status truth
  - scope
  - non-goals
scope: task-local
inputs:
  - selected gap
  - affected surfaces
  - current implementation status
outputs:
  - problem section (exists / expected / diverges)
  - impact section
  - status truth table
  - source-of-truth table
  - scope section (in scope / out of scope / deferred / blocked by)
  - non-goals section
  - skills to load section
  - active design constraints section
  - related artifacts section
  - implementation guidance section
  - success / closure condition
---

# Issue Contract

## Core Principle

A contract section names one truth. It does not argue, qualify, or gesture. Each section answers one question so precisely that α can implement without asking and β can verify without recovering hidden conversation context.

Failure mode: **authority drift** — docs, implementation, and issue body disagree about what exists, what is required, and what is excluded.

---

## 1. Define

### 1.1 Identify the parts

Contract sections:

- **Problem** — the incoherence in three lines
- **Impact** — who is affected and what is blocked
- **Status truth** — what exists vs. what is expected vs. what is draft/planned
- **Source of truth** — one canonical path per load-bearing claim
- **Scope** — the execution boundary (in / out / deferred / blocked)
- **Non-goals** — what must not appear in ACs
- **Skills to load** — Tier 3 skills only
- **Active design constraints** — stable policy vs. volatile decisions
- **Related artifacts** — canonical links
- **Implementation guidance** — branch, identity, known migration debt
- **Success / closure condition** — when γ may close

### 1.2 Articulate how they fit

Problem → Impact establishes why this work matters.
Status truth → Source of truth establishes what is real now.
Scope + Non-goals establish the execution boundary.
Skills + Constraints establish the generation constraints for α.
Closure establishes when β's merge and γ's close are justified.

### 1.3 Name the failure mode

Contract fails through **authority drift**: two sections disagree about whether something is in scope, or the implementation cites a different canonical source than the issue does.

---

## 2. Unfold

### 2.1 Problem format

```markdown
## Problem

What exists:
What is expected:
Where they diverge:
```

Keep the problem 3–5 lines. Link to deeper evidence instead of reproducing design discussion.

- ❌ `The package system is inconsistent and needs cleanup.`
- ✅ `SKILL.md frontmatter has a CTB v0.1 schema in LANGUAGE-SPEC.md, but no CI gate validates the 54 package skill files. Missing fields can drift indefinitely.`

### 2.2 Impact format

Answer:

- Who or what is affected?
- What breaks or drifts?
- Why now?
- What is blocked if this remains unresolved?

- ❌ `This would be nice to have.`
- ✅ `Without this gate, CTB frontmatter conformance depends on reviewer memory. v0.2 promotion inherits unchecked v0.1 drift.`

### 2.3 Status truth

Distinguish shipped reality from expected or draft targets:

| Status | Meaning |
|---|---|
| `Shipped` | Implemented and enforced today |
| `Current spec` | Normative document currently governing conformance |
| `Draft target` | Proposed or draft spec not yet promoted |
| `Planned` | Accepted direction, not implemented |
| `Not in scope` | Explicitly excluded |
| `Unknown` | Not verified yet |

- ❌ Describe draft behavior as runtime behavior.
- ✅ `CTB v0.2 draft defines witnessed close-outs. The shipped runtime does not yet enforce them.`

### 2.4 Source-of-truth table

```markdown
## Source of truth

| Claim / surface | Canonical source | Status | Notes |
|---|---|---|---|
| CTB v0.1 frontmatter fields | docs/alpha/ctb/LANGUAGE-SPEC.md §2, §11 | current spec | governs current skill frontmatter |
```

Every load-bearing claim must link to a path, issue, or upstream document. No "the docs say" without a path.

### 2.5 Scope section

```markdown
## Scope

In scope:
- ...

Out of scope:
- ...

Deferred:
- ...

Blocked by:
- ...
```

Rules:

- Non-goals are mandatory for issues with 3+ ACs.
- Every noun in ACs must be in scope.
- If a non-goal noun appears in an AC, fix the issue before filing.
- If a deferred item is necessary for success, it is a dependency, not deferred.

### 2.6 Non-goals

Name what must not appear in ACs or implementation work. Non-goals exist to close the scope leak failure mode.

- ❌ Non-goal says "no runtime enforcement," but AC4 says "runtime enforces the schema."
- ✅ Non-goal says "no runtime enforcement," and AC4 says "CI rejects malformed frontmatter using a script."

### 2.7 Skills to load

Name Tier 3 skills only. Tier 1 and Tier 2 are mandatory and not repeated.

```markdown
## Skills to load

Tier 3:
- cnos.core/skills/design
- eng/tool

Why:
- design: boundary decisions and source-of-truth clarity.
- tool: shell script standards.
```

Do not list skills as decoration. If a skill changes the issue shape, explain why.

### 2.8 Active design constraints

State constraints in plain language.

- ❌ `Follow project conventions.`
- ✅ `Schema owns field shape/type/enum validation; shell script owns discovery and diagnostics. Do not duplicate schema rules in shell.`

### 2.9 Related artifacts

Link exact paths. If a companion artifact is being drafted in parallel, say how conflicts are resolved.

- ❌ `Related: CTB docs.`
- ✅ `docs/alpha/ctb/LANGUAGE-SPEC.md §2.1`

### 2.10 Implementation guidance

Use for: branch name, git identity, file paths, resolution base, known migration debt.

Do not use to replace a plan when the work requires multi-step execution.

```markdown
## Implementation guidance

Branch: cycle/324
Identity: alpha <alpha@cdd.cnos>
Notes: ...
```

### 2.11 Success / closure condition

State precisely when γ may close:

- all ACs met
- non-goals remain unviolated
- known gaps either resolved or carried as named debt

---

## 3. Rules

3.1. **State what exists before what is expected.** Every problem statement leads with current state, then target, then divergence.

3.2. **Status must be explicit.** Never describe draft behavior as shipped or planned enforcement as current.

3.3. **Source of truth is a path, not a claim.** Every load-bearing fact must link to a file or issue.

3.4. **Non-goals must not appear in ACs.** If a non-goal noun appears in an AC, fix the issue.

3.5. **Contradiction resolution is the author's job.** When two sections contradict, resolve before filing. α does not choose silently between contradictory interpretations.

---

## 4. Kata

### Scenario

You are filing an issue to add a CI gate that validates SKILL.md frontmatter against a CUE schema. The shipped runtime has no frontmatter validation. A draft CTB v0.2 spec proposes a witnessed close-out field.

### Task

Write the contract sections:

- problem (3–5 lines)
- status truth (table)
- source of truth (table with 3 rows)
- scope (in / out / deferred)
- non-goals (2 items)
- success / closure condition

### Expected artifacts

- Problem that names exists/expected/diverges without describing draft behavior as shipped.
- Status truth that distinguishes `Shipped`, `Current spec`, and `Draft target`.
- Source of truth that gives exact file paths.
- Scope that does not list any non-goal noun in "In scope."
- Non-goals that do not contradict ACs.

### Common failures

- Status truth says "CTB enforces witnessed close-outs" (draft claim as shipped).
- Source of truth says "CTB docs" without a path.
- Non-goal says "no v0.2 semantics" while an AC requires a v0.2 check.
- Scope says "in scope: full CI validation suite" while non-goals say "not in scope: CI enforcement."
