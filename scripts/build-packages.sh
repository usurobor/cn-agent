#!/usr/bin/env bash
# build-packages.sh — Build package tarballs and the package index.
#
# Thin wrapper around `cn build`. The Go binary is the real builder;
# this script exists so `scripts/build-packages.sh` works as a
# conventional entry point.
#
# Usage:
#   scripts/build-packages.sh          # build
#   scripts/build-packages.sh --check  # validate only
#
# Outputs:
#   dist/packages/*.tar.gz
#   dist/packages/index.json
#   dist/packages/checksums.txt

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

# Find cn binary: built in repo root, or on PATH
if [ -x "$REPO_ROOT/cn" ]; then
  CN="$REPO_ROOT/cn"
elif command -v cn >/dev/null 2>&1; then
  CN="cn"
else
  echo "ERROR: cn binary not found. Build it first:" >&2
  echo "  cd src/go && go build -o ../../cn ./cmd/cn" >&2
  exit 1
fi

exec "$CN" build "$@"
