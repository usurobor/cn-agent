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
# the repository root, so a local run reported every CLI check green against a
# binary that could predate the change entirely — and a mutation test passed
# after the guard it was testing had been deleted. CI builds first and was
# never affected; this closes the gap between what CI measures and what a local
# run measures, which is the whole point of a shared corpus.
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

repo=$(pwd)

# THE FIXTURE HUB, built before anything invokes the CLI. Skill authority is the
# INSTALLED package root under a hub — never the working directory and never
# this checkout's source tree — so every `cn cell run` below runs from inside it.
#
# It has to exist this early now that a cell declares ONE methodology bundle: the
# runtime loads that bundle before it constructs either seat, so a CLI negative
# run from a directory with no installed packages is rejected for an uninstalled
# skill and never reaches the defect its fixture is named for. Six of them did
# exactly that, and `run_bad`'s reason check is what said so.
#
# The vendored set is the DEFAULT INSTALLED PACKAGE SET (what `cn repo install`
# pins), so this proves a hub the product can actually produce rather than a
# hand-picked fixture. repoinstall.DefaultPackages is the source of truth;
# TestDefaultPackagesCoverShippedCells keeps it and the shipped cells in
# agreement.
hub="$tmpdir/hub"
mkdir -p "$hub/.cn/vendor/packages"
default_packages=$(grep -o 'var DefaultPackages = \[\]string{[^}]*}' src/go/internal/repoinstall/repoinstall.go |
  grep -o '"[^"]*"' | tr -d '"')
if [ -z "$default_packages" ]; then
  echo "✗ could not read DefaultPackages" >&2; exit 1
fi
for pkg in $default_packages; do
  [ -d "src/packages/$pkg/skills" ] || continue
  mkdir -p "$hub/.cn/vendor/packages/$pkg"
  cp -r "src/packages/$pkg/skills" "$hub/.cn/vendor/packages/$pkg/skills" ||
    { echo "✗ could not vendor $pkg into the fixture hub" >&2; exit 1; }
done
# Absolute, for the subshells that run INSIDE the hub; CUE package paths stay
# relative to the repo root, where the rest of this script runs.
CN=$(cd "$(dirname "$CN")" && pwd)/$(basename "$CN")
spec=$repo/schemas/cds/fixtures/code-cell-spec.json

# cn_run invokes the CLI from inside the fixture hub, rewriting every argument
# that names an existing repo-relative file to an absolute path. ONE helper, so
# no caller can accidentally run from a directory where skills do not resolve —
# which is a green tick for the wrong reason, not a failure.
cn_run() {
  local -a argv=()
  local a
  for a in "$@"; do
    if [ -f "$repo/$a" ]; then argv+=("$repo/$a"); else argv+=("$a"); fi
  done
  (cd "$hub" && "$CN" "${argv[@]}")
}

