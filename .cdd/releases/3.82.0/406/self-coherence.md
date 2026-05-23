# α self-coherence — cycle/406 (Sub 1 of cnos#403)

α writes self-coherence before β review. AC-by-AC mechanical check + β Rule 7 implementation-contract conformance.

## AC self-check

### AC1 — Package skeleton exists at canonical path

**Oracle:**

- `src/packages/cnos.cds/cn.package.json` exists with the verbatim schema/name/version/kind/engines from #406.
- `src/packages/cnos.cds/README.md` exists, ≥ 50 lines, follows the `cnos.cdr/README.md` section shape.
- `src/packages/cnos.cds/skills/cds/SKILL.md` exists with valid frontmatter (name, description, artifact_class, governing_question, triggers, calls).
- `src/packages/cnos.cds/docs/extraction-map.md` exists.

**Verification:**

```
$ ls -la src/packages/cnos.cds/ src/packages/cnos.cds/skills/cds/ src/packages/cnos.cds/docs/
src/packages/cnos.cds/:
  cn.package.json (140 bytes)
  README.md (7773 bytes)
  docs/
  skills/

src/packages/cnos.cds/skills/cds/:
  SKILL.md (9337 bytes)
  .gitkeep (179 bytes)

src/packages/cnos.cds/docs/:
  extraction-map.md (32477 bytes)

$ wc -l src/packages/cnos.cds/README.md src/packages/cnos.cds/skills/cds/SKILL.md src/packages/cnos.cds/docs/extraction-map.md
   64 src/packages/cnos.cds/README.md
  121 src/packages/cnos.cds/skills/cds/SKILL.md
  275 src/packages/cnos.cds/docs/extraction-map.md
```

All four named files exist. README is 64 lines (≥ 50 ✓). SKILL.md frontmatter inspected — names `name: cds`, `description: …`, `artifact_class: skill`, `governing_question: …`, `triggers: […]`, `calls: […]` — all six required fields present.

`cn.package.json` shape verified verbatim:

```
$ python3 -c "import json; print(json.dumps(json.load(open('src/packages/cnos.cds/cn.package.json')), indent=2))"
{
  "schema": "cn.package.v1",
  "name": "cnos.cds",
  "version": "0.1.0",
  "kind": "package",
  "engines": {
    "cnos": ">=3.81.0"
  }
}
```

Matches the #406-pinned schema/name/version/kind/engines exactly.

**Status:** PASS.

### AC2 — Package loads via existing `cn` discovery without error

**Oracle:** `cn build --check` treats `cnos.cds` as a valid v1 package. No changes to `cnos.core` discovery code required.

**Verification:**

```
$ go run ./src/go/cmd/cn build --check
✓ cnos.cdd: valid
✓ cnos.cdd.kata: valid
✓ cnos.cdr: valid
✓ cnos.cds: valid
✓ cnos.core: valid
✓ cnos.eng: valid
✓ cnos.kata: valid
✓ All packages valid.
```

`cnos.cds` is discovered alongside the existing packages and reports `valid`. The discovery code in `src/go/internal/pkg/pkg.go` is unchanged in this cycle:

```
$ git diff origin/main..HEAD -- src/go/ | wc -l
0
```

The discovery is purely filesystem-presence-based (per `pkg.go §"ContentClasses"` doc comment: "Content classes are discovered by directory presence"). `cnos.cds` declares a `skills/` directory (one of the seven canonical content-class names: `doctrine`, `mindsets`, `skills`, `extensions`, `commands`, `orchestrators`, `katas`), and the manifest parses cleanly under `FullPackageManifest` (`schema`, `name`, `version`, `kind`, `engines.cnos` — no `commands` map; allowed since `commands` is optional per the JSON shape).

**Status:** PASS.

### AC3 — README declares the cross-protocol relationship correctly

**Oracle:** README declares CDS as software-development realization of CCNF; CDD owns generic recursive kernel; CDR is sibling research realization; architecture choice (a) per #388 inherited.

**Verification:**

```
$ grep -n "software-development realization\|generic recursive\|sibling.*research\|CCNF\|#388" src/packages/cnos.cds/README.md
```

Spot-read of `README.md` (manual):

