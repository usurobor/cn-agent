#!/usr/bin/env bash
# release.sh — Single-command release: stamp, verify, commit, tag, push.
#
# Usage: scripts/release.sh [version]
#   If version is given, writes it to VERSION first.
#   If omitted, reads the current VERSION file.
#
# This is the only way to tag a release. Manual `git tag` is not allowed.
# See: operator/SKILL.md §3.4, CDD.md release gate.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

# 1. Set version if provided, otherwise read VERSION
if [ -n "${1:-}" ]; then
  echo "$1" > VERSION
  echo "→ VERSION set to $1"
fi

VERSION=$(cat VERSION | tr -d '\n')
if [ -z "$VERSION" ]; then
  echo "ERROR: VERSION is empty" >&2
  exit 1
fi

echo "=== Releasing $VERSION ==="

# 2. Ensure we're on main and up to date
BRANCH=$(git branch --show-current)
if [ "$BRANCH" != "main" ]; then
  echo "ERROR: must be on main (currently on $BRANCH)" >&2
  exit 1
fi

git fetch --quiet origin main
LOCAL=$(git rev-parse HEAD)
REMOTE=$(git rev-parse origin/main)
if [ "$LOCAL" != "$REMOTE" ]; then
  echo "ERROR: main is not up to date with origin/main" >&2
  echo "  local:  $LOCAL" >&2
  echo "  remote: $REMOTE" >&2
  exit 1
fi

# 3. Check tag doesn't already exist
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "ERROR: tag $VERSION already exists" >&2
  exit 1
fi

# 4. Stamp all manifests
echo "→ stamping manifests..."
scripts/stamp-versions.sh

# 5. Verify consistency
echo "→ verifying consistency..."
scripts/check-version-consistency.sh

# 6. Stage and commit if anything changed
if ! git diff --quiet; then
  git add -A
  git commit -m "release: $VERSION"
  echo "→ committed release artifacts"
else
  echo "→ no changes to commit (already stamped)"
fi

# 7. Update CHANGELOG if it exists and doesn't have this version
if [ -f CHANGELOG.md ]; then
  if ! grep -q "## $VERSION" CHANGELOG.md; then
    echo ""
    echo "WARNING: CHANGELOG.md does not contain a ## $VERSION entry." >&2
    echo "  Continue anyway? (y/N)"
    read -r answer
    if [ "$answer" != "y" ] && [ "$answer" != "Y" ]; then
      echo "Aborted. Update CHANGELOG.md and re-run." >&2
      exit 1
    fi
  fi
fi

# 8. Tag
git tag "$VERSION"
echo "→ tagged $VERSION"

# 9. Push
git push origin main --tags
echo ""
echo "✅ Released $VERSION — release CI will run."
echo "   Watch: gh run list --workflow release.yml --limit 1"
