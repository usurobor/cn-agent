---
cycle: 524
parent_issue: cnos#524
document_class: w0-design
status: locked
authored_by: α@cdd.cnos (R0)
date: 2026-06-30 (UTC)
scope: W0 design only (operator directive); W1→W4 ACs not in scope for this document's cycle
---

# W0 design — cnos#524: wake-as-skill

This document formalizes the converged, operator-ratified W0 design for the wake-as-skill migration. It is a structured extraction of the design already present in issue #524 body (`## W0 design (locked)`) and the operator W0-calls comment (first @usurobor comment). No new design is introduced here.

---

## §A. Problem statement and invariant

**Problem:** A wake is currently two artifacts in a bespoke JSON+prompt contract system — `wake-provider.json` (the typed manifest) plus `prompt.md` (the prompt body) — compiled by `cn install-wake` into `.github/workflows/cnos-<wake>.yml`. This system runs **beside** the typed SKILL.md contract system (`SKILL.md` frontmatter + body, validated by `schemas/skill.cue` via the I5 gate). Two contract systems means: wake contracts escape CUE typing, every wake change touches two files in a non-standard format, and the discipline imposed by the skill pipeline does not apply to wakes.

**What is expected:** one contract system. A wake is a typed `SKILL.md` module whose frontmatter carries a `wake:` contract block and whose body is the prompt, validated by CUE and compiled to a workflow by the renderer.

**Invariant:** *A wake is a typed SKILL.md module whose body is the prompt and whose frontmatter compiles to workflow substrate. One SKILL.md+CUE contract system, not a second JSON+prompt system.*

---

## §B. The `wake:` block schema (the `#Wake` CUE definition)

### B.1 Complete `wake:` block shape

```yaml
wake:
  role: admin | dispatch                           # required; drives role-shaped output disjunction
  package: string                                  # required; e.g. cnos.core, cnos.cds
  admin_only: bool                                 # required; true for admin wake, false for dispatch
  activation_log_writer: bool                      # required; whether this wake writes channel logs
  activation_state?: live | declaration-only       # dispatch only; optional on admin
  protocol?: string                                # dispatch only (e.g. "cds"); optional on admin
  selector?:                                       # dispatch only; optional on admin
    include: [...string]
    exclude: [...string]
  input:
    triggers: [...string]                          # required; e.g. ["schedule", "issues_opened_title_match"]
    issues_opened_title_pattern?: string           # optional; present on admin wake
  output: <role-shaped disjunction; see §B.2>
  permission_intent: [...string]                   # required; logical permissions
  concurrency:
    serialize: bool                                # required
    group: string                                  # required; e.g. "agent-admin-{agent}"
  agent_variable:
    name: string                                   # required
    default: string | null                         # required; null means operator-required
  surfaces?:
    allowed: [...string]
    disallowed: [...string]
  defer_path?:
    cell_shaped_directive?: string
    off_role_directive?: string
    ambiguous_directive?: string
```

### B.2 Role-shaped output disjunction

The `wake.output` field is role-shaped. The CUE `#Wake` definition MUST express a disjunction based on `wake.role`:

**When `wake.role == "admin"` (agent-admin shape):**
```yaml
output:
  channel_log_convention: string                   # e.g. docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
  writer_surface: string                           # e.g. ".cn-{agent}/logs/YYYYMMDD.md"
  class_taxonomy: [...string]                      # entry class enum values
  cursor_advance: bool
  cursor_field: string
```

**When `wake.role == "dispatch"` (cds-dispatch shape):**
```yaml
output:
  cycle_artifact_root: string                      # e.g. ".cdd/unreleased/{N}/"
  artifact_class_taxonomy: [...string]             # CDD artifact class enum values
  cell_runtime: string                             # e.g. "cnos.cdd"
```

### B.3 Field-by-field cross-reference to current `wake-provider.json` manifests

The following table maps every field in the two current `wake-provider.json` files to its placement in the SKILL.md contract. Fields marked "→ body" are verbose prose that moves to the SKILL.md body per design decision §C-5; they are NOT frontmatter.

