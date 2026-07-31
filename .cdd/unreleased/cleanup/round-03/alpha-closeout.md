# α Closeout — Round 03 (final docs round)

**Role:** α (implementer)
**Branch:** `claude/repo-cleanup-newcomer`
**Input:** `round-03/beta-review.md` (12 target-verified live hits, subset A) + `cnos.core/skills/write/SKILL.md`.
**Scope held:** DOCS/STRUCTURE ONLY. No code, schema, package-skill, or `src/` behavioral content touched. Subset-(C) `docs/gamma/cdd/…` snapshot convention deferred to the code pass.

Falsification gate: re-ran every grep against the live tree and confirmed each replacement target exists on disk before swapping. One minor β/tree discrepancy noted under F14; otherwise the tree matched β.

---

## F12 — MED — dead self-referential Canonical-Path headers (mechanical swap) — DONE

All three targets verified on disk. Swapped each header to its real location.

| File:line | Before | After |
|---|---|---|
| `SUSTAINABILITY.md:6` | `` `docs/beta/SUSTAINABILITY.md` `` | `` `SUSTAINABILITY.md` `` |
| `docs/concepts/lineage/LINEAGE.md:10` | `docs/beta/lineage/LINEAGE.md` | `docs/concepts/lineage/LINEAGE.md` |
| `docs/concepts/lineage/ORIGIN.md:10` | `docs/beta/lineage/ORIGIN.md` | `docs/concepts/lineage/ORIGIN.md` |

---

## F13 — MED — design-index nav at retired dirs + false essays/papers split (coherence edit) — DONE

`docs/development/design/README.md`. Read the whole "Relationship to other doc directories" section per the write skill; made it coherent.

- **L5 intro** — swapped the sibling pointer:
  `These docs sit alongside `` `docs/gamma/essays/` `` (position papers)…` → `` `docs/papers/` `` (position papers). Mechanical.
