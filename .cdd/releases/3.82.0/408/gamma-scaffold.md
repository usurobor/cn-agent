# γ Scaffold — cycle/408

**Issue:** [cnos#408](https://github.com/usurobor/cnos/issues/408) — Sub 3 of #403: Migrate §Selection function + §Development lifecycle to CDS (B-lite thin extract)
**Mode:** design-and-build (sub-issue of cnos#403 wave)
**Wave:** cnos#403 (cnos.cds bootstrap; Subs 1, 2 landed at `987acd04`, `d9829412`)
**Branch:** `cycle/408` (created from `origin/main` at `d9829412`)
**Priority:** P2

## Dispatch shape

**γ+α+β collapsed on δ** per the breadth-2026-05-12 wave manifest precedent (CDS §"Six-field instantiation contract" Field 6 — actor collapse rule, "β-α-collapse-on-δ for skill/docs-class cycles"). This cycle is canonical-content migration with no novel executable surface — the failure modes that α/β independence catches are reduced by the absence of new compiled code; the matter is markdown content moved from one canonical home (CDD.md, mined for the actual rule statements from pre-#402 history) to another (CDS.md), under δ's pinned B-lite contract. The configuration-floor cap on γ-axis and β-axis at A- is acknowledged (per `release/SKILL.md §3.8`); declared formally in `gamma-closeout.md §"Configuration-floor declaration"`.

## Pinned implementation contract (per δ)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None new |
| Package scoping | `src/packages/cnos.cds/skills/cds/CDS.md` (add new top-level sections) + optional thin overlays at `skills/cds/selection/SKILL.md` + `skills/cds/lifecycle/SKILL.md` + `docs/extraction-map.md` Status column for rows 1+2 |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | `cnos.cdd` **NOT modified** (hard rule; Sub 6 sweeps CDD.md). `cnos.cdr` **NOT modified**. |

## Surfaces γ expects α to touch

- `src/packages/cnos.cds/skills/cds/CDS.md` — insert two new top-level sections (`## Selection function`, `## Development lifecycle`) immediately after `## Six-field instantiation contract` and before `## Empirical anchor`. Update the section manifests at the top of the file accordingly.
- `src/packages/cnos.cds/docs/extraction-map.md` — update Status column for rows 1 (§Selection function, §1 of the map) and 2 (§Development lifecycle, §2 of the map). Other rows untouched.
- Optional thin overlay files at `src/packages/cnos.cds/skills/cds/selection/SKILL.md` and/or `src/packages/cnos.cds/skills/cds/lifecycle/SKILL.md` (each ≤ 40 lines) — α decides whether to ship; CDS.md alone suffices for the canonical commitment.

## AC oracle approach

Each AC has a mechanical or read-check oracle (per issue #408):

- AC1, AC2: `grep -n "^## Selection function" CDS.md` / `grep -n "^## Development lifecycle" CDS.md` returns exactly 1 match each; sections non-empty.
- AC3: `grep -ci "<rule-name>" CDS.md` for each of the 10 rule names (within the §Selection block).
- AC4: grep verifies the 5 lifecycle components appear in §Development lifecycle.
- AC5: read-check — each new section ends with an "Operational realization" pointer citing at least one cdd role-skill file.
- AC6: `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns 0 lines (hard rule).
- AC7: `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns 0 lines (hard rule).
- AC8: extraction-map Status column updated for rows 1+2; other rows unchanged.
- AC9: no new files under `skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`; only `skills/cds/{selection,lifecycle}/` permitted as thin overlays (each ≤ 40 lines).

## Source mining map (B-lite — where the canonical rule text comes from)

Per the B-lite scope ruling, canonical rules move into CDS; the operational realization stays in cnos.cdd role/runtime skills as the v0.1 overlay. Source pin for the canonical rule statements is the **pre-#402 CDD.md** (commit `8f06a606^`) which carried the rules inline before the CCNF spine rewrite quarantined them. Specifically:

- §3 Selection Function (lines 671–745 of pre-#402 CDD.md) — the 10 rules (§3.1 P0 override, §3.2 operational infrastructure override, §3.3 assessment commitment default, §3.4 stale backlog re-evaluation, §3.5 MCI freeze check, §3.6 weakest-axis rule, §3.7 maximum leverage, §3.8 dependency order, §3.9 effort-adjusted tie-break, §3.10 no-gap case).
- §4 Development Lifecycle (lines 749–903 of pre-#402 CDD.md) — the 0–13 step table (§4.1), the S0–S12 state table (§4.1a), the branch rule (§4.2), branch pre-flight (§4.3), skill loading tiers (§4.4, tiers 1a/1b/1c/2/3).

α copies these into CDS.md with light editorial framing (CDS-prefix vocabulary; engineering-discipline framing where the pre-#402 CDD.md used generic CDD vocabulary) and ends each new section with an "Operational realization" pointer at the corresponding cdd role-skill files (gamma/alpha/beta/SKILL.md primarily; harness/SKILL.md cited under §Lifecycle).

## Expected diff scope

| Surface | Change shape | Approx. line delta |
|---|---|---|
| `CDS.md` | Insert 2 sections + update section manifests | +250 to +400 lines |
| `extraction-map.md` | Update Status column for rows 1+2 | +5 to +10 lines (Status text added) |
| `skills/cds/selection/SKILL.md` | Optional thin overlay | 0 or ~25 lines |
| `skills/cds/lifecycle/SKILL.md` | Optional thin overlay | 0 or ~25 lines |
| `.cdd/unreleased/408/*.md` | Close-out artifact set | ~7 files |
| `.cdd/iterations/INDEX.md` | One row | +1 line |
| `cnos.cdd/` | Untouched (hard rule AC6) | 0 |
| `cnos.cdr/` | Untouched (hard rule AC7) | 0 |

## Peer enumeration

Per `gamma/SKILL.md §2.2a` peer-enumeration-at-scaffold-time and `alpha/SKILL.md §2.3`:

- **§Selection function** in CDS.md — verified absent: `grep -n "^## Selection function" src/packages/cnos.cds/skills/cds/CDS.md` returns 0 matches as of `d9829412`. This is a wholly new section, not a partial overlap.
- **§Development lifecycle** in CDS.md — verified absent: `grep -n "^## Development lifecycle" src/packages/cnos.cds/skills/cds/CDS.md` returns 0 matches as of `d9829412`. Wholly new.
- **Existing related sections in CDS.md** — Field 4 (δ cadence) names the per-cycle/per-release/per-wave cadence taxonomy (the "Sub-3-vs-Field-4 line" delegates lifecycle detail to Sub 3 = this cycle). Field 4 stays as-is; new §Development lifecycle expands per-cycle into the explicit 0–13 steps + S0–S12 state machine. No conflict; the contract was authored anticipating this delegation.
- **`skills/cds/selection/`** + **`skills/cds/lifecycle/`** directories — verified absent: `ls src/packages/cnos.cds/skills/cds/` shows only `CDS.md` + `SKILL.md`. Thin overlays would be wholly new.
- **`docs/extraction-map.md` rows 1+2** — present at the source pin; Status column currently does not exist as a separate column (rows have a "Note" column; the Status update will add explicit "v0.1 migrated" markers to the Note column or to a new Status indicator at section level).

No additive-or-consolidation framing needed; this is a clean B-lite extract under δ's pinned contract.

## Skill loading

Tier 1a: `CDS.md` (canonical, this cycle's primary edit target), `CDD.md` (kernel, for citation), `cnos.cds/skills/cds/SKILL.md` loader (read), the collapsed-role's matching skill files (`gamma/`, `alpha/`, `beta/SKILL.md`).
Tier 1b: `design/SKILL.md` (this cycle is a canonical-content move under design/SKILL.md §3.2 "one source of truth"); `issue/SKILL.md` (AC interpretation); `review/SKILL.md` (β self-review).
Tier 2: `eng/writing` (markdown authoring); `cnos.core/skills/design` (architecture-level reasoning for the section placement decision).
Tier 3: cnos.cdr/CDR.md (read as parallel-record / structural sibling for sub-section pattern emulation).

## Empirical anchor for the migration shape

The cnos.cdr v0.1 wave (cnos#376) is the structural precedent: it moved canonical research-discipline rules into `cnos.cdr/skills/cdr/CDR.md` while leaving operational realization in the existing CDR-side runtime skills. The same A-decision-applied-thrice pattern (cnos#388 schemas → cnos#376 cdr skills → cnos#403 cds skills) is the doctrinal basis for the B-lite scope; CDR.md is the structural sibling that demonstrates the destination shape.

The pre-#402 CDD.md rule statements are themselves the empirical anchor for the rule taxonomy: every rule named in extraction-map row 1 (the 10 selection rules) and row 2 (the 5 lifecycle components) was authored against real cnos cycle experience. The B-lite migration preserves the authority chain (CDS.md cites the rules; cdd role skills retain operational mechanics; the rules continue to govern the same cycles).
