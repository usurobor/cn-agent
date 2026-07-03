# Glossary – cnos v3.6.0

Short definitions of the main terms used in cnos, the CN whitepaper, and the coherence system documentation.

> **Note:** Document versions (e.g., GLOSSARY v2.0.0) are local to each file. See `CHANGELOG.md` for the overall template version.

---

## Two senses of α, β, γ (and δ, ε, V)

The Greek letters carry two independent, current meanings. Both are live; neither replaces the other.

1. **TSC measurement axes** — α (Pattern), β (Relation), γ (Exit/Process): a three-axis score applied *to* an artifact or cycle after the fact. See **TSC**.
2. **CDD/CDS role actors** — α produces matter, β reviews it, γ coordinates the α↔β loop, δ operates the boundary and dispatch cadence, ε iterates the protocol itself, V validates a receipt against its contract. See **Role (α/β/γ/δ/ε)** and **V** below.

A cell's γ *coordinates a cycle*; that same cycle's output is later *scored* on the γ (Exit) axis. These are different objects that share a Greek letter — role output is not axis score.

**The filesystem is organized for readers. α/β/γ are role grammar and measurement axes, not docs folders.** There is no `docs/alpha/`, `docs/beta/`, or `docs/gamma/` — that filing scheme was retired; content lives under `docs/{quickstart,concepts,guides,reference,architecture,development,papers,evidence}/`, organized by what a reader is trying to do. See [`docs/README.md`](../../README.md).

---

## Coherence System Terms

### Coherence

The degree to which a system's model matches reality. Incoherence is the gap between model and reality. The foundational claim of cnos is that coherence is primary — everything in the system exists to make coherence visible, bounded, reviewable, recoverable, and evolvable.

Measured across three axes: α (Pattern), β (Relation), γ (Exit) — the TSC-axis sense; see "Two senses of α, β, γ" above. See **TSC**.

Defined in: `src/packages/cnos.core/doctrine/COHERENCE.md`, `docs/papers/COHERENCE-SYSTEM.md`.

### Coherence delta

A bounded movement from a less coherent state to a more coherent one. The architectural unit of meaningful change in cnos. Each agent cycle, each development iteration, and each release should produce a measurable coherence delta.

A coherence delta is not merely a feature or a fix — it is the change in coherence itself, of which the feature or fix is the concrete, operator-visible articulation.

Used in: `docs/reference/runtime/CAA.md` §5.6 "The cycle as coherence delta", §10 "Coherence as Direction".

### CAP (Coherent Agent Principle)

The dynamic atom of coherence. When an agent detects a gap between model and reality, there are two coherent responses:

- **MCA** — Most Coherent Action: change reality (act on the world)
- **MCI** — Most Coherent Insight: change the model (update understanding)

Priority rule: **MCA before MCI.** If you can act, act first. If you cannot act, learn. If both are needed, act then learn from the result.

> Gap → MCA or MCI → MCA first

CAP is derived from Friston's Free Energy Principle, reframed for agents. It is doctrine — always-on, non-negotiable.

Defined in: `src/packages/cnos.core/doctrine/CAP.md`.

### MCA (Most Coherent Action)

One of two mechanisms for closing the gap between model and reality. MCA changes reality to match the model. Examples: fix a bug, ship a feature, correct a runtime behavior.

MCA is preferred over MCI when coherent action is possible.

Defined in: `src/packages/cnos.core/doctrine/CAP.md` §2.1.

### MIC (Make It Coherent)

Verb. Cohere an artifact to its canonical contract for the first time. The artifact has content but was never structured to the current spec. "MIC the release skill" = it has good content, DUR contract now exists, give it the right shape.

MIC is authoring coherence where none existed (relative to the current contract). The artifact may have been internally consistent before, but it wasn't aligned to the spec that now governs its class.

### MICA (Make It Coherent Again)

