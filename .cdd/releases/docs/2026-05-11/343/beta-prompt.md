# β Dispatch — Cycle #343

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 343 --json title,body,state,comments
Branch: cycle/343
```

## Context for δ

- Allowed tools: `--allowedTools "Read,Write"` (β must not tag, release, or run arbitrary commands)
- Git identity for β commits on this cycle: `Beta <beta@cdd.cnos>` (elision form per AC2)
- β dispatched after α signals review-readiness: `status: review-ready` in `.cdd/unreleased/343/self-coherence.md` on `origin/cycle/343`.
- β writes review passes incrementally to `.cdd/unreleased/343/beta-review.md`, committing + pushing after each pass.
- On approval: β merges `cycle/343` into main with `Closes #343` in the merge commit, then writes `.cdd/unreleased/343/beta-closeout.md`.

## What β must review

This is a docs-only cycle. The diff will contain:
- A new or updated "Git identity for role actors" prescription section in one canonical cdd file.
- Updated `cdd/review/SKILL.md` Review identity section (cycle #287 doctrine replaced or cross-referenced).
- Updated `cdd/operator/SKILL.md` if it prescribes identities.
- Any SKILL.md worked examples updated to new three-level form.
- Migration paragraph in `cdd/post-release/SKILL.md`.

**Key review targets per AC:**

AC1 oracle: `rg '@cdd\.' src/packages/cnos.cdd/skills/cdd/` returns only matches inside migration/history blocks or deprecated-example rows — no surviving prescription of the old two-level form.
AC1 positive: `rg '\.cdd\.cnos' src/packages/cnos.cdd/skills/cdd/` returns the new prescription site(s).

AC2: One form chosen for cnos self-identity. No "either is fine." A worked example shows it explicitly.

AC3: Worked example table present, 3–5 rows, deprecated row marked.

AC4: `cdd/post-release/SKILL.md` has a "Migration" or "Identity migration" subsection. ≤80 words. Names cycle #343 and cutover date.

AC5: `git log --format='%ae' cycle/343` on the branch shows only `alpha@cdd.cnos` (or whatever AC2 resolved). No old-form `alpha@cdd.cnos` vs `alpha@cdd.{project}` confusion — verify α's commits actually carry the new form.

## §2.5b disconnect note

This is a docs-only cycle (no version bump, no tag). The disconnect is the merge commit on main. β must NOT run `scripts/release.sh`, bump VERSION, or push a tag. Per `release/SKILL.md` §2.5b: the merge commit hash is the disconnect signal. The cycle directory move (`.cdd/unreleased/343/` → `.cdd/releases/docs/{ISO-date}/343/`) is γ's step post-merge.

## Dispatch note

δ: dispatch β with `claude -p` against the cnos repo. β polls `origin/cycle/343` for the review-readiness signal in `.cdd/unreleased/343/self-coherence.md` before beginning. Pass the prompt text above.
