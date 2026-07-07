# alpha close-out — cycle #610

**Issue:** [cnos#610](https://github.com/usurobor/cnos/issues/610) — `cds-install Sub 3: cn repo install --dispatch cds` (dispatch layer). Sub 3 of the #607 wave (`cn repo install`), depends on #608 (base installer, merged) and #609 (renderer generalization, merged via PR #619). **Verdict:** CONVERGE after one iterate round (R0 REQUEST CHANGES → R1 fix → R1 CONVERGE, independently re-verified by β).

Written per `alpha/SKILL.md` §2.8's re-dispatch path, in the same session as β's closeout for operational reasons only (closeouts are retrospective, not review-independent the way R0/R1 implementation-vs-review needed to be) — the content below reflects what α actually built and what R0/R1 actually found, not a review-independent re-derivation.

---

## What α built

`cn repo install --dispatch cds` now routes to the (already-merged, #609) generalized wake renderer instead of unconditionally refusing. Concretely, in `src/go/internal/repoinstall/repoinstall.go` and `src/go/internal/cli/cmd_repo_install.go`:

- `Args`/`Options` gained the identity flags the renderer already accepts (`--agent`, `--workflow-pat-secret`, `--bot-name`, `--bot-id`), named to mirror `cn-install-wake`'s own flag vocabulary 1:1 rather than inventing new names.
- `validateDispatch`/`Run` were rewired off the Mock-B4 hard-refusal stub (`"--dispatch cds requires generalized wake renderer support (#609)"`) onto the real renderer call path.
- AC3's cnos#493 label gap (canonical-label install, still open) is detected and surfaced as an actionable, nonzero-exit error rather than a silent skip — cnos#493 itself is not implemented by this cell (out of scope, per the issue's own Scope text and γ's scaffold).
- AC5's tenant-prose-clean requirement was satisfied by rendering `cds-dispatch/SKILL.md` free of `sigma`/`SIGMA_WORKFLOW_PAT`/`41898282` literals for non-default agents, while keeping `--agent sigma`'s render byte-identical to this cycle's own re-committed golden.
- A stale claim in `docs/guides/INSTALL-CDS.md` ("`--dispatch cds` ... fails explicitly ... writes nothing") was found to be a third sibling of the same now-false claim already fixed in `cmd_repo_install.go` help text and `cds-dispatch/SKILL.md` prose, and was fixed as part of §2.3 peer enumeration.

Net diff: 14 files, +2087/-81, concentrated in `internal/repoinstall/repoinstall.go` (+215 net) and its test file (+409 net) plus the renderer SKILL.md/golden pair and the doc fix above.

## R0 → R1: five findings, one fix round

β's R0 verdict was REQUEST CHANGES on 5 mechanical/procedural findings; both of α's documented interpretive judgment calls (see below) were independently re-verified as accepted/MATCH and were **not** reopened in the fix round. Summary of what was fixed (full evidence trail is in `self-coherence.md` §R1 and independently re-confirmed in `beta-review.md` §R1):

1. **(D) Non-canonical section headers → CI red.** `self-coherence.md` used `## §Gap` / `## §Skills` / `## §ACs` / `## §CDD Trace` instead of the bare `## Gap` / `## Skills` / `## ACs` / `## CDD Trace` form `cn cdd verify`'s `sectionPresent()` requires. Before `beta-review.md` existed, `classifyCycleType` classified the cycle as "small-change" and applied the non-lenient section check, producing a hard FAIL (`Build` CI red at the review SHA, `181 passed, 1 failed`). Fixed by renaming every header to the bare form; verified clean under both the lenient triadic path and the non-lenient small-change path (reproduced by temporarily removing `beta-review.md` in a scratch copy).
2. **(D) Missing `mock_parity` block.** The issue's own "Parity requirement" text mandates a `mock_parity` YAML block (C1–C6 plus an explicit AC5 row, `missed: 0`) — entirely absent from `.cdd/unreleased/610/`. Fixed by adding the block under §CDD Trace, every row's evidence field citing a test or transcript already present in §ACs (no evidence fabricated to close this finding).
3. **(C) Missing `gamma-clarification.md` for Rule 7.** `docs/guides/INSTALL-CDS.md` was touched and disclosed in §CDD Trace as a peer-enumeration addition, but no `gamma-clarification.md` existed to formally ratify extending δ's pinned "Package scoping" row to cover it. Fixed by authoring `.cdd/unreleased/610/gamma-clarification.md` (retroactive ratification, no content revert — β had already verified the doc fix itself was correct).
4. **(C) Wrong test counts.** §Review-readiness claimed "26 tests in `internal/repoinstall`, 39 in `internal/cli`... 65 total" against the actual 27 + 45 = 72. Fixed by re-deriving the counts from runner output (`grep -c '^--- PASS'`) rather than manual enumeration, matching β's independently-derived figure exactly.
5. **(B, non-blocking) Prose redundancy at line 101.** The line-101 sentence in `cds-dispatch/SKILL.md` read as a grammatically redundant tautology for non-sigma agents (`"...as acme (today: acme; ...)"`). Fixed by rephrasing to describe the binding mechanism (`--agent`/`--workflow-pat-secret`/`--bot-name`/`--bot-id`) without repeating the agent name — deliberately not using β's own illustrative wording, which named `sigma` literally and would have regressed AC4's blanket case-insensitive `sigma` grep for non-sigma renders.

β's R1 pass independently re-derived every one of these five findings from scratch (fresh binary, fresh index, fresh scratch repos) rather than trusting the R1 self-coherence appendix, and confirmed CONVERGE with no regressions and no scope creep (`git diff --stat` showing exactly the 5 touched files).

## Two interpretive judgment calls, and how they were handled

Two genuine ambiguities surfaced during implementation that neither the issue body nor `gamma-scaffold.md` fully resolved:

1. **AC1/AC3 relationship** — whether `Run()` must return non-nil (nonzero exit) whenever the cnos#493 label gap is present, even after an otherwise-successful render, or whether a printed warning with exit 0 satisfies AC3's "actionable error... not a silent skip" language. Mock C's console transcript (which predates the AC3 requirement) offered no exit-code signal either way.
2. **AC5 prose-clean vs. the pinned byte-identical-to-golden floor** — for the one line-296 leak (a dated, un-templatable historical incident citation) that γ's own scaffold flagged as impossible to render honestly for any agent but sigma.

Both were resolved by authoring `.cdd/unreleased/610/alpha-clarification-needed.md` at implementation time: stating the ambiguity, the assumption made, the reasoning, and — critically — the exact, narrow fix if the assumption is wrong, rather than stalling or silently picking a reading. β reviewed both at R0, accepted both as-implemented with no override, and did not reopen them at R1. This is the same shape as `gamma-clarification.md`-style handling (name the gap, make a reasoned call, leave an explicit reversal path) but surfaced from α's side before review rather than after a finding forced it — worth naming as a pattern that worked, since it let β's R0 review spend its judgment budget confirming or overriding two pre-named calls instead of discovering them cold.

## Process friction worth naming

**The section-header convention gap was not a novel mistake — it was a repeat of a precedent #608 had already hit and fixed, three cycles earlier, in the same design-doc family.** `.cdd/unreleased/608/self-coherence.md` explicitly documents hitting this identical failure mode (bare `## Gap` vs. `## §Gap`, the `checkSmallChangeArtifacts`/`classifyCycleType` non-lenient-path interaction) and fixing it, in prose, in its own §Self-check section. Cycle #608's `gamma-closeout.md` also logged the underlying `classifyCycleType` misclassification as "already tracked" against cnos#577. Neither artifact is discoverable to a fresh α/γ working #610 except by manually reading a prior cycle's own `self-coherence.md` — the lesson never made it into `alpha/SKILL.md`, `gamma/SKILL.md`, or any Tier 2/3 skill that #610's dispatch actually loaded. #610 re-hit the exact same failure (real CI red at the review SHA, `181 passed, 1 failed`) that #608 had already paid the cost to discover. This is real, avoidable-in-hindsight process debt: a lesson that lives only in one cycle's own retrospective artifact does not propagate to sibling cycles working the same design-doc family, and should not be expected to via manual cross-cycle reading. β's closeout names the same gap from the review side; this entry is the α-side confirmation that the gap was real and cost a full R0→R1 round on this cycle specifically.

No other significant friction: the peer-enumeration and harness-audit disciplines in `alpha/SKILL.md` §2.3/§2.4 worked as intended (the `docs/guides/INSTALL-CDS.md` peer was caught by α's own enumeration, not by β — the R0 finding on it was a procedural ratification gap, not a missed peer).

## Known debt and follow-on work for γ's triage

1. **cnos#493 (canonical label install) is still OPEN and is real follow-on debt, not resolved by this cell.** `cn repo install --dispatch cds` now detects the label gap and fails with an actionable, nonzero-exit error naming cnos#493 — but every current invocation of `--dispatch cds` ends in that failure until cnos#493 ships the actual label-install mechanism. There is currently no way to get a fully successful (exit 0) `--dispatch cds` run against a repo lacking the canonical labels; this cell's AC3 was explicitly scoped to "actionable error," not "install the labels." This is worth naming explicitly for γ's triage: #607's wave is not fully closed end-to-end until #493 lands, even though #608/#609/#610 (Subs 1–3) are each individually complete and merged/converged.
2. **The header-convention lesson from #608 is a skill/process gap, not a code gap** — see β's closeout for the full framing of what γ should triage (whether the lesson needs to move into `alpha/SKILL.md`/`gamma/SKILL.md` directly, or whether some other propagation mechanism is needed so a prior cycle's self-coherence.md observation doesn't silently expire at that cycle's own directory boundary).
3. **`classifyCycleType`'s hard-fail-vs-lenient-warn split** (already tracked at cnos#577, per #608's gamma-closeout) is the same underlying mechanical cause of both #608's and #610's header-format CI failures. Two independent occurrences across two cycles is stronger evidence for prioritizing #577 than either occurrence alone.

Per `alpha/SKILL.md` §2.8's voice constraint, the above is reported as factual pattern/observation, not a recommended disposition — triage is γ's job.

## Final-state evidence

Re-verified at the final cycle HEAD (`012bc2ae1c26e72de27d1fbe49a823b5816b3498`), from `src/go/`:

```
$ go build ./...              # exit 0
$ go vet ./...                # exit 0
$ go test ./internal/repoinstall/... ./internal/cli/... -v 2>&1 | grep -c '^--- PASS'
72
$ go test ./internal/repoinstall/... ./internal/cli/... -v 2>&1 | grep -c '^--- FAIL'
0
```

From repo root: `cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` → `## Summary: 184 passed, 0 failed, 121 warnings (305 total)`, verdict "PASSED with warnings" (0 failed; the 121 warnings are the standard unreleased-cycle in-progress/stale-cycle advisories across the other 36 cycle directories, none new to #610).

## Artifact list

`.cdd/unreleased/610/`: `gamma-scaffold.md`, `CLAIM-REQUEST.yml`, `self-coherence.md` (with R1 fix-round appendix), `alpha-clarification-needed.md`, `gamma-clarification.md`, `beta-review.md` (R0 + R1), this file.
