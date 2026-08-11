# architecture/

**How the parts fit together — the system's structure and its invariants.**

## System

- [ARCHITECTURE.md](ARCHITECTURE.md) — system overview: how the parts relate.
- [CAA.md](../reference/runtime/CAA.md) — coherent agent architecture.
- [INVARIANTS.md](INVARIANTS.md) — architectural invariants.

## Cognition

- [COGNITIVE-SUBSTRATE.md](cognitive-substrate/COGNITIVE-SUBSTRATE.md) — cognitive asset classes (doctrine, mindsets, skills).
- [CAR.md](cognitive-substrate/CAR.md) — cognitive asset resolver: local, versioned cognition.

## Runtime

- [CELL-RUNTIME.md](CELL-RUNTIME.md) — cell classes (WC/PC/CC), matter domains, and the generic cell runner. *Proposed* (#627 / #628); realization peer of `COHERENCE-CELL-NORMAL-FORM.md`.
- [AGENT-DIALOGUE-PROTOCOL.md](AGENT-DIALOGUE-PROTOCOL.md) — agent-to-agent dialogue v0: writer-owned `{agent, locus}` activation streams, thread reconstruction, and the boundaries separating dialogue from #690 memory and from project authority. Canonical design of record, ratified 2026-08-05 (#698).

## Security & observability

- [SECURITY-MODEL.md](security/SECURITY-MODEL.md) — sandbox, FSM enforcement, audit trail.
- [TRACEABILITY.md](security/TRACEABILITY.md) — event stream, state projections, readiness.

## Constraints

- [DESIGN-CONSTRAINTS.md](DESIGN-CONSTRAINTS.md) — system-wide design constraints.
- [HUB-PLACEMENT-MODELS.md](HUB-PLACEMENT-MODELS.md) — hub placement topology models.