# THE HERMETIC CODE REPOSITORY and the live run input, built here for the same
# reason the hub is: every CLI invocation below runs from inside the hub, and a
# run input whose subject names `.` would then pin the hub rather than a
# repository. The committed fixture names `.`, which is what an author writes;
# pinning really opens it, so the live witness points it at this throwaway
# repository. Only the subject is rewritten — the issue and the design are the
# corpus documents both authorities vetted.
coderepo="$tmpdir/coderepo"
mkdir -p "$coderepo"
(
  cd "$coderepo" && git init -q -b main && echo base >README.md && git add -A &&
    GIT_AUTHOR_NAME=t GIT_AUTHOR_EMAIL=t@t GIT_COMMITTER_NAME=t GIT_COMMITTER_EMAIL=t@t \
      git commit -qm base
) >/dev/null 2>&1 || { echo "✗ could not build the code-cell fixture repo" >&2; exit 1; }
ri="$tmpdir/run-input.json"
python3 - "schemas/cds/fixtures/runinput/valid-run-input.json" "$coderepo" "$ri" <<'RIEOF' ||
import json, sys
doc = json.load(open(sys.argv[1]))
doc["subject"] = {"kind": "git.snapshot/0.1", "repo": sys.argv[2], "base_sha": "HEAD"}
json.dump(doc, open(sys.argv[3], "w"))
RIEOF
  { echo "✗ could not build the live run input" >&2; exit 1; }

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
  # Usage: run_bad <fixture> <reason substring> [extra cn args...]
  #
  # Exit 2 is ALSO the runner's missing-contract exit, the missing-run-input
  # exit, and now the unmet-fill-requirement exit, so the CODE proves refusal
  # and never proves CAUSE. A fixture must therefore name the reason it is
  # rejected for. This is not hypothetical: when the subject requirement moved
  # ahead of construction, six of these fixtures began short-circuiting there
  # and stopped exercising their own defect entirely, while this helper went on
  # printing a tick for each. Fixture presence is still checked, for the same
  # reason at one remove.
  local fixture=$1 reason=$2; shift 2
  if ! files_exist "$fixture"; then fail=1; return; fi
  local err; err=$(cn_run cell run --contract "$fixture" "$@" 2>&1 >/dev/null); local code=$?
  if [ "$code" != 2 ]; then
    echo "  ✗ expected CLI exit 2: $fixture (got $code)"; fail=1
  elif ! printf '%s' "$err" | grep -qF -- "$reason"; then
    echo "  ✗ rejected for the wrong reason: $fixture"
    echo "      got:  $err"
    echo "      want: $reason"; fail=1
  else
    echo "  ✓ CLI rejected $fixture ($reason)"
  fi
}
run_vet() { # want-exit, cn args...
  local want=$1; shift
  cn_run cell run "$@" >"$tmp" 2>/dev/null; local code=$?
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
# ONE METHODOLOGY, ONE PLACE. A cds.patch seat has no `skills` key: the cell
# declares the bundle and the seat receives a projection of it. A declaration
# that carries its own list is rejected by the closed overlay here and by the
# fill's exact-key set below — the two authorities delete the field together, or
# the deletion is only a Go convention.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-seat-skills.json -d '#CDSCellSpec'
# NO WORKSPACE ON THE REVIEWING SEAT, and this is the beta-side counterpart of
# the rule above. A reviewer that could name a repository could open the
# workspace it is meant to judge from the outside, which is the independence the
# seat exists to provide — so the key is absent from the closed overlay here and
# from the fill's exact-key set below.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-assess-workspace.json -d '#CDSCellSpec'
# Fill-owned keys are exact and case-sensitive at every depth: encoding/json
# would otherwise decode these while the closed overlay rejects them.
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-case-seat-tag.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-case-top-arg.json -d '#CDSCellSpec'
vet_bad ./schemas/cds:cds schemas/cds/fixtures/invalid/cds-case-nested-arg.json -d '#CDSCellSpec'

echo "# CDS subject corpus (one corpus, two authorities)"
# cellwork.AdmitSubject and #GitSnapshotPinned must accept and reject exactly
# the same documents, so both read THESE files: schemas/cds/fixtures/subject/ is
# vetted here and table-tested by internal/cellwork. Each negative is invalid
# for exactly ONE reason, and the Go test pins WHICH rule that is.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/subject/valid-subject.json -d '#GitSnapshotPinned'
for neg in bad-kind missing-repo empty-base unpinned-base unknown-key mixed-case-key; do
  vet_bad ./schemas/cds:cds "schemas/cds/fixtures/subject/subject-$neg.json" -d '#GitSnapshotPinned'
done
# The authored form is WIDER by exactly one rule: a base that is still a moving
# revision. It is admissible before pinning and inadmissible in a record, and
# this pair is what shows the two definitions are not the same definition.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/subject/subject-unpinned-base.json -d '#GitSnapshotAuthored'

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
for neg in bad-kind empty-id blank-problem-line blank-unicode-whitespace reserved-acceptance-id \
           no-sources source-without-path \
           empty-scope-in missing-scope-out no-acceptance \
           criterion-without-verification duplicate-acceptance-id \
           unknown-key mixed-case-key; do
  vet_bad ./schemas/cds:cds "schemas/cds/fixtures/issue/issue-$neg.json" -d '#CDSIssue'
done

echo "# CDS design corpus (one corpus, two authorities)"
# cdsdesign.Admit and #CDSDesign, same discipline as the issue above. The
# design is a SEPARATE document with a separate schema on purpose: an issue
# that could also state architecture would let the problem be redefined by the
# same act that proposes the change.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/design/valid-design.json -d '#CDSDesign'
for neg in bad-kind blank-approach no-invariants blank-invariant \
           no-impact blank-surface impact-without-why \
           unknown-key mixed-case-key; do
  vet_bad ./schemas/cds:cds "schemas/cds/fixtures/design/design-$neg.json" -d '#CDSDesign'
done

echo "# CDS run-input corpus (one corpus, two authorities)"
# The whole run contract in one document, vetted here and table-tested by
# internal/cdsadmit. The two authorities do NOT report identically and are not
# claimed to: Go distinguishes an absent payload (incomplete) from a malformed
# one (rejected), while CUE has one verdict. What must agree — and is what this
# block checks — is WHICH documents are admissible.
vet_ok ./schemas/cds:cds schemas/cds/fixtures/runinput/valid-run-input.json -d '#CDSRunInput'
for neg in bad-kind no-issue no-design no-subject \
           malformed-issue malformed-design wrong-kind-subject \
           unknown-key mixed-case-key; do
  vet_bad ./schemas/cds:cds "schemas/cds/fixtures/runinput/runinput-$neg.json" -d '#CDSRunInput'
done

echo "# every negative is invalid for exactly ONE reason"
# The rejections above prove only that each negative is inadmissible — not that
# it is inadmissible for the single reason its name claims. A fixture that
# broke two rules would still be rejected, and the Go table's expected
# substring would still match the first rule to fire, so both gates would stay
# green while the corpus quietly stopped isolating defects.
#
# The check is structural and needs no repair table: every negative is derived
# from its directory's positive by ONE mutation, so it must differ from that
# positive in exactly one place. Restoring that one place reproduces the
# positive document itself — which both authorities accept above — so "repair
# the one defect and both authorities pass" is proven rather than asserted.
#
# A renamed key (`kind` → `Kind`) is one defect spelled as a removal plus an
# addition, and is counted as one only when the two names differ by case alone
# AND carry the same value.
if python3 - <<'PYEOF'
import json, os, sys

CORPORA = {
    "schemas/cds/fixtures/issue": "valid-issue.json",
    "schemas/cds/fixtures/design": "valid-design.json",
    "schemas/cds/fixtures/subject": "valid-subject.json",
    "schemas/cds/fixtures/runinput": "valid-run-input.json",
}

def diffs(good, bad, path=""):
    """JSON locations at which bad departs from good."""
    if isinstance(good, dict) and isinstance(bad, dict):
        out, gone, added = [], [], []
        for k in good:
            if k not in bad:
                gone.append(k)
            else:
                out += diffs(good[k], bad[k], f"{path}.{k}")
        added = [k for k in bad if k not in good]
        # A key present under a differently-cased name, same value, is ONE
        # defect: the key was misspelled, not removed and something else added.
        for g in list(gone):
            for a in list(added):
                if g.lower() == a.lower() and good[g] == bad[a]:
                    gone.remove(g); added.remove(a)
                    out.append(f"{path}.{g}~{a}")
        out += [f"{path}.{k} (absent)" for k in gone]
        out += [f"{path}.{k} (extra)" for k in added]
        return out
    if isinstance(good, list) and isinstance(bad, list):
        # Same length: compare element-wise, so a single changed element is one
        # defect. Different length: the list itself is the one defect.
        if len(good) != len(bad):
            return [f"{path} (list)"]
        out = []
        for i, (g, b) in enumerate(zip(good, bad)):
            out += diffs(g, b, f"{path}[{i}]")
        return out
    return [] if good == bad else [path or "."]

bad = False
checked = 0
for d, positive in CORPORA.items():
    good = json.load(open(os.path.join(d, positive)))
    names = sorted(n for n in os.listdir(d)
                   if n.endswith(".json") and not n.startswith("valid-"))
    if not names:
        print(f"    {d}: no negatives found"); bad = True; continue
    for n in names:
        found = diffs(good, json.load(open(os.path.join(d, n))))
        checked += 1
        if len(found) != 1:
            print(f"    {d}/{n}: {len(found)} defects, want exactly 1: {found}")
            bad = True
# A run that compared nothing would report a corpus it never opened.
if checked == 0:
    print("    no negatives were compared"); bad = True
else:
    print(f"    compared {checked} negatives against their positives")
sys.exit(1 if bad else 0)
PYEOF
then
  echo "  ✓ every negative differs from its positive in exactly one place"
else
  echo "  ✗ a negative fixture carries more than one defect"; fail=1
fi

echo "# Go-only negatives (executable authority = the CLI)"
run_bad schemas/cdd/fixtures/invalid/cellspec-dup-required-id.json 'duplicate required_evidence id'
run_bad schemas/cdd/fixtures/invalid/cellspec-bad-producer.json 'producer'
run_bad schemas/cdd/fixtures/invalid/cellspec-missing-fill.json 'fill'
run_bad schemas/cdd/fixtures/invalid/cellspec-unknown-fill.json 'unknown alpha fill'
run_bad schemas/cdd/fixtures/invalid/cellspec-empty-goal.json 'goal'
run_bad schemas/cdd/fixtures/invalid/cellspec-case-alias.json 'unknown key'
run_bad schemas/cdd/fixtures/invalid/cellspec-bad-param-name.json 'parameter'
run_bad schemas/cdd/fixtures/invalid/cellspec-null-skills.json 'null'
run_bad schemas/cds/fixtures/invalid/cds-smuggled-argv.json 'argv' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-modelless-provider.json 'requires a model selector' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-fake-with-model.json 'takes no model' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-codex-held.json 'unknown provider' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-bad-hole-name.json 'malformed'
run_bad schemas/cds/fixtures/invalid/cds-bad-model-hole.json 'malformed'
run_bad schemas/cds/fixtures/invalid/cds-case-seat-tag.json 'seat declaration has no fill'
run_bad schemas/cds/fixtures/invalid/cds-case-top-arg.json 'unknown key' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-case-nested-arg.json 'unknown key' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-seat-skills.json 'unknown key "skills"' --input "$ri"
run_bad schemas/cds/fixtures/invalid/cds-assess-workspace.json 'unknown key "workspace"' --input "$ri"

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
  # ...against the FROZEN PRE-DELETION shape, deliberately. This episode ran
  # while the cds.patch fill still declared its own workspace; its record is
  # covered by a digest that recomputes, so editing the declaration out of it
  # would break that digest and turn evidence into a claim. The artifact stays
  # byte-for-byte and is held to the closed shape it was actually produced in.
  # A cell authored today cannot reach that shape: #CDSCellSpec.alpha is
  # #CDSPatchAlphaAuthored, which has no workspace key at all.
  if python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["alpha"], open(sys.argv[2],"w"))' "$ev" "$evalpha" 2>/dev/null; then
    vet_ok ./schemas/cds:cds "$evalpha" -d '#CDSPatchAlphaResolvedPreWorkspaceDeletion'
    # ...and it really is the pre-deletion shape, or the line above would be
    # vetting a current-shape declaration against a definition that merely
    # tolerates it.
    if "$CUE" vet ./schemas/cds:cds "$evalpha" -d '#CDSPatchAlphaResolved' >/dev/null 2>&1; then
      echo "  ✗ the evidence alpha vets the CURRENT shape: the frozen definition is guarding nothing"; fail=1
    else echo "  ✓ the evidence alpha is the pre-deletion shape the current definition rejects"; fi
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
# The CDS cell against a hermetic throwaway repository: the runtime cuts a
# worktree, the deterministic coder changes a file, the diff in the closure is
# MEASURED from that worktree, the reviewing seat reconstructs the candidate
# from (subject, matter), the closed checker runs against that reconstruction,
# and every catalogue unit is disposed of — and the episode still closes
# needs_repair (exit 1), because nothing was rented, so no acceptance criterion
# was decided by anyone. A cell that ran the whole chain and still refused to
# accept is the honest outcome, and it is part of the corpus.
# `$coderepo` and `$ri` were built near the top of this script: the CLI
# negatives above already need an admissible run input, because the subject is
# pinned before either seat is constructed.
cn_run cell run --contract "$spec" --input "$ri" \
  --param language=cnos.eng:eng/go --param provider=fake >"$tmp" 2>/dev/null
