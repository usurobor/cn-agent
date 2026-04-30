#!/usr/bin/env bash
# tools/validate-skill-frontmatter.sh — validate every SKILL.md frontmatter
# against schemas/skill.cue (#301, I5 coherence-CI job).
#
# The CUE schema owns shape / type / enum constraints. This script owns
# everything outside that: file discovery, frontmatter extraction,
# exception-list handling, and `calls` filesystem-existence checks
# (#301 AC2 surface boundary).
#
# Exit codes:
#   0  every SKILL.md (or fixture) passed
#   1  one or more validation failures
#   2  prerequisite missing (cue, jq) — not a validation failure
#
# Usage:
#   ./tools/validate-skill-frontmatter.sh
#       Validate all SKILL.md under src/packages/.
#   ./tools/validate-skill-frontmatter.sh --root <dir>
#       Validate every SKILL.md under <dir> instead of src/packages/.
#   ./tools/validate-skill-frontmatter.sh --self-test
#       Run schemas/fixtures/skill-frontmatter/{valid,invalid}/ as the
#       built-in positive/negative regression suite (#301 AC8).
#   ./tools/validate-skill-frontmatter.sh --file <path>
#       Validate a single SKILL.md path. Used by --self-test internally.
#
# Honors NO_COLOR. Diagnostics are machine-readable: one finding per line,
# `path :: field :: rule :: reason :: suggested-fix`.

set -euo pipefail

# --- color / symbols ----------------------------------------------------
if [[ -n "${NO_COLOR:-}" || ! -t 1 ]]; then
  RED="" GREEN="" YELLOW="" CYAN="" RESET=""
else
  RED=$'\033[0;31m' GREEN=$'\033[0;32m' YELLOW=$'\033[0;33m' CYAN=$'\033[0;36m' RESET=$'\033[0m'
fi
SYM_OK="✓" SYM_FAIL="✗" SYM_WARN="⚠"

# --- prereqs ------------------------------------------------------------
need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "${RED}${SYM_FAIL}${RESET} prerequisite missing: $1" >&2
    echo "  Install it before running this script." >&2
    exit 2
  }
}
need cue
need jq
need awk
need find

# --- paths --------------------------------------------------------------
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
SCHEMA="${REPO_ROOT}/schemas/skill.cue"
EXCEPTIONS="${REPO_ROOT}/schemas/skill-exceptions.json"
FIXTURE_VALID="${REPO_ROOT}/schemas/fixtures/skill-frontmatter/valid"
FIXTURE_INVALID="${REPO_ROOT}/schemas/fixtures/skill-frontmatter/invalid"
DEFAULT_ROOT="${REPO_ROOT}/src/packages"

[[ -f "$SCHEMA" ]] || {
  echo "${RED}${SYM_FAIL}${RESET} schema not found at: $SCHEMA" >&2
  exit 2
}

# --- args ---------------------------------------------------------------
mode="all"
root="$DEFAULT_ROOT"
single_file=""
while (($#)); do
  case "$1" in
    --root)
      root="$2"; shift 2;;
    --self-test)
      mode="self-test"; shift;;
    --file)
      mode="file"; single_file="$2"; shift 2;;
    -h|--help)
      # Print the leading comment block as usage text.
      awk '
        /^#!/ { next }
        /^#/  { sub(/^# ?/, ""); print; next }
        { exit }
      ' "$0"
      exit 0;;
    *)
      echo "${RED}${SYM_FAIL}${RESET} unknown argument: $1" >&2
      exit 2;;
  esac
done

# --- fields the script enforces beyond CUE -------------------------------
SPEC_REQUIRED_EXCEPTION_BACKED=(artifact_class kata_surface inputs outputs)

# --- temp workspace -----------------------------------------------------
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

