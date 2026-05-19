---
cycle: 379
role: gamma
issue: "https://github.com/usurobor/cnos/issues/379"
date: "2026-05-19"
merge_sha: "a3bf7892"
release_version: "3.78.0"
dispatch_configuration: "§5.1 canonical multi-session (one `claude -p` per role) per gamma-scaffold.md"
sections:
  planned: [Cycle Summary, Dispatch Record, Post-merge Verification, Close-out Triage, §9.1 Trigger Assessment, Cycle Iteration, Cross-repo Lineage, Deferred Outputs, Hub Memory, Closure]
  completed: [Cycle Summary, Dispatch Record, Post-merge Verification, Close-out Triage, §9.1 Trigger Assessment, Cycle Iteration, Cross-repo Lineage, Deferred Outputs, Hub Memory, Closure]
---

# γ Close-out — #379

## Cycle Summary

**Issue:** #379 — agent/activate skill: single source of truth for agent self-activation

**Gap:** cnos exposed activation as procedure-in-Go (`src/go/internal/activate/activate.go::writePrompt`); every other cnos behavior is a skill artifact. Non-cn bodies (Claude Code on the web, Codex, Claude.ai with WebFetch) could not fetch the activation procedure because it did not exist as a fetchable artifact — they fell into the "I wake up incoherent by default" failure named in cn-sigma `threads/adhoc/20260325-session2-learnings.md`.

**Shipped:**
- `src/packages/cnos.core/skills/agent/activate/SKILL.md` (485 lines): body-agnostic six-item canonical load order; three-tier body-capability matrix; paste-testable README router template; explicit disambiguation from `cdd/activation/SKILL.md`.
- `src/go/internal/activate/activate.go`: `writePrompt` no longer carries the section ordering as in-Go literals; iterates a `readFirst` slice produced by parsing the skill's §4.1 marker-bounded block; fallback path tested.
- `src/go/internal/activate/activate_test.go`: `TestSkillIsSourceOfTruthForReadFirstOrder` (two-phase edit-and-swap, `out1 != out2` coherence assertion) + `TestSkillFallback_NotVendored` + parser unit tests.
- (γ step 13a — this PRA commit) `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §pre-merge gate row 3 tightened (enumerated validator-set; R5-activate named as activate-surface kata exemplar; #379 cited as empirical anchor).
- (γ step 13a — this PRA commit) `src/packages/cnos.core/skills/agent/activate/SKILL.md` `calls:` paths corrected to package-skill-root-relative form (I5 validator green: "66 SKILL.md validated; no findings").
- (γ step 13a — this PRA commit) `src/packages/cnos.kata/katas/R5-activate/{run.sh,kata.md}` P10 updated to the new canonical `## Read first` ordering.

**Result:** All seven ACs PASS. β R1 APPROVED with zero in-cycle findings. Post-merge CI red on two wiring-class regressions (I5 frontmatter + R5-activate kata P10) — both patched in this PRA commit; CI expected green on next push. Cycle ships at version 3.78.0.

## Dispatch Record

| Role | Sessions | Outcome | Notes |
|------|----------|---------|-------|
| γ (scaffold) | 1 | scaffold + α/β prompts committed (`ad9c03f9`); branch `cycle/379` created from `origin/main` at `7a7f7152` | Peer enumeration 6-check; 10 pre-flagged failure modes; design-and-build mode declared |
| α | 1 | 11 commits authoring SKILL.md section-by-section + activate.go evolution + activate_test.go; `self-coherence.md` 5 commits per α/SKILL.md §2.5; row-14 retroactive author-amend rebase (worktree carried γ identity from scaffold session); review-ready signal | Clean exit; all 7 ACs evidenced |
| β | 1 (R1 only) | R1 APPROVED, 0 findings of any severity; merge at `a3bf7892` (`Closes #379`); beta-review.md + beta-closeout.md | Single review round; pre-merge gate rows 1–4 documented; **row 3 application gap surfaced post-merge — see Cycle Iteration below** |
| α (close-out re-dispatch) | 1 | `alpha-closeout.md` committed `a6bc3ca9` on main; γ-requested per `CDD.md §1.6a` | Clean |
| γ (close-out) | this session | PRA + gamma-closeout.md + cdd-iteration.md + step-13a patches; cycle-dir move to `.cdd/releases/3.78.0/379/`; CHANGELOG 3.78.0 row; RELEASE.md | This file |

