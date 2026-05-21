---
name: epsilon
description: ε role in CDR. Iterates the CDR protocol itself — observes research-protocol gaps from the receipt stream, applies MCA discipline, writes the CDR-iteration artifact.
artifact_class: skill
parent: cdr
scope: role-local
governing_question: How does ε turn research-protocol-gap findings from the receipt stream into durable CDR-protocol improvements without accumulating undeclared overclaim risk?
triggers:
  - epsilon
  - cdr-iteration
  - protocol-iteration
---

# Epsilon (research ε)

> **This is the CDR-specific instantiation of the generic ε doctrine declared in
> [`ROLES.md §4b`](../../../../../../ROLES.md).** The kernel grammar —
> ε-as-protocol-reviewer in the scope ladder, the watched receipt fields
> (`protocol_gap_count`, `protocol_gap_refs`), the gap-class instantiation
> pattern (`{protocol}-{axis}-gap`), the iteration-artifact rule (required
> only when `protocol_gap_count > 0`), the MCA discipline (ship-now /
> next-MCA / no-patch), and the ε=δ collapse rule for small-protocol
> regimes — is inherited by reference. Only the **trigger classes**
> (research-failure classes per [`CDR.md`](../CDR.md) Field 5, not
> engineering-failure classes) and the **output-artifact location** (the
> CDR-iteration artifact under the project binding's wave surface) diverge
> for the research loss function per
> [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction).

## §1 ε's CDR-side scope

ε's domain in CDR is **protocol-iteration**: the work of observing whether the CDR protocol is itself coherent — whether it selects the right research gaps, closes them with correctly-calibrated typed receipts, and produces a research-claim transmission discipline that learns from its own waves rather than repeating the same overclaim class.

The canonical output artifact is the **CDR-iteration record** (the research analogue of `cdd-iteration.md`). Per the inherited cadence rule ([`ROLES.md §4b.4`](../../../../../../ROLES.md)), the artifact is required only when the wave's receipt has `protocol_gap_count > 0` — i.e. when the wave produced ≥1 finding tagged with one of the `cdr-*-gap` classes below. The receipt's `protocol_gap_count == 0` is the no-gap signal; no iteration file is required when the wave ran cleanly. The artifact's concrete location is project-binding-dependent (typically `<project>/.cdr/waves/{wave-id}/cdr-iteration.md` or an equivalent under the project's wave surface); the protocol-overlay layer declares the artifact's existence and shape, not its file path.

For waves that produce findings, the trigger classes the artifact enumerates are CDR-specific (per [`CDR.md §"Field 5"`](../CDR.md)):

- **`cdr-data-gate-gap`** — receipts shipping `observed` or `computed` claims with under-specified `data_refs` (no manifest, no checksum, no data-use compliance record). Signals the data-policy oracle (Field 2) is not being enforced.
- **`cdr-overclaim-gap`** — receipts whose `claim_status: observed` evidence supports only `inferred`. Signals β's claim/evidence-alignment oracle is too lenient or under-specified; the procedural skill (`cnos.cdr/skills/cdr/beta/SKILL.md`) needs a sharper calibration-determination procedure.
- **`cdr-reproduction-gap`** — receipts with `reproduction.output_match: false`, or receipts that should have a reproduction record (per Field 2 + per `cnos.cdr/skills/cdr/alpha/SKILL.md` §"Evidence-ref completeness rule") but lack one. Signals reproduction-from-clean is not running or is failing silently.
- **`cdr-citation-gap`** — receipts with claims that should cite external work but do not, or citations that do not support the claim they are invoked for. Signals the citation-integrity oracle is drifting.
- **`cdr-oracle-ambiguity-gap`** — the same review-oracle interpretation question appears across multiple waves (β cannot decide; γ has to adjudicate). Signals oracle under-specification at the protocol layer; `CDR.md` or the per-role skills need additional procedural detail.
- **`cdr-construct-drift-gap`** — a key research construct's definition shifts across receipts without explicit deprecation. Signals the project binding (cph or downstream) needs a construct-stability artifact (a glossary, a deprecation registry); the protocol-layer surface ε patches may declare a placeholder convention.

These are the research-failure trigger classes, distinct from engineering's `cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, `cdd-metric-gap`. The MCA-discipline shape (ship now / next-MCA / drop) is identical.

The MCA discipline governs what ε does with each finding (per the generic doctrine):

- **MCA available** (skill patch, gate, schema refinement) → ship it now as an immediate output in the same session.
- **No MCA yet, pattern real** → file a protocol MCI (entry in the project binding's iteration INDEX + an issue against `cnos.cdr` if the patch crosses the protocol boundary).
- **One-off, no pattern** → drop; note the drop explicitly.

ε cross-references [`ROLES.md §1`](../../../../../../ROLES.md#1-the-role-ladder) row 5 (ε iterates the δ-discipline) and [`CDR.md §"Field 5"`](../CDR.md) (the authoring procedure and per-finding shape for the CDR-iteration artifact).

## §2 ε's relationship to δ

ε and δ are often the same actor for research, just as for engineering. In small-protocol regimes — one active research operator running a handful of waves — protocol-iteration volume does not justify a dedicated reviewer of the protocol. The operator (δ) naturally accumulates the longitudinal view ε requires and performs ε work as part of wave close-out triage.

Separation becomes warranted when protocol-iteration volume justifies dedicated attention: when the CDR-iteration artifact is written on most waves with non-empty findings; when the iteration INDEX accumulates faster than one actor can triage; when the operator's operational load crowds out the reflective work. At that point, ε may be a distinct actor — a second agent or a dedicated human — who reads receipts across waves and drives protocol patches independently.

No claim is made here that ε is required as a separate human or agent. ε is a structural role in the scope-escalation ladder ([`ROLES.md §1`](../../../../../../ROLES.md#1-the-role-ladder)), not a headcount requirement. The role may collapse onto δ indefinitely in small-protocol regimes; the obligation is that ε work happens and is attributed, not that it is performed by a distinct person.

See [`ROLES.md §4`](../../../../../../ROLES.md#4-hats-vs-actors-roles-as-behavioral-contracts) for the general principle governing role collapse rules, and [`CDR.md §"Field 6"`](../CDR.md) for CDR's specific collapse rule — note that the α=β collapse is **never** safe for research-claim cycles (overclaim is the precise failure mode α cannot self-detect), while ε=δ and γ=δ are conditionally permitted.

## §3 Persona / protocol / project boundary

This overlay declares **what research ε does at the protocol layer**. It does not declare:

- **Who enacts ε** (persona — layer 1; lives in `<persona-hub>/spec/`).
- **What concrete iteration-INDEX path or wave directory ε writes to** (project — layer 4; lives in `<project>/.cdr/`). The trigger classes are protocol-bound (`cdr-*-gap` enumerated above); the storage path is project-bound.

If a project chooses additional trigger classes (e.g. `cdr-ethics-gap` for waves whose data raises ethics-review issues, or `cdr-stats-gap` for statistical-method drift), those are project-binding extensions; this overlay declares the minimum set per `CDR.md §"Field 5"`.
