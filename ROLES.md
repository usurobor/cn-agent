# The Generic Role-Scope Ladder

This document explains the generic α/β/γ/δ/ε role structure that any
coherence-driven protocol instantiates. It is the single source of truth for
what the roles mean, why they form a scope-escalation ladder, and what an
instantiating protocol must declare to apply the pattern.

The immediate audience is anyone operating a c-d-X protocol (cdd, cdw, cdr,
cda, …) or designing a new one. The governing question this document answers:
what is α, and how does the same role structure apply to writing prose,
analysing data, and developing code without being re-derived each time?

---

## §1 The role ladder

Each role acts on the previous role's output as its frame. The roles are
nested, not parallel.

| Role | Verb        | Domain                                               |
|------|------------|------------------------------------------------------|
| α    | produces    | the matter (deliverable: code, prose, data, plan)    |
| β    | reviews     | α's matter                                           |
| γ    | coordinates | the α↔β loop within one cycle                        |
| δ    | operates    | the sequence of γ-loops across cycles                |
| ε    | iterates    | the δ-discipline itself                              |

The structural property that makes this a ladder rather than a flat list: each
role's frame is the previous role's content. α works on the deliverable. β
works on what α produced. γ works on the loop α and β form together. δ works
on the cadence and selection of γ's loops over time. ε works on whether δ's
discipline is itself coherent — whether the protocol produces better matter
over successive cycles, or whether it has drifted, accrued debt, or stopped
learning.

No role at level N can self-assess at level N+1. β cannot reliably judge
whether the α↔β loop is healthy from inside the loop — that cross-role view
is exactly what γ provides. γ cannot reliably judge whether the sequence of
cycles is coherent from inside a single cycle — that longitudinal view is what
δ holds. δ cannot reliably judge whether the protocol itself is learning —
that is ε's meta-level observation. The separation is structural, not
hierarchical: higher-order roles do not supervise lower-order roles; they
observe the processes those roles form.

The isolation rule — α and β must be separate within any single cycle —
follows directly from this structure. β reviewing its own matter is not order-1
observation; it is order-0 observation masquerading as order-1. The
independence is not a bureaucratic requirement; it is the mechanism by which
review adds information that α could not have produced.

> **Self-application footnote.** This document itself was produced by α,
> reviewed by β, coordinated by γ, dispatched by δ, and ε's work emerged from
> the very findings it surfaces. The pattern is not merely described here — it
> was exercised in the act of writing this doc.

---

## §2 Scope escalation as nested observation

Each role operates at a strictly higher order of observation than the role
below it. Order 0 is the matter itself — the code, prose, or output α
produces. Order 1 is β's review of that matter: β observes what α produced
and judges whether it closes the declared gap. Order 2 is γ's coordination of
the α↔β loop: γ observes the process by which matter is produced and reviewed,
adjusts it when stuck, and ensures the cycle closes cleanly. Order 3 is δ's
operation of the sequence of cycles: δ observes γ's cadence across many loops
and selects which gaps to close in which order. Order 4 is ε's iteration of
the protocol: ε observes whether δ's discipline is itself coherent — whether
it selects the right gaps, closes them durably, and produces a system that
learns from its own cycles rather than repeating the same class of error.

The stacking is not decorative. An intervention by a role at order N can only
improve coherence at order N or below. γ coordinating a poorly-scoped cycle
(order-2 work) cannot fix the matter itself (order-0); it can only improve the
process that produces the next cycle's matter. δ selecting a gap that addresses
a recurring failure pattern (order-3 work) can change which cycles happen, but
cannot retroactively improve γ's coordination of cycles already closed. ε
iterating the protocol (order-4 work) can change the rules under which future
cycles run, but cannot un-produce matter that past cycles shipped. Each order
is causally upstream of the ones below it across time, not within a single
cycle.

