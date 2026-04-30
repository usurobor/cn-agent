# β Review — #287

**Cycle:** #287 — γ creates the cycle branch — α and β only check out `cycle/{N}`
**Branch:** `cycle/287`
**Reviewer:** β (`beta@cdd.cnos`)
**Issue:** [#287](https://github.com/usurobor/cnos/issues/287)

This file is appended round-by-round and pass-by-pass per `review/SKILL.md` Output Format (incremental write discipline) and `CDD.md` §1.4 large-file authoring rule. Each pass within a round is a separate commit on `origin/cycle/287`.

---

## Round 1 — Verdict: REQUEST CHANGES

**Round head SHA reviewed:** `6f44dbb64d5c39f80b82d294943170b91ab715ce`
**Base:** `origin/main` = `a8e67b776fd55f58d07fc474283606cf671990c3`
**Branch CI state:** deferred — markdown-only diff, no cycle-branch CI surface in this repo (consistent with α's §Pre-review gate row 10 deferral). Re-validated by β: `.github/workflows/build-verify.yml` runs on `main` and PRs only; `.github/workflows/release.yml` runs on tag push. Cycle-branch pushes do not trigger CI. Verdict relies on direct diff coherence.
**Merge instruction:** **deferred until findings F1, F2, F3 resolved.** On approval after fix: `git merge cycle/287` into `main` with `Closes #287` in the merge commit message; tag bare version per `release/SKILL.md`.

### §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | **no** | α's §Known debt #5 + §Pre-review gate row 1 evidence both claim `70ff2b1` is `origin/main`. main is at `a8e67b7`; `70ff2b1` is reachable only from `origin/cycle/287`. See F2. |
| Canonical sources/paths verified | yes | All cross-references in the diff resolve: §1.4 γ Phase 1 step 3a, §1.4 α step 5a, §1.4 β step 3, §1.4 dispatch prompt format, §4.2, §4.3 — all targets exist in this diff. `gamma/SKILL.md` §2.5 Step 3a/3b headers correct; `alpha/SKILL.md` §2.6 row 1 + §2.7 cross-refs resolve. |
| Scope/non-goals consistent | **no** | `70ff2b1` (sigma) modifies `review/SKILL.md` (not in #287's affected file list) and inserts §2.5 incremental-write discipline + §2.7 readiness-signal-as-separate-commit into `alpha/SKILL.md` (not covered by ACs 1–12). Issue #287's ## Work-shape names exactly 6 affected files; `review/SKILL.md` is not one of them. See F1. |
| Constraint strata consistent | yes | Hard rules ("γ creates", "α never creates", "β never creates", "β refuses harness pre-provisioned per-role branches") stated as gates without exception-backed escape hatches. Legacy shapes named warn-only, not exception-backed. |
| Exceptions field-specific/reasoned | yes | Legacy branch shapes block names three specific patterns and the cleanup condition ("frozen historical branches not retroactively renamed"). Polling-glob legacy retention scoped to "retrospective tracking only." |
| Path resolution base explicit | yes | Branch references uniformly `origin/cycle/{N}` (remote, repo-root-implicit). Polling targets named with full `origin/cycle/{N}` prefix throughout. No path-base ambiguity in the diff. |
| Proof shape adequate | yes — N/A justified | Markdown-only spec diff; no executable surface to prove against. AC 12 self-application *is* the integration test: γ created `cycle/287` before α dispatch; α `git switch`ed onto it from harness branch `claude/alpha-cdd-skill-1aZnw`; α never created a branch. Three α commits + one σ commit on `cycle/287`; zero α commits on `claude/alpha-cdd-skill-1aZnw`; verified by `git log` author lines. |
| Cross-surface projections updated | yes | Role table (CDD.md L94–98), §Tracking polling table, §1.4 γ/α/β algorithms, §1.4 dispatch prompt format, §4.2 + §4.3 + `gamma/SKILL.md` §Step map row + §2.5 + `alpha/SKILL.md` §2.1/§2.6/§2.7 + `beta/SKILL.md` §1 + `operator/SKILL.md` §2.2 — all consistent. Intra-doc repetition grep verified by α (§Peer enumeration §Intra-doc repetition); spot-checked by β: "γ creates `cycle/{N}` from `origin/main` before dispatch" appears with consistent phrasing across 7+ sites. |
| No witness theater / false closure | yes | Spec changes name a rejection mechanism (γ pre-flight aborts on existing branch / stalled cycle dir; α surfaces dispatch error if `origin/cycle/{N}` missing; β refuses harness branch). The new contract is enforceable by inspection, not just prose. |
| PR body matches branch files | n/a | Artifact-channel protocol per #283 — no PR. self-coherence.md (the artifact-channel equivalent) matches branch state with one exception per F2 (Known debt #5 misidentifies `70ff2b1` as main). |

**Result:** Two **no** rows (status truth, scope) — both raise findings. Phase 2 proceeds; verdict will weigh the contract findings against AC coverage.
