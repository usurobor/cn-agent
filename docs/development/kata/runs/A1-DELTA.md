# Delta Interpretation

**Kata ID:** A1-open-op-registry-conflict
**Cold Run:** A1-cold-001
**Selected Run:** A1-selected-001

## What improved

### α: 2 → 4

The selected-skills run moved from vague conflict handling to explicit invariant-preserving rejection.

### β: 2 → 3

The selected-skills run better respected the system boundary:

- registry validity first
- runtime consequences second

### γ: 2 → 3

The selected-skills run proposed a smaller, more coherent patch and stronger regression strategy.

### Efficiency: 4 → 3

The gain cost more tokens/latency, but it was still economical.

## Failure-shape change

**Cold:**
- local patching
- warning/overwrite thinking
- weak proof depth

**Selected:**
- invariant-first
- invalid-state rejection
- stronger proof family selection

## Provisional conclusion

For this kata family, coding + testing looks like the right governing pair. The important signal is not just score increase; it is the shift from "keep going with warning" to "invalid registry state must be rejected."
