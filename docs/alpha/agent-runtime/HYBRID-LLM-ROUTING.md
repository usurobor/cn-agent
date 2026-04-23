# Hybrid LLM Routing for cnos

**Issue:** TBD
**Version:** 0.1.0
**Status:** Draft
**Mode:** MCA
**Active Skills:** cdd/design, eng/evolve, eng/process-economics
**Engineering Level:** L7

---

## Problem

cnos is moving toward a package-driven, mechanical runtime with:

- package-installed skills
- exact-dispatch commands
- orchestrators as workflow artifacts
- extensions as capability providers

But its current body/runtime model still lacks a coherent policy for model selection. Without an explicit routing layer, the system will drift into one of two bad shapes:

1. **Claude-first everywhere** — expensive, network-dependent, slower, and misaligned with the offline-first / local-first direction of cnos
2. **Local-first everywhere** — cheaper and resilient, but unreliable on high-context, cross-package, or architectural work where current local models still fail too often

The gap is not "support local models." The gap is:

> The runtime has no explicit, inspectable policy for when a task should stay local and when it should escalate to a remote frontier model.

---

## Constraints

- The routing decision must live in the body/runtime policy layer, not in a skill.
- Model providers should be treated as capability providers, not as skills or commands.
- The baseline path should be local-first.
- Escalation to Anthropic must be deterministic and auditable.
- The decision must be possible from data the runtime can observe:
  - task class
  - context size
  - package/dependency breadth
  - file count
  - validation/confidence signals
- Routing policy must not require prompt-level improvisation by the agent.
- The first version should be simple enough to implement in the Go kernel without a large planning engine.

---

## Challenged Assumption

The challenged assumption is:

> "Model selection can be left to the agent's judgment or to free-form prompt instructions."

That is not robust enough for cnos. Model routing is a body-plane policy decision:

- cost
- locality
- resilience
- capability limits
- cross-package reasoning budget

Those are runtime concerns, not skill concerns.

---

## Impact Graph

### Downstream consumers

- runtime contract rendering
- body capability registry
- command execution
- orchestrator execution
- receipts / traceability
- doctor / status
- deployment economics
- local/offline resilience

### Upstream producers

- command metadata
- orchestrator manifests
- runtime capability inventory
- package graph / dependency graph
- packed context estimator
- file-touch set / artifact scope
- local model confidence/validation results

### Authority relationships

- **Routing policy** is core/runtime authority
- **Model providers** are capability providers
- Skills may influence task shape but do not choose the provider directly
- Commands/orchestrators may declare routing hints, but the runtime makes the final decision

---

## Proposal

Implement a hybrid LLM routing layer in the cnos runtime body.

### 1. Architectural placement

#### Core runtime owns:

- routing matrix
- context estimation
- package/dependency breadth estimation
- fallback/escalation policy
- receipts and observability
- provider dispatch

#### Capability providers own:

- local model execution
- remote model execution

Examples:

- `cnos.llm.local`
- `cnos.llm.anthropic`

#### Skills do not own:

- provider choice
- routing policy
- token/depth thresholds

Skills may shape the task. The body chooses the provider.

### 2. Routing model

The runtime classifies every model-bound task into one of three routing outcomes:

- **local**
- **remote**
- **local-then-remote** (attempt local first, escalate on explicit failure signal)

### 3. Task classes

#### Local-default tasks

Use the local model by default when **all** of the following are true:

- task is bounded and deterministic:
  - single-file edit
  - small multi-file edit in one package
  - CI workflow generation
  - manifest / lockfile / config edits
  - isolated code generation or transformation
  - local review/commentary over one artifact
- estimated packed context ≤ 8k tokens
- package/dependency depth ≤ 2
- context file set ≤ 10 files
- touched package count ≤ 1
- no public contract / architecture / release boundary is involved

#### Remote-default tasks

Escalate directly to Anthropic when **any** of the following is true:

- estimated packed context > 12k tokens
- package/dependency depth > 2
- context file set > 15 files
- touched package count > 2
- task class is:
  - architecture design
  - cross-package refactor
  - release
  - post-release assessment
  - policy/spec reconciliation
  - runtime contract / package-system / extension-boundary change
- the task requires macro coherence across multiple registries, packages, or protocol surfaces

#### Local-then-remote tasks

Attempt local first, then escalate if:

- local output fails validation
- local output returns explicit low-confidence / insufficient-context signal
- generated patch fails deterministic checks
- the first pass expands the required context beyond local thresholds

### 4. Thresholds

#### Recommended v1 thresholds

**Local band** (all must hold):

- `context_tokens_estimate` ≤ 8000
- `dependency_depth` ≤ 2
- `context_file_count` ≤ 10

**Hard escalation band** (any triggers remote):

- `context_tokens_estimate` > 12000
- `dependency_depth` > 2
- `context_file_count` > 15
- `package_count` > 2

Anything between these bands is:

- local-first for deterministic implementation tasks
- remote-first for architectural / policy / release tasks

### 5. Dependency depth definition

`dependency_depth` is the maximum package/module-edge distance required to gather the context needed to answer the task correctly.

Examples:

- single package file edit → depth 0–1
- file + direct neighbor package contract → depth 2
- package graph / runtime registry / multi-package refactor → depth 3+

This is intentionally structural, not semantic.

### 6. Receipts and observability

