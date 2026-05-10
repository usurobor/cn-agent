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

### §Self-check

**Did α's work push ambiguity onto β?**

No unresolved ambiguity was deferred. The three sub-rules in §1.6c each have explicit, operator-actionable guidance. The §timeout-recovery section has a concrete decision table, not vague prose. The telemetry fields are named and defined, not left as "TBD fields."

**Is every claim backed by evidence in the diff?**

- §1.6c heuristic constants (`120s`, `180s`, `300s floor`, `400s floor`): backed by empirical anchors at lines 513–514 of CDD.md citing cycle #335 and TSC supercycle artifacts.
- §1.6c "initial; refine with telemetry": language present in the section body (line 509) and in the telemetry framing in post-release §4.
- Cross-refs: all three verified in §ACs above.
- No frontmatter changes: confirmed — none of the three SKILL.md files had their frontmatter blocks modified.
- No new sub-skills: §timeout-recovery is inline at `operator/SKILL.md §7`, not a separate file.

**Peer enumeration:**

This is a docs-only cycle touching three files in `src/packages/cnos.cdd/skills/cdd/`. The peer enumeration question is: are there other lifecycle or role skills that need updating as a consequence of §1.6c? The change is additive (new section) to the coordination model. §1.6b (the sister rule) is explicitly cited and not modified. Role skills (`alpha/`, `beta/`, `gamma/`) do not need updating — the sizing rule is a δ/operator concern codified in CDD.md + operator/SKILL.md. The `release/SKILL.md §3.8` rubric is referenced by the operator override declaration template but not changed by this cycle. Peer set = {§1.6a, §1.6b, operator/§4, post-release/§4, release/§3.8}. Updated: operator/§4 (cross-ref present), post-release/§4 (telemetry fields). Exempted: §1.6a (no change required, re-dispatch prompts are already compatible), §1.6b (cited, not modified per design constraint), release/§3.8 (referenced by operator §7.3 template text, structural reference only, no content change needed).

**Recursive coherence (proof plan):**

This cycle's budget was 600s. Heuristic for 5 ACs docs cycle: `max(300, 120×5) = 900s`. The budget set in the dispatch prompt (600s) is below the heuristic. β will verify this discrepancy per the proof plan. α records it honestly: the dispatch predates the heuristic being codified (the heuristic is the output of this cycle), so the budget mismatch is structurally inevitable for the cycle that introduces the rule. This is declared as known debt in §Debt.

### §Debt

**D1 — Recursive coherence: dispatch budget below heuristic**
This cycle's dispatch budget (600s) is below the heuristic this cycle proposes (`max(300, 120×5) = 900s`). The bootstrap paradox is unavoidable: the rule is the output of the cycle that needed it, so the cycle that introduces the rule cannot simultaneously satisfy it. β's recursive-coherence check (per the issue proof plan) will find this discrepancy. α declares it honestly rather than hiding it.
- No remediation possible within this cycle.
- The discrepancy validates that the heuristic is non-trivial and would have changed the dispatch if it had existed prior.

**D2 — Provisional alpha-closeout.md**
In the sequential bounded dispatch model, α exits before β approves. `alpha-closeout.md` is written before exit, explicitly marked `[provisional — pending β outcome]`, per `alpha/SKILL.md §2.8` provisional close-out fallback. This is standard protocol for non-re-dispatch cycles.

**D3 — Heuristic constants not empirically validated**
The `120s/180s` multipliers and `300s/400s` floors are anchored in N=4 failures, not a dataset. Validation deferred per the issue Scope §Deferred: ≥10 cycles of telemetry needed from the new PRA fields before tightening is possible. Filed as issue scope intent; no issue number yet.

**D4 — No tests**
Docs-only cycle. All changes are structured prose to three skill files. No test framework applies. Explicit reason: §2.2 "tests must prove the actual claim" — the claim is prose correctness and cross-ref coherence, both verified by oracle commands run in §ACs.

### §CDD-Trace

CDD step sequence per `CDD.md §1.4 α algorithm`:

