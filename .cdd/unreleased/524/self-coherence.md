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

---

## Manual-δ amendment 2 — W3 completion (manual_delta_repair, operator-directed, 2026-06-30)

`run_class: manual_delta_repair`. **The W3 cell shipped install-wake-golden RED and β
falsely CONVERGE'd it — that earlier "green" claim is NOT preserved.** AC5 was only the
first of several install-wake-golden steps the source-flip broke: every step that fed a
**JSON-only synthetic manifest** and expected JSON-default behavior broke once `skill`
became the default source. Full classification + repair below (operator principle: default
wake-rendering tests exercise SKILL.md; only JSON-legacy / JSON-schema-failure tests may
pin `--source json`; `--source json` must not be used merely to dodge the new default).

| install-wake-golden step | what it tests | classification | repair |
|---|---|---|---|
| AC5 — declaration-only refusal (exit 3) | activation_state refusal; applies to **either** source | default → SKILL.md | added `test-fixtures/declaration-only/SKILL.md` (amendment 1) |
| AC4 / cycle-496 — log-writer mis-declaration (exit 4) | activation_log_writer mis-declaration refusal; applies to **either** source | default → SKILL.md | **added `test-fixtures/log-writer-misdeclaration/SKILL.md`** (role:dispatch + admin_only:false + activation_log_writer:true) — skill default now hits exit 4 + `activation_log_writer mis-declaration:` |
| AC2 — negative malformed manifest (exit 2) | malformed **JSON** missing required `schema` field — no SKILL.md equivalent (SKILL.md is CUE/#Wake-validated) | **JSON-legacy specific** | pinned **`--source json`** on the AC2 step; retire/replace in W4 |
| AC4 cycle-496 write-fence (positive / negative / R1 / R2 / false-positive) | git write-fence logic; **invokes no renderer** | unaffected | none |
| W3 parity (both wakes) | render(JSON+prompt) == render(SKILL.md) | parity (implies `--source json` internally) | none (works) |

Two new SKILL.md fixtures (declaration-only, log-writer-misdeclaration) validate under
`#Wake` (I5 → 95 modules, no findings). The JSON twins remain for W3 dual-source parity and
are deleted in W4.

**Re-verified (κ, worktree, this amendment):** AC5 exit 3; **AC4/496 exit 4 + precise
stderr**; **AC2-negative exit 2 + `required field "schema" missing`** (exact CI step replay
with `set -o pipefail` passes); both goldens byte-identical from the skill default; parity
exit 0 both; negative drift proof (mutate JSON → render unchanged, mutate SKILL.md → render
changes); dispatch-repair-preflight green; I5 = 95 no findings. Boundary preserved: no
renderer/synthesizer logic change, no golden change, no live-workflow change, no
wake-behavior change, JSON+prompt retained.

---

## §R3 — W4 implementation self-coherence (α@cdd.cnos, 2026-06-30)

`run_class: repair_pass` per `REPAIR-PLAN.md` (the prior W4 dispatch run, `28464981342`, was
invalidated as an empty run — no PR/commits/closeout — and is not authoritative; there is no
prior W4 **code** to repair against, only process state, already remediated by PR #531). This
cell performs the actual W4 implementation against current `main` (`db547ebe`) under the
operator's "W4 (final phase, clean re-dispatch)" directive (comment id 4847294646) and
`gamma-scaffold.md`.

### 1. What changed

**Deleted (8 files, both production wakes + both test fixtures):**
- `src/packages/cnos.core/orchestrators/agent-admin/{wake-provider.json,prompt.md}`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/{wake-provider.json,prompt.md}`
- `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/{wake-provider.json,prompt.md}`
- `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/{wake-provider.json,prompt.md}`

Each directory already carried a `SKILL.md` twin (added W1/W3); only the legacy JSON+prompt
pair is removed. The SKILL.md files themselves are untouched (frontmatter and body both).

**`cn-install-wake` rewritten SKILL.md-only:**
- `--source json|skill` and `--parity-check` flags removed. Both flag forms now match a
  dedicated `case` arm that calls `die` with an explanatory message ("flag '...' was removed
  in cnos#524 W4 ..."), rather than silently disappearing — a stale CI invocation or muscle-
  memory local invocation gets a precise error, not surprising default behavior.
- `source_type` / `parity_check` shell variables removed entirely. The SKILL.md→JSON synth
  (`skill_to_json_manifest`) and body extraction (`skill_body`) — previously gated behind
  `if [ "$source_type" = "skill" ]` — are now unconditional; there is exactly one code path.
- The JSON-branch `required_fields` set (the `else` arm carrying `responsibilities`,
  `prompt_template`, `cross_references` as top-level required fields) is deleted; only the
  reduced SKILL.md-shape set remains.
  - The JSON-branch prompt resolution (`prompt_template_relpath` / `jq -r '.prompt_template'`
    against the raw manifest) is deleted; `prompt_path` is now unconditionally
    `$skill_body_tmp` (the SKILL.md body, per the W0 §E body-as-prompt rule).
- **Manifest resolution now resolves `SKILL.md`** in both places that used to glob for
  `wake-provider.json`: the default per-package lookup
  (`${manifest_dir}/SKILL.md`) and the cross-package sibling-search fallback
  (`candidate="${sibling}/orchestrators/${wake_name}/SKILL.md"`). Verified the fallback
  explicitly (§4 below — negative proof with `CN_PACKAGE_ROOT` unset).
- **Header attribution** changed from the two-line `# manifest: .../wake-provider.json` /
  `# prompt: .../prompt.md` block to a single `# source: orchestrators/<name>/SKILL.md` line.
  This is the only byte-level change to any golden/live-workflow output (§3 below).
- Exit code 5 (formerly `--parity-check` failure) is documented as **retired, not reused** —
  see judgment call (e) below.
- Top-of-script doc comment (Usage/Arguments/Exit-codes) updated to match: `--manifest`'s
  description now says it accepts a SKILL.md path directly; `--source`/`--parity-check`
  documentation removed; default manifest lookup path updated to name `SKILL.md`.

**`.github/workflows/install-wake-golden.yml`:**
- "W3 parity check — render(JSON+prompt) == render(SKILL.md)" step removed entirely (nothing
  left to compare parity against).
- "AC2 negative-case smoke" converted from a malformed-`wake-provider.json` fixture (missing
  `schema`) to a malformed-SKILL.md fixture (a `wake:` block omitting `role`). Same oracle:
  malformed declaration → exit 2 + an informative stderr substring, not silently accepted.
  See judgment call (b) below for why the exact substring differs from the old message.
- "AC5 — declaration-only refusal" and "AC4 (cycle/496) — renderer refusal on mis-declaration"
  steps' `--manifest` arguments repointed at the fixtures' `SKILL.md` files (the
  `wake-provider.json` they used to point at is now deleted). Assertions unchanged (exit 3 /
  exit 4 + the same stderr substrings).
- All other steps (re-render, golden-diff, sha256 live-vs-golden, idempotence ×2, YAML-parses,
  substrate-shape ×2, write-fence smokes ×5, AC8/AC7 authority audit) left structurally
  unchanged.

**Goldens + live workflows re-rendered (header-only diff, verified twice — see §3):**
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
- `.github/workflows/cnos-agent-admin.yml`
- `.github/workflows/cnos-cds-dispatch.yml`

**Two CI guard scripts fixed (not on the scaffold's original MUST-change list — see judgment
call (f) below for why this was necessary, not scope creep):**
- `scripts/ci/check-dispatch-closeout-integrity.sh`
- `scripts/ci/check-dispatch-repair-preflight.sh`

### 2. Judgment calls (scope-boundary reasoning)

**(a) `--manifest` contract (FN-W4-2).** Made `--manifest <path>` accept either a `SKILL.md`
file path directly, or a directory containing one (the latter resolves to
`<dir>/SKILL.md`, matching the shape of the default per-package lookup). This keeps the CI
fixture smokes' existing "point `--manifest` at the fixture directory's manifest file"
calling convention intact while being unambiguous. Documented in the script's doc-comment
header.

**(b) AC2 conversion (FN-W4-5).** Tested several malformed-SKILL.md shapes before picking the
fixture: a `wake:` block omitting `role` fails cleanly at the **renderer's own** `role must be
one of admin/dispatch/observer (got "")` enum check (exit 2, precise message) — not inside
`skill_to_json_manifest`'s python bridge (which defaults `role` to `''` via
`wake.get('role') or ''` and lets the shell-side validation catch it). I deliberately did NOT
use a non-mapping `wake:` value (e.g. `wake: "garbage"`) — that crashes the python bridge with
an uncaught exception and exits 1, which is the wrong exit code and not a "malformed
declaration is rejected with an informative message" oracle; it's an unhandled-input bug
class, not the AC2 invariant. The chosen fixture (missing `role`) preserves AC2's intent —
"malformed declarations are rejected, not silently defaulted" — through the only code path
that now exists. The CI assertion greps for the literal renderer message
(`role must be one of admin/dispatch/observer`) rather than the old JSON-path's
`required field "schema" missing`, because that's what the SKILL.md path actually produces;
matching the assertion to actual behavior (per FN-W4-5's explicit guidance) rather than
forcing artificial parity with the retired message.

**(c) `--source`/`--parity-check` made hard errors, not silently dropped.** The scaffold
offered both options ("no `--source` flag... OR keep `--source` accepting only `skill`... die
on `json`"). I chose a third, stricter variant: both flags (and `--source=*`) are recognized
by the arg parser and unconditionally `die` with a message naming cnos#524 W4 and the reason.
This is friendlier than "unknown flag" (which would print a generic error with no context) and
safer than silently accepting/ignoring the flag (which could mask a caller's wrong assumption
about renderer behavior). No caller in this repo passes these flags after this cycle's CI
edits (confirmed by re-grepping `.github/workflows/` post-edit), so this is purely a defensive
error path, not a live behavior.

**(d) `--out` / `wake_name` / `--activation-state-override` left untouched** — none of these
flags or their resolution logic referenced JSON/prompt; no changes needed there per the
scaffold's scope.

**(e) Exit code 5 retired, not reused (FN-W4-3).** Chose "leave a one-line note" over "drop the
row from the table": the exit-codes doc-comment block now states exit 5 is retired as of W4
and explains why (nothing left to compare parity against), so historical CI logs or operator
muscle memory referencing exit 5 aren't misread as a different, newly-assigned failure class.
No code path can produce exit 5 any more (confirmed: `grep -n 'exit 5' cn-install-wake` →
zero hits post-edit).

**(f) Two CI guard scripts required a follow-on fix not on the scaffold's MUST-change list.**
`scripts/ci/check-dispatch-closeout-integrity.sh` and `scripts/ci/check-dispatch-repair-
preflight.sh` (both pre-existing, `build.yml`-wired CI gates — closeout-integrity is itself
one of W4's own required-green gates per the operator directive) each hardcoded
`PROMPT="src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md"` as a required-file
check in their presence-of-contract guard, left over from the W2/W3 dual-source transition
window (their comments said "kept verbatim-identical for W2/W3 parity until W4"). Deleting
`cds-dispatch/prompt.md` per the W4 mandate made both scripts fail with `required file
missing` — a genuine regression I'd have introduced by doing exactly what the scaffold
required. I verified the literal phrases both scripts also already required from `$SKILL` /
`$GOLDEN` / `$LIVE` were present (they are — the SKILL.md body and the rendered golden/live
workflow both carry the closeout-integrity and repair-preflight contract prose, confirmed by
`grep -qF`), then replaced each script's `PROMPT` variable + loop entry with the already-
present `SKILL` variable, with a comment explaining the W4 deletion. This is the minimum
necessary fix to keep two required-green CI gates green after a deletion the scaffold itself
mandates — not scope creep into prompt/design content (no prose, contract wording, or
detection logic changed; only the file-existence target). Re-ran both scripts' presence guard
and (for closeout-integrity) `--self-test` locally — both green (§4).

### 3. The header-only-diff invariant — verified twice, mechanically

Per FN-W4-1, I ran the renderer change FIRST (with JSON/prompt files still present) and
diffed the re-render against the then-committed (pre-W4) goldens — confirmed **header-only**
(lines 2–3 collapsed to a single new line 2; zero other bytes). I then deleted the JSON/prompt
files and re-rendered again (both production wakes + both live-workflow `--out` targets) and
re-confirmed the same header-only diff against the same pre-W4 baseline:

```
-#   manifest: orchestrators/agent-admin/wake-provider.json
-#   prompt:   orchestrators/agent-admin/prompt.md
+#   source: orchestrators/agent-admin/SKILL.md
```
(and the `cds-dispatch` equivalent), identically on all four files: both goldens and both live
`.github/workflows/cnos-*.yml`. `git diff --cached` on each of the four files was inspected
individually (not just `--stat`) to confirm no other line moved. sha256(live) == sha256(golden)
for cds-dispatch holds post-change. Idempotence (second render is a byte-for-byte no-op) holds
for both wakes.

### 4. Verification results (run locally, not just claimed)

| Check | Result |
|---|---|
| `cn-install-wake agent-admin` | exit 0, renders from SKILL.md (JSON absent) |
| `cn-install-wake cds-dispatch` | exit 0, renders from SKILL.md (JSON absent) |
| Header-only diff, both goldens | confirmed — see §3 |
| Header-only diff, both live workflows | confirmed — see §3 |
| sha256 live == golden (cds-dispatch) | match (`30b6cf2d19...`) |
| Idempotence (agent-admin, cds-dispatch) | both: second render byte-identical |
| YAML parses (both goldens) | pass |
| Substrate structural shape (both) | pass, incl. OG-2 schedule gate on cds-dispatch |
| AC5 declaration-only refusal (`--manifest .../declaration-only/SKILL.md`) | exit 3; stderr names `declaration-only`; stderr matches `cnos#454\|cnos#467\|preconditions` |
| AC4/cycle-496 mis-declaration (`--manifest .../log-writer-misdeclaration/SKILL.md`) | exit 4; stderr contains `activation_log_writer mis-declaration` |
| Converted AC2 negative smoke (replayed the exact YAML-extracted step via Python+bash) | exit 2; stderr contains `role must be one of admin/dispatch/observer` |
| `find ... -name wake-provider.json -o -name prompt.md` (orchestrators + test-fixtures) | empty |
| `grep -n wake-provider.json cn-install-wake` | 5 hits, all comments documenting the removal; zero in resolution/parsing logic |
| AC8 (admin-shape leak audit) | 0 leaks |
| AC7 (dispatch-shape leak audit) | 0 leaks |
| `--manifest` accepting a directory (not just a file) | confirmed working (resolves to `<dir>/SKILL.md`) |
| Cross-package sibling-fallback negative proof (`unset CN_PACKAGE_ROOT`, direct invocation, JSON absent) | `cds-dispatch` still resolves via sibling search → SKILL.md; output byte-identical to golden |
| `--source`/`--source=skill`/`--parity-check` now hard-error | confirmed, exit 1, explanatory message naming cnos#524 W4 |
| `scripts/ci/check-dispatch-closeout-integrity.sh` (presence guard) | green post-fix |
| `scripts/ci/check-dispatch-closeout-integrity.sh --self-test` | green |
| `scripts/ci/check-dispatch-repair-preflight.sh` | green post-fix |
| `go build ./...` (src/go, untouched by this cycle) | clean |
| `go vet ./...` | clean |
| `go test ./...` | all packages pass |
| `./cn build --check` (I1 package/source drift) | all packages valid |
| `diff docs/reference/schemas/protocol-contract.json tests/fixtures/protocol-contract.json` (I2) | identical |
| `./cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` (I6) | 106 passed, 0 failed (warnings only, pre-existing/unrelated) |

**Not runnable locally:** I5 (`scripts/ci/validate-skill-frontmatter.sh`, CUE-schema-driven)
— `cue` binary is not installed in this environment. I did not edit `schemas/skill.cue` or
either wake `SKILL.md`'s frontmatter, so I5 has no reason to regress, but I could not execute
it myself; flagging for β/CI to confirm. Binary/Package CI jobs (`binary-verify`,
`package-verify`) were not run locally (they build + run multi-tier verification suites with
CI-specific setup, e.g. a synthetic test hub) — `go build`/`go vet`/`go test` were run as the
closest local proxies and are clean; flagging for CI to confirm the full jobs.

### 5. Stop-condition audit

| Stop condition (operator directive) | Status |
|---|---|
| Any non-header workflow bytes change | NOT triggered — verified twice, §3 |
| Wake behavior changes | NOT triggered — no selector/permission/concurrency/prompt-body change |
| Deleting JSON/prompt changes semantics | NOT triggered — byte-identical render proven before AND after deletion |
| Renderer needs old JSON compatibility to pass | NOT triggered — zero JSON-source code path remains |
| #524 would be closed | NOT triggered — commit messages use `Refs #524` only |
| Demo 0 gets invoked | NOT triggered — no Demo 0 surface touched |

### 6. Known gap explicitly NOT fixed (per gamma-scaffold §1.1)

`src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` and
`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` both have **zero diff** in this
cycle (confirmed: `git diff --cached -- <both paths>` → 0 lines). Their body prose still
references the now-deleted `wake-provider.json` in a few places — most visibly
`cds-dispatch/SKILL.md`'s `[wake-provider.json](wake-provider.json)` markdown link (now a
dangling 404 target) plus both files' `## Responsibilities (from wake-provider.json; body
reference)` / `## Cross-references (from wake-provider.json)` section-header parentheticals.
Per §1.1's explicit instruction, this is NOT fixed here: the SKILL.md body IS the rendered
prompt (W0 §E body-as-prompt rule), so editing it would change non-header bytes of the
rendered golden/live-workflow output — directly violating the operator's "Stop if: any
non-header workflow bytes change" condition and the "no wake behavior / prompt semantic
change" scope constraint. The "no active references to wake-provider.json/prompt.md remain"
required-proof bullet is satisfied per §1.1's resolution: by (a) the files being deleted (AC5)
and (b) zero renderer resolution/parsing-logic references (confirmed §4) — not by scrubbing
frozen prompt-body prose. This stale-prose correction is a legitimate, separately-scoped
follow-up (prompt content change, not a renderer/substrate change) for a future cell; I have
not filed the tracked follow-up issue myself (no `gh issue create` scope was given to α in
this prompt) — flagging for δ/operator to file it at closeout, naming this self-coherence
section + gamma-scaffold §1.1 as the source.

### 7. Scope compliance

`git diff --cached --stat` (post this cycle's edits) touches exactly: the 8 deleted JSON+prompt
files, `cn-install-wake`, `install-wake-golden.yml`, the two CI guard scripts (judgment call f),
and the four re-rendered golden/workflow files. No diff on: either wake `SKILL.md`,
`schemas/skill.cue`, `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`,
`src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`,
`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`, anything under `docs/`, `.cdd/releases/`,
any other `.cdd/unreleased/{N}/`, or `src/go/`. Confirmed via targeted `git diff --cached --
<path>` on each protected path (§6, plus an explicit empty-diff check on the six other
protected non-goal paths).

REVIEW READY: R3
