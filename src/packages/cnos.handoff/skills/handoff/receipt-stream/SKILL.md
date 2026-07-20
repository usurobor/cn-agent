---
name: receipt-stream
description: Cross-cycle receipt-stream doctrine for protocol-iteration. Use when γ performs post-disconnect terminal close and writes the per-cycle iteration artifact (`cdd-iteration.md`), when γ updates the cross-cycle aggregator at `.cdd/iterations/INDEX.md`, when ε reads the receipt stream across cycles to surface protocol-gap patterns, when a finding's disposition lands a cross-repo patch and the cross-repo trace bundle attaches, or when a courtesy empty-findings stub is written for a `protocol_gap_count == 0` cycle.
artifact_class: skill
kata_surface: embedded
governing_question: How does the protocol carry its own learning signal across cycles — through a per-cycle iteration artifact (with a fixed per-finding shape) and a cross-cycle aggregator (with a fixed row format) — so that ε can read the receipt stream and surface protocol-gap patterns no single cycle can surface?
visibility: public
parent: handoff
triggers:
  - receipt-stream
  - cdd-iteration
  - INDEX.md
  - aggregator
  - protocol_gap_count
  - courtesy stub
  - per-finding shape
  - ε feed
scope: task-local
inputs:
  - closing cycle with γ-closeout filed
  - cycle's receipt with `protocol_gap_count` populated
  - α / β / γ close-out artifacts on main
outputs:
  - per-cycle iteration artifact at `.cdd/unreleased/{N}/cdd-iteration.md` (or courtesy stub)
  - one new row in `.cdd/iterations/INDEX.md`
  - cross-repo trace bundle at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` when a finding lands a cross-repo patch
requires:
  - γ close-out is filed
  - cycle's receipt is final (`protocol_gap_count`, `protocol_gap_refs` populated)
calls:
  - HANDOFF.md
  - artifact-channel/SKILL.md
  - cross-repo/SKILL.md
---

# Receipt-Stream — cross-cycle ε feed (per-cycle iteration artifact + cross-cycle aggregator)

## Core principle

**The receipt-stream is the protocol's own learning signal across cycles. Each cycle that surfaces a protocol gap (`protocol_gap_count > 0`) emits a per-cycle iteration artifact (`cdd-iteration.md`) with a fixed per-finding shape; γ then appends one row to the cross-cycle aggregator (`.cdd/iterations/INDEX.md`). ε reads the aggregator to surface patterns across cycles; the per-cycle artifact carries the structured findings the aggregator points at. Cycles with `protocol_gap_count == 0` may write a courtesy empty-findings stub (per cycle/401 convention) but are not required to.**

The per-cycle artifact is **role-local**: γ authors it as part of the cycle's close-out (with ε review where ε is a distinct actor; ε=δ collapsed in the current cnos cdd operating point per `cnos.cdd/skills/cdd/epsilon/SKILL.md §2`). The cross-cycle aggregator is **persistent**: `.cdd/iterations/INDEX.md` accumulates one row per cycle that produced an iteration artifact. The two together form the receipt-stream surface that ε observes — the per-cycle structured findings *plus* the cross-cycle row-by-row index that lets ε scan many cycles at once.

The failure mode is **receipt-stream drift** — γ writes the iteration artifact but forgets the aggregator row (ε cannot find the cycle when scanning); γ writes the aggregator row but the per-finding shape diverges from the canonical fields (ε's pattern-detection across cycles breaks); a cross-repo `patch-landed` finding lands without the bundle trace (the patch's lineage is lost when the source-side bundle archives); the courtesy stub gets misused as "iteration complete" when no actual triage happened (survivorship bias re-emerges). The protocol closes each failure mode by pinning the per-finding shape, the row format, the cadence rule, and the cross-repo bundle invariant explicitly.

This skill is the **canonical surface** for the receipt-stream wire-format invariants: the per-finding shape inside `cdd-iteration.md`; the INDEX.md row format; the cadence rule (required when `protocol_gap_count > 0`; courtesy stub when 0); the finding-disposition vocabulary (`patch-landed` / `next-MCA` / `no-patch`); the cross-repo trace bundle composition. Consumer protocols cite this skill for the wire-format; the role-local authoring procedure (which session writes, what review pass applies) lives in the consumer (`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` is now a pointer; `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` is ε's role-local authority on classification).

## Authority

- This skill (in `cnos.handoff`) is the only canonical home for the receipt-stream wire-format invariants: the per-finding shape (`Source` / `Class` / `Trigger` / `Description` / `Root cause` / `Disposition` + per-disposition sub-fields); the INDEX.md row format (eight pipe-separated columns: `Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path`); the cadence rule; the finding-disposition vocabulary; the cross-repo trace bundle invariant. If another skill or doctrine surface carries the same wire-format claims, this skill governs.
- The consumer protocol packages carry the role-local authoring procedure. `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` is a pointer to this skill for the wire-format and applies the rule at cycle close-out (γ writes the artifact per this skill's shape; γ appends the INDEX row). `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` is the role-local authority on the CDD-specific gap-class instantiation (`cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap`); the gap-class names are CDD-side. CDR has its own gap-class set per `cnos.cdr/skills/cdr/epsilon/SKILL.md §1` (`cdr-data-gate-gap` / `cdr-overclaim-gap` / etc.); the per-finding shape and aggregator row format are protocol-agnostic.
- The cross-repo trace bundle composes with `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. When a finding's disposition is `patch-landed` with a `Cross-repo` target, γ creates a bundle at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` per the cross-repo skill's case-(b) / case-(c) bundle file set. This skill names the trigger; cross-repo names the bundle composition.

