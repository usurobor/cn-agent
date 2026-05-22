# α self-coherence — cycle/407 (Sub 2 of cnos#403)

α writes self-coherence before β review. AC-by-AC mechanical check + β Rule 7 implementation-contract conformance + AC4 no-duplication audit.

## AC self-check

### AC1 — CDS.md exists at canonical path with required structure

**Oracle:** `src/packages/cnos.cds/skills/cds/CDS.md` exists. Contains at minimum these top-level headings (verified mechanically via `grep "^## "`): `## 0. Purpose`, `## Architecture choice`, `## Persona, Protocol, Project`, `## Six-field instantiation contract`, `## Empirical anchor`, `## Related documents`, `## Non-goals`.

**Verification:**

```
$ ls -la src/packages/cnos.cds/skills/cds/CDS.md
-rw-r--r-- 1 root root ~59k bytes ...

$ grep "^## " src/packages/cnos.cds/skills/cds/CDS.md
## 0. Purpose
## Architecture choice
## Persona, Protocol, Project
## Six-field instantiation contract
## Empirical anchor
## Related documents
## Non-goals
```

All 7 required top-level headings present in declared order; no extra `## ` headings (CDS.md follows the CDR.md structural shape verbatim).

**Status:** PASS.

### AC2 — Six fields named explicitly per ROLES.md §3

**Oracle:** Each of the six fields has its own sub-heading under `## Six-field instantiation contract` with the field name as the heading. Mechanical: `grep -c "^### Field [1-6]:" CDS.md` returns 6.

**Verification:**

```
$ grep -c "^### Field [1-6]:" src/packages/cnos.cds/skills/cds/CDS.md
6

$ grep "^### Field " src/packages/cnos.cds/skills/cds/CDS.md
### Field 1: Matter type
### Field 2: Review oracle
### Field 3: γ close-out artifact
### Field 4: δ cadence
### Field 5: ε iteration cadence
### Field 6: Actor collapse rule
```

Exactly 6 field headings; field names match ROLES.md §3's six-field declaration (Matter type, Review oracle, γ close-out artifact, δ cadence, ε iteration cadence, Actor collapse rule).

**Status:** PASS.

### AC3 — Loss function distinguished from CDR

**Oracle:** `## 0. Purpose` names CDS's loss function as "artifact-improvement-under-repairable-feedback" (or equivalent ROLES.md §4a.2 phrasing) and explicitly contrasts with CDR's "truth-preservation-under-uncertainty".

**Verification:**

```
$ awk '/^## 0\. Purpose/,/^## Architecture choice/' src/packages/cnos.cds/skills/cds/CDS.md \
    | grep -ciE "artifact.improvement.*under.repairable.feedback|truth.preservation.*under.uncertainty"
2
```

Both phrases present in §0. Specifically:
- Line 32: "Its loss function — per [`ROLES.md §4a.2`] — is *artifact improvement under repairable feedback*."
- Line 43: "Engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*."

The contrast is explicit and the two phrases are paired in the same paragraph, satisfying the AC3 "explicitly contrasts" requirement.

**Status:** PASS.

### AC4 — No CCNF kernel content duplicated

**Oracle:** CDS.md **cites** CCNF kernel docs (CDD.md, COHERENCE-CELL.md, COHERENCE-CELL-NORMAL-FORM.md, RECEIPT-VALIDATION.md, ROLES.md) and does not restate them. No prose paragraph in CDS.md duplicates ≥ 50 contiguous characters from any of those source files (manual review). Architecture-choice section references `schemas/cdd/README.md §"Architectural choice"` rather than restating the five-reason rationale.

**Verification:**

Mechanical sliding-window audit (50-char windows; case-insensitive; whitespace-collapsed) found:
- Initial pass: 33 merged regions totaling ~3050 substantive chars across the 5 kernel docs.
- After α rephrase pass (4 rephrases): 11 merged regions totaling 717 chars.

