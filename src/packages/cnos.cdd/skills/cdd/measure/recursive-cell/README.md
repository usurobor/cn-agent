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
# 0b7e26a2fce364c2f89bbd04a72e8e1068cb7339de9e9b9e9d33e520eff75645
```

`assemble-instruction.sh --output <path>` is the exact checked assembly
command. It refuses changed core, supplement, or composite bytes.

## Exact #662 R7 route

The TSC CLI resolves manifest paths from `--root`. Materialize the source CM
inside the frozen target checkout, then run its fail-closed preflight:

```bash
CNOS_SOURCE=/path/to/cnos-cycle-669
TSC_SOURCE=/path/to/tsc-at-26aab5023f03dc7d0abf82e5fdba20134fc6adad
TARGET=/tmp/cnos-662-r7
PROMPT=/tmp/cc662-system.prompt.md
RESPONSE=/tmp/cc662-system.response.json
REPORTS=/tmp/cc662-system-reports
CM_REL=src/packages/cnos.cdd/skills/cdd/measure/recursive-cell
COH="$TSC_SOURCE/engine/ocaml/_build/default/bin/main.exe"

git -C "$CNOS_SOURCE" worktree add --detach "$TARGET" \
  a0d39293a27cfe57b49dacff696345b1ee2cdb40
mkdir -p "$TARGET/$(dirname "$CM_REL")"
cp -a "$CNOS_SOURCE/$CM_REL" "$TARGET/$(dirname "$CM_REL")/"

"$TARGET/$CM_REL/calibration/662/verify-target.sh" "$TARGET"
"$CNOS_SOURCE/scripts/ci/validate-tsc-targets.sh" \
  --registry "$TARGET/$CM_REL/calibration/662/registry.tsc" \
  --target-root "$TARGET" --target cc662-system

"$COH" --mode llm --target cc662-system \
  --registry "$TARGET/$CM_REL/calibration/662/registry.tsc" \
  --instruction "$TARGET/$CM_REL/INSTRUCTION.md" --root "$TARGET" \
  --emit-prompt "$PROMPT"
sha256sum "$PROMPT" >"$PROMPT.sha256"
sha256sum -c "$PROMPT.sha256"

# Produce RESPONSE from exactly PROMPT, then reverify before ingestion.
sha256sum -c "$PROMPT.sha256"
"$COH" --mode hybrid --target cc662-system \
  --registry "$TARGET/$CM_REL/calibration/662/registry.tsc" \
  --instruction "$TARGET/$CM_REL/INSTRUCTION.md" --root "$TARGET" \
  --llm-response "$RESPONSE" --output "$REPORTS"
```

Repeat the mechanical command three times and use `coh consistency-spread`
for three semantic responses before making any consistency claim. The bundled
self target demonstrates prompt emission, digest verification, three raw
mechanical runs, and external-response smoke ingestion.

## Path-base boundary

Authority paths are repository-root-relative; `--root` is the explicit target
root. Installed-package activation is deliberately refused until TSC supports
separate authority and target bases:

```bash
"$CM/resolve-authority.sh" --source-root "$PWD"
./scripts/ci/smoke-recursive-cell-package.sh
```

Arbitrary methodology declaration loading remains unshipped. Do not use #662
as a held-out target: it designed this CM and confers no standing.
