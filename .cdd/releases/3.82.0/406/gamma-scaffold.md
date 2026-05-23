# γ scaffold — cycle/406 (Sub 1 of cnos#403)

**Issue:** [cnos#406](https://github.com/usurobor/cnos/issues/406) — Sub 1 of [cnos#403](https://github.com/usurobor/cnos/issues/403): Bootstrap cnos.cds package skeleton + extraction map.

**Parent:** [cnos#403](https://github.com/usurobor/cnos/issues/403) (cnos.cds bootstrap tracker — extract-by-reference v0.1 wave; this is Sub 1 of 7; gates Sub 2).

**Mode:** design-and-build; γ+α+β collapsed on δ (per the breadth-2026-05-12 wave manifest precedent for skill/docs-class cycles).

## Implementation contract (pinned by δ — issue body verbatim)

| Axis | Pinned value |
|---|---|
| Language | Markdown skills + JSON manifest (mirror `cnos.cdr` v0.1 shape) |
| CLI integration target | None new; package loads via existing `cn` package-discovery |
| Package scoping | `src/packages/cnos.cds/` (new package, peer of `cnos.cdd`, `cnos.cdr`) |
| Existing-binary disposition | N/A; no executables in this sub |
| Runtime dependencies | None |
| JSON/wire contract | `cn.package.json` schema `cn.package.v1`, version `0.1.0` |
| Backward compat | `cnos.cdd` not modified; "pending cds extraction" markers stay until Sub 6 |

## Surface

Files created (D1 + D2 per #406):

1. **NEW** `src/packages/cnos.cds/cn.package.json` — manifest (verbatim shape from `cnos.cdr/cn.package.json` with `name: cnos.cds`).
2. **NEW** `src/packages/cnos.cds/README.md` — purpose / package structure / quick start / cross-protocol relationship (mirrors `cnos.cdr/README.md` section shape; ≥ 50 lines per AC1).
3. **NEW** `src/packages/cnos.cds/skills/cds/SKILL.md` — loader skill (mirrors `cnos.cdr/skills/cdr/SKILL.md` pattern; advisory `calls:` for forthcoming CDS.md and role overlays).
4. **NEW** `src/packages/cnos.cds/skills/cds/.gitkeep` — placeholder until Sub 2 lands CDS.md.
5. **NEW** `src/packages/cnos.cds/docs/extraction-map.md` — the D2 deliverable; one table per surface group from #403's "Source content" list; each row names source CDD §-anchor, destination under `cnos.cds/`, migration sub-issue (Sub 3/4/5), one-line note.

Cycle artifacts authored under `.cdd/unreleased/406/`:

6. `gamma-scaffold.md` (this file).
7. `self-coherence.md` — α self-coherence pass; AC-by-AC mechanical check.
8. `beta-review.md` — β-collapsed-on-δ review; AC verification + β-rigor.
9. `alpha-closeout.md` — α-level findings.
10. `beta-closeout.md` — β-level findings.
11. `gamma-closeout.md` — γ-level closure summary; finding dispositions.
12. `cdd-iteration.md` — courtesy empty-findings stub per cycle/401 rule (`protocol_gap_count: 0`).
13. `.cdd/iterations/INDEX.md` — row appended for cycle 406 (per the courtesy convention).

Files NOT touched (per #406 Non-goals + design discipline):

- `src/packages/cnos.cdd/**` — explicitly out of scope; AC6 mechanical check verifies `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns empty.
- `src/packages/cnos.cdr/**` — explicitly out of scope; AC7 is informational only.
- `schemas/cds/**` — already exists per #388; no edits needed (CDS.md will cite it once Sub 2 lands).
- `src/packages/cnos.cds/skills/cds/CDS.md` — Sub 2 territory.
- `src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` — Subs 3–5 territory.

## AC oracle approach (issue body verbatim)

| AC | Oracle | Surface |
|----|--------|---------|
| AC1 | All four named files exist; README ≥ 50 lines; SKILL.md has valid frontmatter (name, description, artifact_class, governing_question, triggers, calls). | filesystem + `wc -l` + frontmatter inspection |
| AC2 | `go run ./src/go/cmd/cn build --check` reports `cnos.cds: valid`. No changes to `cnos.core` discovery code. | `cn build --check` |
| AC3 | README declares CDS as software-development realization of CCNF; cites kernel docs; cites CDR as sibling; cites #388 architecture choice. | read-check of `README.md` |
| AC4 | SKILL.md mirrors cdr loader pattern: frontmatter, Load order, Rule, Role overlays, Cross-protocol relationship; `calls:` names forthcoming files per cdr precedent. | structural diff vs `cnos.cdr/skills/cdr/SKILL.md` |
| AC5 | `grep -c "^|" docs/extraction-map.md` ≥ 12 (10 rows + table header + separator). Each row's destination resolves under `cnos.cds/`; each row's migration-sub field names Sub 3/4/5. | `grep -c` + read-check |
| AC6 | `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/CDD.md` returns empty. `skills/cds/CDS.md` does not exist. No role-overlay SKILL.md under `skills/cds/{alpha,beta,gamma,delta,epsilon}/`. | `git diff` + `ls` |
| AC7 | CDR README + SKILL.md references to CDS by name now resolve to an existing peer package. | informational read-check; no edits to cdr |

## Branch/identity

- Branch: `cycle/406` from `origin/main` (HEAD == `ecbcb5d5` Merge cycle/402; created).
- Worktree: `worktree-agent-ab215a55af8d3dc27` (harness-created); cycle/406 branched from current HEAD.
- γ identity: `gamma@cdd.cnos`
- α identity: `alpha@cdd.cnos`
- β identity: `beta@cdd.cnos`
- δ identity: `operator@cdd.cnos` (post-merge to main; merge by operator, not by this agent per dispatch).

## Dispatch shape

This is `cdd/operator/SKILL.md §5.2` (γ+α+β-collapsed-on-δ-as-agent). Mode is design-and-build (extraction map = design; package skeleton = build); 7 ACs; ε is in-cycle per the cycle/401 cadence rule (only if `protocol_gap_count > 0`). Per `release/SKILL.md §3.8` configuration-floor clause, γ-axis is capped at A- (γ/δ separation absent in this collapse pattern); β-axis capped at A- (β-α collapse acknowledged per Rule 7).

## Risks and forecasts

- **R1: name-overpromise in `calls:` frontmatter.** SKILL.md's `calls:` lists `CDS.md`, `alpha/SKILL.md`, …, `epsilon/SKILL.md` — none exist yet (Sub 2+ territory). Per cdr-Sub-1 precedent, this is acceptable as "advisory targets"; the v0.1 caveat section in SKILL.md makes the forthcoming nature explicit. β to verify this is not over-promising readiness.
- **R2: extraction-map destinations are commitments, not suggestions.** Per #406 active design constraints: "If a destination is uncertain, name the uncertainty in the row's note column rather than guessing." α records uncertainty in §14 (Open questions) of extraction-map.md (5 open coordination questions: `.cdd/→.cds/` filesystem rename; cross-repo / harness / operator SKILL location; release-effector location). These are not destination uncertainties (the doctrine-surface destinations are stable); they are migration-coordination questions about whether operational realizations also move. β to verify the destinations themselves are stable.
- **R3: cnos.cdd modification.** AC6 hard rule. Mechanical check is `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returning empty; verified PASS at α-commit-time.
- **R4: extraction-map coverage gaps.** AC5 hard rule. Mechanical check is row count ≥ 12 (10 surface groups + header+separator). Verified PASS: 89 |-prefixed lines across 12 surface-group tables. Coverage-verification table in §13 explicitly maps every CDD.md "pending cds extraction" marker to its containing table.
- **R5: cn discovery failure.** AC2 hard rule. Verified PASS: `cn build --check` reports `cnos.cds: valid`.

## Plan order

1. ✅ Read all source files (#406 body, #403 body, cnos.cdr template files, CDD.md "pending cds extraction" markers, design SKILL §3.10/§3.11, cycle 401 close-out templates).
2. ✅ Branch `cycle/406` from `origin/main` HEAD.
3. ✅ Author D1 (package skeleton: cn.package.json, README.md, skills/cds/SKILL.md, skills/cds/.gitkeep).
4. ✅ Author D2 (docs/extraction-map.md with 12 surface-group tables, coverage-verification §13, open-questions §14).
5. ✅ Mechanical AC check (file existence, line counts, grep counts, git diff cdd-untouched, `cn build --check`).
6. ✅ α commit (role tag `α-406`).
7. Self-coherence sweep (`self-coherence.md`).
8. β-collapsed review (`beta-review.md`).
9. Close-outs (α, β, γ).
10. cdd-iteration courtesy stub (per cycle/401 rule for `protocol_gap_count == 0`).
11. INDEX.md row (per cycle/401 courtesy convention).
12. β+γ+cdd-iteration commit (role tag `β-406` blended with `γ-406`).
13. γ-INDEX commit (role tag `γ-406`).
14. Push `cycle/406` to origin.
15. Report back to operator with branch name, commits, AC summary, and merge instruction. **Do NOT merge to main.**

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | #406 body, #403 body, CDD.md "pending cds extraction" markers (lines 122–141), cnos.cdr template files | cdd, design, issue/contract, issue/proof | Inputs read; implementation contract pinned in #406 body |
| 1 Select | cnos#406 | — | Sub 1 of #403; first dispatchable cycle of the wave |
| 2 Branch | `cycle/406` | cdd | Branched from `origin/main` (HEAD `ecbcb5d5`) per CDD §4.2 / #406 dispatch |
| 3 Bootstrap | `.cdd/unreleased/406/` | cdd | Cycle dir created |
| 4 Gap | this file | — | Named: cnos.cds package does not exist; extraction map does not exist; 14 "pending cds extraction" markers in CDD.md await Subs 3–5 dispatch surface |
