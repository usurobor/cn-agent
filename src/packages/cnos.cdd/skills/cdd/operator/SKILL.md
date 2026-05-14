---
name: operator
description: δ role in CDD. Owns external gates, routing, and override authority during an active triad cycle.
artifact_class: skill
kata_surface: embedded
governing_question: What does the operator do — and not do — during an active CDD cycle?
visibility: internal
parent: cdd
triggers:
  - operator
  - dispatch
  - gate
  - unblock
scope: role-local
inputs:
  - γ dispatch prompts
  - external gate requests from α/β/γ
  - cycle state (issue, branch, .cdd/unreleased/{N}/, CI)
outputs:
  - routed prompts to agent sessions
  - external gate executions and confirmations (push on β's behalf, tag, release, deploy, issue filing)
  - override declarations when needed
requires:
  - active CDD cycle exists
  - γ has produced dispatch prompts
calls:
  - gamma/SKILL.md
---

# Operator (δ)

## Core Principle

**Coherent δ operation routes work, holds gates, and stays out of the triad's reasoning.**

δ is not a fourth triad role — δ is not scored on coherence axes. δ owns what agents cannot: external platform actions, session routing, and override authority. The failure mode is **invisible meddling** — the operator adjusts implementation, review, or coordination reasoning without declaring an override, and the triad's coherence record no longer matches what actually happened.

> **First-time operator?** Before running the algorithm below, ensure the repository has been activated under CDD.
> See [`cdd/activation/SKILL.md`](../activation/SKILL.md) for the one-time bootstrap: `.cdd/` scaffold, version
> pin, identity convention setup, and the §24 verification command that confirms activation is complete.

## Algorithm

1. **Dispatch γ** — run γ via `claude -p`; γ reads the issue, creates the cycle branch, and produces α/β prompts.
2. **Dispatch α** — run α via `claude -p` with the prompt γ produced; α implements on the cycle branch and exits after signaling review-readiness.
3. **Dispatch β** — run β via `claude -p` with the prompt γ produced; β reviews, merges, writes β close-out, and exits.
4. **Re-dispatch α for fix rounds** (when β returns RC) — run α via `claude -p` with the fix-round re-dispatch prompt (CDD.md §1.6a); α fixes findings, appends fix-round to self-coherence.md, exits.
5. **Re-dispatch α for close-out** (when γ requests after β merge) — run α via `claude -p` with the close-out re-dispatch prompt (CDD.md §1.6a); α writes alpha-closeout.md, commits to main, exits. **This step is mandatory when γ requests it.** γ cannot complete the closure gate without alpha-closeout.md.
6. **Gate** — execute external actions: push main, tag, release, branch cleanup. **Do not tag/release before `gamma-closeout.md` exists on main.** After tag push, δ runs `gh run list --branch <tag>` and waits for release workflow completion. **δ blocks release completion until CI is green and owns recovery on red:**
   - **CI Green** → δ declares release complete  
   - **CI Red** → δ owns the failure and executes recovery runbook:
     1. **Investigate** release logs to understand the failure
     2. **Classify** failure as release-specific vs pre-existing infrastructure  
     3. **Fix or escalate** — if fixable: fix, re-tag/re-run, poll again; if pre-existing: document, escalate to operator, do NOT declare release complete
     4. **Re-verify** — poll CI again after fix attempts
     5. **Operator override** — explicit operator acceptance required for known pre-existing failures (escape hatch for cases like v3.66.0/v3.67.0 smoke failures)
   
   **The gate does not close until CI is green or operator explicitly accepts the failure.**
7. **Override** — reassign roles or redirect scope only with an explicit declaration.

δ runs one role at a time. This keeps memory pressure low (single `claude -p` process), gives δ direct visibility into each session, and isolates failures — if α dies, δ retries α without losing γ or β state.

---

## Git identity for role actors

Every CDD role actor configures a git identity in the form `{role}@{project}.cdd.cnos` before making any commits on the cycle branch. DNS domains read broad-to-narrow right-to-left: `cnos` is the origin repository where the cdd protocol is defined and versioned, `cdd` is the protocol namespace inside cnos, and `{project}` is the tenant project running the protocol. The role name is the local part. This form makes the protocol's origin repo visible in every commit trailer and leaves namespace room for sibling protocols (`cnav`, `cnobs`) under the same cnos root.

**Special case — cnos itself.** When the project running the cycle is the cnos repo, the literal form would be `{role}@cnos.cdd.cnos` (redundant `cnos`). The canonical elision is `{role}@cdd.cnos`, which reads as "the cdd protocol at cnos." Existing cnos commit trailers already use this form; the redundancy adds no information.

| Role | Project | Canonical identity | Notes |
|------|---------|-------------------|-------|
| alpha | tsc | `alpha@tsc.cdd.cnos` | tsc project actor |
| beta | cnos | `beta@cdd.cnos` | cnos actor — elision form (see above) |
| gamma | acme | `gamma@acme.cdd.cnos` | hypothetical third project |
| beta | * | `beta@cdd.{project}` | **(deprecated)** — cycle #287 form; cycle #343 cutover |

Set identity before the first commit of each dispatch session:

```bash
# general form (non-cnos projects):
git config user.name "{role}"
git config user.email "{role}@{project}.cdd.cnos"

# cnos project (elision form):
git config user.name "{role}"
git config user.email "{role}@cdd.cnos"
```

---

## 1. Route

### 1.1. Dispatch γ first

δ dispatches γ via `claude -p`. γ reads the issue, creates the cycle branch, and returns α/β prompts to δ. γ does not execute dispatch — δ does.

```bash
cat /tmp/gamma-prompt.md | claude -p --allowedTools "Read,Write,Bash" --output-format stream-json --verbose --model <model>
```

### 1.2. Dispatch α and β sequentially

δ dispatches α, waits for completion, then dispatches β. One `claude -p` at a time.

```bash
# α — implements
cat /tmp/alpha-prompt.md | claude -p --allowedTools "Read,Write,Bash" --output-format stream-json --verbose --permission-mode acceptEdits --model <model>

# β — reviews (needs Bash for git/gh read-only commands)
cat /tmp/beta-prompt.md | claude -p --allowedTools "Read,Write,Bash" --output-format stream-json --verbose --permission-mode acceptEdits --model <model>
```

Note: `--output-format stream-json --verbose` is required for all dispatches — without it, δ cannot monitor agent output in real time. (`--verbose` is mandatory when combining `--output-format stream-json` with `-p`; without it, `claude` exits with an error.) `--permission-mode acceptEdits` is required because `claude -p` as a fresh session hits the trust dialog. Without it, agents cannot write files. β gets Bash because it needs `git diff`, `gh issue view`, etc. — role boundaries are enforced by beta/SKILL.md, not tool scoping.

- ❌ Rewrite the prompt to add constraints or context γ didn't include
- ✅ Deliver the prompt verbatim to the `claude -p` session
- ❌ Run α and β concurrently or nest them inside γ's session
- ✅ Run one role at a time, inspect artifacts between dispatches
- ✅ Name the agent-to-role mapping before delivering prompts

---

## 2. Wait

### 2.1. Do not poll internal work

Once dispatched, the triad runs. The operator does not need to monitor branch diffs, `.cdd/unreleased/{N}/` files, or CI runs — γ owns that.

The operator's wake-up signals are:

- **γ requests an external gate** (push on β's behalf, tag push, release, deploy, issue filing, auth action)
- **γ requests an unblock decision** (design ambiguity, scope question, env constraint)
- **γ declares cycle closure** and names deferred outputs that need operator action
- **An agent session dies or stalls** (no activity for an unexpected duration)

Between these signals, the operator's correct action is nothing.

- ❌ Check the branch every 30 minutes and suggest changes (role leak into β)
- ❌ Heartbeat reveals cycle state (tag missing, branches exist) → δ acts on it (observation is not a gate request)
- ✅ Wait for γ to surface a gate or decision request
- ✅ Heartbeat reveals cycle state → δ notes it, waits for γ

### 2.2. Subscribe to the issue

Poll the issue for activity using the same transition-only pattern as the triad (CDD §1.4). δ polls less frequently — the wake-up signals are coarser (gate requests, not per-commit state).

```bash
prev=""; while true; do
  cur="$(cd /path/to/repo && git fetch --quiet origin && gh issue view <N> --json comments --jq '.comments | length')"
  if [ "$cur" != "$prev" ]; then echo "issue-activity: comments=$cur"; prev="$cur"; fi
  sleep 300
done
```

Run under `Monitor` or equivalent. 5-minute interval is sufficient for δ — γ owns the tight loop. Supplement with branch + `.cdd/unreleased/` polling once cycles are active. Canonical cycle branches are `origin/cycle/{N}` (per `CDD.md` §4.2, since #287). The pre-#287 `'origin/claude/*'` glob is **warn-only / retrospective** — retained for tracking historical cycles whose branches predate the rule, never as a discovery surface for new cycles:

```bash
prev_branches=""; declare -A prev_head
while true; do
  cd /path/to/repo && git fetch --quiet origin
  # Canonical: cycle/{N} branches (γ creates these per CDD.md §1.4 γ algorithm Phase 1 step 3a).
  cur_branches="$(git branch -r --list 'origin/cycle/*' 2>/dev/null | sed 's| ||g' | sort)"
  comm -13 <(echo "$prev_branches") <(echo "$cur_branches") | sed 's/^/new-branch: /'
  # Per-branch head SHA — cycle artifacts live on cycle branches, not on main.
  for b in $cur_branches; do
    cur_head="$(git rev-parse "$b" 2>/dev/null)"
    [ "$cur_head" != "${prev_head[$b]:-}" ] && [ -n "$cur_head" ] && echo "branch-update: $b → $cur_head"
    prev_head[$b]="$cur_head"
  done
  prev_branches="$cur_branches"
  sleep 300
done
# To track legacy branches retrospectively, swap the glob to 'origin/claude/*'
# (warn-only — pre-#287 cycles only). Do not use for new cycles.
```

---

## 3. Gate

### 3.1. External actions the operator holds

These actions require platform permissions agents may lack:

| Action | Trigger | Who requests |
|--------|---------|-------------|
| Pre-merge gate validation | Before authorizing β merge, run `scripts/validate-release-gate.sh --mode pre-merge` to verify cycle artifacts exist and are well-formed. See `CDD.md` §5.3b and `gamma/SKILL.md` §2.10. | γ |
| Push β-approved merge to main | β runs `git merge` — δ only pushes when β cannot execute the push directly (env/auth constraint). This is execution of β's integration authority, not δ approval. | β or γ |
| Release-boundary preflight | After β merge + close-outs + γ PRA, δ verifies merge commit, release artifacts, tag/deploy preconditions, and platform readiness. Proceed / request changes / override. See CDD §1.4 Phase 5a. | γ |
| Tag push + release | After δ preflight confirms and γ closes the cycle. **δ is sole tag-author** — β does not tag; only δ creates tags per cycle | γ |
| Branch delete | Cycle closed, merged branches | γ |
| Issue filing on external repos | Cross-project dependency | γ |
| Force push | Rebase required with env constraints | α via γ |
| Auth refresh | Token/permission expiry | any role via γ |

### 3.2. Execute on request, not on observation

Gate actions fire when a role requests them, not when δ notices they're needed. Observing that a tag isn't pushed or a branch exists is not a gate trigger — γ's explicit request is.

- ❌ Heartbeat shows tag not pushed → δ pushes it (role leak: δ decided the gate, not γ)
- ❌ "β asked me to push the merge but I think we should wait for one more review"
- ✅ γ requests tag push → δ pushes it and confirms
- ✅ If you disagree with a gate request, declare an override (§4)

### 3.3. Report completion

After executing a gate action, confirm to the requesting role that the action completed.

- ❌ Execute silently and assume the triad will notice
- ✅ "Tag pushed: `git push origin 3.59.0` — confirmed on remote"

### 3.4. Cut the release — disconnect the triad's final state

After all post-cycle work lands on main (γ's PRA + skill patches, δ's own session patches), δ cuts the release. This is not optional — the release is how δ disconnects the triad's output into a distributable, tagged whole.

The triad's work is not complete until it is tagged. Untagged post-cycle patches on main are an open boundary — the triad's output is still entangled with whatever comes next. The tag is the disconnection point.

**Algorithm:**
1. Confirm all post-cycle commits are on main (γ PRA, γ skill patches, δ session patches)
2. Edit VERSION to the new number
3. Run `scripts/release.sh` — this stamps all manifests, verifies consistency, generates structured tag messages, commits, creates annotated tags, and pushes in one command
4. **Poll release CI** — after tag push, run `gh run list --branch <tag>` and monitor release workflow completion. **δ blocks until CI is green, owns recovery on red** (see §6 Gate for full recovery runbook). δ may NOT declare release complete while CI is red. Requires CI green or explicit operator override.
5. **Branch cleanup** — delete merged cycle branches (`cycle/{N}`) and γ session branches (harness-given `claude/...` or operator-named `gamma/session-{N}` patterns) that were used during the cycle. No orphan γ session branches survive past closure.

**Manual tagging is not allowed.** Do not run `git tag` directly. The release script is the only way to tag. It prevents the class of failures where VERSION, cn.json, and package manifests disagree (see DISPATCH-FAILURE-EVIDENCE.md, cycle #84 failure 3).

**Tag message generation:** The release script automatically generates structured annotated tag messages that include issue metadata, wave context when present, and CDD review artifacts when available. The generated message is deterministically derived from git history, GitHub metadata, and CDD artifacts. δ can inspect the tag message content after tagging using `git show <version>` or `git for-each-ref refs/tags/<version> --format='%(contents)'`.

```bash
# Edit VERSION, then:
scripts/release.sh
# Or pass version directly:
scripts/release.sh 3.67.0
```

- ❌ `git tag 3.67.0 && git push --tags` (skips stamp, skips consistency check)
- ❌ γ's skill patches sit on main untagged across multiple cycles
- ❌ δ defers release "because there are no consumers" (the tag is structural, not consumer-driven)
- ✅ `scripts/release.sh 3.67.0` — stamps, verifies, commits, tags, pushes

### 3.5. The tag is the signal

The disconnect tag (§3.4) is git-observable. γ and all future agents can see it. No separate completion signal is needed — the tag appearing on main IS the proof that all gate actions completed and the cycle is disconnected.

For mid-cycle gate actions (tag push before the disconnect, branch cleanup), confirm completion to the requesting role per §3.3. But the disconnect tag itself needs no announcement — it speaks for itself.

---

## 4. Override

### 4.1. When to override

The operator may override any triad decision when:

- The triad's direction conflicts with project priorities the triad cannot see
- A role is stuck in a loop and γ's unblocking hasn't resolved it
- External constraints force a scope change (e.g. deadline, dependency shift)
- Safety or security requires immediate action

### 4.2. Override protocol

Per CDD §1.4: the reassignment must name the target agent and the reason. No implicit drift.

State:
1. What you are overriding (role assignment, scope, priority, decision)
2. Why (the information or constraint the triad didn't have)
3. What the new state is

- ❌ Edit the issue quietly and let α discover the change
- ✅ "Override: descoping AC4–AC5 from this cycle. Reason: timeline constraint for release by Friday. New scope: AC1–AC3 only. γ to update the issue."

### 4.3. Do not use override for taste

Override is for information asymmetry or hard constraints, not preference.

- ❌ "I think the implementation should use a different approach" → that's α's job
- ❌ "The review was too harsh" → that's β's judgment
- ✅ "The API we're building against is being deprecated next week" → information the triad needs

---

## 5. Dispatch configurations

CDD defines two valid dispatch configurations. Choose one before dispatching; record it in `gamma-closeout.md`.

### 5.1 Canonical multi-session dispatch

One `claude -p` process per role; each has an independent auth context and no shared memory. This is the model described in §1.2 above.

- γ/δ separation is structurally present: the operator (δ) selects and scaffolds; γ coordinates; α and β are separate processes with no access to each other's reasoning or conversation state.
- Sub-agent returns do not apply — each `claude -p` session exits cleanly; the operator reads committed artifacts.
- Branch names are stable: `cycle/{N}` persists through all fix rounds because each role session checks it out fresh.

Use this configuration when the cycle is substantial (see §5.3 escalation criteria).

### 5.2 Single-session δ-as-γ via Agent tool (Claude Code activation)

When the operator is a Claude Code agent (one parent session), α and β are dispatched as sub-agents using the Agent tool rather than via separate `claude -p` processes. Sub-agents run with fresh context per invocation and are functionally equivalent to `claude -p` for role-isolation purposes (each sub-agent reasons independently, cannot see the parent's conversation state, and cannot see the other sub-agent's conversation state). However, sub-agents inherit MCP scope and filesystem access from the parent session.

**Scope of the collapse.** §5.2 collapses **δ↔γ only**. γ↔α↔β remain structurally separate per `CDD.md §1.4` Triadic rule: γ scaffolds and coordinates in the parent session, α implements in its own sub-agent, β reviews and merges in its own sub-agent. The dyad-plus-coordinator structure is preserved; only the operator (δ) and coordinator (γ) functions fuse into one parent session.

A single sub-agent that performs γ-selection plus α-implementation plus β-review is not §5.2 — it is a §1.4 violation. §5.2 requires three execution contexts: the parent session (γ, also δ), a separate α sub-agent, and a separate β sub-agent. Lumping γ+α+β into one sub-agent breaks role-isolation (α gains access to β's reasoning and vice versa) and is rejected.

Three structural consequences follow:

1. **δ=γ collapse.** γ/δ separation is structurally absent: one parent session holds both the selection/scaffolding function (γ) and the external gate authority (δ). The cycle proceeds, but the grading floor in §3.8 of `release/SKILL.md` applies — see §3.8 configuration-floor clause.

For the general principle governing when role collapses are safe, see `ROLES.md §4` (hats-vs-actors: independence as the collapse constraint).

2. **sub-agent return messages are summaries, not full transcripts.** The Agent tool returns a summary message from the sub-agent, not a full conversation transcript. δ-as-γ verifies committed artifacts — specifically `beta-review.md` for β's verdict and `self-coherence.md` for α's review-readiness signal — rather than relying on the sub-agent's return message. This is the protocol invariant that makes §5.2 valid despite the summary-not-transcript limitation: the artifact β commits is canonical; the return message is informational only.

3. **Harness push restrictions surface as branch-name churn under fix-rounds.** Some Claude Code harness environments block updates to existing remote branches (403 on push to a previously-written branch). Under this constraint, a fresh-branch chain (`cycle/{N}` → `cycle/{N}-impl` → `cycle/{N}-impl-r2` → `cycle/{N}-impl-r3` → `cycle/{N}-merged` → `cycle/{N}-final`) is an acceptable workaround. Each link in the chain is a valid cycle branch for that fix-round; the final fast-forward into `main` becomes an external operator action.

**Empirical anchors:** The cnos-tsc supercycle (cycles 24–26, close-outs at `usurobor/tsc:.cdd/releases/{0.5.0,0.6.0,0.7.0}/{24,25,26}/gamma-closeout.md`) ran under §5.2 end-to-end; tsc cycle 26 γ-closeout explicitly records "operator (δ = γ in this two-agent configuration)." Tsc cycle #32 (close-out at `usurobor/tsc:.cdd/releases/docs/2026-05-09/32/gamma-closeout.md`) ran §5.2 and produced a five-link branch trail (`cycle/32` → `cycle/32-impl` → `cycle/32-impl-r2` → `cycle/32-merged` → `cycle/32-final`) due to harness push restrictions; that trail is the source observation for consequence (3) above.

**Reference dependencies:** §5.2 dispatch sizing follows `CDD.md §1.6c` (sub-agent dispatch budget heuristic). The harness push restriction that produces branch-name churn is the same constraint that makes the mechanical pre-merge gate (`release/SKILL.md §2.1`) an operator-side action when β cannot push directly (`CDD.md §1.4 β algorithm step 8`).

### 5.2.1 Parent-session quiescence during sub-agent runs

When a sub-agent dispatch is in flight (via the Agent tool), the parent session MUST enter quiescent mode:

**Prohibited during sub-agent runs:**
- **No file edits** in the working tree
- **No commits** from the parent session
- **No branch switches** (`git checkout`, `git switch`)
- **No `git add` / `git restore --staged`** (index state must remain stable for the sub-agent's view)
- **No `git pull` / `git fetch` that updates the current branch HEAD**

**Permitted during sub-agent runs:**
- Reads: `git status`, `git log`, `git diff`, file reads via Read tool, web fetches, GitHub MCP queries
- Coordination: dispatching additional sub-agents (each gets its own context; isolation is by sub-agent), reading existing branches via `git fetch <branch>` (does not update HEAD)
- **Filesystem operations on paths outside the repo** (e.g., `/tmp`, drafts)

**When the parent must edit:** dispatch a sub-agent to do the edit, or wait for current sub-agents to complete. The parent is the coordinator, not a fourth concurrent writer.

**What goes wrong when this is violated:** Sub-agent A is running α R1, has read the working tree, is about to `git add` and `git commit`. Parent session edits file X concurrently. Sub-agent's `git add .` picks up parent's edits accidentally. Sub-agent's commit message describes only its intended change but the diff includes parent's edits. Result: corrupted commit, branch must be reset, fix-round required.

**Sub-agent parallelism note:** Multiple sub-agents launched in parallel (via `Agent` tool calls in one parent message) have **isolated contexts** from each other but **share the working tree with the parent and with each other**. Concurrent file edits by parallel sub-agents are also a corruption risk and should be avoided — multi-sub-agent parallelism is for independent reads, not concurrent writes.

**Exception:** When the Agent tool runs with `isolation: "worktree"`, parent-session quiescence is unnecessary (the sub-agent operates on a copy of the repo). Default mode requires quiescence.

### 5.3 Escalation criteria

Switch from §5.2 to §5.1 (multi-session) when any of the following is true at γ scaffold time:

- **≥7 ACs** in the cycle (complexity exceeds what a single-session coordinator can reliably hold across sub-agent round-trips)
- **New contract surface or cross-repo deliverables** (γ/δ separation matters for audit when external systems are affected)
- **≥3 β rounds expected** (multiple fix-round re-dispatches compound the branch-name churn problem under §5.2 harness constraints)
- **≥3 γ judgment calls expected mid-cycle** (design ambiguity, scope changes, unblock decisions — each requires full cross-agent context that δ-as-γ holds less reliably than a dedicated γ session)

---

## 6. What the operator does NOT do

These are role boundaries. Crossing them without an override declaration breaks the coherence record.

- **Do not implement.** That's α. If you want to fix something, declare an override and take the α role explicitly.
- **Do not review.** That's β. If you disagree with β's verdict, declare an override.
- **Do not triage findings.** That's γ. If you want a finding handled differently, tell γ.
- **Do not rewrite prompts.** γ owns prompt quality. If the prompt is weak, tell γ to fix the issue.
- **Do not merge without β's approval.** The merge gate fires when β says it's ready, not before.
- **Do not communicate directly with α or β during in-cycle work.** γ is the bridge. δ-to-γ, γ-to-α/β. δ remains external to the triad; γ is the operator's only in-cycle counterparty.

---

## 7. Cycle lifecycle from the operator's view

| Phase | Operator action | Wait for |
|-------|----------------|----------|
| γ dispatch | Run γ via `claude -p`; γ creates branch, returns α/β prompts | γ completion |
| α dispatch | Run α via `claude -p` with γ's prompt | α completion (exits after review-readiness) |
| β dispatch | Run β via `claude -p` with γ's prompt | β completion (merge + β close-out) |
| α fix-round re-dispatch | Run α via `claude -p` with fix-round prompt (CDD.md §1.6a) when β returns RC | α completion (exits after fix-round) |
| α close-out re-dispatch | Run α via `claude -p` with close-out prompt (CDD.md §1.6a) when γ requests | α completion (alpha-closeout.md on main) |
| Release prep | γ writes RELEASE.md, moves cycle dirs; δ holds until complete | γ request |
| δ preflight | Verify merge commit, release artifacts, tag preconditions | γ preflight request |
| Closure | Gate: do not tag before `gamma-closeout.md` exists on main | γ closure declaration (gamma-closeout.md) |
| Disconnect | Cut the release — `scripts/release.sh` after γ closure declaration | γ close-out + δ session patches on main |
| Post-release | Execute deferred operator actions from γ close-out | γ deferred-output list |
| Inter-cycle | Nothing until next γ dispatch | γ next-cycle selection |

---

## 8. Timeout recovery

### §timeout-recovery — What to do when an agent session terminates before committing

An agent dispatched via `claude -p` may be SIGTERM'd by the OS, hit the session timeout, or crash without committing its in-progress work. This section is the executable recovery procedure.

#### 7.1 Inspect the worktree

Run these commands in the repo root to assess what the agent produced before it died:

```bash
# 1. What is staged or modified but not committed?
git status --short

# 2. What files were written during the session?
#    Replace $session_start with the approximate session start time.
find . -newer /tmp/session-start-sentinel -not -path './.git/*' | sort

# 3. Did the agent stash anything?
git stash list

# 4. What is the diff against main?
git diff origin/main..HEAD
git diff HEAD   # staged + unstaged against last commit
```

To create the sentinel file before dispatching, run `touch /tmp/session-start-sentinel` immediately before `claude -p`.

#### 7.2 Decision tree

After inspecting the worktree, choose one path:

| Situation | Action |
|---|---|
| Work is committed — agent wrote one or more commits before SIGTERM | Normal recovery: the agent left durable checkpoints. Read `origin/cycle/{N}` to see what landed. No further action needed unless the cycle is incomplete. |
| Work is staged / unstaged but not committed | Commit under agent identity: `git config user.name "Alpha" && git config user.email "alpha@cdd.cnos"`, then `git add <files> && git commit -m "<msg>"`. Preserve the agent's intent; do not rewrite scope. |
| Work exists only as new/modified files (never staged) | Same as above — stage and commit under agent identity with a commit message that names the recovery context (e.g., `"recovery(338): α partial work — committed by operator after SIGTERM"`). |
| Nothing useful recovered — session started but produced no artifact content | Declare a failed dispatch. File the failure in the cycle's self-coherence.md §Debt. Re-dispatch α with a fresh budget per §1.6c(a). The override is operator-identity if re-dispatch is not available. |
| Stash exists | Inspect with `git stash show -p`. Pop and commit if the content is relevant: `git stash pop && git add <files> && git commit`. |

**Agent-identity vs operator-identity commit:** Prefer committing under the agent's canonical identity (`alpha@cdd.cnos`) to preserve the role-identity-is-git-observable property (`CDD.md §1.4`). If the agent's identity cannot be confirmed, commit under operator identity but declare this as an override in `self-coherence.md §Debt` with the operator-override cross-reference to §4.

#### 7.3 Override declaration

If the operator commits work on behalf of an agent, this is an implicit override. Declare it per §4:

> Override: operator-identity commit for cycle #N. Reason: α session SIGTERM'd before committing. Committed staged/unstaged work under operator identity at `<SHA>`. Grade implication: per `release/SKILL.md §3.8`, a cycle with operator override has γ < A.

#### 7.4 Prevention (dispatch side)

The recovery procedure is the failure path. The prevention path is a correctly-sized dispatch budget and commit-checkpoint instruction per `CDD.md §1.6c`. If this recovery section is being exercised, the dispatch that spawned the agent did not satisfy §1.6c — record the actual budget and AC count in the PRA telemetry fields (`post-release/SKILL.md §4`: `dispatch_seconds_budget`, `dispatch_seconds_actual`, `commit_count_at_termination`) so the heuristic can be refined.

---

## 9. Embedded Kata

### Kata A — Normal cycle

#### Scenario

γ has dispatched α and β prompts for issue #230. You have two agent sessions available.

#### Task

Execute the operator role through the full cycle.

#### Expected actions

1. Confirm agent-to-role mapping
2. Deliver α prompt to session A, β prompt to session B
3. Wait
4. When β requests merge: execute merge
5. When γ requests tag push: execute tag push
6. When γ declares closure: execute branch cleanup
7. Execute any deferred operator actions from γ's close-out

#### Common failures

- Checking the branch during review and suggesting changes (role leak into β)
- Merging before β approves (gate fired early)
- Rewriting γ's prompts before delivering (role leak into γ)

### Kata B — Override

#### Scenario

α has been implementing for 4 hours. γ reports α is stuck on AC3 due to an undocumented API change. γ's unblock attempt (clarifying the issue) didn't resolve it.

#### Task

Decide whether to override and, if so, execute the protocol.

#### Expected answer

- Override: descope AC3 from this cycle. Reason: undocumented API change requires investigation beyond cycle scope. New scope: AC1–AC2 + AC4–AC6. γ to update the issue and file AC3 as a follow-on.
- Or: provide the missing API documentation as an artifact fact (not implementation direction), which γ can add to the issue. This is unblocking, not override.

#### Reflection

- Did you route through γ, or did you talk directly to α? (Should route through γ unless γ is unavailable.)
- Did you declare the override explicitly, or did you quietly adjust scope?

---

## 10. Wave Coordination

**Multi-cycle coordination primitive for sequences of related issues.**

When δ runs a wave (sequence of related cycles), the artifacts, coordination protocol, and dispatch templates below ensure wave state is durable, auditable, and git-committed — not ephemeral like `/tmp` files.

### 10.1. δ Wave Dispatch Template

**Wave dispatch prompt template for δ-as-agent.** Parallel to α/β templates in `gamma/SKILL.md` §2.5. Takes wave manifest as input, produces γ prompts per issue.

```markdown
# δ Wave Dispatch — {wave-name}

You are δ (operator) for CDD wave coordination.

## Load Order (mandatory)

1. Read `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — canonical δ algorithm
2. Read wave manifest at `.cdd/waves/{wave-id}/manifest.md` — issue order, dependencies, permissions
3. Read wave status at `.cdd/waves/{wave-id}/status.md` — current per-issue state

## Context

- **Wave:** {wave-name}
- **Wave ID:** {wave-id} 
- **Repo:** {repo-path}
- **Issues in scope:** {issue-list}

## Algorithm

1. **Initialize wave directory** — create `.cdd/waves/{wave-id}/` if it doesn't exist
2. **Update wave status** — write current state to `status.md` (queued → in-progress for next issue)  
3. **Select next issue** — per manifest order and dependency constraints
4. **Dispatch γ** — create issue-specific γ prompt using manifest permissions and timeouts
5. **Monitor cycle** — poll `origin/cycle/{N}` and `.cdd/unreleased/{N}/` for completion
6. **Update wave status** — mark cycle completed/failed, update rounds, notes
7. **Repeat or close** — continue with next issue or close wave when complete

## Wave permissions from manifest

- Push to cycle branches: {push-cycles}
- Push merges to main: {push-main}  
- Auto-dispatch α fix rounds: {auto-fix} (max {fix-max})
- Tag/release: {tag-release}
- Branch delete after merge: {branch-delete}

## Standing instructions

- Update `.cdd/waves/{wave-id}/status.md` after each cycle completion
- Emit wave status table on operator request  
- Close wave when all issues completed/deferred per manifest
- Write `.cdd/waves/{wave-id}/wave-closeout.md` at wave end

Signal completion when wave is closed or blocked.
```

### 10.2. Wave Manifest Format

**Canonical wave manifest** — `.cdd/waves/{wave-id}/manifest.md`.

Contains issue list, execution order, dependency constraints, standing permissions, timeout budgets.

```markdown
# Wave: {wave-title}

**Date:** {YYYY-MM-DD}
**Dispatcher:** {δ-identity}  
**Repo:** {repo-path}

## Issues (run in order)

| Order | # | Title | ACs | Type | Dependencies |
|-------|---|-------|-----|------|--------------|
| 1 | {N1} | {title1} | {ac-count1} | {type1} | — |
| 2 | {N2} | {title2} | {ac-count2} | {type2} | #{N1} |
| 3 | {N3} | {title3} | {ac-count3} | {type3} | #{N1}, #{N2} |

## Standing permissions

- Push to cycle branches: {yes|no}
- Push merges to main: {yes|no}
- Auto-dispatch α fix rounds on REQUEST CHANGES: {yes|no} (max {N})
- Tag/release: {yes|no} — write to status.md only
- Branch delete after merge: {yes|no}

## Timeout budgets

- γ: {N}s ({configuration-note})
- α: {N}s  
- β: {N}s

## Dependencies

{per-issue dependency description}
```

### 10.3. Wave Status Format

**Wave status tracking** — `.cdd/waves/{wave-id}/status.md`.

Per-issue status table updated by δ after each cycle closes.

```markdown
# Wave Status: {wave-title}

**Last updated:** {timestamp} by δ

| # | Issue | Status | Rounds | Branch | Tag | Notes |
|---|-------|--------|--------|--------|-----|-------|
| {N1} | {title1} | {status1} | {rounds1} | {branch-state1} | {tag-state1} | {notes1} |
| {N2} | {title2} | {status2} | {rounds2} | {branch-state2} | {tag-state2} | {notes2} |

## Status values

- `⬜ queued` — not started
- `🔄 in-progress` — γ/α/β active
- `✅ completed` — merged to main
- `❌ failed` — cycle abandoned
- `⏸️ blocked` — waiting on dependency
- `⭕ deferred` — explicitly moved out of wave
```

### 10.4. Wave Closure Protocol

**Wave completion conditions and closure procedure.**

Wave is closed when:
1. **All issues completed** or explicitly deferred
2. **Status.md final** — all entries show terminal state (completed/failed/deferred)  
3. **Wave-closeout.md written** with aggregate stats and findings

**Closure algorithm:**
1. **Verify completion** — every manifest issue has terminal status  
2. **Write wave-closeout.md** — wave summary, aggregate metrics, cross-wave findings
3. **Finalize status.md** — add wave closure timestamp, mark final
4. **Signal closure** — wave coordination complete, artifacts durable

### 10.5. δ Wave Reporting Format  

**Wave status table emitted by δ.** Parallel to γ's TLDR (§2.11) but wave-scoped.

| # | Issue | Status | Rounds | Notes |
|---|-------|--------|--------|-------|
| {N} | {title} | {status} | {round-count} | {notes} |

**Status icons:**
- `⬜ queued` — not started  
- `🔄 in-progress` — cycle active
- `✅ completed` — merged, tagged
- `❌ failed` — cycle abandoned
- `⏸️ blocked` — dependency wait
- `⭕ deferred` — moved out of wave

**Emit frequency:**
- After each cycle completion  
- On operator request
- At wave closure

### 10.6. Iteration Artifact Lifecycle

**Three-stage lifecycle for wave iteration findings.**

Wave iteration files (like `.cdd/iterations/wave-2026-05-12.md`) track cross-cycle findings and their dispositions through wave execution and into release.

**Stage 1: Per-issue close**
- Wave iteration file's disposition column updated with issue numbers or commit SHAs
- Each finding gets: `filed #{N}`, `patched {SHA}`, `deferred`, or `no-action`
- Updated as each wave issue closes

**Stage 2: Wave close**  
- Wave iteration file moves to `.cdd/waves/{wave-id}/iteration.md`
- Dispositions finalized, no further updates
- Wave-closeout.md cross-references iteration.md
- Iteration.md becomes part of durable wave record

**Stage 3: Release boundary**
- Findings that produced MCAs (Major Change Approvals) get cross-referenced in release PRA
- PRA links back to specific iteration.md findings that drove release changes
- Creates audit trail from wave findings → MCAs → release notes

**File movement example:**
```
.cdd/iterations/wave-2026-05-12.md  
  ↓ (wave closes)
.cdd/waves/hardening-2026-05-12/iteration.md
  ↓ (release time) 
.cdd/releases/3.60.0/post-release-assessment.md (references iteration findings)
```