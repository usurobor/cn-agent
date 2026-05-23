# γ scaffold — cycle/396 (Sub 4 of #376)

**Issue:** [cnos#396](https://github.com/usurobor/cnos/issues/396) — empirical-anchor doc mapping cph CDR practice to cnos.cdr v0.1 surface.

**Mode:** docs-only; γ+α+β-collapsed-on-δ (β-α-collapse acknowledged by δ).

## Surface (single deliverable)

One new file:

- `src/packages/cnos.cdr/docs/empirical-anchor-cph.md`

Touches no existing files in the build/test path. No code, no schemas, no CI. Single Markdown deliverable.

## AC oracle approach (issue body verbatim)

| AC | Oracle (mechanical) | Surface |
|----|---------------------|---------|
| AC1 | `test -s src/packages/cnos.cdr/docs/empirical-anchor-cph.md` | the file |
| AC2 | `rg "PROJECT.md\|CDR.md\|ROADMAP.md\|wave-manifest\|wave-receipt\|wave-closeout\|field-report\|dataset\|diagnostic-oracle\|gate.*verdict\|hash.*pinning\|cph#32" src/packages/cnos.cdr/docs/empirical-anchor-cph.md` returns hits covering ≥12 classes | the file |
| AC3 | `rg "TBD\|no surface\|not mapped" src/packages/cnos.cdr/docs/empirical-anchor-cph.md` returns 0 hits (or only in clear prose-context that explicitly discusses what is NOT a CDR surface concern) | the file |
| AC4 | Closing section names `cnos#376 AC3` as satisfied; merge-time comment on cnos#376 cites this doc | the file + cnos#376 comment |
| AC5 | Doc is structured as mapping table + commentary; no extended quotations of cph-local files; citations are path + SHA only | the file |

## cph artifact-class enumeration plan

The issue body enumerates 12+ classes. Empirical survey of cph (commit `2bf73ad64b7aa338cc6c8bc48067fbbba818cffa`) confirms the enumeration and discloses two additional classes worth extending the table for (per δ's refusal-condition guidance: "extend the mapping table; note the addition"):

| # | Artifact class (issue body) | cph evidence path |
|---|------------------------------|-------------------|
| 1 | `PROJECT.md` (project posture) | `PROJECT.md` |
| 2 | `CDR.md` (project-local CDR doctrine) | **DOES NOT EXIST at cph top-level** — cph carries CDR doctrine via `.cdr/` directory (Stream-B SKILL drafts per cph#32) + `coherence-log.md`. Map to "project-local CDR enactment surface" rather than a single file. |
| 3 | `ROADMAP.md` | `ROADMAP.md` |
| 4 | Wave manifest | `.cdr/waves/<wave-id>/manifest.md` (sampled: `zeroth-pilot-2026-05-15/manifest.md`) |
| 5 | Wave receipt | `.cdr/waves/<wave-id>/receipt.md` (sampled: same wave) |
| 6 | Wave closeout | `.cdr/waves/<wave-id>/wave-closeout.md` (sampled: same wave) |
| 7 | Field reports | `reports/field-report-*.md` (5 reports present) |
| 8 | Dataset documentation | `data/external/opencap-lab-validation.md` |
| 9 | Diagnostic oracles | `scripts/segmentation_diagnostics.py` (+ `segmentation.py`, `features.py`, etc.) |
| 10 | Gate decisions (GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO; partial GO; construct-survives-subject-to-caveats) | embedded in receipt/closeout files; verdict vocabulary observed in `ROADMAP.md` and `wave-closeout.md` |
| 11 | Dataset hash pinning | `data/external/opencap-lab-validation.md` `SHA-256: 3290d485...` |
| 12 | Multi-agent attribution patterns | `coherence-log.md` + cph#32 issue body (Stream-A/Stream-B split rationale) |
| 13 | cph#32 data-mounted-gate failure pattern | cph#32 (Stream-B SKILL drafts; failure pattern named in body) |

Extensions discovered during survey (not in issue body; will be added per δ guidance):

| # | Extension artifact class | cph evidence path | CDR.md v0.1 surface | Rationale |
|---|--------------------------|-------------------|---------------------|-----------|
| 14 | Coherence log (per-version coherence assessment) | `.cdr/coherence-log.md` | Field 5 (ε iteration cadence) — receipt-stream review surface | Cross-receipt protocol-gap signal aggregation; native ε artifact |
| 15 | Intra-wave status tracking | `.cdr/waves/<wave-id>/status.md` | Field 4 (δ cadence) — wave-in-progress sub-state | Distinct from receipt (final) and closeout (post-final); names the live wave dashboard |

## WebFetch strategy

cph is outside cnos's MCP scope. Access strategy (verified working during scaffold):

- **Raw file reads:** `https://raw.githubusercontent.com/usurobor/cph/main/<path>` — works for PROJECT.md, ROADMAP.md, `.cdr/waves/*/manifest.md`, receipt.md, wave-closeout.md, status.md, coherence-log.md, reports, dataset docs.
- **Directory listings:** `https://github.com/usurobor/cph/tree/main/<dir>` (HTML; WebFetch summarises) — works for directory enumeration.
- **api.github.com endpoints:** **403 Forbidden** (no auth). Avoid.
- **Commit SHA:** `https://github.com/usurobor/cph/commits/main` — works; returns latest SHA in HTML.
- **Issues (cph#32):** `https://github.com/usurobor/cph/issues/32` — works.

Working cph commit pin: **`2bf73ad64b7aa338cc6c8bc48067fbbba818cffa`** (main, as of cycle/396 filing).

## Deliverable shape (γ-frozen)

The mapping doc will have these sections:

1. Preamble — what this demonstrates; cnos#376 AC3 obligation.
2. Method — how the mapping was derived; cph commit pin; survey window.
3. Mapping table — one row per artifact class; 15 rows (13 issue-body + 2 extensions).
4. Verdict-vocabulary cross-walk — cph gate names ↔ CDR.md Field 3 verdict vocabulary.
5. Non-surface concerns (if any) — explicit nil if none.
6. Cross-references — CDR.md sections cited, schemas/cdr/receipt.cue fields cited.
7. Closing — cnos#376 AC3 satisfied + named gaps.

## Frozen non-goals (binding on α)

- Do not author CDR doctrine (Sub 1 owns CDR.md).
- Do not author role overlays (Sub 3 owns `cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md`); references to role-overlay files are by name only — role-skill file contents are NOT assumed merged yet.
- Do not copy cph prose verbatim; citations only (path + SHA).
- Do not propose CDR-shape changes; if a cph artifact has no surface, surface it as a finding rather than invent CDR surface.

## Tier 3 skills loaded

- `cnos.cdd/skills/cdd/design` — design discipline (mostly bypassed for docs-only; the design is the mapping table itself, framed by issue body's enumeration).
- `cnos.cdd/skills/cdd/issue/proof` — AC + oracle discipline (oracle structure adopted in §AC table above).

## CDD Trace

| Step | Artifact | Decision |
|------|----------|----------|
| 0 Observe | δ's pinned implementation contract + issue body | cph is the empirical anchor; 13 issue-body classes + 2 survey-discovered extensions |
| 1 Select | cnos#376 AC3 obligation | Sub 4 ships the mapping doc |
| 4 Gap | named incoherence | "cnos.cdr v0.1 has no per-artifact-class cph mapping; AC3 unsatisfied" |
| 5 Mode | docs-only, γ+α+β-collapsed-on-δ | β-α-collapse acknowledged; matter type is docs not research-claim |
| 6 Artifacts | this scaffold | next: design-notes.md → empirical-anchor-cph.md → self-coherence → β-review |
