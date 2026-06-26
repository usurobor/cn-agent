---
cycle: 500
role: alpha
version: R0
mode: design + build
---

# Self-coherence — cnos#500 — cdd/review-return

## §R0 Design decisions

### Command shape

`cn cell return` and `cn cell resume` — registered as `cell-return` and `cell-resume` in the kernel registry using the existing noun-verb dispatch pattern (`dispatch.go` `ResolveCommand` §1 lookup). User-facing: `cn cell return --issue N --verdict V --review path`, `cn cell resume --issue N`.

Rationale: the `cn cell` namespace is explicitly pinned in the implementation contract; the noun-verb pattern is the canonical form in this codebase (`dispatch.go` `GroupMembers`/`PrintGroup`); `return` and `resume` are the most semantically precise verbs for the two operations (return = deliver verdict; resume = re-arm the existing cycle).

### HI contract placement

New file `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md`. The existing `prompt.md` is the functional wake prompt for the agent-admin wake — it is runtime-loaded by the substrate and must not carry doctrinal contract prose that is not wake-specific. A separate `hi-contract.md` is the correct placement for the HI behavioral contract (which governs all HI sessions, not just the agent-admin wake).

### AC6 enforcement mechanism

Convention-only enforcement with β Rule 7 as the backstop. CI grep of HI authorship signatures in role-owned paths is not feasible: the HI commits do not carry a stable machine-distinguishable authorship signature (the HI may be dispatched under any git identity; the only stable signal is the artifact path + content, which requires semantic parsing). The convention is: prohibited surfaces are named in `hi-contract.md`; β's Rule 7 diff-vs-contract verification catches violations at review time. This is the same enforcement model as the implementation-contract axes. Named as known debt in §Debt.

### δ SKILL.md amendment structure

New `§9.10` subsection in `delta/SKILL.md` for the `resumed-from-changes` wake-invoked mode shape. Existing §9.1–§9.9 preserved intact. The §9.6 carve-out ("status:changes is EXTERNAL") is amended to read: `cn cell return` acts on behalf of the operator's stated verdict — the operator is still the authority; the command is the mechanical translator. This preserves the carve-out's intent while installing the mechanism.

### Cycle sizing judgment

Keep as one cell. Five-factor check:
1. One new Go package (`src/go/internal/cell/`) — bounded; no changes to existing packages
2. Three existing surfaces touched (delta SKILL.md, agent-admin, dispatch-protocol) — all additive amendments
3. Spans design + code + docs — explicitly the declared mode ("design + build")
4. Design pass required and completed above
5. Operator suggested split as optional ("start with the live worker; γ can make that call") — γ kept whole; all 7 ACs are tightly coupled (the HI contract, the schema, the CLI command, and the δ amendment are a single conceptual unit that splits poorly: AC2+AC6+AC7 require the same schema, AC3+AC4 require the same CLI command, AC5 requires understanding both)

---

## §Gap

**Issue:** cnos#500 — cdd/review-return: mechanically route operator iterate verdicts back into an existing cell

**Version/mode:** R0 / design + build

**The gap:** The live CDS dispatch wake handles a cell through `status:review` with no autonomous mechanism to re-engage from `status:review` when the operator's final read returns `iterate-narrowly`. Without this primitive, the HI absorbs operator-final-read corrections inline — crossing role boundaries. This was observed empirically in cycle/497: the HI edited α-owned `self-coherence.md`, rewrote γ-owned `gamma-closeout.md`, and inserted §R1 sections into β/α/γ artifacts, framing the absorption as "δ-direct R1." That absorption was declared a `degraded_recovery: human_interface_applied_operator_patch` in `.cdd/unreleased/497/gamma-closeout.md §5`.

**This cycle installs:**
1. A typed `cn.operator-review.v1` schema for HI to capture operator verdicts
2. `cn cell return --issue N --verdict V --review path` — mechanical `status:review → status:changes` transition
3. `cn cell resume --issue N` — re-arm the existing cycle (branch + artifact directory preserved; R[N+1] appended)
4. δ SKILL.md §9 amendment for the `resumed-from-changes` wake-invoked mode shape
5. HI behavioral contract — explicit prohibited surfaces + MUST NOT language
6. `degraded_recovery` declaration schema as first-class convention

