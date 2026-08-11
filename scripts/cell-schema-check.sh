#!/usr/bin/env bash
# Shared Go/CUE/CLI corpus for the cell runner (Pi PR-#718 β D6): one gate that
# vets the cell-spec and episode-envelope schemas against a positive + negative
# fixture corpus AND vets actual `cn cell run` output against the terminal
# schema. Run from the repo root. Requires `cue` and a built `cn` (env CN).
set -uo pipefail
cd "$(dirname "$0")/.."

CUE=${CUE:-cue}
fail=0

tmpdir=$(mktemp -d)
tmp="$tmpdir/envelope.json" # .json so `cue` infers the format, not CUE
trap 'rm -rf "$tmpdir"' EXIT

# The CLI is one of the two authorities this corpus checks, so it must be the
# CLI BUILT FROM THE SOURCE UNDER REVIEW. Unless CN names one explicitly, build
# it here: the previous default of `./cn` ran whatever binary happened to be in
# the repo root, so a local run reported every CLI check green against a binary
# that could predate the change entirely. CI builds first and was never
# affected — this closes the gap between what CI measures and what a local run
# measures, which is the whole point of a shared corpus.
if [ -n "${CN:-}" ]; then
  echo "# using CN=$CN (caller-supplied; not rebuilt)"
else
  CN=$tmpdir/cn
  go -C src/go build -o "$CN" ./cmd/cn || { echo "✗ cannot build cn from source" >&2; exit 1; }
fi

# A missing tool must not read as a corpus that passed. Every negative below is
# an expected NON-ZERO exit, so an absent `cue` or `cn` would satisfy all of
# them vacuously — the failure mode this gate exists to prevent.
for tool in "$CUE" "$CN"; do
  command -v "$tool" >/dev/null 2>&1 || { echo "✗ required tool not found: $tool" >&2; exit 1; }
done

# files_exist guards the negative helpers below. A negative asserts a NON-ZERO
# exit, and a missing or misnamed fixture produces one too — so without this a
# deleted or typo'd fixture reads as the schema correctly rejecting it, which
# is the same vacuity class as an absent tool (Pi #55 C2).
files_exist() {
  local ok=0 a
  for a in "$@"; do
    case "$a" in
      *.json|*.cue) [ -f "$a" ] || { echo "  ✗ fixture not found: $a"; ok=1; } ;;
    esac
  done
  return $ok
}

