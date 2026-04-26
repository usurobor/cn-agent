# Coherence-Driven Development (CDD)

**Version:** 3.15.0
**Status:** Draft
**Date:** 2026-03-25
**Placement:** `src/packages/cnos.cdd/skills/cdd/`
**Audience:** Contributors, reviewers, maintainers, release operators
**Scope:** Canonical algorithm spec for how cnos selects, executes, measures, and closes substantial development cycles

---

## 0. Purpose

CDD is the development method used to evolve cnos coherently. Its purpose is not merely to ship features. Its purpose is to reduce incoherence across the system as a whole:

- doctrine
- architecture
- implementation
- runtime behavior
- operator understanding
- release state
- development process itself

A release is therefore not just a bundle of outputs. It is a measured coherence delta.

> A change is good not merely when it is implemented, but when it reduces incoherence across the system as a whole.

CDD is a triadic coherence method at development scale.

## Invocation model

CDD is invoked as:

- `cdd` at global scope
- one role skill at role-local scope
- zero or more lifecycle sub-skills at task-local scope

This is the execution model the package assumes.

---

## 1. Scope

CDD applies to substantial changes: work that spans design, code, tests, docs, process, packaging, runtime behavior, or release state. CDD also defines a small-change path for changes too small to warrant a full version-directory cycle.

### 1.1 Substantial change

A change is substantial when one or more of the following are true:

- it introduces or changes a subsystem, contract, or protocol
- it changes runtime behavior or ABI
- it changes package, security, or release surfaces
- it requires design, test, docs, and release alignment
- it will likely take more than one day
- it creates future coherence risk if handled informally

### 1.2 Small change

A change may take the small-change path when it is narrowly local, low-risk, and does not need a frozen snapshot. In the small-change path:

- bootstrap does not apply
- the coherence contract may live in the commit message or PR body
- self-coherence is optional
- the author still owes a named incoherence and an explicit scope

### 1.3 First principle

CDD begins from the same first principle as the coherent agent:

> There is a gap between model and reality.

In development terms:

- model = doctrine, architecture, design, operator understanding, intended behavior
- reality = code, runtime behavior, logs, failures, release outcomes, actual operator experience

CDD exists to close that gap through coherent action.

### 1.4 Large-file authoring rule

Any file longer than ~50 lines (close-outs, assessments, design docs, plans) must be written **section by section to disk**, reporting to the operator after each section is written. Do not compose the entire file in memory before writing. This prevents context loss from API timeouts, session interruptions, or compaction. Partial artifacts must be recoverable.

This applies to all roles and all CDD artifact types. It is not repeated in individual role or phase instructions.

### 1.5 Roles

CDD is triadic at the role level. A substantial cycle needs three distinct functions:

- **α** — implements
- **β** — reviews and releases
- **γ** — coordinates and unblocks

These are operational roles. They are not a claim that every cycle always uses three different humans.

| Role | Function | Steps owned | Responsibility | Identity constraint |
|------|----------|-------------|----------------|---------------------|
| **γ (Coordinator)** | Orchestrate | 0–1 + 11–13 + cycle-wide | Observe, select, issue creation, dispatch prompts to α and β, unblocking when stuck, cross-agent context, compliance verification, post-release assessment, cycle closure | Must hold full cycle context |
| **α (Implementer)** | Produce | 2–7a | Branch, bootstrap, gap, mode, artifacts (tests/code/docs), self-coherence, pre-review readiness, PR | Must be separate from β |
| **β (Reviewer + Releaser)** | Judge and integrate | 8–10 | Review (RC/A decision), merge, tag, deploy, β close-out | Must be separate from α |
| **δ (Operator)** | Route and disconnect | gates + 17 | Session routing, external gate execution (merge/tag/branch), post-cycle disconnect release | Owns platform actions agents cannot perform |

#### Triadic rule

- α produces the change.
- β judges the change and integrates it.
- γ orchestrates the cycle and ensures coherence across handoffs.

The structure is a **dyad plus coordinator**: α and β are two workers that interact through artifacts, isolated from each other. γ coordinates the dyad — sees both sides, does neither.

- α cannot see β's review reasoning or conversation state
- β cannot see α's implementation rationale or conversation state
- γ sees both — that is its function

β owns both review and release because the reviewer already has full artifact context when it's time to merge — splitting review from release creates a handoff that adds no value. γ owns coordination because issue quality determines implementation quality, dispatch prompts are the control surface, and unblocking requires cross-agent context that only the coordinator holds.

#### δ and the external boundary

δ (operator) is not a fourth triad role — δ is the boundary between the triad-as-whole and the platform. The triad is one-as-three (per `COHERENCE-FOR-AGENTS.md`): producer, judge, inspectable boundary. δ creates a new one-as-two relation with the triad-as-whole: the triad owns the cycle's internal coherence; δ owns the external boundary (platform gates, session routing, override authority, and the disconnect release that crystallizes the triad's output into a tagged whole).

When the operator serves as γ (two-agent configuration), δ and γ collapse — one person holds both coordination and external boundary authority. The δ skill handles this: "If δ is unavailable, γ may execute gates directly." The structural distinction remains even when the roles are held by the same participant.

#### Default flow

```text
γ (issue + dispatch) → α (implement + PR) → β (review) → RC → α (fix) → β (re-review)
                                                        → A  → β (merge, tag, deploy, assess)
                       γ (unblocks α or β when stuck)
```

#### Tracking: periodic poll, not event subscription

GitHub event notifications are unreliable when agents share a GitHub identity (self-authored PR comments, CI status changes, and review events may not trigger notifications). **All roles must track issue and PR activity via periodic polling, not GitHub event subscriptions.**

Polling has three parts: (a) the **query** that detects new state, (b) the **wake-up mechanism** that returns control to the role's session when state changes, and (c) **reachability verification** that the chosen query form actually works in the current environment. All three are mandatory.

Polling without a wake-up mechanism is silent — the loop runs but the role never reacts. Polling with an unreachable query form is silent in the same way — the loop runs, returns empty or errors, and the role assumes nothing has happened.

**Reachability preflight — run before committing to a query form.** The table below lists three query forms. They are not interchangeable — each has environment constraints:

- `gh` requires shell access and `gh auth` configured
- MCP tools require in-conversation invocation — they cannot run inside `Monitor` or background shell loops
- `git fetch` requires network access to the git remote from the execution environment
- Direct `api.github.com` access may be blocked by sandbox network policy

