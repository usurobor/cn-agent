# CCNF-O — Track A1 survey + decision doc

**Status:** v0.1 — pinned decisions for Track A1 of [cnos#405](https://github.com/usurobor/cnos/issues/405) (CCNF-O orchestration grammar + TSC coherence steering)
**Cycle:** [cnos#421](https://github.com/usurobor/cnos/issues/421) — first dispatchable sub of #405 post-#404 closure
**Date:** 2026-05-23
**Gate satisfied:** [cnos#404](https://github.com/usurobor/cnos/issues/404) handoff extraction wave closed 2026-05-22 (`9fc04a99`); `cnos.handoff` v0.1 ships canonical wire formats for cross-repo / dispatch / mid-flight / artifact-channel / receipt-stream.

This document is the **survey + decision artifact** that Tracks A2–A6 and Track B1 dispatch against. It pins five interlocking decisions and refines the sub-issue queue for the rest of #405. **It is not an implementation; it authors no schemas, no Go types, no validator changes.** Implementation begins at Track A2.

---

## 0. Frame

Per [cnos#405](https://github.com/usurobor/cnos/issues/405) §0 (post-#403 restructure), #405 has two tracks. **Track A — CCNF-O orchestration grammar** types dispatch, implementation contract, findings, issue graph, wave manifest, closeout chain, and composition operators. **Track B — TSC coherence steering** attaches TSC reports to accepted receipts, types `IssueProposal.v1`, types `RiskPolicy.v1`, and defines the coherence controller. The hard gate (#404 handoff extraction must land before Track A types dispatch/handoff surfaces) is now satisfied; Track A1 is dispatchable.

The two tracks share substrate (cells, receipts, V verdicts, δ boundary decisions) and consumer (the autonomous-engineer loop). They are separable along the **reason-to-change** axis named in `cnos.core/skills/design/SKILL.md §3.1`: Track A changes when the **composition grammar** changes; Track B changes when the **measurement function or steering signal** changes. The five pinned decisions below preserve that separation.

### What this survey does NOT do

- It does not author any CUE schemas. `schemas/ccnf-o/` does not exist; Track A2 creates it.
- It does not generate any Go types via `cue exp gengotypes`. Track A2 lifts schemas, then gengotypes.
- It does not modify any role-skill in `cnos.cdd` / `cnos.cds` / `cnos.cdr` / `cnos.handoff`. The handoff package's v0.1 surfaces are read, not re-pinned.
- It does not edit the #405 tracker body. Decisions are produced externally so a follow-on cycle can fold them in cleanly.
- It does not touch the CCNF kernel (`CDD.md` / `COHERENCE-CELL.md` / `COHERENCE-CELL-NORMAL-FORM.md`).

### What this survey DOES do

- Pins the name (`CCNF-O` vs `CCNF-X`).
- Classifies every Track-A-candidate surface as universal CCNF-O or realization-specific (CDS / CDR / future c-d-X).
- Classifies every higher-level form in #405 §9 as universal or realization-derived.
- Confirms Track B1 can design in parallel with Track A2+, names the shared touchpoints, and lists Track B1's dependencies on Track A.
- Pins the package location for CCNF-O.
- Refines the dispatch-brief shape for Tracks A2–A6 + B1.

---

## 1. Name pick

### Candidates

- **CCNF-O** — Coherence Cell Normal Form, **Orchestration overlay**. Names the relationship to the kernel (CCNF) and what the overlay formalizes (orchestration — composition of cells into cycles / waves / roadmaps / gates / joins).
- **CCNF-X** — Coherence Cell Normal Form, **Composition overlay** (X = composition, after the kernel's "no FAIL + accept" rule that composes axes geometrically). Names the relationship to CCNF and the **operation** (composition) rather than the consumer (orchestration).

Both names share the prefix `CCNF-` which correctly signals "an overlay on the kernel" rather than a separate calculus. The difference is what the suffix names: the **consumer** (orchestration) or the **operation** (composition).

### Decision: **CCNF-O** pinned

Pin **CCNF-O**. The operator's mild preference (#405 §0 ruling + #421 dispatch brief) holds; the survey confirms.

**Rationale.** CCNF-O names what the grammar **formalizes for** — orchestration: the runtime activity of taking cells produced by CCNF and composing them into the structures the protocol actually ships (cycles, waves, roadmaps, gates, joins, repairs). The kernel calculus (CCNF) already names the composition primitive at the per-cell scale ("a cell is a 5-tuple `(α, β, γ, V, δ)` that composes via the no-FAIL-and-accept rule"); CCNF-O names the **next stratum up** — the composition of cells. "Orchestration" is the discoverable consumer-facing term: it is what δ does at boundary-decision time, what γ does at wave-join time, what the operator does when launching a roadmap. "Composition" is internal to the grammar; "orchestration" is the externally-visible activity the grammar serves. Naming by the externally-visible activity matches the discoverability principle in `docs/papers/FOUNDATIONS.md` §3.4 ("cnos is the practice layer in which the theory becomes real") and the consumer-naming rule in `cnos.handoff/skills/handoff/dispatch/SKILL.md §1.1` (the dispatch prompt names the **worker role** the prompt is for, not the wire-format internals it carries).

A second reading: CCNF-X already exists as the **pattern** for naming substrate-specific overlays of CCNF (CCNF-CDS, CCNF-CDR are the obvious next siblings if a per-domain stratification were needed). Reserving `-X` for "any X consumer" leaves `-O` free for the specific orchestration stratum. Picking `-X` for orchestration would collide with the generic naming convention and force a rename if a future per-domain stratification needs the `-X` slot.

A third reading: the autonomy-level table in #405 §1 names L5 as "**Roadmap autonomy** — agent decomposes goals, launches subs, joins receipts, repairs failed branches" and L6 as "**Coherence-driven autonomy** — TSC measures every shipment; ε generates next issues from bottleneck leverage; protocol evolves itself". L5 is squarely orchestration. L6 is the coherence-steered orchestration loop. Naming the L5 substrate "CCNF-O" makes the autonomy-level mapping (L5 = CCNF-O completes; L6 = CCNF-O + TSC integrated) one short consonant: `O` for orchestration.

**Binding rule.** All Track A2–A6 and Track B1 dispatch-brief templates, sub-issue bodies, schema file names (`schemas/ccnf-o/...`), package names, and downstream documentation use **CCNF-O** as the canonical name. No `CCNF-X` aliases are emitted. Historical artifacts that named the open question with both candidates (#405 §13 table header "Candidate typing", #405 §14, this survey's §1) remain valid pre-decision records and are not edited.

---

## 2. Surface inventory matrix

Every surface a Track A or Track B sub-issue might touch is classified along two axes: **current home** (where the canonical doctrine lives today, post-#404 closure) and **universal CCNF-O vs realization-specific** (does the surface belong in CCNF-O across CDD/CDS/CDR/future c-d-X, or does it belong in a realization package?).

Three classification verdicts are possible per row:

- **U** — Universal CCNF-O surface. CCNF-O **types** this surface across all consumer realizations. Realization packages bind the type, do not redefine it.
- **R** — Realization-specific. Stays in the realization package (CDD / CDS / CDR / future c-d-X). CCNF-O may name an abstract slot; the realization fills it.
- **H** — Handoff-resident wire format. CCNF-O **may** type the wire format (lift the Markdown table into a CUE schema) but does not own it. `cnos.handoff` is the canonical home; CCNF-O types it as a downstream service to handoff's consumers.

The matrix below uses **Universal CCNF-O?** = `Y` when the verdict is U, `N` otherwise; **Realization-specific?** = `Y` when the verdict is R, `N` otherwise. H-class surfaces have **Universal CCNF-O?** = `Y (types handoff format)` and **Realization-specific?** = `N (handoff owns)`.

| # | Surface | Current home | Universal CCNF-O? | Realization-specific? | Rationale |
|---|---|---|---|---|---|
| 1 | Mode enum (`MCA` / `explore` / `design-and-build` / `docs-only`) | `cnos.cdd/skills/cdd/issue/SKILL.md §"Mode declaration and MCA preconditions"` | **Y** | N | Mode is a property of a cycle's relationship to its design+plan stability — universal across any c-d-X protocol because every domain has the same MCA-precondition question ("is design + plan committed and stable?"). The vocabulary is CCNF-O-typed; the **per-realization narrative** (what counts as "design" in CDS vs CDR) is realization-side. |
| 2 | Sizing predicate + five-factor heuristic | `cnos.cdd/skills/cdd/issue/SKILL.md §"Cycle scope sizing"` | **Y** | N | A cycle is in-scope, at-edge, or over-scope independent of the realization domain — the five factors (AC count; AC heterogeneity; review-round budget; MCA-precondition stability; substrate change cost) are domain-agnostic predicates. CCNF-O types the predicate; realizations bind the soft-AC-count tier per their domain. |
| 3 | Master/sub graph | `cnos.cdd/skills/cdd/issue/SKILL.md §"Master+subs pattern"` + `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (cross-repo variant) | **Y** | N | The DAG-of-cells shape (master → subs; subs gate-in to master close) is universal — every domain composes a tracker into sub-cells the same way. The **cross-repo** variant's wire format is handoff-owned; the **graph shape** itself is CCNF-O. |
| 4 | Dispatch-prompt schema | `cnos.handoff/skills/handoff/dispatch/SKILL.md` (canonical) | **Y (types handoff format)** | N (handoff owns) | The 6-element envelope (role / project / role-skill / issue / branch / Tier-3 skills) is a wire format `cnos.handoff` owns. CCNF-O may type the Markdown template into `schemas/ccnf-o/dispatch_prompt.cue` and emit Go via `cue exp gengotypes` for validator consumption — but the wire format stays at handoff. **Track A2 surface.** |
| 5 | Implementation-contract 7-axis schema | `cnos.handoff/skills/handoff/dispatch/SKILL.md §2.3` (canonical) | **Y (types handoff format)** | N (handoff owns) | The 7-row table (Language / CLI integration target / Package scoping / Existing-binary disposition / Runtime dependencies / JSON-wire contract / Backward-compat invariant) is a handoff wire format. CCNF-O types it as `#ImplementationContract` and exposes Go types to V. The **values** populated in each row are repo-convention-specific; the **schema** is universal. **Track A2 surface.** |
| 6 | Findings state machine + classification | `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6` + `cnos.cdd/skills/cdd/epsilon/SKILL.md §1` + `cnos.handoff/skills/handoff/receipt-stream/SKILL.md §1` | **Partial-Y** | **Partial-Y** | Split surface. The **per-finding shape** (Source / Class / Trigger / Description / Root cause / Disposition + per-disposition sub-fields) and the **disposition vocabulary** (`patch-landed` / `next-MCA` / `no-patch`) are universal CCNF-O — Track A4 types them as `#Finding`. The **Class vocabulary** (`cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap` vs CDR's `cdr-data-gate-gap` / `cdr-overclaim-gap` / etc.) is realization-specific — each c-d-X enumerates its own class set against an open enum slot CCNF-O reserves. |
| 7 | Issue-pack schema | `cnos.cdd/skills/cdd/issue/SKILL.md §"Minimal output pattern"` | **Y** | N | The 14-section issue-pack template (Title / Problem / Impact / Status truth / Source of truth / Scope / Cycle scope sizing / Acceptance criteria / Proof plan / Skills to load / Active design constraints / Related artifacts / Non-goals / Success-closure / Handoff checklist) is universal across c-d-X — the structure of "a quality-gated, dispatchable cycle issue" is domain-agnostic. CCNF-O types as `#IssuePack`. **Track A3 surface.** |
| 8 | Wave manifest schema | empirical (per-tracker, ad-hoc; nearest precedent: `.cdd/iterations/wave-2026-05-19-protocol-hygiene.md` + the wave bundles in `.cdd/iterations/`) | **Y** | N | The wave manifest (parent tracker / sub-cell list / gate predicates / join semantic / per-sub status) is universal — every roadmap composes the same way. Currently smeared across tracker issues and `.cdd/iterations/wave-*.md` files; Track A5 types as `#WaveManifest`. |
| 9 | Closeout / receipt chain | `cnos.cdd/skills/cdd/post-release/SKILL.md` (CDD-side runbook) + `cnos.cds/skills/cds/CDS.md §"Artifact contract"` (CDS Ownership matrix) + `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` (channel substrate) | **Partial-Y** | **Partial-Y** | Split surface. The **chain shape** (α-closeout → β-closeout → γ-closeout; the closeout-receipt → INDEX-row sequence; the cross-cycle linkage) is universal — CCNF-O types as `#CycleCloseout`. The **per-artifact filename set** (`alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `self-coherence.md`, `beta-review.md`, `cdd-iteration.md`, `gamma-scaffold.md`) is realization-specific and lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract"` for software-class cycles; CDR has a different filename set. **Track A4 surface.** |
| 10 | Mid-flight rescue (`gamma-clarification.md`) | `cnos.handoff/skills/handoff/mid-flight/SKILL.md` (canonical) | **Y (types handoff format)** | N (handoff owns) | The rescue file format (date / edit summary / affected ACs-non-goals-constraints / rationale), the cycle-branch SHA-transition wake-up semantic, and the cache-bust contract are handoff wire formats. CCNF-O may type as `#MidFlightClarification`. Operationally a Track A surface to lift into CUE; the wire-format ownership stays at handoff. **Not a v0.1 priority.** |
| 11 | Cross-repo state machine + bundle file sets | `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (canonical) | **Y (types handoff format)** | N (handoff owns) | The 8-event STATUS lifecycle (drafted / submitted / accepted / modified / rejected / landed / withdrawn / revised), the four directional cases (a/b/c/d), and the LINEAGE schemas per case are handoff wire formats. CCNF-O may type the STATUS state machine as `#CrossRepoStatus` and the bundle file sets as `#CrossRepoBundle`. **Not a v0.1 priority.** |
| 12 | Artifact channel (per-role `.cdd/unreleased/{N}/` write ownership) | `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` (canonical) | **Y (types handoff format)** | **Partial-Y** | The channel-shape invariants (directory pattern; per-role write ownership; sequential α→β→γ rule; frozen-snapshot rule; release-time move) are handoff wire formats CCNF-O may type. The **per-artifact filename instantiation** (which files a CDS cycle writes; the 13-stage Ordered flow) is realization-specific and stays at `cnos.cds`. **Not a v0.1 priority.** |
| 13 | Receipt-stream + INDEX.md aggregator | `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (canonical) + `.cdd/iterations/INDEX.md` (live data) | **Y (types handoff format)** | N (handoff owns) | The per-cycle `cdd-iteration.md` shape, the eight-column INDEX row format (`Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path`), the cadence rule (required when `protocol_gap_count > 0`; courtesy stub when 0), and the cross-repo trace bundle invariant are handoff wire formats. CCNF-O may type as `#ReceiptStream` + `#IterationIndexRow`. **Cross-references with the Finding state machine (#6 row); Track A4 + Track B1 share this surface for `TSCReport` attachment.** |
| 14 | TSCReport | external — `docs/papers/DECREASING-INCOHERENCE.md §"Per-shipment artifact contract"` (design-doc-only; not yet typed) | **Y** | N | A coherence measurement on an accepted artifact bundle (α / β / γ scores; bottleneck axis; axis evidence; unresolved ambiguity; next-fixes seed). Universal because TSC measures any CCNF-cell-produced artifact regardless of realization. **Track B1 surface** — types into `schemas/ccnf-o/tsc_report.cue` and binds to the receipt-stream by `receipt_ref`. |
| 15 | IssueProposal.v1 | nascent — `cnos#405 §6` design sketch only | **Y** | N | The L3-autonomy primitive: a receipt-grounded, finding-cited, TSC-report-cited proposal whose `bottleneck_axis` and `kind` (`pattern-repair` / `relation-repair` / `process-repair` / `bug` / `test` / `docs` / `schema` / `refactor` / `tracker`) match. Universal across c-d-X. **Track B2 surface** (gated on B1). |
| 16 | RiskPolicy.v1 | nascent — `cnos#405 §7` design sketch + empirical (the "always human-gated" surfaces list scattered across role-skill rules) | **Y** | N | The L4-autonomy gate: an enumeration of `dispatch_rules` × `protected_surfaces` against which an IssueProposal is auto-dispatched or escalated. Universal because risk is a property of the **proposal shape**, not the realization. **Track B5 surface** (gated on B2). |
| 17 | Composition operators (`seq` / `par` / `gate` / `repair` / `lift` / `join` / `route`) | nascent — `cnos#405 §9 "Higher-level forms"` design sketch only | **Y** | N | The combinators by which cells compose into cycles → waves → roadmaps. Universal CCNF-O by definition — CCNF-O **is** the composition grammar. **Track A5 surface.** |
| 18 | Boundary decision (BoundaryDecision schema) | `schemas/cdd/boundary_decision.cue` (typed; generated Go in `src/go/`); `cnos.cdd/skills/cdd/delta/SKILL.md §1` (doctrine) | **Y (already typed; cite-not-re-type)** | N | Already typed in `schemas/cdd/`. CCNF-O cites the existing schema as a CCNF-O type without re-authoring. Not a Track A2+ priority — listed for inventory completeness. |
| 19 | Receipt schema | `schemas/cdd/receipt.cue` (typed; generated Go) + per-domain `schemas/cds/receipt.cue` + `schemas/cdr/receipt.cue` (typed; dispatched by `protocol_id`) | **Y (already typed; cite-not-re-type)** | **Partial-Y** | The generic-receipt + per-realization-receipt split already exists per #388 (Phase 2.5). CCNF-O cites the existing pattern. The `protocol_gap_count` field on the receipt is what triggers the receipt-stream cadence rule (#13 above). Not a Track A2+ priority — listed for inventory completeness. |
| 20 | Plan validation predicate (`cn cdd validate-plan`) | nascent — does not exist yet | **Y** | N | A CUE-typed plan (an instance of `#WaveManifest` or `#Roadmap`) must `cue vet` before any agent starts. The validator command is new. **Track A6 surface.** |

**Aggregate count:** 20 surfaces. 12 are universal CCNF-O (`U`); 7 are handoff-resident wire formats CCNF-O may type without owning (`H`); 1 (Finding state machine + Closeout chain) is split — universal in shape, realization-specific in vocabulary.

**Surfaces explicitly NOT enumerated** (out of CCNF-O scope):

- Harness substrate (dispatch invocation primitive `cn dispatch`; polling primitives; identity rotation; observability flags) — stays at `cnos.cdd/skills/cdd/harness/SKILL.md`. CCNF-O composes plans; harness executes them.
- Per-realization role-skill procedure (γ's branch-preflight; α's intake; β's review CLP). Stays in each realization's role-skill files.
- CCNF kernel doctrine (`CDD.md` / `COHERENCE-CELL.md` / `COHERENCE-CELL-NORMAL-FORM.md`). Untouched — CCNF-O sits **above** the kernel, not inside it.

---

## 3. Higher-level form classification

The six higher-level forms in [cnos#405](https://github.com/usurobor/cnos/issues/405) §9 are the structures CCNF-O composes from cells. Each is classified along the same Universal vs Realization-derived axis.

### 3.1 SimpleCycle

**Definition.** A single CCNF cell traversed once: `Cell → V → BoundaryDecision`. The 5-tuple `(α, β, γ, V, δ)` produces matter, review, receipt; V verdicts the receipt; δ records the boundary decision. No composition; the smallest unit CCNF-O can name.

**Classification.** **Universal.**

**Rationale.** A SimpleCycle is the CCNF kernel atom — the unit CCNF-O composes. Any c-d-X realization produces SimpleCycles by definition (the kernel is substrate-independent per `docs/papers/CCNF-AND-TYPED-TRUST.md §3.1`). CCNF-O types it as `#SimpleCycle` with three sub-types it composes: `#Cell`, `#Verdict`, `#BoundaryDecision`. The first two already have CUE schemas (`schemas/cdd/`); CCNF-O imports them.

### 3.2 SubstantialCycle

**Definition.** A SimpleCycle augmented with per-role evidence (α-closeout, β-review + β-closeout, γ-scaffold + γ-closeout, self-coherence, cdd-iteration when `protocol_gap_count > 0`) on the artifact channel. The form a real production cycle actually ships.

**Classification.** **Universal (in shape) + realization-derived (in artifact set).**

**Rationale.** The **shape** — that a real-world cycle augments the kernel cell with per-role evidence on a sequential α→β→γ channel that freezes on merge — is universal. CCNF-O types the shape as `#SubstantialCycle` with a `closeout_chain: #CycleCloseout` field. The **per-role artifact filename set** is realization-specific: CDS's set lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix"`; CDR's set lives in `cnos.cdr/skills/cdr/CDR.md`; future c-d-X enumerate theirs. CCNF-O reserves the slot; the realization fills it. This matches the universal-shape + realization-fills split in the surface inventory matrix row #9 above.

### 3.3 Wave

**Definition.** A set of SubstantialCycles dispatched in parallel under one parent tracker, with a join over their receipts at wave close: `par(SubstantialCycle*) → join(receipts) → wave_closeout`. The cnos#388-wave / cnos#403-wave / cnos#404-wave (just closed) are anchors.

**Classification.** **Universal.**

**Rationale.** A wave is a composition operator over SubstantialCycles. The composition primitive (`par`) and the join semantic (collect-receipts; verify-each; emit-wave-closeout when all subs accepted) are universal CCNF-O constructs. CCNF-O types as `#Wave` with `subs: [...#SubstantialCycle]`, `join: #JoinPredicate`, `wave_closeout: #CycleCloseout`. The wave manifest (the document that names the cycles before they ship) is the persistent representation; #405 §13 names `#WaveManifest` as a Track A5 schema; #421 inventory row #8 confirms.

### 3.4 Roadmap

**Definition.** A DAG of Waves connected by gates (per-phase gate predicates that consume one wave's receipt-set and admit the next wave): `DAG(Wave*) + gates + phase_receipts`. The cnos#366 → cnos#388 → cnos#403 → cnos#404 → cnos#405 chain is the canonical example.

**Classification.** **Universal.**

**Rationale.** A Roadmap is a composition operator over Waves with explicit gate semantics (`gate` is one of the seven CCNF-O composition operators per surface inventory row #17). CCNF-O types as `#Roadmap` with `waves: [...#Wave]`, `gates: [...#GatePredicate]`, `phase_receipts: [...#PhaseReceipt]`. The realization binding (a #405 issue body in GitHub, a `.cdd/iterations/wave-*.md` file, an offline planning doc) is the consumer's call. CCNF-O types the structure; consumers serialize.

### 3.5 CoherenceMeasuredCell

**Definition.** An accepted SubstantialCycle augmented with a TSCReport that scores the post-cell artifact bundle on (α_TSC, β_TSC, γ_TSC) axes and names the bottleneck: `AcceptedCell → TSC(target) → TSCReport(α, β, γ, bottleneck_axis, next_fixes)`. The bridge primitive between Track A (composition grammar) and Track B (measurement + steering).

**Classification.** **Universal.**

**Rationale.** TSC is substrate-independent — it measures any artifact bundle, regardless of which c-d-X produced it (per `docs/papers/DECREASING-INCOHERENCE.md §"Core model" → "TSC measures the accepted target"`). CCNF-O types as `#CoherenceMeasuredCell` with `cell: #SubstantialCycle`, `tsc_report: #TSCReport`. The TSCReport schema lives in `schemas/ccnf-o/tsc_report.cue` per #405 §3; Track B1 authors it. **This is the Track A ↔ Track B integration surface.** Note: the **CCNF α/β/γ role stations** (cell-side) and the **TSC α/β/γ measurement axes** (TSC-side) share letters but are distinct concepts per #405 §4.1; CCNF-O disambiguates via the `_CDD` / `_TSC` suffix convention when the two appear in the same expression.

### 3.6 AutonomousLoop

**Definition.** The L6-autonomy form: a CoherenceMeasuredCell feeds a coherence controller that maps the bottleneck to an IssueProposal under RiskPolicy.v1 gating; an accepted proposal dispatches the next cell; the cell ships and the loop continues. `CoherenceMeasuredCell → ε.propose(bottleneck, goal_pressure) → IssueProposal → V(proposal) → RiskPolicy.gate → Cell.dispatch → (loop)`.

**Classification.** **Universal.**

**Rationale.** The AutonomousLoop is the highest-stratum CCNF-O form — it ties the composition grammar (Track A) to the measurement function (Track B) to the issue generator (ε; Track B4) to the risk-gated dispatcher (Track B5). The loop itself is the closure of #405 §1's autonomous-engineer architecture. CCNF-O types as `#AutonomousLoop` referencing `#CoherenceMeasuredCell`, `#IssueProposal`, `#RiskPolicy`. Realizations bind goal-pressure inputs and observation cadences; the **loop structure** is universal. This is the final shape Track A6 + Track B6 jointly verify.

### 3.7 Summary table

| Form | Classification | Owning Track | CCNF-O schema |
|---|---|---|---|
| SimpleCycle | Universal | A2 (imports existing schemas) | `#SimpleCycle` |
| SubstantialCycle | Universal (shape) + realization (artifact set) | A4 | `#SubstantialCycle` + per-realization fill |
| Wave | Universal | A5 | `#Wave` |
| Roadmap | Universal | A5 | `#Roadmap` |
| CoherenceMeasuredCell | Universal | A4 ↔ B1 (joint design) | `#CoherenceMeasuredCell` |
| AutonomousLoop | Universal | A6 ↔ B6 (joint verification) | `#AutonomousLoop` |

All six forms are universal. No form is realization-only. The only realization-specific content surfaces in SubstantialCycle's per-role artifact filename set — and that is the same realization-specific surface row #9 of the inventory matrix called out.

---

## 4. TSC integration v0.1 scope

### Headline decision: **Track B1 can dispatch in parallel with Track A2.**

This confirms the operator's standing ruling (#405 §0 sequence step 6, #405 §12, #421 dispatch brief "Active design constraints"). The survey finds no structural obstruction. Track B1 designs the **TSC report attachment** surface (where post-cell TSC reports live; which targets get measured; how reports cite receipts); Track B2+ executes against Track B1's pins **only after Track A1 has closed** (this cycle). With Track A1 closing, both Track A2 (typing dispatch + implementation-contract schemas) and Track B1 (TSC report attachment design) become dispatchable; they may run in parallel without contention.

### Why parallel design is safe

Track A2 and Track B1 do not share files in their v0.1 scope:

- **Track A2** writes `schemas/ccnf-o/dispatch_prompt.cue` + `schemas/ccnf-o/implementation_contract.cue`, generates Go via `cue exp gengotypes` into `src/go/internal/cdd/` (or per Track A2's package-location pin), and wires the new types into `cn cdd verify`. It does not touch the receipt schema, the post-release runbook, or any TSC artifact.
- **Track B1** designs the **per-shipment artifact contract** (per `docs/papers/DECREASING-INCOHERENCE.md §"Per-shipment artifact contract"`): names where the post-cell TSC report file lives (likely `.tsc/{X.Y.Z}/{N}.json` or a similar parallel-to-`.cdd/iterations/` pattern), which targets get measured (release tag? per-cycle artifact bundle? both?), and how the report cites the cycle's receipt by `receipt_ref`. Track B1 is **design**, not implementation; it does not author the `TSCReport` CUE schema (that is Track B2 under a different operator ruling, or a B1.1 sub).
- The shared substrate (the artifact channel, the receipt's `receipt_ref` field, the receipt-stream INDEX) is already settled in `cnos.handoff/skills/handoff/{artifact-channel,receipt-stream}/SKILL.md` and `schemas/cdd/receipt.cue`. Track B1 cites these; it does not modify them.

### Design touchpoints (Track A ↔ Track B)

The two tracks meet at three identifiable surfaces:

1. **The CoherenceMeasuredCell form** (§3.5 above). A SubstantialCycle (Track A4) joined to a TSCReport (Track B1's design output, eventually B2's typed schema) by `receipt_ref`. Track A4 reserves the `tsc_report?: #TSCReport` field on `#SubstantialCycle` (optional in v0.1; required at the AutonomousLoop closure); Track B1 + B2 fill it.
2. **The receipt-stream / iteration aggregator** (surface matrix row #13). The per-cycle `cdd-iteration.md` and the `.cdd/iterations/INDEX.md` aggregator are the cross-cycle observable for ε (Track B4) and the TSC controller. Track B1's report-attachment design adds a parallel `.tsc/` directory (or extends `cdd-iteration.md`'s schema with a `tsc_report_ref:` field) — Track B1 picks the binding. Track A4's `#Finding` typing must not block Track B1's choice; conversely, Track B1's design must not redefine the receipt-stream wire format (handoff owns it).
3. **The IssueProposal pipeline (B2 → B4)** — gated on Track A1 (this cycle), Track A4 (findings state machine + closeout-receipt chain), and Track B1 (TSC report attachment). The proposal's `source` block (per #405 §6) carries `receipt_ref`, `finding_ref`, `originating_issue`, `tsc_report_ref`, `bottleneck_axis`, `leverage`. CCNF-O types `#IssueProposal` (Track B2); Track A4 makes `#Finding` and `#CycleCloseout` available as schema imports; Track B1 makes `#TSCReport` available; the four schemas compose to type `#IssueProposal`'s `source` block.

### Track B1 dependencies on Track A

Track B1 needs these Track A1 (this cycle) decisions pinned before it can author its design doc:

- **Name pinned** (§1 above) — Track B1 emits TSC reports against artifact bundles. The bundle's identity refers to `CCNF-O cell n`, not `CCNF-X cell n`. Track B1's design citations must use the chosen name. **Pinned: CCNF-O.**
- **Package location pinned** (§5 below) — Track B1 cites the `#TSCReport` schema from the CCNF-O package, e.g. `cnos.ccnf-o/schemas/tsc_report.cue` (depending on §5's pin). **Pinned: see §5.**
- **TSCReport classification confirmed Universal** (§2 row #14) — Track B1 designs against a universal type, not a per-realization type. **Confirmed.**
- **CoherenceMeasuredCell classification confirmed Universal** (§3.5) — Track B1's attachment binds against a universal substrate. **Confirmed.**

Track B1 does **NOT** depend on Track A2 (typed dispatch / implementation-contract). The two tracks are independent v0.1-scope-wise; Track A2 is **not** a Track B1 gate.

### What is explicitly **out of v0.1 scope** for the Track A ↔ Track B integration

- The coherence controller's runtime (Track B4) — gated on Track B3 + Track A4.
- The `RiskPolicy.v1` schema (Track B5) — gated on Track B2.
- The `cn cdd validate-plan` command for CCNF-O plans (Track A6) — gated on Track A5.
- Auto-dispatch of L4+ IssueProposals — gated on Track B5 + RiskPolicy.

The v0.1 deliverable for the integration is: **TSC report exists as an attached artifact on every accepted CCNF-O-shaped cycle; ε can read it; no automatic action is taken on it.** That is sufficient for L3-autonomy testability.

---

## 5. Package location

### Three options

| Option | Form | Verdict |
|---|---|---|
| **(a)** `cnos.cdd/orchestration/` (sub-directory inside cnos.cdd) | Tight coupling; CCNF-O lives inside the CCNF-kernel package | **Rejected.** Rebuilds the boundary mistake the #404 wave just corrected — composition has a different reason to change than the kernel calculus. Per `cnos.core/skills/design/SKILL.md §3.1`, one reason to change per boundary; per #404's experience, accreted-doctrine-inside-the-kernel-package is the recurring failure mode the discipline must structurally resist. |
| **(b)** New `cnos.{ccnf-o\|orchestration\|grammar\|...}` package | Clean peer to `cnos.cdd` / `cnos.cds` / `cnos.cdr` / `cnos.handoff` | **Pinned (see below for name).** |
| **(c)** Fold into `cnos.handoff` | CCNF-O surfaces share namespace with handoff wire formats | **Rejected per #405 §15.** CCNF-O composes cells; handoff transports them. Different reasons to change; the boundary `cnos.handoff/skills/handoff/HANDOFF.md` §"Boundary vs CCNF-O" already declares is structurally load-bearing — folding would re-mix what the package contract just separated. |

### Decision: **option (b) with package name `cnos.ccnf-o`** pinned

Pin **`cnos.ccnf-o`** as the package name.

**Rationale.** Three considerations weighted:

1. **The naming convention.** The repo's existing protocol packages use the `cnos.{protocol-name}` pattern (`cnos.cdd`, `cnos.cds`, `cnos.cdr`, `cnos.handoff`). The name CCNF-O was pinned in §1 above as the canonical protocol-name. The package name follows: `cnos.ccnf-o`. (The hyphen is consistent with `cnos.ccnf-o` reading naturally as "the CCNF-O package".) Alternative names considered:
   - `cnos.orchestration` — descriptive but loses the "CCNF-overlay" relationship. The package's name should signal that it is the orchestration **overlay on the CCNF kernel**, not a generic orchestration concept.
   - `cnos.grammar` — too generic. The package contains the orchestration grammar specifically, not grammar generally.
   - `cnos.compose` — overweights the composition operators (one piece of the package, not its whole shape).
   - `cnos.ccnf-o` — names the protocol; carries the CCNF-overlay relationship in the prefix; consistent with peer packages.

2. **The reason-to-change.** This package will change when the **composition grammar** changes — when a new combinator (e.g., a `retry` or `compensate` operator) is added; when the higher-level forms gain a new shape (e.g., a `DistributedRoadmap` for cross-org coordination); when the wave-manifest schema evolves. It will **not** change when CCNF kernel doctrine changes, when handoff wire formats change, when a domain protocol's lifecycle shape changes, or when TSC's measurement function changes. One reason to change; clean boundary.

3. **The discoverability.** A future agent searching for "where does the wave manifest schema live?" finds `cnos.ccnf-o/schemas/wave_manifest.cue`. Searching for "where is the composition grammar?" finds `cnos.ccnf-o/skills/ccnf-o/SKILL.md`. Searching for "where do I extend the orchestration overlay?" finds the package by its prefix-match-to-the-protocol-name. The package's home is mnemonically recoverable from the name pinned in §1.

### Package shape sketch (Track A2 authors; not this cycle)

```
src/packages/cnos.ccnf-o/
├── README.md                     # package contract
├── cn.package.json               # cn package manifest
├── skills/
│   └── ccnf-o/
│       ├── SKILL.md              # loader
│       ├── CCNF-O.md             # canonical doctrine surface (peer of CDD.md / CDS.md / CDR.md / HANDOFF.md)
│       └── {sub-skill folders as Tracks A2–A6 ship surfaces}
├── schemas/
│   ├── dispatch_prompt.cue       # Track A2
│   ├── implementation_contract.cue  # Track A2
│   ├── issue_pack.cue            # Track A3
│   ├── finding.cue               # Track A4
│   ├── cycle_closeout.cue        # Track A4
│   ├── wave_manifest.cue         # Track A5
│   ├── roadmap.cue               # Track A5
│   ├── tsc_report.cue            # Track B1/B2 (cross-track)
│   └── ...
└── docs/
    └── ...                       # internal design notes
```

**Note on schema directory.** The dispatch brief and #405 §3 reference `schemas/ccnf-o/` as the schema home. The above sketch nests schemas **inside the package** rather than at repo-root `schemas/ccnf-o/`. Track A2 picks one. Both patterns exist in the repo today (per-package schemas under `src/packages/{pkg}/schemas/` vs centralized `schemas/{domain}/`). The dispatch brief's reference to `schemas/ccnf-o/` is non-binding for this survey; the location decision falls to Track A2 under its own implementation-contract pin. **AC11 of #421 only requires that `schemas/ccnf-o/` not exist as a directory at this cycle's close — it does not pre-bind the directory's eventual location.**

### Where the coherence controller lives

Per #405 §15, the coherence controller's eventual location is "likely `src/packages/cnos.core/skills/coherence/SKILL.md` or a new `cnos.coherence` package". This is a Track B4 decision, not a Track A1 decision. This survey **does not** pin the controller's location; it pins only the orchestration-grammar package's location. Track B4's dispatch brief will revisit the question once Track B3 (IssueProposal validation) has landed and the controller's interface contract is clearer.

---

## 6. Sub-issue queue (refined)

Per #405 §14 + §14b. Each paragraph below pins the shape the next-cycle dispatch brief will carry. **Each sub may dispatch only after its gates are satisfied** (named per paragraph). All sub-issue bodies will use **CCNF-O** as the canonical name (§1 pin) and **`cnos.ccnf-o`** as the canonical package name (§5 pin).

### Track A2 — Type dispatch-prompt + implementation-contract schemas

**Pin shape.** Dispatch authors two new CUE schemas: `dispatch_prompt.cue` (the 6-element envelope from `cnos.handoff/skills/handoff/dispatch/SKILL.md §2.1`: role / project / role-skill / issue / branch / Tier-3 skills) and `implementation_contract.cue` (the 7-axis schema from `dispatch/SKILL.md §2.3`: Language / CLI integration target / Package scoping / Existing-binary disposition / Runtime dependencies / JSON-wire contract / Backward-compat invariant). Both lift their Markdown sources verbatim — no schema-shape revision. Run `cue exp gengotypes` to generate `DispatchPrompt` and `ImplementationContract` Go types into the chosen location. Wire the new types into `cn cdd verify` via a new validator predicate that loads any cycle's `## Implementation contract` section and verifies row-population. The dispatch brief pins: schema location (`src/packages/cnos.ccnf-o/schemas/` vs repo-root `schemas/ccnf-o/`); the Go-type integration point; the `cn cdd verify` wiring surface; the test fixture set. **Gate:** Track A1 (this cycle) closed. **Empirical anchor:** the post-#388 / #389 / #391 / #392 / #393 wave the dispatch surface was crystallized through.

### Track A3 — Type issue-pack + cell schemas

**Pin shape.** A3 types the issue-pack schema from `cnos.cdd/skills/cdd/issue/SKILL.md §"Minimal output pattern"` as `#IssuePack` (14 sections: Title / Problem / Impact / Status truth / Source of truth / Scope / Cycle scope sizing / Acceptance criteria / Proof plan / Skills to load / Active design constraints / Related artifacts / Non-goals / Success-closure / Handoff checklist), the Mode enum (`#Mode`: `MCA` | `explore` | `design-and-build` | `docs-only`), and the cycle-sizing predicate (`#CycleSizing`: in-scope / at-edge / over-scope + the five-factor heuristic encoded as predicate inputs). Also types `#Cell` as the kernel-cell projection that CCNF-O composes from, importing `schemas/cdd/receipt.cue` and `schemas/cdd/boundary_decision.cue` (already typed; cite-not-re-type). The dispatch brief pins: the Mode enum's exact case set; whether the sizing heuristic emits a tri-valued predicate or a scalar score; the `#IssuePack` validation rules (`MCA → design-path AND plan-path cited`; `over-scope → master+subs pointer present`); the integration point with `cn cdd verify`. **Gate:** Track A2 closed.

### Track A4 — Type findings state machine + closeout-receipt chain

**Pin shape.** A4 types the `#Finding` schema from `cnos.handoff/skills/handoff/receipt-stream/SKILL.md §1` (per-finding shape: `Source` / `Class` / `Trigger` / `Description` / `Root cause` / `Disposition` + per-disposition sub-fields), with the `Disposition` field as a closed enum (`patch-landed` | `next-MCA` | `no-patch`) and the `Class` field as an **open enum** populated per-realization (`cdd-skill-gap` | `cdd-protocol-gap` | `cdd-tooling-gap` | `cdd-metric-gap` for CDD; `cdr-data-gate-gap` | `cdr-overclaim-gap` | ... for CDR; future realizations enumerate their own). A4 also types `#CycleCloseout` as the closeout-receipt chain (α-closeout → β-closeout → γ-closeout; the closeout-receipt → INDEX-row sequence) with realization-fill slots for per-artifact filename sets. **The realization-fill pattern** is a v0.1 design decision A4 must pin: `#CycleCloseout` defines abstract slots (`alpha_closeout: #Artifact`, `beta_closeout: #Artifact`, etc.); realizations bind via `Artifact: {filename: string, owner: #Role, ...}` instances. The dispatch brief pins: the open-enum mechanism for `#Finding.Class`; the realization-fill pattern for `#CycleCloseout`; the reserved `tsc_report?: #TSCReport` slot on `#SubstantialCycle` for Track B1 ↔ Track A4 integration. **Gate:** Track A3 closed.

### Track A5 — Type wave manifest + master/sub graph + composition operators

**Pin shape.** A5 types the seven CCNF-O composition operators (`#Seq` | `#Par` | `#Gate` | `#Repair` | `#Lift` | `#Join` | `#Route`) per #405 §9 and the higher-level forms that compose from them: `#Wave` (par over `#SubstantialCycle` with a join), `#Roadmap` (DAG over `#Wave` with gates), `#MasterSubGraph` (one `#SubstantialCycle` (master) with `subs: [...#SubstantialCycle]` bound by an open join predicate), `#WaveManifest` (the persistent serialization of a `#Wave` — parent tracker / sub-cell list / gate predicates / per-sub status). A5 is the largest typing-surface Track A ships in v0.1 — it is the **substantive grammar** that justifies the package's existence. The dispatch brief pins: the seven-operator algebra (what `seq(A, par(B, C))` means semantically); the `#GatePredicate` schema (boolean expression over upstream `#PhaseReceipt`s); the `#JoinPredicate` schema (collect-and-emit semantic); the v0.1 expressiveness ceiling (e.g., no `retry` or `compensate` operators in v0.1; reserved for v0.2 if pressured). **Gate:** Track A4 closed.

### Track A6 — Wire CCNF-O validation into `cn cdd verify`; add `cn cdd validate-plan`

**Pin shape.** A6 wires the CCNF-O schemas (from A2 / A3 / A4 / A5) into `cn cdd verify` as a new verifier dispatch — `cn cdd verify` learns to dispatch per `protocol_id` (per the existing #388 generic-domain split) into a CCNF-O verifier for any cycle whose receipt names a CCNF-O-shaped artifact. A6 also adds a new command `cn cdd validate-plan` that takes a `#WaveManifest` or `#Roadmap` instance (in CUE source or rendered JSON) and runs `cue vet` against the CCNF-O schemas. The plan-validate command is the L5-autonomy gate: a plan that does not `cue vet` cannot be dispatched. The dispatch brief pins: the `cn cdd verify` dispatch table extension; the `cn cdd validate-plan` command surface (arguments, output format, exit codes); the integration with `cnos.core/skills/coherence/` if/when the coherence controller lands. **Gate:** Track A5 closed.

### Track B1 — Define TSC report attachment for accepted CCNF receipts

**Pin shape.** B1 is **design**, not implementation — output is a design doc + a small spec that A2-onwards execution can consume. B1 designs **where** the post-cell TSC report file lives (candidates: parallel `.tsc/{X.Y.Z}/{N}.json` directory; embedded `tsc_report:` block in `cdd-iteration.md`; sidecar `.cdd/releases/{X.Y.Z}/{N}/tsc_report.json`); **which targets** get measured (candidates: every accepted shippable cell; every release tag; every cycle close; selectable per realization); **how the report cites its receipt** (`receipt_ref: <path>` ↔ `tsc_report_ref: <path>` bidirectional invariant); and **what the report's per-target schema is** (the seven-field shape from `docs/papers/DECREASING-INCOHERENCE.md §"Per-shipment artifact contract"`: `alpha`, `beta`, `gamma`, `delta_*`, `bottleneck_axis`, `axis_evidence`, `unresolved_ambiguity`, `next_fixes`). The dispatch brief pins: the file-location choice (one option, with rationale); the targets enumeration; the bidirectional citation invariant; whether the v0.1 report schema is authored in B1 or deferred to B2. **Gate:** Track A1 (this cycle) closed. **May dispatch in parallel with Track A2.** Does not gate on Track A2.

---

## 7. Open questions deferred

Track A1 cannot resolve the following — each is flagged for the named downstream Track's dispatch brief.

### Q1 — Schema location: per-package vs repo-root

`src/packages/cnos.ccnf-o/schemas/{dispatch_prompt,implementation_contract,...}.cue` vs `schemas/ccnf-o/{dispatch_prompt,implementation_contract,...}.cue`? The repo has both patterns today (#388-wave used `schemas/cdd/`, `schemas/cds/`, `schemas/cdr/` at repo root; the per-package pattern is newer and not yet dominant). Track A2's dispatch brief decides under its own implementation-contract pin. **Track A2.**

### Q2 — The seven composition operators — full v0.1 or staged

#405 §9 names seven (`seq` / `par` / `gate` / `repair` / `lift` / `join` / `route`). Track A5 may type all seven in v0.1 or stage a subset (e.g., `seq` / `par` / `gate` / `join` in v0.1; `repair` / `lift` / `route` in v0.2 once the first three have empirical anchors). **Track A5.**

### Q3 — Realization-fill mechanism for `#CycleCloseout`

The realization-specific per-artifact filename set (CDS has 7 close-out files; CDR has a different set) needs a CUE-native binding mechanism. Three candidates: (a) `#CycleCloseout` is parameterized — `#CycleCloseout[Realization]`; (b) realizations extend via `#CDSCycleCloseout: #CycleCloseout & {alpha_closeout: "alpha-closeout.md"}`; (c) the abstract slots stay abstract and a runtime registry binds per `protocol_id`. **Track A4.**

### Q4 — `#Finding.Class` open-enum syntax in CUE

CUE supports closed disjunctions out of the box but expressing an **open-enum-with-known-values** (a set including the named CDD classes and CDR classes and a wildcard for future realizations) requires picking a CUE idiom. The mechanism affects how `cn cdd verify` validates findings whose `Class` is a never-before-seen value (warn? fail? log?). **Track A4.**

### Q5 — `#TSCReport` schema authorship — Track B1 or Track B2

Track B1 is design (per #405 §14b). The TSC report's CUE schema could land as part of Track B1 (design doc + first schema) or be deferred to Track B2 (schema first-class deliverable after B1's design pins). The operator (#421 dispatch brief) does not pre-decide; the question lands at **Track B1's dispatch brief**.

### Q6 — Wave manifest authorship surface

A `#WaveManifest` instance has to be **authored somewhere** before `cue vet` can validate it. Candidates: a top-of-tracker section in the parent issue body (machine-extracted); a sidecar `.cdd/iterations/wave-{slug}.yaml`; an in-repo `cnos.ccnf-o/plans/` directory. The choice affects who has write authority and how the manifest stays in sync with the issue body. **Track A5.**

### Q7 — Coherence controller location

#405 §15 names two candidates: `cnos.core/skills/coherence/` or a new `cnos.coherence` package. Track A1 does not pin (this is a Track B4 question, not a Track A question). **Track B4.**

### Q8 — When CCNF-O imports vs duplicates existing schemas

`schemas/cdd/receipt.cue` and `schemas/cdd/boundary_decision.cue` are already typed. CCNF-O composes against them. Does CCNF-O **import** them (CUE `import` directive) or **duplicate** them under `cnos.ccnf-o/schemas/`? Import preserves single-source-of-truth; duplication insulates CCNF-O from per-realization receipt-schema evolution. **Track A3.**

### Q9 — Whether to type the handoff wire formats now or in a later wave

Surfaces #4 / #5 / #10 / #11 / #12 / #13 of the inventory matrix are handoff-resident; CCNF-O **may** type them but does not own them. v0.1 priority per the surface matrix is: #4 (dispatch-prompt) and #5 (implementation-contract) in Track A2 (high pressure, codification cycle #393); #13 (receipt-stream) at Track A4 + Track B1 (cross-track integration); #10 (mid-flight), #11 (cross-repo), #12 (artifact-channel) deferred to v0.2 unless a downstream consumer (V; the coherence controller; an autonomous-engineer runtime) creates pressure to type them. **Deferred — revisit at v0.2.**

### Q10 — Whether the `_CDD` / `_TSC` suffix is the canonical disambiguation

#405 §4.1 names the α / β / γ collision between CDD role stations and TSC measurement axes and proposes the `α_CDD` / `α_TSC` prefix convention. The convention is used throughout the §4 essay and #405 body; this survey adopts it (§3.5). Whether the convention extends to schema field names (`#CoherenceMeasuredCell.tsc_alpha` vs `#CoherenceMeasuredCell.alpha_tsc` vs `#CoherenceMeasuredCell.tsc.alpha`) is a Track A4 ↔ Track B1 joint decision. **Track A4 ↔ Track B1.**

---

## 8. Closeout

This survey pins five decisions and refines the sub-issue queue. The decisions are:

1. **Name:** CCNF-O (§1).
2. **Surface inventory:** 20 surfaces classified along the universal / handoff-resident / realization-specific axis; 12 universal, 7 handoff-resident, 1 split (§2).
3. **Higher-level forms:** all 6 universal in shape; SubstantialCycle has a realization-specific artifact-set fill (§3).
4. **TSC integration v0.1 scope:** Track B1 dispatches in parallel with Track A2; touchpoints at CoherenceMeasuredCell, receipt-stream, IssueProposal pipeline; Track B1 depends on Track A1 (this cycle) only (§4).
5. **Package location:** new `cnos.ccnf-o` package, peer to `cnos.cdd` / `cnos.cds` / `cnos.cdr` / `cnos.handoff` (§5).

Six Tracks (A2 / A3 / A4 / A5 / A6 / B1) have refined dispatch-brief shapes (§6). Ten open questions are deferred to named downstream Tracks (§7).

With this survey landed, Tracks A2 and B1 are dispatchable. The remaining Tracks (A3 / A4 / A5 / A6 / B2 / B3 / B4 / B5 / B6) are gated per the chain in #405 §0 + §14 + §14b and the per-paragraph gates in §6 of this doc.

The cycle merges via `Closes #421`. #405 stays open; #405's tracker body may be updated in a follow-on cycle to fold this survey's decisions in (operator's call).

---

## Related artifacts

- [cnos#421](https://github.com/usurobor/cnos/issues/421) — this cycle's parent issue.
- [cnos#405](https://github.com/usurobor/cnos/issues/405) — Track A + Track B tracker.
- [cnos#404](https://github.com/usurobor/cnos/issues/404) — handoff extraction wave (closed 2026-05-22; the gate this cycle's dispatch satisfies).
- [cnos#393](https://github.com/usurobor/cnos/issues/393) — δ-as-architect; codified dispatch-prompt + implementation-contract (the most mature CCNF-O candidate surface; Track A2 target).
- [cnos#414](https://github.com/usurobor/cnos/issues/414) — `docs/papers/DECREASING-INCOHERENCE.md` merged (`bc29c009`); design-doc counterpart of #405.
- [`docs/papers/CCNF-AND-TYPED-TRUST.md`](../../papers/CCNF-AND-TYPED-TRUST.md) — the typed-trust precursor essay; the substrate decision for CCNF + CUE + V the CCNF-O package extends.
- [`docs/papers/DECREASING-INCOHERENCE.md`](../../papers/DECREASING-INCOHERENCE.md) — the Track B steering essay.
- [`docs/papers/FOUNDATIONS.md`](../../papers/FOUNDATIONS.md) — the C≡ / TSC / CTB / cnos stack.
- [`src/packages/cnos.handoff/`](../../../src/packages/cnos.handoff) — the wire-format package the CCNF-O package will type (without owning).
- [`src/packages/cnos.cdd/skills/cdd/issue/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/issue/SKILL.md) — Mode + sizing + issue-pack source.
- [`src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md) — findings runbook (pointer to handoff).
- [`src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md) — ε's CDD-side scope; finding-class instantiation.