Bateson's learning levels and von Foerster's second-order cybernetics are the
ambient prior art. Bateson's Learning 0 through Learning III map directly to
orders 0–3: Learning 0 is fixed response within a fixed context (α producing
a defined deliverable), Learning I is correction within a fixed frame (β
reviewing against a stated oracle), Learning II is revision of the frame that
Learning I uses (γ adjusting the cycle process), Learning III is change of the
context that generates frames (δ operating the protocol across cycles). ε
corresponds to the reflexive move von Foerster names as second-order
cybernetics: the observer is part of the system observed, and the system that
examines its own observation loop is doing order-4 work. The analogy is not
exact — the roles are operational, not theoretical — but the structural
property is the same.

---

## §3 Instantiation contract

A protocol that instantiates this role ladder must declare six fields. Without
them the pattern is invoked but not grounded — readers cannot verify that the
roles are operating correctly or that the protocol is self-consistent across
cycles.

**Field 1 — Matter type.** What does α produce? The answer is
protocol-specific and concrete: source code, tests, and docs (cdd); editorial
prose (cdw); a research summary with citations (cdr); a quantitative analysis
(cda). The matter type is not vague ("α produces work") — it names the
artifact class, its canonical file or directory, and the authoring constraint
applied as generation constraints. The matter type determines what β reviews,
what γ coordinates the production of, and what ε assesses whether the protocol
is improving at producing.

**Field 2 — Review oracle.** How does β determine that α's matter closes the
declared gap? The oracle must be mechanically checkable — "it looks good" is
not an oracle. In cdd the oracle is acceptance criteria with concrete rg/wc
checks per AC. In cdw the oracle will be clarity checks and source-of-truth
alignment against named authority files. The oracle names the verifiable
criterion and the tooling or procedure used to check it. A protocol without a
stated oracle cannot distinguish "β approved this matter" from "β stopped
looking."

**Field 3 — γ close-out artifact.** What does γ write when a cycle closes?
This artifact carries the triage of α and β findings, the coherence delta the
cycle produced, and the next-move commitment. Naming it makes the close-out
obligation inspectable and prevents the γ role from closing a cycle without
committing to what was learned. In cdd this is `gamma-closeout.md` plus the
post-release assessment. The artifact is not optional — a cycle without γ
close-out has not closed; it has only stopped.

**Field 4 — δ cadence.** How often does δ select a new gap and dispatch a new
cycle? Daily, weekly, milestone-triggered, or demand-driven — the answer is
protocol-specific. The cadence is not required to be fixed, but it must be
stated. An unstated cadence is an unowned protocol: no one is accountable for
the sequence of cycles, and the protocol drifts by default toward whichever
work happens to be loudest.

