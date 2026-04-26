#!/usr/bin/env bash
# 90-release-bootstrap.sh — Prove a freshly released cn binary can
# bootstrap a hub from production endpoints.
#
# Issue #230 AC3–AC5: the encoding-leak chain (lockfile → vendor →
# released binary) was last caught by accident (packages/index.json
# migration leak in the v3.x era). This smoke is the named test that
# would have caught it.
#
# Steps:
#   1. detect host platform (linux-x64 | linux-arm64 | macos-x64 | macos-arm64)
#   2. resolve release tag (arg or latest from GitHub API)
#   3. download cn-<platform>, packages/index.json, and every referenced
#      tarball from the GitHub release
#   4. verify each tarball's SHA-256 against the index.json entry — the
#      same authority `cn deps restore` itself trusts at runtime
#   5. cn init scratch-hub  → cn setup  → cn deps lock → cn deps restore
#   6. assert each pinned package is present under .cn/vendor/packages/<name>/
#      with a cn.package.json whose version matches the lockfile pin
#
# Usage:
#   scripts/smoke/90-release-bootstrap.sh              # latest release
#   scripts/smoke/90-release-bootstrap.sh 3.58.0       # specific tag
#
# Exit:
#   0 — bootstrap chain works end-to-end against the released artifacts
#   1 — failure (chain is broken: missing asset, sha mismatch, install error)
#   2 — skipped (offline, no curl, or no release available — not a failure)
#
# Conventions: eng/tool/SKILL.md (set -euo pipefail, NO_COLOR support,
# prereq checks, idempotent within its own temp dir, machine-readable
# trailing summary line).

set -euo pipefail

REPO="usurobor/cnos"
TAG="${1:-}"

# --- color (NO_COLOR support) ---
if [[ -n "${NO_COLOR:-}" ]] || [[ ! -t 1 ]]; then
  RED="" GREEN="" YELLOW="" BLUE="" RESET=""
else
  RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[0;33m'
  BLUE='\033[0;34m'; RESET='\033[0m'
fi

info()  { printf "%b→%b %s\n"  "$BLUE"   "$RESET" "$*"; }
pass()  { printf "%b✓%b %s\n"  "$GREEN"  "$RESET" "$*"; }
warn()  { printf "%b!%b %s\n"  "$YELLOW" "$RESET" "$*" >&2; }
fail()  { printf "%b✗%b %s\n"  "$RED"    "$RESET" "$*" >&2; }
skip()  { printf "%b○%b %s\n"  "$YELLOW" "$RESET" "$*" >&2; }

# --- prereqs ---
need() {
  command -v "$1" >/dev/null 2>&1 || { fail "missing prereq: $1"; exit 1; }
}
need curl
need tar
need shasum 2>/dev/null || need sha256sum
SHA_TOOL="$(command -v shasum >/dev/null 2>&1 && echo 'shasum -a 256' || echo 'sha256sum')"

# jq is preferred but optional — we have a sed fallback for the index.
HAVE_JQ=0
command -v jq >/dev/null 2>&1 && HAVE_JQ=1

# --- platform detection ---
detect_platform() {
  local os arch
  case "$(uname -s)" in
    Linux)  os="linux"  ;;
    Darwin) os="macos"  ;;
    *)      fail "unsupported OS: $(uname -s)"; exit 1 ;;
  esac
  case "$(uname -m)" in
    x86_64|amd64)   arch="x64"   ;;
    arm64|aarch64)  arch="arm64" ;;
    *)              fail "unsupported arch: $(uname -m)"; exit 1 ;;
  esac
  echo "${os}-${arch}"
}
PLATFORM="$(detect_platform)"
info "platform: $PLATFORM"

# --- network probe (AC5: graceful offline skip) ---
if ! curl -fsS --max-time 5 -o /dev/null https://api.github.com/zen 2>/dev/null; then
  skip "no network reachability to api.github.com — skipping (AC5)"
  echo "RESULT: skipped (offline)"
  exit 2
fi

# --- resolve tag ---
api() { curl -fsS --max-time 30 -H "Accept: application/vnd.github+json" "$@"; }

