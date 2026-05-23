# α design notes — cycle/399 release-effector skill

## 1. Path & frontmatter

**Path:** `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` (matches the issue body's preferred form).

**Frontmatter shape** — modeled on `release/SKILL.md` + `operator/SKILL.md`. `release-effector` is δ-owned, **task-local** in scope (fires once per release at disconnect, plus the recovery runbook when CI is red), Tier-3 (loaded by δ at the disconnect moment, not by every cycle).

```yaml
---
name: release-effector
description: δ-side mechanics for cutting a release — scripts/release.sh invocation, tag-push, release CI polling, branch cleanup. The platform-actions companion to release/SKILL.md (β-side authoring).
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
  - VERSION decision
  - RELEASE.md on main (γ-authored per release/SKILL.md §2.5)
  - cycle directories under .cdd/unreleased/ ready to move
outputs:
  - bare X.Y.Z annotated tag on origin
  - release CI green (or operator-accepted override)
  - merged cycle/{N} branches deleted
  - mid-cycle gate-action confirmations (per operator/SKILL.md §3.3)
requires:
  - δ has read operator/SKILL.md §3 (outward gate policy)
  - γ has shipped RELEASE.md + cycle-dir moves per release/SKILL.md §2.5–§2.5a
  - gamma-closeout.md exists on main (CDD.md §5.3b)
calls:
  - release/SKILL.md
  - operator/SKILL.md
---
```

The wording is deliberately surgical:
- `description` calls out `scripts/release.sh` and names release-effector as δ-side mechanics-companion to release/SKILL.md (β-side authoring).
- `calls` lists `release/SKILL.md` and `operator/SKILL.md` because release-effector reads β-side and δ-policy upstream contexts.
- `scope: task-local` matches release/SKILL.md (the release event is task-local).

## 2. Boundary decisions (β/δ + operator/effector)

| Concern | Owner | Lives in |
|---|---|---|
| RELEASE.md drafting | γ | `release/SKILL.md §2.5` (read-only here) |
| CHANGELOG ledger row | γ + β | `release/SKILL.md §2.4, §3.8` (read-only here) |
| Cycle-directory move at release time | γ | `release/SKILL.md §2.5a` (read-only here) — note: `scripts/release.sh` step 6a also moves cycle dirs, so the script is the durable mechanical fallback; release-effector skill names this |
| Pre-release validator (`validate-release-gate.sh`) | β/CI | `release/SKILL.md §2.1` (read-only here); release-effector cross-references because the validator runs inside `scripts/release.sh` step 4 |
| `scripts/release.sh` invocation (stamp + tag) | δ | **release-effector** (moved from `operator/SKILL.md §3.4`) |
| Annotated tag creation (bare `X.Y.Z`) | δ | **release-effector** (moved from `operator/SKILL.md §3.4`) |
| Release-CI polling (`gh run list --branch <tag>`) | δ | **release-effector** (moved from `operator/SKILL.md §3.4` step 4 + §6-step-Gate sub-runbook) |
| CI-red recovery runbook | δ | **release-effector** (moved from `operator/SKILL.md §6` step 6 substeps 1-5) |
| Merged-cycle-branch cleanup | δ | **release-effector** (moved from `operator/SKILL.md §3.4` step 5) |
| Gate-action completion reports | δ | `operator/SKILL.md §3.3` (stays — it's policy about when/how δ confirms, not effector mechanics) |
| Mid-cycle gate observation policy (§3.1–§3.2) | δ | `operator/SKILL.md` (stays — δ-policy doctrine) |
| Disconnect-tag-as-signal (§3.5) | δ | **stays in `operator/SKILL.md`** — this is **doctrine** ("the tag IS the signal") not mechanics. The mechanical fact of pushing the tag moves; the doctrinal claim that the tag is the signal is δ-policy. Section §3.5 in operator/SKILL.md is updated to cross-reference release-effector for the mechanics. |
| Lifecycle "Disconnect" row in §7 table | δ | **stays in `operator/SKILL.md`** with the cross-ref updated to point at release-effector |
| §3a inward-membrane (implementation-contract enrichment) | δ | **stays in `operator/SKILL.md`** — orthogonal to release |
| §10 wave coordination | δ | **stays in `operator/SKILL.md`** — orthogonal to per-release mechanics |

**The pivot:** §3.4 (the disconnect-release algorithm — 5 numbered steps + manual-tagging prohibition + tag-message generation paragraph + cases) is the heart of the relocation. §6 step 6 sub-runbook (5 recovery steps) goes with it. §3.5 stays as policy.

## 3. New skill content — section outline

```
release-effector/SKILL.md
├── frontmatter (above)
├── # Release Effector (δ)
├── ## Core Principle
│   "The release tags the disconnection point. δ runs the script,
│    polls CI, and cleans the branches. δ executes the mechanics
│    that turn γ-closure-declaration into a tagged, distributable
│    release."
├── ## Preconditions (what must be true before δ runs scripts/release.sh)
│   - gamma-closeout.md on main (CDD §5.3b)
│   - RELEASE.md at repo root (release/SKILL.md §2.5; CDD §5.3b)
│   - cycle dirs under .cdd/unreleased/{N}/ ready to move
│     (or already moved by γ per release/SKILL.md §2.5a — script handles both)
│   - VERSION decided (semver bump per release/SKILL.md §2.2)
│   - β has signaled "release ready for δ tag" in beta-closeout.md
│     (release/SKILL.md §2.6)
├── ## 1. Run scripts/release.sh
│   - Single-command release: edit VERSION, run scripts/release.sh
│     (or scripts/release.sh X.Y.Z to pass version directly)
│   - What the script does (mechanical narration, mirrors scripts/release.sh):
│     1. set VERSION (if arg given)
│     2. confirm on main + up to date with origin/main
│     3. confirm tag X.Y.Z doesn't already exist
│     4. validate release gate (scripts/validate-release-gate.sh)
│     5. stamp manifests (scripts/stamp-versions.sh)
│     6. consistency check (scripts/check-version-consistency.sh)
│     7. move .cdd/unreleased/*/ → .cdd/releases/X.Y.Z/
│     8. stage & commit "release: X.Y.Z" if anything changed
│     9. warn-prompt if CHANGELOG.md missing the version
│    10. generate annotated tag message (scripts/generate-release-tag-message.sh)
│    11. create annotated bare tag X.Y.Z
│    12. push origin main --tags
│   - Manual tagging is not allowed.
│   - Why: prevents version drift (VERSION/cn.json/package manifests
│     disagreeing). See DISPATCH-FAILURE-EVIDENCE.md cycle #84 failure 3.
├── ## 2. Tag policy (bare X.Y.Z, no v-prefix)
│   - All tags are bare X.Y.Z. Canonical reference: CDD.md §5.3a
│     ("Tags are bare X.Y.Z everywhere ... v-prefixed tags are legacy
│     and warn-only").
│   - Same convention applies to: VERSION file, git tag, branch name
│     version segment, CHANGELOG row, RELEASE.md, snapshot dir.
│   - Annotated tags (not lightweight) with structured message generated
│     by scripts/generate-release-tag-message.sh.
├── ## 3. Poll release CI
│   - After tag push, run `gh run list --branch <tag>` and monitor
│     workflow completion.
│   - **δ blocks release completion until CI is green.**
│   - CI Green → release complete; report per operator/SKILL.md §3.3.
│   - CI Red → execute recovery runbook (§4 below).
├── ## 4. CI-red recovery runbook
│   (moved verbatim — except for renumbering — from operator/SKILL.md §6 step 6)
│   1. Investigate release logs
│   2. Classify (release-specific vs pre-existing infra)
│   3. Fix or escalate
│      - Fixable: fix, re-tag/re-run, poll again
│        (re-tag procedure per release/SKILL.md §3.6 — amend, force-push,
│        do not proliferate tags)
│      - Pre-existing: document, escalate to operator-as-human, do NOT
│        declare release complete
│   4. Re-verify — poll CI again after fix attempts
│   5. Operator override — explicit operator acceptance required for
│      known pre-existing failures (escape hatch for v3.66.0/v3.67.0
│      smoke failures)
│   - **The gate does not close until CI is green or operator
│     explicitly accepts the failure.**
├── ## 5. Branch cleanup
│   - After release complete: delete merged cycle/{N} branches.
│   - Also delete γ session branches (harness-given `claude/...`
│     or operator-named `gamma/session-{N}`) that were used during
│     the cycle.
│   - Mechanical procedure (mirrors release/SKILL.md §2.6a):
│     `git branch -r --merged origin/main | grep -v main | grep -v HEAD \
│       | sed 's/origin\///' | xargs -I{} git push origin --delete {}`
│   - No orphan γ session branches survive past closure.
│   - Some harnesses block remote branch deletion (403 on push origin
│     --delete). δ attempts; if blocked, records the failure in the
│     cycle's close-out and proceeds — the merged-into-main state is
│     authoritative.
├── ## 6. Disconnect-release rules
│   - **Do not tag/release before gamma-closeout.md exists on main.**
│   - **Manual `git tag` is not allowed.** scripts/release.sh is the
│     only way to tag.
│   - **One tag per release.** Per release/SKILL.md §3.6: if CI fails
│     post-tag, amend the release commit, delete-and-recreate the
│     bare-version tag, force-push. Do NOT proliferate `3.9.1`,
│     `3.9.1-fix`, `3.9.1-fix2`.
│   - **The tag is the signal.** No separate completion announcement
│     is needed for the disconnect tag itself; the tag appearing on
│     main IS proof that all gate actions completed. Mid-cycle gate
│     actions (per operator/SKILL.md §3.3) still require explicit
│     completion confirmations.
├── ## 7. Docs-only disconnect
│   - When a cycle ships docs-only (no version bump, no tag):
│     - scripts/release.sh is NOT run.
│     - Disconnect signal is the merge commit hash on main.
│     - Cycle dirs move to .cdd/releases/docs/{ISO-date}/{N}/
│       (release/SKILL.md §2.5b — γ does this, not δ).
│     - No release-effector action required.
├── ## 8. Kata (embedded)
│   - Scenario: γ has declared closure, RELEASE.md is on main,
│     cycle dir is on main, VERSION says 3.67.0.
│   - Expected actions:
│     1. Verify gamma-closeout.md on main
│     2. Verify RELEASE.md at repo root
│     3. Run scripts/release.sh
│     4. Poll gh run list --branch 3.67.0
│     5. CI green → confirm to γ; delete merged branches
│     6. CI red → execute §4 runbook
│   - Common failures:
│     - manual `git tag 3.67.0 && git push --tags` (skips stamp/consistency)
│     - tag pushed before gamma-closeout.md on main (gate violated)
│     - CI red but δ declares complete anyway (override discipline broken)
│     - leaving merged branches accumulating
```

This shape is roughly 200 lines once written out — proportionate to release/SKILL.md (~366 lines) and operator/SKILL.md disconnect-release surface area (~50 lines).

## 4. Edits to operator/SKILL.md

The §3.4 body and §6 step 6 sub-runbook relocate. Replacements:

**§3.4 — full replacement:**

Current text: "Cut the release — disconnect the triad's final state" section, ~30 lines including 5-step algorithm + manual-tagging prohibition + tag-message generation paragraph + 4 bullets.

New text: cross-reference stub naming the doctrinal claim (the disconnect happens via release-effector) + pointer:

> ### 3.4. Cut the release — disconnect the triad's final state
>
> After all post-cycle work lands on main (γ's PRA + skill patches, δ's own session patches), δ cuts the release. **This is not optional** — the release is how δ disconnects the triad's output into a distributable, tagged whole.
>
> The triad's work is not complete until it is tagged. Untagged post-cycle patches on main are an open boundary — the triad's output is still entangled with whatever comes next. The tag is the disconnection point.
>
> **The mechanics live in `release-effector/SKILL.md`.** That skill owns: `scripts/release.sh` invocation, the bare-`X.Y.Z` tag policy (per CDD §5.3a), release-CI polling, the CI-red recovery runbook, and merged-cycle-branch cleanup.
>
> Two gate rules δ enforces from this surface:
> - **Do not tag/release before `gamma-closeout.md` exists on main.**
> - **δ blocks release completion until CI is green** (or operator explicitly accepts a known pre-existing failure per the release-effector recovery runbook).

(8 lines total; no algorithm steps, no script invocation, no recovery runbook — just the doctrinal frame and two gate-rule reminders. Algorithm + runbook live in release-effector.)

**§6 step 6 (Gate algorithm step) — collapse the CI-recovery inline runbook:**

Current text has step 6 carrying the full 5-step recovery runbook inline. The runbook moves to release-effector §4. The §6 step 6 text becomes:

> 6. **Gate** — execute external actions: push main, tag, release, branch cleanup. **Do not tag/release before `gamma-closeout.md` exists on main.** Release mechanics (tag/release/CI polling/branch cleanup) and the CI-red recovery runbook live in `release-effector/SKILL.md`; see §3.4 below for the doctrinal frame.

(2 lines; runbook content gone.)

**§3.5 — light edit:**

The §3.5 text says "The disconnect tag (§3.4) is git-observable. γ and all future agents can see it. ..." The §3.4 reference stays valid (§3.4 stub still names the tag-is-disconnect doctrine), so §3.5 is unchanged or gets a one-word update.

**§7 lifecycle table — one row update:**

Current: `| Disconnect | Cut the release — `scripts/release.sh` after γ closure declaration | γ close-out + δ session patches on main |`

New: `| Disconnect | Cut the release — see `release-effector/SKILL.md` | γ close-out + δ session patches on main |`

(`scripts/release.sh` reference removed — it's no longer in operator/SKILL.md per AC2.)

**Kata A §9 — line 520:**

Current: `6. When γ declares closure: execute branch cleanup`

New: `6. When γ declares closure: cut the release per release-effector skill (run scripts/release.sh, poll CI, branch cleanup)`

(The kata expected-action list mentions "branch cleanup" which is one of the AC2 grep terms. Replacement preserves the action while routing the mechanics through the new skill. The phrase "branch cleanup" still appears, but now in a cross-reference context — see AC2 oracle which permits cross-references.)

Wait — AC2 oracle says `rg "scripts/release.sh|tag creation|branch cleanup|release CI" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md returns 0 hits in normative position (only cross-references to release-effector)`. "Cross-references" is permitted. The kata line stays as a cross-reference because it points the reader at release-effector.

Actually — to be safe and to maximally satisfy AC2, I'll rewrite the kata line as: `6. When γ declares closure: cut the release per release-effector/SKILL.md` — drops "branch cleanup" entirely.

## 5. Edits to release/SKILL.md

Cross-references to update:

- **Line 35** (Core Principle): "See `CDD.md` §1.4 β algorithm, γ algorithm Phase 5a, and `operator/SKILL.md` §3.4." → "See `CDD.md` §1.4 β algorithm, γ algorithm Phase 5a, and `release-effector/SKILL.md`."
- **Line 238** (2.7 Wait for release CI): "**δ owns release CI polling** per `operator/SKILL.md` §3.4 step 4 — ..." → "**δ owns release CI polling** per `release-effector/SKILL.md` §3 — ..."
- **Line 195** (2.5b docs-only): "The merge commit hash IS the disconnect signal. No `scripts/release.sh` invocation." — stays as-is. This is about *not* invoking the script; it's a β-side note about which disconnect path applies. The mention of `scripts/release.sh` is contextual ("no script invocation"), not a procedure for invoking it. Keep this — it's β-side accurate.

No other release/SKILL.md content moves. β-side flow (readiness, version decision, CHANGELOG, RELEASE.md, cycle-dir move, validators, scoring) all stays per AC5.

## 6. Edits to gamma/SKILL.md

Cross-references:

- **Line 65** (load order): "load `operator/SKILL.md` — δ owns dispatch execution ..., release-phase gate execution (tag push, branch cleanup, release CI), and the disconnect release (§3.4)" → "load `operator/SKILL.md` — δ owns dispatch execution ... and outward gate policy; `release-effector/SKILL.md` owns the disconnect release (scripts/release.sh, tag push, release CI, branch cleanup)."
- **Line 89** ("Step 17 → δ disconnect release"): "Step 17 → δ disconnect release (`operator/SKILL.md` §3.4) — δ's step, not γ's" → "Step 17 → δ disconnect release (`release-effector/SKILL.md`) — δ's step, not γ's".
- **Line 450**: "δ runs `scripts/release.sh` (stamp + tag) but does not author artifacts." — this is a γ-side rule (γ does authoring; δ runs the script). The `scripts/release.sh` reference here is informational about the script's existence; the canonical mechanics live in release-effector. Update to: "δ runs `scripts/release.sh` per `release-effector/SKILL.md` (stamp + tag) but does not author artifacts."
- **Line 458**: "Both must be committed before γ requests the disconnect release from δ (§2.10 step 15 → operator/SKILL.md §3.4)." → "Both must be committed before γ requests the disconnect release from δ (§2.10 step 15 → release-effector/SKILL.md)."

## 7. Edits to activation/SKILL.md

- **Line 598**: "Tagging is a δ action executed via `scripts/release.sh` per `operator/SKILL.md §3.4`; manual `git tag` is not allowed." → "Tagging is a δ action executed via `scripts/release.sh` per `release-effector/SKILL.md`; manual `git tag` is not allowed."

## 8. Edits to CDD.md

- **Line 1183**: `| F3 | Missing `gamma-closeout.md` before δ tag | `.cdd/unreleased/{N}/gamma-closeout.md` exists on main before δ runs `scripts/release.sh` | §5.3b; `operator/SKILL.md` §3.4 |` → update second cite to `release-effector/SKILL.md`.

## 9. Edits to COHERENCE-CELL.md

- **Line 254**: "Today, `release/SKILL.md` and `operator/SKILL.md` both carry pieces of this. The refactor target is to name release as boundary effection and to relocate the mechanics (`scripts/release.sh`, CI polling, tag-message generation, branch cleanup) to the substrate where they belong." — This text describes the not-yet-done refactor as a target. Cycle 399 is *part of* doing it. Update prose to reflect the cnos#399 partial completion: "Today, `release/SKILL.md` carries β-side authoring; `release-effector/SKILL.md` (cnos#399) carries δ-side mechanics (`scripts/release.sh`, CI polling, tag-message generation, branch cleanup); `operator/SKILL.md` carries δ-policy frame. Phase 4a (δ split) will further relocate δ-policy into `delta/SKILL.md`."

## 10. AC walkthrough at design time

- **AC1**: `release-effector/SKILL.md` exists with frontmatter. ✓ (planned above)
- **AC2**: After edits, `rg "scripts/release.sh|tag creation|branch cleanup|release CI" operator/SKILL.md` returns 0 hits in normative position. The cross-reference §3.4 stub does NOT use any of those phrases ("Release mechanics ... live in `release-effector/SKILL.md`"). The §6 step 6 collapsed text does NOT use the phrases either. The §7 row uses none. Kata line is rewritten to "cut the release per release-effector/SKILL.md". Likely 0 hits. ✓
- **AC3**: Tag policy text in §2 of release-effector with explicit "CDD §5.3a" reference. ✓
- **AC4**: §3.4 doctrinal frame stays in operator/SKILL.md (the "release is the disconnection point" claim); mechanics move to release-effector. Both surfaces carry their share; no content lost. ✓
- **AC5**: release/SKILL.md edits are cross-reference updates only — §2.4 CHANGELOG, §2.5 RELEASE.md, §2.5a cycle-dir move, pre-release validators stay. ✓
- **AC6**: gamma/SKILL.md, release/SKILL.md, activation/SKILL.md, CDD.md, COHERENCE-CELL.md cross-refs all updated. ✓
- **AC7**: `scripts/release.sh` not modified. release-effector skill's §1 narration of the 12 steps matches the script's actual flow (verified by reading scripts/release.sh). ✓

## 11. CDD Trace continuation

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 5 Mode | this file | cdd/design, cdd/operator (read), cdd/release (read) | design-and-build; γ+α+β collapsed on δ-as-agent (§5.2) |
| 6 Artifacts | release-effector/SKILL.md (planned); 7 cross-ref updates planned | — | Section outline mapped, edits enumerated |
