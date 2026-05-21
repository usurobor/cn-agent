# schemas/cdd/

Machine-checkable [CUE](https://cuelang.org/) schemas for the **generic
coherence-cell kernel**. This package types the substrate-independent
parent-facing artifacts every coherence-driven protocol (CDS, CDR, future
cdw / cda) shares: the contract, the receipt skeleton, and the boundary
decision. Domain-specific evidence-ref requirements (the software-process
closure records, the research claim/data/method/result refs) live in
sibling packages:

- [`schemas/cds/`](../cds/) — software protocol overlay
- [`schemas/cdr/`](../cdr/) — research protocol overlay

The schemas here type the shape the doctrine document
([`COHERENCE-CELL-NORMAL-FORM.md`](../../src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md))
specifies as the kernel layer. They are **draft** until Phase 3 of
[#366](https://github.com/usurobor/cnos/issues/366) implements the
validator `V` against them; once Phase 3 lands, schemas+validator become
the binding parent-facing trust surface. This package was split out from
the unified Phase 2 schema by cycle [#388](https://github.com/usurobor/cnos/issues/388)
(Phase 2.5).

The semantic owner of the validation interface remains
[`RECEIPT-VALIDATION.md`](../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md).
The schemas here type the shape the doctrine document specifies. Each
`.cue` file's header comment points back to this README's relevant
section so a reader entering through a CUE file lands on the doctrine
surface in one hop.

## Files

| File | Defines | Schema ID |
|---|---|---|
| `contract.cue` | `#Contract` | `cnos.cdd.contract.v1` |
| `boundary_decision.cue` | `#ValidationVerdict`, `#BoundaryDecision`, `#Override`, `#Transmissibility` | `cnos.cdd.boundary_decision.v1` |
| `receipt.cue` | `#Receipt`, `#ProtocolGapRef`, `#EvidenceRef`, `#EvidenceRefValue` | `cnos.cdd.receipt.v1` |

All three files are in CUE package `cdd`. `receipt.cue` references the
boundary decision types directly (same package; no `import` statement).
Domain packages (`schemas/cds/`, `schemas/cdr/`) import this package via
the repo-level CUE module declaration at `cue.mod/module.cue`. `v1` is
the **schema-artifact** version, not a CDD protocol version.

The fixture corpus is at `schemas/cdd/fixtures/`:

| Fixture | Role |
|---|---|
| `valid-generic-receipt.yaml` | Demonstrates `#Receipt` validates without any domain-specific evidence keys; proves the generic schema is usable on its own |
| `invalid-fail-no-boundary-decision.yaml` | `verdict: FAIL` with no `boundary_decision`; fails vet for the missing required field |
| `invalid-gamma-preflight-authoritative.yaml` | Asserts γ-preflight as authoritative (via a `preflight_only: true` marker) and omits `boundary_decision`; fails vet because the schema requires `boundary_decision` regardless of any γ-preflight signal |
| `invalid-override-masks-verdict.yaml` | Asserts `transmissibility: accepted` on `verdict: FAIL` + `action: override`; fails vet for the transmissibility mismatch |

The CDS-shaped fixture from the pre-cycle/388 corpus
(`valid-receipt.yaml`) moved to `schemas/cds/fixtures/`. The three invalid
fixtures stay here because they exercise *generic* constraints (boundary
decision required, override polarity, γ-preflight non-authoritativeness)
that apply across all protocols.

## Architectural choice — generic vs domain split

This section is the durable decision record for cycle/388's architectural
choice (#388 AC1). It is citable by cnos#376 Sub 1 (which inherits the
same architectural decision-class for skills/role-overlays per
cnos#376 AC7).

### The decision

**Option (a) — Split.** Generic CDD lives in `schemas/cdd/` (this
package). Software-process evidence requirements live in `schemas/cds/`
as `#CDSReceipt: cdd.#Receipt & { ... }`. Research evidence requirements
live in `schemas/cdr/` as `#CDRReceipt: cdd.#Receipt & { ... }`. Receipt
frontmatter carries a `protocol_id` field (e.g.
`"cnos.cdd.cds.receipt.v1"`, `"cnos.cdd.cdr.receipt.v1"`) so Phase 3's V
dispatches validation to the correct schema package without re-deriving
the protocol from field shape.

**Option (b) — Adapter (rejected).** Keep `schemas/cdd/receipt.cue` with
`evidence_refs` as an open typed map; CDS/CDR are CUE subtypes requiring
specific keys in the map. This is rejected for the reasons below.

### Rationale (1)–(5)

1. **CUE's named-field type system is the strong-typed primitive.**
   Required-keys-in-open-map (option b) shifts validation from CUE to V's
   runtime logic. Option (a) keeps the proof in CUE: the unification
   `cdd.#Receipt & { evidence_refs: { self_coherence: cdd.#EvidenceRef, ... } }`
   makes `self_coherence` a required key whose absence fails `cue vet`
   structurally.

2. **Clarity at the protocol boundary.** A CDR receipt author looking at
   `schemas/cdr/receipt.cue` sees exactly the required fields — no
   indirection. Option (b) would require looking outside the schema file
   to know which keys are required for which subtype.

3. **Mechanical generic-vs-domain boundary.** Different files plus CUE
   package names (`cdd`, `cds`, `cdr`) make the boundary structural.
   Option (b) makes it semantic — same file; required-key set varies by
   subtype.

4. **Future c-d-X generalizes mechanically.** A new protocol (cdw, cda,
   …) adds a directory: `schemas/cdw/receipt.cue` declares
   `#CDWReceipt: cdd.#Receipt & { ... }`. The generic surface does not
   need to know about it.

5. **Unifies with cnos#376 AC7.** The same reading ((a) split = common
   constitution + per-protocol procedures) applies to skills and
   role-overlays, not just schemas. The decision is recorded once here
   and applied twice (this cycle for schemas; cnos#376 Sub 1 for skills).

### Negative leverage

- Three CUE package directories under `schemas/` instead of one. Mitigated
  by per-package README files and this section.
- A repo-level CUE module file (`cue.mod/module.cue`) is required so
  `schemas/cds/` and `schemas/cdr/` can `import "cnos.dev/cnos/schemas/cdd"`.
  Mitigated: a 2-line file; does not affect the existing
  `tools/validate-skill-frontmatter.sh` invocations.
- Phase 3's V acquires a small dispatch layer reading `protocol_id`.
  Mitigated: one switch statement; the alternative (field-shape
  inference) is structurally worse.

### Design source

The design rationale is sourced from:

- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`](../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md)
  §"Schema direction" — names the split into `schemas/cdd/` +
  `schemas/cds/` + `schemas/cdr/` as the preferred shape; also names the
  adapter-boundary variant as acceptable if it preserves the same
  invariant ("Generic CDD validates the cell shape. Domain schemas
  validate domain evidence. V dispatches by declared protocol."). The (a)
  decision recorded here is consistent with the essay's preferred shape.
- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`](../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md)
  §2 "Use CUE for typed trust surfaces" — distinguishes CUE structural
  validation (the schemas here) from runtime evidence dereference (V's
  job in Phase 3). The split preserves the distinction.
- [`COHERENCE-CELL-NORMAL-FORM.md`](../../src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md)
  §Kernel — substrate-independent role signatures; generic kernel must
  not name domain-specific evidence fields. The cycle/388 split aligns
  `schemas/cdd/` with this discipline.
- [`ROLES.md §4a`](../../ROLES.md#4a-persona-operator-protocol-project-receipt--the-five-layer-enforcement-chain)
  — five-layer enforcement chain: protocol overlay (layer 3) is where
  domain evidence refs live. §4a.3 names CDS and CDR receipt sketches as
  the two required field sets per domain.
- [#388](https://github.com/usurobor/cnos/issues/388) issue body §"Architecture
  choice" — the L7 cdd/design analysis that recorded this decision
  pre-cycle.

### protocol_id dispatch convention

```text
cnos.cdd.cds.receipt.v1   → schemas/cds/#CDSReceipt
cnos.cdd.cdr.receipt.v1   → schemas/cdr/#CDRReceipt
cnos.cdd.<protocol>.receipt.v1 → schemas/<protocol>/#<Protocol>Receipt
```

Generic-only fixtures (those validating against `#Receipt` directly,
without any domain overlay) declare `protocol_id: "cnos.cdd.generic.receipt.v1"`.
Phase 3's V interprets this as "validate against the generic kernel only;
do not apply any domain-overlay requirements." This is useful for
bootstrap cycles and stub protocols.

### Cross-reference

- cnos#376 AC7 — same architectural decision-class for skills /
  role-overlays. The (a) split-vs-adapter reading recorded here is
  inherited by Sub 1.

## Scope-Lift Invariant

A coherence cell's parent-facing receipt is not just a record of what
happened inside the cell — it is the matter the **parent scope** reads
when the cell crosses the trust boundary. The recursion equation in
[`COHERENCE-CELL.md`](../../src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md)
§Recursion Equation states this in compressed form. This section names
the **three projections** the schemas in this directory preserve under
scope-lift. These are projections — distinct surfaces at scope `n` that
map onto distinct surfaces at scope `n+1` — not a flat renaming of roles
inside a single cell.

### Projection 1 — Closed cell → α-matter

> A **closed α/β/γ cell at scope `n`** is **α-matter at scope `n+1`**.

The cell at scope `n` runs its α/β/γ metabolism — α produces, β
discriminates, γ closes. The closed-cell record (issue, diff, the
per-cycle markdown files, the typed receipt) is what becomes available
to scope `n+1` as input material — the substrate scope-`n+1`'s α reasons
over.

Schema consequence: `#Contract` (`cnos.cdd.contract.v1`) types the
accountability surface scope `n+1` checks the receipt against. `#Receipt`
(`cnos.cdd.receipt.v1`) types the parent-facing artifact itself —
including the typed evidence refs scope `n+1`'s `V` resolves to verify
the receipt's claims. Without these typed surfaces, scope `n+1` would
have to re-derive structure from the in-cell narrative each time —
exactly what the receipt exists to fix.

### Projection 2 — δₙ BoundaryDecision → βₙ₊₁-like discrimination

> The **δₙ `#BoundaryDecision`** (this scope's boundary action) **projects
> to a βₙ₊₁-like discrimination at scope `n+1`** (the parent's review of
> incoming matter).

When δ at scope `n` records `action: accept | release | reject |
repair_dispatch | override` plus a `#Transmissibility`, scope `n+1` reads
that signal as the **discrimination input** for its own β-like layer at
the parent scope. A `transmissibility: not_transmissible` receipt cannot
be used as scope-`n+1` α-matter; a `transmissibility: degraded` receipt
can be used but only under explicit acknowledgment of the degradation; a
`transmissibility: accepted` receipt is the clean-closure case scope
`n+1` inherits without further discrimination work on this cell.

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

> An **εₙ receipt-stream observation** (the protocol-gap signal across
> many cells at scope `n`) **projects to a γₙ₊₁-like
> coordination/evolution at scope `n+1`** (parent-scope protocol
> authoring).

ε at scope `n` reads the receipt stream — many cells' `#ProtocolGapRef`
entries — and surfaces patterns that require protocol patches. Those
patches are scope-`n+1` γ-like work: coordinating which protocol gaps to
close, sequencing the fixes, evolving the rules the next round of
scope-`n` cells will run under. ε's signal at scope `n` is therefore the
**input** to the parent scope's coordination layer at scope `n+1`.

Schema consequence: `#Receipt` requires `protocol_gap_count: int & >= 0`
and `protocol_gap_refs: [...#ProtocolGapRef]`, with the consistency
constraint `protocol_gap_count == len(protocol_gap_refs)`. Both fields
are required in `v1` (not deferred to `v2`) so Phase 6's ε relocation
reads them without forcing a schema bump (see
[#366](https://github.com/usurobor/cnos/issues/366) Phase 6 /
[#369](https://github.com/usurobor/cnos/issues/369) `## Impact`).
`#ProtocolGapRef` is a structured object (`id`, `source` enum, `ref`)
rather than an opaque string, so scope-`n+1` ε-relocation reads typed
observations, not freeform notes.

### Projection-under-scope-lift, not role-renaming

The three projections above are **distinct surfaces at each scope**, not
a flat aliasing of roles. The δ at scope `n` does not *become* the β at
scope `n+1`. The δ at scope `n` performs δ work (boundary action on its
own receipt); the β at scope `n+1` performs β work (discrimination on
incoming matter); the projection maps the **output** of the former to a
**typed input** for the latter. The same applies for
closed-cell→α-matter and ε-stream→γ-coordination.

Reading the projection as flat role-renaming collapses two surfaces that
must stay distinct. The schemas keep them distinct: a `#Receipt` is a
typed artifact (not "a β-review at scope n+1"); a `#BoundaryDecision` is
a boundary action (not "a discrimination verdict from scope n+1's β"); a
`#ProtocolGapRef` is an observation (not "a γ next-commitment from scope
n+1"). The projection is what crosses the scope-lift; the surfaces it
lifts between remain typed at each end.

## How to run

You need [CUE](https://cuelang.org/) v0.10+ on your `PATH`. From the
repo root:

### Schema-only compile (generic kernel)

Checks that the generic schemas parse, type-check, and unify without
concrete values (`-c=false` allows incomplete schema values). Use this
when editing the kernel schemas themselves:

```
cue vet -c=false \
  schemas/cdd/contract.cue \
  schemas/cdd/boundary_decision.cue \
  schemas/cdd/receipt.cue
```

### Fixture validation (generic kernel)

Validates a fixture (a concrete receipt YAML) against `#Receipt` with all
three generic schemas in scope. `-c` requires concrete values for every
required field; `-d '#Receipt'` points the validator at the receipt
definition. Substitute `{fixture}` with the receipt fixture filename:

```
cue vet -c -d '#Receipt' \
  schemas/cdd/contract.cue \
  schemas/cdd/boundary_decision.cue \
  schemas/cdd/receipt.cue \
  schemas/cdd/fixtures/{fixture}.yaml
```

The valid generic fixture `valid-generic-receipt.yaml` succeeds; each of
the three named invalid fixtures fails for its named load-bearing reason.

### Fixture validation (CDS / CDR)

Domain fixtures validate against their domain schema via the CUE module
import (no need to pass the generic files explicitly; CUE resolves the
`import "cnos.dev/cnos/schemas/cdd"` automatically once `cue.mod/` is
in place):

```
cue vet -c -d '#CDSReceipt' \
  schemas/cds/receipt.cue \
  schemas/cds/fixtures/valid-receipt.yaml

cue vet -c -d '#CDRReceipt' \
  schemas/cdr/receipt.cue \
  schemas/cdr/fixtures/valid-cdr-receipt.yaml
```

Phase 3 of [#366](https://github.com/usurobor/cnos/issues/366) will wrap
the same invocations as the operator-facing `cn-cdd-verify` command's
structural reader, dispatching by `protocol_id`; Phase 3 inherits this
invocation pattern, not a re-derived one.

## What this directory does NOT do

This directory is **schema-only**. It does not:

- **Implement `V`.** The validator predicate (V at
  [RECEIPT-VALIDATION.md](../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md)
  §Validation Interface) is Phase 3 of [#366](https://github.com/usurobor/cnos/issues/366).
  These schemas are V's *input contract*, not V itself.
- **Type the full evidence graph.** `#Receipt` types evidence *refs*
  (string references to artifacts and SHAs) so Phase 3 has a pinned
  input shape. The graph's own schema (`evidence_graph.cue`) is Phase 3
  territory.
- **Carry CDS- or CDR-specific evidence-ref requirements.** That is
  what `schemas/cds/` and `schemas/cdr/` are for. The generic kernel
  declares `evidence_refs: [string]: #EvidenceRefValue` as a typed open
  primitive; domain packages add required keys via CUE unification.
- **Modify `cn-cdd-verify`.** The current artifact-presence checker at
  `src/packages/cnos.cdd/commands/cdd-verify/` is unchanged this cycle.
  Phase 3 refactors it into V's operator-facing wrapper.
- **Edit `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, or `CDD.md`.**
  The doctrine and design surfaces remain authoritative; the schemas
  reference them via header comment.
- **Pin a protocol version.** `v1` refers to each schema artifact's
  version (`cnos.cdd.contract.v1`, `cnos.cdd.receipt.v1`,
  `cnos.cdd.boundary_decision.v1`), not to a CDD protocol version. The
  CDD protocol is unversioned at this layer.

## Related issues and surfaces

- [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5; this
  cycle. Established the generic/domain split.
- [#369](https://github.com/usurobor/cnos/issues/369) — Phase 2; the
  unified schemas this cycle refactored.
- [#366](https://github.com/usurobor/cnos/issues/366) — coherence-cell
  executability roadmap; this directory is its Phase 2.5 deliverable
  (post-cycle/388) and Phase 3's input contract.
- [#367](https://github.com/usurobor/cnos/issues/367) — Phase 1; froze
  the parent-facing validator interface in
  [`RECEIPT-VALIDATION.md`](../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md).
- [#364](https://github.com/usurobor/cnos/issues/364) (closed) —
  doctrine;
  [`COHERENCE-CELL.md`](../../src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md).
- [#376](https://github.com/usurobor/cnos/issues/376) — CDR bootstrap;
  AC7 architectural-choice decision inherited from this cycle.
- [`schemas/skill.cue`](../skill.cue) + [`schemas/README.md`](../README.md) +
  [`schemas/fixtures/skill-frontmatter/`](../fixtures/skill-frontmatter/) —
  schema-layer pattern precedent: CUE owns shape, README owns prose,
  fixtures live under `schemas/{subsystem}/fixtures/`.
- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`](../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md) —
  L7 design essay; §"Schema direction" is the design source for this
  cycle's architectural choice.