code=$?
if [ "$code" != 1 ]; then
  echo "  ✗ CDS cell exit=$code want=1 (needs_repair: nothing was rented, so nothing was decided)"; fail=1
else echo "  ✓ CDS cell closes needs_repair from an installed hub"; fi
if ! "$CUE" vet schemas/cdd/episode-closure.cue "$tmp" -d '#EpisodeClosure' >/dev/null 2>&1; then
  echo "  ✗ CDS closure failed #EpisodeClosure"; fail=1
else echo "  ✓ CDS closure vets #EpisodeClosure"; fi
decl="$tmpdir/resolved-alpha.json"
if python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["alpha"], open(sys.argv[2],"w"))' "$tmp" "$decl" 2>/dev/null &&
   "$CUE" vet ./schemas/cds:cds "$decl" -d '#CDSPatchAlphaResolved' >/dev/null 2>&1; then
  echo "  ✓ resolved alpha vets #CDSPatchAlphaResolved (canonical shape, methodology projection, no repository, no skills)"
else echo "  ✗ resolved alpha failed #CDSPatchAlphaResolved"; fail=1; fi
# ...and the reviewing seat records the ADVERSARIAL projection, so a record whose
# two seats were held to the same role — or to no methodology — is rejected by
# the closed shape rather than noticed by a reader.
bdecl="$tmpdir/resolved-beta.json"
if python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["beta"], open(sys.argv[2],"w"))' "$tmp" "$bdecl" 2>/dev/null &&
   "$CUE" vet ./schemas/cds:cds "$bdecl" -d '#CDSAssessBetaResolved' >/dev/null 2>&1; then
  echo "  ✓ resolved beta vets #CDSAssessBetaResolved (adversarial projection, no workspace, no skills)"