| `wake-provider.json` field | Wake | Placement in SKILL.md |
|---|---|---|
| `schema` | both | Not carried; artifact class signals schema version |
| `name` | both | SKILL.md `name:` field (standard skill field) |
| `package` | both | `wake.package` |
| `role` | both | `wake.role` |
| `admin_only` | both | `wake.admin_only` |
| `activation_log_writer` | both | `wake.activation_log_writer` |
| `activation_state` | dispatch | `wake.activation_state` |
| `activation_state_notes` | dispatch | → body |
| `protocol` | dispatch | `wake.protocol` |
| `selector.include` | dispatch | `wake.selector.include` |
| `selector.exclude` | dispatch | `wake.selector.exclude` |
| `description` | both | SKILL.md `description:` field (standard skill field) |
| `responsibilities` | both | → body |
| `input_contract.triggers` | both | `wake.input.triggers` |
| `input_contract.issues_opened_title_pattern` | admin | `wake.input.issues_opened_title_pattern` |
| `input_contract.trigger_descriptions` | both | → body |
| `input_contract.inbound` | both | → body |
| `output_contract.channel_log_convention` | admin | `wake.output.channel_log_convention` |
| `output_contract.writer_surface` | admin | `wake.output.writer_surface` |
| `output_contract.class_taxonomy` | admin | `wake.output.class_taxonomy` |
| `output_contract.class_taxonomy_notes` | admin | → body |
| `output_contract.cursor_advance` | admin | `wake.output.cursor_advance` |
| `output_contract.cursor_field` | admin | `wake.output.cursor_field` |
| `output_contract.cycle_artifact_root` | dispatch | `wake.output.cycle_artifact_root` |
| `output_contract.artifact_class_taxonomy` | dispatch | `wake.output.artifact_class_taxonomy` |
| `output_contract.artifact_class_notes` | dispatch | → body |
| `output_contract.cell_runtime` | dispatch | `wake.output.cell_runtime` |
| `output_contract.cell_runtime_notes` | dispatch | → body |
| `allowed_surfaces` | both | `wake.surfaces.allowed` |
| `disallowed_surfaces` | both | `wake.surfaces.disallowed` |
| `defer_path.cell_shaped_directive` | both | `wake.defer_path.cell_shaped_directive` |
| `defer_path.off_role_directive` | both | `wake.defer_path.off_role_directive` |
| `defer_path.ambiguous_directive` | both | `wake.defer_path.ambiguous_directive` |
| `prompt_template` | both | Not carried; body IS the prompt (§E) |
| `agent_variable.name` | both | `wake.agent_variable.name` |
| `agent_variable.default` | both | `wake.agent_variable.default` |
| `agent_variable.description` | both | → body |
| `permission_intent` | both | `wake.permission_intent` |
| `permission_intent_notes` | both | → body |
| `concurrency_intent.serialize` | both | `wake.concurrency.serialize` |
| `concurrency_intent.group` | both | `wake.concurrency.group` |
| `concurrency_intent.notes` | both | → body |
| `superseded_substrate_artifact` | admin | → body |
| `relationship_to_substrate` | admin | → body |
| `cross_references` | both | → body |

**Notes on field renaming:**
- `input_contract` → `wake.input` (shorter, consistent with YAML block-key convention)
- `output_contract` → `wake.output`
- `concurrency_intent` → `wake.concurrency`
- `allowed_surfaces` / `disallowed_surfaces` → `wake.surfaces.allowed` / `wake.surfaces.disallowed`
- `defer_path.*` → `wake.defer_path.*` (unchanged keys, nested under `wake`)

### B.4 Standard SKILL.md fields carried by the wake

