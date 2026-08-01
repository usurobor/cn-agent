#!/bin/sh
# cnos installer — downloads pre-built cn binary from GitHub Releases
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
#
# Override install directory:
#   BIN_DIR=/usr/bin sh install.sh
#
# If /usr/local/bin requires root:
#   curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sudo sh
#
# Select a release channel (default: stable — unchanged from before
# cnos#618; this is additive/opt-in only):
#   CNOS_CHANNEL=tooling curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
#
# "tooling" resolves the newest tooling-<date>-<sha> prerelease (published
# via .github/workflows/release.yml's workflow_dispatch with an explicit
# tag + prerelease=true input) instead of GitHub's /releases/latest
# redirect, which by definition never returns a prerelease. This ships the
# current main's cn (FSM verbs, cn repo install, ...) ahead of the next
# held feature-version cut — see cnos#618.
#
# NOTE: when piping into `sh`, put the env var on the `sh` side of the
# pipe (or `export` it first) so it is visible to the process that reads
# it — `VAR=val curl ... | sh` only sets VAR for curl, not for sh:
#   curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | CNOS_CHANNEL=tooling sh

set -e

REPO="usurobor/cnos"
BIN_DIR="${BIN_DIR:-/usr/local/bin}"
BINARY_NAME="cn"
MIN_SIZE=1000000  # 1 MB — reject obviously bad downloads

# --- UX helpers ---

# Respect NO_COLOR (https://no-color.org/) and non-TTY
if [ -n "${NO_COLOR:-}" ] || [ ! -t 1 ]; then
  RED="" GREEN="" YELLOW="" RESET=""
else
  RED='\033[0;31m' GREEN='\033[0;32m' YELLOW='\033[0;33m' RESET='\033[0m'
fi

info() { printf '%s\n' "  $1"; }
ok()   { printf "${GREEN}%s${RESET}\n" "✓ $1"; }
warn() { printf "${YELLOW}%s${RESET}\n" "⚠ $1"; }
fail() {
  printf "${RED}%s${RESET}\n" "✗ $1" >&2
  shift
  for line in "$@"; do
    printf '%s\n' "  $line" >&2
  done
  exit 1
}

# --- Cleanup trap ---
TMPFILE=""
cleanup() { [ -n "$TMPFILE" ] && rm -f "$TMPFILE"; return 0; }
trap cleanup EXIT INT TERM

# --- Prerequisites ---
if ! command -v curl >/dev/null 2>&1; then
  fail "Cannot continue — curl is required but not found" \
    "" \
    "Fix by running:" \
    "  apt install curl    (Debian/Ubuntu)" \
    "  brew install curl   (macOS)" \
    "" \
    "Then rerun:" \
    "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
fi

# --- Detect platform ---
OS="$(uname -s)"
ARCH="$(uname -m)"

case "${OS}" in
  Linux)  PLATFORM="linux" ;;
  Darwin) PLATFORM="macos" ;;
  *)
    fail "Cannot continue — no pre-built binary for ${OS}" \
      "" \
      "Supported platforms: Linux, macOS" \
      "" \
      "You can build from source instead:" \
      "  git clone https://github.com/${REPO}.git && cd cnos" \
      "  opam install . --deps-only && dune build"
    ;;
esac

case "${ARCH}" in
  x86_64)  ARCH="x64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)
    fail "Cannot continue — no pre-built binary for ${OS} ${ARCH}" \
      "" \
      "Supported architectures: x86_64, aarch64 (arm64)" \
      "" \
      "You can build from source instead:" \
      "  git clone https://github.com/${REPO}.git && cd cnos" \
      "  opam install . --deps-only && dune build"
    ;;
esac

TARGET="${PLATFORM}-${ARCH}"
ok "Detected platform: ${TARGET}"

# --- Resolve release channel ---
CNOS_CHANNEL="${CNOS_CHANNEL:-}"

case "$CNOS_CHANNEL" in
  ""|stable)
    # Default (unchanged pre-cnos#618 behavior): fetch latest release via
    # redirect, no JSON parsing. GitHub's /releases/latest redirect never
    # returns a prerelease, so this path is structurally unaffected by
    # the tooling channel below.
    LATEST=$(curl -fsSI "https://github.com/${REPO}/releases/latest" \
      | grep -i '^location:' \
      | sed -E 's|.*/tag/([^ ]+).*|\1|' \
      | tr -d '\r') || true

    if [ -z "$LATEST" ]; then
      fail "Cannot continue — could not determine latest release" \
        "" \
        "Fix by checking:" \
        "  1) Your internet connection" \
        "  2) https://github.com/${REPO}/releases has at least one release" \
        "" \
        "Then rerun:" \
        "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
    fi
    ;;
  tooling)
    # Opt-in tooling channel (cnos#618 AC2): resolve the newest
    # tooling-<date>-<sha> prerelease tag via the GitHub API. Prereleases
    # never appear at /releases/latest, so this can't reuse the redirect
    # above. tag_name values are lexicographically sortable by
    # construction (fixed 8-digit UTC date first), so we sort locally
    # rather than depend on any undocumented API response ordering.
    LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases?per_page=100" 2>/dev/null \
      | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"tooling-[^"]*"' \
      | sed -E 's/.*"(tooling-[^"]*)"$/\1/' \
      | LC_ALL=C sort -r \
      | head -1) || true

    if [ -z "$LATEST" ]; then
      fail "Cannot continue — no tooling-channel release found" \
        "" \
        "CNOS_CHANNEL=tooling resolves the newest tooling-<date>-<sha>" \
        "prerelease (published via release.yml's workflow_dispatch with" \
        "an explicit tag + prerelease=true input)." \
        "" \
        "Fix by checking:" \
        "  1) Your internet connection" \
        "  2) https://github.com/${REPO}/releases has a tooling-* prerelease" \
        "" \
        "Or use the stable channel instead:" \
        "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
    fi
    ;;
  *)
    fail "Cannot continue — unknown CNOS_CHANNEL: '${CNOS_CHANNEL}'" \
      "" \
      "Supported values: (unset, default: stable), tooling"
    ;;
