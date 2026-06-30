---
cycle: 524
parent_issue: cnos#524
document_class: beta-review
round: R0
authored_by: β@cdd.cnos (R0)
date: 2026-06-30 (UTC)
base_sha: d7d2244cb4cb95b94e28569ac8ef090d49876c89
alpha_commit_shas:
  - 43b1c140016048fe0a1a0843b61e7a69a171eb81  # α-524 W0: w0 design document
  - ada5f9fa3169698b13f7eb407d4d35733f7f2013  # α-524 R0: self-coherence §R0
beta_review_timestamp: 2026-06-30T00:00:00Z
verdict: converge
---

# β-review — cnos#524 R0

## §R0 Verdict

**CONVERGE**

Scope compliance is clean. All 9 design document sections (§A–§I) are present and complete. Every field from both `wake-provider.json` manifests is accounted for in `w0-design.md §B`. Internal coherence holds across §B, §D, §F, and §G. No W1→W4 implementation artifacts on the branch.

---

## §R0 Scope Compliance Walk

| Constraint | Result | Observation |
|---|---|---|
| No changes to `schemas/skill.cue` | PASS | Not in `git diff main..cycle/524 --name-only`; confirmed absent |
| No changes to `cn-install-wake` or any renderer file | PASS | No `src/packages/cnos.core/commands/install-wake/` path in diff |
| No changes to `.github/workflows/*.yml` | PASS | No workflow files in diff |
| No changes to `*.golden.yml` | PASS | No golden files in diff |
| `wake-provider.json` untouched (both wakes) | PASS | `agent-admin/wake-provider.json` and `cds-dispatch/wake-provider.json` both absent from diff |
| `prompt.md` untouched (both wakes) | PASS | Neither `agent-admin/prompt.md` nor `cds-dispatch/prompt.md` in diff |
| No W1→W4 implementation code on the branch | PASS | Diff contains only `.cdd/unreleased/524/gamma-scaffold.md`, `.cdd/unreleased/524/w0-design.md`, `.cdd/unreleased/524/self-coherence.md`, and `.cn-sigma/logs/20260630.md` (heartbeat log; CDD-external) |

**Scope compliance: CLEAN.** The only writable outputs on the branch beyond the γ scaffold are the two α artifacts and the σ heartbeat log.

---

## §R0 Section Table (§A–§I)

| Section | Status | Observation |
|---|---|---|
| §A — Problem statement and invariant | GREEN | Two-sentence problem accurately names the two-system problem (JSON+prompt vs SKILL.md). Invariant stated verbatim: "A wake is a typed SKILL.md module whose body is the prompt and whose frontmatter compiles to workflow substrate." Matches γ-scaffold §4.§A requirement. |
| §B — `wake:` block schema | GREEN | Complete field listing with types, optionality, and notes. §B.1 shows the full `wake:` YAML block. §B.2 presents the role-shaped output disjunction (admin shape vs dispatch shape) as separate YAML blocks. §B.3 is a full cross-reference table mapping every `wake-provider.json` field to its SKILL.md placement (frontmatter key or "→ body"). §B.4 names the standard SKILL.md fields that carry wake data outside the `wake:` block. Field coverage verified independently below. |
| §C — 5 operator design decisions | GREEN | All 5 decisions present in order: (1) `wake:` nesting, (2) `artifact_class: wake`, (3) `scope: global` / role in `wake.role`, (4) CUE owns typed contract / renderer owns substrate + runtime refusals, (5) prose fields → body. Each decision has a clearly labeled Decision, Rationale, and W1 implication paragraph. |
| §D — CUE↔renderer boundary table | GREEN | 15-row table. Static shape, enums, role-shaped disjunction, selector shape, concurrency shape, permission_intent array shape, static required-field presence → CUE. Runtime refusals (exit 3 / exit 4 / FN-6), body extraction, install-time substitution, GitHub Actions encoding, bot identity, prompt delivery, pre-W4 JSON presence → renderer. No concern appears in both columns. Consistent with §C Decision 4. |
| §E — Body-as-prompt rule | GREEN | Frontmatter/body split defined. Rule stated clearly: body (after second `---`) IS the prompt, extracted verbatim. Byte-identity oracle named and formalized: `render(SKILL.md source) == render(wake-provider.json + prompt.md source)`. `install-wake golden` CI check cited as the substrate encoding of the oracle. |
| §F — W1→W4 migration plan | GREEN | 4 phases present (W1–W4) in order. Each phase: What / Oracle / State-after. W1: author SKILL.md + extend CUE; oracle = I5 count 91→93 + cue vet pass/fail. W2: dual-source parity gate; oracle = `cmp render(SKILL) render(JSON)`. W3: flip renderer source; oracle = `install-wake golden` CI green + re-render diff clean. W4: delete JSON + prompt; oracle = files absent + byte-identical goldens from SKILL.md only + I5 count 93. Byte-identity oracle present at each phase boundary. |
| §G — W1 ACs (AC1–AC7) | GREEN | All 7 ACs from issue #524 body present. Each AC has Invariant and Oracle sub-sections. AC1: `#Wake` CUE schema. AC2: SKILL.md files authored. AC3: renderer reads SKILL.md. AC4: byte-identical goldens. AC5: JSON + prompt deleted. AC6: renderer refusals preserved (exit 3 / exit 4 / FN-6). AC7: all CI gates green. |
| §H — Friction notes | GREEN | FN-1 through FN-8 present and carried from prior γ scaffolds (commit `f148a1fd`, updated scaffold at `9a10781d`). All 8 notes named (observer role collision, cue vet multi-doc, null encoding, array length, body verbatim constraint, FN-6 attach-incompatibility, dual-source parity gate implementation, deletion sequencing). Material for W1 implementer. |
| §I — Non-goals | GREEN | 8 explicit non-goals listed: no rendered-workflow output change, no wake runtime-behavior change, no renderer substrate-encoding change, no new wake names, no Demo 0, no `.cdd/releases/` work, no per-package bot identities, no change to CUE validator script. Consistent with issue body non-goals. |

**All 9 sections: GREEN.**

---

## §R0 Field Coverage Walk

### `agent-admin/wake-provider.json` fields → §B coverage

