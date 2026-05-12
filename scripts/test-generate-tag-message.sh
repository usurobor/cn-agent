#!/usr/bin/env bash
# test-generate-tag-message.sh — Test suite for generate-release-tag-message.sh
#
# Creates temporary test repositories with controlled git history and fixtures
# to verify tag message generation under various scenarios.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Test counter
TESTS_RUN=0
TESTS_PASSED=0

# Test result tracking
test_result() {
  local name="$1"
  local expected="$2"
  local actual="$3"

  TESTS_RUN=$((TESTS_RUN + 1))

  if [ "$actual" = "$expected" ]; then
    echo "✅ $name"
    TESTS_PASSED=$((TESTS_PASSED + 1))
  else
    echo "❌ $name"
    echo "  Expected: $expected"
    echo "  Actual:   $actual"
  fi
}

# Create test repository
setup_test_repo() {
  local test_dir="$1"
  rm -rf "$test_dir"
  mkdir -p "$test_dir"
  cd "$test_dir"

  git init --quiet
  git config user.name "Test User"
  git config user.email "test@example.com"

  # Create initial commit and tag
  echo "initial" > README.md
  git add README.md
  git commit --quiet -m "Initial commit"
  git tag v1.0.0

  # Copy generator script
  cp "$REPO_ROOT/scripts/generate-release-tag-message.sh" ./generate-tag-message.sh
  chmod +x ./generate-tag-message.sh
}

# Test 1: No issues in commit range
test_no_issues() {
  echo "Testing: No issues in commit range"
  local test_dir="/tmp/test-tag-gen-no-issues"

  setup_test_repo "$test_dir"

  # Add commit without issue reference
  echo "change" >> README.md
  git add README.md
  git commit --quiet -m "Update README"

  # Test
  local output
  output=$(bash ./generate-tag-message.sh 1.1.0 v1.0.0 2>/dev/null || echo "ERROR")

  if [ "$output" = "ERROR" ]; then
    echo "❌ Generator script failed in test_no_issues"
    cd - >/dev/null
    return
  fi

  local has_no_issues
  has_no_issues=$(echo "$output" | grep -c "No issues found" || echo "0")

  test_result "No issues found message" "1" "$has_no_issues"

  cd - >/dev/null
}

# Test 2: Single issue with #N format
test_single_issue_hash_format() {
  echo "Testing: Single issue with #N format"
  local test_dir="/tmp/test-tag-gen-hash"

  setup_test_repo "$test_dir"

  # Add commit with #N issue reference
  echo "fix" >> README.md
  git add README.md
  git commit --quiet -m "Fix bug (#123)"

  # Test
  local output
  output=$(bash ./generate-tag-message.sh 1.1.0 v1.0.0)
  local has_issue_123
  has_issue_123=$(echo "$output" | grep -c "#123:" || echo "0")

  test_result "Issue #123 extracted" "1" "$has_issue_123"

  cd - >/dev/null
}

# Test 3: Multiple issues with role prefix format
test_multiple_issues_role_format() {
  echo "Testing: Multiple issues with role prefix format"
  local test_dir="/tmp/test-tag-gen-role"

  setup_test_repo "$test_dir"

  # Add commits with role-prefixed issue references
  echo "change1" >> README.md
  git add README.md
  git commit --quiet -m "α-456: implement feature"

  echo "change2" >> README.md
  git add README.md
  git commit --quiet -m "β-789: review complete"

  # Test
  local output
  output=$(bash ./generate-tag-message.sh 1.1.0 v1.0.0)
  local issue_count
  issue_count=$(echo "$output" | grep -c "^- #" || echo "0")

  test_result "Two issues extracted" "2" "$issue_count"

  cd - >/dev/null
}

# Test 4: Deterministic output (same input, same output)
test_deterministic_output() {
  echo "Testing: Deterministic output"
  local test_dir="/tmp/test-tag-gen-deterministic"

  setup_test_repo "$test_dir"

  echo "change" >> README.md
  git add README.md
  git commit --quiet -m "Fix issue (#999)"

  # Generate twice
  local output1
  local output2
  output1=$(./generate-tag-message.sh 1.1.0 v1.0.0 | grep -v "Generated:")
  output2=$(./generate-tag-message.sh 1.1.0 v1.0.0 | grep -v "Generated:")

  local matches
  if [ "$output1" = "$output2" ]; then
    matches="1"
  else
    matches="0"
  fi

  test_result "Deterministic output (excluding timestamp)" "1" "$matches"

  cd - >/dev/null
}

# Test 5: Wave context handling
test_wave_context() {
  echo "Testing: Wave context handling"
  local test_dir="/tmp/test-tag-gen-wave"

  setup_test_repo "$test_dir"

  # Create wave fixture
  mkdir -p .cdd/waves/test-wave-2026-01-01
  cat > .cdd/waves/test-wave-2026-01-01/manifest.md << 'EOF'
# Test Wave

## Issues
- #111: First issue
- #222: Second issue
EOF

  echo "change" >> README.md
  git add README.md .cdd
  git commit --quiet -m "Add wave and fix (#111)"

  # Test
  local output
  output=$(bash ./generate-tag-message.sh 1.1.0 v1.0.0)
  local has_wave
  has_wave=$(echo "$output" | grep -c "Wave: test-wave-2026-01-01" || echo "0")

  test_result "Wave context detected" "1" "$has_wave"

  cd - >/dev/null
}

# Test 6: Fallback behavior without GitHub access
test_github_fallback() {
  echo "Testing: GitHub fallback behavior"
  local test_dir="/tmp/test-tag-gen-fallback"

  setup_test_repo "$test_dir"

  # Remove gh command temporarily
  export PATH="/bin:/usr/bin"

  echo "change" >> README.md
  git add README.md
  git commit --quiet -m "Fix issue (#888)"

  # Test
  local output
  output=$(bash ./generate-tag-message.sh 1.1.0 v1.0.0)
  local has_unavailable
  has_unavailable=$(echo "$output" | grep -c "unavailable" || echo "0")

  test_result "GitHub fallback with unavailable metadata" "1" "$has_unavailable"

  cd - >/dev/null
}

# Run all tests
echo "=== Tag Message Generator Test Suite ==="
echo ""

cd "$REPO_ROOT"

test_no_issues
test_single_issue_hash_format
test_multiple_issues_role_format
test_deterministic_output
test_wave_context
test_github_fallback

echo ""
echo "=== Test Results ==="
echo "Tests run: $TESTS_RUN"
echo "Tests passed: $TESTS_PASSED"

if [ "$TESTS_PASSED" -eq "$TESTS_RUN" ]; then
  echo "✅ All tests passed!"
  exit 0
else
  echo "❌ Some tests failed!"
  exit 1
fi