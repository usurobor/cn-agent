---
cycle: 367
role: alpha
issue: "https://github.com/usurobor/cnos/issues/367"
date: "2026-05-15"
manifest:
  sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
note: "δ-authored recovery after α SIGTERM. α committed all 8 doc sections (632L RECEIPT-VALIDATION.md) but was killed before writing self-coherence.md. Doc content is α's; this evidence wrapper is δ's."
---

# Self-coherence — #367

## Gap

The parent-facing validation interface (`V`) was named in `COHERENCE-CELL.md` (#364) but not specified. Five Open Questions were seeded and deferred. Without resolving them, Phases 2–7 of the #366 roadmap are speculative.

## Skills

- `alpha/SKILL.md` — role contract
- `design/SKILL.md` — interface-freeze discipline
- `write/SKILL.md` — design-doc authoring
- `COHERENCE-CELL.md` — predecessor doctrine (read-only)
- `CDD.md` — canonical algorithm (read-only)

## ACs

| AC | Description | Met? | Evidence |
|----|-------------|------|----------|
| AC1 | RECEIPT-VALIDATION.md at canonical path | yes | `test -f src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` → exists, 632L |
| AC2 | Draft-design status + predecessor citation | yes | Line 7: "Draft design"; Line 13: "COHERENCE-CELL.md"; Preamble §Status table: "Becomes binding when Phase 3 implements V" |
| AC3 | Q1 — when V fires | yes | §Q1: δ-boundary validation after γ receipt emission, before parent acceptance. γ-preflight explicitly non-authoritative. |
| AC4 | Q2 — capability or command | yes | §Q2: V is a capability; `cn-cdd-verify` is its command wrapper. Capability invoked by δ in-skill. |
| AC5 | Q3 — ε relocation | yes | §Q3: ε relocates to `ROLES.md`. Phase 6 ships the move. |
| AC6 | Q4 — override receipt shape | yes | §Q4: structured `override:` block inside receipt boundary section. Required fields: actor, justification, original_verdict, failed_predicates, degraded flag. Downstream detection rule: `boundary.override != null` → degraded. Override never rewrites ValidationVerdict. |
| AC7 | Q5 — closeouts as evidence-graph inputs | yes | §Q5: five files become evidence-graph inputs to V; receipt is parent-facing surface. Derivation rule: receipt computed at γ close-out from the five files. |
| AC8 | Validation interface frozen | yes | §Validation Interface: input (ContractRef × ReceiptRef × EvidenceRootRef), output (ValidationVerdict with verdict/failed_predicates/warnings/provenance), invocation (δ boundary). Verdict/decision distinction explicit. "Schema syntax is Phase 2" disclaimer present. |
| AC9 | Surface containment | yes | `git diff origin/main..HEAD --stat`: RECEIPT-VALIDATION.md (A) + .cdd/unreleased/367/*.md only. No prohibited files touched. |

## Self-check

- All five Open Questions answered with single chosen positions and rationale
- Validation interface frozen at doctrine level with input/output/invocation contracts
- ValidationVerdict and BoundaryDecision are syntactically distinct throughout
- Override never rewrites V; downstream detection rule named
- No .cue syntax; Phase 2 disclaimer present
- Surface containment verified via git diff

## Debt

1. **δ-authored self-coherence.md.** α was SIGTERM'd before writing this file. Doc content (632L) is α's across 8 commits. This evidence wrapper is δ's recovery per operator/SKILL.md §timeout-recovery.

## CDD Trace

| Step | Artifact | Decision |
|------|----------|----------|
| 0 Observe | issue #367, COHERENCE-CELL.md, CDD.md | Five Open Questions unresolved; validation interface unspecified |
| 1 Scaffold | gamma-scaffold.md | docs-only, design-and-build, 9 ACs, §5.1 multi-session |
| 2 Implement | RECEIPT-VALIDATION.md (632L) | Q1: δ-boundary; Q2: capability; Q3: ROLES.md; Q4: structured block; Q5: evidence-graph inputs |
| 3 Recovery | self-coherence.md (this file) | α SIGTERM'd after 8 commits; δ wrote evidence wrapper |

## Review-readiness

All 9 ACs mapped to evidence above. Review-ready.

- completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
