---
cycle: 369
date: "2026-05-17"
issue: "https://github.com/usurobor/cnos/issues/369"
merge_sha: "ff54f2a045d5b8b8574ef6f335d7f30c5f8b23fd"
findings_count: 3
patches_count: 0
mcas_count: 2
no_patch_count: 1
---

# CDD Iteration — #369

Three cdd-`*`-gap findings surfaced during cycle #369 close-out triage. One is **new** (γ-axis rule 3.11b discoverability — filed as #375); two are **recurrence-confirmations** of classes already triaged in cycle #370's iteration (`#370`'s F1 patched via step 13a `4a0115d2`; `#370`'s F2 next-MCA filed as #373).

Format per `cdd/post-release/SKILL.md` Step 5.6b.

## Finding 1 — γ-axis rule 3.11b discoverability gap (late-authored gamma-scaffold.md)

**Class:** `cdd-protocol-gap` (γ-side pre-dispatch gate missing; β-side gate exists at `review/SKILL.md` rule 3.11b).

**Surfaced by:** β R1 D1 (binding finding; review SHA `6835197d`; rule 3.11b fired on missing `.cdd/unreleased/369/gamma-scaffold.md` on cycle branch; no `## Protocol exemption` in sub-issue body). α close-out F4 corroborates: the late-authoring trigger is γ-axis (the artifact is γ's pre-dispatch coordination surface, not an α deliverable).

**Asymmetry:**

- `review/SKILL.md` rule 3.11b is a **binding β-side** gate: if `gamma-scaffold.md` is absent at review time, β must return REQUEST CHANGES (D-severity, `contract / protocol-compliance`).
- `gamma/SKILL.md` §2.5 Step 3a (cycle branch creation) and Step 3b (α/β dispatch) do **not** block dispatch on the artifact's existence. γ can dispatch α without ever authoring `gamma-scaffold.md`.
- Consequence: every γ cycle that forgets the scaffold pre-dispatch pays one full β review round before recovery path (a) clears the binding gate.

**Empirical anchor (cycle #369):**

- γ-369 passed §2.4 issue-quality gate, ran Step 3a (created `cycle/369` from `origin/main@704365d2` and pushed), then ran Step 3b (dispatched α via `claude -p` with `/tmp/alpha-369-prompt.txt`). The scaffold was not written before Step 3b. α completed implementation, signaled review-readiness at `6835197d` (14-row pre-review gate green), β R1 fired D1.
- Recovery: γ committed `.cdd/unreleased/369/gamma-scaffold.md` at `227d2373` (149 lines, frontmatter + 9 H2 sections matching `#367` precedent + 1 §Process note section flagging the late-authoring trigger). β R2 narrowed to a single-row artifact-presence re-check + CI re-verification on new HEAD; APPROVED unconditionally.
- Cost paid: one full β R1 review round (~6 minutes wall-clock for R1; ~3 minutes for R2 narrow) + one γ recovery scaffold authoring + one CI re-run.

**Class definition:** When a binding β-side artifact-presence gate exists on the cycle branch, the γ-side pre-dispatch workflow must mirror the gate as a pre-dispatch check or every γ cycle that forgets the artifact pays one β round-trip. The gate is symmetric in cost — γ pays once-at-scaffold or β pays once-at-review — but currently asymmetric in enforcement.

**Disposition:** **Next-MCA committed.** The fix is a skill-patch in `gamma/SKILL.md` §2.5 Step 3a/3b OR `CDD.md` §1.4 step 3 adding the pre-dispatch existence check as a binding gate. The fix is out of scope for #369 (rider explicitly excludes edits to `gamma/SKILL.md` and `CDD.md`); filed as a follow-on cycle.

**MCA first AC:** A canonical surface (`gamma/SKILL.md` §2.5 Step 3a/3b or `CDD.md` §1.4 step 3) names the check explicitly as "γ cannot proceed to dispatch (Step 3b) until `.cdd/unreleased/{N}/gamma-scaffold.md` exists on the cycle branch", framed as the dual of `review/SKILL.md` rule 3.11b, citing cycle #369 as the empirical anchor.

**MCA owner:** γ (cdd-protocol owner); filed as **#375** (https://github.com/usurobor/cnos/issues/375) — "γ-side pre-dispatch gate for gamma-scaffold.md (rule 3.11b symmetry; cdd-protocol-gap)", P2, parent #366.

## Finding 2 — `extensions.worktreeConfig=true` identity-leak class (recurrence in #369 α F1)

**Class:** `cdd-skill-gap` + `cdd-protocol-gap` (recurrence; pattern is now confirmed across three role-surfaces in two cycles).

**Surfaced by:** α close-out F1 (mid-cycle commit landed as `gamma@cdd.cnos` between the §Gap commit and the §Skills commit; α applied path (a) — `--worktree`-scoped identity write + `git rebase --exec` rewrite of all six prior α commits + `git push --force-with-lease` to cycle branch).

**Pattern (per α close-out O1):** When the repo has `extensions.worktreeConfig=true` at the shared `.git/config` layer (this repo does) and one or more sibling worktrees carry per-worktree `config.worktree` files, plain `git config user.email X` (no `--worktree` flag) writes to the shared layer. Any subsequent process — sibling worktree, hook, parallel-cycle dispatcher — that writes to the shared layer overwrites the value silently. The next read returns the overwriting role's identity. The detection cost is bounded (one `git config --get extensions.worktreeConfig` returning `true` at session start); the preventive cost is one `--worktree` flag per identity write; the recovery cost (path a) scales asymmetrically with cycle size (force-push + rebase exec across N commits).

**Three confirmed surfaces (across two cycles):**

- β-side merge-test recipe — cycle #370 R1 row 1 (β's merge-test worktree-config-leak under shared `.git/config`); citation #301 O8
- α-side close-out re-dispatch — cycle #370 α close-out F4 (α's close-out re-dispatch identity drift)
- α-side mid-cycle commit sequence — **cycle #369 α close-out F1** (this cycle)

**Disposition:** **MCA already filed by #370 — same-class duplicate.** Cycle #370's iteration filed #373 ("Preventive --worktree identity write across all role skills when extensions.worktreeConfig=true", P2, parent #366). #369's F1 is a same-class recurrence within hours of #373's filing, confirming the class is structural (three role-surfaces in two cycles) and the multi-role skill-patch in #373 is correctly scoped. No new MCA filed for #369; this iteration record adds #369's evidence to #373's surface-count without filing a duplicate issue.

**MCA owner / link:** γ (cdd-protocol owner); existing **#373** (https://github.com/usurobor/cnos/issues/373) absorbs #369's evidence. The third role-surface (α mid-cycle commits) strengthens the multi-surface preventive-write framing #373 already targets.

## Finding 3 — validator-literal vs SKILL-prose drift (`§CDD-Trace` vs `§CDD Trace`) — class-confirmation, not a new finding

**Class:** `cdd-tooling-gap` + `cdd-skill-gap` (recurrence; α's first-draft §CDD-Trace header used hyphen-form per `alpha/SKILL.md` §-shorthand).

**Surfaced by:** α close-out F2 (α's first-draft §CDD-Trace header used hyphen-form; rename commit `ff450f6d` aligned to space form pre-signal). The pre-signal correction was triggered by cross-cycle signal: α read cycle #370's R1 verdict (F1 = CDD-Trace drift) on `origin/main` during polling, then applied the rename pre-emptively.

**Pattern:** Same class #370 F1 surfaced. The validator's grep at `cn-cdd-verify:495,573` is `^## CDD Trace` (space); the SKILL prose at `alpha/SKILL.md:217,352` and the validator's own comment at `cn-cdd-verify:480` use `CDD-Trace` (hyphen). Cycle #370 patched all three drifted surfaces in step 13a (commit `4a0115d2`).

**Status:** Step 13a's `4a0115d2` patched the drift on `origin/main` before #369's merge SHA `ff54f2a0`. #369's α surfaced the friction class but did not edit the drifted surfaces (out of scope per AC9). The pre-patch authoring friction (cross-cycle signal as the load-bearing detection) remains.

**Disposition:** **No new patch / no new MCA.** Cycle #370 step 13a (`4a0115d2`) already aligned the three drift surfaces. The pre-patch authoring friction (validator-literal-vs-SKILL-prose detection happening at review time or via cross-cycle parallel-poll, not via a pre-write check) is bounded — α's `alpha/SKILL.md` §2.7 polling discipline does not include parallel-cycle β-review reads, but the mechanical preventive (a pre-review-gate row reading the validator source for `grep -q "^## "` patterns and cross-checking against `alpha/SKILL.md`'s section enumeration) is a candidate for a future α-side skill expansion. **Recorded as "no-patch / one-off"** because (a) the drift surfaces are now aligned post-step-13a, (b) future drifts of the same class require simultaneous edit-divergence between validator and SKILL (a structurally less-frequent event than this cycle's pre-step-13a state), and (c) filing a third issue (after #373 and #375) for a single observational friction would be issue-inflation without proportional MCA-density gain.

**Reason captured for audit:** The class is now well-defined across two cycles; if a third cycle surfaces the same class after #370 step 13a, the mechanical-preventive issue becomes filable. Until then, the cross-cycle signal pattern (α reads parallel-cycle β-reviews on main during polling) is the de-facto detection.

## Summary

| Finding | Class | Disposition | Patch / MCA path |
|---------|-------|-------------|------------------|
| 1 — γ-axis rule 3.11b discoverability gap | cdd-protocol-gap | next-MCA | **#375** (filed pre-closure; P2; parent #366) |
| 2 — worktree-config identity-leak class (recurrence) | cdd-skill-gap + cdd-protocol-gap | next-MCA (existing) | **#373** (filed by cycle #370; #369 adds 3rd surface evidence) |
| 3 — validator-literal-vs-SKILL-prose drift (recurrence) | cdd-tooling-gap + cdd-skill-gap | no-patch / one-off | Patched in cycle #370 step 13a (`4a0115d2`); #369's α confirms class; deferred to a future occurrence |

## Aggregator update

`.cdd/iterations/INDEX.md` row added at cycle #369 same-day; this artifact's path is `.cdd/releases/docs/2026-05-17/369/cdd-iteration.md` after the cycle dir move.
