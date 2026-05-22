---
name: epsilon
description: ε role in CDD — the CDD-specific instantiation of the generic
  ε doctrine. Reads receipt streams across cycles, surfaces `cdd-*-gap`
  findings, applies MCA discipline, and writes `cdd-iteration.md` when the
  cycle's receipt has `protocol_gap_count > 0`.
artifact_class: skill
parent: cdd
scope: role-local
governing_question: How does ε turn `cdd-*-gap` findings from the receipt
  stream into durable cdd-protocol improvements without accumulating deferred
  debt?
triggers:
  - cdd-iteration
  - protocol-iteration
  - epsilon
---

# Epsilon (CDD ε)

> **This is the CDD-specific instantiation of the generic ε doctrine declared in
> [`ROLES.md §4b`](../../../../../../ROLES.md).** The kernel grammar — what ε
> observes (receipt streams across cells); the watched receipt fields
> (`protocol_gap_count`, `protocol_gap_refs`); the gap-class instantiation
> pattern (`{protocol}-{axis}-gap`); the iteration artifact rule (required only
> when `protocol_gap_count > 0`); the MCA discipline (ship-now / next-MCA /
> no-patch); the ε=δ collapse rule for small-protocol regimes — is inherited
> by reference. Only **CDD-specific gap class names**, the **receipt-stream
> location**, the **iteration artifact path**, and the **aggregator location**
> diverge for the engineering loss function per
> [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction).

## §1 ε's CDD-side scope

ε's domain in cdd is **protocol-iteration**: the work of observing whether the
cdd protocol is itself coherent — whether it selects the right gaps, closes
them durably, and produces a system that learns from its own cycles rather than
repeating the same class of error.

**CDD-specific gap classes** (per the `{protocol}-{axis}-gap` pattern in
[`ROLES.md §4b.3`](../../../../../../ROLES.md)):

- **`cdd-skill-gap`** — a procedural skill (the active skill set loaded during
  the cycle) was underspecified or wrong for the failure pattern that surfaced.
  Patch lands in the affected skill file under
  `src/packages/cnos.cdd/skills/cdd/`.
- **`cdd-protocol-gap`** — CDD doctrine itself (CDD.md, COHERENCE-CELL.md,
  RECEIPT-VALIDATION.md, ROLES.md) drifted from its intended invariant, or two
  surfaces disagree about the same rule. Patch lands in the canonical doctrine
  file.
- **`cdd-tooling-gap`** — tooling absent, wrong, or environmentally unavailable
  blocked or distorted the cycle (the agent had to work around it). Patch
  lands in `scripts/`, CI workflows, or `cnos.eng/` tooling packages.
- **`cdd-metric-gap`** — a measurement the protocol expects is missing or wrong
  (e.g. mechanical-ratio threshold mis-calibrated, a CHANGELOG TSC field
  computed incorrectly). Patch lands in the measurement-emitting surface.

These are the CDD-specific trigger classes, distinct from cdr's
research-failure classes (`cdr-data-gate-gap`, `cdr-overclaim-gap`, etc. per
[`cnos.cdr/skills/cdr/epsilon/SKILL.md §1`](../../../../cnos.cdr/skills/cdr/epsilon/SKILL.md)).
The MCA-discipline shape (ship now / next-MCA / no-patch) is identical and
inherited from [`ROLES.md §4b.5`](../../../../../../ROLES.md).

**Receipt stream ε reads:** `.cdd/releases/{X.Y.Z}/{N}/receipt.yaml` (after
release) and `.cdd/unreleased/{N}/receipt.yaml` (during cycle). ε observes
`protocol_gap_count` + `protocol_gap_refs` across cycles to surface patterns
no single cycle could surface.

**Iteration artifact:** `cdd-iteration.md`, located at
`.cdd/unreleased/{N}/cdd-iteration.md` during the cycle and migrated to
`.cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md` at release time (per
`release/SKILL.md §2.5a`). For docs-only releases, the path uses the §2.5b
form: `.cdd/releases/docs/{ISO-date}/{N}/cdd-iteration.md`. The per-finding
shape is canonical at
[`cnos.handoff/skills/handoff/receipt-stream/SKILL.md §1`](../../../../cnos.handoff/skills/handoff/receipt-stream/SKILL.md);
the cdd-side runbook applies the rule at cycle close-out per
[`post-release/SKILL.md §5.6b`](../post-release/SKILL.md) (pointer).

