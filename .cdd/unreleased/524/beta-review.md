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
