---
name: review
description: Orchestrate β review in phases — contract integrity, then implementation, then verdict. Each phase loads its own sub-skill.
artifact_class: skill
kata_surface: external
governing_question: Is the branch coherent enough to merge, or what specific incoherence must be named and fixed?
visibility: internal
parent: cdd
triggers:
  - review
  - pr
  - approve
  - request changes
scope: task-local
inputs:
  - branch
  - issue
  - diff
  - .cdd/unreleased/{N}/ artifacts
outputs:
  - beta-review.md
requires:
  - β completed intake and skill loading
  - canonical CDD.md loaded
calls:
  - review/contract/SKILL.md
  - review/issue-contract/SKILL.md
  - review/diff-context/SKILL.md
  - review/architecture/SKILL.md
calls_dynamic:
  - source: project design constraints
kata_ref: src/packages/cnos.cdd.kata/katas/M2-review/
---

# Review — Orchestrator

## Core Principle

A review is a witnessed judgment, not a rubber stamp. The reviewer owes the same honesty discipline as the implementer: name what is wrong, trace it to evidence, and state what must change.

```text
Review = tri(contract truth, implementation/projection evidence, witnessed verdict)
```

The review proceeds in three phases. **Each phase loads its own sub-skill.** Do not attempt all phases from memory — the sub-skills contain the detail.

---

## Phases

### Phase 1: Contract integrity

**Load:** `review/contract/SKILL.md`

Before reading the diff for implementation correctness, verify that the work contract (issue + PR body + branch summary) is internally consistent, truthful about status, and non-contradictory.

Complete the Contract Integrity table. If any row is "no," the review cannot approve unless the row is explicitly out of scope and the reviewer names why.

### Phase 2: Implementation review

Phase 2 runs three sub-skills sequentially. Load each in order:

1. **`review/issue-contract/SKILL.md`** — issue contract walk (AC coverage, named doc updates, CDD artifact contract, active skill consistency)
2. **`review/diff-context/SKILL.md`** — diff and context inspection (structural closure, multi-format parity, snapshot consistency, stale paths, authority conflicts, architecture leverage, design constraints)
3. **`review/architecture/SKILL.md`** — architecture and design check (7 questions A–G; load `cnos.core/skills/design/SKILL.md` when active)

Walk the issue contract (ACs, named docs, CDD artifacts). Read the diff and its neighbors. Apply mechanical scans, architecture checks, and evidence-bound findings.

### Phase 3: Verdict

Return to this file for verdict rules and output format.

---

## Verdict Rules

3.1. **Every claim traces to evidence**
  - No "seems wrong." Point to a line, commit, file, artifact, or behavior.

3.2. **Name the severity**
  - D — blocker, demonstrable incoherence
  - C — significant incoherence, non-blocking
  - B — improvement opportunity
  - A — polish

3.3. **All findings must be resolved before merge**
  - There is no "approved with follow-up."
  - D findings block merge. C/B/A findings must be fixed on-branch before merge.
  - Only exception: finding requires a design decision outside issue scope → "deferred by design scope," author files issue before merge.

3.4. **Verdict before details**
  - Lead with APPROVED / REQUEST CHANGES.

3.5. **No phantom blockers**
  - Only block on incoherence you can demonstrate.

3.6. **Approve when coherent, not when perfect**
  - The bar is coherence, not taste.

3.7. **Close the search space on approval**
  - Approval explicitly states that no remaining blocker was found in the relevant contract.

3.8. **Evidence depth matches claim strength**
  - structural claim → unit/schema proof may suffice
  - runtime behavior claim → path/integration proof required
  - operator contract claim → output/projection/state artifact required

3.9. **Specify regression pairs for D-level findings**
  - Every D-level finding includes positive case + negative case.

3.10. **CI / release-gate state**
  - If merge is requested, verify required CI/build checks are green on branch head.
  - If checks are missing/stale/red, approval is provisional: "Do not merge until checks finish green."

3.11. **Merge instruction is explicit**
  - Names the exact branch and merge action with `Closes #N` in the merge commit.

3.12. **Review divergence is a skill gap**
  - When two reviewers diverge, the fix is a patch to the review skill, not "be more careful."

3.13. **Honest-claim verification**

  Documents claim things; β verifies the claims are backed by code, data, or canonical source. Three sub-checks, all binding:

  - **(a) Reproducibility** — Every measurement quoted in a doc must be reproducible from artifacts in this commit. If the doc says "engine output = 0.83", the run that produced 0.83 must be runnable from the diff with provenance attached. A measurement with no reproduction path is a D-level finding.
  - **(b) Source-of-truth alignment** — Every term used in a non-spec doc must trace to its canonical definition. Drift between informal and normative usage is a D-level finding (e.g. "W2 spread" used inconsistently with `cnos.cdd/skills/.../W2` spec definition).
  - **(c) Wiring claims** — If a module documents "X is wired into Y", β grep-checks that X actually appears in Y's call graph. A doc lying about wiring is the most expensive class of bug to find late, because consumers trust the doc and stop reading the code.

  Failure mode this rule catches: α produces a narrative document (release note, self-coherence report, post-release assessment, runtime spec) whose prose is internally consistent but whose claims are not backed by what the cycle actually shipped. The supercycle in cnos-tsc surfaced this pattern across multiple cycles — three of four findings on the v3.2.0 self-coherence report cycle were honest-claim violations (undeclared artifact, score discrepancy, missing provenance attachment).

