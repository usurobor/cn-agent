<!-- sections: [Gap, Design, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap, Design, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->

# Self-Coherence — Cycle #388

## §Gap

**Issue:** [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5 (#366): split generic `schemas/cdd/` from CDS/CDR domain evidence

**Version / mode:** unreleased · design-and-build · γ+α+β-collapsed on δ (single Claude Code session; per breadth-2026-05-12 wave manifest precedent, acceptable for schema-refactor class because β uses grep-mechanical + `cue vet`-mechanical gates against canonical sources)

**Gap:** `schemas/cdd/receipt.cue` (#369, merge `ff54f2a0`) required CDS-shaped named fields (`self_coherence_ref`, `beta_review_ref`, three `*_closeout_ref`, `diff_ref`, optional `ci_refs?`) as required fields of the generic `#Receipt`. The generic kernel must not require domain-specific evidence fields by name (`COHERENCE-CELL-NORMAL-FORM.md §Kernel`, `CCNF-AND-TYPED-TRUST.md §3` "Put domain evidence below the generic kernel", `ROLES.md §4a` five-layer chain). Without the split, Phase 3 V cannot dispatch coherently to CDS and CDR receipts against the same V.

**Goal:** Split `schemas/cdd/` (generic kernel: contract, receipt skeleton with typed `evidence_refs` primitive, boundary_decision, transmissibility, protocol-gap signal) + `schemas/cds/` (software protocol overlay: `#CDSReceipt: cdd.#Receipt & {...}` with closure-record refs) + `schemas/cdr/` (research protocol overlay: `#CDRReceipt: cdd.#Receipt & {...}` per ROLES §4a.3 sketch). Receipt frontmatter carries `protocol_id`; Phase 3's V dispatches on it.

## §Design

See [`design-notes.md`](./design-notes.md) for the full L7 design record (impact graph, alternatives considered, leverage / negative leverage). Decision summary:

- **Option (a) — split** chosen (three CUE packages, three directories under `schemas/`).
- **Option (b) — adapter** (open typed map with required keys per subtype) rejected for five rationale points (1)–(5) enumerated in `schemas/cdd/README.md §"Architectural choice"`.
- Decision applies to both this cycle (schemas) and cnos#376 AC7 (skills/role-overlays). Single decision, applied twice.

CUE inheritance idiom: `#CDSReceipt: cdd.#Receipt & { protocol_id: "cnos.cdd.cds.receipt.v1", evidence_refs: { self_coherence: cdd.#EvidenceRef, ... } }`. Requires a `cue.mod/module.cue` at repo root for cross-package `import` resolution — added this cycle.

Design-finding during build half: the original primitive `evidence_refs: [string]: #EvidenceRef` (string-valued map) was too narrow for CDR's list-shaped fields (`claim`, `data`, `method`, `result` are lists). Widened to `[string]: #EvidenceRefValue` where `#EvidenceRefValue: #EvidenceRef | [...#EvidenceRef]`. Preserves type discipline (both alternatives unify to `#EvidenceRef`).

## §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` — kernel doctrine (the *what* the schemas type)
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — Phase 1 V design (the verdict/decision separation preserved here)

**Tier 2:**
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` — L7 design essay; §"Schema direction" is the design source

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline; L7 output format
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring; β re-checks ACs as oracles
- `src/packages/cnos.core/skills/skill/SKILL.md` — README shape conventions
- `src/packages/cnos.eng/skills/eng/code/SKILL.md` — CUE refactor discipline

## §ACs

Per-AC oracles run against the branch `cycle/388` HEAD.

**AC1 — Architectural choice declared:**
- Oracle: `schemas/cdd/README.md` contains a section "Architectural choice — generic vs domain split" naming option (a), declaring (b) rejected, enumerating rationale (1)–(5), cross-referencing cnos#376 AC7, citing the essay (`CCNF-AND-TYPED-TRUST.md §"Schema direction"`).
- Evidence: `grep -n "Architectural choice" schemas/cdd/README.md` returns line 67. The section spans lines 67-160 and includes all required content (option (a) chosen, option (b) rejected, rationale 1-5, design source citations, protocol_id dispatch convention, cross-reference to cnos#376 AC7).
- **Status: ✓ met**

**AC2 — Generic `schemas/cdd/` carries no CDS-specific field by name:**
- Oracle: `rg "self_coherence_ref|beta_review_ref|alpha_closeout_ref|beta_closeout_ref|gamma_closeout_ref|diff_ref|ci_refs" schemas/cdd/` returns 0 hits.
- Evidence: `rg` exit code 1 (no matches). Verified post-`contract.cue` comment fix (the lone match — a benign comment in `contract.cue` describing the v1 derivation — was rewritten to use generic language pointing to `schemas/cds/`).
- **Status: ✓ met**

**AC3 — `schemas/cds/` exists and pins the CDS-specific evidence shape:**
- Oracle: `cue vet -c -d "#CDSReceipt" schemas/cds/receipt.cue schemas/cds/fixtures/valid-receipt.yaml` exits 0.
- Evidence: exit 0 confirmed; `#CDSReceipt` declared as `cdd.#Receipt & { protocol_id: "cnos.cdd.cds.receipt.v1", evidence_refs: { evidence_root, self_coherence, beta_review, alpha_closeout, beta_closeout, gamma_closeout, diff: #EvidenceRef, ci?: [...#EvidenceRef] } }`. Negative test (missing `self_coherence`): exit 1 with `evidence_refs.self_coherence: incomplete value string`. Negative test (wrong protocol_id): exit 1 with `conflicting values`.
- **Status: ✓ met**

**AC4 — `schemas/cdr/` exists and pins the CDR-specific evidence shape:**
- Oracle: `cue vet -c -d "#CDRReceipt" schemas/cdr/receipt.cue schemas/cdr/fixtures/valid-cdr-receipt.yaml` exits 0. Field set matches ROLES.md §4a.3 (claim/data/method/result refs as lists; claim_status enum; optional limitations + reproduction).
- Evidence: exit 0 confirmed; `#CDRReceipt` requires the four list-shaped evidence_refs keys with min-length-one constraint (`[_, ...#EvidenceRef]`), the `#ClaimStatus` enum (observed | computed | inferred | hypothesized | indeterminate), optional limitations, optional `#Reproduction` record. Negative test (empty data list): exit 1 with `incompatible list lengths (0 and 2)`. Negative test (non-enum claim_status `certain`): exit 1 with `conflicting values` against all 5 enum members.
- **Status: ✓ met**

**AC5 — Existing #369 fixture corpus re-categorised cleanly:**
- Oracle: `ls schemas/{cdd,cds,cdr}/fixtures/` matches expected layout; each fixture validates (or fails) at its target schema.
- Evidence:
  - `schemas/cdd/fixtures/`: `invalid-fail-no-boundary-decision.yaml`, `invalid-gamma-preflight-authoritative.yaml`, `invalid-override-masks-verdict.yaml`, `valid-generic-receipt.yaml` (new).
  - `schemas/cds/fixtures/`: `valid-receipt.yaml` (moved from cdd/, reshaped to nest evidence_refs + protocol_id added).
  - `schemas/cdr/fixtures/`: `valid-cdr-receipt.yaml` (new; minimal CDR receipt; reproduction record populated).
  - All four cdd/ fixtures behave correctly (one valid passes, three invalid fail with appropriate diagnostics). The CDS fixture passes against `#CDSReceipt`. The CDR fixture passes against `#CDRReceipt`.
- **Status: ✓ met**

**AC6 — `cue vet` regression-free:**
- Oracle: for every fixture in `schemas/{cdd,cds,cdr}/fixtures/`, `cue vet` against the target schema; valid-* exits 0, invalid-* exits non-0 with diagnostic.
- Evidence (full suite results):
  - `schemas/cdd/fixtures/valid-generic-receipt.yaml` → exit 0 ✓
  - `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml` → exit 1, diagnostic: `boundary_decision: unresolved disjunction "accept" | "release" | "reject" | "repair_dispatch" | "override" (type string)` (the if-chain in `#Receipt` references `boundary_decision.action`; absent field cascades to that diagnostic — semantically the missing-required-field signal) ✓
  - `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml` → exit 1, same diagnostic (same root cause: omitted boundary_decision) ✓
  - `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml` → exit 1, diagnostic: `transmissibility: conflicting values "degraded" and "accepted"` (the AC4 of #369 table forces FAIL × override → degraded; author asserted "accepted") ✓
  - `schemas/cds/fixtures/valid-receipt.yaml` → exit 0 ✓
  - `schemas/cdr/fixtures/valid-cdr-receipt.yaml` → exit 0 ✓
- Sanity check: `bash tools/validate-skill-frontmatter.sh --self-test` → all positive + negative fixtures behave as expected (the `cue.mod/module.cue` added this cycle does not break existing skill-frontmatter validation).
- **Status: ✓ met**

**AC7 — `protocol_id` enables V dispatch (Phase 3 prerequisite):**
- Oracle: `rg "protocol_id" schemas/` returns matches in every schema file declaring or pinning the field; every receipt fixture file has `protocol_id:` populated; `protocol_id: string` is required in `#Receipt`.
- Evidence:
  - `protocol_id` declared structurally in `schemas/cdd/receipt.cue` line 71 (required `protocol_id: string` at top of `#Receipt`).
  - `protocol_id` pinned in `schemas/cds/receipt.cue` line 33 (`protocol_id: "cnos.cdd.cds.receipt.v1"`).
  - `protocol_id` pinned in `schemas/cdr/receipt.cue` line 56 (`protocol_id: "cnos.cdd.cdr.receipt.v1"`).
  - Every fixture file declares `protocol_id:` (verified by `grep -l "^protocol_id:" schemas/{cdd,cds,cdr}/fixtures/*.yaml` returning all 6 files).
  - Dispatch convention documented in `schemas/cdd/README.md §"Architectural choice" > protocol_id dispatch convention`.
- **Status: ✓ met**

## §Self-check

**Did α's work push ambiguity onto β?** No. Every AC has concrete oracle evidence (grep counts, cue vet exit codes, diagnostic strings). The β-α-collapse is acknowledged below; β review applies the same oracles mechanically.

**Is every claim backed by evidence in the diff?**
- AC1: `schemas/cdd/README.md` §"Architectural choice" is in the diff.
- AC2: `rg` exit 1 (no matches in `schemas/cdd/`).
- AC3: `schemas/cds/receipt.cue` + fixture in the diff; cue vet exit 0.
- AC4: `schemas/cdr/receipt.cue` + fixture in the diff; cue vet exit 0.
- AC5: `git status` shows the moved fixture, the new fixtures, the trimmed invalid fixtures.
- AC6: cue vet suite output captured in evidence.
- AC7: `rg "protocol_id" schemas/` output captured; schema lines cited.

**Did the schema-finding (CDR list-shape) get handled coherently?**
Yes. During the build half, the original `evidence_refs: [string]: #EvidenceRef` failed when applied to CDR's list-shaped fields. The discovery was named and the primitive widened to `#EvidenceRefValue: #EvidenceRef | [...#EvidenceRef]`. Both CDS (string-valued keys) and CDR (list-valued keys) now unify against the generic primitive. The widening is documented in `schemas/cdd/receipt.cue` comments.

**Is the β-α-collapse defensible?**
Yes for this schema-refactor class. β's gates are mechanical:
- `rg` for forbidden CDS-specific names in `schemas/cdd/` (AC2 oracle).
- `cue vet` exit-code suite for every fixture (AC3, AC4, AC5, AC6 oracles).
- `grep -l "protocol_id:"` against fixtures (AC7 oracle).
- Re-read of `schemas/cdd/README.md` for the rationale enumeration (AC1 oracle).
None require subjective judgment that an independent β could provide above what α has already verified. The collapse is named in `beta-review.md`.

## §Debt

Known debt named in this cycle's close-out (per `cdd/CDD.md` §close-out conventions):

1. **`ROLES.md §4a.3` CDS sketch vs `schemas/cds/` realized shape divergence.** §4a.3 uses forward-looking names (`artifact_refs`, `test_refs`, `ci_refs`, `diff_ref`, `debt_refs`); cycle/388 implements what #369 shipped (closure-record shape: `self_coherence`, `beta_review`, the three closeouts, `diff`, optional `ci`). A future cycle may align §4a.3 and `schemas/cds/receipt.cue`. Named in `schemas/cds/README.md §"Forward-looking debt — §4a.3 sketch vs realized shape"`.

2. **CDR schema is a skeleton, not frozen.** cnos#376 Sub 1 / Sub 2 may refine field cardinalities (e.g. make `reproduction` required when `claim_status ∈ {observed, computed}`), add domain-specific fields (e.g. `protocol_compliance`, `review_independence`), or split `data` into `data_input` and `data_derived`. Named in `schemas/cdr/README.md §"Skeleton — not frozen"`.

3. **Essay v0.1.0 open-question #1 resolution not folded back into the essay.** `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Open questions"` lists Q1 as "Should Phase 2.5 prefer split schemas or a generic schema with domain adapters?" — resolved by this cycle as option (a). Essay update is **deferred to a follow-up doc cycle** (out of #388's ACs; expanding scope would inflate the cycle). The deferral is named in `gamma-closeout.md`.

4. **γ skill does not yet emit `protocol_id` automatically.** Receipts authored by hand or by future γ-skill tooling need `protocol_id: "cnos.cdd.cds.receipt.v1"` in their frontmatter; the current γ skill (cycle/385's surface) does not enforce this. Phase 3 or a small γ-skill follow-up may add the guidance.

## §CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 0 Observe | issue #388 | — | Phase 2.5 of #366; option (a) decided pre-cycle at ε session 2026-05-19 |
| 1 Select | gamma-scaffold.md | — | Phase 2.5 schema split — single coherent surface; design-and-build mode |
| 4 Gap | design-notes.md | cdd/design | generic `schemas/cdd/receipt.cue` requires CDS-shaped fields by name (lines 80-86); generic kernel must not name domain evidence |
| 5 Mode | design-notes.md + this self-coherence | cdd/design, cdd/issue/proof, cnos.core/skills/skill, cnos.eng/skills/eng/code | design-and-build; γ+α+β-collapsed on δ; L7 (changes the substrate boundary so future work becomes simpler) |
| 6 Artifacts | gamma-scaffold.md, design-notes.md, schemas/cdd/{README.md, receipt.cue, fixtures/*}, schemas/cds/{receipt.cue, README.md, fixtures/valid-receipt.yaml}, schemas/cdr/{receipt.cue, README.md, fixtures/valid-cdr-receipt.yaml}, cue.mod/module.cue, this self-coherence.md | cdd/design, cnos.eng/skills/eng/code | scaffold pushed; design notes record L7 rationale; build half landed all schemas + fixtures + READMEs; cue vet suite passes |
| 7 Self-check | this artifact | cdd/issue/proof | AC1–AC7 all met; β-α-collapse acknowledged; debt named |

## §Review-readiness

**Implementation SHA (α review-readiness):** [will be filled by the alpha-closeout commit SHA]

**Verification:**
- `rg "self_coherence_ref|beta_review_ref|alpha_closeout_ref|beta_closeout_ref|gamma_closeout_ref|diff_ref|ci_refs" schemas/cdd/` → exit 1 (no hits) ✓
- `cue vet -c -d "#Receipt" schemas/cdd/contract.cue schemas/cdd/boundary_decision.cue schemas/cdd/receipt.cue schemas/cdd/fixtures/valid-generic-receipt.yaml` → exit 0 ✓
- `cue vet -c -d "#Receipt" schemas/cdd/contract.cue schemas/cdd/boundary_decision.cue schemas/cdd/receipt.cue schemas/cdd/fixtures/invalid-*.yaml` → exit 1 for each ✓
- `cue vet -c -d "#CDSReceipt" schemas/cds/receipt.cue schemas/cds/fixtures/valid-receipt.yaml` → exit 0 ✓
- `cue vet -c -d "#CDRReceipt" schemas/cdr/receipt.cue schemas/cdr/fixtures/valid-cdr-receipt.yaml` → exit 0 ✓
- `bash tools/validate-skill-frontmatter.sh --self-test` → all positive + negative fixtures behave as expected ✓
- `grep -l "^protocol_id:" schemas/{cdd,cds,cdr}/fixtures/*.yaml` → 6 of 6 files match ✓

α signals review-readiness. β review proceeds against the same mechanical oracles.
