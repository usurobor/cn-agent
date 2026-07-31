# β Review — Round 03 (final docs-surface pass)

**Role:** β (independent adversarial reviewer)
**Branch:** `claude/repo-cleanup-newcomer`
**Inputs:** `round-02/gamma-closeout.md` (Round-3 target = subset A), round-01/02 beta-review + alpha-closeout, `cnos.core/skills/write/SKILL.md`.
**Method:** Re-derived every hit from the live tree. Ran the greps myself, opened each cited file:line, confirmed each replacement target exists on disk. Edited only this file.

---

## Verdict: **REVISE** — but this is the **LAST docs round.**

Subset (A) is a bounded, enumerable, target-verified sweep of the exact F9 class (live docs citing dead `docs/beta/` / `docs/gamma/` paths). **12 live-authority hits remain.** After α applies F12–F15, the newcomer docs surface reaches **PRISTINE**; subset (C) `docs/gamma/cdd/…` stays with the code pass. No regression in F1–F11.

---

## Final newcomer impression

The primary path reads clean end-to-end. `README.md` opens with the point ("cnos gives AI agents a Git-native home"), separates shipped (`cn` hub/package commands) from planned (agent runtime) without contradiction, and names `src/go/` + `src/packages/` correctly — no `src/agent/` leak. `docs/README.md` is a single intent-nav map (`quickstart/ concepts/ guides/ reference/ development/ architecture/ papers/ evidence/`), consistent with `DOCUMENTATION-SYSTEM.md` and the real tree. `THESIS.md` front-loads coherence as the root concept and defers definitions to the Glossary. Signal/noise is high; authority is explicit.

The **only** thing marring it is subset (A): live docs (position papers, the design index, the CTB vision, three lineage/sustainability headers) still cite retired `docs/beta/` / `docs/gamma/` directories that no longer exist on disk. A sharp reader who clicks a `related:` cite or a Canonical-Path header hits a dead path. That is the residual, and it is the whole residual.

---

## Findings

### F12 — MED — Self-referential dead Canonical-Path headers (swap table)

Three live docs declare their own canonical location as a retired `docs/beta/` path. Highest-signal class: a doc naming its own dead home (write-skill 3.13 "keep authority explicit" — the stated authority is false). All targets confirmed on disk.

| File:line | Reads | Correct on-disk target |
|---|---|---|
| `SUSTAINABILITY.md:6` (repo root) | `docs/beta/SUSTAINABILITY.md` | `SUSTAINABILITY.md` |
| `docs/concepts/lineage/LINEAGE.md:10` | `docs/beta/lineage/LINEAGE.md` | `docs/concepts/lineage/LINEAGE.md` |
| `docs/concepts/lineage/ORIGIN.md:10` | `docs/beta/lineage/ORIGIN.md` | `docs/concepts/lineage/ORIGIN.md` |

**Fix:** swap each to its real path (mechanical).

### F13 — MED — Live design-index nav points at retired dirs + carries a now-false essays/papers split

