---
cycle: 524
parent_issue: cnos#524
protocol: cds
cycle_branch: cycle/524
base_main_sha: 23240e4d
head_sha_at_scaffold: 23240e4d
mode: MCA
role: γ
authored_by: γ@cdd.cnos (wake-invoked δ dispatch)
date: 2026-06-30 (UTC)
output_contract: γ-scaffold + α prompt + β prompt
dispatch_scope: W1-implementation
ac_set: AC1, AC2, AC4, AC7 (W1-relevant; AC3/AC5/AC6 are W2/W3/W4 scope)
w0_locked: true
w0_design_ref: .cdd/unreleased/524/w0-design.md
---

# γ-scaffold — cnos#524: wake-as-skill W1 implementation

**Scope:** W1 — first implementation phase. Author both wake `SKILL.md` modules, extend `schemas/skill.cue` with `#Wake`. The renderer still reads from `wake-provider.json` + `prompt.md`; goldens are unchanged.

**W0 reference:** `.cdd/unreleased/524/w0-design.md` (locked; governs the design). This scaffold interprets the W0 design into W1 implementation instructions. Do not re-derive design decisions from first principles — read W0 first.

---

## §1. W1 scope

**What W1 delivers (the only three writable changes):**
1. `schemas/skill.cue` — additive: `"wake"` added to `artifact_class` enum + `#Wake` CUE definition
2. `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` — new file; `artifact_class: wake`
3. `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — new file; `artifact_class: wake`

**What W1 does NOT deliver** (these are W2/W3/W4 scope, explicitly out of W1):
- Renderer flip (W2/W3)
- Parity gate `render(SKILL) cmp render(JSON)` (W2)
- Goldens change (W3 consequence of renderer flip)
- Deletion of `wake-provider.json` or `prompt.md` (W4)
- Renderer refusal smokes from SKILL.md source (W3/W4)

**State after W1:** Three coexisting sources of truth — `wake-provider.json` (renderer reads), `prompt.md` (renderer reads), `SKILL.md` (I5 validates). No rendered output changes. I5 count advances 91 → 93.

---

## §2. AC oracle list (W1-relevant ACs)

### AC1 — `schemas/skill.cue` adds `"wake"` to `artifact_class` enum + `#Wake` definition

**Oracle (pass conditions):**
- `cue vet` against a conformant wake `SKILL.md` passes with no errors.
- `cue vet` against a wake `SKILL.md` with `wake.role: observer` fails (not in enum).
- `cue vet` against a wake `SKILL.md` with `wake.role: admin` but missing `wake.output.channel_log_convention` fails (role-shaped required field).
- `cue vet` against a wake `SKILL.md` with `wake.role: dispatch` but missing `wake.output.cycle_artifact_root` fails (role-shaped required field).
- `cue vet` against a wake `SKILL.md` with `artifact_class: wake` passes (enum accepts the new value).
- `cue vet` against an existing non-wake `SKILL.md` continues to pass (additive change; no regression).

**Failure scenario:** `cue vet` passes on a `SKILL.md` with `wake.role: observer` → enum is too wide. `cue vet` passes on a dispatch SKILL.md missing `wake.output.cycle_artifact_root` → role-shaped disjunction is not enforced.

**Scope:** additive to `schemas/skill.cue` only. No other file touched.

---

### AC2 — Both wake `SKILL.md` files authored

**Oracle (pass conditions):**
- I5 validator (`scripts/ci/validate-skill-frontmatter.sh`) reports: "93 SKILL.md validated; no findings".
- Each `SKILL.md` has: `artifact_class: wake`, `scope: global`, a `wake:` block carrying all current manifest data per the field mapping in `w0-design.md §B.3`, body = verbatim `prompt.md` content.
- `cue vet` passes on each file independently (AC1 oracle plus AC2 oracle are linked).
- No field from the respective `wake-provider.json` that belongs in frontmatter (per `w0-design.md §B.3`) is absent or misvalued.
- Verbose prose fields (`responsibilities`, `*_notes`, `cross_references`, `trigger_descriptions`, `inbound`, `agent_variable.description`, `concurrency_intent.notes`, `permission_intent_notes`) are in the body, NOT in frontmatter.
- `prompt.md` body section: all text from the respective `prompt.md` file, verbatim, no paraphrase.

**Failure scenario:** I5 validator shows 91 (files missing) or shows 93 but with findings (frontmatter malformed). `cue vet` fails on either file. A frontmatter field is misvalued vs. `wake-provider.json`. A prose field appears in frontmatter instead of body. Body diverges from `prompt.md` verbatim.

**Note on body section order (OB-2):** The SKILL.md body sections MUST appear in this order:
1. Identity
2. Purpose
3. Authority boundary
4. Wake procedure
5. Repair/stop rules
6. Outputs/receipts
7. Non-goals

This body section structure applies to each wake's body. The verbatim `prompt.md` content maps to these sections; α must confirm the mapping before writing.

---

### AC4 — Byte-identical goldens (W1 form: renderer untouched → no golden change)

**Oracle (pass conditions):**
- `install-wake golden` CI check passes: re-render diff clean; sha256 golden == live; idempotent.
- The committed goldens `cnos-agent-admin.golden.yml` and `cnos-cds-dispatch.golden.yml` are byte-identical to what the renderer produces (renderer still reads `wake-provider.json` + `prompt.md`).
- The W1 branch diff contains NO changes to any `*.golden.yml` file.
- The W1 branch diff contains NO changes to `wake-provider.json` or `prompt.md`.