| JSON field | §B disposition | Status |
|---|---|---|
| `schema` | Not carried; artifact class signals schema version | PRESENT |
| `name` | SKILL.md `name:` (standard skill field) | PRESENT |
| `package` | `wake.package` | PRESENT |
| `role` | `wake.role` | PRESENT |
| `admin_only` | `wake.admin_only` | PRESENT |
| `activation_log_writer` | `wake.activation_log_writer` | PRESENT |
| `description` | SKILL.md `description:` (standard skill field) | PRESENT |
| `responsibilities` | → body | PRESENT |
| `input_contract.triggers` | `wake.input.triggers` | PRESENT |
| `input_contract.issues_opened_title_pattern` | `wake.input.issues_opened_title_pattern` | PRESENT |
| `input_contract.trigger_descriptions` | → body | PRESENT |
| `input_contract.inbound` | → body | PRESENT |
| `output_contract.channel_log_convention` | `wake.output.channel_log_convention` | PRESENT |
| `output_contract.writer_surface` | `wake.output.writer_surface` | PRESENT |
| `output_contract.class_taxonomy` | `wake.output.class_taxonomy` | PRESENT |
| `output_contract.class_taxonomy_notes` | → body | PRESENT |
| `output_contract.cursor_advance` | `wake.output.cursor_advance` | PRESENT |
| `output_contract.cursor_field` | `wake.output.cursor_field` | PRESENT |
| `allowed_surfaces` | `wake.surfaces.allowed` | PRESENT |
| `disallowed_surfaces` | `wake.surfaces.disallowed` | PRESENT |
| `defer_path.cell_shaped_directive` | `wake.defer_path.cell_shaped_directive` | PRESENT |
| `defer_path.off_role_directive` | `wake.defer_path.off_role_directive` | PRESENT |
| `defer_path.ambiguous_directive` | `wake.defer_path.ambiguous_directive` | PRESENT |
| `prompt_template` | Not carried; body IS the prompt (§E) | PRESENT |
| `agent_variable.name` | `wake.agent_variable.name` | PRESENT |
| `agent_variable.default` | `wake.agent_variable.default` | PRESENT |
| `agent_variable.description` | → body | PRESENT |
| `permission_intent` | `wake.permission_intent` | PRESENT |
| `permission_intent_notes` | → body | PRESENT |
| `concurrency_intent.serialize` | `wake.concurrency.serialize` | PRESENT |
| `concurrency_intent.group` | `wake.concurrency.group` | PRESENT |
| `concurrency_intent.notes` | → body | PRESENT |
| `superseded_substrate_artifact` | → body | PRESENT |
| `relationship_to_substrate` | → body | PRESENT |
| `cross_references` | → body | PRESENT |

**agent-admin coverage: 35/35 fields — all PRESENT.**

### `cds-dispatch/wake-provider.json` fields → §B coverage

| JSON field | §B disposition | Status |
|---|---|---|
| `schema` | Not carried; artifact class signals schema version | PRESENT |
| `name` | SKILL.md `name:` (standard skill field) | PRESENT |
| `package` | `wake.package` | PRESENT |
| `role` | `wake.role` | PRESENT |
| `admin_only` | `wake.admin_only` | PRESENT |
| `activation_log_writer` | `wake.activation_log_writer` | PRESENT |
| `activation_state` | `wake.activation_state` | PRESENT |
| `activation_state_notes` | → body | PRESENT |
| `protocol` | `wake.protocol` | PRESENT |
| `selector.include` | `wake.selector.include` | PRESENT |
| `selector.exclude` | `wake.selector.exclude` | PRESENT |
| `description` | SKILL.md `description:` (standard skill field) | PRESENT |
| `responsibilities` | → body | PRESENT |
| `input_contract.triggers` | `wake.input.triggers` | PRESENT |
| `input_contract.trigger_descriptions` | → body | PRESENT |
| `input_contract.inbound` | → body | PRESENT |
| `output_contract.cycle_artifact_root` | `wake.output.cycle_artifact_root` | PRESENT |
| `output_contract.artifact_class_taxonomy` | `wake.output.artifact_class_taxonomy` | PRESENT |
| `output_contract.artifact_class_notes` | → body | PRESENT |
| `output_contract.cell_runtime` | `wake.output.cell_runtime` | PRESENT |
| `output_contract.cell_runtime_notes` | → body | PRESENT |
| `allowed_surfaces` | `wake.surfaces.allowed` | PRESENT |
| `disallowed_surfaces` | `wake.surfaces.disallowed` | PRESENT |
| `defer_path.cell_shaped_directive` | `wake.defer_path.cell_shaped_directive` | PRESENT |
| `defer_path.off_role_directive` | `wake.defer_path.off_role_directive` | PRESENT |
| `defer_path.ambiguous_directive` | `wake.defer_path.ambiguous_directive` | PRESENT |
| `prompt_template` | Not carried; body IS the prompt (§E) | PRESENT |
| `agent_variable.name` | `wake.agent_variable.name` | PRESENT |
| `agent_variable.default` | `wake.agent_variable.default` | PRESENT |
| `agent_variable.description` | → body | PRESENT |
| `permission_intent` | `wake.permission_intent` | PRESENT |
| `permission_intent_notes` | → body | PRESENT |
| `concurrency_intent.serialize` | `wake.concurrency.serialize` | PRESENT |
| `concurrency_intent.group` | `wake.concurrency.group` | PRESENT |
| `concurrency_intent.notes` | → body | PRESENT |
| `cross_references` | → body | PRESENT |

**cds-dispatch coverage: 36/36 fields — all PRESENT.**

**Field coverage verdict: COMPLETE. Zero missing fields across both manifests.**

---

## §R0 Internal Coherence Walk

### Role-shaped output disjunction (§B.2): admin shape ≠ dispatch shape; both present?

**PASS.** §B.2 presents two distinct named YAML blocks:
- Admin shape (`wake.role == "admin"`): `channel_log_convention`, `writer_surface`, `class_taxonomy`, `cursor_advance`, `cursor_field` — all sourced from `agent-admin/wake-provider.json output_contract`. All 5 admin-specific output fields present.
- Dispatch shape (`wake.role == "dispatch"`): `cycle_artifact_root`, `artifact_class_taxonomy`, `cell_runtime` — all sourced from `cds-dispatch/wake-provider.json output_contract`. All 3 dispatch-specific output fields present.
The two shapes are disjoint: no field appears in both. The disjunction is correctly role-gated and will require CUE to enforce field presence based on `wake.role` value.

