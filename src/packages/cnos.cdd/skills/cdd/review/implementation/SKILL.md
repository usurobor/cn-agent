---
name: review/implementation
description: Implementation review — issue contract walk, diff inspection, architecture checks, and evidence-bound findings.
artifact_class: skill
kata_surface: external
governing_question: How does β verify that the branch implementation satisfies the issue contract, diff/context invariants, and architecture constraints?
visibility: internal
parent: review
triggers:
  - review implementation
scope: task-local
inputs:
  - branch
  - issue
  - diff
  - contract integrity results from Phase 1
outputs:
  - AC coverage table
  - named doc updates table
  - CDD artifact contract table
  - active skill consistency table
  - findings with severity and type
requires:
  - Phase 1 (contract integrity) completed
calls: []
calls_dynamic:
  - source: project design constraints
---

# Review — Phase 2: Implementation

**PRE-GATE: §2.0.0 Contract Integrity must be completed (Phase 1).** Verify branch is unmerged.

---

## 2.0 Issue contract walk

Build these tables before reading the diff:

### AC Coverage
| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|

### Named Doc Updates
| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|

### CDD Artifact Contract
| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|

### Active Skill Consistency
| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|

Verify that every Tier 3 skill the issue requires (and every active design constraint named in the project) was actually loaded by α and is applied somewhere in the diff. A required skill that is "loaded" but not "applied" is silent process theater. A skill applied without being declared is a discoverability gap.

---

## 2.1 Diff and context inspection

### 2.1.1 Structural closure and input-source enumeration

**a)** When code claims "structurally prevented" or "impossible by construction," enumerate all input sources at security-level rigor and verify the closure claim against every one. An unaudited input source makes any structural-prevention claim unchecked.

- For filters/sanitizers/validators: trace every source through the full pipeline.
- For new data surfaces: verify both write paths AND read paths.

**b)** When one fallback/compatibility path is touched, audit all sibling fallback paths. Compatibility code is a bundle — changes to one path can invalidate assumptions in a sibling.

### 2.1.2 Multi-format semantic parity

If a value is rendered in multiple formats (prose, JSON, YAML, table, code block), verify all formats say the same thing. Format divergence is a common source of stale paths.

### 2.1.3 Snapshot consistency

If the project uses snapshot/expect tests, diff each snapshot change against the PR description.
- ❌ Accept snapshot changes as "expected" without reading them.
- ✅ Verify each snapshot change makes sense given the diff.

### 2.1.4 Stale-path validation

When a file is moved, renamed, or deleted:
- Grep the old path across the full repo (excluding historical audit docs).
- Check live consumers, imports, README links, docs, and source maps.

### 2.1.5 Branch naming and conventions

Verify branch name follows project convention. Check for duplicate list entries, convention violations, or inconsistencies introduced by the diff.

### 2.1.6 Execution timeline for state-changing paths

If code crosses process or binary boundaries, trace which process owns which values at each step.
- ❌ "This reads the version in process, so it's current"
- ✅ "Old process still holds old constants after replacement; re-exec before reading"

### 2.1.7 Derivation vs validation

If an AC claims "single source of truth" or "one-file edit," verify that a generation step exists, not only a check.

### 2.1.8 Authority-surface conflict

When multiple surfaces claim to define the same thing, verify they agree. Common conflict sites:
- canonical doc vs executable skill
- runtime contract vs machine-readable JSON
- issue ACs vs self-coherence.md
- branch summary vs actual branch state

### 2.1.9 Module-truth audit for model-correctness changes

When a change is about correctness of a model, do not limit review to the diff. Scan the full touched module for other assumptions of the same kind.
- Ask: "What else in this module assumes the same thing the diff is fixing?"

### 2.1.10 Contract-implementation confinement check

When a function's contract claims a restricted input domain, verify the implementation actually **rejects** inputs outside that domain. Silent acceptance of unclaimed inputs is a confinement gap.

### 2.1.11 Architecture leverage check

Ask whether the diff merely improves the current architecture or misses a higher-leverage boundary move. Is repeated pressure being treated as a one-off patch?

### 2.1.12 Process overhead check

If the diff adds new docs, artifacts, gates, or procedures: What exact failure does this prevent? Who uses the new artifact? Could this be automated instead?

### 2.1.13 Project design constraints check

If the project maintains a design constraints document, load it. Verify the change preserves, tightens, or explicitly revises each affected constraint. Silent violation is a D finding. Silent revision is a C finding.

---

## 2.2 Architecture and design check

**Activate when the change touches:** package boundaries, command/provider/orchestrator/skill separation, source/artifact/installed flow, registry design, kernel vs package responsibility, transport vs protocol semantics, command dispatch vs domain logic.

**Load:** `src/packages/cnos.core/skills/design/SKILL.md` when active.

| Check | Question |
|---|---|
| A. Reason to change | Does each touched module still have one real reason to change? |
| B. Policy above detail | Does policy remain in the kernel/core? |
| C. Truthful interface | Does each interface promise only what all implementations support? |
| D. Registry normalization | Do different source forms normalize into one runtime descriptor? |
| E. Source/artifact/installed | Is it still clear what is authored, built, and installed? |
| F. Surface separation | Are skills, commands, orchestrators, providers distinct? |
| G. Degraded-path visibility | Is fallback behavior visible and testable? |

Output block (include when active):

```markdown
## Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes / no / n/a | |
| Policy above detail preserved | yes / no / n/a | |
| Interfaces remain truthful | yes / no / n/a | |
| Registry model remains unified | yes / no / n/a | |
| Source/artifact/installed boundary preserved | yes / no / n/a | |
| Runtime surfaces remain distinct | yes / no / n/a | |
| Degraded paths visible and testable | yes / no / n/a | |
```

If any Architecture Check row is "no," the review cannot approve without an explicit issue-backed redesign or scope reduction.

Rules:
- Silent architectural boundary smear is a blocking finding.
- Silent source-of-truth duplication is a blocking finding.
- Intentional constraint revision must be named explicitly.

---

## After Phase 2

Collect all findings with severity and type. Return to orchestrator (`review/SKILL.md`) for verdict rules and output format.
