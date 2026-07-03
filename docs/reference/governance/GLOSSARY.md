# Glossary — cnos

The canonical vocabulary surface for cnos: the coherence system, the CDD/CDS
work model (cells, waves, receipts), the issue taxonomy, wakes, packages, and
the trust/evidence layer. Definitions are short and precise; each points at the
canonical source with an inline path.

> **Version note.** Document versions are local to each file; see `CHANGELOG.md`
> for the template/release version. This glossary tracks the *current* work
> model — reader-intent docs, α/β/γ as role/measurement grammar, CDS cell
> dispatch, and the issue taxonomy. It is not organized alphabetically first;
> it is organized the way people encounter the system.

> **One rule to read everything else by.** The filesystem is organized **for
> readers** (`quickstart/`, `concepts/`, `guides/`, `reference/`,
> `development/`, `architecture/`, `papers/`, `evidence/`). **α/β/γ are role
> grammar and TSC measurement axes — never docs folders.** There is no
> `docs/alpha/`, `docs/beta/`, or `docs/gamma/`; those directories were retired
> in the docs cleanup and their contents now live under intent directories
> (e.g. the foundational essays are in `docs/papers/`). Cognitive packages live
> under `src/packages/` (e.g. `src/packages/cnos.core/`), not a root
> `packages/`. There is no root `threads/` in this repo — threads are a *hub*
> concept (see **Thread**).

---

## 1. Core coherence terms

### Coherence

The degree to which a system's model matches reality. The foundational claim of
cnos: coherence is primary — everything exists to make coherence visible,
bounded, reviewable, recoverable, and evolvable. Measured across three axes:
α (Pattern), β (Relation), γ (Exit/Process). See **TSC**.

Defined in: `src/packages/cnos.core/doctrine/COHERENCE.md`.

### Incoherence

The gap between model and reality. cnos moves along a *path of decreasing
incoherence*: each cell, cycle, and release should close some measurable slice
of that gap. Incoherence is the thing acted on; coherence is the direction.

### Coherence delta

A bounded movement from a less coherent state to a more coherent one — the
architectural unit of meaningful change. Not merely a feature or a fix, but the
*change in coherence itself*, of which the feature or fix is the concrete,
operator-visible articulation.

### Articulation

Any durable expression of coherence at a particular scale: doctrine, code,
skills, packages, receipts, release notes, agents. cnos is a network of
recurrent coherent articulations that generate, constrain, package, execute,
observe, and repair one another.

### TSC (Triadic Self-Coherence)

The framework for *measuring* coherence across three algebraically independent
axes:

- **α — PATTERN**: internal consistency. Are the articulations non-contradictory?
- **β — RELATION**: alignment across views. Do all layers reveal the same system?
- **γ — EXIT / PROCESS**: viable evolution path. Can it change without losing itself?

Composite score: **CΣ = (sα · sβ · sγ)^(1/3)**. TSC is a *measurement recorded
from a report*, never a hand-authored frontmatter score. It is **not** a single
hardcoded pass gate: in current CDS work the verdict is **V**'s `PASS`/`FAIL`
against the cell contract plus **δ**'s boundary decision (see **V**, **δ**), not
a universal "CΣ ≥ 0.80" threshold.

Defined in: `src/packages/cnos.core/doctrine/COHERENCE.md`.

### α / β / γ (as measurement axes)

The three TSC axes — Pattern, Relation, Exit. Used as a *measurement* grammar
and as *role* grammar (see §3), never as a filing taxonomy. If you see α/β/γ,
ask "which axis / which role" — never "which folder."

### CAP (Coherent Agent Principle)

The dynamic atom of coherence. On detecting a gap between model and reality,
there are two coherent responses: **MCA** (change reality) or **MCI** (change
the model). Priority rule: **MCA before MCI** — if you can act, act; if you
cannot, learn. Doctrine — always-on, non-negotiable.

Defined in: `src/packages/cnos.core/doctrine/CAP.md`.

### MCA (Most Coherent Action)

A small action inside an approved contract that moves the system toward
coherence without needing operator interruption. The action side of CAP: change
reality to match the model (fix a bug, ship a feature, correct behavior).
Preferred over MCI when coherent action is possible.