`docs/development/design/README.md` is a live newcomer nav doc. Its "Relationship to other doc directories" section (and L5 intro) cite dead dirs. The essays it calls "position papers" now live in `docs/papers/` (confirmed: `CELL-OF-CELLS.md`, `DECREASING-INCOHERENCE.md`, `BOX-AND-THE-RUNNER.md`, `CCNF-AND-TYPED-TRUST.md`, all `AGENT-*` essays are under `docs/papers/`), so the `docs/gamma/essays/` row is not just a stale path — it duplicates the existing `docs/papers/` row (L26) and asserts a distinction that no longer exists (write-skill 3.3: one home per fact; 3.14: don't keep a decorative distinction).

| File:line | Reads | Fix |
|---|---|---|
| `.../design/README.md:5` | `docs/gamma/essays/` (position papers) | → `docs/papers/` |
| `.../design/README.md:23` | `docs/gamma/essays/` (position papers) | → merge into / point at `docs/papers/`; drop the false essays-vs-papers split |
| `.../design/README.md:24` | `docs/gamma/cdd/` (CDD-class docs + historical PRAs) | CDD-class docs → `docs/development/cdd/`; the historical-PRA snapshot location is subset (C) — coordinate with the code pass, do not blind-swap that clause |

### F14 — MED — Dead pointer to a file that exists nowhere (not a swap — authoring fix)

`docs/reference/ctb/CTB-v4.0.0-VISION.md:341`: "Full evidence at `docs/gamma/essays/SKILLS-LANGUAGE-EVIDENCE.md`." That file exists **nowhere** in the tree (git grep: only CHANGELOG + this doc name it). This is worse than a stale path — it is a reference to a document that was never created. The evidence is the table inline immediately below (L342+), so the pointer is also redundant. Cannot be a mechanical path swap. **Fix:** drop the dead external pointer — the following table *is* the evidence ("surfaced four additional language design requirements:" → table). If the content is genuinely meant to live in a separate doc, that doc must be authored; until then the cite must not name a nonexistent file.

### F15 — MED — Live paper rubric / `related:` cites to retired `docs/gamma/` (swap table)

Current position papers cite retired `docs/gamma/` paths in live rubric lines and `related:` / `References` frontmatter. A reader following the cite lands nowhere. Both targets confirmed on disk (`docs/development/ENGINEERING-LEVELS.md`, `docs/development/design/ccnf-o-track-a1-survey.md`).

| File:line | Reads | Correct target |
|---|---|---|
| `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md:4` | `docs/gamma/ENGINEERING-LEVELS.md` (Rubric) | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md:8` | `docs/gamma/ENGINEERING-LEVELS.md` (Related) | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/CCNF-AND-TYPED-TRUST.md:8` | `docs/gamma/ENGINEERING-LEVELS.md` (related:) | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/CCNF-AND-TYPED-TRUST.md:458` | `docs/gamma/ENGINEERING-LEVELS.md` (References) | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/CELL-OF-CELLS.md:16` | `docs/gamma/design/ccnf-o-track-a1-survey.md` (related:) | `docs/development/design/ccnf-o-track-a1-survey.md` |

### F16 — LOW / policy — Adjudications: leave the frozen / authoring-intent lines

Recorded so α does not over-reach and the loop can close cleanly. **All LEAVE** (consistent with the round-2 frozen-vs-live rule that left the identical CHANGELOG cites untouched):

- `RELEASE.md:18,20` (root) — cite `docs/gamma/essays/` and `docs/gamma/design/` as placement-at-release. **Ruling: FROZEN — leave.** The file is a version-keyed release record (`# 3.82.0`), the same class as a promoted CHANGELOG entry; it describes placement true at the pin, not current authority. Touching it would churn a dated record.
- `docs/papers/CCNF-AND-TYPED-TRUST.md:383` ("Land this document under `docs/gamma/essays/`…") and `docs/papers/DECREASING-INCOHERENCE.md:528` ("Add this file under `docs/gamma/essays/`…") — frozen authoring-intent prose (γ's rule). Leave. (Minor: both instructions are now satisfied — the docs landed at `docs/papers/` — so they are inert; not worth an edit.)
- `DOCUMENTATION-SYSTEM.md:28,121`, `GLOSSARY.md:18` — intentional "these dirs were retired" forbid-statements. Correct as-is; **verified no regression.**

---

## Subset (C) — code-pass-deferred, NOT filed as docs findings

The `docs/gamma/cdd/{X.Y.Z}/…` snapshot-path convention is a live behavioral contract in package skills/schemas — do not blind-swap in docs. Deferred hits: `docs/development/cdd/CDD-PACKAGE-AUDIT.md` (~30), `docs/development/cdd/OVERVIEW.md:6`, `docs/development/cdd/RATIONALE.md:6` (placement headers), and the `src/packages/**` + `schemas/cdd/README.md` cites. `docs/development/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md:452` is a frozen signed record — leave. These stay with the code pass, coordinated, per γ's split.

---

## Regression re-verify (rounds 1–2)

**No regression.** All survivors are frozen or intentional:

- `src/agent/` in docs: every hit is a dated CHANGELOG row, a dated plan (`INVARIANT-HARDENING`, `PLAN-package-system`, `PLAN-v3.22`, `PR-docs-governance`), a "this migration has landed / `src/agent/` is gone" retirement note (`BUILD-AND-DIST.md`, `CORE-REFACTOR.md`), a 404-noted absence (`JFA-cycle-log-dyad.md`), or a `MODULAR-ARCHITECTURE-REFACTOR.md` design record. Zero live-authority `src/agent/` cite in newcomer docs. (Remaining `src/packages/cnos.core/skills/**` self-refs are code-pass behavioral content, out of scope.)
- `docs/alpha/` on live docs: every survivor frozen/intentional — CHANGELOG, dated plans, `WRITER-PACKAGE.md` (Superseded evidence), `EXECUTABLE-SKILLS.md:398` (bibliography), `DESIGN-266` (frozen design), plus the retirement forbid-statements. Clear.
- `cnos.pm`: only CHANGELOG (frozen), `CAR-implementation-plan.md` (Implemented plan), `PACKAGE-RESTRUCTURING.md` (bannered Draft). Clear.
- Root `README.md` / `CLI.md` / `OPERATOR.md`: shipped-vs-planned still agree — README ships the `cn` hub/package commands and explicitly defers the agent runtime; no `src/agent/`. `DOCUMENTATION-SYSTEM.md` intent-layout still matches the tree and `docs/README.md`. Correct.

---

## Convergence assessment

**After α applies F12–F15, the docs reader surface is PRISTINE. Stop the docs loop.**

The residual is exactly what β forecast in round 1 and γ confirmed in round 2: one bounded, target-verified sweep of the F9 class. It contains no new class of drift — no contradiction, no two-jobs file, no orphaned structure — only stale/dead path citations with confirmed correct targets (F12, F13, F15 mechanical; F14 one authoring drop). F16 is leave-decisions, no edits. Once these land:

- every live Canonical-Path header names its real location,
- every live nav/reference/`related:` cite resolves,
- no live doc points at a `docs/beta/` or `docs/gamma/` directory that no longer exists,
- the only remaining `docs/gamma/…` strings are frozen records (CHANGELOG, RELEASE, dated plans, sigma logs), intentional retirement statements, or the subset-(C) `docs/gamma/cdd/` convention owned by the code pass.

No further material docs drift remains after F12–F15. This is the final docs round.