Before starting a polling loop, **probe the chosen query form once synchronously** and confirm it returns real data. If it fails, fall back to the next available form. If no form is reachable, surface the gap to δ/operator before proceeding — do not silently assume polling is working.

**Query — pick whichever is reachable in the environment.** Concrete equivalents:

| Surface | `gh` form (shell envs) | MCP form (MCP-only envs) | git form (clone-aware envs) |
|---|---|---|---|
| PR existence | `gh pr list --search "closes:#{N} OR refs:#{N}" --state open --json number` | `mcp__github__list_pull_requests` (filter by head ref or scan for "closes:#{N}" in body) | `git fetch origin && git branch -r` (compare to prior list, emit on new branch) |
| PR status / CI | `gh pr view {N} --json statusCheckRollup,reviews,state` | `mcp__github__pull_request_read` method=`get_status` / `get_check_runs` | — |
| Issue comments | `gh issue view {N} --comments` | `mcp__github__issue_read` method=`get_comments` | — |
| PR comments | `gh pr view {N} --comments` | `mcp__github__pull_request_read` method=`get_comments` / `get_review_comments` | — |
| Branch head SHA | — | — | `git fetch origin {branch} && git rev-parse origin/{branch}` (compare to prior, emit on change) |

**Wake-up mechanism — name it explicitly in the role's session.** Polling is only effective if the loop's transition produces a notification that wakes the role. Each environment has its own form:

- **`Monitor` (Claude Code on the web):** `Monitor`'s contract is "every stdout line is a notification delivered as a `task-notification` system message that wakes the session." Wrap polling in `until {transition-detected}; do {poll}; sleep 60; done` and emit only on **transition** (new branch, new SHA, new comment-id) — not on every poll iteration, or you flood the session's context.
- **Shell `gh` + harness wake hook:** if the harness wakes the session when a `gh ...` command completes its until-loop, `until {gh-poll}; do sleep 60; done` is sufficient. Verify the harness contract before relying on it.
- **Push (`mcp__github__subscribe_pr_activity`):** delivers `github-webhook-activity` messages once a PR exists. Reliability under shared GitHub identity is unverified — use as a complement to polling, not a replacement.

If neither a `Monitor`-equivalent nor a shell-wake harness exists, the role cannot autonomously detect cycle progression. Surface the gap to the operator before dispatch — do not rely on a wake-up contract that the environment does not provide.

Poll interval: 60 seconds unless the operator specifies otherwise. This applies to α (waiting for review), β (waiting for PR, waiting for α fixes), and γ (tracking full cycle).

**Stdout filter discipline:** transition-only emission is mandatory. A loop that emits on every poll will fill the session with `task-notification` blocks and consume context budget. β's #268 cycle ran four wake events across the full cycle by emitting only on new-branch / new-SHA transitions; the pattern is reproducible.

**Synchronous-baseline-pull is a precondition of transition-only polling.** Transition-only emission is correct on its own terms (avoid context flood) but has a structural blind spot: the loop's first iteration sets `prev` to the empty string (or current state) and silently absorbs whatever already exists. State that exists *before* the Monitor's first iteration will never surface as an event. Every transition-only Monitor must therefore be paired with a synchronous initial-state pull of the same surface (`mcp__github__list_pull_requests state=open` for PRs, `mcp__github__pull_request_read get_comments` for comments, etc.) immediately when the role's session starts — the synchronous channel owns the past, the polling channel owns the future. *Derives from: #274 cycle, where α + β + γ in three independent role sessions all hit first-iteration absorption — β's broad PR-list Monitor absorbed PR #274 as baseline; γ's `*230*` branch glob never matched the harness-encoded `claude/alpha-tier-3-skills-IZOsO` branch; α's git-only Monitor could not see comment activity until the operator prompted synchronously.*

