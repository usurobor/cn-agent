# α closeout — Cycle 423

**Cycle:** [cnos#423](https://github.com/usurobor/cnos/issues/423) — Build-fix for cnos#416 cross-repo stub missing triggers
**Branch:** `cycle/423`
**Author:** α (collapsed γ+α+β on δ)
**Date:** 2026-05-23

## What α shipped

One frontmatter-only edit:

| Path | Disposition | Purpose |
|---|---|---|
| `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` | Modified (+3/-0) | Add `triggers:` field with two list items (`cross-repo-moved`, `handoff-extracted`) so the activation validator stops reporting `IssueEmptyTriggers` against the cnos#416 stub |

## What α did NOT ship

Per the hard rules in `gamma-scaffold.md` and the AC pass-set in `self-coherence.md`:

- **No body changes.** The stub's 18-line body (lines 11–29) is byte-identical to the cnos#416 authored version. The pointer narrative (canonical home, contents, deprecation note) is preserved verbatim.
- **No test modifications.** `src/go/internal/activation/index_test.go` is untouched. The fix patches the data the test runs over, not the test's contract.
- **No validator modifications.** `src/go/internal/activation/validate.go` is untouched. The architecturally-cleaner alternative — exempting `artifact_class: pointer` skills from the empty-triggers check — is recorded as a `next-MCA` disposition in `cdd-iteration.md`, not executed here.
- **No other skill edits.** No edits to other SKILL.md files in any of the four protocol packages. The corpus test surfaces 26 issues total post-fix (conflict / missing rows; all non-IssueEmptyTriggers and therefore non-blocking for this test); any of those that warrant patching is a separate cycle.
- **No CCNF kernel changes.** `CDD.md` / `COHERENCE-CELL.md` / `COHERENCE-CELL-NORMAL-FORM.md` byte-identical to `origin/main`.
- **No #405 / Track A / Track B work.**
- **No tag push.** The v3.82.0 tag has not been pushed; operator handles that post-merge of this cycle's PR.
- **No changes to cnos.core / cnos.eng / cnos.kata / cnos.cdd.kata.**

## Why the chosen trigger values

The two new triggers — `cross-repo-moved` and `handoff-extracted` — were chosen so that:

1. **Neither collides with the canonical's five triggers** at `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (`cross-repo`, `proposal intake`, `outbound iteration`, `bilateral iteration`, `feedback patch`). The activation index's conflict detector treats duplicate trigger keywords across skills as a `IssueConflictingTriggers` defect; the dispatch brief explicitly enumerated the canonical's five so the stub could avoid them.
2. **Both describe what the stub actually announces** — that the cross-repo doctrine has moved, and that handoff was extracted to its own package per cnos#404 / cnos#416. A loader searching the trigger keyword space for `cross-repo-moved` or `handoff-extracted` is plausibly looking for the deprecation notice; the canonical's five describe the live operational surface.
3. **They are unique to this stub.** A second grep across `src/packages/**/SKILL.md` confirms no other skill claims either keyword.

## Verification handoff to β

- `go test ./internal/activation/... -run TestValidate_RealCorpus_NoEmptyTriggers -count=1` exits 0 (AC1).
- `go test ./... -count=1` exits 0 (AC2).
- `git diff origin/main..HEAD --name-only -- src/packages/` lists exactly one path (AC3).
- The unified diff hunk is +3/-0 inside the frontmatter (AC4).
- No collision with the canonical's triggers (AC5).
- `git diff origin/main..HEAD -- src/go/` returns 0 lines (AC6).
- No other skill edits (AC7).

All 7 ACs PASS mechanically per `self-coherence.md`. Handoff to β for adversary review.

## Cdd-iteration finding to be filed

This cycle surfaces one binding `cdd-skill-gap` finding (loaded-skill miss): the dispatch brief for cnos#416 did not load the activation/Validate requirements as a Tier 3 reference, so the stub-authoring step proceeded without the "every SKILL.md needs non-empty `triggers:`" rule in scope. The full per-finding record (root cause, disposition, first AC for the eventual MCA) is in `cdd-iteration.md`. The disposition is `next-MCA` — operator picks the remediation path at next-MCA dispatch.
