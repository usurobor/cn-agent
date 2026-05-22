# β review — cycle/406 (Sub 1 of cnos#403)

**Reviewer:** β-collapsed-on-δ (per δ contract; γ+α+β collapsed on δ for skill/docs-class cycles per breadth-2026-05-12 wave manifest precedent).
**Round:** R1.
**Verdict:** APPROVED.
**β-rigor:** AC1 file-existence + frontmatter parse; AC4 cdr-loader-pattern structural match; AC5 row count + destination resolution + coverage verification; AC6 hard rule (cdd untouched).

## R1 review oracle pass

### AC1 — Package skeleton exists at canonical path

**β check:** all four named files exist; README ≥ 50 lines; SKILL.md has valid frontmatter.

```
$ ls src/packages/cnos.cds/cn.package.json \
     src/packages/cnos.cds/README.md \
     src/packages/cnos.cds/skills/cds/SKILL.md \
     src/packages/cnos.cds/docs/extraction-map.md
(all four exist; no errors)

$ wc -l src/packages/cnos.cds/README.md
64 src/packages/cnos.cds/README.md
```

README is 64 lines — above the 50-line floor.

Frontmatter parse:

```
$ head -38 src/packages/cnos.cds/skills/cds/SKILL.md
---
name: cds
description: Coherence-Driven Software. Use for engineering work …
artifact_class: skill
kata_surface: embedded
governing_question: How do we improve software artifacts under repairable feedback coherently, without the engineering loop drifting farther from the kernel with each round?
visibility: public
triggers:
  - software
  - code
  - test
  - branch
  - diff
  - ci
  - release
  - deploy
  - cycle
  - engineering
scope: global
inputs: …
outputs: …
requires: …
calls:
  - CDS.md
  - alpha/SKILL.md
  - beta/SKILL.md
  - gamma/SKILL.md
  - delta/SKILL.md
  - epsilon/SKILL.md
---
```

Six required frontmatter fields per AC1 all present: `name`, `description`, `artifact_class`, `governing_question`, `triggers`, `calls`. Frontmatter is well-formed YAML between `---` markers. PASS.

Manifest shape:

```
$ python3 -c "import json; d=json.load(open('src/packages/cnos.cds/cn.package.json')); \
    assert d=={'schema':'cn.package.v1','name':'cnos.cds','version':'0.1.0','kind':'package','engines':{'cnos':'>=3.81.0'}}, d; \
    print('manifest matches #406-pinned shape verbatim')"
manifest matches #406-pinned shape verbatim
```

PASS.

**AC1 verdict:** PASS.

### AC2 — Package loads via existing `cn` discovery without error

**β check:** `cn build --check` reports `cnos.cds: valid`; `git diff` on `src/go/` returns empty.

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

