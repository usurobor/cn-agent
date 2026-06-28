# Self-Coherence — cnos#510 Pass 4B (R0)

**Role:** α (collapsed onto δ per β-α-collapse-on-δ / PERSONA.md commitment 5)
**Round:** R0
**Head:** 1433203 (Pass 4B commit)
**Base:** 6b24ab6

---

## AC Verification

### AC1 — git name-status shows R (rename) entries for moved files

**Status: PARTIAL / CALIBRATION REQUIRED**

`git diff HEAD~1 HEAD --name-status -M` shows no `R` entries. Old paths appear as `M` (modified — stub content replaced original); new paths appear as `A` (added — full content at new location). Total: 35 A + 40 M = 75 entries.

**Root cause:** AC1 and AC2 conflict structurally. Git rename detection identifies a rename when the old path is *deleted* and a new path is *created* with similar content (≥50% similarity). When a stub is placed at the old path (AC2), the old path is not deleted — it has new (stub) content. The old content is gone from the old path; git sees it as modified, not renamed. Git cannot produce R entries when stubs occupy the old paths.

**Calibration note:** This tension is inherent to the stub pattern established in Pass 2/3. AC1 as written cannot be satisfied simultaneously with AC2. The effective signal of "files moved" is: A entries at new paths with full original content, M entries at old paths with stub-only content. The reviewer can verify move correctness via `git show HEAD:<new-path>` and `git show HEAD:<old-path>` rather than relying on R-entry classification. AC1's intent (traceable move) is satisfied via the commit structure and stub content; the literal R-entry requirement is architecturally incompatible with AC2.

---

### AC2 — `# Moved` redirect stub at every old path

**Status: PASS**

34 stub files written at old paths across all four bundles:
- `docs/beta/guides/` — 7 stubs (AUTOMATION, BUILD-RELEASE, DOJO, HANDSHAKE, MIGRATION, TROUBLESHOOTING, WRITE-A-SKILL)
- `docs/beta/evidence/` — 8 stubs (AUDIT.md + 7 rca/ files)
- `docs/beta/lineage/` — 3 stubs (LINEAGE, ORIGIN, README)
- `docs/alpha/doctrine/` — 16 stubs (README + 15 essay/log files across 4 subdirs)

Each stub format: `# Moved\n\nThis file has moved to [\`<new-path>\`](<relative-link>).\n`

Spot-check: `git show HEAD:docs/beta/guides/AUTOMATION.md` → stub with relative link to `../../guides/AUTOMATION.md`.

---

### AC3 — active markdown links target new homes

**Status: PASS (after δ accept-with-repair — 4 active links were initially missed; see δ note at end)**

Repointed active markdown links in:
- `docs/guides/README.md` — 4 guide links updated; Pass 1 overlay blockquote removed
- `docs/evidence/README.md` — 2 evidence links updated; Pass 1 overlay blockquote removed
- `docs/concepts/README.md` — doctrine and lineage links updated; Pass 1 overlay blockquote removed
- `OPERATOR.md` — 4 guide links (AUTOMATION ×2, BUILD-RELEASE, HANDSHAKE, TROUBLESHOOTING) updated
- `docs/papers/ACTIVATION-NOT-DEPLOYMENT.md` — 5 absolute GitHub URLs updated (doctrine/README + 4 essay files)
- `docs/papers/AGENT-CONTINUITY.md` — 5 absolute GitHub URLs updated (doctrine/README + 4 essay files)

Relative link depths adjusted:
- `docs/guides/WRITE-A-SKILL.md`: `../../../src/` → `../../src/` (7 occurrences)
- `docs/guides/DOJO.md`: `../../../src/` → `../../src/` (5 occurrences)
- `docs/evidence/AUDIT.md`: `../architecture/ARCHITECTURE.md` → `../beta/architecture/ARCHITECTURE.md`

---

### AC4 — `docs/beta/lineage/3.14.4/` left in place

**Status: PASS**

`docs/beta/lineage/3.14.4/` is a versioned snapshot (`**/X.Y.Z/`) explicitly excluded from cell scope. Verified: snapshot directory untouched. Only LINEAGE.md, ORIGIN.md, README.md moved from `docs/beta/lineage/`.

---

### AC5 — no edits to out-of-scope files

**Status: PASS**

Changed files: `.cdd/unreleased/510/gamma-scaffold.md` (A), `OPERATOR.md` (M), 34 old-path stubs (M), 35 new-path files (A), 3 intent-index READMEs (M), 2 docs/papers/ files (M), 3 moved files with relative link repairs.

No edits to `.github/workflows/`, `.cn-sigma/logs/`, or any file outside the declared surfaces.

---

### AC6 — Pass 1 overlay blockquotes removed from updated intent-index READMEs

**Status: PASS**

All three intent-index READMEs (`docs/guides/README.md`, `docs/evidence/README.md`, `docs/concepts/README.md`) had their Pass 1 overlay blockquotes removed and now contain direct navigation links to the canonical new-path locations.

---

### AC7 — no active stale links to 4B bundle old paths remain

**Status: PASS (after δ accept-with-repair — initially over-claimed; see δ note at end)**

