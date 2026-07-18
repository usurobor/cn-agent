# Cell Runtime Classes — Contracts, FSM, and Routing

**Status:** Draft (realization layer). Not ratified. Formalized by a Planning Cell in PC-D0 mode (cnos#662), under the operator-authorization comment of 2026-07-13T19:14:52Z. Ratification is **not** this cell's job: the exit sequence is external-β review → γ binds SHA_M + review + evidence → a separate Cohering Cell (CC) ratification → operator-final-read → merge of this D0 → a **separate PC-Wave** → CC wave review → operator wave authorization — see the **Exit sequence** block in §16.
**Owns:** How the WC/PC/CC output-telos classes named by `docs/architecture/CELL-RUNTIME.md` (#628) route through a common cell contract envelope, class-specific `V` predicates, the CC↔ε lineage, the cell FSM and the (not-yet-shipped) wave FSM, the command surface, and the human-gate/wake-topology policy that makes class-aware dispatch mechanical.
**Does-not-own:** The kernel algorithm (`COHERENCE-CELL-NORMAL-FORM.md`, #370), the WC/PC/CC class definitions themselves (`CELL-RUNTIME.md`, #628 — cited throughout, never restated), the receipt rule and four surfaces (`COHERENCE-CELL.md`), the CUE schema implementation (`schemas/cdd/`), or any Go runtime/FSM code. This note **names** schemas and FSM extensions (D9); it does not implement them.

---

## Thesis

`CELL-RUNTIME.md` established that CNOS has one CCNF kernel, deployed in three output-telos shapes — WC (artifact), PC (relation graph), CC (process judgment) — executed by one generic runner. That note stops at the class definitions. It does not say what a cell contract looks like as a typed object, what a Cohering Cell consumes and emits, how `cell_class` interacts with the shipped label-driven FSM, or which parts of the "mechanical FSM" sketch are running in production today versus still design. This note is the **operationalization layer**: it takes the classes as given and specifies the routing, contracts, and FSM extension that let a generic runner dispatch them without inventing policy at runtime. It carries ten operator-pinned decisions (D1–D10, §16) as settled input, not as open questions to re-derive.

**What this note bootstraps, and the identity model it assumes.** The larger objective is to **bootstrap a three-cell agent** — Working / Planning / Cohering — that operates this repository autonomously. That mechanical three-cell runtime does **not** exist yet (no `cell_class`-aware routing, no wave FSM, no per-class `V` validators, no scheduled CC pulse); this note *specifies* it. Throughout, two levels are kept distinct:

- **"Sigma" is an activation / agent lineage** — one model identity operating this repo.
- **κ, α, β, γ, δ and the cell classes WC / PC / CC are functional roles.** One activation can carry several roles; a role is not an identity.

This distinction frames the κ/α discussion in §8. The note's target — **State B**, where α is a distinct runtime-executed role and κ sits outside the cell — is the design and holds **unconditionally**. During bootstrap the operator **authorized κ (Sigma) to execute the α role** (State A) only because the mechanical Planning-Cell runtime that would execute α does not exist yet; that is an **operator-authorized, explicitly transitional posture retired by the very migration this note specifies**, not a licensed permanent way to run a Planning Cell (§8).

---

## 1. Relationship to CELL-RUNTIME.md and the kernel

This note is a **sibling** of `docs/architecture/CELL-RUNTIME.md`, not a replacement or a restatement. The relationship is strict:

| Surface | Owns | This note |
|---|---|---|
| `COHERENCE-CELL-NORMAL-FORM.md` (#370) | The five-step kernel `α→β→γ→V→δ`, the four cell outcomes, the two recursion modes, the three scope-lift projections | Cited; §4 below is a direct instance of Projection 3 (ε) |
| `CELL-RUNTIME.md` (#628) | WC/PC/CC as output-telos classes of one kernel; the four orthogonal axes (kernel cell / roles / output-telos / matter domain); CM measures, V gates, δ effects | Consumed as fixed input; this note does not re-derive the classes |
| `CELL-KINDS.md` (#570) | The matter/contract domain vocabulary (axis 4) | Domain vocabulary unaffected; `cell_class` (this note) and `matter_domain` (`CELL-KINDS.md`) are orthogonal fields on the same contract. Its κ/α role note (`:137`) carries the State-A→State-B migration dependency (§8, §14). |
| **This note** (#662) | The cell contract envelope; class-specific `V`; CC↔ε; the distinct-`CellClass` FSM extension point; the unconditional κ≠α invariant (State B) + the operator-authorized transitional bootstrap posture (State A, §8); the wave FSM sketch | — |

**D1 (one CCNF kernel).** WC, PC, and CC are output-telos classes of the *same* kernel — already `CELL-RUNTIME.md`'s own thesis. This note confirms it and does not re-derive it: every class below still runs its own internal, independent `α→β→γ→V→δ`; nothing here forks the kernel, adds a role letter, or changes α/β/γ/δ semantics.

---

## 2. The cell contract envelope (`cn.cell.contract.v1`)

All three classes consume one typed envelope. Only *input admissibility, matter type, class-specific acceptance, and post-closure routing* differ by class (§3, §5). Per **D9 (schema-first destination)**, this note **names** the schema; it does not implement it — no CUE, no Go struct, no validator. The two blocks below share **one canonical set of key paths**: a placeholder template and an illustrative worked instance that **normalize to identical key paths**, with the worked instance **conforming to the type/cardinality constraint model stated below**. An ordinary value-parse of the two blocks proves only **key-path identity**; the **types, cardinalities, and unions are enforced by a constraint checker** (future CUE/Go per D9), not by the parse. To keep the key-path property machine-checkable under an ordinary parser, the template uses **quoted scalar placeholders and comment annotations only** — no bare commas inside flow scalars:

```yaml
schema: cn.cell.contract.v1               # const
cell:
  id: "<issue>"                           # scalar cell id
  class: "working"                        # enum: working | planning | cohering
  mode: "<class-specific>"                # scalar (e.g. d0 | wave for planning)
  protocol: "cds"                         # scalar
  matter_domain: "<CELL-KINDS.md-vocab>"  # scalar
scope:
  repo: "<owner/name>"                    # scalar
  wave: "<parent-wave-issue-or-null>"     # scalar-or-null
  parent_cell: "<parent-cell-issue-or-null>"  # scalar-or-null
intent_ref:
  schema: "cn.intent.v1"                  # const
  id: "<intent-id>"                       # scalar
  carrier:
    kind: "github_issue"                  # scalar
    ref: "<issue>"                        # scalar
inputs:
  required: [ "<dependency-ref>" ]        # seq<scalar>, 1+
  optional: [ "prior_receipts" ]          # seq<scalar>, 0+
requested_output:
  kind: "<artifact | relation_graph | judgment>"   # enum scalar
  path: "<where-the-matter-lands>"        # scalar
acceptance:
  predicates: [ "<class-specific-predicate>" ]     # seq<scalar>, 1+ (§5)
constraints:
  allowed_paths:   [ "<surface-glob>" ]   # seq<scalar>, 1+
  forbidden_paths: [ "<surface-glob>" ]   # seq<scalar>, 0+
  non_goals:       [ "<explicit-non-goal>" ]        # seq<scalar>, 0+
gates:
  operator_authorization_required: false  # bool
  operator_acceptance_required: false     # bool
  reason: null                            # scalar-or-null; required-nonempty IFF a gate bool is true, else null/absent (here both false ⇒ null)
doctrine_affecting: false                 # bool — the doctrine-gate discriminator (§9)
stop_conditions: [ "<typed-STOP-trigger>" ]         # seq<scalar>, 0+
```

**Normative constraint model (key paths, types, cardinalities) — authoritative.** This paragraph, not the illustrative YAML, is the authoritative constraint model: the contract is exactly these paths with these types and cardinality bounds, and nothing normalizes to any other shape: `schema` (const) · `cell.{id, class, mode, protocol, matter_domain}` (scalars; `class` an enum) · `scope.{repo, wave, parent_cell}` (`wave`/`parent_cell` scalar-or-null) · `intent_ref.{schema, id, carrier.{kind, ref}}` (scalars) · `inputs.{required (1+), optional (0+)}` (seq<scalar>) · `requested_output.{kind (enum), path}` (scalars) · `acceptance.predicates` (seq<scalar>, 1+) · `constraints.{allowed_paths (1+), forbidden_paths (0+), non_goals (0+)}` (seq<scalar>) · `gates.{operator_authorization_required, operator_acceptance_required}` (bool) with `gates.reason` (scalar-or-null; see gate invariants) · `doctrine_affecting` (bool) · `stop_conditions` (seq<scalar>, 0+). A future CUE/Go **constraint checker** pins and enforces these types/cardinalities/unions; an ordinary value-parse proves only key-path identity. Verification for this section is **two fixtures**: (i) a key-path diff — parse both blocks, diff their normalized key-path sets, which must be identical; (ii) a constraint-conformance check — the worked instance validated against this table by the future checker.

**Gate invariants (decidable).** The `gates` / `doctrine_affecting` fields satisfy a small decidable invariant set the constraint checker enforces (a truth table over three booleans plus `reason`):
1. `doctrine_affecting: true` ⟹ `gates.operator_acceptance_required: true` — a doctrine-affecting matter always requires operator acceptance.
2. `gates.reason` is **required-and-nonempty iff** at least one gate boolean (`operator_authorization_required`, `operator_acceptance_required`) is `true`; otherwise it is **absent or null**.
3. For the doctrine gate, `doctrine_affecting` is the **authoritative** field and `gates.operator_acceptance_required` is **derived** from it (invariant 1); a checker rejects any contract where `doctrine_affecting: true` but `gates.operator_acceptance_required: false`.

The template above carries both gate booleans `false` (and `doctrine_affecting: false`), so `reason` is `null`; the worked instance below carries both gate booleans `true` (and `doctrine_affecting: true`) with a nonempty `reason` — both satisfy the invariants.

**Intent is referenced, not inlined (D9, §13).** The contract does not *contain* intent; it **references** a first-class `cn.intent.v1` object (§13) via `intent_ref: { schema, id, carrier }`. A GitHub issue is a **carrier / control-plane projection** of that intent — `carrier: { kind: github_issue, ref: <issue> }` — **not** intent identity. Per the D9 dependency order (`cn.intent.v1 → cn.cell.contract.v1 → …`), a durable `cn.intent.v1` object is produced by κ *before any cell exists*, and the cell contract points back at it. **`cn.intent.v1` is named-not-implemented in this note (D9): no durable intent object exists on a pinned surface today.** Materializing a durable `cn.intent.v1` object (and making `intent_ref.id` resolve to it immutably) is therefore a **schema-implementation deliverable** owned by the schema/PC-Wave work (§13, S2 of #627 / #644), not a claim of this Draft.

**Doctrine gate lives in the contract (D5, §9).** The doctrine gate is determinable from this canonical contract *alone*: `matter_domain` plus the top-level `doctrine_affecting` boolean (and, when a gate fires, `gates.reason`) are the discriminator. There is no companion operator-comment block that the mechanical gate selection must also read — the contract is the single authority surface (§9).

**This worked instance is illustrative, not a resolved live object.** It shows the shape **#662's own cell contract takes once `cn.cell.contract.v1` and `cn.intent.v1` are implemented** — it is *not* an assertion that a durable `cn.intent.v1` with id `intent-2026-0711-662` exists in the repository today (it does not; see the paragraph above and §13). It carries the same key paths as the template and conforms to the constraint model above:

```yaml
schema: cn.cell.contract.v1
cell:
  id: "662"
  class: "planning"
  mode: "d0"
  protocol: "cds"
  matter_domain: "doctrine"
scope:
  repo: "usurobor/cnos"
  wave: "627"
  parent_cell: null
intent_ref:
  schema: "cn.intent.v1"                  # illustrative id — not yet a durable object (§13)
  id: "intent-2026-0711-662"
  carrier:
    kind: "github_issue"
    ref: "cnos#662"
inputs:
  required: [ "docs/architecture/CELL-RUNTIME.md", "cnos#627", "cnos#530", "cnos#644" ]
  optional: [ "prior_receipts" ]
requested_output:
  kind: "artifact"
  path: "docs/architecture/CELL-RUNTIME-CLASSES.md"
acceptance:
  predicates: [ "wc_pc_cc_are_one_kernel_classes", "cc_epsilon_reconciliation_preserved", "fsm_states_and_events_specified", "kappa_outside_cell", "current_and_target_state_distinguished" ]
constraints:
  allowed_paths:   [ "docs/architecture/CELL-RUNTIME-CLASSES.md", ".cdd/unreleased/662/**" ]
  forbidden_paths: [ "src/go/**", ".github/workflows/**" ]
  non_goals:       [ "implement_cell_runner", "alter_ccnf", "add_new_role_letters", "replace_cdd_or_cds", "file_or_dispatch_child_issues" ]
gates:
  operator_authorization_required: true
  operator_acceptance_required: true
  reason: "doctrine-affecting architecture matter; operator final-read required"
doctrine_affecting: true
stop_conditions: [ "source_doctrine_conflicts", "required_decision_unresolved" ]
```

Both blocks carry the same key paths at the same nesting: `cell.id` (not `cell.issue`); `requested_output` an **object** `{ kind, path }` (not a list); `non_goals` **under `constraints`** (not `cell.non_goals`); `gates` / `doctrine_affecting` / `stop_conditions` **top-level** (not `cell.*`). The template's placeholders are quoted scalars, so no comma-in-flow-scalar splits a value into extra keys or elements — the two blocks are one shape, not two. (`matter_domain: doctrine` + `doctrine_affecting: true` are the doctrine-gate pair, both now canonical contract fields — §9.)

**Labels vs typed contract.** Per **D6/D9** and #643 (shipped): GitHub **labels** carry only coarse mechanical state — `dispatch:cell`, `protocol:cds`, `status:*` — and answer *is this dispatchable, which protocol, where in the lifecycle*. The **typed contract** (above) carries semantics — `cell.class`, `mode`, `matter_domain`, `requested_output`, `acceptance`, `constraints`, `gates`, `stop_conditions` — and is *not* itself a label. #662 was authorized without a `cell_class:*` label for exactly this reason (operator-authorization comment, 2026-07-13): "the issue contract remains the source" for `cell.class`.

---

## 3. Class contracts: WC, PC, CC

Each class is the *same* kernel cell with a different output telos (`CELL-RUNTIME.md` §5, restated here only enough to anchor the contract-level detail `CELL-RUNTIME.md` leaves implicit).

### 3.1 Working Cell (WC)

*Telos:* increase the implemented repository — matter targets pattern (TSC-α), output is an artifact.

**Input admissible only if executable:** operator-authorized, exactly one protocol, dependencies satisfied, concrete requested output, testable acceptance criteria, allowed/forbidden surfaces explicit, human gates named, no unresolved planning decision, STOP conditions present, claimable lifecycle state (`status:todo`).

**Matter:** code, tests, docs, schemas, fixtures, workflow/config, migration patches. **α** produces the implementation diff plus self-coherence evidence; **β** verifies every acceptance criterion, no forbidden-surface touch, claims match code, tests cover contract predicates, no hidden behavior, no scope creep; **γ** closes with a receipt and triage.

**Result shape:** `{ class: working, matter_ref, review_verdict, validation_verdict: PASS|FAIL, receipt_ref, proposed_next_state: status:review }`.

### 3.2 Planning Cell (PC)

*Telos:* increase executable relational structure — matter targets relation (TSC-β), output is a relation graph. Two modes, not one contract:

- **PC-D0** — `operator intent → product/architecture/spec contract`. Admissible if operator intent is captured, the problem boundary is named, existing state is linked, unresolved product/architecture questions are made explicit, the requested artifact is specified, and **no implementation is requested.** Output passes only if it establishes: who the user/operator is; what relation must become true; State A (current); State B (target); primary surface and fallbacks; authority/permission boundaries; explicit decisions; genuinely open questions; non-goals; and criteria to begin wave planning. **#662 is itself the first bootstrap instance of a PC-D0 contract** (operator-authorization comment, 2026-07-13) — see §11's bootstrap-calibration note.
- **PC-Wave** — `validated D0 → executable wave`. Admissible only once the D0 is validated, the target is sufficiently defined, scope/non-goals are stable, major decisions have owners, and surfaces are known. Output passes only if: the graph is acyclic; **every node embeds or immutably references a complete `cn.cell.contract.v1` (§2)** — the full envelope, carrying `class` + `protocol` + `matter_domain` + `requested_output` + acceptance predicates + input/output refs (dependencies) + gates + `doctrine_affecting` + STOP conditions (a node does **not** re-declare a partial subset of those fields); **wave graph edges are the mechanical projection of producer-output→consumer-input relations across the child contracts, not an independent dependency authority** (§13); the critical path is explicit; parallel cells don't share write surfaces; human gates are located; **no child is auto-dispatched**; STOP conditions exist at cell *and* wave level; the **child-completion and whole-wave-completion predicates (§13) are defined**. The child's `cn.cell.contract.v1` is the **single canonical authority** for its class/protocol/dependencies; the wave graph's `depends_on` edges are **derived from** (or validated for exact parity against) the child input/output refs — the graph never introduces a dependency the child contracts do not carry, and never omits one they do (§13, `cn.wave.v1`). #644 is the open implementation target for PC-Wave mechanization; **this note specifies its contract, it does not build it.**

**Guard, both modes:** a Planning Cell creates and refines its own children; it does not apply `status:todo` to them. Dispatching a PC's own children requires operator authorization (or a typed policy gate at the wave-authorization boundary, §9) — never an automatic action taken by the PC itself. This is the same guard #662 itself operates under: it produces exactly one artifact and files or dispatches nothing.

**Result shape — a tagged union by `mode`, not one shape.** PC-D0 produces a spec/architecture artifact; PC-Wave produces a wave graph. Forcing `wave_ref`/`graph` onto a D0 result would quietly turn every Planning Cell into a wave-producer, which it is not:

```yaml
# mode: d0
{ class: planning, mode: d0,   artifact_ref, readiness: ready_for_coherence_review, requires_operator_gate: true }
# mode: wave
{ class: planning, mode: wave, wave_ref, graph: { nodes, edges }, readiness: ready_for_wave_review, requires_operator_gate: true }
```

A PC-D0 result carries an `artifact_ref` and is `ready_for_coherence_review` (its next step is a Cohering Cell review, then operator-final-read — per the **Exit sequence** block in §16). Only a PC-Wave result carries `wave_ref` + `graph` and is `ready_for_wave_review`. Both are `requires_operator_gate: true`; neither is self-dispatchable. #662, a PC-D0, produces the `artifact_ref` form.

### 3.3 Cohering Cell (CC)

*Telos:* increase justified confidence about state and select the next coherence-preserving action — matter targets process (TSC-γ), output is a coherence judgment plus exactly one next disposition. CC does **not** write implementation or planning matter; it **judges** state and, where more work is needed, names the class of cell that should do it (§4 sharpens the ε relationship this telos depends on).

**Input:** operator intent, scope/issue/PR/CI state, wave graph, receipts, decisions, open risks, **ε cross-cell observations** (§4), prior CC receipt. Admissible if: scope is named, state is mechanically snapshottable, relevant receipts are available, an intent/pulse reason is present, the observation window is explicit, no other CC is active for the same scope/version, and the disposition vocabulary is known.

**Matter:** coherence assessment, gap classification, bottleneck axis, contradiction map, wave-progress judgment, residual risk, next-MCA recommendation (MCA = Minimum Coherence Action).

**Output passes only if it:** names the evaluated scope; names the state snapshot; cites evidence per load-bearing claim; separates observation from inference; identifies the dominant bottleneck; emits **exactly one** primary disposition; emits a typed next-MCA when more work is needed; performs **none** of the recommended work itself; distinguishes `done` from `done-with-residuals`; requests human action only when the mechanics genuinely cannot decide.

**Dispositions:** `request_planning | request_working | hold | request_human | continue_wave | complete | complete_with_residuals | block`.

**Result shape — a tagged union by `disposition`, not one overloaded shape.** A single result carrying `next_mca` *and* `requires_operator_gate` *and* an optional reason lets impossible combinations (e.g. `complete` with a `next_mca`, or `request_human` with no gate) stay contract-valid. Instead, every CC result carries a common head and exactly one disposition tag; the tag fixes which payload fields are **required** and which are **forbidden**, so `V` can prove *exactly one coherent disposition*:

```yaml
# common head (every tag): { class: cohering, judgment, disposition, requires_operator_gate }
```

| Disposition tag | Requires | Forbids |
|---|---|---|
| `request_planning`, `request_working` | `next_mca` (typed `cn.next-mca.v1`, names the class) | `residuals`, `gate` |
| `continue_wave` | `next_mca` (the continuation MCA) | `residuals`, `gate` |
| `hold` | `reason` | `next_mca`, `residuals` |
| `block` | `reason` | `next_mca`, `residuals` |
| `request_human` | `gate` (typed human-gate, §9) | `next_mca`, `residuals` |
| `complete` | — | `next_mca`, `residuals`, `gate` |
| `complete_with_residuals` | `residuals` (typed residual list) | `next_mca`, `gate` |

CC recommends; `V` validates the judgment's structural shape **and** that the emitted tag's required/forbidden field contract holds (positive and negative fixtures: reject `complete` + `next_mca`, reject `request_human` without a gate); the FSM (not CC) applies any resulting transition (**D7**, §7). This tagged union is the **CC class-result shape** — a refinement of how a CC result embeds `cn.next-mca.v1` (§13), pinned during implementation of the four-schema surface; it is **not** a fifth canonical schema.

**Anti-patterns, both directions (`CELL-RUNTIME.md`'s runtime-loop invariant, restated as a guard):** CC writing the architecture note directly blurs CC into PC. PC turning vague intent directly into implementation issues without a D0/product note skips the customer-facing contract PC-D0 exists to produce (this is exactly what #654 hardens against).

---

## 4. CC ↔ ε — settled (D3, F1)

**Carried verbatim per the operator's authorization comment and the issue's own §CC↔ε reconciliation** (not paraphrased, per that section's explicit instruction):

```text
ε is the cross-cell receipt-stream observer / projection function in CCNF.
A Cohering Cell is not ε itself and does not add a new role to the cell closure.
A Cohering Cell is the cell-class realization at scope n+1 that consumes the εₙ
receipt-stream projection as input and produces coherence judgment / next MCA as its matter.

So:
  εₙ:      observes receipt streams across cells at scope n and projects
           protocol/coherence observations upward.
  CCₙ₊₁:   is a normal CCNF cell with internal α→β→γ→V→δ closure;
           its output-telos is coherence judgment;
           its matter may be produced from εₙ observations, current state,
           operator intent, receipts, issues, PRs, and artifacts.

One lineage, not two competing coherence surfaces:
  ε supplies the cross-scope observation.
  CC realizes that observation as an executable cell class.
```

This matches CCNF's own §Scope-Lift Projection 3 (`εₙ receipt-stream observation → γₙ₊₁-like coordination/evolution`, `COHERENCE-CELL-NORMAL-FORM.md`) and its explicit projection-not-renaming framing: `εₙ` is not literally `γₙ₊₁`, and CC is not ε wearing a different hat — CC is a full cell at scope `n+1` with its own α/β/γ/δ, for which εₙ's observation is one *input*, not the cell itself. **#530 owns the typed scope-lift / ε-projection artifact.** This note depends on that lineage; it does not redefine it, and it files no work against #530.

---

## 5. Class-specific `V` (D8)

`V : Contract × Receipt → Verdict` stays the one typed predicate CCNF pins at step 4 (`PASS|FAIL`, no `WARN` — `schemas/cdd/receipt.cue`, `#ValidationVerdict`). Per **D8**, its *predicate content* differs by class, layered over a common kernel floor:

| Layer | Predicates |
|---|---|
| **Common `V`** (all classes) | Receipt complete · evidence present · role ownership valid · contract fields satisfied · no unresolved blocking finding |
| **WC `V`** adds | Executable acceptance criteria pass · tests/checks pass · allowed-surface and non-goal guards hold |
| **PC `V`** adds | Requested planning artifacts exist · relation graph acyclic · every child cell has ACs + dependencies + gates + STOP conditions · unresolved decisions explicitly gated (not hidden) · **no child auto-dispatched** (mechanically checkable: the PC applied no `status:todo` to any child) |
| **CC `V`** adds | Judgment is evidence-backed · exactly one disposition emitted · the emitted disposition tag's required/forbidden field contract holds (§3.3 tagged union — e.g. `complete` carries no `next_mca`, `request_human` carries a typed gate) · **no implementation surface modified** (mechanically checkable as an allowed-surface check on the CC's own matter: matter paths ⊆ judgment artifacts, no code/product diff) |

These predicates are the seed for the future per-class `cell.cue` validators (`CELL-RUNTIME.md`'s Reconciliation map, `#369`/#627 S2–S3). This note names them as the destination `V` extends toward; it implements none of them.

---

## 6. `cell_class` as an FSM dimension (D6, F2)

The runtime eventually evaluates the triple **`(status, protocol, cell_class)`**, not just `(status, protocol)`:

- **`status`** (label-based) answers *where is this cell in its lifecycle?*
- **`protocol`** (label-based) answers *which package's semantics apply?* (`cds`, and eventually `cdr`/others)
- **`cell_class`** (contract-based) answers *what kind of cell is it — what authority and guards apply?* Can it create child issues? Can it dispatch? Can it modify implementation surfaces? Can it emit a next-MCA?

**Dispatch readiness stays solely label-gated:** `dispatch:cell + protocol:{P} + status:todo` — exactly the shipped #643 rule. `cell_class` is **post-claim** routing/authority metadata sourced from the typed contract, not a second readiness gate. If a claimed cell's contract is missing `cell.class`, `requested_output`, or acceptance criteria, the runtime routes it to `status:blocked` with `degraded_reason: cell_contract_incomplete` — a post-claim validation failure, not a pre-claim admission failure. This answers `CELL-RUNTIME.md`'s Open-Q1/Open-Q2 without contradicting #643 or #640.

**Grounding in the shipped seam — a distinct field, not a repurposing.** The shipped FSM carries one observation-only seam: `FactSnapshot.CellKind{Observed, Source, DefaultedTo}`, populated from the `cell_kind:` line of `gamma-scaffold.md` (`src/packages/cnos.issues/commands/issues-fsm/fetch.go`) and locked as *observation, not enforcement* by `TestSeam_CellKindNotEnforced`. **That seam observes `matter_domain`, not class** — `CELL-KINDS.md` §"FSM awareness" is explicit that the legacy `cell_kind` field *is* `matter_domain`. This note therefore does **not** repurpose `CellKind` to evaluate `cell_class`; doing so would collapse the two orthogonal axes (`CELL-RUNTIME.md` Open-Q1: *"add both, keep `cell_kind` as a compat alias"* — add both fields, do not overload one).

The specified extension is a **dual-field migration**:

- **`CellClass`** — a **distinct, new** observation/evaluation field, sourced from `cell.class` in the typed contract (§2). Class guards (authority: can it create/dispatch children, modify implementation surfaces, emit a next-MCA?) evaluate `CellClass`.
- **`CellKind`** — retained **unchanged** as the `MatterDomain` compatibility adapter (legacy alias). Domain refinements continue to evaluate `MatterDomain` (legacy `CellKind`); its observation is *not* lost.

The two fields are populated and evaluated independently, so model tests can vary class and domain orthogonally and prove class guards read `CellClass` while domain refinements read `MatterDomain`. This is the deferred "FSM Phase 2" `CELL-KINDS.md` names — realized as *two* fields, not a promotion of one into the other. The exact source mechanics of `CellClass` (contract field vs. a new observation seam vs. eventual label) are for schema work to pin (§13); the settled decisions here are that `cell_class` is an evaluated FSM dimension **and** that it gets its own field distinct from the matter-domain seam.

---

## 7. The Coherence Loop

```text
κ captures operator intent
  → CC: is there enough coherent planning matter to proceed?
      no  → PC-D0 produces a product/architecture/spec contract
  → CC judges the D0
      coherent → PC-Wave produces an executable graph
  → CC judges wave dispatchability
  → operator authorizes the wave (§9 — wave-boundary gate, not per-child)
  → WC implements one bounded contract
  → CC consumes the WC receipt + updated state → continue | replan | hold | complete
```

Formally: `Jₖ := CC.evaluate(I₀, Stateₖ, Receiptsₖ)`; PC/WC produce on `Jₖ.next_mca`; the loop terminates only when CC emits `complete | complete_with_residuals` **and** `V`/δ/operator accept. This is `CELL-RUNTIME.md` §"Runtime Loop" made executable — its two anti-patterns (CC writing the note directly; PC skipping the D0) are the guards named in §3.2/§3.3 above.

---

## 8. κ remains outside the cell; κ ≠ α (D2, F3) — with an operator-authorized transitional bootstrap posture

κ is a **control-plane** slot, not a cell role — it does not appear beside α/β/γ/δ/ε in CCNF, and this note adds no new role letter. The invariant, stated **without qualification**, is the design:

> **κ supplies operator/control-plane input; the cell's α owns and writes the specification or implementation matter. κ does not author α matter, does not perform β or γ, and role functions must not be silently collapsed.** κ speaks through issues and comments; cells work through typed contracts and receipts.

This is **State B** — the target the three-cell agent is being built toward — and it is **not softened, relativized, or made optional by the bootstrap.** It aligns with #583/#584 (landed doctrine: the mechanism/cognition boundary — mechanical dispatch separated from the cognitive skills that run inside a claimed cell).

**State A — the operator-authorized transitional bootstrap posture.** The mechanical Planning-Cell runtime that would execute α **does not exist yet**. Because there is no runtime to run α mechanically, the operator **authorized κ (Sigma) to execute the α role of the Planning Cell during bootstrap.** This is State A: κ = α **by operator authorization** — a legitimate, explicitly transitional posture, a designed stage on the path to State B, **not** a nonconformance, **not** a violation, and **not** a second valid role topology licensed to persist. It is **retired the moment the Planning Cell runs mechanically** — its own α executed by the runtime — at which point the κ/α collapse is gone and κ≠α holds unconditionally. Because State A is transitional **by construction**, it does **not** soften, relativize, or make optional the State-B invariant above: authorized-transitional is a bootstrap stage, not a licensed *permanent* topology. Using the identity model from the Thesis ("Sigma" = activation lineage; κ/α/β/γ/δ + WC/PC/CC = functional roles): the migration from State A to State B **is** the operationalization of the three-cell (WC/PC/CC) agent — once the Planning Cell runs mechanically, its own α is executed by the runtime. Authorization is **not** self-warrant, and it does **not** relax the firebreak: **actor collapse** — one activation performing two cell roles inside one cell boundary (e.g. acting as α and β at once), the case that defeats CCNF's β-independence firebreak — remains **forbidden** in State A exactly as in State B. What the bootstrap actually carries is the weaker **hosting-identity collapse**: the separate α/β/γ activations run under one shared Sigma lineage (separate activations, one hosting identity), tracked by **#664** — a distinct, weaker limitation than actor collapse, and the reason the State-A *independent* warrant is reduced.

**Provenance / warrant of this bootstrap.** Operator authorization of the State-A posture settles that κ=α is *permitted* during bootstrap; it does **not** by itself supply the cell's *independent* warrant. Because any *internal* β pass is the same Sigma lineage (the hosting-identity limitation tracked by **#664**), **the internal passes do not by themselves constitute independent-β warrant.** The independent warrant for this matter is the **external β — a review entity outside the Sigma lineage — bound to the matter-only R3 commit SHA it reviewed (§11.6)**, plus the later γ receipt that binds that same SHA and the external review. The reduced *independent* warrant of the authorized State-A posture is discharged by that external review, not by any internal attestation: authorized ≠ self-warranting.

**`CELL-KINDS.md:137` is State-A-accurate, scoped, and superseded at State B.** `CELL-KINDS.md:137` ("κ … commonly acts as **α of a planning cell**") **correctly describes the authorized State-A bootstrap posture**: during bootstrap, κ (Sigma) does execute the α role of the Planning Cell, exactly as line 137 says. It is **accurate for State A and scoped to it** — and it is **superseded at State B**, where the Planning Cell runs mechanically, its own α is runtime-executed, and κ≠α holds, so line 137 no longer describes the operating topology. Reconciling that prose — marking line 137 as State-A-bootstrap wording and repointing it to the operational, mechanically-separated α role — is bound to a **bounded downstream doctrine-migration issue** owned by the PC-Wave / later doctrine wave; that issue **repoints and scopes** the line, it does not correct an error. **This PC-D0 does not edit `CELL-KINDS.md`** (§14 records the dependency and owner; §15 non-goal).

**Guard consequence (§12):** "κ editing cell-owned artifacts" is a forbidden action precisely because κ ≠ α/β/γ/δ. κ carries intent into control-plane artifacts (issues, labels, PR-review surfaces); it does not produce, review, or close a cell's matter. In State A, Sigma executes the α role of the Planning Cell under explicit operator authorization — an **operator-authorized transitional posture disclosed on the receipts**, not a relaxation of the guard and not a disclosed nonconformance.

---

## 9. Human gates (D5)

Human authority sits at **irreversible or scope-expanding boundaries**, not between internal phases:

**Default gates (always human):** intent acceptance · wave authorization (first dispatch of a wave, not every child in it) · production/release boundary · final acceptance of doctrine-affecting or otherwise system-doctrine-changing matter · explicit hold/block/escalation resolution. **The doctrine gate is determinable from the canonical `cn.cell.contract.v1` (§2) alone** — it fires off the contract's own **typed fields**: `matter_domain` plus the top-level `doctrine_affecting` boolean (with `gates.reason` recording why), exactly as §2's worked instance carries (`matter_domain: doctrine`, `doctrine_affecting: true`). There is **no companion operator-comment contract block** the gate selection must also read; the discriminator lives in the contract. Mechanical systems consume these declarations — they never classify prose to decide whether a gate applies.

**No default gate between:** `CC → requested PC/WC` · PC's internal α/β/γ · WC's internal α/β/γ · safe mechanical recovery (stale-claim resume, repair-dispatch re-entry). A wave, once operator-authorized at its boundary, may pre-authorize its own dependency-respecting internal nodes — `operator → authorize wave → wave executable → FSM routes children by dependencies without further operator intervention`, re-entering only at a named boundary (human gate, blocked, review, release). **The operator must not become the scheduler for every child** — that failure mode is exactly what wave-boundary (not per-child) authorization exists to prevent.

---

## 10. Wake topology (D4)

**v0:** one generic package runner per protocol. The `cds-dispatch` wake claims `protocol:cds` work and, in the specified (not-yet-shipped) target, reads `cell.class` and routes Working/Planning/Cohering internally. There is **no separate PC or CC wake provider in v0** — that would multiply providers, schedules, concurrency groups, secrets, and recovery paths before evidence requires it. **A CC runs as a claimed cell in v0**, exactly like WC/PC (`dispatch:cell + protocol:{P} + status:todo`); a **scheduled** coherence pulse (CC firing on a timer or on wave/receipt events rather than on an explicit claim) is a **future extension**, not part of the v0 architecture. This is the mode #662 itself runs in: a bootstrap Planning Cell realization on the currently shipped generic CDS/CCNF runner, **not** evidence that a `cell_class`-routing runtime already exists (§11.1's bootstrap-calibration note is the same fact stated for the record).

**State B (specified, not shipped):** one generic Cell Runner, potentially several wake-provider *profiles* over it — not because there are three kernels, but because trigger shapes differ (WC: issue/event-driven; PC: operator-authorized, infrequent; CC: scheduled | receipt-triggered | wave-triggered | operator-triggered). Whether that is realized as one manifest that reads `cell.class` and routes, or as three provider manifests sharing one runner binary, is an **open sub-decision** (§16 Q1) — not settled by this note. **Scope tag:** this manifest-count question is a **State-B (post-v0) sub-decision only; it does not reopen the v0 D4 decision.** v0 is closed — one generic provider per protocol, routing by `cell.class` **after claim**, with **no class-split provider manifest** (the paragraph above). What remains open is only the State-B manifest count, kept open exactly as the operator left it (§16 Q1).

---

## 11. FSM — State A (shipped) vs specified vs illustrative-future (D10)

Per **D10**, this section is deliberately three-way partitioned. Nothing below presents a not-yet-shipped command or mechanism as though it already runs in production.

### 11.1 State A — shipped cell FSM (`src/packages/cnos.cds/skills/cds/fsm/transitions.json`)

The actual shipped FSM is data-driven, package-owned, and evaluated by a generic engine (`src/packages/cnos.issues/commands/issues-fsm/table.go`) that never hardcodes a CDS state name. It is **not** the illustrative table `CELL-RUNTIME.md`'s embedded draft sketched — this subsection grounds the actual shipped table, and §11.5 relabels that draft's simplified FSM-events sketch as illustrative-future.

**Shipped declared states:** `ready, todo, in-progress, review, changes`. **`blocked` is reachable but not a declared top-level state** — it is a `target_state` from `in-progress` when `block_request_present` is true (a nuance load-bearing enough to name: the states array is `["ready","todo","in-progress","review","changes"]`; `blocked` is a transition target, not an enum member, in the shipped table today).

**Shipped guard vocabulary:** `run_active`, `branch_exists`, `branch_has_commits`, `pr_exists`, `pr_has_commits`, `review_request_present`, `repair_contract_present`, `cdd_artifacts_present`, `checks_passing`, `claim_request_present`, `block_request_present`, `release_request_present`.

**The shipped guard mechanism is a request-marker-file pattern**, not a bare state check: a cell (or the dispatch wake acting on its behalf) writes a typed marker file under `.cdd/unreleased/{issue}/` *before* calling `cn issues fsm evaluate --issue {N} --apply`, and the evaluator only proposes the corresponding transition when that marker is present **and** the evidentiary guards it requires (PR exists, commits beyond base, etc.) also hold:

| Marker file | Requests | Guard(s) also required |
|---|---|---|
| `CLAIM-REQUEST.yml` | `todo → in-progress` | no competing active run on the would-be `cycle/{issue}` branch (`run_active` false) |
| `REVIEW-REQUEST.yml` | `in-progress → review` | `pr_exists` and `pr_has_commits` (branch-commits-only no longer qualifies, #574) |
| `BLOCK-REQUEST.yml` | `in-progress → blocked` | none beyond the marker — this is the worker's own synchronous STOP/escalation request |
| `RELEASE-REQUEST.yml` | `in-progress → todo` (pre-work release) or delta-recovery | branches on whether matter exists since claim — no matter ⇒ requeue; matter exists ⇒ delta-recovery, never blind-requeue (#368 protection) |

Also shipped: `changes → todo` repair re-entry, gated on `repair_contract_present` (prior beta/delta findings or a repair contract present, #516) — this is the FSM-level instance of CCNF's **within-scope repair-dispatch recursion mode** (`COHERENCE-CELL-NORMAL-FORM.md` §Recursion Modes): the cell stays at the same scope, a repair contract drives re-entry, and no scope-lift occurs. And an unguarded dead-run reconciliation sweep (`cn issues fsm scan --protocol P`) that requeues or proposes delta-recovery for abandoned `in-progress` cells with no live request marker — the *sibling reconciler* to `evaluate`, not a mode of it.

### 11.2 State A — shipped command surface

Confirmed against the shipped command definitions in `src/go/internal/cli` and `src/go/internal/cell` (`cmd_cell.go`, `cmd_issues_dispatch.go`, `cmd_issues_fsm*.go`, and the `cmd_help_test.go` command-registry assertions), this cycle:

```text
cn cell return    Deliver operator verdict; status:review → status:changes on iterate/reject
cn cell resume    Re-arm an existing cycle after status:changes (preserves branch and artifacts)
cn cell finalize  Mechanical checkpoint + idempotent draft-PR open/update when a cell has matter
cn issues fsm evaluate --issue N [--apply]   Evaluate/apply one issue's FSM transition
cn issues fsm scan --protocol P [--apply]    Sweep dead in-progress cells for reconciliation
cn issues dispatch                            Authorize one design-first issue: status:ready → status:todo
```

`cn issues dispatch` is the shipped mechanism behind the "authorize dispatch: `ready → todo`" event the note's illustrative FSM-events sketch names (§11.5) — it is not illustrative, it ships today. `cn cell return/resume/finalize` are shipped (#500/#593); **`run`, `pulse`, `measure`, `bundle`, `act` are not shipped commands** (§11.5).

**Review-return (#500, closed/shipped):** `status:review → status:changes → status:in-progress → status:review` — the operator-iterate path is live.
**Stale-claim recovery (#504, open):** genuinely open, tracked as "Sub C of #583." The dead-run reconciliation rules in §11.1 (`cn issues fsm scan`) exist and cover the mechanical sweep; #504's fuller resume-or-escalate design remains unshipped.

**Protocol Package state truth (F4).** The draft's §4.4 "Protocol Package" term — the package that owns concrete protocol semantics for a matter domain — has an uneven shipped status that this note's own State-A grounding must carry, not just the `cn cell`/#500/#504 material above: **`cnos.cds`** (software) is shipped (#403); **`cnos.cdr`** (research) is shipped at v0.1 (#376); **`cnos.cdw`** (writing) is **illustrative only** — `CDD.md` v4.0.0 names `cdo`/`cdh` as *future* domain bindings, and `cdw` is not a shipped package today. Every worked example in this note that names a protocol package (§2's envelope, §3's class contracts) uses `protocol: cds` — the one shipped protocol this note's own bootstrap instance (#662) runs under; `cdr` and `cdw` are named here only for State-A completeness, not as protocols this note's contracts were validated against.

### 11.3 Specified (this note) — a distinct `CellClass` field alongside `CellKind`

§6 above is this note's own specified extension: `cell_class` becomes an *evaluated* FSM dimension via a **new, distinct `CellClass` field**, sourced from `cell.class` in the typed contract — **not** a repurposing of the shipped `CellKind` (matter-domain) seam, which is retained unchanged as the `MatterDomain` compatibility adapter. Class guards evaluate `CellClass`; domain refinements evaluate `MatterDomain`/legacy `CellKind`; the two axes stay orthogonal (a dual-field migration). None of this is shipped: `TestSeam_CellKindNotEnforced` currently locks even the existing `CellKind` seam as observation-only, and no `CellClass` field exists yet.

### 11.4 Specified (this note) — the wave FSM

The cell FSM (§11.1) governs one cell's lifecycle. It has **no wave-level counterpart today** — this is the largest gap between State A and the destination architecture, and this note specifies it without implementing it:

```text
intent → planning-required → d0-review → wave-planning → wave-review
       → dispatchable → executing → { holding | replanning | completing }
       → complete | complete-with-residuals
```

**D7 — CC judgment vs FSM transition.** CC owns the wave-state *judgment* (as typed matter, §3.3); the mechanical FSM owns the wave-state *transition*. This is the same CM/V/δ separation `CELL-RUNTIME.md` pins within one cell, lifted one scope: observations and receipts feed CC; CC emits a typed `cn.next-mca`; `V` checks the judgment's structural shape (§5); the wave FSM applies the authorized transition. CC never mutates a wave-state label itself. A wave-transition request is typed and guarded, mirroring the cell FSM's request-marker pattern (§11.1):

```yaml
wave_transition_request:
  wave: <id>
  from: wave-planning
  to:   wave-review
  requested_by_receipt: .cdd/unreleased/<issue>/
  guard: { d0_validated: true, graph_acyclic: true, child_contracts_complete: true }
```

No Go code, no schema, and no wake trigger implements this today. It is named so that #644 (PC-Wave mechanization) and #654 (PC D0-contract hardening) have a citable target rather than each re-deriving wave-level state independently.

### 11.5 Illustrative-future — command shape

The following are **not shipped** and must not be read as shipped. They illustrate the target surface the specified mechanics above (§11.3, §11.4) would eventually need a command wrapper for:

```text
cn cell pulse --class cohering --wave <id>    → { next_mca: { class, reason, requested_output } }
cn cell run   --class planning  --issue <id>  → note + wave graph + subissues + ACs + gates
cn cell run   --class working   --issue <id>  → implementation PR + receipt + status:review
```

Exact naming, argv shape, and sequencing (§16 Q3 — should a pulse command land before or after #504 stale-claim recovery) are explicitly unresolved and not decided by this note.

### 11.6 Review evidence binds to an immutable matter SHA (CCNF evidence-binding)

A β review is evidence about *a specific revision of the matter*, so it must bind to the exact matter commit SHA it reviewed. Per the CCNF evidence-binding rule: **a β review references the immutable commit SHA of the matter it reviewed; γ binds that same SHA together with the review artifact; and any matter edit after the reviewed SHA invalidates that β and requires a fresh review against the new SHA.** Mechanically this means the review-return / closeout chain must show, on hashes not prose, `matter commit (SHA_M) → β review naming SHA_M → γ binding SHA_M + review` — a review record that names no reviewed SHA/tree, or a γ receipt whose `head_sha` is a pointer phrase ("branch HEAD at push") rather than an actual SHA, does not satisfy it. This is why a matter revision is committed **on its own** (spec-only), reviewed at that exact SHA, and only then bound: co-committing the matter, the review, and the closeout in one commit destroys the ability to prove the review saw the ratified bytes and that nothing changed afterward. (State-A note: the shipped FSM does not yet enforce this binding mechanically; it is honored by commit discipline and receipt fields until a schema/FSM check pins it — a specified, not shipped, guarantee.)

---

## 12. Mechanical guards against role/class collapse

The FSM, once it evaluates `cell_class` (§6, §11.3), must prevent — mechanically where possible (label state, branch state, artifact-path ownership, role-owned output files, operator-review schema, receipt validation, `V` verdict, δ decision), not merely narratively:

- PC dispatching its own child issues without operator authorization (§3.2, §9) — enforced by the PC `V` predicate "no child auto-dispatched" (§5).
- CC implementing the work it recommends (§3.3) — enforced by the CC `V` predicate "no implementation surface modified" (§5).
- WC changing its own contract mid-cycle.
- **κ editing cell-owned artifacts** (§8) — forbidden precisely because κ ≠ α/β/γ/δ.
- α reviewing its own matter, β modifying matter, γ inventing missing evidence, δ implementing matter — all inherited unchanged from CCNF's firebreak (`COHERENCE-CELL.md` §"β Independence"; `CELL-RUNTIME.md` §"WC/PC/CC are deployment shapes, not role-relabelings").

This is the FSM-level answer to `CELL-RUNTIME.md`'s note-level AC7 (mechanical guards prevent role/class collapse).

---

## 13. Schema-first destination (D9)

This note **names, without implementing**, four schemas in a stated dependency order:

```text
cn.intent.v1  →  cn.cell.contract.v1  →  cn.next-mca.v1  →  cn.wave.v1
```

- `cn.intent.v1` — the typed object κ produces from operator intent **before any cell exists** (`{ id, source: operator, captured_by: kappa, statement, scope, constraints, desired_outcome }`). It has its own identity; a GitHub issue is one *carrier* of it, not the intent itself. The dependency order above (`cn.intent.v1 → cn.cell.contract.v1`) is why §2's envelope carries `intent_ref` rather than an inline `intent` block. **Named-not-implemented (D9):** no durable `cn.intent.v1` object exists on a pinned surface today, so §2's `intent_ref.id` is illustrative, not a resolvable reference. **Materializing a durable `cn.intent.v1` object (on a pinned, immutable surface) so `intent_ref.id` resolves is a schema-implementation deliverable** owned by the schema/PC-Wave work (S2 of #627 / #644) — it must land *before* `cn.cell.contract.v1` per the dependency order, and this note does not build it.
- `cn.cell.contract.v1` — §2 above. It **references** `cn.intent.v1` via `intent_ref: { schema, id, carrier }`; it does not restate or own intent identity (the §2 ↔ §13 reconciliation). It carries the doctrine-gate discriminator (`matter_domain` + `doctrine_affecting`) so the gate is determinable from the contract alone (§9).
- `cn.next-mca.v1` — the typed handoff CC emits (`{ class, mode, reason_code, requested_output, input_refs, requires_operator_authorization }`), consumed by whichever class the disposition names.
- `cn.wave.v1` — a typed, mechanically inspectable wave graph in which **every node carries a complete child contract**: `{ wave, goal, nodes: [ { id: <stable cell id>, contract_ref: <immutable ref to a complete cn.cell.contract.v1> } ], edges: [ { from: <cell id>, to: <cell id>, kind: depends_on } ], gates, completion }`. Each node's `contract_ref` resolves the full §2 envelope (class, protocol, matter_domain, requested_output, acceptance, dependencies, gates, doctrine_affecting, stop_conditions) — the node does **not** duplicate a partial subset of those fields. GitHub sub-issues *mirror* this graph, they do not define it.

  **Single dependency authority + projection rule (normative).** The child contract's immutable **input/output refs are the one canonical dependency authority**. A `depends_on` edge `A→B` means B's `inputs.required` references an output that A's `requested_output` produces. `cn.wave.v1` edges are therefore **mechanically derived** from producer-output→consumer-input relations across the child contracts; if edges are authored rather than derived, a validator must prove **exact parity** between the edge set and the producer/consumer relations implied by the child input/output refs. The graph never carries a dependency the child contracts do not, and never omits one they do. (`depends_on` is thus an edge, not a per-node partial field — the child's own `cn.cell.contract.v1` remains the single source of truth for each child.)

  **Completion predicates (normative).** *Child completion* is the predicate: the child's **requested output is produced** · its **acceptance predicates all pass** · its **`V` verdict is PASS** · a **bound receipt** exists for it. *Whole-wave completion* is the predicate: **every child-completion predicate holds** · **and the wave-level completion predicate holds** (the wave goal's own closure condition, recorded in `cn.wave.v1.completion`). Neither is a new §2 field — child completion is a predicate **derived from** the child contract's existing `acceptance` / `V` / receipt surfaces, not a `completion` field on the contract.

  **Verification sketch (fixtures named, not implemented).** Positive: a two-node **A→B** fixture where A produces B's required input **closes only after both bound child-completion predicates pass** and the wave-level predicate holds. Negative (each rejected): a **missing edge** (B consumes A's output but no `A→B` edge); a **spurious edge** (`A→B` with no producer/consumer relation between their refs); an **unresolved producer-input ref** (a node whose `inputs.required` names an output no sibling produces); and an **issue-closed node lacking `V`+receipt** (a node marked done whose child-completion predicate does not hold). A schema fixture additionally rejects a node with an unresolved or incomplete `contract_ref`, and a graph resolver proves every node has a complete executable contract before wave authorization.

**CC class-result shape (a refinement of the four-schema surface, not a fifth schema).** The tagged CC-disposition union (§3.3) is **not** a separate canonical schema; it is the **shape/refinement of how a CC result embeds `cn.next-mca.v1`** — a common head `{ class: cohering, judgment, disposition, requires_operator_gate }` plus a per-`disposition` payload with required/forbidden fields (`complete` forbids `next_mca`; `request_human` requires a typed gate; `hold`/`block` require a reason; `complete_with_residuals` requires residuals; the `next_mca` payload, when present, is a `cn.next-mca.v1`). This CC class-result shape is a **constraint pinned during implementation of the authorized four-schema surface above**, not a fifth canonical schema — the canonical set and dependency chain stay exactly four (`cn.intent.v1 → cn.cell.contract.v1 → cn.next-mca.v1 → cn.wave.v1`).

A PC-Wave cycle (#644, once dispatched) derives the exact implementation dependency order for these schemas from the finalized version of this note; **this note does not file, order, or dispatch that implementation work.**

---

## 14. Reconciliation map (what this note is coherent with, and how)

| Surface | State | This note's relation |
|---|---|---|
| `COHERENCE-CELL-NORMAL-FORM.md` (#370) | Landed kernel doctrine | Cited, not restated. §4 is a direct instance of Scope-Lift Projection 3; §11.1's repair re-entry is a direct instance of the within-scope repair-dispatch recursion mode. |
| `CELL-RUNTIME.md` (#628) | Landed as a proposed, not-yet-ratified architecture note (merged #629) | **Parent for this note's classes.** This note operationalizes it into contracts/FSM; it does not restate or contradict the four-orthogonal-axes frame, the WC/PC/CC deployment-shape framing, or the CM/V/δ separation. |
| `#627` (parent wave, open) | Open | Cited as parent. Its AC-INT (one generic runner routes WC/PC/CC by `cell_class`, no separate PC/CC provider) is the source of D4/§10 and resolves `CELL-RUNTIME.md`'s Open-Q3. |
| `#530` (scope-lift/ε projection artifact, open) | Open | §4 depends on it; does not redefine it. |
| `#500` (review-return) | Closed/shipped | Cited as shipped in §11.2; `review → changes → in-progress → review` is live. |
| `#504` (stale-claim recovery) | Open | Cited as genuinely open in §11.2, "Sub C of #583." Not presented as shipped. |
| `#644` (PC mechanization, S6 of #627) | Open | §3.2/§11.4 name this as PC-Wave's and the wave FSM's implementation target; this note specifies, #644 would build. |
| `#654` (PC D0-contract hardening) | Open | §3.2's D0 admissibility criteria are the D0 contract #654 hardens; §3.3's anti-pattern list names exactly the failure #654 exists to prevent. |
| `#583` (fully-mechanical dispatch) | Closed/landed | Cited in §8/§10 as the doctrine establishing the mechanism/cognition boundary this note's κ-outside-the-cell rule depends on. |
| `#584` (mechanism/cognition boundary codified) | Closed/landed | Same as #583; cited alongside it in §8. |
| κ boundary (κ ≠ α, κ outside the cell) | Invariant (D2/F3) — unconditional target (State B) | §8 states it without qualification. During bootstrap the operator **authorized κ (Sigma) to execute the α role** (State A, no mechanical Planning-Cell runtime yet) — an **operator-authorized, explicitly transitional posture**, not a nonconformance and not a licensed permanent topology; it is retired at State B and does not soften the invariant. Authorization ≠ self-warrant: the reduced *independent* warrant is discharged by the **external β bound to the matter-only R3 SHA (§11.6)** + the later γ receipt, not by the internal passes. Actor collapse stays forbidden; the bootstrap carries only the weaker hosting-identity collapse across separate Sigma activations, tracked by #664. |
| `CELL-KINDS.md:137` ("κ commonly acts as α of a planning cell") | **State-A-accurate, scoped, superseded at State B** | Correctly describes the authorized State-A bootstrap posture (κ executes α during bootstrap); accurate for State A and superseded at State B, where κ≠α once the PC runs mechanically. The **bounded downstream doctrine-migration issue** (PC-Wave / later doctrine wave) **repoints and scopes** line 137 to State-A-bootstrap wording pointing at the mechanically-separated α role — it does not correct an error. **This PC-D0 does not edit `CELL-KINDS.md`** (§8, §15). |
| Go mechanical-orchestration target | Partially shipped, partially specified | §11.1–§11.2 name exactly what's shipped (`transitions.json`, the evaluator, `cn cell`/`cn issues` commands); §11.3–§11.5 name what's specified or illustrative; none of it is implemented by this cycle. |

---

## 15. Non-goals

Redefining CCNF or its role semantics; adding a new role letter; replacing CDD or CDS; implementing the WC/PC/CC runtime, the wave FSM, or any of the named schemas in Go, CUE, or any other code; changing α/β/γ/δ semantics; filing, labeling, or dispatching any child issue as a side effect of writing this note (§3.2's own guard applies reflexively to the cell that produced this note); a runner-per-class or a wake-per-class in v0 (D4); auto-merge or gate-crossing authority; a global cross-repo scheduler; rewriting `CDD.md` (Phase 7 of #366, unaffected); **editing `CELL-KINDS.md`** (its `:137` κ-as-α line is State-A-accurate, scoped, superseded-at-State-B prose whose reconciliation is a bounded downstream doctrine-migration dependency, §8/§14 — not this cell's work); new agent-personality doctrine.

---

## 16. Open questions

### Exit sequence (operator-authorized)

Ratification is **not** this cell's job. The operator-authorized ordered exit sequence for this PC-D0 matter is:

1. **External-β review** — a review entity outside the Sigma lineage reviews the immutable matter SHA (SHA_M, §11.6).
2. **γ binds SHA_M + review + evidence** — γ records the reviewed SHA, the external review artifact, and the cell's evidence, on hashes not prose.
3. **Separate Cohering Cell (CC) ratification** — a distinct CC judges coherence and emits its disposition.
4. **Operator final-read** — the operator reads and accepts (doctrine-affecting gate, §9).
5. **Merge of this D0** — the D0 matter merges.
6. **Separate PC-Wave** — a distinct Planning Cell in wave mode plans the executable wave from this validated D0. This is the **immediate post-merge successor**.
7. **CC wave review** — a Cohering Cell judges the wave's dispatchability.
8. **Operator wave authorization** — the operator authorizes the wave at its boundary (§9).

The downstream doctrine-migration of `CELL-KINDS.md:137` (§8, §14) is **not** a step in this sequence; it is a dependency **owned by the PC-Wave / later doctrine wave** (step 6 onward), not the immediate post-merge action.

Most of the embedded draft's original open-question list (its own §19) is resolved by D1–D10 and is folded into the body above rather than left dangling: `cell_class` location (§2, §6 — contract-first, per D6/F2), CC pulse triggering (§10 — claimed cell in v0, scheduled pulse a stated future extension), PC's dispatch provider (§10 — one generic runner, per D4/#627 AC-INT), the minimum D0 contract (§3.2, owned going forward by #654), and the minimal `V` predicate for PC/CC (§5, per D8). What remains genuinely open after D1–D10:

1. **Wake-provider realization (§10).** Whether v0's "one generic runner" eventually surfaces as a single manifest that reads `cell.class` and routes, or as several provider manifests (distinct triggers/cadence/permissions/concurrency) sharing one runner binary. Operator framing leans toward "several trigger shapes, one runner" but the manifest-count question is explicitly left open. **Scope:** this is a **State-B (post-v0) sub-decision only and does not reopen the v0 D4 decision** — v0 is settled (one generic provider per protocol, route by `cell.class` after claim, no class-split provider manifest, §10); only the State-B manifest count remains open, per operator framing.
2. **What receipts prove a wave graph is dispatchable.** §11.4 names the wave-transition-request guard shape (`d0_validated`, `graph_acyclic`, `child_contracts_complete`); the exact evidence each guard dereferences is unspecified.
3. **Sequencing of the illustrative command surface against #504.** Whether a `cn cell pulse` (or equivalent) should land before or after stale-claim recovery ships — left as a sequencing question for whoever builds §11.4/§11.5, not decided here.
4. **Exact schema field types** for `cn.intent.v1`, `cn.cell.contract.v1`, `cn.next-mca.v1`, `cn.wave.v1` (§13) — named, not typed; typing is explicitly deferred schema work.
5. **Concurrency and idempotence at wave scope** (one active CC per wave/snapshot; one PC-D0 per intent; one PC-Wave per validated D0; multiple WCs only when the graph and write-surfaces allow it) — stated as a requirement in the operator's State-B framing, not yet reduced to a mechanically checkable predicate the way §5's per-class `V` layers are.

---

*Authoring note.* Formalized under cnos#662, a Planning Cell in PC-D0 mode, dispatched by the `cds-dispatch` wake in δ wake-invoked mode (bootstrap realization on the currently shipped generic CDS/CCNF runner — §10, §11 — not evidence that `cell_class`-aware routing already ships). α reconciled the issue's embedded κ-corrected draft (§1–§20 of #662) against: the ten operator-pinned decisions (D1–D10) from the 2026-07-13T19:14:52Z authorization comment and the two preceding pinned-decisions comments (2026-07-11, 2026-07-12); `CELL-RUNTIME.md` (#628) for structure and framing; `COHERENCE-CELL-NORMAL-FORM.md` for the kernel and scope-lift vocabulary; and the shipped `transitions.json` / the shipped `cn` command surface (`src/go/internal/cli`, `src/go/internal/cell`) for State-A ground truth (§11), correcting the draft's illustrative FSM-events table against the actual request-marker-file mechanism rather than restating it as shipped. **R3 repaired the typed-contract and reconciliation surface against an external-β ITERATE (PR #667):** canonical envelope now parses to the worked instance's shape (§2); the worked instance is marked illustrative and a durable `cn.intent.v1` is named a schema deliverable (§2, §13); `cell_class` gets a distinct `CellClass` field, not a repurposing of the matter-domain `CellKind` seam (§6, §11.3); `cn.wave.v1` nodes carry a complete child contract by `contract_ref` (§13); review evidence binds to an immutable matter SHA (§11.6); the doctrine gate is determinable from the canonical contract alone (§2, §9); the CC result is a tagged disposition union (§3.3, §13); and §8 states κ≠α as an unconditional State-B target while framing the bootstrap κ=α as the **operator-authorized transitional posture (State A)** — a legitimate, explicitly transitional bootstrap stage the operator settled directly, not a nonconformance and not a licensed permanent topology, still carrying a reduced *independent* warrant (hosting-identity, #664) that the external β discharges, with `CELL-KINDS.md:137` marked State-A-accurate/scoped/superseded-at-State-B and its repointing bound to a downstream doctrine-migration issue. The R3 matter is committed **spec-only** at an immutable SHA for exact-SHA external-β review (§11.6). **R4 repaired the matter against an external-β ITERATE verdict (PR #667):** a single canonical wave dependency authority — the child contract input/output refs are canonical, `cn.wave.v1` edges are their mechanical projection (or parity-validated) — plus explicit child- and whole-wave-completion predicates and a positive/negative verification sketch (§3.2, §13); the canonical schema set restored to **exactly four**, the CC tagged-disposition union reframed as the **CC class-result shape/refinement** rather than a fifth `cn.cc-result.v1` schema (§3.3, §13); the §2 claim narrowed to key-path identity plus a **normative type/cardinality constraint model** with a **decidable gate-invariant set**, and the template YAML made internally consistent (`reason: null` when both gate booleans are false) (§2); an explicit **Exit sequence** block added to §16 with the header, §3.2, and this note repointed to it and the immediate post-merge successor corrected to a **separate PC-Wave**; and §10/§16 Q1 hardened so the State-B manifest-count is scoped as a post-v0 sub-decision that does **not** reopen the v0 D4 decision (the disputed finding-2 hardened, not reversed). No implementation, schema, FSM code, or child issue was produced or dispatched by this cell; `CELL-KINDS.md` was not edited, and the operator-settled §8 authorized-transitional framing was preserved unchanged. Draft, not ratified — the remaining steps are the **Exit sequence** block in §16: external-β review → γ binds SHA_M + review + evidence → a separate Cohering Cell ratification → operator-final-read → merge of this D0 → a **separate PC-Wave** (the immediate post-merge successor) → CC wave review → operator wave authorization.
