# Run Record

**Run ID:** A1-cold-001
**Date:** 2026-03-25
**Evaluator:** Claude
**Kata ID:** A1-open-op-registry-conflict
**Family:** Local implementation + invariant proof
**Model:** GPT-5.4 Pro
**Temperature / settings:** default
**Fresh session:** yes

## Arm

- cold

## Skills Loaded

- none beyond baseline doctrine/runtime context

## Governing Skills Hypothesis

- none (cold baseline)

## Cost

- prompt tokens: 2,100
- completion tokens: 980
- latency: 16s
- rounds: 1
- extra artifacts created: none

## Outcome

- partial

## Scores

- α: 2
- β: 2
- γ: 2
- efficiency: 4

## Failure Classification

- selection-limited

## Notes

- Agent correctly identified that duplicate op kinds are bad.
- Agent proposed "last one wins with warning" as one fallback, which violates the invariant.
- Agent suggested example tests for duplicate-name and duplicate-op cases.
- Agent did not propose property/order-independence testing.
- Agent did not clearly separate:
  - registry construction
  - dispatch/runtime execution
- Failure shape: local patching mindset, weak invariant discipline.