A wake SKILL.md uses the standard `#Skill` fields for the following (they do NOT move under `wake:`):
- `name:` — the wake name (e.g. `agent-admin`, `cds-dispatch`)
- `description:` — the short description (from `wake-provider.json` `description`)
- `governing_question:` — authored per skill discipline
- `triggers:` — install triggers (distinct from `wake.input.triggers` which are substrate event triggers; these are skill-level invocation triggers)
- `scope: global` — per design decision §C-3
- `artifact_class: wake` — per design decision §C-2
- `kata_surface: none`
- `inputs:` — logical inputs (human-readable)
- `outputs:` — logical outputs (human-readable)

---

## §C. The 5 operator design decisions

### Decision 1: Use `wake:` nesting (single block, not flattened)

**Decision:** Wake-specific contract fields live under a single `wake:` block in frontmatter. Do not flatten into top-level keys like `wake_role:`, `wake_selector:`, etc.

**Rationale:** The SKILL.md frontmatter stays readable and the wake contract has a clear boundary. Flattening would pollute the standard skill namespace with wake-specific keys, making it impossible for validators to distinguish wake fields from skill fields without special-casing.

**W1 implication:** When authoring both wake SKILL.md files, all wake-specific fields go under `wake:`. The renderer reads from `wake.*` keys. The I5 validator type-checks the `wake:` block as the `#Wake` CUE definition.

### Decision 2: Use `artifact_class: wake`

**Decision:** A wake's `artifact_class` is `wake`, not `skill`, `runbook`, or `reference`. The SKILL.md discovery mechanism (which walks the repo for all `SKILL.md` files) discovers wakes automatically; the `artifact_class` tells tooling what kind of module it is.

**Rationale:** Cleaner than pretending a wake is an ordinary `skill`. Cleaner than inventing a separate discovery pipeline. The I5 validator already discovers SKILL.md files; adding `wake` to the `artifact_class` enum reuses that pipeline without modification. The I5 count advances from 91 to 93 when the two wake SKILL.md files are validated.

**W1 implication:** AC1 adds `"wake"` to the `artifact_class` enum in `schemas/skill.cue`. AC2 sets `artifact_class: wake` on both wake SKILL.md files. I5 must show 93 validated SKILL.md files (no findings) after W2.

### Decision 3: Use existing `scope: global`; role is carried by `wake.role`

**Decision:** A wake uses the existing `scope: global` vocabulary. The wake role (`admin`, `dispatch`) is carried by `wake.role`, not by a new `scope:` value.

**Rationale:** A wake is repo/workflow-level, which maps exactly to `scope: global`. Adding new scope vocabulary (e.g. `scope: wake`) for a concept already captured by `artifact_class: wake` + `wake.role` would inflate the scope enum with no discriminating value.

**W1 implication:** Both wake SKILL.md files set `scope: global`. No new scope vocabulary is added to `schemas/skill.cue`.

### Decision 4: CUE owns typed contract; renderer owns substrate compilation + runtime refusals

**Decision:** There is a strict boundary between what CUE validates at I5 time and what the renderer enforces at install/run time.

**CUE (`schemas/skill.cue`) validates:**
- `wake:` block shape and field types
- Enums: `wake.role`, `wake.activation_state`, `wake.permission_intent` values
- Role-shaped required fields (admin-output vs dispatch-output disjunction)
- `wake.selector` shape (`include`/`exclude` arrays)
- `wake.permission_intent` array
- `wake.concurrency` shape

**Renderer (`cn-install-wake`) enforces:**
- `activation_state` refusal (exit 3 when `declaration-only`)
- `activation_log_writer` mis-declaration refusal (exit 4 when `role: dispatch + admin_only: false` declares `activation_log_writer: true`)
- The FN-6 `attach`-incompatibility refusal
- Body extraction (second `---` delimiter; the body IS the prompt — see §E)
- Install-time substitution (`{agent}` → operator-supplied value)
- GitHub Actions YAML encoding (substrate-specific `permissions:`, `concurrency:`, `on:` triggers, cron slots)
- Bot identity resolution (mapping `agent_variable.default` to a substrate bot identity)

