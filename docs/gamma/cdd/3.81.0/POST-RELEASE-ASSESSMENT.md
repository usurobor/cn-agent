<!-- sections: [Header, §1 Coherence, §2 Lag, §3 Process, §4 Review Quality, §4a Self-Coherence, §4b Cycle Iteration, §5 Production Verification, §6 CDD Closeout, §7 Next Move, §8 Hub Memory] -->
<!-- completed: [Header, §1 Coherence, §2 Lag, §3 Process, §4 Review Quality, §4a Self-Coherence, §4b Cycle Iteration, §5 Production Verification, §6 CDD Closeout, §7 Next Move, §8 Hub Memory] -->

## Post-Release Assessment — 3.81.0

**CI status on merge SHA:** red — https://github.com/usurobor/cnos/actions/runs/26191235776 (merge commit `45dbcc47`; pre-existing failure pattern: all recent main commits show `conclusion=failure`; not cycle-introduced; same class as 3.79.0. Per `operator/SKILL.md` §3.4 step 4, release proceeds. γ axis caps at C per `release/SKILL.md` §3.8 CI-red cap clause.)

---

### 1. Coherence Measurement

- **Baseline:** 3.79.0 — α B+, β A, γ C, C_Σ B− (3.80.0 had no PRA — α-direct small-change override; 3.79.0 is last independently-scored baseline)
- **This release:** 3.81.0 — α B+, β A, γ C, C_Σ B−
- **Delta:**
  - **α (PATTERN): Held at B+.** Nine ACs met cleanly. The absorb-into-cap structure is architecturally sound: operational rules now live in one scan (cap §4–§6) instead of six files. Non-absorbed boundary notes (root-finding sequence stays in mca; scope/specificity rules + MCI vs MCA table stay in mci; Coherence Modes + Anti-Patterns tables stay in coherent) are concrete and named in each section's closing note. Cap section renumbering (old §4–§7 → new §8–§10) is a structural consequence of adding §4–§6, handled cleanly with no external cross-reference rot (confirmed by grep). Provisional close-out (acceptable per §1.4 α step 10 fallback; no re-dispatch occurred) is carried as known debt. Minor: no version directory created for this cycle — acceptable for a skill-doctrine cycle with no frozen snapshot required.
  - **β (RELATION): Held at A.** Zero findings at any severity. All 9 ACs verified with grep-reproducible evidence. Architecture check (7 questions) passed. Pre-merge gate 4/4 rows passed. Contract integrity 11/11 rows yes. One non-finding observation logged appropriately (no integration test for `{cap,clp}` rendered string — correct by inspection, not raised as formal finding per Rule 3.5/3.6).
  - **γ (EXIT/PROCESS): Held at C (CI-red cap).** One review round (≤2 target ✓). Zero superseded cycles ✓. Mechanical ratio N/A (0 findings). Cycle economics clean. CI red on merge commit is pre-existing (same I4 link-validation + notify pattern as 3.79.0; β verified identical failure shape on bare `origin/main` and merge tree). CI-red cap is mechanical per §3.8; γ caps at C regardless of cause. The right response remains fixing I4, not refining §3.8.
  - **C_Σ:** B− (geometric mean of B+ · A · C ≈ B−; same as 3.79.0)

- **Coherence contract closed?** Yes. The cycle's named gap — `agent/activate/SKILL.md §2.1 step 2` loading 6 CA skills when 4 were redundant or deprecated — is closed. The activation soul is now: KERNEL + cap + clp (3 surfaces). All operational rules previously scattered across mca/mci/coherent are absorbed into cap §4–§6 with concrete non-absorbed boundary notes. agent-ops is marked as describing the OCaml-era daemon contract (not the current Go activation path). `cnos.core/AGENTS.md` is deleted. The renderer and kata fixture reflect the 2-skill path. Hub-side downstream cleanup (cn-sigma RULES.md, spec/AGENTS.md, spec/TOOLS.md, spec/HEARTBEAT.md) remains deferred per scope — named in issue §Hub-side downstream as separate per-hub work.

