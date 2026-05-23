<!-- sections: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->
<!-- completed: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->

# γ Scaffold — Cycle #388

**Date:** 2026-05-21
**Issue:** [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5 (#366): split generic `schemas/cdd/` from CDS/CDR domain evidence
**Branch:** `cycle/388`
**Base SHA:** `7cdcd8d4e4d17340734d0c9df3786d499655c0b0` (origin/main at session start; tag `3.81.0`)
**γ identity:** gamma / gamma@cdd.cnos
**Dispatch config:** γ+α+β-collapsed-on-δ (single-session schema-refactor; per breadth-2026-05-12 wave manifest precedent, validated empirically by 375/377/378). Acceptable for schema-refactor class because β uses grep-mechanical + `cue vet`-mechanical gates against canonical sources.

---

## Issue

**Gap:** `schemas/cdd/receipt.cue` (#369, merge `ff54f2a0`) requires CDS-shaped named fields (`self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref`, optional `ci_refs?`) as part of the generic `#Receipt` definition. The generic kernel schema must not require domain-specific evidence fields by name (per `CCNF-AND-TYPED-TRUST.md §3` "Put domain evidence below the generic kernel" and `COHERENCE-CELL-NORMAL-FORM.md §Kernel`'s substrate-independence discipline).

**Goal:** Split into three CUE packages (option (a) per issue body L7 design): `schemas/cdd/` (generic kernel, with typed `evidence_refs` primitive), `schemas/cds/` (software-specific receipt extension), `schemas/cdr/` (research-specific receipt extension). Receipt frontmatter carries `protocol_id` (e.g. `cnos.cdd.cds.receipt.v1`) so Phase 3's V dispatches to the correct schema package.

**Priority:** P1 — Phase 3 (`cn-cdd-verify` rewrite, V implementation) gates on this. CDR bootstrap (#376) Sub 1 also gates on this architectural-choice decision.

**Work-shape:** small-to-medium cycle (7 ACs; single coherent surface: schema dir refactor + fixture re-categorisation + README updates).

**Mode:** design-and-build with γ+α+β collapsed. The architectural choice ((a) split vs (b) adapter) was decided pre-cycle at the ε session L7 cdd/design analysis (recorded in issue body). The design half this cycle records the rationale durably (in `schemas/README.md`) and authors the per-package CUE inheritance idiom; the build half executes the refactor.

---

## Surfaces γ expects α to touch

1. `schemas/cdd/receipt.cue` — refactor: remove CDS-specific named fields; declare `#Receipt` with `protocol_id`, generic evidence_refs primitive, preserve transmissibility-derivation if-chain.
2. `schemas/cdd/contract.cue` — preserved as-is (already generic).
3. `schemas/cdd/boundary_decision.cue` — preserved as-is (already generic).
4. `schemas/cdd/README.md` — replace #369 content with three-package layout description; add "Architectural choice" section (AC1) recording option (a) decision, citing the essay and rejecting (b).
5. `schemas/cdd/fixtures/valid-receipt.yaml` — moved to `schemas/cds/fixtures/valid-receipt.yaml` (it is CDS-shaped).
6. `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml` — stays under `schemas/cdd/fixtures/` (generic constraint); content trimmed to remove CDS-specific evidence_refs.
7. `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml` — stays under `schemas/cdd/fixtures/` (generic constraint); content trimmed.
8. `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml` — stays under `schemas/cdd/fixtures/` (generic constraint); content trimmed.
9. `schemas/cdd/fixtures/valid-generic-receipt.yaml` — **new**: proves `#Receipt` validates without CDS-specific keys.
10. `schemas/cds/receipt.cue` — **new**: declares `#CDSReceipt: #Receipt & {...}` with required CDS evidence fields, pins `protocol_id: "cnos.cdd.cds.receipt.v1"`.
11. `schemas/cds/README.md` — **new**: package-level description, link back to schemas/cdd/.
12. `schemas/cds/fixtures/valid-receipt.yaml` — moved from `schemas/cdd/fixtures/` (and protocol_id added).
13. `schemas/cdr/receipt.cue` — **new**: declares `#CDRReceipt: #Receipt & {...}` per ROLES §4a.3 sketch, pins `protocol_id: "cnos.cdd.cdr.receipt.v1"`.
14. `schemas/cdr/README.md` — **new**.
15. `schemas/cdr/fixtures/valid-cdr-receipt.yaml` — **new**: minimal CDR receipt; empirical anchor in cph receipt shape.
16. `.cdd/unreleased/388/{gamma-scaffold,design-notes,self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout}.md` — cycle evidence.

---

## Peer enumeration (§2.2a — grep evidence)

- "schemas/cdd/receipt.cue requires CDS-shaped fields by name" — confirmed: `grep -n "self_coherence_ref\|beta_review_ref\|alpha_closeout_ref\|beta_closeout_ref\|gamma_closeout_ref\|diff_ref\|ci_refs" schemas/cdd/receipt.cue` returns 7 hits at lines 80-86 (the required-fields block).
- "schemas/cds/ does not exist" — confirmed: `ls schemas/cds/` → does not exist.
- "schemas/cdr/ does not exist" — confirmed: `ls schemas/cdr/` → does not exist.
- "schemas/cdd/fixtures/valid-receipt.yaml is CDS-shaped" — confirmed: it populates `self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref`.
- "the three invalid fixtures exercise generic constraints (boundary decision, override polarity, γ-preflight authoritativeness)" — confirmed by reading their header comments; each exercises an invariant declared in `schemas/cdd/boundary_decision.cue` / `schemas/cdd/receipt.cue`'s structural if-chain.
- "ROLES.md §4a.3 CDR sketch names claim/data/method/result + claim_status + limitations + reproduction" — confirmed: ROLES.md lines 244-252.
- "essay §Schema direction names the same split (cdd/cds/cdr) as preferred" — confirmed: docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §Schema direction (lines 310-329).
- "cph empirical anchor exists" — referenced in issue body §Impact and §Active design constraints; cph repo not present in this worktree, treated as design-source by citation.

No claim of "X does not exist" asserted without ls/grep evidence.

---

## Scope boundary

**In scope (Phase 2.5):** schema refactor + fixture re-categorisation + README decision-record.

**Out of scope:**
- V implementation (Phase 3 of #366)
- cn-cdd-verify code changes (Phase 3)
- Full CDR field set beyond ROLES §4a.3 sketch (cnos#376 Sub 1/Sub 2 owns refinement)
- Migrating historical CDS receipts in `.cdd/releases/` to new schema location (they validate against `schemas/cds/` as-is)
- cdw / cda schema skeletons (not yet declared protocols)
- A `cn schemas verify` CLI surface (file separately if needed)

**Sister-cycle awareness:**
- cnos#366: this cycle marks Phase 2.5 shipped, Phase 3 unblocked. γ-closeout posts comment on #366.
- cnos#376 AC7: architectural-choice decision is inherited from this cycle's (a) decision. γ-closeout posts comment on #376.
- cnos#387 just shipped (no overlap; touches cn-cdd-verify CLI surface).
- Essay v0.1.0 (CCNF-AND-TYPED-TRUST.md, commit `7cdcd8d4`) open-question #1 is resolved by this cycle. Essay update is **optional**; default decision is to defer to a follow-up doc cycle (essay update is not in the issue ACs and would expand scope).

---

## AC oracle approach

α should verify each AC using the oracles embedded in the issue ACs:

**AC1 — Architectural choice declared:**
Oracle: `schemas/cdd/README.md` (NB: issue body names this as the canonical artifact; a new `schemas/README.md` section is acceptable too but not required) contains a section titled "Architectural choice" or "Generic vs domain split"; the section names option (a), declares (b) explicitly rejected, and enumerates rationale (1)–(5) from issue body; cross-reference to cnos#376 AC7.

**AC2 — Generic schemas/cdd/ carries no CDS-specific field by name:**
Oracle: `rg "self_coherence_ref|beta_review_ref|alpha_closeout_ref|beta_closeout_ref|gamma_closeout_ref|diff_ref|ci_refs" schemas/cdd/` returns 0 hits (or only matches in explicit "see schemas/cds/" comments).

**AC3 — schemas/cds/ exists and pins the CDS-specific evidence shape:**
Oracle: `cue vet -c -d '#CDSReceipt' schemas/cdd/*.cue schemas/cds/*.cue schemas/cds/fixtures/valid-receipt.yaml` exits 0. A negative fixture omitting a required CDS field would fail.

**AC4 — schemas/cdr/ exists and pins the CDR-specific evidence shape:**
Oracle: `cue vet -c -d '#CDRReceipt' schemas/cdd/*.cue schemas/cdr/*.cue schemas/cdr/fixtures/valid-cdr-receipt.yaml` exits 0. Field set matches ROLES.md §4a.3.

**AC5 — Existing #369 fixture corpus re-categorised cleanly:**
Oracle: `ls schemas/{cdd,cds,cdr}/fixtures/` matches the expected layout (CDS-shaped `valid-receipt.yaml` → `schemas/cds/fixtures/`; three invalid fixtures stay under `schemas/cdd/fixtures/`; new `valid-generic-receipt.yaml` added under `schemas/cdd/fixtures/`). Each fixture validates (or fails) against its target schema.

**AC6 — cue vet regression-free:**
Oracle: for f in schemas/{cdd,cds,cdr}/fixtures/*.yaml; do `cue vet` against the correct schema; valid-* exits 0, invalid-* exits non-0 with diagnostic.

**AC7 — protocol_id enables V dispatch (Phase 3 prerequisite):**
Oracle: `rg "protocol_id" schemas/` returns matches in each schema file declaring or pinning the field. Every receipt fixture file has `protocol_id:` populated. `protocol_id: string` is required in `#Receipt`; the CDS/CDR refinements pin specific values (`"cnos.cdd.cds.receipt.v1"`, `"cnos.cdd.cdr.receipt.v1"`).

**Validation tool:** `cue` v0.10.0 (installed at `/root/.local/bin/cue`). Pre-cycle baseline verified: `cue vet -c -d '#Receipt' schemas/cdd/contract.cue schemas/cdd/boundary_decision.cue schemas/cdd/receipt.cue schemas/cdd/fixtures/valid-receipt.yaml` exits 0; the three invalid fixtures fail.

---

## Design-source citations

The cycle's design half cites:

1. `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` (commit `7cdcd8d4`, v0.1.0):
   - §"Schema direction" — preferred split is `schemas/cdd/` + `schemas/cds/` + `schemas/cdr/` (matches option (a)).
   - §2 "Use CUE for typed trust surfaces" — distinguishes CUE structural validation from V evidence-dereference; this cycle preserves the distinction by keeping `evidence_refs` typed-but-open in `#Receipt` and pinning the domain-specific named fields in `schemas/cds/` / `schemas/cdr/`.
   - §"Open questions" #1: this cycle resolves the question; essay update is deferred to follow-up.

2. `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` §Kernel:
   - V signature `Contract × Receipt → ValidationVerdict` (step 4).
   - Evidence-binding rule: γ binds evidence into receipt as typed references; V dereferences.
   - Substrate-independence: kernel names only roles, verdicts, decisions, receipts, evidence, scopes — no substrate-specific names. The generic `schemas/cdd/` mirrors this discipline.

3. `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`:
   - Phase 1 design — V at δ-boundary; verdict ≠ decision separation.
   - This cycle preserves the structural separation: `#Receipt.validation: #ValidationVerdict` and `#Receipt.boundary_decision: #BoundaryDecision` remain in the generic schema.

4. `ROLES.md §4a` (especially §4a.3):
   - Five-layer enforcement chain — protocol overlay (layer 3) is where domain-specific evidence refs live.
   - §4a.3 CDS sketch: `artifact_refs`, `test_refs`, `ci_refs`, `diff_ref`, `debt_refs`.
   - §4a.3 CDR sketch: `claim_refs`, `data_refs`, `method_refs`, `result_refs`, `claim_status`, `limitations`, `reproduction`.

   **Note on §4a.3 CDS vs #369 receipt shape:** §4a.3 uses `artifact_refs/test_refs/ci_refs/diff_ref/debt_refs` (forward-looking sketch). The #369 `valid-receipt.yaml` populates `self_coherence_ref/beta_review_ref/{alpha,beta,gamma}_closeout_ref/diff_ref` (CDD-process-evidence shape, which is what #369 actually shipped). This cycle preserves what #369 actually shipped — `#CDSReceipt` requires the closure-record fields (`self_coherence`, `beta_review`, `alpha_closeout`, `beta_closeout`, `gamma_closeout`, `diff`, optional `ci`), naming them at the evidence_refs map under their type-stripped names. The §4a.3 forward-looking sketch is a debt (named in close-out): a future cycle may align §4a.3 and `schemas/cds/` if the field-set should converge.

5. `schemas/cdd/README.md` and `schemas/cdd/{contract,receipt,boundary_decision}.cue` (#369, merge `ff54f2a0`):
   - The Phase 2 deliverable being refactored. Preserve all generic invariants (transmissibility derivation, override iff, protocol_gap count==len(refs), required `boundary_decision`).

---

## CUE inheritance idiom (design half)

CUE's `&` (unification) operator composes definitions. The idiom for "CDS receipt is a generic Receipt with these extra required fields" is:

```cue
package cds

import "cnos.dev/schemas/cdd:cdd"

#CDSReceipt: cdd.#Receipt & {
    protocol_id: "cnos.cdd.cds.receipt.v1"
    evidence_refs: {
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

**Decision on import vs same-pkg with -d:** CUE supports both. For simplicity (no module file needed; matches existing schemas/cdd/ pattern), the implementation uses separate packages (`package cdd`, `package cds`, `package cdr`) and `cue vet` is invoked with multiple `.cue` files providing all needed definitions. This requires that `#EvidenceRef` and `#Receipt` are accessible via cross-file unification — CUE allows passing multiple files from different packages to `cue vet` as long as the `-d` definition path points to the right one. The simpler alternative (also viable): keep all three packages under a single CUE package `cdd` namespace by reusing `package cdd` — but that conflates the three protocols.

**Settled approach:** Each domain package (`cds`, `cdr`) re-states the generic `#Receipt` definition locally (or uses CUE's `_` placeholder), and the domain-specific definitions (`#CDSReceipt`, `#CDRReceipt`) wrap it with the required evidence_refs structure. To avoid duplication, the alpha build will use CUE's module/import facility if it works simply; if not, the build half pre-tests an alternative (e.g. file-level cross-package via `cue vet` taking files from both `schemas/cdd/` and `schemas/cds/` paths and using fully-qualified definition references).

**Pre-test:** α writes a tiny throwaway fixture before committing to the import approach; if CUE module resolution requires a `cue.mod/module.cue` file at repo root and one doesn't exist, α falls back to "domain pkg files literally include the generic field structure inline, with a comment 'mirrors schemas/cdd/receipt.cue #Receipt'." This is duplication but is the safe fallback. **α decides at build time; β audits the chosen idiom.**

---

## Expected diff scope

| Surface | Expected delta |
|---|---|
| `schemas/cdd/receipt.cue` | -10 lines (remove CDS named-field block) + 5 lines (add `protocol_id`, `#EvidenceRef`, `evidence_refs` open map) |
| `schemas/cdd/contract.cue` | 0 lines (unchanged) |
| `schemas/cdd/boundary_decision.cue` | 0 lines (unchanged) |
| `schemas/cdd/README.md` | full rewrite (~150-200 lines): three-package layout, architectural-choice section, dispatch convention |
| `schemas/cdd/fixtures/*` | 3 invalid fixtures trimmed (remove CDS-specific evidence refs); 1 new generic fixture |
| `schemas/cds/receipt.cue` | new (~50 lines) |
| `schemas/cds/README.md` | new (~50 lines) |
| `schemas/cds/fixtures/valid-receipt.yaml` | moved from cdd/fixtures + protocol_id added + nested under evidence_refs |
| `schemas/cdr/receipt.cue` | new (~70 lines) |
| `schemas/cdr/README.md` | new (~60 lines) |
| `schemas/cdr/fixtures/valid-cdr-receipt.yaml` | new (~30 lines) |
| `.cdd/unreleased/388/*.md` | cycle evidence (7 files, ~600 lines total) |

Total: ~600-800 line net change across ~15 files. Mostly CUE + YAML + Markdown.

---

## Dispatch configuration

**γ+α+β collapsed on δ.** Single Claude Code session. The β-α-collapse is acknowledged: α ≠ β within a session is structurally compromised but acceptable for schema-refactor class because β uses grep-mechanical + `cue vet`-mechanical gates against canonical sources (the existing #369 schemas, the issue body's ACs, the essay text). The beta-review.md will name this collapse explicitly and run mechanical AC oracles.

Pre-flight check result:
```
γ pre-dispatch check — 2026-05-21:
  origin/cycle/388 exists: YES (just created)
  .cdd/unreleased/388/gamma-scaffold.md exists locally: YES (this file)
  issue #388 is open: YES
  branch pre-flight: PASS (base SHA 7cdcd8d4 matches origin/main HEAD)
  peer enumeration: PASS — all referenced surfaces confirmed by grep/ls
  scope boundary: documented (Phase 2.5 only; Phase 3 deferred)
  cross-repo intake: n/a
  issue quality gate: PASS (formal ACs, design source cited, architectural decision recorded)
  cue tool available: YES (v0.10.0 at /root/.local/bin/cue)
  baseline cue vet: PASS (valid fixture passes; invalid fixtures fail)
  dispatch config: γ+α+β-collapsed on δ; small/medium refactor; mechanical β gates
  timeout budget: n/a (synchronous session)
```

---

## Tier 3 skills

Named explicitly:
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline; L7 output format informs `schemas/cdd/README.md` "Architectural choice" section
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring; β re-checks ACs as oracles
- `src/packages/cnos.core/skills/skill/SKILL.md` — for README shape
- `src/packages/cnos.eng/skills/eng/code/SKILL.md` — CUE refactor discipline; engineering quality
