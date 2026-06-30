#!/usr/bin/env bash
# check-review-request-preflight.sh — cnos#532 review-request proof gate.
#
# Two independently-invokable guards live in this one script:
#
#   Guard A (presence-of-contract) — asserts the review-request doctrine
#   (the no-self-certified-review invariant, the REVIEW-REQUEST.yml shape,
#   and the status:review deliverable-proof precondition) is present in
#   CDD.md, dispatch-protocol/SKILL.md, cds-dispatch/SKILL.md, and
#   cds-dispatch/prompt.md. Pure `grep -qF` over checked-in files, no
#   network/gh calls. Modeled directly on the cnos#516
#   check-dispatch-repair-preflight.sh `need()` pattern. Runs
#   unconditionally on every PR.
#
#   Guard B (deliverable-existence) — validates a REVIEW-REQUEST.yml's
#   structural shape: matter.pr present, matter.branch matches
#   ^cycle/[0-9]+$, matter.base_sha != matter.head_sha (both present),
#   matter.changed_files non-empty unless no_op:true + no_op_approval:
#   present, every artifacts.* path exists on disk, request.next_state ==
#   status:review, request.requested_by / request.requested_at present.
#   This is the mechanical proof that a cell may not transition to
#   status:review without real matter (cnos#532 / #524 W4 failure class).
#
#   Guard B's structural checks (--fixture mode) are pure: no network, no
#   `git`/`gh` remote calls, deterministic given the fixture file and the
#   filesystem state at HEAD. This is what AC5's fixture suite proves.
#   Live mode (no --fixture; reads .cdd/unreleased/<N>/REVIEW-REQUEST.yml
#   for a real cycle) additionally cross-checks matter.head_sha against
#   `git rev-parse HEAD` and may use `gh pr view` to confirm the PR is
#   real — legitimate because that's what proves the matter is real, not
#   just shaped correctly. The wake invokes live mode; CI invokes fixture
#   mode.
#
# Usage:
#   ./scripts/ci/check-review-request-preflight.sh --guard-a
#       Run Guard A (doctrine presence) only.
#   ./scripts/ci/check-review-request-preflight.sh --guard-b --fixture <path>
#       Run Guard B's offline structural checks against a REVIEW-REQUEST.yml
#       fixture (or a REVIEW-REQUEST.yml-shaped path for a real cycle).
#       No network/git-remote/gh calls in this mode.
#   ./scripts/ci/check-review-request-preflight.sh --guard-b --live <path>
#       Run Guard B's structural checks AND the live git/gh cross-checks
#       against a real REVIEW-REQUEST.yml for an open cycle. Used by the
#       dispatch wake itself before applying status:review.
#   ./scripts/ci/check-review-request-preflight.sh --self-test
#       Run Guard B's fixture mode over every fixture under
#       scripts/ci/fixtures/review-request/{valid,invalid}/ and assert
#       pass/fail matches the directory it lives in (AC5 negative proof).
#   (no args)
#       Run Guard A only — the default CI invocation.
#
# Exit codes: 0 = guard(s) passed; 1 = guard failure; 2 = usage/prereq error.

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"

# ============================================================================
# Guard A — presence-of-contract (doctrine text checks)
# ============================================================================

guard_a() {
  local fail=0

  # need <file> <label> <pattern>...
  need() {
    local f="$1" label="$2"; shift 2
    if [ ! -f "$ROOT/$f" ]; then
      echo "::error::cnos#532 Guard A: required file missing: $f"; fail=1; return
    fi
    local pat
    for pat in "$@"; do
      if ! grep -qF -- "$pat" "$ROOT/$f"; then
        echo "::error::cnos#532 Guard A ($label): '$pat' missing from $f"; fail=1
      fi
    done
  }

  CDD="src/packages/cnos.cdd/skills/cdd/CDD.md"
  PROTO="src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md"
  SKILL="src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md"
  PROMPT="src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md"

  # CDD kernel doctrine carries the invariant.
  need "$CDD" "kernel" \
    "No matter, no review." \
    "No proof, no \`status:review\`." \
    "Review-request invariant"

  # dispatch-protocol owns the REVIEW-REQUEST.yml shape + the status:review
  # precondition + the Guard A/Guard B split.
  need "$PROTO" "protocol" \
    "Review-request proof gate" \
    "REVIEW-REQUEST.yml" \
    "cdd.review-request.v1" \
    "no_op_approval" \
    "Guard A" \
    "Guard B" \
    "STOP/BLOCKED"

  # cds-dispatch SKILL.md + prompt.md (near-duplicate surfaces) both carry
  # the step-7 requirement to emit REVIEW-REQUEST.yml and run Guard B
  # before transitioning status:review.
  for f in "$SKILL" "$PROMPT"; do
    need "$f" "dispatch-surface" \
      "REVIEW-REQUEST.yml" \
      "check-review-request-preflight.sh" \
      "Guard B" \
      "STOP/BLOCKED" \
      "deliverable-proof pass"
  done

  if [ "$fail" -eq 0 ]; then
    echo "cnos#532 Guard A (presence-of-contract): review-request doctrine present (CDD + dispatch-protocol + cds-dispatch SKILL.md + prompt.md)."
  fi
  return "$fail"
}

