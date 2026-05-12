# Wave: CDD Protocol Hardening

**Date:** 2026-05-12
**Dispatcher:** δ-as-agent (single claude -p session)
**Repo:** usurobor/cnos

## Issues (run in order)

| Order | # | Title | ACs | Type |
|-------|---|-------|-----|------|
| 1 | 352 | CI-green gate (β + γ PRA) | 5 | docs |
| 2 | 353 | Parent-session quiescence | 3 | docs |
| 3 | 351 | γ peer-enumeration | 4 | docs |
| 4 | 350 | Wave primitive | 6 | docs+design |
| 5 | 348 | ROLES.md §4 hats-vs-actors | TBD | docs |

## Standing permissions

- Push to cycle branches: yes
- Push merges to main: yes
- Auto-dispatch α fix rounds on REQUEST CHANGES: yes (max 3)
- Tag/release: NO — write to status.md only
- Branch delete after merge: yes

## Timeout budgets

- γ: 1200s (§5.2 full-cycle)
- α: 900s
- β: 900s

## Dependencies

- #352 before #350 (wave primitive can reference CI gate)
- #353 before #350 (wave primitive can reference quiescence)
- #351 independent
- #348 independent
