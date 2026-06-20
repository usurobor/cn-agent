---
name: dispatch-protocol
description: Label-gated issue queue protocol for implementation cells. Defines the canonical GitHub label set, claim operation, concurrency discipline, lifecycle transitions, wave semantics, and drift handling by which a human authorizes a wake worker to claim and execute exactly one implementation cell per run.
artifact_class: skill
kata_surface: embedded
governing_question: How does a human authorize exactly one implementation cell per wake run, with no double execution, no unauthorized claiming, and explicit operator control over the dispatch gate?
visibility: public
parent: agent
triggers:
  - dispatch
  - dispatch protocol
  - label gate
  - dispatch:cell
  - status:todo
  - claim issue
  - implementation cell
scope: task-local
inputs:
  - GitHub repository with labeled issues
  - wake worker with repository read/write access
  - operator authorization (status:todo label applied)
outputs:
  - at most one executable issue claimed per run
  - claim comment written with run id + head sha
  - status labels transitioned (todo → in-progress)
  - cell runtime launched on claimed issue
requires:
  - GitHub repository with label set defined per AC1
  - operator-applied status:todo label on target issue
  - repository-level workflow concurrency group
calls:
  - (dispatches into CDD cell runtime: γ/α/β inside the claimed issue's cycle)
---

<!--
section-manifest:
  planned: [frontmatter, core-principle, algorithm, define, unfold-label-set, unfold-claim-operation, unfold-concurrency, unfold-transitions, unfold-wave-semantics, unfold-drift-handling, rules, verify, failure-modes, references]
  completed: [frontmatter, core-principle, algorithm, define, unfold-label-set, unfold-claim-operation, unfold-concurrency, unfold-transitions, unfold-wave-semantics, unfold-drift-handling, rules, verify, failure-modes, references]
-->

# Dispatch Protocol

## Core Principle

**An issue is executable only when ALL of: open, carries `dispatch:cell`, carries `status:todo`.** The wake claims at most one executable issue per run. Lifecycle (`status:*`) ≠ issue type (`dispatch:cell`).

The two-label composition is intentional. `status:todo` alone is not enough because master issues, RFC review issues, meta-planning issues, and implementation cells can all be "todo" in a human board sense. Only **implementation cells** should be executable — `dispatch:cell` is the type marker. `status:todo` is the authorization gate. Both required; each necessary; neither sufficient alone.

The failure mode this skill prevents is **unauthorized claiming** — a wake worker grabbing a master issue, a review RFC, or a planning issue and launching a cell runtime on it. The `dispatch:cell` label is the structural lock: humans apply it deliberately to issues that are scope-bounded, independently shippable, and ready for autonomous execution. Its absence on master/RFC/planning issues is the signal the worker respects.

## Algorithm

1. **Query** — find open issues with `dispatch:cell` + `status:todo`. On a sweep run, find all; take the lowest issue number. On an `issues: labeled` event run, check only the event issue.
2. **Re-read** — do not trust the query result or event payload alone. Re-read the issue's current labels before claiming.
3. **Verify** — confirm open, `dispatch:cell` present, `status:todo` present, `status:in-progress` absent, `status:blocked` absent, `status:review` absent.
4. **Claim** — atomically (within the GitHub label API's best-effort): remove `status:todo`, add `status:in-progress`, write claim comment.
5. **Re-verify** — re-read labels after claim. If `status:in-progress` is not present, abort and report.
6. **Launch** — pass claimed issue to cell runtime (γ → α → β → γ-receipt inside the issue's cycle scope).

---

## 1. Define

### 1.1. Label set (canonical, v0)

Two label namespaces: `status:*` (lifecycle) and type markers (`dispatch:cell`, `kind/*`). No per-org overrides in v0.

**Lifecycle labels** — exactly one per open issue:

| Label | Meaning |
|-------|---------|
| `status:backlog` | Identified but not yet ready for implementation |
| `status:ready` | Specified and reviewed; not yet authorized for dispatch |
| `status:todo` | Human-authorized for dispatch; executable when combined with `dispatch:cell` |
| `status:in-progress` | Claimed by a wake; cell runtime running |
| `status:review` | Cell runtime complete; awaiting human review |
| `status:changes` | Human requested changes; awaiting repair |
| `status:blocked` | Blocked on external dependency or human decision |

Closed issue = done. No `status:done` label in v0.

**Type markers** — zero or one per issue:

| Label | Meaning |
|-------|---------|
| `dispatch:cell` | This issue is an implementation cell — independently shippable, scope-bounded, executable by wake |
| `kind/wave` | Optional; marks a master issue tracking a wave of dependent sub-issues |
| `kind/rfc` | Optional; marks an RFC or review issue (MUST NOT carry `dispatch:cell`) |

**Invariants:**
- Exactly one `status:*` label per open issue.
- `dispatch:cell` MUST NOT appear on master issues, RFC issues, or planning issues.
- Only `status:todo` + `dispatch:cell` can be claimed.

### 1.2. Human dispatch gate

The distinction between `status:ready` and `status:todo` encodes human intent:

- `status:ready` — spec is complete, review is done, but the human has not yet authorized autonomous execution. A planner or operator reviewing the spec applies `status:ready`.
- `status:todo` — the human explicitly authorizes dispatch. No wake may claim an issue without this label.

Who may apply `status:todo`:
- **Directly**: any operator with write access to the repo.
- **Via planner (interactive Sigma)**: only when quoting an explicit human approval comment (e.g. "Go — dispatch #449"). Planner MUST NOT apply `status:todo` on its own judgment.

The `status:ready` → `status:todo` transition is the operator's authorization event. It is not automated. It is the gate this skill protects.

- ❌ Wake applies `status:todo` to an issue it thinks is ready
- ✅ Human applies `status:todo`; wake claims on next run

### 1.3. Failure modes

Five named failure modes specific to dispatch. Each has a structural fix in §2.

- **D1 — Type confusion.** Wake claims a master issue or RFC because it carries `status:todo` without `dispatch:cell`. Fix: §3.1 — both labels required; `dispatch:cell` MUST NOT appear on non-cell issues.
- **D2 — Double execution.** Two wakes (event-driven + scheduled sweep) race; both see `status:todo` and both claim. Fix: §2.3 — repository concurrency group + claim-time re-read + claim comment with run id.
- **D3 — Stale event payload.** Wake receives `issues: labeled` event; trusts the event's label snapshot; issue labels changed since the event fired. Fix: §2.2 — always re-read current labels before claiming; never trust event payload alone.
- **D4 — Ghost claim.** Wake writes claim comment but label transition fails (API error). Cell launches without `status:in-progress` set; next sweep re-claims the same issue. Fix: §2.2 — re-verify `status:in-progress` after claim; abort if absent.
- **D5 — Unauthorized transition.** Planner applies `status:todo` without an explicit human approval quote. Operator intent is bypassed. Fix: §1.2 human gate rule + §3.4.

---

## 2. Unfold

### 2.1. Label set administration

Labels are created at the repository level before any dispatch can occur. Wake template (cnos#450) includes a setup step that creates the canonical label set if missing. Label names are lowercase-with-colon namespace (`status:todo`, `dispatch:cell`) — no spaces.

Creating missing labels (in setup or first run):
```bash
# status labels
gh label create "status:backlog" --color "#EEEEEE" --repo <owner>/<repo>
gh label create "status:ready"   --color "#BFD4F2" --repo <owner>/<repo>
gh label create "status:todo"    --color "#0075CA" --repo <owner>/<repo>
gh label create "status:in-progress" --color "#E4E669" --repo <owner>/<repo>
gh label create "status:review"  --color "#D4C5F9" --repo <owner>/<repo>
gh label create "status:changes" --color "#E11D48" --repo <owner>/<repo>
gh label create "status:blocked" --color "#B60205" --repo <owner>/<repo>
# type markers
gh label create "dispatch:cell"  --color "#0E8A16" --repo <owner>/<repo>
gh label create "kind/wave"      --color "#C5DEF5" --repo <owner>/<repo>
gh label create "kind/rfc"       --color "#F9D0C4" --repo <owner>/<repo>
```

Colors are recommendations; operators may change colors without affecting protocol behavior.

### 2.2. Claim operation

The claim sequence. Every step is required. No step may be skipped.

```
1. Query for executable issues:
   - On scheduled sweep: gh issue list --repo <repo> --label "dispatch:cell" --label "status:todo" --state open --json number,labels
     → sort by issue number ascending; take lowest
   - On issues:labeled event: gh issue view <event-issue-number> --repo <repo> --json number,labels,state

2. Re-read labels:
   gh issue view <issue-number> --repo <repo> --json labels,state
   → Do not proceed with event payload labels alone

3. Verify (all must pass; abort on any failure):
   - issue.state == "open"
   - "dispatch:cell" in labels
   - "status:todo" in labels
   - "status:in-progress" NOT in labels
   - "status:blocked" NOT in labels
   - "status:review" NOT in labels
   → On failure: log "not claimable: <reason>"; exit 0 (do not error; scheduled sweep will retry)

4. Claim:
   gh issue edit <issue-number> --repo <repo> \
     --remove-label "status:todo" \
     --add-label "status:in-progress"

   gh issue comment <issue-number> --repo <repo> --body "$(cat <<EOF
   claimed_by: <home-hub>/<run-id>
   claimed_at: <utc-timestamp>
   head: <current-main-sha>
   EOF
   )"

5. Re-verify:
   gh issue view <issue-number> --repo <repo> --json labels
   → Confirm "status:in-progress" is present
   → If absent: abort; log "claim failed: status:in-progress not set after edit"; exit 1

6. Launch cell runtime (only after step 5 confirms):
   → pass issue number to cell runtime (γ → α → β → γ-receipt)
```

- ❌ Wake trusts the query result; skips re-read
- ✅ Wake re-reads labels immediately before every claim

### 2.3. Concurrency and anti-double-execution

GitHub label API is not a true compare-and-swap primitive. Anti-double-execution discipline is a composition of three layers:

**Layer 1 — Repository concurrency group (GH Actions serialization)**

The wake workflow declares a single repository-level concurrency group:
```yaml
concurrency:
  group: ${{ github.repository }}-dispatch
  cancel-in-progress: false
```

At most one dispatch job runs at a time per repository. A second trigger queues behind the first, not racing.

**Layer 2 — Claim-time re-read + label transition**

Step 3 (verify) + step 4 (claim) + step 5 (re-verify) form a best-effort CAS. If two runs reach step 4 concurrently despite layer 1 (e.g. different workflow files), only one can successfully transition `status:todo` → `status:in-progress` before the other re-reads in step 5 and finds `status:in-progress` already set. The second run aborts.

**Layer 3 — Claim comment with run id**

The claim comment records `claimed_by: <hub>/<run-id>` and `head: <sha>`. If a double-claim occurs and both transitions appear to succeed (race between two separate GitHub Actions runs), the claim comments provide an audit trail. Operator inspects; manual deduplication; the cell that ran furthest is the canonical one.

**Sweep is backstop, not primary**

Issue-labeled events fire on each `status:todo` label application. Scheduled sweep (e.g. every 15 minutes) backstops missed or canceled events. Sweep takes the lowest-numbered executable issue. Event-driven claim takes only the event issue. The two paths do not race in normal operation (one is scoped to one issue; the other to all); they race only if an event-driven run lags past the next sweep interval.

### 2.4. Lifecycle transitions

Complete transition map. All transitions require explicit action (label change + comment); no implicit transitions.

```
[created as status:backlog]
    │
    │ planner / operator: spec complete, reviewed
    ▼
status:ready
    │
    │ operator: explicit "go" authorization
    ▼
status:todo  ← ── ── [blocked or changes resolved; re-authorized]
    │
    │ wake: claim operation (§2.2)
    ▼
status:in-progress
    │
    ├─────────────────────────────────────┐
    │ cell completes                      │ cell blocked
    ▼                                     ▼
status:review                         status:blocked
    │                                     │
    ├── human accepts                     │ human/planner resolves
    │       ▼                             ▼
    │   close issue                  status:ready | status:todo
    │
    └── human rejects
            ▼
        status:changes
            │
            │ planner/human repairs
            ▼
        status:ready | status:todo
```

**Transition comments** (required on state changes):

| Transition | Required comment |
|------------|------------------|
| `todo → in-progress` | claim comment (§2.2 step 4) |
| `in-progress → review` | PR/receipt link |
| `review → changes` | specific change requests |
| `changes → ready/todo` | repair summary |
| `in-progress → blocked` | names the blocking dependency and required human decision |
| `close` | "Accepted. Closes #N." |

### 2.5. Wave semantics

A wave is a coordinated group of implementation cells with a shared tracker (master issue).

Master issue rules:
- Carries `kind/wave` label.
- MUST NOT carry `dispatch:cell`.
- Lists all sub-issues with expected ordering / dependency notes.
- Does not carry `status:*` lifecycle labels (master is not dispatchable).

Sub-issue rules:
- Each sub carries `dispatch:cell`.
- Each sub is independently dispatchable only when labeled `status:todo`.
- Default: sequential dispatch — sub N carries `status:todo` only after sub N-1 closes.
- Exception: if sub N and sub N+1 are explicitly independent AND merge-safe (`dispatch:cell` issues that can land in either order without conflict), both may carry `status:todo` simultaneously for parallel dispatch. Parallel is the exception; justification required.

Lowest-issue-number sweep ordering: when multiple subs are simultaneously `status:todo` (parallel), sweep claims the lowest issue number first. Priority labels are deferred to v1.

### 2.6. Drift handling

A drift state is an issue with an unexpected or inconsistent label combination. Do not claim drifted issues.

| Drift pattern | Degraded reason | Repair |
|---------------|-----------------|--------|
| Multiple `status:*` labels | `dispatch_label_drift` | Remove extra status labels; leave one |
| No `status:*` label | `dispatch_label_drift` | Add correct status label |
| `dispatch:cell` on a master/RFC issue | `dispatch_type_mismatch` | Remove `dispatch:cell` from non-cell issue |
| `status:in-progress` without a claim comment | `dispatch_ghost_claim` | Clear `status:in-progress`; re-apply `status:todo` if still authorized |

On drift detection:
```
outcome: degraded
degraded_reason: dispatch_label_drift | dispatch_type_mismatch | dispatch_ghost_claim
operator_action: repair instruction (specific labels to add/remove)
```

Post a comment naming the drift and the repair instruction. Do NOT claim the drifted issue.

---

## 3. Rules

### 3.1. Both labels required for claim

A wake MUST NOT claim an issue missing either `dispatch:cell` or `status:todo`. No exceptions. Master issues, RFC issues, and planning issues MUST NOT carry `dispatch:cell` — this is the structural lock against type confusion.

- ❌ "This issue has status:todo; I'll claim it for the wave tracker"
- ✅ "No dispatch:cell label; skip"

### 3.2. At most one claim per run

A single wake run claims at most one issue. Scheduled sweep queries all `dispatch:cell + status:todo` issues and takes the lowest-numbered one. It does not claim multiple.

- ❌ Wake claims all three `dispatch:cell + status:todo` issues in one run
- ✅ Wake claims issue #23 (lowest); leaves #31 and #45 for subsequent runs

### 3.3. Re-read before claim; re-verify after

Never trust the query result or event payload as the final label source. Always re-read immediately before claiming. Always verify after the label transition.

- ❌ Event says "labeled status:todo"; wake skips re-read and claims immediately
- ✅ Re-read; confirm status:todo still present and not already in-progress

### 3.4. Only operator authorizes dispatch gate

The `status:todo` label is the operator's authorization. Planner (interactive Sigma acting as δ) may apply it only when quoting an explicit human approval. Planner MUST NOT apply `status:todo` based on its own judgment that the issue "seems ready."

- ❌ Planner: "Issue #449 looks complete; I'll apply status:todo"
- ✅ Planner: "Operator said 'Step 4 GO' in D35; applying status:todo to #449 and #454 per authorization"

### 3.5. Claim comment is mandatory

Every claim must produce a claim comment with `claimed_by`, `claimed_at`, and `head`. The comment is the audit trail for double-execution detection and the coordination point for operator oversight.

- ❌ Wake transitions labels without writing the claim comment
- ✅ Claim comment is part of the atomic claim sequence in §2.2 step 4

### 3.6. Blocked transition stops the cell

When a cell encounters a dependency it cannot resolve (missing operator input, external gate, unresolvable ambiguity), the cell transitions `in-progress → blocked` with a comment naming the block. The wake does NOT continue with other work after declaring blocked. The blocked issue is NOT re-claimed on next sweep — `status:blocked` is not claimable.

- ❌ Cell hits a block; continues anyway; produces incomplete output
- ✅ Cell transitions to blocked; names the blocking dependency; stops

---

## 4. Verify

### 4.1. Label set completeness check

Confirm all seven `status:*` labels and three type-marker labels are defined and distinct. Confirm no `status:done` label exists. Confirm closed issue is the terminal state.

### 4.2. Claim operation check

Walk through §2.2 steps 1–6. Confirm each step is concrete enough to execute without consulting another document. Confirm the re-verify step (step 5) is present and can detect a failed claim.

### 4.3. Concurrency composition check

Confirm the three layers (workflow concurrency group + claim-time re-read + claim comment) together prevent double execution under normal operating conditions. Identify the residual race condition (concurrent runs in separate workflow files) and confirm the audit trail (claim comment) covers it.

### 4.4. Transition map completeness check

Confirm every state in the lifecycle table (§2.4) has at least one forward transition. Confirm no state is terminal except closed issue. Confirm `blocked` can return to `ready` or `todo`.

### 4.5. Wave-semantics check

Confirm master issues CANNOT carry `dispatch:cell` (structural, not documented-only). Confirm parallel-dispatch exception requires justification. Confirm lowest-issue-number sweep ordering is deterministic.

### 4.6. Drift-handling coverage check

Confirm the four drift patterns in §2.6 cover the named failure modes. Confirm each produces a degraded report with a repair instruction comment, not a silent skip.

---

## 5. Failure modes catalogue

- **D1 — Type confusion.** _Symptom:_ Wake claims a master issue and launches a cell runtime on it. Cell has no bounded scope; produces unbounded work. _Fix:_ §3.1 + `dispatch:cell` absent on master/RFC issues.
- **D2 — Double execution.** _Symptom:_ Two runs claim the same issue; same cell runs twice; duplicate PR or conflicting changes. _Fix:_ §2.3 three-layer anti-double-execution.
- **D3 — Stale event payload.** _Symptom:_ Event fires; issue label changed between event and wake; wake acts on stale state. _Fix:_ §3.3 — always re-read.
- **D4 — Ghost claim.** _Symptom:_ Label edit fails silently; claim comment written; cell launches; issue still shows `status:todo`; next sweep re-claims. _Fix:_ §2.2 step 5 re-verify; abort if `status:in-progress` absent after edit.
- **D5 — Unauthorized planner transition.** _Symptom:_ Planner applies `status:todo` without explicit operator authorization; autonomous cell runs on unapproved work. _Fix:_ §3.4 — planner requires explicit human approval quote before applying `status:todo`.

---

## 6. References

### Upstream protocol

- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — channel convention; claim protocol is the implementation-cell analog of the activation-log write discipline.

### Skills and issues that depend on this protocol

- `agent/wake-template/SKILL.md` (`cnos#450`) — renders `on: issues` labeled trigger, scheduled sweep trigger, concurrency group, dispatch query, and claim/report steps. Cites this skill as source of truth for dispatch behavior.
- `agent/cohere/SKILL.md` (`cnos#444`) — cohere orchestrator dispatches cells per this protocol. `dispatch:cell` + `status:todo` composition is enforced at cohere's dispatch layer.

### Related issues

- `cnos#454` — this skill's authoring issue. ACs in §1–§5 correspond to AC1–AC11 of cnos#454.
- `cnos#450` — wake-template (cites this protocol for template content).
- `cnos#444` — cohere skill (dispatches using this protocol).
- `cnos#449` — registration skill (cells for registration-class issues follow this protocol).

### Empirical evidence

- cn-sigma `threads/adhoc/20260616-coherer-console-vs-wake-worker.md` — interface boundary doctrine; this protocol is the concrete dispatch mechanism for the wake worker side of the boundary.
- D35 directive (2026-06-20T09:30Z) — first operational use of this protocol: operator applied `status:todo` to cnos#449 and cnos#454 via "Step 4 GO" authorization; planner cited explicit human approval per §3.4.

### Authority and stability

This skill is infrastructure-level: every autonomous implementation cell in cnos passes through this protocol. Changes to the claim operation (§2.2), label set (§1.1), or concurrency discipline (§2.3) require operator review — they affect the authorization gate for all autonomous work. The two-label composition (`dispatch:cell` + `status:todo`) is a constitutive constraint; weakening either requirement breaks the authorization model.
