# α Close-Out — Cycle #268 (PR #270)

## Cycle summary

| | |
|---|---|
| Issue | #268 — *Converge cnos.cdd skill-program contracts from CDD package audit*. Filed by γ from `docs/gamma/cdd/CDD-PACKAGE-AUDIT.md` D/C findings. Closes on merge of PR #270 (`Closes #268`). |
| PR | #270 — *cdd: converge skill-program contracts from package audit (#268)*. Squash-merged to main as `bfddcf22`. |
| Dispatch | γ → α: project `cnos`, issue #268, Tier 3 skills `cnos.core/skills/skill` (skill-program/frontmatter coherence), `eng/tool` (cn-cdd-verify), `eng/document` (CDD/skill prose), `eng/test` (verifier tests). Branch `claude/fix-skill-frontmatter-coherence-qdzch`. |
| Work shape | Substantial documentation/contract convergence (MCA, L7-shaped). 14 files in PR #270 (13 modified + 1 new test fixture); +518 / −242 net; 2 commits on-branch (`6c77aee` initial, `15e8366` round-1 fix). One file deletion (`review/checklist.md`, the orphan checklist that contradicted release/closure cleanup). |
| Tier 2 skills applied | `eng/document` (durable docs/skill prose), `eng/tool` (cn-cdd-verify rewrite), `eng/test` (fixture-based test script + 17 assertions). |
| Tier 3 skills applied | `cnos.core/skills/skill` (frontmatter coherence — visibility, triggers, classification, authority hierarchy across `cdd/SKILL.md`, `design/SKILL.md`, the eng/README path drift). |
| MCA selection | The audit named ~21 findings; the issue distilled them into 16 ACs with five operator decisions (D1–D5). MCA scope was therefore largely pre-decided by γ; α's job was to land the named patches coherently and synchronously across CDD.md, six role/lifecycle skills, eng/README.md, and the verifier. The L7 leverage move is the new §5.3a Artifact Location Matrix in CDD.md plus its mechanical enforcement in `cn-cdd-verify` — turning "where does artifact X live?" from a contract-drift problem into a verifier-failure if violated. |
| Review rounds | **2** — R1 REQUEST CHANGES + 1 finding (F1, C-mechanical, stale backtick to deleted `review/checklist.md` in `cnos.core/mindsets/ENGINEERING.md:5`); R2 APPROVED final after 1-line repoint commit. β merged at `bfddcf22`. |
| Findings against α | 1 total — F1 (C, mechanical, deletion left one live cross-reference). Resolved on-branch per `review/SKILL.md` §7.0. No code or semantic findings. |
| Mechanical ratio | 1/1 = 100%. N=1 well below §9.1 10-finding floor; not an automatic process-debt trigger. |
| Local verification | `bash src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh` → 17/17 assertions pass on every review-round head (initial commit, F1 fix). Six fixture cases cover canonical accept, v-tag warn, legacy-PRA-only fail, missing-bare-tag fail, triadic canonical accept, and triadic missing-γ fail. CI on merged head `bfddcf22`: 7/7 green (go, kata-tier1, kata-tier2, I1 package/source-drift, I2 protocol-contract-schema-sync, 2× notify). |
| Concurrent main activity | None during review (review opened and closed within ~minutes; CI completed before β posted). The 3.58.0 PRA on main (`a4f1b7ea`) names #268 as the next MCA, so this cycle is the planned implementation of that commitment. |

## Findings (α-side observations)

### Finding F1 — stale backtick path after deleting `review/checklist.md` (C, mechanical)

**What happened.** AC9 directed deletion of `src/packages/cnos.cdd/skills/cdd/review/checklist.md` (the orphan/contradicted checklist). The deletion landed in `6c77aee`. β's R1 grep across the merged tree found exactly one live consumer of the path:

```
src/packages/cnos.core/mindsets/ENGINEERING.md:5:**Verify `skills/review/checklist.md` P0 + P1 before pushing.**
```