- Line 3: "CDS … is the **software-development realization** of CCNF — the generic recursive coherence-cell algorithm that lives in cnos.cdd."
- Line 3: "It sits as a peer to `cnos.cdr` (research realization); both inherit the generic kernel from `cnos.cdd` without re-deriving it."
- Line 5: "The kernel/realization split is the option (a) decision recorded at cnos#388 (schemas) and extended to skills at cnos#376 AC7 (research realization, the structural precedent for this package)."
- Line 9 ("What CDS Does"): "CDS is the engineering discipline used to improve artifacts under repairable feedback. Its loss function … is artifact improvement under repairable feedback."
- Line 11: "CDS shares the role-cell kernel (α/β/γ/δ/ε) with cnos.cdr (research). The kernel grammar — COHERENCE-CELL, COHERENCE-CELL-NORMAL-FORM, role grammar, generic receipt-validation interface — lives in cnos.cdd and is cited by reference, never restated."
- Line 13: "The architectural split — common kernel in cnos.cdd, per-protocol procedures here — inherits the option (a) decision from cnos#388 …"
- Multiple Quick Start references to CDD.md, COHERENCE-CELL-NORMAL-FORM.md, COHERENCE-CELL.md as kernel pointer destinations.
- "Looking for the research-side sibling?" section explicitly distinguishes CDS from CDR: "CDR is **not** CDS-with-different-words: engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*."

All four AC3 sub-clauses verified PASS.

**Status:** PASS.

### AC4 — SKILL.md loader pattern mirrors cdr

**Oracle:** `skills/cds/SKILL.md` mirrors `cnos.cdr/skills/cdr/SKILL.md` structure: frontmatter, Load order, Rule, Role overlays, Cross-protocol relationship. `calls:` frontmatter names `CDS.md` and forthcoming role overlays as advisory targets.

**Verification:**

Structural diff (section presence):

| Section | cdr/SKILL.md | cds/SKILL.md |
|---|---|---|
| Frontmatter | ✓ (lines 1–38) | ✓ (lines 1–38) |
| `# CDR` / `# CDS` heading | ✓ | ✓ |
| `## Load order` | ✓ | ✓ |
| `## Rule` | ✓ | ✓ |
| `## Role overlays` | ✓ | ✓ |
| `## Cross-protocol relationship` | ✓ | ✓ |
| `## Conflict rule` | ✓ | ✓ |
| Frontmatter `calls:` names CDS.md + 5 role overlays | ✓ (cdr.md + 5 overlays) | ✓ (CDS.md + 5 overlays) |
| `calls:` targets do not exist at v0.1 ship time | ✓ (cdr Sub 1 shipped loader before role overlays at Sub 3) | ✓ (cds Sub 1 ships loader before CDS.md at Sub 2 and overlays at Subs 3–5) |

cds/SKILL.md adds one section beyond cdr/SKILL.md's shape: `## v0.1 caveat`. This documents the forthcoming nature of CDS.md and the role overlays, points readers at `docs/extraction-map.md` as the interim navigation surface, and explicitly cites the cdr-Sub-1 precedent for the advisory `calls:` pattern. The addition is consistent with the cdr precedent (cdr/SKILL.md does not need this caveat because by the time cdr/SKILL.md shipped, CDR.md was also being shipped in the same cycle).

Frontmatter `calls:` list:

```
calls:
  - CDS.md
  - alpha/SKILL.md
  - beta/SKILL.md
  - gamma/SKILL.md
  - delta/SKILL.md
  - epsilon/SKILL.md
```

Five role overlays + CDS.md. None exist at v0.1 ship time; all are forthcoming per Sub 2 (CDS.md) and Subs 3–5 (overlays) of #403. The advisory nature is acceptable per the cdr-Sub-1 precedent and is explicitly documented in `## v0.1 caveat`.

**Status:** PASS.

### AC5 — Extraction map covers every named software surface

**Oracle:** `docs/extraction-map.md` contains one row per surface in #403's "Source content" list (10 surface groups; the map may expand per-surface sub-lists as multiple rows). Each row's destination resolves under `cnos.cds/`. Each row's migration-sub field names Sub 3, 4, or 5.

**Mechanical check:** `grep -c "^|" docs/extraction-map.md` returns ≥ 10 + 2 (rows + table header + separator).

**Verification:**

```
$ grep -c "^|" src/packages/cnos.cds/docs/extraction-map.md
89
```

89 |-prefixed lines, well above the 12-line floor. Distribution across the 12 surface-group tables (§1–§12) plus the §13 coverage-verification table:

| Surface group (table §) | Source surface | Sub assigned | # rows |
|---|---|---|---|
| §1 Selection function | §Selection (+ §Inputs folded) | Sub 3 | 11 |
| §2 Development lifecycle | §Lifecycle (+ §Roles folded) | Sub 3 | 8 |
| §3 Coordination surfaces | §Tracking | Sub 4 | 4 |
| §4 Artifact contract | §Artifacts | Sub 4 | 9 |
| §5 Mechanical vs judgment | §Mechanical | Sub 5 | 2 |
| §6 Review CLP | §Review | Sub 5 | 2 |
| §7 Gate + closure verification | §Gate | Sub 5 | 2 |
| §8 Assessment + cycle iteration | §Assessment | Sub 5 | 4 |
| §9 Closure | §Closure | Sub 5 | 3 |
| §10 Retro-packaging rule | §Retro-packaging | Sub 5 | 1 |
| §11 Non-goals | §Non-goals | Sub 5 | 1 |
| §12 Large-file authoring rule | §Large-file | Sub 5 | 1 |
| §13 Coverage verification | (cross-reference table) | — | 14 |

Total surface-row count (§1–§12): 48. Plus §13 (14 marker-mapping rows) + per-table header/separator overhead = 89 total |-lines. All 10 #403-named surface groups have dedicated tables; all 14 CDD.md "pending cds extraction" markers are covered (4 of the 14 — §Inputs, §Roles, §Tracking already as its own §3, §Large-file — are folded into related tables with cross-surface origin called out in the row's note column; §13 makes the marker-to-table mapping explicit).

Destination check: spot-read of row destinations confirms all destinations resolve under `cnos.cds/` (e.g. `skills/cds/CDS.md §"Selection function" → "P0 override"`, `skills/cds/lifecycle/SKILL.md`, `skills/cds/<role>/SKILL.md`, etc.). No destination points outside `cnos.cds/`.

Sub-assignment check: every row's migration-sub field names Sub 3, Sub 4, or Sub 5 (e.g. §1 rows = Sub 3; §3–§4 rows = Sub 4; §5–§12 rows = Sub 5). No row references Sub 1, Sub 2, Sub 6, or Sub 7 (those subs have different missions: Sub 1 = this cycle's skeleton; Sub 2 = CDS.md; Sub 6 = marker cleanup; Sub 7 = empirical anchor).

**Status:** PASS.

### AC6 — No content migrated yet

**Oracle:**

- `CDD.md` not modified in this cycle. `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/CDD.md` returns empty.
- `cnos.cds/skills/cds/CDS.md` does not exist.
- No role-overlay SKILL.md files under `skills/cds/{alpha,beta,gamma,delta,epsilon}/`.

**Verification:**

```
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/CDD.md
(empty)
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/ | wc -l
0
```

cnos.cdd is entirely untouched (not just CDD.md — the whole package directory diff is empty).

```
$ test -e src/packages/cnos.cds/skills/cds/CDS.md && echo "EXISTS — FAIL" || echo "DOES NOT EXIST — PASS"
DOES NOT EXIST — PASS
$ ls src/packages/cnos.cds/skills/cds/
SKILL.md  .gitkeep
```

CDS.md does not exist. `skills/cds/` contains only `SKILL.md` and `.gitkeep` — no role-overlay subdirectories.

**Status:** PASS.

### AC7 — Cross-references from cnos.cdr work