**Field 5 — ε iteration cadence.** How often does ε assess whether the
protocol itself is coherent? What event triggers ε work? In cdd, ε writes
`cdd-iteration.md` when the cycle triage produces at least one finding tagged
`cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, or `cdd-metric-gap`.
The cadence is finding-triggered, not calendar-triggered. Other protocols may
choose differently, but the trigger must be named. A protocol where ε's work
is perpetually deferred is a protocol that cannot learn from its own cycles.

**Field 6 — Actor collapse rule.** Which roles may be held by the same person
or agent, and at what scale does separation become warranted? Collapsing ε
onto δ is common in small-protocol regimes where protocol-iteration volume
does not justify a dedicated reviewer of the protocol. Collapsing γ onto δ
(the operator-as-coordinator pattern) is normal when only two agents are
available. The one collapse that is never permitted: α and β in the same
cycle. The independence of review is the mechanism, not a nicety. The collapse
rule names which collapses are permitted, which are prohibited, and what
signal — finding volume, cycle count, or protocol complexity — warrants
splitting a collapsed role into two.

---

## §4 Hats vs actors: roles as behavioral contracts

Roles are behavioral contracts (hats), not entity slots (actors). Any agent can wear multiple hats across cycles. The constraint is which hats cannot be worn simultaneously *within a single cycle* — and the reason is always structural.

**Independence, not headcount.** The constraint that α≠β is not ceremonial separation or bureaucratic headcount. The constraint is independence: review that is not independent of implementation is not review. When β and α are the same actor, β cannot observe α's matter from outside α's frame — the essential property that makes review add information rather than repeat α's internal validation.

**Actor collapse rules are derivable.** Each protocol's collapse rules follow from asking: "which independence guarantees does this collapse destroy?" If the answer is "none that matter for this cycle's risk profile," the collapse is safe.

**Worked examples:**

- **δ=γ collapse (safe).** When the operator *is* the coordinator, collapsing γ onto δ removes a communication hop without removing a judgment gate. δ already holds external boundary authority; adding coordination authority does not compromise the independence of α↔β review. The operator-as-coordinator sees both α and β sides (γ's function) and executes release gates (δ's function) — but does not duplicate α's implementation work or β's review judgment.

- **α=β collapse (never safe).** Collapsing β onto α destroys review independence entirely. α reviewing α's own matter is order-0 observation masquerading as order-1. No cycle complexity or risk profile makes this collapse structural sound — the mechanism by which review adds information is eliminated.

- **ε=δ collapse (safe in small-protocol regimes).** When protocol-iteration volume is low, collapsing ε onto δ does not compromise any independence boundary. δ operates the cycle sequence; ε observes whether the protocol governing those cycles is coherent. The same actor can hold both functions when the ε workload (typically finding-triggered) does not justify separate specialization.

The pattern applies regardless of whether actors are humans, agents, or sessions. The structural property is the independence of observation, not the nature of the observer.

## §4a Persona, operator, protocol, project, receipt — the five-layer enforcement chain

A protocol that operates on real work needs more than role definitions. It needs a chain of artifacts that together determines (i) what kind of agent does the work, (ii) what the agent refuses to do, (iii) what counts as valid work in the domain, (iv) what concrete gates apply for this project, and (v) what evidence must exist before the work is transmissible to a parent scope.

The role ladder (§1) defines the abstract grammar (α/β/γ/δ/ε); the enforcement chain instantiates the grammar at a specific operating point.

| Layer | Question | Canonical home |
|---|---|---|
| **1. Persona** | What kind of mind is this agent? | `<hub>/spec/PERSONA.md` |
| **2. Operator contract** | What does this agent refuse to do? | `<hub>/spec/OPERATOR.md` |
| **3. Protocol overlay** | What counts as valid work in this domain? | `cnos.<protocol>/skills/<protocol>/` |
| **4. Project binding** | What concrete gates apply for this project? | `<project>/.<protocol>/` or `<project>/<protocol>/` |
| **5. Receipt / validator** | What evidence makes the work transmissible? | per-protocol receipt schema; V dereferences |

These five layers map onto a four-question grouping that names the boundaries:

- **Who** — persona + operator contract (layers 1–2)
- **What** — protocol overlay (layer 3)
- **Where** — project binding (layer 4)
- **How** — receipt + validator (layer 5)

A persona hub is not the protocol package, and the protocol package is not the project binding. Conflating any two layers produces drift: persona that smuggles project-specific data paths becomes unreusable across projects; protocol overlay that bakes in one project's gates produces a protocol that can only operate on that project; project binding that carries reusable role doctrine traps that doctrine inside the project.

### §4a.1 Discipline profile (required persona section)

The persona's `## Discipline` section is the standing artifact that encodes the agent's loss function. It is **not** equivalent to a task-prompt instruction ("be careful," "verify before claiming"). Task prompts are weak — they may or may not be loaded, may or may not be remembered. The discipline profile is a standing characteristic of the agent that survives prompt-to-prompt context churn and that the operator contract, protocol overlay, and receipt validator can each reference.

Required fields:

```markdown
## Discipline

Primary virtue: <what this agent optimises for>
Primary error: <the failure mode this agent must structurally resist>
Default tempo: <fast / slow / variable, and the constraint>
Claim/artifact boundary: <when this agent advances vs holds>
Refusal conditions: <conditions that force "deferred" / "blocked" / "rejected" instead of any output>
Receipt requirements: <what the agent's receipts must carry>
```

### §4a.2 Loss-function distinction

The deepest distinction across instantiations is not the role grammar — α/β/γ/δ/ε are the same shape everywhere. It is the loss function that the discipline profile encodes:

