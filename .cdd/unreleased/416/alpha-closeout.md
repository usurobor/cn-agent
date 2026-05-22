# α close-out — Cycle 416

**Cycle:** [cnos#416](https://github.com/usurobor/cnos/issues/416) — Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Branch:** `cycle/416` from `origin/main @ 92a7442d`.
**Mode:** γ+α+β collapsed on δ.

## Disposition

α attests cnos#416 deliverables D1–D5 are filed; AC1–AC10 PASS per `self-coherence.md`; handoff to β for adversary review.

## Files touched

```
new file:   .cdd/unreleased/416/gamma-scaffold.md           (γ commit)
new file:   .cdd/unreleased/416/self-coherence.md           (α commit)
new file:   .cdd/unreleased/416/alpha-closeout.md           (α commit; this file)
new file:   src/packages/cnos.handoff/skills/handoff/HANDOFF.md            (D3; 62 lines)
new file:   src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md   (D1; 643 lines)
deleted:    src/packages/cnos.handoff/skills/handoff/.gitkeep              (replaced by HANDOFF.md per AC1)
modified:   src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md           (D2; 644 → 28 lines, compatibility stub)
modified:   src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md                (D4; 2 citation re-points)
modified:   src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md         (D4; 2 citation re-points)
modified:   src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md              (D4; 2 citation re-points)
modified:   src/packages/cnos.cds/skills/cds/CDS.md                        (D4; 6 citation re-points / prose updates)
modified:   src/packages/cnos.handoff/docs/extraction-map.md               (D5; Sub 2 rows updated)
```

## AC verification — α witness

| AC | Status | Evidence |
|---|---|---|
| AC1 (3 files exist) | PASS | `test -f` on all 3 paths returns OK |
| AC2 (≥ 600 lines + ≥ 6 keywords) | PASS | 643 lines; 6/6 keyword categories present (cross-repo state machine=1, LINEAGE=40, STATUS=52, feedback patch=6, bundle archival=2, hat-collapse=11) |
| AC3 (stub ≤ 50 lines + ≥ 3 hits) | PASS | 28 lines; 7 hits on cnos.handoff/handoff/cross-repo/moved/canonical |
| AC4 (STATUS-canonical-home flipped) | PASS | "The STATUS vocabulary is canonical in this skill (§2.3); CDS / CDR / CDD bind or consume it, but do not own it." (§2.3 of new SKILL.md) |
| AC5 (HANDOFF.md 50–150 lines) | PASS | 62 lines |
| AC6 (no old-path-as-canonical) | PASS | `rg "cnos\.cdd/skills/cdd/cross-repo/SKILL\.md.*canonical\|cdd/cross-repo.*canonical"` returns 0 hits |
| AC7 (≥ 3 new-path cites in consumers) | PASS | 11 cites across cdd (5 files) + cds (CDS.md) |
| AC8 (≥ 5 structural elements preserved verbatim) | PASS | 6/6 spot-checked: 8 STATUS events; 4 directional cases; LINEAGE schemas per case; feedback-patch header; hat-collapse rules; protocol edge cases |
| AC9 (no schemas / cdd-verify / src/go / scripts/release.sh) | PASS | `test ! -d` + `git diff` all clean |
| AC10 (extraction-map Sub 2 updated) | PASS | §1 preamble + 7 row notes carry `Status: v0.1 migrated` |

## Hard-rule attestation

- ✅ **No behavioral redesign.** STATUS vocabulary, directional cases, LINEAGE schemas, feedback-patch format, archival rule, hat-collapse — all transport unchanged. Only edits to the moved file are: frontmatter (`parent`, `requires`, `calls`); Authority paragraph; §2.3 STATUS-canonical-home declaration; §2.3.3 source-pin cite update; §2.5 / §2.9 / §5 / §6 cross-package cite updates (all expanding relative paths to absolute).
- ✅ **No `schemas/handoff/` or `schemas/ccnf-o/`.** No CUE schemas authored.
- ✅ **No `cn cdd verify` changes.** `src/packages/cnos.cdd/commands/cdd-verify/` untouched.
- ✅ **No runtime/harness changes.** `src/go/`, harness scripts, `scripts/release.sh` untouched.
- ✅ **No CCNF-O work.**
- ✅ **Old cdd/cross-repo file not deleted.** Replaced with compatibility stub.

## Open observations for β

- Per scaffold, citations in cnos.cdd/skills/cdd/CDD.md, operator/SKILL.md, issue/SKILL.md, activation/SKILL.md that mention "cross-repo" were left as-is because they describe `.cdd/iterations/cross-repo/` directory shape (not the doctrine skill) or coordinate at a higher abstraction (CDS.md as aggregator). β should sanity-check this scope-of-sweep decision.
- The cnos.cds extraction-map (`cnos.cds/docs/extraction-map.md`) still contains 4 references framing `cross-repo/SKILL.md` location as "open". This is a historical / pre-Sub-2 document outside D4's scope ("D4 targets cdd/cds/cdr skill files"); the scaffold flagged this explicitly. Operator decides whether a future cycle sweeps that doc.
- The cnos.handoff loader (`skills/handoff/SKILL.md`) still describes HANDOFF.md and sub-skills as "forthcoming" in places — that v0.1-caveat framing is now partially out of date (cross-repo and HANDOFF.md landed). A future Sub 3+ cycle should tighten that prose; not in Sub 2 scope.

## Handoff

β: read `self-coherence.md` + this file; run the AC suite independently as adversary; file `beta-review.md` with verdict.