vet_ok() {
  if ! "$CUE" vet "$@" >/dev/null 2>&1; then echo "  ✗ expected PASS: cue vet $*"; fail=1; else echo "  ✓ vet $*"; fi
}
vet_bad() {
  if ! files_exist "$@"; then fail=1; return; fi
  if "$CUE" vet "$@" >/dev/null 2>&1; then echo "  ✗ expected FAIL: cue vet $*"; fail=1; else echo "  ✓ rejected $*"; fi
}
run_bad() { # Go-only negatives: the CLI is the executable authority (exit 2).
  # Exit 2 is ALSO the runner's missing-contract exit, so the fixture has to be
  # proven present or this helper cannot tell rejection from absence.
  if ! files_exist "$1"; then fail=1; return; fi
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
# The CDS patch cell (fill-owned seats) vets against BOTH oracles: the generic
# tagged envelope and the CDS overlay's closed seat shapes.
vet_ok schemas/cdd/spec.cue schemas/cds/fixtures/code-cell-spec.json -d '#CellSpec'
vet_ok ./schemas/cds:cds schemas/cds/fixtures/code-cell-spec.json -d '#CDSCellSpec'
# The authored cognition shape admits what is STRUCTURALLY POSSIBLE before
# resolution (Pi #59 B1) — wider than what will succeed, since a provider hole
# may resolve to claude-cli and then fail construction. A literal fake whose
# model is a hole resolving to empty, and a provider hole with the model
# omitted, are both structurally possible; both were rejected before this arm.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/fake-model-hole-cell-spec.json -d '#CDSCellSpec'
vet_ok ./schemas/cds:cds schemas/cds/fixtures/provider-hole-cell-spec.json -d '#CDSCellSpec'
# Case 3 differs from Case 2 by ONE field: beta.fill. Both must vet.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/reviewed-code-cell-spec.json -d '#CDSCellSpec'

echo "# positive closures"
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-accepted.json -d '#EpisodeClosure'
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-needs-repair.json -d '#EpisodeClosure'
vet_ok schemas/cdd/episode-closure.cue schemas/cdd/fixtures/episode-closure-simulated.json -d '#EpisodeClosure'
# A committed SCHEMA EXAMPLE of a Case-2 closure — deliberately not a proof.
# Episode ids are minted per invocation, so a stored file can never be
# byte-reproduced; the authority is the live run below, which executes the
# committed spec from an installed hub and vets both the closure and its
# resolved declaration.
vet_ok schemas/cdd/episode-closure.cue schemas/cds/fixtures/episode-closure-cds-case2.example.json -d '#EpisodeClosure'

echo "# negative cell specs (must be rejected)"
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-bad-producer.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-missing-fill.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-empty-goal.json -d '#CellSpec'
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-case-alias.json -d '#CellSpec'
# One identifier grammar: a param name legal in Go must be legal in CUE, or a
# spec resolves in one authority and is rejected by the other (Pi #55 C1).
vet_bad schemas/cdd/spec.cue schemas/cdd/fixtures/invalid/cellspec-bad-param-name.json -d '#CellSpec'
vet_bad schemas/cdd/episode-closure.cue schemas/cdd/fixtures/invalid/episode-closure-null-arrays.json -d '#EpisodeClosure'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-no-diff.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-diff-not-first.json -d '#CDSCellSpec'
# Fill-owned strictness: null skill lists and smuggled provider argv are
# rejected by the CDS overlay (and by the fill decoder below).
vet_bad ./schemas/cds:cds schemas/cdd/fixtures/invalid/cellspec-null-skills.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-smuggled-argv.json -d '#CDSCellSpec'
# A real provider without a model selector must be rejected by BOTH authorities,
# and so must a fake carrying a model it would ignore.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-modelless-provider.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-fake-with-model.json -d '#CDSCellSpec'
# codex-cli is held out of the admitted provider set until its ambient-context
# suppression can be proven by a real run (Pi #55 D1).
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-codex-held.json -d '#CDSCellSpec'
# A hole-capable CDS field must not accept a malformed hole through an
# unrestricted string arm: Go reads every $... as a hole (Pi #56 C1).
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-bad-hole-name.json -d '#CDSCellSpec'
# ...including in the cognition model position, which was bare `string`.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-bad-model-hole.json -d '#CDSCellSpec'
# A reviewer's canonical input is (contract, matter). Declaring a workspace
# would let it reach the very worktree it judges, so the key is not admitted.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-review-with-workspace.json -d '#CDSCellSpec'
# Fill-owned keys are exact and case-sensitive at every depth: encoding/json
# would otherwise decode these while the closed overlay rejects them.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-case-seat-tag.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-case-top-arg.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-case-nested-arg.json -d '#CDSCellSpec'
# A declared task must be an admissible issue, and knowing what one IS belongs
# to the CDS overlay: the same malformed file vets clean against the generic
# #CellSpec, which is what shows the boundary sits where it is claimed to sit
# rather than having leaked into the kernel's schema.
#
# Absence is a different case and is deliberately NOT checked here — see
# #CDSCellSpec on why `task` is optional in this schema and mandatory at the
# door. The runtime witness for it is below, with the live cells.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-malformed-task.json -d '#CDSCellSpec'
vet_ok schemas/cdd/spec.cue schemas/cds/fixtures/invalid/cds-malformed-task.json -d '#CellSpec'

echo "# CDS issue corpus (one corpus, two authorities)"
# cdsissue.Admit and #CDSIssue must accept and reject exactly the same
# documents, so both read THESE files: schemas/cds/fixtures/issue/ is vetted
# here and table-tested by internal/cdsissue. Each negative is invalid for
# exactly ONE reason — a fixture breaking two rules cannot show which fired —
# and the Go test additionally pins WHICH rule that is.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/issue/valid-issue.json -d '#CDSIssue'
# scope.out PRESENT but empty is admissible: non-goals are load-bearing, so an
# empty list says "considered, none" where an absent key says nothing at all.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/issue/valid-empty-scope-out.json -d '#CDSIssue'
# blank-unicode-whitespace carries EVERY rune unicode.IsSpace covers in one
# field. It is the witness that #NonBlank and cdsissue's nonBlankPattern are
# the same class: a rune missing from the CUE enumeration would make this
# document vet clean while Go rejects it, which is the divergence the two
# authorities exist to catch.
for neg in bad-kind empty-id blank-problem-line blank-unicode-whitespace \
           no-sources source-without-path \
           empty-scope-in missing-scope-out no-acceptance \
           criterion-without-verification duplicate-acceptance-id \
           unknown-key mixed-case-key; do
  vet_bad ./schemas/cds:cds "schemas/cds/fixtures/issue/issue-$neg.json" -d '#CDSIssue'
done

echo "# Go-only negatives (executable authority = the CLI)"
run_bad schemas/cdd/fixtures/invalid/cellspec-dup-required-id.json
run_bad schemas/cdd/fixtures/invalid/cellspec-bad-producer.json
run_bad schemas/cdd/fixtures/invalid/cellspec-missing-fill.json
run_bad schemas/cdd/fixtures/invalid/cellspec-unknown-fill.json
run_bad schemas/cdd/fixtures/invalid/cellspec-empty-goal.json
run_bad schemas/cdd/fixtures/invalid/cellspec-case-alias.json
run_bad schemas/cdd/fixtures/invalid/cellspec-bad-param-name.json
run_bad schemas/cdd/fixtures/invalid/cellspec-null-skills.json
run_bad schemas/cds/fixtures/invalid/cds-smuggled-argv.json
run_bad schemas/cds/fixtures/invalid/cds-modelless-provider.json
run_bad schemas/cds/fixtures/invalid/cds-fake-with-model.json
run_bad schemas/cds/fixtures/invalid/cds-codex-held.json
run_bad schemas/cds/fixtures/invalid/cds-bad-hole-name.json
run_bad schemas/cds/fixtures/invalid/cds-bad-model-hole.json
run_bad schemas/cds/fixtures/invalid/cds-case-seat-tag.json
run_bad schemas/cds/fixtures/invalid/cds-case-top-arg.json
run_bad schemas/cds/fixtures/invalid/cds-case-nested-arg.json

echo "# committed rented-Claude evidence (one-off receipt, NOT a provider run)"
# The live corpus rents `fake`, so the cognitive path has no runtime witness
# here. The committed closure below IS that witness, and these checks make it
# accountable rather than narrative (Pi #58 D1): both CUE oracles are
# recomputed from the artifact, and the measurement is recomputed from the
# diff the artifact carries.
#
# Scoped exactly (Pi #60 D1): a one-byte edit INSIDE `matter.data`, or deleting
# the artifact, fails this block. Metadata and JSON whitespace are not covered,
# and nothing here claims they are. Nothing in this gate invokes a RENTED
# cognition provider — the corpus does invoke the deterministic `fake`, which
# is the point of it being deterministic.
ev=docs/architecture/evidence/cds-case2-claude-closure.json
ev_diff_sha=3826a7e883a9fb78769d1ef99ca54a16bad631aea244620412e2d5be58261766
ev_diff_bytes=2479
ev_touched="CONTRIBUTING.md README.md"
if ! files_exist "$ev"; then
  fail=1
else
  vet_ok schemas/cdd/episode-closure.cue "$ev" -d '#EpisodeClosure'
  evalpha="$tmpdir/evidence-alpha.json"
  if python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["alpha"], open(sys.argv[2],"w"))' "$ev" "$evalpha" 2>/dev/null; then
    vet_ok ./schemas/cds:cds "$evalpha" -d '#CDSPatchAlphaResolved'
  else
    echo "  ✗ evidence closure has no resolved alpha"; fail=1
  fi
  if python3 - "$ev" "$ev_diff_sha" "$ev_diff_bytes" "$ev_touched" <<'PYEOF'
