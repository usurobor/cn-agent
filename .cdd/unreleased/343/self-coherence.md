---
cycle: 343
issue: "cdd: Canonical git identity convention for cdd role actors ({role}@{project}.cdd.cnos)"
branch: cycle/343
mode: design-and-build
status: dispatched
---

# Self-Coherence — Cycle #343

## Gap

The cdd identity convention `{role}@cdd.{project}` inverts the DNS namespace hierarchy. `cdd` is a protocol hosted by cnos, not a subdomain of the tenant project. The correct three-level form is `{role}@{project}.cdd.cnos`.

## Mode

`design-and-build` — design is fixed in the issue body; α implements the patch in this cycle.

## ACs

1. AC1 — Three-level identity convention named in one canonical cdd file.
2. AC2 — Special case for cnos resolved (elision form `{role}@cdd.cnos` recommended).
3. AC3 — Worked example table present (3–5 rows, deprecated form annotated).
4. AC4 — Migration paragraph in `cdd/post-release/SKILL.md`.
5. AC5 — Patch-landing cycle uses new form in its own commit trailers.

## CDD Trace

| Step | Action | Skill | Decision |
|------|--------|-------|----------|
| 1 | Cycle branch created | gamma/SKILL.md §2.5a | cycle/343 from origin/main |
| 2 | Scaffold created | gamma/SKILL.md §2.5a | .cdd/unreleased/343/ |
| 3 | α prompt written | gamma/SKILL.md §2.5b | dispatched to δ |
| 4 | β prompt written | gamma/SKILL.md §2.5b | dispatched to δ |

## Review-readiness signal

α will set `status: review-ready` in this file when implementation is complete and β may begin review.

## Fix rounds

R0 — initial dispatch.