---

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| #385 | Streamline activation soul: 6→2 CA skills; deprecate AGENTS.md | feature | issue body (CLP r2) | shipped (this release) | none |
| #383 | `cn activate`: collapse rendered prompt to Read-first only | feature | issue body | not started (unblocked by #385) | growing |
| #373 | Preventive `--worktree` identity write across all role skills | process | #373 issue body | not started | **escalating** (4th cycle confirming class: β Row 1 + γ startup fix in #385; #380 F1; #379 F-item-1; #370 F4) |
| (I4 debt) | I4 link validation — stale file:// refs in `.cdd/releases/docs/2026-05-17/369/self-coherence.md` | process | scope known | not started | **forcing function** (γ-axis cap every release until fixed) |
| (hub) | Hub-side downstream cleanup — cn-sigma RULES.md, spec/AGENTS.md, spec/TOOLS.md, spec/HEARTBEAT.md | feature | #385 §Hub-side downstream | not started | growing |
| (carried) | Hub README router adoption in cn-sigma | feature | `agent/activate/SKILL.md` §2.3 template | not started | growing |
| (carried) | Hub README router adoption in cn-pi | feature | `agent/activate/SKILL.md` §2.3 template | not started | growing |
| (carried) | `//go:embed` activate skill in cn binary | feature | named in #379 α §Debt 3 | not started | growing |
| (carried) | `cn doctor` enforcement of activation invariants | feature | named in #379 α §Debt 4 | not started | growing |
| (carried) | `cnos.xyz/activate/<hub>` rendering service | feature | named in #379 §Scope | not started | growing |
| (carried) | `--cursor` / `--aider` / other AI CLI flags | feature | named in #380 §Non-goals | not started | low |
| (carried) | `--auto` body-detection + `$CN_DEFAULT_BODY` env-var default | feature | named in #380 §Deferred | not started | low |

**MCI/MCA balance:** **Freeze MCI** — 6+ items at growing lag; #373 escalating (4 cycles); I4 forcing function active.

**Rationale:** The freeze declared at 3.78.0 remains active. #385 shipped one MCA (soul streamlining) but did not reduce the hub-adoption lag or address the structural worktree identity gap (#373). The fastest unfreeze move is now #373 (escalating, structural fix, directly blocks productivity on every cycle start) followed by I4 remediation. Continue holding new design commitments.

---

### 3. Process Learning

**What went wrong:**

1. **Worktree identity inheritance — γ session startup fix.** γ's session started with `config.worktree` carrying β's identity (`beta@cdd.cnos`) from the prior role session. γ had to apply the worktree override at session startup before committing. Same class as #380 F1 (β worktree leak from α session), #379 F-item-1, #370 F4. This is the 4th cycle-level confirmation of the worktree identity leak class. The structural fix (#373 — preventive `--worktree` write across all role skills when `extensions.worktreeConfig=true`) is filed and escalating but not yet implemented. **No new skill patch this cycle** — the class is already fully documented in #373's issue body and in three prior PRAs. Additional cycle-level patching of individual role skills without the structural fix is churn.

2. **Provisional α close-out — re-dispatch not executed.** Per CDD §1.4 α step 10, α may write a provisional close-out at review-readiness time, explicitly marked `[provisional — pending β outcome]`. The provisional form was used for cycle #385. This is acceptable per protocol but represents a gap in the re-dispatch mechanism: δ did not re-dispatch α after β merge. α's close-out remains provisional. γ supplements with PRA observations per the fallback contract. **No skill patch** — the fallback is correctly specified; the gap is a dispatch mechanics gap (δ didn't re-dispatch), not a skill gap.

3. **Missing 3.80.0 PRA.** Cycle #381 (small-change α-direct override) shipped as 3.80.0 with no post-release assessment. Per CDD §1.2, PRA is required for every release including small-change. This PRA debt from 3.80.0 should be resolved: either write a retroactive PRA for 3.80.0, or formally record it as waived with operator authorization. **Observation only** — not patching protocol here; the protocol already requires PRA. Gap is operator process, not skill spec.

