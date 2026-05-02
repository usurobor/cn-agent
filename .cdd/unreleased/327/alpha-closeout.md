# α Close-out — #327

## Summary

Implementation complete. `scripts/validate-release-gate.sh` added and wired into `release.sh` before stamp/move/tag. All 3 ACs met in first pass. Test suite covers triadic (5 required files), small-change (1 required file), and no-unreleased-dir cases.

## Artifacts

- `scripts/validate-release-gate.sh` — new validation script
- `scripts/release.sh` — wired gate call at step 4
- `tests/test-validate-release-gate.sh` — test suite

## Findings received

0 findings from β R1. No fix-rounds needed.

## Debt

None introduced.
