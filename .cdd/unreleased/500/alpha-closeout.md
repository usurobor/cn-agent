---
cycle: 500
role: alpha
status: provisional — pending β outcome
---

# α close-out — cnos#500 — cdd/review-return

[provisional — pending β outcome per alpha/SKILL.md §2.8 bounded dispatch model]

## Cycle summary

cycle/500 installs the review-return primitive: the missing mechanism that allows the HI to route an operator `iterate` verdict back into an existing cell at `status:review` without crossing role boundaries.

Seven artifacts shipped:
1. `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` — `cn.operator-review.v1` schema + `degraded_recovery` declaration schema
2. `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md` — HI behavioral contract; prohibited surfaces; bootstrap-exception escape hatch
3. `src/go/internal/cell/cell.go` — `cn cell return` + `cn cell resume` Go implementation
4. `src/go/internal/cell/cell_test.go` — 24 tests; all pass
5. `src/go/internal/cli/cmd_cell.go` — CLI command registration
6. `src/go/cmd/cn/main.go` — kernel registry registration
7. `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.10` — `resumed-from-changes` wake-invoked mode shape; §9.6 reconciliation note

## Friction log

1. **Design-pass judgment on cycle sizing.** The issue body suggested splitting into 500A/500B/500C. α's design pass produced a keep-whole justification: the 7 ACs are a single conceptual unit (schema + contract + CLI + δ amendment); they split poorly because each AC references the others. The judgment was documented in `self-coherence.md §R0 §Design decisions` before implementation.

2. **gitignore interaction.** `src/go/cmd/cn/` is listed in `.gitignore` but `main.go` is already tracked. `git add` on the directory path fails with "paths are ignored"; `git add -f` on the file path works. Belt-and-suspenders for future cycles: use file paths, not directory paths, when staging files in gitignored directories with pre-tracked files.

3. **AC6 enforcement gap.** CI grep of HI authorship signatures is not feasible without a stable machine-distinguishable HI signature. Named as known debt; convention-based enforcement with β Rule 7 is the backstop. This is an honest scope boundary.

4. **§9.6 carve-out reconciliation.** The existing δ SKILL.md §9.6 carve-out ("status:changes is EXTERNAL") needed a reconciliation note to explain how `cn cell return` co-exists with the carve-out without contradicting it. The note was added in the same edit as §9.10 to avoid two separate structural edits to adjacent prose.

## Observations

- The `operator-review/SKILL.md` governing-question structure (one file, one question) maps cleanly to write/SKILL.md's principle; each section answers one sub-question.
- cycle/497's `degraded_recovery` declaration in `gamma-closeout.md §5` was already present at scaffold time (per the γ-provided friction note in the dispatch prompt). AC7 only required schema formalization, not a cross-branch patch. The pre-checked friction note was accurate.
- The noun-verb dispatch pattern (`cell-return`, `cell-resume`) means `cn cell` automatically works as a group listing command without any additional code (the dispatch.go `GroupMembers` function finds all `cell-*` commands). Clean extension of the existing pattern.
