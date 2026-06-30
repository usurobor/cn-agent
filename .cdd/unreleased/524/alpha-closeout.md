---
cycle: 524
role: alpha
verdict: converge
date: 2026-06-30 (UTC)
authored_by: α@cdd.cnos (R0 closeout)
parent_issue: cnos#524
round: R0
beta_verdict: converge
---

# α-closeout — cnos#524 R0: wake-as-skill W0 design document

## What was implemented

This cycle delivered a single primary artifact: `.cdd/unreleased/524/w0-design.md`, a standalone W0 design document formalizing the wake-as-skill contract. No source code, schema, workflow, or golden file was touched.

The design document is a structured extraction and formal presentation of the design already present in cnos#524 (the `## W0 design (locked)` section and the operator W0-calls comment). It introduces no new design. It covers nine sections (§A through §I):

- §A: Problem statement and invariant (two-system problem; one-system resolution)
- §B: Complete `wake:` block schema with role-shaped output disjunction and field-by-field cross-reference to both `wake-provider.json` manifests
- §C: The 5 operator design decisions (nesting, artifact_class, scope, CUE/renderer boundary, body-as-prompt)
- §D: CUE ↔ renderer responsibility boundary table (15 concerns, cleanly partitioned)
- §E: Body-as-prompt rule and byte-identity oracle
- §F: W1→W4 migration plan (4 phases, oracle at each phase boundary)
- §G: W1 acceptance criteria (AC1–AC7, each with invariant and oracle)
- §H: Friction notes FN-1 through FN-8 for the W1 implementer
- §I: Explicit non-goals (8 items)

CDD support artifacts produced: `self-coherence.md §R0` and `beta-review.md §R0`.

Field coverage: all 35 fields from `agent-admin/wake-provider.json` and all 36 fields from `cds-dispatch/wake-provider.json` are accounted for in `w0-design.md §B`. Zero missing.

---

## What worked well

**Scope discipline.** The operator directive restricting this cycle to W0 design only was unambiguous, and the γ scaffold encoded it mechanically in §5 and §7. There was no temptation to begin W1 implementation — the hard scope list made it clear what was out of bounds and why. The scope guardrail confirmation in `self-coherence.md §1` was a clean pass.

**Extraction from a locked design.** Because the W0 design was already operator-ratified in the issue body, α's task was formalization, not invention. This made the work more verifiable and less prone to drift. The field-by-field cross-reference table (§B.3) was the most mechanical part of the work and produced the highest-confidence output — every row is traceable to a specific JSON key in one or both manifests.

**Role-shaped disjunction clarity.** Separating the admin output shape from the dispatch output shape as two distinct YAML blocks in §B.2 (rather than a flat union with optional fields) made the CUE design intent unambiguous. β's internal coherence walk confirmed the two shapes are fully disjoint.

**Nine-section structure from γ scaffold.** The γ scaffold prescribing the exact nine sections (§A–§I) in order eliminated structural ambiguity. β's section-by-section review confirmed all nine sections GREEN without requiring iteration.

**Friction notes as first-class content.** Including FN-1 through FN-8 verbatim in §H, rather than just citing them by reference, means the W1 implementer reads one document and gets the full friction picture without needing to trace back to prior γ scaffolds.

---

## What was challenging

**Naming precision at the CUE/renderer boundary (§D).** The 15-row boundary table required careful judgment on each row: is this a static shape concern (CUE) or a runtime/substrate concern (renderer)? The FN-6 attach-incompatibility refusal was the least clear-cut entry — it is a cross-field runtime condition the renderer enforces, but its trigger condition is not fully documented in the W0 source material. It was correctly assigned to the renderer column, but the exact trigger condition is deferred to the W1 implementer (see FN-6).

**Body authoring precision requirement (§E / AC2).** The body-as-prompt rule is stated clearly but its implementation involves a subtle constraint: the SKILL.md body must equal `prompt.md` verbatim, with verbose prose fields appended after — but the ordering of the appended prose is unspecified. β's Observation 2 correctly flags this as something the byte-identity oracle will enforce rather than the W0 design specifying. The ambiguity is acceptable at W0 but the W1 implementer should be aware.

