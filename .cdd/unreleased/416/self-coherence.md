# α self-coherence — Cycle 416

**Cycle:** [cnos#416](https://github.com/usurobor/cnos/issues/416) (Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404)).
**Branch:** `cycle/416` from `origin/main @ 92a7442d`.
**Mode:** γ+α+β collapsed on δ. This document is α's witnessed verification of the AC suite at α-close-out time.

## Deliverables filed

- **D1** — [`src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md) — 643 lines, verbatim doctrine transport from the 644-line source, with frontmatter changes (`parent: cdd` → `parent: handoff`; `requires:` reworded; `calls:` paths rewritten as absolute cross-package), Authority paragraph reworded to name cnos.handoff as canonical home, §2.3 STATUS-canonical-home declaration flipped from "canonical in `CDS.md`" → "canonical in this skill (§2.3)", §2.3.3 source-pin updated to cite §2.3.1 above instead of `CDS.md`, §2.5 Case (c) ε work product cite updated to absolute `cnos.cdd/skills/cdd/...` path, §2.9 ε epsilon SKILL cite updated, §5 Cross-references updated to absolute cnos.cdd paths + new "consumer" framing, §6 kata gamma SKILL cite updated to absolute path.
- **D2** — [`src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md) — 28-line compatibility pointer (≤ 50 lines). Frontmatter declares `artifact_class: pointer; status: moved; canonical: cnos.handoff/skills/handoff/cross-repo/SKILL.md`. Body lists the 7 doctrine elements covered by the new canonical home and notes that `cnos.cdd` is a consumer, not an owner. Backward-compatibility framing preserved.
- **D3** — [`src/packages/cnos.handoff/skills/handoff/HANDOFF.md`](../../../src/packages/cnos.handoff/skills/handoff/HANDOFF.md) — 62-line package contract (50–150 lines). Verbatim shape from cnos#416 issue body: preamble + package-contract section (owns / does NOT own / consumed by) + sub-surfaces (cross-repo landed v0.1; dispatch/mid-flight/artifact-channel/receipt-stream forthcoming) + related documents + non-goals. **Not** a synthesis doc. The old `.gitkeep` was removed.
- **D4** — Citation repair across cdd and cds (cdr verified no hits):
  - `cnos.cdd/skills/cdd/gamma/SKILL.md` — 2 sites: §"Cross-repo proposal intake" line 109 (`§2.3.3` cite + protocol-details cite); §"Cross-repo proposal close-out" line 318 (`§"STATUS state machine"` cite).
  - `cnos.cdd/skills/cdd/post-release/SKILL.md` — 1 site: §"Cross-repo trace" line 338 (2 cites in same paragraph).
  - `cnos.cdd/skills/cdd/epsilon/SKILL.md` — 2 sites: §"Cross-repo trace" line 99; §3 Cross-references line 137.
  - `cnos.cds/skills/cds/CDS.md` — 6 sites: line 814 (`§2.3.3` cite); line 954 (link block); lines 1506–1516 prose paragraph rewritten to reflect Sub 2 resolution; line 1540 (`§2.3` cite + "current home" prose rewritten); line 1547 (`§2.1` cite).
  - `cnos.cdd/skills/cdd/CDD.md` line 129 — left as-is (cites CDS.md as the coordination-surface aggregator, not cnos.cdd/cross-repo as canonical authority).
  - `cnos.cdd/skills/cdd/{operator,issue,activation}/SKILL.md` — left as-is (descriptive prose about `.cdd/iterations/cross-repo/` directory, not citations of `cdd/cross-repo/SKILL.md` as canonical).
  - `cnos.cds/docs/extraction-map.md` — left as-is (historical / pre-Sub-2 doc; out of scope per scaffold).
  - `cnos.cdr/skills/cdr/` — `rg` confirmed no cross-repo/SKILL.md citations; no edits required.
- **D5** — [`src/packages/cnos.handoff/docs/extraction-map.md`](../../../src/packages/cnos.handoff/docs/extraction-map.md) §1 Sub 2 rows — preamble updated to declare `Status: v0.1 migrated; canonical at cnos.handoff/skills/handoff/cross-repo/SKILL.md`; each of the 7 row notes carries the migrated status; the open-question paragraph was converted to "resolved open question" recording Q1's wholesale-move decision. Other sections (§2–§11) untouched.

## AC matrix — verification

### AC1: Required files exist — PASS

```
$ test -f src/packages/cnos.handoff/skills/handoff/HANDOFF.md && echo OK
OK
$ test -f src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md && echo OK
OK
$ test -f src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md && echo OK
OK
```

`.gitkeep` removed (HANDOFF.md fulfills its purpose).

### AC2: cnos.handoff cross-repo SKILL.md preserves doctrine — PASS

```
$ wc -l src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md
643 src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md      # ≥ 600 ✓
$ rg -c "cross-repo state machine|LINEAGE|STATUS|feedback patch|bundle archival|hat-collapse" \
   src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md
99                                                                     # ≥ 6 ✓
```

Per-keyword counts: `cross-repo state machine`=1, `LINEAGE`=40, `STATUS`=52, `feedback patch`=6, `bundle archival`=2, `hat-collapse`=11. All 6 keywords present (distinct ≥ 6 ✓).

### AC3: cnos.cdd cross-repo SKILL.md is a compatibility stub — PASS

```
$ wc -l src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md
28 src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md                # ≤ 50 ✓
$ rg -c "cnos.handoff|handoff/cross-repo|moved|canonical" \
   src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md
7                                                                      # ≥ 3 ✓
```

### AC4: STATUS-canonical-home declaration updated — PASS

```
$ rg "STATUS vocabulary is canonical|canonical home for the STATUS|STATUS vocabulary lives here|owns the STATUS vocabulary" \
   src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md
**The STATUS vocabulary is canonical in this skill (§2.3); CDS / CDR / CDD bind or consume it, but do not own it.**
```

1 hit. Declaration flipped from "canonical in `CDS.md`" → "canonical in this skill (§2.3)". ✓

### AC5: HANDOFF.md is minimal — PASS

```
$ wc -l src/packages/cnos.handoff/skills/handoff/HANDOFF.md
62 src/packages/cnos.handoff/skills/handoff/HANDOFF.md                 # 50–150 ✓
```

### AC6: No old-path-as-canonical citations remain — PASS

```
$ rg "cnos\.cdd/skills/cdd/cross-repo/SKILL\.md.*canonical|cdd/cross-repo.*canonical" \
   src/packages/cnos.cdd/skills/cdd/ src/packages/cnos.cds/skills/cds/ src/packages/cnos.cdr/skills/cdr/
(no matches)                                                           # 0 hits ✓
```

### AC7: Cross-references re-pointed in cdd/cds/cdr — PASS

```
$ rg -c "cnos.handoff/skills/handoff/cross-repo" \
   src/packages/cnos.cdd/ src/packages/cnos.cds/ src/packages/cnos.cdr/
src/packages/cnos.cds/skills/cds/CDS.md:6
src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md:1
src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md:2
src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md:2
src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md:3                 # ≥ 3 across packages ✓
```

11 distinct citations of the new canonical path in consumer files (cdd: 5 skill files cite; cds: 6 cites in CDS.md). Each citation that previously named the old cdd path as canonical authority now names the new cnos.handoff path. The stub's 3 self-references are inside the compatibility pointer itself.

### AC8: No behavioral redesign — PASS

Structural elements preserved verbatim (spot-check on the moved file):

1. **8 STATUS events listed** — `drafted`, `submitted`, `accepted`, `modified`, `rejected`, `landed`, `withdrawn`, `revised` all present with verbatim definitions (table at §2.3.1) and verbatim transition graph (§2.3.2). ✓
2. **4 directional cases with subcases** — §2.1 Case (a) with sub-shapes a.1 + a.2; Case (b) with b.1; Case (c); Case (d) with d.1 + d.2 + d.3. Verbatim definitions and empirical anchors. ✓
3. **LINEAGE.md schema sections per case** — §2.6 has 4 case-specific subsection lists (a / b / c / d), each with its required sections preserved verbatim. ✓
4. **Feedback-patch header form** — §2.7 header template (`From:` / `Date:` / `Subject:` / context + `---` / `<unified diff>`) preserved verbatim; `Apply command convention` + `Patch filename convention` + `Content rules` blocks preserved. ✓
5. **Hat-collapse attribution rules** — §2.9 carries the two-places rule (bundle LINEAGE.md `## ε actor` + cdd-iteration.md `## Hat-collapse acknowledgment`), the boundary-preservation paragraph (`ROLES.md §4a` persona/protocol/project boundary), and Rule 2.9.a verbatim. Only edit: §2.9's `epsilon/SKILL.md §2` cite became `cnos.cdd/skills/cdd/epsilon/SKILL.md §2` (absolute cross-package form). ✓
6. **Known protocol edge cases** — §2.10 carries all 6 verbatim edge-case bullets with their empirical anchors: cn-sigma/agent-activate-skill direct-acceptance; cph/bootstrap-cdr repo rename; cn-sigma post-filing refinement; cn-sigma asymmetric source posture; cph/coherence-drift-sweep-followup case-c dual-purpose; cph/issue-32-tightening case d.3; cn-rho/bootstrap-2026-05-19 case d.1. ✓

≥ 5 of 6 spot-checked elements pass (in fact all 6). ✓

### AC9: No new schemas, no cdd-verify changes, no runtime/harness changes — PASS

```
$ test ! -d schemas/handoff && echo OK
OK
$ test ! -d schemas/ccnf-o && echo OK
OK
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/commands/cdd-verify/ | wc -l
0
$ git diff origin/main..HEAD -- src/go/ | wc -l
0
$ git diff origin/main..HEAD -- scripts/release.sh | wc -l
0
```

All four invariants hold. ✓

### AC10: Extraction-map status update — PASS

```
$ rg "v0.1 migrated.*cross-repo" src/packages/cnos.handoff/docs/extraction-map.md
**Status: v0.1 migrated; canonical at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`.** Landed via [cnos#416] ...
| ... | Sub 2 | **Status: v0.1 migrated; canonical at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`.** Verbatim move; ...
| ... | Sub 2 | **Status: v0.1 migrated; canonical at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`.** Inside the wholesale move; ...
(... 7 row hits + 1 preamble hit) ✓
```

§1 preamble + all 7 row notes carry the `v0.1 migrated` Status; other sections (§2–§11) untouched. ✓

## Out-of-scope changes (none)

- No `schemas/handoff/` or `schemas/ccnf-o/` directories created.
- No `cn cdd verify` changes.
- No `src/go/` or `scripts/release.sh` changes.
- No CCNF-O work.
- Old `cnos.cdd/skills/cdd/cross-repo/SKILL.md` not deleted (compatibility stub remains).
- No behavioral redesign — STATUS vocabulary, directional cases (a / b / c / d.1 / d.2 / d.3), LINEAGE schemas, feedback-patch format, archival rule, hat-collapse attribution all transport unchanged.

## Mid-cycle clarifications

None. The dispatch metadata pinned every axis; no `gamma-clarification.md` was needed.

## Verdict

AC1–AC10 all PASS. α attests Sub 2 of #404 is mechanically complete and ready for β-review.
