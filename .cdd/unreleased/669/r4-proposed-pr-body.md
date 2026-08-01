# Proposed PR body after R4

Closes #669. Draft pending fresh exact-SHA R4 β/γ/CC review.

This publishes the recursive-cell CM as a source-checkout-only, zero-standing
candidate and ships the deterministic #662 State-A route described by the
methodology:

- emits six registered prompts plus a separate H01-H13 assessment prompt and
  binds all prompt and response bytes;
- validates exact standard TSC v3.2.4 witnesses, including exact typed
  nonempty `next_fixes` entries;
- runs six hybrid targets, computes the unweighted L0-L4 aggregate, and carries
  the selected level's canonical TSC bottleneck axis into the result;
- atomically emits one immutable prompt bundle, snapshots it with all responses
  and the invariant assessment, runs coh only from snapshots, and CUE-validates
  a digest-binding marker before atomically creating immutable `publication/`;
- refuses reuse, symlinked or locked destinations, response replay, malformed
  witness shapes, mid-target failure, and final-schema failure without exposing
  partial canonical output;
- binds the active runner, coh executable, schema, `SKILL.md`, instruction and
  assembly path, assessment, registry/preflight/manifests, prompts, responses,
  reports, target, and caller-supplied time;
- adds deterministic no-live-model CI and separately passes six-target
  compatibility against pinned `coh 0.12.0 (26aab50)` with tracked synthetic
  witnesses; neither is a semantic calibration sample; and
- preserves prior pinned-engine evidence while keeping its same-author N=1
  superseded after the instrument change and ineligible for future k=3.

State-B arbitrary methodology/target loading and installed activation remain
unshipped. #662 remains a design/calibration source, not a held-out target, and
this PR claims no independent, admissibility, consistency, held-out, hosting-
independent, or off-house standing.

External β R3 comment `5016263226` requested iteration over exact SHA
`bc6f6aa946cbbb555bb409c8bc180f5d49685b6e`. R4 is a focused successor commit
and requires a fresh exact-SHA β review, followed only on convergence by fresh γ
binding and CC re-adjudication. No merge or control-plane transition is
authorized by this draft.
