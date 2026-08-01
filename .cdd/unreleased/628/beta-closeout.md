# β close-out — cnos#628

**Issue:** #628 — S1 (doctrine): supersede cell-kind ontology with the CCNF-kernel + WC/PC/CC
deployment classes; settle CM↔V.
**cell_kind:** `doctrine`
**run_class:** `repair_pass` (R1).
**Rounds:** R0 (`verdict: iterate`, one dispositive blocker) → R1 (`verdict: converge`).

## Round summary

**R0** (`.cdd/unreleased/628/beta-review.md §R0`): independent review of PR #629 against AC1–AC8.
Found the substantive doctrine content independently sound on the merits, but the branch was not
actually rebased onto current `main` as its own body claimed, and silently deleted landed
cnos#639 content (`delta/SKILL.md` §9.13; `.cdd/unreleased/639/*`, 8 files; 4 dispatch-substrate
files) undisclosed in the PR's own "Changes" summary. Dispositive — blocked convergence. Findings
posted as a plain PR comment on #629 (GitHub refused a formal "request changes" review since PR
#629 and this dispatch wake execute under the same `sigma` account — a live substrate gap in the
α≠β firebreak, flagged for operator attention as a distinct follow-up, not this cell's own scope).
Cell moved to `status:blocked` (external dependency) — this cell cannot self-repair PR #629
without collapsing the independence it exists to preserve, and has no in-band α to re-dispatch.

**R1** (`.cdd/unreleased/628/beta-review.md §R1`, `repair_pass` per
`cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step A check 3): κ rebased PR #629 (new head
`68797cf9`) and re-applied `status:todo`. This round independently re-verified the rebase —
`git merge-base --is-ancestor` now succeeds; the three-dot diff against current `main` shows
exactly the 5 originally-intended doctrine files with the same `+149/-20` net change; all
previously-deleted content (§9.13, `.cdd/unreleased/639/*`, the 4 dispatch-substrate files) is
confirmed present and unchanged. AC1–AC8 substance, already sound at R0, is confirmed unchanged
via spot-checks against the rebased tip (not re-derived from scratch, per this cell's own
"reviews the artifact, does not re-derive the doctrine" instruction and `REPAIR-PLAN.md`'s scope
discipline). **Verdict: converge.**

## Findings disposition

| Finding | Round | Disposition |
|---|---|---|
| Stale base / silent deletion of cnos#639 content | R0 | **Resolved** — κ's rebase (68797cf9), independently re-verified at R1. Not a δ override; the finding was fully addressed, not waived. |
| `human_gate` / `waiting_human_gate` terminology overlap in `CELL-RUNTIME.md`'s runner-routing enums | R0 | **Carried forward, non-blocking.** Does not block AC6; named again at R1 as still open. Candidate follow-up tidy for the operator, not this cell's scope (AC6 only requires the invariants stated, not enum-name consistency). |
| AC7 (#366 reconciliation) verified only as far as the PR body's own text, not independently re-derivable | R0 | **Carried forward, non-blocking.** Unaffected by the R1 rebase; still "satisfied as far as diff-checkable." |
| AC8 vocabulary nit ("constitutive change"/"migration" wording lives in the doctrine file, not the PR body proper) | R0 | **Carried forward, non-blocking.** Substance present; nit stands, does not block. |
| Same-account α/β firebreak gap (GitHub cannot structurally distinguish κ-as-PR-author from this wake's independent β) | R0 | **Out of this cell's scope — flagged for operator.** Procedurally (not structurally) enforced today; a separate bot identity for dispatch-wake execution is the named candidate mechanism. Not fixed by this cell (non-goal: no wakes, no CDS behavior change). |

No finding required a δ override at R1 — every blocking finding from R0 was genuinely resolved,
independently re-verified, not overridden or waived.

## Repair evidence

```yaml
repair_evidence:
  prior_rejection: "https://github.com/usurobor/cnos/issues/628#issuecomment (status:blocked, 2026-07-09T08:16:10Z) + .cdd/unreleased/628/beta-review.md §R0 (verdict: iterate)"
  repairs_required:
    - headline_finding: "PR #629 branch not rebased onto current main; silently deletes cnos#639 content (delta/SKILL.md §9.13, .cdd/unreleased/639/* [8 files], 4 dispatch-substrate files) with no disclosure in the PR's own Changes summary"
  repairs_completed:
    - headline_finding: "κ rebased sigma/cell-runtime-arch-note onto current main (new head 68797cf9, parent 44aa9f84). Independently re-verified at R1 (beta-review.md §R1): git merge-base --is-ancestor now succeeds; three-dot diff vs current main shows exactly the 5 intended doctrine files (+149/-20, matching R0's reviewed scope); delta/SKILL.md §9.13, all 8 .cdd/unreleased/639/* files, and all 4 dispatch-substrate files confirmed present and unchanged on the rebased branch tip."
  repairs_not_completed: []
  delta_overrides: []
  new_state_differs_from_rejected: "PR #629 head SHA changed from the R0-reviewed (non-rebased, ancestor-check-failing) state to 68797cf9 (rebased, ancestor-check-passing, zero false deletions) — verified via git merge-base --is-ancestor and git diff --stat three-dot comparison, not by trusting the PR body or the operator's comment."
```

Per `cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step E, this block is complete (no open
`repairs_not_completed` entries, no `delta_overrides`), so this cycle may advance to
`status:review`.
