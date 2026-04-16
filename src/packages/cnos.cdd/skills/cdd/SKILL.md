---
name: cdd
description: Coherence-Driven Development. Use for substantial changes that require explicit selection, artifact flow, review, release, and post-release closure.
artifact_class: skill
kata_surface: embedded
governing_question: How do we evolve a system through a substantial change without losing coherence across selection, implementation, review, release, and closure?
triggers: [review, PR, release, issue, design, plan, assess, post-release, ship, tag]
---

# CDD

## Load order

This skill is a loader entrypoint. When CDD applies:

1. Load `CDD.md` in this directory as the canonical algorithm
2. Load the role skill for the active role:
   - α: `alpha/SKILL.md`
   - β: `beta/SKILL.md`
   - γ: `gamma/SKILL.md`
3. Load lifecycle sub-skills as directed by the role and work shape
4. Load Tier 2 and Tier 3 skills as directed by the issue

## Rule

`CDD.md` is the only normative source for:

- roles (§1.4)
- lifecycle (§3)
- selection (§2)
- artifact contract (§5)
- skill loading tiers (§4.4)
- gates (§4)
- closeout (§10)
- cycle iteration (§9.1)
- dispatch prompt format (§1.4)

Role skills may add:
- evidence requirements
- role-local gates
- checklists
- examples and katas
- execution detail for the steps they own

Role skills may **not** redefine:
- selection rule order
- ordered lifecycle steps
- dispatch artifact contract
- closeout obligations already named in `CDD.md`

## Role skills

- `alpha/` — α role: implementation, self-coherence, pre-review gate
- `beta/` — β role: review, release, post-release assessment
- `gamma/` — γ role: issue quality, dispatch, unblocking, close-out triage, process iteration

## Lifecycle sub-skills

- `design/` — design artifact production
- `issue/` — issue creation and structure
- `plan/` — implementation planning
- `review/` — code review protocol
- `release/` — merge, tag, deploy
- `post-release/` — assessment and closure

## Conflict rule

If this file and `CDD.md` disagree, `CDD.md` governs.

If a role skill and `CDD.md` disagree on ordered steps, selection rules, or artifact contract, `CDD.md` governs.

If a role skill adds execution detail, gates, or evidence requirements for a `CDD.md` step **without** changing that step's order or meaning, the role skill governs for that local execution detail.