Verb. Restore coherence that was lost through drift. An artifact was coherent to its contract — then context evolved (new contracts, new specs, new standards) and the artifact didn't move with it. MICA is the act of bringing it back.

The "again" does real work: MICA asserts a prior coherent state. You can't MICA something that was never coherent to the current contract (that's MIC). You can't MICA something already aligned (that's a no-op). MICA only applies in the gap between "was coherent" and "context moved."

Related: MIC (first-time coherence), DUR (the contract being cohered to), CLP (the process used to verify convergence).

### MCI (Most Coherent Insight)

The second mechanism for closing the gap between model and reality. MCI changes the model to match reality. Examples: update documentation, revise an assumption, change the design.

MCI is used when MCA is blocked or when the model itself is wrong.

Defined in: `src/packages/cnos.core/doctrine/CAP.md` §2.2.

### MCP (Most Coherent Picture)

The current best picture the system can form of itself and its world. Not omniscience — the best currently available picture across α, β, and γ.

An MCP includes:
- the relevant articulations
- the known relations among them
- the visible tensions and gaps
- the current path pressures / likely exits

MCP is what Sense → Compare (FOUNDATIONS §3) produces implicitly at the agent scale, named explicitly at the system scale.

Defined in: `docs/papers/COHERENCE-SYSTEM.md` §3.1; used in `docs/reference/runtime/CAA.md` §5.2, `docs/development/cdd/CDD.md` §2.1.

### CMP (Coherence Mapping Pass)

The operation that produces or refreshes the MCP. CMP is the system-scale expression of the Sense → Compare phase.

At the agent scale, CMP is the sensing / comparison phase that precedes the MCA-or-MCI choice. At the development scale, CMP maps the relevant artifacts (α), their relations (β), and the pressures on future movement (γ) before choosing a development intervention.

CMP asks: what are the relevant articulations right now? How do they relate? Where are the strongest incoherences? What is the weakest axis?

Defined in: `docs/papers/COHERENCE-SYSTEM.md` §3.2; used in `docs/reference/runtime/CAA.md` §5.1, `docs/development/cdd/CDD.md` §4.2.

### CLP (Coherence Ladder Process)

The reflective law that governs movement. CLP is the review rhythm that keeps CAP from drifting. Structure:

1. **Seed** — state the gap clearly
2. **Bohmian reflection** — dialog that seeks the real structure of the problem
3. **Triadic check** — score across α (Pattern), β (Relation), γ (Exit)
4. **Patch weakest axis** — do not polish the strongest axis; fix the weakest one
5. **Repeat** — until coherent enough to proceed

CLP also provides a dialogue structure for reviews and discussions:
- **TERMS** — what is being discussed; shared vocabulary
- **POINTER** — where the tension is; what's incoherent
- **EXIT** — what changed; what's next

Never publish cold. Always run at least one CLP cycle first.

Defined in: `src/packages/cnos.core/doctrine/COHERENCE.md`.

### TSC (Triadic Self-Coherence)

A framework for measuring coherence across three algebraically independent axes:

- **α (Alpha) — PATTERN**: Internal consistency. Are the articulations non-contradictory? Do the parts hold together?
- **β (Beta) — RELATION**: Alignment across views and articulations. Do all layers reveal the same system?
- **γ (Gamma) — EXIT / PROCESS**: Viable evolution path. Can the system change without losing itself?

Composite score: `C_Σ = (s_α · s_β · s_γ)^(1/3)`. Default PASS threshold: **Θ ≥ 0.75** (teams may set a stricter override, e.g. 0.90, for self-application or safety-critical use).

This corrects the earlier-cited `C_Σ ≥ 0.80` threshold, per the 2026-06-23 v3.0.1 errata recorded in `docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md`: *"Corrects the referenced TSC normative threshold from `PASS ≥ 0.80` to `PASS ≥ 0.75`... No CN protocol semantics changed."*

