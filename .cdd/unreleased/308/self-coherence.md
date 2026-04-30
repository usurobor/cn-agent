---
cycle: 308
role: α
issue: "skill(cdd/review): split phase 2 into review/issue-contract, review/diff-context, review/architecture siblings"
---

# Self-Coherence — Cycle #308

## §Gap

**Issue:** #308 — skill(cdd/review): split phase 2 into review/issue-contract, review/diff-context, review/architecture siblings

**Version/mode:** content-preserving split; no content rewriting; CTB v0.1 frontmatter required on each new sub-skill.

**Gap being closed:** `review/implementation/SKILL.md` bundles three cognitive modes (issue-contract walk §2.0, diff/context inspection §2.1.1–§2.1.13, architecture/design check §2.2) in one file. The orchestrator (`review/SKILL.md`) already names the deferred-split as the correct decomposition. This cycle lands that decomposition: three new siblings each owning one reason to change, the monolith deleted, orchestrator updated, M2-review kata reference rewired.

**Invariant:** One reason to change per skill (design-skill discipline). Each mode changes for a different reason — issue-contract walks change when CDD artifact rules change; diff/context heuristics change when implementation review patterns change; architecture checks change when design boundaries shift.

---

## §Skills

**Tier 1:**
- `CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2:**
- `src/packages/cnos.eng/skills/eng/` (always-applicable engineering skills) — not loaded explicitly; no language-specific implementation work in this cycle (pure skill authoring / YAML/Markdown)

**Tier 3 (issue-declared):**
- `cnos.core/skills/design/SKILL.md` — split design with one reason to change per skill; governs the boundary call for this decomposition
- `cnos.cdd/skills/cdd/review/SKILL.md` — parent skill being modified; loaded to verify the orchestrator's phase contract stays intact

**Applied constraints from Tier 3:**
- design/SKILL.md §2.1: decompose by reason to change → each new sub-skill owns one cognitive mode
- design/SKILL.md §3.2: one source of truth per fact → orchestrator owns phase order + verdict; sub-skills own their cognitive mode
- design/SKILL.md §2.4: truthful interfaces → each sub-skill's `inputs`/`outputs` declares exactly what it consumes/produces
- CTB v0.1 `LANGUAGE-SPEC.md` §2: hard-gate frontmatter fields on every new sub-skill — `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, `triggers`, `scope`

---

## §ACs

### AC1: Three new sub-skills exist with valid CTB v0.1 frontmatter

**Status: MET**

Evidence:
- `review/issue-contract/SKILL.md` created; frontmatter check: name ✓, description ✓, artifact_class ✓, kata_surface ✓, governing_question ✓, triggers ✓, scope ✓
- `review/diff-context/SKILL.md` created; frontmatter check: all 7 hard-gate fields present ✓
- `review/architecture/SKILL.md` created; frontmatter check: all 7 hard-gate fields present ✓

Python frontmatter parse confirmed all three pass at commit `d27db7b1`.

### AC2: review/implementation/SKILL.md deleted

**Status: MET**

Evidence:
```
find src/packages/cnos.cdd/skills/cdd/review -type f -name "SKILL.md" | sort
# Result:
# src/packages/cnos.cdd/skills/cdd/review/SKILL.md
# src/packages/cnos.cdd/skills/cdd/review/architecture/SKILL.md
# src/packages/cnos.cdd/skills/cdd/review/contract/SKILL.md
# src/packages/cnos.cdd/skills/cdd/review/diff-context/SKILL.md
# src/packages/cnos.cdd/skills/cdd/review/issue-contract/SKILL.md
```
Exactly 5 files. `implementation/SKILL.md` absent. Deleted via `git rm` in commit `d9331f19`.

### AC3: Orchestrator updated

**Status: MET**

Evidence:
- `calls` array now enumerates: `review/contract/SKILL.md`, `review/issue-contract/SKILL.md`, `review/diff-context/SKILL.md`, `review/architecture/SKILL.md`
- Phase 2 body replaced: single "Load: review/implementation/SKILL.md" with three sequential sub-skill loads in original §2.0 → §2.1 → §2.2 order
- Deferred-split note removed. Verified: `grep "Future split|deferred-split|grows past" review/SKILL.md` → 0 hits

### AC4: Each sub-skill is independently loadable

**Status: MET**

Evidence:
- `review/issue-contract/SKILL.md`: 63 lines (well under 250 target); self-contained with own `inputs`, `outputs`, `requires`; no cross-references to sibling sub-skills' internal state
- `review/diff-context/SKILL.md`: 108 lines; self-contained; `inputs` declares dependency on Phase 2a results
- `review/architecture/SKILL.md`: 77 lines; self-contained; `inputs` declares dependency on Phase 2b results; `calls: src/packages/cnos.core/skills/design/SKILL.md` declared for the design skill it loads

