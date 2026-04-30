---
cycle: 309
role: beta
version: 3.65.0
---

# β Close-out — Cycle 309

## Review

**Round:** R1
**Verdict:** APPROVED — no findings

All 9 ACs met. Contract integrity clean. No stale paths, no scope drift, no phantom blockers.

- AC1: artifact classification explicit (skill, formula, governing question, named failure mode) ✅
- AC2: five Tier 3 skills loaded before drafting (α attestation in §Skills) ✅
- AC3: all four external sources adapted — IBM, Google SRE, Red Hat, CompTIA ✅
- AC4: eleven-step triage algorithm mapped to six evidence classes ✅
- AC5: hypothesis discipline explicit (state before test, oracle before test, cheapest first, one change, verify) ✅
- AC6: three worked examples — OOM kill, `gh` GraphQL error, background-process kill ✅
- AC7: five RCA handoff triggers; §3.7 prohibits starting RCA during live diagnosis ✅
- AC8: embedded kata with OOM scenario, all required fields ✅
- AC9: frontmatter manually validated against schemas/skill.cue (cue CLI not installed on host — exit 2, not exit 1) ✅

Disclosed debt (D1): `cue` CLI not installed on host; validator exits 2 (prerequisite missing), not 1 (schema failure). Manual validation confirmed schema shape.

## Release

**Version:** 3.65.0 (minor bump — new skill, backwards compatible)
**Merge commit:** main HEAD after `git merge cycle/309 --no-ff`
**Tag:** 3.65.0
**Branch:** cycle/309 → main (`Closes #309`)

Pre-merge gate:
- Row 1 (identity truth): `beta@cdd.cnos` verified before merge commit ✅
- Row 2 (canonical-skill freshness): `origin/main = 0e976257` fetched synchronously before review; no change during session ✅
- Row 3 (non-destructive merge test): throwaway worktree at `/tmp/cnos-merge-test/wt` confirmed zero unmerged paths; validator exit 2 (same as on cycle branch) ✅

Release artifacts:
- `VERSION` bumped to 3.65.0
- `scripts/stamp-versions.sh` updated cn.json and all package manifests
- `scripts/check-version-consistency.sh` passed (7/7 ok)
- `CHANGELOG.md` new ledger row added at top of table
- `RELEASE.md` written
- `.cdd/unreleased/309/` moved to `.cdd/releases/3.65.0/309/`
- Tag `3.65.0` pushed to origin

## For γ

**What the cycle did:** Added `eng/troubleshoot/SKILL.md` — a live diagnosis skill for active failures. Single new file, no boundary change, no schema change. Clean single-round review.

**Evidence:** `beta-review.md` in this directory carries the full R1 verdict.

**Debt to evaluate:** D1 (cue CLI not installed) — environmental; may need CI remediation separately. Not introduced by this cycle.

**No §9.1 triggers fired** — 1 review round, 0 findings, no loaded skill failed to prevent a finding.

**γ owns the PRA.** β's release chain is complete.
