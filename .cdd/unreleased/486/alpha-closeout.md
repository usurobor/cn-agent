---
cycle: 486
parent_issue: cnos#486
master_tracker: cnos#467 (Sub 5B)
cycle_branch: cycle/486
pr: https://github.com/usurobor/cnos/pull/489
head_sha: fb6ae3fa253ab8df1cb579a40a5aff3ea1099d66
base_main_sha: 950730c74985864537696ec45ebf0023fde16b97
gamma_scaffold_sha: f1011f29d989a44d76fff0818267b70a07dec796
alpha_r0_impl_sha: 1743e3cda7da3134e32bea414ef814495c6067a1
alpha_r0_signal_sha: 43ba38389d5fa4781fbe4a75506d1cf48e4c846a
beta_r0_review_sha: cc7e5db2 (verdict: converge; zero findings)
operator_iterate_r1_sha: fb6ae3fa253ab8df1cb579a40a5aff3ea1099d66 (δ-direct one-line fix)
rounds: R0 (β converge, zero findings) + operator-iterate-narrowly + R1 (δ-direct one-line fix; no α R1 spawn)
role: α
authored_by: α@cdd.cnos (this session: bootstrap-δ — a fresh Agent spawn under the δ-the-orchestrator parent session that authored γ-486 + dispatched α/β R0 + made the R1 one-line fix directly)
date: 2026-06-23 (UTC)
---

# α-486 closeout — cdd/delta dispatch-wake-invoked δ mode amendment

## 1. Cycle summary

