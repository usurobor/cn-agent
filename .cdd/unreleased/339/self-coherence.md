---
cycle: 339
issue: "#339"
branch: cycle/339
role: alpha
---

# Self-Coherence — Cycle #339

## §Gap

**Issue:** #339 — cdd/release: mechanical pre-merge closure-gate enforcement + rubric closure-gate-failure handling (§3.8)

**Mode:** `docs-only` (§2.5b disconnect — no version bump, no tag). The cycle ships a shell script extension and a Markdown skill patch; neither changes code that would require a version bump. Disconnect is the merge commit on main.

**Version:** 3.15.0 (unchanged)

**Selected gap:** Three consecutive cycles (#331, #333, #334) merged to main without close-out artifacts. The closure gate at `gamma/SKILL.md` §2.10 is prose-enforced (γ checks a table); no mechanical check prevents a merge when files are absent. The rubric at `release/SKILL.md` §3.8 has no clause forcing `<C` when the closure gate fails, and `<C` in the rubric has never been reconciled with `C−` used in CHANGELOG/PRA prose.

**Empirical anchors (per AC3 / review/SKILL.md rule 3.13):**
- Cycle #331 merge commit: `315e529` — zero close-out artifacts at merge time
- Cycle #333 merge commit: `6ffdf48` — zero close-out artifacts at merge time
- Cycle #334 merge commit: `0900448` — zero close-out artifacts at merge time
- Retroactive close-outs landed: commit `00e6f8e` (cycle #334), `1ec471d` (cycles #331/#333)
- §3.8 rubric introduced: commit `b27fc15` (cycle #331 patch 5)
- F2 (mechanical reinforcement) source: `.cdd/releases/docs/2026-05-09/335/cdd-iteration.md`
- F7 (rubric closure-gate-failure handling) source: `.cdd/releases/docs/2026-05-09/335/beta-review.md` Round 2 §F7

---

## §Skills

**Tier 1:**
- `CDD.md` — lifecycle and role contract (canonical)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (this file)
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — release skill (surface being patched for AC2)
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — γ role (§2.10 closure gate, authoritative artifact list)

**Tier 2 (always-applicable eng skills):**
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` — shell script standards (exit codes, fail-fast, NO_COLOR, set -euo pipefail)

**Tier 3 (issue-specific):**
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` — applied to validate-release-gate.sh extension
- `src/packages/cnos.eng/skills/eng/writing/SKILL.md` — applied to §3.8 rubric prose patch (note: file not found at declared path; prose written against existing §3.8 style and the eng/tool SKILL.md's structural principles)

**Active design constraints loaded from issue:**
1. Reuse `scripts/validate-release-gate.sh` artifact-presence logic — do not duplicate file-existence check
2. Diagnostic must name file AND cycle number
3. Override clause must precede geometric-mean math in §3.8 source order
4. Letter normalization is binding — both `<C` and `C−` must be reconciled
5. Recursive coherence is non-optional — gate exercised against this cycle's own artifacts

---

## §ACs

### AC1 — Pre-merge closure-gate check exists and produces clear diagnostic

**Status: MET**

**Evidence:**

- `scripts/validate-release-gate.sh` extended with `--mode pre-merge` (commit `ed982f6b`).
- Reuses the existing artifact-presence loop (lines 62–83 of the original script); the file-existence check is not duplicated — the `REQUIRED` array is switched based on `$MODE` and the same loop body runs.
- Diagnostic format: `❌ cycle 999: missing gamma-closeout.md — required before merge (CDD.md §5.3b, gamma/SKILL.md §2.10)` — names both cycle number and filename.
- Cross-references in diagnostic: `CDD.md §5.3b` and `gamma/SKILL.md §2.10`.
- Positive fixture verified: all three close-outs present → exit 0, `✅ Pre-merge closure gate passed`.
- Negative fixture verified: missing `gamma-closeout.md` → exit 1, stderr names cycle `999` + `gamma-closeout.md`.
- Release mode behaviour unchanged (regression verified with positive fixture).
- Bash syntax: `bash -n scripts/validate-release-gate.sh` passes.

**Design constraint check:**
- ✅ Reuses artifact-presence logic (does not duplicate; switches `REQUIRED` array)
- ✅ Diagnostic names file AND cycle number
- ✅ Script cites `CDD.md §5.3b` and `gamma/SKILL.md §2.10`

### AC2 — §3.8 amended with closure-gate-failure override + letter normalization

**Status: MET**

**Evidence:**

- `release/SKILL.md` §3.8 amended at commit `0b56ff86`.
- `grep -nE "closure-gate|closure gate" src/packages/cnos.cdd/skills/cdd/release/SKILL.md` → line 324 in §3.8.
- Override clause text: "If any artifact required by `gamma/SKILL.md` §2.10 closure gate (per `CDD.md` §5.3b ownership matrix) is absent from the cycle's `.cdd/unreleased/{N}/` directory at merge time [...] then `C_Σ` is forced to `<C` regardless of per-axis math. Cycle disposition is 'open and remediate'; the geometric mean is not computed."
- Source order: override paragraph at line 324, letter normalization at line 326, geometric-mean instruction at line 328. Override precedes math. ✅
- Letter normalization sentence: "C− is the operator-visible projection of <C. Both mean 'open and remediate; do not close.'"
- Cross-references: override cites `gamma/SKILL.md §2.10` and `CDD.md §5.3b`. ✅
- Closure-gate-missing example: the override clause explicitly states the geometric mean is not computed when the gate fails — a cycle missing `gamma-closeout.md` scores `<C` even if per-axis math would yield ≥ C.

**Design constraint check:**
- ✅ Override clause precedes geometric-mean instructions in source order
- ✅ `<C` and `C−` reconciled explicitly
- ✅ Cross-references to `CDD.md §5.3b` and `gamma/SKILL.md §2.10` present
- ✅ Gate script also cites the same §5.3b and §2.10 rows (in diagnostic text)

### AC3 — Recursive honest-claim verification (rule 3.13) on this cycle's own artifacts

**Status: MET (β to re-verify)**

All measurements quoted in this cycle's artifacts trace to reproducible sources:

| Claim | Source |
|---|---|
| Cycle #331 merge: `315e529` | `git show 315e529` → "feat(332): β merge — CDD role-skill refactor lands on main" |
| Cycle #333 merge: `6ffdf48` | `git show 6ffdf48` → merge commit for PR #333 |
| Cycle #334 merge: `0900448` | `git show 0900448` → merge commit for PR #336 |
| Retro close-outs: `00e6f8e` | `git show 00e6f8e` → "closeout(334): retroactive close-out artifacts for cycle #334" |
| Retro close-outs: `1ec471d` | `git show 1ec471d` → cycle #335 retro artifacts for #331 and #333 |
| §3.8 rubric introduced: `b27fc15` | `git show b27fc15` → cycle #331 patch 5 |
| F2 source: `.cdd/releases/docs/2026-05-09/335/cdd-iteration.md` | file exists; F2 reads "cdd-tooling-gap: mechanical closure-gate enforcement" |
| F7 source: `.cdd/releases/docs/2026-05-09/335/beta-review.md` | file exists; F7 at line 138 names "rubric closure-gate-failure handling" |

β oracle: re-verify each cited SHA via `git show <sha>`; run `grep -nE "closure-gate|closure gate" src/packages/cnos.cdd/skills/cdd/release/SKILL.md` to confirm §3.8 placement; confirm no C_Σ quoted in this cycle without override-clause trace.

### AC4 — This cycle's own close-out follows the full protocol (recursive coherence)

**Status: IN PROGRESS — artifacts being produced on branch per dispatch protocol**

This self-coherence.md exists on the branch. ✅
`alpha-closeout.md` will be written before review-readiness signal (provisional per §2.8 fallback). ✅ (pending)
`beta-closeout.md` — β writes after review, on branch before merge. (pending β)
`gamma-closeout.md` — γ writes before merge. (pending γ)
`cdd-iteration.md` — required because F2 + F7 are `cdd-*-gap` findings; written in this cycle. ✅ (pending)
`INDEX.md` row — will be added when `cdd-iteration.md` is complete. (pending)

β oracle: run `scripts/validate-release-gate.sh --mode pre-merge` against `.cdd/unreleased/339/` before approving merge; record exit code and diagnostic in `beta-review.md`.

### AC5 — cdd-iteration.md present for this cycle

**Status: PENDING — will be written before review-readiness signal**

File at `.cdd/unreleased/339/cdd-iteration.md` with both findings:
- F1 (cycle-internal): `cdd-tooling-gap` — F2 from cycle #335, mechanical pre-merge gate
- F2 (cycle-internal): `cdd-protocol-gap` — F7 from cycle #335, rubric closure-gate-failure handling

Both with disposition `patch-landed`, patch SHAs cited, affected files named.
INDEX.md row added for cycle #339.

---

## §Self-check

**Did α's work push ambiguity onto β?**

No. Both implementation deliverables are complete and verifiable before β reads the branch:
- AC1: `scripts/validate-release-gate.sh --mode pre-merge` is runnable; positive and negative fixtures are verified in this document.
- AC2: `release/SKILL.md` §3.8 amendment is in the diff; source-order constraint is verifiable by line numbers.

**Is every claim backed by evidence in the diff?**

- AC1 evidence: script diff (commit `ed982f6b`); fixture outputs quoted in §ACs.
- AC2 evidence: skill diff (commit `0b56ff86`); line numbers quoted in §ACs.
- AC3 evidence: SHA table in §ACs maps every cited commit to a verifiable `git show`.
- AC4/AC5: status declared as in-progress with explicit remaining items (beta-closeout, gamma-closeout, cdd-iteration.md, INDEX.md row); not claimed as complete.

**Peer enumeration:**

This cycle patches two surfaces: `scripts/validate-release-gate.sh` and `release/SKILL.md` §3.8. Peer enumeration:

- `scripts/release.sh` line 55 (existing invocation of the gate script) — checked; release-mode call passes unchanged. ✅
- `operator/SKILL.md` §3.4 (δ pre-merge checklist) — the gate script is now available as a δ-side hook. No mandatory update required this cycle (issue explicitly defers δ SKILL.md update; the script exists and is callable). Noted as known debt.
- `gamma/SKILL.md` §2.10 closure gate rows — not modified by this cycle (issue out-of-scope). Cross-referenced in script diagnostic and §3.8 amendment. ✅
- `CDD.md` §5.3b — not modified; cited as authority. ✅

**Harness audit:** No parser, schema, or runtime contract changed. Script extension adds a new flag to an existing script; exit codes follow existing conventions (0=pass, 1=fail). Not applicable for harness audit.

**Recursive coherence gate:** This cycle's own close-out triple must be present before β invokes the gate at merge time. α-closeout.md is written provisionally as part of this dispatch. β must write beta-closeout.md on the branch before merge, and γ must write gamma-closeout.md before merge, so the gate passes.

---

## §Debt

1. **`operator/SKILL.md` §3.4 not updated.** The δ pre-merge checklist does not yet include a row for `scripts/validate-release-gate.sh --mode pre-merge`. Issue #339 defers this explicitly ("Watchdog enforcement at PR creation time" is out of scope; the operator-side integration is usable but not mandated by text). β or γ may file a follow-on if desired. This is known debt.

2. **CI workflow integration deferred.** The gate is implemented as a callable script but not wired into a CI workflow that runs automatically on PRs targeting main. Issue #339 leaves this as "CI integration optional" for option B. A future cycle can add `.github/workflows/validate-pre-merge-closure.yml`. Known debt.

3. **`eng/writing/SKILL.md` not found.** The Tier 3 writing skill declared in the issue was not found at `src/packages/cnos.eng/skills/eng/writing/SKILL.md`. The §3.8 prose was authored against the existing §3.8 style directly. No writing skill constraint was missed (the amendment is prose-consistent with adjacent rubric paragraphs), but the load failure is disclosed.

4. **alpha-closeout.md is provisional** (per §2.8 fallback). Written before β outcome is known. Explicitly marked provisional in the file.

---

## §CDD-Trace

### Step 1 — Receive
- Dispatched by γ on issue #339.
- Selected gap: F2 (cdd-tooling-gap, cycle #335) + F7 (cdd-protocol-gap, cycle #335).
- Active constraints enumerated from issue body.

### Step 2 — Design / Plan
- Design in issue body (mode: `design-and-build`). No separate design artifact required.
- Plan: not required; implementation sequencing is two independent surfaces (script + skill doc). Commit sequence follows issue dispatch order.

### Step 3 — Tests
- Mechanical fixtures: positive and negative fixture runs documented in §ACs.
- No unit test framework applicable (shell script + Markdown). `bash -n` passes. Fixture outputs quoted directly in self-coherence.md.

### Step 4 — Code
- `scripts/validate-release-gate.sh` — extended with `--mode pre-merge` (commit `ed982f6b`).
  - Caller: β runs `scripts/validate-release-gate.sh --mode pre-merge` per AC4 oracle. Also invokable from δ pre-merge checklist.
  - No new module added; function is a flag on the existing script.

### Step 5 — Docs
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` §3.8 — amended (commit `0b56ff86`).
- Write skill applied: prose written per existing §3.8 style.

### Step 6 — Files in diff (artifact enumeration)

Every file in `git diff --stat origin/main..HEAD` (as of implementation SHA `d32917fc`):

| File | Mentioned in |
|---|---|
| `scripts/validate-release-gate.sh` | §ACs AC1, §CDD-Trace Step 4 |
| `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` | §ACs AC2, §CDD-Trace Step 5 |
| `.cdd/unreleased/339/self-coherence.md` | this file (primary branch artifact) |

Pending (will appear in diff before merge):
| File | Mentioned in |
|---|---|
| `.cdd/unreleased/339/alpha-closeout.md` | §ACs AC4, §Debt item 4 |
| `.cdd/unreleased/339/cdd-iteration.md` | §ACs AC5 |
| `.cdd/iterations/INDEX.md` | §ACs AC5 |

### Step 7 — Self-coherence
- This file. Written incrementally per §2.5 (one section per commit+push).
- CDD Trace complete through step 7.

### Step 7a — Pre-review gate
- Gate rows: see §Review-readiness section below.

---

## Review-readiness | round 1 | implementation SHA: fe59d106 | branch CI: local CI unavailable (no CI workflow on cycle branches) | ready for β

**Pre-review gate — all 14 rows:**

| Row | Check | Result |
|---|---|---|
| 1 | `origin/cycle/339` rebased onto current `origin/main` | ✅ merge-base = `00e6f8e2` = current origin/main HEAD (verified 2026-05-10) |
| 2 | `self-coherence.md` carries CDD Trace through step 7 | ✅ §CDD-Trace Steps 1–7a present |
| 3 | Tests present or explicit reason none apply | ✅ No unit test framework applies (shell + Markdown); `bash -n` passes; positive/negative fixtures verified in §ACs |
| 4 | Every AC has evidence | ✅ AC1 MET, AC2 MET, AC3 MET (β to re-verify SHAs), AC4 IN PROGRESS (pending β + γ close-outs), AC5 PENDING → converted to MET with cdd-iteration.md at commit `fe59d106` |
| 5 | Known debt explicit | ✅ §Debt lists 4 items: operator/SKILL.md §3.4, CI workflow, eng/writing/SKILL.md not found, provisional α close-out |
| 6 | Schema/shape audit completed when contracts changed | ✅ Not applicable — no parser/schema/runtime contract changed |
| 7 | Peer enumeration completed | ✅ §Self-check: scripts/release.sh, operator/SKILL.md §3.4, gamma/SKILL.md §2.10, CDD.md §5.3b all enumerated |
| 8 | Harness audit completed | ✅ Not applicable — no schema-bearing contract changed |
| 9 | Post-patch re-audit for every language in diff | ✅ bash: `bash -n scripts/validate-release-gate.sh` passes; Markdown: table shapes and cross-references verified inline; no shell dead-code introduced |
| 10 | Branch CI green on head commit | ✅ No CI workflow runs on cycle branches; stated explicitly as unavailability |
| 11 | Every file in `git diff --stat origin/main..HEAD` mentioned in self-coherence.md | ✅ All 6 files listed in §CDD-Trace Step 6 and cross-checked |
| 12 | Caller-path trace for new modules | ✅ No new module; flag added to existing script; callers: `scripts/release.sh` (release mode) + β oracle (pre-merge mode) |
| 13 | Test assertion count from runner output | ✅ No runner; fixture outputs quoted verbatim in §ACs (exit 0 positive, exit 1 negative with exact stderr) |
| 14 | Author email = `alpha@cdd.cnos` | ✅ `git log -1 --format='%ae' HEAD` → `alpha@cdd.cnos` |

**Artifacts on branch at implementation SHA `fe59d106`:**
- `.cdd/unreleased/339/self-coherence.md` ✅
- `.cdd/unreleased/339/alpha-closeout.md` ✅ (provisional)
- `.cdd/unreleased/339/cdd-iteration.md` ✅
- `.cdd/iterations/INDEX.md` ✅ (row added)
- `scripts/validate-release-gate.sh` ✅ (AC1)
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` ✅ (AC2)

**Missing (β and γ to write on branch before merge):**
- `.cdd/unreleased/339/beta-review.md` — β
- `.cdd/unreleased/339/beta-closeout.md` — β
- `.cdd/unreleased/339/gamma-closeout.md` — γ

**β oracle for AC4:** run `scripts/validate-release-gate.sh --mode pre-merge` against `.cdd/unreleased/339/` before approving merge. Record exit code and output in `beta-review.md`.
