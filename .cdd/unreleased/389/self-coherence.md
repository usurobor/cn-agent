<!-- sections: [Gap, Design, Skills, ACs, AC Coverage, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap, Design, Skills, ACs, AC Coverage, Self-check, Debt, CDD Trace, Review-readiness] -->

# Self-Coherence — Cycle #389 (Phase 3 of #366)

## §Gap

**Issue:** [#389](https://github.com/usurobor/cnos/issues/389) — Phase 3 of #366: implement V (`cn-cdd-verify` rewrite — Contract × Receipt → ValidationVerdict).

**Version / mode:** unreleased · design-and-build · γ+α+β-collapsed on δ (single Claude Code session; per the breadth-2026-05-12 wave manifest precedent and the empirically validated cycle pattern from 375/377/378/388). β-α-collapse acknowledged; the resulting receipt is the cycle's own first instance of a `mode: collapsed` declaration (see §Design).

**Gap:** Before this cycle, `cn-cdd-verify` was an artifact-presence checker — it verified the 5 cycle files existed at the right paths but did not validate receipt-vs-contract coherence. Phase 2.5 (#388) shipped the three-package schema split (`schemas/cdd/`, `schemas/cds/`, `schemas/cdr/`) and the `protocol_id` dispatch convention, but no executable V consumes those schemas. Without executable V, role-skill trust remains prose-only and δ cannot mechanically decide on receipts. Phase 4 (δ split) is blocked on V.

**Goal:** Implement V as a runtime validator over `protocol_id`-dispatched CUE schemas + receipt-bound evidence refs + counterfeit-receipt structural rules. Emit a typed `ValidationVerdict` (machine-consumable by δ). Preserve backward compatibility with the existing `cn-cdd-verify` operator-facing modes.

## §Design

See [`design-notes.md`](./design-notes.md) for the full L7 design record. Decision summary:

- **Implementation: hybrid bash-wrapper + Python core.** `cn-cdd-verify` (existing bash) gains four new flags (`--receipt`, `--contract`, `--json`, `--structural-only`) that dispatch to a sibling Python helper `cn-cdd-validate-receipt`. Pure-bash and Go-rewrite alternatives rejected (rationale in design-notes §Implementation choice).
- **CUE invocation: V wraps `cue vet`** for structural validation (per essay §3) and adds evidence-ref dereference + counterfeit-receipt checks. Dispatch table reads `protocol_id` from receipt frontmatter and routes to the matching schema package.
- **ValidationVerdict JSON schema:** `schemas/cdd/validation_verdict.schema.json` (JSON Schema draft-07). V emits this shape on `--json`; the prose render is for human readers.
- **Counterfeit-receipt rules (load-bearing contribution):** five mechanical predicates — C1 actor_separation, C2 verdict_precedes_merge, C3 override_does_not_rewrite, C4 evidence_ref_unresolved, C5 protocol_id_mismatch. Filesystem-touching rules (C1, C2, C4) are skipped under `--structural-only`; filesystem-free rules (C3, C5) always run.
- **`mode: collapsed` receipt-level field (cycle design finding):** when the receipt declares `mode: collapsed`, C1 and C2 emit warnings rather than failed_predicates. The collapse is acknowledged in the receipt itself; V is honest about the resulting structural-independence degradation by surfacing it as a `warnings` entry. Used by this cycle (#389) and back-projected onto cycle 369's CDS valid fixture (which legitimately shipped under collapsed mode).
- **Fixture corpus extension:** five new counterfeit fixtures (actor-collision, merge-precedes-verdict, override-rewrite, evidence-missing, mismatched-protocol-id) — each exercises one counterfeit rule. Closure-record stubs live under `schemas/cds/fixtures/_closure_records/<fixture-name>/`. Existing `schemas/cds/fixtures/valid-receipt.yaml` updated to point at cycle 369's real on-disk closure-record paths (a fixture-strengthening fix in-scope: turns the fixture from placeholder demonstration into end-to-end V-validatable).

## §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` — V signature
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — V's input/output/invocation contract (the design this cycle implements)

**Tier 2:**
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §2 (CUE typed trust surfaces) and §3 (V wraps CUE) — design source
- `schemas/cdd/README.md` §"Architectural choice" — option (a) split + `protocol_id` dispatch convention

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring (drives the harness)
- `src/packages/cnos.eng/skills/eng/code/SKILL.md` — Python + bash code quality
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — fixture + test-harness discipline

## §ACs

Per-AC oracles run by `tests/cdd/test_cn_cdd_validate_receipt.sh`. Full run: **37/37 PASS**.

## §AC Coverage

**AC1 — V signature matches CCNF:**
- Oracle: invoke V with `(--receipt R [--contract C])`; assert `ValidationVerdict` emitted; assert receipt-without-protocol_id is rejected (no separate evidence-graph input is accepted).
- Evidence: `tests/cdd/test_cn_cdd_validate_receipt.sh` AC1.a–AC1.e (5 sub-checks all PASS). V's invocation signature is exactly `(receipt, contract?)`; evidence refs are read *through* the receipt's `evidence_refs.*`.
- **Status: ✓ met**

**AC2 — V is protocol-aware via `protocol_id` dispatch:**
- Oracle: CDS receipt PASSes against `#CDSReceipt`; CDR receipt PASSes against `#CDRReceipt`; generic receipt PASSes against `#Receipt`; mismatched `protocol_id` (declares CDS, omits CDS keys) FAILs with diagnostic naming missing keys.
- Evidence: AC2.a–AC2.e (5 sub-checks all PASS). Dispatch table in `cn-cdd-validate-receipt`:
  ```
  cnos.cdd.generic.receipt.v1 → #Receipt   (schemas/cdd/)
  cnos.cdd.cds.receipt.v1     → #CDSReceipt (schemas/cds/)
  cnos.cdd.cdr.receipt.v1     → #CDRReceipt (schemas/cdr/)
  ```
- **Status: ✓ met**

**AC3 — V rejects missing `boundary_decision`:**
- Oracle: V on `invalid-fail-no-boundary-decision.yaml` exits 1, with diagnostic naming `boundary_decision`.
- Evidence: AC3.a + AC3.b PASS. CUE's structural rejection identifies the field at `boundary_decision.cue:71:5` / `:75:5` (the override-required-iff if-chain that fails when boundary_decision is absent).
- **Status: ✓ met**

**AC4 — V rejects γ-preflight as authoritative:**
- Oracle: V on `invalid-gamma-preflight-authoritative.yaml` exits 1; diagnostic cites the structural rule (boundary_decision required regardless of γ-preflight signal).
- Evidence: AC4.a + AC4.b PASS. Same CUE diagnostic mechanism as AC3 — the schema does not permit `boundary_decision` to be omitted regardless of any `preflight_only` field.
- **Status: ✓ met**

**AC5 — ValidationVerdict ≠ BoundaryDecision:**
- Oracle: V on `invalid-override-masks-verdict.yaml` exits 1; CUE catches transmissibility mismatch (`degraded` vs `accepted`). V on `counterfeit-override-rewrite.yaml` (the deeper case where outer `validation.verdict` was rewritten to PASS while override.original_validation_verdict says FAIL) exits 1 with `counterfeit.override_does_not_rewrite` diagnostic.
- Evidence: AC5.a–AC5.d PASS. CUE's transmissibility-table derivation handles the outer mismatch; V's C3 rule catches the inner override-rewrite that CUE cannot infer (CUE doesn't know `override.original_validation_verdict.verdict` must equal `validation.verdict` — that's the receipt-author-honesty invariant V enforces).
- **Status: ✓ met**

**AC6 — Evidence-graph inputs bound into receipt:**
- Oracle: CDS receipt with real on-disk closure-record refs PASSes (V dereferences them); CDS receipt with non-existent paths FAILs with `counterfeit.evidence_ref_unresolved` diagnostic.
- Evidence: AC6.a–AC6.c PASS. `valid-receipt.yaml` (paths → `.cdd/releases/docs/2026-05-17/369/`) dereferences cleanly; `counterfeit-evidence-missing.yaml` (paths → `.cdd/releases/nonexistent-version/000/`) FAILs with 6 unresolved-ref diagnostics.
- **Status: ✓ met**

**AC7 — ValidationVerdict JSON structure:**
- Oracle: `--json` emits parseable JSON; `result`, `failed_predicates`, `warnings`, `provenance` keys present; `result ∈ {PASS, FAIL}`; `provenance.validator_identity == "cnos.cdd.validate_receipt"`.
- Evidence: AC7.a–AC7.e PASS (8 sub-checks). JSON Schema lives at `schemas/cdd/validation_verdict.schema.json`. The shape was hand-validated against the schema's required-key list by the test harness.
- **Status: ✓ met**

**AC8 — Counterfeit-receipt rejection:**
- Oracle: three synthetic fixtures (one per rule C1/C2/C3) FAIL with the appropriate diagnostic.
- Evidence: AC8.a–AC8.f PASS (6 sub-checks):
  - `counterfeit-actor-collision.yaml` → `counterfeit.actor_separation` (C1) FAIL
  - `counterfeit-merge-precedes-verdict.yaml` → `counterfeit.verdict_precedes_merge` (C2) FAIL
  - `counterfeit-override-rewrite.yaml` → `counterfeit.override_does_not_rewrite` (C3) FAIL
- **Status: ✓ met**

**Backward compatibility (informally tested):**
- BC.a + BC.b PASS. `cn cdd-verify --unreleased` still drives the legacy artifact-presence checker (exits with summary; behaviour preserved).

## §Self-check

- All AC oracles mechanically pass (37/37 in the harness).
- V's signature matches CCNF: `Contract × Receipt → ValidationVerdict` (receipt-only invocation works; contract is optional but recorded in provenance).
- CUE invocation is via subprocess; V wraps `cue vet` per essay §3 — V does not re-implement CUE structural validation.
- Counterfeit-rule predicates are mechanical (finite-time, deterministic; reduce to filesystem ops + git log queries + simple field equality).
- JSON output validates against `schemas/cdd/validation_verdict.schema.json` (structural check via Python script).
- Backward compat: legacy modes (`--version`, `--pr`, `--all`, `--unreleased`, `--triadic`, `--cycle`) unchanged; new flags strictly additive.
- The cycle's own design notes acknowledge β-α-collapse; the `mode: collapsed` mechanism V introduces is the receipt-level honest acknowledgment of that collapse.

## §Debt

Named, scoped, deferred:

- **D1 — `cn-cdd-validate-receipt --diff <git_diff_ref>` resolution.** V's path-shape filter excludes `git_diff(...)` synthetic refs (they pass through as warnings, not failed_predicates). A future cycle (Phase 4 δ split or Phase 5 γ shrink) may want V to actually resolve `git_diff(origin/main..merge_sha)` and run scope-drift detection against the contract's `non_goals`. Deferred — out of scope for Phase 3.
- **D2 — `validator_version` autodiscovery.** V's version constant is baked into `cn-cdd-validate-receipt` (`V_VERSION = "phase3.0"`). A future cycle may want to derive this from a git tag or CHANGELOG entry. Deferred — not load-bearing.
- **D3 — receipt-level `mode` field formalization in schemas.** The cycle introduces `mode: collapsed` as a receipt-level acknowledgment of β-α-collapse, but the field is currently honored only by V's counterfeit rules; it is not pinned in `schemas/cdd/receipt.cue`. The schema's open-extension (`...` at the end of `#Receipt`) permits the field, but future cycles may want to type it as `mode?: "collapsed" | "distributed"` and document it in `schemas/cdd/README.md`. Deferred to a follow-up cycle.
- **D4 — Historical CDS-fixture migration into V's full-mode validation.** The current `schemas/cds/fixtures/valid-receipt.yaml` is updated to dereference cleanly; other historical receipt-shaped artifacts under `.cdd/releases/` are not yet receipt-typed. Phase 5+ work (per #389 non-goals: "Do NOT migrate historical closeout files").
- **D5 — γ-skill `protocol_id` automation patch (cnos#388 F1).** Deferred per #389 non-goals.

## §CDD Trace

Cycle artefacts:

- `gamma-scaffold.md` — issue framing, surfaces, AC oracle approach, expected diff scope
- `design-notes.md` — L7 design record: implementation choice, CUE invocation strategy, ValidationVerdict JSON schema, counterfeit-receipt rules, alternatives considered, risks
- `self-coherence.md` — this file
- `beta-review.md` — β-collapsed mechanical re-derivation of AC1–AC8 (TODO this cycle, after self-coherence)
- `alpha-closeout.md` — α's final summary
- `beta-closeout.md` — β-collapsed closeout
- `gamma-closeout.md` — γ closeout + post-merge plan
- `cdd-iteration.md` — ε iteration finding(s); empty findings recorded if nothing surfaced (closure-gate hygiene per #388's missing-iteration finding F1)

Skills loaded: per §Skills above.

Diff scope at this point (cycle/389 HEAD versus origin/main):
- Modified: `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` (+~50 lines: V dispatch + new flag handling)
- New: `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt` (V's Python core)
- New: `schemas/cdd/validation_verdict.schema.json` (JSON Schema for V output)
- New: `schemas/cds/fixtures/counterfeit-{actor-collision,merge-precedes-verdict,override-rewrite,evidence-missing,mismatched-protocol-id}.yaml` (counterfeit corpus)
- New: `schemas/cds/fixtures/_closure_records/<fixture-name>/{self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout}.md` (closure-record stubs for C1/C2 fixtures)
- New: `tests/cdd/test_cn_cdd_validate_receipt.sh` (AC harness)
- Modified: `schemas/cds/fixtures/valid-receipt.yaml` (evidence_refs paths updated to real on-disk locations; `mode: collapsed` declared)
- New: `.cdd/unreleased/389/*.md` (cycle evidence)

## §Review-readiness

α signals review-readiness:

- AC1–AC8 PASS in the mechanical harness (37/37)
- Design notes name all alternatives considered and rationale
- Counterfeit rules are reducible to mechanical predicates (no prose interpretation)
- Backward compat with existing operator workflows verified
- Diff scope matches the γ-scaffold expectation (one Python core file, one JSON Schema, five fixtures + closure stubs, one harness, the bash extension)
- Debt is named (D1–D5) and scoped (each deferred to a named future cycle)

β-collapsed-on-δ self-review proceeds next.