**What went right:**

1. **CLP r2 convergence at issue creation reduced issue-level ambiguity to zero.** Nine ACs mapped directly to implementable changes. α reported no design discovery needed. β found zero findings. This is the first cycle in recent history achieving 0 findings at R1 on a 9-AC code cycle. The CLP-prior-to-dispatch discipline is demonstrably effective for complex skill-authoring changes.

2. **Non-absorbed boundary notes in cap §4/§5/§6 make the absorption boundary concrete.** Each new section's closing note names exactly what stays in the standalone file and why. This pattern (first confirmed in `activate/SKILL.md §2.4` disambiguation) is now instantiated in three cap sections, making the boundary a maintained artifact rather than a tacit assumption.

3. **`activation_status:` frontmatter field as a new convention.** The four standalone files (mca, mci, coherent, agent-ops) each received `activation_status: on-demand — not activation-loaded since cycle/385`. α identified this as potentially formalizable as a standard frontmatter field for any skill that transitions from activation-loaded to on-demand. This is a useful new convention; it could be codified in `skill/SKILL.md` (the skill-authoring skill) as a standard field for downgraded activation skills.

**Skill patches:** None committed this cycle.

**Active skill re-evaluation:**

| Finding | Skill it should have caught | Underspecified or application gap? | Disposition |
|---------|-----------------------------|------------------------------------|-------------|
| Worktree identity (observation, not β finding) | `gamma/SKILL.md` (no explicit worktree identity check at γ startup) | **Skill gap** — role skills don't prescribe `--worktree` identity write; structural fix is #373 | No patch this cycle — structural fix (#373) is the right vehicle; patchwork per-role fixes are churn until #373 ships |
| Provisional α close-out | n/a — protocol correctly specifies the fallback | Application gap — δ didn't re-dispatch α | No patch |
| Missing 3.80.0 PRA | `CDD.md §1.2` — PRA required for all releases | Application gap — operator skipped for α-direct cycle | No patch (protocol already requires it) |
| β non-finding: no integration test for `{cap,clp}` path | `review/SKILL.md` Rule 3.5/3.6 correctly applied by β | Application — β noted and deferred correctly | No patch; follow-up: augment `TestReadFirstSection_OrderedSigma` with `strings.Contains(out, "{cap,clp}")` in a future small-change cycle |

