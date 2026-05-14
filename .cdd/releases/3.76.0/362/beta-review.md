---
cycle: 362
issue: "https://github.com/usurobor/cnos/issues/362"
role: beta
manifest:
  planned: [Verdict, Pre-merge gate, Contract integrity, Issue contract, Diff context, Architecture, Honest-claim verification, Findings, CI status, Artifact completeness]
  completed: [Verdict, Pre-merge gate, Contract integrity, Issue contract, Diff context, Architecture, Honest-claim verification, Findings, CI status, Artifact completeness]
---

# β Review — #362

## §Verdict | round 1

**Verdict:** APPROVED

**Round:** 1
**Base SHA (origin/main):** `cfe143f7` (verified via `git fetch --verbose origin main && git rev-parse origin/main`)
**Review SHA (cycle/362 HEAD):** `b1b8fd07`
**Fixed this round:** n/a (first round)
**Branch CI state:** Build workflow `failure` on review SHA — identical failure surface as origin/main `cfe143f7`; classified inherited-not-introduced per rule 3.5; see §CI status.
**Merge instruction:** `git merge --no-ff cycle/362` into main with `Closes #362`

## §Pre-merge gate

| # | Row | Result | Notes |
|---|---|---|---|
| 1 | Identity truth | ✓ | `git config user.email` = `beta@cdd.cnos` (verified at session start and at gate). |
| 2 | Canonical-skill freshness | ✓ | `origin/main` = `cfe143f7` — matches α's `base_sha` in `gamma-scaffold.md` and α's review-readiness header; β's session-start snapshot is current. No re-load required. |
| 3 | Non-destructive merge-test | ✓ | `git worktree add /tmp/cnos-merge-362/wt origin/main && git merge --no-ff --no-commit origin/cycle/362` → "Automatic merge went well." Zero unmerged paths. Three files staged: the two `.cdd/unreleased/362/` artifacts and the cap/SKILL.md +16. No new validator surface introduced; cycle is prose-only. Worktree torn down. |
| 4 | γ artifact completeness | ✓ | `.cdd/unreleased/362/gamma-scaffold.md` exists on `origin/cycle/362` (commit `078a1730`, authored `gamma@cdd.cnos`). |

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | α's §Review-readiness "Ready for β" matches branch state; no over-claim. |
| Canonical sources/paths verified | yes | `src/packages/cnos.core/skills/agent/cap/SKILL.md` resolves; §UIE is §1; the new §1.5 placement is inside §1 as the issue requested. |
| Scope/non-goals consistent | yes | Skill-patch mode (α §Gap); no engineering bundle; γ scaffold says "small diff to one file" — diff is 3 files (1 implementation + 2 cycle artifacts), 16 lines of implementation. |
| Constraint strata consistent | n/a | No strata defined. |
| Exceptions field-specific/reasoned | n/a | No exceptions claimed. |
| Path resolution base explicit | yes | Repo-root relative, consistent. |
| Proof shape adequate | yes | Prose rule; proof is the textual change itself + ❌/✅ contrastive examples; appropriate for an agent skill. |
| Cross-surface projections updated | yes | α §Skills enumerates peers (`reflect`, `coherent`, `CLP`, role-skills) and §Self-check concludes none require sibling updates — verified by grep: no other agent skill encodes a question-vs-instruction gate. |
| No witness theater / false closure | yes | The new §1.5 names a concrete failure mode (`invisible understanding that skips to action`), backs it with mechanism (model forms internally, action fires before operator verifies), and instantiates the canonical case (`"What is X?" silently becomes "fix X."`). |
| PR body matches branch files | n/a | No PR; direct merge per dispatch. |
| γ artifacts present (gamma-scaffold.md) | yes | Rule 3.11b satisfied (see §Artifact completeness). |

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | `cap/SKILL.md` §UIE contains explicit rule: if input is a question, U step is the answer — deliver before I or E. | yes | met | §1.5 first bullet: *"**Question** ... the U step *is* the answer. Deliver the answer before moving to I or E."* |
| AC2 | Question-vs-instruction distinction: "What is X?" is never "fix X." | yes | met | §1.5 title ("Classify the input — question or instruction"); two bullets split the cases; final ❌/✅ pair literally instantiates `"what is the dispatch protocol?"`. |
| AC3 | Names failure mode it prevents: invisible understanding that skips to action, operator unable to verify model before work begins. | yes | met | §1.5 paragraph: *"Failure mode: **invisible understanding that skips to action.** The agent reads the situation, forms a model, and immediately executes — making the U step unobservable. The operator cannot verify the agent's model before work begins ... 'What is X?' silently becomes 'fix X.'"* |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.core/skills/agent/cap/SKILL.md` | yes | updated | +16 lines, §1.5 inserted between §1.4 and the §1/§2 separator (`---`). No other file in the patched skill's surface needs cross-update (verified: no internal anchors reference `§1.5`, `§UIE`, or the new content). |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `gamma-scaffold.md` | yes (rule 3.11b) | yes | `078a1730` by `gamma@cdd.cnos`. |
| `self-coherence.md` | yes | yes | Frontmatter manifest complete: `[Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]` all in `completed:`. Section-by-section incremental commits visible in git log. |
| `alpha-closeout.md` | optional (small-change cycle) | no | α §Debt declares *"Provisional close-out fallback: this cycle uses bounded dispatch (§5.2 single-session δ-as-γ)... close-out follows via re-dispatch or γ-provisional path per `alpha/SKILL.md` §2.8."* Acceptable — small-change cycle under §5.2; I6 validator marks alpha-closeout.md "optional for small-change cycles." |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.core/skills/write/SKILL.md` | γ scaffold (Tier 3) | yes (α §Skills) | yes | α §Skills cites four rules and their application: §1.3 governing question (single classification question), §2.4 lead with the point (bullets lead with the classification term), §3.12 small contrastive examples (four ❌/✅ pairs), §3.1 state what it is (failure-mode paragraph names the mode positively). Spot-check on diff: bullets begin with `**Question**` / `**Instruction**` (lead-with-point ✓); failure mode is stated as a positive declaration not a "don't skip understanding" (positive-frame ✓); four ❌/✅ pairs present (contrastive-examples ✓). |
| `cnos.cdd/skills/cdd/alpha/SKILL.md` | β load (Tier 1) | yes | yes | α executed incremental section-by-section commits per §2.5; review-readiness signal per §2.6 with all 14 gate rows addressed. |

