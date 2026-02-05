# tool-writing

Standards for writing tools in `tools/` — mechanical scripts that run without AI.

---

## TERMS

1. The task is mechanical (no judgment required).
2. AI would burn tokens on clockwork work.
3. A script can do it reliably and cheaply.

---

## Principles

### Zero runtime dependencies

```bash
#!/bin/bash
# Only use: bash, git, coreutils, jq (if JSON needed)
# No: npm, pip, ruby, etc.
```

Why: Tools must work on any Unix system without setup.

### NO_COLOR support

```bash
# At top of script
if [[ -n "${NO_COLOR:-}" ]]; then
  RED="" GREEN="" YELLOW="" CYAN="" RESET=""
else
  RED='\033[0;31m' GREEN='\033[0;32m' 
  YELLOW='\033[0;33m' CYAN='\033[0;36m' RESET='\033[0m'
fi
```

Why: Accessibility. Some users pipe output or use screen readers.

### Autoconf-style prereq checks

```bash
check_prereqs() {
  local missing=()
  command -v git &>/dev/null || missing+=(git)
  command -v jq &>/dev/null || missing+=(jq)
  
  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "Missing: ${missing[*]}" >&2
    exit 1
  fi
}
```

Why: Fail fast with clear message, not cryptic errors mid-run.

### Idempotent

```bash
# Safe to run multiple times
# Don't append if already exists
grep -q "pattern" file || echo "pattern" >> file
```

Why: Users will run tools repeatedly. Don't corrupt state.

### Machine-readable output

```bash
# Exit codes
# 0 = success
# 1 = error
# 2 = nothing to do (already up-to-date)

# Parseable stdout
echo "FETCHED:cn-pi:3"  # not "Fetched 3 new items from cn-pi!"
```

Why: Other tools/scripts may consume output.

### Semantic logging

```bash
log_ok()   { echo -e "${GREEN}✓${RESET} $1"; }
log_err()  { echo -e "${RED}✗${RESET} $1" >&2; }
log_warn() { echo -e "${YELLOW}⏸${RESET} $1"; }
log_info() { echo -e "${CYAN}→${RESET} $1"; }
```

Why: Visual consistency. Users learn the symbols.

### Fail safely

```bash
set -euo pipefail  # Exit on error, undefined var, pipe fail

# Trap for cleanup
trap 'echo "Error on line $LINENO" >&2' ERR
```

Why: No silent failures. Know exactly where it broke.

### Header comment

```bash
#!/bin/bash
# tool-name.sh — one-line description
# Usage: ./tools/tool-name.sh [args]
# Output: what it produces
# Exit: 0=success, 1=error, 2=no-op
```

Why: Self-documenting. `head -5 tool.sh` explains it.

---

## Template

```bash
#!/bin/bash
# tool-name.sh — description
# Usage: ./tools/tool-name.sh [hub-dir]
# Output: description of output
# Exit: 0=success, 1=error, 2=no-op

set -euo pipefail

# Colors (NO_COLOR support)
if [[ -n "${NO_COLOR:-}" ]]; then
  RED="" GREEN="" YELLOW="" CYAN="" RESET=""
else
  RED='\033[0;31m' GREEN='\033[0;32m'
  YELLOW='\033[0;33m' CYAN='\033[0;36m' RESET='\033[0m'
fi

log_ok()   { echo -e "${GREEN}✓${RESET} $1"; }
log_err()  { echo -e "${RED}✗${RESET} $1" >&2; }
log_warn() { echo -e "${YELLOW}⏸${RESET} $1"; }
log_info() { echo -e "${CYAN}→${RESET} $1"; }

# Prereq check
check_prereqs() {
  local missing=()
  command -v git &>/dev/null || missing+=(git)
  if [[ ${#missing[@]} -gt 0 ]]; then
    log_err "Missing: ${missing[*]}"
    exit 1
  fi
}

main() {
  check_prereqs
  
  # Your logic here
  
  log_ok "Done"
}

main "$@"
```

---

## Anti-patterns

| Don't | Do |
|-------|-----|
| `curl \| bash` | Download, inspect, then run |
| Hardcoded paths | Use `$1` or `${HUB_DIR:-$PWD}` |
| Silent failures | `set -euo pipefail` |
| Fancy deps (python, node) | Pure bash + coreutils |
| ANSI without NO_COLOR check | Always check `$NO_COLOR` |

---

## NOTES

- Tools go in `tools/`, not `skills/`
- Skills = AI judgment required
- Tools = mechanical, no AI needed
- When in doubt: if it's just git/file ops, it's a tool
