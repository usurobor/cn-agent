# Recursive-cell hard-invariant assessment

Assess H01-H13 only after reading all six emitted TSC prompts and all six standard
TSC witness responses. Do not add invariant fields to a TSC witness. Produce one
separate JSON object with this exact shape:

```json
{
  "schema": "cnos-recursive-cell-invariants/0.2",
  "target": "cc662-system",
  "assessment_prompt_sha256": "<sha256 of this materialized prompt>",
  "prompt_sha256": {
    "cc662-system": "<sha256>",
    "cc662-l0": "<sha256>",
    "cc662-l1": "<sha256>",
    "cc662-l2": "<sha256>",
    "cc662-l3": "<sha256>",
    "cc662-l4": "<sha256>"
  },
  "response_sha256": {
    "cc662-system": "<sha256 of exact witness bytes>",
    "cc662-l0": "<sha256 of exact witness bytes>",
    "cc662-l1": "<sha256 of exact witness bytes>",
    "cc662-l2": "<sha256 of exact witness bytes>",
    "cc662-l3": "<sha256 of exact witness bytes>",
    "cc662-l4": "<sha256 of exact witness bytes>"
  },
  "items": [
    {
      "id": "H01",
      "status": "pass | fail | unknown",
      "evidence": "specific evidence",
      "level": "system | L0 | L1 | L2 | L3 | L4",
      "primary_axis": "alpha | beta | gamma",
      "next_mca": null
    }
  ]
}
```

Include H01 through H13 exactly once and in numeric order. `pass` requires
specific evidence and a null `next_mca`. `fail` and `unknown` require a concrete
`next_mca`; both are nonaccepting. A failed item also requires a systemic defect
card in the standard witness for the stated level (or system), on the stated
primary axis, with `[Hxx]` in its evidence or summary. This cross-reference does
not change the standard TSC witness schema.

Fill `response_sha256` only after all six witness files are frozen. Hash the
exact bytes ingested by the runner. The assessment is invalid if either its
six-prompt map or six-response map differs from the run inputs.

Assess these stable invariants:

- H01: one CCNF kernel for all classes;
- H02: class determined by canonical output telos;
- H03: class and matter domain orthogonal;
- H04: alpha and beta are distinct for the reviewed matter;
- H05: kappa is outside and unequal to alpha, beta, gamma, and delta;
- H06: epsilon is an observation/projection, not a cell role or CC itself;
- H07: immutable matter, review, CM, target-bundle, measurement, and receipt links;
- H08: CM/V/delta separation with a type-correct data path;
- H09: single dependency authority;
- H10: single transition authority per FSM edge;
- H11: complete wave-FSM semantics before executable claims;
- H12: honest shipped/specified/illustrative partition;
- H13: CC judgment distinct from FSM transition and operator acceptance.
