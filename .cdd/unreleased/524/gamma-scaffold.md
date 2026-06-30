---
cycle: 524
parent_issue: cnos#524
protocol: cds
cycle_branch: cycle/524
base_main_sha: d7d2244cb4cb95b94e28569ac8ef090d49876c89
head_sha_at_scaffold: d7d2244cb4cb95b94e28569ac8ef090d49876c89
mode: MCA
role: γ
authored_by: γ@cdd.cnos (wake-invoked δ dispatch)
date: 2026-06-30 (UTC)
output_contract: γ-scaffold + α prompt + β prompt
ac_set: cnos#524 body ACs 1–7
w0_locked: true
---

# γ-scaffold — cnos#524: wake-as-skill migration W1→W4

## §1. Issue framing

**The problem.** Wakes today use two artifact systems in parallel: `wake-provider.json` (machine-readable manifest) + `prompt.md` (prompt body), compiled by `cn install-wake` into `.github/workflows/cnos-<wake>.yml`. The SKILL.md contract system (frontmatter + body + CUE validation by I5) runs beside this — not over it. Two contract systems means two maintenance surfaces, two validation pipelines, and no single source of truth for what a wake is.

**The goal.** One contract system. A wake becomes a typed SKILL.md module: `artifact_class: wake`, a `wake:` frontmatter block validated by CUE carrying exactly what `wake-provider.json` carries today, and a body that IS `prompt.md` verbatim. The renderer reads SKILL.md; JSON + prompt are deleted after byte-identical proof.

**Migration shape.** Four phases (W0 operator-ratified and locked):

- **W1:** Author both SKILL.md files beside existing JSON/prompt (body = prompt.md verbatim). Add `#Wake` + `artifact_class: wake` to `skill.cue`. Renderer still reads JSON → goldens unchanged.
- **W2:** Teach `cn-install-wake` to read SKILL.md; dual-source parity gate: render(SKILL) byte-identical to render(JSON).
- **W3:** Flip renderer source to SKILL.md; re-render (no-op if byte-identical).
- **W4:** Delete `wake-provider.json` + `prompt.md` for both wakes after byte-identical proof; I5 now covers 91→93 modules.

**Byte-identity is the migration oracle at every step.**

## §2. Source-of-truth table

| Surface | Path |
|---|---|
| CUE schema | `schemas/skill.cue` |
| I5 validator | `scripts/ci/validate-skill-frontmatter.sh` |
| Renderer | `src/packages/cnos.core/commands/install-wake/cn-install-wake` |
| Golden CI | `.github/workflows/install-wake-golden.yml` |
| agent-admin manifest | `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` |
| agent-admin prompt | `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` |
| agent-admin golden | `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` |
| cds-dispatch manifest | `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` |
| cds-dispatch prompt | `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` |
| cds-dispatch golden | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` |
| Live admin workflow | `.github/workflows/cnos-agent-admin.yml` |
| Live dispatch workflow | `.github/workflows/cnos-cds-dispatch.yml` |
| agent-admin SKILL.md (to create) | `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` |
| cds-dispatch SKILL.md (to create) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` |

## §3. Surfaces α will touch

### Files to READ (source inputs):
1. `schemas/skill.cue` — current `#Skill` definition; where to graft `#Wake` and extend `artifact_class`
2. `src/packages/cnos.core/commands/install-wake/cn-install-wake` — every field consumed from `wake-provider.json` (lines 283–376); every renderer validation path
3. `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — agent-admin manifest data
4. `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` — agent-admin prompt body (verbatim into SKILL.md body)
5. `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — cds-dispatch manifest data
6. `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — cds-dispatch prompt body (verbatim into SKILL.md body)
7. `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — current golden (byte-identity anchor)
8. `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — current golden (byte-identity anchor)
9. `.github/workflows/install-wake-golden.yml` — CI oracle; where to add dual-source parity gate and W4 deletion proof
10. `scripts/ci/validate-skill-frontmatter.sh` — I5 validator; confirms discovery pattern (`find src/packages -name SKILL.md`) picks up orchestrator SKILL.mds automatically
11. `schemas/skill-exceptions.json` — check if any exception entries affect the new SKILL.mds
12. An existing SKILL.md with a full field set (e.g., `src/packages/cnos.cds/skills/cds/SKILL.md`) for format reference

### Files to MODIFY:
1. `schemas/skill.cue` — add `"wake"` to `artifact_class` enum; add `#Wake` definition
2. `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` — **CREATE** (W1 deliverable)
3. `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — **CREATE** (W1 deliverable)
4. `src/packages/cnos.core/commands/install-wake/cn-install-wake` — W2: SKILL.md reader path; W2: dual-source parity; W3: flip source; W4: operate from SKILL.md only
5. `.github/workflows/install-wake-golden.yml` — add: dual-source parity step (W2), SKILL.md-source render verification (W3), deletion-proof step (W4), I5 count step (91→93)
6. `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — W3 re-render (no-op expected if byte-identical); W4 confirm unchanged
7. `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — W3 re-render (no-op expected); W4 confirm unchanged
8. `.github/workflows/cnos-agent-admin.yml` — W3/W4 re-render for live substrate (no-op expected)
9. `.github/workflows/cnos-cds-dispatch.yml` — W3/W4 re-render for live substrate (no-op expected)

### Files to DELETE (W4 only, after byte-identical proof):
1. `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`
2. `src/packages/cnos.core/orchestrators/agent-admin/prompt.md`
3. `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`
4. `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md`

## §4. AC oracle table

### AC1 — `#Wake` schema in CUE

