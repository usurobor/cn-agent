#!/usr/bin/env bash
# test-validate-release-gate.sh — prove validate-release-gate.sh enforces the
# pre-tag release gate from CDD.md §5.3b and §1.2.
#
# Tests AC1, AC2, AC3 from issue #327:
#   AC1: substantial (triadic) cycle requires all 5 lifecycle files
#   AC2: small-change cycle (no beta-review.md) collapses requirements per §1.2
#   AC3: RELEASE.md must exist at repo root before tag
#
# Run: bash scripts/test-validate-release-gate.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VALIDATE="$SCRIPT_DIR/validate-release-gate.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

PASS=0
FAIL=0
TESTS=0

assert_exit() {
  local name="$1"
  local expected="$2"
  local actual="$3"
  TESTS=$((TESTS + 1))
  if [[ "$expected" == "$actual" ]]; then
    PASS=$((PASS + 1))
    printf "  ✅ %s (exit %d)\n" "$name" "$actual"
  else
    FAIL=$((FAIL + 1))
    printf "  ❌ %s (expected exit %d, got %d)\n" "$name" "$expected" "$actual"
  fi
}

assert_grep() {
  local name="$1"
  local pattern="$2"
  local file="$3"
  TESTS=$((TESTS + 1))
  if grep -q "$pattern" "$file"; then
    PASS=$((PASS + 1))
    printf "  ✅ %s\n" "$name"
  else
    FAIL=$((FAIL + 1))
    printf "  ❌ %s (pattern not found: %s)\n" "$name" "$pattern"
    printf "     output:\n"
    sed 's/^/       /' "$file"
  fi
}

# ----- Fixture builders -----

make_root() {
  local root="$1"
  mkdir -p "$root"
}

write_release_md() {
  local root="$1"
  echo "# Release" > "$root/RELEASE.md"
}

# write_triadic_cycle DIR N [missing_file]
# Creates a cycle dir with all 5 required triadic files, optionally omitting one.
write_triadic_cycle() {
  local unreleased="$1"
  local N="$2"
  local missing="${3:-}"
  local dir="$unreleased/$N"
  mkdir -p "$dir"
  for f in self-coherence.md beta-review.md alpha-closeout.md beta-closeout.md gamma-closeout.md; do
    if [[ "$f" != "$missing" ]]; then
      echo "# $f stub" > "$dir/$f"
    fi
  done
}

# write_small_change_cycle DIR N
# Creates a cycle dir with only self-coherence.md (no beta-review.md = small-change).
write_small_change_cycle() {
  local unreleased="$1"
  local N="$2"
  mkdir -p "$unreleased/$N"
  echo "# small-change self-coherence" > "$unreleased/$N/self-coherence.md"
}

# ----- Test cases -----

# AC1 positive: all 5 required files present + RELEASE.md → gate passes
echo "## AC1 positive — triadic cycle with all required artifacts"
R="$TMP/ac1_pass"
make_root "$R"
write_release_md "$R"
write_triadic_cycle "$R/.cdd/unreleased" "200"
OUT="$TMP/ac1_pass.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 0 when all artifacts present" 0 "$EC"
assert_grep "cycle 200 triadic pass" "✅ cycle 200 (triadic)" "$OUT"
assert_grep "gate passes" "Release gate passed" "$OUT"

# AC1 negative: missing alpha-closeout.md → gate fails
echo ""
echo "## AC1 negative — missing alpha-closeout.md"
R="$TMP/ac1_no_alpha"
make_root "$R"
write_release_md "$R"
write_triadic_cycle "$R/.cdd/unreleased" "201" "alpha-closeout.md"
OUT="$TMP/ac1_no_alpha.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 1 when alpha-closeout.md missing" 1 "$EC"
assert_grep "error names alpha-closeout.md" "alpha-closeout.md" "$OUT"
assert_grep "gate failed" "Release gate FAILED" "$OUT"

# AC1 negative: missing gamma-closeout.md → gate fails
echo ""
echo "## AC1 negative — missing gamma-closeout.md"
R="$TMP/ac1_no_gamma"
make_root "$R"
write_release_md "$R"
write_triadic_cycle "$R/.cdd/unreleased" "202" "gamma-closeout.md"
OUT="$TMP/ac1_no_gamma.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 1 when gamma-closeout.md missing" 1 "$EC"
assert_grep "error names gamma-closeout.md" "gamma-closeout.md" "$OUT"

# AC2 positive: small-change cycle (no beta-review.md) → gate passes
echo ""
echo "## AC2 positive — small-change cycle passes without triadic files"
R="$TMP/ac2_small"
make_root "$R"
write_release_md "$R"
write_small_change_cycle "$R/.cdd/unreleased" "203"
OUT="$TMP/ac2_small.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 0 for small-change cycle" 0 "$EC"
assert_grep "small-change classification shown" "small-change" "$OUT"
assert_grep "gate passes for small-change" "Release gate passed" "$OUT"

# AC2 negative: small-change cycle NOT blocked for beta-closeout.md absence
echo ""
echo "## AC2 negative — small-change cycle not blocked for missing triadic files"
assert_grep "no error for missing beta-closeout in small-change" "Release gate passed" "$TMP/ac2_small.out"

# AC3 positive: RELEASE.md present → check passes
echo ""
echo "## AC3 positive — RELEASE.md present"
R="$TMP/ac3_pass"
make_root "$R"
write_release_md "$R"
OUT="$TMP/ac3_pass.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 0 when RELEASE.md present" 0 "$EC"
assert_grep "RELEASE.md check passes" "RELEASE.md present" "$OUT"

# AC3 negative: RELEASE.md absent → gate fails
echo ""
echo "## AC3 negative — RELEASE.md absent"
R="$TMP/ac3_fail"
make_root "$R"
# no RELEASE.md
OUT="$TMP/ac3_fail.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 1 when RELEASE.md absent" 1 "$EC"
assert_grep "error names RELEASE.md" "RELEASE.md missing" "$OUT"
assert_grep "gate failed for missing RELEASE.md" "Release gate FAILED" "$OUT"

# Bonus: no unreleased dir → gate still passes (no cycles to validate)
echo ""
echo "## Bonus — no .cdd/unreleased/ dir + RELEASE.md present → passes"
R="$TMP/no_unreleased"
make_root "$R"
write_release_md "$R"
OUT="$TMP/no_unreleased.out"
set +e
"$VALIDATE" --repo-root "$R" >"$OUT" 2>&1
EC=$?
set -e
assert_exit "exit 0 when no unreleased dir" 0 "$EC"
assert_grep "gate passes with no unreleased dir" "Release gate passed" "$OUT"

# ----- Summary -----

echo ""
echo "================================================================"
echo "ran $TESTS assertions: $PASS passed, $FAIL failed"
if [[ $FAIL -gt 0 ]]; then
  exit 1
fi
exit 0