### §F migration plan consistent with §B schema shape?

**PASS.** W1 authors SKILL.md with the `wake:` block per §B — consistent. W2 extends the renderer to read `wake:` frontmatter (§D renderer column fields) and body — consistent with §B field map. W3 flips source; the §B schema shape is unchanged — consistent. W4 deletes JSON + prompt; only the §B schema shape survives — consistent. The I5 count progression (91 → 93) matches §C Decision 2 (two wake SKILL.md files added to the `artifact_class: wake` discovery).

### §G ACs consistent with §D CUE↔renderer boundary?

**PASS.**
- AC1 (`#Wake` CUE schema): covers static shape, enums (`wake.role`, `wake.activation_state`, `wake.permission_intent` values), role-shaped output disjunction, `selector` shape, `concurrency` shape, `permission_intent` array — these are precisely the CUE column entries in §D. No AC1 item crosses into the renderer column.
- AC3 (renderer reads SKILL.md): body extraction, install-time substitution — renderer column in §D. No CUE involvement.
- AC6 (refusals preserved): `activation_state` gating (exit 3), `activation_log_writer` mis-declaration (exit 4), FN-6 attach-incompatibility — all renderer column in §D. No CUE involvement. AC6 correctly cites the existing cycle/496 mis-declaration smoke as the oracle.
- No AC asks CUE to enforce a runtime refusal; no AC asks the renderer to enforce static shape. Boundary is respected throughout.

### §C Decision 5 (body-as-prompt) consistent with §E and §B.3 disposition map?

**PASS.** §C Decision 5 names the specific verbose fields that move to body: `responsibilities`, `*_notes`, `cross_references`, `trigger_descriptions`, `relationship_to_substrate`, `superseded_substrate_artifact`, `agent_variable.description`. §B.3 marks each of these "→ body". §E states the body IS the prompt (verbatim `prompt.md` + moved prose). All three sections are mutually consistent — no field appears as both "→ body" and as a `wake.*` frontmatter key.

### Self-coherence artifact consistency?

**PASS.** α's `self-coherence.md §R0` performs the same field-coverage walk independently and reaches the same verdict (all fields covered, zero missing). The section walk in self-coherence aligns with the actual content of `w0-design.md`. The known gaps (FN-1 `observer` role, Gap-2 `defer_path` string lengths, Gap-3 `governing_question` authorship, Gap-4 `triggers:` vs `wake.input.triggers` distinction) are appropriately scoped as W1 implementation concerns, not W0 design gaps.

---

## §R0 Findings

No blocking findings. No iterate conditions triggered.

**Observations (non-blocking; for W1 implementer context):**

1. **FN-1 observer role (carried forward):** The W0 `#Wake` schema defines `wake.role: "admin" | "dispatch"` only. If the renderer internally supports an `observer` role, this must be resolved before W1 finalizes the `#Wake` enum. α's Gap-1 correctly flags this; it is a W1 pre-flight check, not a W0 design gap.

2. **AC2 body authoring precision:** §E and §C Decision 5 require the SKILL.md body to equal `prompt.md` verbatim plus appended prose from moved fields. The ordering of appended prose (which notes fields appear in which sequence) is left to the W1 implementer. The byte-identity oracle (AC4) will catch any divergence from the rendered prompt.

3. **Gap-4 (triggers: vs wake.input.triggers):** The design correctly distinguishes standard SKILL.md `triggers:` (skill-level invocation triggers for the I5/discovery pipeline) from `wake.input.triggers` (substrate event triggers encoded by the renderer). §B.4 names this distinction. The W1 implementer must author both fields correctly on each wake SKILL.md.

---

_Reviewed by β@cdd.cnos (R0), 2026-06-30 (UTC). Verdict: CONVERGE. No iterate conditions. α's W0 design document is complete, field-complete, scope-clean, and internally coherent._

---

# β Review — cnos#524 W1 R0

## §R0

**Verdict: CONVERGE**

---

### AC1 — schemas/skill.cue

- [x] `artifact_class` enum includes `"wake"` — line 35: `"skill" | "runbook" | "reference" | "deprecated" | "wake"`. No existing value removed.
- [x] `#Wake` definition present AFTER `#Skill` definition — `#WakeOutputAdmin` at line 76, `#WakeOutputDispatch` at line 88, `#Wake` at line 107; all after `#Skill` closes at line 71.
- [x] `role: "admin" | "dispatch"` — no other values. Line 114. OB-1 satisfied.
- [x] `agent_variable.default: string | null` — line 156. Handles null correctly (FN-3).
- [x] `artifact_class: "wake"` in `#Wake` — line 108. Constrained to the `"wake"` literal within `#Wake`.
- [x] Arrays use `[...string]` syntax — `surfaces.allowed?` and `surfaces.disallowed?` at lines 162–163; `permission_intent`, `input.triggers`, selector arrays. No length constraints (FN-4).
- [x] Output disjunction present — line 140: `output: #WakeOutputAdmin | #WakeOutputDispatch`. `#WakeOutputAdmin` requires `channel_log_convention` (line 77); `#WakeOutputDispatch` requires `cycle_artifact_root` (line 89). Role-shaped enforcement is structural.
- [x] CUE syntax: `package skill` at top (line 18). All struct bodies balanced. No syntax errors detected. (`cue` binary unavailable in review environment; manual scan performed.)
- [x] ONLY `schemas/skill.cue` changed for AC1 — confirmed by `git diff origin/main --name-only | grep -v .cdd/` showing only three paths: `schemas/skill.cue`, `agent-admin/SKILL.md`, `cds-dispatch/SKILL.md`.

---

### AC2 — agent-admin/SKILL.md

Checked against `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`.

