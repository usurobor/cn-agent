---
manifest:
  sections:
    - Gap
    - Skills
    - ACs
    - Self-check
    - Debt
    - CDD-Trace
    - Review-readiness
  completed: [Gap, Skills, ACs]
---

# Self-Coherence — #365 (I6 CDD artifact validator era-mismatch on v3.75.0/v3.76.0)

## Gap

`src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` enforces current
artifact requirements on cycles released before those requirements were in
force. The CI job `cdd-artifact-check` (`I6 CDD artifact ledger validation`)
fails on every push to every branch with **87 errors / 426 total checks**
(reproduced locally on `cycle/365@4fa0bd05`; the issue body cited 77 — drift
is because more v3.75/v3.76 cycles have shipped since the issue was filed).
Every one of the 87 errors is on a `v3.75.0` (126) or `v3.76.0` (2) cycle:

```
src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --all 2>&1 \
  | grep "❌" | grep -oE "v3\.[0-9]+\.[0-9]+" | sort | uniq -c
    126 v3.75.0
      2 v3.76.0
```

The validator already has an `is_legacy_version` helper that downgrades
errors to warnings for older releases. The current cutoff is `minor < 74`
(i.e. cycles pre-`v3.74.0` warn; `v3.74.0`+ strict), but the file header
comment at line 253 claims the cutoff is `v3.65.0`. Both the threshold and
the comment need to move forward. The actual era boundary is `v3.77.0`:
`v3.74.0` already has every required artifact (cycle 327 ships clean),
`v3.75.0` and `v3.76.0` predate the artifact lockdown, and `v3.77.0` is
the next un-shipped minor where strict enforcement becomes correct.

**Mode (selected):** cutoff bump rather than per-artifact era policy. γ's
peer enumeration confirms the failure shape is sharp at a single boundary —
α-closeout/γ-closeout absence collapses to one threshold (`v3.77.0`). A
per-artifact era table would add three knobs all set to the same value;
that is complexity without payoff. The single-function change keeps the
validator simple and matches γ's "lightweight path" framing.

## Skills

| Tier | Skill | Use |
|------|-------|-----|
| 1 | `cdd/CDD.md` | Lifecycle contract; artifact-location matrix §5.3a |
| 1 | `cdd/alpha/SKILL.md` | Role contract; incremental self-coherence §2.5 |
| 2 | `cnos.eng/skills/eng/*` | Tier 2 always-applicable engineering bundle |
| 3 | `cnos.core/skills/write/SKILL.md` | Doc + commit-message style (CHANGELOG row, self-coherence prose) |
| 3 | `cnos.eng/skills/eng/tool/SKILL.md` | Bash standards — `set -euo pipefail`, idempotence, exit codes |
| 3 | `cnos.eng/skills/eng/test/SKILL.md` | Test invariant first; oracle = exit code + warn/fail line for each era boundary |
| 3 | `cnos.eng/skills/eng/ship/SKILL.md` | Bug-fix flow: failing test → fix → passing test → no regressions |

No `eng/{language}` skill is loaded — the validator is a single bash file
and γ confirmed no other language surface is involved. The diff is
bash + bash test harness + Markdown (CHANGELOG row + this file).

## ACs

### AC1 — CI green

**Claim:** `cn-cdd-verify --all` exits 0 on the branch state, and the
`cdd-artifact-check` job in `.github/workflows/build.yml` passes on
`cycle/365`.

**Evidence:**

- Baseline (before fix): 280 passed / 87 failed / 60 warn → exit 1.
- After cutoff bump + fixture relocation: 280 passed / 0 failed / 147
  warn → exit 0 (`src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify
  --all 2>&1 | tail -5`). The 147 warns are 87 era-downgraded misses
  on v3.75/v3.76 cycles plus existing legacy-tag / pre-v3.65.0 warns.
- The CI job runs `cn-cdd-verify --all` then `test-fixtures.sh`. The
  fixture script also passes after the `incomplete-triadic` fixture
  was relocated from v3.75.0 → v3.78.0 (post-cutoff). Both commands
  exit 0 → job passes.

**Caller-path trace (α/SKILL.md §2.6 row 12):** `cn-cdd-verify` is
invoked from `.github/workflows/build.yml` line 269 (`Check CDD
artifacts`) and `test-fixtures.sh` is invoked from line 272 (`Run test
fixtures`). Both are non-test callers in the cdd-artifact-check job.

### AC2 — Era-appropriate handling for v3.75.0 / v3.76.0

**Claim:** Every cycle directory under `.cdd/releases/3.75.0/` and
`.cdd/releases/3.76.0/` missing `alpha-closeout.md`,
`gamma-closeout.md`, or `CDD Trace` section now warns instead of
failing.

**Evidence:**

- `is_legacy_version` cutoff changed from `minor < 74` to `minor < 77`
  (`src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` lines
  251–279). Both downstream warn-string sites updated from `(pre-v3.65.0)`
  to `(pre-v3.77.0)` (lines 388 + 509).
- Concrete v3.75.0 case: cycle 286's `alpha-closeout.md` now reports:
  `⚠️  alpha-closeout.md (issue #286 (v3.75.0)) — missing in legacy
  cycle (pre-v3.77.0) — expected .cdd/unreleased/286 (v3.75.0)/alpha-closeout.md`.