**Rationale:** CUE is a static shape validator; it cannot enforce cross-field runtime conditions (e.g. "if this wake is installed on substrate X, check activation_state"). The renderer is already the single enforcement surface for substrate concerns; keeping runtime refusals in the renderer avoids duplicating that logic in CUE and preserves the "package owns role behavior, renderer owns substrate encoding" split already documented in `cn-install-wake`.

**W1 implication:** AC1 adds ONLY static shape/type/enum validation to `schemas/skill.cue`. AC6 verifies that the renderer's existing refusal smokes still fire from `SKILL.md`-sourced input after the renderer flip in W3.

### Decision 5: `responsibilities` / `*_notes` / `cross_references` → body, not frontmatter

**Decision:** Verbose prose fields from `wake-provider.json` — `responsibilities`, `*_notes` (e.g. `activation_state_notes`, `class_taxonomy_notes`, `cell_runtime_notes`, `permission_intent_notes`, `concurrency_intent.notes`), `cross_references`, `trigger_descriptions`, `relationship_to_substrate`, `superseded_substrate_artifact`, `agent_variable.description` — move to the SKILL.md body, not to frontmatter.

**Rationale:** Frontmatter is contract (machine-readable fields consumed mechanically by the renderer or validator). Body is prompt/role doctrine (consumed by the agent executing the wake). Verbose prose that is neither type-checked by CUE nor consumed by the renderer belongs in the body where it informs the agent's execution.

**W1 implication:** AC2 requires that the SKILL.md body for each wake equals `prompt.md` verbatim (the existing prose prompt) PLUS the verbose reference data that is moving from `wake-provider.json` into the body. The byte-identity oracle (§F) proves that the renderer extracting the body produces the same prompt the renderer previously read from `prompt.md`.

---

## §D. CUE ↔ renderer responsibility boundary

| Concern | CUE owns | Renderer owns |
|---|---|---|
| `wake:` block field types and shapes | yes | no |
| Enum validation (`wake.role`, `wake.activation_state`, `wake.permission_intent[*]`) | yes | no |
| Role-shaped output disjunction (admin-output vs dispatch-output fields required) | yes | no |
| `wake.selector` shape (`include`/`exclude` arrays) | yes | no |
| `wake.concurrency` shape | yes | no |
| `wake.permission_intent` array shape | yes | no |
| Static presence of required fields | yes | no |
| Cross-field runtime refusal: `activation_state == "declaration-only"` → exit 3 | no | yes |
| Cross-field runtime refusal: `role:dispatch + admin_only:false + activation_log_writer:true` → exit 4 | no | yes |
| FN-6 attach-incompatibility refusal | no | yes |
| Body extraction from SKILL.md (second `---` delimiter rule) | no | yes |
| Install-time substitution (`{agent}` → operator-supplied value) | no | yes |
| GitHub Actions YAML encoding (`permissions:`, `concurrency:`, `on:` triggers, cron slots) | no | yes |
| Bot identity resolution (`agent_variable.default` → substrate bot name/id) | no | yes |
| Prompt text delivered to the agent (body of SKILL.md) | no | yes |
| Existence of `wake-provider.json` or `prompt.md` (pre-W4) | no | yes |

**Design note:** CUE asserts static shape; the renderer is the single enforcement surface for all substrate, runtime, and defense-in-depth concerns. This is the "package owns role behavior, renderer owns substrate encoding" split documented in `cn-install-wake`.

---

## §E. Body-as-prompt rule

**Rule:** The SKILL.md body (everything after the second `---` frontmatter delimiter) IS the prompt. The renderer extracts the body verbatim and uses it as the prompt delivered to the agent executing the wake.

**Implication:** `prompt.md` is the current source of truth for the prompt. When authoring the wake SKILL.md (W1), the body MUST equal the current `prompt.md` content verbatim — no paraphrase, no reformatting, no omissions. The body may additionally include the verbose reference data moving from `wake-provider.json` (responsibilities, cross_references, notes), appended after the verbatim `prompt.md` content.

**Frontmatter / body split:**
- Frontmatter (`---` … `---`): machine-readable contract fields only. Consumed by the renderer and the CUE validator. No prose.
- Body: the agent's execution prompt + role doctrine + verbose reference data. Not consumed by the renderer structurally (extracted as a single text block). Not CUE-typed.