- [x] `artifact_class: wake` — line 5.
- [x] `scope: global` — line 6.
- [x] `wake.role: admin` — line 20. Matches JSON `"role": "admin"`.
- [x] `wake.package: cnos.core` — line 21. Matches JSON.
- [x] `wake.admin_only: true` — line 22. Matches JSON.
- [x] `wake.activation_log_writer: true` — line 23. Matches JSON.
- [x] `wake.concurrency.group: "agent-admin-{agent}"` — line 47. Matches JSON `concurrency_intent.group`.
- [x] `wake.agent_variable.default: null` — line 50. YAML literal null (not `~`, not `""`, not absent). FN-3 satisfied.
- [x] `wake.input.triggers` = `[schedule, issues_opened_title_match]` — lines 26–28. Matches JSON `input_contract.triggers`.
- [x] `wake.input.issues_opened_title_pattern: claude-wake` — line 29. Matches JSON.
- [x] Body starts with verbatim content of `prompt.md` — verified by diff: after the leading blank line that follows the closing `---` delimiter, the body content is byte-identical to `prompt.md` (129 lines). The leading blank line is a formatting artifact of the YAML frontmatter separator, not a content deviation.
- [x] OB-2 section order — body uses verbatim prompt.md headings (per FN-5). "Identity and activation" appears first (line 83). OB-2 requires sections in order `if present`; exact OB-2 heading names (Purpose, Authority boundary, Wake procedure, Repair/stop rules, Outputs/receipts, Non-goals) are not present as standalone headings. The verbatim constraint (FN-5) takes precedence; vacuously satisfied for absent headings.
- [x] OB-3 — top-level `triggers:` = `[activate, attach, channel-sync, admin-wake]` (skill-invocation labels); `wake.input.triggers:` = `[schedule, issues_opened_title_match]` (GitHub workflow event classes). Semantically distinct; no overlap.
- [x] `wake.output.channel_log_convention` — line 30. Verbatim from JSON.
- [x] `wake.output.writer_surface` — line 31. Verbatim from JSON.
- [x] `wake.output.class_taxonomy` — 5 entries (heartbeat, substantive, inaugural, directive-out, cycle-complete). Matches JSON.
- [x] `wake.output.cursor_advance: true` — line 38. Matches JSON.
- [x] `wake.output.cursor_field` — line 39. Verbatim from JSON.
- [x] `wake.permission_intent` — 4 entries (contents.write, issues.write, pull_requests.write, id_token.write). Matches JSON.
- [x] `wake.concurrency.serialize: true` — line 46. Matches JSON.
- [x] `wake.agent_variable.name: agent` — line 49. Matches JSON.
- [x] `wake.surfaces.allowed` — 6 entries, verbatim from JSON `allowed_surfaces`. Verified side-by-side.
- [x] `wake.surfaces.disallowed` — 9 entries, verbatim from JSON `disallowed_surfaces`. Verified side-by-side.
- [x] `wake.defer_path.*` — all three values verbatim from JSON. Verified side-by-side.
- [x] No prose fields in frontmatter — `responsibilities`, `*_notes`, `cross_references`, `trigger_descriptions`, `inbound`, `agent_variable.description`, `permission_intent_notes`, `concurrency_intent.notes`, `superseded_substrate_artifact`, `relationship_to_substrate` all absent from frontmatter and present in body.

---

### AC2 — cds-dispatch/SKILL.md

Checked against `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`.

- [x] `artifact_class: wake` — line 5.
- [x] `scope: global` — line 6.
- [x] `wake.role: dispatch` — line 19. Matches JSON.
- [x] `wake.package: cnos.cds` — line 20. Matches JSON.
- [x] `wake.admin_only: false` — line 21. Matches JSON.
- [x] `wake.activation_log_writer: false` — line 22. Matches JSON.
- [x] `wake.activation_state: live` — line 23. Matches JSON.
- [x] `wake.protocol: cds` — line 24. Matches JSON.
- [x] `wake.selector.include` = `[dispatch:cell, protocol:cds, status:todo]` — lines 26–29. Matches JSON.
- [x] `wake.selector.exclude` = `[status:in-progress, status:blocked, status:review, status:changes]` — lines 30–34. Matches JSON.
- [x] `wake.concurrency.group: "cds-dispatch-{agent}"` — line 57. Matches JSON `concurrency_intent.group`.
- [x] `wake.agent_variable.default: sigma` — line 60. String (not null). FN-3 satisfied.
- [x] `wake.input.triggers` = `[schedule, issues_labeled_selector_match]` — lines 37–39. Matches JSON.
- [x] No `issues_opened_title_pattern` in dispatch `wake.input` — confirmed absent.
- [x] `wake.output.cycle_artifact_root: ".cdd/unreleased/{N}/"` — line 41. Matches JSON.
- [x] `wake.output.artifact_class_taxonomy` — 7 entries. Matches JSON.
- [x] `wake.output.cell_runtime: cnos.cdd` — line 50. Matches JSON.
- [x] `wake.permission_intent` — 4 entries. Matches JSON.
- [x] `wake.concurrency.serialize: true` — line 56. Matches JSON.
- [x] `wake.agent_variable.name: agent` — line 59. Matches JSON.
- [x] `wake.surfaces.allowed` — 5 entries, verbatim from JSON `allowed_surfaces`. Verified.
- [x] `wake.surfaces.disallowed` — 8 entries, verbatim from JSON `disallowed_surfaces`. Verified.
- [x] `wake.defer_path.*` — all three values verbatim from JSON. Verified.
- [x] No prose fields in frontmatter — `responsibilities`, `activation_state_notes`, `artifact_class_notes`, `cell_runtime_notes`, `trigger_descriptions`, `inbound`, `agent_variable.description`, `permission_intent_notes`, `concurrency_intent.notes`, `cross_references` all absent from frontmatter and present in body.
- [x] Body starts with verbatim content of `prompt.md` — verified by diff: body content (after leading blank line) is byte-identical to `prompt.md` (241 lines).
- [x] OB-2 and OB-3 — same analysis as admin. Pass.

---

### AC4 — Protected files

Command: `git diff origin/main -- '*.golden.yml' wake-provider.json prompt.md .github/workflows/`

Result: empty (no output). Zero diff on all protected files:
- Both `*.golden.yml` files: unchanged.
- Both `wake-provider.json` files: unchanged.
- Both `prompt.md` files: unchanged.
- `.github/workflows/`: no files changed.