else echo "  ✗ resolved beta failed #CDSAssessBetaResolved"; fail=1; fi

# THE ASSESSMENT, LIVE. The catalogue is the admitted issue's acceptance ids
# plus the two check units, in that order, and every unit that did not pass
# states a reason. Asserted on the EMITTED CLOSURE rather than on the source,
# because what a reader can re-derive from the receipt is the whole point of
# putting the assessment in the record.
#
# The expected ids are read from the run input the episode actually consumed, so
# a corpus issue that grows a third criterion moves this check with it instead of
# failing it.
if python3 - "$tmp" "$ri" <<'PYEOF'
import json, sys
rec = json.load(open(sys.argv[1]))["receipt"]["record"]
issue = json.load(open(sys.argv[2]))["issue"]
want = [c["id"] for c in issue["acceptance"]] + ["check:matter-nonempty", "check:project-verify"]
units = rec["review"].get("assessment")
if units is None:
    print("    the record carries no assessment"); sys.exit(1)
got = [u["unit"] for u in units]
if got != want:
    print(f"    assessment covers {got}, want {want}"); sys.exit(1)
by = {u["unit"]: u for u in units}
# The matter unit was decided from a diff the runtime measured, so it passes;
# nothing was rented, so no acceptance criterion was decided by anyone.
if by["check:matter-nonempty"]["disposition"] != "pass":
    print(f"    the measured matter did not pass its unit: {by['check:matter-nonempty']}"); sys.exit(1)
