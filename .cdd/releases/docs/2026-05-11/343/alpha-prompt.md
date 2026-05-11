# α Dispatch — Cycle #343

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 343 --json title,body,state,comments
Branch: cycle/343
Tier 3 skills:
  - src/packages/cnos.core/skills/write/SKILL.md
```

## Context for δ

- Allowed tools: `--allowedTools "Read,Write,Bash"`
- Git identity for α commits on this cycle: `Alpha <alpha@cdd.cnos>`
  (AC5 requires the *new* form — `alpha@cnos.cdd.cnos` — to appear in the commit trailers α produces on cycle/343. AC2 resolves the special-case for cnos-itself; per the issue recommendation, the elision form `{role}@cdd.cnos` is canonical for cnos-side actors. So α should use `alpha@cdd.cnos` — which is already the elision form and matches what AC5 expects once AC2 confirms the elision choice.)
- α writes to `cycle/343` only — never to main.
- On completion α sets `status: review-ready` in `.cdd/unreleased/343/self-coherence.md` and commits + pushes to `cycle/343`.

## What α must implement (per issue #343)

AC1: Add a "Git identity for role actors" prescription section to one canonical cdd file (most likely `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` or `src/packages/cnos.cdd/skills/cdd/CDD.md`). The section must:
  - Prescribe `{role}@{project}.cdd.cnos` as the three-level form.
  - Include a 2–4 sentence rationale (DNS hierarchy + cnos-as-protocol-origin).
  - Replace or supersede the cycle #287 `beta@cdd.{project}` doctrine in `cdd/review/SKILL.md`.

AC2: Resolve the cnos self-identity special case. Per issue recommendation: `{role}@cdd.cnos` (elision of `cnos.cdd.cnos`) is canonical for cnos-side actors. Document the rationale (existing trailers match; redundancy adds no information).

AC3: Add a worked example table (3–5 rows) in the same file as AC1:
  - `alpha@tsc.cdd.cnos` — tsc project actor
  - `beta@cdd.cnos` — cnos actor (elision form per AC2)
  - `gamma@<future-project>.cdd.cnos` — hypothetical third project
  - One row showing the deprecated `{role}@cdd.{project}` form, marked "(deprecated)".

AC4: Add identity-migration paragraph to `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md`. ≤80 words. Must name: cycle #343 as cutover; date of merge; history immutable / forward-only switch; transition-window tolerance (mixed trailers on in-flight cycles are acceptable).

AC5: All commits α makes on `cycle/343` in this cycle must use the new form in the git identity — specifically `alpha@cdd.cnos` per AC2 elision. Oracle: `git log --format='%ae' cycle/343` shows only `alpha@cdd.cnos`.

## Patch targets (known sites to update)

1. `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — "Review identity" section (cycle #287 doctrine — `beta@cdd.{project}`). Update or cross-reference to new canonical site.
2. `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — §1.2 or identity setup section if it sets `user.email`. Align to new form.
3. Any `alpha/`, `beta/`, `gamma/` SKILL.md that shows a worked example of the identity form — update examples.
4. `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — add migration paragraph (AC4).

## Non-goals (do not touch)

- Historical commit trailer rewriting.
- Branch naming (`cycle/{N}`).
- Cross-repo identity authentication or GPG signing.
- Any file outside `src/packages/cnos.cdd/skills/cdd/` (plus `post-release/SKILL.md`).

## Dispatch note

δ: dispatch α with `claude -p` against the cnos repo on branch `cycle/343`. Pass the prompt text above. α loads its SKILL.md first and reads the full issue via `gh issue view 343 --json title,body,state,comments`.