import hashlib, json, sys
c = json.load(open(sys.argv[1]))
r = c["receipt"]["record"]
d = r["matter"]["data"]
want_sha, want_bytes, want_touched = sys.argv[2], int(sys.argv[3]), sys.argv[4].split()
got_touched = sorted({l.split(" b/")[-1] for l in d.splitlines() if l.startswith("diff --git")})
checks = [
    ("execution_mode", r["execution_mode"], "cognitive"),
    ("status", c["status"], "needs_repair"),
    # UTF-8 BYTES, not code points: len(d) undercounted this diff by 4
    # because it carries two em dashes, and Go — which wrote the record —
    # measures bytes (Pi #59 D1).
    ("diff bytes", len(d.encode()), want_bytes),
    ("diff sha256", hashlib.sha256(d.encode()).hexdigest(), want_sha),
    ("touched files", got_touched, sorted(want_touched)),
]
bad = [(n, g, w) for n, g, w in checks if g != w]
for n, g, w in bad:
    print(f"    {n}: got {g!r}, want {w!r}")
sys.exit(1 if bad else 0)
PYEOF
  then
    echo "  ✓ evidence measurement recomputes (mode, status, $ev_diff_bytes UTF-8 diff bytes, diff digest, touched files)"
  else
    echo "  ✗ evidence closure measurement does NOT recompute"; fail=1
  fi
