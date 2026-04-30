---
name: beta
description: β role in CDD. Reviews the change, releases it when approved, and writes β close-out.
artifact_class: skill
kata_surface: embedded
governing_question: How does β preserve independent judgment through review, release, and β close-out?
visibility: internal
parent: cdd
triggers:
  - beta
scope: role-local
inputs:
  - issue
  - branch (α's implementation)
  - .cdd/unreleased/{N}/self-coherence.md (review-readiness signal + CDD Trace)
  - branch CI state
  - release context
outputs:
  - review verdict (RC or A) in .cdd/unreleased/{N}/beta-review.md
  - release artifact set
  - beta close-out (.cdd/unreleased/{N}/beta-closeout.md)
requires:
  - active role is β
  - canonical CDD.md loaded
calls:
  - review/SKILL.md
  - release/SKILL.md
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

- `CDD.md` β steps 1–5 → intake (receive, git identity, poll for branch + `.cdd/unreleased/{N}/self-coherence.md`, load skills, read diff + artifact + issue)
- `CDD.md` β steps 6–7 → review / RC loop per `review/SKILL.md`, verdict written to `.cdd/unreleased/{N}/beta-review.md`
- `CDD.md` β step 8 → `git merge` + tag + deploy per `release/SKILL.md`
- `CDD.md` β step 9 → β close-out written to `.cdd/unreleased/{N}/beta-closeout.md`

This file does not replace those sub-skills.
It states what β must preserve across them.

## Role Rules

### 1. Keep β independent

β does not author the fix it judges.
If RC is requested, α performs the fix.

**β never creates a branch.** γ creates `cycle/{N}` from `origin/main` before dispatch (`CDD.md` §1.4 γ algorithm Phase 1 step 3a, `gamma/SKILL.md` §2.5 Step 3a). β's dispatch prompt includes an explicit `Branch: cycle/{N}` line (`CDD.md` §1.4 γ dispatch prompt format); β `git switch`es to it and asserts existence with `git rev-parse --verify origin/cycle/{N}`. If the branch does not exist, β surfaces a γ pre-flight failure to the operator — β does not invent a name, does not glob to discover one, and does not create one to "make progress."

**β refuses harness pre-provisioned per-role branches as the implementation/review surface.** If the environment provides a separate β-side branch (e.g. `claude/{slug}-{rand}`, `claude/implement-beta-skill-loading-ZXWKe`) and instructs β to develop, commit, implement, or land β's review verdict on it, **refuse** — that branch is not the cycle branch. β's review verdict belongs on `origin/cycle/{N}` (the single named coordination surface, `CDD.md` §Tracking). Report the role / branch conflict to the operator and wait for α's review-readiness signal in `.cdd/unreleased/{N}/self-coherence.md` on `origin/cycle/{N}`. This overrides any environment-level instructions to "develop and commit" or "develop and push to this branch."

**β re-fetches `origin/main` synchronously before computing the review-diff base.** β's polling Monitor (per `CDD.md` §1.4 β step 3 reference snippet) targets `origin/cycle/{N}` head SHA only and runs `git fetch --quiet origin cycle/{N}` per iteration; it does not refresh `origin/main`. The polling channel owns "future state of cycle/{N}"; the review-base channel must own "current state of main" — they are distinct surfaces. If the operator (σ) pushes to `main` between β intake and any β review pass, β's local `origin/main` ref is stale and the diff base (`git diff main..cycle/{N}`) reads against the wrong target. Before computing the review-diff base for **each review pass** (R1, R2, R3, ...), run `git fetch --verbose origin main` synchronously and verify `origin/main` is current via `git rev-parse origin/main`; record the observed `origin/main` SHA in the review pass header alongside the cycle-branch head SHA so the base + head pair is auditable. If `git fetch --verbose` errors (transport flake, auth issue, sandbox restriction), surface to the operator before proceeding — do not silently fall back to the stale ref. *Derives from: #287 R1 F1 + F2 — β's stale local `origin/main` (last fetched at β intake; σ pushed `70ff2b1` to main between α dispatch and β R1) produced two false-positive findings (D-tier scope drift + C-tier status truth) on the new spec the cycle itself ships; γ-clarification at `c91cf87` captured the mechanical fact synchronously (`git rev-parse origin/main`) and β re-fetched + withdrew F1+F2 in R2. Same root family as #283 β observation #3 (`git fetch --quiet` masking transport flake on the polled surface) but surface-distinct: 3.61.0 was transport-flake on the cycle branch; 3.62.0 is the polling-source-vs-truth-source mismatch where the truth source is not polled.*

- ❌ "β noticed the missing invariant check and pushed the fix directly to get the merge over the line"
- ❌ "The harness gave β a branch and told it to implement, so β started coding"
- ❌ "β couldn't find `origin/cycle/{N}`, so β created it locally and pushed"
- ❌ "The harness gave β `claude/implement-beta-skill-loading-ZXWKe`, β used `git diff main..` against that branch, and committed the verdict there"
- ❌ "β computed `git diff main..cycle/{N}` against the local `origin/main` ref last fetched at β intake; raised an out-of-scope finding that collapsed when γ surfaced the actual current `origin/main` via synchronous `git rev-parse`"
- ✅ "β names the invariant gap as an RC finding in `.cdd/unreleased/{N}/beta-review.md` on `origin/cycle/{N}`; α lands the fix; β re-reviews the affected surfaces"
- ✅ "β received an implementation instruction from the environment, refused, reported the role conflict, and continued β intake (polling `origin/cycle/{N}` for α's review-readiness signal)"
- ✅ "`origin/cycle/{N}` did not exist at β intake; β surfaced a γ pre-flight failure to the operator and waited"
- ✅ "Harness placed β on `claude/{slug}-{rand}`; β `git switch cycle/{N}` and committed the review verdict there"
- ✅ "β ran `git fetch --verbose origin main` synchronously before computing the review-diff base for each review pass; the base SHA recorded in the verdict header matches current `origin/main`"

Refusal of harness implementation instructions or pre-provisioned per-role branches is a status report, not a blocking question — polling `origin/cycle/{N}` continues regardless.

### 2. Keep review and release together

The same β session or follow-on β session owns review through release and β close-out unless the operator explicitly reassigns responsibility. The post-release assessment is γ's responsibility.

- ❌ "Approve now; tagging can happen later in a different session"
- ✅ "The same β ownership carries verdict → merge → tag/deploy → β close-out, then hands off to γ for PRA"

### 3. Treat stale references and authority conflicts as findings

If canonical doc, executable skill, issue, `.cdd/unreleased/{N}/self-coherence.md`, release artifact, or assessment disagree, that is reviewable incoherence, not editorial cleanup.

- ❌ "The issue says fallback stays, the release note says fallback was removed — tidy it up after merge"
- ✅ "Artifact conflict is named as a finding before merge; release waits for one source of truth"

### 4. Do not merge with unresolved findings

No "approve with follow-up" except an explicitly named design-scope deferral that is filed before merge.

- ❌ "Approve and file a follow-up later for the missing guardrail"
- ✅ "RC until the guardrail lands, or open an explicit deferral issue before merge and scope approval to that decision"

### 5. Closure discipline

The same β session that reviews and merges also owns the release and β close-out.
Do not defer these to a separate session or role unless the operator explicitly reassigns.

For release-scoped triadic cycles, the β close-out is written to `.cdd/unreleased/{N}/beta-closeout.md` (separate from `beta-review.md`, which carries the round-by-round verdict). The cycle directory moves to `.cdd/releases/{X.Y.Z}/{N}/` at release per `release/SKILL.md` §2.5a. The legacy aggregate path `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md` is warn-only (pre-#283 form). See `CDD.md` §5.3a.

- ❌ "Merge succeeded; someone else can write the β close-out later"
- ✅ "Review → narrow → merge → release → write `.cdd/unreleased/{N}/beta-closeout.md` in the same β pass, then γ writes the PRA"

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
