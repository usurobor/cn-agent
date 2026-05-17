<!-- sections: [Round1] -->
<!-- completed: [Round1] -->

# β review — #370

## Round 1

**Verdict:** REQUEST CHANGES

**Round:** 1
**Round opened:** 2026-05-17
**Review-base SHA (origin/main):** `704365d23378fcbfcf1e33679025809af6b81100`
**Review SHA (cycle/370 head):** `aa10f902a9446d6a0253a49e71a1896d9605e86f`
**Branch CI state:** red — `Build` workflow failing on review SHA (run 25990518085, conclusion `failure`, event `push`, 2026-05-17T12:13:00Z)
**β identity:** `beta@cdd.cnos` (worktree-config; re-asserted after merge-test pollution; verified `git config --get user.email` returns `beta@cdd.cnos` before the verdict commit)
**Merge instruction (deferred until R2 verdict APPROVED):** `git merge --no-ff cycle/370` into `main` with `Closes #370` in the merge commit body.

### §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Issue body labels this as "doctrine cycle, companion to #364" and "produces COHERENCE-CELL-NORMAL-FORM.md"; doc's frontmatter declares `**Status:** Draft doctrine. Phase 1.5 of #366`. Status framing is consistent across surfaces. |
| Canonical sources/paths verified | yes | All cited paths resolve: `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, `CDD.md`, `schemas/cdd/`, `cn-cdd-verify`, the four `## Kernel`/`## Cell Outcomes`/`## Recursion Modes`/`## Scope-Lift` headers. |
| Scope/non-goals consistent | yes | §Non-goals enumerates the same prohibited-edit set the issue body's `## Non-goals` lists; `git diff --stat` confirms (AC9). |
| Constraint strata consistent | yes | Active design constraints in issue body (substrate-independence, two-layer separation, exact header names, 200–400 line target, no angle-bracket placeholders) match the doc's authoring discipline (§Preamble §"Authoring discipline"). |
| Exceptions field-specific/reasoned | n/a | Doctrine cycle; no exceptions invoked. |
| Path resolution base explicit | yes | All paths under `src/packages/cnos.cdd/skills/cdd/` are package-relative; cited consistently. |
| Proof shape adequate | yes | Per-AC oracles in issue body executed by α and recorded in `self-coherence.md` §ACs; β re-executed AC1, AC7, AC9 oracles independently (see Findings table). |
| Cross-surface projections updated | n/a | Docs-only kernel doc; no operator-visible projection surfaces to update. |
| No witness theater / false closure | yes | Each AC's evidence is content-bound (named section + verbatim quotes + oracle output), not structural-only. |
| PR body matches branch files | yes | No PR; cycle branch surface containment verified via `git diff origin/main..HEAD --stat` (see AC9). |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/370/gamma-scaffold.md` exists on `origin/cycle/370` (140 lines; γ pre-flagged nine failure modes, all of which α executed against). Rule 3.11b compliance: satisfied. |

### §2.0 Issue Contract

#### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| 1 | File exists at `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` | yes | PASS | β re-ran `test -f`: file exists, 435 lines, status `A` in diff stat. |
| 2 | Draft-doctrine status + predecessor + peer citations | yes | PASS | β counted token occurrences in the doc: `draft doctrine` ×2 (lines 14, 435), `COHERENCE-CELL.md` ×4, `RECEIPT-VALIDATION.md` ×6. Predecessor framing explicit at §Preamble §"Companion, not replacement" (lines 20–26). |
| 3 | Kernel as five-step recursion + evidence-binding rule | yes | PASS | §Kernel (lines 50–142) states five-step closure. β's signature `βₙ.review(contractₙ, matterₙ)` (line 60, 81) — evidence excluded. δ's signature `δₙ.decide(receiptₙ, verdictₙ)` (line 63, 123) — evidence excluded. Evidence-binding rule pinned as blockquote at line 98. γ.close and V remain distinct steps. (See observation 1 below for α's signature-variant decision.) |
| 4 | Four cell outcomes + verdict × decision preconditions | yes | PASS | §Cell Outcomes (lines 146–192) names exactly the four outcomes (`accepted`/`degraded`/`blocked`/`invalid`) with the preconditions the issue body specifies. `degraded` correctly bound to `verdict ≠ PASS ∧ decision = override` (line 154). `invalid` non-terminal at §"`invalid` is non-terminal" (line 167–173). Alignment with #369 AC4 explicit (line 176–179). |
| 5 | Two recursion modes — within-scope vs cross-scope | yes | PASS | §Recursion Modes (lines 195–263). Within-scope mode is the repair-dispatch sub-case (line 199–217): same scope `n`, parent γₙ re-emits after child accept. Cross-scope mode covers accept/release (PASS) and override (non-PASS), with scope advance (line 219–236). Same-operator framing at line 238–254. Non-mode guardrails at line 256–263. |
| 6 | Three scope-lift projections + projection-not-renaming + β/γ-no-upward-projection | yes | PASS | §Scope-Lift (lines 267–321). Three projections named at line 273–277 (closed cell → α-matter, δ → β-like discrimination, ε → γ-like coordination/evolution). Projection-not-renaming framing as blockquote at line 295–297. β/γ-no-upward-projection clause at §"β and γ have no upward projection" (line 301–307). Alignment with #369 AC3 explicit (line 309–311). |
| 7 | Substrate-independence (section-bounded oracle) | yes | PASS | β re-ran the issue body's awk + rg recipe verbatim. `wc -l /tmp/ccnf-kernel.md` → 275 (non-zero; section headers match awk pattern exactly). Headers extracted: `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, `## Scope-Lift` (exactly the canonical four). `rg -i '\b(github\|cue\|cn-cdd-verify\|cn dispatch\|claude\|gh)\b' /tmp/ccnf-kernel.md` → no matches. |
| 8 | Two-layer separation + four realization peers cited | yes | PASS | §Two-Layer Separation (lines 325–370). Kernel-is-the-*what* / realization-is-the-*how-on-this-substrate* framing as blockquote at line 329. Four realization peers cited under §Realization peers (lines 337–347): `RECEIPT-VALIDATION.md`, `schemas/cdd/` (Phase 2, #369), `cn-cdd-verify` (Phase 3, deferred), `CDD.md` (Phase 7 rewrite, deferred). Direction-of-dependency rule stated at line 362. |
| 9 | Non-goals respected — surface containment proved | yes | PASS | `git diff origin/main..HEAD --stat` shows exactly: `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (A, +435), `.cdd/unreleased/370/{alpha-codex-prompt,beta-codex-prompt,gamma-scaffold,self-coherence}.md`. No edits to `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, `CDD.md`, `schemas/cdd/`, `cn-cdd-verify`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `ROLES.md`, CI workflows. Doc body restates non-goals (§Non-goals, lines 374–408). |

#### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` | yes | A (+435) | The new kernel-layer doctrine doc this cycle produces. |
| `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` | no | unchanged | Predecessor doctrine — non-goal honoured. |
| `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` | no | unchanged | Receptor design — non-goal honoured. |
| `src/packages/cnos.cdd/skills/cdd/CDD.md` | no | unchanged | Phase 7 rewrite target — non-goal honoured. |

#### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/370/gamma-scaffold.md` | yes (3.11b binding gate) | yes | Created by γ at scaffold time; 140 lines; carries γ's gap, ACs reference, dispatch configuration, nine failure-modes-to-guard. |
| `.cdd/unreleased/370/self-coherence.md` | yes (β intake input) | yes | Authored by α across the cycle; final form at HEAD (`aa10f902`) carries the §Review-readiness signal. **Section-header drift detected — see Finding F1 below.** |
| `.cdd/unreleased/370/alpha-codex-prompt.md` | yes (γ dispatch artifact) | yes | γ-authored dispatch prompt for α. |
| `.cdd/unreleased/370/beta-codex-prompt.md` | yes (γ dispatch artifact) | yes | γ-authored dispatch prompt for β. |
| `.cdd/unreleased/370/beta-review.md` | yes (β output) | yes (this file) | Round 1 written here. |

#### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cdd/CDD.md` | β intake | yes | yes | Canonical algorithm; β β-algorithm steps 1–9 applied. |
| `cdd/beta/SKILL.md` | β intake | yes | yes | β role surface; pre-merge gate rows 1–4 exercised. |
| `cdd/review/SKILL.md` | β phase 1–3 | yes | yes | Three review phases run; verdict rules + verdict-shape lint applied. |
| `cdd/release/SKILL.md` | β step 8–9 (deferred to R2 APPROVE) | yes | n/a | Loaded for completeness; not exercised this round (no merge until findings clear). |
| `cnos.core/skills/design/SKILL.md` | review architecture check | yes | yes | Kernel-articulation discipline verified: substrate-independence enforced (AC7), two-layer separation declared (AC8), no operational-content leakage in kernel sections. |
| `cdd/COHERENCE-CELL.md` | predecessor doctrine | yes | yes | Read as source-of-truth for receipt rule, four-way separation, role-as-cell-function — predecessor unedited (verified). |
| `cdd/RECEIPT-VALIDATION.md` | typed-interface design peer | yes | yes | Read as source-of-truth for V's interface, verdict/decision distinction, δ-authoritative validation — receptor design unedited (verified). |

### §2.1 Pre-merge gate (`beta/SKILL.md`)

| # | Row | Result | Notes |
|---|---|---|---|
| 1 | **Identity truth** — `git config --get user.email` returns `beta@cdd.cnos` | PASS (after re-assertion) | At β intake, `git config user.name beta && git config user.email beta@cdd.cnos` succeeded. Mid-review, a non-destructive merge-test worktree (`/tmp/cnos-merge-test/wt`) enabled `extensions.worktreeConfig` on the *shared* repo config; after the merge-test was torn down, `git config --get user.email` on the main worktree began reading `alpha@cdd.cnos` from the shared `.git/config` (because the worktree-config layer had no `user.email` set for the main worktree). β re-asserted at worktree scope: `git config --worktree user.name beta && git config --worktree user.email beta@cdd.cnos`. Verified `git config --get user.email` returns `beta@cdd.cnos` before any β-authored commit on this branch. Audit trail will show `beta@cdd.cnos` on every β commit. |
| 2 | **Canonical-skill freshness** — `origin/main` matches β session-start snapshot | PASS | `git fetch --verbose origin main` returns `[up to date]`; `git rev-parse origin/main` → `704365d2`, matching γ-scaffold's recorded base SHA. No advance since β intake; no skill re-load required. |
| 3 | **Non-destructive merge-test** — build merge tree, run cycle's validator on merge tree | **FAIL** | Merge tree built cleanly in `/tmp/cnos-merge-test/wt` (`git merge --no-ff --no-commit origin/cycle/370` → "Automatic merge went well; stopped before committing"; zero unmerged paths; diff-stat matches `git diff origin/main..HEAD --stat`). AC7 awk-extract + substrate-scan on the merge tree → PASS (275-line kernel slice, no substrate matches). However, **running the cycle's own validator (`cn-cdd-verify --unreleased`) on the merge tree reports `❌ self-coherence.md sections — missing required sections: CDD Trace` and prints `❌ Cycle artifact verification FAILED`.** This is the same failure CI reports against review SHA `aa10f902`. Merge-test worktree torn down cleanly. See **Finding F1**. |
| 4 | **γ artifact completeness** — `gamma-scaffold.md` exists on cycle branch | PASS | `.cdd/unreleased/370/gamma-scaffold.md` present at HEAD (`140 +`); γ coordinated this cycle through the canonical triadic protocol. Rule 3.11b compliance satisfied. |

### §2.2 CI status (rule 3.10)

Required workflows on `cycle/370`: `Build` (only workflow that runs; no branch-protection rules configured, so fallback per rule 3.10 is "every workflow that runs on the cycle branch is required").

| Workflow | Conclusion on review SHA | Notes |
|---|---|---|
| `Build` (run id `25990518085`, headSha `aa10f902`) | **failure** | Job `CDD artifact ledger validation (I6)` step `Check CDD artifacts` runs `cn-cdd-verify --all` and exits the workflow on `❌ Cycle artifact verification FAILED`. Single failed check: `self-coherence.md sections — missing required sections: CDD Trace` (cycle #370 path). See **Finding F1**. |

Per rule 3.10: required workflow red on review SHA → verdict RC, B-severity, classification `ci-status`. This is the binding driver for RC R1.

### §2.3 Architecture / design (review/architecture/SKILL.md Q's A–G)

| Q | Reading | Notes |
|---|---|---|
| A — Boundaries clear | yes | Kernel sections are role/artifact/verdict/decision/scope only; realization-layer concerns confined to §Two-Layer Separation. |
| B — Single source of truth | yes | Kernel doctrine doc is additive; predecessor surfaces remain authoritative for their layers (receipt rule, receptor design); the kernel doc names the algorithmic layer the predecessor implies. Direction-of-dependency rule explicit (line 362): realization cites kernel; kernel does not cite realization. |
| C — Coupling appropriate | yes | The doc cites four realization peers but does not couple to their internal structure — it only positions them (`V` at step 4, δ at step 5, schemas type the four outcomes, `CDD.md` will expand the kernel). |
| D — Failure modes named | yes | γ-scaffold §"Failure modes to guard against" pre-flagged nine; the doc's authoring respected all nine, and α's §Self-check explicitly notes the substrate-leakage scan caught a `[#369](https://github.com/...)` URL in §Cell Outcomes that was replaced with a bare reference — the AC7 mechanism worked as designed during authoring. |
| E — Testability | yes | Per-AC oracles (AC1, AC2, AC7, AC9) are mechanically executable; β re-ran them. AC3–AC6, AC8 verified by prose evidence with verbatim quotes. |
| F — Observability | yes | The doc carries an explicit phase-inheritance table (§Closure, lines 426–434) — Phases 2–7 will cite specific kernel sections as their input contract. |
| G — Extensibility | yes | Substrate-independence in the kernel (AC7) keeps the doctrine reusable across substrates; the kernel does not anticipate this substrate's tooling beyond the labelled realization peers. |

### §2.4 Honest-claim verification (rule 3.13)

- **(a) Reproducibility.** α's §ACs cites four executable oracles with captured stdout; β re-executed three (AC1 `test -f` → exists; AC7 awk + rg → 275 lines / 4 headers / no substrate matches; AC9 `git diff --stat` → 5 files, all named). All three reproduce α's reported values. PASS.
- **(b) Source-of-truth alignment.** Terms used in the doc (closed cell, receipt, verdict, decision, evidence, validator predicate `V`, repair-dispatch, scope-lift, projection-under-scope-lift) align with `COHERENCE-CELL.md` (receipt rule, four-way separation), `RECEIPT-VALIDATION.md` (verdict vs decision, δ-authoritative validation), and `#369`'s in-flight verdict × action × transmissibility framing. No drift from canonical predecessor terms detected. PASS.
- **(c) Wiring claims.** The doc claims four realization peers cite the kernel (`RECEIPT-VALIDATION.md`, `schemas/cdd/`, `cn-cdd-verify`, `CDD.md`). β grep-confirmed the four are real files / planned phases — `RECEIPT-VALIDATION.md` exists at the named path; `CDD.md` exists; `schemas/cdd/` is in flight per #369; `cn-cdd-verify` exists. Direction-of-dependency rule (realization → kernel) is asserted; reverse-citation (kernel → realization tooling) is absent from kernel sections per AC7. PASS.
- **(d) Gap claims.** γ-scaffold §Gap enumerates four predecessor surfaces and asserts the algorithm is "implicit across these four surfaces but stated nowhere." α's empirical anchor in γ-scaffold lines 25–43 grep-verified this with three commands (file absence + `normal form` keyword absence in package skills + predecessor surface presence). β re-ran `rg -i 'normal form|coherence-cell normal|CCNF' src/packages/cnos.cdd/skills/cdd/` (excluding the new file): no prior matches confirmed. PASS.

### Findings

| # | Finding | Evidence | Severity | Type |
|---|---|---|---|---|
| F1 | `self-coherence.md` uses `## CDD-Trace` (hyphen) as the section header; `cn-cdd-verify` requires `## CDD Trace` (space). The Build CI workflow fails on review SHA `aa10f902` as a direct consequence (`Build` run `25990518085`, job `CDD artifact ledger validation (I6)`, step `Check CDD artifacts`). The validator's grep at `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify:495` is literal `^## CDD Trace`. Convention precedent: `#367`'s self-coherence (released at `.cdd/releases/docs/2026-05-15/367/self-coherence.md` on `origin/main`) uses `## CDD Trace`. The manifest comment at the top of `self-coherence.md` (line 1) also enumerates `CDD-Trace` and should be aligned. | β re-ran `cn-cdd-verify --unreleased` locally; output: `❌ self-coherence.md sections — missing required sections: CDD Trace` followed by `❌ Cycle artifact verification FAILED`. CI run `25990518085` shows the same failure on the workflow surface. `git show origin/main:.cdd/releases/docs/2026-05-15/367/self-coherence.md \| grep '^## '` shows `## CDD Trace` (the convention). | B | mechanical (grep-detectable convention drift); classification `ci-status` per rule 3.10. |

### Required fix for F1

In `.cdd/unreleased/370/self-coherence.md`:

1. Replace the section header `## CDD-Trace` with `## CDD Trace` (one occurrence, around line 242).
2. Update the manifest comment at line 1 from `<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->` to `<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->` and the matching `<!-- completed: ... -->` line for internal consistency.
3. Commit on `cycle/370` with `git commit -m "α-370: rename §CDD-Trace → §CDD Trace (align with cn-cdd-verify convention)"`.
4. Push to `origin/cycle/370`.

After R2 push, β will re-fetch `origin/main`, recompute the diff base, re-run `cn-cdd-verify --unreleased`, and verify Build CI passes on the new review SHA before re-emitting verdict. No other change is requested.

### Regressions Required (D-level only)

None this round.

### Observations (not findings)

These are recorded for transparency; they do not block merge and α need not address them in R2.

1. **AC3 signature variant (γ.close and V arities).** The issue body's AC3 oracle lists `receiptₙ := γₙ.close(contractₙ, matterₙ, reviewₙ)` (3 args) and `verdictₙ := V(contractₙ, receiptₙ, evidenceₙ)` (3 args). The doc commits to `γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)` (4 args) and `V(contractₙ, receiptₙ)` (2 args). The doc justifies the variant at lines 110–116 ("V's signature shows two inputs, not three. The receipt carries the evidence references; V does not take an unbound `evidenceₙ` argument. This is the canonical form of evidence-binding…"). The variant strengthens the evidence-binding rule by making evidence a typed input to γ (which binds it) and removing the redundant evidence input to V (which dereferences from the receipt). The dispatch-prompt's three binding requirements for AC3 — β excludes evidence, δ excludes evidence, γ.close-and-V-not-collapsed — are all satisfied. AC3's own Positive case ("Five steps present in order; evidence-binding rule stated; β's signature explicitly excludes evidence; δ's signature explicitly excludes evidence") is satisfied. AC3's Negative case lists no failure mode the variant triggers. β accepts the variant as a coherent strengthening; α need not revisit unless γ requests realignment at PRA. *Recorded so the divergence from the literal issue-body signatures is auditable.*

2. **Doc length: 435 vs 200–400 target.** Doc lands at 435 lines (9% over upper target). Kernel slice (the four kernel sections) is 275 lines — well within target. Non-kernel sections (preamble + two-layer separation + non-goals + closure) account for 160 lines. β assessment: the non-kernel sections carry AC2 (predecessor citations), AC8 (realization-peer enumeration with framing), AC9 (non-goals enumeration), and the phase-inheritance table required by the cycle's Phase 7 spine purpose. None of the non-kernel content is operational drift; it is scaffolding around the kernel. Per the β dispatch prompt, drift past 400 lines is a process signal, not a binding finding. α's §Debt #1 self-noted this overage. Recorded as observation only.

3. **Substrate-independence held under stress.** α's §Self-check notes that the AC7 kernel-slice scan caught a single substrate-term leak during authoring — a `[#369](https://github.com/...)` URL in §Cell Outcomes — which α corrected to a bare `#369` reference before review-readiness. The mechanism worked as designed and α's authoring discipline executed against γ-scaffold's failure-mode #1 (substrate leakage). β's re-execution of the AC7 oracle returns no matches.

### Notes

This cycle ships a doctrine surface only; the cycle's surface containment is tight. The single finding is a mechanical convention drift (`CDD-Trace` → `CDD Trace`) detected by the CI validator. R2 should be a one-character (and one whitespace) edit on α's self-coherence file plus the manifest comment alignment, after which the Build CI workflow should turn green and β can approve.

The signature variant in AC3 is recorded as an observation rather than a finding because (a) it is internally consistent and strengthens kernel coherence, (b) all binding AC3 properties from the dispatch prompt and the issue body's Positive/Negative oracle hold, and (c) the doc is explicit about the variant and its justification.

α may proceed with the F1 fix in R2 without re-running the full AC bank — only the validator-relevant rows and AC7 (to confirm the section-header rename doesn't break the awk pattern, which it cannot since the rename is in a different file). β will re-run AC1, AC7, AC9 oracles and the validator on the new review SHA.