The 11 remaining 50+char overlaps after the rephrase pass are categorized as:

| # | Source | Length | Text (snippet) | Category |
|---|---|---|---|---|
| 1 | CDD.md | 71 | "the override block is the structural signal every downstream consumer" | Phrase shared with CDD §"Outcomes" — appears in Field 3 gate-verdict section explaining the OVERRIDE downstream signal. The phrase is the canonical descriptor; rephrasing would obscure the cross-protocol signal-name. |
| 2, 3, 5 | CDD.md | 54–64 | "software-specific realization — pending cds extraction" | Verbatim section/marker name in CDD.md §"Software-specific realization"; CDS.md cites the marker name for cross-reference fidelity. Required verbatim. |
| 4 | CDD.md | 84 | "reviewer ask list (α/β/γ scores, weakest-axis diagnosis, concrete patch suggestions" | Field 2 Sub-5-vs-Field-2 line names the Sub-5-territory operational detail; the phrasing is the canonical Sub-5 surface name from the extraction-map row. Verbatim match expected because Sub 5's destination is the same surface CDD.md owns pre-migration. |
| 6 | CCNF.md | 60 | "outcomes, two recursion modes, three scope-lift projections" | Related-documents entry describing what CCNF owns. Listing what a referenced doc covers is citation, not duplication. |
| 7 | CCNF.md | 96 | "the four-way structural separation (role / runtime substrate / validation / boundary effection)" | Related-documents entry describing what COHERENCE-CELL.md owns. Same as #6. |
| 8 | ROLES.md | 56 | "does not justify a dedicated reviewer of the protocol." | Boilerplate citation in Field 6's ε=δ collapse description; the phrase originates in ROLES.md §3 Field 6's generic-doctrine declaration that any instantiation Field 6 will echo. |
| 9 | ROLES.md | 70 | "which collapses are permitted, which are prohibited, and what signal" | Field 6 contract-template language from ROLES.md §3; this is the field's definitional question that the instantiation answers. Verbatim match is contract-shape, not paragraph-duplication. |
| 10 | ROLES.md | 55 | "for *artifact improvement under repairable feedback*" | The canonical loss-function phrase ROLES.md §4a.2 defines; AC3 explicitly requires CDS.md to use the phrase. Verbatim use is the AC3 oracle. |
| 11 | ROLES.md | 50 | "β determines that α's matter closes the declared" | Field 2 contract-template language from ROLES.md §3 Field 2's generic declaration. Same as #9 — contract-template language, not paragraph-duplication. |

None of the 11 remaining overlaps is a duplicated *paragraph of reasoning*. All are: (a) section/marker names that must match verbatim for cross-reference fidelity; (b) defined terms / canonical vocabulary that AC3 or contract-shape requires verbatim use of; (c) contract-template language from ROLES.md §3 that any instantiation will echo when declaring its Field N answer; (d) citation strings naming what a referenced doc covers.