**Oracle (informational only per #406):** CDR README and SKILL.md references to CDS by name continue to resolve coherently — `cnos.cds` now exists as a peer package, so the CDR statements are no longer forward-references to a non-existent package.

**Verification:**

```
$ grep -n "cds\|CDS" src/packages/cnos.cdr/skills/cdr/SKILL.md src/packages/cnos.cdr/README.md
```

Per the cdr SKILL.md text read earlier in scaffold step: "CDR is not CDS-with-different-words" appears in the Cross-protocol relationship section. Before this cycle, that statement was a forward-reference to a non-existent peer package. After this cycle, `src/packages/cnos.cds/` exists as a peer of `src/packages/cnos.cdr/`, so the statement resolves to an existing surface that itself reciprocates the distinction (cds/SKILL.md "Cross-protocol relationship" section says "CDS is not CDR-with-different-words" — the reciprocation closes the cross-reference loop).

No edits to cnos.cdr in this sub. AC7 is informational; verification is read-only.

**Status:** PASS (informational).

## β Rule 7 implementation-contract conformance

Per #393 Rule 7, α self-verifies that each axis of the implementation contract is satisfied before β review.

| Axis | Pinned value | Conformance |
|---|---|---|
| Language | Markdown skills + JSON manifest (mirror `cnos.cdr` v0.1 shape) | PASS — only Markdown + JSON; manifest is verbatim from cdr shape (name swap only). |
| CLI integration target | None new; package loads via existing `cn` package-discovery | PASS — no new CLI; `cn build --check` reports `cnos.cds: valid` without code changes. |
| Package scoping | `src/packages/cnos.cds/` (new package, peer of cnos.cdd, cnos.cdr) | PASS — all five new files are under `src/packages/cnos.cds/`; zero files touched outside this directory tree (except the cycle artifacts in `.cdd/unreleased/406/` which are the cycle's own bookkeeping). |
| Existing-binary disposition | N/A; no executables in this sub | PASS — no executables created. |
| Runtime dependencies | None | PASS — no runtime dependencies added. |
| JSON/wire contract | `cn.package.json` schema `cn.package.v1`, version `0.1.0` | PASS — manifest declares exactly these values. |
| Backward compat | `cnos.cdd` not modified; "pending cds extraction" markers stay until Sub 6 | PASS — `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns empty; all 14 markers in CDD.md untouched. |

**Surface containment:** files touched in this cycle (5 in cnos.cds + N cycle artifacts):

1. `src/packages/cnos.cds/cn.package.json` (new)
2. `src/packages/cnos.cds/README.md` (new)
3. `src/packages/cnos.cds/skills/cds/SKILL.md` (new)
4. `src/packages/cnos.cds/skills/cds/.gitkeep` (new)
5. `src/packages/cnos.cds/docs/extraction-map.md` (new)
6. `.cdd/unreleased/406/gamma-scaffold.md`
7. `.cdd/unreleased/406/self-coherence.md` (this file)
8. `.cdd/unreleased/406/beta-review.md` (next)
9. `.cdd/unreleased/406/alpha-closeout.md` (next)
10. `.cdd/unreleased/406/beta-closeout.md` (next)
11. `.cdd/unreleased/406/gamma-closeout.md` (next)
12. `.cdd/unreleased/406/cdd-iteration.md` (next; courtesy stub)
13. `.cdd/iterations/INDEX.md` (next; appended row)

Files NOT touched (per issue Non-goals + design discipline):

- `src/packages/cnos.cdd/**` — out of scope; AC6 verified empty diff.
- `src/packages/cnos.cdr/**` — out of scope; AC7 informational only.
- `schemas/cds/**` — already exists per #388; out of scope per #406 Non-goals.
- `src/packages/cnos.cds/skills/cds/CDS.md` — Sub 2 territory.
- `src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` — Subs 3–5 territory.
- `src/go/**`, `src/ocaml/**`, `bin/**` — no CLI changes; AC2 verified discovery works unchanged.

## Forecasts for β

- **Likely binding finding 0.** All seven ACs pass mechanically; cn discovery verified; cdd untouched; extraction-map coverage explicit at §13.
- **Likely advisory finding 1: SKILL.md adds a `## v0.1 caveat` section beyond cdr's shape.** α defends: the caveat is necessary because cds Sub 1 (this cycle) ships the loader without the canonical contract (CDS.md is Sub 2); cdr Sub 1 shipped both, so cdr/SKILL.md didn't need it. The caveat makes the v0.1 incompleteness explicit and gives readers a navigation pointer (`docs/extraction-map.md`). β may comment; α holds.
- **Likely advisory finding 2: extraction-map's §14 (Open questions) names 5 coordination questions.** β may surface as scope creep. α defends: these are explicitly migration-coordination questions (whether operational realizations also move, not whether destinations are uncertain), and the issue #406 active design constraint says "If a destination is uncertain, name the uncertainty in the row's note column rather than guessing." None of the per-surface table destinations are uncertain; §14 records the broader coordination questions Sub 4/5 will face. Recording them here saves the next cycle's δ from rediscovering them.
- **Likely advisory finding 3: §13 coverage-verification table is technically beyond the strict AC5 wording.** AC5 names the 10-row floor and the destination-resolves-under-cnos.cds check; it does not require a coverage-verification table. α added §13 because the CDD.md marker family has 14 markers (vs #403's 10 surface groups), and folding 4 markers into related tables creates a citation gap that Sub 6 needs to close. The §13 table makes the mapping explicit so Sub 6's marker sweep is mechanical. β may note; α defends as Sub-6-enabling.

No binding findings forecast. β-collapsed review proceeds.
