<!-- sections: [Cycle, Verdict, Merge Evidence, Review Evidence, Findings, Debt, Skills-That-Helped] -->
<!-- completed: [Cycle, Verdict, Merge Evidence, Review Evidence, Findings, Debt, Skills-That-Helped] -->

# β Closeout — Cycle #390

**Cycle:** [#390](https://github.com/usurobor/cnos/issues/390) — Sub 1 of [#376](https://github.com/usurobor/cnos/issues/376): CDR six-field instantiation contract + architecture-choice declaration
**Date:** 2026-05-21
**Branch:** `cycle/390`
**Pre-merge cycle SHA:** `4eaa11443f8bf8a7efc59b6ed53bdb3d462e41aa` (α + scaffold + design-notes + self-coherence) plus subsequent close-out commits
**Base SHA:** `417b6227ba7ce7c47e02ec4e8b9614feb70b6f64`
**β identity:** beta / beta@cdd.cnos (γ+α+β-collapsed on δ, single Claude Code session — collapse acknowledged in `beta-review.md`)

---

## §Verdict (β statement)

**APPROVED — Round 1.** All six ACs pass on mechanical oracles. β-α-collapse acknowledged and structurally appropriate for docs-only contract-authoring class. β proceeds to merge.

## §Merge Evidence

**Merge action:** `git merge --no-ff cycle/390 -m "Merge cycle/390: CDR six-field contract + architecture-choice declaration (#390)"` with `Closes #390` in the merge body, then `git push origin main`.

**Merge SHA:** [filled by γ post-merge in gamma-closeout.md]

**Branch state pre-merge:**
- `origin/cycle/390` HEAD: post-close-outs (this file + alpha-closeout, gamma-closeout, cdd-iteration land on the cycle branch before merge per CDD.md §close-out conventions).
- `origin/main` HEAD: `417b6227` (cycle 388 ε artifact merge); no rebase required.

**Pre-merge integrity verification:**
- All AC oracles pass (recorded in `beta-review.md §AC Re-Check`).
- No forbidden tokens in normative CDR.md sections.
- Cited paths resolve on branch.
- No prior in-flight cycles touch `src/packages/cnos.cdr/` (peer-enumeration confirmed at scaffold time).

## §Review Evidence

| AC | Verdict | Oracle | Line(s) |
|---|---|---|---|
| AC1 | ✓ | `rg "^### Field" CDR.md \| wc -l` → 6; `rg -c "TBD"` → 0 | CDR.md 244–494 |
| AC2 | ✓ | `## Architecture choice` present; option (a) named; option (b) rejected; cnos#388 cited | CDR.md 55–149 |
| AC3 | ✓ | `## Persona, Protocol, Project` present; three layers + canonical homes | CDR.md 154–230 |
| AC4 | ✓ | cph cited; shape-compatibility claim; Sub 4 deferral | CDR.md 495–547 |
| AC5 | ✓ | forbidden-token sweep classified; persona-identity absent | CDR.md (entire) |
| AC6 | ✓ | Sub 2/3/4 read-through walk; no re-derivation needed | CDR.md (entire) |

**Narrowing pattern across rounds:** N/A (R1 APPROVE, no fix-round).

**Mechanical-vs-judgment ratio:** mechanical = 5/6 oracles (AC1 wc/rg + AC2/3/4 grep + AC5 classified sweep); judgment = 1/6 (AC6 read-through walk, which β performed independently from α's walk). The β-α-collapse is defensible because the judgment portion is read-through over a stable doctrinal target whose mapping to Sub 2/3/4 scopes is structurally direct.

## §Findings

See `beta-review.md §Findings` (Observations 1–7 + the borderline observation on α=β collapse vs Field 6). Summary:
- 7 observations — all confirming the contract is held.
- 1 borderline (this cycle's β-α-collapse vs Field 6) — reconciled by class-identification: docs-only contract authoring is CDS-class under repairable feedback, not research-claim authoring; Field 6 governs research-claim cycles.
- 0 findings requiring revision (R1 APPROVE).

**β-side findings (factual only):**

### Bf1 — Architectural-choice rationale (1)–(5) transposition is well-distributed

Each of the five points is given a skill-specific example without padding. The transposition is dense and load-bearing; no point is filler. This validates the design-notes §ArchitecturalChoice approach (transpose, do not re-derive).

### Bf2 — Field 5 trigger-class enumeration matches the issue body precisely

The six classes (missing data gates, overclaiming, unreproducible numbers, weak citation discipline, recurring oracle ambiguity, construct drift) appear in CDR.md Field 5 in the same order and with the same names as the issue body's AC5 enumeration. This makes Sub 3's `epsilon/SKILL.md` authoring mechanical.

### Bf3 — The disavowal-prose discipline is effective at preventing surface fusion

Multiple normative sentences carry explicit disavowals: "CDR does not produce software artifacts as matter" (Field 1); "Engineering's 'compiles + tests pass' oracle does not apply" (Field 2); "not release-shaped" (Field 4); "Engineering-class collapse precedents do **not** transfer" (Field 6). The disavowal pattern is a load-bearing discipline mechanism — it prevents a reader who knows CDS from importing CDS-discipline reflexes when reading CDR.

## §Debt

(β-side; same list as α's debt; preserved here.)

Acknowledged and not raised as fix-blockers:
- Sub 2 package metadata not authored — out of scope.
- Sub 3 role overlays not authored — out of scope.
- Sub 4 cph mapping not exhaustive — out of scope; shape-compatibility claim is sufficient for Sub 1.
- CDR gate-verdict enum not pinned in schema — future cycle.
- Project-specific stricter floor template not provided — out of scope.

## §Skills-That-Helped

- `cdd/issue/proof` — β oracle authoring; each AC re-checked against a concrete mechanical oracle with line-numbered evidence.
- `cdd/review/issue-contract` (loaded by extension when reviewing contract-shaped matter) — contract-section verification discipline; "one section names one truth" check applied to each `### Field N` subsection.

β complete.
