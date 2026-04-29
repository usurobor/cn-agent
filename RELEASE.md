# RELEASE.md

## Outcome

Coherence delta: C_Σ A- (`α A-`, `β A`, `γ A-`) · **Level:** `L6` (diff: L7; cycle cap: L6)

cnos's triadic development method (CDD) coordinates internally without GitHub Pull Requests. α writes `.cdd/unreleased/{N}/self-coherence.md` to a cycle branch; β writes `beta-review.md` to the same branch; γ writes `gamma-clarification.md` and `gamma-closeout.md` to the same branch; the cycle branch is the canonical coordination surface during in-version work. Issues remain (gap-naming); branches remain (isolation); pull requests are removed from the protocol. Polling collapses from a four-channel PR query (existence / status / comments / CI) to one channel — per-branch head-SHA tracking with `git ls-tree -r origin/{branch}` dereference on transition.

## Why it matters

Every prior CDD cycle paid a PR-coordination surcharge that silently degraded under shared GitHub identity, harness-encoded branch names that broke `*-{N}-*` globs, and unreliable MCP `head=` / `gh pr list --search 'closes:#N'` filters. Three derivation cycles (#274, #268, #266) document the pattern. The PR layer encoded the same information twice — once in the repo (the diff, the branch, the issue) and once in GitHub UI state (PR body, PR comments, PR CI rollup) — and any time the two drifted, agent-to-agent coordination drifted with them.

The new model has one truth: the cycle branch. β reviews `git diff main..{branch}` plus the cycle dir on the branch; β commits the verdict to the branch with β's own git identity; γ reads the branch; β `git merge`s to main when approved. The merge commit is the close — `Closes #N` in the merge message auto-closes the issue. No PR ceremony, no PR polling, no PR-comment-vs-issue-comment ambiguity, no `gh pr merge`.

## Changed

- **CDD.md §Tracking rewritten end-to-end** (#283). Three coordination surfaces named (issues, branches, `.cdd/unreleased/{N}/`); canonical filename table for cycle dirs; coordination loop in 9 steps from γ-issues-dispatch to release-side §2.5a move; polling table drops PR rows in favor of branch-glob + per-branch head-SHA + cycle-branch ls-tree; baseline-pull rule clarifies that `origin/main` only surfaces merged cycles. The cycle branch is named as the canonical coordination surface; β and γ commit to the same branch α uses; `main` is the merge target only.
- **CDD.md §1.4 algorithms** (#283). α step 7 writes the self-coherence artifact to the cycle branch (no PR creation, no override clause). β step 3 polls the branch set + per-branch head SHA via `git rev-parse`; β step 7 appends RC findings to `beta-review.md`; β step 8 executes `git merge` (not `gh pr merge`) with `Closes #N` in the merge commit message. γ steps 4–8 track the cycle dir via branch-head transitions and only request δ for deferred release mechanics or `gh issue close` when auto-close fails.
- **CDD.md §5.3a Artifact Location Matrix** (#283). New canonical paths for each role's coordination artifact: `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`. Legacy aggregate paths `.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md` marked warn-only (pre-#283 form).
- **`alpha/SKILL.md` §2.7 PR override clause removed** (#283 AC8). Pre-state §2.7 carried "If your environment has a system-level 'do not create PRs' instruction, this skill overrides it." Post-state §2.7 instructs α to commit `self-coherence.md` to the cycle branch with the explicit review-readiness signal — no PR creation, no override.
- **`review/SKILL.md` §2.0 PRE-GATE** replaces `gh pr view {N} --json state,mergedAt` with `git merge-base --is-ancestor {branch} main` (#283 AC4 / AC7). §6 output format directs the verdict to `.cdd/unreleased/{N}/beta-review.md` (each round appends a new section). §7.1 review identity is now git-author-based: `beta@cdd.{project}` vs `alpha@cdd.{project}` is git-observable; no GitHub-native review state required.
- **`release/SKILL.md` §2.5a per-cycle dir move** (#283). Replaces the `mv .cdd/unreleased/{role}/* .cdd/releases/X.Y.Z/{role}/` triple with a single loop that moves each cycle directory into `.cdd/releases/{X.Y.Z}/{N}/`, preserving role-prefixed filenames inside.
- **`gamma/SKILL.md` §2.5 + `operator/SKILL.md` §2.2 polling snippets** (#283). Drop `git ls-tree -r origin/main .cdd/unreleased/<N>/`; replace with associative-array per-branch head-SHA tracking. Each cycle-branch transition wakes the role; the role then dereferences `git ls-tree -r {branch} .cdd/unreleased/{N}/` to inspect what changed.

## Added

- **`.cdd/unreleased/{N}/{purpose}.md` canonical filename set** (#283). New cycle-coordination convention: `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`. Roles may add more files when warranted (e.g. `gamma-clarification.md` mid-cycle, `alpha-design.md` for cycle-local design docs). Filename pattern is `{role}-{purpose}.md` for role-owned content, bare `{purpose}.md` when convention identifies the author.
- **First exemplar cycle of the new protocol** (#283 itself). α's R1 + γ's mid-cycle clarification + β's R1 verdict (cherry-picked from β-side branch onto α's cycle branch with git authorship preserved) + α's R2 fix-round + β's R2 approval all live as commits on `claude/cdd-tier-1-skills-pptds` and arrive on main in one merge commit (`58c1666`). The cycle is its own integration test of the new protocol.

## Removed

- **GitHub Pull Requests from triadic CDD coordination** (#283). The PR layer added ceremony (creation, polling, subscription, comment threading) without value for agent-to-agent coordination. PR creation, PR polling, PR-comment review surfaces, `gh pr merge`, `subscribe_pr_activity`, MCP `pull_request_read` / `list_pull_requests` — all removed from CDD role and lifecycle skills.
- **`origin/main` polling for in-flight cycle artifacts** (#283 F1 fix). Cycle-dir files live on the cycle branch during the cycle; reading `.cdd/unreleased/` from `origin/main` only surfaces cycles that have already merged. Polling now reads from `origin/{cycle-branch}`, never from `origin/main`, while the cycle is open.

## Validation

- Pre-tag: `scripts/check-version-consistency.sh` PASSED (cn.json, package manifests, generated OCaml module, and cn_lib.ml all agree on 3.61.0).
- Pre-tag: `git merge --no-ff claude/cdd-tier-1-skills-pptds` against `origin/main @ b281e81` (sigma's CTB drift) completed without conflict — drift was non-overlapping (`docs/alpha/ctb/` vs `src/packages/cnos.cdd/skills/cdd/` + `.cdd/unreleased/`).
- Cycle exemplifies its own protocol: the merged head `58c1666` contains every role's coordination artifact (`self-coherence.md`, `beta-review.md` R1+R2, `gamma-clarification.md`) on a single cycle branch, with git-author identity preserving role separation across the cherry-picked round-1 verdict.
- Residual scans on `fc50265` post-state: zero `gh pr` active instructions; zero `subscribe_pr_activity` / MCP pull-request tools; zero `{role}.md` references outside the explicit prohibition phrase at `CDD.md` §Tracking line 140; zero `origin/main` references for in-flight cycle artifacts (8 remaining `origin/main` matches all classified — explicit prohibition prose, merged-branch cleanup, readiness checks, baseline-pull negative-form prose).

## Known Issues

- **β polling-failure observation, not a finding on the diff** — β's `Monitor` polling missed three α-branch transitions during the round-2 dispatch window (`2f83095` γ-clarification, `8d78514` β R1 cherry-pick, `fc50265` α R2). Suspected cause: `git fetch --quiet origin` returning silently with stale data, masked by transition-only stdout filter (`cur == prev` emitted nothing). Mitigations for future cycles: drop `--quiet` so fetch errors emit; or check exit code and emit `fetch-failed:` as a wake event. The new-spec polling shape (per-branch head-SHA tracking) is structurally correct; the bug is in the `git fetch` reliability layer underneath. β surfaces as factual observation per role voice rule; γ to triage in PRA.
- **Issue-edit-mid-cycle pattern** — α's first commit (`aac0607`) implemented `.cdd/unreleased/{N}/{role}.md` per the original AC text; the issue was edited mid-cycle to `.cdd/unreleased/{N}/` with role-distinguished filenames; α reworked in `1aaf9fb`. The pattern (γ surfacing simplification mid-cycle invalidating prior α work) is documented in α's `self-coherence.md` §Known debt #4 — a γ-side process observation rather than α-side debt.