Every routed model call must leave a receipt with:

- **route decision:** local, remote, local-then-remote
- **reason set:**
  - token threshold
  - dependency depth
  - file count
  - package count
  - task class
  - validation failure
- **provider chosen**
- **model name**
- **fallback/escalation status**
- **latency**
- **token/input estimate**
- **confidence / validation outcome** if available

These receipts are required to tune thresholds later.

### 7. Runtime contract integration

Expose provider availability and current routing policy under `body.capabilities` or a dedicated body sub-surface:

```json
"body": {
  "capabilities": { ... },
  "llm": {
    "providers": [
      { "name": "local", "available": true, "kind": "local" },
      { "name": "anthropic", "available": true, "kind": "remote" }
    ],
    "routing_policy": {
      "default": "local-first",
      "local_max_tokens": 8000,
      "remote_min_tokens": 12000,
      "max_dependency_depth_local": 2,
      "max_files_local": 10,
      "max_files_remote_direct": 15
    }
  }
}
```

### 8. Failure policy

If the remote provider is unavailable:

- remote-default tasks **fail closed** unless the operator explicitly allows degraded local fallback

If the local provider is unavailable:

- local-default tasks may escalate directly to remote if credentials/policy allow it

**Do not silently swap provider classes without a visible receipt.**

---

## Leverage

This design:

- keeps the baseline execution layer local and resilient
- reduces unnecessary Claude spend
- preserves frontier capacity for macro-coherence tasks
- turns routing into inspectable runtime policy instead of prompt folklore
- aligns with cnos's existing separation of:
  - core policy
  - capability providers
  - skills
  - commands
  - orchestrators

---

## Negative Leverage

This adds:

- one more runtime policy surface
- one more body capability family
- context estimation logic
- tuning burden for thresholds
- a new class of routing receipts and doctor/status exposure

If done badly, it could create:

- brittle threshold gaming
- hidden provider switching
- duplicated policy between commands/orchestrators and core

---

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Claude-first always | Simple, high quality | Expensive, network-dependent, weakens local autonomy | Rejected |
| Local-first always | Cheap, resilient, local | Unreliable on architectural and high-context work | Rejected |
| Agent decides provider ad hoc | Flexible | Not auditable, not deterministic, drifts into prompt behavior | Rejected |
| Core routing policy + provider capability split | Explicit, inspectable, resilient, tunable | More runtime policy work | **Chosen** |

---

## Process Cost / Automation Boundary

**Automate:**

- token/context estimation
- file-count estimation
- package/dependency breadth estimation
- task-class defaults where declared explicitly
- validation-triggered escalation
- receipt generation

**Leave as policy / human tuning:**

- exact threshold tuning over time
- when a task class should permanently move bands
- whether degraded local fallback is allowed when remote is unavailable

---

## Non-goals

This design does not:

- define the full orchestrator IR
- define CTB compilation
- replace the skill system
- make trigger words the routing engine
- require every model backend to be a package in v1
- solve multi-agent distributed scheduling

---

## File Changes

**Create:**

- `docs/alpha/agent-runtime/HYBRID-LLM-ROUTING.md` (this file)

**Edit (future, when implementation begins):**

- `docs/alpha/agent-runtime/RUNTIME-CONTRACT-v2.md` — add body-level provider/routing visibility
- `docs/alpha/agent-runtime/AGENT-RUNTIME.md` — clarify provider routing as body-plane policy
- `docs/alpha/runtime-extensions/RUNTIME-EXTENSIONS.md` — clarify that LLM backends can be delivered as capability providers if desired
- Go runtime kernel: routing policy module, provider interface, receipts

---

## Acceptance Criteria

- [ ] AC1: runtime can route to local or Anthropic provider by deterministic policy
- [ ] AC2: local-default tasks use the local provider when thresholds remain inside local band
- [ ] AC3: remote-default tasks escalate directly when hard thresholds or task-class rules apply
- [ ] AC4: validation failure or explicit low-confidence from local can trigger escalation
- [ ] AC5: every model call leaves a routing receipt with decision and reason
- [ ] AC6: routing policy is visible in runtime contract/status
- [ ] AC7: no silent provider swap occurs without a receipt
- [ ] AC8: thresholds are configurable but have sane runtime defaults
- [ ] AC9: package/registry/architecture tasks route remote by default
- [ ] AC10: single-file deterministic coding tasks route local by default

---

## Known Debt

- token estimation quality will need one tuning pass after real use
- dependency depth heuristics may need refinement for monorepo/package-graph reality
- local confidence/validation signals need a concrete first implementation
- provider packaging vs built-in provider adapters can be decided later without changing the routing model

---

## Tuning Strategy

Collect routing receipts for a few weeks. Tune from actual failures and unnecessary escalations instead of trying to guess perfectly up front. That is the simplest robust first cut.

---

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Current runtime has no explicit model-routing policy even though capability routing belongs in body/runtime policy |
| 1 Select | — | — | Selected gap: model routing is currently implicit and unauditable |
| 4 Gap | this artifact | — | Named incoherence: provider choice is not yet a deterministic runtime decision |
| 5 Mode | this artifact | cdd/design, eng/evolve, eng/process-economics | L7 MCA; runtime/body architecture decision |
| 6 Artifacts | this artifact | — | Design drafted; implementation planning deferred |