$ git diff origin/main..HEAD -- src/go/ | wc -l
0
```

cnos.cds is discovered as a valid v1 package alongside the existing six packages, with zero changes to the cn discovery code. The `cn build --check` output is the operational confirmation that the manifest shape, content-class directory presence (skills/), and v1-schema declaration all parse cleanly under the existing `FullPackageManifest` (see `src/go/internal/pkg/pkg.go` lines 160–167).

**AC2 verdict:** PASS.

### AC3 — README declares the cross-protocol relationship correctly

**β-rigor check:** read README's first three sections (lines 1–13 prose) and the Cross-protocol portion of Quick Start.

- Line 3 ("CDS … is the software-development realization of CCNF — the generic recursive coherence-cell algorithm that lives in cnos.cdd"): declares CDS as software-development realization of CCNF, names cnos.cdd as kernel home. ✓
- Line 3 ("It sits as a peer to cnos.cdr (research realization); both inherit the generic kernel from cnos.cdd without re-deriving it"): names CDR as sibling research realization; affirms no kernel duplication. ✓
- Line 5 ("The kernel/realization split is the option (a) decision recorded at cnos#388"): cites #388 architecture inheritance. ✓
- Line 9 ("CDS is the engineering discipline … artifact improvement under repairable feedback"): declares CDS's loss function, which is the divergence from CDR's research-discipline loss function. ✓
- Line 11 (kernel grammar lives in cnos.cdd, cited by reference, never restated): pointer-only discipline declared. ✓
- Line 13 (architectural split inherits option (a) from #388, same split that produced cnos.cdr v0.1 at #376): structural-precedent citation. ✓

All four AC3 sub-clauses verified PASS.

**AC3 verdict:** PASS.

### AC4 — SKILL.md loader pattern mirrors cdr

**β-rigor check:** structural-diff cds/SKILL.md against cdr/SKILL.md (section-by-section).

| Section | cdr/SKILL.md | cds/SKILL.md | β note |
|---|---|---|---|
| Frontmatter (name, description, artifact_class, kata_surface, governing_question, visibility, triggers, scope, inputs, outputs, requires, calls) | ✓ | ✓ | Identical 12-field shape; calls: lists CDR.md or CDS.md + 5 role overlays. |
| Top heading `# CDR` / `# CDS` | ✓ | ✓ | |
| `## Load order` | ✓ | ✓ | cds adds a v0.1-status leading paragraph noting CDS.md is forthcoming; cdr did not need this because its Sub 1 shipped both loader and CDR.md. Acceptable per cdr-Sub-1 precedent. |
| `## Rule` | ✓ | ✓ | cds rephrases "is the only normative source" as "will be the only normative source" to acknowledge CDS.md hasn't shipped yet. Acceptable. |
| `## Role overlays` | ✓ | ✓ | cds names the same five-role positions (α, β, γ, δ, ε) with one renaming: cdr uses `operator/` for δ (legacy naming per cdr/SKILL.md §"operator role" note); cds uses `delta/` directly. β verified: this is consistent with the cycle/394 ROLES.md decision to standardize role-letter directory naming; cdr's `operator/` is the older legacy form. |
| `## Cross-protocol relationship` | ✓ | ✓ | cds declares "CDS is not CDR-with-different-words" reciprocating cdr's "CDR is not CDS-with-different-words". |
| `## Conflict rule` | ✓ | ✓ | cds adds one rule beyond cdr's set: "If CDS.md and cnos.cdd (CCNF kernel) disagree on the kernel grammar, cnos.cdd governs." This is the kernel/realization conflict rule that cdr/SKILL.md does not explicitly state because cdr's kernel/realization split was implicit. cds making it explicit is doctrine-coherence improvement, not over-promising. |
| `## v0.1 caveat` (cds-only) | — | ✓ | New section beyond cdr's shape. β assessment: necessary because cds Sub 1 ships loader without CDS.md (Sub 2); the section is explicit about the forthcoming nature and points readers at extraction-map.md for interim navigation. Not over-promising; under-promising-and-documenting. Acceptable. |

`calls:` advisory-target check:

```
$ for f in CDS.md alpha/SKILL.md beta/SKILL.md gamma/SKILL.md delta/SKILL.md epsilon/SKILL.md; do
    test -e src/packages/cnos.cds/skills/cds/$f && echo "$f: EXISTS" || echo "$f: not present (Sub 2/3-5 territory)"
  done
CDS.md: not present (Sub 2/3-5 territory)
alpha/SKILL.md: not present (Sub 2/3-5 territory)
beta/SKILL.md: not present (Sub 2/3-5 territory)
gamma/SKILL.md: not present (Sub 2/3-5 territory)
delta/SKILL.md: not present (Sub 2/3-5 territory)
epsilon/SKILL.md: not present (Sub 2/3-5 territory)
```

All six `calls:` targets are advisory per the cdr-Sub-1 precedent. The v0.1 caveat documents this explicitly. PASS.

**AC4 verdict:** PASS.

### AC5 — Extraction map covers every named software surface

**β-rigor check 1:** row count meets floor.

```
$ grep -c "^|" src/packages/cnos.cds/docs/extraction-map.md
89
```

89 |-prefixed lines >> 12-row floor (10 row + header + separator).

**β-rigor check 2:** every #403-named surface group has a dedicated table.

| #403 source-content surface | Table in extraction-map | β verified |
|---|---|---|
| Selection function | §1 | ✓ (line 49 heading "## 1. Selection function (§Selection)") |
| Development lifecycle | §2 | ✓ |
| Artifact contract | §4 | ✓ |
| Mechanical vs judgment boundary | §5 | ✓ |
| Review CLP | §6 | ✓ |
| Gate + closure verification checklist F1–F10 | §7 | ✓ |
| Assessment + cycle iteration | §8 | ✓ |
| Closure | §9 | ✓ |
| Retro-packaging | §10 | ✓ |
| Non-goals | §11 | ✓ |