for c in issue["acceptance"]:
    if by[c["id"]]["disposition"] != "unverified":
        print(f"    {c['id']} = {by[c['id']]['disposition']}; a seat that rented nothing decided nothing"); sys.exit(1)
for u in units:
    if u["disposition"] != "pass" and not u.get("reason", "").strip():
        print(f"    unit {u['unit']} is {u['disposition']} with no reason"); sys.exit(1)
# ...and the checker really ran against the reconstruction: its unit names the
# recipe, whatever it observed.
if "cnos.project-verify.v0" not in by["check:project-verify"].get("reason", ""):
    print(f"    the checker unit does not name the recipe: {by['check:project-verify']}"); sys.exit(1)
if rec["review"]["pass"]:
    print("    a review carrying non-passing units reported pass"); sys.exit(1)
PYEOF
then
  echo "  ✓ the closure carries an exact catalogue, each non-pass unit with a reason, and the checker's own unit"
else echo "  ✗ the live assessment is not the catalogue, or a judgement carries no reason"; fail=1; fi
# ONE REPOSITORY DECLARATION, live. The seat measured its base from the worktree
# it actually cut, and the contract's subject is the only place the run was told
# which repository and which commit that is — so these two must be the same
# string. While the fill resolved its own workspace they were two independent
# resolutions, and a closure recording a repository the episode never acted on
# still self-verified. Asserted on the emitted closure, not on the source.
if python3 - "$tmp" <<'PYEOF'
import json, re, sys
r = json.load(open(sys.argv[1]))["receipt"]["record"]
subject = json.loads(r["contract"]["subject"]) if isinstance(r["contract"]["subject"], str) else r["contract"]["subject"]
measured = [a["text"] for a in r["alpha"]["artifacts"] if a["id"] == "base_sha"]
if len(measured) != 1:
    print(f"    the record carries {len(measured)} measured base_sha artifacts, want exactly 1"); sys.exit(1)
