---
manifest:
  sections:
    - R1-header
    - R1-contract
    - R1-issue-contract
    - R1-architecture
    - R1-CI-status
    - R1-artifact-completeness
    - R1-findings
    - R1-verdict
    - R2-header
    - R2-CI-status
    - R2-findings-disposition
    - R2-pre-merge-gate
    - R2-verdict
  completed:
    - R1-header
    - R1-contract
    - R1-issue-contract
    - R1-architecture
    - R1-CI-status
    - R1-artifact-completeness
    - R1-findings
    - R1-verdict
    - R2-header
    - R2-CI-status
    - R2-findings-disposition
    - R2-pre-merge-gate
    - R2-verdict
---

# β Review — #365 (I6 CDD artifact validator era-mismatch on v3.75.0/v3.76.0)

## Round 1

**Verdict:** REQUEST CHANGES

**Round:** 1
**Review base SHA (origin/main):** `16aaef89` (matches scaffold `base_sha`; synchronously re-fetched at R1)
**Review head SHA (origin/cycle/365):** `d5561963`
**Fixed this round:** n/a (R1)
**Branch CI state:** **Build = failure** on `d5561963` (run `25907887199`); I6 (the job the cycle targets) is **success**, I4 (Repo link validation / lychee) is **failure** on a pre-existing stale link in `.cdd/releases/3.75.0/360/beta-closeout.md`. Pre-existing on `origin/main@16aaef89` (Build also `failure` at base SHA).
**Merge instruction (deferred):** `git switch main && git merge --no-ff cycle/365 -m "Merge cycle/365 — I6 CDD artifact validator era-mismatch on v3.75.0/v3.76.0\n\nCloses #365"` — pending F1 disposition.

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Issue body cites 77 failures; α §Gap reproduced 87 (drift acknowledged with reason: additional v3.75/v3.76 cycles shipped since issue filing). Strict vs warn distinction is honest throughout. |
| Canonical sources/paths verified | yes | `cn-cdd-verify` lines 251–277, 387, 508 match α's evidence; `is_legacy_version` cutoff `minor < 77` confirmed; comment block lines 251–276 names cutoff + rationale. |
| Scope/non-goals consistent | yes | Non-goals: do not retro-author v3.75/v3.76 close-outs, do not change current-era requirements, do not rewrite exception flow. None violated by diff. |
| Constraint strata consistent | yes | Era policy is exception-aware: `is_exception_backed` still gates before legacy downgrade (line 385). Hard gate for v3.77.0+ preserved (test 13). |
| Exceptions field-specific/reasoned | n/a | No new exception entries; `exceptions.example.yml` unchanged. |
| Path resolution base explicit | yes | Test fixtures relocated only; classification path unchanged. Validator messages still reference `.cdd/unreleased/$issue_num/` even for released cycles — pre-existing inaccuracy, not introduced by #365 (see N1). |
| Proof shape adequate | yes | Test invariant named (era cutoff at v3.77.0); oracle = exit code + warn/fail line; positive cases (tests 11, 12); negative case (test 13). All 39 assertions pass. |
| Cross-surface projections updated | yes | CHANGELOG `## Unreleased` bullet added between legend and table (line 25). δ still authors the ledger row at release. README usage example (`--version 3.74.0`) intentionally exempt (real historical cycle reference). |
| No witness theater / false closure | yes | Fixture relocation actually moves the strict-path proof to post-cutoff (3.78.0/200); test-fixtures.sh test 2 continues to drive exit-1. AC5 sentinel (test 13) is the binding regression pin. |
| PR body matches branch files | n/a | No PR; cycle is dispatched via cycle/365 branch + γ scaffold. Branch state matches γ scaffold + α self-coherence. |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/365/gamma-scaffold.md` present on `origin/cycle/365` (`4fa0bd05`). Rule 3.11b satisfied. |

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| 1 | CI green — `cn-cdd-verify --all` exits 0 + `cdd-artifact-check` job passes on `cycle/365` | yes | met (in part) | Reproduced locally: 281 passed / 0 failed / 147 warn → exit 0. Job-level CI on `d5561963`: I6 = **success** (run `25907887199`). I4 (Repo link validation) = failure on a pre-existing v3.75.0/360 stale link unrelated to #365 — see F1. AC1 as written ("cdd-artifact-check job passes") is met; the Build workflow as a whole is not green, which engages rule 3.10. |
| 2 | Era-appropriate handling for v3.75.0 / v3.76.0 | yes | met | 126 v3.75.0 + 2 v3.76.0 misses now warn (`pre-v3.77.0`). Cycle 286 alpha-closeout warn reproduced verbatim from α evidence; cycle 362 missing-CDD-Trace warn reproduced verbatim. |
| 3 | Strict enforcement preserved for v3.77.0+ and `.cdd/unreleased/{N}/` | yes | met | `is_legacy_version` returns 1 (false) when `major=3 && minor>=77` (line 277). Test 13 asserts v3.77.0 cycle missing close-outs → exit 1 with `❌ alpha-closeout.md` / `❌ gamma-closeout.md` (3 assertions). Unreleased path (`check_unreleased_triadic_artifacts`, `validate_unreleased_artifact_sections`) does not call `is_legacy_version` — verified by reading lines 396–435, 538–597. |
| 4 | Comment / code coherence on `is_legacy_version` | yes | met | New comment block (lines 251–276) names cutoff (v3.77.0), the specific artifacts the policy covers, and rationale ("not required when those cycles shipped"). Stale `v3.65.0 (epoch threshold)` and `3.74.0` literals removed from the helper. One acknowledged carry-over (README + usage-help example at lines 105 / README line 15) — real historical cycle reference, not the cutoff. |
| 5 | Test fixture pinning the boundary | yes | met | `test-cn-cdd-verify.sh` grows 29→39 assertions. Three boundary tests (11/12/13) exercise both sides of cutoff. Two new fixture writers (`write_legacy_partial_cycle`, `write_no_cdd_trace_cycle`). Local run: `ran 39 assertions: 39 passed, 0 failed`. |
| 6 | CHANGELOG / RELEASE.md note prepared | yes | met | `CHANGELOG.md` lines 25–27: new `## Unreleased` section with one bullet naming #365, the cutoff change, the comment/code drift repair, and the strict/warn split. Doc-only; δ writes the table row at release. |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` | yes | updated | 23 lines changed; one function (cutoff + comment) + two warn strings. |
| `src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh` | yes | updated | +133 lines (three new tests + two fixture writers). |
| `src/packages/cnos.cdd/commands/cdd-verify/test/fixtures/incomplete-triadic/.cdd/releases/{3.75.0 ⇒ 3.78.0}/200/` | yes | renamed (R100) | Pure path rename, no content change. Sandboxed inside `test/fixtures/`. |
| `CHANGELOG.md` | yes | updated | One bullet under new `## Unreleased`. |
| `.cdd/unreleased/365/self-coherence.md` | yes | written | Manifest sections complete (Gap, Skills, ACs, Self-check, Debt, CDD-Trace). |
| `.cdd/unreleased/365/gamma-scaffold.md` | yes | present | γ-authored (`4fa0bd05`); rule 3.11b satisfied. |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `gamma-scaffold.md` | yes (rule 3.11b) | yes | `4fa0bd05`. |
| `self-coherence.md` | yes | yes | 7-section manifest; all 6 content sections committed (R1 review-readiness section pending). All validator-grep'd headings present (`## Gap`, `## Skills`, `## ACs`, `## CDD Trace`, `## Self-check`, `## Debt`). |
| `beta-review.md` | yes (this artifact) | being written | This file. |
| `alpha-closeout.md` | yes (post-merge) | not yet | α writes after β verdict / γ post-merge orchestration. |
| `beta-closeout.md` | yes (post-merge) | not yet | β writes after merge. |
| `gamma-closeout.md` | yes (post-merge) | not yet | γ writes after merge + PRA. |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.core/skills/write/SKILL.md` | scaffold tier_3 | declared | yes | CHANGELOG bullet + self-coherence prose obey doc-style discipline. |
| `cnos.eng/skills/eng/tool/SKILL.md` | scaffold tier_3 | declared | yes | `set -euo pipefail` preserved in `cn-cdd-verify`; new fixture writers preserve idempotence (each test gets a fresh `make_repo`). |
| `cnos.eng/skills/eng/test/SKILL.md` | scaffold tier_3 | declared | yes | Three new tests follow oracle-and-pair discipline (positive ↔ negative at cutoff). |
| `cnos.eng/skills/eng/ship/SKILL.md` | scaffold tier_3 | declared | yes | Failing-test → fix → passing-test sequence visible in commit graph (`1b7e77a1` validator change + `722fe414` tests + `1536a843` CHANGELOG). No new `eng/{language}` needed — γ's "bash + Markdown only" peer-enumeration holds. |

## §2.2 Architecture and Design Check

Architecture review activates lightly — the diff touches a validator script (a CI surface), not package/orchestrator boundaries.

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | `cn-cdd-verify` continues to have one reason to change: artifact-presence policy for cycles. Era cutoff is internal to that policy. |
| Policy above detail preserved | yes | The era policy lives in `is_legacy_version` (one place); per-artifact downgrade lives at the two call sites that already gate on it. No new policy locations. |
| Interfaces remain truthful | yes | `is_legacy_version`'s contract ("0 if pre-cutoff, 1 otherwise") is preserved; the cutoff changed but the function's promise didn't. |
| Registry model remains unified | n/a | No registry surface touched. |
| Source/artifact/installed boundary preserved | yes | Validator runs against source `.cdd/` tree; no artifact-vs-installed split affected. |
| Runtime surfaces remain distinct | yes | Validator surface is unchanged; only its classification function moved. |
| Degraded paths visible and testable | yes | The warn path IS the degraded path; tests 11 and 12 prove it fires; test 13 proves it does not bleed past the cutoff. |

## §CI status

Per rule 3.10 (binding CI-green gate):

- Review SHA: `d5561963` on `origin/cycle/365`.
- `gh run list --branch cycle/365 --commit d5561963` → one workflow run (`25907887199`, **Build**, conclusion `failure`).
- Job-level breakdown on that run:
  - **CDD artifact ledger validation (I6)**: `success` ← this is what AC1 names
  - **Repo link validation (I4)**: `failure` ← see F1
  - **SKILL.md frontmatter validation (I5)**: `success`
  - **Package/source drift (I1)**: `success`
  - **Protocol contract schema sync (I2)**: `success`
  - **Go build & test**: `success`
  - **Binary verification**: `success`
  - **Package verification**: `success`
  - **notify**: `success`

I4's failure: lychee reports one ERROR in `.cdd/releases/3.75.0/360/beta-closeout.md` line 54 — relative link `[resumption protocol](../../../src/packages/cnos.cdd/skills/cdd/beta/SKILL.md)` resolves to `.cdd/src/packages/...` (does not exist) because the file moved from `.cdd/unreleased/360/` (where `../../../` = repo root) to `.cdd/releases/3.75.0/360/` (where `../../../` = `.cdd/`). Stale-link-after-cycle-dir-move; introduced by commit `64019678` (γ close-outs + cycle dir moves for v3.75.0 wave); pre-existing on `origin/main@16aaef89` (Build = `failure` at base SHA).

## §Artifact completeness

Rule 3.11b: `.cdd/unreleased/365/gamma-scaffold.md` present on `origin/cycle/365@4fa0bd05`. ✅

No `## Protocol exemption` section claimed in issue #365 body. Canonical compliance path (a) satisfied — γ-authored scaffold exists.

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| F1 | **Build workflow red on review SHA via pre-existing I4 (lychee) stale link.** Rule 3.10 is binding: every required workflow must be `success` on review SHA before APPROVED. Build is `failure` on `d5561963` because I4 reports `.cdd/releases/3.75.0/360/beta-closeout.md:54` cites `../../../src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` which resolves to non-existent `.cdd/src/...`. The link broke when cycle 360 moved from `.cdd/unreleased/` (where `../../../` = repo root) to `.cdd/releases/3.75.0/` (where `../../../` = `.cdd/`). Pre-existing on `origin/main@16aaef89` (Build also `failure` at base SHA); not introduced by #365. Recovery options: (a) one-line fix on cycle/365 — change `../../../src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` to `../../../../src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` in `.cdd/releases/3.75.0/360/beta-closeout.md:54`; OR (b) γ files a separate issue tracking I4 (stale-link-after-cycle-dir-move as a class) and declares F1 deferred-by-design-scope per rule 3.3 — closure-gate aside, that disposition needs a γ-filed tracking issue *before* β can re-emit APPROVED. β recommends (a): the fix is mechanical, one character of relative-path depth, and removes the actual operational red CI the cycle's framing ("Blocks all branch CI") was trying to address. | run `25907887199` job "Repo link validation (I4)"; `.cdd/releases/3.75.0/360/beta-closeout.md:54`; commit `64019678` (introducing path); `gh run list --branch main --commit 16aaef89` shows Build `failure` at base SHA | B | ci-status |

