# quickstart/

**Runnable first experiences — act before reading the whole doctrine.**

> **Pass 1 overlay.** This area is being populated. Its flagship artifact —
> a Hello World cell that emits a real receipt and TSC coherence report
> (Demo 0) — is pending the two-wakes migration (`cnos#467`). Until then,
> the closest runnable entry points live in the agent runtime and operator
> guide.

## Available now

- [Activate an agent in a repo](../../src/packages/cnos.core/skills/agent/activate/SKILL.md) — the `activate` skill.
- [Operator handshake](../guides/HANDSHAKE.md) — connect an operator to a repo and exchange the first coordination.
- [Automation](../guides/AUTOMATION.md) — wake/dispatch automation surfaces.

## Coming (gated on Demo 0)

- `HELLO-WORLD-WITH-RECEIPTS.md` — the smallest end-to-end cell: issue → cell → validation → receipt → TSC report → degraded release, captured live (not illustrative).
- `RUN-A-CELL.md` — drive one CDD cell start to finish.

The captured receipt and coherence report from Demo 0 will land in
[`evidence/`](../evidence/README.md); the conceptual argument for why the
receipt matters is [Dumb Models, Smart Cells](../papers/DUMB-MODELS-SMART-CELLS.md).
