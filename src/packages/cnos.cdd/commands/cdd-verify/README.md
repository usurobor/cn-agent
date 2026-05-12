# CDD Artifact Checker

`cn-cdd-verify` validates CDD cycle artifact completeness against the canonical contract declared in `src/packages/cnos.cdd/skills/cdd/CDD.md` §5.3a (Artifact Location Matrix).

## Quick Start

```bash
# Check the entire repository CDD artifact ledger
src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --all

# Check only unreleased/active cycles  
src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --unreleased

# Check specific release version
src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --version 3.74.0
```

## Modes

- **Repository-wide validation** (`--all`, `--unreleased`): Scans directories for CDD cycles and validates artifact completeness
- **Single version validation** (`--version`): Checks canonical release artifacts for one version  
- **Cycle-specific validation** (`--version X --cycle N`): Validates current artifact layout for one cycle
- **Legacy validation** (`--version X --triadic`): Also checks pre-#283 close-out paths

## Cycle Classification

The checker automatically classifies cycles based on artifact presence:

- **Triadic cycles**: Have `beta-review.md` → requires all 5 artifacts (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`)
- **Small-change cycles**: No `beta-review.md` → optional artifacts per CDD.md §1.2 collapse rules
- **Unknown cycles**: No recognizable artifacts → warning

## Exception Handling

Historical gaps that cannot be repaired can be documented in a YAML exceptions file:

```bash  
cn-cdd-verify --all --exceptions .cdd/exceptions.yml
```

Exception format:
```yaml
exceptions:
  - path: ".cdd/unreleased/286/alpha-closeout.md"
    missing_artifacts: ["alpha-closeout.md"]
    reason: "historical cycle predating close-out requirement"
    repair_possible: false
    follow_up: "none"
```

## Testing

Run the test fixtures to verify behavior:

```bash
src/packages/cnos.cdd/commands/cdd-verify/test-fixtures.sh
```

## CI Integration

This checker runs automatically in CI as job `I6: CDD artifact ledger validation` on every push and pull request.

## Exit Codes

- `0`: All required checks pass (warnings allowed)
- `1`: One or more required artifacts missing (unless exception-backed)