## Scope

In scope:

- The per-finding shape inside `cdd-iteration.md` (F1, F2, … with `Source` / `Class` / `Trigger` / `Description` / `Root cause` / `Disposition` and per-disposition sub-fields: `Patch` + `Affects` + `Cross-repo` for `patch-landed`; `Issue filed` + `First AC` for `next-MCA`; `Reason` for `no-patch`).
- The cross-cycle aggregator (`.cdd/iterations/INDEX.md`) row format: eight pipe-separated columns naming the cycle / issue / date / findings count / patches count / MCAs count / no-patch count / iteration-file path.
- The cadence rule: per-cycle artifact required when `protocol_gap_count > 0`; courtesy empty-findings stub permitted (not required) when `protocol_gap_count == 0` (cycle/401 convention).
- The finding-disposition vocabulary: `patch-landed` | `next-MCA` | `no-patch` (no other dispositions are valid).
- The aggregator location: `.cdd/iterations/INDEX.md` at the repo root (the live, persistent aggregator; data, not doctrine — content stays put; doctrine about it lives here).
- The cross-repo trace bundle invariant: when a finding's disposition is `patch-landed` with `Cross-repo` target, γ creates a bundle at the canonical path per `cnos.handoff/skills/handoff/cross-repo/SKILL.md`.
- Empirical anchors: cycle/335 (aggregator initialized); cycle/401 (courtesy stub convention landed); every cycle since cnos#364 carries an INDEX row + a `cdd-iteration.md`.

Out of scope:

