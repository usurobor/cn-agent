# α R4 Self-Coherence and Repair Receipt — cnos#669

## Binding and authority

- Exact R4 implementation base:
  `bc6f6aa946cbbb555bb409c8bc180f5d49685b6e`
- Triggering external β verdict: **ITERATE**, comment `5016263226`
- β verdict URL:
  `https://github.com/usurobor/cnos/pull/670#issuecomment-5016263226`
- R4 head binding: the immutable α PR comment posted after the focused R4
  commit; a tracked receipt cannot contain its own commit SHA without
  circularity.

This is α's same-author repair receipt. It is not β review, γ closeout, CC
judgment, independent warrant, or standing promotion.

## Coherence and repair account

R4 makes standard TSC witness validation exact at the previously underchecked
`next_fixes` boundary. Entries must contain only `axis` and `fix`; axes are the
three TSC axes and fixes are nonblank strings. The positive fixture exercises a
nonempty list, while three negative fixtures cover an invalid axis, an extra
field, and a missing fix.

Emission and publication are now fresh-run-safe. Emission requires a unique run
root and builds the prompt bundle under a private stage and exclusive lock
before one rename creates immutable `emission/`. Ingest locks publication,
snapshots that complete emission plus all six response objects and the invariant
assessment, then validates, hashes, and runs coh only from those snapshots. All
replay inputs, six canonical reports, aggregate, and success marker exist in the
stage before CUE validation and one rename creates immutable `publication/`.
Existing, symlinked, and concurrently locked destinations refuse; successful
reuse cannot change bytes; mid-target and final-CUE failures clean stage and
lock and create no canonical output. The typed success marker binds every
replay input, result, and report path and digest.

The separate H01-H13 assessment now binds the exact six response digests, so a
valid assessment cannot be replayed against changed response bytes. Each level
record also carries its canonical TSC `bottleneck_axis`; after choosing the
lowest numeric L0-L4 `C_sigma` (level-order tie break), the aggregate uses that
reported axis rather than recomputing an axis from raw scores. A fixture makes
those values deliberately disagree so the contract is observable. Finally,
the authority bundle includes `SKILL.md` and the assembly script in addition to
the previously bound instrument and route files.

## Evidence boundary

The deterministic runner suite covers the exact witness shape, prompt/response
bindings, aggregate math, invariant gates, TSC bottleneck propagation, marker
digests, immutable emission/publication reuse, lock/symlink refusal, staged coh
input use, and cleanup after mid-target
and final-schema failures. The finished route is also checked against the real
pinned coh executable using tracked synthetic witnesses. This proves CLI and
schema compatibility, not semantic calibration, independence, admissibility,
or standing. R4 creates no semantic sample; the pre-R3 N=1 sample remains
superseded and ineligible for a later k=3 claim.

Final exact commands, commit SHA, push/PR workflow URLs, and their conclusions
belong in α's immutable PR comment after the focused commit is pushed. Fresh β,
then γ and CC actions remain required.
