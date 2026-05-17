You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 370 --json title,body,state,comments
Branch: cycle/370
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.core/skills/design/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/SKILL.md

Predecessor doctrine + design peer (read as source of truth, do not modify):
  src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md — receipt rule, four-way separation, role-as-cell-function; the organism this kernel describes algorithmically
  src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md — V's interface, verdict/decision distinction, δ-authoritative validation; the receptor positioned at steps 4–5 of the kernel
  src/packages/cnos.cdd/skills/cdd/CDD.md — canonical executable algorithm; unchanged by this cycle (Phase 7's rewrite target)

Target artifact:
  src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md (NEW; target 200–400 lines; sibling to COHERENCE-CELL.md and RECEIPT-VALIDATION.md)

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

Scaffold reference: .cdd/unreleased/370/gamma-scaffold.md carries γ's gap, peer enumeration, mode, ACs reference, acceptance posture, dispatch configuration, and the nine failure modes γ pre-flagged for this cycle. Read it first — it is the authoritative starting contract alongside the issue body.

Issue invariants (compressed; issue body is binding):
- 9 ACs, all docs-only. AC9 forbids any file outside {new doc + cycle evidence}.
- AC2 requires literal tokens: "draft doctrine", "COHERENCE-CELL.md", "RECEIPT-VALIDATION.md", and explicit declaration of this doc as kernel-layer companion to the predecessor doctrine.
- AC3 kernel as five-step recursion at scope n, with evidence-binding rule. β's signature explicitly excludes evidence; δ's signature explicitly excludes evidence; γ binds evidence into the receipt; V dereferences from the receipt.
- AC4 four cell outcomes (accepted, degraded, blocked, invalid) with verdict × decision preconditions matching #369 AC4; invalid is non-terminal (δ must re-decide).
- AC5 two recursion modes (within-scope repair-dispatch; cross-scope accept/degraded); same recursion operator, different scope behavior.
- AC6 three scope-lift projections (closed cell → α-matter; δ → β-discrimination; ε → γ-coordination/evolution) plus projection-not-renaming framing plus β/γ-no-upward-projection clause.
- AC7 substrate-independence in kernel sections — section-bounded oracle via awk-extract + word-bounded rg. Section headers must be EXACTLY `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, `## Scope-Lift` (capitalized, hyphenated, no qualifiers).
- AC8 two-layer separation with four realization-peer citations: RECEIPT-VALIDATION.md, schemas/cdd/ (#369), cn-cdd-verify (Phase 3, deferred), CDD.md (Phase 7 rewrite, deferred).
- AC9 surface-containment: only COHERENCE-CELL-NORMAL-FORM.md (A) + .cdd/unreleased/370/*.md. No edits to COHERENCE-CELL.md, RECEIPT-VALIDATION.md, CDD.md, schemas/cdd/, cn-cdd-verify, operator/SKILL.md, gamma/SKILL.md, epsilon/SKILL.md, ROLES.md, CI workflows.

Authoring discipline:
- Apply src/packages/cnos.cdd/skills/cdd/CDD.md §1.4 large-file authoring rule: write section by section to disk with an HTML-comment section manifest at the top of COHERENCE-CELL-NORMAL-FORM.md; update completed: as each section lands.
- Kernel discipline: substrate-independent. No GitHub, CUE, cn-cdd-verify, cn dispatch, claude, gh in kernel sections. Realization references live only in §Realization Layer or §Two-Layer Separation, clearly labeled.
- Section header discipline: use EXACTLY `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, `## Scope-Lift`. The AC7 awk extracts via these literal headers; "## Kernel Algorithm" or "## The Kernel" or "## Scope-lift" silently fails the extraction.
- Doctrine over operational: kernel names the *what*; do not preempt Phase 7's operational rewrite. Target 200–400 lines.
- No angle-bracket placeholders (`<placeholder>`); GitHub renderer strips them. Use `{placeholder}` or literal forms (matterₙ, receiptₙ, etc.).
- Evidence-binding rule must be stated and visible: β consumes matter only; γ binds evidence into the receipt; V dereferences refs from receipt; δ decides on receipt + verdict only.
- Recursion modes must be distinguished: within-scope repair-dispatch (same scope index, parent γₙ re-emits) vs cross-scope accept/degraded (scope index advances). Do not conflate.
- Scope-lift: three projections, projection-not-renaming framing, β/γ-no-upward-projection clause.
- #369 alignment notes: AC4 (four outcomes) cites #369 AC4; AC6 (scope-lift) cites #369 AC3. Reference the alignment explicitly in doc body.

Commit checkpoints: commit after each major section lands on disk (file create with manifest + intro; then §Kernel; then §Cell Outcomes; then §Recursion Modes; then §Scope-Lift; then §Two-Layer Separation; then self-coherence.md updates). First commit must land within the first 25% of the dispatch budget. Do not batch all work into a single end-of-cycle commit.

Per-AC oracle execution (record results inline in self-coherence.md):
- AC1: `test -f src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`
- AC2: `rg 'draft doctrine' src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (≥1) and same for `COHERENCE-CELL\.md` and `RECEIPT-VALIDATION\.md`
- AC3: prose check — five steps in order, evidence-binding rule named, β-no-evidence and δ-no-evidence stated
- AC4: prose check — four outcomes named with verdict × decision; invalid non-terminal stated
- AC5: prose check — both modes named; within-scope vs cross-scope distinction explicit
- AC6: prose check — three projections named; projection-not-renaming framing present; β/γ-no-upward-projection clause present
- AC7: run the literal awk + rg snippet from the issue body against the authored file:
    awk '/^## (Kernel|Cell Outcomes|Recursion Modes|Scope-Lift)/ {keep=1} /^## / && !/^## (Kernel|Cell Outcomes|Recursion Modes|Scope-Lift)/ {keep=0} keep {print}' src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md > /tmp/ccnf-kernel.md
    ! rg -i '\b(github|cue|cn-cdd-verify|cn dispatch|claude|gh)\b' /tmp/ccnf-kernel.md
- AC8: prose check + rg for `RECEIPT-VALIDATION\.md`, `schemas/cdd/`, `cn-cdd-verify`, `CDD\.md` (all four in §Two-Layer Separation or §Realization Layer)
- AC9: `git diff origin/main..HEAD --stat` — exactly {COHERENCE-CELL-NORMAL-FORM.md (A), .cdd/unreleased/370/*.md}

Review-readiness: when AC1–AC9 are mapped to evidence in .cdd/unreleased/370/self-coherence.md and all per-AC oracles have been run and recorded inline, signal review-readiness per alpha/SKILL.md.
