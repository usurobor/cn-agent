#!/usr/bin/env bash
# tests/cdd/test_cn_cdd_validate_receipt.sh
#
# AC1–AC8 oracle harness for #389 (Phase 3 — V implementation).
# Run from the repo root: `bash tests/cdd/test_cn_cdd_validate_receipt.sh`.
#
# Exit 0 if all AC oracles pass; exit 1 on first failure with a diagnostic.
#
# This harness is the mechanical β-review surface for cycle/389: each AC's
# named oracle from the issue body is replayed; the harness reports PASS/FAIL
# per AC.  The harness must remain runnable from repo root after the cycle
# merges so future regressions are caught.

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$REPO_ROOT"

V="$REPO_ROOT/src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify"

PASS_CT=0
FAIL_CT=0
declare -a FAILURES

red()    { printf "\033[31m%s\033[0m" "$*"; }
green()  { printf "\033[32m%s\033[0m" "$*"; }
yellow() { printf "\033[33m%s\033[0m" "$*"; }

ac_pass() {
  PASS_CT=$((PASS_CT + 1))
  printf "  %s %s\n" "$(green PASS)" "$1"
}
ac_fail() {
  FAIL_CT=$((FAIL_CT + 1))
  FAILURES+=("$1: $2")
  printf "  %s %s — %s\n" "$(red FAIL)" "$1" "$2"
}

# Helper: run V and capture stdout, stderr, exit.
run_v() {
  local out err rc tmpdir
  tmpdir="$(mktemp -d)"
  "$V" "$@" > "$tmpdir/out" 2> "$tmpdir/err"
  rc=$?
  echo "RC=$rc"
  echo "OUT<<<"
  cat "$tmpdir/out"
  echo ">>>OUT"
  echo "ERR<<<"
  cat "$tmpdir/err"
  echo ">>>ERR"
  rm -rf "$tmpdir"
}

assert_exit() {
  local label="$1" want="$2" got="$3"
  if [[ "$got" -eq "$want" ]]; then ac_pass "$label"
  else ac_fail "$label" "expected exit $want, got $got"
  fi
}
assert_contains() {
  local label="$1" needle="$2" haystack="$3"
  if echo "$haystack" | grep -qF -- "$needle"; then ac_pass "$label"
  else ac_fail "$label" "expected substring '$needle' in output"
  fi
}

echo "=================================================================="
echo "AC1 — V signature matches CCNF"
echo "  V : Contract × Receipt → ValidationVerdict (no separate evidence input)"
echo "=================================================================="

# Receipt-only invocation works (refs are bound into the receipt).
out=$("$V" --receipt schemas/cdd/fixtures/valid-generic-receipt.yaml --structural-only --json 2>&1)
rc=$?
assert_exit "AC1.a receipt-only invocation succeeds" 0 $rc
assert_contains "AC1.b ValidationVerdict emitted" '"result"' "$out"
assert_contains "AC1.c validator_identity pinned" '"validator_identity": "cnos.cdd.validate_receipt"' "$out"

# Receipt without protocol_id is dispatch-rejected.
tmpfile=$(mktemp /tmp/v-no-protocol-XXXXXX.yaml)
cat > "$tmpfile" <<'EOF'
validation:
  verdict: PASS
  failed_predicates: []
  warnings: []
  provenance:
    validator_identity: x
    validator_version: x
    checked_at: x
    input_refs: {contract_ref: x, receipt_ref: x, evidence_root_ref: x}
boundary_decision:
  actor: x
  action: accept
transmissibility: accepted
protocol_gap_count: 0
protocol_gap_refs: []
evidence_refs: {}
EOF
out=$("$V" --receipt "$tmpfile" --structural-only 2>&1)
rc=$?
rm -f "$tmpfile"
assert_exit "AC1.d missing protocol_id rejected" 1 $rc
assert_contains "AC1.e diagnostic names missing protocol_id" "missing_protocol_id" "$out"

echo ""
echo "=================================================================="
echo "AC2 — protocol_id dispatch"
echo "  CDS receipt → schemas/cds/  |  CDR receipt → schemas/cdr/  |  generic → schemas/cdd/"
echo "=================================================================="