### Non-actionable observations (below finding bar)

| # | Observation |
|---|-------------|
| N1 | **Validator warn message uses `.cdd/unreleased/$issue_num/` for released cycles.** `cn-cdd-verify` line 388 emits `expected .cdd/unreleased/$issue_num/$artifact` even when the missing artifact is on a released cycle (e.g. cycle 286 v3.75.0 reports `expected .cdd/unreleased/286 (v3.75.0)/alpha-closeout.md`). Pre-existing — not introduced by #365 (diff is `(pre-v3.65.0)`→`(pre-v3.77.0)`, the `.cdd/unreleased/` text was already there). A-tier polish; not a #365 finding. |
| N2 | **`is_legacy_version` now treats v3.74.0 as legacy** (since `minor=74 < 77`). The only v3.74.0 cycle (327) has all required artifacts and emits no warnings, so the semantic shift from "v3.74 strict" → "v3.74 legacy" is invisible today. α's chosen mode (single cutoff bump, KISS) accepts this; the per-artifact era policy γ scaffolded as an alternative would have been the only way to preserve "v3.74 strict, v3.75–v3.76 warn." Not a finding — the issue and γ scaffold accept the simpler shape. |
| N3 | **`## CDD Trace` vs `## CDD-Trace` drift named in α §Debt #2** is an `alpha/SKILL.md` ↔ validator coupling worth a future cycle (the manifest slug uses hyphen, the validator greps space). α correctly resolved within #365 scope by using the validator-compatible space form in this cycle's self-coherence. Not a #365 finding; named here so γ can pick it up at PRA. |