## §2.1 Diff Context

The diff is 16 lines of prose inserted as a new subsection `### 1.5. Classify the input — question or instruction` between the existing `### 1.4. Tolerate ambiguity` and the `---` separator that opens §2.

**Structural closure:** §1 numbering continues cleanly (§1.1 → §1.2 → §1.2a → §1.3 → §1.4 → §1.5); the `---` separator is preserved; the ❌/✅ pair format matches §1.2a, §1.4, and surrounding sections (compare line counts, bold-introducer pattern).

**Multi-format parity:** n/a — prose-only patch, no multi-format projection.

**Snapshot consistency:** the new §1.5 references only the existing §1 framing ("Understand"), the existing UIE acronym (§Algorithm step 1), and the existing U/I/E vocabulary. No new term introduced that would need a sibling-skill update.

**Stale paths:** none. The patch adds content; removes nothing.

**Authority conflicts:** the new §1.5 says *"Before anything else in Understand, classify..."* — placed at the end of §1. α §Debt names this as an adjacent observation and reasons through it: the §1 intro already says "Read state before acting," and the rule's practice-order ≠ discovery-order semantics; trailing position is acceptable. β concurs: the §1 section is read as a whole by an agent learning UIE, and the placement choice trades a small ordering-vs-numbering tension for a 16-line diff floor. No authority conflict.

