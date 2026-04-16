#!/usr/bin/env bash
# R3 — Doctor (broken) kata.
#
# Proves: `cn doctor` detects broken installed-package state and
# reports a `✗` check with non-zero exit.
#
# Sequence:
#   1. build packages (if not already present at $HUB_PARENT/dist/)
#   2. create temp hub, install cnos.core
#   3. assert doctor is clean (pre-break baseline)
#   4. break an installed command entrypoint (chmod -x)
#   5. assert doctor now reports broken and exits non-zero
#
# The break is a known doctor check: command integrity per CHANGELOG
# 3.52.0 ("cn doctor validates command integrity (missing/non-exec
# entrypoints, duplicates)").

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../../lib.sh"

echo "=== R3: Doctor (broken) ==="
echo ""

require_cn

# Locate the repo's dist/packages/index.json BEFORE cd-ing away from
# the caller's starting directory. The runner is dispatched from the
# CI hub (or a developer's hub inside a repo); walking up from there
# finds the repo dist.
REPO_DIST="$(find_repo_dist "$PWD")" || {
  fail "could not find dist/packages/index.json by walking up from $PWD"
  info "hint: run 'cn build' from the repo root before invoking kata-runtime"
  kata_summary
}
info "repo dist: $REPO_DIST"

setup_temp_hub
cd "$KATA_HUB_WORK"

# Make dist/ reachable by walk-up from $KATA_HUB_WORK/cn-r3-hub.
ln -s "$REPO_DIST" "$KATA_HUB_WORK/dist"

cn init r3-hub >/dev/null 2>&1 || { fail "cn init failed"; kata_summary; }
cd cn-r3-hub
cn setup >/dev/null 2>&1 || { fail "cn setup failed"; kata_summary; }

CORE_VER="$(pkg_version_from_source cnos.core "$REPO_DIST")" || {
  fail "could not read cnos.core version from src/packages/cnos.core/cn.package.json"
  kata_summary
}
write_deps_json "cnos.core:$CORE_VER"

if cn deps lock >/dev/null 2>&1 && cn deps restore >/dev/null 2>&1; then
  pass "cnos.core installed (precondition)"
else
  fail "restore failed — cannot establish precondition"
  kata_summary
fi

ENTRYPOINT=".cn/vendor/packages/cnos.core/commands/daily/cn-daily"
if [ ! -f "$ENTRYPOINT" ]; then
  fail "expected entrypoint $ENTRYPOINT missing — precondition violated"
  kata_summary
fi

# Baseline: doctor should be clean on a freshly-restored hub. Capture
# stdout+stderr and real rc (avoid `|| true`, which masks rc).
if CLEAN_OUT="$(cn doctor 2>&1)"; then
  CLEAN_RC=0
else
  CLEAN_RC=$?
fi
if echo "$CLEAN_OUT" | grep -q '^✗'; then
  fail "pre-break doctor already reports ✗ — cannot prove break detection"
  info "pre-break output:"
  info "$CLEAN_OUT"
  kata_summary
elif [ "$CLEAN_RC" -ne 0 ]; then
  fail "pre-break doctor exited $CLEAN_RC with no ✗ — unexpected"
  info "$CLEAN_OUT"
  kata_summary
else
  pass "pre-break doctor is clean (rc=0, no ✗)"
fi

# Break: strip exec bit on the command entrypoint.
chmod -x "$ENTRYPOINT"
info "broke: chmod -x $ENTRYPOINT"

# Assert doctor catches it.
if BROKEN_OUT="$(cn doctor 2>&1)"; then
  BROKEN_RC=0
else
  BROKEN_RC=$?
fi

if [ "$BROKEN_RC" -ne 0 ]; then
  pass "post-break cn doctor exited non-zero ($BROKEN_RC)"
else
  fail "post-break cn doctor exited 0 — broken state invisible"
  info "output: $BROKEN_OUT"
fi

if echo "$BROKEN_OUT" | grep -q '^✗'; then
  pass "post-break cn doctor reports ✗"
else
  fail "post-break cn doctor emits no ✗ marker"
  info "output: $BROKEN_OUT"
fi

echo ""
kata_summary