**Byte-identity oracle:** The oracle that proves correctness at each W1→W4 step is:
```
render(SKILL.md source) == render(wake-provider.json + prompt.md source)
```
Equivalently: the renderer extracting the SKILL.md body produces output byte-identical to the renderer reading `prompt.md` directly. The `install-wake golden` CI check (sha256 golden == live, diff clean, idempotence) is the substrate encoding of this oracle.

**Why this matters:** The body-as-prompt rule is what makes the migration a no-op for the rendered artifact. If the body is not verbatim `prompt.md`, the rendered workflow's prompt section will differ, breaking byte-identity. The rule is enforced by the oracle, not by a CUE constraint.

---

## §F. W1→W4 migration plan

Each phase keeps both `install-wake golden` (byte-identical CI) and I5 (91→93 typed modules) green. The byte-identity oracle is checked at every phase boundary.

### W1 — Author wake SKILL.md files; extend `schemas/skill.cue`

**What:** Author both `agent-admin/SKILL.md` and `cds-dispatch/SKILL.md` beside their existing `wake-provider.json` and `prompt.md`. Add `#Wake` definition and `"wake"` to `artifact_class` in `schemas/skill.cue`.

**Oracle:** I5 validates both new SKILL.md files (count: 91 → 93, no findings). `cue vet` passes on conformant wake SKILL.md; fails on a wake SKILL.md with a bad enum or missing role-required field. The renderer still reads from `wake-provider.json` + `prompt.md` — goldens are unchanged.

**State after W1:** Three sources of truth coexist briefly: `wake-provider.json` (renderer reads this), `prompt.md` (renderer reads this), and the new `SKILL.md` (I5 validates this). No rendered output has changed.

### W2 — Teach the renderer to read SKILL.md; dual-source parity gate

**What:** Extend `cn-install-wake` to resolve its source from SKILL.md (frontmatter `wake:` block → the fields it currently reads from `wake-provider.json`; body → the prompt it currently reads from `prompt.md`). Run a dual-source parity gate: `render(SKILL.md)` must be byte-identical to `render(wake-provider.json + prompt.md)`.

**Oracle:** `cmp render(SKILL) render(JSON)` passes for both wakes. The `install-wake golden` CI check remains green (renderer can still fall back to JSON, or the parity gate is run before the CI golden check).

**State after W2:** The renderer can read from either source. Parity is proved. Goldens are still produced from whichever source is active (JSON by default, unless the renderer is already flipped).

### W3 — Flip renderer source to SKILL.md; re-render

**What:** Configure `cn-install-wake` to use SKILL.md as the primary source (not `wake-provider.json`). Re-render both wakes.

**Oracle:** `install-wake golden` CI check green (rendered output byte-identical to committed goldens). The re-render is expected to be a no-op if W2 parity holds. If the re-render produces any diff, the diff is a bug — iterate to fix before proceeding.

**State after W3:** Renderer reads only SKILL.md. `wake-provider.json` and `prompt.md` are still present but no longer consumed.

### W4 — Delete `wake-provider.json` + `prompt.md`; verify I5 count

**What:** After the W3 byte-identical proof holds in CI, delete `wake-provider.json` and `prompt.md` for both wakes. Verify the renderer still renders byte-identical goldens with no `wake-provider.json` or `prompt.md` present.

**Oracle:** Files absent; `cn install-wake agent-admin` and `cn install-wake cds-dispatch` render byte-identical goldens from SKILL.md only. I5 count = 93 (91 pre-migration + 2 wake SKILL.md files). The migration is complete.

**State after W4:** One contract system. Two typed `SKILL.md` modules in the skill pipeline. No `wake-provider.json` or `prompt.md`.

---

## §G. W1 acceptance criteria (AC1–AC7)

These are the ACs from issue #524 body, restated as the acceptance surface for the W1→W4 build sequence. This cycle (R0) is W0 design only; these ACs govern the next cycle's implementation.