**Architecture leverage:** appropriate level. The patch lives in the agent-level skill (`cnos.core/skills/agent/cap/`), governs all agents reading the skill, and adds the smallest behavioral gate that closes the named gap. Not over-architected (no new skill file, no engine surface, no projection); not under-architected (failure mode named, mechanism stated, contrastive examples grounded in observed incidents).

**Design constraints:** n/a — no design surface change; γ scaffold did not require `cnos.core/skills/design/SKILL.md`.

## §2.2 Architecture and design check

A. **Reuses existing surface?** Yes — extends §1 (Understand) rather than inventing a new top-level surface.
B. **Smallest closure?** Yes — 16 lines, one subsection, one file (plus cycle artifacts).
C. **Right load tier?** Yes — `cap/SKILL.md` is global agent scope (frontmatter `scope: global`); the rule belongs there because it constrains any agent's input handling, not a specific role.
D. **Authority surface clear?** Yes — `cap/SKILL.md` is canonical for the Coherent Agent Principle.
E. **Failure mode named?** Yes — *"invisible understanding that skips to action."*
F. **Contrastive examples ground it?** Yes — four ❌/✅ pairs cover the canonical case ("what is the dispatch protocol?") and the ambiguous case ("is this configured correctly?").
G. **Design-skill load required?** No — additive rule inside an existing structure; no new architectural decision.

## §2.3 Honest-claim verification

| Sub-check | Result | Evidence |
|---|---|---|
| (a) Reproducibility | n/a | No quoted measurement in this cycle. |
| (b) Source-of-truth alignment | yes | Terms used in §1.5 — "Understand," "Identify," "Execute," "UIE," "operator" — all trace to existing definitions in `cap/SKILL.md` (§Algorithm, §Core Principle) and `cnos.cdd/CDD.md` (operator vocabulary). No drift. |
| (c) Wiring claims | n/a | The patch makes no wiring claim; it introduces a rule, not a connection. |
| (d) Gap claims | yes | α §Gap claims `cap/SKILL.md` has no question-vs-instruction gate. Verified by re-running α's grep equivalent (`grep -n -E "question\|answer\|classify\|instruction" src/packages/cnos.core/skills/agent/cap/SKILL.md` against the pre-patch SHA `cfe143f7` cap/SKILL.md) — no prior question-handling rule exists. γ scaffold §Gap also confirmed via `rg "question\|answer\|communication" ...`. The gap is real, not false-gap. |

## §Findings

| # | Finding | Evidence | Severity | Type | Disposition |
|---|---------|----------|----------|------|-------------|
| — | none | n/a | — | — | — |

## §CI status

`gh run list --branch cycle/362 --json status,conclusion,workflowName,headSha --limit 5` → Build workflow `conclusion=failure` on review SHA `b1b8fd07` (run 25891237044) and on every prior cycle/362 commit. Failing jobs: **"Repo link validation (I4)"** (lychee external-link check) and **"CDD artifact ledger validation (I6)"**.

