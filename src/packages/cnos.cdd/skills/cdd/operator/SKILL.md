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
  - cycle state (issue, PR, CI)
outputs:
  - routed prompts to agent sessions
  - external gate decisions (merge, tag, issue filing)
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

1. **Route** — deliver γ's dispatch prompts to agent sessions.
2. **Wait** — monitor for gate requests. Do not poll the triad's internal work.
3. **Gate** — execute external actions when requested by a role.
4. **Override** — reassign roles or redirect scope only with an explicit declaration.

---

## 1. Route

### 1.1. Receive dispatch from γ

γ produces α and β prompts (CDD §1.4). The operator delivers each prompt to the correct agent session.

- ❌ Rewrite the prompt to add constraints or context γ didn't include
- ✅ Deliver the prompt verbatim to the target agent session

### 1.2. Confirm session assignment

Before dispatch, confirm which agent session runs α and which runs β. If only two agents are available, the operator serves as γ (CDD §1.4 minimum configuration).

- ❌ Dispatch both prompts to the same session without declaring the configuration
- ✅ Name the agent-to-role mapping before delivering prompts

---

## 2. Wait

### 2.1. Do not poll internal work

Once dispatched, the triad runs. The operator does not need to monitor PR diffs, review comments, or CI runs — γ owns that.

The operator's wake-up signals are:

- **γ requests an external gate** (merge, tag push, issue filing, auth action)
- **γ requests an unblock decision** (design ambiguity, scope question, env constraint)
- **γ declares cycle closure** and names deferred outputs that need operator action
- **An agent session dies or stalls** (no activity for an unexpected duration)

Between these signals, the operator's correct action is nothing.

- ❌ Check the PR every 30 minutes and suggest changes
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

Run under `Monitor` or equivalent. 5-minute interval is sufficient for δ — γ owns the tight loop. Supplement with PR polling once a PR exists:

```bash
prev=""; while true; do
  cur="$(gh pr list --repo owner/repo --search 'is:open' --json number,title --jq '.[].number' | sort)"
  comm -13 <(echo "$prev") <(echo "$cur") | sed 's/^/new-pr: /'
  prev="$cur"; sleep 300
done
```

---

## 3. Gate

### 3.1. External actions the operator holds

These actions require platform permissions agents may lack:

| Action | Trigger | Who requests |
|--------|---------|-------------|
| PR merge | β approves, merge ready | β or γ |
| Tag push | Release tagged | β or γ |
| Branch delete | Cycle closed, merged branches | γ |
| Issue filing on external repos | Cross-project dependency | γ |
| Force push | Rebase required with env constraints | α via γ |
| Auth refresh | Token/permission expiry | any role via γ |

### 3.2. Execute on request, not on observation

Gate actions fire when a role requests them, not when δ notices they're needed. Observing that a tag isn't pushed or a branch exists is not a gate trigger — γ's explicit request is.

- ❌ Heartbeat shows tag not pushed → δ pushes it (role leak: δ decided the gate, not γ)
- ❌ "β asked me to merge but I think we should wait for one more review"
- ✅ γ requests tag push → δ pushes it and confirms
- ✅ If you disagree with a gate request, declare an override (§4)

### 3.3. Report completion

After executing a gate action, confirm to the requesting role that the action completed.

- ❌ Execute silently and assume the triad will notice
- ✅ "Tag pushed: `git push origin 3.59.0` — confirmed on remote"

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
| Pre-dispatch | Receive γ prompts, confirm agent mapping | γ dispatch |
| Dispatch | Deliver prompts to agent sessions | — |
| Implementation | Nothing | α PR or γ unblock request |
| Review | Nothing | β verdict or γ unblock request |
| Release | Gate actions if requested (merge, tag) | β or γ gate request |
| Closure | Gate actions if requested (branch delete, issue close) | γ closure declaration |
| Post-closure | Execute deferred operator actions from γ close-out | γ deferred-output list |

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

- Checking the PR during review and suggesting changes (role leak into β)
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
