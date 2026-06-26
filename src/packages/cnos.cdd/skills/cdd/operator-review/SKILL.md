---
name: operator-review
description: Schema skill for the cn.operator-review.v1 artifact. HI uses this schema to capture a human operator's verdict on a cell at status:review and translate it into a typed durable input for the cycle's mechanical re-entry path.
artifact_class: skill
governing_question: How does the HI record an operator verdict as a typed artifact that the mechanical runtime can act on without role-boundary crossing?
visibility: public
parent: cdd
triggers:
  - operator-review
  - operator verdict
  - iterate
  - converge
  - reject
  - clarify
  - review-return
  - degraded_recovery
scope: task-local
inputs:
  - human operator's verbal or written verdict on a cell at status:review
  - the cell's PR and issue
outputs:
  - a conforming cn.operator-review.v1 YAML-frontmatter + Markdown findings artifact
  - (on iterate/reject) input to cn cell return for the status:review → status:changes transition
  - (on iterate) input to cn cell resume for cycle re-entry
requires:
  - cell is at status:review
  - HI has read the PR and the operator's verdict
calls:
  - cnos.cdd/skills/cdd/delta/SKILL.md §9.6 (return token carve-out)
  - cnos.cdd/skills/cdd/alpha/SKILL.md §2.1 (α dispatch intake)
---

# Operator-review schema — cn.operator-review.v1

## Core principle

**The HI's role in the operator-iterate flow is: read the PR; translate the operator's verdict into a typed artifact; commit and push; trigger the runtime transition. Stop there.**

The HI may turn human intent into an artifact. It must not turn itself into the worker that artifact instructs.

This skill defines the `cn.operator-review.v1` schema, the authoring rules for the artifact, and the `degraded_recovery` declaration schema used when the HI applies an operator patch under bootstrap-exception conditions.

---

## 1. cn.operator-review.v1 schema

### 1.1. Required frontmatter fields

Every `cn.operator-review.v1` artifact MUST carry a YAML frontmatter block with the following fields:

| Field | Type | Description |
|---|---|---|
| `schema` | string | Fixed value: `cn.operator-review.v1` |
| `issue` | integer | The cell's GitHub issue number |
| `pr` | integer or null | The PR associated with the cell at review time; null if no PR exists |
| `verdict` | enum | One of: `converge` / `iterate` / `reject` / `clarify` |
| `reviewer` | string | Identity of the reviewing actor (typically `human-operator`) |
| `captured_by` | string | Identity of the HI session that authored this artifact (e.g. `gamma-interface (HI)`) |
| `captured_at` | string | UTC date (ISO 8601) when the artifact was authored |
| `findings_count` | integer | Count of findings listed in the body (0 on `converge`; ≥1 on `iterate`/`reject`) |

Optional fields:

| Field | Type | Description |
|---|---|---|
| `worker_pr_head_at_review` | string | Short SHA of the PR HEAD commit the reviewer read (improves audit trail) |
| `round` | integer | The R[N] round this review corresponds to (default 1 if omitted) |

### 1.2. Verdict enum values

| Verdict | Semantics | Runtime action |
|---|---|---|
| `converge` | Operator approves the cell; PR is mergeable | No label transition; operator merges PR; standard close path |
| `iterate` | Operator requests changes; cell must re-engage | `cn cell return` → `status:review → status:changes`; `cn cell resume` → cycle re-armed |
| `reject` | Operator rejects the cell outright; cycle terminates | `cn cell return` → `status:review → status:changes`; cycle does not resume automatically |
| `clarify` | Operator needs clarification before deciding | HI surfaces the clarification request; no label transition until clarification received |

### 1.3. Findings array

On `iterate` or `reject`, the artifact body MUST list one or more findings. Each finding is a Markdown section with the following sub-fields:

| Sub-field | Required | Description |
|---|---|---|
| `id` | yes | Short identifier (e.g. `O1`, `F1`, `R1`) unique within this artifact |
| `surface` | yes | The file path or named surface the finding addresses |
| `problem` | yes | Precise description of the problem |
| `expected_change` | yes | What the operator expects the next round to address |
| `class` | no | Classification (e.g. `semantic-precision`, `doctrinal-vocabulary`, `citation-depth`) |

### 1.4. HI attribution rule

The `captured_by` field MUST be distinct from α/β/γ/δ artifact attribution. Valid values include:
- `gamma-interface (HI)` — when the HI session is the γ-interface mouthpiece
- `sigma (HI)` — when the HI session is the Sigma agent-admin
- `human-operator-direct` — when the operator authors the artifact directly (rare; no HI translation layer)

The `captured_by` value MUST NOT be `alpha`, `beta`, `gamma`, or `delta` — those identities are reserved for dispatched role sessions.

---

## 2. Example artifact

```yaml
---
schema: cn.operator-review.v1
issue: 500
pr: 501
verdict: iterate
reviewer: human-operator
captured_by: gamma-interface (HI)
captured_at: 2026-07-01 (UTC)
worker_pr_head_at_review: abc1234f
findings_count: 2
---
```

```markdown
# Operator-review — PR #501 (cnos#500 cycle/500 R1)

## Findings

### O1 — Missing test for converge path

**Surface:** `src/go/internal/cell/return_test.go`

**Problem:** The test suite covers iterate and reject verdicts but not converge. The `converge` verdict must not apply the label transition; this is untested.

**Expected change:** Add a test that calls `return` with `--verdict converge` and verifies no label mutation occurs.

**Class:** test coverage gap.

### O2 — Help text missing example

**Surface:** `src/go/internal/cli/cmd_cell.go`

**Problem:** The help text for `cn cell return` does not include a usage example with the `--review` flag.

**Expected change:** Add one concrete example to the help text: `cn cell return --issue 500 --verdict iterate --review .cdd/unreleased/500/operator-review.md`.

**Class:** UX clarity.
```

