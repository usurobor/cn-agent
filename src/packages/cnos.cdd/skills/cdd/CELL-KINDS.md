# CELL-KINDS — Cell-Kind Taxonomy

**Status:** Doctrine-first (Phase A of [#570](https://github.com/usurobor/cnos/issues/570)). Deferred: CUE schema, verifier enforcement, UI display, automatic inference, full wave/cleanup execution, `table.go` transition consumption.
**Placement:** `src/packages/cnos.cdd/skills/cdd/`
**Parent:** [`CDD.md`](CDD.md) — this document refines the generic kernel; it does not replace or re-derive it.

## Design rule: one kernel, multiple cell kinds

[`CDD.md`](CDD.md) names the generic recursive coherence-cell kernel:

```text
contractₙ -> matterₙ -> reviewₙ -> receiptₙ -> verdictₙ -> decisionₙ
```

Every closed cell — regardless of what it produces — runs this same five-step closure. **Cell kind does not fork the kernel.** A cell kind is a *typed refinement* of the kernel: it names what matter the kernel's `α.produce` step is expected to emit, what `β.review` means for that matter, what the `γ.close` receipt projects, and what boundary decisions `δ` may take at that cell's scope. Do not create a separate lifecycle per cell kind — bind the existing kernel to a specific matter shape instead.

```text
cell kind -> expected matter -> review surface -> closeout/projection
```

## `issue kind` != `cell kind`

These are two different, non-interchangeable classifications:

- **`issue kind`** (e.g. `kind/process`, `kind/bug`, `kind/feature` — see `docs/development/issues/TAXONOMY.md`) is a **GitHub issue taxonomy label**. It classifies *what an issue is about*.
- **`cell_kind`** (e.g. `implementation`, `issue_authoring`, `wave`) is a **lifecycle declaration**. It classifies *what the coherence cell dispatched against that issue is expected to produce as matter*.

An issue labeled `kind/feature` can be worked by an `implementation` cell (the common case) or, earlier in its life, be *itself* the matter of an `issue_authoring` cell that drafted it. The GitHub label and the `cell_kind` value travel on different axes and must not be conflated.

## The 10 cell kinds

### 1. `issue_authoring`

- **Purpose:** turn a selected gap into an executable issue contract.
- **Required input:** a selected gap, active design constraints, the `cdd/issue` skill family.
- **Matter produced:** the GitHub issue body / issue pack — problem, impact, status truth, source of truth, scope, ACs, proof plan, non-goals. **The issue itself is the matter** (see `issue/SKILL.md` §"An issue can be a cell's matter").
- **Review surface:** β reviews the issue for executability and lack of ambiguity — the `issue/SKILL.md` handoff checklist, not a code diff.
- **Closeout / projection:** γ confirms process and closes; the closed cell projects the finished issue (usually `status:ready`) as α-matter for the *next*, separate `implementation` cell — it does not itself dispatch implementation work.
- **FSM implication:** may create or update an issue. Does not imply the created issue should be dispatched immediately unless the issue contract itself says so.

### 2. `implementation`

- **Purpose:** produce software/docs/package/runtime matter for an already-approved issue.
- **Required input:** an executable issue at `status:todo` or later, active skills, branch state, CI state.
- **Matter produced:** branch, commits, PR, changed files, validation evidence, receipts.
- **Review surface:** β reviews the diff against the issue's ACs — the CDS dispatch review path.
- **Closeout / projection:** γ closes via `alpha-closeout.md` / `beta-closeout.md` / `gamma-closeout.md`; the closed cell projects as accepted matter into the release stream. Maps directly to existing CDS dispatch (`cnos.cds/skills/cds/CDS.md` §"Development lifecycle").
- **FSM implication:** may request `status:review` only with deliverable proof (PR/diff/commits). This is the FSM's current default assumption (`FactSnapshot.CellKind.DefaultedTo == "implementation"` when nothing else is observed).

### 3. `repair`

- **Purpose:** repair a rejected or changes-requested cell.
- **Required input:** repair contract, prior β/δ findings, existing branch state.
- **Matter produced:** targeted fixes addressing the prior findings, plus a fix-round appendix to `self-coherence.md`.
- **Review surface:** β re-reviews only the findings named as addressed, not the whole diff again from scratch.
- **Closeout / projection:** rejoins the same cell's closeout once β approves; does not open a new cell.
- **FSM implication:** `status:changes -> status:todo` requires repair context (a repair contract or an "iterate" verdict present) — see `RepairContractPresent` in the FSM fact model.

### 4. `recovery`
- **Purpose:** recover a stranded or partially completed cell.
- **Required input:** observed branch/run state (dead run, orphaned branch, ambiguous label state).
- **Matter produced:** a recovery PR, a STOP comment, a requeue decision, or an explicit blocked state.
- **Review surface:** δ/operator judgment over the observed facts, not a normal β content review.
- **Closeout / projection:** either requeues the original cell (`status:todo`) or hands off to δ-recovery; does not itself close as `accepted`.
- **FSM implication:** dead run with matter (commits beyond base) routes to δ-recovery and must never silently propose `status:todo` (cnos#368 regression class); dead run with no matter may requeue; ambiguous state blocks.

### 5. `wave`

- **Purpose:** coordinate several child cells under one approved contract.
- **Required input:** an approved wave contract naming the child cells and their independent shippability.
- **Matter produced:** child-cell receipts, child projections, a wave receipt, a wave handoff, follow-up issues.
- **Review surface:** parent β reviews wave coherence across children, not any single child's diff.
- **Closeout / projection:** parent γ closes the wave; child receipts *project* upward into the wave receipt — **projection, not role-renaming** (`CDD.md` §"Scope-lift": a wave's δ is not literally a child's β; the wave has its own full α/β/γ/δ/ε at its own scope).
- **FSM implication:** a wave must not be treated as one normal `implementation` PR unless its own contract says so; each child cell keeps its own `cell_kind` and its own FSM evaluation.

### 6. `cleanup`

- **Purpose:** improve repo/system coherence by measuring drift and producing cleanup work.
- **Required input:** a coherence-measurement pass (e.g. TSC) or an operator-flagged drift signal.
- **Matter produced:** a TSC report, a drift inventory, a cleanup PR, or cleanup issues.
- **Review surface:** β reviews the drift claims and any resulting cleanup diff for correctness and scope.
- **Closeout / projection:** may spawn `issue_authoring` or `implementation` child cells rather than closing with its own code diff.
- **FSM implication:** may create `issue_authoring` or `implementation` cells as children; must not bypass issue/review/receipt for the cleanup work itself.

### 7. `audit`

- **Purpose:** inspect system state and produce a report, not necessarily a code change.
- **Required input:** the surface to be audited and the invariant being checked.
- **Matter produced:** an inventory, a classification, a drift report, a risk report, a proof table.
- **Review surface:** β reviews the report's evidence and conclusions, not a diff.
- **Closeout / projection:** may close with docs/evidence only; a PR may be optional if an issue comment or artifact is the matter.
- **FSM implication:** does not require a PR to reach `status:review`; the review surface is the report artifact itself.

### 8. `release`

- **Purpose:** move accepted matter across a release/deploy boundary.
- **Required input:** an accepted, closed implementation cell (or several) ready to cross the release boundary.
- **Matter produced:** a tag, release notes, a deploy branch, a release receipt.
- **Review surface:** δ verifies release-readiness preconditions (`cnos.cds/skills/cds/CDS.md` §"Gate").
- **Closeout / projection:** closes the release cycle; projects as the shipped version other cells build on next.
- **FSM implication:** requires explicit boundary authority (δ); must not be inferred from an `implementation` cell's own transitions.

### 9. `doctrine`

- **Purpose:** update system doctrine, skill docs, glossary, taxonomy, or design contracts.
- **Required input:** the doctrine gap and the surfaces it must reconcile (cross-references, terminology, examples).
- **Matter produced:** documentation, `SKILL.md` changes, glossary entries, schema doctrine (not code behavior).
- **Review surface:** β reviews for doctrinal coherence — cross-reference correctness, no orphaned new doctrine, no silent behavior change smuggled in under a docs label.
- **Closeout / projection:** closes like any other cell — a docs-only diff still requires review and receipt. **This cell (#570) is itself a worked example of the `doctrine` kind** — see its `gamma-scaffold.md` `cell_kind:` line.
- **FSM implication:** may be docs-only but still requires review and receipt; must not claim `status:review` without the doctrine diff actually present.

### 10. `experiment`

- **Purpose:** explore an uncertain design or technical path without committing to product behavior.
- **Required input:** an open question with no converged design yet.
- **Matter produced:** a spike result, a prototype, a recommendation, or a follow-up issue.
- **Review surface:** β reviews whether the exploration answered the question, not whether it shipped production code.
- **Closeout / projection:** projects a recommendation or a follow-up `issue_authoring` cell — an experiment cell's output must not silently become production implementation.
- **FSM implication:** should not silently become production implementation; a follow-on `implementation` cell is a separate, explicitly dispatched cell.

## FSM awareness: `cell kind -> valid matter -> allowed transition request`

The issue FSM (`cnos.issues/commands/issues-fsm`) should evaluate transitions differently depending on which cell kind is running against an issue. This table names the mapping at the doctrine level; it is **not yet consumed by any transition rule** — see "Observation, not enforcement" below.

| Cell kind | Valid matter | Allowed transition request |
|---|---|---|
| `issue_authoring` | issue body / issue pack draft | `selected gap -> issue draft -> issue reviewed -> status:ready` (not `status:todo -> status:in-progress -> status:review`) |
| `implementation` | branch / PR / diff / review-request | `status:todo -> status:in-progress -> status:review` |
| `repair` | targeted fix + repair contract | `status:changes -> status:todo` (requires `RepairContractPresent`) |
| `recovery` | recovery PR / STOP comment / requeue decision | `status:in-progress -> status:todo` (no matter) or `status:in-progress -> delta-recovery` (matter present, cnos#368) |
| `wave` | child-cell receipts / wave receipt | wave-contract-defined; not a single `todo -> in-progress -> review` step |
| `cleanup` | TSC report / drift inventory / cleanup PR or issues | may request child `issue_authoring` or `implementation` cells; the cleanup cell's own request follows `audit` or `implementation` shape depending on its matter |
| `audit` | inventory / drift report / proof table | `status:todo -> status:review` (PR optional; report artifact is the matter) |
| `release` | tag / release notes / deploy branch / release receipt | boundary-only; requires δ, not a normal β content review |
| `doctrine` | doc / SKILL.md / glossary / schema doctrine diff | `status:todo -> status:in-progress -> status:review` (same shape as `implementation`, docs-only matter) |
| `experiment` | spike result / prototype / recommendation | `status:todo -> status:review` with a recommendation or follow-up issue as matter, not a merge-to-main implementation |

## Recording point (AC8)

`cell_kind` is first recorded as a line in `.cdd/unreleased/{N}/gamma-scaffold.md`:

```markdown
**cell_kind:** `doctrine`
```

or the plain form:

```text
cell_kind: doctrine
```

This cell's own scaffold, `.cdd/unreleased/570/gamma-scaffold.md`, is the worked example — γ names the cell kind at dispatch time, before α ever starts producing matter. This is the **first canonical recording point** per the issue's own recommendation: "document the field and require new cells to record it in the CDD receipt / self-coherence surface. Do not make verifier enforcement mandatory yet." No existing cells are required to backfill `cell_kind`; its absence defaults to `implementation` (the current CDS dispatch cell kind) wherever it is observed.

**Future typed field (deferred, Phase B).** The generic-kernel + typed-refinement shape here mirrors the CUE unification pattern already used across `schemas/cdd/receipt.cue` (generic `#Receipt`) and `schemas/cds/receipt.cue` (`#CDSReceipt: cdd.#Receipt & {...}`). A future `schemas/cdd/cell.cue` could define `#Cell` and per-kind refinements as a closed disjunction discriminated by `cell_kind`, validated via `cdd-verify`'s `cuevet.go`. That CUE schema, its verifier enforcement, and any codegen of Go types from it are explicitly **out of scope for this doctrine-first cell** (see #570 Deferred list).

## Observation, not enforcement

`src/packages/cnos.issues/commands/issues-fsm` carries an observation seam for `cell_kind`: `FactSnapshot.CellKind{Observed, Source, DefaultedTo}` (`snapshot.go`). As of this cell, `fetch.go`'s live-assembly path (`assembleLive`) parses the `cell_kind:` line out of `.cdd/unreleased/{N}/gamma-scaffold.md` when present and sets `CellKind.Observed` / `CellKind.Source = "cdd_artifact"`. This is **observation only** — no transition rule in `table.go` consumes `CellKind`, so its value cannot change any FSM decision. `TestSeam_CellKindNotEnforced` locks this: evaluating the same facts under every defined cell kind must yield a byte-identical decision. Consuming `cell_kind` in transition guards (the "valid transition" column above becoming enforced behavior) is deferred to a future, explicitly scoped FSM Phase 2 issue.

## Cross-references

- [`CDD.md`](CDD.md) — the generic kernel this taxonomy refines.
- [`issue/SKILL.md`](issue/SKILL.md) — issue-authoring doctrine; states an issue can itself be a cell's matter.
- `docs/reference/governance/GLOSSARY.md` §"Cell" — glossary entry for the generic term.
- `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` — the `CellKind` observation struct.
- `src/packages/cnos.issues/commands/issues-fsm/fetch.go` — the `assembleLive` parse path.
