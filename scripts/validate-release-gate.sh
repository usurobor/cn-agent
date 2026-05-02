#!/usr/bin/env bash
# validate-release-gate.sh — pre-tag release gate for scripts/release.sh
#
# Validates two conditions before a release tag may be cut:
#   1. RELEASE.md exists at repo root  (CDD.md §1.2, §5.3b, release/SKILL.md §3.7)
#   2. Each substantial cycle dir in .cdd/unreleased/ has required lifecycle artifacts
#
# Cycle classification follows CDD.md §1.2 small-change collapse table:
#   Triadic/substantial — beta-review.md present in the cycle dir:
#     required: self-coherence.md, beta-review.md, alpha-closeout.md,
#               beta-closeout.md, gamma-closeout.md  (CDD.md §5.3b)
#   Small-change — no beta-review.md present:
#     no cycle-dir artifacts required (all collapsed per CDD.md §1.2)
#
# Usage:
#   scripts/validate-release-gate.sh [--repo-root DIR] [--unreleased-dir DIR]
#
# Exit: 0=all conditions met, 1=one or more required artifacts missing

set -euo pipefail

REPO_ROOT="."
UNRELEASED_DIR_OVERRIDE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo-root)      REPO_ROOT="$2";          shift 2 ;;
    --unreleased-dir) UNRELEASED_DIR_OVERRIDE="$2"; shift 2 ;;
    *) echo "Unknown option: $1" >&2; exit 1 ;;
  esac
done

UNRELEASED_DIR="${UNRELEASED_DIR_OVERRIDE:-$REPO_ROOT/.cdd/unreleased}"

ERRORS=0

# 1. RELEASE.md must exist at repo root before tag.
if [[ -f "$REPO_ROOT/RELEASE.md" ]]; then
  echo "  ✅ RELEASE.md present"
else
  echo "  ❌ RELEASE.md missing at repo root — required before tag (CDD.md §1.2, release/SKILL.md §3.7)" >&2
  ERRORS=$((ERRORS + 1))
fi

# 2. Validate cycle dirs in .cdd/unreleased/.
REQUIRED_TRIADIC=(self-coherence.md beta-review.md alpha-closeout.md beta-closeout.md gamma-closeout.md)

if [[ -d "$UNRELEASED_DIR" ]]; then
  shopt -s nullglob
  for dir in "$UNRELEASED_DIR"/*/; do
    [[ -d "$dir" ]] || continue
    N="$(basename "$dir")"
    if [[ -f "$dir/beta-review.md" ]]; then
      # Triadic/substantial cycle: all 5 lifecycle files required (CDD.md §5.3b).
      CYCLE_ERRORS=0
      for f in "${REQUIRED_TRIADIC[@]}"; do
        if [[ ! -f "$dir/$f" ]]; then
          echo "  ❌ cycle $N (triadic): missing $f" >&2
          ERRORS=$((ERRORS + 1))
          CYCLE_ERRORS=$((CYCLE_ERRORS + 1))
        fi
      done
      if [[ $CYCLE_ERRORS -eq 0 ]]; then
        echo "  ✅ cycle $N (triadic): all required artifacts present"
      fi
    else
      # Small-change cycle: artifact requirements collapse per CDD.md §1.2.
      echo "  ✅ cycle $N (small-change): no required cycle-dir artifacts (CDD.md §1.2)"
    fi
  done
  shopt -u nullglob
fi

# Summary
echo ""
if [[ $ERRORS -gt 0 ]]; then
  echo "❌ Release gate FAILED: $ERRORS missing required artifact(s)" >&2
  exit 1
fi

echo "✅ Release gate passed"
exit 0