**Distinguishing `triggers:` (skill-level) from `wake.input.triggers` (substrate-level) (Gap-4).** These two fields share a name stem but serve entirely different purposes. The standard SKILL.md `triggers:` field is the I5/discovery invocation trigger; `wake.input.triggers` carries the GitHub Actions `on:` event set for the renderer. Both must appear on a wake SKILL.md. The W0 design notes the distinction (§B.4), but the potential for W1 authoring confusion is real — this is Gap-4.

---

## Lessons for the W1 implementer

The following eight friction notes (FN-1 through FN-8) from the γ scaffold and `w0-design.md §H` are the primary guidance. Key points amplified here:

**FN-1 (observer role):** Before finalizing `#Wake` in `schemas/skill.cue`, audit `cn-install-wake` for any `observer` role handling. The W0 schema defines `admin | dispatch` only. If `observer` exists in the renderer and live wakes, add it to the enum before the I5 gate validates anything. Do not assume the enum is complete.

**FN-2 (`cue vet` with multi-doc frontmatter):** The I5 extractor must strip the SKILL.md body (everything after the second `---`) before running `cue vet`. Confirm this explicitly — do not rely on `cue vet` tolerating Markdown in the YAML stream.

**FN-3 (`agent_variable.default: null`):** The admin wake's `agent_variable.default` is `null` (operator-required; no default). YAML `null` must round-trip through the frontmatter extractor and through `cue vet` as CUE `null`, not as the empty string or absent. Test this explicitly with a fixture before claiming AC1.

**FN-4 (long `surfaces` arrays):** The `allowed_surfaces` and `disallowed_surfaces` arrays are 8–9 entries each. The CUE schema should use `[...string]` (no length constraint). Confirm the frontmatter extractor handles multi-line YAML arrays of this size without truncation.

**FN-5 (body verbatim constraint):** Do not paraphrase, reformat, or summarize `prompt.md` when authoring the SKILL.md body. Every trailing newline matters for the byte-identity oracle. If `prompt.md` has a trailing newline, the body must too.

**FN-6 (attach-incompatibility refusal):** Before W3 (renderer flip), confirm the exact trigger condition for this refusal from `cn-install-wake` source. AC6 requires it to fire from SKILL.md source. Author a synthetic mis-declared wake SKILL.md fixture to validate it.

**FN-7 (dual-source parity gate):** Design the W2 parity gate before implementing it. Options: a `--source=skill` flag on `cn install-wake`, a temporary parallel render path, or a CI-only comparison step. The gate must not break the existing `install-wake golden` CI check.

**FN-8 (deletion sequencing):** W4 must delete `wake-provider.json` and `prompt.md` in a single commit after byte-identical proof holds in CI — not across multiple commits. If the renderer uses a `SKILL.md then JSON fallback` pattern, the deletion is safe in a single commit. If the fallback path is removed separately, the commit sequence must be designed carefully to avoid a state where neither source is complete.

**Gap-4 reminder:** On each wake SKILL.md, author BOTH `triggers:` (standard skill-level field, describing skill invocation) AND `wake.input.triggers` (renderer-consumed field, encoding the GitHub Actions `on:` event set). They are different fields with different consumers.

---

## The four known gaps (non-blocking)

These gaps were identified in `self-coherence.md §6` and confirmed non-blocking by β. They are recorded here for the W1 implementer.

**Gap-1 — `observer` role:** The W0 `#Wake` schema defines `wake.role: "admin" | "dispatch"`. The renderer may accept an `observer` role internally. The W1 implementer must audit `cn-install-wake` to determine whether `observer` belongs in the `#Wake` CUE enum or remains renderer-internal. This must be resolved before the I5 gate validates the first wake SKILL.md.

