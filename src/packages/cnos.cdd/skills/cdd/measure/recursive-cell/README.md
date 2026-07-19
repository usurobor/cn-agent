# cnos recursive-cell CM

Candidate coherence methodology produced from the external CC analysis of
cnos#662. `SKILL.md` is the authority. `INSTRUCTION.md` is the one executable
semantic instrument: deterministic assembly binds the exact pinned TSC
Self-Measure v3.2.4 bytes to the cnos supplement.

Publication issue: cnos#669. Parent runtime wave: cnos#627.

## Frozen-instrument check

From the cnos source root:

```bash
CM=src/packages/cnos.cdd/skills/cdd/measure/recursive-cell
"$CM/instruction/assemble-instruction.sh" --check
sha256sum "$CM/INSTRUCTION.md"
# de0414bac1a086d6ed3ac5c7d84cd6293517eb83c920f6c4a0723f5dfac06a3a
```

`assemble-instruction.sh --output <path>` is the exact checked assembly
command. It refuses changed core, supplement, or composite bytes.

## Exact #662 R7 State-A route

The TSC CLI resolves manifest paths from `--root`. Materialize the source CM
inside the frozen target checkout. The runner then emits all six registered
prompts plus the separate H01-H13 assessment prompt in one clean step:

```bash
CNOS_SOURCE=/path/to/cnos-cycle-669
TSC_SOURCE=/path/to/tsc-at-26aab5023f03dc7d0abf82e5fdba20134fc6adad
TARGET=/tmp/cnos-662-r7
RUN_ROOT=/tmp/cnos-662-r7-run
RESPONSES=/tmp/cnos-662-r7-responses
INVARIANTS=/tmp/cnos-662-r7-invariant-assessment.json
CM_REL=src/packages/cnos.cdd/skills/cdd/measure/recursive-cell
COH="$TSC_SOURCE/engine/ocaml/_build/default/bin/main.exe"
CM_REVISION="$(git -C "$CNOS_SOURCE" rev-parse HEAD)"
TSC_REVISION=26aab5023f03dc7d0abf82e5fdba20134fc6adad

git -C "$CNOS_SOURCE" worktree add --detach "$TARGET" \
  a0d39293a27cfe57b49dacff696345b1ee2cdb40
mkdir -p "$TARGET/$(dirname "$CM_REL")"
cp -a "$CNOS_SOURCE/$CM_REL" "$TARGET/$(dirname "$CM_REL")/"
RUNNER="$TARGET/$CM_REL/runner/recursive-cell-runner.py"

"$RUNNER" emit \
  --coh "$COH" --root "$TARGET" --output "$RUN_ROOT" \
  --engine-revision "$TSC_REVISION" --cm-revision "$CM_REVISION" \
  --timestamp 2026-07-19T00:00:00Z
```

Produce one exact standard TSC v3.2.4 response for each emitted target prompt,
named `cc662-system.json` and `cc662-l0.json` through `cc662-l4.json` under
`$RESPONSES`. Produce `$INVARIANTS` from
`$RUN_ROOT/invariant-assessment-prompt.md`; it is a separate typed document,
not an extension to any TSC witness. Then ingest all seven inputs together:

```bash
"$RUNNER" ingest \
  --coh "$COH" --root "$TARGET" --output "$RUN_ROOT" \
  --engine-revision "$TSC_REVISION" --cm-revision "$CM_REVISION" \
  --timestamp 2026-07-19T00:00:00Z \
  --responses "$RESPONSES" --invariants "$INVARIANTS"
```

The output root contains the six prompts, `prompt-digests.json`, the generated
invariant-assessment prompt, six canonical reports under `reports/`, and one
schema-validated `recursive-cell-run.json`. The latter binds the actual coh
executable, active runner, schema, template, prompts, responses, reports,
assessment, and target bytes. `declared_cm_revision` and
`declared_engine_revision` are explicitly caller assertions; byte identity is
carried by individual SHA-256 fields and `cm_authority_bundle_sha256`. That
bundle hashes, in sorted CM-relative path order, each path, NUL, byte length,
NUL, file bytes, NUL across the instruction, runner, schema, assessment
template, registry, preflight, and six manifests. The result also computes the
unweighted L0-L4 geometric aggregate, applies H01-H13, and emits exactly one
deterministic bottleneck and disposition.

`scripts/ci/test-recursive-cell-runner.sh` covers this orchestration, schema,
math, gating, provenance, and refusal behavior with deterministic fixture
responses and a strict fake CLI. It does not claim real-engine integration or
semantic calibration. The earlier bundled self target separately preserves one
real pinned-engine external-response smoke ingestion and three mechanical runs.

## Path-base boundary

Authority paths are repository-root-relative; `--root` is the explicit target
root. Installed-package activation is deliberately refused until TSC supports
separate authority and target bases:

```bash
"$CM/resolve-authority.sh" --source-root "$PWD"
./scripts/ci/smoke-recursive-cell-package.sh
```

Arbitrary State-B methodology and target loading remains unshipped. Do not use
#662 as a held-out target: it designed this CM and confers no standing. The
historical N=1 semantic sample was produced with the prior instrument and is
superseded after R3; it cannot count toward any later k=3 consistency claim.
