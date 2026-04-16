#!/usr/bin/env bash
# lib.sh — shared helpers for Tier 2 (runtime/package) kata scripts.
#
# These kata run AFTER at least one package is installed. Helpers are
# small on purpose: Tier 2 proves post-install behavior, not scoring.

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

FAILURES=0
SKIPS=0
PASSES=0

pass() { echo -e "${GREEN}PASS${NC}: $1"; PASSES=$((PASSES + 1)); }
fail() { echo -e "${RED}FAIL${NC}: $1"; FAILURES=$((FAILURES + 1)); }
skip() { echo -e "${YELLOW}SKIP${NC}: $1"; SKIPS=$((SKIPS + 1)); }
info() { echo "     $1"; }

# require_cn aborts the kata if `cn` is not on PATH.
require_cn() {
  if ! command -v cn &>/dev/null; then
    echo "ERROR: cn binary not found in PATH"
    exit 1
  fi
}

# setup_temp_hub creates an isolated workdir under $TMPDIR, registers
# a cleanup trap, sets $KATA_HUB_WORK, and cd's into it. The caller
# runs `cn init <name>` from here, then cd's into cn-<name>/.
setup_temp_hub() {
  KATA_HUB_WORK=$(mktemp -d "${TMPDIR:-/tmp}/kata-tier2-XXXXXX")
  trap 'rm -rf "$KATA_HUB_WORK"' EXIT
  info "temp workdir: $KATA_HUB_WORK"
}

# find_repo_dist echoes the absolute path to the repo's `dist/`
# directory (the one that contains `packages/index.json`). It walks up
# from the given directory (default: the caller's starting $PWD) in the
# same spirit as `cn deps restore`'s FindIndexPath. Returns non-zero
# with no output if not found.
find_repo_dist() {
  local start="${1:-$PWD}"
  local dir
  dir="$(cd "$start" && pwd)"
  while [ "$dir" != "/" ]; do
    if [ -f "$dir/dist/packages/index.json" ]; then
      echo "$dir/dist"
      return 0
    fi
    dir="$(dirname "$dir")"
  done
  return 1
}

# write_deps_json writes a minimal deps.json into $PWD/.cn/deps.json.
# Caller must be inside a hub (post-`cn init`). Takes "pkg:version" args.
#
# Schema is cn.deps.v1: `packages` is an ARRAY of {name, version}
# objects (see src/go/internal/pkg/pkg.go and src/ocaml/lib/cn_package.ml).
# Pre-#250 the lockfile dumped the entire index regardless of manifest
# content, so a malformed (object-shaped) packages field still produced
# an installable lockfile. Now `cn deps lock` honors the manifest, so
# the JSON shape must match the parser.
write_deps_json() {
  local out=".cn/deps.json"
  {
    echo '{'
    echo '  "schema": "cn.deps.v1",'
    echo '  "profile": "engineer",'
    echo '  "packages": ['
    local first=1
    for spec in "$@"; do
      local name="${spec%%:*}"
      local ver="${spec##*:}"
      if [ "$first" -eq 1 ]; then first=0; else echo '    ,'; fi
      echo "    {\"name\": \"$name\", \"version\": \"$ver\"}"
    done
    echo '  ]'
    echo '}'
  } > "$out"
}

# pkg_version_from_source echoes the version declared in a package's
# source manifest. Args: <pkg_name> <repo_dist_path>. The repo source
# lives one directory above $REPO_DIST. Returns non-zero with no output
# if the manifest cannot be read or the version field is absent.
#
# Used by katas that need to pin to "whatever version cn build just
# produced" — hardcoded versions go stale and break silently when the
# parent repo bumps src/packages/<pkg>/cn.package.json.
pkg_version_from_source() {
  local pkg="$1"
  local dist="$2"
  local repo_root
  repo_root="$(dirname "$dist")"
  local manifest="$repo_root/src/packages/$pkg/cn.package.json"
  [ -f "$manifest" ] || return 1
  local ver
  ver=$(grep -m1 '"version"' "$manifest" \
    | sed -E 's/.*"version":[[:space:]]*"([^"]+)".*/\1/')
  [ -n "$ver" ] || return 1
  echo "$ver"
}

# kata_summary prints a one-line pass/fail summary and exits non-zero
# if any assertion failed. Called at the end of every kata script.
kata_summary() {
  echo ""
  if [ "$FAILURES" -gt 0 ]; then
    echo -e "${RED}$FAILURES failure(s)${NC}, $PASSES pass(es), $SKIPS skip(s)"
    exit 1
  elif [ "$SKIPS" -gt 0 ]; then
    echo -e "${GREEN}All passed${NC} ($PASSES pass(es), $SKIPS skip(s))"
  else
    echo -e "${GREEN}All passed${NC} ($PASSES pass(es))"
  fi
}
