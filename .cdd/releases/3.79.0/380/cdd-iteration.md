---
cycle: 380
date: "2026-05-19"
issue: "https://github.com/usurobor/cnos/issues/380"
merge_sha: "770ea1b4"
findings_count: 1
patches_count: 1
mcas_count: 0
no_patch_count: 0
---

# CDD Iteration — #380

One `cdd-skill-gap` finding surfaced during the cycle and from γ's PRA triage: F1 (honest-claim — pre-rebase SHA citations in `self-coherence.md` §CDD Trace + §ACs) was caused by `alpha/SKILL.md` §2.6 row 14 path (a) being mechanically complete but downstream-incomplete. The patch lands in the same PRA commit; no findings deferred to a next-MCA.

Format per `cdd/post-release/SKILL.md` Step 5.6b.

## F1: `alpha/SKILL.md` §2.6 row 14 path (a) under-specifies SHA-citation invalidation as a downstream consequence

- **Source:** β-review R1 F1 (`.cdd/unreleased/380/beta-review.md` §"Findings"); α-closeout §"Friction log" item 1 and §"Patterns observed" P1; PRA §3 "What went wrong" item 1 and §4b §9.1 trigger 4 (loaded skill failed to prevent a finding) — `docs/gamma/cdd/3.79.0/POST-RELEASE-ASSESSMENT.md`
- **Class:** `cdd-skill-gap`
- **Trigger:** §9.1 trigger 4 — loaded skill failed to prevent a finding it covers (`alpha/SKILL.md` §2.6 row 14 path (a) was loaded by α to remediate a worktree-inherited identity leak; the rebase invalidated 9 SHA citation sites in `self-coherence.md` that α had already authored)
- **Description:** Row 14 of `alpha/SKILL.md` §2.6 (pre-review gate) defines path (a) as the retroactive fix for a non-canonical α commit-author email: `git config user.email "alpha@{project}.cdd.cnos" && git rebase -i {merge-base} --exec 'git commit --amend --reset-author --no-edit' && git push --force-with-lease origin cycle/{N}`. The rebase rewrites every α commit's SHA. Row 14's prescription is mechanically complete (the rebase works; the email is corrected; β + γ commits are preserved). What row 14 does *not* name is the downstream consequence: any SHA-bearing artifact authored *before* the rebase (typically `self-coherence.md` §CDD Trace step table, §ACs header references like "Per-AC oracles run against branch HEAD `<SHA>`", §Debt "at HEAD `<SHA>`" citations) references pre-rebase SHAs that no longer exist on `origin/cycle/{N}` after the rebase. Those SHAs live only in α's local reflog. β's `review/SKILL.md` rule 3.13(a) reproducibility check (a future reviewer must be able to resolve cited pointers) catches the invalidation at the next round.
- **Root cause:** Row 14 was authored to close the cycle #287 R1 F3 finding (session-start email drift surviving all 10 prior pre-review-gate rows). Its empirical anchor was an end-of-cycle identity correction; the patterns it considered did not include "α has already authored SHA-bearing artifacts before the rebase fires." The downstream consequence is a class of artifact-integrity violations that the row's mechanical prescription enables without naming. Combined with `alpha/SKILL.md` §2.3 intra-doc repetition rule (which catches additional occurrences post-fact), row 14's mechanics are sound but the role-skill's coverage of the path-(a) timing-vs-authoring-order is incomplete.
- **Disposition:** **patch-landed**
- **Patch:** this PRA commit — `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 — added a new auxiliary paragraph after row 14's main text titled **"SHA citations across path (a) rebase."** The paragraph names the SHA-citation-invalidation consequence and prescribes two coherent resolution paths: (i) run identity verification at session-start, before any SHA-bearing artifact is authored — the rebase has nothing to rewrite (preferred); (ii) re-stamp every SHA citation in `self-coherence.md` and other authored artifacts to the current-branch SHA immediately after the rebase, applying the §2.3 intra-doc repetition rule to catch every occurrence. Empirical anchor cited inline (`#380 R1 F1`).
- **Affects:** `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 (auxiliary paragraph following row 14)

## Summary

| Finding | Class | Disposition | Patch / MCA path |
|---------|-------|-------------|------------------|
| F1 — alpha/SKILL.md §2.6 row 14 under-specifies SHA-citation invalidation after path (a) rebase | cdd-skill-gap | patch-landed | this PRA commit (alpha/SKILL.md §2.6 auxiliary paragraph following row 14) |

## Related observations (non-cdd-gap; recorded for traceability, not aggregated)

These items surfaced during this cycle's triage but do not qualify as `cdd-*-gap` findings; they are project-level or non-CDD-protocol-level and live in the PRA's §2 lag table and §7 deferred outputs respectively:

- **Pre-existing I4 / notify CI red on merge commit** (PRA §3 "What went wrong" item 2; §4b §9.1 trigger CI-red-on-merge): inherited from cycle #369; project-debt remediation MCA filed via PRA §7 next-move list. Not a cdd-protocol-gap — the §3.8 cap clause is operating as a forcing function as designed.
- **#380 P5 — worktree-config identity-leak class is repo-level** (#379 F-item-1, #370 F4, #380 F1 + β-closeout §"Identity hygiene"): converged in issue #373 as a preventive `--worktree` identity write across role skills. Escalated from "growing" to "escalating" lag in PRA §2; γ recommends prioritizing #373 in a near-term cycle.

## Aggregator update

`.cdd/iterations/INDEX.md` row added at cycle #380 close (this commit); after the cycle dir move per `gamma/SKILL.md §2.6`, this artifact's path is `.cdd/releases/3.79.0/380/cdd-iteration.md`.
