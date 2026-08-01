#!/usr/bin/env bash
# verify-channel-reconstruction.sh — AC6 content-preservation verification (cnos#684)
#
# STATUS: SOURCE ONLY. NOT EXECUTED BY THIS CELL.
#
# This script is authored, not run, by α under cnos#684. It exists purely as a
# specified, reviewable procedure. It is NOT executed against `.cn-sigma/**`
# content by any role in this cell (δ/γ/α/β), for two independent reasons:
#
#   1. This checkout's sparse-checkout config excludes `.cn-sigma/` entirely
#      (`git sparse-checkout list` shows `!/.cn-sigma` with no counter-include).
#      `.cn-sigma` is not present in this working tree at all.
#   2. `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.12 (cnos#626) is a
#      standing doctrine constraint, independent of checkout visibility:
#      cell-execution cognition (δ/γ/α/β under wake-invoked-δ dispatch) does
#      not read or write `.cn-{agent}/logs/` regardless of what the working
#      tree happens to expose.
#
# The intended runner is the OPERATOR (or another actor with genuine
# `.cn-sigma/` access — a full, non-sparse checkout of `cnos` plus access to
# wherever the orphan-ref import (relocate step) lands, e.g. `cn-sigma` or a
# staging clone). Running this script is a PRECONDITION GATE that MUST pass
# before the strip step named in `dry-run-migration-plan.md` is proposed for
# operator execution — see that file's §"Sequencing" for where this fits.
#
# WHAT IT VERIFIES (AC6):
#   "Every `.cn-sigma/logs/**` payload currently in `main` history is
#    reconstructable from the orphan ref before any strip is proposed."
#
# APPROACH: Git blob SHAs are already content-addressed digests — two blobs
# with identical content always have the identical SHA, and vice versa. So
# "every payload is reconstructable" reduces to a set-membership check: does
# every blob SHA that has EVER existed at path `.cn-sigma/logs/**` in `main`'s
# history also exist as a reachable blob in the orphan-ref import? No re-hash
# scheme is needed; `git`'s own object model IS the content digest.
#
# Usage:
#   ./verify-channel-reconstruction.sh <path-to-main-checkout> <path-to-orphan-import-checkout-or-remote> <orphan-ref-name>
#
# Example (illustrative — NOT run by this cell):
#   ./verify-channel-reconstruction.sh \
#       /path/to/full/cnos-checkout \
#       /path/to/orphan-import-checkout \
#       refs/heads/channels/sigma/cnos-to-home
#
# Exit codes:
#   0 — PASS: every historical `.cn-sigma/logs/**` blob is present in the
#       orphan-ref import. Safe to proceed to the operator-gated strip step.
#   1 — FAIL: at least one historical blob is missing from the orphan-ref
#       import. Strip MUST NOT proceed until resolved.
#   2 — usage / environment error (bad args, not a git repo, ref not found).

set -euo pipefail

MAIN_CHECKOUT="${1:?usage: $0 <main-checkout> <orphan-import-checkout> <orphan-ref>}"
ORPHAN_CHECKOUT="${2:?usage: $0 <main-checkout> <orphan-import-checkout> <orphan-ref>}"
ORPHAN_REF="${3:?usage: $0 <main-checkout> <orphan-import-checkout> <orphan-ref>}"

TARGET_GLOB='.cn-sigma/logs/'
REPORT_DIR="$(mktemp -d)"
SOURCE_BLOBS="${REPORT_DIR}/source-blobs.txt"
ORPHAN_BLOBS="${REPORT_DIR}/orphan-blobs.txt"
MISSING_BLOBS="${REPORT_DIR}/missing-blobs.txt"

echo "== AC6 verification: .cn-sigma/logs/** reconstructability ==" >&2
echo "main checkout:   ${MAIN_CHECKOUT}" >&2
echo "orphan checkout: ${ORPHAN_CHECKOUT}" >&2
echo "orphan ref:      ${ORPHAN_REF}" >&2
echo "report dir:      ${REPORT_DIR}" >&2