**Failure scenario:** α edits `wake-provider.json` or `prompt.md`, causing a golden mismatch. α changes the golden files directly. Renderer behavior changes unexpectedly because α touched renderer-adjacent code.

**Note:** In W1 the "byte-identical" oracle is trivially satisfied because the renderer still reads JSON+prompt. The oracle becomes load-bearing in W2 (parity gate) and W3 (post-flip golden check). AC4 in W1 is a guardrail confirming α stayed inside scope.

---

### AC7 — All CI gates green

**Oracle (pass conditions, W1-specific):**
- I1 (schema coherence): passes — `schemas/skill.cue` additive change is backward-compatible.
- I2 (skill-graph): passes — two new SKILL.md nodes added; no broken `calls:` edges introduced.
- I4 (unknown): passes — no regression.
- I5 (frontmatter validation): 93 SKILL.md validated, no findings (the W1 oracle from AC2).
- I6 (unknown): passes — no regression.
- `install-wake golden`: passes — byte-identical goldens (AC4 oracle).
- `dispatch-repair-preflight (cnos#516)`: passes — no change to dispatch contract.
- Go: passes — no Go source changes.
- Package: passes — no package registry changes.
- Binary: passes — no binary changes.

**Failure scenario:** I5 reports findings on the new wake SKILL.md files (AC2 gap). `install-wake golden` fails (AC4 violation). I1 fails because the `#Wake` addition contains a CUE syntax error or conflicts with `#Skill`.

---

## §3. Source-of-truth table

| Surface | Path | Role in W1 |
|---|---|---|
| W0 design (locked) | `.cdd/unreleased/524/w0-design.md` | Primary design reference; field mapping in §B.3 is authoritative for what goes in SKILL.md frontmatter |
| CUE skill schema | `schemas/skill.cue` | α extends this (additive only) |
| Admin wake manifest | `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` | Read-only reference for SKILL.md frontmatter values |
| Admin wake prompt | `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` | Read-only reference; body = verbatim content of this file |
| Admin wake golden | `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` | Must be unchanged by W1 |
| Dispatch wake manifest | `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | Read-only reference for SKILL.md frontmatter values |
| Dispatch wake prompt | `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | Read-only reference; body = verbatim content of this file |
| Dispatch wake golden | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | Must be unchanged by W1 |
| I5 validator script | `scripts/ci/validate-skill-frontmatter.sh` | Must NOT be touched; expected to report 93 after W1 |
| Install-wake golden CI | `.github/workflows/install-wake-golden.yml` | Must NOT be touched; must stay green |
| Admin SKILL.md (new) | `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` | α authors this; the primary W1 output for AC2 |
| Dispatch SKILL.md (new) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | α authors this; the primary W1 output for AC2 |

---

## §4. α dispatch prompt

