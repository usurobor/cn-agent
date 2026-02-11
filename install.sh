#!/bin/sh
# cnos installer - native OCaml binary
# Usage: curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh

set -e

REPO="https://github.com/usurobor/cnos"
INSTALL_DIR="/usr/local/lib/cnos"
BIN_DIR="/usr/local/bin"

echo "Installing cnos (native binary)..."

# Clone or pull
if [ -d "$INSTALL_DIR" ]; then
  echo "Updating existing installation..."
  cd "$INSTALL_DIR"
  git pull --ff-only
else
  echo "Cloning cnos..."
  git clone --depth 1 "$REPO.git" "$INSTALL_DIR"
  cd "$INSTALL_DIR"
fi

# Build native binary
if ! command -v opam >/dev/null 2>&1; then
  echo "Error: opam required. Install OCaml toolchain first:"
  echo "  https://ocaml.org/docs/installing-ocaml"
  exit 1
fi

echo "Building native binary..."
eval $(opam env)
dune build tools/src/cn/cn.exe

# Install binary
cp _build/default/tools/src/cn/cn.exe "$BIN_DIR/cn"
chmod +x "$BIN_DIR/cn"

# Install cn-cron
cp bin/cn-cron "$BIN_DIR/cn-cron"
chmod +x "$BIN_DIR/cn-cron"

echo ""
echo "✓ cnos installed successfully"
echo ""
cn --version
