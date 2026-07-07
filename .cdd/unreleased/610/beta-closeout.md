# beta close-out — cycle #610

**Issue:** [cnos#610](https://github.com/usurobor/cnos/issues/610) — `cn repo install --dispatch cds` (dispatch layer), Sub 3 of the #607 wave. **Verdict:** CONVERGE after one iterate round (R0 REQUEST CHANGES → R1 fix → R1 CONVERGE).

Written per `beta/SKILL.md`'s Phase map ("β close-out → written to `.cdd/unreleased/{N}/beta-closeout.md`"), in the same session as α's closeout for operational reasons only — closeouts are retrospective narrative, not the review-independence R0/R1 required, so authoring both in sequence here does not compromise either one. Both round verdicts themselves (`beta-review.md` §R0, §R1) were written and independently re-derived before this closeout, and are not revised here.

---

## Review Summary

Two review passes. **R0** (head `b140664`, base `f7e9aaad34793dbea80c603e315e0ecc0760fdfa`) returned REQUEST CHANGES on 5 findings (2×D, 2×C, 1×B) plus two accepted judgment calls. **R1** (head `387b01fbde25fcc100e862ccafed29332ab00e39`, same base, re-confirmed unchanged) independently re-verified all 5 fixes from scratch — fresh `cn` binary, fresh `dist/packages/index.json`, fresh scratch git repos for every e2e run, nothing copied from α's transcript — and returned CONVERGE. The final review-verdict commit is `012bc2ae1c26e72de27d1fbe49a823b5816b3498`.

Review quality: one iterate round for 5 findings, all mechanically small (header rename, one YAML block, one paperwork artifact, two corrected numbers, one prose tightening), no re-derivation of already-collected AC evidence required. Fix-round diff was scoped to exactly 5 files with no regressions and no scope creep (`git diff --stat 3fb8e47..387b01f`). The R1 re-verification held itself to the same standard as R0 — real CI was polled to completion (not assumed) at the R1 head, both `Build` and `install-wake golden` reporting `conclusion=success`; mock_parity rows were spot-checked against real test names, not the citing text; the line-101 prose fix was confirmed sigma-literal-free via direct grep on a freshly rendered file, not by reading source.

## Was the RC round avoidable?

**Yes — for the D-severity finding that actually caused CI red (Finding 1), this was an avoidable-in-hindsight class of finding, not novel risk.** The root cause is not that section-header formatting is inherently hard to get right; it is that this exact failure mode — `## §Gap` vs. the bare `## Gap` form `cn cdd verify`'s `sectionPresent()` requires, and its interaction with `classifyCycleType`'s non-lenient "small-change" path before `beta-review.md` exists — was already discovered, diagnosed, and fixed three cycles earlier in `.cdd/unreleased/608/self-coherence.md`, working the same `cn repo install` design-doc family (#608 is Sub 1 of the same #607 wave #610 is Sub 3 of). #608's own self-coherence.md documents the fix explicitly enough to be unambiguous: "Section headers use the bare form (`## Gap` not `## §Gap`) ... required for `cn cdd verify`'s I6 gate to actually pass." #608's `gamma-closeout.md` additionally logged the underlying `classifyCycleType` misclassification as tracked against cnos#577.

Neither #610's γ-authored scaffold nor α's implementation consulted or propagated that precedent. The result: #610 hit the identical failure (real `Build` CI red at the R0 review SHA, `181 passed, 1 failed`), costing a full R0→R1 round for what should have been a zero-round non-issue if the precedent had been consulted. This is the single real, avoidable-in-hindsight finding class in this cycle — everything else (the `mock_parity` block, the missing `gamma-clarification.md`, the test-count arithmetic, the prose redundancy) is ordinary review-catches-a-gap-α-should-have-caught territory, not a repeat of a documented prior-cycle lesson.

## Release Evidence

This cell runs in **wake-invoked mode**: per the R1 conclusion in `beta-review.md`, "δ next dispatches γ+α+β for closeouts (β is not responsible for closeouts in this pass, and does not merge to main — merge happens later via draft PR #620 after external review)." Accordingly:

- **No `git merge` was executed by β for this cycle.** The standard `release/SKILL.md` §2.6 flow ("β signals 'release ready for δ tag' in `beta-closeout.md`... δ tags") does not apply as written here, because the merge itself has not happened yet — this closeout is written at the review-CONVERGE boundary, not the post-merge boundary.
- **Draft PR [#620](https://github.com/usurobor/cnos/pull/620) exists and carries the work.** It is open and marked draft at the time of this writing. The actual merge to `origin/main`, the close-keyword pre-merge gate row (`beta/SKILL.md` pre-merge gate row 5), tag/deploy, and the cycle-directory move to `.cdd/releases/{X.Y.Z}/610/` are all external to this session and later — not claimed as done here.
- **Final review-verified state at HEAD `012bc2ae1c26e72de27d1fbe49a823b5816b3498`:** from `src/go/`, `go build ./...` and `go vet ./...` exit 0; `go test ./internal/repoinstall/... ./internal/cli/...` → 72 tests, 72 PASS, 0 FAIL. From repo root, `cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` → `## Summary: 184 passed, 0 failed, 121 warnings (305 total)`, "PASSED with warnings" (0 failed). This is the release-readiness evidence for whoever executes the actual merge/tag externally; it does not itself constitute the merge.
- Issue #610 remains **OPEN** at the time of this writing (the close-keyword merge that will close it has not happened).

## Process-gap observation for γ's triage

γ owns the PRA; this is not that document. But one process-gap observation is worth naming explicitly so γ can decide where it belongs, since it is the direct cause of the one avoidable round in this cycle:

**Cycle #608's `self-coherence.md` documented a real, specific header-convention lesson (bare `## Gap` vs. `## §Gap`, and its interaction with `classifyCycleType`'s non-lenient small-change path) that was never consulted by #610's γ (at scaffold time) or α (at implementation time) — three cycles later, in the same design-doc family, working the same wave.** The lesson lived only in #608's own cycle-scoped artifact (`self-coherence.md` §Self-check, echoed in `gamma-closeout.md`'s PRA table against cnos#577); it was never promoted into a durable skill surface (`alpha/SKILL.md`, `gamma/SKILL.md`) that a later cycle's dispatch would actually load. `alpha/SKILL.md` §2.3 already has a mechanism for exactly this shape of problem — the "*Derives from: #NNN F#*" empirical-anchor pattern used throughout that file to convert a prior cycle's specific finding into a durable, loaded-at-dispatch-time rule (e.g. the #283, #600, #266 anchors already in §2.3). This lesson never went through that promotion step.

This is a real gap in the propagation mechanism, not a one-off oversight: a `self-coherence.md` or `gamma-closeout.md` observation is scoped to its own cycle directory and is not part of any role's Tier 1 load order (`alpha/SKILL.md` §Load Order loads `CDD.md` + `alpha/SKILL.md` + lifecycle sub-skills — not prior cycles' artifacts). Without an explicit promotion step, a documented lesson is discoverable only by an agent that happens to read the specific prior cycle's directory, which is not guaranteed and did not happen here. γ should triage whether the fix is:

- promoting this specific lesson into `alpha/SKILL.md` §2.2 (produce-in-artifact-order) or §2.5 (self-coherence authoring) as a named empirical anchor, the same way #283/#600/#266 already are, and/or
- promoting it into `gamma/SKILL.md` (scaffold-authoring time, since γ authors the cycle directory before α ever touches it, and a scaffold template could name the canonical header convention directly), and/or
- treating the underlying `classifyCycleType` non-lenient-path mechanical cause (cnos#577, already tracked per #608's own gamma-closeout) as the actual fix — if `cn cdd verify` simply accepted the `§`-prefixed form or applied the lenient path pre-`beta-review.md`, neither #608 nor #610 would have had a header-format finding to make in the first place, which would close the class at the tool level instead of at the skill-propagation level.

Naming the tool-level option alongside the skill-level ones because two independent occurrences (#608, #610) of the same finding, three cycles apart, in the same wave, is stronger evidence than either occurrence alone that this is worth γ prioritizing rather than deferring again.

## For a future reviewer of a similar cycle

Before reviewing a cycle that is a later Sub of an already-in-progress wave (as #610 is Sub 3 of #607, following merged Subs #608 and #609), grep the prior Subs' own `self-coherence.md`/`gamma-closeout.md`/`beta-closeout.md` for any finding class that could recur mechanically (header conventions, test-count derivation discipline, mock_parity block presence) before starting R0 — not as a substitute for independent review, but because a wave's later cells inherit the same CDD-tooling surface as its earlier cells, and a prior cell's already-paid discovery cost should not be paid twice.

## Summary for γ

- Review quality: 1 iterate round, 5 findings (2D/2C/1B), clean CONVERGE at R1, no scope creep in the fix diff.
- Root cause of the RC round: an avoidable repeat of a documented #608 precedent (header-convention + `classifyCycleType` interaction) that did not propagate to #610's γ scaffold or α implementation.
- Release evidence: wake-invoked mode, no merge executed by β; draft PR #620 carries the work; final state at HEAD `012bc2ae1c26e72de27d1fbe49a823b5816b3498` is build/vet/test-clean (72/72 passing) and `cn cdd verify` clean (0 failed, 121 pre-existing warnings unrelated to #610).
- Triage item for γ: promote the #608 header-convention lesson into a durable skill surface (`alpha/SKILL.md` and/or `gamma/SKILL.md`) or fix the underlying `classifyCycleType` mechanism (cnos#577) — either closes the recurrence class; leaving the lesson only in #608's own artifact does not.