### AC1 — `#Wake` schema in CUE
**Invariant:** `schemas/skill.cue` adds `"wake"` to the `artifact_class` enum and a `#Wake` definition validating the `wake:` block. The definition covers: `role` enum (`admin | dispatch`); `admin_only` / `activation_log_writer` bool; `activation_state` optional enum (`live | declaration-only`); `protocol` optional string; `selector` optional with `include`/`exclude` arrays; `permission_intent` array; `concurrency` shape with `serialize` bool and `group` string; `input.triggers` array; `input.issues_opened_title_pattern` optional string; role-shaped `output` disjunction (admin shape vs dispatch shape).
**Oracle:** `cue vet` a conformant wake SKILL.md → pass; a wake SKILL.md with a bad enum or missing role-required field → fail.

### AC2 — Wake SKILL.md modules authored
**Invariant:** `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` and `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` exist. Each has: `scope: global`, `artifact_class: wake`, a `wake:` block carrying the exact current manifest data (per §B field mapping), and a body equal to the current `prompt.md` verbatim. Verbose `responsibilities` / `*_notes` / `cross_references` move into the body; frontmatter is contract only.
**Oracle:** I5 validator (`scripts/ci/validate-skill-frontmatter.sh`) validates both → "93 SKILL.md validated; no findings".

### AC3 — Renderer reads SKILL.md
**Invariant:** `cn-install-wake` resolves its source from SKILL.md (`wake:` frontmatter block → fields; body → prompt), preserving install-time `{agent}` / `cn-{agent}` substitution.
**Oracle:** `cn install-wake <name>` renders from SKILL.md.

### AC4 — Byte-identical goldens (the binding proof)
**Invariant:** Rendering from SKILL.md produces output byte-identical to the committed goldens. During transition (W2), `render(SKILL)` byte-identical to `render(JSON)`.
**Oracle:** `install-wake golden` CI green (re-render diff clean + sha256 golden == live + idempotence) for both wakes; one-shot dual-source parity check (`cmp render(SKILL) render(JSON)`) passes before deletion.

### AC5 — JSON + prompt deleted
**Invariant:** After AC4 holds, `wake-provider.json` and `prompt.md` are removed for both wakes. The renderer reads only SKILL.md.
**Oracle:** Files absent; `cn install-wake` still renders byte-identical goldens with no `wake-provider.json` / `prompt.md` present.

### AC6 — Renderer refusals preserved (defense-in-depth)
**Invariant:** The renderer's runtime refusals still fire from SKILL.md source — `activation_log_writer` mis-declaration (exit 4), the FN-6 `attach`-incompatibility refusal, and `activation_state` gating (exit 3 on `declaration-only`).
**Oracle:** Existing cycle/496 mis-declaration smoke in `install-wake-golden.yml` still passes against a wake SKILL.md; a synthetic mis-declared wake SKILL.md → exit 4.

### AC7 — All CI gates green
**Invariant:** No inherited-cap, no new red.
**Oracle:** I1 / I2 / I4 / I5 / I6 / `install-wake golden` / `dispatch-repair-preflight (cnos#516)` / Go / Package / Binary all green on the cell's PR.

---

## §H. Friction notes for the W1 implementer

These friction notes are carried from the prior γ scaffolds (original R0 scaffold at commit `f148a1fd`, updated W0-scope scaffold at commit `9a10781d`). They are material for the W1 implementer.

**FN-1: Renderer role-enum collision.** The renderer (`cn-install-wake`) currently accepts `admin | dispatch | observer` in its internal role check (line 342 area). The W0 `#Wake` schema defines `admin | dispatch` only. The `observer` role, if present in the renderer, is a renderer-internal extension and NOT a valid `wake.role` value in the CUE schema. The W1 implementer must confirm whether `observer` appears in any live manifest and whether it needs to be added to the `#Wake` role enum or kept renderer-internal.