**Gap-2 — `defer_path.*` string lengths:** The `defer_path.*` values in the current manifests are long prose strings (multiple sentences). The W0 schema types them as `string` with no length constraint. The W1 implementer must confirm these strings pass `cue vet` cleanly and that the renderer does not parse their content structurally (i.e., they are opaque strings to both CUE and the renderer's structural parser).

**Gap-3 — `governing_question` authorship:** The W0 design does not specify the `governing_question` text for either wake SKILL.md. This is standard skill-authoring work for W1, per the skill-authoring discipline. No design decision is needed; the W1 implementer authors these.

**Gap-4 — `triggers:` (standard) vs `wake.input.triggers` (renderer):** The standard SKILL.md `triggers:` field (skill-level invocation trigger, consumed by the I5/discovery pipeline) is distinct from `wake.input.triggers` (substrate event triggers encoded into the GitHub Actions `on:` block by the renderer). Both fields must appear on each wake SKILL.md. The W1 implementer must author both correctly for each wake.

All four gaps are scoped as W1 implementation concerns. None is a W0 design gap.

---

## Scope violation confirmation

No scope violations. The only diffs on `cycle/524` relative to main that belong to this cycle are under `.cdd/unreleased/524/`. The following files were explicitly NOT touched:

- `schemas/skill.cue`
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` (renderer)
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md`
- `.github/workflows/*.yml`
- `*.golden.yml`

No W1, W2, W3, or W4 implementation artifacts were created. β's scope compliance walk confirmed all six scope constraints PASS.

---

_Authored by α@cdd.cnos (R0 closeout), 2026-06-30 (UTC). R0 β-verdict: CONVERGE. Cycle/524 W0 design phase complete; W1→W4 build phases not yet scheduled._

---

# α-closeout — cnos#524 W1 R0: CUE #Wake schema + wake SKILL.md modules

---
cycle: 524
role: alpha
verdict: converge
date: 2026-06-30 (UTC)
authored_by: α@cdd.cnos (W1 R0 closeout)
parent_issue: cnos#524
round: W1-R0
beta_verdict: converge
---

## What was implemented

Three files changed on `cycle/524` relative to `main` for the W1 scope:

**AC1 — `schemas/skill.cue`:**
- `artifact_class` enum extended: `"wake"` added as the fifth value (was `"skill" | "runbook" | "reference" | "deprecated"`)
- `#WakeOutputAdmin` definition: required fields `channel_log_convention`, `writer_surface`, `class_taxonomy`, `cursor_advance`, `cursor_field`; open struct
- `#WakeOutputDispatch` definition: required fields `cycle_artifact_root`, `artifact_class_taxonomy`, `cell_runtime`; open struct
- `#Wake` definition: embeds `#Skill`; constrains `artifact_class: "wake"`, `scope: "global"`; full `wake:` block with `role: "admin" | "dispatch"` enum (OB-1); `agent_variable.default: string | null` (FN-3); `surfaces.allowed?/disallowed?` as `[...string]` (FN-4); output disjunction `#WakeOutputAdmin | #WakeOutputDispatch`

**AC2 — `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md`:**
- New file: 6 allowed + 9 disallowed surfaces; `agent_variable.default: null` (literal null, operator-required); `activation_log_writer: true`; body = verbatim `prompt.md` content + verbose prose from JSON

**AC2 — `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`:**
- New file: `activation_state: live`; `protocol: cds`; selector include/exclude matching `wake-provider.json` exactly; `agent_variable.default: sigma`; `activation_log_writer: false`; body = verbatim `prompt.md` content (242 lines) + verbose prose from JSON

No renderer, golden, workflow, or wake-provider.json files were touched.

## Scope compliance

Zero diff on all constrained paths: `*.golden.yml`, `wake-provider.json`, `prompt.md`, `.github/workflows/`, renderer source.

## Friction notes encountered

FN-1 (observer role): OB-1 check from W0 β. Audited before authoring `#Wake` — no `observer` role found in any live `wake-provider.json`; enum correctly constrained to `"admin" | "dispatch"`.

FN-3 (null default): YAML `literal null` written for `agent_variable.default` on admin SKILL.md. Confirmed not empty string.

FN-5 (body verbatim): Both bodies verified line-by-line against respective `prompt.md` files before committing.

_Authored by α@cdd.cnos (W1 R0 closeout), 2026-06-30 (UTC). W1 R0 β-verdict: CONVERGE. AC1 + AC2 delivered._
