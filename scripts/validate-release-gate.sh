#!/usr/bin/env bash
# validate-release-gate.sh — phase-aware CDD lifecycle artifact gate
#
# Modes:
#   pre-merge  Review/merge readiness. A triadic cycle must expose only the
#              artifacts available before merge: gamma-scaffold.md,
#              self-coherence.md, and beta-review.md.
#   post-merge Post-merge closeout after matter has merged. In addition to the
#              pre-merge record, alpha-closeout.md, beta-closeout.md, and
#              gamma-closeout.md are required. gamma-closeout.md must carry
#              the exact marker "CDD-Post-Merge-Closeout: complete".
#              This is a pre-release declaration, not terminal cycle closure;
#              RELEASE.md is not required at this phase.
#   release    Pre-tag/release gate (default). Requires RELEASE.md and the
#              complete post-merge lifecycle set and closeout marker.
#
# A scaffold's canonical **Mode:** declaration controls classification. For an
# exact --cycle target, the scaffold and one recognized declaration are
# mandatory: only explicit small-change/immediate-output may collapse, while
# substantial (including the landed `substantial-cycle` spelling used by CDS)
# takes the triadic path. Repository-wide legacy scans retain the
# historical inference for unscaffolded directories so old unrelated records
# do not suddenly block a release; scaffolded unknown modes still fail closed.
#
# Usage:
#   scripts/validate-release-gate.sh [--repo-root DIR] [--unreleased-dir DIR]
#     [--mode pre-merge|post-merge|release] [--cycle N]
#
# --cycle validates exactly one numeric cycle and fails if its directory is
# absent. Without --cycle, every directory under the unreleased root is read.
#
# Exit: 0=all conditions met, 1=invalid input or a missing/invalid artifact

set -euo pipefail

