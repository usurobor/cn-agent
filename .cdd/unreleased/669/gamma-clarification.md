# γ Clarification — cnos#669 R6 base-drift oracle

**Date:** 2026-07-19
**Trigger:** first α replay completed; AC3's whole-tree oracle encountered
pre-existing `origin/main` advancement.
**Affected scaffold surfaces:** `## ACs` item 3 and `## AC oracle approach`.

## Observation

The superseded lineage branched at
`90287522ebf55b54f875973631506500a5f4f578`. Canonical R6 correctly starts
from the current `origin/main` base
`e8ba9954764d58e7b808104d633504e25aa615cc`. Between those bases, main changed
three paths unrelated to cnos#669:

- `.cn-sigma/logs/20260719.md`
- `docs/development/board/board-data.json`
- `docs/development/board/index.html`

Therefore literal final-tree equality to superseded R5 is neither reachable
nor desirable: it would revert current-main state. α correctly left those
paths untouched.

## Re-pinned oracle

Scaffold AC3 is clarified from whole-tree equality to **cycle-patch
equivalence outside `.cdd/unreleased/669/`**. At the first replay boundary:

- old patch: `90287522..9cc9c3bb`
- R6 patch: `e8ba9954..6be13f40`
- exact `git diff --binary` SHA-256 for both:
  `5024e9392d22ac76d3d152035a66dec401a98ea6ce87e2178976df0b59067c84`

The same comparison must be repeated at the final technical matter boundary:

```bash
git diff --binary 90287522ebf55b54f875973631506500a5f4f578..ccaf35607520ffb43d9450e35759e30d742355d5 \
  -- . ':(exclude).cdd/unreleased/669/**'

git diff --binary e8ba9954764d58e7b808104d633504e25aa615cc..<R6-final-alpha-matter> \
  -- . ':(exclude).cdd/unreleased/669/**'
```

The two byte streams must hash identically. This proves R6 replays the exact
technical change while preserving unrelated current-main evolution.

## Handoff

No product scope, CM contract, State-A/B claim, or role assignment changes.
The next α activation continues from this clarification and never edits γ's
scaffold or this file.
