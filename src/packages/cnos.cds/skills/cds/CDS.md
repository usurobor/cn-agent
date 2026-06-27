<!-- sections: [preamble, architecture-choice, persona-protocol-project, six-field-contract, selection-function, development-lifecycle, coordination-surfaces, artifact-contract, mechanical-vs-judgment, review-clp, gate, assessment, closure, retro-packaging, large-file, empirical-anchor, related-documents, non-goals] -->
<!-- completed: [preamble, architecture-choice, persona-protocol-project, six-field-contract, selection-function, development-lifecycle, coordination-surfaces, artifact-contract, mechanical-vs-judgment, review-clp, gate, assessment, closure, retro-packaging, large-file, empirical-anchor, related-documents, non-goals] -->

# Coherence-Driven Software (CDS)

> CDS is a software-engineering-protocol instantiation of the generic
> role-scope ladder pattern. The pattern (α/β/γ/δ/ε roles, scope-escalation
> contract, six-field instantiation contract) is documented at
> [`ROLES.md`](../../../../../ROLES.md) at the repo root. CDS does not claim
> to have originated the role structure; it instantiates the structure for
> software-engineering work — code, tests, documentation, schemas, releases,
> deployments — under the engineering loss function (artifact improvement
> under repairable feedback).

**Version:** 0.1 (Subs 2–4 of [cnos#403](https://github.com/usurobor/cnos/issues/403))
**Status:** Draft — instantiation contract (Sub 2) + §Selection function and
§Development lifecycle (Sub 3 / cnos#408) + §Coordination surfaces and
§Artifact contract (Sub 4 / cnos#409, B-lite extracts). Sub 5 (review,
gate, assessment, closure, retro-packaging, non-goals, mechanical,
large-file), CDD.md marker cleanup (Sub 6), and empirical-anchor doc (Sub 7)
are downstream.
**Date:** 2026-05-22
**Placement:** `src/packages/cnos.cds/skills/cds/`
**Audience:** Engineering operators, engineering reviewers, project
maintainers binding CDS to a concrete software repo
**Scope:** The doctrinal contract for what CDS is. Declares the six
instantiation-contract fields, records the architectural-choice inheritance
from [cnos#388](https://github.com/usurobor/cnos/issues/388), and names the
persona/protocol/project boundary.

---

## 0. Purpose

CDS is the engineering discipline used to improve software artifacts
coherently under repairable feedback. Its loss function — per
[`ROLES.md §4a.2`](../../../../../ROLES.md#4a2-loss-function-distinction) —
is *artifact improvement under repairable feedback*. The primary failure
mode CDS must structurally resist is **the stalled loop**: code that drifts
farther from coherence with each round; tests that grow but do not bind;
documentation that lags the implementation; a working artifact that never
ships because each repair-round opens more gaps than it closes.

CDS shares the role-cell grammar (α/β/γ/δ/ε) with CDR (research), CDW
(writing), and CDA (analysis). What CDS does not share is the discipline
profile. Engineering optimises for *artifact improvement under repairable
feedback*; research optimises for *truth preservation under uncertainty*.
The role names are the same; the loss functions diverge; the receipts and
oracles diverge to enforce the divergent disciplines.

The asymmetry that drives the divergence — engineering's strong
correction-surface (build, test, CI, deployment-probe loops) versus
research's slow-truth-check surface — is the subject of
[`ROLES.md §4a.2`](../../../../../ROLES.md#4a2-loss-function-distinction);
that section is the source-of-truth for the distinction and is not
restated here. CDS's authoring discipline inherits the engineering side
of that distinction: the shared role grammar covers both regimes; the
per-protocol loss function determines which failure mode the discipline
optimises against; CDS optimises against the stalled loop.

This document is the instantiation contract for CDS per `ROLES.md §3`. It
declares the six fields any c-d-X protocol must declare. It does not
specify how to operate any individual role — that is the work of the
per-role skills downstream (cnos#403 Subs 3–5). It does not specify
the per-lifecycle-step procedure — that is also downstream (Subs 3–5,
guided by `docs/extraction-map.md`). It does not author the `cnos.cds`
package metadata — that landed in Sub 1 (cnos#406).

---

## Architecture choice

CDS inherits the architectural decision recorded at
[`schemas/cdd/README.md §"Architectural choice"`](../../../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split)
from [cnos#388](https://github.com/usurobor/cnos/issues/388) (Phase 2.5 of
cnos#366). cnos#388 decided the choice for **schemas**; cnos#376 AC7
extended the same decision-class to **skills/role-overlays**, and cnos#403
applies it a second time on the skills/role-overlay side for the engineering
realization (the cnos#376 application produced cnos.cdr; cnos#403 produces
cnos.cds). This section records the inheritance from CDS's side.

### The decision

**Option (a) — common constitution + per-protocol procedures.**

```text
common constitution                    per-protocol procedures
─────────────────────                  ────────────────────────
cnos.cdd / skills / cdd /              cnos.cds / skills / cds /
  CDD.md (doctrinal contract)            CDS.md (this file)
  COHERENCE-CELL.md                       (cites CCNF + CC by reference)
  COHERENCE-CELL-NORMAL-FORM.md
  RECEIPT-VALIDATION.md                    schemas/cds/receipt.cue
  alpha/SKILL.md (engineering α)           alpha/SKILL.md (engineering α — Sub 3)
  beta/SKILL.md  (engineering β)           beta/SKILL.md  (engineering β — Sub 3)
  gamma/SKILL.md (engineering γ)           gamma/SKILL.md (engineering γ — Sub 3)
  delta/SKILL.md (engineering δ)           delta/SKILL.md (engineering δ — Sub 3)
  epsilon/SKILL.md (engineering ε)         epsilon/SKILL.md (engineering ε — Sub 3)
```

The **common constitution** lives in `cnos.cdd`: the coherence-cell normal
form, the recursion equation, the role-ladder grammar, the generic
receipt-validation interface, the generic schemas. The common constitution
does **not** name engineering evidence by name; it types the shape.

The **per-protocol procedures** live in `cnos.cds` (and the sibling
`cnos.cdr`, plus future `cnos.cdw`, `cnos.cda`): the matter-type-specific
α/β/γ/δ/ε procedures, the review oracles, the close-out artifacts. Each
protocol package imports the `cnos.cdd` doctrine by reference and adds its
own discipline profile.

A subtlety in CDS's case: pre-cycle/402, the engineering-realization
procedures were resident inline in `cnos.cdd/skills/cdd/CDD.md`. Cycle/402
compressed CDD.md to the CCNF spine and marked the engineering-realization
content as `pending cds extraction`. cnos#403 (this sub's parent tracker)
is the extraction wave: Sub 1 shipped the cnos.cds package skeleton plus
`docs/extraction-map.md`; Sub 2 (this cycle) ships the doctrinal contract;
Subs 3–5 migrate the lifecycle/artifacts/review/gate procedures onto the
contract's surfaces; Sub 6 removes the `pending cds extraction` markers
from CDD.md once Subs 3–5 land. The (a) split is therefore not a fresh
authoring act for CDS — it is the structural endpoint of a multi-cycle
wave whose first decision was made at cnos#388 and whose first skill-side
application was cnos#376.

### Option (b) — common protocol-agnostic skill + domain overlay (rejected)

The rejected alternative would put α/β/γ/δ/ε as protocol-agnostic skills in
`cnos.cdd` with thin per-domain overlays (CDS/CDR/CDW). This is rejected
for the same five reasons cnos#388 used for schemas, transposed to skills.
The CDS-side reading of each rationale:

1. **The role grammar is the constitution; the loss function is the
   procedure.** Engineering α and research α share α's verb ("produces")
   but the failure modes diverge sharply: engineering's stalled-loop
   failure is α-internal (α can detect it by running the tests); research's
   overclaim failure is α-blind (α cannot detect that its claim outran its
   evidence without an outside frame). The discipline profile cannot be
   common because the failure-mode geometry differs.

2. **Clarity at the protocol boundary.** An engineering α reading
   `cnos.cds/skills/cds/alpha/SKILL.md` (Sub 3 deliverable) sees only the
   engineering-α discipline — no indirection through a generic skill plus a
   domain overlay. The corresponding research α reading
   `cnos.cdr/skills/cdr/alpha/SKILL.md` sees only research-α. The
   boundary is the file location, not a configuration flag.

3. **Mechanical generic-vs-domain boundary.** Different files plus different
   package directories make the boundary structural, not semantic. Reuse
   happens by reading + reference, not by inheritance gymnastics. For
   software-class matter this matters especially because the engineering
   substrate (CI, build tools, deployment effectors) already has strong
   filesystem conventions; the skill layout honours those conventions
   rather than hiding them behind an overlay-resolution mechanism.

4. **Future c-d-X generalizes mechanically.** Adding cdw (writing) means
   adding `cnos.cdw/skills/cdw/CDW.md` + its overlays. CDS's authorship
   does not need to know about cdw, and cdw's authorship does not need to
   know about CDS. The common constitution stays clean.

5. **Decision-once-applied-thrice.** cnos#388 decided (a) for schemas;
   cnos#376 applied the same reading to skills for the research case;
   cnos#403 applies the same reading to skills for the engineering case.
   The decision is recorded once and applied N times. CDS.md does not
   re-derive the rationale — it cites the decision-record and notes the
   inheritance.

### Design source

The decision rationale and the schema-level precedent are sourced from:

- [`schemas/cdd/README.md §"Architectural choice"`](../../../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split)
  — the cnos#388 decision-record. Rationale (1)–(5) for the schema case;
  inherited here by reference.
- [`docs/papers/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol,
  and project"`](../../../../../docs/papers/CCNF-AND-TYPED-TRUST.md)
  — the doctrinal source for treating persona, protocol, and project as
  three distinct layers.
- [`docs/papers/CCNF-AND-TYPED-TRUST.md §"Wave 5 — CDS extraction by
  reference"` and `§"Wave 7 — final CDD.md rewrite"`](../../../../../docs/papers/CCNF-AND-TYPED-TRUST.md)
  — declare the wave-shape for the engineering-realization extraction. The
  CCNF spine landed at Phase 7 of cnos#366 (cycle/402); the cds extraction
  is the post-CCNF wave at cnos#403; CDS.md is its doctrinal-contract step.
- [`ROLES.md §4a.2`](../../../../../ROLES.md#4a2-loss-function-distinction)
  — the loss-function distinction that makes engineering-discipline and
  research-discipline differ at the procedural layer while sharing the
  role-cell kernel.
- [`src/packages/cnos.cdr/skills/cdr/CDR.md §"Architecture choice"`](../../../cnos.cdr/skills/cdr/CDR.md)
  — the parallel record on the research side. CDR.md recorded the same
  inheritance from cnos#388 + cnos#376; CDS.md records it from cnos#388 +
  cnos#403. The two records are siblings; neither supersedes the other.

### Cross-reference

- [cnos#388](https://github.com/usurobor/cnos/issues/388) — schema-side
  decision. The (a) decision recorded here for engineering skills is its
  third application.
- [cnos#376](https://github.com/usurobor/cnos/issues/376) — research-side
  skills application (first skill-side application).
  [cnos#376 AC7](https://github.com/usurobor/cnos/issues/376) names the
  inheritance pattern Sub 1 of that wave used.
- [cnos#403](https://github.com/usurobor/cnos/issues/403) — engineering-side
  skills application (this wave). Sub 2 (cnos#407) is the cycle landing
  this document.

---

## Persona, Protocol, Project

CDS is a **protocol overlay** — layer 3 of the five-layer enforcement chain
declared at
[`ROLES.md §4a`](../../../../../ROLES.md#4a-persona-operator-protocol-project-receipt--the-five-layer-enforcement-chain).
This section names CDS's relationship to the layers above and below it.
The boundary is doctrinal, not operational; operational mechanics (dispatch,
polling, repo wiring, branch creation, merge effection) belong in role-overlay
skills (Sub 3) and in persona/operator contracts (separate hubs).

### Layer 1 — Persona (canonical home: `cn-sigma/spec/PERSONA.md` + `cn-sigma/spec/OPERATOR.md`)

A persona is *what kind of mind is doing the work*. For engineering, the
named exemplar is **Sigma** — the engineering persona declared at
[`ROLES.md §4a.4`](../../../../../ROLES.md#4a4-worked-example--sigma-engineering-and-rho-research).
Sigma's discipline profile is sourced from that section (action-biased,
correction-surface-driven, debt-recording rather than block-on-perfection);
CDS.md does not restate the profile, only cites it. Other engineering
personas are admissible provided they declare a compatible discipline
profile in their hub's `PERSONA.md`.

**Sigma may play δ for CDS cycles.** The persona enacts the role; the role
does not enact the persona. Sigma may also play δ for non-CDS engineering
coordination work, or play γ for a project that runs CDS + CDR in parallel.

**Sigma is not CDS.** The persona hub `cn-sigma` lives outside this package
(and is itself forthcoming — the current cnos engineering work is done by
operator-class agents whose persona contract is the
`cnos.cdd/skills/cdd/operator/SKILL.md` surface, which is generic-substrate
material slated to migrate into Sigma's hub when the hub bootstraps). CDS.md
does not author persona content; CDS.md does not state "I am Sigma" or "my
voice." If a future engineering persona is named alongside Sigma, the
protocol overlay remains unchanged — only the new persona hub's
`PERSONA.md` is added.

### Layer 2 — Protocol overlay (canonical home: `cnos.cds/skills/cds/`)

The protocol overlay specifies *what counts as valid work in the
engineering domain*. CDS's protocol overlay is this package: CDS.md (this
file) declares the doctrinal contract; the role-specific skills under
`cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` (Sub 3) will
declare each role's procedural discipline.

**CDS is not a persona.** Anyone (or any agent) can author a CDS cycle if
they enact the engineering-loss-function discipline; the role-cell shape
is portable. CDS's overlay specifies what the role-cell looks like for
software-engineering work; it does not specify who plays the roles.

**CDS defines α/β/γ/δ/ε engineering role overlays independent of any one
persona hub.** The CDS α/β/γ/δ/ε skills (Sub 3) describe the engineering
discipline at each role; they cite Sigma's discipline profile by reference
when describing what kind of persona is typical for CDS work, but they do
not embed persona-identity content.

### Layer 3 — Project binding (canonical home: `<project>/.cds/` or `<project>/cds/`)

A project binds CDS's roles to *concrete code, tests, branches, gates, CI
pipelines, and release effectors*. The canonical example — and CDS's
primary empirical anchor (see §Empirical anchor below) — is `usurobor/cnos`
itself: cnos's repo carries the wave of CDS cycles (cycle/364 through
cycle/406 and the in-flight wave); cnos's `.cdd/unreleased/{N}/` and
`.cdd/releases/{X.Y.Z}/{N}/` directories carry the per-cycle close-out
artifacts; cnos's CI configuration, `cn` build tool, release-effector
scripts, and merge gates are the project-level operational substrate.

The `.cdd/`-rooted directory naming is historical (the convention predates
the cds extraction); the eventual rename to `.cds/` is a separate post-#403
coordination problem (see `docs/extraction-map.md §14 "Open questions"`).
CDS.md uses the destination naming `<project>/.cds/` to declare the
canonical home; current cnos cycles continue to use `.cdd/unreleased/{N}/`
until the filesystem rename lands.

**CDS.md does not author project-specific bindings.** Project-specific gate
names, CI thresholds, deployment targets, release-cadence rules, and
configuration-floor caps live in the project's `.cds/` directory (or
`.cdd/` until the rename). CDS.md declares the doctrinal shape; projects
bind concrete artifacts.

### Why this matters

Conflating any two of the three layers produces drift, exactly as
[`ROLES.md §4a`](../../../../../ROLES.md#4a-persona-operator-protocol-project-receipt--the-five-layer-enforcement-chain)
warns. Specifically for CDS:

- **Persona content inside CDS.md** would make CDS un-reusable across
  engineering personas — a future engineering persona that is not Sigma
  would have to fork the protocol.
- **Project-specific gate names inside CDS.md** would make CDS
  un-reusable across engineering projects — a downstream engineering repo
  that is not cnos would have to fork the protocol. (The cnos repo is the
  bootstrap binding; downstream tenants must be able to bind CDS without
  inheriting cnos-specific gate names.)
- **Doctrinal content inside cnos's project binding** (the `.cdd/` or
  `.cds/` directory) would trap the engineering discipline inside one
  project — Sub 7's empirical-anchor mapping work would become an
  authoring task rather than a binding task.

The three canonical homes (`cn-sigma/spec/`, `cnos.cds/skills/cds/`,
`<project>/.cds/`) are the structural separation that prevents the drift.

---

## Six-field instantiation contract

Per [`ROLES.md §3`](../../../../../ROLES.md#3-instantiation-contract), every
c-d-X protocol must declare six fields. CDS declares them below, each in
engineering-loss-function language. The shape follows the role-ladder
contract; the language is software-engineering, not research.

### Field 1: Matter type

α produces **software artifacts under repairable feedback**: source code,
tests, documentation, schemas, CI definitions, runtime contracts, and
design notes.

- **Source code** — the executable artifact. Programs, libraries, scripts,
  CLI binaries, services. Source is α-matter when it implements a contract
  step, refactors an existing artifact under a contract, or generates an
  artifact that downstream tooling consumes. Source carries its build
  invocation by reference (in `cn.package.json`, `Makefile`, `go.mod`, etc.)
  — the build artifact itself is derived, not authored.
- **Tests** — the executable oracle. Unit tests, integration tests, golden
  tests, fixture suites. Tests are α-matter when they encode an AC, when
  they pin a regression, or when they establish a behavioural contract a
  downstream consumer relies on. Tests are *not* a separate artifact class
  from source; they are source whose role is to discriminate other source.
- **Documentation** — the prose surface. Skill files (`SKILL.md`),
  doctrine files (`CDD.md`, this file), READMEs, design essays, RFCs,
  inline comments at load-bearing decision points. Documentation is
  α-matter when it declares the contract a future cycle binds against, when
  it records a design decision whose rationale is non-obvious, or when it
  resolves an ambiguity surfaced by ε.
- **Schemas** — the typed-trust surface. CUE schemas (`schemas/cds/*.cue`),
  JSON schemas (`tools/validate-*/`), CUE module declarations, frontmatter
  schemas. Schemas are α-matter when they declare what V (the validator)
  consumes, when they pin a wire format, or when they ratchet a structural
  invariant a future cycle inherits.
- **CI definitions** — the automation surface. GitHub Actions workflows,
  build scripts, pre-commit hooks, lint configurations, mechanical-check
  scripts. CI is α-matter when it adds a mechanical check that catches a
  class of finding ε surfaced, when it tightens an existing check's
  enforcement, or when it adds a check the contract requires β to verify.
- **Runtime contracts** — the loadable manifest surface. `cn.package.json`,
  command manifests, hook definitions, package-discovery schemas. Runtime
  contracts are α-matter when they declare what the runtime substrate
  loads, when they pin a discovery convention, or when they extend the
  loadable-surface schema.
- **Design notes** — the contract-shape surface. `cdd/design/SKILL.md`
  artifacts (invariant/volatile/boundary triplets), issue-body
  implementation contracts (per cnos#393 δ-as-architect), L7 essays in
  `docs/papers/`. Design notes are α-matter when they pin a
  contract-shape decision a future cycle binds against; the design itself
  is matter, not metadata about matter.

The matter type determines what β reviews (Field 2), what γ writes to close
out (Field 3), and what ε observes across many cycles (Field 5).

**CDS does not produce research artifacts as matter.** Research produced
in service of engineering work — a coherence-delta measurement, an
empirical anchor, a falsifiable hypothesis about whether a refactor will
land — is **research-class evidence** consumed by reference (cited from
the cycle's design notes), not authored under CDS. A cycle whose primary
deliverable is a knowledge claim under uncertainty is structurally a CDR
cycle, not a CDS cycle; the protocol dispatch boundary (per
[`ROLES.md §4a.2`](../../../../../ROLES.md#4a2-loss-function-distinction))
routes the work to the matching discipline. The reverse routing also
holds: a CDR wave that ships executable measurement scripts as its
primary deliverable is structurally a CDS cycle whose result feeds the CDR
wave by reference.

The Sub-3-vs-Field-1 line: Field 1 owns the matter-type taxonomy (what α
produces); Sub 3 owns the per-step lifecycle procedure (how α produces it
across the 0–13 cycle steps). The taxonomy here is what Sub 3 binds
operational detail against; the operational detail does not live here.

### Field 2: Review oracle

How β determines that α's matter closes the declared engineering gap. Each
oracle is mechanically checkable or procedurally enforceable; the primary
failure the oracle catches is **the stalled loop** — α producing matter
that drifts farther from coherence with each round, rather than binding on
the contract β can mechanically verify.

- **Compilation, type-check, and build** — the artifact compiles, types
  resolve, the build invocation declared by the package manifest succeeds.
  β rejects matter that does not build (no merge across a broken build is
  permitted). For non-compiled languages, the equivalent is the language's
  static-check pass (lint clean at the cycle's declared lint floor;
  type-stub validation where applicable).
- **Tests pass** — the cycle's declared test surface passes from a clean
  environment. At minimum the cycle-touched test suites pass; where
  mechanically tractable the full suite passes (the project binding sets
  the floor — small cycles may exempt unrelated long-running suites; the
  exemption is recorded in the cycle's receipt). β rejects matter whose
  declared tests do not pass.
- **AC verification** — each AC has a mechanical or read-check oracle
  named in the issue body (per cnos#393 δ-as-architect: ACs are
  contract-shape, not implementation discretion). β verifies each
  oracle ran and recorded its result. β rejects a cycle that closes
  without per-AC oracle evidence.
- **No regressions** — pre-existing tests pass; pre-existing mechanical
  checks pass; documented behaviours are preserved (or their change is
  named in the contract as intentional). β rejects matter whose
  side-effects degrade adjacent surfaces.
- **Implementation-contract coherence** — per
  [cnos#393](https://github.com/usurobor/cnos/issues/393) Rule 7, α's
  matter satisfies the implementation-contract axes pinned by δ at
  dispatch: language, CLI integration target, package scoping,
  existing-binary disposition, runtime dependencies, JSON/wire contract,
  backward compatibility. β verifies each axis. β rejects matter that
  drifts from a pinned axis without an explicit contract-amendment record.
- **Evidence-binding** — per the CCNF kernel's evidence-binding rule (see
  [`COHERENCE-CELL-NORMAL-FORM.md §"Kernel"`](../../../cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md)),
  the receipt carries typed evidence-refs to the cycle's artifacts; γ
  binds the evidence at close-out; V dereferences the refs at validation.
  β verifies the receipt's refs resolve (the artifacts exist at the cited
  paths; the diff ref matches the cycle's diff). β rejects a receipt whose
  refs do not resolve or whose evidence-binding is incomplete.

Research's "claim/evidence alignment" oracle does **not** apply at this
layer. Engineering's truth is provable by execution (within the engineering
substrate's correction window): the artifact works or it does not; the test
passes or it does not. The β oracle is execution-and-contract alignment,
not evidential alignment. Engineering matter that *makes* an empirical
claim (a benchmark assertion, a behavioural measurement) crosses into
research-class matter and is routed per the dispatch boundary above; CDS β
does not adjudicate research-class claims directly.

The Sub-5-vs-Field-2 line: Field 2 owns the oracle taxonomy (what β
catches); Sub 5 will detail the review CLP (the TERMS / POINTER / EXIT
form), the reviewer ask list (α/β/γ scores, weakest-axis diagnosis,
concrete patch suggestions), the iterate-or-converge verdict shape, and
the per-finding disposition workflow. Field 2 owns the oracle taxonomy;
Sub 5 owns the per-oracle operational detail.

### Field 3: γ close-out artifact

What γ writes when an engineering cycle closes: a **per-cycle close-out
artifact set** plus the typed **`#CDSReceipt`** that ties the set to V
(the validator).

**The close-out artifact set** (per the cnos cycle convention; lifecycle
detail is Sub 4 territory):

- `self-coherence.md` — α's pre-review self-check; AC-by-AC mechanical pass
  + β-Rule-7 implementation-contract conformance.
- `beta-review.md` — β's review verdict (typically APPROVED, fix-needed
  rounds R1/R2/…, or NO-GO); per-AC oracle verification; β's findings
  (binding and non-binding).
- `alpha-closeout.md` — α-level findings from the cycle (often empty for
  mechanical cycles; populated when α surfaces a skill-gap or
  protocol-gap during production).
- `beta-closeout.md` — β-level findings from the cycle (typically
  reviewer-quality observations; sometimes protocol gaps surfaced during
  review).
- `gamma-closeout.md` — γ's closure summary: finding triage, disposition,
  configuration-floor declaration, coherence delta, next-action handoff
  to the operator.
- `cdd-iteration.md` — ε's per-cycle protocol-iteration artifact. Per the
  cycle/401 cadence rule, required only when `protocol_gap_count > 0`; a
  courtesy empty-findings stub is acceptable when count is 0. (Naming:
  the file is currently `cdd-iteration.md` for compatibility with the
  pre-rename empirical anchor; Sub 5 will land the rename to
  `cds-iteration.md` as part of the lifecycle migration.)

**The typed receipt** is `#CDSReceipt` declared at
[`schemas/cds/receipt.cue`](../../../../../schemas/cds/receipt.cue). It is
the parent-facing typed surface CDS Field 3 names. The receipt carries
the typed evidence-refs the close-out set provides — `self_coherence`,
`beta_review`, `alpha_closeout`, `beta_closeout`, `gamma_closeout`,
`evidence_root`, plus the cycle's diff and any CI refs — and inherits the
generic kernel's required fields (`protocol_id` pinned to
`cnos.cdd.cds.receipt.v1`, `boundary_decision`, `verdict`,
`protocol_gap_count`, `protocol_gap_refs`) from `schemas/cdd/#Receipt`
via CUE unification.

**The relationship between the two surfaces:** the narrative close-out
files are the cycle-local close-out artifact set; the typed
`#CDSReceipt` is the parent-facing typed surface that ties the close-out
files to V. The two surfaces compose: the typed receipt's `evidence_refs`
point at the close-out files; V dereferences the refs to validate the
receipt against the cycle's contract. Without the close-out files the
receipt's refs do not resolve; without the typed receipt the close-out
files have no parent-facing structural surface for V to read.

**Gate verdict vocabulary** carried by the receipt's `boundary_decision`
and the cycle's narrative close-out:

- **GO / accept** — the cycle is mergeable; the closed cell projects as
  α-matter at the parent scope (the next cycle / the release). Maps to
  `action: accept` + `transmissibility: accepted` in the receipt.
- **REVISE / repair_dispatch** — the cycle is not yet mergeable; β
  identified one or more binding findings; a fix round is required. Maps
  to `action: repair_dispatch` in the receipt; γ either re-dispatches α
  for the fix or escalates to δ for contract amendment.
- **NO-GO / reject** — the cycle does not close at this scope. Further
  work, contract revision, or cancellation is required. Maps to
  `action: reject`.
- **OVERRIDE / override** — the cycle merges despite a non-PASS verdict
  under explicit δ-authority. Maps to `action: override`;
  `transmissibility` is computed as `degraded` per the structural override
  rule in `schemas/cdd/receipt.cue`. The override block is the structural
  signal every downstream consumer (next cycle's α; release's gate) must
  detect.

The receipt is what scope-`n+1` reads when the cycle crosses the trust
boundary — per
[`schemas/cdd/README.md §"Scope-Lift Invariant"`](../../../../../schemas/cdd/README.md#scope-lift-invariant),
the closed cycle is α-matter at the parent scope. A cycle that has not
produced a typed `#CDSReceipt` (or its narrative analogue when V is not
yet wired in for that cycle's project binding) has not closed; it has only
stopped.

The Sub-4-vs-Field-3 line: Field 3 owns the artifact-set taxonomy + the
typed-receipt citation; Sub 4 will detail the artifact contract — the
Artifact Location Matrix (canonical paths per surface), the role/artifact
ownership matrix, the CDS Trace format, the bootstrap procedure, the
frozen-snapshot rule. Field 3 owns the artifact taxonomy; Sub 4 owns the
per-artifact operational detail.

### Field 4: δ cadence

How δ selects a new engineering gap and dispatches a new cycle.

**The cadence is per-cycle, with release-shaped bundling.** Unlike CDR's
research cadence (which is strictly gate-transition-shaped because research
has no release artifact), CDS's cadence has three nested time-scales:

- **Per-cycle:** each cycle opens when δ dispatches a gap and closes when
  δ records the boundary decision on the receipt. The cycle is the unit
  of CDS work; merge is the boundary effection per cycle; cycles that
  fail the boundary do not merge.
- **Per-release:** when the unreleased bundle (one or more merged cycles
  accumulating under `.cds/unreleased/`, or `.cdd/unreleased/` until the
  rename) reaches a coherent shipping point — PRA-readiness,
  RELEASE.md author-ability, no in-flight cycles blocking the release
  — δ tags a release; the release-effector mechanics handle the
  tag/build/deploy.
- **Per-wave:** for multi-cycle waves (like cnos#403 itself), δ
  dispatches the wave's sub-cycles in dependency order, tracks the wave's
  shape against the wave manifest, and closes the wave when its
  sub-cycles' close-outs satisfy the wave's collective AC set.

**The inward-membrane discipline** (per cnos#393 δ-as-architect) runs at
each cycle's dispatch: δ pins the implementation contract (the seven
axes — language, CLI target, package scoping, existing-binary disposition,
runtime deps, JSON/wire, backward compat) in the issue body before α
begins. The pinning is the contract-shape act; α binds against the pinned
axes; β verifies the binding per Field 2's implementation-contract
coherence oracle.

**The outward-membrane discipline** runs at merge: δ records the boundary
decision on the receipt; the merge into main (or the cycle's target
branch) is the boundary effection. The kernel rules governing what δ
may consult, the verdict/decision composition, and the repair-vs-reinvoke
paths when δ doubts V are sourced from
[`COHERENCE-CELL-NORMAL-FORM.md §"Kernel"` step 5 and §"Cell Outcomes"](../../../cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md);
CDS Field 4 inherits them without restatement. The CDS-side qualifier:
the boundary effection on engineering substrate is concretely a git
merge into the project's target branch (typically `main`) plus any
release-effector mechanics the project binding wires in.

Cadence triggers — what opens the next cycle:

- A **GO** verdict on a wave's prerequisite cycle may unblock dispatchable
  sub-cycles (the wave's dependency edges).
- A **REVISE** verdict triggers an in-flight repair round (same cycle, new
  α turn) rather than a new cycle.
- A **NO-GO** verdict may trigger a follow-up cycle designed to close the
  gap that produced the NO-GO (with contract revisions per the NO-GO's
  findings), or may close the work-stream entirely if the gap is judged
  not worth closing.
- An **OVERRIDE** decision is followed up by a debt-cycle (filed at
  override-time) that closes the override's gap durably; the debt-cycle's
  receipt is what removes the override's degraded mark from the
  downstream-consumer signal.
- A **release-readiness** signal (PRA-ready, no in-flight cycles, release
  bundle coherent) triggers the per-release δ work (tag, build, deploy).

The cadence is not required to be fixed (per `ROLES.md §3`), but δ records
the trigger that opened each cycle. An unstated trigger sequence is an
unowned protocol; receipts without parent-cycle / parent-wave trigger
references break the receipt-stream signal ε reads (Field 5).

The Sub-3-vs-Field-4 line: Field 4 owns the cadence taxonomy (per-cycle +
per-release + per-wave; trigger classes); Sub 3 will detail the lifecycle
state machine S0–S12, the branch rule (`cycle/{N}` canonical; γ creates
from `origin/main` before dispatch), the branch pre-flight, and the
skill-loading tiers. Field 4 owns the cadence taxonomy; Sub 3 owns the
per-step lifecycle operational detail.

### Field 5: ε iteration cadence

When ε assesses whether CDS itself is coherent, and what triggers ε work.

**Cadence: receipt-stream review over engineering protocol gaps.** ε reads
across many `#CDSReceipt` instances (or their narrative close-out analogues
where V is not yet wired in) and surfaces patterns that require protocol
patches. The cadence is **finding-triggered**, not calendar-triggered, per
the generic ε doctrine at
[`ROLES.md §4b`](../../../../../ROLES.md#4b-generic-%CE%B5--the-protocol-iteration-role).

**Gap classes** — per the
[`ROLES.md §4b.3`](../../../../../ROLES.md#4b3-gap-class-taxonomy-instantiation-pattern)
`{protocol}-{axis}-gap` naming convention, CDS declares four classes:

- **`cds-skill-gap`** — a procedural skill is underspecified or wrong. A
  CDS skill file (the SKILL.md for a role, lifecycle step, or per-area
  procedure) did not give the operator enough to avoid the finding. The
  patch is a skill edit; the gap is closed when the edited skill prevents
  the same class of finding in a future cycle.
- **`cds-protocol-gap`** — CDS doctrine itself drifted. The contract
  (this file — CDS.md) needs an edit: a field's taxonomy missed a case;
  a cadence rule did not anticipate an actor configuration; an oracle
  was under-specified. The patch is a CDS.md edit.
- **`cds-tooling-gap`** — tooling absent, wrong, or unavailable. A
  mechanical check the cycle needed did not exist (or existed but did not
  run, or ran but did not catch). The patch is tooling work: a new
  validator, a new lint rule, a new pre-commit hook, a new mechanical
  check in CI.
- **`cds-metric-gap`** — measurement missing or wrong. A coherence-delta
  measurement the cycle relied on did not produce the right signal: the
  α-axis / β-axis / γ-axis scores were ambiguous; the cycle-iteration
  trigger thresholds (rounds > 2; mechanical ratio > 20%) misfired; the
  receipt-stream view ε reads lacked a derivable field. The patch is a
  measurement fix or a metric definition.

**Naming-class transition note:** the current empirical-anchor cycles (the
cnos #364–#406 wave) use `cdd-*-gap` because they predate the cds
extraction. CDS.md declares the new `cds-*-gap` names as the canonical
CDS taxonomy. The transition is organic: future CDS cycles use
`cds-*-gap`; pre-#403 cycles' `cdd-*-gap` markers stay as historical
record. The two name-sets cover the same gap geometry; the rename is a
package-attribution rename, not a doctrine change.

**Iteration artifact** — per
[`ROLES.md §4b.4`](../../../../../ROLES.md#4b4-iteration-artifact-and-cadence-rule):

- File: `cds-iteration.md` (canonical name; the empirical-anchor cycles
  use the pre-rename `cdd-iteration.md` until Sub 5 lands the rename).
- In-cycle path: `<project>/.cds/unreleased/{N}/cds-iteration.md` (or
  `.cdd/unreleased/{N}/cdd-iteration.md` for the pre-rename binding).
- Post-release path:
  `<project>/.cds/releases/{X.Y.Z}/{N}/cds-iteration.md`.
- Aggregator: `<project>/.cds/iterations/INDEX.md` (one row per cycle
  that produced findings; the cycle/401 courtesy convention appends a
  row even for zero-finding cycles).

The cadence rule (per
[`ROLES.md §4b.4`](../../../../../ROLES.md#4b4-iteration-artifact-and-cadence-rule)):
the iteration artifact is **required only when `protocol_gap_count > 0`**.
The cycle/401 courtesy convention permits writing a zero-findings stub for
traceability; the stub records `protocol_gap_count: 0` explicitly and is
linked in the aggregator with the cycle's finding count of 0.

**Trigger classes — what surfaces a CDS protocol gap:**

- **Review churn** — review rounds > 2 on a single cycle. Signals that
  the contract was under-specified, the AC oracle was ambiguous, or the
  α-skill was under-specified for the matter class.
- **Mechanical-finding overload** — mechanical-class findings > 20% of
  total findings on a single cycle. Signals tooling-class gap: a check
  that should be in CI is being run by β manually.
- **Avoidable tooling / environment failure** — the cycle was blocked by
  a tooling or environment failure that a CDS tooling artifact could
  have prevented. Signals `cds-tooling-gap`.
- **Loaded skill failed to prevent a finding** — a finding surfaced
  despite the relevant skill being loaded; the skill's procedure did
  not catch the class of error. Signals `cds-skill-gap`.
- **AC oracle ambiguity recurrence** — the same per-AC oracle
  interpretation question appears across multiple cycles. Signals
  `cds-skill-gap` (per-AC oracle authoring skill) or `cds-protocol-gap`
  (Field 2 oracle taxonomy under-specifies the class).
- **CI red post-merge** — a CI failure on the merge commit that the
  cycle's β oracle (Field 2 "no regressions") should have caught.
  Signals `cds-tooling-gap` (CI check missing) or `cds-skill-gap`
  (β procedure missed it).

ε's output is the iteration artifact above. Sub 3's `epsilon/SKILL.md`
declares the per-finding shape, the disposition workflow (ship-now /
next-MCA / no-patch per
[`ROLES.md §4b.5`](../../../../../ROLES.md#4b5-mca-discipline)), and the
aggregator-update procedure.

The Sub-5-vs-Field-5 line: Field 5 owns the gap-class taxonomy + the
cadence rule + the trigger classes; Sub 5 will detail the per-class
operational workflow, the per-finding shape, the MCA discipline mechanics,
the aggregator format. Field 5 owns the class taxonomy; Sub 5 owns the
per-class operational detail.

### Field 6: Actor collapse rule

Which roles may be held by the same actor; which collapses are permitted,
which are prohibited, and what signal warrants splitting.

**Never permitted for substantive software work:**

- **α = β within a single cycle.** Always prohibited. Engineering β
  reviews α's matter for stalled-loop and contract-coherence failures.
  The order-of-observation argument for why α cannot self-review (the
  structural independence that makes review add information α could
  not produce) is the subject of
  [`ROLES.md §1`](../../../../../ROLES.md#1-the-role-ladder) and
  [`ROLES.md §4`](../../../../../ROLES.md#4-hats-vs-actors-roles-as-behavioral-contracts);
  CDS inherits the argument without restatement. The engineering-side
  qualifier: the engineering substrate's fast correction surface
  (build, test, CI) catches some failure modes structurally, but the
  contract-coherence failures Field 2 names (implementation-contract
  oracle, no-regressions oracle on adjacent surfaces, evidence-binding
  oracle) require an outside frame; the correction surface does not
  eliminate the need for β independence on substantive software work.

**Permitted with conditions:**

- **β-α-collapse-on-δ for skill/docs-class cycles** (per the
  breadth-2026-05-12 wave manifest precedent; per the breadth of
  cycles #388 through #406 dispatched under this pattern). When the
  matter is structural-mirror, migration, or contract-authoring — no novel
  substantive code; mechanical correctness verifiable from the spec; the
  failure modes the α-β independence catches are reduced by the absence
  of novel executable surface — δ may dispatch γ + α + β collapsed on a
  single agent. The
  [`release/SKILL.md §3.8`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  configuration-floor clause then caps γ-axis and β-axis at A-; the
  collapse is acknowledged in the cycle's `gamma-scaffold.md` "Dispatch
  shape" section. **This collapse is not permitted for substantive
  software work** — new features, bug fixes, refactors of executable
  surface, schema changes that affect runtime contracts, or any cycle
  whose primary matter is novel α-internal code.
- **γ = δ for small-project regimes.** Allowed when the project is small
  enough that one actor reasonably holds both cycle-coordination (γ) and
  gap-selection-across-cycles (δ). The collapse is safe when γ's
  coordination authority does not compromise β's independence — γ sees
  both α and β but does not author either's matter. The signal for
  splitting: cycle-volume crosses a threshold where γ's per-cycle
  attention is incompatible with δ's cross-cycle gap-selection load.
- **ε = δ until receipt-stream volume warrants split.** Allowed in
  small-protocol regimes where ε's receipt-stream review work does not
  justify a dedicated reviewer of the protocol. Per
  [`ROLES.md §4b.6`](../../../../../ROLES.md#4b6-%CE%B5s-relationship-to-%CE%B4-collapse-rule),
  the collapse is one of the *safe* collapses. The signal for splitting:
  finding volume crosses a threshold where the protocol-iteration load
  competes with cross-cycle gap-selection load.

**Project-specific stricter floors are permitted.** A project may impose
stricter rules in its `<project>/.cds/POLICY.md` (or `.cdd/POLICY.md`
until the rename) — for example, requiring that for any
production-release cycle, β be a distinct human reviewer (not merely a
separate session of the same actor) regardless of cycle class; or
prohibiting the β-α-collapse-on-δ pattern entirely for repos with strict
audit obligations. Such floors are project-binding concerns; they belong
in the project's `.cds/POLICY.md` or equivalent, not in CDS.md.

**The collapse signal is observable.** If a cycle's receipt cannot
distinguish α and β authorship (single signature, single session, no
review record), the cycle is structurally invalid for substantive
software work. The project binding (or δ as fallback) must reject the
receipt and re-run the cycle with separated actors. For
β-α-collapse-on-δ skill/docs-class cycles, the receipt explicitly
records the collapse in `gamma-scaffold.md §"Dispatch shape"` and the
configuration-floor cap in `gamma-closeout.md §"Configuration-floor
declaration"`; the collapse is acknowledged, not hidden.

The Sub-3-vs-Field-6 line: Field 6 owns the actor-shape taxonomy
(prohibited vs permitted vs conditional); Sub 3 will detail the
dispatch model — the sequential bounded dispatch mechanism, the
re-dispatch prompt forms (initial / patch / full-review), the initial
dispatch sizing rules, the commit-checkpoint discipline. Field 6 owns
the actor-shape taxonomy; Sub 3 owns the per-shape operational detail.

---

## Selection function

CDS selection is coherence-driven. The next substantial engineering gap is
selected by the following function, in order. Each rule names the
circumstance that fires it, the action it requires, and the relationship to
the surrounding rules; collectively they discriminate which candidate cycle
δ dispatches next. The function is the canonical content for **what δ
selects** (paired with §Field 4 "δ cadence", which names **when δ runs**).

The 10 rules below are the canonical statements. Operational realization —
the per-rule mechanics γ runs at observation time (read CHANGELOG TSC table,
read encoding lag table, read doctor/status, read last PRA, build the
candidate table, name the decisive clause) — is the v0.1 overlay in
[`cnos.cdd/skills/cdd/gamma/SKILL.md`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
§2.1–§2.2 until the v1 role rewrite (Subs 3+ of a post-#403 wave) lands a
CDS-side `gamma/SKILL.md`. The CDS-side `selection/SKILL.md` thin overlay
(if present) delegates mechanics to the cdd `gamma/SKILL.md` surface.

### Inputs to selection

The selection function reads four observation surfaces (per
[`cnos.cdd/skills/cdd/gamma/SKILL.md §2.1`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md);
folded here from extraction-map row 1's §Inputs cross-cut):

- **Last post-release assessment (PRA)** — read first. The prior cycle's
  next-MCA commitment, deferred outputs, cycle-iteration findings, MCI
  freeze state. Binding, not optional.
- **CHANGELOG TSC table** — the technical-state-of-the-codebase signal
  recorded at each release.
- **Encoding lag table** — the cross-cycle backlog of gaps that have
  been named but not closed.
- **Doctor / status / operational-health surface** — the substrate-level
  health probe (build state, CI state, deploy state, runtime contracts).
- **Cross-repo proposal intake** — scan `.cdd/iterations/cross-repo/{target}/*/STATUS`
  for `submitted` (or `drafted` with source-delegated filing-authority per
  `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.3.3`) proposals; each
  candidate is a potential gap.

Operational mechanics (the `gh`/`MCP`/`git` queries that read each surface,
the candidate-table format) stay in the cdd `gamma/SKILL.md` v0.1 overlay.

### P0 override

If doctor/status shows a P0 bug — crash, data loss, silent failure, or
analogous safety-class incoherence — that is the gap. No further selection
logic applies until it is addressed.

**New-vs-known rule.** If the P0 was already visible when the last
assessment was written and the assessment committed a different next MCA,
the assessment decision governs unless the P0 has escalated (e.g. is now
causing active data loss or blocking all development). A known P0 that was
weighed and deferred is not an override — it is a prioritized backlog item.

### Operational-infrastructure override

If core operational paths are broken, fix them before feature work.
Examples: self-update broken; logging absent; health checks missing;
deployment path incoherent; system cannot observe or maintain itself. These
are not "nice to have" — they are preconditions for coherent engineering
work; subsequent cycles depend on the operational substrate to be running.

**Sizing rule.** Before selecting infrastructure debt as a full cycle, ask
whether the fix is cycle-sized or immediate-output-sized. If the fix is
executable in minutes (a script, a hook, a one-line config change), execute
it as an immediate output (per §Closure once Sub 5 lands; or per the
operational v0.1 overlay until then) and continue to the next selection
rule. Only select infrastructure debt as the cycle gap when the fix requires
design, multiple files, tests, or review.

### Assessment-commitment default

If the last PRA named a next MCA and no stronger override fires, that MCA is
selected by default. Observation may override it, but the override must be
stated explicitly — the prior assessment was authored under more recent
context than this cycle's observation, and silently abandoning a committed
next-MCA is incoherence between cycles.

### Stale-backlog re-evaluation

If P0 / operational / assessment-commitment produce no actionable candidate
(P0s exist but have no clear fix path, the PRA doesn't commit a next MCA, no
operational debt), re-evaluate stale issues before selecting new work:

- For each stale issue: is it still a real gap, or has the system moved
  past it?
- **Descope** issues that are no longer coherence gaps (close with
  rationale).
- **Consolidate** issues that overlap or could be addressed by one MCA.
- **Commit** the stale issue with the clearest fix path as the next MCA.
- If no stale issue has a clear fix path either, enter observation mode
  (no-gap case, below).

Stale backlog accumulating across multiple cycles without re-evaluation is
itself an incoherence.

### MCI freeze check

If the lag table contains stale issues (Master Coherence Issues that have
been open across multiple cycles), the next substantial MCA must come from
the stale set. New design work is **frozen** until at least one stale item
ships. The freeze prevents the lag table from monotonically growing while
δ-the-operator's attention is captured by newer-feeling work.

### Weakest-axis rule

If no stronger rule decides selection, choose work that addresses the
weakest current axis:

- **α** → structural / consistency work (the artifact-production discipline
  has a gap a future cycle's α will inherit).
- **β** → alignment / integration work (the review-oracle discipline has a
  gap that lets matter through without binding on a contract).
- **γ** → process / evolution work (the cycle-coordination discipline has a
  gap that produces friction across cycles).

The axis health signal is the receipt-stream view ε reads (§Field 5);
weakest-axis is what longitudinal observation discriminates.

### Maximum-leverage rule

Among candidates that address the weakest axis, choose the one that moves
the most lag entries. "Leverage" here is the count of downstream gaps a
candidate's closure unblocks; high-leverage candidates produce coherence
delta across multiple lag rows, low-leverage candidates close one row.

### Dependency order

If candidate A blocks candidate B blocks candidate C, choose A regardless
of local excitement or presentation value. Dependency-order is the lattice
the candidate table is sorted against; the topological prefix is
non-negotiable.

### Effort-adjusted tie-break

Between candidates with equal leverage and axis effect, choose the smaller
one. Ship sooner, observe sooner, correct sooner — the engineering loss
function (artifact improvement under repairable feedback) rewards shorter
cycles when other axes are tied.

### No-gap case

If:

- no P0 exists;
- no operational-infrastructure override exists;
- no stale lag item exists;
- no prior assessment commitment forces a next MCA;
- axes are healthy or tied;

then **do not start a new substantial cycle.** Remain in observation, or
choose a small-change path (immediate output executed inline; no new cycle
opened). The no-gap case is the structural permission to *not* dispatch —
the function does not always produce a candidate, and forcing one when no
coherence gap exists is itself incoherence (cycle-for-cycle's-sake).

### Operational realization

The 10 rules above are CDS's canonical selection-function statement. The
v0.1 operational overlay — how γ reads each input surface, builds the
candidate table, applies the rules in order, names the decisive clause,
sizes the intervention, and dispatches — lives in the existing cdd role
skills as the temporary v0.1 overlay until the v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.1`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — Step 1a: observe and build candidates (the four observation surfaces,
  the candidate-table format).
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.2`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — Step 1b: select by rule order, name the decisive clause, size the
  intervention (immediate-output / small-change / substantial-cycle).
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.2a`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — peer enumeration at scaffold time (the §Gap honest-claim discipline
  that gates a selected candidate before issue authoring).
- [`cnos.cdd/skills/cdd/alpha/SKILL.md §2.1`](../../../cnos.cdd/skills/cdd/alpha/SKILL.md)
  — dispatch-intake: how α reads the selected candidate's issue + active
  skills, the receipt α gives the selection function.
- [`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md)
  — the cross-repo proposal STATUS state machine that feeds the
  cross-repo-intake input surface.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{gamma,alpha,…}/SKILL.md`; the
v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Development lifecycle

CDS owns the full arc from observation back to observation. The lifecycle
discriminates **when** each role acts on a cycle (paired with §Field 4
"δ cadence", which names **what triggers** each cycle, and §Selection
function above, which names **what** δ selects). The lifecycle is the
canonical content for the engineering-cycle structure.

This section names the canonical 0–13 step list, the S0–S12 state machine
that the steps compose into under sequential bounded dispatch, the branch
rule (`cycle/{N}` canonical), the γ-owned branch pre-flight, and the skill
loading tier structure. Operational realization stays in the cdd role/runtime
skills as the v0.1 overlay until the v1 CDS-side role rewrite.

### Step table

A complete CDS cycle composes 14 ordered steps. Each step has a single
owner (the role that produces the step's required output), an inputs set
(the prior steps' outputs the step reads), and a required output (the
artifact or state change the step lands before the next step fires).

| # | Step | Owner | Purpose | Required output |
|---|------|-------|---------|-----------------|
| 0 | Observe | γ | Read current coherence state | Selection inputs read |
| 1 | Select | γ | Choose the next gap by §Selection function | Named selected gap + decisive rule clause |
| 2 | Branch | γ | Create the cycle branch (`cycle/{N}` from `origin/main`, §Branch rule + §Branch pre-flight) | `origin/cycle/{N}` exists |
| 3 | Bootstrap | γ | Author `gamma-scaffold.md` + issue pack | Issue #N with ACs + scaffold on `cycle/{N}` |
| 4 | Gap | α | Name the incoherence precisely | `self-coherence.md §Gap` (the coherence contract) |
| 5 | Mode | α | Choose MCA/MCI and active skills | Mode declaration + Tier 1/2/3 skill set |
| 6 | Artifacts | α | Produce in artifact order: design → contract → plan → tests → code → docs | Aligned implementation artifacts on `cycle/{N}` |
| 7 | Self-coherence | α | α audits own work against ACs and the triad | `self-coherence.md` complete through §CDD Trace |
| 8 | Review | β | β review CLP — verdict A / RC / NO-GO | `beta-review.md` round verdict; merge on A |
| 9 | Gate | δ | Verify release-readiness; record boundary decision | δ preflight verdict (Proceed / RC / Override) |
| 10 | Release | δ | Tag, publish, deploy (release-effector mechanics) | Release artifacts exist; tag on main |
| 11 | Observe | γ | Confirm runtime matches design | Observation result (CI green; runtime probe) |
| 12 | Assess | γ | Author the PRA (post-release assessment) | `POST-RELEASE-ASSESSMENT.md` |
| 13 | Close | γ | Execute immediate outputs; commit deferred outputs; declare closure | `gamma-closeout.md`; closure declaration |

Step 13 feeds back to Step 0 of the next cycle. The lifecycle is the loop
the CCNF kernel's recursion modes (`COHERENCE-CELL-NORMAL-FORM.md §Recursion
Modes`) realize on the engineering substrate; cross-scope accept advances
the scope index; within-scope repair-dispatch (a fix round in Step 8) keeps
the cell open at scope `n`.

### State machine

The states below are the canonical state machine for a substantial triadic
cycle under the **sequential bounded dispatch model** (the dispatch model
named in §Field 6 actor-collapse-rule and detailed in the v0.1 overlay at
[`cnos.cdd/skills/cdd/CDD.md §1.6`](../../../cnos.cdd/skills/cdd/CDD.md)).
Each state has one owner, required inputs, required outputs, a next state,
and a failure / retry path.

| State | Owner | Required inputs | Required outputs | Next state | Failure / retry |
|---|---|---|---|---|---|
| **S0: Observed** | γ | CHANGELOG TSC, lag table, doctor/status, last PRA | Named selected gap + decisive §Selection clause | S1 | Re-observe if no clear gap (no-gap case) |
| **S1: Issue filed** | γ | Selected gap; issue-quality gate passed | Issue #N with ACs, Tier 3 skills, non-goals, related artifacts | S2 | Fix issue until quality gate passes |
| **S2: Branch created** | γ | Issue #N open; `origin/cycle/{N}` does not exist | `origin/cycle/{N}` exists; pre-flight passed | S3 | Pre-flight fail → fix and retry |
| **S3: α dispatched** | δ | α prompt (from γ); `origin/cycle/{N}` exists | α session running | S4 | α session fails → re-dispatch α |
| **S4: α implementing** | α | Branch, issue, active skills loaded | Code, tests, docs on `cycle/{N}` | S5 | Blocker → γ unblocks; RC from β → re-dispatch α for fix round (within-scope repair) |
| **S5: α signaled review-ready** | α | Pre-review gate passed (`alpha/SKILL.md §2.6`) | `self-coherence.md` review-readiness section on `cycle/{N}` | S6 | Gate fails → fix and re-signal |
| **S6: β reviewing** | β | `self-coherence.md` review-readiness, diff, CI green | `beta-review.md` round verdict (RC or A) | RC → S4; A → S7 | β blocked → γ unblocks |
| **S7: β merged** | β | Approved verdict; pre-merge gate passed (`beta/SKILL.md §"Pre-merge gate"`) | Merge commit on main with `Closes #N`; `beta-closeout.md` | S8 | Merge conflict → β resolves in throwaway worktree |
| **S8: α close-out** | α (re-dispatched) | `beta-review.md` (approved); merged state | `alpha-closeout.md` on main | S9 | Re-dispatch unavailable → provisional close-out at review-readiness (declared as debt) |
| **S9: γ triaging** | γ | `alpha-closeout.md`, `beta-closeout.md`, `RELEASE.md` | PRA at canonical path; γ close-out triage; skill patches landed | S10 | Missing close-out → request re-dispatch; missing RELEASE.md → γ writes it |
| **S10: δ release-boundary preflight** | δ | PRA present; RELEASE.md present; `.cdd/unreleased/{N}/` not yet moved; merge on main | Proceed / Request changes / Override | Proceed → S11; RC → route to β/α/override | δ blocks → γ routes change |
| **S11: γ closing** | γ | δ preflight passed; closure gate rows pass (`gamma/SKILL.md §2.10`) | `gamma-closeout.md`; cycle-dir move to `.cdd/releases/{X.Y.Z}/{N}/`; closure declaration | S12 | Missing artifact → obtain before closing |
| **S12: δ disconnect** | δ | γ closure declaration; `gamma-closeout.md` exists on main | Tag + release commit; branch cleanup | S0 (next cycle) | Script fails → fix and retry |

The state machine maps onto the 0–13 step table as follows: Steps 0–3
compose S0–S3 (γ-owned setup); Steps 4–7 compose S4–S5 (α implementation
and review-ready signal); Step 8 composes S6–S7 (β review and merge); Step
9–10 compose S8–S10 (close-out collection and δ preflight); Steps 11–12
compose S9 + S11 (γ triage, PRA, closure declaration); Step 10's release
mechanics compose S12 (δ disconnect). The step table is the canonical
ordering; the state machine is the canonical ownership-and-transition shape.

### Branch rule

Every substantial CDS cycle is developed on its own dedicated branch. No
substantial CDS work is performed directly on `main`.

**Canonical branch format:**

```text
cycle/{N}
```

where `{N}` is the issue number. One cycle = one branch = one named target
for all polling. The cycle branch is the canonical coordination surface
during the cycle (§"Coordination surfaces" pending Sub 4; currently sourced
from `cnos.cdd/skills/cdd/CDD.md §Tracking` v0.1 overlay).

**Ownership.** γ creates `cycle/{N}` from `origin/main` **before** dispatch
(Step 2 of §"Step table"). α and β **never** create branches — they `git
switch cycle/{N}` from the name in their dispatch prompt and refuse to
invent or accept any other name. This is the role-side enforcement of the
single-named-target rule.

**γ session branch.** γ's session is on a separate γ session branch
(harness-given as `claude/...` or operator-named as `gamma/session-{N}`);
this branch is γ's pre-publication drafting surface for work that cannot
land directly on `cycle/{N}`. **Default rule:** γ commits cycle artifacts
directly to `cycle/{N}`. The session branch is for genuine drafts only
(e.g. PRA work-in-progress that requires γ-only review before publication).
Any γ-session-branch commit must merge into `cycle/{N}` or `main` by closure
declaration time, or be discarded. No orphan γ session branches survive
past closure.

**Legacy shapes — warn-only.** Pre-#287 shapes (`{agent}/{version}-{issue}-{scope}`,
`{agent}/{issue}-{scope}`, `claude/{slug}-{rand}`) are warn-only — verifiers
may flag them; nothing else changes. Frozen historical branches under those
shapes are not retroactively renamed. Polling glob `'origin/claude/*'` is
retained only for retrospective tracking and is not a discovery surface for
new cycles.

### Branch pre-flight

γ runs branch pre-flight **before** creating `cycle/{N}` (Step 2 of the
step table). Verify:

- `origin/cycle/{N}` does not yet exist (`git rev-parse --verify
  origin/cycle/{N}` returns non-zero — fail loud if the branch already
  exists, since one cycle = one branch);
- no stalled `.cdd/unreleased/{N}/` directory exists on `origin/main`
  (would indicate a previous cycle for `{N}` did not complete its
  release-time move per the §Artifacts location-matrix rule, pending
  Sub 4 — currently sourced from
  [`cnos.cdd/skills/cdd/release/SKILL.md §2.5a`](../../../cnos.cdd/skills/cdd/release/SKILL.md));
- the issue's intended scope is declared in the issue body before the
  branch is created;
- base SHA is known (`git rev-parse origin/main`);
- the issue is open.

The pre-flight is γ-owned because γ creates the branch. α's pre-review
gate retains its own row (verify that `origin/cycle/{N}` is rebased onto
current `origin/main` at review-readiness time —
[`cnos.cdd/skills/cdd/alpha/SKILL.md §2.6 row 1`](../../../cnos.cdd/skills/cdd/alpha/SKILL.md)),
but α does not create the branch and does not own pre-flight at branch
creation.

### Skill loading tiers

Skills are loaded in tiers. For substantial cycles, all applicable tiers
are mandatory. The tier structure is **executable**: every role can answer
"did I load what this tier requires?" without ambiguity.

**Tier 1a — CDS authority (always loaded by every role).**

- This file — `CDS.md` (the canonical instantiation contract).
- The kernel — [`cnos.cdd/skills/cdd/CDD.md`](../../../cnos.cdd/skills/cdd/CDD.md)
  (the CCNF spine; loaded as cited reference, not re-restated).
- The loader entrypoint — [`SKILL.md`](SKILL.md).
- The current role skill (when v1 lands): `alpha/SKILL.md`, `beta/SKILL.md`,
  or `gamma/SKILL.md` under `cnos.cds/skills/cds/`. Until v1 lands, the
  cdd-side role skills serve as the v0.1 overlay
  ([`cnos.cdd/skills/cdd/{alpha,beta,gamma}/SKILL.md`](../../../cnos.cdd/skills/cdd/)).

Tier 1a is loaded unconditionally. Roles do **not** load peer role skills.

**Tier 1b — Lifecycle phase skills (load per current phase or work shape).**

The lifecycle sub-skills under `cnos.cdd/skills/cdd/` (until Sub 5 of #403
migrates them to `cnos.cds`):

- `issue/SKILL.md` — when interpreting issue ACs or quality;
- `design/SKILL.md` — when producing or judging design-required work;
- `plan/SKILL.md` — when implementation sequencing is non-trivial;
- `review/SKILL.md` — when reviewing;
- `release/SKILL.md` — when merging / tagging / deploying;
- `post-release/SKILL.md` — when assessing or closing.

Load by phase. Roles that operate across multiple phases load the matching
set.

**Tier 1c — β closure bundle (always loaded by β).**

β always loads:

- `review/SKILL.md`;
- `release/SKILL.md`.

γ always loads `post-release/SKILL.md` because γ owns the PRA (see §Field 5
ε iteration cadence, §Field 3 close-out artifact set).

**Tier 2 — General engineering (load the applicable bundle).**

Pick the engineering-substrate bundle for the work shape from
[`src/packages/cnos.eng/skills/eng/README.md`](../../../../cnos.eng/skills/eng/README.md):
coding / review / design / runtime-and-platform / tooling / writing.

The engineering package README is the source of truth for which bundle
covers which work shape and which skills it includes. CDS does not
enumerate language- or platform-specific bundles here — that surface lives
in the engineering package and changes independently of the CDS protocol.

The skill that owns skill-program / frontmatter coherence is
[`cnos.core/skills/skill/SKILL.md`](../../../../cnos.core/skills/skill/SKILL.md)
(load when authoring or modifying skills). The skill that owns
architecture / design reasoning is
[`cnos.core/skills/design/SKILL.md`](../../../../cnos.core/skills/design/SKILL.md)
(load when reviewing or producing architecture-level decisions).

**Tier 3 — Issue-specific (named per issue).**

Skills that depend on what the cycle's work touches. The issue's
`## Skills to load` (or equivalent) section names these. Examples: a
language skill (`eng/{language}`); a domain skill (`eng/ux-cli`,
`eng/performance-reliability`, `eng/tool`); an architecture skill
(`eng/evolve`, `eng/write-functional`); a skill-authoring skill
(`cnos.core/skills/skill`).

γ names only Tier 3 in issues. Tier 1 and Tier 2 are mandatory and not
repeated in the issue body.

**Read each `SKILL.md` file before beginning any work step.** Naming a
skill without reading it is not loading it. Loaded skills are **hard
generation constraints** — not post-hoc review checklists.

When in doubt about mode (MCA vs MCI), apply CAP: if the answer is already
in the system, cite it (MCA) — don't reinvent it (MCI). If two paths close
the same gap, take the lighter one unless the heavier one buys durability
the lighter one cannot.

Review (Step 8) checks whether the implementation is consistent with all
loaded tiers. Findings that a loaded skill would have prevented are
process debt (an §Field 5 ε-class `cds-skill-gap` signal).

### Operational realization

The step table, state machine, branch rule, pre-flight, and tier structure
above are CDS's canonical lifecycle statement. The v0.1 operational overlay
— how each role executes each step's mechanics, how the dispatch model
wires roles together across the state machine, how the branch pre-flight
runs in shell, how each tier is loaded by the harness — lives in the
existing cdd role/runtime skills as the temporary v0.1 overlay until the
v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/gamma/SKILL.md`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ algorithm across Steps 0–3 + 11–13 (observation, selection, issue
  pack, branch creation, scaffold, dispatch coordination, close-out
  triage, PRA, cycle iteration, closure declaration).
- [`cnos.cdd/skills/cdd/alpha/SKILL.md`](../../../cnos.cdd/skills/cdd/alpha/SKILL.md)
  — α algorithm across Steps 4–7 (intake, artifact-order production,
  self-coherence, pre-review gate, review-ready signal) + Step 8's α-side
  fix-round mechanics + Step 9's α close-out re-dispatch.
- [`cnos.cdd/skills/cdd/beta/SKILL.md`](../../../cnos.cdd/skills/cdd/beta/SKILL.md)
  — β role contract across Step 8 (review CLP, pre-merge gate, merge,
  β close-out).
- [`cnos.cdd/skills/cdd/delta/SKILL.md`](../../../cnos.cdd/skills/cdd/delta/SKILL.md)
  — δ role contract across Steps 9 + 12 (release-boundary preflight,
  disconnect release).
- [`cnos.cdd/skills/cdd/harness/SKILL.md`](../../../cnos.cdd/skills/cdd/harness/SKILL.md)
  — runtime-substrate mechanics (dispatch, polling, branch creation
  shell, session lifecycle, identity wiring).
- [`cnos.cdd/skills/cdd/release-effector/SKILL.md`](../../../cnos.cdd/skills/cdd/release-effector/SKILL.md)
  — tag / release / deploy mechanics δ invokes at Step 10 / S12.
- [`cnos.cdd/skills/cdd/operator/SKILL.md`](../../../cnos.cdd/skills/cdd/operator/SKILL.md)
  — δ-the-operator's session-routing and re-dispatch contract (the
  identity-rotation primitive at S3 / S6 / S8).

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{gamma,alpha,beta,delta,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Coordination surfaces

CDS coordination is artifact-driven on repo-native surfaces. The same three
surfaces that the engineering substrate already carries — GitHub Issues,
git branches, and per-cycle directories on the cycle branch — serve as the
coordination surfaces for α, β, γ, and δ. GitHub Pull Requests are **not**
part of the CDS coordination protocol: every action a PR records can be
expressed as a write to a file in the per-cycle directory plus a branch
commit; issues remain (gap-naming); branches remain (isolation); PRs add
ceremony without value for the artifact-driven model. This section names
the four canonical sub-surfaces — cycle-state evidence, polling primitives,
mid-flight clarification, and cross-repo proposals — and what each one
binds.

The section is paired with §Artifact contract below: §Coordination surfaces
names *where the cycle state is observed and emitted*; §Artifact contract
names *what the cycle artifacts contain and where they freeze*. The two
sections together specify the in-cycle coordination model that γ, α, and β
operate inside.

Operational realization stays in the cdd runtime/role skills as the v0.1
overlay until the v1 CDS-side role rewrite.

### Cycle-state evidence

Three observable surfaces carry the cycle state during in-version work:

- **Issue activity.** The cycle's GitHub issue is the gap-naming surface
  and the operator-facing audit log. γ files; α/β/γ subscribe. State
  observed: comments (operator-class clarifications, β verdicts when
  posted as comments), label changes, issue body edits (cache-busted via
  §Mid-flight clarification below).
- **Cycle branch state.** The canonical cycle branch is `origin/cycle/{N}`
  per §Development lifecycle → §Branch rule. The branch's head SHA and
  the per-file blob SHAs under the cycle directory carry the in-cycle
  state — α's commits (review-ready signal, fix-round appendices), β's
  commits (review verdicts, merge commit), γ's commits (scaffold,
  clarifications, close-out triage). All role artifacts live on the
  cycle branch; `main` is the merge target only, never the in-cycle
  coordination surface.
- **Cycle directory state.** The per-cycle directory at
  `.cdd/unreleased/{N}/` (or `.cds/unreleased/{N}/` post-re-rooting; see
  the re-rooting note below) carries the role-distinguished-by-filename
  artifact set: `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`,
  `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md` (or
  `cds-iteration.md` post-rename), and any cycle-local additions
  (`gamma-scaffold.md`, `gamma-clarification.md`, etc.). File additions
  and blob-SHA transitions are the in-cycle progress signal.

**The three surfaces are observed transitionally.** Roles read each
surface synchronously at session start (the baseline pull) and then poll
for transitions (the wake-up loop) — see §Polling primitives below for the
mechanics that prevent the polling-looks-healthy-but-truth-is-stale
failure mode.

**Cycle directory naming — planned re-rooting (documented, not performed).**
The current project binding at `usurobor/cnos` carries per-cycle
directories under `.cdd/unreleased/{N}/`. The destination naming under the
post-#403-wave project binding is `<project>/.cds/unreleased/{N}/` —
software-class cycles belong under `.cds/`, not under `.cdd/`. The
filesystem rename is **out of scope for this section's authoring
cycle**; it is itself a separate post-#403 coordination problem (it must
coordinate with all in-flight cycle directories on `origin/main`, with
the `cn cdd verify` schema paths, and with the historical
`.cdd/releases/{X.Y.Z}/{N}/` directories that hold post-release moved
artifacts). Until the rename lands, every path in this section and in
§Artifact contract below uses the current `.cdd/` form and notes the
destination form once. Subsequent paths in the same context read
`.cdd/`-form-now / `.cds/`-form-after.

### Polling primitives

Polling has three parts, all mandatory: (a) the **query** that detects new
state, (b) the **wake-up mechanism** that returns control to the role's
session on transition, and (c) **reachability verification** that the
chosen query form works in the current environment. Polling without a
wake-up is silent — the loop runs but the role never reacts. Polling with
an unreachable query is silent in the same way — the loop runs, returns
empty or errors, and the role assumes nothing has happened.

**Query forms — pick whichever is reachable in the environment.**
Different environments expose different surfaces; the query form is
discovered at session start by reachability preflight (below).

| Surface | `gh` form (shell envs) | MCP form (MCP-only envs) | git form (clone-aware envs) |
|---|---|---|---|
| Issue comments | `gh issue view {N} --comments` | `mcp__github__issue_read` method=`get_comments` | — |
| Cycle branch existence (γ pre-flight) | — | — | `git fetch --quiet origin cycle/{N} && git rev-parse --verify origin/cycle/{N}` |
| Cycle branch head SHA | — | — | `git fetch --quiet origin cycle/{N} && git rev-parse origin/cycle/{N}` (compare to prior, emit on change) |
| Cycle directory state | — | — | `git fetch --quiet origin cycle/{N} && git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` (compare blob SHAs; emit on add or change) |

γ owns the tight loop on a single named branch (per §Development lifecycle
→ §Branch rule, one cycle = one branch = one named target); polling
`origin/main` for in-flight cycle dirs is silent — the files live on
`origin/cycle/{N}` until the release-time move per §Artifact contract →
§Location matrix.

**Wake-up mechanism — name it explicitly in the role's session.** Polling
is only effective if the loop's transition produces a notification that
wakes the role. Each environment has its own form:

- **`Monitor` (Claude Code on the web):** wrap polling in a transition-loop
  whose stdout lines deliver as `task-notification` system messages.
  Transition-only emission is mandatory (emitting on every iteration
  floods the session's context budget).
- **Shell wake hook (custom harness):** the loop's exit signals the
  session. Verify the harness contract before relying on it.

If neither a `Monitor`-equivalent nor a shell-wake harness exists, the
role cannot autonomously detect cycle progression — the gap is surfaced
to the operator before dispatch.

**Reachability preflight — run before committing to a query form.** Query
forms are not interchangeable: `gh` requires shell + `gh auth`; MCP tools
cannot run inside `Monitor` or background shell loops; `git fetch`
requires network access to the remote from the execution environment;
direct `api.github.com` access may be blocked by sandbox network policy.
Before starting a polling loop, the role probes the chosen query form
once synchronously and confirms it returns real data. If it fails, the
role falls back to the next available form. If no form is reachable, the
role surfaces the gap to the operator before proceeding — it does not
silently assume polling is working.

**Transition-only emission is mandatory.** A loop that emits on every
poll fills the session with `task-notification` blocks and consumes
context budget. The loop computes `cur != prev`, emits only on
transition, then updates `prev := cur` and sleeps.

**Synchronous baseline pull is a precondition of transition-only polling.**
Transition-only emission is correct on its own terms (avoid context
flood) but has a structural blind spot: the loop's first iteration sets
`prev` to the empty string and silently absorbs whatever already exists.
State that exists *before* the polling loop's first iteration will never
surface as an event. Every transition-only polling loop must therefore be
paired with a synchronous initial-state pull of the same surface
immediately when the role's session starts — the synchronous channel
owns the past, the polling channel owns the future. For the `cycle/{N}`
model, the baseline pull is per-cycle: `git rev-parse --verify
origin/cycle/{N}` to confirm the branch exists, `git rev-parse
origin/cycle/{N}` for the head SHA, `git ls-tree -r origin/cycle/{N}
.cdd/unreleased/{N}/` for the cycle artifact set, and the gh/MCP form
for issue activity. Reading `.cdd/unreleased/` from `origin/main` only
surfaces cycles that have already merged; in-flight cycles live on
`origin/cycle/{N}` and will be invisible to a `main`-only baseline.

**`git fetch` reliability is an explicit dependency.** The `git fetch
--quiet` form silently swallows transport flake (DNS hiccup, proxy 502,
expired token, transient 4xx) — fetch returns 0, the local ref does not
advance, and the per-iteration comparison sees `cur == prev` and emits
nothing. The transition loop then drops every commit landing during the
flake window. **Mitigation:** after **N successive empty iterations**
(canonical N = 10, ≈ 10 minutes at 60s interval), the loop does a
synchronous reachability re-probe with explicit stderr capture (`git
fetch --verbose ...` with stderr to a log). If the re-probe succeeds,
the transition loop continues. If it fails, the role surfaces the failure
to the operator immediately — silently looping with a broken transport
is the failure mode the rule catches.

**Single named branch — no globs for new cycles.** Polling targets
`origin/cycle/{N}` directly. There is no glob discovery step for new
cycles — γ tells α and β the branch name in the dispatch prompt. The
legacy glob `'origin/claude/*'` is retained only for retrospective
tracking of historical cycles whose branches predate the `cycle/{N}`
rule (§Development lifecycle → §Branch rule "Legacy shapes"); it must
not be used as a discovery surface for new cycles.

### Mid-flight clarification

When γ edits the cycle's issue body mid-cycle (because α surfaced an
ambiguity, because δ refined an axis, because operator clarified a
non-goal), γ writes a `gamma-clarification.md` entry to
`.cdd/unreleased/{N}/` on the cycle branch *before* signaling the edit.
The entry names: date, edit summary, and which ACs / non-goals /
constraints / artifacts changed.

**The cycle-branch SHA transition is the cache-bust signal.** Roles
polling the cycle branch (per §Polling primitives) see the SHA advance;
their next intake re-fetches the issue body from the live source rather
than trusting cached state. The signal is the cycle-branch SHA transition,
not the GitHub issue mtime — which may be invisible to MCP-cached
issue reads.

**Empirical anchor — cnos#391.** Wrong-shape implementation (wrong
package scoping + a separate-binary axis γ did not pin at dispatch);
γ recovered by editing the issue body to pin the missing axes and
committing `gamma-clarification.md` to the cycle branch; α picked up
the cache-bust via cycle-branch polling and re-shaped the implementation.
The cycle was rescued mid-flight without abandoning the branch. The
empirical-anchor doctrine cited at §Field 4 (δ cadence) — δ-as-architect
per cnos#393 — names this rescue path as the structural escape valve
when γ under-specifies and α improvises.

**γ MUST NOT silently re-pin a contract axis without logging in
`gamma-clarification.md`.** Mid-cycle contract changes are recorded in
the clarification artifact; the file is part of the cycle directory and
travels with the cycle through release.

**Wire-format canonical home.** The wire-format invariant for the
mid-flight rescue mechanism (file path, authoring role, reader role,
trigger conditions, cache-bust function, resumption protocol,
spec-staleness propagation) lives at
[`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../../../cnos.handoff/skills/handoff/mid-flight/SKILL.md),
landed under Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404)
(cnos#418). This section retains the v0.1 software-cycle operational
realization (the `gh issue edit` + `git commit` + `git push` shell
sequence γ runs at clarification time; the polling-primitive mapping);
the wire-format invariants are owned at cnos.handoff and cited from here.

### Cross-repo proposals

Cycles that originate work in one repo and land it in another use the
cross-repo proposal lifecycle. The proposal carries an event log
(`STATUS`) whose vocabulary tracks the proposal across the
source-target handshake. The same vocabulary applies to inbound
proposals (an upstream agent hub proposing work to cnos), outbound
proposals (cnos proposing work to a downstream repo), and bilateral
iterations.

**Eight events constitute the canonical vocabulary:**

| Event | Meaning |
|---|---|
| `drafted` | Source has written the proposal but has not requested target action. Pre-intake. |
| `submitted` | Source requests target intake. This is the only event required for target intake. |
| `accepted` | Target γ will act substantially as proposed and has filed a target reference. |
| `modified` | Target γ accepts the governing gap but changes scope, split, wording, implementation, proof, or patch application materially. Carries a `Delta` field in the target issue's `## Source Proposal` block. May fire post-`accepted` to record a refinement. |
| `landed` | Target work merged or otherwise became target truth. For 1:1 proposals: one event total. For master/sub: one per sub merge + one terminal master-close event. |
| `rejected` | Target γ declines the proposal. Terminal. |
| `withdrawn` | Source retracts the request. Terminal. |
| `revised` / `corrected` | Optional audit events for post-submission revisions or corrections. Append; do not rewrite history. May fire from any non-terminal state without changing lifecycle state. |

**The STATUS state machine** (canonical transition graph):

- `(start) → drafted | submitted` (source authors directly into either state)
- `drafted → submitted` (source γ requests intake)
- `drafted → accepted | modified | rejected` (permitted when source
  explicitly delegates filing-authority to target without intermediate
  `submitted` — the agent-hub direct-acceptance path)
- `drafted → withdrawn`
- `submitted → accepted | modified | rejected | withdrawn`
- `accepted → modified` (post-filing refinement; Delta updated)
- `modified → modified` (further refinement)
- `accepted → landed` and `modified → landed` (target work merges)
- `* → revised | corrected` (audit-only; lifecycle state unchanged)

**Illegal transitions:** `rejected → *` (terminal); `landed → *` (terminal
for 1:1; master/sub permits multiple `landed` rows per the master/sub
rule); `withdrawn → *` (terminal); `submitted → landed` (must pass
through `accepted` or `modified`); `modified → accepted` (cannot un-modify
once Delta is recorded).

**Bundle-state phases** (`open | converging | closed`) derive from
STATUS: `open` ↔ `drafted` or `submitted`; `converging` ↔ `accepted`
or `modified` without terminal `landed`; `closed` ↔ terminal `landed`
(1:1 or master-close) or `rejected` or `withdrawn`. Audit events
(`revised`, `corrected`) do not change the bundle phase.

**Master/sub rule for `landed`:** for master/sub-shaped proposals
(a parent issue with N sub-issues), one `landed` event fires per sub
merge plus one terminal master-close event when the master issue closes.
A wave with 4 subs across 3 releases produces 5 `landed` rows: 4 per-sub
+ 1 master-close. For 1:1 proposals: one `landed` event total; no
separate master-close.

**Filing-decision rule:** after the target γ's filing decision, source
STATUS MUST NOT remain at `submitted` (or at `drafted` once target γ
has taken action on a delegated direct-acceptance path). The decision
is recorded — directly or via `FEEDBACK.patch` — within the same target
session that made the decision.

**Bundle layout** (canonical source-side path):

```text
{source-repo}:.cdd/iterations/cross-repo/cnos/{slug}/
  ISSUE.md
  STATUS
  PATCH.diff        # optional
  LINEAGE.md        # source-side trace
```

The matching cnos-side mirror lives at
`cnos:.cdd/iterations/cross-repo/{source-repo}/{slug}/`. The path is
direction-agnostic — outbound iteration traces and bilateral iterations
live at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` on
whichever side carries the bundle.

**Operational realization location.** The canonical home for the
cross-repo proposal lifecycle is
[`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md),
landed by Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404)
(cnos#416). The skill was relocated from its v0.1 home at
`cnos.cdd/skills/cdd/cross-repo/SKILL.md` (where a compatibility pointer
remains for backward compatibility) because the cross-repo coordination
mechanism is wire-format / handoff-class — it coordinates proposals
across any c-d-X realization — and cnos.handoff is the canonical owner
of that surface class. CDS binds and consumes the STATUS state machine
declared there; it does not own the vocabulary.

### Operational realization

The three observation surfaces, the polling primitives, the mid-flight
clarification protocol, and the cross-repo proposal lifecycle above are
CDS's canonical coordination-surfaces statement. The v0.1 operational
overlay — the shell snippets that realize each polling query, the
Monitor-wrapped transition loops, the `gh issue edit` + `commit` +
`push` sequence that γ runs at mid-flight clarification time, the
per-event STATUS-write sequence that source and target γs run across
the proposal lifecycle — lives in the existing cdd role/runtime skills
as the temporary v0.1 overlay until the v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/harness/SKILL.md §5.4`](../../../cnos.cdd/skills/cdd/harness/SKILL.md)
  — Single-named-branch polling under `Monitor` (the transition loop
  with baseline sync + reachability re-probe at N=10 empty iterations).
- [`cnos.cdd/skills/cdd/harness/SKILL.md §5.1–§5.5`](../../../cnos.cdd/skills/cdd/harness/SKILL.md)
  — Polling and wake-up: δ's wake-up signals, issue activity polling,
  cycle branch polling, reachability probe.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.5`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ's coordination loop across the cycle branch: dispatch, polling
  cross-reference (citing `harness/SKILL.md §5.4`), the
  `gamma-clarification.md` issue-edit cache-bust procedure.
- [`cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.3`](../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md)
  — STATUS state machine: 8-event vocabulary (codified above), full
  transition graph (codified above), emitter-per-event rules, master/sub
  `landed` rule, bundle-state phase mapping, direct-acceptance
  (`drafted → accepted`) path. The skill's canonical home is
  `cnos.handoff/skills/handoff/cross-repo/` per Sub 2 of
  [cnos#404](https://github.com/usurobor/cnos/issues/404) (cnos#416); a
  compatibility pointer remains at `cnos.cdd/skills/cdd/cross-repo/` for
  backward compatibility.
- [`cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.1`](../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md)
  — Directional cases (1:1; master/sub; inbound; outbound) and the
  bundle file set per case.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{gamma,alpha,beta,harness,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Artifact contract

CDS is artifact-driven. Every substantial cycle must produce inspectable
artifacts at canonical paths with named owners and verification gates.
This section names the canonical artifact contract for CDS: the
terminology, the bootstrap rule, the ordered artifact flow across the
cycle, the per-step manifest, the canonical paths (the Location Matrix),
the role/artifact ownership matrix, the CDS Trace format, the supporting
rules, and the frozen-snapshot rule.

The section is paired with §Coordination surfaces above and with §Field 3
(γ close-out artifact). §Coordination surfaces names *where the cycle
state is observed*; §Artifact contract names *what each cycle artifact
contains, who owns it, where it lives, and when it freezes*. §Field 3
names the close-out artifact set as part of the six-field instantiation
contract; §Artifact contract here is the canonical operational realization
of that field.

Operational realization stays in the cdd role/runtime skills as the v0.1
overlay until the v1 CDS-side role rewrite.

**Wire-format canonical home.** The channel wire-format invariants
(channel-directory path pattern `.cdd/unreleased/{N}/`; per-role write
ownership pattern; sequential α→β→γ rule with γ-scaffold cycle-start
exception; frozen-snapshot rule on merge; release-time directory move)
live at
[`cnos.handoff/skills/handoff/artifact-channel/SKILL.md`](../../../cnos.handoff/skills/handoff/artifact-channel/SKILL.md),
landed under Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404)
(cnos#418). This section retains the CDS-specific per-artifact contract
(the Location matrix's specific filename set; the Ownership matrix's
per-role artifact assignment with verification gates; the 13-stage
Ordered flow); the channel-shape wire-format invariants are owned at
cnos.handoff and cited from here.

### Terminology

Four terms are used precisely across §Artifact contract and the
downstream sections (§Field 3, §Field 5, §Coordination surfaces):

- **Post-release** — the umbrella phase after release. Covers lifecycle
  Steps 11–13 of §Development lifecycle → §Step table (observe, assess,
  close).
- **Assessment** (a.k.a. post-release assessment, PRA) — the γ-owned
  repo artifact at the canonical path declared in §Location matrix
  below. One PRA per release.
- **Close-out** — a role-local findings record written by α, β, and γ at
  the canonical paths in §Location matrix. Close-outs feed γ's PRA and
  γ's close-out triage. They are not a substitute for the PRA.
- **Closure** — the final cycle state: PRA committed, all close-outs on
  main, immediate outputs executed, deferred outputs committed, hub
  memory updated. γ declares closure in `gamma-closeout.md`.

### Bootstrap

A CDS cycle bootstraps in two phases, both before α begins:

1. **γ scaffold + issue contract load.** γ authors `gamma-scaffold.md`
   on the cycle branch naming: the issue, the mode (substantial /
   small-change / immediate-output), the wave (if applicable per the
   wave-manifest pattern), the surfaces γ expects α to touch, the AC
   oracle approach, the empirical anchor (when framed), and the
   expected diff scope. The scaffold is the executable rendering of
   what γ already decided at selection time, not new analysis.
2. **Branch creation.** γ creates `cycle/{N}` from `origin/main` per
   §Development lifecycle → §Branch rule and runs the §Branch pre-flight
   gate before publishing the branch.

For substantial release-shipping cycles, α opens the cycle by creating a
**version directory** for the bundle that will receive the frozen
snapshot:

```text
docs/{tier}/{bundle}/{X.Y.Z}/
```

Each version directory contains a `README.md` snapshot manifest and one
stub per declared deliverable. Artifacts outside version directories
(e.g. files in `.cdd/unreleased/{N}/`, navigation updates) are not
required as bootstrap stubs. Small-change cycles may skip the version
directory; the exemption is recorded in the cycle's close-out per
§Field 6 (small-change collapse rule).

**Pre-dispatch γ scaffold check (binding gate).** γ MUST NOT proceed to
α dispatch until `gamma-scaffold.md` exists on `origin/cycle/{N}`. If
absent when γ is about to produce the α prompt, γ authors it first,
commits and pushes it to the cycle branch, then continues. The scaffold
makes the cycle's intent legible to α and to β before any matter is
produced.

### Ordered flow

A CDS cycle produces artifacts in a canonical 13-stage order. Each stage
has a single owner, a required output, and a stage-specific format spec
(per §Manifest below):

1. **design** — α authors the design artifact (when required by mode);
   names the invariant/volatile/boundary decomposition.
2. **contract** — α authors the coherence contract: gap, mode, active
   skills, ACs with oracle approach, known debt.
3. **plan** — α authors the implementation plan (when sequencing is
   non-trivial); names the order in which surfaces will be touched.
4. **tests** — α writes the tests that encode the ACs (or names "no
   tests required" with rationale).
5. **code** — α produces the implementation diff.
6. **docs** — α updates the doc surfaces affected by the implementation.
7. **self-coherence** — α audits the work against ACs and the triad;
   writes `self-coherence.md` to completion.
8. **review** — β runs the review CLP (terms / pointer / exit); writes
   `beta-review.md`; emits A / RC / NO-GO verdict per round.
9. **gate** — δ verifies release-readiness preconditions before tagging;
   records boundary decision (Proceed / Request changes / Override).
10. **release** — δ runs the release-effector mechanics (tag, build,
    deploy); release artifacts land on main.
11. **observe** — γ confirms runtime matches design; CI green; runtime
    probe passes.
12. **assess** — γ authors the PRA at the canonical version-directory
    path.
13. **close** — γ executes immediate outputs; commits deferred outputs;
    writes `gamma-closeout.md`; declares closure.

The stages compose with the 14-step (0–13) lifecycle table in
§Development lifecycle → §Step table: Steps 0–3 are pre-α γ work
(observe, select, branch, bootstrap); Steps 4–7 are α's stages 1–7 here;
Step 8 is β's stage 8; Steps 9–10 are δ's stages 9–10; Step 11 is γ's
stage 11; Step 12 is γ's stage 12; Step 13 is γ's stage 13.

### Manifest

For substantial changes, each artifact stage has a manifest row naming
the artifact, the role (α/β/γ/δ), the format spec, the canonical owner
location, and the required-or-conditional flag. The manifest is the
master reference for all stage attributes; the per-stage operational
detail (what to write, what skill to load, what the row looks like in
the file) lives in the cdd role skills as the v0.1 overlay.

The manifest below is the canonical CDS shape (transposed from the
pre-#402 CDD §5.3 with software-engineering vocabulary preserved):

| Stage | Phase | Role | Required output | Format spec | Required |
|---|---|---|---|---|---|
| design | build | α | design artifact OR explicit "not required" | `design/SKILL.md §3.1` | substantial only |
| contract | build | α | named incoherence + AC oracle | `self-coherence.md §Gap` + `§Mode` + `§ACs` | always |
| plan | build | α | sequencing artifact OR explicit "not required" | `docs/gamma/cdd/PLAN-TEMPLATE.md` | L7 / cycle-sized |
| tests | build | α | test files OR explicit reason none apply | diff | always |
| code | build | α | implementation diff OR "docs/process only" | diff | always |
| docs | build | α | changed canonical docs / specs / READMEs | diff | when docs affected |
| self-coherence | build | α | review-readiness signal complete | `self-coherence.md` carrying CDS Trace through stage 7 | substantial only |
| review | review | β | verdict + findings (round-by-round) | `review/SKILL.md` output format | always |
| gate | release | δ | release-readiness preflight verdict | `docs/gamma/cdd/GATE-TEMPLATE.md` | always |
| release | release | δ | tag + release notes + version-snapshot | `release-effector/SKILL.md` | always |
| observe | close | γ | post-release observation result | `post-release/SKILL.md` | always |
| assess | close | γ | `POST-RELEASE-ASSESSMENT.md` | `post-release/SKILL.md` output template | always |
| close | close | γ | immediate outputs + deferred committed + closure declaration | `gamma-closeout.md` | always |

**Manifest rules:**

- "Not required" is valid only when stated explicitly.
- An omitted stage with no explicit note is incomplete, not implicit.
- Small-change mode may collapse stages "design"–"self-coherence" into
  commit-message evidence; the same distinctions still apply.
- Skills loaded shape generation; record the skill that shaped each
  artifact in the §Trace format below.

### Location matrix

Every named artifact has exactly one canonical location. Verifiers
(e.g. `cn cdd-verify`, the forthcoming `cn cds-verify`) enforce these
paths as canonical and treat any other location as legacy/warn-only.

**Path notation.** Paths below use the current `.cdd/`-rooted form
(`.cdd/unreleased/{N}/`, `.cdd/releases/{X.Y.Z}/{N}/`). The destination
naming under the post-#403-wave project binding is the `.cds/`-rooted
form (`.cds/unreleased/{N}/`, `.cds/releases/{X.Y.Z}/{N}/`). The
**`.cdd/` → `.cds/` re-rooting is documented here as planned, not
performed** — the filesystem migration is its own separate post-#403
cycle (it must coordinate with all in-flight cycle directories, with
`cn cdd verify` schema paths, and with the historical
`.cdd/releases/{X.Y.Z}/{N}/` post-release directories). Until the
rename lands, every path in this section uses the `.cdd/` form.

| Artifact | Canonical repo location | CDS package default | Noncanonical / legacy / scratch |
|---|---|---|---|
| Version snapshot directory | `docs/{tier}/{bundle}/{X.Y.Z}/` | `docs/gamma/cdd/{X.Y.Z}/` | `.cdd/releases/{X.Y.Z}/` is **not** the frozen snapshot — it is triadic protocol/scratch space |
| POST-RELEASE-ASSESSMENT.md (PRA) | `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `.cdd/releases/{X.Y.Z}/beta/POST-RELEASE-ASSESSMENT.md` and `.cdd/releases/{X.Y.Z}/beta/ASSESSMENT.md` are legacy/warn-only |
| α self-coherence (primary branch artifact) | `.cdd/unreleased/{N}/self-coherence.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/self-coherence.md` at release | same | none — required for every substantial cycle |
| β review record | `.cdd/unreleased/{N}/beta-review.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/beta-review.md` at release | same | none — required for every substantial cycle |
| α close-out | `.cdd/unreleased/{N}/alpha-closeout.md` (in-version), moved at release | same | `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` (legacy aggregate form) is warn-only |
| β close-out | `.cdd/unreleased/{N}/beta-closeout.md` (in-version), moved at release | same | `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md` (legacy aggregate form) is warn-only |
| γ close-out | `.cdd/unreleased/{N}/gamma-closeout.md` (in-version), moved at release | same | `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md` (legacy aggregate form) is warn-only |
| cdd-iteration (per-cycle) | `.cdd/unreleased/{N}/cdd-iteration.md` (in-version), moved at release | same | none — required when `protocol_gap_count > 0`; courtesy stub permitted when count is 0 (per cnos#401 cadence rule cited in §Field 3 and §Field 5) |
| cdd-iteration aggregator | `.cdd/iterations/INDEX.md` (root, persistent) | same | one row per cycle that produced a `cdd-iteration.md`; γ updates at close |
| cross-repo trace | `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` | same | persistent until target PR merges; LINEAGE.md preserved in the target repo's `cdd-iteration.md` |
| γ kata verdict (optional) | `.cdd/releases/{X.Y.Z}/gamma/KATA-VERDICT.md` when kata is available | same | warn-only when kata unavailable |
| CHANGELOG ledger row | `CHANGELOG.md` (Release Coherence Ledger) — Version, C_Σ, α, β, γ, Level, coherence note | same | none |
| RELEASE.md | `RELEASE.md` at repo root, included in the release commit | same | CI auto-generated body is not an acceptable substitute |
| Hub memory | external agent-hub state (not a repo artifact) | external agent-hub state | the PRA records path / commit-sha / unavailable-reason; verifiers do not inspect the external hub directly |

**Path rules:**

- `.cdd/unreleased/{N}/` (destination: `.cds/unreleased/{N}/`) is the
  per-cycle coordination directory, keyed by issue number. Files inside
  are role-distinguished by **filename**, not by directory.
- `.cdd/releases/{X.Y.Z}/{N}/` (destination: `.cds/releases/{X.Y.Z}/{N}/`)
  holds the moved-at-release form of each cycle's coordination directory.
  Multiple cycle directories may live under one release directory when
  several issues ship in the same release.
- `.cdd/` (destination: `.cds/`) is triadic protocol space and role-local
  close-out evidence storage. It is **not** the canonical frozen
  post-release snapshot. The frozen snapshot lives under
  `docs/{tier}/{bundle}/{X.Y.Z}/`.
- Tags are bare `X.Y.Z` everywhere (VERSION file, git tag, branch-name
  version segment, CHANGELOG row, RELEASE.md, snapshot directory).
  `v`-prefixed tags are legacy and warn-only.
- All substantial cycles use the in-version cycle-directory
  artifact-exchange surface. GitHub Pull Requests are not used for
  CDS coordination (see §Coordination surfaces above).

### Ownership matrix

Every required cycle artifact has one owner, one verification gate, and
one consequence if missing. γ's closure gate (`gamma/SKILL.md §2.10` in
the v0.1 overlay) checks every row marked "Required before γ closure."

| Artifact | Owner | Written when | Verified by | Required before | Missing means |
|---|---|---|---|---|---|
| `self-coherence.md` | α | During α session, incrementally; review-readiness section last | β at review intake | β review | β waits; no review until review-readiness present |
| `beta-review.md` | β | During β review session, incrementally per round | γ at close-out triage | γ closure | γ cannot triage; requests β re-dispatch |
| `alpha-closeout.md` | α (re-dispatched after merge) | After β merge, via δ re-dispatch. **Provisional fallback:** α may write a provisional close-out at review-readiness (marked `[provisional]`); the pre-merge gate accepts this form. Full close-out via δ re-dispatch remains the normative path for tagged releases. | γ before PRA; γ closure gate; pre-merge gate (provisional accepted) | γ closure | γ closure gate blocks; γ requests δ to re-dispatch α |
| `beta-closeout.md` | β | Before β exits (same β session as merge) | γ before PRA; γ closure gate | γ closure | γ closure gate blocks; γ requests β re-dispatch |
| `gamma-closeout.md` | γ | After all closure gate rows pass | δ before tag/release (implicit: tag requires γ closure declaration, which is `gamma-closeout.md`) | δ tag/release | δ must not tag; γ has not declared closure |
| `cdd-iteration.md` | γ (with ε review) | Same session as `gamma-closeout.md`; required when `protocol_gap_count > 0`; courtesy stub permitted when count is 0 (per cnos#401) | γ closure gate; aggregator update (`.cdd/iterations/INDEX.md`) | γ closure | γ closure gate blocks; if `cdd-*-gap` findings exist, the artifact must exist before closure |
| `.cdd/iterations/INDEX.md` row | γ | Same session as `cdd-iteration.md` | γ closure gate | γ closure | γ closure gate blocks if `cdd-iteration.md` was written but INDEX.md was not updated |
| `RELEASE.md` | γ | Before requesting δ tag/release; committed to main in release commit | δ at release-boundary preflight | δ tag/release | δ must not tag; CI auto-generates sparse notes |
| `.cdd/releases/{X.Y.Z}/{N}/` (the cycle-dir move) | γ (at release) | Before γ requests δ tag; included in release commit | γ closure gate; δ preflight | δ tag/release | Stale unreleased dir; γ closure gate blocks until moved |
| POST-RELEASE-ASSESSMENT.md (PRA) | γ | After β merge + close-outs | γ closure gate; δ at release-boundary preflight | γ closure | γ closure gate blocks |

**Ownership rules:**

- `gamma-closeout.md` is the closure declaration artifact. δ's obligation
  to not tag before closure is satisfied when `gamma-closeout.md` exists
  and the closure-declaration commit is on main.
- For small-change cycles, `beta-closeout.md` and `gamma-closeout.md` may
  not apply per the small-change collapse rule (cited from §Field 6);
  `alpha-closeout.md`, `RELEASE.md`, and the cycle-directory move remain
  required.
- The cnos#401 cadence rule applies to the iteration artifact: required
  when `protocol_gap_count > 0`; courtesy stub permitted when count is 0
  for traceability (cited in §Field 3 and §Field 5).

### Trace format

Every substantial cycle must carry a lightweight execution trace named
the **CDS Trace** (renamed from the pre-extraction "CDD Trace"; the
format itself is verbatim — the rename reflects that the trace is
software-cycle-specific). The trace uses lifecycle step numbers from
§Development lifecycle → §Step table, not section numbers.

For stages design → release (0–10), the trace lives in the primary
branch artifact. For stages observe → close (11–13), closure lives in
the post-release assessment.

The **primary branch artifact** is the artifact that owns the named
incoherence, mode, active skills, and acceptance criteria. For triadic
cycles this is `.cdd/unreleased/{N}/self-coherence.md`. For
governance/process work, the governing doc being changed may carry the
trace inline. When a separate design artifact exists, the design
artifact carries the trace and `self-coherence.md` references it.

**Required format:**

```markdown
## CDS Trace
| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Observation inputs read; selected signal |
| 1 Select | — | — | Selected gap |
| 2 Branch | branch | cds | Branch created / verified against §Branch rule / §Branch pre-flight |
| 3 Bootstrap | version dir | cds | Bootstrap stubs created or explicit small-change exemption |
| 4 Gap | primary artifact | — | Named incoherence / coherence contract |
| 5 Mode | primary artifact | skill1, skill2 | Work shape, level (if used), mode, active skills |
| 6 Artifacts | design / plan / tests / docs | — | Artifact progress or explicit "not required" |
| 7 Self-coherence | `.cdd/unreleased/{N}/self-coherence.md` | cds | AC-by-AC self-check completed |
| 7a Pre-review | `.cdd/unreleased/{N}/self-coherence.md` | cds | Pre-review gate passed; review-readiness signaled |
| 8 Review | `.cdd/unreleased/{N}/beta-review.md` | review | CLP review result |
| 9 Gate | `.cdd/unreleased/{N}/beta-review.md` or release surface | release | Release-readiness decision |
| 10 Release | release surface | release | Tag / changelog / release decision |
```

**Trace rules:**

- One row per completed lifecycle step.
- The Step column carries both the number and the name for readability.
- "Skills loaded" is required when skills shaped generation or
  lifecycle execution.
- If a lifecycle skill is used later (review, release, writing,
  post-release), record it when it becomes active.
- Missing rows mean the stage is not yet evidenced.
- Contradictory rows are findings.

### Supporting rules

The following supporting rules govern the artifact contract; each is the
condensed CDS-side form of the pre-#402 CDD §5.5 supporting rules
(verbatim move with engineering-loss-function framing preserved):

- **One source of truth per fact.** A given fact lives in exactly one
  canonical surface; downstream surfaces cite, not duplicate.
- **Derive, do not duplicate.** Derived content (CHANGELOG rows from
  release commits, INDEX.md aggregator rows from per-cycle iteration
  artifacts) is generated, not hand-authored at the derivation site.
- **Update docs before release.** Documentation surfaces that describe
  current behaviour are updated in the same cycle as the behaviour
  change; documentation lag is a `cds-protocol-gap` per §Field 5.
- **Write tests before or alongside the code they validate.** The order
  is producer-discipline (α): the test surface and the code surface
  land in the same cycle, with tests authored before or alongside
  code, never after.
- **Build-sync source asset changes before commit.** Build-generated
  artifacts (compiled binaries, generated docs, schema-derived outputs)
  are regenerated from source before commit; the commit's state is
  internally consistent.
- **Enumerate affected files before implementation begins.** α names
  the surface set the cycle will touch in the contract (Stage 2) before
  producing matter; surfaces discovered mid-cycle are recorded in
  `gamma-clarification.md` per §Coordination surfaces → §Mid-flight
  clarification.
- **Every AC must map to evidence before review.** β's review intake
  requires per-AC oracle evidence to be present in
  `self-coherence.md §ACs`; ACs without oracles are RC findings.
- **All review findings must be resolved before merge.** The author
  fixes every finding (A/B/C/D) on the branch before merge. No
  "approved with follow-up." The only exception is a finding that
  requires a design decision outside the issue's scope, which the
  reviewer explicitly names as "deferred by design scope" and the
  author files as an issue before merge.

### Frozen snapshot rule

After release, version directories are frozen by repository policy.
Once a release ships:

- The version-snapshot directory at `docs/{tier}/{bundle}/{X.Y.Z}/`
  is read-only by convention.
- The moved cycle directory at `.cdd/releases/{X.Y.Z}/{N}/` (destination:
  `.cds/releases/{X.Y.Z}/{N}/` post-re-rooting) is read-only by
  convention.
- The PRA at `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`
  is read-only by convention.
- The CHANGELOG row for the release is read-only by convention.

**Only path-reference repairs are allowed after freeze:** markdown
links and backtick paths may be updated to fix stale references. No
semantic content may change. A finding that the frozen content is wrong
is a `cds-protocol-gap` per §Field 5 and is closed by a follow-up cycle
that authors a corrective artifact at a new path (or marks the frozen
artifact's defect in the next cycle's PRA), not by editing the frozen
content directly.

The frozen-snapshot rule is the CDS-side realization of the CCNF kernel
scope-lift invariant: once a cycle has projected as α-matter at the
parent scope (the release), the cycle's local-scope artifacts are
witnesses, not editable matter. Subsequent corrections happen at
scope `n+1`, not by rewriting scope `n`.

### Operational realization

The terminology, bootstrap rule, ordered flow, manifest, location
matrix, ownership matrix, CDS Trace format, supporting rules, and
frozen-snapshot rule above are CDS's canonical artifact-contract
statement. The v0.1 operational overlay — the per-artifact authoring
mechanics, the closure-gate row-by-row check sequence, the
release-time directory-move command sequence, the
`scripts/validate-release-gate.sh --mode pre-merge` mechanical check —
lives in the existing cdd role/runtime skills as the temporary v0.1
overlay until the v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/release/SKILL.md §2.5a`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — Release-time cycle-directory move (`.cdd/unreleased/{N}/` →
  `.cdd/releases/{X.Y.Z}/{N}/`); the operational realization of the
  §Location matrix release-time-move column and the §Frozen snapshot
  rule lock-in moment.
- [`cnos.cdd/skills/cdd/release/SKILL.md §3.8`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — Configuration-floor caps tied to actor configuration (cited from
  §Field 6 actor collapse rule); the closure-gate override that forces
  `C_Σ` to `<C` when required artifacts are absent.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.10`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ's closure gate (the row-by-row check of every "Required before
  γ closure" entry in §Ownership matrix above).
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.6–§2.9`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ's release-preparation steps (`RELEASE.md` author, cycle-directory
  move, close-out triage, PRA author, cycle-iteration trigger
  assessment).
- [`cnos.cdd/skills/cdd/alpha/SKILL.md §2.6`](../../../cnos.cdd/skills/cdd/alpha/SKILL.md)
  — α's pre-review gate (the readiness signal in `self-coherence.md`;
  the CDS Trace through Stage 7; AC-to-evidence binding; branch CI
  green on head commit).
- [`cnos.cdd/skills/cdd/beta/SKILL.md`](../../../cnos.cdd/skills/cdd/beta/SKILL.md)
  — β's review CLP, pre-merge gate, merge mechanics, β close-out
  authoring.
- [`cnos.cdd/skills/cdd/release-effector/SKILL.md`](../../../cnos.cdd/skills/cdd/release-effector/SKILL.md)
  — Tag policy (bare `X.Y.Z`); release-time mechanics; the operational
  realization of §Location matrix's tag-format rule and the
  `gamma-closeout.md gates tag` ownership-matrix row.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into
`cnos.cds/skills/cds/{gamma,alpha,beta,release,release-effector,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Mechanical vs judgment

CDS distinguishes between two classes of cycle-coherence check:
**mechanical axes** that tools can enforce structurally (a script returns
pass / fail; no human judgment required) and **judgment axes** that
require an independent reviewer's frame (the failure mode is not visible
to the author at production time; only an outside observer can call it).
The boundary matters because conflating the two produces both directions
of incoherence: mechanizing a judgment axis ships a check that misses
the real failure mode; relying on judgment for a mechanical axis lets
tool-detectable errors accumulate as review debt.

The class boundary is **CDS doctrine**, not implementation discretion.
Field 6 (actor collapse rule) names α=β as never permitted for substantive
software work; the deeper reason is that the judgment axes below require
an outside frame the producer cannot supply. The mechanical axes are
configurable as project-binding CI; the judgment axes are configurable
only as oracle quality, not as tool quality.

### Mechanical axes

Tool-enforceable; the project binding wires a mechanical check that the
cycle either passes or fails. β verifies the check ran; β does not
re-derive the verdict.

- **Branch naming** — `cycle/{N}` canonical per §Development lifecycle →
  §Branch rule. A pre-flight grep / shell-test on the branch name
  (regex `^cycle/[0-9]+$`) is the mechanical oracle. Legacy shapes are
  warn-only at the same surface.
- **Version-directory presence** — for substantial release-shipping
  cycles, `docs/{tier}/{bundle}/{X.Y.Z}/` exists at the bootstrap step
  (Step 3 of §Step table) with a `README.md` snapshot manifest and the
  per-deliverable stubs. A shell directory-existence check is the
  mechanical oracle.
- **AC accounting** — every AC named in the issue body resolves to an
  evidence row in `self-coherence.md §ACs` before β review intake. The
  count of AC headings in the issue body must equal the count of AC
  evidence rows in `self-coherence.md`; mismatched count is a mechanical
  finding β surfaces. The oracle is a grep + count.
- **Stale-cross-reference detection** — links and `§`-anchors in the
  cycle's touched files resolve. A mechanical link-checker (e.g. the
  pre-commit hook implied by §Field 2's "no regressions" oracle on
  documentation surfaces) catches dead references. The oracle is the
  link-checker's exit code.
- **Gate checks** — the closure-gate row-by-row checks per §Gate →
  §Closure verification checklist below (F1–F10). Each F-row is a
  mechanical predicate: file exists / file absent at canonical path.
  `scripts/validate-release-gate.sh --mode pre-merge` is the v0.1
  realization of the gate-check oracle.

### Judgment axes

Non-mechanical; require an independent reviewer's frame. The producer
cannot self-verify these axes at production time — the failure mode is
structurally invisible from the production side. β verifies through
reading, not through running.

- **Real incoherence** — the cycle's §Gap names a coherence-class
  failure the cycle actually closes. A cycle may pass every mechanical
  check while shipping matter that closes no real gap (cycle-for-its-own-sake;
  the no-gap-case judgment from §Selection function). β's judgment
  question: does the merged artifact reduce a coherence delta that was
  named honestly, or does it close a manufactured gap?
- **MCA vs MCI** — the cycle's mode declaration distinguishes between
  shipping an MCA (Master Coherence Action; the gap closes durably this
  cycle) and recording an MCI (Master Coherence Issue; the gap is named
  for a future cycle). The boundary is judgment: an MCA that ships
  surface-only is structurally an MCI mis-labelled. β's judgment
  question: does the matter ship the durability the mode claims?
- **Scoring soundness** — the α/β/γ scores assigned at PRA time per the
  §Field 4 cadence's release artifact. The numeric grade is mechanical;
  the soundness of the grade against what shipped is judgment. β's
  judgment question: do the scores reflect the cycle's actual coherence
  delta, or are they inflated / deflated against an external reference?
- **Review convergence** — when to call review "done" versus "iterate
  another round." The verdict-shape lint (per `review/SKILL.md §3.4a`
  v0.1 overlay) catches mechanical convergence-failure (no `APPROVED`
  with unresolved findings); the substantive convergence question is
  whether one more round would surface incoherence the current set
  missed. β's judgment question: has the review's evidence depth
  matched the cycle's claim strength (per `review/SKILL.md §3.8` v0.1
  overlay)?
- **Design coherence** — whether the cycle's design choices fit the
  surrounding system's invariant / volatile / boundary decomposition
  (the design artifact's structural shape, per `cnos.core/skills/design/SKILL.md`
  v0.1 overlay). A mechanically passing diff may still ship a design
  decision that misshapes the codebase's structure. β's judgment
  question: does the design honour the invariants the surrounding
  system depends on?
- **When to stop iterating** — the cycle-iteration trigger thresholds
  (review rounds > 2; mechanical ratio > 20%; per §Assessment →
  §Cycle iteration triggers below) are mechanical; the judgment of
  whether iterating one more round is structurally productive (versus
  signalling a `cds-skill-gap` or `cds-protocol-gap`) is not.
  γ's judgment question: does the cycle's iteration profile reveal a
  protocol gap, or is the cycle on a healthy trajectory?

### Class invariants

- A mechanical axis that requires human judgment to interpret its
  output is mis-classified — the oracle is incomplete. The fix is
  tooling, not judgment instruction.
- A judgment axis that two reviewers consistently disagree on is a
  `cds-skill-gap` per §Field 5 — the review-skill is under-specified
  for the class of judgment. The fix is a skill patch (per
  `review/SKILL.md §3.12` v0.1 overlay: "review divergence is a skill
  gap"), not "be more careful."
- A finding can carry both classes (e.g. a stale path is a mechanical
  finding; the design coherence the stale path obscures is a judgment
  finding). The finding taxonomy in `review/SKILL.md` (mechanical /
  judgment / contract / honest-claim) permits both classes on a single
  finding row.
- Mechanical findings reaching review are **process bugs**. If >20% of
  findings in a cycle are mechanical, the §Assessment → §Cycle
  iteration triggers fires (Trigger 2: mechanical overload); the
  cycle's process learning fires a `cds-tooling-gap` per §Field 5.

### Operational realization

The mechanical/judgment axis taxonomy above is CDS's canonical statement
of where tools may enforce and where independent judgment is required.
The v0.1 operational overlay — the specific mechanical-check scripts
that realize each axis, the per-finding mechanical-vs-judgment
disposition workflow β runs at review time, the divergence-as-skill-gap
mechanics γ runs at PRA time — lives in the existing cdd role/runtime
skills as the temporary v0.1 overlay until the v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/review/SKILL.md §"Finding Taxonomy"`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — the mechanical / judgment / contract / honest-claim taxonomy β
  applies at review time; the rule that mechanical findings > 20% are
  process bugs.
- [`cnos.cdd/skills/cdd/review/SKILL.md §3.12`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — "review divergence is a skill gap" (judgment-axis divergence as
  `cds-skill-gap` signal).
- [`cnos.cdd/skills/cdd/review/SKILL.md §3.13`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — honest-claim verification (a/b/c/d sub-checks); the boundary between
  mechanical-claim verification (grep-checkable: rule 3.13c wiring
  claims) and judgment-claim verification (rule 3.13b source-of-truth
  alignment, rule 3.13d gap claims) that lives along the mechanical /
  judgment axis-class boundary.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §3.9`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — managerial-residue sweep (the binding γ check before any γ-skill
  patch); the judgment-axis enforcement that catches "be more careful"
  recommendations and rejects them.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{review,gamma,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Review CLP

CDS's review surface uses a **Closed-Loop Predicate (CLP)** form — a
contract whose three parts (TERMS, POINTER, EXIT) make the review's
scope, evidence frame, and termination condition explicit before β reads
a single line of diff. The CLP form is software-class — research-class
CDR uses a different review CLP shape because the loss function differs
(engineering review terminates on contract-coherence + execution-pass;
research review terminates on claim-evidence alignment under uncertainty).

The CLP is paired with the **reviewer ask list** — the structured output
β emits at the end of each review round. The CLP names what β reads;
the ask list names what β returns.

### CLP form

Every β review round opens with a CLP statement that pins three fields
before evidence reading begins:

- **TERMS** — the contract β reads against. Concretely: the issue body's
  AC set; the implementation-contract axes pinned by δ at dispatch
  (per §Field 4 inward-membrane discipline; per
  [cnos#393](https://github.com/usurobor/cnos/issues/393) Rule 7);
  the non-goals enumerated in the issue body; the loaded-skill set the
  cycle's α declared in the §Mode section of `self-coherence.md`.
  TERMS is the **closed surface** β reviews against — findings that
  reach beyond TERMS are out-of-scope and are filed as separate issues,
  not as findings on this round.
- **POINTER** — the evidence surface β reads. Concretely: the cycle
  branch's diff against `origin/main`; the cycle directory's artifact
  set (`self-coherence.md`, prior `beta-review.md` rounds if any,
  `gamma-scaffold.md`, `gamma-clarification.md` if mid-flight pinning
  fired); the CI status on review-SHA; the loaded skills' content; any
  cross-referenced canonical surfaces the cycle touches. POINTER names
  the **evidence frame** — claims grounded outside POINTER are
  phantom-blocker findings (rejected per `review/SKILL.md §3.5` v0.1
  overlay).
- **EXIT** — the termination condition. Concretely: every AC has an
  evidence row resolving to PASS or named as a binding finding; every
  implementation-contract axis pinned by δ resolves to coherent or
  drift-flagged; β emits a **single terminal verdict** per round
  (`APPROVED` or `REQUEST CHANGES`); split verdicts and conditional
  qualifiers are forbidden (per `review/SKILL.md §3.4a` v0.1 overlay).
  EXIT names the **shape** of the round's verdict — a round that
  produces matter outside the EXIT shape is structurally invalid and
  is re-emitted.

The CLP is stated explicitly in `beta-review.md` per round (the round
header carries the TERMS / POINTER / EXIT triple). Stating the CLP is
not a ceremony — it is the contract β binds against. A review that
reads without a stated CLP is structurally a re-derivation of the
contract from the diff, not a verification of the diff against a
declared contract; β cannot detect drift from a contract β did not
state.

### Reviewer ask list

The ask list is β's structured output per round, written into
`beta-review.md` alongside the per-AC evidence rows and the findings
table. The list carries six fields, in this order:

- **α/β/γ scores (per-axis)** — the round's α / β / γ grade per the
  rubric in `release/SKILL.md §3.8` v0.1 overlay (A / A- / B+ / B / C+ /
  C / `<C`). The provisional scores are written here; γ finalises them
  at PRA time per the provisional-vs-final scoring rule
  (`release/SKILL.md §2.4`).
- **Weakest-axis diagnosis** — α, β, or γ named explicitly as the
  weakest axis for this round, with one-paragraph rationale citing the
  finding(s) that lowered the axis. The weakest-axis signal is what
  §Selection function reads at the next cycle's observation step
  (selection rule "Weakest-axis rule").
- **Concrete patch suggestions** — for each binding finding (D / C),
  a named patch path: a file + line range + the substantive change the
  fix requires. The patch suggestion is structural, not a code snippet
  — α produces the code; β names the gap. "Be more careful" is
  rejected per `review/SKILL.md §3.12` v0.1 overlay (review divergence
  is a skill gap, not a careful-reading exhortation).
- **Iterate-or-converge verdict** — `iterate` (RC with named fix
  rounds) or `converge` (APPROVED + merge instruction). The verdict
  carries the round number so subsequent rounds can chain
  (`Round 2 of 2`). A `converge` verdict closes the search space per
  `review/SKILL.md §3.7` v0.1 overlay; an `iterate` verdict names the
  next-round entry condition (which findings, in which file, with
  which oracle).
- **Configuration-floor flags** — if the cycle ran under
  configuration-floor conditions (`operator/SKILL.md §5.2` single-session
  δ-as-γ; β-α-collapse-on-δ from §Field 6), the floor caps applied
  this round (γ ≤ A-, β ≤ A- under collapse; γ ≤ A- under §5.2;
  γ ≤ C under CI-red; γ ≤ B- under CI-skip; γ ≤ B under false-gap
  peer-enumeration). The flag carries the rule clause that fired
  (`release/SKILL.md §3.8` configuration-floor / CI-red / false-gap
  clause).
- **Empirical-anchor cite (optional)** — when a finding's pattern
  reproduces a known empirical-anchor cycle's failure mode, β cites
  the anchor (e.g. cnos#391 for wrong-shape γ implementation;
  cnos-tsc supercycle for honest-claim drift). The cite is structural
  signal for ε's receipt-stream review (per §Field 5 cadence) —
  recurring anchor citations across multiple cycles surface as
  `cds-skill-gap` patterns.

### Operational realization

The CLP form and the reviewer ask list above are CDS's canonical
review-surface statement. The v0.1 operational overlay — the per-phase
β procedure (Phase 1 contract integrity, Phase 2 implementation review,
Phase 3 verdict), the finding taxonomy (mechanical / judgment / contract
/ honest-claim) β tags each finding with, the severity rubric (D / C /
B / A) β assigns, the per-round output template, the merge mechanics
β executes on `converge`, the post-merge close-out authoring β runs —
lives in the existing cdd role/runtime skills as the temporary v0.1
overlay until the v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/review/SKILL.md`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — the review orchestrator; the three-phase procedure
  (`review/contract/SKILL.md`, `review/issue-contract/SKILL.md`,
  `review/diff-context/SKILL.md`, `review/architecture/SKILL.md`); the
  verdict rules (3.1–3.13); the §Output Format template; the §Checklist.
- [`cnos.cdd/skills/cdd/review/SKILL.md §3.4a`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — verdict-shape lint (`APPROVED` + unresolved findings; conditional
  qualifier; split verdicts); the recovery path that reformats
  conditions as required-fix findings.
- [`cnos.cdd/skills/cdd/review/SKILL.md §3.10`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — CI-green gate (β verifies required CI workflows green on review-SHA
  before emitting APPROVED).
- [`cnos.cdd/skills/cdd/review/SKILL.md §3.11b`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — γ-artifact-completeness gate (β verifies `gamma-scaffold.md` exists
  on cycle branch before APPROVED).
- [`cnos.cdd/skills/cdd/review/SKILL.md §3.13`](../../../cnos.cdd/skills/cdd/review/SKILL.md)
  — honest-claim verification (a/b/c/d sub-checks for reproducibility,
  source-of-truth alignment, wiring claims, gap claims).
- [`cnos.cdd/skills/cdd/beta/SKILL.md`](../../../cnos.cdd/skills/cdd/beta/SKILL.md)
  — β role contract; pre-merge gate; merge mechanics; β close-out
  authoring.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{review,beta,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Gate

The **gate** is the release-boundary check that protects a cycle's
transition from "merged" to "tagged + released." A cycle that merges
matter into `main` but bypasses the gate ships unreleased — the trust
boundary on which downstream consumers depend (per the CCNF kernel's
Scope-Lift Invariant: a closed cycle projects as α-matter at the parent
scope only after the boundary effection records the boundary decision)
has not been crossed. The gate names the **preconditions** that must be
true before the release-effector mechanics run, and the **closure
verification checklist** — the F1–F10 failure modes the gate catches.

The gate is δ-owned (δ runs the release-boundary preflight per §Field 4
outward-membrane discipline) but γ-prepared (γ stages the artifacts the
gate verifies per §Step 12 of §Step table). The actor split matters:
γ cannot self-verify the gate, because γ produced the matter γ is
gating; δ verifies as an outside frame.

### Release-readiness preconditions

Before δ runs the release-effector mechanics, **every** precondition
below must be observable on `origin/main`:

- **Merge commit present.** β's `git merge` into `main` with `Closes #N`
  has landed; the merge commit's SHA is recoverable; CI ran green on the
  merge commit (per Field 2's "no regressions" oracle on adjacent
  surfaces).
- **CI green on merge SHA.** Required CI workflows have
  `conclusion == "success"` on the merge commit's SHA. A red or pending
  required workflow blocks the gate; the recovery is the release-effector
  CI-red runbook (per `release-effector/SKILL.md` v0.1 overlay).
- **`RELEASE.md` committed to `main`.** The GitHub release body lives at
  the repo root, included in the release commit. Missing RELEASE.md
  triggers F5 (below) and the gate blocks.
- **Cycle directories moved.** Every `.cdd/unreleased/{N}/` for cycles
  closed in this release has been moved to `.cdd/releases/{X.Y.Z}/{N}/`
  and the move is part of the release commit (per §Artifact contract
  → §Location matrix). Stale `.cdd/unreleased/{N}/` after the release
  triggers F4 and the gate blocks.
- **All close-outs present.** `alpha-closeout.md`, `beta-closeout.md`,
  `gamma-closeout.md` all live at the canonical paths on `main`. Missing
  any close-out triggers F1 / F2 / F3 respectively.
- **PRA committed.** The `POST-RELEASE-ASSESSMENT.md` exists at the
  canonical version-directory path (`docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`).
  Missing PRA triggers F8.
- **Cycle iteration artifact present when triggered.** If the cycle's
  receipt has `protocol_gap_count > 0`, `.cdd/unreleased/{N}/cdd-iteration.md`
  (or `cds-iteration.md` post-rename) exists with each finding
  structured per §Field 5's per-finding shape, and
  `.cdd/iterations/INDEX.md` has a row for cycle N. Missing iteration
  artifact under triggered conditions fires F9. The courtesy empty-findings
  stub per cnos#401 cadence rule is permitted when count is 0.
- **No unresolved triage.** Every finding from `alpha-closeout.md`,
  `beta-closeout.md`, and the PRA carries a disposition (immediate MCA
  / project MCI / agent MCI / drop). Unresolved triage triggers F10.
- **γ closure declaration on `main`.** `gamma-closeout.md` is the
  closure-declaration artifact; its commit hash on `main` is the
  signal that γ has declared closure. δ does **not** tag before this
  signal (per the operator/SKILL.md §3.4 doctrine; per F6).

The preconditions are conjunctive: all must hold simultaneously.
A preflight that finds any precondition false returns "Request changes"
(δ routes the gap back to γ); a preflight that finds all preconditions
true returns "Proceed" (δ runs the release-effector mechanics); δ may
override per the override doctrine (§Field 4 boundary decision) with
the override block recorded on the receipt.

### Closure verification checklist

The 10 failure modes below — F1 through F10 — are the gate's named
checks. Each F-row is a mechanical predicate (file exists / file absent
at canonical path; commit ordering; field present) whose oracle is the
gate-check script (`scripts/validate-release-gate.sh --mode pre-merge`
in the v0.1 realization). The F-anchors are **stable cross-reference
targets**: cdd skill files (and post-#403, the CDS-side role overlays)
cite the F-anchors when naming which closure-gate row a finding maps
to; Sub 6 of cnos#403 re-points the citations from
`gamma/SKILL.md §2.10` / `release/SKILL.md` at these anchors.

#### F1: Missing α close-out

`.cdd/unreleased/{N}/alpha-closeout.md` does not exist on `origin/main`
at gate-check time. α's role-local findings record is absent; γ cannot
complete the close-out triage (per §Closure below); the cycle's α-axis
signal is unrecorded.

**Recovery:** γ requests δ to re-dispatch α via the close-out prompt
(per `cnos.cdd/skills/cdd/CDD.md §1.6a` v0.1 overlay). The provisional
fallback (α writes `[provisional]` at review-readiness; pre-merge gate
accepts the form) is acceptable for non-tagged docs-only cycles but
not for tagged releases — full close-out via δ re-dispatch is the
normative path. F7 is the related anchor for the re-dispatch mechanism.

#### F2: Missing β close-out

`.cdd/unreleased/{N}/beta-closeout.md` does not exist on `origin/main`
at gate-check time. β's review-context and merge-evidence record is
absent; γ cannot complete the close-out triage; the cycle's β-axis
signal is unrecorded.

**Recovery:** γ requests δ to re-dispatch β for close-out authoring.
β authors `beta-closeout.md` in the same session as the merge commit
under the canonical path; re-dispatch is the exception, not the rule.

#### F3: Missing γ close-out (no closure declaration)

`.cdd/unreleased/{N}/gamma-closeout.md` does not exist on `origin/main`
at gate-check time. γ has not declared closure; the cycle has not closed
in the structural sense (it has only stopped). δ MUST NOT tag before
γ closure declaration.

**Recovery:** γ runs the closure gate (§Closure → §Closure rule below);
if every row passes, γ writes `gamma-closeout.md` with the closure
declaration ("Cycle #N closed. Next: #M.") and commits to main. If
any closure-gate row fails, γ remediates the underlying gap (re-dispatch
α/β; author missing artifact) before re-running the gate.

#### F4: Stale `.cdd/unreleased/{N}/` after release

The per-cycle directory at `.cdd/unreleased/{N}/` still exists on
`origin/main` after the release commit. The cycle-directory move
(`.cdd/unreleased/{N}/` → `.cdd/releases/{X.Y.Z}/{N}/`) was skipped or
landed after the release commit (instead of in it). The cycle loses
its version association; the §Frozen snapshot rule (§Artifact contract)
is structurally broken.

**Recovery:** γ moves the directory in a follow-up commit on `main`
referenced from the PRA; the corrective commit cites F4 explicitly.
For future cycles, γ executes the move **before** requesting δ tag,
per `release/SKILL.md §2.5a` v0.1 overlay.

#### F5: Missing `RELEASE.md`

`RELEASE.md` does not exist at the repo root in the release commit.
The GitHub release CI workflow falls back to auto-generated sparse
notes (commit titles only); the coherence-delta framing the release
body is supposed to carry (per `release/SKILL.md §2.5` v0.1 overlay)
is lost; the release ships under-described.

**Recovery:** γ authors `RELEASE.md` per the format spec
(§Outcome / §Why it matters / §Fixed / §Added / §Changed / §Removed /
§Validation / §Known Issues) and commits to `main` **before** requesting
δ tag. The release body is overwritten via `gh release edit` if the
tag has already pushed; the F5 finding is recorded in the cycle's
`gamma-closeout.md`.

#### F6: δ tag ordering violation

δ pushed a tag (`X.Y.Z`) before `gamma-closeout.md` existed on `main`.
The release boundary fired before γ declared closure; the cycle's
structural closure (per §Field 4 outward-membrane discipline) is
inverted — δ recorded the boundary decision before the cycle was
ready to cross the boundary. F6 is the most-severe gate violation
because it cannot be reversed cleanly (a pushed tag is git-observable
to every downstream consumer).

**Recovery:** δ does not delete the tag (force-deletion of pushed tags
breaks every downstream `cn update` consumer). Instead, γ writes the
missing `gamma-closeout.md` immediately, files an F6 finding in
`cdd-iteration.md` per §Field 5's per-finding shape, and the override
block on the cycle's receipt records `transmissibility: degraded` per
the structural override rule (per `schemas/cdd/receipt.cue`). The
debt-cycle that closes the override (per §Field 4 cadence triggers)
re-establishes the structural ordering for future cycles.

#### F7: Missing α close-out re-dispatch mechanism

The cycle's actor configuration cannot re-dispatch α post-merge (e.g.
sequential-bounded-dispatch terminated α's session at review-readiness;
operator-class agents have no re-dispatch primitive available). F1
(missing α close-out) cannot be remediated by re-dispatch because the
re-dispatch mechanism itself is absent. F7 is the **meta-finding**
that F1 is structurally un-remediable at this actor configuration.

**Recovery:** γ accepts the provisional fallback (α writes `[provisional]`
close-out at review-readiness; the pre-merge gate accepts the form
per §Ownership matrix); the cycle is marked as carrying a provisional
close-out in `gamma-closeout.md`; the next cycle's α-axis signal carries
the provisional flag. For tagged release cycles, F7 forces a
configuration-change cycle that lands re-dispatch primitives before
the next release.

#### F8: Missing PRA

`POST-RELEASE-ASSESSMENT.md` does not exist at the canonical
version-directory path
(`docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`). The cycle
has no measured coherence delta; the encoding-lag-table update is
missing; the next cycle's §Selection function has no PRA to read at
its "Assessment-commitment default" rule. The §Closure rule cannot
fire because §Immediate outputs cannot resolve without the PRA's
process-learning section.

**Recovery:** γ authors the PRA per §Assessment → §PRA contents below;
the PRA commit precedes the closure declaration. For docs-only releases,
the PRA lives at `docs/gamma/cdd/docs/{ISO-date}/POST-RELEASE-ASSESSMENT.md`
per the `release/SKILL.md §2.5b` v0.1 overlay form; the canonical path
shape is unchanged.

#### F9: Missing `cdd-iteration.md` when triggers fired

The cycle's receipt has `protocol_gap_count > 0` (≥1 finding tagged
`cds-skill-gap` / `cds-protocol-gap` / `cds-tooling-gap` /
`cds-metric-gap`; or the pre-rename `cdd-*-gap` taxonomy for legacy
cycles), but `.cdd/unreleased/{N}/cdd-iteration.md` is absent or
`.cdd/iterations/INDEX.md` has no row for cycle N. ε's per-cycle
protocol-iteration artifact (per §Field 5 cadence) is missing under
triggered conditions; the protocol cannot measure its own learning
rate.

**Recovery:** γ authors `cdd-iteration.md` (or `cds-iteration.md`
post-rename) per §Field 5's per-finding shape; the aggregator
`.cdd/iterations/INDEX.md` row is appended in the same commit; if
any finding's disposition is `patch-landed` with a cross-repo target,
the cross-repo bundle at `.cdd/iterations/cross-repo/{counterpart}/{slug}/`
is created. The cnos#401 courtesy empty-findings stub is **not** a
remedy for F9 — F9 fires only when triggers fired (count > 0); the
stub applies only when count is 0.

#### F10: Unresolved triage

A finding from `alpha-closeout.md`, `beta-closeout.md`, or the PRA
carries no disposition (no `immediate MCA` / `project MCI` / `agent MCI`
/ `drop` row in γ's triage table). Silence is not a disposition per
`gamma/SKILL.md §3.7` v0.1 overlay; the cycle has unowned coherence
gaps at closure time.

**Recovery:** γ adds a disposition row for every finding before
writing `gamma-closeout.md`. "Drop" is a valid disposition only when
explicit and reasoned — the drop row carries the finding ID and a
one-sentence rationale. A `noted for next cycle` row that does not
file an issue is structurally a `drop` row mis-labelled and fires
F10 on re-check.

### Operational realization

The release-readiness preconditions and the F1–F10 closure verification
checklist above are CDS's canonical gate statement. The v0.1 operational
overlay — the row-by-row gate-check script, the per-precondition
verification command sequence, the δ release-boundary preflight
procedure, the F-anchor recovery runbooks, the merged-branch deletion
mechanics — lives in the existing cdd role/runtime skills as the
temporary v0.1 overlay until the v1 CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/release/SKILL.md`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — the release-readiness-check procedure (§2.1); the version-decision
  procedure (§2.2); the release-commit + tag-push sequence (§2.6); the
  CI-wait + deploy + validate sequence (§2.7–§2.9); the
  configuration-floor clause (§3.8) that ties closure-gate failures to
  the C_Σ override.
- [`cnos.cdd/skills/cdd/release/SKILL.md §2.5a`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — cycle-directory move (`.cdd/unreleased/{N}/` →
  `.cdd/releases/{X.Y.Z}/{N}/`); the F4 mechanical realization.
- [`cnos.cdd/skills/cdd/release/SKILL.md §2.5b`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — docs-only disconnect (no tag); the docs-only release path that
  shares F1–F10 verification but bypasses the tag-ordering F6 check.
- [`cnos.cdd/skills/cdd/release/SKILL.md §3.8`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — closure-gate override (forces `C_Σ` to `<C` when any close-out
  artifact is absent at merge time); the F1 / F2 / F3 mechanical
  realization. Verified by `scripts/validate-release-gate.sh --mode
  pre-merge` per cnos#339 AC1.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.10`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ's closure gate; the 14-row row-by-row check (rows 1–14) that
  realizes F1–F10 across the cycle's closure preconditions; the
  closure-declaration commit that signals F3 PASS.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.6–§2.9`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ's release-preparation steps (`RELEASE.md` author per F5; cycle-directory
  move per F4; close-out triage per F10; PRA author per F8;
  cycle-iteration trigger assessment per F9).
- [`cnos.cdd/skills/cdd/operator/SKILL.md §3`](../../../cnos.cdd/skills/cdd/operator/SKILL.md)
  — δ-as-operator gate doctrine; "Do not tag/release before
  `gamma-closeout.md` exists on main" (F6); "δ blocks release completion
  until CI is green or operator explicitly accepts a known pre-existing
  failure"; the gate-action-on-request-not-observation rule.
- [`cnos.cdd/skills/cdd/operator/SKILL.md §3.1`](../../../cnos.cdd/skills/cdd/operator/SKILL.md)
  — the external-actions table (pre-merge gate validation, push merge,
  release-boundary preflight, tag push + release, branch delete,
  issue filing on external repos, force push, auth refresh).

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{release,gamma,delta,operator,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Assessment

The **post-release assessment (PRA)** is γ's per-cycle observation
artifact at the canonical version-directory path
(`docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`). One PRA
per release; no exceptions. The PRA measures what the release actually
shipped against what the cycle's contract claimed, records the
encoding-lag table, runs the cycle-iteration trigger check, and
commits the next-move concretely. The PRA feeds §Selection function's
"Assessment-commitment default" rule at the next cycle's observation
step; a release without a PRA is structurally a release without a
next-cycle handoff.

The PRA is **γ-owned**, not β-owned, by structural argument: the PRA
measures β's review quality, and β assessing β's own review quality
is the self-grading problem the role separation is designed to avoid
(per §Field 6 actor collapse rule). γ holds the cycle-level observational
authority — γ sees both α's matter and β's review, but γ does not
author either. γ writes the PRA after β's merge and β close-out have
landed.

### PRA contents

The PRA carries seven canonical sections plus the trace + hub memory
row. Each section is required; explicit "not applicable" with rationale
is permitted only on sections that name themselves as conditional. The
section vocabulary below is the canonical CDS shape; the per-section
authoring detail (what each field measures, what the table cells
contain) lives in the v0.1 operational overlay.

- **§1 Coherence Measurement** — baseline (previous release's α / β /
  γ scores from the CHANGELOG TSC table); this release's α / β / γ
  scores; the **coherence delta** (which axes improved / held /
  regressed and why); the coherence-contract-closed verdict (was the
  cycle's named gap actually closed? if not, what remains?). This
  section anchors the §Selection function's weakest-axis rule for the
  next cycle.
- **§2 Encoding Lag** — the cross-cycle backlog of converged-but-not-shipped
  design commitments. Every open design issue and every open process
  issue appears as a row (Issue / Title / Type / Design / Impl / Lag).
  Lag levels: `none` (shipped this or prior release), `low`
  (implementation in progress), `growing` (design converged, no
  implementation started), `stale` (design aging without implementation
  plan). The **MCI/MCA balance decision** (balanced / freeze MCI /
  resume MCI) is mandatory; every release states the balance with
  rationale.
- **§3 Process Learning** — three questions: what went wrong, what
  went right, what skill patches are needed. Active skill re-evaluation:
  for each review finding, would the loaded skills (as written) have
  prevented it? Underspecified → patch the skill (§Step 5 of the v0.1
  procedure); application gap → note it. Explicit CDS-improvement
  disposition is mandatory: patch landed (description + commit) OR
  no-patch-needed (justification per the three valid justifications:
  application gaps with adequate spec, zero findings, or environmental
  failure with no spec-level fix).
- **§4 Review Quality** — review-rounds-per-cycle counts (target ≤1
  docs, ≤2 code); superseded cycles count (target 0); finding-class
  breakdown (mechanical / wiring / honest-claim / judgment / contract);
  mechanical-ratio (threshold 20% with ≥10 findings — fires §Cycle
  iteration triggers Trigger 2); honest-claim ratio (target <30%; high
  ratio means α docs drifting from artifacts). Section §4 closes the
  loop on β's per-cycle review-axis signal.
- **§4a CDS Self-Coherence** — CDS-axis scoring of the cycle itself:
  CDS α (artifact integrity — required artifacts present? bootstrap /
  frozen snapshot complete?); CDS β (surface agreement — canonical doc,
  executable skill, cycle artifacts, changelog, and assessment agree?
  authority conflicts or stale references?); CDS γ (cycle economics —
  review rounds within target? mechanical ratio under threshold?
  immediate outputs executed?). Weakest-axis named; action stated
  (patch skill / patch doc / automate check / none).
- **§5 Production Verification** — a concrete executable scenario that
  demonstrates the new capability (or blocked failure mode) in a real
  environment. "CI passes" is not production verification — it is
  build verification. If the change is structural (L7), the scenario
  demonstrates the boundary move, not just a single path. Result:
  pass / fail / deferred (with commitment to when/how).
- **§7 Next Move** — the concrete next-cycle commitment: next MCA
  issue number, owner, target branch name, first AC to close, MCI
  freeze/resume state, rationale. This is the §Closure → §Deferred
  outputs surface; the PRA's §7 is what the next cycle's §Selection
  function reads at "Assessment-commitment default" time. Closure
  evidence: immediate outputs executed (with links/commits); deferred
  outputs committed (issue # / owner / branch / first AC / freeze
  state per outpost).

Two additional rows are part of the PRA's structural surface:

- **§4b Cycle Iteration** — the §9.1 trigger assessment (see §Cycle
  iteration triggers below). Required when any trigger fired; permitted
  as "No §9.1 trigger fired" when none did.
- **§8 Hub Memory** — the operational state the next session loads to
  orient (daily reflection path + commit SHA; adhoc thread(s) updated
  + commit SHA). Hub memory is external to the repo but the PRA records
  the path / commit-sha / unavailable-reason; verifiers do not inspect
  the external hub directly.

### Cycle iteration triggers

A §9.1 trigger is a per-cycle process-learning signal that the cycle's
protocol substrate (its skills, its tooling, its operator configuration)
needs patching. When a trigger fires, the PRA's §4b Cycle Iteration
section is mandatory; γ records the trigger, the root cause, the
disposition (patch landed / next MCA / no-patch with reason), and the
evidence (commit / issue / note). The four canonical triggers below
are the **§9.1 anchor set** — citing skill files (and post-#403, the
CDS-side role overlays) reference them by trigger number; Sub 6 of
cnos#403 re-points the `CDD.md §9.1` citations at this section's
anchors.

1. **Review churn — review rounds > 2.** The cycle required more than
   two rounds of β review before merge. Signal: the contract was
   under-specified (the issue body's ACs were ambiguous), or the AC
   oracle was insufficient, or α's skill was under-specified for the
   matter class. Disposition: γ patches the issue-quality gate (per
   `gamma/SKILL.md §2.4` v0.1 overlay) when the gap is clear, or
   files a next-MCA naming the specific gap.

2. **Mechanical ratio > 20% (with ≥10 findings).** Mechanical-class
   findings (per the §Mechanical vs judgment axis-class above) exceed
   20% of total findings on a single cycle, and total findings ≥10
   (the floor prevents noise on small cycles). Signal: tooling-class
   gap — a check that should be in CI is being run by β manually; a
   pre-commit hook is missing; a mechanical validator does not exist
   for the surface class. Disposition: γ files a `cds-tooling-gap`
   per §Field 5, with concrete next-MCA naming the missing check
   (validator name, input surface, oracle).

3. **Avoidable tooling / environmental failure.** The cycle was
   blocked by a tooling or environment failure that a CDS tooling
   artifact could have prevented (e.g. the harness's polling
   `git fetch` reliability re-probe per §Coordination surfaces →
   §Polling primitives; the dispatch timeout budget per `CDD.md §1.6c`
   v0.1 overlay; a missing CI workflow that the cycle's δ-pinned
   contract required). Signal: `cds-tooling-gap`. Disposition: γ patches
   the tooling guard now when the fix is clear (shell snippet, hook,
   config change), files a next-MCA when the fix requires design.

4. **Loaded skill failed to prevent a finding.** A finding surfaced
   despite the relevant skill being loaded; the skill's procedure did
   not catch the class of error. Signal: `cds-skill-gap` — the
   skill's procedure is under-specified for the failure-mode
   geometry. Disposition: γ patches the skill in the same commit as
   the PRA (per CDD §13a's "patch in the same commit" rule) when the
   patch is clear; otherwise files a concrete next-MCA naming the
   specific skill, the specific procedure step, and the failure-mode
   geometry the skill missed.

Each fired trigger must end in exactly one of three states: **patch
landed now**, **concrete next MCA committed**, or **explicit no-patch
decision with reason**. Silence is not a disposition; "noted for next
cycle" is structurally a `next-MCA` row mis-labelled (and fires F10
of §Gate's closure verification checklist on re-check).

The four-trigger set is intentionally narrow. Two additional patterns
are valid signals but are routed differently:

- **CI red on merge commit (post-merge)** — appears as a §4 Review
  Quality observation (per the post-release/SKILL.md §4b row) and as
  a §Gate F-anchor recovery path (red on merge SHA fires the gate's
  CI-green precondition). The CI-red configuration-floor cap
  (γ ≤ C per `release/SKILL.md §3.8`) fires regardless of whether a
  §9.1 trigger fires.
- **AC oracle ambiguity recurrence across cycles** — signal that
  emerges from ε's receipt-stream review (§Field 5), not from a single
  cycle's PRA. The PRA records the per-cycle observation; ε aggregates
  across cycles to detect recurrence.

### Friction log

The PRA carries a **friction log** entry for each non-trigger
process-friction observation γ surfaces during cycle close-out. The
friction log is the prose-shaped record that feeds the next cycle's
§Selection function (stale-backlog re-evaluation rule) and ε's
receipt-stream review.

Each friction-log row carries four fields:

- **Root cause classification** — what kind of friction this is:
  process (cycle-coordination friction; γ-axis signal); skill
  (the loaded skill set did not cover the case; α-axis or β-axis
  signal); tooling (a mechanical check is missing or wrong;
  γ-axis signal); environmental (sandbox, network, identity wiring;
  out-of-system signal); contract (the issue body's contract was
  ambiguous; γ-axis signal).
- **Skill impact** — which skill(s) the friction implicates, named
  by path (e.g. `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5`). A
  friction-log row that names no skill is structurally an environmental
  finding mis-classified; environmental findings explicitly carry
  "no skill" with a rationale.
- **MCA (next master coherence action)** — the concrete next step:
  issue # (filed or to-be-filed); owner; first AC. If the friction
  does not warrant an MCA (one-off; environmental with no spec-level
  fix), the row carries "no MCA" with rationale.
- **Disposition** — the same four-state shape as cycle-iteration
  triggers (patch landed / next MCA / no-patch / drop). The friction
  log is structurally lighter than §4b — it surfaces below-threshold
  observations that did not trip a §9.1 trigger but still warrant
  γ's attention.

### Engineering levels

The PRA records the cycle-level **engineering level** (L5 / L6 / L7)
per the canonical framework at
[`docs/gamma/ENGINEERING-LEVELS.md`](../../../../../docs/gamma/ENGINEERING-LEVELS.md).
The Level column appears in the CHANGELOG TSC ledger row (per
§Artifact contract → §Ownership matrix); the PRA names the level and
its rationale.

The three levels (cited from `ENGINEERING-LEVELS.md`):

- **L5 — Local Correctness.** The cycle ships matter that is locally
  correct: the implementation produces the right output for the named
  inputs; the test suite passes; the contract's ACs are met. L5 is
  the floor for any CDS cycle — a cycle that does not reach L5 has
  not shipped.
- **L6 — System-Safe Execution.** The cycle ships matter that is
  safe-to-execute in the surrounding system: no regressions on
  adjacent surfaces; documented behaviours preserved; CI green on
  merge commit; the implementation-contract axes pinned by δ at
  dispatch are honoured. L6 is the floor for any tagged release —
  a sub-L6 cycle ships as `<C` per `release/SKILL.md §3.8`
  configuration-floor.
- **L7 — System-Shaping Leverage.** The cycle ships matter that
  changes the surrounding system's structure — boundary moves, new
  invariants pinned, deprecated surfaces removed, architectural
  decompositions revised. L7 is the level of substantial design
  cycles (per the §Selection function's maximum-leverage rule); L7
  cycles' PRAs carry an L7 essay in `docs/papers/` when the
  structural change is large enough to warrant standalone
  rationale-record.

The PRA's Level row names the level and a one-paragraph rationale
citing the structural reach of the cycle. Score-the-release-not-the-intent
applies: an L7-intended cycle that ships only L5 matter records L5
with rationale.

### Operational realization

The PRA contents, the four §9.1 cycle-iteration triggers, the friction
log, and the L5/L6/L7 framework above are CDS's canonical assessment
statement. The v0.1 operational overlay — the per-section authoring
procedure, the scoring rubric and geometric-mean computation, the
encoding-lag-table population workflow, the MCI/MCA balance decision
procedure, the hub-memory write sequence — lives in the existing cdd
role/runtime skills as the temporary v0.1 overlay until the v1 CDS-side
role rewrite:

- [`cnos.cdd/skills/cdd/post-release/SKILL.md`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — the PRA procedure (Steps 1–7); the output-template (§1 Coherence
  Measurement through §8 Hub Memory); the pre-publish gate (the
  mechanical checklist that two agents independently verifying the
  same template must produce the same list of missing fields).
- [`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — CDD self-coherence scoring (CDS α / β / γ axes); the weakest-axis
  action selection.
- [`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6a`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — §9.1 cycle-iteration per-trigger authoring (name trigger / root
  cause / disposition / evidence).
- [`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — `cdd-iteration.md` (or `cds-iteration.md` post-rename) authoring;
  the per-finding shape; the aggregator update; the cross-repo trace
  rule.
- [`cnos.cdd/skills/cdd/post-release/SKILL.md §5.7`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — production-verification scenario authoring; the executable-not-hypothetical
  discipline.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.7`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — close-out triage; the four-disposition shape (immediate MCA /
  project MCI / agent MCI / drop); the post-merge CI-verification
  precondition.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.8`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — cycle-iteration-trigger enforcement (per-trigger γ action; closure
  rule); the four-trigger table that realizes the §9.1 triggers
  operationally.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.9`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — independent γ process-gap check (even when no §9.1 trigger fired,
  γ asks the four questions; states why-not in one sentence when no
  patch is needed).
- [`docs/gamma/ENGINEERING-LEVELS.md`](../../../../../docs/gamma/ENGINEERING-LEVELS.md)
  — the L5 / L6 / L7 framework cited above; the per-level definitions,
  typical behaviors, strengths, limits, and signals.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{post-release,gamma,assessment,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Closure

**Closure** is the final cycle state — the structural transition from
"merged + released" to "cycle complete." A cycle that has merged matter
and crossed the release boundary has not yet *closed* in the CDS sense
— closure requires that **immediate outputs** are executed, **deferred
outputs** are committed, and (when triggered) cycle-iteration artifacts
are written. The Closure rule below names the conjunction; γ declares
closure in `gamma-closeout.md` only after the conjunction holds.

Closure is the discipline that prevents the "stopped, not closed"
failure mode — a cycle whose matter ships but whose downstream
coherence work (skill patches, hub-memory updates, next-cycle handoff,
process learning) trails into next sessions and accumulates as
process debt. The §Closure rule is structurally the same rule
§Gate F10 (unresolved triage) catches at the gate; F10 is the
release-time mechanical realization of the Closure rule's "no
unresolved triage" precondition.

### Immediate outputs

Outputs executed **in the same session** as the PRA / `gamma-closeout.md`
commit. Immediate outputs close the per-cycle learning loop before the
cycle records itself as closed; deferring them into the next cycle is
structurally a `cds-protocol-gap` (silence is not a disposition; the
§Closure rule does not fire on a cycle whose immediate outputs are
"noted for next cycle").

The seven canonical immediate-output classes:

- **Changelog corrections.** Any TSC ledger row drift between β's
  release-time provisional scores and γ's PRA-time final scores;
  any missing Level / Rounds column entry; any provisional marker
  still present after γ's PRA. The change lands in the same commit
  as the PRA per the provisional-vs-final scoring rule
  (`release/SKILL.md §2.4` v0.1 overlay).
- **Missing documentation.** Documentation surfaces that describe
  current behaviour but lagged the cycle's behaviour change. Updated
  in the same commit as the PRA (per §Artifact contract → §Supporting
  rules: "update docs before release" — when documentation lag is
  caught at closure rather than at merge, the immediate output is
  the patch that closes the lag).
- **Skill / process micro-patches.** Patches to cdd or cds skills
  whose under-specification caused a finding this cycle. Micro-patches
  are landed when the gap is clear; larger patches become next-MCA
  rows in §Deferred outputs. The micro-patch is synced across the
  canonical source (`src/packages/`), any package-visible loader
  entrypoint (`SKILL.md`), and any human-facing pointer surface that
  exposes the changed rule (per `post-release/SKILL.md §5` Step 5).
- **Skill patches from cycle iteration.** When §Assessment → §Cycle
  iteration triggers fired and the disposition is `patch landed now`,
  the patch lands in the same commit as the PRA (per CDD §13a's
  same-commit rule). A `patch landed now` disposition with the patch
  deferred to a separate commit is structurally a `next-MCA`
  disposition mis-labelled.
- **Issue filing.** Any next-MCA whose first-AC is concrete enough
  to file as a GitHub issue is filed in the same session, with the
  filed issue number recorded in §Deferred outputs below. An MCA
  named in the PRA without a filed issue is structurally a deferred
  output without commitment (and fires F10 if it lacks disposition).
- **Lag-table updates.** The encoding-lag table at the PRA's §2 is
  populated with every open design / process issue, with current
  lag levels (none / low / growing / stale); shipped-this-release
  issues drop off; new converged-but-unimplemented designs are
  added. The lag-table update happens in the PRA commit, not in a
  separate session.
- **Hub-memory writes.** Daily reflection (the operational state
  the next session loads to orient) and adhoc-thread updates (which
  ongoing thread(s) this release advances). Hub memory is external
  to the repo; the PRA records the path + commit-sha + unavailable-reason
  per §Assessment → §PRA contents §8. Skipping hub memory creates
  a compaction gap that the next session cannot recover from.

### Deferred outputs

Outputs **committed** to a downstream cycle. Deferred outputs are not
deferred decisions — they are decisions made this cycle with execution
scheduled for the named downstream surface. The commitment is structural:
a deferred output must carry a downstream identity (issue number,
owner, branch) and a per-cycle first-AC.

The five canonical deferred-output fields:

- **Next-MCA issue number.** The specific GitHub issue number that
  encodes the next master coherence action. If the issue does not yet
  exist, filing it is an immediate output (per §Immediate outputs
  above); the §Deferred outputs row carries the filed-now issue
  number.
- **Owner.** The role-or-persona that owns the next-MCA's execution.
  Typically `γ-the-operator` for cycle-coordination MCAs;
  `cds-skill-owner` for skill patches; a named persona hub for
  persona-class work. An MCA without an owner is structurally an
  un-dispatched cycle (no one runs it).
- **Target branch name.** The cycle/{N} branch name the next-MCA's
  α will work on. `cycle/{N}` per §Development lifecycle → §Branch
  rule (where N is the next-MCA's issue number). A deferred output
  with `pending branch creation` is acceptable when the next-MCA's
  γ has not yet created the branch — γ creates the branch at the
  next cycle's Step 2.
- **First AC.** The concrete first acceptance criterion that ships
  in the next-MCA. The first-AC is what the next cycle's α binds
  against at intake; a next-MCA with an under-specified first-AC
  fires the §Selection function's "issue-quality gate" at the next
  cycle's bootstrap step.
- **MCI freeze / resume state.** The state of the §Selection function
  MCI freeze flag after this cycle (frozen / resumed / no-change).
  `Frozen` means the next substantial MCA must come from the stale
  set (per the §MCI freeze check rule). `Resumed` means implementation
  has caught up to the design frontier. `No-change` means the prior
  state continues.

### Closure rule

A cycle closes **only when all three conjuncts hold**:

1. **Immediate outputs executed** — every item in §Immediate outputs
   above has either landed in this cycle's commits or been explicitly
   ruled out with rationale. "Noted for next cycle" without a filed
   next-MCA row is not "ruled out."
2. **Deferred outputs committed** — every §Deferred outputs field
   is populated; the next-MCA issue exists; the owner is named; the
   target branch name is declared (or `pending branch creation` with
   the next cycle's γ named); the first-AC is concrete; the MCI
   freeze state is declared.
3. **Cycle iteration present if triggered** — if any §Assessment →
   §Cycle iteration triggers fired, the PRA's §4b Cycle Iteration
   section is authored and `cdd-iteration.md` (or `cds-iteration.md`
   post-rename) exists at the canonical path with each finding
   structured per §Field 5's per-finding shape, and
   `.cdd/iterations/INDEX.md` has a row for cycle N. The cnos#401
   cadence rule applies: required only when `protocol_gap_count > 0`;
   courtesy empty-findings stub permitted when count is 0.

When all three conjuncts hold, γ writes `gamma-closeout.md` with the
**closure declaration** — the explicit statement *"Cycle #N closed.
Next: #M."* This is γ's last commit on the cycle's surfaces;
the closure-declaration commit's SHA on `main` is the structural
signal that the cycle has crossed from "released" to "closed." δ's
disconnect-release tag appearing on `main` after the closure-declaration
commit is the **observable proof** that all gate actions completed
and the cycle is fully closed.

A cycle that has stopped without closing (matter shipped; releases
tagged; but immediate outputs unexecuted, deferred outputs uncommitted,
or cycle-iteration absent under triggered conditions) is **structurally
in-flight** until the conjunction holds. The next cycle's γ inherits
the unfinished close-out work as immediate outputs of the next
cycle — but the deferred-output cost compounds: the lag-table grows;
the next-MCA handoff is unanchored; ε's receipt-stream review surfaces
the closure gap as a `cds-protocol-gap` finding.

### Operational realization

The immediate outputs, deferred outputs, and closure rule above are
CDS's canonical closure statement. The v0.1 operational overlay — the
per-output execution mechanics, the closure-gate row-by-row check,
the closure-declaration authoring procedure, the next-MCA filing
workflow, the hub-memory write sequence — lives in the existing cdd
role/runtime skills as the temporary v0.1 overlay until the v1 CDS-side
role rewrite:

- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.10`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — γ's closure gate (the 14-row row-by-row check that realizes the
  §Closure rule conjunction); the closure-declaration commit
  ("Cycle #N closed. Next: #M.") that signals structural closure.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.7`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — close-out triage; the four-disposition shape that realizes the
  "every finding gets a disposition" precondition of §Closure rule
  conjunct 1.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.8`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — cycle-iteration-trigger enforcement; the three-state shape (patch
  landed / next MCA / no-patch) that realizes §Closure rule conjunct 3.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §3.6`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — "Land immediate process fixes in the same cycle when possible";
  the discipline that makes §Immediate outputs structural rather than
  optional.
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §3.7`](../../../cnos.cdd/skills/cdd/gamma/SKILL.md)
  — "Do not close the cycle with unresolved triage"; the §Closure
  rule conjunct 1 enforcement at gate-check time (also fires §Gate
  F10).
- [`cnos.cdd/skills/cdd/post-release/SKILL.md §6`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — Decide next move (next-MCA issue / owner / branch / first AC /
  MCI freeze); the §Deferred outputs authoring procedure.
- [`cnos.cdd/skills/cdd/post-release/SKILL.md §7`](../../../cnos.cdd/skills/cdd/post-release/SKILL.md)
  — Hub memory write (daily reflection + adhoc thread); the
  immediate-output execution that prevents the compaction-gap
  failure mode.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{gamma,post-release,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Retro-packaging

**Retro-packaging** is the exception path for matter that landed
directly on `main` outside the canonical cycle/{N} branch model.
The direct-to-main exception covers the case where a change shipped
without a cycle branch (typically: emergency patches; tiny
operator-class fixes; matter that landed before the canonical
cycle/{N} model was in force). The exception's structural cost is
**lost cycle-shape evidence** — there is no cycle directory, no
self-coherence, no β review record at the canonical paths. Retro-packaging
re-creates the shape after the fact, preserving the §Frozen snapshot
rule while acknowledging the structural gap.

Retro-packaging is **rare** (default rule: every substantial change
opens a cycle/{N} branch per §Development lifecycle → §Branch rule),
and the structural cost is recorded explicitly — the retro-packaged
cycle's PRA carries a `<C` configuration-floor cap on γ (per
`release/SKILL.md §3.8` v0.1 overlay) reflecting the absence of
cycle-shape evidence. Retro-packaging is not a way to escape the
cycle/{N} discipline; it is a way to honour the §Frozen snapshot
rule for matter that already shipped.

### Direct-to-main exception

When matter landed on `main` without a cycle branch, γ executes
retro-packaging as a docs-only release per `release/SKILL.md §2.5b`
v0.1 overlay. The retro-packaged cycle creates three artifacts that
re-establish the cycle-shape evidence:

- **Retro-snapshot in version directory** — γ creates the
  `docs/{tier}/{bundle}/{X.Y.Z}/` version directory (or
  `docs/{tier}/{bundle}/docs/{ISO-date}/` for docs-only releases per
  the `release/SKILL.md §2.5b` form). The snapshot directory carries
  a `README.md` snapshot manifest naming the deliverables and the
  retro-packaging context (the merge commit SHA of the matter being
  packaged; the rationale for the direct-to-main path).
- **Self-coherence covering the change** — γ authors a retroactive
  `self-coherence.md` at `.cdd/releases/{X.Y.Z}/{N}/self-coherence.md`
  (or `.cdd/releases/docs/{ISO-date}/{N}/self-coherence.md` for
  docs-only). The retroactive self-coherence is **structurally
  honest**: it names the gap (the matter shipped without a cycle
  branch); records the ACs the matter satisfies (with evidence binding
  to the merge commit); names the contract the matter binds against
  retroactively. The `[retroactive]` marker is mandatory in the
  self-coherence's §Mode section.
- **Version-history entry** — γ adds a CHANGELOG TSC ledger row at
  the canonical bare-version format (per `release/SKILL.md §2.4`
  v0.1 overlay) with honest scoring. The γ-axis score reflects the
  retro-packaging cost (configuration-floor cap at C+ per
  `release/SKILL.md §3.8`'s "score the release, not the intent"
  discipline; the canonical example is `usurobor/tsc` cycle 27 at
  α B / β C+ / γ C / C_Σ C+). The coherence note records the
  retro-packaging rationale.

The retro-packaged cycle's PRA carries an explicit §Retro-packaging
note in §3 Process Learning naming **why** the direct-to-main path was
taken (operator-class emergency; pre-cycle-model historical matter;
operational-substrate fix) and **what** prevents recurrence (a tooling
patch that catches direct-to-main commits; an operator-skill patch
that names the cycle/{N} discipline; an explicit acceptance that the
direct-to-main path is allowed for the specific class of matter).

### Operational realization

The direct-to-main exception above is CDS's canonical retro-packaging
statement. The v0.1 operational overlay — the docs-only disconnect
mechanics, the retroactive `self-coherence.md` authoring procedure,
the CHANGELOG honest-grading workflow under retro conditions, the
configuration-floor cap mechanics — lives in the existing cdd
role/runtime skills as the temporary v0.1 overlay until the v1
CDS-side role rewrite:

- [`cnos.cdd/skills/cdd/release/SKILL.md §2.5b`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — docs-only disconnect (no tag); the operational form that retro-packaged
  cycles use (`.cdd/releases/docs/{ISO-date}/{N}/` cycle-dir move;
  `docs/gamma/cdd/docs/{ISO-date}/POST-RELEASE-ASSESSMENT.md` PRA
  path).
- [`cnos.cdd/skills/cdd/release/SKILL.md §3.8`](../../../cnos.cdd/skills/cdd/release/SKILL.md)
  — configuration-floor clause; the C+ / C / `<C` cap mechanics that
  apply to retro-packaged cycles whose γ-axis lacks cycle-shape
  evidence; the "score the release, not the intent" discipline that
  ensures the retro-packaging cost is recorded honestly.

When the v1 CDS-side role overlays land (post-#403 wave), the operational
realization moves into `cnos.cds/skills/cds/{release,gamma,…}/SKILL.md`;
the v0.1 cdd cite-points above are the temporary v0.1 home.

---

## Large-file authoring rule

CDS authoring discipline imposes a **section-by-section write rule**
for files larger than 50 lines. The rule's failure mode is the
session-end / stream-timeout / context-truncation gap: a generation
that produces a single 1000-line write loses everything if the session
terminates mid-stream. Section-by-section writing with a manifest
header lets the next session resume from the last committed section
without rebuilding the entire file in memory.

The rule is self-referential — `CDS.md` itself (this document) follows
the rule, with the section-manifest header at lines 1–2 carrying both
the `sections:` list (the canonical section ordering) and the
`completed:` list (which sections have been fully written to disk).
The Sub 5 (this cycle) write process incrementally appended the 8 new
top-level sections per the resumption protocol below; the post-Sub-5
header reflects the full canonical section set.

### File-size threshold

Files **larger than 50 lines** of authored content (excluding the
manifest header and trailing whitespace) trigger the section-by-section
rule. Below the threshold, single-write authoring is acceptable; above
it, the structural cost of single-write authoring (lost work on
session termination; context-window overflow during generation;
review-comment locality failures) outweighs the convenience.

The 50-line floor is the canonical empirical threshold (the cnos
cycle history shows session-termination at 60–80 line generations as
a recurring `cds-tooling-gap`). Project bindings may impose a stricter
floor (some downstream repos use 30 lines for safety); the floor is
**only** raised via project-binding policy, never lowered.

### Section-manifest HTML-comment header

Files subject to the rule open with two HTML comment lines:

```html
<!-- sections: [section-id-1, section-id-2, section-id-3, …] -->
<!-- completed: [section-id-1, section-id-2] -->
```

The two lists are typed pairs:

- **`sections:`** — the canonical ordering of section identifiers
  the file commits to author. The identifiers are slugs (kebab-case,
  lowercase) of the section headings, not the heading text itself.
  The list is the file's structural contract: a section that does
  not appear in `sections:` cannot be authored under the file's
  contract; a section in `sections:` that has no corresponding
  `## ` heading in the file's body is incomplete.

- **`completed:`** — the prefix of `sections:` that has been fully
  written to disk. A section appears in `completed:` only when its
  prose-body is finished and the next section's heading exists
  immediately after. Partial sections do not appear in `completed:`.

The two lists' relationship is structural: `completed:` is a prefix
of `sections:` for files under canonical write; `completed:` may be
empty (no sections completed yet); `completed:` may equal `sections:`
(file is complete). A `completed:` list that is not a prefix of
`sections:` signals a manifest-protocol violation (e.g. a section
landed out of order; the file was edited outside the canonical
write process).

### Resumption protocol

When a session resumes work on a file under the rule, the resumption
protocol below runs as α's first action on the file:

1. **Read the manifest header** (lines 1–2). Parse `sections:` and
   `completed:`.
2. **Compute the next section.** The next section to author is the
   first element of `sections:` that is not in `completed:`.
3. **Verify the file body** — every section identifier in `completed:`
   has a corresponding `## ` heading in the file body, in the order
   given by `sections:`. If any verification fails, surface the
   manifest-protocol violation before authoring the next section.
4. **Author the next section.** Write the section's prose-body
   immediately before the next-section anchor (or at the file's
   end if the next section is the last in `sections:`).
5. **Commit the section.** Each section gets its own commit (per the
   incremental-write discipline; cnos cycles routinely show
   per-section α commits on `cycle/{N}` branches). The commit
   message names the section authored.
6. **Update `completed:`** in the manifest header. Append the
   just-authored section identifier to the `completed:` list. The
   manifest update is part of the section's commit.

A session that terminates mid-section (the section's prose-body is
partially written; the section identifier does not yet appear in
`completed:`) is recoverable: the next session reads the manifest,
computes the next section (which is the partially-written one,
because it does not appear in `completed:`), reads what exists on the
branch for that section, and continues authoring from the last
written paragraph. The branch becomes the recovery surface; the
manifest is the index.

### Anti-patterns

- **Single-write authoring of a 200-line section.** The session-end
  failure mode dominates; one terminated session loses the whole
  section. The fix is splitting the section into smaller pieces or
  committing partial writes (with the partial state visible in a
  `<!-- WIP: paragraph 3 of 7 -->` note inside the section's body).
- **Out-of-order section authoring.** `sections:` declares the
  canonical ordering; authoring section 5 before section 3 violates
  the manifest's prefix invariant on `completed:`. The fix is
  honouring the declared order, or amending the manifest to declare
  the new ordering before authoring.
- **Manifest header omitted.** A file > 50 lines without a manifest
  header is structurally not under the rule's protection — the
  next session has no resumption index. The fix is adding the
  manifest header in the next commit; the file's prior commit
  history reconstructs the inferred `completed:` list.
- **Mismatched `sections:` and body headings.** A `sections:` entry
  with no corresponding `## ` heading in the body, or a `## `
  heading with no `sections:` entry, signals manifest drift. The
  fix is reconciling the two; the reconciliation is its own commit.

### Operational realization

The 50-line threshold, the section-manifest header, the resumption
protocol, and the anti-patterns above are CDS's canonical large-file
authoring rule. The rule is **self-referential**: its operational
realization is its own usage pattern. Every file in `cnos.cds`,
`cnos.cdd`, and `cnos.cdr` that exceeds the 50-line threshold follows
the rule; this `CDS.md` document is the canonical exemplar (Sub 2,
Sub 3, Sub 4, and Sub 5 of cnos#403 incrementally authored the
sections via the resumption protocol, with each sub's manifest update
landing in the sub's α commits).

The v0.1 operational overlay — the per-file resumption-protocol
script, the manifest-validator that enforces the prefix invariant,
the integration with γ's branch pre-flight that checks for manifest
drift before α dispatch — does **not** yet exist as a separate skill
file; the rule lives in this section's body as the canonical
statement. A future cycle may extract the operational overlay into
`cnos.cdd/skills/cdd/large-file/SKILL.md` (or the CDS-side equivalent)
once the rule has accumulated enough cycle history to warrant a
dedicated skill; until then, the rule's authority is this section.

The self-referential operational-realization pointer is therefore
this section itself — the file you are reading is the operational
exemplar. Cross-references from other skills (when they cite the
rule) cite `CDS.md §Large-file authoring rule`; the citation resolves
to this section's body, which is both the canonical statement and
the operational realization.

---

## Empirical anchor

CDS's empirical anchor is [`usurobor/cnos`](https://github.com/usurobor/cnos)
— this repository itself. Unlike CDR (whose empirical anchor `usurobor/cph`
is an external research repo), CDS's primary empirical anchor is the cnos
repo itself: every cnos cycle since cnos#364 has been a CDS cycle in
practice, executing the engineering discipline (artifact improvement under
repairable feedback) on the engineering substrate (the cnos codebase, its
tests, its build tooling, its CI, its release effectors).

The arrangement is mildly recursive: the CDS protocol is authored and
maintained inside the project that is also CDS's primary empirical anchor.
This is intentional. The CDS doctrine and the engineering work happen on
the same substrate so the doctrine learns from the work and the work binds
against the doctrine in tight feedback loops. The recursion is bounded
by the
[`ROLES.md §4a`](../../../../../ROLES.md#4a-persona-operator-protocol-project-receipt--the-five-layer-enforcement-chain)
five-layer chain: the protocol overlay (this file) is doctrinally
independent of the project binding (cnos's `.cdd/` directory) even when
they live in the same repository.

### Shape-compatibility claim

The six-field shape declared above describes cnos's CDS practice without
contradiction. Spot-checks against the empirical cycle wave:

- **Field 1 (Matter type)** — cnos cycles produce source code (the `cn`
  binary at `src/go/cmd/cn/`, the `cnos.cdd` doctrine, the `cnos.cdr`
  protocol overlay, this very file as a `cnos.cds` overlay), tests (the
  unit/integration suites under `src/go/`, the CUE fixture suites under
  `schemas/*/fixtures/`), documentation (every SKILL.md, every README,
  every doctrine file), schemas (`schemas/cdd/`, `schemas/cds/`,
  `schemas/cdr/`), CI definitions (the GitHub Actions workflows), runtime
  contracts (`cn.package.json` per package, the package-discovery
  conventions), and design notes (`docs/papers/`, the
  `cdd/design/SKILL.md` artifacts produced in cycle-authoring). The
  matter-type taxonomy decomposes the cnos-cycle output set without loss.
- **Field 2 (Review oracle)** — cnos cycles use compilation +
  type-check + build (`go run ./src/go/cmd/cn build --check`,
  `cue vet`, schema validators), tests pass (per the per-cycle test
  surface declared in each cycle's contract), AC verification
  (mechanical grep/wc/ls per #393 contract discipline; read-checks for
  prose-judgment ACs), no regressions (the full cn build + all
  schema-fixture validations remain green), implementation-contract
  coherence (the cnos#393 Rule 7 enforcement that landed in #393's
  wave), and evidence-binding (the close-out artifact set per CDD
  §5.5b, ratcheting toward the typed `#CDSReceipt` schema). The oracle
  taxonomy describes cnos β practice without contradiction.
- **Field 3 (γ close-out artifact)** — cnos cycles produce the
  six-artifact close-out set (`self-coherence.md`, `beta-review.md`,
  `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`,
  `cdd-iteration.md`) under `.cdd/unreleased/{N}/`, plus the
  forthcoming typed `#CDSReceipt` per `schemas/cds/receipt.cue` once
  Phase 3 of cnos#366 wires V into the release flow. The Field 3
  declaration above matches the per-cycle artifact convention cnos has
  used since #335 (the cycle that initialised the iteration aggregator).
- **Field 4 (δ cadence)** — cnos cycles are per-cycle (one cycle = one
  branch = one issue close), with release-shaped bundling (PRA per
  release; RELEASE.md per release; release tags on main; release-effector
  mechanics in `release-effector/SKILL.md`), and wave-shaped grouping for
  multi-cycle waves (cnos#366 phases; cnos#376 cdr wave; cnos#403 cds
  wave). The cadence taxonomy describes cnos δ practice without
  contradiction. The inward-membrane discipline (cnos#393 δ-as-architect)
  is the canonical example: every recent cycle's issue body carries a
  pinned implementation contract.
- **Field 5 (ε trigger classes)** — cnos cycles' recurring protocol-gap
  signals match the trigger classes named in Field 5. The
  `.cdd/iterations/INDEX.md` aggregator (37 rows as of cycle/406) carries
  the receipt-stream view ε reads; the `cdd-*-gap` markers are the
  pre-rename forms of the `cds-*-gap` classes declared in Field 5; the
  per-cycle finding counts and patch counts in the aggregator are exactly
  the longitudinal view ε's discipline produces.
- **Field 6 (Actor collapse rule)** — cnos cycles' actor configurations
  match the Field 6 collapse table. The breadth-2026-05-12 wave manifest
  introduced β-α-collapse-on-δ for skill/docs-class cycles; cycles #388
  through #406 (and this cycle #407) have dispatched under that pattern;
  the γ=δ and ε=δ collapses are the cnos default for small-protocol
  regimes (operator-as-coordinator pattern). No cnos cycle has used
  α=β within a single substantive-code cycle.

The compatibility claim is sufficient for Sub 2's scope: the six-field
shape can host cnos's CDS practice without re-derivation. The
surface-by-surface mapping — every cnos engineering surface to the
corresponding CDS Field — is Sub 7 territory (see below).

### Representative cycle milestones

A representative subset of the cnos cycle wave that anchors CDS:

- **cnos#364** (closed) — `COHERENCE-CELL.md` doctrine landed at merge
  `32b126e4`. The point where cnos's doctrine surface began to be
  authored as a recursive cell model rather than a flat
  α/β/γ/δ/ε description.
- **cnos#366** (in-flight) — coherence-cell executability roadmap. The
  multi-phase wave that wired the CCNF kernel through to a working V
  predicate and ratcheted the doctrine through Phase 7 (the CDD.md
  rewrite that compressed CDD.md to the CCNF spine and marked the
  software-realization content as pending cds extraction).
- **cnos#376** — `cnos.cdr` v0.1 wave. The structural precedent for
  cnos#403: the same option-(a) split applied to skills/role-overlays
  for the research realization.
- **cnos#388** — Phase 2.5 of cnos#366; schemas generic/domain split.
  The architectural-choice decision-record CDS inherits.
- **cnos#393** — δ-as-architect; introduction of Rule 7
  (implementation-contract coherence) as a β oracle. The cycle that
  pinned the inward-membrane discipline cited in Field 4.
- **cnos#402** — CCNF spine rewrite; CDD.md compressed to the kernel
  layer; the "Software-specific realization — pending cds extraction"
  section that quarantined the content cnos#403 extracts.
- **cnos#403** — cnos.cds bootstrap tracker; the parent of this cycle.
  Sub 1 (cnos#406) shipped the package skeleton + extraction map at
  merge `987acd04`; Sub 2 (cnos#407) ships this file.

A fuller anchor — every cnos cycle's per-Field mapping; the per-AC
coverage matrix; the gap-class distribution across the cycle wave — is
Sub 7 territory.

### Sub 7 deferred surface-by-surface mapping

Sub 7 of [cnos#403](https://github.com/usurobor/cnos/issues/403) owns the
full empirical-anchor doc:
[`docs/empirical-anchor-cdd.md`](../../docs/empirical-anchor-cdd.md) (to
be authored). The doc will perform a surface-by-surface mapping of cnos's
`.cdd/` cycle-artifact set to CDS's six-field structure and to
`#CDSReceipt`'s typed keys, mirroring
[`cnos.cdr/docs/empirical-anchor-cph.md`](../../../cnos.cdr/docs/empirical-anchor-cph.md)'s
structural precedent. Sub 2 (this cycle) establishes the citation and the
shape-compatibility claim; Sub 7 verifies the compatibility on every
cnos cycle-artifact and records the mapping table.

The artifact's name (`empirical-anchor-cdd.md` rather than
`empirical-anchor-cnos.md`) reflects the historical anchor: the cycles
being mapped were authored as `cdd-*` artifacts under `.cdd/`. Post-
filesystem-rename (the open coordination question in `docs/extraction-map.md
§14`), the anchor doc may rename to `empirical-anchor-cnos.md` or
similar; the v0.1 naming follows the cycle/406 README precedent.

---

## Related documents

Inherits, cites, or extends:

- [`ROLES.md`](../../../../../ROLES.md) — the role-cell grammar and the
  six-field instantiation contract that CDS realises. Specifically:
  - `§3` — six-field instantiation contract (the structure mirrored
    above).
  - `§4a` — five-layer enforcement chain (persona/protocol/project
    cited above).
  - `§4a.2` — loss-function distinction (CDS's engineering discipline).
  - `§4a.3` — receipt-enforces-discipline (the CDS receipt sketch
    realised in the schema).
  - `§4a.4` — Sigma persona exemplar.
  - `§4b` — generic ε doctrine (gap class taxonomy pattern, iteration
    artifact, cadence rule, MCA discipline, collapse rule).
  - `§7` — naming convention reserving `cds` for software.
- [`schemas/cds/receipt.cue`](../../../../../schemas/cds/receipt.cue) — the
  typed γ close-out surface; `#CDSReceipt` is the parent-facing artifact
  CDS Field 3 names.
- [`schemas/cdd/README.md §"Architectural choice"`](../../../../../schemas/cdd/README.md#architectural-choice--generic-vs-domain-split)
  — the upstream decision-record for the (a) split inherited above.
- [`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../../cnos.cdd/skills/cdd/CDD.md)
  — the CCNF kernel doctrine. CDS instantiates this kernel; the kernel
  is cited by reference, not restated. The "Software-specific
  realization — pending cds extraction" section of CDD.md is the
  source-of-truth for the lifecycle/artifacts/review/gate content
  migrating to CDS in Subs 3–5.
- [`src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`](../../../cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md)
  — the substrate-independent kernel algorithm (five-step recursion,
  four outcomes, two recursion modes, three scope-lift projections).
  Cited by Field 2 (evidence-binding oracle) and Field 4 (boundary
  decision per step 5).
- [`src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md`](../../../cnos.cdd/skills/cdd/COHERENCE-CELL.md)
  — the predecessor doctrine; the four-way structural separation (role
  / runtime substrate / validation / boundary effection) that CDS's
  role overlays (Sub 3) realise on the engineering substrate.
- [`src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`](../../../cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md)
  — V's design surface; the parent-facing validator interface CDS's
  Field 3 typed receipt feeds.
- [`src/packages/cnos.cdr/skills/cdr/CDR.md`](../../../cnos.cdr/skills/cdr/CDR.md)
  — the research realization. Structural sibling of this file: same
  section structure; same six-field contract shape; divergent discipline
  profile per `ROLES.md §4a.2`. CDS.md cites CDR.md as the parallel
  record at §"Architecture choice" and §"Six-field instantiation
  contract" (where each field's CDS-side declaration is paired against
  CDR's research-side declaration as the language contrast).
- [`src/packages/cnos.cds/docs/extraction-map.md`](../../docs/extraction-map.md)
  — the Sub 1 deliverable; the surface-by-surface migration plan that
  Subs 3–5 dispatch against. CDS.md is the *contract*; the extraction
  map is the *migration plan that lands operational content against the
  contract's surfaces*.
- [`src/packages/cnos.cds/README.md`](../../README.md) — the package
  README authored at Sub 1; declares CDS as the software-development
  realization of CCNF; cites kernel docs; cites CDR as sibling.
- [`src/packages/cnos.cds/skills/cds/SKILL.md`](SKILL.md) — the loader
  skill; loads CDS.md as Step 1 of its Load order (per Sub 2's D2
  edit); names the role-overlay positions as advisory targets until
  Subs 3–5 land them.
- [cnos#403](https://github.com/usurobor/cnos/issues/403) — parent
  tracker; this file is Sub 2's deliverable.
- [cnos#406](https://github.com/usurobor/cnos/issues/406) — Sub 1;
  package skeleton + extraction map (landed at merge `987acd04`).
- [cnos#407](https://github.com/usurobor/cnos/issues/407) — Sub 2;
  this cycle.
- [cnos#388](https://github.com/usurobor/cnos/issues/388) — schema-side
  architectural-choice precedent.
- [cnos#376](https://github.com/usurobor/cnos/issues/376) — research-side
  skills application (structural precedent for cnos#403).
- [cnos#393](https://github.com/usurobor/cnos/issues/393) — δ-as-architect
  / Rule 7; the implementation-contract-coherence β oracle cited in
  Field 2 and the inward-membrane discipline cited in Field 4.
- [cnos#401](https://github.com/usurobor/cnos/issues/401) — the
  cdd-iteration cadence rule (required only when `protocol_gap_count > 0`;
  courtesy stub otherwise) cited in Field 3 and Field 5.
- [cnos#402](https://github.com/usurobor/cnos/issues/402) — CCNF spine
  rewrite; the cycle that quarantined the software-realization content
  this wave (cnos#403) extracts.

---

## Non-goals

CDS's non-goals partition into two scopes: **sub-level non-goals** that
scope the cnos#403 wave's authoring sub-cycles (Subs 2–5; this document's
authoring discipline) and **software-cycle non-goals** that apply to
every CDS cycle as protocol-level doctrine (the engineering-loss-function
discipline's structural anti-patterns).

### Sub-level non-goals (Subs 2–5 authoring scope discipline)

This document does **not**:

- Migrate any software-lifecycle content from CDD.md — that is cnos#403
  Subs 3–5 (Selection, Lifecycle, Artifacts, Review, Gate, Assessment,
  Closure, Retro-packaging, Non-goals, Large-file per the extraction
  map's 12 surface-group tables).
- Modify CDD.md or remove the "pending cds extraction" markers from
  CDD.md — that is cnos#403 Sub 6.
- Author role-overlay skills (`alpha/SKILL.md`, `beta/SKILL.md`,
  `gamma/SKILL.md`, `delta/SKILL.md`, `epsilon/SKILL.md` for CDS) —
  those are cnos#403 Subs 3–5.
- Author the empirical-anchor mapping doc
  (`docs/empirical-anchor-cdd.md`) — that is cnos#403 Sub 7.
- Modify `cnos.cdr` (the research realization) or `cnos.cdd` (the
  kernel). Hard rule per the cnos#407 implementation contract; the AC7
  and AC8 mechanical checks verify `git diff` empty for both packages.
- Modify `schemas/cds/` or any other schema package. The typed receipt
  `#CDSReceipt` was authored at cnos#388 (Phase 2.5); CDS.md cites the
  schema once and does not edit it.
- Redefine `ROLES.md` or the role ladder. CDS instantiates; ROLES
  governs.
- Restate any CCNF kernel doctrine. CDD.md, COHERENCE-CELL.md,
  COHERENCE-CELL-NORMAL-FORM.md, and RECEIPT-VALIDATION.md own the
  kernel; CDS.md cites them. The pointer discipline is the one source
  of truth principle from `design/SKILL.md §3.2`.
- Author persona-identity content for `cn-sigma`. Persona drafts live
  in the persona hub (`cn-sigma/spec/PERSONA.md`), not here. The
  persona hub itself is forthcoming.
- Author project-binding content for cnos. Project bindings live in
  `<project>/.cds/` (or `.cdd/` until the rename); CDS.md declares the
  doctrinal shape, not project-specific gate names or CI thresholds.
- Specify runtime mechanics (dispatch, polling, git wiring, CI hooks,
  release-effector tagging, deployment probes). Those belong in
  role-overlay skills (Sub 3), runtime-substrate skills
  (`harness/SKILL.md`, `release-effector/SKILL.md` — currently in
  `cnos.cdd`; relocation considered as part of Sub 5 per
  `docs/extraction-map.md §14`), and the persona/operator contracts.
- Author a `cn-cds-verify` command or any tooling. The validator
  interface (V) is `cnos.cdd`'s scope, applied to `#CDSReceipt` via
  `protocol_id` dispatch per
  [`schemas/cdd/README.md §"protocol_id dispatch convention"`](../../../../../schemas/cdd/README.md#protocol_id-dispatch-convention).
- Author a `cdr-vs-cds-vs-cdw-vs-cda` comparison table. CDS is a sibling
  to CDR, not the comparison-curator; CDR.md and CDS.md each cite the
  other as a sibling, but the comparison-curating job belongs in
  generic-substrate documentation (`ROLES.md §4a.2` already names the
  loss-function distinction).

### Software-cycle non-goals

CDS as a protocol does **not**:

- Do NOT optimize primarily for speed. The engineering loss function is
  artifact improvement under repairable feedback, not throughput; a
  cycle that ships fast but ships incoherent matter has not improved
  the artifact. The §Selection function's effort-adjusted tie-break
  prefers smaller cycles among candidates of equal leverage and axis
  effect, but speed is a tie-break term, not a primary objective.
- Do NOT treat issue queues as self-justifying. An open issue does not
  prove a coherence gap exists; the §Selection function's stale-backlog
  re-evaluation rule explicitly requires re-evaluating whether stale
  issues remain real gaps before selecting them as next-MCAs. A queue
  whose age is its own justification is structural debt.
- Do NOT reduce review to local diff reading. β's review CLP names
  TERMS / POINTER / EXIT scopes that extend beyond the diff (the issue
  body's ACs, the implementation-contract axes pinned by δ, the
  loaded-skill set, cross-referenced canonical surfaces). A review that
  reads only the diff misses the contract-coherence and honest-claim
  classes the broader CLP catches.
- Do NOT treat release as "tag and hope". The §Gate's release-readiness
  preconditions are mechanical preconditions, not aspirational targets;
  the F1–F10 closure verification checklist catches the failure modes
  the "tag and hope" pattern produces. A release that skips the gate
  ships unreleased in the CCNF kernel sense — the boundary effection
  did not happen even though the tag pushed.
- Do NOT confuse a shipped feature with a closed coherence cycle. A
  feature that merged into main is structurally **stopped**, not
  **closed**, until the §Closure rule's three-conjunct condition holds
  (immediate outputs executed; deferred outputs committed; cycle
  iteration present if triggered). Closure is the transition from
  released to closed; conflating "released" with "closed" is the
  structural failure mode the §Closure rule prevents.
