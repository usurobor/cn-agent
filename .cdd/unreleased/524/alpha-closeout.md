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

---

# α-closeout — cnos#524 W2 R1: renderer dual-source parity

---
cycle: 524
role: alpha
verdict: converge
date: 2026-06-30 (UTC)
authored_by: α@cdd.cnos (W2 R1 closeout)
parent_issue: cnos#524
round: W2-R1
beta_verdict: pending
---

## What was implemented

Four files changed on `cycle/524` relative to `main` for the W2 scope:

**`src/packages/cnos.core/commands/install-wake/cn-install-wake` (renderer extension):**
- `skill_to_json_manifest()` helper: converts SKILL.md YAML frontmatter `wake:` block → synthesized
  `cn.wake-provider.v1` JSON with the same jq paths as `wake-provider.json`. Preserves
  `activation_log_writer` has()-vs-absent distinction. Uses python3 + pyyaml.
- `skill_body()` helper: Python-based extractor for the verbatim-prompt portion of SKILL.md body.
  Strips the leading blank line (standard blank after closing `---` delimiter) and stops at
  `\n\n---\n\n## .+(from wake-provider.json|body reference)` — the W1 boundary marker for
  reference-data sections appended beyond the verbatim `prompt.md` content.
- `--source skill` flag: routes manifest and prompt through `skill_to_json_manifest()` +
  `skill_body()` instead of `wake-provider.json` + `prompt.md`.
- `--parity-check` flag: implies `--source skill`; renders to temp file; compares with committed
  golden (stripping `^#` lines from both sides per FN-W2-1); exits 0 on match, 5 on diff.
- Header stability fix: `display_manifest_path` / `display_prompt_path` always use the canonical
  JSON-source paths (`json_manifest_path` / `json_prompt_path`) regardless of `--source` flag,
  so the golden's comment header is stable across invocations.
- Exit code 5 documented in header comments.

**`.github/workflows/install-wake-golden.yml` (parity guard step):**
- One new step added BEFORE the AC8 authority audit step:
  ```
  W2 parity check — render(SKILL.md) == render(JSON+prompt)
  cn-install-wake agent-admin --parity-check
  cn-install-wake cds-dispatch --parity-check
  ```

**`.cdd/unreleased/524/gamma-scaffold.md`:**
- W2-scoped γ scaffold superseding the W1 scope document.

**`.cdd/unreleased/524/self-coherence.md §R1`:**
- W2 implementation self-coherence section appended.

## Parity verification results

```
cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for agent-admin  EXIT: 0
cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for cds-dispatch EXIT: 0
cn-install-wake: agent-admin → …/cnos-agent-admin.golden.yml (unchanged)  (default render)
cn-install-wake: cds-dispatch → …/cnos-cds-dispatch.golden.yml (unchanged) (default render)
cn-install-wake: agent-admin → …/cnos-agent-admin.golden.yml (unchanged)  (--source skill)
cn-install-wake: cds-dispatch → …/cnos-cds-dispatch.golden.yml (unchanged) (--source skill)
```

## Scope compliance

| Constraint | Status |
|---|---|
| `schemas/skill.cue` | Not touched |
| `wake-provider.json` (both wakes) | Not touched |
| `prompt.md` (both wakes) | Not touched |
| `cnos-agent-admin.golden.yml` | Not touched (`git diff` confirmed) |
| `cnos-cds-dispatch.golden.yml` | Not touched (`git diff` confirmed) |
| `.github/workflows/cnos-agent-admin.yml` | Not touched |
| `.github/workflows/cnos-cds-dispatch.yml` | Not touched |
| W3/W4 implementation artifacts | Not created |
| `#524` issue state | Open (no close/fix/resolve linkage in commits) |

## Friction notes encountered

**FN-W2-1 (header-line exclusion):** Parity comparison strips `^#` prefix lines before `cmp`.
Verified working: parity reports OK for both wakes despite differing `# manifest:` / `# prompt:`
header comments between JSON and skill renders.