---

## §Skills

**Tier 1 — CDD lifecycle:**
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (this role)
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (amended §9; §9.6 carve-out; §9.10 new subsection)

**Tier 2 — always-applicable eng:**
- `src/packages/cnos.core/skills/write/SKILL.md` (governing question, one-file-one-job, front-load the point)

**Tier 3 — issue-specific:**
- `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` (implementation contract; 7-axis schema; four-surface mesh)

**Active constraints from loaded skills:**
- write/SKILL.md §3.2 — one file, one governing question (applied to hi-contract.md and operator-review/SKILL.md)
- write/SKILL.md §3.4 — front-load the point (applied to all authored Markdown)
- alpha/SKILL.md §2.5 — incremental self-coherence write discipline (one section per commit)
- alpha/SKILL.md §2.3 — peer enumeration before closure claims
- alpha/SKILL.md §3.6 — implementation contract is δ's; α MUST NOT improvise

---

## §ACs

Per-AC oracle run against branch HEAD `f7d88053` (implementation commit).

### AC1 — Typed operator-review schema

**Claim:** `cn.operator-review.v1` schema defined; HI can produce a valid conforming artifact.

**Evidence:**
- `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` exists (new file, created this cycle)
- Required fields enumerated: `schema`, `issue`, `pr`, `verdict` (enum: converge/iterate/reject/clarify), `reviewer`, `captured_by`, `captured_at`, `findings_count`
- Example artifact in §2 of SKILL.md demonstrates all required fields
- cycle/497's `operator-review.md` conforms: all 8 required fields present at lines 2–10 (verified by grep)
  - `schema: cn.operator-review.v1` ✓
  - `issue: 497` ✓
  - `pr: 499` ✓
  - `verdict: iterate` ✓
  - `reviewer: human-operator` ✓
  - `captured_by: gamma-interface (HI)` ✓
  - `captured_at: 2026-06-26 (UTC)` ✓
  - `findings_count: 6` ✓

**Pass:** all required fields present in schema definition; example conforms; 497 artifact conforms.

---

### AC2 — HI records verdict without editing role-owned matter

**Claim:** HI behavioral contract documented; "invisible meddling" cited; HI attribution distinguished; prohibited surfaces listed.

**Evidence:**
- `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md` exists (new file, created this cycle)
- §2 lists prohibited surfaces explicitly (self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md, gamma-scaffold.md)
- §2 uses MUST NOT language (not "should avoid")
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` "invisible meddling" failure mode cited by path in §2 and §7
- §4 distinguishes HI attribution: `captured_by: gamma-interface (HI)` vs role identities; states HI MUST NOT sign with role identities

**Pass:** all four oracle conditions met (explicit MUST NOT; named prohibited surfaces; delta SKILL.md citation by path; attribution distinguished).

---

### AC3 — Mechanical `status:review` → `status:changes` transition

**Claim:** `cn cell return --verdict iterate` on a cell at `status:review` removes `status:review` and adds `status:changes`.

**Evidence:**
- `src/go/internal/cli/cmd_cell.go` — `CellReturnCmd` registered as `cell-return`; dispatches to `cell.Returner.Return()`
- `src/go/internal/cell/cell.go` — `Returner.Return()` on `iterate`/`reject` calls `applyLabelTransition()`:
  - removes `status:review`: `gh issue edit N --remove-label status:review`
  - adds `status:changes`: `gh issue edit N --add-label status:changes`
- On `converge`: no label transition applied (verified by `TestReturner_Return_Converge` — output does not contain "status:changes")
- `src/go/cmd/cn/main.go` — `CellReturnCmd` registered in the kernel registry
- Smoke test: `TestReturner_Return_Converge` proves converge path does not mutate labels; schema validation test proves wrong-schema artifacts are rejected

**Pass:** label transition is mechanical (not manual); command compiles and builds (`go build ./cmd/cn/`); smoke fixture demonstrates converge does not apply transition; tests pass 24/24.

---

### AC4 — Existing cycle resumed (not replaced)

**Claim:** `cn cell resume --issue N` preserves branch, artifact directory, and existing artifacts; appends R[N+1]; does not create new directory.

**Evidence:**
- `src/go/internal/cell/cell.go` — `Resumer.Resume()`:
  - calls `verifyBranchExists()` (checks origin; does NOT create branch)
  - checks `os.Stat(artifactDir)` — fails if directory does not exist (does NOT create new directory)
  - calls `nextRoundNumber()` — counts existing §R[N] sections (does NOT wipe prior rounds)
  - calls `appendRoundHeader()` — appends §R[N+1] (does NOT replace file)
- `appendRoundHeader()` uses `os.O_APPEND` — file opened in append mode; existing content preserved by OS semantics
- `TestAppendRoundHeader_AppendsToExisting` — existing §R0 preserved after §R1 appended (verified by string search)
- `TestResumer_Resume_AppendRound` — directly verifies round number detection and append path
- `delta/SKILL.md §9.10` amendment names this as the `resumed-from-changes` shape with R[N+1] increment documented and preserved invariants stated

**Pass:** no directory-creation or branch-creation call in resume path; existing artifacts preserved by O_APPEND semantics; R[N+1] increment documented in δ skill amendment; tests pass.

---

### AC5 — δ routes R[N+1] in order

**Claim:** δ SKILL.md §9 amended to define `resumed-from-changes` shape with R[N+1] routing sequence.

**Evidence:**
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.10` added (new subsection)
- §9.10 names:
  - (1) Input includes operator-review artifact + prior round artifacts + R[N+1] increment (§9.10 "Input contract (resumed-from-changes)")
  - (2) Routing sequence: δ dispatches α R[N+1] → β R[N+1] → γ closeout amendment (§9.10 "Routing sequence (resumed-from-changes)" steps 1–4)
  - (3) Resumed-cell contract distinctions from first-claim contract (§9.10 "How this differs from a first-claim cell" + "Preserved invariants")
