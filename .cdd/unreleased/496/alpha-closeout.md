---
cycle: 496
parent_issue: cnos#496
umbrella: cnos#495 (Sub 1 — first concrete mechanical enforcement primitive crossing the long-arc partition)
master_tracker: cnos#495
pr: cnos#498
merge_sha: b15143d2f2bd4d6aab4af6a607c20b8fcc2bd97c
cycle_branch: cycle/496
base_main_sha_at_cycle_start: 3f57210d95f765ce1884e0f2d6a0868e25b7e1b0
final_head_at_close: 7c57b8748eac378e41d81b686af44f815f455a09
role: α
authored_by: α@cdd.cnos (bootstrap-δ via δ-interface session)
date: 2026-06-25 (UTC)
rounds:
  - R0 (198e795f signal at 2d2e9c9e substantive head) — 8 step commits + signal; D1 (c) / D2 (a) / D3 (a) / D4 (a) forks applied per γ-scaffold §5; β converge at f2569a0d with FN-β-1 advisory; FN-α-1 (pathspec literal globbing) surfaced + fixed inline in step 7
  - R1 (δ-direct; T-486-7 pattern; signal at 14838e6d) — operator iterate-narrowly absorbed: `HEAD@{1}..HEAD` unsafe for multi-commit work phases; explicit pre-work `CN_WAKE_BASE_SHA` baseline step added; 3 commits (renderer + multi-commit CI fixture + self-coherence §R1)
  - R2 (δ-direct; T-486-7 pattern; signal at 7c57b874) — operator iterate-narrowly absorbed: `--diff-filter=ACMR` excludes committed deletions; flag removed entirely (operator preferred form); 2 commits (renderer + golden + substrate + 5 fixture call-sites + new R2 delete-fixture in one commit + self-coherence §R2 in second commit)
---

# α-496 closeout — cnos#496 Sub 1 (cnos#495 umbrella): first mechanical enforcement primitive crossing the long-arc partition

## §1 Cycle close summary

cycle/496 shipped the **mechanical write fence** that prevents package-dispatch wakes from writing `.cn-{agent}/logs/`. R0 landed 8 step commits + a signal commit and β converged on the first review (`f2569a0d`) with zero load-bearing findings and one advisory FN-β-1. The cycle then absorbed **two operator iterate-narrowly rounds** as **δ-direct R1 + R2** (T-486-7 pattern; no α/β re-spawn) — each catching a load-bearing edge case the γ-scaffold's β-prompt fence-audit checklist didn't methodologically anticipate. R1 swapped the commit-graph layer's baseline from `HEAD@{1}` (unsafe for multi-commit work phases) to an explicit pre-work `CN_WAKE_BASE_SHA` step recorded BEFORE the wake's work phase. R2 removed `--diff-filter=ACMR` from the fence's `git log` commands entirely (operator preferred form), restoring coverage to all six change types (A/C/D/M/R/T) so that committed deletions can no longer escape the fence.

This cycle is the **first concrete mechanical-enforcement primitive crossing the cnos#495 long-arc partition** (*"Use intelligence where meaning is unresolved. Use mechanics where the rule is known."*). The rule is known: package-dispatch wakes are non-writers of `.cn-{agent}/logs/`. The mechanic is now installed: belt-and-suspenders — render-time refusal at `cn install-wake` (new exit code 4) for mis-declared providers, plus run-time workflow-level write fence emitting `dispatch_activation_log_write_violation` on detected local writes. Two-layer enforcement matches operator KISS guardrails (no new `cn` subcommands; no Go-stack migration; no `attach` skill split) per γ-scaffold §5 fork choices.

The cycle was the **5th operator iterate-narrowly across the wake-orchestration / wake-logging wave** (cycles 485, 486, 487, 491, 496 — multiple within 496 itself). T-486-12 P1 (operator-final-read defense-in-depth) continues to validate empirically with each firing. The cycle's deliverable is small in surface (one new field; one new exit code; one new failure class; one new convention sub-section) but load-bearing in framing: the long-arc partition between mechanics and intelligence now has its first mechanical witness in production substrate.

