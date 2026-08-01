# CELL-KINDS — Legacy Cell-Kind Taxonomy, now Domain Vocabulary

**Status:** Doctrine-first (Phase A of [#570](https://github.com/usurobor/cnos/issues/570)). **Demoted (proposed, [#628](https://github.com/usurobor/cnos/issues/628)):** the "cell kinds" below are **matter/contract domains**, not canonical cell types — this document is now the **domain vocabulary**. Deferred: CUE schema, verifier enforcement, UI display, automatic inference, full wave/cleanup execution, `table.go` transition consumption.
**Placement:** `src/packages/cnos.cdd/skills/cdd/`
**Parent:** [`CDD.md`](CDD.md) — this document refines the generic kernel; it does not replace or re-derive it.

> **⚠ Framing demoted (proposed — [#628](https://github.com/usurobor/cnos/issues/628); ratifies the architecture note `docs/architecture/CELL-RUNTIME.md`).**
> The canonical unit of work is the **one CCNF kernel cell** (`COHERENCE-CELL-NORMAL-FORM.md`), not a "cell kind." The names below are **matter/contract domains** (the *fourth* orthogonal axis) — **not** cell types, and **not** the operating-scale **classes** WC/PC/CC (output-telos deployment shapes of the one kernel; see `docs/architecture/CELL-RUNTIME.md`). They are heterogeneous:
> - **matter domains** — `implementation`, `doctrine`, `audit`, `experiment`, `cleanup`, `planning`
> - **contract modes** — `repair`, `recovery`
> - **coordination / boundary** — `wave` (a scope), `release` (a δ boundary action)
>
> **Compat (legacy / adapter only):** legacy `cell_kind: X` maps to `cell_class` (`working|planning|cohering`) + `matter_domain: X`; the FSM `CellKind` observation seam keeps `cell_kind` as an alias. **New contracts must not introduce new canonical cell kinds** — a contract names a `cell_class` and a `matter_domain`. The file is kept in place (link stability); a physical rename to `CELL-DOMAINS.md` may follow once inbound links are swept. The definitions below stand as the **domain vocabulary**.

## Design rule: one kernel, multiple domains

[`CDD.md`](CDD.md) names the generic recursive coherence-cell kernel:

```text
contractₙ -> matterₙ -> reviewₙ -> receiptₙ -> verdictₙ -> decisionₙ
```

Every closed cell — regardless of what it produces — runs this same five-step closure. **A domain does not fork the kernel** (this rule was formerly stated as *"cell kind does not fork the kernel"*). A domain is a *typed refinement* of the kernel: it names what matter the kernel's `α.produce` step is expected to emit, what `β.review` means for that matter, what the `γ.close` receipt projects, and what boundary decisions `δ` may take. Do not create a separate lifecycle per domain — bind the existing kernel to a specific matter shape instead. A domain is **not** the cell's **class** (WC/PC/CC), which names its output telos; see `docs/architecture/CELL-RUNTIME.md`.

```text
domain -> expected matter -> review surface -> closeout/projection
```

## `issue kind` != `cell kind`

These are two different, non-interchangeable classifications:

- **`issue kind`** (e.g. `kind/process`, `kind/bug`, `kind/feature` — see `docs/development/issues/TAXONOMY.md`) is a **GitHub issue taxonomy label**. It classifies *what an issue is about*.
- **`matter_domain`** (legacy `cell_kind`; e.g. `implementation`, `issue_authoring`, `wave`) is a **domain declaration**. It classifies *what matter the coherence cell dispatched against that issue is expected to produce*. `cell_kind` is retained **only** as a compatibility alias; new contracts name a `cell_class` (`working|planning|cohering`) plus a `matter_domain`.

An issue labeled `kind/feature` can be worked in the `implementation` domain (the common case) or, earlier in its life, be *itself* the matter of an `issue_authoring`-domain cell that drafted it. The GitHub label and the `matter_domain` value travel on different axes and must not be conflated.

## The 11 legacy domains

The eleven entries below are the **matter/contract domain vocabulary** (legacy name: "cell kinds"). Each names an expected matter shape, review surface, and closeout/projection — a *typed refinement* of the one kernel, not a distinct cell type. A cell also carries a **class** (WC/PC/CC) orthogonal to its domain.

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

### 11. `planning`

- **Purpose:** produce issue/wave contracts as matter. Planning is work, not "chat before work" — therefore a cell.
- **Required input:** a scope, sequencing question, or wave-shaped gap that has not yet been decomposed into executable issue(s).
- **Matter produced:** issue body · wave master issue · subissue list · sequencing · scope · ACs · STOP conditions · dispatch plan.
- **Lifecycle (α/β/γ):** α produces the plan / issue / wave. β reviews the plan for coherence, scope, proof, and design boundary. γ closes the planning cell: posts/links issues, verifies labels, records the sequence, captures decisions, and binds learning/follow-ups into the receipt.
- **Role-vs-model note:** the Greek letter names the **role**, not the model/agent. κ (human-interface herald) commonly acts as **α of a planning cell**; a reviewing activation acts as **β**; either may be a different activation. Do **not** mint a new Greek identity per model — the role is the invariant.
- **Relationship to neighbors:** `issue_authoring` produces one executable issue; `wave` composes child-cell receipts; `planning` produces the *plan/wave contract + sequencing* and may spawn `issue_authoring` + `implementation` children. β resolves exact boundaries during review.
- **Worked example:** operator review of #570 surfaced two gaps (this cell kind, and the mandatory learning section) that were not cleanly captured by `issue_authoring` (a single issue) or `wave` (child-receipt composition) — a planning-shaped pass would have filed a scoped issue with explicit AC/scope up front, then handed off to a `doctrine`-kind `implementation` cell, illustrating the planning → issue_authoring/implementation handoff this entry formalizes.
- **FSM implication:** a planning cell may **create or update issues**, but does **not** imply immediate dispatch unless the plan explicitly says so (mirrors `issue_authoring`).

## FSM awareness: legacy `cell_kind` (= `matter_domain`) -> valid matter -> allowed transition request

The issue FSM (`cnos.issues/commands/issues-fsm`) evaluates transitions using the observed legacy `cell_kind` field (now read as the `matter_domain`). This table names the mapping at the doctrine level; it is **not yet consumed by any transition rule** — see "Observation, not enforcement" below.

| Domain (legacy `cell_kind`) | Valid matter | Allowed transition request |
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
| `planning` | issue body / wave master issue / subissue list / sequencing / dispatch plan | `selected gap -> plan contract -> status:ready` (mirrors `issue_authoring`); does not imply `status:todo -> status:in-progress` unless the plan explicitly says so |

## Recording point (AC8)

The legacy `cell_kind` field (now the `matter_domain`) is recorded as a line in `.cdd/unreleased/{N}/gamma-scaffold.md`:

```markdown
**cell_kind:** `doctrine`
```

or the plain form:

```text
cell_kind: doctrine
```

This cell's own scaffold, `.cdd/unreleased/570/gamma-scaffold.md`, is the worked example — γ names the domain at dispatch time, before α ever starts producing matter. This is the **first canonical recording point** per the issue's own recommendation: "document the field and require new cells to record it in the CDD receipt / self-coherence surface. Do not make verifier enforcement mandatory yet." No existing cells are required to backfill `cell_kind`; its absence defaults to `implementation` (the current CDS dispatch domain) wherever it is observed. `cell_kind` is retained as the on-disk field name for compatibility; new contracts additionally carry a `cell_class`.

**Future typed field (deferred, Phase B).** The generic-kernel + typed-refinement shape here mirrors the CUE unification pattern already used across `schemas/cdd/receipt.cue` (generic `#Receipt`) and `schemas/cds/receipt.cue` (`#CDSReceipt: cdd.#Receipt & {...}`). A future `schemas/cdd/cell.cue` could define `#Cell` and per-domain refinements as a closed disjunction discriminated by `matter_domain` (plus the orthogonal `cell_class`), validated via `cdd-verify`'s `cuevet.go`. That CUE schema, its verifier enforcement, and any codegen of Go types from it are explicitly **out of scope for this doctrine-first cell** (see #570 Deferred list).

## Mandatory terminal learning section

Every terminal cell closeout MUST include a `learning` (equivalently `epsilon_observations`) section — a general rule across **all** domains, not just `planning`. This is not a giant postmortem; it is a small, always-present artifact:

```yaml
learning:
  observations:        # what surprised us
    - ...
  process_deltas:      # what should change next time
    - ...
  reusable_patterns:   # what can become doctrine / tooling
    - ...
  followups:           # issues to file / link
    - issue: ...
  operator_burden:     # what required unnecessary operator intervention
    - ...
```

This is **ε captured by γ**: ε observes patterns across runs; γ binds the observation into the receipt at closeout (see `gamma/SKILL.md` — γ's closeout responsibility now states this explicitly). Every terminal closeout answers: *what did we produce? what did we validate? what did we decide? what did we learn? what should change?* Learning becomes part of the system, not a chat side effect.

**Required by doctrine; mechanical verifier enforcement deferred to a follow-up** — consistent with this document's own "document first, enforce later" stance (see "Observation, not enforcement" below, and #570's original doctrine-first framing). No existing cells are required to backfill a `learning` section retroactively.

## Process self-improvement loop (cnos#640)

Capturing a `learning:` block (above) is not the same as preventing the failure it names from recurring. [cnos#614](https://github.com/usurobor/cnos/issues/614) captured an RCA for a design-first issue's body-hold prose drifting out of sync with its `status:*` label; [cnos#633](https://github.com/usurobor/cnos/issues/633) hit the **identical** failure in the same session because the #614 learning lived only as a captured observation, never as a mechanism — the Kernel's own §1.2 names this precisely: *"'Won't repeat' without a mechanism is not a fix."* [cnos#640](https://github.com/usurobor/cnos/issues/640) is the mechanization of the #614/#633 recurrence itself, and its seed content is κ's meta-comment on that issue (["Meta: how does process self-improvement actually work here?"](https://github.com/usurobor/cnos/issues/640#issuecomment-4920935086), 2026-07-09) — the honest self-audit this section formalizes rather than re-derives.

The full loop has **five** named steps (κ's comment folds "detect-recurrence" in as step 3, between capture and mechanize):

```text
observe -> capture -> detect-recurrence -> mechanize -> verify-non-recurrence
```

| Step | What happens | Owner today |
|---|---|---|
| **Observe** | Friction surfaces during a firing, a review, or supervision. | Whoever encounters it — any role, any session. No single owner by design; observation is incidental to doing the work, not a scheduled pass. |
| **Capture** | The friction is written into the terminal closeout's `learning:` block (`observations` / `process_deltas` / `reusable_patterns` / `followups` / `operator_burden`) — see "Mandatory terminal learning section" above. | **γ** — already doctrine, landed by [cnos#614](https://github.com/usurobor/cnos/issues/614)/[cnos#630](https://github.com/usurobor/cnos/issues/630). γ binds the observation into the receipt at close-out; this is "ε captured by γ" per the section above. |
| **Detect recurrence** | A pass scans captured learnings + issue history for a repeated failure signature (e.g. the same `process_deltas` text across ≥2 closeouts, or a rising `operator_burden` count) and flags it before a third occurrence costs another firing. | **No current owner.** κ's comment names this explicitly: *"nobody/nothing actually runs it."* ε is the candidate role (`ROLES.md §4b`, "cross-cell receipt-stream observer") but ε is a role name, not a running process today — there is no scheduled pass, no CI gate, and no reconciler that reads `learning:` blocks across cells. **This is the gap.** Next step: **ε made concrete** — a periodic pass that scans closeout `learning:` blocks and issue history for repeated signatures. Filing the issue that builds that pass is a follow-up action for γ/κ, not part of this cell (building a full ε pass was explicitly out of scope for cnos#640's sizing — see its γ scaffold's scope guardrails); this section names the gap and the shape of its fix rather than asserting it is solved. |
| **Mechanize** | A repeated or high-burden signature is promoted to a `kind/process`-labeled issue whose acceptance bar is an MCA (a gate, a tool, a single-source-of-truth fix) — never a second prose reminder. | Whoever files and dispatches the promoted issue (operator, κ, or a future detect-recurrence pass once one exists). The bar itself — MCA required, "try harder" rejected — is Kernel §1.2/§2.3 doctrine, already binding; what is missing is the trigger (detect-recurrence, above), not the bar. |
| **Verify non-recurrence** | The shipped mechanism carries a test or gate that would fail if the pattern reappeared. Recurrence after mechanization is itself a signal the mechanism was wrong, not that the fix should be re-explained. | Whoever ships the mechanizing cell — this is standard α/β test-coverage discipline (`alpha/SKILL.md` §2.2 "tests must prove the actual claim"), applied to the specific failure signature the mechanism targets. cnos#640's own `cn issues dispatch` fixture tests (`src/packages/cnos.issues/commands/issues-dispatch/dispatch_test.go`) are a worked instance: they prove the #614/#633 body/label contradiction is corrected and cannot silently reappear unnoticed. |

**Honest state as of cnos#640:** observe and capture are established; mechanize's bar and verify's discipline are established; **detect-recurrence has no running owner.** cnos#640 mechanizes one specific recurring failure (body/label authorization drift) directly — it does not close the general detect-recurrence gap. A future occurrence of a *different* repeated failure signature will not be automatically caught until a detect-recurrence pass exists. This section documents that gap rather than asserting it is closed — an assertion that ε "now runs this" without a shipped ε pass would be the exact false-closure failure mode `issue/SKILL.md` names.

## Observation, not enforcement

`src/packages/cnos.issues/commands/issues-fsm` carries an observation seam for the legacy `cell_kind` field: `FactSnapshot.CellKind{Observed, Source, DefaultedTo}` (`snapshot.go`). As of this cell, `fetch.go`'s live-assembly path (`assembleLive`) parses the `cell_kind:` line out of `.cdd/unreleased/{N}/gamma-scaffold.md` when present and sets `CellKind.Observed` / `CellKind.Source = "cdd_artifact"`. This is **observation only** — no transition rule in `table.go` consumes `CellKind`, so its value cannot change any FSM decision. `TestSeam_CellKindNotEnforced` locks this: evaluating the same facts under every defined domain must yield a byte-identical decision. Consuming the field in transition guards (the "valid transition" column above becoming enforced behavior) is deferred to a future, explicitly scoped FSM Phase 2 issue.

## Cross-references

- [`CDD.md`](CDD.md) — the generic kernel this taxonomy refines.
- `docs/architecture/CELL-RUNTIME.md` — the WC/PC/CC deployment **classes** + generic runner these **domains** sit under; the orthogonal-axes correction (proposed; #627/#628).
- [`issue/SKILL.md`](issue/SKILL.md) — issue-authoring doctrine; states an issue can itself be a cell's matter; its "false closure" failure mode is what §"Process self-improvement loop" refuses to commit (asserting the detect-recurrence gap is closed when it isn't).
- `docs/reference/governance/GLOSSARY.md` §"Cell" — glossary entry for the generic term.
- `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` — the `CellKind` observation struct.
- `src/packages/cnos.issues/commands/issues-fsm/fetch.go` — the `assembleLive` parse path.
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §1.2 "Labels are the sole source of truth for dispatch readiness" + its D13 failure-mode entry ([cnos#640](https://github.com/usurobor/cnos/issues/640)) — the mechanize/verify worked example §"Process self-improvement loop" cites.
- `ROLES.md` §4b — ε's generic role doctrine ("cross-cell receipt-stream observer"); §"Process self-improvement loop" names ε as the candidate (not-yet-concrete) owner of the detect-recurrence step.
