#!/usr/bin/env bash
# test-validate-release-gate.sh — phase-split regression suite for the CDD
# pre-merge, post-merge closure, and release gates.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VALIDATE="$SCRIPT_DIR/validate-release-gate.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

PASS=0
FAIL=0
TESTS=0

assert_exit() {
  local name="$1" expected="$2" actual="$3"
  TESTS=$((TESTS + 1))
  if [[ "$expected" == "$actual" ]]; then
    PASS=$((PASS + 1)); printf "  ✅ %s (exit %d)\n" "$name" "$actual"
  else
    FAIL=$((FAIL + 1)); printf "  ❌ %s (expected exit %d, got %d)\n" "$name" "$expected" "$actual"
  fi
}

assert_grep() {
  local name="$1" pattern="$2" file="$3"
  TESTS=$((TESTS + 1))
  if grep -q "$pattern" "$file"; then
    PASS=$((PASS + 1)); printf "  ✅ %s\n" "$name"
  else
    FAIL=$((FAIL + 1)); printf "  ❌ %s (pattern not found: %s)\n" "$name" "$pattern"
    sed 's/^/       /' "$file"
  fi
}

make_root() { mkdir -p "$1"; }
write_release_md() { printf '# Release\n' > "$1/RELEASE.md"; }

write_triadic_cycle() {
  local unreleased="$1" cycle="$2" missing="${3:-}" terminal="${4:-yes}"
  local dir="$unreleased/$cycle"
  mkdir -p "$dir"
  for f in gamma-scaffold.md self-coherence.md beta-review.md alpha-closeout.md beta-closeout.md gamma-closeout.md; do
    if [[ "$f" != "$missing" ]]; then
      printf '# %s stub\n' "$f" > "$dir/$f"
    fi
  done
  if [[ "$terminal" == "yes" && "$missing" != "gamma-closeout.md" ]]; then
    printf '\nCDD-Cycle-Closure: terminal\n' >> "$dir/gamma-closeout.md"
  fi
}

write_small_change_cycle() {
  mkdir -p "$1/$2"
  printf '# small-change self-coherence\n' > "$1/$2/self-coherence.md"
}

run_gate() {
  local output="$1"; shift
  set +e
  "$VALIDATE" "$@" >"$output" 2>&1
  RUN_EXIT=$?
  set -e
}

echo "## pre-merge accepts only structurally available triadic artifacts"
R="$TMP/pre_pass"; make_root "$R"
write_triadic_cycle "$R/.cdd/unreleased" 200 alpha-closeout.md no
rm "$R/.cdd/unreleased/200/beta-closeout.md" "$R/.cdd/unreleased/200/gamma-closeout.md"
OUT="$TMP/pre_pass.out"; run_gate "$OUT" --repo-root "$R" --mode pre-merge --cycle 200
assert_exit "pre-merge passes without future close-outs" 0 "$RUN_EXIT"
assert_grep "pre-merge pass diagnostic" "Pre-merge gate passed" "$OUT"

echo "## pre-merge refuses a missing review-time artifact"
R="$TMP/pre_fail"; make_root "$R"
write_triadic_cycle "$R/.cdd/unreleased" 201 beta-review.md
OUT="$TMP/pre_fail.out"; run_gate "$OUT" --repo-root "$R" --mode pre-merge --cycle 201
assert_exit "pre-merge fails without beta-review" 1 "$RUN_EXIT"
assert_grep "pre-merge names beta-review" "missing beta-review.md" "$OUT"

echo "## pre-merge refuses a missing gamma scaffold"
R="$TMP/pre_no_scaffold"; make_root "$R"
write_triadic_cycle "$R/.cdd/unreleased" 212 gamma-scaffold.md
OUT="$TMP/pre_no_scaffold.out"; run_gate "$OUT" --repo-root "$R" --mode pre-merge --cycle 212
assert_exit "pre-merge fails without gamma-scaffold" 1 "$RUN_EXIT"
assert_grep "pre-merge names gamma-scaffold" "missing gamma-scaffold.md" "$OUT"

echo "## gamma scaffold classifies an in-flight cycle as triadic"
R="$TMP/pre_scaffold"; make_root "$R"
mkdir -p "$R/.cdd/unreleased/202"; printf '# scaffold\n' > "$R/.cdd/unreleased/202/gamma-scaffold.md"
OUT="$TMP/pre_scaffold.out"; run_gate "$OUT" --repo-root "$R" --mode pre-merge --cycle 202
assert_exit "scaffold-only cycle fails triadic gate" 1 "$RUN_EXIT"
assert_grep "scaffold-only cycle names self-coherence" "missing self-coherence.md" "$OUT"

echo "## post-merge requires close-outs but not RELEASE.md"
R="$TMP/post_pass"; make_root "$R"; write_triadic_cycle "$R/.cdd/unreleased" 203
OUT="$TMP/post_pass.out"; run_gate "$OUT" --repo-root "$R" --mode post-merge --cycle 203
assert_exit "post-merge passes complete lifecycle without RELEASE.md" 0 "$RUN_EXIT"
assert_grep "post-merge pass diagnostic" "Post-merge closure gate passed" "$OUT"