The path's only resolvable target was the file deleted in this PR. ENGINEERING.md is shipped mindset doctrine, referenced live by `docs/alpha/cli/SETUP-INSTALLER.md`, `docs/alpha/cognitive-substrate/CAR.md`, `docs/alpha/package-system/PACKAGE-ARTIFACTS.md`, `docs/beta/guides/WRITE-A-SKILL.md`, `docs/gamma/checklists/engineering.md`, and `src/ocaml/cmd/cn_assets.ml` — a live consumer, not vestigial.

**Pattern.** `review/SKILL.md` §2.2.5 names the rule for **file moves**: "If files move or are renamed, grep for old paths across all live docs. Any live stale path is a D-level blocker." The rule applies symmetrically to file **deletions** — a deletion is a path-going-from-live-to-nothing, structurally identical to a path-changing-from-A-to-B for any consumer that referenced the deleted name. I read the skill text literally as "moves" and did not generalize to "deletions."

**Resolution chosen.** β suggested two paths: (a) repoint at the surviving authority (`cdd/review/SKILL.md` §4 Checklist + §7 After Review) or (b) delete the line. I took (a) — the engineer-facing "before pushing" mindset cue is still useful, and §4 of the surviving review skill already supersedes the deleted file's content (P0/P1 governance + process items, plus the unicode hygiene and CI-gate items the deleted checklist lacked). Repoint commit `15e8366` reads:

```
**Verify `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` §4 (Checklist) and §7 (After Review) before pushing.**
```

Post-fix `git grep "skills/review/checklist\|review/checklist" -- ':!docs/gamma/cdd/CDD-PACKAGE-AUDIT.md'` → zero hits.

**Surfaces involved.** `review/SKILL.md` §2.2.5 (cross-reference validation, currently scoped to file *moves*); `src/packages/cnos.core/mindsets/ENGINEERING.md`; the deletion in PR #270.

## Cycle Iteration (CDD §9.1)

### Triggers fired

- [ ] review rounds > 2 — actual: 2. Below the threshold.
- [ ] mechanical ratio > 20% with ≥ 10 findings — actual: 100% (1/1) but N=1 ≪ 10-finding floor; below the §9.1 automatic-trigger threshold.
- [ ] avoidable tooling / environmental failure — none. The verifier-test fixture initially failed once on GPG signing in throwaway repos; I added `commit.gpgsign=false` / `tag.gpgsign=false` to the fixture builder before requesting review, so β saw a passing test run. Not "avoidable failure during the cycle"; just a one-off in fixture authoring caught at write-time.
- [x] **loaded skill failed to prevent a finding** — `review/SKILL.md` §2.2.5 was nominally loaded as a Tier 1 lifecycle constraint via the role's load order. Its principle covers F1 by generalization (any live stale path after a path-changing edit is a finding), but its examples and stated trigger phrase ("If files move or are renamed") biased my application toward moves only. I deleted a file and did not run the §2.2.5 grep against the deleted path.

### Friction log (α-side)

1. **Deletion-as-path-change scope mismatch.** `review/SKILL.md` §2.2.5 names "moves and renames" but applies equally to deletions for any consumer holding the old path. My pre-review pass enumerated the `Artifact Location Matrix` schema-bearing change correctly (alpha §2.4 harness audit), but treated the `review/checklist.md` deletion as a closed-loop AC9 item — checklist-deletion has no consumers in the same skill program, which is true within the package, false across `src/packages/cnos.core/mindsets/`. The pre-review-gate row for "schema/shape audit" was satisfied for the matrix change; no row was satisfied for the deletion.

2. **Cross-package live-doc reach not enumerated.** `src/packages/cnos.core/mindsets/ENGINEERING.md` is outside the issue's declared scope, but it is a live consumer of paths declared by the CDD package's skill program. Per `eng/go §2.12` ("sibling harnesses are contract surfaces too"), live docs in adjacent packages are also harnesses for path contracts in the package being changed. β's grep over the full repo caught what α's grep over the package didn't.

### Root cause (α-side reading)

**Peer-enumeration scope mismatch on file deletion.** I enumerated peers for the *additive* contract change (Artifact Location Matrix) — every consumer that needed to reference §5.3a was updated (alpha, beta, gamma, release, post-release, cn-cdd-verify, plus the test fixture). I did not enumerate peers for the *subtractive* contract change (`review/checklist.md` deletion) — every consumer that referenced the now-deleted path. The same `peer enumeration before closure claims` rule (`alpha/SKILL.md` §2.3) covers both directions; I applied it asymmetrically.

