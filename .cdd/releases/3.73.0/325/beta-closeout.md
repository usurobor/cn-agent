## Review context

**Issue:** #325 — cdd: systematic lifecycle audit and refactor for role, artifact, and release coherence
**Cycle branch:** `cycle/325`
**β session:** R1, R2, R3 (3 review rounds)
**Merge commit:** `1f07f7ba` on main — `merge(#325): CDD systematic lifecycle audit and refactor`
**Closes #325** (auto-close via `Closes #325` in merge commit message)

**Skill set loaded for review:**
- CDD.md (canonical lifecycle and role contract)
- beta/SKILL.md (β role surface)
- review/SKILL.md (review orchestrator)
- release/SKILL.md (release sub-skill)
- cnos.core/skills/design/SKILL.md (architecture/design check — active per review/SKILL.md §2.2.14)
- cnos.cdd/skills/cdd/issue/SKILL.md (issue contract audit)

---

## Merge evidence

| Item | Value |
|---|---|
| Merge commit SHA | `1f07f7ba` |
| Merge strategy | `--no-ff` (ort) |
| Branch merged | `cycle/325` → `main` |
| origin/main base SHA at β intake | `0ff6d427275b68998c8413ab8f26079928222c78` |
| cycle/325 approved HEAD SHA | `0b512c18` (R3 verdict commit) |
| Merge commit author | `beta@cdd.cnos` |
| Branch CI state at merge | no CI workflow for docs/skills-only changes |
| Diff summary | 9 files changed, 882 insertions, 19 deletions |

---

## Narrowing pattern across rounds

**Round 1 (R1):**
- 3 findings: F1 (D), F2 (D), F3 (C)
- F1: §5.3b PRA "Written when" column inverted the PRA/δ-preflight ordering
- F2: gamma/SKILL.md §2.10 missing δ release-boundary preflight gate row — γ could close without δ preflight
- F3: gamma/SKILL.md §2.10 "Then:" block instructed "write RELEASE.md" after gamma-closeout.md, contradicting gate row 11
- All three are intrastructural: the cycle was fixing CDD incoherence, and the new surfaces (§5.3b, §2.10) contained incoherence within themselves

**Round 2 (R2):**
- 1 finding: R2-F1 (C)
- R2-F1: §5.3b PRA "Verified by" omitted δ at release-boundary preflight — missed when F1 was fixed
- R1 fixes had corrected the "Written when" cell but did not update the adjacent "Verified by" cell
- Pattern: the F1 fix was correct but local; the affected invariant (δ checks PRA at preflight) spans two columns

**Round 3 (R3):**
- 0 findings
- R2-F1 fix verified: "Verified by" updated to `γ closure gate; δ at release-boundary preflight`
- Pre-merge gate passed: identity, canonical-skill freshness, merge-tree zero-conflict
- Approved

**Narrowing shape:** 3 D/C → 1 C → 0. The narrowing is clean. Each round targeted only what was wrong; no scope drift; no phantom blockers.

---

## β-side findings

### Observation 1: Self-referential incoherence pattern

All three R1 findings and R2-F1 are incoherence within the new surfaces the cycle itself ships — not in pre-existing text. The cycle added §5.3b (ownership matrix) and reinforced §2.10 (closure gate), and both contained inconsistencies with the adjacent §4.1a state table the cycle also added. This pattern — new surfaces that are internally consistent but mutually inconsistent — is specific to multi-surface refactors and is distinct from a single-surface error.

The cycle's correctness architecture is sound (design decisions D1–D7 are correct). The incoherence was cross-table: one table's "Written when" / "Verified by" cells did not match another table's "Required inputs" column for the same artifact.

### Observation 2: Three-round cycle for a docs-only change

The cycle took 3 review rounds for a P1 spec/skill-only change. The findings were real (two D-level in R1, one C-level in R2), so the rounds were not avoidable by better pre-review checking alone. The root cause is the multi-table consistency problem: when a single lifecycle event (δ preflight verifying PRA) is represented in three tables simultaneously (§4.1a Required inputs, §5.3b Written when, §5.3b Verified by), all three cells must agree and the pre-review gate has no mechanical check that forces cross-table consistency before signaling readiness.

### Observation 3: Validator environment constraint

`tools/validate-skill-frontmatter.sh` requires `cue` (not installed in the review environment). The SKILL.md frontmatter did not change in this cycle, so the validator was not blocking. But the constraint is observable: if a future cycle modifies frontmatter and `cue` is unavailable, the validator cannot run. This is an environment constraint, not a cycle failure.

---

## Voice check

The above findings are factual observations and patterns. No dispositions recommended. Triage is γ's job.
