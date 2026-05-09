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

## Mode declaration and MCA preconditions

Every issue declares a mode in its body header. The mode tells α how much design work has already converged and how much α must do before implementation begins.

| Mode | Meaning | When to use |
|---|---|---|
| **MCA** (Make-Coherent-Acceleration) | Design + plan already converged; α executes a known sequence | Both design and plan are committed at stable paths and will not be edited during the cycle |
| **explore** | Gap is real but boundary not yet bounded | The issue names a problem that needs investigation before implementation |
| **design-and-build** | Design needed before code; both happen in this cycle | Small enough that staging a separate design cycle would be overhead |
| **docs-only** | No code change, no version bump | Retroactive close-out, self-coherence report, doc cleanup, skill patches as the cycle's only output |

### MCA preconditions (all three must hold)

A cycle qualifies as MCA only when γ can answer **yes** to all of:

1. **Design committed** at a stable path (not in the issue body) — `docs/{tier}/{bundle}/{X.Y.Z}/DESIGN.md`, an architecture doc, or a converged proposal. The path must resolve and the document must be the source of truth, not a draft.
2. **Plan committed** with explicit step ordering — `docs/{tier}/{bundle}/{X.Y.Z}/PLAN.md` or equivalent, listing the sequenced work items α will execute. The plan names *what* α does in *what order*, not just *that* α will do something.
3. **Both stable** — neither will be edited during the cycle. If a step would require revising design or plan, the cycle is not MCA.

If any of the three fails, γ declares the cycle as `explore` or `design-and-build` and (where applicable) stages a separate design cycle first.

### MCA cycle citation

When mode = MCA, the issue's *Source of truth* table MUST cite the design path and the plan path:

```markdown
| Claim / surface | Canonical source | Status |
|---|---|---|
| Design | `docs/.../DESIGN.md` | Shipped |
| Plan | `docs/.../PLAN.md` | Shipped |
```

Without these rows, the mode is mis-labeled — γ either re-declares the mode or stages the missing artifact.

### Failure mode when mis-labeled MCA

α re-derives the plan from the issue body, which inflates context load and review pressure. The supercycle in `usurobor/tsc` master #23 produced direct evidence:

- Cycles where MCA discipline was honored (design + plan stable, cited) ran in **1–2 review rounds**
- Cycles where the design surface was ambiguous (no stable plan, α re-deriving) ran in **3 review rounds**

Mis-labeling MCA is not free — it costs α context and β review patience.

---

## Cycle scope sizing

γ sizes a cycle at scoping time — not after dispatch. The sizing decision appears explicitly in the issue body and is verified at handoff.

### Soft AC-count guideline

| Band | AC count | Signal |
|---|---|---|
| Small | 1–4 | Single-focus; artifact-class collapse optional (see `cdd/CDD.md` §1.2) |
| Typical | 5–7 | Standard cycle; monitor round-count target (≤2 code, ≤1 docs) |
| At-edge | 8–10 | Allowed; requires five-factor check and written justification if kept whole |
| Split-or-justify | ≥11 | Must split into master+subs or provide written justification |

The guideline is soft — γ's judgment governs. The heuristic forces the judgment to be visible.

### Five-factor split-decision heuristic

For any cycle at the at-edge (8–10) or split-or-justify (≥11) band, γ applies:

| Factor | Question | Splitting signal |
|---|---|---|
| (a) New code surface | How many new modules, packages, or runtime contracts does this introduce? | ≥2 new modules → split |
| (b) Cross-module breadth | How many existing modules are touched in non-trivial ways? | ≥3 modules touched → split |
| (c) Lifecycle span | Does the issue span multiple lifecycle phases (design → code → docs → infra) in ways that would serialize α's context? | Phases can't overlap → split |
| (d) MCA-precondition stability | Would mid-cycle design discovery force re-scoping? | Design not yet stable → explore or design-and-build cycle first |
| (e) Independent shippability of AC groups | Can any AC subset be shipped and used independently of the rest? | Yes → subs are the right shape |

A "yes" on any single factor is a signal; a "yes" on two or more factors is a strong split indicator.

### Master+subs pattern

When γ decides to split, the master+subs pattern applies:

1. **Independent shippability of subs** — each sub-issue is executable, verifiable, and shippable on its own. A chain where sub B blocks on sub A's merge is an anti-pattern; avoid it.
2. **Master cites subs** via the repo's native sub-issue mechanism (cnos: `sub-issue` field; cnos-tsc: `mcp__github__sub_issue_write`).
3. **Master closes** when all subs are closed, or any deferred sub is explicitly named as tracked debt in the master's closure condition.

### Empirical anchor

Evidence from `usurobor/tsc#23` (cnos-tsc supercycle, cycles 24–29):

- **Cycle 25** (12 ACs, 2 rounds): ran at the at-edge; α's 20-minute reading phase before first commit was operator-ambiguous ("stuck vs thinking"). Smaller issues converge faster and produce a clearer signal.
- **Cycle 24** (7 ACs, 3 rounds): same AC count as a "typical" cycle but touched 3+ modules with new runtime contracts. AC count alone did not predict round count — factors (a) and (b) were the load driver.
- **Cycle 26** (6 ACs, 1 round) vs **cycle 29** (6 ACs, 2 rounds): identical AC count, different difficulty. Factor (b) diverged — cycle 29 touched a cross-cutting surface that cycle 26 did not.

The five-factor heuristic encodes what AC count alone cannot capture.

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

## Cycle scope sizing

| Factor | Reading | Splitting signal? |
|---|---|---|
| (a) New code surface | | |
| (b) Cross-module breadth | | |
| (c) Lifecycle span | | |
| (d) MCA preconditions | | |
| (e) Independent shippability | | |

Decision: keep whole / split into master+subs — <one-line justification>

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
- [ ] Mode declared: MCA / explore / design-and-build / docs-only.
- [ ] If MCA: design and plan paths cited in source-of-truth table; both stable.
- [ ] Cycle scope-sizing decision recorded; if at-edge (8–10) or above (≥11), justification or master+subs pointer present.

---

## External kata

Practice and evaluation for this issue skill live in:

`src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`

That kata exercises status-truth split, source-of-truth alignment, constraint strata, AC negative space, proof-plan oracle, and non-goal discipline.
