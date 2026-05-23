# α Closeout — Cycle #389

Actor: alpha@cdd.cnos (δ-as-agent, collapsed-mode)

## §Summary

Implemented V (the executable validator predicate) per Phase 3 of #366. V's signature matches CCNF: `Contract × Receipt → ValidationVerdict`. V dispatches on `protocol_id` declared in receipt frontmatter, wraps `cue vet` for structural validation, dereferences evidence refs bound into the receipt, applies counterfeit-receipt rules, and emits a typed `ValidationVerdict` (JSON Schema at `schemas/cdd/validation_verdict.schema.json`).

37/37 AC oracles pass mechanically.

## §Deliverables

- `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt` (new) — V's Python core
- `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` (modified) — bash wrapper now dispatches into V via `--receipt`
- `schemas/cdd/validation_verdict.schema.json` (new) — JSON Schema for V's output
- `schemas/cds/fixtures/counterfeit-{actor-collision,merge-precedes-verdict,override-rewrite,evidence-missing,mismatched-protocol-id}.yaml` (new) — counterfeit corpus
- `schemas/cds/fixtures/_closure_records/<fixture-name>/{self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout}.md` (new) — closure-record stubs
- `tests/cdd/test_cn_cdd_validate_receipt.sh` (new) — AC oracle harness
- `schemas/cds/fixtures/valid-receipt.yaml` (modified) — paths point at real on-disk closure records; `mode: collapsed` acknowledged
- `.cdd/unreleased/389/*.md` — cycle evidence

## §Findings

**One design-finding addressed in-cycle (F0 in beta-review.md):** receipt-level `mode: collapsed` field added. V emits structural-review-independence rule violations as warnings (not failed_predicates) when the receipt declares it. This makes V honest about γ+α+β-collapsed-on-δ mode (the wave-manifest precedent) — the trust degradation is preserved as a warning surface rather than silently hidden or noisily failed.

## §Debt named (deferred per design-notes §Backward compatibility / §Risks)

- D1: `git_diff(...)` ref resolution (scope-drift detection) — Phase 4/5 work
- D2: `validator_version` autodiscovery
- D3: `mode` field formalization in `schemas/cdd/receipt.cue` (the field works today via open-extension; pinning is a follow-up)
- D4: historical CDS-fixture migration into V full-mode validation (not in scope; per #389 non-goals)
- D5: γ-skill `protocol_id` automation patch (#388 F1; deferred)

## §Closure
α work complete. β-collapsed self-review APPROVED.