- **Essays-vs-papers split (former L23 + L25)** — DROPPED the false distinction. The essays it called "position papers" now live in `docs/papers/`, so the `docs/gamma/essays/` bullet duplicated the existing `docs/papers/` bullet (write-skill 3.3 one home per fact; 3.14 don't keep a decorative distinction). Merged into a single `docs/papers/` bullet:
  `- **`` `docs/papers/` ``** — position papers (C≡ / TSC / CTB / cnos stack). Read papers for the *why* and the long-form vision; read design docs for the *what and where* a downstream cycle must pin. Cited from γ-class designs.`
  This **preserves the load-bearing why/vision-vs-what/where distinction** (design docs vs papers — a real class difference) while cutting the decorative essays-vs-papers one.
- **CDD bullet (former L24)** — split the two conflated facts:
  `- **`` `docs/development/cdd/` ``** — CDD-class docs. Read for archival precedent; do not edit historical records. Historical post-release assessments (PRAs) are snapshotted under `` `docs/gamma/cdd/` ``.`
  The **CDD-class-docs pointer** was fixed to its real home `docs/development/cdd/`. The **historical-PRA snapshot clause was NOT blind-swapped** — it retains the `docs/gamma/cdd/` snapshot convention, which is subset-(C) code-pass-owned. See subset-(C) note below.

Result: L3's "not an essay (position paper), not an implementation" class-defining contrast was LEFT — it carries no dead path and defines the decision-artifact class against its two neighboring classes (load-bearing, not decorative). β did not flag it.

### F13 subset-(C) PRA clause — PARTIALLY FIXED (deliberate)

The line mixed two facts. I fixed only the CDD-class-docs pointer (`docs/gamma/cdd/` → `docs/development/cdd/` for the doc class) and kept the PRA snapshot clause referencing the `docs/gamma/cdd/` convention untouched, because that path is the live behavioral snapshot convention owned by package skills/schemas (subset C). This honors the "do not blind-swap that clause" instruction while still making the nav bullet coherent and pointing the reader at the real CDD-class-docs home. The surviving `docs/gamma/cdd/` string in this file is therefore an intentional subset-(C) reference, permitted by the round's verification rule.

---

## F14 — MED — dead pointer to a nonexistent file (authoring drop) — DONE

`docs/reference/ctb/CTB-v4.0.0-VISION.md` §8.5.1.

- Before: `…surfaced four additional language design requirements. Full evidence at `` `docs/gamma/essays/SKILLS-LANGUAGE-EVIDENCE.md` ``.`
- After: `…surfaced four additional language design requirements:`

Dropped the dead external pointer and turned the sentence-final period into a colon so the prose flows directly into the evidence table (the CTB-requirement / practice-evidence / workaround table immediately below **is** the evidence). No file invented.

**β/tree discrepancy (minor, non-blocking):** β stated `git grep -l "SKILLS-LANGUAGE-EVIDENCE"` "should hit only CHANGELOG + this doc." The live tree shows it hits **this doc + three `.cdd/` cleanup artifacts** (`508/pass4-path-inventory.txt`, `508/pass4-triad-token-inventory.txt`, `round-02/gamma-closeout.md`, plus this round's `beta-review.md`) — and NOT CHANGELOG. The load-bearing claim holds regardless: no real doc named `SKILLS-LANGUAGE-EVIDENCE.md` exists anywhere in the tree; the only live-doc occurrence was this dead pointer, now removed. The `.cdd/` hits are review-process inventories, not doc content.

---

## F15 — MED — paper rubric / related cites to retired docs/gamma/ (mechanical swap) — DONE

Both targets verified on disk (`docs/development/ENGINEERING-LEVELS.md`, `docs/development/design/ccnf-o-track-a1-survey.md`).

| File:line | Before | After |
|---|---|---|
| `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md:4` (Rubric) | `docs/gamma/ENGINEERING-LEVELS.md` | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md:8` (Related) | `docs/gamma/ENGINEERING-LEVELS.md` | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/CCNF-AND-TYPED-TRUST.md:8` (related:) | `docs/gamma/ENGINEERING-LEVELS.md` | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/CCNF-AND-TYPED-TRUST.md:458` (References) | `docs/gamma/ENGINEERING-LEVELS.md` | `docs/development/ENGINEERING-LEVELS.md` |
| `docs/papers/CELL-OF-CELLS.md:16` (related:) | `docs/gamma/design/ccnf-o-track-a1-survey.md` | `docs/development/design/ccnf-o-track-a1-survey.md` |

---

## F16 — LOW/policy — LEAVE decisions (no edits) — RECORDED

All left untouched; each re-verified against the live tree.

- `RELEASE.md:18,20` — version-keyed release record (`docs/gamma/essays/`, `docs/gamma/design/` as placement-at-release). **LEFT.** Frozen dated record; describes placement true at the pin, not current authority. Verified present, unchanged.
- `docs/papers/CCNF-AND-TYPED-TRUST.md:383` ("Land this document under `docs/gamma/essays/`…") and `docs/papers/DECREASING-INCOHERENCE.md:528` ("Add this file under `docs/gamma/essays/`…"). **LEFT.** Frozen authoring-intent prose (γ's rule); now-inert but not worth churning. Both verified present, unchanged.
- `docs/reference/governance/DOCUMENTATION-SYSTEM.md:28,121` and `docs/reference/governance/GLOSSARY.md:18` — intentional "these dirs were retired" forbid-statements. **LEFT — no regression.** (Note: these live at `docs/reference/governance/`, not repo root as the finding text abbreviated; content matches β's description exactly.)

---

## Subset (C) — NOT TOUCHED

`docs/gamma/cdd/{X.Y.Z}/…` snapshot convention left untouched in `docs/development/cdd/CDD-PACKAGE-AUDIT.md`, `OVERVIEW.md`, `RATIONALE.md`, `ISSUE-CONSOLIDATION-ANALYSIS.md`, `src/packages/**`, `schemas/cdd/README.md`. The only subset-(C) `docs/gamma/cdd/` string I authored is the deliberate PRA-snapshot clause in `docs/development/design/README.md` (see F13), which references — does not modify — the code-pass-owned convention.

---

## Verification evidence

**(A) Touched-file sweep** — no live-authority dead path remains in any file I edited:

```
grep -n "docs/beta\|docs/gamma/essays\|docs/gamma/design\|docs/gamma/ENGINEERING" \
  SUSTAINABILITY.md docs/concepts/lineage/LINEAGE.md docs/concepts/lineage/ORIGIN.md \
  docs/development/design/README.md docs/reference/ctb/CTB-v4.0.0-VISION.md \
  docs/papers/RELEASE-LEVEL-CLASSIFICATION.md docs/papers/CCNF-AND-TYPED-TRUST.md docs/papers/CELL-OF-CELLS.md
→ single hit: docs/papers/CCNF-AND-TYPED-TRUST.md:383  (F16 FROZEN authoring-intent — expected LEAVE)
```

**(B) Repo-wide docs sweep** — every surviving `docs/beta` / `docs/gamma` (non-cdd) string across `docs/` + root is frozen or intentional:
- Frozen records: `CHANGELOG.md`, `RELEASE.md:18,20`, dated plans (`PLAN-package-system`, `PLAN-v3.22.0-eng-lane-clarity`, `INVARIANT-HARDENING-v1`), audit/smoke evidence logs (`docs/evidence/AUDIT.md:79`, `docs/evidence/smoke/cds-dispatch-smoke-20260623.md`).
- Intentional statements: retirement/redirect notes (`CORE-REFACTOR.md:452,494`, `PACKAGE-SYSTEM.md:9`, `PLAN-package-system.md:6`), forbid-statements (`DOCUMENTATION-SYSTEM.md:28,121`, `GLOSSARY.md:18`).
- F16 authoring-intent: `CCNF-AND-TYPED-TRUST.md:383`, `DECREASING-INCOHERENCE.md:528`.
- No live-authority newcomer-nav / `related:` / Canonical-Path citation to a retired dir remains.

**(C) Subset-(C) survivors** (`docs/gamma/cdd`, expected, code-pass owned): `PLAN-v3.13.0-docs-governance.md`, `PR-docs-governance-v3.13.0.md`, `docs/development/design/README.md` (my PRA clause), `docs/development/cdd/{RATIONALE,OVERVIEW,ISSUE-CONSOLIDATION-ANALYSIS,CDD-PACKAGE-AUDIT}.md`, `CHANGELOG.md`.

**(D) Scope** — `git diff --stat`: 8 files, all docs (+root `SUSTAINABILITY.md`), 12 insertions / 13 deletions. No `src/`, no `schemas/`, no package-skill, no code touched.

---

## Result

F12, F13, F14, F15 remediated; F16 leave-decisions recorded and re-verified. Subset (A) closed: the newcomer docs surface no longer cites any retired `docs/beta/` or `docs/gamma/{essays,design}/` or `docs/gamma/ENGINEERING-LEVELS` path with live authority. Remaining `docs/beta`/`docs/gamma` strings are all frozen records, intentional retirement/forbid statements, or the subset-(C) `docs/gamma/cdd/` convention owned by the code pass. One minor β/tree discrepancy noted (F14 git-grep expectation), non-blocking. Docs-only scope held. Not committed — orchestrator to commit matter + receipt.
