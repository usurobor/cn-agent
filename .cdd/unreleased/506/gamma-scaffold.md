---
issue: 506
protocol: cds
main_sha: 1ce2d27dd097e00fd4b791e984547b6735aeb464
cycle_branch: cycle/506
wake_run_id: 28277178937
role: gamma
round: R0
---

# γ scaffold — issue #506

## Source of truth

| Surface | Canonical source | Status |
|---|---|---|
| Canonical essay home | `docs/papers/README.md` | Shipped (Pass 3 / `09414da`) |
| Moved papers inventory | `docs/papers/*.md` | Shipped — target paths for repoint |
| Stale `src/` citations | grep enumerated below | Current — surface to fix |
| `CDS.md` essay-home convention | `src/packages/cnos.cds/skills/cds/CDS.md` §design-notes | Stale — names deprecated `docs/gamma/essays/` |

## Scope guardrails

In scope: repoint live `src/` citations from `docs/{alpha,gamma}/essays/<F>.md` (and `docs/essays/agent-first.md`) to `docs/papers/<F>.md`; update `CDS.md` convention statements naming the essay home.

Out of scope: file moves; CI golden changes; stub removal; `.github/workflows/` edits; `.cdd/` frozen records; kata fictional scenario.

## AC oracle list

| AC | What to verify | Oracle |
|---|---|---|
| AC1 | No live `src/` citation at pre-move essay path | `grep -rn "docs/alpha/essays/\|docs/gamma/essays/\|docs/essays/" src/` returns only allowed exceptions (extraction-map.md + kata.md) |
| AC2 | Repointed targets resolve to real `docs/papers/` files | All repointed markdown links resolve to existing files under `docs/papers/`; no stub paths |
| AC3 | `CDS.md` essay-home convention names `docs/papers/` | `grep -n "docs/gamma/essays/" src/packages/cnos.cds/skills/cds/CDS.md` returns no convention/home statements |
| AC4 | No moves, no goldens, no workflows touched | `git diff --name-only origin/main..HEAD` lists only `src/**/*.md` and `.cdd/unreleased/506/**`; no `.golden.yml` or `.github/workflows/` paths |
| AC5 | Frozen records untouched | No diff under `.cdd/`, no version-stamped snapshot paths, no kata changes |

## Files to change

Derived from `grep -rn "docs/alpha/essays/\|docs/gamma/essays/\|docs/essays/" src/`:

| File | Line(s) | Change |
|---|---|---|
| `src/packages/cnos.handoff/skills/handoff/HANDOFF.md` | 49–50 | CCNF-AND-TYPED-TRUST + DECREASING-INCOHERENCE → `docs/papers/` |
| `src/packages/cnos.cdr/docs/empirical-anchor-cph.md` | 181 | CCNF-AND-TYPED-TRUST link → `docs/papers/` |
| `src/packages/cnos.cds/skills/cds/CDS.md` | 174–179, 357, 2784, 3312 | CCNF-AND-TYPED-TRUST links + 3 convention statements → `docs/papers/` |
| `src/packages/cnos.cdr/skills/cdr/CDR.md` | 130–131, 134, 501, 579 | CCNF-AND-TYPED-TRUST links → `docs/papers/` |
| `src/packages/cnos.eng/skills/eng/README.md` | 9 | ENGINEERING-LEVEL-ASSESSMENT → `docs/papers/` |
| `src/packages/cnos.cdd/commands/cdd-verify/README.md` | 100 | CCNF-AND-TYPED-TRUST → `docs/papers/` |
| `src/packages/cnos.cdd/skills/cdd/CDD.md` | 120, 146 | CCNF-AND-TYPED-TRUST → `docs/papers/` |
| `src/packages/cnos.cdd/skills/cdd/review/contract/SKILL.md` | 80 | MANIFESTO → `docs/papers/` |
| `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | 325 | BOX-AND-THE-RUNNER → `docs/papers/` |

## Files to leave unchanged (deliberate)

| File | Reason |
|---|---|
| `src/packages/cnos.handoff/docs/extraction-map.md` | Historical Q&A sweep record — do not falsify |
| `src/packages/cnos.cdd.kata/katas/M2-review/kata.md` | Fictional PR scenario — intentional training content |

## Known gaps

- Redirect stubs remain (by design — grace period not ended).
- `docs/gamma/essays/CDD-OVERVIEW.pdf` out of scope.

## α prompt

Act as α on cycle/506 for issue #506.

Mode: docs-only. Collapse allowed (Persona commitment 5).

Make the file changes listed in the "Files to change" table above. For each:
- Replace `docs/gamma/essays/<F>.md` → `docs/papers/<F>.md` in prose mentions and markdown link targets.
- Replace `docs/alpha/essays/<F>.md` → `docs/papers/<F>.md` (specifically MANIFESTO.md in review/contract/SKILL.md).
- For markdown links, update both display path mentions in code/prose AND the `(…)` href target.
- For relative-path hrefs, preserve the correct depth (number of `../` segments).
- For the CDS.md convention statements naming `docs/gamma/essays/` as the design-note home, change them to `docs/papers/`.
- Do NOT change extraction-map.md or kata.md.
- After changes, append `self-coherence.md §R0`.

## β prompt

Act as β on cycle/506 for issue #506. Walk each AC oracle independently:

AC1: run `grep -rn "docs/alpha/essays/\|docs/gamma/essays/\|docs/essays/" src/` — verify only allowed exceptions remain.
AC2: for each repointed file, verify the target `docs/papers/<F>.md` exists.
AC3: run `grep -n "docs/gamma/essays/" src/packages/cnos.cds/skills/cds/CDS.md` — verify no convention/home statements remain.
AC4: run `git diff --name-only origin/main..HEAD` — verify only `src/**/*.md` and `.cdd/unreleased/506/**` paths.
AC5: verify no `.cdd/`, version-snapshot, or kata diffs.

Write `.cdd/unreleased/506/beta-review.md §R0` with verdict and findings.
