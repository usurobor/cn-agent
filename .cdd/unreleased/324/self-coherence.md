# Self-Coherence — Issue #324

## §Gap

**Issue:** [#324 — skill(cdd/issue): split issue skill into focused subskills and add label taxonomy](https://github.com/usurobor/cnos/issues/324)

**Version/mode:** CDD 3.15.0 — MCA (skill refactor + new label taxonomy subskill)

**Named incoherence:**

`cdd/issue/SKILL.md` has grown into a large all-in-one skill carrying label taxonomy, AC/proof rules, constraint strata, exception discipline, path resolution, cross-surface updates, handoff checks, and katas inline. The root skill loads too much context when only one concern is needed and has multiple reasons to change. New label rules would make it larger if added directly.

**Target state:** Root `issue/SKILL.md` becomes a small orchestrator that delegates to focused subskills: `labels/`, `contract/`, `proof/`, `constraints/`. Label taxonomy lives in `issue/labels/SKILL.md` and requires exactly one kind label and one priority label per issue.

**Design:** Not required — subskill boundaries are specified in the issue. No novel architecture decisions.

**Plan:** Not required — five files in sequence, no complex dependencies, no multi-step execution ordering.

---

## §Skills

**Tier 1a — CDD authority:**
- `CDD.md` v3.15.0
- `src/packages/cnos.cdd/skills/cdd/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`

**Tier 1b — CDD lifecycle:**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (interpreting issue ACs and quality)

**Tier 2 — Writing bundle:**
- Applied: write/SKILL.md constraints — one governing question per artifact, no repetition, front-load the point

**Tier 3 — Issue-specific:**
- `src/packages/cnos.core/skills/write/SKILL.md`
- `src/packages/cnos.core/skills/design/SKILL.md`
- `src/packages/cnos.core/skills/skill/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` (loaded per issue body §Skills to load)

---

## §ACs

### AC1: Root issue skill becomes an orchestrator

**Evidence:** `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` is now 200 lines (was ~786 lines). It no longer contains the full label enum, full AC/proof-plan rules, full constraint-strata/exception/path-resolution rules, or any kata cases inline.

It contains: core principle, tri formula, when to load each subskill, minimal output pattern, closure rule, handoff checklist, external kata reference.

**Positive:** Root skill explains when to load `issue/labels`, `issue/contract`, `issue/proof`, `issue/constraints` in a dedicated "When to load each subskill" section.

**Negative check:** grep confirms the label enum (`bug`, `enhancement`, `docs`, `design`, `test`, `refactor`, `chore`, `tracking`) does NOT appear as a taxonomy table in root — only in `labels/SKILL.md`.

```
grep -c "| \`bug\`" src/packages/cnos.cdd/skills/cdd/issue/SKILL.md → 0
grep -c "| \`bug\`" src/packages/cnos.cdd/skills/cdd/issue/labels/SKILL.md → 1
```

✅ PASS

---

### AC2: Label taxonomy subskill exists

**Evidence:** `src/packages/cnos.cdd/skills/cdd/issue/labels/SKILL.md` exists.

**Contents verified:** kind label (8 values), priority label (P0–P3), optional surface/domain labels, examples, label discipline rule.

✅ PASS

---

### AC3: Required kind-label enum is defined

**Evidence:** `issue/labels/SKILL.md` defines all 8 required kind labels with definitions: `bug`, `enhancement`, `docs`, `design`, `test`, `refactor`, `chore`, `tracking`.

✅ PASS

---

### AC4: Required priority-label enum is defined

**Evidence:** `issue/labels/SKILL.md` defines P0–P3 with meanings. Rule states "Exactly one priority label per issue."

✅ PASS

---

### AC5: Surface/domain labels are optional and minimal

**Evidence:** `issue/labels/SKILL.md` §Surface/domain labels defines 9 optional labels and the rule that surface labels must not duplicate kind labels, with ❌/✅ examples.

✅ PASS

---

### AC6: Existing issue-contract sections are preserved

**Evidence:** After the split, the issue skill family still teaches all required sections:

| Section | Location |
|---|---|
| problem | `issue/contract/SKILL.md` §2.1 |
| impact | `issue/contract/SKILL.md` §2.2 |
| status truth | `issue/contract/SKILL.md` §2.3 |
| source of truth | `issue/contract/SKILL.md` §2.4 |
| scope | `issue/contract/SKILL.md` §2.5 |
| non-goals | `issue/contract/SKILL.md` §2.6 |
| constraint strata | `issue/constraints/SKILL.md` §2.1 |
| acceptance criteria | `issue/proof/SKILL.md` §2.1 |
| proof plan | `issue/proof/SKILL.md` §2.2 |
| skills to load | `issue/contract/SKILL.md` §2.7 |
| active design constraints | `issue/contract/SKILL.md` §2.8 |
| related artifacts | `issue/contract/SKILL.md` §2.9 |
| non-goals | `issue/contract/SKILL.md` §2.6 |
| success / closure condition | `issue/contract/SKILL.md` §2.11 + root §Closure rule |

No section disappeared from the family. ✅ PASS

---

### AC7: Detailed rules move to focused subskills

**Evidence:** Subskill ownership:

| Subskill | Owns |
|---|---|
| `issue/labels` | label taxonomy and examples |
| `issue/contract` | problem, impact, status truth, source map, scope, non-goals, skills, constraints, related artifacts, implementation guidance, success/closure |
| `issue/proof` | AC shape, invariant/oracle/positive/negative/surface, proof plan |
| `issue/constraints` | constraint strata, exceptions ledger, path resolution, cross-surface projections |

No rule is duplicated across subskills. ✅ PASS

---

### AC8: Root skill calls subskills

**Evidence:** Root `issue/SKILL.md` frontmatter declares:

```yaml
calls:
  - issue/labels/SKILL.md
  - issue/contract/SKILL.md
  - issue/proof/SKILL.md
  - issue/constraints/SKILL.md
```

Body has "When to load each subskill" section naming all four. ✅ PASS

---

### AC9: New subskills have valid frontmatter

**Evidence:** Each new subskill declares all 10 required fields. Verified by reading each file:

| Field | labels | contract | proof | constraints |
|---|---|---|---|---|
| name | ✅ | ✅ | ✅ | ✅ |
| description | ✅ | ✅ | ✅ | ✅ |
| artifact_class | ✅ (reference) | ✅ (skill) | ✅ (skill) | ✅ (skill) |
| kata_surface | ✅ (none) | ✅ (embedded) | ✅ (embedded) | ✅ (embedded) |
| governing_question | ✅ | ✅ | ✅ | ✅ |
| visibility | ✅ | ✅ | ✅ | ✅ |
| triggers | ✅ | ✅ | ✅ | ✅ |
| scope | ✅ | ✅ | ✅ | ✅ |
| inputs | ✅ | ✅ | ✅ | ✅ |
| outputs | ✅ | ✅ | ✅ | ✅ |

✅ PASS

---

### AC10: Examples reflect label rules

**Evidence:** `issue/labels/SKILL.md` contains 6 examples demonstrating correct label selection:

```
docs: README activation step       → docs, P1
feat(cli): cn activate             → enhancement, P1, cli
package(cdw): writing package      → enhancement, P2
tracking umbrella                  → tracking, P2
skill(cdd/issue): split subskills  → refactor, P2, cdd
bug: cn doctor reports false green → bug, P0
```

All examples use exactly one kind + one priority label. ✅ PASS

---

### AC11: External kata remains available

**Evidence:** Root `issue/SKILL.md` frontmatter has `kata_surface: external` and `kata_ref: src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`.

`src/packages/cnos.cdd.kata/katas/M5-issue-authoring/` exists and contains: `baseline.prompt.md`, `cdd.prompt.md`, `kata.md`, `rubric.json`. ✅ PASS

---

### AC12: No semantic expansion beyond labels

**Evidence:** Diff adds exactly 5 files (4 subskills + root rewrite). No new issue-contract doctrine was introduced except the label taxonomy. No GitHub label automation, CI checks, structured issue frontmatter, or new internal kata subskill was added. ✅ PASS
