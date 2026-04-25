#!/usr/bin/env bash
# test-cn-cdd-verify.sh — verify cn-cdd-verify enforces the canonical CDD
# contract from CDD.md §5.3a and treats legacy paths/tags as warn-only.
#
# This test builds throwaway fixture repos under a temp directory, runs the
# verifier with --repo-root pointing at them, and asserts the expected exit
# code and output substrings.
#
# Run:
#   bash src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERIFY="$SCRIPT_DIR/cn-cdd-verify"
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

make_repo() {
  local repo="$1"
  mkdir -p "$repo"
  git -C "$repo" init -q -b main
  git -C "$repo" config user.email "test@cnos.local"
  git -C "$repo" config user.name "test"
  git -C "$repo" config commit.gpgsign false
  git -C "$repo" config tag.gpgsign false
  : > "$repo/.gitkeep"
  git -C "$repo" add .gitkeep
  git -C "$repo" commit -q -m "init"
}

write_changelog() {
  local repo="$1"
  local version="$2"
  cat > "$repo/CHANGELOG.md" <<EOF
# CHANGELOG

| $version | C_Σ | α | β | γ | L6 | first canonical entry |
EOF
  git -C "$repo" add CHANGELOG.md >/dev/null
  git -C "$repo" commit -q -m "changelog $version"
}

write_canonical_pra() {
  local repo="$1"
  local version="$2"
  mkdir -p "$repo/docs/gamma/cdd/$version"
  cat > "$repo/docs/gamma/cdd/$version/POST-RELEASE-ASSESSMENT.md" <<EOF
## Post-Release Assessment — $version

(test fixture)
EOF
  git -C "$repo" add "docs/gamma/cdd/$version/POST-RELEASE-ASSESSMENT.md" >/dev/null
  git -C "$repo" commit -q -m "PRA $version"
}

write_legacy_pra() {
  local repo="$1"
  local version="$2"
  mkdir -p "$repo/.cdd/releases/$version/beta"
  echo "legacy" > "$repo/.cdd/releases/$version/beta/POST-RELEASE-ASSESSMENT.md"
  git -C "$repo" add ".cdd/releases/$version/beta/POST-RELEASE-ASSESSMENT.md" >/dev/null
  git -C "$repo" commit -q -m "legacy PRA $version"
}

write_canonical_closeouts() {
  local repo="$1"
  local version="$2"
  for role in alpha beta gamma; do
    mkdir -p "$repo/.cdd/releases/$version/$role"
    echo "$role close-out for $version" > "$repo/.cdd/releases/$version/$role/CLOSE-OUT.md"
  done
  git -C "$repo" add ".cdd/releases/$version" >/dev/null
  git -C "$repo" commit -q -m "close-outs $version"
}

# ----- Test cases -----

echo "## test 1 — canonical contract passes"
R="$TMP/canonical"
make_repo "$R"
write_changelog "$R" "9.0.0"
write_canonical_pra "$R" "9.0.0"
git -C "$R" tag "9.0.0"
OUT="$TMP/canonical.out"
set +e
"$VERIFY" --version 9.0.0 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "exits 0 on canonical contract" 0 "$EC"
assert_grep "canonical PRA pass" "Canonical PRA" "$OUT"
assert_grep "bare tag pass" "Git tag (bare 9.0.0)" "$OUT"
assert_grep "no legacy PRA warning" "Cycle artifact verification PASSED" "$OUT"

echo ""
echo "## test 2 — v-prefixed tag warns but does not fail"
R="$TMP/legacy_tag"
make_repo "$R"
write_changelog "$R" "9.1.0"
write_canonical_pra "$R" "9.1.0"
git -C "$R" tag "9.1.0"
git -C "$R" tag "v9.1.0"
OUT="$TMP/legacy_tag.out"
set +e
"$VERIFY" --version 9.1.0 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "exits 0 with v-tag warning" 0 "$EC"
assert_grep "v-prefixed tag warns" "Legacy v-prefixed tag" "$OUT"

echo ""
echo "## test 3 — legacy PRA path warns, canonical missing fails"
R="$TMP/legacy_pra_only"
make_repo "$R"
write_changelog "$R" "9.2.0"
write_legacy_pra "$R" "9.2.0"
git -C "$R" tag "9.2.0"
OUT="$TMP/legacy_pra_only.out"
set +e
"$VERIFY" --version 9.2.0 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "fails when only legacy PRA exists" 1 "$EC"
assert_grep "canonical PRA missing fail" "❌ Canonical PRA" "$OUT"
assert_grep "legacy PRA warn" "Legacy PRA path" "$OUT"

echo ""
echo "## test 4 — missing bare tag fails (only v-prefixed present)"
R="$TMP/v_only"
make_repo "$R"
write_changelog "$R" "9.3.0"
write_canonical_pra "$R" "9.3.0"
git -C "$R" tag "v9.3.0"
OUT="$TMP/v_only.out"
set +e
"$VERIFY" --version 9.3.0 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "fails when bare tag absent" 1 "$EC"
assert_grep "bare tag fail message" "❌ Git tag (bare 9.3.0)" "$OUT"

echo ""
echo "## test 5 — triadic canonical close-outs pass"
R="$TMP/triadic_canon"
make_repo "$R"
write_changelog "$R" "9.4.0"
write_canonical_pra "$R" "9.4.0"
write_canonical_closeouts "$R" "9.4.0"
git -C "$R" tag "9.4.0"
OUT="$TMP/triadic_canon.out"
set +e
"$VERIFY" --version 9.4.0 --triadic --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "triadic exits 0 on canonical close-outs" 0 "$EC"
assert_grep "alpha CLOSE-OUT pass" "✅ alpha/CLOSE-OUT.md" "$OUT"
assert_grep "beta CLOSE-OUT pass"  "✅ beta/CLOSE-OUT.md" "$OUT"
assert_grep "gamma CLOSE-OUT pass" "✅ gamma/CLOSE-OUT.md" "$OUT"

echo ""
echo "## test 6 — triadic missing γ close-out fails"
R="$TMP/triadic_no_gamma"
make_repo "$R"
write_changelog "$R" "9.5.0"
write_canonical_pra "$R" "9.5.0"
mkdir -p "$R/.cdd/releases/9.5.0/alpha" "$R/.cdd/releases/9.5.0/beta"
echo "α" > "$R/.cdd/releases/9.5.0/alpha/CLOSE-OUT.md"
echo "β" > "$R/.cdd/releases/9.5.0/beta/CLOSE-OUT.md"
git -C "$R" add .cdd >/dev/null && git -C "$R" commit -q -m "partial close-outs"
git -C "$R" tag "9.5.0"
OUT="$TMP/triadic_no_gamma.out"
set +e
"$VERIFY" --version 9.5.0 --triadic --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "triadic fails without γ close-out" 1 "$EC"
assert_grep "γ close-out fail" "❌ gamma/CLOSE-OUT.md" "$OUT"

# ----- Summary -----

echo ""
echo "================================================================"
echo "ran $TESTS assertions: $PASS passed, $FAIL failed"
if [[ $FAIL -gt 0 ]]; then
  exit 1
fi
exit 0