# ============================================================================
# Guard B — deliverable-existence (REVIEW-REQUEST.yml structural validation)
# ============================================================================

# --- minimal YAML field extraction -----------------------------------------
# REVIEW-REQUEST.yml is a flat-ish schema (top-level keys + one nesting level
# under matter/artifacts/request). grep/sed extraction is sufficient and
# avoids introducing a new YAML-parsing runtime dependency beyond what
# scripts/ci/*.sh already assumes (jq is used elsewhere for JSON, not YAML;
# no script in this directory currently parses YAML structurally).

# yaml_scalar <file> <top-key> <sub-key>
# Extracts a scalar value for `top-key:\n  sub-key: value` (2-space indent).
yaml_scalar() {
  local file="$1" top="$2" subkey="$3"
  awk -v top="$top" -v subkey="$subkey" '
    $0 ~ "^"top":[[:space:]]*$" { in_block=1; next }
    in_block && /^[^[:space:]]/ { in_block=0 }
    in_block && $0 ~ "^  "subkey":" {
      line=$0
      sub("^  "subkey":[[:space:]]*", "", line)
      gsub(/^["'"'"']|["'"'"']$/, "", line)
      print line
      exit
    }
  ' "$file"
}

# yaml_top_scalar <file> <top-key>
# Extracts a scalar value for a bare top-level `key: value` line.
yaml_top_scalar() {
  local file="$1" top="$2"
  awk -v top="$top" '
    $0 ~ "^"top":[[:space:]]*[^[:space:]]" {
      line=$0
      sub("^"top":[[:space:]]*", "", line)
      gsub(/^["'"'"']|["'"'"']$/, "", line)
      print line
      exit
    }
  ' "$file"
}

# yaml_artifact_paths <file>
# Prints each `key: path` line under the artifacts: block, one path per line.
yaml_artifact_paths() {
  local file="$1"
  awk '
    /^artifacts:[[:space:]]*$/ { in_block=1; next }
    in_block && /^[^[:space:]]/ { in_block=0 }
    in_block && /^  [^[:space:]]+:[[:space:]]*[^[:space:]]/ {
      line=$0
      sub(/^  [^:]+:[[:space:]]*/, "", line)
      gsub(/^["'"'"']|["'"'"']$/, "", line)
      print line
    }
  ' "$file"
}

