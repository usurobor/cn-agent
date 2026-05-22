<!-- sections: [preamble, architecture-choice, persona-protocol-project, six-field-contract, empirical-anchor, related-documents, non-goals] -->
<!-- completed: [preamble, architecture-choice, persona-protocol-project, six-field-contract, empirical-anchor, related-documents, non-goals] -->

# Coherence-Driven Software (CDS)

> CDS is a software-engineering-protocol instantiation of the generic
> role-scope ladder pattern. The pattern (α/β/γ/δ/ε roles, scope-escalation
> contract, six-field instantiation contract) is documented at
> [`ROLES.md`](../../../../../ROLES.md) at the repo root. CDS does not claim
> to have originated the role structure; it instantiates the structure for
> software-engineering work — code, tests, documentation, schemas, releases,
> deployments — under the engineering loss function (artifact improvement
> under repairable feedback).

**Version:** 0.1 (Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403))
**Status:** Draft — instantiation contract only. Lifecycle migration (Subs
3–5), CDD.md marker cleanup (Sub 6), and empirical-anchor doc (Sub 7) are
downstream.
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
- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol,
  and project"`](../../../../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md)
  — the doctrinal source for treating persona, protocol, and project as
  three distinct layers.
- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Wave 5 — CDS extraction by
  reference"` and `§"Wave 7 — final CDD.md rewrite"`](../../../../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md)
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
  `docs/gamma/essays/`. Design notes are α-matter when they pin a
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
  conventions), and design notes (`docs/gamma/essays/`, the
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
