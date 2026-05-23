# γ scaffold — Cycle 422

**Cycle:** [cnos#422](https://github.com/usurobor/cnos/issues/422) — Release-hygiene v3.82.0 (CCNF package-architecture baseline)
**Branch:** `cycle/422` (from `origin/main` @ `92038de4` — the merge that closed cycle/421 / Track A1 of #405 / CCNF-O survey)
**Mode:** docs-only / release-hygiene
**Collapse pattern:** γ+α+β collapsed on δ. Commits: `γ-422:`, `α-422:`, `β-422:`.
**Date:** 2026-05-23

## Intent

cnos has crossed a stable architecture boundary. CDD is the compact CCNF kernel (#402); CDS v0.1 ships (#403 wave: #406–#412 closed); CDR v0.1 ships (#376 wave: subs 1–4 closed); cnos.handoff v0.1 ships (#404 wave: #415–#420 closed); V is executable (`cn cdd verify`); domain evidence has homes. The operator (2026-05-22) ruled: **cut release v3.82.0 = CCNF package-architecture baseline; then stop expanding the protocol.**

This is **one bounded release-hygiene cycle, not a feature cycle**. Three READMEs reflect reality; VERSION bumps; RELEASE.md is authored; the protocol pauses post-tag.

## Surface plan (D1–D5)

### D1 — Rewrite `src/packages/cnos.cds/README.md` for v0.1 complete

Replace "v0.1 skeleton" / "Pending Sub 2 of cnos#403" / "Forthcoming surfaces" / "in flight" language with v0.1-complete framing. Mirror `cnos.handoff/README.md` post-#420 shape. Mark `CDS.md` (#407), `SKILL.md` (#406), `extraction-map.md` (#406), the operational sub-skill overlays at `skills/cds/lifecycle/` + `skills/cds/selection/` (landed across #408–#410), and `docs/empirical-anchor-cdd.md` (#412) as **Landed** with sub-issue references. Replace status block with `**v0.1 complete — cnos#403 wave closed 2026-05-22.**` plus a one-paragraph wave-shape narrative naming subs 1–7.

Reality note (does not require skill edits): the v0.1 wave shipped **operational sub-area overlays** (`lifecycle/SKILL.md`, `selection/SKILL.md`) rather than per-role overlays (`alpha/beta/gamma/delta/epsilon`). The old README anticipated role overlays; reality landed sub-area overlays plus canonical CDS.md sections (per the extraction map's Sub 3/4/5 split). The new README describes what actually shipped.

### D2 — Rewrite `src/packages/cnos.cdr/README.md` for v0.1 complete

Replace "v0.1 skeleton" / "Role overlays (Sub 3) and empirical-anchor doc (Sub 4) are in flight" with v0.1-complete framing. Mark `CDR.md` (#390), `SKILL.md` (#394), role overlays at `skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md` (#395), and `docs/empirical-anchor-cph.md` (#396) as Landed under the #376 wave. Note: the role overlays directory is `operator/` (the δ role's research-side directory name), not `delta/`.

### D3 — Confirm `cnos.handoff/README.md` is v0.1 complete

Pre-dispatch read confirms it already declares `**v0.1 complete — cnos#404 wave closed 2026-05-22.**` (line 63). No edit required; verify only.

### D4 — Bump `VERSION` from `3.81.0` → `3.82.0`

Single-line edit; bare semver, no `v` prefix (per `cnos.cdd/skills/cdd/release/SKILL.md §2.6` "tag naming convention").

### D5 — Author `RELEASE.md` for v3.82.0 at `.cdd/releases/3.82.0/RELEASE.md`

Per AC5 of cnos#422 (explicit canonical path: `test -f .cdd/releases/3.82.0/RELEASE.md`). The release/SKILL.md §2.5 default location is repo root; the issue overrides to the per-release directory. Required sections:

- **Title:** v3.82.0 — CCNF package-architecture baseline
- **Includes:** CCNF kernel (#402); CDS v0.1 (#403 wave: #406–#412); CDR v0.1 (#376 wave: subs 1–4); cnos.handoff v0.1 (#404 wave: #415–#420); essays `CCNF-AND-TYPED-TRUST.md` + `DECREASING-INCOHERENCE.md` (#414); Sigma activation bundle staged + operator-applied (#413); CCNF-O / TSC steering roadmap (#405) open but not executed; Track A1 survey doc (#421).
- **Does NOT include:** TSC report attachment; `IssueProposal.v1`; `RiskPolicy.v1`; CCNF-O schemas; field-trial results; wave-executor / new runtime / autonomy work.
- **Stop condition:** Post-tag, pause protocol evolution. Next phase = field application of CDS/CDR + handoff/memory-return testing.

The release/SKILL.md §2.5 RELEASE.md format (Outcome / Why it matters / Added / Changed / Fixed / Removed / Validation / Known Issues) is the shape; this release is the CCNF package-architecture baseline, so Outcome / Why it matters lead, then Added (the four package families + essays + bundle + roadmap), Changed (CDD kernel compression to 160 lines), Removed (none), Validation (post-merge δ-side tag + scripts/release.sh), Known Issues (stale anticipatory wording in `skills/cds/SKILL.md` v0.1-caveat block — flagged as post-v0.1 follow-up; not blocking).

## Implementation contract

| Axis | Pinned value |
|---|---|
| Language | Markdown + plain-text VERSION file |
| CLI integration target | None; `scripts/release.sh` is operator-side post-merge |
| Package scoping | **Edits only:** `src/packages/cnos.cds/README.md`, `src/packages/cnos.cdr/README.md`, `VERSION`, new `.cdd/releases/3.82.0/RELEASE.md`, cycle-close artifacts in `.cdd/unreleased/422/`, INDEX.md row |
| Existing-binary disposition | N/A; tag itself (`v3.82.0`) is operator-side post-merge action |
| Runtime dependencies | None |
| JSON/wire contract | N/A — no schemas; no new packages |
| Backward compat | All existing functionality preserved; documentation + version bump only |

## Acceptance criteria

AC1–AC11 per [cnos#422](https://github.com/usurobor/cnos/issues/422). All mechanical. Verified in `self-coherence.md` post-α.

- AC1: `cat VERSION` → `3.82.0`
- AC2: cds/README "v0.1 complete" hit ≥ 1; "skeleton|in flight|pending sub|forthcoming sub-deliverables" → 0 current-state hits
- AC3: cdr/README "v0.1 complete" hit ≥ 1; "in flight|skeleton" → 0 current-state hits
- AC4: handoff/README "v0.1 complete" hit ≥ 1 (already true; verify)
- AC5: `test -f .cdd/releases/3.82.0/RELEASE.md` exits 0; contains "Includes" + "Does NOT include" + "Stop condition" sections
- AC6: `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/{CDD.md,COHERENCE-CELL.md,COHERENCE-CELL-NORMAL-FORM.md}` → 0 lines
- AC7: `test ! -d schemas/ccnf-o`, `test ! -d schemas/handoff`, `test ! -d src/packages/cnos.ccnf-o`, `test ! -d src/packages/cnos.coherence`; `git diff origin/main..HEAD -- schemas/` → 0 lines
- AC8: no diff in `src/packages/cnos.cdd/commands/cdd-verify/` `src/go/` `scripts/release.sh` `src/packages/cnos.cdd/skills/cdd/harness/` `src/packages/cnos.cdd/skills/cdd/release-effector/`
- AC9: no new files matching `track-a|track-b|ccnf-o|tsc-report|issue-proposal|risk-policy|coherence-controller` patterns (outside this scaffold mentioning them as out-of-scope)
- AC10: no diff in `src/packages/cnos.cdd/skills/`, `src/packages/cnos.cds/skills/`, `src/packages/cnos.cdr/skills/`, `src/packages/cnos.handoff/skills/`
- AC11: no diff in `src/packages/cnos.core/` `src/packages/cnos.eng/` `src/packages/cnos.kata/` `src/packages/cnos.cdd.kata/`

## Hard rules

1. **No CCNF kernel changes.** CDD.md / COHERENCE-CELL.md / COHERENCE-CELL-NORMAL-FORM.md untouched.
2. **No new schemas.** `schemas/handoff/`, `schemas/ccnf-o/` MUST NOT exist after this cycle.
3. **No new packages.** No `cnos.ccnf-o/`, no `cnos.coherence/`, etc.
4. **No skill content changes.** Only README + VERSION + RELEASE.md edits permitted; `skills/` subdirectories not touched beyond reading. The stale "Pending Subs 3–5" / "v0.1 caveat" wording inside `cnos.cds/skills/cds/SKILL.md` is left in place and noted as a Known Issue in RELEASE.md (post-v0.1 follow-up).
5. **No `cn cdd verify` / runtime / harness / release-effector changes.**
6. **No #405 / Track A / Track B work.** Track A1 closed at #421; remaining Tracks are deferred until post-tag field evidence.
7. **Do NOT push the tag.** `scripts/release.sh` is operator-side post-merge.
8. **cnos.core / cnos.eng / cnos.kata / cnos.cdd.kata untouched.**

## Non-goals

- Do NOT dispatch #405 Track A2 / Track B1 / any other #405 sub.
- Do NOT define `IssueProposal.v1` / `RiskPolicy.v1` / `TSCReport` / `CoherenceController`.
- Do NOT add `schemas/ccnf-o/` or `schemas/handoff/` CUE.
- Do NOT add new package families.
- Do NOT modify CCNF kernel.
- Do NOT launch CDS/CDR field trials (deferred to post-tag phase).
- Do NOT push tag (operator-side via `scripts/release.sh` post-merge).
- Do NOT modify skill files beyond what `release/SKILL.md` may require for the release artifact (it does not require any).
- Do NOT modify Sigma activation bundle (already operator-applied to cn-sigma).
- Do NOT modify `CHANGELOG.md` ledger row in this scaffold — the release-side artifact authority pins ledger-row authorship to the β-release commit; this cycle's `RELEASE.md` carries the canonical release narrative; CHANGELOG ledger row authoring under §2.4 may be a δ-side post-merge step (operator decides; both options are coherent with the issue scope, which does not name a CHANGELOG change in D1–D5).

## Operator action on close

After all 11 ACs PASS, push: `git push -u origin cycle/422`. Then merge to main with `--no-ff` and message `Merge cycle/422: Release-hygiene v3.82.0 — CCNF package-architecture baseline. Closes #422.`. Post-merge: push the `v3.82.0` tag via `scripts/release.sh` (the canonical δ-side release-effector flow per `cnos.cdd/skills/cdd/release/SKILL.md §2.6` + `release-effector/SKILL.md`). Post-tag: pause protocol evolution.