**Branch-glob templates must not assume harness-encoded issue numbers.** The reference glob shape `'origin/claude/*-<N>-*'` assumes the harness encodes the issue number `<N>` in the branch path. Real harnesses may instead encode scope words and a random suffix (e.g. `claude/alpha-tier-3-skills-IZOsO` for issue #230) with no issue number anywhere. The defensive default is the broader glob `'origin/claude/*'` with downstream synchronous filtering against the issue/PR surface (the issue number lives reliably in the PR body's `Closes #N` line, not in the branch name). The narrow glob is acceptable only when the dispatching γ pins a canonical branch name in the dispatch prompt itself. *Derives from: #274 cycle, where β + γ both armed `*230*`-shaped globs that matched zero branches.*

#### γ algorithm

The compact algorithm is here; `gamma/SKILL.md` expands each phase into executable detail with gates, katas, and selection mechanics. When they diverge on role-local execution detail, the skill governs, provided it does not alter role ownership, lifecycle order, selection rules, dispatch contract, artifact contract, or closure obligations (those remain CDD.md's authority — see `cdd/SKILL.md` Conflict rule).

**Phase 1 — Dispatch**

1. Configure git identity using the project name: `git config user.name "gamma"` and `git config user.email "gamma@cdd.{project}"`
2. Observe and select the gap (§2)
3. Create the issue with full implementation guidance, including Tier 3 skills (§4.4)
4. Begin polling the issue immediately — γ must track the full cycle. Run the §Tracking reachability preflight first, then start the transition loop. When α opens a PR, poll that too.
5. Write α and β dispatch prompts (see format below). Both can be dispatched at the same time — β begins polling the issue and starts intake while α implements. β's skill handles waiting for the PR.
6. If α or β is blocked, diagnose and unblock: clarify requirements, resolve ambiguity, provide missing context

**Phase 2 — Release support**

7. If β deferred tag push or other release mechanics (env constraint per β step 8), request δ to execute the gate action. δ signals γ on completion (see `operator/SKILL.md` §3.5). If δ is unavailable, γ may execute directly.
8. If the issue did not auto-close on merge (missing `Closes #N`), close it: `gh issue close {number}`

**Phase 3 — Close-out triage**

9. Collect close-outs from both α and β. Both must exist on main before proceeding. If either is missing, request it.
10. Read both close-outs. Write the post-release assessment per `post-release/SKILL.md` — γ owns the PRA as the cycle-level observer. For each finding (from close-outs and the PRA), triage using CAP:
   - MCA available (skill patch, gate, mechanization) → ship it now as immediate output
   - No MCA yet, pattern real → MCI. Two kinds:
     - **Project MCI** (future cycles on this project need to know) → `.cdd/` in the repo
     - **Agent MCI** (future sessions of this agent need to know, any project) → agent hub (`cn-{agent}/threads/adhoc/`)
   - One-off, no pattern → drop

**Phase 4 — CDD iteration**

11. Check §9.1 triggers against this cycle's data:
    - Review rounds > 2?
    - Mechanical ratio > 20% (with ≥ 10 findings)?
    - Avoidable tooling/environmental failure?
    - Loaded skill failed to prevent a finding?
12. If any trigger fired: verify the assessment contains a Cycle Iteration section with root cause and MCA disposition.
13. Independently assess: did this cycle reveal a CDD process gap — a recurring friction, a missing gate, an underspecified step, or a skill that should have caught something but didn't? If yes, write and commit the skill/spec patch now. If no, state why not (one sentence). **This step applies even when no §9.1 trigger fired.** Triggers catch mechanical failures; this step catches process drift that triggers miss.

**Phase 5 — Hub memory and closure**

14. Update hub memory:
    - Daily reflection: cycle summary, scoring, MCI freeze status, next move
    - Adhoc thread: update or create the thread this cycle advances
15. Request δ to delete merged remote branches. If δ is unavailable, γ may execute directly: `git branch -r --merged origin/main | grep -v main | grep -v HEAD | sed 's/origin\///' | xargs -I{} git push origin --delete {}`
16. Cycle is closed. Commit the closure declaration to main and **signal δ**: *"Cycle #N closed. Next: #M. δ: disconnect."* This declaration must be γ's **last commit** for the cycle — no post-closure pushes. The signal to δ is what triggers the disconnect (step 17). If γ cannot reach δ directly, the closure commit on main is the fallback signal (δ polls for it per `operator/SKILL.md` §2.2).

**Phase 6 — Disconnect (δ)**

17. δ receives γ's closure declaration (step 16) as the trigger. δ lands any remaining δ session patches, then cuts the disconnect release: bump version, tag, push. The tag is the coherence boundary — it crystallizes the triad's output into an inspectable, immutable whole. Without it, the cycle's output bleeds into the next cycle with no named edge. See `operator/SKILL.md` §3.4. This is not optional. **Nothing is committed to main between the closure declaration and the disconnect tag except δ's own patches.**

#### γ dispatch prompt format

**To α:**
```
You are α. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md.
Issue: gh issue view <number>
Tier 3 skills: <list issue-specific skills>
```

**To β:**
```
You are β. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view <number>
```

**To γ (operator or γ-agent):**
```
You are γ. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md.
Issue: gh issue view <number>
```

Parameters: `{project}` is the project name (e.g. `cnos`, `myapp`). Git identity uses `{role}@cdd.{project}` (e.g. `alpha@cdd.cnos`). `{number}` is the **issue number** for dispatch prompts; β and α discover the PR from issue activity, α's review request, or the polling commands in §Tracking. A PR number may be used only for post-hoc review/recovery workflows, in which case the prompt must say explicitly that it is PR-scoped.

The prompt names the role, provides parameters, and points to the issue or PR. The CDD skill tells each role what to load (§4.4) and what to do (§1.4). γ does not enumerate skills or steps in the prompt — that is the skill's job. If the prompt needs to restate the algorithm, the algorithm is not clear enough — fix the skill.

#### α algorithm

The compact algorithm is here; `alpha/SKILL.md` expands each step into executable detail with self-coherence, pre-review gate, peer enumeration, and harness audit procedures. When they diverge on role-local execution detail, the skill governs, provided it does not alter role ownership, lifecycle order, selection rules, dispatch contract, artifact contract, or closure obligations (those remain CDD.md's authority — see `cdd/SKILL.md` Conflict rule).

1. Receive dispatch prompt from γ
2. Configure git identity using the project name from the dispatch prompt: `git config user.name "alpha"` and `git config user.email "alpha@cdd.{project}"`
3. Subscribe to the issue (`gh issue edit {number} --add-assignee @me` or equivalent) so you receive PR and update notifications
4. Load CDD skill, load all Tier 1 + Tier 2 skills (§4.4), load Tier 3 skills from the issue
5. Read the issue fully, read source files referenced in implementation guidance
6. Implement: branch, code, tests, self-coherence
7. Open PR (draft if CI unavailable locally), wait for CI green. PR body or commit message must include `Closes #N` or `Fixes #N` to auto-close the issue on merge.
8. Subscribe to PR notifications
9. Request review from β
10. If β returns RC: fix findings, push, re-request review
11. When β approves: write α close-out (cycle findings or "no findings"). **Commit the close-out to main directly** (not on the PR branch) — squash-merge destroys branch-only files.
12. Done

**α close-out:** Report cycle-level learnings to γ. Concrete findings (skill gaps, process friction, things to mechanize) or "no new findings" — explicitly stated, not omitted. This is α's input to γ's cycle iteration decision (§9.1). **Voice: factual observations and patterns only.** Do not recommend dispositions ("patch now", "file issue") — triage is γ's job. State what happened, what pattern it matches, and what surfaces were affected. Let γ decide the action.

#### β algorithm

The compact algorithm is here; `beta/SKILL.md` defines β's role boundary, load order, and closure discipline. The detailed review, release, and post-release procedures live in the lifecycle sub-skills (`review/`, `release/`, `post-release/`).

1. Receive dispatch prompt from γ (or pick up from α's review request)
2. Configure git identity using the project name from the dispatch prompt: `git config user.name "beta"` and `git config user.email "beta@cdd.{project}"`
3. Immediately begin polling the issue and PR (see §Tracking above for query forms, wake-up mechanism, and reachability preflight) — do not ask, just do it. **The poll requires a query, a wake-up mechanism, and a reachability probe.** Run the §Tracking reachability preflight first: probe your chosen query form synchronously, confirm it returns real data, and fall back if it doesn't. Then pick the wake-up form your harness provides (`Monitor` stdout-as-notification, shell-wake-on-loop-exit, or push subscription). Emit only on transition to avoid context flood. If the PR does not exist yet, **poll until α opens a PR that references this issue.** Reference shape (shell + `gh`):
   ```bash
   until gh pr list --search "closes:#<N> OR refs:#<N>" --state open --json number -q '.[0].number' 2>/dev/null | grep -q .; do
     sleep 60
   done
   ```
   Reference shape (MCP-only, `Monitor`-wrapped, transition-only stdout):
   ```bash
   # Baseline sync — run BEFORE the transition loop.
   # The transition loop absorbs first-iteration state silently (prev=∅ → cur=...).
   # Any branch or PR that already exists at loop-start is invisible to it.
   # Read current state synchronously at intake so the past is not lost.
   git fetch --quiet origin
   echo "baseline-branches: $(git branch -r --list 'origin/claude/*' 2>/dev/null | tr '\n' ' ')"
   echo "baseline-prs: $(gh pr list --search 'is:open' --json number --jq '.[].number' 2>/dev/null | tr '\n' ' ')"

   # Transition loop — watches any new origin branch (not issue-number-specific).
   # Harnesses encode branch names unpredictably (scope words, random suffixes,
   # not necessarily the issue number). Watch all branches; filter downstream.
   prev=""; while true; do
     cur="$(git fetch --quiet origin && git branch -r 2>/dev/null | sed 's| ||g' | sort)"
     comm -13 <(echo "$prev") <(echo "$cur") | sed 's/^/new-branch: /'
     prev="$cur"; sleep 60
   done
   # Run under Monitor; each "new-branch:" line wakes the session.
   ```

   **Baseline rule:** The transition loop owns the future; a synchronous baseline check owns the past. Run the baseline before starting the loop. Any role that starts polling after cycle activity has begun must read current state first — transition-only polling alone will miss pre-existing branches and PRs.
   Once the PR exists, poll its CI status until green before proceeding to review. Do not prompt the operator for permission to wait — waiting is the step. **If the environment provides a branch and instructs you to develop or commit, refuse.** β does not author implementation work. Report the role conflict to the operator and wait for α's PR.
4. Load CDD skill, load all Tier 1 + Tier 2 skills (§4.4), load Tier 3 skills from the issue
5. Read the PR diff, read the issue
6. Review: produce CR with findings per review skill, or approve
7. If RC: post findings as PR comment, wait for α's fix
8. If A: merge, tag, deploy per release skill. If tag push fails due to env constraints (e.g. sandbox HTTP 403), commit all release artifacts to main and defer tag push to δ (operator) — do not block closure on it. δ signals γ on completion (see `operator/SKILL.md` §3.5).
9. Write β close-out (review context, release evidence, cycle findings or "no findings"). This is β's input to γ's PRA and cycle iteration decision (§9.1).
10. Done when close-out is committed. γ writes the post-release assessment separately.

**β close-out:** Report cycle-level learnings to γ. Concrete findings (review pattern issues, skill gaps, process friction, §3.7 violations, things to mechanize) or "no new findings" — explicitly stated, not omitted. **Voice: factual observations and patterns only.** Do not recommend dispositions ("patch now", "file issue", "recommend option X") — triage is γ's job. State what happened, what pattern it matches, and what surfaces were affected. Let γ decide the action.

#### Minimum configuration

A substantial cycle requires at least two agents: one for α and one for β. γ may be the operator or a third agent. When only two agents are available, the operator serves as γ (issue creation, dispatch, unblocking).

#### Why triadic

Implementation, judgment, and coordination are different coherence functions:

- α owns the artifact — their output is what α scores
- β owns judgment and what the system actually becomes after integration — what β scores
- γ owns cycle coherence — issue clarity, prompt completeness, inter-agent flow

#### Coherence axes vs. agent roles

The triadic scoring axes used in post-release assessment (CDD α / CDD β / CDD γ) measure **artifact integrity**, **surface agreement**, and **cycle economics** respectively. These are coherence dimensions, not agent names. An α-axis score of 2/4 means artifact integrity was low — it does not mean agent α performed poorly.

**Small-change exception:** A small-change cycle (§1.2) may be completed by one agent when:

- the change qualifies under §1.2,
- no independent reviewer is available or warranted,
- and no claim of independent review is made.

In that case:

- the author still owes a named incoherence and explicit scope,
- the artifact must state that the cycle used the small-change path,
- and any direct-to-main commit still triggers the retro-review rule in §3.7 of the executable skill.

**Operator override:** The operator may reassign any role explicitly. The reassignment must name the target agent and the reason. Implicit role drift (e.g., reviewer starts authoring fixes mid-review) is not permitted — if β requests changes, α executes the fix.

---

## 2. Inputs

CDD selection begins from observation, not preference. Every substantial cycle reads these inputs before selecting work:

### 2.1 CHANGELOG TSC table

Read the current α / β / γ release state.

Question:
- which axis is weakest?

### 2.2 Encoding lag table

Read the lag state of open feature and process issues.

Questions:
- what is stale?
- what is growing?
- what is blocked by something else?

### 2.3 Doctor / status

Read the health of the running system.

Questions:
- is there a P0?
- is operational infrastructure broken?
- is the system able to observe, update, and maintain itself?

### 2.4 Last post-release assessment

Read the prior cycle's output.

Questions:
- what MCA was committed as next?
- what MCIs remain unresolved?
- what process debt was identified?

If no prior assessment exists, skip this input and select from §2.1–§2.3 alone.

---

## 3. Selection Function

CDD selection is coherence-driven. The next substantial gap is selected by the following function, in order.

### 3.1 P0 override

If doctor/status shows a P0 bug such as crash, data loss, or silent failure, that is the gap. No further selection logic applies until it is addressed.

**New-vs-known rule:** If the P0 was already visible when the last assessment was written and the assessment committed a different next MCA, the assessment decision governs unless the P0 has escalated (e.g. now causing active data loss or blocking all development). A known P0 that was weighed and deferred is not an override — it is a prioritized backlog item.

### 3.2 Operational infrastructure override

If core operational paths are broken, fix them before feature work. Examples:

- self-update broken
- logging absent
- health checks missing
- deployment path incoherent
- system cannot observe or maintain itself

These are not "nice to have." They are preconditions for coherent development.

**Sizing rule:** Before selecting infrastructure debt as a full cycle, ask whether the fix is cycle-sized or immediate-output-sized. If the fix is executable in minutes (a script, a hook, a one-line config change), execute it as an immediate output (§10.1) and continue to §3.3. Only select infrastructure debt as the cycle gap when the fix requires design, multiple files, tests, or review.

### 3.3 Assessment commitment default

If the last assessment named a next MCA and no stronger override fires, that MCA is selected by default. Observation may override it, but the override must be stated explicitly.

### 3.4 No clear winner — stale backlog re-evaluation

If §3.1–§3.3 produce no actionable candidate (e.g. P0s exist but have no clear fix path, assessment doesn't commit a next MCA, no operational debt), re-evaluate stale issues before selecting new work:

- For each stale issue: is it still a real gap, or has the system moved past it?
- **Descope** issues that are no longer coherence gaps (close with rationale).
- **Consolidate** issues that overlap or could be addressed by one MCA.
- **Commit** the stale issue with the clearest fix path as the next MCA.
- If no stale issue has a clear fix path either, enter observation mode (§3.9).

Stale backlog accumulating across multiple cycles without re-evaluation is itself an incoherence.

### 3.5 MCI freeze check

If the lag table contains stale issues, the next substantial MCA must come from the stale set. New design work is frozen until at least one stale item ships.

### 3.6 Weakest-axis rule

If no stronger rule decides selection, choose work that addresses the weakest current axis:

- α → structural / consistency work
- β → alignment / integration work
- γ → process / evolution work

### 3.7 Maximum leverage

Among candidates that address the weakest axis, choose the one that moves the most lag entries.

### 3.8 Dependency order

If A blocks B blocks C, choose A regardless of local excitement or presentation value.

### 3.9 Effort-adjusted tie-break

Between candidates with equal leverage and axis effect, choose the smaller one. Ship sooner, observe sooner, correct sooner.

### 3.10 No-gap case

If:

- no P0 exists
- no operational-debt override exists
- no stale lag item exists
- no prior assessment commitment forces a next MCA
- axes are healthy or tied

then do not start a new substantial cycle. Remain in observation or choose a small-change path.

---

## 4. Development Lifecycle

CDD owns the full arc from observation back to observation.

### 4.1 Lifecycle steps

| # | Step | Purpose | Required output |
|---|------|---------|-----------------|
| 0 | Observe | Read current coherence state | Selection inputs read |
| 1 | Select | Choose the next gap | Named selected gap |
| 2 | Branch | Create a dedicated branch | Valid branch name |
| 3 | Bootstrap | Create snapshot skeleton | Version dir + stubs |
| 4 | Gap | Name the incoherence precisely | Coherence contract draft |
| 5 | Mode | Choose MCA/MCI and governing skills | Mode + active skill set |
| 6 | Artifacts | Design → contract → plan → tests → code → docs (or delegated handoff §2.5a) | Aligned implementation artifacts |
| 7 | Self-coherence | Author checks own work against ACs and triad | Self-coherence report |
| 8 | Review | CLP with another CA or human reviewer | Converged review outcome |
| 9 | Gate | Verify release readiness | Gate checklist passes |
| 10 | Release | Tag, publish, announce | Release artifacts exist |
| 11 | Observe | Confirm runtime matches design | Observation result |
| 12 | Assess | Post-release assessment | Assessment artifact |
| 13 | Close | Execute immediate outputs and commit deferred outputs | Cycle closed |

Step 13 feeds back to step 0.

### 4.2 Branch rule

Every substantial change must be developed on its own dedicated branch. No substantial CDD work is performed directly on main.

Canonical branch format:

```text
{agent}/{version}-{issue}-{scope}
```

Version may be omitted when not yet known:

```text
{agent}/{issue}-{scope}
```

### 4.3 Branch pre-flight

Before creating the branch, verify:

- version segment, if present, is greater than the latest released/tagged version
- no remote branch already claims the same issue
- no open PR already covers the same issue
- branch name matches the canonical format
- current CI / main state is known
- the intended scope is declared before implementation begins

### 4.4 Skill loading

Skills are loaded in tiers. For substantial changes, all applicable tiers are mandatory. The tier structure is executable: every role can answer "did I load what this tier requires?" without ambiguity.

**Tier 1a — CDD authority (always loaded by every role):**

- `CDD.md` (this document)
- `src/packages/cnos.cdd/skills/cdd/SKILL.md` (the loader entrypoint)
- the current role skill: `alpha/SKILL.md`, `beta/SKILL.md`, or `gamma/SKILL.md`

These are loaded unconditionally. Roles do not load peer role skills.

**Tier 1b — CDD lifecycle phase skills (load per current phase or work shape):**

Lifecycle sub-skills under `cdd/`:

- `issue/` — when interpreting issue ACs or quality
- `design/` — when producing or judging design-required work
- `plan/` — when implementation sequencing is non-trivial
- `review/` — when reviewing
- `release/` — when merging/tagging/deploying
- `post-release/` — when assessing or closing

Load by phase. Roles that operate across multiple phases load the matching set.

**Tier 1c — β closure bundle (always loaded by β):**

β always loads:

- `review/SKILL.md`
- `release/SKILL.md`

γ always loads `post-release/SKILL.md` because γ owns the post-release assessment (PRA, see §5.3a Artifact Location Matrix).

**Tier 2 — General engineering (load the applicable bundle from `src/packages/cnos.eng/skills/eng/README.md`):**

Pick the bundle for the work shape:

- coding
- review
- design
- runtime / platform
- tooling
- writing

The engineering package README is the source of truth for which bundle covers which work shape and which skills it includes. CDD does not enumerate language- or platform-specific bundles here — that surface lives in the engineering package and changes independently of the method.

The skill that owns skill-program/frontmatter coherence is `src/packages/cnos.core/skills/skill/SKILL.md` (load when authoring or modifying skills). The skill that owns architecture/design reasoning is `src/packages/cnos.core/skills/design/SKILL.md` (load when reviewing or producing architecture-level decisions).

**Tier 3 — Issue-specific (named per issue):**

Skills that depend on what the work touches. The issue's "Skills and constraints" section names these. Examples:

- Language: `eng/{language}`
- Domain: `eng/ux-cli`, `eng/performance-reliability`, `eng/tool`
- Architecture: `eng/evolve`, `eng/write-functional`
- Skill authoring: `cnos.core/skills/skill`

γ names only Tier 3 in issues. Tier 1 and Tier 2 are mandatory and not repeated.

**Read each SKILL.md file before beginning any work step.** Naming a skill without reading it is not loading it. Loaded skills are **hard generation constraints** — not post-hoc review checklists.

The Tier 3 skill set must be stated alongside mode. Example:

```text
Mode: MCA
Tier 3 skills: eng/{language}, eng/ux-cli
```

When in doubt about mode, apply CAP: if the answer is already in the system, cite it (MCA) — don't reinvent it (MCI). If two paths close the same gap, take the lighter one unless the heavier one buys durability the lighter one cannot.

Review (step 8) checks whether the implementation is consistent with all loaded tiers. Findings that a loaded skill would have prevented are process debt (§6.1).

---

## 5. Artifact Contract

CDD is artifact-driven. Every substantial cycle must produce inspectable artifacts.

### 5.0 Terminology: post-release, assessment, close-out, closure

- **Post-release** is the umbrella phase after release. It covers lifecycle steps 11–13: observe, assess, close.
- **Assessment** (a.k.a. post-release assessment, PRA) is the γ-owned repo artifact at the canonical path declared in §5.3a Artifact Location Matrix.
- **Close-out** is a role-local findings record written by α, β, and γ at the canonical paths in §5.3a. Close-outs feed γ's PRA and triage. They are not a substitute for the PRA.
- **Closure** is the final cycle state: PRA committed, all close-outs on main, immediate outputs executed, deferred outputs committed, hub memory updated.

### 5.1 Bootstrap

The first diff on the branch must create a version directory for every bundle that will receive a frozen snapshot.

Path convention:

```text
docs/{tier}/{bundle}/{X.Y.Z}/
```

Each version directory must contain:

- README.md — snapshot manifest
- one stub per declared deliverable

Artifacts outside version directories, such as PR body files or navigation updates, are not required as bootstrap stubs.

### 5.2 Ordered artifact flow

The canonical artifact order is:

1. design
2. coherence contract
3. plan
4. tests
5. code
6. docs
7. self-coherence
8. review
9. gate
10. release
11. observe
12. assess
13. close

### 5.3 Artifact manifest

CDD is artifact-driven. For substantial changes, each lifecycle step must leave one inspectable evidence surface. This table is the master reference for all step attributes.

| Step | Name | Phase | Role (§1.4) | Evidence | Format spec | Owner | Producer | Required | Skill |
|------|------|-------|-------------|----------|-------------|-------|----------|----------|-------|
| 0 | Observe | observe | γ | CDD Trace row: inputs read, selected signal | §5.4 | primary branch artifact | agent | always | cdd |
| 1 | Select | observe | γ | CDD Trace row: selected gap + issue | §5.4 (trace) + `.github/ISSUE_TEMPLATE/cdd-issue.md` (issue) | primary branch artifact + issue tracker | agent | always | cdd |
| 2 | Branch | build | α | valid branch name | — | branch + CDD Trace row | mechanical | always | cdd |
| 3 | Bootstrap | build | α | version directory + manifest README + declared stubs | §5.1 | branch diff | mechanical | substantial only | cdd |
| 4 | Gap | build | α | named incoherence / coherence contract | PR template §Gap or design/SKILL.md §3.1 | primary branch artifact | agent | always | cdd |
| 5 | Mode | build | α | mode + active skills (+ bundle/level if used) | PR template §Mode or design/SKILL.md §3.1 | primary branch artifact | agent | always | cdd + applicable eng bundle from `src/packages/cnos.eng/skills/eng/README.md` |
| 6a | Design | build | α | design artifact or explicit "not required" | design/SKILL.md §3.1 | primary branch artifact | agent | substantial only | design |
| 6b | Plan | build | α | plan artifact or explicit "not required" | `docs/gamma/cdd/PLAN-TEMPLATE.md` | primary branch artifact or linked plan | agent | L7 / cycle-sized | plan |
| 6c | Tests | build | α (or delegated implementer) | test files or explicit reason none apply | — (diff) | diff / primary branch artifact | agent | always | eng/test |
| 6d | Code | build | α (or delegated implementer) | implementation diff or "docs/process only" | — (diff) | diff / primary branch artifact | agent | always | active loaded Tier 2/Tier 3 generation skills (language skill, runtime/platform skill, tooling skill — see `src/packages/cnos.eng/skills/eng/README.md`) |
| 6e | Docs | build | α (or delegated implementer) | changed canonical docs / specs / READMEs | — (diff) | diff | agent | when docs affected | eng/document for durable docs; eng/write-functional for functional/dataflow prose; `cnos.core/skills/skill` when authoring or modifying skills |
| 6f | Delegated handoff | build | α → implementer | implementation prompt with: active skills, test requirements per module, engineering conventions, artifact order + self-verification report from implementer | alpha/SKILL.md §2.2 | prompt + report | delegator + implementer | when implementation is delegated | cdd |
| 7 | Self-coherence | build | α (or delegated implementer) | SELF-COHERENCE.md | `docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md` | version directory | agent | substantial only | cdd |
| 7a | Pre-review | build | α | branch rebased onto current `main`; PR body carries CDD Trace through step 7; tests reference ACs; known debt explicit; **schema/shape audit across all touched files** when contracts change — when introducing or changing a canonical form, verify (a) the new form is present across all relevant files AND (b) any superseded form has been removed; cite the verifying command (e.g. grep that returns exactly one match per file); **workspace-global library-name uniqueness check** when adding a new `(library (name X))` stanza; **post-patch re-audit** — after any mid-cycle code change, re-read the PR body top-to-bottom and verify every CDD Trace / invariant / self-coherence row still matches HEAD; **CI green on head commit** before requesting review (draft-until-green when local verification unavailable) | alpha/SKILL.md §2.6 | PR body | mechanical | always | cdd |
| 8 | Review | review | β | review artifact / PR review / comment link | review/SKILL.md output format | review surface | reviewer | always | review |
| 9 | Gate + merge | release | β | gate result / release-readiness evidence + PR merge | `docs/gamma/cdd/GATE-TEMPLATE.md` | release or review surface | mechanical + reviewer | always | release |
| 10 | Release | release | β | CHANGELOG row, tag, release note | CHANGELOG.md ledger + release/SKILL.md | release surface | agent + mechanical | always | release + eng/document when authoring release notes/changelog prose |
| 11 | Observe | close | γ | post-release observation result | post-release/SKILL.md | post-release assessment | γ | always | post-release |
| 12 | Assess | close | γ | POST-RELEASE-ASSESSMENT.md | post-release/SKILL.md output template | version directory | γ | always | post-release |
| 12a | Skill patch | close | γ | skill/spec patches for recurring failure modes identified in close-out triage or PRA §3; synced across all affected surfaces under src/packages/ | post-release/SKILL.md §3 + gamma/SKILL.md §2.7–§2.10 | same commit as PRA when identified during assessment, otherwise γ close-out follow-on commit | γ | when close-outs or PRA identify a recurring failure or skill gap | post-release, cdd |
| 13 | Close | close | γ | immediate outputs executed (incl. 12a patches) + deferred committed | post-release/SKILL.md §6 CDD Closeout | post-release assessment | γ | always | post-release |

**Primary branch artifact:** the PR body (`.github/PULL_REQUEST_TEMPLATE.md`) for L5/L6 changes, or the design artifact (design/SKILL.md §3.1) for larger changes.

**Role key (§1.4):** *γ (coordinator)* = steps 0–1 (observe, select, issue creation, dispatch) + steps 11–13 (PRA, cycle iteration, closure) + cycle-wide coordination, *α (implementer)* = steps 2–7a (branch, bootstrap, gap, mode, artifacts, self-coherence, pre-review), *β (reviewer + releaser)* = steps 8–10 (review RC/A decision, merge, deploy, β close-out). Delegated implementer is α-side. Merge is part of step 9 (gate + merge).

**Producer key:** *agent* = judgment required, *mechanical* = automatable by cnos (#94), *reviewer* = produced by the review process.

Rules:
- "Not required" is valid only when stated explicitly.
- An omitted step with no explicit note is incomplete, not implicit.
- Small-change mode may collapse steps 4–7 into commit/PR-body evidence, but the same distinctions still apply.

### 5.3a Artifact Location Matrix

For release-scoped triadic cycles, every named artifact has exactly one canonical location. Verifiers (e.g. `cn cdd-verify`) enforce these paths as canonical and treat any other location as legacy/warn-only.

| Artifact | Canonical repo location | CDD package default | Noncanonical / legacy / scratch |
|---|---|---|---|
| Version snapshot directory | `docs/{tier}/{bundle}/{X.Y.Z}/` | `docs/gamma/cdd/{X.Y.Z}/` | `.cdd/releases/{X.Y.Z}/` is **not** the frozen snapshot — it is triadic protocol/scratch space |
| POST-RELEASE-ASSESSMENT.md (PRA) | `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `.cdd/releases/{X.Y.Z}/beta/POST-RELEASE-ASSESSMENT.md` and `.cdd/releases/{X.Y.Z}/beta/ASSESSMENT.md` are legacy/warn-only |
| α close-out | `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` | same | PR comment is acceptable evidence only for PR-scoped, unreleased, non-triadic cycles |
| β close-out | `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md` | same | PR comment is acceptable evidence only for PR-scoped, unreleased, non-triadic cycles |
| γ close-out | `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md` | same | none — γ close-out is always a repo artifact |
| γ kata verdict (optional) | `.cdd/releases/{X.Y.Z}/gamma/KATA-VERDICT.md` when kata is available | same | warn-only when kata unavailable |
| CHANGELOG ledger row | `CHANGELOG.md` (Release Coherence Ledger) — must include Version, C_Σ, α, β, γ, **Level**, coherence note | same | none |
| RELEASE.md | `RELEASE.md` at repo root, included in the release commit | same | CI auto-generated body is not an acceptable substitute |
| Hub memory | external agent-hub state (not a repo artifact) | external agent-hub state | the PRA must record path / commit-sha / unavailable-reason; verifiers do not inspect the external hub directly |

Rules:

- `.cdd/` is triadic protocol space and role-local close-out / close-out evidence storage. It is **not** the canonical frozen post-release snapshot.
- `docs/{tier}/{bundle}/{X.Y.Z}/` is the canonical frozen post-release snapshot directory. The PRA lives there; close-outs do not.
- Tags are bare `X.Y.Z` everywhere (VERSION file, git tag, branch name version segment, CHANGELOG row, RELEASE.md, snapshot directory). `v`-prefixed tags are legacy and warn-only.
- For PR-scoped, unreleased, non-triadic cycles, close-outs may be PR comments; release-scoped triadic cycles require the canonical repo paths above.

### 5.4 CDD Trace

Every substantial cycle must carry a lightweight execution trace. Use lifecycle step numbers, not section numbers. For steps 0–10, the trace lives in the primary branch artifact. For steps 11–13, closure lives in the post-release assessment.

The primary branch artifact is the artifact that owns the named incoherence, mode, active skills, and acceptance criteria. In most cycles this is the design artifact. For governance/process work it may be the governing doc being changed. For smaller substantial changes it may be the PR body when no separate design doc is required.

Required format:

```markdown
## CDD Trace
| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Observation inputs read; selected signal |
| 1 Select | — | — | Selected gap |
| 2 Branch | branch | cdd | Branch created / verified against §4.2 / §4.3 |
| 3 Bootstrap | version dir | cdd | Bootstrap stubs created or explicit small-change exemption |
| 4 Gap | primary artifact | — | Named incoherence / coherence contract |
| 5 Mode | primary artifact | skill1, skill2 | Work shape, level (if used), mode, active skills |
| 6 Artifacts | design / plan / tests / docs | — | Artifact progress or explicit "not required" |
| 7 Self-coherence | SELF-COHERENCE.md / PR body | cdd | AC-by-AC self-check completed |
| 7a Pre-review | PR body | cdd | Pre-review gate passed (alpha/SKILL.md §2.6) |
| 8 Review | review surface | review (+ others if loaded) | CLP review result |
| 9 Gate | review or release surface | release (+ eng/document if release prose authored) | Release-readiness decision |
| 10 Release | release surface | release (+ eng/document if release prose authored) | Tag / changelog / release decision |
```

Rules:
- One row per completed lifecycle step.
- Step column carries both the number and the name for readability.
- "Skills loaded" is required when skills shaped generation or lifecycle execution.
- If a lifecycle skill is used later (review, release, writing, post-release), record it when it becomes active.
- Missing rows mean the step is not yet evidenced.
- Contradictory rows are findings.

### 5.5 Supporting rules

- one source of truth per fact
- derive, do not duplicate
- update docs before release
- write tests before or alongside the code they validate
- build-sync source asset changes before commit
- enumerate affected files before implementation begins
- every AC must map to evidence before review
- **all review findings must be resolved before merge** — the author fixes every finding (A/B/C/D) on the branch before merge. No "approved with follow-up." The only exception is a finding requiring a design decision outside the issue's scope, which the reviewer must explicitly name as "deferred by design scope" and the author must file as an issue before merge. See review skill §7.0.

### 5.6 Frozen snapshot rule

After release, version directories are frozen by repository policy. Only path-reference repairs are allowed after freeze:

- markdown links
- backtick paths

No semantic content may change.

---

## 6. Mechanical vs Judgment Boundary

CDD is rigorous, but not fully mechanical.

### 6.1 Mechanical

These may be enforced by tools or checklists:

- branch naming
- branch uniqueness
- version-directory presence
- required artifact presence
- stale cross-reference detection
- AC accounting
- frozen snapshot integrity
- gate checks
- release artifact presence
- review-quality metrics
- process-debt filing when thresholds trigger

### 6.2 Judgment

These remain judgment-bearing:

- what the real incoherence is
- whether MCA or MCI is the right intervention
- whether α / β / γ scoring is substantively sound
- whether a review has truly converged
- whether a design is coherent enough to implement
- whether iteration should stop

Tools may validate the existence of judgment artifacts. They do not replace the judgment itself.

---

## 7. Review

CDD review uses CLP. Every substantial review should answer:

- TERMS — what are we talking about?
- POINTER — where is the incoherence?
- EXIT — what changed, or what still blocks closure?

Every reviewer should be asked for:

- α / β / γ scores
- weakest-axis diagnosis
- concrete patch suggestions
- iterate or converge verdict

The review skill owns the detailed protocol. CDD owns when review is required and what it must decide.

---

## 8. Gate

Release may proceed only when:

- the selected gap was actually addressed
- required artifacts exist
- self-coherence exists for substantial change
- review converged
- CI/build/test requirements pass
- docs, code, and release artifacts agree
- the previous release has an assessment
- known debt is explicit, not implicit

A passing gate means:

- structurally ready to release

It does not mean:

- intellectually perfect

---

## 9. Assessment

Post-release assessment is mandatory for substantial releases. It must record:

- measured coherence delta
- encoding lag table
- MCA/MCI balance
- process learning (including active skill re-evaluation against review findings)
- review quality metrics
- CDD self-coherence (α artifact integrity, β surface agreement, γ cycle economics)
- cycle iteration (see §9.1)
- next move commitment

The CHANGELOG TSC entry at release time is provisional. Assessment governs the final judgment of the cycle.

### 9.1 Cycle iteration

When a cycle exceeded expected effort — extra review rounds, avoidable mechanical errors, tooling failures, or environmental surprises — the assessment must include a cycle iteration section.

**Trigger:** Any of:
- review rounds exceeded target (default: 2)
- mechanical finding ratio exceeded 20%
- avoidable tooling or environmental failure delayed the cycle
- a skill that was loaded failed to prevent a finding it covers

**Required artifact:** A `## Cycle Iteration` section in the post-release assessment (located in `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`). Must contain these fields:

```markdown
## Cycle Iteration

### Triggers fired
- [ ] review rounds > 2 (actual: N)
- [ ] mechanical ratio > 20% (actual: N%)
- [ ] avoidable tooling/environmental failure
- [ ] loaded skill failed to prevent a finding

### Friction log
What went wrong in the cycle itself (process, not code).

### Root cause
One of: design flaw | skill gap | tooling gap | environmental | one-off

### Skill impact
Which active skill should have prevented the friction?
If a loaded skill failed: name it and patch it as immediate output (§10.1).

### MCA
System change shipped or proposed. "Won't repeat" without a mechanism is not an MCA.

### Cycle level
L5 | L6 | L7 (per ENGINEERING-LEVELS.md §6 — level = lowest miss)
Justification: [one line explaining the level]
```

The cycle level is also recorded in the CHANGELOG TSC table as a suffix: e.g. `(cycle: L6)`.

**Cycle level assessment:** After the friction log and MCA, assess the executed engineering level of the cycle. This applies the level framework from `docs/gamma/ENGINEERING-LEVELS.md` to cycles (change sets), not just individual diffs. Use the §6 diagnostic questions adapted for cycle scope:

- **L5: local correctness** — Was the code locally correct before review? Did it compile, pass tests, follow current patterns? If mechanical errors (compilation failures, type errors, broken assertions, missing tests) reached review, L5 was not met. Cycle caps at L5.
- **L6: system-safe execution** — Did the change stay coherent across docs, runtime, artifacts, tests, and operator truth? Were failure modes bounded and visible? If cross-surface drift (package sync, authority-sync, doc/code mismatch, test coverage gaps) reached review, L6 was not met. Cycle caps at L6.
- **L7: system-shaping leverage** — Did the cycle change the system boundary so the friction class gets easier or disappears? If friction was fixed locally but no system change prevents recurrence (no MCA shipped, no skill patched, no gate added), L7 was not met. Cycle caps at L6. Not every cycle needs L7 — only assess when friction occurred and a system-shaping response was available.
- **L7 achieved** — The cycle shipped an MCA that eliminates or reduces a friction class for future cycles.

The cycle level is the lowest miss. A cycle with L5 misses caps at L5 regardless of whether it also shipped an L7 MCA. Each level must be earned cleanly.

Record the cycle level in the assessment. Over time, cycle levels track whether the development process itself is improving.

**Gate:** If the trigger fires and the cycle iteration section is absent, the cycle cannot close (§10.3).

If the cycle went cleanly (no trigger fires), cycle iteration may be omitted.

---

## 10. Closure

A cycle is not closed merely because code merged.

### 10.1 Immediate outputs

These must be executed within the same cycle:

- changelog corrections
- missing documentation
- skill/process micro-patches
- skill patches identified by cycle iteration (§9.1) — if a loaded skill failed to prevent a finding it covers, patch the skill now, not later
- issue filing required by the assessment
- lag-table updates
- metadata fixes
- hub memory writes: daily reflection (cycle state, scoring, MCI status) and adhoc thread update (which ongoing thread this release advances)

Skill/spec patches produced as immediate outputs must pass CLP β: does this change create a mismatch with any canonical or derived surface? If the edited artifact has a paired authority surface (executable skill ↔ canonical spec), both must be updated in the same commit.

### 10.2 Deferred outputs

These may become the next cycle's work, but must be committed concretely:

- next MCA issue number
- owner, if known
- target branch name, if known
- first AC
- freeze/resume state for MCI backlog

### 10.3 Closure rule

A cycle closes only when:

- all immediate outputs are executed (including hub memory writes)
- all deferred outputs are captured as explicit next-cycle commitments
- cycle iteration (§9.1) is present if any trigger fired

That is the handoff from step 13 back to step 0.

---

## 11. Related documents

### 11.1 Executable summary

`src/packages/cnos.cdd/skills/cdd/SKILL.md` is the package-visible loader entrypoint for this spec. It is not a second fact source.

### 11.2 Companion rationale

`docs/gamma/cdd/RATIONALE.md` explains why CDD takes this shape. Frozen release snapshot directories under `docs/gamma/cdd/{X.Y.Z}/` may also contain version-local rationale files.

---

## 12. Retro-packaging rule

If a substantial change lands direct-to-main as an exception to the branch/bootstrap discipline, it must be followed immediately by:

- a retro-snapshot in the appropriate version directory (frozen copies of the changed canonical artifacts)
- a self-coherence artifact covering the change
- a version-history entry in the bundle README

This closes the loophole where substantial method rewrites bypass their own packaging discipline. The method must eat its own cooking.

---

## 13. Non-goals

CDD does not:

- optimize primarily for speed
- treat issue queues as self-justifying
- reduce review to local diff reading
- treat release as "tag and hope"
- confuse a shipped feature with a closed coherence cycle
