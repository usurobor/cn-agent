# Wave: 2026-05-19 Protocol Hygiene

**Date:** 2026-05-19
**Launched by:** ε (acting as wave-planner) → δ (operator-side dispatcher)
**Repo:** usurobor/cnos
**Mode:** §5.2 wave-mode (per `operator/SKILL.md` §5.2 + `cdd/CDD.md` §1.4 γ=δ collapse)
**Branch policy:** each cycle on its own `cycle/{N}` branch from `origin/main`; merges to `main` per cycle close-out

## Purpose

Three independent skill-patch cycles, runnable in parallel. All three close named protocol gaps that surfaced during 2026-05-17 / 2026-05-18 work. Wave is **protocol hygiene** — small, fast, low-risk skill patches that reduce per-cycle friction before larger work (Phase 2.5, Phase 3) dispatches.

## Issues (parallel; no ordering)

| # | Title | Mode | Surfaces touched |
|---|---|---|---|
| 375 | γ-side pre-dispatch gate for gamma-scaffold.md (rule 3.11b symmetry) | skill-patch (docs-only disconnect) | `cdd/gamma/SKILL.md` §2.5 OR `cdd/CDD.md` step 3 |
| 378 | rule 3.11b discoverability under §5.2 wave-mode | skill-patch (docs-only disconnect) | `cdd/review/SKILL.md` §3.11b + `cdd/alpha/SKILL.md` §2.6 + `cdd/operator/SKILL.md` §5.2/§10 |
| 377 | codify cross-repo coordination protocol (intake / outbound / bilateral) | design-and-build | new `cdd/cross-repo/SKILL.md` + `cdd/gamma/SKILL.md` §2.1/§2.7 + `cdd/post-release/SKILL.md` §5.6b + `.cdd/iterations/cross-repo/README.md` |

## Standing permissions

- Push to cycle branches: yes
- Push merges to main: yes (per cycle β close-out)
- Auto-dispatch α fix rounds on REQUEST CHANGES: yes (max 3)
- Tag/release: NO — operator gate (docs-only disconnect class for #375 + #378; #377 is also docs-only since it ships skill files only)
- Branch delete after merge: yes

## Timeout budgets

- γ: 600s (issue read + branch creation + scaffold)
- α: 1200s (#377 is design-and-build, larger; #375 and #378 small skill-patch)
- β: 900s
- per cycle ceiling: 3 review rounds (α R1, β R1, α R2, β R2, α R3, β R3 — fail-the-cycle at R3 RC)

## Wave-level invariants

These hold for every cycle in this wave:

1. **No CI / runtime / release surface change.** All three cycles are skill-patch class; `git diff origin/main..HEAD --stat` must show changes only under `src/packages/cnos.cdd/skills/cdd/`, `src/packages/cnos.cdd/skills/cdd/issue/`, `src/packages/cnos.cdd/skills/cdd/review/`, `.cdd/iterations/cross-repo/`, or `.cdd/unreleased/{N}/`.
2. **Empirical anchors must be cited.** Each cycle's patched skill text cites the specific prior cycle/wave that surfaced the gap (per `gamma/SKILL.md` §2.2a precedent). Cycle anchors per issue:
   - #375: cycle #369 β R1 D1 → recovery path (a) at `227d2373` → β R2 APPROVED → merged at `ff54f2a0`
   - #378: cph cdr-refactor wave 2026-05-18 (master cph#11; subs cph#12, #13, #14, #15)
   - #377: bootstrap-cdr inbound proposal (cnos#376 filing, 2026-05-18) + tsc/cdd-supercycle outbound trace (cnos#331/#332, 2026-05-09) + sigma agent-activate-skill inbound (cnos#379, 2026-05-19) + cph coherence-drift-sweep-followup outbound (Rule 6 patch on `claude/file-cnos-cdr-issue-fi9Ld`, 2026-05-19)
3. **rule 3.11b: γ-scaffold required at dispatch.** Per the rule itself and #375's preventive (which this wave dispatches). Each cycle's γ must author `.cdd/unreleased/{N}/gamma-scaffold.md` before α dispatch. Wave-level γ-as-δ collapse does not exempt — these are per-cycle subs, not wave-internal subs of a parent (cph#11 was the cph wave-internal pattern; this wave's cycles are independent cnos cycles).
4. **Cross-cycle coordination.** #378 patches the same `review/SKILL.md` surface that #375 references and that #377's `cdd/cross-repo/SKILL.md` may reference. If two cycles touch the same file, β raises a cross-cycle finding at review and γ coordinates merge order (probably #375 first since it's smallest and γ-side-only, then #378, then #377; the order is suggestive — last-merge-wins on file conflicts; γ resolves at merge time).

## Known constraints

- γ-as-spawner does not work in this environment (nested `claude -p` fails silently); δ dispatches γ→α→β sequentially per cycle, in parallel across cycles.
- α must be reminded to commit implementation files (not just `self-coherence.md`) — known recurring issue.
- The β-skill Rule 6 patch from 2026-05-19 (cph iteration) is on `claude/file-cnos-cdr-issue-fi9Ld`, awaiting merge to main. If that branch merges before this wave dispatches, β cycles in this wave inherit Rule 6 (code-first oracle anchoring); otherwise β operates without it. Either is fine; the wave's cycles don't depend on Rule 6.

## ε iteration target

Findings from this wave land in `.cdd/iterations/wave-2026-05-19-protocol-hygiene.md` (per `cdd/post-release/SKILL.md` §5.6b convention; pattern matches cph's `.cdd/iterations/wave-2026-05-18.md`). γ-axis, β-axis, α-axis observations across the three cycles consolidate there.
