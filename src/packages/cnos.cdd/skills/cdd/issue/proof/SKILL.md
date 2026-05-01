---
name: issue/proof
description: Write acceptance criteria and a proof plan for a cnos issue — invariant, oracle, positive/negative cases, surface, and known gaps.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ write ACs and a proof plan so α can prove closure and β can verify without recovering hidden context?
visibility: internal
parent: issue
triggers:
  - acceptance criteria
  - proof plan
  - AC
  - oracle
scope: task-local
inputs:
  - problem statement
  - scope
  - non-goals
outputs:
  - numbered acceptance criteria (invariant / oracle / positive / negative / surface)
  - proof plan (invariant / surface / oracle / positive / negative / operator-visible projection / known gap)
---

# Issue Proof

## Core Principle

An AC names done precisely enough that there is no ambiguity between a passing and a failing state. A proof plan names the oracle that decides the outcome.

Failure mode: **false closure** — the AC has good prose but no mechanism by which any reviewer can say "this passes" or "this fails."

---

## 1. Define

### 1.1 Identify the parts

An acceptance criterion has:

- **Invariant** — the property that must hold
- **Oracle** — the specific mechanism that determines pass or fail
- **Positive case** — what a passing artifact looks like
- **Negative case** — what a failing artifact looks like
- **Surface** — the file, command, or runtime surface where the oracle runs

A proof plan has:

- **Invariant** — the umbrella property the whole plan proves
- **Surface** — the files, commands, or runtime surfaces
- **Oracle** — how pass/fail is determined
- **Positive case** — what success looks like across surfaces
- **Negative case** — what failure looks like
- **Operator-visible projection** — how an operator sees the result (CI status, command output, error message)
- **Known gap** — what this proof plan does not cover

### 1.2 Articulate how they fit

ACs name discrete, testable properties. The proof plan aggregates them into a single verifiable claim about the system.

A missing oracle is a false closure: the AC looks like done, but no one can test it.

### 1.3 Name the failure mode

Proof fails through:

- **No oracle** — AC says "it should work" without naming who checks what
- **Missing negative case** — only happy-path coverage; no evidence that violations are caught
- **Scope leak** — AC tests something named in non-goals
- **Oracle drift** — proof plan references a surface that does not exist

---

## 2. Unfold

### 2.1 AC shape

```markdown
### AC1: <title>

Invariant: <the property that must hold>
Oracle: <what checks it and how>
Positive: <what a passing case looks like>
Negative: <what a failing case looks like>
Surface: <file, command, or CI job>
```

Rules:

- ACs are numbered and independently testable.
- At least one AC must cover a negative case for validation/checker/tooling issues.
- ACs must not duplicate the full implementation plan.
- ACs must not include out-of-scope nouns.

- ❌ `AC1: Schema works.`
- ✅ `AC1: Missing scope causes scripts/check-skill-frontmatter.sh to exit non-zero naming the file and missing field.`

**Shortened AC shape** (acceptable when the oracle is obvious):

```markdown
### AC1: <title>

<one sentence naming what is true when this AC passes, with a testable claim>
```

Use the full shape for complex/tooling/runtime issues. Use the short shape only when the oracle is unambiguous from context.

### 2.2 Proof plan

```markdown
## Proof plan

Invariant: <umbrella property>
Surface: <files and commands>
Oracle: <how pass/fail is decided>
Positive case: <what success looks like>
Negative case: <what failure looks like>
Operator-visible projection: <how an operator sees the result>
Known gap: <what this plan does not cover>
```

Rules:

- Existing CI passing is the floor, not the proof.
- Every new surface requires at least one CI-gated or documented assertion.
- If proof is partial, say what remains unproven in "Known gap."
- The oracle must be a specific check, not "it works."

- ❌ `Oracle: test passes`
- ✅ `Oracle: script exits non-zero and names the missing field when given a fixture missing 'scope'`

### 2.3 Negative space requirement

For issues that add validators, checkers, parsers, CI jobs, runtime behavior, or schemas: at least one AC or the proof plan must cover the negative case.

A validator that is never tested on bad input is not proven.

- ❌ All ACs show the happy path. No AC asks "what happens when it fails?"
- ✅ AC2 uses a fixture missing a required field and asserts the check rejects it.

### 2.4 Operator-visible projection

When operators rely on a check (CI output, command error message, status signal), name how they see the result.

- ❌ Proof plan proves the check exists but does not say how operators know it failed.
- ✅ `Operator-visible projection: CI job I5 status appears in the notify aggregation; failure names the file and field.`

---

## 3. Rules

3.1. **Every AC is independently testable.** One AC should be able to pass while another fails.

3.2. **Every AC names an oracle.** "It should work" is not an oracle.

3.3. **Negative case is mandatory for checker/validator/CI issues.** At least one AC must show what happens when the constraint is violated.

3.4. **Proof plan must name known gaps.** If proof is partial, say so.

3.5. **ACs must not contain out-of-scope nouns.** If a non-goal appears in an AC, fix the issue.

3.6. **Oracle drift is a finding.** If the proof plan references a surface that does not exist or was renamed, β will flag it.

---

## 4. Kata

### Scenario

You are adding a CI job that validates SKILL.md frontmatter. Required fields are: `name`, `description`, `artifact_class`. The check script exits non-zero on missing fields.

### Task

Write:

- AC1: valid frontmatter passes
- AC2: missing `artifact_class` fails
- proof plan (all seven fields)

### Expected artifacts

- AC1 with oracle (script exit code) and positive case (fixture with all fields).
- AC2 with oracle (script exit code + diagnostic) and negative case (fixture missing `artifact_class`).
- Proof plan with non-empty "Negative case" and "Known gap."

### Common failures

- AC1 oracle says "tests pass" without naming the command.
- AC2 has no negative case.
- Proof plan says "Known gap: none" when the issue explicitly defers v0.2 semantics.
- Oracle drift: proof plan names a script path that the issue has not yet specified as final.
