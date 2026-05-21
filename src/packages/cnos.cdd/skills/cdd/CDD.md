# CDD — Recursive Coherence Cell Protocol

**Version:** 4.0.0
**Status:** Draft (Phase 7 of [#366](https://github.com/usurobor/cnos/issues/366) — CCNF spine; terminal phase)
**Placement:** `src/packages/cnos.cdd/skills/cdd/`

CDD owns the **generic recursive coherence-cell algorithm.** This document is the algorithm's canonical statement: it names the kernel (CCNF), the four cell outcomes, the two recursion modes, the three scope-lift projections, and the domain packages that bind the kernel to substrates. CDD does not own substrate-specific realization — software-lifecycle realization lives in `cnos.cds` (pending bootstrap, tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403)); research-lifecycle realization lives in `cnos.cdr`.

The kernel is stated verbatim from [`COHERENCE-CELL-NORMAL-FORM.md`](COHERENCE-CELL-NORMAL-FORM.md) (CCNF). That document is the citable source of the recursion equation; this document is the operational doctrine that names domain bindings, points at canonical surfaces, and quarantines the substrate-specific lifecycle content (currently still resident in §"Software-specific realization") until cds extraction completes.

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

- **CDS — `cnos.cds`** — software-development realization. CDS binds CCNF to source code, tests, releases, deployments. Receipt fixtures live at `schemas/cds/`; the CDS package itself is pending bootstrap (tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403)). Until CDS bootstraps, software-lifecycle realization lives in this file's §"Software-specific realization — pending cds extraction" section.
- **CDR — `cnos.cdr`** — research-domain realization. CDR binds CCNF to research artifacts (claims, evidence, replication, falsification). CDR v0.1 shipped at [cnos#376](https://github.com/usurobor/cnos/issues/376); receipts validate at `schemas/cdr/`.
- **Future c-d-X protocols** — additional domain bindings (e.g., a `cnos.cdo` operations protocol, a `cnos.cdh` hardware protocol). Each c-d-X protocol follows the same pattern: the kernel stays generic; the c-d-X domain package owns the substrate-specific evidence, validators, and decision shapes.

The binding pattern is: a CDS / CDR / c-d-X domain package cites the CDD kernel as input contract, defines the domain's typed receipt schema under `schemas/{cdd-X}/`, names the substrate-specific α/β/γ/δ/ε realizations, and consumes `V` (the generic validator) through that schema. CDS, CDR, and any future c-d-X protocol do not re-derive the kernel; they extend it.

## Pointers

The canonical surfaces CDD cites are grouped by what they own:

**Generic doctrine (the kernel and its companions):**
- [`COHERENCE-CELL.md`](COHERENCE-CELL.md) — predecessor doctrine; receipt rule; four-way structural separation
- [`COHERENCE-CELL-NORMAL-FORM.md`](COHERENCE-CELL-NORMAL-FORM.md) — CCNF; the recursion equation source for this document
- [`RECEIPT-VALIDATION.md`](RECEIPT-VALIDATION.md) — V's design surface; input refs, output verdict shape, invocation rule
- `ROLES.md` §1, §3, §4, §4a, §4b — generic role doctrine (§4a: δ; §4b: ε protocol-iteration)

**Schemas (typed evidence by domain):**
- `schemas/cdd/` — generic kernel schemas (`contract.cue`, `receipt.cue`, `boundary_decision.cue`, `validation_verdict.schema.json`)
- `schemas/cds/` — software domain schemas (receipt + fixtures; package pending per [cnos#403](https://github.com/usurobor/cnos/issues/403))
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
- `src/packages/cnos.cds/` — software realization (pending bootstrap per [cnos#403](https://github.com/usurobor/cnos/issues/403))
- `src/packages/cnos.cdr/` — research realization (shipped per [cnos#376](https://github.com/usurobor/cnos/issues/376))

**Loader and rationale:**
- [`SKILL.md`](SKILL.md) — package-visible loader entrypoint (not a second fact source)
- `docs/gamma/cdd/RATIONALE.md` — companion rationale for shape decisions
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` — the essay that pins the CCNF spine and the kernel/realization separation

## Software-specific realization — pending cds extraction

The content below is **software-lifecycle realization detail** that the pre-CCNF CDD.md carried inline. Per the essay's Wave 5 ("CDS extraction by reference") and Wave 7 ("final CDD.md rewrite") framing, this content migrates to `cnos.cds` when that package bootstraps. **Tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403).** Until cds bootstraps, the named surfaces below are still resident in this file and remain the citable home for cross-references from the cdd/ skill files. **No surface listed below is silently dropped; the citing skill files (e.g. `gamma/SKILL.md`, `alpha/SKILL.md`, `release/SKILL.md`, `post-release/SKILL.md`, `harness/SKILL.md`, `operator/SKILL.md`, `activation/SKILL.md`) continue to resolve their `CDD.md §X` citations against this section family until #403 lands and re-points them at the cds package.**

- **Inputs (§Inputs)** — selection inputs the software-cycle observer reads: CHANGELOG TSC table, encoding lag table, doctor/status health surface, last post-release assessment. Operational realization in `gamma/SKILL.md`.
- **Selection function (§Selection)** — P0 override, operational-infrastructure override, assessment-commitment default, stale-backlog re-evaluation, MCI freeze check, weakest-axis rule, maximum-leverage rule, dependency order, effort-adjusted tie-break, no-gap case. Operational realization in `gamma/SKILL.md`.
- **Development lifecycle (§Lifecycle)** — the 0–13 step table; the lifecycle state machine S0–S12; the branch rule (`cycle/{N}` canonical, γ creates from `origin/main` before dispatch); branch pre-flight; skill loading tier structure (1a CDD authority, 1b lifecycle phase skills, 1c β closure bundle, 2 general engineering, 3 issue-specific). Operational realization in `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`.
- **Roles and dispatch (§Roles)** — the triadic rule (α produces, β judges and merges, γ orchestrates); the dyad-plus-coordinator framing; δ at the external boundary; γ/α/β/δ algorithms; γ dispatch prompt format; named operator-decision points; sequential bounded dispatch model (§1.6); re-dispatch prompt formats (§1.6a, §1.6b); initial dispatch sizing, prompt scope, and commit checkpoints (§1.6c). Operational realization in the role SKILL files plus `operator/SKILL.md` and `harness/SKILL.md`.
- **Coordination surfaces (§Tracking)** — issue activity; cycle branch state; `.cdd/unreleased/{N}/` directory state; polling query forms (gh / MCP / git); wake-up mechanism; reachability preflight; transition-only emission; synchronous baseline pull; `git fetch` reliability re-probe; issue-edit cache-bust via `gamma-clarification.md`; cross-repo proposal lifecycle and STATUS state machine. Operational realization in `harness/SKILL.md §5.4`, `gamma/SKILL.md`, `cross-repo/SKILL.md`.
- **Artifact contract (§Artifacts)** — terminology (post-release, assessment, close-out, closure); bootstrap; ordered artifact flow (design → contract → plan → tests → code → docs → self-coherence → review → gate → release → observe → assess → close); artifact manifest with per-step format spec; **Artifact Location Matrix** (canonical paths for `.cdd/unreleased/{N}/self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md`, `RELEASE.md`, the version-snapshot directory, the PRA, the cross-repo trace dir); **role/artifact ownership matrix**; CDD Trace format; supporting rules; frozen snapshot rule. Operational realization in `release/SKILL.md`, `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`.
- **Mechanical vs judgment boundary (§Mechanical)** — what may be enforced by tools (branch naming, version-directory presence, AC accounting, stale-cross-reference detection, gate checks) vs what remains judgment-bearing (the real incoherence, MCA vs MCI, scoring soundness, review convergence, design coherence, when to stop iterating).
- **Review (§Review)** — CLP (TERMS / POINTER / EXIT); reviewer ask list (α/β/γ scores, weakest-axis diagnosis, concrete patch suggestions, iterate-or-converge verdict). Operational realization in `review/SKILL.md`.
- **Gate (§Gate)** — release-readiness preconditions; the closure verification checklist F1–F10 (missing α/β/γ close-outs, stale `.cdd/unreleased/{N}/` after release, missing RELEASE.md, δ tag ordering, α close-out re-dispatch mechanism, PRA presence). Operational realization in `release/SKILL.md`, `gamma/SKILL.md §2.10`, `operator/SKILL.md`.
- **Assessment (§Assessment)** — PRA contents (measured coherence delta, encoding lag, MCA/MCI balance, process learning, review quality metrics, CDD self-coherence on α/β/γ axes, cycle iteration, next-move commitment); **§9.1 cycle iteration** triggers (review rounds > 2; mechanical ratio > 20%; avoidable tooling/environmental failure; loaded skill failed to prevent a finding); friction log; root cause classification; skill impact; MCA; **cycle-level L5/L6/L7 framework** (per `docs/gamma/ENGINEERING-LEVELS.md`). Operational realization in `post-release/SKILL.md`.
- **Closure (§Closure)** — immediate outputs (changelog corrections, missing documentation, skill/process micro-patches, skill patches identified by cycle iteration, issue filing, lag-table updates, hub-memory writes); deferred outputs (next-MCA issue number, owner, target branch name, first AC, MCI freeze/resume state); closure rule (cycle closes only when immediate outputs executed, deferred outputs committed, cycle iteration present if triggered). Operational realization in `gamma/SKILL.md`, `post-release/SKILL.md`.
- **Retro-packaging rule (§Retro-packaging)** — direct-to-main exception handling (retro-snapshot in version directory, self-coherence covering the change, version-history entry). Operational realization in `release/SKILL.md`.
- **Non-goals (§Non-goals)** — software-cycle non-goals: do not optimize primarily for speed; do not treat issue queues as self-justifying; do not reduce review to local diff reading; do not treat release as "tag and hope"; do not confuse a shipped feature with a closed coherence cycle.
- **Large-file authoring rule (§Large-file)** — files > 50 lines written section-by-section to disk with the section-manifest HTML-comment header and the resumption protocol. Operational realization in this section's own usage pattern.

**Anchor convention for cross-references:** Citing skill files use specific anchor forms (e.g. `CDD.md §1.4 γ algorithm Phase 1 step 3a`, `CDD.md §5.3a`, `CDD.md §5.3b`, `CDD.md §9.1`, `CDD.md §1.6a`, `CDD.md §1.6c`, `CDD.md §Tracking`, `CDD.md Phase 6 step 17`). Until [cnos#403](https://github.com/usurobor/cnos/issues/403) re-points each citation at the cds package, those anchors resolve to the named family above (Selection / Lifecycle / Roles and dispatch / Coordination surfaces / Artifact contract / Mechanical / Review / Gate / Assessment / Closure / Retro-packaging / Non-goals / Large-file). A reader following an old anchor reaches the family that owns the cited content; the operational expansion lives in the named role / runtime-substrate SKILL.md file (per the pointers section above).

## Hard rule

The essay `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` pins: *"Do not finalize CDD.md until V works and domain evidence has somewhere else to live."* Both preconditions hold as of this rewrite:

- **V is executable.** The cn binary at `src/packages/cnos.cdd/commands/cdd-verify/` implements `V : Contract × Receipt → ValidationVerdict`. The operator-facing wrapper `cn cdd verify --receipt <path>` dispatches into V. Shipped under [cnos#392](https://github.com/usurobor/cnos/issues/392) (Phase 3 of #366).
- **Domain evidence has homes.** `schemas/cdd/` (generic), `schemas/cds/` (software), and `schemas/cdr/` (research) all exist on `origin/main` per [cnos#388](https://github.com/usurobor/cnos/issues/388) (Phase 2.5 — generic/domain schema split). `cnos.cdr` v0.1 shipped per [cnos#376](https://github.com/usurobor/cnos/issues/376). `cnos.cds` extraction is tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403) — the software-realization content currently in this file's §"Software-specific realization" section migrates there when the package bootstraps.

No domain-specific evidence requirements appear in the kernel sections of this document (§Kernel, §Outcomes, §Recursion modes, §Scope-lift, §Domain packages, §Pointers). Software-specific vocabulary (test, code, branch, deploy, CI, release) is quarantined inside §"Software-specific realization — pending cds extraction" and resolves out of CDD.md when #403 lands.

## Non-goals

CDD's kernel doctrine does not:

- Name any particular tooling, platform, dispatch mechanism, schema language, or invocation surface in the kernel sections (§Kernel through §Scope-lift) — those names belong in domain packages and runtime substrate.
- Formalize CDD orchestration grammar (mode enum, sizing predicate, master+sub graph, dispatch-prompt schema, findings state machine) — that is the CCNF-X follow-on, not this document.
- Define `cnos.cds` package contents — the cds bootstrap is the [cnos#403](https://github.com/usurobor/cnos/issues/403) cycle, not this one.
- Re-derive `COHERENCE-CELL.md` doctrine or `RECEIPT-VALIDATION.md`'s typed interface — those are the input doctrines this document cites.
