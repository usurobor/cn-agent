---
name: issues/triage
description: Issue triage doctrine pointer. Use when classifying a new issue, deciding dispatch eligibility, or deciding whether to close/supersede/defer.
artifact_class: skill
kata_surface: none
governing_question: How does a raw issue become a coherently labeled, correctly dispatched (or deferred, or closed) item?
visibility: public
triggers:
  - triage
  - classify issue
  - dispatch eligibility
  - close issue
  - supersede issue
scope: task-local
inputs:
  - a raw or partially labeled issue
outputs:
  - a triage decision (label set, dispatch eligibility, close/supersede/defer disposition)
---

# Issue triage

## Core Principle

Triage operational rules are defined once, canonically, at
[`docs/development/issues/TRIAGE.md`](../../../../../docs/development/issues/TRIAGE.md).
This skill **cites** that document; it does not fork or restate the
operational rules here. If this skill and `TRIAGE.md` ever disagree,
`TRIAGE.md` governs.

`TRIAGE.md` is the *how-to-apply* companion to `../taxonomy/SKILL.md` (the
*what-the-labels-mean* reference) — triage decides which labels a given
issue gets and whether it is dispatchable; taxonomy defines what each label
means.

## What `cnos.issues` needs from triage

Board mapping (`skills/map/SKILL.md`) assumes issues arriving at
`cn issues map` are triaged well enough to carry a coherent label set. The
board does not re-triage issues — it renders whatever labels are present,
flagging gaps (unlabeled kind, unestimated effort) rather than guessing.

## Rule

- Do not add board-specific triage rules inside `cnos.issues`. If the board
  needs a new triage rule, propose it as a `TRIAGE.md` change, not a
  package-local fork.
- This package does not alter dispatch/wake/CDD/label behavior — it only
  reads the taxonomy and triage doctrine to render a board.