if measured[0] != subject["base_sha"]:
    print(f"    measured base {measured[0]} != contract.subject.base_sha {subject['base_sha']}"); sys.exit(1)
# ...and the declaration names no repository at all, or "one source" would be
# a claim about which of two the seat happened to prefer.
alpha = r["resolved_spec"]["alpha"]
if "workspace" in alpha:
    print("    the resolved alpha still declares a workspace"); sys.exit(1)
# ONE METHODOLOGY, live. The seat declares no skills of its own and records the
# projection it was handed: the role, and the digest of the cell's one bundle.
# Asserted on the emitted closure, not on the source, for the same reason the
# base is.
if "skills" in alpha:
    print("    the resolved alpha still declares its own skills"); sys.exit(1)
m = alpha.get("methodology")
if not isinstance(m, dict) or m.get("role") != "constructive" or \
   not re.fullmatch(r"[0-9a-f]{64}", m.get("sha256", "")):
    print(f"    the resolved alpha records no methodology projection: {m!r}"); sys.exit(1)
PYEOF
then
  echo "  ✓ the measured base equals contract.subject.base_sha; the seat declares no repository and no skills, and records the cell's methodology"
else echo "  ✗ the record carries two repository declarations, or two methodologies, or they disagree"; fail=1; fi

echo "# run input: admitted at the door, pinned once, bound into the record"
# `$ri` was built above, before the cds.patch run that consumes it.

# A refused run input mints no episode. Asserted by REASON and by SHAPE, not by
# exit code alone: exit 4 is this refusal's own code, but a receipt on stdout
# and the ABSENCE of a closure are what show that nothing was run.
"$CN" cell run --contract schemas/cdd/fixtures/empty-cell-spec.json \
  --input schemas/cds/fixtures/runinput/runinput-malformed-issue.json \
  >"$tmpdir/refused.json" 2>"$tmpdir/refused.err"; refcode=$?
if [ "$refcode" != 4 ]; then
  echo "  ✗ a malformed issue must be refused at the door, got exit $refcode: $(head -c 200 "$tmpdir/refused.err")"; fail=1
elif ! grep -q '"outcome": "rejected"' "$tmpdir/refused.json" ||
     ! grep -q 'problem.diverges is required' "$tmpdir/refused.json" ||
     grep -q 'closure_schema' "$tmpdir/refused.json"; then
  echo "  ✗ a refusal must emit an admission receipt naming its reason and no closure: $(head -c 300 "$tmpdir/refused.json")"; fail=1
else echo "  ✓ a malformed issue is refused at the door with its own reason and no episode"; fi
# The receipt is vetted by BOTH authorities. It had a Go decoder and no CUE
# half until #CDSAdmissionReceipt, which is the single-authority state this
# schema pair exists to end — and `semantic_adequacy` is exactly the field that
# must not quietly stop being emitted, since it is the receipt's own statement
# that this door judged structure and not executability (Pi #81 C2).
vet_ok ./schemas/cds:cds "$tmpdir/refused.json" -d '#CDSAdmissionReceipt'

