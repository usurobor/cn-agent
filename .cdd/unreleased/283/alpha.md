# α Artifact — Cycle #283

This file is the primary branch artifact for cycle #283 per CDD.md §5.4. It carries the gap, mode, active skills, impact graph, CDD Trace, AC mapping, known debt, review-readiness signal, fix-round appendices, and final α close-out.

## Gap

CDD's triadic coordination model uses GitHub Pull Requests for agent-to-agent communication: α creates a PR (`gh pr create`), all three roles poll PR state (`gh pr view`, comment threads, review state, CI checks), β posts review findings as PR comments, β merges via `gh pr merge`, and close-outs reference PR numbers throughout. PRs are a GitHub UI surface, not a repo artifact — every PR action consumes context budget on polling/subscription mechanics that produce no durable repo evidence. The `.cdd/unreleased/` directory already exists for cycle-local artifacts; the natural artifact-exchange surface is `.cdd/unreleased/{N}/{role}.md`, one cycle directory per issue keyed by issue number, one role file per role inside. Multiple `Derives from` annotations (#274, #268, #266) document repeated failures of the PR-polling model: first-iteration absorption of pre-existing PRs, branch-glob templates that assume issue-number encoding the harness doesn't produce, transient PR-state rows that drift between write and read, and shared-GitHub-identity review-state degradation.

The three failure modes the PR coordination layer creates:

1. **Polling ceremony.** Every cycle pays a PR-existence/PR-status/PR-comments/PR-CI quadruple-poll that silently degrades when any of: `gh auth` is unavailable, MCP and shell tools have unequal reachability, branch-globs miss harness-encoded names, or shared identity defeats GitHub-native review state. These are observable failures because they happened, not theoretical risk.
2. **Surface duplication.** A PR body restates the gap; a PR comment restates a finding; a release CI auto-generates from PR titles; a close-out references a PR number. Each is a copy of information already (or potentially) in a repo artifact. Drift between copies is silent.
3. **Coordination over a non-repo surface.** A PR is not a git artifact. It cannot be `git diff`ed, cannot be merged via `git merge`, cannot be moved into the release directory by a `mv`. Every coordination action on a PR is a side effect outside the repo's audit trail.

The change replaces PR-based coordination with `.cdd/unreleased/{N}/{role}.md` — one cycle directory per issue, one role file per role, all coordination expressed as commits to those files. Issues remain (gap-naming); branches remain (isolation); PRs are removed from the triadic protocol.

## Mode

| | |
|---|---|
| Mode | MCA |
| Cycle level (target) | L7 — system-shaping leverage. Replaces a coordination protocol with a leaner one that eliminates a recurring friction class (PR polling). |
| Tier 1 skills | `cdd/SKILL.md`, `cdd/CDD.md`, `cdd/alpha/SKILL.md` (this is CDD self-modification — no Tier 2 or Tier 3 beyond Tier 1 per the issue) |
| Tier 2 skills | not applicable per the issue's "Tier 3 skills: none beyond CDD Tier 1" declaration; Tier 2 `eng/document` would normally apply for prose discipline but is not loaded since CDD's lifecycle sub-skills (design, plan, issue) are themselves Tier 1 here |
| Tier 3 skills | none beyond CDD Tier 1 (per issue declaration) |
| Work shape | Substantial — 10 files declared in issue + the migration of existing in-version cycle artifacts |
| Affected files (modified) | `src/packages/cnos.cdd/skills/cdd/{CDD.md, SKILL.md, alpha/SKILL.md, beta/SKILL.md, gamma/SKILL.md, review/SKILL.md, release/SKILL.md, post-release/SKILL.md, operator/SKILL.md}` |
| Affected files (file moves) | `.cdd/unreleased/{alpha,beta,gamma}/{268,278}.md` → `.cdd/unreleased/{268,278}/{alpha,beta,gamma}.md` (6 renames; existing in-version cycles migrated to the new layout) |
| Affected files (declared but unchanged) | `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — issue lists it but a full PR-ref scan returned zero hits. Leaving unchanged is correct. |
| Branch | `claude/cdd-tier-1-skills-pptds` (carries the dispatch's pre-set branch name; the cycle scope claims this branch) |
| Design artifact | not required — the issue body carries Problem, Impact, ACs, Non-goals, related artifacts, and dependencies; this α artifact carries the impact graph and CDD Trace. A separate design doc would duplicate without adding leverage. |
| Plan | not required — implementation order is `CDD.md` (canonical authority) → role skills (α/β/γ) → lifecycle sub-skills (review/release/post-release/operator) → top-level SKILL.md → in-version artifact migration → α artifact write. The order is dictated by authority hierarchy + the natural flow of the new coordination model; no novel sequencing is needed. |
