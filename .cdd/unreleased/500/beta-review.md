---
cycle: 500
role: beta
round: R0
verdict: REQUEST CHANGES
origin_main_sha: 3095fa2b44145490c8e5241bd347165a53ace827
cycle_head_sha: 958fb0ac2a8361d0c184013f873a939ee78434ad
---

# β review — cnos#500 — cdd/review-return

## §R0

**Verdict: REQUEST CHANGES**

**Base:** `origin/main` @ `3095fa2b` (re-fetched synchronously before review; unchanged from scaffold time)
**Cycle HEAD:** `958fb0ac` (alpha-500: self-coherence review-readiness signal R1)

---

## Implementation contract conformance

| Axis | Pinned value | Diff result | Status |
|---|---|---|---|
| Language | Go (CLI) + Markdown (schema, HI contract, δ amendment, bootstrap schema) | `cell.go`, `cell_test.go`, `cmd_cell.go` in Go; all doc artifacts in Markdown | PASS |
| CLI integration target | `cn cell {return,resume}` namespace | Registered as `cell-return` / `cell-resume` in noun-verb dispatch pattern; `cn cell return` and `cn cell resume` user-facing forms via `ResolveCommand` noun-verb resolution | PASS |
| Package scoping | `src/go/internal/cli/cmd_cell.go` + `src/go/internal/cell/`; Markdown at operator-review SKILL.md, delta §9, agent-admin hi-contract | All six file paths present in diff at exact pinned locations | PASS |
| Existing-binary disposition | Preserve `cn`; `cn cell` additive | Two lines added to `main.go`; no existing commands modified | PASS |
| Runtime dependencies | None beyond existing | `gh` CLI assumed; no new language runtimes or Go deps | PASS |
| JSON/wire contract preservation | N/A — all artifacts new | No existing schemas or wire contracts modified | PASS |
| Backward-compat invariant | All existing lifecycle transitions preserved | `cell.go` adds new commands; existing `delta/SKILL.md` sections unchanged; `main.go` additions only | PASS |

Implementation contract: **all 7 axes conform.**

---

## Per-AC verdict

### AC1 — Typed operator-review schema

**Verdict: PASS**

**Evidence verified independently:**

`src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` exists (new file). Required frontmatter fields confirmed present in §1.1:

| Field | Present | Value/type |
|---|---|---|
| `schema` | yes | fixed string `cn.operator-review.v1` |
| `issue` | yes | integer |
| `pr` | yes | integer or null |
| `verdict` | yes | enum: converge/iterate/reject/clarify |
| `reviewer` | yes | string |
| `captured_by` | yes | string |
| `captured_at` | yes | string (ISO 8601) |
| `findings_count` | yes | integer |

**Design observation on `findings_count` vs `findings[]`:** The issue body AC1 lists `findings[]` as a required frontmatter field. The implementation substitutes `findings_count` (integer) in frontmatter + structured body Markdown sections for the findings array. This is a valid structural design decision (YAML frontmatter with large arrays is unwieldy; Markdown sections are the natural format for rich finding descriptions). The γ-scaffold AC1 oracle approach does not prescribe frontmatter array vs body Markdown — it asks that findings carry `id`, `surface`, `problem`, `expected_change`, which the body sections do. The schema change (array → count+body) does not lose information; `cn cell return` only reads `schema:` from frontmatter, not findings. β accepts this design decision.

**497 conformance:** `.cdd/unreleased/497/operator-review.md` frontmatter confirmed present: `schema: cn.operator-review.v1`, `issue: 497`, `pr: 499`, `verdict: iterate`, `reviewer: human-operator`, `captured_by: gamma-interface (HI)`, `captured_at: 2026-06-26 (UTC)`, `findings_count: 6`. All 8 required fields present. **Conforms.**

---

### AC2 — HI records verdict without editing role-owned matter

**Verdict: PASS**

**Evidence verified independently:**

`src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md` exists (new file).

- §2 prohibited surfaces table: all six listed explicitly (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `gamma-scaffold.md`). ✓
- §2 uses "MUST NOT" language: "The HI MUST NOT author or amend any of the following role-owned artifact surfaces." ✓
- §2 names "invisible meddling" failure mode: "Named failure mode: 'invisible meddling' — per `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`." ✓
- §4 HI attribution: `captured_by` values (`gamma-interface (HI)`, `sigma (HI)`, `human-operator-direct`) distinguished from role identities. ✓

