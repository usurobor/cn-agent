#!/usr/bin/env bash
# check-dispatch-closeout-integrity.sh — cnos#524 W4 RCA regression guard.
#
# Two jobs:
#  (1) presence-of-contract: assert the rendered dispatch surface carries the
#      closeout-integrity contract (deliverable proof before status:review), so
#      a future edit cannot silently remove it from the prompt/protocol (and
#      thus from the re-rendered golden + live workflow).
#  (2) --self-test: exercise the empty-review detector — the mechanical failure
#      mode from cnos#524 W4 (a cell set status:review with no PR and no commits).
#      Confirms the detector flags the bad case and passes the good cases.
#
# Background: in cnos#524 W4 a dispatch run claimed the issue, ran ~25 min, and
# set status:review while pushing NO PR, NO commits, NO closeout, NO stop comment
# — a false-complete indistinguishable from a finished cell. dispatch-protocol
# §2.9 makes "status:review without a deliverable" a named violation; this guard
# keeps that contract present and the detector honest.
#
# Exit 0 = contract present + self-test passes; 1 = a surface lost it / detector broke.

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
fail=0

# --- the empty-review detector (the mechanical core of §2.9) -----------------
# A status:review transition is a cnos#524 violation when no deliverable backs
# it: status == "review" AND (no PR OR no commits-beyond-base).
# Echoes "violation" or "ok"; pure function, no I/O.
closeout_violation() {
  status="$1" has_pr="$2" has_commits="$3"
  if [ "$status" = "review" ] && { [ "$has_pr" != "yes" ] || [ "$has_commits" != "yes" ]; }; then
    echo "violation"
  else
    echo "ok"
  fi
}

self_test() {
  local rc=0
  # case: empty review (W4 failure mode) -> MUST be a violation
  [ "$(closeout_violation review no no)" = "violation" ] || { echo "::error::self-test: empty review (review/no-pr/no-commits) not flagged"; rc=1; }
  # case: review with PR but no commits -> violation
  [ "$(closeout_violation review yes no)" = "violation" ] || { echo "::error::self-test: review with PR but no commits not flagged"; rc=1; }
  # case: review with PR + commits -> ok (a real deliverable)
  [ "$(closeout_violation review yes yes)" = "ok" ] || { echo "::error::self-test: real deliverable (review/pr/commits) wrongly flagged"; rc=1; }
  # case: in-progress with nothing -> ok (not a review transition)
  [ "$(closeout_violation in-progress no no)" = "ok" ] || { echo "::error::self-test: non-review status wrongly flagged"; rc=1; }
  if [ "$rc" -eq 0 ]; then
    echo "cnos#524 closeout-integrity self-test: empty-review detector correct (flags review-without-deliverable; passes real deliverables)."
  fi
  return "$rc"
}

if [ "${1:-}" = "--self-test" ]; then
  self_test; exit $?
fi

# --- presence-of-contract guard ---------------------------------------------
need() {
  local f="$1" label="$2"; shift 2
  if [ ! -f "$ROOT/$f" ]; then
    echo "::error::cnos#524 guard: required file missing: $f"; fail=1; return
  fi
  local pat
  for pat in "$@"; do
    if ! grep -qF -- "$pat" "$ROOT/$f"; then
      echo "::error::cnos#524 guard ($label): '$pat' missing from $f"; fail=1
    fi
  done
}

PROTO="src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md"
PROMPT="src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md"
SKILL="src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md"
GOLDEN="src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml"
LIVE=".github/workflows/cnos-cds-dispatch.yml"

# Protocol surface defines the rule (§2.9 + §4.9 + D12).
need "$PROTO" "protocol" \
  "Closeout integrity" \
  "deliverable proof before" \
  "deliverable_evidence" \
  "cnos#524" \
  "check-dispatch-closeout-integrity.sh"

# Prompt + SKILL body + rendered substrate (golden + live workflow) carry the
# executable preflight. SKILL.md is the W3 default render source; prompt.md is
# kept verbatim-identical for W2/W3 parity until W4.
for f in "$PROMPT" "$SKILL" "$GOLDEN" "$LIVE"; do
  need "$f" "dispatch-surface" \
    "Closeout integrity preflight" \
    "deliverable_evidence" \
    "commits beyond its base" \
    "No-deliverable rule" \
    "status:review"
done

# The named failure mode must be explicitly forbidden, not merely implied.
need "$GOLDEN" "empty-review-guard" "without a complete \`deliverable_evidence\` block"

# Always run the detector self-test as part of the presence guard too.
self_test || fail=1

if [ "$fail" -eq 0 ]; then
  echo "cnos#524 closeout-integrity guard: dispatch surface carries the deliverable-proof contract (protocol + prompt + SKILL + golden + live), and the empty-review detector is correct."
fi
exit "$fail"
