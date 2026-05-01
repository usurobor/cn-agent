#!/usr/bin/env bash
# R5 — Activate kata.
#
# Proves: `cn activate` generates a bootstrap prompt from local hub
# state, with the Kernel/Persona/Operator triad introduced in 3.71.0
# (#321) and the structural guarantees described in kata.md (P1–P11).
#
# Each pass condition builds its own fixture under an isolated temp
# workdir so the assertions are independent. Fixtures use `cn init`
# where possible and manual files where the kata stipulates a specific
# shape (sigma-shape, manifest-only, none).

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../../lib.sh"

echo "=== R5: Activate ==="
echo ""

require_cn

# Activate does not need installed packages, but the manifest-only and
# vendored fixtures need to know a real cnos.core version. Resolve via
# the repo dist (same walk-up the other R-katas use).
REPO_DIST="$(find_repo_dist "$PWD")" || {
  fail "could not find dist/packages/index.json by walking up from $PWD"
  info "hint: run 'cn build' from the repo root before invoking kata-runtime"
  kata_summary
}
info "repo dist: $REPO_DIST"

CORE_VER="$(pkg_version_from_source cnos.core "$REPO_DIST")" || {
  fail "could not read cnos.core version from src/packages/cnos.core/cn.package.json"
  kata_summary
}
info "cnos.core version: $CORE_VER"

setup_temp_hub
cd "$KATA_HUB_WORK"

# ---- Build a freshly-initialized hub for P1–P6. ----
cn init r5-hub >/dev/null 2>&1 || { fail "cn init failed (precondition)"; kata_summary; }
HUB_DIR="$KATA_HUB_WORK/cn-r5-hub"

# --- P1: cwd hub discovery ---
info "P1: cn activate from inside hub"
cd "$HUB_DIR"
P1_OUT="$(cn activate 2>/dev/null)"; P1_RC=$?
if [ "$P1_RC" -eq 0 ]; then
  pass "P1: cn activate exits 0 from inside hub"
else
  fail "P1: cn activate exited $P1_RC from inside hub"
fi
if echo "$P1_OUT" | grep -qF "You are activating a cnos hub."; then
  pass "P1: stdout contains activation header"
else
  fail "P1: stdout missing activation header"
fi
if echo "$P1_OUT" | grep -qF "$HUB_DIR"; then
  pass "P1: stdout contains hub path"
else
  fail "P1: stdout missing hub path"
fi

# --- P2: explicit HUB_DIR from outside ---
info "P2: cn activate <hub> from outside"
cd /tmp
P2_OUT="$(cn activate "$HUB_DIR" 2>/dev/null)"; P2_RC=$?
if [ "$P2_RC" -eq 0 ] \
   && echo "$P2_OUT" | grep -qF "You are activating a cnos hub." \
   && echo "$P2_OUT" | grep -qF "$HUB_DIR"; then
  pass "P2: explicit HUB_DIR from /tmp emits prompt with hub path"
else
  fail "P2: explicit HUB_DIR did not produce expected prompt (rc=$P2_RC)"
fi

# --- P3: stdout/stderr separation ---
info "P3: stdout is prompt-only"
cd /tmp
cn activate "$HUB_DIR" >"$KATA_HUB_WORK/p3.out" 2>"$KATA_HUB_WORK/p3.err"
if grep -qE '^(→|✓|✗|ERROR|WARN|Generating)' "$KATA_HUB_WORK/p3.out"; then
  fail "P3: stdout contains diagnostic markers (should go to stderr)"
  info "stdout head: $(head -3 "$KATA_HUB_WORK/p3.out")"
else
  pass "P3: stdout carries prompt only (no diagnostic markers)"
fi

# --- P4: missing hub fails clearly ---
info "P4: no hub, no arg → non-zero with diagnostic"
cd /tmp
if P4_OUT="$(cn activate 2>"$KATA_HUB_WORK/p4.err")"; then
  P4_RC=0
else
  P4_RC=$?
fi
if [ "$P4_RC" -ne 0 ]; then
  pass "P4: cn activate exits non-zero outside a hub"
else
  fail "P4: cn activate succeeded with no hub (rc=0)"
fi
if [ -z "$P4_OUT" ]; then
  pass "P4: stdout empty on failure"
else
  fail "P4: stdout non-empty on failure"
