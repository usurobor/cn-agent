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
run_bad() { # Go-only negatives: the CLI is the executable authority (exit 2).
  "$CN" cell run --contract "$1" >/dev/null 2>&1; local code=$?
  if [ "$code" != 2 ]; then echo "  ✗ expected CLI exit 2: $1 (got $code)"; fail=1; else echo "  ✓ CLI rejected $1"; fi
}
run_vet() { # want-exit, cn args...
  local want=$1; shift
  "$CN" cell run "$@" >"$tmp" 2>/dev/null; local code=$?
  if [ "$code" != "$want" ]; then echo "  ✗ cn cell run exit=$code want=$want ($*)"; fail=1; fi
  if ! "$CUE" vet schemas/cdd/episode-closure.cue "$tmp" -d '#EpisodeClosure' >/dev/null 2>&1; then
    echo "  ✗ CLI output failed #EpisodeClosure ($*)"; fail=1
  else echo "  ✓ CLI output vets ($*)"; fi
}

echo "# positive cell specs"
vet_ok schemas/cdd/spec.cue schemas/cdd/fixtures/empty-cell-spec.json -d '#CellSpec'
vet_ok schemas/cdd/spec.cue schemas/cdd/fixtures/bool-cell-spec.json -d '#CellSpec'
vet_ok ./schemas/cds:cds schemas/cds/fixtures/valid-cell-spec.json -d '#CDSCellSpec'

echo "# positive closures"
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-accepted.json -d '#EpisodeClosure'
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-needs-repair.json -d '#EpisodeClosure'
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-simulated.json -d '#EpisodeClosure'

echo "# negative cell specs (must be rejected)"
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-bad-producer.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-missing-profile.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-empty-goal.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-missing-skills.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-case-alias.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-null-skills.json -d '#CellSpec'
vet_bad schemas/cdd/episode-closure.cue schemas/cdd/fixtures/invalid/episode-closure-null-arrays.json -d '#EpisodeClosure'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-no-diff.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-diff-not-first.json -d '#CDSCellSpec'

echo "# Go-only negatives (executable authority = the CLI)"
run_bad schemas/cdd/fixtures/invalid/cellspec-dup-required-id.json
run_bad schemas/cdd/fixtures/invalid/bool-missing-value.json
run_bad schemas/cdd/fixtures/invalid/cellspec-bad-producer.json
run_bad schemas/cdd/fixtures/invalid/cellspec-missing-profile.json
run_bad schemas/cdd/fixtures/invalid/cellspec-empty-goal.json
run_bad schemas/cdd/fixtures/invalid/cellspec-missing-skills.json
run_bad schemas/cdd/fixtures/invalid/cellspec-case-alias.json
run_bad schemas/cdd/fixtures/invalid/cellspec-null-skills.json

echo "# SIGINT terminates a blocked stdin reader (Pi round-5 D3)"
mkfifo "$tmpdir/stdin.fifo"
# set -m: without job control, a non-interactive shell starts background jobs
# with SIGINT ignored, which the child inherits across exec — the kill below
# would silently test nothing.
set -m
"$CN" cell run --contract - <"$tmpdir/stdin.fifo" >/dev/null 2>&1 & cnpid=$!
set +m
exec 9>"$tmpdir/stdin.fifo" # hold the writer open: no EOF, reader stays blocked
sleep 1
# The child must still be blocked when signaled — an early exit would make the
# regression pass vacuously (Pi round-6 B1).
if ! kill -0 "$cnpid" 2>/dev/null; then
  echo "  ✗ blocked reader exited before SIGINT was sent"; fail=1
fi
kill -INT "$cnpid" 2>/dev/null
deadline=$((SECONDS + 5))
while kill -0 "$cnpid" 2>/dev/null && [ "$SECONDS" -lt "$deadline" ]; do sleep 0.1; done
if kill -0 "$cnpid" 2>/dev/null; then
  echo "  ✗ SIGINT did not terminate blocked cn cell run"; kill -9 "$cnpid" 2>/dev/null; fail=1
else echo "  ✓ SIGINT terminates a blocked reader"; fi
wait "$cnpid" 2>/dev/null
exec 9>&-

echo "# actual CLI output vetted against the terminal schema"
run_vet 0 --contract schemas/cdd/fixtures/bool-cell-spec.json --param value=true
run_vet 1 --contract schemas/cdd/fixtures/bool-cell-spec.json --param value=false
run_vet 3 --contract schemas/cdd/fixtures/empty-cell-spec.json

if [ "$fail" = 0 ]; then echo "✓ cell schema/CLI corpus OK"; else echo "✗ cell schema check FAILED"; exit 1; fi
