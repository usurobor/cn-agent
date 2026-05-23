# γ scaffold — Cycle 421

**Cycle:** [cnos#421](https://github.com/usurobor/cnos/issues/421) — Track A1 of [cnos#405](https://github.com/usurobor/cnos/issues/405)
**Branch:** `cycle/421` (from `origin/main` @ `9fc04a99` — the merge that closed cycle/420 / Sub 6 of #404 / the handoff extraction wave)
**Mode:** explore / design — survey + decision cycle, not implementation
**Collapse pattern:** γ+α+β collapsed on δ. Commits: `γ-421:`, `α-421:`, `β-421:`.
**Date:** 2026-05-23

## Intent

First dispatchable sub of #405 post-#404 closure. The cnos.handoff package landed v0.1-complete on 2026-05-22 (`9fc04a99`), satisfying the Track A gate. Track A1 is the survey + name-pick + decision cycle that A2–A6 + B1 dispatch against. This cycle produces a **single design doc** with **five pinned decisions** + a **refined sub-issue queue** + a **deferred-questions list**. No implementation, no schemas, no kernel changes.

## Surface plan (D1–D2)

### D1 — `docs/gamma/design/ccnf-o-track-a1-survey.md`

A single markdown doc, ~400–700 lines, with seven required sections per the dispatch brief:

1. **Name pick** — pin CCNF-O or CCNF-X. Operator's mild preference is CCNF-O; the survey confirms or overrides on a 1-paragraph rationale.
2. **Surface inventory matrix** — table with ≥ 15 rows × 5 columns covering every surface in #405 §13 and the five surfaces #404 landed in `cnos.handoff`. For each surface: current home / universal CCNF-O? / realization-specific? / rationale.
3. **Higher-level form classification** — SimpleCycle / SubstantialCycle / Wave / Roadmap / CoherenceMeasuredCell / AutonomousLoop. Universal vs realization-derived, with rationale.
4. **TSC integration v0.1 scope** — confirm Track B1 can design in parallel with Track A2+. Name design touchpoints + dependencies.
5. **Package location** — pin one of three options (cnos.cdd/orchestration/ rejected; new cnos.{ccnf-o|orchestration|grammar|…}; fold into cnos.handoff rejected per #405 §15). Pin a candidate name if option 2.
6. **Sub-issue queue (refined)** — paragraph per Track A2–A6 + Track B1 (≥ 6 paragraphs) naming what each dispatch brief will pin.
7. **Open questions deferred** — anything Track A1 can't resolve; flagged for Track A2's dispatch brief.

### D2 — `docs/gamma/design/README.md` (new directory)

`docs/gamma/design/` does not yet exist. Create the directory with a minimal README carrying a Document Map row for the survey doc. Do NOT put the survey in `docs/gamma/essays/` — this is a decision-survey artifact (design-class), not an essay (position-paper-class). The categorical match is correct as authored. The `docs/gamma/essays/README.md` is left untouched.

## Implementation contract

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None |
| Package scoping | New doc at `docs/gamma/design/ccnf-o-track-a1-survey.md`; new directory README at `docs/gamma/design/README.md` |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | N/A — Track A2+ types the wire formats CCNF-O surfaces; this cycle does not author schemas |
| Backward compat | Hard rules below; no src/packages/ diff, no schemas/ diff, no CCNF kernel diff, no cnos.handoff diff |

## Acceptance criteria

AC1–AC11 per [cnos#421](https://github.com/usurobor/cnos/issues/421). All mechanical or read-checks. Verified in `self-coherence.md` post-α.

- AC1: survey doc exists; ≥ 300 lines
- AC2: §1 names one of CCNF-O / CCNF-X with rationale
- AC3: §2 contains ≥ 15-row × 5-column matrix
- AC4: §3 covers all 6 higher-level forms with classification + rationale
- AC5: §4 states whether Track B1 can dispatch in parallel + lists touchpoints + dependencies
- AC6: §5 pins ONE of the three location options with rationale
- AC7: §6 has paragraph-level descriptions for Track A2, A3, A4, A5, A6, B1 (≥ 6 paragraphs)
- AC8: `git diff origin/main..HEAD -- src/packages/` returns 0 lines; same for schemas/ and src/go/
- AC9: CCNF kernel files (CDD.md, COHERENCE-CELL.md, COHERENCE-CELL-NORMAL-FORM.md) unchanged
- AC10: `docs/gamma/design/README.md` has a row for the survey doc
- AC11: no `schemas/ccnf-o/` directory; no `cnos.handoff/` diff; no `cnos.cdr/` or `cnos.cds/` diff

## Hard rules

1. **No implementation.** No CUE schemas. No Go types. No `cn cdd verify` changes.
2. **No CCNF kernel changes.** CDD.md / COHERENCE-CELL.md / COHERENCE-CELL-NORMAL-FORM.md untouched.
3. **No cnos.handoff edits.** Read it; don't modify it.
4. **No src/packages/ edits at all.** No cnos.cdd, cnos.cds, cnos.cdr, cnos.handoff, cnos.core, cnos.eng diff.
5. **No schemas/ edits.** `schemas/ccnf-o/` is created in Track A2, not here.
6. **No #405 body edits in this sub.** Produce decisions externally; a follow-on cycle may fold them in.

## Non-goals

- Do NOT author CUE schemas.
- Do NOT generate Go types.
- Do NOT modify `cn cdd verify`.
- Do NOT touch `src/packages/`.
- Do NOT change CCNF kernel doctrine.
- Do NOT preempt Track A2's typing work.
- Do NOT preempt Track B1's design work; this survey *enables* B1 to start, it does not do B1's work.
- Do NOT redesign cnos.handoff surfaces (cross-repo / dispatch / mid-flight / artifact-channel / receipt-stream).
- Do NOT modify the #405 tracker body.

## Operator action on close

Merge `cycle/421` to `main` with `--no-ff` and message `Merge cycle/421: Track A1 of #405 — CCNF-O survey + name-pick + sub-issue queue. Closes #421.`. #405 stays open; Track A2 + Track B1 are dispatchable thereafter.
