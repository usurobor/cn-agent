# α closeout — cycle/426

**Issue:** [cnos#426](https://github.com/usurobor/cnos/issues/426) — Publish v3.82.0 GH release directly (second use of remote-runner-delegation doctrine; sidesteps broken release.yml build).

**Mode:** workflow + receipt only (doctrine already on main from cnos#425); γ+α+β-collapsed-on-δ.

## Matter produced

Two-deliverable bundle (workflow + receipt), in commit `5d237d2e`:

1. **`.github/workflows/publish-3.82.0-release.yml`** — new file, 80 lines. One-shot, push-triggered workflow that publishes the v3.82.0 GitHub Release directly via `softprops/action-gh-release@v1`, sidestepping `release.yml`'s broken `build` job. Triggers on `push: branches: [main], paths: ['.github/workflows/publish-3.82.0-release.yml']` so it fires only when the cycle/426 merge lands (doctrine already on main from cnos#425). Header comments (lines 1–32) cite essay + skill section + receipt + issue + cnos#425 cdd-iteration findings as motivation. Four substantive steps: (1) `actions/checkout@v4` with `ref: '3.82.0'` (load-bearing — pins the body source to the tag commit, not main HEAD which has CELL-OF-CELLS + BOX-AND-THE-RUNNER post-tag); (2) RELEASE.md verification guard (`head -1 RELEASE.md` + `grep -q '^# v3.82.0'`); (3) `softprops/action-gh-release@v1` with `tag_name: '3.82.0'`, `name: '3.82.0 — CCNF package-architecture baseline'`, `body_path: RELEASE.md`, `draft: false`, `prerelease: false`; (4) self-delete via `git fetch origin main` + `git checkout main` + `git rm` + `git commit` + `git push origin main` (worktree switch is required because step 1 checks out the tag rather than main — pattern difference vs cnos#425's workflow).

2. **`.cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md`** — new file, 196 lines. Second instantiated 6-field receipt under the remote-runner doctrine (first use: cnos#425). All six fields filled: §1 Who asked (operator + cnos#425 follow-on directive after release.yml build-failure diagnosis); §2 What runs (4 substantive job steps with specific git/softprops commands); §3 Where (`ubuntu-latest` GH-hosted ephemeral VM); §4 Authority (`GITHUB_TOKEN` with `contents: write` scoped to `usurobor/cnos`); §5 Evidence (post-run-fillable per hard rule 1 with 6-field evidence-shape table declared in advance); §6 Who may accept (operator with release-body-match acceptance criterion). Plus expected-effect (5 steps), failure-modes (7 rows with mitigations including release-already-exists idempotency case + wrong-ref recovery), acceptance-criteria (7 numbered + partial/rejection branches), V/δ composition note, and a Relationship-to-cnos#425 section explicitly framing the two cycles as a doctrine-governed sequence.

## Authoring discipline

The issue body's implementation contract pinned all values at dispatch time. α's role was mechanical instantiation under the pinned values:

- **Workflow YAML shape verbatim** from the issue body D1 spec, modulo header comments (33-line header expansion citing doctrine + skill + receipt + issue + cnos#425 motivation + self-delete worktree-switch rationale).
- **Receipt 6-field structure verbatim** from the issue body D2 spec, with each field substantively populated rather than summarized.
- **Tag SHA preserved** — `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e` (verified at branch creation, inherited from cnos#425 retarget). This cycle does NOT touch the tag.
- **RELEASE.md header preserved** — `git show 3.82.0:RELEASE.md | head -1` returns `# v3.82.0 — CCNF package-architecture baseline` (verified at branch creation). The workflow's step-2 verification is a defense-in-depth guard against drift.
- **No doctrine modification** — `BOX-AND-THE-RUNNER.md`, `delta/SKILL.md §8`, and the README are all byte-identical to origin/main; they are inherited from cnos#425.

The 33-line header in the workflow file (longer than cnos#425's 23-line header) is justified because the second-use needs to explain:

- *Why* a second one-shot workflow exists (the cnos#425 retarget worked but release.yml did not retrigger and its build job is broken);
- *Why* the workflow checks out at the tag rather than main (main HEAD has post-tag content that would publish the wrong body);
- *Why* the self-delete needs a worktree switch (consequence of the at-tag checkout);
- *Why* no binaries are attached (release.yml fix is a separate cycle, filed as cdd-iteration F1).

This is doctrine-discipline verbosity (per the cnos#425 cdd-iteration "doctrine essay tense is descriptive, not prescriptive" observation), not bloat.

## Mechanical AC verification

All 10 ACs PASS per `self-coherence.md`:

| AC | Verdict | Evidence |
|----|---------|----------|
| AC1 | PASS | workflow has `name`, `on:`, `permissions:`, `jobs:`; 4 named/uses steps |
| AC2 | PASS | `grep "ref:.*3.82.0"` returns 1 hit; checkout pins the tag |
| AC3 | PASS | `grep "body_path: RELEASE.md"` returns 1 hit |
| AC4 | PASS | receipt has 6 numbered field sections, all populated; §5 post-run-fillable per hard rule 1 |
| AC5 | PASS | workflow header (lines 1–32) cites essay + skill §8 + receipt path + issue + cnos#425 context |
| AC6 | PASS | final step has `git rm` + `git commit` + `git push origin main` (plus worktree switch prereq) |
| AC7 | PASS | CDD.md 0 lines diff; cds/cdr/handoff 0 lines diff; entire cnos.cdd/ 0 lines diff |
| AC8 | PASS | schemas/, src/go/, cnos.core/, cnos.eng/, cnos.kata/, cnos.cdd.kata/, scripts/ all 0 lines diff |
| AC9 | PASS | docs/ 0 files modified (essay + README inherited from cnos#425) |
| AC10 | PASS | release.yml 0 lines diff (build-job fix deferred to next-MCA per cdd-iteration F1) |

## Refusal conditions held

All hard rules from the dispatch brief held:

- **Receipt before execution.** 6 fields filled; evidence post-run-fillable (§5).
- **One-shot self-deleting workflow.** Final step is `git rm` + `git commit` + `git push origin main` (with worktree switch).
- **No release.yml fix in this cycle.** Filed as cdd-iteration F1 with disposition `next-MCA`.
- **No CCNF kernel / CDS / CDR / handoff / cnos.cdd changes.** Verified by mechanical diff (0 lines across CDD.md, cds/, cdr/, handoff/, cnos.cdd/).
- **Workflow checks out at the TAG (3.82.0), not main HEAD.** `ref: '3.82.0'` pinned in step 1; load-bearing because main HEAD has CELL-OF-CELLS + BOX-AND-THE-RUNNER post-tag content that would publish the wrong body.
- **No essay or README edits.** Verified by mechanical diff (0 docs/ files modified).

## Commit

- **α-426** (`5d237d2e`): "α-426: publish-3.82.0-release workflow + remote-runner receipt" — both deliverables in a single commit per the authoring-order rule (`delta/SKILL.md §8.2`). The doctrine inherited from cnos#425 is already on main, so the bundle does not need to land doctrine simultaneously (this is the doctrine-already-on-main precedence pattern this cycle establishes for second-and-later uses of any cnos primitive — first use ships doctrine + first artifact + first receipt atomically; second-and-later uses ship artifact + receipt atomically with doctrine inherited).

## Handoff

Self-coherence complete; β-review (R1) APPROVED. γ to file gamma-scaffold, gamma-closeout, cdd-iteration (substantive: 2 findings — release.yml build-auth gap + tag-force-update non-trigger, both `next-MCA`), and INDEX.md row, then push branch.

Filed by α@cnos (γ+α+β-collapsed-on-δ) on 2026-05-24.
