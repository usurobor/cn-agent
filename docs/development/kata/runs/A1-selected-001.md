# Run Record

**Run ID:** A1-selected-001
**Date:** 2026-03-25
**Evaluator:** Claude
**Kata ID:** A1-open-op-registry-conflict
**Family:** Local implementation + invariant proof
**Model:** GPT-5.4 Pro
**Temperature / settings:** default
**Fresh session:** yes

## Arm

- selected-skills

## Skills Loaded

- coding
- testing

## Governing Skills Hypothesis

1. coding — should force the invariant into the data/model and reject invalid states explicitly
2. testing — should force stronger proof than a few examples, especially order-independence and negative-space coverage

## Cost

- prompt tokens: 5,200
- completion tokens: 1,430
- latency: 28s
- rounds: 1
- extra artifacts created: none

## Outcome

- pass

## Scores

- α: 4
- β: 3
- γ: 3
- efficiency: 3

## Failure Classification

- skill-effective

## Notes

- Agent explicitly named the invariant:
  - effective op kinds must resolve uniquely
  - duplicate extension names must reject
  - duplicate op kinds must reject
- Agent proposed changing registry construction so conflicts accumulate into a validation error, not shadowing behavior.
- Agent proposed table-driven tests for:
  - unique manifests succeed
  - duplicate name rejects
  - duplicate op kind rejects
  - registry result does not depend on manifest order
- Agent suggested property-style testing for order-independence, which is the right stronger proof form here.
- Failure shape changed from "warn and keep going" to "invalid state must be rejected."
