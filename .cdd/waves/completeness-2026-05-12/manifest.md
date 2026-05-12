# Wave: CDD Completeness + Infrastructure Hardening

**Date:** 2026-05-12
**Dispatcher:** δ (Sigma)
**Repo:** usurobor/cnos

## Issues (run in order)

| Order | # | Title | ACs | Type |
|-------|---|-------|-----|------|
| 1 | 354 | δ polls release CI after tag push | 5 | docs |
| 2 | 349 | Activation §11a notification channels | TBD | docs |
| 3 | 328 | CI artifact ledger checker | TBD | test+ci |
| 4 | 273 | Rebase-collision integrity guard | TBD | process |
| 5 | 323 | fix(activate): scanner misses threads/inbox/ | TBD | fix |

## Standing permissions

- Push to cycle branches: yes
- Push merges to main: yes
- Auto-dispatch α fix rounds on REQUEST CHANGES: yes (max 3)
- Tag/release: NO — operator gate
- Branch delete after merge: yes

## Timeout budgets

- γ: 1200s (§5.2 full-cycle)
- α: 900s
- β: 900s

## Dependencies

- #354 before #328 (ledger checker can reference release CI gate)
- #349 independent
- #273 independent
- #323 independent
