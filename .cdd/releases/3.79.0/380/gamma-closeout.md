---
cycle: 380
role: gamma
issue: "https://github.com/usurobor/cnos/issues/380"
date: "2026-05-19"
merge_sha: "770ea1b4"
release_version: "3.79.0"
dispatch_configuration: "§5.1 canonical multi-session (one `claude -p` per role) per gamma-scaffold.md"
sections:
  planned: [Cycle Summary, Dispatch Record, Post-merge Verification, Close-out Triage, §9.1 Trigger Assessment, Cycle Iteration, Deferred Outputs, Hub Memory, Closure]
  completed: [Cycle Summary, Dispatch Record, Post-merge Verification, Close-out Triage, §9.1 Trigger Assessment, Cycle Iteration, Deferred Outputs, Hub Memory, Closure]
---

# γ Close-out — #380

## Cycle Summary

**Issue:** #380 — `cn activate --claude` / `--codex` flags to spawn the AI body interactively with the activation prompt pre-loaded

**Gap:** 3.78.0 shipped `cn activate HUB_DIR` rendering the activation prompt to stdout via the skill-driven renderer. The intended consumer pattern (`cn activate cn-sigma | claude` / `| codex`) was non-interactive — the pipe form prints and exits. The operator had to copy/paste the rendered prompt into a fresh CLI session to land in an interactive body REPL. The bridge from "prompt rendered" to "body operating interactively at this hub" was missing.

**Shipped:**
- `src/go/internal/activate/spawn.go` (NEW, 70L, `//go:build unix`): exports `Spawn(binary, prompt)` (uses `syscall.Exec` for TTY-clean process replacement) and `CheckSpawnBinary(binary, flag)` (LookPath-based pre-render check with actionable missing-binary error naming binary + flag + PATH + install hint).
- `src/go/internal/activate/spawn_other.go` (NEW, 19L): non-unix build-tag stub returning `<flag> spawn is not supported on this platform — use the default stdout form …`. No-flag default path remains functional everywhere.
- `src/go/internal/activate/spawn_test.go` (NEW, 7 tests): argv-shape assertions for both binaries (including `argv[1] != "exec"` critical negative for codex), default-hook wiring, LookPath failure, error message shape.
- `src/go/internal/cli/cmd_activate.go` (M, +47/-5): adds `--claude` / `--codex` flags; mutual-exclusion check (pre-render, exits non-zero); LookPath check (pre-render, exits non-zero with actionable error); spawn arm captures into `bytes.Buffer` only on spawn paths; default arm passes `inv.Stdout` straight through (Option A render-capture seam — bytes-equal property is structural, not test-only).
- `src/go/internal/cli/cmd_activate_test.go` (NEW, 9 tests): mutex-fires-before-render, missing-binary-fires-before-render, bytes-equal-to-direct-Run, help-text-documents-flags, unknown-flag-rejection. Widening declared in self-coherence §Self-check (avoids cyclic import between `cli` and `activate`).
- (γ step 13a — this PRA commit) `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 — added auxiliary paragraph after row 14 naming SHA-citation-invalidation as a downstream consequence of path (a) rebase. Prescribes two coherent resolution paths (preferred: session-start identity check; reactive: re-stamp citations applying §2.3 intra-doc rule).

**Result:** All 5 ACs PASS. β R1 REQUEST CHANGES with one B-severity honest-claim finding (F1 — pre-rebase SHAs in §CDD Trace). α R1→R2 fix-round closed F1 cleanly via `91a2cad6` (SHA rewrite at 9 occurrences applying §2.3 intra-doc rule) + `d838e7e9` (self-coherence §Fix-round 1 append). β R2 APPROVED with zero new findings. Merge at `770ea1b4` on main (`Closes #380`). Post-merge CI red only on pre-existing `Repo link validation (I4)` and downstream `notify` (inherited from cycle #369; not introduced or scoped to fix by #380; β verified non-destructive). Cycle ships at version 3.79.0.