---

## 3. degraded_recovery declaration schema

### 3.1. Purpose

The `degraded_recovery` declaration records when the HI applied an operator patch under bootstrap-exception conditions — i.e., when no mechanical review-return primitive existed and the HI was forced to absorb operator corrections inline, crossing role boundaries.

The declaration is NOT a cover story for arbitrary HI overreach. It is an honest acknowledgment that:
- the underlying cause was a missing mechanical primitive
- the override was necessary to close the cycle
- the exception must not become the default

### 3.2. Required fields

A `degraded_recovery` declaration MUST carry:

| Field | Required | Description |
|---|---|---|
| `degraded_recovery` | yes | Fixed key identifying the declaration type; value is the recovery type string: `human_interface_applied_operator_patch` |
| `reason` | yes | Narrative reason: why the HI was forced to overstep |
| `scope` | yes | Precise description of what the HI did that crossed role boundaries |
| `recovery_actions` | yes | Ordered list of actions taken to restore role attribution |
| `status` | yes | One of: `accepted as bootstrap exception` / `pending recovery` / `closed` |
| `governing_doctrine` | yes | Path(s) to the doctrine surface(s) whose failure mode this violates |
| `target_state` | yes | Description of the future state once the missing primitive exists |

### 3.3. Placement

The declaration lives in the cycle's `gamma-closeout.md` under a clearly named section (e.g., `## §N. degraded_recovery declaration`). It is γ's artifact — γ records the HI's boundary-crossing and the recovery sequence.

### 3.4. Detection

The declaration is detectable by:

```bash
grep -rn "degraded_recovery" .cdd/
```

A cycle that carried a `degraded_recovery` declaration surfaces in this grep. The key name `degraded_recovery` is stable; tooling can detect and count declarations across cycles.

### 3.5. Canonical witness

`.cdd/unreleased/497/gamma-closeout.md §5` is the first conforming witness for this schema. It was authored before this schema was formally defined; it conforms to the required fields as defined here (all required fields are present with correct semantics).

---

## 4. Authoring rules for the HI

### 4.1. What the HI does

1. Read the PR and the operator's verdict carefully.
2. Translate the verdict into a `cn.operator-review.v1` artifact: populate all required fields; translate each finding into the `id / surface / problem / expected_change` structure.
3. Commit the artifact to the cycle branch at `.cdd/unreleased/{N}/operator-review.md`.
4. Push to `origin/cycle/{N}`.
5. Invoke `cn cell return --issue N --verdict V --review .cdd/unreleased/{N}/operator-review.md` to apply the mechanical label transition.
6. (On `iterate`) Invoke `cn cell resume --issue N` to re-arm the cycle.
7. Report to the human that the transition is complete and describe the next step.

### 4.2. What the HI does NOT do

The HI MUST NOT author or amend the following role-owned artifact surfaces:

- `self-coherence.md` — α's artifact
- `beta-review.md` — β's artifact
- `alpha-closeout.md` — α's artifact
- `beta-closeout.md` — β's artifact
- `gamma-closeout.md` — γ's artifact
- `gamma-scaffold.md` — γ's artifact

These surfaces belong to the dispatched roles. The HI's artifact is `operator-review.md`. The HI may read any surface; it may write only `operator-review.md` and apply mechanical transitions via `cn cell return`/`cn cell resume`.

**Failure mode:** "invisible meddling" — per `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` named failure mode. When the HI authors role-owned matter, role attribution in the durable record is corrupted. Future verification, retrospective, and TSC grading cannot determine who claimed what.

**Bootstrap-exception escape hatch:** The HI may patch role-owned matter ONLY when (a) the missing mechanical primitive is the underlying cause AND (b) the override is explicitly declared as `degraded_recovery: human_interface_applied_operator_patch` in the cycle's `gamma-closeout.md` with reason + scope + recovery_actions + status fields. The escape hatch is a doctrinal acknowledgment that the bootstrap path exists; it must not become the default.

---

## 5. Lifecycle context

The `cn.operator-review.v1` artifact is the input surface for the `cn cell return` command (AC3) and the `cn cell resume` command (AC4). The flow:

```text
Operator reads PR at status:review
  → HI authors operator-review.md (this schema)
  → cn cell return --issue N --verdict iterate --review operator-review.md
      → status:review label removed
      → status:changes label added
  → cn cell resume --issue N
      → cycle branch preserved (cycle/{N})
      → artifact directory preserved (.cdd/unreleased/{N}/)
      → R[N+1] section appended to self-coherence.md
  → δ routes α R[N+1] → β R[N+1] → γ closeout amendment
  → HI reports result; HI does not author role-owned matter
```

---

## 6. Related surfaces

- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §5` — "δ does not produce matter"; "invisible meddling" failure mode
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.6` — `status:changes` carve-out; `cn cell return` is the mechanical translator of operator authority
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.10` — `resumed-from-changes` wake-invoked mode shape (this cycle's §9 amendment)
- `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md` — HI behavioral contract (explicit prohibited surfaces)
- `src/go/internal/cell/` — Go implementation of `cn cell return` and `cn cell resume`
- `.cdd/unreleased/497/operator-review.md` — canonical first-use witness (conforming to this schema; filed before schema was formally written)
