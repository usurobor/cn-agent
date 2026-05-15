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

write_cycle_artifacts() {
  local repo="$1"
  local version="$2"
  local cycle="$3"
  mkdir -p "$repo/.cdd/releases/$version/$cycle"

  # Write 5 hard-gate artifacts with required sections
  cat > "$repo/.cdd/releases/$version/$cycle/self-coherence.md" <<EOF
# Self-Coherence — Issue #$cycle

## Gap
Test gap description

## Mode
Test mode description

## AC Coverage
Test AC coverage

## CDD Trace
Test trace

## Role self-check
Test self-check

## Known debt
No debt
EOF

  cat > "$repo/.cdd/releases/$version/$cycle/beta-review.md" <<EOF
# Beta Review — Issue #$cycle

## Round 1
**Verdict:** APPROVED
**Findings:** None
EOF

  for artifact in "alpha-closeout" "beta-closeout" "gamma-closeout"; do
    cat > "$repo/.cdd/releases/$version/$cycle/${artifact}.md" <<EOF
# ${artifact^} — Issue #$cycle

## Summary
Test summary

## Findings
No findings
EOF
  done

  git -C "$repo" add ".cdd/releases/$version/$cycle" >/dev/null
  git -C "$repo" commit -q -m "cycle artifacts $version/$cycle"
}

