# Issue label taxonomy

How open issues in this repo are labeled so the board can be reasoned about at a
glance and dispatched safely. Every open issue should carry a coherent label set:

- **exactly one** primary `kind/*`
- **one or more** `area/*`
- **one** priority (`P0`–`P3`, or `priority/deferred`)
- **one** `status:*` when the issue is actionable
- `protocol:*` **only** when the issue is dispatchable through a protocol
- `dispatch:cell` **only** when the issue is actually executable by CDS/CDD
- a `resolution/*` on close, when the close reason is not a plain completion

The single most important split is the primary `kind/*` — it answers "what *kind*
of work is this?" and keeps feature work, cleanup, doctrine, and process
improvements from blurring together.

## Primary kind (choose exactly one)

| Label | Meaning |
|---|---|
| `kind/bugfix` | Something is broken relative to current intended behavior. |
| `kind/cleanup` | No new capability; removes stale structure, terminology, dead code, drift, or repo noise. |
| `kind/process` | Improves how CDD/CDS/cells/waves/δ/β/γ operate — the work engine itself. |
| `kind/feature` | New user- or product-facing capability (web console, Telegram, voice, package source, …). |
| `kind/tooling` | CLI, CI, validators, checkers, generators, developer/command machinery. |
| `kind/doctrine` | Spec/design/theory/doc contract that changes how the system should be understood. |
| `kind/audit` | Inventory, drift check, historical gap review, baseline assessment. |
| `kind/tracking` | Umbrella/tracker issue; does not itself implement unless separately scoped. |
| `kind/research` | Open question / design exploration / unresolved theory. |
| `kind/skill` | A skill artifact itself is the main deliverable. |
| `kind/spike` | Timeboxed exploration/prototype; not intended as final architecture. |
| `kind/chore` | Routine maintenance (kept from the prior scheme; prefer a more specific kind when one fits). |

### The conceptual split (memorize this)

- `kind/process` = improves **how** cnos builds / works / thinks.
- `kind/feature` = adds a **new** product/user capability.
- `kind/cleanup` = **removes** drift / noise / stale surface.
- `kind/bugfix` = **fixes** broken intended behavior.
- `kind/doctrine` = **changes the written contract** of the system.
- `kind/tooling` = command / CI / dev machinery.

`kind/audit`, `kind/research`, `kind/spike`, `kind/tracking`, and `kind/skill` are
for inventory, open questions, timeboxed prototypes, umbrellas, and skill-artifact
deliverables respectively.

## Area (choose one or more)

Where the work lives. Non-exhaustive; add new `area/*` labels as the system grows.

`area/agent` · `area/cdd` · `area/cds` · `area/core` · `area/runtime` ·
`area/protocol` · `area/transport` · `area/ci` · `area/cli` · `area/docs` ·
`area/package` · `area/security` · `area/ui` · `area/telegram` · `area/skills` ·
`area/wake` · `area/coherence` · `area/identity` · `area/handoff` ·
`area/diagnostics` · `area/observability` · `area/ux` · `area/memory` ·
`area/continuity` · `area/distribution`

### Transition from unprefixed labels

Older unprefixed labels — `core`, `cdd`, `cds`, `ci`, `cli`, `package`, `handoff`,
`process`, `docs` — are being retired in favor of `area/*` (and, for `process`, the
`kind/*` scheme). **Do not block on this.** Both may coexist during the transition;
add the prefixed `area/*` when relabeling, and drop the unprefixed twin when it is
convenient rather than as a gating step.

## Priority

`P0` · `P1` · `P2` · `P3`, plus `priority/deferred` for work intentionally parked.
Priority is orthogonal to `kind` and `area`.

## Status (when actionable)

`status:ready` · `status:todo` · `status:in-progress` · `status:review` ·
`status:changes`.

Only `status:todo` is claimable by a dispatch wake (per the standard dispatch
selector `dispatch:cell + protocol:{p} + status:todo`, minus the excluded
statuses). `status:ready` means "shaped and approved, but held" — the operator
flips `status:ready → status:todo` to release it for dispatch. A non-dispatch
design/tracking issue carries no `status:*`.

## Dispatch and protocol

- `dispatch:cell` — the issue is a cell that a dispatch wake can execute. Apply
  **only** when the issue is genuinely executable by CDS/CDD (a scoped, proof-carrying
  cell), never to a bare design/tracking/research issue.
- `protocol:cds` / `protocol:cdd` — the concrete protocol the cell dispatches through.
  Apply **only** together with `dispatch:cell` on a dispatchable issue.

## Resolution (on close)

Apply on close when the reason is not a plain "done":

`resolution/completed` · `resolution/superseded` · `resolution/duplicate` ·
`resolution/deferred` · `resolution/wontfix`.

These are especially useful for retiring historical experiment issues (early
wake/OAuth test runs) and for consolidating duplicates.

## Applying the taxonomy — rules of thumb

1. Every open issue gets exactly one primary `kind/*` and at least one `area/*`.
2. Add a priority; add a `status:*` only if the issue is actionable now.
3. `dispatch:cell` + `protocol:*` go on together, and only on genuinely dispatchable cells.
4. Non-dispatch design/tracking/research issues do **not** carry `dispatch:cell`.
5. Close stale historical experiment issues with `resolution/superseded` (or
   `resolution/completed` when the work actually shipped elsewhere).
6. Consolidate duplicates: keep the canonical issue, close the rest with
   `resolution/duplicate`.
7. Retag legacy `enhancement`-only issues into the right `kind/*` (usually
   `kind/feature`, `kind/process`, `kind/doctrine`, or `kind/cleanup`).

## Label lifecycle

Labels are created lazily — applying a `kind/*` or `area/*` label that does not yet
exist creates it. Keep the set small and prefixed; prefer an existing label to a
near-synonym. Retire ambiguous unprefixed labels opportunistically, not as a blocking
migration.
