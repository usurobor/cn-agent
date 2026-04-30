---
name: review/diff-context
description: Diff and context inspection — structural closure, multi-format parity, snapshot consistency, stale-path validation, execution timeline, derivation, authority-surface conflict, module-truth audit, contract confinement, architecture leverage, process overhead, and design constraints checks.
artifact_class: skill
kata_surface: external
governing_question: How does β verify that the diff is internally consistent, closes its claimed gap, and does not violate invariants across the full diff context?
visibility: internal
parent: review
triggers:
  - review diff
  - review context
scope: task-local
inputs:
  - branch
  - diff
  - issue contract results from Phase 2a
outputs:
  - findings with severity and type (mechanical / judgment / contract)
requires:
  - Phase 1 (contract integrity) completed
  - Phase 2a (issue contract walk) completed
  - review orchestrator loaded
calls: []
---

# Review — Phase 2b: Diff and Context Inspection

**PRE-GATE: Phase 2a (issue contract walk) must be completed.**

---

## §2.1 Diff and context inspection

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

## After Phase 2b

Collect all diff/context findings with severity and type. Pass to Phase 2c (`review/architecture/SKILL.md`) for architecture check.

Return to orchestrator (`review/SKILL.md`) after all Phase 2 sub-skills complete for verdict rules and output format.
