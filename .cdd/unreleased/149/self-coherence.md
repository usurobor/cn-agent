## §Gap

**Issue:** #149 — UIE must include skill loading as part of Understand before Execute
**Version/mode:** MCA (small-change path — single-file prose addition, §1.2)

SOUL.md §2.1 Observation had no explicit requirement to identify and load required skills before acting. An agent could complete the Observation step and proceed to Action using stale memory of a skill rather than the current skill file. Three concrete failures on 2026-04-02 demonstrated the gap: heartbeat without inbox skill (escalation leak), release without release skill (skipped CDD bookends), review without review skill (operator had to prompt mid-task).

The fix: add a **skill-loading gate** to §2.1 Observation, adjacent to the falsification gate, that makes skill loading an explicit Observation sub-requirement with visible-load convention and ❌/✅ examples.

---

## §Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 3 (issue-named):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — meta-skill discipline: imperative form, ❌/✅ pairs, governing question
- `src/packages/cnos.core/skills/write/SKILL.md` — house style for operator-facing load-bearing prose

---

## §ACs

### AC1: Skill-loading rule in SOUL.md §2.1 Observation

**Invariant:** An agent reading SOUL.md §2 sees the skill-loading requirement before reaching §2.2 Action.
**Oracle:** §2.1 contains a sub-rule stating "Before acting, identify which skills this action requires and load them."
**Positive:** Sub-rule explicitly names that memory of a skill ≠ loading the skill.
**Negative:** A rule that only says "consider loading skills" does not satisfy the AC.

**Evidence:** SOUL.md §2.1 now contains:

> **Skill-loading gate:** Before acting, identify which skills this action requires and load them. Memory of a skill is not the same as loading it — skill files are the constraint, and their content drifts from agent memory over time. Name the loaded skills before proceeding to Action.

The rule is imperative ("identify … and load them"), appears inside §2.1 before §2.2, and explicitly states "Memory of a skill is not the same as loading it." AC1: **met**.

### AC2: ❌/✅ examples for skill loading

**Invariant:** The rule shows what does and does not satisfy it.
**Oracle:** At least one ❌/✅ pair shows acting on memory of a skill versus loading the skill file.

**Evidence:** Added immediately after the skill-loading gate:

```
- ❌ Act on memory of what the release skill says.
- ✅ Load `release/SKILL.md`, name it as loaded, then act.
```

Same situation (release skill), contrastive (memory vs file), imperative. AC2: **met**.

### AC3: Visible-load convention named

**Invariant:** When the agent loads skills, the load is observable to the operator.
**Oracle:** Sub-rule states the agent names loaded skills before proceeding to Action.
**Negative:** Silent load does not satisfy this AC.

**Evidence:** The skill-loading gate includes: "Name the loaded skills before proceeding to Action." AC3: **met**.

---

## §Self-check

**Did α push ambiguity onto β?** No. All three ACs have direct evidence quotes from the diff. The rule is imperative, the examples are contrastive, and the visible-load convention is explicit. No judgment is deferred.

**Is every claim backed by evidence in the diff?** Yes. Each AC cites the verbatim added text from SOUL.md §2.1.

**Non-goals verified unviolated:**
- No auto-detection logic added.
- No loading-all-skills-on-every-action logic added.
- No changes outside §2.1 — §2.2 through §3.6 are untouched.
- No changes to USER.md.
- No linter/runtime enforcement added.
- Sigma's derived SOUL not touched.

**Consistency with #277:** The skill-loading gate is added as a named block adjacent to the falsification gate inside §2.1. It does not anchor to any section number. If #277 renumbers sections, the prose survives without contradiction — the gate is paragraph-level, not number-anchored.

---

## §Debt

None. The change is narrowly scoped per the issue's explicit non-goals. No deferred items.

---

## §CDD-Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Observation inputs read; selected signal: SOUL.md §2.1 missing explicit skill-loading requirement |
| 1 Select | — | — | Selected gap: #149 — UIE must include skill loading as part of Understand before Execute |
| 2 Branch | `cycle/149` | cdd | Branch verified on `origin/cycle/149`; base `origin/main` |
| 3 Bootstrap | — | cdd | Small-change path (§1.2) — no version directory required |
| 4 Gap | `.cdd/unreleased/149/self-coherence.md` §Gap | — | Named incoherence: SOUL §2.1 lacked explicit skill-loading gate |
| 5 Mode | `.cdd/unreleased/149/self-coherence.md` §Skills | cdd, skill, write | MCA (small-change). Tier 3: `cnos.core/skills/skill`, `cnos.core/skills/write` |
| 6 Artifacts | `src/packages/cnos.core/templates/SOUL.md` §2.1 diff | skill, write | Design: not required (single-file prose addition, explicit by issue). Plan: not required (single-file, single-section). Tests: not applicable (prose rule, no executable surface). Code: `SOUL.md` §2.1 edited — skill-loading gate added. Docs: SOUL.md is the doc. |
| 7 Self-coherence | `.cdd/unreleased/149/self-coherence.md` | cdd | AC-by-AC self-check completed; all three ACs met with direct evidence quotes |
| 7a Pre-review | `.cdd/unreleased/149/self-coherence.md` | cdd | Pre-review gate passed (see §Review-readiness below) |

---

## Review-readiness | round 1 | implementation SHA: 56438d40 | base SHA: 3cef674f (origin/main at 2026-04-30T15:50 UTC) | branch CI: not available for branch-push CI (CI runs on main pushes only; β to verify on merge) | ready for β

**Pre-review gate rows (alpha/SKILL.md §2.6):**

| Row | Status | Notes |
|-----|--------|-------|
| 1. Branch rebased onto current origin/main | ✅ | base SHA `3cef674f` = origin/main at observation time |
| 2. Self-coherence carries CDD Trace through step 7 | ✅ | §CDD-Trace rows 0–7 present |
| 3. Tests present or explicit reason none apply | ✅ | Prose rule change — no executable surface |
| 4. Every AC has evidence | ✅ | AC1–AC3 all cite verbatim diff text |
| 5. Known debt explicit | ✅ | §Debt: "None" |
| 6. Schema/shape audit | ✅ | Not applicable — prose-only change |
| 7. Peer enumeration | ✅ | Issue explicitly scopes to §2.1 only; USER.md and Sigma SOUL explicitly out of scope per issue |
| 8. Harness audit | ✅ | Not applicable — no schema-bearing contract changed |
| 9. Post-patch re-audit | ✅ | No mid-cycle patches |
| 10. Branch CI green | ⚠️ | CI runs on main pushes only; no run available for branch — β verifies before merge |
| 11. Commit author email | ✅ | `alpha@cdd.cnos` confirmed via `git log -1 --format='%ae' HEAD` |
