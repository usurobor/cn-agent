# β close-out — cnos#593

## Review summary

Converge, zero blocking findings, two noted-but-non-blocking observations (see beta-review.md). Every AC independently re-verified by re-running the actual test/gate commands, not by re-reading α's prose. The hub-detection bug fix was independently reproduced in a fresh sandbox before I accepted it as real rather than an artifact of α's own environment.

## R0 verdict — reconfirmed, still holds

Converge. All 11 ACs pass under independent re-verification. Implementation-contract conformance (Go-only, extends existing `issues-fsm`/`cell` packages, no new binary, injectable-dependency testing pattern mirrored from `cell.Finalizer`, zero diff to the FSM table/engine, protocol-agnostic renderer wiring) confirmed by direct diff inspection.

## The hub-gate fix — independent verification

This is the one change in this cycle that touches shared kernel-command infrastructure (`CommandSpec.NeedsHub`) rather than being scoped purely to the new scanner. I scrutinized it specifically for blast radius:

- Confirmed via `git diff` that only `CellFinalizeCmd`'s `Spec()`/`Run()` changed — `CellResumeCmd` and `CellReturnCmd` (the other two commands in the same file) are byte-for-byte unchanged.
- Confirmed `CellReturnCmd` already sets `NeedsHub: false` with the same justification ("operates via gh CLI; no hub required") — this fix brings `CellFinalizeCmd` into line with an existing precedent in the same file, not a novel pattern.
- Reproduced the bug myself independently (separate sandbox, separate `mktemp -d`, not α's) to confirm it is real and not an artifact of α's particular environment.
- Confirmed the fallback (`git rev-parse --show-toplevel`) only activates when `inv.HubPath` is already empty — a real downstream hub (with `.cn/` present) gets `inv.HubPath` from the existing `discoverHub()` walk exactly as before, unaffected.

## Process observations for δ

None requiring escalation. The scope guardrails (no new labels, no daemon, no claim/resume-guard change, no FSM weakening) all held under my own read of the diff, not just α's self-report.

## Release notes

`cn issues fsm scan --protocol P [--apply]`: new mechanical recovery scanner for stranded CDS dispatch cells (cnos#593). `cn cell finalize` no longer requires a `.cn/` hub directory (cnos#593 bug fix, unblocks the scanner and the pre-existing cnos#591 finalizer step in this repository).