All three sub-skills declare their inputs and do not rely on undeclared shared state.

### AC5: M2-review kata cross-references updated

**Status: MET**

Evidence:
- `grep -rn "review/implementation/SKILL.md" src/packages/cnos.cdd.kata/` → 0 matches
- `kata.md` line 16: architecture check now points at `review/architecture/SKILL.md §2.2` (commit `690caf14`)

### AC6: Phase order in orchestrator unchanged

**Status: MET**

Evidence:
- `grep -n "Phase 1|Phase 2|Phase 3" review/SKILL.md` → Phase 1 (contract integrity) L51, Phase 2 (implementation review) L59, Phase 3 (verdict) L69
- Sub-skill load order within Phase 2: issue-contract → diff-context → architecture matches original §2.0 → §2.1 → §2.2 sequence

### Stale-path audit (beyond issue ACs)

Peer enumeration: all live surfaces referencing `review/implementation/SKILL.md` were searched and addressed:
- `review/contract/SKILL.md` L194: "After Phase 1" footer — updated to reference `review/issue-contract/SKILL.md` (commit `7c97d257`)
- `schemas/skill.cue` L25: illustrative comment — updated to use `review/issue-contract` as example (commit `7c97d257`)
- `M2-review/kata.md` L16: architecture check reference — updated (commit `690caf14`)
- Historical release artifacts (`.cdd/releases/3.63.0/301/`) — exempt; audit history, not live code

---

## §Self-check

**Did α push ambiguity onto β?**

No. Every AC has concrete evidence (grep results, commit SHAs, line counts). The diff is content-preserving: cognitive content moved without rewriting. The three reasons to change per mode are named and match the issue's stated rationale.

**Is every claim backed by evidence in the diff?**

Yes:
- AC1 frontmatter presence: Python check run and confirmed
- AC2 deletion: `git rm` confirmed; find oracle confirmed
- AC3 orchestrator: calls array verified against SKILL.md frontmatter; deferred-split note grep confirmed 0 hits
- AC4 independence: line counts verified; inputs/outputs fields present and correct
- AC5 kata refs: grep confirmed 0 stale refs
- AC6 phase order: grep confirmed Phase 1 → 2 → 3 ordering intact

**Scope violations?** None. Content was not rewritten. Phase ordering not changed. Verdict rules untouched. CTB v0.2 not promoted. No per-mode katas added.

**Design skill applied?** Yes. Each new sub-skill owns one cognitive mode (one reason to change). Orchestrator owns phase order and verdict format (source of truth separation). Sub-skills' inputs/outputs are truthful about what they consume/produce.

---

## §Debt

**Known debt: None.**

The `schemas/skill.cue` comment was a documentation example using `review/implementation` as an illustrative name. Updated to `review/issue-contract` in the same commit. No functional schema change required (the regex already admits slash-compound names; new names all conform).

The `.cdd/releases/3.63.0/301/` artifacts retain `review/implementation` references — this is correct: they are immutable historical records of the #301 cycle.

---

## §CDD-Trace

| Step | Evidence |
|------|---------|
| 1. Receive | Dispatched on issue #308, `cycle/308` branch, identity `alpha@cdd.cnos` |
| 2. Produce | Artifact order: design (not required — structural split, no design artifact needed), plan (not required — single-class split, clear sequencing), tests (not required — skill files have no executable test harness; correctness proved by AC oracles), code/content (3 new sub-skills + orchestrator update + deletion + kata rewire), docs (self-coherence), pre-review |
| 3. Prove | AC1–AC6 each mapped to concrete oracle evidence above |
| 4. Gate | Pre-review gate checklist below |
| 5. Review loop | Pending |
| 6. Close-out | Pending |

**Design not required:** This is a structural relocation. The issue itself names the exact decomposition (issue-contract / diff-context / architecture), the content origin (§2.0 / §2.1 / §2.2 verbatim), and the boundary rationale (one reason to change per mode). No novel design decision is required — the decision was made in the orchestrator's deferred-split note (PR #304).

**Plan not required:** Five file operations in a deterministic order (create 3, update 1, delete 1, update 2 cross-refs). No non-trivial sequencing.

**Tests not required:** Skill files are SKILL.md artifacts. There is no executable test harness for skill content. Correctness is proved by AC oracle evidence (frontmatter parse, find oracle, grep, line counts). The existing CI job (I5, `tools/validate-skill-frontmatter.sh` via `schemas/skill.cue`) validates frontmatter conformance on changed files.
