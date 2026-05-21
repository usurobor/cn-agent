# γ scaffold — cycle/399 (Phase 4c of cnos#366)

**Issue:** [cnos#399](https://github.com/usurobor/cnos/issues/399) — release-effector skill: tag/release/branch-cleanup mechanics extracted from `operator/SKILL.md`.

**Parent:** cnos#366 (Phase 4: δ split). This is the third of three Phase 4 cycles:
- Phase 4a (δ-role; not yet shipped) — `operator/SKILL.md` → `delta/SKILL.md` (boundary policy)
- Phase 4b (harness; not yet shipped) — dispatch mechanics → harness substrate
- **Phase 4c (this cycle)** — release-effector mechanics → `release-effector/SKILL.md`

Note on ordering: 4c ships first per δ-as-agent dispatch. 4a and 4b have not yet started; their cycle branches do not exist on `origin`. This cycle therefore leaves the rest of `operator/SKILL.md` (δ-role + harness mechanics + §3a inward membrane) entirely untouched and edits only the release-effector surface.

**Mode:** design-and-build; γ+α+β-collapsed-on-δ (cnos#393 Rule 7 β-rigor on implementation-contract conformance applies).

## Implementation contract (pinned by δ — issue body verbatim)

| Axis | Pinned value |
|---|---|
| Language | Markdown (skill file) + shell (existing `scripts/release.sh` referenced; not reimplemented) |
| CLI integration target | `scripts/release.sh` and release CI workflows referenced; this cycle codifies the contract |
| Package scoping | New skill at `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` (or equivalent; agent picks) + edits to existing `release/SKILL.md` and `operator/SKILL.md` |
| Existing-binary disposition | `scripts/release.sh` preserved; this cycle codifies the contract around it |
| Runtime dependencies | None new |
| JSON/wire contract | N/A |
| Backward compat | All existing release workflows continue to work; effector contract documented, not changed |

## Surface

Files touched:

1. **NEW** — `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` (~150–200 lines): the release-effector skill, owned by δ, mechanics-of-the-disconnect-release.
2. **EDIT** — `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`: extract §3.4 disconnect-release body + the embedded recovery runbook in §6 step "Gate"; replace with cross-references to release-effector.
3. **EDIT** — `src/packages/cnos.cdd/skills/cdd/release/SKILL.md`: cross-references already point at `operator/SKILL.md §3.4`. Update to point at release-effector.
4. **EDIT** — `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`: existing refs to `operator/SKILL.md §3.4` and `scripts/release.sh` updated to point at release-effector.
5. **EDIT** — `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md`: one ref to `scripts/release.sh per operator/SKILL.md §3.4`.
6. **EDIT** — `src/packages/cnos.cdd/skills/cdd/CDD.md`: one ref to `operator/SKILL.md §3.4` in §F3 dispatch-failure-evidence table.
7. **EDIT** — `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md`: prose ref to mechanics relocation.

Files NOT touched (Phase 4a/4b scope):
- Anywhere in `operator/SKILL.md` outside the release-effector mechanics (algorithm steps 1–5+7, §1 Route, §2 Wait, §3.1–3.3 gate policy, §3a inward-membrane, §4 Override, §5 Dispatch configurations, §6 What operator does NOT, §7 lifecycle table rows other than Disconnect, §8 timeout recovery, §9 Kata, §10 Wave coordination).
- `scripts/release.sh` itself (referenced, not modified — per issue Non-goals).
- `release/SKILL.md` β-side content (artifact set authoring, RELEASE.md drafting, pre-release validators) — per issue AC5.
- `cnos#366` parent issue — comment only on close.

## AC oracle approach (issue body verbatim)

| AC | Oracle | Surface |
|----|--------|---------|
| AC1 | `test -s src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` and frontmatter parses | new skill file |
| AC2 | `rg "scripts/release.sh\|tag creation\|branch cleanup\|release CI" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` returns 0 hits outside cross-reference position (cross-refs to release-effector skill are permitted) | operator/SKILL.md edits |
| AC3 | Tag policy text (bare `X.Y.Z` no `v` prefix per CDD §5.3a) appears in release-effector skill with the canonical CDD §5.3a reference | release-effector skill |
| AC4 | §3.4 disconnect-release content is in release-effector skill (preferred path) with operator/SKILL.md carrying a cross-reference stub; no content lost between source + destination | both files |
| AC5 | `release/SKILL.md` retains: §2.5 RELEASE.md authoring, §2.4 CHANGELOG, §2.5a cycle-dir move, pre-release validators, β-side flow. release-effector does NOT absorb these. | release/SKILL.md unchanged in those sections |
| AC6 | `rg "operator/SKILL.md §3.4\|operator/SKILL.md §3\\.4" src/packages/cnos.cdd/skills/cdd/` returns 0 hits or only legacy/empirical-anchor citations; cross-refs to release mechanics point at release-effector | grep across skills |
| AC7 | `scripts/release.sh` git diff vs main is empty; the new skill cites the script as the canonical binding and matches the script's actual behavior (RELEASE.md gate validation, stamp, consistency check, dir move, commit, tag-message generation, annotated tag, push) | git diff + script vs skill |

## Branch/identity

- Branch: `cycle/399` from `origin/main` (created)
- γ identity: `gamma@cdd.cnos`
- α identity: `alpha@cdd.cnos`
- β identity: `beta@cdd.cnos`
- δ identity: `operator@cdd.cnos` (post-merge to main)

## Dispatch shape

This is `cdd/operator/SKILL.md §5.2` (single-session δ-as-γ with γ+α+β collapsed on δ-as-agent). The cycle has 7 ACs which exceeds the §5.3 escalation threshold (`≥7 ACs`), but mode is design-and-build with one new file + targeted edits to existing files; the design surface is small. Per §5.2 configuration-floor clause (`release/SKILL.md §3.8`), γ-axis is capped at A- (γ/δ separation absent).

## Risks and forecasts

- **R1: operator/SKILL.md concurrent edits from cycles/397/398.** No such branches exist on origin (`git branch -r | grep -E "cycle/(397|398)"` → empty). Risk: none at scaffold time.
- **R2: scripts/release.sh drift.** Script has not changed since last release per git log. Risk: low; this cycle reads the script as-is and codifies the contract around it.
- **R3: Ambiguity between δ-role policy vs effector mechanics.** Boundary: anything that runs a command, reads CI output, deletes a branch, or pushes a tag = effector. Anything that decides whether the gate fires = δ-role policy. The disconnect-release §3.4 algorithm IS mechanics (run script, poll CI, recovery runbook, branch delete); it moves. The §3.1-3.3 gate policy (when δ executes vs observes, how δ reports completion) stays in operator/SKILL.md (Phase 4a's relocation target).
- **R4: §3a inward membrane has nothing to do with releases.** Stays in operator/SKILL.md untouched.

## Plan order

1. Read all source files (done).
2. Draft release-effector/SKILL.md (α design phase output: `design-notes.md`).
3. Implement file changes (α build phase).
4. Self-coherence sweep with AC-by-AC verification (`self-coherence.md`).
5. β-collapsed review (`beta-review.md`).
6. Close-outs.
7. cdd-iteration.md (empty-findings expected; this is a mechanical extract).
8. Merge to main.
9. Comment on cnos#366; update INDEX.md.
10. Push & attempt branch delete.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | issue body, source skills | — | Inputs read; release-effector boundary identified |
| 1 Select | cnos#399 | — | Phase 4c selected as next 4-series cycle |
| 2 Branch | `cycle/399` | cdd | Branched from `origin/main` per CDD §4.2 |
| 3 Bootstrap | `.cdd/unreleased/399/` | cdd | Cycle dir created |
| 4 Gap | this file | — | Named: `operator/SKILL.md` mixes effector mechanics with δ-policy; release-effector skill is the surface needed |
