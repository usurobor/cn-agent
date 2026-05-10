## Alpha Close-Out — cycle #338 [provisional — pending β outcome]

**Issue:** #338 — cdd: Add §1.6c — dispatch sizing, prompt scope, and commit checkpoints
**Branch:** `cycle/338`
**Mode:** docs-only (`design-and-build`; §2.5b)
**Implementation SHA:** `275bb1a4` (self-coherence §CDD-Trace commit, before review-readiness signal)

This close-out is provisional because α exits after signaling review-readiness in the sequential bounded dispatch model (§1.6). It will be finalized via re-dispatch after β approves and merges, per `alpha/SKILL.md §2.8`. This provisional status is declared as D2 in `self-coherence.md §Debt`.

---

### Summary

Three skill files patched in one cycle, adding a feedback loop that was missing from the CDD coordination model: (1) a timeout budget heuristic that scales with AC count, (2) an operator-executable recovery procedure for the case when an agent is SIGTERMed without committing, and (3) per-cycle telemetry fields in the PRA to validate and tighten the heuristic over time.

The three changes are a single feedback loop — the heuristic (§1.6c) needs the telemetry (post-release §4) to be validated, and the recovery procedure (operator §timeout-recovery) is the failure path when the heuristic is violated. They belong in one cycle.

---

### Friction log

**F1 — Recursive coherence bootstrap paradox:** This cycle's own dispatch budget (600s) is below the heuristic it proposes (900s for 5 ACs docs cycle). The mismatch is structurally unavoidable: the rule doesn't exist until this cycle ships it. β's recursive-coherence check (per the proof plan) will surface this. α recorded it in §Debt D1 and in self-coherence.md §Self-check. The heuristic is not wrong because of this — the mismatch confirms the rule is non-trivial.

**F2 — AC2 oracle case sensitivity:** The initial section heading was "Timeout Recovery" (capital R) but the AC2 oracle uses case-sensitive grep for "Timeout recovery" (lowercase r). Caught during §ACs verification, corrected immediately (changed heading to sentence case). No β round needed.

---

### Patterns observed

- **Docs-only cycles with ≥3 cross-file edits benefit from explicit commit ordering in the dispatch prompt.** The recommended order (CDD.md → operator → post-release) was correct: each file edit is independently useful, and the CDD.md change is the anchor others reference.
- **The "initial; refine with telemetry" framing is load-bearing.** The heuristic constants are guesses. Presenting them as final would be an honest-claim violation. The framing is the key to avoiding that violation.
- **Cross-ref drift risk:** four surfaces (§1.6c, §1.6b, operator §7, post-release §4) form a ring. If any one section name changes in a future cycle, three cross-refs break. The telemetry field names (`dispatch_seconds_*`, `commit_count_at_termination`) are the most stable anchors — they are concrete and appear verbatim in three places.

---

### Engineering reading

The dispatch sizing problem is a known class in agent scheduling literature: task complexity cannot be reliably predicted before execution, so the timeout must be set from a prior estimate (AC count proxy here) with empirical refinement. The `120s/ac` multiplier is conservative relative to observed cycle times but aggressive relative to the prior fixed-budget approach. The real test is whether the telemetry shows `dispatch_seconds_actual` systematically above or below `dispatch_seconds_budget` across 10+ cycles.

The operator recovery section (`§timeout-recovery`) is the more immediately useful artifact — it converts an undocumented operator memory dependency into a checklist. The prior state required the operator to know to run `git status --short && find . -newer ...`; now that knowledge is in the skill file.
