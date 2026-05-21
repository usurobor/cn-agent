#!/usr/bin/env bash
# test-fixtures.sh — validate cdd-verify behavior against test fixtures

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../../.." && pwd)"

# cycle/392: cdd-verify is now the cn kernel subcommand; invoke via the
# built cn binary using the canonical noun-verb form `cn cdd verify`.
CN_BIN="${CN_BIN:-}"
if [ -z "$CN_BIN" ]; then
  CN_BIN="$REPO_ROOT/.cdd-cache/cn-test"
  mkdir -p "$(dirname "$CN_BIN")"
  ( cd "$REPO_ROOT" && go build -o "$CN_BIN" ./src/go/cmd/cn ) >&2
fi
CDD_VERIFY=("$CN_BIN" cdd verify)

echo "Testing CDD artifact checker fixtures..."
echo

# Test 1: Valid triadic cycle should pass
echo "=== Test 1: Valid triadic cycle ==="
if "${CDD_VERIFY[@]}" --unreleased --repo-root "$SCRIPT_DIR/test/fixtures/valid-triadic" | grep -q "Cycle artifact verification PASSED"; then
    echo "✅ Valid triadic cycle test PASSED"
else
    echo "❌ Valid triadic cycle test FAILED"
    "${CDD_VERIFY[@]}" --unreleased --repo-root "$SCRIPT_DIR/test/fixtures/valid-triadic"
    exit 1
fi

# Test 2: Incomplete triadic cycle should fail
echo ""
echo "=== Test 2: Incomplete triadic cycle ==="
set +e
OUTPUT=$("${CDD_VERIFY[@]}" --all --repo-root "$SCRIPT_DIR/test/fixtures/incomplete-triadic" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -ne 0 ]] && echo "$OUTPUT" | grep -q "Cycle artifact verification FAILED"; then
    echo "✅ Incomplete triadic cycle test PASSED (correctly detected missing artifacts)"
else
    echo "❌ Incomplete triadic cycle test FAILED"
    echo "$OUTPUT"
    exit 1
fi

# Test 3: Valid small-change cycle should pass
echo ""
echo "=== Test 3: Valid small-change cycle ==="
OUTPUT=$("${CDD_VERIFY[@]}" --unreleased --repo-root "$SCRIPT_DIR/test/fixtures/valid-small-change")
if echo "$OUTPUT" | grep -q "Checking small-change cycle #300"; then
    echo "✅ Valid small-change cycle test PASSED"
else
    echo "❌ Valid small-change cycle test FAILED"
    echo "Expected: 'Checking small-change cycle #300'"
    echo "Actual output:"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo "🎉 All fixture tests passed!"