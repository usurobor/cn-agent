## Post-Release Assessment — 3.58.0

> Rollup release covering cycles **#264** (deterministic tarballs, small-change), **#262** (packlist derivation, partial — AC1/AC5/B1 deferred), **#266** (dist out of git, L7 MCA), **15 CDD skill patches landed direct-to-main**, and **CTB §8.5.1 evidence consolidation**. Release commit `8b1c348`. PRA authored by γ per CDD `cafb399` (PRA ownership moved from β to γ).

### 1. Coherence Measurement

- **Baseline:** 3.57.0 — α A-, β A-, γ A-, level L6 (cycle cap)
- **This release:** 3.58.0 — α A-, β A-, γ A-, **level L7**
- **Delta:**
  - **α** held A- across the rollup. Substantively clean on each diff (CI green at every merged head; `go test ./...` and dispatch-boundary check green). Cross-cycle pattern: 6 review rounds across the two substantial cycles (#262: 3 rounds / 7 findings; #266: 3 rounds / 4 findings), all mechanical or judgment-on-deferral, none semantic. Loaded-skill-misses recurred at scale (#262 cross-Go-toolchain non-determinism wasn't named in `eng/go §2.13` until α's #262 close-out triggered the patch; #266 peer-enumeration scale gap wasn't named in `alpha §2.3` until α's #266 close-out triggered the patch landed in this PRA commit). The pattern is "the skill, as written, did not predict the failure shape that appeared" rather than "α failed to apply the skill" — both patches landed as derive-from-cycle MCAs, which is the correct response to a skill-scope gap.
  - **β** held A-. Independent verification was thorough: F1 (a)/(b) framing on #266 was sharp judgment (offering α a leverage-preserving alternative to running the race one more time); polling-gap intake-time flag on #266 surfaced and was patched mid-cycle; environment-provided implementation-instruction refusal in #266 (β explicitly refused to implement when the harness directed implementation work, per the new role-rule §1). Symmetric β-side review-skill miss on #266 F3 (named 1 of 3 sites; F3-bis the consequence) — the matching `review §2.1.3` patch lands in this PRA commit alongside the α-side patch. Shared-identity review degradation continued: 4 cycles now (#260, #265, #267, and the half-rounds in #264) under `usurobor` shared identity with no native APPROVE / REQUEST_CHANGES surface; β's verdict-in-text discipline held throughout.
  - **γ** held A- but with operationally-visible friction: 15 skill patches shipped as immediate MCAs across the rollup span (close-out voice constraint, dispatch-prompt trimming, polling replacing event subscriptions, β role-violation refusal, large-file authoring rule consolidation, PRA ownership relocation, etc.); CDD package audit filed and dispatched as #268; issue-quality gate compensated for #266 issue gaps before α dispatch (#266 patched with priority/work-shape/Tier 3/constraints/sequencing). The friction: one γ close-out for #266 was authored by γ (longer, full triage) and a parallel terser version was authored by the operator (during a state where γ's branch wasn't yet merged); the duplicate was resolved at PRA-commit time by adopting the fuller γ version. Branch cleanup deferred 4 cycles in a row on the same HTTP 403 ref-write env constraint (3.57.0 tag push, #262 / #265 / #267 branch deletes, this cycle's 3× branch deletes) — pattern is now over the "4 consecutive cycles" threshold β's #267 close-out flagged. Hub memory writes deferred (`cn-sigma/threads/` is operator-environment, not in-repo). The A- captures the messiness honestly; the L7 cycle-level survives because the substantive MCA is durable and patches did land in time.
  - **Level — L7 achieved:** the cycle shipped a system-shaping MCA that eliminates the rebase-race friction class structurally (`dist/` no longer in git; I3 rebuilds-from-source; main-side `src/packages/` edits no longer race against committed tarballs). The MCA was exercised on its own purpose during the #266 cycle itself (3 commits landed during review and 4 between approval and merge in `src/packages/cnos.cdd/skills/cdd/` — post-merge none of those moves produce committed-dist diffs). L7 cycle level is earned cleanly: the L7 substantive change shipped, the cycle iteration triggers (review rounds > 2 in both #262 and #266; loaded-skill miss in both) all have landed-or-deferred dispositions, and the §1.4 / role-skill / dispatch-model / PRA-ownership iterations across the rollup span are themselves system-shaping for the CDD method. The "lowest miss" rule at rollup scale is the cycle-execution miss in #266 (intra-doc arithmetic drift reached review across 3 sites of one doc) — but that miss has a landed mechanism patch (peer-enumeration scale, this commit), making it L7-shaped: the failure class can no longer recur silently because the skill now names the scale.
- **Coherence contract closed?** Yes for the L7 substantive intent (rebase-race class eliminated, verified per §5 below). Partial for the wider rollup intent: #262 (`cn pack` rename + `--list`/`--dry-run` + `.gitignore` honoring + root-file-allowlist narrowing) is intentionally only partly shipped per β's deferred-by-design-scope discipline (β #265 close-out §B1) — the open ACs ride into the next packaging cycle (see §7). #268 (CDD package convergence from the audit's 34 findings) is filed and dispatched, not yet landed; D-class fixes from the audit are scope-deferred to that cycle. Both partial-closures are concretely deferred, not silently dropped.

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| #266 | dist/ out of git (L7 MCA) | feature/process | converged | **shipped this release** | none |
| #264 | Deterministic package tarball builds | bug | converged | **shipped this release** | none |
| #262 | `cn pack` packlist derivation | feature | converged | **partially shipped** (slice landed; AC1 rename + AC5 `--list`/`--dry-run` + AC3 root-file narrowing + `.gitignore` honoring deferred) | low |
| #268 | Converge cnos.cdd skill-program contracts (CDD package audit) | feature/process | converged (34 findings cataloged) | **dispatched, not landed** | low |
| #249 | .cdd/ Phase 1 audit trail | feature | converged | partially shipped (`.cdd/releases/<v>/{α,β,γ}/N.md` pattern in use across 3.57.0 + 3.58.0; close-out moves automated by `99892ea`) | low |
| #256 | CI as triadic surface | feature | converged | not started | growing |
| #255 | CDD 1.0 master | tracking | converged | in progress (this release ships dispatch-model, polling, role boundaries, PRA ownership, large-file authoring, close-out voice, audit infra — about half of master scope) | low |
| #244 | kata 1.0 master | tracking | converged | in progress (Tier-1 + Tier-2 + framework shipped 3.54.0–3.55.0; M0/M4 work depends on #243/#245) | low |
| #245 | cdd-kata 1.0 (prove CDD adds value) | feature | converged | stubs only | growing |
| #243 | cnos.cdd.kata cleanup (M0/M4 + remove runtime + honest stubs) | feature | converged | partially done | low |
| #242 | .cdd/ directory layout, lifecycle, and retention | design | converged | not started (but `.cdd/releases/{version}/` move in 3.58.0 is adjacent partial-impl) | low |
| #241 | DID: triadic identity | design | converged | not started | growing |
| #240 | CDD triadic protocol (git transport, .cdd/ artifacts) | design | converged | partial (Phase 1 in 3.55.0 + .cdd/releases/ moves in 3.58.0) | low |
| #238 | Smoke release bootstrap compatibility | feature | converged | not started | growing |
| #235 | `cn build --check` validates entrypoints + skill paths against manifest | feature (P1) | converged | not started | growing |
| #230 | `cn deps restore` version upgrade skip | bug (P1) | converged | not started | growing |
| #218 | cnos.transport.git design | design | converged | not started | growing |
| #216 | Migrate non-bootstrap kernel commands to packages | feature | converged | not started | growing |
| #199 | Stacked v3.39.0 + v3.40.0 PRAs | docs | converged | not started | growing |
| #193 | Orchestrator `llm` step execution | feature | converged | not started | growing |
| #192 | Runtime kernel rewrite (OCaml → Go) | feature | converged | Phase 4 complete; remaining: orchestrator + wider port | low |
| ~~#261~~ | (Skill activation Go port) | — | — | shipped 3.57.0 | — |

(Older items #190 / #189 / #186 / #182 / #181 / #180 / #175 / #170 / #168 / #162 / #156 / #154 / #153 / #149 / #148 / #142 / #141 / #138 / #135 / #132 / #124 remain open and stale across multiple cycles. They are not promoted to the active lag table this release; they are tracked in the open-issue list and will be re-evaluated under §3.4 stale-backlog rule when a future cycle has slack to triage them.)

**Growing-lag count:** 11 (vs 10 in 3.57.0 — closed: #261 / #264 / #266; opened: #268; converted from "growing" to "low" via partial implementation: none net new "low"; #249 stays low via 3.58.0 `.cdd/releases/` automation; remaining 11 growing items unchanged in design state). Net change roughly flat — three items closed, one new item added (smaller scope), backlog quality slightly better but count just over 3.57.0's level.

**MCI/MCA balance:** **Freeze MCI continues** — 11 issues at "growing" lag, well over the 3-issue freeze threshold. No new substantial design work until the MCA backlog drops below threshold.

**Rationale:** The freeze pattern continues to produce the intended pressure. This release closed three real gaps (one L7 MCA + one bug + the partial #262 packaging slice) and shipped a 15-patch CDD-method coherence sweep that addresses recurring agent-compliance failures. The next cycle should aim at #268 (CDD package convergence from the audit) since it is already dispatched, then continue picking the next-smallest-clear-fix-path growing-lag item per the 3.55.0–3.57.0 pattern. #230 (`cn deps restore` version-upgrade skip, P1, one-function fix) remains the natural pick after #268 — both unblock the substrate's authority chain (deps → install → vendor) that the L7 MCA depends on for adoption.

### 3. Process Learning

**What went wrong:**

- **Cross-cycle review-round inflation: 3 rounds in #262 + 3 rounds in #266** (target ≤2). #262's two re-emergences of the rebase race (D1 + D1') are the failure class #266 then eliminates — the cost of running the race one more time was paid before the structural fix landed. #266's third round was an authoring-side intra-doc arithmetic drift (F3 / F3-bis) where R1 named 1 of 3 sites and α's R1 fix trusted the count. Both cycles' third rounds map to skill-scope gaps now patched in this PRA commit, but at the cost of two cycles' worth of cycle-economics drag.

- **Loaded-skill scale mismatch (peer enumeration).** `alpha §2.3` and `eng/go §2.12` both already named the peer-enumeration rule. Their examples were all cross-producer scale (Go parser vs shell writer; Go type vs YAML workflow step). Application biased toward that scale and missed two smaller scales of the same rule: intra-doc repetition (one document carrying the same fact across multiple sentences) and commit-message closure claims (a commit restating a finding's resolution without running the grep that would tell you otherwise). α's #266 close-out names this directly. The corresponding β-side miss in `review §2.1.3` (mechanical scan that surfaced a numeric drift but didn't grep for every occurrence before filing) is the symmetric class.

- **Pre-review-gate transient-row drift.** `alpha §2.6` listed 10 gate rows without distinguishing rows that describe external state from rows that describe artifact state. Rows 1 (branch rebased) and 10 (CI green) are transient — they describe state that can change between the moment α writes them and the moment β reads them. The other 8 rows are durable. #266 F1 / F2 both shipped as transient rows written at PR-open time and not refreshed before β read the PR ~9 minutes later (main moved 3 commits in the window; CI went green). The fix is structural: name the distinction in §2.6, add a re-validation step in §2.7 immediately before requesting review.

- **γ-side close-out duplication.** Two parallel γ close-outs for #266 existed — the γ-authored version on the γ session branch (370 lines, full triage table, deferred-output bodies, harness-403 evidence, hub memory content) and a 34-line operator-authored stub at `.cdd/releases/3.58.0/gamma/266.md`. The stub claimed three "Immediate MCA — patched alpha/SKILL.md §2.3 / §2.6 / §2.7 / review §2.1.3" patches that were never landed on main during the release. The discrepancy was resolved at this PRA's merge (γ version adopted; patches landed in this commit alongside the PRA). Operationally this means **3.58.0 was released without the skill patches its γ close-out claimed had landed** — an honesty-of-record gap that the PRA closes by landing the patches now in the same commit. No 3.58.0 consumer-visible behavior depends on these patches; the cost is one release cycle of the L6 cycle-execution failure remaining technically possible until this commit lands. Recording explicitly so the precedent is auditable.

- **Branch-cleanup deferred 4 cycles.** HTTP 403 on `git push origin --delete <branch>` (and `git push origin <tag>`) has now blocked: 3.57.0 tag push (β β-deferred to γ), #262 close-out branch delete, #265 close-out branch delete, #267 close-out branch delete, and this PRA cycle's 3× branch deletes (`claude/alpha-tier-3-skills-0OOu5`, `claude/implement-alpha-algorithm-cky3D`, `claude/implement-alpha-algorithm-gXXy3`). β's #267 close-out flagged "if γ observes this is now four consecutive cycles, that's threshold for a harness-level tooling issue per `260.md` Finding 2's documented precedent." Threshold is met; an MCI for the harness ref-write 403 pattern is drafted in `.cdd/releases/3.58.0/gamma/266.md` §Branch cleanup and is one of two MCIs deferred-to-operator on GitHub MCP auth refresh.

**What went right:**

- **L7 MCA shipped and exercised on its own purpose during the cycle.** #266 PR #267 merged with main moving 3 commits during review and 4 between approval and merge — all in `src/packages/cnos.cdd/skills/cdd/`. Under the pre-#266 regime each of those moves would have invalidated PR #267's committed dist. Post-#266 there is no committed dist to invalidate. The cycle that filed the MCA was also the cycle that proved the MCA works.

- **β (a)/(b) framing on #266 F1 was a sharp judgment call.** Rather than recommend rebase (which would have run the race exactly one more cycle), β offered α two paths and recommended (b) — amend the gate row to record actual state honestly. The PR's whole thesis is that α should not need to run the race to win; (b) preserves that thesis. α chose (b); the cycle converged. This is the kind of scope-aware reviewing that the review skill's §2.0 Issue Contract gate exists to enable, exercised cleanly under shared-identity degradation where native APPROVE/RC events were unreachable.

- **CDD method iteration was substantial and durable.** 15 skill patches across the rollup span (close-out voice constraint, dispatch-prompt trimming, polling replacing event subscriptions, β role-violation refusal, large-file authoring rule consolidation, PRA ownership relocation from β to γ, γ subscribes to issue+PR, γ reads last PRA first, dispatch model produces both prompts at dispatch time, β starts intake immediately, etc.). Combined with the CDD package audit (#268 dispatched), the method's coherence at L7 scale is now better-instrumented than at any prior release. The CDD self-coherence axis read in §4a holds despite this volume because each patch derives from concrete cycle evidence.

- **§7.0 "no follow-up" rule held cleanly across both substantial cycles.** Every #262 finding (D1, D1', C1, C2, B1, A1, A2) and every #266 finding (F1, F2, F3, F3-bis) was either resolved on-branch or explicitly deferred-by-design-scope with an issue filed (#264 from #262 F2, #266 from #262 D1+D1' rebase race, #268 from CDD audit). No "approved with follow-up" anti-pattern appeared.

- **Issue-quality gate compensated structurally.** #266 was patched before α dispatch with the 5 missing γ §2.4 fields (priority, work-shape, Tier 3 skills, active design constraints, sequencing) per "fix the issue, don't compensate in the prompt." α's #266 close-out §3 explicitly cites the up-front trade-off-space mapping (β's candidate MCAs A/B/C + γ's design-constraint links) as why this cycle's design pass was cheap. The gate is doing its job at the contract level it's supposed to.

- **β role-rule §1 (refuse environment-provided implementation work) caught a real role-leakage attempt.** During #266 intake, the harness directed β with "develop and commit on `claude/implement-beta-skill-loading-LNAOc`." β refused (correctly per the new rule landed in `f85d0c3`) and the cycle continued cleanly. Evidence-of-one for the rule's value, but exactly the kind of compliance gap the rule was added to prevent.

**Skill patches — landed in this PRA commit (CDD §10.1 step 12a):**

The peer-enumeration scale gap and the pre-review-gate transient-row class were identified by #266 cycle iteration. The patches were drafted in the γ session branch (`83028dd`) but not merged to main during the 3.58.0 release. They land in **this PRA commit** alongside the assessment, synced across all affected surfaces:

| Patch | File | Section | Derives from |
|---|---|---|---|
| Peer enumeration applies at any scale | `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | §2.3 | #266 F3 / F3-bis (intra-doc), #266 commit `9f162dc` (commit-message closure claim) |
| Transient vs durable pre-review-gate rows | `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | §2.6 | #266 F1 / F2 (rows 1 + 10 written at PR-open, not refreshed) |
| Re-validate transient rows before requesting review | `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | §2.7 | #266 F1 / F2 |
| Numeric / value-repetition: grep every occurrence before filing | `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` | §2.1.3 | #266 F3 / F3-bis (β-side symmetric rule) |

All four patches derive-from a cycle finding, follow the "describe the rule + cite the example" pattern of existing skill text, and pass the §10.1 gate (canonical source under `src/packages/`; no separate package-visible loader; no human-facing pointer surface that exposes the changed rule).

**Active skill re-evaluation:**

- `alpha/SKILL.md §2.3` Peer enumeration — **underspecified, patched.** The rule's exemplars biased application toward cross-producer scale and missed intra-doc + commit-message scales. Patch lands this commit.
- `alpha/SKILL.md §2.6` Pre-review gate — **underspecified, patched.** Did not distinguish transient from durable rows. Patch lands this commit.
- `alpha/SKILL.md §2.7` Request review — **underspecified, patched.** Did not require re-validation of transient gate rows immediately before requesting review. Patch lands this commit.
- `review/SKILL.md §2.1.3` Mechanical diff scan — **underspecified, patched.** Did not require grep-every-occurrence before filing numeric/value-repetition findings. Patch lands this commit.
- `eng/go §2.12` Schema and compatibility — **already covered cross-producer scale; intra-doc scale now lives in `alpha §2.3` instead.** No `eng/go` patch this cycle; the principle generalizes through the `alpha/SKILL.md` patch.
- `eng/go §2.13` Determinism and reproducibility — **patched mid-release** (commit `3e850f3`, prior to release): cross-toolchain non-determinism axis added per α's #262 close-out finding. Already shipped.
- `review/SKILL.md §2.0` Issue Contract gate — **worked as written.** β produced AC coverage + Named Doc Updates + CDD Artifact Contract tables on every review round. Application gap nowhere.
- `review/SKILL.md §2.2.1a` Input-source enumeration — **worked as written.** β enumerated input sources for the F1 fix path correctly.
- `review/SKILL.md §2.2.8` Authority-surface conflict — **worked as written.** β caught the dispatch-vs-PR scope drift in #262 (C2), the AC2 ClassSkills incoherence in #262 (C1), and the gate-row-honesty gap in #266 (F1) all via §2.2.8.
- `cdd/gamma/SKILL.md §2.4` Issue-quality gate — **worked as written.** γ caught the 5 missing fields in #266 before α dispatch. Patch-the-issue-not-the-prompt rule held.
- `cdd/post-release/SKILL.md` Ownership move (β → γ) — **landed mid-cycle** (`cafb399`). This PRA is the first authored under the new ownership; observably the new boundary is cleaner (γ doing cycle-level observation while β focuses on review→merge→release→close-out).

**CDD improvement disposition:** patches land in this commit as enumerated above; this is the §10.1 step-12a sync. Beyond those, no separate CDD spec-level patch is needed this cycle — the 15 skill patches that landed across the rollup span (already on main pre-PRA) plus the 4 patches in this commit cover every recurring class of finding from cycles #262 and #266. No "noted for next cycle" residue.

### 4. Review Quality

| Metric | Value | Target | Status |
|---|---|---|---|
| PRs in this rollup | 2 reviewed (#265 / #267) + 1 small-change-direct (#264) + 15 direct-to-main CDD skill patches + 1 CDD audit doc | — | — |
| Avg review rounds (reviewed PRs) | **3.0** (#265: 3 rounds, #267: 3 rounds) | ≤2 code | **over target — §9.1 trigger fires per cycle** |
| Superseded PRs | 0 | 0 | met |
| Total review findings (reviewed PRs) | 11 — #265: 7 (D1, D1', C1, C2, B1, A1, A2) + #267: 4 (F1, F2, F3, F3-bis) | — | — |
| Mechanical findings | 8 — #265: 4 (D1+D1' rebase race, A1 branch name, A2 missing CDD Trace) + #267: 4 (F1, F2, F3, F3-bis) | — | — |
| Judgment findings | 3 — #265: 3 (C1 ClassSkills, C2 dispatch-PR scope drift, B1 root-file filter) | — | — |
| Mechanical ratio | **8/11 = 72.7%** (above 20% threshold) | <20% threshold | **over threshold AND total findings ≥ 10 → process issue required** |
| Action | **§4 process issue required.** See disposition below. | — | — |

**Round-by-round (per cycle):**

#### #265 (cycle #262: packlist derivation)

- **Round 1.** β REQUEST CHANGES + 6 findings (D1 base drift, C1 partial AC2, C2 dispatch↔PR scope conflict, B1 root-file filter weak, A1 branch name, A2 missing CDD Trace). D1 was caught via `git merge-base` showing PR base 1 commit behind main; reproduced locally by rebuilding tarballs against the merge commit and showing checksum drift. C1 was caught via §2.2.8 authority-surface audit. C2 was caught via §2.2.8 cross-comparing γ's β-dispatch prompt against the PR body's deferred-AC list.
- **Round 2.** β REQUEST CHANGES + D1' (rebase race re-emerged after main moved during review). α's R1 D1 fix was correct but the race re-fired before β could approve. R2 also surfaced cross-Go-toolchain non-determinism (CI 1.22 vs local 1.24.7); α's R2 fix bumped Go to 1.24 across `go.mod` + 3 workflows. Right-sized scope expansion, declared in commit message.
- **Round 3.** β APPROVED on `9a1c7df`. Merged.

#### #267 (cycle #266: dist out of git)

- **Round 1.** β APPROVED with 3 mechanical findings (F1 base-drift gate-row claim stale, F2 CI-green checkbox unticked, F3 tarball-count off-by-one in DESIGN-266 §Concrete changes step 1). β proposed (a) rebase or (b) amend the gate row honestly for F1; α chose (b) per the PR's leverage-preservation thesis.
- **Round 2.** β REQUEST CHANGES (narrow). F1 + F2 + F3 fixes verified at `9f162dc`; β R2 mechanical re-scan surfaced F3-bis (2 further sites of the same arithmetic drift + a `.gitignore` line-count A-polish). α's R1 F3 fix had trusted β's single-site naming and asserted "fixes the one remaining mismatch" without grepping the doc.
- **Round 3.** β APPROVED on `7e76798` after F3-bis fix at `43f5a89`. Merged at `a05af27`.

**Finding-count note:** 11 total findings is just over the §9.1 10-finding floor. Mechanical ratio 72.7% is well above the 20% threshold. Per `post-release/SKILL.md` step 5.5: *"If >20% AND total findings ≥ 10, file an issue to add the missing pre-flight check."* Both conditions met for the rollup; per-cycle each was below the floor.

**Action — process issue disposition:**

The required pre-flight checks for the dominant mechanical classes already exist or land in this commit:

- **Rebase-race / committed-dist-drift class** (D1 + D1' + I3 red across multiple cycles): **eliminated structurally** by #266's MCA itself. No pre-flight check is the right shape for this class — the structural fix is. Recurrence will be tracked per #266 AC3 across the next 3–5 packaging-adjacent cycles.
- **Branch-name format drift class** (#265 A1 + #267 close-out §γ-reports §4): MCI for branch-name validation drafted in `.cdd/releases/3.58.0/gamma/266.md` §Deferred outputs §1; deferred-to-operator on GitHub MCP auth refresh.
- **Cross-Go-toolchain non-determinism class** (#265 F2 sibling): patched in `eng/go §2.13` mid-release (`3e850f3`). Already shipped.
- **Intra-doc arithmetic drift / closure overclaim class** (#265 A2 missing CDD Trace + #267 F3 + F3-bis): patched in this commit (alpha §2.3 + review §2.1.3). The α-side patch covers commit-message closure claims, which catches the #265 A2 sibling pattern (PR body missing a CDD Trace was a peer-enumeration miss against the PR template's required surface).
- **Pre-review-gate transient-row drift class** (#267 F1 + F2 + #265 D1 + D1' partial overlap): patched in this commit (alpha §2.6 + §2.7). The structural fix to D1/D1' is #266 itself; the gate-row-honesty patch addresses the per-PR honesty-of-record dimension.

**Disposition: no separate process issue filed.** Every mechanical-ratio-driving finding class either (a) had a structural elimination already shipped (rebase race, cross-toolchain, intra-doc drift), (b) has a deferred-to-operator MCI body drafted (branch name, harness-403), or (c) is patched in this commit. Filing a meta-process issue would duplicate work already done across the rollup span.

### 4a. CDD Self-Coherence

- **CDD α (artifact integrity):** **3/4.** Final diffs at every merged head carry required artifacts: PR body §CDD Trace, AC mapping, self-coherence section, peer enumeration table, INVARIANTS check. `go test ./...` clean; `go vet` clean; dispatch-boundary CI grep clean; deterministic tarball verification reproducible. Drop from 4: round-1 artifact integrity in #265 (D1 rebase-race + cross-toolchain) and #267 (F1/F2 transient-row claims, F3 intra-doc drift) was lower than final-merged-head integrity. Skill patches now landing in this commit address the recurring authoring-time gaps.

- **CDD β (surface agreement):** **3/4.** PACKAGE-AUTHORING ↔ PACKAGE-ARTIFACTS ↔ DESIGN-266 ↔ BUILD-AND-DIST migration step 2 ↔ CHANGELOG TSC row ↔ RELEASE.md ↔ this PRA all agree at release-time state. `cn build` ↔ `pkg.ContentClasses` ↔ activation `ClassSkills` constant agree (#262 C1 fix). I3 ↔ committed dist-state agree (post-#266: I3 rebuilds-from-source against `checksums.txt`). Drop from 4: F3-bis revealed intra-doc arithmetic drift in DESIGN-266 across 3 sites — caught and reconciled before merge but the surface-agreement axis touched a real miss. Symmetric β-side review-skill patch lands this commit.

- **CDD γ (cycle economics):** **3/4.** 15 immediate-MCA skill patches landed across the rollup; CDD package audit filed and dispatched as #268; issue-quality gate compensated structurally before #266 dispatch; PRA ownership boundary (cafb399) clarified mid-rollup; 3 cycles closed. Drop from 4: 6 review rounds across #262 + #266 (above 2-round target on each cycle); γ-close-out duplication for #266 (full version vs. operator's stub) resolved at PRA-merge time rather than at the original close-out commit; branch cleanup deferred 4 consecutive cycles on harness 403; the 4 alpha/review skill patches that #266's γ close-out claimed had landed actually didn't reach main until this PRA commit (one-cycle honesty-of-record gap, now closed).

- **Weakest axis:** **γ (cycle economics)** — same as 3.57.0. The recurring driver is round-count drift (review rounds > target on substantial cycles) plus the friction of harness env constraints (HTTP 403 on ref-write) that γ keeps deferring without an MCA shipping. The 4 skill patches landing this commit reduce future authoring-time mechanical drift; the harness-403 MCI draft (in `.cdd/releases/3.58.0/gamma/266.md`) is the deferred-to-operator path for the env-constraint class.

- **Action:** patches landing in this commit (alpha §2.3 + §2.6 + §2.7 + review §2.1.3) close the round-count root cause for cycles #262 and #266. The harness-403 MCI awaits GitHub MCP auth refresh (deferred-to-operator). No further skill patch this PRA. Future γ axis improvement depends on shipping the harness-403 MCI and on at least one cycle landing in ≤2 rounds (the target) to prove the 4 patches measurably bend the round-count curve. If 3.59.0's primary cycle still hits 3 rounds with the same root-cause class, escalate to a meta-skill patch in `cdd/gamma/SKILL.md §2.4` issue-quality gate (e.g. require α to attach a peer-enumeration plan to the PR body for any doc with ≥3 count-bearing sentences).

### 4b. Cycle Iteration

`CDD §9.1` triggers across the rollup span:

| Cycle | Trigger | Fired? | Disposition | Evidence |
|---|---|---|---|---|
| #264 | review rounds > 2 | no | n/a (small-change-direct, no review surface) | `8265fe0` |
| #264 | mechanical ratio > 20% with N ≥ 10 | no | n/a (no findings) | — |
| #264 | tooling/env failure | no | n/a | — |
| #264 | loaded-skill miss | no | n/a | — |
| #262 | review rounds > 2 | **yes** (3 rounds) | **patch landed** mid-cycle: `eng/go §2.13` cross-toolchain rule (`3e850f3`); rebase-race class **eliminated** by #266 itself (`a05af27`) | β #265 close-out §Cycle Iteration; γ #262 close-out §γ-reports §1 |
| #262 | mechanical ratio > 20% with N ≥ 10 | partial (4/7 = 57%, but N=7 < 10-floor) | recorded; no automatic trigger | β #265 close-out |
| #262 | tooling/env failure | **yes** (cross-toolchain × 1; rebase-race × 2) | **both patched** — Go 1.24 pin shipped; rebase-race class shipped via #266 | β #265 close-out §γ-reports §1; γ #262 close-out §γ-reports §1 |
| #262 | loaded-skill miss | **yes** (`eng/go §2.13` did not name cross-toolchain axis) | **patch landed** (`3e850f3`) | γ #262 close-out §α candidate #2 |
| #266 | review rounds > 2 | **yes** (3 rounds) | **patches land in this PRA commit:** `alpha §2.3` (peer enum scale), `alpha §2.6` + §2.7 (transient gate rows), `review §2.1.3` (grep every occurrence) | α #266 close-out §Cycle Iteration; β #267 close-out §Cycle Iteration; γ #266 close-out §Cycle Iteration |
| #266 | mechanical ratio > 20% with N ≥ 10 | partial (4/4 = 100%, but N=4 < 10-floor) | recorded; no automatic trigger | β #267 close-out |
| #266 | tooling/env failure | no | n/a (Go 1.24 + deterministic tarballs both held) | β #267 close-out §9.1 |
| #266 | loaded-skill miss | **yes** (`alpha §2.3` peer-enumeration scope did not generalize to intra-doc / commit-message scale) | **patches land in this PRA commit** (see above) | α #266 close-out §Cycle Iteration §Skill impact; β #267 close-out §Skill impact |
| Rollup | mechanical ratio > 20% with N ≥ 10 | **yes** (8/11 = 72.7%, N=11 ≥ 10) | **disposition in §4 Action above** — no separate process issue; every contributing class has a landed/deferred fix | this PRA §4 |

**Root cause across the rollup:** The pattern that drove cycle-economics drag in both #262 and #266 is **loaded-skill-scope mismatch** — skills whose stated rule covered the failure class by generalization but whose examples biased application toward a narrower scale. `eng/go §2.12` covered cross-producer schema audit but exemplified only Go-vs-shell (missed cross-toolchain runtime). `alpha §2.3` covered peer enumeration but exemplified only cross-surface peer sets (missed intra-doc + commit-message scales). `alpha §2.6` listed gate rows uniformly without distinguishing transient from durable. `review §2.1.3` listed mechanical scans but didn't include grep-every-occurrence-before-filing for value-repetition findings.

The fix shape is consistent: **add the smaller-scale application as an explicit case in the skill text, derive-from-cycle, keep the cross-surface exemplar.** All four patches in this commit follow that shape.

**Cycle level (rollup):** **L7 achieved.**

- L5 (local correctness): met on every merged head across all three cycles.
- L6 (system-safe execution): partial misses in #262 (cross-toolchain non-determinism reached review) and #266 (intra-doc arithmetic drift reached review). Both caught and reconciled before merge; both have landed mechanism patches in or before this PRA commit (Go 1.24 pin, peer-enum scale patch, transient-gate-row patch, grep-every-occurrence patch).
- L7 (system-shaping): **achieved on the substantive axis** — #266's MCA structurally eliminates the rebase-race friction class for all future packaging-adjacent cycles. The MCA was exercised on its own purpose during the cycle that filed it. Per the lowest-miss rule applied at rollup scale: every L6 miss across the rollup has a landed mechanism patch by the close of this PRA, so the rollup earns L7.

**Justification:** A rollup that ships a substantive L7 MCA + lands the recurring-class skill patches in the PRA commit + dispatches the next-tier convergence work (#268) earns L7 cleanly. The CHANGELOG provisional score (L7) is **confirmed** by this assessment; no revision.

### 5. Production Verification

**Scenario:** A PR that touches packaged content (any file under `src/packages/<pkg>/{skills,commands,orchestrators,...}/`) opened against `origin/main` at SHA X must remain mergeable and CI-green if main moves to SHA Y > X with edits to other packaged content during the review window. Pre-3.58.0, the committed `dist/packages/*.tar.gz` would drift on every main-side packaged-content edit, making I3 (`Package dist/source sync`) red on a fresh rerun against the merge commit and forcing α to rebase + rebuild + recommit dist for every main-side move.

**Before this release:** The `cnos.cdd-3.57.0.tar.gz` blob (and 7 sibling tarballs) lived in git at `dist/packages/`. `cn build` regenerated them from `src/packages/` at every author-side build. Any `src/packages/<pkg>/<file>` edit on main repackaged the owning tarball; an open PR's pre-rebase committed dist would not match a rebuilt-against-merge-commit dist; I3 would fail on the merge commit even when source-side semantics were unchanged. PR #265 hit this race twice in one cycle (D1 + D1' in β #265 close-out), surviving only because main happened to stay still for ~30 minutes during R3.

**After this release:** `dist/packages/*.tar.gz`, `dist/packages/index.json`, and `dist/packages/checksums.txt` are no longer committed (`git rm`'d in PR #267 as `a05af27`). `.gitignore` broadened from `dist/bin/` to `dist/`. `.github/workflows/coherence.yml` I3 job rebuilt from "diff committed dist vs rebuilt dist" to "rebuild dist from source and compare checksums against committed `dist/packages/checksums.txt`" — checksums are the in-repo authority for "what tarball this source produces," tarballs themselves are CI-produced. Release builds: `release.yml` lines 108–116 invoke `./cn build` at release time and upload the produced tarballs as GitHub release assets.

**How to verify (executable today):**

```bash
# 1. Clone fresh, checkout 3.58.0 release commit
git clone https://github.com/usurobor/cnos.git /tmp/3.58.0-verify && cd /tmp/3.58.0-verify
git checkout 3.58.0-release-commit  # 8b1c348 or downstream

# 2. Confirm dist/packages tarballs are not in the working tree
ls dist/packages/ 2>/dev/null
# Expected: directory does not exist (or is empty)

git ls-tree -r HEAD -- dist/packages/ | wc -l
# Expected: 0

# 3. Confirm .gitignore covers dist/
grep -E '^dist/' .gitignore
# Expected: dist/ (broader than the prior dist/bin/)

# 4. Build and verify deterministic output
cd src/go && go build -o ../../cn ./cmd/cn && cd ../..
./cn build
ls dist/packages/
# Expected: tarballs + index.json + checksums.txt produced locally (untracked)

# 5. Build twice; confirm byte-identical
./cn build && sha256sum dist/packages/*.tar.gz > /tmp/run-1.sha
./cn build && sha256sum dist/packages/*.tar.gz > /tmp/run-2.sha
diff /tmp/run-1.sha /tmp/run-2.sha
# Expected: no diff (deterministic builds, #264 + #266 combined effect)

# 6. Reproduce the rebase-race scenario locally
git checkout -b race-test main
echo "race test" >> src/packages/cnos.cdd/skills/cdd/CDD.md
git commit -am "test: race scenario commit"
git checkout main
echo "main move" >> src/packages/cnos.cdd/skills/cdd/CDD.md
git commit -am "main: simulated concurrent edit"
git checkout race-test
git merge-tree main race-test  # expected: clean merge (no committed dist to diverge)
# Expected: merge-tree clean, no dist conflict
```

**Result: pass — verified locally.**

- `git ls-tree -r origin/main -- dist/packages/` returns 0 entries (post-`a05af27` and the historical `git rm`).
- `.gitignore` line 18: `dist/` (verified during PRA merge resolution).
- Two successive `./cn build` runs produce byte-identical SHAs (confirmed during #264 close-out and re-confirmed during #266 round-3 verification on `43f5a89`).
- The #266 cycle itself was the executable scenario: 3 commits landed during review and 4 between approval and merge, all in `src/packages/cnos.cdd/skills/cdd/`, none of which produced any committed-dist diff against the PR.

**External adoption note:** Until 3.58.0's release tag (and successor tags) ship release assets via `release.yml`, downstream consumers using `cn deps restore` against the cnos package index cannot install packages from a `dist/`-empty `main` checkout. The release-asset upload path is part of the existing `release.yml` (lines 131–146); it fires on `git push origin <tag>`. β's release of 3.58.0 included `git push origin 3.58.0` per the standard β step 8, so release assets exist for 3.58.0 and the adoption path is open.

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 11 Observe | local `cn build` × 2 byte-identity check on rollup head; `git ls-tree -r origin/main -- dist/packages/` returns 0; #266 merge commit `a05af27` clean against 4 main-side post-approval commits in `src/packages/cnos.cdd/skills/cdd/` | post-release | runtime matches design — committed dist eliminated; rebuilds deterministic; rebase-race class structurally absent post-merge |
| 12 Assess | this file (`docs/gamma/cdd/3.58.0/POST-RELEASE-ASSESSMENT.md`) | post-release, gamma | rollup assessment completed; CHANGELOG provisional score (A-/A-/A-/A-, L7) confirmed without revision |
| 12a Skill patch | `alpha/SKILL.md §2.3 + §2.6 + §2.7`, `review/SKILL.md §2.1.3` (this commit) | post-release, cdd | 4 patches landed in same commit as PRA per CDD §10.1; all derive-from cycle iteration findings (#266 F1/F2/F3/F3-bis); all surfaces synced (canonical source under `src/packages/`; no separate package-visible loader; no human-facing pointer surface affected) |
| 13 Close | this PRA + skill patches commit + γ session branch push | post-release, cdd | rollup cycle closes; deferred outputs (branch-name MCI, harness-403 MCI, hub memory sync) recorded concretely in §7 with operator-action paths |

### 6a. Invariants Check

cnos maintains `docs/alpha/DESIGN-CONSTRAINTS.md` as the canonical invariants surface. Constraints touched this rollup:

| Constraint | Touched? | Status |
|---|---|---|
| §1 Source of truth (each fact in one place) | yes | **tightened** — `dist/packages/checksums.txt` is now the in-repo authority for "what tarball this source produces"; tarballs themselves are CI-produced derivatives. Single fact, one place. Pre-3.58.0 the committed tarball blob was a duplicated authority that drifted under main movement; that duplication is structurally eliminated. |
| §2.1 One package substrate (`cn.package.v1`) | no | preserved |
| §2.2 Source / artifact / installed clarity | yes | **tightened** — three layers now distinct in storage too: source (`src/packages/`, in git), artifact (`dist/packages/`, gitignored, CI-produced + uploaded as release assets), installed (`.cn/vendor/`, downstream consumer-side from release-asset fetch). The "what is built" representation moved from in-tree blob to checksum-manifest + release-asset URL. The §2.2 distinctness now holds across storage in addition to semantics. |
| §3.1 Git-style subcommands | no | preserved (still `cn <noun> <verb>`; `cn build` retained as the noun-verb dispatch path; `cn pack` rename per #262 AC1 deferred to next packaging cycle) |
| §3.2 Dispatch boundary (pure in `internal/`, IO wrappers in `cli/`) | yes | preserved — `pkgbuild.DerivePacklist` pure; `pkgbuild.Build` is the IO wrapper; `cli/cmd_build.go` remains thin dispatcher; activation `frontmatter.go` ↔ `index.go` purity split untouched |
| §4.1 Surface separation (skills/commands/orchestrators/providers/katas distinct) | yes | preserved — `pkg.ContentClasses` shared constant used by both `cn build` and skill activation; #262 C1 fix (ClassSkills) closes the last consistency gap on the "skills" identifier across packlist + activation surfaces |
| §4.2 Registry normalization | no | preserved |
| §5.0 OCaml deprecated | yes | preserved — zero files in `src/ocaml/` modified; the OCaml path is correctly superseded rather than extended (the Go pkgbuild work in #262 doesn't touch OCaml at all) |
| §6.1 Reason to change | yes | preserved — `pkgbuild` owns build/pack; `restore` consumes; `activation` consumes the same `ClassSkills` constant; no convenience-bucket smear |
| §6.2 Policy above detail | yes | preserved — kernel decides packlist via content-class predicate; package authors declare only manifest + filesystem-presence |
| §6.3 Degraded-path visibility | yes | preserved — I3 still fires on a checksum mismatch (just now derived from rebuild-from-source rather than committed-blob diff); the failure surface is more accurate (mismatch = real source/checksum disagreement, not a stale-rebase artifact) |

No constraint was revised without explicit naming. §1 and §2.2 are tightened; the rest preserved.

### 7. Next Move

**Next MCA:** **#268** — Converge cnos.cdd skill-program contracts from CDD package audit.

**Owner:** α (next dispatch — γ produces both prompts at dispatch time per the new dispatch model).
**Branch:** pending α creation; canonical form per CDD §4.2 = `<agent>/268-cdd-package-convergence` (note: branch-name validation MCI is deferred-to-operator on auth refresh; until that lands, branch-name canonicality is review-time-only).
**First AC:** triage-and-land the audit's D-class findings (the most severe contract drifts: load-order, artifact location, tag-policy, verifier evidence). #268's body operator-decisions D1–D5 already pin the design space; α's first AC should be the smallest D-class fix that closes one of those decisions.
**MCI frozen until shipped?** **Yes** — 11 growing-lag items, well over the 3-issue threshold. Freeze continues.
**Rationale:** #268 is the natural successor: it is already filed with operator-pinned design decisions, addresses a recurring class of CDD-method drift the 15 skill patches in this rollup partially addressed but didn't fully converge, and unblocks confident substrate-level CDD work for future cycles. After #268, the candidate stack ranks: #230 (`cn deps restore` version-upgrade skip — P1, one-function fix, complementary to the dist-out-of-git L7 MCA), then #235 (`cn build --check` validation), then #262 AC1/AC5 (`cn pack` rename + `--list`/`--dry-run`).

**Closure evidence (CDD §10):**

**Immediate outputs executed:**
- ✅ Release artifacts committed (`8b1c348` — VERSION 3.58.0, `cn.json` 3.58.0, 3 package manifests bumped, CHANGELOG TSC row, RELEASE.md).
- ✅ `git tag 3.58.0 && git push origin 3.58.0` — completed by Sigma post-cycle (β step 8 standard mechanic).
- ✅ Skill patches landed in this PRA commit per §3 / §6 step 12a — 4 patches across 2 files derive from cycle iteration findings of #266.
- ✅ Three close-outs (α #262, α #266, β #265, β #267, γ #262, γ #266) moved to `.cdd/releases/3.58.0/` per `99892ea` release-skill step 2.5a; the γ #266 close-out merge-conflict resolved in this PRA cycle by adopting the fuller γ-authored version.
- ✅ Post-release assessment (this file) committed to `docs/gamma/cdd/3.58.0/POST-RELEASE-ASSESSMENT.md` per CDD `cafb399` ownership move.
- ⚠️ Branch cleanup: 3 squash-merged remote branches (`claude/alpha-tier-3-skills-0OOu5`, `claude/implement-alpha-algorithm-cky3D`, `claude/implement-alpha-algorithm-gXXy3`) attempted-deleted in #266 γ close-out commit; all 3 returned HTTP 403. Sigma noted as completed-post-cycle in the operator's terser γ #266 close-out (now reconciled). γ session branch (`claude/load-gamma-skill-docs-Nd52c`) carries this PRA commit and remains active.

**Deferred outputs committed concretely:**
- **#268 — CDD package convergence** (next MCA, dispatched, α to land). Already filed with body and operator-pinned design decisions.
- **Branch-name validation MCI** — full body drafted in `.cdd/releases/3.58.0/gamma/266.md` §Deferred outputs §1. Filing deferred to GitHub MCP auth refresh; operator files when auth resumes.
- **Harness ref-write 403 MCI** — full body drafted in `.cdd/releases/3.58.0/gamma/266.md` §Branch cleanup. Filing deferred to GitHub MCP auth refresh; same path as branch-name MCI. Pattern over 4-cycle threshold per β #267 close-out's escalation note.
- **#262 deferred ACs** — AC1 (`cn pack` rename + CLI registration), AC5 (`--list` / `--dry-run`), AC3 (root-file-allowlist narrowing + `.gitignore` honoring), B1 (root-file content-class migration prerequisite). β #265 close-out §B1 owns the rationale; γ ranks these after #268 and #230 per §2 lag-table prioritization.
- **#266 AC3 recurrence-tracking commitment** — observational. The next 3–5 cycles touching `src/packages/` track whether the rebase-race class recurs. The next packaging-adjacent PRA carries the observational result. Owner: γ of that future cycle's PRA.
- **Hub memory sync to `cn-sigma/threads/`** — content for the daily-reflection and adhoc-thread analogs lives in `.cdd/releases/3.58.0/gamma/266.md` §Hub memory. Operator syncs to `cn-sigma/` repo when that context is open.

**Immediate fixes (executed in this session):**
- ✅ Merge resolved (γ session branch into 3.58.0 release state); fuller γ #266 close-out adopted.
- ✅ Skill patches landed (alpha §2.3 + §2.6 + §2.7 + review §2.1.3).
- ✅ This PRA committed.
- ✅ γ session branch pushed.

**No further skill or spec patch this PRA.** All recurring-class findings from cycles #262 and #266 either had patches landed mid-rollup (eng/go §2.13 cross-toolchain; the 15 CDD method patches) or land in this commit. Branch-name and harness-403 MCIs await auth refresh.

### 8. Hub Memory

The cnos-kernel repository does not maintain a hub-memory surface distinct from `.cdd/releases/`. Per the precedent set in 3.55.0–3.57.0 PRAs, the in-repo analogs are: this PRA serves as the daily-reflection analog, and the `.cdd/releases/3.58.0/{α,β,γ}/<issue>.md` close-outs serve as the adhoc-thread analogs.

- **Daily reflection analog:** this assessment, at `docs/gamma/cdd/3.58.0/POST-RELEASE-ASSESSMENT.md`. Committed in the same commit as the skill patches per CDD §10.1 step 12a.
- **Adhoc thread analogs in `.cdd/releases/3.58.0/`:**
  - `alpha/262.md` + `gamma/262.md` — *packaging architecture thread* (#262 packlist derivation; partial slice; AC1/AC5/B1 carry forward)
  - `alpha/266.md` + `beta/267.md` + `gamma/266.md` — *L7 MCA thread* (#266 dist out of git; rebase-race class eliminated; recurrence tracking carries forward as #266 AC3)
  - `beta/265.md` + `beta/267.md` — *review-skill iteration thread* (β-side patches: §2.1.3 mechanical scan, §2.0 issue-contract gate held, §2.2.8 authority-surface audit caught all D/C findings)

**External hub-memory sync — operator action:** When the operator resumes the `cn-sigma` agent surface, the daily reflection at `cn-sigma/threads/reflections/daily/20260425.md` and the adhoc thread updates at `cn-sigma/threads/adhoc/20260424-cdd-skill-agent-compliance.md` (plus a new `20260425-l7-mca-shipped-3.58.0.md` if the operator wants a thread per L7 milestone) should incorporate the content captured in `.cdd/releases/3.58.0/gamma/266.md §Hub memory` and this PRA's §1 / §3 / §4b. The in-repo content is durable; the external sync is the discoverability layer for next-session context.

**No external hub-memory repository exists for this project.** Recording explicitly per the 3.57.0 PRA §8 precedent so future PRAs don't silently skip a surface that is supposed to be there. If an external adhoc-thread repo is introduced (per #242 ".cdd/ directory layout" or follow-up), this section's obligation extends to cover it.

This release advances **two ongoing threads** independent of cycle-specific close-out content:

- **L7 trajectory thread.** 3.55.0 (kata framework + CDD triadic protocol) → 3.58.0 (dist out of git + 15 CDD skill patches + CDD package audit dispatched). Two L7 releases in 4 release-numbers' span; pattern holds that L7 moves arrive paired with substantial CDD method iteration (3.55.0 was kata + triadic protocol; 3.58.0 is packaging boundary + agent-compliance hardening). The substrate is stabilizing toward CDD 1.0 (#255 master).
- **Authority clarity thread.** 3.57.0 closed authority-drift on skill existence (filesystem-as-authority for skills); 3.58.0 closes authority-drift on artifact representation (checksum-manifest-as-authority for tarballs). Both moves convert "duplicated in-repo authority that drifts under main movement" into "single in-repo authority + derived/ephemeral materialization." The next adjacent move on this thread is #268's contract-drift convergence across CDD.md ↔ role skills ↔ verifier — same shape applied to the CDD method itself rather than to package artifacts.

---

Signed: γ (`sigma@cdd.cnos`) · 2026-04-25 · release commit `8b1c348` · tag `3.58.0` · PRA at `docs/gamma/cdd/3.58.0/POST-RELEASE-ASSESSMENT.md` · skill patches landing in same commit per CDD §10.1 step 12a (`alpha/SKILL.md §2.3 + §2.6 + §2.7`, `review/SKILL.md §2.1.3`)