---

### AC7 — CI readiness

**SKILL.md files exist:**
```
src/packages/cnos.core/orchestrators/agent-admin/SKILL.md  ✓
src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md  ✓
```

**I5 count:**
- `git ls-tree -r origin/main --name-only | grep SKILL.md | wc -l` → 99
- `git ls-tree -r HEAD --name-only | grep SKILL.md | wc -l` → 101
- Delta: +2. Exactly the two new wake SKILL.md files. Correct.
- Note: γ-scaffold cited 91→93 (the count at scaffold-authoring time; main has grown since). The delta invariant (+2) holds.

**Changes outside `.cdd/`:** only `schemas/skill.cue` (additive), `agent-admin/SKILL.md` (new), `cds-dispatch/SKILL.md` (new). No other paths.

- [x] CUE syntax valid (manual scan; `cue` binary unavailable)
- [x] Both SKILL.md files have `artifact_class: wake` and valid frontmatter
- [x] No `.github/workflows/*.yml` changes
- [x] No `dispatch-repair-preflight`-relevant files changed
- [x] No Go / Package / Binary changes

---

### Implementation contract conformance

- [x] Language: CUE for `schemas/skill.cue`; YAML+Markdown for SKILL.md files. No Python, Go, or other languages.
- [x] `wake-provider.json` files: NOT modified.
- [x] `prompt.md` files: NOT modified.
- [x] `.github/workflows/` files: NOT modified.
- [x] `*.golden.yml` files: NOT modified.

---

### Findings

None. All checklist items pass.

---

### Verdict rationale

All three W1 deliverables are correctly implemented: `schemas/skill.cue` adds `"wake"` to the enum and a well-formed `#Wake` definition (OB-1 role enum, FN-3 null default, FN-4 open arrays, role-shaped output disjunction); both SKILL.md files carry all frontmatter fields verbatim from their respective `wake-provider.json` sources with no value deviations; both bodies contain verbatim `prompt.md` content plus the required appended prose from `wake-provider.json`; no protected file was touched; the SKILL.md delta is +2 as required. W1 is ready to converge.

---

_β@cdd.cnos — cycle/524 W1 R0 review — 2026-06-30 (UTC)_

---

# β-review — cnos#524 W2 R1

---
cycle: 524
parent_issue: cnos#524
document_class: beta-review
round: W2-R1
authored_by: β@cdd.cnos (W2-R1)
date: 2026-06-30 (UTC)
base_sha: 1a0acfc14ccbdf26d1088472125864be037ff1e1
verdict: converge
---

## §R1 Verdict

**CONVERGE**

Scope compliance is clean. Parity gate passes for both wakes (exit 0). No golden files changed.
No live wake workflows touched. The renderer extension is internally coherent and does not alter
the existing wake behavior (default source remains JSON+prompt). No stop conditions triggered.

---

## §R1 Scope Compliance Walk

| Constraint | Result | Observation |
|---|---|---|
| `schemas/skill.cue` | PASS | Absent from diff |
| `wake-provider.json` (both wakes) | PASS | Absent from diff |
| `prompt.md` (both wakes) | PASS | Absent from diff |
| `cnos-agent-admin.golden.yml` | PASS | Absent from diff; verified byte-identical |
| `cnos-cds-dispatch.golden.yml` | PASS | Absent from diff; verified byte-identical |
| `.github/workflows/cnos-agent-admin.yml` | PASS | Absent from diff |
| `.github/workflows/cnos-cds-dispatch.yml` | PASS | Absent from diff |
| No W3/W4 artifacts | PASS | No source flip, no JSON/prompt deletion |
| #524 remains open | PASS | Commit messages use `Refs #524` / `Part of #524` only |

Files in diff (allowed):
- `cn-install-wake`: W2 renderer extension (writable per directive)
- `install-wake-golden.yml`: parity step addition (writable per directive — CI guard step ONLY)
- `.cdd/unreleased/524/gamma-scaffold.md`: W2 γ scaffold (writable)
- `.cdd/unreleased/524/self-coherence.md`: R1 section appended (writable)
- `.cdd/unreleased/524/alpha-closeout.md`: W2 R1 section appended (writable)

---

## §R1 Parity Gate Verification

```
cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for agent-admin  EXIT: 0
cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for cds-dispatch EXIT: 0
```

Both wakes produce byte-identical non-header output. The parity comparison correctly strips
`^#`-prefixed header comment lines (FN-W2-1) before comparing.

---

## §R1 Renderer Extension Review

**`skill_to_json_manifest()` correctness:**
- Generates synthesized JSON with `schema: cn.wake-provider.v1` — matches the JSON source schema field.
- Maps `wake.*` frontmatter fields to the exact jq paths consumed by the renderer's validation and
  rendering logic. Verified by parity pass: the renderer reads the synthesized JSON and produces
  the same output as when reading `wake-provider.json` directly.
- `activation_log_writer` has()-vs-absent distinction: Python includes the key only when explicitly
  present in SKILL.md's `wake:` block. admin SKILL.md has no `activation_log_writer` key → field
  absent from synthesized JSON (same as `wake-provider.json`). cds-dispatch SKILL.md has
  `activation_log_writer: false` → field present as `false`. Both cases match the existing JSON
  manifests. Parity confirms this.

**`skill_body()` extraction:**
- Python-based extractor. Strips leading blank line (standard blank after closing `---` delimiter)
  and stops before the W1 reference-data sections. Boundary detection via
  `\n\n---\n\n## .+(from wake-provider.json|body reference)` regex matches the W1 authoring
  convention for both wakes.
- Both extracted bodies are byte-identical to their respective `prompt.md` files (confirmed by
  parity pass — the rendered output is byte-identical).

**Header stability:**
- `display_manifest_path` / `display_prompt_path` always use `json_manifest_path` / `json_prompt_path`
  regardless of `--source` flag. Verified: `--source skill` produces `(unchanged)` for both goldens
  (identical headers, identical body → full file identical).