**Oracle:** `cue vet -d '#Wake' schemas/skill.cue <wake-frontmatter.yaml>` passes for a conformant wake SKILL.md frontmatter; fails with a diagnostic naming the bad field for a manifest with a bad enum or missing required field.

**Testable boundary:** `schemas/skill.cue` adds:
- `"wake"` to the `artifact_class` enum in `#Skill`
- A new `#Wake` definition with the W0-locked shape (role enum; bool fields; enums; nested objects per §W0 design block)

**Filesystem check:**
```sh
grep -c '"wake"' schemas/skill.cue  # must return >= 1
cue vet -d '#Wake' schemas/skill.cue /tmp/conformant-wake.yaml  # exit 0
cue vet -d '#Wake' schemas/skill.cue /tmp/bad-enum-wake.yaml    # exit nonzero + diagnostic
```

**Key `#Wake` fields (from W0 design + renderer field consumption at lines 283–523 of cn-install-wake):**

| Field | Type | Notes |
|---|---|---|
| `wake.role` | `"admin" \| "dispatch"` | renderer line 342: enum check; cross-field invariants |
| `wake.package` | `string` | renderer line 281 schema check |
| `wake.admin_only` | `bool` | renderer line 309/353 |
| `wake.activation_log_writer` | `bool` | renderer line 452 |
| `wake.activation_state` | optional `"live" \| "declaration-only"` | renderer line 386 |
| `wake.protocol` | optional `string` | required when `role == "dispatch"` per renderer line 367 |
| `wake.selector` | optional `{include: [...string], exclude: [...string]}` | required when `role == "dispatch"` per renderer lines 369-372 |
| `wake.input.triggers` | `[...string]` | renderer line 338 |
| `wake.output` | role-shaped disjunction | admin: `{channel_log_convention: string}`; dispatch: `{cycle_artifact_root: string, artifact_class_taxonomy: [...string], cell_runtime: string}` |
| `wake.permission_intent` | `[...string]` | renderer lines 614-628 |
| `wake.concurrency` | `{serialize: bool, group: string}` | renderer lines 493-495 |
| `wake.agent_variable` | `{name: string, default: string \| null}` | renderer lines 85-89 |

**Non-goal:** `#Wake` does NOT carry `responsibilities`, `*_notes`, `cross_references` — those move to body per W0 call #5.

### AC2 — Wake SKILL.md modules authored

**Oracle:** `scripts/ci/validate-skill-frontmatter.sh` validates both new SKILL.md files — total count advances from 91 to 93; self-test green.

**Testable boundary:**
- `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` exists with `scope: global`, `artifact_class: wake`, `wake:` block carrying exact current `wake-provider.json` data; body = current `prompt.md` verbatim.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` exists with same shape; dispatch-specific `wake:` fields.
- Verbose fields (`responsibilities`, `*_notes`, `cross_references`) are in the body, not frontmatter.

**Filesystem check:**
```sh
ls src/packages/cnos.core/orchestrators/agent-admin/SKILL.md  # exists
ls src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md  # exists
./scripts/ci/validate-skill-frontmatter.sh 2>&1 | grep 'SKILL.md validated'  # shows 93
./scripts/ci/validate-skill-frontmatter.sh --self-test  # exits 0
```

**I5 discovery note:** `validate-skill-frontmatter.sh` runs `find src/packages -name SKILL.md`. Both orchestrator dirs are under `src/packages/`; SKILL.md files placed there are discovered automatically — no exception-list changes needed.

### AC3 — Renderer reads SKILL.md

**Oracle:** `cn install-wake <name>` when given a SKILL.md-bearing directory (W2 dual-source) produces output byte-identical to the JSON-sourced render; after W3 flip, renders from SKILL.md only.

**Testable boundary:** renderer code path in `cn-install-wake`:
- W2: add `--source skill-md` or auto-detect SKILL.md beside manifest; render from SKILL.md; compare bytes to render-from-JSON
- W3: flip the default source to SKILL.md; the JSON path becomes a fallback or is removed

**Filesystem check:**
```sh
# W2 parity (before JSON deletion):
cn install-wake agent-admin --source skill-md --out /tmp/from-skill.yml
cn install-wake agent-admin --source json      --out /tmp/from-json.yml
cmp /tmp/from-skill.yml /tmp/from-json.yml     # must exit 0

# W3 (after flip):
cn install-wake agent-admin --out /tmp/from-default.yml
cmp /tmp/from-default.yml src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

### AC4 — Byte-identical goldens

**Oracle:** `install-wake golden` CI green: re-render diff clean + sha256 golden==live + idempotence; one-shot dual-source parity check `render(SKILL) cmp render(JSON)` passes before deletion.

**Testable boundary:**
- Both goldens render byte-identically from SKILL.md versus from JSON
- sha256 of live `.github/workflows/cnos-*.yml` equals sha256 of the package golden
- Second render is a no-op (idempotence)

**Filesystem check:**
```sh
sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
sha256sum .github/workflows/cnos-agent-admin.yml
# must match

sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
sha256sum .github/workflows/cnos-cds-dispatch.yml
# must match
```

### AC5 — JSON + prompt deleted

**Oracle:** `wake-provider.json` and `prompt.md` absent for both wakes; `cn install-wake` still renders byte-identical goldens.

