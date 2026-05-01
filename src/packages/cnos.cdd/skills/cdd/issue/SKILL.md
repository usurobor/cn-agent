---
name: issue
description: Write an executable issue pack that names the incoherence, source-of-truth, constraints, proof obligations, required skills, and implementation guidance without creating ambiguity or false closure.
artifact_class: runbook
kata_surface: external
governing_question: How does γ turn a selected gap into an issue pack that α and β can execute, verify, and close without hidden assumptions?
visibility: internal
parent: cdd
triggers:
  - issue
scope: task-local
inputs:
  - selected gap
  - mode
  - active design constraints
  - Tier 3 skills
  - affected surfaces
  - related artifacts
  - current implementation status
outputs:
  - executable issue pack
requires:
  - γ completed observe/select
  - canonical CDD.md loaded
calls:
  - issue/labels/SKILL.md
  - issue/contract/SKILL.md
  - issue/proof/SKILL.md
  - issue/constraints/SKILL.md
kata_ref: src/packages/cnos.cdd.kata/katas/M5-issue-authoring/
---

# Issue

## Core Principle

A coherent issue is an executable work contract. It names:

1. the incoherence;
2. the impact;
3. the status of the affected system;
4. the source of truth;
5. the implementation boundary;
6. the acceptance criteria;
7. the proof or rejection mechanism;
8. the non-goals.

A coherent issue lets an engineer who was not in the conversation act without asking clarifying questions.

Failure modes:

- **ambiguity** — engineer must ask what the issue means;
- **overclaiming** — issue states target/draft/spec behavior as if it already ships;
- **false closure** — issue has good prose but no proof or rejection mechanism;
- **source drift** — issue cites the wrong canonical doc, path, or runtime surface;
- **scope leak** — non-goals reappear inside acceptance criteria or work steps;
- **exception theater** — exceptions hide the rule instead of documenting temporary debt.

Operational shape:

```text
Issue = tri(gap, execution boundary, verification surface)
```

The problem names the gap.
The scope and constraints name the execution boundary.
The acceptance criteria and proof plan name the verification surface.

---

## When to load each subskill

Load subskills as the issue shape requires:

- **`issue/labels`** — when selecting or reviewing kind, priority, or surface labels for the issue title and header.
- **`issue/contract`** — when writing problem, impact, status truth, source of truth, scope, non-goals, skills, constraints, related artifacts, implementation guidance, or closure condition.
- **`issue/proof`** — when writing acceptance criteria, proof plan, or oracle/positive/negative cases.
- **`issue/constraints`** — when writing constraint strata, exceptions ledger, path resolution rules, or cross-surface projection lists.

For a simple issue with few ACs, loading `issue/labels` and `issue/contract` may be sufficient. For complex tooling/CI/runtime issues, load all four.

---

## Minimal output pattern

```markdown
# Title
Labels: <kind>, <priority>[, surface…]

Priority: <P0–P3> — <one-line rationale>
Status: <one-line summary>

## Problem

What exists:
What is expected:
Where they diverge:

## Impact

## Status truth

## Source of truth

| Claim / surface | Canonical source | Status | Notes |
|---|---|---|---|

## Scope

In scope:
Out of scope:
Deferred:

## Acceptance criteria

### AC1: …

Invariant:
Oracle:
Positive:
Negative:
Surface:

## Proof plan

Invariant:
Surface:
Oracle:
Positive case:
Negative case:
Operator-visible projection:
Known gap:

## Skills to load

Tier 3:
- …

Why:
- …

## Active design constraints

## Related artifacts

## Non-goals

## Success / closure condition

This issue is closeable when:
- all ACs are met;
- the proof plan passes;
- non-goals remain unviolated;
- known gaps are either resolved or explicitly carried as named debt.
```

For small issues, compress the shape but preserve: problem, impact, ACs, non-goals (if substantial), source of truth, and proof surface.

---

## Closure rule

An issue is closed when:

- every AC maps to evidence in the branch diff;
- the proof plan oracle produces a passing result;
- no non-goal noun appears in the diff;
- known gaps are named explicitly in `self-coherence.md`.

---

## Handoff checklist

Before filing or dispatching:

- [ ] Problem states exists / expected / divergence.
- [ ] Impact says who cares and what is blocked.
- [ ] Status truth is explicit (no draft behavior described as shipped).
- [ ] Source-of-truth paths resolve.
- [ ] Scope and non-goals do not contradict ACs.
- [ ] Contradictions between issue sections are resolved explicitly.
- [ ] Hard gates do not appear in exception examples.
- [ ] ACs are numbered and independently testable.
- [ ] Proof plan has oracle, positive case, negative case.
- [ ] New surfaces include operator-visible projections.
- [ ] CI additions include notification/status implications.
- [ ] Related artifacts include canonical docs with exact paths.
- [ ] Examples obey the rules they are demonstrating.
- [ ] Known gaps are named honestly.
- [ ] Labels: exactly one kind label and one priority label (load `issue/labels`).

---

## External kata

Practice and evaluation for this issue skill live in:

`src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`

That kata exercises status-truth split, source-of-truth alignment, constraint strata, AC negative space, proof-plan oracle, and non-goal discipline.
