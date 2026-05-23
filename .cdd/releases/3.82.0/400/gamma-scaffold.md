# γ scaffold — cycle #400

**Issue:** cnos#400 — Phase 5 of #366: γ shrink (gamma/SKILL.md → coordination + closure + triage)
**Mode:** substantial; design-and-build; γ+α+β-collapsed-on-δ (breadth-2026-05-12 precedent)
**Branch:** `cycle/400` from `origin/main`
**Date:** 2026-05-21

## Surfaces this cycle touches

Primary:
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — restructured (target ≤523 lines from 748)

Secondary (cross-reference updates + F1/F2 absorption):
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — refresh Phase 4b forward-refs (F1); fix `requires.1` CUE type mismatch (F2); update cross-references where they cite γ
- `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` — verify "called by γ" sections still name correct γ functions; potentially absorb γ-side polling-loop content (single-named-branch transition loop) if not already covered
- `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` — verify cross-refs to γ remain accurate post-shrink
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — verify cross-refs to γ remain accurate post-shrink

Tertiary (likely no edit, but checked):
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — no edit (out of scope per Phase 7)
- `INDEX.md` — INDEX.md row for cycle #400 cdd-iteration

## AC oracle approach

| AC | Oracle |
|----|--------|
| AC1 | `wc -l gamma/SKILL.md` ≤ 0.70 × 748 = 523 lines |
| AC2 | `rg "issue quality\|artifact coordination\|closure\|triage\|repair.dispatch" gamma/SKILL.md` returns hits in normative sections |
| AC3 | `rg "polling loop\|claude -p\|cn dispatch.*invocation\|CI polling\|wake mechanics\|branch preflight" gamma/SKILL.md` returns 0 hits in normative position (cross-refs OK) |
| AC4 | Cross-references in harness/delta/release-effector to γ still accurate; harness "called by γ" sections name right γ functions |
| AC5 | `self-coherence.md §Managerial-residue sweep` lists every γ verb pre-shrink, classifies KEEP / MOVE / DROP. At least one DROP entry |
| AC6 | F1: `rg "Phase 4b.*pending\|harness mechanics.*operator" delta/SKILL.md` returns 0; F2: `bash tools/validate-skill-frontmatter.sh` reports 0 findings on delta/SKILL.md |
| AC7 | No broken cross-refs; mesh remains discoverable |

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | N/A |
| Package scoping | Primary `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`; secondary cross-reference updates listed above |
| Existing-binary disposition | gamma/SKILL.md restructured; runtime supervision moves to harness/release-effector cross-references; coordination + closure + triage stays |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | All existing γ workflows continue; location of mechanics doctrine moves, not the mechanics themselves |

## Empirical anchor

The issue body and `COHERENCE-CELL.md §γ and δ Managerial-Residue Sweep` are the binding anchors. The sweep rule:

```text
If a skill says "monitor", "supervise", "oversee", or "manage",
ask what artifact, receipt, or boundary decision that verb produces.
If none, it is probably managerial residue.
If yes, rename the verb to the actual biological function.
```

Pre-cycle: gamma/SKILL.md = 748 lines (issue cites ~688; actual is 748).
Target: ≤523 lines (70%).

Phase 4a/4b/4c shipped:
- delta/SKILL.md (Phase 4a, cnos#397) — δ-role boundary
- harness/SKILL.md (Phase 4b, cnos#398) — harness substrate (dispatch invocation, observability contract, polling, identity, parallel worktrees, timeout recovery)
- release-effector/SKILL.md (Phase 4c, cnos#399) — release mechanics

γ can now shed runtime supervision and shrink to coordination + closure + triage.

## Expected diff scope

- gamma/SKILL.md: ~−250 lines (shrink polling §2.5 transition loop, dispatch shell snippets, release-prep mechanics narration, kata bloat); restructure to keep coordination/closure/triage tight
- delta/SKILL.md: ~−10 lines (refresh Phase 4b forward-refs to past-tense) + 1-line frontmatter fix (F2)
- harness/SKILL.md, release-effector/SKILL.md, operator/SKILL.md: small cross-reference touches if any γ-section anchors changed
- INDEX.md: +1 row for cycle #400 cdd-iteration

## Sections of gamma/SKILL.md (current 748 lines) — inventory plan

To be expanded in `design-notes.md` with the KEEP/MOVE/DROP classification. The 11 normative sections are:
- §1.1–§1.3 Define (parts, fit, failure)
- §2.1 Step 1a Observe
- §2.2 Step 1b Select
- §2.2a Peer enumeration
- §2.3–§2.4 Issue pack + quality gate
- §2.5 Steps 3a–5 Create branch + dispatch + unblock
- §2.6 Steps 6–7 Release prep
- §2.7 Steps 8–9 Close-out triage
- §2.8–§2.9 Steps 10–13 Cycle iteration + process gap
- §2.10 Steps 13–15 Closure gate
- §2.11 γ as autonomous coordinator
- §3 Rules
- §4 Embedded Kata (A, B, C) + Resumption
