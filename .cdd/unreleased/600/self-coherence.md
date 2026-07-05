<!-- section-manifest
completed: [Gap, Skills]
remaining: [ACs, Self-check, Debt, CDD Trace, guard-inventory, Review-readiness]
-->

# self-coherence — cycle #600

## §Gap

**Issue:** [cnos#600](https://github.com/usurobor/cnos/issues/600) — "cds/runtime: consolidate strand-era guards and scaffold now that the mechanical lifecycle is closed (Sub D of #583)".

**Mode:** `design-and-build` (per the issue's own `**Mode:**` field).

**cell_kind:** `audit` — recorded per γ's scaffold; observed-only, does not gate any FSM transition (see AC4 below).

**Parent:** #583 (Sub D). **Preconditions (all merged, per issue body):** #584, #575, #591, #593.

**Base SHA at branch creation:** `1c5dd993fcd25b5bdca3843555d240895e3212ee`. `origin/main` advanced twice while this cycle ran (to `65a55ae5`, then to `eb94445b`, both trivial `activate+attach` heartbeat-log commits touching only `.cn-sigma/logs/*` — no overlap with this cycle's surfaces). The cycle branch was rebased onto `eb94445b` and force-pushed; see §CDD Trace / §Review-readiness for the exact SHAs.

**The gap this cycle closes:** after Sub B (`cn cell finalize`, #591) and Sub C (`cn issues fsm scan`, #593) made checkpoint/PR-open and dead-run recovery mechanical, a layer of CI guards and scaffold accreted *while the strand problem was still unsolved* (cnos#516 repair re-entry preflight, cnos#524 W4 closeout integrity, the cnos#574 review-guard tightening, the cnos#568/#570 `cell_kind` seam). This cycle is an **audit-first consolidation**: classify every strand-era guard KEEP/REMOVE/FOLD/NARROW, cite a live replacement for anything removed/folded, and leave no lifecycle invariant weakened. It is explicitly **not** a feature cycle — no new FSM behavior, no new labels, no `cell_kind` enforcement, no Demo 0.

**Dispatch intake:** `git fetch origin cycle/600 && git switch cycle/600` — branch pre-existed (γ created it), no branch was created by α, per the constraint. γ's `gamma-scaffold.md` (already on the branch at dispatch time) was read in full, along with the live issue body via `gh issue view 600`, before any classification work began. γ's guard-surface inventory table was treated as a starting hypothesis to verify, not settled fact, per the dispatch's explicit instruction.

## §Skills

**Tier 1 (loaded in full):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical kernel/role-contract doctrine.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — the α role surface governing this cycle: dispatch intake (§2.1), produce-in-artifact-order (§2.2), peer enumeration (§2.3), harness audit for schema-bearing changes (§2.4), incremental self-coherence discipline (§2.5), pre-review gate (§2.6, including the 15-row checklist), request-review (§2.7), the implementation-contract rule (§3.6).

**Tier 2/3, per the scaffold's "Skills to load" list and the actual diff surface:**
- `eng/process-economics` — the "guard-accretion vs boundary-fix" lens the issue names as its own governing frame: is a guard still doing boundary-fix work, or has the boundary moved (to the FSM/finalizer/reconciler) leaving the guard as pure accretion? Applied directly to the classification of both CI scripts.
- `go` — Tier 3, applied when reading/verifying `issuesfsm.go`, `scan.go`, `table.go`, `issuesfsm_test.go`, `scan_test.go`, and running `go test`/`go vet`.
- `eng/test` — applied to the AC2/AC3 oracle discipline (cite the specific test function, run it, prove green — not "the tests should cover this").
- Bash conventions (no dedicated Tier-3 `eng/bash` skill was loaded as a distinct file; the two guard scripts were edited following the existing script's own house style — `set -euo pipefail`, a `need()` helper, `::error::` GitHub Actions annotations — rather than introducing a new style).

**Not loaded:** `design/SKILL.md`, `plan/SKILL.md` as full artifacts — see §ACs / design-and-plan-not-required justification below. β/γ role skills were not loaded per α's Tier-1 restriction (`alpha/SKILL.md` §2.1 step 6: "do not load β or γ role skills").

**Active-skill-as-generation-constraint note (rule 3.1):** `eng/process-economics`'s guard-accretion lens directly shaped the classification method used throughout §ACs below — every KEEP/NARROW/FOLD/REMOVE decision below is argued in terms of "has the boundary that used to require this guard moved to a mechanical enforcer," not "does this guard look old."
