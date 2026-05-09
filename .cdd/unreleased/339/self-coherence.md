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
