# γ scaffold — Cycle 423

**Cycle:** [cnos#423](https://github.com/usurobor/cnos/issues/423) — Build-fix for cnos#416 cross-repo stub missing triggers
**Branch:** `cycle/423` (from `origin/main` @ `7a1f7024` — the merge that closed cycle/422 / v3.82.0 release-hygiene)
**Mode:** docs-only / build-fix (P0)
**Collapse pattern:** β-α-collapse-on-δ. Commits: `α-423:`, `β-423:`, `γ-423:`.
**Date:** 2026-05-23

## Intent

cnos#416 (Sub 2 of #404 / handoff extraction) replaced the canonical cross-repo SKILL.md inside `cnos.cdd` with a 28-line compatibility pointer to the new canonical home at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. The stub's frontmatter declared `artifact_class: pointer` and `status: moved` but omitted the `triggers:` field. The activation validator at `src/go/internal/activation/Validate` treats any SKILL.md without non-empty `triggers:` as an `IssueEmptyTriggers` defect, and `TestValidate_RealCorpus_NoEmptyTriggers` (`src/go/internal/activation/index_test.go:312`) treats any `IssueEmptyTriggers` against the real corpus as a hard test failure. CI is red on `main` post-v3.82.0 merge; the tag has not yet been pushed.

This cycle is **one single-file frontmatter edit** to restore green CI before the v3.82.0 tag goes out. Nothing else.

## Surface plan (D1)

### D1 — Add `triggers:` field to `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` frontmatter

Insert two triggers, neither colliding with the canonical's five (`cross-repo`, `proposal intake`, `outbound iteration`, `bilateral iteration`, `feedback patch`):

```yaml
triggers:
  - cross-repo-moved
  - handoff-extracted
```

Both names are descriptively accurate for the pointer's content (the file announces that cross-repo doctrine has moved + that handoff was extracted to its own package) and are distinct from the canonical's keyword set, so the activation conflict-detector will not produce duplicate-trigger warnings.

The body of the SKILL.md is **unchanged byte-for-byte**. Only the frontmatter receives 3 new lines.

## Implementation contract

| Axis | Pinned value |
|---|---|
| Language | Markdown frontmatter (YAML) |
| CLI integration target | None |
| Package scoping | **One file edit:** `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md`; plus cycle-close artifacts in `.cdd/unreleased/423/` and `.cdd/iterations/INDEX.md` row |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | YAML triggers list (string array of length 2) |
| Backward compat | Stub remains a pointer; `artifact_class: pointer` / `status: moved` / `canonical:` unchanged; body verbatim |

## Acceptance criteria

AC1–AC7 per [cnos#423](https://github.com/usurobor/cnos/issues/423). All mechanical. Verified in `self-coherence.md` post-α.

- AC1: `cd src/go && go test ./internal/activation/... -run TestValidate_RealCorpus_NoEmptyTriggers -count=1` exits 0.
- AC2: `cd src/go && go test ./... -count=1` exits 0.
- AC3: `git diff origin/main..HEAD --name-only -- src/packages/` lists exactly one path: `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md`.
- AC4: The diff hunk against the stub adds exactly 3 frontmatter lines (`triggers:` + two list items) and changes no body line; `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` shows `+3 -0` (or equivalent `wc -l` on the unified diff hunk).
- AC5: Neither new trigger collides with the canonical's keyword set (`grep -E "^  - (cross-repo|proposal intake|outbound iteration|bilateral iteration|feedback patch)$" src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` returns 0 matches).
- AC6: `git diff origin/main..HEAD -- src/go/` returns 0 lines (no test modifications; no validator modifications).
- AC7: `git diff origin/main..HEAD --name-only -- src/packages/` lists no other skill edits (no protocol skill content changes beyond the stub frontmatter addition).

## Hard rules

1. **Single file edit.** Only `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` (plus standard close-out artifacts in `.cdd/unreleased/423/` and the iterations INDEX.md row).
2. **Frontmatter-only addition.** Body verbatim.
3. **No test modifications.** Fix the data, not the contract. `internal/activation/index_test.go` is unchanged.
4. **No other skill edits.** If other trigger-missing skills exist (the test reports 26 corpus issues; an unknown number may be additional empty-triggers), leave them for a follow-on cycle. The cdd-iteration.md may name the broader question as a next-MCA candidate.
5. **No protocol evolution.** No CCNF kernel changes; no schemas; no #405 work; no new packages.

## Non-goals

- Do NOT modify `internal/activation/Validate` to exempt `artifact_class: pointer` skills from the triggers check (architecturally cleaner but invasive; recorded as next-MCA disposition in `cdd-iteration.md`).
- Do NOT modify `internal/activation/index_test.go`.
- Do NOT patch other trigger-missing skills surfaced by the corpus test (the 26 surfaced issues include conflict / missing rows that are corpus-level questions; this cycle is bounded to the cnos#416 stub).
- Do NOT push the v3.82.0 tag. Operator handles that post-merge per `cnos.cdd/skills/cdd/release/SKILL.md §2.6`.
- Do NOT amend cycle/416's commit history; this cycle's α commit explicitly cites cnos#416 as the upstream stub-authoring cycle whose dispatch brief did not enumerate the triggers requirement.

## Operator action on close

After all 7 ACs PASS, push: `git push -u origin cycle/423`. Then merge to main with `--no-ff` and message `Merge cycle/423: Build-fix for cnos#416 cross-repo stub missing triggers. Closes #423.`. Post-merge: CI returns green; the v3.82.0 tag may then be pushed via `scripts/release.sh`.