REPO_ROOT="."
UNRELEASED_DIR_OVERRIDE=""
MODE="release"
CYCLE=""
POST_MERGE_CLOSEOUT_MARKER="CDD-Post-Merge-Closeout: complete"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo-root)
      [[ $# -ge 2 ]] || { echo "Missing value for --repo-root" >&2; exit 1; }
      REPO_ROOT="$2"
      shift 2
      ;;
    --unreleased-dir)
      [[ $# -ge 2 ]] || { echo "Missing value for --unreleased-dir" >&2; exit 1; }
      UNRELEASED_DIR_OVERRIDE="$2"
      shift 2
      ;;
    --mode)
      [[ $# -ge 2 ]] || { echo "Missing value for --mode" >&2; exit 1; }
      MODE="$2"
      shift 2
      ;;
    --cycle)
      [[ $# -ge 2 ]] || { echo "Missing value for --cycle" >&2; exit 1; }
      CYCLE="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

case "$MODE" in
  release|pre-merge|post-merge) ;;
  *) echo "Unknown --mode: $MODE (expected: pre-merge, post-merge, or release)" >&2; exit 1 ;;
esac

if [[ -n "$CYCLE" && ! "$CYCLE" =~ ^[1-9][0-9]*$ ]]; then
  echo "Invalid --cycle: $CYCLE (expected a positive integer)" >&2
  exit 1
fi

UNRELEASED_DIR="${UNRELEASED_DIR_OVERRIDE:-$REPO_ROOT/.cdd/unreleased}"
ERRORS=0

if [[ "$MODE" == "release" ]]; then
  if [[ -f "$REPO_ROOT/RELEASE.md" ]]; then
    echo "  ✅ RELEASE.md present"
  else
    echo "  ❌ RELEASE.md missing at repo root — required before tag (CDD.md §1.2, release/SKILL.md §3.7)" >&2
    ERRORS=$((ERRORS + 1))
  fi
fi

REQUIRED_TRIADIC_PRE_MERGE=(gamma-scaffold.md self-coherence.md beta-review.md)
REQUIRED_TRIADIC_POST_MERGE=(gamma-scaffold.md self-coherence.md beta-review.md alpha-closeout.md beta-closeout.md gamma-closeout.md)

if [[ -n "$CYCLE" ]]; then
  CYCLE_DIR="$UNRELEASED_DIR/$CYCLE"
  if [[ ! -d "$CYCLE_DIR" ]]; then
    echo "  ❌ cycle $CYCLE: directory missing at $CYCLE_DIR" >&2
    ERRORS=$((ERRORS + 1))
    CYCLE_DIRS=()
  else
    CYCLE_DIRS=("$CYCLE_DIR/")
  fi
elif [[ -d "$UNRELEASED_DIR" ]]; then
  shopt -s nullglob
  CYCLE_DIRS=("$UNRELEASED_DIR"/*/)
  shopt -u nullglob
else
  CYCLE_DIRS=()
fi

for dir in "${CYCLE_DIRS[@]}"; do
  [[ -d "$dir" ]] || continue
  N="$(basename "$dir")"

  SCAFFOLD_MODE=""
  if [[ -f "$dir/gamma-scaffold.md" ]]; then
    SCAFFOLD_MODE="$(sed -nE 's/^\*\*Mode:\*\*[[:space:]]*(substantial-cycle|substantial|small-change|immediate-output)([^[:alnum:]_-].*)?$/\1/p' "$dir/gamma-scaffold.md")"
    SCAFFOLD_MODE="${SCAFFOLD_MODE%%$'\n'*}"
    # CDS names the intervention class `substantial-cycle`; γ's work-shape
    # shorthand is `substantial`. They are one triadic class at this gate.
    [[ "$SCAFFOLD_MODE" == "substantial-cycle" ]] && SCAFFOLD_MODE="substantial"
  fi

  if [[ -n "$CYCLE" && ! -f "$dir/gamma-scaffold.md" ]]; then
    echo "  ❌ cycle $N: gamma-scaffold.md missing — exact-cycle validation requires an explicit canonical **Mode:** declaration" >&2
    ERRORS=$((ERRORS + 1))
    continue
  fi

  if [[ -f "$dir/gamma-scaffold.md" && -z "$SCAFFOLD_MODE" ]]; then
    echo "  ❌ cycle $N: gamma-scaffold.md lacks a recognized canonical **Mode:** declaration (expected substantial/substantial-cycle, small-change, or immediate-output)" >&2
    ERRORS=$((ERRORS + 1))
    continue
  fi

  if [[ -f "$dir/beta-review.md" || "$SCAFFOLD_MODE" == "substantial" ]]; then
    CYCLE_ERRORS=0
    if [[ "$MODE" == "pre-merge" ]]; then
      REQUIRED=("${REQUIRED_TRIADIC_PRE_MERGE[@]}")
      PHASE_LABEL="before merge"
    else
      REQUIRED=("${REQUIRED_TRIADIC_POST_MERGE[@]}")
      PHASE_LABEL="for post-merge closeout"
    fi

    for f in "${REQUIRED[@]}"; do
      if [[ ! -f "$dir/$f" ]]; then
        echo "  ❌ cycle $N: missing $f — required $PHASE_LABEL" >&2
        ERRORS=$((ERRORS + 1))
        CYCLE_ERRORS=$((CYCLE_ERRORS + 1))
      fi
    done

    if [[ "$MODE" != "pre-merge" && -f "$dir/gamma-closeout.md" ]]; then
      if ! grep -Fxq "$POST_MERGE_CLOSEOUT_MARKER" "$dir/gamma-closeout.md"; then
        echo "  ❌ cycle $N: gamma-closeout.md lacks post-merge closeout marker: $POST_MERGE_CLOSEOUT_MARKER" >&2
        ERRORS=$((ERRORS + 1))
        CYCLE_ERRORS=$((CYCLE_ERRORS + 1))
      fi
    fi

    if [[ $CYCLE_ERRORS -eq 0 ]]; then
      echo "  ✅ cycle $N (triadic): $MODE artifact set complete"
    fi
  elif [[ "$SCAFFOLD_MODE" == "small-change" || "$SCAFFOLD_MODE" == "immediate-output" || -z "$CYCLE" ]]; then
    echo "  ✅ cycle $N (small-change): no triadic lifecycle artifacts required"
  else
    echo "  ❌ cycle $N: internal classification error for mode '$SCAFFOLD_MODE'" >&2
    ERRORS=$((ERRORS + 1))
  fi
done

echo ""
if [[ $ERRORS -gt 0 ]]; then
  case "$MODE" in
    pre-merge)  echo "❌ Pre-merge gate FAILED: $ERRORS artifact error(s) — merge blocked" >&2 ;;
    post-merge) echo "❌ Post-merge closeout gate FAILED: $ERRORS artifact error(s) — release remains blocked" >&2 ;;
    release)    echo "❌ Release gate FAILED: $ERRORS artifact error(s)" >&2 ;;
  esac
  exit 1
fi

case "$MODE" in
  pre-merge)  echo "✅ Pre-merge gate passed" ;;
  post-merge) echo "✅ Post-merge closeout gate passed (release pending)" ;;
  release)    echo "✅ Release gate passed" ;;
esac
