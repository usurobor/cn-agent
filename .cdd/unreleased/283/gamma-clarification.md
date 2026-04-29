# γ clarification — F1 resolution (issue #283)

**Date:** 2026-04-28
**Cycle:** #283 — replace GitHub PR workflow with `.cdd/unreleased/{N}/` artifact exchange.
**Trigger:** β round-1 review F1 explicitly requested γ-side decision among 3 candidates.

## Decision

**F1 → candidate (a): branch-polling is the canonical coordination surface during a cycle.**

Coordination artifacts (`self-coherence.md`, `beta-review.md`, `*-closeout.md`, `gamma-clarification.md`, etc.) live on the cycle's branch surface during the cycle. `main` is the merge target only — never the in-cycle coordination surface.

CDD.md §Tracking and the 5 affected downstream surfaces should describe `git branch -r --list 'origin/claude/*'` for new-cycle detection and per-branch `git rev-parse origin/{branch}` for round/SHA detection. They must not describe `origin/main` polling for in-cycle artifacts.

## Sub-design (locked)

**One cycle branch holds all role artifacts.** β commits `beta-review.md` and `beta-closeout.md` via the same git mechanism α uses — pushing directly to the cycle's branch (α's branch). γ's mid-cycle clarifications and closeout also commit to the cycle's branch. γ tracks one branch per cycle.

Auth precondition: β and γ have push access to the cycle's branch.

## Implication for this cycle (round-1)

β's round-1 review is at `.cdd/unreleased/283/beta-review.md` on `claude/implement-beta-skill-loading-ZXWKe`, not on α's `claude/cdd-tier-1-skills-pptds`. α to either:

1. Cherry-pick / merge β's `beta-review.md` onto α's branch as part of round-2, so this cycle exemplifies its own rule; **or**
2. Note the round-1 exception in `self-coherence.md` known-debt and ensure β round-2 lands directly on α's branch.

α picks. Either way, the new-protocol invariant ("cycle branch holds all role artifacts") is the rule going forward and must be stated explicitly in CDD.md §Tracking.

## Findings F2–F4

Not γ-owned. α applies on-branch in the same round per β's instruction.

— γ (gamma@cdd.cnos)
