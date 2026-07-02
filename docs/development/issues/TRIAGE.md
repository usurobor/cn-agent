# Issue triage

Operational rules for turning a raw issue into a coherently-labeled,
dispatchable (or deferred, or closed) item. The label *definitions* live in
[TAXONOMY.md](TAXONOMY.md); this document is *how to apply them*.

## Classifying a new issue

1. **Pick the primary `kind/*`.** Ask "what kind of work is this?" — see the
   decision guide below. Exactly one primary kind.
2. **Add `area/*`.** Where does the work live? One or more.
3. **Add a priority** (`P0`–`P3`, or `priority/deferred`).
4. **Add a `status:*`** only if the issue is actionable now. Design/tracking/
   research issues that aren't ready to execute carry no `status:*`.
5. **Add `dispatch:cell` + `protocol:*`** only if the issue is a genuine,
   proof-carrying cell the dispatch system can execute (see below).

## Distinguishing the kinds

The most common triage mistake is blurring these. Use the verb:

| If the issue… | …it is |
|---|---|
| **fixes** behavior that is broken vs. current intended behavior | `kind/bugfix` |
| **removes** stale structure / terminology / dead code / drift / noise (no new capability) | `kind/cleanup` |
| improves **how** the work engine (CDD/CDS/cells/waves/δ/β/γ) operates | `kind/process` |
| adds a **new** user- or product-facing capability | `kind/feature` |
| builds command / CI / validator / generator / dev machinery | `kind/tooling` |
| **changes the written contract** (spec/design/theory/doc) of the system | `kind/doctrine` |
| inventories / drift-checks / reviews a baseline | `kind/audit` |
| is an **umbrella** coordinating other issues | `kind/tracking` |
| is an **open question** / exploration with no committed answer | `kind/research` |
| **is a skill artifact** as its main deliverable | `kind/skill` |
| is a **timeboxed prototype**, not final architecture | `kind/spike` |

Tie-breakers:

- **bugfix vs cleanup** — bugfix restores *behavior*; cleanup removes *structure/
  noise* without changing behavior. A stale doc reference is cleanup; a validator
  that passes when it should fail is a bugfix.
- **process vs feature** — process improves the engine that builds things; feature
  is a thing the engine builds. "Require review-request proof before
  `status:review`" is process; "web console" is feature.
- **doctrine vs process** — doctrine changes what the system *should be understood
  to be* (a written contract); process changes *how the work happens*. They often
  co-occur; when they do, tag the primary and, if truly necessary, a secondary.
- **cleanup vs doctrine** — deleting a retired term is cleanup; rewriting the
  contract that defined it is doctrine.

A single secondary `kind/*` is acceptable when an issue genuinely spans two (e.g.
a `kind/process` change that is also a `kind/bugfix`), but prefer one primary.

## When to add `dispatch:cell` (and when not)

Add `dispatch:cell` + `protocol:cds` (or `protocol:cdd`) **only** when the issue
is an executable cell: scoped, with acceptance criteria and a proof oracle, that a
dispatch wake can claim and run end-to-end. A dispatch wake claims on
`dispatch:cell + protocol:{p} + status:todo` (minus the excluded statuses), so:

- `status:ready` — shaped and approved, but held. The operator flips
  `status:ready → status:todo` to release it. **Not** claimed while `ready`.
- `status:todo` — claimable now.

**Do NOT** add `dispatch:cell` to:

- a bare design, tracking, research, or doctrine issue with no proof oracle;
- an umbrella/roadmap issue (label it `kind/tracking`, no dispatch);
- anything whose destination or ownership is unclear.

## When to close, supersede, or defer

- **Close as `resolution/completed`** when the work shipped (verify against the
  repo — e.g. the file exists, the workflow is gone, the test passes).
- **Close as `resolution/superseded`** when a *different* issue or a later design
  delivered the capability. Name the superseding issue.
- **Close as `resolution/duplicate`** (with `duplicate_of`) after confirming the
  issues are the same; keep the richer/canonical one open.
- **Close as `resolution/wontfix`** when the work is intentionally declined.
- **`priority/deferred`** (kept open) when the work is real but parked.

Always **verify before closing** — read the issue and check the repo state.
Closing is reversible, but a wrong close on live work is friction.

## When to ask the operator

Bring a decision back rather than deciding it yourself when it is **STOP-class**:

- behavior / authority / schema-vocabulary / label-or-dispatch-semantics change;
- a non-prose (semantic) golden diff, or a release/deploy boundary;
- deleting current operational data;
- an unclear destination or ownership, or a conflict between two approved invariants.

## When δ may relabel under MCA

A parent/wave δ operating inside an **approved** contract may relabel issues as a
minor coherent action (MCA) — and must record it — without a per-item interrupt.
In scope for MCA relabeling:

- normalizing an issue to `one kind + area(s) + priority` per this taxonomy;
- repointing a stale `area/*` or retiring an unprefixed twin;
- adding `resolution/*` and closing an obviously-superseded historical experiment
  **after verifying** it against the repo;
- consolidating a confirmed duplicate.

Out of scope for MCA (STOP — ask the operator): changing an issue's **dispatch
eligibility** (`dispatch:cell` / `protocol:*` / `status:todo`) in a way that would
release it for execution, closing anything whose completeness is uncertain, or
re-scoping an issue's meaning. This mirrors the wave-execution autonomy policy —
see #539 (parent/wave δ MCA) and the CAP skill (`agent/cap/SKILL.md`).

## Board hygiene

Target state: every open issue carries exactly one primary `kind/*`, at least one
`area/*`, and a priority; actionable issues carry a `status:*`; only genuine cells
carry `dispatch:cell` + `protocol:*`; closed issues that weren't a plain
completion carry a `resolution/*`. Umbrella/tracking issues may omit priority.
