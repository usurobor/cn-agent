# α close-out — cycle/397

**Issue:** cnos#397 — Phase 4a of #366
**Mode:** design-and-build; γ+α+β-collapsed-on-δ
**Identity:** delta@cdd.cnos

## What was built

1. **NEW** `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (340 lines) — δ-role skill carrying:
   - Two-sided membrane doctrine (outward §1 + inward §2)
   - Override semantics (§3) with cnos#367 verdict-vs-decision freeze
   - Composition with V (§4) per `RECEIPT-VALIDATION.md` §Validation Interface
   - δ-as-role boundary constraints (§5)
   - Cross-references and Phase tracking (§6, §7)

2. **EDIT** `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (109 lines net delta):
   - §3.1-§3.3 → role-policy moved to delta/; mechanics-execution discipline retained with delta/ redirects
   - §3.4 → role-policy paragraph moved to delta/; `scripts/release.sh` runbook retained (Phase 4c relocation target)
   - §3.5 → role-policy moved to delta/; in-place redirect added
   - §3a → entire section replaced by one-paragraph redirect pointing at delta/SKILL.md §2
   - §4 → entire section replaced by one-paragraph redirect pointing at delta/SKILL.md §3
   - Algorithm step 7 enriched with delta/SKILL.md §3 ref
   - §timeout-recovery §7.3 override-declaration ref updated from §4 → delta/SKILL.md §3

3. **EDIT** cross-references (4 path-only swaps; no content changes):
   - `gamma/SKILL.md` line 363, line 370 → `delta/SKILL.md` §2
   - `alpha/SKILL.md` line 355 → `delta/SKILL.md` §2
   - `beta/SKILL.md` line 175 → `delta/SKILL.md` §2

## Lifecycle artifacts

- `.cdd/unreleased/397/gamma-scaffold.md` — surfaces, AC oracle approach, refusal-condition check
- `.cdd/unreleased/397/design-notes.md` — full inventory of operator/SKILL.md sections + classification (δR/OC/HM/RE) + cross-reference plan + content-loss audit plan + implementation-contract conformance plan
- `.cdd/unreleased/397/self-coherence.md` — AC1-AC7 oracle results, implementation-contract conformance, content-loss audit, review-readiness
- `.cdd/unreleased/397/beta-review.md` — R1 APPROVE; AC re-runs; β-rigor checks (implementation-contract, no-content-loss, no-broken-refs, Phase 4b/4c untouched)
- `.cdd/unreleased/397/alpha-closeout.md` (this file)

## Implementation-contract conformance

All 7 axes pinned by issue body satisfied per self-coherence.md §2.

## Observations for ε

- The β-α-collapsed-on-δ pattern carried cleanly for this extract-by-restructure cycle: the content-loss audit (β-rigor) was self-administered without role-isolation pressure because the operator/SKILL.md → delta/SKILL.md mapping is purely mechanical (every removed paragraph maps to a corresponding new paragraph; verification is row-by-row).
- One observation worth surfacing: when the cycle has a precursor whose body explicitly anchored the content being moved (here cnos#393 §3a authored the doctrine that this cycle relocates), the move becomes near-mechanical. Future Phase-4b and Phase-4c cycles will face the same shape: the substance is anchored; the relocation is restructure.

No protocol gaps (`protocol_gap_count = 0`).