echo "## post-merge refuses absent alpha close-out"
R="$TMP/post_no_alpha"; make_root "$R"; write_triadic_cycle "$R/.cdd/unreleased" 204 alpha-closeout.md
OUT="$TMP/post_no_alpha.out"; run_gate "$OUT" --repo-root "$R" --mode post-merge --cycle 204
assert_exit "post-merge fails without alpha-closeout" 1 "$RUN_EXIT"
assert_grep "post-merge names alpha-closeout" "missing alpha-closeout.md" "$OUT"

echo "## post-merge distinguishes assurance receipt from terminal gamma closure"
R="$TMP/post_marker"; make_root "$R"; write_triadic_cycle "$R/.cdd/unreleased" 205 "" no
OUT="$TMP/post_marker.out"; run_gate "$OUT" --repo-root "$R" --mode post-merge --cycle 205
assert_exit "post-merge fails without terminal marker" 1 "$RUN_EXIT"
assert_grep "marker failure is explicit" "lacks terminal closure marker" "$OUT"

echo "## exact-cycle selector ignores unrelated incomplete cycles"
R="$TMP/selector"; make_root "$R"
write_triadic_cycle "$R/.cdd/unreleased" 206
write_triadic_cycle "$R/.cdd/unreleased" 207 beta-closeout.md
OUT="$TMP/selector.out"; run_gate "$OUT" --repo-root "$R" --mode post-merge --cycle 206
assert_exit "selected complete cycle passes" 0 "$RUN_EXIT"
if grep -q "cycle 207" "$OUT"; then
  FAIL=$((FAIL + 1)); TESTS=$((TESTS + 1)); printf "  ❌ selector inspected unrelated cycle 207\n"
else
  PASS=$((PASS + 1)); TESTS=$((TESTS + 1)); printf "  ✅ selector excludes unrelated cycle 207\n"
fi

echo "## exact-cycle selector fails when target directory is absent"
OUT="$TMP/selector_missing.out"; run_gate "$OUT" --repo-root "$R" --mode post-merge --cycle 999
assert_exit "missing selected cycle fails" 1 "$RUN_EXIT"
assert_grep "missing cycle diagnostic names target" "cycle 999: directory missing" "$OUT"

echo "## release mode preserves RELEASE.md + full lifecycle enforcement"
R="$TMP/release_pass"; make_root "$R"; write_release_md "$R"; write_triadic_cycle "$R/.cdd/unreleased" 208
OUT="$TMP/release_pass.out"; run_gate "$OUT" --repo-root "$R" --cycle 208
assert_exit "release passes complete lifecycle" 0 "$RUN_EXIT"
assert_grep "release reports RELEASE.md" "RELEASE.md present" "$OUT"
assert_grep "release pass diagnostic" "Release gate passed" "$OUT"

echo "## release mode preserves full lifecycle refusal"
R="$TMP/release_no_alpha"; make_root "$R"; write_release_md "$R"
write_triadic_cycle "$R/.cdd/unreleased" 213 alpha-closeout.md
OUT="$TMP/release_no_alpha.out"; run_gate "$OUT" --repo-root "$R" --cycle 213
assert_exit "release fails without alpha-closeout" 1 "$RUN_EXIT"
assert_grep "release names alpha-closeout" "missing alpha-closeout.md" "$OUT"

echo "## release mode still refuses missing RELEASE.md"
R="$TMP/release_no_notes"; make_root "$R"; write_triadic_cycle "$R/.cdd/unreleased" 209
OUT="$TMP/release_no_notes.out"; run_gate "$OUT" --repo-root "$R" --cycle 209
assert_exit "release fails without RELEASE.md" 1 "$RUN_EXIT"
assert_grep "release names RELEASE.md" "RELEASE.md missing" "$OUT"

echo "## release mode requires terminal gamma marker"
R="$TMP/release_marker"; make_root "$R"; write_release_md "$R"; write_triadic_cycle "$R/.cdd/unreleased" 210 "" no
OUT="$TMP/release_marker.out"; run_gate "$OUT" --repo-root "$R" --cycle 210
assert_exit "release fails without terminal marker" 1 "$RUN_EXIT"
assert_grep "release marker failure is explicit" "lacks terminal closure marker" "$OUT"

echo "## small-change collapse remains intact"
R="$TMP/small"; make_root "$R"; write_release_md "$R"; write_small_change_cycle "$R/.cdd/unreleased" 211
OUT="$TMP/small.out"; run_gate "$OUT" --repo-root "$R" --cycle 211
assert_exit "small-change release gate passes" 0 "$RUN_EXIT"
assert_grep "small-change classification shown" "small-change" "$OUT"

echo ""
echo "================================================================"
echo "ran $TESTS assertions: $PASS passed, $FAIL failed"
[[ $FAIL -eq 0 ]]
