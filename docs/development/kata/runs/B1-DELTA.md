# Delta Interpretation

**Kata ID:** B1-runtime-contract-v2-parity-review
**Cold Run:** B1-cold-001
**Selected Run:** B1-selected-001

## What improved

### α: 3 → 4

The selected-skills run produced a cleaner review structure and better evidence discipline.

### β: 1 → 4

This is the real gain. The selected-skills run caught the authoritative cross-surface mismatches:

- runtime code vs AGENT-RUNTIME example
- runtime code vs CAA wake strata
- issue promise vs doctor depth

### γ: 2 → 3

The selected-skills run proposed a much better finishing patch set and cleaner issue-closure decision.

### Efficiency: 4 → 3

The run became more expensive in prompt tokens and latency, but the gain was worth it.

## Failure-shape change

**Cold:**
- local implementation review
- weak sibling/authority-surface reasoning

**Selected:**
- issue-contract-first review
- strong cross-surface parity reasoning
- better minimal-patch judgment

## Provisional conclusion

For this kata family, review + documenting appears to be a real and useful governing-skill pair. The next step is transfer testing on a different same-family kata.

## How to use this immediately

For the next run, do:

- Task A cold
- Task A selected
- Task B selected

where Task B is another review/parity kata of the same shape:

- markdown vs JSON mismatch
- canonical doc vs executable skill mismatch
- issue ACs vs implementation-plan mismatch

If the selected-skills arm catches the same class of contradictions again, you have evidence of transfer, not just repair.
