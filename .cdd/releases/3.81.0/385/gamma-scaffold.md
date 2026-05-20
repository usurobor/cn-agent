<!-- sections: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, empirical-anchor, diff-scope, dispatch-config, tier3-skills] -->
<!-- completed: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, empirical-anchor, diff-scope, dispatch-config, tier3-skills] -->

# γ Scaffold — Cycle #385

**Date:** 2026-05-20
**Issue:** [#385](https://github.com/usurobor/cnos/issues/385) — Streamline activation soul: collapse 6 CA skills → 2 (cap + clp); deprecate cnos.core/AGENTS.md and hub pre-skill scaffolding
**Branch:** `cycle/385`
**Base SHA:** `6a187c6219602fa7f67ae609f0f2aba7947c6d46`
**γ identity:** gamma / gamma@cdd.cnos
**Dispatch config:** §5.1 canonical multi-session (escalation criteria fired: 9 ACs, multi-skill contract surface)

---

## Issue

**Gap:** The `agent/activate/SKILL.md §2.1 step 2` loads 6 CA skills. Four of the six (`mca`, `mci`, `coherent`, `agent-ops`) are partially or fully covered by `cap`; `agent-ops` describes the deprecated OCaml-era daemon runtime. Loading 6 skills on every activation adds context overhead for rules already present in cap and propagates a deprecated surface. The pre-skill scaffolding in `cnos.core/AGENTS.md` references `SOUL.md`/`USER.md` (pre-cycle/379 naming), the OCaml daemon runtime, and an "OPERATIONS mindset" not loaded by activate skill.

**Goal:** Lean soul: KERNEL + cap + clp (3 surfaces). Absorb operational rules from mca/mci/coherent into cap; standalone files kept for on-demand citation. AGENTS.md deleted or replaced with a one-pager. Renderer and kata fixture updated accordingly.

**Priority:** P2 — doctrine cleanup.

**Work-shape:** substantial cycle (9 ACs, multi-surface skill authoring + Go renderer + test fixture).

**Mode:** design-and-build. Design is the issue body (CLP'd at r2 against codex gpt-5.5 high-reasoning; two rounds, converged). No separate DESIGN.md or PLAN.md at a stable path — issue body is the authoritative source. α executes from the ACs; design discovery scope is minimal.

---

## Surfaces γ expects α to touch

1. `src/packages/cnos.core/skills/agent/activate/SKILL.md` — §1.1 (parts), §1.3 (F2 failure mode), §2.1 step 2 (load list 6→2), §3.1 (load order rule), §4.1 (read-first-order block), §4.2 (ca-skills token), §6 (failure modes catalogue), §7 (kata step 4), §8 (references). Also `calls:` frontmatter — drop 4 removed callees.
2. `src/packages/cnos.core/skills/agent/cap/SKILL.md` — absorb rules from mca/mci/coherent; rewrite §4 Boundary to Adjacent Skills.
3. `src/packages/cnos.core/skills/agent/mca/SKILL.md` — frontmatter: add clarification that it is no longer activation-loaded.
4. `src/packages/cnos.core/skills/agent/mci/SKILL.md` — frontmatter: same.
5. `src/packages/cnos.core/skills/agent/coherent/SKILL.md` — frontmatter: same.
6. `src/packages/cnos.core/skills/agent/agent-ops/SKILL.md` — frontmatter: not activation-loaded; document OCaml-era daemon contract scope.
7. `src/packages/cnos.core/AGENTS.md` — delete (preferred per AC5(a)) or replace with one-pager; adjust `pkgbuild/build.go:161,183` and `pkgbuild/build_test.go:428,435,444` if deleted.
8. `src/go/internal/activate/activate.go:505-506` — `ca-skills` token string: update from 6-skill to 2-skill path.
9. `src/packages/cnos.kata/katas/R5-activate/kata.md:115` — update skill list from 6 to 2.
10. Citation audit: `src/packages/cnos.core/doctrine/CA-CONDUCT.md:108`, `src/packages/cnos.core/skills/agent/ca-conduct/SKILL.md:116`, plus any other prose that implies mca/mci/coherent/agent-ops are activation-loaded.

---

## Peer enumeration (§2.2a — grep evidence)

All claims verified before dispatch:

- "activate/SKILL.md references 6-skill set" — confirmed: `grep -n "cap.*clp.*mca\|{cap,clp" activate/SKILL.md` returns ≥15 hits (lines 78, 118-124, 327, 435, 466-468, plus prose at 3, 53, 62, 89, 112, 150, 243, 250).
- "activate.go has ca-skills hardcoded as 6 skills" — confirmed: `grep -n "ca-skills" activate.go` returns line 505: `return "cnos.core/skills/agent/{cap,clp,mca,mci,coherent,agent-ops}/SKILL.md (CA skill set)"`.
- "R5-activate kata.md references 6 skills" — confirmed: `grep -n "cap.*clp.*mca\|CA skill" kata.md` returns line 115: `2) CA skills (cap/clp/mca/mci/coherent/agent-ops)`.
- "AGENTS.md is load-bearing in pkgbuild" — confirmed: `grep -n "AGENTS.md" build.go` returns lines 161, 183; `grep -n "AGENTS.md" build_test.go` returns lines 428, 435, 444.
- "mca/mci/coherent/agent-ops `calls:` references only in activate/SKILL.md" — confirmed: `rg -n 'mca/SKILL\.md|mci/SKILL\.md|coherent/SKILL\.md|agent-ops/SKILL\.md' src/packages` returns only `activate/SKILL.md` (AC3 verified pre-dispatch: no other `calls:` entries in cnos.core/cnos.eng/cnos.cdd).

No gaps: no claim of "X does not exist" asserted without grep evidence.

---

## Scope boundary: #383 vs #385

**cnos#383** (renderer prompt collapse, two-phase): touches `activate/SKILL.md §4.4` (observable-output preservation contract), `activate.go` legacy-section emit deletions, R5-activate **renderer-output** fixture.

**cnos#385** (this cycle): touches `activate/SKILL.md §2.1 step 2` and `§4.2` ca-skills token, `cap/SKILL.md` content, standalone skill frontmatter, AGENTS.md, `activate.go` **ca-skills string** (line 505), R5-activate **skill-list** fixture (kata.md line 115).

**Overlap surface:** Both touch `activate.go` and R5-activate. The concerns are different (renderer prompt structure vs soul-skill path), but the files intersect.

**Dependency note:** #385 should land first. #383 AC2 states "§4.2 interpolation surface still describes how each of the 6 tokens emits its line" — after #385 lands, the `ca-skills` token will reference 2 skills, not 6. #383 will need to reconcile this when it runs. γ will note this in the #383 issue at close-out.

**Cross-repo bundle (activate-foreign-body):** STATUS = `2026-05-19 drafted sigma` — drafted, not submitted. Out of scope per dispatch note. No intake action.

---

## AC oracle approach

α should verify each AC using the oracles embedded in the issue ACs:

- **AC1:** `calls:` frontmatter in activate/SKILL.md has exactly 2 entries: `agent/cap/SKILL.md` and `agent/clp/SKILL.md`; the load list in §2.1 step 2 names 2 skills only.
- **AC2:** `cap/SKILL.md` contains absorbed rules from mca lines 15-19, 34-54, 60-67; mci lines 25-37, 58-69; coherent lines 21-39, 45-50, 52-64, 70-99, 103-112. Verify by grep that the content is present in cap; verify by inspection that the non-absorbed content is named as "remains in standalone file."
- **AC3:** mca, mci, coherent, agent-ops standalone files exist on disk; their frontmatter states they are not activation-loaded; their `calls:` entries from activate/SKILL.md are removed. `grep -rn "mca/SKILL.md\|mci/SKILL.md\|coherent/SKILL.md\|agent-ops/SKILL.md" src/packages --include="*.md"` returns zero `calls:` references.
- **AC4:** activate/SKILL.md §1.1, §1.3, §3.1, §4.1, §4.2, §6, §7, §8 — all 6-skill references collapsed to 2. Verify by grep; no instance of the 6-item list remains.
- **AC5:** AGENTS.md deleted (preferred) — `ls src/packages/cnos.core/AGENTS.md` returns non-zero; `go test ./src/go/internal/pkgbuild/...` passes with adjusted test.
- **AC6:** `activate.go:505` returns `cnos.core/skills/agent/{cap,clp}/SKILL.md` (2-skill path); `go test ./src/go/internal/activate/...` passes.
- **AC7:** `R5-activate/kata.md` line that listed 6 skills now lists 2; `src/packages/cnos.kata/katas/R5-activate/run.sh` passes.
- **AC8:** Prose citation audit: `rg -rn 'mca/SKILL\.md|mci/SKILL\.md|coherent/SKILL\.md|agent-ops/SKILL\.md' src/packages` — every match is in a context describing these files as on-demand/reference, not activation-loaded; any "loaded at activation" prose is rewritten.
- **AC9:** `cap/SKILL.md §4 Boundary to Adjacent Skills` rewritten to reflect: CAP owns activation-loaded coherent-output basics; CLP remains adjacent; mca/mci/coherent remain on-disk; agent-ops removed from adjacent-skill mention.

**Validation script:** `scripts/validate-skill-frontmatter.sh` (or equivalent) must pass; `go test ./src/go/...` must pass.

---

## Empirical anchor

No specific "CDD tsc #36 §Gap" style anchor for this cycle — the gap is self-evident from activation-session friction. CLP r2 convergence provides the design-quality baseline.

---

## Expected diff scope

| Surface | Expected delta |
|---|---|
| `activate/SKILL.md` | ~40-80 line edit (§2.1 step 2 reduction, §4.2 token, 6 reference collapses, `calls:` frontmatter) |
| `cap/SKILL.md` | ~80-120 line expansion (absorb from 3 skills) + §4 rewrite (~10 lines) |
| `mca/SKILL.md`, `mci/SKILL.md`, `coherent/SKILL.md`, `agent-ops/SKILL.md` | ~5-10 line frontmatter addition each |
| `AGENTS.md` | deleted (preferred) |
| `pkgbuild/build.go` | minor (if AGENTS.md deleted: remove comment reference) |
| `pkgbuild/build_test.go` | minor (if AGENTS.md deleted: remove AGENTS.md from expected-root-files list) |
| `activate.go` | 1 line change (ca-skills string) |
| `R5-activate/kata.md` | 1 line change |
| `CA-CONDUCT.md`, `ca-conduct/SKILL.md` | ~2 lines each (prose describing as on-demand) |

Total: ~200-300 line net change across 10-12 files. Mostly Markdown. One Go line change.

---

## Dispatch configuration

**§5.1 canonical multi-session.** Escalation criteria fired:
- ≥7 ACs (9 ACs present)
- Multi-skill contract surface (activate skill §2.x/§4.x, cap skill §3+§4, mca/mci/coherent/agent-ops frontmatter, cnos.core/AGENTS.md, renderer activate.go, R5-activate kata fixture)

γ/δ separation structurally present. Each role runs as a separate `claude -p` process.

---

## Tier 3 skills

Named explicitly for dispatch:
- `src/packages/cnos.core/skills/write/SKILL.md` — always included (every α output is a written artifact)
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-file authoring; modifying ≥6 SKILL.md files
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go engineering; activate.go line 505 + pkgbuild test changes

Note: The issue body does not have a formal "## Skills to load" section. Tier 3 skills are derived from the issue's impact graph. This is a minor issue-quality gap (γ observation only — not a gate-blocker given CLP r2 convergence and explicit surface references throughout the ACs).

---

## Pre-dispatch check result

```
γ pre-dispatch check — 2026-05-20:
  origin/cycle/385 exists: YES (just created)
  .cdd/unreleased/385/gamma-scaffold.md exists on cycle/385: YES (this file, will be pushed)
  issue #385 is open: YES
  branch pre-flight: PASS
  peer enumeration: PASS — all 6 referenced surfaces confirmed by grep
  scope boundary with #383: DOCUMENTED (different sections; #385 before #383)
  cross-repo intake: no submitted proposals pending
  issue quality gate: PASS (minor: no formal Tier 3 section in issue body; compensated in dispatch)
  dispatch config: §5.1 confirmed
  timeout budget: 1800s (code cycle, 9 ACs: max(400, 180×9) = 1620s, rounded to 1800s)
```