**Adjacent failure: `review/SKILL.md` §2.2.5 phrasing.** The skill says "If files move or are renamed, grep for old paths." The principle is path-going-stale, regardless of mechanism (move, rename, delete). The phrasing biased my reading toward moves and renames. This is the kind of cross-scale generalization the surrounding text already does for `eng/go §2.12` (intra-doc → cross-producer); §2.2.5 is the same shape and benefits from the same generalization, but the existing text only names two of the three trigger conditions.

### Engineering level (α-side reading, for γ to adjudicate)

- **L5 (local correctness):** met on the diff itself. Verifier tests 17/17 green on every review-round head. CI 7/7 green on the merged head. No code semantic finding. The verifier rewrite enforces canonical paths; the test fixture demonstrates accept/reject. Diff-level mechanical scans (no duplicate list entries; no hidden/bidi unicode; no whitespace anomalies) clean.

- **L6 (system-safe execution):** partial miss. F1 reached β review — a live cross-surface stale path. The miss is real even though containment was clean (1 round, 1 line, no semantic ripple). L6's criterion is "cross-surface drift did not reach review"; this reached review. Caught at R1, fixed in `15e8366`, re-verified by R2.

- **L7 (system-shaping leverage):** achieved on the substantive axis. The Artifact Location Matrix + `cn-cdd-verify` enforcement together turn "release looks valid but is noncanonical" from a class of recurring drift into a verifier-failure-or-canonical-pass. Future cycles cannot land a release with an assessment in `.cdd/releases/.../ASSESSMENT.md` while leaving the canonical `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` empty — the verifier fails fast. Same for `v`-prefixed tags (warn-only) and triadic close-out paths.

Per §9.1 "lowest miss" rule, these combine to **L6 (with an L7 MCA shipped, capped by the L6 deletion-peer-enumeration miss)**. β's close-out reads the same axes independently.

## What worked

1. **Audit-driven AC structure made closure mechanical.** γ distilled a ~21-finding audit into 16 ACs with five operator decisions (D1–D5) pinned in the issue body. Each AC was a single concrete edit with named files. My self-coherence section in the PR body was a 16-row AC → evidence map; β's R1 review used the same row structure and verified each line. No interpretation gap between issue, PR body, and review.

