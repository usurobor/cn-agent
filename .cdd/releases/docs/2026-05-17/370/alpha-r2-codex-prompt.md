You are α (round 2). Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 370 --json title,body,state,comments
Branch: cycle/370

Round 1 verdict: REQUEST CHANGES (β at SHA d0ea77f9). Read `.cdd/unreleased/370/beta-review.md` §Findings + §"Required fix for F1" for the binding finding.

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

R2 scope (single binding finding, mechanical):

F1 — `cn-cdd-verify` requires `## CDD Trace` (space), not `## CDD-Trace` (hyphen). The Build CI workflow is red on review SHA `aa10f902` for this reason (`Build` run `25990518085`, job `CDD artifact ledger validation (I6)`).

Required edits in `.cdd/unreleased/370/self-coherence.md`:
1. Rename the section header `## CDD-Trace` → `## CDD Trace` (single occurrence near line 242).
2. Align the manifest comment at line 1: `<!-- sections: [...CDD-Trace...] -->` → `<!-- sections: [...CDD Trace...] -->`, and the same in the `<!-- completed: ... -->` line.
3. Add an §R2 entry to self-coherence (or append to existing §Review-readiness) noting: R1 finding F1, the single-character fix, oracle re-run confirmation (`cn-cdd-verify --unreleased` PASS locally; CI re-run expected green).
4. Commit on `cycle/370` with message `α-370 R2: rename §CDD-Trace → §CDD Trace (align with cn-cdd-verify convention; closes F1)`.
5. Push to `origin/cycle/370`.

After the rename, run locally:
  src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --unreleased

It must report PASS. If anything else turns red, surface immediately — but β's review confirms F1 is the only binding finding and AC1, AC7, AC9 oracles all pass at HEAD.

Do NOT:
- modify `COHERENCE-CELL-NORMAL-FORM.md` (it passed all 9 ACs; R2 is mechanical evidence-doc fix only)
- re-run the full AC bank (β will re-run AC1, AC7, AC9 + validator on the new review SHA)
- touch any other file (AC9 surface containment still binding)

Signal review-readiness for R2 in self-coherence.md once the validator passes locally and the commit is pushed.
