---
name: troubleshoot
description: Diagnose and resolve active environmental, runtime, process, or tool failures while work is blocked.
artifact_class: skill
kata_surface: embedded
governing_question: How do you isolate an active failure through the cheapest discriminating hypothesis before speculating about application or model behavior?
triggers:
  - troubleshoot
  - live failure
  - process died
  - silent exit
  - tool error
  - OOM
  - diagnosis
  - debugging
scope: task-local
inputs:
  - blocked task context
  - symptom description
  - available evidence (logs, stream output, exit codes, process state)
outputs:
  - resolved symptom or blocked/escalated troubleshooting record
  - triage record (hypotheses tested, results, change made)
  - RCA handoff trigger (if recurrence risk or systemic root cause)
---

# Troubleshoot

This skill teaches how to isolate an active failure through the cheapest discriminating hypothesis before speculating about application or model behavior.

## Core Principle

**Coherent live troubleshooting preserves evidence first, tests hypotheses cheapest-to-most-expensive in a fixed triage order, makes one controlled change, and verifies against the original symptom.**

The failure mode is **premature hypothesis**: jumping to application-layer or model-layer speculation before ruling out process state, external kill, tool error, resource pressure, and lifecycle failures. Every dispatch-test failure in 2026-04-30 shared this pattern.

Live troubleshooting ends when the original symptom is resolved. Root cause analysis (`eng/rca`) begins after.

## Algorithm

1. **Define** — name the symptom, the evidence at hand, and the failure class under investigation.
2. **Unfold** — triage cheapest-to-most-expensive; test one hypothesis at a time; make one change.
3. **Rules** — preserve evidence, state hypotheses before tests, verify the original symptom.

---

## 1. Define

### 1.0. Classify the artifact

This is a **skill**, not a runbook:

- It teaches judgment — what to check first, how to form a hypothesis, what counts as verification.
- Platform-specific diagnostic commands belong in runbooks that reference this skill.
- Reference pages listing diagnostic commands are lookup material, not live diagnosis.

A runbook is appropriate when the steps are fully ordered and platform-specific. This skill is appropriate when the judgment of *which step first* is the problem.

### 1.1. Identify the parts

Live troubleshooting has six parts:

- **Symptom** — the observable failure that blocked work
- **Evidence** — what is already available without touching the system (logs, stream output, exit codes, process state)
- **Hypothesis** — a testable claim about the failure class
- **Discriminating test** — the cheapest test that rules the hypothesis in or out
- **Change** — one controlled modification to address the confirmed hypothesis
- **Verification** — proof that the original symptom is gone after the change

- ❌ "Something broke, let me try things"
- ✅ "Symptom: silent exit. Evidence: no stderr. Hypothesis: OOM kill. Test: dmesg."

### 1.2. Articulate how they fit

The six parts work in sequence: preserve evidence → describe symptom → form hypothesis → test cheapest discriminating → make one change → verify.

If the symptom returns after a change, the hypothesis was wrong or incomplete. Record the negative result and advance to the next hypothesis in triage order.

### 1.3. Name the failure mode

Troubleshooting fails through **premature hypothesis**:

- Jumping from symptom to application-layer explanation without checking process, kernel, tool, resource, or lifecycle evidence
- Testing the most expensive hypothesis first (model behavior, token budget) when cheaper tests (dmesg, stderr, process state) would have discriminated immediately
- Making multiple simultaneous changes and losing the signal
- Restarting the system before preserving evidence

- ❌ "Process disappeared → token limit exceeded"
- ✅ "Process disappeared → check process state → check dmesg → OOM kill confirmed"

---

## 2. Unfold

### 2.1. Preserve evidence first

Before any action, record: timestamp, symptom text, all available output, environment state, and process identifiers. Evidence that exists now may not exist after touching the system.

- ❌ Restart the process immediately to see if it reproduces
- ✅ Read logs, capture stderr, check process state before restarting anything

### 2.2. Describe the symptom precisely

State the symptom as a complete problem description before forming a hypothesis. A complete description eliminates an entire class of wrong hypotheses.

Required fields (drawn from IBM troubleshooting problem-description practice):