2. **Test-alongside-implementation for `cn-cdd-verify`.** I wrote `test-cn-cdd-verify.sh` while rewriting the verifier, with six fixture-based cases covering canonical pass, v-tag warn, legacy PRA path warn, missing bare tag fail, triadic canonical pass, missing-γ-close-out fail. The first run surfaced two issues immediately: GPG signing failures in throwaway repos (fixed in fixture builder), and a too-permissive triadic close-out fallback (originally allowed PR-comment α/β as pass; AC15's "required canonical α/β/γ close-out paths" mandates a stricter read). The latter caught a literal AC misinterpretation before β did. Post-patch re-audit per `alpha/SKILL.md` §3.4 surfaced this; tests still 17/17 after the tightening.

3. **Schema-bearing change audit for the matrix worked.** Adding §5.3a Artifact Location Matrix is a schema-bearing change (per `alpha/SKILL.md` §2.4): every consumer of "where does X live?" is now bound to the matrix. I enumerated consumers by name — `alpha/SKILL.md` §2.8 (close-out path), `beta/SKILL.md` §5 (close-out path) and load-order header, `gamma/SKILL.md` §2.7 + load-order, `release/SKILL.md` header, `post-release/SKILL.md` header + Output, `cn-cdd-verify` header — and added matrix-pointer references to each. β's review confirmed each pointer existed by grep.

4. **Section-by-section authoring for ≥50-line edits.** CDD.md §4.4 rewrite, §5.3a matrix, gamma §2.7 expansion, alpha §2.2 reorder, post-release Output rewrite, the verifier rewrite — each over ~50 lines and each written section-by-section to disk per `CDD.md` §1.4 large-file rule. No partial-recovery problem; no compaction-loss anxiety.

5. **β's narrow R1 + suggested fix in two forms.** β filed F1 once, gave both the failing grep and two patch options (repoint vs delete) with rationale for each. The fix was 1 line, deterministic, and self-evident given the suggested patch text. R2 closed in one round.

## What didn't

1. **Asymmetric peer enumeration on AC9 (deletion).** I peer-enumerated the additive matrix change exhaustively. I did not peer-enumerate the subtractive deletion at all — I treated `review/checklist.md` as having "no consumers in the package" (true) without grep-ing the rest of the repo. F1 is the direct consequence. `alpha/SKILL.md` §2.3 (peer enumeration before closure claims) covers both directions of contract change; my application was direction-asymmetric.

2. **`review/SKILL.md` §2.2.5 read literally as "moves" only.** The skill names file moves and renames; deletions are the same class for any consumer holding the old path. My read was within the literal scope. The skill text would catch more deletion cases if it named the trigger as "any path-going-stale event" rather than "moves and renames." (Pattern observation only — patching that skill is γ's call, and would be a separate cycle's MCA if the pattern recurs.)

3. **Hub memory not re-checked at cycle close.** Per `CDD.md` §1.4 γ algorithm Phase 5, hub memory writes are γ's responsibility. α does not own hub memory. Noting because my last cycle-touch is this close-out commit; γ closes the loop.

## Observations and patterns (α-side)

### 1. Deletion-as-path-change in `review/SKILL.md` §2.2.5

**Pattern.** §2.2.5 names "If files move or are renamed, grep for old paths across all live docs." A deletion is structurally a move-to-nothing — every consumer of the deleted path now holds a stale reference. The current text covers two of three trigger events (move, rename); deletion is the third and was the one this cycle hit.

**Surfaces this pattern touches.** `review/SKILL.md` §2.2.5 (trigger phrasing); `alpha/SKILL.md` §2.3 (peer enumeration on subtractive changes); `eng/go §2.12` (the same path-contract-as-harness principle at the larger scale).

**Evidence.** N=1 — this cycle. Pattern observation only; not a claim of generality.

### 2. Test-alongside-implementation as a closure-overclaim guard

**Pattern.** Writing the verifier tests during the verifier rewrite (not after) caught two issues at write-time that would otherwise have surfaced in review. The first (GPG fixture failure) was mechanical. The second (triadic close-out fallback too permissive against AC15's literal text) was a closure-overclaim risk: the original implementation read AC15 as "require γ canonical, allow α/β PR-comment fallback as pass," which is wrong against the AC's "required canonical α/β/γ close-out paths" wording. Re-running the test suite against the AC text caught the disagreement.

**Surfaces this pattern touches.** Not a skill-text gap — the discipline is already named in `alpha/SKILL.md` §3.4 (post-patch re-audit). The observation is that the discipline pays for itself in a cycle whose AC depends on a verifier behavior; the test was the audit surface.

**Evidence.** N=1 — this cycle.

### 3. Audit-pre-distilled-into-ACs reduces α scope-drift risk

**Pattern.** γ ran a separate audit cycle (`docs/gamma/cdd/CDD-PACKAGE-AUDIT.md`) with severity classification (D/C/B/A) and patch-level recommendations. The issue's 16 ACs are a direct map from the audit's D/C-class findings to landable patches. Every AC was a concrete edit with named files; I never had to interpret "what does this audit point mean as a patch."

**Observation.** When γ files an issue distilled from an upstream audit (rather than naming the audit + asking α to triage), α's interpretation surface shrinks to "land each AC's named patch coherently." The cost paid is γ's audit-distillation work; the benefit is α's faster convergence on closed scope.

**Surfaces this pattern touches.** `issue/SKILL.md` §2.4 (skills and constraints) and §2.5 (related artifacts) — the issue's "Related artifacts" section linked the audit; the body distilled it. Not an α-side skill change; noting because the hand-off is why this cycle's scope was unambiguous.

**Evidence.** N=1 — this cycle (#268 ← `CDD-PACKAGE-AUDIT.md`).

---

Signed: α (`alpha@cdd.cnos`) · 2026-04-25 · merge commit `bfddcf22` · cycle level per §9.1: L6 (L7 MCA shipped, capped by L6 deletion-peer-enumeration miss).
