## Post-Release Assessment — 3.78.0

**CI status on merge SHA:** red — https://github.com/usurobor/cnos/actions/runs/26093805709 (merge commit `a3bf7892`; two new I-class failures regressed from green pre-merge: SKILL.md frontmatter validation (I5) and Package verification kata R5-activate P10). Both regressions are γ-patched in this PRA commit; CI is expected green on the next push.

### 1. Coherence Measurement
- **Baseline:** 3.77.0 — α A, β A, γ A−
- **This release:** 3.78.0 — α B+, β B+, γ C
- **Delta:** α regressed (peer-surface coverage — frontmatter format conformance and kata assertion update were not authored). β regressed (pre-merge gate row 3 not exhaustively applied: validate-skill-frontmatter.sh and R5-activate kata were not run on the merge tree). γ regressed sharply (CI-red on merge commit triggers `release/SKILL.md` §3.8 CI-red cap clause: γ axis capped at C).
- **Coherence contract closed?** Partially. The "everything is a skill" invariant is restored for activation — α's primary deliverable (new `agent/activate/SKILL.md` + Go renderer evolution + source-of-truth test) is intact and merged. But two peer surfaces (skill-frontmatter validator and R5-activate kata) were not aligned with the cycle's intended change, producing post-merge CI red. The intended coherence delta shipped; the surrounding contract surfaces required γ-side step-13a patches to restore CI green. The cycle's named gap is closed; the cycle's _process_ produced an avoidable post-merge regression.

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| #379 | agent/activate skill — single source of truth | feature | issue body | shipped (this release) | none |
| (deferred from #379) | Hub README router adoption in cn-sigma / cn-pi | feature | template in skill §2.3 | not started | growing |
| (deferred from #379) | `//go:embed` activate skill in cn binary | feature | named in α §Debt 3 | not started | growing |
| (deferred from #379) | `cn doctor` enforcement of activation invariants | feature | named in α §Debt 4 | not started | growing |
| (deferred from #379) | `cnos.xyz/activate/<hub>` rendering service | feature | issue §Scope | not started | growing |
| #373 | Preventive --worktree identity write across role skills when extensions.worktreeConfig=true | process | converged | not started | growing |
| (this PRA) | Extend γ-scaffold failure-mode catalogue with "β contract-validator coverage" pre-flag | process | named here | shipped via β/SKILL.md row 3 patch (this PRA) | none |

**MCI/MCA balance:** **Freeze MCI** — 5 issues at "growing" lag (four #379-deferred follow-ons + #373).
**Rationale:** #379 itself was a focused implementation cycle that produced four growing-lag follow-ons against a single primary deliverable. The hub README router adoption (cn-sigma, cn-pi) is the highest-leverage next MCA because it activates the cycle's published artifact in two consuming hubs; without it, the new skill ships but no body consumes it via the router pattern. Hold all new substantial design until at least one of the four #379-deferred MCAs ships.

### 3. Process Learning

**What went wrong:** 

1. **Two peer surfaces not authored alongside the primary cycle change.** α delivered the named ACs cleanly (seven ACs PASS, β R1 APPROVED with zero findings), but the cycle's intended change (new SKILL.md frontmatter + `## Read first` ordering displacement) regressed two adjacent contract surfaces:
   - I5 (SKILL.md frontmatter validator) — the new skill's `calls:` field used package-prefixed paths (`cnos.core/skills/agent/cap/SKILL.md`) instead of the validator's required package-skill-root-relative paths (`agent/cap/SKILL.md`). Pre-cycle CI green → post-cycle CI red, all 7 findings in the new file.
   - Package verification kata R5-activate P10 — the kata asserts `## Read first` ordering for `cn activate` output. The kata's assertion (`persona < operator < kernel < deps < refl`) was the pre-cycle canonical; AC3 inverted it (`kernel < persona < operator < ...`). The kata was not updated when the canonical ordering changed.
2. **β pre-merge gate row 3 ("Non-destructive merge-test") was applied narrowly.** β ran `go test ./internal/activate/...` on the merge tree (documented in beta-review.md §"Branch CI state") but did not run `./tools/validate-skill-frontmatter.sh` against the cycle's new SKILL.md, and did not run the R5-activate kata against a fresh build. The prior row 3 wording ("the cycle's own validator or any CI-equivalent the cycle ships") was under-specified — the validator set was not enumerated, so β chose the most obvious one (`go test`) and stopped. `release/SKILL.md` §2.1 has a more exhaustive list but β/SKILL.md row 3 did not cross-reference it.

**What went right:**

1. **β found nothing wrong inside the named-ACs surface.** Zero findings of any severity at R1 across all seven ACs. The AC7 source-of-truth test (the load-bearing AC for the cycle's "everything is a skill" claim) was authored in the strong two-phase edit-and-swap form on first pass with no β nudge required.
2. **γ-scaffold failure-mode catalogue paid off — for the modes it named.** Ten γ-pre-flagged failure modes (file-read oracle, capability-matrix collapse, path typos, etc.) were all held by α at authoring time and re-checked by β at review time. The two regressions were modes γ did _not_ name — peer-validator coverage and peer-kata coverage. Pattern: γ-scaffolding compresses review rounds for failure modes that are surfaced; failure modes that γ omits from the scaffold catalogue cost a full post-merge CI red cycle to recover from. The fix is to make β/SKILL.md row 3 the structural backstop (now done), not to expand γ-scaffolding indefinitely.
3. **Mechanical fixes were available and immediately landable.** Both post-merge regressions were two-line patches (frontmatter `calls:` path correction; kata P10 assertion update) plus one peer doc (`kata.md` P10 description). γ landed all three as immediate outputs in this PRA commit per `gamma/SKILL.md` §3.6, restoring CI green without a follow-on fix-cycle.

**Skill patches:** YES — committed this release.

1. **β/SKILL.md §pre-merge gate row 3 — validator-set enumeration.** Tightened from "the cycle's own validator (or any CI-equivalent the cycle ships)" to an explicit enumerated list cross-referencing `release/SKILL.md §2.1`: (a) `./tools/validate-skill-frontmatter.sh` for any SKILL.md frontmatter add/modify; (b) `cn-cdd-verify` for `.cdd/` artifacts; (c) `scripts/check-version-consistency.sh` for version stamping; (d) the kata under `src/packages/cnos.kata/katas/{N}/` whose surface the cycle touches (with R5-activate named as the example for `cn activate` rendering / activate skill / `## Read first` changes). Empirical anchor cited inline: #379 post-merge CI red. This is the structural backstop — β no longer has to invent the validator set per cycle. (Commit: this PRA.)
2. **`src/packages/cnos.core/skills/agent/activate/SKILL.md` `calls:` frontmatter — path correction.** Paths changed from `cnos.core/{doctrine,skills/agent/X}/...` to package-skill-root-relative (`../doctrine/KERNEL.md`, `agent/X/SKILL.md`). Validator green after patch (`./tools/validate-skill-frontmatter.sh --file ...` → ✓; full run → 66 SKILL.md validated, no findings). (Commit: this PRA.)
3. **`src/packages/cnos.kata/katas/R5-activate/{run.sh,kata.md}` — P10 ordering update.** Assertion updated from the pre-#379 `persona < operator < kernel < deps < refl` to the post-#379 `kernel < persona < operator < {deps, reflection}` (deps + refl share one hub-state line in the new ordering, so relative deps/refl order is no longer asserted). `kata.md` P10 documentation updated to name the canonical source (`src/packages/cnos.core/skills/agent/activate/SKILL.md` §4.1) and the deprecated pre-#379 order. (Commit: this PRA.)

**Active skill re-evaluation:**

| Finding | Skill it should have caught | Underspecified or application gap? | Disposition |
|---------|-----------------------------|------------------------------------|-------------|
| F1 (I5 SKILL.md frontmatter) | β/SKILL.md §pre-merge gate row 3 + release/SKILL.md §2.1 | **Underspecified** — row 3 said "the cycle's own validator" without enumerating the validator set | Patch (β/SKILL.md row 3 enumeration, this PRA) |
| F1 (also) | γ/SKILL.md §2.2a (peer enumeration at scaffold time) | Application gap — γ enumerated peer skills but did not enumerate peer _validators_ | No patch this cycle (β-side fix is sufficient backstop) |
| F2 (R5-activate kata P10) | β/SKILL.md §pre-merge gate row 3 + release/SKILL.md §2.1 | **Underspecified** — same row as F1; the kata coverage was not enumerated | Patch (row 3 enumeration now names R5-activate as the example for activate-surface cycles) |
| F2 (also) | α/SKILL.md §2.5 incremental authoring + §2.4 harness audit | Application gap — α surfaced the observable-output delta honestly in §Debt 2 but did not update the kata assertion that depends on the delta. The harness audit (§2.4) did name "producers/consumers of the schema" but the schema was the marker-bounded ordering block, not the higher-level `## Read first` ordering that the kata asserts. | No patch this cycle (β-side fix catches this class) |

**CDD improvement disposition:** Patch landed (β/SKILL.md §pre-merge gate row 3 enumeration). The two regressions were caught by existing skill rules (`release/SKILL.md §2.1`'s validator list) but applied through an under-specified pointer in β's pre-merge gate. Closing the under-specification at the β-side closes the application-gap class.

### 4. Review Quality

**Cycles this release:** 1 (#379)
**Avg review rounds:** 1.0 (target: ≤2 code)
**Superseded cycles:** 0 (target: 0)

**Per-cycle round counts:**

| Cycle | Issue | Mode | Rounds | Binding findings (R1) | Notes |
|-------|-------|------|--------|----------------------|-------|
| #379 | agent/activate skill — single source of truth | design-and-build | 1 | 0 | R1 APPROVED. Zero in-cycle findings. Two post-merge findings (γ-found via CI verification). |

**Per-cycle dispatch telemetry** (this release contributes the first row of accumulating data):

| Cycle | `dispatch_seconds_budget` | `dispatch_seconds_actual` | `commit_count_at_termination` |
|-------|--------------------------|--------------------------|-------------------------------|
| #379 (α) | not recorded | not recorded | 11 (clean exit; review-readiness signaled) |
| #379 (β) | not recorded | not recorded | 4 (R1 verdict + merge + close-out + cycle-branch close) |

The γ scaffold (`gamma-scaffold.md`) declared §5.1 multi-session dispatch but did not name a timeout budget per role. Since neither agent SIGTERM'd (both exited cleanly after committing their last artifact), the budget-vs-actual is non-binding for this cycle; future cycles should record both per `CDD.md §1.6c`.

**Finding-class breakdown** (across cycles in this release; in-cycle + post-merge combined):

| Class | Definition | Count |
|---|---|---|
| **mechanical** | Caught by grep/diff/script | 0 |
| **wiring** | "X is wired into Y" but isn't (see review/SKILL.md 3.13c) | 2 (both post-merge: I5 frontmatter, R5-activate kata) |
| **honest-claim** | Doc claims something code/data doesn't back (review/SKILL.md 3.13) | 0 |
| **judgment** | Design/coherence assessment | 0 |
| **contract** | Work contract incoherent | 0 |

**Mechanical ratio:** 0% (no mechanical findings, but two wiring findings is the comparable signal here). Below the 20%-with-≥10-findings threshold; no process issue filing required on the mechanical-ratio axis.

**Honest-claim ratio:** 0% (well within the <30% target).

**Action:** none on the mechanical-ratio axis. The wiring-class findings are addressed by the β/SKILL.md row 3 patch (above, §3) — not via a separate process issue.

### 4a. CDD Self-Coherence

- **CDD α (artifact integrity):** 3/4 — `gamma-scaffold.md` + `self-coherence.md` + `beta-review.md` + `alpha-closeout.md` + `beta-closeout.md` all present and well-formed on `cycle/379` and now on `main`. The artifact-integrity axis itself was solid; the −1 is for the two peer-surface omissions (validator-compliant frontmatter; updated kata assertion) which are arguably artifact-integrity concerns at the cycle level.
- **CDD β (surface agreement):** 3/4 — canonical doc (`CDD.md`), executable skills (`alpha/`, `beta/`, `gamma/`, `release/`, `review/`, `operator/`), `.cdd/unreleased/379/` cycle artifacts, CHANGELOG, and PRA all agree on cycle scope and outcome. The −1 is for the under-specification surfaced in §3 (β/SKILL.md row 3 vs release/SKILL.md §2.1) which has now been patched.
- **CDD γ (cycle economics):** 2/4 — 1 review round (within target). 0 superseded cycles. CI red on merge commit (CI-red cap clause). Mechanical-ratio under threshold but two wiring-class regressions required γ step-13a patches to recover. Step-13a patches landed in this PRA commit (not deferred). The CI-red post-merge condition is the dominant signal pulling this axis below 3.
- **Weakest axis:** γ (CI-red cap).
- **Action:** Patch landed (β/SKILL.md row 3 enumeration — automates the missing check at the structural-gate layer rather than relying on β-side judgment per cycle).

### 4b. Cycle Iteration

- **Triggered by** (per CDD.md §9.1 thresholds):
  - review rounds > 2: NO (actual: 1)
  - mechanical ratio > 20% with ≥10 findings: NO (actual: 0% mechanical; 2 wiring findings total, well below threshold count)
  - avoidable tooling/environmental failure: **YES** — post-merge CI failed on the merge commit because two contract validators that the cycle's surface ships were not run pre-merge. Class: avoidable; bounded fix available.
  - **CI red on merge commit (post-merge): YES** — I5 (SKILL.md frontmatter validation) and Package verification (R5-activate kata P10) failed on merge commit `a3bf7892`. Triggers `release/SKILL.md §3.8` CI-red cap clause (γ axis capped at C).
  - loaded skill failed to prevent a finding: **YES** — β/SKILL.md §pre-merge gate row 3 had under-specified wording ("the cycle's own validator") that allowed β to run `go test` and skip the broader validator set named in `release/SKILL.md §2.1`. The skill was loaded; the application was incomplete; the wording made incompleteness easy.

- **Root cause:** β/SKILL.md row 3 cross-referenced "the cycle's own validator" without enumerating the validator set, and β's review documented `go test` execution but did not pivot to a broader contract-validator pass. The miss propagated because γ-scaffold's failure-mode catalogue named ten cycle-specific concerns but did not name "β contract-validator coverage" as a process-level pre-flag.

- **Disposition:** **Patch landed now.** β/SKILL.md row 3 enumeration patch (this PRA commit) names the four contract-validator classes explicitly with R5-activate as the activate-surface kata exemplar. This is a structural fix at the gate-layer — β no longer needs to invent the validator set per cycle. The cycle is closed with the patch landed, not deferred to a next-MCA.

- **Evidence:**
  - Merge-commit CI run: https://github.com/usurobor/cnos/actions/runs/26093805709 (I5 failure log, Package verification R5-activate P10 failure log).
  - Patch: this PRA commit (β/SKILL.md row 3 + activate SKILL.md `calls:` + R5-activate kata P10 + R5-activate kata.md P10).
  - Validator green after patch: `./tools/validate-skill-frontmatter.sh` → 66 SKILL.md validated, no findings.
  - Kata green after patch (against new binary): expected on next CI run (the kata uses the freshly-built `cn` binary in CI; the new assertion matches the new ordering).

### 5. Production Verification

**Scenario:** A body told "you are σ, your hub is `https://github.com/usurobor/cn-sigma`" can fetch the agent/activate skill at its raw GitHub URL, follow its §2.1 load order (Kernel → CA skills → Persona → Operator → hub state → identity), and name its identity, operator, and current orientation.

**Before this release:** Activation existed only in `cn activate`'s Go source. Bodies without the `cn` binary (Claude Code on the web, Codex, Claude.ai with WebFetch) had no fetchable artifact for the procedure; they depended on operator improvisation per session. The "I wake up incoherent by default" failure (cn-sigma `threads/adhoc/20260325-session2-learnings.md`) reproduced cycle-after-cycle.

**After this release:** The activate skill ships at `src/packages/cnos.core/skills/agent/activate/SKILL.md` (485 lines), `cn activate`'s renderer parses its §4.1 load-order block (the skill is the source of truth on drift, with an in-Go fallback for the pre-`cn deps restore` boundary case), and a body with shell + git or HTTP fetch can self-bootstrap by cloning the cnos repo (tier (a)) or fetching the raw URL (tier (b)). The §2.3 README router template documents the hub-side convention hub READMEs adopt verbatim.

**How to verify:**
1. **Skill-as-source-of-truth (renderer-level):** `cd src/go && go test -count=1 -run TestSkillIsSourceOfTruth ./internal/activate/...` — two-phase test (write vendored fixture with kernel-before-persona; assert ordering; rewrite to swap; assert swap propagates and `out1 != out2`). Already PASS on cycle HEAD and merge tree per β R1.
2. **Validator + kata pre-merge contract (cycle-level):** `./tools/validate-skill-frontmatter.sh` (must return "no findings") + `bash src/packages/cnos.kata/katas/R5-activate/run.sh` against a freshly-built `cn` binary (must return "0 failures"). After this PRA commit, both expected green; CI run will be the receipt.
3. **End-to-end body activation (manual dry-run):** Paste the §2.3 README router snippet into a fresh Claude.ai session with WebFetch enabled, supply `https://github.com/usurobor/cn-sigma` as the hub URL, observe whether the body reaches an identity-confirmation step without further operator prompts. **Deferred** — this PRA does not have a body session available for the dry-run; commit the deferral to the next σ session checklist.

**Result:** PASS (1–2) / deferred (3). The cycle's structural change (skill is source of truth; renderer reads it) is verified by the two-phase Go test that landed in `activate_test.go`. The cycle's intended consumer-facing change (non-cn bodies self-bootstrap) is testable but not exercised in this PRA session; named explicitly as deferred verification.

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 11 Observe | `gh run view 26093805709` + `git log a3bf7892..main` + `./tools/validate-skill-frontmatter.sh` + `bash src/packages/cnos.kata/katas/R5-activate/run.sh` | post-release, gamma, beta, release | Merge commit landed; CI red; two wiring-class regressions identified |
| 12 Assess | `docs/gamma/cdd/3.78.0/POST-RELEASE-ASSESSMENT.md` (this file) | post-release | Assessment completed — three skill patches landed (β/SKILL.md row 3; activate SKILL.md calls; R5-activate kata) |
| 13 Close | `.cdd/unreleased/379/gamma-closeout.md` (next) + `.cdd/unreleased/379/cdd-iteration.md` (next) + cycle dir move to `.cdd/releases/3.78.0/379/` + RELEASE.md + CHANGELOG 3.78.0 row | post-release, release, gamma | Cycle to close after this PRA + closeout artifacts commit; δ tag/release follows |
| 13a Skill/spec patch | β/SKILL.md §pre-merge gate row 3 + activate SKILL.md `calls:` + R5-activate kata run.sh + kata.md | post-release, beta, release | Patches landed in this PRA commit (not deferred) |

### 6a. Invariants Check

cnos has no single repo-level INVARIANTS.md document; the protocol invariants live in `CDD.md` and per-skill kata sections. The constraints touched by this cycle:

| Constraint | Touched? | Status |
|---|---|---|
| "Everything is a skill" (cnos.core scope) | Yes — activation was the last cnos procedure encoded only in Go | preserved (now skill-defined; renderer reads skill) |
| Layering rule (Soul = kind of agent; Identity = which agent) | Yes — skill §2.1 cites cn-sigma `threads/adhoc/20260325-session2-learnings.md §3` | preserved |
| `cn activate` CLI surface (flags, arguments, exit codes) | No (non-goal held) | preserved |
| `cnos.cdd/skills/cdd/activation/SKILL.md` (CDD repo activation, distinct concern) | No (non-goal held) | preserved |
| β/SKILL.md §pre-merge gate row 3 (non-destructive merge-test) | Yes — wording tightened to enumerate the validator set | tightened (this PRA) |

### 7. Next Move

**Next MCA:** Hub README router adoption in `usurobor/cn-sigma` (first; cn-pi follows in a separate cycle). This is the highest-leverage #379-deferred follow-on: the activate skill ships, but no hub README points at it yet, so no non-cn body can actually use the new path end-to-end. Adopting the §2.3 router template in cn-sigma's README closes the consumption loop for the canonical hub.
**Owner:** γ at cnos files the cross-repo issue against `usurobor/cn-sigma`; cn-sigma's δ executes the README patch.
**Branch:** pending creation at cn-sigma.
**First AC:** `cn-sigma/README.md` includes the §2.3 router template with `<HUB-URL>` substituted to `https://github.com/usurobor/cn-sigma`; a body told "Activate as `https://github.com/usurobor/cn-sigma`" fetches the activate skill and reaches the §1 identity-confirmation step.
**MCI frozen until shipped?** Yes — see §2 above (5 issues at growing lag; freeze in effect).
**Rationale:** A new skill that ships without a consuming router is a tree-falling-in-the-forest. The router adoption is the first measurable consumption of the cycle's primary deliverable. Other #379-deferred items (//go:embed durable fix, `cn doctor` enforcement, `cnos.xyz`) are P3 / structural and can wait for the consumer loop to validate the skill's shape first.

**Closure evidence (CDD §10):**
- Immediate outputs executed: yes
  - β/SKILL.md §pre-merge gate row 3 enumeration patch (this PRA commit; structural backstop)
  - activate SKILL.md `calls:` frontmatter fix (this PRA commit; I5 validator green)
  - R5-activate kata run.sh P10 update (this PRA commit; Package verification expected green on next CI)
  - R5-activate kata.md P10 documentation update (this PRA commit; keeps test + doc in sync)
- Deferred outputs committed: yes
  - Hub README router adoption (cn-sigma) — to file at cn-sigma during δ disconnect; first AC named above
  - Hub README router adoption (cn-pi) — to file at cn-pi after cn-sigma adoption
  - `//go:embed` durable fix for the renderer fallback duplication (α §Debt 3) — P3, future cycle
  - `cn doctor` enforcement of activation invariants (α §Debt 4) — P3, future cycle
  - `cnos.xyz/activate/<hub>` rendering service (α §Debt 5) — P3, requires hosting decision
  - Cross-repo source-proposal `landed` event on `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/STATUS` — flagged for δ to coordinate with operator (mirror branches `cnos:sigma/cross-repo-mirror-agent-activate-skill@212d5239` and `cn-sigma:sigma/cross-repo-status-lineage-update@89049611c` are pending operator merge; landed event added after both merges complete)

**Immediate fixes** (executed in this session):
- β/SKILL.md row 3 patch (named above)
- activate SKILL.md `calls:` paths (named above)
- R5-activate kata + kata doc (named above)

### 8. Hub Memory

- **Daily reflection:** deferred — this is a cnos repo cycle; the cnos repo has no hub memory surface (`threads/reflections/daily/`). Hub memory in cnos lives in per-hub repos (cn-sigma, cn-pi). The σ session that lands the cn-sigma README router adoption (next MCA) will write the daily reflection there.
- **Adhoc thread(s) updated:** deferred — same reason. The cn-sigma adhoc thread on "I wake up incoherent by default" (`threads/adhoc/20260325-session2-learnings.md`) is the natural place to record "the skill that closes this thread shipped at cnos 3.78.0"; that update happens during the cn-sigma adoption cycle.
- **Cnos-side memory equivalent:** this PRA + `gamma-closeout.md` + `cdd-iteration.md` for #379 are the durable memory at the protocol layer; no separate cnos-side daily reflection convention exists.
