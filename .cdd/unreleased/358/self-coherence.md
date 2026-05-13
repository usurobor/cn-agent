<!-- sections: [gap, skills, ac-evidence, self-check, debt, cdd-trace, review-readiness] -->
<!-- completed: [gap, skills, ac-evidence, self-check, debt, cdd-trace, review-readiness] -->

# Self-Coherence — Issue #358

## Gap

Issue #358 closes a docs/skill gap: cross-repo proposal issue packs can be submitted to cnos, but CDD did not define a durable proposal `STATUS` lifecycle or require target-side feedback when a proposal is accepted, modified, rejected, or landed.

Mode: docs-only. Branch: `cycle/358-gpt`.

## Skills

- Tier 1 / role: `CDD.md`, `gamma/SKILL.md`, `issue/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`.
- Lifecycle: `post-release/SKILL.md`.
- Tier 3: skill authoring via the existing CDD skill files.

## AC Evidence

AC1: STATUS format and lifecycle added.
Evidence: `src/packages/cnos.cdd/skills/cdd/CDD.md` now defines proposal layouts, append-only `STATUS` format, event vocabulary, no-rewrite rule, source lineage, and `landed` feedback.

AC2: gamma observation/selection scans cross-repo proposals.
Evidence: `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` Step 1a now scans both proposal path families, reads submitted proposals, checks target state, decides accepted/modified/rejected, creates source blocks, and updates source `STATUS` or emits patches.

AC3: gamma close-out covers landed proposals.
Evidence: `gamma/SKILL.md` Step 8-9 close-out now requires landed status or feedback patch for accepted/modified source proposals touched by the cycle.

AC4: post-release checklist item added.
Evidence: `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` pre-publish gate contains `Cross-repo proposal status updated or feedback patch emitted.`

AC5: optional Source Proposal block added.
Evidence: `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` minimal output pattern includes an optional `## Source Proposal` block, and the handoff checklist verifies it when a target issue comes from a source proposal.

AC6: minimal docs/skill scope preserved.
Evidence: no activation scaffolding, generated index, CI validator, command spec, migration, or runtime code was added.

## Self-Check

Peer enumeration found existing `STATUS` mentions for cross-repo trace bundles (`open|converging|closed`) in activation artifacts and `.cdd/iterations/cross-repo/README.md`. This cycle does not alter that trace-bundle format; it adds a separate proposal-specific event log contract.

All files in the implementation diff are named in the AC evidence above. No test runner is applicable because the change is Markdown skill/spec text only.

## Debt

Known gap: no migration of existing source proposal directories and no validator for malformed proposal `STATUS` files. Both are explicit non-goals in #358.

## CDD Trace

| Step | Evidence |
|---|---|
| 1-2 | Issue #358 loaded; issue validation recorded in `gamma-scaffold.md`. |
| 3 | Existing `cycle/358-gpt` branch used. |
| 4-6 | Docs/skill implementation landed in `CDD.md`, `gamma/SKILL.md`, `issue/SKILL.md`, and `post-release/SKILL.md`. |
| 7 | This self-coherence maps each AC to diff evidence. |

## Review Readiness

Ready for beta review on branch HEAD after the alpha implementation commit. Branch CI is not run locally for this Markdown-only change.
