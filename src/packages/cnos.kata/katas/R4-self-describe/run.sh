#!/usr/bin/env bash
# R4 — Self-describe kata.
#
# Proves: `cn status` surfaces installed packages and discovered
# commands, not just the empty hub skeleton. The 3.53.0 status rewrite
# made this machine-testable.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../../lib.sh"

echo "=== R4: Self-describe ==="
echo ""

require_cn

REPO_DIST="$(find_repo_dist "$PWD")" || {
  fail "could not find dist/packages/index.json by walking up from $PWD"
  info "hint: run 'cn build' from the repo root before invoking kata-runtime"
  kata_summary
}
info "repo dist: $REPO_DIST"

setup_temp_hub
cd "$KATA_HUB_WORK"

ln -s "$REPO_DIST" "$KATA_HUB_WORK/dist"

cn init r4-hub >/dev/null 2>&1 || { fail "cn init failed"; kata_summary; }
cd cn-r4-hub
cn setup >/dev/null 2>&1 || { fail "cn setup failed"; kata_summary; }

CORE_VER="$(pkg_version_from_source cnos.core "$REPO_DIST")" || {
  fail "could not read cnos.core version from src/packages/cnos.core/cn.package.json"
  kata_summary
}
write_deps_json "cnos.core:$CORE_VER"

if cn deps lock >/dev/null 2>&1 && cn deps restore >/dev/null 2>&1; then
  pass "cnos.core installed (precondition)"
else
  fail "restore failed"
  kata_summary
fi

info "running: cn status"
STATUS_OUT=$(cn status 2>&1)
STATUS_RC=$?

if [ "$STATUS_RC" -eq 0 ]; then
  pass "cn status exits 0"
else
  fail "cn status exited $STATUS_RC"
  info "output: $STATUS_OUT"
  kata_summary
fi

if echo "$STATUS_OUT" | grep -qi 'cn hub'; then
  pass "cn status shows hub identity"
else
  fail "cn status omits hub identity"
fi

if echo "$STATUS_OUT" | grep -q 'cnos.core'; then
  pass "cn status lists cnos.core (installed package name)"
else
  fail "cn status omits cnos.core"
  info "output: $STATUS_OUT"
fi

if echo "$STATUS_OUT" | grep -qE '(^|\s)daily(\s|$)'; then
  pass "cn status surfaces 'daily' command from installed package"
else
  fail "cn status omits the 'daily' command"
  info "output: $STATUS_OUT"
fi

echo ""
kata_summary