## §2 R0 implementation retrospective

### §2.1 R0 step commits

| Step | Surface | Commit SHA | AC |
|------|---------|------------|----|
| 1 | Convention amendment — `AGENT-ACTIVATION-LOG-v0.md` §0.1 ("Wake-class writer ownership (same-repo)") | `07f8e52c` | AC1 |
| 2 | Provider field declarations — cds-dispatch `activation_log_writer:false` + agent-admin `activation_log_writer:true` (explicit, not default) | `1894f526` | AC2 |
| 3 | cds-dispatch prompt scrub — name mechanical enforcement in disallowed-surfaces section | `3127b8f6` | AC3 (prompt surface) |
| 4 | Renderer guard — `cn install-wake` parses field; refuses mis-declarations with exit code 4; omits `attach` when `activation_log_writer:false`; appends fence step to rendered output | `2abf1fb6` | AC4 |
| 5 | Re-rendered cds-dispatch golden + production substrate (sha256 `f2f1970c…` — broken-pathspec render; superseded by step 7) | `446c4491` | AC3 + AC5 surrogate |
| 6 | Dispatch-protocol skill failure-class taxonomy extension — `dispatch_activation_log_write_violation` (§1.3 D10 + §2.7 + §5 D10) | `f9441103` | AC7 secondary |
| 7 | CI fixtures (positive + negative-staged + negative-committed + false-positive resistance + renderer refusal + mis-declaration fixture) + pathspec correction (FN-α-1 fix; rendered sha256 → `d8b77e5a…`) | `2d2e9c9e` | AC4 + AC7 oracles |
| 8 | Self-coherence §R0 + signal commit | `30a6d254` + `198e795f` | signal |

Eight step commits + signal. The signal commit (`198e795f`) is zero file changes; β R0 begins at that SHA per cycle/487 precedent (`e107b7e4`).

### §2.2 Fork choices applied (per γ-scaffold §5)

All four design forks applied per γ recommendation:

- **D1 (c) — KISS:** dispatch wakes don't load `attach` at all. The renderer's emission omits `attach` from the loaded-skills list when `activation_log_writer:false`. The `attach` skill source is unchanged (no split into `attach-identity` + `attach-channel`); operator KISS guardrail honored.
- **D2 (a) — workflow-level fence step:** the renderer appends a `Write fence — verify no .cn-*/logs/ writes` step to the rendered workflow. ~30 lines of bash in the rendered YAML; no new `cn` subcommand; no Go-stack migration. Operator "keep cnos#496 narrow" guardrail honored.
- **D3 (a) — render-time refusal:** `cn install-wake` exits with new exit code 4 when a `role:dispatch + admin_only:false` manifest declares `activation_log_writer:true` (mis-declaration). Belt-and-suspenders: render-time validation is primary defense; run-time fence catches drift if hand-edited workflows bypass the renderer.
- **D4 (a) — job summary only:** no-op dispatch firings surface in the GitHub Actions job summary (UI-visible, structured, ephemeral); zero commits to `.cn-{agent}/logs/`. The post-merge AC5 ground-truth verification will witness this.

No deviations from γ's fork recommendations. No FN-α was filed for fork-choice friction; the operator's KISS preferences and the γ-scaffold's recommendations were mutually reinforcing.

### §2.3 R0 AC verdict (β R0 walk)