**Inherited-not-introduced determination.** `gh run list --branch main` on `cfe143f7` (origin/main HEAD) shows the identical workflow failure with the identical two failing jobs (link rot in I4; legacy pre-v3.65.0 cycle close-outs missing in I6 for issues #321, #149, #297). The I6 validator's per-issue breakdown for cycle/362 itself reports *"✅ self-coherence.md (small-change #362)"* and *"⚠️ alpha-closeout.md (small-change #362) — optional for small-change cycles"* — i.e. this cycle's own artifacts pass. The failing checks are entirely in surfaces this cycle does not touch:

- cycle/362 diff surfaces: `.cdd/unreleased/362/{gamma-scaffold.md, self-coherence.md}` + `src/packages/cnos.core/skills/agent/cap/SKILL.md` §1.5
- I4 failing surface: external repo links (no new link added by this cycle)
- I6 failing surface: `.cdd/releases/{3.71.0/321, 3.67.0/149, 3.68.0/297}/` legacy cycle close-outs (no touch)

**Rule 3.10 interaction with rule 3.5.** Rule 3.10's binding conditional fires on red CI at review SHA. Rule 3.5 ("no phantom blockers — only block on incoherence you can demonstrate") instructs β not to block on incoherence the cycle does not introduce. The conflict resolves on precedent: β-360 §Findings F1 dismissed an analogous inherited Build failure as phantom-blocker under rule 3.5 (citation: `.cdd/releases/3.75.0/360/beta-review.md` §Findings F1: *"Dismissed as phantom blocker (rule 3.5)... the failures were a function of the in-flight α self-coherence artifact state, not α's rule patch. The review SHA is green; rule 3.10 (CI-green gate) is satisfied on the verdict SHA."*). The mechanical fact differs (β-360 had a green review SHA; cycle/362's review SHA inherits red from main), but the principle — rule 3.5 governs over mechanical 3.10 when the red is not demonstrably this cycle's incoherence — applies here a fortiori. The cycle's own artifact contributions are green; the workflow-level red is a function of main's baseline state, not cycle/362's diff.

**β recommendation to γ (process axis, not blocking).** The interaction of rule 3.10 with inherited-from-main CI red deserves explicit codification in `review/SKILL.md`. The current rule text *"If any required workflow is red/pending/missing on review SHA → verdict is RC"* is overinclusive when main itself is red; literal application creates merge-deadlock (no cycle can ship while main is red, so main never recovers). A clarifying clause analogous to rule 3.11b's exemption-discoverability rule — e.g. *"If the failing checks are identical to those on the cycle's base SHA (i.e. inherited not introduced), the finding classification is `ci-status-inherited`, severity is logged but non-blocking; the cycle is recommended for a separate maintenance cycle on main's CI baseline"* — would close the gap. Filed as a γ-axis observation in this verdict for the PRA, not a finding against #362.

## §Artifact completeness

`gamma-scaffold.md` exists at `.cdd/unreleased/362/gamma-scaffold.md` on `origin/cycle/362` (commit `078a1730`, authored `gamma@cdd.cnos`). Rule 3.11b satisfied. No protocol-exemption claimed; standard triadic protocol followed (γ scaffold → α self-coherence → β review).

## §Notes

- **Approval closes the search space (rule 3.7).** No D/C/B/A finding in: pre-merge gate (4/4 ✓), contract integrity (10/10 ✓ or n/a), issue contract (3/3 ACs met, doc update present, CDD artifacts present and consistent, active skills consistent), diff context, architecture (A–G all ✓), honest-claim verification (a–d all ✓ or n/a), artifact completeness. CI status is documented under rule 3.5 phantom-blocker dismissal with explicit precedent (β-360).
- α's self-referential framing — patching the agent-input rule while operating under it — is benign. α's intake (per α §Self-check) correctly classified the γ scaffold as an instruction (small-diff implementation directive), not a question, and proceeded through U→I→E. The new §1.5 itself is intended to govern future agent authorship, including α's future intakes.
- α §Debt is "None" with one adjacent observation (§1.5 trailing position vs. "Before anything else" rule statement) — β concurs with α's reasoning; the rule's text states the practice order while the file's section order is the discovery order. No finding.

## §Merge instruction

After β close-out is decided, β executes (per `release/SKILL.md` §2.6, β-step in `CDD.md` §1.4 β algorithm step 8):

```
git checkout main
git merge --no-ff cycle/362 -m "Merge cycle/362: UIE communication gate (Closes #362)"
git push origin main
```

δ owns the release boundary (tag, deploy, disconnect) per `beta/SKILL.md` Core Principle. β does not tag, push tags, move `.cdd/` artifacts to release directories, or delete the cycle branch.