Defined in: `src/packages/cnos.core/doctrine/CAP.md`.

### MCI (Most Coherent Insight)

The insight side of CAP: change the model to match reality (update a doc, revise
an assumption, change the design). Used when MCA is blocked or the model itself
is wrong.

Defined in: `src/packages/cnos.core/doctrine/CAP.md`.

### MCP (Most Coherent Picture)

The current best picture the system can form of itself and its world across α,
β, and γ — the relevant articulations, their known relations, the visible
tensions, and the likely exits. Not omniscience; the best currently available
picture. (Distinct from the unrelated "MCP" of Model Context Protocol.)

---

## 2. Work lifecycle

### Cell

A **bounded unit of work** that moves through **contract → matter → review →
receipt → verdict → decision**. The atom of dispatchable, proof-carrying work.
A cell is scoped, carries acceptance criteria and a proof oracle, and closes
with a receipt.

Canonical: `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md`,
`src/packages/cnos.cdd/skills/cdd/CDD.md`.

Cells come in typed kinds (e.g. `implementation`, `issue_authoring`, `wave`,
`cleanup`, `doctrine`) that refine the same generic kernel with a specific
matter shape — see
`src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md`.

### Wave

A **parent cell that coordinates several child cells under one approved
contract**. A wave lets an agent execute a planned campaign without asking the
operator about every minor coherent action. Cell vs wave: a cell is one bounded
unit; a wave is a cell-of-cells with a boundary authority (**wave δ**) that runs
the campaign inside the approved contract.

Canonical: `src/packages/cnos.cdd/skills/cdd/` (recursive cell framework);
wave records under `.cdd/waves/`.

### Matter

The **produced work body of a cell**: code, docs, artifacts, receipts,
generated files, or other concrete output. Matter is what α produces and what β
reviews.

### Receipt

The **durable closeout artifact** that records what happened, what evidence was
bound, what validation ran, what decision was made, and what can be inspected
later. A CDS cell's receipt requires the closure-record evidence set — the five
per-cycle artifacts (`self_coherence`, `beta_review`, `alpha_closeout`,
`beta_closeout`, `gamma_closeout`), the evidence root, the cycle diff, and
optional CI refs — which **V** dereferences when validating.

Schema: `schemas/cdd/receipt.cue`, `schemas/cds/receipt.cue`. Contract:
`src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`.

### Review request

A **proof artifact showing that matter exists and is ready for β review.** It is
**not** a review verdict — it asserts readiness, not acceptance. The verdict
comes from β (converge/iterate) and, at the boundary, from δ.

### Projection

The **typed interface by which a closed child cell becomes input to a parent
cell.** A child cell does **not** become a parent role; it *projects* matter,
verdicts, decisions, risks, and handoff upward. See **projection, not
role-renaming** (§3).

### Validator (V)

See §3. The predicate that checks a receipt/evidence against the contract and
emits a `PASS`/`FAIL` verdict.

### Delta (δ) / Epsilon (ε)

See §3. δ is the boundary authority (accept/reject/repair/override/block); ε is
the stream observer across receipts and cells.

---

## 3. Roles / Greek letters

The Greek letters are **roles** (and TSC axes), not folders and not renames of
each other. A role is a function a cell plays, not a directory it lives in.

### α / alpha — Pattern role

Produces matter and establishes the **shape** of the work (the implementer).

### β / beta — Relation / review role

Reviews matter against contract, evidence, and neighboring artifacts. Emits a
`converge` (advance) or `iterate` (send back to α) verdict. β iteration is an
**internal** cell loop, not an external lifecycle event.

### γ / gamma — Process / closeout role

Coordinates **closure**: receipts, handoff, and continuity. Authors the cell's
scaffold and lands the closeout artifact set.

### δ / delta — Boundary authority

Decides **what crosses the boundary**: `accept`, `release`, `reject`,
`repair_dispatch`, or `override`. An `override` accepts a non-`PASS` receipt and
closes the cell in **degraded state** (see **Degraded outcome**). δ never
rewrites V's verdict; an override is recorded separately.

Schema: `schemas/cdd/boundary_decision.cue`. Contract:
`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`.

### ε / epsilon — Stream observer