# AN ENVELOPE REFUSAL TAKES THE REFUSAL PATH. A wrong `kind` is decisively
# inadmissible, not a usage error: it must exit 4 with a receipt exactly as a
# malformed payload does. It exited 2 with empty stdout and no receipt while the
# runner decoded the envelope itself — one question answered through two paths,
# and the one an operator hit told them their file could not be read.
for env in bad-kind unknown-key mixed-case-key; do
  "$CN" cell run --contract schemas/cdd/fixtures/empty-cell-spec.json \
    --input "schemas/cds/fixtures/runinput/runinput-$env.json" \
    >"$tmpdir/env.json" 2>"$tmpdir/env.err"; envcode=$?
  if [ "$envcode" != 4 ]; then
    echo "  ✗ envelope refusal runinput-$env exited $envcode, want 4"; fail=1
  elif ! grep -q '"outcome": "rejected"' "$tmpdir/env.json" ||
       grep -q 'closure_schema' "$tmpdir/env.json"; then
    echo "  ✗ envelope refusal runinput-$env emitted no receipt: $(head -c 200 "$tmpdir/env.json")"; fail=1
  # O3: the receipt names WHICH document it refused. Nothing is frozen and no
  # closure exists, so this digest is the only record of the artifact decided on.
  elif ! python3 - "$tmpdir/env.json" "schemas/cds/fixtures/runinput/runinput-$env.json" <<'PYEOF'
import hashlib, json, sys
got = json.load(open(sys.argv[1]))["input_digest"]
want = hashlib.sha256(open(sys.argv[2], "rb").read()).hexdigest()
sys.exit(0 if got == want else 1)
PYEOF
  then
    echo "  ✗ envelope refusal runinput-$env carries the wrong input_digest"; fail=1
  else echo "  ✓ runinput-$env refuses through the receipt path, naming its input digest"; fi
done

# An absent payload is a DIFFERENT outcome from a malformed one. Two documents
# refused identically would make the vocabulary decorative.
"$CN" cell run --contract schemas/cdd/fixtures/empty-cell-spec.json \
  --input schemas/cds/fixtures/runinput/runinput-no-design.json \
  >"$tmpdir/incomplete.json" 2>/dev/null
if ! grep -q '"outcome": "incomplete"' "$tmpdir/incomplete.json"; then
  echo "  ✗ an absent design must be incomplete, not rejected: $(head -c 200 "$tmpdir/incomplete.json")"; fail=1
else echo "  ✓ an absent payload is incomplete where a malformed one is rejected"; fi

# The admitted path, and the live half of pinning: the run input was authored
# with `base_sha: HEAD`, a moving name, and the RECORDED contract subject must
# name a commit. This is the witness that pinning happened once, before the
# stations, rather than being a claim in a schema.
"$CN" cell run --contract schemas/cdd/fixtures/empty-cell-spec.json --input "$ri" \
  >"$tmpdir/bound.json" 2>/dev/null; bcode=$?
if [ "$bcode" != 3 ]; then
  echo "  ✗ an admitted stub run must close simulated (exit 3), got $bcode"; fail=1
else echo "  ✓ an admitted run input closes an episode"; fi
vet_ok schemas/cdd/episode-closure.cue "$tmpdir/bound.json" -d '#EpisodeClosure'
for slot in issue design subject; do
  if ! python3 -c 'import json,sys; json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["contract"][sys.argv[2]], open(sys.argv[3],"w"))' \
       "$tmpdir/bound.json" "$slot" "$tmpdir/bound-$slot.json" 2>/dev/null; then
    echo "  ✗ the record carries no contract $slot"; fail=1
  fi
done
vet_ok ./schemas/cds:cds "$tmpdir/bound-issue.json" -d '#CDSIssue'
vet_ok ./schemas/cds:cds "$tmpdir/bound-design.json" -d '#CDSDesign'
vet_ok ./schemas/cds:cds "$tmpdir/bound-subject.json" -d '#GitSnapshotPinned'
# ...and the authored document really did name a moving revision, or the line
# above would hold for an input that arrived pinned.
if ! grep -q '"base_sha": "HEAD"' "$ri"; then
  echo "  ✗ the live run input was not the moving-revision case"; fail=1
else echo "  ✓ authored HEAD reached the bound contract as an exact commit"; fi

if [ "$fail" = 0 ]; then echo "✓ cell schema/CLI corpus OK"; else echo "✗ cell schema check FAILED"; exit 1; fi
