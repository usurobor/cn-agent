## Self-Coherence — cycle #338

### §Gap

**Issue:** #338 — cdd: Add §1.6c — dispatch sizing, prompt scope, and commit checkpoints
**Version / mode:** docs-only (`design-and-build`; §2.5b disconnect path)
**Branch:** `cycle/338`
**Gap being closed:** `CDD.md §1.6` had no rule for scaling the initial dispatch timeout budget with work complexity, no prompt-scope guidance matching load to task, and no commit-checkpoint mandate. This created a recoverable-only-by-luck failure class: agents SIGTERMed without committing lose all in-progress work. Pattern observed in N=4 instances (cnos #335, TSC supercycle 3/5 re-dispatches).

This cycle also validates its own heuristic: 5 ACs × 120s + 300s floor = 900s. Budget set to 600s per the dispatch prompt — this is below the heuristic. β will apply recursive coherence check per the proof plan. Known debt: budget discrepancy declared in §Debt.

### §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle; §1.6 coordination model is the change target
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface and load order (this file)

**Tier 2 (always-applicable eng):**
- `src/packages/cnos.eng/skills/eng/writing` — prose patches to skill files (docs-only cycle)

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-program/frontmatter coherence on the three modified files (constraint: no frontmatter changes)

**Active design constraints loaded:**
- No frontmatter changes to any of the three modified SKILL.md files
- No new sub-skills — §timeout-recovery lives inline in operator/SKILL.md
- Heuristic constants marked "initial; refine with telemetry"
- Cross-refs are normative: §1.6c ↔ §1.6b ↔ operator §timeout-recovery ↔ post-release §4
- Empirical citations reproducible (close-out paths)
- No retroactive changes to §1.6a or §1.6b

### §ACs

**AC1: §1.6c added to CDD.md with three sub-rules**

Oracle: `grep -nE "^### §1\.6c|^#### \(a\)|^#### \(b\)|^#### \(c\)" src/packages/cnos.cdd/skills/cdd/CDD.md`

Result (4 matches in correct order):
```
496:### §1.6c Initial dispatch sizing, prompt scope, and commit checkpoints
500:#### (a) Timeout budget heuristic
516:#### (b) Prompt scope
525:#### (c) Commit checkpoints
```

Evidence: commit `69de7ef8` — `docs(338): CDD.md §1.6c — dispatch sizing, prompt scope, commit checkpoints`. Section present; three sub-rules present; numerical heuristic in (a) is explicit (`max(300s, 120 × ac_count)` / `max(400s, 180 × ac_count)`); prompt-scope guidance in (b) cites §1.6b; commit-checkpoint instruction in (c) is quotable text. **AC1: PASS**

---

**AC2: operator/SKILL.md gains timeout-recovery section**

Oracle: `grep -nE "Timeout recovery|timeout recovery" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` → ≥1 match; section content includes `git status --short` and `git stash list`.

Result:
```
268:## 7. Timeout recovery
280:git status --short
287:git stash list
```

Evidence: commit `b2f5ee3b` — `docs(338): operator/SKILL.md §7 timeout-recovery`. Section present at §7; `git status --short` at line 280; `git stash list` at line 287; decision tree present (§7.2 table with 5 rows); override declaration template present (§7.3). **AC2: PASS**

---

**AC3: PRA telemetry fields added to post-release/SKILL.md §4**

Oracle: `grep -nE "dispatch_seconds_budget|dispatch_seconds_actual|commit_count_at_termination" src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` → 3 matches in or near §4.

Result:
```
90: | Cycle | `dispatch_seconds_budget` | `dispatch_seconds_actual` | `commit_count_at_termination` |
93: - `dispatch_seconds_budget` — the timeout budget...
94: - `dispatch_seconds_actual` — wall-clock seconds...
95: - `commit_count_at_termination` — number of commits...
```

All three field names present at lines 90–95, within §4 (Review Quality). Framing: "optional initially; mandatory after ~10 cycles of data accumulate". Evidence: commit `482f1c81`. **AC3: PASS**

---

**AC4: Cross-references between §1.6c, §1.6b, operator §timeout-recovery, post-release §4 are intact**

Oracle: manual review — three explicit cross-refs in §1.6c body.

Evidence (all in `CDD.md` §1.6c, commit `69de7ef8`):
1. Line 498: "Sister rule to §1.6b (which covers re-dispatch)" — §1.6c → §1.6b ✓
2. Line 509: "Per-cycle `dispatch_seconds_budget` and `dispatch_seconds_actual` are recorded in the PRA telemetry fields (post-release §4)" — §1.6c(a) → post-release §4 ✓
3. Lines 531, 533: "`operator/SKILL.md §timeout-recovery`" cited twice in §1.6c(c) — §1.6c(c) → operator §timeout-recovery ✓

**AC4: PASS**

---

**AC5: Empirical anchor cited in §1.6c body**

Oracle: `grep -E "cycle #335|tsc-supercycle|usurobor/tsc|600s.*0 commits|18 files" src/packages/cnos.cdd/skills/cdd/CDD.md` → matches in §1.6c.

Result:
```
513: - **cnos cycle #335** (close-out artifacts at `.cdd/releases/docs/2026-05-09/335/`): 9 ACs, 600 s budget, 0 commits at SIGTERM, 18 files recovered from worktree by operator.
514: - **TSC supercycle** (close-outs at `usurobor/tsc:.cdd/releases/{0.5.0,0.6.0,0.7.0}/{25,24,26}/`): 3 of 5 α close-out re-dispatches failed under full α skill load.
```

Both empirical anchors present with reproducible close-out paths. Anyone can navigate to `.cdd/releases/docs/2026-05-09/335/` in the cnos repo or to `usurobor/tsc:.cdd/releases/{0.5.0,0.6.0,0.7.0}/{25,24,26}/` in the TSC repo. **AC5: PASS**



