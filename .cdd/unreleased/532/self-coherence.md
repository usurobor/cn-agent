# self-coherence — cycle/532

manifest:
  completed: [Gap, Skills]

## Gap

Issue: [#532](https://github.com/usurobor/cnos/issues/532) — "cdd/cds: require review-request proof before status:review"

A dispatch cell was able to set an issue to `status:review` without producing reviewable matter (observed in #524 W4: `status:review` + no PR + no commits beyond base + no diff + no receipt + no closeout + no STOP comment). CDD's closure sequence (`contract → matter → review → receipt → verdict → decision`) does not currently gate the `status:review` label transition on any mechanical proof that matter exists — the label is writable by the cell on a self-assertion alone.

The operator's pre-claim clarification comment (2026-06-30T19:29:09Z) is the authoritative, binding interpretation of AC4/AC5 and overrides any ambiguity in the issue body: the guard must be a mechanical state-integrity check, not prompt guidance alone; deliverable proof for the current GitHub/CDS surface means a PR (or approved no-op), commits beyond base (or no-op approval), a non-empty diff (or declared no-op), a `REVIEW-REQUEST.yml`, required closeout artifacts, and a review request that names PR/branch/head SHA/checks; on missing proof the guard must fail and the run must post STOP/BLOCKED, not leave `status:review` standing unproven.

**Version / mode:** `design-and-build` per the issue header (no stable `DESIGN.md`/`PLAN.md` existed at scoping time; γ confirmed MCA preconditions did not hold). This cycle both designs the `REVIEW-REQUEST.yml` shape (drafted as a strong proposal in the issue body, not yet a converged schema doc) and builds the guard in the same cycle, per the issue's own "Cycle scope sizing" decision (keep whole as one design-and-build issue).

**δ's binding correction to γ's scaffold (followed, not γ's mis-citation):** γ's scaffold claimed the cnos#516 guard is "invoked from inside `cds-dispatch/prompt.md`'s own steps" — δ verified this is factually wrong: `check-dispatch-repair-preflight.sh` runs as a standalone CI job (`dispatch-repair-preflight` in `.github/workflows/build.yml`), triggered on every PR, unrelated to the dispatch wake's own run. This correction does not change γ's underlying two-guard-split design (which δ pinned as final), only where Guard A is wired (a new CI job mirroring the #516 job's exact shape) and confirms Guard B's integration point (the dispatch wake's own run, immediately before the `status:in-progress → status:review` label transition) on its own independent merits — not because of the #516 mis-citation. See §Self-check below for how this resolution was followed exactly as δ pinned it.

## Skills

Tier 1 (loaded per dispatch §2.1, always):

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface; governed this cycle's execution (§2.2 artifact order, §2.3 peer enumeration, §2.6 pre-review gate, §2.7 review-readiness signal).

Tier 3 (named in the dispatch prompt; all loaded):

- `src/packages/cnos.core/skills/write/SKILL.md` — governed prose density in the doctrine additions (CDD.md, dispatch-protocol/SKILL.md §"Review-request proof gate", cds-dispatch step-7 edits): front-loaded points, one governing question per section, no throat-clearing.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — read to interpret AC boundaries and verify the issue's own handoff-checklist completeness (confirmed complete by γ's scaffold; re-verified independently during implementation — ACs numbered/testable, non-goals present and not leaking into ACs, source-of-truth paths all resolve).
- `src/packages/cnos.cdd/skills/cdd/issue/contract/SKILL.md` — governed reading the issue's Problem/Status-truth/Scope/Non-goals sections precisely (e.g. confirming "Status truth" table's "Missing" claims before trusting them — re-verified via `rg -i review-request` repo-wide, 0 real hits beyond γ's own scaffold).
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — governed AC-to-oracle mapping; every AC below names a concrete, runnable oracle (script invocation or grep), not "it works."
- `src/packages/cnos.cdd/skills/cdd/issue/constraints/SKILL.md` — governed the constraint-strata reading: AC8's "no label taxonomy change" is a hard gate (verified via diff grep for `gh label create`/rename, 0 hits beyond gamma-scaffold's own discussion text); cross-surface projections enumerated in §Self-check below.

Not found (noted, not blocking per dispatch instruction): `src/packages/cnos.core/skills/eng/SKILL.md` does not exist at that path — confirmed via `ls` before starting; no Go/CLI surface was touched this cycle (δ's pinned contract: "No Go source changes"), so this absence had no material effect on the work.

Role-skill-only constraint honored: did not load β/γ/δ role skills as primary governance (per `alpha/SKILL.md` §2.1 "do not load β or γ role skills"); `delta/SKILL.md` and `gamma-scaffold.md` were read as **input artifacts** (the dispatch prompt and the scaffold), not as α's own governing skill — consistent with the role boundary.