fi

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
) >/dev/null 2>&1 || { echo "  ✗ could not build the code-cell fixture repo"; fail=1; }
# Skill authority is the INSTALLED package root under a hub — never the
# working directory and never this checkout's source tree. Vendor the DEFAULT
# INSTALLED PACKAGE SET (what `cn repo install` pins) into a throwaway hub and
# run from inside it, so this proves a hub the product can actually produce
# rather than a hand-picked fixture. repoinstall.DefaultPackages is the source
# of truth; TestDefaultPackagesCoverShippedCells keeps it and the shipped
# cells in agreement.
hub="$tmpdir/hub"
mkdir -p "$hub/.cn/vendor/packages"
default_packages=$(grep -o 'var DefaultPackages = \[\]string{[^}]*}' src/go/internal/repoinstall/repoinstall.go |
  grep -o '"[^"]*"' | tr -d '"')
if [ -z "$default_packages" ]; then
  echo "  ✗ could not read DefaultPackages"; fail=1
fi
for pkg in $default_packages; do
  [ -d "src/packages/$pkg/skills" ] || continue
  mkdir -p "$hub/.cn/vendor/packages/$pkg"
  cp -r "src/packages/$pkg/skills" "$hub/.cn/vendor/packages/$pkg/skills" ||
    { echo "  ✗ could not vendor $pkg into the fixture hub"; fail=1; }
done
# Absolute paths for the subshell that runs INSIDE the hub; CUE package paths
# stay relative to the repo root, where the rest of this script runs.
CN=$(cd "$(dirname "$CN")" && pwd)/$(basename "$CN")
spec=$(pwd)/schemas/cds/fixtures/code-cell-spec.json
(
  cd "$hub" || exit 1
  "$CN" cell run --contract "$spec" \
    --param language=cnos.eng:eng/go --param provider=fake \
    --param base_sha=HEAD --param repo="$coderepo" >"$tmp" 2>/dev/null
  echo $? >"$tmpdir/code.exit"
)
code=$(cat "$tmpdir/code.exit")
if [ "$code" != 1 ]; then
  echo "  ✗ cds.patch cell exit=$code want=1 (needs_repair: beta cannot judge the goal)"; fail=1
