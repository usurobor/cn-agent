---
name: review/architecture
description: Architecture and design check — verify reason-to-change, policy-above-detail, interface truthfulness, registry normalization, source/artifact/installed boundary, surface separation, and degraded-path visibility.
artifact_class: skill
kata_surface: external
governing_question: How does β verify that the branch preserves architectural boundaries and does not smear runtime surfaces, violate policy placement, or hide degraded paths?
visibility: internal
parent: review
triggers:
  - review architecture
  - review design
scope: task-local
inputs:
  - branch
  - diff
  - diff/context findings from Phase 2b
outputs:
  - architecture check table (7 questions A–G with yes/no/n/a results)
  - findings with severity and type
requires:
  - Phase 1 (contract integrity) completed
  - Phase 2b (diff/context inspection) completed
  - review orchestrator loaded
calls: []
---

# Review — Phase 2c: Architecture and Design Check

**PRE-GATE: Phase 2b (diff/context inspection) must be completed.**

**Activate when the change touches:** package boundaries, command/provider/orchestrator/skill separation, source/artifact/installed flow, registry design, kernel vs package responsibility, transport vs protocol semantics, command dispatch vs domain logic.

**Load:** `src/packages/cnos.core/skills/design/SKILL.md` when active.

---

## §2.2 Architecture and design check

| Check | Question |
|---|---|
| A. Reason to change | Does each touched module still have one real reason to change? |
| B. Policy above detail | Does policy remain in the kernel/core? |
| C. Truthful interface | Does each interface promise only what all implementations support? |
| D. Registry normalization | Do different source forms normalize into one runtime descriptor? |
| E. Source/artifact/installed | Is it still clear what is authored, built, and installed? |
| F. Surface separation | Are skills, commands, orchestrators, providers distinct? |
| G. Degraded-path visibility | Is fallback behavior visible and testable? |

Output block (include when active):

```markdown
## Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes / no / n/a | |
| Policy above detail preserved | yes / no / n/a | |
| Interfaces remain truthful | yes / no / n/a | |
| Registry model remains unified | yes / no / n/a | |
| Source/artifact/installed boundary preserved | yes / no / n/a | |
| Runtime surfaces remain distinct | yes / no / n/a | |
| Degraded paths visible and testable | yes / no / n/a | |
```

If any Architecture Check row is "no," the review cannot approve without an explicit issue-backed redesign or scope reduction.

Rules:
- Silent architectural boundary smear is a blocking finding.
- Silent source-of-truth duplication is a blocking finding.
- Intentional constraint revision must be named explicitly.

---

## After Phase 2c

Collect architecture findings. Return to orchestrator (`review/SKILL.md`) for verdict rules and output format.