**CDD improvement disposition:** No patches committed this cycle. No formal CDD §9.1 triggers fired. All identified gaps are either already filed as issues (#373) or are application gaps with adequate existing spec. One new convention observed (`activation_status:` field) — not patching skill/SKILL.md this session; the observation is worth a small-change cycle to codify once the convention is confirmed stable.

---

### 4. Review Quality

**Cycles this release:** 1 (#385)
**Avg review rounds:** 1.0 (target: ≤2 code — ✓, at target floor)
**Superseded cycles:** 0 (target: 0 ✓)

**Per-cycle round counts:**

| Cycle | Issue | Mode | Rounds | Binding findings (R1) | Notes |
|-------|-------|------|--------|----------------------|-------|
| #385 | Streamline activation soul: 6→2 CA skills; deprecate AGENTS.md | design-and-build | 1 | 0 | Zero findings at any severity. 9 ACs, 14 files, 527 ins/200 del. APPROVED R1. |

**Per-cycle dispatch telemetry:**

| Cycle | `dispatch_seconds_budget` | `dispatch_seconds_actual` | `commit_count_at_termination` |
|-------|--------------------------|--------------------------|-------------------------------|
| #385 | not recorded | not recorded | 7 (α commits landed before review-readiness) |

(Budget and actual not recorded for this cycle; operator-dispatched without explicit timeout. Target: code cycles ≥ `max(400s, 180×9) = 1620s`. Data point needed for heuristic calibration.)

**Finding-class breakdown:**

| Class | Definition | Count |
|---|---|---|
| **mechanical** | Caught by grep/diff/script | 0 |
| **wiring** | "X is wired into Y" but isn't | 0 |
| **honest-claim** | Doc claims something code/data doesn't back | 0 |
| **judgment** | Design/coherence assessment | 0 |
| **contract** | Work contract incoherent | 0 |

**Mechanical ratio:** N/A (0 findings total — ratio is undefined)
**Honest-claim ratio:** N/A (0 findings)
**Action:** none

### 4a. CDD Self-Coherence

- **CDD α:** 3/4 — Required artifacts present (self-coherence, gamma-scaffold, beta-review, provisional alpha-closeout). All 9 ACs met with grep-reproducible evidence. No bootstrap/version directory — acceptable for skill-doctrine cycle with no frozen snapshot required. Provisional close-out is known debt (re-dispatch not executed). Minor deduction for provisional close-out form.
- **CDD β:** 4/4 — Contract integrity 11/11 rows. AC coverage 9/9 with reproducible evidence. Architecture check 7/7. Pre-merge gate 4/4. Zero findings. Non-finding observation logged appropriately. Surface agreement complete.
- **CDD γ:** 2/4 — One review round (clean ✓). Zero superseded cycles (✓). Cycle economics clean. CI red on merge commit (pre-existing, mechanical γ-axis cap per §3.8). Missing 3.80.0 PRA from prior release is documented process debt. Worktree identity issue required manual fix at γ startup. Deductions for CI cap (mechanical) and process continuity gaps.
- **Weakest axis:** γ
- **Action:** γ axis — fix I4 (next MCA candidate), implement #373 (structural worktree fix). No skill patch needed this cycle.

### 4b. Cycle Iteration

- **Triggered by:** none
  - review rounds > 2: No (1 round)
  - mechanical ratio > 20% with ≥ 10 findings: No (0 findings)
  - avoidable tooling/environmental failure: No (CI red is pre-existing infrastructure debt, not cycle-introduced)
  - loaded skill failed to prevent a finding: No (zero findings)

- **Independent γ process-gap check (CDD §13, step 13):** No new process gates identified that an immediate patch would close. Three friction points were observed (worktree identity, provisional close-out, missing 3.80.0 PRA) but all map to either an already-filed issue (#373) or an application gap with existing protocol coverage. No recurring friction that a skill patch would fix this cycle.

- **Root cause:** N/A — no trigger fired.
- **Disposition:** No patch — no §9.1 trigger fired; independent check found no actionable gap requiring immediate patch. #373 remains the structural fix for the worktree class.
- **Evidence:** zero β findings, 1 round, clean economics.

---

### 5. Production Verification

**Scenario:** `cn activate` renders a Read-first prompt that lists `cnos.core/skills/agent/{cap,clp}/SKILL.md` (2-skill path) — not the 6-skill path.

**Before this release:** `cn activate HUB_DIR` rendered a prompt listing `cnos.core/skills/agent/{cap,clp,mca,mci,coherent,agent-ops}/SKILL.md` (6 skills).

**After this release:** `cn activate HUB_DIR` renders a prompt listing `cnos.core/skills/agent/{cap,clp}/SKILL.md` (2 skills).

**How to verify:**
```bash
# Build
cd /root/cnos && go build -o /tmp/cn ./src/go/cmd/cn/
# Run activation (requires a hub dir with cnos vendored)
/tmp/cn activate /root/cn-sigma 2>/dev/null | grep "cap,clp"
# Or: inspect the renderer directly
grep "ca-skills\|cap,clp" src/go/internal/activate/activate.go
```

**Result:** Deferred — `cn activate` requires a hub dir and vendored skills; environment not set up for end-to-end invocation in this session. Code-level verification confirmed by β: `activate.go:506` returns `cnos.core/skills/agent/{cap,clp}/SKILL.md (CA skill set)`. Go unit tests pass (`go test ./internal/activate/...`). R5-activate kata passes 27/27 (parser-level verification). The `{cap,clp}` path is present in the compiled renderer.

---

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 11 Observe | CHANGELOG TSC, lag table, β/α close-outs, beta-review.md | post-release | CI red pre-existing; 0 findings; 1 round; 9 ACs met |
| 12 Assess | `docs/gamma/cdd/3.81.0/POST-RELEASE-ASSESSMENT.md` | post-release | Assessment completed; no immediate skill patches; #373 escalated |
| 13 Close | gamma-closeout.md, RELEASE.md, CHANGELOG 3.81.0 row, cdd-iteration.md | post-release | Immediate outputs executed in this session; deferred outputs committed (#373, hub cleanup) |

### 6a. Invariants Check

| Constraint | Touched? | Status |
|---|---|---|
| KERNEL is always loaded first at activation | yes | preserved — activate/SKILL.md §2.1 step 1 unchanged |
| CA skills are activation-loaded (cap, clp) | yes | tightened — from 6 to 2; mca/mci/coherent/agent-ops are now on-demand |
| Hub-side content (PERSONA, OPERATOR) stays at hub | no | N/A |
| Renderer `ca-skills` token reflects skill list | yes | preserved — activate.go:506 updated to `{cap,clp}` |
| pkgbuild root files include all load-bearing roots | yes | tightened — AGENTS.md removed from build list |

---

### 7. Next Move

**Next MCA:** #373 — Preventive `--worktree` identity write across all role skills when `extensions.worktreeConfig=true`
**Owner:** α/δ
**Branch:** cycle/373 (to be created by γ at dispatch)
**First AC:** All role skill files (alpha/SKILL.md, beta/SKILL.md, gamma/SKILL.md) prescribe `git config --worktree user.name {role} && git config --worktree user.email {role}@{project}.cdd.cnos` at the top of the identity step — before any commit or SHA-bearing artifact authoring.
**MCI frozen until shipped?** Yes — freeze continues. No new substantial design until growing-lag queue is reduced.
**Rationale:** #373 is the structural fix for a class that has now fired 4 times across 4 cycles (#370, #379, #380, #385). Each cycle pays a manual workaround. The workaround is well-understood but the structural fix has not been implemented. This is the highest-leverage near-term MCA: it unblocks clean role startup on every future cycle without manual intervention. I4 fix (link validation) is the follow-up immediately after #373 — it removes the CI cap that has suppressed γ to C for 3 consecutive releases.

**Closure evidence (CDD §10):**
- Immediate outputs executed: yes
  - PRA written: `docs/gamma/cdd/3.81.0/POST-RELEASE-ASSESSMENT.md` (this artifact)
  - gamma-closeout.md: `.cdd/unreleased/385/gamma-closeout.md`
  - RELEASE.md: `/root/cnos/RELEASE.md`
  - CHANGELOG 3.81.0 entry: `CHANGELOG.md`
  - cdd-iteration.md: `.cdd/unreleased/385/cdd-iteration.md`
  - Cycle dir move: `.cdd/unreleased/385/` → `.cdd/releases/3.81.0/385/` (at release time)
- Deferred outputs committed: yes
  - #373 (worktree identity structural fix) — open issue, escalating lag, named as next MCA
  - Hub-side downstream cleanup (cn-sigma) — deferred per #385 scope, tracked in issue body
  - Missing 3.80.0 PRA — operator process debt, not a cycle artifact

**Immediate fixes (executed in this session):**
- None beyond PRA + close-out artifacts (no skill patches needed)

---

### 8. Hub Memory

- **Daily reflection:** `cn-sigma/threads/daily/2026-05-20-gamma.md` — to be committed to cn-sigma after γ close-out
- **Adhoc thread:** `cn-sigma/threads/adhoc/activation-soul/` — update with #385 landing, soul now 3-surface (KERNEL + cap + clp), hub cleanup deferred