`git grep -nE 'docs/(beta/guides|beta/evidence|beta/lineage|alpha/doctrine)/'` run over `docs/`, `OPERATOR.md`. All remaining hits classified:

| Source | Reference | Classification |
|--------|-----------|----------------|
| `docs/alpha/design/WRITER-PACKAGE.md` | Backtick path citations to `docs/alpha/doctrine/` | AC7 exempt: backtick citations in design doc |
| `docs/alpha/ctb/` files | Backtick path citations to `docs/alpha/doctrine/` | AC7 exempt: backtick citations in design docs |
| `docs/beta/governance/GLOSSARY.md` | `docs/alpha/doctrine/COHERENCE-SYSTEM.md` (not in 4B bundle) | Not a 4B-moved file; pre-existing |
| `docs/beta/governance/3.14.4/` | All references | Versioned snapshot — out of scope |
| `docs/beta/lineage/3.14.4/` | All references | Versioned snapshot — out of scope |
| `docs/concepts/lineage/LINEAGE.md`, `ORIGIN.md` | Canonical-Path metadata rows | AC10 frozen content; AC7 exempt |
| `docs/concepts/doctrine/*/` | Absolute GitHub URLs within moved essays to other doctrine files | AC10 prohibits prose edits; AC7 frozen/historical |
| `docs/evidence/AUDIT.md:79` | Plain-text path `docs/beta/evidence/rca/` in table cell | AC10 frozen content; not a markdown link |
| `docs/THESIS.md:703` | `docs/alpha/doctrine/FOUNDATIONS.md` (not in 4B bundle) | Pre-existing; not a 4B-moved file |

No active navigation links to 4B old paths remain in non-frozen active docs.

---

### AC8 — `docs/beta/guides/MIGRATION.md` present as stub

**Status: PASS**

Stub exists at `docs/beta/guides/MIGRATION.md`. Content: redirect to `../../guides/MIGRATION.md`.

---

### AC9 — no bidi/hidden characters in changed files

**Status: PASS**

`LC_ALL=C grep -rlP '[\x00-\x08\x0b\x0c\x0e-\x1f\x7f\x{200b}-\x{200d}\x{feff}\x{202a}-\x{202e}\x{2066}-\x{2069}]'` over all changed files: no matches.

---

### AC10 — moved files' content unchanged except required relative-link repairs

**Status: PASS**

Three files received relative-link repairs (not prose edits):
- `docs/guides/WRITE-A-SKILL.md`: `../../../src/` → `../../src/` — mechanical depth correction from move
- `docs/guides/DOJO.md`: `../../../src/` → `../../src/` — mechanical depth correction from move
- `docs/evidence/AUDIT.md`: `../architecture/` → `../beta/architecture/` — mechanical target correction from move

All other 32 moved files: content byte-for-byte identical to original (stubs at old paths preserve original content at new paths only). No prose reflow, no semantic edits.

---

## Summary

9/10 ACs pass cleanly. AC1 is architecturally incompatible with AC2 (the stub pattern). The intent of AC1 (traceability of moves) is satisfied via commit structure and stub content; the literal R-entry requirement cannot be satisfied when stubs occupy old paths. This is a known limitation of the stub pattern established in Pass 2/3 and should be addressed in the cell spec as a known tension.

**Review-ready signal: R0 — ready for β.**

---

## δ accept-with-repair (post-β, operator-authorized 2026-06-28)

β converged claiming AC3/AC7/AC8 PASS. An independent δ boundary check (operator-directed) found the AC7 oracle was too narrow: it grepped `docs/(beta/guides|beta/evidence|beta/lineage|alpha/doctrine)/` (the `docs/`-prefixed form) and therefore **missed relative markdown links** of the form `../beta/...` and `../../beta/...`. Four active stale links to old 4B paths remained:

- `docs/quickstart/README.md:14` — `[Operator handshake](../beta/guides/HANDSHAKE.md)` (Pass 1 portal index → AC8)
- `docs/quickstart/README.md:15` — `[Automation](../beta/guides/AUTOMATION.md)` (Pass 1 portal index → AC8)
- `docs/alpha/cli/DAEMON.md:125` — `[AUTOMATION.md](../../beta/guides/AUTOMATION.md)`
- `docs/alpha/security/TRACEABILITY.md:743` — `[AUDIT.md](../../beta/evidence/AUDIT.md)`

These resolved via the redirect stubs (not broken), but were active stale links, so AC3/AC7/AC8 were over-claimed.

**Disposition: δ accept-with-repair (NOT status:changes).** A 4-line repair commit repointed all four to their new homes (`../guides/...`, `../../guides/...`, `../../evidence/...`). The AC7 grep was re-run including the relative-link form: **0 active stale links to old 4B paths remain.** No semantic/prose edits, no frozen records, no goldens, no do-not-touch paths touched.

**Corrected AC status:** AC3 PASS (after repair), AC7 PASS (after repair), AC8 PASS (portal index repaired). The "9/10 clean" summary above is superseded by this record. The narrow AC7 grep pattern (missing relative `../` link forms) is a cell-process lesson for 4C–4E.
