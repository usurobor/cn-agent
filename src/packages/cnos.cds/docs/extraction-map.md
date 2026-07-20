# CDS Extraction Map — v0.1

> **Status:** Sub 6 ([cnos#411](https://github.com/usurobor/cnos/issues/411))
> complete — `CDD.md` "pending cds extraction" markers replaced with CDS
> pointers; cross-references in 11 cdd skill files + 2 YAML templates
> re-pointed at CDS canonical homes (AC1–AC9 PASS).

**Version:** 0.1 (Sub 1 of [cnos#403](https://github.com/usurobor/cnos/issues/403); shipped under [cnos#406](https://github.com/usurobor/cnos/issues/406))
**Date:** 2026-05-22
**Placement:** `src/packages/cnos.cds/docs/extraction-map.md`
**Audience:** δ operators dispatching Subs 3–5 of [cnos#403](https://github.com/usurobor/cnos/issues/403); reviewers verifying that every "pending cds extraction" surface in [`cnos.cdd/skills/cdd/CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) has a named CDS destination.
**Scope:** A surface-by-surface migration plan from the pre-#402 software-lifecycle content currently quarantined in [`cnos.cdd/skills/cdd/CDD.md` §"Software-specific realization — pending cds extraction"](../../cnos.cdd/skills/cdd/CDD.md) onto its planned CDS canonical home, with the migration sub-issue named for each surface.

> This document is a **mapping table + commentary**, not a re-statement of CDD
> or of CDS. It does not embed CDD prose; citations carry the source path and
> the section marker at which the content currently resides. Subs 3–5 of
> [cnos#403](https://github.com/usurobor/cnos/issues/403) execute the
> migration against this map; Sub 6 removes the "pending cds extraction"
> markers from `CDD.md` once Subs 3–5 land.

---

## 0. Method

### 0.1 Source pin

CDD evidence cited herein refers to the cnos repository's `main` branch at
the head of `cycle/406`'s parent commit
(`ecbcb5d5ec4f16a72328836d2e756a3e9bc278b8` — the Merge cycle/402 commit
that compressed CDD.md to the CCNF spine). The "Software-specific
realization — pending cds extraction" section spans lines 122–141 of
[`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../cnos.cdd/skills/cdd/CDD.md)
at that commit; each marker named in the table below appears as a bullet
in that section family.

### 0.2 CDS v0.1 target surfaces

The migration lands into three surface classes within `cnos.cds`:

1. **`skills/cds/CDS.md`** — the canonical instantiation contract. Sections under `CDS.md` are named `§"Selection function"`, `§"Development lifecycle"`, etc., mirroring the family-naming convention CDD.md used. Shipped by Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403); deeper per-surface SKILL.md files are shipped by Subs 3–5.
2. **`skills/cds/<area>/SKILL.md`** — sub-area skill files for operationally-heavy surfaces (lifecycle, artifacts, review, gate, assessment). Each sub-area skill is the operational expansion of the CDS.md section that names it. Pointer-only overlays delegating to the existing `cnos.cdd/skills/cdd/<role>/SKILL.md` operational realizations are acceptable in v0.1; deep rewrites are deferred to post-v0.1 cycles.
3. **`schemas/cds/`** — already exists per [cnos#388](https://github.com/usurobor/cnos/issues/388); CDS Field 3 (γ close-out artifact) will cite `schemas/cds/receipt.cue`'s `#CDSReceipt` schema for the typed receipt shape. No schema-side migration in Subs 3–5; the migration is doctrine-surface only.

### 0.3 Sub naming

The migration sub-issues follow [cnos#403](https://github.com/usurobor/cnos/issues/403)'s wave shape:

- **Sub 3** — "Migrate selection function + development lifecycle to CDS" (P2; depends on Sub 2). Owns rows whose source is §Selection or §Lifecycle.
- **Sub 4** — "Migrate artifact contract + evidence surfaces to CDS" (P2; depends on Sub 2). Owns rows whose source is §Artifacts (the artifact contract family), §Inputs (selection inputs the observer reads, which is an evidence surface), or §Tracking (coordination evidence).
- **Sub 5** — "Migrate review CLP + gate + closure + assessment + retro-packaging to CDS" (P2; depends on Sub 2). Owns rows whose source is §Review, §Gate, §Assessment, §Closure, §Retro-packaging, §Mechanical, §Non-goals, or §Large-file.

The Sub-3/Sub-4/Sub-5 split follows the natural cohesion of the surfaces; cross-cutting bullets (e.g. the §Roles family, which §Lifecycle and §Artifacts both reference) are noted in the row's commentary.

### 0.4 Coverage commitment

Per [cnos#406](https://github.com/usurobor/cnos/issues/406) AC5 ("exactly one row for each surface in #403's 'Source content' list (10 surface groups; the map may expand the per-surface sub-list as multiple rows)"), the tables below carry one row per top-level surface (with sub-list expansions where the source CDD.md bullet itself enumerates multiple distinct mechanical surfaces). The mechanical check `grep -c "^|" docs/extraction-map.md` is ≥ 10 + per-table header/separator overhead.

Every "pending cds extraction" marker in [`cnos.cdd/skills/cdd/CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" (lines 122–141 at the source-pin commit) appears as at least one row below. The fourteen marker families are: Inputs, Selection, Lifecycle, Roles and dispatch, Coordination surfaces (Tracking), Artifacts, Mechanical, Review, Gate, Assessment, Closure, Retro-packaging, Non-goals, Large-file. Of those, the ten surface groups named in [cnos#403](https://github.com/usurobor/cnos/issues/403)'s "Source content" list (Selection, Lifecycle, Artifacts, Mechanical, Review, Gate, Assessment, Closure, Retro-packaging, Non-goals) each receive a dedicated table below; the four additional CDD.md marker families (Inputs, Roles and dispatch, Coordination surfaces, Large-file) are folded into the Artifact-contract or Lifecycle tables with the cross-surface origin called out in the row commentary, so the marker sweep at Sub 6 will find every marker covered.

---

## 1. Selection function (§Selection) — Sub 3

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Selection function (§Selection)" (line 127 at source pin).
**Source content:** P0 override, operational-infrastructure override, assessment-commitment default, stale-backlog re-evaluation, MCI freeze check, weakest-axis rule, maximum-leverage rule, dependency order, effort-adjusted tie-break, no-gap case. Operational realization currently in `gamma/SKILL.md`.
**Migration sub:** Sub 3 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Selection function"`](../skills/cds/CDS.md).** Landed via [cnos#408](https://github.com/usurobor/cnos/issues/408) (Sub 3, B-lite thin extract). The 10 named rules + the §Inputs cross-cut now live as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/gamma/SKILL.md §2.1–§2.2` (cited from `CDS.md §"Selection function" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Selection — P0 override | `skills/cds/CDS.md §"Selection function" → "P0 override"` | Sub 3 | Verbatim move; P0-override is the gate-floor every selection function must respect. |
| §Selection — operational-infrastructure override | `skills/cds/CDS.md §"Selection function" → "Operational-infrastructure override"` | Sub 3 | Verbatim move; software-class operational overrides (CI red, build broken). |
| §Selection — assessment-commitment default | `skills/cds/CDS.md §"Selection function" → "Assessment-commitment default"` | Sub 3 | Verbatim move; ties to PRA outputs (Field 5 cadence). |
| §Selection — stale-backlog re-evaluation | `skills/cds/CDS.md §"Selection function" → "Stale-backlog re-evaluation"` | Sub 3 | Verbatim move. |
| §Selection — MCI freeze check | `skills/cds/CDS.md §"Selection function" → "MCI freeze check"` | Sub 3 | Verbatim move; MCI is a software-class concept (master coherence issue). |
| §Selection — weakest-axis rule | `skills/cds/CDS.md §"Selection function" → "Weakest-axis rule"` | Sub 3 | Verbatim move; α/β/γ axis scoring is the CDS evidence-summary shape. |
| §Selection — maximum-leverage rule | `skills/cds/CDS.md §"Selection function" → "Maximum-leverage rule"` | Sub 3 | Verbatim move. |
| §Selection — dependency order | `skills/cds/CDS.md §"Selection function" → "Dependency order"` | Sub 3 | Verbatim move; dependency-order respects software-class wave shapes. |
| §Selection — effort-adjusted tie-break | `skills/cds/CDS.md §"Selection function" → "Effort-adjusted tie-break"` | Sub 3 | Verbatim move. |
| §Selection — no-gap case | `skills/cds/CDS.md §"Selection function" → "No-gap case"` | Sub 3 | Verbatim move; gates the cycle on whether selection has any non-trivial output. |
| §Inputs — selection inputs (CHANGELOG TSC table, encoding lag table, doctor/status health surface, last PRA) | `skills/cds/CDS.md §"Selection function" → "Inputs"` | Sub 3 | Cross-cuts §Inputs (line 126 of source pin); §Inputs marker is folded here because §Inputs is "what the selection observer reads." Operational realization stays in `gamma/SKILL.md` (CDS sub-area `gamma/SKILL.md` pointer-overlay acceptable for v0.1). |

---

## 2. Development lifecycle (§Lifecycle) — Sub 3

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Development lifecycle (§Lifecycle)" (line 128 at source pin).
**Source content:** the 0–13 step table; the lifecycle state machine S0–S13; the branch rule (`cycle/{N}` canonical, γ creates from `origin/main` before dispatch); branch pre-flight; skill loading tier structure (1a CDD authority, 1b lifecycle phase skills, 1c β closure bundle, 2 general engineering, 3 issue-specific). Operational realization in `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`.
**Migration sub:** Sub 3 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Development lifecycle"`](../skills/cds/CDS.md).** Landed via [cnos#408](https://github.com/usurobor/cnos/issues/408) (Sub 3, B-lite thin extract). The 0–13 step table, the S0–S13 state machine, the `cycle/{N}` branch rule, the γ-owned branch pre-flight, and the tier-1a/1b/1c/2/3 skill loading structure now live as canonical content in `CDS.md`. The §Roles cross-cut content (triadic rule, dyad-plus-coordinator, dispatch model, dispatch-prompt formats) named in this row's table remains **pending** — it cross-cuts the Roles family (extraction-map §2 row "§Roles — …") whose deep-role-rewrite is deferred per the B-lite scope ruling; the v0.1 operational overlay for the role-side mechanics continues to live in `cnos.cdd/skills/cdd/{gamma,alpha,beta,delta}/SKILL.md` + `harness/SKILL.md` + `operator/SKILL.md` (cited from `CDS.md §"Development lifecycle" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Lifecycle — 0–13 step table | `skills/cds/CDS.md §"Development lifecycle" → "Step table"` (or `skills/cds/lifecycle/SKILL.md §"Steps"` if Sub 3 opts to expand to a sub-skill) | Sub 3 | Rewrite with CDS-specific framing; "code", "tests", "diff", "branch", "merge" are software-class vocabulary. CDS.md is the canonical home; a `skills/cds/lifecycle/SKILL.md` sub-skill may host the per-step procedural detail if Sub 3 chooses to expand. |
| §Lifecycle — lifecycle state machine S0–S13 | `skills/cds/CDS.md §"Development lifecycle" → "State machine"` | Sub 3 | Software-class states retained; cycle #669 split post-disconnect archive (S12) from terminal close (S13). |
| §Lifecycle — branch rule (`cycle/{N}` canonical; γ creates from `origin/main`) | `skills/cds/CDS.md §"Development lifecycle" → "Branch rule"` | Sub 3 | Verbatim move. The branch rule is CDS-specific (research-class CDR uses wave-shaped branches; software-class CDS uses `cycle/{N}`). |
| §Lifecycle — branch pre-flight | `skills/cds/CDS.md §"Development lifecycle" → "Branch pre-flight"` | Sub 3 | Verbatim move. |
| §Lifecycle — skill loading tier structure (1a/1b/1c/2/3) | `skills/cds/CDS.md §"Development lifecycle" → "Skill loading tiers"` | Sub 3 | Rewrite: tier 1a was "CDD authority" in pre-#402 doctrine; in CDS the tier-1a authority is `CDS.md`, with the generic CCNF kernel cited from `cnos.cdd` as a tier-1a-also-loaded reference. |
| §Roles — triadic rule (α produces, β judges and merges, γ orchestrates) | `skills/cds/CDS.md §"Roles and dispatch" → "Triadic rule"` | Sub 3 | Cross-cuts §Roles (line 129 of source pin). Verbatim move; the triadic rule is the CDS-specific form of the role-cell grammar inherited from CCNF. |
| §Roles — dyad-plus-coordinator framing; δ at external boundary | `skills/cds/CDS.md §"Roles and dispatch" → "Dyad-plus-coordinator"` | Sub 3 | Verbatim move. |
| §Roles — γ/α/β/δ algorithms; γ dispatch prompt format; named operator-decision points | `skills/cds/<role>/SKILL.md` (pointer overlays acceptable in v0.1) | Sub 3 (skeletal); deep rewrites deferred | The per-role algorithms live in role-overlay SKILLs; v0.1 may ship pointer overlays delegating to `cnos.cdd/skills/cdd/<role>/SKILL.md`. |
| §Roles — sequential bounded dispatch model (§1.6); re-dispatch prompt formats (§1.6a, §1.6b); initial dispatch sizing, prompt scope, commit checkpoints (§1.6c) | `skills/cds/CDS.md §"Roles and dispatch" → "Dispatch model"` (or `skills/cds/dispatch/SKILL.md` if expanded) | Sub 3 | Verbatim move; dispatch-prompt formats are software-class operational shape; the §1.6a/§1.6b/§1.6c anchors migrate as `CDS.md §"Dispatch model" → §"Initial"`, `→ §"Re-dispatch with patch"`, `→ §"Re-dispatch full-review"`. |

---

## 3. Coordination surfaces (§Tracking) — Sub 4

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Coordination surfaces (§Tracking)" (line 130 at source pin).
**Source content:** issue activity; cycle branch state; `.cdd/unreleased/{N}/` directory state; polling query forms (gh / MCP / git); wake-up mechanism; reachability preflight; transition-only emission; synchronous baseline pull; `git fetch` reliability re-probe; issue-edit cache-bust via `gamma-clarification.md`; cross-repo proposal lifecycle and STATUS state machine. Operational realization currently in `harness/SKILL.md §5.4`, `gamma/SKILL.md`, `cross-repo/SKILL.md`.
**Migration sub:** Sub 4 of [cnos#403](https://github.com/usurobor/cnos/issues/403). (Coordination surfaces are evidence-class for the software cycle.)
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Coordination surfaces"`](../skills/cds/CDS.md).** Landed via [cnos#409](https://github.com/usurobor/cnos/issues/409) (Sub 4, B-lite thin extract). The four sub-surfaces — §Cycle-state evidence, §Polling primitives, §Mid-flight clarification, §Cross-repo proposals — now live as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/harness/SKILL.md §5.4`, `gamma/SKILL.md §2.5`, and `cross-repo/SKILL.md §2.3` (cited from `CDS.md §"Coordination surfaces" → "Operational realization"`) until the v1 CDS-side role rewrite. **`.cdd/` → `.cds/` re-rooting documented in §Cycle-state evidence as planned (not performed; deferred to a separate post-#403 cycle).** Cross-repo `SKILL.md` long-term home left open per [cnos#404](https://github.com/usurobor/cnos/issues/404); skill currently cited from `cnos.cdd/skills/cdd/cross-repo/SKILL.md`.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Tracking — issue activity; cycle branch state; `.cdd/unreleased/{N}/` directory state | `skills/cds/CDS.md §"Coordination surfaces" → "Cycle-state evidence"` | Sub 4 | The `.cdd/unreleased/{N}/` directory will likely migrate to `.cds/unreleased/{N}/` once the project binding switches; the v0.1 doc records the planned re-rooting without performing it. |
| §Tracking — polling query forms (gh / MCP / git); wake-up mechanism; reachability preflight; transition-only emission; synchronous baseline pull; `git fetch` reliability re-probe | `skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"` | Sub 4 | Verbatim move; the operational realization stays in `harness/SKILL.md §5.4` (runtime substrate) and is cited by `CDS.md` rather than duplicated. |
| §Tracking — issue-edit cache-bust via `gamma-clarification.md` | `skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` | Sub 4 | Verbatim move. |
| §Tracking — cross-repo proposal lifecycle and STATUS state machine | `skills/cds/CDS.md §"Coordination surfaces" → "Cross-repo proposals"` | Sub 4 | Verbatim move; operational realization stays in `cross-repo/SKILL.md` (currently in `cnos.cdd`; may move to `cnos.cds` or stay in `cnos.cdd` as a generic substrate — Sub 4 decides). |

---

## 4. Artifact contract (§Artifacts) — Sub 4

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Artifact contract (§Artifacts)" (line 131 at source pin).
**Source content:** terminology (post-release, assessment, close-out, closure); bootstrap; ordered artifact flow (design → contract → plan → tests → code → docs → self-coherence → review → gate → release → observe → assess → close); artifact manifest with per-step format spec; **Artifact Location Matrix** (canonical paths for `.cdd/unreleased/{N}/self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md`, `RELEASE.md`, the version-snapshot directory, the PRA, the cross-repo trace dir); **role/artifact ownership matrix**; CDD Trace format; supporting rules; frozen snapshot rule. Operational realization in `release/SKILL.md`, `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`.
**Migration sub:** Sub 4 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Artifact contract"`](../skills/cds/CDS.md).** Landed via [cnos#409](https://github.com/usurobor/cnos/issues/409) (Sub 4, B-lite thin extract). The nine sub-surfaces — §Terminology, §Bootstrap, §Ordered flow, §Manifest, §Location matrix, §Ownership matrix, §Trace format, §Supporting rules, §Frozen snapshot rule — now live as canonical content in `CDS.md`. The CDD Trace is renamed **CDS Trace** in the new home (format verbatim). The v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/release/SKILL.md §2.5a` + `§3.8`, `gamma/SKILL.md §2.6–§2.10`, `alpha/SKILL.md §2.6`, `beta/SKILL.md`, and `release-effector/SKILL.md` (cited from `CDS.md §"Artifact contract" → "Operational realization"`) until the v1 CDS-side role rewrite. **`.cdd/` → `.cds/` re-rooting documented in §Location matrix as planned (not performed; deferred to a separate post-#403 cycle).**

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Artifacts — terminology (post-release, assessment, close-out, closure) | `skills/cds/CDS.md §"Artifact contract" → "Terminology"` | Sub 4 | Verbatim move; the terminology is software-class. |
| §Artifacts — bootstrap | `skills/cds/CDS.md §"Artifact contract" → "Bootstrap"` | Sub 4 | Verbatim move. |
| §Artifacts — ordered artifact flow (13 stages) | `skills/cds/CDS.md §"Artifact contract" → "Ordered flow"` | Sub 4 | Verbatim move; the 13-stage flow is the CDS-specific cycle shape (CDR's research flow differs). |
| §Artifacts — artifact manifest with per-step format spec | `skills/cds/CDS.md §"Artifact contract" → "Manifest"` (or `skills/cds/artifacts/SKILL.md §"Manifest"` if expanded) | Sub 4 | Verbatim move; per-step format specs may expand into `skills/cds/artifacts/SKILL.md` if the manifest is large. |
| §Artifacts — Artifact Location Matrix (canonical paths under `.cdd/unreleased/{N}/`) | `skills/cds/CDS.md §"Artifact contract" → "Location matrix"` | Sub 4 | Re-rooted: paths migrate from `.cdd/unreleased/{N}/` to `.cds/unreleased/{N}/` once the project binding switches. The v0.1 map records the planned re-rooting; the actual `.cdd/`-to-`.cds/` filesystem migration is out of scope for Sub 4 and may itself be a separate post-403 cycle. |
| §Artifacts — role/artifact ownership matrix | `skills/cds/CDS.md §"Artifact contract" → "Ownership matrix"` | Sub 4 | Verbatim move; the ownership matrix is software-class (α/β/γ/δ/ε ownership of self-coherence, review, close-outs, etc.). |
| §Artifacts — CDD Trace format | `skills/cds/CDS.md §"Artifact contract" → "Trace format"` (renamed "CDS Trace") | Sub 4 | Rename: "CDD Trace" → "CDS Trace" since the trace is software-cycle-specific. The trace format itself is verbatim. |
| §Artifacts — supporting rules | `skills/cds/CDS.md §"Artifact contract" → "Supporting rules"` | Sub 4 | Verbatim move. |
| §Artifacts — frozen snapshot rule | `skills/cds/CDS.md §"Artifact contract" → "Frozen snapshot rule"` | Sub 4 | Verbatim move; the rule pins that artifacts in `.cds/unreleased/{N}/` freeze on merge. |

---

## 5. Mechanical vs judgment boundary (§Mechanical) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Mechanical vs judgment boundary (§Mechanical)" (line 132 at source pin).
**Source content:** what may be enforced by tools (branch naming, version-directory presence, AC accounting, stale-cross-reference detection, gate checks) vs what remains judgment-bearing (the real incoherence, MCA vs MCI, scoring soundness, review convergence, design coherence, when to stop iterating).
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Mechanical vs judgment"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The mechanical axes (branch naming, version-directory presence, AC accounting, stale-cross-reference detection, gate checks) and judgment axes (real incoherence, MCA vs MCI, scoring soundness, review convergence, design coherence, when-to-stop-iterating) now live as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/review/SKILL.md §"Finding Taxonomy"`, `§3.12`, `§3.13`, and `cnos.cdd/skills/cdd/gamma/SKILL.md §3.9` (cited from `CDS.md §"Mechanical vs judgment" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Mechanical — tool-enforceable axes (branch naming, version-dir presence, AC accounting, stale-xref detection, gate checks) | `skills/cds/CDS.md §"Mechanical vs judgment" → "Mechanical axes"` | Sub 5 | Verbatim move; each tool-enforceable check has a corresponding `cn cds <check>` aspiration (out of scope for migration; CDS Field 3 cites schemas/cds/ for the typed receipt). |
| §Mechanical — judgment-bearing axes (real incoherence, MCA vs MCI, scoring soundness, review convergence, design coherence, when-to-stop-iterating) | `skills/cds/CDS.md §"Mechanical vs judgment" → "Judgment axes"` | Sub 5 | Verbatim move; the judgment axes anchor the configuration-floor (α=β is prohibited because the judgment axes require independent oracles). |

---

## 6. Review CLP (§Review) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Review (§Review)" (line 133 at source pin).
**Source content:** CLP (TERMS / POINTER / EXIT); reviewer ask list (α/β/γ scores, weakest-axis diagnosis, concrete patch suggestions, iterate-or-converge verdict). Operational realization in `review/SKILL.md`.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Review CLP"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The CLP form (TERMS / POINTER / EXIT) and the reviewer ask list (α/β/γ scores, weakest-axis diagnosis, concrete patch suggestions, iterate-or-converge verdict, configuration-floor flags, optional empirical-anchor cite) now live as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/review/SKILL.md` (cited from `CDS.md §"Review CLP" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Review — CLP (TERMS / POINTER / EXIT) | `skills/cds/CDS.md §"Review CLP" → "CLP form"` (or `skills/cds/review/SKILL.md §"CLP"` if expanded) | Sub 5 | Verbatim move; the CLP form is software-class (research review uses a different CLP shape). |
| §Review — reviewer ask list (α/β/γ scores, weakest-axis diagnosis, concrete patch suggestions, iterate-or-converge verdict) | `skills/cds/CDS.md §"Review CLP" → "Reviewer ask list"` (or `skills/cds/review/SKILL.md §"Asks"`) | Sub 5 | Verbatim move; the ask list is the β oracle's structured output. |

---

## 7. Gate + closure verification (§Gate) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Gate (§Gate)" (line 134 at source pin).
**Source content:** release-readiness preconditions; the closure verification checklist F1–F10 (missing α/β/γ close-outs, stale `.cdd/unreleased/{N}/` after release, missing RELEASE.md, δ tag ordering, α close-out re-dispatch mechanism, PRA presence). Operational realization in `release/SKILL.md`, `gamma/SKILL.md §2.10`, `operator/SKILL.md`.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Gate"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The release-readiness preconditions and the F1–F10 closure verification checklist (F1 missing α close-out; F2 missing β close-out; F3 missing γ close-out; F4 stale `.cdd/unreleased/{N}/`; F5 missing RELEASE.md; F6 δ tag ordering violation; F7 missing α close-out re-dispatch; F8 missing PRA; F9 missing `cdd-iteration.md` when triggers fired; F10 unresolved triage) now live as canonical content in `CDS.md` with F1–F10 preserved as identifiable `#### F{N}:` sub-headings (stable cross-reference targets for Sub 6); the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/release/SKILL.md §2.5a + §2.5b + §3.8`, `gamma/SKILL.md §2.6–§2.10`, and `operator/SKILL.md §3.1` (cited from `CDS.md §"Gate" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Gate — release-readiness preconditions | `skills/cds/CDS.md §"Gate" → "Release-readiness preconditions"` (or `skills/cds/gate/SKILL.md §"Preconditions"`) | Sub 5 | Verbatim move; release is software-class (research has no "release"; it has "publication"). |
| §Gate — closure verification checklist F1–F10 | `skills/cds/CDS.md §"Gate" → "Closure verification checklist"` (or `skills/cds/gate/SKILL.md §"F1–F10"`) | Sub 5 | Verbatim move; the F1–F10 anchors are stable citation targets used by `gamma/SKILL.md` and `release/SKILL.md`. Sub 6 of [cnos#403](https://github.com/usurobor/cnos/issues/403) re-points the citations to the new CDS surface. |

---

## 8. Assessment + cycle iteration (§Assessment) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Assessment (§Assessment)" (line 135 at source pin).
**Source content:** PRA contents (measured coherence delta, encoding lag, MCA/MCI balance, process learning, review quality metrics, CDD self-coherence on α/β/γ axes, cycle iteration, next-move commitment); **§9.1 cycle iteration** triggers (review rounds > 2; mechanical ratio > 20%; avoidable tooling/environmental failure; loaded skill failed to prevent a finding); friction log; root cause classification; skill impact; MCA; **cycle-level L5/L6/L7 framework** (per `docs/gamma/ENGINEERING-LEVELS.md`). Operational realization in `post-release/SKILL.md`.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Assessment"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The PRA contents (seven sections: coherence measurement, encoding lag, process learning, review quality, CDS self-coherence, production verification, next move), the four §9.1 cycle-iteration triggers (review rounds > 2; mechanical ratio > 20%; avoidable tooling/environmental failure; loaded skill failed) preserved verbatim as anchor targets for Sub 6, the friction log (root cause classification / skill impact / MCA / disposition), and the L5/L6/L7 engineering-levels framework now live as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6 / §5.6a / §5.6b / §5.7`, `gamma/SKILL.md §2.7–§2.9`, and `docs/gamma/ENGINEERING-LEVELS.md` (cited from `CDS.md §"Assessment" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Assessment — PRA contents (coherence delta, encoding lag, MCA/MCI balance, process learning, review quality, self-coherence, cycle iteration, next-move) | `skills/cds/CDS.md §"Assessment" → "PRA contents"` (or `skills/cds/assessment/SKILL.md §"PRA"`) | Sub 5 | Verbatim move; PRA is CDS-specific (research realization has a different post-publication assessment shape). |
| §Assessment — §9.1 cycle iteration triggers (rounds > 2; mechanical ratio > 20%; avoidable tooling failure; loaded skill failed) | `skills/cds/CDS.md §"Assessment" → "Cycle iteration triggers"` (or `skills/cds/assessment/SKILL.md §"Triggers"`) | Sub 5 | Verbatim move; the triggers anchor §9.1 are stable cross-reference targets. |
| §Assessment — friction log; root cause classification; skill impact; MCA | `skills/cds/CDS.md §"Assessment" → "Friction log"` | Sub 5 | Verbatim move. |
| §Assessment — cycle-level L5/L6/L7 framework (per `docs/gamma/ENGINEERING-LEVELS.md`) | `skills/cds/CDS.md §"Assessment" → "Engineering levels"` | Sub 5 | Verbatim move with citation; the L5/L6/L7 framework is software-class. |

---

## 9. Closure (§Closure) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Closure (§Closure)" (line 136 at source pin).
**Source content:** immediate outputs (changelog corrections, missing documentation, skill/process micro-patches, skill patches identified by cycle iteration, issue filing, lag-table updates, hub-memory writes); deferred outputs (next-MCA issue number, owner, target branch name, first AC, MCI freeze/resume state); closure rule (cycle closes only when immediate outputs executed, deferred outputs committed, cycle iteration present if triggered). Operational realization in `gamma/SKILL.md`, `post-release/SKILL.md`.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Closure"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The immediate outputs (changelog corrections, missing documentation, skill/process micro-patches, skill patches from cycle iteration, issue filing, lag-table updates, hub-memory writes), the deferred outputs (next-MCA issue number, owner, target branch name, first AC, MCI freeze/resume state), and the Closure rule (three-conjunct condition: immediate outputs executed AND deferred outputs committed AND cycle iteration present if triggered) now live as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/gamma/SKILL.md §2.7–§2.10`, `§3.6 + §3.7`, and `post-release/SKILL.md §6 + §7` (cited from `CDS.md §"Closure" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Closure — immediate outputs (changelog, docs, micro-patches, skill patches, issue filing, lag-table updates, hub-memory writes) | `skills/cds/CDS.md §"Closure" → "Immediate outputs"` | Sub 5 | Verbatim move; outputs are software-class (changelog = CHANGELOG.md; hub-memory = δ-the-operator's substrate). |
| §Closure — deferred outputs (next-MCA issue number, owner, target branch name, first AC, MCI freeze/resume state) | `skills/cds/CDS.md §"Closure" → "Deferred outputs"` | Sub 5 | Verbatim move. |
| §Closure — closure rule (cycle closes only when …) | `skills/cds/CDS.md §"Closure" → "Closure rule"` | Sub 5 | Verbatim move; the closure rule is the cycle-level coherence guarantee. |

---

## 10. Retro-packaging rule (§Retro-packaging) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Retro-packaging rule (§Retro-packaging)" (line 137 at source pin).
**Source content:** direct-to-main exception handling (retro-snapshot in version directory, self-coherence covering the change, version-history entry). Operational realization in `release/SKILL.md`.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Retro-packaging"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The direct-to-main exception (retro-snapshot in version directory; retroactive self-coherence covering the change with `[retroactive]` marker; version-history entry with honest scoring under configuration-floor caps) now lives as canonical content in `CDS.md`; the v0.1 operational overlay continues to live in `cnos.cdd/skills/cdd/release/SKILL.md §2.5b + §3.8` (cited from `CDS.md §"Retro-packaging" → "Operational realization"`) until the v1 CDS-side role rewrite.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Retro-packaging — direct-to-main exception (retro-snapshot, self-coherence, version-history entry) | `skills/cds/CDS.md §"Retro-packaging" → "Direct-to-main exception"` | Sub 5 | Verbatim move; the retro-packaging rule is software-class (research realization has no "direct-to-main" — publications are not branches). |

---

## 11. Non-goals (§Non-goals) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Non-goals (§Non-goals)" (line 138 at source pin).
**Source content:** software-cycle non-goals — do not optimize primarily for speed; do not treat issue queues as self-justifying; do not reduce review to local diff reading; do not treat release as "tag and hope"; do not confuse a shipped feature with a closed coherence cycle.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403).
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Non-goals" → "Software-cycle non-goals"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The 5 software-cycle non-goals (speed; issue queue self-justification; review reduction to local diff; release as tag-and-hope; shipped feature ≠ closed cycle) now live as canonical content in `CDS.md` under the `### Software-cycle non-goals` sub-heading (the existing `## Non-goals` was split into "Sub-level non-goals" — scope-discipline of Subs 2–5 — and "Software-cycle non-goals" — the 5 protocol-level items). No separate operational overlay; the non-goals are doctrine, not procedure.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Non-goals — software-cycle non-goals (5 items) | `skills/cds/CDS.md §"Non-goals"` | Sub 5 | Verbatim move; these are CDS-specific non-goals (CDR's non-goals differ; both inherit kernel-class non-goals from `cnos.cdd/skills/cdd/CDD.md §"Non-goals"`). |

---

## 12. Large-file authoring rule (§Large-file) — Sub 5

**Source family:** [`CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" bullet "Large-file authoring rule (§Large-file)" (line 139 at source pin).
**Source content:** files > 50 lines written section-by-section to disk with the section-manifest HTML-comment header and the resumption protocol. Operational realization in this section's own usage pattern.
**Migration sub:** Sub 5 of [cnos#403](https://github.com/usurobor/cnos/issues/403). (Folded into §Mechanical / §Artifacts family; the large-file rule is operational discipline that supports both.)
**Status:** **v0.1 migrated; canonical at [`CDS.md §"Large-file authoring rule"`](../skills/cds/CDS.md).** Landed via [cnos#410](https://github.com/usurobor/cnos/issues/410) (Sub 5, B-lite thin extract). The 50-line threshold, the section-manifest HTML-comment header (`<!-- sections: [...] -->` + `<!-- completed: [...] -->`), the resumption protocol (6-step), and the anti-patterns now live as canonical content in `CDS.md`. The rule is **self-referential** — its operational realization is its own usage pattern (every file > 50 lines in `cnos.cds`, `cnos.cdd`, and `cnos.cdr` follows the rule; `CDS.md` is the canonical exemplar). No separate v0.1 overlay; the rule's authority is the §Large-file section's body.

| Source (CDD.md §-anchor or family) | Destination (path under cnos.cds/) | Sub | Note |
|---|---|---|---|
| §Large-file — files > 50 lines: section-by-section write + section-manifest header + resumption protocol | `skills/cds/CDS.md §"Large-file authoring rule"` (or `skills/cds/large-file/SKILL.md` if expanded; pointer to `cnos.cdd` operational realization acceptable in v0.1) | Sub 5 | Verbatim move; the rule is software-class operational discipline (research-class CDR may inherit a similar rule from the generic substrate, but its citation site differs). |

---

## 13. Coverage verification

Every "pending cds extraction" marker in [`cnos.cdd/skills/cdd/CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) §"Software-specific realization" (lines 122–141 at source pin) is represented in the tables above:

| CDD.md marker | Covered in table § |
|---|---|
| §Inputs (line 126) | §1 (folded into Selection) |
| §Selection (line 127) | §1 (Selection function) |
| §Lifecycle (line 128) | §2 (Development lifecycle) |
| §Roles (line 129) | §2 (folded into Development lifecycle) |
| §Tracking (line 130) | §3 (Coordination surfaces) |
| §Artifacts (line 131) | §4 (Artifact contract) |
| §Mechanical (line 132) | §5 (Mechanical vs judgment boundary) |
| §Review (line 133) | §6 (Review CLP) |
| §Gate (line 134) | §7 (Gate + closure verification) |
| §Assessment (line 135) | §8 (Assessment + cycle iteration) |
| §Closure (line 136) | §9 (Closure) |
| §Retro-packaging (line 137) | §10 (Retro-packaging rule) |
| §Non-goals (line 138) | §11 (Non-goals) |
| §Large-file (line 139) | §12 (Large-file authoring rule) |

Sub 6 of [cnos#403](https://github.com/usurobor/cnos/issues/403) will, on completion of Subs 3–5, sweep the "pending cds extraction" markers in `CDD.md` against this map: each marker either resolves (content migrated) or updates (citation re-pointed at the CDS surface named above). No marker is silently dropped.

## 14. Open questions

These are noted for the operator dispatching Subs 3–5 against this map; none are blockers for Sub 1.

- **`.cdd/` → `.cds/` filesystem migration.** The Artifact Location Matrix (§4) names paths under `.cdd/unreleased/{N}/`. Subs 3–5 may either retain `.cdd/` paths (citing them as "the current project binding's CDD-style cycle directory") or plan a `.cds/`-rooted rename. The rename is its own coordination problem (existing `.cdd/unreleased/` directories must continue to validate against the post-rename schema). Recommend: Sub 4 keeps `.cdd/` paths and notes the planned rename as a separate post-403 cycle.
- **`cross-repo/SKILL.md` location.** Currently in `cnos.cdd/skills/cdd/cross-repo/SKILL.md`. The cross-repo proposal lifecycle is software-class (it coordinates software-cycle proposals across repos); however, the generic mechanism (proposal STATUS state machine) is arguably kernel-level. Sub 4 may either move the skill to `cnos.cds/skills/cds/cross-repo/SKILL.md` (with `cnos.cdd` retaining a pointer) or leave it in `cnos.cdd` as a generic substrate. Recommend: leave in `cnos.cdd` for v0.1; revisit at Sub 6.
- **`harness/SKILL.md` location.** Same situation as `cross-repo/SKILL.md`. The harness coordinates software-cycle dispatch but its primitives (polling, branch creation, session lifecycle) are arguably generic. Recommend: leave in `cnos.cdd` for v0.1; cite from `CDS.md §"Coordination surfaces"` rather than duplicate.
- **`operator/SKILL.md` location.** Same situation. The δ-the-operator session-routing contract is software-class operational; the generic δ doctrine lives in `ROLES.md §4a`. Recommend: leave in `cnos.cdd` for v0.1; CDS instantiates δ via `skills/cds/delta/SKILL.md` (pointer overlay acceptable).
- **`release-effector/SKILL.md` and `release/SKILL.md` location.** These are software-class (research-class CDR does not "release"; it publishes). They are strong candidates to move into `cnos.cds/skills/cds/release/SKILL.md`. Recommend: Sub 5 moves them.

These are migration-coordination questions, not destination-uncertainty questions. The destinations named in the per-surface tables above are stable commitments for Subs 3–5; the open questions are about whether *operational realizations* (currently in role/runtime-substrate SKILLs under `cnos.cdd`) also move, or stay and are cited from the new CDS doctrine surface.

---

**End of extraction map.**