**Testable boundary:** files absent + renderer still produces correct output from SKILL.md only.

**Filesystem check:**
```sh
ls src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json  # must NOT exist
ls src/packages/cnos.core/orchestrators/agent-admin/prompt.md            # must NOT exist
ls src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json   # must NOT exist
ls src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md            # must NOT exist
cn install-wake agent-admin   # exit 0; golden unchanged
cn install-wake cds-dispatch  # exit 0; golden unchanged
```

### AC6 — Renderer refusals preserved

**Oracle:** existing cycle/496 mis-declaration smoke in `install-wake-golden.yml` still passes; a synthetic mis-declared wake SKILL.md → exit 4.

**Testable boundary:** after W3/W4, the renderer reads the `wake:` block from SKILL.md and applies the same refusal logic previously applied against JSON fields:
- `activation_state != "live"` → exit 3
- `role:dispatch + admin_only:false + activation_log_writer:true` → exit 4
- `activation_log_writer:false + attach in cross_references` → exit 4 (FN-6 defense-in-depth)

**Filesystem check:**
```sh
# Existing AC4/AC5 smoke in install-wake-golden.yml must still pass:
# Step "AC4 (cycle/496) — renderer refusal on mis-declaration" exits 4.
# Step "AC5 — declaration-only refusal" exits 3.
# If synthetic fixture JSONs are deleted (or kept alongside), verify the smoke
# is wired to SKILL.md-sourced equivalents post-W4.
```

### AC7 — All gates green

**Oracle:** I1/I2/I4/I5/I6/`install-wake golden`/`dispatch-repair-preflight`/Go/Package/Binary all green on the cell's PR.

**Testable boundary:** no inherited-cap failures; no new red gates.

**Filesystem check:** CI run on the PR passes all jobs.

## §5. Implementation contract (for α dispatch prompt)

| Axis | Value |
|---|---|
| **Branch** | `cycle/524` |
| **Base SHA** | `d7d2244cb4cb95b94e28569ac8ef090d49876c89` |
| **Issue** | cnos#524 |
| **Mode** | MCA — W0 operator-ratified; α executes W1→W4 sequence |
| **AC set** | AC1–AC7 (all) |
| **W0 non-goals** | No rendered-workflow output change; no runtime-behavior change; no renderer substrate-encoding change; no new wake names; no schema changes outside `schemas/skill.cue`; no I5 exception-list additions |
| **Byte-identity oracle** | `cmp render(SKILL) render(JSON)` must pass at every intermediate step before deletion |
| **Deletion gate** | W4 deletions ONLY after byte-identical proof holds in CI |

## §6. α dispatch prompt

```text
You are α@cdd.cnos, running as the implementer for cycle/524.

Branch: cycle/524 (at base main SHA d7d2244cb4cb95b94e28569ac8ef090d49876c89).
Parent issue: cnos#524 ("wake-as-skill: migrate wakes to typed SKILL.md modules").
Scaffold: .cdd/unreleased/524/gamma-scaffold.md (READ THIS FIRST in full before any edits).

## What you implement

Implement R0 of the W1→W4 migration per the operator-ratified W0 design (locked; no design reframing without operator approval). The 7 ACs are the acceptance surface. The byte-identity oracle is the migration correctness proof at every step.

## Skills to load before implementing

Load these skills in order:
1. `src/packages/cnos.cds/skills/cds/SKILL.md` — CDS protocol discipline
2. `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract
3. `src/packages/cnos.cdd/skills/cdd/SKILL.md` — cell-runtime framework

## Implementation contract

| Axis | Value |
|---|---|
| Branch | cycle/524 |
| Base SHA | d7d2244cb4cb95b94e28569ac8ef090d49876c89 |
| Issue | cnos#524 |
| Mode | MCA — execute known W1→W4 sequence; W0 is locked |
| AC set | AC1–AC7 (all) |
| W0 non-goals | No rendered-workflow output change; no runtime-behavior change; no new wake names; no schema changes outside schemas/skill.cue |
| Byte-identity oracle | cmp render(SKILL) render(JSON) must pass at every step before deletion |
| Deletion gate | W4 deletions ONLY after byte-identical proof CI-green |

## Implementation order (W1 → W4; minimizes inter-step risk)

### W1 — Schema + both SKILL.md files (AC1 + AC2)

