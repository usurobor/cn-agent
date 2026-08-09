#!/usr/bin/env bash
# Shared Go/CUE/CLI corpus for the cell runner (Pi PR-#718 β D6): one gate that
# vets the cell-spec and episode-envelope schemas against a positive + negative
# fixture corpus AND vets actual `cn cell run` output against the terminal
# schema. Run from the repo root. Requires `cue` and a built `cn` (env CN).
set -uo pipefail
cd "$(dirname "$0")/.."

CUE=${CUE:-cue}
CN=${CN:-./cn}
fail=0
tmpdir=$(mktemp -d)
tmp="$tmpdir/envelope.json" # .json so `cue` infers the format, not CUE
trap 'rm -rf "$tmpdir"' EXIT

vet_ok() {
  if ! "$CUE" vet "$@" >/dev/null 2>&1; then echo "  ✗ expected PASS: cue vet $*"; fail=1; else echo "  ✓ vet $*"; fi
}
vet_bad() {
  if "$CUE" vet "$@" >/dev/null 2>&1; then echo "  ✗ expected FAIL: cue vet $*"; fail=1; else echo "  ✓ rejected $*"; fi
}
run_vet() { # want-exit, cn args...
  local want=$1; shift
  "$CN" cell run "$@" >"$tmp" 2>/dev/null; local code=$?
  if [ "$code" != "$want" ]; then echo "  ✗ cn cell run exit=$code want=$want ($*)"; fail=1; fi
  if ! "$CUE" vet schemas/cdd/episode-envelope.cue "$tmp" -d '#EpisodeEnvelope' >/dev/null 2>&1; then
    echo "  ✗ CLI output failed #EpisodeEnvelope ($*)"; fail=1
  else echo "  ✓ CLI output vets ($*)"; fi
}

echo "# positive cell specs"
vet_ok schemas/cdd/spec.cue schemas/cdd/fixtures/empty-cell-spec.json -d '#CellSpec'
vet_ok schemas/cdd/spec.cue schemas/cdd/fixtures/bool-cell-spec.json -d '#CellSpec'
vet_ok ./schemas/cds:cds schemas/cds/fixtures/valid-cell-spec.json -d '#CDSCellSpec'

echo "# positive envelopes"
vet_ok schemas/cdd/episode-envelope.cue schemas/cdd/fixtures/episode-envelope-accepted.json -d '#EpisodeEnvelope'
vet_ok schemas/cdd/episode-envelope.cue schemas/cdd/fixtures/episode-envelope-needs-repair.json -d '#EpisodeEnvelope'

echo "# negative cell specs (must be rejected)"
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-bad-producer.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-missing-profile.json -d '#CellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-no-diff.json -d '#CDSCellSpec'

echo "# actual CLI output vetted against the terminal schema"
run_vet 0 --contract schemas/cdd/fixtures/bool-cell-spec.json --param value=true
run_vet 1 --contract schemas/cdd/fixtures/bool-cell-spec.json --param value=false

if [ "$fail" = 0 ]; then echo "✓ cell schema/CLI corpus OK"; else echo "✗ cell schema check FAILED"; exit 1; fi
