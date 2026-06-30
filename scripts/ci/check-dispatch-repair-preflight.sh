#!/usr/bin/env bash
# check-dispatch-repair-preflight.sh — cnos#516 regression guard.
#
# Asserts the dispatch surface carries the repair re-entry preflight contract,
# so a re-claimed `status:changes → status:todo` cell cannot silently re-certify
# rejected work (the Pass 4D / cnos#514 failure mode). This is mechanically
# checkable: if the preflight is removed from the prompt or protocol — and thus
# from the re-rendered golden + live workflow — this check fails.
#
# Scope: presence-of-contract only. It does not run a wake; it guards that the
# rendered dispatch substrate still tells the body to load the repair context
# and forbids closeouts-without-repair_evidence on a repair re-entry.
#
# Exit 0 = contract present everywhere it must be; 1 = a surface lost it.

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
fail=0

# need <file> <label> <pattern>...
need() {
  local f="$1" label="$2"; shift 2
  if [ ! -f "$ROOT/$f" ]; then
    echo "::error::cnos#516 guard: required file missing: $f"; fail=1; return
  fi
  local pat
  for pat in "$@"; do
    if ! grep -qF -- "$pat" "$ROOT/$f"; then
      echo "::error::cnos#516 guard ($label): '$pat' missing from $f"; fail=1
    fi
  done
}

PROTO="src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md"
SKILL="src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md"
GOLDEN="src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml"
LIVE=".github/workflows/cnos-cds-dispatch.yml"

# Protocol surface defines the rule.
need "$PROTO" "protocol" \
  "Repair re-entry detection and preflight" \
  "cnos#516" \
  "REPAIR-PLAN" \
  "repair_evidence" \
  "run_class" \
  "first_pass" "repair_pass" "manual_delta_repair" "blocked"

# SKILL body + rendered substrate (golden and live workflow, which must be
# byte-identical per install-wake-golden) carry the executable preflight.
# The SKILL.md body is the wake's only prompt source, so this loop checks
# the SKILL.md, its golden, and the live workflow.
for f in "$SKILL" "$GOLDEN" "$LIVE"; do
  need "$f" "dispatch-surface" \
    "Repair re-entry preflight" \
    "REPAIR-PLAN" \
    "repair_evidence" \
    "run_class" \
    "first_pass" "repair_pass"
done

# The named failure mode must be explicitly forbidden (not merely implied):
# a closeout on a repair re-entry without a complete repair_evidence block.
need "$GOLDEN" "closeout-guard" "without a complete \`repair_evidence\` block"

if [ "$fail" -eq 0 ]; then
  echo "cnos#516 repair-preflight guard: dispatch surface carries the repair re-entry contract (protocol + prompt + golden + live workflow)."
fi
exit "$fail"