- The role-local authoring procedure (which session writes; the close-out gate ordering; the ε-vs-γ collapse rule). Lives in `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` (cdd-side runbook, now pointer to this skill) and `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` (ε's CDD-side role-local authority).
- The CDD-specific gap-class names (`cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap`). Lives in `cnos.cdd/skills/cdd/epsilon/SKILL.md §1`. CDR has different gap classes; future c-d-X may have other gap classes. The per-finding `Class` field is consumer-specific in vocabulary; the wire-format `Class:` slot is protocol-agnostic.
- The receipt schema itself (the typed `#Receipt` that carries `protocol_gap_count` and `protocol_gap_refs`). Lives in `schemas/cdd/receipt.cue`. This skill consumes the receipt's `protocol_gap_count` as the cadence-rule input; it does not define the receipt's shape.
- The per-cycle artifact channel (where `cdd-iteration.md` sits before the post-disconnect archive). The receipt-stream owns content/aggregation; artifact-channel owns path, role ownership, and archive timing.
- The cross-repo bundle composition (LINEAGE.md schemas per case; STATUS state machine; feedback-patch format; archival rule). Lives in `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. This skill names the *trigger* for a bundle (a `patch-landed` finding with `Cross-repo` target); cross-repo names the *bundle*.
- The activation-side cadence declaration (per-repo `cdd-iteration.md` cadence; D/C/B/A severity scale; auto-spawn MCA trigger). Lives in `cnos.cdd/skills/cdd/activation/SKILL.md §22`. The activation skill declares per-repo cadence overlays; this skill carries the generic wire-format invariants.
- Tooling that auto-verifies receipt-stream discipline (`cn cdd verify` cycle-iteration checks). Deferred; doctrine first, tooling later.

---

## 1. Per-finding shape (`cdd-iteration.md`)

Each finding inside `cdd-iteration.md` carries the following fields, in this order. The shape is preserved verbatim from the pre-migration canonical home (`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`):

```markdown
### F1: <title>

- **Source:** α-closeout / β-review / γ-triage / PRA §3 / PRA §4b — name the artifact and section
- **Class:** `cdd-skill-gap` | `cdd-protocol-gap` | `cdd-tooling-gap` | `cdd-metric-gap`
- **Trigger:** §9.1 trigger N | "γ process-gap check" | "review pattern across cycles"
- **Description:** one paragraph
- **Root cause:** one paragraph
- **Disposition:** `patch-landed` | `next-MCA` | `no-patch`

If `patch-landed`:
- **Patch:** `<commit-sha>` or `<file-path>`
- **Affects:** `<cdd skill file or section>` (e.g. `cdd/review/SKILL.md` §3.13)
- **Cross-repo:** `<target-repo>` PR `#NN` if patch landed in a different repo than this cycle (otherwise omit)

If `next-MCA`:
- **Issue filed:** `<repo>` `#NN`
- **First AC:** ...

If `no-patch`:
- **Reason:** required justification per `cnos.cds/skills/cds/CDS.md` §"Closure" → §"Closure rule"
```

**Field meaning:**

- **Source.** Which close-out artifact (α-closeout / β-review / γ-triage / PRA §3 / PRA §4b) surfaced this finding. Names the artifact + section; lets ε trace the finding back to its origin.
- **Class.** The CDD-side gap class per `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` — `cdd-skill-gap` (procedural skill underspecified / wrong); `cdd-protocol-gap` (CDD doctrine drifted from invariant); `cdd-tooling-gap` (tooling absent / wrong / unavailable); `cdd-metric-gap` (measurement missing / wrong). CDR uses different classes per `cnos.cdr/skills/cdr/epsilon/SKILL.md §1`; future c-d-X may use other classes.
- **Trigger.** Which trigger fired — a `cnos.cds/skills/cds/CDS.md §"Assessment" → §"Cycle iteration triggers"` row (e.g. review rounds > 2; mechanical ratio > 20%; avoidable tooling/environmental failure; loaded skill failed to prevent a finding), or the γ process-gap check, or a multi-cycle pattern surfaced by ε's review-pattern-across-cycles read.
- **Description.** One paragraph naming what happened, in concrete cycle-local terms.
- **Root cause.** One paragraph naming why it happened — what protocol invariant was unmet, what skill was underspecified, what tool was missing.
- **Disposition.** One of three values: `patch-landed` (the patch shipped in this cycle), `next-MCA` (issue filed for a follow-on cycle), `no-patch` (explicit no-patch with required reason). No other dispositions are valid.

**Per-disposition sub-fields:**

- **`patch-landed`** — `Patch:` names the commit SHA or file path; `Affects:` names the cdd skill file or section that received the patch; `Cross-repo:` names the target-repo PR if the patch landed in a different repo (otherwise omit).
- **`next-MCA`** — `Issue filed:` names the repo + issue number; `First AC:` names the first acceptance criterion the follow-on cycle ships against.
- **`no-patch`** — `Reason:` carries the required justification per `cnos.cds/skills/cds/CDS.md §"Closure" → §"Closure rule"` (one of: all findings were application gaps with adequate existing spec; zero review findings this cycle; environmental failure with no spec-level fix available).

The shape is **not redesigned by this skill**; the migration is wholesale. Future tightening (adding fields, changing field semantics) is a follow-on cycle; this skill carries the v0.1 shape.

---

## 2. Aggregator (`.cdd/iterations/INDEX.md`)

The cross-cycle aggregator at `.cdd/iterations/INDEX.md` carries one row per cycle that produced an iteration artifact. The row format is **eight pipe-separated columns**, preserved verbatim from the pre-migration canonical home:

```markdown
| Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path |
|-------|-------|------|----------|---------|------|----------|------|
| {N}   | #N    | YYYY-MM-DD | M | P | A | Z | .cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md |
```

**Column meanings:**

| Column | Meaning |
|---|---|
| **Cycle** | The cycle identifier (typically the issue number; may include sub-cycle suffix like `344-b` for multi-phase docs cycles). |
| **Issue** | The GitHub issue number (`#NNN`). Usually equals the cycle but may differ for follow-on cycles re-opening an earlier issue. |
| **Date** | ISO date (`YYYY-MM-DD`) of cycle close. |
| **Findings** | Count of findings inside the cycle's `cdd-iteration.md` (M). Zero for courtesy stubs (cycle/401 convention). |
| **Patches** | Count of findings with disposition `patch-landed` (P). Includes both same-cycle patches and cross-repo `patch-landed` shipments. |
| **MCAs** | Count of findings with disposition `next-MCA` (A). Each represents one follow-on cycle scheduled. |
| **No-patch** | Count of findings with disposition `no-patch` (Z). Each carries a recorded `Reason`. |
| **Path** | Relative path to the cycle's `cdd-iteration.md` artifact. Uses `.cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md` for version-tagged releases; uses `.cdd/releases/docs/{ISO-date}/{N}/cdd-iteration.md` for docs-only releases; uses `.cdd/unreleased/{N}/cdd-iteration.md` for in-flight or unreleased cycles. |

**Row update procedure:**

After writing `cdd-iteration.md` (or the courtesy empty-findings stub), γ updates `.cdd/iterations/INDEX.md` with one row appended at the bottom of the data table:

1. Open `.cdd/iterations/INDEX.md` on the cycle branch.
2. Append one row matching the format above, with the cycle's actual values (`{N}` / `#N` / ISO date / M / P / A / Z / path).
3. The four count columns sum: `Findings = Patches + MCAs + No-patch` (M = P + A + Z). This is a derived consistency check; ε's pattern-detection relies on it.
4. Commit on the cycle branch (same commit as `cdd-iteration.md` if both are written together).

The aggregator file is **data**, not doctrine. Its content rows are not moved by this migration; this skill names the *format* and the *update procedure*. The live file at `.cdd/iterations/INDEX.md` is the authoritative aggregator for the cnos repo at the current scale.

---

## 3. Cadence rules

The per-cycle iteration artifact's cadence is **conditional on the cycle's receipt**:

### 3.1. Required: `protocol_gap_count > 0`

When the cycle's receipt has `protocol_gap_count > 0` — i.e. the close-out triage in `gamma-closeout.md` produced ≥ 1 finding tagged `cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, or `cdd-metric-gap` — the per-cycle iteration artifact is **required**. γ MUST author `.cdd/unreleased/{N}/cdd-iteration.md` with each finding structured per the per-finding shape in §1; γ MUST append one INDEX.md row per §2; γ closure gate blocks closure until both exist.

The receipt's `protocol_gap_count` field is **consistency-pinned** to `len(protocol_gap_refs)` by `schemas/cdd/receipt.cue`. The receipt is the always-present record of *whether* an iteration file is required; the file is the conditional record of *what was found*. This is the inherited rule from `ROLES.md §4b.4`.

### 3.2. Permitted: courtesy empty-findings stub when `protocol_gap_count == 0`

When the cycle's receipt has `protocol_gap_count == 0` — i.e. the cycle ran cleanly with no `cdd-*-gap` findings — the per-cycle iteration artifact is **optional**. The cycle's gamma-closeout names the no-gap result; no iteration file is required.

However, γ MAY write a courtesy empty-findings stub at `.cdd/unreleased/{N}/cdd-iteration.md` for traceability. The stub records `protocol_gap_count: 0` explicitly, names the no-findings result, and (optionally) explains why no protocol gap surfaced. The stub variant is the convention landed in **cycle/401** (the empirical anchor for the courtesy-stub rule). Backward compatibility: existing empty-findings `cdd-iteration.md` files written under the prior every-cycle rule remain valid artifacts — the rule is "required when `> 0`", not "forbidden when `== 0`".

When a courtesy stub is written, γ MAY also append the INDEX.md row (with `Findings = 0; Patches = 0; MCAs = 0; No-patch = 0`). This is the dominant cnos cdd convention at current scale: nearly every recent cycle has an INDEX row, even for `protocol_gap_count == 0` cycles, because the row supports ε's per-cycle scan-without-omission discipline.

### 3.3. Activation-side overlay

The per-repo cadence declaration (whether courtesy stubs are mandatory in a given project; the D/C/B/A + `info` severity scale for findings; the sliding-window auto-spawn MCA trigger) is set by `cnos.cdd/skills/cdd/activation/SKILL.md §22`. The activation skill overlays per-repo policy on top of the generic wire-format invariants this skill carries; the two are layered.

---

## 4. Aggregator location

The aggregator file lives at:

```text
.cdd/iterations/INDEX.md
```

at the repo root. The file is **data**, not doctrine: its content rows accumulate cycle by cycle and are not moved by this migration. The doctrine *about* the file (the row format; the per-row update procedure; the eight column meanings) is what this skill carries.

The filename will eventually become consumer-specific. CDS's analogue is `<project>/.cds/iterations/INDEX.md` per `cnos.cds/skills/cds/CDS.md §"Closure"`. CDR's analogue is the per-project coherence-log per `cnos.cdr/skills/cdr/epsilon/SKILL.md §1` and `cnos.cdr/docs/empirical-anchor-cph.md` row #14. The wire-format pattern (one persistent aggregator at a known per-protocol-package path; one row per cycle that produced an iteration artifact; eight columns naming cycle / issue / date / findings / patches / MCAs / no-patch / path) is what this skill records; the per-protocol filename is consumer-side.

For the cnos repo at the current operating point, the file is at `.cdd/iterations/INDEX.md`. The bootstrap row was added by cycle/335 (per `cnos.cds/docs/empirical-anchor-cdd.md` #335 row); every cycle since cnos#364 carries one row.

---

## 5. Finding-disposition discipline

Each finding carries one of three dispositions; the vocabulary is **closed** (no other dispositions are valid):

- **`patch-landed`** — the patch shipped in this cycle. The finding is closed by the same commit set that produced the cycle's release. The patch lands in a CDD skill file, a doc, a tool, or (rarely) a runtime surface. If the patch landed in a different repo, the `Cross-repo:` sub-field names the target-repo + PR; the cross-repo trace bundle attaches per §6.
- **`next-MCA`** — the patch is deferred to a follow-on cycle. The cycle files an issue (the MCA), records the issue number under `Issue filed:`, and pins the `First AC:` the follow-on cycle ships against. This is the dominant disposition for findings that are not safe to ship inside the current cycle (cross-cycle redesigns; surfaces with > 1 dependent cycle).
- **`no-patch`** — the finding is closed without a patch. This requires explicit justification per `cnos.cds/skills/cds/CDS.md §"Closure" → §"Closure rule"`: one of (a) all findings were application gaps with adequate existing spec (the skill was right; the agent didn't follow it; patch the operator, not the spec); (b) zero review findings this cycle (the finding is a γ-side observation, not a review finding); (c) environmental failure with no spec-level fix available (tooling constraint; cross-cutting infrastructure). The `Reason:` sub-field is **required**; silence is not an acceptable disposition.

The vocabulary is **classification-grade**: ε's pattern-detection across cycles requires that the three dispositions are mutually exclusive and exhaustive. Adding a fourth disposition (e.g., "drop") is a follow-on cycle; this skill carries the v0.1 vocabulary.

The role-local authority for the disposition discipline is `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` (CDD-side). The MCA-discipline shape (ship-now / next-MCA / no-patch) is inherited from `ROLES.md §4b.5`; the CDD-specific gap-class names (`cdd-skill-gap` / `cdd-protocol-gap` / etc.) live in `cnos.cdd/skills/cdd/epsilon/SKILL.md §1`.

---

## 6. Cross-repo trace

When a finding's `Disposition` is `patch-landed` with a `Cross-repo:` target — i.e. the patch shipped in a different repo than this cycle — γ creates a **cross-repo trace bundle** at the canonical path:

```text
.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/
```

per `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. The bundle composition depends on the case (per `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2` directional cases):

- **Outbound iteration trace (case b)** — the bundle file set is `LINEAGE.md` with per-patch confirmation; one outbound destination patch.
- **Bilateral iteration (case c)** — the bundle additionally carries `cdd-iteration.md` and one feedback patch per outbound destination patch.

**Archival.** The source-side (cnos-side, for outbound) bundle MAY be archived once the counterpart-side mirror confirms receipt per `cnos.handoff/skills/handoff/cross-repo/SKILL.md §"Bundle archival rule"`; target-side mirrors are preserved indefinitely. The asymmetric archival rule preserves the bundle's lineage on the target side past the source-side cleanup.

This skill names the trigger (a `patch-landed` finding with `Cross-repo` target); `cnos.handoff/skills/handoff/cross-repo/SKILL.md` names the bundle's composition and archival rule. The two skills compose: the receipt-stream's per-finding shape exposes the trigger; cross-repo carries the bundle.

---

## 7. Empirical anchors

Every cycle since cnos#364 carries an INDEX.md row + a `cdd-iteration.md` (file or courtesy stub). Representative anchors:

- **cycle/335** — `.cdd/iterations/INDEX.md` initialized; bootstrap for the receipt-stream artifact. Per `cnos.cds/docs/empirical-anchor-cdd.md` #335: "Initialized `.cdd/iterations/INDEX.md`; bootstrap for the receipt-stream artifact". The 8-column row format crystallized here.
- **cycle/401** — the courtesy empty-findings stub convention landed. Before cycle/401, the every-cycle `cdd-iteration.md` rule produced files for `protocol_gap_count == 0` cycles too (the survivorship-bias ambiguity: "is the missing file a clean cycle or a skipped close-out?"). cycle/401 resolved the ambiguity structurally — the receipt's `protocol_gap_count` field is the always-present record of *whether* an iteration file is required; the file is the conditional record of *what was found*. Empty-findings files written under the prior rule remain valid artifacts (backward compatibility); future cycles with `protocol_gap_count == 0` may write a courtesy stub or omit the file entirely.
- **cycle/369** — Phase 2 receipt schemas; first cycle to record `cdd-protocol-gap` findings under the typed receipt + iteration shape. 3 findings, 0 patches, 2 MCAs, 1 no-patch.
- **cycle/388** — Phase 2.5 schemas split; demonstrated the `cdd-protocol-gap` + `cdd-skill-gap` mixed-class triage in one cycle. 3 findings, 0 patches, 2 MCAs.
- **cycle/393** — δ-as-architect protocol patches; 4 patches landed on the iteration surface in one cycle (the `patch-landed` disposition exercised at scale). 0 findings new; 4 patches landed against prior MCAs.
- **cycle/410** — Sub 5 of cnos#403 (eight-section migration); closure with courtesy `cdd-iteration.md` stub (zero findings). Demonstrates the closure-stub option as the standard pattern at current operating point.
- **The #406–#412 wave** — post-codification cycles that exercised the receipt-stream in steady-state. Every cycle in the wave produced a complete `cdd-iteration.md` (or courtesy stub) + an INDEX.md row; the receipt-stream wire-format is operationally stable.

The full row history is at `.cdd/iterations/INDEX.md` — 22+ rows at the source pin, accumulating one per cycle. The empirical anchor index for CDS-class cycles is at `cnos.cds/docs/empirical-anchor-cdd.md` (which cross-references the receipt-stream's per-cycle artifacts as the canonical receipt instances per CDS Field 3).

---

## 8. Related documents

- `cnos.handoff/skills/handoff/HANDOFF.md` — package contract; this skill is one of five sub-surfaces.
- `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` — the per-cycle channel substrate. The `cdd-iteration.md` file lives on this channel as a γ-owned file (per artifact-channel §2.2); the channel governs the file path, write ownership, and freeze-on-merge rule; this skill governs the file's content shape and the cross-cycle aggregation.
- `cnos.handoff/skills/handoff/cross-repo/SKILL.md` — the cross-repo bundle composition. When a `patch-landed` finding has a `Cross-repo` target, the bundle attaches per the cross-repo skill's case-(b) / case-(c) bundle file set + asymmetric archival rule.
- `cnos.handoff/skills/handoff/dispatch/SKILL.md` — the dispatch wire format; the dispatch prompt initiates the cycle whose close-out feeds this receipt-stream.
- `cnos.handoff/skills/handoff/mid-flight/SKILL.md` — the rescue mechanism; orthogonal to receipt-stream (rescue happens during the cycle; receipt-stream consumes the cycle's close-out artifact).
- `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` — the cdd-side runbook (now a pointer to this skill). post-release/SKILL.md applies the rule at cycle close-out: γ writes `cdd-iteration.md` per the canonical shape in §1; γ appends the INDEX row per §2; γ runs the cross-repo bundle composition per §6 if applicable.
- `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` — ε's CDD-side role-local authority on the gap-class instantiation (`cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap`). The `Class:` field in this skill's per-finding shape takes one of these four values for CDD-side findings; CDR's `Class:` values are different per `cnos.cdr/skills/cdr/epsilon/SKILL.md §1`.
- `cnos.cdd/skills/cdd/activation/SKILL.md §22` — the per-repo cadence declaration (whether courtesy stubs are mandatory; D/C/B/A + `info` severity scale; sliding-window auto-spawn MCA trigger). The activation skill overlays per-repo policy on top of the generic wire-format invariants this skill carries.
- `cnos.cdd/skills/cdd/COHERENCE-CELL.md §"ε as Protocol Evolution" + §"ε Artifact Rule"` — the cell-doctrinal grounding for ε's role; the receipt-stream is the artifact surface ε reads.
- `cnos.cds/skills/cds/CDS.md §"Closure"` — CDS-class closure gate; F9 ("Missing `cdd-iteration.md` when triggers fired") names the gate that blocks closure when `protocol_gap_count > 0` but no iteration artifact exists. CDS.md cites this skill (post-#419) for the per-finding shape; the operational realization in `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` is a pointer to this skill.
- `cnos.cds/skills/cds/CDS.md §"Artifact contract" → §"Location matrix"` — declares the `cdd-iteration.md` file's canonical path in the CDS-class artifact set; the file lives on the artifact-channel and its content shape is governed by this skill.
- `schemas/cdd/receipt.cue` (`#Receipt`) — the typed receipt that carries `protocol_gap_count` and `protocol_gap_refs`. This skill consumes the receipt's `protocol_gap_count` as the cadence-rule input (§3.1).
- `ROLES.md §4b` — the generic ε doctrine; the authoritative source for the watched-fields invariant, the MCA discipline, and the ε=δ collapse rule. The cadence rule (`required only when protocol_gap_count > 0`) is inherited from §4b.4.
- [cnos#335](https://github.com/usurobor/cnos/issues/335) — INDEX.md aggregator bootstrap empirical anchor.
- [cnos#364](https://github.com/usurobor/cnos/issues/364) — first cycle in the steady-state INDEX.md every-cycle period.
- [cnos#401](https://github.com/usurobor/cnos/issues/401) — courtesy empty-findings stub convention empirical anchor.
- [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker (handoff/coordination extraction wave).
- [cnos#419](https://github.com/usurobor/cnos/issues/419) — the codification cycle for this surface (Sub 5 of cnos#404).

---

## 9. Non-goals

- This skill does NOT redesign the per-finding shape. The six top-level fields (`Source` / `Class` / `Trigger` / `Description` / `Root cause` / `Disposition`) and the per-disposition sub-fields are preserved verbatim from the pre-migration canonical home (`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`). Future tightening is a follow-on cycle.
- This skill does NOT redesign the INDEX.md row format. The eight pipe-separated columns (`Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path`) are preserved verbatim. Future tightening (e.g., adding a sub-cycle suffix column, separating same-cycle from cross-repo patch counts) is a follow-on cycle.
- This skill does NOT modify `.cdd/iterations/INDEX.md` content. The file's existing rows stay where they are; only the standard close-out row for this cycle is added (per the artifact-channel cycle close-out rule and the post-release §5.6b pointer this skill now governs).
- This skill does NOT own the role-local authoring procedure. Which session writes the artifact, the close-out gate ordering, the ε=δ collapse rule for the current operating point — those live in `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b` (cdd-side runbook, now pointer) and `cnos.cdd/skills/cdd/epsilon/SKILL.md` (ε role-local).
- This skill does NOT own the CDD-specific gap-class names (`cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap`). Those live in `cnos.cdd/skills/cdd/epsilon/SKILL.md §1`. The `Class:` slot in this skill's per-finding shape is protocol-agnostic; the values are consumer-specific.
- This skill does NOT own the cross-repo bundle composition (LINEAGE.md schemas; STATUS state machine; feedback-patch format; archival rule). Those live in `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. This skill names the trigger (a `patch-landed` finding with `Cross-repo` target); cross-repo names the bundle.
- This skill does NOT own the polling primitives or harness substrate. Those live in `cnos.cdd/skills/cdd/harness/SKILL.md §5.4` and `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → §"Polling primitives"`. This skill is an artifact-surface for a cross-cycle stream; the read mechanism (how ε polls the aggregator across cycles) is harness-side.
- This skill does NOT lift the per-finding shape or the INDEX row format into a typed schema (e.g. `schemas/handoff/iteration-finding.cue` / `schemas/handoff/iteration-index.cue`). The Markdown invariants are the v0.1 canonical rendering; a typed schema is deferred per the extraction-map.md §6 note (Sub 3's Q2 ruling precedent — schema stays Markdown for v0.1). Re-evaluate if a downstream consumer (CCNF-O Track A; `cn cdd verify`; cross-repo bundle validators) needs a typed shape.

---

## 10. Kata

### Scenario

δ has just reported the cycle's tagged disconnect and green release CI. The role close-outs and γ's nonterminal marker are on main, and the cycle still lives under `.cdd/unreleased/{N}/`. Its receipt carries `protocol_gap_count: 2` with two refs naming `cdd-skill-gap` findings surfaced during review. γ has not yet archived or declared terminal closure.

### Expected reasoning

1. **Cadence check (§3.1).** The receipt has `protocol_gap_count > 0` (specifically, 2). The per-cycle iteration artifact is **required**, not optional. γ MUST author `.cdd/unreleased/{N}/cdd-iteration.md` and MUST append one INDEX.md row.

2. **Per-finding authoring (§1).** γ opens `.cdd/unreleased/{N}/cdd-iteration.md` and writes one F1 / F2 section per finding. Each section carries the six top-level fields verbatim per the per-finding shape:
   - `Source:` names the artifact + section (e.g. `β-review §3 R1`).
   - `Class:` is `cdd-skill-gap` (per the receipt's refs).
   - `Trigger:` names the CDS §"Assessment" trigger that fired (e.g. "mechanical ratio > 20% with ≥ 10 findings" if applicable).
   - `Description:` one paragraph of cycle-local detail.
   - `Root cause:` one paragraph naming the unmet invariant.
   - `Disposition:` one of `patch-landed` / `next-MCA` / `no-patch`. γ picks based on whether the patch shipped this cycle, is deferred to a follow-on, or is closed with explicit no-patch reason.

3. **Per-disposition sub-fields (§1).** For each finding:
   - If `patch-landed`: γ records the commit SHA / file path under `Patch:`, names the affected skill file/section under `Affects:`, and (if cross-repo) names the target-repo + PR under `Cross-repo:`.
   - If `next-MCA`: γ files the issue (e.g. `gh issue create`), records the issue number under `Issue filed:`, and pins the `First AC:`.
   - If `no-patch`: γ records the `Reason:` per `cnos.cds/skills/cds/CDS.md §"Closure" → §"Closure rule"`.

4. **Aggregator row (§2).** γ opens `.cdd/iterations/INDEX.md` and appends one row:
   ```
   | N | #N | YYYY-MM-DD | 2 | P | A | Z | .cdd/unreleased/N/cdd-iteration.md |
   ```
   where P + A + Z = 2 (the consistency check per §2). γ commits both files in the same post-disconnect learning commit before the archive move.

5. **Cross-repo bundle (§6).** If any finding had `Disposition: patch-landed` with a `Cross-repo:` target, γ creates a bundle at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` per `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. The bundle's LINEAGE.md names the source-side patch + the destination-side PR.

6. **Archive and terminal seal.** After the PRA, dispositions, and aggregation are complete, γ moves the cycle directory to `.cdd/releases/{X.Y.Z}/{N}/` in its own archive commit. γ then appends the terminal declaration to the archived `gamma-closeout.md` in a later commit binding the release tag and archive SHA. Only that final commit closes the cycle.

6. **Release-time move (artifact-channel §2.4.2).** Before γ requests the tag from δ, γ moves `.cdd/unreleased/{N}/` → `.cdd/releases/{X.Y.Z}/{N}/` (or `.cdd/releases/docs/{ISO-date}/{N}/` for docs-only releases). The `cdd-iteration.md` file rides the move; the INDEX.md row's `Path` column is updated to reflect the new path in a subsequent close-out commit if not already pinned at write time.

### Common failures

- γ writes `cdd-iteration.md` but forgets the INDEX row — ε's pattern-detection across cycles skips this cycle when scanning the aggregator. CDS gate F9 catches this.
- γ writes the INDEX row but the per-finding shape diverges (e.g. uses `Status:` instead of `Disposition:`; uses a new class name not in the four-value vocabulary). ε's class-by-class scan across cycles breaks.
- γ writes a `no-patch` disposition without a `Reason:` — closure rule is violated; the silence is not an acceptable disposition.
- γ writes a courtesy stub (`protocol_gap_count == 0`) and labels it "iteration complete" with the no-findings result *despite* there being findings. The survivorship-bias ambiguity re-emerges; ε cannot distinguish clean cycles from skipped triage.
- A `patch-landed` finding with a `Cross-repo:` target lands without the bundle at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`. The patch's lineage is lost when the source-side bundle archives.

### Verification

- `.cdd/unreleased/{N}/cdd-iteration.md` exists if `protocol_gap_count > 0`.
- `.cdd/iterations/INDEX.md` has one new row for cycle N with the eight columns populated.
- For each finding, the six top-level fields are present + the per-disposition sub-fields are present.
- `Findings = Patches + MCAs + No-patch` in the INDEX row (M = P + A + Z).
- For each `patch-landed` + `Cross-repo` finding, `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` exists with the bundle.

### Reflection

A second γ reading the cycle's record across cycles (via `.cdd/iterations/INDEX.md` row-by-row) should be able to surface protocol-gap patterns no single cycle could surface — e.g., "five cycles in a row had `cdd-skill-gap` findings against `review/SKILL.md §3.13`; the skill is underspecified and needs a patch". The receipt-stream is the cross-cycle learning signal; without it, ε would have to re-derive cycle context from gh logs and per-cycle assessment files for every triage. The artifact-channel carries the per-cycle record; the receipt-stream carries the cross-cycle pattern surface.

---

**End of skill.**
