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

---

## Self-check

**Did α's work push ambiguity onto β?** No. Every AC maps to a concrete test assertion
with the exact command that proves it. The validator script and its test are self-contained
and independently runnable.

**Is every claim backed by evidence in the diff?** Yes:
- AC1 claim → `scripts/validate-release-gate.sh` REQUIRED_TRIADIC array + test cases
- AC2 claim → small-change branch in validate-release-gate.sh + AC2 test assertions
- AC3 claim → RELEASE.md check block in validate-release-gate.sh + AC3 test assertions

**Peer enumeration:** The diff touches one concern: pre-tag validation in `scripts/`.
No sibling scripts share the same validation contract. `cn-cdd-verify` is post-release
verification (a different tool with different scope); it is not a peer of
`validate-release-gate.sh`. No peer enumeration required beyond these two files.

**Harness audit:** No schema-bearing contracts changed. `validate-release-gate.sh` reads
filenames from the filesystem; no parser, schema type, or manifest shape is introduced.

---

## Debt

None. The gate correctly exempts small-change cycles per CDD.md §1.2. The existing
cn-cdd-verify tool checking legacy `{role}/CLOSE-OUT.md` paths (rather than the new
per-cycle `{role}-closeout.md` form) is a pre-existing issue outside the scope of #327.

---

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Selected signal: P1 bug — release.sh tags without validating required cycle artifacts; 3.73.0 shipped incomplete |
| 1 Select | — | — | Selected gap: #327 — release gate must validate required cycle artifacts before tag |
| 2 Branch | `cycle/327` | cdd | Branch created by γ; verified `origin/cycle/327` exists |
| 3 Bootstrap | — | cdd | Version directory not required — no CDD docs bundle modified; scripts-only change |
| 4 Gap | `.cdd/unreleased/327/self-coherence.md` §Gap | — | Named incoherence: `release.sh` moves cycle dirs without validating required files |
| 5 Mode | `.cdd/unreleased/327/self-coherence.md` §Skills | cdd, design, write, test, tool | MCA; active skills listed |
| 6 Artifacts | tests (commit `d63830f2`), code (commit `d63830f2`), docs (n/a) | test, tool | design: not required; plan: not required; tests: `scripts/test-validate-release-gate.sh` (19 assertions); code: `scripts/validate-release-gate.sh` + `scripts/release.sh`; docs: no doc surface changed |
| 7 Self-coherence | `.cdd/unreleased/327/self-coherence.md` | cdd | AC-by-AC check completed; self-check done; debt explicit |
| 7a Pre-review | `.cdd/unreleased/327/self-coherence.md` | cdd | Pre-review gate below |

---

## Review-readiness | round 1 | implementation SHA: d63830f2 | branch CI: local gate green | ready for β

**Pre-review gate (alpha/SKILL.md §2.6):**

1. ✅ `cycle/327` rebased onto `origin/main` at `ce315831` (verified at write time; main SHA = merge-base SHA)
2. ✅ `.cdd/unreleased/327/self-coherence.md` carries CDD Trace through step 7
3. ✅ Tests present: `scripts/test-validate-release-gate.sh` — 19 assertions, 19 passed, 0 failed (re-run against HEAD)
4. ✅ Every AC has evidence (§ACs section maps each to concrete test assertions)
5. ✅ Known debt explicit: none beyond pre-existing cn-cdd-verify legacy-path issue
6. ✅ Schema/shape audit: no schema-bearing contracts changed (filesystem path reads only)
7. ✅ Peer enumeration: no peer family affected; `cn-cdd-verify` explicitly scoped out (different concern, different lifecycle phase)
8. ✅ Harness audit: no parser/schema/manifest introduced; n/a
9. ✅ Polyglot re-audit: `bash -n` on all 3 diff'd scripts — syntax-clean; no YAML/Go/Markdown surfaces touched
10. ✅ Branch CI: local gate green (`bash scripts/test-validate-release-gate.sh` → 19/19 pass); remote CI unavailable in this environment — β should wait for CI green before merge
11. ✅ Commit author email: `alpha@cdd.cnos` on all cycle commits (verified via `git log --format='%ae'`)

Branch is review-ready on branch HEAD.