fi
if grep -qiE 'hub' "$KATA_HUB_WORK/p4.err"; then
  pass "P4: stderr mentions hub discovery"
else
  fail "P4: stderr lacks hub-discovery diagnostic"
fi

# --- P5: explicit bad path fails ---
info "P5: explicit bad path → non-zero with diagnostic"
if P5_OUT="$(cn activate /nonexistent/r5/path 2>"$KATA_HUB_WORK/p5.err")"; then
  P5_RC=0
else
  P5_RC=$?
fi
if [ "$P5_RC" -ne 0 ] && [ -z "$P5_OUT" ] && [ -s "$KATA_HUB_WORK/p5.err" ]; then
  pass "P5: bad path → non-zero, empty stdout, diagnostic on stderr"
else
  fail "P5: bad path did not fail cleanly (rc=$P5_RC, stdout=${#P5_OUT}b, stderr=$(stat -c %s "$KATA_HUB_WORK/p5.err" 2>/dev/null || echo 0)b)"
fi

# --- P6: no secrets in prompt ---
info "P6: secrets are not surfaced in prompt"
SECRET_TOKEN="r5-kata-secret-$(date +%s%N)"
mkdir -p "$HUB_DIR/.cn"
echo "API_KEY=$SECRET_TOKEN" > "$HUB_DIR/.cn/secrets.env"
echo "DOTENV_KEY=$SECRET_TOKEN-2" > "$HUB_DIR/.env"
P6_OUT="$(cn activate "$HUB_DIR" 2>/dev/null)"
if echo "$P6_OUT" | grep -qF "$SECRET_TOKEN"; then
  fail "P6: prompt leaks .cn/secrets.env contents"
elif echo "$P6_OUT" | grep -qF "$SECRET_TOKEN-2"; then
  fail "P6: prompt leaks .env contents"
else
  pass "P6: prompt does not contain secrets file contents"
fi
rm -f "$HUB_DIR/.cn/secrets.env" "$HUB_DIR/.env"

# ---- Sigma-shape fixture for P7, P8a, P9a, P10, P11 ----
cd "$KATA_HUB_WORK"
SIGMA_DIR="$KATA_HUB_WORK/cn-sigma-hub"
cn init sigma-hub >/dev/null 2>&1 || { fail "cn init sigma-hub failed"; kata_summary; }
mkdir -p "$SIGMA_DIR/spec" \
         "$SIGMA_DIR/.cn/vendor/packages/cnos.core/doctrine" \
         "$SIGMA_DIR/threads/reflections/daily"
echo "# Persona" > "$SIGMA_DIR/spec/PERSONA.md"
echo "# Operator" > "$SIGMA_DIR/spec/OPERATOR.md"
echo "# Kernel" > "$SIGMA_DIR/.cn/vendor/packages/cnos.core/doctrine/KERNEL.md"
cat > "$SIGMA_DIR/.cn/vendor/packages/cnos.core/cn.package.json" <<JSON
{"name": "cnos.core", "version": "$CORE_VER"}
JSON
cat > "$SIGMA_DIR/.cn/deps.json" <<JSON
{
  "schema": "cn.deps.v1",
  "profile": "engineer",
  "packages": [
    {"name": "cnos.core", "version": "$CORE_VER"}
  ]
}
JSON
echo "# 2026-04-29" > "$SIGMA_DIR/threads/reflections/daily/2026-04-29.md"
echo "# 2026-05-01" > "$SIGMA_DIR/threads/reflections/daily/2026-05-01.md"

SIGMA_OUT="$(cn activate "$SIGMA_DIR" 2>/dev/null)"

# --- P7: triad split ---
info "P7: ## Kernel / ## Persona / ## Operator; no ## Identity"
if echo "$SIGMA_OUT" | grep -qF "## Kernel" \
   && echo "$SIGMA_OUT" | grep -qF "## Persona" \
   && echo "$SIGMA_OUT" | grep -qF "## Operator"; then
  pass "P7: prompt contains Kernel/Persona/Operator sections"
else
  fail "P7: prompt missing one of Kernel/Persona/Operator sections"
fi
if echo "$SIGMA_OUT" | grep -qF "## Identity"; then
  fail "P7: prompt still contains legacy ## Identity bucket"
else
  pass "P7: prompt does not contain ## Identity"
