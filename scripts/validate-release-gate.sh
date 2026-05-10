#!/usr/bin/env bash
# validate-release-gate.sh — pre-tag release gate and pre-merge closure gate
#
# Validates two modes:
#
#   --mode release (default): pre-tag conditions before a release tag may be cut
#     1. RELEASE.md exists at repo root  (CDD.md §1.2, §5.3b, release/SKILL.md §3.7)
#     2. Each substantial cycle dir in .cdd/unreleased/ has required lifecycle artifacts
#        (CDD.md §5.3b: self-coherence.md, beta-review.md, alpha-closeout.md,
#         beta-closeout.md, gamma-closeout.md for triadic cycles)
#
#   --mode pre-merge: pre-merge closure gate (CDD.md §5.3b, gamma/SKILL.md §2.10)
#     For each cycle dir in .cdd/unreleased/ that has beta-review.md present:
#       required: alpha-closeout.md, beta-closeout.md, gamma-closeout.md
#     RELEASE.md check skipped (release-time only, not pre-merge).
#     Exit non-zero with diagnostic naming each missing file and cycle number.
#
# Cycle classification follows CDD.md §1.2 small-change collapse table:
#   Triadic/substantial — beta-review.md present in the cycle dir
#   Small-change — no beta-review.md present (artifact requirements collapse)
#
# Usage:
#   scripts/validate-release-gate.sh [--repo-root DIR] [--unreleased-dir DIR] [--mode MODE]
#
# Exit: 0=all conditions met, 1=one or more required artifacts missing

set -euo pipefail

REPO_ROOT="."
UNRELEASED_DIR_OVERRIDE=""
MODE="release"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo-root)      REPO_ROOT="$2";               shift 2 ;;
    --unreleased-dir) UNRELEASED_DIR_OVERRIDE="$2"; shift 2 ;;
    --mode)           MODE="$2";                    shift 2 ;;
    *) echo "Unknown option: $1" >&2; exit 1 ;;
  esac
done

case "$MODE" in
  release|pre-merge) ;;
  *) echo "Unknown --mode: $MODE (expected: release or pre-merge)" >&2; exit 1 ;;
esac

UNRELEASED_DIR="${UNRELEASED_DIR_OVERRIDE:-$REPO_ROOT/.cdd/unreleased}"

ERRORS=0

# ── Release-mode-only checks ──────────────────────────────────────────────────

if [[ "$MODE" == "release" ]]; then
  # 1. RELEASE.md must exist at repo root before tag.
  if [[ -f "$REPO_ROOT/RELEASE.md" ]]; then
    echo "  ✅ RELEASE.md present"
  else
    echo "  ❌ RELEASE.md missing at repo root — required before tag (CDD.md §1.2, release/SKILL.md §3.7)" >&2
    ERRORS=$((ERRORS + 1))
  fi
fi

# ── Cycle-dir checks (both modes) ─────────────────────────────────────────────

# release mode: all 5 lifecycle files required for triadic cycles (CDD.md §5.3b)
REQUIRED_TRIADIC_RELEASE=(self-coherence.md beta-review.md alpha-closeout.md beta-closeout.md gamma-closeout.md)

# pre-merge mode: close-out triple required for triadic cycles (CDD.md §5.3b, gamma/SKILL.md §2.10)
REQUIRED_TRIADIC_PRE_MERGE=(alpha-closeout.md beta-closeout.md gamma-closeout.md)

if [[ -d "$UNRELEASED_DIR" ]]; then
  shopt -s nullglob
  for dir in "$UNRELEASED_DIR"/*/; do
    [[ -d "$dir" ]] || continue
    N="$(basename "$dir")"
    if [[ -f "$dir/beta-review.md" ]]; then
      # Triadic/substantial cycle: required artifact set depends on mode.
      CYCLE_ERRORS=0
      if [[ "$MODE" == "release" ]]; then
        REQUIRED=("${REQUIRED_TRIADIC_RELEASE[@]}")
      else
        REQUIRED=("${REQUIRED_TRIADIC_PRE_MERGE[@]}")
      fi
      for f in "${REQUIRED[@]}"; do
        if [[ ! -f "$dir/$f" ]]; then
          echo "  ❌ cycle $N: missing $f — required before merge (CDD.md §5.3b, gamma/SKILL.md §2.10)" >&2
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

# ── Summary ───────────────────────────────────────────────────────────────────

echo ""
if [[ $ERRORS -gt 0 ]]; then
  if [[ "$MODE" == "pre-merge" ]]; then
    echo "❌ Pre-merge closure gate FAILED: $ERRORS missing required artifact(s) — merge blocked (release/SKILL.md §3.8, CDD.md §5.3b)" >&2
  else
    echo "❌ Release gate FAILED: $ERRORS missing required artifact(s)" >&2
  fi
  exit 1
fi

if [[ "$MODE" == "pre-merge" ]]; then
  echo "✅ Pre-merge closure gate passed"
else
  echo "✅ Release gate passed"
fi
exit 0
