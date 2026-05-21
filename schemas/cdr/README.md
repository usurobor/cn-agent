# schemas/cdr/

Machine-checkable [CUE](https://cuelang.org/) schema for the **CDR
(research) protocol overlay** receipt shape. This package extends the
generic coherence-cell kernel ([`schemas/cdd/`](../cdd/)) with the
required evidence-ref keys [`ROLES.md §4a.3`](../../ROLES.md) names for
research discipline: claim refs (which claims this receipt asserts),
data refs (dataset / mount / manifest / checksum), method refs (script
paths + commit SHA), result refs (output file paths), a claim-status
enum, optional limitations, and an optional reproduction record (β
re-ran the producing commands, output matched).

Established by cycle [#388](https://github.com/usurobor/cnos/issues/388)
(Phase 2.5 of [#366](https://github.com/usurobor/cnos/issues/366)) as a
**skeleton**. The minimum CDR field set per §4a.3 is shipped; cnos#376
Sub 1 / Sub 2 may refine field types or cardinalities as part of the
CDR-package bootstrap.

The architectural-choice decision (option (a) split) is recorded in
[`schemas/cdd/README.md` §"Architectural choice"](../cdd/README.md#architectural-choice--generic-vs-domain-split).

## Files

| File | Defines | Schema ID |
|---|---|---|
| `receipt.cue` | `#CDRReceipt` (extends `cdd.#Receipt`), `#ClaimStatus`, `#Reproduction` | `cnos.cdd.cdr.receipt.v1` |

The file is in CUE package `cdr`. It imports the generic kernel via
`import "cnos.dev/cnos/schemas/cdd"`; resolution is handled by the
repo-level CUE module declaration at `cue.mod/module.cue`.

The fixture corpus is at `schemas/cdr/fixtures/`:

| Fixture | Role |
|---|---|
| `valid-cdr-receipt.yaml` | Demonstrates a minimal CDR receipt: observed-claim status, four populated evidence-ref lists, reproduction record. Passes `cue vet -c -d '#CDRReceipt'` |

## What `#CDRReceipt` requires

`#CDRReceipt: cdd.#Receipt & { ... }` inherits the generic Receipt's
structural invariants (transmissibility derivation, override-iff,
required `boundary_decision`, `protocol_gap_count == len(refs)`) and
adds:

- `protocol_id: "cnos.cdd.cdr.receipt.v1"` — pinned. Phase 3's V
  dispatches on this string.
- `evidence_refs.claim: [_, ...#EvidenceRef]` — ≥1 claim ref. A
  receipt that asserts no claim is not a research receipt.
- `evidence_refs.data: [_, ...#EvidenceRef]` — ≥1 data source ref.
  Per `ROLES.md §4a.3`: "no `data_refs` → V rejects research receipt."
- `evidence_refs.method: [_, ...#EvidenceRef]` — ≥1 method ref
  (script path + commit SHA).
- `evidence_refs.result: [_, ...#EvidenceRef]` — ≥1 result ref
  (output file path).
- `claim_status: #ClaimStatus` — enum of `observed | computed |
  inferred | hypothesized | indeterminate`. The calibration of the
  receipt's asserted claims.
- `limitations?: [...string]` — optional caveat list.
- `reproduction?: #Reproduction` — optional record of β having re-run
  the producing commands and verified output match.

## Empirical anchor

The CDR receipt shape is anchored in `usurobor/cph` (the research
empirical anchor per [cnos#376](https://github.com/usurobor/cnos/issues/376)).
cph's `.cdr/waves/<wave>/receipt.md` exemplifies the field set this
schema types. The cycle/388 fixture (`valid-cdr-receipt.yaml`) is a
synthetic minimal example modeled on that shape.

## Skeleton — not frozen

Cycle/388 ships the minimum CDR field set per `ROLES.md §4a.3` sketch:

- `#ClaimStatus` enum and `#Reproduction` struct are typed at v0.1
  cardinalities.
- `evidence_refs.{claim,data,method,result}` are list-typed with
  minimum-length-one.
- `limitations?` and `reproduction?` are optional at v0.1.

cnos#376 Sub 1 / Sub 2 may refine:

- Make `reproduction` required when `claim_status ∈ {observed, computed}`.
- Add `evidence_refs.protocol_compliance` for cph's
  protocol-compliance evidence.
- Add `evidence_refs.review_independence` (β re-ran independently of
  α's exact session).
- Split `data` into `data_input` and `data_derived` for clarity.

These refinements are out of scope for cycle/388 (named here as known
debt). A future cycle (Phase 2.5.1 or absorbed into cnos#376 Sub 2)
addresses them as the CDR package matures.

## How to run

You need [CUE](https://cuelang.org/) v0.10+ on your `PATH`. From the
repo root:

```
cue vet -c -d '#CDRReceipt' \
  schemas/cdr/receipt.cue \
  schemas/cdr/fixtures/valid-cdr-receipt.yaml
```

The valid fixture succeeds. A fixture asserting a claim with an empty
`data:` list fails with `incompatible list lengths (0 and 2)`. A fixture
with a non-enum `claim_status` (e.g. `certain`) fails with `conflicting
values` against each enum member.

## What this directory does NOT do

- **Implement V.** Phase 3 of [#366](https://github.com/usurobor/cnos/issues/366)
  owns the runtime validator and the per-protocol dispatch layer.
- **Define CDR doctrine or role overlays.** Doctrine and the CDR role
  lifecycle live in [`cnos#376`](https://github.com/usurobor/cnos/issues/376)'s
  package bootstrap; this package only types the receipt shape.
- **Define the cph project binding.** cph's project binding (data
  mounts, project-specific gates, harness contracts) is project-local
  per `ROLES.md §4a` layer 4 — out of scope for this generic schema
  overlay.
- **Carry generic invariants.** Those live in
  [`schemas/cdd/`](../cdd/). The CDR schema unifies with the generic
  schema and adds only the CDR-required keys and top-level fields.

## Related issues and surfaces

- [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5;
  established this package as a skeleton.
- [#376](https://github.com/usurobor/cnos/issues/376) — CDR package
  bootstrap; will refine this schema.
- [#366](https://github.com/usurobor/cnos/issues/366) — coherence-cell
  executability roadmap; Phase 3 (V implementation) consumes this
  package via `protocol_id` dispatch.
- [`schemas/cdd/README.md`](../cdd/README.md) — the generic kernel; this
  package extends it.
- [`schemas/cds/`](../cds/) — sibling package for software receipts.
- [`ROLES.md §4a.3`](../../ROLES.md) — receipts-enforce-discipline
  principle; names the CDR receipt sketch this skeleton types.