**FN-W2-2 (body boundary detection):** The Python `skill_body()` extractor detects the boundary
using the `\n\n---\n\n## .+(from wake-provider.json|body reference)` pattern. Both wakes'
reference-data sections are correctly excluded from the parity comparison.

**Awk-to-Python migration:** The initial `skill_body()` awk implementation had a subtle bug:
`/^---/{c++; next}` consumed ALL bare `---` lines (both frontmatter delimiters and body
horizontal-rule separators), silently dropping the latter. The fix required switching to Python
for reliable body extraction with look-ahead boundary detection.

_Authored by α@cdd.cnos (W2 R1 closeout), 2026-06-30 (UTC). W2 R1 β-verdict: pending β review._

---

# α-closeout — cnos#524 W3 R2: renderer source-flip

---
cycle: 524
role: alpha
verdict: converge
date: 2026-06-30 (UTC)
authored_by: α@cdd.cnos (W3 R2 closeout)
parent_issue: cnos#524
round: W3-R2
beta_verdict: converge
---

## What was implemented

Three files changed on `cycle/524` for the W3 scope (on top of the W2 baseline):

**`src/packages/cnos.core/commands/install-wake/cn-install-wake` (6 targeted edits):**
- Line 298: `source_type="json"` → `source_type="skill"` — default source is now SKILL.md
- Line 415-416: `--parity-check` now implies `--source json` (was `--source skill`) — the
  parity gate now proves render(JSON+prompt) == render(SKILL.md)
- Lines 55-63: `--source` doc comment updated — `skill` listed as default; `json` as override
- Lines 63-75: `--parity-check` doc comment updated — now describes JSON+prompt-sourced render
  against the SKILL.md-produced golden; "Implies --source json"; "W3 CI parity gate"
- Lines 102-105: exit-5 doc comment updated — "render(JSON+prompt) differs from render(SKILL.md)"
- Lines 1125-1128: parity success/failure messages updated — "render(JSON+prompt) ==
  render(SKILL.md)" in both success and failure forms

**`.github/workflows/install-wake-golden.yml` (one step update):**
- Step renamed from "W2 parity check — render(SKILL.md) == render(JSON+prompt)" to
  "W3 parity check — render(JSON+prompt) == render(SKILL.md)"
- Step comment updated to reflect W3 semantics (goldens now from SKILL.md; `--parity-check`
  implies `--source json`). Run block unchanged.

**`.cdd/unreleased/524/` CDD artifacts:**
- `gamma-scaffold.md`: W3 γ scaffold (replaced W2 scaffold)
- `self-coherence.md §R2`: W3 implementation self-coherence section appended
- `beta-review.md §R2`: W3 β review section appended
- `alpha-closeout.md §W3`: this section

## Scope compliance

| Constraint | Status |
|---|---|
| `schemas/skill.cue` | Not touched |
| `wake-provider.json` (both wakes) | Not touched |
| `prompt.md` (both wakes) | Not touched |
| `SKILL.md` (both wakes) | Not touched |
| `cnos-agent-admin.golden.yml` | Not touched |
| `cnos-cds-dispatch.golden.yml` | Not touched |
| `.github/workflows/cnos-agent-admin.yml` | Not touched |
| `.github/workflows/cnos-cds-dispatch.yml` | Not touched |
| W4 implementation artifacts | Not created |
| Issue #524 remains open | Commit uses `Refs #524` / `Part of #524` only |

`git diff --stat HEAD` confirms exactly 3 files changed: `cn-install-wake`,
`install-wake-golden.yml`, `.cdd/unreleased/524/gamma-scaffold.md` (plus this and other
`.cdd/` artifacts added during this session).

## AC oracle status

| AC | Status |
|---|---|
| AC3 (default reads SKILL.md) | SATISFIED — `source_type="skill"` at init; code inspection confirms; CI re-render expected to pass |
| AC4 (byte-identical goldens) | SATISFIED — no golden in diff; W2 parity transitivity holds |
| AC6 (refusals preserved) | SATISFIED — no gate code touched; CI smokes unaffected |
| AC7/AC8 (CI green; no role strings) | EXPECTED PASS — no role-decision strings added |

