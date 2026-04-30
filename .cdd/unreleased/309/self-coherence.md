---
cycle: 309
role: alpha
---

# Self-Coherence — Cycle 309

## §Gap

**Issue:** #309 — skill(eng): create eng/troubleshoot — live diagnosis skill for environmental and runtime failures

**Version / mode:** cnos · CDD triadic cycle · α implementation

**Gap:** No live troubleshooting skill existed. When an agent hits a live failure mid-task, diagnosis was improvised. The 2026-04-30 identity-rotation dispatch test produced five β dispatch failures across three root-cause classes (OOM kill, `gh` GraphQL error, background-process lifecycle kill), each diagnosed ad-hoc with a wrong first hypothesis. A structured skill for live diagnosis would have forced process state and kernel log checks before model/token speculation.

**Target artifact:** `src/packages/cnos.eng/skills/eng/troubleshoot/SKILL.md`

**Not in scope:** Automated diagnostic tooling, platform-specific runbooks, changes to `eng/rca`, new commands, dispatch behavior changes.

## §Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract
- `CDD.md` — canonical lifecycle

**Tier 2:**
- `src/packages/cnos.core/skills/write/SKILL.md` — writing standard (top-down, point first, no clutter)

**Tier 3 (issue-specified):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — classification, formula, kata surface requirements
- `src/packages/cnos.core/skills/write/SKILL.md` — governing question, levels of structure, direct rules
- `src/packages/cnos.core/skills/design/SKILL.md` — boundary discipline; skill vs runbook vs reference; surface separation
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — invariant, oracle, positive/negative proof model
- `src/packages/cnos.eng/skills/eng/rca/SKILL.md` — adjacent sibling boundary; RCA handoff definition

All five Tier 3 skills were loaded and read before drafting began.

## §ACs

### AC1: Artifact classification is explicit

**Evidence:** SKILL.md frontmatter declares `artifact_class: skill`. §1.0 of the skill body explicitly classifies the artifact as a skill (not runbook, not reference) with justification: it teaches judgment rather than platform-specific command sequences.

The skill includes:
- frontmatter: `name`, `description`, `governing_question`, `triggers`, `scope`, `artifact_class`, `kata_surface`, `inputs`, `outputs`
- domain-specific coherence formula (§Core Principle)
- one governing question (stated in frontmatter and §1.3)
- named failure mode: **premature hypothesis** (§1.3)

**Met:** fully.

### AC2: Required skills loaded before writing

**Evidence:** All five Tier 3 skills were loaded before drafting. See §Skills above. The α session loaded them in this conversation before the SKILL.md was written.

Loaded: skill ✓, write ✓, design ✓, test ✓, rca ✓.

**Met:** fully.

### AC3: External troubleshooting practices incorporated

**Evidence:** All four sources appear in the skill body (adapted, not copied wholesale):

- **Google SRE effective troubleshooting (hypothesis-driven):** §2.4 hypothesis discipline paragraph; §2.3 triage order footnote
- **IBM problem-description questions (what/where/when/conditions/reproducibility/error messages/effect/environment):** §2.2 full problem-description field list
- **Red Hat one-change / retest discipline:** §2.4 paragraph; §2.5; §2.6; Rules 3.4, 3.6
- **CompTIA sequence (identify/theory/test/plan/implement-escalate/verify/document):** §2.3 triage order footnote

Each source is credited inline; none is reproduced wholesale.

**Met:** fully.

### AC4: Live triage algorithm is explicit

**Evidence:** §2.3 lists the eleven-step sequence from the AC (mapped to six triage classes):

| AC step | Skill location |
|---------|---------------|
| 1. Preserve evidence | §2.1 |
| 2. Describe symptom | §2.2 |
| 3. Check process state | §2.3 step 1 |
| 4. Check external kill / kernel log | §2.3 step 2 |
| 5. Check command/tool output | §2.3 step 3 |
| 6. Check resources | §2.3 step 4 |
| 7. Check lifecycle / parent process | §2.3 step 5 |
| 8. Check application/model behavior | §2.3 step 6 |
| 9. Apply one controlled fix | §2.5 |
| 10. Verify against original symptom | §2.6 |
| 11. Decide RCA handoff | §2.7 |

