# γ-scaffold — Cycle 418

**Cycle:** [cnos#418](https://github.com/usurobor/cnos/issues/418) — Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404)
**Branch:** `cycle/418` from `origin/main @ 0bbd312b` (Sub 3 / cnos#417 landed).
**Mode:** γ+α+β collapsed on δ; docs-only / package migration (no behavioral redesign).
**Date:** 2026-05-22.

## Scope ruling (three-ruling pattern, same as Subs 2 / 3 — cnos#416 / cnos#417)

1. **Move two parallel surfaces into cnos.handoff in one cycle.** Two new sub-skills authored from empirical practice (mid-flight rescue, anchored at cnos#391) and from extraction of wire-format invariants resident at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` and `§"Artifact contract"`. The two surfaces share consumers (cdd/alpha, cdd/beta, cdd/gamma, dispatch/SKILL.md) and HANDOFF.md update; collapsing them is the natural cohesion.
2. **Update source pointers** — `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.10` (issue-edit cache-bust now cites mid-flight as canonical); `cnos.handoff/skills/handoff/dispatch/SKILL.md` (drop "pending Sub 4" / "forthcoming" qualifiers; cite Sub 4 as landed); `cnos.handoff/skills/handoff/SKILL.md` loader (v0.1 caveat update).
3. **Update consumer cross-references** to cite the new canonical paths — cdd/alpha/SKILL.md, cdd/beta/SKILL.md, cdd/gamma/SKILL.md (issue-edit cache-bust + spec-staleness + artifact channel references), cnos.cds/CDS.md (mid-flight clarification section now points back to mid-flight/SKILL.md as wire-format canonical), README files.

Plus: update `cnos.handoff/skills/handoff/HANDOFF.md` (two-row edit; "mid-flight" + "artifact-channel" → Landed) and `cnos.handoff/docs/extraction-map.md` (§4, §5, §8 rows marked migrated).

## Open-question rulings (extraction-map.md Q3 + §4 + §5)

- **Q3 (Sub 4 unification):** **split** into two skills per recommendation. Different rates of fire (mid-flight is asynchronous γ-to-in-flight-α; artifact-channel is sequential α→β→γ); per-mechanism cohesion is high; combined file would conflate two protocols.
- **§4 row 3 (cnos.cds mid-flight clarification):** option **(b)** — keep canonical in cnos.cds for the operational realization (the `gh issue edit` + `commit` + `push` sequence γ runs at clarification time); add `cnos.handoff/skills/handoff/mid-flight/SKILL.md` as the wire-format canonical (file path, authoring role, reader role, trigger conditions, cache-bust function). cnos.cds becomes a pointer for the wire-format invariant; the operational overlay stays. The empirical anchor (cnos#391) is the rescue mechanism; the cache-bust function is the wire-format invariant.
- **§5 rows (Location matrix / Ownership matrix / Ordered flow / Frozen snapshot rule):** mixed per row — Location matrix + Ownership matrix + Ordered flow stay canonical in `cnos.cds/skills/cds/CDS.md §"Artifact contract"` (consumer-protocol-specific; CDR has different artifacts); the wire-format **invariant** (sequential α→β→γ over `.cdd/unreleased/{N}/`; per-role write ownership pattern; frozen-snapshot rule) is canonical in `cnos.handoff/skills/handoff/artifact-channel/SKILL.md`. The new handoff skill cites cnos.cds for the CDS-specific instantiation; cnos.cds gets a pointer line acknowledging the new canonical home for the wire-format invariant. The frozen-snapshot rule moves verbatim (wire-format), with a cite back from cds.
- **§4 + §5 + §8 folded:** the spec-staleness propagation invariant (γ writes a coordination note when spec changes mid-flight) is folded into mid-flight/SKILL.md as recommended.

## Surfaces γ expects α to touch

| Surface | Action |
|---|---|
| `src/packages/cnos.handoff/skills/handoff/mid-flight/SKILL.md` | **author** (new file; synthesis of empirical anchor cnos#391 + cds §"Mid-flight clarification" wire-format invariant + cdd/gamma §"Spec-staleness propagation" wire-format invariant) |
| `src/packages/cnos.handoff/skills/handoff/artifact-channel/SKILL.md` | **author** (new file; synthesis of cds §"Artifact contract" wire-format invariants — sequential rule, per-role write ownership pattern, frozen-snapshot rule — with pointers back to cds for CDS-specific instantiation) |
| `src/packages/cnos.handoff/skills/handoff/HANDOFF.md` | move "mid-flight" + "artifact-channel" from Forthcoming → Landed (two-row edit + non-goals row update) |
| `src/packages/cnos.handoff/skills/handoff/SKILL.md` | update v0.1 caveat: mid-flight + artifact-channel now landed under cnos#418; receipt-stream remains pending (Sub 5) |
| `src/packages/cnos.handoff/skills/handoff/cross-repo/SKILL.md` | §2.10 update — issue-edit cache-bust mention now cites mid-flight/SKILL.md as canonical |
| `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` | drop "pending Sub 4" / "forthcoming" qualifiers across §"Scope" non-goals + §"Spec-staleness propagation" + §3.3 + §"Related documents" + §"Non-goals"; cite mid-flight + artifact-channel as landed |
| `src/packages/cnos.handoff/docs/extraction-map.md` | mark Sub 4 rows in §4 + §5 + §8 as `**Status: v0.1 migrated; canonical at cnos.handoff/skills/handoff/{mid-flight,artifact-channel}/SKILL.md**` |
| `src/packages/cnos.handoff/README.md` | update Sub 4 status (mid-flight + artifact-channel land) |
| `src/packages/cnos.cds/skills/cds/CDS.md` | §"Mid-flight clarification" — append pointer line citing mid-flight/SKILL.md as wire-format canonical; §"Artifact contract" — append pointer line citing artifact-channel/SKILL.md as wire-format invariant home |
| `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | §2.5 → "Spec-staleness propagation" — add pointer to mid-flight/SKILL.md for the wire-format invariant (the consumer-specific file list stays cdd-side per extraction-map.md §8); §2.5 → "Issue-edit cache-bust" line — cite mid-flight/SKILL.md as canonical |
| `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | issue-edit cache-bust mention now cites mid-flight/SKILL.md as canonical |
| `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` | issue-edit cache-bust mention now cites mid-flight/SKILL.md as canonical |

## AC oracle approach

11 ACs from cnos#418 issue body are mechanical (file existence / line counts / grep counts / negative greps). γ runs the AC matrix in `self-coherence.md` after authoring.

- AC1: `test -f` on the two new SKILL.md files + parse-check
- AC2: ≥ 5 mid-flight keyword hits; line count 100–350
- AC3: ≥ 7 artifact-channel keyword hits; line count 150–400
- AC4: HANDOFF.md "mid-flight" + "artifact-channel" rows under Landed (not Forthcoming)
- AC5: ≥ 3 new-canonical cites in cdd consumers (alpha/beta/gamma referring to mid-flight/SKILL.md)
- AC6: cnos.cdr untouched (`git diff origin/main..HEAD -- src/packages/cnos.cdr/` empty)
- AC7: cross-repo §2.10 cites mid-flight as canonical for issue-edit cache-bust
- AC8: dispatch/SKILL.md drops "pending Sub 4" / "forthcoming" qualifiers; cites Sub 4 as landed
- AC9: CDD.md kernel sections byte-identical; ≤ 5 lines insertion total (D6 likely skipped — artifact channel discoverability runs through CDS.md pointer list already)
- AC10: no schemas/handoff, no schemas/ccnf-o, no cdd-verify changes
- AC11: rescue mechanism + artifact-file names spot-check unchanged (gamma-clarification.md path; self-coherence.md / alpha-closeout.md / beta-review.md / beta-closeout.md / gamma-closeout.md / gamma-scaffold.md / cdd-iteration.md filenames preserved verbatim)

## Empirical anchor (precedent)

Sub 2 / [cnos#416](https://github.com/usurobor/cnos/issues/416) (cross-repo wholesale move; same δ+γ+α+β collapse pattern; single-file output). Sub 3 / [cnos#417](https://github.com/usurobor/cnos/issues/417) (dispatch synthesis from three source sections; same collapse pattern). Sub 4 differs in shape: two parallel surfaces in one cycle (mid-flight + artifact-channel share consumers + close-out overhead). The precedent's `cross-repo/SKILL.md` (Sub 2) and `dispatch/SKILL.md` (Sub 3) are the **structural precedents** for the two new sub-skill files' frontmatter shape and section ordering.

The **rescue mechanism** itself anchors on [cnos#391](https://github.com/usurobor/cnos/issues/391) — the cycle where γ recovered from a wrong-shape α implementation by editing the issue body to pin missing axes and committing `gamma-clarification.md` to the cycle branch; α picked up the cache-bust via cycle-branch polling and re-shaped the implementation. The mechanism crystallized as a γ→in-flight-α channel here.

## Expected diff scope

- 2 new files: `cnos.handoff/skills/handoff/mid-flight/SKILL.md` (~200 lines) + `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` (~250 lines)
- 1 two-row edit in `cnos.handoff/skills/handoff/HANDOFF.md`
- ~4 row updates in `cnos.handoff/docs/extraction-map.md` (§4 + §5 + §8 marked migrated)
- ~5–10 line edit in `cnos.handoff/skills/handoff/SKILL.md` (loader v0.1 caveat update)
- ~3 line edit in `cnos.handoff/README.md`
- ~5–10 line edit in `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (§2.10 issue-edit cache-bust pointer)
- ~10–20 line edits across `cnos.handoff/skills/handoff/dispatch/SKILL.md` (drop "pending" qualifiers; cite Sub 4 as landed)
- ~2 pointer-line additions in `cnos.cds/skills/cds/CDS.md` (§"Mid-flight clarification" + §"Artifact contract")
- ~3 small edits in `cnos.cdd/skills/cdd/{alpha,beta,gamma}/SKILL.md` (issue-edit cache-bust cite + gamma spec-staleness pointer)
- 0 changes in cnos.cdr (per Sub 2 / Sub 3 verification; confirm at AC time)
- 0 schema / runtime / cdd-verify / src/go changes
- 0 changes in `cnos.cdd/skills/cdd/CDD.md` (artifact-channel discoverability already runs through CDS.md pointer list; no kernel edit needed)

## Mode declaration

γ+α+β collapsed on δ (single-session). Commits will be authored under `γ-418:`, `α-418:`, `β-418:` per Sub 2 / Sub 3 precedent. γ does not self-merge — operator merges with `Closes #418`.
