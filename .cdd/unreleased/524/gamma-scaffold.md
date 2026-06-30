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
dispatch_scope: W0-design-only
scope_override_reason: operator-directive-2026-06-30
ac_set: W0 design document (not issue AC1–AC7; those are W1→W4 scope)
w0_locked: true
---

# γ-scaffold — cnos#524: wake-as-skill W0 design document

## §1. Scope

**This dispatch is W0 design only.** The W1→W4 build ACs (AC1–AC7) in the issue body are explicitly NOT in scope for this run. The operator DISPATCH DIRECTIVE (issue #524 comment "⛳ DISPATCH DIRECTIVE") governs this cycle.

**Scope override reason:** The operator directive restricts this dispatch to W0 design work:
> "This is W0 design only. Do not start W1. Do not change renderer code. Do not change CUE schema. Do not change workflows. Do not change goldens. Do not delete wake-provider.json or prompt.md."

**What this cycle delivers:**
- `.cdd/unreleased/524/w0-design.md` — a standalone W0 design document formalizing the wake-as-skill contract
- Normal CDD closeout records (`self-coherence.md`, `beta-review.md`, closeouts)
- A PR containing ONLY the design artifact and CDD artifacts (no source/schema/workflow changes)

**What this cycle explicitly does NOT deliver:**
- Any changes to `schemas/skill.cue`
- Any changes to `src/.../cn-install-wake`
- Any `.github/workflows/*.yml` changes
- Any `*.golden.yml` changes
- Deletion or editing of `wake-provider.json` or `prompt.md`

## §2. W0 design context (what α will formalize)

The W0 design is already complete and operator-ratified. It lives in two places:
1. **Issue body** (`cnos#524 §"W0 design (locked)")** — the locked design block with the `wake:` block shape, CUE↔renderer boundary, and W1→W4 migration plan
2. **Operator W0 calls comment** (issue #524, first comment by @usurobor) — the 5 design decisions

The W0 design document α authors is a FORMAL EXTRACTION and PRESENTATION of this design, not a new design. α reads the issue body and the operator W0 calls comment, then structures them into a design document that:
- Can be read and reviewed independently of the issue
- Names each design decision explicitly with its rationale
- Defines the schema precisely enough for a future W1 implementer
- Is internally coherent (no drift between the schema shape and the migration plan)

## §3. Source-of-truth table

| Surface | Path | Role |
|---|---|---|
| Issue body (W0 design section) | GitHub issue #524 body | Primary source; contains locked W0 design |
| Operator W0 calls comment | First @usurobor comment on #524 | The 5 operator decisions |
| w0-design.md (to create) | `.cdd/unreleased/524/w0-design.md` | α's output artifact |
| Current wake manifests (read-only) | `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`<br>`src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | Reference: what the `wake:` block must capture |
| Current skill CUE schema (read-only) | `schemas/skill.cue` | Reference: what `#Skill` looks like; where `#Wake` will extend it |

## §4. W0 design document structure (α deliverable)

The `w0-design.md` α produces MUST contain these sections in order:

### §A. Problem statement and invariant
- Two-sentence problem: two contract systems, should be one
- The invariant: "A wake is a typed SKILL.md module whose body is the prompt and whose frontmatter compiles to workflow substrate"

### §B. The `wake:` block schema (the `#Wake` CUE definition)
- Complete field listing with types, enums, optionality, and notes
- Source: issue body `## W0 design → wake: block` + operator W0 calls comment
- Must include ALL fields the renderer currently reads from `wake-provider.json` (cross-referenced against the source manifests)
- Must show the role-shaped output disjunction (admin vs dispatch shapes)

### §C. The 5 operator design decisions
1. `wake:` nesting (single block, not flattened)
2. `artifact_class: wake` (reuse SKILL.md discovery + CUE pipeline)
3. `scope: global` (role in `wake.role`)
4. CUE owns typed contract; renderer owns substrate compilation + runtime refusals
5. `responsibilities`/`*_notes`/`cross_references` → body, not frontmatter

Each decision: state the decision, the rationale, and the implication for W1 implementation.

### §D. CUE ↔ renderer responsibility boundary
Concrete table:
| Boundary | CUE owns | Renderer owns |
- Static shape, field types, enums, role-output disjunction → CUE
- Cross-field refusals (activation_log_writer mis-declaration, activation_state gating) → renderer
- Substrate encoding (GitHub Actions YAML, secrets, cron slots, bot identity table) → renderer
- Body extraction (second `---` delimiter; prompt substitution) → renderer
- Install-time substitution (`{agent}` → `$agent`) → renderer

### §E. Body-as-prompt rule
- Frontmatter = contract (machine-readable fields only)
- Body = prompt/role doctrine (verbatim prompt.md content + verbose reference data from wake-provider.json)
- The byte-identity oracle proves correctness: renderer extracting body → same output as reading prompt.md

### §F. W1→W4 migration plan (the 4 phases)
- W1: Author both SKILL.md beside JSON/prompt; add `#Wake` to skill.cue; goldens unchanged
- W2: Teach renderer to read SKILL.md; dual-source parity gate
- W3: Flip renderer source to SKILL.md; re-render (no-op expected)
- W4: Delete JSON + prompt after byte-identical proof; I5 count 91 → 93
- Byte-identity oracle at each step

### §G. W1 acceptance criteria (scoped)
The 7 ACs from the issue body, restated as the acceptance criteria for W1 specifically:
- AC1: `#Wake` schema in `schemas/skill.cue`
- AC2: Both wake SKILL.md files authored at the orchestrator paths
- AC3: Renderer reads SKILL.md
- AC4: Byte-identical goldens
- AC5: JSON + prompt deleted
- AC6: Renderer refusals preserved
- AC7: All CI gates green

### §H. Friction notes carried forward (for W1 implementer)
The 8 friction notes (FN-1 through FN-8) from the prior γ scaffold are material for the W1 implementer. α includes the key ones here.

### §I. Non-goals
Explicit non-goals: no rendered-workflow output change; no runtime-behavior change; no new wake names; no Demo 0.

## §5. α dispatch prompt

```text
You are α@cdd.cnos, running as the implementer for cycle/524 R0.

Branch: cycle/524 (at base main SHA d7d2244cb4cb95b94e28569ac8ef090d49876c89).
Parent issue: cnos#524 ("wake-as-skill: migrate wakes to typed SKILL.md modules").
Scaffold: .cdd/unreleased/524/gamma-scaffold.md (READ THIS FIRST in full before any work).

## Scope of this run (MANDATORY — read before any file write)

**This run is W0 DESIGN ONLY.**

The operator DISPATCH DIRECTIVE (see issue #524 comments — "⛳ DISPATCH DIRECTIVE") overrides the full issue scope (W1→W4 ACs) for this run. Your task is to produce a standalone W0 design document — NOT to implement W1→W4.

**Hard scope constraints (enforced by directive):**
- Do NOT touch `schemas/skill.cue`
- Do NOT touch `src/packages/.../cn-install-wake`
- Do NOT touch any `.github/workflows/*.yml`
- Do NOT touch any `*.golden.yml`
- Do NOT delete or edit `wake-provider.json` or `prompt.md`
- Do NOT start W1, W2, W3, or W4 implementation

**Your ONLY writable outputs:**
- `.cdd/unreleased/524/w0-design.md` — the W0 design document (see structure in γ-scaffold §4)
- `.cdd/unreleased/524/self-coherence.md §R0` — your self-coherence artifact
- Commits on `cycle/524` branch

## What you produce

Author `.cdd/unreleased/524/w0-design.md` per the structure in γ-scaffold §4 (§A through §I).

This document FORMALIZES the W0 design already in the issue body. It does NOT design anything new. Read:
1. Issue #524 body (especially `## W0 design (locked)` section)
2. The operator W0 calls comment (first @usurobor comment on #524)
3. `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — REFERENCE ONLY (what `wake:` block must capture for admin wake)
4. `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — REFERENCE ONLY (what `wake:` block must capture for dispatch wake)
5. `schemas/skill.cue` — REFERENCE ONLY (what `#Skill` looks like; where `#Wake` fits)

Do NOT modify any of the reference files above.

## Skills to load before implementing

1. `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract
2. `src/packages/cnos.cdd/skills/cdd/SKILL.md` — cell-runtime framework

## Self-coherence (after w0-design.md is complete)

Write `.cdd/unreleased/524/self-coherence.md §R0`:
- Confirm scope guardrails: which files were NOT touched
- Walk the §4 structure sections (§A through §I): each section present? Internal coherence of `wake:` block fields?
- Cross-reference check: does `w0-design.md §B` cover ALL fields from both `wake-provider.json` manifests?
- Confirm: no W1→W4 implementation artifacts exist on the branch

## Commits

1. Commit `.cdd/unreleased/524/w0-design.md` as: `α-524 W0: w0 design document`
2. Commit `.cdd/unreleased/524/self-coherence.md` as: `α-524 R0: self-coherence §R0`
3. Push all commits to `cycle/524`
```

## §6. β dispatch prompt

```text
You are β@cdd.cnos, running as the reviewer for cycle/524 R0.

Branch: cycle/524 (read α's R0 commits).
Parent issue: cnos#524 ("wake-as-skill: migrate wakes to typed SKILL.md modules").
Scaffold: .cdd/unreleased/524/gamma-scaffold.md (read in full).
α self-coherence: .cdd/unreleased/524/self-coherence.md §R0 (read in full).

## Scope reminder

This run is W0 design only. The ONLY artifact α should have produced is `.cdd/unreleased/524/w0-design.md` plus CDD records. If α touched any source/schema/workflow files, that is an immediate R0 iterate finding.

## Review checklist

**Scope compliance (check first — iterate immediately on any failure):**
- [ ] No changes to `schemas/skill.cue`
- [ ] No changes to `cn-install-wake` or any renderer file
- [ ] No changes to `.github/workflows/*.yml`
- [ ] No changes to `*.golden.yml`
- [ ] `wake-provider.json` and `prompt.md` untouched for both wakes
- [ ] No W1→W4 implementation code on the branch

**`w0-design.md` completeness (walk each section):**
- [ ] §A: Problem statement and invariant — accurate?
- [ ] §B: `wake:` block schema — covers ALL fields from BOTH `wake-provider.json` manifests? Role-shaped output disjunction present?
- [ ] §C: All 5 operator decisions present with decision + rationale + W1 implication?
- [ ] §D: CUE↔renderer boundary table — complete? No mis-assignments?
- [ ] §E: Body-as-prompt rule — clear? Byte-identity oracle cited?
- [ ] §F: W1→W4 migration plan — 4 phases with oracle at each step?
- [ ] §G: W1 ACs (AC1–AC7) present and accurate?
- [ ] §H: Friction notes cited or included?
- [ ] §I: Non-goals listed?

**Field completeness (load-bearing):**
Read both `wake-provider.json` files. Walk every field. Is each field covered in `w0-design.md §B`?
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`

**Internal coherence:**
- Role-shaped output disjunction in §B: admin shape ≠ dispatch shape; both present?
- §F migration plan consistent with §B schema shape?
- §G ACs consistent with §D CUE↔renderer boundary?

**Verdict:**
- **converge** if: scope compliance clean; all 9 sections present and complete; all JSON fields covered in §B; internal coherence holds
- **iterate** if: any scope compliance failure; any section absent; any load-bearing field missing from §B; internal incoherence

Write `.cdd/unreleased/524/beta-review.md §R0` with verdict.
Commit as: `β-524 R0 review: <converge|iterate>`
Push to `cycle/524`.
```

## §7. Scope guardrails (mechanical — operator directive)

The cell MUST NOT touch any of:
- `src/.../cn-install-wake` (renderer)
- `schemas/skill.cue`
- Any `.github/workflows/*.yml`
- Any `*.golden.yml`
- `wake-provider.json` or `prompt.md` (either wake)

The ONLY diff in the PR branch should be under `.cdd/unreleased/524/`.

## §8. Cross-references

- Issue #524 body: `## W0 design (locked)` — the locked W0 design
- Issue #524 comments: operator W0 calls (first @usurobor comment) — the 5 decisions
- Issue #524 comments: ⛳ DISPATCH DIRECTIVE — scope override governing this run
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — admin wake manifest (reference only)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — dispatch wake manifest (reference only)
- `schemas/skill.cue` — current skill schema (reference only)
- Prior γ scaffold (cycle/524 commit `f148a1fd`): FN-1 through FN-8 friction notes — material for W1 implementer

---

Filed by γ@cdd.cnos (wake-invoked δ dispatch, cycle/524), 2026-06-30 (UTC).
Scope override: operator DISPATCH DIRECTIVE (W0 design only; W1→W4 ACs not in scope for this run).
