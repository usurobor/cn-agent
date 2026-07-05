#!/usr/bin/env bash
# check-dispatch-closeout-integrity.sh — cnos#524 W4 RCA regression guard.
#
# One job: presence-of-contract. Assert the rendered dispatch surface carries
# the closeout-integrity contract (deliverable proof before status:review), so
# a future edit cannot silently remove it from the prompt/protocol (and thus
# from the re-rendered golden + live workflow).
#
# Background: in cnos#524 W4 a dispatch run claimed the issue, ran ~25 min, and
# set status:review while pushing NO PR, NO commits, NO closeout, NO stop comment
# — a false-complete indistinguishable from a finished cell. dispatch-protocol
# §2.9 makes "status:review without a deliverable" a named violation; this guard
# keeps that contract present in the prompt/protocol/golden/live-workflow text.
#
# cnos#600 consolidation note: this script previously also carried a
# `--self-test` mode exercising a hand-rolled `closeout_violation(status,
# has_pr, has_commits)` shell predicate against 4 cases. That was a duplicate,
# shell-reimplemented copy of an invariant now proven directly against the
# REAL FSM transition table (not a reimplemented shim) by
# src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go's
# TestAC3_EmptyReviewBlocked (empty review -> blocked),
# TestApply_EmptyReviewStateBlocked (same, through the --apply CLI path),
# TestAC574_ReviewPartialEvidenceBlocked (PR present, no commits -> blocked),
# and TestAC574_ReviewWithPRStillValid (PR + commits -> valid) — all green
# as of this consolidation (`go test ./src/packages/cnos.issues/commands/issues-fsm/...`).
# The bash self-test was folded out in favor of those tests; this script now
# only proves the PROMPT still says the words, not that the FSM behaves
# correctly — the FSM's actual behavior is proven by the Go tests above and by
# transitions.json's `review_request_present && pr_exists && pr_has_commits`
# all_true guard (cnos#574 AC2/AC3), not by this script.
#
# Exit 0 = contract present everywhere it must be; 1 = a surface lost it.

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
fail=0

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

# SKILL body + rendered substrate (golden + live workflow) carry the
# executable preflight. The SKILL.md body is the wake's only prompt source,
# so this loop checks the SKILL.md, its golden, and the live workflow.
for f in "$SKILL" "$GOLDEN" "$LIVE"; do
  need "$f" "dispatch-surface" \
    "Closeout integrity preflight" \
    "deliverable_evidence" \
    "commits beyond its base" \
    "No-deliverable rule" \
    "status:review"
done

# The named failure mode must be explicitly forbidden, not merely implied.
need "$GOLDEN" "empty-review-guard" "without a complete \`deliverable_evidence\` block"

if [ "$fail" -eq 0 ]; then
  echo "cnos#524 closeout-integrity guard: dispatch surface carries the deliverable-proof contract (protocol + prompt + SKILL + golden + live). The empty-review detector itself is proven live by the Go FSM test suite (see header note) rather than by a bash self-test."
fi
exit "$fail"
