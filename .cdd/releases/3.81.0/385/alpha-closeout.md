[provisional — pending β outcome]

# α Close-out — Cycle #385

**Issue:** #385 — Streamline activation soul: collapse 6 CA skills → 2 (cap + clp); deprecate cnos.core/AGENTS.md
**Branch:** cycle/385
**α identity:** alpha / alpha@cdd.cnos
**Date:** 2026-05-20

---

## Summary

Collapsed the activation soul from 6 CA skills to 2 (cap + clp). Absorbed operational rules from mca/mci/coherent into cap/SKILL.md. Marked mca/mci/coherent/agent-ops as on-demand (not activation-loaded). Deleted AGENTS.md. Updated renderer and kata fixture. All 9 ACs satisfied; Go tests pass; frontmatter validation passes.

---

## Friction log

- **Issue design quality:** The issue body is well-structured (CLP r2 converged). No design discovery needed; all 9 ACs map directly to implementable changes. The minor gap (no formal Tier 3 section in issue body) was compensated by the gamma-scaffold.
- **calls: frontmatter vs prose references:** The distinction between activation-loaded (calls:) and on-demand (prose citation) was the core discipline required throughout. The absorption boundary notes in cap §4/§5/§6 make this explicit for every reader.
- **cap section numbering:** Adding three new major sections (§4 MCA, §5 MCI, §6 Coherent output) shifted the prior §4-§7 to §8-§10. This is load-bearing structural change to the primary activation-loaded behavioral skill. No external prose references to cap section numbers were found in the codebase, so no cross-reference rot resulted.
- **AGENTS.md deletion:** Clean. The pkgbuild test was easy to update; the build.go comment references were localized to two sites.

---

## Engineering observations

- Pattern: activation-loaded skill bloat is a recurring failure mode. The six-skill list in activate/SKILL.md had grown incrementally; each new skill was reasonable at add time, but collectively they doubled activation context. The absorb-into-cap approach concentrates the behavioral rules into one scan instead of six.
- The `activation_status:` frontmatter field added to the four standalone files is a new convention. It could be formalized as a standard frontmatter field for any skill that was once activation-loaded and was later downgraded to on-demand.
- The non-absorbed boundary notes in cap §4/§5/§6 follow the same pattern as `activate/SKILL.md §2.4` (disambiguation from cdd/activation). This is the correct form: name the boundary concretely, say what stays where, point to the other file.

---

## Known debt carried forward

- Hub-side downstream cleanup (cn-sigma RULES.md, spec/AGENTS.md, spec/TOOLS.md, spec/HEARTBEAT.md) is out of scope; tracked in #385 issue body under "Hub-side downstream."
- Close-out is provisional pending β merge. α will be re-dispatched for final close-out if β merge occurs.
