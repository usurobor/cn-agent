---
name: wake-provider
description: The wake contract a cnos package declares to provide a wake. A wake is a typed `SKILL.md` module (`artifact_class: wake`) whose frontmatter `wake:` block carries the contract (identity, role, input/output, surfaces, defer-path, selector) and whose body is the prompt. Defines what the block MUST contain, the split between what the package declares (the wake skill) and what the substrate renderer materializes (the workflow), and the split between what CUE validates (`#Wake` / I5) and what the renderer enforces at run time. The agent-admin wake (cnos.core, `orchestrators/agent-admin/SKILL.md`) is the reference instance.
artifact_class: skill
kata_surface: none
governing_question: What does a cnos package declare to provide a wake, and which parts of the declaration are the package's authority vs the substrate renderer's authority?
visibility: public
parent: agent
triggers:
  - wake provider
  - wake declaration
  - wake skill
  - cn install-wake
  - package-owned wake
scope: task-local
inputs:
  - none (contract reference)
outputs:
  - the wake contract (the `wake:` frontmatter block + body-as-prompt)
  - the substrate-rendering boundary (package authority vs renderer authority)
  - the typed-validation split (`#Wake` in `schemas/skill.cue` / I5 vs the renderer's runtime refusals)
requires:
  - none (foundational; consumed by `cn install-wake` and by every package-owned wake `SKILL.md`)
---

# Wake contract

A cnos package declares a **wake** to ship a long-form work shape that an external substrate (today: GitHub Actions) materializes into a real workflow file. A wake is a **typed `SKILL.md` module** (`artifact_class: wake`, `scope: global`): its frontmatter `wake:` block is the machine-consumable contract, and its body is the prompt the agent runs. The package authors the *abstract* declaration (identity, role, surfaces, prompt); the renderer authors the *substrate-specific* artifact (today: YAML; tomorrow possibly other schedulers).

This skill is the contract: what a wake `SKILL.md` MUST carry, what it MAY carry, where the `#Wake` CUE definition draws the line, and where the boundary lies between the package's authority and the renderer's authority.

## Core Principle

**A wake is substrate-agnostic data (the `wake:` frontmatter block) + a substrate-agnostic prompt (the body). The substrate is the renderer's authority, not the package's.**

Two coupled invariants govern the contract:

1. **Package authority.** The package declares *what the wake is*: identity, role, allowed/disallowed surfaces, defer-path, the prompt the agent runs. These are fixed per the package's design — changing the substrate does not change them.
2. **Renderer authority.** The renderer materializes the declaration into the current substrate (today: GitHub Actions YAML in `.github/workflows/`). Permission allowlist *shape*, secret-name *binding*, trigger *encoding*, run identity *binding*, concurrency *encoding* — these are substrate-specific and the renderer owns them.

A package that emits substrate-specific YAML in its `wake:` block or its body has confused the boundary; a renderer that decides what responsibilities the wake's prompt enumerates has confused the boundary in the other direction. The two errors compound: package-side substrate emission locks the package to one scheduler; renderer-side responsibility decision means the package's role is defined by the renderer's defaults, not the package's design.

The failure modes this contract structurally prevents are listed in §6.

## Algorithm

A wake is **one typed `SKILL.md` module**: the frontmatter `wake:` block is the contract; the body is the prompt. Both live in a single file the package nominates under a content class. The renderer:

1. **Discovers** the wake by resolving the wake name to `orchestrators/{wake-name}/SKILL.md` under the installed package's content classes (`commands/`, `orchestrators/`, or any future class declared in `cn.package.json`).
2. **Reads** the frontmatter `wake:` block — the typed contract the `#Wake` definition validates at I5 time — and projects it into the internal `cn.wake-provider.v1` shape the renderer queries via `jq`; validates that required fields are present and well-typed.
3. **Inlines** the body (the prompt) into the rendered substrate artifact's `prompt:` field (or the substrate-specific equivalent).
4. **Materializes** the substrate artifact at the substrate's canonical path (today: `.github/workflows/cnos-{wake-name}.yml`); the committed golden (`orchestrators/{wake-name}/cnos-{wake-name}.golden.yml`) is the reviewed copy of that same output.
5. **Reports** render idempotently — a repeated render is a no-op when the materialized artifact already matches (the `install-wake golden` byte-identity oracle: re-render diff clean + `sha256(golden)==sha256(live)` + idempotence).

The contract is what the renderer reads. The package declares; the renderer materializes. This skill defines what the renderer is allowed to expect from the declaration and what the package is allowed to expect from the renderer.

---

## 1. Define

### 1.1. Identify the parts

A wake declaration is **one `SKILL.md` module** (`artifact_class: wake`) with two parts:

- **The `wake:` frontmatter block** — the typed, machine-consumable contract (identity, role, input/output, surfaces, defer-path, selector). Validated by the `#Wake` definition in `schemas/skill.cue` at I5 time. Machine-consumable.
- **The body** — the prompt the renderer inlines into the rendered substrate artifact's prompt field. Substrate-agnostic prose. Verbose reference material (responsibilities narrative, cross-references, doctrine) lives here — after the prompt, appended below a `## …` boundary heading — not in the frontmatter.

A substrate-rendered artifact (the renderer's output) is **not** authored by the package. The declaration is what the package writes by hand; the renderer emits both the live workflow (`.github/workflows/cnos-{wake}.yml`) and its committed golden (`orchestrators/{wake}/cnos-{wake}.golden.yml`) at install time. The live workflow MUST NOT be hand-edited; the golden is the renderer-emitted, reviewed compiled artifact that CI diffs against on every change.

- ❌ Package hand-edits `.github/workflows/cnos-{wake}.yml` (the substrate is rendered, not authored)
- ✅ Package ships `orchestrators/{wake}/SKILL.md` (frontmatter `wake:` + body); the renderer emits `.github/workflows/cnos-{wake}.yml` and the committed `cnos-{wake}.golden.yml` at install time

### 1.2. Articulate how they fit

The `wake:` block is **the schema-bound data** the renderer can dispatch on. The body is **the prompt** the renderer inlines into the substrate's prompt-receiving field. Together, in one typed `SKILL.md`, they are *the wake* — a package-owned, substrate-agnostic declaration validated by the same CUE / I5 pipeline as every other skill.

The package authority / renderer authority split (§Core Principle) makes the split concrete:

| Lives in the declaration (package authority) | Lives in the renderer (substrate authority) |
|---|---|
| Wake name, role, target agent variable | Workflow `name:` field encoding |
| Responsibilities (in the body — what the agent does) | Permission allowlist syntax (e.g. GitHub Actions `permissions:` block) |
| Allowed / disallowed surfaces (`wake.surfaces`) | Secret name binding (e.g. `${{ secrets.X }}`) |
| Input contract (`wake.input`: logical trigger taxonomy) | Substrate trigger encoding (e.g. `on: { schedule: [...], issues: { types: [opened] } }`) |
| Output contract (`wake.output`: channel-log/cursor for admin; cell-artifact for dispatch) | Run-identity binding (e.g. `bot_name` / `bot_id` in claude-code-action) |
| Defer-path for off-role directives (`wake.defer_path`) | Concurrency / serialization encoding |
| The body (prompt) | Inlining the body into the substrate's prompt field |
| Cross-references (in the body: skills the wake invokes; doctrine it cites) | (none — references are package authority) |

The contract is the rule that makes the split machine-checkable: if a field belongs in the right column above, it MUST NOT appear in the `wake:` block; if a field belongs in the left column, it MUST appear in the declaration.

- ❌ `wake:` block carries a `permissions:` field with GitHub Actions YAML keys
- ✅ `wake:` block carries a `permission_intent:` field that names the *logical* permission (e.g. `contents.write`, `issues.write`); the renderer maps to the substrate's encoding

### 1.3. Name the failure mode

The wake contract fails through **substrate leakage**: the package declares substrate-specific content in its `wake:` block or body (GitHub Actions YAML, runner names, secret names with `${{ }}` interpolation, hard-coded workflow paths). When the substrate changes (a future scheduler; a different runner; a different secret manager), every package's declaration must be edited. The package no longer owns its wake — the substrate owns the package.

The opposite failure is **renderer overreach**: the renderer decides what responsibilities the wake enumerates, what defer-path applies, what the prompt body says. The package's intent is reduced to a thin declaration of "just install the wake" and the renderer fills in defaults. Two packages get the same wake regardless of their distinct design.

Both failures are observable:

- **F1 — Substrate leakage.** _Symptom:_ `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action' {wake SKILL.md}` returns non-zero (except in an explicit "relationship to existing substrate-bound wake" section, if one exists). _Fix:_ remove substrate-specific tokens; replace with logical equivalents the renderer maps.
- **F2 — Renderer overreach.** _Symptom:_ The `wake:` block omits one of the required fields below (§2.1) and the renderer inserts a default. _Fix:_ the renderer's contract rejects declarations with missing required fields rather than defaulting silently; the package must declare.
- **F3 — Contract drift.** _Symptom:_ A field changes meaning without the `#Wake` definition changing; older renderers misparse newer declarations. _Fix:_ the typed contract is `#Wake` in `schemas/skill.cue` (the I5 gate); the renderer queries a versioned internal shape (`cn.wake-provider.v1`); breaking changes are made to both.
- **F4 — Hand-edited substrate artifact.** _Symptom:_ The live workflow under `.github/workflows/` is edited by hand and diverges from the golden. _Fix:_ §3.4 — the renderer emits the workflow and its golden; the `install-wake golden` oracle fails on any drift.

---

## 2. Unfold

### 2.1. Required `wake:` block fields

Every wake `SKILL.md` carries the standard skill frontmatter (`name`, `description`, `governing_question`, `artifact_class: wake`, `scope: global`, `kata_surface`, `triggers`, `inputs`, `outputs`) plus a `wake:` block. The `#Wake` definition in `schemas/skill.cue` validates the block's shape at I5; `cn-install-wake` enforces the cross-field refusals at render time. Every `wake:` block MUST declare these fields:

| Field | Type | Purpose |
|---|---|---|
| `role` | string, enum: `admin` \| `dispatch` | The wake's role class. `admin` = channel sync + admin-only directive handling, never executes cells. `dispatch` = claims cells matching its protocol selector and runs the protocol's runtime. Drives the role-shaped `output` disjunction. (There is no `observer` role.) |
| `package` | string | The owning package (e.g. `cnos.core` for the agent-admin wake; `cnos.cds` for the CDS dispatch wake). Authoritative for ownership; the renderer verifies the `SKILL.md` lives under the named package's directory tree. Dispatch wakes belong to concrete protocol packages (cds/cdr/cdw); `cnos.cdd` is the generic cell-runtime framework and owns no dispatch wake. |
| `admin_only` | bool | When `true`, the wake MUST NOT execute cells. The body enforces this; the renderer surfaces it. `role: admin` implies `admin_only: true`; `role: dispatch` implies `admin_only: false`. |
| `activation_log_writer` | bool | Writer-ownership of the activation channel log (`AGENT-ACTIVATION-LOG-v0` §0.1). Admin wakes are the sole writers (`true`); package-dispatch wakes are non-writers (`false`) and MUST declare it explicitly. The renderer refuses a `role: dispatch + admin_only: false` wake that declares (or defaults to) `activation_log_writer: true` with exit 4 (cycle/496). |
| `input` | object | Logical trigger taxonomy. Required key: `triggers` (array; recognized values today: `schedule`, `issues_opened_title_match`, `issues_labeled_selector_match`). Optional key: `issues_opened_title_pattern` (the literal title substring the `issues_opened_title_match` trigger matches; defaults to the wake `name`). The renderer maps to substrate-specific trigger encoding; unknown trigger values are rejected per §3.5. |
| `output` | object, **role-shaped** | Validated by a role disjunction (`#WakeOutputAdmin` \| `#WakeOutputDispatch`). For `role: admin`: `channel_log_convention` (canonical doc path), `writer_surface` (the log path the wake writes), `class_taxonomy` (allowed channel-entry `class:` values), `cursor_advance` (bool), `cursor_field` (how the cursor is recorded). For `role: dispatch`: `cycle_artifact_root` (path under which cell artifacts land, e.g. `.cdd/unreleased/{N}/`), `artifact_class_taxonomy` (canonical artifact class names the wake's cells emit), `cell_runtime` (the cell framework the wake invokes, e.g. `cnos.cdd`). |
| `permission_intent` | array of strings | Logical permissions the wake needs (e.g. `contents.write`, `issues.write`, `pull_requests.write`, `id_token.write`). The renderer maps to substrate-specific permission encoding. |
| `concurrency` | object | Keys: `serialize` (bool; when true the substrate is asked to serialize firings to prevent surface races), `group` (logical concurrency-group name; the renderer maps to substrate syntax). |
| `agent_variable` | object | The per-install agent identity binding. Keys: `name` (the variable the body/substitution uses; today `agent`), `default` (`string \| null`; `null` means operator-required — the admin-wake pattern; a string means operator-defaultable — the dispatch-wake pattern). |

### 2.2. Optional `wake:` block fields

These fields MAY appear; the dispatch-shaped ones are **required when `role: dispatch`**.

| Field | Type | Applies | Purpose |
|---|---|---|---|
| `activation_state` | enum: `live` \| `declaration-only` | dispatch (default `live`) | Names whether the wake is runnable. `declaration-only`: shipped as a contract specimen — the renderer refuses to render (or renders with a never-fire guard) per §3.10. `live`: the renderer-emitted workflow is the production wake. |
| `protocol` | string, kebab-case | **required when `role: dispatch`** | The concrete protocol qualifier the wake claims cells for (e.g. `cds`, `cdr`, `cdw`); matches the `protocol:{id}` label on claimed cells. Owned by the declaring package (cnos#468 §2). |
| `selector` | object `{ include: [...], exclude: [...] }` | **required when `role: dispatch`** | The claim selector — the label-set filter applied when claiming cells (§3.9). |
| `surfaces` | object `{ allowed: [...], disallowed: [...] }` | any | What the wake MAY / MUST NOT write. Patterns use `{agent}` for the per-install agent name. `disallowed` MUST name `.github/workflows/`, and (for admin wakes) `cell_execution` (§3.8). |
| `defer_path` | object | any | Routing doctrine: `cell_shaped_directive`, `off_role_directive`, `ambiguous_directive` — what the wake does when a directive is cell-shaped, off-role, or ambiguous. |

The `#Wake` block is **open** (per `LANGUAGE-SPEC §11`): renderer fields not yet named in `#Wake` pass through validation rather than failing it, so a renderer feature can land ahead of a schema amendment.

### 2.3. Body (prompt) requirements

The body of the wake `SKILL.md` is the prompt the renderer inlines verbatim into the substrate artifact's prompt-receiving field (e.g. claude-code-action's `prompt:`). The body MUST:

1. **Be substrate-agnostic.** No GitHub Actions YAML, no `${{ }}` interpolation, no hard-coded runner names, no `.github/workflows/` paths (except in a documented "relationship to existing substrate-bound wake" section when the wake supersedes one — and even there, the reference is descriptive, not interpolative).
2. **Enumerate the responsibilities.** The body is the agent's view; the `wake:` block is the renderer's view; the two MUST agree.
3. **Enumerate the defer-path** named in `wake.defer_path`. The wake's behavior when a directive falls outside its role is the most important behavioral commitment; it must be explicit where the agent reads it.
4. **For admin wakes (`role: admin`, `admin_only: true`):** explicitly forbid cell execution in unambiguous language (e.g. "MUST NOT execute cells under any circumstance"). This is what makes the admin/dispatch boundary observable in the prompt.
5. **Cite the consumed skills + conventions.** The agent reading the prompt should be able to load every skill the prompt expects without consulting another document.
6. **Substitute the `agent_variable`** (or its default) where the prompt names the agent identity (e.g. "Activate and attach as https://github.com/usurobor/cn-{agent}"). The renderer performs this substitution at install time using the value the operator passed (e.g. `--agent sigma` → `cn-sigma`).

The renderer does NOT rewrite the body. It substitutes named variables (today: `{agent}`); it inlines the result; it does not edit prose. The prompt portion is the body up to the reference-appendix boundary heading; verbose reference material after that heading is not inlined.

### 2.4. Substrate-rendering target

Today's substrate is **GitHub Actions** with `anthropics/claude-code-action@v1` as the agent-execution carrier. Specifically:

- **Substrate artifact path:** `.github/workflows/cnos-{wake-name}.yml` (live) + `orchestrators/{wake-name}/cnos-{wake-name}.golden.yml` (committed golden)
- **Substrate trigger encoding:** YAML `on:` block (e.g. `on: { schedule: [...], issues: { types: [opened] } }`)
- **Substrate permission encoding:** YAML `permissions:` block at workflow level
- **Substrate run-identity binding:** `bot_name` + `bot_id` inputs to `claude-code-action`
- **Substrate secret binding:** `${{ secrets.* }}` interpolations
- **Substrate prompt-receiving field:** `prompt:` input to `claude-code-action`

Tomorrow's substrate may be different. The package's declaration does not change when the substrate changes — only the renderer changes. This is the contract's leverage: future substrate migrations are renderer-local, not package-tree-wide.

The mapping from `wake:` fields to substrate-specific encoding is the renderer's responsibility; this skill names what the declaration MUST provide so the renderer can produce a working substrate artifact without operator-side improvisation.

### 2.5. The package-declares-vs-renderer-materializes split (canonical table)

This split is the contract's core invariant. The table makes it grep-able. The left column is the wake `SKILL.md` (its `wake:` block or body); the right column is renderer / substrate authority.

| Declared in the wake `SKILL.md` (package authority) | Emitted by the renderer (substrate authority) |
|---|---|
| `wake.role`, `wake.package`, `wake.admin_only`, `wake.activation_log_writer` | `name:` field of the substrate artifact |
| `wake.protocol` (dispatch-required) | (n/a — runtime gating; the body uses it to verify cell ownership; the selector encodes it as a label filter) |
| `wake.selector.include` / `.exclude` (label-set filter, dispatch-required) | substrate event filter and/or runtime claim gate (`if:` on `labeled` events; runtime label-presence checks at claim time) |
| Responsibilities (enumerated in the body) | (n/a — responsibilities are prose for the agent) |
| `wake.input.triggers` (logical taxonomy) | `on:` block (substrate trigger encoding) |
| `wake.input.issues_opened_title_pattern` (literal title substring) | substituted into the job-level `if:` for the `issues_opened_title_match` trigger |
| `wake.output.channel_log_convention` / `.writer_surface` / `.class_taxonomy` / `.cursor_*` (admin-shape) | (n/a — the agent reads/writes per the cited convention) |
| `wake.output.cycle_artifact_root` / `.artifact_class_taxonomy` / `.cell_runtime` (dispatch-shape) | (n/a — the agent writes per the cell framework's contract) |
| `wake.surfaces.allowed` / `.disallowed` | (n/a — declaration only; the agent enforces via the body) |
| `wake.defer_path` | (n/a — declaration only; the agent enforces via the body) |
| `wake.permission_intent` (logical permissions) | `permissions:` block (substrate permission encoding) |
| `wake.concurrency` | `concurrency:` block (substrate concurrency encoding) |
| `wake.agent_variable` | (substituted by the renderer at install time using `--agent`) |
| The body (prompt) | inlined into the substrate's prompt field |
| Cross-references (in the body) | (n/a — declaration only; for traceability) |
| (nothing) | substrate-bound secret names, runner names, action versions, `${{ }}` interpolations, default-branch policy, OIDC permission encoding |

The right column is the renderer's territory. A wake `SKILL.md` that declares any of the right column has confused the boundary.

### 2.6. Authoring a new wake

To declare a new wake in package `{pkg}`:

1. **Pick a directory under the package's content classes.** `orchestrators/{wake-name}/` is the standard surface today (per cnos#467 active design constraint — no new `wakes/` class until `commands/` + `orchestrators/` prove insufficient).
2. **Create `SKILL.md`** with the standard skill frontmatter, `artifact_class: wake`, `scope: global`, and a `wake:` block populating every required field from §2.1 (plus the dispatch-shaped fields from §2.2 when `role: dispatch`).
3. **Author the body as the prompt** per §2.3. Move verbose responsibilities, `*_notes`, and cross-references into the body's reference appendix (after the prompt-boundary heading), not the frontmatter — the frontmatter is the typed contract; the body is prompt + doctrine.
4. **Do NOT touch the substrate.** No hand-edits under `.github/workflows/`. The renderer emits the live workflow and its committed golden at `cn install-wake` time.
5. **Verify:** `cue vet` (I5) validates the `wake:` block; the substrate-agnostic check (§4.3) passes; `cn install-wake {wake}` renders a golden byte-identical to the committed one (§4.6).

The agent-admin wake (cnos.core, `orchestrators/agent-admin/SKILL.md`) is the reference instance of this contract. The cnos.cds dispatch wake (`orchestrators/cds-dispatch/SKILL.md`) is the reference `role: dispatch` instance — same shape, declaring `role: dispatch` with `protocol`, `selector`, and a dispatch-shaped `output`. (cnos.cds is the concrete software protocol; cnos.cdd is the generic cell-runtime framework that cds and the other concrete protocol packages — cnos.cdr, cnos.cdw — use under the hood. Dispatch wakes belong to concrete protocols, not to the framework.)

---

## 3. Rules

### 3.1. The typed contract is CUE-validated; the renderer enforces runtime refusals

The `#Wake` definition in `schemas/skill.cue` is the typed contract surface (the I5 gate): it validates the `wake:` block's shape — field types, enums, the role-shaped `output` disjunction, dispatch-required fields. The renderer (`cn install-wake`) owns the *runtime* refusals CUE cannot express statically (activation-state gating, `activation_log_writer` mis-declaration, body resolution, install-time substitution, exit codes) and queries a versioned internal shape (`cn.wake-provider.v1`) synthesized from the frontmatter. A breaking change to the contract is made in both surfaces together; the open `#Wake` block (`…`) lets a renderer field land ahead of its schema amendment.

- ❌ Add a required field to the renderer's expectations without adding it to `#Wake`; I5 cannot catch a declaration that omits it
- ✅ Add the required field to `#Wake` (I5 rejects declarations that omit it) and to the renderer (which enforces any cross-field rule)

### 3.2. The `wake:` block is data; the body is prose

The `wake:` block is machine-consumable structured data. The body is human-readable prose for the agent to read. Do not put prose in the `wake:` block (a free-form `description` in the standard frontmatter is the exception); do not put structured contract data in the body (use named variables only at the `{agent}`-placeholder positions).

- ❌ `wake:` block carries a paragraph of explanation in a `notes` field
- ✅ `wake:` block carries fielded data; prose lives in the body or a `README.md`

### 3.3. Substrate-specific content does not appear in the declaration

The declaration (the `wake:` block and the body) is substrate-agnostic by contract. Substrate-specific content (GitHub Actions YAML, runner names, secret names with `${{ }}`, hard-coded workflow paths, action versions) belongs in the renderer.

The single carve-out: a "relationship to existing substrate-bound wake" prose section MAY name a substrate-bound artifact this wake supersedes — a *descriptive* reference, not substrate emission. The grep gate (§4.3) excludes this carve-out by the surrounding section header.

- ❌ `wake.permission_intent` carries GitHub Actions permission keys
- ❌ The body inlines `${{ secrets.GITHUB_TOKEN }}`
- ✅ `wake.permission_intent` lists logical permissions; the renderer maps to the substrate
- ✅ The body names the agent's behavioral commitments; the renderer handles secret binding

### 3.4. The package does not hand-edit rendered artifacts

The package's source tree carries the declaration (the wake `SKILL.md`). The renderer emits the substrate artifact — the live workflow at `.github/workflows/cnos-{wake}.yml` and its committed golden `orchestrators/{wake}/cnos-{wake}.golden.yml`. The golden is renderer-output committed for review and as the byte-identity oracle; it is regenerated by re-running the renderer, never hand-edited. The live workflow is never hand-edited.

- ❌ Someone edits `.github/workflows/cnos-agent-admin.yml` (or the golden) by hand; it drifts from a fresh render
- ✅ `cn install-wake agent-admin --agent sigma` regenerates both; the `install-wake golden` oracle confirms `sha256(golden)==sha256(live)`

### 3.5. The renderer rejects malformed declarations rather than defaulting silently

If a required `wake:` field is missing, the wrong type, or names a path that does not exist, `cn install-wake` MUST fail with a precise error naming the field. Silent defaulting hides design gaps. (Static shape violations are caught earlier by `#Wake` at I5; the renderer is the second gate for cross-field and path rules.)

- ❌ The renderer fills in a default `output` when the `wake:` block omits it
- ✅ The renderer errors naming the missing required field and exits non-zero

### 3.6. Cross-references are declared in the body, not inferred

The body names the architecture (cnos#467), predecessor cycles, consumed skills (e.g. `cnos.core/skills/agent/activate`), and consumed conventions (e.g. `AGENT-ACTIVATION-LOG-v0.md`) the prompt depends on. These are the package's authority; the renderer does not infer them from the prompt text.

- ❌ The renderer greps the body for skill paths and constructs a derived cross-reference view
- ✅ The body declares its cross-references explicitly; a reviewer can verify the prompt cites everything it depends on

### 3.7. Admin-only wakes enforce the boundary in the body

A wake declared with `role: admin` + `admin_only: true` MUST carry a body that explicitly forbids cell execution and names the defer-path for cell-shaped directives. The `wake:` block makes the boundary structurally declared; the body makes it operational at agent-runtime.

- ❌ `admin_only: true` declared but the body never names "no cell execution"
- ✅ `admin_only: true` declared and the body contains "MUST NOT execute cells" + the defer-path

### 3.8. Disallowed surfaces explicitly name `.github/workflows/` and cell execution

For *any* wake (admin or dispatch), `.github/workflows/` is a disallowed surface — the substrate is rendered, not hand-edited by the wake. For admin wakes, `cell_execution` is named as a disallowed surface (a logical "surface," not a path) so the boundary is grep-able from `wake.surfaces.disallowed` without reading the body.

- ❌ `wake.surfaces.disallowed: [".github/", "branch_protection"]` (cell execution implicit only in the body)
- ✅ `wake.surfaces.disallowed: [".github/workflows/", "cell_execution", "branch protection rules", …]`

### 3.9. Dispatch wakes declare their claim selector explicitly

A wake declared with `role: dispatch` MUST declare both `protocol` (the concrete protocol qualifier; matches `protocol:{id}` labels) and `selector` (the label-set filter applied when claiming cells). Per cnos#454 dispatch-protocol the selector is the contract between the dispatch protocol and the wake — it names exactly which cells the wake claims. A dispatch wake that omits the selector is under-specified; per §3.5 the renderer rejects rather than silently defaulting.

- ❌ `role: dispatch` declared but `selector` omitted; the renderer infers from label doctrine and "just runs"
- ✅ `role: dispatch` with explicit `selector.include: [dispatch:cell, protocol:cds, status:todo]` + `selector.exclude: [status:in-progress, status:blocked, status:review, status:changes]`

The selector's `include` protocol entry MUST be `protocol:{protocol}` where `{protocol}` matches `wake.protocol` — the wake claims only cells whose qualifier matches its package's owned protocol (cross-protocol claims are a label-doctrine violation per cnos#468 §2.1).

### 3.10. Declaration-only wakes are not silently activated

A wake declared with `activation_state: declaration-only` is a contract specimen — it documents the intended wake shape but is not runnable yet (typically because downstream skill dependencies or renderer-extension work has not landed). The renderer MUST NOT silently produce a substrate artifact an operator installs without seeing the activation-state declaration; it MUST either:

- **Refuse to install:** exit non-zero with a precise error naming the `activation_state` value; OR
- **Render with explicit warning:** emit the substrate artifact with a top-of-file `# DO NOT INSTALL: activation_state = declaration-only` comment + a stderr warning, AND set the workflow's job-level `if:` gate to a permanently-false condition so the rendered artifact cannot fire even if committed.

For declaration-only wakes, refusal is the safest default; the warned-render is the carve-out when the operator explicitly wants the rendered artifact for diff review. When the preconditions hold, a future cycle flips `activation_state` to `live`, the renderer succeeds without the guard, and the substrate artifact is committed and activated. The flip is a package-authority decision (the `wake:` block changes), not a renderer decision.

- ❌ The renderer reads `activation_state: declaration-only` and renders a normal-looking workflow (the operator commits and it fires without realizing the contract is incomplete)
- ✅ The renderer refuses (`activation_state == declaration-only; refusing to render`), OR emits with a leading `# DO NOT INSTALL` banner + a never-fire `if:` gate

---

## 4. Verify

### 4.1. Contract-shape check

Confirm the wake `SKILL.md` carries `artifact_class: wake`, `scope: global`, and a `wake:` block, and that `cue vet` against `#Wake` passes (the I5 gate). Missing or mistyped fields are §3.1 / §3.5 failures.

### 4.2. Required-field check

Confirm every required field in §2.1 is present and well-typed, plus the dispatch-shaped fields (§2.2) when `role: dispatch`. Missing required fields are §3.5 failures.

### 4.3. Substrate-agnostic check

Run: `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action' {wake SKILL.md}`. The result MUST be 0 *except* for occurrences within an explicit "relationship to existing substrate-bound wake" prose section.

### 4.4. Admin-only enforcement check

For wakes declared with `admin_only: true`: grep the body for explicit "no cell execution" language and for the defer-path. Both MUST be present in unambiguous form.

### 4.5. Cross-reference completeness check

For each skill / convention the body cites as consumed, confirm the body names the path. If a skill is depended on but not named, the body is under-citing; both directions are findings.

### 4.6. Golden byte-identity check

Confirm `cn install-wake {wake}` renders a live workflow byte-identical to the committed golden (`sha256(golden)==sha256(live)`), and that re-rendering is a no-op. The `install-wake golden` CI job is the oracle. The golden is the reviewed compiled artifact; a diff means either the declaration changed (re-render + commit) or the substrate artifact was hand-edited (§3.4 violation).

---

## 5. Cross-references

This skill is part of the wake-orchestration cluster (cnos#467). It defines the contract consumed by:

- **cnos#467** (`agent/wake-orchestration`) — the master architecture; this contract operationalizes the "package-owned wakes" doctrine declared there.
- **cnos#468** (`agent/label-doctrine`) — the label control plane any admin wake's body cites for label-application discipline.
- **cnos#470** — authored this skill and the agent-admin wake as the reference instance.
- **cnos#450** (`agent/wake-template`, amended) — the renderer (`cn install-wake`) consumes this contract; required-field validation and substrate emission are its authority.
- **cnos#454** (`agent/dispatch-protocol`, amended) — defines the claim mechanics for *dispatch*-role wakes; cited by dispatch wakes alongside this contract.
- **cnos#524** (`wake-as-skill`) — migrated wakes to the typed `SKILL.md` model this skill describes; added `#Wake` + `artifact_class: wake` to `schemas/skill.cue` and made `cn-install-wake` read the wake `SKILL.md` as its only source.

Adjacent skills and conventions:

- `cnos.core/skills/agent/activate/SKILL.md` — the activate skill any wake's body invokes for identity load.
- `cnos.core/skills/agent/attach/SKILL.md` — the attach skill any wake's body invokes for channel sync.
- `cnos.core/skills/agent/label-doctrine/SKILL.md` — the label doctrine admin wakes cite for label-application discipline (cnos#468).
- `cnos.core/skills/skill/SKILL.md` — the skill-format meta-skill this skill (and every wake `SKILL.md`) conforms to.
- `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` — the channel log convention an admin wake's `output` cites.

Reference instances:

- `cnos.core/orchestrators/agent-admin/SKILL.md` — the agent-admin wake; the reference `role: admin` instance authored against this contract (cnos#470).
- `cnos.cds/orchestrators/cds-dispatch/SKILL.md` — the CDS dispatch wake; the reference `role: dispatch` instance (cnos#467 Sub 4).

---

## 6. Failure modes

The four named failure modes from §1.3, plus two implementation-level ones:

- **F1 — Substrate leakage.** _Symptom:_ The `wake:` block or body carries GitHub Actions YAML / `${{ }}` interpolations / hard-coded workflow paths. _Fix:_ §3.3; replace with logical equivalents the renderer maps.
- **F2 — Renderer overreach.** _Symptom:_ The renderer fills in defaults for required fields the declaration omitted. _Fix:_ §3.5; the renderer errors on missing required fields.
- **F3 — Contract drift.** _Symptom:_ A field's meaning changes without `#Wake` (and the renderer's internal shape) changing; declarations misparse. _Fix:_ §3.1; the change is made to both the CUE contract and the renderer.
- **F4 — Hand-edited substrate artifact.** _Symptom:_ The live workflow or golden diverges from a fresh render. _Fix:_ §3.4; the renderer emits both; the byte-identity oracle fails on drift.
- **F5 — Admin/dispatch boundary erosion.** _Symptom:_ A wake declared `admin_only: true` carries a body that does not forbid cell execution. _Fix:_ §3.7; admin wakes enforce the boundary in the body.
- **F6 — Cross-reference drift.** _Symptom:_ The body depends on a skill it does not name, or names one it does not use. _Fix:_ §4.5; cross-references are explicit in the body.

---

## 7. Authority and stability

This skill is doctrine-adjacent: it defines the contract every cnos package follows when shipping a wake. Future changes follow the constitutive-change approval discipline that governs other doctrine-adjacent skills (compare `cnos.core/skills/agent/label-doctrine/SKILL.md §8 Versioning`).

Drift between this skill and the renderer (`cn install-wake`) is resolved in favor of this skill — the contract is the canonical surface; the renderer follows. Drift between this skill and the typed `#Wake` definition is resolved by amending both together (this skill states the contract in prose; `#Wake` states it in CUE; they must agree). Drift between this skill and a concrete wake `SKILL.md` (e.g. `agent-admin/SKILL.md`) is resolved in favor of this skill — the declaration is one *instance* of the contract; the contract governs.