**Step 1: Extend `schemas/skill.cue` (AC1)**
- Add `"wake"` to the `artifact_class` enum in `#Skill` (line 35).
- Add `#Wake` definition after `#Skill`. The `#Wake` definition must type-check the `wake:` block per the W0 design (see γ-scaffold §4 AC1 field table and §5 for the complete W0 wake: block shape).
- Key design constraints for `#Wake`:
  - `wake.role: "admin" | "dispatch"` (not observer — the renderer's enum check at line 342 of cn-install-wake accepts admin/dispatch/observer but the W0 schema is admin|dispatch only)
  - Role-shaped output disjunction: if `wake.role == "admin"` then `wake.output.channel_log_convention: string` required; if `wake.role == "dispatch"` then `wake.output.{cycle_artifact_root, artifact_class_taxonomy, cell_runtime}` required.
  - `wake.selector` is required when `role == "dispatch"` (sub-fields: `include: [...string]`, `exclude: [...string]`).
  - `wake.protocol` is optional string (required for dispatch semantically but CUE enforces it as optional to not break admin shapes).
  - Keep `#Wake` open (`...`) per the LANGUAGE-SPEC §11 principle from `#Skill`.

**Step 2: Author `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` (AC2)**

Frontmatter:
```yaml
---
name: agent-admin
description: "<from wake-provider.json description field>"
governing_question: "How does the agent-admin wake orchestrate the admin role for a given agent installation?"
triggers:
  - "cn install-wake agent-admin --agent <agent>"
  - "schedule (cron; substrate-encoded by renderer)"
  - "issues:opened with title matching issues_opened_title_pattern"
scope: global
artifact_class: wake
kata_surface: none
inputs:
  - "home thread (.cn-{agent}/threads/activations/{this-hub}/)"
  - "open issues in this repo (read-only)"
outputs:
  - ".cn-{agent}/logs/YYYYMMDD.md (channel log entry)"
  - ".cn-{agent}/state/ (cursor advancement)"
wake:
  role: admin
  package: cnos.core
  admin_only: true
  activation_log_writer: true
  input:
    triggers:
      - schedule
      - issues_opened_title_match
    issues_opened_title_pattern: claude-wake
  output:
    channel_log_convention: docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
  permission_intent:
    - contents.write
    - issues.write
    - pull_requests.write
    - id_token.write
  concurrency:
    serialize: true
    group: agent-admin-{agent}
  agent_variable:
    name: agent
    default: null
---
```

Body: the VERBATIM content of `src/packages/cnos.core/orchestrators/agent-admin/prompt.md`, followed (optionally, in a trailing `## Reference` section) by the verbose data from wake-provider.json that moves to body per W0 call #5: `responsibilities`, `*_notes`, `cross_references`, `allowed_surfaces`, `disallowed_surfaces`, `defer_path`, `relationship_to_substrate`.

**Step 3: Author `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (AC2)**

Frontmatter:
```yaml
---
name: cds-dispatch
description: "<from wake-provider.json description field>"
governing_question: "How does the cds-dispatch wake claim and route CDS cells to the δ role contract?"
triggers:
  - "cn install-wake cds-dispatch --agent <agent>"
  - "schedule (cron; substrate-encoded by renderer)"
  - "issues:labeled with label in selector.include"
scope: global
artifact_class: wake
kata_surface: none
inputs:
  - "open issues matching selector (dispatch:cell + protocol:cds + status:todo, minus exclude labels)"
outputs:
  - ".cdd/unreleased/{N}/ (cell artifact root via δ/γ/α/β)"
  - "issue lifecycle labels on claimed cell"
  - "pull request for claimed cell"
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
    cycle_artifact_root: .cdd/unreleased/{N}/
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
    group: cds-dispatch-{agent}
  agent_variable:
    name: agent
    default: sigma
---
```

Body: VERBATIM content of `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md`, followed by the verbose cross_references, responsibilities, *_notes, allowed/disallowed surfaces, and defer_path from wake-provider.json.

**W1 verification (before proceeding to W2):**
```sh
cue vet -d '#Skill' schemas/skill.cue src/packages/cnos.core/orchestrators/agent-admin/SKILL.md
cue vet -d '#Skill' schemas/skill.cue src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
./scripts/ci/validate-skill-frontmatter.sh 2>&1 | grep 'SKILL.md validated'  # should show 93
./scripts/ci/validate-skill-frontmatter.sh --self-test  # exit 0
# Goldens UNCHANGED — renderer still reads JSON:
cn install-wake agent-admin   # no diff to golden
cn install-wake cds-dispatch  # no diff to golden
```

### W2 — Teach renderer to read SKILL.md; dual-source parity gate (AC3 + AC4)

**Step 4: Extend `cn-install-wake` to parse SKILL.md**

The renderer currently reads fields from `wake-provider.json` via `jq`. For W2, add a SKILL.md parsing path:
- Detect `SKILL.md` beside the manifest path.
- Extract the YAML frontmatter (between the two `---` delimiters).
- Map `wake.*` frontmatter keys to the variables the renderer currently reads from JSON:
  - `name` ← `name` (top-level in both)
  - `package` ← `wake.package`
  - `role` ← `wake.role`
  - `admin_only` ← `wake.admin_only`
  - `activation_log_writer` ← `wake.activation_log_writer`
  - `activation_state` ← `wake.activation_state` (default: `"live"` when absent)
  - `protocol` ← `wake.protocol`
  - `selector` ← `wake.selector`
  - `input_contract.triggers` ← `wake.input.triggers`
  - `input_contract.issues_opened_title_pattern` ← `wake.input.issues_opened_title_pattern`
  - `output_contract.*` ← `wake.output.*`
  - `permission_intent` ← `wake.permission_intent`
  - `concurrency_intent.serialize` ← `wake.concurrency.serialize`
  - `concurrency_intent.group` ← `wake.concurrency.group`
  - `agent_variable.default` ← `wake.agent_variable.default`
  - **prompt body** ← SKILL.md body (everything after the second `---`)
- Expose this as `--source skill-md` flag (or auto-detect: if SKILL.md present beside manifest, available as alternative source). Keep `--source json` (or default) path unchanged.

**Step 5: Add dual-source parity gate to `install-wake-golden.yml` (AC4)**

Add a new CI step after the existing "Re-render" steps:
```yaml
- name: W2 parity gate — render(SKILL) == render(JSON) for both wakes
  run: |
    set -euo pipefail
    for wake in agent-admin cds-dispatch; do
      ./src/packages/cnos.core/commands/install-wake/cn-install-wake "$wake" \
        --source skill-md --out /tmp/${wake}-from-skill.yml
      ./src/packages/cnos.core/commands/install-wake/cn-install-wake "$wake" \
        --source json     --out /tmp/${wake}-from-json.yml
      if ! cmp -s /tmp/${wake}-from-skill.yml /tmp/${wake}-from-json.yml; then
        echo "::error::W2 parity FAILED for $wake: render(SKILL) != render(JSON)"
        diff /tmp/${wake}-from-skill.yml /tmp/${wake}-from-json.yml || true
        exit 1
      fi
      echo "W2 parity OK ($wake): render(SKILL) == render(JSON)"
    done
```

### W3 — Flip renderer source to SKILL.md; re-render (AC3 + AC4)

**Step 6: Flip default source in `cn-install-wake`**

Change the resolver so that when SKILL.md is present, it is the primary source (not the JSON). The JSON fallback path remains for backward compatibility until W4.

**Step 7: Re-render goldens and live workflows**

```sh
cn install-wake agent-admin   --out src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
cn install-wake cds-dispatch  --out src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
cn install-wake agent-admin   --out .github/workflows/cnos-agent-admin.yml
cn install-wake cds-dispatch  --out .github/workflows/cnos-cds-dispatch.yml
```

These MUST produce zero diff if byte-identical — the oracle passes. If any diff appears: stop, file an FN, do NOT proceed to W4.

### W4 — Delete JSON + prompt; I5 count confirms 93 (AC5 + AC7)

**Step 8: Confirm byte-identity in CI, then delete source files**

Only after W3 re-renders are verified byte-identical in CI:
```sh
rm src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json
rm src/packages/cnos.core/orchestrators/agent-admin/prompt.md
rm src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json
rm src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md
```

**Step 9: Update renderer to remove JSON fallback path**

After deletion, remove the JSON parsing path from `cn-install-wake` (the dead code for `wake-provider.json` lookup). The manifest path `--manifest <path>` flag still works but resolves against SKILL.md instead.

**Step 10: Update `install-wake-golden.yml` for W4**

- Remove or update the `--manifest` path comments that reference `wake-provider.json`.
- Add an I5 count assertion step:
  ```yaml
  - name: I5 count check (93 SKILL.md post-W4)
    run: |
      n=$(find src/packages -name SKILL.md | wc -l)
      if [ "$n" -ne 93 ]; then
        echo "::error::Expected 93 SKILL.md files post-W4; found $n"
        exit 1
      fi
      echo "I5 count OK: $n SKILL.md files"
  ```
- Add a W4 deletion-proof step confirming JSON + prompt absent for both wakes.
- Verify `cn install-wake` still renders byte-identical goldens.

## Commit message pattern

```
α-524 W1: <step name summary>
α-524 W2: <step name summary>
α-524 W3: <step name summary>
α-524 W4: <step name summary>
α-524 R0 signal: ready for β review
```

## Signaling β-readiness

After all steps land + push, write `.cdd/unreleased/524/self-coherence.md` §R0, commit as `α-524 R0: self-coherence §R0`, push. Then add a signal commit `α-524 R0 signal: ready for β review` (zero file changes).

## Scope guardrails (hard constraints)

- **NO rendered-workflow output change** unless caused by a genuine field delta between JSON and SKILL.md (any such delta is an R0 blocker — stop, file FN, do not proceed to W3).
- **NO runtime-behavior change** — the renderer's refusal logic (exits 2/3/4), the write-fence step, the permission encoding, the cron slots, the bot identity bindings — none of these change in this migration.
- **NO renderer substrate-encoding change** — the YAML structure the renderer emits is unchanged; only the SOURCE it reads from changes.
- **NO new wake names** — this migration covers agent-admin and cds-dispatch exactly.
- **NO schema changes outside `schemas/skill.cue`** — do not modify I5 validator logic.
- **NO I5 exception-list additions** — the new SKILL.mds must pass cue vet without exception entries.
- **Byte-identity at every step** — if `cmp render(SKILL) render(JSON)` fails at any point in W2 or later, stop and file an FN before continuing.
- **Deletion is W4-only** — do NOT delete JSON or prompt until CI confirms byte-identical goldens from SKILL.md source.

## Re-iteration discipline

If β R0 returns iterate:
- Read β's review; do NOT re-derive ACs from scratch.
- R[N] section in self-coherence.md per T-486-7 off-by-one prevention.
- Apply δ-direct R1 pattern if the iterate is narrow + mechanical (single file/line fix without architectural reframing).

## Hard do-not-cross lines

- Do NOT pre-empt or modify any cycle outside cycle/524's scope.
- Do NOT modify the test-fixtures at `src/packages/cnos.core/commands/install-wake/test-fixtures/` unless a new mis-declared SKILL.md fixture is needed for AC6 verification (in which case: file an FN naming the new fixture and why it's needed).
- Do NOT touch `schemas/skill-exceptions.json` — new SKILL.mds must be clean.
- Do NOT modify `.cn-sigma/logs/` — this cycle has no activation-log write scope.
```

## §7. β dispatch prompt

```text
You are β@cdd.cnos, running as the reviewer for cycle/524 R0.

Branch: cycle/524 (read α's R0 commits up to and including the "α-524 R0 signal" commit).
Parent issue: cnos#524 ("wake-as-skill: migrate wakes to typed SKILL.md modules").
Scaffold: .cdd/unreleased/524/gamma-scaffold.md (read in full).
α self-coherence: .cdd/unreleased/524/self-coherence.md §R0 (read in full).

## Your output

Write `.cdd/unreleased/524/beta-review.md` with §R0 review. Final verdict: **converge** OR **iterate**.

## Pre-AC checklist (mechanical; walk every item in order)

**AC1 — `#Wake` schema in CUE:**
- `schemas/skill.cue` adds `"wake"` to `artifact_class` enum.
- `#Wake` definition present with correct field types/enums per W0 design.
- Run: `cue vet -d '#Wake' schemas/skill.cue /tmp/conformant-wake.yaml` → exit 0.
- Run a bad-enum fixture → exit nonzero + diagnostic. Record both results.

**AC2 — Wake SKILL.md modules authored:**
- Both SKILL.md files exist at the paths in §2 Source-of-truth table.
- Frontmatter carries `scope: global`, `artifact_class: wake`, `wake:` block with exact manifest data.
- Body contains prompt.md content verbatim (no paraphrase; no summary; literal copy).
- Verbose fields (`responsibilities`, `*_notes`, `cross_references`) are in body ONLY, not frontmatter.
- `./scripts/ci/validate-skill-frontmatter.sh` exits 0 and reports 93 SKILL.md.

**AC3 — Renderer reads SKILL.md:**
- `cn-install-wake` has a SKILL.md parsing path.
- The field mapping from `wake.*` to renderer variables is complete (no missing fields; see γ-scaffold §4 AC1 table).
- After W3 flip: `cn install-wake agent-admin` uses SKILL.md as primary source.

**AC4 — Byte-identical goldens:**
- `cmp render(SKILL) render(JSON)` passes for BOTH wakes (record sha256 of each pair).
- sha256(golden) == sha256(live workflow) for both wakes.
- Second render is a no-op (idempotence: sha256 before == sha256 after).
- The W2 parity CI step is present in `install-wake-golden.yml` and passes.

**AC5 — JSON + prompt deleted:**
- All four deletion-target files are absent in the W4 commit.
- `cn install-wake` still renders byte-identical goldens post-deletion.
- `install-wake-golden.yml` W4 deletion-proof step present and passes.

**AC6 — Renderer refusals preserved:**
- Existing cycle/496 refusal smoke in `install-wake-golden.yml` still passes: exit 4 + "activation_log_writer mis-declaration" in stderr; exit 3 on declaration-only fixture.
- Refusal logic fires correctly whether source is JSON (W2) or SKILL.md (W3/W4).
- If α created a synthetic mis-declared SKILL.md fixture: verify it exits 4.

**AC7 — All gates green:**
- No inherited-cap failures; no new red gates on the PR.
- I5 count step in CI reports 93.

## Independence requirement

β must reach its verdict independently. Do NOT ask α to pre-certify; do NOT accept α's self-assessment as a substitute for β's own verification. Run the oracles yourself (or simulate them against the committed state).

## Variable consistency walk (T-486-1 + T-487-1)

Walk every variable that changed across: SKILL.md frontmatter → renderer field consumption → rendered workflow YAML → CI oracle → prompt body → PR body.

Key variables to check:
| Variable | SKILL.md `wake:` block | Renderer parsing | Rendered YAML | CI assertion | Body prose |
|---|---|---|---|---|---|
| `wake.role` | `admin` / `dispatch` | line ~342 enum check | implicit in shape | AC1 CUE vet | body mentions role |
| `wake.activation_log_writer` | `true` / `false` | line ~452 mis-declaration check | write-fence step presence | AC6 refusal smoke | body names "non-writer" |
| `wake.concurrency.group` | `agent-admin-{agent}` / `cds-dispatch-{agent}` | line ~498 substitution | `group:` in YAML | golden shape check | body mentions concurrency |
| prompt body | verbatim in SKILL.md body | body extraction after `---` | `prompt: \|` block | sha256 parity | N/A (IS the body) |

Flag any drift between columns.

## Per-CI-step bash-e audit (cnos#478 mechanical-injection)

For every new `run:` step α added to `install-wake-golden.yml`:
- Verify `set -eu` / `set -euo pipefail` is in effect.
- Verify no `grep -c ... | true` pipefail-trap class.
- Verify fixture cleanup (no orphaned temp files after step).

## Byte-identity audit (load-bearing)

1. Record sha256 of both goldens from JSON source (pre-W3 baseline).
2. Record sha256 of both goldens from SKILL.md source (W2 parity render).
3. Verify all four sha256s match.
4. Record sha256 of live workflows; verify they match goldens.
5. After W4 deletion: re-verify goldens unchanged.

If any pair mismatches: P0 blocker. Iterate immediately. Byte-identity is the migration oracle.

## Output structure for `.cdd/unreleased/524/beta-review.md`

- **§R0 frontmatter.** YAML with cycle, base SHA, α R0 signal commit SHA, β review timestamp.
- **§R0 verdict.** `converge` OR `iterate` at the top, before narrative.
- **§R0 AC table.** AC1–AC7 with per-AC green/red/ambiguous + oracle command + observation.
- **§R0 variable consistency walk.** Every variable × every column.
- **§R0 per-CI-step bash-e audit.** Every modified/added `run:` step + result.
- **§R0 byte-identity audit.** sha256 pairs for all four files.
- **§R0 friction notes (if any).** FN-β-N entries.

## When to iterate vs converge

- **Iterate** if: any byte-identity pair mismatches; any AC oracle returns red; SKILL.md body is NOT verbatim prompt.md (paraphrase = red); renderer skips a field from the `wake:` block; I5 count != 93; any existing refusal smoke breaks.
- **Converge** if: byte-identical across all pairs; all AC oracles green; I5 count = 93; existing CI smokes all pass; variable consistency clean.

## Hard constraints

- **Do NOT iterate for taste** — β iterates only when an oracle returns red or a load-bearing invariant is violated.
- **Do NOT modify any α-written file** — β reviews; does not implement.
- **Honor T-486-12 (operator-final-read P1)** — β converge is not the cycle's final gate.
- **Commit message:** `β-524 R0 review: <converge|iterate>`. Co-Authored-By trailer matches prior cycles.
```

## §8. Scope guardrails (W0 non-goals)

These are hard guardrails α and β must both enforce:

1. **No rendered-workflow output change.** The YAML emitted by `cn install-wake` for both wakes must be byte-identical before and after the migration. Any diff is an R0 blocker.
2. **No runtime-behavior change.** The renderer's refusals (exits 2/3/4), the write-fence step logic, the permission block, the cron-slot encoding, the bot identity table — none change. This migration is a source-swap, not a behavior change.
3. **No renderer substrate-encoding change.** The `gha_permission_line` function, the cron-slot list (`8 23 38 53`), the `anthropics/claude-code-action@v1` pinning, the `SIGMA_WORKFLOW_PAT` secret name, the `cancel-in-progress: false` policy — all unchanged.
4. **No new wake names.** Migration scope is exactly: `agent-admin` (cnos.core) and `cds-dispatch` (cnos.cds).
5. **No other schema changes.** `schemas/skill.cue` changes are limited to: adding `"wake"` to the `artifact_class` enum + adding the `#Wake` definition. No other field or constraint changes.

## §9. Friction notes (anticipated classes for α + β)

**FN-1: Frontmatter YAML quoting.** The wake-provider.json contains string values with characters that may require quoting in YAML frontmatter (colons in descriptions, curly braces in `{agent}` substitution tokens). `{agent}` tokens in strings like `"agent-admin-{agent}"` must be quoted in YAML to avoid being parsed as flow mappings. α must verify: `cue vet` passes with the actual frontmatter text, not a simplified version.

**FN-2: Prompt body verbatim discipline.** The SKILL.md body must contain prompt.md content verbatim — not paraphrased, not reformatted, not summarized. The byte-identity oracle (AC4) is the enforcement mechanism: if the renderer's body-extraction produces a different string than the raw `prompt.md`, the rendered output will differ. Any transform applied to the body (e.g., indentation normalization, line-ending normalization) will break byte-identity. The renderer currently uses `sed "s/{agent}/${agent}/g" "$prompt_path"` — the SKILL.md body extraction path must apply the same substitution and the same trailing-whitespace strip (`sed 's/[[:space:]]*$//'`) to produce identical output.

**FN-3: The renderer's `--manifest` flag interaction.** `cn-install-wake` accepts `--manifest <path>` to override the default JSON lookup. The manifest-override path resolves `prompt_path` relative to the manifest's directory (`manifest_dir`). When the source is SKILL.md, the body extraction replaces both: the `prompt_template` JSON field lookup AND the `$manifest_dir/$prompt_template_relpath` file read. The W2 SKILL.md path must handle `--manifest` correctly — if `--manifest` points to a SKILL.md (or if there's a SKILL.md beside a `--manifest`-specified JSON), the precedence is well-defined. α should document the flag interaction in the renderer comments.

**FN-4: Renderer field name mapping (`input_contract` vs `wake.input`).** The JSON manifest uses `input_contract: {triggers: [...]}` while the W0 SKILL.md shape uses `wake: {input: {triggers: [...]}}`. The renderer currently reads `jq -r '.input_contract.triggers'`. The SKILL.md parser must read `wake.input.triggers` and expose it as `input_contract_triggers` (or however the renderer variable is named) to avoid duplication of the rendering logic. α must map ALL fields — including the less-obvious ones like `input_contract.issues_opened_title_pattern` ← `wake.input.issues_opened_title_pattern` and `output_contract.channel_log_convention` ← `wake.output.channel_log_convention` (admin shape) and `output_contract.{cycle_artifact_root,artifact_class_taxonomy,cell_runtime}` ← `wake.output.{...}` (dispatch shape).

**FN-5: `cross_references.consumed_skills` FN-6 defense-in-depth check.** The renderer currently (at line 477-484) checks `cross_references.consumed_skills` from the JSON for the FN-6 attach-incompatibility refusal. After W3/W4, this data moves to the SKILL.md body (per W0 call #5: `cross_references` → body). The renderer must extract this from the SKILL.md body OR the `calls` / `requires` frontmatter field. Design choice: either (a) keep `cross_references.consumed_skills` in the SKILL.md frontmatter as a `requires:` list (LANGUAGE-SPEC §11 allows this), or (b) move the FN-6 check to a body-prose check (fragile). Option (a) is preferred. α must document this choice as a design fork and record the resolution in self-coherence.md.

**FN-6: `activation_state_notes` field.** The renderer reads `activation_state_notes` from the JSON (line 387) for the exit-3 refusal message. The W0 `wake:` block shape does not include this field. α must decide: add `wake.activation_state_notes` as an optional string to `#Wake`, OR absorb the notes into the body. If absorbed into body, the renderer's exit-3 message will not include the notes content. This is a narrow behavioral change (refusal message text only; not the refusal itself). α should file an FN if the notes are dropped from the exit-3 message.

**FN-7: I5 discovery path — are orchestrator/ dirs under `src/packages/`?** The I5 validator runs `find src/packages -name SKILL.md`. Both orchestrator dirs (`src/packages/cnos.core/orchestrators/agent-admin/` and `src/packages/cnos.cds/orchestrators/cds-dispatch/`) ARE under `src/packages/`. Discovery is automatic — confirmed. However, α should verify the `calls:` resolution in `validate_skill_file` for orchestrator SKILL.mds: `package_skill_root_of` uses the path's `.../packages/<pkg>/skills/` pattern to find the root. Orchestrator SKILL.mds are under `.../packages/<pkg>/orchestrators/`, not `.../packages/<pkg>/skills/`. The `package_skill_root_of` function falls back to `dirname "$file"` for out-of-pattern paths. This means `calls:` entries in orchestrator SKILL.mds resolve relative to the orchestrator dir, not the package skill root. If the new SKILL.mds use `calls:`, α must set them to avoid the path-resolution trap. Safest: use `requires:` instead of `calls:` for the cross_references.consumed_skills data.

**FN-8: `cue export` vs `cue vet` for `#Wake`.** The I5 validator uses `cue export --out json -d '#Skill'` (line 176). When a SKILL.md has `artifact_class: wake`, the `cue export` call passes `-d '#Skill'`, not `-d '#Wake'`. The `#Skill` definition must include `artifact_class: "wake"` in its enum for the export to succeed. The `#Wake` definition is separate and used for the wake-specific type-check. These are two separate CUE operations: `cue vet -d '#Skill'` (general SKILL.md validation by I5) and `cue vet -d '#Wake'` (AC1 oracle). Both must pass for a wake SKILL.md. This is the correct architecture — the `#Skill` definition is the I5 gate; `#Wake` is the wake-contract gate. α must not conflate them.

## §10. Cross-references + carryforwards

| Carryforward | Source | How applied this cycle |
|---|---|---|
| T-486-1 (variable consistency table) | cycle/486 closeouts | α writes seed table in self-coherence.md; β walks every variable × every surface |
| T-486-7 (R[N] off-by-one prevention) | cycle/486 closeouts | α's iteration discipline cited in §6 α prompt |
| T-486-12 (operator-final-read defense-in-depth, P1) | cycle/486 + cycle/487 promotion | β explicitly defers "final gate" to operator-final-read |
| T-487-1 (variable consistency table extended to ALL surfaces) | cycle/487 closeouts | β prompt §"Variable consistency walk" covers 5 columns including prompt body |
| Byte-identity oracle | cycle/487 Sub 3 (cycle/476 pattern) | The migration oracle at every step; stricter than cycle/476 (which only verified within a single source) |
| Wake-provider contract (AC4 refusal patterns) | cycle/496 | Renderer refusals PRESERVED post-migration; AC6 confirms they fire from SKILL.md source |
| AC8/AC7 renderer-side authority audit | cycle/476 / cycle/485 | Not explicitly repeated here; the audit constraint is inherited: renderer must not encode role decisions. The migration does not change role decision encoding — it changes source format only. |

## §11. Non-goals + out-of-scope

- No new wake providers beyond agent-admin and cds-dispatch.
- No changes to the rendered workflow structure (no new steps, no changed triggers, no permission changes).
- No changes to the `cn wake install` operator-facing subcommand contract.
- No migration of test-fixtures at `src/packages/cnos.core/commands/install-wake/test-fixtures/` to SKILL.md format (those fixtures are JSON-native synthetic cases; migrating them is a separate cycle if ever needed).
- No schema changes to `wake-provider.json` format (the JSON is being replaced, not evolved).
- No other orchestrator directories (future packages: cnos.cdr, cnos.cdw; those are out of scope for this cycle).
- No admin dispatch-summary work (cnos#495 Sub 2; HELD).
- No artifact-root rename (cnos#497; design-held).

## §12. Cycle-shape forecast

- **R0 converge target.** The migration is mechanical: author SKILL.md (body = verbatim prompt), extend CUE, teach renderer, verify byte-identity, delete source. Four phases but each phase has a clear oracle. R0 converge is achievable if α is rigorous about byte-identity at each step and FN-2 (verbatim body) + FN-4 (field mapping completeness) are handled cleanly.
- **Most likely R0 iterate trigger.** FN-2 (prompt body not verbatim → byte-identity fails) or FN-4 (field mapping gap → renderer skips a field → wrong YAML output). Both are detectable by the W2 parity gate before proceeding to W3.
- **Operator-final-read expectation.** Operator is likely to verify: (a) SKILL.md body is actually verbatim — not summarized; (b) `#Wake` CUE definition correctly types the role-output disjunction; (c) the W4 deletion proof CI step is load-bearing, not ceremonial.
- **PR title forecast.** `cycle/524: wake-as-skill migration W1→W4 — typed SKILL.md modules replace JSON+prompt pair`.

---

Filed by γ@cdd.cnos (wake-invoked δ dispatch, cycle/524), 2026-06-30 (UTC).