- **What** is failing (process, command, tool call, stream)
- **Where** does the failure appear (which surface, which output stream)
- **When** did it start (after what action, at what step in the sequence)
- **Conditions** — what changed before the failure appeared
- **Reproducibility** — fails every time, or intermittently
- **Error messages** — exact text (not paraphrased)
- **Effect** — what work is blocked
- **Environment** — host, OS, memory, parent process, network state

- ❌ "The agent died"
- ✅ "Beta process exited after 9 stream events. No stderr. No result event. tool_result is_error: true, message: 'GraphQL field deprecated'. Reproducible every run. Host: 2GB VPS. Parent: foreground shell."

### 2.3. Triage order

Test in this order. Each class is cheaper and faster to eliminate than the next. Do not skip a cheaper class without accounting for it (see §3.8).

1. **Process state** — is the process running? Check exit code, PID, process list.
2. **External kill / kernel log** — was the process killed externally? Check dmesg, kernel log, OOM killer output.
3. **Command / tool output** — what did the failed command or tool return? Read stderr, stdout, tool_result, or event stream for the exact error message.
4. **Resource pressure** — memory headroom, swap availability, disk space, open file descriptors, CPU saturation.
5. **Lifecycle / parent process** — was the parent process still alive? Was the process backgrounded? Did the shell session end?
6. **Application / model behavior** — only after 1–5 are ruled out: application logic, model configuration, token budget, API behavior.

*Triage order derived from: Google SRE effective troubleshooting (hypothesis-driven process), CompTIA sequence (identify → theory → test → plan → implement/escalate → verify → document), and three root-cause classes from the 2026-04-30 dispatch test.*

- ❌ Check model token budget before checking process state
- ✅ Work down the triage order; stop at the first class that produces discriminating evidence

### 2.4. Hypothesis discipline

Form one hypothesis at a time, test it, record the result, and advance.

Required steps per hypothesis:

1. **State the hypothesis** before running any test: "I believe the process was OOM-killed."
2. **State the oracle**: "If true, dmesg will show an OOM event for this PID or process name."
3. **Run the cheapest discriminating test** — the test that most cheaply confirms or rules out this hypothesis.
4. **Record the result** — positive or negative. A negative result is data, not failure.
5. **Advance** to the next hypothesis only after the current one is confirmed or ruled out.

*Derived from: Red Hat one-change discipline (one variable at a time, test original symptoms), Google SRE hypothesis-driven troubleshooting.*

- ❌ "Let me add swap, switch to foreground, and reduce workers simultaneously"
- ✅ "Hypothesis: OOM kill. Oracle: dmesg shows OOM for this PID. Test. Confirmed. One change: add swap. Verify."

### 2.5. One controlled change

When a hypothesis is confirmed, make exactly one change. Then verify.

If the symptom persists, the hypothesis was wrong or incomplete. Record the result and advance to the next hypothesis in triage order.

- ❌ Two changes at once to "fix faster"
- ✅ One change, verify original symptom gone, then consider whether a second change is needed

### 2.6. Verify against the original symptom

After the change, re-run the original failing scenario. Confirmation means the original symptom no longer occurs under the same conditions — not "it seems better."

*Derived from: Red Hat problem-solving — always test original symptoms after any change.*

- ❌ "Added swap, looks good"
- ✅ "Re-ran β review cycle under same conditions; process completed without OOM; no silent exit"

### 2.7. RCA handoff

Hand off to `eng/rca` when the live symptom is resolved and any of these apply:

- Recurrence risk remains (the fix is a workaround, not structural prevention)
- Multi-step incident affected dispatch, release, or runtime
- Root cause is systemic, not one-off
- Prevention action is needed
- Evidence should be preserved for future agents

Do not start an RCA while still diagnosing. `eng/rca` requires a stable system.

- ❌ Starting Five Whys while the process is still failing
- ✅ Resolve the live symptom; then trigger RCA if recurrence risk or systemic cause applies

### 2.8. Escalate or close blocked

If no hypothesis can be safely tested, no confirmed fix is available, or the required action exceeds current authority, stop live troubleshooting and return a blocked troubleshooting record.

