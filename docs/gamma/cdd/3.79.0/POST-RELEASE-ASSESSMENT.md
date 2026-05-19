## Post-Release Assessment — 3.79.0

**CI status on merge SHA:** red — https://github.com/usurobor/cnos/actions/runs/26102342968 (merge commit `770ea1b4`; failing jobs are `Repo link validation (I4)` and `notify` (downstream of I4). Both are **pre-existing infrastructure** carried forward from cycle #369, listed as a Known Issue in `3.78.0/RELEASE.md` §"Known Issues", and confirmed by β-run dual verification on bare `origin/main` and the merge tree showing identical failure shape on both. No regression introduced or scoped to fix by #380. Per `operator/SKILL.md` §3.4 step 4 (operator override for known pre-existing failures), the release proceeds; per `release/SKILL.md` §3.8 CI-red cap clause, the γ axis caps at **C** mechanically.

### 1. Coherence Measurement

- **Baseline:** 3.78.0 — α B+, β B+, γ C, C_Σ B−
- **This release:** 3.79.0 — α B+, β A, γ C, C_Σ B−
- **Delta:** β improved (clean review work; F1 caught at R1 with actionable evidence; pre-merge gate rows 1 and 3 preventively applied; β-closeout names worktree-config leak class as repo-level). α held (1 B-severity honest-claim at R1, cleanly closed at R2). γ held at C (same CI-red cap, but from a different cause — 3.78.0's red was introduced regressions γ patched in the PRA; 3.79.0's red is purely inherited #369 infrastructure debt). C_Σ held at B−.
- **Coherence contract closed?** Yes. The cycle's named gap — `cn activate HUB_DIR` rendered the activation prompt to stdout but did not deliver it into an interactive AI-CLI session — is closed. `cn activate --claude HUB_DIR` and `cn activate --codex HUB_DIR` replace the cn process with the AI CLI via `syscall.Exec`, dropping the operator into a live REPL with the body pre-activated. Default no-flag stdout behavior is preserved bytes-equal to 3.78.0 by construction (Option A render-capture seam: no buffer indirection on the default arm). Mutual exclusion and missing-binary errors fire pre-render, proven by the "nonexistent HUB_DIR" oracle that asserts the pre-render diagnostic plus the *absence* of the renderer's own error in stderr. The operator-visible projection is: one command (`cn activate --claude cn-sigma`) bootstraps a body session end-to-end on a workstation with `claude` installed.

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| #380 | `cn activate --claude` / `--codex` flags | feature | issue body | shipped (this release) | none |
| (carried from 3.78.0) | Hub README router adoption in cn-sigma | feature | template in `agent/activate/SKILL.md` §2.3 | not started | growing |
| (carried from 3.78.0) | Hub README router adoption in cn-pi | feature | template in `agent/activate/SKILL.md` §2.3 | not started | growing |
| (carried from 3.78.0) | `//go:embed` activate skill in cn binary | feature | named in #379 α §Debt 3 | not started | growing |
| (carried from 3.78.0) | `cn doctor` enforcement of activation invariants | feature | named in #379 α §Debt 4 | not started | growing |
| (carried from 3.78.0) | `cnos.xyz/activate/<hub>` rendering service | feature | named in #379 §Scope | not started | growing |
| #373 | Preventive `--worktree` identity write across role skills when `extensions.worktreeConfig=true` | process | converged | not started | **growing → escalating** (#380 P5 confirms the class is repo-level — three cycles, three roles: #379 F-item-1, #370 F4, #380 F1 + β-closeout) |
| (deferred from #380) | I4 link validation — remediate stale file:// references in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` | process | scope known | not started | growing (forcing function active — every cycle until fixed pays γ-axis cost per §3.8 cap) |
| (deferred from #380) | `--cursor` / `--aider` / other AI CLI flags | feature | named in #380 issue §Non-goals | not started | low (held by explicit non-goal until friction emerges) |
| (deferred from #380) | `--auto` body-detection flag + `$CN_DEFAULT_BODY` env-var default | feature | named in #380 issue §Deferred | not started | low (held by explicit non-goal until friction emerges) |

**MCI/MCA balance:** **Freeze MCI** — 6 issues at "growing" lag (5 carried from 3.78.0 + escalating #373). The freeze that 3.78.0 declared remains in effect. Two new low-lag items (`--cursor` / `--auto`) are held by explicit non-goal scope, not by freeze.

**Rationale:** 3.78.0 declared freeze because #379 produced four growing-lag follow-ons against one primary deliverable. 3.79.0 ships the operator-facing handover that closes 3.78.0's primary consumer-loop gap on the local-dev flow, but the hub-side adoption (cn-sigma README router, cn-pi README router) still has not happened. The fastest unfreeze move is now the cross-repo cn-sigma README router adoption: it activates the cycle's already-shipped artifact in a consuming hub and validates the cycle's "one command bootstraps a body" outcome end-to-end. Continue holding all new substantial design.

### 3. Process Learning

**What went wrong:**

1. **F1 honest-claim — pre-rebase SHA citations in §CDD Trace.** α applied `alpha/SKILL.md` §2.6 row 14 path (a) mid-cycle to remediate a worktree-inherited `gamma@cdd.cnos` identity leak (per α-closeout §Friction log item 1, same class as #379 F-item-1 and #370 F4). The `git rebase --exec 'git commit --amend --reset-author --no-edit'` rewrote every α commit's SHA. α had already authored `self-coherence.md` §CDD Trace step table + §ACs header reference citing the pre-rebase reflog SHAs (5 + 1 = 6 occurrences). After the rebase, those SHAs lived only in α's local reflog — not reachable from `origin/cycle/380`. β R1 caught via `review/SKILL.md` rule 3.13(a) reproducibility (a future reviewer cannot resolve the cited pointers). The cycle's loaded skill (`alpha/SKILL.md` §2.6 row 14 path (a)) prescribed the rebase mechanically but did not name SHA-citation invalidation as a downstream consequence. The skill was loaded; the gate fired; row 14's prescription was incomplete. **Class:** §9.1 trigger 4 — loaded skill failed to prevent a finding it covers.

2. **Pre-existing CI red on merge commit (I4 / notify) caps γ at C per §3.8 mechanically.** `Repo link validation (I4)` has been red since cycle #369 — lychee fails on file:// links in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` (cycle #369's own artifact). `notify` is downstream of I4 (sends a Telegram notification only on Build success). 3.78.0 documented this in RELEASE.md §"Known Issues" without applying the §3.8 cap (3.78.0's own regressions were what triggered the cap). 3.79.0's diff is non-destructive against the I4 failure (β verified identical failure shape on bare `origin/main` and merge tree per `beta-review.md` §"Branch CI state"); the cycle did not introduce or scope to fix the failure. The §3.8 CI-red cap clause is mechanical and does not distinguish "CI red introduced by this cycle" from "CI red purely pre-existing", so γ caps at C regardless of cause. **This is the rule operating as a forcing function**: every cycle until I4 is remediated pays the γ-axis cost. The right response is to fix I4, not to refine §3.8.

**What went right:**

1. **γ-scaffold failure-mode register served as α-side authoring checklist (third repo-level confirmation).** α-closeout P4 records that γ's pre-flagged six failure modes were held by α as generation constraints during authoring: #1 → `syscall.Exec` primitive (not `os/exec.Cmd`); #2 → `argv[1] != "exec"` test assertion (disproves `codex exec PROMPT` non-interactive subcommand drift); #4 + #6 → nonexistent-HUB_DIR pre-render-ordering tests; #5 → error message names both binary AND flag; #7 → bytes-equal default arm via Option A render-capture seam. Same pattern named in #367 O3, #379 P1, #370 O6, and now #380 P4 — four repo-level confirmations that γ scaffold's failure-mode catalogue compresses review rounds for the modes it names. The cycle's round count (2, at target) and finding density (1 B-severity, 0 D-severity, 0 in-AC-surface findings) reflect the discipline.

2. **Option A render-capture seam makes AC3 structurally true, not merely tested.** α's choice to capture into a `bytes.Buffer` only in the spawn arm and pass `inv.Stdout` straight through in the default arm is the minimum-surface implementation of "preserve default-path bytes-equal to 3.78.0". The bytes-equal property holds *by construction*, not by accident — even if the unit oracle (`TestActivate_DefaultNoFlag_BytesEqualToDirectRun`) were absent, the default arm cannot regress AC3 unless `activate.Run` itself is modified, and `git diff origin/main..origin/cycle/380 -- activate.go` is empty. α-closeout P2 generalizes the pattern: when a "preserve X exactly" invariant admits both a minimum-surface and a uniform-surface implementation, the minimum-surface implementation moves the invariant from oracle-checked to structurally-preserved.

3. **β pre-merge gate row 3 applied with explicit dual verification.** β ran the full Build workflow against the cycle/380 HEAD AND against bare `origin/main`, and named the identical failure on both as evidence that the merge is non-destructive with respect to the failing validators. The 3.78.0 PRA's β/SKILL.md row 3 enumeration patch (the four contract-validator classes) made the validator set legible; β's run exercised every class (Go build & test, SKILL.md frontmatter (I5), CDD ledger (I6), Package verification kata) and named the two pre-existing failures explicitly. The structural backstop from 3.78.0 paid off in 3.79.0: β did not have to invent which validators to run.

4. **Intra-doc repetition rule (α/SKILL.md §2.3) caught three β-blindspot SHA occurrences.** β R1 F1 enumerated 6 SHA citation sites in `self-coherence.md`. α grepped the file for every occurrence of the six pre-rebase SHAs and found 9 hits total — 3 beyond β's enumeration. α fixed all 9 in one commit (`91a2cad6`); R2 surfaced no F1-bis. The rule (derives from #266 F3-bis) extends naturally from numeric/named-value drift to honest-claim SHA citations as a peer class. This is the second cnos surface confirming the transfer.

**Skill patches:** YES — committed this release.

1. **`src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 row 14 — SHA-citation-after-rebase invariant.** Added a new auxiliary paragraph after row 14's path (a) prescription naming the SHA-citation invalidation as a downstream consequence and prescribing two coherent resolution paths: (i) run identity verification at session-start, before any SHA-bearing artifact is authored (preferred — the rebase has nothing to rewrite); (ii) re-stamp every SHA citation in `self-coherence.md` and other authored artifacts to the current-branch SHA immediately after the rebase, applying the §2.3 intra-doc repetition rule to catch every occurrence. Empirical anchor cited inline (#380 R1 F1). Class: structural backstop — α no longer has to derive the SHA-invalidation consequence per cycle. (Commit: this PRA.)

**Active skill re-evaluation:**

| Finding | Skill it should have caught | Underspecified or application gap? | Disposition |
|---------|-----------------------------|------------------------------------|-------------|
| F1 (§CDD Trace pre-rebase SHAs) | `alpha/SKILL.md` §2.6 row 14 path (a) | **Underspecified** — row 14 names the rebase mechanically but not the SHA-citation invalidation downstream | Patch (alpha/SKILL.md §2.6 row 14 auxiliary paragraph, this PRA) |
| F1 (also) | `alpha/SKILL.md` §2.3 intra-doc repetition rule | Application — the rule caught the additional 3 occurrences post-fact; it does not prevent the initial leak | Already-applied (α used it to find 9 occurrences); no patch |
| F1 (also) | `review/SKILL.md` rule 3.13(a) reproducibility | Already covers — β caught F1 at R1 via this rule cleanly | No patch needed (rule worked) |
| Pre-existing I4 / notify CI red | `release/SKILL.md` §3.8 CI-red cap clause | Working as designed — mechanical forcing function. Open question: should the clause distinguish introduced vs pre-existing failures? Filed as a γ process-gap observation but not patched this cycle (one cycle of data is not enough; need ≥2 cycles where the cap fires purely on inherited failures before refining) | No patch this cycle; observation recorded for next assessment |
| Worktree-config identity leak (#380 P5 + α F1 root cause) | `alpha/SKILL.md` §2.6 row 14 path (a) + `beta/SKILL.md` row 1 worktree user.email | **Skill gap is repo-level, not role-specific** — #379 F-item-1, #370 F4, #380 F1 all carry the same class. Issue #373 (preventive `--worktree` identity write across role skills) is the converged-but-unimplemented design at growing lag | Issue #373 escalates from growing to **escalating** lag per §2; γ recommends prioritizing #373 in a near-term cycle as the structural fix |

**CDD improvement disposition:** Patch landed (alpha/SKILL.md §2.6 row 14 SHA-citation-after-rebase clause). One §9.1 trigger fired (loaded-skill miss); the patch closes the loop at the role-skill layer. The pre-existing-CI-red γ-axis cap is the rule operating as designed; remediation is via I4 fix as a project-level next-MCA, not via §3.8 refinement.

### 4. Review Quality

**Cycles this release:** 1 (#380)
**Avg review rounds:** 2.0 (target: ≤2 code)
**Superseded cycles:** 0 (target: 0)

**Per-cycle round counts:**

| Cycle | Issue | Mode | Rounds | Binding findings (R1) | Notes |
|-------|-------|------|--------|----------------------|-------|
| #380 | `cn activate --claude` / `--codex` flags | design-and-build | 2 | 1 (B-severity honest-claim, F1) | R1 RC on F1 (pre-rebase SHAs in §CDD Trace); R2 APPROVED after one self-coherence-only fix commit (`91a2cad6`) + fix-round-1 append (`d838e7e9`). No production-code delta R1→R2. |

**Per-cycle dispatch telemetry** (second row of accumulating data; first row in 3.78.0 PRA):

| Cycle | `dispatch_seconds_budget` | `dispatch_seconds_actual` | `commit_count_at_termination` |
|-------|--------------------------|--------------------------|-------------------------------|
| #380 (γ scaffold) | not recorded | not recorded | 1 (γ-scaffold.md committed; clean exit) |
| #380 (α R1) | not recorded | not recorded | 11 (incremental self-coherence sections + 5 implementation steps + review-readiness signal) |
| #380 (β R1) | not recorded | not recorded | 1 (R1 RC verdict commit `2dac48ba`) |
| #380 (α fix-round 1) | not recorded | not recorded | 2 (`91a2cad6` SHA rewrite + `d838e7e9` fix-round-1 append) |
| #380 (β R2) | not recorded | not recorded | 4 (R2 APPROVED commit `4a7cda79` + merge `770ea1b4` + β-closeout `b53ba6a4`) |
| #380 (α close-out) | not recorded | not recorded | 1 (alpha-closeout `f15c27cc`) |

The γ scaffold (`gamma-scaffold.md` §Dispatch configuration) declared §5.1 multi-session dispatch but did not name a per-role timeout budget. No agent SIGTERM'd; all exited cleanly. Budget-vs-actual remains non-binding. Future cycles should record both per `CDD.md` §1.6c so the heuristic constants can be tightened.

**Finding-class breakdown** (across cycles in this release; in-cycle + post-merge combined):

| Class | Definition | Count |
|---|---|---|
| **mechanical** | Caught by grep/diff/script | 0 |
| **wiring** | "X is wired into Y" but isn't (review/SKILL.md 3.13c) | 0 |
| **honest-claim** | Doc claims something code/data doesn't back (review/SKILL.md 3.13) | 1 (F1) |
| **judgment** | Design/coherence assessment | 0 |
| **contract** | Work contract incoherent | 0 |

**Mechanical ratio:** 0% (below 20%-with-≥10-findings threshold; no process issue filing required on the mechanical-ratio axis).

**Honest-claim ratio:** 100% (1/1). The denominator is too small to be meaningful (target <30% applies when N is large enough to be a ratio rather than a coin flip). Note for cumulative tracking: 3.78.0 had 0 honest-claim findings (2 wiring-class post-merge); 3.79.0 has 1 honest-claim in-cycle (caught by β R1, fixed at R2). Trend across the two-release window: 1 honest-claim in 3 total findings = 33% cumulative, slightly over the <30% guideline. One more cycle of data will determine whether this is a trend or noise.

**Action:** none on the mechanical-ratio axis. The honest-claim finding is addressed by the alpha/SKILL.md §2.6 row 14 patch (above, §3).

### 4a. CDD Self-Coherence

- **CDD α (artifact integrity):** 3/4 — All required cycle artifacts present on `cycle/380` and `main`: `gamma-scaffold.md`, `self-coherence.md` (7 sections + fix-round-1 append), `beta-review.md` (R1 + R2), `alpha-closeout.md`, `beta-closeout.md`. F1 was an artifact-integrity violation at R1 (pre-rebase SHAs unreachable) but α closed it cleanly at R2 with intra-doc rule application. The −1 is for F1's existence at R1, not for the artifact set's completeness.
- **CDD β (surface agreement):** 4/4 — Canonical doc (`CDD.md`), executable skills (`alpha/`, `beta/`, `gamma/`, `release/`, `review/`, `operator/`), `.cdd/unreleased/380/` cycle artifacts, CHANGELOG, RELEASE.md, and this PRA all agree on cycle scope, deliverables, and outcome. The alpha/SKILL.md row 14 SHA-citation patch landing in this PRA is itself a β-axis tightening (closes the under-specification surfaced by F1).
- **CDD γ (cycle economics):** 2/4 — 2 review rounds (at target, not exceeded). 0 superseded cycles. CI red on merge commit (pre-existing infrastructure; §3.8 cap mechanical). Mechanical-ratio under threshold. Step-13a patch landed in this PRA commit (not deferred). The CI-red post-merge condition is the dominant signal pulling this axis below 3, same shape as 3.78.0 — but for 3.79.0 the red is inherited infrastructure rather than introduced regressions, so the *coordination quality* signal is stronger than the rubric grade reflects.
- **Weakest axis:** γ (CI-red cap, mechanical).
- **Action:** Patch landed at the role-skill layer (alpha/SKILL.md §2.6 row 14). The CI-red cap is addressed at the project-debt layer (deferred MCA for I4 link validation remediation), not at the rule layer.

### 4b. Cycle Iteration

- **Triggered by** (per CDD.md §9.1 thresholds):
  - review rounds > 2: **NO** (actual: 2; at target, not exceeded by the strict reading of §9.1)
  - mechanical ratio > 20% with ≥10 findings: NO (1 finding total; below threshold count)
  - avoidable tooling/environmental failure: NO (the F1 cause is a skill gap, not tooling failure; the pre-existing I4 red is infrastructure debt, not a cycle-blocking environmental failure)
  - **CI red on merge commit (post-merge): YES** — `Repo link validation (I4)` and `notify` failed on merge commit `770ea1b4` (pre-existing from #369; not introduced by #380; documented in 3.78.0 RELEASE.md Known Issues; β verified non-destructive). Triggers `release/SKILL.md` §3.8 CI-red cap clause (γ axis capped at C).
  - **loaded skill failed to prevent a finding: YES** — `alpha/SKILL.md` §2.6 row 14 path (a) was loaded and applied by α to remediate the worktree-inherited identity leak; the rebase invalidated 9 SHA citation sites in `self-coherence.md` that α had already authored; β R1 caught via rule 3.13(a) reproducibility. Row 14 named the rebase mechanically but did not name SHA-citation invalidation as a downstream consequence.

- **Root cause:** Two distinct causes, one trigger fired each:
  1. **F1 — alpha/SKILL.md §2.6 row 14 under-specification.** Path (a) prescription (`git rebase --exec 'git commit --amend --reset-author --no-edit' && git push --force-with-lease`) rewrote every α commit's SHA but did not name the SHA-citation-in-artifacts downstream consequence. α had already authored §CDD Trace + §ACs with pre-rebase reflog SHAs; the rebase invalidated all references.
  2. **CI-red — pre-existing I4 infrastructure debt.** Cycle #369's `self-coherence.md` (archived at `.cdd/releases/docs/2026-05-17/369/self-coherence.md`) carries file:// links lychee cannot resolve. Every cycle merged into main since #369 has inherited the red. #380 did not scope to fix it; the §3.8 cap applies mechanically.

- **Disposition:** **Patch landed now** for cause 1; **deferred MCA** for cause 2.
  1. (F1) alpha/SKILL.md §2.6 row 14 SHA-citation invariant clause — patch lands in this PRA commit.
  2. (I4) — file a follow-on issue to remediate `.cdd/releases/docs/2026-05-17/369/self-coherence.md` file:// references (per next MCA list below).

- **Evidence:**
  - Merge-commit CI run: https://github.com/usurobor/cnos/actions/runs/26102342968 (I4 + notify failure; all other jobs green).
  - F1 finding: `.cdd/unreleased/380/beta-review.md` §"Findings" R1 F1; closure at `.cdd/unreleased/380/self-coherence.md` §"Fix-round 1" + α commit `91a2cad6`.
  - Patch (cause 1): this PRA commit (`src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 auxiliary paragraph after row 14).
  - MCA (cause 2): GitHub issue to be filed naming I4 link validation remediation as a structural fix (clears the §3.8 cap for future cycles); first AC: `lychee` returns "0 errors" against the `.cdd/releases/docs/2026-05-17/369/` directory.

### 5. Production Verification

**Scenario:** An operator running `cn activate --claude cn-sigma` on a workstation with the `claude` CLI installed lands in a live interactive `claude` REPL with the Sigma body already activated — naming its identity (Sigma), its operator (usurobor / Axiom), and its current orientation — within one command. The same `cn activate --codex cn-sigma` shape works for codex.

**Before this release:** 3.78.0 shipped the renderer (`cn activate HUB_DIR` writes the rendered prompt to stdout). The operator had to (a) capture the stdout, (b) start a fresh `claude` / `codex` REPL, and (c) paste the prompt manually. The pipe form `cn activate cn-sigma | claude` was non-interactive (claude `-p`-equivalent prints output and exits). The bridge from "prompt rendered" to "body operating interactively at this hub" was missing — exactly the gap #380 §Problem named.

**After this release:** `cn activate --claude HUB_DIR` and `cn activate --codex HUB_DIR` render the prompt via the existing 3.78.0 path (`activate.Run`), then replace the cn process image with `claude "$prompt"` or `codex "$prompt"` via `syscall.Exec` (process replacement preserves the TTY cleanly). The operator lands in a live REPL with the body pre-activated. Default `cn activate HUB_DIR` (no flag) is unchanged from 3.78.0 by construction. Pre-render mutual exclusion and missing-binary errors fire before any render call, proven by the "nonexistent HUB_DIR" oracle shape.

**How to verify:**
1. **Bytes-equal default path (cycle-level):** `go test -count=1 -run TestActivate_DefaultNoFlag_BytesEqualToDirectRun ./internal/cli/...` and `git diff origin/main..origin/cycle/380 -- src/go/internal/activate/activate.go` empty. Already PASS on cycle HEAD and merge tree per β R1 + R2 verification.
2. **Argv shape (cycle-level):** `go test -count=1 -run TestSpawnWith ./internal/activate/...` — both `TestSpawnWith_ClaudeArgvShape` and `TestSpawnWith_CodexArgvShape` pass; the codex test explicitly disproves `argv[1] == "exec"` (non-interactive subcommand drift).
3. **Pre-render ordering (cycle-level):** `go test -count=1 -run TestActivate_(MutualExclusion|MissingBinary)_FiresBeforeRender ./internal/cli/...` — stderr contains the pre-render diagnostic AND does NOT contain `Hub path not found` AND does NOT contain `Generating activation prompt` for both mutual-exclusion and missing-binary cases.
4. **End-to-end interactive bootstrap (manual dry-run):** open a fresh terminal on a workstation with `claude` installed, run `cn activate --claude cn-sigma`, observe whether the resulting `claude` REPL has Sigma already activated (names identity + operator + current orientation) without further operator prompts. **Deferred** — this PRA does not have an interactive workstation session available for the dry-run; commit the deferral to the next σ session checklist.

**Result:** PASS (1–3) / deferred (4). The cycle's structural change (default-path bytes-equal by construction; argv shape verified by unit tests; pre-render ordering proven by absence-of-render-diagnostic) is verified by the Go tests that landed in `spawn_test.go` and `cmd_activate_test.go`. The cycle's intended operator-visible change (one command lands in a live REPL) is testable but not exercised in this PRA session; named explicitly as deferred verification.

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 11 Observe | `gh run view 26102342968` (CI on merge SHA) + `git log 770ea1b4..main` + cycle-dir close-out reads | post-release, gamma | Merge commit landed; CI red but exclusively on pre-existing I4 / notify infrastructure (not introduced by #380); β verified non-destructive; one in-cycle finding (F1) cleanly closed. |
| 12 Assess | `docs/gamma/cdd/3.79.0/POST-RELEASE-ASSESSMENT.md` (this file) | post-release | Assessment completed — one skill patch landed (alpha/SKILL.md §2.6 row 14 SHA-citation invariant); one deferred MCA filed (I4 link validation remediation). |
| 13 Close | `.cdd/unreleased/380/gamma-closeout.md` (this PRA commit) + `.cdd/unreleased/380/cdd-iteration.md` (this PRA commit) + cycle dir move to `.cdd/releases/3.79.0/380/` + `RELEASE.md` + CHANGELOG 3.79.0 row | post-release, release, gamma | Cycle closed in this PRA commit; δ tag/release follows. |
| 13a Skill/spec patch | `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 row 14 — SHA-citation-after-rebase invariant | post-release, alpha | Patch lands in this PRA commit (not deferred). |

### 6a. Invariants Check

cnos has no single repo-level INVARIANTS.md document; protocol invariants live in `CDD.md` and per-skill kata sections. The constraints touched by this cycle:

| Constraint | Touched? | Status |
|---|---|---|
| `cn activate HUB_DIR` default stdout behavior preserved bytes-equal to 3.78.0 | Yes (AC3) — preserved by construction via Option A render-capture seam | preserved |
| `cn activate --claude` / `--codex` mutually exclusive (AC4) | Yes — added in this cycle; pre-render check | new invariant; preserved |
| Pre-render ordering for mutex + missing-binary errors (AC4 + AC5) | Yes — added in this cycle; proven by nonexistent-HUB_DIR oracle | new invariant; preserved |
| `syscall.Exec` over `os/exec.Cmd` for TTY preservation (active design constraint) | Yes — `spawn.go:24` `defaultExecFn = syscall.Exec`; `os/exec.Cmd` not imported in `spawn.go` | preserved |
| Bare-positional argv (`[binary, prompt]`, not `[binary, -p, prompt]` or `[binary, exec, prompt]`) | Yes — `spawn.go:68` builds `argv := []string{binary, prompt}`; `TestSpawnWith_CodexArgvShape` explicitly asserts `argv[1] != "exec"` | preserved |
| `src/go/internal/activate/activate.go` byte-identical to 3.78.0 (Option A) | Yes — `git diff origin/main..origin/cycle/380 -- activate.go` empty | preserved |
| `src/packages/cnos.core/skills/agent/activate/SKILL.md` byte-identical to 3.78.0 (non-goal) | Yes — skill is body-agnostic; cycle is `cn`-side concern | preserved |
| β/SKILL.md §pre-merge gate row 3 enumeration (from 3.78.0) | Yes — β applied the enumerated validator set; structural backstop paid off | tightened-and-exercised |
| `alpha/SKILL.md` §2.6 row 14 (path (a) rebase) | Yes — extended with SHA-citation-after-rebase invariant clause in this PRA | tightened (this PRA) |

### 7. Next Move

**Next MCA:** **I4 link validation remediation** (file:// references in `.cdd/releases/docs/2026-05-17/369/self-coherence.md`).

This is the highest-leverage near-term move because:
1. It clears the §3.8 CI-red cap that has been in effect since #379 merged (and likely earlier for prior cycles whose own scope didn't introduce regressions); every cycle until this is fixed pays the γ-axis cost.
2. The fix is mechanical and bounded — replace file:// references with relative or `https://` references; expected to be a single small-change cycle.
3. It is independent of the larger Hub README router adoption work (which has been at growing lag since 3.78.0); the two can proceed in parallel.

The Hub README router adoption (cn-sigma first, cn-pi second) remains the highest-leverage *cross-repo* MCA and is carried forward unchanged from 3.78.0's next-move commitment. The I4 fix is the next *cnos-internal* MCA. γ's recommendation: prioritize I4 (small, fast, clears the cap) before resuming the cn-sigma adoption (cross-repo, requires coordination).

**Owner:** γ at cnos files the I4 issue; α executes the fix in a single small-change cycle.
**Branch:** pending creation as `cycle/{N+1}` once the issue is filed.
**First AC:** `lychee` returns "0 errors" against `.cdd/releases/docs/2026-05-17/369/`; CI Build workflow's `Repo link validation (I4)` job exits green on the resulting merge commit; the `notify` job (downstream of Build success) exits green.
**MCI frozen until shipped?** Yes — see §2 (6 issues at growing lag; freeze in effect from 3.78.0 PRA, continuing).
**Rationale:** The §3.8 CI-red cap clause is a forcing function for infrastructure debt. The right response is to clear the debt, not to refine the rule. A single small-change cycle on I4 unblocks the γ-axis grade for every subsequent cycle.

**Closure evidence (CDD §10):**
- Immediate outputs executed: yes
  - alpha/SKILL.md §2.6 row 14 SHA-citation-after-rebase invariant clause (this PRA commit)
- Deferred outputs committed: yes
  - I4 link validation remediation — to file as a new cnos issue; first AC named above
  - Hub README router adoption (cn-sigma) — carried forward from 3.78.0; γ at cnos files cross-repo issue at `usurobor/cn-sigma`; cn-sigma's δ executes the README patch
  - Hub README router adoption (cn-pi) — carried forward from 3.78.0
  - `//go:embed` durable fix for the renderer fallback duplication (#379 α §Debt 3) — P3, future cycle
  - `cn doctor` enforcement of activation invariants (#379 α §Debt 4) — P3, future cycle
  - `cnos.xyz/activate/<hub>` rendering service (#379 §Scope) — P3, requires hosting decision
  - `--cursor` / `--aider` / other AI CLI flags — held by explicit non-goal; revisit when friction emerges
  - `--auto` / `$CN_DEFAULT_BODY` — held by explicit non-goal; revisit when friction emerges
  - #373 (preventive `--worktree` identity write) — escalating from growing per §2; γ recommends prioritizing in a near-term cycle alongside or after I4 fix

**Immediate fixes** (executed in this session):
- alpha/SKILL.md §2.6 row 14 SHA-citation-after-rebase invariant clause (named above)

### 8. Hub Memory

- **Daily reflection:** deferred — cnos is the protocol repo, not a hub. Hub memory in cnos lives in per-hub repos (cn-sigma, cn-pi). The σ session that lands the cn-sigma README router adoption (the cross-repo next MCA carried forward from 3.78.0) will write the daily reflection there; the σ session that lands the I4 fix (the next cnos-internal MCA) will write the daily reflection if a σ session is involved, otherwise the cycle's α writes their own observations into `alpha-closeout.md`.
- **Adhoc thread(s) updated:** deferred — same reason. The cnos-internal "rebase-with-reset-author SHA-citation hazard" is now codified in `alpha/SKILL.md` §2.6 row 14; the protocol-layer durable memory is the skill patch itself. No separate adhoc thread surface is needed at cnos for protocol learnings.
- **Cnos-side memory equivalent:** this PRA + `gamma-closeout.md` + `cdd-iteration.md` for #380 are the durable memory at the protocol layer; no separate cnos-side daily reflection convention exists.
