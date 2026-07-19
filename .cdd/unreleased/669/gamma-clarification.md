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

---

# γ Clarification — cnos#669 R6/R2 oracle correction

**Date:** 2026-07-19
**Trigger:** independent β R6/R1 receipt
`62135222a4de8d078cfb5cb1108532252dd19366` found that the opening scaffold
still stated final-tree/direct-diff equality after the first clarification.
**Affected scaffold surfaces:** `## ACs` item 3, `## Implementation contract`
row `Backward-compat invariant`, and `## AC oracle approach` item 3.

## Supersession

Those three scaffold statements are superseded. Neither final-tree identity
with `ccaf35607520ffb43d9450e35759e30d742355d5` nor an empty direct diff from
that commit is an R6 acceptance criterion. They are false because the R6 base
inherits three unrelated main-side paths that the superseded base lacks:

- `.cn-sigma/logs/20260719.md`
- `docs/development/board/board-data.json`
- `docs/development/board/index.html`

The scaffold remains immutable historical input. This γ-owned clarification
corrects its oracle without rewriting it or transferring γ authorship to α.

## Binding R6 oracle

The sole replay-equivalence criterion outside `.cdd/unreleased/669/**` is
byte equality of these two binary patch streams:

```bash
git diff --binary \
  90287522ebf55b54f875973631506500a5f4f578..ccaf35607520ffb43d9450e35759e30d742355d5 \
  -- . ':(exclude).cdd/unreleased/669/**'

git diff --binary \
  e8ba9954764d58e7b808104d633504e25aa615cc..<R6-final-alpha-matter> \
  -- . ':(exclude).cdd/unreleased/669/**'
```

At the pre-repair R6 boundary
`b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544`, both streams have SHA-256
`4c42826034a446e0953e0934085dbfd8ee6bb71febff39254c7a34c44a65d9c0`.
After α's symlink-boundary repair, β must compare the new R6 patch stream to
the intended R6 matter contract and review the additive repair directly; the
pre-repair equality digest remains provenance evidence, not a constraint that
forbids the required fix.

## R6/R2 handoff

α owns only the D1 runtime repair and its tests/claim projections. α appends a
new fix round to its own `self-coherence.md`, cites this corrected oracle, and
does not edit γ's scaffold, this clarification, or β's review. A fresh β round
must bind the resulting exact α SHA. No γ closeout, CC judgment, standing
promotion, merge, ready mark, or FSM transition is authorized before β
converges.