# --- exception lookup ---------------------------------------------------
# Returns a JSON array of allowed_missing for the given path, or [] if no
# exception entry matches. Path is repo-relative so the JSON is portable.
get_allowed_missing() {
  local rel="$1"
  if [[ ! -f "$EXCEPTIONS" ]]; then
    echo "[]"
    return
  fi
  jq --arg p "$rel" '
    [ .[] | select(.path == $p) | .allowed_missing[]? ] | unique
  ' "$EXCEPTIONS"
}

# --- diagnostics --------------------------------------------------------
findings=0
emit_finding() {
  local path="$1" field="$2" rule="$3" reason="$4" fix="$5"
  printf "%s%s%s %s :: %s :: %s :: %s :: %s\n" \
    "$RED" "$SYM_FAIL" "$RESET" "$path" "$field" "$rule" "$reason" "$fix"
  findings=$((findings + 1))
}

# --- per-skill validation ----------------------------------------------
# Returns 0 (pass) / 1 (fail). Adds findings as side effect.
validate_skill_file() {
  local file="$1"
  local rel="${file#"$REPO_ROOT"/}"

  # 1. Extract frontmatter. Spec form: file starts with `---\n`, second
  # `---\n` closes the block. Reject anything that doesn't match the shape
  # so we never silently skip a malformed file (#301 AC2).
  local first_line
  first_line=$(head -n 1 "$file" 2>/dev/null || true)
  if [[ "$first_line" != "---" ]]; then
    emit_finding "$rel" "(frontmatter)" "extract" \
      "no opening '---' delimiter on line 1" \
      "add YAML frontmatter starting with '---' on line 1"
    return 1
  fi

  local fm
  fm=$(awk '
    /^---$/ { n++; if (n==1) next; if (n==2) exit }
    n==1 { print }
  ' "$file")

  if [[ -z "$fm" ]]; then
    emit_finding "$rel" "(frontmatter)" "extract" \
      "frontmatter block is empty or has no closing '---'" \
      "add at least name/description/governing_question/triggers/scope between two '---' delimiters"
    return 1
  fi

  local yaml_path="${TMPDIR}/skill.yaml"
  printf '%s\n' "$fm" > "$yaml_path"

  # 2. cue vet — schema-level type/enum/hard-gate enforcement.
  local cue_err
  if ! cue_err=$(cue vet -d '#Skill' "$SCHEMA" "$yaml_path" 2>&1); then
    # Translate cue's multi-line output into per-line findings. cue's
    # diagnostics name the field, so we forward them as-is rather than
    # restating the rule in shell.
    while IFS= read -r line; do
      [[ -z "$line" ]] && continue
      # Skip cue's source-pointer continuations (lines starting with
      # whitespace or with our tmp path).
      if [[ "$line" =~ ^[[:space:]] || "$line" == *"$yaml_path"* ]]; then
        continue
      fi
      emit_finding "$rel" "(cue)" "schema" "$line" \
        "fix the field per schemas/skill.cue or quote ambiguous YAML values"
    done <<<"$cue_err"
    return 1
  fi

  # 3. Convert YAML to JSON for the script-side checks.
  local json_path="${TMPDIR}/skill.json"
  if ! cue export --out json -d '#Skill' "$SCHEMA" "$yaml_path" >"$json_path" 2>/dev/null; then
    # cue vet already passed, so this would be a tooling bug; surface it.
    emit_finding "$rel" "(cue export)" "tooling" \
      "cue export failed despite cue vet passing" \
      "rerun with CUE_DEBUG=1 to investigate"
    return 1
  fi

  # 4. Spec-required-but-exception-backed: each must be present in the
  # frontmatter OR listed in skill-exceptions.json's allowed_missing.
  local allowed
  allowed=$(get_allowed_missing "$rel")
  local local_fail=0
  local field
  for field in "${SPEC_REQUIRED_EXCEPTION_BACKED[@]}"; do
    local present
    present=$(jq --arg k "$field" 'has($k)' "$json_path")
    if [[ "$present" == "true" ]]; then
      continue
    fi
    local excepted
    excepted=$(echo "$allowed" | jq --arg k "$field" 'index($k) // false')
    if [[ "$excepted" == "false" ]]; then
      emit_finding "$rel" "$field" "required-or-excepted" \
        "field is missing and not listed in schemas/skill-exceptions.json" \
        "add '${field}: ...' to frontmatter or add an exception entry with reason and spec_ref"
      local_fail=1
    fi
  done

  # 5. Static `calls`: every entry must resolve to an existing file under
  # the package skill root. The base is package-skill-root, NOT the
  # caller's directory (#301 AC6).
  local pkg_root
  pkg_root=$(package_skill_root_of "$file")
  local calls_count
  calls_count=$(jq '(.calls // []) | length' "$json_path")
  local i
  for ((i = 0; i < calls_count; i++)); do
    local target
    target=$(jq -r --argjson i "$i" '
      (.calls[$i] | if type == "string" then . else .path end)
    ' "$json_path")
    if [[ -z "$target" || "$target" == "null" ]]; then
      emit_finding "$rel" "calls[$i]" "calls-shape" \
        "entry has no resolvable path" \
        "use bare path 'sub/SKILL.md' or mapping with 'path:' per LANGUAGE-SPEC §2.4.1"
      local_fail=1
      continue
    fi
    local resolved="${pkg_root}/${target}"
    if [[ ! -e "$resolved" ]]; then
      emit_finding "$rel" "calls[$i]" "calls-target-exists" \
        "target '${target}' resolves to '${resolved#"$REPO_ROOT"/}' which does not exist" \
        "fix the path (resolution base is the package skill root, not the caller directory) or remove the entry"
      local_fail=1
    fi
  done

  return "$local_fail"
}

# Compute the "package skill root" for a skill file: the directory under
# packages/<X>/skills/ whose subtree contains every SKILL.md owned by that
# package. Per LANGUAGE-SPEC §2.4.1 example, `calls: design/SKILL.md` from
# `cnos.cdd/skills/cdd/alpha/SKILL.md` resolves to `cnos.cdd/skills/cdd/
# design/SKILL.md`, so the base is `cnos.cdd/skills/cdd/`, not the caller
# directory `cnos.cdd/skills/cdd/alpha/`.
package_skill_root_of() {
  local file="$1"
  # Match path of the form ".../packages/<pkg>/skills/" and find the
  # common parent of every SKILL.md in that subtree.
  if [[ "$file" =~ (.*)/skills/ ]]; then
    local pkg_skills="${BASH_REMATCH[1]}/skills"
    # Common ancestor of all SKILL.md in pkg_skills/.
    local common=""
    while IFS= read -r p; do
      local d
      d="$(dirname "$p")"
      if [[ -z "$common" ]]; then
        common="$d"
      else
        # strip until common is a prefix of d
        while [[ "$d" != "$common"* ]]; do
          common="$(dirname "$common")"
          [[ "$common" == "/" ]] && break
        done
      fi
    done < <(find "$pkg_skills" -name SKILL.md 2>/dev/null)
    # If common ended up as the leaf parent (single-skill package), use it.
    # If common is below pkg_skills, prefer common.
    echo "${common:-$pkg_skills}"
  else
    # Fixture or out-of-tree file: resolve relative to the file's parent.
    dirname "$file"
  fi
}

# --- modes --------------------------------------------------------------
run_all() {
  local n=0
  while IFS= read -r f; do
    n=$((n + 1))
    validate_skill_file "$f" || true
  done < <(find "$root" -name SKILL.md | sort)
  echo
  if (( findings == 0 )); then
    printf "%s%s%s %d SKILL.md validated; no findings.\n" "$GREEN" "$SYM_OK" "$RESET" "$n"
    return 0
  fi
  printf "%s%s%s %d findings across %d SKILL.md files.\n" "$RED" "$SYM_FAIL" "$RESET" "$findings" "$n" >&2
  return 1
}

run_one() {
  validate_skill_file "$1" || true
  if (( findings == 0 )); then
    printf "%s%s%s %s\n" "$GREEN" "$SYM_OK" "$RESET" "${1#"$REPO_ROOT"/}"
    return 0
  fi
  return 1
}

run_self_test() {
  local positive_fail=0 negative_fail=0

  # Positive: every fixture under valid/ MUST pass.
  echo "${CYAN}-- self-test: positive (must pass) --${RESET}"
  if [[ ! -d "$FIXTURE_VALID" ]]; then
    echo "${RED}${SYM_FAIL}${RESET} fixture dir missing: $FIXTURE_VALID" >&2
    return 2
  fi
  while IFS= read -r f; do
    findings=0
    if validate_skill_file "$f"; then
      printf "%s%s%s positive: %s\n" "$GREEN" "$SYM_OK" "$RESET" "${f#"$REPO_ROOT"/}"
    else
      printf "%s%s%s positive: %s (expected to PASS)\n" "$RED" "$SYM_FAIL" "$RESET" "${f#"$REPO_ROOT"/}" >&2
      positive_fail=$((positive_fail + 1))
    fi
  done < <(find "$FIXTURE_VALID" -name SKILL.md | sort)

  # Negative: every fixture under invalid/ MUST fail. Each fixture has a
  # sidecar `.expect` file naming a substring of the expected diagnostic;
  # the test asserts that substring appears.
  echo
  echo "${CYAN}-- self-test: negative (must fail) --${RESET}"
  if [[ ! -d "$FIXTURE_INVALID" ]]; then
    echo "${RED}${SYM_FAIL}${RESET} fixture dir missing: $FIXTURE_INVALID" >&2
    return 2
  fi
  while IFS= read -r f; do
    local expect_file="${f%.md}.expect"
    local expect=""
    [[ -f "$expect_file" ]] && expect="$(cat "$expect_file")"
    findings=0
    # Capture stdout to a temp file so we can also read findings (which is
    # mutated by emit_finding inside validate_skill_file). Using $(…)
    # would run the call in a subshell and lose the counter.
    local out="${TMPDIR}/negative.out"
    : > "$out"
    if validate_skill_file "$f" >"$out" 2>&1; then
      :  # rc=0 → did not fail
    fi
    local capture
    capture="$(cat "$out")"
    cat "$out"
    if (( findings == 0 )); then
      printf "%s%s%s negative: %s (expected to FAIL but passed)\n" "$RED" "$SYM_FAIL" "$RESET" "${f#"$REPO_ROOT"/}" >&2
      negative_fail=$((negative_fail + 1))
      continue
    fi
    if [[ -n "$expect" ]] && ! echo "$capture" | grep -qF "$expect"; then
      printf "%s%s%s negative: %s (failed but expected substring '%s' missing)\n" \
        "$RED" "$SYM_FAIL" "$RESET" "${f#"$REPO_ROOT"/}" "$expect" >&2
      negative_fail=$((negative_fail + 1))
      continue
    fi
    printf "%s%s%s negative: %s (failed as expected: '%s')\n" \
      "$GREEN" "$SYM_OK" "$RESET" "${f#"$REPO_ROOT"/}" "$expect"
  done < <(find "$FIXTURE_INVALID" -name SKILL.md | sort)

  echo
  if (( positive_fail == 0 && negative_fail == 0 )); then
    printf "%s%s%s self-test: all positive + negative fixtures behave as expected.\n" \
      "$GREEN" "$SYM_OK" "$RESET"
    return 0
  fi
  printf "%s%s%s self-test: positive_fail=%d negative_fail=%d\n" \
    "$RED" "$SYM_FAIL" "$RESET" "$positive_fail" "$negative_fail" >&2
  return 1
}

case "$mode" in
  all)        run_all ;;
  self-test)  run_self_test ;;
  file)
    [[ -n "$single_file" ]] || { echo "${RED}${SYM_FAIL}${RESET} --file requires a path" >&2; exit 2; }
    run_one "$single_file"
    ;;
esac