Reads patterns **across** receipts/cells and surfaces process-level
observations. ε watches the stream; it does not produce or gate a single cell's
matter.

Contract: `src/packages/cnos.cdd/skills/cdd/epsilon/`.

### V — Validator predicate

Checks the receipt/evidence against the contract and emits a **verdict**
(`PASS` | `FAIL`). V dispatches on the receipt's `protocol_id` and dereferences
the required evidence refs. δ acts on V's verdict; it does not overwrite it.

Schema: `schemas/cdd/boundary_decision.cue` (`#ValidationVerdict`).

### Projection, not role-renaming

The governing rule for how cells compose: a closed child cell **projects** its
outputs (matter, verdicts, decisions, risks, handoff) into a parent cell through
a typed interface. It does **not** get "promoted" or "renamed" into a parent
role. Composition is by projection across a boundary, not by relabeling a role.

### MCA / Wave δ (parent δ)

**MCA** (Most Coherent Action, §1) is the unit of autonomy a boundary authority
may take inside an approved contract. **Wave δ / parent δ** is the boundary
authority *for a wave*: it may take minor coherent actions inside the approved
wave contract and **must stop** on scope, authority, semantic, risk, or
release-boundary changes.

---

## 4. Issues and dispatch

GitHub Issues is the single source of truth for the backlog; labels make it
legible and dispatchable. Full model: `docs/development/issues/TAXONOMY.md`;
application rules: `docs/development/issues/TRIAGE.md`.

### Issue

A backlog item. A coherent issue carries exactly one primary `kind/*` and one or
more `area/*`. An **actionable** issue also carries one priority and a
`status:*`; `kind/tracking` issues are containers and may be priority-exempt.
Only a genuine executable cell also carries `dispatch:cell` + `protocol:*`.

### dispatch:cell

Marks an issue as a **cell a dispatch wake can execute**: scoped, with
acceptance criteria and a proof oracle. Applied **only** to genuinely executable
cells — never to a bare design/tracking/research issue.

### protocol:cds / protocol:cdd

The concrete protocol the cell dispatches through (`cds` = the software-
development realization of `cdd`). Applied **only** together with
`dispatch:cell`.

### status:* (lifecycle)

`status:ready` (shaped and approved, but held) · `status:todo` (claimable now) ·
`status:in-progress` (claimed, cycle running) · `status:review` (matter landed,
awaiting operator/δ) · `status:changes` (rejected; repair required). Only
`status:todo` is claimable by a dispatch wake, via the selector
`dispatch:cell + protocol:{p} + status:todo`. The operator flips
`status:ready → status:todo` to release a cell.

### kind/* · area/* · priority · effort/*

- **`kind/*`** — the primary split (exactly one): `bugfix`, `cleanup`,
  `process`, `feature`, `tooling`, `doctrine`, `audit`, `tracking`, `research`,
  `skill`, `spike`, `chore`.
- **`area/*`** — where the work lives (one or more): `agent`, `cdd`, `cds`,
  `core`, `runtime`, `docs`, `ci`, `cli`, `wake`, `coherence`, … (extensible).
- **priority** — `P0`–`P3`, or `priority/deferred`. One per actionable issue;
  `kind/tracking` containers may be priority-exempt. Orthogonal to kind/area.
- **`effort/*`** — ordinal size estimate (at most one): `S`, `M`, `L`, `XL`.
  Missing effort means *unestimated*, **not** `S`.
- **`resolution/*`** — on close when not a plain completion: `completed`,
  `superseded`, `duplicate`, `deferred`, `wontfix`.

---

## 5. Agent / wake terms

### Wake

One firing of an agent from its installed local files — the reconstitution of a
coherent agent at the start of a cycle (identity → doctrine → mindsets →
reflections → skills → capabilities → inbound). A coherent agent never wakes as
an empty chat window. Wakes are defined **as skills** (wake-as-skill): a package
ships a wake declaration and the renderer materializes the substrate artifact.

Defined in: `docs/reference/runtime/CAA.md`;
`src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`.

### Admin wake

The wake that handles channel sync, status reporting, and label routing. It
activates/attaches and **never executes cells**. Substrate:
`.github/workflows/cnos-agent-admin.yml`.

### CDS dispatch wake