if [[ -z "$TAG" ]]; then
  info "resolving latest release tag from GitHub API"
  if [[ "$HAVE_JQ" -eq 1 ]]; then
    TAG="$(api "https://api.github.com/repos/${REPO}/releases/latest" | jq -r .tag_name)"
  else
    # Fallback: grep the tag_name line. Conservative — bare alphanumerics + dots/dashes.
    TAG="$(api "https://api.github.com/repos/${REPO}/releases/latest" \
      | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name":[[:space:]]*"([^"]+)".*/\1/')"
  fi
fi
if [[ -z "$TAG" || "$TAG" == "null" ]]; then
  skip "no release available on $REPO (empty tag_name) — skipping"
  echo "RESULT: skipped (no-release)"
  exit 2
fi
info "release tag: $TAG"

DL="https://github.com/${REPO}/releases/download/${TAG}"

# --- workspace ---
WORKDIR="$(mktemp -d -t cn-bootstrap-XXXXXX)"
trap 'rm -rf "$WORKDIR"' EXIT
info "workspace: $WORKDIR"

# --- download cn binary ---
BIN="$WORKDIR/cn"
info "downloading cn-$PLATFORM"
if ! curl -fsSL --max-time 120 -o "$BIN" "${DL}/cn-${PLATFORM}"; then
  fail "could not download cn-${PLATFORM} from ${DL}"
  exit 1
fi
chmod +x "$BIN"

# Sanity: the binary executes. `cn status` exits non-zero outside a
# hub, which is expected; the `|| true` swallows the rc and the line
# only fails if the binary cannot be invoked at all (missing loader,
# wrong arch, exec bit lost in transit).
"$BIN" status >/dev/null 2>&1 || true
pass "binary cn-$PLATFORM downloaded and executable"

# --- download index + tarballs into a dist/packages/ tree ---
# restore.FindIndexPath walks up from the hub looking for
# dist/packages/index.json. Lay the assets out so the hub at
# $WORKDIR/scratch/cn-scratch-hub/ resolves the index at
# $WORKDIR/dist/packages/index.json with relative tarball URLs
# adjacent.
PKG_DIR="$WORKDIR/dist/packages"
mkdir -p "$PKG_DIR"

info "downloading packages/index.json"
if ! curl -fsSL --max-time 60 -o "$PKG_DIR/index.json" "${DL}/index.json"; then
  fail "could not download index.json from ${DL}"
  exit 1
fi
# index.json is the integrity authority for tarball SHA-256 — it is what
# `cn deps restore` itself trusts at runtime. The release also publishes
# a checksums.txt covering the cn binary assets (release.yml `Generate
# binary checksums`); that file is not the tarball authority and is not
# fetched here, keeping the smoke aligned with one source of truth.

# Enumerate (name,version,filename) triples from the index. The release
# pipeline writes URLs as bare filenames (cn build's relative form), so
# we treat .url as the filename to fetch from the same release.
list_index_entries() {
  if [[ "$HAVE_JQ" -eq 1 ]]; then
    jq -r '
      .packages
      | to_entries[]
      | .key as $name
      | .value | to_entries[]
      | "\($name)\t\(.key)\t\(.value.url)\t\(.value.sha256)"
    ' "$PKG_DIR/index.json"
  else
    # jq-free fallback: parse the canonical structure with python3 if
    # available; otherwise warn and bail.
    if command -v python3 >/dev/null 2>&1; then
      python3 - "$PKG_DIR/index.json" <<'PY'
import json, sys
with open(sys.argv[1]) as f:
    idx = json.load(f)
for name, vers in (idx.get("packages") or {}).items():
    for v, e in vers.items():
        print(f"{name}\t{v}\t{e.get('url','')}\t{e.get('sha256','')}")
PY
    else
      fail "neither jq nor python3 available to parse index.json"
      exit 1
    fi
  fi
}

ENTRIES="$(list_index_entries)"
if [[ -z "$ENTRIES" ]]; then
  fail "index.json has no package entries — release is empty?"
  exit 1
fi

# Download each tarball and verify SHA-256 against the index entry.
N_PKGS=0
while IFS=$'\t' read -r name ver url sha; do
  [[ -z "$name" ]] && continue
  N_PKGS=$((N_PKGS + 1))
  fname="$url"
  # If the index ever moves to absolute URLs, take the basename.
  case "$fname" in
    http://*|https://*) fname="$(basename "$fname")" ;;
  esac
  dest="$PKG_DIR/$fname"
  if [[ ! -f "$dest" ]]; then
    info "downloading $fname"
    if ! curl -fsSL --max-time 180 -o "$dest" "${DL}/${fname}"; then
      fail "could not download $fname"
      exit 1
    fi
  fi
  actual="$($SHA_TOOL "$dest" | awk '{print $1}')"
  if [[ "$actual" != "$sha" ]]; then
    fail "sha256 mismatch for $name@$ver: index=$sha got=$actual"
    exit 1
  fi
done <<< "$ENTRIES"
pass "downloaded + verified $N_PKGS tarball(s)"

# --- bootstrap chain ---
SCRATCH="$WORKDIR/scratch"
mkdir -p "$SCRATCH"
cd "$SCRATCH"

info "cn init scratch-hub"
if ! "$BIN" init scratch-hub >/dev/null; then
  fail "cn init failed"; exit 1
fi
HUB="$SCRATCH/cn-scratch-hub"
[[ -d "$HUB" ]] || { fail "cn init did not create cn-scratch-hub"; exit 1; }
cd "$HUB"

info "cn setup"
if ! "$BIN" setup >/dev/null; then
  fail "cn setup failed"; exit 1
fi

# AC3 strict reading: the chain `init + setup + deps restore` must work
# against the production package source. `cn setup` writes a default
# deps.json pinning cnos.core/cnos.eng to the binary's compiled-in
# version. Verify each pinned (name, version) is present in the index
# we downloaded — if not, the released binary cannot self-bootstrap and
# the smoke must fail honestly.
DEPS_JSON="$HUB/.cn/deps.json"
[[ -f "$DEPS_JSON" ]] || { fail ".cn/deps.json not written by cn setup"; exit 1; }

index_has() {
  # index_has <name> <version> → exit 0 iff index.json packages[name][version] exists.
  local n="$1" v="$2"
  if [[ "$HAVE_JQ" -eq 1 ]]; then
    jq -e --arg n "$n" --arg v "$v" \
      '(.packages // {}) | (.[$n] // {}) | has($v)' \
      "$PKG_DIR/index.json" >/dev/null 2>&1
  else
    python3 - "$PKG_DIR/index.json" "$n" "$v" <<'PY'
import json, sys
with open(sys.argv[1]) as f:
    d = json.load(f)
pkgs = d.get("packages") or {}
sys.exit(0 if sys.argv[2] in pkgs and sys.argv[3] in pkgs[sys.argv[2]] else 1)
PY
  fi
}

verify_default_deps_resolvable() {
  local missing=""
  while IFS=$'\t' read -r name ver; do
    [[ -z "$name" ]] && continue
    if ! index_has "$name" "$ver"; then
      missing="${missing}${missing:+, }${name}@${ver}"
    fi
  done < <(default_deps_pins)
  if [[ -n "$missing" ]]; then
    fail "cn setup pinned packages not in released index: $missing"
    fail "  → released binary's compiled version disagrees with released package set"
    exit 1
  fi
}

default_deps_pins() {
  if [[ "$HAVE_JQ" -eq 1 ]]; then
    jq -r '.packages[] | "\(.name)\t\(.version)"' "$DEPS_JSON"
  elif command -v python3 >/dev/null 2>&1; then
    python3 - "$DEPS_JSON" <<'PY'
import json, sys
with open(sys.argv[1]) as f:
    d = json.load(f)
for p in d.get("packages", []):
    print(f"{p.get('name','')}\t{p.get('version','')}")
PY
  else
    fail "neither jq nor python3 available to parse deps.json"; exit 1
  fi
}

verify_default_deps_resolvable
pass "cn setup default deps resolvable in released index"

info "cn deps lock"
if ! "$BIN" deps lock >/dev/null; then
  fail "cn deps lock failed"; exit 1
fi

info "cn deps restore"
if ! "$BIN" deps restore >/dev/null; then
  fail "cn deps restore failed"; exit 1
fi

# --- verify installation ---
ok=0; bad=""
while IFS=$'\t' read -r name ver; do
  [[ -z "$name" ]] && continue
  pkg_dir="$HUB/.cn/vendor/packages/$name"
  m="$pkg_dir/cn.package.json"
  if [[ ! -f "$m" ]]; then
    bad="${bad}${bad:+, }${name}(no manifest)"
    continue
  fi
  installed_ver=""
  if [[ "$HAVE_JQ" -eq 1 ]]; then
    installed_ver="$(jq -r .version "$m")"
  else
    installed_ver="$(grep -m1 '"version"' "$m" | sed -E 's/.*"version":[[:space:]]*"([^"]+)".*/\1/')"
  fi
  if [[ "$installed_ver" != "$ver" ]]; then
    bad="${bad}${bad:+, }${name}(installed=${installed_ver}, want=${ver})"
    continue
  fi
  ok=$((ok + 1))
done < <(default_deps_pins)

if [[ -n "$bad" ]]; then
  fail "vendor verification failed: $bad"
  echo "RESULT: failed"
  exit 1
fi

pass "vendor verified: $ok package(s) match deps.json pins"
echo "RESULT: ok ($TAG / $PLATFORM, $ok package(s) verified)"
