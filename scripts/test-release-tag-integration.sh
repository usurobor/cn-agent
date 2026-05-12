#!/usr/bin/env bash
# test-release-tag-integration.sh — Test annotated tag creation in release script
#
# Creates a temporary repository to test that the release script properly
# generates annotated tags with structured tag messages.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== Release Script Tag Integration Test ==="

# Create test repository
TEST_DIR="/tmp/test-release-integration"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

echo "→ Setting up test repository..."
git init --quiet
git config user.name "Test Release"
git config user.email "release@test.com"

# Create VERSION file and initial commit
echo "1.0.0" > VERSION
echo "# Test Project" > README.md
git add VERSION README.md
git commit --quiet -m "Initial release"
git tag v1.0.0

# Create required release infrastructure files
mkdir -p scripts
cp "$REPO_ROOT/scripts/generate-release-tag-message.sh" scripts/
cp "$REPO_ROOT/scripts/stamp-versions.sh" scripts/
cp "$REPO_ROOT/scripts/check-version-consistency.sh" scripts/
cp "$REPO_ROOT/scripts/validate-release-gate.sh" scripts/

# Create a minimal release script for testing (just the tag creation part)
cat > scripts/release.sh << 'EOF'
#!/usr/bin/env bash
set -euo pipefail

VERSION="$1"
echo "=== Test Release: $VERSION ==="

# Generate tag message and create annotated tag
echo "→ generating tag message..."
TAG_MESSAGE_FILE=$(mktemp)
if ! scripts/generate-release-tag-message.sh "$VERSION" > "$TAG_MESSAGE_FILE"; then
  echo "ERROR: Failed to generate tag message" >&2
  rm -f "$TAG_MESSAGE_FILE"
  exit 1
fi

echo "→ creating annotated tag $VERSION..."
git tag -a "$VERSION" -F "$TAG_MESSAGE_FILE"
rm -f "$TAG_MESSAGE_FILE"
echo "→ tagged $VERSION (annotated)"

# Show tag details
echo "→ tag verification:"
echo "  Tag type: $(git cat-file -t "$VERSION")"
echo "  Tag message preview:"
git for-each-ref refs/tags/"$VERSION" --format='%(contents)' | head -3
EOF

chmod +x scripts/release.sh

# Create RELEASE.md to satisfy release gate
cat > RELEASE.md << 'EOF'
# Test Release 1.1.0

## Outcome
Test release to verify tag message generation.

## Changes
- #456: Test feature implementation

EOF

# Add a commit with issue reference
echo "Test change" >> README.md
git add README.md RELEASE.md
git commit -m "Add test feature (#456)"

# Test the release process
echo "→ Running test release..."
if scripts/release.sh 1.1.0; then
  echo "✅ Release script succeeded"
else
  echo "❌ Release script failed"
  exit 1
fi

# Verify the tag
echo ""
echo "→ Verifying tag properties..."

# Check tag exists
if git rev-parse 1.1.0 >/dev/null 2>&1; then
  echo "✅ Tag 1.1.0 exists"
else
  echo "❌ Tag 1.1.0 not found"
  exit 1
fi

# Check tag type
TAG_TYPE=$(git cat-file -t 1.1.0)
if [ "$TAG_TYPE" = "tag" ]; then
  echo "✅ Tag is annotated (type: $TAG_TYPE)"
else
  echo "❌ Tag is not annotated (type: $TAG_TYPE)"
  exit 1
fi

# Check tag message content
TAG_MESSAGE=$(git for-each-ref refs/tags/1.1.0 --format='%(contents)')
if echo "$TAG_MESSAGE" | grep -q "1.1.0 — Release Tag Message"; then
  echo "✅ Tag message contains expected header"
else
  echo "❌ Tag message missing expected header"
  exit 1
fi

if echo "$TAG_MESSAGE" | grep -q "#456:"; then
  echo "✅ Tag message contains issue reference"
else
  echo "❌ Tag message missing issue reference"
  exit 1
fi

echo ""
echo "✅ All tag integration tests passed!"
echo ""
echo "Generated tag message:"
echo "----------------------------------------"
echo "$TAG_MESSAGE"
echo "----------------------------------------"

cd - >/dev/null