fi
if echo "$SIGMA_OUT" | grep -qF "no identity files found"; then
  fail "P7: prompt still emits 'no identity files found'"
else
  pass "P7: prompt does not emit 'no identity files found'"
fi

# --- P8a: kernel vendored ---
info "P8a: vendored kernel"
if echo "$SIGMA_OUT" | grep -qE "vendored at .*doctrine/KERNEL\.md@$CORE_VER"; then
  pass "P8a: ## Kernel reports vendored at path@$CORE_VER"
else
  fail "P8a: ## Kernel missing 'vendored at … @$CORE_VER'"
fi
# stdout should not suggest cn deps restore when kernel is vendored
KERNEL_BLOCK="$(echo "$SIGMA_OUT" | awk '/^## Kernel/{flag=1; next} /^## /{flag=0} flag')"
if echo "$KERNEL_BLOCK" | grep -qF "cn deps restore"; then
  fail "P8a: vendored ## Kernel block suggests 'cn deps restore'"
else
  pass "P8a: vendored ## Kernel block does not suggest 'cn deps restore'"
fi

# --- P8b: kernel manifest-only ---
info "P8b: manifest-only kernel"
cd "$KATA_HUB_WORK"
MANIFEST_DIR="$KATA_HUB_WORK/cn-manifest-hub"
cn init manifest-hub >/dev/null 2>&1 || { fail "cn init manifest-hub failed"; kata_summary; }
cat > "$MANIFEST_DIR/.cn/deps.json" <<JSON
{
  "schema": "cn.deps.v1",
  "profile": "engineer",
  "packages": [
    {"name": "cnos.core", "version": "$CORE_VER"}
  ]
}
JSON
MANIFEST_OUT="$(cn activate "$MANIFEST_DIR" 2>/dev/null)"
if echo "$MANIFEST_OUT" | grep -qF "dependency manifest declares cnos.core; not restored — run cn deps restore"; then
  pass "P8b: ## Kernel reports manifest-only state"
else
  fail "P8b: ## Kernel missing manifest-only diagnostic"
fi
if echo "$MANIFEST_OUT" | grep -qF "vendored at"; then
  fail "P8b: manifest-only prompt incorrectly contains 'vendored at'"
else
  pass "P8b: manifest-only prompt does not contain 'vendored at'"
fi

# --- P8c: kernel none ---
info "P8c: no kernel reference"
cd "$KATA_HUB_WORK"
NONE_DIR="$KATA_HUB_WORK/cn-none-hub"
cn init none-hub >/dev/null 2>&1 || { fail "cn init none-hub failed"; kata_summary; }
NONE_OUT="$(cn activate "$NONE_DIR" 2>/dev/null)"
NONE_KERNEL="$(echo "$NONE_OUT" | awk '/^## Kernel/{flag=1; next} /^## /{flag=0} flag')"
if echo "$NONE_KERNEL" | grep -qF "no kernel reference"; then
  pass "P8c: ## Kernel reports 'no kernel reference'"
else
  fail "P8c: ## Kernel missing 'no kernel reference'"
fi
if echo "$NONE_KERNEL" | grep -qF "cn deps restore"; then
  fail "P8c: no-kernel ## Kernel block suggests 'cn deps restore'"
else
  pass "P8c: no-kernel ## Kernel block does not suggest 'cn deps restore'"
fi

# --- P9a: deps restored ---
info "P9a: deps restored"
DEPS_BLOCK_SIGMA="$(echo "$SIGMA_OUT" | awk '/^## Dependencies/{flag=1; next} /^## /{flag=0} flag')"
if echo "$DEPS_BLOCK_SIGMA" | grep -qE "restored:.*cnos\.core"; then
  pass "P9a: ## Dependencies lists restored cnos.core"
else
  fail "P9a: ## Dependencies missing 'restored: …cnos.core'"
fi

# --- P9b: deps manifest-only ---
info "P9b: deps manifest-only"
DEPS_BLOCK_MANIFEST="$(echo "$MANIFEST_OUT" | awk '/^## Dependencies/{flag=1; next} /^## /{flag=0} flag')"
if echo "$DEPS_BLOCK_MANIFEST" | grep -qF "dependency manifest present; packages not restored — run cn deps restore"; then
  pass "P9b: ## Dependencies reports manifest-only state"
else
  fail "P9b: ## Dependencies missing manifest-only diagnostic"
fi

