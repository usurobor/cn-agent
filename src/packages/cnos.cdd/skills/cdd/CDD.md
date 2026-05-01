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
- the coherence contract may live in the commit message or `.cdd/unreleased/{N}/self-coherence.md`
- self-coherence is optional
- the author still owes a named incoherence and an explicit scope

**Artifact collapse for small-change cycles.** The following table states what applies, what collapses, and what is never waived:

| Artifact | Small-change status | Reason |
|---|---|---|
| `self-coherence.md` | Optional — commit message may substitute | Scope is narrow enough to state in prose |
| `alpha-closeout.md` | Optional — write if there are findings; one-liner "no findings" otherwise | No substantial cycle findings expected |
| `beta-closeout.md` | Not applicable | No independent β in single-author small change |
| `gamma-closeout.md` | Not applicable unless γ coordinated | No triadic coordination |
| `beta-review.md` | Not applicable | No independent β review |
| `RELEASE.md` | **Required before tag/release — no exception** | Release notes are required regardless of change size; CI uses them |
| `.cdd/releases/{X.Y.Z}/{N}/` move | **Required at release time** | Version association is required regardless of change size |
| Post-release assessment (PRA) | **Required after every release** | Assessment is mandatory for every release, patch or major |
| Version directory / bootstrap stubs | Not required | Frozen snapshot is not needed for small changes |

A small-change cycle may not use the small-change path to avoid `RELEASE.md`, cycle-directory movement, or the post-release assessment. Those three are unconditional.

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
| **γ (Coordinator)** | Orchestrate | 0–2 + 11–13 + cycle-wide | Observe, select, issue creation, **branch creation (`cycle/{N}` from `origin/main`, see §1.4 γ algorithm Phase 1 step 3a)**, dispatch prompts to α and β, unblocking when stuck, cross-agent context, compliance verification, post-release assessment, cycle closure | Must hold full cycle context |
| **α (Implementer)** | Produce | 3–7a (uses the branch γ created at step 2) | Check out `cycle/{N}` (never creates), bootstrap, gap, mode, artifacts (tests/code/docs), self-coherence, pre-review readiness, `.cdd/unreleased/{N}/self-coherence.md` review-readiness signal | Must be separate from β |
| **β (Reviewer + Integrator)** | Judge and integrate | 8–10 | Polls `origin/cycle/{N}` directly (never creates a branch and refuses harness pre-provisioned per-role branches; see `beta/SKILL.md` §1), Review (RC/A decision), merge into main, β close-out | Must be separate from α |
| **δ (Operator)** | Route and disconnect | gates + 17 | Session routing, external gate execution (tag/release/deploy), post-cycle disconnect release | Owns platform actions agents cannot perform |

#### Triadic rule

- α produces the change.
- β judges the change and integrates it.
- γ orchestrates the cycle and ensures coherence across handoffs.

The structure is a **dyad plus coordinator**: α and β are two workers that interact through artifacts, isolated from each other. γ coordinates the dyad — sees both sides, does neither.

- α cannot see β's review reasoning or conversation state
- β cannot see α's implementation rationale or conversation state
- γ sees both — that is its function

β owns review and merge because the reviewer already has full artifact context when it's time to integrate — splitting review from merge creates a handoff that adds no value. δ owns the release boundary (tag, deploy, disconnect) because those are external platform actions that may require permissions the triad does not hold. γ owns coordination because issue quality determines implementation quality, dispatch prompts are the control surface, and unblocking requires cross-agent context that only the coordinator holds.

#### δ and the external boundary

δ (operator) is not a fourth triad role — δ is the boundary between the triad-as-whole and the platform. The triad is one-as-three (per `COHERENCE-FOR-AGENTS.md`): producer, judge, inspectable boundary. δ creates a new one-as-two relation with the triad-as-whole: the triad owns the cycle's internal coherence; δ owns the external release boundary (tag, deploy, disconnect release) and session routing.

**Merge is β's authority, not δ's.** β merges `cycle/{N}` into `main` on approval — this is the natural conclusion of review, not a platform gate. δ's authority begins at the release boundary: tagging, deploying, and the disconnect release that crystallizes the triad's output into an inspectable whole. δ may request changes at the release boundary if something was missed, but the merge itself requires no δ authorization.

When the operator serves as γ (two-agent configuration), δ and γ collapse — one person holds both coordination and external boundary authority. The δ skill handles this: "If δ is unavailable, γ may execute gates directly." The structural distinction remains even when the roles are held by the same participant.

#### Default flow

```text
γ (issue + dispatch) → α (implement + branch + .cdd/unreleased/{N}/self-coherence.md) → β (review) → RC → α (fix) → β (re-review)
                                                        → A  → β (merge) → α/β close-outs → γ (PRA + triage) → δ (release-boundary preflight) → γ (closure) → δ (tag/release/deploy)
                       γ (unblocks α or β when stuck)
```

#### Tracking: artifact-driven coordination via `.cdd/unreleased/`

Triadic CDD coordination uses three repo-native surfaces:

1. **GitHub Issues** — gap-naming surface. γ files; α/β/γ subscribe.
2. **Git branches** — α's implementation isolation. β reads via `git diff main..{branch}`.
3. **`.cdd/unreleased/{N}/`** — per-cycle directory. The single coordination surface during in-version work. `{N}` is the issue number. One directory per cycle; multiple files inside named for what they are.

GitHub Pull Requests are **not** part of the triadic protocol. PR creation, polling, subscription, and comment threading add ceremony without value for agent-to-agent coordination — every action a PR records can be expressed as a write to a file in `.cdd/unreleased/{N}/` plus a branch commit. Issues remain (gap-naming); branches remain (isolation); PRs are removed.

**Path layout:** `.cdd/unreleased/{N}/` (e.g. `.cdd/unreleased/283/`). Per-cycle directory keyed by issue number. **Role distinguished by filename, not by directory or by a rigid `{role}.md` shape.** A role may write multiple files (e.g. α writes `self-coherence.md` AND `alpha-closeout.md`); each file's name describes what it is. After release, `release/SKILL.md` §2.5a moves each cycle directory into `.cdd/releases/{X.Y.Z}/`.

**Canonical filenames — what each role writes.**