# --- Step 1: enumerate every blob SHA that has ever existed at the target
# path across main's full history (every commit, every ancestor, not just
# HEAD — the append-only stream's full historical content is what AC6
# requires be reconstructable, not merely the current tip state).
#
# `git rev-list --objects --all -- <path>` walks every ref (restrict to
# `main` if other refs are noisy) and lists every object (commit, tree, blob)
# that is reachable AND touches the given pathspec, alongside its path.
# We then filter to blob lines only.
(
  cd "${MAIN_CHECKOUT}"
  git rev-list --objects main -- "${TARGET_GLOB}" \
    | git cat-file --batch-check='%(objectname) %(objecttype) %(rest)' \
    | awk '$2 == "blob" { print $1 }' \
    | sort -u > "${SOURCE_BLOBS}"
)
SOURCE_COUNT=$(wc -l < "${SOURCE_BLOBS}")
echo "source blobs (every historical .cn-sigma/logs/** blob in main): ${SOURCE_COUNT}" >&2

# --- Step 2: enumerate every blob SHA reachable from the orphan-ref import.
# After the relocate step (a `git filter-repo --path .cn-sigma/logs/
# --path-rename ...` style import — see dry-run-migration-plan.md §"Relocate
# step") lands content onto ORPHAN_REF, every blob that survived the import
# should be a blob reachable from that ref's history.
(
  cd "${ORPHAN_CHECKOUT}"
  git rev-list --objects "${ORPHAN_REF}" \
    | git cat-file --batch-check='%(objectname) %(objecttype) %(rest)' \
    | awk '$2 == "blob" { print $1 }' \
    | sort -u > "${ORPHAN_BLOBS}"
)
ORPHAN_COUNT=$(wc -l < "${ORPHAN_BLOBS}")
echo "orphan blobs (every blob reachable from ${ORPHAN_REF}): ${ORPHAN_COUNT}" >&2

# --- Step 3: set difference. Every source blob must appear in the orphan set.
# (The orphan set may be a strict superset — e.g. it may also carry blobs
# from renamed paths, README content, etc. — that is fine; the requirement is
# one-directional: source ⊆ orphan.)
comm -23 "${SOURCE_BLOBS}" "${ORPHAN_BLOBS}" > "${MISSING_BLOBS}" || true
MISSING_COUNT=$(wc -l < "${MISSING_BLOBS}")

echo "" >&2
echo "== Result ==" >&2
echo "missing blobs (present in main history, absent from orphan ref): ${MISSING_COUNT}" >&2

if [ "${MISSING_COUNT}" -gt 0 ]; then
  echo "" >&2
  echo "FAIL — the following blob SHAs exist in main's .cn-sigma/logs/** history" >&2
  echo "but are NOT reachable from ${ORPHAN_REF}. Do not proceed to the strip step." >&2
  echo "Missing blob list: ${MISSING_BLOBS}" >&2
  cat "${MISSING_BLOBS}" >&2
  exit 1
fi

echo "" >&2
echo "PASS — every historical .cn-sigma/logs/** blob (${SOURCE_COUNT} blobs) is" >&2
echo "reachable from ${ORPHAN_REF}. Content-preservation precondition satisfied." >&2
echo "" >&2
echo "-- Recommended secondary spot-check (manual, not automated by this script) --" >&2
echo "For a random sample of source commits, diff the file content directly:" >&2
echo "  git -C \"${MAIN_CHECKOUT}\" show <commit>:<path-under-.cn-sigma/logs/> \\" >&2
echo "    | diff - <(git -C \"${ORPHAN_CHECKOUT}\" show ${ORPHAN_REF}:<corresponding-path>)" >&2
echo "This confirms not just blob presence but that the import preserved the" >&2
echo "path-to-content mapping the migration plan intends (not merely that the" >&2
echo "bytes exist somewhere reachable from the ref)." >&2

exit 0
