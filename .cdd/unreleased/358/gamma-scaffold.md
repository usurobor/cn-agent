<!-- sections: [intake, issue-validation, extracted-acceptance-criteria, peer-enumeration, implementation-plan, non-goals-guard] -->
<!-- completed: [intake, issue-validation, extracted-acceptance-criteria, peer-enumeration, implementation-plan, non-goals-guard] -->

# Gamma Scaffold — Issue #358

## Intake

- Issue: #358, `cdd(cross-repo): proposal lifecycle — STATUS file + feedback loop for satellite repos`
- Labels: `enhancement`, `P2`, `cdd`
- Mode: docs-only, inferred from the issue body and required changes.
- Work shape: small docs/skill cycle on the existing cycle branch.
- Branch: `cycle/358-gpt`

## Issue Validation

`issue/SKILL.md` was loaded after `CDD.md` and `gamma/SKILL.md`.

Validation result: dispatchable with normalization.

- Problem / impact: present in the R3 proposal. Cross-repo proposal issue packs can remain stale at `submitted`, and target decisions do not reliably flow back to the source repo.
- Status truth: present. The issue states the current state as proposal bodies without a durable lifecycle status protocol.
- Source of truth: present by issue body. The converged R3 proposal is the source for this cycle.
- Scope: present. Required skill changes are limited to `cdd/gamma`, `cdd/post-release`, and `cdd/issue`; this scaffold also includes `CDD.md` because the STATUS format is lifecycle-contract material.
- Acceptance criteria: implicit, extracted from "Required Skill Changes" and "Minimum Viable Protocol" rather than numbered issue AC headings.
- Proof surface: branch diff against `main`, plus peer grep for proposal/STATUS/Source Proposal surfaces.
- Non-goals: present. No activation scaffolding, generated indexes, CI validators, new commands, forced migration, target mirrors, YAML frontmatter, or review-board states.

Normalization note: the issue does not declare a `Mode:` header, cycle sizing table, or numbered AC headings. The body is nevertheless concrete enough for a docs-only implementation because it names exact files, exact text blocks, and explicit non-goals.

## Extracted Acceptance Criteria

1. Add a minimal proposal `STATUS` file format and lifecycle contract to the canonical CDD skill surfaces.
2. Update `cdd/gamma/SKILL.md` so observation/selection scans known cross-repo proposal paths, reads submitted proposals, checks target state, decides accepted/modified/rejected, creates or links target issues with source lineage, and updates source `STATUS` or emits a patch.
3. Update `cdd/gamma/SKILL.md` so close-out requires every accepted or modified source proposal touched by the cycle to receive a `landed` event or feedback patch.
4. Update `cdd/post-release/SKILL.md` with the checklist item: `Cross-repo proposal status updated or feedback patch emitted.`
5. Update `cdd/issue/SKILL.md` with an optional `## Source Proposal` target issue block.
6. Keep the change minimal and docs/skill-only; do not add activation scaffolding, generated indexes, CI gates, command specs, or broad migration.

## Peer Enumeration

Directory enumeration:

- `find src/packages/cnos.cdd/skills/cdd -maxdepth 2 -type f` lists the CDD role and lifecycle skill files, including `CDD.md`, `gamma/SKILL.md`, `issue/SKILL.md`, `post-release/SKILL.md`, `release/SKILL.md`, `review/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `operator/SKILL.md`, `activation/SKILL.md`, `design/SKILL.md`, `plan/SKILL.md`, `epsilon/SKILL.md`, and the package-visible `SKILL.md`.

Grep evidence:

- `rg "\bSTATUS\b|Source Proposal|cross-repo proposal|submitted proposal|proposal status|landed event|feedback patch|accepted or modified" src/packages/cnos.cdd/skills/cdd .cdd ...`
- Existing matches are limited to activation / cross-repo trace bundle `STATUS` files with `open|converging|closed`, and historical release artifacts.
- No existing CDD skill defines the proposal lifecycle event-log format, the optional issue `Source Proposal` block, or gamma proposal intake / landed close-out obligation.

Framing: this cycle adds a new proposal-specific `STATUS` protocol and must not overwrite the older cross-repo trace bundle `STATUS` meaning.

## Implementation Plan

1. Patch `CDD.md` near artifact-driven coordination with the proposal lifecycle contract and `STATUS` event-log format.
2. Patch `gamma/SKILL.md` in observation/selection and close-out triage with intake and landed-status obligations.
3. Patch `issue/SKILL.md` with the optional `## Source Proposal` block and checklist coverage.
4. Patch `post-release/SKILL.md` pre-publish gate with proposal status / feedback patch verification.
5. Write `self-coherence.md`, verify `git diff --stat main..HEAD` includes implementation files, then push.

## Non-Goals Guard

- No path migration from `.cdd/iterations/proposals/`.
- No mandatory target mirror directory.
- No generated index or CI validator.
- No new `cn proposal` command.
- No broad migration of existing proposal files.