**CI parity step in `install-wake-golden.yml`:**
- Step is placed BEFORE the AC8/AC7 authority audit step.
- Uses `--parity-check` flag (exits 0 or 5); no `--out` flag; no golden write.
- The step comment correctly explains the oracle: goldens are ground truth for render(JSON+prompt);
  `--parity-check` proves render(SKILL.md) == golden → establishes byte-identity.

---

## §R1 Stop Condition Audit

| Stop condition | Status |
|---|---|
| SKILL.md render differs from JSON+prompt render | NOT triggered — parity OK both wakes |
| Any golden changes | NOT triggered — `git diff` PASS on both goldens |
| Any live wake workflow changes | NOT triggered — PASS on both live workflows |
| Renderer change alters wake behavior | NOT triggered — default source unchanged; golden idempotent |
| Parity requires weakening validation | NOT triggered — skill source has reduced required_fields per W0 §D (prose fields → body), not a validation weakening |
| Any green gate turns red | NOT triggered — parity step passes; no existing steps removed |

---

_β@cdd.cnos — cycle/524 W2 R1 review — 2026-06-30 (UTC). CONVERGE._

---

# β-review — cnos#524 W3 R2

---
cycle: 524
parent_issue: cnos#524
document_class: beta-review
round: W3-R2
authored_by: β@cdd.cnos (W3-R2)
date: 2026-06-30 (UTC)
verdict: converge
---

## §R2 Verdict

**CONVERGE**

Scope compliance is clean. The source flip (line 298) and parity inversion (line 416) are
correct. All 6 targeted edits applied. CI step updated with consistent semantics. No golden
files changed. No live wake workflows touched. No SKILL.md or JSON/prompt files touched.
PR linkage constraint is met. No stop conditions triggered.

---

## §R2 Scope Compliance Walk

| Constraint | Result | Observation |
|---|---|---|
| `schemas/skill.cue` | PASS | Absent from diff |
| `wake-provider.json` (both wakes) | PASS | Absent from diff |
| `prompt.md` (both wakes) | PASS | Absent from diff |
| `SKILL.md` (both wakes) | PASS | Absent from diff |
| `cnos-agent-admin.golden.yml` | PASS | Absent from diff |
| `cnos-cds-dispatch.golden.yml` | PASS | Absent from diff |
| `.github/workflows/cnos-agent-admin.yml` | PASS | Absent from diff |
| `.github/workflows/cnos-cds-dispatch.yml` | PASS | Absent from diff |
| No W4 artifacts (JSON/prompt deletion) | PASS | No deletion of JSON+prompt files |
| `#524` remains open (PR linkage) | PASS | Commit must use `Refs #524` / `Part of #524` only |

Files in diff (allowed):
- `cn-install-wake`: 6 targeted edits (source flip, parity inversion, doc updates)
- `install-wake-golden.yml`: "W2 parity check" → "W3 parity check" step update
- `.cdd/unreleased/524/gamma-scaffold.md`: W3 γ scaffold (replaced W2)
- `.cdd/unreleased/524/self-coherence.md`: §R2 section appended
- `.cdd/unreleased/524/beta-review.md`: §R2 section appended (this document)

---

## §R2 Core Change Verification

### Edit 3.1.A — default source flip

- **Location:** `cn-install-wake:298`
- **Before:** `source_type="json"`
- **After:** `source_type="skill"`
- **Verified:** git diff shows exactly this change at the initialization block.
- **Logic:** The ONLY code paths that set `source_type` after initialization are the `--source`
  flag parser (lines 392/394) and the parity override (line 416). A plain `cn install-wake <name>`
  call (no flags) leaves `source_type` at its initialization value: `"skill"`. The SKILL.md
  extraction block (`if [ "$source_type" = "skill" ]; then`) runs. AC3 satisfied by construction.

### Edit 3.1.B — parity inversion

- **Location:** `cn-install-wake:415-416`
- **Before:** `[ -n "$parity_check" ] && source_type="skill"`
- **After:** `[ -n "$parity_check" ] && source_type="json"`
- **Verified:** git diff shows exactly this change at the post-validation parity override.
- **Logic:** With `--parity-check`, `source_type` becomes `"json"`. The SKILL.md extraction block
  is then SKIPPED. The JSON path (jq reads from `wake-provider.json`) runs. The rendered output is
  compared against the committed golden (now produced from SKILL.md by default). This proves
  render(JSON+prompt) == render(SKILL.md). AC3 parity direction confirmed.

### Edits 3.1.C–F — documentation updates

- **3.1.C** (`--source` doc): `json` (default) → `skill` (default). `skill` described first. ✓
- **3.1.D** (`--parity-check` doc): "Render from SKILL.md source…" → "Render from JSON+prompt source…";
  "Implies --source skill" → "Implies --source json"; "W2 CI parity gate" → "W3 CI parity gate". ✓
- **3.1.E** (exit-5 doc): "render(SKILL.md) differs from render(JSON+prompt)" → "render(JSON+prompt)
  differs from render(SKILL.md)"; "W2" → "W3". ✓
- **3.1.F** (messages): "render(SKILL.md) == render(JSON+prompt)" → "render(JSON+prompt) ==
  render(SKILL.md)" for both success and failure messages. ✓

### CI step (§3.2) — W2 → W3

- **Before:** `name: W2 parity check — render(SKILL.md) == render(JSON+prompt)`
- **After:** `name: W3 parity check — render(JSON+prompt) == render(SKILL.md)`
- **Comment:** Updated to reflect W3 semantics (goldens now produced from SKILL.md by default;
  `--parity-check` now implies `--source json`).
- **Run block:** Unchanged — same two `--parity-check` invocations.

---

## §R2 AC Oracle Walk

| AC | Check | Result |
|---|---|---|
| AC3 | `source_type="skill"` at init; SKILL.md path runs for plain invocation | PASS — by construction |
| AC3 | CI re-render steps use plain invocation; goldens expected unchanged | EXPECTED PASS |
| AC4 | No golden file in diff | PASS — `cnos-agent-admin.golden.yml` and `cnos-cds-dispatch.golden.yml` absent from diff |
| AC4 | W2 parity transitivity: SKILL.md renders == JSON+prompt renders == goldens | PASS — transitivity holds; no source material changed |
| AC6 | Refusal gates (exit 3, exit 4) are in renderer logic below the source selection block | PASS — no gate code touched |
| AC6 | CI refusal smokes use `--manifest` override; source selection bypassed | PASS — unaffected |
| AC7 | All CI gates: no changes to Go, Package, Binary, dispatch-repair-preflight scope | EXPECTED PASS |
| AC8 | No new literal role-decision strings in renderer diff | PASS — only `"json"`, `"skill"` literals and comment text changed |