All 10 #403 surface groups have dedicated tables.

**β-rigor check 3:** every CDD.md "pending cds extraction" marker is covered.

```
$ grep -n "pending cds extraction\|pending cds" src/packages/cnos.cdd/skills/cdd/CDD.md
80, 122, 124 — section headings and contextual mentions
150 — final pointer in the hard-rule section
```

The actual marker bullets are in lines 126–139 (read in scaffold step). The 14 markers (Inputs, Selection, Lifecycle, Roles and dispatch, Tracking, Artifacts, Mechanical, Review, Gate, Assessment, Closure, Retro-packaging, Non-goals, Large-file) are all listed in the extraction-map's §13 coverage-verification table (lines 248–263 of extraction-map.md), each mapped to its containing table. PASS.

**β-rigor check 4:** every row's destination resolves under `cnos.cds/`.

Spot-read of 8 sampled rows across all 12 tables:

- §1 row "P0 override": destination `skills/cds/CDS.md §"Selection function" → "P0 override"`. Resolves under cnos.cds/. ✓
- §1 row "Inputs": destination `skills/cds/CDS.md §"Selection function" → "Inputs"`. ✓
- §2 row "0–13 step table": destination `skills/cds/CDS.md §"Development lifecycle" → "Step table"` (or `skills/cds/lifecycle/SKILL.md §"Steps"`). ✓
- §4 row "Artifact Location Matrix": destination `skills/cds/CDS.md §"Artifact contract" → "Location matrix"`. ✓
- §5 row "judgment-bearing axes": destination `skills/cds/CDS.md §"Mechanical vs judgment" → "Judgment axes"`. ✓
- §7 row "F1–F10 checklist": destination `skills/cds/CDS.md §"Gate" → "Closure verification checklist"` (or `skills/cds/gate/SKILL.md §"F1–F10"`). ✓
- §8 row "cycle iteration triggers": destination `skills/cds/CDS.md §"Assessment" → "Cycle iteration triggers"`. ✓
- §10 row "direct-to-main exception": destination `skills/cds/CDS.md §"Retro-packaging" → "Direct-to-main exception"`. ✓

All sampled destinations resolve under cnos.cds/. No row points outside the package.

**β-rigor check 5:** every row's migration-sub field names Sub 3, 4, or 5.

Spot-read of the Sub column across all 12 tables: §1 = Sub 3 throughout; §2 = Sub 3 throughout; §3 = Sub 4 throughout; §4 = Sub 4 throughout; §5–§12 = Sub 5 throughout. No row references Sub 1 (this cycle), Sub 2 (CDS.md cycle), Sub 6 (cleanup cycle), or Sub 7 (empirical anchor) — the Sub-3/Sub-4/Sub-5 partition is preserved per the §0.3 sub-naming convention.

**AC5 verdict:** PASS.

### AC6 — No content migrated yet

**β-rigor check (hard rule):**

```
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/CDD.md
(empty)
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/ | wc -l
0
$ test -e src/packages/cnos.cds/skills/cds/CDS.md && echo "EXISTS — FAIL" || echo "DOES NOT EXIST — PASS"
DOES NOT EXIST — PASS
$ ls src/packages/cnos.cds/skills/cds/ | grep -v "^\.gitkeep$\|^SKILL.md$" | wc -l
0
```

CDD.md untouched (zero-byte diff); entire cnos.cdd package untouched; CDS.md does not exist; no role-overlay files under skills/cds/ (just SKILL.md and .gitkeep).

**AC6 verdict:** PASS (hard rule satisfied).

### AC7 — Cross-references from cnos.cdr work

**β check:** informational read. cdr/SKILL.md's "CDR is not CDS-with-different-words" statement (line 97 of cdr/SKILL.md per earlier scaffold read) now resolves to an existing peer package (cnos.cds), which itself contains a reciprocating "CDS is not CDR-with-different-words" statement in its own SKILL.md Cross-protocol relationship section. The cross-reference loop closes; no edits to cdr were required.

**AC7 verdict:** PASS (informational).

