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

## Algorithm

1. **Dispatch γ** — run γ via `claude -p`; γ reads the issue, creates the cycle branch, and produces α/β prompts.
2. **Dispatch α** — run α via `claude -p` with the prompt γ produced; α implements on the cycle branch and exits after signaling review-readiness.
3. **Dispatch β** — run β via `claude -p` with the prompt γ produced; β reviews, merges, writes β close-out, and exits.
4. **Re-dispatch α for fix rounds** (when β returns RC) — run α via `claude -p` with the fix-round re-dispatch prompt (CDD.md §1.6a); α fixes findings, appends fix-round to self-coherence.md, exits.
5. **Re-dispatch α for close-out** (when γ requests after β merge) — run α via `claude -p` with the close-out re-dispatch prompt (CDD.md §1.6a); α writes alpha-closeout.md, commits to main, exits. **This step is mandatory when γ requests it.** γ cannot complete the closure gate without alpha-closeout.md.
6. **Gate** — execute external actions: push main, tag, release, branch cleanup. **Do not tag/release before `gamma-closeout.md` exists on main.**
7. **Override** — reassign roles or redirect scope only with an explicit declaration.

δ runs one role at a time. This keeps memory pressure low (single `claude -p` process), gives δ direct visibility into each session, and isolates failures — if α dies, δ retries α without losing γ or β state.

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

- ❌ Check the branch every 30 minutes and suggest changes
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
| Push β-approved merge to main | β runs `git merge` — δ only pushes when β cannot execute the push directly (env/auth constraint). This is execution of β's integration authority, not δ approval. | β or γ |
| Release-boundary preflight | After β merge + close-outs + γ PRA, δ verifies merge commit, release artifacts, tag/deploy preconditions, and platform readiness. Proceed / request changes / override. See CDD §1.4 Phase 5a. | γ |
| Tag push + release | After δ preflight confirms and γ closes the cycle | γ |
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
3. Run `scripts/release.sh` — this stamps all manifests, verifies consistency, commits, tags, and pushes in one command

**Manual tagging is not allowed.** Do not run `git tag` directly. The release script is the only way to tag. It prevents the class of failures where VERSION, cn.json, and package manifests disagree (see DISPATCH-FAILURE-EVIDENCE.md, cycle #84 failure 3).

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

## 5. What the operator does NOT do

These are role boundaries. Crossing them without an override declaration breaks the coherence record.

- **Do not implement.** That's α. If you want to fix something, declare an override and take the α role explicitly.
- **Do not review.** That's β. If you disagree with β's verdict, declare an override.
- **Do not triage findings.** That's γ. If you want a finding handled differently, tell γ.
- **Do not rewrite prompts.** γ owns prompt quality. If the prompt is weak, tell γ to fix the issue.
- **Do not merge without β's approval.** The merge gate fires when β says it's ready, not before.

---

## 6. Cycle lifecycle from the operator's view

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

## 7. Embedded Kata

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
