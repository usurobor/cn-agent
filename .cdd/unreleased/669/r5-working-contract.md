# R5 Working Contract and Fresh Review Request — cnos#669

## Control-plane binding

- Issue: `#669`
- PR: `#670`
- Exact reviewed R4 head / R5 implementation base:
  `8baf9ead0ff773a23812db1153d1bcfa27bedf2f`
- External β R4 verdict: **ITERATE**
- β comment ID: `5016473801`
- β comment URL:
  `https://github.com/usurobor/cnos/pull/670#issuecomment-5016473801`

R5 is one focused successor commit and does not amend R4. All earlier receipts
and reviews remain immutable historical evidence; none is a verdict over R5.

## Exact scope

This remains a zero-standing, source-checkout-only WC. R5 repairs only the two
standard-contract mismatches identified by β:

1. treat the six defect-card fields `id`, `primary_axis`, `category`,
   `severity`, `evidence`, and `summary` as required, with `secondary_axes`
   optional and semantically defaulted to `[]` when omitted;
2. when `secondary_axes` is present, require a list of unique valid axes that
   excludes the primary axis, and continue refusing every undeclared card key;
3. exercise one nonempty positive card that omits `secondary_axes` through both
   the State-A runner and pinned real coh, plus malformed optional-field
   refusals;
4. declare all `methodology.execution.emits` paths relative to the output root
   under `emission/`, direct the README assessor to the exact emitted assessment
   path, and make the path-base smoke snippet self-contained; and
5. regression-check the authoritative path declaration and documentation
   against the complete tree created by a fresh emit.

R5 does not change State-A custody, atomic emission/publication, invariant or
response bindings, aggregate/disposition semantics, State-B status, installed
activation, semantic evidence, standing, #662 status, or any control-plane
state.

## Acceptance and review request

Acceptance requires a single commit based on the exact R4 SHA, deterministic
and real pinned-coh compatibility for the omitted optional field, all local and
shallow checks green, updated draft PR truth, and an immutable α receipt naming
the final R5 SHA and both workflow results. Unrelated #429 material must remain
untouched.

After α stops writing, dispatch a fresh context-isolated functional β review of
the exact R5 SHA. Shared OpenAI Codex hosting and GitHub identity must be
disclosed; functional separation is not hosting independence or off-house
standing. No merge, ready mark, label/FSM transition, or standing promotion is
authorized.
