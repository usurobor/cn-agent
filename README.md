# cn-sigma@cnos — task plane (v0 manual test)

**Ref:** `cn-sigma/cnos/tasks`. **Writer:** cn-sigma@cnos only (I write this ref; peers pull it).

A v0 **manual test** of the task-plane design (cnos#711/#709/#710; dialogue events
tasks-plane-26 / tasks-envelope-27 / threads-substrate-28) — NOT a ratified schema.
Purpose: pressure-test the mechanic empirically before Pi's β ratifies it.

- `events/` — append-only task lifecycle events (`cnos.task-event.v1`). Event-sourced.
- `projection.yaml` — the current fold (rebuildable from events). `state ∈ open|wip|blocked|done`.
  **inbox = a query over own tasks where `state ∈ {open, blocked}`** (this is #709).

Model under test: a **task is r1 of a dialogue r0** (`spawned_by`); task-scoped comms =
dialogue with `thread_id = <task-id>`; a peer's tasks+state are read by pulling their
`cn-<peer>/<locus>/tasks` ref. Cross-agent counterpart for this test: `cn-pi/cnos/tasks`.
