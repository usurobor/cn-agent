# Remote-runner receipt — v3.82.0 GH release publication (cnos#426)

**Doctrine:** [`docs/gamma/essays/BOX-AND-THE-RUNNER.md`](../../../docs/gamma/essays/BOX-AND-THE-RUNNER.md)
**Skill:** [`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/delta/SKILL.md) §8 (Remote-runner delegation)
**Artifact:** [`.github/workflows/publish-3.82.0-release.yml`](../../../.github/workflows/publish-3.82.0-release.yml)
**Cycle:** cnos#426 — Publish v3.82.0 GH release directly (second use of remote-runner-delegation doctrine; sidesteps broken release.yml build)
**Date:** 2026-05-24

This receipt is the **second instantiation** of the 6-field remote-runner receipt convention (first use: cnos#425 — tag retarget). The doctrine that governs it (`BOX-AND-THE-RUNNER.md` essay + `delta/SKILL.md §8`) already lives on main from cnos#425; this cycle is a pure workflow + receipt + close-outs cycle that does not modify the doctrine. The receipt is authored in the same commit-set as the workflow it accompanies, per the authoring-order rule in `delta/SKILL.md §8.2`. The `evidence` field (§5) is a post-run-fillable placeholder; the run cannot happen until the cycle merges and the push trigger fires.

---

## 1. Who asked for it

**Operator** (cnos γ session, 2026-05-24), via the release-boundary directive following cnos#425's first-use diagnosis:

- "execute as delta or gamma" (carried forward from the cnos#425 session) — operator authorizing the agent to act under δ-as-role authority for the release-fix boundary.
- The cnos#425 first-use exposed two infrastructure findings (recorded in this cycle's `cdd-iteration.md`): (F1) `release.yml`'s `build` job fails on tag-push with a git-auth prompt, blocking the standard publish path; (F2) `on.push.tags` does not reliably re-fire on tag force-update, so the cycle-425 tag retarget did not cascade into a fresh `release.yml` run.
- Operator's follow-on directive for cycle/426: publish the v3.82.0 GH release directly via the remote-runner primitive, sidestepping the broken `release.yml` build path; defer the `release.yml` fix to a separate cycle.

The cnos#426 issue body codifies this directive as the cycle's mandate. No agent-initiated escalation; the move is operator-directed and explicitly second-use of the doctrine.

---

## 2. What it will run

The workflow's four job steps, in order:

```text
1. actions/checkout@v4 with:
     ref: '3.82.0'           (checkout AT THE TAG, not main HEAD —
                              the tag commit fd1d654e carries the
                              canonical v3.82.0 RELEASE.md body)
     fetch-depth: 0          (full history so the self-delete step
                              can switch to main and push)
     token: GITHUB_TOKEN     (substrate-provided; contents: write)

2. Verify root RELEASE.md is v3.82.0:
     head -1 RELEASE.md
     grep -q '^# v3.82.0' RELEASE.md || (echo ERROR; exit 1)
   Guard against a wrong-tag checkout publishing a wrong body.

3. softprops/action-gh-release@v1 with:
     tag_name: '3.82.0'
     name: '3.82.0 — CCNF package-architecture baseline'
     body_path: RELEASE.md   (sourced from the tag's checked-out tree)
     draft: false
     prerelease: false
   Creates (or updates) the GH release at refs/tags/3.82.0 with the
   v3.82.0 body. No binary artifacts are attached (release.yml fix
   is a separate cycle; see cdd-iteration F1).

4. Self-delete the workflow file:
     git fetch origin main
     git checkout main
     git config user.name  "cnos-remote-runner[bot]"
     git config user.email "cnos-remote-runner@cnos.local"
     git rm .github/workflows/publish-3.82.0-release.yml
     git commit -m "release: remove one-shot publish-3.82.0-release workflow (cnos#426)"
     git push origin main
```

The release-creation step (step 3) is the load-bearing effect: it publishes (or updates) the GitHub Release at the `3.82.0` tag with the canonical v3.82.0 body. Because the checkout in step 1 pins `ref: '3.82.0'` (commit `fd1d654e`), the `RELEASE.md` body_path resolves to the v3.82.0 baseline body, NOT to main HEAD (which now has CELL-OF-CELLS + BOX-AND-THE-RUNNER content that landed post-tag and would publish the wrong body).

The self-delete step (step 4) closes the latent-execution authority of the workflow itself, per `delta/SKILL.md §8.3` (one-shot workflows self-delete). It switches to a fresh `main` worktree to commit + push because the job's working tree is on the tag, not main.

---

## 3. Where it will run

**GitHub Actions runner: `ubuntu-latest`** (GitHub-hosted ephemeral VM).

Not a self-hosted runner. Not a third-party CI provider. The runner is a single-use VM provisioned by GitHub for this workflow's job; it is destroyed after the job completes. Same execution surface as cnos#425's first-use workflow (`repoint-3.82.0.yml`) — the doctrine treats two consecutive uses on the same runner class as evidence-class identical.

---

## 4. What authority it has

**`GITHUB_TOKEN`** with **`contents: write`** permission, scoped to `usurobor/cnos`.

No other secrets are accessed by this workflow. No deploy keys. No cloud OIDC role assumptions. No third-party tokens. The workflow does not invoke `gh` CLI commands that would require additional scopes; both the release-creation step (via `softprops/action-gh-release@v1`, which uses `GITHUB_TOKEN` automatically) and the self-delete commit-and-push succeed under `contents: write` alone.

The `contents: write` scope grants:
- create / update GitHub Releases (used by step 3)
- push commits to `main` (used by step 4 self-delete)
- read repository contents at any ref (used by step 1 checkout at tag)

The scope does NOT grant:
- write access to other repositories under `usurobor/`
- access to repository secrets beyond `GITHUB_TOKEN`
- workflow-editing capability beyond this repo (the self-delete commit edits only this workflow's own file)

Branch-protection on `main` is enforced by the substrate; the workflow's self-delete commit is subject to whatever rules apply to direct pushes to `main` from `GITHUB_TOKEN`-authenticated commits. If branch-protection rejects the self-delete push, the release publication from step 3 still succeeds and the workflow file persists on `main` until manually removed; this is a named failure mode (§Failure modes below) and is the same pattern resolved manually in cnos#425's first use (commit `144db7f1` — `release: remove one-shot repoint-3.82.0 workflow (cnos#425)`).

---

## 5. Evidence

**Post-run-fillable placeholder.** The receipt is authored before the run; the evidence shape is named in advance so a downstream consumer can detect a missing receipt rather than reading silence as success.

Required evidence (to be filled after the run):

| Field | Shape | Value (post-run) |
|---|---|---|
| Workflow run URL | `https://github.com/usurobor/cnos/actions/runs/<run-id>` | _(post-run)_ |
| Published release URL | `https://github.com/usurobor/cnos/releases/tag/3.82.0` | _(post-run)_ |
| Release body matches v3.82.0 | First line of release body equals `# v3.82.0 — CCNF package-architecture baseline` | _(post-run)_ |
| Release tag SHA | `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e8d6361f0db2a6407f19e573a96d1054d` (unchanged from cnos#425 retarget) | _(post-run)_ |
| Self-delete commit SHA | SHA of the `release: remove one-shot publish-3.82.0-release workflow (cnos#426)` commit on `main` | _(post-run)_ |
| Workflow file gone | `.github/workflows/publish-3.82.0-release.yml` no longer present on `main` | _(post-run)_ |

The acceptor (§6) fills these fields when declaring acceptance. The receipt then moves with the cycle into `.cdd/releases/3.82.0/` (or the next release directory) per standard cycle-close artifact handling.

---

## 6. Who may accept the result

**Operator.** Acceptance criterion:

> Operator inspects `https://github.com/usurobor/cnos/releases/tag/3.82.0` and confirms that the release body matches `.cdd/releases/3.82.0/RELEASE.md` (the canonical v3.82.0 release notes — the "CCNF package-architecture baseline" body), which is mirrored at root `RELEASE.md` on the v3.82.0 tag commit (`fd1d654e`). The first line of the body must read `# v3.82.0 — CCNF package-architecture baseline`.

The operator's δ-decision is the move from "the run completed" to "the run is accepted." A run that completes but produces a release body that does not match `.cdd/releases/3.82.0/RELEASE.md` (e.g. because the checkout pinned the wrong ref or the body file changed between authoring and run) is a completed-but-rejected outcome; the operator records the rejection and invokes the manual recovery path (§Failure modes).

No agent — α, β, γ, δ-as-role, or ε — has acceptance authority for this remote-runner move. The acceptance is at the project-binding / release-boundary layer (`ROLES.md` §5-layer chain) and is held by the operator at the outermost cell (per `CELL-OF-CELLS.md §13`). Inherited unchanged from cnos#425's receipt §6.

---

## Expected effect (5 steps, in order)

1. **Merge of cycle/426 to main.** Operator merges the cycle's branch; the merge commit lands `.github/workflows/publish-3.82.0-release.yml` (plus this receipt and the close-out artifacts) onto `main`.
2. **Push trigger fires.** The push to `main` includes a change to `.github/workflows/publish-3.82.0-release.yml` (because the file is newly added); the path filter matches; the workflow's `publish` job is scheduled.
3. **Checkout at tag + body verification.** Step 1 of the job checks out the repository at `ref: '3.82.0'` (commit `fd1d654e`). Step 2 verifies the first line of root `RELEASE.md` is `# v3.82.0` and exits non-zero if it is not.
4. **GitHub Release published.** Step 3 invokes `softprops/action-gh-release@v1` with `tag_name: '3.82.0'`, `name: '3.82.0 — CCNF package-architecture baseline'`, `body_path: RELEASE.md`. The release is created (or its body updated if a stub already exists) at `https://github.com/usurobor/cnos/releases/tag/3.82.0` with the canonical v3.82.0 body.
5. **Self-delete.** Step 4 of the job switches to a fresh `main` worktree, runs `git rm` + `git commit` + `git push origin main`, removing `.github/workflows/publish-3.82.0-release.yml` from `main`. The latent-execution authority of this workflow closes. The self-delete commit on `main` is the structural signal that the one-shot ran.

---

## Failure modes + mitigations

| Failure | Detection | Mitigation |
|---|---|---|
| Branch-protection on `main` rejects the workflow's self-delete push | Workflow logs show `push origin main` failure in step 4 | Operator manually removes `.github/workflows/publish-3.82.0-release.yml` via PR (a normal cycle, not under this receipt). Release publication from step 3 already succeeded, so the release is correct; only cleanup remains. This is the same pattern cnos#425 first-use resolved manually via commit `144db7f1`. |
| RELEASE.md verification fails (wrong tag or body drift) | Workflow logs show step 2 exits non-zero with `ERROR: RELEASE.md is not v3.82.0` | Workflow fails fast before publishing a wrong body. Operator investigates whether the tag commit is `fd1d654e` (use `git ls-remote origin refs/tags/3.82.0`); if the tag moved unexpectedly, run cnos#425's retarget pattern; if RELEASE.md at the tag is wrong, file a fix-cycle to correct the body at the tag. |
| softprops/action-gh-release@v1 fails (action unavailable, network, rate limit) | Workflow logs show step 3 failure with action-side error | Operator re-runs the workflow via `gh run rerun <run-id>` (the file is still on main until step 4 completes, so the trigger can be re-fired by an empty-commit touch to the workflow path — or by direct re-run from the Actions UI). If the action is permanently unavailable, fall back to `gh release create 3.82.0 --notes-file RELEASE.md` invoked from operator's local under operator credentials. |
| GITHUB_TOKEN lacks `contents: write` at runtime (org policy override) | Workflow logs show 403 on release create or self-delete push | Operator escalates to GitHub org settings; if temporarily unavailable, falls back to manual `gh release create` under operator credentials. |
| Release already exists (release.yml partial run from cnos#425 era left a stub) | Step 3 attempts to create a release that exists | `softprops/action-gh-release@v1` will update the existing release in-place (its default behavior) — body and name are overwritten with the values in this workflow. This is the desired behavior; not a failure. If the existing release has shipped binaries that should be preserved, the action's default does not delete them (it only updates body/name) — verify this in step 3's logs post-run. |
| Workflow does not trigger (path filter mismatch, file permissions) | Operator merges the cycle but no workflow run appears in Actions tab | Operator inspects the path filter against the file path; if mismatch is found, file a follow-up cycle to fix the trigger or invoke manually via `workflow_dispatch` (this workflow has no `workflow_dispatch` trigger; would require an edit-cycle to add one). Release publication can also be done manually per the operator-credential fallback above. |
| Wrong-ref checkout publishes wrong body | Operator inspects release body and finds it does not match `.cdd/releases/3.82.0/RELEASE.md` | Operator declares the run completed-but-rejected; investigates which ref `softprops/action-gh-release@v1` sourced RELEASE.md from (step 1 logs show the checkout's ref); re-runs the workflow with a corrected ref, or manually edits the GH release body. The step-2 verification guard should prevent this failure mode in practice. |

---

## Acceptance criteria

This remote-runner move is accepted when **all** of the following hold:

1. The workflow run for `.github/workflows/publish-3.82.0-release.yml` completed successfully (all 4 steps green).
2. The GitHub Release at `https://github.com/usurobor/cnos/releases/tag/3.82.0` exists and is not draft / not prerelease.
3. The release name is `3.82.0 — CCNF package-architecture baseline`.
4. The release body matches `.cdd/releases/3.82.0/RELEASE.md` (mirrored at root `RELEASE.md` on tag `fd1d654e`). The first line of the body reads `# v3.82.0 — CCNF package-architecture baseline`.
5. `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e8d6361f0db2a6407f19e573a96d1054d` (unchanged from cnos#425's retarget — this cycle does not touch the tag).
6. `.github/workflows/publish-3.82.0-release.yml` is no longer present on `main` (self-delete succeeded).
7. Operator records acceptance by filling the §5 evidence table in this receipt and moving the receipt into the appropriate release directory.

Partial acceptance (1, 2, 3, 4, 5 pass but 6 fails) is **completed-with-cleanup-debt**: the release is correct; the workflow file lingers. Operator records partial acceptance + opens a one-line cleanup cycle to remove the file (same pattern as cnos#425's `144db7f1` cleanup commit).

Rejection (4 fails) is **completed-but-rejected**: a remote run happened but the result is not what was asked for. Operator records rejection + invokes the appropriate recovery path from the failure-modes table.

---

## Composition with V/δ and the cycle's full receipt

This remote-runner receipt is **one piece of evidence** inside cnos#426's full cycle receipt. Per `delta/SKILL.md §8.4`:

- V (per `RECEIPT-VALIDATION.md`) validates the cycle's full receipt; the remote-runner receipt is a typed evidence ref V dereferences when the cycle's matter includes a remote-runner-triggering artifact.
- δ-as-role records the `BoundaryDecision` against V's verdict per `delta/SKILL.md §1.5` as normal.
- This receipt does NOT replace V or the BoundaryDecision — it is the surface that makes the *authoring* of the effect artifact legible.
- If post-run evidence is missing (§5 fields unfilled past a reasonable window) and the operator has not declared acceptance, δ records the cycle as degraded with the override block populated and `failed_predicates_overridden: [remote_runner_evidence_missing]`.

---

## Relationship to cnos#425 first use

cnos#425 was the first instantiation of this doctrine (tag retarget). cnos#426 is the second (release publication). The two together exercise the same primitive on two adjacent effect surfaces inside the same release boundary:

- **cnos#425** authored `.github/workflows/repoint-3.82.0.yml` to move the `3.82.0` tag from its prior commit to `fd1d654e`. Run was successful; tag moved correctly. **But** `release.yml` did not re-trigger (cdd-iteration F2 of this cycle), and even if it had, its `build` job is broken on tag-push (cdd-iteration F1). So the release was not auto-published.
- **cnos#426** (this cycle) authors `.github/workflows/publish-3.82.0-release.yml` to publish the GH release directly using the same `softprops/action-gh-release@v1` action `release.yml` already uses, sidestepping the broken `build` job. The pattern is identical: one-shot, push-triggered on own path, self-deleting after work completes; 6-field receipt with operator as acceptor.

Together the two cycles complete the v3.82.0 release boundary as a doctrine-governed sequence. The pattern of "infrastructure breaks; one-shot workaround under doctrine; doctrine'd workaround leaves a finding for proper repair" is exactly the trick-vs-protocol distinction `BOX-AND-THE-RUNNER.md §"What not to celebrate"` names: the second use is not "we keep escaping," it is "we keep recording, and we file the repair-the-base findings (cdd-iteration F1 + F2) as next-MCA so the workaround does not become the steady-state."

---

Filed by γ@cnos (γ+α+β-collapsed-on-δ) on 2026-05-24, authored under δ-as-role authority for the remote-runner artifact per `delta/SKILL.md §8`.