_Authored by α@cdd.cnos (W3 R2 closeout), 2026-06-30 (UTC). W3 R2 β-verdict: CONVERGE._

---

# α-closeout — cnos#524 W4 R3: delete wake-provider.json + prompt.md; renderer is SKILL.md-only (final phase)

---
cycle: 524
role: alpha
verdict: converge
date: 2026-06-30 (UTC)
authored_by: α@cdd.cnos (W4 R3 closeout)
parent_issue: cnos#524
round: W4-R3
beta_verdict: converge
run_class: repair_pass
---

## What was implemented

This is the **final phase** of the W0→W4 wake-as-skill migration. `run_class: repair_pass` per
`REPAIR-PLAN.md`: a prior W4 dispatch run (substrate run `28464981342`) was invalidated as an
**empty run** — it reached `status:review` with no PR, no commits, no closeout, and no deliverable.
That process gap was remediated by a separate, already-merged cell (PR #531, the
closeout-integrity guard). This cycle performed the actual W4 implementation, against current
`main` (`db547ebe5b408e4c74092ad3ed56509e605894ef`), under the operator's "W4 (final phase, clean
re-dispatch)" directive and `gamma-scaffold.md`.

**Deleted (8 files — the legacy JSON+prompt contract system, in full):**
- `src/packages/cnos.core/orchestrators/agent-admin/{wake-provider.json,prompt.md}`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/{wake-provider.json,prompt.md}`
- `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/{wake-provider.json,prompt.md}`
- `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/{wake-provider.json,prompt.md}`

Each directory already carried a `SKILL.md` twin (added across W1/W3); only the legacy pair was
removed. Both wake SKILL.md files themselves carry **zero diff** in this cycle (frontmatter and
body both untouched — see "Known gap," below).

**`cn-install-wake` made SKILL.md-only:** `--source json|skill` and `--parity-check` removed as
live flags — both (and `--source=*`) now hit a dedicated `case` arm that `die`s with a message
naming cnos#524 W4, rather than silently disappearing or falling through to a now-broken default.
The `source_type`/`parity_check` shell variables are gone entirely; the SKILL.md→JSON synthesis
(`skill_to_json_manifest`) and body extraction (`skill_body`) are now unconditional — there is
exactly one code path. Manifest resolution (both the default per-package lookup and the
cross-package sibling-search fallback) now resolves `SKILL.md` instead of `wake-provider.json`.
Header attribution changed from the two-line `# manifest: .../wake-provider.json` /
`# prompt: .../prompt.md` block to a single `# source: orchestrators/<name>/SKILL.md` line — the
**only** byte-level change to any golden or live-workflow output. Exit code 5 (`--parity-check`
failure) is documented as retired, not reused.

**`.github/workflows/install-wake-golden.yml`:** the W3 parity-check step is removed (nothing left
to compare). The AC2 negative-case smoke converts from a malformed-JSON fixture to a
malformed-SKILL.md fixture (`wake:` block omitting `role`); same oracle (malformed declaration →
exit 2 + informative stderr), different message text because the code path is different (see
judgment call below). The AC5 and AC4/cycle-496 refusal-smoke steps' `--manifest` arguments are
repointed at the fixtures' `SKILL.md` files; assertions (exit 3 / exit 4, same stderr substrings)
are unchanged.

**Goldens + live workflows re-rendered, header-only diff** (verified twice — see AC verification
below): both `cnos-agent-admin.golden.yml` / `cnos-cds-dispatch.golden.yml` and both
`.github/workflows/cnos-*.yml`.

## AC verification summary (W4-relevant ACs)