**Cadence rule:** `cdd-iteration.md` is **required only when the cycle's
receipt has `protocol_gap_count > 0`** — i.e. when the close-out triage
produced ≥1 finding tagged `cdd-skill-gap`, `cdd-protocol-gap`,
`cdd-tooling-gap`, or `cdd-metric-gap`. The receipt's
`protocol_gap_count == 0` is the no-gap signal; no iteration file is required
when the cycle ran cleanly. This is the inherited rule from
[`ROLES.md §4b.4`](../../../../../../ROLES.md) — the receipt is the
always-present record of *whether* an iteration file is required; the file is
the conditional record of *what was found*. The canonical wire-format home
for the cadence rule + per-finding shape + aggregator-update procedure is
[`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`](../../../../cnos.handoff/skills/handoff/receipt-stream/SKILL.md);
[`post-release/SKILL.md §5.6b`](../post-release/SKILL.md) is the cdd-side
runbook (now a pointer to the canonical home).

Backward compatibility: existing empty-findings `cdd-iteration.md` files
(written under the prior every-cycle rule) remain valid artifacts. They are
no longer *required* under the new rule but are not invalid; future cycles
with `protocol_gap_count == 0` may simply omit the file.

**Aggregator:** `.cdd/iterations/INDEX.md` carries one row per cycle that
produced an iteration artifact. The row format (eight pipe-separated columns:
Cycle / Issue / Date / Findings / Patches / MCAs / No-patch / Path) and the
per-row update procedure are canonical at
[`cnos.handoff/skills/handoff/receipt-stream/SKILL.md §2`](../../../../cnos.handoff/skills/handoff/receipt-stream/SKILL.md).

**Cross-repo trace:** when a finding's disposition is `patch-landed` against
a different repo, ε also writes a cross-repo bundle at
`.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` per
[`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md).

## §2 ε's relationship to δ in CDD

The generic ε=δ collapse rule is declared in
[`ROLES.md §4b.6`](../../../../../../ROLES.md); this section names the
CDD-specific operating point.

In the current cnos cdd operating point — one active operator running a
handful of cycles per day across a small number of repos — protocol-iteration
volume does not justify a dedicated reviewer of the protocol. The operator
(δ) naturally accumulates the longitudinal view ε requires and performs ε
work as part of γ's Phase 4 (CDD iteration) and Phase 3 close-out triage.
This is the **ε=δ collapsed** mode, and it is the default for cdd at current
scale.

Separation becomes warranted when:

- `cdd-iteration.md` is written with non-empty findings on most cycles, or
- `.cdd/iterations/INDEX.md` accumulates faster than one actor can triage, or
- the operator's operational load (δ) crowds out the reflective work (ε).

At that point, ε may be a distinct actor — a second agent or a dedicated
human — who reads close-outs across cycles and drives protocol patches
independently.

The role-collapse is one of the *safe* collapses
([`ROLES.md §4`](../../../../../../ROLES.md)); contrast α=β (never safe).

## §3 Cross-references

- [`ROLES.md §4b`](../../../../../../ROLES.md) — generic ε doctrine; the
  authoritative source for the watched-fields invariant, the MCA discipline,
  and the ε=δ collapse rule.
- [`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`](../../../../cnos.handoff/skills/handoff/receipt-stream/SKILL.md) —
  canonical wire-format for the per-finding shape, INDEX.md aggregator-update
  procedure, cadence rule, finding-disposition vocabulary, and cross-repo
  trace bundle invariant. Migrated from `post-release/SKILL.md §5.6b` in Sub 5
  of [cnos#404](https://github.com/usurobor/cnos/issues/404).
- [`post-release/SKILL.md §5.6b`](../post-release/SKILL.md) — cdd-side
  runbook (pointer); applies the receipt-stream rule at cycle close-out.
- [`activation/SKILL.md §22`](../activation/SKILL.md) — per-repo cadence
  declaration, severity scale (D/C/B/A + `info`), auto-spawn MCA trigger.
- [`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md) — cross-repo iteration
  bundle structure.
- [`COHERENCE-CELL.md`](../COHERENCE-CELL.md) §"ε as Protocol Evolution" and
  §"ε Artifact Rule" — the cell-doctrinal grounding for ε's role.
- [`../../../../../../schemas/cdd/receipt.cue`](../../../../../../schemas/cdd/receipt.cue)
  `#Receipt` — the typed receipt that carries `protocol_gap_count` and
  `protocol_gap_refs`.
