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

## WCC 0.1 — the coding cell

- [WCC-0.1-RUNBOOK.md](WCC-0.1-RUNBOOK.md) — how to run one episode, how to extract the patch, and how to check it outside the closure that reports it.
- [WCC-ADVERSARIAL-SUITE.md](WCC-ADVERSARIAL-SUITE.md) — seven cases that try to break the cell, and what each one establishes.
- [WCC-CODE-QUALITY-AUDIT.md](WCC-CODE-QUALITY-AUDIT.md) — eight measured findings against `eng/go`, `eng/write-functional` and `eng/evolve`.

## Security & observability

- [SECURITY-MODEL.md](security/SECURITY-MODEL.md) — sandbox, FSM enforcement, audit trail.
- [TRACEABILITY.md](security/TRACEABILITY.md) — event stream, state projections, readiness.

## Constraints

- [DESIGN-CONSTRAINTS.md](DESIGN-CONSTRAINTS.md) — system-wide design constraints.
- [HUB-PLACEMENT-MODELS.md](HUB-PLACEMENT-MODELS.md) — hub placement topology models.
