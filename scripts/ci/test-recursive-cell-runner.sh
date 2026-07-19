#!/usr/bin/env bash
# Exercise the #662 State-A runner without a live model. The fake coh process
# owns only deterministic prompt/report plumbing; runner validation and math
# execute exactly as they do with the pinned TSC engine route.

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
revision="a0d39293a27cfe57b49dacff696345b1ee2cdb40"
cm_revision="$(git -C "$repo_root" rev-parse HEAD)"
test_root="$(mktemp -d)"
pinned_root="$test_root/r7"
output="$test_root/output"
refusal_log="$test_root/refusal.log"

cleanup() {
  git -C "$repo_root" worktree remove --force "$pinned_root" >/dev/null 2>&1 || true
  rm -rf "$test_root"
}
trap cleanup EXIT

git -C "$repo_root" cat-file -e "${revision}^{commit}"
git -C "$repo_root" worktree add --detach "$pinned_root" "$revision" >/dev/null

cm_rel="src/packages/cnos.cdd/skills/cdd/measure/recursive-cell"
mkdir -p "$pinned_root/$(dirname "$cm_rel")"
cp -a "$repo_root/$cm_rel" "$pinned_root/$(dirname "$cm_rel")/"

runner="$pinned_root/$cm_rel/runner/recursive-cell-runner.py"
fixtures="$repo_root/$cm_rel/runner/fixtures"
fake_coh="$fixtures/fake-coh.py"
common=(
  --coh "$fake_coh"
  --root "$pinned_root"
  --output "$output"
  --engine-revision "$revision"
  --cm-revision "$cm_revision"
  --timestamp "2026-07-19T00:00:00Z"
)

expect_refusal() {
  local label="$1"
  shift
  if PYTHONDONTWRITEBYTECODE=1 "$@" >"$refusal_log" 2>&1; then
    echo "FAIL runner accepted refusal fixture: $label" >&2
    exit 1
  fi
  grep -Fq "REFUSE:" "$refusal_log" || {
    echo "FAIL runner refusal was not explicit: $label" >&2
    cat "$refusal_log" >&2
    exit 1
  }
}

PYTHONDONTWRITEBYTECODE=1 "$runner" emit "${common[@]}"

# Mirror the real coh path contract: authority arguments are root-relative.
if "$fake_coh" --target cc662-system \
  --registry "$pinned_root/$cm_rel/calibration/662/registry.tsc" \
  --instruction "$cm_rel/INSTRUCTION.md" --root "$pinned_root" \
  --mode llm --emit-prompt "$test_root/absolute-path.prompt" >/dev/null 2>&1; then
  echo "FAIL fake coh accepted an absolute authority path" >&2
  exit 1
fi

# Independently check every prompt digest and the generated assessment prompt.
jq -e '.schema == "cnos-recursive-cell-prompts/0.1" and (.targets | length) == 6' \
  "$output/prompt-digests.json" >/dev/null
while IFS=$'\t' read -r relative expected; do
  actual="$(sha256sum "$output/$relative" | awk '{print $1}')"
  [[ "$actual" == "$expected" ]] || {
    echo "FAIL prompt digest mismatch: $relative" >&2
    exit 1
  }
done < <(jq -r '.targets[] | [.path, .sha256] | @tsv' "$output/prompt-digests.json")
assessment_path="$(jq -r '.invariant_assessment.path' "$output/prompt-digests.json")"
assessment_sha="$(jq -r '.invariant_assessment.sha256' "$output/prompt-digests.json")"
[[ "$(sha256sum "$output/$assessment_path" | awk '{print $1}')" == "$assessment_sha" ]]

PYTHONDONTWRITEBYTECODE=1 "$runner" ingest "${common[@]}" \
  --responses "$fixtures/responses-pass" \
  --invariants "$fixtures/invariants-pass.json"

jq -e '
  .hard_invariant_gate.passed == true and
  (.hard_invariant_gate.items | map(.id)) ==
    ["H01","H02","H03","H04","H05","H06","H07","H08","H09","H10","H11","H12","H13"] and
  .bottleneck == {
    "axis":"beta",
    "level":"L4",
    "tie_break":"lowest-level-index-then-alpha-beta-gamma"
  } and
  .disposition == {
    "reason":"zero-standing-calibration-source",
    "standing":"none",
    "value":"hold"
  } and
  (.targets | map(.target)) ==
    ["cc662-system","cc662-l0","cc662-l1","cc662-l2","cc662-l3","cc662-l4"] and
  (.cross_level.C_sigma_cross_num > 0.566762129988) and
  (.cross_level.C_sigma_cross_num < 0.566762129989) and
  (.provenance.runner_sha256 | test("^[0-9a-f]{64}$")) and
  (.provenance.output_schema_sha256 | test("^[0-9a-f]{64}$")) and
  (.provenance.coh_executable_sha256 | test("^[0-9a-f]{64}$")) and
  (.provenance.cm_authority_bundle_sha256 | test("^[0-9a-f]{64}$")) and
  .provenance.declared_cm_revision == $cm_revision
