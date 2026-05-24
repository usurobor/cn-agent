# ε cdd-iteration — cycle/427

**Issue:** [cnos#427](https://github.com/usurobor/cnos/issues/427) — Rewrite v3.82.0 release notes per write/SKILL.md (third use of remote-runner-delegation doctrine).
**Mode:** docs + workflow + receipt; γ+α+β-collapsed-on-δ. **Substantive** entry: records one `cdd-skill-gap` finding with disposition `next-MCA`, surfaced by the existence of this cycle (the cycle had to be authored because cnos#422's release-notes authoring did not load `cnos.core/skills/write/SKILL.md`; the gap is in `release/SKILL.md`'s discipline around RELEASE.md / CHANGELOG authoring).

## Context

This cycle is the third use of the `remote-runner delegation` primitive landed in cnos#425. The first use (cnos#425, `.github/workflows/repoint-3.82.0.yml`) moved the `3.82.0` tag to `fd1d654e`. The second use (cnos#426, `.github/workflows/publish-3.82.0-release.yml`) published the GH release with the v3.82.0 baseline body. This third use (cnos#427, `.github/workflows/republish-3.82.0-notes.yml`) updates the GH release body to a rewritten tight version that follows `cnos.core/skills/write/SKILL.md` discipline.

The structural fact making this cycle necessary: the cnos#422 release-hygiene cycle authored a 109-line v3.82.0 `RELEASE.md` that deviated from the 3.80.0/3.81.0 template (which run 30–60 lines following an "Outcome / Why it matters / Added / Changed / Unchanged / γ note" shape). The cnos#422 dispatch brief specified verbose section content (Stop condition, Does NOT include, Validation, multi-paragraph Why it matters) without referencing `write/SKILL.md`. The result was a release-notes body that violates `write/SKILL.md` §3.2 (one file = one governing question; the body mixed "what changed" with "why we paused next" and "what's deliberately deferred"), §3.3 (say a fact once; stop-condition language repeated 3+ times), §3.4 (front-load the point; opening sentence was an abstraction), §3.5 (no throat-clearing; "This is a release-hygiene cycle, not a feature cycle" is decorative contrast), §3.8 (no abstractions without evidence; "the structural pause point" was abstract framing), §3.14 (brevity earned; 109 lines vs ~50 for the template), and §3.15 (end when point delivered; Stop condition restarted the argument at file end).

The patch (this cycle) was authored under `write/SKILL.md` discipline. The gap (one cycle later) is at the *upstream* skill that should have made `write/SKILL.md` a required load for the original cycle/422 dispatch.

## Findings

### F1 — cdd-skill-gap: release/SKILL.md must require write/SKILL.md load for release-notes authoring

**Class:** `cdd-skill-gap`
**Severity:** binding (every future release-notes authoring cycle would re-discover the gap if release/SKILL.md does not require write/SKILL.md load)
**Disposition:** `next-MCA` — file a follow-on issue with the first AC below; sequence after cnos#427's merge so the rewritten v3.82.0 notes ship first, then close the upstream skill gap.

**The gap.** `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` §2.4 line 107 says:

> "If release notes or CHANGELOG wording are being authored, load the write skill and record it in the CDD Trace."

This is a soft directive: (a) phrased as a conditional rather than a hard requirement; (b) buried mid-section under the CHANGELOG sub-section (§2.4), where a γ dispatching a release-hygiene cycle is unlikely to encounter it without already-loading release/SKILL.md AND already-reading deep into §2.4; (c) NOT surfaced in the skill's loader frontmatter (`description`, `triggers`, `calls`); (d) NOT surfaced in the §"Authoring discipline" or §"What to load before writing release notes" sections (which do not exist as discrete sections — the directive appears mid-flow). The result: a cycle/422-style dispatch brief that does not enumerate `write/SKILL.md` as required can ship a release-notes body that violates write/SKILL.md without the skill graph alarming on the gap.

**Empirical anchor.** cycle/422 (v3.82.0 release-hygiene cycle, commit shipped 2026-05-23) authored a 109-line `RELEASE.md` and a 11-line `CHANGELOG.md ## 3.82.0` row, both of which violate `write/SKILL.md` §§3.2 / 3.3 / 3.4 / 3.5 / 3.8 / 3.14 / 3.15 per the discipline analysis in §AC8 of this cycle's `self-coherence.md`. The cycle/422 dispatch brief (recoverable from the cycle/422 issue body) did not enumerate `write/SKILL.md` as a required load; the release/SKILL.md soft directive at §2.4 line 107 did not trigger an alarm because the γ-as-author for cycle/422 was operating on the dispatch brief's specification rather than on a release/SKILL.md required-load list.

**Why this is a `cdd-skill-gap` rather than other classes.**
- `cdd-protocol-gap` (the cnos cell algorithm itself has an undefined edge case) is wrong — the protocol shape (γ authors release notes; β reviews; δ ships) is correctly specified. The gap is in WHAT γ is required to load before authoring, which is skill-content territory.
- `cdd-tooling-gap` (a substrate / pipeline mechanism is broken) is wrong — there is no broken substrate. The release-effector mechanics, the release CI, the tag-publishing flow all work as designed (modulo cnos#426's F1/F2 which are tooling gaps). This is a discipline-content gap in the release skill itself.
- `cdd-skill-gap` is correct — `release/SKILL.md` exists, is loaded by release-hygiene cycle dispatches, and is mostly correct. But one specific load discipline (require `write/SKILL.md` for any RELEASE.md or CHANGELOG.md authoring) is soft-cited mid-section rather than declared as a required Tier 3 load. The patch is content-level edit to release/SKILL.md.

**Why disposition `next-MCA` is correct.**
- `patch-landed-in-cycle` is wrong because this cycle's matter does not modify `release/SKILL.md` (hard rule 4 explicitly forbids modifying cnos.cdd skills; the cycle's scope is the release-notes rewrite + publication, not the upstream skill patch).
- `no-patch` is wrong because the gap is real and recurring (every future release-hygiene cycle would re-hit the same gap if release/SKILL.md does not require write/SKILL.md load).
- `next-MCA` is correct: a follow-on cycle with the scope "elevate write/SKILL.md citation in release/SKILL.md to a required Tier 3 load with named triggers" is required. The cycle's first AC is given below; the issue scope is contained (single skill file; mechanical verification via grep for the elevated citation in release/SKILL.md's loader frontmatter or §"Loads" section).

**First AC for the MCA follow-on (verbatim from cnos#427 issue body AC9):** "`release/SKILL.md` cites `write/SKILL.md` as required Tier 3 for any RELEASE.md or CHANGELOG.md authoring."

**Recommended additional scope hints for the MCA follow-on** (not binding; just signal for the next γ):

- Add a §"Required loads for release-notes authoring" section near the top of `release/SKILL.md` (or expand the existing §2.4) that explicitly lists `write/SKILL.md` as a required Tier 3 load whenever a release-hygiene cycle's dispatch involves RELEASE.md or CHANGELOG.md authoring.
- Update `release/SKILL.md`'s loader frontmatter `calls:` (if present) or equivalent to enumerate `write` as a dependency for the release-notes authoring path.
- Add a "release-notes authoring discipline" sub-section that summarizes the write/SKILL.md rules most relevant to release notes: §3.3 (say a fact once — no repeated stop-condition / non-goal language), §3.4 (front-load the outcome — opening sentence is concrete WHAT-changed, not abstract WHY), §3.5 (no throat-clearing — no "This is a release-hygiene cycle, not a feature cycle" framing), §3.14 (brevity earned — target 30–60 lines mirroring 3.80.0/3.81.0).
- Add a γ-side gate to the release-hygiene dispatch template that requires the dispatch brief to enumerate `write/SKILL.md` as a required load if the cycle authors RELEASE.md or CHANGELOG.md content.
- Add a regression check (small kata or doc-quality script) that flags release-notes bodies exceeding 90 lines or missing the write-skill γ-note acknowledgment.
- Consider cross-citing this finding in `write/SKILL.md`'s `triggers:` field by adding "writing release notes" or "writing CHANGELOG entries" as triggers so the skill graph cross-references the release path.

## Protocol-gap signals (across receipt-stream)

This cycle files 1 `cdd-skill-gap` finding; not a `cdd-protocol-gap`, not a `cdd-tooling-gap`. Looking across the receipt stream (per `cnos.handoff/skills/handoff/receipt-stream/SKILL.md`):

- cnos#425 filed 1 `cdd-protocol-gap` (boundary-model needed remote-runner delegation naming) with disposition `patch-landed-in-cycle` (essay + skill §8 + first artifact + first receipt atomically).
- cnos#426 filed 2 `cdd-tooling-gap` (F1 release.yml build-job git-auth failure; F2 GH Actions on.push.tags non-trigger on tag force-update), both `next-MCA`.
- cnos#427 (this cycle) files 1 `cdd-skill-gap` (release/SKILL.md write-skill-load discipline), `next-MCA`.

The pattern across the three release-boundary cycles: cnos#425's first-use surfaced a protocol gap (the doctrine itself was missing); cnos#426's second-use surfaced downstream substrate gaps (now that the doctrine is real, its dependencies expose substrate cracks); cnos#427's third-use surfaces an upstream skill gap (the release-notes hygiene gap that motivated this cycle's existence). This is healthy ε behavior: each cycle's iteration produces a different class of finding because each cycle is at a different position in the dependency graph — the protocol layer (cnos#425), the substrate layer (cnos#426), the skill layer (cnos#427). ε notes the pattern: when a wave of cycles works the same release boundary on adjacent surfaces, the cdd-iteration findings naturally walk through the dependency layers in surface-order.

## Non-findings (worth recording)

- **Verbatim-content discipline is a positive structural property.** This cycle's docs surfaces are operator-pre-authored verbatim (D1 + D2 from the issue body); α applied them without redraft per hard rule 1. The verbatim discipline removes a class of α-side authorial drift (each γ would otherwise draft slightly different release notes for the same factual content). The cost: AC1's "60–90 lines" lower bound was an authoring expectation that the verbatim content's natural shape (45 lines, long-bulleted condensed prose) does not satisfy; β recorded the precision-mismatch as a non-binding observation. ε observation: dispatch briefs that pre-author verbatim content should state mechanical bounds as one-sided constraints (`≤ N`) rather than ranges, because verbatim content cannot stretch to meet a lower bound.

- **Third-use of doctrine on adjacent surfaces is doctrine-class durable.** cnos#425/#426/#427 exercise the same primitive (remote-runner delegation) on three adjacent effect surfaces (tag SHA / release existence / release body) of the same release boundary (v3.82.0). The workflow shape (one-shot + push-triggered on own path + softprops + self-delete) and receipt shape (6 fields + post-run-fillable evidence + expected-effect + failure-modes + acceptance-criteria + relationship-to-prior-cycles) are stable across all three uses with only the target-effect-surface choice varying. Per `BOX-AND-THE-RUNNER.md §"What this enables"`, three consecutive uses on the same primitive class is evidence-class durable rather than one-off escape. ε observation, not a finding; could be added to a future doctrine update that names the "three-uses-makes-it-a-pattern" rule.

- **Tag SHA invariant across the three cycles.** `git ls-remote origin refs/tags/3.82.0` returns `fd1d654e` after cnos#425's retarget; cnos#426 does not touch the tag; cnos#427 does not touch the tag. The tag SHA is the structural invariant binding the three release-boundary cycles. Worth naming in a future doctrine update as the "tag-as-release-boundary-anchor" pattern: when multiple cycles work the same release boundary on different effect surfaces, the tag SHA is the shared identity that holds the boundary together.

- **Checkout ref walk reflects "publish current source of truth" rather than "pin historical snapshot."** cnos#425 implicitly-main (its effect IS the main snapshot of the tag move); cnos#426 at-tag (publishing the *historical* baseline body that the tag carries); cnos#427 at-main (publishing the *current* main body, which is the new rewritten version landed by this same cycle). The choice between at-tag and at-main is "which ref carries the load-bearing content for this cycle's effect." Pattern: at-tag for historical body / at-main for current-source-of-truth body. ε observation; could inform future remote-runner cycle dispatch briefs.

- **softprops/action-gh-release@v1 idempotency is a positive invariant for all three uses.** All three remote-runner cycles to date have idempotent target actions (`git tag -f`, softprops in create-or-update mode, softprops in update mode). Idempotent target actions are easier to reason about than non-idempotent ones (a failure-and-retry produces the same correct end state). The pattern "remote-runner cycles should prefer idempotent target actions" is now supported by three uses' worth of evidence; worth naming explicitly in a future doctrine update.

- **β-independence collapse named, not papered over.** This cycle ran as `γ+α+β-collapsed-on-δ`, inheriting cnos#425/#426/cycle-414/cycle-424 collapse precedents. The receipt is closed-as-degraded at the structural-independence axis. `beta-closeout.md` names the collapse explicitly with the inheritance and the dispatch authorization. The discipline of naming the collapse rather than hiding it is good ε hygiene and inherited from prior cycles.

## Verdict

`tooling_gap_count: 0`. `protocol_gap_count: 0`. `skill_gap_count: 1`.

One `cdd-skill-gap` finding (F1: `release/SKILL.md` must require `write/SKILL.md` load for release-notes authoring), with disposition `next-MCA`. First AC verbatim per the issue body. No `patch-landed-in-cycle` or `no-patch` dispositions.

INDEX.md row: `findings=1, patches=0, MCAs=1, no-patch=0`.

Filed by ε@cnos (γ+α+β-collapsed-on-δ) on 2026-05-24.
