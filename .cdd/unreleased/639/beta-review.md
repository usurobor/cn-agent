# beta-review.md — cnos#639 (β)

## §R0

**Verdict:** CONVERGE (APPROVED — all issue ACs met, zero unresolved findings)

**Round:** 0
**Base SHA:** `f2959e150c0a308db460afe4fcec4b6c4429ee34` (`origin/main`)
**Head SHA reviewed:** `649c8649c2db6d0cc87e78768e9e912151e7cb14` (`origin/cycle/639`)
**Branch CI state:** green — `Build` `success` at HEAD (`649c8649`, run 28999353785, all 10 jobs green: Go build & test, Binary verification, Package verification, Package/source drift I1, Protocol contract schema sync I2, Repo link validation I4, SKILL.md frontmatter validation I5, CDD artifact ledger validation I6, Dispatch repair-preflight guard cnos#516, Dispatch closeout-integrity guard cnos#524); `install-wake golden` `success` at `beccf0987` (the last commit that touched `cds-dispatch/SKILL.md`; HEAD's only delta over `beccf0987` is a prose-only append to `self-coherence.md`, confirmed via `git diff beccf0987...649c8649 --stat`, so the golden invariant established at `beccf0987` still holds at HEAD — independently re-verified below, not merely assumed from the older SHA).
**Merge instruction:** `git merge cycle/639` into `main` with `Closes #639` (not yet executed — this round is the review verdict only, per δ's dispatch scope for this pass).

This review independently re-derived every AC oracle against the actual diff and cited artifacts — it does not take α's `self-coherence.md` framing on trust. All citations below were re-read from source, not copied from α's paraphrase.

---

## §Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Issue #639 `status:in-progress`, claimed cleanly by `cds-dispatch` at run 28998099166; no status/label contradiction observed. |
| Canonical sources/paths verified | yes | All three cited doctrine surfaces (`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md`) exist at the paths cited and the cited line ranges match. |
| Scope/non-goals consistent | yes | Diff is doctrine-only; verified below (§Findings / diff footprint). |
| Constraint strata consistent | yes | Implementation contract (doctrine/Markdown only, no code/FSM/label) matches the actual diff shape exactly. |
| Exceptions field-specific/reasoned | n/a | No exceptions claimed. |
| Path resolution base explicit | yes | All cross-references use relative paths that resolve from their source file's location; spot-checked `cds-dispatch/SKILL.md`'s links to `delta/SKILL.md` §9.11/§9.13 and vice versa. |
| Proof shape adequate | yes | AC1's mapping table cites artifact path + line range per shape; AC4's grep + CI-guard re-run is reproducible; AC5's diff-stat claim is reproducible. |
| Cross-surface projections updated | yes | `dispatch-protocol/SKILL.md` §2.8/§4.8, `cds-dispatch/SKILL.md` Step A + enum, `delta/SKILL.md` §9.13 all updated consistently; no surface left stating a stale/narrower enum as complete. |
| No witness theater / false closure | yes | α's `self-coherence.md` §Debt discloses I4/I5/I6 as "not independently runnable in this session" rather than claiming false-green; this review confirms the real CI actually ran all three green (see §Findings). |
| PR body matches branch files | n/a | No PR is open yet for `cycle/639` (`gh pr list --head cycle/639` → empty) — β reviewed the branch diff directly (`origin/main...HEAD`), which is the equivalent evidence surface; not a blocker for this review pass. |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/639/gamma-scaffold.md` present on `origin/cycle/639` (rule 3.11b satisfied). |

---

## §Issue Contract — AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | One canonical, exhaustive, mutually-exclusive `run_class` set; every observed shape maps to exactly one class | yes | **Met** | See §Findings AC1 below — independently re-walked all five shapes and the mutual-exclusivity construction. |
| AC2 | Scope-continuation shape named with defined re-entry behavior | yes | **Met** | `cds-dispatch/SKILL.md` enum + Step A check 2; `delta/SKILL.md` §9.13 (new, additive, does not touch §9.10/§9.11/§9.12). |
| AC3 | Recovery-vs-resume and reset-branch-vs-first-pass resolved explicitly | yes | **Met** | Explicit paragraphs at Step A check 1 (recovery/resume = one class) and check 2 (issue-level history gates continuation, cites #626 R2); procedure order re-verified, not just prose. |
| AC4 | Divergent enumerations reconciled; no surface carries a divergent list | yes | **Met** | `dispatch-protocol/SKILL.md` L445 now states the full 6-value set; `cds-dispatch/SKILL.md` is the single canonical source; `delta/SKILL.md` narrates shapes only, no competing bracketed enum. CI guard re-run green. |
| AC5 | Doctrine-only; no code/FSM/label/dispatch-behavior change; gates green | yes | **Met** | Zero `.go` files, zero `transitions.json` diff, zero label changes in diff; golden/live workflow byte-identical to a fresh `cn install-wake` render; all 10 Build-workflow jobs green at HEAD. |

## §Issue Contract — Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | yes | updated | Step A rewritten to a 4-check ordered procedure + mutual-exclusivity table + 5-shape mapping table + enum extended to 6 values. |
| `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | yes | updated | §2.8 "Detection" + "Receipt" reconciled to point at the canonical enum rather than restate a narrower one; §4.8 verify-language re-checked, still accurate post-edit. |
| `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | yes | updated | New §9.13, additive at EOF; §9.10/§9.11/§9.12 untouched (confirmed via diff hunk boundaries — only one hunk, pure append). |
| `cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml` | yes | regenerated | Verified below to be byte-identical to a fresh `cn install-wake cds-dispatch` render — not hand-edited. |
| `CELL-KINDS.md` / `docs/development/cdd/CDD.md` | no | confirmed not required | α's §Debt item 4 confirms γ's read (`run_class` is a δ/dispatch-runtime concern, distinct from `cell_kind`'s `FactSnapshot.CellKind` seam) and re-greps both files for `run_class` (0 hits) before declining to add a cross-reference. β independently re-ran the same grep and confirms 0 hits in both files — no cross-reference needed, decision honestly disclosed rather than silently skipped. |

## §Issue Contract — CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `gamma-scaffold.md` | yes | yes | Present, per-AC oracle list + source-of-truth table + friction notes all substantive. |
| `self-coherence.md` | yes | yes | Contains `§Gap, §Skills, §ACs, §Self-check, §Debt, §CDD Trace, §Review-readiness` — canonical section headers present verbatim. |
| `CLAIM-REQUEST.yml` | yes (claim marker) | yes | Pre-existing per γ's Friction note 7; `run_class: first_pass` for this cell's own dispatch (correctly self-referential — the cell about the taxonomy is itself an ordinary first pass). |
| `beta-review.md` | this round | yes | This file. |

## §Issue Contract — Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `beta/SKILL.md` | β role contract | yes | yes | Load order followed: `CDD.md` (canonical lifecycle, referenced), `beta/SKILL.md`, `review/SKILL.md`. |
| `review/SKILL.md` | review orchestration | yes | yes | Contract-integrity → issue-contract → verdict phases followed; output format matches. |
| `cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md` | issue-specific Tier 3 | yes | yes | All three read in full (not excerpted) — the doctrine surfaces under review. |

---

## §Diff Context / Architecture

**Diff footprint (AC5, independently re-derived).** `git diff --stat origin/main...HEAD`: 8 files — three doctrine `.md` files (`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md`), two mechanically-regenerated non-`.md` files (`cnos-cds-dispatch.golden.yml`, `.github/workflows/cnos-cds-dispatch.yml`), and three cell-artifact `.md`/`.yml` files under `.cdd/unreleased/639/` (`CLAIM-REQUEST.yml`, `gamma-scaffold.md`, `self-coherence.md`). `git diff origin/main...HEAD -- '*.go'` → empty (0 lines). `git diff origin/main...HEAD -- '**/transitions.json'` → empty (0 lines; the file does not appear in the diff at all). Grepped the full diff for `status:`/`dispatch:`/`protocol:`/`label` occurrences — every hit is inside prose (doctrine text describing existing label semantics, or `.cdd/unreleased/639/` artifact narrative); zero label *definitions* added, removed, or redefined.

**`cn install-wake` regeneration verified independently, not assumed.** Ran the actual renderer (`CN_PACKAGE_ROOT=src/packages/cnos.cds sh src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --out /tmp/regen-live.yml`) against the reviewed cycle branch's `cds-dispatch/SKILL.md`. `diff /tmp/regen-live.yml .github/workflows/cnos-cds-dispatch.yml` → empty (byte-identical). `diff /tmp/regen-live.yml cnos-cds-dispatch.golden.yml` → empty (byte-identical). Both committed rendered artifacts are confirmed pure tool output, not hand-edited.

**CI guards re-run locally, both green.** `./scripts/ci/check-dispatch-repair-preflight.sh` → exit 0. `./scripts/ci/check-dispatch-closeout-integrity.sh` → exit 0. Both also ran and passed as jobs inside the `Build` workflow at HEAD (`gh run view 28999353785 --json jobs` confirms `Dispatch repair-preflight guard (cnos#516): success`, `Dispatch closeout-integrity guard (cnos#524): success`).

**Go modules build/vet clean in all five `go.work` members** (`./src/go`, `cnos.cdd/commands/cdd-verify`, `cnos.issues/commands/issues-map`, `cnos.issues/commands/issues-fsm`, `cnos.issues/commands/issues-dispatch`) — expected no-op since zero `.go` files changed, but independently confirmed rather than assumed.

**#642 confirmed not dispatched.** `gh issue view 642` shows 0 comments, still open, no in-progress/review status label — untouched.

**#626 / sparse-checkout / write-fence content untouched.** `delta/SKILL.md`'s diff is a single append-only hunk at EOF (line 622 onward); §9.12 (the cell/substrate-identity-boundary content) is above that hunk boundary and shows zero changes. `cds-dispatch/SKILL.md`'s "Disallowed surfaces" section (sparse-checkout / write-fence narrative) is likewise outside both of that file's two diff hunks (152–204, 210–252).

---

## §AC1 — independent re-derivation (mutual exclusivity + five-shape mapping)

Re-read `cds-dispatch/SKILL.md` Step A directly (not the self-coherence.md paraphrase) and re-derived each shape's classification from the cited source artifacts:

- **#593** (`.cdd/unreleased/593/gamma-scaffold.md` — read in full): "A prior `cds-dispatch` firing had written `.cdd/unreleased/593/CLAIM-REQUEST.yml` and created `cycle/593`, but the FSM claim transition never applied... No `status:changes` history, no beta/delta findings, no closeouts." → check 1 fails (no scanner comment), check 2 fails (no operator continuation comment, no prior merged round), check 3 fails (no rejection evidence) → `first_pass`. Matches the diff's table row 1.
- **#630** (`.cdd/unreleased/630/gamma-scaffold.md` — read in full): "no prior `status:changes` history..., no prior `cycle/630` branch..., no prior `.cdd/unreleased/630/` artifacts..., zero issue comments beyond the operator's dispatch-authorization comment." → all three checks fail trivially → `first_pass`. Matches row 2.
- **#614** (`.cdd/unreleased/614/alpha-closeout.md` — read in full): explicit self-disclosure quoted in the source: "the operator's 2026-07-08 'Dispatch authorized' comment explicitly frames this as 'a clean resume,' not a repair." No scanner comment, no continuation-authorization comment framing a *new bounded scope on merged matter* (the issue had no prior merged round at all), no rejection evidence → all three checks correctly fail despite branch/artifact matter pre-existing → `first_pass` with disclosure. Matches row 3, and is the sharpest test of the new doctrine's tightened `repair_pass` discriminator (old pre-#639 doctrine's "branch exists with commits" bullet would have risked misrouting this to `repair_pass`; new doctrine correctly does not).
- **#626 R0** (`self-coherence.md` §R0 + `gamma-closeout.md`): "run_class: first_pass" verbatim, clean claim → matches row 4.
- **#626 R1/R2** (`self-coherence.md` L432–448, L647–649; `gamma-closeout.md` L309–315, L382–385 — read in full): confirmed the "run_class note" sections state explicitly: no scanner comment, no rejection evidence, but an operator/CAP comment authorizing a bounded new AC scope (R1: AC3/AC4; R2: AC4) on an issue whose R0 (PR #635) and R1 (PR #637) had already merged before the respective claims, with R2's branch reset to a clean `main` base — check 2 passes on exactly this evidence (operator comment + issue-level merged-round evidence, not branch-commit-count) → `scope_continuation`. Matches row 5, and is the concrete empirical case AC3(b) requires be resolved by issue-level history rather than branch state.

**Mutual exclusivity re-verified structurally, not just by the stated pairwise table.** Each non-terminal check's own definition embeds the negative conditions that guarantee no overlap: check 2 (`scope_continuation`) explicitly requires check 1's positive trigger ("MECHANICAL reversion" comment) to be *absent*; check 3 (`repair_pass`) is reachable only when both check 1 and check 2 have already failed, and both of those checks require rejection evidence to be absent as one of their AND-conditions — so if rejection evidence is present, checks 1 and 2 necessarily fail regardless of what else is present, and check 3 is the only one that can fire. This means mutual exclusivity is enforced by the discriminator definitions themselves, not merely by evaluation order — a stronger property than the oracle asked for. No adversarial construction (e.g., an issue with both a scanner comment *and* an operator continuation comment *and* rejection evidence all simultaneously present) produces an ambiguous double-match: the ordered AND-gated conditions collapse it to exactly one class every time (rejection evidence wins over both matter-shapes; scanner comment wins over continuation when no rejection evidence is present).

**Negative-check re-verified:** `grep -c "Positive discriminator:" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` → 3 (checks 1–3), plus `manual_delta_repair`/`blocked` each carry a named positive event in the Step-A note (repair commits landed out-of-band; a posted STOP/BLOCKED comment) rather than being defined by elimination. Only `first_pass` (check 4) is exclusion-defined, exactly as the oracle requires.

## §AC2/AC3 — independent re-verification

`grep -n "does NOT re-scaffold\|never treated as a conflict"` against `delta/SKILL.md` hits §9.13's routing step 1 verbatim: "δ does NOT dispatch γ for a fresh R0 scaffold... prior `CDDArtifacts` are read as valid closed history, never treated as a conflict requiring resolution" — structurally parallel to §9.11's own "does NOT overwrite it and does NOT dispatch γ for a new R0 scaffold" phrasing, confirmed by direct comparison of both sections' text.

Re-read the "Recovery vs. resume" paragraph (Step A check 1) and the "reset-branch" paragraph (Step A check 2) directly, not via self-coherence.md's summary. Both are present as stated, and the ordered procedure — re-read top to bottom, not merely trusted from prose elsewhere — genuinely places check 2 (issue-level, `scope_continuation`) before check 4 (`first_pass`, terminal fallthrough). There is no branch-state-only check positioned ahead of it.

## §AC4 — independent re-verification

`grep -n "run_class" src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` re-run: `dispatch-protocol/SKILL.md` L445 states the full 6-value set with an explicit pointer to `cds-dispatch/SKILL.md` as the canonical source (not a competing restatement); `cds-dispatch/SKILL.md` is the canonical home (both the ordered procedure and the flat enum, self-consistent — 6 values match in both places); `delta/SKILL.md` narrates shapes without a bracketed enum anywhere (spot-checked §9.10/§9.11/§9.13 — none restate `run_class ∈ {...}`).

`grep -rln "run_class" --include="*.md" src/packages/` (repo-wide, not limited to the three named files) → exactly the same three files. No fourth active doctrine surface under `src/packages/` mentions `run_class`. (The one other repo-wide `.md` hit, `docs/evidence/rca/2026-06-30-cnos524-w4-empty-review.md`, states its own cell's *value* — `manual_delta_repair` — as a historical record, not an enum restatement; not a divergence.)

CI guard hard constraint re-verified: all literal substrings the CI script `grep -qF`s for (`run_class`, `first_pass`, `repair_pass`, `manual_delta_repair`, `blocked` in `dispatch-protocol/SKILL.md`; `run_class`, `first_pass`, `repair_pass` in `cds-dispatch/SKILL.md` + golden + live workflow) are present verbatim — only `resumed_from_matter` and `scope_continuation` were added.

## §AC5 — independent re-verification

Covered above under §Diff Context / Architecture and the CI job table in §Contract Integrity's header. All 10 named gates (Go build & test, Binary verification, Package verification, I1, I2, I4, I5, I6, dispatch-repair-preflight, dispatch-closeout-integrity) ran and passed inside the `Build` workflow at HEAD `649c8649` — this closes the gap α honestly disclosed as "not independently runnable in this session" for I4/I5/I6: real CI ran and passed all three.

---

## §Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|

No findings at any severity. Every AC oracle re-derived independently from source artifacts (not from α's self-coherence.md framing) confirms the claimed result. Scope guardrails (no `.go`, no `transitions.json` diff, no label changes, no #626/sparse-checkout/write-fence touch, #642 not dispatched, `scripts/ci/*.sh` untouched, golden/workflow mechanically regenerated and byte-verified against a fresh render) all hold. Implementation contract (Rule 7 — Markdown/doctrine-only, no CLI/package/binary/runtime/wire-contract/backward-compat axis violated) confirmed conformant.

## §Regressions Required (D-level only)

None.

## §Notes

- No PR is currently open for `cycle/639`; this review evaluated the branch diff (`origin/main...HEAD`) directly, which carries the same evidentiary weight for this review pass. δ's finalizer step (or the merge itself) will open/associate the PR per the cell's normal flow.
- Close-keyword presence (β pre-merge gate row 5) and the merge itself are out of scope for this review pass per the dispatching δ's explicit instruction for this round ("write beta-review.md... do NOT spawn α or γ... report back to δ"); β's pre-merge gate will be re-run in full at the actual merge pass.
- Search space closed: no remaining blocker found in any of AC1–AC5, the scope guardrails, the implementation contract, or the CI-guard hard constraint.
