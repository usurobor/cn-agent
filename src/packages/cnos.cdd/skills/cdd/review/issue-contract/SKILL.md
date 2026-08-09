---
name: review/issue-contract
description: Issue contract walk — verify AC coverage, named doc updates, CDD artifact contract, and active skill consistency before diff inspection.
artifact_class: skill
kata_surface: external
governing_question: How does β verify that the branch satisfies the issue contract (ACs, named docs, CDD artifacts, active skills) before reading the diff?
visibility: internal
parent: review
triggers:
  - review issue contract
  - review AC coverage
scope: task-local
inputs:
  - contract / acceptance criteria (issue text, PR body, or equivalent)
  - diff (or artifact set) at a stated base
  - protocol-supplied artifact locations, when the invoking protocol names them (optional context per review/SKILL.md rule 3.11b)
  - contract integrity results from Phase 1
outputs:
  - AC coverage table
  - named doc updates table
  - CDD artifact contract table
  - active skill consistency table
requires:
  - Phase 1 (contract integrity) completed
  - review orchestrator loaded
calls: []
---

# Review — Phase 2a: Issue Contract Walk

**PRE-GATE: §2.0.0 Contract Integrity (Phase 1) must be completed.** When the review target is a branch, verify it is unmerged at the stated base; when reviewing a diff or artifact set directly, this check is n/a.

---

## §2.0 Issue contract walk

Build these tables before reading the diff:

### AC Coverage
| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|

### Named Doc Updates
| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|

### CDD Artifact Contract
| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|

This table covers only artifacts the contract under review explicitly names as deliverables. Protocol-supplied artifact locations (a coordination scaffold, a wave manifest, a closure record) are optional context per `review/SKILL.md` rule 3.11b — verify them against the contract when supplied; when none are supplied, the table is n/a.

### Active Skill Consistency
| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|

Verify that every Tier 3 skill the issue requires (and every active design constraint named in the project) was actually loaded by α and is applied somewhere in the diff. A required skill that is "loaded" but not "applied" is silent process theater. A skill applied without being declared is a discoverability gap.

---

## After Phase 2a

Collect table results. Pass to Phase 2b (`review/diff-context/SKILL.md`) for diff and context inspection.

Return to orchestrator (`review/SKILL.md`) after all Phase 2 sub-skills complete for verdict rules and output format.