# --- P9c: deps none ---
info "P9c: no deps manifest"
DEPS_BLOCK_NONE="$(echo "$NONE_OUT" | awk '/^## Dependencies/{flag=1; next} /^## /{flag=0} flag')"
if echo "$DEPS_BLOCK_NONE" | grep -qF "no dependency manifest"; then
  pass "P9c: ## Dependencies reports 'no dependency manifest'"
else
  fail "P9c: ## Dependencies missing 'no dependency manifest'"
fi

# --- P10: read-first ordering ---
info "P10: ## Read first ordering (PERSONA → OPERATOR → KERNEL → deps → reflection)"
READ_BLOCK="$(echo "$SIGMA_OUT" | awk '/^## Read first/{flag=1; next} /^## /{flag=0} flag')"
if [ -n "$READ_BLOCK" ]; then
  pass "P10: ## Read first section present"
else
  fail "P10: ## Read first section missing"
fi
PERSONA_LINE=$(echo "$READ_BLOCK" | grep -nF "PERSONA.md" | head -1 | cut -d: -f1)
OPERATOR_LINE=$(echo "$READ_BLOCK" | grep -nF "OPERATOR.md" | head -1 | cut -d: -f1)
KERNEL_LINE=$(echo "$READ_BLOCK"   | grep -nF "KERNEL.md"   | head -1 | cut -d: -f1)
DEPS_LINE=$(echo "$READ_BLOCK"     | grep -nF "deps.json"   | head -1 | cut -d: -f1)
REFL_LINE=$(echo "$READ_BLOCK"     | grep -nF "reflection"  | head -1 | cut -d: -f1)
if [ -n "$PERSONA_LINE" ] && [ -n "$OPERATOR_LINE" ] && [ -n "$KERNEL_LINE" ] \
   && [ -n "$DEPS_LINE" ] && [ -n "$REFL_LINE" ] \
   && [ "$PERSONA_LINE" -lt "$OPERATOR_LINE" ] \
   && [ "$OPERATOR_LINE" -lt "$KERNEL_LINE" ] \
   && [ "$KERNEL_LINE" -lt "$DEPS_LINE" ] \
   && [ "$DEPS_LINE" -lt "$REFL_LINE" ]; then
  pass "P10: read-first order is persona < operator < kernel < deps < reflection"
else
  fail "P10: read-first order incorrect (persona=$PERSONA_LINE operator=$OPERATOR_LINE kernel=$KERNEL_LINE deps=$DEPS_LINE refl=$REFL_LINE)"
fi

# --- P11a: latest reflection pointer ---
info "P11a: latest reflection points at most-recent file"
if echo "$SIGMA_OUT" | grep -qF "2026-05-01.md"; then
  pass "P11a: ## Read first names the latest reflection (2026-05-01.md)"
else
  fail "P11a: ## Read first does not name 2026-05-01.md as latest"
fi
LATEST_LINE=$(echo "$SIGMA_OUT" | grep -nF "latest reflection" | head -1 | cut -d: -f1)
if [ -n "$LATEST_LINE" ]; then
  CONTEXT=$(echo "$SIGMA_OUT" | sed -n "${LATEST_LINE}p")
  if echo "$CONTEXT" | grep -qF "2026-04-29.md"; then
    fail "P11a: 2026-04-29.md (older) listed on the latest-reflection line"
  else
    pass "P11a: older reflection not listed as latest"
  fi
fi

# --- P11b: latest reflection absent → omitted, not 'not present' ---
info "P11b: empty reflections directory → no 'not present' marker"
cd "$KATA_HUB_WORK"
P11B_DIR="$KATA_HUB_WORK/cn-p11b-hub"
cn init p11b-hub >/dev/null 2>&1 || { fail "cn init p11b-hub failed"; kata_summary; }
rm -rf "$P11B_DIR/threads/reflections/daily"/*
P11B_OUT="$(cn activate "$P11B_DIR" 2>/dev/null)"
P11B_READ="$(echo "$P11B_OUT" | awk '/^## Read first/{flag=1; next} /^## /{flag=0} flag')"
if echo "$P11B_READ" | grep -qE "reflection.*not present|not present.*reflection"; then
  fail "P11b: read-first section says 'not present' for absent reflection"
else
  pass "P11b: absent reflection is omitted, not flagged 'not present'"
fi

echo ""
kata_summary