| File | Owner | Purpose | Minimum contents |
|---|---|---|---|
| `self-coherence.md` | α | Primary branch artifact: gap, mode, AC mapping, CDD Trace, review-readiness signal, fix-round appendices | gap, mode, active skills, impact graph, AC-by-AC evidence, CDD Trace through step 7a, known debt, review-readiness section per round (base SHA, head SHA, branch CI state, "ready for β" line), fix-round appendices when β returns RC |
| `beta-review.md` | β | Review record across rounds | round-by-round verdict (RC / A), findings with severity + type, merge instruction (the `git merge` action β will execute on approval), provisional-vs-green CI state |
| `alpha-closeout.md` | α | α close-out narrative (after merge) | cycle summary, findings α-side, friction log, observations and patterns, cycle-level engineering level reading; **factual observations only — no dispositions** (γ's job) |
| `beta-closeout.md` | β | β close-out narrative + merge evidence | review context, narrowing pattern across rounds, merge evidence (merge commit SHA, branch state), β-side findings; **factual observations only — no dispositions** |
| `gamma-closeout.md` | γ | γ close-out triage + cycle iteration | close-out triage table (every finding from α + β + PRA → disposition), §9.1 trigger assessment, cycle iteration section if any trigger fired, deferred outputs, next-MCA commitment |

Roles may add additional files when the cycle warrants it (e.g. `alpha-design.md` if a separate design artifact lives in the cycle dir, `gamma-dispatch.md` if dispatch evidence is large enough to deserve its own file). The filename pattern is `{role}-{purpose}.md` (e.g. `alpha-closeout.md`) for role-owned content, and bare `{purpose}.md` (e.g. `self-coherence.md`) for cycle-shared content where the convention identifies the author.

**Coordination loop.**

1. γ files issue #N, dispatches α and β
2. α implements on a branch; writes `.cdd/unreleased/{N}/self-coherence.md` carrying the review-readiness signal
3. β polls `.cdd/unreleased/{N}/`; on `self-coherence.md` review-readiness, reads `git diff main..{branch}` + the cycle dir + the issue
4. β writes verdict to `.cdd/unreleased/{N}/beta-review.md`
5. If RC: α appends fix-round to `.cdd/unreleased/{N}/self-coherence.md`, β reads and re-reviews (appends new round to `beta-review.md`)
6. On A: β `git merge`s the branch into main with a commit message containing `Closes #N` (no `gh pr merge`) and pushes. If β cannot push, δ may execute the push on β/γ request — this is execution of β's integration authority, not δ approval.
7. After merge: α writes `.cdd/unreleased/{N}/alpha-closeout.md`; β writes `.cdd/unreleased/{N}/beta-closeout.md` (merge evidence, review evidence, cycle findings)
8. γ reads `.cdd/unreleased/{N}/{alpha-closeout,beta-closeout}.md`, writes PRA, triages findings
9. δ performs release-boundary preflight (see Phase 5a below). If δ requests changes, γ routes to β, α, or operator override as appropriate.
10. γ closes the cycle. δ cuts the disconnect release.
9. `release/SKILL.md` §2.5a moves `.cdd/unreleased/{N}/` → `.cdd/releases/{X.Y.Z}/{N}/` for every cycle that shipped in the release

**Issue auto-close.** Use `Closes #N` or `Fixes #N` in the merge commit message (`git merge -m "Closes #N: ..."`) to auto-close on merge. If the convention is missed, γ closes via `gh issue close {N}`.

**Polling — in-session only.** Polling applies during each role's active session. A role that has completed its phase and exited is not polling. The instructions below describe what each role monitors while it is running; they do not imply a background daemon loop. All roles poll three surfaces:

- **issue activity** — comments, label changes
- **branch state** — `origin/cycle/{N}` head-SHA changes (β waits for α's `self-coherence.md` review-readiness signal on the branch γ pre-created; γ tracks the full cycle)
- **`.cdd/unreleased/{N}/` directory state** — new files, file mtime / blob-SHA changes (α waits for `beta-review.md` to land/update; β waits for `self-coherence.md` to update with a new fix-round; γ watches the whole dir)

Polling has three parts: (a) the **query** that detects new state, (b) the **wake-up mechanism** that returns control to the role's session when state changes, and (c) **reachability verification** that the chosen query form works in the current environment. All three are mandatory.

Polling without a wake-up mechanism is silent — the loop runs but the role never reacts. Polling with an unreachable query form is silent in the same way — the loop runs, returns empty or errors, and the role assumes nothing has happened.

**Reachability preflight — run before committing to a query form.** The table below lists query forms. They are not interchangeable — each has environment constraints:

- `gh` requires shell access and `gh auth` configured (used only for issue activity, not PRs)
- MCP tools require in-conversation invocation — they cannot run inside `Monitor` or background shell loops
- `git fetch` requires network access to the git remote from the execution environment
- Direct `api.github.com` access may be blocked by sandbox network policy

Before starting a polling loop, **probe the chosen query form once synchronously** and confirm it returns real data. If it fails, fall back to the next available form. If no form is reachable, surface the gap to δ/operator before proceeding — do not silently assume polling is working.

**Query — pick whichever is reachable in the environment.**

| Surface | `gh` form (shell envs) | MCP form (MCP-only envs) | git form (clone-aware envs) |
|---|---|---|---|
| Issue comments | `gh issue view {N} --comments` | `mcp__github__issue_read` method=`get_comments` | — |
| Cycle branch existence (γ pre-flight) | — | — | `git fetch --quiet origin cycle/{N} && git rev-parse --verify origin/cycle/{N}` (γ runs before dispatch; α and β receive the branch name in the dispatch prompt and assert existence at intake — they do not glob and do not invent) |
| Cycle branch head SHA | — | — | `git fetch --quiet origin cycle/{N} && git rev-parse origin/cycle/{N}` (compare to prior, emit on change — this catches α's review-readiness commit, α's fix-round commits, β's verdict commits, γ's clarifications, all of which land on the cycle's branch) |
| `.cdd/unreleased/{N}/` directory state on the cycle's branch | — | — | `git fetch --quiet origin cycle/{N} && git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` (compare blob SHAs across the directory; emit on any file added or any blob SHA changed) |
| Legacy `'origin/claude/*'` glob | — | — | `git branch -r --list 'origin/claude/*'` is **warn-only / retrospective**: retained for tracking historical cycles whose branches predate the `cycle/{N}` rule. New cycles do not poll this glob. |

**The cycle's branch is the canonical coordination surface during in-version work.** Canonical branch name: `cycle/{N}` where `{N}` is the issue number (one cycle = one branch = one named target for all polling — see §4.2). **γ creates the branch from `origin/main` before dispatching α and β** (see §1.4 γ algorithm Phase 1 step 3a). All role artifacts (`self-coherence.md`, `beta-review.md`, `*-closeout.md`, `gamma-*.md`) live on the cycle branch. α and β commit their artifacts to **the same branch γ created** (push directly to `origin/cycle/{N}` with their own git identity), not to separate branches and not to `origin/main`. `main` is the merge target only, never the in-cycle coordination surface; cycle-dir files arrive on `main` as part of the eventual `git merge` (β step 8). Polling `origin/main` for in-flight cycle artifacts is silent — the spec must always poll `origin/cycle/{N}` while the cycle is open.

Auth precondition for the dyad-on-one-branch model: α, β, and γ all have push access to `origin/cycle/{N}`. If the operator's harness pre-provisions a per-role branch (e.g. `claude/{slug}-{rand}`), the role refuses (and surfaces to operator) any instruction to commit cycle-dir files to that pre-provisioned branch — those files belong on `origin/cycle/{N}`. **α and β never create branches.** If `origin/cycle/{N}` does not exist when α or β starts, the role surfaces the dispatch error to γ and refuses to invent a name. *Derives from #283 β round-1 F1: γ resolved to candidate (a) "branch-polling canonical, one cycle branch holds all role artifacts." Refined in #287: γ owns branch creation under a canonical name; α and β own only check-out and commit-to-the-branch.*

**Wake-up mechanism — name it explicitly in the role's session.** Polling is only effective if the loop's transition produces a notification that wakes the role. Each environment has its own form:

- **`Monitor` (Claude Code on the web):** `Monitor`'s contract is "every stdout line is a notification delivered as a `task-notification` system message that wakes the session." Wrap polling in `until {transition-detected}; do {poll}; sleep 60; done` and emit only on **transition** (new branch, new SHA, new file in cycle dir, changed blob SHA on existing file) — not on every poll iteration, or you flood the session's context.
- **Shell `gh` / `git` + harness wake hook:** if the harness wakes the session when a polling command completes its until-loop, `until {poll}; do sleep 60; done` is sufficient. Verify the harness contract before relying on it.

If neither a `Monitor`-equivalent nor a shell-wake harness exists, the role cannot autonomously detect cycle progression. Surface the gap to the operator before dispatch — do not rely on a wake-up contract that the environment does not provide.

Poll interval: 60 seconds unless the operator specifies otherwise. This applies to α (waiting for β verdict in `.cdd/unreleased/{N}/beta-review.md`), β (waiting for α's review-readiness signal in `.cdd/unreleased/{N}/self-coherence.md` on `origin/cycle/{N}`, then waiting for fix-round appendices), and γ (tracking issue + `origin/cycle/{N}` + every file in `.cdd/unreleased/{N}/` across the full cycle).

**Stdout filter discipline:** transition-only emission is mandatory. A loop that emits on every poll will fill the session with `task-notification` blocks and consume context budget.

**Synchronous-baseline-pull is a precondition of transition-only polling.** Transition-only emission is correct on its own terms (avoid context flood) but has a structural blind spot: the loop's first iteration sets `prev` to the empty string (or current state) and silently absorbs whatever already exists. State that exists *before* the Monitor's first iteration will never surface as an event. Every transition-only Monitor must therefore be paired with a synchronous initial-state pull of the same surface immediately when the role's session starts — the synchronous channel owns the past, the polling channel owns the future. For the `cycle/{N}` model, the baseline pull is per-cycle: `git rev-parse --verify origin/cycle/{N}` to confirm the branch exists, `git rev-parse origin/cycle/{N}` for the head SHA, and `git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` for the cycle artifact set, plus `gh issue view {N} --comments` for issue activity. Reading `.cdd/unreleased/` from `origin/main` only surfaces cycles that have already merged; in-flight cycles live on `origin/cycle/{N}` and will be invisible to a `main`-only baseline. *Derives from: #274 cycle, where α + β + γ in three independent role sessions all hit first-iteration absorption — β's broad PR-list Monitor absorbed PR #274 as baseline; γ's `*230*` branch glob never matched the harness-encoded `claude/alpha-tier-3-skills-IZOsO` branch; α's git-only Monitor could not see comment activity until the operator prompted synchronously. #283 β round-1 F1 surfaced the symmetric polling-source bug for the new artifact-exchange model. #287 closed the source: γ creates `cycle/{N}` before dispatch and the dispatch prompt names it explicitly, so polling targets a single named branch instead of a glob.*

**Single named branch — no globs for new cycles.** Polling targets `origin/cycle/{N}` directly: `git rev-parse origin/cycle/{N}` for the head SHA, `git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` for the cycle artifact set. There is no glob discovery step for new cycles — γ tells α and β the branch name in the dispatch prompt (see §1.4 γ dispatch prompt format), and α / β assert `origin/cycle/{N}` exists at intake. The legacy glob `'origin/claude/*'` is **retained only for retrospective tracking** of historical cycles whose branches predate this rule (see §4.2 Branch rule); it must not be used as a discovery surface for new cycles. *Derives from: #274 + #283 R1 F1 branch-discovery friction. The narrow glob (`'origin/claude/*-<N>-*'`) and the broader glob (`'origin/claude/*'`) both fail in different ways; the named-branch shape (`origin/cycle/{N}`) eliminates the friction class.*

**`git fetch` reliability is an explicit dependency.** The polling loops above use `git fetch --quiet origin cycle/{N}` to refresh the remote ref, then read the local copy with `git rev-parse` / `git ls-tree`. The `--quiet` flag suppresses progress output, but it also silences network and auth flake (DNS hiccup, proxy 502, expired token, transient 4xx from the server) — fetch returns 0, the local ref does not advance, and the per-iteration comparison sees `cur == prev` and emits nothing. The transition loop then drops every commit landing during the flake window. **Mitigation rule:** after **N successive polling iterations with no transition** (suggest **N = 10**, ≈ 10 minutes at 60s interval), do a synchronous reachability re-probe with explicit stderr capture:

```bash
git fetch --verbose origin cycle/{N} 2>&1 | tee /tmp/cycle-{N}-fetch.log
# inspect for: "Could not resolve host", "fatal: unable to access", "remote: HTTP",
# "fatal: Authentication failed", "Couldn't connect to server", non-zero exit.
```

If the re-probe succeeds (no stderr errors, ref advances or is confirmed unchanged), continue the transition loop. If the re-probe fails, **surface the failure to the operator immediately** — the role cannot autonomously detect cycle progression with a broken transport. Do not silently keep looping. *Derives from #283 β observation #3 / 3.61.0 PRA §4b: β's Monitor polling silently dropped α-branch transitions during round-2 dispatch; suspected `git fetch --quiet` masking auth/network flake.*

#### γ algorithm

The compact algorithm is here; `gamma/SKILL.md` expands each phase into executable detail with gates, katas, and selection mechanics. When they diverge on role-local execution detail, the skill governs, provided it does not alter role ownership, lifecycle order, selection rules, dispatch contract, artifact contract, or closure obligations (those remain CDD.md's authority — see `cdd/SKILL.md` Conflict rule).

**Phase 1 — Dispatch**

1. Configure git identity using the project name: `git config user.name "gamma"` and `git config user.email "gamma@cdd.{project}"`
2. Observe and select the gap (§2)
3. Create the issue with full implementation guidance, including Tier 3 skills (§4.4)
3a. **Create the cycle branch.** From `origin/main`, create `cycle/{N}` and push it before any dispatch:
   ```bash
   git fetch --quiet origin main
   git switch -c cycle/{N} origin/main
   git push -u origin cycle/{N}
   ```
   Run γ's branch pre-flight (§4.3) before pushing. The branch must exist on `origin` *before* α and β are dispatched — α and β do not create branches and the dispatch prompts name `cycle/{N}` explicitly (§1.4 γ dispatch prompt format). One cycle = one branch = one named target for all polling.
4. Begin polling the issue and `.cdd/unreleased/{N}/` on `origin/cycle/{N}` immediately — γ must track the full cycle. Run the §Tracking reachability preflight first, then start the transition loop against the single named branch (`git rev-parse origin/cycle/{N}` for head SHA, `git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` for the artifact set). No glob discovery — γ created the branch in step 3a and knows its name.
5. Write α and β dispatch prompts (see format below). Both prompts include a `Branch: cycle/{N}` line. Both can be dispatched at the same time — β begins polling the issue and `origin/cycle/{N}` and starts intake while α implements.
6. If α or β is blocked, diagnose and unblock: clarify requirements, resolve ambiguity, provide missing context

**Phase 2 — Post-merge support**

7. If β could not push the merge (env constraint, auth), request δ to execute the push on β's behalf. This is execution of β's integration authority, not δ approval. δ confirms completion to the requesting role (see `operator/SKILL.md` §3.3). If δ is unavailable, γ may execute directly.
8. If the issue did not auto-close on merge (missing `Closes #N` in the merge commit message), close it: `gh issue close {number}`

**Phase 3 — Close-out triage**

9. Collect close-outs from both α and β. Both must exist on main (`.cdd/unreleased/{N}/alpha-closeout.md` and `.cdd/unreleased/{N}/beta-closeout.md`) before proceeding. If either is missing, request it.
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

**Phase 5a — δ release-boundary preflight**

After β merges `cycle/{N}` into main and α/β close-outs exist, γ requests δ release-boundary preflight before declaring the cycle closed.

δ reviews the merged release boundary — not the implementation judgment. δ verifies:
- merge commit is present on main;
- required release artifacts are present;
- tag/release/deploy preconditions are satisfied;
- no external platform/auth/deployment blocker is visible;
- any requested release action is executable.

Outcomes:
- **Proceed** — δ confirms the release boundary is ready. γ may close the cycle, after which δ cuts the disconnect release.
- **Request changes** — δ names the release-boundary gap and returns control to γ. No tag is cut. γ routes the change to β, α, or operator override as appropriate.
- **Override** — if δ blocks a release for reasons outside the triad's internal judgment, δ declares the override and names the affected boundary.

15. γ requests δ release-boundary preflight. If δ confirms, proceed to closure.

**Phase 5b — Hub memory and closure**

16. Update hub memory:
    - Daily reflection: cycle summary, scoring, MCI freeze status, next move
    - Adhoc thread: update or create the thread this cycle advances
17. Request δ to delete merged remote branches. If δ is unavailable, γ may execute directly: `git branch -r --merged origin/main | grep -v main | grep -v HEAD | sed 's/origin\///' | xargs -I{} git push origin --delete {}`
18. Cycle is closed. Commit the closure declaration to main: *"Cycle #N closed. Next: #M."* This must be γ's **last commit** for the cycle — no post-closure pushes.

**Phase 6 — Disconnect (δ)**

19. δ lands any remaining δ session patches, then cuts the disconnect release: bump version, tag, push. The tag is the coherence boundary — it crystallizes the triad's output into an inspectable, immutable whole. The tag is git-observable: γ and all future agents can see it. Without it, the cycle's output bleeds into the next cycle with no named edge. See `operator/SKILL.md` §3.4. This is not optional. **Nothing is committed to main between the closure declaration and the disconnect tag except δ's own patches.**

#### γ dispatch prompt format

**To α:**
```
You are α. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md.
Issue: gh issue view <number>
Branch: cycle/<number>
Tier 3 skills: <list issue-specific skills>
```

**To β:**
```
You are β. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view <number>
Branch: cycle/<number>
```

**To γ (operator or γ-agent):**
```
You are γ. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md.
Issue: gh issue view <number>
```

Parameters: `{project}` is the project name (e.g. `cnos`, `myapp`). Git identity uses `{role}@cdd.{project}` (e.g. `alpha@cdd.cnos`). `{number}` is the **issue number** for dispatch prompts. The α and β prompts include a `Branch:` line because γ creates `cycle/{N}` from `origin/main` before dispatch (§1.4 γ algorithm Phase 1 step 3a) and α/β `git switch` to it — they do not invent a name and they do not glob to discover it. The γ prompt does not include `Branch:` because γ owns its own session and creates the branch as part of Phase 1.

The prompt names the role, provides parameters, points to the issue, and (for α and β) names the branch. The CDD skill tells each role what to load (§4.4) and what to do (§1.4). γ does not enumerate skills or steps in the prompt — that is the skill's job. If the prompt needs to restate the algorithm, the algorithm is not clear enough — fix the skill.

#### α algorithm

The compact algorithm is here; `alpha/SKILL.md` expands each step into executable detail with self-coherence, pre-review gate, peer enumeration, and harness audit procedures. When they diverge on role-local execution detail, the skill governs, provided it does not alter role ownership, lifecycle order, selection rules, dispatch contract, artifact contract, or closure obligations (those remain CDD.md's authority — see `cdd/SKILL.md` Conflict rule).

1. Receive dispatch prompt from γ. The prompt names the cycle branch explicitly (`Branch: cycle/{N}`).
2. Configure git identity using the project name from the dispatch prompt: `git config user.name "alpha"` and `git config user.email "alpha@cdd.{project}"`
3. Subscribe to the issue (`gh issue edit {number} --add-assignee @me` or equivalent) so you receive comment notifications
4. Load CDD skill, load all Tier 1 + Tier 2 skills (§4.4), load Tier 3 skills from the issue
5. Read the issue fully, read source files referenced in implementation guidance
5a. **Check out the cycle branch** named in the dispatch prompt: `git fetch origin cycle/{N} && git switch cycle/{N}`. **α never creates a branch.** If `origin/cycle/{N}` does not exist, surface the dispatch error to γ and refuse to invent a name — γ owns branch creation (γ algorithm Phase 1 step 3a). If the harness placed α on a different branch (e.g. `claude/{slug}-{rand}`), switch off it before any code change; that branch is not the cycle branch and must not receive cycle commits.
6. Implement: code, tests, docs, self-coherence (the branch already exists — α only commits to it)
7. Write `.cdd/unreleased/{N}/self-coherence.md` carrying the review-readiness signal (gap, mode, active skills, AC mapping, CDD Trace, known debt). Wait for branch CI green before signaling readiness in the artifact.
8. Begin polling `.cdd/unreleased/{N}/beta-review.md` for β's verdict (per §Tracking)
9. If β returns RC: fix findings on the branch, append a fix-round section to `.cdd/unreleased/{N}/self-coherence.md`, re-poll
10. **α close-out — re-dispatch mechanism.** In the sequential bounded dispatch model, α has exited before β approves and merges. γ requests δ to re-dispatch α after β merge using the α close-out prompt format (see §1.6a Re-dispatch prompt formats). α is re-dispatched, reads `.cdd/unreleased/{N}/beta-review.md` (the approved verdict) and the merged branch state, writes `.cdd/unreleased/{N}/alpha-closeout.md` (cycle findings or "no findings"), commits to main, and exits. The release skill's §2.5a move carries the cycle directory into the release directory at release time — do not duplicate close-outs elsewhere. **Provisional close-out fallback:** if re-dispatch is unavailable, α may write a provisional close-out at review-readiness time (before exit), explicitly marked `[provisional — pending β outcome]`. γ supplements with PRA observations. The provisional form must be declared as known debt in self-coherence.md §Debt.
11. Done

**α close-out:** Report cycle-level learnings to γ. Concrete findings (skill gaps, process friction, things to mechanize) or "no new findings" — explicitly stated, not omitted. This is α's input to γ's cycle iteration decision (§9.1). **Voice: factual observations and patterns only.** Do not recommend dispositions ("patch now", "file issue") — triage is γ's job. State what happened, what pattern it matches, and what surfaces were affected. Let γ decide the action.

#### β algorithm

The compact algorithm is here; `beta/SKILL.md` defines β's role boundary, load order, and closure discipline. The detailed review, release, and post-release procedures live in the lifecycle sub-skills (`review/`, `release/`, `post-release/`).

1. Receive dispatch prompt from γ. The prompt names the cycle branch explicitly (`Branch: cycle/{N}`).
2. Configure git identity using the project name from the dispatch prompt: `git config user.name "beta"` and `git config user.email "beta@cdd.{project}"`
3. Immediately begin polling the issue and `origin/cycle/{N}` (see §Tracking above for query forms, wake-up mechanism, and reachability preflight) — do not ask, just do it. **The poll target is a single named branch — no glob.** Run the §Tracking reachability preflight first: probe `git rev-parse --verify origin/cycle/{N}` synchronously, confirm it returns a SHA, and fall back to the operator if it does not (γ creates the branch in γ algorithm Phase 1 step 3a; if it is missing at β intake, that is a γ pre-flight failure to surface). Then pick the wake-up form your harness provides (`Monitor` stdout-as-notification or shell-wake-on-loop-exit). Emit only on transition to avoid context flood. **β never creates a branch.** Reference shape (MCP-only, `Monitor`-wrapped, transition-only stdout):
   ```bash
   # Baseline sync — run BEFORE the transition loop.
   # The transition loop absorbs first-iteration state silently (prev=∅ → cur=...).
   # Read current state synchronously at intake so the past is not lost.
   git fetch --quiet origin cycle/<N>
   git rev-parse --verify origin/cycle/<N>     # confirm γ created the branch; fail loud if missing
   echo "baseline-cycle-dir: $(git ls-tree -r --name-only origin/cycle/<N> .cdd/unreleased/<N>/ 2>/dev/null | tr '\n' ' ')"

   # Transition loop — single named branch, no glob.
   prev_head=""; empty_iters=0
   while true; do
     if ! git fetch --quiet origin cycle/<N>; then
       # transport failure — surface and continue; the reliability re-probe below will catch persistent failures
       echo "fetch-error: cycle/<N>"
     fi
     cur_head="$(git rev-parse origin/cycle/<N> 2>/dev/null)"
     if [ "$cur_head" != "$prev_head" ] && [ -n "$cur_head" ]; then
       echo "branch-update: cycle/<N> → $cur_head"
       prev_head="$cur_head"
       empty_iters=0
     else
       empty_iters=$((empty_iters + 1))
       # §Tracking git fetch reliability rule: every 10 empty iterations, do a synchronous reachability re-probe.
       if [ "$empty_iters" -ge 10 ]; then
         git fetch --verbose origin cycle/<N> 2>&1 | tee /tmp/cycle-<N>-fetch.log >&2 || \
           echo "reachability-fail: cycle/<N> — surface to operator"
         empty_iters=0
       fi
     fi
     sleep 60
   done
   # Run under Monitor; each transition line wakes the session.
   # The single named branch catches α's review-readiness commit, α's fix-round
   # commits, β's verdict commits, and γ's clarifications — all on
   # origin/cycle/<N>, all observable as branch-head SHA transitions.
   ```

   **Baseline rule:** The transition loop owns the future; a synchronous baseline check owns the past. Run the baseline before starting the loop. Any role that starts polling after cycle activity has begun must read current state first — transition-only polling alone will miss pre-existing branches and artifacts.
   Once `.cdd/unreleased/{N}/self-coherence.md` is present on `origin/cycle/{N}` and signals review-readiness, poll branch CI status until green before proceeding to review. Do not prompt the operator for permission to wait — waiting is the step. **β refuses harness pre-provisioned per-role branches and refuses any instruction to develop, commit, or implement on them.** β does not author implementation work, and β's verdict commits belong on `origin/cycle/{N}`, not on a separate β-side branch (per §Tracking branch-polling rule + `beta/SKILL.md` §1 Role Rule 1). Report the role / branch conflict to the operator and wait for α's review-readiness signal.
4. Load CDD skill, load all Tier 1 + Tier 2 skills (§4.4), load Tier 3 skills from the issue
5. Read `git diff main..{branch}`, read every file in `.cdd/unreleased/{N}/` (start with `self-coherence.md`), read the issue
6. Review: produce CR with findings per review skill, or approve
7. If RC: append findings round to `.cdd/unreleased/{N}/beta-review.md` (commit on main or push the branch with the file), wait for α's fix-round in `.cdd/unreleased/{N}/self-coherence.md`
8. If A: append approval verdict to `.cdd/unreleased/{N}/beta-review.md`, then `git merge` the branch into main with a commit message containing `Closes #N` (no `gh pr merge`) and push. If β cannot execute the push (env constraint, auth), δ may execute the push on β/γ request — this is execution of β's integration authority, not δ approval. Tag, deploy, and release are δ's responsibility at the release boundary (see γ algorithm Phase 5a).
9. Write `.cdd/unreleased/{N}/beta-closeout.md` (review context, merge evidence, cycle findings or "no findings"). This is β's input to γ's PRA and cycle iteration decision (§9.1). The release skill's §2.5a move carries the cycle directory into the release directory at release time.
10. Done when `beta-closeout.md` is committed to main. γ writes the post-release assessment separately.

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

### 1.6 Coordination model

**Current model: sequential bounded role invocations.**

Each role is dispatched one at a time via `claude -p` or equivalent, completes its phase, and exits. The triad runs sequentially (γ → α → β → γ for re-dispatch), not concurrently.

| Property | Value |
|---|---|
| Dispatch mechanism | `claude -p` one role at a time |
| Role session lifetime | Bounded — exits when phase is complete |
| In-session polling | Applies during each role's active session window (§Tracking) — not a background daemon loop |
| Fix rounds | δ re-dispatches α per γ's instruction when β returns RC |
| α close-out | δ re-dispatches α after β merge, per γ's request (§1.4 α step 10; see §1.6a for prompt format) |
| Cross-session coordination | Cycle branch + `.cdd/unreleased/{N}/` artifacts are the coordination surface; roles read them at session start |

**Not current:** persistent polling role sessions where α, β, γ remain alive indefinitely between phases. That model requires dispatch-daemon infrastructure (#310, #295) not yet implemented. When such infrastructure exists, the polling instructions in §Tracking are already compatible — they are written for in-session use and work in either model.

**§1.6a Re-dispatch prompt formats.** Additional bounded re-dispatch prompts beyond the standard α/β/γ dispatches (§1.4 γ dispatch prompt format):

**α close-out re-dispatch (γ requests, δ executes after β merge):**
```
You are α. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view <number> --json title,body,state,comments
Branch: cycle/<number>
Close-out re-dispatch: β has merged and approved. Read .cdd/unreleased/<number>/beta-review.md for the approved verdict. Write .cdd/unreleased/<number>/alpha-closeout.md (cycle findings or "no findings"), commit to main, and exit.
```

**α fix-round re-dispatch (γ requests, δ executes when β returns RC):**
```
You are α. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view <number> --json title,body,state,comments
Branch: cycle/<number>
Fix-round re-dispatch: β returned RC in .cdd/unreleased/<number>/beta-review.md. Fix the named findings on cycle/<number>, append a fix-round section to .cdd/unreleased/<number>/self-coherence.md, push, and exit.
```

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
| 2 | Branch | Create a dedicated branch (γ creates `cycle/{N}` per §4.2 Branch rule + §1.4 γ algorithm Phase 1 step 3a; α/β check out, never create) | Valid branch name |
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

### 4.1a Lifecycle state table

The table below is the state machine for a substantial triadic cycle under the sequential bounded dispatch model (§1.6). Each state has one owner, required inputs, required outputs, a next state, and a failure/retry path.

| State | Owner | Required inputs | Required outputs | Next state | Failure / retry |
|---|---|---|---|---|---|
| **S0: Observed** | γ | CHANGELOG TSC, lag table, doctor/status, last PRA | Named selected gap + decisive CDD §3 clause | S1: Issue filed | Re-observe if no clear gap |
| **S1: Issue filed** | γ | Selected gap, issue-quality gate passed | Issue #N with ACs, Tier 3 skills, non-goals, related artifacts | S2: Branch created | Fix issue until quality gate passes |
| **S2: Branch created** | γ | Issue #N open; `origin/cycle/{N}` does not exist | `origin/cycle/{N}` exists, pre-flight passed | S3: α dispatched | Pre-flight fail → fix and retry |
| **S3: α dispatched** | δ | α prompt (from γ); `origin/cycle/{N}` exists | α session running | S4: α implementing | α session fails → re-dispatch α |
| **S4: α implementing** | α | Branch, issue, active skills loaded | Code, tests, docs on `cycle/{N}` | S5: α signaled review-ready | Blocker → γ unblocks; RC from β → re-dispatch α for fix round |
| **S5: α signaled review-ready** | α | Pre-review gate passed (§2.6) | `self-coherence.md` with review-readiness section on `cycle/{N}` | S6: β dispatched (if not already) | Gate fails → fix and re-signal |
| **S6: β reviewing** | β | `self-coherence.md` review-readiness, diff, CI green | `beta-review.md` round verdict (RC or A) | RC→S4; A→S7 | β blocked → γ unblocks |
| **S7: β merged** | β | Approved verdict, pre-merge gate passed | Merge commit on main with `Closes #N`; `beta-closeout.md` | S8: α close-out re-dispatch | Merge conflict → β resolves in throwaway worktree |
| **S8: α close-out** | α (re-dispatched) | `beta-review.md` (approved), merged branch state | `alpha-closeout.md` on main | S9: γ PRA | Re-dispatch unavailable → provisional close-out at review-readiness (declare as debt) |
| **S9: γ triaging** | γ | `alpha-closeout.md`, `beta-closeout.md`, RELEASE.md | PRA at canonical path, γ close-out triage, skill patches landed | S10: δ preflight | Missing close-out → request re-dispatch; missing RELEASE.md → γ writes it |
| **S10: δ release-boundary preflight** | δ | PRA present, RELEASE.md present, `.cdd/unreleased/{N}/` not yet moved, merge on main | Proceed / Request changes / Override | Proceed→S11; RC→route to β/α/override | δ blocks → γ routes change |
| **S11: γ closing** | γ | δ preflight passed; all closure gate rows pass (§gamma/SKILL.md §2.10) | `gamma-closeout.md`, cycle-dir move to `.cdd/releases/{X.Y.Z}/{N}/`, closure declaration | S12: δ disconnect | Missing artifact → obtain before closing |
| **S12: δ disconnect** | δ | γ closure declaration; `gamma-closeout.md` exists on main | Tag + release commit; branch cleanup | S0 (next cycle) | Script fails → fix and retry |

### 4.2 Branch rule

Every substantial change must be developed on its own dedicated branch. No substantial CDD work is performed directly on main.

**Canonical branch format:**

```text
cycle/{N}
```

where `{N}` is the issue number. One cycle = one branch = one named target for all polling.

**Ownership:** γ creates the branch from `origin/main` before dispatch (§1.4 γ algorithm Phase 1 step 3a). α and β never create branches — they `git switch cycle/{N}` from the name in their dispatch prompt and refuse to invent or accept any other name (§1.4 α step 5a, §1.4 β step 3, `alpha/SKILL.md`, `beta/SKILL.md` §1 Role Rule 1).

**Legacy shapes — warn-only.** The pre-#287 shapes are warn-only (verifiers may flag them; nothing else changes):

```text
{agent}/{version}-{issue}-{scope}        # canonical pre-#283
{agent}/{issue}-{scope}                  # canonical pre-#283 with version omitted
claude/{slug}-{rand}                     # harness-encoded (#274 / #283 friction class)
```

Frozen historical branches under those shapes are not retroactively renamed. Polling glob `'origin/claude/*'` is retained only for retrospective tracking (§Tracking) and is not a discovery surface for new cycles.

### 4.3 Branch pre-flight

γ runs branch pre-flight before creating `cycle/{N}` in §1.4 γ algorithm Phase 1 step 3a. Verify:

- `origin/cycle/{N}` does not yet exist (`git rev-parse --verify origin/cycle/{N}` returns non-zero — fail loud if the branch already exists, since one cycle = one branch)
- no stalled `.cdd/unreleased/{N}/` directory exists on `origin/main` (would indicate a previous cycle for `{N}` did not complete its release move per `release/SKILL.md` §2.5a)
- the issue's intended scope is declared in the issue body before the branch is created
- base SHA is known (`git rev-parse origin/main`)
- the issue is open

Pre-#287 ownership rested with α (α verified branch uniqueness, name format, CI state). Under the `cycle/{N}` rule γ owns pre-flight because γ creates the branch; α's pre-review gate retains its own row (`alpha/SKILL.md` §2.6 row 1) for verifying that `origin/cycle/{N}` is rebased onto current `origin/main` at review-readiness time, but α no longer creates the branch and does not own pre-flight at branch creation.

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

Artifacts outside version directories, such as `.cdd/unreleased/{N}/` files or navigation updates, are not required as bootstrap stubs.

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
| 2 | Branch | build | γ (creates `cycle/{N}` per §1.4 γ algorithm Phase 1 step 3a; α/β check out, never create) | `origin/cycle/{N}` exists with `main` as base | §4.2 / §4.3 | branch + CDD Trace row | mechanical | always | cdd |
| 3 | Bootstrap | build | α | version directory + manifest README + declared stubs | §5.1 | branch diff | mechanical | substantial only | cdd |
| 4 | Gap | build | α | named incoherence / coherence contract | `.cdd/unreleased/{N}/self-coherence.md` §Gap or design/SKILL.md §3.1 | primary branch artifact | agent | always | cdd |
| 5 | Mode | build | α | mode + active skills (+ bundle/level if used) | `.cdd/unreleased/{N}/self-coherence.md` §Mode or design/SKILL.md §3.1 | primary branch artifact | agent | always | cdd + applicable eng bundle from `src/packages/cnos.eng/skills/eng/README.md` |
| 6a | Design | build | α | design artifact or explicit "not required" | design/SKILL.md §3.1 | primary branch artifact | agent | substantial only | design |
| 6b | Plan | build | α | plan artifact or explicit "not required" | `docs/gamma/cdd/PLAN-TEMPLATE.md` | primary branch artifact or linked plan | agent | L7 / cycle-sized | plan |
| 6c | Tests | build | α (or delegated implementer) | test files or explicit reason none apply | — (diff) | diff / primary branch artifact | agent | always | eng/test |
| 6d | Code | build | α (or delegated implementer) | implementation diff or "docs/process only" | — (diff) | diff / primary branch artifact | agent | always | active loaded Tier 2/Tier 3 generation skills (language skill, runtime/platform skill, tooling skill — see `src/packages/cnos.eng/skills/eng/README.md`) |
| 6e | Docs | build | α (or delegated implementer) | changed canonical docs / specs / READMEs | — (diff) | diff | agent | when docs affected | eng/document for durable docs; eng/write-functional for functional/dataflow prose; `cnos.core/skills/skill` when authoring or modifying skills |
| 6f | Delegated handoff | build | α → implementer | implementation prompt with: active skills, test requirements per module, engineering conventions, artifact order + self-verification report from implementer | alpha/SKILL.md §2.2 | prompt + report | delegator + implementer | when implementation is delegated | cdd |
| 7 | Self-coherence | build | α (or delegated implementer) | `.cdd/unreleased/{N}/self-coherence.md` carrying gap, mode, ACs, evidence, known debt | `docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md` | `.cdd/unreleased/{N}/self-coherence.md` | agent | substantial only | cdd |
| 7a | Pre-review | build | α | branch rebased onto current `main`; `.cdd/unreleased/{N}/self-coherence.md` carries CDD Trace through step 7; tests reference ACs; known debt explicit; **schema/shape audit across all touched files** when contracts change — when introducing or changing a canonical form, verify (a) the new form is present across all relevant files AND (b) any superseded form has been removed; cite the verifying command (e.g. grep that returns exactly one match per file); **workspace-global library-name uniqueness check** when adding a new `(library (name X))` stanza; **post-patch re-audit** — after any mid-cycle code change, re-read `.cdd/unreleased/{N}/self-coherence.md` top-to-bottom and verify every CDD Trace / invariant / self-coherence row still matches HEAD; **branch CI green on head commit** before signaling review-readiness in the artifact | alpha/SKILL.md §2.6 | `.cdd/unreleased/{N}/self-coherence.md` | mechanical | always | cdd |
| 8 | Review | review | β | `.cdd/unreleased/{N}/beta-review.md` verdict + findings | review/SKILL.md output format | `.cdd/unreleased/{N}/beta-review.md` | reviewer | always | review |
| 9 | Merge | release | β | `git merge` of `cycle/{N}` into main (β authority — no δ required) | `docs/gamma/cdd/GATE-TEMPLATE.md` | release surface or `.cdd/unreleased/{N}/beta-review.md` | mechanical + reviewer | always | release |
| 10 | Release | release | β | CHANGELOG row, tag, release note | CHANGELOG.md ledger + release/SKILL.md | release surface | agent + mechanical | always | release + eng/document when authoring release notes/changelog prose |
| 11 | Observe | close | γ | post-release observation result | post-release/SKILL.md | post-release assessment | γ | always | post-release |
| 12 | Assess | close | γ | POST-RELEASE-ASSESSMENT.md | post-release/SKILL.md output template | version directory | γ | always | post-release |
| 12a | Skill patch | close | γ | skill/spec patches for recurring failure modes identified in close-out triage or PRA §3; synced across all affected surfaces under src/packages/ | post-release/SKILL.md §3 + gamma/SKILL.md §2.7–§2.10 | same commit as PRA when identified during assessment, otherwise γ close-out follow-on commit | γ | when close-outs or PRA identify a recurring failure or skill gap | post-release, cdd |
| 13 | Close | close | γ | immediate outputs executed (incl. 12a patches) + deferred committed | post-release/SKILL.md §6 CDD Closeout | post-release assessment | γ | always | post-release |

**Primary branch artifact:** `.cdd/unreleased/{N}/self-coherence.md` for triadic cycles, or the design artifact (design/SKILL.md §3.1) when a separate design doc is required for larger changes. The primary branch artifact owns the named incoherence, mode, active skills, AC mapping, CDD Trace, known debt, and the review-readiness signal.

**Role key (§1.4):** *γ (coordinator)* = steps 0–2 (observe, select, issue creation, **branch creation**, dispatch) + steps 11–15 (PRA, cycle iteration, δ preflight request, closure) + cycle-wide coordination, *α (implementer)* = steps 3–7a (bootstrap, gap, mode, artifacts, self-coherence, pre-review — α checks out the branch γ created in step 2 and never creates one), *β (reviewer + integrator)* = steps 8–10 (review RC/A decision, merge into main, β close-out — β polls `origin/cycle/{N}` directly and never creates a branch), *δ (operator)* = release-boundary preflight (step 15) + disconnect release (step 19: tag, deploy). Delegated implementer is α-side. Merge is β's authority (step 9). Tag/release/deploy is δ's authority (step 19).

**Producer key:** *agent* = judgment required, *mechanical* = automatable by cnos (#94), *reviewer* = produced by the review process.

Rules:
- "Not required" is valid only when stated explicitly.
- An omitted step with no explicit note is incomplete, not implicit.
- Small-change mode may collapse steps 4–7 into commit-message evidence, but the same distinctions still apply.

### 5.3a Artifact Location Matrix

For release-scoped triadic cycles, every named artifact has exactly one canonical location. Verifiers (e.g. `cn cdd-verify`) enforce these paths as canonical and treat any other location as legacy/warn-only.

| Artifact | Canonical repo location | CDD package default | Noncanonical / legacy / scratch |
|---|---|---|---|
| Version snapshot directory | `docs/{tier}/{bundle}/{X.Y.Z}/` | `docs/gamma/cdd/{X.Y.Z}/` | `.cdd/releases/{X.Y.Z}/` is **not** the frozen snapshot — it is triadic protocol/scratch space |
| POST-RELEASE-ASSESSMENT.md (PRA) | `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `.cdd/releases/{X.Y.Z}/beta/POST-RELEASE-ASSESSMENT.md` and `.cdd/releases/{X.Y.Z}/beta/ASSESSMENT.md` are legacy/warn-only |
| α self-coherence (primary branch artifact) | `.cdd/unreleased/{N}/self-coherence.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/self-coherence.md` at release per `release/SKILL.md` §2.5a | same | none — required for every triadic cycle |
| β review record | `.cdd/unreleased/{N}/beta-review.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/beta-review.md` at release | same | none — required for every triadic cycle |
| α close-out | `.cdd/unreleased/{N}/alpha-closeout.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/alpha-closeout.md` at release | same | `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` (legacy aggregate form, pre-#283) is warn-only |
| β close-out | `.cdd/unreleased/{N}/beta-closeout.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/beta-closeout.md` at release | same | `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md` (legacy aggregate form, pre-#283) is warn-only |
| γ close-out | `.cdd/unreleased/{N}/gamma-closeout.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/gamma-closeout.md` at release | same | `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md` (legacy aggregate form, pre-#283) is warn-only |
| γ kata verdict (optional) | `.cdd/releases/{X.Y.Z}/gamma/KATA-VERDICT.md` when kata is available | same | warn-only when kata unavailable |
| CHANGELOG ledger row | `CHANGELOG.md` (Release Coherence Ledger) — must include Version, C_Σ, α, β, γ, **Level**, coherence note | same | none |
| RELEASE.md | `RELEASE.md` at repo root, included in the release commit | same | CI auto-generated body is not an acceptable substitute |
| Hub memory | external agent-hub state (not a repo artifact) | external agent-hub state | the PRA must record path / commit-sha / unavailable-reason; verifiers do not inspect the external hub directly |

Rules:

- `.cdd/unreleased/{N}/` is the per-cycle coordination directory, keyed by issue number. Files inside are role-distinguished by **filename** — see §Tracking for the canonical filename set (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`) and the convention for adding more (`{role}-{purpose}.md` for role-owned content; bare `{purpose}.md` when convention identifies the author).
- `.cdd/releases/{X.Y.Z}/{N}/` holds the moved-at-release form of each cycle's coordination directory. Multiple cycle directories may live under one release directory when several issues ship in the same release.
- `.cdd/` is triadic protocol space and role-local close-out / close-out evidence storage. It is **not** the canonical frozen post-release snapshot.
- `docs/{tier}/{bundle}/{X.Y.Z}/` is the canonical frozen post-release snapshot directory. The PRA lives there; close-outs do not.
- Tags are bare `X.Y.Z` everywhere (VERSION file, git tag, branch name version segment, CHANGELOG row, RELEASE.md, snapshot directory). `v`-prefixed tags are legacy and warn-only.
- All triadic cycles use the in-version `.cdd/unreleased/{N}/` artifact-exchange surface. GitHub Pull Requests are not used for triadic coordination (see §Tracking).

### 5.3b Role/artifact ownership matrix

Every required cycle artifact has one owner, one verification gate, and one consequence if missing. γ's closure gate (§gamma/SKILL.md §2.10) checks every row marked "Required before γ closure."

| Artifact | Owner | Written when | Verified by | Required before | Missing means |
|---|---|---|---|---|---|
| `self-coherence.md` | α | During α session, incrementally; review-readiness section last | β at review intake | β review | β waits; no review until present with readiness signal |
| `beta-review.md` | β | During β review session, incrementally per round | γ at close-out triage | γ closure | γ cannot triage; requests β re-dispatch |
| `alpha-closeout.md` | α (re-dispatched after merge) | After β merge, via δ re-dispatch (§1.4 α step 10, §1.6a) | γ before PRA; γ closure gate | γ closure | γ closure gate blocks; γ requests δ to re-dispatch α |
| `beta-closeout.md` | β | Before β exits (same β session as merge) | γ before PRA; γ closure gate | γ closure | γ closure gate blocks; γ requests β re-dispatch |
| `gamma-closeout.md` | γ | After all closure gate rows pass (§gamma/SKILL.md §2.10) | δ before tag/release (implicit: tag requires γ closure declaration which is gamma-closeout.md) | δ tag/release | δ must not tag; γ has not declared closure |
| `RELEASE.md` | γ | Before requesting δ tag/release; committed to main in release commit | δ at release-boundary preflight | δ tag/release | δ must not tag; CI auto-generates sparse notes |
| `.cdd/releases/{X.Y.Z}/{N}/` | γ (at release, via `release/SKILL.md` §2.5a) | Before γ requests δ tag; included in release commit | γ closure gate row 12; δ preflight | δ tag/release | Stale unreleased dir; γ closure gate blocks until moved |
| POST-RELEASE-ASSESSMENT.md | γ | After β merge + close-outs | γ closure gate; δ at release-boundary preflight | γ closure | γ closure gate blocks |
| `beta-review.md` fix-round appendix | β | After α fixes RC findings | α at next fix-round | α re-signal | α waits for β's round verdict |

Notes:
- `gamma-closeout.md` is the closure declaration artifact. δ's obligation to not tag before closure is satisfied when `gamma-closeout.md` exists and the closure declaration commit is on main.
- For small-change cycles: see §1.2 artifact collapse table — `beta-closeout.md` and `gamma-closeout.md` may not apply; `alpha-closeout.md`, `RELEASE.md`, and the cycle-directory move remain required per §1.2.

### 5.4 CDD Trace

Every substantial cycle must carry a lightweight execution trace. Use lifecycle step numbers, not section numbers. For steps 0–10, the trace lives in the primary branch artifact. For steps 11–13, closure lives in the post-release assessment.

The primary branch artifact is the artifact that owns the named incoherence, mode, active skills, and acceptance criteria. For triadic cycles this is `.cdd/unreleased/{N}/self-coherence.md`. For governance/process work the governing doc being changed may carry the trace inline. When a separate design artifact exists (design/SKILL.md §3.1), the design artifact carries the trace and `self-coherence.md` references it.

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
| 7 Self-coherence | `.cdd/unreleased/{N}/self-coherence.md` | cdd | AC-by-AC self-check completed |
| 7a Pre-review | `.cdd/unreleased/{N}/self-coherence.md` | cdd | Pre-review gate passed (alpha/SKILL.md §2.6); review-readiness signaled |
| 8 Review | `.cdd/unreleased/{N}/beta-review.md` | review (+ others if loaded) | CLP review result |
| 9 Gate | `.cdd/unreleased/{N}/beta-review.md` or release surface | release (+ eng/document if release prose authored) | Release-readiness decision |
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

### 8.1 Closure verification checklist

This checklist catches the failure modes that motivated #325. Any reviewer (β, γ, δ) may run through it mechanically before the cycle closes. Each row names the failure mode, the check, and where the authoritative gate lives.

| # | Failure mode | Check | Authoritative gate |
|---|---|---|---|
| F1 | Missing `alpha-closeout.md` | `.cdd/unreleased/{N}/alpha-closeout.md` exists on main | §5.3b; `gamma/SKILL.md` §2.10 row 1 |
| F2 | Missing `beta-closeout.md` | `.cdd/unreleased/{N}/beta-closeout.md` exists on main | §5.3b; `gamma/SKILL.md` §2.10 row 2 |
| F3 | Missing `gamma-closeout.md` before δ tag | `.cdd/unreleased/{N}/gamma-closeout.md` exists on main before δ runs `scripts/release.sh` | §5.3b; `operator/SKILL.md` §3.4 |
| F4 | Stale `.cdd/unreleased/{N}/` after release | `git ls-tree -r origin/main .cdd/unreleased/ --name-only` returns empty (or only in-progress cycles) after tag | `release/SKILL.md` §2.5a; `gamma/SKILL.md` §2.6 + §2.10 row 12 |
| F5 | Missing `RELEASE.md` before tag | `RELEASE.md` exists at repo root, committed in or before the release commit | `release/SKILL.md` §3.7; `gamma/SKILL.md` §2.6; §5.3b |
| F6 | δ tag before γ closure declaration | `gamma-closeout.md` exists before δ runs release script | §5.3b; §4.1a state S12 requires S11 (γ closing) complete |
| F7 | δ tag before δ preflight | δ preflight (state S10) was run and returned Proceed before tag | §4.1a S10→S11→S12 ordering |
| F8 | α close-out relies on already-exited α | α close-out mechanism used: either re-dispatch (§1.6a) or provisional fallback declared in §Debt | §5.3b; §1.4 α step 10; `alpha/SKILL.md` §2.8 |
| F9 | Polling-required role exits early | Role exited per bounded-dispatch model (§1.6); any fix-round or close-out needed from exited α goes through re-dispatch (§1.6a) | §1.6; `operator/SKILL.md` §1 |
| F10 | PRA missing after release | `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` exists on main | `post-release/SKILL.md` pre-publish gate; `gamma/SKILL.md` §2.10 row 3 |

**Positive test:** A substantial cycle that reaches S12 (δ disconnect) with all 10 rows above passing has closed correctly.

**Negative test:** Any row that fails when γ reaches state S11 (γ closing) must block closure. γ must not write the closure declaration commit until all rows pass.

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
