# β review — cycle/407 (Sub 2 of cnos#403)

**Reviewer:** β-collapsed-on-δ (per δ contract; γ+α+β collapsed on δ for skill/docs-class cycles per breadth-2026-05-12 wave manifest precedent).
**Round:** R1.
**Verdict:** APPROVED.
**β-rigor:** AC1–AC9 mechanical checks + AC4 no-duplication audit (highest-rigor focus per #407 active design constraints) + doctrine-coherence-with-cdr-template structural diff + name-overpromise audit on Sub 3/4/5 references.

## R1 review oracle pass

### AC1 — Top-level structure

**β check:** `grep "^## " CDS.md` returns the seven required headings in declared order; no extra `## ` headings.

```
$ grep "^## " src/packages/cnos.cds/skills/cds/CDS.md
## 0. Purpose
## Architecture choice
## Persona, Protocol, Project
## Six-field instantiation contract
## Empirical anchor
## Related documents
## Non-goals
```

Mirrors CDR.md's section structure verbatim. PASS.

**AC1 verdict:** PASS.

### AC2 — Six fields named explicitly

**β check:** `grep -c "^### Field [1-6]:" CDS.md` returns 6.

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

Six fields; field names match ROLES.md §3's declaration. Heading form (`### Field N: <name>`) is the exact form #407 AC2 specified.

**AC2 verdict:** PASS.

### AC3 — Loss function distinguished from CDR

**β check:** §0 Purpose names both phrases.

```
$ awk '/^## 0\. Purpose/,/^## Architecture choice/' src/packages/cnos.cds/skills/cds/CDS.md \
    | grep -ciE "artifact.improvement.*under.repairable.feedback|truth.preservation.*under.uncertainty"
2
```

Both phrases present. β-rigor read of §0:
- Line 32: "Its loss function — per [`ROLES.md §4a.2`] — is *artifact improvement under repairable feedback*."
- Line 43: "Engineering optimises for *artifact improvement under repairable feedback*; research optimises for *truth preservation under uncertainty*."

The phrasing matches ROLES.md §4a.2 verbatim for the engineering-loss-function side (required per AC3 oracle); the contrast with CDR is explicit and paired in the same paragraph.

**AC3 verdict:** PASS.

### AC4 — No CCNF kernel content duplicated

**β-rigor check (load-bearing per #407 active design constraints).** Mechanical sliding-window scan (50-char windows, case-insensitive, whitespace-collapsed) against CDD.md, COHERENCE-CELL.md, COHERENCE-CELL-NORMAL-FORM.md, RECEIPT-VALIDATION.md, ROLES.md.

**Initial scan result (pre-α-rephrase):** 33 merged regions totaling ~3050 substantive chars.

**Final scan result (post-α-rephrase):** 11 merged regions totaling 717 chars. β review per-region:

| # | Source | Length | Snippet (lowercased, whitespace-collapsed) | β verdict |
|---|---|---|---|---|
| 1 | CDD.md | 71 | "the override block is the structural signal every downstream consumer" | **Acceptable** — descriptor of the OVERRIDE downstream signal; same wording in CDD.md §"Outcomes" and CDS.md Field 3 gate-verdict section. Rephrasing would obscure cross-protocol signal-name consistency. |
| 2, 3, 5 | CDD.md | 54–64 | "software-specific realization — pending cds extraction" | **Acceptable** — verbatim section/marker name. Cross-reference fidelity requires verbatim quote. |
| 4 | CDD.md | 84 | "reviewer ask list (α/β/γ scores, weakest-axis diagnosis, concrete patch suggestions" | **Acceptable** — Field 2 Sub-5-vs-Field-2 line names the Sub-5-territory operational detail using the canonical surface-name from the extraction-map row. |
| 6 | CCNF.md | 60 | "outcomes, two recursion modes, three scope-lift projections" | **Acceptable** — Related-documents entry describing what CCNF owns. Citation, not duplication. |
| 7 | CCNF.md | 96 | "the four-way structural separation (role / runtime substrate / validation / boundary effection)" | **Acceptable** — Related-documents entry describing what COHERENCE-CELL.md owns. Citation. |
| 8 | ROLES.md | 56 | "does not justify a dedicated reviewer of the protocol." | **Acceptable** — Field 6 ε=δ collapse description; the phrase originates in ROLES.md §3 Field 6's generic-doctrine declaration; instantiation Field 6 will echo. |
| 9 | ROLES.md | 70 | "which collapses are permitted, which are prohibited, and what signal" | **Acceptable** — Field 6 contract-template language from ROLES.md §3. Contract-shape, not paragraph-duplication. |
| 10 | ROLES.md | 55 | "for *artifact improvement under repairable feedback*" | **Required** — the canonical loss-function phrase; AC3 explicitly requires verbatim use. |
| 11 | ROLES.md | 50 | "β determines that α's matter closes the declared" | **Acceptable** — Field 2 contract-template language from ROLES.md §3. Contract-shape. |

**β verdict on AC4:** PASS. None of the 11 remaining overlaps is a paragraph of duplicated kernel reasoning. All fall into one of four defensible categories: (a) section/marker names requiring verbatim quote for cross-reference fidelity; (b) defined terms / canonical vocabulary whose verbatim use is required by AC3 or contract-shape; (c) contract-template language from ROLES.md §3 that any instantiation will echo when answering its Field N declaration; (d) citation-string descriptions naming what a referenced doc covers.

The four substantive paragraphs that *did* duplicate kernel reasoning in α's first-pass draft were rephrased: §0 Purpose's correction-surface paragraph, §Persona/Protocol/Project's Sigma discipline-profile description, Field 4's outward-membrane δ-holds-gate-authority paragraph, and Field 6's α=β-collapse "order-0 observation masquerading as order-1" argument. Each was replaced with an explicit citation to the source-of-truth section and a CDS-specific qualifier only (engineering-substrate, project-binding, or contract-shape addition).

**Architecture-choice section spot-check:** CDS.md §"Architecture choice" cites `schemas/cdd/README.md §"Architectural choice"` and `cnos#388` rather than restating the five-reason rationale verbatim. The CDS-side rationales (1)–(5) are CDS-perspective rephrasings, not copies of CDR.md's research-side prose. β verified by reading both versions side-by-side: rationale (1) on the CDS side names "engineering's stalled-loop failure is α-internal; research's overclaim failure is α-blind" — a CDS-perspective rephrasing of the same point CDR's rationale (1) makes from the research side. β confirms this is the per-protocol authorship discipline #407 requires.

**AC4 verdict:** PASS.

### AC5 — SKILL.md updated to load CDS.md

**β check:** structural-diff cds/SKILL.md vs the cycle/406-shipped version + read of the four removed/edited passages.

```
$ grep "^1\. Load \[\`CDS\.md\`\]" src/packages/cnos.cds/skills/cds/SKILL.md
1. Load [`CDS.md`](CDS.md) in this directory as the canonical instantiation contract. CDS.md declares the six fields per [`ROLES.md §3`](...) ...
```

Step 1 of Load order loads CDS.md as the canonical instantiation contract. The Step 1 description names what CDS.md declares (six fields + architectural-choice inheritance + persona/protocol/project boundary + empirical-anchor citation), parallel to cdr/SKILL.md's Step 1 description.

```
$ grep -c "Pending Sub 2\|will be the only" src/packages/cnos.cds/skills/cds/SKILL.md
0
```

Zero remaining "Pending Sub 2 / will be the only" notes about CDS.md.

Edited passages verified:
- **Frontmatter:** `outputs:` line changed from "canonical CDS contract loaded (when CDS.md lands at Sub 2)" → "canonical CDS contract loaded (CDS.md)". `requires:` removes "extraction map exists (this Sub 1) until CDS.md lands (Sub 2)".
- **Load order paragraph:** "**v0.1 status.** `CDS.md` is **not yet present**..." → "`CDS.md` is shipped in this directory by Sub 2..."; reorientation describes role overlays (still forthcoming) rather than CDS.md (now shipped).
- **Load order Step 1:** "**Pending Sub 2...**" appended note → removed; Step 1 now describes CDS.md as the shipped canonical contract with summary of what it declares.
- **Rule section:** "[`CDS.md`](CDS.md), once shipped by Sub 2, will be the only normative source for:" → "[`CDS.md`](CDS.md) is the only normative source for:"; also expanded the empirical-anchor item to name `usurobor/cnos` and the actor-collapse item to name both the α=β prohibition and the β-α-collapse-on-δ permission.
- **Conflict rule:** "If this file and [`CDS.md`](CDS.md) disagree (once Sub 2 lands `CDS.md`)..." → "If this file and [`CDS.md`](CDS.md) disagree..."; the parenthetical removed since CDS.md has shipped.
- **v0.1 caveat:** entire paragraph rewritten: "Until Sub 2 of cnos#403 lands CDS.md..." → "`CDS.md` is shipped (Sub 2 of cnos#403; cycle/407). The role-overlay files... remain **forthcoming** until Subs 3–5 land..." The caveat now describes the accurate post-Sub-2 / pre-Subs-3-5 state.

Role-overlay references in `calls:` frontmatter remain (5 overlays + CDS.md = 6 calls entries; CDS.md exists; 5 role overlays still advisory per cdr precedent).

**AC5 verdict:** PASS.

### AC6 — Empirical anchor cites cnos cycles

**β-rigor check:** read §Empirical anchor; verify named cycle subset includes representative milestones and that `docs/empirical-anchor-cdd.md` is named as the Sub 7 deliverable.

```
$ awk '/^## Empirical anchor/,/^## Related documents/' src/packages/cnos.cds/skills/cds/CDS.md \
    | grep -cE 'cnos#3[0-9]+|empirical-anchor-cdd'
12
```

12 cnos-cycle / empirical-anchor-cdd references in §Empirical anchor. β read of the three sub-sections:

- **§"Shape-compatibility claim"** — per-Field spot-check naming concrete cnos surfaces (`go run ./src/go/cmd/cn build --check`, `cue vet`, `schemas/cds/receipt.cue`, `.cdd/iterations/INDEX.md` aggregator with 37 rows, breadth-2026-05-12 wave manifest, etc.). The check binds the six fields to observed cnos practice; β verified each Field-N spot-check is non-vacuous (names actual cnos artifacts, not generic descriptions).
- **§"Representative cycle milestones"** — 7 cycles named (#364, #366, #376, #388, #393, #402, #403, plus the in-flight #406 + #407 as bracket points). β verified: each named cycle has a one-line context naming what it shipped + how it anchors CDS. The subset spans the full empirical wave (#364 doctrine landing through #406 cds-skeleton landing) without padding.
- **§"Sub 7 deferred surface-by-surface mapping"** — names `docs/empirical-anchor-cdd.md` as the forthcoming Sub 7 deliverable; cites `cnos.cdr/docs/empirical-anchor-cph.md` as the structural precedent; explains the artifact-name choice (cdd vs cnos) with reference to the historical naming + the cycle/406 README precedent.

**AC6 verdict:** PASS.

### AC7 — cnos.cdd untouched

**β-rigor check (hard rule):**

```
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/ | wc -l
0
$ git diff --stat origin/main..HEAD -- src/packages/cnos.cdd/
(empty)
```

Zero-byte diff. The entire `cnos.cdd` package directory is untouched. All 14 "pending cds extraction" markers in CDD.md remain in place (cleanup is Sub 6 territory per the implementation contract).

**AC7 verdict:** PASS (hard rule).

### AC8 — cnos.cdr untouched

**β-rigor check (hard rule):**

```
$ git diff origin/main..HEAD -- src/packages/cnos.cdr/ | wc -l
0
$ git diff --stat origin/main..HEAD -- src/packages/cnos.cdr/
(empty)
```

Zero-byte diff. The entire `cnos.cdr` package directory is untouched. CDS.md cites CDR.md as the structural sibling (§Related documents) and as a parallel record (§Architecture choice), but does not edit it.

**AC8 verdict:** PASS (hard rule).

### AC9 — No role overlays authored

**β check:**

```
$ find src/packages/cnos.cds/skills/cds/ -mindepth 1 -type d
(empty)

$ find src/packages/cnos.cds/skills/cds/ -type f | sort
src/packages/cnos.cds/skills/cds/CDS.md
src/packages/cnos.cds/skills/cds/SKILL.md
```

No subdirectories under `skills/cds/`; no `alpha/`, `beta/`, `gamma/`, `delta/`, `epsilon/`, or `operator/` directories. Only `CDS.md` and `SKILL.md` files. The `.gitkeep` was removed (per AC9 explicit permission since CDS.md now occupies the directory). The pre-existing `cdr` precedent (`cnos.cdr/skills/cdr/`) similarly has CDR.md + SKILL.md + role overlays at parity points; CDS at v0.1 (post-Sub-2) has the contract and loader, with role overlays slated for Subs 3–5.

**AC9 verdict:** PASS.

## β-rigor: doctrine-coherence-with-cdr-template structural diff

**Check:** does CDS.md mirror CDR.md's section structure such that a reader comparing the two side-by-side sees parallel structure?

| Section | CDR.md | CDS.md | β note |
|---|---|---|---|
| HTML manifest comments | ✓ (sections / completed) | ✓ (sections / completed) | Identical 7-section manifest. |
| Top heading + blockquote pointer to ROLES.md | ✓ | ✓ | Both name the c-d-X letter + the matter-type domain. |
| Version / Status / Date / Placement / Audience / Scope block | ✓ | ✓ | Same 6-field block; CDS.md cites Sub 2 of #403 where CDR.md cites Sub 1 of #376. |
| `## 0. Purpose` | ✓ | ✓ | Both name loss function via `ROLES.md §4a.2`; both contrast with sibling discipline; both name the role-grammar-shared / discipline-diverges principle. |
| `## Architecture choice` (with §The decision, §Option (b), §Design source, §Cross-reference) | ✓ | ✓ | Same four sub-section shape; CDS-side rationales (1)–(5) are CDS-perspective rephrasings, not copies of CDR's prose. |
| `## Persona, Protocol, Project` (with Layer 1 / Layer 2 / Layer 3 / Why this matters) | ✓ | ✓ | Same four sub-section shape; CDS names Sigma + cnos.cds + cnos itself as the three layer exemplars, where CDR names Rho + cnos.cdr + cph. |
| `## Six-field instantiation contract` (with one sub-section per field) | ✓ | ✓ | Six `### Field N: <name>` headings; per-field structure mirrors CDR's (each field declares its taxonomy, names what it determines, includes a sibling-protocol contrast, and includes a Sub-N-vs-Field-M scope-discipline line). |
| `## Empirical anchor` (with shape-compatibility claim + paths + deferred mapping) | ✓ | ✓ | CDR cites cph (external); CDS cites cnos itself; both include shape-compatibility claim + representative-paths/cycles + Sub-N-deferred mapping pointer. |
| `## Related documents` | ✓ | ✓ | Both group by inherits/cites/extends; CDS cites CDR as the sibling structural exemplar. |
| `## Non-goals` | ✓ | ✓ | Both enumerate what the file does not do; CDS-specific non-goals list 12 items (CDR's lists 8), reflecting the wider Subs 3–7 territory CDS defers to. |

CDS.md mirrors CDR.md's section structure exactly. The per-field structure (taxonomy / sibling-contrast / Sub-N-line) is a CDS-specific refinement: CDR's per-field sections do not always include the explicit Sub-N-vs-Field-M scope-discipline line (CDR has only 4 sub-deliverables total in #376; the Sub-N boundaries are less load-bearing). CDS has 7 sub-deliverables (Subs 1–7 in #403), so the Sub-N-vs-Field-M lines are load-bearing for #407's scope-discipline.

**β verdict:** doctrine-coherence-with-cdr-template PASS.

## β-rigor: name-overpromise audit on Sub 3/4/5 references

**Check:** does CDS.md frame Sub 3/4/5 references as "Sub N owns operational detail; CDS.md owns the contract" — not "Sub N will say X" (which would over-commit Sub N's authoring)?

β scan of all Sub-N references in CDS.md (`grep "Sub [3-5]"`):

- **Field 1 Sub-3-vs-Field-1 line:** "Field 1 owns the matter-type taxonomy (what α produces); Sub 3 owns the per-step lifecycle procedure (how α produces it across the 0–13 cycle steps)." — **acceptable** (declares Sub 3 owns operational, doesn't pre-author Sub 3's content).
- **Field 2 Sub-5-vs-Field-2 line:** "Field 2 owns the oracle taxonomy; Sub 5 owns the per-oracle operational detail." — **acceptable**.
- **Field 3 Sub-4-vs-Field-3 line:** "Field 3 owns the artifact-set taxonomy + the typed-receipt citation; Sub 4 will detail the artifact contract..." — **acceptable** ("will detail" names scope, not content).
- **Field 4 Sub-3-vs-Field-4 line:** "Field 4 owns the cadence taxonomy; Sub 3 owns the per-step lifecycle operational detail." — **acceptable**.
- **Field 5 Sub-5-vs-Field-5 line:** "Field 5 owns the class taxonomy; Sub 5 owns the per-class operational detail." — **acceptable**.
- **Field 6 Sub-3-vs-Field-6 line:** "Field 6 owns the actor-shape taxonomy; Sub 3 owns the per-shape operational detail." — **acceptable**.

No "Sub N will say X" framings found. All Sub-N references are scope-discipline statements, not content commitments. The extraction-map cited at every relevant point is the right pointer for "what content lands at the Sub-N surfaces"; CDS.md does not duplicate the extraction-map's per-row destinations.

**β verdict:** name-overpromise audit PASS.

## β observations (non-blocking, not findings)

**Obs-1:** CDS.md is 1040 lines vs the issue body's 500–700 target. β assessment: the issue body explicitly says "Lines budget is a target, not a hard ceiling." The excess length comes from substantive scope-discipline content (per-Field Sub-N-vs-Field-M lines, per-Field empirical-anchor spot-checks, per-rationale CDS-side reframings). β verified each excess passage is non-padding. Acceptable; not a finding.

**Obs-2:** CDS.md Field 5 declares `cds-*-gap` as the new canonical gap-class naming, replacing the empirical-anchor `cdd-*-gap` names. β assessment: the rename follows ROLES.md §4b.3's `{protocol}-{axis}-gap` convention mechanically; the four class names (skill/protocol/tooling/metric) are preserved across the rename. The transition is described as organic (future CDS cycles use new names; pre-#403 cycles keep historical names). β verified the rename is package-attribution only, not class-set change. Acceptable; not a finding.

**Obs-3:** CDS.md Field 6 formally declares the β-α-collapse-on-δ pattern as a permitted collapse for skill/docs-class cycles. β assessment: this is the first canonical declaration of the pattern (previously a wave-manifest precedent + dispatch-shape convention). The declaration includes stated conditions (matter is structural-mirror/migration/contract-authoring, not novel substantive code) and configuration-floor consequences (γ-axis and β-axis capped at A- per release/SKILL.md §3.8). β verified the conditions exclude substantive software work and verified the configuration-floor citation resolves. Acceptable; arguably the declaration *was overdue* (the empirical practice was operating without doctrine). Not a finding.

**Obs-4:** CDS.md uses post-rename paths (`<project>/.cds/`) for project binding citations while acknowledging current cnos practice uses `.cdd/`. β assessment: this is forward-looking contract authorship; the rename is named as the relevant open coordination question in `docs/extraction-map.md §14`. β verified the acknowledgment is explicit (every `.cds/` path mention is paired with the `.cdd/` historical-state note). Acceptable; not a finding.

**Obs-5:** CDS.md Field 4 introduces three-time-scale cadence (per-cycle, per-release, per-wave) where CDR.md Field 4 only declares per-wave. β assessment: this is a structural divergence from CDR driven by the loss-function difference — engineering has releases (CDR does not); engineering has per-merge boundary effections (CDR's research-wave close is the only boundary). The taxonomy distinguishes CDS from CDR cleanly. Acceptable; correctly cited as the engineering-side specialization.

**Obs-6:** CDS.md Field 3 declares both the narrative close-out artifact set (six files per CDD §5.5b) AND the typed `#CDSReceipt` schema; CDR Field 3 declared only the typed `#CDRReceipt`. β assessment: this is the engineering substrate's status — pre-Phase 3 of cnos#366, V is wired in but the CDS-cycle wave is still authoring close-outs as narrative files (cycle/335 onward). The two-surface declaration accurately describes the current state and the target state simultaneously. The "relationship between the two surfaces" paragraph explains the composition: typed receipt's `evidence_refs` point at narrative close-outs. Acceptable; describes the actual empirical practice.

**Obs-7:** CDS.md §"Architecture choice" rationales (1)–(5) include one CDS-side framing improvement over CDR.md's research-side framing: rationale (5) is renamed "Decision-once-applied-thrice" (CDR uses "Decision-once-applied-twice") to reflect that cnos#388 → cnos#376 → cnos#403 is the third application of the same architectural decision-class. β assessment: this is a numeric update reflecting empirical reality; not a doctrine change. Acceptable.

## Trigger assessment

| Trigger | Fire condition | Fired? | β note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVED. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 0 binding findings; 7 advisory observations. |
| Avoidable tooling / environment failure | environment blocked the cycle | **No** | `go run ./src/go/cmd/cn build --check` ran cleanly post-edit. |
| CI red on merge commit | CI fails post-merge | **N/A** | merge not yet executed (operator authority). |
| Loaded skill failed to prevent a finding | skill underspecified | **No** | no findings; the loaded skills (cdd, design, issue/contract, issue/proof) all guided the work successfully. |

No §9.1 triggers fired.

## Verdict

**R1 APPROVED.** AC1–AC9 PASS mechanically. AC4 no-duplication audit verified (the highest-rigor focus per #407 active design constraints): 11 remaining 50+char overlaps are all citation strings, marker names, defined vocabulary, or contract templates — no paragraph of duplicated kernel reasoning. Doctrine-coherence-with-cdr-template structural diff PASS. Name-overpromise audit on Sub 3/4/5 references PASS. No binding findings. Seven non-blocking observations recorded above. No fix-round needed.

The cycle is ready for close-out and operator-facing merge instruction.
