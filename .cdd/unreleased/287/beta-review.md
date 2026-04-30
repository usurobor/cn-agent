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

### §2.0 Issue Contract

#### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| 1 | `CDD.md` §1.4 γ Phase 1 step 3a inserted (creation snippet, "before α and β are dispatched", α/β do not create) | yes | **met** | CDD.md L235–240 — step 3a with snippet `git switch -c cycle/{N} origin/main && git push -u origin cycle/{N}`; gating sentence "branch must exist on `origin` *before* α and β are dispatched"; cross-ref to §1.4 γ dispatch prompt format. Step 4 (polling) retargeted to `origin/cycle/{N}`; step 5 prompts gain `Branch:` line. |
| 2 | `CDD.md` §1.4 α step 5a inserted; α never creates; surfaces dispatch error if branch missing | yes | **met** | CDD.md L312–315 — "5a. **Check out the cycle branch** named in the dispatch prompt: `git fetch origin cycle/{N} && git switch cycle/{N}`. **α never creates a branch.** ... surface the dispatch error to γ and refuse to invent a name". Step 6 also rewritten to drop "branch" from the implement list. Step 1 updated: "The prompt names the cycle branch explicitly". |
| 3 | `CDD.md` §1.4 β step 3 polling target single-named-branch; β never creates; refuses harness branches | yes | **met** | CDD.md L335–375 — step 3 fully rewritten: single-named-branch poll (`git rev-parse --verify origin/cycle/{N}`), N=10 reachability re-probe embedded in shell snippet, "**β never creates a branch.**" stated, "**β refuses harness pre-provisioned per-role branches**" stated, refusal-as-status-report retained. |
| 4 | `CDD.md` §1.4 γ dispatch prompt format `Branch:` line for α and β; γ prompt unchanged | yes | **met** | CDD.md L290 (α prompt) + L298 (β prompt) — `Branch: cycle/<number>` line added to both. γ prompt at L302–306 unchanged. Parameters paragraph at L313 rewritten to explain why α/β prompts include `Branch:` and γ does not. |
| 5 | `CDD.md` §4.2 canonical = `cycle/{N}`; legacy three patterns warn-only; frozen historical not retroactively renamed | yes | **met** | CDD.md §4.2 (L568–598) — canonical block: `cycle/{N}` with one-cycle-one-branch rule; Ownership clause names γ; Legacy shapes block lists `{agent}/{version}-{issue}-{scope}`, `{agent}/{issue}-{scope}`, `claude/{slug}-{rand}` — all warn-only; "Frozen historical branches under those shapes are not retroactively renamed" present verbatim. |
| 6 | `CDD.md` §4.3 ownership shifts to γ; collapsed checks (branch not yet exist, no stalled cycle dir, scope declared, base SHA known, issue open) | yes | **met** | CDD.md §4.3 (L600–608) — first sentence names γ as owner; five collapsed checks listed verbatim. Pre-#287 historical note retained for α's pre-review-gate row. |
| 7 | `CDD.md` §Tracking — Cycle branch existence + Cycle branch head SHA rows; legacy glob warn-only | yes | **met** | CDD.md §Tracking polling table (L189–197) — "Cycle branch existence (γ pre-flight)" row with `git rev-parse --verify origin/cycle/{N}`; "Cycle branch head SHA" row with `git rev-parse origin/cycle/{N}`; "Legacy `'origin/claude/*'` glob" row marked warn-only/retrospective; old "New branches" row removed. `.cdd/unreleased/{N}/` row updated to `origin/cycle/{N}`. |
| 8 | `CDD.md` §Tracking — `git fetch --quiet` reliability dependency; N=10 re-probe rule | yes | **met** | CDD.md §Tracking new paragraph "**`git fetch` reliability is an explicit dependency**" (L210–225) — names the failure mode, states **N = 10** rule with sizing, shows the synchronous re-probe shape, names on-failure escalation ("surface the failure to the operator immediately"). Embedded in β step 3 shell snippet via `empty_iters` counter. Embedded in `gamma/SKILL.md` §2.5 Step 3b polling snippet. |
| 9 | `gamma/SKILL.md` §2.5 split into Step 3a + Step 3b; ## Step map gains Step 3a row | yes | **met** | gamma/SKILL.md L201 — section header changed: "### 2.5. Steps 3a–5 — Create the cycle branch, dispatch, unblock without leakage"; Step 3a subsection (L203–230) with γ pre-flight + creation snippet; Step 3b subsection (L232–275) with single-named-branch polling. ## Step map at L77 gains "Step 3a → create the cycle branch" row; "Steps 4–6" row collapsed to "Steps 3b–6". |
| 10 | `alpha/SKILL.md` drops α-creates-branch references; pre-review gate row 1 verifies rebased | yes | **met** | alpha/SKILL.md §2.1 step 2 inserted with full check-out / never-creates / surface-error / harness-switch-off rule. Step renumbering 2→3, 3→4, 4→5, 5→6 mechanical. §2.6 row 1 fully rewritten: "α verifies `origin/cycle/{N}` is rebased onto current `origin/main`" with rebase command. §2.7 push instruction updated to `origin/cycle/{N}`. Transient-row paragraph updated. β grep `α opens\|α creates a branch\|alpha creates\|alpha opens` against alpha/SKILL.md returns 0 hits — confirmed. |
| 11 | `beta/SKILL.md` §1 Role Rule 1 expanded: forbid β creating branch; refuse harness pre-provisioned; new ❌/✅ examples | yes | **met** | beta/SKILL.md §1 (L76–96) — opening "β does not author the fix it judges. If RC is requested, α performs the fix." retained. New bold paragraph "**β never creates a branch.**" with creation forbidden, surface-failure-to-operator rule. New bold paragraph "**β refuses harness pre-provisioned per-role branches as the implementation/review surface.**" with concrete examples (`claude/{slug}-{rand}`, `claude/implement-beta-skill-loading-ZXWKe`). ❌/✅ examples expanded to 4 ❌ + 4 ✅. Closing refusal-as-status-report sentence retained, polling target updated to `origin/cycle/{N}`. |
| 12 | Self-application: cycle/287; γ creates branch *before* dispatch; cycle exemplifies its own rule | yes | **met** | β verified by execution: `origin/cycle/287` existed at β intake (`git rev-parse origin/cycle/287` returned `7533978...` at γ's pre-dispatch state, the head of "γ iteration proposal — after 3.61.0"). α's three commits (`6a897a4`, `3336fe3`, `6f44dbb`) all land on `origin/cycle/287`; α never created a separate branch. The dispatch prompt β received included `Branch: cycle/287` per AC 4. β `git switch`ed from harness branch `claude/cnos-skill-module-x9jTE` to `cycle/287` at intake (this verdict commits land on `cycle/287`, not on the harness branch — exemplifying the new β contract). |

**AC summary:** 12/12 met. Spec text faithfully implements all twelve AC clauses. Peer enumeration in self-coherence.md is exhaustive (4 role-skill peers all updated; 6 lifecycle-skill peers reviewed; 4 intra-doc-repetition facts grepped). The diff coherence is high.

#### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.cdd/skills/cdd/CDD.md` | yes | met | 147 lines net change across §1.4 γ/α/β algorithms, §1.4 dispatch prompt format, §Tracking polling table + reliability paragraph, §4.2 Branch rule, §4.3 Branch pre-flight, Role table |
| `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | yes | met | §2.5 split Step 3a / Step 3b; ## Step map row added |
| `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | yes | met (with σ ride-along) | α-authored: §2.1 step 2 insertion + step renumbering; §2.6 row 1 rewrite + transient-row note; §2.7 push target. **σ-authored (`70ff2b1`):** §2.5 incremental-write discipline + §2.7 readiness-signal-as-separate-commit. The σ portion is out of #287 scope — see F1. |
| `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` | yes | met | §1 Role Rule 1 expanded: never-creates rule, harness-refusal rule, ❌/✅ examples |
| `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` | yes | met | §2.2 polling snippet — canonical glob updated to `origin/cycle/*`; legacy `origin/claude/*` named warn-only/retrospective in inline comment |
| `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` | yes — out of scope | scope drift | σ-authored only (`70ff2b1`): Output Format paragraph rewritten to mandate incremental-write discipline. Not in issue #287's affected file list. See F1. |
| `src/packages/cnos.cdd/skills/cdd/SKILL.md` (frontmatter) | no | omitted but not load-bearing | Issue's ## Work-shape names "cdd/SKILL.md frontmatter" as one of the 6 affected files. The file's frontmatter was not modified in this diff. β grep against `cdd/SKILL.md` finds no obsolete branch-shape reference; the frontmatter exposes loader entrypoint only. **Disposition: A-polish at most** — the canonical authority for branch ownership is `CDD.md` + role skills, which are fully updated. The frontmatter omission does not affect the role contract; flagged only for completeness. |

#### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/287/self-coherence.md` | required (CDD.md §Tracking canonical filenames; α step 7) | yes | 333 lines (pre-readiness) → with §Pre-review gate + §Review-readiness signal at HEAD. CDD Trace through step 7a present. AC-by-AC mapping for all 12 ACs. §Peer enumeration with role-skill + lifecycle-skill + intra-doc tables. §Pre-review gate fills all 10 rows (transient row caveats explicit). §Review-readiness signal explicit with base SHA / head SHA / branch CI deferral / "ready for β" line. Two factual issues (F2 + F3) flagged below. |
| `.cdd/unreleased/287/beta-review.md` | required (this file) | yes | being written incrementally per pass — Pass 1 already committed (`cefb4a8`); Pass 2 in this commit; Pass 3 with findings table to follow. |
| `.cdd/unreleased/287/alpha-closeout.md` | post-merge | n/a | written after β merge per α step 10 |
| `.cdd/unreleased/287/beta-closeout.md` | post-merge | n/a | written after β merge per β step 9 |
| `.cdd/unreleased/287/gamma-closeout.md` | post-PRA | n/a | written after γ PRA per γ Phase 3 |

#### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| Tier 1a — `cdd/SKILL.md`, `CDD.md`, `alpha/SKILL.md` | CDD.md §4.4 | yes (per α §Mode + CDD Trace) | yes | role-surface skills loaded; α stayed in α voice (factual, no dispositions); did not load β / γ peer skills |
| Tier 1b — `issue/`, `design/` | CDD.md §4.4 (per phase) | yes (consulted) | yes | α §Mode names both as consulted; design "not required" justified per CDD Trace step 6a (issue body already encodes design); issue/SKILL.md consulted for AC interpretation |
| Tier 2 — `cnos.eng/skills/eng/*` | bundle per `eng/README.md` (docs/skills bundle) | n/a | n/a | markdown-only diff; no language-specific generator skill applies; α correctly classified |
| Tier 3 — issue-specific | dispatch prompt | n/a (none beyond Tier 1) | n/a | dispatch prompt + issue ## Tier 3 skills both state "none beyond CDD Tier 1 (CDD self-modification only)"; α correctly carried forward |

**β-side validation:** β re-grepped α's claim "alpha/SKILL.md after the patch contains zero hits for `α opens|α creates a branch|alpha creates|alpha opens`" — confirmed. β re-grepped "γ creates `cycle/{N}` from `origin/main`" — found at 7 sites (CDD.md Role table, CDD.md §Tracking ¶ "The cycle's branch...", CDD.md §1.4 γ Phase 1 step 3a, CDD.md §1.4 dispatch prompt format parameters paragraph, CDD.md §4.2 Ownership clause, CDD.md §4.3 first sentence, gamma/SKILL.md §2.5 Step 3a body) — consistent phrasing across all sites. β re-grepped "warn-only / retrospective" — found at the expected sites (CDD.md §Tracking polling table row, CDD.md §Tracking ¶ "Single named branch", CDD.md §4.2 Legacy block, operator/SKILL.md §2.2 inline comment) — consistent.