The wake that **claims software-protocol cells** from the open-issue queue and
runs them via the δ role contract. Its inbound is the issue queue, not a home
thread; its selector is `dispatch:cell + protocol:cds + status:todo`; it claims
one cell per firing (FIFO by creation), invokes δ, and opens a PR.

Package: `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`. Substrate:
`.github/workflows/cnos-cds-dispatch.yml`.

### SKILL.md

The file that defines a skill (and, for orchestrators, a wake). A skill's
`SKILL.md` carries the DUR contract (Define / Unfold / Rules) and, for wakes,
`wake:` frontmatter (role, protocol, selector, `activation_state`,
`wake.output`). See **DUR**, **Skill**.

Location: `src/packages/{pkg}/skills/{category}/{name}/SKILL.md`;
orchestrators at `src/packages/{pkg}/orchestrators/{name}/SKILL.md`.

### wake.output

A wake's **output contract**. Admin-shape: channel-log convention + cursor.
Dispatch-shape: the cell-artifact set the cycle emits. It names what the wake is
accountable for producing, machine-readably.

Defined in: `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`.

### golden

The **renderer-emitted, reviewed, committed compiled artifact** a wake install
produces alongside the live workflow — `orchestrators/{wake}/cnos-{wake}.golden.yml`.
CI byte-diffs the golden against the live workflow on every change; any drift
fails the `install-wake golden` oracle. More generally, a *golden* is a
ring-fenced reference output that a validator compares against and that must not
be silently modified.

### live workflow

The **materialized substrate artifact** at the substrate's canonical path —
`.github/workflows/cnos-{wake}.yml`. It **MUST NOT be hand-edited**; it is
emitted by the renderer from the wake declaration and must match its golden.

### activation_state

Frontmatter enum on a wake: `live` means the renderer emits an installable
workflow. `declaration-only` means the wake is a contract specimen and the
renderer refuses to render it as a live wake.

---

## 6. Package / command terms

### Package

A distributable cognitive unit containing doctrine, mindsets, skills, commands,
and/or orchestrators. Source lives under `src/packages/{pkg}/`; installed into a
hub's `.cn/vendor/packages/`.

### cnos.core

The core package: doctrine (`CAP`, `COHERENCE`, `CBP`, `CA-Conduct`,
`AGENT-OPS`), mindsets, and core agent skills (activate, wake-provider,
dispatch-protocol). `src/packages/cnos.core/`.

### cnos.cdd / cnos.cds / cnos.issues

- **`cnos.cdd`** — the generic Coherence-Driven Development kernel: the
  recursive cell framework, the δ/β/γ/V role skills, receipt validation.
- **`cnos.cds`** — the **software** realization of CDD (the concrete engineering
  protocol) and its dispatch orchestrator.
- **`cnos.issues`** — the issue/backlog package boundary (taxonomy, triage,
  board tooling).

### Skill (module)

A bounded, selected, instrumental cognitive module under `skills/{category}/{name}/`
with a `SKILL.md` (DUR contract). Skills are situational amplifiers that compete
for bounded slots at wake-up — not identity. Doctrine never competes for skill
slots.

### Command

Executable machinery a package ships under `commands/` (e.g.
`src/packages/cnos.cdd/commands/cdd-verify/`), invoked through the `cn` CLI.

### Kernel

The generic, protocol-independent core of a system that specific protocols
overlay. E.g. the CDD generic receipt kernel (`schemas/cdd/`) that the CDS and
CDR receipts unify onto.

### cn

The capability CLI / runtime: the agent's body. The agent proposes typed ops;
`cn` validates against policy, executes within budget, records receipts, and
feeds evidence back. "The agent is the brain, `cn` is the body, git is the
nervous system." Reference: `docs/reference/cli/CLI.md`.

---

## 7. Trust / evidence

### Trust vs coherence

**Coherence** asks whether the work body hangs together with its stated model
strongly enough to be judged. **Trust** asks whether an actor or artifact can be
relied on, and what evidence supports that reliance. A coherent artifact can
still be wrong. A trusted actor can still produce incoherent work.

### Trust claim

An assertion that some artifact or actor warrants reliance — carried by
evidence (a signature, an attestation, a validated receipt), not by assertion
alone. A trust claim is only as good as the evidence bound to it.