The record must include:
- Original symptom
- Evidence preserved
- Hypotheses tested
- Negative results
- Current best hypothesis, if any
- Reason live troubleshooting cannot continue
- Required next authority, resource, tool, or operator action

Do not keep poking after the evidence or authority boundary is exhausted.

- ❌ Trying increasingly risky fixes because "something has to work"
- ✅ "Cannot test further without root SSH access. Best hypothesis: filesystem full on host. Escalating to operator with evidence."

---

## 3. Rules

### 3.1. Preserve evidence before touching anything

A restart or a change can destroy the evidence you need. Capture it first.

### 3.2. State the hypothesis before the test

If you cannot state the hypothesis before running the test, you are not testing — you are poking.

- ❌ Run dmesg and see what it says
- ✅ Hypothesis: OOM kill. Oracle: dmesg shows OOM event for this PID. Test: `dmesg | grep -i "killed process"`.

### 3.3. Test the cheapest discriminating hypothesis first

The triage order in §2.3 is not optional. Process state and kernel log are cheaper than API debug logs. API debug logs are cheaper than model behavior speculation. But evidence already in hand can satisfy a class without a dedicated test (see §3.8).

- ❌ "Let me check if the token budget was exceeded" (before checking process state)
- ✅ "Process state first. No process running, no exit code. Step 2: kernel log."

### 3.4. Make one change at a time

Multiple simultaneous changes make it impossible to attribute the fix.

- ❌ Add swap, switch to foreground, and reduce workers in the same action
- ✅ Add swap. Verify. If symptom recurs, then investigate the next confirmed hypothesis.

### 3.5. Record negative results

A hypothesis tested and ruled out is evidence. Record it. Future agents should not re-test a ruled-out hypothesis.

### 3.6. Verify against the original symptom

Proxy verification is not verification. Run the same scenario that originally failed.

- ❌ "Process ran for 30 seconds without crashing — must be fixed"
- ✅ "Re-ran the full beta review scenario; process reached a result event; symptom does not recur"

### 3.7. Do not start RCA during live diagnosis

RCA and live diagnosis are distinct processes. Mixing them produces noise in both.

### 3.8. Do not skip triage steps without accounting for them

Do not skip a cheaper class without accounting for it. For each earlier class, either:
- Test it
- Mark it already satisfied by preserved evidence
- Mark it not applicable with reason
- Record why testing it would destroy evidence or exceed authority

This preserves discipline while avoiding ritual checks when discriminating evidence is already in hand.

- ❌ "Tool error is obvious — skip process/kernel/resource checks entirely"
- ✅ "Process: alive (ps confirms). Kernel: no kill (dmesg clean). Tool: error message says 'deprecated GraphQL field' — this is the discriminating evidence. Remaining classes: not tested, not needed."

---

## 4. Worked Examples

The following examples are from the 2026-04-30 identity-rotation dispatch test: five β dispatch failures, three root-cause classes, each diagnosed ad-hoc with a wrong first hypothesis. Each example shows: symptom → wrong first hypothesis → better triage path → root cause → corrective action → verification.

### Example 1 — OOM kill

**Symptom:** β process disappeared mid-review. No output, no error, no result event.

**Wrong first hypothesis:** Token limit exceeded.

**Better triage path:**
1. Process state: no process running, no exit code available
2. Kernel log: `dmesg` → `Out of memory: Kill process <PID> (node) score <N> or sacrifice child`

**Root cause:** 2 GB VPS, no swap, two Node processes (OpenClaw gateway + Claude CLI) competing for memory. Kernel OOM-killed the largest.

**Corrective action:** Add 2 GB swap.

**Verification:** Re-ran β review under same conditions. Process completed without OOM kill. Original symptom gone.

**What this skill forces:** "Process disappeared silently → check process state (step 1) → check kernel log for OOM before assuming application-level failure (step 2)."

---

### Example 2 — `gh` GraphQL error

**Symptom:** Session died after 9 stream events. tool_result showed `is_error: true`.

**Wrong first hypothesis:** GitHub API rate limit.

