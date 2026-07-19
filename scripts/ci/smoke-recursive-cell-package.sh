#!/usr/bin/env bash
# Build cnos.cdd, inspect its public recursive-cell layout, and prove that the
# source-checkout-only activation contract fails closed after installation.

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
version="$(jq -r '.version' "$repo_root/src/packages/cnos.cdd/cn.package.json")"
tarball="$repo_root/dist/packages/cnos.cdd-${version}.tar.gz"

(cd "$repo_root" && cn build >/dev/null)
[[ -f "$tarball" ]] || { echo "FAIL package not built: $tarball" >&2; exit 1; }

required=(
  "skills/cdd/measure/recursive-cell/SKILL.md"
  "skills/cdd/measure/recursive-cell/INSTRUCTION.md"
  "skills/cdd/measure/recursive-cell/resolve-authority.sh"
  "skills/cdd/measure/recursive-cell/instruction/assemble-instruction.sh"
  "skills/cdd/measure/recursive-cell/calibration/662/registry.tsc"
  "skills/cdd/measure/recursive-cell/runner/recursive-cell-runner.py"
  "skills/cdd/measure/recursive-cell/runner/recursive-cell-run.schema.cue"
  "skills/cdd/measure/recursive-cell/runner/invariant-assessment-template.md"
)
listing="$(tar -tzf "$tarball")"
for path in "${required[@]}"; do
  grep -Fxq "$path" <<<"$listing" || { echo "FAIL package layout missing: $path" >&2; exit 1; }
done

install_root="$(mktemp -d)"
trap 'rm -rf "$install_root"' EXIT
tar -xzf "$tarball" -C "$install_root"
resolver="$install_root/skills/cdd/measure/recursive-cell/resolve-authority.sh"
set +e
refusal="$($resolver --installed-root "$install_root" 2>&1)"
rc=$?
set -e
[[ $rc -eq 3 ]] || { echo "FAIL installed activation returned $rc, expected 3" >&2; exit 1; }
grep -Fq "REFUSE installed activation" <<<"$refusal" || {
  echo "FAIL installed activation did not state the source-only refusal" >&2
  exit 1
}

printf 'PASS cnos.cdd package layout and installed activation refusal: %s\n' "$tarball"
