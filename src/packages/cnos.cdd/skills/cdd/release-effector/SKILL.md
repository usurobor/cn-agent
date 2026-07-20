---
name: release-effector
description: δ-side mechanics for cutting a release — scripts/release.sh invocation, tag push, release CI polling, CI-red recovery runbook, branch cleanup. The platform-actions companion to release/SKILL.md (β-side authoring).
artifact_class: skill
kata_surface: embedded
governing_question: When γ declares closure, what does δ run, in what order, and how does δ verify the release disconnected cleanly?
visibility: internal
parent: cdd
triggers:
  - release effector
  - scripts/release.sh
  - tag push
  - release CI
  - branch cleanup
  - disconnect release
scope: task-local
inputs:
  - gamma-closeout.md on main
  - VERSION decision (per release/SKILL.md §2.2)
  - RELEASE.md on main (per release/SKILL.md §2.5)
  - complete cycle records under .cdd/unreleased/{N}/ for release validation
outputs:
  - bare X.Y.Z annotated tag on origin
  - release CI green (or operator-accepted override)
  - merged cycle/{N} branches deleted (or attempt recorded if harness blocks delete)
  - gate-action completion confirmations per operator/SKILL.md §3.3
requires:
  - δ has read operator/SKILL.md §3 (outward gate policy)
  - γ has shipped RELEASE.md; cycle directories remain under unreleased until after disconnect
  - "gamma-closeout.md on main carries `CDD-Post-Merge-Closeout: complete`"
  - β has signaled "release ready for δ tag" in beta-closeout.md (per release/SKILL.md §2.6)
calls:
  - release/SKILL.md
  - operator/SKILL.md
---

# Release Effector (δ)

## Core Principle

**The tag is the disconnection point. δ runs the script, polls CI, cleans the branches.**

`release-effector` is the **mechanics-of-the-disconnect** skill. β-side authoring (RELEASE.md, CHANGELOG, version decision, cycle-dir move) lives in `release/SKILL.md`. δ-policy frame (when the release gate fires, what δ does NOT do, override authority) lives in `operator/SKILL.md` (Phase 4a will relocate to `delta/SKILL.md`). This skill is the executable surface between them: the exact commands δ runs at the disconnect moment, in order, with the failure-mode runbook for red CI.

Failure mode: **manual tagging.** A `git tag X.Y.Z && git push --tags` skips version stamping, skips consistency checks, skips structured tag-message generation, and produces a release whose VERSION, `cn.json`, package manifests, and tag disagree. The single-command release script is the only sanctioned path.

---

## Preconditions

Before δ runs `scripts/release.sh`, all of the following must be true on `main`:

| Precondition | Owner | Verified where |
|---|---|---|
| `gamma-closeout.md` on main carries `CDD-Post-Merge-Closeout: complete` | γ | `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Ownership matrix"; `gamma/SKILL.md` §2.10 closeout gate |
| `RELEASE.md` exists at repo root | γ | `release/SKILL.md` §2.5; `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Ownership matrix" |
| Cycle directories remain at `.cdd/unreleased/{N}/` for validation | γ | `release/SKILL.md` §2.5a |
| VERSION decided (semver bump) | γ + β | `release/SKILL.md` §2.2 |
| β signaled "release ready for δ tag" in `beta-closeout.md` | β | `release/SKILL.md` §2.6 |
| Local main up to date with origin/main | δ | `scripts/release.sh` step 2 (auto-check) |
| Tag X.Y.Z does not already exist | δ | `scripts/release.sh` step 3 (auto-check) |

The script enforces the last two mechanically and errors out if they fail. The first five are γ/β responsibilities; δ confirms they are met before invoking the script.

---

## 1. Run scripts/release.sh

**Single-command release.** Edit `VERSION`, then run the script:

```bash
# Edit VERSION first (canonical: bare X.Y.Z, no v-prefix), then:
scripts/release.sh

# Or pass version directly:
scripts/release.sh 3.67.0
```

