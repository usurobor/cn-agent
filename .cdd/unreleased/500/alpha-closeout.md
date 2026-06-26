---
cycle: 500
role: alpha
status: final
beta_verdict: APPROVE (converge) — R1
---

# α close-out — cnos#500 — cdd/review-return

## Cycle summary

cycle/500 installs the review-return primitive: the missing mechanism that allows the HI to route an operator `iterate` verdict back into an existing cell at `status:review` without crossing role boundaries.

Seven artifacts shipped:
1. `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` — `cn.operator-review.v1` schema + `degraded_recovery` declaration schema
2. `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md` — HI behavioral contract; prohibited surfaces; bootstrap-exception escape hatch
3. `src/go/internal/cell/cell.go` — `cn cell return` + `cn cell resume` Go implementation; 26 tests pass at R1 HEAD
4. `src/go/internal/cell/cell_test.go` — 26 tests (24 R0 + 2 R1 positive-path injection tests); all pass
5. `src/go/internal/cli/cmd_cell.go` — CLI command registration
6. `src/go/cmd/cn/main.go` — kernel registry registration
7. `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.10` — `resumed-from-changes` wake-invoked mode shape; §9.6 reconciliation note

β R1 (converge) confirmed: all 7 ACs pass. 26/26 tests pass. Three R0 findings resolved (F1 B, F3 B, F2 C addressed). No new findings at R1.

## Friction log

1. **Design-pass judgment on cycle sizing.** The issue body suggested splitting into 500A/500B/500C. α's design pass produced a keep-whole justification: the 7 ACs are a single conceptual unit (schema + contract + CLI + δ amendment); they split poorly because each AC references the others. The judgment was documented in `self-coherence.md §R0 §Design decisions` before implementation.

2. **gitignore interaction.** `src/go/cmd/cn/` is listed in `.gitignore` but `main.go` is already tracked. `git add` on the directory path fails with "paths are ignored"; `git add -f` on the file path works. Belt-and-suspenders for future cycles: use file paths, not directory paths, when staging files in gitignored directories with pre-tracked files.

3. **AC6 enforcement gap.** CI grep of HI authorship signatures is not feasible without a stable machine-distinguishable HI signature. Named as known debt; convention-based enforcement with β Rule 7 is the backstop. This is an honest scope boundary.

4. **§9.6 carve-out reconciliation.** The existing δ SKILL.md §9.6 carve-out ("status:changes is EXTERNAL") needed a reconciliation note to explain how `cn cell return` co-exists with the carve-out without contradicting it. The note was added in the same edit as §9.10 to avoid two separate structural edits to adjacent prose.

5. **Citation accuracy (F1, β R0).** `hi-contract.md` and `operator-review/SKILL.md` both cited "invisible meddling" as a named failure mode in `delta/SKILL.md`. The phrase is in `operator/SKILL.md §Core Principle` (line 37). Fixed at R1: citations corrected to `operator/SKILL.md §Core Principle` in all four affected lines (hi-contract.md §2, §7; operator-review/SKILL.md §4.2, §6).

6. **Frontmatter compliance (F3, β R0).** `operator-review/SKILL.md` calls: entries contained inline section anchors that fail `validate-skill-frontmatter.sh` path resolution; `kata_surface:` field absent. Fixed at R1: bare paths + `kata_surface: none` added.

7. **Test injection (F2, β R0 C-severity).** `applyLabelTransition()` had no injection point. Fixed at R1: `RunGH func` field added to `Returner`; two new positive-path tests confirm iterate and reject transitions. Total test count advanced from 24 to 26.

## Observations

- The `operator-review/SKILL.md` governing-question structure (one file, one question) maps cleanly to write/SKILL.md's principle; each section answers one sub-question.
- cycle/497's `degraded_recovery` declaration in `gamma-closeout.md §5` was already present at scaffold time (per the γ-provided friction note in the dispatch prompt). AC7 only required schema formalization, not a cross-branch patch. The pre-checked friction note was accurate.
- The noun-verb dispatch pattern (`cell-return`, `cell-resume`) means `cn cell` automatically works as a group listing command without any additional code (the dispatch.go `GroupMembers` function finds all `cell-*` commands). Clean extension of the existing pattern.
- Citation accuracy in authoritative contract documents (hi-contract.md, operator-review/SKILL.md) is a class of finding where both the author and β Rule 7 need to verify against the actual file, not from memory. The same `grep` discipline that applies to code-first oracle anchoring applies to doctrine citations.
- The provisional close-out mechanism (alpha/SKILL.md §2.8 fallback) correctly covered the R0 state. At β APPROVE, the provisional note is resolved and this final version records the full outcome.
