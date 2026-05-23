# α self-coherence — cycle/412

**Cycle:** 412
**Issue:** [cnos#412](https://github.com/usurobor/cnos/issues/412)
**Author:** α (collapsed with γ+β on δ)
**Date:** 2026-05-22

## Pre-flight

- Branch `cycle/412` checked out at `71b25672` (head of `origin/main` post-cycle/410).
- Sibling sub cnos#411 (Sub 6) is on disjoint files (`CDD.md` cleanup) — no conflict.
- Hard rules pinned: cnos.cdd, cnos.cdr, cnos.cds/skills/cds/CDS.md all read-only for this cycle.
- Structural precedent: `src/packages/cnos.cdr/docs/empirical-anchor-cph.md` read end-to-end.
- CDS.md surface enumeration read (top-level `^## ` headers: 19 sections, of which surfaces 4–16 are the load-bearing operational surfaces this doc maps cycles to).
- INDEX.md read (32 rows as of cycle/410); representative-cycle close-out paths spot-checked under `.cdd/unreleased/`.

## Design decisions

1. **Mapping unit = CDS surface, not cycle.** The cph anchor's unit is "cph artifact class → CDR field." The CDS anchor inverts: "CDS surface → cycle table." This matches the issue body's "Cycle-to-surface mapping" framing and produces tables a reader can pick by surface, not by cycle.

2. **Table columns.** `Cycle | Issue | Date | What it exercised | Artifact citation` per the issue body's example shape. Each row carries a one-line "What it exercised" and a path-or-SHA citation under "Artifact citation."

3. **Twelve surfaces mapped.** §Six-field, §Selection function, §Development lifecycle, §Coordination surfaces, §Artifact contract, §Mechanical vs judgment, §Review CLP, §Gate, §Assessment, §Closure, §Retro-packaging, §Large-file authoring rule. Doctrinal surfaces (§Purpose, §Architecture choice, §Persona/Protocol/Project, §Empirical anchor, §Related documents, §Non-goals) noted but not given separate cycle tables — they are exercised at the document layer.

4. **Cycle reuse across surfaces.** A cycle (e.g. #402, #392, #393) appears under multiple surfaces. This is faithful to the empirical claim: cycles are compositional, not single-axis.

5. **Pin commit.** Anchored against cnos `main` at `71b25672` (the head when cycle/412 was filed). Mirrors the cph anchor's pin discipline (`2bf73ad6`).

6. **AC4 generous coverage.** 34 distinct cycle numbers cited (issue requires ≥ 20). Cycles span #331 through #412; includes the parent (#403), the sibling (#411), and this cycle (#412) where they appear as artifact citations in their own right.

7. **No new content beyond synthesis.** All cited cycles already exist on `origin/main` with their commits and close-out artifacts. The doc synthesizes the mapping; it does not author new doctrine or new cycle history.

## Hard-rule preservation

- `git diff origin/main..HEAD -- src/packages/cnos.cdd/` → 0 lines (verified).
- `git diff origin/main..HEAD -- src/packages/cnos.cdr/` → 0 lines (verified).
- `git diff origin/main..HEAD -- src/packages/cnos.cds/skills/cds/CDS.md` → 0 lines (verified).
- Only files added: `src/packages/cnos.cds/docs/empirical-anchor-cdd.md` + close-out artifacts under `.cdd/unreleased/412/`.

## Mechanical AC pre-check (α-side)

- AC1: file at canonical path, 355 lines ≥ 100. PASS.
- AC2: `grep -c "^## " ...` = 6 ≥ 5. PASS.
- AC3: 12 surface mapping tables (`###` with mapping table) ≥ 5. PASS.
- AC4: 34 distinct cycle numbers cited ≥ 20. PASS.
- AC5/6/7: hard-rule diffs all zero. PASS.
- AC8: §Related documents cites CDS.md, CDD.md, cph empirical anchor, INDEX.md. PASS.

## Known limits

- The mapping is **representative**, not exhaustive. New cycles will continue to exercise these surfaces; this doc does not claim to enumerate them all.
- Some date columns inherit the issue-open date rather than the merge date (deliberately, to make the cycle's "filing context" legible).
- A few citation paths are conventional (e.g. `.cdd/releases/docs/.../<N>/`) rather than explicit per-cycle paths; the INDEX aggregator is the authoritative resolver.
- The "Closure cycles and INDEX rows" section compresses ~32 INDEX rows into a 14-row significance table. This is a curation choice, not a claim of completeness.

α hands off to β for review.