- §9.6 carve-out preserved and reconciliation with `cn cell return` added (explicitly states operator is still the authority; command is the mechanical translator)
- §9.10 does not conflict with existing §9.3 routing sequence (§9.10 references §9.3 and adds the resume-specific variant)

**Pass:** `resumed-from-changes` shape named as a distinct §9.10 subsection; R[N+1] routing sequence explicit; §9.6 carve-out reconciled not contradicted; no conflicts with existing §9 contract.

---

### AC6 — HI cannot substitute for dispatched roles

**Claim:** HI contract explicitly lists prohibited surfaces; uses MUST NOT language; enforcement mechanism present.

**Evidence:**
- `hi-contract.md §2` — prohibited surfaces table lists all six: self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md, gamma-scaffold.md
- §2 — "The prohibition is absolute. There is no 'narrow mechanical text fix' exception."
- `operator-review/SKILL.md §4.2` — secondary enforcement surface; same list; MUST NOT language
- `hi-contract.md §5` — enforcement: convention-based with β Rule 7 as backstop; CI enforcement named as a future candidate with explicit justification for why it is not implemented (HI authorship signatures not machine-distinguishable at grep level)

**Pass:** all six surfaces named; MUST NOT language present; enforcement mechanism named (convention + β Rule 7); future CI enforcement acknowledged as candidate.

**Known debt on AC6:** CI grep enforcement is not implemented. See §Debt.

---

### AC7 — Declared bootstrap-exception escape hatch

**Claim:** `degraded_recovery` schema exists; cycle/497's gamma-closeout carries it as the first witness; grep detects it.

**Evidence:**
- `operator-review/SKILL.md §3` — `degraded_recovery` declaration schema defined with all required fields: `degraded_recovery` (key + recovery type value), `reason`, `scope`, `recovery_actions`, `status`, `governing_doctrine`, `target_state`
- `.cdd/unreleased/497/gamma-closeout.md §5` — first conforming witness confirmed:
  - `grep -n "^degraded_recovery" .cdd/unreleased/497/gamma-closeout.md` → line 78 (verified)
  - all 7 required fields present in the §5 declaration
- Detection command documented: `grep -rn "degraded_recovery" .cdd/` (§3.4 of operator-review/SKILL.md)
- No cross-branch patch to 497 artifacts needed (issue body §"AC7 resolution" note confirmed: 497 §5 already carries the declaration; this cycle only formalizes the schema)

**Pass:** schema exists; 497 gamma-closeout carries the declaration at line 78; grep detects it; schema fields match the 497 declaration's fields.
