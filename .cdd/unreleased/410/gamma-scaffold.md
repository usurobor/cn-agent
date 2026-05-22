# γ Scaffold — cycle/410 (Sub 5 of cnos#403 wave)

**Issue:** [cnos#410](https://github.com/usurobor/cnos/issues/410)
**Branch:** `cycle/410` (from `origin/main` at `4a87cdf9`)
**Mode:** substantial (design-and-build; B-lite migration)
**Wave:** cnos#403 sub-wave (Sub 5 of 7; final canonical-migration sub before Sub 6 cleanup)
**Dispatch shape:** **β-α-collapse-on-δ** — γ+α+β collapsed on a single agent under operator δ supervision. Justified per CDS Field 6 actor collapse rule (B-lite migration / contract-authoring class; no novel executable surface). Configuration-floor cap applies: γ-axis ≤ A-, β-axis ≤ A- per `cnos.cdd/skills/cdd/release/SKILL.md §3.8` configuration-floor clause.

## Source contract

cnos#410 issue body is the authoritative scope. 17 ACs (AC1–AC17). 10 deliverables (D1–D10). 8 new top-level sections in CDS.md, plus extraction-map.md Status updates for rows 5–12, plus optional thin overlays at `skills/cds/{review,gate,assessment}/`.

## Surfaces α expects to touch

1. `src/packages/cnos.cds/skills/cds/CDS.md` — **edited**. 8 new top-level sections inserted between `## Artifact contract` (ends line ~1951) and `## Empirical anchor` (line 1955). Manifest header (lines 1–2) updated to reflect the 8 new sections.
2. `src/packages/cnos.cds/docs/extraction-map.md` — **edited**. Status column updated for tables §5–§12 (rows 5–12) to `**v0.1 migrated; canonical at CDS.md §{section name}**`. Rows 1–4 untouched.
3. Optional thin overlays — `src/packages/cnos.cds/skills/cds/{review,gate,assessment}/SKILL.md`. Deferred decision: skip unless a clean reason to add. Decision recorded at α self-coherence time.
4. `.cdd/unreleased/410/` — γ scaffold (this file), α self-coherence, β review, α/β/γ close-outs, cdd-iteration stub if applicable.
5. `.cdd/iterations/INDEX.md` — appended row for cycle 410.

## Surfaces α MUST NOT touch (hard rules from #410)

- `src/packages/cnos.cdd/` — UNTOUCHED. `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns 0 lines (AC14). Sub 6 sweeps CDD.md pending-cds markers, not this sub.
- `src/packages/cnos.cdr/` — UNTOUCHED. `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns 0 lines (AC15).
- `skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/` — NO new files (AC17 — no deep role rewrites).
- `schemas/` — out of scope.

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None new |
| Package scoping | `src/packages/cnos.cds/skills/cds/CDS.md` + `src/packages/cnos.cds/docs/extraction-map.md` + optional `skills/cds/{review,gate,assessment}/SKILL.md` (≤ 40 lines each) |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | hard rules above; F1–F10 anchors preserved; §9.1 trigger phrases preserved; 5 software-cycle non-goals preserved verbatim |

## AC oracle approach

All 17 ACs are mechanical (grep / file existence / `git diff` empty) or read-checks:

- **AC1–AC8:** `grep -c "^## {section-name}" CDS.md` returns 1 for each of the 8 sections.
- **AC9:** Per-section sub-heading coverage via `grep -E "^### "`.
- **AC10 (F1–F10):** All 10 F-anchors appear as identifiable items in §Gate § Closure verification checklist.
- **AC11 (§9.1 triggers):** All 4 trigger phrases appear in §Assessment § Cycle iteration triggers.
- **AC12 (non-goals):** 5 software-cycle non-goals named verbatim.
- **AC13:** Each new top-level section (D1–D6, D8) ends with `### Operational realization` citing ≥1 cdd skill file.
- **AC14:** `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns 0 lines.
- **AC15:** `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns 0 lines.
- **AC16:** extraction-map.md Status for rows 5–12 updated; rows 1–4 unchanged.
- **AC17:** No new files under `skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`.

## F1–F10 anchor design (canonical commitment for AC10)

§Gate § Closure verification checklist names 10 failure-mode anchors (F1–F10). These are stable cross-reference targets Sub 6 re-points. The 10 anchors codify the closure-gate failure modes named in the issue body's parenthetical (the §Gate source content) and in `gamma/SKILL.md §2.10`'s 14-row closure gate. Mapping:

| Anchor | Failure mode caught | Source (cdd v0.1 overlay) |
|---|---|---|
| F1 | Missing α close-out (`alpha-closeout.md` absent at γ closure) | `gamma/SKILL.md §2.10` row 1 |
| F2 | Missing β close-out (`beta-closeout.md` absent at γ closure) | `gamma/SKILL.md §2.10` row 2 |
| F3 | Missing γ close-out (`gamma-closeout.md` absent — closure not declared) | `gamma/SKILL.md §2.10` (closure-declaration commit) |
| F4 | Stale `.cdd/unreleased/{N}/` directory after release (cycle-directory move skipped) | `gamma/SKILL.md §2.10` row 12; `release/SKILL.md §2.5a` |
| F5 | Missing `RELEASE.md` at release commit (CI auto-generates sparse notes) | `gamma/SKILL.md §2.10` row 11; `release/SKILL.md §3.7` |
| F6 | δ tag ordering violation (tag pushed before γ closure declaration on main) | `operator/SKILL.md §3.4` "Do not tag/release before `gamma-closeout.md` exists on main"; `release/SKILL.md §2.6` |
| F7 | Missing α close-out re-dispatch (γ failed to request δ re-dispatch of α for post-merge close-out) | `gamma/SKILL.md §2.7`; `cnos.cdd/skills/cdd/CDD.md §1.6a` |
| F8 | Missing PRA (`POST-RELEASE-ASSESSMENT.md` absent) | `gamma/SKILL.md §2.10` row 3; `post-release/SKILL.md` |
| F9 | Missing `cdd-iteration.md` when triggers fired (`protocol_gap_count > 0` but no iteration artifact) | `gamma/SKILL.md §2.10` row 14; `post-release/SKILL.md §5.6b` |
| F10 | Unresolved triage (cycle closes with `Noted` / no disposition on a finding) | `gamma/SKILL.md §3.7` "Do not close the cycle with unresolved triage" |

The F1–F10 ordering matches the issue body's parenthetical sequence (α close-out → β close-out → γ close-out → stale dir → missing RELEASE.md → tag ordering → re-dispatch → PRA → iteration → triage). Sub 6 will re-point any `gamma/SKILL.md §2.10` / `release/SKILL.md` citations of the F-anchors at `CDS.md §Gate § Closure verification checklist § F{N}`.

## §9.1 cycle iteration trigger anchor design (canonical commitment for AC11)

§Assessment § Cycle iteration triggers names the 4 §9.1 triggers verbatim from `post-release/SKILL.md §4b` and `gamma/SKILL.md §2.8`:

1. Review churn — review rounds > 2
2. Mechanical overload — mechanical ratio > 20% (with ≥10 findings)
3. Avoidable tooling / environmental failure
4. Loaded skill failed to prevent a finding

The four-trigger anchor preserves §9.1 compatibility for citing skill files. Sub 6 re-points any `CDD.md §9.1` citations at `CDS.md §Assessment § Cycle iteration triggers`.

## §Non-goals integration design (D7)

CDS.md already has `## Non-goals` (line 2175) carrying **sub-level non-goals** (the scope-discipline of Sub 2's authoring cycle: "do not modify CDD.md", "do not author role overlays", etc.). The 5 software-cycle non-goals from the issue body are a different class — they apply to **every CDS cycle** as doctrinal non-goals, not just to this sub.

**Integration:** rename the existing section structure to split:

- `## Non-goals` becomes the umbrella heading
- `### Sub-level non-goals (Subs 2–5 authoring scope discipline)` — the existing bullet list
- `### Software-cycle non-goals` — the 5 new items, verbatim from the issue body

Both are discoverable under one top-level anchor; the split is content-grouped, not new-section.

## Section ordering inside CDS.md (insertion plan)

Insert in this order after line 1953 (`---` separator) and before line 1955 (`## Empirical anchor`):

1. `## Mechanical vs judgment` (D1)
2. `## Review CLP` (D2)
3. `## Gate` (D3)
4. `## Assessment` (D4)
5. `## Closure` (D5)
6. `## Retro-packaging` (D6)
7. `## Large-file authoring rule` (D8)

§Non-goals (D7) augments the existing §Non-goals at line 2175; no insertion in the §Artifact contract → §Empirical anchor band.

Manifest header (lines 1–2) updated to reflect:
- new sections in manifest list
- updated completion list

## Empirical anchor

cnos#410 itself; the cnos#403 wave's largest sub. Sub 4 (#409) is the prior empirical anchor for B-lite migration shape — same δ-pinned contract, same B-lite scope ruling, same anchor-preservation discipline (Sub 4 preserved §Tracking sub-surfaces; Sub 5 preserves F1–F10, §9.1, and the 5 non-goals).

## Expected diff scope

- `src/packages/cnos.cds/skills/cds/CDS.md`: ~600–900 lines added (8 sections + non-goals split); manifest header updated.
- `src/packages/cnos.cds/docs/extraction-map.md`: Status column updates for §5–§12 (8 small edits).
- `.cdd/unreleased/410/`: γ-scaffold.md (this file), self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md, cdd-iteration.md (courtesy stub).
- `.cdd/iterations/INDEX.md`: 1 row appended.

## Commit plan

Per the pacing-note in the dispatch (incremental α commits for reviewability):

- `γ-410: scaffold for cycle/410 (Sub 5 of #403 — Mechanical/Review/Gate/Assessment/Closure/Retro/Non-goals/Large-file migration to CDS)` — this scaffold + cycle dir scaffolding.
- `α-410: §Mechanical vs judgment + §Review CLP migrated to CDS.md` — D1 + D2 (small adjacent sections).
- `α-410: §Gate migrated to CDS.md (F1–F10 preserved)` — D3 (biggest single section; F-anchors load-bearing).
- `α-410: §Assessment migrated to CDS.md (§9.1 triggers preserved; engineering levels cited)` — D4.
- `α-410: §Closure + §Retro-packaging migrated to CDS.md` — D5 + D6.
- `α-410: §Large-file authoring rule + §Non-goals split (5 software-cycle non-goals added)` — D7 + D8.
- `α-410: extraction-map.md Status updates for rows 5–12` — D10.
- `β-410: AC1–AC17 PASS; close-outs; courtesy cdd-iteration stub` — review pass + close-outs.

## CDS Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | dispatch + cnos#410 + cnos.cds/CDS.md + extraction-map + cdd source skills | cds | Sub 5 of #403 wave; B-lite migration class; 8 surface families |
| 1 Select | dispatch (pre-selected by δ) | cds | cnos#410 selected by parent-wave dependency order |
| 2 Branch | `cycle/410` | cds | Branch created from `origin/main@4a87cdf9` |
| 3 Bootstrap | this file | cds | Scaffold committed; α dispatch ready (γ+α+β collapsed) |
