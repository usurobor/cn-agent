# α design notes — cycle/396

**Issue:** [cnos#396](https://github.com/usurobor/cnos/issues/396).
**Scope:** Author `src/packages/cnos.cdr/docs/empirical-anchor-cph.md`.
**γ scaffold:** [`gamma-scaffold.md`](./gamma-scaffold.md) (binding).

## Per-artifact-class mapping (α-frozen pre-authoring)

Each row names the cph artifact class, its cph evidence path (commit-pinned to `2bf73ad64b7aa338cc6c8bc48067fbbba818cffa`), and the cnos.cdr v0.1 surface element it maps onto. The CDR.md surface citations target sections of `src/packages/cnos.cdr/skills/cdr/CDR.md` (Sub 1, merged at c0eb9656).

### Class 1 — `PROJECT.md` (project posture)

- **cph evidence:** `PROJECT.md` (top-level). Snapshot of current stage, blocker, next action, last field report.
- **CDR surface:** `CDR.md §"Persona, Protocol, Project" → Layer 3 — Project binding (canonical home: <project>/.cdr/ or <project>/cdr/)`. Project posture is project-binding content; CDR.md explicitly relegates it to the project layer.
- **Mechanism:** Project-binding layer 4 of the five-layer enforcement chain (`ROLES.md §4a`); CDR.md does not author project posture; cph owns it.

### Class 2 — Project-local CDR doctrine (cph has NO `CDR.md` at top-level)

- **cph evidence:** No top-level `CDR.md` in cph (verified: HTTP 404 on `https://raw.githubusercontent.com/usurobor/cph/main/CDR.md`). cph's CDR doctrine is enacted via the `.cdr/` directory (Stream-B α + β SKILL drafts per cph#32) and `.cdr/coherence-log.md`. The "predecessor of cnos.cdr/skills/cdr/CDR.md" referenced in cnos#396 issue body is a hypothetical artifact — the actual cph realisation is structural (a `.cdr/` directory and Stream-B SKILL drafts), not a single doctrinal file.
- **CDR surface:** `CDR.md §"Empirical anchor"` — anchors cnos.cdr to cph's `.cdr/` directory; project-local doctrinal extensions are project-binding concerns.
- **Mechanism:** The doctrine has been **uplifted** from cph to cnos.cdr by this whole cycle's parent (#376); the absence of `CDR.md` in cph at top-level is the post-uplift state. Map to the §Empirical anchor section (which names cph as the anchor and `.cdr/` as the project binding).

### Class 3 — `ROADMAP.md` (research roadmap)

- **cph evidence:** `ROADMAP.md` (top-level). Phases R0–R6 with gate verdicts per phase.
- **CDR surface:** `CDR.md §"Field 4 — δ cadence"` (wave/phase sequencing under gate-transition cadence) + `CDR.md §"Layer 3 — Project binding"` (project-local roadmap).
- **Mechanism:** δ cadence is gate-transition-shaped; cph's phase-shaped roadmap is the project-binding realisation. Each R-phase corresponds to a δ-opened wave sequence; phase verdicts (GO/partial GO/etc.) are Field 3 receipt verdicts at the phase boundary.

### Class 4 — Wave manifest (`.cdr/waves/<wave-id>/manifest.md`)

- **cph evidence:** `.cdr/waves/zeroth-pilot-2026-05-15/manifest.md` (sampled). Carries: date, dispatcher, repo, parent-issue, sub-issue dependency chain, standing permissions, timeout budgets, known constraints.
- **CDR surface:** `CDR.md §"Field 4 — δ cadence" → Wave open` (δ selects gap + dispatches α+β) + receipt-schema generic-kernel field `protocol_id` (the manifest declares which protocol governs the wave).
- **Mechanism:** Manifest is the wave-open artifact; δ writes it; α+β operate within its bounds. Maps to the Wave-open sub-state of Field 4.

### Class 5 — Wave receipt (`.cdr/waves/<wave-id>/receipt.md`)

- **cph evidence:** `.cdr/waves/zeroth-pilot-2026-05-15/receipt.md`. Records WAVE OUTCOME, PER-ISSUE SUMMARY, BRANCHES & TAGS, ESCALATIONS, KNOWN GAPS, embeds the wave-closeout decision. Verdict: REVISE.
- **CDR surface:** `CDR.md §"Field 3 — γ close-out artifact"` → `schemas/cdr/receipt.cue` `#CDRReceipt` (typed γ artifact). Receipt fields map onto `claim_refs`/`data_refs`/`method_refs`/`result_refs`/`claim_status`/`limitations`/`reproduction` + generic-kernel `boundary_decision`/`verdict`/`protocol_gap_refs`.
- **Mechanism:** The wave receipt is the parent-facing typed receipt at the wave boundary; it carries the gate verdict (Field 3 vocabulary).

### Class 6 — Wave closeout (`.cdr/waves/<wave-id>/wave-closeout.md`)

- **cph evidence:** `.cdr/waves/zeroth-pilot-2026-05-15/wave-closeout.md`. Records: Final wave decision (REVISE), per-issue verdicts, wave-level findings (Pattern 1–4: cdd-*-gap patterns), authority bounds, escalations, known wave-level gaps.
- **CDR surface:** `CDR.md §"Field 5 — ε iteration cadence"` (the "wave-level findings (cdd-*-gap patterns across cycles)" section in the closeout is structurally identical to ε's receipt-stream protocol-gap signal). Also `CDR.md §"Field 3"` for the final-decision verdict line.
- **Mechanism:** Wave-closeout is dual-loaded — it carries Field 3's final verdict AND Field 5's cross-cycle protocol-gap findings. The "Pattern N" sub-sections in cph's closeout map directly onto ε's trigger classes.

### Class 7 — Field reports (`reports/field-report-*.md`)

- **cph evidence:** `reports/field-report-01-existing-data-zeroth-pilot.md` (and -00, -02, -03, -04). Each carries: claim, decision gate, status labels (✓/✗/PASS/FAIL/Testable/Blocked), limitations, data inventory, segmentation method revisions, hypothesis status, falsification assessment.
- **CDR surface:** `CDR.md §"Field 1 — Matter type" → Report` (synthesis surface: "field notes, wave reports, project briefs"). Field 1 explicitly names "field notes" as a report sub-class. Reports cite the typed receipt (`#CDRReceipt`) that backs them.
- **Mechanism:** Report is α-matter (Field 1). Status labels (PASS/FAIL/Testable/Blocked) map onto `claim_status` enum (`observed | computed | inferred | hypothesized | indeterminate`) and onto receipt `limitations`.

### Class 8 — Dataset documentation (`data/external/opencap-lab-validation.md`)

- **cph evidence:** `data/external/opencap-lab-validation.md`. Carries: dataset name, source (SimTK 2385/6688), archive path, SHA-256, mount paths, license, data-use policy, provenance, version pinning, trial inventory.
- **CDR surface:** `CDR.md §"Field 1 — Matter type" → Dataset` ("a body of observations or measurements with a manifest and provenance: mount point, checksum, source attribution, data-use policy compliance"). `#CDRReceipt.data_refs` typed field.
- **Mechanism:** Dataset manifest is the per-`data_refs`-entry artifact; receipt cites the manifest path.

### Class 9 — Diagnostic oracles (`scripts/segmentation_diagnostics.py` etc.)

- **cph evidence:** `scripts/segmentation_diagnostics.py`, `scripts/segmentation.py`, `scripts/segmentation_contralateral.py`, `scripts/features.py`, `scripts/comparison.py`, `scripts/io_opencap.py`, `scripts/measure-coherence.sh`.
- **CDR surface:** `CDR.md §"Field 2 — Review oracle" → Diagnostic oracles` ("where claims rest on measurements, the measurement procedure has a known failure-detection step: calibration check, signal-to-noise floor, control comparison, sentinel test. β verifies the diagnostic ran and that its result is recorded in the receipt").
- **Mechanism:** Diagnostic scripts are the β-oracle implementation; `#CDRReceipt.method_refs` cites the script + commit SHA.

### Class 10 — Gate decisions (verdict vocabulary)

- **cph evidence:** Verdicts observed across cph: **GO**, **partial GO**, **indeterminate**, **survives-subject-to-caveats** (`ROADMAP.md` §"R4"), **REVISE** (`zeroth-pilot-2026-05-15/wave-closeout.md` final decision), **NO-GO** (in `PROJECT.md` rationales), **triggered/NOT-triggered** (falsification table in field-report-01).
- **CDR surface:** `CDR.md §"Field 3" → "Gate verdict vocabulary"` — declares the vocabulary GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO with explicit mapping rules (cph's "partial GO" → BOUNDED-GO; cph's "construct-survives-subject-to-caveats" → BOUNDED-GO with bound in `limitations`).
- **Mechanism:** Each cph verdict has a CDR vocabulary equivalent; the mapping is enumerated in CDR.md §"Empirical anchor" → "Shape-compatibility claim". The mapping doc spells out the full cross-walk in §Verdict-vocabulary cross-walk.

### Class 11 — Dataset hash pinning

- **cph evidence:** `data/external/opencap-lab-validation.md` carries `SHA-256: 3290d485124fd12c85dd3bc9ee851f3a0530ad0ff58bc396973e665dd6d28187` for the OpenCap archive; verified via `unzip -t`.
- **CDR surface:** `CDR.md §"Field 1 — Matter type" → Dataset` ("manifest and provenance: mount point, **checksum**, source attribution, data-use policy compliance") + `CDR.md §"Field 2 — Review oracle" → Data-policy compliance` ("β rejects a receipt whose data refs ... are under-specified (**no manifest, no checksum**)"). Schema field: `#CDRReceipt.data_refs[i].checksum` (per schemas/cdr/receipt.cue).
- **Mechanism:** Hash pinning is the structural mechanism that the data-policy oracle enforces. cph's `SHA-256:` declaration in the manifest is what β reads when verifying data-policy compliance.

### Class 12 — Multi-agent attribution patterns

- **cph evidence:** cph#32 (Stream-A/Stream-B split rationale): "Single persona set (α/β/γ/δ/ε) conflated code-authoring cycles with field-report-authoring cycles, blurring accountability"; `coherence-log.md` records per-version finding attribution; `wave-closeout.md` records per-issue verdict authorship.
- **CDR surface:** `CDR.md §"Field 6 — Actor collapse rule"` — declares α=β NEVER permitted for research claims; γ=δ permitted with conditions; ε=δ permitted until receipt-stream volume warrants split. **`CDR.md §"Field 6" → "The collapse signal is observable"** explicitly names the cph#32 failure mode: "If a cycle's receipt cannot distinguish α and β authorship (single signature, single session, no review record), the cycle is structurally invalid for research-class claim transmission."
- **Mechanism:** Field 6's actor-collapse-rule is the doctrinal answer to cph's discovered need for Stream-A/Stream-B separation; the cnos.cdr-side overlay (Sub 3, not yet merged) will declare the Stream-A/Stream-B split as a role-overlay procedure.

### Class 13 — cph#32 data-mounted-gate failure pattern

- **cph evidence:** cph#32 issue body: "α-as-agent authored field reports citing numerical results ('60/60 trials, 7 BH-significant features') without the required data mount present"; "Inference smuggled into measurement language without boundary marking; reports paraphrased rather than quoted script outputs."
- **CDR surface:** `CDR.md §"Field 5 — ε iteration cadence" → "Missing data gates"` ("receipts shipping `observed` or `computed` claims with under-specified `data_refs` (no manifest, no checksum, no data-use compliance record)") AND `CDR.md §"Field 5" → "Overclaiming"` ("receipts whose `claim_status: observed` evidence supports only `inferred`") AND `CDR.md §"Field 5" → "Unreproducible numbers"` ("receipts with `reproduction.output_match: false`, or receipts that should have a reproduction record... but lack one").
- **Mechanism:** cph#32 is the empirical motivator for three of Field 5's six trigger classes. The cnos.cdr v0.1 surface absorbs the failure pattern as ε-trigger-classes that surface across receipt streams.

### Extension class 14 — Coherence log (`.cdr/coherence-log.md`)

- **cph evidence:** `.cdr/coherence-log.md`. Per-version coherence assessment; tracks cycle IDs, trigger classes, findings, protocol gaps across cph#16, #21, #26, #27, #28; records Stream-A scaffolding caveat ("empirical claims pending Stream-B validation; OpenCap archive not mounted in dispatch container for versions 0.3.0–0.4.1").
- **CDR surface:** `CDR.md §"Field 5 — ε iteration cadence"` — receipt-stream review surface. The coherence-log is structurally the cross-receipt ε surface for cph.
- **Mechanism:** ε reads across many `#CDRReceipt` instances and surfaces protocol-gap patterns; cph's `coherence-log.md` is exactly that surface, project-bound. cnos.cdr's ε role-overlay (Sub 3, not yet merged) declares the doctrinal surface; the project-local realisation is cph's coherence-log.

### Extension class 15 — Intra-wave status tracking (`.cdr/waves/<wave-id>/status.md`)

- **cph evidence:** `.cdr/waves/zeroth-pilot-2026-05-15/status.md`. Intra-wave dashboard; tracks sub-task approval status, round counts, blockers. Distinct from manifest (wave-open) and receipt (wave-close).
- **CDR surface:** `CDR.md §"Field 4 — δ cadence" → "Wave in progress"` ("α produces matter (Field 1); β reviews against the review oracle (Field 2); γ coordinates the loop"). The status.md file is γ's in-flight coordination surface.
- **Mechanism:** γ coordinates the in-progress wave; status.md is the project-local realisation of γ's live wave dashboard. CDR.md Field 4's wave-in-progress sub-state is the doctrinal surface.

## Surface-coverage check (pre-authoring)

Every class 1–15 maps to a concrete CDR.md section, a `#CDRReceipt` schema field, or both. No "no surface" rows. AC3 will pass.

## Citation discipline (AC5)

Each row in the final doc cites cph evidence as `<path>@<sha>` where the SHA is `2bf73ad64b7aa338cc6c8bc48067fbbba818cffa` (cph main as of cycle filing). No cph-local prose is reproduced beyond a 1-line identifying snippet per row. Section headings can be cited verbatim; multi-line paragraphs cannot.

## Verdict vocabulary cross-walk (pre-authoring)

| cph verdict (project) | CDR.md verdict (protocol) | Mapping notes |
|------------------------|----------------------------|---------------|
| GO | GO | direct |
| partial GO | BOUNDED-GO | named bound goes in `limitations` |
| REVISE | REVISE | direct (Field 3 vocabulary) |
| NO-GO | NO-GO | direct |
| INDETERMINATE | INDETERMINATE | direct |
| survives-subject-to-caveats | BOUNDED-GO | bound named in `limitations` |
| triggered / NOT-triggered (falsification table) | (oracle-internal, not gate verdict) | falsification-check states; not transmissibility verdicts |

CDR.md §"Empirical anchor" → "Shape-compatibility claim" (Field 3 bullet) names exactly this mapping. The empirical-anchor doc tabulates it row-by-row for AC2 coverage.

## Surface non-goals (explicit)

This doc does NOT:

- Author CDR doctrine. Sub 1 owns CDR.md; the mapping doc cites CDR.md sections, does not extend them.
- Author role-overlays. Sub 3 owns `cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md`. References are by name; role-overlay contents not assumed merged.
- Author cph project content. cph owns `PROJECT.md`, `.cdr/`, etc.
- Embed cph prose. Citations carry path + SHA; section headings only.
- Propose CDR surface extensions. Extension classes 14 + 15 map onto **existing** CDR.md surfaces; they are not new surfaces.

## CDD Trace

| Step | Artifact | Decision |
|------|----------|----------|
| 0 Observe | γ scaffold + cph survey | 15 classes (13 issue + 2 extensions) all have CDR surface |
| 1 Select | cnos#376 AC3 | Sub 4 ships mapping doc; this is the design |
| 4 Gap | unsatisfied AC3 | 15 mappings to author |
| 5 Mode | docs-only collapse | proceed to authoring phase |
| 6 Artifacts | design-notes.md | next: empirical-anchor-cph.md |
