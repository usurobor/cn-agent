# γ scaffold — cycle/426

**Issue:** [cnos#426](https://github.com/usurobor/cnos/issues/426) — Publish v3.82.0 GH release directly (second use of remote-runner-delegation doctrine; sidesteps broken release.yml build).

**Mode:** workflow + receipt only (doctrine already on main from cnos#425); γ+α+β-collapsed-on-δ (per `ROLES.md §4`-precedent for mechanical-AC + receipt-discipline cycles; α=β permitted because the primary product is YAML + Markdown with grep-checked and shape-checked invariants, paralleling cnos#425's first-use precedent which itself inherited the cycle-414 / cycle-424 essay precedents).

**Branch:** `cycle/426` from `origin/main` (`144db7f1` — the post-cnos#425-merge `release: remove one-shot repoint-3.82.0 workflow (cnos#425)` commit).

## Surface (2-deliverable bundle + close-out artifacts)

Two-deliverable bundle authored under the **doctrine-already-on-main precedence pattern** (the second-use specialization of the doctrine-before-first-use rule from cnos#425 — when the doctrine is already on main from a prior cycle, the second-use cycle ships only the artifact + receipt atomically; doctrine inheritance from main is implicit):

- **D1:** `.github/workflows/publish-3.82.0-release.yml` (new; 80 lines). One-shot, push-triggered workflow. Triggers on `push: branches: [main], paths: ['.github/workflows/publish-3.82.0-release.yml']`. Four substantive job steps (checkout at tag `3.82.0` / RELEASE.md verification / softprops/action-gh-release@v1 invocation / self-delete with worktree switch). Header comments (lines 1–32) cite doctrine, skill section, receipt, issue, and the cnos#425 cdd-iteration findings (F1 release.yml build-auth + F2 tag-force-update non-trigger) as second-use motivation.

- **D2:** `.cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md` (new; 196 lines). Second instantiated 6-field receipt under the remote-runner doctrine. All 6 fields populated (§5 Evidence is post-run-fillable per hard rule 1; the shape is named in advance as a 6-row evidence-shape table). Plus expected-effect (5 steps), failure-modes (7 rows with mitigations including release-already-exists idempotency case + wrong-ref recovery), acceptance-criteria (7 numbered + partial/rejection branches), V/δ composition note, and a Relationship-to-cnos#425 section that explicitly frames the two cycles as a doctrine-governed sequence.

Close-out artifacts at `.cdd/unreleased/426/`:

- `gamma-scaffold.md` (this file)
- `self-coherence.md` (α-side AC verification, all 10)
- `beta-review.md` (β-side R1 review verdict; 10 binding findings, all PASS)
- `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` (per-role closeouts)
- `cdd-iteration.md` (**substantive** — records 2 `cdd-tooling-gap` findings: F1 release.yml build-job git-auth failure on tag-push, disposition `next-MCA`; F2 GH Actions `on.push.tags` does not reliably trigger on tag force-update, disposition `next-MCA`)
- INDEX.md row appended at `.cdd/iterations/INDEX.md`

## Implementation contract (pinned by δ at dispatch; verified at this scaffold)

| Axis | Pinned value | Conforms? |
|---|---|---|
| Language | YAML (workflow) + Markdown (receipt) | yes |
| CLI integration target | None new; uses softprops/action-gh-release@v1 (already used by release.yml's publish job) | yes |
| Package scoping | New workflow `.github/workflows/publish-3.82.0-release.yml`; new receipt `.cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md`; close-out artifacts | yes (exact 2-deliverable + standard close-out bundle) |
| Runtime dependencies | GitHub Actions ubuntu-latest runner; `GITHUB_TOKEN` with `contents: write` | yes (encoded in workflow `permissions:` + receipt §3 + §4) |
| Backward compat | No existing files modified | yes (mechanical diff: 0 lines changed across CDD.md, cds/, cdr/, handoff/, cnos.cdd/, docs/, schemas/, runtime/, scripts/, release.yml) |
| Doctrine-already-on-main precedence | Doctrine (essay + skill §8) inherited from cnos#425; this cycle ships only artifact + receipt atomically; workflow trigger path-filtered on own path so it fires only post-merge | yes (commit `5d237d2e` atomic; trigger pinned; doctrine commit `7720c441` already on main) |

## AC oracle approach (issue body verbatim)

| AC | Oracle (mechanical) | Result |
|----|---------------------|--------|
| AC1 | workflow has `name`, `on:`, `permissions:`, `jobs:`; 4 substantive steps | PASS (4 named/uses steps) |
| AC2 | `grep "ref:.*3.82.0"` ≥ 1 hit | PASS (1) |
| AC3 | `grep "body_path: RELEASE.md"` returns 1 hit | PASS (1) |
| AC4 | receipt has 6 numbered field rows; evidence post-run-fillable | PASS (6 sections, §5 placeholder) |
| AC5 | workflow header cites essay + skill §8 + receipt path | PASS (lines 1–32) |
| AC6 | final step has `git rm` + `git commit` + `git push origin main` | PASS (+ worktree-switch prereq) |
| AC7 | CDD.md / cds/ / cdr/ / handoff/ / cnos.cdd/ all 0 lines diff | PASS (0/0/0/0/0) |
| AC8 | schemas/ / src/go/ / cnos.core/ / cnos.eng/ / cnos.kata/ / cnos.cdd.kata/ / scripts/ all 0 lines diff | PASS (0) |
| AC9 | docs/ 0 files modified | PASS (empty) |
| AC10 | release.yml 0 lines diff | PASS (0) |

## Branch + commit shape

- **α-426** (`5d237d2e`): "α-426: publish-3.82.0-release workflow + remote-runner receipt" — both deliverables in a single atomic commit per the doctrine-already-on-main precedence pattern.
- **β-426** (`4ec96450`): "β-426: R1 review APPROVED (10/10 PASS) + role closeouts (α, β)" — beta-review.md (10 binding findings, all PASS) + alpha-closeout.md + beta-closeout.md.
- **γ-426** (forthcoming this commit): "γ-426: close-outs (γ + scaffold + self-coherence) + substantive cdd-iteration + INDEX row" — gamma-scaffold.md (this file) + self-coherence.md + gamma-closeout.md + cdd-iteration.md (substantive; 2 findings) + INDEX.md row.

Push to `origin/cycle/426`; do NOT merge to main (operator's call). Merge instruction reported in γ closeout.

## Critical refusal conditions surfaced during authoring

- **Doctrine-already-on-main precedence held.** Both D1 + D2 land in commit `5d237d2e`; doctrine (essay `BOX-AND-THE-RUNNER.md`, skill `delta/SKILL.md §8`, README pointer) is byte-identical to origin/main (verified by `git diff origin/main -- docs/ src/packages/cnos.cdd/` returning 0). Workflow's push trigger is path-filtered on own path so it fires only on the cycle/426 merge to main — at which point all components are on main.

- **Checkout at tag, not main HEAD.** The workflow's step 1 pins `ref: '3.82.0'` (the literal tag, which resolves to commit `fd1d654e` per `git ls-remote origin refs/tags/3.82.0`). Critical because:
  - main HEAD (`144db7f1`) has CELL-OF-CELLS post-tag (via cycle/424) AND BOX-AND-THE-RUNNER doctrine post-tag (via cycle/425);
  - main HEAD's root `RELEASE.md` is NOT the v3.82.0 baseline body — it would publish the wrong content;
  - the tag commit `fd1d654e`'s root `RELEASE.md` begins with `# v3.82.0 — CCNF package-architecture baseline` (verified by `git show 3.82.0:RELEASE.md | head -1`);
  - the workflow's step 2 verification (`head -1 RELEASE.md` + `grep -q '^# v3.82.0'`) is a defense-in-depth guard.

- **Self-delete worktree switch is load-bearing.** Because step 1 checks out the tag, the job's working tree is detached at the tag, not on main. Step 4 (self-delete) must `git fetch origin main` + `git checkout main` BEFORE `git rm` + `git commit` + `git push origin main`, otherwise the commit lands on a detached HEAD and the push fails. cnos#425's workflow did NOT need this switch (its checkout was implicitly at main). The pattern difference is noted in β-review B6 and the relationship-to-cnos#425 section of the receipt.

- **No CCNF kernel / CDS / CDR / handoff / cnos.cdd changes.** Verified by mechanical diff. `src/packages/cnos.cdd/` entirely byte-identical to origin/main (delta/SKILL.md §8 inherited from cnos#425).

- **No essay or README edits.** Verified by mechanical diff (`docs/` 0 files modified). The doctrine essay `BOX-AND-THE-RUNNER.md` and the README pointer are inherited unchanged from cnos#425.

- **No release.yml modification.** Verified by mechanical diff (0 lines). The release.yml build-job fix is deferred to a separate cycle and recorded as cdd-iteration F1 with disposition `next-MCA`.

- **No binaries attached to the release.** Pure release-notes publication. The `softprops/action-gh-release@v1` invocation has only `tag_name`, `name`, `body_path`, `draft`, `prerelease` — no `files:` key. Binary post-attachment is deferred to a follow-on cycle after cdd-iteration F1 (release.yml build-job repair) ships, if needed.

- **Two findings filed substantively in cdd-iteration.** F1 (release.yml build-job git-auth failure on tag-push) and F2 (GH Actions `on.push.tags` non-trigger on tag force-update) are both `cdd-tooling-gap`-class findings with disposition `next-MCA`. Each has a first AC specified per the issue body's requirement; both will appear in `.cdd/iterations/INDEX.md` as `findings=2, patches=0, MCAs=2, no-patch=0`.

- **Citation closed loop across the 2 deliverables.** Workflow header → cites doctrine + skill + receipt + issue + cnos#425 motivation. Receipt header → cites doctrine + skill + workflow + cycle. The Relationship-to-cnos#425 section in the receipt makes the two cycles' connection explicit. β-review B5 verified the workflow-side citation; the receipt-side citation is verified by direct read of the receipt's frontmatter.

Filed by γ@cnos (γ+α+β-collapsed-on-δ) on 2026-05-24.
