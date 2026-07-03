# self-coherence — cycle/569

manifest:
  planned_sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
  completed: [Gap, Skills]

## §Gap

- **Issue:** [#569](https://github.com/usurobor/cnos/issues/569) — "cds/issues: FSM Phase 2 — authority flip (FSM applies labels; workers request transitions)"
- **Parent/wave:** #567 (master). Depends on #568 (Phase 1 read-only reconciler, closed/merged) and #570 (cell-kind taxonomy, closed/merged).
- **Mode:** design-and-build. γ's scaffold (`.cdd/unreleased/569/gamma-scaffold.md`) names one real open design decision left to α: the exact guard combination for the new `in-progress → review` proposal rule, and the 404-tolerance choice for the label-remove call. Both are resolved and documented below (§ACs, §Self-check).
- **Version/mode:** unreleased cycle, no release version cut yet; PR not opened by α per dispatch instructions (δ opens it after β converges).
- **Base:** cycle branch `cycle/569`, created by γ from `origin/main@0520235e1285c078eb3bc9d7eeba191b0413c53b`.
- **Binding scope constraint (operator comment, 2026-07-03T21:34:35Z on #569):** Phase 2 enforces FSM-controlled transitions **only** for the implementation-cell lifecycle already in `transitions.json` (`ready/todo/in-progress/review/changes`). No `cell_kind`-based enforcement branching; `FactSnapshot.CellKind` stays observed-only. Honored — see §ACs / scope-guardrail evidence below.

## Gap statement

Phase 1 (#568) gave the FSM the ability to *observe and reason about* issue state, read-only. It had no authority: workers (concretely, the δ-driven cds-dispatch wake) still wrote `status:*` labels directly, which is the root cause of the empty-review / stranded-in-progress / label-drift failure class the issue names (cnos#524 W4). This cycle moves that authority: the FSM applies status labels; workers produce matter (PR, commits, `REVIEW-REQUEST.yml`) and *request* transitions via a new guard-gated `--apply` flag on the existing `cn issues fsm evaluate` verb.

## §Skills

**Tier 1:** `CDD.md` (canonical lifecycle) + `alpha/SKILL.md` (this role surface).

**Tier 2 (`eng/*`, always-applicable):** `cnos.eng/skills/eng/SKILL.md` (coding bundle); `cnos.eng/skills/eng/go/SKILL.md` (Go conventions — package co-location per cnos#568/#556 precedent, dispatch-boundary rule INVARIANTS.md T-002 / `eng/go` §2.18, respected by keeping `cmd_issues_fsm.go` a thin doc-comment-only edit); `cnos.core/skills/write/SKILL.md`.

**Tier 3 (issue-specific):** `cnos.cdd/skills/cdd/issue/SKILL.md` (issue-pack contract); `cnos.cdd/skills/cdd/issue/proof/SKILL.md` (proof-plan discipline — AC oracle list was γ-authored, α proves each with fixture/CLI evidence below); `cnos.cdd/skills/cdd/issue/constraints/SKILL.md` (constraints discipline — the operator's cell-kind-deferral comment and the "no broad label redesign" non-goal are both hard constraints honored, not suggestions).

**Consumed but not loaded as role skills (read for context per §2.1.5 of `alpha/SKILL.md`, artifact enumeration):** `cnos.cds/skills/cds/fsm/{table,transitions}.json`-adjacent Go source (`table.go`, `snapshot.go`) to understand the guard-registry / rule-matching engine before extending its declarative data; `cnos.cds/orchestrators/cds-dispatch/SKILL.md`, `cnos.cdd/skills/cdd/delta/SKILL.md` §9.5–§9.6, `cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.4/§2.9 — the AC2 prose-edit targets named in the γ scaffold; `cnos.core/commands/install-wake/cn-install-wake` (the golden/live-workflow renderer, since editing cds-dispatch's SKILL.md body requires re-rendering `cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml` to stay in sync per `install-wake-golden` CI).
