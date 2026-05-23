<!-- sections: [Issue, Mode, Surfaces, AC oracle approach, Expected diff scope] -->

# γ-Scaffold — Cycle #389 (Phase 3 of #366)

## §Issue

[#389](https://github.com/usurobor/cnos/issues/389) — Phase 3 of #366: implement V (`cn-cdd-verify` rewrite — Contract × Receipt → ValidationVerdict).

Parent roadmap: #366 (Phase 3 gates Phase 4 — δ split).

Inherits from #388 (Phase 2.5 — three-package schema split with `protocol_id` dispatch) and #367 (Phase 1 — validation-surface design).

## §Mode

**design-and-build, γ+α+β-collapsed on δ.**

Per the breadth-2026-05-12 wave manifest precedent (375/377/378/388). β-α-collapse acknowledged. Acceptable for this cycle class because:

- AC1–AC8 are mechanical predicates with file-system + `cue vet` + JSON-schema oracles
- AC2 (multi-protocol dispatch) and AC8 (counterfeit rejection) are the load-bearing surfaces — β-mechanical fixture-based checks suffice
- No design-novel surface beyond what's in #367's `RECEIPT-VALIDATION.md`, #370's `COHERENCE-CELL-NORMAL-FORM.md`, #388's `schemas/cdd/README.md`, and the `CCNF-AND-TYPED-TRUST.md` essay
- Empirically validated as a viable mode for L7-shape-already-fixed implementation work

## §Surfaces

In-scope (cycle diff):

- `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` (bash) — extend with `--receipt <path>` and `--contract <path>` flags; preserve existing modes (`--version`, `--pr`, `--cycle`, `--all`, `--unreleased`, `--triadic`) untouched
- `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt` (new, Python) — V's core implementation: dispatch on `protocol_id`, invoke `cue vet`, dereference `evidence_refs`, apply counterfeit-receipt rules, emit `ValidationVerdict` JSON
- `schemas/cdd/validation_verdict.schema.json` (new) — JSON Schema for the emitted `ValidationVerdict` structure (machine-consumable by δ side)
- `schemas/cds/fixtures/counterfeit-actor-collision.yaml` (new)
- `schemas/cds/fixtures/counterfeit-merge-precedes-verdict.yaml` (new)
- `schemas/cds/fixtures/counterfeit-override-rewrite.yaml` (new)
- `schemas/cds/fixtures/counterfeit-evidence-missing.yaml` (new) — AC6 negative
- `schemas/cds/fixtures/counterfeit-mismatched-protocol-id.yaml` (new) — AC2 negative
- `tests/cdd/test_cn_cdd_validate_receipt.sh` (new) — driver invoking V against every fixture + asserting outcome
- `.cdd/unreleased/389/{gamma-scaffold,design-notes,self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout,cdd-iteration}.md` (cycle evidence)

Out of scope (this cycle):

- δ-skill split (Phase 4 — #366 next)
- γ-skill shrink (Phase 5)
- ε-skill relocation (Phase 6)
- `CDD.md` rewrite (Phase 7)
- γ-skill `protocol_id` automation patch (#388 F1; deferred per #389 non-goals)
- Historical closeout migration in `.cdd/releases/` (back-fill is separate work)
- New top-level commands (`cn schemas verify`, `cn validate`); Phase 3 extends `cn-cdd-verify` only

## §AC oracle approach

| AC | Oracle (mechanical) |
|----|---------------------|
| AC1 | Invoke `cn-cdd-verify --receipt R --contract C` against valid fixtures; assert `ValidationVerdict` emitted. Assert no separate `--evidence` input exists (refs are in receipt). |
| AC2 | Run V on `valid-generic-receipt.yaml` (PASS against `#Receipt`), `valid-receipt.yaml` (PASS against `#CDSReceipt`), `valid-cdr-receipt.yaml` (PASS against `#CDRReceipt`), `counterfeit-mismatched-protocol-id.yaml` (FAIL — declares `cnos.cdd.cds.receipt.v1` but missing CDS required keys). |
| AC3 | Run V on `invalid-fail-no-boundary-decision.yaml`; assert FAIL + `failed_predicates` cites `boundary_decision` absence. |
| AC4 | Run V on `invalid-gamma-preflight-authoritative.yaml`; assert FAIL + diagnostic cites γ-preflight non-authoritativeness (delegated to `cue vet` since missing `boundary_decision`). |
| AC5 | Run V on `invalid-override-masks-verdict.yaml`; assert FAIL + diagnostic shows transmissibility derivation mismatch (override does not rewrite verdict). |
| AC6 | Run V on `valid-receipt.yaml` (passes — evidence refs resolve), and on `counterfeit-evidence-missing.yaml` (FAIL — `evidence_refs.self_coherence` points at non-existent path). |
| AC7 | `cn-cdd-verify --receipt R --json` emits valid JSON; assert `jq` parses; assert structure has `{result, failed_predicates, warnings, provenance}` keys. |
| AC8 | Run V on three synthetic counterfeit fixtures (actor-collision, merge-precedes-verdict, override-rewrite); assert FAIL + diagnostic naming violated rule. |

## §Expected diff scope

```text
M  src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify          (+~80 lines: --receipt/--contract/--json flag handling)
A  src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt (~400 lines Python)
A  schemas/cdd/validation_verdict.schema.json                       (~60 lines JSON Schema)
A  schemas/cds/fixtures/counterfeit-actor-collision.yaml            (+ fixture closure-record files)
A  schemas/cds/fixtures/counterfeit-merge-precedes-verdict.yaml     (+ fixture closure-record files)
A  schemas/cds/fixtures/counterfeit-override-rewrite.yaml           (+ fixture closure-record files)
A  schemas/cds/fixtures/counterfeit-evidence-missing.yaml
A  schemas/cds/fixtures/counterfeit-mismatched-protocol-id.yaml
A  schemas/cds/fixtures/_closure_records/<fixture-name>/{self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout}.md
A  tests/cdd/test_cn_cdd_validate_receipt.sh                        (~150 lines bash driver)
A  .cdd/unreleased/389/*.md                                         (cycle evidence)
```

Total estimated diff: ~1000 lines net new, ~80 lines modified.

## §Review-readiness signal

α writes `self-coherence.md` with §Review-readiness section once AC1–AC8 all pass on their oracles; β-collapsed self-review then runs the same oracles plus mechanical re-derivation.