**Manual tagging is not allowed.** Do not run `git tag` directly. The release script is the only way to tag. It prevents the class of failures where VERSION, `cn.json`, and package manifests disagree (see `DISPATCH-FAILURE-EVIDENCE.md`, cycle #84 failure 3).

**What the script does, in order** (mechanical narration; matches `scripts/release.sh` as of cycle #399):

1. Set VERSION (writes `$1` to `VERSION` if a version arg was passed).
2. Confirm current branch is `main` and `HEAD` matches `origin/main`. Errors out otherwise.
3. Confirm the tag does not already exist locally. Errors out if it does (use `release/SKILL.md` §3.6 amend-don't-re-tag procedure when CI fails post-tag).
4. Run `scripts/validate-release-gate.sh` — RELEASE.md + cycle-artifact completeness check (see `release/SKILL.md` §2.1 for the validator surface).
5. Run `scripts/stamp-versions.sh` — derives `cn.json` and package manifests from VERSION.
6. Run `scripts/check-version-consistency.sh` — validates all version-stamped files agree.
7. Stage + commit `release: $VERSION` if anything changed in the working tree. Cycle directories are deliberately not moved.
8. Warn-prompt (interactive `y/N`) if `CHANGELOG.md` does not contain a `## $VERSION` heading. Abort path is `N`.
9. Generate the annotated tag message via `scripts/generate-release-tag-message.sh`.
10. Create the **annotated** tag at the bare version.
11. `git push origin main --tags`.

After the script returns, monitor release CI per §3 below.

- ❌ `git tag 3.67.0 && git push --tags` (skips stamp, skips consistency check, skips structured tag message)
- ❌ γ's skill patches sit on main untagged across multiple cycles
- ❌ δ defers release "because there are no consumers" (the tag is structural, not consumer-driven)
- ✅ `scripts/release.sh 3.67.0` — stamps, verifies, commits, tags, pushes

---

## 2. Tag policy (bare X.Y.Z, no v-prefix)

**All tags are bare `X.Y.Z`.** Canonical reference: `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix" — *"Tags are bare `X.Y.Z` everywhere (VERSION file, git tag, branch-name version segment, CHANGELOG row, RELEASE.md, snapshot directory). `v`-prefixed tags are legacy and warn-only."*

The bare-version convention applies to every surface in the chain:

| Surface | Form |
|---|---|
| `VERSION` file content | `3.67.0` |
| Git tag | `3.67.0` (annotated, not lightweight) |
| Branch name version segment | `cycle/399`, `claude/3.67.0-22-...` (release directory `3.67.0/`) |
| `CHANGELOG.md` ledger row | `## 3.67.0` |
| `RELEASE.md` body | references the bare version |
| Snapshot dir | `docs/gamma/cdd/3.67.0/` |
| Cycle-dir move target | `.cdd/releases/3.67.0/{N}/` |

**Annotated tags only.** `scripts/release.sh` creates `git tag -a $VERSION -F <message-file>` — never lightweight tags. The annotated form carries the structured message generated by `scripts/generate-release-tag-message.sh` (issue metadata, wave context, CDD review artifacts, deterministic fallback). δ can inspect the tag-message body after the fact with:

```bash
git show <version>
git for-each-ref refs/tags/<version> --format='%(contents)'
```

**One tag per release.** Per `release/SKILL.md` §3.6: if CI fails post-tag, amend the release commit, delete-and-recreate the same bare-version tag, force-push. Do NOT proliferate `3.9.1`, `3.9.1-fix`, `3.9.1-fix2`.

- ❌ `v3.67.0` (legacy v-prefix; warn-only per `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix")
- ❌ Lightweight tag (`git tag $VERSION` without `-a` or `-F`) — loses the structured message
- ❌ Tag proliferation across CI-red retries
- ✅ One bare `3.67.0` annotated tag, amended until CI green

---

## 3. Poll release CI

After tag push, the release CI workflow builds binaries (linux-x64, macos-x64, macos-arm64) and attaches them to the GitHub release. δ blocks release completion until that workflow reports green.

```bash
# Quick view:
gh run list --branch <version>

# Or watch a specific run to completion:
gh run watch <run-id> --exit-status

# Or via the workflow surface:
gh run list --workflow release.yml --limit 1
```

**Decision rule:**

- **CI Green** → δ declares the release complete; report completion to γ per `operator/SKILL.md` §3.3.
- **CI Red** → execute the recovery runbook in §4 below. δ MUST NOT declare release complete while CI is red.

**The gate does not close until CI is green or operator explicitly accepts a known pre-existing failure (override per §4 step 5).**

- ❌ "Tag pushed; we're done" — release CI may still be running or red
- ❌ Deploy from binaries before CI completes (stale binary from previous release)
- ✅ `gh run watch <id> --exit-status` then verify assets attached, then report completion

---

## 4. CI-red recovery runbook

When release CI reports red, δ owns the recovery. The runbook is mandatory; declaring the release complete while CI is red is a discipline break (see §6 below).

1. **Investigate** the release CI logs to identify which workflow step failed and why.
2. **Classify** the failure:
   - **Release-specific** — the failure is caused by something this release introduced (a broken stamp, a manifest drift, a binary-build error new to this version).
   - **Pre-existing infrastructure** — the failure was present before this release (flaky smoke test, transient registry outage, expired token, environment regression unrelated to release content).
3. **Fix or escalate** based on classification:
   - **Fixable / release-specific** — fix the cause, then amend the release commit and re-tag per `release/SKILL.md` §3.6 (`git commit --amend --no-edit && git tag -d $VERSION && git tag -a $VERSION -F <message-file> && git push --force origin main $VERSION`). Then re-poll CI. Do NOT proliferate tags.
   - **Pre-existing / infrastructural** — document the failure, escalate to operator-as-human, and do NOT declare release complete on the agent's authority.
4. **Re-verify** — after any fix attempt, re-poll CI (§3 above) and reapply this runbook if CI is still red.
5. **Operator override** — explicit operator-as-human acceptance is the only path to declare release complete with red CI. This is the escape hatch for known pre-existing failures (canonical examples: v3.66.0 and v3.67.0 smoke failures). The override must be declared per `operator/SKILL.md` §4 (override protocol) and recorded in `gamma-closeout.md` or `cdd-iteration.md` so the next release sees the precedent.

**The gate does not close until CI is green or operator explicitly accepts the failure.**

- ❌ Declare release complete because the tag is on origin, ignoring red CI
- ❌ Proliferate `3.67.0-fix`, `3.67.0-fix2` instead of amend-and-re-tag
- ❌ Apply operator override to a release-specific failure ("just ship it")
- ✅ Investigate → classify → fix-or-escalate → re-verify → operator-override-if-pre-existing

---

## 5. Branch cleanup

After the release is declared complete (CI green or accepted override), δ deletes merged cycle branches. This is mechanical cleanup, not judgment — if it's merged into main, it's dead.

**Targets:**
- Cycle branches: `cycle/{N}` (post-#287 canonical form per `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch rule")
- γ session branches: harness-given `claude/...` or operator-named `gamma/session-{N}` patterns
- Any other branches merged into main during this release

No orphan γ session branches survive past closure.

**Procedure** (mirrors `release/SKILL.md` §2.6a):

```bash
git branch -r --merged origin/main \
  | grep -v main \
  | grep -v HEAD \
  | sed 's/origin\///' \
  | xargs -I{} git push origin --delete {}
```

For targeted deletion of a single cycle branch:

```bash
git push origin --delete cycle/399
```

**Harness 403 / push-restriction case.** Some Claude Code harness environments block remote branch deletion (`403` on `git push origin --delete <branch>`; same constraint as `operator/SKILL.md` §5.2 consequence 3). δ attempts the delete; if the harness blocks it, δ records the failure in the cycle's close-out and proceeds — the merged-into-main state is authoritative. Branches blocked from delete in-cycle become an inter-cycle cleanup task (operator-as-human or a follow-on agent session with delete permission).

- ❌ Leave merged branches accumulating ("someone might need them")
- ❌ Treat a 403 on delete as a release failure (the merge is what matters; the branch is residue)
- ✅ Attempt delete immediately after release; record any 403s in close-out; merged-to-main state is the source of truth

---

## 6. Disconnect-release rules

These are the hard rules δ enforces from the release-effector surface. Each maps to a precondition or a discipline failure mode evidenced by prior cycles.

1. **Do not tag/release before `gamma-closeout.md` on main carries `CDD-Post-Merge-Closeout: complete`.** Filename existence alone may represent pre-operator assurance. The marker gates the tag but is nonterminal.
2. **Manual `git tag` is not allowed.** `scripts/release.sh` is the only way to tag a release.
3. **One tag per release.** Per `release/SKILL.md` §3.6: amend the release commit and reuse the same bare-version tag on CI-red retry. Do not proliferate `3.9.1-fix`, `3.9.1-fix2`.
4. **Annotated tags only.** Generated message via `scripts/generate-release-tag-message.sh`. Lightweight tags lose the structured metadata.
5. **δ blocks release completion until CI is green** (or operator explicitly accepts a known pre-existing failure per §4 step 5).
6. **The tag is the signal.** No separate completion announcement is needed for the disconnect tag itself — the tag appearing on main IS proof that all gate actions completed. Mid-cycle gate actions (per `operator/SKILL.md` §3.3) still require explicit completion confirmations; the disconnect tag does not. (The doctrinal claim lives in `operator/SKILL.md` §3.5; the mechanical fact that the tag is what δ pushes lives here.)
7. **No mixing v-prefixed and bare tags.** Bare `X.Y.Z` only, per `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix".
8. **The release is structural, not consumer-driven.** Untagged post-cycle patches on main are an open boundary; δ does not defer release because "no one is downstream right now." (The doctrinal claim lives in `operator/SKILL.md` §3.4.)

---

## 7. Docs-only disconnect (no tag)

When a cycle ships docs-only — no code change, no version bump — the disconnect is the merge commit on main, not a tag. Release-effector is **not invoked**:

- `scripts/release.sh` is NOT run.
- VERSION is unchanged.
- The merge commit hash IS the disconnect signal.
- Cycle directories move to `.cdd/releases/docs/{ISO-date}/{N}/` (γ does this per `release/SKILL.md` §2.5b, not δ).
- No CHANGELOG ledger row; no tag-message generation; no release CI; no branch cleanup driven from release-effector (any merged-branch cleanup runs as ordinary inter-cycle housekeeping under δ-policy, not under this skill).

See `release/SKILL.md` §2.5b for the docs-only flow and `issue/SKILL.md` for the `docs-only` mode declaration on the issue.

- ❌ Run `scripts/release.sh` on a docs-only cycle "just to be safe" (version drift; synthetic bump)
- ❌ Skip the cycle-directory move because no tag is being cut (loses cycle-to-release association)
- ✅ γ moves cycle dirs to `.cdd/releases/docs/{ISO-date}/{N}/`; merge commit IS the disconnect

---

## 8. What this skill does NOT do

These are intentionally out of scope. Crossing them is a role-boundary violation or an authority claim release-effector does not hold.

- **Do not author `RELEASE.md`.** That's γ per `release/SKILL.md` §2.5.
- **Do not write the CHANGELOG ledger row.** That's β/γ per `release/SKILL.md` §2.4 and §3.8.
- **Do not decide the version bump.** That's γ + β per `release/SKILL.md` §2.2.
- **Do not run `scripts/validate-release-gate.sh` outside the script's flow.** β/CI may run it independently for pre-merge validation; release-effector triggers it transitively via the script's step 4 only.
- **Do not deploy.** Deployment (binary distribution, host installation, daemon restart) is a separate effector surface — see `release/SKILL.md` §2.8 for the current deploy step (β/operator).
- **Do not declare the cycle closed.** δ establishes disconnect; γ performs the subsequent archive move. Terminal closure requires both.
- **Do not override the gate.** Override is `operator/SKILL.md` §4 territory. Release-effector executes the mechanics; override is a δ-policy decision.

---

## 9. Cross-references

- `operator/SKILL.md` §3 (outward gate policy: when the gate fires, who requests, completion-reporting protocol)
- `operator/SKILL.md` §3.4 (doctrinal frame: the tag is the disconnection point)
- `operator/SKILL.md` §3.5 (doctrinal frame: the tag is the signal)
- `operator/SKILL.md` §4 (override protocol — used in §4 step 5 above)
- `operator/SKILL.md` §7 (cycle lifecycle table — Disconnect row points here for mechanics)
- `release/SKILL.md` §2.1 (pre-release validator; invoked transitively by script step 4)
- `release/SKILL.md` §2.4 (CHANGELOG — γ/β surface; release-effector does not author)
- `release/SKILL.md` §2.5 (RELEASE.md — γ surface; release-effector does not author)
- `release/SKILL.md` §2.5a (post-disconnect cycle-dir archive — γ surface; release script does not move it)
- `release/SKILL.md` §2.5b (docs-only disconnect — release-effector is not invoked)
- `release/SKILL.md` §2.6 (release-readiness signaling — β surface; precondition above)
- `release/SKILL.md` §2.7 (release CI polling — β-side text says δ owns this per release-effector §3)
- `release/SKILL.md` §3.6 (amend-don't-re-tag — invoked by §4 step 3 above)
- `release/SKILL.md` §3.8 (TSC scoring — γ/β; release-effector does not score)
- `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix" (canonical tag policy)
- `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Ownership matrix" (gamma-closeout.md gates tag)
- `gamma/SKILL.md` §2.10 (γ closure gate — verifies preconditions above)
- `gamma/SKILL.md` §2 release-prep (γ writes RELEASE.md before disconnect and archives cycle dirs afterward)
- `scripts/release.sh` (the script — δ does not reimplement)
- `scripts/validate-release-gate.sh` (script step 4)
- `scripts/stamp-versions.sh` (script step 5)
- `scripts/check-version-consistency.sh` (script step 6)
- `scripts/generate-release-tag-message.sh` (script step 10)

---

## 10. Embedded Kata

### Scenario

γ has declared post-merge closeout complete. `gamma-closeout.md` carries `CDD-Post-Merge-Closeout: complete` on main. `RELEASE.md` is at repo root. The complete cycle directory remains under `.cdd/unreleased/{N}/`. VERSION still reads the previous release's version. The current release will be `3.67.0`.

### Task

Execute the disconnect release.

### Expected actions

1. **Verify preconditions and the negative case.** Exact marker in `gamma-closeout.md`; an assurance-only file without it must fail; `RELEASE.md`; β's release-ready signal; local main matches `origin/main`.
2. **Run the script.** `scripts/release.sh 3.67.0` (or edit VERSION then `scripts/release.sh`).
3. **Watch the script's output.** Expect: VERSION set, gate validated against `unreleased/`, manifests stamped, consistency OK, release commit, tag annotated, push complete. Assert the script did not move cycle dirs.
4. **Poll release CI.** `gh run list --branch 3.67.0`, then `gh run watch <id> --exit-status`.
5. **CI Green path.** Confirm completion to γ. Delete merged cycle branches. γ then moves the cycle directory to `.cdd/releases/3.67.0/{N}/` and commits the terminal archive.
6. **CI Red path.** Execute §4 runbook: investigate, classify, fix-or-escalate, re-verify. Operator override only for pre-existing infrastructural failures.

### Common failures

- Manual `git tag 3.67.0 && git push --tags` instead of `scripts/release.sh` (skips stamp/consistency check; produces version drift).
- Pushing the tag before the exact post-merge closeout marker exists on main.
- Declaring release complete while CI is red without operator override (discipline failure; the gate has not actually closed).
- Treating a 403 on `git push origin --delete cycle/{N}` as a release failure (the merge is authoritative; the branch is residue — record and move on).
- Letting merged branches accumulate ("someone might need them") instead of deleting at release time.
- Proliferating tags (`3.67.0`, `3.67.0-fix`) instead of amend-and-re-tag on CI-red retry (`release/SKILL.md` §3.6).
