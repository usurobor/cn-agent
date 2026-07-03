---
name: issues/taxonomy
description: Issue label taxonomy pointer. Use when classifying an issue's kind/area/priority/effort/status/protocol/dispatch labels, or when the board needs to interpret them.
artifact_class: skill
kata_surface: none
governing_question: What label vocabulary does an open issue carry, and what does each label mean?
visibility: public
triggers:
  - taxonomy
  - kind label
  - area label
  - effort label
  - issue labels
scope: task-local
inputs:
  - raw GitHub issue labels
outputs:
  - a classified label set (kind, area, priority, effort, status, protocol, dispatch)
---

# Issue taxonomy

## Core Principle

The issue label taxonomy is defined once, canonically, at
[`docs/development/issues/TAXONOMY.md`](../../../../../docs/development/issues/TAXONOMY.md).
This skill **cites** that document; it does not fork or restate the label
vocabulary here. If this skill and `TAXONOMY.md` ever disagree, `TAXONOMY.md`
governs.

## What `cnos.issues` reads from the taxonomy

The board mapping (`skills/map/SKILL.md`) consumes exactly the labels
`TAXONOMY.md` defines:

- **exactly one** primary `kind/*`
- **one or more** `area/*`
- **one** priority (`P0`–`P3`, or `priority/deferred`)
- **at most one** `effort/*` (`S`|`M`|`L`|`XL`, dashboard weight 1/2/4/8;
  missing effort means *unestimated*, never treated as small)
- **one** `status:*` when actionable
- `protocol:*` and `dispatch:cell` only when the issue is dispatchable

`kind/tracking` issues are containers and are excluded from effort rollups
by default.

## Rule

- Do not introduce a second definition of any label's meaning inside
  `cnos.issues`. Extend `TAXONOMY.md` first, then let the board mapping
  consume the extension.
- For *how* to apply the taxonomy to a raw issue (classification,
  dispatch-eligibility, close/supersede/defer), see `../triage/SKILL.md`
  and [`docs/development/issues/TRIAGE.md`](../../../../../docs/development/issues/TRIAGE.md) — a separate concern from
  label *definition*.
