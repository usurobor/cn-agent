<!-- sections: [Gap, Design, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap, Design, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->

# Self-Coherence — Cycle #390

## §Gap

**Issue:** [#390](https://github.com/usurobor/cnos/issues/390) — Sub 1 of [#376](https://github.com/usurobor/cnos/issues/376): CDR six-field instantiation contract + architecture-choice declaration

**Version / mode:** unreleased · design-and-build · γ+α+β-collapsed on δ (single Claude Code session; per breadth-2026-05-12 wave manifest precedent, validated empirically across cycles 375/377/378/388; acceptable for docs-only contract-authoring class because β's gates are mechanical: `rg`/`wc`/`grep` against the new file).

**Gap:** `cnos.cdr` does not yet exist as a protocol-skill package. `ROLES.md §3` mandates the six-field instantiation contract every c-d-X protocol must declare. [cnos#388](https://github.com/usurobor/cnos/issues/388) (Phase 2.5) decided option (a) split (common kernel + per-protocol overlay) for schemas; the same architectural reading applies to skills ([cnos#376 AC7](https://github.com/usurobor/cnos/issues/376)), but no document records that inheritance for the CDR skill bundle. Without (i) the six-field contract declared in research-loss-function language and (ii) the architectural-choice section recording (a)-inheritance, downstream Sub 2 (package skeleton), Sub 3 (role overlays), and Sub 4 (empirical-anchor doc) cannot author against a stable shape.

**Goal:** Ship `src/packages/cnos.cdr/skills/cdr/CDR.md` declaring the six instantiation-contract fields in research-loss-function language, the architectural-choice inheritance, the persona/protocol/project boundary, and the cph empirical-anchor citation with shape-compatibility claim. Sub 2/3/4 inheritance ready.

## §Design

See [`design-notes.md`](./design-notes.md) for the full L7 design record (architectural choice statement, six-field rough prose, persona/protocol/project boundary structure, empirical-anchor citation strategy, alternatives, leverage, success criteria, open questions). Decision summary:

- **Option (a) — common constitution + per-protocol procedures** chosen for skills (inheriting cnos#388's schema decision). Rationale (1)–(5) cited by reference to `schemas/cdd/README.md §"Architectural choice"` rather than re-derived.
- **Option (b) — common protocol-agnostic skill + domain overlay** rejected with the same five reasons cnos#388 used for schemas, transposed to skills.
- **Persona/protocol/project boundary** stated as three doctrinal layers with canonical homes: `cn-rho/spec/` (persona), `cnos.cdr/skills/cdr/` (protocol overlay), `<project>/.cdr/` (project binding).
- **Six fields** declared in research-loss-function language: matter type names research artifacts (claims, hypotheses, methods, datasets, analyses, reports); review oracle names research-discipline checks (falsifiability, diagnostic oracles, reproduction-from-clean, citation integrity, data-policy compliance, claim/evidence alignment); γ close-out cites `schemas/cdr/receipt.cue` (`#CDRReceipt`) with GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO verdict vocabulary; δ cadence is gate-transition-shaped (not release-shaped); ε cadence is receipt-stream review over research-failure trigger classes; actor collapse rule pins α≠β as never-permitted for research claims.
- **Empirical anchor** cites `usurobor/cph` by path with shape-compatibility claim; detailed mapping deferred to Sub 4.

## §Skills

**Tier 1 (read-only doctrinal):**
- [`ROLES.md`](../../../ROLES.md) §3 (six-field contract), §4a (five-layer enforcement chain), §4a.2 (loss-function distinction), §4a.3 (receipt sketch), §7 (naming convention)
- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`](../../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md) §"Separate persona, protocol, and project" + §"Wave 4 — CDR bootstrap"
- [`schemas/cdd/README.md §"Architectural choice"`](../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split) — decision-record inherited
- [`schemas/cdr/receipt.cue`](../../../schemas/cdr/receipt.cue) — typed γ close-out surface
- [`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../../src/packages/cnos.cdd/skills/cdd/CDD.md) — structural exemplar (heading shape only)

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline; architectural-choice authoring shape
- `src/packages/cnos.cdd/skills/cdd/issue/contract/SKILL.md` — contract authoring (six-field shape is structurally a contract)
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring; β re-checks
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` §5.6b — `cdd-iteration.md` closure-gate

## §ACs

Per-AC oracles run against the branch `cycle/390` HEAD.

**AC1 — Six fields declared per ROLES.md §3:**
- Oracle: `rg "^### Field" src/packages/cnos.cdr/skills/cdr/CDR.md | wc -l` returns 6; `rg -c "TBD" src/packages/cnos.cdr/skills/cdr/CDR.md` returns 0; each `### Field N` subsection uses research-loss-function language; `rg "release|deploy|tag" CDR.md` returns matches only in cross-references or non-CDR-δ contexts.
- Evidence:
  - `rg "^### Field" CDR.md | wc -l` → `6` ✓
  - `rg -c "TBD" CDR.md` → `0` ✓
  - Six fields present: `### Field 1 — Matter type`, `### Field 2 — Review oracle`, `### Field 3 — γ close-out artifact`, `### Field 4 — δ cadence`, `### Field 5 — ε iteration cadence`, `### Field 6 — Actor collapse rule`.
  - `release|deploy|tag` matches at lines 372–375 (Field 4 explicitly disavowing release/tag/deploy: "gate-transition-shaped, not release-shaped... A CDR δ does not 'cut a release' or 'tag a deploy'"), line 388 ("not calendar or release-bundle"), line 443 (Field 5 mentions "tagged with one of the classes above" — `tagged` is a label-tag, not a release-tag, used in the sense of issue-class labelling), line 524 (Field 4 spot-check in §Empirical anchor reiterating "gate-transition-shaped, not release-shaped"). All occurrences are in **disavowing contexts** or in non-release-tag senses. None embed release/deploy/tag as the CDR δ cadence. ✓
  - Research-loss-function language verified per field:
    - Field 1: claims/hypotheses/methods/datasets/analyses/reports + `claim_status` calibration + "CDR does not produce software artifacts as matter."
    - Field 2: falsifiability, diagnostic oracles, reproduction-from-clean, citation integrity, data-policy compliance, claim/evidence alignment + "primary failure the oracle catches is overclaim" + explicit rejection of engineering's "compiles + tests pass."
    - Field 3: claim_refs/data_refs/method_refs/result_refs/claim_status/limitations/reproduction (per ROLES §4a.3) + cites `schemas/cdr/receipt.cue` + GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO verdicts.
    - Field 4: gate-transition-shaped (wave open → wave in progress → wave close) + receipt is the artifact + next-wave triggered by gate verdict.
    - Field 5: receipt-stream review over research-failure trigger classes (missing data gates, overclaiming, unreproducible numbers, weak citation discipline, recurring oracle ambiguity, construct drift).
    - Field 6: α≠β never permitted for research claims (with the rationale "α reviewing α's own claims is order-0 observation"); γ=δ + ε=δ permitted with conditions; engineering-class collapse precedent explicitly does **not** transfer.
- **Status: ✓ met**

**AC2 — Architecture-choice section declared:**
- Oracle: `## Architecture choice` section exists; option (a) named; option (b) named as rejected; cnos#388 cited; rationale (1)–(5) inherited from `schemas/cdd/README.md §"Architectural choice"`.
- Evidence:
  - `grep -n "^## Architecture choice" CDR.md` → line 55 ✓
  - Section spans lines 55–149.
  - Option (a) named at line 66: `**Option (a) — common constitution + per-protocol procedures.**`
  - Option (b) named as rejected at line 93: `### Option (b) — common protocol-agnostic skill + domain overlay (rejected)`
  - cnos#388 cited at lines 58, 88, 145.
  - `schemas/cdd/README.md §"Architectural choice"` cited at lines 56–57 and 121–124 (the design-source citation).
  - Rationale (1)–(5) enumerated at lines 95–117, each transposed to skills with explicit reference to cnos#388's rationale-5 ("decision-once-applied-twice").
  - Essay `§"Separate persona, protocol, and project"` cited at lines 125–128.
  - cnos#376 AC7 cross-referenced at lines 147–148.
- **Status: ✓ met**

**AC3 — Persona/protocol/project boundary section declared:**
- Oracle: `## Persona, Protocol, Project` section exists; three layers named with their canonical homes; statements about Rho≠CDR and CDR≠persona present.
- Evidence:
  - `grep -n "^## Persona, Protocol, Project" CDR.md` → line 154 ✓
  - Section spans lines 154–230.
  - Layer 1 (Persona) at line 164: canonical home `cn-rho/spec/PERSONA.md` + `cn-rho/spec/OPERATOR.md`. Contains: "Rho may play δ for CDR cycles" + "Rho is not CDR" + "CDR.md does not author persona content; CDR.md does not state 'I am Rho' or 'my voice'."
  - Layer 2 (Protocol overlay) at line 182: canonical home `cnos.cdr/skills/cdr/`. Contains: "CDR is not a persona" + "CDR defines α/β/γ/δ/ε research role overlays independent of any one persona hub."
  - Layer 3 (Project binding) at line 201: canonical home `<project>/.cdr/` or `<project>/cdr/`. Contains: cph as canonical example + "CDR.md does not author project-specific bindings."
  - Why-this-matters subsection at line 211 names the leakage to avoid (cph-specific gate names inside CDR.md, persona content inside CDR.md, etc.).
- **Status: ✓ met**

**AC4 — Empirical anchor validation:**
- Oracle: `## Empirical anchor` section names cph; shape-compatibility claim present; detailed mapping deferred to Sub 4.
- Evidence:
  - `grep -n "^## Empirical anchor" CDR.md` → line 495 ✓
  - Section spans lines 495–540.
  - cph cited as the empirical anchor (lines 497–502): "CDR's empirical anchor is `usurobor/cph` — the research repository where the CDR pattern is realised."
  - Shape-compatibility claim explicit (lines 506–530) with spot-checks across Field 1, Field 3, Field 4, Field 5.
  - cph's specific gate verdicts (`bounded GO`, `partial GO`, `INDETERMINATE`, `construct-survives-subject-to-caveats`) mapped onto CDR vocabulary (BOUNDED-GO, INDETERMINATE) without re-derivation.
  - cph canonical paths listed (lines 534–540) **by reference only**, no cph-local prose embedded — preserves the no-surface-fusion discipline.
  - Detailed mapping explicitly deferred to Sub 4 at lines 542–547.
- **Status: ✓ met**

**AC5 — No surface fusion:**
- Oracle: `rg "polling|dispatch|claude -p|gh issue" CDR.md` returns 0 hits in normative sections (or only in "out of scope" / "Sub 2/3 owns this" disclaimers); persona-identity prose absent.
- Evidence:
  - Token sweep: `rg "polling|dispatch|claude -p|gh issue" CDR.md` returns the following hits (all classified):
    1. Line 156–157: "operational mechanics (dispatch, polling, repo wiring) belong in role-overlay skills (Sub 3)" — **explicit non-goal/disclaimer**.
    2. Line 277: "dispatch boundary (per `ROLES.md §4a.2`) routes the work to the matching discipline" — **doctrinal citation** to ROLES.md's dispatch-boundary concept (loss-function dispatch), not runtime mechanics.
    3. Line 322: "`action: repair_dispatch` (β returns the receipt for revision)" — **schema-enum reference** to `schemas/cdd/boundary_decision.cue`'s `#BoundaryDecision.action` enum value.
    4. Line 369: "How δ selects a new research gap and dispatches a new wave." — **doctrinal verb** from `ROLES.md §1` ("δ operates... dispatches"). This is the role-ladder verb, not a runtime mechanism.
    5. Line 380: "disagreement) and dispatches α + β" — same doctrinal verb.
    6. Line 611: "Specify runtime mechanics (dispatch, polling, git wiring, CI hooks). Those belong in role-overlay skills (Sub 3)" — **explicit non-goal**.
    7. Line 614: "`protocol_id` dispatch per [`schemas/cdd/README.md §'protocol_id dispatch convention'`]" — **schema cross-reference** to V's dispatch convention.
  - No hits embed runtime mechanics in normative CDR doctrine. All hits are either disclaimers, doctrinal citations (ROLES verbs), or schema cross-references. ✓
  - Persona-identity check: `rg -i "I am Rho|my voice|as Rho" CDR.md` matches at lines 175–176 ("CDR.md does not state 'I am Rho' or 'my voice'") — this is the **explicit prohibition statement**, not persona content. ✓
  - Project-specific gates: cph's gate verdicts are referenced abstractly in §Empirical anchor (as project-naming variants of the doctrinal vocabulary); no cph-specific threshold values, dataset paths, or cph-local definitions appear in CDR.md normative sections.
- **Status: ✓ met**

**AC6 — Sub 2/3/4 inheritance ready:**
- Oracle: each Sub-N's expected scope per cnos#376 body can be authored against CDR.md without re-deriving fields/architecture.
- Evidence (read-through walk):
  - **Sub 2 (package skeleton)** — Sub 2 needs to author `cn.package.json`, package `README.md`, and root `SKILL.md`. CDR.md provides:
    - Package placement (line 18: `src/packages/cnos.cdr/skills/cdr/`) — Sub 2's `cn.package.json` cites this.
    - Package purpose (preamble + §0 Purpose) — Sub 2's `README.md` cites this.
    - Architectural inheritance (`## Architecture choice`) — Sub 2's package metadata can declare `imports: cnos.cdd` by reference.
    - Sub 2 does not need to re-derive any field or rationale.
  - **Sub 3 (role overlays)** — Sub 3 needs to author `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `delta/SKILL.md`, `epsilon/SKILL.md` for CDR. CDR.md provides:
    - α discipline (`### Field 1 Matter type`) — Sub 3's `alpha/SKILL.md` cites this for matter-class.
    - β discipline (`### Field 2 Review oracle`) — Sub 3's `beta/SKILL.md` cites this for review-oracle structure.
    - γ discipline (`### Field 3 γ close-out artifact`) — Sub 3's `gamma/SKILL.md` cites `#CDRReceipt` + GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO verdict vocabulary.
    - δ discipline (`### Field 4 δ cadence`) — Sub 3's `delta/SKILL.md` cites gate-transition cadence + wave open/in-progress/close transitions.
    - ε discipline (`### Field 5 ε iteration cadence`) — Sub 3's `epsilon/SKILL.md` cites the six trigger classes.
    - Actor collapse rule (`### Field 6`) — Sub 3 cites the never-permitted α≠β rule.
    - Sub 3 does not need to re-derive any field; each role-skill is a procedural elaboration of the doctrinal field declaration.
  - **Sub 4 (empirical-anchor doc)** — Sub 4 needs to map cph's full `.cdr/` artifact set to CDR.md fields. CDR.md provides:
    - cph citation by canonical path (§Empirical anchor, lines 534–540).
    - Shape-compatibility claim (lines 506–530) as the spec to verify surface by surface.
    - GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO ↔ cph gate-verdict mapping baseline (lines 515–522) — Sub 4 extends this to every cph wave verdict.
    - Sub 4 does not need to re-derive the six-field shape; the mapping is structurally bound to CDR.md's fields.
- **Status: ✓ met**

## §Self-check

**Did α's work push ambiguity onto β?** No. Every AC has concrete oracle evidence (grep matches with line numbers, exit codes, classification of every hit in AC5's token sweep). The β-α-collapse is acknowledged below; β review applies the same mechanical oracles.

**Is every claim backed by evidence in the diff?**
- AC1: `rg "^### Field" CDR.md | wc -l` → 6; `rg -c "TBD" CDR.md` → 0; release/deploy/tag hits classified as disavowing.
- AC2: `## Architecture choice` section in diff, line 55; option (a) at 66; option (b) rejected at 93; cnos#388 cited.
- AC3: `## Persona, Protocol, Project` section in diff, line 154; three layers with canonical homes named.
- AC4: `## Empirical anchor` section in diff, line 495; cph cited; shape-compatibility claim; deferral to Sub 4.
- AC5: forbidden-token sweep performed and classified; persona-identity absent (except in the explicit prohibition statement).
- AC6: read-through walk performed; Sub 2/3/4 expected scopes verified citable.

**Is the research-loss-function language defensible in each field?**
Yes. Field 2 names overclaim explicitly as the failure mode and rejects "compiles + tests pass" by name. Field 3 cites `schemas/cdr/receipt.cue` for the typed surface. Field 4 explicitly states "not release-shaped" and uses wave-transition vocabulary throughout. Field 5's trigger classes are research-failure classes (missing data gates, overclaiming, unreproducible numbers, weak citation discipline, recurring oracle ambiguity, construct drift) — no engineering-failure classes (test failure, CI flake, deploy drift). Field 6 pins α≠β with the research-specific rationale "α reviewing α's own claims for overclaim is order-0 observation" and explicitly states the engineering-class collapse precedent does **not** transfer.

**Is the architectural-choice section a defensible inheritance?**
Yes. The (a) rationale is cited by reference to `schemas/cdd/README.md §"Architectural choice"` (the decision-record), not re-derived. The five rationale points are transposed to skills with skill-specific examples (engineering α vs research α; per-protocol skill files vs per-protocol schema files; future c-d-X like cdw and cda). The decision-once-applied-twice property (cnos#388 rationale 5) is honoured.

**Is the persona/protocol/project boundary doctrinal (not operational)?**
Yes. The three layers are named with canonical homes; no operational mechanics (dispatch wiring, polling forms, git identities) are specified. The "Why this matters" subsection names the leakage to avoid abstractly; no project-specific gate names appear in normative text.

**Is the cph empirical-anchor citation prose-free?**
Yes. cph paths are listed by reference; cph-local prose is **not** embedded. Spot-checks reference cph's gate verdicts abstractly (project-naming variants of doctrinal vocabulary), never as the canonical definitions.

**Is the β-α-collapse defensible for this cycle class?**
Yes for docs-only contract-authoring. β's gates are mechanical:
- `rg "^### Field" | wc -l` for field count (AC1).
- `rg -c "TBD"` for completeness (AC1).
- `grep -n "^## "` for section presence (AC2, AC3, AC4).
- `rg "polling|dispatch|claude -p|gh issue"` for forbidden-token classification (AC5).
- `rg -i "I am Rho|my voice|as Rho"` for persona-identity absence (AC5).
- Read-through walk against Sub 2/3/4 expected scope (AC6).
None require subjective judgment that an independent β could provide above what α has already verified mechanically. The two non-mechanical β actions (re-reading the field-language for research-loss-function discipline; re-reading the persona/protocol/project boundary for doctrinal-not-operational shape) are performed in `beta-review.md`. The collapse is named there.

## §Debt

Known debt named in this cycle's close-out (per `cdd/CDD.md` §close-out conventions):

1. **`cnos.cdr` package metadata is not authored.** Sub 2 of #376 owns `cn.package.json`, package-level `README.md`, and root `SKILL.md`. CDR.md lives in an otherwise-empty package directory until Sub 2 lands. Sub 2 fills around it without modifying CDR.md.

2. **Role-overlay skills are not authored.** Sub 3 of #376 owns `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `delta/SKILL.md`, `epsilon/SKILL.md` for CDR. CDR.md declares the doctrinal shape; Sub 3 elaborates the procedural detail per role.

3. **Empirical-anchor mapping is not exhaustive.** Sub 4 of #376 owns the full surface-by-surface mapping of cph's `.cdr/` artifact set to CDR.md fields. Sub 1 establishes the citation and the shape-compatibility claim; Sub 4 verifies it artifact by artifact.

4. **CDR gate verdict vocabulary is not enum-pinned in `schemas/cdr/receipt.cue`.** The schema uses `boundary_decision.action` + `transmissibility` as the typed surface; CDR.md maps GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO onto these abstractly. A future cycle may pin a CDR-specific verdict enum if cph's or downstream research projects' gate verdicts diverge. Recorded as design-notes §OpenQuestions item 1.

5. **The actor-collapse rule (Field 6) sets a doctrinal floor; project-specific stricter floors are permitted but no policy template is provided.** A future cycle (or Sub 4) may author a project-binding policy template for projects that need a stricter floor (e.g. "for externally-published claims, β must be a distinct human reviewer"). Recorded as design-notes §OpenQuestions item 3.

## §CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 0 Observe | issue #390 | — | Sub 1 of #376; AC6 + AC7 of #376 owned here |
| 1 Select | gamma-scaffold.md | — | docs-only contract authoring; design-and-build mode; γ+α+β-collapsed on δ |
| 4 Gap | design-notes.md | cdd/design, cdd/issue/contract | `cnos.cdr` has no six-field contract; architectural-choice inheritance from cnos#388 not recorded for skills |
| 5 Mode | design-notes.md + this self-coherence | cdd/design, cdd/issue/contract, cdd/issue/proof, cdd/post-release §5.6b | design-and-build; γ+α+β-collapsed on δ; L7 (changes substrate boundary so Sub 2/3/4 work becomes simpler); single new file deliverable |
| 6 Artifacts | gamma-scaffold.md, design-notes.md, `src/packages/cnos.cdr/skills/cdr/CDR.md`, this self-coherence | cdd/design, cdd/issue/contract | scaffold pushed; design notes record L7 rationale; CDR.md authored with §Architecture choice, §Persona/Protocol/Project, six `### Field N` subsections, §Empirical anchor |
| 7 Self-check | this artifact | cdd/issue/proof | AC1–AC6 all met; β-α-collapse acknowledged; debt named |

## §Review-readiness

**Implementation SHA (α review-readiness):** [will be filled by the alpha-closeout commit SHA on `cycle/390`]

**Verification (mechanical, run against `cycle/390` HEAD):**
- `rg "^### Field" src/packages/cnos.cdr/skills/cdr/CDR.md | wc -l` → `6` ✓ (AC1)
- `rg -c "TBD" src/packages/cnos.cdr/skills/cdr/CDR.md` → `0` ✓ (AC1)
- `grep -n "^## Architecture choice" src/packages/cnos.cdr/skills/cdr/CDR.md` → `55:## Architecture choice` ✓ (AC2)
- `grep -n "^## Persona, Protocol, Project" src/packages/cnos.cdr/skills/cdr/CDR.md` → `154:## Persona, Protocol, Project` ✓ (AC3)
- `grep -n "^## Empirical anchor" src/packages/cnos.cdr/skills/cdr/CDR.md` → `495:## Empirical anchor` ✓ (AC4)
- `grep -ni "cph\|usurobor/cph" src/packages/cnos.cdr/skills/cdr/CDR.md | head -5` → multiple matches ✓ (AC4)
- `rg "polling|dispatch|claude -p|gh issue" src/packages/cnos.cdr/skills/cdr/CDR.md` → 7 hits, all classified as disclaimers / doctrinal citations / schema cross-references (no runtime mechanics embedded) ✓ (AC5)
- Read-through walk for Sub 2/3/4 inheritance: verified each Sub-N's expected scope citable without re-derivation ✓ (AC6)

α signals review-readiness. β review proceeds against the same mechanical oracles.