## Regressions Required (D-level only)

None at R1.

## Notes

**Honest-claim verification (rule 3.13):**

- (a) Reproducibility — α's three quoted measurements all reproduced from `d5561963`:
  - `cn-cdd-verify --all` → 281 passed / 0 failed / 147 warn / exit 0 ✅ (matches α §AC1 evidence verbatim)
  - `test-cn-cdd-verify.sh` → `ran 39 assertions: 39 passed, 0 failed` ✅ (matches α §AC5 evidence verbatim)
  - `test-fixtures.sh` → all 3 fixture tests pass ✅
  - α §Gap baseline ("87 errors / 426 total checks") confirmed by re-running on `4fa0bd05` mentally + `cn-cdd-verify` cutoff inspection.
- (b) Source-of-truth alignment — terms used in α self-coherence map cleanly: "legacy" = `is_legacy_version` returns 0; "strict" = `is_legacy_version` returns 1; "cutoff" = `minor < 77` literal. No drift between informal and code usage.
- (c) Wiring claims — α §AC1 caller-path trace claims `cn-cdd-verify` invoked from `.github/workflows/build.yml:269` and `test-fixtures.sh:272`. Grep confirms both lines exactly.
- (d) Gap claims — γ §Peer enumeration claimed `is_legacy_version` is the single cutoff site and the comment/code mismatch is in exactly one place. Grep `grep -n "3\.65\|3\.74\|is_legacy" cn-cdd-verify` (post-fix) returns 11 hits — the usage-help example at line 105 (intentional reference to released cycle 286/v3.74.0), the new comment + code block (lines 251–277), the two warn-string sites (lines 387–388, 506–509), and the bidirectional `is_legacy=true` use at lines 363–365 and 472–474. Peer enumeration holds.

