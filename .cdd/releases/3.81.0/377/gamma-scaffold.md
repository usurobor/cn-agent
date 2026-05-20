# γ scaffold — cycle #377

**Issue:** cnos#377 — design(cdd): codify cross-repo coordination protocol (intake / outbound / bilateral)
**Branch:** `cycle/377` from `origin/main` at `dd5a36d`
**Mode:** design-and-build (per issue body + wave manifest)
**Wave:** 2026-05-19 protocol hygiene (`.cdd/waves/2026-05-19-protocol-hygiene/manifest.md`)
**Role collapse:** γ=δ=α=β collapsed on δ for a single design-and-build cycle under §5.2 wave-mode with γ=δ collapse permitted. α≠β within a session is structurally compromised but acceptable for the design phase converging on the build per wave manifest precedent. Named explicitly in receipt and self-coherence.

## Selection
Per `CDD.md` §3: this is a wave-dispatched cycle from the 2026-05-19 protocol hygiene wave. The selection authority is the wave manifest. Decisive clause: wave manifest §"Issues" row #377 — design-and-build, surfaces `cdd/cross-repo/SKILL.md` (new) + `cdd/gamma/SKILL.md` §2.1/§2.7 + `cdd/post-release/SKILL.md` §5.6b + `.cdd/iterations/cross-repo/README.md`.

## §Peer enumeration at scaffold time (gamma §2.2a)

Per `gamma/SKILL.md §2.2a`, γ must list/grep before asserting "X does not exist". Performed:

- `ls src/packages/cnos.cdd/skills/cdd/` → no `cross-repo/` subdirectory present. **Confirmed**: `cdd/cross-repo/SKILL.md` does not exist.
- `rg "cross-repo" src/packages/cnos.cdd/skills/cdd/` → matches only in `gamma/SKILL.md §2.1, §2.7`, `post-release/SKILL.md §5.6b`, and one mention in CDD.md §5.3b closure-gate. **Confirmed**: protocol fragments are scattered across these three skills only; no central surface.
- `ls .cdd/iterations/cross-repo/` → 8 anchor directories exist (tsc, cph, gait-support-paths, cn-sigma×2, cn-rho, cph×3 inc. coherence-drift-sweep-followup, issue-32-tightening). **Confirmed**: 8 empirical anchors available for retroactive validation.

The §Gap framing is **additive** — this cycle introduces a new canonical surface and converts existing fragments to one-line references; no surface deletion required.

## Design surface

Five files in scope:
- **C**reate: `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md`
- **E**dit: `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (§2.1 + §2.7 — reduce inline doctrine to one-line summary + reference)
- **E**dit: `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` (§5.6b — reduce inline doctrine to one-line summary + reference)
- **E**dit: `.cdd/iterations/cross-repo/README.md` (align bundle-state phases with the canonical state machine; cross-reference the new skill)
- Cycle artifacts under `.cdd/unreleased/377/` (gamma-scaffold, design-notes, self-coherence, beta-review, alpha/beta/gamma-closeout)

## AC oracle approach per AC

| AC | Invariant | Oracle |
|---|---|---|
| AC1 | New skill exists with standard CDD frontmatter + ≥2 referencing matches | Read frontmatter; `rg "cdd/cross-repo" src/packages/cnos.cdd/skills/cdd/{gamma,post-release}/SKILL.md` ≥2 |
| AC2 | STATUS state machine: 5 events + `drafted` reconciliation + transition graph + emitters + master/sub `landed` rule | Read skill body §"STATUS state machine"; verify table/graph + master/sub rule named |
| AC3 | Directional cases × bundle file set + LINEAGE schema per case; both precedents validate without contradiction | Retroactive validation table in design-notes.md against all 8 anchors |
| AC4 | One canonical path per directional case; gamma §2.1 intake-scan updated to match; no stale path candidates | `rg "\.cdd/iterations/proposals\|\.cdd/proposals/" src/packages/cnos.cdd/skills/cdd/` returns zero or deprecation-only |
| AC5 | Feedback-patch format + bundle archival rule codified in skill body; doctrine removed from bundle READMEs | Skill has "Feedback patch" + "Bundle archival" subsections; `.cdd/iterations/cross-repo/README.md` carries bundle-shape but not protocol doctrine |
| AC6 | Gamma §2.1/§2.7 + post-release §5.6b reference new skill; `issue/SKILL.md` Source Proposal block unchanged | `rg "cdd/cross-repo" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` ≥3; `git diff` shows no change to `cdd/issue/SKILL.md` |

## Empirical anchor validation plan

The new protocol must validate against all 8 anchors without contradiction. Design-phase output `design-notes.md` includes a retroactive validation table with one row per anchor naming: directional case, bundle file set actually used, LINEAGE.md schema, STATUS state machine compliance, and observed gap (the "gap" column drives which protocol features the skill must encode). If any anchor contradicts the sketch, the sketch is revised.

Anchors:
1. `tsc/cdd-supercycle` — outbound iteration trace (6 patches, 1:N)
2. `cph/bootstrap-cdr` — inbound proposal master/sub (post-rename canonical)
3. `gait-support-paths/bootstrap-cdr` — inbound proposal master/sub (historical mirror, pre-rename)
4. `cn-sigma/agent-activate-skill` — inbound proposal 1:1 (with `drafted` event gap)
5. `cph/coherence-drift-sweep-followup-2026-05-18` — bilateral iteration (cross-repo ε; patches in both repos; hat-collapse)
6. `cn-rho/bootstrap-2026-05-19` — operator-pending bundle for non-existent target repo
7. `cn-sigma/discipline-section-2026-05-19` — feedback patch for existing target repo
8. `cph/issue-32-tightening-2026-05-19` — proposal as GitHub issue comment (not file patch)

## Expected diff scope

`git diff origin/main..HEAD --stat` should show approximately:
- `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` | +600..800 / -0 (new skill)
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | +6..12 / -30..50 (extract §2.1 + §2.7 doctrine, replace with reference)
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` | +3..6 / -8..12 (extract §5.6b doctrine, replace with reference)
- `.cdd/iterations/cross-repo/README.md` | +10..20 / -3..5 (align phases, reference new skill)
- `.cdd/unreleased/377/*.md` | +400..700 / -0 (cycle artifacts)

Per wave manifest invariant 1: changes only under `src/packages/cnos.cdd/skills/cdd/`, `.cdd/iterations/cross-repo/`, or `.cdd/unreleased/377/`. **Confirmed in scope.**

## Cross-cycle coordination

Per wave manifest §"Cross-cycle coordination": #375 touches `gamma/SKILL.md` §2.5; #377 touches `gamma/SKILL.md` §2.1 + §2.7 (different sections). #378 touches `review/SKILL.md`, `alpha/SKILL.md`, `operator/SKILL.md` (different files entirely). No direct file collision expected; possible merge-time conflict on `gamma/SKILL.md` if #375 lands first — γ resolves by integrating both patches at merge time.

## Dispatch

γ-as-δ collapsed on δ; α and β phases run in the same session sequentially per wave precedent. α produces design-notes.md + skill files + self-coherence.md; β-collapsed-on-α produces beta-review.md.

Authored 2026-05-20 by γ (collapsed on δ).
