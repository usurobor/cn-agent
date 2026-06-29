# Score Sheet

**Kata ID:** A1-open-op-registry-conflict
**Run ID:** A1-cold-001

| Dimension | Score (0–4) | Why |
|-----------|-------------|-----|
| α Pattern | 2 | It recognized the conflict but proposed a non-deterministic or semantically weak resolution (warn + overwrite / keep first). |
| β Relation | 2 | It partially respected the extension architecture but blurred registry validation with runtime dispatch. |
| γ Exit | 2 | It proposed tests and a patch, but not the strongest invariant-preserving move. |
| Efficiency | 4 | Cheap, fast, one round. |

## Positive Signals

- [x] Saw that duplicate names and duplicate op kinds both matter.
- [x] Proposed at least some regression tests.

## Negative Signals

- [x] Failed to enforce deterministic rejection.
- [x] Did not treat silent shadowing as categorically wrong.
- [x] Did not suggest stronger proof like order-independence/property testing.

## Transfer Assessment

- none yet

## Cost Assessment

- cheap

## Final Classification

- selection-limited