**Configuration:** §5.1 canonical multi-session via `claude -p` per `gamma-scaffold.md` (escalated from hub default §5.2 per `operator/SKILL.md §5.3` because: 7 ACs at threshold; new contract surface — net-new skill + Go renderer; cross-repo proposal lineage). δ dispatched γ → α → β sequentially. Configuration-floor (`release/SKILL.md §3.8`) A− γ-cap does NOT apply (§5.1 full γ/δ separation).

## Post-merge Verification

**Required by `gamma/SKILL.md §2.7` "Post-merge CI verification (mandatory)":**

- CI run on merge commit `a3bf7892`: https://github.com/usurobor/cnos/actions/runs/26093805709
- **Result on merge commit:** RED. Two new regressions caused by #379 (both wiring-class, both surface peer-validators not aligned with the cycle's intended change):
  - **I5 (SKILL.md frontmatter validation)** — 7 findings, all in the new `src/packages/cnos.core/skills/agent/activate/SKILL.md`'s `calls:` field; paths used `cnos.core/...` prefix instead of validator's required package-skill-root-relative form. Pre-cycle CI green; post-cycle CI red.
  - **Package verification (R5-activate kata, P10)** — kata asserted pre-#379 `## Read first` ordering (`persona < operator < kernel < deps < refl`); AC3 displaced this to `kernel < ca-skills < persona < operator < hub-state < identity`. Pre-cycle CI green; post-cycle CI red.
