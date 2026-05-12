# Wave: CDD Enforcement + Docs Breadth

**Date:** 2026-05-12
**Dispatcher:** δ-as-agent
**Repo:** usurobor/cnos

## Issues (run in order)

| Order | # | Title | Type |
|-------|---|-------|------|
| 1 | 356 | δ release gate blocks until green | docs |
| 2 | 355 | Structural compliance gates | enhancement |
| 3 | 292 | Resumption protocol for partial artifacts | docs |
| 4 | 341 | CDD README | docs |
| 5 | 340 | Reconcile front/docs README | docs |

## Standing permissions

- Push to cycle branches: yes
- Push merges to main: yes
- Auto-dispatch α fix rounds on REQUEST CHANGES: yes (max 3)
- Tag/release: NO — operator gate
- Branch delete after merge: yes

## Timeout budgets

- α: 900s
- β: 900s

## Known constraints

- γ-as-spawner does not work (nested claude -p fails silently)
- δ dispatches α→β directly, handles γ responsibilities (scaffold, close-out)
- α must be reminded to commit implementation files (not just self-coherence.md)
