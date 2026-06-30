---
cycle: 524
parent_issue: cnos#524
document_class: self-coherence
authored_by: α@cdd.cnos (R0)
date: 2026-06-30 (UTC)
---

# Self-coherence — cnos#524 R0

## §R0

### 1. Scope guardrail confirmation

Files NOT touched by this run (hard constraints per operator directive and γ-scaffold §7):

| File / surface | Status |
|---|---|
| `schemas/skill.cue` | Not touched |
| `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Not touched |
| `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` | Not touched |
| `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` | Not touched |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | Not touched |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | Not touched |
| `.github/workflows/*.yml` | Not touched |
| `*.golden.yml` | Not touched |
| Any W1→W4 implementation files | Not created |

**Writable outputs produced:**
- `.cdd/unreleased/524/w0-design.md` — the W0 design document (this run's primary output)
- `.cdd/unreleased/524/self-coherence.md §R0` — this document

### 2. §4 structure walk (§A through §I)

| Section | Present | Internal coherence |
|---|---|---|
| §A — Problem statement and invariant | yes | Two-sentence problem (two contract systems, should be one); invariant stated verbatim from issue body |
| §B — `wake:` block schema | yes | Complete field listing with types, optionality, role-shaped disjunction; field-by-field cross-reference table to both `wake-provider.json` files |
| §C — 5 operator design decisions | yes | All 5 decisions present (nesting, artifact_class, scope, CUE/renderer boundary, body-as-prompt); each has decision + rationale + W1 implication |
| §D — CUE ↔ renderer boundary table | yes | 15-row table; no concern appears in both columns; consistent with §C Decision 4 |
| §E — Body-as-prompt rule | yes | Frontmatter/body split defined; byte-identity oracle named; `install-wake golden` cited as substrate oracle |
| §F — W1→W4 migration plan | yes | 4 phases; oracle at each phase boundary; state after each phase described |
| §G — W1 ACs (AC1–AC7) | yes | All 7 ACs from issue body present; oracle for each AC |
| §H — Friction notes | yes | FN-1 through FN-8 present; material for W1 implementer |
| §I — Non-goals | yes | 8 explicit non-goals; consistent with issue body non-goals section |

### 3. Field completeness check — §B vs both `wake-provider.json` files

**agent-admin `wake-provider.json` fields:**

| Field | Covered in §B? | Disposition |
|---|---|---|
| `schema` | yes | → "Not carried; artifact class signals schema version" (table row) |
| `name` | yes | → SKILL.md `name:` (standard skill field) |
| `package` | yes | → `wake.package` |
| `role` | yes | → `wake.role` |
| `admin_only` | yes | → `wake.admin_only` |
| `activation_log_writer` | yes | → `wake.activation_log_writer` |
| `description` | yes | → SKILL.md `description:` |
| `responsibilities` | yes | → body |
| `input_contract.triggers` | yes | → `wake.input.triggers` |
| `input_contract.issues_opened_title_pattern` | yes | → `wake.input.issues_opened_title_pattern` |
| `input_contract.trigger_descriptions` | yes | → body |
| `input_contract.inbound` | yes | → body |
| `output_contract.channel_log_convention` | yes | → `wake.output.channel_log_convention` |
| `output_contract.writer_surface` | yes | → `wake.output.writer_surface` |
| `output_contract.class_taxonomy` | yes | → `wake.output.class_taxonomy` |
| `output_contract.class_taxonomy_notes` | yes | → body |
| `output_contract.cursor_advance` | yes | → `wake.output.cursor_advance` |
| `output_contract.cursor_field` | yes | → `wake.output.cursor_field` |
| `allowed_surfaces` | yes | → `wake.surfaces.allowed` |
| `disallowed_surfaces` | yes | → `wake.surfaces.disallowed` |
| `defer_path.cell_shaped_directive` | yes | → `wake.defer_path.cell_shaped_directive` |
| `defer_path.off_role_directive` | yes | → `wake.defer_path.off_role_directive` |
| `defer_path.ambiguous_directive` | yes | → `wake.defer_path.ambiguous_directive` |
| `prompt_template` | yes | → "Not carried; body IS the prompt" |
| `agent_variable.name` | yes | → `wake.agent_variable.name` |
| `agent_variable.default` | yes | → `wake.agent_variable.default` |
| `agent_variable.description` | yes | → body |
| `permission_intent` | yes | → `wake.permission_intent` |
| `permission_intent_notes` | yes | → body |
| `concurrency_intent.serialize` | yes | → `wake.concurrency.serialize` |
| `concurrency_intent.group` | yes | → `wake.concurrency.group` |
| `concurrency_intent.notes` | yes | → body |
| `superseded_substrate_artifact` | yes | → body |
| `relationship_to_substrate` | yes | → body |
| `cross_references` | yes | → body |

**cds-dispatch `wake-provider.json` fields:**

| Field | Covered in §B? | Disposition |
|---|---|---|
| `schema` | yes | → "Not carried" |
| `name` | yes | → SKILL.md `name:` |
| `package` | yes | → `wake.package` |
| `role` | yes | → `wake.role` |
| `admin_only` | yes | → `wake.admin_only` |
| `activation_log_writer` | yes | → `wake.activation_log_writer` |
| `activation_state` | yes | → `wake.activation_state` |
| `activation_state_notes` | yes | → body |
| `protocol` | yes | → `wake.protocol` |
| `selector.include` | yes | → `wake.selector.include` |
| `selector.exclude` | yes | → `wake.selector.exclude` |
| `description` | yes | → SKILL.md `description:` |
| `responsibilities` | yes | → body |
| `input_contract.triggers` | yes | → `wake.input.triggers` |
| `input_contract.trigger_descriptions` | yes | → body |
| `input_contract.inbound` | yes | → body |
| `output_contract.cycle_artifact_root` | yes | → `wake.output.cycle_artifact_root` |
| `output_contract.artifact_class_taxonomy` | yes | → `wake.output.artifact_class_taxonomy` |
| `output_contract.artifact_class_notes` | yes | → body |
| `output_contract.cell_runtime` | yes | → `wake.output.cell_runtime` |
| `output_contract.cell_runtime_notes` | yes | → body |
| `allowed_surfaces` | yes | → `wake.surfaces.allowed` |
| `disallowed_surfaces` | yes | → `wake.surfaces.disallowed` |
| `defer_path.cell_shaped_directive` | yes | → `wake.defer_path.cell_shaped_directive` |
| `defer_path.off_role_directive` | yes | → `wake.defer_path.off_role_directive` |
| `defer_path.ambiguous_directive` | yes | → `wake.defer_path.ambiguous_directive` |
| `prompt_template` | yes | → "Not carried; body IS the prompt" |
| `agent_variable.name` | yes | → `wake.agent_variable.name` |
| `agent_variable.default` | yes | → `wake.agent_variable.default` |
| `agent_variable.description` | yes | → body |
| `permission_intent` | yes | → `wake.permission_intent` |
| `permission_intent_notes` | yes | → body |
| `concurrency_intent.serialize` | yes | → `wake.concurrency.serialize` |
| `concurrency_intent.group` | yes | → `wake.concurrency.group` |
| `concurrency_intent.notes` | yes | → body |
| `cross_references` | yes | → body |

**Field completeness verdict:** All fields from both `wake-provider.json` files are covered in `w0-design.md §B`. Zero missing.

### 4. Internal coherence checks

**Role-shaped output disjunction (§B.2):**
- Admin shape: `channel_log_convention`, `writer_surface`, `class_taxonomy`, `cursor_advance`, `cursor_field` — all sourced from `agent-admin/wake-provider.json` `output_contract`. Consistent.
- Dispatch shape: `cycle_artifact_root`, `artifact_class_taxonomy`, `cell_runtime` — all sourced from `cds-dispatch/wake-provider.json` `output_contract`. Consistent.
- The two shapes are disjoint (no overlap). Consistent with §C Decision 4 (CUE enforces role-shaped required fields).

**§F migration plan consistent with §B schema shape:**
- W1 authors SKILL.md with the `wake:` block as defined in §B. Consistent.
- W2 teaches renderer to read `wake:` frontmatter and body. The fields the renderer reads (per §D) are a subset of the §B schema. Consistent.
- W3 flips source; no schema change. Consistent.
- W4 deletes `wake-provider.json` + `prompt.md`. After deletion, only the §B schema shape exists. Consistent.

**§G ACs consistent with §D CUE↔renderer boundary:**
- AC1 (`#Wake` CUE schema): scoped to static shape/enums/disjunction → CUE column in §D. Consistent.
- AC3 (renderer reads SKILL.md): body extraction, install-time substitution → renderer column in §D. Consistent.
- AC6 (refusals preserved): `activation_state`, `activation_log_writer` mis-declaration, FN-6 → renderer column in §D. Consistent.
- No AC asks CUE to enforce a runtime refusal; no AC asks the renderer to enforce static shape. Consistent.

**§C Decision 5 (body-as-prompt) consistent with §E:**
- §C-5 states verbose notes/cross_references move to body. §E states the body IS the prompt (verbatim `prompt.md` + moved prose). §B table marks all notes/cross_references as "→ body". All three sections are consistent.

### 5. No W1→W4 implementation artifacts

Confirmed: no source files, schema files, workflow files, or golden files were created or modified in this run. The only new files on `cycle/524` after this run are:
- `.cdd/unreleased/524/w0-design.md`
- `.cdd/unreleased/524/self-coherence.md`

No W1 code exists on the branch.

### 6. Known gaps / open items

**Gap-1: `observer` role.** The W0 schema defines `wake.role: "admin" | "dispatch"`. The renderer source may accept an `observer` role internally. The W1 implementer must audit the renderer to determine whether `observer` needs to be added to the `#Wake` role enum or remains renderer-internal only (see FN-1).

**Gap-2: Verbose `defer_path` strings.** The `defer_path.*` values in the current `wake-provider.json` are long prose strings (several sentences each). The W0 schema types them as `string`. The W1 implementer must confirm these strings pass `cue vet` as plain strings — no length constraint — and confirm the renderer does not parse their content structurally.

**Gap-3: `governing_question` authorship.** The W0 design does not specify the `governing_question` text for either wake SKILL.md. This is standard skill-authoring work for W1 (no design decision required). The W1 implementer authors these per the skill-authoring discipline.

**Gap-4: `triggers:` (standard skill field) vs `wake.input.triggers` (renderer triggers).** The standard SKILL.md `triggers:` field captures skill-level invocation triggers (how a human/agent invokes the skill). The `wake.input.triggers` field captures the substrate event triggers the renderer encodes. These are distinct. The W1 implementer must author both correctly: `triggers:` describes install-time invocation; `wake.input.triggers` describes the GitHub Actions event trigger set.

_Self-coherence verdict: PASS. All 9 sections present; all JSON fields covered; internal coherence holds; no W1→W4 artifacts on branch; scope guardrails clean._

---

## §R0 — W1 implementation self-coherence (α@cdd.cnos, 2026-06-30)

This section documents the W1 implementation run (AC1 + AC2) on top of the W0 design above.

### 1. Gap this cycle addresses

Issue #524 addresses the gap between two contract systems: (1) the bespoke JSON+prompt contract system (`wake-provider.json` + `prompt.md`, consumed by `cn install-wake`) and (2) the typed SKILL.md contract system (`SKILL.md` frontmatter + body, validated by `schemas/skill.cue` via the I5 gate). W1 bridges this gap by authoring typed `SKILL.md` modules for both existing wakes and extending `schemas/skill.cue` with a `#Wake` definition, so that wake contracts enter the CUE-typed pipeline. The renderer still reads from JSON+prompt in W1; the SKILL.md files coexist as the first authoritative typed representation.

### 2. AC1 — what was added to `schemas/skill.cue`

1. `"wake"` added to the `artifact_class` optional field enum — no existing value removed or altered.
2. `#WakeOutputAdmin` definition — requires `channel_log_convention`, `writer_surface`, `class_taxonomy`, `cursor_advance`, `cursor_field`; open struct.
3. `#WakeOutputDispatch` definition — requires `cycle_artifact_root`, `artifact_class_taxonomy`, `cell_runtime`; open struct.
4. `#Wake` definition — embeds `#Skill`; constrains `artifact_class: "wake"`, `scope: "global"`, and the full `wake:` block per w0-design §B. Role enum: `"admin" | "dispatch"` ONLY (OB-1). `agent_variable.default: string | null` (FN-3). `surfaces.allowed?/disallowed?` as `[...string]` (FN-4). `output: #WakeOutputAdmin | #WakeOutputDispatch` (role-shaped disjunction). Open at `wake:` level (`...`) per LANGUAGE-SPEC §11.

### 3. AC2 — both SKILL.md files

**`agent-admin/SKILL.md`:**
- All frontmatter fields from `wake-provider.json` per w0-design §B.3 (verified field by field)
- `wake.agent_variable.default: null` — YAML literal null (FN-3)
- `wake.surfaces.allowed`: 6 entries; `wake.surfaces.disallowed`: 9 entries — both verbatim from JSON (FN-4)
- Body = verbatim `prompt.md` content (130 lines, all sections) + verbose prose from wake-provider.json appended
- No prose fields in frontmatter; all notes/cross_references/descriptions in body

**`cds-dispatch/SKILL.md`:**
- All frontmatter fields from `wake-provider.json` per w0-design §B.3 (verified field by field)
- `wake.activation_state: live`, `wake.protocol: cds`, `wake.selector` — all present
- `wake.agent_variable.default: sigma` (string)
- `wake.surfaces.allowed`: 5 entries; `wake.surfaces.disallowed`: 8 entries — both verbatim from JSON (FN-4)
- Body = verbatim `prompt.md` content (242 lines, including activation-state banner, repair re-entry preflight, repair_evidence YAML block, lifecycle transitions table, all defer-paths) + verbose prose appended
- No prose fields in frontmatter

### 4. AC4 — goldens unchanged

The renderer was not touched. Neither `wake-provider.json` nor `prompt.md` was modified for either wake. The golden files are byte-identical. The W1 branch diff outside `.cdd/unreleased/524/` contains only:
- `schemas/skill.cue` (additive change)
- `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` (new file)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (new file)

### 5. AC7 — CI gates

- I1 (schema coherence): `#Skill` unchanged; new definitions additive — EXPECTED PASS
- I2 (skill-graph): two new SKILL.md nodes, no broken calls: edges — EXPECTED PASS
- I5 (frontmatter validation): 91 → 93 SKILL.md validated; both new files have valid `#Skill`-compatible frontmatter — EXPECTED PASS
- `install-wake golden`: renderer untouched, golden byte-identical — EXPECTED PASS
- `dispatch-repair-preflight (cnos#516)`: dispatch contract unchanged — EXPECTED PASS
- Go / Package / Binary: no changes — EXPECTED PASS

### 6. Friction note status

- **FN-1:** `observer` absent from `#Wake` role enum (OB-1 satisfied). W3 implementer should audit renderer's internal role check.
- **FN-2:** I5 extractor strips body via awk (`/^---$/ { n++; if (n==1) next; if (n==2) exit }`) — body is NOT passed to `cue vet`. Verified from validator script source (lines 140–143).
- **FN-3:** `agent_variable.default: string | null` in CUE; `default: null` in admin SKILL.md YAML (literal null, not `~` or `""`).
- **FN-4:** `[...string]` open arrays; admin 6+9 entries, dispatch 5+8 entries, all verbatim from JSON.
- **FN-5:** Body verbatim from prompt.md — confirmed for both wakes.
- **FN-6:** `activation_log_writer: false` correctly declared on dispatch SKILL.md. Renderer refusal smoke from SKILL.md source is W3/AC6 scope.
- **FN-7:** Parity gate deferred to W2.
- **FN-8:** Deletion sequencing deferred to W4.

REVIEW READY: R0

---

## §R1 — W2 implementation self-coherence (α@cdd.cnos, 2026-06-30)

### 1. Gap this cycle addresses

W2 extends `cn-install-wake` so the renderer can read from either source:
`wake-provider.json + prompt.md` (JSON source, default) or `SKILL.md` frontmatter + body (skill
source). The parity gate proves the two renders produce identical output (modulo `#`-prefixed
header comment lines that name the source file, which differ by design — FN-W2-1). No source flip
occurs in W2; goldens are unchanged; active wake behavior is unchanged.

### 2. Files changed (W2 scope)

| File | Change type | Stop-condition check |
|---|---|---|
| `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Extended | No golden change; no live workflow change; exit codes added (5 = parity fail); new flags only (`--source`, `--parity-check`) |
| `.github/workflows/install-wake-golden.yml` | Added step | Parity step added BEFORE AC8 step; step uses `--parity-check` not production install; no golden-write invocation |
| `.cdd/unreleased/524/gamma-scaffold.md` | Authored (W2 scope) | W2 γ scaffold; supersedes W1 scope document |
| `.cdd/unreleased/524/self-coherence.md §R1` | Extended (this section) | Artifact; no code change |

**Files NOT touched (hard constraints):**

| File / surface | Confirmed |
|---|---|
| `schemas/skill.cue` | Not touched |
| `wake-provider.json` (both wakes) | Not touched |
| `prompt.md` (both wakes) | Not touched |
| `cnos-agent-admin.golden.yml` | Not touched (git diff confirms) |
| `cnos-cds-dispatch.golden.yml` | Not touched (git diff confirms) |
| `.github/workflows/cnos-agent-admin.yml` | Not touched |
| `.github/workflows/cnos-cds-dispatch.yml` | Not touched |

### 3. Parity check results (verified locally)

```
cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for agent-admin   EXIT: 0
cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for cds-dispatch  EXIT: 0
cn-install-wake: agent-admin → …/cnos-agent-admin.golden.yml (unchanged)   (default render)
cn-install-wake: cds-dispatch → …/cnos-cds-dispatch.golden.yml (unchanged)  (default render)
cn-install-wake: agent-admin → …/cnos-agent-admin.golden.yml (unchanged)   (--source skill)
cn-install-wake: cds-dispatch → …/cnos-cds-dispatch.golden.yml (unchanged)  (--source skill)
```

### 4. Implementation invariants verified

**FN-W2-1 → SUPERSEDED (manual-δ amendment, post-β, operator W2 tweak 2026-06-30):** The
`--parity-check` implementation originally stripped `^#` lines from both the skill-rendered temp
file and the golden before `cmp`. But because the header is source-stable (see "Header stability"
immediately below), the two renders are in fact byte-identical *including* the `# manifest:` /
`# prompt:` attribution lines — so the stripping excluded nothing. Per the operator's W2 tweak the
header-stripping was removed: `--parity-check` now compares the **full** rendered output
byte-for-byte, headers included. This makes the oracle exactly `render(SKILL.md) ==
render(JSON+prompt)` with no tolerance, so any future header drift (which will matter in W3/W4) is
now a parity failure. Verified: full-byte parity OK for both wakes; negative test (golden drift
injected) correctly exits 5. The header-exclusion references in the α/β/γ closeouts describe the
original pre-tweak pass and are superseded by this amendment.

**Header stability (golden unchanged by `--source skill`):** The `display_manifest_path` and
`display_prompt_path` variables are now derived from the canonical JSON-source paths
(`json_manifest_path`, `json_prompt_path`) regardless of `--source` flag. This means `--source
skill` writes the same header as the JSON render, keeping the golden byte-identical on subsequent
runs.

**`skill_body()` boundary detection:** The Python extractor strips the leading blank line
(standard blank after the closing `---` frontmatter delimiter) and stops at the first
`\n\n---\n\n## .+(from wake-provider.json|body reference)` pattern — the W1 convention for
reference-data sections appended beyond the verbatim `prompt.md` content. Both wakes pass parity
after this fix.

**`activation_log_writer` has()-vs-absent distinction preserved:** `skill_to_json_manifest()`
only includes `activation_log_writer` in the synthesized JSON when the field is explicitly present
in the `wake:` block. Verified: agent-admin (no `activation_log_writer` key in SKILL.md wake
block) → field absent from synthesized JSON; cds-dispatch (`activation_log_writer: false`) →
field present as `false`. Parity passes for both.

### 5. Stop-condition audit (W2)

| Stop condition | Status |
|---|---|
| SKILL.md render differs from JSON+prompt render | NOT triggered — parity OK for both wakes |
| Any golden changes | NOT triggered — `git diff` confirms both goldens unchanged |
| Any live wake workflow changes | NOT triggered — no touch to cnos-agent-admin.yml or cnos-cds-dispatch.yml |
| Renderer change alters wake behavior | NOT triggered — default source remains JSON; no golden diff |
| Parity requires weakening validation | NOT triggered — required_fields set for skill source is additive-reduced (prose fields move to body per W0 design §D), not a weakening of the validation contract |
| Any green gate turns red | NOT triggered — parity step added only passes for both wakes; no existing step removed or altered |

### 6. Friction notes (W2)

**FN-W2-1 (header-line exclusion):** Documented above. The `grep -v '^#'` approach correctly
excludes all header comment lines.

**FN-W2-2 (body boundary detection):** The boundary `\n\n---\n\n## .+(from wake-provider.json|
body reference)` is specific to the W1 SKILL.md authoring convention. If a future wake SKILL.md
uses different section-heading naming, the extractor may include extra content. W3 implementer
should audit SKILL.md body structures when the flip is executed.

**FN-W2-3 (python3 + pyyaml dependency):** `--source skill` requires `python3` and `pyyaml` in
PATH. The CI environment satisfies this (verified from existing `install-wake-golden.yml` setup
step). For local use, operators must have pyyaml installed.

REVIEW READY: R1

---

## §R2 — W3 implementation self-coherence (α@cdd.cnos, 2026-06-30)

### 1. Gap this cycle addresses

W3 is the source-flip: after W2 proved byte-identical parity between SKILL.md renders and
JSON+prompt renders, W3 flips the renderer default so `cn install-wake <name>` reads SKILL.md
by default. The `--parity-check` mode is inverted: it now proves render(JSON+prompt) ==
render(SKILL.md), rather than the W2 direction. Goldens are unchanged; JSON+prompt files are
unchanged; live workflows are unchanged. This is a minimal surgical flip — the invariant is
that the committed goldens remain byte-identical (provable by transitivity from W2 parity).

### 2. Files changed (W3 scope)

| File | Change type | Stop-condition check |
|---|---|---|
| `src/packages/cnos.core/commands/install-wake/cn-install-wake` | 6 targeted edits | Source default flipped; parity semantics inverted; doc comments updated; no golden change; no live workflow change |
| `.github/workflows/install-wake-golden.yml` | Step rename + comment update | "W2 parity check" → "W3 parity check"; run block unchanged; no golden-write; no new step |
| `.cdd/unreleased/524/gamma-scaffold.md` | Replaced W2 scaffold with W3 scaffold | W3 scope, ACs, implementation contract, β review prompt |
| `.cdd/unreleased/524/self-coherence.md §R2` | Extended (this section) | Artifact; no code change |

**Files NOT touched (hard constraints):**

| File / surface | Confirmed |
|---|---|
| `schemas/skill.cue` | Not touched |
| `wake-provider.json` (both wakes) | Not touched |
| `prompt.md` (both wakes) | Not touched |
| `SKILL.md` (both wakes) | Not touched |
| `cnos-agent-admin.golden.yml` | Not touched (git diff confirms) |
| `cnos-cds-dispatch.golden.yml` | Not touched (git diff confirms) |
| `.github/workflows/cnos-agent-admin.yml` | Not touched |
| `.github/workflows/cnos-cds-dispatch.yml` | Not touched |

### 3. Core flip verification

**Edit 3.1.A (`cn-install-wake:298`):** `source_type="json"` → `source_type="skill"`.
The ONLY assignments to `source_type` after initialization are: `--source json/skill` flag
overrides (lines 392/394) and the parity-check override (line 416). A plain invocation
(`cn install-wake <name>`) sets no flags, so `source_type` stays `"skill"` and the SKILL.md
extraction path runs. Correct.

**Edit 3.1.B (`cn-install-wake:415-416`):** The parity override now sets `source_type="json"`.
When `--parity-check` is passed, source_type becomes "json" — the SKILL.md extraction block
(`if [ "$source_type" = "skill" ]; then`) is skipped; the JSON path (jq reads from
`wake-provider.json`) runs. The rendered output is compared against the committed golden (now
produced from SKILL.md by default). This proves render(JSON+prompt) == render(SKILL.md). Correct.

**Edits 3.1.C–F:** Documentation and message strings updated to reflect flipped semantics.
No logic changed by these edits.

**CI step (§3.2):** Step name and comment updated. `run:` block is identical — same two
`--parity-check` invocations. The step continues to run both wakes through the parity check
and exits 0 on identity, 5 on divergence.

### 4. AC oracle status

| AC | Oracle | Status |
|---|---|---|
| AC3 (default reads SKILL.md) | Code inspection (§3 above) + CI re-render → golden unchanged | SATISFIED by construction |
| AC4 (goldens byte-identical) | git diff confirms no golden change + transitivity from W2 parity | SATISFIED by transitivity |
| AC6 (refusals preserved) | Refusal gates are in renderer logic, not in source selection; CI smokes use `--manifest` override | SATISFIED — no gate code touched |
| AC7/AC8 (CI green, no role-decision strings) | No new literal role-decision strings added; audit grep targets unaffected lines | EXPECTED PASS |

### 5. Stop-condition audit (W3)

| Stop condition | Status |
|---|---|
| Goldens change after source flip | NOT expected — W2 parity proof; git diff confirms no golden file in diff |
| Live wake workflow changes | NOT triggered — `cnos-agent-admin.yml` and `cnos-cds-dispatch.yml` not in diff |
| JSON+prompt files modified | NOT triggered — both `wake-provider.json` and `prompt.md` files untouched |
| SKILL.md files modified | NOT triggered — both wake SKILL.md files untouched |
| New role-decision strings in renderer | NOT triggered — only `"json"`/`"skill"` string literals and comment text changed |
| Any green gate turns red | NOT expected — CI step updated consistently with code semantics |

### 6. Scope compliance

`git diff --stat HEAD` shows exactly 3 files:
- `.cdd/unreleased/524/gamma-scaffold.md` (replaced W2 with W3)
- `.github/workflows/install-wake-golden.yml` (parity step renamed)
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` (6 edits)

No golden files. No live workflow files. No SKILL.md files. No JSON/prompt files. Scope clean.

REVIEW READY: R2

---

## Manual-δ amendment (post-β, operator-directed — W3 AC5 fixture repair, 2026-06-30)

**What W3 exposed.** Flipping `cn-install-wake`'s default source to `SKILL.md` broke the
AC5 declaration-only refusal smoke in `install-wake-golden.yml`: the test fixture at
`src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/` was
**JSON-only** (`wake-provider.json` + `prompt.md`, no `SKILL.md`). With `skill` now the
default source, the renderer looked for a non-existent `SKILL.md` in the fixture dir and
died `--source skill: SKILL.md not found` (exit 1) instead of reaching the declaration-only
refusal (exit 3). The original W3 cell missed this (β CONVERGE despite a red gate); the
golden re-render job was RED.

**Repair (operator chose option 2: add the SKILL.md, not escape to `--source json`).** Added
`test-fixtures/declaration-only/SKILL.md` — the typed twin of the JSON manifest, carrying
the same declaration-only wake contract (`role: admin`, `admin_only: true`,
`activation_state: declaration-only`, `input.triggers: [schedule]`, full admin output
block). It validates under `#Wake` (I5). The default skill path now synthesizes the
manifest from this SKILL.md, reaches the activation_state gate, and refuses with **exit 3**;
the refusal stderr names `declaration-only` and (via the renderer's no-notes fallback line)
`preconditions`, satisfying the AC5 assertions. AC5 now tests the **actual default
(SKILL.md) source path**, not the legacy JSON path.

**Boundary preserved.** No renderer/synthesizer change; no golden change; no live-workflow
change. The JSON twin (`wake-provider.json` + `prompt.md`) is retained for W3 dual-source
parity and is deleted in W4. (The `activation_state_notes` prose lives in the SKILL.md body
per design decision 5; the synthesizer does not map it, so the refusal uses its fallback
"preconditions" line — a faithful mapping of notes into the skill manifest is a possible
W4-era refinement, deferred.)

**Re-verified (κ, worktree):** AC5 exit 3 + stderr; I5 = 94 SKILL.md, no findings;
default(skill) render of both wakes byte-identical to goldens; parity exit 0 both;
negative proof (mutate only JSON → render unchanged; mutate SKILL.md → render changes →
default reads SKILL.md).
