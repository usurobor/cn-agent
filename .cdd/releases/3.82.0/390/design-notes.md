<!-- sections: [Context, ArchitecturalChoice, SixFieldDraft, PersonaProtocolProject, EmpiricalAnchorStrategy, Alternatives, Leverage, SuccessCriteria, OpenQuestions] -->
<!-- completed: [Context, ArchitecturalChoice, SixFieldDraft, PersonaProtocolProject, EmpiricalAnchorStrategy, Alternatives, Leverage, SuccessCriteria, OpenQuestions] -->

# Design Notes — Cycle #390 (Sub 1 of #376)

**Cycle:** [#390](https://github.com/usurobor/cnos/issues/390) — CDR six-field instantiation contract + architecture-choice declaration
**Date:** 2026-05-21
**Mode:** L7 design half preceding the build half (single docs-only deliverable: `src/packages/cnos.cdr/skills/cdr/CDR.md`).

---

## §Context

This is Sub 1 of the cnos#376 wave (CDR v0.1 master). The wave shape is:

- **Sub 1 (this cycle):** CDR.md instantiation contract + architectural-choice declaration. Owns cnos#376 AC6 (persona/protocol/project boundary) + AC7 (architectural choice).
- **Sub 2 (downstream):** `cnos.cdr` package skeleton (`cn.package.json`, package `README.md`, root `SKILL.md`).
- **Sub 3 (downstream):** role overlays (`alpha/SKILL.md` … `epsilon/SKILL.md` for CDR).
- **Sub 4 (downstream):** empirical-anchor doc mapping cph's `.cdr/` artifact set to CDR.md field structure.

cnos#388 (Phase 2.5 of #366) shipped `schemas/cdr/receipt.cue` and recorded the architectural-choice decision for **schemas** at `schemas/cdd/README.md §"Architectural choice"`. The decision: **option (a) — common kernel + per-protocol overlay**. The same decision-class applies to **skills**, but the application was deferred to cnos#376 Sub 1 (this cycle).

This cycle records the (a)-inheritance for skills and authors the six-field contract.

---

## §ArchitecturalChoice — statement

**Option (a) — common constitution + per-protocol procedures.** Inherited from cnos#388's schema decision. Applied to skills means:

```text
common constitution                    per-protocol procedures
───────────────────────                ────────────────────────
cnos.cdd / skills / cdd /              cnos.cdr / skills / cdr /
  CDD.md (instantiation contract)        CDR.md (instantiation contract)
  COHERENCE-CELL.md (doctrine)            (uses cnos.cdd doctrine by reference)
  COHERENCE-CELL-NORMAL-FORM.md           (uses CCNF by reference)
  alpha/SKILL.md (engineering α)         alpha/SKILL.md (research α — Sub 3)
  beta/SKILL.md  (engineering β)         beta/SKILL.md  (research β — Sub 3)
  gamma/SKILL.md (engineering γ)         gamma/SKILL.md (research γ — Sub 3)
  delta/SKILL.md (engineering δ)         delta/SKILL.md (research δ — Sub 3)
  epsilon/SKILL.md (engineering ε)       epsilon/SKILL.md (research ε — Sub 3)
```

The **common constitution** is `cnos.cdd`: the coherence-cell normal form, the recursion equation, the role-ladder grammar, the receipt-validation discipline, the generic schemas. These do **not** name software evidence by name; they type the shape.

The **per-protocol procedures** are `cnos.cdr` (and future `cnos.cds`, `cnos.cdw`, `cnos.cda`): the matter-type-specific α/β/γ/δ/ε procedures, the review oracles, the close-out artifacts. Each protocol package imports `cnos.cdd` doctrine by reference and adds its own discipline.

**Option (b) — common protocol-agnostic skill + domain overlay (rejected).** The alternative would put α/β/γ/δ/ε as protocol-agnostic skills in `cnos.cdd` with thin per-domain overlays (CDS/CDR/CDW). This is rejected for the same five reasons cnos#388 used for schemas:

1. **The role grammar is the constitution; the loss function is the procedure.** Engineering α and research α share α's verb ("produces") but diverge sharply on what counts as production discipline. Engineering ships under repairable feedback; research holds under irreparable claim transmission. The discipline profile is per-protocol, not common.

2. **Clarity at the protocol boundary.** A research α reading `cnos.cdr/skills/cdr/alpha/SKILL.md` should see only the research-α discipline — no indirection through a generic skill plus a domain overlay.

3. **Mechanical generic-vs-domain boundary.** Different files plus different package directories make the boundary structural, not semantic. Reuse happens by reading + reference, not by inheritance gymnastics.

4. **Future c-d-X generalizes mechanically.** Adding cdw (writing) means adding `cnos.cdw/skills/cdw/CDW.md` + its overlays. The common constitution does not need to know about it.

5. **Decision-once-applied-twice.** cnos#388 decided (a) for schemas; the same reading applies here without re-derivation. This is the property cnos#388 named in its rationale point 5: "the decision is recorded once and applied twice."

The architectural-choice section in CDR.md will state (a) explicitly, name (b) as rejected, cite cnos#388 + `schemas/cdd/README.md §"Architectural choice"`, and reference `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol, and project"`.

---

## §SixFieldDraft — rough prose before formal authoring

Each field must use **research-loss-function** language per `ROLES.md §4a.2`: research optimises for *truth-preserving claim transmission under uncertainty*; primary failure mode is *overclaim*. The fields below are drafted as research-discipline declarations, not software-discipline declarations.

### Field 1 — Matter type (draft)

α produces: **research claims, hypotheses, methods, datasets, analyses, and reports.**

- **Claim**: a stated proposition about the world or about a system under study. Carries a calibration: observed, computed, inferred, hypothesized, or indeterminate.
- **Hypothesis**: a proposition not yet calibrated — a candidate claim awaiting measurement.
- **Method**: a procedure that produces evidence — script, protocol, experimental design, analysis pipeline. Method refs include script path + commit SHA.
- **Dataset**: a body of observations or measurements with a manifest and provenance (mount, checksum, source).
- **Analysis**: the application of method to dataset, producing result artifacts (output files).
- **Report**: the synthesis surface — field notes, wave reports, project briefs — assembling claims, methods, data, and results into a transmissible narrative.

The matter type determines what β reviews. CDR does **not** produce software artifacts as matter; software produced in support of research is software-evidence (a CDS-shaped artifact whose CDR-side use is by reference, not authorship).

**Not in scope:** engineering deliverables (those are CDS matter); editorial prose unrelated to research claims (that is CDW matter).

### Field 2 — Review oracle (draft)

How β determines that α's matter closes the declared gap. Each oracle is mechanically checkable or procedurally enforceable:

- **Falsifiability**: claims are stated so that some observation would falsify them. β rejects a claim too vague to falsify.
- **Diagnostic oracles**: where claims rest on measurements, the measurement procedure has a known failure-detection step (calibration check, signal-to-noise floor, control comparison). β verifies the diagnostic ran and passed.
- **Reproduction-from-clean**: β re-runs the producing command from a clean environment; output matches what the receipt claims. Required when `claim_status ∈ {observed, computed}`.
- **Citation integrity**: every claim derived from external work cites a resolvable reference; every cited reference supports the claim it is invoked for.
- **Data-policy compliance**: data refs comply with the project's data-use policy (consent, anonymisation, retention). β rejects a receipt whose data refs violate policy.
- **Claim/evidence alignment**: the claim's strength is no stronger than the evidence supports. An `observed`-status claim has data refs; a `computed`-status claim has method refs + script SHA; an `inferred`-status claim names the inferential step; a `hypothesized`-status claim is labelled as such.

The primary failure the oracle catches is **overclaim**: a claim made stronger than its evidence. Engineering's "it compiles + tests pass" oracle does not apply here.

### Field 3 — γ close-out artifact (draft)

What γ writes when a research cycle (a "wave") closes. This is the parent-facing **research receipt** typed by `schemas/cdr/receipt.cue` (`#CDRReceipt`).

The receipt carries:
- `claim_refs`: which claims this receipt asserts (≥1).
- `data_refs`: dataset / mount / manifest / checksum (≥1).
- `method_refs`: script paths + commit SHA (≥1).
- `result_refs`: output file paths (≥1).
- `claim_status`: enum (observed | computed | inferred | hypothesized | indeterminate).
- `limitations` (optional): explicit caveats.
- `reproduction` (optional but required for observed/computed): β re-ran, command + output match.
- plus the generic-kernel fields inherited from `schemas/cdd/#Receipt` (`protocol_id`, `boundary_decision`, `protocol_gap_count`, `protocol_gap_refs`, `verdict`).

**Gate verdict vocabulary** (carried in `boundary_decision.action` + `transmissibility`):
- **GO** — claim transmissible at full strength.
- **REVISE** — claim transmissible after named revisions.
- **NO-GO** — claim not transmissible; further work required.
- **INDETERMINATE** — evidence insufficient to determine transmissibility; carries `claim_status: indeterminate`.
- **BOUNDED-GO** — claim transmissible within named bounds (subject-to-caveats; partial GO with explicit scope).

The schema-side mapping: `boundary_decision.action ∈ {accept, release, reject, repair_dispatch, override}` + `transmissibility ∈ {accepted, degraded, not_transmissible}`. The CDR vocabulary above maps onto these structural primitives; cycle-side display vocabulary may differ from schema enum, but the typed surface is what V dereferences.

### Field 4 — δ cadence (draft)

How often δ selects a new research gap and dispatches a new cycle.

**Cadence shape: gate-transition-shaped (not release-shaped).** A CDR δ does **not** "cut a release" or "tag a deploy" — a research wave concludes when a research-gate verdict is recorded, not when an artifact ships to production. The wave-transitions are:

- **Wave open**: δ selects a research gap (an unanswered question, an unverified claim, a needed dataset) and dispatches.
- **Wave in progress**: α (the research producer) writes methods, runs them, produces data and analyses; β reviews against the review oracle (Field 2).
- **Wave close**: γ writes the research receipt (Field 3) carrying the gate verdict (GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO). δ records the boundary decision; the wave is closed.

The trigger for δ to open the next wave is **gate-transition**, not calendar or release-bundle: a NO-GO verdict may trigger a follow-up wave to close the gap; a GO verdict may trigger a downstream synthesis wave; a BOUNDED-GO may trigger a scope-expansion wave.

**Not:** "δ cuts a tag every two weeks." Research does not have a tag/deploy artifact in the engineering sense; the receipt is the artifact.

### Field 5 — ε iteration cadence (draft)

When ε assesses whether the protocol itself is coherent and what triggers ε work.

**Cadence: receipt-stream review over protocol gaps.** ε reads across receipts (many `#CDRReceipt` instances over many waves) and surfaces patterns that require protocol patches. The trigger classes are research-failure classes (not engineering-failure classes):

- **Missing data gates**: receipts shipping observed/computed claims without `data_refs` (or with under-specified manifests). Pattern signals the data-policy gate is not enforced.
- **Overclaiming**: receipts with `claim_status: observed` whose evidence supports only `inferred`. Pattern signals β's claim/evidence-alignment oracle is too lenient or under-specified.
- **Unreproducible numbers**: receipts with `reproduction.output_match: false` (or missing reproduction records when they were required). Pattern signals reproduction-from-clean is not running.
- **Weak citation discipline**: receipts with claims that cite no external source where one should exist. Pattern signals citation-integrity oracle drift.
- **Recurring oracle ambiguity**: the same review-oracle interpretation question reappears across cycles (β can't decide; γ has to adjudicate). Pattern signals oracle under-specification.
- **Construct drift**: a key research construct's definition changes across receipts without explicit deprecation. Pattern signals the project binding (cph or downstream) needs a glossary/construct-stability artifact.

ε's output is a CDR-shaped `cdd-iteration.md` (or CDR-specific equivalent, declared by Sub 3 — `epsilon/SKILL.md` owns the artifact spec). Trigger is **finding-triggered**, not calendar-triggered.

### Field 6 — Actor collapse rule (draft)

Which roles may be held by the same actor; which collapses are permitted, which are prohibited, what signal warrants splitting.

**Never permitted:**
- **α = β within a single cycle.** Always prohibited. Research β reviews α's claim/evidence alignment; α reviewing α's own claims for overclaim is order-0 observation masquerading as order-1. The structural independence is the mechanism by which review catches overclaim. **No research-class waiver exists** for this collapse — engineering-class precedent (collapsed-on-δ for mechanical refactor cycles) does **not** transfer to research-class work, because research's failure mode (overclaim) is precisely the mode α cannot self-detect.

**Permitted (with conditions):**
- **γ = δ for small project waves.** Allowed when the project is small enough that one actor reasonably holds both wave-coordination (γ) and gap-selection-across-waves (δ). The collapse is safe when γ's coordination authority does not compromise β's independence (it does not — γ sees both α and β but does not author either).
- **ε = δ until receipt-stream volume warrants split.** Allowed in small-protocol regimes where ε's receipt-stream review work does not justify a dedicated reviewer of the protocol. Split when finding volume crosses a threshold (TBD per project; cph and downstream may set their own).

**Signal for splitting collapsed roles:**
- Number of receipts per quarter passing some threshold (suggests ε load justifies its own actor).
- Recurrence of oracle-ambiguity findings (suggests γ/ε collapse is masking protocol-drift signal).
- A specific overclaim incident — δ should record the cycle, the receipt, and the structural cause; if α=β were collapsed, the cycle is **invalid** and must be re-run with separated actors.

---

## §PersonaProtocolProject — boundary structure

This section structure for CDR.md will name three doctrinal layers and their canonical homes. The boundary is **doctrinal**, not operational — it does not describe how to wire dispatch, polling, or repo plumbing. Those operational mechanics belong in role-overlay skills (Sub 3) and in persona/operator contracts (separate hubs).

### Layer 1 — Persona (canonical home: `cn-rho/spec/`)

A persona is what kind of mind is doing the work. For research, that mind is **Rho** (or any research persona — Rho is the named exemplar; cnos accepts other research personas if they declare the research discipline profile per `ROLES.md §4a.2`).

**Rho is not CDR.** Rho can play δ for CDR cycles; Rho can also play δ for CDW cycles or for non-protocol research-coordination work. The persona's discipline profile lives in `cn-rho/spec/PERSONA.md` + `cn-rho/spec/OPERATOR.md`. The persona hub is independent of any one protocol package.

**CDR is not a persona.** Anyone can author a CDR cycle if they enact the research-loss-function discipline; the role-cell shape is portable across personas. The CDR protocol overlay specifies what the role-cell looks like for research; it does not specify who plays the roles.

### Layer 2 — Protocol overlay (canonical home: `cnos.cdr/skills/cdr/`)

The protocol overlay specifies what counts as valid work in the research domain. It is the **role-cell discipline for research**, declared in CDR.md (this file) and authored as role-specific skills in Sub 3.

CDR's α/β/γ/δ/ε overlays are independent of any one persona hub. They specify:
- α's discipline (Field 1 matter type + research-α verbs).
- β's review oracle (Field 2).
- γ's close-out artifact (Field 3).
- δ's cadence (Field 4).
- ε's iteration cadence (Field 5).
- The actor collapse rule (Field 6).

CDR.md is the doctrinal declaration; Sub 3 authors the per-role procedural skills.

### Layer 3 — Project binding (canonical home: `<project>/.cdr/` or `<project>/cdr/`)

A project (e.g. cph) binds CDR's roles to concrete data, scripts, reports, and gates via its own `.cdr/` directory. The project binding answers:
- Which datasets does this project's α use?
- Which scripts produce evidence?
- Which gate verdicts apply for this project's research domain?
- Where do waves' receipts live?

cph (the empirical anchor; see §EmpiricalAnchorStrategy below) is the canonical example: its `.cdr/waves/` directory carries wave-shaped receipts; its `PROJECT.md` carries the project-level binding rules; its field reports + dataset docs are the per-wave evidence.

**CDR.md does not author project-specific bindings.** Project-specific gate names, dataset paths, and report templates live in the project's `.cdr/` directory, not in this doctrinal file.

### Why this matters (the leakage to avoid)

Conflating any two layers produces drift, exactly as `ROLES.md §4a` warns:
- Persona content (e.g. "I am Rho") inside CDR.md makes CDR un-reusable across personas.
- CDR.md gate names (e.g. cph's specific `bounded GO` definition with cph-local thresholds) makes CDR un-reusable across projects.
- Project content (cph's specific datasets) inside CDR.md makes downstream research projects re-derive what should be inherited doctrinally.

The boundary section in CDR.md states each layer's canonical home and explicitly names what CDR.md does **not** author.

---

## §EmpiricalAnchorStrategy

The empirical anchor for CDR is `usurobor/cph` — the research repo where the CDR pattern is realised in concrete waves, datasets, and field reports. CDR.md's `## Empirical anchor` section:

1. **Cites cph by path**, not by embedded prose. `usurobor/cph:.cdr/waves/`, `usurobor/cph:PROJECT.md`, `usurobor/cph:.cdr/field-reports/` — names the canonical paths.
2. **Declares the shape-compatibility claim**: cph's CDR practice is describable by the six-field shape declared in CDR.md without contradiction. Spot-checks named:
   - cph's wave-shape (observation → translation → analysis → measurement) fits Field 1's matter-type vocabulary (claims, hypotheses, methods, datasets, analyses, reports).
   - cph's gate verdicts (`bounded GO`, `partial GO`, `INDETERMINATE`, `construct-survives-subject-to-caveats`) map onto Field 3's gate vocabulary (GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO) — `partial GO` and `bounded GO` are project-naming variants of BOUNDED-GO; `construct-survives-subject-to-caveats` is a BOUNDED-GO with the bound named.
3. **Defers detailed mapping to Sub 4.** Sub 4 of #376 owns the full artifact-by-artifact mapping (cph's specific files → CDR.md fields → schemas/cdr/receipt.cue keys). Sub 1 establishes the citation and the compatibility claim; Sub 4 verifies it surface by surface.

**Refusal-condition discipline**: cph-local prose is not embedded in CDR.md. If a claim in the §Empirical anchor section requires copying cph-local content, the claim is restated abstractly + cph path is cited as the locus. The boundary preserves protocol/project separation.

---

## §Alternatives considered

### Alt 1 — Defer the six-field declaration to Sub 2

The alternative would be: Sub 1 ships only the architectural-choice section + persona/protocol/project boundary; Sub 2 ships the six-field contract alongside the package skeleton.

**Rejected.** The issue (#390) explicitly bundles the six-field contract with the architecture-choice declaration as Sub 1's deliverable, because Sub 2's package skeleton needs to cite the contract surface as the doctrinal anchor. Splitting them inverts the dependency order. cnos#376 AC6 + AC7 are owned by Sub 1.

### Alt 2 — Author the six fields as full procedural skills inline in CDR.md

The alternative would be: each `### Field N` subsection contains the full role-skill procedural detail (α's exact algorithm, β's exact gate-execution steps).

**Rejected.** Sub 3 owns role-overlay skills. CDR.md is the **doctrinal contract** — the field declarations name the shape, not the procedure. The procedural detail lives in `cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` (Sub 3). Embedding procedural detail in CDR.md would expand scope and pre-empt Sub 3's design space.

### Alt 3 — Inline the cph artifact-mapping in §Empirical anchor

The alternative would be: §Empirical anchor enumerates every cph file and maps it to a CDR.md field.

**Rejected.** Sub 4 owns this. Sub 1's compatibility claim is enough to validate the six-field shape; full mapping is Sub 4's deliverable. Inlining it here would expand scope and pre-empt Sub 4. Also: prose embedding from cph would violate the no-surface-fusion discipline (refusal condition 4 in the cycle directive).

### Alt 4 — Re-derive the architectural choice from first principles

The alternative would be: CDR.md re-runs cnos#388's rationale (1)–(5) in CDR-specific language.

**Rejected.** Decision-once-applied-twice (rationale 5 of cnos#388). CDR.md cites cnos#388 + `schemas/cdd/README.md §"Architectural choice"` by reference. Re-running the rationale is duplicative and risks drift if cnos#388's reasoning is later refined.

### Alt 5 — Use software-protocol verbs for δ cadence

The alternative would be: Field 4 uses "release" / "tag" / "deploy" verbs because they are familiar from CDS.

**Rejected.** This is exactly the surface fusion the issue (AC5) and the cycle directive (refusal condition 2) prohibit. Research's δ cadence is **gate-transition-shaped**: a wave concludes when a gate verdict is recorded. The receipt is the artifact; there is no tag/deploy. Field 4's draft uses "wave open / wave in progress / wave close" + "boundary decision" vocabulary.

---

## §Leverage

This is L7 work in the cdd/design vocabulary — it changes the substrate boundary so future work becomes simpler.

**L7 leverage realised here:**
- Sub 2 inherits a stable doctrinal target. The package skeleton can cite CDR.md without scope expansion.
- Sub 3 inherits the field structure as the spec for each role overlay. Each role's SKILL.md cites CDR.md's `### Field N` subsections.
- Sub 4 inherits the citation strategy. The full cph mapping is built against a fixed §Empirical anchor reference.
- Future c-d-X protocols (cdw, cda) follow the same instantiation-contract shape — CDR.md serves as the second example after CDD.md, reinforcing the pattern.
- The architectural decision-once-applied-twice rationale (cnos#388 rationale 5) is now applied twice. This validates the rationale as a reusable design principle for c-d-X work.

**Negative leverage:**
- Authoring CDR.md before Sub 2's package skeleton means the file lives in an otherwise-empty package directory. Mitigated: Sub 2 fills around it without modifying CDR.md.
- The six-field draft language must be defensible — if a future research-class practitioner finds the field declarations too software-shaped or too cph-specific, the file may need revision. Mitigated: design notes record the loss-function discipline explicitly; cph mapping is deferred to Sub 4 so cph specifics do not pollute the doctrinal surface.

---

## §SuccessCriteria

CDR.md ships when:

1. `rg "^### Field" CDR.md | wc -l` returns 6 (AC1).
2. `rg -c "TBD" CDR.md` returns 0 (AC1).
3. `## Architecture choice` section exists, names (a), rejects (b), cites cnos#388 (AC2).
4. `## Persona, Protocol, Project` section exists with all three layers + canonical homes (AC3).
5. cph cited as empirical anchor; shape-compatibility claim present; detailed mapping deferred (AC4).
6. `rg "polling|dispatch|claude -p|gh issue"` returns 0 hits in normative sections (AC5).
7. Read-through walk confirms Sub 2/3/4 can author against CDR.md without re-derivation (AC6).
8. β-α-collapsed review records APPROVE with mechanical AC re-check.

---

## §OpenQuestions

These are recorded but not blockers for Sub 1; surfaced for Sub 2/3/4 / future ε work:

1. **Should the gate verdict vocabulary (GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO) be enum-pinned in `schemas/cdr/receipt.cue`?** Currently the schema uses `boundary_decision.action` (accept/release/reject/repair_dispatch/override) + `transmissibility` (accepted/degraded/not_transmissible) as the typed surface. The CDR vocabulary in Field 3 is mapped onto these but not pinned. A future cycle may pin a CDR-specific verdict enum if cph's gate verdicts (or those of downstream research projects) drift. Recorded; not blocking Sub 1.

2. **Where does the cph PROJECT.md ↔ CDR.md cross-reference live structurally?** cph's `PROJECT.md` carries project-level binding. CDR.md cites cph by path. Should cph also carry a back-reference to `cnos.cdr/skills/cdr/CDR.md`? Probably yes, but that is a cross-repo concern owned by Sub 4 (the empirical-anchor doc) or by cph-side maintenance. Sub 1 declares the citation; Sub 4 handles the reciprocity.

3. **Does CDR's actor collapse rule (Field 6) need a stricter floor than ROLES.md §3 / §4 for cross-personnel research?** ROLES.md §4 names α=β as never-safe. CDR's research class arguably needs a stricter floor: for any externally-published claim, β must be a distinct human reviewer, not a same-human "second pass." This is a discipline-profile question that arguably belongs in `cn-rho/spec/PERSONA.md §Discipline`, not in CDR.md. Sub 1 records Field 6 at the doctrinal floor (α=β never; γ=δ and ε=δ permitted with conditions). Project-specific stricter floors live in project bindings (e.g. cph's `.cdr/POLICY.md`).

4. **Is `cnos.cdr` package versioned independently from `cnos.cdd`?** Sub 2 question. Sub 1 does not pin a version. The architectural-choice section says CDR procedures import CDD doctrine by reference; whether the import is version-pinned is Sub 2's design space.

5. **Should the six-field declarations cite the typed schema fields directly?** Field 3 cites `schemas/cdr/receipt.cue` (`#CDRReceipt`). Fields 1, 2, 4, 5, 6 cite ROLES + the essay. Should there be a deeper schema-side typing for, e.g., the actor collapse rule? Probably not at this layer — collapse rules are doctrinal, not receipt-typed. Recorded; not blocking.
