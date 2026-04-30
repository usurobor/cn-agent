---
cycle: 309
role: beta
round: 1
---

**Verdict:** APPROVED

**Round:** R1
**Fixed this round:** n/a (R1 — no prior findings)
**origin/main SHA (review base):** `0e9762576aef170ff9aac851d79dd380d91e9046` (fetched synchronously before review)
**cycle/309 head SHA:** `5fd1eb14cdb025fbe495e83da43dbc3f9f4e2e60`
**Branch CI state:** provisional — `build.yml` triggers on `main` only; I5 frontmatter validation will run post-merge. `cue` not installed on host (exit code 2 = prerequisite missing, not schema failure); manual validation confirms schema shape.
**Merge instruction:** `git merge cycle/309 --no-ff -m "Merge cycle/309: eng/troubleshoot skill (#309)"` into `main` with `Closes #309`

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | "No live troubleshooting skill existed" matches repo state; no overclaims |
| Canonical sources/paths verified | yes | `eng/rca` cross-reference resolves; no stale paths |
| Scope/non-goals consistent | yes | No automated tooling, no runbook, no eng/rca changes introduced |
| Constraint strata consistent | yes | Five Tier 3 skills loaded per issue requirement (α attestation in §Skills) |
| Exceptions field-specific/reasoned | yes | D1 (cue unavailable) disclosed and reasoned in §Debt; no false closure |
| Path resolution base explicit | yes | Single file; no multi-file path chain |
| Proof shape adequate | yes | Invariant + oracle + positive + negative case; kata embodies both |
| Cross-surface projections updated | n/a | Single new file; no existing surface updated |
| No witness theater / false closure | yes | AC9 acknowledges cue unavailable; manual validation confirmed and disclosed |
| PR body matches branch files | yes | Diff = two files (SKILL.md + self-coherence.md); ACs map to diff evidence |

---

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | Artifact classification explicit (skill, not runbook; formula + governing question) | yes | **met** | `artifact_class: skill` in frontmatter; §1.0 explicit classification with justification; §Core Principle formula; `governing_question` in frontmatter and §1.3; named failure mode "premature hypothesis" |
| AC2 | Required skills loaded before writing (skill, write, design, test, rca) | yes | **met** | α §Skills lists all five Tier 3 skills loaded before drafting; CDD attestation mechanism |
| AC3 | External troubleshooting practices incorporated (Google SRE, IBM, Red Hat, CompTIA) | yes | **met** | IBM §2.2 (full 8-field problem description); Google SRE §2.3 footnote + §2.4; Red Hat §2.4 + §2.6; CompTIA §2.3 footnote; none copied wholesale |
| AC4 | Live triage algorithm explicit and numbered | yes | **met** | §2.1–§2.7 covers all 11 required steps; §2.3 six-step triage class order is numbered and sequential |
| AC5 | Hypothesis testing stated (cheapest first, one variable, oracle before test) | yes | **met** | §2.4 five-step hypothesis discipline; §2.5 one-change rule; §2.6 original-symptom verification; §3.2–§3.6 reinforce each element |
| AC6 | Three worked examples from dispatch test (OOM, GraphQL, background kill) | yes | **met** | §4 Examples 1–3; each shows symptom → wrong first hypothesis → better triage path → root cause → MCA → verification |
| AC7 | RCA handoff defined with triggers | yes | **met** | §2.7 five explicit triggers; §3.7 "do not start RCA during live diagnosis" |
| AC8 | Kata surface with scenario/symptoms/evidence/wrong hypothesis/expected path | yes | **met** | §5 has all required fields plus MCA, verification, common failures, reflection |
| AC9 | Frontmatter validates | yes | **met** | cue not installed (exit 2); manual validation against schemas/skill.cue confirms all hard-gate and spec-required-exception-backed fields present with correct types and enum values |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.eng/skills/eng/troubleshoot/SKILL.md` | yes | present | New file; the only deliverable |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `self-coherence.md` | yes | yes | 9-AC mapping, §Debt, §CDD-Trace, pre-review gate table |
| `beta-review.md` | yes | yes (this file) | |
| `beta-closeout.md` | yes | pending | Written after merge per release/SKILL.md §2.10 |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.core/skills/skill/SKILL.md` | issue §Skills | yes (α) | yes | artifact_class, formula, kata surface present |
| `cnos.core/skills/write/SKILL.md` | issue §Skills | yes (α) | yes | governing_question, point-first structure, no clutter |
| `cnos.core/skills/design/SKILL.md` | issue §Skills | yes (α) | yes | §1.0 explicit skill/runbook boundary; non-goals respected |
| `cnos.eng/skills/eng/test/SKILL.md` | issue §Skills | yes (α) | yes | invariant + oracle + positive/negative cases in proof plan and kata |
| `cnos.eng/skills/eng/rca/SKILL.md` | issue §Skills | yes (α) | yes | §2.7 RCA handoff boundary clearly defined; eng/rca not rewritten |

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| — | No findings | All 9 ACs met; contract integrity clean; no stale paths, no scope drift, no phantom debt | — | — |

## Regressions Required (D-level only)

None — no D-level findings.

## Notes

**N1 — cue prerequisite (D1 from α §Debt):** The frontmatter validator exits with code 2 (prerequisite missing) rather than code 1 (schema failure). This means the I5 CI job will also return exit 2 post-merge if `cue` is not installed in CI. This is an environmental constraint predating this cycle, not a schema error. Manual validation against `schemas/skill.cue` confirms all hard-gate fields (`name`, `description`, `governing_question`, `triggers`, `scope`) and all spec-required-exception-backed fields (`artifact_class`, `kata_surface`, `inputs`, `outputs`) are present with correct types and enum values. No exception entry required.

**N2 — Kata scenario reuses Example 1:** The embedded kata uses the OOM scenario from §4 Example 1. This is intentional and appropriate: OOM is the most instructive case (process-disappeared + no stderr + no error visible → dmesg check). The kata does not need to be distinct from the worked examples; it exercises the triage algorithm against a concrete case. This design choice is coherent.

**N3 — "What this skill forces" annotations:** Each worked example ends with a "What this skill forces:" paragraph explicitly mapping the example back to the triage step. This is a pedagogical addition beyond the AC requirement. It makes the triage algorithm concrete and actionable without adding clutter.

**CI state:** `build.yml` triggers on `main` only; cycle branches do not run CI. Post-merge, the I5 frontmatter validation job will run. The only failure mode is exit 2 (cue not installed), not exit 1 (schema violation). Approval is provisional on this basis; the schema shape is verified manually.
