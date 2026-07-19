#!/usr/bin/env bash
# Enforce the v0.2 authority path-base and activation contract.

set -euo pipefail

usage() {
  echo "usage: $0 --source-root REPOSITORY | --installed-root PACKAGE" >&2
  exit 2
}

mode=""
root=""
case "${1:-}" in
  --source-root) mode="source"; root="${2:-}" ;;
  --installed-root) mode="installed"; root="${2:-}" ;;
  *) usage ;;
esac
[[ -n "$root" ]] || usage

relative="src/packages/cnos.cdd/skills/cdd/measure/recursive-cell"
if [[ "$mode" == "installed" ]]; then
  [[ -f "$root/skills/cdd/measure/recursive-cell/SKILL.md" ]] || {
    echo "FAIL installed package layout lacks recursive-cell/SKILL.md" >&2
    exit 1
  }
  echo "REFUSE installed activation: methodology.activation=source-checkout-only; TSC has no separate authority-root/target-root contract" >&2
  exit 3
fi

[[ -d "$root/.git" ]] || { echo "FAIL source root is not a repository checkout: $root" >&2; exit 1; }
cm="$root/$relative"
[[ -f "$cm/SKILL.md" ]] || { echo "FAIL missing source authority: $cm/SKILL.md" >&2; exit 1; }
[[ -f "$cm/calibration/662/registry.tsc" ]] || { echo "FAIL missing registry" >&2; exit 1; }
[[ -f "$cm/INSTRUCTION.md" ]] || { echo "FAIL missing instruction" >&2; exit 1; }
[[ -x "$cm/calibration/662/verify-target.sh" ]] || { echo "FAIL preflight is not executable" >&2; exit 1; }
"$cm/instruction/assemble-instruction.sh" --check
printf 'PASS source authority base=repository-root cm=%s\n' "$relative"
