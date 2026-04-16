#!/usr/bin/env bash
# R2 — Round-trip kata.
#
# Proves: author → build → install → dispatch in a hub freshly set up
# from nothing but the `cn` binary.
#
# Isolated: creates a tempdir, authors a fixture package, builds and
# installs there. Does not touch the hub/repo the runner was invoked
# from. Cleans up the tempdir on exit (including on failure).

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../../lib.sh"

echo "=== R2: Round-trip ==="
echo ""

require_cn

WORKDIR=$(mktemp -d "${TMPDIR:-/tmp}/kata-r2-XXXXXX")
trap 'rm -rf "$WORKDIR"' EXIT
info "workdir: $WORKDIR"

# `cn build` walks up from cwd looking for .git to locate the repo
# root. Make the workdir itself a git repo so discovery stays local.
( cd "$WORKDIR" && git init -q ) || { fail "git init failed"; kata_summary; }

# 1. Author the fixture package.
FIXTURE="$WORKDIR/src/packages/kata-rt-fixture"
mkdir -p "$FIXTURE/commands/rt-test"

cat > "$FIXTURE/cn.package.json" <<'JSON'
{
  "schema": "cn.package.v1",
  "name": "kata-rt-fixture",
  "version": "0.1.0",
  "kind": "package",
  "engines": { "cnos": ">=3.50.0" },
  "commands": {
    "kata-rt-test": {
      "entrypoint": "commands/rt-test/cn-rt-test",
      "summary": "Round-trip fixture: emits a marker with dispatch env vars"
    }
  }
}
JSON

cat > "$FIXTURE/commands/rt-test/cn-rt-test" <<'SH'
#!/usr/bin/env bash
echo "ROUNDTRIP_OK hub=${CN_HUB_PATH:-} pkg=${CN_PACKAGE_ROOT:-} cmd=${CN_COMMAND_NAME:-}"
SH
chmod +x "$FIXTURE/commands/rt-test/cn-rt-test"

# 2. Build.
info "running: cn build (in $WORKDIR)"
cd "$WORKDIR"
if cn build >/dev/null 2>&1; then
  pass "cn build succeeded"
else
  fail "cn build failed"
  kata_summary
fi

if [ -f "$WORKDIR/dist/packages/index.json" ]; then
  pass "dist/packages/index.json produced"
else
  fail "dist/packages/index.json missing"
  kata_summary
fi

TARBALL="$(find "$WORKDIR/dist/packages" -maxdepth 1 -name 'kata-rt-fixture-*.tar.gz' | head -1)"
if [ -n "$TARBALL" ] && [ -f "$TARBALL" ]; then
  pass "fixture tarball produced ($(basename "$TARBALL"))"
else
  fail "fixture tarball missing"
  kata_summary
fi

# 3. Create a sub-hub under $WORKDIR so FindIndexPath walks up to
# $WORKDIR/dist/packages/index.json.
cd "$WORKDIR"
cn init rt-hub >/dev/null 2>&1 || { fail "cn init failed"; kata_summary; }
cd cn-rt-hub
cn setup >/dev/null 2>&1 || { fail "cn setup failed"; kata_summary; }

write_deps_json "kata-rt-fixture:0.1.0"

info "running: cn deps lock"
if cn deps lock >/dev/null 2>&1; then
  pass "cn deps lock exits 0"
else
  fail "cn deps lock failed"
  kata_summary
fi

info "running: cn deps restore"
if cn deps restore >/dev/null 2>&1; then
  pass "cn deps restore exits 0"
else
  fail "cn deps restore failed"
  kata_summary
fi

if [ -f ".cn/vendor/packages/kata-rt-fixture/cn.package.json" ]; then
  pass "kata-rt-fixture installed under .cn/vendor/packages/"
else
  fail "kata-rt-fixture not present in .cn/vendor/packages/"
  kata_summary
fi

# 4. Dispatch.
info "running: cn kata rt-test"
RT_OUT=$(cn kata rt-test 2>&1 || true)
if echo "$RT_OUT" | grep -q 'ROUNDTRIP_OK'; then
  pass "cn kata rt-test dispatched and produced marker"
else
  if echo "$RT_OUT" | grep -qi 'unknown command'; then
    fail "cn kata rt-test: unknown command — round-trip broken"
  else
    fail "cn kata rt-test: unexpected output"
  fi
  info "output: $RT_OUT"
  kata_summary
fi

# Env vars on dispatch. Both checks must be able to fail (else: fail
# branch); the regex requires a non-empty value (\S after `=`) so that
# an unset/empty CN_HUB_PATH cannot tautologically satisfy the grep.
# Review F2 (PR #248): bare `hub=` matched regardless of value because
# the fixture echoes `hub=${CN_HUB_PATH:-}`; tighten to require content.
if echo "$RT_OUT" | grep -qE 'hub=\S'; then
  pass "CN_HUB_PATH passed to command (non-empty)"
else
  fail "CN_HUB_PATH absent or empty in dispatched command output"
  info "output: $RT_OUT"
fi
if echo "$RT_OUT" | grep -q 'cmd=kata-rt-test'; then
  pass "CN_COMMAND_NAME passed to command"
else
  fail "CN_COMMAND_NAME absent or wrong in dispatched command output"
  info "output: $RT_OUT"
fi

echo ""
kata_summary