- **Engineering protocols (CDS):** optimise for *artifact improvement under repairable feedback*. Primary failure mode is the stalled loop — no working artifact ships, the system does not improve. Engineering has a strong immediate-feedback surface (compilers, tests, schemas, CI) that produces fast correction, so the discipline can support shipping bounded artifacts and repairing them.
- **Research protocols (CDR):** optimise for *truth-preserving claim transmission under uncertainty*. Primary failure mode is overclaim — a claim becomes stronger than its evidence; false knowledge propagates; the system believes something it did not measure. Research often lacks an immediate truth-check (a fabricated or overstated claim may look coherent for a long time), so the discipline must be more conservative at the claim boundary.

The same kernel runs both. The discipline profile is what makes one agent action-biased under repairable feedback and another epistemically conservative under irreparable claim transmission.

**Mode-switching one agent across discipline regimes is not safe.** Engineering "ship-and-repair" cadence leaks into research evidence discipline; research caution leaks into engineering cadence. The dispatch boundary — the operator routes work to the persona whose discipline matches the work — is the safe form of switching.

### §4a.3 Receipts enforce discipline mechanically

The receipt schema makes the discipline profile mechanically checkable. V (the validator, from the coherence-cell normal form) reads the receipt and verifies that the discipline-required evidence is present. This is the difference between "the agent says it was careful" and "the receipt cannot validate without the evidence the discipline requires."

A CDS (engineering) receipt requires, at minimum:

```
artifact_refs:    paths of changed files
test_refs:        test results
ci_refs:          CI results
diff_ref:         the cycle's diff
debt_refs:        explicit follow-up obligations
```

A CDR (research) receipt requires, at minimum:

```
claim_refs:       which claims this receipt asserts
data_refs:        dataset / mount / manifest / checksum
method_refs:      script paths + commit SHA
result_refs:      output file paths
claim_status:     observed | computed | inferred | hypothesized | indeterminate
limitations:      explicit caveats
reproduction:     β re-ran, command + output match
```

V's job is to reject receipts that fail the schema. Discipline becomes mechanical: "no `data_refs`" → V rejects research receipt; "no `test_refs`" → V rejects engineering receipt. The persona's `## Discipline` section names the loss function; the receipt schema makes it enforceable.

### §4a.4 Worked example — Sigma (engineering) and Rho (research)

Two persona hubs, one shared role grammar, two different discipline profiles:

- **`cn-sigma`** — engineering persona. Plays δ for CDS cycles. Discipline: action-biased; ship bounded artifacts; use tests/CI/review as the correction surface; record debt rather than blocking on perfection.
- **`cn-rho`** — research persona. Plays δ for CDR cycles. Discipline: epistemically conservative; no data → no empirical claim; every number cites its producing command; preserve negative and indeterminate findings; refuse rather than fill evidential gaps with plausible prose.

