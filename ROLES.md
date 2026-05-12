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

## §4 cdd as the reference instantiation

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

## §5 cdw as the planned sibling

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

## §6 Naming convention for new c-d-X protocols

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

## §7 Role-name stability

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

## §8 Glossary

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
