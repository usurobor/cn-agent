# Remote-runner receipt — v3.82.0 GH release-notes republication (cnos#427)

**Doctrine:** [`docs/gamma/essays/BOX-AND-THE-RUNNER.md`](../../../docs/gamma/essays/BOX-AND-THE-RUNNER.md)
**Skill:** [`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/delta/SKILL.md) §8 (Remote-runner delegation)
**Artifact:** [`.github/workflows/republish-3.82.0-notes.yml`](../../../.github/workflows/republish-3.82.0-notes.yml)
**Cycle:** cnos#427 — Rewrite v3.82.0 release notes per write/SKILL.md (third use of remote-runner-delegation doctrine)
**Date:** 2026-05-24

This receipt is the **third instantiation** of the 6-field remote-runner receipt convention. First use: cnos#425 (tag retarget). Second use: cnos#426 (GH release publication). This third use updates the GH release **body** so the GitHub-facing release notes match the rewritten in-repo canonical body. The doctrine (`BOX-AND-THE-RUNNER.md` essay + `delta/SKILL.md §8`) already lives on main from cnos#425; this cycle is a docs + workflow + receipt cycle that does not modify the doctrine. The receipt is authored in the same commit-set as the workflow it accompanies, per the authoring-order rule in `delta/SKILL.md §8.2`. The `evidence` field (§5) is a post-run-fillable placeholder; the run cannot happen until the cycle merges and the push trigger fires.

---

## 1. Who asked for it

**Operator** (cnos γ session, 2026-05-24), via the following directives in the dispatch chain:

- "rewrite release notes" — operator directing the docs rewrite at the release-notes layer.
- "And" — operator linking the rewrite to a parallel-applied publish step (i.e., not just commit the new body to repo, also push it to the GH release surface).
- Implicit "load write/SKILL.md" — operator's authorial intent is to bring the v3.82.0 release-notes under the `cnos.core/skills/write/SKILL.md` discipline (the prior v3.82.0 RELEASE.md was 5× longer than 3.80.0/3.81.0; mixed multiple jobs per write/SKILL.md §3.2; repeated stop-condition language per §3.3; opened with abstraction per §3.4; contained throat-clearing per §3.5; exceeded reasonable length per §3.14).

The cnos#427 issue body codifies these directives as the cycle's mandate. No agent-initiated escalation; the move is operator-directed and explicitly third-use of the doctrine. The "load write/SKILL.md" intent surfaces as cdd-iteration F1 (`cdd-skill-gap`: release/SKILL.md must require write/SKILL.md load for release-notes authoring) so future cycles do not re-discover this gap.

---

## 2. What it will run

The workflow's four job steps, in order:

```text
1. actions/checkout@v4 with:
     ref: 'main'             (checkout at main HEAD — the rewritten
                              tight body lives on main as root
                              RELEASE.md after this cycle merges; the
                              tag's old RELEASE.md at fd1d654e is NOT
                              the source for this cycle)
     fetch-depth: 0          (full history so the self-delete step
                              can push to main)
     token: GITHUB_TOKEN     (substrate-provided; contents: write)

2. Verify root RELEASE.md is the new tight version:
     head -1 RELEASE.md
     grep -q "^# 3.82.0$" RELEASE.md || (echo ERROR; exit 1)
     LINES=$(wc -l < RELEASE.md)
     [ "$LINES" -le 90 ] || (echo ERROR: ...; exit 1)
   Guard against publishing the old verbose body (which started
   "# v3.82.0 — CCNF package-architecture baseline") or any body
   exceeding the 90-line budget.

3. softprops/action-gh-release@v1 with:
     tag_name: '3.82.0'
     name: '3.82.0'
     body_path: RELEASE.md   (sourced from the main-checked-out tree)
     draft: false
     prerelease: false
   Updates the existing GH release at refs/tags/3.82.0 in place
   (softprops/action-gh-release@v1 is idempotent — body and name are
   overwritten with the values supplied). The tag SHA (fd1d654e) is
   NOT touched; binaries (none, per cnos#426 cdd-iteration F1) are
   NOT touched.

4. Self-delete the workflow file:
     git config user.name  "cnos-remote-runner[bot]"
     git config user.email "cnos-remote-runner@cnos.local"
     git fetch origin main
     git checkout main
     git rm .github/workflows/republish-3.82.0-notes.yml
     git commit -m "release: remove one-shot republish-3.82.0-notes workflow (cnos#427)"
     git push origin main
```

The body-update step (step 3) is the load-bearing effect: it transitions the GitHub Release body from the prior verbose 109-line body (published by cnos#426's workflow) to the rewritten ~45-line tight body. Because the checkout in step 1 pins `ref: 'main'`, the `RELEASE.md` body_path resolves to the post-merge main tree, which carries the new tight version.

The self-delete step (step 4) closes the latent-execution authority of the workflow itself, per `delta/SKILL.md §8.3` (one-shot workflows self-delete). cnos#426's workflow needed a worktree switch because its checkout was at the tag; this workflow's checkout is at main already, so the `git fetch origin main` + `git checkout main` are no-op equivalents at the worktree level (they ensure the current ref is current main rather than a slightly-stale post-checkout state). The commit + push proceed normally.

---

## 3. Where it will run

**GitHub Actions runner: `ubuntu-latest`** (GitHub-hosted ephemeral VM).

Not a self-hosted runner. Not a third-party CI provider. The runner is a single-use VM provisioned by GitHub for this workflow's job; it is destroyed after the job completes. Same execution surface as cnos#425's first-use and cnos#426's second-use workflows — three consecutive uses on the same runner class is evidence-class identical (the doctrine treats `ubuntu-latest` as the default remote-runner surface for cnos one-shots).

---

## 4. What authority it has

**`GITHUB_TOKEN`** with **`contents: write`** permission, scoped to `usurobor/cnos`.

No other secrets are accessed by this workflow. No deploy keys. No cloud OIDC role assumptions. No third-party tokens. Same authority shape as cnos#425 and cnos#426. The workflow does not invoke `gh` CLI commands that require additional scopes; both the release-body-update step (via `softprops/action-gh-release@v1`, which uses `GITHUB_TOKEN` automatically) and the self-delete commit-and-push succeed under `contents: write` alone.

The `contents: write` scope grants:
- update GitHub Releases body/name (used by step 3)
- push commits to `main` (used by step 4 self-delete)
- read repository contents at `main` (used by step 1 checkout)

The scope does NOT grant:
- write access to other repositories under `usurobor/`
- access to repository secrets beyond `GITHUB_TOKEN`
- workflow-editing capability beyond this repo (the self-delete commit edits only this workflow's own file)
- tag modification (this workflow does NOT touch the tag — the tag SHA fd1d654e set by cnos#425 is preserved)
- release-binary attach/delete (this workflow does NOT touch binaries — softprops/action-gh-release@v1 in update mode preserves existing assets)

Branch-protection on `main` is enforced by the substrate; the workflow's self-delete commit is subject to whatever rules apply to direct pushes to `main` from `GITHUB_TOKEN`-authenticated commits. If branch-protection rejects the self-delete push, the release body update from step 3 still succeeds and the workflow file persists on `main` until manually removed; this is a named failure mode (§Failure modes below) and is the same pattern resolved manually in cnos#425's first use (commit `144db7f1` — `release: remove one-shot repoint-3.82.0 workflow (cnos#425)`).

---

## 5. Evidence

**Post-run-fillable placeholder.** The receipt is authored before the run; the evidence shape is named in advance so a downstream consumer can detect a missing receipt rather than reading silence as success.

Required evidence (to be filled after the run):

| Field | Shape | Value (post-run) |
|---|---|---|
| Workflow run URL | `https://github.com/usurobor/cnos/actions/runs/<run-id>` | _(post-run)_ |
| Updated release URL | `https://github.com/usurobor/cnos/releases/tag/3.82.0` | _(post-run)_ |
| Release body starts with new header | First line of release body equals `# 3.82.0` | _(post-run)_ |
| Release body line count ≤ 90 | `gh release view 3.82.0 --json body --jq .body \| wc -l` ≤ 90 | _(post-run)_ |
| Self-delete commit SHA | SHA of the `release: remove one-shot republish-3.82.0-notes workflow (cnos#427)` commit on `main` | _(post-run)_ |
| Workflow file gone | `.github/workflows/republish-3.82.0-notes.yml` no longer present on `main` | _(post-run)_ |

The acceptor (§6) fills these fields when declaring acceptance. The receipt then moves with the cycle into `.cdd/releases/3.82.0/` (or the next release directory) per standard cycle-close artifact handling.

---

## 6. Who may accept the result

**Operator.** Acceptance criterion:

> Operator inspects `https://github.com/usurobor/cnos/releases/tag/3.82.0` and confirms that (a) the release body's first line is `# 3.82.0` (the new tight version), (b) the release body line count is ≤ 90, and (c) the body content matches root `RELEASE.md` on main (which equals `.cdd/releases/3.82.0/RELEASE.md`). The tag SHA must remain `fd1d654e` (this cycle does NOT touch the tag).

The operator's δ-decision is the move from "the run completed" to "the run is accepted." A run that completes but produces a release body that does not match root `RELEASE.md` (e.g. because main drifted between merge and run, or the verify step's grep was bypassed) is a completed-but-rejected outcome; the operator records the rejection and invokes the manual recovery path (§Failure modes).

No agent — α, β, γ, δ-as-role, or ε — has acceptance authority for this remote-runner move. The acceptance is at the project-binding / release-boundary layer (`ROLES.md` §5-layer chain) and is held by the operator at the outermost cell (per `CELL-OF-CELLS.md §13`). Inherited unchanged from cnos#425 and cnos#426 receipts §6.

---

## Expected effect (5 steps, in order)

1. **Merge of cycle/427 to main.** Operator merges the cycle's branch; the merge commit lands root `RELEASE.md` (rewritten), `.cdd/releases/3.82.0/RELEASE.md` (mirrored), `CHANGELOG.md` (3.82.0 row tightened), `.github/workflows/republish-3.82.0-notes.yml` (new), and `.cdd/unreleased/427/*` (this receipt + close-outs) onto `main`.
2. **Push trigger fires.** The push to `main` includes a change to `.github/workflows/republish-3.82.0-notes.yml` (because the file is newly added); the path filter matches; the workflow's `republish` job is scheduled.
3. **Checkout at main + body verification.** Step 1 checks out main HEAD. Step 2 verifies the first line of root `RELEASE.md` is `# 3.82.0` and the file is ≤ 90 lines; exits non-zero on mismatch.
4. **GH release body updated.** Step 3 invokes `softprops/action-gh-release@v1` in update mode with `tag_name: '3.82.0'`, `name: '3.82.0'`, `body_path: RELEASE.md`. The release at `https://github.com/usurobor/cnos/releases/tag/3.82.0` has its body and name overwritten with the new tight content. Tag SHA and binary assets (if any) untouched.
5. **Self-delete.** Step 4 runs `git rm` + `git commit` + `git push origin main`, removing `.github/workflows/republish-3.82.0-notes.yml` from `main`. The latent-execution authority of this workflow closes. The self-delete commit on `main` is the structural signal that the one-shot ran.

---

## Failure modes + mitigations

| Failure | Detection | Mitigation |
|---|---|---|
| Branch-protection on `main` rejects the workflow's self-delete push | Workflow logs show `push origin main` failure in step 4 | Operator manually removes `.github/workflows/republish-3.82.0-notes.yml` via PR (a normal cycle, not under this receipt). Release body update from step 3 already succeeded, so the release is correct; only cleanup remains. Same pattern as cnos#425's `144db7f1` cleanup commit. |
| RELEASE.md verification fails (wrong version header or line-count overrun) | Workflow logs show step 2 exits non-zero with `ERROR: RELEASE.md not the new version` or `ERROR: RELEASE.md is N lines (>90)` | Workflow fails fast before publishing a wrong body. Operator investigates whether main has the rewritten body (`head -1 RELEASE.md` should match `# 3.82.0`); if main drifted, file a fix-cycle; if body is correct but line-count overran, re-tighten the body to ≤ 90 lines. |
| softprops/action-gh-release@v1 fails (action unavailable, network, rate limit) | Workflow logs show step 3 failure with action-side error | Operator re-runs the workflow via `gh run rerun <run-id>` (the file is still on main until step 4 completes, so the trigger can be re-fired by an empty-commit touch to the workflow path — or by direct re-run from the Actions UI). If the action is permanently unavailable, fall back to `gh release edit 3.82.0 --notes-file RELEASE.md` invoked from operator's local under operator credentials. |
| GITHUB_TOKEN lacks `contents: write` at runtime (org policy override) | Workflow logs show 403 on release update or self-delete push | Operator escalates to GitHub org settings; if temporarily unavailable, falls back to manual `gh release edit` under operator credentials. |
| Release does not exist at tag 3.82.0 | Step 3 reports release-not-found | Should not occur (cnos#426 published the release). If it does, the action's default behavior is to CREATE the release rather than update — this is the desired fallback; verify in step 3 logs that the release was created (not updated) and that `name` and `body` are correct post-run. |
| Workflow does not trigger (path filter mismatch, file permissions) | Operator merges the cycle but no workflow run appears in Actions tab | Operator inspects the path filter against the file path; if mismatch is found, file a follow-up cycle to fix the trigger or invoke manually via `workflow_dispatch` (this workflow has no `workflow_dispatch` trigger; would require an edit-cycle to add one). Body update can also be done manually per the operator-credential fallback above. |
| Tag accidentally moved | Operator inspects `git ls-remote origin refs/tags/3.82.0` and finds a SHA other than `fd1d654e` | Should not occur (this workflow does NOT touch the tag — softprops/action-gh-release@v1 in update mode does not alter the tag SHA). If detected, run cnos#425's retarget pattern to restore. |

---

## Acceptance criteria

This remote-runner move is accepted when **all** of the following hold:

1. The workflow run for `.github/workflows/republish-3.82.0-notes.yml` completed successfully (all 4 steps green).
2. The GitHub Release at `https://github.com/usurobor/cnos/releases/tag/3.82.0` exists and is not draft / not prerelease.
3. The release name is `3.82.0` (tight; the prior `3.82.0 — CCNF package-architecture baseline` suffix is gone, matching the new RELEASE.md's `# 3.82.0` header per write/SKILL.md §3.4).
4. The release body's first line is `# 3.82.0` and the body line count is ≤ 90.
5. The release body content matches root `RELEASE.md` on main (which equals `.cdd/releases/3.82.0/RELEASE.md`).
6. `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e8d6361f0db2a6407f19e573a96d1054d` (unchanged — this cycle does not touch the tag).
7. `.github/workflows/republish-3.82.0-notes.yml` is no longer present on `main` (self-delete succeeded).
8. Operator records acceptance by filling the §5 evidence table in this receipt and moving the receipt into the appropriate release directory.

Partial acceptance (1, 2, 3, 4, 5, 6 pass but 7 fails) is **completed-with-cleanup-debt**: the release is correct; the workflow file lingers. Operator records partial acceptance + opens a one-line cleanup cycle to remove the file (same pattern as cnos#425's `144db7f1` cleanup commit).

Rejection (4 or 5 fails) is **completed-but-rejected**: a remote run happened but the result is not what was asked for. Operator records rejection + invokes the appropriate recovery path from the failure-modes table.

---

## Composition with V/δ and the cycle's full receipt

This remote-runner receipt is **one piece of evidence** inside cnos#427's full cycle receipt. Per `delta/SKILL.md §8.4`:

- V (per `RECEIPT-VALIDATION.md`) validates the cycle's full receipt; the remote-runner receipt is a typed evidence ref V dereferences when the cycle's matter includes a remote-runner-triggering artifact.
- δ-as-role records the `BoundaryDecision` against V's verdict per `delta/SKILL.md §1.5` as normal.
- This receipt does NOT replace V or the BoundaryDecision — it is the surface that makes the *authoring* of the effect artifact legible.
- If post-run evidence is missing (§5 fields unfilled past a reasonable window) and the operator has not declared acceptance, δ records the cycle as degraded with the override block populated and `failed_predicates_overridden: [remote_runner_evidence_missing]`.

---

## Relationship to cnos#425 + cnos#426

cnos#425 was the first instantiation of this doctrine (tag retarget). cnos#426 was the second (release publication). cnos#427 (this cycle) is the third — and notably the first use where the *target effect surface* is the release-notes body rather than the tag or the release-existence. The three together complete the v3.82.0 release boundary as a doctrine-governed sequence on three adjacent effect surfaces:

- **cnos#425** — moved the `3.82.0` tag from its prior commit to `fd1d654e`. Tag SHA effect.
- **cnos#426** — published the GH release at `3.82.0` with the v3.82.0 baseline body. Release-existence effect.
- **cnos#427** (this cycle) — updates the GH release body to the rewritten tight version. Release-body effect.

The pattern is identical across all three: one-shot, push-triggered on own path, self-deleting after work completes; 6-field receipt with operator as acceptor. The third use differs from the second use only in:

1. **Checkout ref.** cnos#426 checked out at the tag (because main HEAD had post-tag content that would publish the wrong body). cnos#427 checks out at main (because this cycle's whole purpose is to publish the *current* main RELEASE.md, which is the rewritten tight version landed by this same cycle).
2. **Target action mode.** cnos#426 effectively created-or-updated the release. cnos#427 updates it in place (the release exists from cnos#426). softprops/action-gh-release@v1 idempotency handles both modes transparently.
3. **Self-delete worktree.** cnos#426 needed a worktree switch (checkout was at tag). cnos#427's checkout is at main, so the worktree switch is no-op-equivalent (kept in the recipe for symmetry; safe to execute).

The pattern of "infrastructure works; doctrine governs three adjacent uses on the same release boundary; each use leaves a 6-field receipt" is exactly what `BOX-AND-THE-RUNNER.md §"What this enables"` named: the primitive is reusable, and three uses on adjacent surfaces of the same release boundary is the first evidence that it is durable as a primitive rather than a one-off escape.

cdd-iteration F1 of this cycle (release/SKILL.md must require write/SKILL.md load for release-notes authoring) addresses the *upstream* gap: cnos#427 had to be authored because cnos#422's release-notes authoring did not surface write/SKILL.md as a required load. Recording this as a `cdd-skill-gap` `next-MCA` finding ensures the gap is closed at the skill layer rather than re-discovered at the next release.

---

Filed by γ@cnos (γ+α+β-collapsed-on-δ) on 2026-05-24, authored under δ-as-role authority for the remote-runner artifact per `delta/SKILL.md §8`.