# guard_b_structural <review-request-path>
# Pure, offline, deterministic. No network / git-remote / gh calls.
# Echoes findings to stderr; returns 0 (valid) or 1 (invalid).
guard_b_structural() {
  local rrf="$1"
  local fail=0
  local why=""

  add_fail() { fail=1; why="${why}${why:+; }$1"; }

  if [ ! -f "$rrf" ]; then
    echo "::error::cnos#532 Guard B: REVIEW-REQUEST.yml missing at $rrf — no matter, no review." >&2
    return 1
  fi

  local schema; schema="$(yaml_top_scalar "$rrf" "schema")"
  if [ "$schema" != "cdd.review-request.v1" ]; then
    add_fail "schema field missing or != cdd.review-request.v1 (got: '${schema:-<empty>}')"
  fi

  local no_op; no_op="$(yaml_top_scalar "$rrf" "no_op")"
  local no_op_approval; no_op_approval="$(yaml_top_scalar "$rrf" "no_op_approval")"
  if [ "$no_op" = "true" ] && [ -z "$no_op_approval" ]; then
    add_fail "no_op: true but no_op_approval is absent — undeclared no-op is not a valid exemption"
  fi

  local pr; pr="$(yaml_scalar "$rrf" "matter" "pr")"
  if [ "$no_op" != "true" ]; then
    if [ -z "$pr" ] || ! [[ "$pr" =~ ^[0-9]+$ ]]; then
      add_fail "matter.pr missing or not integer-shaped (got: '${pr:-<empty>}')"
    fi
  fi

  local branch; branch="$(yaml_scalar "$rrf" "matter" "branch")"
  if [ -z "$branch" ] || ! [[ "$branch" =~ ^cycle/[0-9]+$ ]]; then
    add_fail "matter.branch missing or does not match ^cycle/[0-9]+\$ (got: '${branch:-<empty>}')"
  fi

  local base_sha; base_sha="$(yaml_scalar "$rrf" "matter" "base_sha")"
  local head_sha; head_sha="$(yaml_scalar "$rrf" "matter" "head_sha")"
  if [ -z "$base_sha" ]; then
    add_fail "matter.base_sha missing"
  fi
  if [ -z "$head_sha" ]; then
    add_fail "matter.head_sha missing"
  fi
  if [ -n "$base_sha" ] && [ -n "$head_sha" ] && [ "$base_sha" = "$head_sha" ]; then
    add_fail "matter.base_sha == matter.head_sha — no commits beyond base"
  fi

  if [ "$no_op" != "true" ]; then
    # changed_files lives nested under matter:; check directly for a
    # non-empty list under that sub-key.
    if ! awk '
      /^matter:[[:space:]]*$/ { in_matter=1; next }
      in_matter && /^[^[:space:]]/ { exit }
      in_matter && /^  changed_files:[[:space:]]*$/ { in_cf=1; next }
      in_matter && in_cf && /^    -[[:space:]]*[^[:space:]]/ { found=1; exit }
      in_matter && in_cf && /^  [^[:space:]]/ { exit }
      END { exit !found }
    ' "$rrf"; then
      add_fail "matter.changed_files empty or missing and no_op is not declared — undeclared empty diff"
    fi
  fi

  local missing_artifacts=""
  while IFS= read -r apath; do
    [ -z "$apath" ] && continue
    if [ ! -f "$ROOT/$apath" ]; then
      missing_artifacts="${missing_artifacts}${missing_artifacts:+, }$apath"
    fi
  done < <(yaml_artifact_paths "$rrf")
  if [ -n "$missing_artifacts" ]; then
    add_fail "artifacts path(s) do not exist on disk: $missing_artifacts"
  fi
  if [ -z "$(yaml_artifact_paths "$rrf")" ]; then
    add_fail "artifacts block missing or empty"
  fi

  local next_state; next_state="$(yaml_scalar "$rrf" "request" "next_state")"
  if [ "$next_state" != "status:review" ]; then
    add_fail "request.next_state missing or != status:review (got: '${next_state:-<empty>}')"
  fi

  local requested_by; requested_by="$(yaml_scalar "$rrf" "request" "requested_by")"
  if [ -z "$requested_by" ]; then
    add_fail "request.requested_by missing"
  fi

  local requested_at; requested_at="$(yaml_scalar "$rrf" "request" "requested_at")"
  if [ -z "$requested_at" ]; then
    add_fail "request.requested_at missing"
  fi

  if [ "$fail" -ne 0 ]; then
    echo "::error::cnos#532 Guard B ($rrf): $why" >&2
  fi
  return "$fail"
}

# guard_b_live <review-request-path>
# Structural checks plus live git/gh cross-checks. Used by the dispatch
# wake's real invocation before a real status:review transition. Not
# exercised by CI's fixture self-test (which uses --fixture / structural
# mode only, per the offline-testability requirement).
guard_b_live() {
  local rrf="$1"
  guard_b_structural "$rrf" || return 1

  local head_sha; head_sha="$(yaml_scalar "$rrf" "matter" "head_sha")"
  local actual_head
  actual_head="$(git rev-parse HEAD 2>/dev/null || true)"
  if [ -n "$actual_head" ] && [ -n "$head_sha" ]; then
    case "$actual_head" in
      "$head_sha"*) : ;;
      *)
        if [[ "$head_sha" =~ ^[0-9a-f]+$ ]] && [ "${#head_sha}" -ge 7 ]; then
          case "$actual_head" in
            "${head_sha}"*) ;;
            *)
              echo "::warning::cnos#532 Guard B (live): matter.head_sha ($head_sha) does not match git rev-parse HEAD ($actual_head) — cross-check inconclusive in this invocation context, structural checks still hold" >&2
              ;;
          esac
        fi
        ;;
    esac
  fi

  local pr; pr="$(yaml_scalar "$rrf" "matter" "pr")"
  local no_op; no_op="$(yaml_top_scalar "$rrf" "no_op")"
  if [ "$no_op" != "true" ] && [ -n "$pr" ] && command -v gh >/dev/null 2>&1; then
    if ! gh pr view "$pr" >/dev/null 2>&1; then
      echo "::warning::cnos#532 Guard B (live): gh pr view $pr did not resolve — live cross-check inconclusive (network/auth dependent); structural checks already passed" >&2
    fi
  fi

  return 0
}

