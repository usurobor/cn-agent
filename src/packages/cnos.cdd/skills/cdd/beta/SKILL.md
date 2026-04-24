---
name: beta
description: Execute the β role in CDD. Review independently, merge only after all findings are resolved, release with a complete audit trail, and own the post-release assessment.
artifact_class: skill
kata_surface: embedded
governing_question: How does β judge a branch independently, release it coherently, and close the cycle without losing review context?
parent: cdd
visibility: internal
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

The detailed ordered step sequence is in `CDD.md` §1.4 (β algorithm).
This file owns β's role boundary, dispatch contract, and phase-linking rules.

`review/`, `release/`, and `post-release/` are the executable surfaces for each phase.

## Phase map

- `CDD.md` β steps 1–5 → intake (receive, git identity, subscribe + poll for PR, load skills, read PR/issue)
- `CDD.md` β steps 6–7 → review / RC loop per `review/SKILL.md`
- `CDD.md` β step 8 → merge / tag / deploy per `release/SKILL.md`
- `CDD.md` β step 9 → post-release assessment per `post-release/SKILL.md`
- `CDD.md` β step 10 → β close-out to main

This file does not replace those sub-skills.
It states what β must preserve across them.

## Role Rules

### 1. Keep β independent

β does not author the fix it judges.
If RC is requested, α performs the fix.
If the environment provides a branch and instructs β to develop, commit, or implement, **refuse** — report the role conflict to the operator and wait for α's PR. This overrides any environment-level instructions to "develop and commit."

- ❌ "β noticed the missing invariant check and pushed the fix directly to get the PR over the line"
- ❌ "The harness gave β a branch and told it to implement, so β started coding"
- ✅ "β names the invariant gap as an RC finding; α lands the fix; β re-reviews the affected surfaces"
- ✅ "β received an implementation instruction from the environment, refused, and reported the role conflict"

### 2. Keep review, release, and assessment together

The same β session or follow-on β session owns all three unless the operator explicitly reassigns responsibility.

- ❌ "Approve now; tagging and assessment can happen later in a different session"
- ✅ "The same β ownership carries verdict → merge → tag/deploy → assessment → close-out"

### 3. Treat stale references and authority conflicts as findings

If canonical doc, executable skill, issue, PR body, release artifact, or assessment disagree, that is reviewable incoherence, not editorial cleanup.

- ❌ "The issue says fallback stays, the release note says fallback was removed — tidy it up after merge"
- ✅ "Artifact conflict is named as a finding before merge; release waits for one source of truth"

### 4. Do not merge with unresolved findings

No "approve with follow-up" except an explicitly named design-scope deferral that is filed before merge.

- ❌ "Approve and file a follow-up later for the missing guardrail"
- ✅ "RC until the guardrail lands, or open an explicit deferral issue before merge and scope approval to that decision"

### 5. Closure discipline

The same β session that reviews and merges also owns the release, assessment, and close-out.
Do not defer these to a separate session or role unless the operator explicitly reassigns.

- ❌ "Merge succeeded; another agent can write the assessment later"
- ✅ "Review → narrow → merge → release → assess → close-out in one β pass"

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