**Better triage path:**
1. Process state: process exited cleanly with exit code 0
2. Kernel log: no kill event
3. Tool output: tool_result `is_error: true`, message: `"GitHub Projects Classic deprecated GraphQL field"`

**Root cause:** `gh issue view` without `--json` hits a deprecated GraphQL field. Tool error, not rate limit.

**Corrective action:** Use `gh issue view N --json title,body,state,comments` in all dispatch prompts.

**Verification:** Subsequent `gh issue view --json` calls return data without error.

**What this skill forces:** "Tool call returned an error → read the error message before assuming surrounding context (step 3: tool output). Rate limit is a step-6 application hypothesis; step-3 tool output was the discriminating evidence."

---

### Example 3 — Background process kill

**Symptom:** Process exited silently. Partial stream-json output, no result event, no error. Same pattern as OOM but dmesg was clean.

**Wrong first hypothesis:** Rate limit / token budget.

**Better triage path:**
1. Process state: no process running
2. Kernel log: no OOM event
3. Tool output: no tool_result error
4. Resources: memory normal
5. Lifecycle / parent process: parent exec session ended → background process (launched with `&`) was cleaned up by shell when parent exited

**Root cause:** Shell backgrounding (`&`) attaches the process to the parent's session. When the parent exec session ends, the shell cleans up the process tree.

**Corrective action:** Run `claude -p` as a foreground direct child process only. Never background the session process with `&`.

**Verification:** Subsequent foreground runs complete without silent exit. Original symptom gone.

**What this skill forces:** "Silent exit, no error, no OOM → work down the triage order through lifecycle/parent (step 5) before reaching application/model behavior (step 6). Each step-1–4 ruling out is recorded."

---

## 5. Kata

### Scenario

An agent running a long beta review goes silent. The stream stops mid-output. No error appears. No result event is emitted. The operator sees partial output and no final status. Host: 2 GB RAM, no swap.

### Symptoms

- Stream output stops after N events
- No result event
- No stderr output
- No error visible in the last tool_result

### Available evidence

- Partial stream-json log
- Host: 1.8 GB used / 2.0 GB total RAM, no swap configured
- `dmesg` access available
- Process list available

### Wrong tempting hypothesis

"The agent hit its token limit or the model stopped generating."

### Expected triage path

1. **Process state:** confirm the process is no longer running
2. **Kernel log:** check dmesg for OOM kill event matching the process PID or name
3. **Tool output:** read the last few stream events for error signals in tool_result
4. **Resources:** confirm memory headroom and absence of swap
5. **Lifecycle:** confirm whether the process was foreground or backgrounded
6. **Application/model:** (only if 1–5 yield nothing — skip here given the evidence)

### Expected root cause

OOM kill: host memory exhausted, no swap, kernel killed the largest process.

### Corrective action

Add swap or reduce concurrent memory pressure. Verify the process completes the same scenario without OOM.

### Verification

Re-run the same beta review scenario. Process reaches a result event. Original symptom (silent exit mid-stream) does not recur.

### Common failures

- Jumping to step 6 (application/model) before step 2 (kernel log) — the most repeated pattern from the dispatch test
- Restarting the process before capturing dmesg — evidence is gone
- Making multiple changes simultaneously (add swap + switch foreground + reduce workers) — cannot isolate which change fixed it
- Verifying "process ran longer" instead of verifying "process completed the original scenario"

### Reflection

- Did the triage order force you to check process state and kernel log before reaching application-layer explanations?
- Did you state your hypothesis before running each test?
- Did you make exactly one change and verify against the original symptom, not a proxy?
- Did you record negative results as evidence for future agents?

---

## 6. Sources

This skill adapts:

- **Google SRE, _Effective Troubleshooting_** — hypothesis-driven troubleshooting, telemetry/log examination, test-and-treat loop, and common pitfalls.
- **IBM troubleshooting techniques** — complete problem description, conditions, reproducibility, and environment questions.
- **Red Hat problem-solving** — symptom vs problem, one change at a time, and retesting the original symptom.
- **CompTIA troubleshooting methodology** — identify, form theory, test, plan, implement/escalate, verify, document.
