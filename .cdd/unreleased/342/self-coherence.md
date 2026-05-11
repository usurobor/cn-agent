---
cycle: 342
role: alpha
status: in-progress
---

# Self-Coherence — Cycle #342

## Gap

**Issue:** cdd/operator: Add §5 — Dispatch configurations (single-session δ-as-γ via Agent tool, Claude Code activation)
**Mode:** design-and-build
**Version:** CDD 3.15.0

`cdd/operator/SKILL.md` names only the canonical `claude -p` multi-session dispatch model. The single-session δ-as-γ configuration (one parent Claude Code agent dispatching α and β as sub-agents) is in active use but ungoverned — no canonical text names what is preserved, what is lost, or what grading implications follow. This cycle adds `operator/SKILL.md §5 Dispatch configurations` and amends `release/SKILL.md §3.8` with a configuration-floor clause.

## Skills

**Tier 1:**
- `cdd/CDD.md` (v3.15.0) — canonical lifecycle, artifact contract, role algorithm
- `cdd/alpha/SKILL.md` — α role surface and execution detail

**Tier 2:**
- (docs-only cycle; no eng/* tier-2 code skills apply)

**Tier 3:**
- `cnos.core/skills/write/SKILL.md` — prose authoring: one governing question per section, front-load the point, cut throat-clearing
- `cnos.core/skills/skill/SKILL.md` — skill-program/frontmatter coherence; no frontmatter changes allowed per issue constraint

## ACs

_Populated in §ACs section below._

## Self-check

_Populated after implementation._

## Debt

_Populated after implementation._

## CDD-Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| scaffold | .cdd/unreleased/342/self-coherence.md | — | cycle branch created by γ |
| α intake | self-coherence.md §Gap + §Skills | Tier 1 + Tier 3 | loaded; no ambiguity; proceeding |
