<!-- sections: [problem, constraints, challenged-assumption, impact-graph, proposal, alternatives, leverage, negative-leverage, non-goals, file-changes, acs] -->
<!-- completed: [problem, constraints, challenged-assumption, impact-graph, proposal, alternatives, leverage, negative-leverage, non-goals, file-changes, acs] -->

# Design Notes — Cycle #388

**Issue:** [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5 (#366): split generic schemas/cdd/ from CDS/CDR domain evidence
**Mode:** design-and-build (γ+α+β collapsed on δ)
**Active Skills:** cdd/design, cdd/issue/proof, cnos.core/skills/skill, cnos.eng/skills/eng/code
**Engineering Level:** L7 (changes the system boundary — generic kernel vs domain overlay — so future work becomes simpler)

---

## Problem

The Phase 2 schema (#369) ships `schemas/cdd/receipt.cue` with the following required fields in `#Receipt`:

```
self_coherence_ref, beta_review_ref, alpha_closeout_ref, beta_closeout_ref,
gamma_closeout_ref, diff_ref, ci_refs? (optional)
```

These names are the CDS process-evidence shape (the closure-record set #369 itself shipped against). They are not generic. A CDR (research) receipt does not have a `self_coherence_ref` or `beta_review_ref` in the same sense — it has `claim_refs`, `data_refs`, `method_refs`, `result_refs` (per `ROLES.md §4a.3`).

If Phase 3's V validates against `#Receipt`-as-written, then:
- a CDR receipt would either have to fake the CDS field names (incoherent), or
- V would need a different validator for CDR (defeats Phase 3's purpose of a single V over a single typed surface).

The named incoherence: **the generic kernel schema (`schemas/cdd/`) requires domain-specific evidence fields by name**, conflating CCNF's substrate-independent kernel (per `COHERENCE-CELL-NORMAL-FORM.md §Kernel` and `CCNF-AND-TYPED-TRUST.md §3` "Put domain evidence below the generic kernel") with the CDS protocol overlay (per `ROLES.md §4a` layer 3).

Evidence:
- `grep -n "self_coherence_ref\|beta_review_ref\|alpha_closeout_ref\|beta_closeout_ref\|gamma_closeout_ref\|diff_ref\|ci_refs" schemas/cdd/receipt.cue` returns 7 hits at lines 80-86.
- `schemas/cdd/fixtures/valid-receipt.yaml` is a CDS-shaped receipt populating the closure-record fields; it would not validate a CDR research receipt that has `claim_refs`/`data_refs`/etc.
- cph's `.cdr/waves/<wave>/receipt.md` (cited in #388 §Impact as the empirical anchor) carries research evidence fields; it would fail `cue vet` against the current `schemas/cdd/receipt.cue`.

---

## Constraints

**What can't change:**
- Generic invariants from #369: transmissibility derivation if-chain (verdict × action → transmissibility); required `boundary_decision`; `protocol_gap_count == len(protocol_gap_refs)`; override-iff (`#BoundaryDecision`); `#ValidationVerdict` shape; `#Override.degraded_state: true`. All three load-bearing invalid fixtures (`invalid-fail-no-boundary-decision`, `invalid-gamma-preflight-authoritative`, `invalid-override-masks-verdict`) must continue to fail `cue vet`.
- `#Receipt` continues to require `boundary_decision`, `validation`, `transmissibility`, `protocol_gap_count`, `protocol_gap_refs` — these are generic (any protocol's receipt has them).
- CCNF kernel discipline: `schemas/cdd/` names only roles, verdicts, decisions, receipts, evidence, scopes. No substrate-specific names. ([source: `COHERENCE-CELL-NORMAL-FORM.md §Two-Layer Separation`])
- CUE/V boundary per essay §2: CUE validates structural invariants; V validates evidence dereference + provenance.

**Existing contracts:**
- `schemas/cdd/README.md §Scope-Lift Invariant` documents three projections; the kernel sections of that README must remain intact (or be moved into a deeper section while the new top-level reflects the three-package layout).
- `tools/validate-skill-frontmatter.sh` invokes `cue vet -d '#Skill' schemas/skill.cue ...` — must not regress.

**Abstraction level:**
- The split is at the CUE-package/directory level (per ROLES.md §4a layer 3). Same-file subtype-by-required-keys (option b) is a semantic split; same-pkg different-files (option c, not in #388) is a file-level split; option (a) is directory-level — matches the protocol-overlay layer.

**What governs on disagreement:**
- `COHERENCE-CELL-NORMAL-FORM.md` is the kernel doctrine (the *what*); `schemas/cdd/` is the realization layer typing the kernel (the *how-on-this-substrate*). Where this cycle and CCNF disagree on shape, CCNF wins.
- `RECEIPT-VALIDATION.md` is the parent-facing validator design (Phase 1); the schema's verdict/decision separation is anchored there.

---

## Challenged Assumption

**Challenged:** "The CDD coherence-cell receipt schema names CDS-process-evidence fields (`self_coherence_ref`, `beta_review_ref`, the three closeout refs) because that's what every CDD receipt has."

**Counter:** CDD is the **generic** kernel; CDS is the **software-process** realization. The kernel's receipt names *what* every protocol's receipt has (validation, boundary decision, transmissibility, protocol gaps, typed evidence refs as a primitive). The CDS protocol overlay names *which* evidence_refs CDS requires (self-coherence, β-review, the three closeouts, diff, optional ci). Conflating the two surfaces freezes CDS-process structure into CCNF.

This matches `ROLES.md §4a`'s five-layer chain: persona → operator contract → **protocol overlay** → project binding → receipt/validator. The protocol overlay is the layer that names domain evidence; receipt/validator V dispatches by protocol.

---

## Impact Graph

**Downstream consumers of `schemas/cdd/receipt.cue` (the refactor target):**
- `schemas/cdd/fixtures/valid-receipt.yaml` (CDS-shaped) → moves to `schemas/cds/fixtures/`.
- Three invalid fixtures → stay under `schemas/cdd/fixtures/` (generic constraints) but trim CDS-specific keys from their bodies since they exercise generic constraints, not CDS evidence presence.
- `schemas/cdd/README.md` (declaration of the three .cue files + the load-bearing scope-lift invariant) → rewritten to reflect three-package layout + architectural-choice section.
- Phase 3 V (deferred) → will read `protocol_id` to dispatch validation; this cycle introduces the field structurally.
- `cn-cdd-verify` command (deferred) → reads via Phase 3's V; not touched this cycle.

**Upstream producers:**
- γ's receipt-authoring (a future Phase) produces receipts that pin `protocol_id`. This cycle does not change γ skill; close-out can name a debt (γ skill should be updated to emit `protocol_id` when it ships the CDS receipt).

**Copies and embeddings:**
- `ROLES.md §4a.3` carries CDS and CDR receipt sketches (prose). The CDS sketch in §4a.3 uses `artifact_refs/test_refs/ci_refs/diff_ref/debt_refs` (forward-looking shape); the #369 valid-receipt fixture uses the closure-record shape (`self_coherence_ref` etc). This cycle's `#CDSReceipt` matches what #369 actually shipped (the closure-record shape). Resolution of the divergence (§4a.3 sketch vs realized schema) is named as known debt; a future cycle aligns them.
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §Schema direction` describes the preferred split — matches option (a). Essay open-question #1 is resolved; essay update is deferred (optional).
- `schemas/cdd/README.md §Scope-Lift Invariant` documents three projections that the cycle/388 refactor preserves; the section moves into the new README structure unchanged.

**Authority relationships:**
- Where the kernel (`COHERENCE-CELL-NORMAL-FORM.md`) and the schemas disagree, the kernel statement is load-bearing and the schemas align to it (per CCNF §Cell Outcomes / §Two-Layer Separation).
- Where the essay (`CCNF-AND-TYPED-TRUST.md`) and the schemas disagree on shape, the essay is design rationale (not doctrine) — the schemas are the realization peer; both must compose coherently. In this cycle the essay's preferred shape and the schemas' implementation match.
- Where `ROLES.md §4a.3` and the schemas disagree on CDS field naming, `ROLES.md` is the canonical source-of-truth for the discipline-receipt mapping; the schema authors what #369 shipped (a debt to align noted).

---

## Proposal

### Types (starting from the model)

**`schemas/cdd/` — generic kernel package (`package cdd`):**

```cue
#EvidenceRef: string  // typed evidence reference primitive

#Receipt: {
    protocol_id: string                // dispatch key for V (REQUIRED, new this cycle)
    validation: #ValidationVerdict     // V's output (UNCHANGED)
    boundary_decision: #BoundaryDecision  // δ's record (UNCHANGED)
    transmissibility: #Transmissibility   // derived; UNCHANGED
    protocol_gap_count: int & >=0      // UNCHANGED
    protocol_gap_refs: [...#ProtocolGapRef]  // UNCHANGED
    protocol_gap_count: len(protocol_gap_refs)  // UNCHANGED

    evidence_refs: [string]: #EvidenceRef  // typed open primitive (NEW shape)
    // No CDS-specific named fields. Domain packages add required keys via & unification.

    // Transmissibility derivation if-chain — UNCHANGED.
    if validation.verdict == "PASS" { ... }
    if validation.verdict == "FAIL" { ... }

    ...  // open struct
}
```

Removed from `#Receipt`: the CDS-specific named fields (`self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref`, `ci_refs?`) and the `evidence_root_ref` field (since it points at a process-evidence root that's CDS-shaped). `evidence_root_ref` migrates into the CDS `evidence_refs` map.

**`schemas/cds/` — software protocol overlay (`package cds`, imports `cdd`):**

```cue
package cds

import "cnos.dev/cnos/schemas/cdd"

#CDSReceipt: cdd.#Receipt & {
    protocol_id: "cnos.cdd.cds.receipt.v1"
    evidence_refs: {
        evidence_root: cdd.#EvidenceRef
        self_coherence: cdd.#EvidenceRef
        beta_review:    cdd.#EvidenceRef
        alpha_closeout: cdd.#EvidenceRef
        beta_closeout:  cdd.#EvidenceRef
        gamma_closeout: cdd.#EvidenceRef
        diff:           cdd.#EvidenceRef
        ci?: [...cdd.#EvidenceRef]
    }
}
```

The CDS receipt requires all six closure-record refs plus the diff plus optional ci. Field names are stripped of the `_ref` suffix because the map's keys already imply "this is a ref": `evidence_refs.self_coherence` is a `#EvidenceRef`. This is cleaner than `evidence_refs.self_coherence_ref` (double-naming).

**`schemas/cdr/` — research protocol overlay (`package cdr`, imports `cdd`):**

```cue
package cdr

import "cnos.dev/cnos/schemas/cdd"

#ClaimStatus: "observed" | "computed" | "inferred" | "hypothesized" | "indeterminate"

#Reproduction: {
    reviewer_rerun: bool          // β re-ran the producing commands
    command:        string         // the command run
    output_match:   bool           // observed output matches receipt-claimed output
    notes?:         string
}

#CDRReceipt: cdd.#Receipt & {
    protocol_id: "cnos.cdd.cdr.receipt.v1"
    evidence_refs: {
        claim:  [...cdd.#EvidenceRef]   // which claims this receipt asserts
        data:   [...cdd.#EvidenceRef]   // dataset / mount / manifest / checksum
        method: [...cdd.#EvidenceRef]   // script paths + commit SHA
        result: [...cdd.#EvidenceRef]   // output file paths
    }
    claim_status: #ClaimStatus
    limitations?: [...string]
    reproduction?: #Reproduction
}
```

Matches `ROLES.md §4a.3` CDR sketch. Field cardinalities are lists (a CDR receipt has many claims, data sources, methods, results — unlike CDS where each closure record is singular).

### Module file

A new `cue.mod/module.cue` at the repo root with:

```
module: "cnos.dev/cnos"
language: version: "v0.10.0"
```

This enables `import "cnos.dev/cnos/schemas/cdd"` from `schemas/cds/receipt.cue` and `schemas/cdr/receipt.cue`. Tested pre-cycle: the addition of `cue.mod/` does **not** break the existing `tools/validate-skill-frontmatter.sh --self-test` (verified empirically; all positive + negative skill-frontmatter fixtures behave as expected post-module-add).

### protocol_id dispatch convention

V (Phase 3) reads `protocol_id` from receipt frontmatter and dispatches to the schema package:

```
"cnos.cdd.cds.receipt.v1" → schemas/cds/#CDSReceipt
"cnos.cdd.cdr.receipt.v1" → schemas/cdr/#CDRReceipt
```

Future protocols (cdw, cda, etc.) follow the same convention: `cnos.cdd.<protocol>.receipt.v1` → `schemas/<protocol>/#<Protocol>Receipt`. The dispatch is mechanical because the schema-package lookup is by string match on `protocol_id`.

### Fixture re-categorisation

| Fixture (current location) | New location | Reason |
|---|---|---|
| `schemas/cdd/fixtures/valid-receipt.yaml` | `schemas/cds/fixtures/valid-receipt.yaml` | CDS-shaped (populates closure-record fields); add `protocol_id` and nest evidence refs under `evidence_refs:` map |
| `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml` | unchanged location | Exercises generic constraint (`boundary_decision` required); trim CDS-specific evidence fields |
| `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml` | unchanged location | Generic constraint; trim CDS-specific fields |
| `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml` | unchanged location | Generic constraint (transmissibility table); trim CDS-specific fields |
| **new:** `schemas/cdd/fixtures/valid-generic-receipt.yaml` | `schemas/cdd/fixtures/` | Proves `#Receipt` validates without CDS keys (just generic fields + minimal `evidence_refs`) |
| **new:** `schemas/cdr/fixtures/valid-cdr-receipt.yaml` | `schemas/cdr/fixtures/` | Proves `#CDRReceipt` validates a research receipt with `claim/data/method/result` refs |

### Architectural-choice section (AC1)

`schemas/cdd/README.md` acquires an "Architectural choice" section recording:
- Option (a) split chosen.
- Option (b) adapter rejected.
- Rationale (1)–(5) from issue body.
- Cross-reference to cnos#376 AC7 (Sub 1 inherits this decision).
- Citation to `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §Schema direction`.

---

## Alternatives Considered

| Option | Pros | Cons | Decision |
|---|---|---|---|
| **(a) Split into three CUE pkgs/dirs** | Required-key validation in CUE; clear boundary at protocol layer; future cdw/cda generalize by adding a directory; mechanical generic-vs-domain split | Three dirs instead of one; CUE module file needed at repo root | **chosen** |
| (b) Adapter (open typed map, required keys per subtype in CDS/CDR subtypes within one file) | One directory | Validation shifts from CUE to V runtime logic; required-key set varies by subtype within same file (semantic, not structural); harder for new protocol authors to discover required fields | rejected |
| (c) Multi-pkg in one dir (`cdd` + `cds` + `cdr` files in `schemas/cdd/`) | No new directory; no module file | Conflates files; CUE package per directory is the idiomatic boundary; future protocols still need separate files | rejected |

The decision was reached in the ε session 2026-05-19 L7 cdd/design analysis (per issue body); this cycle records the rationale durably and implements the split.

---

## Leverage

- **Phase 3 (V implementation) gains a clean target.** V dispatches by `protocol_id`; the schema package set is the lookup table.
- **CDR bootstrap (#376 Sub 1) inherits this decision** — option (a) applied to skills/role-overlays as well as schemas. Single decision, applied twice (cnos#376 AC7).
- **Future protocols (cdw, cda, etc.) generalize mechanically.** Add `schemas/cdw/`, `schemas/cda/` directories with their own receipt packages; the generic schema doesn't need to know.
- **`schemas/cdd/` becomes citable as the substrate-independent kernel realization.** Future doctrine/essay work can cite it as the realization peer of the CCNF kernel without inheriting CDS-specific content.

---

## Negative Leverage

- **Three directories** under `schemas/` instead of one. Slightly heavier discovery; offset by per-package README and the README's "Three-package layout" section.
- **CUE module file** (`cue.mod/module.cue`) added at the repo root. Standard CUE practice but is a new top-level file; offset by enabling clean `import` semantics and not breaking existing skill-frontmatter validation.
- **V (Phase 3) acquires a dispatch layer** reading `protocol_id`. Small (one switch statement); much smaller than V having to infer protocol from field-shape.

---

## Non-goals (V0.1 of the split)

- Do NOT implement V (Phase 3 of #366).
- Do NOT define the full CDR field set beyond ROLES §4a.3 sketch (cnos#376 Sub 1/Sub 2 owns refinement).
- Do NOT migrate historical CDS receipts in `.cdd/releases/` to new schema location (they validate against `schemas/cds/` as-is; their location does not change).
- Do NOT rename `cn-cdd-verify` (Phase 3 may).
- Do NOT add cdw / cda schema skeletons.
- Do NOT change the existing invalid fixtures' load-bearing content — only re-shape their bodies to use `evidence_refs:` map structure (trimming the CDS-specific named fields).
- Do NOT add a `cn schemas verify` CLI surface.
- Do NOT realign `ROLES.md §4a.3` CDS sketch with what #369 shipped (named debt).
- Do NOT update the essay (`CCNF-AND-TYPED-TRUST.md`) open-questions section (deferred to a follow-up doc cycle).

---

## File Changes

**Create:**
- `cue.mod/module.cue` — repo-level CUE module declaration.
- `schemas/cds/receipt.cue` — declares `#CDSReceipt`.
- `schemas/cds/README.md` — package-level README.
- `schemas/cds/fixtures/valid-receipt.yaml` — moved from cdd/ and reshaped (evidence_refs nested map + protocol_id).
- `schemas/cdr/receipt.cue` — declares `#CDRReceipt`, `#ClaimStatus`, `#Reproduction`.
- `schemas/cdr/README.md` — package-level README.
- `schemas/cdr/fixtures/valid-cdr-receipt.yaml` — minimal CDR receipt.
- `schemas/cdd/fixtures/valid-generic-receipt.yaml` — proves `#Receipt` validates without CDS keys.
- `.cdd/unreleased/388/{gamma-scaffold,design-notes,self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout}.md` — cycle evidence.

**Edit:**
- `schemas/cdd/receipt.cue` — remove CDS-specific required fields; add `#EvidenceRef`; add `protocol_id`; replace `evidence_root_ref`/CDS-fields with `evidence_refs: [string]: #EvidenceRef`.
- `schemas/cdd/README.md` — full rewrite: three-package layout, architectural-choice section (AC1), protocol_id dispatch convention, preserved scope-lift invariant section.
- `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml` — trim CDS-specific evidence fields; add `protocol_id`; nest remaining refs under `evidence_refs:` map.
- `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml` — same trimming.
- `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml` — same trimming.

**Delete:**
- `schemas/cdd/fixtures/valid-receipt.yaml` — moved to schemas/cds/fixtures/.

**Unchanged:**
- `schemas/cdd/contract.cue` — already generic.
- `schemas/cdd/boundary_decision.cue` — already generic.
- `schemas/skill.cue`, `schemas/skill-exceptions.json`, `schemas/README.md` (top-level), `tools/validate-skill-frontmatter.sh` — orthogonal to this cycle.

---

## Acceptance Criteria (AC1–AC7 from issue #388)

(Restated here for completeness; full oracles in `self-coherence.md §ACs`.)

- **AC1**: `schemas/cdd/README.md` has "Architectural choice" section naming option (a), rejecting (b), enumerating rationale (1)–(5), referencing cnos#376 AC7, citing the essay.
- **AC2**: `rg "self_coherence_ref|beta_review_ref|alpha_closeout_ref|beta_closeout_ref|gamma_closeout_ref|diff_ref|ci_refs" schemas/cdd/` returns 0 hits (or only in explanatory comments referencing schemas/cds/).
- **AC3**: `cue vet -c -d '#CDSReceipt' schemas/cds/receipt.cue schemas/cds/fixtures/valid-receipt.yaml` exits 0; CDS-required field omission fails.
- **AC4**: `cue vet -c -d '#CDRReceipt' schemas/cdr/receipt.cue schemas/cdr/fixtures/valid-cdr-receipt.yaml` exits 0; matches ROLES §4a.3.
- **AC5**: fixture corpus layout matches expected; each fixture validates (or fails) at its target schema.
- **AC6**: `cue vet` regression-free; all valid fixtures pass, all invalid fixtures fail with diagnostic.
- **AC7**: `rg "protocol_id" schemas/` returns matches in each schema; every receipt fixture pins `protocol_id`; `protocol_id: string` declared structurally in `#Receipt`.

---

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 0 Observe | issue #388 | — | Phase 2.5 of #366; option (a) decided pre-cycle; design rationale needs durable record |
| 1 Select | this artifact | — | Phase 2.5 schema split — single coherent surface |
| 4 Gap | this artifact | — | generic `schemas/cdd/` requires CDS-shaped fields by name (lines 80-86 of receipt.cue) |
| 5 Mode | this artifact | cdd/design, cdd/issue/proof, cnos.core/skills/skill, cnos.eng/skills/eng/code | design-and-build; γ+α+β-collapsed on δ; L7 |
| 6 Artifacts | gamma-scaffold.md + this design-notes.md + self-coherence.md (in progress) + the CUE/README/fixture refactor | — | gamma-scaffold pushed; design-notes records L7 rationale; build half follows |