# Generic receipt PASSes against schemas/cdd/.
out=$("$V" --receipt schemas/cdd/fixtures/valid-generic-receipt.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC2.a generic receipt PASSes against #Receipt" 0 $rc

# CDS receipt PASSes against schemas/cds/.
out=$("$V" --receipt schemas/cds/fixtures/valid-receipt.yaml 2>&1)
rc=$?
assert_exit "AC2.b CDS receipt PASSes against #CDSReceipt" 0 $rc

# CDR receipt PASSes against schemas/cdr/.
out=$("$V" --receipt schemas/cdr/fixtures/valid-cdr-receipt.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC2.c CDR receipt PASSes against #CDRReceipt" 0 $rc

# Mismatched protocol_id (declares cds but missing CDS-required keys) FAILs
# with diagnostic naming the missing fields.
out=$("$V" --receipt schemas/cds/fixtures/counterfeit-mismatched-protocol-id.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC2.d mismatched protocol_id FAILs" 1 $rc
assert_contains "AC2.e diagnostic names missing keys" "missing required CDS evidence_refs keys" "$out"

echo ""
echo "=================================================================="
echo "AC3 — V rejects missing boundary_decision"
echo "=================================================================="

out=$("$V" --receipt schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC3.a missing boundary_decision FAILs" 1 $rc
assert_contains "AC3.b diagnostic mentions boundary_decision" "boundary_decision" "$out"

echo ""
echo "=================================================================="
echo "AC4 — V rejects γ-preflight as authoritative"
echo "=================================================================="

out=$("$V" --receipt schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC4.a γ-preflight-only FAILs" 1 $rc
# CUE rejects the missing boundary_decision (same structural rule); the
# γ-preflight non-authoritativeness manifests as the receipt being incomplete.
assert_contains "AC4.b diagnostic cites schema rule" "boundary_decision" "$out"

echo ""
echo "=================================================================="
echo "AC5 — ValidationVerdict ≠ BoundaryDecision (override does not rewrite)"
echo "=================================================================="

out=$("$V" --receipt schemas/cdd/fixtures/invalid-override-masks-verdict.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC5.a override-masks-verdict FAILs" 1 $rc
assert_contains "AC5.b CUE flags transmissibility mismatch" "transmissibility" "$out"

# C3 rule independently catches the deeper case: outer validation.verdict
# was rewritten while override block kept the original.
out=$("$V" --receipt schemas/cds/fixtures/counterfeit-override-rewrite.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC5.c override-rewrites-verdict FAILs" 1 $rc
assert_contains "AC5.d C3 names the rewrite signal" "counterfeit.override_does_not_rewrite" "$out"

echo ""
echo "=================================================================="
echo "AC6 — Evidence-graph inputs bound into receipt"
echo "  V dereferences refs through the receipt; non-existent paths FAIL."
echo "=================================================================="

# Valid CDS receipt with real on-disk closure-record refs → PASS.
out=$("$V" --receipt schemas/cds/fixtures/valid-receipt.yaml 2>&1)
rc=$?
assert_exit "AC6.a valid CDS receipt dereferences PASSes" 0 $rc

# Counterfeit-evidence-missing → FAIL (paths point at nonexistent files).
out=$("$V" --receipt schemas/cds/fixtures/counterfeit-evidence-missing.yaml 2>&1)
rc=$?
assert_exit "AC6.b evidence-missing FAILs" 1 $rc
assert_contains "AC6.c diagnostic names unresolved ref" "counterfeit.evidence_ref_unresolved" "$out"

echo ""
echo "=================================================================="
echo "AC7 — ValidationVerdict JSON structure"
echo "=================================================================="

json_out=$("$V" --receipt schemas/cds/fixtures/valid-receipt.yaml --json)
rc=$?
assert_exit "AC7.a --json invocation succeeds" 0 $rc
echo "$json_out" | jq . > /dev/null 2>&1
rc=$?
assert_exit "AC7.b output is parseable JSON" 0 $rc
for key in result failed_predicates warnings provenance; do
  if echo "$json_out" | jq -e ".${key}" > /dev/null 2>&1; then
    ac_pass "AC7.c[$key] key present"
  else
    ac_fail "AC7.c[$key] key present" "missing key '$key'"
  fi
done
result=$(echo "$json_out" | jq -r .result)
if [[ "$result" == "PASS" || "$result" == "FAIL" ]]; then
  ac_pass "AC7.d result is PASS|FAIL"
else
  ac_fail "AC7.d result is PASS|FAIL" "got '$result'"
fi
v_id=$(echo "$json_out" | jq -r .provenance.validator_identity)
if [[ "$v_id" == "cnos.cdd.validate_receipt" ]]; then
  ac_pass "AC7.e validator_identity pinned"
else
  ac_fail "AC7.e validator_identity pinned" "got '$v_id'"
fi

echo ""
echo "=================================================================="
echo "AC8 — Counterfeit-receipt rejection"
echo "  C1 actor-separation | C2 verdict-precedes-merge | C3 override-rewrite"
echo "=================================================================="

# C1 — actor collision.
out=$("$V" --receipt schemas/cds/fixtures/counterfeit-actor-collision.yaml 2>&1)
rc=$?
assert_exit "AC8.a C1 actor-collision FAILs" 1 $rc
assert_contains "AC8.b diagnostic names C1" "counterfeit.actor_separation" "$out"

# C2 — merge predates β verdict.
out=$("$V" --receipt schemas/cds/fixtures/counterfeit-merge-precedes-verdict.yaml 2>&1)
rc=$?
assert_exit "AC8.c C2 merge-precedes-verdict FAILs" 1 $rc
assert_contains "AC8.d diagnostic names C2" "counterfeit.verdict_precedes_merge" "$out"

# C3 — override rewrites verdict.
out=$("$V" --receipt schemas/cds/fixtures/counterfeit-override-rewrite.yaml --structural-only 2>&1)
rc=$?
assert_exit "AC8.e C3 override-rewrite FAILs" 1 $rc
assert_contains "AC8.f diagnostic names C3" "counterfeit.override_does_not_rewrite" "$out"

echo ""
echo "=================================================================="
echo "Backward compatibility — existing cn-cdd-verify modes still work"
echo "=================================================================="

# --unreleased mode still drives the legacy artifact-presence checker.
out=$("$V" --unreleased 2>&1)
rc=$?
# Exit 0 OR 1 acceptable (cycles may have warnings/missing artifacts); the
# check is that --unreleased does NOT fall through to V dispatch.
if [[ $rc -eq 0 || $rc -eq 1 ]]; then
  ac_pass "BC.a --unreleased mode runs legacy checker"
else
  ac_fail "BC.a --unreleased mode runs legacy checker" "unexpected exit $rc"
fi
assert_contains "BC.b legacy summary format preserved" "Summary:" "$out"

echo ""
echo "=================================================================="
echo "Summary"
echo "=================================================================="
printf "  %s passed, %s failed\n" "$(green $PASS_CT)" "$(red $FAIL_CT)"
if [[ $FAIL_CT -gt 0 ]]; then
  echo ""
  echo "Failures:"
  for f in "${FAILURES[@]}"; do
    echo "  - $f"
  done
  exit 1
fi
exit 0
