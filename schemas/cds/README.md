# schemas/cds/

Machine-checkable [CUE](https://cuelang.org/) schema for the **CDS
(software development) protocol overlay** receipt shape. This package
extends the generic coherence-cell kernel ([`schemas/cdd/`](../cdd/))
with the required evidence-ref keys every CDS cycle's γ emits at
close-out: the closure-record artifacts (self-coherence, β-review, the
three role close-outs), the cycle diff, and optional CI references.

Established by cycle [#388](https://github.com/usurobor/cnos/issues/388)
(Phase 2.5 of [#366](https://github.com/usurobor/cnos/issues/366)),
which split the unified Phase 2 schema (#369) into a generic kernel
(`schemas/cdd/`) and per-protocol overlays (`schemas/cds/`, `schemas/cdr/`).

The architectural-choice decision (option (a) split) is recorded in
[`schemas/cdd/README.md` §"Architectural choice"](../cdd/README.md#architectural-choice--generic-vs-domain-split).

## Files

| File | Defines | Schema ID |
|---|---|---|
| `receipt.cue` | `#CDSReceipt` (extends `cdd.#Receipt`) | `cnos.cdd.cds.receipt.v1` |

The file is in CUE package `cds`. It imports the generic kernel via
`import "cnos.dev/cnos/schemas/cdd"`; resolution is handled by the
repo-level CUE module declaration at `cue.mod/module.cue`.

The fixture corpus is at `schemas/cds/fixtures/`:

| Fixture | Role |
|---|---|
| `valid-receipt.yaml` | Demonstrates a typical CDS PASS / `accept` / `accepted` receipt; passes `cue vet -c -d '#CDSReceipt'` |

This fixture was the Phase 2 (#369) valid fixture for the (then-unified)
`#Receipt`; it moved here as part of cycle/388's split and was reshaped
to nest the closure-record refs under `evidence_refs:` (matching the
generic kernel's typed-open primitive) and to pin `protocol_id` to
`"cnos.cdd.cds.receipt.v1"`.

## What `#CDSReceipt` requires

`#CDSReceipt: cdd.#Receipt & { ... }` inherits the generic Receipt's
structural invariants (transmissibility derivation, override-iff,
required `boundary_decision`, `protocol_gap_count == len(refs)`) and
adds:

- `protocol_id: "cnos.cdd.cds.receipt.v1"` — pinned. Phase 3's V
  dispatches on this string.
- `evidence_refs.evidence_root` — root of the cycle evidence corpus
  (typically `.cdd/releases/<version>/<cycle>/` or
  `.cdd/unreleased/<cycle>/`).
- `evidence_refs.self_coherence` — the cycle's `self-coherence.md`.
- `evidence_refs.beta_review` — the cycle's `beta-review.md`.
- `evidence_refs.alpha_closeout` — the cycle's `alpha-closeout.md`.
- `evidence_refs.beta_closeout` — the cycle's `beta-closeout.md`.
- `evidence_refs.gamma_closeout` — the cycle's `gamma-closeout.md`.
- `evidence_refs.diff` — the cycle's diff against base SHA.
- `evidence_refs.ci?` — optional list of CI artifact refs.

## Forward-looking debt — §4a.3 sketch vs realized shape

[`ROLES.md §4a.3`](../../ROLES.md) names a CDS receipt sketch using
forward-looking field names: `artifact_refs`, `test_refs`, `ci_refs`,
`diff_ref`, `debt_refs`. Cycle/388 implements the field set #369 actually
shipped (the closure-record shape): `self_coherence`, `beta_review`, the
three closeouts, `diff`, optional `ci`. A future cycle may align §4a.3
and `schemas/cds/receipt.cue` if the field names should converge.

## How to run

You need [CUE](https://cuelang.org/) v0.10+ on your `PATH`. From the
repo root:

```
cue vet -c -d '#CDSReceipt' \
  schemas/cds/receipt.cue \
  schemas/cds/fixtures/valid-receipt.yaml
```

The valid fixture succeeds. A fixture omitting any required
`evidence_refs.*` key (e.g. missing `self_coherence`) fails with a
diagnostic naming the missing field. A fixture asserting a different
`protocol_id` value (e.g. `cnos.cdd.cdr.receipt.v1`) fails with a
`conflicting values` diagnostic.

## What this directory does NOT do

- **Implement V.** Phase 3 of [#366](https://github.com/usurobor/cnos/issues/366)
  owns the runtime validator and the per-protocol dispatch layer.
- **Define CDS doctrine.** Doctrine and the role lifecycle live in
  [`src/packages/cnos.cdd/skills/cdd/`](../../src/packages/cnos.cdd/skills/cdd/).
  This package only types the receipt shape.
- **Migrate historical CDS receipts.** Existing receipts in
  `.cdd/releases/<version>/<cycle>/` continue to validate against the
  CDS schema as-is (the fixture corpus is the migration template). No
  bulk rewrite of historical receipts is in scope for cycle/388.
- **Carry generic invariants.** Those live in
  [`schemas/cdd/`](../cdd/). The CDS schema unifies with the generic
  schema and adds only the CDS-required keys.

## Related issues and surfaces

- [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5;
  established this package.
- [#369](https://github.com/usurobor/cnos/issues/369) — Phase 2; the
  unified schema this package was extracted from.
- [#366](https://github.com/usurobor/cnos/issues/366) — coherence-cell
  executability roadmap; Phase 3 (V implementation) consumes this
  package via `protocol_id` dispatch.
- [`schemas/cdd/README.md`](../cdd/README.md) — the generic kernel; this
  package extends it.
- [`schemas/cdr/`](../cdr/) — sibling package for research receipts.
- [`ROLES.md §4a.3`](../../ROLES.md) — receipts-enforce-discipline
  principle; names the CDS receipt sketch as the discipline anchor.