Originated by usurobor. Formal spec: `tsc/spec/tsc-core.md` (https://github.com/usurobor/tsc).

Defined in: `src/packages/cnos.core/doctrine/COHERENCE.md`.

### Release Coherence Ledger

The authoritative per-release record in `CHANGELOG.md`. Each row captures two dimensions of a release cycle:

- **Coherence quality** — TSC grades (C_Σ, α, β, γ) measuring how cleanly the release shipped
- **Architectural scope** — Engineering level (L5/L6/L7) measuring what boundary moved

The ledger is append-only. TSC grades are provisional at release time and may be revised by the post-release assessment. Engineering levels are classified per `docs/development/ENGINEERING-LEVELS.md`, capped by CDD §9.1 cycle execution quality (mechanical ratio, review rounds).

Historical classification: `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md`.

Defined in: `CHANGELOG.md`.

### The Core Coherence Algorithm

The recurrent algorithm of the system at every scale:

> **CMP → CAP (MCA/MCI) → CLP → update articulations → CMP**

1. CMP — build the MCP
2. CAP — choose the coherent move (MCA first)
3. CLP — review the result across α / β / γ
4. Update — record the result
5. CMP — build the next picture

At the agent scale: Sense → Compare → MCA/MCI → Repeat (FOUNDATIONS §4).
At the system scale: CMP → CAP → CLP → update → CMP (COHERENCE-SYSTEM §3.3).

### Bohmian Reflection

Dialog that seeks the real structure of a problem before acting. Not brainstorming. Core questions:
- What is the system trying to become?
- What are we implicitly assuming?
- Where are two truths coexisting?
- What would make this simpler and truer?

Used in CLP step 2 and CDD §4.2.

### Coherence Walk

The practice of rebalancing between coherence axes after scoring:

1. Score α, β, γ
2. Reflect on what contributed to each score
3. Set a goal for the next cycle, investing in the lower axis

If α < β, invest in PATTERN. If β < α, invest in RELATION. Balance two, let the third emerge.

### Articulation

Any durable expression of coherence at a particular scale. Doctrine, code, skills, packages, traces, release notes, and agents are all articulations. cnos is a network of recurrent coherent articulations that generate, constrain, package, execute, observe, and repair one another.

Defined in: `docs/papers/COHERENCE-SYSTEM.md` §5 "Articulation Types".

---

## CDD/CDS Roles, Cells, and Waves

The CDD/CDS role-actor sense of α/β/γ, plus δ/ε/V. See "Two senses of α, β, γ" above for how this differs from the TSC-axis sense.

### Role (α, β, γ, δ, ε)

The generic scope-escalation role ladder: α **produces** the matter; β **reviews** α's matter; γ **coordinates** the α↔β loop within one cycle; δ **operates** the sequence of γ-loops across cycles; ε **iterates** the δ-discipline itself. Roles are structural positions, not job titles — the same actor may hold different roles across protocols, but α and β must be distinct actors within one cycle.

Defined in: `ROLES.md` §1 "The role ladder", §9 Glossary.

### δ (delta)

The cell's boundary actor. Reads the receipt γ emits at close-out and V's validation verdict, records a `BoundaryDecision` (accept / release / override / reject / repair_dispatch), operates wave cadence (see **Wave**), and dispatches every role. Also enriches a dispatch prompt's implementation contract when γ leaves it under-specified.

Defined in: `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`; `ROLES.md` §1.

### ε (epsilon)

The CDD-specific instantiation of the generic ε doctrine (`ROLES.md` §4b). Reads the cross-cycle receipt stream, surfaces `cdd-*-gap` findings, applies MCA discipline, and writes `cdd-iteration.md` when a cycle's receipt carries `protocol_gap_count > 0`. ε iterates the protocol across cycles; it does not review any single cell — that is β's job.

Defined in: `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`.

### V (validator)

The pure validation predicate at a cell's parent-facing boundary: `V : Contract × Receipt × EvidenceGraph → ValidationVerdict`. V emits PASS / FAIL / WARN. V does not decide what happens next — that composition is δ's `BoundaryDecision`, which reads V's verdict alongside the receipt.

Defined in: `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`.

### Cell

A bounded unit of coherence-cell work: it receives a contract, closes internal work through the α (produce) → β (review) → γ (close) triad, has its receipt validated by V, and crosses (or does not cross) the boundary under δ's decision. One issue, one α↔β↔γ loop, one receipt. See **Cell vs wave** below.

Defined in: `docs/papers/CELL-OF-CELLS.md` §3 "Cell as recursive type"; `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (the cell-shaped issue contract); `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` (draft recursion doctrine — defers to `CDD.md` for current executable behavior).

### Wave

A δ-run sequence of related cells under one durable, git-committed manifest at `.cdd/waves/{wave-id}/manifest.md`. A wave groups sub-issues that share dispatch order, standing permissions, and scope so γ does not have to author a fresh scaffold for every sub-issue.

Defined in: `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` §10 "Wave Coordination"; `docs/development/issues/TAXONOMY.md` ("...when a wave δ may relabel under MCA").

**Cell vs wave.** A cell is one bounded unit of work — one issue, one α↔β↔γ loop, one receipt. A wave is an ordered *sequence* of cells under one δ-owned manifest; the wave as a whole is not reviewed by a single β pass — each cell inside it is. "The wave shipped" and "one cell in the wave shipped" are not the same claim.

### Matter

The deliverable α produces against a contract: `matterₙ := αₙ.produce(contractₙ)`. The matter type is protocol-specific — source code and tests in CDD, editorial prose in CDW. An accepted child receipt can itself become matter at the enclosing scope. This is the CDD-cell sense of "deliverable" — distinct from the CN Shell runtime substrate (see **CN Shell**).

Defined in: `docs/papers/CELL-OF-CELLS.md` §3; `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`; `ROLES.md` §9 Glossary.

### Receipt

Two current, distinct senses — do not conflate them:

1. **CN Shell runtime receipt** — a record of one CN Shell execution/mutation, stored under `state/receipts/` in an agent's hub, enabling crash recovery, auditability, and idempotency. See **Receipts** under CN Shell Terms below.
2. **CDD cell receipt** — the artifact γ emits at a cell's close-out, carrying the contract, evidence, review, and closure blocks. V validates this receipt against its contract; δ reads V's verdict plus the receipt to record a `BoundaryDecision`. An accepted receipt can become matter for the parent cell.

Sense 1 is per-mutation runtime bookkeeping; sense 2 is per-cell close-out evidence consumed by V and δ.

Defined in (sense 2): `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`; `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`.

### Review request vs review verdict

Not two formally separate named artifacts today — one signal, one decision:

- **Review request** — α's outbound signal that a cell is ready for β: the review-readiness section α appends to `.cdd/unreleased/{N}/self-coherence.md` once the pre-review gate passes.
- **Review verdict** — β's rendered judgment: RC (request changes) or APPROVE, written to `.cdd/unreleased/{N}/beta-review.md`.

The request asks for review; the verdict is the review's outcome. There is no third, separately-named "review request" artifact type in the current skills — treat any claim of one as overclaiming.

Defined in: `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.7 "Request review"; `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md`.

### Projection vs role-renaming

A **projection** is a typed object — a boundary decision or a receipt-stream observation — that crosses a scope boundary and becomes input to the next scope's roles. **Role-renaming** is the rejected idea that a role at scope `n` is simply relabeled as a different role at scope `n+1`.

δₙ's boundary decision projects onto scope `n+1` as input to βₙ₊₁'s discrimination; εₙ's receipt-stream observation projects onto scope `n+1` as input to γₙ₊₁'s coordination. But δₙ is not literally βₙ₊₁, and εₙ is not literally γₙ₊₁ — the scope-`n+1` cell has its own full α, β, γ, δ, ε; only the *decision* or *observation* crosses, not the role.

Defined in: `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`: *"This is a projection under scope-lift. It is not a flat role-renaming inside a single cell. δₙ is not literally βₙ₊₁; εₙ is not literally γₙ₊₁."*

### Trust claim vs coherence witness

Two independent, composable properties of a receipt:

- **Trust claim** — the strength of attribution carried by a receipt: is it signed, attested, or unsigned?
- **Coherence witness** — whether the artifacts still describe one coherent unit of work — what TSC measures.

A receipt can carry a strong trust claim and a weak (or absent) coherence witness, or vice versa; neither implies the other. The system must not claim more trust than the artifact carries, and must not claim more coherence than it measured.

Defined in: `docs/papers/DUMB-MODELS-SMART-CELLS.md`.

### Signature vs attestation

Two distinct trust-claim components:

- **Signature** — cryptographic proof that an artifact is unaltered and tied to an identity.
- **Attestation** — an external or third-party witness statement about the event.

A receipt may carry a signature, an attestation, both, or neither. If signed, the signature should be visible; if attested, the attestation should be visible.

Defined in: `docs/papers/DUMB-MODELS-SMART-CELLS.md`.

### Issue (cell-shaped contract)

An executable work contract that names the incoherence, impact, source of truth, scope, acceptance criteria, and proof plan for one cell. γ authors it; α executes it; β verifies against it. Not every GitHub issue is cell-shaped — see **dispatch:cell**.

Defined in: `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md`.

### dispatch:cell

An issue label applied **only** when the issue is genuinely executable by CDS/CDD — a scoped, proof-carrying cell — never to a bare design/tracking/research issue. Paired with a `protocol:{cds,cdd}` label naming the concrete protocol. The standard dispatch selector is `dispatch:cell + protocol:{p} + status:todo`.

Defined in: `docs/development/issues/TAXONOMY.md` §"Dispatch and protocol".

### effort/*

An ordinal size-estimate label (`effort/S`, `effort/M`, `effort/L`, `effort/XL`) on an issue — implementation/review size, **not** priority and not calendar time. Missing effort means *unestimated*, not small. `kind/tracking` issues are excluded from effort rollups by default.

Defined in: `docs/development/issues/TAXONOMY.md` §"Effort labels".

---

## Doctrinal Terms

### CBP (Coherent Behavior Principle)

The ethical complement to CAP. While CAP governs internal coherence (MCA/MCI), CBP governs relational coherence with others. Core values: **Peace, Love, Unity, Respect (PLUR).**

CAP without CBP produces effective but harmful agents. CBP without CAP produces kind but ineffective agents.

Defined in: `src/packages/cnos.core/doctrine/CBP.md`.

### CA-Conduct

How every Coherent Agent must behave. Built on PLUR as the absolute foundation, with operational principles: Ship (done beats perfect, bias for action), Own (radical ownership, surface MCAs), Truth (radical candor), Learn (kaizen, go and see), Prepare (search before asking, use existing skills).

Defined in: `src/packages/cnos.core/doctrine/CA-CONDUCT.md`.

### AGENT-OPS

Runtime emission discipline. Always-on rules for agent output format. Defines the output contract (frontmatter with id, coordination ops, typed capability ops, response body), RACI discipline, and the principle that if you see something, you capture and track it.

Defined in: `src/packages/cnos.core/doctrine/AGENT-OPS.md`.

### Doctrine

Always-on, constitutive, non-competitive cognitive substrate. Doctrine defines first principles, conduct, runtime grammar, and review law. Loaded at every wake-up. Never competes for skill slots. Never dropped for context budget. Doctrine is identity, not tooling.

The four doctrinal layers (FOUNDATIONS.md):
1. **CAP** — the dynamic atom (when/how to act or learn)
2. **COHERENCE** — the review geometry (whether the result remained whole)
3. **CBP + CA-Conduct** — the relational boundary (trust preservation)
4. **AGENT-OPS** — the runtime grammar (how the agent emits output)

### Mindset

Always-on, orienting, deterministic cognitive frame. Mindsets shape craft and style. They are not optional enrichments — they are structural orientation. Like doctrine, mindsets never compete for skill slots.

Examples: ENGINEERING, PM, WISDOM, OPERATIONS.

### DUR (Define / Unfold / Rules)

The canonical skill contract. The three-section structure every skill MUST follow:

1. **Define** — identify the parts, articulate how they fit, name the failure mode
2. **Unfold** — expand each part into domain-specific structure with ❌/✅ pairs
3. **Rules** — numbered, stable-ID rules with ❌/✅ pairs showing incoherent/coherent behavior

DUR is to skills what TSC is to coherence scoring — the structural invariant that makes the class recognizable and reviewable.

"MIC X to DUR" = cohere skill X to Define/Unfold/Rules for the first time. "MICA X" = restore DUR coherence lost through drift.

Defined in: `src/packages/cnos.core/skills/skill/SKILL.md` §1–§3 (Define / Unfold / Rules).

### Skill / SKILL.md

Bounded, selected, instrumental cognitive module. Skills are situational amplifiers, not identity. Each skill is authored as a `SKILL.md` file at `skills/{category}/{name}/SKILL.md` in its package source (e.g., `src/packages/cnos.core/skills/agent/self-cohere/SKILL.md`), installed into `.cn/vendor/packages/{pkg}/skills/{category}/{name}/`. `SKILL.md` is both the literal file name and cnos's shorthand for "a skill module" — its frontmatter (`name`, `description`, `triggers`, `scope`, …) makes it selectable and scoreable at wake-up.

If doctrine competes with task skills, the architecture is already confused.

---

## Architecture Terms

### Coherent Agent (CA)

An articulation of coherence that can sense, compare, choose, act or learn, and remain itself while evolving. Specifically, an agent that can:
- perform CMP to form an MCP
- detect a gap between model and reality
- choose MCA or MCI (MCA first)
- execute a bounded move
- review the result via CLP
- continue the loop

Defined in: `docs/reference/runtime/CAA.md`, `docs/papers/COHERENCE-SYSTEM.md` §7.

### CAA (Coherent Agent Architecture)

The design document that specifies what a coherent agent is structurally: definition, first principle, doctrinal layers, cognitive strata at wake-up, the agent loop, runtime embodiment, invariants, and failure modes.

Document: `docs/reference/runtime/CAA.md`.

### CDD (Coherence-Driven Development)

Also a verb. "CDD 47" = take issue #47 through the full pipeline: gap → design → plan → tests → code → docs → review → release → observe. Each step produces a measured coherence delta.

The development method in which every meaningful change is treated as an intervention on incoherence. CDD applies CAP to the development process itself. CDD is γ at the development scale — the expression of evolution applied to cnos.

Each substantial release is a **measured coherence delta**. Features are the operator-facing articulation of that movement.

Canonical algorithm: `src/packages/cnos.cdd/skills/cdd/CDD.md` (the recursive coherence-cell protocol — CCNF kernel, cell outcomes, recursion modes, scope-lift). `docs/development/cdd/CDD.md` is a pointer to that file, not a duplicate.

### CN Shell

The capability runtime that mediates between the agent and the world. The agent proposes typed ops; CN Shell validates against policy, executes within budget, records receipts, and feeds evidence back. Enforces the N-pass bind loop:

- **Observe pass**: agent requests observe ops → runtime gathers evidence, defers effects
- **Effect pass**: agent proposes effect ops → runtime executes governed effects
- **Terminal pass**: no ops → final projection

The loop is bounded by `max_passes` (default 5), `max_total_ops`, and `max_total_artifact_bytes`. Each pass is one packed context → one LLM call → one execution step. This is CAP made runtime-real: sensing is first-class, action is governed.

Defined in: `docs/reference/runtime/AGENT-RUNTIME.md`.

### N-pass bind loop

The runtime's enforcement of "observe before effect" via bounded iteration. Each pass classifies its typed ops as observe-class or effect-class:

- **Observe-class pass**: gathers evidence (file reads, state checks, searches), defers effects with receipts
- **Effect-class pass**: executes governed effects (writes, patches, commits)
- **Terminal pass**: no ops remaining → final projection to user

This generalizes the original two-pass structure (`max_passes=2` reproduces it exactly). The loop is the runtime expression of CAP's Sense → Compare → Act cycle, extended to support deeper reasoning chains (observe → effect → verify → adapt).

### Receipts

Records of CN Shell execution. Every mutation produces a receipt. Receipts enable crash recovery, auditability, and idempotency. Stored in `state/receipts/`.

This is the CN Shell runtime sense of "receipt." See **Receipt** under CDD/CDS Roles above for the distinct CDD-cell close-out sense.

### Wake-up / Reconstitution

The process by which a coherent agent reconstitutes itself from local, versioned, installed packages at the start of each cycle. A coherent agent must not wake as an empty chat window. Wake-up loads:

1. Identity (SOUL, USER)
2. Core Doctrine (FOUNDATIONS, CAP, COHERENCE, CBP, CA-Conduct, AGENT-OPS)
3. Mindsets (role and operating frames)
4. Reflections (recent learning)
5. Task Skills (situationally selected, scored)
6. Capabilities (what the runtime supports)
7. Conversation + inbound message

No network. No sibling checkout. Local files only.

Defined in: `docs/reference/runtime/CAA.md` §4 "Cognitive Strata at Wake-Up".

Not to be confused with **Wake / wake-as-skill** below — a distinct, newer sense of "wake" naming a typed `SKILL.md` module rendered into a scheduled workflow.

### Wake / wake-as-skill

A typed `SKILL.md`-shaped module (`artifact_class: wake`, `scope: global`) whose frontmatter `wake:` block declares identity, role, input/output, surfaces, defer-path, and selector, and whose body is the agent prompt. A substrate renderer (`cn install-wake`) materializes the wake into a real workflow file — today a `*.golden.yml` GitHub Actions workflow (see **Golden**). The `#Wake` CUE definition in `schemas/skill.cue` validates the frontmatter contract at I5 time; the renderer enforces the rest at run time.

Defined in: `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`; `schemas/skill.cue` (`#Wake`).

### CAR (Cognitive Asset Resolver)

The package distribution system. Defines how cognitive assets (doctrine, mindsets, skills) are packaged, versioned, installed, and resolved locally. CAR ensures wake-up is deterministic: same lockfile + same hub state → same packed context.

Document: `docs/architecture/cognitive-substrate/CAR.md`.

### Coherence Contract

A lightweight declaration that accompanies every substantial change under CDD. Answers: what gap is being closed? Is this MCA or MCI? Which layer is affected? What is the expected triadic effect (α/β/γ)? What fails if skipped? What is the expected coherence delta?

Current realization: the CCNF kernel's `contractₙ` object (the issue body, its acceptance criteria, scope, and active design constraints — see `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` kernel step 1), operationalized per cell in `.cdd/unreleased/{N}/self-coherence.md` §Gap (`src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` §2.5).

---

## Protocol Terms

### CN (Coherence Network)

A network of agents that use git repositories as their primary surface for specs, threads, and state. Git is the transport. Files are the state. Commits are the communication.

### Hub

A git repository that is an agent's home. Holds identity (`spec/SOUL.md`, `spec/USER.md`), threads, state, and installed cognitive packages. Created by `cn init` + `cn setup`.

Cognition is local: doctrine, mindsets, and skills are installed into `.cn/vendor/packages/` at setup time. The hub is wake-ready from its own files — no external checkout required.

### cnos (Coherence Network OS)

The source repo and runtime for CN agents. Contains the `cn` CLI, cognitive packages (doctrine, mindsets, skills), design documents, and tests. cnos is the current lowest durable software articulation of the coherence system. Agents consume cnos through installed packages, not through a checkout dependency.

### Thread

A Markdown file under `threads/` **in an agent's hub** (not in this cnos source repo — see **Hub**; cnos itself has no root `threads/` directory) that represents a conversation, reflection, or topic.

**Structure:** Within a hub, threads are organized by type under `threads/`:
- `threads/reflections/daily/YYYYMMDD.md` — daily reflections using α/β/γ format
- `threads/reflections/weekly/YYYY-WNN.md` — weekly reflections
- `threads/adhoc/` — topic conversations, reviews, discussions
- `threads/mail/inbox/`, `threads/mail/outbox/` — multi-party threads

### Peer

Another agent or hub that this hub tracks in `state/peers.md`. Peers are discovered via `cn peer add` and synchronized via `cn sync`.

### Agent

A system (usually an AI assistant + host runtime) that:
- has a CN hub
- reads `state/input.md` and writes `state/output.md`
- proposes typed ops; the runtime executes them
- is a pure function: input → output. Never touches files or git directly — `cn` handles all I/O

The agent is the brain. `cn` is the body. Git is the nervous system.

---

## Operational Terms

### Kata

A practice exercise that walks an agent or human through concrete steps to learn or exercise a behavior.

**Location:** `skills/{category}/{name}/kata.md` — katas live alongside the skill they exercise (e.g., `src/packages/cnos.core/skills/agent/self-cohere/kata.md`).

### Coherent Reflection

The structured practice of assessing coherence at regular cadences (daily, weekly, monthly, quarterly, half-yearly, yearly) using TSC's α/β/γ framework. Output is a reflection thread under the appropriate cadence directory in the agent's hub (e.g., `threads/reflections/daily/YYYYMMDD.md`, `threads/reflections/weekly/YYYY-WNN.md`).

Each periodic thread:
1. Scores PATTERN (α), RELATION (β), EXIT (γ)
2. Identifies what contributed to each score
3. Sets a rebalancing goal for the next cycle (Coherence Walk)

### State

Files under `state/` that record the current situation for this hub. Unlike specs, state is expected to change frequently. Includes queue, receipts, conversation history, input/output, and lock files.

### Package

A distributable cognitive unit. Packages contain doctrine, mindsets, and/or skills. Installed locally under `.cn/vendor/packages/`. Current packages under `src/packages/`: `cnos.core` (doctrine, mindsets, core skills), `cnos.eng` (engineering skills), `cnos.cdd` (the coherence-driven-development protocol), `cnos.cds` (its software-lifecycle realization), `cnos.handoff` (inter-agent wire formats).

### Golden

A rendered, checked-in fixture treated as source-of-truth to diff against. `cn install-wake` renders a wake `SKILL.md` (see **Wake / wake-as-skill**) into a `*.golden.yml` GitHub Actions workflow marked "DO NOT EDIT" — it is regenerated from source, not hand-edited; CI fails if the checked-in golden drifts from a fresh render.

Example: `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`.

---

## Document Hierarchy

| Layer | Document | Question |
|-------|----------|----------|
| Meta-model | COHERENCE-SYSTEM.md | How does coherence unfold into cnos? |
| Doctrine | FOUNDATIONS.md, CAP.md, COHERENCE.md, CBP.md, CA-CONDUCT.md, AGENT-OPS.md | Why — first principles |
| Architecture | CAA.md | What is a coherent agent structurally? |
| Runtime | AGENT-RUNTIME.md | How does it execute? |
| Distribution | CAR.md | How does cognition arrive locally? |
| Development | CDD.md | How does cnos evolve coherently? |
| Observability | TRACEABILITY.md | How is it observed? |