' --arg cm_revision "$cm_revision" "$output/recursive-cell-run.json" >/dev/null

fail_responses="$test_root/fail-responses"
mkdir -p "$fail_responses"
cp "$fixtures"/responses-pass/*.json "$fail_responses/"
cp "$fixtures/cc662-l3-fail.json" "$fail_responses/cc662-l3.json"
PYTHONDONTWRITEBYTECODE=1 "$runner" ingest "${common[@]}" \
  --responses "$fail_responses" \
  --invariants "$fixtures/invariants-fail.json"
jq -e '
  .hard_invariant_gate.passed == false and
  .disposition.value == "request_planning"
' "$output/recursive-cell-run.json" >/dev/null

# Unknown is typed but nonaccepting and must carry a next MCA.
unknown_invariants="$test_root/invariants-unknown.json"
jq '
  (.items[] | select(.id == "H05") | .status) = "unknown" |
  (.items[] | select(.id == "H05") | .next_mca) = "Resolve H05 evidence."
' "$fixtures/invariants-pass.json" >"$unknown_invariants"
PYTHONDONTWRITEBYTECODE=1 "$runner" ingest "${common[@]}" \
  --responses "$fixtures/responses-pass" \
  --invariants "$unknown_invariants"
jq -e '
  .hard_invariant_gate.passed == false and
  .disposition.value == "hold"
' "$output/recursive-cell-run.json" >/dev/null

# Refuse a failed invariant whose declared L3/gamma witness has no matching
# systemic defect card.
expect_refusal "failed invariant without scoped card" \
  "$runner" ingest "${common[@]}" \
  --responses "$fixtures/responses-pass" \
  --invariants "$fixtures/invariants-fail.json"

# Refuse incomplete per-level ingestion.
missing_responses="$test_root/missing-responses"
mkdir -p "$missing_responses"
for target in cc662-system cc662-l0 cc662-l1 cc662-l2 cc662-l3; do
  cp "$fixtures/responses-pass/$target.json" "$missing_responses/"
done
expect_refusal "missing L4 response" \
  "$runner" ingest "${common[@]}" \
  --responses "$missing_responses" \
  --invariants "$fixtures/invariants-pass.json"

# Refuse a witness whose standard target identity moved.
moved_responses="$test_root/moved-responses"
mkdir -p "$moved_responses"
cp "$fixtures"/responses-pass/*.json "$moved_responses/"
jq '.target = "cc662-l3"' "$fixtures/responses-pass/cc662-l2.json" \
  >"$moved_responses/cc662-l2.json"
expect_refusal "moved witness target" \
  "$runner" ingest "${common[@]}" \
  --responses "$moved_responses" \
  --invariants "$fixtures/invariants-pass.json"

# Refuse prompt mutation between emit and ingest, then restore the bundle.
printf '\nmutation\n' >>"$output/prompts/cc662-l2.md"
expect_refusal "mutated emitted prompt" \
  "$runner" ingest "${common[@]}" \
  --responses "$fixtures/responses-pass" \
  --invariants "$fixtures/invariants-pass.json"
PYTHONDONTWRITEBYTECODE=1 "$runner" emit "${common[@]}" >/dev/null

# Refuse unshaped source provenance.
bad_revision_common=(
  --coh "$fake_coh"
  --root "$pinned_root"
  --output "$test_root/bad-revision"
  --engine-revision "$revision"
  --cm-revision "mutable-main"
  --timestamp "2026-07-19T00:00:00Z"
)
expect_refusal "mutable CM revision" "$runner" emit "${bad_revision_common[@]}"

# Refuse an active runner that does not match the target-root authority copy.
source_runner="$repo_root/$cm_rel/runner/recursive-cell-runner.py"
printf '\n# authority mutation\n' >>"$runner"
expect_refusal "active/authority runner mismatch" \
  "$source_runner" emit "${common[@]}"

echo "PASS recursive-cell State-A runner: six prompts/reports, aggregate, gate, disposition, provenance, refusals"