**Finding F1 applies here** (see Findings section): "invisible meddling" is cited as being in `delta/SKILL.md` but the phrase does not exist in that file — it is in `operator/SKILL.md §Core Principle`. The citation in `hi-contract.md §2` and `operator-review/SKILL.md §4.2` is imprecise. The failure mode name and concept are correct; only the citation path is wrong. This is classified as a B-severity finding (should-fix: the citation in authoritative contract documents must resolve accurately, especially since β Rule 7 relies on file+section citations for audit; a miscited source is hard to verify in future reviews).

AC2 substance passes; citation correction required.

---

### AC3 — Mechanical `status:review` → `status:changes` transition

**Verdict: PASS with finding**

**Evidence verified independently:**

`src/go/internal/cell/cell.go` `Returner.Return()`:
- On `iterate`/`reject`: calls `applyLabelTransition()` → `gh issue edit {N} --remove-label status:review` + `gh issue edit {N} --add-label status:changes`. Label removal precedes addition. Correct two-step sequence. ✓
- On `converge`/`clarify`: no label transition. ✓
- Schema validation: reads `schema:` from frontmatter; rejects non-`cn.operator-review.v1` values. ✓

**Build:** `go build ./src/go/cmd/cn/` — no errors. ✓
**go vet:** no findings. ✓
**Tests:** 24 PASS, 0 FAIL. ✓

**Finding F2 applies here** (see Findings section): the smoke fixture only covers the negative case — `TestReturner_Return_Converge` proves converge does NOT apply `status:changes`. There is no unit test that exercises `applyLabelTransition()` (the positive iterate/reject path). `runGH` is a package-level function with no injection point on `Returner`, so the iterate path cannot be unit-tested without the `gh` CLI present. This is a C-severity finding (nice-to-have: adding a `runGHFunc` field on `Returner` or an interface would make the critical label-mutation path testable; the current code is correct but the positive path is untested). The self-coherence documents this as known (Debt §4), so β records it but does not require it as a blocking fix.

AC3 passes. F2 is C-severity (not blocking).

---

### AC4 — Existing cycle resumed (not replaced)

**Verdict: PASS**

**Evidence verified independently:**

`src/go/internal/cell/cell.go` `Resumer.Resume()`:
- `verifyBranchExists()` uses `git ls-remote --exit-code origin refs/heads/{branch}` — check only, no create. ✓
- `os.Stat(artifactDir)` — fails if missing; no `os.MkdirAll()` call. ✓
- `nextRoundNumber()` scans for `## §R<digits>` headers; returns `max+1`. ✓
- `appendRoundHeader()` opens with `os.O_WRONLY|os.O_CREATE|os.O_APPEND` — append-only. ✓

`TestAppendRoundHeader_AppendsToExisting`: existing `§R0` preserved after `§R1` appended — verified in test. ✓
`TestResumer_Resume_AppendRound`: `§R0` preserved, `§R1` appended. ✓

`delta/SKILL.md §9.10` names the `resumed-from-changes` shape; "Preserved invariants" section explicitly lists all required preservation rules. ✓

---

### AC5 — δ routes R[N+1] in order

**Verdict: PASS**

**Evidence verified independently:**

`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.10` added as a new subsection (54 lines, including routing sequence, detection criteria, input contract, preserved invariants, γ closeout amendment shape, empirical anchor, cross-references).

- Input contract names: operator-review.md + prior self-coherence §R[0..N] + prior beta-review §R[0..N] + R[N+1] increment. ✓
- Routing sequence steps 1–4 explicit: (1) δ dispatches α R[N+1]; (2) δ dispatches β R[N+1]; (3) δ routes on β verdict; (4) δ dispatches γ closeout amendment. ✓
- §9.6 reconciliation paragraph added: "operator is still the authority; `cn cell return` is the mechanical translator." The existing carve-out semantics are preserved. ✓
- §9.10 references §9.3 (first-claim routing sequence) as baseline and names distinctions. No conflict with existing §9 contract rows. ✓
- Existing §9.1–§9.9 unchanged (confirmed by diff: only one deletion in §9.6 + 54 new lines at end). ✓

---

### AC6 — HI cannot substitute for dispatched roles

**Verdict: PASS with finding**

**Evidence verified independently:**

