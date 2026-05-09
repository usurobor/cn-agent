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
