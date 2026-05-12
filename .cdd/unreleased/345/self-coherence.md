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

<!-- α populates with evidence during implementation -->

## CDD Trace

<!-- α populates during implementation -->

## Known Debt

<!-- α populates during implementation -->