`hi-contract.md §2` table: all six prohibited surfaces listed with owner and reason. ✓
`hi-contract.md §2`: "The prohibition is absolute. There is no 'narrow mechanical text fix' exception." ✓
`hi-contract.md §5`: enforcement mechanism named — convention + β Rule 7 backstop; future CI candidate with explicit justification. ✓
`operator-review/SKILL.md §4.2`: same six surfaces, same MUST NOT language. ✓

Two surfaces state the same prohibition list (hi-contract.md §2 and operator-review/SKILL.md §4.2). Both lists carry the same 6 surfaces with no drift. ✓

**Finding F1 applies** (imprecise citation of "invisible meddling" from `delta/SKILL.md` — see AC2). Both hi-contract.md and operator-review/SKILL.md cite this as a named failure mode in `delta/SKILL.md`. It is actually in `operator/SKILL.md §Core Principle` line 37. The substance of the prohibition is sound; the citation path is wrong.

---

### AC7 — Declared bootstrap-exception escape hatch

**Verdict: PASS**

**Evidence verified independently:**

`operator-review/SKILL.md §3` defines the `degraded_recovery` declaration schema. Required fields confirmed present: `degraded_recovery` (key + recovery type value `human_interface_applied_operator_patch`), `reason`, `scope`, `recovery_actions`, `status`, `governing_doctrine`, `target_state`. ✓

`.cdd/unreleased/497/gamma-closeout.md §5` confirmed as first canonical witness:
- `grep -n "^degraded_recovery" .cdd/unreleased/497/gamma-closeout.md` → line 78. ✓
- All 7 required fields present in the §5 YAML block (`degraded_recovery`, `reason`, `scope`, `recovery_actions`, `status`, `governing_doctrine`, `target_state`). ✓
- Detection command documented: `grep -rn "degraded_recovery" .cdd/` in §3.4 of operator-review/SKILL.md. ✓

No cross-branch write to 497 artifacts required: the 497 gamma-closeout already carries the declaration on `main`.

---

## Findings

### F1 — "invisible meddling" citation points to wrong file

**Severity: B (should-fix)**
**Classification: citation-accuracy**

**Surface:** `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md §2` (line: `"Named failure mode: 'invisible meddling' — per src/packages/cnos.cdd/skills/cdd/delta/SKILL.md."`) and `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md §4.2` (line: `"Failure mode: 'invisible meddling' — per src/packages/cnos.cdd/skills/cdd/delta/SKILL.md named failure mode."`) and `§6` cross-reference line `"delta/SKILL.md §5 — 'δ does not produce matter'; 'invisible meddling' failure mode"`.

**Problem:** The phrase "invisible meddling" and its definition ("the operator adjusts implementation, review, or coordination reasoning without declaring an override, and the triad's coherence record no longer matches what actually happened") appears in `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md §Core Principle` (line 37), NOT in `delta/SKILL.md`. Running `grep -n "invisible meddling" src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` returns zero results. The phrase is also not in the diff additions to `delta/SKILL.md`. The issue body references it as "operator doctrine; named failure mode: 'invisible meddling'" — which correctly attributes it to the operator doctrine surface, but hi-contract.md and operator-review SKILL.md both cite it as if it were in delta/SKILL.md.

**Impact:** Contract documents that rely on auditable cross-references cite a non-existent anchor. β Rule 7 (and future readers) following the citation will not find the named failure mode at the stated path. This weakens the audit trail for the contract the cycle is designed to formalize.

**Expected change (α):** Update hi-contract.md §2 citation to: `per src/packages/cnos.cdd/skills/cdd/operator/SKILL.md §Core Principle`. Update operator-review/SKILL.md §4.2 citation to the same. Update operator-review/SKILL.md §6 cross-reference from `delta/SKILL.md §5 — "δ does not produce matter"; "invisible meddling" failure mode` to two separate lines: `delta/SKILL.md §5 — "δ does not produce matter"` and `operator/SKILL.md §Core Principle — "invisible meddling" failure mode`.

---

### F2 — `applyLabelTransition()` has no unit-testable injection point

**Severity: C (nice-to-have)**
**Classification: test-coverage**

**Surface:** `src/go/internal/cell/cell.go` `Returner` struct + `runGH` function.

**Problem:** `runGH` is a package-level function invoked directly from `applyLabelTransition()`. `Returner` has no field for injecting a test double. The iterate and reject verdict paths — the core mechanical behavior this cycle ships — cannot be unit-tested without a live `gh` CLI. The existing smoke tests only verify the converge path (no `gh` call) and the schema-rejection path.

