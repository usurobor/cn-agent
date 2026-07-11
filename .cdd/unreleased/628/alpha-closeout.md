# α close-out — cnos#628

**Issue:** #628 — S1 (doctrine): supersede cell-kind ontology with the CCNF-kernel + WC/PC/CC
deployment classes; settle CM↔V.
**cell_kind:** `doctrine`
**run_class:** `repair_pass` (R1; R0 was `first_pass`).

## No in-band α round

Per `gamma-scaffold.md` "Shape of this cell," this cell has no in-band α to dispatch: the
α-matter under review is externally authored as PR #629 (`sigma/cell-runtime-arch-note` →
`main`, written by κ/Sigma-as-author outside the dispatch flow), precisely because κ cannot
author this cycle's own independent β-review without violating the α≠β firebreak
(`COHERENCE-CELL.md` §"β Independence"). This document stands in place of a true
`alpha-closeout.md` to keep the cycle's canonical six-artifact set complete; it is authored by
γ/δ, not by an implementer role, and asserts nothing about AC correctness — that assertion is
β's (`beta-review.md`).

## R0 → R1 transition (what changed in the reviewed matter)

R0's independent β found one dispositive blocker: PR #629's branch was not actually rebased onto
current `main` and silently deleted landed doctrine (`delta/SKILL.md` §9.13, cnos#639;
`.cdd/unreleased/639/*`; 4 dispatch-substrate files). This cell itself did not — and could not —
repair PR #629 directly (rebasing it via this wake's own hand would collapse the α/β independence
the cell exists to preserve; see the R0 `status:blocked` comment). κ performed the rebase
externally (new PR #629 head `68797cf9`, parent `44aa9f84`) and re-applied `status:todo` to
re-trigger this wake's claim. No commit on `cycle/628` itself modifies PR #629's branch at any
point — the "repair" this cycle performs is entirely re-verification, per `REPAIR-PLAN.md`.

## Self-check

Not applicable in the conventional α sense (no implementation diff to self-assess). The
equivalent check performed by this cycle — did the repair genuinely resolve the rejected finding,
independently re-verified rather than trusted — is `beta-review.md §R1`'s job, and it is
independent of this document by construction (satisfying the same α≠β firebreak this cell's own
subject matter is about).
