# schemas/cdd/

Machine-checkable [CUE](https://cuelang.org/) schemas for the CDD coherence
cell's parent-facing artifacts. This directory types the validator surface
that [`src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`](../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md)
froze in prose. The schemas are **draft** until Phase 3 of [#366](https://github.com/usurobor/cnos/issues/366)
implements the validator `V` against them; once Phase 3 lands, schemas+validator
become the binding parent-facing trust surface. This cycle is [#369](https://github.com/usurobor/cnos/issues/369),
Phase 2 of the coherence-cell executability roadmap.

The semantic owner of the validation interface remains `RECEIPT-VALIDATION.md`.
The schemas here type the shape the doctrine document specifies. Each `.cue`
file's header comment points back to this README's `## Scope-Lift Invariant`
section so a reader entering through a CUE file lands on the doctrine surface
in one hop.

## Files

| File | Defines | Schema ID |
|---|---|---|
| `contract.cue` | `#Contract` | `cnos.cdd.contract.v1` |
| `boundary_decision.cue` | `#ValidationVerdict`, `#BoundaryDecision`, `#Override`, `#Transmissibility` | `cnos.cdd.boundary_decision.v1` |
| `receipt.cue` | `#Receipt`, `#ProtocolGapRef` | `cnos.cdd.receipt.v1` |

All three files are in CUE package `cdd`. `receipt.cue` references the boundary
decision types directly (same package; no `import` statement). `v1` is the
**schema-artifact** version, not a CDD protocol version.

The fixture corpus is at `schemas/cdd/fixtures/`:

| Fixture | Role |
|---|---|
| `valid-receipt.yaml` | Demonstrates a typical PASS / `accept` / `accepted` receipt; passes `cue vet -c -d '#Receipt'` |
| `invalid-override-masks-verdict.yaml` | Asserts `transmissibility: accepted` on `verdict: FAIL` + `action: override`; fails vet for the transmissibility mismatch |
| `invalid-fail-no-boundary-decision.yaml` | `verdict: FAIL` with no `boundary_decision`; fails vet for the missing required field |
| `invalid-gamma-preflight-authoritative.yaml` | Asserts γ-preflight as authoritative (via a `preflight_only: true` marker) and omits `boundary_decision`; fails vet because the schema requires `boundary_decision` regardless of any γ-preflight signal |

## Scope-Lift Invariant

A coherence cell's parent-facing receipt is not just a record of what
happened inside the cell — it is the matter the **parent scope** reads
when the cell crosses the trust boundary. The recursion equation in
[`COHERENCE-CELL.md`](../../src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md)
§Recursion Equation states this in compressed form. This section names the
**three projections** the schemas in this directory preserve under
scope-lift. These are projections — distinct surfaces at scope `n` that map
onto distinct surfaces at scope `n+1` — not a flat renaming of roles inside
a single cell.

### Projection 1 — Closed cell → α-matter

> A **closed α/β/γ cell at scope `n`** is **α-matter at scope `n+1`**.

The cell at scope `n` runs its α/β/γ metabolism — α produces, β discriminates,
γ closes. The closed-cell record (issue, diff, the five per-cycle markdown
files, the typed receipt) is what becomes available to scope `n+1` as input
material — the substrate scope-`n+1`'s α reasons over.

Schema consequence: `#Contract` (`cnos.cdd.contract.v1`) types the
accountability surface scope `n+1` checks the receipt against. `#Receipt`
(`cnos.cdd.receipt.v1`) types the parent-facing artifact itself — including
the typed evidence refs scope `n+1`'s `V` resolves to verify the receipt's
claims. Without these typed surfaces, scope `n+1` would have to re-derive
structure from the in-cell narrative each time — exactly what the receipt
exists to fix.

### Projection 2 — δₙ BoundaryDecision → βₙ₊₁-like discrimination

> The **δₙ `#BoundaryDecision`** (this scope's boundary action) **projects to a
> βₙ₊₁-like discrimination at scope `n+1`** (the parent's review of incoming
> matter).

When δ at scope `n` records `action: accept | release | reject |
repair_dispatch | override` plus a `#Transmissibility`, scope `n+1` reads
that signal as the **discrimination input** for its own β-like layer at the
parent scope. A `transmissibility: not_transmissible` receipt cannot be used
as scope-`n+1` α-matter; a `transmissibility: degraded` receipt can be used
but only under explicit acknowledgment of the degradation; a
`transmissibility: accepted` receipt is the clean-closure case scope `n+1`
inherits without further discrimination work on this cell.

Schema consequence: `#BoundaryDecision` and `#Transmissibility` live in
`boundary_decision.cue`. The verdict × action → transmissibility table is
enforced **structurally** in `receipt.cue` (the if-chain at the bottom of
`#Receipt`). Authors cannot relitigate the mapping — a fixture asserting
`verdict: FAIL, action: override, transmissibility: accepted` (override
masking the non-PASS verdict) fails `cue vet`. This is the load-bearing
mechanism that prevents projection-2 drift: scope `n+1`'s read of the
boundary decision is reliable because the decision's transmissibility
projection is computed, not asserted.

### Projection 3 — εₙ receipt-stream → γₙ₊₁-like coordination

> An **εₙ receipt-stream observation** (the protocol-gap signal across many
> cells at scope `n`) **projects to a γₙ₊₁-like coordination/evolution at scope
> `n+1`** (parent-scope protocol authoring).

ε at scope `n` reads the receipt stream — many cells' `#ProtocolGapRef`
entries — and surfaces patterns that require protocol patches. Those
patches are scope-`n+1` γ-like work: coordinating which protocol gaps to
close, sequencing the fixes, evolving the rules the next round of scope-`n`
cells will run under. ε's signal at scope `n` is therefore the **input**
to the parent scope's coordination layer at scope `n+1`.

Schema consequence: `#Receipt` requires `protocol_gap_count: int & >= 0`
and `protocol_gap_refs: [...#ProtocolGapRef]`, with the consistency
constraint `protocol_gap_count == len(protocol_gap_refs)`. Both fields are
required in `v1` (not deferred to `v2`) so Phase 6's ε relocation reads them
without forcing a schema bump (see [#366](https://github.com/usurobor/cnos/issues/366)
Phase 6 / [#369](https://github.com/usurobor/cnos/issues/369) `## Impact`).
`#ProtocolGapRef` is a structured object (`id`, `source` enum, `ref`) rather
than an opaque string, so scope-`n+1` ε-relocation reads typed observations,
not freeform notes.

### Projection-under-scope-lift, not role-renaming

The three projections above are **distinct surfaces at each scope**, not a
flat aliasing of roles. The δ at scope `n` does not *become* the β at scope
`n+1`. The δ at scope `n` performs δ work (boundary action on its own
receipt); the β at scope `n+1` performs β work (discrimination on
incoming matter); the projection maps the **output** of the former to a
**typed input** for the latter. The same applies for closed-cell→α-matter
and ε-stream→γ-coordination.

Reading the projection as flat role-renaming collapses two surfaces that
must stay distinct. The schemas keep them distinct: a `#Receipt` is a
typed artifact (not "a β-review at scope n+1"); a `#BoundaryDecision` is a
boundary action (not "a discrimination verdict from scope n+1's β"); a
`#ProtocolGapRef` is an observation (not "a γ next-commitment from scope
n+1"). The projection is what crosses the scope-lift; the surfaces it lifts
between remain typed at each end.

## How to run

You need [CUE](https://cuelang.org/) v0.10+ on your `PATH`. From the repo
root:

### Schema-only compile

Checks that all three schemas parse, type-check, and unify without concrete
values (`-c=false` allows incomplete schema values). Use this when editing
the schemas themselves:

```
cue vet -c=false \
  schemas/cdd/contract.cue \
  schemas/cdd/boundary_decision.cue \
  schemas/cdd/receipt.cue
```

### Fixture validation

Validates a fixture (a concrete receipt YAML) against `#Receipt` with all
three schemas in scope. `-c` requires concrete values for every required
field; `-d '#Receipt'` points the validator at the receipt definition.
Substitute `{fixture}` with the receipt fixture filename:

```
cue vet -c -d '#Receipt' \
  schemas/cdd/contract.cue \
  schemas/cdd/boundary_decision.cue \
  schemas/cdd/receipt.cue \
  schemas/cdd/fixtures/{fixture}.yaml
```

The valid fixture `valid-receipt.yaml` succeeds; each of the three named
invalid fixtures fails for its named load-bearing reason. Phase 3 of
[#366](https://github.com/usurobor/cnos/issues/366) will wrap the same
invocation as the operator-facing `cn-cdd-verify` command's structural
reader; Phase 3 inherits this invocation pattern, not a re-derived one.

## What this directory does NOT do

This directory is **schema-only**. It does not:

- **Implement `V`.** The validator predicate (V at
  [RECEIPT-VALIDATION.md](../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md)
  §Validation Interface) is Phase 3 of [#366](https://github.com/usurobor/cnos/issues/366).
  These schemas are V's *input contract*, not V itself.
- **Type the full evidence graph.** `#Receipt` types evidence *refs* (string
  references to artifacts and SHAs) so Phase 3 has a pinned input shape.
  The graph's own schema (`evidence_graph.cue`) is Phase 3 territory.
- **Modify `cn-cdd-verify`.** The current artifact-presence checker at
  `src/packages/cnos.cdd/commands/cdd-verify/` is unchanged this cycle.
  Phase 3 refactors it into V's operator-facing wrapper.
- **Edit `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, or `CDD.md`.** The
  doctrine and design surfaces remain authoritative; the schemas reference
  them via header comment.
- **Pin a protocol version.** `v1` refers to each schema artifact's version
  (`cnos.cdd.contract.v1`, `cnos.cdd.receipt.v1`, `cnos.cdd.boundary_decision.v1`),
  not to a CDD protocol version. The CDD protocol is unversioned at this layer.

## Related issues and surfaces

- [#369](https://github.com/usurobor/cnos/issues/369) — this work; Phase 2.
- [#366](https://github.com/usurobor/cnos/issues/366) — coherence-cell
  executability roadmap; this directory is its Phase 2 deliverable.
- [#367](https://github.com/usurobor/cnos/issues/367) — Phase 1; froze the
  parent-facing validator interface in
  [`RECEIPT-VALIDATION.md`](../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md).
- [#364](https://github.com/usurobor/cnos/issues/364) (closed) — doctrine;
  [`COHERENCE-CELL.md`](../../src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md).
- [`schemas/skill.cue`](../skill.cue) + [`schemas/README.md`](../README.md) +
  [`schemas/fixtures/skill-frontmatter/`](../fixtures/skill-frontmatter/) —
  schema-layer pattern precedent: CUE owns shape, README owns prose,
  fixtures live under `schemas/{subsystem}/fixtures/`.
