# cdd-iteration — cycle #398

**Author:** γ (collapsed on δ)
**Date:** 2026-05-21
**Cycle:** cnos#398 — Phase 4b of cnos#366; harness substrate

## Findings

### F1: delta/SKILL.md (Phase 4a) carries stale forward-references to Phase 4b

**Class:** `cdd-skill-gap`
**Surface:** `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` lines 59, 61, 81, 173, 316, 336, 340.

**What:** Phase 4a (cycle/397) authored delta/SKILL.md before Phase 4b (this cycle) shipped harness/SKILL.md. The 4a doctrine, written when 4b was pending, says things like:
- "harness mechanics" live in operator/SKILL.md
- "Phase 4b — harness substrate (pending)"
- "the harness mechanics currently living in `operator/SKILL.md` ... are slated to relocate into a harness substrate surface in a subsequent cycle"

Post-cycle-#398 these references are factually stale: harness/SKILL.md exists; the mechanics moved.

**Why this cycle didn't patch:** Issue cnos#398's scope explicitly forbids touching δ-role boundary content (Phase 4a — cycle/397). The strict scope-control rule kept this cycle from sweeping delta/SKILL.md's prose.

**Disposition:** next MCA — file as part of Phase 5 (γ shrink) sweep, or a dedicated cycle that sweeps "Phase 4b (pending)" / "harness mechanics currently in operator" prose across the doctrine. The patch is small (~6–8 line edits in delta/SKILL.md) and mechanical.

**Concrete next MCA:** Phase 5 (γ shrink) cycle's selection should include a `gamma-scaffold.md` row noting delta/SKILL.md stale-reference sweep as in-scope. First AC: `rg "Phase 4b.*pending|harness mechanics.*operator" src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` returns 0 hits.

### F2: delta/SKILL.md frontmatter `requires.1` CUE schema type mismatch (pre-existing)

**Class:** `cdd-tooling-gap`
**Surface:** `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` frontmatter.

**What:** `bash tools/validate-skill-frontmatter.sh` reports `requires.1: conflicting values string and {or:"active γ dispatch prompt awaiting δ's pre-routing enrichment"} (mismatched types string and struct)`. The schema (`schemas/skill.cue`) expects each `requires` item to be a string; one item appears to have been authored as a YAML mapping (`or: "..."` form) instead of a flat string.

**Why this cycle didn't patch:** Pre-existing on origin/main (introduced by Phase 4a, cycle/397). Scope-control rule prevents touching Phase 4a content; the frontmatter issue is non-blocking (the doctrine is still readable; only the schema validator complains).

**Disposition:** no-patch in this cycle; surface to Phase 5 or a frontmatter sweep cycle. Note for triage: the fix is changing the YAML mapping form to a flat string (likely something like `- "active γ dispatch prompt awaiting δ's pre-routing enrichment" (or alt-form)`); a domain reader of delta/SKILL.md should make the call about the intended meaning.

**Concrete next MCA:** sweep delta/SKILL.md frontmatter to match `schemas/skill.cue`. First AC: `bash tools/validate-skill-frontmatter.sh` reports 0 findings on `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (15 remaining are pre-existing CDR-side findings out of scope).

## Patches landed in this cycle

None. All cycle output is the harness/SKILL.md authoring + cross-reference edits, not skill-gap patches per se.

## No-patch findings

None.

## Summary

- Findings: 2 (`cdd-skill-gap` ×1, `cdd-tooling-gap` ×1)
- Patches: 0
- MCAs (concrete next MCA committed): 2 (F1 and F2 — both deferred to Phase 5 γ shrink)
- No-patch: 0
