#!/usr/bin/env bash
# scripts/ci/validate-repo-state.sh — validate `.cn/repo.state.json` against
# schemas/repo_state.cue (cnos#656, P1 CI validation gate).
#
# The CUE schema owns shape/type/enum/closedness constraints (including the
# structural timestamp-free and no-lock-duplication invariants — see
# schemas/repo_state.cue's header comment). This script owns discovery and
# the self-test regression suite, mirroring
# scripts/ci/validate-skill-frontmatter.sh's schema-owns-shape /
# script-owns-discovery split for schemas/skill.cue.
#
# Exit codes:
#   0  every checked file (or the self-test suite) passed
#   1  one or more validation failures
#   2  prerequisite missing (cue) — not a validation failure
#
# Usage:
#   ./scripts/ci/validate-repo-state.sh
#       Validate ./.cn/repo.state.json if present; no-op (exit 0) if absent
#       (cnos itself is not a `cn repo install` consumer repo).
#   ./scripts/ci/validate-repo-state.sh --file <path>
#       Validate a single repo.state.json path.
#   ./scripts/ci/validate-repo-state.sh --self-test
#       Run schemas/fixtures/repo-state/{valid,invalid}/*.json as the
#       built-in positive/negative regression suite.

set -euo pipefail

if [[ -n "${NO_COLOR:-}" || ! -t 1 ]]; then
  RED="" GREEN="" RESET=""
else
  RED=$'\033[0;31m' GREEN=$'\033[0;32m' RESET=$'\033[0m'
fi
SYM_OK="✓" SYM_FAIL="✗"

command -v cue >/dev/null 2>&1 || {
  echo "${RED}${SYM_FAIL}${RESET} prerequisite missing: cue" >&2
  exit 2
}

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
SCHEMA="${REPO_ROOT}/schemas/repo_state.cue"
FIXTURE_VALID="${REPO_ROOT}/schemas/fixtures/repo-state/valid"
FIXTURE_INVALID="${REPO_ROOT}/schemas/fixtures/repo-state/invalid"

[[ -f "$SCHEMA" ]] || {
  echo "${RED}${SYM_FAIL}${RESET} schema not found at: $SCHEMA" >&2
  exit 2
}

mode="default"
single_file=""
while (($#)); do
  case "$1" in
    --self-test) mode="self-test"; shift;;
    --file) mode="file"; single_file="$2"; shift 2;;
    -h|--help)
      awk '/^#!/{next} /^#/{sub(/^# ?/,""); print; next} {exit}' "$0"
      exit 0;;
    *)
      echo "${RED}${SYM_FAIL}${RESET} unknown argument: $1" >&2
      exit 2;;
  esac
done

findings=0

# Runs `cue vet` against $1 and prints its outcome. Never touches the global
# `findings` counter itself — callers decide what counts as a failure (a
# valid-fixture caller wants "did NOT vet clean" to count; an
# invalid-fixture caller wants "DID vet clean" to count).
cue_vet_ok() {
  local file="$1"
  cue vet -c -d '#RepoState' "$SCHEMA" "$file"
}

case "$mode" in
  file)
    [[ -f "$single_file" ]] || { echo "${RED}${SYM_FAIL}${RESET} no such file: $single_file" >&2; exit 2; }
    rel="${single_file#"$REPO_ROOT"/}"
    if err=$(cue_vet_ok "$single_file" 2>&1); then
      echo "${GREEN}${SYM_OK}${RESET} $rel"
    else
      echo "${RED}${SYM_FAIL}${RESET} $rel"
      echo "$err" | sed 's/^/    /'
      findings=$((findings + 1))
    fi
    ;;
  self-test)
    echo "=== valid fixtures (expect pass) ==="
    for f in "$FIXTURE_VALID"/*.json; do
      [[ -e "$f" ]] || continue
      rel="${f#"$REPO_ROOT"/}"
      if err=$(cue_vet_ok "$f" 2>&1); then
        echo "${GREEN}${SYM_OK}${RESET} $rel"
      else
        echo "${RED}${SYM_FAIL}${RESET} $rel — expected cue vet to PASS this fixture, but it failed:"
        echo "$err" | sed 's/^/    /'
        findings=$((findings + 1))
      fi
    done
    echo
    echo "=== invalid fixtures (expect reject) ==="
    for f in "$FIXTURE_INVALID"/*.json; do
      [[ -e "$f" ]] || continue
      rel="${f#"$REPO_ROOT"/}"
      if cue_vet_ok "$f" >/dev/null 2>&1; then
        echo "${RED}${SYM_FAIL}${RESET} $rel — expected cue vet to REJECT this fixture, but it passed"
        findings=$((findings + 1))
      else
        echo "${GREEN}${SYM_OK}${RESET} $rel (correctly rejected)"
      fi
    done
    ;;
  default)
    default_file="${REPO_ROOT}/.cn/repo.state.json"
    if [[ ! -f "$default_file" ]]; then
      echo "${GREEN}${SYM_OK}${RESET} no .cn/repo.state.json in this repo — nothing to validate"
      exit 0
    fi
    rel="${default_file#"$REPO_ROOT"/}"
    if err=$(cue_vet_ok "$default_file" 2>&1); then
      echo "${GREEN}${SYM_OK}${RESET} $rel"
    else
      echo "${RED}${SYM_FAIL}${RESET} $rel"
      echo "$err" | sed 's/^/    /'
      findings=$((findings + 1))
    fi
    ;;
esac

if ((findings > 0)); then
  echo
  echo "${RED}${SYM_FAIL}${RESET} $findings finding(s)."
  exit 1
fi

echo
echo "${GREEN}${SYM_OK}${RESET} all checks passed."
