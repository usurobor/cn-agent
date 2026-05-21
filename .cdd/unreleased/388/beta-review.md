<!-- sections: [Verdict, β-α-Collapse Acknowledgment, Contract Integrity, Issue Contract, Diff Context, Architecture, Findings] -->
<!-- completed: [Verdict, β-α-Collapse Acknowledgment, Contract Integrity, Issue Contract, Diff Context, Architecture, Findings] -->

**Verdict:** APPROVED

**Round:** 1
**Fixed this round:** n/a (first review)
**Branch CI state:** local — `cue vet` suite passes on all 6 fixtures with correct exit codes; `tools/validate-skill-frontmatter.sh --self-test` passes (cue.mod/ addition does not regress); no CI workflow runs cue vet against schemas/cdd/cds/cdr fixtures yet (Phase 3 of #366 will wire that).
**origin/main SHA at review time:** `7cdcd8d4e4d17340734d0c9df3786d499655c0b0` (matches session-start; no rebase required)
**Cycle branch head:** `7294b473d849740f9fbae252d267b3396fe27497`
**Implementation SHA (α review-readiness):** `7294b473`
**Merge instruction:** `git merge --no-ff cycle/388` into main with `Closes #388`

---

## §β-α-Collapse Acknowledgment

This cycle ran γ+α+β-collapsed on δ — single Claude Code session, no independent β process. Per breadth-2026-05-12 wave manifest precedent (validated empirically across cycles 375/377/378), this is acceptable for schema-refactor class **iff** β's gates are mechanical (grep / `cue vet` / file-existence) against canonical sources, not subjective judgment.

This cycle's β oracles are all mechanical:

| AC | Oracle | Mechanism |
|---|---|---|
| AC1 | "Architectural choice" section in schemas/cdd/README.md present, names option (a), rejects (b), enumerates rationale 1-5 | `grep` + re-read |
| AC2 | rg returns 0 hits for the seven CDS-specific field names in schemas/cdd/ | mechanical |
| AC3 | `cue vet -c -d "#CDSReceipt" schemas/cds/receipt.cue schemas/cds/fixtures/valid-receipt.yaml` exit 0; negative tests exit 1 | mechanical |
| AC4 | `cue vet -c -d "#CDRReceipt" schemas/cdr/receipt.cue schemas/cdr/fixtures/valid-cdr-receipt.yaml` exit 0; negative tests exit 1; field set matches ROLES §4a.3 | mechanical + re-read |
| AC5 | `ls schemas/{cdd,cds,cdr}/fixtures/` matches expected layout; each fixture's target schema validates | mechanical |
| AC6 | cue vet exit-code suite across all 6 fixtures | mechanical |
| AC7 | `rg "protocol_id" schemas/` matches every schema; `grep -l "^protocol_id:" schemas/{cdd,cds,cdr}/fixtures/*.yaml` matches all 6 fixtures | mechanical |

No AC requires β to exercise judgment α has not already exercised. The β-α-collapse here is structurally appropriate for the work class.

The two non-mechanical β actions are also performed below:
- Re-reading the architectural-choice section to confirm it cites the design-source documents (essay, COHERENCE-CELL-NORMAL-FORM, ROLES §4a.3, RECEIPT-VALIDATION).
- Re-reading the CDR schema to confirm it matches ROLES §4a.3 sketch (and the divergence between sketch and shipped is named as known debt).

---

## §Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Issue #388 OPEN; cycle in-progress; no false-closure signals |
| Canonical sources/paths verified | yes | All cited paths (`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`, `COHERENCE-CELL-NORMAL-FORM.md`, `RECEIPT-VALIDATION.md`, `ROLES.md §4a.3`) resolve on branch |
| Scope/non-goals consistent | yes | Phase 3 (V) explicitly deferred; CDR-skeleton-vs-frozen named; historical CDS receipts not migrated; cdw/cda not added |
| Constraint strata consistent | yes | ACs concrete and distinguishable; no overlap |
| Exceptions field-specific/reasoned | n/a | No exceptions claimed |
| Path resolution base explicit | yes | Base SHA `7cdcd8d4` in gamma-scaffold; HEAD = `7294b473` |
| Proof shape adequate | yes | All 7 ACs have mechanical oracle evidence |
| Cross-surface projections updated | yes | schemas/cdd/README.md "Architectural choice" cites cnos#376 AC7; gamma-closeout plans the #366 + #376 close-out comments |
| No witness theater / false closure | yes | `cue vet` actually runs and emits exit codes; rg actually counts hits; sanity check (skill-frontmatter self-test) confirms no regression |
| PR body matches branch files | n/a | CDD triadic protocol — no PR |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/388/gamma-scaffold.md` present on `origin/cycle/388` after push |

---

## §Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| 1 | `schemas/cdd/README.md` "Architectural choice" section names option (a), rejects (b), enumerates rationale 1-5, cross-references cnos#376 AC7, cites essay | yes | met | Section spans lines 67-160; all required content present; design-source citations include essay §"Schema direction", COHERENCE-CELL-NORMAL-FORM, ROLES §4a, RECEIPT-VALIDATION, issue body |
| 2 | `rg "self_coherence_ref\|beta_review_ref\|alpha_closeout_ref\|beta_closeout_ref\|gamma_closeout_ref\|diff_ref\|ci_refs" schemas/cdd/` returns 0 hits | yes | met | rg exit code 1 (no matches); the one pre-fix match in `contract.cue` was rewritten to generic language (`evidence_refs.diff`) |
| 3 | `schemas/cds/receipt.cue` exists, declares `#CDSReceipt: cdd.#Receipt & {...}` with closure-record evidence keys, pins `protocol_id` | yes | met | File exists; cue vet passes; missing-key fixture fails with `evidence_refs.self_coherence: incomplete value string`; wrong-protocol_id fixture fails with `conflicting values` |
| 4 | `schemas/cdr/receipt.cue` exists, declares `#CDRReceipt: cdd.#Receipt & {...}` with claim/data/method/result lists + claim_status + optional limitations + optional reproduction per ROLES §4a.3 | yes | met | File exists; cue vet passes; empty-data fixture fails with `incompatible list lengths (0 and 2)`; non-enum claim_status fails with `conflicting values` across all 5 enum members |
| 5 | Fixture corpus re-categorised: CDS-shaped valid moved to cds/; 3 invalid stay in cdd/; 1 new generic in cdd/; 1 new CDR in cdr/ | yes | met | Layout verified by `ls`; each fixture validates (or fails) at correct schema |
| 6 | `cue vet` regression-free; all valid fixtures pass; all invalid fixtures fail with diagnostic | yes | met | Full suite results: 3 valid pass (generic, cds, cdr), 3 invalid fail with appropriate diagnostics; skill-frontmatter self-test unaffected by `cue.mod/` addition |
| 7 | `protocol_id` declared structurally in `#Receipt`; pinned in `#CDSReceipt`, `#CDRReceipt`; populated in every fixture | yes | met | `protocol_id: string` is a required top-level field of `#Receipt` (line 71); pinned in CDS line 33, CDR line 56; all 6 fixtures declare `protocol_id:` |

### Named Doc Updates

- `schemas/cdd/README.md` — full rewrite reflecting three-package layout and architectural-choice decision-record.
- `schemas/cds/README.md` — new; describes CDS package; names §4a.3 vs realized-shape debt.
- `schemas/cdr/README.md` — new; describes CDR package; names "skeleton — not frozen" debt for cnos#376 Sub 1/Sub 2.

### Cross-references for Phase 2.5 close-out

- `schemas/cdd/README.md §"Architectural choice"` cross-references cnos#376 AC7 (which will inherit this decision in Sub 1).
- `schemas/cdd/README.md §"Cross-reference"` and `§"Related issues and surfaces"` cite cnos#376 and the essay.
- `gamma-closeout.md` plans the post-merge comments on cnos#366 (Phase 2.5 shipped, Phase 3 unblocked) and cnos#376 (architectural-choice decision inherited).

---

## §Diff Context

```
.cdd/unreleased/388/design-notes.md       | 318 +++  (new)
.cdd/unreleased/388/gamma-scaffold.md     | 225 +++  (new)
.cdd/unreleased/388/self-coherence.md     | 159 +++  (new)
cue.mod/module.cue                        |  16 +++  (new, repo-level)
schemas/cdd/README.md                     | 357 +/-  (rewritten)
schemas/cdd/contract.cue                  |   3 +/-  (comment language only)
schemas/cdd/fixtures/invalid-fail-*.yaml  |  26 +/-  (trim CDS keys, add protocol_id)
schemas/cdd/fixtures/invalid-gamma-*.yaml |  24 +/-  (same)
schemas/cdd/fixtures/invalid-override-*.yaml | 30 +/-  (same)
schemas/cdd/fixtures/valid-generic-receipt.yaml | 45 +  (new)
schemas/cdd/fixtures/valid-receipt.yaml   |  35 ---  (moved to cds/)
schemas/cdd/receipt.cue                   |  70 +/-  (remove CDS fields; add #EvidenceRef, #EvidenceRefValue, protocol_id)
schemas/cdr/README.md                     | 137 +++  (new)
schemas/cdr/fixtures/valid-cdr-receipt.yaml | 63 +++  (new)
schemas/cdr/receipt.cue                   |  74 +++  (new)
schemas/cds/README.md                     | 114 +++  (new)
schemas/cds/fixtures/valid-receipt.yaml   |  48 +++  (moved from cdd/, reshaped)
schemas/cds/receipt.cue                   |  57 +++  (new)
                                          ----------
                                          18 files, 1617 ins / 184 del
```

---

## §Architecture

- **Generic kernel discipline preserved.** `schemas/cdd/receipt.cue` names only roles, verdicts, decisions, receipts, evidence, scopes — no domain-specific names. Substrate-independence per `COHERENCE-CELL-NORMAL-FORM.md §Two-Layer Separation`.
- **CUE/V boundary preserved.** Per essay §2: CUE validates structure (the transmissibility table, the required-keys-per-protocol, the enum domains); V (deferred to Phase 3) validates evidence dereference + provenance. The cycle/388 schemas type structure only; no evidence-dereference logic creeps into CUE.
- **Verdict / decision separation preserved.** `#Receipt.validation: #ValidationVerdict` (V-emitted) and `#Receipt.boundary_decision: #BoundaryDecision` (δ-recorded) remain distinct structures in the generic kernel. The override-iff and degraded-state semantics are unchanged.
- **Five-layer chain (ROLES §4a) honored.** `schemas/cdd/` = the cell kernel (substrate-independent). `schemas/cds/`, `schemas/cdr/` = layer 3 protocol overlays. Project bindings (cph for CDR; the cnos repo itself for CDS) remain layer 4 — out of scope for this schema work. Receipt/V is layer 5 — Phase 3 inherits the cycle/388 schemas as input contract.
- **Future c-d-X generalization.** A new protocol (cdw, cda) adds `schemas/cdw/`, `schemas/cda/` with their own `#CDWReceipt`/`#CDAReceipt` definitions. The generic surface needs no change.

---

## §Findings

**No R1 findings — verdict APPROVE.** All ACs met on the first oracle pass.

**Non-blocking observations (for record, not fix-round):**

1. **CUE module file at repo root is a new top-level addition.** `cue.mod/module.cue` (16 lines) is the minimum CUE requires for cross-package import resolution. The file is documented inline. Sanity-checked: skill-frontmatter self-test passes unchanged.

2. **CDR `evidence_refs` list min-length-1 constraint uses `[_, ...#EvidenceRef]` pattern.** Confirmed empirically: empty list fails with `incompatible list lengths (0 and 2)` diagnostic. The pattern matches what ROLES §4a.3 names ("no data_refs → V rejects research receipt") even though §4a.3 doesn't specify the CUE encoding.

3. **AC2 oracle was initially false-positive on `contract.cue`'s comment line mentioning `diff_ref`.** Fixed pre-commit: comment rewritten to generic language (`evidence_refs.diff (per schemas/cds/)`). AC2 now returns 0 hits.

4. **`#Receipt.evidence_refs` widened from `[string]: #EvidenceRef` (string-only) to `[string]: #EvidenceRefValue` (union of single ref or list of refs).** This was a design finding during build half — the original primitive was too narrow for CDR's list-shaped fields. Widening preserves type discipline (both alternatives unify to `#EvidenceRef`). Documented in receipt.cue comments and design-notes.md §Proposal.

5. **Essay v0.1.0 open-question #1 resolution deferred to a follow-up doc cycle.** Folding the (a) decision back into the essay's open-questions section is out of #388's ACs and would expand scope. Named in self-coherence §Debt point 3 and gamma-closeout.

6. **γ skill does not yet emit `protocol_id` automatically.** Receipts authored by hand need `protocol_id:` in frontmatter; the current γ skill does not enforce this. Phase 3 or a γ-skill follow-up may address. Named in self-coherence §Debt point 4.

All six observations are debt to be tracked, not blockers.

---

**β signs off.** α may proceed to alpha-closeout; γ may merge.