### Signature vs attestation

- **Signature** — a cryptographic binding proving *who* produced or approved an
  artifact (authenticity/integrity of origin).
- **Attestation** — a statement *about* an artifact or process (what ran, what
  was checked, under what conditions). A signature says "I made this"; an
  attestation says "this is what happened."

### Coherence witness (TSC witness)

Evidence that a coherence measurement actually occurred — the recorded TSC
report (scores + provenance) that a receipt or ledger row points at. It is the
*witness* to a coherence claim, distinguishing a measured CΣ from an
asserted one.

### CDD receipt

See **Receipt** (§2). The trust/evidence view: the receipt is the durable object
against which **V** validates and on which **δ** decides — the unit that makes a
cell's outcome inspectable and transmissible later.

### Degraded outcome

The state a cell closes in when **δ records an `override`** — accepting a
non-`PASS` receipt anyway. `degraded_state` is pinned true on the override; the
cell is closed but flagged as degraded, with the original (non-`PASS`) verdict
preserved and the override rationale + authority recorded. A degraded outcome is
still transmissible, but it is honestly labeled as degraded rather than clean.

Schema: `schemas/cdd/boundary_decision.cue` (`#Override`).

---

## Appendix — doctrine, architecture, protocol (reference)

Terms that remain load-bearing but sit beneath the day-to-day work vocabulary.

### Doctrine

Always-on, constitutive cognitive substrate: first principles, conduct, runtime
grammar, review law. Loaded at every wake; never competes for skill slots. Four
layers: **CAP** (dynamic atom), **COHERENCE** (review geometry), **CBP +
CA-Conduct** (relational boundary), **AGENT-OPS** (runtime grammar).
`src/packages/cnos.core/doctrine/`.

### DUR (Define / Unfold / Rules)

The canonical skill contract every `SKILL.md` follows: **Define** the parts,
**Unfold** each with ❌/✅ pairs, state numbered **Rules**. DUR is to skills what
TSC is to coherence — the structural invariant that makes the class recognizable.

### CLP (Coherence Ladder Process)

The review rhythm that keeps CAP from drifting: seed the gap → Bohmian
reflection → triadic check (α/β/γ) → patch the *weakest* axis → repeat. Never
publish cold. `src/packages/cnos.core/doctrine/COHERENCE.md`.

### CAA (Coherent Agent Architecture)

The design document specifying what a coherent agent is structurally: doctrinal
layers, cognitive strata at wake-up, the agent loop, invariants, failure modes.
`docs/reference/runtime/CAA.md`.

### CN Shell / N-pass bind loop

The capability runtime that mediates agent↔world: the agent proposes typed ops;
the shell validates against policy, executes within budget, records receipts.
The **N-pass bind loop** enforces observe-before-effect across bounded passes
(observe → effect → terminal). `docs/reference/runtime/AGENT-RUNTIME.md`.

### CAR (Cognitive Asset Resolver)

The package distribution system: how cognitive assets are packaged, versioned,
installed, and resolved locally so wake-up is deterministic.
`docs/architecture/cognitive-substrate/CAR.md`.

### CN (Coherence Network) / Hub / Peer / Thread

- **CN** — a network of agents using git repositories as the primary surface for
  specs, state, and messages. Git is the transport; files are the state.
- **Hub** — a git repository that is an agent's home (identity, state, installed
  packages). Wake-ready from its own files.
- **Peer** — another hub this hub tracks and syncs with.
- **Thread** — a Markdown conversation/reflection file **in a hub** (under that
  hub's `threads/`). This is a *hub-side* concept; the cnos source repo has no
  root `threads/`.

Protocol: `docs/reference/protocol/cn/PROTOCOL.md`.

### cnos (Coherence Network OS)

The source repo and runtime for CN agents: the `cn` CLI, cognitive packages
(`src/packages/`), design docs (`docs/`), schemas, and tests. cnos is the
current lowest durable software articulation of the coherence system; agents
consume it through installed packages, not a checkout dependency.

---

*Canonical spec for the doc system these terms live in:*
`docs/reference/governance/DOCUMENTATION-SYSTEM.md`. *Naming conventions:*
`docs/reference/governance/NAMING.md`.