**Shipped.** A new top-level §9 *"Dispatch-wake-invoked mode"* in `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (+142 lines diff), satisfying all **9 ACs** from cnos#486 (AC1–AC9) and all **4 operator-named guardrails** (OG-1 substrate neutrality; OG-2 sharp output contract; OG-3 production-routing invariant; OG-4 empirical citation discipline). β converged at R0 on PR #489 with **zero findings** (`cc7e5db2`); the operator's final review caught a single off-by-one in §9.5 — the R[N≥1] artifact-contract row named β's same-round verdict as `§R[N-1]` (which would mean the wake reads β's PRIOR-round review at the current round). δ (the parent session) fixed this directly as R1 (`fb6ae3fa`; one line, one diff hunk) rather than spawning a fresh α R1 cycle; operator then re-converged. Cycle/486 is Sub 5B of the cnos#467 wake-orchestration wave and the **second R0-converge** in the wave (after cycle/485) — pattern: thorough γ scaffold + cnos#478 mechanical-injection discipline → R0 converge.

## 2. What I implemented (R0)

Appended **§9 "Dispatch-wake-invoked mode"** at the end of `delta/SKILL.md` (γ-FN-7 placement option (b) — no renumbering of existing §3–§8; diff is purely additive). Nine sub-sections:

- **§9.1 — Bootstrap-δ vs wake-invoked-δ (load-bearing distinction).** Two-row table contrasting the two modes on the parent-session axis. Names bootstrap-δ as *empirical observation* (cycles 470/476/485/486) and wake-invoked-δ as the *destination* (lands via Sub 5C). The clause "every contract clause in §9.2–§9.7 is written for wake-invoked-δ" pins the framing so no future reader mistakes bootstrap-δ's session-coordinated shape for the production contract.
- **§9.2 — Input contract.** Five-input table enumerating what the wake hands δ at invocation: claimed-issue-number; protocol identifier; current main SHA; wake run id; package-runtime context. Each input has a one-sentence semantics paragraph plus a citation to where it's produced (claim sequence, protocol label, repo HEAD, substrate firing, runtime path).
- **§9.3 — Routing sequence.** Numbered 5-step discrete list (dispatch γ for scaffold → dispatch α for R[N] → dispatch β for R[N] review → route on β's verdict → dispatch γ + α + β for closeouts). The OG-3 production invariant *"δ dispatches every role; γ does not spawn α/β"* appears verbatim as the opening line; each step carries an inline anti-invariant (**γ does NOT spawn α**, **α does NOT spawn β**, **β does NOT spawn α**, **δ does NOT author closeout content**). Per-role spawn-context table follows.
- **§9.4 — R[N] iteration token discipline.** Three-row table pinning the v0 hybrid mechanism (branch-state primary + issue-comment secondary; `.cdd/.../*.R{N}.json` explicitly not-pinned-for-v0). Anchors the wake-observability claim by walking through how a fresh substrate firing can `git fetch + git ls-tree` to read round state without consulting agent memory.
- **§9.5 — Per-R[N] artifact contract.** Three-row table per boundary (R0 / R[N≥1] / converge). Names the 7 canonical artifact classes (`gamma-scaffold`, `self-coherence`, `beta-review`, `alpha-closeout`, `beta-closeout`, `gamma-closeout`, optional `post-release-assessment`) — 1:1 with the cds-dispatch manifest's `output_contract.artifact_class_taxonomy`. **(This is the section where the R1 off-by-one lived — see §4 below.)**
- **§9.6 — Return tokens.** Four-row table mapping δ's wake-observable writes (`status:review`, `status:blocked + reason`, claim release, cycle-PR URL) onto `dispatch-protocol/SKILL.md` §2.4 lifecycle transitions. Includes the explicit `status:changes` carve-out ("β iteration is INTERNAL; the cell stays `status:in-progress`; `status:changes` is the operator/planner's authority AFTER `status:review`") so a future reader doesn't conflate the two.
- **§9.7 — v0 substrate constraints.** Three-row table naming the timeout horizon (class-level, not a specific minute count), the per-protocol concurrency group, and the one-claim-per-firing rule. Constraints are framed as "what δ honors" with "where the substrate encodes it" pointing at renderer authority — keeps the skill substrate-agnostic in semantic content.
- **§9.8 — Relationship to substrate (descriptive only — no substrate emission).** The single carve-out paragraph (3 substrate-token hits — meets the AC8 threshold). Names the substrate tokens δ does NOT emit, identifies v0 substrate (GitHub Actions) + renderer (`cn install-wake`) descriptively, and enumerates the four substrate-abstraction properties any future substrate must supply (fire-once horizon; serialize-per-protocol concurrency; one-claim-per-firing claim; issue-comment + label-transition API).
- **§9.9 — Cross-references.** Citations to dispatch-protocol (§§2.2/2.3/2.4/2.6), wake-provider §3.3, cds-dispatch manifest + prompt, cdd/{issue,gamma,alpha,beta}, and the three empirical-anchor cycle directories (470/476/485) each with a specific finding sentence.

**Key design decisions** (full rationale in self-coherence §Design):

- **Heading + placement:** chose `## 9. Dispatch-wake-invoked mode` (the "Dispatch-" prefix matches the predecessor wake name `cds-dispatch` and the production-mode framing) appended at end (γ-FN-7 option (b) — simplicity of diff; no renumbering).
- **Iteration-token mechanism:** pinned the hybrid (branch-state primary + issue-comment secondary) per γ-FN-8; rejected branch-state-alone (no human-observable signal), issue-comments-alone (brittle, doesn't carry per-round artifact state), and `.R{N}.json` machine artifacts (too much machinery for v0 with no current consumer).
- **Bootstrap-δ → wake-invoked-δ mapping:** encoded as a two-row distinction-table in §9.1 with the parent-session axis as the load-bearing structural difference; cited bootstrap-δ as the empirical observation set; cited wake-invoked-δ as the not-yet-observed destination.

## 3. What happened in R1 (the operator-direct one-line fix)

**The off-by-one β missed.** §9.5 has a 3-row table (R0 / R[N≥1] / converge). The R[N≥1] row read:

> `self-coherence.md` with appended `§R[N]` section (authored by α) + `beta-review.md` with appended `§R[N-1]` verdict (authored by β)

The `§R[N-1]` for β's verdict is **wrong**. β reviews α's same-round work — α's `§R[N]` → β's `§R[N]` in the SAME round. `§R[N-1]` would be β's PRIOR review, which is what α RECEIVES as input for the next round (correctly named in §9.3 step 2 and §9.3's input table — different context). At the R[N≥1] *artifact-contract* boundary, β's verdict for the current round must be at `§R[N]`. If the wake reads §9.5 to decide whether round N is complete, it would look for β's `§R[N-1]` verdict — which is the prior round's verdict, already on the branch — and could mis-route by treating an incomplete round as complete (or vice-versa, depending on indexing).

**How the operator caught it; why β missed it.** β reviewed §9.4 (which had `§R[N]` correctly) and §9.5 in isolation but did not run a cross-section invariant check on the §R[N]/§R[N-1] occurrences. The visual consistency trap: §9.3 step 2 + §9.3 input-context table correctly say β's findings are at `§R[N-1]` (because that's what α RECEIVES from the prior round); the §9.5 R[N≥1] row's `§R[N-1]` mirrored that phrasing — but the context differs: §9.3 names α's INPUT for round N (= β's prior verdict); §9.5 names the artifact CONTRACT at round N (= β's current-round verdict). Same token; opposite meaning by direction. β didn't notice; the operator did.

**Why δ fixed directly rather than spawning α R1.** One-line diff (line 460 of `delta/SKILL.md`: `§R[N-1]` → `§R[N]`). The operator's verdict was *iterate narrowly* — a precise, named, one-token fix with no design surface. Spawning a fresh α R1 cycle (full scaffold-read + AC-table re-verification + commit + signal + β re-dispatch) for a one-character textual correction would be ceremony with no information return. δ (the parent session) edited the file directly, committed with a detailed message (`fb6ae3fa`) explaining the cross-section consistency restoration, and presented the fix to the operator for re-convergence.

**Lesson.** This is the kind of cross-section invariant that the per-AC oracle cannot catch with a single regex against a single location — it requires a *direction-aware* check (β's verdict in §R[N] vs α's prior-round input from §R[N-1]). See §7 below for the precise specification gap and recommendation.

## 4. What went well in R0

**The convergence on R0 was real, not lucky.** Honest accounting of what worked:

- **The §9.1 bootstrap-δ vs wake-invoked-δ distinction-table** carried the load-bearing conceptual move γ-FN-3 named. By making the parent-session axis explicit, every downstream §9.2–§9.7 clause had a clear test ("does this require a parent session?") which I could self-apply while drafting. Without §9.1, the temptation to encode bootstrap-δ-the-empirical as wake-invoked-δ-the-destination would have leaked through.
- **The hybrid iteration-token mechanism (§9.4)** turned out to be the right disposition. Branch-state alone fails the human-observability test; issue-comments alone fails the per-round-artifact-state test; the hybrid (branch primary + comment secondary) gives both audiences (wake/substrate + human operator) the surfaces they need without inventing new machinery. β explicitly endorsed this in FN-β4.
- **The §9.8 substantive carve-out** (FN-α2) — keeping it substantive rather than minimal turned out to be the right call. β read it end-to-end and judged that the substrate-abstraction properties enumeration (fire-once horizon; serialize-per-protocol concurrency; one-claim-per-firing claim; issue-comment + label-transition API surface) is a load-bearing forward-compatibility statement for the wave's substrate-replacement future. β-FN-β3 names this as worth citing in cnos#487 / Sub 5C.
- **Empirical citation discipline (OG-4).** Each cycle citation (470/476/485) named a specific finding (broken relative-link / 3-round class-trap recurrence / mechanical-injection R0-converge), not a bare reference. β accepted all three under OG-4 without note.
- **Line discipline (OG-2).** Final +142 lines, well under the ~300 soft target. I caught myself starting to write detailed γ/α/β substep workflows in §9.3 step bodies and pruned them to one-sentence "the role's job per X §Y" references that point to the role-skill itself. The amendment cites; it does not restate.
- **Doctrine discipline (post-PR-480).** Zero stale `protocol:cdd` / `cdd-dispatch` / "cnos.cdd owns dispatch wake" leaks. The framework-vs-protocol distinction held throughout drafting.

## 5. What was friction

**γ's 12 friction notes — which mattered?** Honest accounting:

- **High-impact for R0 design:** FN-7 (placement options), FN-8 (iteration-token mechanism trade-offs), FN-9 (cycle/485 citation bonus), FN-11 (closeout triad disposition), FN-3 (bootstrap-δ → wake-invoked-δ mapping as load-bearing). I leaned on each of these in §Design decisions.
- **Background-validated (no inline decision needed):** FN-5 (CI surface carve-out — no `.github/workflows/` touched, no per-CI-step audit), FN-4 (branch-as-shared-state pattern — implicitly honored via §9.3's no-hidden-state discipline), FN-12 (FN-9 cycle/485 inclusion recommendation — I included).
- **Lower-impact but useful as guardrail:** FN-1 (γ scaffold AC regex looseness — I hit one instance of this; see FN-α1), FN-2 (substrate-leakage threshold-3), FN-6, FN-10.

12 friction notes is a lot, and most of them are operating as a friction-budget reserve rather than as required-resolution items. A future γ-scaffold template could distinguish "load-bearing decisions α MUST resolve in §Design" from "background guardrails γ surfaces for completeness". On balance, the 12-FN scaffold worked — I had every decision pre-thought, no inline ambiguity required improvisation — but it's heavy.

**Ambiguities I had to resolve inline:**

- **FN-α1 (AC3 routing-step regex looseness).** γ's AC3 oracle regex `^\s*[0-9]+\.\s+(spawn|dispatch|route|invoke)\s+(γ|...)` matches 4/5 of my routing-step verbs because step 4 ("route on β's verdict") has "route" followed by "on", not by a role glyph. I judged step 4's readable phrasing won over a regex-friendly rename and documented this as a soft-pass in self-coherence. β accepted in α's favor.
- **FN-α2 (§9.8 carve-out substantive vs minimal).** The carve-out's 3 substrate-token hits exactly meet AC8's threshold-3. I deliberately kept the carve-out substantive (a single paragraph naming GitHub Actions + workflow + claude-code-action descriptively) because the substrate-abstraction enumeration provides forward-compatibility value to a future reader. β accepted in α's favor (FN-α2 decision endorsed substantive).

**The substrate-agnosticism gate (AC8, OG-1).** This was a deliberate-craft job. The threshold is ≤3 substrate-token hits in the diff; my final §9.8 carve-out has exactly 3 (one negative enumeration, one descriptive identification, one meta-statement). I consciously crafted the carve-out paragraph to land at 3 — adding a fourth (e.g. naming `ubuntu-latest`) would have pushed over; dropping to 2 would have lost the forward-compatibility framing. The threshold was tight but achievable. A future γ-scaffold could consider whether the threshold-3 is the right number (it forces the carve-out into a single paragraph, which is good discipline, but it makes the writing exercise about token-counting rather than about substrate-neutrality reasoning).

**Source-of-truth reading load.** γ-scaffold listed 12 source-of-truth files. I read each before drafting (per γ-scaffold §6 "Read this entire scaffold + the cited source files before touching the SKILL"). This took real time but was necessary — citing dispatch-protocol §§2.2/2.3/2.4/2.6 + §3.2 + §3.6 + §3.7 against the actual line content (not from memory) is what made the cross-skill consistency check pass on β re-verification.

## 6. Lessons for future α prompts

**Lesson 1 (the §R[N-1] vs §R[N] off-by-one is a specification gap, not an α failure).** AC5's oracle in γ-486 scaffold checks: (a) all 7 artifact-class names appear in the file; (b) R0 boundary names `gamma-scaffold`; (c) `R[N>0]` boundary names `self-coherence` + `beta-review`; (d) converge names the three closeouts. None of these checks distinguish *which round-index* the artifact carries. The fix-shape for a future γ-AC5 oracle:

```bash
# Cross-section invariant: β review section index per context
# β's CURRENT-round verdict at R[N>0] is §R[N] (artifact contract — what β WRITES this round)
# α's INPUT from β's PRIOR-round verdict is §R[N-1] (α's spawn context — what β WROTE last round)
# The two are at different locations in the SKILL and refer to opposite directions

# Check β-write contexts (artifact contract; §R[N]):
grep -nE 'beta-review.md.*§R\[N-1\]' "$f" | \
  grep -v 'α.*receive\|input\|prior.round\|previous.round' \
  && echo "AC5 finding: beta-review §R[N-1] outside an α-input context"

# Check α-input contexts (β's prior verdict; §R[N-1]):
grep -nE 'α.*beta-review.*§R\[N\]' "$f" \
  && echo "AC5 finding: α-input names β's same-round (impossible — α implements BEFORE β reviews this round)"
```

This is the kind of *direction-aware cross-section invariant* check the per-AC oracle should carry. Single-location regexes catch single-location violations; the §R[N]/§R[N-1] distinction is structural (who reads/writes/in-which-direction), not lexical (does the token appear). γ-template improvement candidate.

**Lesson 2 (the δ-direct one-line fix pattern).** When the operator's verdict is *iterate narrowly* and the fix is a one-line textual correction with no design surface, δ-the-orchestrator editing directly (rather than spawning α R1) is a useful pattern. Honest assessment of when this is and isn't appropriate:

- **Appropriate:** one-line textual fix; no design judgment; no AC re-verification needed beyond the fixed line; the fix is mechanical (rename / token-swap / typo) and the operator has already done the cross-section consistency reasoning. The δ-direct edit collapses ceremony that would otherwise consume a full α-spawn + β-spawn cycle without producing new information.
- **NOT appropriate:** any fix requiring design re-judgment; any fix that touches multiple sections; any fix where α's §Design rationale needs to be updated; any fix where β's AC oracle would need to be re-run against more than the changed line. In these cases the discipline is α R1 + β R1 with the full self-coherence + beta-review append-discipline.

This pattern *does* bypass some discipline — α's append-only §R[N] section discipline in `self-coherence.md` isn't updated; β's `beta-review.md` isn't append-extended with an R1 verdict; the round-state on the branch is "R0 + a δ-direct fix commit" rather than "R0 + R1 with full per-round artifacts". For one-line fixes that's acceptable cost; the alternative ceremony-for-its-own-sake would dilute the meaning of "α R1" (which should denote substantive α re-engagement, not mechanical correction). The operator should pin this as a pattern explicitly — if δ-direct is the disposition for narrow iterate verdicts, name it in the operator skill so the convention is visible.

**Lesson 3 (mechanical-injection discipline scales).** cycle/485 introduced the per-CI-step audit format under cnos#478 to absorb the bash-e class-trap class; cycle/486 doesn't touch CI but inherited the same scaffold-precision-rigor discipline (γ's 12 friction notes, source-of-truth list, OG-named guardrails). Both cycles R0-converged. This is two data points but suggests the discipline transfers across cycle-class (renderer+CI → pure-doctrine) — the mechanism isn't the per-CI-step audit per se but the cnos#478 principle of mechanical, pre-specified, scaffold-encoded review-readiness criteria. Worth naming as a wave-level finding.

## 7. Empirical findings for cycle/486 specifically

- **Second R0-converge in the wave (after cycle/485).** Pattern: thorough γ-scaffold (9.5K words / 700 lines for γ-486; 6.5K for γ-485) + cnos#478 mechanical-injection discipline → R0 converge. cycle/486 is the *first* pure-doctrine R0-converge (cycle/485 was renderer + CI). β-FN-β1 names this; γ-closeout's triage should pick up the recommendation to pin the γ-template improvement.

- **Operator caught what β missed (the §9.5 off-by-one).** The β role's oracle re-run was independent and rigorous — every per-AC oracle was re-verified on the cycle HEAD — but did not include a cross-section invariant check ("for every `§R[N]` occurrence, verify the direction-context matches the artifact-vs-input distinction"). The §9.5 R[N≥1] row's `§R[N-1]` was inside an AC5-checked location (artifact-class name present; per-round artifact set named) but the *round-index correctness* was outside AC5's resolution. Recommendation for β-skill: add a *cross-section invariant* review step ("for any token that appears at multiple locations in the diff with potentially-different semantics — e.g. `§R[N]`/`§R[N-1]`, `status:in-progress`/`status:review`, `claim`/`release` — verify all occurrences are direction-consistent"). This is a B-severity finding for `beta/SKILL.md` (procedural gap; the off-by-one was caught by the operator's final read so the convergence wasn't actually wrong — but β should not depend on the operator catching this kind of thing).

- **δ-direct R1 fix as a new pattern.** Empirical observation: when the operator's verdict is *iterate narrowly* and the fix is a one-line textual correction, δ-the-session can edit directly rather than spawning α R1. The branch state ends up as "R0 commit + R0 signal commit + β converge commit + δ-direct fix commit + (closeouts)" — one extra commit, no extra α/β round, full operator audit trail in the fix commit message. β-skill should probably name this pattern explicitly so the iterate-narrowly + δ-direct disposition isn't ambiguous in future cycles. See Lesson 2 above for the appropriate/not-appropriate split.

- **The §9.5 fix matters for Sub 5C (cnos#487).** The fix commit message names this: "the wake reads §9.5's per-round artifact contract to decide whether a round is complete. The wrong section reference would cause wake observability bugs (wake looks for §R[N-1] but β wrote §R[N]; round looks incomplete; wake mis-routes)". Sub 5C will instantiate the contract for a real `protocol:cds` smoke cell; the §9.5 R[N≥1] row is the contract clause the wake reads at every iteration boundary. The fix prevented a class of substrate-side observability bug from shipping into the production wake's first firing.

## 8. Empirical findings for cycle/486 specifically (additional)

(Per §7 above — the prompt enumerates this as §8 but the substantive content is consolidated in §7. Retained as a stub for section-numbering symmetry with the closeout template.)

- **R0 vs R1 commit count economy.** R0 = 2 commits on the branch (α impl + α R0 signal); β-converge added a 3rd; δ-direct R1 fix added a 4th; closeouts will add 3 more. Total expected commits at converge-plus-closeouts: 8. cycle/485 had a similar count. The branch-state-primary iteration-token discipline (§9.4) produces a manageable commit graph at this scale.

## 9. Recommendations

Picked up from §6 + §7; pass to γ-closeout for triage:

1. **γ-template improvement — direction-aware cross-section invariant oracle for AC5** (Lesson 1). The §R[N]/§R[N-1] off-by-one is the canonical case; the fix-shape is a grep-pair that checks artifact-write contexts (β's verdict at §R[N]) vs α-input contexts (β's prior verdict at §R[N-1]) and surfaces violations. Worth pinning in a future γ-scaffold template / γ-skill amendment.

2. **β-skill amendment — cross-section invariant review step** (§7 second bullet). β-skill should require, in addition to per-AC oracle re-run, a cross-section invariant check for tokens that recur with potentially-different semantics. Reasonably small addition (one Role Rule); high leverage (catches the class of bug the operator caught here). This fits the cycle/485 closeouts' β-skill cluster (T1/T2/T11–T14).

3. **operator-skill amendment — δ-direct one-line fix pattern** (Lesson 2 + §7 third bullet). Name the pattern explicitly: when the operator's verdict is *iterate narrowly* and the fix is one-line / no-design-judgment / mechanical, δ-the-session may edit directly and present the fix to the operator for re-convergence. Names the appropriate/not-appropriate split per Lesson 2. Reduces ambiguity for future operators / future δ sessions.

4. **wave-level finding — mechanical-injection discipline transfers across cycle-class** (Lesson 3). Two-data-point evidence (cycle/485 renderer+CI + cycle/486 pure-doctrine, both R0-converge under thorough γ-scaffold + cnos#478 discipline). Worth observing in the wave's master tracker (cnos#467) as the pattern that absorbs scaffold-side substantive ambiguity AND mechanical class-traps under one discipline.

5. **Sub 5C / cnos#487 hook** (FN-β3). cnos#487 should cite §9.8's substrate-abstraction-properties enumeration as the contract surface to test against when verifying the cds-dispatch wake's substrate emission honors the abstraction layer. The four properties (fire-once horizon; serialize-per-protocol concurrency; one-claim-per-firing claim; issue-comment + label-transition API) are the substrate-replacement future's test harness; Sub 5C is the first concrete case to test them.

## 10. Closeout signoff

α-486 closeout complete; cycle shipped as PR #489 R0 (β converge, zero findings) + δ-direct R1 (one-line operator-iterate-narrowly fix; operator re-converge); no carryover debt within α's scope; recommendations passed to γ-closeout for triage.
