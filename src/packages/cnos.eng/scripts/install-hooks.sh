#!/bin/bash
# install-hooks.sh — Install cnos engineering git hooks
# Usage: ./src/packages/cnos.eng/scripts/install-hooks.sh
# Exit: 0=success, 1=error, 2=no-op (already installed)

set -euo pipefail

if [[ -n "${NO_COLOR:-}" ]]; then
  RED="" GREEN="" YELLOW="" RESET=""
else
  RED='\033[0;31m' GREEN='\033[0;32m' YELLOW='\033[1;33m' RESET='\033[0m'
fi

check_prereqs() {
  command -v git &>/dev/null || { echo "Missing: git" >&2; exit 1; }

  # Check if we're in a git repository
  if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo -e "${RED}✗${RESET} Not in a git repository" >&2
    exit 1
  fi

  # Check if hooks directory exists
  if [[ ! -d "src/packages/cnos.eng/hooks" ]]; then
    echo -e "${RED}✗${RESET} Hook source directory not found: src/packages/cnos.eng/hooks" >&2
    exit 1
  fi
}

main() {
  check_prereqs

  local hooks_path="src/packages/cnos.eng/hooks"
  local current_hooks_path
  current_hooks_path=$(git config --get core.hooksPath || echo "")

  # Check if hooks are already configured
  if [[ "$current_hooks_path" == "$hooks_path" ]]; then
    echo -e "${GREEN}✓${RESET} Hooks already installed at $hooks_path"
    exit 2
  fi

  # Configure git to use our hooks directory
  git config core.hooksPath "$hooks_path"
  echo -e "${GREEN}✓${RESET} Installed git hooks from $hooks_path"

  # List available hooks
  echo "Available hooks:"
  find "$hooks_path" -type f -executable -printf "  %f\n" | sort

  # Verify hook is executable
  local pre_push_hook="$hooks_path/pre-push"
  if [[ -f "$pre_push_hook" && -x "$pre_push_hook" ]]; then
    echo -e "${GREEN}✓${RESET} pre-push hook is executable"
  else
    echo -e "${YELLOW}⚠${RESET} pre-push hook may not be executable"
  fi

  echo ""
  echo "Hooks installed successfully."
  echo "To bypass rebase integrity check: ALLOW_CONTENT_LOSS=1 git push"
}

main "$@"