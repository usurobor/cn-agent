# Empirical anchor — cph CDR practice → cnos.cdr v0.1 surface mapping

**Version:** 0.1 (Sub 4 of [cnos#376](https://github.com/usurobor/cnos/issues/376))
**Date:** 2026-05-21
**Placement:** `src/packages/cnos.cdr/docs/empirical-anchor-cph.md`
**Audience:** Reviewers verifying that cnos.cdr v0.1's surface is sufficient to retroactively describe [usurobor/cph](https://github.com/usurobor/cph)'s CDR practice; future research projects binding their own `.cdr/` artifacts to cnos.cdr.
**Scope:** A surface-by-surface mapping from every load-bearing cph CDR artifact class onto its corresponding cnos.cdr v0.1 surface element (a section of [`CDR.md`](../skills/cdr/CDR.md), a role overlay file, or a typed-receipt schema field per [`schemas/cdr/receipt.cue`](../../../../schemas/cdr/receipt.cue)). Demonstrates that cnos.cdr v0.1 can host cph's current practice without contradiction. Closes [cnos#376](https://github.com/usurobor/cnos/issues/376) AC3.

> This document is a **mapping table + commentary**, not a re-statement of cph
> or of CDR.md. It does not embed cph prose; citations carry the cph file path
> and the commit SHA at which the cited content existed. Per
> [`CDR.md §"Empirical anchor" → "Detailed mapping deferred to Sub 4"`](../skills/cdr/CDR.md):
> Sub 1 ships the shape-compatibility claim; Sub 4 (this document) ships the
> surface-by-surface verification.

---

## 0. Method

### 0.1 Source pin

cph evidence cited herein refers to cph `main` at commit
**`2bf73ad64b7aa338cc6c8bc48067fbbba818cffa`** (the head of cph `main` at the
time cycle/396 was filed). Citations are written `<path>@2bf73ad6` (short SHA
sufficient for human-readable cross-reference; the full SHA above is the
authoritative pin).

### 0.2 cnos.cdr v0.1 target surfaces

Three surface classes carry the cnos.cdr v0.1 protocol overlay:

1. [`CDR.md`](../skills/cdr/CDR.md) — the doctrinal contract. Six fields per
   `ROLES.md §3`: Matter type, Review oracle, γ close-out artifact, δ cadence,
   ε iteration cadence, Actor collapse rule.
2. Role overlays at `cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md`
   (the deliverable of [cnos#395](https://github.com/usurobor/cnos/issues/395)).
   This mapping references role overlays by name; the procedural content of
   each overlay is the property of its owning Sub.
3. [`schemas/cdr/receipt.cue`](../../../../schemas/cdr/receipt.cue) —
   `#CDRReceipt` and its typed fields (`claim_refs`, `data_refs`,
   `method_refs`, `result_refs`, `claim_status`, `limitations`,
   `reproduction`, generic-kernel `protocol_id` / `boundary_decision` /
   `verdict` / `protocol_gap_refs`).

### 0.3 Survey window

cph was surveyed top-down: top-level files (`PROJECT.md`, `README.md`,
`ROADMAP.md`, `CHANGELOG.md`); `.cdr/` (`coherence-log.md`, `waves/`,
`iterations/`, `unreleased/`); a sampled wave directory
(`.cdr/waves/zeroth-pilot-2026-05-15/`: `manifest.md`, `receipt.md`,
`status.md`, `wave-closeout.md`); the `reports/` directory (five field
reports); the `data/external/` directory; the `scripts/` directory; and
[cph#32](https://github.com/usurobor/cph/issues/32) (the data-mounted-gate
failure pattern).

### 0.4 Survey-discovered extensions

The cnos#396 issue body enumerates 13 cph artifact classes. The survey
discovered two further classes that load-bear under cph's current practice;
both map onto existing cnos.cdr v0.1 surfaces (no new surface required) and
appear in this mapping as **extension rows** (rows 14 and 15). No
cph artifact class was discovered that lacks a CDR surface.

---

## 1. Mapping table

Each row carries: the cph artifact class (with evidence path and short SHA),
the cnos.cdr v0.1 surface it maps onto, and the surface mechanism.

| # | cph artifact class | cph evidence (`<path>@<sha>`) | cnos.cdr v0.1 surface | Mechanism |
|---|---------------------|-------------------------------|------------------------|-----------|
| 1 | `PROJECT.md` (project posture) | `PROJECT.md@2bf73ad6` | [`CDR.md §"Persona, Protocol, Project" → Layer 3 — Project binding`](../skills/cdr/CDR.md) | Project posture is **project-binding** content (layer 4 of the five-layer enforcement chain per `ROLES.md §4a`); CDR.md relegates it to the project layer and does not author it. |
| 2 | Project-local CDR doctrine | (no top-level `CDR.md` in cph; doctrine enacted via `.cdr/` directory structure + Stream-B SKILL drafts per cph#32) | [`CDR.md §"Empirical anchor"`](../skills/cdr/CDR.md) | cph carries no top-level `CDR.md`; its CDR doctrine is **structural** (`.cdr/` layout) and **procedural** (Stream-B α/β SKILL drafts). The post-uplift state is: cnos.cdr owns the protocol-overlay doctrine; cph owns the project binding. |
| 3 | `ROADMAP.md` (research roadmap, phases R0–R6) | `ROADMAP.md@2bf73ad6` | [`CDR.md §"Field 4 — δ cadence"`](../skills/cdr/CDR.md) + [`CDR.md §"Layer 3 — Project binding"`](../skills/cdr/CDR.md) | δ cadence is **gate-transition-shaped**; cph's phase-shaped roadmap is the project-binding realisation. Each R-phase corresponds to a δ-opened wave sequence; phase verdicts are Field 3 receipt verdicts at the phase boundary. |
| 4 | Wave manifest / wave-manifest (`.cdr/waves/<wave-id>/manifest.md`) | `.cdr/waves/zeroth-pilot-2026-05-15/manifest.md@2bf73ad6` | [`CDR.md §"Field 4 — δ cadence" → "Wave open"`](../skills/cdr/CDR.md) + [`schemas/cdr/receipt.cue`](../../../../schemas/cdr/receipt.cue) generic-kernel `protocol_id` field | The wave-manifest is the **wave-open** artifact; δ writes it and dispatches α+β within its bounds (sub-issue dependency chain, standing permissions, timeout budgets). Maps to the Wave-open sub-state of Field 4. |
| 5 | Wave receipt / wave-receipt (`.cdr/waves/<wave-id>/receipt.md`) | `.cdr/waves/zeroth-pilot-2026-05-15/receipt.md@2bf73ad6` | [`CDR.md §"Field 3 — γ close-out artifact"`](../skills/cdr/CDR.md) → [`#CDRReceipt`](../../../../schemas/cdr/receipt.cue) | The wave-receipt is the **parent-facing typed γ artifact** at the wave boundary. Its fields project onto `#CDRReceipt.{claim_refs, data_refs, method_refs, result_refs, claim_status, limitations, reproduction}` and the generic kernel `{protocol_id, boundary_decision, verdict, protocol_gap_refs}`. |
| 6 | Wave closeout / wave-closeout (`.cdr/waves/<wave-id>/wave-closeout.md`) | `.cdr/waves/zeroth-pilot-2026-05-15/wave-closeout.md@2bf73ad6` | Dual-loaded: [`CDR.md §"Field 3"`](../skills/cdr/CDR.md) (final decision) + [`CDR.md §"Field 5 — ε iteration cadence"`](../skills/cdr/CDR.md) (wave-level findings) | The wave-closeout's "Final wave decision" carries Field 3's verdict; the closeout's "Wave-level findings (cdd-*-gap patterns across cycles)" section (Pattern 1, Pattern 2, ...) is structurally identical to ε's receipt-stream protocol-gap signal. Each pattern row populates `protocol_gap_refs` in subsequent receipts. |
| 7 | Field reports (`reports/field-report-*.md`) | `reports/field-report-01-existing-data-zeroth-pilot.md@2bf73ad6` (representative; `-00` through `-04` present) | [`CDR.md §"Field 1 — Matter type" → "Report"`](../skills/cdr/CDR.md) | Field 1 names "field notes, wave reports, project briefs" as report sub-classes; reports cite the typed receipt (`#CDRReceipt`) that backs them. Status labels (PASS/FAIL/Testable/Blocked, observed/inferred) project onto `claim_status ∈ {observed, computed, inferred, hypothesized, indeterminate}` and onto receipt `limitations`. |
| 8 | Dataset documentation (`data/external/<dataset>.md`) | `data/external/opencap-lab-validation.md@2bf73ad6` | [`CDR.md §"Field 1 — Matter type" → "Dataset"`](../skills/cdr/CDR.md) → [`#CDRReceipt.data_refs`](../../../../schemas/cdr/receipt.cue) | Field 1 declares Dataset as a body of observations "with a manifest and provenance: mount point, checksum, source attribution, data-use policy compliance" — exactly the columns cph's dataset doc carries. Each receipt cites the manifest path under `data_refs`. |
| 9 | Diagnostic oracles / diagnostic-oracle scripts (`scripts/segmentation_diagnostics.py` etc.) | `scripts/segmentation_diagnostics.py@2bf73ad6` (+ `segmentation.py`, `segmentation_contralateral.py`, `features.py`, `comparison.py`, `io_opencap.py`, `measure-coherence.sh`) | [`CDR.md §"Field 2 — Review oracle" → "Diagnostic oracles"`](../skills/cdr/CDR.md) → [`#CDRReceipt.method_refs`](../../../../schemas/cdr/receipt.cue) | Field 2's "Diagnostic oracles" bullet names "calibration check, signal-to-noise floor, control comparison, sentinel test" — the function class cph's `scripts/segmentation_diagnostics.py` (the canonical diagnostic-oracle script) implements. β verifies the diagnostic-oracle ran and that its result is recorded in the receipt; `method_refs` cites the script + commit SHA. |
| 10 | Gate-verdict vocabulary (GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO; partial GO; construct-survives-subject-to-caveats) | `ROADMAP.md@2bf73ad6` (verdicts per phase); `.cdr/waves/*/wave-closeout.md@2bf73ad6` (verdict at wave boundary); `PROJECT.md@2bf73ad6` (NO-GO rationales) | [`CDR.md §"Field 3" → "Gate verdict vocabulary"`](../skills/cdr/CDR.md) → [`#CDRReceipt.boundary_decision` + `transmissibility`](../../../../schemas/cdr/receipt.cue) | Full cross-walk in §2 below. Every cph verdict has an explicit cnos.cdr verdict mapping; no vocabulary extension required. |
| 11 | Dataset hash pinning (SHA-256 in dataset manifest) | `data/external/opencap-lab-validation.md@2bf73ad6` declares `SHA-256: 3290d485...` for the OpenCap archive (full hash: `3290d485124fd12c85dd3bc9ee851f3a0530ad0ff58bc396973e665dd6d28187`) | [`CDR.md §"Field 1 → Dataset"`](../skills/cdr/CDR.md) (the "checksum" requirement) + [`CDR.md §"Field 2 → Data-policy compliance"`](../skills/cdr/CDR.md) (β rejects receipts whose data refs are under-specified — "no manifest, no checksum") | Hash pinning is the structural mechanism the data-policy oracle enforces. cph's `SHA-256:` declaration in the manifest is the field β reads when verifying data-policy compliance; schema-side it appears in `data_refs[i].checksum`. |
| 12 | Multi-agent attribution patterns (Stream-A / Stream-B persona separation) | cph#32 (Stream-B SKILL drafts rationale); `.cdr/coherence-log.md@2bf73ad6` (per-version finding attribution); `.cdr/waves/*/wave-closeout.md@2bf73ad6` (per-issue verdict authorship) | [`CDR.md §"Field 6 — Actor collapse rule"`](../skills/cdr/CDR.md) — explicitly names α=β as **never permitted** for research claims; "The collapse signal is observable" sub-section names the cph#32 failure mode | Field 6's actor-collapse-rule is the doctrinal answer to cph's discovered need for Stream-A/Stream-B separation. The cnos.cdr role-overlay deliverable (Sub 3, cnos#395) declares the Stream-A/Stream-B procedural split in `cnos.cdr/skills/cdr/{alpha,beta}/SKILL.md`. |
| 13 | cph#32 data-mounted-gate failure pattern (numbers cited without `/opt/gait-data/` mount) | [cph#32](https://github.com/usurobor/cph/issues/32) issue body | Triple-mapped onto [`CDR.md §"Field 5"`](../skills/cdr/CDR.md): "Missing data gates" (under-specified `data_refs`), "Overclaiming" (`claim_status: observed` evidence supports only `inferred`), "Unreproducible numbers" (missing or false `reproduction.output_match`) | cph#32 is the **empirical motivator** for three of Field 5's six trigger classes. The cnos.cdr surface absorbs the failure pattern as ε-trigger-classes that surface across receipt streams; the project-binding side (cph) responds with the Stream-A/Stream-B persona split. |
| 14 | **(extension)** Coherence log (`.cdr/coherence-log.md`) | `.cdr/coherence-log.md@2bf73ad6` (per-version coherence assessment; tracks cycle IDs cph#16/#21/#26/#27/#28, trigger classes, findings, protocol gaps; records the Stream-A scaffolding caveat for versions 0.3.0–0.4.1) | [`CDR.md §"Field 5 — ε iteration cadence"`](../skills/cdr/CDR.md) — receipt-stream review surface | ε reads across many `#CDRReceipt` instances and surfaces protocol-gap patterns; cph's `coherence-log.md` is exactly that surface, project-bound. The cnos.cdr ε role-overlay (Sub 3) declares the doctrinal surface; cph's coherence-log is the project-local realisation. |
| 15 | **(extension)** Intra-wave status tracking (`.cdr/waves/<wave-id>/status.md`) | `.cdr/waves/zeroth-pilot-2026-05-15/status.md@2bf73ad6` (intra-wave dashboard; tracks sub-task approval status, round counts, blockers; distinct from manifest and receipt) | [`CDR.md §"Field 4 — δ cadence" → "Wave in progress"`](../skills/cdr/CDR.md) — γ coordinates the loop | γ's in-flight wave-coordination surface. CDR.md Field 4's wave-in-progress sub-state is the doctrinal home; cph's `status.md` is the project-local realisation. (Distinguished from the receipt because receipt is the wave-close artifact; status.md is the wave-in-progress dashboard.) |

---

## 2. Verdict-vocabulary cross-walk

cph's observed gate verdicts and their cnos.cdr v0.1 mapping (per
[`CDR.md §"Field 3" → "Gate verdict vocabulary"`](../skills/cdr/CDR.md) and
[`CDR.md §"Empirical anchor" → "Shape-compatibility claim"`](../skills/cdr/CDR.md)):

| cph verdict (project naming) | cnos.cdr v0.1 verdict (protocol naming) | Mapping mechanism | cph evidence |
|------------------------------|------------------------------------------|-------------------|--------------|
| GO | GO | direct; `boundary_decision: accept` + `transmissibility: accepted` | `ROADMAP.md@2bf73ad6` (R0, R2 closures) |
| partial GO | BOUNDED-GO | named bound enumerated in `#CDRReceipt.limitations`; `transmissibility: degraded` | `ROADMAP.md@2bf73ad6` (R3 "Partial GO on R-side"); `reports/field-report-01-existing-data-zeroth-pilot.md@2bf73ad6` ("GO with bounded scope (path (a) inferred bilateral)") |
| REVISE | REVISE | direct; `action: repair_dispatch` (β returns the receipt for revision) | `.cdr/waves/zeroth-pilot-2026-05-15/wave-closeout.md@2bf73ad6` ("REVISE for friend pre-pilot") |
| NO-GO | NO-GO | direct; `action: reject` | `PROJECT.md@2bf73ad6` (NO-GO rationale form: "close Sub C as NO-GO with rationale per protocol") |
| INDETERMINATE | INDETERMINATE | direct; `claim_status: indeterminate`; receipt records what would change the determination | `ROADMAP.md@2bf73ad6` (R3 "indeterminate on bilateral inferred-bilateral surface") |
| construct-survives-subject-to-caveats | BOUNDED-GO | bound named in `limitations`; same mapping as "partial GO" but expressed as construct-survives | `ROADMAP.md@2bf73ad6` (R4 "Construct survives R4 contact subject to anchor caveats") |
| triggered / NOT-triggered (falsification table) | (oracle-internal state, not a transmissibility verdict) | These are **falsification-check states** internal to the diagnostic oracle (Field 2); they populate `#CDRReceipt.method_refs.result` rather than `boundary_decision` | `reports/field-report-01-existing-data-zeroth-pilot.md@2bf73ad6` (6-condition falsification table) |

No cph verdict requires a cnos.cdr vocabulary extension. The shape-compatibility
claim recorded in [`CDR.md §"Empirical anchor" → "Shape-compatibility claim"`
Field 3 bullet](../skills/cdr/CDR.md) is verified row-by-row.

---

## 3. Mechanism-by-CDR-field summary

Reading the mapping by cnos.cdr v0.1 field (a transpose of §1):

- **Field 1 — Matter type** is exercised by cph rows 7 (Reports), 8 (Datasets),
  9 (Methods via diagnostic oracles), and the matter-type-vocabulary terms
  (Claim, Hypothesis, Analysis) embedded in row 7's report structure.
- **Field 2 — Review oracle** is exercised by cph rows 9 (Diagnostic oracles)
  and 11 (Data-policy compliance via hash pinning). Citation integrity is
  exercised across cph's report-cites-receipt-cites-manifest chain.
- **Field 3 — γ close-out artifact** is exercised by cph rows 5 (wave receipt
  → `#CDRReceipt`), 6 (wave closeout → final decision), and 10
  (verdict vocabulary).
- **Field 4 — δ cadence** is exercised by cph rows 3 (ROADMAP phases as wave
  sequences), 4 (wave manifest as wave-open), and 15 (status.md as
  wave-in-progress).
- **Field 5 — ε iteration cadence** is exercised by cph rows 6 (closeout's
  wave-level findings), 13 (cph#32 failure-pattern source), and 14
  (coherence-log as receipt-stream surface).
- **Field 6 — Actor collapse rule** is exercised by cph row 12 (multi-agent
  attribution patterns, Stream-A/Stream-B separation).
- **Persona / Protocol / Project layering** (`CDR.md §"Persona, Protocol,
  Project"`) is exercised by cph rows 1 (PROJECT.md as layer-3 project
  binding) and 2 (no cph-local CDR.md, doctrine uplifted to cnos.cdr).

All six fields plus the persona-protocol-project layering are exercised by
at least one cph artifact class. No CDR.md field is unexercised by the
mapping.

---

## 4. Non-surface concerns (explicit nil with prose-context disclosure)

This section discusses what is **not** a CDR surface concern, per AC3's prose-context allowance. None of the 15 cph artifact classes lacks a CDR surface; the following are listed only to forestall mis-mapping.

- **cph-local data paths (`/opt/gait-data/...`)** are project-binding
  artifacts, not CDR surface. CDR.md does not encode mount paths; cph's
  dataset manifest does. The mapping in row 8 covers the doctrinal home
  (Field 1 → Dataset, with checksum + mount-point); the path string itself
  is project-binding content.
- **cph-local script implementations** are project-binding; CDR's surface is
  the `method_refs` schema field (row 9), not the script source.
- **cph-local persona content (Stream-A, Stream-B persona names)** is
  persona-layer content (`cn-rho` and downstream); CDR.md Field 6 names the
  collapse-rule structure, not the persona identities. Row 12 maps the
  pattern, not the persona names themselves.
- **cph#32's specific implementation patches** (the Stream-B SKILL drafts in
  cph) are project-local realisations; the cnos.cdr role-overlay deliverable
  (Sub 3) declares the doctrinal Stream-A/Stream-B procedural split. Row 12
  + row 13 cover the mapping; the cph-local SKILL files themselves are
  outside cnos.cdr scope.

These are *not* AC3 violations — every artifact class in §1 has a concrete
surface. This section disambiguates project-binding content that the
mapping deliberately does not absorb.

---

## 5. Cross-references

This document cites:

- [`src/packages/cnos.cdr/skills/cdr/CDR.md`](../skills/cdr/CDR.md) —
  the cnos.cdr v0.1 doctrinal contract (Sub 1, [cnos#390](https://github.com/usurobor/cnos/issues/390)). All Field N citations target this file. The "Shape-compatibility claim" sub-section of §"Empirical anchor" is the seed that this document verifies row-by-row.
- [`schemas/cdr/receipt.cue`](../../../../schemas/cdr/receipt.cue) —
  the `#CDRReceipt` schema. Schema field citations target this file.
- [`ROLES.md`](../../../../ROLES.md) — the six-field instantiation
  contract (`§3`) and the five-layer enforcement chain (`§4a`).
- [`docs/papers/CCNF-AND-TYPED-TRUST.md §"Wave 4 — CDR bootstrap"`](../../../../docs/papers/CCNF-AND-TYPED-TRUST.md) — the doctrinal declaration that cph is anchored but not definitional.
- The role overlays at `cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` (Sub 3, [cnos#395](https://github.com/usurobor/cnos/issues/395)) — referenced by file path only; this document does not depend on their merge state.
- [usurobor/cph @ `2bf73ad6`](https://github.com/usurobor/cph) — the
  empirical anchor. Specific file references are noted in §1 and §2.
- [cph#32](https://github.com/usurobor/cph/issues/32) — the data-mounted-gate failure pattern (row 13).

---

## 6. Closing — cnos#376 AC3 satisfied

[cnos#376](https://github.com/usurobor/cnos/issues/376)'s AC3 names the
**empirical-anchor validation** obligation: cnos.cdr v0.1's design must be
demonstrated sufficient to retroactively describe cph's CDR practice. This
document discharges that obligation:

- All 13 cph artifact classes enumerated in the [cnos#396](https://github.com/usurobor/cnos/issues/396) issue body (rows 1–13) map onto a concrete cnos.cdr v0.1 surface.
- Two additional artifact classes discovered during survey (rows 14–15) also map onto existing surfaces; no new CDR-side surface is required.
- The gate-verdict vocabulary cross-walk (§2) verifies every cph verdict has a cnos.cdr equivalent.
- Each row of §1 cites cph evidence with a path + commit SHA (no verbatim cph prose embedded).
- All six CDR.md fields plus the persona/protocol/project layering are exercised by at least one cph artifact class (§3).
- No "no surface" rows; no TBD; no "not mapped" entries.

**cnos#376 AC3 is therefore satisfied** by this document at cph
commit `2bf73ad64b7aa338cc6c8bc48067fbbba818cffa`.

### 6.1 Named gaps

Per [cnos#396](https://github.com/usurobor/cnos/issues/396)'s Proof Plan ("Known gap: The mapping reflects cph's state at this cycle's filing time"), the gaps surfaced by this mapping are:

- **Forward-looking:** future cph waves may introduce artifact classes not
  surveyed here; those become follow-on extensions to §1, not
  contradictions of cnos.cdr v0.1.
- **Dependency-on-downstream:** the role-overlay references (rows 12, 14)
  target files owned by Sub 3 ([cnos#395](https://github.com/usurobor/cnos/issues/395)). When Sub 3 ships, the role-overlay file paths cited here become directly resolvable. The mapping is structurally valid regardless of Sub 3's merge timing because the references are by name, not by content dependency.
- **No CDR-side gaps:** the survey did not discover any cph artifact class
  lacking a cnos.cdr v0.1 surface. If a future cph wave introduces such a
  class, that finding requires Sub 1 ([cnos#390](https://github.com/usurobor/cnos/issues/390)) reopen, not a patch to this document.
