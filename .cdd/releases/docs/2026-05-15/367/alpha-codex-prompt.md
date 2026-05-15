You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 367 --json title,body,state,comments
Branch: cycle/367
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.core/skills/design/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/SKILL.md

Predecessor doctrine (read as source of truth, do not modify):
  src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md — names V, receipt rule, four-way separation, the five Open Questions verbatim in §Open Questions
  src/packages/cnos.cdd/skills/cdd/CDD.md — canonical executable algorithm; unchanged by this cycle

Target artifact:
  src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md (NEW; ≈300–600 lines; sibling to COHERENCE-CELL.md)

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

Scaffold reference: .cdd/unreleased/367/gamma-scaffold.md carries γ's gap, peer enumeration, mode, ACs reference, acceptance posture, dispatch configuration, and the eight failure modes γ pre-flagged for this cycle. Read it first — it is the authoritative starting contract alongside the issue body.

Issue invariants (compressed; issue body is binding):
- 9 ACs, all docs-only. AC9 forbids any file outside {new doc + cycle evidence}.
- AC2 requires literal tokens: "draft design", "COHERENCE-CELL.md", and an explicit binding-under-Phase-3 clause.
- AC3–AC7 each resolves one of the five Open Questions seeded by #364 — single chosen answer with rationale, never "TBD" or option-enumeration-without-choice.
- AC8 freezes the validation interface at doctrine level (input / output / invocation), distinguishes ValidationVerdict from BoundaryDecision, and explicitly defers schema syntax to Phase 2. Do not write .cue syntax.
- AC9 surface-containment: only RECEIPT-VALIDATION.md (A) + .cdd/unreleased/367/*.md. No edits to COHERENCE-CELL.md, CDD.md, ROLES.md, operator/SKILL.md, gamma/SKILL.md, epsilon/SKILL.md, or commands/cdd-verify/.

Authoring discipline:
- Apply src/packages/cnos.cdd/skills/cdd/CDD.md §1.4 large-file authoring rule: write section by section to disk with an HTML-comment section manifest at the top of RECEIPT-VALIDATION.md; update completed: as each section lands.
- Design over checklist: the surface is the doctrine prose that justifies each chosen position, not five one-line answers stapled together (#364 PRA observation).
- γ-preflight is non-authoritative; δ-boundary validation is the authoritative firing point (γ scaffold §Failure mode 5).
- Override never rewrites or masks ValidationVerdict; override is detectable by downstream consumers via the rule named in §Q4 (γ scaffold §Failure mode 7).

Commit checkpoints: commit after each major section lands on disk (file create with manifest + intro; then per-Q sections; then §Validation Interface; then self-coherence.md updates). First commit must land within the first 25% of the dispatch budget. Do not batch all work into a single end-of-cycle commit.

Review-readiness: when AC1–AC9 are mapped to evidence in .cdd/unreleased/367/self-coherence.md and all per-AC oracles have been run and recorded inline, signal review-readiness per alpha/SKILL.md.
