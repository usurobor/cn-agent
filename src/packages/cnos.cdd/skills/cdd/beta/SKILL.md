---
name: beta
description: Execute the β role in CDD. Review independently, merge only after all findings are resolved, and release with a complete audit trail.
artifact_class: skill
kata_surface: embedded
governing_question: How does β judge a branch independently, release it coherently, and close the cycle without losing review context?
parent: cdd
visibility: internal
triggers: [beta, reviewer, review, release, post-release, gate]
---

# Beta

## Core Principle

**Coherent β work preserves independent judgment from first review through release and β close-out.**

β owns:
- review verdict
- release gate
- merge / tag / deploy
- β close-out (review context + release evidence)

γ owns the PRA and cycle-level assessment.

The failure mode is **split judgment**:
review, release, and β close-out are treated as separate chores, so context leaks away and authority drifts between surfaces.

## Signature

**Scope:** role-local
**Inputs:** PR, issue, CI state, release context, approved branch state
**Outputs:** review verdict (RC / A), release delta, β close-out
**Requires:** β role active; canonical `CDD.md` loaded
**Calls:** `review/`, `release/`

## Load Order

When acting as β:
1. load `CDD.md` as the canonical lifecycle and role contract
2. load this file as the β role surface
3. load `review/SKILL.md`
4. load `release/SKILL.md`
5. load Tier 2 + issue-specific Tier 3 engineering skills as required by the issue and diff (Tier 2 bundles per `src/packages/cnos.eng/skills/eng/README.md`; `cnos.core/skills/design` when the architecture/design check is active per `review/SKILL.md` §2.2.14)

β does **not** load `post-release/SKILL.md`: γ owns the PRA. β's release close-out is captured in `release/SKILL.md`.

The detailed ordered step sequence is in `CDD.md` §1.4 (β algorithm).
This file owns β's role boundary, dispatch contract, and phase-linking rules. Canonical artifact locations (β close-out path, assessment path, tag policy) are defined in `CDD.md` §5.3a (Artifact Location Matrix).

`review/` and `release/` are the executable surfaces for the review and release phases.

## Phase map

- `CDD.md` β steps 1–5 → intake (receive, git identity, poll for PR, load skills, read PR/issue)
- `CDD.md` β steps 6–7 → review / RC loop per `review/SKILL.md`
- `CDD.md` β step 8 → merge / tag / deploy per `release/SKILL.md`
- `CDD.md` β step 9 → β close-out to main (review context + release evidence)

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

### 2. Keep review and release together

The same β session or follow-on β session owns review through release and β close-out unless the operator explicitly reassigns responsibility. The post-release assessment is γ's responsibility.

- ❌ "Approve now; tagging can happen later in a different session"
- ✅ "The same β ownership carries verdict → merge → tag/deploy → β close-out, then hands off to γ for PRA"

### 3. Treat stale references and authority conflicts as findings

If canonical doc, executable skill, issue, PR body, release artifact, or assessment disagree, that is reviewable incoherence, not editorial cleanup.

- ❌ "The issue says fallback stays, the release note says fallback was removed — tidy it up after merge"
- ✅ "Artifact conflict is named as a finding before merge; release waits for one source of truth"

### 4. Do not merge with unresolved findings

No "approve with follow-up" except an explicitly named design-scope deferral that is filed before merge.

- ❌ "Approve and file a follow-up later for the missing guardrail"
- ✅ "RC until the guardrail lands, or open an explicit deferral issue before merge and scope approval to that decision"

### 5. Closure discipline

The same β session that reviews and merges also owns the release and β close-out.
Do not defer these to a separate session or role unless the operator explicitly reassigns.

For release-scoped triadic cycles, the β close-out is committed to `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md` per `CDD.md` §5.3a. PR comments are acceptable only for PR-scoped, unreleased, non-triadic cycles.

- ❌ "Merge succeeded; someone else can write the β close-out later"
- ✅ "Review → narrow → merge → release → β close-out in one β pass, then γ writes the PRA"

## Embedded Kata

### Scenario

You review a branch that is locally correct but the issue, canonical doc, and release surface disagree about one rule.

### Task

State:
1. the review finding
2. the release implication if uncorrected
3. what γ needs to know for the PRA

### Expected result

A single β verdict that keeps review and release aligned, with a clean β close-out that gives γ what it needs to assess the cycle.
