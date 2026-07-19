#!/usr/bin/env bash
# Deterministically assemble the single instruction accepted by TSC's CLI.

set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
core="$here/TSC-SELF-MEASURE-v3.2.4.md"
supplement="$here/CNOS-SUPPLEMENT-v0.3.0.md"
canonical="$here/../INSTRUCTION.md"
expected_core="ca26d7a1a4dc6bd73e0afff558ed0342d42daeee193d4169de2c51b762759391"
expected_supplement="3ad7c80ad318b2f736738a444da189b72712dc61914a52f183c7c7c584e377ee"
expected_composite="de0414bac1a086d6ed3ac5c7d84cd6293517eb83c920f6c4a0723f5dfac06a3a"

mode="emit"
output=""
while (($#)); do
  case "$1" in
    --check) mode="check"; shift ;;
    --output) output="$2"; shift 2 ;;
    *) echo "usage: $0 [--check] [--output PATH]" >&2; exit 2 ;;
  esac
done

check_digest() {
  local expected="$1" path="$2" actual
  actual="$(sha256sum "$path" | awk '{print $1}')"
  [[ "$actual" == "$expected" ]] || {
    printf 'FAIL digest %s: expected %s, got %s\n' "$path" "$expected" "$actual" >&2
    exit 1
  }
}

check_digest "$expected_core" "$core"
check_digest "$expected_supplement" "$supplement"

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
{
  command cat "$core"
  printf '\n---\n\n'
  command cat "$supplement"
} >"$tmp"
check_digest "$expected_composite" "$tmp"

if [[ "$mode" == "check" ]]; then
  cmp -s "$tmp" "$canonical" || {
    echo "FAIL canonical composite differs from deterministic assembly: $canonical" >&2
    exit 1
  }
  check_digest "$expected_composite" "$canonical"
  printf 'PASS composite instruction sha256=%s\n' "$expected_composite"
elif [[ -n "$output" ]]; then
  command cp "$tmp" "$output"
else
  command cat "$tmp"
fi
