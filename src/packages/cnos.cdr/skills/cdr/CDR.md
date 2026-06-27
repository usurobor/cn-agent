<!-- sections: [preamble, architecture-choice, persona-protocol-project, six-field-contract, empirical-anchor, related-documents, non-goals] -->
<!-- completed: [preamble, architecture-choice, persona-protocol-project, six-field-contract, empirical-anchor, related-documents, non-goals] -->

# Coherence-Driven Research (CDR)

> CDR is a research-protocol instantiation of the generic role-scope ladder
> pattern. The pattern (α/β/γ/δ/ε roles, scope-escalation contract, six-field
> instantiation contract) is documented at
> [`ROLES.md`](../../../../../ROLES.md) at the repo root. CDR does not claim
> to have originated the role structure; it instantiates the structure for
> research work — investigation, synthesis, citation, dataset stewardship,
> claim transmission under uncertainty.

**Version:** 0.1 (Sub 1 of [cnos#376](https://github.com/usurobor/cnos/issues/376))
**Status:** Draft — instantiation contract only. Package skeleton (Sub 2),
role overlays (Sub 3), and empirical-anchor doc (Sub 4) are downstream.
**Date:** 2026-05-21
**Placement:** `src/packages/cnos.cdr/skills/cdr/`
**Audience:** Research operators, research reviewers, project maintainers
binding CDR to a concrete research repo
**Scope:** The doctrinal contract for what CDR is. Declares the six
instantiation-contract fields, records the architectural-choice inheritance
from [cnos#388](https://github.com/usurobor/cnos/issues/388), and names the
persona/protocol/project boundary.

---

## 0. Purpose

CDR is the research discipline used to transmit claims about the world (or
about a system under study) coherently. Its loss function — per
[`ROLES.md §4a.2`](../../../../../ROLES.md#4a2-loss-function-distinction) —
is *truth-preserving claim transmission under uncertainty*. The primary
failure mode CDR must structurally resist is **overclaim**: a claim becoming
stronger than its evidence; false knowledge propagating; the system
believing something it did not measure.

CDR shares the role-cell grammar (α/β/γ/δ/ε) with CDS (engineering), CDW
(writing), and CDA (analysis). What CDR does not share is the discipline
profile. Engineering optimises for *artifact improvement under repairable
feedback*; research optimises for *truth preservation under uncertainty*.
The role names are the same; the loss functions diverge; the receipts and
oracles diverge to enforce the divergent disciplines.

This document is the instantiation contract for CDR per `ROLES.md §3`. It
declares the six fields any c-d-X protocol must declare. It does not
specify how to operate any individual role — that is the work of the
per-role skills downstream (cnos#376 Sub 3). It does not specify how cph
(the empirical anchor) binds its specific artifacts to CDR's surfaces —
that is the work of the empirical-anchor doc downstream (cnos#376 Sub 4).
It does not author the `cnos.cdr` package metadata — that is cnos#376 Sub 2.

---

## Architecture choice

CDR inherits the architectural decision recorded at
[`schemas/cdd/README.md §"Architectural choice"`](../../../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split)
from [cnos#388](https://github.com/usurobor/cnos/issues/388) (Phase 2.5 of
cnos#366). cnos#388 decided the choice for **schemas**; cnos#376 AC7
extended the same decision-class to **skills/role-overlays**. This section
records the inheritance.

### The decision

**Option (a) — common constitution + per-protocol procedures.**

```text
common constitution                    per-protocol procedures
─────────────────────                  ────────────────────────
cnos.cdd / skills / cdd /              cnos.cdr / skills / cdr /
  CDD.md (doctrinal contract)            CDR.md (this file)
  COHERENCE-CELL.md                       (cites CCNF + CC by reference)
  COHERENCE-CELL-NORMAL-FORM.md
  RECEIPT-VALIDATION.md                    schemas/cdr/receipt.cue
  alpha/SKILL.md (engineering α)           alpha/SKILL.md (research α — Sub 3)
  beta/SKILL.md  (engineering β)           beta/SKILL.md  (research β — Sub 3)
  gamma/SKILL.md (engineering γ)           gamma/SKILL.md (research γ — Sub 3)
  delta/SKILL.md (engineering δ)           delta/SKILL.md (research δ — Sub 3)
  epsilon/SKILL.md (engineering ε)         epsilon/SKILL.md (research ε — Sub 3)
```

The **common constitution** lives in `cnos.cdd`: the coherence-cell normal
form, the recursion equation, the role-ladder grammar, the generic
receipt-validation interface, the generic schemas. The common constitution
does **not** name engineering evidence by name; it types the shape.

The **per-protocol procedures** live in `cnos.cdr` (and future `cnos.cds`,
`cnos.cdw`, `cnos.cda`): the matter-type-specific α/β/γ/δ/ε procedures, the
review oracles, the close-out artifacts. Each protocol package imports the
`cnos.cdd` doctrine by reference and adds its own discipline profile.

### Option (b) — common protocol-agnostic skill + domain overlay (rejected)

The rejected alternative would put α/β/γ/δ/ε as protocol-agnostic skills in
`cnos.cdd` with thin per-domain overlays (CDS/CDR/CDW). This is rejected
for the same five reasons cnos#388 used for schemas, transposed to skills:

1. **The role grammar is the constitution; the loss function is the
   procedure.** Engineering α and research α share α's verb ("produces") but
   diverge sharply on what counts as production discipline. Engineering
   ships under repairable feedback; research holds under irreparable claim
   transmission. The discipline profile cannot be common.

2. **Clarity at the protocol boundary.** A research α reading
   `cnos.cdr/skills/cdr/alpha/SKILL.md` (Sub 3 deliverable) sees only the
   research-α discipline — no indirection through a generic skill plus a
   domain overlay.

3. **Mechanical generic-vs-domain boundary.** Different files plus different
   package directories make the boundary structural, not semantic. Reuse
   happens by reading + reference, not by inheritance gymnastics.

4. **Future c-d-X generalizes mechanically.** Adding cdw (writing) means
   adding `cnos.cdw/skills/cdw/CDW.md` + its overlays. The common
   constitution does not need to know about it.

5. **Decision-once-applied-twice.** cnos#388 decided (a) for schemas; the
   same reading applies here without re-derivation. This is the property
   cnos#388 named in its rationale 5: "the decision is recorded once and
   applied twice."

### Design source

The decision rationale and the schema-level precedent are sourced from:

- [`schemas/cdd/README.md §"Architectural choice"`](../../../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split)
  — the cnos#388 decision-record. Rationale (1)–(5) for the schema case;
  inherited here by reference.
- [`docs/papers/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol,
  and project"`](../../../../../docs/papers/CCNF-AND-TYPED-TRUST.md)
  — the doctrinal source for treating persona, protocol, and project as
  three distinct layers.
- [`docs/papers/CCNF-AND-TYPED-TRUST.md §"Wave 4 — CDR bootstrap"`](../../../../../docs/papers/CCNF-AND-TYPED-TRUST.md)
  — declares the wave-shape: "Create `cnos.cdr` as the research overlay,
  anchored in cph but not defined by cph. CDR owns research role overlays.
  Rho owns the research persona. cph owns the project binding." CDR.md
  realises this declaration.
- [`ROLES.md §4a.2`](../../../../../ROLES.md#4a2-loss-function-distinction)
  — the loss-function distinction that makes engineering-discipline and
  research-discipline differ at the procedural layer while sharing the
  role-cell kernel.

### Cross-reference

- [cnos#388](https://github.com/usurobor/cnos/issues/388) — schema-side
  decision. The (a) decision recorded here for skills is its skill-side
  application.
- [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376) — names this
  cycle (Sub 1) as the place where the inheritance is recorded for skills.

---

## Persona, Protocol, Project

CDR is a **protocol overlay** — layer 3 of the five-layer enforcement chain
declared at
[`ROLES.md §4a`](../../../../../ROLES.md#4a-persona-operator-protocol-project-receipt--the-five-layer-enforcement-chain).
This section names CDR's relationship to the layers above and below it.
The boundary is doctrinal, not operational; operational mechanics (dispatch,
polling, repo wiring) belong in role-overlay skills (Sub 3) and in
persona/operator contracts (separate hubs).

### Layer 1 — Persona (canonical home: `cn-rho/spec/PERSONA.md` + `cn-rho/spec/OPERATOR.md`)

A persona is *what kind of mind is doing the work*. For research, the named
exemplar is **Rho** — the research persona whose discipline profile encodes
research's loss function (truth-preserving claim transmission under
uncertainty) per `ROLES.md §4a.4`. Other research personas are admissible
provided they declare a compatible discipline profile in their hub's
`PERSONA.md`.

**Rho may play δ for CDR cycles.** The persona enacts the role; the role
does not enact the persona. Rho may also play δ for non-CDR research-
coordination work, or play γ for a project that runs CDR + CDS in parallel.

**Rho is not CDR.** The persona hub `cn-rho` lives outside this package.
CDR.md does not author persona content; CDR.md does not state "I am Rho"
or "my voice." If a future research persona is named, the protocol overlay
remains unchanged — only the persona hub's `PERSONA.md` is added.

### Layer 2 — Protocol overlay (canonical home: `cnos.cdr/skills/cdr/`)

The protocol overlay specifies *what counts as valid work in the research
domain*. CDR's protocol overlay is this package: CDR.md (this file)
declares the doctrinal contract; the role-specific skills under
`cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` (Sub 3) will
declare each role's procedural discipline.

**CDR is not a persona.** Anyone can author a CDR cycle if they enact the
research-loss-function discipline; the role-cell shape is portable. CDR's
overlay specifies what the role-cell looks like for research; it does not
specify who plays the roles.

**CDR defines α/β/γ/δ/ε research role overlays independent of any one
persona hub.** The CDR α/β/γ/δ/ε skills (Sub 3) describe the research
discipline at each role; they cite Rho's discipline profile by reference
when describing what kind of persona is typical for CDR work, but they do
not embed persona-identity content.

### Layer 3 — Project binding (canonical home: `<project>/.cdr/` or `<project>/cdr/`)

A project binds CDR's roles to *concrete data, scripts, reports, and
gates*. The canonical example is cph (see §Empirical anchor below): cph's
`PROJECT.md` carries the project-level binding rules; cph's `.cdr/waves/`
directory carries the wave-shaped research receipts; cph's field reports
and dataset docs are the per-wave evidence.

**CDR.md does not author project-specific bindings.** Project-specific
gate names, dataset paths, threshold definitions, and report templates
live in the project's `.cdr/` directory. CDR.md declares the doctrinal
shape; projects bind concrete artifacts.

### Why this matters

Conflating any two of the three layers produces drift, exactly as
`ROLES.md §4a` warns. Specifically for CDR:

- **Persona content inside CDR.md** would make CDR un-reusable across
  research personas — a future research persona that is not Rho would have
  to fork the protocol.
- **Project-specific gate names inside CDR.md** would make CDR
  un-reusable across research projects — a downstream research repo that
  is not cph would have to fork the protocol.
- **Doctrinal content inside cph** (or any project) would trap the
  research discipline inside one project — Sub 4's mapping work would
  become an authoring task rather than a binding task.

The three canonical homes (`cn-rho/spec/`, `cnos.cdr/skills/cdr/`,
`<project>/.cdr/`) are the structural separation that prevents the drift.

---

## Six-field instantiation contract

Per `ROLES.md §3`, every c-d-X protocol must declare six fields. CDR
declares them below, each in research-loss-function language. The shape
follows CDD.md's structural precedent; the language is research, not
software.

### Field 1 — Matter type

α produces **research claims, hypotheses, methods, datasets, analyses, and
reports**.

- **Claim** — a stated proposition about the world or about a system under
  study. Every claim carries a calibration (`claim_status`): observed,
  computed, inferred, hypothesized, or indeterminate (per
  [`schemas/cdr/receipt.cue`](../../../../../schemas/cdr/receipt.cue)
  `#ClaimStatus`). A claim without a calibration is incomplete matter.
- **Hypothesis** — a proposition not yet calibrated. A candidate claim
  awaiting measurement or analysis. The hypothesis is α-matter when the
  hypothesis itself is the cycle's deliverable (e.g. a "what should we
  measure next" wave).
- **Method** — a procedure that produces evidence: a script, an
  experimental protocol, an analysis pipeline, a measurement procedure.
  Methods carry script paths + commit SHA per `ROLES.md §4a.3`.
- **Dataset** — a body of observations or measurements with a manifest and
  provenance: mount point, checksum, source attribution, data-use policy
  compliance.
- **Analysis** — the application of method to dataset, producing
  result-file artifacts. Analyses are α-matter when they generate new
  claims or refine existing ones.
- **Report** — the synthesis surface: field notes, wave reports, project
  briefs. Reports assemble claims, methods, data, and results into a
  transmissible narrative; reports cite the typed receipt
  (`#CDRReceipt`) that backs them.

The matter type determines what β reviews (Field 2), what γ writes to close
out (Field 3), and what ε observes across many cycles (Field 5).

**CDR does not produce software artifacts as matter.** Software produced
in service of research is **software-evidence** — a CDS-shaped artifact
whose CDR-side use is by reference (cited from `method_refs` with its
commit SHA), not authorship. A research wave that ships code as its primary
deliverable is structurally a CDS cycle, not a CDR cycle; the protocol
dispatch boundary (per `ROLES.md §4a.2`) routes the work to the matching
discipline.

### Field 2 — Review oracle

How β determines that α's matter closes the declared research gap. Each
oracle is mechanically checkable or procedurally enforceable; the primary
failure the oracle catches is **overclaim** — a claim made stronger than
its evidence supports.

- **Falsifiability** — claims are stated such that some specifiable
  observation would falsify them. β rejects a claim too vague to falsify
  (e.g. "the system behaves well"). A falsifiable restatement (e.g. "the
  measurement decreases monotonically under condition X") is required.
- **Diagnostic oracles** — where claims rest on measurements, the
  measurement procedure has a known failure-detection step: calibration
  check, signal-to-noise floor, control comparison, sentinel test. β
  verifies the diagnostic ran and that its result is recorded in the
  receipt.
- **Reproduction-from-clean** — β re-runs the producing command from a
  clean environment (per
  [`schemas/cdr/receipt.cue`](../../../../../schemas/cdr/receipt.cue)
  `#Reproduction`); the recorded output matches what the receipt claims.
  Required when `claim_status ∈ {observed, computed}`. The reproduction
  record carries the canonical command (with args + SHA) and the
  output-match boolean.
- **Citation integrity** — every claim derived from external work cites a
  resolvable reference; every cited reference supports the claim it is
  invoked for. β rejects receipts whose citations do not resolve or whose
  invocation overstates the cited result.
- **Data-policy compliance** — `data_refs` comply with the project's
  data-use policy (consent, anonymisation, retention, redistribution
  rights). β rejects a receipt whose data refs violate policy or whose
  data refs are under-specified (no manifest, no checksum).
- **Claim/evidence alignment** — the claim's strength is no stronger than
  the evidence supports. An `observed`-status claim has `data_refs` and a
  reproduction record; a `computed`-status claim has `method_refs` and a
  reproduction record; an `inferred`-status claim names the inferential
  step explicitly; a `hypothesized`-status claim is labelled as such.

Engineering's "compiles + tests pass" oracle does **not** apply. Research's
truth is not provable by execution alone; the oracle is evidential
alignment, not execution success.

### Field 3 — γ close-out artifact

What γ writes when a research cycle (a "wave") closes: the parent-facing
**research receipt**, typed by
[`schemas/cdr/receipt.cue`](../../../../../schemas/cdr/receipt.cue)
(`#CDRReceipt`). The typed receipt is the load-bearing γ artifact;
narrative wave reports are separate human-facing surfaces that cite the
receipt.

The receipt carries (per
[`ROLES.md §4a.3`](../../../../../ROLES.md#4a3-receipts-enforce-discipline-mechanically)
and the schema):

- `claim_refs` — which claims this receipt asserts (≥1).
- `data_refs` — dataset / mount / manifest / checksum (≥1).
- `method_refs` — script paths + commit SHA (≥1).
- `result_refs` — output file paths (≥1).
- `claim_status` — enum: `observed | computed | inferred | hypothesized |
  indeterminate`.
- `limitations` (optional) — explicit caveats.
- `reproduction` (optional, but required when `claim_status ∈ {observed,
  computed}`) — β re-ran the producing command; output matched.
- generic-kernel fields inherited from `schemas/cdd/#Receipt`:
  `protocol_id` (pinned to `cnos.cdd.cdr.receipt.v1`), `boundary_decision`,
  `verdict`, `protocol_gap_count`, `protocol_gap_refs`.

**Gate verdict vocabulary** carried by the receipt's `boundary_decision`
and `transmissibility` fields:

- **GO** — the claim is transmissible at full strength. Maps to
  `action: accept` + `transmissibility: accepted`.
- **REVISE** — the claim is transmissible after named revisions. Maps to
  `action: repair_dispatch` (β returns the receipt for revision).
- **NO-GO** — the claim is not transmissible. Further work required. Maps
  to `action: reject`.
- **INDETERMINATE** — evidence is insufficient to determine
  transmissibility. The receipt carries `claim_status: indeterminate` and
  the wave records what would change the determination.
- **BOUNDED-GO** — the claim is transmissible within named bounds: a
  "partial GO" or "subject-to-caveats" verdict. Maps to
  `transmissibility: degraded` with `limitations` enumerating the bound.

The receipt is what scope-`n+1` reads when the wave crosses the trust
boundary — per
[`schemas/cdd/README.md §"Scope-Lift Invariant"`](../../../../../schemas/cdd/README.md#scope-lift-invariant),
the closed wave is α-matter at the parent scope. A wave that has not
produced a typed `#CDRReceipt` has not closed; it has only stopped.

### Field 4 — δ cadence

How δ selects a new research gap and dispatches a new wave. **The cadence
is gate-transition-shaped, not release-shaped.** A CDR δ does not "cut a
release" or "tag a deploy" — a research wave concludes when a gate verdict
is recorded in the receipt. The receipt is the artifact; there is no
release-bundle artifact in the engineering sense.

The wave-transitions are:

- **Wave open** — δ selects a research gap (an unanswered question, an
  unverified claim, a needed dataset, a recurring measurement
  disagreement) and dispatches α + β.
- **Wave in progress** — α produces matter (Field 1); β reviews against
  the review oracle (Field 2); γ coordinates the loop.
- **Wave close** — γ writes the research receipt (Field 3) carrying the
  gate verdict. δ records the boundary decision; the wave is closed.

The trigger for δ to open the next wave is **gate-transition**, not
calendar or release-bundle:

- A **NO-GO** verdict may trigger a follow-up wave designed to close the
  gap that produced the NO-GO.
- A **GO** verdict may trigger a downstream synthesis wave (combining the
  newly-transmissible claim with prior claims).
- A **BOUNDED-GO** verdict may trigger a scope-expansion wave aimed at
  closing the bound.
- An **INDETERMINATE** verdict may trigger a measurement-design wave aimed
  at producing the missing evidence.

The cadence is not required to be fixed (per `ROLES.md §3`), but δ records
the trigger that opened each wave. An unstated trigger sequence is an
unowned protocol; receipts without parent-wave trigger references break
the receipt-stream signal ε reads (Field 5).

### Field 5 — ε iteration cadence

When ε assesses whether CDR itself is coherent, and what triggers ε work.

**Cadence: receipt-stream review over protocol gaps.** ε reads across many
`#CDRReceipt` instances and surfaces patterns that require protocol
patches. The cadence is **finding-triggered**, not calendar-triggered. The
trigger classes are research-failure classes (not engineering-failure
classes):

- **Missing data gates** — receipts shipping `observed` or `computed`
  claims with under-specified `data_refs` (no manifest, no checksum, no
  data-use compliance record). Pattern signals the data-policy gate is
  not being enforced; CDR's β oracle (Field 2 data-policy compliance)
  needs procedural sharpening.
- **Overclaiming** — receipts whose `claim_status: observed` evidence
  supports only `inferred`. Pattern signals β's claim/evidence-alignment
  oracle is too lenient or under-specified; the procedural skill (Sub 3
  `beta/SKILL.md`) needs a sharper claim-status determination procedure.
- **Unreproducible numbers** — receipts with `reproduction.output_match:
  false`, or receipts that should have a reproduction record (per Field 2)
  but lack one. Pattern signals reproduction-from-clean is not running or
  is failing silently; the procedural skill needs a stricter pre-publish
  gate.
- **Weak citation discipline** — receipts with claims that should cite
  external work but do not, or citations that do not support the claim
  they are invoked for. Pattern signals the citation-integrity oracle is
  drifting; the procedural skill needs explicit citation-validation
  checks.
- **Recurring oracle ambiguity** — the same review-oracle interpretation
  question appears across multiple cycles (β cannot decide; γ has to
  adjudicate). Pattern signals oracle under-specification; CDR.md or the
  per-role skills need additional procedural detail.
- **Construct drift** — a key research construct's definition shifts
  across receipts without explicit deprecation. Pattern signals the
  project binding (cph or downstream) needs a construct-stability artifact
  (a glossary, a deprecation registry).

ε's output is a CDR-shaped iteration artifact — Sub 3's `epsilon/SKILL.md`
declares the artifact spec. The trigger is "at least one finding tagged
with one of the classes above"; receipts with `protocol_gap_refs` populated
flow into ε's input.

### Field 6 — Actor collapse rule

Which roles may be held by the same actor; which collapses are permitted,
which are prohibited, and what signal warrants splitting.

**Never permitted for research claims:**

- **α = β within a single cycle.** Always prohibited. Research β reviews
  α's claim/evidence alignment for overclaim; α reviewing α's own claims
  is order-0 observation masquerading as order-1 (per
  [`ROLES.md §1`](../../../../../ROLES.md#1-the-role-ladder) +
  [`ROLES.md §4`](../../../../../ROLES.md#4-hats-vs-actors-roles-as-behavioral-contracts)).
  The structural independence is the mechanism by which review catches
  overclaim. No research-class waiver exists. Engineering-class collapse
  precedents (γ+α+β-collapsed-on-δ for mechanical refactor cycles) do
  **not** transfer to research-class claim transmission — research's
  failure mode (overclaim) is precisely the mode α cannot self-detect.

**Permitted with conditions:**

- **γ = δ for small project waves.** Allowed when the project is small
  enough that one actor reasonably holds both wave-coordination (γ) and
  gap-selection-across-waves (δ). The collapse is safe when γ's
  coordination authority does not compromise β's independence — γ sees
  both α and β but does not author either's matter. The signal for
  splitting: wave-volume crosses a threshold where γ's per-wave attention
  is incompatible with δ's cross-wave gap-selection load.
- **ε = δ until receipt-stream volume warrants split.** Allowed in
  small-protocol regimes where ε's receipt-stream review work does not
  justify a dedicated reviewer of the protocol. The signal for splitting:
  finding volume crosses a threshold where the protocol-iteration load
  competes with cross-wave gap-selection load.

**Project-specific stricter floors are permitted.** A project may impose
stricter rules in its `<project>/.cdr/` binding — for example, requiring
that for any externally-published claim, β be a distinct human reviewer
(not merely a separate session of the same actor). Such floors are
project-binding concerns; they belong in the project's `.cdr/POLICY.md` or
equivalent, not in CDR.md.

**The collapse signal is observable.** If a cycle's receipt cannot
distinguish α and β authorship (single signature, single session, no
review record), the cycle is structurally invalid for research-class
claim transmission. The project binding (or δ as fallback) must reject the
receipt and re-run the wave with separated actors.

---

## Empirical anchor

CDR's empirical anchor is [`usurobor/cph`](https://github.com/usurobor/cph)
— the research repository where the CDR pattern is realised in concrete
research waves, datasets, and field reports. CDR is anchored in cph but
not defined by cph (per
[`docs/papers/CCNF-AND-TYPED-TRUST.md §"Wave 4 — CDR bootstrap"`](../../../../../docs/papers/CCNF-AND-TYPED-TRUST.md)):
cph is the project binding (layer 4 of the five-layer enforcement chain);
CDR is the protocol overlay (layer 3); the layers are doctrinally
independent.

### Shape-compatibility claim

The six-field shape declared above describes cph's CDR practice without
contradiction. Spot-checks:

- **Field 1 (Matter type)** — cph's wave-shape (observation proposes →
  OpenCap translates → AI-analysis sorts → measurement decides → TSC
  measures the project itself) decomposes into the matter-type vocabulary
  (claims, hypotheses, methods, datasets, analyses, reports) without loss.
- **Field 3 (γ close-out gate verdicts)** — cph's gate verdicts include
  `bounded GO`, `partial GO`, `INDETERMINATE`, and
  `construct-survives-subject-to-caveats`. These map onto the CDR
  vocabulary as follows: `bounded GO` and `partial GO` are project-naming
  variants of **BOUNDED-GO**; `INDETERMINATE` maps directly;
  `construct-survives-subject-to-caveats` is a **BOUNDED-GO** with the
  bound named in `limitations`. No cph gate verdict requires a CDR
  vocabulary extension.
- **Field 4 (δ cadence)** — cph's wave sequencing is gate-transition-
  shaped, not release-shaped. Waves close on receipts; next-wave dispatch
  is driven by prior-wave verdicts. The compatibility holds without
  caveat.
- **Field 5 (ε trigger classes)** — cph's recurring protocol-gap signals
  (missing data manifests, citation drift, construct-shift across field
  reports) match the trigger classes named in Field 5.

The compatibility claim is sufficient for Sub 1's scope: the six-field
shape can host cph's practice without re-derivation.

### Canonical cph paths (by reference)

- `usurobor/cph:PROJECT.md` — project-level binding rules, gate-verdict
  definitions, data-use policy.
- `usurobor/cph:.cdr/waves/<wave>/` — wave-shaped research receipts and
  evidence; each wave's `receipt.md` exemplifies the typed
  `#CDRReceipt` shape per the schema's header comment.
- `usurobor/cph:.cdr/field-reports/` — narrative reports that cite the
  receipts.
- `usurobor/cph:.cdr/datasets/` — dataset manifests and provenance
  records.

These paths are cited; cph-local prose is **not** embedded in this file.
Embedding would violate the protocol/project separation declared in the
§Persona, Protocol, Project section.

### Detailed mapping deferred to Sub 4

Sub 4 of cnos#376 owns the full empirical-anchor doc: a surface-by-surface
mapping of cph's `.cdr/` artifact set to CDR.md's six-field structure and
to `#CDRReceipt`'s typed keys. Sub 1 (this cycle) establishes the citation
and the shape-compatibility claim; Sub 4 verifies the compatibility on
every cph artifact and records the mapping table.

---

## Related documents

Inherits, cites, or extends:

- [`ROLES.md`](../../../../../ROLES.md) — the role-cell grammar and the
  six-field instantiation contract that CDR realises. Specifically:
  - `§3` — six-field instantiation contract (the structure mirrored
    above).
  - `§4a` — five-layer enforcement chain (persona/protocol/project
    cited above).
  - `§4a.2` — loss-function distinction (CDR's research discipline).
  - `§4a.3` — receipt-enforces-discipline (the CDR receipt sketch
    realised in the schema).
  - `§7` — naming convention reserving `cdr` for research.
- [`schemas/cdr/receipt.cue`](../../../../../schemas/cdr/receipt.cue) — the
  typed γ close-out surface; `#CDRReceipt` is the parent-facing artifact
  CDR Field 3 names.
- [`schemas/cdd/README.md §"Architectural choice"`](../../../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split)
  — the upstream decision-record for the (a) split inherited above.
- [`docs/papers/CCNF-AND-TYPED-TRUST.md`](../../../../../docs/papers/CCNF-AND-TYPED-TRUST.md)
  — the L7 design essay; specifically `§"Separate persona, protocol, and
  project"` and `§"Wave 4 — CDR bootstrap"`.
- [`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../../../cnos.cdd/skills/cdd/CDD.md)
  — the reference instantiation (CDS-shaped). Structural exemplar for
  this file's shape; the research-discipline content here is **not**
  inherited from CDD.md, only the heading conventions.
- [cnos#376](https://github.com/usurobor/cnos/issues/376) — parent master;
  this file is Sub 1's deliverable.
- [cnos#388](https://github.com/usurobor/cnos/issues/388) — schema-side
  architectural-choice precedent.

---

## Non-goals

This document does **not**:

- Author the `cnos.cdr` package skeleton (`cn.package.json`, package
  `README.md`, root `SKILL.md`) — that is cnos#376 Sub 2.
- Author role-overlay skills (`alpha/SKILL.md`, `beta/SKILL.md`,
  `gamma/SKILL.md`, `delta/SKILL.md`, `epsilon/SKILL.md` for CDR) — that
  is cnos#376 Sub 3.
- Map cph's full `.cdr/` artifact set surface-by-surface — that is
  cnos#376 Sub 4.
- Redefine `ROLES.md` or the role ladder. CDR instantiates; ROLES
  governs.
- Author persona-identity content for `cn-rho`. Persona drafts live in
  the persona hub (`cn-rho/spec/PERSONA.md`), not here.
- Author project-binding content for cph. Project bindings live in
  `cph/PROJECT.md` and `cph/.cdr/`, not here.
- Specify runtime mechanics (dispatch, polling, git wiring, CI hooks).
  Those belong in role-overlay skills (Sub 3) and the persona/operator
  contracts.
- Author a `cn-cdr-verify` command or any tooling. The validator
  interface (V) is `cnos.cdd`'s scope, applied to `#CDRReceipt` via
  `protocol_id` dispatch per
  [`schemas/cdd/README.md §"protocol_id dispatch convention"`](../../../../../schemas/cdd/README.md#protocol_id-dispatch-convention).