---

## Finding Taxonomy

| Type | Definition | Examples |
|------|-----------|----------|
| **mechanical** | Caught by grep/diff/script | stale path, wrong branch name, broken link |
| **judgment** | Requires design/coherence assessment | missing AC, authority conflict, design trade-off |
| **contract** | Work contract is incoherent | issue contradiction, PR overclaim, draft-as-current, exception contradicts hard gate, proof plan missing |
| **honest-claim** | Doc claims something code/data doesn't back (rule 3.13) | non-reproducible measurement, term used inconsistently with spec, wiring claim that grep disproves |

Contract and honest-claim findings may overlap with mechanical or judgment — tag both when applicable.

Mechanical findings reaching review are **process bugs**. If >20% of findings in a cycle are mechanical, file a process issue.

---

## Severity

| Sev | Meaning | Merge readiness |
|-----|---------|-----------------|
| D | Demonstrable incoherence | not merge-ready |
| C | Real incoherence, locally non-blocking | not merge-ready until fixed |
| B | Improvement opportunity | not merge-ready until fixed |
| A | Polish | not merge-ready until fixed |

---

## Output Format

Written to `.cdd/unreleased/{N}/beta-review.md` **incrementally**. Each review pass (contract, implementation, verdict) is a separate commit+push to the cycle branch. Do not write the entire review in one generation — stream timeouts will discard partial work.

**Incremental write discipline:**
1. Write each pass as a separate operation (§2.0.0 Contract → §2.1 Implementation → Verdict)
2. Commit and push after each pass
3. If resuming after a failure, read what exists on the branch and continue from the last committed pass

Each round appends a new section.

```markdown
**Verdict:** APPROVED / REQUEST CHANGES

**Round:** N
**Fixed this round:** {commit hashes} closes {prior findings}
**Branch CI state:** green / provisional
**Merge instruction:** `git merge {branch}` into main with `Closes #{issue}`

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes / no / n/a | |
| Canonical sources/paths verified | yes / no / n/a | |
| Scope/non-goals consistent | yes / no / n/a | |
| Constraint strata consistent | yes / no / n/a | |
| Exceptions field-specific/reasoned | yes / no / n/a | |
| Path resolution base explicit | yes / no / n/a | |
| Proof shape adequate | yes / no / n/a | |
| Cross-surface projections updated | yes / no / n/a | |
| No witness theater / false closure | yes / no / n/a | |
| PR body matches branch files | yes / no / n/a | |

## §2.0 Issue Contract

### AC Coverage
| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|

### Named Doc Updates
| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|

### CDD Artifact Contract
| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|

### Active Skill Consistency
| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|

## Regressions Required (D-level only)

## Notes
```

---

## Checklist

Before submitting a review:

- [ ] Phase 1 (contract integrity) completed before Phase 2
- [ ] §2.0.0 Contract integrity table filled
- [ ] Status truth checked: shipped/current/draft/planned/non-goal not conflated
- [ ] Source-of-truth paths resolve and match canonical docs
- [ ] Issue/PR examples obey their own rules
- [ ] Hard gates do not appear in exception examples
- [ ] Path resolution base verified where paths are validated
- [ ] Proof plan includes oracle, positive, negative where required
- [ ] New CI/check/status surfaces update operator-visible projections
- [ ] PR body matches corrected branch files
- [ ] No witness theater: structure backed by rejection mechanism or honest caveat
- [ ] Every issue AC verified (met, partial, missing, deferred)
- [ ] Required named docs/files checked
- [ ] CDD artifacts exist and are internally consistent
- [ ] Mechanical diff scan: duplicates, branch names, snapshot plausibility
- [ ] Every claim traces to evidence
- [ ] Honest-claim verification (3.13a): every quoted measurement reproducible from this commit
- [ ] Honest-claim verification (3.13b): every term used in a non-spec doc traces to canonical source
- [ ] Honest-claim verification (3.13c): every wiring claim grep-verified
- [ ] Severity assigned to every finding
- [ ] Type assigned to every finding (mechanical / judgment / contract / honest-claim)
- [ ] D-level findings include regression test pairs
- [ ] CI/build checks green, or approval marked provisional
- [ ] Approval explicitly closes the search space
- [ ] Merge instruction names branch and merge action
- [ ] Verdict stated first

---

## After Review

- **Approved:** β merges branch into main with `Closes #N`, pushes, proceeds to release per `release/SKILL.md`.
- **Changes requested:** α fixes on branch, appends to `self-coherence.md`; β narrows on next round.

### Review identity

β uses a different git identity than α — different role names in the email local part. The canonical form is `{role}@{project}.cdd.cnos` (or the cnos elision `{role}@cdd.cnos`); see `operator/SKILL.md` §Git identity for role actors for the full prescription, rationale, and worked examples. The role separation is git-observable. The two-level form from cycle #287 is deprecated as of cycle #343.

---

## External kata

Practice and evaluation for this review skill live in:

`src/packages/cnos.cdd.kata/katas/M2-review/`

That kata exercises the contract-integrity preflight, implementation review, architecture check, finding taxonomy, active-skill consistency, and evidence-depth rules. The frontmatter `kata_ref` field above carries the same path for machine-readable linkage.