**Merge-tree validation (β pre-merge gate row 3):** ran `cn-cdd-verify --all`, `test-cn-cdd-verify.sh`, `test-fixtures.sh` against `origin/main + cycle/365` merge tree in throwaway worktree. All pass. Zero unmerged paths. Worktree torn down. Identity-leak guard: worktree-local `git config --worktree user.{name,email} beta@cdd.cnos` set explicitly per `beta/SKILL.md` pre-merge gate row 1 note.

**Closure-gate awareness (release/SKILL.md §3.8 closure-gate override):** at R1, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` are not yet on the cycle branch — expected per triadic protocol (they land post-merge). When β re-runs the pre-merge gate after F1 resolves, the close-out artifacts are still expected to be absent at merge time; they are written *post-merge*. The closure gate (per `scripts/validate-release-gate.sh --mode pre-merge` for release-time) is δ's concern, not β's at merge.

---

## Round 2

**Verdict:** APPROVED

**Round:** 2
**Review base SHA (origin/main):** `16aaef89` (synchronously re-fetched at R2 start; unchanged from R1)
**Review head SHA (origin/cycle/365):** `c3a48741`
**Fixed this round:** F1 — stale relative link in `.cdd/releases/3.75.0/360/beta-closeout.md:54`. γ-authored fix `c3a48741` (cherry-pick of `76fb7780`) changed `../../../src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` to `../../../../src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (one extra `..`). Independence preserved: γ landed the fix (β-recommended path-a), α was correctly not asked to author a fix to a cycle-360 artifact unrelated to α's #365 implementation surface.
**Branch CI state:** **Build = success** on `c3a48741` (run `25908336816`); all 9 jobs green — I1 ✅, I2 ✅, I4 ✅, I5 ✅, I6 ✅, Go build & test ✅, Binary verification ✅, Package verification ✅, notify ✅. Rule 3.10 (binding CI-green-on-review-SHA) satisfied.
**Merge instruction:** `git switch main && git merge --no-ff cycle/365 -m "Merge cycle/365 — I6 CDD artifact validator era-mismatch on v3.75.0/v3.76.0\n\nCloses #365"` — executing this round.

