#!/usr/bin/env bash
# generate-release-tag-message.sh — Generate structured annotated tag messages
#
# Usage: scripts/generate-release-tag-message.sh <version> [previous-tag]
#   version: Target version for the tag (e.g., "3.69.0")
#   previous-tag: Previous tag to compare from (defaults to latest tag)
#
# Outputs plain text suitable for `git tag -a -F <file>` or equivalent.
# Discovers issues from commits, enriches with GitHub metadata when available,
# includes CDD review artifacts when present, and handles wave context.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

# Parse arguments
VERSION="${1:-}"
PREVIOUS_TAG="${2:-}"

if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version> [previous-tag]" >&2
  echo "  version: Target version for the tag (e.g., '3.69.0')" >&2
  echo "  previous-tag: Previous tag to compare from (defaults to latest tag)" >&2
  exit 1
fi

# Find previous tag if not provided
if [ -z "$PREVIOUS_TAG" ]; then
  PREVIOUS_TAG=$(git tag --sort=-version:refname | head -1)
  if [ -z "$PREVIOUS_TAG" ]; then
    echo "ERROR: No previous tag found and none specified" >&2
    exit 1
  fi
fi

echo "# Generating tag message for $VERSION (since $PREVIOUS_TAG)" >&2

# Get commit range
COMMIT_RANGE="${PREVIOUS_TAG}..HEAD"

# Function to extract issue numbers from commit messages
extract_issues() {
  git log --oneline "$COMMIT_RANGE" | \
    grep -oE '(#[0-9]+|[αβγδ]-[0-9]+)' | \
    sed -E 's/^[αβγδ]-//; s/#//' | \
    sort -u
}

# Function to get GitHub issue metadata with fallback
get_issue_metadata() {
  local issue_num="$1"

  # Try GitHub API via gh command
  if command -v gh >/dev/null 2>&1; then
    if gh issue view "$issue_num" --json title,labels,state 2>/dev/null; then
      return 0
    fi
  fi

  # Fallback: return minimal structure
  echo '{"title":"unavailable","labels":[],"state":"unknown"}'
}

# Function to find wave context
find_wave_context() {
  if [ ! -d ".cdd/waves" ]; then
    return 1
  fi

  # Look for most recent wave directory
  local wave_dir
  wave_dir=$(find .cdd/waves -maxdepth 1 -type d -name '*-*' | sort | tail -1)

  if [ -z "$wave_dir" ] || [ ! -f "$wave_dir/manifest.md" ]; then
    return 1
  fi

  echo "$wave_dir"
}

# Function to get CDD review artifacts for an issue
get_review_artifacts() {
  local issue_num="$1"
  local rounds=0
  local findings=""

  # Check in both unreleased and releases directories
  for base_dir in .cdd/unreleased .cdd/releases/*; do
    local issue_dir="$base_dir/$issue_num"

    if [ -d "$issue_dir" ]; then
      # Count review rounds from beta-review.md
      if [ -f "$issue_dir/beta-review.md" ]; then
        rounds=$(grep -c "^## Round" "$issue_dir/beta-review.md" 2>/dev/null || echo "1")
      fi

      # Extract findings from beta-closeout.md
      if [ -f "$issue_dir/beta-closeout.md" ]; then
        findings=$(grep -A 5 -i "findings\|highlights" "$issue_dir/beta-closeout.md" | head -3 | tail -2 | tr '\n' ' ' 2>/dev/null || echo "")
      fi

      break
    fi
  done

  echo "$rounds|$findings"
}

# Function to get wave info if available
get_wave_info() {
  local wave_dir
  wave_dir=$(find_wave_context)

  if [ -z "$wave_dir" ]; then
    echo "unavailable"
    return
  fi

  local wave_name=""
  local issue_count=0

  if [ -f "$wave_dir/manifest.md" ]; then
    wave_name=$(basename "$wave_dir")
    issue_count=$(grep -c '^- #[0-9]' "$wave_dir/manifest.md" 2>/dev/null || echo "0")
  fi

  echo "$wave_name|$issue_count"
}

# Generate timestamp
TIMESTAMP=$(date -u +"%Y-%m-%d %H:%M:%S UTC")

# Start generating tag message
echo "$VERSION — Release Tag Message"
echo ""

# Get wave context
WAVE_INFO=$(get_wave_info)
if [ "$WAVE_INFO" != "unavailable" ]; then
  WAVE_NAME=$(echo "$WAVE_INFO" | cut -d'|' -f1)
  WAVE_COUNT=$(echo "$WAVE_INFO" | cut -d'|' -f2)
  echo "Wave: $WAVE_NAME ($WAVE_COUNT issues)"
  echo ""
fi

# Extract and process issues
ISSUES=$(extract_issues)

if [ -z "$ISSUES" ] || [ "$ISSUES" = "" ]; then
  echo "No issues found in commit range $COMMIT_RANGE"
  echo ""
  echo "Generated: $TIMESTAMP"
  exit 0
fi

ISSUE_COUNT=$(echo "$ISSUES" | grep -c '^[0-9]' || echo "0")

echo "Changes:"

# Process each issue
for issue in $ISSUES; do
  echo "# Processing issue #$issue..." >&2

  # Get GitHub metadata
  METADATA=$(get_issue_metadata "$issue")

  # Parse JSON safely
  if command -v jq >/dev/null 2>&1; then
    TITLE=$(echo "$METADATA" | jq -r '.title // "unavailable"' 2>/dev/null || echo "unavailable")
    LABELS=$(echo "$METADATA" | jq -r '.labels[]?.name // empty' 2>/dev/null | tr '\n' ',' | sed 's/,$//' || echo "")
  else
    # Fallback without jq
    TITLE="unavailable"
    LABELS=""
  fi

  # Get review artifacts
  REVIEW_INFO=$(get_review_artifacts "$issue")
  ROUNDS=$(echo "$REVIEW_INFO" | cut -d'|' -f1)
  FINDINGS=$(echo "$REVIEW_INFO" | cut -d'|' -f2)

  # Format issue entry
  echo -n "- #$issue: $TITLE"

  # Add labels if available
  if [ -n "$LABELS" ] && [ "$LABELS" != "" ]; then
    echo -n " [$LABELS]"
  fi

  # Add review info if available
  if [ "$ROUNDS" -gt 0 ]; then
    echo -n " ($ROUNDS rounds"
    if [ -n "$FINDINGS" ] && [ "$FINDINGS" != "" ]; then
      echo -n ", $FINDINGS"
    fi
    echo -n ")"
  fi

  echo ""
done

echo ""
echo "Summary: $ISSUE_COUNT issues since $PREVIOUS_TAG"
echo ""
echo "Generated: $TIMESTAMP"