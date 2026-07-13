# Cell Runtime Classes — Contracts, FSM, and Routing

**Status:** Proposed architecture note (realization layer). Not ratified. Formalized by a Planning Cell in PC-D0 mode (cnos#662), under the operator-authorization comment of 2026-07-13T19:14:52Z. Ratification is **not** this cell's job: the exit sequence is a separate Cohering Cell (CC) review, then operator-final-read, then merge, then a separate CDD `doctrine` cell — see §16.
**Owns:** How the WC/PC/CC output-telos classes named by `docs/architecture/CELL-RUNTIME.md` (#628) route through a common cell contract envelope, class-specific `V` predicates, the CC↔ε lineage, the cell FSM and the (not-yet-shipped) wave FSM, the command surface, and the human-gate/wake-topology policy that makes class-aware dispatch mechanical.
**Does-not-own:** The kernel algorithm (`COHERENCE-CELL-NORMAL-FORM.md`, #370), the WC/PC/CC class definitions themselves (`CELL-RUNTIME.md`, #628 — cited throughout, never restated), the receipt rule and four surfaces (`COHERENCE-CELL.md`), the CUE schema implementation (`schemas/cdd/`), or any Go runtime/FSM code. This note **names** schemas and FSM extensions (D9); it does not implement them.

---

## Thesis

`CELL-RUNTIME.md` established that CNOS has one CCNF kernel, deployed in three output-telos shapes — WC (artifact), PC (relation graph), CC (process judgment) — executed by one generic runner. That note stops at the class definitions. It does not say what a cell contract looks like as a typed object, what a Cohering Cell consumes and emits, how `cell_class` interacts with the shipped label-driven FSM, or which parts of the "mechanical FSM" sketch are running in production today versus still design. This note is the **operationalization layer**: it takes the classes as given and specifies the routing, contracts, and FSM extension that let a generic runner dispatch them without inventing policy at runtime. It carries ten operator-pinned decisions (D1–D10, §16) as settled input, not as open questions to re-derive.

---

## 1. Relationship to CELL-RUNTIME.md and the kernel

This note is a **sibling** of `docs/architecture/CELL-RUNTIME.md`, not a replacement or a restatement. The relationship is strict:

| Surface | Owns | This note |
|---|---|---|
| `COHERENCE-CELL-NORMAL-FORM.md` (#370) | The five-step kernel `α→β→γ→V→δ`, the four cell outcomes, the two recursion modes, the three scope-lift projections | Cited; §4 below is a direct instance of Projection 3 (ε) |
| `CELL-RUNTIME.md` (#628) | WC/PC/CC as output-telos classes of one kernel; the four orthogonal axes (kernel cell / roles / output-telos / matter domain); CM measures, V gates, δ effects | Consumed as fixed input; this note does not re-derive the classes |
| `CELL-KINDS.md` (#570) | The matter/contract domain vocabulary (axis 4) | Unaffected; `cell_class` (this note) and `matter_domain` (`CELL-KINDS.md`) are orthogonal fields on the same contract |
| **This note** (#662) | The cell contract envelope; class-specific `V`; CC↔ε; the FSM extension point for `cell_class`; State A/State B honesty; the wave FSM sketch | — |

**D1 (one CCNF kernel).** WC, PC, and CC are output-telos classes of the *same* kernel — already `CELL-RUNTIME.md`'s own thesis. This note confirms it and does not re-derive it: every class below still runs its own internal, independent `α→β→γ→V→δ`; nothing here forks the kernel, adds a role letter, or changes α/β/γ/δ semantics.

---

## 2. The cell contract envelope (`cn.cell.contract.v1`)

All three classes consume one typed envelope. Only *input admissibility, matter type, class-specific acceptance, and post-closure routing* differ by class (§3, §5). Per **D9 (schema-first destination)**, this note **names** the schema; it does not implement it — no CUE, no Go struct, no validator.

```yaml
schema: cn.cell.contract.v1
cell:   { id: <issue>, class: working | planning | cohering, mode: <class-specific>, protocol: cds, matter_domain: <CELL-KINDS.md vocabulary> }
scope:  { repo: <owner/name>, wave: <parent wave issue, if any>, parent_cell: <parent cell issue, if any> }
intent: { source: issue, source_ref: <issue>, operator_directive_ref: <issue comment anchor>, summary: <one line> }
inputs:
  required: [ <artifacts/issues this cell depends on> ]
  optional: [ prior_receipts ]
requested_output: { kind: <artifact|relation_graph|judgment>, path: <where the matter lands> }
acceptance:
  predicates: [ <class-specific acceptance predicates, §5> ]
constraints:
  allowed_paths:   [ <surfaces this cell's matter may touch> ]
  forbidden_paths: [ <surfaces this cell's matter must not touch> ]
  non_goals: [ <explicit non-goals> ]
gates: { operator_authorization_required: <bool>, operator_acceptance_required: <bool> }
stop_conditions: [ <typed STOP triggers, e.g. source_doctrine_conflicts, required_decision_unresolved> ]
```

This is a worked, not hypothetical, template: #662's own γ-scaffold (`.cdd/unreleased/662/gamma-scaffold.md §1`) instantiates it —

```yaml
cell:
  class: planning
  protocol: cds
  matter_domain: doctrine
  issue: 662
  requested_output: [ "docs/architecture/CELL-RUNTIME-CLASSES.md" ]
  non_goals: [ implement_cell_runner, alter_ccnf, add_new_role_letters, replace_cdd_or_cds, file_or_dispatch_child_issues ]
  gates: { operator_authorization_required: true, operator_acceptance_required: true }
  stop_conditions: [ source_doctrine_conflicts, required_decision_unresolved ]
```

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
- **PC-Wave** — `validated D0 → executable wave`. Admissible only once the D0 is validated, the target is sufficiently defined, scope/non-goals are stable, major decisions have owners, and surfaces are known. Output passes only if: the graph is acyclic; every node carries `class` + `protocol` + `matter_domain` + `requested_output` + acceptance criteria; dependencies are explicit; the critical path is explicit; parallel cells don't share write surfaces; human gates are located; **no child is auto-dispatched**; STOP conditions exist at cell *and* wave level; a wave-completion predicate is defined. #644 is the open implementation target for PC-Wave mechanization; **this note specifies its contract, it does not build it.**

**Guard, both modes:** a Planning Cell creates and refines its own children; it does not apply `status:todo` to them. Dispatching a PC's own children requires operator authorization (or a typed policy gate at the wave-authorization boundary, §9) — never an automatic action taken by the PC itself. This is the same guard #662 itself operates under: it produces exactly one artifact and files or dispatches nothing.

**Result shape:** `{ class: planning, mode: d0|wave, wave_ref, graph: { nodes, edges }, dispatchable: false, requires_operator_gate: true }`.

### 3.3 Cohering Cell (CC)

*Telos:* increase justified confidence about state and select the next coherence-preserving action — matter targets process (TSC-γ), output is a coherence judgment plus exactly one next disposition. CC does **not** write implementation or planning matter; it **judges** state and, where more work is needed, names the class of cell that should do it (§4 sharpens the ε relationship this telos depends on).

**Input:** operator intent, scope/issue/PR/CI state, wave graph, receipts, decisions, open risks, **ε cross-cell observations** (§4), prior CC receipt. Admissible if: scope is named, state is mechanically snapshottable, relevant receipts are available, an intent/pulse reason is present, the observation window is explicit, no other CC is active for the same scope/version, and the disposition vocabulary is known.

**Matter:** coherence assessment, gap classification, bottleneck axis, contradiction map, wave-progress judgment, residual risk, next-MCA recommendation (MCA = Minimum Coherence Action).

**Output passes only if it:** names the evaluated scope; names the state snapshot; cites evidence per load-bearing claim; separates observation from inference; identifies the dominant bottleneck; emits **exactly one** primary disposition; emits a typed next-MCA when more work is needed; performs **none** of the recommended work itself; distinguishes `done` from `done-with-residuals`; requests human action only when the mechanics genuinely cannot decide.

**Dispositions:** `request_planning | request_working | hold | request_human | continue_wave | complete | complete_with_residuals | block`.

**Result shape:** `{ class: cohering, judgment, disposition, next_mca, requires_operator_gate }`. CC recommends; `V` validates the judgment's structural shape; the FSM (not CC) applies any resulting transition (**D7**, §7).

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
| **CC `V`** adds | Judgment is evidence-backed · exactly one disposition emitted · a typed next-MCA is present **or** a typed human-gate/hold is named · **no implementation surface modified** (mechanically checkable as an allowed-surface check on the CC's own matter: matter paths ⊆ judgment artifacts, no code/product diff) |

These predicates are the seed for the future per-class `cell.cue` validators (`CELL-RUNTIME.md`'s Reconciliation map, `#369`/#627 S2–S3). This note names them as the destination `V` extends toward; it implements none of them.

---

## 6. `cell_class` as an FSM dimension (D6, F2)

The runtime eventually evaluates the triple **`(status, protocol, cell_class)`**, not just `(status, protocol)`:

- **`status`** (label-based) answers *where is this cell in its lifecycle?*
- **`protocol`** (label-based) answers *which package's semantics apply?* (`cds`, and eventually `cdr`/others)
- **`cell_class`** (contract-based) answers *what kind of cell is it — what authority and guards apply?* Can it create child issues? Can it dispatch? Can it modify implementation surfaces? Can it emit a next-MCA?

**Dispatch readiness stays solely label-gated:** `dispatch:cell + protocol:{P} + status:todo` — exactly the shipped #643 rule. `cell_class` is **post-claim** routing/authority metadata sourced from the typed contract, not a second readiness gate. If a claimed cell's contract is missing `cell.class`, `requested_output`, or acceptance criteria, the runtime routes it to `status:blocked` with `degraded_reason: cell_contract_incomplete` — a post-claim validation failure, not a pre-claim admission failure. This answers `CELL-RUNTIME.md`'s Open-Q1/Open-Q2 without contradicting #643 or #640.

**Grounding in the shipped seam.** The shipped FSM already carries an observation-only seam for exactly this purpose: `FactSnapshot.CellKind{Observed, Source, DefaultedTo}`, locked as *observation, not enforcement* by `TestSeam_CellKindNotEnforced` (see `src/packages/cnos.cds/skills/cds/fsm/transitions.json` and its adjacent Go evaluator). This note's `cell_class` dimension is the **specified promotion of that seam from observation to evaluation** — the deferred "FSM Phase 2" `CELL-KINDS.md` names — keyed on `cell_class` for behavior/authority, with `matter_domain` remaining the separate domain refinement (`CELL-KINDS.md`, unaffected). The exact source mechanics (contract field vs. observation seam vs. eventual label) are for schema work to pin (§13); the decision that `cell_class` is an evaluated FSM dimension is settled here.

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

## 8. κ remains outside the cell (D2, F3)

κ is a **control-plane** slot, not a cell role — it does not appear beside α/β/γ/δ/ε in CCNF and this note adds no new role letter. **κ supplies operator/control-plane input; the cell's α owns and writes the specification or implementation matter.** κ must not author the specification as α. If a single model body performs both κ and α in a bootstrap case (as happened structurally across #662's own drafting and dispatch history), that is **actor collapse** and must be **declared explicitly** as a configuration-floor consequence — never silently assumed as equivalence. κ speaks through issues and issue comments; cells work through typed contracts and receipts. This aligns with #583 and #584 (both landed doctrine: the mechanism/cognition boundary — mechanical dispatch is fully separated from the cognitive skills that run inside a claimed cell).

**Guard consequence (§12):** "κ editing cell-owned artifacts" is a forbidden action precisely because κ ≠ α/β/γ/δ. κ carries intent into control-plane artifacts (issues, labels, PR-review surfaces); it does not produce, review, or close a cell's matter.

---

## 9. Human gates (D5)

Human authority sits at **irreversible or scope-expanding boundaries**, not between internal phases:

**Default gates (always human):** intent acceptance · wave authorization (first dispatch of a wave, not every child in it) · production/release boundary · final acceptance of doctrine-affecting or otherwise system-doctrine-changing matter · explicit hold/block/escalation resolution. The doctrine-affecting condition fires off a **typed contract flag** (`matter_domain: doctrine` / `doctrine_affecting: true`, as #662's own contract carries — §2) — mechanical systems consume declarations, they never classify prose to decide whether a gate applies.

**No default gate between:** `CC → requested PC/WC` · PC's internal α/β/γ · WC's internal α/β/γ · safe mechanical recovery (stale-claim resume, repair-dispatch re-entry). A wave, once operator-authorized at its boundary, may pre-authorize its own dependency-respecting internal nodes — `operator → authorize wave → wave executable → FSM routes children by dependencies without further operator intervention`, re-entering only at a named boundary (human gate, blocked, review, release). **The operator must not become the scheduler for every child** — that failure mode is exactly what wave-boundary (not per-child) authorization exists to prevent.

---

## 10. Wake topology (D4)

**v0:** one generic package runner per protocol. The `cds-dispatch` wake claims `protocol:cds` work and, in the specified (not-yet-shipped) target, reads `cell.class` and routes Working/Planning/Cohering internally. There is **no separate PC or CC wake provider in v0** — that would multiply providers, schedules, concurrency groups, secrets, and recovery paths before evidence requires it. **A CC runs as a claimed cell in v0**, exactly like WC/PC (`dispatch:cell + protocol:{P} + status:todo`); a **scheduled** coherence pulse (CC firing on a timer or on wave/receipt events rather than on an explicit claim) is a **future extension**, not part of the v0 architecture. This is the mode #662 itself runs in: a bootstrap Planning Cell realization on the currently shipped generic CDS/CCNF runner, **not** evidence that a `cell_class`-routing runtime already exists (§11.1's bootstrap-calibration note is the same fact stated for the record).

**State B (specified, not shipped):** one generic Cell Runner, potentially several wake-provider *profiles* over it — not because there are three kernels, but because trigger shapes differ (WC: issue/event-driven; PC: operator-authorized, infrequent; CC: scheduled | receipt-triggered | wave-triggered | operator-triggered). Whether that is realized as one manifest that reads `cell.class` and routes, or as three provider manifests sharing one runner binary, is an **open sub-decision** (§17) — not settled by this note.

---

## 11. FSM — State A (shipped) vs specified vs illustrative-future (D10)

Per **D10**, this section is deliberately three-way partitioned. Nothing below presents a not-yet-shipped command or mechanism as though it already runs in production.

### 11.1 State A — shipped cell FSM (`src/packages/cnos.cds/skills/cds/fsm/transitions.json`)

The actual shipped FSM is data-driven, package-owned, and evaluated by a generic engine (`src/packages/cnos.issues/commands/issues-fsm/table.go`) that never hardcodes a CDS state name. It is **not** the illustrative table `CELL-RUNTIME.md`'s embedded draft sketched (§11.6 corrects that table explicitly).

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

Confirmed directly against the built binary (`cn --help`, this cycle):

```text
cn cell return    Deliver operator verdict; status:review → status:changes on iterate/reject
cn cell resume    Re-arm an existing cycle after status:changes (preserves branch and artifacts)
cn cell finalize  Mechanical checkpoint + idempotent draft-PR open/update when a cell has matter
cn issues fsm evaluate --issue N [--apply]   Evaluate/apply one issue's FSM transition
cn issues fsm scan --protocol P [--apply]    Sweep dead in-progress cells for reconciliation
cn issues dispatch                            Authorize one design-first issue: status:ready → status:todo
```

`cn issues dispatch` is the shipped mechanism behind the "authorize dispatch: `ready → todo`" event the note's illustrative FSM-events sketch names (§11.6) — it is not illustrative, it ships today. `cn cell return/resume/finalize` are shipped (#500/#593); **`run`, `pulse`, `measure`, `bundle`, `act` are not shipped commands** (§11.6).

**Review-return (#500, closed/shipped):** `status:review → status:changes → status:in-progress → status:review` — the operator-iterate path is live.
**Stale-claim recovery (#504, open):** genuinely open, tracked as "Sub C of #583." The dead-run reconciliation rules in §11.1 (`cn issues fsm scan`) exist and cover the mechanical sweep; #504's fuller resume-or-escalate design remains unshipped.

**Protocol Package state truth (F4).** The draft's §4.4 "Protocol Package" term — the package that owns concrete protocol semantics for a matter domain — has an uneven shipped status that this note's own State-A grounding must carry, not just the `cn cell`/#500/#504 material above: **`cnos.cds`** (software) is shipped (#403); **`cnos.cdr`** (research) is shipped at v0.1 (#376); **`cnos.cdw`** (writing) is **illustrative only** — `CDD.md` v4.0.0 names `cdo`/`cdh` as *future* domain bindings, and `cdw` is not a shipped package today. Every worked example in this note that names a protocol package (§2's envelope, §3's class contracts) uses `protocol: cds` — the one shipped protocol this note's own bootstrap instance (#662) runs under; `cdr` and `cdw` are named here only for State-A completeness, not as protocols this note's contracts were validated against.

### 11.3 Specified (this note) — promoting `CellKind` from observation to evaluation

§6 above is this note's own specified extension: `cell_class` becomes an *evaluated* FSM dimension, not merely an observed one. This is not shipped; it is what `TestSeam_CellKindNotEnforced` currently locks as future work.

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

Exact naming, argv shape, and sequencing (§17 Q6 — should a pulse command land before or after #504 stale-claim recovery) are explicitly unresolved and not decided by this note.

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

- `cn.intent.v1` — the typed object κ produces from operator intent before any cell exists (`{ id, source: operator, captured_by: kappa, statement, scope, constraints, desired_outcome }`).
- `cn.cell.contract.v1` — §2 above.
- `cn.next-mca.v1` — the typed handoff CC emits (`{ class, mode, reason_code, requested_output, input_refs, requires_operator_authorization }`), consumed by whichever class the disposition names.
- `cn.wave.v1` — a typed, mechanically inspectable wave graph (`{ wave, goal, nodes: [{issue, class, protocol, depends_on}], edges, gates, completion }`) that GitHub sub-issues *mirror*, not define.

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
| κ boundary | Settled by D2/F3 | §8 carries it forward without re-litigation; #662's own dispatch history (κ authored the issue and its comments; α, in a separate cell instance, authors this file) is the worked example of declared, non-collapsed actor separation. |
| Go mechanical-orchestration target | Partially shipped, partially specified | §11.1–§11.2 name exactly what's shipped (`transitions.json`, the evaluator, `cn cell`/`cn issues` commands); §11.3–§11.5 name what's specified or illustrative; none of it is implemented by this cycle. |

---

## 15. Non-goals

Redefining CCNF or its role semantics; adding a new role letter; replacing CDD or CDS; implementing the WC/PC/CC runtime, the wave FSM, or any of the named schemas in Go, CUE, or any other code; changing α/β/γ/δ semantics; filing, labeling, or dispatching any child issue as a side effect of writing this note (§3.2's own guard applies reflexively to the cell that produced this note); a runner-per-class or a wake-per-class in v0 (D4); auto-merge or gate-crossing authority; a global cross-repo scheduler; rewriting `CDD.md` (Phase 7 of #366, unaffected); new agent-personality doctrine.

---

## 16. Open questions

Most of the embedded draft's original open-question list (its own §19) is resolved by D1–D10 and is folded into the body above rather than left dangling: `cell_class` location (§2, §6 — contract-first, per D6/F2), CC pulse triggering (§10 — claimed cell in v0, scheduled pulse a stated future extension), PC's dispatch provider (§10 — one generic runner, per D4/#627 AC-INT), the minimum D0 contract (§3.2, owned going forward by #654), and the minimal `V` predicate for PC/CC (§5, per D8). What remains genuinely open after D1–D10:

1. **Wake-provider realization (§10).** Whether v0's "one generic runner" eventually surfaces as a single manifest that reads `cell.class` and routes, or as several provider manifests (distinct triggers/cadence/permissions/concurrency) sharing one runner binary. Operator framing leans toward "several trigger shapes, one runner" but the manifest-count question is explicitly left open.
2. **What receipts prove a wave graph is dispatchable.** §11.4 names the wave-transition-request guard shape (`d0_validated`, `graph_acyclic`, `child_contracts_complete`); the exact evidence each guard dereferences is unspecified.
3. **Sequencing of the illustrative command surface against #504.** Whether a `cn cell pulse` (or equivalent) should land before or after stale-claim recovery ships — left as a sequencing question for whoever builds §11.4/§11.5, not decided here.
4. **Exact schema field types** for `cn.intent.v1`, `cn.cell.contract.v1`, `cn.next-mca.v1`, `cn.wave.v1` (§13) — named, not typed; typing is explicitly deferred schema work.
5. **Concurrency and idempotence at wave scope** (one active CC per wave/snapshot; one PC-D0 per intent; one PC-Wave per validated D0; multiple WCs only when the graph and write-surfaces allow it) — stated as a requirement in the operator's State-B framing, not yet reduced to a mechanically checkable predicate the way §5's per-class `V` layers are.

---

*Authoring note.* Formalized under cnos#662, a Planning Cell in PC-D0 mode, dispatched by the `cds-dispatch` wake in δ wake-invoked mode (bootstrap realization on the currently shipped generic CDS/CCNF runner — §10, §11 — not evidence that `cell_class`-aware routing already ships). α reconciled the issue's embedded κ-corrected draft (§1–§20 of #662) against: the ten operator-pinned decisions (D1–D10) from the 2026-07-13T19:14:52Z authorization comment and the two preceding pinned-decisions comments (2026-07-11, 2026-07-12) that supersede the draft's own open-question wording; `CELL-RUNTIME.md` (#628) for structure and framing; `COHERENCE-CELL-NORMAL-FORM.md` for the kernel and scope-lift vocabulary; and the shipped `transitions.json` / built `cn` binary for State-A ground truth (§11), correcting the draft's illustrative FSM-events table against the actual request-marker-file mechanism rather than restating it as shipped. No implementation, schema, FSM code, or child issue was produced or dispatched by this cell. Proposed, not ratified — a separate Cohering Cell review, operator-final-read, and a CDD `doctrine` cell are the remaining steps in the exit sequence the authorization comment names.