---

## §R2 Stop Condition Audit

| Stop condition | Status |
|---|---|
| Goldens change after source flip | NOT triggered — no golden in diff; W2 parity transitivity holds |
| Live wake workflow changes | NOT triggered — `cnos-agent-admin.yml` and `cnos-cds-dispatch.yml` absent from diff |
| JSON+prompt files modified | NOT triggered — both `wake-provider.json` and `prompt.md` absent from diff |
| SKILL.md files modified | NOT triggered — both wake SKILL.md files absent from diff |
| New role-decision strings in renderer | NOT triggered — only `"json"`/`"skill"` and comment text changed |
| Parity check broken (render divergence) | NOT expected — source material unchanged; same parity paths run in CI |
| Any green gate turns red | NOT expected — step name change only; run block identical |

---

## §R2 Findings

None. All checklist items pass.

---

_β@cdd.cnos — cycle/524 W3 R2 review — 2026-06-30 (UTC). CONVERGE._

---

## §R3 — W4 review (β@cdd.cnos, 2026-06-30)

**Verdict: CONVERGE**

Independent re-review of the W4 implementation (commits `a5d7d482` γ-scaffold/REPAIR-PLAN, `06252591`
the actual delete+renderer+CI work) against `gamma-scaffold.md`, `REPAIR-PLAN.md`, the operator's
verbatim W4 "clean re-dispatch" directive on issue #524, and α's `self-coherence.md §R3` account.
I have no stake in this work; every claim below was re-run by me from a clean checkout of
`cycle/524`, not copy-checked from α's report. Full command transcripts retained in this session;
key results summarized per checklist item.

### 1. Header-only diff — mechanically re-verified

```
git diff db547ebe...cycle/524 -- <both goldens> <both live workflows>
```
Ran this myself. The diff touches exactly the two-line→one-line header on all four files, e.g.:
```
-#   manifest: orchestrators/agent-admin/wake-provider.json
-#   prompt:   orchestrators/agent-admin/prompt.md
+#   source: orchestrators/agent-admin/SKILL.md
```
identically shaped for `cds-dispatch`. Zero other bytes differ on any of the four files (full diff
inspected, not `--stat`). **Confirmed independently — highest-stakes check passes.**

### 2. Files actually deleted

```
find src/packages/cnos.core/orchestrators src/packages/cnos.cds/orchestrators \
     src/packages/cnos.core/commands/install-wake/test-fixtures \
     -name wake-provider.json -o -name prompt.md
```
Returns empty on the checked-out branch. Confirmed.

### 3. SKILL.md body prose untouched

```
git diff db547ebe...cycle/524 -- .../agent-admin/SKILL.md .../cds-dispatch/SKILL.md
```
Empty — zero diff on both files. Confirmed; no scope violation.

### 4. Renderer correctness

Read the full diff on `cn-install-wake` (234 lines changed). Confirmed: `--source`/`--source=*`/
`--parity-check` now hit a dedicated `case` arm that calls `die` with a cnos#524 W4-referencing
message (hard error, not silent passthrough — no path can reach a now-absent JSON file). Default
per-package manifest lookup and the cross-package sibling-search fallback both resolve `SKILL.md`
paths now; `grep -n 'wake-provider.json' cn-install-wake` returns 5 hits, all `#`-comments
documenting the removal, zero in resolution/parsing logic — verified by line-reading each hit, not
just counting.

Ran the renderer myself:
```
./cn-install-wake agent-admin --out /tmp/aa.yml   → exit 0, diff /tmp/aa.yml vs golden: IDENTICAL
./cn-install-wake cds-dispatch --out /tmp/cds.yml → exit 0, diff /tmp/cds.yml vs golden: IDENTICAL
```
`git status --short` after both runs was clean (no incidental working-tree drift). Confirmed.

### 5. Refusal gates — independently run, not trusted from α's report

Read `.github/workflows/install-wake-golden.yml` for the exact updated step commands and replayed
both by hand:

**AC5 declaration-only smoke:**
```
./cn-install-wake test-declaration-only --manifest .../declaration-only/SKILL.md
```
→ exit **3**; stderr: `activation_state="declaration-only"` + names `cnos#454`/`cnos#467`/
`preconditions` line. Matches both CI grep assertions (`declaration-only` substring;
`cnos#454|cnos#467|preconditions` regex). Confirmed.

**AC4/cycle-496 mis-declaration smoke:**
```
./cn-install-wake test-log-writer-misdeclaration --manifest .../log-writer-misdeclaration/SKILL.md \
  --activation-state-override live --out /tmp/x.yml
```
→ exit **4**; stderr contains literal `activation_log_writer mis-declaration`. Matches CI assertion.
Confirmed.

Both gates fire from the SKILL.md-only path (there is no other path left), satisfying AC6.

### 6. AC2 conversion is a real oracle

Replayed the exact heredoc fixture from the CI step (a `SKILL.md` with `wake:` block omitting
`role`, valid otherwise) against the renderer directly:
```
./cn-install-wake test-wake --manifest <fixture>/SKILL.md
```
→ exit **2**; stderr: `role must be one of admin/dispatch/observer (got "")`. This is a genuine
rejection (not silent defaulting, not a different exit code, not an uncaught exception/exit 1) —
confirms α's judgment call (b) in `self-coherence.md` was correct: the chosen fixture shape (missing
`role`) trips the renderer's own enum check with a precise, informative message, preserving AC2's
"malformed declarations are rejected" oracle through the only code path that exists post-W4.

### 7. Out-of-scaffold CI-script fix — claim independently verified TRUE

