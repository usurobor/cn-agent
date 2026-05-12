---
cycle: 345
role: alpha
status: in-progress
---

# Self-Coherence — Cycle #345

## Gap

**Issue:** cnos: Document the generic α/β/γ/δ/ε role-scope ladder as a cnos-level pattern (cdd, cdw, c-d-X all instantiate)
**Mode:** design-and-build (design in issue body; no prior design doc)
**Version:** CDD 3.15.0

The α/β/γ/δ/ε role system is generic — a scope-escalation ladder where each role's domain is the previous role's frame. This structure applies to any coherence-driven discipline. It is currently documented only as a cdd-internal artifact. This cycle produces `ROLES.md` at repo root as the cnos-level pattern doc, stubs `cdd/epsilon/SKILL.md`, cross-references from `CDD.md`, and re-attributes cdd-iteration to ε.

## Skills

**Tier 1:**
- `cdd/CDD.md` (v3.15.0) — canonical lifecycle, artifact contract, role algorithm
- `cdd/alpha/SKILL.md` — α role surface and execution detail

**Tier 2:**
- (docs-only cycle; no eng/* tier-2 code skills apply)

**Tier 3:**
- `cnos.core/skills/write/SKILL.md` — prose authoring: one governing question per section, front-load the point, cut throat-clearing
- `cnos.core/skills/skill/SKILL.md` — skill-program/frontmatter coherence; for authoring `cdd/epsilon/SKILL.md`

## ACs

**AC1 — `ROLES.md` at repo root, §§1–8 present.**

- Oracle `wc -l ROLES.md` → 319 (within 250–500 ✓)
- Oracle `rg '^## §' ROLES.md` → 8 hits ✓
- §1 table: 5 rows (α/β/γ/δ/ε) with verbs produces/reviews/coordinates/operates/iterates ✓
- Self-application footnote present in §1 ✓
- §2: exactly 3 paragraphs; orders 0–4 named with one-sentence examples each ✓
- §2: Bateson / von Foerster cited ✓
- §3: exactly 6 fields (matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule) ✓
- §4: link to `src/packages/cnos.cdd/skills/cdd/CDD.md` present ✓
- §5: ≤200 words; four one-liners (α-prose, β-clarity/coherence/source-of-truth, γ-cycle, δ-pipeline); "separate issue / future cycle" present ✓
- §6: cdd and cdw named as first two letters ✓
- §7: α/β/γ/δ/ε fixed; verbs fixed; matter and oracles vary ✓
- §8: 8 terms (role, matter, frame, instantiation, scope-escalation, order-of-observation, oracle, cycle) ≥6 ✓
- §§1–3 protocol-agnostic (no instantiation-specific language) ✓
- Placement decision: repo root per issue open question #1 recommendation ✓

**AC2 — `cdd/CDD.md` top-of-file pointer.**

- Oracle `head -20 src/packages/cnos.cdd/skills/cdd/CDD.md | grep ROLES.md` → hit ✓
- One pointer, within first 8 lines ✓
- No claim that role structure is cdd-original ✓

**AC3 — `cdd/epsilon/SKILL.md` stub.**

- Oracle `wc -l epsilon/SKILL.md` → 64 (within 30–100 ✓)
- Frontmatter: name, description, artifact_class, parent, scope, governing_question ✓ (6 fields)
- §1: ε's cdd-side scope — protocol-iteration via cdd-iteration.md + MCA discipline ✓
- §1: cross-references `ROLES.md §1` row 5 and `cdd/post-release/SKILL.md` Step 5.6b ✓
- §2: ε's relationship to δ — often same actor; separation when volume warrants ✓
- Negative invariant: no claim ε is required as a separate human/agent ✓

**AC4 — `post-release/SKILL.md` Step 5.6b re-attribution.**

- Oracle `rg 'ε' src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` → hit at Step 5.6b ✓
- Single opening-sentence amendment: "γ authors … — ε's work product … —" ✓
- No other prose changes in post-release/SKILL.md ✓

**AC5 — Pattern self-applies (γ's close-out obligation, not α implementation).**

- Noted in §Known Debt below; no α action required during implementation ✓

**AC6 — `ROLES.md §5` names cdw as planned-not-shipped.**

- Oracle `rg 'separate issue|future cycle' ROLES.md` → hit in §5 ✓
- §5 ≤200 words ✓
- No `cdw/` directory created in this cycle ✓

## CDD Trace

<!-- α populates during implementation -->

## Known Debt

<!-- α populates during implementation -->
