## Gap

**Issue:** #327 — bug(cdd): release gate must validate required cycle artifacts before tag
**Version/mode:** MCA — bug fix to `scripts/release.sh` (pre-tag validation gate)
**Priority:** P1

`release.sh` moves `.cdd/unreleased/{N}/` to `.cdd/releases/{X.Y.Z}/{N}/` but does not
validate whether the directory contains the required lifecycle artifacts before tagging.
3.73.0 shipped with `.cdd/releases/3.73.0/325/` missing `alpha-closeout.md` and
`gamma-closeout.md` — both required by the CDD lifecycle that 3.73.0 itself introduced.

The fix: add `scripts/validate-release-gate.sh` that blocks the release if required
artifacts are missing, and call it from `release.sh` before stamp/move/tag.

Design artifact: not required — single-concern tooling addition (2 new scripts +
3-line change to release.sh); no impact graph beyond the 3 scripts it touches.

Plan artifact: not required — implementation sequencing is trivial (test → implement →
wire in).

---

## Skills

**Tier 1a (CDD authority):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2 (general engineering):**
- Tooling bundle applies: all outputs are shell scripts in `scripts/`

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/design/SKILL.md` — design principles applied to
  validation contract boundary (one reason to change, truthful interface)
- `src/packages/cnos.core/skills/write/SKILL.md` — script header comments and
  self-coherence prose
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — test from invariants, not just
  examples; negative space mandatory; oracle explicit
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` — shell script standards:
  set -euo pipefail, fail-fast, idempotent, zero runtime deps, machine-readable output

---

## ACs

### AC1: release.sh validates required artifacts before tag

**Status: MET**

Evidence: `scripts/validate-release-gate.sh` (commit `d63830f2`) implements the gate.
Called by `scripts/release.sh` at step 4, before stamp/move/tag.

Positive: `test-validate-release-gate.sh` test "AC1 positive — triadic cycle with all
required artifacts" → exit 0. Output contains `✅ cycle 200 (triadic): all required
artifacts present` and `✅ Release gate passed`.

Negative (missing alpha-closeout.md): test "AC1 negative — missing alpha-closeout.md"
→ exit 1. Output contains `❌ cycle 201 (triadic): missing alpha-closeout.md` and
`❌ Release gate FAILED`.

Negative (missing gamma-closeout.md): test "AC1 negative — missing gamma-closeout.md"
→ exit 1. Output contains `❌ cycle 202 (triadic): missing gamma-closeout.md` and
`❌ Release gate FAILED`.

All 19 assertions in `test-validate-release-gate.sh` pass: `19 passed, 0 failed`.

### AC2: small-change path explicitly collapses artifact set

**Status: MET**

Evidence: CDD.md §1.2 already contains the explicit artifact collapse table for
small-change cycles. `scripts/validate-release-gate.sh` follows it: a cycle dir without
`beta-review.md` is classified as small-change and no cycle-dir artifacts are required
(header comment references CDD.md §1.2).

Positive: test "AC2 positive — small-change cycle passes without triadic files" → exit 0.
Output contains `✅ cycle 203 (small-change): no required cycle-dir artifacts (CDD.md §1.2)`.

Negative: test "AC2 negative — small-change cycle not blocked for missing triadic files"
re-uses the AC2 positive output and asserts `Release gate passed` — confirming the
small-change cycle is not blocked for absent `beta-closeout.md`.

### AC3: missing RELEASE.md blocks release

**Status: MET**

Evidence: `scripts/validate-release-gate.sh` checks `RELEASE.md` existence at repo root
as its first check; exits 1 if missing. Called by `release.sh` step 4, before tag.

Positive: test "AC3 positive — RELEASE.md present" → exit 0. Output contains
`✅ RELEASE.md present`.

Negative: test "AC3 negative — RELEASE.md absent" → exit 1. Output contains
`❌ RELEASE.md missing at repo root — required before tag` and `❌ Release gate FAILED`.
