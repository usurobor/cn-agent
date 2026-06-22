---
name: dispatch-protocol
description: Label-gated, protocol-qualified issue queue for implementation cells. Defines the canonical generic label set (cnos.core-owned), the per-package protocol-qualifier label set (concrete-protocol-package-owned per cnos#468 label-doctrine), the claim operation, concurrency discipline, lifecycle transitions, wave semantics, and drift handling by which a human authorizes a package-owned dispatch wake to claim and execute exactly one implementation cell per run. The dispatch wake claims only cells whose `protocol:{P}` qualifier matches its owning package, and launches the matching package's runtime (cnos.cds for software cells, cnos.cdr for research cells, future cnos.cdw for writing cells; the generic cnos.cdd cell-runtime framework is invoked by each concrete package internally).
artifact_class: skill
kata_surface: embedded
governing_question: How does a human authorize exactly one implementation cell per package-owned dispatch wake run, with no double execution, no unauthorized claiming, no cross-protocol claiming, and explicit operator control over the dispatch gate?
visibility: public
parent: agent
triggers:
  - dispatch
  - dispatch protocol
  - label gate
  - protocol qualifier
  - dispatch:cell
  - protocol:cds
  - status:todo
  - claim issue
  - implementation cell
scope: task-local
inputs:
  - GitHub repository with labeled issues (generic labels per cnos#468 + per-package protocol qualifiers)
  - package-owned dispatch wake with repository read/write access (per cnos#467 / cnos#470 wake-provider contract)
  - operator authorization (status:todo label applied)
outputs:
  - at most one executable issue claimed per run
  - claim comment written with wake id + run id + protocol + head sha
  - status labels transitioned (todo → in-progress)
  - matching package runtime launched on the claimed issue
requires:
  - GitHub repository with generic label set per cnos#468 (cnos.core-owned) and the wake's owning protocol's `protocol:{P}` qualifier (per-package-owned)
  - operator-applied status:todo label on target issue
  - repository-level workflow concurrency group (per-protocol scoped per cnos#470)
calls:
  - the matching package's runtime (concrete-protocol-package — e.g. cnos.cds for `protocol:cds` cells); each concrete protocol invokes the generic cnos.cdd cell-runtime framework's γ/α/β/δ contracts internally
---

<!--
section-manifest:
  planned: [frontmatter, core-principle, algorithm, define, unfold-label-set, unfold-claim-operation, unfold-concurrency, unfold-transitions, unfold-wave-semantics, unfold-drift-handling, rules, verify, failure-modes, references]
  completed: [frontmatter, core-principle, algorithm, define, unfold-label-set, unfold-claim-operation, unfold-concurrency, unfold-transitions, unfold-wave-semantics, unfold-drift-handling, rules, verify, failure-modes, references]
-->

# Dispatch Protocol

## Core Principle

**An issue is executable only when ALL of: open, carries `dispatch:cell`, carries exactly one `protocol:{P}` label whose protocol matches the claiming dispatch wake's owning package, carries `status:todo`, AND carries none of `status:in-progress` / `status:blocked` / `status:review` / `status:changes`.** A package-owned dispatch wake claims at most one executable issue per run. Lifecycle (`status:*`) ≠ issue type (`dispatch:cell`) ≠ protocol qualifier (`protocol:{P}`).

The three-label composition is intentional. `status:todo` alone is not enough (master issues, RFC issues, planning issues can all be "todo"). `dispatch:cell + status:todo` alone is not enough either (a cross-protocol cell would still match a research-protocol wake's blunt query). The protocol qualifier — owned by the concrete protocol package per cnos#468 label-doctrine §2 — is what makes the dispatch routable: cnos.cds claims `protocol:cds` cells; cnos.cdr claims `protocol:cdr` cells; future cnos.cdw claims `protocol:cdw` cells. Cnos.cdd is the generic cell-runtime **framework** (γ/α/β/δ contracts) that each concrete protocol package invokes internally — it does not own a protocol qualifier or a dispatch wake.

The three-label composition prevents:
- **Unauthorized claiming** — `dispatch:cell` (type marker) keeps master/RFC/planning issues out of the queue.
- **Premature claiming** — `status:todo` (authorization gate) requires explicit operator intent.
- **Cross-protocol claiming** — `protocol:{P}` (routing qualifier) keeps a software-protocol wake from claiming a research cell (and vice versa).

Each label is necessary; none is sufficient alone. All three are required at claim time, plus none of the excluded statuses.

## Algorithm

1. **Query** — find open issues with `dispatch:cell` + `protocol:{P}` + `status:todo`, where `{P}` is the dispatch wake's owning protocol qualifier (per its wake-provider manifest's `protocol` field, per cnos#470 wake-provider/SKILL.md §3.9). On a sweep run, find all matches; take the lowest issue number. On an `issues: labeled` event run, check only the event issue.
2. **Re-read** — do not trust the query result or event payload alone. Re-read the issue's current labels before claiming.
3. **Verify** — confirm: issue open; `dispatch:cell` present; **exactly one** `protocol:*` label is present AND equals `protocol:{P}` (the wake's owning protocol); **exactly one** `status:*` label is present AND equals `status:todo`; `status:in-progress`, `status:blocked`, `status:review`, `status:changes` all absent.
4. **Claim** — within the GitHub label API's best-effort: remove `status:todo`, add `status:in-progress`, write claim comment with wake id, run id, protocol, and head sha.
5. **Re-verify** — re-read labels after claim. If `status:in-progress` is not present, OR `protocol:{P}` was changed during the claim, abort and release (`status:in-progress → status:todo`); report.
6. **Launch** — pass the claimed issue to the **matching package runtime** (the concrete protocol package owning `protocol:{P}` — e.g. cnos.cds for `protocol:cds`). The package runtime invokes the generic cnos.cdd cell-runtime framework (γ → α → β → δ) internally; the dispatch protocol does not name a specific framework here, since runtime selection is implicit in the protocol qualifier match per cnos#468 §2.2.

---

## 1. Define

### 1.1. Label set (canonical, v0)

Three label namespaces, two owners (per cnos#468 label-doctrine):

| Namespace | Owner | Examples |
|---|---|---|
| `status:*` (lifecycle) | cnos.core | `status:backlog`, `status:ready`, `status:todo`, `status:in-progress`, `status:review`, `status:changes`, `status:blocked` |
| Type markers + dispatchability | cnos.core | `dispatch:cell`, `kind/wave`, `kind/rfc` |
| `protocol:{P}` (protocol qualifier) | concrete protocol package | `protocol:cds` (cnos.cds), `protocol:cdr` (cnos.cdr), future `protocol:cdw` (cnos.cdw) |

The two-layer ownership is doctrine: cnos.core MUST NOT define `protocol:{P}` labels; protocol packages MUST NOT define generic labels (per cnos#468 §3). `cnos.cdd` is the generic cell-runtime **framework** (γ/α/β/δ contracts + cell mechanics) — it does NOT own a protocol qualifier or a dispatch wake; the concrete protocol packages (cnos.cds, cnos.cdr, future cnos.cdw) own theirs.

**Lifecycle labels** — exactly one per open issue:

| Label | Meaning |
|-------|---------|
| `status:backlog` | Identified but not yet ready for implementation |
| `status:ready` | Specified and reviewed; not yet authorized for dispatch |
| `status:todo` | Human-authorized for dispatch; executable when combined with `dispatch:cell` + a `protocol:{P}` |
| `status:in-progress` | Claimed by a package-owned dispatch wake; matching package runtime running |
| `status:review` | Cell complete; awaiting external human/planner review |
| `status:changes` | External (human/planner) review requested changes; awaiting repair before re-dispatch |
| `status:blocked` | Blocked on external dependency or human decision |

Closed issue = done. No `status:done` label in v0.

**Type markers + dispatchability** — zero or one per issue:

| Label | Meaning |
|-------|---------|
| `dispatch:cell` | This issue is an implementation cell — independently shippable, scope-bounded, claimable by a matching-protocol dispatch wake |
| `kind/wave` | Optional; marks a master issue tracking a wave of dependent sub-issues |
| `kind/rfc` | Optional; marks an RFC or review issue (MUST NOT carry `dispatch:cell`) |

**Protocol qualifiers** — required on every `dispatch:cell` issue; exactly one per issue:

| Label | Owning package | Asserts |
|-------|---|---|
| `protocol:cds` | cnos.cds | Software-development cell (CDS protocol; cnos.cds invokes cnos.cdd cell-framework) |
| `protocol:cdr` | cnos.cdr | Research cell (CDR protocol; cnos.cdr invokes cnos.cdd cell-framework) |
| `protocol:cdw` (future) | cnos.cdw | Writing cell (CDW protocol; cnos.cdw invokes cnos.cdd cell-framework) |

Each concrete protocol package's `cn install {pkg}` step creates its `protocol:{P}` label per cnos#468 §3. cnos.core's install does NOT create them. Cross-protocol mismatches (a software-protocol wake encountering `protocol:cdr`) are drift per §2.6.

**Invariants:**
- Exactly one `status:*` label per open issue.
- `dispatch:cell` MUST NOT appear on master issues, RFC issues, or planning issues.
- Every `dispatch:cell` issue MUST carry **exactly one** `protocol:{P}` label.
- A wake claims only cells whose `protocol:{P}` matches its owning package's protocol.

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

Seven named failure modes specific to dispatch. Each has a structural fix in §2.

- **D1 — Type confusion.** Wake claims a master issue or RFC because it carries `status:todo` without `dispatch:cell`. Fix: §3.1 — all three labels required; `dispatch:cell` MUST NOT appear on non-cell issues.
- **D2 — Double execution.** Two wakes (event-driven + scheduled sweep) race; both see `status:todo` and both claim. Fix: §2.3 — per-protocol concurrency group + claim-time re-read + claim comment with run id.
- **D3 — Stale event payload.** Wake receives `issues: labeled` event; trusts the event's label snapshot; issue labels changed since the event fired. Fix: §2.2 — always re-read current labels before claiming; never trust event payload alone.
- **D4 — Ghost claim.** Wake writes claim comment but label transition fails (API error). Cell launches without `status:in-progress` set; next sweep re-claims the same issue. Fix: §2.2 — re-verify `status:in-progress` after claim; abort if absent.
- **D5 — Unauthorized transition.** Planner applies `status:todo` without an explicit human approval quote. Operator intent is bypassed. Fix: §1.2 human gate rule + §3.4.
- **D6 — Protocol missing.** A `dispatch:cell` issue carries `status:todo` but no `protocol:{P}` (or zero protocol qualifiers). Under-specified for dispatch — no wake knows it owns the cell. Fix: §2.2 verify gate rejects; §2.6 `dispatch_protocol_missing` repair instruction comment.
- **D7 — Cross-protocol claim attempt.** A dispatch wake owning protocol `{Q}` encounters a cell labeled `protocol:{P}` where `P ≠ Q`. On scheduled sweep: silent skip (the matching wake claims it). On targeted/attempted claim (event-driven dispatch on this specific issue via a labeled-event trigger or explicit dispatch handle): the wake rejects and reports `dispatch_protocol_mismatch`. Fix: §2.6 `dispatch_protocol_mismatch` handling; the scheduled-vs-targeted distinction matters because a sweep over a heterogeneous queue is expected; a targeted claim attempt landing on a wrong-protocol cell signals upstream label drift or misrouted trigger.

---

## 2. Unfold

### 2.1. Label set administration

Labels are created at the repository level by their owning package's install step. Owners are split per cnos#468 label-doctrine §3:

- **`cn install cnos.core`** creates the generic lifecycle + type-marker labels (`status:*`, `dispatch:cell`, `kind/*`). Manifest: `src/packages/cnos.core/labels.json`. cnos.core MUST NOT create `protocol:{P}` labels — those are package-owned.
- **`cn install cnos.cds`** creates `protocol:cds`. Per-package; cnos.cds owns the label semantics and the dispatch wake (`cds-dispatch`).
- **`cn install cnos.cdr`** creates `protocol:cdr`. Same shape per the cnos.cdr package.
- Future protocol packages (e.g. cnos.cdw) follow the same pattern.

Crossing the split is a doctrine violation per cnos#468 §3:
- ❌ cnos.core creates `protocol:cds` (right-column label out of left-column package)
- ❌ cnos.cds creates `status:todo` (left-column label out of right-column package)
- ✅ Each install creates only its owned labels; installs are idempotent

Creating missing labels (example shape — package installs use the substrate equivalent):
```bash
# cnos.core's install (idempotent; runs on `cn install cnos.core`):
gh label create "status:backlog" --color "#EEEEEE" --repo <owner>/<repo>
gh label create "status:ready"   --color "#BFD4F2" --repo <owner>/<repo>
gh label create "status:todo"    --color "#0075CA" --repo <owner>/<repo>
gh label create "status:in-progress" --color "#E4E669" --repo <owner>/<repo>
gh label create "status:review"  --color "#D4C5F9" --repo <owner>/<repo>
gh label create "status:changes" --color "#E11D48" --repo <owner>/<repo>
gh label create "status:blocked" --color "#B60205" --repo <owner>/<repo>
gh label create "dispatch:cell"  --color "#0E8A16" --repo <owner>/<repo>
gh label create "kind/wave"      --color "#C5DEF5" --repo <owner>/<repo>
gh label create "kind/rfc"       --color "#F9D0C4" --repo <owner>/<repo>

# cnos.cds's install (idempotent; runs on `cn install cnos.cds`):
gh label create "protocol:cds"   --color "#1D76DB" --repo <owner>/<repo>

# cnos.cdr's install:
gh label create "protocol:cdr"   --color "#5319E7" --repo <owner>/<repo>
```

Colors are recommendations; operators may change colors without affecting protocol behavior.

### 2.2. Claim operation

The claim sequence. Every step is required. No step may be skipped. `{P}` is the wake's owning protocol (declared in the wake-provider manifest's `protocol` field per cnos#470 §3.9).

```
1. Query for executable issues matching this wake's protocol:
   - On scheduled sweep: gh issue list --repo <repo> \
       --label "dispatch:cell" --label "protocol:{P}" --label "status:todo" \
       --state open --json number,labels
     → sort by issue number ascending; take lowest
   - On issues:labeled event: gh issue view <event-issue-number> --repo <repo> --json number,labels,state

2. Re-read labels:
   gh issue view <issue-number> --repo <repo> --json labels,state
   → Do not proceed with event payload labels alone

3. Verify (all must pass; abort on any failure):
   - issue.state == "open"
   - "dispatch:cell" in labels
   - exactly one label matching "protocol:*" AND it equals "protocol:{P}"
     • zero protocol labels → dispatch_protocol_missing (§2.6); do NOT claim
     • multiple protocol labels → dispatch_label_drift (§2.6); do NOT claim
     • protocol:{Q} where Q ≠ P → dispatch_protocol_mismatch (§2.6) on targeted claim;
       silent skip on sweep (the matching wake claims it)
   - exactly one label matching "status:*" AND it equals "status:todo"
     • zero status labels → dispatch_label_drift; do NOT claim
     • multiple status labels → dispatch_label_drift; do NOT claim
   - "status:in-progress" NOT in labels
   - "status:blocked" NOT in labels
   - "status:review" NOT in labels
   - "status:changes" NOT in labels
   → On any verify failure: log "not claimable: <reason>"; on a drift case post the
     repair-instruction comment (§2.6); exit 0 (do not error; scheduled sweep
     re-attempts on the next firing if the drift is repaired).

4. Claim:
   gh issue edit <issue-number> --repo <repo> \
     --remove-label "status:todo" \
     --add-label "status:in-progress"

   gh issue comment <issue-number> --repo <repo> --body "$(cat <<EOF
   claimed_by: <wake-name>/<run-id>      # e.g. cds-dispatch/<actions-run-id>
   protocol: {P}                         # e.g. cds
   claimed_at: <utc-timestamp>
   head: <current-main-sha>
   EOF
   )"

5. Re-verify:
   gh issue view <issue-number> --repo <repo> --json labels
   → Confirm "status:in-progress" is present AND "protocol:{P}" is still present
   → If absent OR protocol qualifier changed during claim: release
     (gh issue edit --add-label "status:todo" --remove-label "status:in-progress");
     comment with race-detection note; exit 1 (do not launch the runtime)

6. Launch the matching package runtime (only after step 5 confirms):
   → pass the claimed issue to the concrete protocol package matching protocol:{P}
     (e.g. protocol:cds → cnos.cds; protocol:cdr → cnos.cdr; future
     protocol:cdw → cnos.cdw). The package runtime invokes the generic cnos.cdd
     cell-runtime framework (γ → α → β → δ) internally. The dispatch protocol
     does not name a specific framework here — runtime selection is implicit in
     the protocol qualifier match per cnos#468 §2.2.
```

- ❌ Wake trusts the query result; skips re-read
- ✅ Wake re-reads labels immediately before every claim
- ❌ Wake hard-codes the launch target as cnos.cdd
- ✅ Launch target is the concrete protocol package matching `protocol:{P}`; the package internally uses cnos.cdd as the generic cell framework

### 2.3. Concurrency and anti-double-execution

GitHub label API is not a true compare-and-swap primitive. Anti-double-execution discipline is a composition of three layers:

**Layer 1 — Per-protocol concurrency group (substrate serialization)**

Each package-owned dispatch wake declares a per-protocol concurrency group (per cnos#470 wake-provider/SKILL.md §2.2 `concurrency_intent`). The renderer maps to substrate-specific syntax; for GitHub Actions:
```yaml
concurrency:
  group: ${{ github.repository }}-cds-dispatch-${agent}   # for cnos.cds's wake
  cancel-in-progress: false
```

At most one dispatch job per (repository × wake) runs at a time. A second trigger queues behind the first, not racing. Per-protocol scoping (rather than repository-wide) means a cds-dispatch firing does not block a cdr-dispatch firing — orthogonal protocols can dispatch concurrently because they claim disjoint label sets.

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
- Each sub carries exactly one `protocol:{P}` label naming the concrete protocol that will execute it. A wave may be homogeneous (all subs `protocol:cds`) or heterogeneous (some subs `protocol:cds`, some `protocol:cdr`) — the master `kind/wave` issue coordinates across protocols but does not itself carry `protocol:{P}`.
- Each sub is independently dispatchable only when labeled `status:todo`.
- Default: sequential dispatch — sub N carries `status:todo` only after sub N-1 closes.
- Exception: if sub N and sub N+1 are explicitly independent AND merge-safe (`dispatch:cell` issues that can land in either order without conflict), both may carry `status:todo` simultaneously for parallel dispatch. Parallel is the exception; justification required.
- Cross-protocol subs (sub N is `protocol:cds`, sub N+1 is `protocol:cdr`) parallelize naturally — different dispatch wakes claim them; per-protocol concurrency groups don't block each other.

Lowest-issue-number sweep ordering: when multiple subs are simultaneously `status:todo` and match a wake's protocol, the wake claims the lowest issue number first. Priority labels are deferred to v1.

### 2.6. Drift handling

A drift state is an issue with an unexpected or inconsistent label combination. Do not claim drifted issues.

| Drift pattern | Degraded reason | Repair |
|---------------|-----------------|--------|
| Multiple `status:*` labels | `dispatch_label_drift` | Remove extra status labels; leave one |
| No `status:*` label | `dispatch_label_drift` | Add correct status label |
| Multiple `protocol:*` labels on a `dispatch:cell` issue | `dispatch_label_drift` | Remove extra protocol labels; leave the one matching the cell's actual runtime |
| `dispatch:cell` on a master/RFC issue | `dispatch_type_mismatch` | Remove `dispatch:cell` from non-cell issue |
| `dispatch:cell` issue without any `protocol:*` label | `dispatch_protocol_missing` | Apply the correct `protocol:{P}` for the cell's runtime (e.g. `protocol:cds` for software cells); under-specified for dispatch until applied |
| Cell labeled `protocol:{Q}` where Q ≠ this wake's protocol (targeted/attempted claim) | `dispatch_protocol_mismatch` | Do NOT claim from this wake; the matching-protocol wake claims it. Sweep over heterogeneous queue: silent skip. Targeted claim landing on wrong-protocol cell: comment naming the mismatch; signals upstream label drift or misrouted trigger |
| `status:in-progress` without a claim comment | `dispatch_ghost_claim` | Clear `status:in-progress`; re-apply `status:todo` if still authorized |

On drift detection (other than the silent-skip sweep case):
```
outcome: degraded
degraded_reason: dispatch_label_drift | dispatch_type_mismatch |
                 dispatch_protocol_missing | dispatch_protocol_mismatch |
                 dispatch_ghost_claim
operator_action: repair instruction (specific labels to add/remove)
```

Post a comment naming the drift and the repair instruction. Do NOT claim the drifted issue. The matching-protocol wake's next firing reattempts after operator repair.

**Scheduled-sweep-vs-targeted distinction** — important for `dispatch_protocol_mismatch`:
- **Sweep over heterogeneous queue:** expected behavior; the wake's `protocol:{P}` selector simply doesn't match — the wake silently moves on. Not a drift event. The matching-protocol wake claims it on its own firing.
- **Targeted claim attempt** (event-driven dispatch on this specific issue via labeled-event trigger or explicit dispatch handle): the wake was steered to this issue but the protocol doesn't match. This signals upstream label drift or a misrouted trigger and MUST surface for repair. Drift event.

The distinction matters because the two paths reach the same code (claim attempt → verify gate fail) but call for different operator responses.

---

## 3. Rules

### 3.1. All three labels required for claim

A wake MUST NOT claim an issue missing any of: `dispatch:cell`, exactly one matching `protocol:{P}`, `status:todo`. No exceptions. Master issues, RFC issues, and planning issues MUST NOT carry `dispatch:cell` — this is the structural lock against type confusion. Every `dispatch:cell` issue MUST carry exactly one `protocol:{P}` — the routing qualifier; missing or multiple are drift per §2.6.

- ❌ "This issue has status:todo; I'll claim it for the wave tracker"
- ❌ "This issue has status:todo + dispatch:cell but no protocol qualifier; I'll claim it"
- ❌ "This issue has status:todo + dispatch:cell + protocol:cdr; I'll claim it (I'm cds-dispatch)"
- ✅ "No dispatch:cell label; skip"
- ✅ "No protocol qualifier; surface dispatch_protocol_missing"
- ✅ "Protocol:cdr cell — not my protocol; sweep silently skips, targeted claim posts dispatch_protocol_mismatch comment"

### 3.2. At most one claim per run

A single wake run claims at most one issue. Scheduled sweep queries all `dispatch:cell + protocol:{P} + status:todo` issues matching the wake's owning protocol, and takes the lowest-numbered one. It does not claim multiple.

- ❌ Wake claims all three matching `dispatch:cell + protocol:{P} + status:todo` issues in one run
- ✅ Wake claims issue #23 (lowest matching); leaves #31 and #45 for subsequent runs

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

### 3.7. Wake claims only its owned protocol

A dispatch wake owned by package `{P}` claims **only** issues whose `protocol:{P}` matches. Cross-protocol claims are a label-doctrine violation per cnos#468 §2.1 + this skill's §1.1 invariant. The wake's owning protocol is fixed at install time per its wake-provider manifest's `protocol` field (per cnos#470 §3.9); a wake cannot dynamically claim across protocols.

A multi-protocol repo (cnos cells of both `protocol:cds` and `protocol:cdr`) runs MULTIPLE dispatch wakes — one per concrete protocol package. Each wake's selector is scoped to its protocol; together they cover the multi-protocol queue.

- ❌ A single cds-dispatch wake claims `protocol:cdr` cells "because the operator forgot to install cnos.cdr"
- ✅ Without an installed cdr-dispatch wake, `protocol:cdr` cells sit in the queue; the admin wake surfaces this to the operator per cnos#470's defer-path for missing dispatch wakes

### 3.8. Launch the matching package runtime, not a fixed framework

The dispatch wake's launch step (§Algorithm step 6) passes the claimed cell to the **matching package's runtime** — the concrete protocol package that owns `protocol:{P}`. Each concrete protocol package invokes the generic cnos.cdd cell-runtime framework's γ/α/β/δ contracts internally, but the dispatch protocol skill MUST NOT hard-code cnos.cdd as the runtime target — that would conflate the framework (cnos.cdd) with the concrete protocol (cnos.cds / cnos.cdr / cnos.cdw).

- ❌ Wake's launch step: "pass issue to cnos.cdd's δ"
- ✅ Wake's launch step: "pass issue to the matching package runtime (e.g. cnos.cds for protocol:cds); the package invokes cnos.cdd framework internally"

This is the routing invariant cnos#467 / cnos#468 / cnos#470 / cnos#483 (Sub 4) all enforce. cnos.cdd is the cell-runtime framework; cnos.cds (and siblings) are the concrete protocols that use it.

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
- **D6 — Protocol missing.** _Symptom:_ A `dispatch:cell` issue carries `status:todo` but no `protocol:{P}` label. Under-specified for dispatch — no wake knows which package's runtime should claim it. Every dispatch wake's selector rejects it (silent skip on sweep, since no wake's selector matches). _Fix:_ §2.6 — `dispatch_protocol_missing` repair-instruction comment; operator applies the correct `protocol:{P}`; the matching wake claims on next firing.
- **D7 — Cross-protocol claim attempt.** _Symptom:_ A dispatch wake owning `{Q}` is steered to a cell labeled `protocol:{P}` (P ≠ Q) via a labeled-event trigger or explicit dispatch handle. Signals upstream label drift or a misrouted trigger. _Fix:_ §2.6 + §3.7 — targeted claim attempts reject with `dispatch_protocol_mismatch` and a comment naming the mismatch; scheduled sweep silently skips (the matching wake claims). The wake never claims outside its owned protocol.
- **D8 — Hard-coded runtime in launch step.** _Symptom:_ A dispatch wake's launch step names `cnos.cdd` as the cell runtime (or any other fixed framework), conflating the generic cell-runtime framework with the concrete protocol package. Future protocols can't be added without amending the dispatch skill. _Fix:_ §3.8 — launch step passes to the matching package runtime; the matching package invokes cnos.cdd's contracts internally. cnos.cdd is framework, not concrete protocol.

---

## 6. References

### Upstream protocol

- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — channel convention; claim protocol is the implementation-cell analog of the activation-log write discipline.

### Skills and issues that depend on this protocol

- `agent/label-doctrine/SKILL.md` (`cnos#468`, MERGED) — defines the canonical generic label set (cnos.core-owned) and the per-protocol qualifier rule (concrete-package-owned). This skill's §1.1 cites label-doctrine as the doctrinal source; the dispatch protocol operationalizes label-doctrine's contract at claim time.
- `agent/wake-provider/SKILL.md` (`cnos#470`, MERGED) — `cn.wake-provider.v1` declaration contract. Dispatch wakes (`role: dispatch`) declare `protocol`, `selector`, and dispatch-shape `output_contract` per its §3.9; this skill's claim algorithm reads those fields at install time.
- `agent/wake-template/SKILL.md` (`cnos#450`) — renders package-owned dispatch wakes from `cn.wake-provider.v1` manifests; renders `on: issues` labeled trigger filtered by `protocol:{P}`, scheduled sweep trigger, per-protocol concurrency group, dispatch query, and claim/report steps. Cites this skill as source of truth for claim behavior.
- `cnos.cds/orchestrators/cds-dispatch/wake-provider.json` (`cnos#467 Sub 4 / PR #483`, MERGED, declaration-only) — first reference dispatch wake; declares `protocol: cds`, the standard `selector` (`dispatch:cell + protocol:cds + status:todo` minus excluded statuses), and the dispatch-shape `output_contract` invoking cnos.cdd as `cell_runtime`. cnos.cdd is the framework; cnos.cds is the concrete protocol.
- `agent/cohere/SKILL.md` (`cnos#444`) — cohere orchestrator dispatches cells per this protocol. The three-label composition is enforced at cohere's dispatch layer.

### Related issues

- `cnos#454` — this skill's authoring issue. ACs in §1–§5 correspond to AC1–AC11 of cnos#454, amended for the protocol-qualifier doctrine landed via PR #480.
- `cnos#468` — label doctrine; the two-layer ownership of labels (generic by cnos.core; protocol qualifiers by concrete-protocol packages).
- `cnos#470` — wake-provider contract; dispatch wakes declare `protocol` + `selector` + dispatch-shape `output_contract`.
- `cnos#480` — doctrine correction: cnos.cdd is framework, not concrete protocol; cnos.cds is the concrete software-development protocol.
- `cnos#483` — Sub 4 cnos.cds dispatch wake provider (first dispatch reference instance); declaration-only until renderer + δ wake-invoked mode land.
- `cnos#479` — cutover-A (admin wake live in production); the other half of the two-wake architecture cnos#467 establishes.
- `cnos#450` — wake-template (cites this protocol for template content).
- `cnos#444` — cohere skill (dispatches using this protocol).
- `cnos#449` — registration skill (cells for registration-class issues follow this protocol).

### Empirical evidence

- cn-sigma `threads/adhoc/20260616-coherer-console-vs-wake-worker.md` — interface boundary doctrine; this protocol is the concrete dispatch mechanism for the wake worker side of the boundary.
- D35 directive (2026-06-20T09:30Z) — first operational use of this protocol: operator applied `status:todo` to cnos#449 and cnos#454 via "Step 4 GO" authorization; planner cited explicit human approval per §3.4.

### Authority and stability

This skill is infrastructure-level: every autonomous implementation cell in cnos passes through this protocol. Changes to the claim operation (§2.2), label set (§1.1), or concurrency discipline (§2.3) require operator review — they affect the authorization gate for all autonomous work. The three-label composition (`dispatch:cell` + `protocol:{P}` + `status:todo`) is a constitutive constraint; weakening any of the three requirements breaks the model (`dispatch:cell` keeps non-cells out; `protocol:{P}` routes cells to their owning runtime; `status:todo` gates on operator authorization). cnos.cdd as the generic cell-runtime framework MUST NOT be hard-coded as a launch target by this skill or by any dispatch wake's prompt — the concrete protocol packages (cnos.cds / cnos.cdr / cnos.cdw) own the launch targets, and they invoke cnos.cdd's contracts internally.
