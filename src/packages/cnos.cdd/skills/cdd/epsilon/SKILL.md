---
name: epsilon
description: ε role in CDD. Iterates the cdd protocol itself: collects
  cdd-*-gap findings from cycle close-outs, applies MCA discipline, and
  writes cdd-iteration.md as the first-class artifact for protocol
  self-improvement.
artifact_class: skill
parent: cdd
scope: role-local
governing_question: How does ε turn cdd-*-gap findings from cycle close-outs
  into durable protocol improvements without accumulating deferred debt?
---

# Epsilon

## §1 ε's cdd-side scope

ε's domain in cdd is protocol-iteration: the work of observing whether the
cdd protocol is itself coherent — whether it selects the right gaps, closes
them durably, and produces a system that learns from its own cycles rather
than repeating the same class of error.

The canonical output artifact is `cdd-iteration.md` (see
`cdd/post-release/SKILL.md` Step 5.6b). ε writes `cdd-iteration.md` when
the cycle close-out triage produces at least one finding tagged
`cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, or
`cdd-metric-gap`. Empty cycles produce no file; the signal stays high only
when it carries signal.

The MCA discipline governs what ε does with each finding:
- **MCA available** (skill patch, gate, mechanization) → ship it now as
  an immediate output in the same session.
- **No MCA yet, pattern real** → file a protocol MCI. Two kinds:
  - *Protocol MCI* (future cdd cycles need to know) → entry in
    `.cdd/iterations/INDEX.md` + issue.
  - *Agent MCI* (future sessions of this agent need to know) → agent
    hub update.
- **One-off, no pattern** → drop; note the drop explicitly.

ε cross-references `ROLES.md §1` row 5 (ε iterates the δ-discipline) and
`cdd/post-release/SKILL.md` Step 5.6b (the authoring procedure and
per-finding shape for cdd-iteration.md).

## §2 ε's relationship to δ

ε and δ are often the same actor. In small-protocol regimes — one active
operator running a handful of cycles per week — protocol-iteration volume
does not justify a dedicated reviewer of the protocol. The operator (δ)
naturally accumulates the longitudinal view that ε requires and performs ε
work as part of γ's Phase 4 (CDD iteration) and Phase 3 close-out triage.

Separation becomes warranted when protocol-iteration volume justifies
dedicated attention: when `cdd-iteration.md` is written on most cycles,
when the INDEX accumulates faster than one actor can triage, or when the
operator's operational load (δ) crowds out the reflective work (ε). At
that point, ε may be a distinct actor — a second agent or a dedicated
human — who reads close-outs across cycles and drives protocol patches
independently.

No claim is made here that ε is required as a separate human or agent. ε
is a structural role in the scope-escalation ladder (`ROLES.md §1`), not
a headcount requirement. The role may collapse onto δ indefinitely in
small-protocol regimes; the obligation is that ε work happens and is
attributed, not that it is performed by a distinct person.