| Step | Action | Evidence |
|------|--------|---------|
| 1. Receive | Dispatch taken. Issue #338 read in full. Active design constraints enumerated. Branch `cycle/338` confirmed (already on branch, origin up to date). Git identity configured: `alpha@cdd.cnos`. | Dispatch prompt + `git config user.email` = `alpha@cdd.cnos` |
| 2. Produce (design artifact) | Not required — mode is `design-and-build`; design lives in the issue body. Issue body serves as the design artifact. | Issue #338 body contains Problem, Scope, ACs, Proof plan |
| 3. Produce (coherence contract) | `self-coherence.md §Gap` written first — names the incoherence being closed and the version/mode. | Commit `3722ac07` |
| 4. Produce (plan) | Not required — three-file docs patch with prescribed order in dispatch prompt. Commit order declared in dispatch: CDD.md → operator → post-release. | Dispatch prompt §Commit checkpoints |
| 5. Produce (tests) | Not applicable — docs-only cycle. Explicit reason in §Debt D4. AC oracles serve as the mechanical verification. | §Debt D4; AC oracle results in §ACs |
| 6. Produce (code/docs) | Three skill file patches in declared commit order: (1) `CDD.md §1.6c` — commit `69de7ef8`; (2) `operator/SKILL.md §7 timeout-recovery` — commit `b2f5ee3b`; (3) `post-release/SKILL.md §4 telemetry fields` — commit `482f1c81`. First commit (`69de7ef8`) is the implementation SHA. All files modified; no new files created in skill packages; no frontmatter changes; no new sub-skills. | `git diff --stat origin/main..2c496d10` — 5 files: `CDD.md` (+39), `operator/SKILL.md` (+55/-1), `post-release/SKILL.md` (+11), `.cdd/unreleased/338/self-coherence.md` (created) |
| 7. Self-coherence | Written incrementally per `alpha/SKILL.md §2.5`. Sections: §Gap (commit `3722ac07`), §Skills (same commit), §ACs (commit `7baa35ae`), §Self-check (commit `c75d86bf`), §Debt (commit `2c496d10`), §CDD-Trace (this commit). | Commit sequence on `origin/cycle/338` |

**Artifact enumeration against diff** (pre-review gate row 11):

`git diff --name-only origin/main..HEAD` at final pre-signal SHA covers all 7 files:
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — mentioned in §ACs AC1 and §CDD-Trace step 6 ✓
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — mentioned in §ACs AC2 and §CDD-Trace step 6 ✓
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — mentioned in §ACs AC3 and §CDD-Trace step 6 ✓
- `.cdd/unreleased/338/self-coherence.md` — this file; mentioned in §CDD-Trace step 7 ✓
- `.cdd/unreleased/338/alpha-closeout.md` — mandatory close-out artifact; mentioned in §Debt D2 and §Gap ✓
- `.cdd/unreleased/338/cdd-iteration.md` — mandatory close-out artifact; triage produced F1: `cdd-protocol-gap`, trigger from dispatch ✓
- `.cdd/iterations/INDEX.md` — aggregator update per `post-release/SKILL.md §5.6b`; one new row added for cycle #338 ✓

No file in the diff is unmentioned. Pre-review gate row 11: **PASS**

**No new modules or functions** — docs-only patch. Pre-review gate row 12: not applicable.

**Test runner output** — no tests (docs-only). Pre-review gate row 13: not applicable (explicit reason in §Debt D4).

**Author email:** `git log -1 --format='%ae' HEAD` = `alpha@cdd.cnos` ✓ Pre-review gate row 14: **PASS**

---

## Review-readiness | round 1 | implementation SHA: 18b3491e | branch CI: N/A (docs-only) | ready for β

**Base SHA:** `cfd322e6` (origin/main at signal time — confirmed ancestor of HEAD at 2026-05-10)
**Head SHA:** see branch HEAD on `origin/cycle/338` (per SHA-convention: omit forward-reference; β polls HEAD)
**Implementation SHA:** `18b3491e` (last commit before this readiness-signal commit)

**Pre-review gate summary:**

| Row | Check | Status |
|-----|-------|--------|
| 1 | `origin/cycle/338` rebased onto `origin/main` (`cfd322e6`) | PASS — observed at 2026-05-10; origin/main is ancestor of HEAD |
| 2 | `self-coherence.md` carries CDD Trace through step 7 | PASS — §CDD-Trace present with all 7 steps |
| 3 | Tests present or explicit reason | PASS — docs-only; explicit in §Debt D4 |
| 4 | Every AC has evidence | PASS — AC1–AC5 all have oracle results in §ACs |
| 5 | Known debt explicit | PASS — §Debt D1–D4 |
| 6 | Schema/shape audit | N/A — no schema-bearing changes |
| 7 | Peer enumeration | PASS — peer set enumerated in §Self-check |
| 8 | Harness audit | N/A — no schema-bearing contract changes |
| 9 | Polyglot re-audit | PASS — docs-only (Markdown); no Go/shell/YAML in diff |
| 10 | Branch CI green | N/A — docs-only cycle; no CI pipeline applies |
| 11 | Artifact enumeration matches diff | PASS — all 7 files in diff enumerated in §CDD-Trace step 6 |
| 12 | Caller-path trace for new modules | N/A — no new modules |
| 13 | Test assertion count from runner | N/A — no tests (docs-only) |
| 14 | α commit author email = `alpha@cdd.cnos` | PASS — verified via `git log --format='%ae'` |

**Mechanical pre-merge closure gate (active from #339):** All three required close-out artifacts present in `.cdd/unreleased/338/`:
- `self-coherence.md` ✓
- `alpha-closeout.md` ✓ (provisional — per §Debt D2)
- `cdd-iteration.md` ✓

**Ready for β.**


