# γ-scaffold — Cycle 417

**Cycle:** [cnos#417](https://github.com/usurobor/cnos/issues/417) — Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404)
**Branch:** `cycle/417` from `origin/main @ f2430329` (Sub 2 / cnos#416 landed).
**Mode:** γ+α+β collapsed on δ; docs-only / package migration (no behavioral redesign).
**Date:** 2026-05-22.

## Scope ruling (three-ruling pattern, same as Sub 2 / cnos#416)

1. **Move dispatch + implementation-contract + δ-inward-membrane doctrine wholesale into cnos.handoff.** Verbatim section extraction from three source surfaces, synthesized into one cohesive `cnos.handoff/skills/handoff/dispatch/SKILL.md`. Doctrine changes path, not substance.
2. **Replace source sections with compatibility pointers/stubs** — each of the three source sections (`gamma/SKILL.md §2.5` dispatch-prompts + implementation-contract block; `operator/SKILL.md §3a`; `delta/SKILL.md §2`) becomes a short pointer paragraph + 1–2 sentences of role-local framing.
3. **Update cdd/cds/cdr cross-references** so consumers cite `cnos.handoff/skills/handoff/dispatch/SKILL.md` as canonical authority for the dispatch-prompt template, the 7-axis implementation-contract schema, and δ-as-inward-membrane doctrine.

Plus: update `cnos.handoff/skills/handoff/HANDOFF.md` (single-row edit; "dispatch" → Landed) and `cnos.handoff/docs/extraction-map.md` (Sub 3 rows in §2 + §3 marked migrated).

## Surfaces γ expects α to touch

| Surface | Action |
|---|---|
| `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` | **author** (new file; synthesis of three source sections) |
| `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | §2.5 dispatch-prompts + implementation-contract block: rewrite as pointer (≤ 30 lines for the pointer block; rest of §2.5 — branch pre-flight, scaffold gate, polling cross-ref, spec-staleness, Step 5 unblock — preserved as cdd-resident cycle-lifecycle doctrine) |
| `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` | §3a: rewrite as pointer (already a redirect to delta §2 post-#397; now re-point to cnos.handoff) |
| `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | §2 + §2.1 + §2.3: rewrite as pointer (canonical doctrine moves to cnos.handoff; §2.2 Phase 4a landing note + role-local framing retained) |
| `src/packages/cnos.handoff/skills/handoff/HANDOFF.md` | move "dispatch" from Forthcoming → Landed (one-row edit) |
| `src/packages/cnos.handoff/docs/extraction-map.md` | mark Sub 3 rows in §2 + §3 as `**Status: v0.1 migrated; canonical at cnos.handoff/skills/handoff/dispatch/SKILL.md**` |
| Consumer cross-references | re-point cdd/beta/SKILL.md Rule 7, cdd/alpha/SKILL.md §3.6, cdd/delta/SKILL.md `inputs:` frontmatter, cdd/CDS.md (cds) cite (where canonical) to cnos.handoff/skills/handoff/dispatch/SKILL.md |

## AC oracle approach

11 ACs are largely mechanical (file/section line counts, grep counts, negative greps). γ runs the AC matrix in `self-coherence.md` after authoring.

- AC1: `test -f` on the new SKILL.md + parse-check the three source files
- AC2: ≥ 5 keyword hits (dispatch-prompt, implementation contract, 7 axes, inward membrane, gamma-clarification)
- AC3: section line counts ≤ 30 lines + ≥ 3 new-canonical cites
- AC4: 7 distinct axis names present in dispatch/SKILL.md
- AC5: ≥ 3 δ-inward-membrane authority hits
- AC6: dispatch row in HANDOFF.md under Landed (not Forthcoming)
- AC7: source role-skills ≥ 80% pre-migration line count
- AC8: negative grep — no `gamma/SKILL.md §2.5.*canonical|operator/SKILL.md §3a.*canonical|delta/SKILL.md §2.*canonical` hits
- AC9: cnos.cdr untouched (per Sub 2 verification — confirm by `git diff origin/main..HEAD -- src/packages/cnos.cdr/`)
- AC10: no schemas/handoff, no schemas/ccnf-o, no cdd-verify changes
- AC11: extraction-map.md Sub 3 rows mark migrated

## Empirical anchor (precedent)

[cnos#416](https://github.com/usurobor/cnos/issues/416) — Sub 2 / cross-repo wholesale move. Same three-ruling pattern; same δ+γ+α+β collapse pattern; differs in shape (Sub 2 was a single-file wholesale move with a compatibility pointer at the old full-file path; Sub 3 is a three-source synthesis with three section-level pointers in mid-file sections).

The Sub 2 precedent's `cnos.handoff/skills/handoff/cross-repo/SKILL.md` is the **structural precedent** for the new `dispatch/SKILL.md` frontmatter shape.

## Expected diff scope

- 1 new file: `cnos.handoff/skills/handoff/dispatch/SKILL.md` (~300–500 lines synthesizing 3 source sections)
- 3 section rewrites in cdd role skills (gamma/operator/delta SKILL.md; each section ≤ 30 lines post-rewrite; rest of each file preserved)
- 1 single-row edit in `cnos.handoff/skills/handoff/HANDOFF.md`
- ~2 row updates in `cnos.handoff/docs/extraction-map.md`
- Cross-ref repairs in consumers (cnos.cdd beta/SKILL.md, alpha/SKILL.md, cdd/delta/SKILL.md frontmatter, cnos.cds/CDS.md citation lines)
- 0 changes in cnos.cdr (per Sub 2 verification; confirm at AC time)
- 0 schema / runtime / cdd-verify / src/go changes

## Mode declaration

γ+α+β collapsed on δ (single-session). Commits will be authored under `γ-417:`, `α-417:`, `β-417:` per Sub 2 precedent. γ does not self-merge — operator merges with `Closes #417`.