## Dispatch Record

| Role | Sessions | Outcome | Notes |
|------|----------|---------|-------|
| γ (scaffold) | 1 | scaffold + α/β prompts committed (`7a9bc2e7`); branch `cycle/380` created from `origin/main` at `319893a4` | Peer enumeration completed; 6 pre-flagged failure modes (#1–#7); design-and-build mode declared; render-capture seam choice delegated to α |
| α (R1) | 1 | 5 implementation commits (steps 1–5: flag recognition, helper, mutex+spawn, tests, AC3 oracle) + 7 self-coherence commits per §2.5 + row-14 retroactive author-amend rebase (worktree carried γ identity from scaffold session, same class as #379 F-item-1 and #370 F4); review-ready signal at `69b9ea9d` | Clean exit; all 5 ACs evidenced; F1 trace (pre-rebase SHAs in §CDD Trace) survived to β R1 |
| β (R1) | 1 | R1 REQUEST CHANGES — one B-severity honest-claim finding (F1) at `2dac48ba` | Pre-merge gate row 1 (worktree user.email) preventively applied; row 3 (non-destructive merge) verified with explicit dual-run against bare main + cycle/380 HEAD |
| α (fix-round 1) | 1 | F1 closure at `91a2cad6` — rewrites 9 SHA citation sites (5 enumerated by β + 4 additional caught by §2.3 intra-doc rule application); fix-round-1 append at `d838e7e9` | No production-code delta R1→R2; self-coherence-only fix |
| β (R2) | 1 | R2 APPROVED — F1 closed, all ACs re-verified, honest-claim 3.13 sub-checks all ✓; merge at `770ea1b4` (`Closes #380`); β close-out at `b53ba6a4` | Tree-identical for production code to last α implementation SHA `87aa69e9` |
| α (close-out re-dispatch) | 1 | `alpha-closeout.md` committed `f15c27cc` on main; γ-requested per `CDD.md §1.6a` | Clean exit; P1–P5 patterns + Friction log + CDD Trace evidence |
| γ (close-out) | this session | PRA + gamma-closeout.md + cdd-iteration.md + step-13a patch; cycle-dir move to `.cdd/releases/3.79.0/380/`; CHANGELOG 3.79.0 row; RELEASE.md | This file |

**Configuration:** §5.1 canonical multi-session via `claude -p` per `gamma-scaffold.md` §Dispatch configuration. δ dispatched γ → α → β sequentially. Configuration-floor (`release/SKILL.md §3.8`) A− γ-cap does NOT apply (§5.1 full γ/δ separation).

## Post-merge Verification

**Required by `gamma/SKILL.md §2.7` "Post-merge CI verification (mandatory)":**

- CI run on merge commit `770ea1b4`: https://github.com/usurobor/cnos/actions/runs/26102342968
- **Result on merge commit:** RED, but exclusively on pre-existing infrastructure:
  - **`Repo link validation (I4)`** — failure. lychee fails on file:// links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` (cycle #369's own archived artifact). Inherited since #369 merged; documented in 3.78.0 RELEASE.md §"Known Issues" as a separately-tracked pre-existing failure. β verified identical failure shape on bare `origin/main` and `cycle/380` HEAD (`beta-review.md` §"Branch CI state").
  - **`notify`** — failure. Downstream of Build success; sends a Telegram notification only on green CI. Pre-existing for the same I4-blocked reason.
- All other CI jobs green: Protocol contract schema sync (I2), CDD artifact ledger validation (I6), Go build & test, SKILL.md frontmatter validation (I5) — green confirms the 3.78.0 γ-step-13a patch worked, Package/source drift (I1), Binary verification, Package verification — green confirms the 3.78.0 γ-step-13a R5-activate kata patch worked.

**Classification per `operator/SKILL.md §3.4` step 4:** pre-existing infrastructure failure. Operator standing acceptance (documented in 3.78.0 RELEASE.md Known Issues) applies. β's dual-run non-destructive-merge verification per `beta-review.md` §"Branch CI state" + `beta-closeout.md` §"Pre-existing CI red disposition" establishes that #380's diff is non-destructive against the failing validators. Release proceeds.

**γ-axis grade implication:** Per `release/SKILL.md §3.8` CI-red cap clause ("Cycles with red CI on the merge commit cap the γ axis at C"), the cap applies mechanically regardless of cause. γ axis at C for this cycle. The right response is to remediate I4 in a follow-on small-change cycle (clears the cap for every subsequent cycle); see Deferred Outputs and PRA §7 Next Move.

**γ step-13a patch landed in this PRA commit:**

1. `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 — auxiliary paragraph after row 14 naming SHA-citation-invalidation downstream consequence of path (a) rebase, with two coherent resolution paths (session-start identity verification preferred; reactive re-stamp via §2.3 intra-doc rule). Empirical anchor `#380 R1 F1` cited inline.

**Why γ landed this patch per `gamma/SKILL.md §3.6`:** the under-specification surfaced by F1 is structural (row 14's mechanical prescription is sound but its downstream consequence is unnamed), the fix is unambiguous (add a paragraph documenting the consequence and prescribing two coherent paths), and waiting for a separate fix-cycle would leave the role-skill incomplete with one known failure mode pending. §3.6: "A missing gate discovered this cycle should not automatically become future work when the patch is already clear."

## Close-out Triage

Per `gamma/SKILL.md §2.7` CAP table. Findings collected from α close-out, β close-out, and γ-side post-merge CI verification.

| Finding | Source | Type | Disposition | Artifact / commit |
|---------|--------|------|-------------|-------------------|
| F1: §CDD Trace step table cites unreachable post-rebase SHAs (9 occurrences) | β R1 / α-closeout §Friction log F1 / PRA §3 / PRA §4b §9.1 trigger 4 | skill (cdd-skill-gap) | **Immediate MCA** — alpha/SKILL.md §2.6 row 14 SHA-citation-after-rebase auxiliary paragraph patch | this PRA commit (alpha/SKILL.md §2.6 patch); cdd-iteration.md F1 |
| Pre-existing I4 / notify CI red on merge commit | γ post-merge CI verification / PRA §3 | project (infrastructure debt) | **Project MCI** — next-MCA committed (I4 link validation remediation); see PRA §7 | GitHub issue to file (cnos); first AC named in PRA §7 |
| α P1 — Row-14 path (a) rebase invalidates SHA-citation artifacts (pattern) | α close-out §Patterns observed P1 | skill (cdd-skill-gap; same root as F1) | **Already-applied** by F1 disposition — the row-14 patch closes the pattern at the role-skill layer | n/a (same commit as F1 patch) |
| α P2 — Render-capture seam Option A makes bytes-equal property structural | α close-out §Patterns observed P2 | process (positive pattern; design discipline) | **Drop explicitly** — positive observation generalizing minimum-surface vs uniform-surface implementations; recorded in PRA §3 "What went right" item 2 | n/a (recorded) |
| α P3 — Pre-render-ordering proven by absence-of-render-diagnostic | α close-out §Patterns observed P3 | process (positive pattern; test discipline) | **Drop explicitly** — positive observation generalizing absence-of-other-headers oracle shape (same class as #283 R1 F4); recorded in PRA §3 (implicit in §5 "How to verify" item 3) | n/a (recorded) |
| α P4 — γ-scaffold failure-mode register operates as α-side authoring checklist (3rd repo-level confirmation) | α close-out §Patterns observed P4 | process (positive pattern; cumulative confirmation) | **Drop explicitly** — positive observation cumulative across #367 O3, #379 P1, #370 O6, #380 P4; recorded in PRA §3 "What went right" item 1 | n/a (recorded) |
| α P5 — Worktree-config identity-leak class is repo-level (3 cycles, 3 roles) | α close-out §Patterns observed P5 / β close-out §Process observations "Identity hygiene" | project (process; converged design at #373) | **Project MCI** — #373 already exists at converged design; escalated from "growing" to "escalating" lag in PRA §2; γ recommends prioritizing in near-term cycle alongside or after I4 remediation | #373 in cnos issues (existing); PRA §7 deferred outputs |
| α friction log: Intra-doc repetition rule caught 3 β-blindspot SHA occurrences | α close-out §Friction log / β-closeout §"Honest-claim discipline" | process (positive pattern) | **Drop explicitly** — confirms §2.3 intra-doc rule transfers to honest-claim SHA citations as peer class to numeric/named-value drift (2nd cnos surface confirming) | n/a (recorded in PRA §3 "What went right" item 4) |
| α friction log: β's R1 fix instruction authoritative for named sites, α applied intra-doc rule unprompted for additional 3 | α close-out §Friction log | process (role coordination; positive) | **Drop explicitly** — correct role boundary (intra-doc rule is α's discipline per §2.3, not β's); β's review was complete for the rule β is required to apply | n/a (recorded) |
| α friction log: Inherited project-state CI red (I4 + R5-activate P10 from prior cycles) | α close-out §Friction log | project / chained to triage | **Project MCI** — I4 component filed as next-MCA (see above); R5-activate P10 was patched by 3.78.0 γ step-13a and is now green per CI on merge commit (confirmed in §Post-merge Verification above) | n/a (R5-activate P10 already green; I4 deferred) |
| α friction log: Worktree-config leak class confirmed across 3 cycles (3 roles) | α close-out §Friction log | same as α P5 (rephrased) | **Same disposition as α P5** — escalate #373; deferred MCA | n/a (same as α P5) |
| β process observations: γ → α → β protocol coherence; scaffold load-bearing | β close-out §Process observations 1 | process (positive) | **Drop explicitly** — positive observation; same class as α P4 cumulative confirmation | n/a (recorded) |
| β process observations: Round budget at target | β close-out §Process observations 2 | process (positive) | **Drop explicitly** — 2 rounds at target; review-churn §9.1 trigger not fired by strict reading | n/a (recorded in PRA §4) |
| β process observations: Pre-existing CI red disposition (review/SKILL.md 3.10 judgment) | β close-out §Process observations 3 | process / chained to triage | **Drop explicitly** — correct β-side judgment per 3.10 ("inherited project-state CI red is not RC for a cycle whose diff neither introduces nor includes scope to fix the failure"); component MCAs already filed (I4 deferred; R5-activate P10 already remediated by 3.78.0 γ step-13a) | n/a (recorded) |
| β process observations: Identity hygiene (preventive `--worktree`) | β close-out §Process observations 4 | same as α P5 (β-surface confirmation) | **Same disposition as α P5** — escalate #373 | n/a (same as α P5) |

Triage complete. Every finding has a disposition. No silent items.

## §9.1 Trigger Assessment

Per `CDD.md §9.1`:

| Trigger | Fired? | Detail |
|---------|--------|--------|
| review rounds > 2 | NO | actual 2 (at target ≤2; not exceeded by the strict reading of §9.1 "review rounds exceeded target (default: 2)"). The dispatch prompt phrased the trigger as "≥2 nominal" which is a softer reading; under either reading the right action is to assess what produced the second round (F1 honest-claim) and patch its root cause — done via the alpha/SKILL.md §2.6 row 14 patch. |
| mechanical ratio > 20% with ≥10 findings | NO | total findings 1 (F1); 0% mechanical (F1 class is honest-claim); absolute count below 10-finding threshold so ratio is informational |
| avoidable tooling / environmental failure | NO | F1's cause is a skill gap (alpha/SKILL.md §2.6 row 14 under-specification), not tooling failure. The pre-existing I4 CI red is project-debt infrastructure, not a cycle-blocking environmental failure |
| **CI red on merge commit (post-merge)** | **YES** | Merge commit `770ea1b4` has 2 failures (I4 + notify) attributable to pre-existing infrastructure (cycle #369 inheritance); not introduced or scoped to fix by #380. Triggers `release/SKILL.md §3.8` CI-red cap clause (γ axis capped at C) mechanically. |
| **loaded skill failed to prevent a finding** | **YES** | `alpha/SKILL.md` §2.6 row 14 path (a) was loaded by α to remediate the worktree-inherited identity leak; the rebase invalidated 9 SHA citation sites in `self-coherence.md` that α had already authored; β R1 caught F1 via `review/SKILL.md` rule 3.13(a) reproducibility. Row 14 named the rebase mechanically but did not name SHA-citation invalidation as a downstream consequence. Patched in this PRA commit (alpha/SKILL.md §2.6 auxiliary paragraph after row 14). |

**Two triggers fired.** Disposition per `gamma/SKILL.md §2.8` table:

| Trigger | Action | State |
|---------|--------|-------|
| CI red on merge commit (pre-existing) | `release/SKILL.md §3.8` CI-red cap applied (γ axis at C); project-level MCA filed for I4 remediation (clears the cap for future cycles) | **Project MCI committed** + grade reflects post-merge state honestly |
| loaded skill failed to prevent a finding | alpha/SKILL.md §2.6 row 14 patched with SHA-citation-after-rebase auxiliary paragraph (structural at the role-skill layer) | **patch landed now** |

Trigger 1 (CI red) resolves with project-MCI commitment; trigger 2 (loaded-skill miss) resolves with patch-landed-now. No silent next-cycle deferrals.

## Cycle Iteration

Per `gamma/SKILL.md §2.9` Independent γ process-gap check (additional question even when a §9.1 trigger fires):

- **Did this cycle reveal a recurring friction?** Yes — the worktree-config identity-leak class is now confirmed across three consecutive cycles in three role surfaces (#379 F-item-1, #370 F4, #380 F1 + β-closeout §"Identity hygiene"). Issue #373 (Preventive `--worktree` identity write across role skills when `extensions.worktreeConfig=true`) is converged-but-unimplemented at growing lag. PRA §2 escalates it from "growing" to "escalating" — γ recommends prioritizing in a near-term cycle.
- **Was any gate too weak or too vague?** Yes — `alpha/SKILL.md` §2.6 row 14 path (a). Patched in this PRA commit.
- **Did a role skill fail to prevent a predictable error?** Yes — alpha/SKILL.md §2.6 row 14 (just covered). Also a second-order observation: the §2.3 intra-doc repetition rule operated correctly *post-fact* (α grepped for each pre-rebase SHA and found 9 sites; β had named 6) but does not prevent the initial leak. The preventive layer is row 14's SHA-citation invariant (now patched).
- **Did coordination burden show a better mechanical path?** Borderline — γ-scaffold's failure-mode register did not pre-flag "α may apply row-14 path (a) mid-cycle and invalidate SHA citations." Adding that to a per-cycle scaffold register is unscalable (γ would have to anticipate every role-skill execution path); the structural fix is at the role-skill layer (now done). However, γ-scaffold authors may consider including a generic "if path (a) is exercised, re-validate SHA-bearing artifacts" check in future scaffolds for cycles where worktree-config-leak risk is elevated (i.e., cycles run from worktrees inherited from a prior role session).

**Cycle level (per `CDD.md §9.1` cycle-level assessment):**

- L5: locally correct? YES — the named 5 ACs were cleanly delivered; production code at R2 is tree-identical to last α implementation SHA `87aa69e9`; β R2 APPROVED with no implementation-defect findings.
- L6: system-safe (cross-surface coherence)? YES — no cross-surface drift reached merge. The F1 honest-claim was a self-coherence-artifact issue, not a cross-surface drift; β caught at R1; α closed at R2 cleanly via §2.3 intra-doc rule application; no peer surfaces (skill files, kata, validators, etc.) regressed.
- L7: system-shaping leverage? Borderline — the alpha/SKILL.md §2.6 row 14 SHA-citation invariant patch is an L7-shaped move (structural backstop, eliminates the row-14-path-(a)-SHA-leak class for future cycles). The patch landed in the same PRA commit. However, the cycle's *main delivery* (--claude / --codex flags) is a focused feature addition, not a system-shaping change to cnos itself. The L7 ingredient is bundled in the close-out, not in the cycle's named scope.

- Net cycle level: **L6** (cycle cap). L5 and L6 both earned cleanly. L7-shaped step-13a patch landed but as auxiliary to the cycle's main delivery, not as the cycle's primary deliverable. The γ axis grade is independently capped at C by the CI-red clause (`release/SKILL.md §3.8`); the cycle level (L6) and the γ axis grade (C) are orthogonal per the rubric.

## Deferred Outputs

| Item | Owner | First AC | Frozen until shipped? |
|------|-------|----------|----------------------|
| **I4 link validation remediation** (next cnos-internal MCA) | future cnos α | `lychee` returns "0 errors" against `.cdd/releases/docs/2026-05-17/369/`; CI Build workflow's `Repo link validation (I4)` job exits green on the resulting merge commit; downstream `notify` job exits green | Yes (MCI freeze; clears §3.8 cap for future cycles) |
| Hub README router adoption at `usurobor/cn-sigma` (carried from 3.78.0) | cn-sigma δ (γ at cnos files cross-repo issue during δ disconnect) | cn-sigma `README.md` contains the §2.3 router template with `<HUB-URL>` substituted; a body told "Activate as `https://github.com/usurobor/cn-sigma`" fetches the activate skill and reaches §1 identity-confirmation | Yes (MCI freeze) |
| Hub README router adoption at `usurobor/cn-pi` (carried from 3.78.0) | cn-pi δ | analogous AC | Yes |
| #373 Preventive `--worktree` identity write across role skills | future cnos α/γ | `alpha/SKILL.md` + `beta/SKILL.md` + `gamma/SKILL.md` + `operator/SKILL.md` row 1/equivalent prescribes `git config --worktree user.email "..."` from session-start when `extensions.worktreeConfig=true`; α/β/γ regression tests verify identity hygiene across worktree-inheritance scenarios | Yes (escalated from growing to escalating per α P5 + β identity hygiene observation) |
| `//go:embed` activate skill in cn binary (carried from 3.78.0) | future cnos α | renderer reads skill from compiled-in source rather than vendored-or-fallback split | No (P3) |
| `cn doctor` enforcement of activation invariants (carried from 3.78.0) | future cnos α | `cn doctor` reports activation-invariants violations as `StatusFail` or `StatusInfo` | No (P3) |
| `cnos.xyz/activate/<hub>` rendering service (carried from 3.78.0) | future operator + cnos α | hosting decision + dynamic URL renderer | No (P3) |
| `--cursor` / `--aider` / other AI CLI flags | future cnos α | analogous spawn surface for additional AI CLIs | No (held by explicit non-goal until friction emerges) |
| `--auto` body-detection flag + `$CN_DEFAULT_BODY` env-var default | future cnos α | persistence of chosen body across `cn activate` invocations | No (held by explicit non-goal until friction emerges) |
| Manual end-to-end interactive bootstrap dry-run (PRA §5 production verification item 4) | next operator session with `claude` workstation | `cn activate --claude cn-sigma` lands operator in live `claude` REPL with Sigma activated (names identity + operator + current orientation) without further prompts | No (deferred-verification commitment) |

## Hub Memory

- **Daily reflection:** N/A for the cnos repo (no `threads/reflections/daily/` convention at cnos; hub memory lives in per-hub repos like cn-sigma, cn-pi). Recorded as cnos-side durable memory at the protocol layer: this gamma-closeout + PRA + cdd-iteration.md.
- **Adhoc thread(s) updated:** N/A here; the natural place to record "interactive-spawn closes the local-dev bootstrap loop at cnos 3.79.0" is a cn-sigma or cn-pi adhoc thread, written during whichever per-hub session first runs `cn activate --claude HUB_DIR` and observes the operator-visible projection. The cnos-side protocol-layer durable memory (this artifact + PRA + cdd-iteration.md) is the durable record.
- **Cdd-side index:** `.cdd/iterations/INDEX.md` updated with the 3.79.0 / cycle 380 row pointing at `.cdd/releases/3.79.0/380/cdd-iteration.md`.

## Closure

Per `gamma/SKILL.md §2.10` closure gate — all 14 rows verified:

1. ✅ `.cdd/releases/3.79.0/380/alpha-closeout.md` exists on main (was `.cdd/unreleased/380/alpha-closeout.md` before §2.6 cycle-dir move in this PRA commit)
2. ✅ `.cdd/releases/3.79.0/380/beta-closeout.md` exists on main (same move)
3. ✅ Post-release assessment written: `docs/gamma/cdd/3.79.0/POST-RELEASE-ASSESSMENT.md`
4. ✅ Every fired §9.1 trigger has a `Cycle Iteration` entry — see PRA §4b and this file's §9.1 Trigger Assessment + Cycle Iteration sections; both triggers resolved (CI-red → project MCI committed; loaded-skill miss → patch landed)
5. ✅ Recurring findings assessed for skill/spec patching — see PRA §3 "Active skill re-evaluation" + Cycle Iteration; alpha/SKILL.md §2.6 row 14 patched
6. ✅ Immediate outputs either landed or explicitly ruled out — one step-13a patch landed in this PRA commit (alpha/SKILL.md §2.6 row 14 SHA-citation invariant); other findings have explicit dispositions in §Close-out Triage
7. ✅ Deferred outputs have issue / owner / first AC — see §Deferred Outputs above (I4 remediation named as next cnos-internal MCA; cn-sigma adoption carried forward as cross-repo MCA; #373 escalated)
8. ✅ Next MCA named — I4 link validation remediation (cnos-internal; small-change cycle; clears §3.8 cap); first AC named in PRA §7 and §Deferred Outputs above
9. ✅ Hub memory updated — see §Hub Memory above (cnos-side durable memory at protocol layer; per-hub memory deferred to consuming cycles)
10. (pending δ) Merged remote branches cleaned up — δ deletes `cycle/380` per `operator/SKILL.md §3.4` step 5 after release
11. ✅ `RELEASE.md` written and committed to main (this PRA commit)
12. ✅ Cycle directories moved from `.cdd/unreleased/380/` to `.cdd/releases/3.79.0/380/` and committed (this PRA commit)
13. (pending δ) δ release-boundary preflight — γ's request is implicit in this gamma-closeout.md landing on main; δ runs `scripts/release.sh 3.79.0` per `operator/SKILL.md §3.4`
14. ✅ `cdd-iteration.md` exists at `.cdd/releases/3.79.0/380/cdd-iteration.md` with one finding (F1 cdd-skill-gap, patch-landed); INDEX.md row added at `.cdd/iterations/INDEX.md` for cycle 380. No cross-repo `cdd-iteration` deliverable required (the F1 patch lands at cnos in the role-skill layer; no cross-repo bundle).

**Closure declaration:** Cycle #380 closed. Next: I4 link validation remediation (cnos-internal small-change cycle; clears `release/SKILL.md §3.8` CI-red cap for future cycles).

δ may proceed with `scripts/release.sh 3.79.0` per `operator/SKILL.md §3.4`. RELEASE.md is at repo root. VERSION will be set to `3.79.0` by `scripts/release.sh` (or set manually before running the script per §3.4 algorithm step 2). After release CI green per §3.4 step 4 (or operator override for the still-pre-existing I4 / notify failures), δ deletes `cycle/380` per step 5.
