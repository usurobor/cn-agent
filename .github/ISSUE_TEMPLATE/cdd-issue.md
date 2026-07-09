---
name: CDD Issue
about: Gap, bug, or feature tracked through CDD
labels: ''
---

## Gap

**What:** <!-- named incoherence — what is wrong or missing -->

**Why it matters:** <!-- consequence if left unfixed -->

**What fails if skipped:** <!-- concrete failure scenario -->

## Evidence

<!-- How was this discovered? Link to logs, traces, sessions, adhoc threads, or essays. -->

## Priority

<!-- P0 / P1 / P2 / P3 -->

## Classification

<!-- Label per docs/development/issues/TAXONOMY.md: one primary kind/*, one or more area/*,
     and (only if an executable cell) dispatch:cell + protocol:*. See docs/development/issues/TRIAGE.md. -->

- **kind:** <!-- kind/bugfix | kind/cleanup | kind/process | kind/feature | kind/tooling | kind/doctrine | kind/audit | kind/tracking | kind/research | kind/skill | kind/spike -->
- **area:** <!-- area/* (one or more) -->
- **dispatchable:** <!-- yes → dispatch:cell + protocol:cds/cdd + status:ready|todo ; no → design/tracking/research -->
  <!-- The status:* label alone gates dispatch readiness (cnos#640). Do NOT add a body
       sentence like "Not dispatched — status:ready ... dispatch on explicit operator
       authorization." — that duplicate hold-state prose is what cnos#614/#633 showed can
       drift out of sync with the label. Use `cn issues dispatch --issue N` to authorize
       dispatch; see dispatch-protocol/SKILL.md §1.2. -->

## Scope

- **Work shape:** <!-- substantial / small-change / immediate-output -->
- **Estimated level:** <!-- L5 / L6 / L7 -->
- **Non-goals:** <!-- explicit scope boundary — what this issue does NOT attempt -->

## Skills and Constraints

- **Tier 3 skills:** <!-- only issue-specific skills, not Tier 1 or Tier 2 -->
- **Active design constraints:** <!-- linked governing constraints (invariants, transition, process) in plain language -->

## Acceptance Criteria

- [ ] **AC1:**
- [ ] **AC2:**

## Related Artifacts

- **Design:** <!-- design doc or "none" -->
- **Plan:** <!-- plan doc or "none" -->
- **Related issues:** <!-- #N, #M or "none" -->
- **Dependency notes:** <!-- blockers or sequencing constraints, or "none" -->
