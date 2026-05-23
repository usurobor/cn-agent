# γ Scaffold — cycle/409

**Issue:** [cnos#409](https://github.com/usurobor/cnos/issues/409) — Sub 4 of #403: Migrate §Coordination surfaces + §Artifact contract to CDS (B-lite thin extract)
**Mode:** design-and-build (sub-issue of cnos#403 wave)
**Wave:** cnos#403 (cnos.cds bootstrap; Subs 1, 2, 3 landed at `987acd04`, `d9829412`, `5f13f61c`)
**Branch:** `cycle/409` (created from `origin/main` at `5f13f61c`)
**Priority:** P2

## Dispatch shape

**γ+α+β collapsed on δ** per the breadth-2026-05-12 wave manifest precedent (CDS §"Six-field instantiation contract" Field 6 — actor collapse rule, "β-α-collapse-on-δ for skill/docs-class cycles"). This cycle is canonical-content migration with no novel executable surface — the failure modes that α/β independence catches are reduced by the absence of new compiled code; the matter is markdown content moved from one canonical home (pre-#402 CDD.md, mined from `8f06a606^`) to another (CDS.md), under δ's pinned B-lite contract. The configuration-floor cap on γ-axis and β-axis at A- is acknowledged (per `release/SKILL.md §3.8`); declared formally in `gamma-closeout.md §"Configuration-floor declaration"`.

## Pinned implementation contract (per δ)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None new |
| Package scoping | `src/packages/cnos.cds/skills/cds/CDS.md` (add two new top-level sections) + optional thin overlays at `skills/cds/tracking/SKILL.md` + `skills/cds/artifacts/SKILL.md` + `docs/extraction-map.md` Status column for rows 3+4 |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | `cnos.cdd` **NOT modified** (hard rule; Sub 6 sweeps CDD.md). `cnos.cdr` **NOT modified**. No `.cdd/` → `.cds/` filesystem migration (document planned re-rooting only). No relocation of `cross-repo/SKILL.md` (defer to #404). |

## Surfaces γ expects α to touch

- `src/packages/cnos.cds/skills/cds/CDS.md` — insert two new top-level sections (`## Coordination surfaces`, `## Artifact contract`) immediately after `## Development lifecycle` (the Sub 3 deliverable) and before `## Empirical anchor`. Update the section manifest at the top of the file accordingly.
- `src/packages/cnos.cds/docs/extraction-map.md` — update Status field for rows 3 (§Coordination surfaces) and 4 (§Artifact contract). Other rows untouched.
- Optional thin overlay files at `src/packages/cnos.cds/skills/cds/tracking/SKILL.md` and/or `src/packages/cnos.cds/skills/cds/artifacts/SKILL.md` (each ≤ 40 lines) — γ's call to ship; CDS.md alone suffices for the canonical commitment. **Decision: skip both optional overlays.** Rationale: CDS.md's two new sections plus the `Operational realization` pointers at the end of each section already satisfy the canonical-home commitment; introducing thin overlays would add a second pointer-only file class without changing the binding (the overlay would just delegate back to the cdd-side operational realization that CDS.md already cites). Sub 3 made the same call for `selection/` and `lifecycle/`.

## AC oracle approach

Each AC has a mechanical or read-check oracle (per issue #409):

- **AC1, AC2:** `grep -c "^## Coordination surfaces" CDS.md` / `grep -c "^## Artifact contract" CDS.md` returns exactly 1 each; sections non-empty.
- **AC3:** `grep -c "^### Cycle-state evidence\|^### Polling primitives\|^### Mid-flight clarification\|^### Cross-repo proposals" CDS.md` returns exactly 4 within the §Coordination surfaces block.
- **AC4:** `grep -c "^### Terminology\|^### Bootstrap\|^### Ordered flow\|^### Manifest\|^### Location matrix\|^### Ownership matrix\|^### Trace format\|^### Supporting rules\|^### Frozen snapshot rule" CDS.md` returns exactly 9 within §Artifact contract.
- **AC5:** read-check — each new top-level section ends with an `### Operational realization` sub-heading citing ≥1 cdd skill file.
- **AC6:** read-check — `.cdd/` → `.cds/` re-rooting mentioned as **planned (not performed)** in §Cycle-state evidence AND §Location matrix. `git diff origin/main..HEAD -- '.cdd/'` shows no `.cdd/` → `.cds/` rename/move ops.
- **AC7:** `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns 0 lines (hard rule).
- **AC8:** `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns 0 lines (hard rule).
- **AC9:** extraction-map Status updated for rows 3+4; rows 1+2 (Sub 3) preserved; other rows unchanged.
- **AC10:** `find src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon,operator} -type f 2>/dev/null` returns empty. Only `skills/cds/tracking/` and/or `skills/cds/artifacts/` permitted (γ skips both).

## Source mining map (B-lite — where the canonical rule text comes from)

Per the B-lite scope ruling, canonical rules move into CDS; the operational realization stays in cnos.cdd role/runtime skills as the v0.1 overlay. Source pin for the canonical rule statements is the **pre-#402 CDD.md** (commit `8f06a606^`) which carried the rules inline before the CCNF spine rewrite quarantined them.

**For §Coordination surfaces:**
- §Tracking subsection of pre-#402 CDD.md §1.5 (lines 188–322) — the canonical content for issue activity, cycle branch state, `.cdd/unreleased/{N}/` directory state, the polling query table, wake-up mechanism, reachability preflight, transition-only emission, synchronous baseline pull, `git fetch` reliability re-probe, issue-edit cache-bust via `gamma-clarification.md`, cross-repo proposal lifecycle (8-event vocabulary + STATUS state machine).
- `harness/SKILL.md §5.4` (single-named-branch polling under `Monitor`) — operational realization. Cited from §Polling primitives.
- `gamma/SKILL.md §2.5` (cycle branch + polling cross-reference) + `§2.5 gamma-clarification.md procedure` — operational realization. Cited from §Polling primitives and §Mid-flight clarification.
- `cross-repo/SKILL.md §2.3` (STATUS state machine with full transition graph) — operational realization. Cited from §Cross-repo proposals.
- Empirical anchor for mid-flight clarification: **cnos#391** — wrong-shape implementation rescued mid-cycle via `gamma-clarification.md` cache-bust; cited explicitly in §Mid-flight clarification.

**For §Artifact contract:**
- §5 Artifact Contract of pre-#402 CDD.md (lines 907–1098) — the canonical content. Sub-sections §5.0 (terminology), §5.1 (bootstrap), §5.2 (ordered artifact flow — 13 stages), §5.3 (artifact manifest with per-step format spec), §5.3a (Artifact Location Matrix), §5.3b (role/artifact ownership matrix), §5.4 (CDD Trace format), §5.5 (supporting rules), §5.6 (frozen snapshot rule).
- `release/SKILL.md §2.5a` (release-time directory move) — operational realization. Cited from §Frozen snapshot rule and §Location matrix.
- `gamma/SKILL.md §2.10` (closure gate) — operational realization. Cited from §Ownership matrix and §Operational realization.
- `alpha/SKILL.md §2.6` (pre-review gate) + `beta/SKILL.md §"Pre-merge gate"` — operational realization. Cited from §Operational realization.

α copies these into CDS.md with light editorial framing (CDS-prefix vocabulary; engineering-discipline framing where the pre-#402 CDD.md used generic CDD vocabulary; "CDS Trace" rename in lieu of "CDD Trace") and ends each new section with an `### Operational realization` pointer at the corresponding cdd skill files.

## Planned `.cdd/` → `.cds/` re-rooting (documented; not performed)

Documented in two places per AC6:
1. §Coordination surfaces → §Cycle-state evidence — mentions the `.cdd/unreleased/{N}/` directory and notes the planned re-rooting once the project binding switches.
2. §Artifact contract → §Location matrix — names the canonical paths under `.cdd/unreleased/{N}/` with a callout that the destination naming is `<project>/.cds/unreleased/{N}/`.

The actual filesystem rename is **out of scope** for Sub 4 and is itself a separate post-#403 cycle (it requires coordinating with all in-flight `.cdd/unreleased/{N}/` directories on origin/main, and with `cn cdd verify` schema paths).

## Cross-repo SKILL.md location (open question; deferred)

Currently lives at `cnos.cdd/skills/cdd/cross-repo/SKILL.md`. CDS.md §Coordination surfaces → §Cross-repo proposals cites it from there (the operational realization stays put). The long-term home is open per #404 (the broader handoff/coordination extraction tracker). Sub 4 does NOT relocate the file; CDS.md notes the open question.

## Expected diff scope

| Surface | Change shape | Approx. line delta |
|---|---|---|
| `CDS.md` | Insert 2 sections + update section manifests | +400 to +600 lines |
| `extraction-map.md` | Update Status for rows 3+4 (Status field added inline) | +6 to +12 lines |
| `skills/cds/tracking/SKILL.md` | Optional thin overlay (γ skipping) | 0 |
| `skills/cds/artifacts/SKILL.md` | Optional thin overlay (γ skipping) | 0 |
| `.cdd/unreleased/409/*.md` | Close-out artifact set | ~7 files |
| `.cdd/iterations/INDEX.md` | One row | +1 line |
| `cnos.cdd/` | Untouched (hard rule AC7) | 0 |
| `cnos.cdr/` | Untouched (hard rule AC8) | 0 |

## Peer enumeration

Per `gamma/SKILL.md §2.2a` peer-enumeration-at-scaffold-time and `alpha/SKILL.md §2.3`:

- **§Coordination surfaces** in CDS.md — verified absent: `grep -n "^## Coordination surfaces" src/packages/cnos.cds/skills/cds/CDS.md` returns 0 matches as of `5f13f61c`. Wholly new section.
- **§Artifact contract** in CDS.md — verified absent: `grep -n "^## Artifact contract" src/packages/cnos.cds/skills/cds/CDS.md` returns 0 matches as of `5f13f61c`. Wholly new section.
- **Existing related sections in CDS.md** — Field 3 (γ close-out artifact) names the artifact-set taxonomy (per the "Sub-4-vs-Field-3 line"). Field 4 (δ cadence) mentions `.cds/unreleased/{N}/` paths. Field 5 (ε iteration cadence) declares the iteration-artifact paths. New §Artifact contract expands the canonical artifact contract that Fields 3 + 4 + 5 reference; no conflict. The Sub-3 §Development lifecycle's `cycle/{N}` branch rule references §Coordination surfaces "pending Sub 4"; the new §Coordination surfaces closes that pending reference.
- **`skills/cds/tracking/`** + **`skills/cds/artifacts/`** directories — verified absent: `ls src/packages/cnos.cds/skills/cds/` shows only `CDS.md` + `SKILL.md`. γ skips both per the dispatch-shape decision above.
- **`docs/extraction-map.md` rows 3+4** — present at the source pin (lines 99–133); current Status field absent (the Status was added inline in the section header by Sub 3 for rows 1+2). Sub 4 follows the same Status-update convention.

No additive-or-consolidation framing needed; this is a clean B-lite extract under δ's pinned contract.

## Skill loading

Tier 1a: `CDS.md` (canonical, this cycle's primary edit target), `CDD.md` (kernel, for citation), `cnos.cds/skills/cds/SKILL.md` loader (read), the collapsed-role's matching skill files (`gamma/`, `alpha/`, `beta/SKILL.md`).
Tier 1b: `design/SKILL.md` (this cycle is a canonical-content move under design/SKILL.md §3.2 "one source of truth"); `issue/SKILL.md` (AC interpretation); `review/SKILL.md` (β self-review).
Tier 2: `eng/writing` (markdown authoring); `cnos.core/skills/design` (architecture-level reasoning for the section placement decision).
Tier 3: cnos.cdr/CDR.md (read as parallel-record / structural sibling for sub-section pattern emulation); `harness/SKILL.md §5.4`, `gamma/SKILL.md §2.5`, `cross-repo/SKILL.md §2.3`, `release/SKILL.md §2.5a` (source mining).

## Empirical anchor for the migration shape

The Sub 3 (cnos#408) precedent: it landed §Selection function + §Development lifecycle into CDS.md as canonical content, with operational realization staying in the existing cdd role/runtime skills under each new section's `### Operational realization` sub-heading. Sub 4 follows the identical B-lite shape for §Coordination surfaces + §Artifact contract.

The pre-#402 CDD.md content (mined at `8f06a606^`) is itself the empirical anchor for the rule taxonomy: every rule named in extraction-map row 3 (the 4 coordination-surface sub-rules) and row 4 (the 9 artifact-contract sub-rules) was authored against real cnos cycle experience. The B-lite migration preserves the authority chain — CDS.md becomes the canonical home; cdd role skills retain operational mechanics; the rules continue to govern the same cycles.
