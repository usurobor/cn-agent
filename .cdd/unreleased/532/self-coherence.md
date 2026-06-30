# self-coherence — cycle/532

manifest:
  completed: [Gap]

## Gap

Issue: [#532](https://github.com/usurobor/cnos/issues/532) — "cdd/cds: require review-request proof before status:review"

A dispatch cell was able to set an issue to `status:review` without producing reviewable matter (observed in #524 W4: `status:review` + no PR + no commits beyond base + no diff + no receipt + no closeout + no STOP comment). CDD's closure sequence (`contract → matter → review → receipt → verdict → decision`) does not currently gate the `status:review` label transition on any mechanical proof that matter exists — the label is writable by the cell on a self-assertion alone.

The operator's pre-claim clarification comment (2026-06-30T19:29:09Z) is the authoritative, binding interpretation of AC4/AC5 and overrides any ambiguity in the issue body: the guard must be a mechanical state-integrity check, not prompt guidance alone; deliverable proof for the current GitHub/CDS surface means a PR (or approved no-op), commits beyond base (or no-op approval), a non-empty diff (or declared no-op), a `REVIEW-REQUEST.yml`, required closeout artifacts, and a review request that names PR/branch/head SHA/checks; on missing proof the guard must fail and the run must post STOP/BLOCKED, not leave `status:review` standing unproven.

**Version / mode:** `design-and-build` per the issue header (no stable `DESIGN.md`/`PLAN.md` existed at scoping time; γ confirmed MCA preconditions did not hold). This cycle both designs the `REVIEW-REQUEST.yml` shape (drafted as a strong proposal in the issue body, not yet a converged schema doc) and builds the guard in the same cycle, per the issue's own "Cycle scope sizing" decision (keep whole as one design-and-build issue).

**δ's binding correction to γ's scaffold (followed, not γ's mis-citation):** γ's scaffold claimed the cnos#516 guard is "invoked from inside `cds-dispatch/prompt.md`'s own steps" — δ verified this is factually wrong: `check-dispatch-repair-preflight.sh` runs as a standalone CI job (`dispatch-repair-preflight` in `.github/workflows/build.yml`), triggered on every PR, unrelated to the dispatch wake's own run. This correction does not change γ's underlying two-guard-split design (which δ pinned as final), only where Guard A is wired (a new CI job mirroring the #516 job's exact shape) and confirms Guard B's integration point (the dispatch wake's own run, immediately before the `status:in-progress → status:review` label transition) on its own independent merits — not because of the #516 mis-citation. See §Self-check below for how this resolution was followed exactly as δ pinned it.