β R0 walked all 7 ACs (with AC7 in the amended form per [comment #4792858087](https://github.com/usurobor/cnos/issues/496#issuecomment-4792858087)). Every AC's mechanical oracle returned green. AC5 was marked **green (provisional)** per γ-scaffold §4 interpretation: the CI surrogate (positive write-fence fixture + false-positive resistance fixture) satisfies the mechanical surrogate at PR-review time; operator-final-read post-merge is the ground truth per T-486-12 P1. (AC5 ground-truth verification is deferred to γ-closeout per the task framing — the next post-merge cds-dispatch firing has not been observed at α-closeout authoring time.)

### §2.4 FN-α-1 — git pathspec literal globbing

α surfaced one friction note during R0 implementation, resolved inline before the signal commit. The initial fence used `'.cn-*/logs/'` as a `git status --porcelain` pathspec; git treats bare pathspecs LITERALLY by default (`.cn-*` does NOT expand to `.cn-sigma`). Discovered during Step 7 CI-fixture dry-run when all negative tests reported "no violation detected" against fixtures that clearly contained `.cn-sigma/logs/<file>.md` writes. **Fix:** switched all 6 surfaces (renderer-emitted + 4 CI fixture call-sites + renderer source) to `':(glob).cn-*/logs/**'`. The `:(glob)` magic prefix enables shell-style globbing; the trailing `/**` is required because git treats `.cn-*/logs/` (bare directory) as "directory itself, no contents."

The correction landed in commit `2d2e9c9e` (Step 7). Golden + substrate sha256 went from `f2f1970c…` (Step 5; broken-pathspec render) to `d8b77e5a…` (Step 7; corrected; β R0 review state). γ-scaffold §6 step 4 item 5's example pseudocode used `'.cn-*/logs/'` (same bug) — but the scaffold's example was indicative not authoritative; β R0 verified the fix was NOT drift from the scaffold per scaffold's own framing.

## §3 R1 iteration retrospective — δ-direct (T-486-7 pattern)

### §3.1 The operator iterate-narrowly

β R0 verdict was **converge** at `f2569a0d`. Operator-final-read on PR #498 (T-486-12 P1) returned **iterate narrowly**: `HEAD@{1}..HEAD` is not a reliable pre-run baseline for multi-commit work phases — the reflog's `HEAD@{1}` is the IMMEDIATELY-previous HEAD, so a wake creating `commit 1 (violation) + commit 2 (unrelated)` leaves the violation buried in commit 1, which `HEAD@{1}..HEAD` misses.

The R0 fence used `HEAD@{1}..HEAD` as primary with `$GITHUB_SHA..HEAD` as fallback (per the amended AC7's acceptable-mechanisms enumeration option 2 + 3 mix). The R1 fix inverts the chain to put an explicit pre-work baseline at the top:

- **Primary:** `CN_WAKE_BASE_SHA..HEAD` — captured by a new explicit pre-work step (`Write fence — record pre-work baseline SHA`) that emits `CN_WAKE_BASE_SHA=$(git rev-parse HEAD)` to `$GITHUB_ENV` BEFORE the `claude-code-action` step. Stable for any commit shape (single-commit, multi-commit, no-commit, force-push).
- **Fallback:** `$GITHUB_SHA..HEAD` — the commit checked out at job start; still local-scoped.
- **Last resort:** `HEAD@{1}..HEAD` — kept as a last-resort fallback if neither of the above is available, but no longer the primary mechanism.

### §3.2 R1 step commits

| Step | Surface | Commit SHA |
|------|---------|------------|
| 1 | Renderer source adds pre-work baseline step (emitted BEFORE the work phase); fence's commit-graph layer preference chain → `CN_WAKE_BASE_SHA` → `GITHUB_SHA` → `HEAD@{1}`; re-rendered golden + production substrate; header comment updated. | `0e1078b9` |
| 2 | New CI fixture `AC4 (cycle/496) — write-fence multi-commit baseline (R1)` in `install-wake-golden.yml`: baseline → commit 1 writes `.cn-sigma/logs/20260624-r1-violation.md` → commit 2 unrelated → fence catches the buried commit-1 violation. Informational sub-check compares against old `HEAD@{1}-only` to confirm empirical motivator. | `6135d6f7` |
| 3 | self-coherence.md §R1 records the iterate, the fix, the new sha256, the AC7 oracle update, and the δ-direct absorption verdict. | `14838e6d` |

### §3.3 New rendered sha256 (R1)

- cds-dispatch golden + production substrate: `88d917c81761e72a2fef4efd37aeedbe0d606fad94c9ceea28027f9ca29373f8` (was `d8b77e5a62a98b97a54d02c395ab10bb168f6b71452271b18946c7efe5ead189` at R0).
- agent-admin: **byte-identical to pre-cycle** (FN-2 defensive check still holds; the baseline step is correctly conditional on `activation_log_writer:false` and so does NOT alter agent-admin's rendered output).

### §3.4 Why δ-direct (no α respawn)

Per T-486-7 pattern (cycle/486 precedent for narrow mechanical correctness fixes): the iterate was narrow + mechanical (single field semantic change in the renderer's fence emission + one new CI fixture step); no architectural reframing required. δ-interface (operator's bootstrap session) applied the fix inline as δ-direct R1; α's R0 framing stood; β's R0 converge stood except for the AC7 amended baseline mechanism which R1 corrected mechanically.

β R0 converge was methodologically honest within γ-scaffold §7's three fence checks (local-scope / working-tree / false-positive resistance — did NOT explicitly require a multi-commit case verification, which was a γ-scaffold completeness miss). Operator-final-read defense-in-depth (T-486-12 P1) caught the gap — exactly the empirical case the P2→P1 promotion in cycle/487 was responding to.

## §4 R2 iteration retrospective — δ-direct (T-486-7 pattern; second narrow iterate this cycle)

### §4.1 The operator iterate-narrowly (second firing)

R1 stood until operator-final-read on the R1-amended PR. The second firing of T-486-12 P1 returned **iterate narrowly**: `--diff-filter=ACMR` in the commit-graph layer excludes committed DELETIONS. A dispatch wake could `git rm` a tracked `.cn-{agent}/logs/<file>.md`, commit the deletion, leave the working tree clean, and the fence would return empty.

The R0 + R1 fence used `--diff-filter=ACMR` (A/C/M/R only) on its `git log --name-only` commit-graph layer. Per AGENT-ACTIVATION-LOG-v0 §0.1 + cnos#496 framing, *any path under `.cn-{agent}/logs/` changed by the dispatch wake* is a violation — including ADD, COPY, MODIFY, RENAME, **DELETE**, and TYPE-CHANGE. Deletions and type-changes escape with the ACMR filter.

**Operator's preferred fix:** remove `--diff-filter` entirely (any committed path under `.cn-*/logs/**` in this wake's local commit graph is a violation). Default `git log --name-only` (no diff-filter) covers all six change types — A/C/D/M/R/T. No edge cases; one fewer flag to maintain.

### §4.2 R2 step commits

| Step | Surface | Commit SHA |
|------|---------|------------|
| 1 (code surfaces; single commit; small enough not to split) | Renderer source (header + 3 emission sites) + re-rendered golden + re-rendered production substrate + 5 CI fixture call-sites + new R2 delete-fixture step `AC4 (cycle/496) — write-fence catches committed delete (R2)` (baseline → tracked `.cn-sigma/logs/delete-me.md` → `git rm` + commit → unrelated 2nd commit → fence catches; informational sub-check confirms old `--diff-filter=ACMR` returns empty for the same fixture) | `c2afb90f` |
| 2 (doctrine surface) | self-coherence.md §R2 records the iterate, the fix, the new sha256, the AC7 oracle update, and the δ-direct absorption verdict + FN-γ-R1-1 reinforced (γ-closeout will carry the change-type-coverage extension) | `7c57b874` |

### §4.3 New rendered sha256 (R2 final)

- cds-dispatch golden + production substrate: `cc266be1eb0a61e6c322a88820d6bbefdd2dd890c9c5007c96aa29f72d460508` (was `88d917c8…` at R1; was `d8b77e5a…` at R0).
- agent-admin: **byte-identical to pre-cycle** across R0 + R1 + R2 (FN-2 defensive check held through all three iterations; the wake-class conditional gating is mechanically correct — the renderer's diff-filter and baseline-step emissions are correctly gated by `activation_log_writer:false`).

### §4.4 Local sanity test before commit

A `mktemp -d` fixture confirmed before commit: the new no-filter fence returns `.cn-sigma/logs/delete-me.md` for the committed-delete scenario; the old `--diff-filter=ACMR` returns empty. The new CI fixture preserves this informational comparison in `install-wake-golden.yml` so the empirical motivator stays observable in CI run logs.

### §4.5 Intentional `diff-filter` references remaining

`grep -rn "diff-filter"` shows residual references in commentary + informational-comparison code only (renderer header comment + R2-fixture documentation + R2-fixture's `old_acmr_check` empirical-proof preservation). No `--diff-filter=ACMR` remains in any active fence code path.

## §5 Per-AC verification table (R2 final form)

The R0 oracle column captures β R0's observed result; R1 and R2 columns capture how the iterates changed the oracle (when they did); Final column declares the post-R2 state at PR #498 merge.

| AC | Mechanical oracle | R0 (β R0 walk) | R1 (δ-direct) | R2 (δ-direct) | Final |
|----|-------------------|----------------|----------------|----------------|-------|
| **AC1** Doctrine amendment | `AGENT-ACTIVATION-LOG-v0.md` §0.1 names admin sole writer; package-dispatch non-writer; §0 carve-out | **green** (§0.1 present at lines 37–63; cross-refs cnos#495 + cnos#496 + 2026-06-24 mixed log) | unchanged | unchanged | **green** |
| **AC2** Provider field | `jq -r .activation_log_writer` returns `false` (cds-dispatch) + `true` (agent-admin) | **green** | unchanged | unchanged | **green** |
| **AC3** Dispatch wake does not run attach | rendered `cnos-cds-dispatch.yml` has no `cn attach` invocation; no `attach` in loaded-skills | **green** (only documentation-prose mentions remain) | unchanged | unchanged | **green** |
| **AC4** Renderer guard | mis-declaration fixture → exit 4 + stderr `activation_log_writer mis-declaration:`; cds-dispatch render → exit 0 | **green** (CI fixture encodes oracle) | unchanged | unchanged (the renderer's refusal path is independent of the fence's change-type filter) | **green** |
| **AC5** Live cds-dispatch produces zero `.cn-*/logs/` commits | post-merge: no-op cds-dispatch firing produces zero new commits | **green (provisional)** (CI surrogate satisfies PR-review-time oracle; ground truth deferred to operator-final-read post-merge per T-486-12 P1) | unchanged | unchanged | **green (provisional; awaiting post-merge firing — γ-closeout will record run-id)** |
| **AC6** Historical evidence preserved | `git log cycle/496 -- .cn-sigma/logs/20260624.md` shows no new commits attributable | **green** (empty output; no cycle/496 commit touched the file) | unchanged | unchanged | **green** |
| **AC7** Mechanical write fence (AMENDED) | (a) fence step present with local-scoped check + pathspec `:(glob).cn-*/logs/**`; (b) no `git fetch` / no `origin/main` comparison; (c) CI fixtures: positive + negative (staged) + negative (committed) + false-positive resistance pass | **green** (β walked all 4 fixtures in `mktemp` sandboxes; all passed) | **amended baseline** — primary commit-graph layer baseline now `CN_WAKE_BASE_SHA` (pre-work explicit) → `GITHUB_SHA` → `HEAD@{1}`; new CI fixture (multi-commit baseline) added | **amended change-type filter** — `--diff-filter=ACMR` removed entirely; default no-filter covers A/C/D/M/R/T; new CI fixture (committed delete) added; informational old-ACMR comparison preserved | **green** (R2 final form; local-scope + working-tree + false-positive resistance + multi-commit + all-change-types all covered) |

The R0 oracle status held across all 7 ACs; R1 and R2 amended ONLY AC7's mechanical-implementation details without changing the amended AC7's broader invariants (local-scoped; not remote-state delta; explicit false-positive guard against concurrent admin-wake writes). The iterates strengthened the implementation against edge cases the amendment's acceptable-mechanisms enumeration didn't exhaustively constrain.

## §6 Friction notes (FN-α cluster)

### FN-α-1 — git pathspec literal globbing (R0; resolved inline before β R0 review)

See §2.4. Initial fence implementation used bare `'.cn-*/logs/'` pathspec; git treats it literally; switched to `':(glob).cn-*/logs/**'`. Landed in commit `2d2e9c9e` (Step 7). Pathspec correction was made consistent across all 6 surfaces (rendered workflow + 4 CI fixture call-sites + renderer source). β R0 verified the fix and noted γ-scaffold §6 step 4 item 5's example was indicative not authoritative.

### FN-α-R1-1 — HEAD@{1} multi-commit escape (R1; operator-caught at operator-final-read)

The `HEAD@{1}..HEAD` baseline mechanism in the fence's commit-graph layer is unsafe for multi-commit work phases: a wake that creates `commit 1 (violation) + commit 2 (unrelated)` leaves the violation buried in commit 1, which `HEAD@{1}..HEAD` misses. Fix: explicit pre-work `CN_WAKE_BASE_SHA` baseline step recorded BEFORE the wake's work phase; fence's preference chain becomes `CN_WAKE_BASE_SHA → GITHUB_SHA → HEAD@{1}`.

**Root cause:** γ-scaffold §7 β-prompt audit clause incomplete on **baseline shapes**. The β-prompt's local-write-fence audit had three checks (local-scope / working-tree scope / false-positive resistance) and did NOT explicitly require a multi-commit case verification. β's R0 converge was methodologically honest given the scope; the gap was at the checklist level.

### FN-α-R2-1 — `--diff-filter=ACMR` deletion escape (R2; operator-caught at operator-final-read)

The fence's commit-graph layer used `--diff-filter=ACMR` on its `git log --name-only` commands, excluding committed deletions. A wake that `git rm` + commits a tracked `.cn-{agent}/logs/<file>.md` evades the fence. Fix: remove `--diff-filter` entirely (operator preferred form); default `git log --name-only` covers all six change types (A/C/D/M/R/T).

**Root cause:** γ-scaffold §7 β-prompt audit clause incomplete on **change types**. The β-prompt's audit clause named scopes and shapes implicitly but did NOT enumerate which change types the fence must catch. The amended AC7's acceptable-mechanisms enumeration listed `git diff --name-only` as an option without specifying a diff-filter, so the choice fell into implementation territory; ACMR was a reasonable-but-incomplete default.

### FN-α-defensive-1 — agent-admin renders byte-identical across all 3 iterations

The renderer's wake-class conditional gating on `activation_log_writer` is mechanically correct: the pre-work baseline step (R1) and the no-diff-filter fence (R2) are both emitted ONLY when `activation_log_writer:false`. agent-admin's rendered golden is byte-identical to pre-cycle across R0 + R1 + R2 (FN-2 defensive check from γ-scaffold held through all three iterations). The cleanest empirical confirmation that the wake-class partition is enforced at the renderer's emission boundary, not just at the manifest's declaration boundary.

### FN-α cluster summary

R0 surfaced one FN inline (FN-α-1); both iterates surfaced one each (FN-α-R1-1, FN-α-R2-1); FN-α-defensive-1 is a positive observation. The R1 + R2 FNs share a common root cause — **γ-scaffold's β-prompt fence-audit clause is under-specified on AC-oracle shape and type coverage**. §7 extracts this into a single P1 carryforward.

## §7 Recommendations forwarded to γ-closeout

What α saw across R0 + R1 + R2 that γ-closeout should triage for downstream cycles + γ-template / cnos.cdd skills / wave follow-ups:

1. **FN-γ-R1-1 (PROMOTED to P1 carryforward) — γ-scaffold's β-prompt fence-audit clause must enumerate shape + type coverage explicitly.** R1 + R2 both surfaced because γ-scaffold §7's fence audit was under-specified on **baseline shapes** (single-commit / multi-commit / no-commit / force-push / detached-HEAD) and **change types** (add / copy / modify / rename / **delete** / type-change). Recommend extending T-487-1's "variable consistency table cross-surface scope" with a parallel **AC-oracle SHAPE + TYPE coverage extension** applied specifically to mechanical-guard ACs. Where T-487-1 extended *semantic claim* coverage across surfaces, FN-γ-R1-1 extends *mechanical guard* coverage across edge-case shapes and types. Same mechanism (enumerate-then-walk); different axis.

2. **T-486-12 P1 (REINFORCED 3x in this cycle alone)** — operator-final-read defense-in-depth continues to validate as the right framing. The cycle/487 P2→P1 promotion was responding to exactly this empirical class; cycle/496 is the strongest single-cycle reinforcement to date (two iterates within one PR; both narrow + correct catches; both absorbed δ-direct without re-spawn cost).

3. **AC5 post-merge ground-truth verification deferred to γ-closeout** — at α-closeout authoring time, the next post-merge cds-dispatch firing has not yet been observed. γ-closeout will append the firing run-id + observation outcome once landed. (Surrogate-vs-ground-truth pattern; cycle/487 Stage-2 carryforward precedent.)

4. **Operating model correction — bootstrap-δ-claim conservatism vs live-wake dispatch** — cycle/496 was bootstrap-δ-claimed per cycle/487 Sub 5C precedent (self-application paradox). On reflection: cycle/487's bootstrap-δ-claim of cnos#487 was genuinely self-applicating (the cycle's work was the activation of the very wake that would have claimed it). Cycle/496's situation was subtly different — the live cds-dispatch wake's *no-op exit path* writes to `.cn-{agent}/logs/`, but the live wake's *claim path* invokes δ wake-invoked mode and does the work via δ + γ/α/β; the in-progress cycle would have been routed through δ regardless. The self-application paradox justification was over-cautious. **Recommendation:** future cycles default to live-wake dispatch; reserve bootstrap-δ-claim for cycles whose work directly modifies the wake's own claim sequence.

5. **The `:(glob)` pathspec pattern is load-bearing for any future fence touching `.cn-*` paths** — FN-α-1's correction. Recommend documenting the pattern in `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` §3.10 (or dispatch-protocol §2.7) so the next cycle's α doesn't re-discover it. (β R0 §R0.7 echoed this; carrying forward.)

6. **AC7 amended-form "acceptable mechanisms" enumeration should be expanded with explicit implementation-parameter constraints** — comment #4792858087 named acceptable mechanisms without constraining their implementation parameters (which baselines; which diff-filters). R1 + R2 both surfaced because the constraints were under-specified at amendment time. Recommend an amendment-to-the-amendment shape (or follow-up doctrine micro-cycle) analog to how cnos#454 §2.2 enumerated the claim-sequence steps with implementation-parameter specificity. γ-closeout-eligible.

## §8 Closeout signoff

α-496 closeout complete. **First mechanical-enforcement primitive crossing the cnos#495 long-arc partition is now installed in production substrate.** The wake-class writer-ownership invariant is enforced at two layers: render-time refusal (exit code 4) at the renderer boundary; run-time write fence at the workflow boundary. Belt-and-suspenders held through 3 iterations of edge-case discovery; the final form (R2) covers local-scoped check + working-tree scope + false-positive resistance + multi-commit baseline shapes + all six change types.

The cycle absorbed two operator iterate-narrowly rounds δ-direct (T-486-7 pattern; no α/β re-spawn cost) — the 5th operator iterate-narrowly across the wake-orchestration / wake-logging wave. T-486-12 P1 (operator-final-read defense-in-depth) continues to validate empirically. FN-γ-R1-1 is the single most important carryforward this cycle: γ-scaffold β-prompt fence-audit clauses must enumerate shape + type coverage explicitly, parallel to T-487-1's cross-surface variable consistency extension.

Standing by for β-496 + γ-496 closeouts. AC5 ground-truth verification (post-merge cds-dispatch firing producing zero `.cn-*/logs/` commits) is deferred to γ-closeout; the firing has not yet been observed at α-closeout authoring time.

— α@cdd.cnos, 2026-06-25 (UTC) — post-PR-#498 merge at `b15143d2`; final head at close `7c57b874`.