- Concrete v3.76.0 case: cycle 362's missing `CDD Trace` section now
  reports `⚠️  self-coherence.md sections — missing required sections
  in legacy cycle (pre-v3.77.0): CDD Trace`.
- Boundary regression test: `test-cn-cdd-verify.sh` test 11 (v3.75.0)
  + test 12 (v3.76.0), 7 new assertions, all passing.

### AC3 — Strict enforcement preserved for current cycles

**Claim:** v3.77.0+ release cycles and in-flight `.cdd/unreleased/{N}/`
cycles continue to fail strict when artifacts are missing. Cycle 365
itself satisfies strict enforcement at release.

**Evidence:**

- `is_legacy_version` returns 1 (false) for `major == 3 && minor >= 77`
  by code inspection (`cn-cdd-verify` line 274). Cycle 365 will release
  as v3.77.0+ (next minor); it currently sits at `.cdd/unreleased/365/`
  where `check_unreleased_cycle` runs without the legacy downgrade.
- Negative test: `test-cn-cdd-verify.sh` test 13 (v3.77.0 cycle missing
  α/γ close-outs) — asserts exit 1 and `❌ alpha-closeout.md`,
  `❌ gamma-closeout.md`. 3 new assertions, all passing. This is the
  regression sentinel: any future change that bumps the cutoff to
  v3.77.0+ or weakens strict-path enforcement will surface here.
- The relocated `incomplete-triadic` fixture sits at v3.78.0/200 and
  continues to drive `test-fixtures.sh` test 2's exit-1 expectation
  unchanged.

### AC4 — Comment / code coherence

**Claim:** The `is_legacy_version` header comment accurately describes
the cutoff and names a rationale.

**Evidence:**

- New comment block (`cn-cdd-verify` lines 251–276) names the cutoff
  (v3.77.0), the specific artifacts the era policy covers
  (alpha-closeout.md, gamma-closeout.md, "CDD Trace" section), and
  the rationale ("these artifacts were not required when those cycles
  shipped"). Stale `(epoch threshold)` and `v3.65.0` references
  removed.
- Code threshold (`minor < 77` at line 274) matches header. Both
  warn-string sites use `(pre-v3.77.0)` (lines 388, 509). No occurrence
  of the stale `v3.65.0` or `v3.74` literal remains: `grep -n '3\.65\|3\.74'
  src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` → 1 hit, a
  usage-help example `cn cdd-verify --version 3.74.0 --cycle 286` which
  is a documentation reference to a real released cycle (not the
  legacy threshold).

### AC5 — Test fixture pinning the boundary

**Claim:** A test under `test/` exercises the era-appropriate path so a
future regression bumps surface in `test-fixtures.sh` /
`test-cn-cdd-verify.sh`.

**Evidence:**

- `test-cn-cdd-verify.sh` grows 29 → 39 assertions. Three new boundary
  tests (11/12/13) each exercise one side of the cutoff with positive
  + negative proof per `eng/test/SKILL.md` §2.10.
- Two new fixture writers added: `write_legacy_partial_cycle`
  (self-coherence + beta-review only, mirrors v3.75.0 shape) and
  `write_no_cdd_trace_cycle` (full close-outs but no `CDD Trace`,
  mirrors v3.76.0/362 shape).
- Runner: `bash src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh`
  → `ran 39 assertions: 39 passed, 0 failed` (full output captured
  during gate row 13).

### AC6 — CHANGELOG / RELEASE.md note prepared

**Claim:** A line names the validator change so δ does not have to
retro-author release notes at tag time.

**Evidence:**

- `CHANGELOG.md` lines 24–28 — new `## Unreleased` section with one
  bullet identifying #365, the cutoff change, the failure surface
  closed, and the strict / warn split. Doc-only; the table row remains
  δ's to author at tag time. Placement is between the legend (lines
  1–23) and the release table (line 30+), so δ's row insertion remains
  the same operation as prior cycles.

### Artifact enumeration vs. diff (α/SKILL.md §2.6 row 11)

`git diff --stat origin/main..HEAD` lists 9 files. Each is mapped to
an AC or to the §CDD-Trace step 6 entry below:

| File | Covered by |
|------|-----------|
| `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` | AC1, AC2, AC3, AC4 |
| `src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh` | AC5, AC3 (test 13) |
| `src/packages/cnos.cdd/commands/cdd-verify/test/fixtures/incomplete-triadic/.cdd/releases/3.78.0/200/{self-coherence,beta-review}.md` (rename target) | AC1 (fixture relocation required to keep `test-fixtures.sh` test 2 driving an exit-1 path after the cutoff bump) |
| `CHANGELOG.md` | AC6 |
| `.cdd/unreleased/365/self-coherence.md` | this artifact (§Gap..§Review-readiness) |
| `.cdd/unreleased/365/gamma-scaffold.md` | γ scaffold (not authored by α; included in diff because it was added on this cycle branch) |
| `.cdd/unreleased/365/alpha-codex-prompt.md`, `beta-codex-prompt.md` | δ dispatch prompts (not authored by α; included on this cycle branch by δ) |
