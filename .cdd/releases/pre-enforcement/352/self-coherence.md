---
cycle: 352
issue: "#352"
type: self-coherence
date: 2026-05-12
---

# Self-Coherence — Cycle #352

## §Gap

**Issue:** #352 — cdd/review+gamma: Add CI-green gate — β refuses merge on red CI, γ PRA verifies post-merge

**Version/Mode:** docs-only, design-and-build

**Selected gap:** CDD protocol has no CI-green verification mechanism. β can APPROVE on red CI; γ can close-out without verifying post-merge CI status. This creates a latent failure mode where structurally correct workflows (readable by β) can still fail at runtime, and neither role catches it. Issue emerged from tsc #36 cycle which shipped a CI gate but had no mechanism for verifying its *own* CI ran green post-merge.

**Coherence target:** Add symmetric CI-green gates on both sides of merge:
1. β refuses APPROVED without CI green on review SHA  
2. γ PRA polls CI post-merge before close-out

**Mode justification:** docs-only appropriate because this is rule addition to existing skills, not new code surfaces. design-and-build because it requires coordination across 4 skill files with a cohesive verdict-rule + close-out-rule pair.

## §Skills

**Tier 1:** CDD.md (canonical lifecycle), alpha/SKILL.md (α role surface)  
**Tier 2:** eng/README.md (engineering practices)  
**Tier 3:** write/SKILL.md (docs authoring discipline)

## §Design

**Not required** — rule addition to existing skill surfaces, not new architecture. The design is explicit in issue ACs: symmetric CI-green gates (β refuses APPROVED on red CI; γ polls post-merge). No impact graph needed for rule additions within established skill boundaries.

## §Plan  

**Required** — 4-file coordination with dependency order:

1. **review/SKILL.md** — add β verdict gate rule (new rule 3.14 or amend existing 3.12)
2. **gamma/SKILL.md** — add post-merge CI verification step to close-out process  
3. **post-release/SKILL.md** — add PRA template row + §9.1 trigger amendment
4. **release/SKILL.md** — add §3.8 grade cap rules for CI violations

Order rationale: review rules establish the pre-merge gate; gamma/post-release establish the post-merge discipline; release establishes grading consequences. All four must be consistent on "required workflow" definition.

## §ACs

**AC1 — β refuses APPROVED on red/pending CI.**
- *Evidence:* `cnos:src/packages/cnos.cdd/skills/cdd/review/SKILL.md` rule 3.10 replaced with binding CI-green gate
- *Oracle:* `rg 'CI.*green.*binding|gh run list|ci-status' cnos:src/packages/cnos.cdd/skills/cdd/review/SKILL.md` returns hits in verdict rules
- *Status:* ✅ MET — rule 3.10 now mandates `gh run list --branch <review-SHA>` check before APPROVED; RC on red/pending CI

**AC2 — γ post-merge CI verification mandatory.**  
- *Evidence:* `cnos:src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.7 + `cnos:src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` PRA template
- *Oracle:* `rg 'Post.*merge.*verification|post.*merge.*CI' cnos:src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md cnos:src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` returns hits
- *Status:* ✅ MET — gamma/SKILL.md §2.7 adds mandatory post-merge CI verification before close-out; PRA template carries CI status row

**AC3 — §3.8 grade caps named.**
- *Evidence:* `cnos:src/packages/cnos.cdd/skills/cdd/release/SKILL.md` §3.8 CI-red cap clause added
- *Oracle:* `rg 'CI.*red|CI.*verification' cnos:src/packages/cnos.cdd/skills/cdd/release/SKILL.md` returns hit in §3.8 grading section
- *Status:* ✅ MET — §3.8 CI-red cap clause: red CI = γ axis max C; unverified CI = γ axis max B−

**AC4 — §9.1 trigger amended.**
- *Evidence:* `cnos:src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` §4b trigger list updated
- *Oracle:* `rg 'CI.*red.*merge|post.*merge.*red' cnos:src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` returns hit in §9.1 list
- *Status:* ✅ MET — §9.1 trigger list includes "CI red on merge commit (post-merge)"

**AC5 — Self-applied on patch-landing cycle.**
- *Invariant:* This cycle demonstrates the new discipline in its own β review and γ close-out
- *Status:* 🔄 PENDING — will be verified when β reviews this branch and γ closes the cycle

## §Self-check

**Did α's work push ambiguity onto β?**

No. Each AC maps to specific, greppable evidence. The changes are rule additions/amendments to existing skill surfaces within established boundaries. β can verify each AC mechanically using the provided oracle patterns.

**Is every claim backed by evidence in the diff?**

Yes:
- Issue states "β refuses APPROVED without CI green" → review/SKILL.md rule 3.10 implements binding CI-green gate
- Issue states "γ PRA polls CI post-merge" → gamma/SKILL.md adds mandatory post-merge verification + PRA template updated
- Issue states "§3.8 grade caps" → release/SKILL.md CI-red cap clause added with specific grade limits
- Issue states "§9.1 trigger amendment" → post-release/SKILL.md trigger list includes "CI red on merge commit"

All four surfaces updated consistently with "required workflow" definition pattern (GitHub branch protection → every workflow on cycle branch).

## §Debt

**Known debt:** None. Rule additions are complete as specified in issue ACs. 

**Future considerations (not blocking for this cycle):**
- MCP integration for CI polling (currently documented as manual `gh run list`)
- `.cdd/CI-REQUIRED` file mechanism for repos without GitHub branch protection
- Cross-protocol inheritance to cdw and future c-d-X protocols

## §CDD-Trace

**CDD Trace through Step 7:**

1. **Issue creation** — #352 created by γ wave-2026-05-12 finding #6
2. **Dispatch** — α checked out cycle/352, configured git identity as alpha@cnos.cdd.cnos
3. **Gap** — CI-green gate absence in CDD protocol (latent failure: structurally correct workflows can fail at runtime)
4. **Mode** — docs-only, design-and-build (4-file coordination)
5. **Active skills** — CDD.md + alpha/SKILL.md + write/SKILL.md + eng/README.md
6. **Artifacts produced:**
   - `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — rule 3.10 replaced with binding CI-green gate
   - `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — §2.7 post-merge CI verification added
   - `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — PRA template CI status row + §9.1 trigger amendment  
   - `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — §3.8 CI-red grade cap clause added
   - `.cdd/unreleased/352/self-coherence.md` — this file
7. **Self-coherence** — All 5 ACs mapped to evidence; no ambiguity pushed to β; all claims backed by diff

## Review-readiness | round 1 | base SHA: 27b6520c | head SHA: 3b9fe770 | branch CI: no remote CI configured for cycle branches | ready for β

## Fix-round | round 2 | β R1 → RC (1 finding)

**β R1 verdict:** REQUEST CHANGES  
**Finding addressed:** F1 (B-severity, contract) — CI-green gate logical gap  

**F1 resolution:** Updated `.github/workflows/build.yml` triggers to include `'cycle/*'` branches alongside `main`. 
- **Commit SHA:** `6f211798`  
- **Incoherence fixed:** CI-green gate rule 3.10 was structurally void because no workflows executed on cycle branches  
- **Implementation:** Push triggers now include `branches: [main, 'cycle/*']` per β recommendation Option 1  
- **Gate now operative:** Cycle branches will have required CI runs for rule 3.10 to verify  

**Re-verification needed:** β should confirm CI runs on this head commit demonstrate the gate is now operative.

## Review-readiness | round 2 | base SHA: a5865e4e | head SHA: 6f211798 | branch CI: triggered on push | ready for β