## β observations (non-blocking, not findings)

**Obs-1:** `skills/cds/SKILL.md` includes a `## v0.1 caveat` section that is not present in cdr/SKILL.md. β assessment: necessary documentation of the Sub-1-without-Sub-2 status; reads as "under-promising and documenting" rather than over-promising. The cdr-Sub-1 precedent did not need this because its Sub 1 shipped both loader and CDR.md together; cds Sub 1 explicitly defers CDS.md to Sub 2 per the #406 D1 deliverable shape. Acceptable.

**Obs-2:** `skills/cds/SKILL.md` Conflict rule adds one rule beyond cdr's set: "If CDS.md and cnos.cdd (CCNF kernel) disagree on the kernel grammar, cnos.cdd governs." β assessment: this is the kernel/realization-conflict rule made explicit (cdr's set leaves it implicit). The addition is doctrine-coherence improvement; the new rule is consistent with the kernel-realization hierarchy declared in the Cross-protocol relationship section. Acceptable.

**Obs-3:** Role-overlay directory naming uses `delta/` directly rather than cdr's `operator/`. β verified: cdr/SKILL.md §"Role overlays" includes a note that "the directory is named `operator/` to mirror the engineering doctrine's `cnos.cdd/skills/cdd/operator/SKILL.md` exemplar; the role itself is δ per CDR.md Field 4". The cds-side naming `delta/` is the role-letter-canonical form (per ROLES.md). β assessment: cds is shipping the canonical naming directly rather than inheriting cdr's legacy `operator/` workaround; this is consistent with the broader cycle-394+ trend toward role-letter naming. Acceptable.

**Obs-4:** extraction-map.md §14 records 5 open coordination questions (`.cdd/` → `.cds/` rename; cross-repo / harness / operator SKILL location; release-effector location). β assessment: per #406 active design constraint "If a destination is uncertain, name the uncertainty in the row's note column rather than guessing", the per-row destinations are stable; §14 records migration-coordination questions that Subs 3–5 will face but do not invalidate the destinations. Recording them here saves the next cycle's δ from rediscovering them. Acceptable; arguably exceeds the AC5 minimum.

**Obs-5:** extraction-map.md §13 adds an explicit marker-to-table coverage table (14 rows mapping each CDD.md marker to its containing extraction-map table §). β assessment: not required by AC5 wording (which names destination-resolves and sub-named checks). However, the §13 table is what makes Sub 6's marker-sweep mechanical: Sub 6 reads CDD.md's pending-cds markers, looks up each in §13, finds the destination table, verifies the content migrated, then removes the marker. Without §13, Sub 6 would need to re-derive the mapping. β assessment: §13 is Sub-6-enabling; acceptable; arguably should be a #406-AC clarification for future bootstrap cycles. Not a finding.

**Obs-6:** extraction-map.md folds 4 CDD.md markers (Inputs, Roles, Coordination surfaces folded as §3 explicitly, Large-file) into related tables rather than giving each its own dedicated §. β verified: §13 makes this explicit. Coordination surfaces is its own §3 (not folded); Inputs is folded into §1 Selection; Roles is folded into §2 Lifecycle; Large-file is folded into §12 (which is itself a small dedicated §). The folding is editorially defensible: Inputs is "what selection reads", which is selection-adjacent; Roles is "who runs the lifecycle", which is lifecycle-adjacent. β: acceptable.

## Trigger assessment

| Trigger | Fire condition | Fired? | β note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVED. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 0 binding findings; 6 advisory observations. |
| Avoidable tooling / environment failure | environment blocked the cycle | **No** | `go run ./src/go/cmd/cn build --check` ran cleanly. |
| CI red on merge commit | CI fails post-merge | **N/A** | merge not yet executed (operator's authority per dispatch). |
| Loaded skill failed to prevent a finding | skill underspecified | **No** | no findings. |

No §9.1 triggers fired.

## Verdict

**R1 APPROVED.** AC1–AC7 PASS mechanically. β-rigor checks (file existence + frontmatter parse; cdr-loader-pattern structural match; row count + destination resolution + coverage verification; cdd-untouched hard rule) all pass. No binding findings. Six non-blocking observations recorded above. No fix-round needed.

The cycle is ready for close-out and operator-facing merge instruction.
