#!/usr/bin/env bash
# Fail-closed preflight for the historical #662 R7 calibration corpus.

set -euo pipefail

target_root="${1:-.}"
cm_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
expected_receipt_head="a0d39293a27cfe57b49dacff696345b1ee2cdb40"
matter_revision="2d6b93cc4e69e5b413a80bd8e352cb0a004da460"
matter_path="docs/architecture/CELL-RUNTIME-CLASSES.md"
expected_matter_sha256="80e0d8c68a3d8affabdd4bd14848cbb3f0bee27b078e435f61433e42a0ee89e0"

"$cm_root/instruction/assemble-instruction.sh" --check

actual_head="$(git -C "$target_root" rev-parse HEAD)"
if [[ "$actual_head" != "$expected_receipt_head" ]]; then
  printf 'FAIL receipt head: expected %s, got %s\n' "$expected_receipt_head" "$actual_head" >&2
  exit 1
fi

matter_sha256="$(git -C "$target_root" show "$matter_revision:$matter_path" | sha256sum | awk '{print $1}')"
if [[ "$matter_sha256" != "$expected_matter_sha256" ]]; then
  printf 'FAIL matter bytes: expected %s, got %s\n' "$expected_matter_sha256" "$matter_sha256" >&2
  exit 1
fi

if ! git -C "$target_root" diff --quiet "$matter_revision..$expected_receipt_head" -- "$matter_path"; then
  printf 'FAIL matter moved between %s and %s\n' "$matter_revision" "$expected_receipt_head" >&2
  exit 1
fi

printf 'PASS cnos#662 calibration target: receipt=%s matter=%s sha256=%s\n' \
  "$expected_receipt_head" "$matter_revision" "$matter_sha256"
