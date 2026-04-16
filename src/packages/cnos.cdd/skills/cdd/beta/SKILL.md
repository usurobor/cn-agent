---
name: beta
description: Execute the β role in CDD. Review independently, merge only after all findings are resolved, release with a complete audit trail, and own the post-release assessment.
artifact_class: skill
kata_surface: embedded
governing_question: How does β judge a branch independently, release it coherently, and close the cycle without losing review context?
parent: cdd
triggers: [beta, reviewer, review, release, post-release, gate]
---

# Beta

## Core Principle

**Coherent β work preserves independent judgment from first review through release and assessment.**

β owns:

- review verdict
- release gate
- merge / tag / deploy
- post-release assessment
- cycle closeout on the β side

The failure mode is **split judgment**:
review, release, and assessment are treated as separate chores, so context leaks away and authority drifts between surfaces.

## Load Order

When acting as β:

1. load `CDD.md` as the canonical lifecycle and role contract
2. load this file as the β role surface
3. load `review/SKILL.md`
4. load `release/SKILL.md`
5. load `post-release/SKILL.md`
6. load Tier 2 + issue-specific Tier 3 engineering skills as required by the issue and diff

The detailed step sequence is in CDD.md §1.4 (β algorithm). This file owns β's role boundary and dispatch contract. `review/`, `release/`, and `post-release/` are the detailed executable surfaces for each phase.

## Algorithm

1. **Review** — read the issue and PR independently; produce a verdict per `review/SKILL.md`.
2. **Narrow** — if RC, wait for α's fix; re-review only the affected surfaces.
3. **Merge** — when approved, squash-merge the PR.
4. **Release** — tag, deploy, verify per `release/SKILL.md`. If tag push fails due to env constraints, commit all release artifacts to main and defer tag push to γ/operator.
5. **Assess** — write post-release assessment per `post-release/SKILL.md`.
6. **Close-out** — write β close-out to main directly. This is β's input to γ's cycle iteration decision.

## Role Rules

### 1. Keep β independent

β does not author the fix it judges.
If RC is requested, α performs the fix.

### 2. Keep review, release, and assessment together

The same β session or follow-on β session owns all three unless the operator explicitly reassigns responsibility.

### 3. Treat stale references and authority conflicts as findings

If canonical doc, executable skill, issue, PR body, release artifact, or assessment disagree, that is reviewable incoherence, not editorial cleanup.

### 4. Do not merge with unresolved findings

No "approve with follow-up" except an explicitly named design-scope deferral that is filed before merge.

### 5. Closure discipline

The same β session that reviews and merges also owns the release, assessment, and close-out. Do not defer these to a separate session or role unless the operator explicitly reassigns. This was the recurring failure mode in cycles 3.55.0–3.56.0; the fix is: review → narrow → merge → release → assess → close-out, all in one pass.

## Embedded Kata

### Scenario

You review a branch that is locally correct but the issue, canonical doc, and release surface disagree about one rule.

### Task

State:

1. the review finding
2. the release implication if uncorrected
3. the assessment implication if merged anyway

### Expected result

A single β verdict that keeps review, release, and assessment aligned instead of treating them as three disconnected passes.