write_incomplete_cycle_artifacts() {
  local repo="$1"
  local version="$2"
  local cycle="$3"
  mkdir -p "$repo/.cdd/releases/$version/$cycle"

  # Write only some artifacts (missing gamma-closeout.md)
  cat > "$repo/.cdd/releases/$version/$cycle/self-coherence.md" <<EOF
# Self-Coherence — Issue #$cycle

## Gap
Test gap

## Mode
Test mode

## CDD Trace
Test trace
EOF

  for artifact in "beta-review" "alpha-closeout" "beta-closeout"; do
    echo "Test content" > "$repo/.cdd/releases/$version/$cycle/${artifact}.md"
  done

  git -C "$repo" add ".cdd/releases/$version/$cycle" >/dev/null
  git -C "$repo" commit -q -m "incomplete cycle artifacts $version/$cycle"
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

echo ""
echo "## test 7 — cycle mode complete artifacts pass"
R="$TMP/cycle_complete"
make_repo "$R"
write_changelog "$R" "10.0.0"
write_canonical_pra "$R" "10.0.0"
write_cycle_artifacts "$R" "10.0.0" "100"
git -C "$R" tag "10.0.0"
OUT="$TMP/cycle_complete.out"
set +e
"$VERIFY" --version 10.0.0 --cycle 100 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "cycle mode exits 0 on complete artifacts" 0 "$EC"
assert_grep "self-coherence pass" "✅ self-coherence.md" "$OUT"
assert_grep "beta-review pass" "✅ beta-review.md" "$OUT"
assert_grep "alpha-closeout pass" "✅ alpha-closeout.md" "$OUT"
assert_grep "beta-closeout pass" "✅ beta-closeout.md" "$OUT"
assert_grep "gamma-closeout pass" "✅ gamma-closeout.md" "$OUT"

echo ""
echo "## test 8 — cycle mode missing artifact fails"
R="$TMP/cycle_incomplete"
make_repo "$R"
write_changelog "$R" "10.1.0"
write_canonical_pra "$R" "10.1.0"
write_incomplete_cycle_artifacts "$R" "10.1.0" "101"
git -C "$R" tag "10.1.0"
OUT="$TMP/cycle_incomplete.out"
set +e
"$VERIFY" --version 10.1.0 --cycle 101 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "cycle mode fails with missing artifact" 1 "$EC"
assert_grep "gamma-closeout missing" "❌ gamma-closeout.md" "$OUT"

echo ""
echo "## test 9 — cycle mode requires version"
OUT="$TMP/cycle_no_version.out"
set +e
"$VERIFY" --cycle 123 > "$OUT" 2>&1
EC=$?
set -e
assert_exit "cycle mode fails without version" 1 "$EC"
assert_grep "cycle requires version error" "Error: --cycle requires --version" "$OUT"

echo ""
echo "## test 10 — cycle mode nonexistent directory fails"
R="$TMP/cycle_missing_dir"
make_repo "$R"
write_changelog "$R" "10.2.0"
write_canonical_pra "$R" "10.2.0"
git -C "$R" tag "10.2.0"
OUT="$TMP/cycle_missing_dir.out"
set +e
"$VERIFY" --version 10.2.0 --cycle 999 --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "cycle mode fails with missing directory" 1 "$EC"
assert_grep "missing artifacts" "❌.*expected.*999" "$OUT"

# ----- Era-boundary tests (#365) -----
#
# is_legacy_version cutoff = v3.77.0. Pre-v3.77.0 cycles missing
# alpha-closeout.md / gamma-closeout.md or the "CDD Trace" section in
# self-coherence.md must warn (not fail) so CI does not block on
# historical state. v3.77.0+ cycles must continue to fail strict.
#
# Each test exercises one side of the boundary:
#   test 11 — v3.75.0 cycle missing close-outs  → warn, exit 0
#   test 12 — v3.76.0 self-coherence missing CDD Trace → warn, exit 0
#   test 13 — v3.77.0 cycle missing close-outs  → fail, exit 1
# The fixture writers below intentionally produce *only* the artifacts the
# era under test had; the strict path is exercised by test 13 to confirm
# the era policy does not bleed past the cutoff (AC3 — strict enforcement
# preserved for current cycles).

write_legacy_partial_cycle() {
  # Self-coherence + beta-review only; close-outs absent. Mirrors the
  # shape of every v3.75.0 cycle directory in this repo.
  local repo="$1"
  local version="$2"
  local cycle="$3"
  mkdir -p "$repo/.cdd/releases/$version/$cycle"
  cat > "$repo/.cdd/releases/$version/$cycle/self-coherence.md" <<EOF
# Self-Coherence — Issue #$cycle (v$version)

## Gap
Era-boundary test fixture (v$version).

## Mode
Test fixture for is_legacy_version boundary.

## ACs
AC1: cycle pre-dates close-out requirement.

## Self-check
No ambiguity pushed forward.

## Debt
None.

## CDD Trace
| Step | Decision |
|------|----------|
| 7 | Test fixture |
EOF
  cat > "$repo/.cdd/releases/$version/$cycle/beta-review.md" <<EOF
# β Review — #$cycle
**Verdict:** APPROVED
EOF
  git -C "$repo" add ".cdd/releases/$version/$cycle" >/dev/null
  git -C "$repo" commit -q -m "legacy partial cycle $version/$cycle"
}

write_no_cdd_trace_cycle() {
  # Full close-outs but no "CDD Trace" section in self-coherence.md.
  # Mirrors the shape of v3.76.0 cycle 362 in this repo.
  local repo="$1"
  local version="$2"
  local cycle="$3"
  mkdir -p "$repo/.cdd/releases/$version/$cycle"
  cat > "$repo/.cdd/releases/$version/$cycle/self-coherence.md" <<EOF
# Self-Coherence — Issue #$cycle (v$version)

## Gap
Era-boundary test fixture (v$version).

## Mode
Test fixture.

## ACs
AC1: cycle pre-dates "CDD Trace" section requirement.

## Self-check
None.

## Debt
None.
EOF
  cat > "$repo/.cdd/releases/$version/$cycle/beta-review.md" <<EOF
# β Review — #$cycle
**Verdict:** APPROVED
EOF
  for role in alpha beta gamma; do
    echo "# ${role^} close-out — #$cycle (Summary, Findings: none)" > "$repo/.cdd/releases/$version/$cycle/${role}-closeout.md"
  done
  git -C "$repo" add ".cdd/releases/$version/$cycle" >/dev/null
  git -C "$repo" commit -q -m "no-trace legacy cycle $version/$cycle"
}

echo ""
echo "## test 11 — v3.75.0 cycle missing close-outs warns (legacy era)"
R="$TMP/era_v3_75"
make_repo "$R"
write_legacy_partial_cycle "$R" "3.75.0" "375"
OUT="$TMP/era_v3_75.out"
set +e
"$VERIFY" --all --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "v3.75.0 missing close-outs exits 0" 0 "$EC"
assert_grep "alpha-closeout warns on v3.75.0" "⚠️  alpha-closeout.md.*#375.*missing in legacy cycle (pre-v3.77.0)" "$OUT"
assert_grep "gamma-closeout warns on v3.75.0" "⚠️  gamma-closeout.md.*#375.*missing in legacy cycle (pre-v3.77.0)" "$OUT"
assert_grep "no failures emitted for v3.75.0"  "0 failed" "$OUT"

echo ""
echo "## test 12 — v3.76.0 self-coherence missing CDD Trace warns (legacy era)"
R="$TMP/era_v3_76"
make_repo "$R"
write_no_cdd_trace_cycle "$R" "3.76.0" "376"
OUT="$TMP/era_v3_76.out"
set +e
"$VERIFY" --all --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "v3.76.0 missing CDD Trace exits 0" 0 "$EC"
assert_grep "CDD Trace warns on v3.76.0" "⚠️  self-coherence.md sections.*missing required sections in legacy cycle (pre-v3.77.0).*CDD Trace" "$OUT"
assert_grep "no failures emitted for v3.76.0" "0 failed" "$OUT"

echo ""
echo "## test 13 — v3.77.0 cycle missing close-outs still fails (strict)"
R="$TMP/era_v3_77"
make_repo "$R"
write_legacy_partial_cycle "$R" "3.77.0" "377"
OUT="$TMP/era_v3_77.out"
set +e
"$VERIFY" --all --repo-root "$R" > "$OUT" 2>&1
EC=$?
set -e
assert_exit "v3.77.0 missing close-outs exits 1" 1 "$EC"
assert_grep "alpha-closeout fails on v3.77.0" "❌ alpha-closeout.md.*#377" "$OUT"
assert_grep "gamma-closeout fails on v3.77.0" "❌ gamma-closeout.md.*#377" "$OUT"

# ----- Summary -----

echo ""
echo "================================================================"
echo "ran $TESTS assertions: $PASS passed, $FAIL failed"
if [[ $FAIL -gt 0 ]]; then
  exit 1
fi
exit 0