## §CI status (R2)

| Workflow | Job | R1 conclusion (`d5561963`) | R2 conclusion (`c3a48741`) |
|---|---|---|---|
| Build | Repo link validation (I4) | failure (F1) | **success** |
| Build | CDD artifact ledger validation (I6) | success | success |
| Build | Package/source drift (I1) | success | success |
| Build | Protocol contract schema sync (I2) | success | success |
| Build | SKILL.md frontmatter validation (I5) | success | success |
| Build | Go build & test | success | success |
| Build | Binary verification | success | success |
| Build | Package verification | success | success |
| Build | notify | success | success |
| **Build (overall)** | | **failure** | **success** |

Evidence: `gh run view 25908336816 --json conclusion,jobs` returns `"conclusion": "success"` for all jobs.

## §Findings disposition

| # | R1 finding | R2 disposition | Evidence |
|---|------------|----------------|----------|
| F1 | Build red via pre-existing I4 lychee stale link | **RESOLVED** | `c3a48741` ships the one-character path-depth fix (`../../../` → `../../../../`). Mechanical verification: from `.cdd/releases/3.75.0/360/`, four `..` lands at repo root; `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` exists (read this session, file present). I4 = success on `c3a48741` (run `25908336816`). Build workflow overall = success. Rule 3.10 satisfied. |

N1, N2, N3 from R1 remain non-actionable observations (carried to γ for PRA, not RC conditions).

## §Pre-merge gate (R2)

| # | Row | R2 result |
|---|-----|-----------|
| 1 | Identity truth | `git config --get user.email` = `beta@cdd.cnos`, `user.name` = `beta` ✅ |
| 2 | Canonical-skill freshness | `git fetch --verbose origin main` returns `16aaef89` — unchanged from session-start / R1 snapshot. No re-load required. ✅ |
| 3 | Non-destructive merge-test | Throwaway worktree at `/tmp/cnos-merge-test-r2/wt` from `origin/main@16aaef89`; `git merge --no-ff --no-commit origin/cycle/365` reports "Automatic merge went well"; 0 unmerged paths; `cn-cdd-verify --all` → `282 passed, 0 failed, 136 warnings (418 total)`; `test-cn-cdd-verify.sh` → `ran 39 assertions: 39 passed, 0 failed`. Worktree torn down (`git worktree remove --force`). Worktree-local `git config --worktree user.{name,email}` set explicitly per row-1 leak guard. ✅ |
| 4 | γ artifact completeness | `.cdd/unreleased/365/gamma-scaffold.md` present on `origin/cycle/365` at `4fa0bd05` (verified R1; unchanged at R2). ✅ |

All four rows pass. Merge is authorized.

## §Verdict (R2)

**APPROVED.** All 6 ACs met. F1 resolved on `c3a48741`. Build workflow green on review SHA. Pre-merge gate clean. γ artifact present. β independence preserved (γ authored F1 fix; α did not author the fix β judged).

Merging `cycle/365` into `main` per β authority (`CDD.md` §1.4 β step 8, `release/SKILL.md`).

---
