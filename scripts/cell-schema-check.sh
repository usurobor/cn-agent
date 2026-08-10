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
# A resolved seat declaration must survive its fill's RESOLVED schema: no holes
# left, base pinned to a commit, skills carrying content digests. This proves
# resolution actually happened instead of trusting that it did.
vet_resolved_alpha() { # definition, cn args...
  local def=$1; shift
  "$CN" cell run "$@" >"$tmp" 2>/dev/null
  local decl="$tmpdir/resolved-alpha.json"
  if ! python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["alpha"], open(sys.argv[2],"w"))' "$tmp" "$decl" 2>/dev/null; then
    echo "  ✗ could not extract the resolved alpha declaration ($*)"; fail=1; return
  fi
  if ! "$CUE" vet ./schemas/cds:cds "$decl" -d "$def" >/dev/null 2>&1; then
    echo "  ✗ resolved alpha failed $def ($*)"; fail=1
  else echo "  ✓ resolved alpha vets $def"; fi
}

echo "# positive cell specs"
vet_ok schemas/cdd/spec.cue schemas/cdd/fixtures/empty-cell-spec.json -d '#CellSpec'
vet_ok schemas/cdd/spec.cue schemas/cdd/fixtures/bool-cell-spec.json -d '#CellSpec'
# The CDS patch cell (fill-owned seats) vets against BOTH oracles: the generic
# tagged envelope and the CDS overlay's closed seat shapes.
vet_ok schemas/cdd/spec.cue schemas/cds/fixtures/code-cell-spec.json -d '#CellSpec'
vet_ok ./schemas/cds:cds schemas/cds/fixtures/code-cell-spec.json -d '#CDSCellSpec'

echo "# positive closures"
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-accepted.json -d '#EpisodeClosure'
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-needs-repair.json -d '#EpisodeClosure'
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-simulated.json -d '#EpisodeClosure'
# The committed Case-2 closure, reproducible from its committed input:
#   cn cell run --contract schemas/cds/fixtures/code-cell-spec.json \
#     --param language=cnos.eng:eng/go --param provider=fake --param base_sha=<head>
vet_ok schemas/cdd/episode-closure.cue schemas/cds/fixtures/episode-closure-cds-case2.json -d '#EpisodeClosure'

echo "# negative cell specs (must be rejected)"
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-bad-producer.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-missing-fill.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-empty-goal.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-case-alias.json -d '#CellSpec'
vet_bad schemas/cdd/episode-closure.cue schemas/cdd/fixtures/invalid/episode-closure-null-arrays.json -d '#EpisodeClosure'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-no-diff.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-diff-not-first.json -d '#CDSCellSpec'
# Fill-owned strictness: null skill lists and smuggled provider argv are
# rejected by the CDS overlay (and by the fill decoder below).
vet_bad ./schemas/cds:cds schemas/cdd/fixtures/invalid/cellspec-null-skills.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-smuggled-argv.json -d '#CDSCellSpec'
# A real provider without an exact model must be rejected by BOTH authorities.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-modelless-provider.json -d '#CDSCellSpec'

echo "# Go-only negatives (executable authority = the CLI)"
run_bad schemas/cdd/fixtures/invalid/cellspec-dup-required-id.json
run_bad schemas/cdd/fixtures/invalid/cellspec-bad-producer.json
run_bad schemas/cdd/fixtures/invalid/cellspec-missing-fill.json
run_bad schemas/cdd/fixtures/invalid/cellspec-unknown-fill.json
run_bad schemas/cdd/fixtures/invalid/cellspec-empty-goal.json
run_bad schemas/cdd/fixtures/invalid/cellspec-case-alias.json
run_bad schemas/cdd/fixtures/invalid/cellspec-null-skills.json
run_bad schemas/cds/fixtures/invalid/cds-smuggled-argv.json
run_bad schemas/cds/fixtures/invalid/cds-modelless-provider.json

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
# The cds.patch cell against a hermetic throwaway repository: the runtime cuts
# a worktree, the fake coder changes a file, the diff in the closure is
# MEASURED from that worktree — and the episode still closes needs_repair
# (exit 1), because the mechanical-unmet beta cannot judge the goal. Case-2
# honesty is part of the corpus.
coderepo="$tmpdir/coderepo"
mkdir -p "$coderepo"
(
  cd "$coderepo" && git init -q -b main && echo base >README.md && git add -A &&
    GIT_AUTHOR_NAME=t GIT_AUTHOR_EMAIL=t@t GIT_COMMITTER_NAME=t GIT_COMMITTER_EMAIL=t@t \
      git commit -qm base
) >/dev/null 2>&1 || { echo "  ✗ could not build the code-profile fixture repo"; fail=1; }
run_vet 1 --contract schemas/cds/fixtures/code-cell-spec.json \
  --param language=cnos.eng:eng/go --param provider=fake --param base_sha=HEAD --param repo="$coderepo"
vet_resolved_alpha '#CDSPatchAlphaResolved' --contract schemas/cds/fixtures/code-cell-spec.json \
  --param language=cnos.eng:eng/go --param provider=fake --param base_sha=HEAD --param repo="$coderepo"

if [ "$fail" = 0 ]; then echo "✓ cell schema/CLI corpus OK"; else echo "✗ cell schema check FAILED"; exit 1; fi
