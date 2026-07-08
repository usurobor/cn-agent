# self-coherence — cnos#630

## Gap

**Issue:** [usurobor/cnos#630](https://github.com/usurobor/cnos/issues/630) — "partial-matter in-progress wedge: give the reconciler a mechanical path out."

**Mode:** design-and-build. **Cell kind:** `implementation` (reconciler gap fix).

**The gap.** A `status:in-progress` dispatch cell whose run dies *after* producing matter (a `cycle/{N}` branch + a checkpointed draft PR, per the cnos#591 finalizer) but *before* `REVIEW-REQUEST.yml` is written had **no mechanical path out**: it could not requeue (the cnos#575 release-back-to-queue guard blocks over existing matter — the cnos#368 blind-requeue protection), could not advance (no `REVIEW-REQUEST.yml`, so `in-progress -> review` is blocked), and was never re-claimed (the dispatch selector only claims `status:todo`). Concrete evidence: cnos#614 sat wedged in `status:in-progress` for over a day across many scheduled sweeps, recoverable only by manual κ (operator) intervention. The literal wedge was `scan.go`'s `else` branch of the `propose_delta_recovery` case (the pre-fix `scan.go:252-256`), which re-evaluated to the identical no-op decision on every scan tick once matter was checkpointed.

**Round:** R0 (first implementation pass on `cycle/630`, `run_class: first_pass` per the γ scaffold — clean first dispatch, no prior claim history).

**Base:** `cycle/630`, created from `main@6143b53c9098e415c4e7e83791843bcad2313314` per `.cdd/unreleased/630/gamma-scaffold.md`.