Both hubs load the same generic kernel (`cnos.cdd`'s coherence-cell normal form). They differ in which protocol overlay they load (`cnos.cds` vs `cnos.cdr`) and in their discipline profile. When a repo needs both kinds of work — Sigma ships a measurement tool, Rho runs it under a data gate and writes the field report — the dispatch boundary mediates: the operator routes each unit of work to the persona whose discipline matches it; the repo carries the handoff via commits, scripts, receipts, and reports.

## §4b Generic ε — the protocol-iteration role

§§1–4a declare α/β/γ/δ as the four roles that do *ordinary cell metabolism*: α produces matter, β reviews, γ coordinates the loop, δ operates the cycle sequence and effects the parent-scope boundary. §4b declares ε — the protocol-iteration role — separately, because ε's matter is different in kind from the matter the other four roles act on.

This section is the **generic ε doctrine**. Each c-d-X instantiation (cdd, cdr, cds, future cdw / cda / …) declares a CDD/CDR/CDS-specific ε surface that inherits the kernel grammar below by reference. Instantiations may extend the gap class taxonomy and bind the iteration artifact to a protocol-specific path; they do not re-declare the watched fields, the relationship to δ, or the MCA discipline.

### §4b.1 What ε observes

ε observes **the receipt streams emitted by closed coherence cells across many cycles**. Inside a single cell, α/β/γ/δ do ordinary metabolism and γ closes the cell with a parent-facing receipt. ε reads receipts from outside the cell, across cells, and asks whether the *protocol governing those cells* is itself coherent — whether it selects the right gaps, closes them durably, and produces a system that learns from its own cycles rather than repeating the same class of error.

ε's matter is **the protocol itself**. Where α produces code or prose or analysis, ε produces protocol patches: skill edits, schema refinements, validator additions, harness-contract clarifications. The patch lands wherever the gap surfaces; ε is not bound to a single package.

The structural property that distinguishes ε from δ: δ operates the *sequence* of cycles (selects which gap to close next, dispatches, releases) and is therefore inside the cycle-level frame. ε operates *across* sequences (reads receipts longitudinally, observes whether the protocol is improving) and is therefore outside any single cycle's frame. The order-of-observation distinction from §2 is what makes ε work add information δ could not have produced alone.

### §4b.2 Watched receipt fields

Every cell emits a parent-facing receipt (`schemas/cdd/receipt.cue#Receipt` is the generic kernel; domain extensions in `schemas/cdr/`, `schemas/cds/` unify against it). The receipt carries two fields ε reads on every cycle:

```cue
protocol_gap_count: int & >=0
protocol_gap_refs:  [...#ProtocolGapRef]
// consistency: protocol_gap_count == len(protocol_gap_refs)
```

`#ProtocolGapRef` is the typed structured reference to one protocol gap observed during the cycle (id, source — one of `receipt | artifact | issue | review | closeout` — and ref). The two fields together make "no protocol gaps this cycle" a positive signal carried by every receipt — not an absence inferred from missing artifacts. When `protocol_gap_count == 0`, the cycle ran cleanly with no protocol-iteration findings. When `protocol_gap_count > 0`, the receipt's `protocol_gap_refs` is the typed index into the iteration artifact (§4b.4) where the findings are dispositioned.

The consistency constraint (`count == len(refs)`) is structurally enforced by the schema; authors asserting drifted values fail `cue vet`. This makes the ε signal mechanically checkable at the receipt boundary — ε does not depend on prose-level audit to discover whether the cycle surfaced findings.

### §4b.3 Gap class taxonomy (instantiation pattern)

Each instantiation declares its own gap classes named to its loss function. The naming pattern is `{protocol}-{axis}-gap` where the axis names the protocol-internal surface the gap touches. The generic doctrine does **not** mandate a specific class count or axis set — it mandates that classes are named in the instantiation, typed against the instantiation's loss function, and that ε reads `protocol_gap_refs` typed against the instantiation's class set.

Reference instantiations:

- **cdd (engineering, repairable-feedback regime per §4a.2):** `cdd-skill-gap` (procedural skill underspecified or wrong), `cdd-protocol-gap` (CDD doctrine itself drifted), `cdd-tooling-gap` (tooling absent, wrong, or unavailable), `cdd-metric-gap` (measurement missing or wrong). See [`src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`](src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md).
- **cdr (research, truth-preserving regime per §4a.2):** `cdr-data-gate-gap`, `cdr-overclaim-gap`, `cdr-reproduction-gap`, `cdr-citation-gap`, `cdr-oracle-ambiguity-gap`, `cdr-construct-drift-gap`. The class set is wider than cdd because research-discipline failure modes are wider; the principle is the same. See [`src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md`](src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md).
- **cds, cdw, cda, future:** each instantiation declares its own class names against its discipline profile (§4a.2). The class set must be enumerable, the names stable across cycles, and each class oracle-defined (β can mechanically decide whether a finding fits the class).

A protocol that declares no gap classes has no machinery for protocol-iteration. The receipt's `protocol_gap_count` would have nowhere to type its `protocol_gap_refs` against; ε's matter (the protocol) would be invisible. The gap class declaration is therefore a Field-5-equivalent obligation: a protocol that instantiates the role ladder must declare its gap classes before live cycles run.

### §4b.4 Iteration artifact and cadence rule

For each instantiation, ε produces an **iteration artifact** when the cycle's receipt has `protocol_gap_count > 0`. The artifact dispositions each finding (per the per-finding shape declared by the instantiation) and is the durable record of what was found and what was done.

The artifact is **conditional**: required only when `protocol_gap_count > 0`. When the receipt records zero gaps, no iteration artifact is required; the receipt itself is the always-present record of *whether* an artifact is required, and its absence is unambiguous (no findings, no file). This resolves a class of survivorship-bias ambiguity that obtains when the iteration file's absence could mean either "clean cycle" or "skipped close-out": with the watched-field invariant, the receipt commits to one reading.

Each instantiation declares the artifact's name and path. References:

- **cdd:** `cdd-iteration.md` at `.cdd/unreleased/{N}/cdd-iteration.md` (in-cycle) → `.cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md` (after release); aggregator at `.cdd/iterations/INDEX.md`. The cadence rule is declared in [`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`](src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md). See [`cnos.cdd/skills/cdd/epsilon/SKILL.md`](src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md).
- **cdr:** project-binding-dependent path (typically `<project>/.cdr/waves/{wave-id}/cdr-iteration.md`); the protocol-overlay declares existence and shape, the project binding declares path. See [`cnos.cdr/skills/cdr/epsilon/SKILL.md`](src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md).

The generic doctrine does not mandate a fixed path. Single-binding instantiations (cdd in cnos itself) may bind to a fixed path; multi-binding instantiations (cdr over many possible project repos) leave the path to the project binding. Both are valid; the protocol-overlay layer declares the artifact's existence and per-finding shape.

### §4b.5 MCA discipline

For each finding ε dispositions one of three branches:

- **Ship now (MCA available)** — a skill patch, schema refinement, gate, or mechanization exists that closes the gap durably. Ship it as an immediate output in the same session ε surfaces the finding. This is the preferred branch when a patch is available.
- **Next-MCA (pattern real, no MCA yet)** — the gap is structural but the patch is non-trivial. File a protocol MCI (a tracked issue) and link the finding to its eventual MCA. The instantiation's iteration aggregator (cdd's `.cdd/iterations/INDEX.md`; cdr's project-binding equivalent) carries the row.
- **No-patch (one-off, no pattern)** — the finding is real but does not generalize; no protocol patch would prevent it without overreach. Disposition explicitly and note the drop. "No-patch" is not the default; it requires a justification (per CDD §13a for the cdd instantiation, and equivalent in other instantiations).

Silence — a finding without a disposition — is the failure mode this discipline is designed to prevent. A protocol where findings accumulate without dispositioning is a protocol that has lost the loop; that is the system-level signal ε's existence is meant to catch.

### §4b.6 ε's relationship to δ (collapse rule)

ε and δ are often the same actor in small-protocol regimes — one active operator running a handful of cycles per week or per cadence unit. Protocol-iteration volume does not justify dedicated specialization at that scale. The operator (δ) naturally accumulates the longitudinal view ε requires and performs ε work as part of cycle close-out triage.

Separation becomes warranted when protocol-iteration volume justifies dedicated attention: when the iteration artifact is written on most cycles with non-empty findings, when the iteration aggregator accumulates faster than one actor can triage, or when the operator's operational load (δ) crowds out the reflective work (ε). At that point, ε may be a distinct actor — a second agent or a dedicated human — who reads receipts across cycles and drives protocol patches independently.

The structural argument (per §4 hats-vs-actors): collapsing ε onto δ in small-protocol regimes does not destroy any independence boundary. δ operates the cycle sequence; ε observes whether the protocol governing those cycles is coherent. The same actor can hold both functions when the workload does not justify separate specialization. The collapse is one of the *safe* collapses (§4); contrast α=β (never safe).

ε is therefore a **structural role** in the scope-escalation ladder (§1 row 5), not a headcount requirement. The obligation is that ε work happens and is attributed, not that ε is performed by a distinct person.

### §4b.7 Instantiation declaration (what an instantiation must declare)

A protocol that instantiates the role ladder must declare, per §3 Field 5 + this section:

1. **The instantiation's gap class names** — enumerable, stable across cycles, oracle-defined (β can mechanically decide class membership). Examples: cdd's four classes, cdr's six classes.
2. **The location of receipt streams ε reads** — where closed-cell receipts land for the instantiation. cdd: `.cdd/releases/{X.Y.Z}/{N}/`. cdr: project-binding-dependent wave surface.
3. **The iteration artifact name and per-finding shape** — what ε writes when `protocol_gap_count > 0`. cdd: `cdd-iteration.md` per the per-finding shape in `cdd/post-release/SKILL.md §5.6b`. cdr: `cdr-iteration` artifact per `cnos.cdr/skills/cdr/epsilon/SKILL.md`.
4. **The iteration aggregator location** — the INDEX-like record that maintains the cross-cycle view. cdd: `.cdd/iterations/INDEX.md`. cdr: project-binding equivalent.
5. **The cadence rule** — when the iteration artifact is required. The generic rule is "when `protocol_gap_count > 0`"; instantiations may tighten (require more often) but should not loosen.

An instantiation that omits any of these five declarations has ε's role-letter named in its grammar but no operational machinery for protocol-iteration. The role is invoked but not grounded.

## §5 cdd as the reference instantiation

Coherence-Driven Development (cdd) is the first instantiation of this pattern
and the reference implementation. It applies the role ladder to software
development: α produces source code, tests, and docs; β reviews and
integrates; γ coordinates the cycle from issue creation to close-out; δ
operates the release sequence and external boundary; ε iterates the protocol
via `cdd-iteration.md`, the first-class artifact for cdd-self-improvement
findings.

The full cdd instantiation — lifecycle algorithm, instantiation-contract
fields, artifact specifications, role algorithms, and the six-field contract
as declared for cdd — is at
[`src/packages/cnos.cdd/skills/cdd/CDD.md`](src/packages/cnos.cdd/skills/cdd/CDD.md).

cdd does not claim to have invented this role structure. cdd made it explicit
for software development; this document generalises it to any
coherence-driven discipline. Readers wanting the full cdd lifecycle should
read `CDD.md` directly; readers wanting only the generic pattern stay here.

---

## §6 cdw as the planned sibling

cdw (coherence-driven writing) is the planned sibling protocol. It applies the
same role ladder to editorial prose work. cdw will be specified in a separate
issue and a future cycle; no `cdw/` directory exists in this cycle.

The four role mappings for cdw:
- **α** produces prose: drafts, articles, specifications, documentation.
- **β** reviews for clarity, coherence, and source-of-truth alignment.
- **γ** coordinates a writing cycle: scopes the issue, dispatches, and closes.
- **δ** operates a writing pipeline: selects which prose gaps to close in order.

ε iterates the writing protocol, assessing whether cdw is producing better
prose across cycles and patching the cdw skill bundle when findings accumulate.

The instantiation contract for cdw (all six fields per §3) will be declared
when `cdw/CDW.md` ships. This section is a forward pointer only. cdw is named
as a planned protocol in a separate issue / future cycle; readers should not
begin cdw work on the basis of this stub.

---

## §7 Naming convention for new c-d-X protocols

New protocols follow the pattern `c-d-{letter}` where the letter names the
matter type:

- **cdd** — development (source code, tests, docs)
- **cdw** — writing (prose, editorial, specification)
- **cdr** — research (investigation, synthesis, citation)
- **cda** — analysis (data, metrics, modelling)

The namespace is free: any tenant may claim a letter by filing a protocol
skill bundle issue. The issue establishes the letter, names the matter type,
and commits to declaring the six instantiation-contract fields before the
protocol is used in a live cycle. Conflict resolution follows issue-filing
order. No formal registry exists until a conflict arises.

Higher-order roles (ζ, η, and beyond the five defined here) are reserved but
undefined. The pattern as observed collapses at ε for most work. If reality
demands a sixth order, it may be named, but no c-d-X protocol has reached
that scale yet.

The letter indicates the matter, not the domain. A team writing software
documentation uses cdw for the prose work and cdd for the code work. The same
role actors may operate both protocols concurrently with distinct git
identities per `operator/SKILL.md`: `alpha@{project}.cdd.cnos` for cdd work,
`alpha@{project}.cdw.cnos` for cdw work.

---

## §8 Role-name stability

The role names α, β, γ, δ, ε are fixed across all c-d-X instantiations. The
verbs — produces, reviews, coordinates, operates, iterates — are fixed. Only
the matter and the oracles vary per protocol.

This stability is the property that makes the pattern portable. A person who
has operated as δ in a cdd context knows immediately what δ means in a cdw
context: select gaps, sequence cycles, dispatch. They do not re-derive the
role from scratch. The shared vocabulary collapses the onboarding cost for
every new protocol to the delta between "what α produces in cdd" and "what α
produces in cdw" — not the entire role structure.

The git identity convention follows the protocol namespace:
`{role}@{project}.{protocol}.cnos`, elided to `{role}@{protocol}.cnos` for
the cnos project (e.g. `alpha@cdd.cnos`, `beta@cdw.cnos`). The protocol
namespace in the identity disambiguates when the same agent operates across
protocols in the same week. Role names are not job titles and do not imply a
hierarchy of authority within the triad — γ coordinates the α↔β loop but does
not supervise α or β; δ operates the cycle sequence but does not approve β's
merge decisions.

The verb "operates" for δ and "iterates" for ε are precise choices. δ does
not merely schedule — δ selects which gaps are worth closing, sequences them
against coherence risk, and dispatches with full context. ε does not merely
improve — ε iterates the protocol under the same discipline the protocol
applies to its matter: observe the gap, close it with a verifiable artifact,
measure the delta.

---

## §9 Glossary

**role** — A structural position in the scope-escalation ladder, identified by
a Greek letter (α through ε). A role is defined by what it observes (its
domain), what it produces (its output), and its relationship to the roles
above and below it. Roles are not job titles; the same person or agent may
hold different roles in different protocols. Within a single cycle, α and β
must be held by distinct actors — the independence is the mechanism by which
review adds information.

**matter** — The deliverable that α produces in a given protocol. The matter
type is protocol-specific and concrete: source code and tests in cdd, editorial
prose in cdw, a research summary in cdr. The matter type determines what β
reviews, what the review oracle tests against, and what ε assesses whether the
protocol is improving at producing over time.

**frame** — The context within which a role operates and which it treats as
given. Each role's frame is the previous role's content. β's frame is α's
matter — β reviews what α produced against the declared oracle. γ's frame is
the α↔β loop — γ observes whether the loop is closing coherently, not
re-adjudicating β's verdict. A role cannot step outside its own frame — that
is precisely what the next role provides.

**instantiation** — A protocol that applies the generic role-scope ladder to a
specific matter type by declaring the six fields of the instantiation contract
(§3): matter type, review oracle, γ close-out artifact, δ cadence, ε
iteration cadence, and actor collapse rule. cdd is the reference instantiation.
cdw is the planned sibling. An instantiation without all six fields declared
is incomplete — the roles are named but the protocol is not operational.

**scope-escalation** — The structural property by which each role operates at
a strictly higher order of observation than the role below it. Scope escalates
because each role's domain is the previous role's process, not just its
output. β observes the output; γ observes the process that produced it; δ
observes the sequence of processes; ε observes whether the protocol governing
those sequences is coherent.

**order-of-observation** — The level at which a role observes the system,
numbered 0 through 4. Order 0 is the matter itself. Order 1 is β's review.
Order 2 is γ's coordination of the α↔β loop. Order 3 is δ's operation of the
γ-loop sequence. Order 4 is ε's iteration of the δ-discipline. Each order
requires being outside the lower order's frame — which is why α and β must be
separate within any single cycle, and why ε cannot be the same actor as α in
the cycle whose protocol ε is iterating.

**oracle** — The mechanically checkable criterion by which β determines that
α's matter closes the declared gap. An oracle names the verifiable check and
the tooling or procedure used to run it. In cdd, oracles are
acceptance-criteria checks expressed as `rg`, `wc`, `head`, or test runner
commands. A protocol without a stated oracle cannot distinguish "approved"
from "no longer looking."

**cycle** — One instance of the α↔β loop, coordinated by γ, resulting in a
closed gap. A cycle begins when δ selects a gap and dispatches; it closes when
γ writes the close-out triage and the next-move commitment is recorded. The
cycle is the unit of coherence measurement: each cycle should leave the system
more coherent than it found it.
