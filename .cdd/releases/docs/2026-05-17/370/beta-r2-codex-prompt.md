You are β (round 2). Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 370 --json title,body,state,comments
Branch: cycle/370

Round 1 verdict: REQUEST CHANGES (your own R1 review at `.cdd/unreleased/370/beta-review.md`). Single finding F1 (mechanical, B-severity, ci-status): `## CDD-Trace` → `## CDD Trace` rename in `.cdd/unreleased/370/self-coherence.md`.

α R2 commit: `846800e8` on `cycle/370`. α reports `cn-cdd-verify --unreleased` PASS locally and adds an §R2 entry to self-coherence.

Git identity:
  git config user.name "beta"
  git config user.email "beta@cdd.cnos"
  # also assert at worktree scope per R1 §2.1 row 1 note (extensions.worktreeConfig may still be set from R1 merge-test):
  git config --worktree user.name beta 2>/dev/null || true
  git config --worktree user.email beta@cdd.cnos 2>/dev/null || true

R2 review surface (narrow, per R1 §Notes):
- Re-run AC1: `test -f src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`
- Re-run AC7: awk + rg recipe from issue body
- Re-run AC9: `git diff origin/main..HEAD --stat` (must still be exactly {new doc + cycle evidence})
- Re-run `cn-cdd-verify --unreleased` on cycle branch
- Verify F1 closed: `## CDD Trace` (with space) in `.cdd/unreleased/370/self-coherence.md` §CDD Trace section; manifest comment aligned
- Verify Build CI green on new review SHA (run id discoverable via `gh run list --branch cycle/370 --limit 5`)
- Run pre-merge gate row-by-row (`beta/SKILL.md` §pre-merge gate) — especially row 3 non-destructive merge-test, since R1 row 3 failed

If all pass: APPROVE R2. Merge `cycle/370` into `main` per `beta/SKILL.md` §1 (`git merge --no-ff cycle/370` with `Closes #370` in merge body). Write `beta-review.md` §"Round 2" and `beta-closeout.md` to `.cdd/unreleased/370/`. Exit. Do NOT tag — γ-acting-as-δ owns the release boundary this cycle.

If anything else fails (not F1; not the merge-test): return RC R2 with the specific finding. R3 is uncommon for docs-only but legitimate if a new issue surfaces.
