---
name: label-doctrine
description: cnos's GitHub issue-label control plane — generic lifecycle and dispatchability labels owned by cnos.core, plus the per-package protocol qualifier rule that namespaces multi-protocol dispatch.
artifact_class: skill
kata_surface: none
governing_question: Which labels does cnos use, who owns each, and what may Sigma do with them?
visibility: public
parent: agent
triggers:
  - label
  - status
  - dispatch
  - protocol qualifier
  - issue queue
scope: task-local
inputs:
  - none (doctrine reference)
outputs:
  - canonical generic label set
  - protocol qualifier rule
  - ownership split between cnos.core and protocol packages
  - Sigma admin boundary on labels
requires:
  - none (foundational; consumed by `cnos.core/skills/agent/dispatch-protocol/SKILL.md` per cnos#454, the wake-template renderer per cnos#450, and the agent-admin wake provider per cnos#467)
---

# Label doctrine

cnos's issue-label control plane. Two layers, two owners.

## Core Principle

**Labels are a control plane, not free text.** Each label has an owner (a package), a meaning (what carrying it asserts about the issue), and a role in the dispatch protocol. Without ownership, label semantics drift across packages; without semantics, the dispatch protocol cannot select work safely.

Labels split into two layers:

- **Generic** — lifecycle and dispatchability. Owned by `cnos.core`. The same set in every cnos-using repo.
- **Per-protocol qualifier** — names which package's runtime owns a given cell. Owned by the protocol package. New packages add their qualifier as part of their install.

Sigma (the human-facing admin/console actor) may *apply* labels when authorized by the operator. Sigma does **not** own label semantics, define new labels, or execute the cells the labels gate.

---

## 1. Generic label set (owned by cnos.core)

The generic labels live at `cnos.core`. They are installed in a repo by `cn install cnos.core`. The data manifest is at `src/packages/cnos.core/labels.json`.

### 1.1. Lifecycle labels

Every open issue carries **exactly one** `status:*` label. The label names where the issue is in its lifecycle.

| Label | Meaning | Allowed transitions |
|---|---|---|
| `status:backlog` | well-formed scope but not yet refined to dispatch readiness | `status:ready` |
| `status:ready` | spec'd; ACs converged; awaiting operator authorization | `status:todo`, `status:backlog` |
| `status:todo` | operator-authorized; dispatch wake may claim | `status:in-progress` |
| `status:in-progress` | claimed by a dispatch wake; cycle running | `status:review`, `status:blocked` |
| `status:review` | cell complete; awaiting external/operator review of the receipt/result | `status:changes`, closed |
| `status:changes` | external review requested changes; issue must be revised before re-dispatch | `status:ready`, `status:todo` |
| `status:blocked` | gated on external input (operator, infra, external authority) | `status:ready`, `status:todo` |

Closed issue = done. There is no `status:done` (closing the issue is the terminal state).

### 1.2. Dispatchability label

| Label | Meaning |
|---|---|
| `dispatch:cell` | this issue is an executable implementation cell. Eligible to be claimed by a dispatch wake when paired with `status:todo` and a matching `protocol:{id}`. |

Issues that are NOT cells (master trackers, RFCs, planning issues, design notes, conversations) MUST NOT carry `dispatch:cell`. A `kind/wave` (or `tracking`) master that coordinates cells is not itself a cell — sub-issues are.

---

## 2. Protocol qualifier (owned per package)

Every `dispatch:cell` issue carries **exactly one** `protocol:{id}` label. The qualifier names which package's runtime owns this cell's execution.

| Qualifier | Owner | What its presence asserts |
|---|---|---|
| `protocol:cds` | `cnos.cds` | runtime = CDS software-cell shape (per cnos.cds's spec) |
| `protocol:cdr` | `cnos.cdr` | runtime = CDR research-cell shape (per cnos.cdr's spec) |
| `protocol:cdw` | `cnos.cdw` (when the package exists) | runtime = CDW writing-cell shape (per cnos.cdw's spec) |
| `protocol:{name}` | future package | each future protocol package registers its own qualifier at install |

**Note on cnos.cdd:** `cnos.cdd` is the generic *cell-runtime framework* — it defines the γ/α/β/δ role contracts, the cell mechanics, and the dispatch protocol that the concrete protocol packages (`cnos.cds`, `cnos.cdr`, `cnos.cdw`, ...) all use under the hood. It does **not** own a `protocol:cdd` qualifier of its own; cells are always labeled with the qualifier of the concrete protocol whose runtime executes them (a software cell carries `protocol:cds`, not `protocol:cdd`).

The qualifier label is created in the repo when the owning package is installed via `cn install {pkg}`. cnos.core does NOT create these labels — that would conflate the layers.

### 2.1. Required for dispatch

The authoritative dispatchable selector (defined in detail by `cnos.core/skills/agent/dispatch-protocol/SKILL.md` per cnos#454):

```
open
+ dispatch:cell
+ protocol:{P}            ← matches the wake's owning package
+ status:todo
- status:in-progress
- status:blocked
- status:review
- status:changes
```

A wake owned by package `{P}` claims **only** issues whose `protocol:{P}` matches. Cross-protocol claims are a protocol violation. A `dispatch:cell + status:todo` issue WITHOUT a `protocol:{id}` label is under-specified and is rejected by every dispatch wake's selector.

### 2.2. Where the cell goes

Once claimed, the issue is passed to the owning package's runtime — whichever package owns the matched `protocol:{id}`. The generic dispatch protocol does not name any specific runtime; runtime selection is implicit in the qualifier match.

---

## 3. Ownership split

The two layers have distinct owners and distinct install paths.

| Layer | Owner package | Install responsibility |
|---|---|---|
| Generic lifecycle + dispatchability | `cnos.core` (or a future `cnos.agent` sub-package) | `cn install cnos.core` creates `status:backlog`, `status:ready`, `status:todo`, `status:in-progress`, `status:review`, `status:changes`, `status:blocked`, `dispatch:cell` |
| Per-protocol qualifier | each concrete protocol package | `cn install cnos.cds` creates `protocol:cds`; `cn install cnos.cdr` creates `protocol:cdr`; per-package (cnos.cdd is the framework, not a concrete protocol; it does not create a protocol qualifier) |

Crossing the split is a doctrine violation:

- cnos.core MUST NOT create `protocol:{id}` labels — those are package-owned
- A protocol package MUST NOT create generic lifecycle labels — those are cnos.core-owned
- A protocol package MAY depend on cnos.core (its install assumes the generic set exists in the repo)

Installs are idempotent on label creation. Repeated `cn install cnos.core` runs are no-ops on generic label state; repeated `cn install cnos.cds` (or `cnos.cdr`, future `cnos.cdw`) runs are no-ops on each concrete protocol's `protocol:{id}` label state.

---

## 4. Sigma admin boundary

Sigma is the human-facing admin/console actor. Per `cn-sigma:.cn-sigma/spec/OPERATOR.md §"CDD role assignment"`, Sigma plays δ in CDD cycles and the agent-admin role in channel/wake work. Sigma's relationship with labels is bounded.

### 4.1. Sigma may

- **apply** `status:*` labels when instructed by the operator (e.g., flip `status:ready → status:todo` after operator approval)
- **apply** `protocol:{id}` labels when the issue body clearly cites a single protocol; ask the operator when ambiguous
- **apply** `dispatch:cell` when an issue clearly meets the cell criteria (concrete ACs, single-protocol scope, ready for implementation)
- **report** which issues are dispatchable under each protocol
- **suggest** label changes to the operator

### 4.2. Sigma may not

- **define** new labels or alter the semantics of existing labels — that is package-owned doctrine (this skill + the protocol packages' contracts)
- **invent** new lifecycle states or dispatchability markers
- **execute** the cells the labels gate — execution is the dispatch wake's job, and each dispatch wake belongs to a protocol package
- **own** the issue queue or its selection algorithm — that's the dispatch protocol (cnos#454)
- **own** any protocol's `protocol:{id}` label — that's the protocol package

Sigma's role with labels is **route / refine / report**, not **define / execute**. This matches Sigma's broader admin/console scope: read directives, surface state, prepare work for the right runtime, never run the runtime itself.

---

## 5. Manifest

The generic labels' data manifest is at `src/packages/cnos.core/labels.json`. The manifest is data-only — a `cn install cnos.core` consumer (out of scope here; downstream cycle) reads it and creates the labels in the repo. The manifest lives in `cnos.core` so it travels with the package.

Each entry has at minimum:

| Field | Purpose |
|---|---|
| `name` | label name (e.g., `status:todo`) |
| `description` | one-line meaning |
| `color` | hex color (no `#`) |
| `owner` | always `cnos.core` for generic labels |
| `group` | `lifecycle` or `dispatchability` |

Implementers MAY add fields (e.g., `constraints`, `notes`); the manifest schema is forward-compatible.

---

## 6. Cross-references

This skill is part of the wake-orchestration cluster (cnos#467). It defines doctrine consumed by:

- **cnos#454** (`agent/dispatch-protocol`) — defines the claim mechanics, lifecycle transitions, drift handling, and concurrency discipline; consumes this skill's label set and ownership rules
- **cnos#450** (`agent/wake-template`) — renders package-owned wake providers; each rendered wake pins its `protocol:{P}` filter at render time per this skill's qualifier rule
- **cnos#467** (`agent/wake-orchestration`) — the master architecture; this skill is its first sub (cnos#468)

Adjacent operator/persona doctrine:

- `cn-sigma:.cn-sigma/spec/OPERATOR.md §"CDD role assignment"` — Sigma's CDD role; admin boundary context
- `cn-sigma:.cn-sigma/spec/PERSONA.md` — Sigma identity

The skill format:

- `cnos.core/skills/agent/cap/SKILL.md` — example skill format reference

---

## 7. Failure modes

- **Cross-layer label creation.** A package creates a label it does not own (e.g., cnos.cds creates `status:todo`, or cnos.core creates `protocol:cds`). _Fix:_ per-package install commands create only labels in the package's owned set; `cn install` checks ownership before write.
- **Missing protocol qualifier.** A `dispatch:cell` issue without `protocol:{id}`. _Fix:_ dispatch wakes reject the issue with `degraded_reason: dispatch_protocol_missing` and a repair-instruction comment (per cnos#454 AC9).
- **Cross-protocol mismatch.** A wake owning `{P}` encounters a cell labeled `protocol:{Q}` where `Q ≠ P`. The handling differs by the wake's entry path:
  - **Scheduled sweep:** the wake passes over the cell silently (it is not eligible under the wake's selector). The matching wake claims it on its own sweep. Not a drift event.
  - **Attempted claim / event-targeted dispatch** (e.g., the wake was invoked on this specific issue via a label-trigger or explicit dispatch handle and discovers the protocol mismatch at claim-time): the wake **rejects** the claim and **reports** `degraded_reason: dispatch_protocol_mismatch` with a comment naming the mismatch (per cnos#454 AC9 drift-handling catalogue). This is a drift event.
  The distinction matters because a sweep over a heterogeneous queue is expected behavior, but a targeted claim attempt landing on a wrong-protocol cell signals upstream label drift or a misrouted trigger and must surface for repair.
- **Sigma defines new label.** Sigma adds a label not declared by any package's doctrine. _Fix:_ this is a Sigma-admin boundary violation; new labels must originate from a package's doctrine update + manifest change.
- **Lifecycle skip.** An issue moves from `status:ready` directly to `status:in-progress` without `status:todo`. _Fix:_ dispatch wakes require `status:todo` for claim; transition discipline is enforced by the dispatch protocol (cnos#454 AC4 + AC7).
- **Body/label authorization drift.** A design-first issue's body carries stale hold-state prose ("Not dispatched — status:ready ... dispatch on explicit operator authorization.") that contradicts its actual `status:*` label (cnos#614 → cnos#633 recurrence). _Fix:_ the `status:*` label alone is the source of truth for dispatch readiness — see `dispatch-protocol/SKILL.md` §1.2 "Labels are the sole source of truth for dispatch readiness" (this skill does not restate that doctrine; §1.2 is its one canonical home) and its D13 failure-mode entry. `cn issues dispatch --issue N` (cnos#640) performs the authorization label flip and the legacy body cleanup as one operator-invoked operation.

---

## 8. Versioning

This skill is doctrine-adjacent. Future changes follow the constitutive-change approval discipline that governs other doctrine-adjacent skills. Drift between this skill and `cnos.core/labels.json` is resolved in favor of this skill (the manifest must match the doctrine).

If a future cycle introduces a new lifecycle state, a new dispatchability marker, or a new label-ownership tier, this skill is the canonical source; the manifest follows.