- Two pre-existing failures (NOT caused by #379, confirmed by baseline run `25992298774` on `7a7f7152` pre-merge): I4 (Repo link validation — dead links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md`); `notify` (telegram bot token issue, per #365 close-out N3).

**γ step-13a patches landed in this PRA commit:**

1. `src/packages/cnos.core/skills/agent/activate/SKILL.md` `calls:` rewritten to package-skill-root-relative paths (`../doctrine/KERNEL.md`, `agent/cap/SKILL.md`, etc.).
   - Locally verified: `./tools/validate-skill-frontmatter.sh --file src/packages/cnos.core/skills/agent/activate/SKILL.md` → ✓
   - Full corpus: `./tools/validate-skill-frontmatter.sh` → "66 SKILL.md validated; no findings"
2. `src/packages/cnos.kata/katas/R5-activate/run.sh` P10 assertion updated to `kernel < persona < operator < {deps, refl}` (the deps + refl relative order is no longer asserted because they share one hub-state line in the new ordering).
3. `src/packages/cnos.kata/katas/R5-activate/kata.md` P10 documentation updated to name the six-item canonical order and reference the source-of-truth (`src/packages/cnos.core/skills/agent/activate/SKILL.md §4.1`).
4. `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §pre-merge gate row 3 tightened (structural backstop preventing the application-gap class from recurring).

CI on the next push (this PRA commit) is expected green for I5 and Package verification. I4 and `notify` remain red from pre-existing baseline; out of #379 scope.

**Why γ landed mechanical fixes per `gamma/SKILL.md §3.6`:** the patches were two-line corrections to peer surfaces, the fix was unambiguous, and the cycle would otherwise close with CI red — leaving a follow-on fix-cycle to restore green. §3.6: "A missing gate discovered this cycle should not automatically become future work when the patch is already clear."

## Close-out Triage

Per `gamma/SKILL.md §2.7` CAP table. Findings collected from α close-out, β close-out, and γ-side post-merge CI verification.

| Finding | Source | Type | Disposition | Artifact / commit |
|---------|--------|------|-------------|-------------------|
| F1: β/SKILL.md row 3 under-specified validator set | γ post-merge / PRA §3 | skill (cdd-skill-gap) | **Immediate MCA** — β/SKILL.md row 3 enumeration patch | this PRA commit (β/SKILL.md §pre-merge gate row 3 tighten); cdd-iteration.md F1 |
| F2: agent/activate skill `calls:` paths violate I5 resolution | γ post-merge / PRA §3 | tooling (cdd-tooling-gap; wiring class) | **Immediate MCA** — activate SKILL.md frontmatter `calls:` rewrite | this PRA commit (activate SKILL.md `calls:` fix); cdd-iteration.md F2 |
| F3: R5-activate kata P10 not updated for new ordering | γ post-merge / PRA §3 | tooling (cdd-tooling-gap; wiring class) | **Immediate MCA** — R5-activate run.sh + kata.md P10 update | this PRA commit (R5-activate kata + doc); cdd-iteration.md F3 |
| α §Debt 1: Hub README router not yet adopted in any hub | α close-out §Debt 1 | project | **Project MCI** — next-MCA committed (cn-sigma adoption); see PRA §7 | Will be filed at cn-sigma during δ disconnect |
| α §Debt 2: Observable-output delta in `## Read first` ordering | α close-out §Debt 2 | project | **Drop explicitly** — intended consequence of AC3 + AC7; documented in `RELEASE.md` §Changed and CHANGELOG note; consumers parsing by section header unaffected | n/a (already documented) |
| α §Debt 3: Renderer fallback duplicates canonical ordering as in-Go constants | α close-out §Debt 3 | project | **Project MCI** — `//go:embed` durable fix deferred; P3 follow-on cycle | PRA §7 Deferred Outputs |
| α §Debt 4: No `cn doctor` enforcement of activation invariants | α close-out §Debt 4 | project | **Project MCI** — P3 follow-on cycle; per issue body Non-goals | PRA §7 Deferred Outputs |
| α §Debt 5: `cnos.xyz/activate/<hub>` deferred | α close-out §Debt 5 | project | **Project MCI** — requires hosting decision; P3 future cycle | PRA §7 Deferred Outputs |
| β process observation 1: γ-scaffolding paid off in round-count = 1 | β close-out §Process observations 1 | process | **Drop explicitly** — positive observation; recorded in PRA §3 "What went right" | n/a (recorded) |
| β process observation 2: α treated AC7 as highest-leverage | β close-out §Process observations 2 | process | **Drop explicitly** — positive observation; recorded in PRA §3 "What went right" | n/a (recorded) |
| β process observation 3: Observable-output delta surfaced honestly | β close-out §Process observations 3 | process | **Drop explicitly** — positive observation | n/a (recorded) |
| β process observation 4: Cross-repo proposal faithfulness | β close-out §Process observations 4 | process | **Project MCI** — cross-repo source-proposal `landed` event addition (this gamma-closeout flags for δ; see §Cross-repo Lineage below) | δ-coordinated; see §Cross-repo Lineage |
| β process observation 5: No `gamma-clarification.md` mid-cycle | β close-out §Process observations 5 | process | **Drop explicitly** — positive observation (well-scoped initial issue) | n/a (recorded) |
| α friction log: Worktree identity inheritance from γ scaffold session | α close-out §Friction log | process | **Drop explicitly** — recovery via α/SKILL.md §2.6 row 14 path (a) was clean; same root family as #287 R1 F3 + cycle #370 F2 → already filed as #373 (deferred for the multi-role `--worktree` preventive patch) | n/a (already filed at #373) |
| α friction log: Pre-existing tests updated, not duplicated | α close-out §Friction log | process | **Drop explicitly** — correct call, β confirmed | n/a (positive) |
| α friction log: γ scaffold failure-mode catalogue as α-side checklist | α close-out §Friction log | process | **Drop explicitly** — positive pattern; named in PRA §3 "What went right" item 2 | n/a (recorded) |

Triage complete. Every finding has a disposition. No silent items.

## §9.1 Trigger Assessment

Per `CDD.md §9.1`:

| Trigger | Fired? | Detail |
|---------|--------|--------|
| review rounds > 2 | NO | actual 1 (within target ≤2 code) |
| mechanical ratio > 20% with ≥10 findings | NO | total findings 3 (post-merge); ratio 0% mechanical / 100% wiring; absolute count below 10-finding threshold so ratio is informational |
| avoidable tooling / environmental failure | **YES** | Post-merge CI red on I5 + Package verification; both regressions are wiring-class (peer surfaces not updated for the cycle's intended change); both bounded mechanical fixes available; patched in this PRA commit |
| CI red on merge commit (post-merge) | **YES** | Merge commit `a3bf7892` has 2 new failures attributable to #379; triggers `release/SKILL.md §3.8` CI-red cap clause (γ axis capped at C) |
| loaded skill failed to prevent a finding | **YES** | β/SKILL.md §pre-merge gate row 3 was loaded; the gate fired (β ran the merge-test recipe in a throwaway worktree per row 3); the application stopped at `go test` because row 3 said "the cycle's own validator (or any CI-equivalent the cycle ships)" without enumerating which validators apply when. `release/SKILL.md §2.1` has the exhaustive list but β/SKILL.md row 3 did not cross-reference it. Patched in this PRA commit (row 3 enumeration). |

**Three triggers fired.** Disposition per `gamma/SKILL.md §2.8` table:

| Trigger | Action | State |
|---------|--------|-------|
| avoidable tooling/environmental failure | `release/SKILL.md §3.8` CI-red cap applied (γ axis at C); two mechanical fixes landed in this PRA commit (F2 frontmatter, F3 kata) | **patch landed now** |
| CI red on merge commit | F2 + F3 mechanical fixes restore CI green on next push; merge commit's red CI is immutable (history is); γ-axis cap recorded in CHANGELOG | **patch landed now** + grade reflects post-merge state honestly |
| loaded skill failed to prevent a finding | β/SKILL.md row 3 patched with enumerated validator-set + per-cycle kata-coverage rule (structural backstop) | **patch landed now** |

All three triggers resolved with patch-landed-now disposition; no next-MCA deferral required for any §9.1 trigger.

## Cycle Iteration

Per `gamma/SKILL.md §2.9` Independent γ process-gap check (additional question even when a §9.1 trigger fires):

- **Did this cycle reveal a recurring friction?** Yes — peer-validator coverage was missed at three layers (α authoring, γ scaffold failure-mode catalogue, β pre-merge gate). The cycle was the first to surface the gap at all three layers simultaneously; previous cycles had only one or two layers active for the surface in question.
- **Was any gate too weak or too vague?** Yes — β/SKILL.md §pre-merge gate row 3. Patched in this PRA commit.
- **Did a role skill fail to prevent a predictable error?** Yes — β/SKILL.md row 3 (just covered) AND α/SKILL.md §2.4 harness audit (could have extended the audit from parsing-schema producers/consumers to the kata that asserts the rendered output). The α-side fix would expand §2.4 wording materially; deferred as no-patch this cycle because β-side row 3 enumeration is the cheaper structural fix and catches the same class. If a future cycle ships another peer-kata regression that β/SKILL.md row 3 enumeration does not catch, the α-side §2.4 expansion becomes the next move.
- **Did coordination burden show a better mechanical path?** Yes — γ-scaffolding failure-mode catalogue could include "β contract-validator coverage" as a generic pre-flag for cycles shipping a new SKILL.md or changing user-visible output. Not patched in this PRA commit (the β-side row 3 fix is the structural backstop; expanding γ-scaffolding indefinitely is unscalable). If the pattern recurs, file the γ-side scaffold-template patch.

**Cycle level (per `CDD.md §9.1` cycle-level assessment):**
- L5: locally correct? YES — the named ACs were cleanly delivered; β R1 APPROVED with zero in-cycle findings.
- L6: system-safe (cross-surface coherence)? NO — two peer surfaces (validator-compliant frontmatter; updated kata assertion) regressed. Cycle caps at L6 miss → L5 by the rubric's "lowest miss" rule. But the cycle DID ship an MCA in the same PRA commit (β/SKILL.md row 3 enumeration) that closes the application-gap class going forward. Per `CDD.md §9.1`: "L7: system-shaping leverage — Did the cycle change the system boundary so the friction class gets easier or disappears?" The row 3 enumeration is an L7-shaped move (structural backstop, eliminates the class).
- Net cycle level: **L6** (cycle cap). L6 because cross-surface drift reached merge (capping at L6), but the structural backstop landed in the same close-out commit (the L7 move is shipped, not deferred), so the cycle does not cap at L5 — it earns the L6 cap honestly with the L7 ingredient bundled in. The γ axis grade is independently capped at C by the CI-red clause; the cycle level (L6) and the γ axis grade (C) are orthogonal per the rubric.

## Cross-repo Lineage

This cycle accepted a cross-repo proposal from `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/` (bundle commit `1a4e25f75`; AC4 capability-matrix refinement folded in post-filing from cn-sigma branch HEAD `bdda457f5` per `threads/adhoc/20260519-git-read-and-untested-limits.md`). Disposition: `accepted`. Delta recorded in the issue body §Source Proposal.

**Pending operator action (flagged for δ):**

The cross-repo source-proposal lifecycle requires a `landed` event on the source-side STATUS after the cnos target work merges (per `gamma/SKILL.md §2.7` and §2.1 cross-repo proposal intake rules). The mirror branches that carry the source-side updates are:

- `cnos:sigma/cross-repo-mirror-agent-activate-skill` @ `212d5239` — source-side mirror at cnos pending operator merge to cnos main (authored by δ on this session prior to γ's invocation, per dispatch prompt)
- `cn-sigma:sigma/cross-repo-status-lineage-update` @ `89049611c` — source-side STATUS + LINEAGE updates at cn-sigma pending operator merge to cn-sigma main

The source-side STATUS at cn-sigma already names:
- `accepted gamma@cnos cnos#379`
- `modified gamma@cnos cnos#379 ac4-capability-matrix-foldin@bdda457f5`

After cycle close, the proposal needs a `landed` event added to the source-side STATUS naming:
- Target cycle: `cnos#379`
- Landed commit: `a3bf7892` (cnos main merge)
- Landed artifact path: `src/packages/cnos.core/skills/agent/activate/SKILL.md`
- Landed release: `3.78.0`

**δ to coordinate the two mirror-branch merges with the operator (both branches pending operator merge); once both land, δ or the cn-sigma project owner adds the `landed` event.** γ does not have authority to merge cross-repo branches and does not have a cn-sigma working copy in this session. This is the only outstanding cross-repo obligation for the cycle.

## Deferred Outputs

| Item | Owner | First AC | Frozen until shipped? |
|------|-------|----------|----------------------|
| Hub README router adoption at `usurobor/cn-sigma` | cn-sigma δ (γ at cnos files cross-repo issue during δ disconnect) | cn-sigma `README.md` contains the §2.3 router template with `<HUB-URL>` substituted; a body told "Activate as `https://github.com/usurobor/cn-sigma`" fetches the activate skill and reaches §1 identity-confirmation | Yes (MCI freeze; see PRA §2) |
| Hub README router adoption at `usurobor/cn-pi` | cn-pi δ | analogous AC | Yes |
| `//go:embed` activate skill in cn binary (durable fix for renderer fallback duplication; α §Debt 3) | future cnos α | renderer reads skill from compiled-in source rather than vendored-or-fallback split | No (P3) |
| `cn doctor` enforcement of activation invariants (α §Debt 4) | future cnos α | `cn doctor` reports activation-invariants violations as `StatusFail` or `StatusInfo` | No (P3) |
| `cnos.xyz/activate/<hub>` rendering service (α §Debt 5) | future operator + cnos α | hosting decision + dynamic URL renderer | No (P3) |
| Cross-repo source-proposal `landed` event addition (cn-sigma STATUS) | δ to coordinate after both mirror-branch merges | source-side `STATUS` has a `landed` event naming target cycle, landed commit, artifact path, release | n/a (mechanical; pending operator merge of two mirror branches) |
| Manual end-to-end body activation dry-run (PRA §5 production verification) | next σ session | body told `https://github.com/usurobor/cn-sigma` self-bootstraps to identity-confirmation via WebFetch | No (deferred-verification commitment) |

## Hub Memory

- **Daily reflection:** N/A for the cnos repo (no `threads/reflections/daily/` convention at cnos; hub memory lives in per-hub repos like cn-sigma, cn-pi). Recorded as cnos-side durable memory at the protocol layer: this gamma-closeout + PRA + cdd-iteration.md.
- **Adhoc thread(s) updated:** N/A here; the cn-sigma adhoc thread on "I wake up incoherent by default" (`threads/adhoc/20260325-session2-learnings.md`) is the natural place to record "the skill that closes this thread shipped at cnos 3.78.0." That update happens during the cn-sigma README router adoption cycle (next MCA) when a body operating at cn-sigma writes the post-adoption reflection.
- **Cdd-side index:** `.cdd/iterations/INDEX.md` updated with the 3.78.0 / cycle 379 row pointing at `.cdd/releases/3.78.0/379/cdd-iteration.md`.

## Closure

Per `gamma/SKILL.md §2.10` closure gate — all 14 rows verified:

1. ✅ `.cdd/releases/3.78.0/379/alpha-closeout.md` exists on main (was `.cdd/unreleased/379/alpha-closeout.md` before §2.6 cycle-dir move in this PRA commit)
2. ✅ `.cdd/releases/3.78.0/379/beta-closeout.md` exists on main (same move)
3. ✅ Post-release assessment written: `docs/gamma/cdd/3.78.0/POST-RELEASE-ASSESSMENT.md`
4. ✅ Every fired §9.1 trigger has a `Cycle Iteration` entry — see PRA §4b and this file's §9.1 Trigger Assessment + Cycle Iteration sections; all three triggers resolved patch-landed
5. ✅ Recurring findings assessed for skill/spec patching — see PRA §3 "Active skill re-evaluation" + Cycle Iteration; β/SKILL.md row 3 patched
6. ✅ Immediate outputs either landed or explicitly ruled out — three step-13a patches landed in this PRA commit
7. ✅ Deferred outputs have issue / owner / first AC — see §Deferred Outputs above (cn-sigma adoption named as next MCA; P3 deferrals named honestly)
8. ✅ Next MCA named — Hub README router adoption at cn-sigma; first AC named in PRA §7
9. ✅ Hub memory updated — see §Hub Memory above (cnos-side durable memory at protocol layer; per-hub memory deferred to consuming cycles)
10. (pending δ) Merged remote branches cleaned up — δ deletes `cycle/379` per `operator/SKILL.md §3.4` step 5 after release
11. ✅ `RELEASE.md` written and committed to main (this PRA commit)
12. ✅ Cycle directories moved from `.cdd/unreleased/379/` to `.cdd/releases/3.78.0/379/` and committed (this PRA commit)
13. (pending δ) δ release-boundary preflight — γ's request is implicit in this gamma-closeout.md landing on main; δ runs `scripts/release.sh 3.78.0` per `operator/SKILL.md §3.4`
14. ✅ `cdd-iteration.md` exists at `.cdd/releases/3.78.0/379/cdd-iteration.md` with three findings (F1 skill-gap + F2/F3 tooling-gaps); INDEX.md row added at `.cdd/iterations/INDEX.md` for cycle 379. No cross-repo `cdd-iteration` deliverable required (the F1 patch lands at cnos; F2/F3 patches land at cnos; the cross-repo *source proposal* lineage is tracked separately via §Cross-repo Lineage above and is not a `cdd-iteration.md` cross-repo bundle)

**Closure declaration:** Cycle #379 closed. Next: cn-sigma README router adoption (cross-repo issue to file during δ disconnect).

δ may proceed with `scripts/release.sh 3.78.0` per `operator/SKILL.md §3.4`. RELEASE.md is at repo root. VERSION will be set to `3.78.0` by `scripts/release.sh` (or set manually before running the script per §3.4 algorithm step 2). After release CI green per §3.4 step 4, δ deletes `cycle/379` per step 5 and coordinates the two cross-repo mirror-branch merges with the operator (see §Cross-repo Lineage above).