**FN-2: `cue vet` invocation against SKILL.md with multi-doc frontmatter.** The I5 validator extracts frontmatter before running `cue vet`. The W1 implementer must confirm the extractor handles the SKILL.md body (second `---` delimiter) correctly — specifically that the body is stripped before `cue vet` runs, so `cue vet` sees only the frontmatter YAML/JSON (not the full SKILL.md file).

**FN-3: `agent_variable.default: null` encoding.** The admin wake has `agent_variable.default: null` (no default; operator must supply). YAML `null` must round-trip correctly through the frontmatter extractor and through `cue vet`. The CUE schema must represent `null` as `null` (not `""` or absent). The W1 implementer must test the `null` case explicitly.

**FN-4: `surfaces.allowed` / `surfaces.disallowed` array length.** Both wakes have long `allowed_surfaces` and `disallowed_surfaces` arrays (8–9 entries each). The CUE schema should accept open arrays (`[...string]`); no length constraint. The W1 implementer should confirm the extractor handles multi-line YAML arrays of this length without truncation.

**FN-5: Body verbatim constraint.** AC2 requires the SKILL.md body to equal `prompt.md` verbatim. The W1 implementer must not paraphrase, reformat, or summarize. Any trailing newline difference in `prompt.md` must be preserved exactly, or the byte-identity oracle will fail at W2.

**FN-6: FN-6 attach-incompatibility refusal (defense-in-depth).** The renderer has a refusal that fires when a wake declares `activation_log_writer: false` but its role is not `dispatch` (or some analogous incompatibility with the `attach` surface). The exact trigger condition for this refusal must be confirmed from `cn-install-wake` source before W3. AC6 requires this refusal to still fire from SKILL.md source; the W1 implementer must include a synthetic SKILL.md fixture for this smoke.

**FN-7: Dual-source parity gate implementation.** W2 requires a parity gate: `render(SKILL) cmp render(JSON)`. The W1 implementer must design how this gate is implemented — either as a new `cn install-wake --source=skill` flag, a temporary parallel render path, or a CI-only comparison step. The gate must not break the existing `install-wake golden` CI check.

**FN-8: Deletion sequencing.** W4 deletes `wake-provider.json` and `prompt.md`. The deletion must happen in a single commit after the byte-identical proof holds in CI — not spread across multiple commits where the renderer might be caught in a state where neither source is complete. If the renderer's source-resolution logic is `SKILL.md then JSON fallback`, the deletion is safe in a single commit. If the fallback path is removed in a separate commit, the W1 implementer must sequence the commits carefully.

---

## §I. Non-goals

The following are explicitly out of scope for the W1→W4 migration:

- **No rendered-workflow output change.** The `.github/workflows/cnos-agent-admin.yml` and `cnos-cds-dispatch.yml` files must be byte-identical before and after the migration. No functional change to the rendered substrate.
- **No wake runtime-behavior change.** The agent executing either wake reads the same prompt, applies the same logic, and produces the same outputs after the migration as before.
- **No renderer substrate-encoding change.** The renderer's substrate-encoding logic (GitHub Actions YAML structure, cron slot encoding, permissions block, concurrency block) is unchanged. Only the source the renderer reads changes (from JSON+prompt to SKILL.md).
- **No new wake names.** Only `agent-admin` and `cds-dispatch` are in scope. No new wake roles beyond `admin` and `dispatch` are introduced.
- **No Demo 0.** This migration is pre-Demo-0 infrastructure; it does not implement Demo 0 behavior.
- **No `.cdd/releases/` work.** This is a pre-release infrastructure cycle; no release artifacts.
- **No per-package bot identities.** The `cnos#449` follow-up (per-package bot accounts) is deferred. The migration uses the existing `sigma` substrate identity.
- **No change to the CUE validator script.** `scripts/ci/validate-skill-frontmatter.sh` logic is unchanged; only `schemas/skill.cue` is extended.

---

_Authored by α@cdd.cnos (R0), 2026-06-30 (UTC). Formalizes the W0 design from cnos#524 body (`## W0 design (locked)`) and the operator W0-calls comment. No new design introduced._
