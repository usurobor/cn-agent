# CDD — Recursive Coherence Cell Protocol

**Version:** 4.0.0
**Status:** Draft (Phase 7 of [#366](https://github.com/usurobor/cnos/issues/366) — CCNF spine; terminal phase)
**Placement:** `src/packages/cnos.cdd/skills/cdd/`

CDD owns the **generic recursive coherence-cell algorithm.** This document is the algorithm's canonical statement: it names the kernel (CCNF), the four cell outcomes, the two recursion modes, the three scope-lift projections, and the domain packages that bind the kernel to substrates. CDD does not own substrate-specific realization — software-lifecycle realization lives in [`cnos.cds`](../../../cnos.cds/skills/cds/CDS.md) (v0.1 shipped per the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave); research-lifecycle realization lives in [`cnos.cdr`](../../../cnos.cdr/).

The kernel is stated verbatim from [`COHERENCE-CELL-NORMAL-FORM.md`](COHERENCE-CELL-NORMAL-FORM.md) (CCNF). That document is the citable source of the recursion equation; this document is the operational doctrine that names domain bindings and points at canonical surfaces. Software-specific realization has been extracted to [`cnos.cds`](../../../cnos.cds/skills/cds/CDS.md) per cnos#403 — see §"Software-specific realization → cnos.cds" below.

---

## Kernel (CCNF)

The coherence cell at scope `n` is a **five-step closed loop.** At scope `n`, a `contractₙ` is given. The kernel composes five steps over that contract:

```text
1.  matterₙ      := αₙ.produce(contractₙ)
2.  reviewₙ      := βₙ.review(contractₙ, matterₙ)
3.  receiptₙ     := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)
4.  verdictₙ     := V(contractₙ, receiptₙ)
5.  decisionₙ    := δₙ.decide(receiptₙ, verdictₙ)
```

Each step is a typed function with a fixed signature (see `COHERENCE-CELL-NORMAL-FORM.md §Kernel`). The five steps compose into a closed cell at scope `n`:

```text
closed_cellₙ := { contractₙ, matterₙ, reviewₙ, receiptₙ, verdictₙ, decisionₙ }
```

The **evidence-binding rule** is load-bearing: *evidence accumulates during α and β work; γ binds it into the receipt at close-out as typed references; V dereferences the references to validate; β never consumes evidence directly; δ never re-reads evidence.* The rule has four enforcement points — α's signature excludes the receipt, β's signature excludes evidence, γ's signature includes evidence as the input γ binds, and δ's signature excludes evidence.

V is a typed predicate, not a role. δ holds gate authority over what crosses the boundary; V emits a verdict δ trusts. If δ doubts the verdict, δ has two paths: re-invoke V or repair-dispatch the cell. δ does not read evidence directly.

The kernel sections of `COHERENCE-CELL-NORMAL-FORM.md` are **substrate-independent.** They name only roles (α, β, γ, δ, ε), the validator predicate (V), the artifacts the kernel reasons over (`contract`, `matter`, `review`, `receipt`, `verdict`, `decision`, `evidence`), the scope index (`n`, `n+1`), and the verdicts and decisions the algorithm composes. They do not name any particular tooling, platform, or invocation surface.

## Outcomes

A closed cell at scope `n` terminates in **exactly one of four outcomes,** determined by `(verdictₙ, decisionₙ)`:

```text
accepted := (verdictₙ = PASS)  ∧ (decisionₙ ∈ {accept, release})
degraded := (verdictₙ ≠ PASS)  ∧ (decisionₙ = override)
blocked  := (decisionₙ ∈ {reject, repair_dispatch})        — any verdict
invalid  := (verdictₙ = PASS   ∧ decisionₙ = override)
          ∨ (verdictₙ ≠ PASS   ∧ decisionₙ ∈ {accept, release})
```

`accepted` is the clean exit — the closed cell projects as α-matter at scope `n+1`. `degraded` is the override exit — the cell projects under explicit override; the override block is the structural signal every downstream consumer must detect. `blocked` is the no-projection exit — under `reject` the cell terminates at scope `n`; under `repair_dispatch` the cell stays open at scope `n` and a child cell runs under a repair contract. `invalid` is **non-terminal** — δ must re-decide until the `(verdict, decision)` pair satisfies one of the three terminal preconditions. See `COHERENCE-CELL-NORMAL-FORM.md §Cell Outcomes` for the full statement.

## Recursion modes

The kernel's recursion has **two modes,** distinguished by `decisionₙ`:

```text
decisionₙ = repair_dispatch          → within-scope: same n, re-emit, re-fire
decisionₙ ∈ {accept, release}        → cross-scope: closed_cellₙ → α-matter at n+1
decisionₙ = override                 → cross-scope: closed_cellₙ → α-matter at n+1 (degraded)
decisionₙ = reject                   → cross-scope: closed_cellₙ terminates at n, no projection
```

Both modes share the same recursion operator (the five-step closure); they differ in how the scope index advances. Within-scope repair-dispatch keeps the cell open at scope `n` with a child cell running under a derived repair contract; cross-scope accept/degraded projects the closed cell to scope `n+1`; cross-scope reject closes at scope `n` with no projection. See `COHERENCE-CELL-NORMAL-FORM.md §Recursion Modes` for the full statement.

## Scope-lift

The cross-scope mode advances the scope index from `n` to `n+1` via **three projections under scope-lift:**

```text
Projection 1:  closed (αₙ, βₙ, γₙ) cell        →   αₙ₊₁ matter
Projection 2:  δₙ boundary decision             →   βₙ₊₁-like discrimination
Projection 3:  εₙ receipt-stream observation    →   γₙ₊₁-like coordination / evolution
```

The framing is **projection-not-renaming.** δₙ is not literally βₙ₊₁; εₙ is not literally γₙ₊₁. The scope-`n+1` cell has its own full α, β, γ, δ, ε at scope `n+1`; what projects from scope `n` is the closed cell (carrying βₙ's review and γₙ's receipt as components), δₙ's decision, and εₙ's receipt-stream observation. βₙ and γₙ have no separately typed upward projection — their work is intra-cell and is carried by the closed cell itself. See `COHERENCE-CELL-NORMAL-FORM.md §Scope-Lift` for the full statement.

## Domain packages

The CDD kernel binds to specific substrates through **domain packages** that extend the generic algorithm. Each domain package owns the substrate-specific realization of the kernel — the artifacts the cell produces, the validators the parent scope invokes, the boundary decisions δ records — while the kernel itself remains substrate-independent.

- **CDS — `cnos.cds`** — software-development realization. CDS binds CCNF to source code, tests, releases, deployments. Receipt fixtures live at `schemas/cds/`; the CDS package v0.1 shipped under the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave (subs cnos#406–#412). Canonical doctrine: [`cnos.cds/skills/cds/CDS.md`](../../../cnos.cds/skills/cds/CDS.md). See §"Software-specific realization → cnos.cds" below for the per-section pointer list.
- **CDR — `cnos.cdr`** — research-domain realization. CDR binds CCNF to research artifacts (claims, evidence, replication, falsification). CDR v0.1 shipped at [cnos#376](https://github.com/usurobor/cnos/issues/376); receipts validate at `schemas/cdr/`.
- **Future c-d-X protocols** — additional domain bindings (e.g., a `cnos.cdo` operations protocol, a `cnos.cdh` hardware protocol). Each c-d-X protocol follows the same pattern: the kernel stays generic; the c-d-X domain package owns the substrate-specific evidence, validators, and decision shapes.

The binding pattern is: a CDS / CDR / c-d-X domain package cites the CDD kernel as input contract, defines the domain's typed receipt schema under `schemas/{cdd-X}/`, names the substrate-specific α/β/γ/δ/ε realizations, and consumes `V` (the generic validator) through that schema. CDS, CDR, and any future c-d-X protocol do not re-derive the kernel; they extend it.

## Pointers

The canonical surfaces CDD cites are grouped by what they own:

**Generic doctrine (the kernel and its companions):**
- [`COHERENCE-CELL.md`](COHERENCE-CELL.md) — predecessor doctrine; receipt rule; four-way structural separation
- [`COHERENCE-CELL-NORMAL-FORM.md`](COHERENCE-CELL-NORMAL-FORM.md) — CCNF; the recursion equation source for this document
- [`RECEIPT-VALIDATION.md`](RECEIPT-VALIDATION.md) — V's design surface; input refs, output verdict shape, invocation rule
- [`CELL-KINDS.md`](CELL-KINDS.md) — the cell-kind taxonomy: one kernel, multiple typed refinements (`issue_authoring`, `implementation`, `repair`, `recovery`, `wave`, `cleanup`, `audit`, `release`, `doctrine`, `experiment`), each naming its matter, review surface, and FSM implication ([cnos#570](https://github.com/usurobor/cnos/issues/570))
- `ROLES.md` §1, §3, §4, §4a, §4b — generic role doctrine (§4a: δ; §4b: ε protocol-iteration)

**Schemas (typed evidence by domain):**
- `schemas/cdd/` — generic kernel schemas (`contract.cue`, `receipt.cue`, `boundary_decision.cue`, `validation_verdict.schema.json`)
- `schemas/cds/` — software domain schemas (receipt + fixtures; package v0.1 shipped per [cnos#403](https://github.com/usurobor/cnos/issues/403))
- `schemas/cdr/` — research domain schemas (receipt + fixtures; v0.1 per [cnos#376](https://github.com/usurobor/cnos/issues/376))

**Roles (the substrate realizations of α/β/γ/δ/ε):**
- [`alpha/SKILL.md`](alpha/SKILL.md) — α (produces matter)
- [`beta/SKILL.md`](beta/SKILL.md) — β (reviews and merges; merge is β's authority)
- [`gamma/SKILL.md`](gamma/SKILL.md) — γ (coordinates and closes; emits receipt)
- [`delta/SKILL.md`](delta/SKILL.md) — δ (boundary; decides at the membrane)
- [`epsilon/SKILL.md`](epsilon/SKILL.md) — ε (cross-cell receipt-stream observer; generic doctrine at `ROLES.md §4b`)

**Runtime substrate (the operational surfaces that execute the protocol on this platform):**
- [`harness/SKILL.md`](harness/SKILL.md) — dispatch, polling, branch creation, session lifecycle
- [`release-effector/SKILL.md`](release-effector/SKILL.md) — tag/release/deploy mechanics δ invokes at the boundary
- [`operator/SKILL.md`](operator/SKILL.md) — δ-the-operator's session-routing and re-dispatch contract

**Realization peers (the domain packages that bind this kernel):**
- `src/packages/cnos.cds/` — software realization (v0.1 shipped per the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave; canonical doctrine at [`cnos.cds/skills/cds/CDS.md`](../../../cnos.cds/skills/cds/CDS.md))
- `src/packages/cnos.cdr/` — research realization (shipped per [cnos#376](https://github.com/usurobor/cnos/issues/376))

**Loader and rationale:**
- [`SKILL.md`](SKILL.md) — package-visible loader entrypoint (not a second fact source)
- `docs/development/cdd/RATIONALE.md` — companion rationale for shape decisions
- `docs/papers/CCNF-AND-TYPED-TRUST.md` — the essay that pins the CCNF spine and the kernel/realization separation

## Software-specific realization → cnos.cds

Software-cycle realization detail lives in [`cnos.cds/skills/cds/CDS.md`](../../../cnos.cds/skills/cds/CDS.md) as of the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave closure (Subs 2–5: cnos#407, cnos#408, cnos#409, cnos#410). CDS is the software-development realization peer of CCNF; `cnos.cdr` is the research realization peer. CDD.md retains the generic kernel; CDS owns the software-cycle surfaces. Citing skill files in this package (`alpha/`, `beta/`, `gamma/`, `operator/`, `release/`, `release-effector/`, `post-release/`, `harness/`, `activation/`, `review/`) cite CDS canonical homes directly per the pointer list below.

- [`cnos.cds/skills/cds/CDS.md §"Six-field instantiation contract"`](../../../cnos.cds/skills/cds/CDS.md) — matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule (sequential bounded dispatch lives under Field 6)
- [`cnos.cds/skills/cds/CDS.md §"Selection function"`](../../../cnos.cds/skills/cds/CDS.md) — P0 override, operational-infrastructure override, assessment-commitment default, stale-backlog re-evaluation, MCI freeze check, weakest-axis rule, maximum-leverage rule, dependency order, effort-adjusted tie-break, no-gap case, selection inputs
- [`cnos.cds/skills/cds/CDS.md §"Development lifecycle"`](../../../cnos.cds/skills/cds/CDS.md) — 0–13 step table, S0–S12 state machine, branch rule (`cycle/{N}` canonical; γ creates from `origin/main`), branch pre-flight, skill loading tiers
- [`cnos.cds/skills/cds/CDS.md §"Coordination surfaces"`](../../../cnos.cds/skills/cds/CDS.md) — cycle-state evidence, polling primitives (gh/MCP/git), wake-up mechanism, reachability preflight, transition-only emission, `git fetch` reliability re-probe, mid-flight clarification, cross-repo proposals
- [`cnos.cds/skills/cds/CDS.md §"Artifact contract"`](../../../cnos.cds/skills/cds/CDS.md) — terminology, bootstrap, ordered artifact flow, manifest, **§"Location matrix"** (canonical paths for `.cdd/unreleased/{N}/`, releases, PRA, snapshot dirs, bare-tag rule), **§"Ownership matrix"**, **§"Frozen snapshot rule"**, supporting rules
- [`cnos.cds/skills/cds/CDS.md §"CDS Trace"`](../../../cnos.cds/skills/cds/CDS.md) — the trace format (lifecycle-step ordered; renamed from "CDD Trace")
- [`cnos.cds/skills/cds/CDS.md §"Mechanical vs judgment"`](../../../cnos.cds/skills/cds/CDS.md) — mechanical axes vs judgment axes; class invariants
- [`cnos.cds/skills/cds/CDS.md §"Review CLP"`](../../../cnos.cds/skills/cds/CDS.md) — CLP form (TERMS / POINTER / EXIT); reviewer ask list
- [`cnos.cds/skills/cds/CDS.md §"Gate"`](../../../cnos.cds/skills/cds/CDS.md) — release-readiness preconditions, **§"Closure verification checklist"** (F1–F10 preserved as `#### F{N}:` sub-anchors)
- [`cnos.cds/skills/cds/CDS.md §"Assessment"`](../../../cnos.cds/skills/cds/CDS.md) — PRA contents, **§"Cycle iteration triggers"** (review rounds > 2, mechanical ratio > 20%, avoidable failure, loaded-skill-failed-to-prevent), friction log, engineering levels (L5/L6/L7)
- [`cnos.cds/skills/cds/CDS.md §"Closure"`](../../../cnos.cds/skills/cds/CDS.md) — immediate outputs, deferred outputs, closure rule
- [`cnos.cds/skills/cds/CDS.md §"Retro-packaging"`](../../../cnos.cds/skills/cds/CDS.md) — direct-to-main exception handling (retro-snapshot, self-coherence, version-history entry)
- [`cnos.cds/skills/cds/CDS.md §"Large-file authoring rule"`](../../../cnos.cds/skills/cds/CDS.md) — files > 50 lines section-by-section, section-manifest HTML-comment header, resumption protocol
- [`cnos.cds/skills/cds/CDS.md §"Empirical anchor"`](../../../cnos.cds/skills/cds/CDS.md) — shape-compatibility claim, representative cycle milestones, sub-7 deferred surface-by-surface mapping
- [`cnos.cds/skills/cds/CDS.md §"Non-goals" → §"Software-cycle non-goals"`](../../../cnos.cds/skills/cds/CDS.md) — five software-cycle non-goals (do not optimise primarily for speed; do not treat issue queues as self-justifying; do not reduce review to local diff reading; do not treat release as "tag and hope"; do not confuse a shipped feature with a closed coherence cycle)

**Anchor convention for cross-references.** Citing skill files cite CDS sub-anchors directly (e.g. `cnos.cds/skills/cds/CDS.md §"Artifact contract" → §"Location matrix"`, `§"Gate" → §"Closure verification checklist" → §"F6"`). The pre-#402 CDD.md anchor forms (`§1.4`, `§1.6a`, `§1.6c`, `§5.2`, `§5.3a`, `§5.3b`, `§9.1`, `§Tracking`, `Phase 6 step 17`) are retired: their content has migrated to CDS canonical homes and citing files use the new anchors. The cdd skill files (alpha/beta/gamma/operator/release/release-effector/post-release/harness/activation/review) carry the re-pointed citations as of [cnos#411](https://github.com/usurobor/cnos/issues/411). Where a v0.1 overlay still resides in this package (`alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `operator/SKILL.md` for the operational expansion of CDS sections that have not yet had a CDS-side role rewrite), the overlay file is named explicitly; the doctrinal home remains CDS.

## Hard rule

The essay `docs/papers/CCNF-AND-TYPED-TRUST.md` pins: *"Do not finalize CDD.md until V works and domain evidence has somewhere else to live."* Both preconditions hold as of this rewrite:

- **V is executable.** The cn binary at `src/packages/cnos.cdd/commands/cdd-verify/` implements `V : Contract × Receipt → ValidationVerdict`. The operator-facing wrapper `cn cdd verify --receipt <path>` dispatches into V. Shipped under [cnos#392](https://github.com/usurobor/cnos/issues/392) (Phase 3 of #366).
- **Domain evidence has homes.** `schemas/cdd/` (generic), `schemas/cds/` (software), and `schemas/cdr/` (research) all exist on `origin/main` per [cnos#388](https://github.com/usurobor/cnos/issues/388) (Phase 2.5 — generic/domain schema split). `cnos.cdr` v0.1 shipped per [cnos#376](https://github.com/usurobor/cnos/issues/376). `cnos.cds` v0.1 shipped under the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave — the canonical software-realization doctrine lives at [`cnos.cds/skills/cds/CDS.md`](../../../cnos.cds/skills/cds/CDS.md).

No domain-specific evidence requirements appear in the kernel sections of this document (§Kernel, §Outcomes, §Recursion modes, §Scope-lift, §Domain packages, §Pointers). Software-specific vocabulary (test, code, branch, deploy, CI, release) has been extracted to [`cnos.cds/skills/cds/CDS.md`](../../../cnos.cds/skills/cds/CDS.md) per cnos#403; CDD.md's §"Software-specific realization → cnos.cds" section retains only the pointer list into CDS.

## Non-goals

CDD's kernel doctrine does not:

- Name any particular tooling, platform, dispatch mechanism, schema language, or invocation surface in the kernel sections (§Kernel through §Scope-lift) — those names belong in domain packages and runtime substrate.
- Formalize CDD orchestration grammar (mode enum, sizing predicate, master+sub graph, dispatch-prompt schema, findings state machine) — that is the CCNF-X follow-on, not this document.
- Define `cnos.cds` package contents — the cds canonical doctrine lives at [`cnos.cds/skills/cds/CDS.md`](../../../cnos.cds/skills/cds/CDS.md) per the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave, not in this document.
- Re-derive `COHERENCE-CELL.md` doctrine or `RECEIPT-VALIDATION.md`'s typed interface — those are the input doctrines this document cites.