esac

ok "Resolved release (channel=${CNOS_CHANNEL:-stable}): ${LATEST}"

# --- Probe BIN_DIR writability (early fail with actionable error) ---
if [ ! -d "$BIN_DIR" ]; then
  mkdir -p "$BIN_DIR" 2>/dev/null || \
    fail "Cannot continue — ${BIN_DIR} does not exist and cannot be created" \
      "" \
      "Fix by running with elevated permissions:" \
      "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sudo sh" \
      "" \
      "Or override the install directory:" \
      "  BIN_DIR=\$HOME/.local/bin curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
fi

# --- Download to temp file (same filesystem as BIN_DIR for atomic mv) ---
ARTIFACT="${BINARY_NAME}-${TARGET}"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ARTIFACT}"

TMPFILE="$(mktemp "${BIN_DIR}/.cn.XXXXXX" 2>/dev/null)" || \
  fail "Cannot continue — cannot create temp file in ${BIN_DIR}" \
    "" \
    "Fix by running with elevated permissions:" \
    "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sudo sh" \
    "" \
    "Or override the install directory:" \
    "  BIN_DIR=\$HOME/.local/bin curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"

if ! curl -fsSL -o "$TMPFILE" "$URL"; then
  fail "Cannot download binary — HTTP request failed" \
    "" \
    "URL: $URL" \
    "" \
    "Fix by running:" \
    "  1) Check your internet connection" \
    "  2) Verify the release exists: https://github.com/${REPO}/releases" \
    "" \
    "Then rerun:" \
    "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
fi

# --- Size check ---
FILE_SIZE=$(wc -c < "$TMPFILE" | tr -d ' ')
if [ "$FILE_SIZE" -lt "$MIN_SIZE" ]; then
  fail "Cannot continue — downloaded file is too small (${FILE_SIZE} bytes)" \
    "" \
    "This usually means the release asset was not found (HTML error page)" \
    "or the download was truncated." \
    "" \
    "Expected: > 1 MB for a compiled binary" \
    "URL: $URL" \
    "" \
    "Fix by verifying the release asset exists:" \
    "  https://github.com/${REPO}/releases/tag/${LATEST}" \
    "" \
    "Then rerun:" \
    "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
fi

SIZE_MB=$((FILE_SIZE / 1048576))
ok "Downloaded ${ARTIFACT} (${SIZE_MB} MB)"

# --- Checksum verification ---
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${LATEST}/checksums.txt"
CHECKSUMS="$(curl -fsSL "$CHECKSUM_URL" 2>/dev/null)" || true

if [ -n "$CHECKSUMS" ]; then
  EXPECTED=$(printf '%s\n' "$CHECKSUMS" | grep "${ARTIFACT}$" | awk '{print $1}')
  if [ -n "$EXPECTED" ]; then
    if command -v sha256sum >/dev/null 2>&1; then
      ACTUAL=$(sha256sum "$TMPFILE" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
      ACTUAL=$(shasum -a 256 "$TMPFILE" | awk '{print $1}')
    else
      warn "Cannot verify checksum — neither sha256sum nor shasum found"
      ACTUAL=""
    fi

    if [ -n "$ACTUAL" ]; then
      if [ "$ACTUAL" != "$EXPECTED" ]; then
        fail "Cannot continue — checksum mismatch (download may be corrupted)" \
          "" \
          "Expected: $EXPECTED" \
          "Actual:   $ACTUAL" \
          "" \
          "Fix by retrying:" \
          "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
      fi
      ok "Checksum verified (SHA-256)"
    fi
  else
    warn "No checksum found for ${ARTIFACT} — skipping verification"
  fi
else
  warn "Could not download checksums — skipping verification"
fi

# --- Atomic install ---
chmod +x "$TMPFILE"

if ! mv "$TMPFILE" "${BIN_DIR}/${BINARY_NAME}" 2>/dev/null; then
  fail "Cannot install — failed to write to ${BIN_DIR}" \
    "" \
    "Fix by running with elevated permissions:" \
    "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sudo sh" \
    "" \
    "Or override the install directory:" \
    "  BIN_DIR=\$HOME/.local/bin curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh"
fi
TMPFILE=""  # prevent cleanup from removing installed binary

ok "Installed to ${BIN_DIR}/${BINARY_NAME}"

# --- Verify ---
echo ""
"${BIN_DIR}/${BINARY_NAME}" --version
