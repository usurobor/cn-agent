<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs, Self-check, Debt, CDD Trace] -->

# Self-Coherence — Cycle #385

## §Gap

**Issue:** [#385](https://github.com/usurobor/cnos/issues/385) — Streamline activation soul: collapse 6 CA skills → 2 (cap + clp); deprecate cnos.core/AGENTS.md and hub pre-skill scaffolding

**Version / mode:** 3.15.0 · design-and-build. Issue body is the authoritative design (CLP r2 converged). No separate DESIGN.md or PLAN.md. α executes from the 9 ACs.

**Gap:** `agent/activate/SKILL.md §2.1 step 2` loads 6 CA skills. Four of the six (`mca`, `mci`, `coherent`, `agent-ops`) are partially or fully covered by `cap`; `agent-ops` describes the deprecated OCaml-era daemon runtime. Loading 6 skills on every activation adds context overhead for rules already present in cap and propagates a deprecated surface. The pre-skill scaffolding in `cnos.core/AGENTS.md` references `SOUL.md`/`USER.md` (pre-cycle/379 naming), the OCaml daemon runtime, and an "OPERATIONS mindset" not loaded by activate skill.

**Goal:** Lean soul — KERNEL + cap + clp (3 surfaces). Absorb operational rules from mca/mci/coherent into cap. Standalone files kept for on-demand citation. Renderer and kata fixture updated accordingly. AGENTS.md deleted.

## §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2 (always-applicable engineering):**
- `src/packages/cnos.eng/skills/eng/README.md` (not loaded separately; Tier 3 covers the relevant surface)

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/write/SKILL.md` — prose and skill authoring quality
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-file authoring (modifying ≥6 SKILL.md files)
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go engineering (activate.go + pkgbuild changes)

## §ACs

Per-AC oracles run against branch HEAD (implementation SHA: `382a978a`).

**AC1 — calls: frontmatter collapsed to 2:**
Evidence: `grep -A5 "^calls:" src/packages/cnos.core/skills/agent/activate/SKILL.md` → 3 entries: `../doctrine/KERNEL.md`, `agent/cap/SKILL.md`, `agent/clp/SKILL.md`. §2.1 step 2 lists cap and clp only, with prose noting mca/mci/coherent/agent-ops remain on-disk for on-demand load. ✓

**AC2 — cap/SKILL.md absorbs MCA/MCI/coherent-output rules:**
Evidence: cap/SKILL.md expanded from 309 to 464 lines. Added:
- §4 MCA rules: "Minimum" = scope not speed; workaround failure mode; §4.1 Scope (solve/patch/structural-over-behavioral); §4.2 Act (do it/surface it); §4.3 Quick Test.
- §5 MCI rules: §5.1 Identify (transferable-lesson tests); §5.2 Capture (immediately/imperative/trigger); §5.3 Migrate (Daily→Weekly→Skill).
- §6 Coherent output: §6.1 parts/sources/dependencies+failure-modes; §6.2 L0/L1/L2 levels; §6.3 verification+coherence-vs-correctness; §6.4 rules (terms/trace/resolve/propagate/verifiable/orphans); §6.5 pre-ship checklist.
Non-absorbed content explicitly named: root-finding sequence stays in mca (§4 note); scope rules + MCI vs MCA table stay in mci (§5 note); Coherence Modes + Anti-Patterns table stay in coherent (§6 note). ✓

**AC3 — mca/mci/coherent/agent-ops frontmatter updated:**
Evidence: each file's frontmatter has `activation_status:` field declaring not activation-loaded since cycle/385, with pointer to what was absorbed.
`rg -rn 'mca/SKILL\.md|mci/SKILL\.md|coherent/SKILL\.md|agent-ops/SKILL\.md' src/packages --include="*.md"` → 9 hits; all describe these as on-demand (ca-conduct/SKILL.md, CA-CONDUCT.md use mci as on-demand; activate/SKILL.md §8 lists all 4 as "on-demand reference"; cap/SKILL.md §4/§5/§6 notes explicitly name them as on-disk for on-demand detail). No hit says "loaded at activation." ✓

**AC4 — all 6-skill references in activate/SKILL.md collapsed to 2:**
`grep -n "cap,clp,mca\|mca,mci\|mci,coherent\|coherent,agent-ops" src/packages/cnos.core/skills/agent/activate/SKILL.md` → 0 hits.
Updated sites: calls: frontmatter (l.31-34), §1.1 CA skills path, §2.1 step 2 list, §4.2 ca-skills token, §7 kata step 4, §8 references. ✓

**AC5 — AGENTS.md deleted; pkgbuild tests pass:**
`ls src/packages/cnos.core/AGENTS.md` → exit 2 (not found).
`/usr/local/go/bin/go test ./internal/pkgbuild/...` → ok (cached).
build.go comments updated to remove AGENTS.md reference (lines 161, 183).
build_test.go TestDerivePacklistIncludesRootFiles updated: AGENTS.md removed from file list and assertion list. ✓

**AC6 — activate.go returns 2-skill path; tests pass:**
`grep "ca-skills" src/go/internal/activate/activate.go` → line 506: `return "cnos.core/skills/agent/{cap,clp}/SKILL.md (CA skill set)"`.
`/usr/local/go/bin/go test ./internal/activate/...` → ok. ✓

**AC7 — R5-activate/kata.md updated:**
`grep "CA skills" src/packages/cnos.kata/katas/R5-activate/kata.md` → line 115: `2) CA skills (cap/clp)`. ✓

**AC8 — citation audit passes:**
`rg -rn 'mca/SKILL\.md|mci/SKILL\.md|coherent/SKILL\.md|agent-ops/SKILL\.md' src/packages --include="*.md"` — every match is in a context describing these as on-demand/reference. No prose says or implies "loaded at activation."
CA-CONDUCT.md:108 and ca-conduct/SKILL.md:116 reference mci as on-demand use ("when the result is a behavior-changing insight") — no change required.
`bash tools/validate-skill-frontmatter.sh` → ✓ 67 SKILL.md validated; no findings. ✓

**AC9 — cap/SKILL.md §7 Boundary rewritten:**
§7 (was §4) now states: CAP owns activation-loaded rules (§1-§6). CLP remains separately loaded. mca/mci/coherent remain on-disk for on-demand load. agent-ops removed from adjacent-skill mention (runtime-specific, not behavioral). ✓

## §Self-check

**Did α's work push ambiguity onto β?** No. Every AC has concrete evidence (grep output, test pass/fail, file existence check). The absorption boundary in cap/SKILL.md explicitly names what was absorbed and what remains in each standalone file — β can verify any absorption claim by reading the named section of the standalone skill. The non-absorbed content note in cap §4/§5/§6 is the boundary document.

**Is every claim backed by evidence in the diff?**
- AC1: calls: frontmatter is in the diff; §2.1 step 2 list is in the diff.
- AC2: cap/SKILL.md §4/§5/§6 are in the diff; non-absorbed notes name the standalone files.
- AC3: mca/mci/coherent/agent-ops frontmatter `activation_status:` fields are in the diff.
- AC4: The five changed sites in activate/SKILL.md are in the diff; grep confirms 0 remaining hits.
- AC5: AGENTS.md deletion is in the diff; pkgbuild diffs show comment and test removal.
- AC6: activate.go line 506 is in the diff; tests pass.
- AC7: kata.md line 115 is in the diff.
- AC8: No changes required to CA-CONDUCT.md or ca-conduct/SKILL.md; frontmatter validation passes.
- AC9: cap/SKILL.md §7 is in the diff.

**Peer enumeration:** The family of 6-skill references was enumerated before coding (per gamma-scaffold §Peer enumeration). All 5 sites in activate/SKILL.md were changed. No calls: reference to mca/mci/coherent/agent-ops exists in cnos.core/cnos.eng/cnos.cdd (verified by grep during AC3). ✓

**Harness audit:** The Go renderer (activate.go) was updated. The kata fixture (kata.md) was updated. The pkgbuild build.go comment and test were updated. No other harness surfaces write the ca-skills string. ✓

## §Debt

- **Hub-side downstream cleanup is out of scope** (explicit per gamma-scaffold §Scope boundary). cn-sigma's RULES.md, spec/AGENTS.md, spec/TOOLS.md, spec/HEARTBEAT.md cleanup is a separate per-hub issue downstream of this cnos change.
- **cap/SKILL.md section number shift (§4→§7 for Boundary):** Existing prose in other skills that references "cap §4" by number will now point at MCA rules instead of the Boundary section. A grep found no external prose references to "cap §4" in the current codebase, so no live cross-references are broken. The shift is tracked here as minor debt for completeness.
- **Provisional close-out:** α writes `alpha-closeout.md` at review-readiness time per §2.8 fallback, explicitly marked `[provisional — pending β outcome]`.
- γ-artifact: `gamma-scaffold.md` present at canonical §5.1 path `origin/cycle/385:.cdd/unreleased/385/gamma-scaffold.md`. ✓

## §CDD Trace

| Step | Description | Status |
|------|-------------|--------|
| 1. Receive | Issue #385 read; gamma-scaffold.md loaded; active skill set explicit (Tier 1: CDD.md + alpha/SKILL.md; Tier 3: write, skill, go) | ✓ |
| 2. Produce — design | Design is issue body (CLP r2 converged); no separate DESIGN.md required | ✓ |
| 3. Produce — coherence contract | self-coherence.md §Gap written first | ✓ |
| 4. Produce — plan | Not required; AC sequence is direct; implementation sequencing is trivial | ✓ |
| 5. Produce — tests | Go tests: `go test ./...` → all 13 packages pass. `validate-skill-frontmatter.sh` → 67 SKILL.md validated, no findings. | ✓ |
| 6. Produce — code | Files changed: `activate/SKILL.md`, `cap/SKILL.md`, `mca/SKILL.md`, `mci/SKILL.md`, `coherent/SKILL.md`, `agent-ops/SKILL.md`, `AGENTS.md` (deleted), `activate.go`, `kata.md`, `build.go`, `build_test.go`. All callers of new/changed interfaces verified. | ✓ |
| 7. Produce — docs | `activate/SKILL.md` §8 References updated. `cap/SKILL.md` self-documents absorbed and non-absorbed content. This self-coherence.md is the primary branch artifact. | ✓ |

**Files in `git diff --stat origin/main..HEAD`** (as required by pre-review gate row 11):
- `.cdd/unreleased/385/gamma-scaffold.md` — γ-artifact (not α-authored; present as γ created the cycle branch) ✓
- `.cdd/unreleased/385/self-coherence.md` — this file ✓
- `src/packages/cnos.core/skills/agent/activate/SKILL.md` — AC1+AC4 ✓
- `src/packages/cnos.core/skills/agent/cap/SKILL.md` — AC2+AC9 ✓
- `src/packages/cnos.core/skills/agent/mca/SKILL.md` — AC3 ✓
- `src/packages/cnos.core/skills/agent/mci/SKILL.md` — AC3 ✓
- `src/packages/cnos.core/skills/agent/coherent/SKILL.md` — AC3 ✓
- `src/packages/cnos.core/skills/agent/agent-ops/SKILL.md` — AC3 ✓
- `src/packages/cnos.core/AGENTS.md` — AC5 (deleted) ✓
- `src/go/internal/activate/activate.go` — AC6 ✓
- `src/packages/cnos.kata/katas/R5-activate/kata.md` — AC7 ✓
- `src/go/internal/pkgbuild/build.go` — AC5 ✓
- `src/go/internal/pkgbuild/build_test.go` — AC5 ✓

All 12 changed files mentioned in §ACs or this trace. Pre-review gate row 11 satisfied.

**Caller-path trace for new content (row 12):** cap/SKILL.md §4-§6 content is activation-loaded (called by activate/SKILL.md step 2, which is fetched by any body following the skill). The non-test caller is `activate/SKILL.md §2.1 step 2` referencing `agent/cap/SKILL.md`. ✓
