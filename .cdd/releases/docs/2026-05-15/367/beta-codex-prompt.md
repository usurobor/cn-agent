You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 367 --json title,body,state,comments
Branch: cycle/367

Predecessor doctrine (read as source of truth for coverage judgement):
  src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md — V definition, receipt rule, four-way separation, the five Open Questions verbatim
  src/packages/cnos.cdd/skills/cdd/CDD.md — canonical algorithm; verify this cycle does not touch it

Target artifact under review:
  src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md (NEW; α-authored on cycle/367)

Git identity:
  git config user.name "beta"
  git config user.email "beta@cdd.cnos"

Scaffold reference: .cdd/unreleased/367/gamma-scaffold.md carries γ's gap, peer enumeration, mode, acceptance posture, the eight γ-flagged failure modes, and dispatch configuration. Read it alongside α's self-coherence.md when judging AC coverage. Issue body remains the binding contract.

Review surface (compressed; β skill is authoritative):
- All 9 ACs must map to evidence in α's self-coherence.md and to inspectable content in RECEIPT-VALIDATION.md.
- AC8 verdict/decision distinction is load-bearing: ValidationVerdict (emitted by V) must be syntactically and semantically distinct from BoundaryDecision (emitted by δ); override never rewrites the verdict. If the doc collapses them, return RC.
- AC3 (Q1) must place authoritative firing at δ-boundary, not γ-preflight. If γ-preflight is allowed, it must be explicitly non-authoritative.
- AC6 (Q4) override receipt must include downstream-consumer detection rule. An override that masquerades as PASS is a binding finding.
- AC9 surface-containment: run `git diff origin/main..HEAD --stat` and confirm exactly {RECEIPT-VALIDATION.md (A), .cdd/unreleased/367/*.md}. Any other file is a binding finding.
- No .cue syntax in §Validation Interface; the "schema syntax is Phase 2" disclaimer must appear.
- AC2 literal tokens "draft design" and "COHERENCE-CELL.md" plus a binding-under-Phase-3 clause must be present.

Target round count: ≤1 review round (docs-only convention; cycle #364 precedent).

Pre-merge gate (beta/SKILL.md): run the full row-by-row check before merging. The cycle is docs-only; gate row 2 (canonical-skill staleness fetch) and row 3 (non-destructive merge-test against current origin/main) still apply.

Integration: on APPROVE, merge cycle/367 into main per beta/SKILL.md §1, write beta-review.md + beta-closeout.md to .cdd/unreleased/367/, exit. Do not tag — δ owns the release boundary.