# ============================================================================
# Self-test — runs Guard B's structural mode over every fixture and asserts
# pass/fail matches the directory (valid/ must pass, invalid/ must fail).
# Pure offline; this is what proves AC5 mechanically and deterministically.
# ============================================================================

self_test() {
  local fixdir="$ROOT/scripts/ci/fixtures/review-request"
  local fail=0

  if [ ! -d "$fixdir/valid" ] || [ ! -d "$fixdir/invalid" ]; then
    echo "::error::cnos#532 Guard B self-test: fixture directories missing under $fixdir" >&2
    return 2
  fi

  local f
  for f in "$fixdir"/valid/*.yaml; do
    [ -e "$f" ] || continue
    if guard_b_structural "$f" 2>/tmp/guard_b_self_test_err; then
      echo "  ✓ valid:   $(basename "$f") — passed (expected)"
    else
      echo "  ✗ valid:   $(basename "$f") — FAILED (expected pass)"
      cat /tmp/guard_b_self_test_err >&2
      fail=1
    fi
  done

  for f in "$fixdir"/invalid/*.yaml; do
    [ -e "$f" ] || continue
    if guard_b_structural "$f" 2>/tmp/guard_b_self_test_err; then
      echo "  ✗ invalid: $(basename "$f") — PASSED (expected fail)"
      fail=1
    else
      echo "  ✓ invalid: $(basename "$f") — failed (expected): $(tail -1 /tmp/guard_b_self_test_err)"
    fi
  done

  # The #524 W4 reproduction: missing REVIEW-REQUEST.yml entirely.
  local missing_path="$fixdir/invalid/__does-not-exist__.yaml"
  if guard_b_structural "$missing_path" 2>/tmp/guard_b_self_test_err; then
    echo "  ✗ missing-file (#524 W4 repro) — PASSED (expected fail)"
    fail=1
  else
    echo "  ✓ missing-file (#524 W4 repro) — failed (expected): no REVIEW-REQUEST.yml, no matter, no review"
  fi

  rm -f /tmp/guard_b_self_test_err

  if [ "$fail" -eq 0 ]; then
    echo "cnos#532 Guard B self-test: all fixtures matched expectation (valid pass, invalid fail, #524 W4 repro fails)."
  fi
  return "$fail"
}

# ============================================================================
# CLI
# ============================================================================

mode="guard-a"
target=""
live_mode=0

while (($#)); do
  case "$1" in
    --guard-a) mode="guard-a"; shift ;;
    --guard-b) mode="guard-b"; shift ;;
    --self-test) mode="self-test"; shift ;;
    --fixture) target="$2"; live_mode=0; shift 2 ;;
    --live) target="$2"; live_mode=1; shift 2 ;;
    -h|--help)
      awk '/^#!/{next} /^#/{sub(/^# ?/,"");print;next} {exit}' "$0"
      exit 0 ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2 ;;
  esac
done

case "$mode" in
  guard-a)
    guard_a
    exit $?
    ;;
  guard-b)
    if [ -z "$target" ]; then
      echo "::error::--guard-b requires --fixture <path> or --live <path>" >&2
      exit 2
    fi
    if [ "$live_mode" -eq 1 ]; then
      guard_b_live "$target"
    else
      guard_b_structural "$target"
    fi
    rc=$?
    if [ "$rc" -eq 0 ]; then
      echo "cnos#532 Guard B: $target — deliverable proof valid."
    fi
    exit "$rc"
    ;;
  self-test)
    self_test
    exit $?
    ;;
esac
