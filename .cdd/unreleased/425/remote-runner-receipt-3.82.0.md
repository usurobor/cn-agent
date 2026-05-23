# Remote-runner receipt — v3.82.0 tag retarget (cnos#425)

**Doctrine:** [`docs/gamma/essays/BOX-AND-THE-RUNNER.md`](../../../docs/gamma/essays/BOX-AND-THE-RUNNER.md)
**Skill:** [`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/delta/SKILL.md) §8 (Remote-runner delegation)
**Artifact:** [`.github/workflows/repoint-3.82.0.yml`](../../../.github/workflows/repoint-3.82.0.yml)
**Cycle:** cnos#425 — Capture remote-runner-delegation primitive + first use (v3.82.0 tag retarget via GH Actions)
**Date:** 2026-05-23

This receipt is the first instantiation of the 6-field remote-runner receipt convention. It is authored in the same commit-set as the workflow it accompanies, per the authoring-order rule in `delta/SKILL.md §8.2`. The receipt's `evidence` field is a post-run-fillable placeholder; the run cannot happen until the cycle merges and the push trigger fires.

---

## 1. Who asked for it

**Operator** (cnos γ session, 2026-05-23), via the following directives in the dispatch chain:

- "execute as delta or gamma" — operator authorizing the agent to act under δ-as-role authority for the release-fix boundary
- "you can create gh action and then run it" — operator authorizing the remote-runner move specifically

The cnos#425 issue body codifies these directives as the cycle's mandate. No agent-initiated escalation; the move is operator-directed.

---

## 2. What it will run

The workflow's three job steps, in order:

```text
1. actions/checkout@v4 with fetch-depth: 0 (clone full history so the
   tag move can resolve fd1d654e)

2. Configure git identity:
     git config user.name  "cnos-remote-runner[bot]"
     git config user.email "cnos-remote-runner@cnos.local"

3. Force-move tag:
     git tag -f 3.82.0 fd1d654e
     git push --force origin 3.82.0

4. Self-delete:
     git rm .github/workflows/repoint-3.82.0.yml
     git commit -m "release: remove one-shot repoint-3.82.0 workflow (cnos#425)"
     git push origin main
```

The tag-push step (step 3) is the load-bearing effect: it transitions the `3.82.0` tag from its prior commit to `fd1d654e`. The tag push triggers `release.yml` (which is path-triggered on `tags: [0-9]*.[0-9]*.[0-9]*`); release.yml then publishes the GH release with the body sourced from root `RELEASE.md` at `fd1d654e` (which is the canonical v3.82.0 release notes).

The self-delete step (step 4) closes the latent-execution authority of the workflow itself, per `delta/SKILL.md §8.3` (one-shot workflows self-delete).

---

## 3. Where it will run

**GitHub Actions runner: `ubuntu-latest`** (GitHub-hosted ephemeral VM).

Not a self-hosted runner. Not a third-party CI provider. The runner is a single-use VM provisioned by GitHub for this workflow's job; it is destroyed after the job completes.

---

## 4. What authority it has

**`GITHUB_TOKEN`** with **`contents: write`** permission, scoped to `usurobor/cnos`.

No other secrets are accessed by this workflow. No deploy keys. No cloud OIDC role assumptions. No third-party tokens. The workflow does not invoke `gh release` or any commands that require additional scopes; the tag push and the self-delete commit both succeed under `contents: write` alone.

The `contents: write` scope grants:
- push to refs (including `--force` — required for the tag move)
- create/delete tags
- push commits to `main` (required for the self-delete)

The scope does NOT grant:
- write access to other repositories under `usurobor/`
- access to repository secrets beyond `GITHUB_TOKEN`
- workflow-editing capability beyond this repo (the self-delete commit edits only this workflow's own file)

Branch-protection on `main` is enforced by the substrate; the workflow's self-delete commit is subject to whatever rules apply to direct pushes to `main` from `GITHUB_TOKEN`-authenticated commits. If branch-protection rejects the self-delete push, the workflow logs the failure and the workflow file persists on `main` until manually removed; this is a named failure mode (§Failure modes below).

---

## 5. Evidence

**Post-run-fillable placeholder.** The receipt is authored before the run; the evidence shape is named in advance so a downstream consumer can detect a missing receipt rather than reading silence as success.

Required evidence (to be filled after the run):

| Field | Shape | Value (post-run) |
|---|---|---|
| Workflow run URL | `https://github.com/usurobor/cnos/actions/runs/<run-id>` | _(post-run)_ |
| Tag move proof | `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e` | _(post-run)_ |
| Release re-publish URL | `https://github.com/usurobor/cnos/releases/tag/3.82.0` | _(post-run)_ |
| Self-delete commit SHA | SHA of the `release: remove one-shot repoint-3.82.0 workflow (cnos#425)` commit on `main` | _(post-run)_ |
| Workflow file gone | `.github/workflows/repoint-3.82.0.yml` no longer present on `main` | _(post-run)_ |

The acceptor (§6) fills these fields when declaring acceptance. The receipt then moves with the cycle into `.cdd/releases/3.82.0/` (or the next release directory) per standard cycle-close artifact handling.

---

## 6. Who may accept the result

**Operator.** Acceptance criterion:

> Operator inspects `https://github.com/usurobor/cnos/releases/tag/3.82.0` and confirms that the release body matches `.cdd/releases/3.82.0/RELEASE.md` (the canonical v3.82.0 release notes — the "CCNF package-architecture baseline" body, 110 lines).

The operator's δ-decision is the move from "the run completed" to "the run is accepted." A run that completes but produces a release body that does not match `.cdd/releases/3.82.0/RELEASE.md` (e.g. because `release.yml` sourced `RELEASE.md` from a commit other than `fd1d654e`) is a completed-but-rejected outcome; the operator records the rejection and invokes the manual recovery path (§Failure modes).

No agent — α, β, γ, δ-as-role, or ε — has acceptance authority for this remote-runner move. The acceptance is at the project-binding / release-boundary layer (`ROLES.md` §5-layer chain) and is held by the operator at the outermost cell (per `CELL-OF-CELLS.md §13`).

---

## Expected effect (5 steps, in order)

1. **Merge of cycle/425 to main.** Operator merges the cycle's branch; the merge commit lands `.github/workflows/repoint-3.82.0.yml` (plus the essay, the skill section, this receipt, and the close-out artifacts) onto `main`.
2. **Push trigger fires.** The push to `main` includes a change to `.github/workflows/repoint-3.82.0.yml` (because the file is newly added); the path filter matches; the workflow's `repoint` job is scheduled.
3. **Tag force-move.** Step 3 of the job runs `git tag -f 3.82.0 fd1d654e` followed by `git push --force origin 3.82.0`. The remote `3.82.0` tag now points at `fd1d654e`.
4. **release.yml triggers on the new tag.** The tag push (under the `[0-9]*.[0-9]*.[0-9]*` pattern in `.github/workflows/release.yml`) starts the `build` → `publish` → `smoke` → `notify` chain. The `publish` job sources `RELEASE.md` from the repository at `fd1d654e` (where root `RELEASE.md` is the canonical v3.82.0 body) and creates/updates the GH release.
5. **Self-delete.** Step 4 of the repoint job runs `git rm` + `git commit` + `git push origin main`, removing `.github/workflows/repoint-3.82.0.yml` from `main`. The latent-execution authority of this workflow closes. The self-delete commit on `main` is the structural signal that the one-shot ran.

---

## Failure modes + mitigations

| Failure | Detection | Mitigation |
|---|---|---|
| Branch-protection on `main` rejects the workflow's self-delete push | Workflow logs show `push origin main` failure in step 4 | Operator manually removes `.github/workflows/repoint-3.82.0.yml` via PR (a normal cycle, not under this receipt). Tag move from step 3 already succeeded, so the release is correct; only cleanup remains. |
| Tag force-push rejected by remote (branch/tag protection) | Workflow logs show `push --force origin 3.82.0` failure in step 3 | Operator falls back to manual tag move via release-effector runbook (`release-effector/SKILL.md`): locally `git tag -f 3.82.0 fd1d654e && git push --force origin 3.82.0` under operator's own credentials. Then operator manually triggers release.yml via `workflow_dispatch` if needed. |
| `release.yml` triggers but sources `RELEASE.md` from wrong commit | Operator inspects release body and finds it does not match `.cdd/releases/3.82.0/RELEASE.md` | Operator declares the run completed-but-rejected; investigates which commit `release.yml` checked out (`ref: main` in the publish job — if main has drifted further, the body may be from a later commit). Mitigation: re-run with `workflow_dispatch` after pinning `ref: fd1d654e` in a follow-up PR, or edit the GH release body manually. |
| Workflow does not trigger (path filter mismatch, file permissions) | Operator merges the cycle but no workflow run appears in Actions tab | Operator inspects the path filter against the file path; if mismatch is found, file a follow-up cycle to fix the trigger or invoke manually via `workflow_dispatch`. The tag move can also be done manually per the release-effector runbook. |
| GITHUB_TOKEN lacks contents: write at runtime (org policy override) | Workflow logs show 403 on tag push or commit push | Operator escalates to GitHub org settings; if temporarily unavailable, falls back to manual tag move under operator credentials. |
| release.yml infrastructure red on retriggered run (CI flake, runner issue) | release.yml job fails post-tag-move | This is a known δ-side recovery scenario per `delta/SKILL.md §1.1` (CI Red → δ owns the failure). Operator investigates release logs; if pre-existing infra failure, override per §3 of delta/SKILL.md with `failed_predicates_overridden` populated. |

---

## Acceptance criteria

This remote-runner move is accepted when **all** of the following hold:

1. The workflow run for `.github/workflows/repoint-3.82.0.yml` completed successfully (all 4 steps green).
2. `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e8d6361f0db2a6407f19e573a96d1054d`.
3. `release.yml` produced a workflow run on the retriggered `3.82.0` tag, and the `publish` job completed (release exists with shipped binaries + checksums).
4. The GH release body at `https://github.com/usurobor/cnos/releases/tag/3.82.0` matches `.cdd/releases/3.82.0/RELEASE.md` (108 lines starting "# v3.82.0 — CCNF package-architecture baseline").
5. `.github/workflows/repoint-3.82.0.yml` is no longer present on `main` (self-delete succeeded).
6. Operator records acceptance by filling the §5 evidence table in this receipt and moving the receipt into the appropriate release directory.

Partial acceptance (1, 2, 3, 4 pass but 5 fails) is **completed-with-cleanup-debt**: the release is correct; the workflow file lingers. Operator records partial acceptance + opens a one-line cleanup cycle to remove the file.

Rejection (4 fails) is **completed-but-rejected**: a remote run happened but the result is not what was asked for. Operator records rejection + invokes the appropriate recovery path from the failure-modes table.

---

## Composition with V/δ and the cycle's full receipt

This remote-runner receipt is **one piece of evidence** inside cnos#425's full cycle receipt. Per `delta/SKILL.md §8.4`:

- V (per `RECEIPT-VALIDATION.md`) validates the cycle's full receipt; the remote-runner receipt is a typed evidence ref V dereferences when the cycle's matter includes a remote-runner-triggering artifact.
- δ-as-role records the `BoundaryDecision` against V's verdict per `delta/SKILL.md §1.5` as normal.
- This receipt does NOT replace V or the BoundaryDecision — it is the surface that makes the *authoring* of the effect artifact legible.
- If post-run evidence is missing (§5 fields unfilled past a reasonable window) and the operator has not declared acceptance, δ records the cycle as degraded with the override block populated and `failed_predicates_overridden: [remote_runner_evidence_missing]`.

---

Filed by γ@cnos (γ+α+β-collapsed-on-δ) on 2026-05-23, authored under δ-as-role authority for the remote-runner artifact per `delta/SKILL.md §8`.