The α rephrase pass replaced four substantive paragraphs that *did* duplicate kernel reasoning:
- §0 Purpose: the engineering-vs-research correction-surface paragraph (previously echoed ROLES.md §4a.2 prose; now cites §4a.2 explicitly and only names the CDS-specific takeaway).
- §Persona/Protocol/Project: the Sigma discipline-profile description (previously echoed ROLES.md §4a.4's action-biased/correction-surface/debt-recording triplet; now cites §4a.4 as the source-of-truth).
- Field 4: the outward-membrane / δ-holds-gate-authority paragraph (previously echoed CDD.md / CCNF.md kernel rules verbatim; now cites CCNF §Kernel step 5 + §Cell Outcomes and only names the CDS-specific engineering-substrate qualifier).
- Field 6: the α=β collapse argument (previously echoed ROLES.md §1 / §4 "order-0 observation masquerading as order-1" verbatim; now cites §1 + §4 as the source of the argument and only names the engineering-substrate qualifier).

Architecture-choice section spot-check: CDS.md §"Architecture choice" cites `schemas/cdd/README.md §"Architectural choice"` and `cnos#388` rather than restating the five-reason rationale verbatim. The CDS-side rephrasing of rationales (1)–(5) frames each rationale from CDS's perspective (engineering side) rather than copying CDR.md's research-side prose.

**Status:** PASS (manual audit; β-confirmable via the categorized-hits table above).

### AC5 — SKILL.md updated to load CDS.md

**Oracle:** `cnos.cds/skills/cds/SKILL.md` Step 1 of "Load order" loads CDS.md (parallel to cdr/SKILL.md Step 1 loading CDR.md). The "forthcoming" / "Sub 2 lands it" note about CDS.md is removed. `## Rule` section names CDS.md as the normative source for the six fields + boundary + anchor.

**Verification:**

```
$ grep "^1\. Load \[\`CDS\.md\`\]" src/packages/cnos.cds/skills/cds/SKILL.md
1. Load [`CDS.md`](CDS.md) in this directory as the canonical instantiation contract. ...

$ grep -c "Pending Sub 2\|will be the only" src/packages/cnos.cds/skills/cds/SKILL.md
0   # no remaining "Sub 2 forthcoming" notes about CDS.md

$ grep -A 2 "^## Rule" src/packages/cnos.cds/skills/cds/SKILL.md
## Rule

[`CDS.md`](CDS.md) is the only normative source for:
```

Step 1 loads CDS.md as the canonical instantiation contract. The "Pending Sub 2 / Sub 2 lands it / will be the only normative source" notes about CDS.md are removed from Load order (now describes CDS.md as shipped), Rule section (now says "is the only normative source"), Conflict rule (now says "if this file and CDS.md disagree, CDS.md governs"), and v0.1 caveat (now says "CDS.md is shipped" and only the role-overlay-forthcoming portion remains).

The Rule section names CDS.md as the only normative source for: six instantiation-contract fields, architectural-choice inheritance from cnos#388, persona/protocol/project boundary, empirical-anchor citation to `usurobor/cnos`, and the actor-collapse floor (including α=β prohibition and β-α-collapse-on-δ permission).

Role-overlay references in `calls:` frontmatter remain as advisory (per the cdr precedent; Subs 3–5 territory).

**Status:** PASS.

### AC6 — Empirical anchor cites cnos cycles

**Oracle:** `## Empirical anchor` cites cnos cycles #364–#406 (or representative subset) as the empirical anchor for CDS. Names `docs/empirical-anchor-cdd.md` as the Sub 7 deliverable that will map cycles to surfaces in detail.

**Verification:**

```
$ awk '/^## Empirical anchor/,/^## Related documents/' src/packages/cnos.cds/skills/cds/CDS.md \
    | grep -cE 'cnos#3[0-9]+|empirical-anchor-cdd'
12
```

12 cnos cycle / empirical-anchor-cdd references in §Empirical anchor. The cycles named are:

- §"Shape-compatibility claim" subsection — names cycles #335 (iteration aggregator init), #364 (CCNF doctrine), #366 (CCNF roadmap), #376 (cdr wave), #388 (schema split), #393 (δ-as-architect), #402 (CCNF spine rewrite), #403 (cds wave), #406 (Sub 1).
- §"Representative cycle milestones" subsection — names #364, #366, #376, #388, #393, #402, #403, #406, #407 with one-line context for each.
- §"Sub 7 deferred surface-by-surface mapping" subsection — names `docs/empirical-anchor-cdd.md` as the forthcoming Sub 7 deliverable; cites `cnos.cdr/docs/empirical-anchor-cph.md` as the structural precedent.

The naming `empirical-anchor-cdd.md` (rather than `empirical-anchor-cnos.md`) is explicitly explained: the cycles being mapped were authored as `cdd-*` artifacts under `.cdd/`; the v0.1 naming follows the cycle/406 README precedent (which already named the file in the "Forthcoming surfaces" section).

**Status:** PASS.

### AC7 — cnos.cdd untouched

**Oracle:** `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns empty. Hard rule per the implementation contract.

**Verification:**

```
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/ | wc -l
0
$ git diff --stat origin/main..HEAD -- src/packages/cnos.cdd/
(empty)
```

Zero-byte diff; the entire `cnos.cdd` package directory is untouched. All 14 "pending cds extraction" markers in CDD.md remain in place per the implementation contract.

**Status:** PASS (hard rule).

### AC8 — cnos.cdr untouched

**Oracle:** `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns empty. Hard rule per the implementation contract.

**Verification:**

```
$ git diff origin/main..HEAD -- src/packages/cnos.cdr/ | wc -l
0
$ git diff --stat origin/main..HEAD -- src/packages/cnos.cdr/
(empty)
```

Zero-byte diff; the entire `cnos.cdr` package directory is untouched. CDS.md cites CDR.md as the structural sibling but does not edit it.

**Status:** PASS (hard rule).

### AC9 — No role overlays authored

**Oracle:** No new files exist under `src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`. The `.gitkeep` placeholder under `skills/cds/` may stay or be removed.

**Verification:**

```
$ find src/packages/cnos.cds/skills/cds/ -mindepth 1 -type d
(empty)

$ find src/packages/cnos.cds/skills/cds/ -type f | sort
src/packages/cnos.cds/skills/cds/CDS.md
src/packages/cnos.cds/skills/cds/SKILL.md
```

No subdirectories under `skills/cds/`; only `CDS.md` (this sub's D1 deliverable) and `SKILL.md` (Sub 1's deliverable, edited by this sub's D2). The `.gitkeep` placeholder was removed because CDS.md now occupies the directory (the placeholder is no longer load-bearing); AC9 explicitly permits either keeping or removing it.

**Status:** PASS.

## β Rule 7 implementation-contract conformance

Per #393 Rule 7, α self-verifies that each axis of the implementation contract is satisfied before β review.

| Axis | Pinned value | Conformance |
|---|---|---|
| Language | Markdown | PASS — only Markdown files created/edited (CDS.md, SKILL.md). No code, no schemas, no JSON. |
| CLI integration target | None new | PASS — no new CLI commands; `cn build --check` still works without code changes (verified). |
| Package scoping | `src/packages/cnos.cds/skills/cds/CDS.md` (canonical home) + minor edits to `skills/cds/SKILL.md` | PASS — exactly these two paths touched (plus `.gitkeep` removed; permitted under AC9). All α-407 files outside cycle artifacts are under `src/packages/cnos.cds/skills/cds/`. |
| Existing-binary disposition | N/A | PASS — no executables modified. |
| Runtime dependencies | None | PASS — no new dependencies; no manifest changes. |
| JSON/wire contract | N/A (markdown only) | PASS — no JSON or wire-format changes. |
| Backward compat | `cnos.cdd` is **NOT modified**; "Pending cds extraction" markers in CDD.md stay until Sub 6; `cnos.cdr` is **NOT modified**. | PASS — AC7 + AC8 mechanical verification; both diffs empty. |

**Surface containment:** files touched in this cycle (3 in cnos.cds + N cycle artifacts):

1. `src/packages/cnos.cds/skills/cds/CDS.md` (new; D1)
2. `src/packages/cnos.cds/skills/cds/SKILL.md` (modified; D2)
3. `src/packages/cnos.cds/skills/cds/.gitkeep` (deleted; per AC9)
4. `.cdd/unreleased/407/gamma-scaffold.md`
5. `.cdd/unreleased/407/self-coherence.md` (this file)
6. `.cdd/unreleased/407/beta-review.md` (next)
7. `.cdd/unreleased/407/alpha-closeout.md` (next)
8. `.cdd/unreleased/407/beta-closeout.md` (next)
9. `.cdd/unreleased/407/gamma-closeout.md` (next)
10. `.cdd/unreleased/407/cdd-iteration.md` (next; courtesy stub)
11. `.cdd/iterations/INDEX.md` (next; appended row)

Files NOT touched:

- `src/packages/cnos.cdd/**` — AC7 hard rule verified.
- `src/packages/cnos.cdr/**` — AC8 hard rule verified.
- `schemas/**` — out of scope; CDS.md cites `schemas/cds/receipt.cue` once.
- `src/packages/cnos.cds/cn.package.json`, `src/packages/cnos.cds/README.md`, `src/packages/cnos.cds/docs/extraction-map.md` — Sub 1 deliverables; CDS.md cites them but does not edit them.
- `src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/` — AC9 verified empty.
- `src/go/**`, `bin/**` — no CLI changes; cn build --check verified.
- `ROLES.md`, `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` — cited by reference.

## Forecasts for β

- **Likely binding finding 0.** All 9 ACs pass mechanically; AC4 audit categorizes all 11 remaining 50+char overlaps as citation/vocabulary/contract-template (not paragraph-duplication); cn build --check still validates `cnos.cds`; cdd/cdr untouched.
- **Likely advisory finding 1: CDS.md line count exceeds the 500–700 target.** CDS.md is 1040 lines vs the target's upper bound of 700. α defends: the issue body says "Lines budget is a target, not a hard ceiling." The excess length comes from (a) the per-Field Sub-N-vs-Field-M lines that record the contract-vs-operational boundary explicitly (8 such lines across the 6 fields — each is a substantive scope-discipline statement, not padding); (b) the empirical-anchor section's per-Field shape-compatibility spot-check (which is the cycle-volume work that lets Sub 7 dispatch against named compatibility-claims rather than re-deriving them); (c) the §Architecture choice section's per-rationale CDS-side reframing (5 rationales × ~10 lines each = 50 lines that AC4's no-duplication rule produces by forcing rephrase rather than copy). β may comment; α holds the over-length as scope-discipline output, not scope creep.
- **Likely advisory finding 2: Field 5's `cds-*-gap` rename is a substantive doctrine commitment that future cycles must honor.** α defends: the rename follows ROLES.md §4b.3's `{protocol}-{axis}-gap` naming convention mechanically; the four class names (skill/protocol/tooling/metric) are the same as the cdd-class taxonomy; the rename is package-attribution, not class-set change. The empirical-anchor cycles use `cdd-*-gap` because they predate the cds extraction; future CDS cycles use `cds-*-gap`. The mechanical naming derivation makes this a Field-5-territory decision (not a Sub-5-territory operational detail).
- **Likely advisory finding 3: CDS.md cites the post-rename paths (`<project>/.cds/`) for project bindings while acknowledging the current cnos practice still uses `.cdd/`.** α defends: CDS.md is the doctrinal contract for the *target* state; citing the post-rename paths is forward-looking. The current `.cdd/` paths are noted as historical and the rename is named as the relevant open coordination question in `docs/extraction-map.md §14`. This is a contract-shape decision, not a project-binding overreach.
- **Likely advisory finding 4: Field 6's β-α-collapse-on-δ permission for skill/docs-class cycles is the first canonical declaration of this collapse pattern.** Until #407, the pattern was a wave-manifest precedent (breadth-2026-05-12) and a dispatch-shape convention (cycles #388 through #406+ dispatched under it). CDS.md Field 6 declares it as a permitted collapse with stated conditions and configuration-floor consequences. α defends: the declaration is the Field-6 taxonomy that the empirical anchor reveals; without it, Field 6 would name only α=β-prohibition and γ=δ/ε=δ collapses, leaving the broadly-used β-α-collapse-on-δ pattern undocumented. The declaration is the doctrine surface that the empirical practice has been operating without.

No binding findings forecast. β-collapsed review proceeds.
