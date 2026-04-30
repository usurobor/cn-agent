# Self-Coherence — #287

**Cycle:** #287 — γ creates the cycle branch — α and β only check out `cycle/{N}`
**Branch:** `cycle/287` (γ pre-provisioned per AC 12 self-application)
**Mode:** MCA — concrete spec patch closes a recurring friction class.
**Author:** α (`alpha@cdd.cnos`)
**Issue:** [#287](https://github.com/usurobor/cnos/issues/287)

## Gap

CDD's branch-creation step is α-owned today. α opens a branch using the harness-given name (e.g. `claude/alpha-cdd-skill-1aZnw`, scope-words + random suffix, no issue number). γ and β must glob-match `origin/claude/*` to discover it, and the discovery silently fails across cycles:

- γ's polling glob templates assume `*<N>*` encoding the issue number; the harness encodes scope words and a random suffix. The glob matches zero branches.
- β receives a *different* harness branch with α-style instructions, refuses (per `beta/SKILL.md` §1), and then has nowhere canonical to land its review verdict — β's verdict lands on the harness β-branch and α has to cherry-pick.
- γ has no canonical way to name a single branch in the dispatch prompt; the prompt has to say "the branch α opens" and rely on downstream discovery.

The fix: γ creates `cycle/{N}` from `origin/main` *before* dispatch, and dispatch prompts name the branch explicitly. α and β never create branches — they `git switch cycle/{N}` and refuse to invent or accept a per-role pre-provisioned branch.

## Mode

- **Mode:** MCA (skill / spec patch landing in the same cycle that surfaces the gap, as committed in 3.61.0 PRA §7 next-MCA).
- **Active skills (Tier 1):** `cdd/SKILL.md`, `cdd/CDD.md`, `cdd/alpha/SKILL.md`, `cdd/issue/SKILL.md` (read for AC interpretation), `cdd/design/SKILL.md` (consulted; design "not required" — the change is a contract surface rewrite with the design already named in the issue body).
- **Tier 2:** none — diff is markdown-only across `src/packages/cnos.cdd/skills/cdd/`. No language-specific generator skill applies.
- **Tier 3:** none beyond CDD Tier 1 (per dispatch prompt and issue ## Tier 3 skills).

## Active design constraints (from issue)

- The cycle branch `cycle/{N}` is the canonical coordination surface during in-version work (per #283 §Tracking).
- γ owns branch creation; α and β own only check-out and commit-to-the-branch.
- Dispatch prompts name the branch explicitly; agents do not invent names.
- One cycle = one branch = one named target for all polling.

## Self-application (AC 12)

This cycle uses `cycle/287`. γ created `origin/cycle/287` before dispatching α. α `git switch`ed to it from the harness-given branch (`claude/alpha-cdd-skill-1aZnw`) per the new α step 5a. The cycle exemplifies its own rule.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | 3.61.0 PRA §4b cycle-iteration trigger named #287 as next-MCA; γ filed and dispatched. |
| 1 Select | — | cdd | Selected by 3.61.0 PRA §7 commitment + §3.3 assessment-default rule. |
| 2 Branch | `cycle/287` | cdd | γ pre-created `origin/cycle/287` per AC 1 (self-application); α `git switch cycle/287` per AC 2. α did **not** create the branch. |
| 3 Bootstrap | — | cdd | Not required — no version directory; small markdown-only spec change inside an existing skill package. |
| 4 Gap | this file §Gap | cdd | Branch-creation friction class named (recurrence in #274 / #283 R1 F1). |
| 5 Mode | this file §Mode | cdd | MCA, no Tier 3 skills. |
| 6a Design | — | cdd, design | Not required — the issue body already encodes the design (12 numbered ACs against named files). |
| 6b Plan | this file §Implementation plan | cdd, plan (consulted) | Not required as a separate artifact — the per-AC implementation plan below is sufficient for a markdown-only diff with a flat AC structure. |
| 6c Tests | — | cdd | Not required — markdown / spec-only diff. There is no executable surface in this change. The cycle itself is the integration test (AC 12 self-application: the cycle uses `cycle/287`, exemplifying its own rule). |
| 6d Code | diff against `cycle/287` | cdd | Spec text patches across CDD.md + 4 role skills. |
| 6e Docs | diff | cdd | Same as 6d — the spec text *is* the doc. |
| 7 Self-coherence | this file | cdd | AC-by-AC mapping below. |
| 7a Pre-review | this file §Pre-review gate | cdd | Pre-review gate run; review-readiness signaled at the bottom of this file. |

## Implementation plan

For a flat AC structure across one canonical file (`CDD.md`) plus 4 role skill files, the plan is literally per-AC. Order:

1. **AC 1, 7, 8 (CDD.md §Tracking + §1.4 γ Phase 1 step 3a)** — the heaviest single change. Lands first; downstream changes reference it.
2. **AC 2 (CDD.md §1.4 α step 5a)** — depends on AC 1's branch-creation contract being named.
3. **AC 3 (CDD.md §1.4 β step 3 polling target)** — depends on AC 1's `cycle/{N}` canonical name.
4. **AC 4 (CDD.md §1.4 dispatch prompt format Branch: line)** — sits alongside the algorithms and references `cycle/{N}` named in AC 1.
5. **AC 5, 6 (CDD.md §4.2 / §4.3)** — branch rule + pre-flight; ownership shifts to γ.
6. **AC 9 (gamma/SKILL.md §2.5 split + step map row)** — mirrors AC 1 + AC 7's polling rewrite at the role-skill layer.
7. **AC 10 (alpha/SKILL.md drops branch creation; pre-review gate row 1)** — mirrors AC 2 at the α role-skill layer.
8. **AC 11 (beta/SKILL.md §1 Role Rule 1 expansion)** — mirrors AC 3 at the β role-skill layer.
9. **operator/SKILL.md** — small polling-snippet update naming the legacy glob as warn-only retrospective form (per AC 5 + AC 7 cascading downstream).

## Acceptance Criteria — evidence

### AC 1 — `CDD.md` §1.4 γ algorithm Phase 1 step 3a

**Claim:** γ creates `cycle/{N}` from `origin/main` and pushes it before any dispatch.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §1.4 γ algorithm Phase 1, step 3a inserted with the exact creation snippet (`git switch -c cycle/{N} origin/main && git push -u origin cycle/{N}`) and explicit "before α and β are dispatched" + "α and β do not create branches" gating.
- Step 4 (begin polling) updated to target `origin/cycle/{N}` directly, not glob.
- Step 5 (dispatch prompts) updated to note the `Branch: cycle/{N}` line is included for both prompts.
- The `## Related artifacts` rolled-back commit `eb48e17` referenced in the issue body is unavailable in this repo (rolled back from γ's local environment); spec text was reconstructed from the issue body's AC text + 3.61.0 PRA §4b/§7 next-MCA notes.

**Self-application (AC 12 cross-link):** This cycle proved AC 1 by execution — γ created `origin/cycle/287` before dispatching α; α `git switch`ed onto it from the harness-given `claude/alpha-cdd-skill-1aZnw` branch.

### AC 2 — `CDD.md` §1.4 α algorithm step 5a

**Claim:** α `git switch cycle/{N}` from the dispatch prompt; α never creates a branch; if `origin/cycle/{N}` does not exist, α surfaces the dispatch error and refuses to invent a name.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §1.4 α algorithm: step 5a inserted between the existing step 5 (read the issue) and step 6 (implement). Step 5a explicitly says "**α never creates a branch.**" with the surface-to-γ rule and the harness-branch switch-off instruction.
- α step 1 updated to note the prompt names the cycle branch explicitly.
- α step 6 (formerly "Implement: branch, code, tests, self-coherence") rewritten to drop "branch" — α now only commits to the branch γ created.

**Self-application:** α (this session) executed step 5a verbatim — switched off `claude/alpha-cdd-skill-1aZnw` (harness pre-provisioned) onto `cycle/287` before any code change.

### AC 3 — `CDD.md` §1.4 β algorithm step 3 polling target

**Claim:** β polls `git rev-parse origin/cycle/{N}` directly (single named target, no glob); β never creates a branch; β refuses harness pre-provisioned per-role branches.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §1.4 β algorithm step 3 fully rewritten:
  - polling target is `origin/cycle/{N}` directly with `git rev-parse --verify origin/cycle/{N}` reachability probe at intake
  - reference shell snippet uses single named-branch transition loop (no glob, `comm -13` removed since there is only one branch), with the §Tracking 10-empty-iters reachability re-probe rule embedded
  - "**β never creates a branch.**" stated explicitly
  - "**β refuses harness pre-provisioned per-role branches and refuses any instruction to develop, commit, or implement on them.**" added (folded into the step 3 narrative; expanded fully in `beta/SKILL.md` §1 per AC 11)
- β step 1 updated to note the prompt names the cycle branch explicitly.

### AC 4 — `CDD.md` §1.4 γ dispatch prompt format `Branch:` line

**Claim:** α and β prompts gain a `Branch: cycle/<number>` line; γ prompt does not; the parameters paragraph explains why.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §1.4 γ dispatch prompt format:
  - α prompt: `Branch: cycle/<number>` added between `Issue:` and `Tier 3 skills:` lines
  - β prompt: `Branch: cycle/<number>` added between `Issue:` and the closing fence
  - γ prompt: no `Branch:` line (γ owns its own session and creates the branch as part of Phase 1)
- Parameters paragraph rewritten to explain that α and β prompts include `Branch:` because γ creates `cycle/{N}` from `origin/main` before dispatch, and α/β do not invent or glob to discover the branch.
- The dispatch-prompt format used in this very dispatch matches the new format (`Branch: cycle/287`), exemplifying AC 12 self-application.

### AC 5 — `CDD.md` §4.2 Branch rule canonical format `cycle/{N}`

**Claim:** Canonical format becomes `cycle/{N}`. Legacy shapes (`{agent}/{issue}-{scope}`, harness-encoded `claude/{slug}-{N}-{rand}`) are warn-only. Frozen historical branches not retroactively renamed.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §4.2 fully rewritten:
  - canonical format named: `cycle/{N}` with one-cycle-one-branch-one-named-target rule
  - ownership clause names γ + the cross-references to §1.4 γ step 3a, alpha/SKILL.md, beta/SKILL.md §1
  - legacy shapes block lists the three legacy patterns (`{agent}/{version}-{issue}-{scope}`, `{agent}/{issue}-{scope}`, `claude/{slug}-{rand}`) and marks them warn-only
  - "Frozen historical branches under those shapes are not retroactively renamed" stated verbatim
  - polling glob `'origin/claude/*'` flagged as retrospective-only
- Issue non-goal "Retroactively renaming historical `claude/*` branches; legacy shapes warn-only forever" honored: existing illustrative `claude/*` references in `review/SKILL.md` line 73 (`origin/claude/cdd-X` example), `post-release/SKILL.md` line 380 (PRA template Branch field), `release/SKILL.md` line 178 (tag-naming convention rationale) are intentionally left as frozen historical context.

### AC 6 — `CDD.md` §4.3 Branch pre-flight ownership shifts to γ

**Claim:** Ownership of pre-flight shifts from α to γ. Collapsed checks: branch does not yet exist, no stalled `.cdd/unreleased/{N}/` on main, scope declared in issue body, base SHA known, issue open.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §4.3 fully rewritten:
  - first sentence names γ as the owner: "γ runs branch pre-flight before creating `cycle/{N}` in §1.4 γ algorithm Phase 1 step 3a."
  - the five collapsed checks listed verbatim from the AC text (branch does not yet exist; no stalled `.cdd/unreleased/{N}/` on main; scope declared; base SHA known; issue open)
  - explicit "Pre-#287 ownership rested with α" historical note + "α's pre-review gate retains its own row (`alpha/SKILL.md` §2.6 row 1) for verifying that `origin/cycle/{N}` is rebased onto current `origin/main` at review-readiness time, but α no longer creates the branch and does not own pre-flight at branch creation"
- Cross-referenced in `gamma/SKILL.md` §2.5 Step 3a (per AC 9).

### AC 7 — `CDD.md` §Tracking polling rows

**Claim:** "Cycle branch existence (γ pre-flight)" + "Cycle branch head SHA" rows replace the old `New branches` glob row. Polling becomes `git rev-parse origin/cycle/{N}` directly. Legacy `'origin/claude/*'` glob retained only for retrospective tracking.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §Tracking polling-query table fully rewritten:
  - old row "New branches | — | — | `git fetch --quiet origin && git branch -r --list 'origin/claude/*'`" — **removed**
  - new row "Cycle branch existence (γ pre-flight) | — | — | `git fetch --quiet origin cycle/{N} && git rev-parse --verify origin/cycle/{N}`" — added with γ-runs-pre-flight note
  - new row "Cycle branch head SHA | — | — | `git fetch --quiet origin cycle/{N} && git rev-parse origin/cycle/{N}`" — added; old row updated to reference single named branch
  - new row "Legacy `'origin/claude/*'` glob | — | — | `git branch -r --list 'origin/claude/*'` is **warn-only / retrospective**: retained for tracking historical cycles whose branches predate the `cycle/{N}` rule. New cycles do not poll this glob." — added explicitly
  - `.cdd/unreleased/{N}/` directory state row updated to reference `origin/cycle/{N}` (not generic `origin/{branch}`)
- The pre-#287 paragraphs ("Synchronous-baseline-pull is a precondition...", "Branch-glob templates must not assume harness-encoded issue numbers...") rewritten to:
  - direct the per-cycle baseline at `origin/cycle/{N}` rather than the broad `'origin/claude/*'` glob
  - explicitly name the new "Single named branch — no globs for new cycles" rule with the legacy glob marked as retrospective-only
  - cite the *Derives from #274 + #283 R1 F1* lineage and the #287 closure ("γ creates `cycle/{N}` before dispatch and the dispatch prompt names it explicitly, so polling targets a single named branch instead of a glob")
- The `cycle/287` polling targets used by α and γ in this very cycle — `git rev-parse origin/cycle/287` — exemplify the new shape (AC 12).

### AC 8 — `CDD.md` §Tracking `git fetch --quiet` reliability dependency

**Claim:** Name the `git fetch --quiet` reliability layer as an explicit dependency. Add rule: on N=10 successive empty polling iterations, do a synchronous reachability re-probe (`git fetch --verbose origin {branch}` + check stderr) and surface to operator on failure.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §Tracking new paragraph "**`git fetch` reliability is an explicit dependency**" added immediately after the "Single named branch" paragraph:
  - names the failure mode: `--quiet` silences network/auth flakes, fetch returns 0 with stale data, transition loop drops every commit during the flake window
  - states the **N = 10** rule explicitly with "≈ 10 minutes at 60s interval" sizing
  - shows the synchronous re-probe shape: `git fetch --verbose origin cycle/{N} 2>&1 | tee /tmp/cycle-{N}-fetch.log` + the four diagnostic patterns to scan for in stderr
  - states the on-failure escalation: "**surface the failure to the operator immediately** — the role cannot autonomously detect cycle progression with a broken transport"
  - cites *Derives from #283 β observation #3 / 3.61.0 PRA §4b*
- The β step 3 reference shell snippet (`CDD.md` §1.4 β algorithm) embeds the same N=10 rule via `empty_iters` counter — when 10 iterations pass with no transition, the loop runs `git fetch --verbose origin cycle/{N}` and emits `reachability-fail: cycle/{N}` if it errors.
- `gamma/SKILL.md` §2.5 Step 3b polling snippet embeds the same N=10 rule.

### AC 9 — `gamma/SKILL.md` §2.5 split into Step 3a + Step 3b

**Claim:** §2.5 split into Step 3a (create branch, with creation snippet and γ-owned pre-flight) and Step 3b (subscribe + dispatch). Polling snippet uses single named-branch `git rev-parse origin/cycle/{N}` directly. Step map row added for Step 3a.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5:
  - section header changed: "### 2.5. Steps 3a–5 — Create the cycle branch, dispatch, unblock without leakage" (was "### 2.5. Steps 3–5 — Dispatch and unblock without leakage")
  - **#### Step 3a — Create the cycle branch** new subsection:
    - cross-references `CDD.md` §Tracking + §4.2
    - lists γ-owned branch pre-flight checks verbatim from `CDD.md` §4.3
    - shows the creation snippet (`git switch -c cycle/<N> origin/main && git push -u origin cycle/<N>`) with both pre-flight assertions inline (branch does not yet exist; no stalled `.cdd/unreleased/<N>/` on main)
    - states the gating rule: "The branch must exist on `origin` before dispatch. α and β never create branches"
  - **#### Step 3b — Subscribe and dispatch α and β** subsection:
    - polling snippet rewritten to single named branch (`git rev-parse origin/cycle/<N>`), no glob, no associative `prev_head` array since there is exactly one branch
    - reachability re-probe with N=10 empty-iteration counter embedded
    - α prompt + β prompt both gain `Branch: cycle/<N>` line
    - rules paragraph adds "include the explicit `Branch: cycle/<N>` line so α and β never have to invent or glob-discover the branch"
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` ## Step map gains a new row: "Step 3a → create the cycle branch (`§2.5` Step 3a — `cycle/{N}` from `origin/main`, with γ-owned branch pre-flight per `CDD.md` §4.3, before any dispatch)". The "Steps 4–6" row collapsed to "Steps 3b–6 → subscribe, dispatch + unblocking (`§2.5` Step 3b)".

### AC 10 — `alpha/SKILL.md` drops α-creates-branch references; pre-review gate row 1

**Claim:** α/SKILL.md drops references to α creating a branch; pre-review gate row 1 (branch rebased) becomes "α verifies `origin/cycle/{N}` is rebased onto current `main`" rather than "α rebases own branch."

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.1 Dispatch intake: step 2 inserted: "**check out the cycle branch named in the dispatch prompt.** [...] **α never creates a branch.** If `origin/cycle/{N}` does not exist, surface the dispatch error to γ and refuse to invent a name [...] If the harness placed α on a per-role pre-provisioned branch (e.g. `claude/{slug}-{rand}`), switch off it before any code change". Step numbering updated (old steps 2–5 became 3–6).
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 Pre-review gate row 1 fully rewritten: "**α verifies `origin/cycle/{N}` is rebased onto current `origin/main`** (γ created the branch from `origin/main` at dispatch time per `CDD.md` §1.4 γ algorithm Phase 1 step 3a; α does *not* create or rebase a different branch — α only verifies that the cycle branch γ created has not drifted behind `main` while α was working). If `origin/main` advanced, α rebases the cycle branch onto it: `git fetch origin main && git rebase origin/main && git push --force-with-lease origin cycle/{N}`."
- Transient-vs-durable rows paragraph updated: "Rows 1 (cycle branch rebased) and 10 (branch CI green)..."
- §2.7 review-request push instruction updated: "Push the branch to `origin/cycle/{N}`" (was generic `origin/{branch}`); "Roles poll `origin/cycle/{N}`, not `origin/main`"
- α-side branch-creation language search (`grep "α opens|α creates a branch|alpha creates|alpha opens"`) returns zero hits in `alpha/SKILL.md` after the patch.

### AC 11 — `beta/SKILL.md` §1 Role Rule 1 expansion

**Claim:** β/SKILL.md §1 Role Rule 1 expanded to forbid β creating a branch and to refuse harness pre-provisioned per-role branches as the implementation/review surface. New ❌/✅ examples.

**Evidence:**

- `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §1 Keep β independent fully rewritten:
  - opening paragraphs: "β does not author the fix it judges. If RC is requested, α performs the fix."
  - new bold paragraph: "**β never creates a branch.** γ creates `cycle/{N}` from `origin/main` before dispatch [...] β `git switch`es to it and asserts existence with `git rev-parse --verify origin/cycle/{N}`. If the branch does not exist, β surfaces a γ pre-flight failure to the operator — β does not invent a name, does not glob to discover one, and does not create one to 'make progress.'"
  - new bold paragraph: "**β refuses harness pre-provisioned per-role branches as the implementation/review surface.**" with the concrete examples (`claude/{slug}-{rand}`, `claude/implement-beta-skill-loading-ZXWKe`) and the explicit override rule
  - ❌/✅ examples expanded to 4 ❌ + 4 ✅ rows covering: implementation override, harness implementation, β-creates-locally, β-uses-harness-branch (the #283 R1 F1 mode)
  - closing line on refusal-as-status-report retained, polling target updated to `origin/cycle/{N}`

### AC 12 — Self-application: cycle/287

**Claim:** The cycle implementing this AC must itself use `cycle/287`. γ creates the branch *before* dispatching α and β. The cycle exemplifies its own rule.

**Evidence:**

- γ created `origin/cycle/287` from `origin/main` before dispatching α (verified by `git fetch && git rev-parse origin/cycle/287` returning `7533978d57030faf17f1febb2d55a341fca9d8ca` at the start of α's session — that head SHA is γ's pre-dispatch commit "γ iteration proposal — after 3.61.0").
- The α dispatch prompt this session received literally read:
  ```
  You are α. Project: cnos.
  Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md.
  Issue: gh issue view 287
  Branch: cycle/287
  Tier 3 skills: none beyond CDD Tier 1 (CDD self-modification only)
  ```
  This matches the new dispatch-prompt format AC 4 specifies.
- α (this session) executed `git switch cycle/287` from the harness pre-provisioned branch `claude/alpha-cdd-skill-1aZnw` — exemplifying α step 5a (AC 2) and the harness-branch refusal added to alpha/SKILL.md §2.1 step 2 (AC 10).
- All α commits this cycle land on `cycle/287` with author `alpha@cdd.cnos`; nothing on `main` until β's `git merge`.
- β's intake (when it begins) will `git switch cycle/287` from whatever harness branch the operator places it on, exemplifying β/SKILL.md §1 (AC 11). This cycle's `beta-review.md` will land on `cycle/287`, not on a separate harness branch — a structural correction to the pattern that #283 hit.

## Peer enumeration

Per `alpha/SKILL.md` §2.3, when the change touches a family of peers, enumerate the family before claiming closure. The contract change here is "branch creation moves from α to γ; canonical name becomes `cycle/{N}`." Two peer classes apply:

### Role-skill peers

The diff modifies the role contract. Role-skill peer set: `{ alpha/, beta/, gamma/, operator/ }`.

| Peer | Updated? | Surface | Verifying ref |
|------|----------|---------|---------------|
| `gamma/SKILL.md` | ✅ | §2.5 split into Step 3a (create branch) + Step 3b (subscribe + dispatch); ## Step map row 3a added; polling snippet uses single named branch | AC 9 evidence above |
| `alpha/SKILL.md` | ✅ | §2.1 step 2 (check out, never create); §2.6 row 1 (verify rebased, do not create); §2.7 push target named `origin/cycle/{N}` | AC 10 evidence above |
| `beta/SKILL.md` | ✅ | §1 Role Rule 1 expanded — never creates, refuses harness pre-provisioned per-role branches | AC 11 evidence above |
| `operator/SKILL.md` | ✅ | §2.2 polling snippet — canonical glob updated to `origin/cycle/*`; legacy `origin/claude/*` named warn-only/retrospective | issue ## Related artifacts entry "δ branch-cleanup rule continues to sweep merged branches; no functional change but warn-only legacy is named" |

### Lifecycle-skill peers

`alpha/SKILL.md` §2.3 names lifecycle skills (`review/`, `release/`, `post-release/`, `design/`, `plan/`, `issue/`) as a distinct enumeration class from role skills, *because* lifecycle skills encode role contracts operationally and drift downstream when role contracts change (per #283 R1 F2/F3/F4 in `release/`/`post-release/`).

For #287, the contract change is **branch ownership** (γ creates) and **canonical branch name** (`cycle/{N}`). Lifecycle skills were audited:

| Lifecycle skill | Reviewed? | Branch-creation contract impact | Disposition |
|-----------------|-----------|----------------------------------|-------------|
| `issue/SKILL.md` | grep `branch\|create` returns no creation-ownership references | None — issue skill defines issue quality, not branch ownership | No change required |
| `design/SKILL.md` | grep `branch\|create` returns no creation-ownership references | None — design skill defines design-artifact contract, not branch ownership | No change required |
| `plan/SKILL.md` | grep `branch\|create` returns no creation-ownership references | None | No change required |
| `review/SKILL.md` | line 73 — `git merge-base --is-ancestor origin/claude/cdd-X main` example | One illustrative example uses `origin/claude/cdd-X` to demonstrate redirect-to-main when branch is merged. AC 5 says "frozen historical branches under those shapes are not retroactively renamed" and the issue's non-goal explicitly excludes retroactive renaming. | Intentionally retained as legacy/illustrative |
| `release/SKILL.md` | line 178 — `claude/3.15.0-22-...` reference | Tag-naming convention rationale references the legacy branch shape as a frozen historical example to motivate bare-tag convention. Same disposition as `review/`. | Intentionally retained as legacy/illustrative |
| `post-release/SKILL.md` | line 380 — `Branch: claude/runtime-contract-v2-VWKUT` PRA example | PRA template example shows the legacy branch shape in a worked PRA example. Same disposition as `review/`. | Intentionally retained as legacy/illustrative |

**Disposition rationale (legacy retention):** the issue's non-goals section names "Retroactively renaming historical `claude/*` branches; legacy shapes warn-only forever" — the three illustrative examples in `review/`, `release/`, `post-release/` are exactly that: legacy examples documenting how earlier cycles named their branches. They are not contradicted by the new rule (which governs new cycles); they would actively mislead if rewritten as `cycle/{N}` (since the *historical* cycles those examples document predate the rule). Per AC 5 + non-goals, they stay. The grep confirms there are exactly three such retained references; no other lifecycle-skill drift was found.

### Intra-doc repetition

Per `alpha/SKILL.md` §2.3, intra-doc repetition is a peer enumeration class. Sites in this diff that carry the *same fact* across multiple sentences and require all-occurrences-update:

- **"γ creates the branch"** asserted across 7+ distinct sites in `CDD.md` (Role table, §Tracking ¶ "The cycle's branch...", §Tracking ¶ "Single named branch...", §1.4 γ Phase 1 step 3a, §1.4 dispatch prompt format parameters paragraph, §4.2 Ownership clause, §4.3 first sentence) and in `alpha/SKILL.md` (§2.1 step 2, §2.6 row 1) and in `beta/SKILL.md` (§1 paragraph 2) and in `gamma/SKILL.md` (Step 3a header + body, §Step map). All sites use consistent language: "γ creates `cycle/{N}` from `origin/main` before dispatch (§1.4 γ algorithm Phase 1 step 3a)" — verified by grep.
- **"α and β never create branches"** asserted across 6+ distinct sites (CDD.md Role table α + β rows, CDD.md §Tracking auth precondition paragraph, CDD.md §1.4 α step 5a, CDD.md §1.4 β step 3, CDD.md §4.2 Ownership, alpha/SKILL.md §2.1 step 2, beta/SKILL.md §1). All consistent.
- **`origin/cycle/{N}` polling target** appears as the canonical polling form across CDD.md §Tracking, §1.4 β step 3, gamma/SKILL.md §2.5 Step 3b, operator/SKILL.md §2.2 (with legacy `origin/claude/*` named warn-only). Verified consistent.
- **Legacy `origin/claude/*` warn-only naming** appears at exactly 3 sites in spec text (CDD.md §Tracking polling table row, CDD.md §Tracking ¶ "Single named branch...", CDD.md §4.2 Branch rule legacy block) plus 1 site in `operator/SKILL.md` polling snippet comment. All consistent: marked "warn-only / retrospective" with the same lineage note.

No drift detected. Closure claim is exhaustive across the named peer set.

## Harness audit

Per `alpha/SKILL.md` §2.4, harness audit applies when the diff "changes a parser, schema-bearing type, manifest shape, or runtime contract." This diff is markdown-only spec text — no parser changes, no schema changes, no manifest changes, no runtime contract changes. **Harness audit not applicable** — there is no schema for non-primary-language harnesses to drift from.

The shell snippets embedded in the spec text (γ creation snippet in `CDD.md` §1.4 step 3a + `gamma/SKILL.md` §2.5 Step 3a; β polling reference in `CDD.md` §1.4 β step 3; γ polling reference in `gamma/SKILL.md` §2.5 Step 3b; δ polling in `operator/SKILL.md` §2.2) are documentation snippets, not executable harnesses. They were inspected for shape consistency: all use `origin/cycle/{N}` (or `cycle/<N>` placeholder), all embed the N=10 reachability re-probe rule where the spec text says they should (β + γ), and the δ snippet correctly omits the reachability re-probe (δ polls coarsely, 5-min interval, the wake-up signal is gate-coarse not per-commit).

## Post-patch re-audit (`alpha/SKILL.md` §2.6 row 9)

Diff languages present: **Markdown only**. No Go, no shell-as-program, no YAML workflows, no test fixtures.

Per `alpha/SKILL.md` §2.6 row 9 polyglot re-audit guidance — for a Markdown-only diff, the matching toolchain is "table-shape + cross-reference check":

- **Table shape:** `CDD.md` §5.3 Artifact manifest table row 2 (Branch) updated — the row now has 10 columns same as siblings (`Step | Name | Phase | Role (§1.4) | Evidence | Format spec | Owner | Producer | Required | Skill`). Verified by inspection — pipe count matches sibling rows.
- **Cross-reference check:** every new section reference resolves:
  - "§1.4 γ algorithm Phase 1 step 3a" → `CDD.md` §1.4 γ algorithm Phase 1 step 3a (created in this diff)
  - "§1.4 α step 5a" → `CDD.md` §1.4 α algorithm step 5a (created in this diff)
  - "§1.4 β step 3" → `CDD.md` §1.4 β algorithm step 3 (rewritten in this diff)
  - "§1.4 γ dispatch prompt format" → `CDD.md` §1.4 γ dispatch prompt format (modified in this diff)
  - "§4.2 Branch rule" → `CDD.md` §4.2 (rewritten in this diff)
  - "§4.3 Branch pre-flight" → `CDD.md` §4.3 (rewritten in this diff)
  - "`gamma/SKILL.md` §2.5 Step 3a" → `gamma/SKILL.md` §2.5 Step 3a (created in this diff)
  - "`alpha/SKILL.md` §2.6 row 1" → `alpha/SKILL.md` §2.6 row 1 (updated in this diff)
  - "`beta/SKILL.md` §1" → `beta/SKILL.md` §1 (rewritten in this diff)
- **Re-read on HEAD:** `self-coherence.md` re-read top-to-bottom against HEAD; CDD Trace, AC mapping, peer enumeration all match HEAD diff.

## Role self-check

Per `alpha/SKILL.md` §2.5: "did α's work push ambiguity onto β? Is every claim backed by evidence in the diff?"

- **Did α push ambiguity onto β?** No, with one explicitly disclosed exception (see Known debt #1 below). The 12 ACs each have concrete evidence rows above, every cross-reference resolves to a target in this diff, and every claim of "no change required" (lifecycle skills) cites both the grep result and the issue's non-goals justification. The contract change is universal ("γ creates the branch") and the audit is exhaustive across the named peer set; no closure-overclaim risk.
- **Is every claim backed by evidence in the diff?** Yes. Every AC evidence row names the specific file + section + nature of change. Peer enumeration tables enumerate the role-skill peer set (4 peers, 4 updated) and the lifecycle-skill peer set (6 peers, 3 reviewed-no-change + 3 retained-as-legacy with disposition rationale). Intra-doc repetition is grepped: "γ creates the branch" / "α and β never create branches" / `origin/cycle/{N}` / legacy `origin/claude/*` — all four facts traced across all sites with consistent language.
- **Did α stay in voice?** Yes. The α close-out voice rule (factual observations and patterns; no dispositions) does not yet apply (close-out is post-merge); this self-coherence.md is α's authoring artifact and uses authoring voice. AC evidence rows are factual ("inserted X", "rewrote Y", "verified Z"); no speculation, no β-territory dispositions, no recommendations to γ.
- **Did α load only Tier 1?** Yes. Tier 3 skills are `none beyond CDD Tier 1` per dispatch prompt + issue ## Tier 3 skills. Tier 2 (engineering bundles) — none applied because the diff is markdown-only. The active skill set is `cdd/SKILL.md`, `cdd/CDD.md`, `cdd/alpha/SKILL.md` (loaded as the role surface), `cdd/issue/SKILL.md` (consulted for AC interpretation), `cdd/design/SKILL.md` (consulted; design "not required" justification documented in §Mode + CDD Trace step 6a).
- **Did α follow the new rules this diff introduces?** Yes (this is AC 12 self-application, internalized):
  - α never created a branch this cycle: `git switch cycle/287` from the harness branch was the first action.
  - α's pre-review gate (below) row 1 evaluates "verify `origin/cycle/287` is rebased onto current `origin/main`" — not "α rebases own branch."
  - α's review-readiness signal (below) commits to `origin/cycle/287`, the canonical coordination surface, with the explicit base SHA + head SHA + branch CI line per `alpha/SKILL.md` §2.7.
  - α did not push to `main`. α did not push to `claude/alpha-cdd-skill-1aZnw` (the harness branch). All commits land on `cycle/287`.

## Known debt

1. **`eb48e17` referenced in issue body is unavailable in this repo.** The issue's ## Work-shape section says "Spec text is already drafted in #283's chat record + the rolled-back local commit `eb48e17`'s message; α can ingest both as the authoritative spec source and refine." The commit `eb48e17` is local-only to γ's prior session (rolled back per 3.61.0 PRA §3 γ self-observation) and is not present in this repo's git object store (`git show eb48e17` returns "ambiguous argument: unknown revision"). α reconstructed spec text from the issue body's AC text (which already encodes the substantive content) plus the 3.61.0 PRA's §4b/§7 cycle-iteration + next-MCA notes plus #283's `gamma-clarification.md`. The spec changes match the AC text verbatim where the AC text is prescriptive (e.g. AC 1's `git switch -c cycle/{N} origin/main && git push -u origin cycle/{N}` snippet, AC 4's exact dispatch-prompt block); where the AC text is descriptive, α produced spec prose consistent with the existing CDD.md style. Risk: if `eb48e17`'s commit message contained a richer formulation than the issue body, α's reconstruction may have lost nuance. Mitigation: γ can review the diff against the rolled-back text; if drift is found, β returns RC and α refines.
2. **No CI-equivalent exists for markdown-only spec diffs.** This package has no CI step that verifies spec-internal cross-references resolve, no markdown table-shape linter, no grep-based peer-enumeration check. α ran the cross-reference + table-shape checks manually (post-patch re-audit section above). β should re-run them on review for an independent verification. Spec-internal coherence checking is absent from CI; that is an existing condition (#274 surfaced it, no MCA shipped) not introduced by this cycle. Out of scope for #287.
3. **Lifecycle table at `CDD.md` §4.1 step 2 row was not modified in this diff.** §4.1 row 2 says "Step 2 | Branch | Create a dedicated branch | Valid branch name." The role / ownership for §4.1 is the master role table at §1.4 (which I updated to name γ). §5.3 row 2 (which carries Role + Producer columns) was updated. §4.1 is intentionally minimal — it lists step name, purpose, required output — and does not encode role ownership. No drift.
4. **Operator's δ branch-cleanup rule unchanged.** The issue says "δ branch-cleanup rule continues to sweep merged branches; no functional change but warn-only legacy is named." δ's branch cleanup at γ §2.10 step 15 (`git branch -r --merged origin/main | grep -v main | grep -v HEAD | sed 's/origin\///' | xargs -I{} git push origin --delete {}`) is shape-agnostic — it sweeps any merged branch matching `origin/<name>` regardless of name. So `cycle/287` will be cleaned up identically to legacy `claude/*` branches once merged. No skill text touches the cleanup snippet itself; this is correct (the cleanup is universal, not branch-shape-specific).
5. **Mid-cycle rebase onto upstream `incremental write discipline` commit (`70ff2b1`).** While α was authoring this cycle, `origin/main` advanced with `70ff2b1 skill(cdd): incremental write discipline for α and β artifacts` (sigma — touches `alpha/SKILL.md` §2.5 + §2.7). α's pre-review gate row 1 ("verify `origin/cycle/{N}` is rebased onto current `origin/main`") fired correctly, α rebased and resolved the §2.7 conflict (HEAD's "append the review-readiness section as a separate commit" wording preserved + α's `origin/cycle/{N}` explicit naming retained). The new §2.5 incremental-write discipline is being followed prospectively from this commit forward: this §Pre-review gate section is a separate commit; the §Review-readiness signal will be a separate commit. The first 5 sections of this self-coherence.md (Gap, Skills, ACs, Self-check, Debt) were authored pre-rebase as a single coherent block in the commit `ae53183` (now rebased to `6a897a4`); a strict reading of the new §2.5 would call for those to have been 5 separate commits. α adopts the new discipline from this section forward. Risk: low — the new rule is about resilience to stream timeouts; α's pre-rebase write completed without interruption, so no partial-state recovery was needed.

## Pre-review gate

Per `alpha/SKILL.md` §2.6, ten rows. Rows 1 and 10 are transient (re-validated immediately before signaling review-readiness per `alpha/SKILL.md` §2.7).

| # | Row | State | Evidence (observed at 2026-04-30 00:18:24 UTC) |
|---|-----|-------|------------------------------------------------|
| 1 | α verifies `origin/cycle/287` is rebased onto current `origin/main` | ✅ | base SHA: `70ff2b1b80e49a30a4d7dddded49a5bd33669b32` (`origin/main`); cycle/287 rebased onto it; HEAD: `6a897a44b1bb6f85f14d1d60ee804a404ff7dc1f` (pre-this-commit). Rebase resolved one conflict in `alpha/SKILL.md` §2.7 — HEAD's "append the review-readiness section as a separate commit" wording preserved alongside α's explicit `origin/cycle/{N}` naming. **Transient — will be re-validated immediately before §Review-readiness commit (per `alpha/SKILL.md` §2.7 transient-row rule, derives from #266 F1/F2).** |
| 2 | `.cdd/unreleased/287/self-coherence.md` carries CDD Trace through step 7 | ✅ | CDD Trace table at top of this file; rows 0 through 7a present; step 7a "Pre-review" row present (this section is the evidence). |
| 3 | Tests are present, or explicit reason none apply | ✅ — explicit reason | CDD Trace step 6c row: "Not required — markdown / spec-only diff. There is no executable surface in this change. The cycle itself is the integration test (AC 12 self-application: the cycle uses `cycle/287`, exemplifying its own rule)." Re-stated in §Mode (mode = MCA, no Tier 3 skills, markdown-only diff). |
| 4 | Every AC has evidence | ✅ | AC 1 through AC 12 each have a §Acceptance Criteria — evidence subsection above with explicit file + section + nature-of-change citations. AC 12 (self-application) cites the head SHA `7533978` as γ's pre-dispatch state and the dispatch prompt block α received. |
| 5 | Known debt is explicit | ✅ | §Known debt above lists 5 explicit debt entries, each with a risk + mitigation. No "no debt" silent omission; no implicit debt. |
| 6 | Schema / shape audit completed when contracts changed | ✅ — N/A justified | The contract change is a *role-ownership* contract (γ creates branch, α/β check out), not a schema or runtime contract. Markdown-only diff; no schema-bearing types, parsers, or manifest shapes touched. §Harness audit above documents this explicitly: "Harness audit not applicable — there is no schema for non-primary-language harnesses to drift from." |
| 7 | Peer enumeration completed when closure claim touches a family of surfaces | ✅ | §Peer enumeration above: role-skill peer set `{ alpha/, beta/, gamma/, operator/ }` (4 peers, 4 updated); lifecycle-skill peer set `{ issue/, design/, plan/, review/, release/, post-release/ }` (6 peers, 3 reviewed-no-change + 3 retained-as-legacy with disposition rationale citing AC 5 + non-goals); intra-doc repetition (4 facts, all sites verified consistent by grep). |
| 8 | Harness audit completed when a schema-bearing contract changed | ✅ — N/A | Same as row 6: no schema-bearing contract changed. §Harness audit above documents the audit-not-applicable reasoning + inspects embedded shell snippets for shape consistency. |
| 9 | Post-patch re-audit completed after any mid-cycle patch — covering every language present in the diff | ✅ | §Post-patch re-audit above: diff languages = Markdown only; matching toolchain = "table-shape + cross-reference check"; both performed; 9 cross-references resolved; §5.3 row 2 column count verified consistent with sibling rows; full file re-read against HEAD. The mid-cycle rebase + conflict resolution count as a patch event; post-rebase re-audit performed (this very pre-review gate section is the post-rebase re-audit artifact). |
| 10 | Branch CI is green on the head commit | **deferred — explicit per `alpha/SKILL.md` §2.6 row 10 escape hatch** | This package has no CI workflow that runs on cycle branch pushes — the only CI surface (`build-verify.yml`) runs on `main` and on PRs, not on direct cycle-branch pushes (and the new artifact-channel protocol per #283 does not use PRs). Local CI is unavailable in this α session (no test surface to run; the diff is markdown-only and there is no markdown linter / cross-reference checker in this repo's CI). Per `alpha/SKILL.md` §2.6 row 10: "if local CI is unavailable, the artifact's review-readiness section says so explicitly and β waits for green before merge." α records this here and again in §Review-readiness below. β should not block on a green-CI signal that this branch / this diff cannot produce; β should instead verify diff coherence directly. **Transient — will be re-validated immediately before §Review-readiness commit.** |

**Gate verdict at 2026-04-30 00:18:24 UTC:** all 10 rows pass (rows 1 + 10 with transient/explicit-deferral notes; rows 6 + 8 with N/A justifications). Ready to write the §Review-readiness signal as a separate commit per `alpha/SKILL.md` §2.7.

## Review-readiness | round 1 | base SHA: 70ff2b1b80e49a30a4d7dddded49a5bd33669b32 | head SHA: 3336fe38801f9187ec04fd01e4907ecb096755b6 | branch CI: deferred (no cycle-branch CI surface in this repo — see §Pre-review gate row 10) | ready for β

**Re-validation of transient pre-review-gate rows at 2026-04-30 00:19:22 UTC** (per `alpha/SKILL.md` §2.7, derives from #266 F1/F2):

- **Row 1 (rebase) re-validated:** `git fetch --quiet origin main && git rev-parse origin/main` returns `70ff2b1b80e49a30a4d7dddded49a5bd33669b32` — unchanged since pre-review-gate observation (00:18:24 UTC). `git merge-base origin/main HEAD` returns `70ff2b1b80e49a30a4d7dddded49a5bd33669b32` (= origin/main exactly), confirming `cycle/287` is rebased onto current `origin/main` with no drift. `git rev-list --left-right --count origin/main...HEAD` returns `0	2` — 2 cycle commits ahead, 0 commits behind. Pass.
- **Row 10 (CI) re-validated:** Same as initial observation — this repo has no CI workflow that runs on cycle/* branch pushes (verified by inspection of `.github/workflows/`: `build-verify.yml` runs on `main` + PRs only; `release.yml` runs on tag push). The new artifact-channel coordination protocol (#283) deliberately does not use PRs, so PR-triggered CI does not apply to cycle branches in flight. β should verify diff coherence directly rather than wait for a green-CI signal that this configuration cannot produce; once β `git merge`s into main, `build-verify.yml` will run on the merge commit. Pass with deferral note carried forward.

**Cycle-branch state at signal time:**

- `origin/cycle/287` head: `3336fe38801f9187ec04fd01e4907ecb096755b6` (this self-coherence.md without the §Review-readiness section that this commit will add — i.e. the §Review-readiness section will be in the *next* SHA). After this commit lands, `origin/cycle/287` head will be the new SHA which β should fetch.
- Two cycle commits ahead of `origin/main`:
  - `6a897a4 skill(cdd): γ creates the cycle branch — α and β only check out cycle/{N}` (α R1 — the spec change + initial self-coherence draft)
  - `3336fe3 self-coherence(287): §Pre-review gate + §Debt update post-rebase` (α R1 — incremental commit per new §2.5)
  - (this commit, when it lands, will be the third)

**β intake checklist** (per `beta/SKILL.md` §1 + `CDD.md` §1.4 β step 3):

1. β receives the dispatch prompt with `Branch: cycle/287` line — must be the new format AC 4 introduces. If the operator paste-dispatches β with a pre-#287 prompt (no `Branch:` line), surface to operator.
2. β `git fetch origin cycle/287 && git switch cycle/287`. **β never creates a branch.** If the harness places β on a per-role pre-provisioned branch (e.g. `claude/{slug}-{rand}`), β refuses per `beta/SKILL.md` §1 Role Rule 1 and switches to `cycle/287`.
3. β reads `git diff origin/main..cycle/287` (5 files: `CDD.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `operator/SKILL.md`; +`.cdd/unreleased/287/self-coherence.md`).
4. β reads `.cdd/unreleased/287/self-coherence.md` end-to-end and verifies AC-by-AC mapping resolves against the diff.
5. β writes verdict to `.cdd/unreleased/287/beta-review.md` on `cycle/287`. **β's verdict commits land on `origin/cycle/287` directly** — not on a separate harness branch (this is the structural correction to the #283 R1 F1 pattern).

α now begins polling `.cdd/unreleased/287/beta-review.md` and the `cycle/287` head SHA per `CDD.md` §Tracking. 60-second interval. Will respond to β's RC findings on this branch with a fix-round appendix per `alpha/SKILL.md` §2.7.

— α (`alpha@cdd.cnos`) at 2026-04-30 00:19:22 UTC