**Expected change (α if fixing, or operator decision):** Add a `RunGH func(ctx context.Context, args []string, w io.Writer) error` field to `Returner`, defaulting to the package-level `runGH` in production use. Tests can then inject a mock that records calls and returns success/failure. This makes the critical label-mutation sequence unit-testable. If the operator decides this is out of scope, declare it as a named debt item (it is already in §Debt as item 4 noting the `Repo` flag gap; the test injection gap should be added or noted there explicitly).

**Note:** The self-coherence AC3 evidence counts `TestReturner_Return_Converge` as the smoke fixture. That test only demonstrates converge does NOT apply `status:changes`. β widens this observation: the oracle asked that the "smoke fixture demonstrates the correct before/after label state" — i.e., the positive case. Currently no test demonstrates that iterate DOES trigger the label transition sequence. The test is a negative test, not a positive smoke. For a C-severity finding this does not block converge, but it is a gap in the stated oracle evidence.

---

### F3 — `operator-review/SKILL.md` calls field will fail `validate-skill-frontmatter.sh` CI check

**Severity: B (should-fix)**
**Classification: ci-compliance**

**Surface:** `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` frontmatter `calls:` field.

**Problem:** The `calls` entries in operator-review/SKILL.md include inline section anchors:

```yaml
calls:
  - cnos.cdd/skills/cdd/delta/SKILL.md §9.6 (return token carve-out)
  - cnos.cdd/skills/cdd/alpha/SKILL.md §2.1 (α dispatch intake)
```

The `validate-skill-frontmatter.sh` CI script resolves `calls` entries as filesystem paths using `jq` to extract the string literally, then appends it to the package skill root. The full string `"cnos.cdd/skills/cdd/delta/SKILL.md §9.6 (return token carve-out)"` (including the space and annotation text) is used as the path segment. This path will not exist on the filesystem. The existing SKILL.md convention uses bare relative paths (`delta/SKILL.md`, `design/SKILL.md`) without inline annotations.

**Secondary:** `operator-review/SKILL.md` is missing the `kata_surface` frontmatter field, which is "spec-required-but-exception-backed" per `tools/validate-skill-frontmatter.sh`. No exception entry for this file exists in `schemas/skill-exceptions.json`. This will trigger a validator finding.

**Expected change (α):** Fix `calls:` entries to use bare paths without section anchors: `delta/SKILL.md` and `alpha/SKILL.md` (relative paths from the package skill root `src/packages/cnos.cdd/skills/cdd/`). Move section references to inline prose within the skill body. Add `kata_surface: none` to the frontmatter (this skill has no embedded kata), OR add an exception entry in `schemas/skill-exceptions.json` with reason.

---

## Pre-merge gate assessment

| Row | Status |
|---|---|
| 1. Identity truth | β: `git config --get user.email` = `beta@cdd.cnos` ✓ |
| 2. Canonical-skill freshness | `origin/main` still at `3095fa2b` — same as scaffold time. No advancement. Skills loaded at intake remain current. ✓ |
| 3. Non-destructive merge-test | NOT executed pending RC resolution. β does not merge with open B-severity findings (F1, F3). |
| 4. γ-artifact completeness | `gamma-scaffold.md` present on `origin/cycle/500` ✓ |

---

## Overall verdict

**REQUEST CHANGES**

Two B-severity findings must be resolved before converge:

- **F1 (B):** "invisible meddling" cited as `delta/SKILL.md` but the phrase is in `operator/SKILL.md §Core Principle`. Fix citations in hi-contract.md and operator-review/SKILL.md.
- **F3 (B):** `operator-review/SKILL.md` `calls:` entries with inline section anchors will fail `validate-skill-frontmatter.sh`; `kata_surface` missing without exception entry.

One C-severity finding:

- **F2 (C):** `applyLabelTransition()` has no test injection point; iterate/reject paths untested. Advisory only — does not block converge. α decision: add `RunGH` field or explicitly extend §Debt.

**Substance is sound.** All 7 ACs pass on the core behavioral claims. The Go implementation is correct, clean, and all 24 tests pass. The δ §9.10 amendment is well-structured. The HI contract and schema are complete in content. The two B-severity findings are precision/compliance corrections that do not require rethinking the design.

**Routing:** α addresses F1 + F3 (and optionally F2); appends §R1 to self-coherence.md; signals review-readiness. β reviews §R1.