**Met:** fully.

### AC5: Hypothesis testing stated

**Evidence:** §2.4 lists all five required elements:
- Cheapest discriminating test first: §2.4 step 3; §3.3
- One variable changed at a time: §2.5; §3.4
- Hypothesis stated before test: §2.4 step 1; §3.2
- Oracle stated before test: §2.4 step 2; §3.2
- Negative result recorded: §2.4 step 4; §3.5
- Fix verified against original symptom: §2.6; §3.6

**Met:** fully.

### AC6: Worked examples exist

**Evidence:** §4 contains three examples from the dispatch test:

- **Example 1 — OOM kill:** symptom (silent exit) → wrong hypothesis (token limit) → triage path (process state → dmesg OOM) → root cause → MCA (add swap) → verification (re-ran β review)
- **Example 2 — `gh` GraphQL error:** symptom (9-event die) → wrong hypothesis (rate limit) → triage path (tool output error message) → root cause → MCA (`--json` flag) → verification
- **Example 3 — Background process kill:** symptom (silent exit, no OOM) → wrong hypothesis (rate limit/token) → triage path (process → dmesg → tool → resources → lifecycle) → root cause → MCA (foreground only) → verification

Each example includes all six required elements: symptom, wrong first hypothesis, better test order, root cause, MCA, verification.

**Met:** fully.

### AC7: RCA handoff defined

**Evidence:** §2.7 defines five handoff triggers:
- Recurrence risk remains
- Multi-step incident affected dispatch/release/runtime
- Root cause is systemic
- Prevention action needed
- Evidence should be preserved for future agents

Rule §3.7 reinforces: do not start RCA during live diagnosis.

**Met:** fully.

### AC8: Kata surface exists

**Evidence:** §5 contains an embedded kata with all required elements:
- Scenario ✓
- Symptoms ✓
- Available evidence ✓
- Wrong tempting hypothesis ✓
- Expected triage path ✓
- Expected root cause ✓
- MCA + Verification ✓
- Common failures ✓
- Reflection ✓

**Met:** fully.

### AC9: Frontmatter validates

**Evidence:** `cue` CLI is not installed in this environment; `tools/validate-skill-frontmatter.sh` returned exit code 2 (prerequisite missing). Manual validation against `schemas/skill.cue`:

- Hard-gate fields present: `name` (`troubleshoot`, matches `^[a-z][a-z0-9_/-]*$`), `description` (non-empty), `governing_question` (non-empty), `triggers` (list of strings), `scope` (`task-local` — valid enum)
- Spec-required-but-exception-backed fields present: `artifact_class` (`skill` — valid enum), `kata_surface` (`embedded` — valid enum), `inputs` (list), `outputs` (list)
- No `calls` entries (no filesystem path resolution needed)
- No unknown field type errors

No exception entry needed; all required fields are present.

**Met:** fully (cue CLI not available; manual validation against schema confirms shape).

## §Self-check

**Did α push ambiguity onto β?**

No. Every AC maps to concrete evidence in the diff. The single artifact (SKILL.md) is the scope; all ACs are fully addressed. The frontmatter validator could not run (`cue` missing) — this is disclosed in AC9 with manual schema confirmation and noted below in §Debt.

**Is every claim backed by evidence in the diff?**

Yes:
- AC1 (classification): frontmatter `artifact_class: skill` + §1.0 body — visible in diff
- AC2 (skills loaded): disclosed in §Skills; five skills loaded before drafting in this session
- AC3 (external practices): §2.2, §2.3, §2.4, §2.6 with inline citations — visible in diff
- AC4 (triage algorithm): §2.1–§2.7 — visible in diff
- AC5 (hypothesis discipline): §2.4, §2.5, §2.6, §3.2–§3.6 — visible in diff
- AC6 (worked examples): §4 Examples 1–3 — visible in diff
- AC7 (RCA handoff): §2.7, §3.7 — visible in diff
- AC8 (kata): §5 — visible in diff
- AC9 (frontmatter): manual validation disclosed; no exception entry required

**Peer enumeration:** The diff is a single new file. No sibling skill was modified. No peer update is needed.

**Harness audit:** No schema-bearing contract changed. No harness audit required.
