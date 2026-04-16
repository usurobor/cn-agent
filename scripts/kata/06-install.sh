#!/usr/bin/env bash
# 06-install.sh — Tier 1 kata 06.
#
# Proves: `cn deps restore` installs packages from dist/.
# Pass condition: at least one package appears under .cn/vendor/packages/
# with a cn.package.json.
#
# Depends on 05-build.sh: this kata does not build artifacts itself; it
# requires $REPO_ROOT/dist/packages/index.json to already exist. run-all.sh
# runs 05 before 06 in filename order, which satisfies the dependency.
# Running 06 in isolation without first running 05 (or `cn build`) will
# fail fast with a clear message.
#
# restore.FindIndexPath walks up from the hub looking for
# dist/packages/index.json, so the temp hub is created under a workdir
# that symlinks dist/ from the repo root.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/lib.sh"

echo "=== Kata 06: Install ==="
echo ""

require_cn

REPO_ROOT="$(repo_root)"
if [ ! -f "$REPO_ROOT/dist/packages/index.json" ]; then
  fail "dist/packages/index.json missing — run 05-build.sh first"
  kata_summary
fi

setup_temp_hub
# Expose dist/ to FindIndexPath: the hub will be at
# $KATA_HUB/cn-kata-hub/, walking up reaches $KATA_HUB/dist/.
ln -s "$REPO_ROOT/dist" "$KATA_HUB/dist"
cd "$KATA_HUB"

cn init kata-hub >/dev/null 2>&1 || { fail "cn init failed (precondition)"; kata_summary; }
cd cn-kata-hub

cn setup >/dev/null 2>&1 || { fail "cn setup failed (precondition)"; kata_summary; }

# `cn setup` writes a default deps.json pinning cnos.core/cnos.eng to
# the binary version. In CI the binary is built without -ldflags, so
# the version is "dev" — not a real key in the index. Overwrite with
# a manifest pinned to the actually-built version of cnos.core (read
# from the package source, the single authority for package versions).
# This keeps the kata focused on lock/restore rather than version
# negotiation between binary build and package source.
CORE_VER=$(grep -m1 '"version"' "$REPO_ROOT/src/packages/cnos.core/cn.package.json" \
  | sed -E 's/.*"version":[[:space:]]*"([^"]+)".*/\1/')
if [ -z "$CORE_VER" ]; then
  fail "could not extract cnos.core version from src/packages/cnos.core/cn.package.json"
  kata_summary
fi
cat > .cn/deps.json <<JSON
{
  "schema": "cn.deps.v1",
  "profile": "engineer",
  "packages": [
    {"name": "cnos.core", "version": "$CORE_VER"}
  ]
}
JSON
info "deps.json pinned to cnos.core@$CORE_VER"

info "running: cn deps lock"
if cn deps lock >/dev/null 2>&1; then
  pass "cn deps lock exits 0"
else
  fail "cn deps lock failed"
  kata_summary
fi

if [ -f ".cn/deps.lock.json" ]; then
  pass ".cn/deps.lock.json written"
else
  fail ".cn/deps.lock.json missing after lock"
  kata_summary
fi

info "running: cn deps restore"
if cn deps restore >/dev/null 2>&1; then
  pass "cn deps restore exits 0"
else
  fail "cn deps restore failed"
  kata_summary
fi

# Count installed packages with a cn.package.json manifest.
INSTALLED=0
INSTALLED_NAMES=""
if [ -d ".cn/vendor/packages" ]; then
  for d in .cn/vendor/packages/*/; do
    if [ -f "${d}cn.package.json" ]; then
      INSTALLED=$((INSTALLED + 1))
      INSTALLED_NAMES="$INSTALLED_NAMES $(basename "$d")"
    fi
  done
fi

# Per #250: lock+restore must install only what deps.json pinned.
# We pinned exactly cnos.core, so exactly one package must be installed.
if [ "$INSTALLED" -eq 1 ] && [ -f ".cn/vendor/packages/cnos.core/cn.package.json" ]; then
  pass "exactly cnos.core installed (deps.json pin honored)"
else
  fail "expected exactly cnos.core installed, got $INSTALLED package(s):$INSTALLED_NAMES"
fi

echo ""
kata_summary