else echo "  ✓ cds.patch cell closes needs_repair from an installed hub"; fi
if ! "$CUE" vet schemas/cdd/episode-closure.cue "$tmp" -d '#EpisodeClosure' >/dev/null 2>&1; then
  echo "  ✗ cds.patch closure failed #EpisodeClosure"; fail=1
else echo "  ✓ cds.patch closure vets #EpisodeClosure"; fi
decl="$tmpdir/resolved-alpha.json"
if python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["alpha"], open(sys.argv[2],"w"))' "$tmp" "$decl" 2>/dev/null &&
   "$CUE" vet ./schemas/cds:cds "$decl" -d '#CDSPatchAlphaResolved' >/dev/null 2>&1; then
  echo "  ✓ resolved alpha vets #CDSPatchAlphaResolved (canonical shape, pinned base, digested skills)"
else echo "  ✗ resolved alpha failed #CDSPatchAlphaResolved"; fail=1; fi

# Case 3 end to end from the same installed hub: alpha produces, and a
# cds.review beta is CONSTRUCTED and INVOKED rather than stubbed out. The
# reviewer rents `fake`, which never passes, so this proves wiring — seat
# construction, skill loading, prompt, verdict decode, closure shape — not
# semantic review. Semantic review is proven by the rented evidence below.
c3spec=$(pwd)/schemas/cds/fixtures/reviewed-code-cell-spec.json
c3out="$tmpdir/case3.json"
(
  cd "$hub" || exit 1
  "$CN" cell run --contract "$c3spec" \
    --param language=cnos.eng:eng/go --param provider=fake \
    --param base_sha=HEAD --param repo="$coderepo" >"$c3out" 2>/dev/null
  echo $? >"$tmpdir/c3.exit"
)
c3code=$(cat "$tmpdir/c3.exit")
if [ "$c3code" != 1 ]; then
  echo "  ✗ case-3 cell exit=$c3code want=1 (fake reviewer never passes)"; fail=1
else echo "  ✓ case-3 cell closes needs_repair with a constructed cds.review beta"; fi
vet_ok schemas/cdd/episode-closure.cue "$c3out" -d '#EpisodeClosure'
# The runtime half of issue admission, and the reason #CDSCellSpec can leave
# `task` optional: the SAME spec with its task removed vets clean and does not
# run. Asserted by REASON rather than by exit code — exit 2 is also the
# missing-contract and unresolvable-base exit, so a bare code would prove
# refusal without proving cause. The seat must refuse before renting anything.
nt="$(pwd)/schemas/cds/fixtures/invalid/cds-no-task.json"
if ! files_exist "$nt"; then fail=1; fi
vet_ok ./schemas/cds:cds "$nt" -d '#CDSCellSpec'
(
  cd "$hub" || exit 1
  "$CN" cell run --contract "$nt" \
    --param language=cnos.eng:eng/go --param provider=fake \
    --param base_sha=HEAD --param repo="$coderepo" >/dev/null 2>"$tmpdir/nt.err"
  echo $? >"$tmpdir/nt.exit"
)
if [ "$(cat "$tmpdir/nt.exit")" = 0 ] || ! grep -q "contract carries no task" "$tmpdir/nt.err"; then
  echo "  ✗ a taskless CDS cell must be refused at admission, got exit $(cat "$tmpdir/nt.exit"): $(head -c 200 "$tmpdir/nt.err")"; fail=1
else echo "  ✓ a taskless CDS cell vets and is refused at the door"; fi

c3beta="$tmpdir/case3-beta.json"
if python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["beta"], open(sys.argv[2],"w"))' "$c3out" "$c3beta" 2>/dev/null; then
  vet_ok ./schemas/cds:cds "$c3beta" -d '#CDSReviewBetaResolved'
else
  echo "  ✗ case-3 closure has no resolved beta"; fail=1
fi


if [ "$fail" = 0 ]; then echo "✓ cell schema/CLI corpus OK"; else echo "✗ cell schema check FAILED"; exit 1; fi