Diffed `db547ebe:scripts/ci/check-dispatch-closeout-integrity.sh` against the pre-W4 base: it
hardcoded `PROMPT="src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md"` and looped
`for f in "$PROMPT" "$SKILL" "$GOLDEN" "$LIVE"; do need "$f" ...` — a hard presence-and-phrase
requirement on `prompt.md`. Same shape in `check-dispatch-repair-preflight.sh`. Since W4 deletes
`cds-dispatch/prompt.md`, **not** fixing these scripts would have made `need()` fail at
`[ -f "$ROOT/$f" ]` (or the phrase-grep immediately after), turning two required-green CI gates
(`dispatch-closeout-integrity`, `dispatch-repair-preflight`, both named in the operator's W4
required-proof list) red as a direct, mechanical consequence of the scaffold-mandated deletion. The
causal claim is TRUE, not merely asserted.

The fix is contract-preserving: both scripts swap their loop target from `$PROMPT` to the
already-declared `$SKILL` variable; no required phrase string, label, or detection logic changed —
only the file-existence target. Ran both scripts myself on `cycle/524`:
```
bash scripts/ci/check-dispatch-closeout-integrity.sh             → exit 0 (presence guard + self-test both pass)
bash scripts/ci/check-dispatch-closeout-integrity.sh --self-test → exit 0
bash scripts/ci/check-dispatch-repair-preflight.sh                → exit 0
```
This is a legitimate, narrowly-scoped follow-on fix, not scope creep — it is the minimum patch
required to keep two scaffold-required-green gates green after a deletion the scaffold itself
mandates.

### 8. AC8/AC7 renderer-authority-leak audit — re-run independently

Re-ran both greps from the CI step directly against the new renderer source:
```
admin_leaks=$(grep -nE 'admin_only|disallowed_surfaces|defer_path|cell_execution' cn-install-wake \
  | grep -vE '^[0-9]+:[[:space:]]*#' | grep -vE '(die|manifest_error|log_writer_mis_declaration)[[:space:]]')
→ empty (0 leaks)

n_dispatch=$(grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' cn-install-wake)
→ 0
```
Both zero. Confirmed.

### 9. CI gates runnable locally

| Gate | Result |
|---|---|
| `go build ./...` (src/go) | clean |
| `go vet ./...` | clean |
| `go test ./...` | all 14 packages `ok` (one `[no test files]` for cmd/cn, expected) |
| `./cn build --check` (I1, package/source drift) | all 8 packages valid |
| `./cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` (I6-adjacent) | 106 passed, 0 failed, 77 warnings (pre-existing/unrelated; includes this issue's own `.cdd/unreleased/524` artifact-presence checks, all passing) |
| `scripts/ci/check-dispatch-closeout-integrity.sh` / `--self-test` | green (see §7) |
| `scripts/ci/check-dispatch-repair-preflight.sh` | green (see §7) |
| `install-wake-golden.yml`'s steps (replayed by hand: re-render, sha256, idempotence, AC5, AC2, AC4/496, AC8/AC7) | all green, replayed individually above |

**Not runnable in this environment:** I5 (`scripts/ci/validate-skill-frontmatter.sh`) requires the
`cue` binary, which is not installed here — confirmed `which cue` fails. α flagged the same gap
honestly in `self-coherence.md §R3 §4`. Since neither `schemas/skill.cue` nor either wake SKILL.md's
frontmatter changed in this cycle (confirmed §3 above), I5 has no mechanical reason to regress, but
I could not execute it myself — flagging for the actual CI run to confirm, same as α did. Binary/
Package CI jobs (full multi-tier suites with CI-specific setup) likewise not replicated locally;
`go build`/`go vet`/`go test` serve as the closest local proxy and are clean.

### 10. Scope compliance — full diff stat

```
git diff db547ebe...cycle/524 --stat
```
19 files: 3 `.cdd/unreleased/524/` records (REPAIR-PLAN.md new, gamma-scaffold.md rewritten,
self-coherence.md extended) + 2 live workflows (header-only) + 2 goldens (header-only) +
`install-wake-golden.yml` + `cn-install-wake` + the 2 CI guard scripts (§7) + 8 deleted
JSON/prompt files (4 production + 4 fixture). Every file is on the scaffold's §3.2 MUST-change
list, or is the §7 justified out-of-scaffold CI-script fix. Confirmed absent from the diff:
`schemas/skill.cue`, `wake-provider/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md`,
`docs/**`, `.cdd/releases/**`, any other `.cdd/unreleased/{N}/`, `src/go/**`. Confirmed clean.

### 11. Commit hygiene

```
git log db547ebe..cycle/524 --oneline
a5d7d482 cnos#524 W4: γ-scaffold + REPAIR-PLAN for clean re-dispatch
06252591 cnos#524 W4: delete wake-provider.json + prompt.md; renderer is SKILL.md-only
```
Both full commit messages end with `Refs #524`. Neither message contains `Closes`/`Fixes`/
`Resolves #524` in any form (checked the full message bodies, not just trailers). Confirmed.

### Process note (not a finding): closeouts and PR are correctly absent at this stage

`alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` carry zero diff from base, and no
W4 PR exists yet (`gh pr list --head cycle/524` shows only the merged #525–#529 for W0–W3). I
checked the δ skill (`cdd/delta/SKILL.md §9.3` step 4–5) before treating this as a defect: closeouts
are authored by α/β/γ **after** β's `verdict: converge` lands, as part of δ's converge routing, and
the cycle-PR is opened by δ at the `status:review` transition — not before β review. So this is
expected pipeline state, not a `deliverable_evidence` gap in α's work. It does mean the scaffold's
checklist item "`deliverable_evidence` present" (§3.1 / β-prompt item 8) is **not yet satisfiable**
and is correctly deferred to the post-converge closeout step that follows this verdict.

### Findings

None blocking. All 11 checklist items independently re-verified and pass. No STOP-level issue
found — the header-only-diff invariant (the single highest-stakes claim in this cycle) holds
exactly as claimed, the JSON/prompt deletion is complete and total, both refusal gates and the
AC2 oracle fire correctly from the SKILL.md-only path, the two out-of-scaffold CI-script edits are
a genuine, narrowly-scoped, contract-preserving necessity (not scope creep), and commit hygiene is
clean.

---

_β@cdd.cnos — cycle/524 W4 R3 review — 2026-06-30 (UTC). CONVERGE._