| AC | Result |
|---|---|
| AC4 (revised — header-only diff) | Verified twice: renderer-change-alone diff against pre-W4 goldens, then post-deletion diff against the same baseline. Both confirmed header-only — `git diff --cached` inspected per-file, not just `--stat`. sha256(live)==sha256(golden) holds (cds-dispatch). Idempotence holds for both wakes. |
| AC5 (JSON+prompt deleted) | `find` over both orchestrator dirs + both test-fixture dirs for `wake-provider.json`/`prompt.md` → empty. Renderer reads only SKILL.md; cross-package sibling-fallback verified explicitly with `CN_PACKAGE_ROOT` unset. |
| AC6 (refusals preserved) | AC5 declaration-only smoke → exit 3, correct stderr. AC4/cycle-496 mis-declaration smoke → exit 4, correct stderr. Both fire from the SKILL.md-only path (no other path exists). |
| AC7 (all gates green) | Go build/vet/test clean; I1 (`./cn build --check`) clean; I2 clean; I6-adjacent (`cn cdd verify`) 106 passed/0 failed; `install-wake-golden` steps replayed individually, all green; both CI guard scripts (see judgment call below) green. I5 not runnable locally (no `cue` binary in this environment) — flagged for CI; neither `schemas/skill.cue` nor either wake SKILL.md's frontmatter changed this cycle, so I5 has no mechanical reason to regress. |

Full verification detail, including exact commands and outputs, is in
`self-coherence.md §R3 §3–§4`; β independently re-ran every check rather than trusting this report
(see `beta-review.md §R3`).

## Judgment calls

**Out-of-scaffold CI-script fix (`check-dispatch-closeout-integrity.sh` /
`check-dispatch-repair-preflight.sh`).** Both scripts hardcoded a required-file check against
`cds-dispatch/prompt.md`, left over from the W2/W3 dual-source transition window. Deleting
`prompt.md` per the W4 mandate would have made both scripts fail with `required file missing` —
turning two required-green CI gates red as a direct, mechanical consequence of the deletion the
scaffold itself required. I verified the scripts' required phrase-checks against `$SKILL` /
`$GOLDEN` / `$LIVE` were already satisfied independent of `$PROMPT`, then swapped each script's
loop target from the now-deleted `$PROMPT` variable to the already-present `$SKILL` variable — no
phrase, label, or detection logic changed, only the file-existence target. This was not on the
scaffold's original MUST-change list, but it is the minimum necessary fix to keep two
scaffold-required-green gates green after a deletion the scaffold itself mandates, not scope
creep. β independently verified the causal claim (would these scripts actually have broken?) was
true, not merely asserted, and confirmed the fix is contract-preserving.

Other judgment calls (`--manifest` accepting a directory or file; the AC2 fixture shape choice;
making `--source`/`--parity-check` hard errors rather than silent no-ops; retiring exit code 5
with a documentation note rather than dropping the row) are recorded in full, with rationale, in
`self-coherence.md §R3 §2`.

## Known gap (explicitly out of scope, named for follow-up)

Both wake SKILL.md files — `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` and
`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — have **zero diff** in this cycle.
Their body prose still references the now-deleted `wake-provider.json` in a few places: most
visibly `cds-dispatch/SKILL.md`'s `[wake-provider.json](wake-provider.json)` markdown link (now a
dangling 404 target), plus both files' `## Responsibilities (from wake-provider.json; body
reference)` / `## Cross-references (from wake-provider.json)` section-header parentheticals. This
was a deliberate scope decision, not an oversight: the SKILL.md body **is** the rendered prompt
(W0 design §E, body-as-prompt rule), so editing it would change non-header bytes of the rendered
golden/live-workflow output — directly violating the operator's "Stop if: any non-header workflow
bytes change" condition. The "no active references to wake-provider.json/prompt.md remain"
required-proof bullet is satisfied per the scaffold's own resolution (§1.1): by (a) the files
being deleted (AC5) and (b) zero renderer resolution/parsing-logic references remaining — not by
scrubbing frozen prompt-body prose. This is a legitimate, separately-scoped follow-up (a prompt
content change, not a renderer/substrate change) for a future cell. A tracked follow-up issue is
filed at closeout (see `gamma-closeout.md`'s process-gap audit for the issue number).

---

_Authored by α@cdd.cnos (W4 R3 closeout), 2026-06-30 (UTC). W4 R3 β-verdict: CONVERGE. This closes
the W0→W4 wake-as-skill migration: one contract system (SKILL.md + CUE), the JSON+prompt system
fully retired._
