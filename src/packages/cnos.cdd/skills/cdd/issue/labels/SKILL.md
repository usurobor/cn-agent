---
name: issue/labels
description: Issue label taxonomy — kind, priority, and optional surface labels. Use when selecting or reviewing labels for a cnos issue.
artifact_class: reference
kata_surface: none
governing_question: Which labels does a cnos issue require, and what do they mean?
visibility: internal
parent: issue
triggers:
  - label
  - kind label
  - priority label
scope: task-local
inputs:
  - issue title and problem statement
outputs:
  - one kind label
  - one priority label
  - zero or more surface/domain labels
---

# Issue Labels

## Rule

Every issue requires exactly one kind label and exactly one priority label. Surface/domain labels are optional.

---

## Kind labels

| Label | Use when |
|---|---|
| `bug` | Shipped behavior is broken, regressed, misleading, or incoherent with its contract. |
| `enhancement` | Adds or improves capability, behavior, package surface, workflow, or contract. |
| `docs` | Primary artifact is documentation, source map, README, spec wording, operator guidance, or status-truth correction. |
| `design` | Primary work is deciding architecture, doctrine, contract shape, or boundary before implementation. |
| `test` | Primary work is tests, fixtures, katas, validation, checker coverage, or proof surface. |
| `refactor` | Internal restructuring with no intended behavior or contract change. |
| `chore` | Maintenance, cleanup, dependency, CI hygiene, release plumbing, or non-user-facing upkeep. |
| `tracking` | Umbrella or coordination issue that links child issues and does not itself implement the work. |

Exactly one kind label per issue.

- ❌ An issue with no kind label, multiple kind labels, or an unlisted kind without justification.
- ✅ An issue with exactly one of the eight kind labels above.

---

## Priority labels

| Label | Meaning |
|---|---|
| `P0` | Emergency. Active breakage, data/security risk, release-blocking failure, or no safe workaround. |
| `P1` | High. Blocks important planned work, onboarding, correctness, status truth, or core architecture. Should be scheduled soon. |
| `P2` | Normal. Valuable scoped work, coherent enhancement, cleanup, package/doc improvement, or non-blocking system refinement. Default for most real work. |
| `P3` | Low. Nice-to-have, polish, exploratory, backlog, or work safe to defer indefinitely. |

Exactly one priority label per issue.

- ❌ An issue with no priority label, multiple priority labels, or P0/P1 without a blocking/emergency rationale.
- ✅ An issue with exactly one priority label that matches the urgency described in the problem and impact.

---

## Surface/domain labels (optional)

Surface labels help routing. They do not replace source-of-truth paths or ACs.

Suggested surface labels:

| Label | Scope |
|---|---|
| `cli` | Command-line interface, operator-facing commands |
| `core` | Core runtime, kernel, bootstrap |
| `cdd` | CDD lifecycle, roles, skills, protocol |
| `ctb` | CTB language spec, frontmatter, module system |
| `cdw` | CDW writing/doc system |
| `runtime` | Runtime behavior, execution, platform |
| `package` | Package system, distribution, install |
| `ci` | CI/CD, workflows, automation |
| `security` | Security, auth, permissions |

Rules:

- Multiple surface labels are allowed.
- Keep surface labels minimal — only what helps routing.
- Surface labels must not duplicate a kind label. For a CTB docs issue, use `docs, P2, ctb` rather than `docs, P2, docs`.

- ❌ Surface label `docs` on an issue that already carries kind label `docs`.
- ✅ Surface label `ctb` on an issue with kind label `docs` (the kind is docs; the surface is ctb).

---

## Discipline

Labels do not replace issue content. An issue with correct labels but no problem statement, ACs, or scope is still incoherent.

- ❌ Rely on labels to signal scope instead of writing a scope section.
- ✅ Labels route; the issue body governs.

---

## Examples

```
docs: README activation step       → docs, P1
feat(cli): cn activate             → enhancement, P1, cli
package(cdw): writing package      → enhancement, P2
tracking umbrella                  → tracking, P2
skill(cdd/issue): split subskills  → refactor, P2, cdd
bug: cn doctor reports false green → bug, P0
```