```text
You are α@cdd.cnos, running as the implementer for cycle/524 W1 R0.

Branch: cycle/524 (at main HEAD 23240e4d).
Parent issue: cnos#524 ("wake-as-skill: migrate wakes to typed SKILL.md modules").
Scaffold: .cdd/unreleased/524/gamma-scaffold.md (READ THIS FIRST — the full file — before any work).
W0 design: .cdd/unreleased/524/w0-design.md (READ THIS SECOND — the full file — before any work).

## Scope of this run (MANDATORY — read before any file write)

This run is W1 IMPLEMENTATION. You deliver exactly three file changes:
1. `schemas/skill.cue` — additive: add `"wake"` to `artifact_class` enum + author `#Wake` CUE definition
2. `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` — new file
3. `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — new file

**Hard scope constraints (enforced by guardrail — see §5 of this scaffold):**
- Do NOT touch `wake-provider.json` (either wake)
- Do NOT touch `prompt.md` (either wake)
- Do NOT touch any `*.golden.yml`
- Do NOT touch `scripts/ci/validate-skill-frontmatter.sh`
- Do NOT touch any `.github/workflows/*.yml`
- Do NOT touch `src/packages/cnos.core/commands/install-wake/` (the renderer)
- Do NOT touch any file outside: `schemas/skill.cue`, the two new SKILL.md paths, `.cdd/unreleased/524/`

## Step 1: Read reference files

Before writing anything, read the following (all read-only):
1. `.cdd/unreleased/524/gamma-scaffold.md` — this scaffold (full)
2. `.cdd/unreleased/524/w0-design.md` — W0 design (full); §B.3 field mapping is your authoring key
3. `schemas/skill.cue` — current CUE schema; understand `#Skill` before extending
4. `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — admin manifest values
5. `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` — admin prompt (verbatim body)
6. `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — dispatch manifest values
7. `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — dispatch prompt (verbatim body)

Do NOT skip any of these reads. You will need all of them.

## Step 2: Extend `schemas/skill.cue` (AC1)

Add the following to `schemas/skill.cue`:
- Add `"wake"` to the `artifact_class` enum: `artifact_class?: "skill" | "runbook" | "reference" | "deprecated" | "wake"`
- Add a `#Wake` definition BELOW the `#Skill` definition. The `#Wake` definition MUST embed `#Skill` and constrain the wake-specific fields.

### `#Wake` definition specification

`#Wake` must validate the full `wake:` block per the W0 design (§B of `w0-design.md`). Required fields and types:

```
#Wake: #Skill & {
    artifact_class: "wake"
    scope:          "global"

    wake: {
        role:                 "admin" | "dispatch"
        package:              string
        admin_only:           bool
        activation_log_writer: bool
        activation_state?:    "live" | "declaration-only"
        protocol?:            string
        selector?: {
            include: [...string]
            exclude: [...string]
        }
        input: {
            triggers: [...string]
            issues_opened_title_pattern?: string
        }
        output: #WakeOutputAdmin | #WakeOutputDispatch
        permission_intent: [...string]
        concurrency: {
            serialize: bool
            group:     string
        }
        agent_variable: {
            name:    string
            default: string | null
        }
        surfaces?: {
            allowed:    [...string]
            disallowed: [...string]
        }
        defer_path?: {
            cell_shaped_directive?: string
            off_role_directive?:    string
            ambiguous_directive?:   string
        }
        ...
    }
}

#WakeOutputAdmin: {
    channel_log_convention: string
    writer_surface:         string
    class_taxonomy: [...string]
    cursor_advance: bool
    cursor_field:   string
    ...
}

#WakeOutputDispatch: {
    cycle_artifact_root:     string
    artifact_class_taxonomy: [...string]
    cell_runtime:            string
    ...
}
```

**Key constraints from operator directives and friction notes:**
- OB-1: `wake.role` enum = `"admin" | "dispatch"` ONLY. No `"observer"`.
- FN-3: `wake.agent_variable.default` = `string | null`. The admin wake has `null`; CUE must accept `null`.
- FN-4: `surfaces.allowed` and `surfaces.disallowed` are `[...string]` (open arrays; no length constraint).
- OB-3: `wake.input.triggers` = workflow event triggers (e.g. `schedule`, `issues_opened_title_match`). These are distinct from the top-level `triggers:` skill field (skill-discovery triggers). No overlap in meaning.
- The `#Wake` definition is open at the `wake:` block level (`...`) to pass through renderer fields not yet in the schema — this avoids breaking `cue vet` if wake-provider.json has fields not yet in `#Wake`.
- The role-shaped `output` disjunction must be expressed so that `cue vet` enforces the correct output shape per role. If CUE's open-type disjunction makes this enforcement lossy, use a structural constraint that at minimum requires `channel_log_convention` for admin and `cycle_artifact_root` for dispatch.

**CUE syntax reminder:** Use `...` (not `_`) for open structs. Place `#WakeOutputAdmin` and `#WakeOutputDispatch` as top-level definitions above or below `#Wake`. Keep the file's package declaration (`package skill`) at the top.

After writing, verify the CUE is syntactically valid. If a `cue` tool is available, run:
```
cue vet schemas/skill.cue
```

## Step 3: Author `agent-admin/SKILL.md` (AC2)

Path: `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md`

### Frontmatter (between `---` delimiters)

Use the field mapping in `w0-design.md §B.3` to populate frontmatter. Every field value MUST come from `wake-provider.json` — do not invent or paraphrase.

Required frontmatter structure:
```yaml
---
name: agent-admin
description: <exact value from wake-provider.json "description">
governing_question: <author a single crisp question that the agent executing this wake must answer by the end of the run; consistent with the wake's purpose>
triggers:
  - "wake is invoked by the scheduler or a manual 'claude-wake' issue-open trigger"
scope: global
artifact_class: wake
kata_surface: none
inputs:
  - "home thread inbound directives since last cursor"
  - "open issues in this repo (read-only)"
outputs:
  - "channel log entry (.cn-{agent}/logs/YYYYMMDD.md)"
  - "cursor advance"
wake:
  role: admin
  package: cnos.core
  admin_only: true
  activation_log_writer: true
  input:
    triggers:
      - schedule
      - issues_opened_title_match
    issues_opened_title_pattern: "claude-wake"
  output:
    channel_log_convention: "docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md"
    writer_surface: ".cn-{agent}/logs/YYYYMMDD.md (per the convention's §2; today's file in the current hub)"
    class_taxonomy:
      - heartbeat
      - substantive
      - inaugural
      - directive-out
      - cycle-complete
    cursor_advance: true
    cursor_field: "cursor_out (entry frontmatter); 'agent@<sha>' where <sha> is the home thread HEAD observed during this wake; equal to cursor_in when no advance (heartbeat case)"
  permission_intent:
    - contents.write
    - issues.write
    - pull_requests.write
    - id_token.write
  concurrency:
    serialize: true
    group: "agent-admin-{agent}"
  agent_variable:
    name: agent
    default: null
  surfaces:
    allowed:
      - ".cn-{agent}/logs/ (writer-locality: only the current hub's per-agent log directory)"
      - ".cn-{agent}/state/ (subject to per-installation policy; e.g. cursor advancement for home-side syncs)"
      - "issue comments (on issues in this repo; admin-tone comments only — never implementation work)"
      - "pull request comments (same constraint)"
      - "label application (status:* per operator authorization; protocol:{id} per cnos#468 label-doctrine §4.1; dispatch:cell when an issue clearly meets the cell criteria per the same)"
      - "issue creation (when operator-authorized or standing-posture-aligned; issue refinement same constraint)"
    disallowed:
      - "cell_execution"
      - ".github/workflows/ (substrate is rendered, not hand-edited — the wake never modifies its own carrier or other wake artifacts)"
      - "code files outside .cn-{agent}/ and .cdd/ (cell-shaped work belongs to dispatch wakes; admin wake is read-only on code)"
      - "test files outside .cn-{agent}/ and .cdd/ (same)"
      - "documentation files outside .cn-{agent}/ and .cdd/ (same; admin wake may read any doc but writes only to channel surfaces)"
      - "branch protection rules (operator-owned; admin wake never modifies)"
      - "repository settings (operator-owned; admin wake never modifies)"
      - "label definition (admin wake APPLIES labels per cnos#468 §4.1 but never DEFINES new labels per cnos#468 §4.2)"
      - "other agents' .cn-{other}/ surfaces (writer-locality per AGENT-ACTIVATION-LOG-v0 §0)"
  defer_path:
    cell_shaped_directive: "Defer to the relevant protocol:{P} dispatch wake if installed in this repo (e.g. cnos.cds's cds-dispatch wake for protocol:cds cells; cnos.cdr's cdr-dispatch wake for protocol:cdr cells). If no dispatch wake for the protocol is installed, surface to operator via the channel log entry: name the directive, name the missing dispatch wake, name the cycle the operator should install (e.g. 'cn wake install cds-dispatch'). The admin wake MUST NOT execute the cell inline under any circumstance — doing so collapses the agent-admin / dispatch boundary that cnos#467 establishes."
    off_role_directive: "When a directive falls outside the admin role (e.g. a code-change request, a renderer invocation, a release-cut), the admin wake appends a 'directive-out' or 'substantive' channel entry naming the misroute and the appropriate role/wake/operator action; the admin wake does NOT attempt the work. The channel entry is the operator's signal to re-route."
    ambiguous_directive: "When the directive's role is ambiguous (could be admin or could be cell-shaped), defer to operator per cnos.core/skills/agent/attach §3.8 ('defer to operator on ambiguity, do not guess'); the channel entry names the ambiguity and the candidate interpretations."
---
```

**Note on `agent_variable.default: null`:** YAML `null` is the correct representation. Do NOT use `~`, `""`, or omit the key. This is load-bearing for FN-3 (CUE `null` round-trip).

**Note on `surfaces`:** These values must be copied verbatim from `wake-provider.json` `allowed_surfaces` and `disallowed_surfaces`. Do not paraphrase or truncate.

### Body (after the closing `---`)

The body MUST be structured per OB-2 section order:
1. `## Identity`
2. `## Purpose`
3. `## Authority boundary`
4. `## Wake procedure`
5. `## Repair/stop rules`
6. `## Outputs/receipts`
7. `## Non-goals`

**The body MUST include the verbatim content of `prompt.md`.** Read `prompt.md` carefully and place its content in the appropriate sections. Do not paraphrase, summarize, or reformat the prose.

**Additionally**, include the verbose reference data that moves from `wake-provider.json` to the body (per `w0-design.md §C Decision 5`):
- `responsibilities` array → under `## Wake procedure` or `## Identity` as appropriate
- `cross_references` → at the end of the body
- `class_taxonomy_notes` → near the Outputs/receipts section
- `permission_intent_notes` → near Authority boundary or as an annotated list
- `concurrency_intent.notes` → near Wake procedure or Repair/stop rules
- `superseded_substrate_artifact` → under Non-goals or as a historical note
- `relationship_to_substrate` → under Non-goals or after cross_references
- `agent_variable.description` → under Identity or Wake procedure
- `input_contract.trigger_descriptions` → under Wake procedure
- `input_contract.inbound` → under Wake procedure

**FN-5 (critical):** The `prompt.md` content must be copied verbatim — including all whitespace, line endings, and any trailing newline. Open the file, read it, copy it exactly. Any deviation will break the byte-identity oracle at W2.

## Step 4: Author `cds-dispatch/SKILL.md` (AC2)

Path: `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`

### Frontmatter

Use the same approach as Step 3. Values from `cds-dispatch/wake-provider.json`.

Required frontmatter structure:
```yaml
---
name: cds-dispatch
description: <exact value from wake-provider.json "description">
governing_question: <author a single crisp question consistent with the dispatch wake's purpose>
triggers:
  - "wake is invoked by the scheduler or an issue-labeled event matching the selector"
scope: global
artifact_class: wake
kata_surface: none
inputs:
  - "open issues matching selector (dispatch:cell + protocol:cds + status:todo)"
outputs:
  - "claimed cell driven to completion (.cdd/unreleased/{N}/ artifact set)"
  - "lifecycle label transitions on the claimed cell"
wake:
  role: dispatch
  package: cnos.cds
  admin_only: false
  activation_log_writer: false
  activation_state: live
  protocol: cds
  selector:
    include:
      - dispatch:cell
      - protocol:cds
      - status:todo
    exclude:
      - status:in-progress
      - status:blocked
      - status:review
      - status:changes
  input:
    triggers:
      - schedule
      - issues_labeled_selector_match
  output:
    cycle_artifact_root: ".cdd/unreleased/{N}/"
    artifact_class_taxonomy:
      - gamma-scaffold
      - self-coherence
      - beta-review
      - alpha-closeout
      - beta-closeout
      - gamma-closeout
      - post-release-assessment
    cell_runtime: cnos.cdd
  permission_intent:
    - contents.write
    - issues.write
    - pull_requests.write
    - id_token.write
  concurrency:
    serialize: true
    group: "cds-dispatch-{agent}"
  agent_variable:
    name: agent
    default: sigma
  surfaces:
    allowed:
      - ".cdd/unreleased/{N}/ (the cell artifact root; γ/α/β write here per the cnos.cdd artifact contract)"
      - "cell-cycle branches (e.g. cycle/{N}; α/β commit code/test/doc changes here per the cell's design surface)"
      - "pull request creation and updates (each cell ships its work via a PR scoped to the claimed cell)"
      - "issue comments on the claimed cell (status updates, dispatch_protocol_missing / dispatch_protocol_mismatch diagnostics)"
      - "label application on the claimed cell only (status:todo → status:in-progress at claim; lifecycle transitions per responsibilities #5 above)"
    disallowed:
      - ".github/workflows/ (substrate is rendered, not hand-edited — the dispatch wake never modifies its own carrier or other wake artifacts)"
      - ".cn-{agent}/ surfaces (channel logs are admin's writer-locality per AGENT-ACTIVATION-LOG-v0 §0; dispatch wake does NOT write channel entries)"
      - "label definition (dispatch wake APPLIES lifecycle labels on its claimed cell per cnos#468 §4.1 but never DEFINES new labels per cnos#468 §4.2)"
      - "branch protection rules (operator-owned; dispatch wake never modifies)"
      - "repository settings (operator-owned; dispatch wake never modifies)"
      - "cells outside its protocol (cross-protocol claims are a protocol violation per cnos#468 §2.1; a wake owning protocol:cds rejects protocol:cdr / protocol:cdw issues per cnos#454 AC9)"
      - "issues not matching the selector (the wake never edits, labels, or comments on issues that do not match its claim criteria, except for diagnostic comments on dispatch_protocol_missing cases)"
      - "other agents' .cn-{other}/ surfaces (writer-locality per AGENT-ACTIVATION-LOG-v0 §0)"
  defer_path:
    cell_shaped_directive: "The dispatch wake's primary job is cell execution; cell-shaped directives that match the selector are not deferred — they are claimed and run. For cell-shaped directives that arrive via a channel that is NOT the open-issue queue (e.g. a comment on an issue not labeled for dispatch, or a direct prompt at wake firing time), the wake surfaces to the admin wake via an issue comment: 'Cell-shaped directive arrived outside the standard claim channel; please re-route via dispatch:cell + protocol:cds + status:todo labels if intended for execution.'"
    off_role_directive: "When a directive falls outside the dispatch role (admin work: channel sync, label policy decisions, repo settings, agent identity), the dispatch wake appends a comment on the claimed cell (if any) or surfaces to the admin wake via a comment on the relevant issue: name the misroute and the appropriate role/wake/operator action. The dispatch wake does NOT attempt admin work inline."
    ambiguous_directive: "When the directive's role is ambiguous (could be dispatch or could be admin / off-role), defer to operator: post an issue comment naming the ambiguity and the candidate interpretations; release the claim if one was taken; do not guess. Mirrors cnos.core/skills/agent/attach §3.8 ('defer to operator on ambiguity, do not guess') for the dispatch context."
---
```

### Body

Same approach as Step 3. Body per OB-2 section order; verbatim `prompt.md` content in the appropriate sections; verbose reference data from `wake-provider.json` (responsibilities, cross_references, notes fields) distributed into body sections.

Note: `activation_state_notes`, `artifact_class_notes`, `cell_runtime_notes`, `trigger_descriptions`, `inbound`, `agent_variable.description`, `concurrency_intent.notes`, `permission_intent_notes` all move to body. The body should reference the activation_state liveness clearly (the dispatch wake's `activation_state: live` is significant — the body must communicate the live-state context the same way `prompt.md` does).

## Step 5: Self-coherence

Write `.cdd/unreleased/524/self-coherence.md §R0`:
- Confirm scope guardrails: list which files were NOT touched
- For each of AC1 / AC2 / AC4 / AC7: state the oracle and whether it is satisfied
- Walk the `#Wake` definition: confirm `role` enum is `admin | dispatch` only (OB-1); confirm `agent_variable.default: string | null` (FN-3); confirm `surfaces.allowed/disallowed` use `[...string]` (FN-4)
- Walk each SKILL.md: confirm body section order matches OB-2; confirm `prompt.md` content is verbatim (FN-5); confirm OB-3 (no overlap between top-level `triggers:` and `wake.input.triggers`)
- Note FN-1 status: confirm no `observer` in the schema; note whether `observer` appears in the renderer (do NOT touch the renderer — just note the finding for W3 implementer)
- Note FN-2 status: note that I5 validator extracts frontmatter before `cue vet`; confirm or flag whether the second `---` delimiter handling is observable from the validator output
- Note FN-6: confirm `activation_log_writer: false` on dispatch wake + role=dispatch is correctly declared; the renderer refusal smoke (exit 4 on mis-declaration) is W3 scope — just note it
- Note FN-7 and FN-8 as deferred to W2/W4

## Commits

1. Commit `schemas/skill.cue` as: `α-524 W1: add #Wake to skill.cue (AC1)`
2. Commit both wake SKILL.md files as: `α-524 W1: author agent-admin + cds-dispatch SKILL.md (AC2)`
3. Commit `.cdd/unreleased/524/self-coherence.md §R0` as: `α-524 W1 R0: self-coherence §R0`
4. Push all commits to `cycle/524`
```

---

## §5. β review prompt

```text
You are β@cdd.cnos, running as the reviewer for cycle/524 W1 R0.

Branch: cycle/524 (read α's R0 commits; full diff vs. main HEAD 23240e4d).
Parent issue: cnos#524 ("wake-as-skill: migrate wakes to typed SKILL.md modules").
Scaffold: .cdd/unreleased/524/gamma-scaffold.md (read in full).
W0 design: .cdd/unreleased/524/w0-design.md (read in full).
α self-coherence: .cdd/unreleased/524/self-coherence.md §R0 (read in full).

## Scope reminder

W1 delivers exactly three writable changes: `schemas/skill.cue` (additive), `agent-admin/SKILL.md` (new), `cds-dispatch/SKILL.md` (new). Any other source/schema/workflow change is an immediate **iterate**.

## Review checklist

### Scope compliance (check first — iterate immediately on any failure)

- [ ] No changes to `wake-provider.json` (either wake)
- [ ] No changes to `prompt.md` (either wake)
- [ ] No changes to any `*.golden.yml`
- [ ] No changes to `scripts/ci/validate-skill-frontmatter.sh`
- [ ] No changes to any `.github/workflows/*.yml`
- [ ] No changes to the renderer (`src/packages/cnos.core/commands/install-wake/`)
- [ ] Changes outside `.cdd/unreleased/524/` are ONLY: `schemas/skill.cue`, `agent-admin/SKILL.md`, `cds-dispatch/SKILL.md`

### AC1 — `schemas/skill.cue` `#Wake` definition

- [ ] `"wake"` added to `artifact_class` enum (alongside existing values; no existing value removed or changed)
- [ ] `#Wake` definition present and syntactically valid CUE
- [ ] `wake.role` enum = `"admin" | "dispatch"` ONLY — no `"observer"` (OB-1)
- [ ] `wake.agent_variable.default` type = `string | null` — not `string` alone, not `string | ""` (FN-3)
- [ ] `wake.surfaces.allowed` and `wake.surfaces.disallowed` typed as `[...string]` — no length constraint (FN-4)
- [ ] `wake.output` expresses a role-shaped disjunction: admin shape requires `channel_log_convention`; dispatch shape requires `cycle_artifact_root`
- [ ] `wake.input.triggers: [...string]` — distinct from top-level `triggers:` (OB-3)
- [ ] `#Wake` embeds or extends `#Skill` — does not duplicate `#Skill` fields
- [ ] `schemas/skill.cue` `package skill` declaration still at top; no other existing definition altered

Load the two wake SKILL.md files and mentally run `cue vet` against `#Wake`:
- [ ] Both files would pass `cue vet` with no errors
- [ ] A synthetic wake SKILL.md with `wake.role: observer` would fail `cue vet`

### AC2 — `agent-admin/SKILL.md`

Read `agent-admin/SKILL.md` and `agent-admin/wake-provider.json` side by side:

- [ ] `artifact_class: wake` and `scope: global` present in frontmatter
- [ ] `wake.role: admin` — value matches `wake-provider.json` `"role": "admin"`
- [ ] `wake.package: cnos.core` — matches JSON `"package": "cnos.core"`
- [ ] `wake.admin_only: true` — matches JSON `"admin_only": true`
- [ ] `wake.activation_log_writer: true` — matches JSON `"activation_log_writer": true`
- [ ] `wake.input.triggers` = `[schedule, issues_opened_title_match]` — matches JSON `input_contract.triggers`
- [ ] `wake.input.issues_opened_title_pattern: "claude-wake"` — matches JSON `input_contract.issues_opened_title_pattern`
- [ ] `wake.output.channel_log_convention` matches JSON `output_contract.channel_log_convention` verbatim
- [ ] `wake.output.writer_surface` matches JSON `output_contract.writer_surface` verbatim
- [ ] `wake.output.class_taxonomy` = 5 values, matches JSON `output_contract.class_taxonomy` exactly
- [ ] `wake.output.cursor_advance: true` matches JSON
- [ ] `wake.output.cursor_field` matches JSON verbatim
- [ ] `wake.permission_intent` = 4 values matching JSON `permission_intent` exactly
- [ ] `wake.concurrency.serialize: true` matches JSON `concurrency_intent.serialize`
- [ ] `wake.concurrency.group: "agent-admin-{agent}"` matches JSON `concurrency_intent.group`
- [ ] `wake.agent_variable.name: agent` matches JSON
- [ ] `wake.agent_variable.default: null` — YAML null (not `~`, not `""`, not absent) (FN-3)
- [ ] `wake.surfaces.allowed` = 6 entries matching JSON `allowed_surfaces` verbatim
- [ ] `wake.surfaces.disallowed` = 9 entries matching JSON `disallowed_surfaces` verbatim
- [ ] `wake.defer_path.*` values match JSON `defer_path.*` verbatim
- [ ] No prose fields in frontmatter: `responsibilities`, `*_notes`, `cross_references`, `trigger_descriptions`, `inbound`, `agent_variable.description`, `permission_intent_notes`, `concurrency_intent.notes`, `superseded_substrate_artifact`, `relationship_to_substrate` — these must be ABSENT from frontmatter and PRESENT in body
- [ ] Body section order matches OB-2: Identity, Purpose, Authority boundary, Wake procedure, Repair/stop rules, Outputs/receipts, Non-goals
- [ ] Body contains verbatim `prompt.md` content — diff `agent-admin/prompt.md` against the body prose; no paraphrase, no omission (FN-5)
- [ ] OB-3: top-level `triggers:` is skill-discovery triggers; `wake.input.triggers` is workflow event triggers; they are semantically distinct

### AC2 — `cds-dispatch/SKILL.md`

Read `cds-dispatch/SKILL.md` and `cds-dispatch/wake-provider.json` side by side:

- [ ] `artifact_class: wake` and `scope: global` present in frontmatter
- [ ] `wake.role: dispatch` — matches JSON
- [ ] `wake.package: cnos.cds` — matches JSON
- [ ] `wake.admin_only: false` — matches JSON
- [ ] `wake.activation_log_writer: false` — matches JSON
- [ ] `wake.activation_state: live` — matches JSON
- [ ] `wake.protocol: cds` — matches JSON
- [ ] `wake.selector.include` = `[dispatch:cell, protocol:cds, status:todo]` — matches JSON
- [ ] `wake.selector.exclude` = `[status:in-progress, status:blocked, status:review, status:changes]` — matches JSON
- [ ] `wake.input.triggers` = `[schedule, issues_labeled_selector_match]` — matches JSON
- [ ] `wake.input` has NO `issues_opened_title_pattern` (dispatch wake does not have this field)
- [ ] `wake.output.cycle_artifact_root: ".cdd/unreleased/{N}/"` — matches JSON
- [ ] `wake.output.artifact_class_taxonomy` = 7 values matching JSON exactly
- [ ] `wake.output.cell_runtime: cnos.cdd` — matches JSON
- [ ] `wake.permission_intent` = 4 values matching JSON exactly
- [ ] `wake.concurrency.serialize: true` and `wake.concurrency.group: "cds-dispatch-{agent}"`
- [ ] `wake.agent_variable.name: agent` and `wake.agent_variable.default: sigma` (string, not null)
- [ ] `wake.surfaces.allowed` = 5 entries matching JSON `allowed_surfaces` verbatim
- [ ] `wake.surfaces.disallowed` = 8 entries matching JSON `disallowed_surfaces` verbatim
- [ ] `wake.defer_path.*` values match JSON verbatim
- [ ] No prose fields in frontmatter (same rule as admin wake)
- [ ] Body section order matches OB-2
- [ ] Body contains verbatim `prompt.md` content (FN-5); no paraphrase
- [ ] OB-3: same check as admin wake

### AC4 — Byte-identical goldens (W1 form)

- [ ] `cnos-agent-admin.golden.yml` is unchanged in the branch diff
- [ ] `cnos-cds-dispatch.golden.yml` is unchanged in the branch diff
- [ ] `wake-provider.json` unchanged for both wakes
- [ ] `prompt.md` unchanged for both wakes

### AC7 — CI gate readiness

- [ ] `schemas/skill.cue` CUE syntax is valid (no syntax error would block I1/I5)
- [ ] Both SKILL.md files have `artifact_class: wake` and valid frontmatter structure (I5 readiness)
- [ ] No `.github/workflows/*.yml` changes (install-wake-golden CI not disturbed)
- [ ] No `dispatch-repair-preflight`-relevant files changed (cnos#516 gate should stay green)
- [ ] No Go / Package / Binary source changes

### Friction note checks

- [ ] FN-1: `observer` absent from `#Wake` role enum; noted in self-coherence
- [ ] FN-2: self-coherence notes the I5 extractor's second-`---` behavior
- [ ] FN-3: `agent_variable.default: null` correctly typed in CUE (`string | null`) AND correctly spelled in YAML (literal `null`)
- [ ] FN-4: surfaces arrays use `[...string]`; admin has 6 allowed + 9 disallowed; dispatch has 5 allowed + 8 disallowed
- [ ] FN-5: body verbatim from `prompt.md`; no paraphrase detected
- [ ] FN-6: dispatch wake's `activation_log_writer: false` correctly declared; refusal smoke noted as W3 scope
- [ ] FN-7: parity gate noted as W2 scope in self-coherence
- [ ] FN-8: deletion sequencing noted as W4 scope in self-coherence

## Verdict

- **converge** if: scope compliance clean; AC1 / AC2 / AC4 / AC7 checklists fully green; all FN checks satisfied
- **iterate** if: any scope violation; any field value mismatch vs. `wake-provider.json`; `observer` in role enum; wrong `agent_variable.default` type; any `*.golden.yml` change; any prose field in frontmatter; body not verbatim from `prompt.md`; CUE syntax error; OB-2 / OB-3 violated

Write `.cdd/unreleased/524/beta-review.md §R0` with verdict.
Commit as: `β-524 W1 R0 review: <converge|iterate>`
Push to `cycle/524`.
```

---

## §6. Scope guardrails (mechanical — W1 hard constraints)

α MUST NOT touch any of:
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md`
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
- `src/packages/cnos.core/commands/install-wake/` (the renderer; any file therein)
- `scripts/ci/validate-skill-frontmatter.sh`
- Any `.github/workflows/*.yml`
- Any file in `src/` outside the two new SKILL.md paths
- `.cdd/unreleased/524/w0-design.md` (locked; read-only)

**The only permitted new file writes outside `.cdd/unreleased/524/` are:**
1. `schemas/skill.cue` — modified (additive only)
2. `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` — new file
3. `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — new file

β's first review action: run `git diff main...HEAD -- ':!.cdd/'` and confirm only these three paths appear.

---

## §7. Implementation contract (7 axes)

| Axis | W1 contract |
|---|---|
| **Schema** | `schemas/skill.cue` receives `"wake"` in `artifact_class` enum and a `#Wake` definition; additive; no existing definition altered |
| **Modules** | Two new SKILL.md files authored: `agent-admin/SKILL.md` and `cds-dispatch/SKILL.md`; no existing SKILL.md altered |
| **Renderer** | Unchanged; still reads `wake-provider.json` + `prompt.md`; the SKILL.md files are inert from the renderer's perspective in W1 |
| **Goldens** | Byte-identical; `install-wake-golden.yml` CI green; no golden file touched |
| **I5 count** | 91 → 93 (two new SKILL.md files validated by the existing I5 validator with no script changes) |
| **Friction notes** | FN-1 through FN-8 all acknowledged; FN-1/2/3/4/5 are W1 implementation concerns; FN-6/7/8 are W3/W2/W4 scope (noted but not implemented) |
| **Operator directives** | OB-1 (`admin\|dispatch` only in role enum), OB-2 (body section order), OB-3 (triggers separation) are all W1 implementation constraints |

---

## §8. Friction notes — W1 relevance

### FN-1: Renderer role-enum collision
**W1 relevance: ACTIVE — schema decision.** The renderer (`cn-install-wake`) may accept `observer` in its internal role check. Per OB-1, the `#Wake` CUE schema defines `wake.role` as `"admin" | "dispatch"` only. α must NOT include `observer` in the CUE enum. If `observer` appears in any live wake-provider.json or in the renderer's role check, it is a renderer-internal extension — it is NOT a valid `wake.role` in the CUE schema. The W1 implementer does not touch the renderer. Note the discrepancy in self-coherence for the W3 implementer.

### FN-2: I5 extractor and second `---` delimiter
**W1 relevance: VERIFICATION.** The I5 validator script (`validate-skill-frontmatter.sh`) extracts frontmatter (content between the first and second `---` delimiters) before running `cue vet`. The SKILL.md body (after the closing `---`) must not be passed to `cue vet`. α should confirm this behavior is correct by running the I5 validator after authoring both SKILL.md files, or at minimum note in self-coherence that the extractor behavior was verified. If the extractor strips the body correctly, the CUE schema does not need to accommodate the body's prose. If it does not strip the body, `cue vet` will fail on the prose content.

### FN-3: `agent_variable.default: null`
**W1 relevance: LOAD-BEARING — admin wake YAML + CUE.** The admin wake has `agent_variable.default: null`. Two things must be correct: (a) the YAML in `agent-admin/SKILL.md` must spell this as `default: null` (not `~`, not `""`, not absent); (b) the CUE `#Wake` definition must type this as `string | null`, not just `string`. If (b) is typed as `string` only, `cue vet` on `agent-admin/SKILL.md` will fail with a type error. The dispatch wake has `agent_variable.default: sigma` (a string), so (a) and (b) apply only to admin, but the CUE type must accommodate both. Test case: `cue vet` a synthetic wake with `default: null` → must pass.

### FN-4: Long `surfaces` arrays
**W1 relevance: ACTIVE — CUE type + YAML authoring.** Both wakes have 6–9 entries in `allowed_surfaces`/`disallowed_surfaces`. The CUE type must be `[...string]` (open array; no length constraint). In YAML, these must be copied verbatim from `wake-provider.json`. The admin wake has 6 allowed and 9 disallowed surfaces; the dispatch wake has 5 allowed and 8 disallowed surfaces. Count the entries when authoring SKILL.md and again when reviewing.

### FN-5: Body verbatim from `prompt.md`
**W1 relevance: LOAD-BEARING — AC2 and future W2 byte-identity oracle.** The SKILL.md body for each wake must equal the content of `prompt.md` verbatim. α must open `prompt.md`, copy the content exactly (including trailing newline, if any), and embed it in the SKILL.md body. No paraphrase, no reformatting, no summarization. The body additionally includes verbose reference data from `wake-provider.json`, but the `prompt.md` prose must be identical. Any deviation from verbatim will cause the W2 parity gate (`render(SKILL) cmp render(JSON)`) to fail, requiring a W1 re-open to fix the body.

### FN-6: Dispatch wake `activation_log_writer: false` + renderer refusal
**W1 relevance: DECLARATION ONLY.** The dispatch SKILL.md correctly declares `wake.activation_log_writer: false`. The renderer refusal (exit code 4 when `role:dispatch + admin_only:false + activation_log_writer:true` is mis-declared) is preserved in the renderer, but the refusal smoke from SKILL.md source is AC6, which is W3 scope. In W1, α must ensure the declaration is correct (`false`). The renderer still reads from `wake-provider.json`, so the W1 refusal smoke is unchanged. Note this in self-coherence.

### FN-7: Dual-source parity gate
**W1 relevance: DEFERRED — W2 scope.** The parity gate (`render(SKILL) cmp render(JSON)`) is not implemented in W1. The SKILL.md body verbatim constraint (FN-5) is the W1 preparation for this gate. Note in self-coherence that the parity gate design is deferred to W2.

### FN-8: Deletion sequencing
**W1 relevance: DEFERRED — W4 scope.** No files are deleted in W1. Note in self-coherence that deletion of `wake-provider.json` and `prompt.md` is deferred to W4 after the byte-identical proof holds.

---

## §9. Cross-references

- Issue #524 body: `## W0 design (locked)` — the locked design governing W1 implementation
- `.cdd/unreleased/524/w0-design.md` — W0 design document (read before authoring); §B.3 field mapping is the authoritative lookup table for frontmatter values
- `schemas/skill.cue` — file α extends (current `#Skill` definition is the CUE idiom to follow)
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — admin wake manifest (frontmatter value source)
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` — admin wake prompt (verbatim body source)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — dispatch wake manifest (frontmatter value source)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — dispatch wake prompt (verbatim body source)
- `.github/workflows/install-wake-golden.yml` — CI gate α must not break (read-only)
- `scripts/ci/validate-skill-frontmatter.sh` — I5 validator α must not touch (read-only; expected output: 93 validated)
- Operator W1 directives: OB-1 (`admin|dispatch` only), OB-2 (body section order), OB-3 (trigger separation)
- Prior γ-scaffold (W0 scope, now superseded by this file): archived in git history

---

_Filed by γ@cdd.cnos (wake-invoked δ dispatch, cycle/524), 2026-06-30 (UTC)._
_Scope: W1 implementation (AC1 + AC2 + AC4 + AC7). Supersedes the W0-scope γ-scaffold._
