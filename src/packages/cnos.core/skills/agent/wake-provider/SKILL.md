---
name: wake-provider
description: The wake-provider declaration contract. Defines what a cnos package declares to provide a wake (identity, responsibilities, input/output contract, allowed/disallowed surfaces, defer-path, prompt template) and the split between what the package declares vs what the substrate renderer materializes. The agent-admin wake provider (cnos.core, sibling `orchestrators/agent-admin/wake-provider.json`) is the reference instance; cnos.cds's dispatch wake provider (Sub 4 of cnos#467) — the first concrete-protocol dispatch wake — will be authored against this contract alone.
artifact_class: skill
kata_surface: none
governing_question: What does a cnos package declare to provide a wake, and which fields of the declaration are the package's authority vs the substrate renderer's authority?
visibility: public
parent: agent
triggers:
  - wake provider
  - wake declaration
  - cn wake install
  - package-owned wake
scope: task-local
inputs:
  - none (contract reference)
outputs:
  - the wake-provider declaration contract (required + optional fields)
  - the substrate-rendering boundary (package authority vs renderer authority)
  - the schema name (`cn.wake-provider.v1`) consumed by the Sub 3 renderer
requires:
  - none (foundational; consumed by `cn wake install` per cnos#450, by the agent-admin wake provider per cnos#467 Sub 2, and by future per-package dispatch wake providers per cnos#467 Subs 4+)
---

# Wake-provider declaration contract

A cnos package declares a **wake provider** to ship a long-form work shape that an external substrate (today: GitHub Actions) materializes into a real workflow file. The package authors the *abstract* declaration (identity, responsibilities, surfaces, prompt template); the renderer authors the *substrate-specific* artifact (today: YAML; tomorrow possibly other schedulers).

This skill is the contract: what a wake provider declaration MUST contain, what it MAY contain, and where the boundary lies between the package's authority and the renderer's authority.

## Core Principle

**A wake provider declaration is substrate-agnostic data + a substrate-agnostic prompt template. The substrate is the renderer's authority, not the package's.**

Two coupled invariants govern the contract:

1. **Package authority.** The package declares *what the wake is*: identity, role, allowed/disallowed surfaces, defer-path, the prompt the agent runs. These are fixed per the package's design — changing the substrate does not change them.
2. **Renderer authority.** The renderer materializes the declaration into the current substrate (today: GitHub Actions YAML in `.github/workflows/`). Permission allowlist *shape*, secret-name *binding*, trigger *encoding*, run identity *binding*, concurrency *encoding* — these are substrate-specific and the renderer owns them.

A package that emits substrate-specific YAML in its declaration has confused the boundary; a renderer that decides what responsibilities the wake's prompt enumerates has confused the boundary in the other direction. The two errors compound: package-side substrate emission locks the package to one scheduler; renderer-side responsibility decision means the package's role is defined by the renderer's defaults, not the package's design.

The failure modes this contract structurally prevents are listed in §6.

## Algorithm

A wake provider declaration is **one machine-consumable manifest + one substrate-agnostic prompt template**, both under a single directory the package nominates. The renderer:

1. **Discovers** the declaration by walking the installed package's content classes (`commands/`, `orchestrators/`, or any future class declared in `cn.package.json`) for files matching the wake-provider schema name.
2. **Parses** the manifest (schema `cn.wake-provider.v1`); validates required fields are present and well-typed.
3. **Inlines** the prompt template (a separate `.md` file alongside the manifest) into the rendered substrate artifact's `prompt:` field (or the substrate-specific equivalent).
4. **Materializes** the substrate artifact at the substrate's canonical path (today: `.github/workflows/cnos-{wake-name}.yml`).
5. **Reports** install / uninstall idempotently — repeated `cn wake install {wake}` is a no-op when the materialized artifact already matches.

The contract is what the renderer reads. The package declares; the renderer materializes. This skill defines what the renderer is allowed to expect from the declaration and what the package is allowed to expect from the renderer.

---

## 1. Define

### 1.1. Identify the parts

A wake-provider declaration has these parts:

- **Manifest** — a JSON file declaring the schema (`cn.wake-provider.v1`), identity, intent, contracts, and surfaces. Machine-consumable.
- **Prompt template** — a Markdown file the renderer inlines into the rendered substrate artifact's prompt field. Substrate-agnostic prose.
- **Optional README** — human-readable cross-references; never read by the renderer.

A substrate-rendered artifact (the renderer's output) is **not** part of the declaration. The declaration is what the package authors; the artifact is what the renderer emits. The two MUST NOT live in the same directory of the package's source tree (the package authors declarations, not artifacts).

- ❌ Package ships `orchestrators/{wake}/cnos-{wake}.yml` alongside the manifest (the package is now substrate-coupled)
- ✅ Package ships `orchestrators/{wake}/wake-provider.json` + `orchestrators/{wake}/prompt.md` only; renderer emits `.github/workflows/cnos-{wake}.yml` at install time

### 1.2. Articulate how they fit

The manifest is **the schema-bound data** the renderer can dispatch on. The prompt template is **the inline body** the renderer substitutes into the substrate's prompt-receiving field. Together they are *the wake provider* — a package-owned, substrate-agnostic declaration of what wake gets installed.

The package authority/renderer authority split (§Core Principle) makes the split concrete:

| Lives in the declaration (package authority) | Lives in the renderer (substrate authority) |
|---|---|
| Wake name, role, target agent variable | Workflow `name:` field encoding |
| Responsibilities enumeration (what the agent does) | Permission allowlist syntax (e.g. GitHub Actions `permissions:` block) |
| Allowed / disallowed surfaces | Secret name binding (e.g. `${{ secrets.X }}`) |
| Input contract (logical trigger taxonomy: schedule, issues-opened-with-title-match, ...) | Substrate trigger encoding (e.g. `on: { schedule: [...], issues: { types: [opened] } }`) |
| Output contract (channel log entry per AGENT-ACTIVATION-LOG-v0; cursor advance) | Run-identity binding (e.g. `bot_name` / `bot_id` in claude-code-action) |
| Defer-path for off-role directives | Concurrency / serialization encoding |
| Prompt template body | Inlining the prompt into the substrate's prompt field |
| Cross-references (skills the wake invokes; doctrine the wake cites) | (none — references are package authority) |

The contract is the rule that makes the split machine-checkable: if a field belongs in the right column above, it MUST NOT appear in the manifest; if a field belongs in the left column, it MUST appear in the manifest.

- ❌ Manifest carries a `permissions:` field with GitHub Actions YAML keys
- ✅ Manifest carries a `permission_intent:` field that names the *logical* permission (e.g. `contents.write`, `issues.write`); renderer maps to the substrate's encoding

### 1.3. Name the failure mode

The wake-provider contract fails through **substrate leakage**: the package declares substrate-specific content in its manifest or prompt (GitHub Actions YAML, runner names, secret names with `${{ }}` interpolation, hard-coded workflow paths). When the substrate changes (a future scheduler; a different runner; a different secret manager), every package's declaration must be edited. The package no longer owns its wake — the substrate owns the package.

The opposite failure is **renderer overreach**: the renderer decides what responsibilities the wake enumerates, what defer-path applies, what the prompt body says. The package's intent is reduced to a thin manifest of "just install the wake" and the renderer fills in defaults. Two packages get the same wake regardless of their distinct design.

Both failures are observable:

- **F1 — Substrate leakage.** _Symptom:_ `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action' {manifest} {prompt}` returns non-zero (except in the AC7-equivalent "relationship to existing substrate-bound wake" section if one exists). _Fix:_ remove substrate-specific tokens; replace with logical equivalents the renderer maps.
- **F2 — Renderer overreach.** _Symptom:_ The manifest omits one of the required fields below (§2.1) and the renderer inserts a default. _Fix:_ the renderer's contract rejects manifests with missing required fields rather than defaulting silently; the package must declare.
- **F3 — Schema drift.** _Symptom:_ A field changes meaning without changing schema name; older renderers misparse newer manifests. _Fix:_ schema is versioned (`cn.wake-provider.v1`); breaking changes bump the suffix.
- **F4 — Substrate-rendered artifact in the source tree.** _Symptom:_ The package ships a rendered `.yml` alongside the manifest. _Fix:_ §1.1 boundary — the source tree carries the declaration; the renderer emits at install time at the substrate's canonical path.

---

## 2. Unfold

### 2.1. Required manifest fields

Every wake-provider manifest MUST declare these fields, in this rough order (the renderer parses by name, not position, but consistent ordering aids human review):

| Field | Type | Purpose |
|---|---|---|
| `schema` | string, MUST equal `"cn.wake-provider.v1"` | Identifies this file as a wake-provider declaration consumable by `cn wake install`. The renderer dispatches by schema. |
| `name` | string, kebab-case, repo-unique | The wake's identifier (e.g. `agent-admin`, `cds-dispatch`). The renderer derives the substrate artifact name from this (today: `.github/workflows/cnos-{name}.yml`). |
| `package` | string | The owning package (e.g. `cnos.core` for the agent-admin wake; `cnos.cds` for the CDS dispatch wake; cnos.cdr / cnos.cdw for their respective dispatch wakes). Authoritative for ownership; the renderer verifies the manifest lives under the named package's directory tree. Note: dispatch wakes belong to concrete protocol packages (cds/cdr/cdw); `cnos.cdd` is the generic cell-runtime framework and does not own a dispatch wake itself. |
| `role` | string, enum: `admin` \| `dispatch` \| `observer` | The wake's role class. `admin` = channel sync + admin-only directive handling, never executes cells. `dispatch` = claims cells matching its protocol selector and runs the protocol's runtime. `observer` = reads-only; emits reports but does not act. The agent-admin wake declares `admin`. |
| `responsibilities` | array of strings, non-empty | The enumerated capabilities the wake's prompt commits to. The renderer does not interpret these (they are prose for the agent); the package authors them so they are observable in the declaration without reading the prompt. |
| `admin_only` | boolean | When `true`, the wake MUST NOT execute cells. The prompt template is required to enforce this; the renderer is required to surface it (e.g. via permission allowlist restriction in the materialized substrate artifact). `role: admin` implies `admin_only: true`. |
| `input_contract` | object | Logical trigger taxonomy. Required keys: `triggers` (array; e.g. `["schedule", "issues_opened_title_match"]`), `inbound` (string; logical description, e.g. `"home thread + open issues"`). Renderer maps to substrate-specific trigger encoding. |
| `output_contract` | object | Required keys: `channel_log_convention` (string; canonical doc path, e.g. `"docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md"`), `class_taxonomy` (array of allowed `class:` values for the channel entry's frontmatter), `cursor_advance` (boolean; whether the wake advances the channel cursor). |
| `allowed_surfaces` | array of strings | What the wake MAY write. Each entry is a path glob or named surface. Patterns use `{agent}` placeholder for the per-install agent name. |
| `disallowed_surfaces` | array of strings | What the wake MUST NOT write. Must explicitly include `.github/workflows/` and (for admin wakes) the literal phrase `cell execution` as a named surface. |
| `defer_path` | object | Required keys: `cell_shaped_directive` (string; what to do when a cell-shaped directive arrives — e.g. `"defer to relevant protocol:{P} dispatch wake if installed; surface to operator if not"`), `off_role_directive` (string; what to do when a directive falls outside the wake's role; admin wakes defer to dispatch, dispatch wakes defer to admin or operator). |
| `prompt_template` | string | Relative path to the prompt-template markdown file (e.g. `"prompt.md"`). The renderer reads this file and inlines its body into the substrate artifact's prompt field. |
| `cross_references` | object | Required keys for traceability. Each value is an array of citation strings (package-path, issue ref, or canonical-doc path). At minimum: `architecture` (e.g. `["cnos#467"]`), `predecessors` (e.g. `["cnos#468"]`), `consumed_skills` (the skills the prompt invokes, e.g. `["cnos.core/skills/agent/activate", "cnos.core/skills/agent/attach"]`), `consumed_conventions` (e.g. `["docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md"]`). |

### 2.2. Optional manifest fields

These fields MAY appear; the renderer treats their absence as the documented default.

| Field | Type | Default | Purpose |
|---|---|---|---|
| `description` | string | (none) | Human-readable one-line summary; ignored by the renderer; helpful in code review. |
| `agent_variable` | object | `{ "name": "agent", "default": null }` | The per-install agent identity binding (e.g. `--agent sigma` at `cn wake install` time). Object keys: `name` (the variable name the prompt template substitutes), `default` (the fallback if not provided at install time; `null` means required). |
| `input_contract.issues_opened_title_pattern` | string | (wake's `name`) | The literal title substring the substrate matches when the `issues_opened_title_match` trigger fires. When omitted, the renderer uses the wake's `name`. Declare explicitly to preserve a legacy trigger phrase across a cutover, or whenever the title pattern intentionally differs from the wake name. |
| `permission_intent` | array of strings | (renderer default per substrate) | Logical permissions the wake needs (e.g. `["contents.write", "issues.write", "pull_requests.write", "id_token.write"]`). Renderer maps to substrate-specific permission encoding. When omitted, renderer uses substrate's minimum-for-role defaults. |
| `concurrency_intent` | object | `{ "serialize": false }` | Object keys: `serialize` (boolean; if true, the substrate is asked to serialize firings to prevent channel-surface races); `group` (string; logical concurrency group name; renderer maps to substrate-specific syntax). |
| `superseded_substrate_artifact` | string | (none) | A path to an existing hand-written substrate-bound artifact this wake replaces at renderer cutover (e.g. `.github/workflows/claude-wake.yml`). Documents the migration boundary; the renderer MAY emit a warning if the named artifact still exists after install. |
| `relationship_to_substrate` | string | (none) | Free-form prose describing this declaration's relationship to existing substrate-bound artifacts. Used in conjunction with `superseded_substrate_artifact` to document Sub 3 cutover semantics for cnos#467 Sub 2 (agent-admin) and parallel cycles. |

### 2.3. Prompt template requirements

The prompt template is a Markdown file the renderer inlines verbatim into the substrate artifact's prompt-receiving field (e.g. claude-code-action's `prompt:`). The template MUST:

1. **Be substrate-agnostic.** No GitHub Actions YAML, no `${{ }}` interpolation syntax, no hard-coded runner names, no `.github/workflows/` paths (except in a documented "relationship to existing substrate-bound wake" section when the wake supersedes one — and even there, the reference is descriptive, not interpolative).
2. **Enumerate the responsibilities** named in the manifest's `responsibilities` field. The prompt is the agent's view; the manifest is the renderer's view; the two MUST agree.
3. **Enumerate the defer-path** named in the manifest's `defer_path` field. The wake's behavior when a directive falls outside its role is the most important behavioral commitment; it must be explicit in the prompt where the agent reads it.
4. **For admin wakes (`role: admin`, `admin_only: true`):** explicitly forbid cell execution in unambiguous language (e.g. "MUST NOT execute cells under any circumstance"; "Never execute cell-shaped directives inline"). This is what makes the admin/dispatch boundary observable in the prompt.
5. **Cite the consumed skills + conventions** named in the manifest's `cross_references` field. The agent reading the prompt should be able to load every skill the prompt expects without consulting another document.
6. **Substitute the `agent_variable`** (or its default) where the prompt names "the agent identity" (e.g. `"Activate and attach as https://github.com/usurobor/cn-{agent}"`). The renderer performs this substitution at install time using the value the operator passed (e.g. `--agent sigma` → `cn-sigma`).

The renderer does NOT rewrite the prompt body. It substitutes named variables (today: `{agent}`); it inlines the result; it does not edit prose. Prose changes belong in the package's source tree.

### 2.4. Substrate-rendering target

Today's substrate is **GitHub Actions** with `anthropics/claude-code-action@v1` as the agent-execution carrier. Specifically:

- **Substrate artifact path:** `.github/workflows/cnos-{wake-name}.yml`
- **Substrate trigger encoding:** YAML `on:` block (e.g. `on: { schedule: [...], issues: { types: [opened] } }`)
- **Substrate permission encoding:** YAML `permissions:` block at workflow level
- **Substrate run-identity binding:** `bot_name` + `bot_id` inputs to `claude-code-action`
- **Substrate secret binding:** `${{ secrets.* }}` interpolations
- **Substrate prompt-receiving field:** `prompt:` input to `claude-code-action`

Tomorrow's substrate may be different (a different scheduler, a different agent-execution carrier, a different secret manager). The package's declaration does not change when the substrate changes — only the renderer changes. This is the contract's leverage: future substrate migrations are renderer-local, not package-tree-wide.

The mapping from declaration fields to substrate-specific encoding is the renderer's responsibility (Sub 3 of cnos#467); this skill names what the declaration MUST provide so the renderer can produce a working substrate artifact without operator-side improvisation.

### 2.5. The package-authors-vs-renderer-materializes split (canonical table)

This split is the contract's core invariant. The table makes it grep-able:

| Declared in manifest / prompt (package authority) | Emitted by renderer (substrate authority) |
|---|---|
| `name`, `package`, `role`, `admin_only` | `name:` field of substrate artifact |
| `responsibilities` (enumerated array) | (n/a — responsibilities are prose for the agent) |
| `input_contract.triggers` (logical taxonomy) | `on:` block (substrate trigger encoding) |
| `input_contract.issues_opened_title_pattern` (literal title substring) | substituted into the job-level `if:` for the `issues_opened_title_match` trigger (substrate gating encoding) |
| `input_contract.inbound` (logical description) | (n/a — agent reads inbound from the surface itself) |
| `output_contract.channel_log_convention` (path to canonical doc) | (n/a — agent reads convention from the cited path) |
| `output_contract.class_taxonomy` (array of allowed `class:` values) | (n/a — agent emits per the prompt + convention) |
| `output_contract.cursor_advance` (boolean) | (n/a — agent behavior, not substrate behavior) |
| `allowed_surfaces`, `disallowed_surfaces` | (n/a — declaration only; agent enforces via prompt) |
| `defer_path` | (n/a — declaration only; agent enforces via prompt) |
| `permission_intent` (logical permissions) | `permissions:` block (substrate permission encoding) |
| `concurrency_intent` | `concurrency:` block (substrate concurrency encoding) |
| `agent_variable` | (substituted by renderer at install time using `--agent` flag) |
| Prompt template body | inlined into substrate's prompt field |
| `cross_references` | (n/a — declaration only; for traceability) |
| (nothing) | substrate-bound secret names, runner names, action versions, `${{ }}` interpolations, default-branch policy, OIDC permission encoding |

The right column is the renderer's territory. A package that declares any of the right column has confused the boundary.

### 2.6. Authoring a new wake provider against this contract

To declare a new wake provider in package `{pkg}`:

1. **Pick a directory under the package's content classes.** `commands/install-{wake-name}/` or `orchestrators/{wake-name}/` are the two valid surfaces today (per cnos#467 active design constraint — no new `wakes/` class until `commands/` + `orchestrators/` prove insufficient).
2. **Create `wake-provider.json`** in that directory; populate every required field from §2.1; populate optional fields from §2.2 as needed.
3. **Create `prompt.md`** alongside; author the prompt per §2.3 requirements; substitute named variables only at the `{agent}`-placeholder positions.
4. **Optionally create `README.md`** for human-readable cross-references; the renderer ignores it.
5. **Do NOT touch the substrate.** No edits under `.github/workflows/`. The renderer emits the substrate artifact at `cn wake install` time (Sub 3 of cnos#467).
6. **Verify substrate-agnostic:** `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action' wake-provider.json prompt.md` MUST return 0 *except* for an explicit `relationship_to_substrate` field or a "relationship to existing claude-wake.yml" section that names the artifact being superseded (descriptive reference, not substrate emission).

The agent-admin wake provider (cnos.core, sibling `orchestrators/agent-admin/wake-provider.json`) is the reference instance of this contract. The cnos.cds dispatch wake provider (Sub 4 of cnos#467) will pattern-copy from it: `cnos.cds/orchestrators/cds-dispatch/wake-provider.json` + `prompt.md`, declaring `role: dispatch` instead of `role: admin`. (cnos.cds is the concrete software protocol; cnos.cdd is the generic cell-runtime framework that cds and the other concrete protocol packages — cnos.cdr, cnos.cdw — use under the hood. Dispatch wakes belong to concrete protocols, not to the framework.)

---

## 3. Rules

### 3.1. Schema is versioned; breaking changes bump the suffix

The schema `cn.wake-provider.v1` is the contract version. A breaking change (renaming a required field, changing a field's type, removing a required field) bumps to `cn.wake-provider.v2` and the renderer adds a v2 handler. Older v1 manifests continue to work via the v1 handler; the renderer dispatches by the `schema` field value.

- ❌ Add a new required field to v1 silently; existing renderers reject manifests they previously accepted
- ✅ Add a new required field as v2; existing manifests declare v1 and continue to work

### 3.2. The manifest is data; the prompt is prose

The manifest is machine-consumable structured data. The prompt is human-readable prose for the agent to read. Do not put prose in the manifest (a free-form `description` is the only exception); do not put structured data in the prompt (use named variables only at the `{agent}`-placeholder positions).

- ❌ Manifest carries a paragraph of explanation in a `notes` field; renderer treats as data
- ✅ Manifest carries fielded data; prose lives in `prompt.md` or `README.md`

### 3.3. Substrate-specific content does not appear in the declaration

The declaration is substrate-agnostic by contract. Substrate-specific content (GitHub Actions YAML, runner names, secret names with `${{ }}`, hard-coded workflow paths, action versions) belongs in the renderer.

The single carve-out: a `relationship_to_substrate` field or "relationship to existing claude-wake.yml" prose section MAY name a substrate-bound artifact this declaration supersedes — that is a *descriptive* reference, not substrate emission. The grep gate (§2.6 step 6) excludes this carve-out by reading the field's name + the surrounding section header.

- ❌ Manifest field `permissions:` carries GitHub Actions permission keys
- ❌ Prompt template inlines `${{ secrets.GITHUB_TOKEN }}`
- ✅ Manifest field `permission_intent:` lists logical permissions; renderer maps to substrate
- ✅ Prompt template names the agent's behavioral commitments; renderer handles secret binding

### 3.4. The package does not ship rendered artifacts

The package's source tree carries the declaration. The renderer emits the substrate artifact at install time at the substrate's canonical path. The two MUST NOT coexist in the package source tree.

- ❌ `cnos.core/orchestrators/agent-admin/cnos-agent-admin.yml` shipped alongside the manifest
- ✅ Renderer emits `.github/workflows/cnos-agent-admin.yml` at `cn wake install agent-admin --agent sigma` time

### 3.5. The renderer rejects malformed declarations rather than defaulting silently

If a required field is missing, the wrong type, or names a path that does not exist (e.g. `prompt_template: "prompt.md"` but the file is absent), `cn wake install` MUST fail with a precise error naming the field and the file. Silent defaulting hides design gaps.

- ❌ Renderer fills in default `responsibilities` when the manifest omits the field
- ✅ Renderer errors `wake-provider.json: required field "responsibilities" missing` and exits non-zero

### 3.6. Cross-references are declared, not inferred

The `cross_references` object names the architecture (cnos#467), predecessor cycles (cnos#468 label doctrine for any wake whose prompt cites label-application discipline), consumed skills (e.g. `cnos.core/skills/agent/activate`), and consumed conventions (e.g. `AGENT-ACTIVATION-LOG-v0.md`). These are the package's authority — they document what the declaration depends on. The renderer does not infer them from the prompt text.

- ❌ Renderer greps the prompt for skill paths and constructs a derived `cross_references` view
- ✅ Manifest declares `cross_references` explicitly; renderer can verify the prompt cites everything declared

### 3.7. Admin-only wakes enforce the boundary in the prompt

A wake declared with `role: admin` + `admin_only: true` MUST carry a prompt template that explicitly forbids cell execution and explicitly names the defer-path for cell-shaped directives. The contract makes this enforceable via the manifest's structural declaration; the prompt's prose makes it operational at agent-runtime.

- ❌ `admin_only: true` declared but the prompt body never names "no cell execution"
- ✅ `admin_only: true` declared and the prompt body contains "MUST NOT execute cells" + the defer-path

### 3.8. Disallowed surfaces explicitly name `.github/workflows/` and cell execution

For *any* wake (admin or dispatch), `.github/workflows/` is a disallowed surface — the substrate is rendered, not hand-edited by the wake. For admin wakes, `cell execution` is named as a disallowed surface (a logical "surface," not a path) so the boundary is grep-able from the manifest without reading the prompt.

- ❌ `disallowed_surfaces: [".github/", "branch_protection"]` (cell execution implicit only in prompt)
- ✅ `disallowed_surfaces: [".github/workflows/", "branch_protection", "repo_settings", "cell_execution", ...]`

---

## 4. Verify

### 4.1. Schema-presence check

Confirm the manifest's first key is `"schema"` and its value is `"cn.wake-provider.v1"`. The renderer dispatches by schema; missing or mistyped schema means the renderer cannot identify the manifest.

### 4.2. Required-field check

Confirm every required field in §2.1 is present and well-typed. Missing required fields are §3.5 failures.

### 4.3. Substrate-agnostic check

Run: `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action' wake-provider.json prompt.md`. The result MUST be 0 *except* for occurrences within an explicit `relationship_to_substrate` field or a "relationship to existing {substrate-bound artifact}" prose section.

### 4.4. Admin-only enforcement check

For wakes declared with `admin_only: true`: grep the prompt body for explicit "no cell execution" language and for the defer-path. Both MUST be present in unambiguous form.

### 4.5. Cross-reference completeness check

For each entry in `cross_references.consumed_skills` and `cross_references.consumed_conventions`, confirm the prompt body names the path. If a skill is declared but not named in the prompt, either the declaration is over-stating or the prompt is under-citing — both are findings.

### 4.6. Substrate-rendered-artifact absence check

Confirm no `.yml` files live under the package's source tree at or under the wake-provider directory. The package authors declarations; the renderer emits substrate artifacts.

---

## 5. Cross-references

This skill is part of the wake-orchestration cluster (cnos#467). It defines the contract consumed by:

- **cnos#467** (`agent/wake-orchestration`) — the master architecture; this contract operationalizes the "package-owned wake providers" doctrine declared there.
- **cnos#468** (`agent/label-doctrine`) — the label control plane any admin wake's prompt cites for label-application discipline. The agent-admin wake provider's prompt template carries this citation.
- **cnos#470** (this cycle; Sub 2 of #467) — authors this skill and the agent-admin wake provider as the reference instance.
- **cnos#450** (`agent/wake-template`, amended) — the Sub 3 renderer (`cn wake install`) consumes this contract. The renderer dispatches by `schema`; required-field validation is its responsibility; substrate emission is its sole authority.
- **cnos#454** (`agent/dispatch-protocol`, amended) — defines the claim mechanics for *dispatch*-role wakes; not consumed by admin-role wakes. Future per-package dispatch wake providers (cnos.cds Sub 4; future cnos.cdr, cnos.cdw) cite both this contract and #454. (cnos.cdd is the generic cell-runtime framework, not a concrete dispatch wake owner.)

Adjacent skills and conventions:

- `cnos.core/skills/agent/activate/SKILL.md` — the activate skill any wake's prompt invokes for identity load.
- `cnos.core/skills/agent/attach/SKILL.md` — the attach skill any wake's prompt invokes for channel sync.
- `cnos.core/skills/agent/label-doctrine/SKILL.md` — the label doctrine admin wakes cite for label-application discipline (cnos#468 merged predecessor).
- `cnos.core/skills/skill/SKILL.md` — the skill-format meta-skill this skill conforms to.
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — the channel log convention any wake's output contract cites.

Reference instance:

- `cnos.core/orchestrators/agent-admin/wake-provider.json` + `prompt.md` — the agent-admin wake provider; the reference instance authored against this contract (cnos#470).

---

## 6. Failure modes

The four named failure modes from §1.3, plus two implementation-level ones:

- **F1 — Substrate leakage.** _Symptom:_ Manifest or prompt carries GitHub Actions YAML / `${{ }}` interpolations / hard-coded workflow paths. _Fix:_ §3.3; replace with logical equivalents the renderer maps.
- **F2 — Renderer overreach.** _Symptom:_ Renderer fills in defaults for required fields the manifest omitted. _Fix:_ §3.5; renderer errors on missing required fields.
- **F3 — Schema drift.** _Symptom:_ A field's meaning changes without a schema-version bump; older renderers misparse. _Fix:_ §3.1; breaking changes bump the schema suffix.
- **F4 — Substrate-rendered artifact in the source tree.** _Symptom:_ Package ships a rendered `.yml` alongside the manifest. _Fix:_ §3.4; renderer emits at install time; source tree carries declarations only.
- **F5 — Admin/dispatch boundary erosion.** _Symptom:_ A wake declared `admin_only: true` carries a prompt that does not forbid cell execution. _Fix:_ §3.7; admin wakes enforce the boundary in the prompt prose.
- **F6 — Cross-reference drift.** _Symptom:_ Manifest declares a `consumed_skill` not named in the prompt, or the prompt cites a skill the manifest does not declare. _Fix:_ §4.5; cross-references are explicit on both sides.

---

## 7. Authority and stability

This skill is doctrine-adjacent: it defines the contract every cnos package follows when shipping a wake provider. Future changes follow the constitutive-change approval discipline that governs other doctrine-adjacent skills (compare `cnos.core/skills/agent/label-doctrine/SKILL.md §8 Versioning`).

Drift between this skill and the Sub 3 renderer (`cn wake install`) is resolved in favor of this skill — the contract is the canonical surface; the renderer follows. If a future renderer needs a field this contract does not declare, the contract is amended (with a schema-version bump if breaking) before the renderer ships the field-dependent behavior.

Drift between this skill and a concrete wake provider declaration (e.g. `cnos.core/orchestrators/agent-admin/wake-provider.json`) is resolved in favor of this skill — the declaration is one *instance* of the contract; the contract governs.
