# α R5 Self-Coherence and Repair Receipt — cnos#669

## Binding and role boundary

- Exact R5 implementation base:
  `8baf9ead0ff773a23812db1153d1bcfa27bedf2f`
- Triggering external β verdict: **ITERATE**, comment `5016473801`
- β verdict URL:
  `https://github.com/usurobor/cnos/pull/670#issuecomment-5016473801`
- Final R5 head: bound by α's immutable PR comment after the one successor
  commit; embedding a commit's own SHA in its tracked bytes would be circular.

This is same-author α implementation evidence, not β review, γ closeout, CC
judgment, independent warrant, ratification, or standing promotion.

## Repair account

The runner now matches the pinned TSC v3.2.4 defect-card contract. Every card
must contain exactly the six required fields and may additionally contain
`secondary_axes`. Omission maps to an empty list for validation. When present,
the field must be a list of unique valid axes excluding the primary axis;
unknown card fields still refuse. The positive L1 fixture contains a cosmetic
card with all six required fields and deliberately omits `secondary_axes`.
Runner tests prove that witness publishes successfully and separately refuse a
non-list value, the primary axis, an invalid axis, duplicate axes, and an extra
card field.

The exact State-A path contract is also coherent again. README assessment
instructions name
`$RUN_ROOT/emission/invariant-assessment-prompt.md`; authoritative frontmatter
declares all eight output-root-relative `emission/...` paths; and the path-base
smoke block defines `CM` locally. The runner regression asserts every declared
path exists in a fresh emission, no undeclared file is present, and the README
and SKILL declarations contain the exact executable paths.

## Evidence boundary

The final runner, real pinned-coh, frontmatter, package, historical target, and
depth-one results are recorded in α's immutable post-commit comment. Real-coh
acceptance is CLI/schema compatibility over a tracked synthetic witness, not a
new semantic sample. R5 leaves immutable custody and every State-A/B and zero-
standing claim unchanged. The historical same-author N=1 remains superseded
and ineligible for a later k=3 claim.

Fresh exact-SHA β review, then only on convergence new γ binding and CC
re-adjudication, remain required.
