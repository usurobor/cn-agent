#!/usr/bin/env bash
# Exercise the #662 State-A runner without a live model. The fake coh process
# owns only deterministic prompt/report plumbing; runner validation, atomic
# publication, and math execute exactly as they do on the pinned route.

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
revision="a0d39293a27cfe57b49dacff696345b1ee2cdb40"
cm_revision="$(git -C "$repo_root" rev-parse HEAD)"
test_root="$(mktemp -d)"
pinned_root="$test_root/r7"
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
source_runner="$repo_root/$cm_rel/runner/recursive-cell-runner.py"
fixtures="$repo_root/$cm_rel/runner/fixtures"
fake_coh="$fixtures/fake-coh.py"
common=(
  --coh "$fake_coh"
  --root "$pinned_root"
  --engine-revision "$revision"
  --cm-revision "$cm_revision"
  --timestamp "2026-07-19T00:00:00Z"
)
targets=(cc662-system cc662-l0 cc662-l1 cc662-l2 cc662-l3 cc662-l4)

emit_run() {
  local output="$1"
  PYTHONDONTWRITEBYTECODE=1 "$runner" emit "${common[@]}" --output "$output"
}

ingest_run() {
  local output="$1" responses="$2" invariants="$3"
  PYTHONDONTWRITEBYTECODE=1 "$runner" ingest "${common[@]}" \
    --output "$output" --responses "$responses" --invariants "$invariants"
}

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

assert_unpublished() {
  local output="$1" label="$2"
  if [[ -e "$output/publication" || -L "$output/publication" ]]; then
    echo "FAIL canonical publication appeared after refusal: $label" >&2
    exit 1
  fi
  if find "$output" -maxdepth 1 -name '.publication-stage-*' -print -quit | grep -q .; then
    echo "FAIL staging directory survived refusal: $label" >&2
    exit 1
  fi
  if find "$output" -maxdepth 1 \( -name '.publication-lock' -o -name '.emission-lock' -o -name '.emission-stage-*' \) -print -quit | grep -q .; then
    echo "FAIL run lock or staging path survived refusal: $label" >&2
    exit 1
  fi
}

publication_digest() {
  local publication="$1"
  (
    cd "$publication"
    find . -type f -print | sort | while IFS= read -r path; do
      printf '%s\0' "$path"
      sha256sum "$path" | awk '{print $1}'
    done
  ) | sha256sum | awk '{print $1}'
}

pass_output="$test_root/pass"
[[ -d "$fixtures/responses-pass" && ! -L "$fixtures/responses-pass" ]]
for target in "${targets[@]}"; do
  [[ -f "$fixtures/responses-pass/$target.json" && ! -L "$fixtures/responses-pass/$target.json" ]]
done
[[ -f "$fixtures/invariants-pass.json" && ! -L "$fixtures/invariants-pass.json" ]]
emit_run "$pass_output"

# The authoritative output-root-relative emission declaration matches the
# complete fresh tree, README assessor path, and self-contained smoke snippet.
declared_emits=(
  emission/prompts/cc662-system.md
  emission/prompts/cc662-l0.md
  emission/prompts/cc662-l1.md
  emission/prompts/cc662-l2.md
  emission/prompts/cc662-l3.md
  emission/prompts/cc662-l4.md
  emission/invariant-assessment-prompt.md
  emission/prompt-digests.json
)
for path in "${declared_emits[@]}"; do
  [[ -f "$pass_output/$path" ]] || {
    echo "FAIL declared emission path missing: $path" >&2
    exit 1
  }
done
[[ "$(find "$pass_output/emission" -type f | wc -l)" -eq "${#declared_emits[@]}" ]]
expected_emits='    emits: ["emission/prompts/cc662-system.md", "emission/prompts/cc662-l0.md", "emission/prompts/cc662-l1.md", "emission/prompts/cc662-l2.md", "emission/prompts/cc662-l3.md", "emission/prompts/cc662-l4.md", "emission/invariant-assessment-prompt.md", "emission/prompt-digests.json"]'
grep -Fxq "$expected_emits" "$repo_root/$cm_rel/SKILL.md"
grep -Fq '`$RUN_ROOT/emission/invariant-assessment-prompt.md`' "$repo_root/$cm_rel/README.md"
sed -n '/^## Path-base boundary$/,$p' "$repo_root/$cm_rel/README.md" |
  grep -Fxq 'CM=src/packages/cnos.cdd/skills/cdd/measure/recursive-cell'

# Positive standard witness: one nonempty card omits optional secondary_axes.
jq -e '
  (.defect_cards | length) == 1 and
  (.defect_cards[0].id == "OPTIONAL-OMITTED") and
  (.defect_cards[0] | has("secondary_axes") | not)
' "$fixtures/responses-pass/cc662-l1.json" >/dev/null

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
  "$pass_output/emission/prompt-digests.json" >/dev/null
while IFS=$'\t' read -r relative expected; do
  actual="$(sha256sum "$pass_output/emission/$relative" | awk '{print $1}')"
  [[ "$actual" == "$expected" ]] || {
    echo "FAIL prompt digest mismatch: $relative" >&2
    exit 1
  }
done < <(jq -r '.targets[] | [.path, .sha256] | @tsv' "$pass_output/emission/prompt-digests.json")
assessment_path="$(jq -r '.invariant_assessment.path' "$pass_output/emission/prompt-digests.json")"
assessment_sha="$(jq -r '.invariant_assessment.sha256' "$pass_output/emission/prompt-digests.json")"
[[ "$(sha256sum "$pass_output/emission/$assessment_path" | awk '{print $1}')" == "$assessment_sha" ]]

FAKE_COH_REQUIRE_SNAPSHOT=1 \
  ingest_run "$pass_output" "$fixtures/responses-pass" "$fixtures/invariants-pass.json"
publication="$pass_output/publication"
result="$publication/recursive-cell-run.json"
success="$publication/publication-success.json"

jq -e --arg cm_revision "$cm_revision" '
  .schema == "cnos-recursive-cell-run/0.2" and
  .hard_invariant_gate.passed == true and
  (.hard_invariant_gate.items | map(.id)) ==
    ["H01","H02","H03","H04","H05","H06","H07","H08","H09","H10","H11","H12","H13"] and
  .bottleneck == {
    "axis":"alpha",
    "level":"L4",
    "tie_break":"lowest-level-index-then-TSC-bottleneck-axis"
  } and
  (.cross_level.levels[] | select(.level == "L4") |
    .bottleneck_axis == "alpha" and .beta < .alpha) and
  .disposition == {
    "reason":"zero-standing-calibration-source",
    "standing":"none",
    "value":"hold"
  } and
  (.targets | map(.target)) ==
    ["cc662-system","cc662-l0","cc662-l1","cc662-l2","cc662-l3","cc662-l4"] and
  (.targets | map(.report_path)) ==
    ["reports/cc662-system.json","reports/cc662-l0.json","reports/cc662-l1.json","reports/cc662-l2.json","reports/cc662-l3.json","reports/cc662-l4.json"] and
  (.cross_level.C_sigma_cross_num > 0.566762129988) and
  (.cross_level.C_sigma_cross_num < 0.566762129989) and
  .provenance.response_sha256 == $responses and
  (.provenance.skill_sha256 | test("^[0-9a-f]{64}$")) and
  (.provenance.cm_authority_bundle_sha256 | test("^[0-9a-f]{64}$")) and
  .provenance.declared_cm_revision == $cm_revision
' --argjson responses "$(jq '.response_sha256' "$fixtures/invariants-pass.json")" "$result" >/dev/null

jq -e '
  .schema == "cnos-recursive-cell-publication/0.2" and
  .target == "cc662-r7-recursive" and .status == "complete" and
  .result.path == "recursive-cell-run.json" and
  (.reports | map(.target)) ==
    ["cc662-system","cc662-l0","cc662-l1","cc662-l2","cc662-l3","cc662-l4"] and
  (.reports | map(.path)) ==
    ["reports/cc662-system.json","reports/cc662-l0.json","reports/cc662-l1.json","reports/cc662-l2.json","reports/cc662-l3.json","reports/cc662-l4.json"] and
  .emission.prompt_manifest.path == "emission/prompt-digests.json" and
  .emission.invariant_assessment_prompt.path == "emission/invariant-assessment-prompt.md" and
  (.emission.prompts | map(.path)) ==
    ["emission/prompts/cc662-system.md","emission/prompts/cc662-l0.md","emission/prompts/cc662-l1.md","emission/prompts/cc662-l2.md","emission/prompts/cc662-l3.md","emission/prompts/cc662-l4.md"] and
  (.inputs.responses | map(.path)) ==
    ["inputs/responses/cc662-system.json","inputs/responses/cc662-l0.json","inputs/responses/cc662-l1.json","inputs/responses/cc662-l2.json","inputs/responses/cc662-l3.json","inputs/responses/cc662-l4.json"] and
  .inputs.invariant_assessment.path == "inputs/invariant-assessment.json"
' "$success" >/dev/null
[[ "$(sha256sum "$result" | awk '{print $1}')" == "$(jq -r '.result.sha256' "$success")" ]]
while IFS=$'\t' read -r path expected; do
  [[ "$(sha256sum "$publication/$path" | awk '{print $1}')" == "$expected" ]]
done < <(jq -r '[.reports[], .emission.prompt_manifest, .emission.invariant_assessment_prompt, .emission.prompts[], .inputs.responses[], .inputs.invariant_assessment] | .[] | [.path,.sha256] | @tsv' "$success")

# Emission is also immutable and requires a unique fresh run root.
before_emission="$(publication_digest "$pass_output/emission")"
expect_refusal "successful emission reuse" emit_run "$pass_output"
[[ "$(publication_digest "$pass_output/emission")" == "$before_emission" ]]
emit_locked_output="$test_root/emit-locked"
mkdir -p "$emit_locked_output/.emission-lock"
expect_refusal "emission lock contention" emit_run "$emit_locked_output"
rmdir "$emit_locked_output/.emission-lock"
assert_unpublished "$emit_locked_output" "emission lock contention"

# A completed publication is immutable and reuse refuses without changing it.
before_reuse="$(publication_digest "$publication")"
expect_refusal "successful output reuse" \
  ingest_run "$pass_output" "$fixtures/responses-pass" "$fixtures/invariants-pass.json"
[[ "$(publication_digest "$publication")" == "$before_reuse" ]]

# A symlink cannot claim or redirect the immutable publication path.
symlink_output="$test_root/symlink-publication"
emit_run "$symlink_output"
ln -s "$test_root/redirect-target" "$symlink_output/publication"
expect_refusal "symlink publication path" \
  ingest_run "$symlink_output" "$fixtures/responses-pass" "$fixtures/invariants-pass.json"
[[ -L "$symlink_output/publication" ]]
[[ ! -e "$test_root/redirect-target" ]]

# The output root and every existing ancestor are lexical custody boundaries:
# neither may redirect canonical emission before the runner creates anything.
symlink_root_target="$test_root/symlink-root-target"
mkdir "$symlink_root_target"
symlink_root="$test_root/symlink-root"
ln -s "$symlink_root_target" "$symlink_root"
expect_refusal "symlink output root" emit_run "$symlink_root"
[[ -L "$symlink_root" ]]
[[ -z "$(find "$symlink_root_target" -mindepth 1 -print -quit)" ]]

symlink_ancestor_target="$test_root/symlink-ancestor-target"
mkdir "$symlink_ancestor_target"
symlink_ancestor="$test_root/symlink-ancestor"
ln -s "$symlink_ancestor_target" "$symlink_ancestor"
expect_refusal "symlink output ancestor" emit_run "$symlink_ancestor/run"
[[ -L "$symlink_ancestor" ]]
[[ -z "$(find "$symlink_ancestor_target" -mindepth 1 -print -quit)" ]]

# Every one of the six external response paths is preflighted before staging.
for linked_target in "${targets[@]}"; do
  linked_responses="$test_root/symlink-response-$linked_target"
  mkdir "$linked_responses"
  for target in "${targets[@]}"; do
    if [[ "$target" == "$linked_target" ]]; then
      ln -s "$fixtures/responses-pass/$target.json" "$linked_responses/$target.json"
    else
      cp "$fixtures/responses-pass/$target.json" "$linked_responses/$target.json"
    fi
  done
  linked_output="$test_root/symlink-response-output-$linked_target"
  emit_run "$linked_output"
  expect_refusal "symlink response $linked_target" \
    ingest_run "$linked_output" "$linked_responses" "$fixtures/invariants-pass.json"
  [[ -L "$linked_responses/$linked_target.json" ]]
  assert_unpublished "$linked_output" "symlink response $linked_target"
done

# A symlinked response-directory component is also refused before staging.
regular_response_parent="$test_root/regular-response-parent"
mkdir "$regular_response_parent"
cp "$fixtures"/responses-pass/*.json "$regular_response_parent/"
linked_response_parent="$test_root/linked-response-parent"
ln -s "$regular_response_parent" "$linked_response_parent"
linked_parent_output="$test_root/linked-response-parent-output"
emit_run "$linked_parent_output"
expect_refusal "symlink response ancestor" \
  ingest_run "$linked_parent_output" "$linked_response_parent" "$fixtures/invariants-pass.json"
assert_unpublished "$linked_parent_output" "symlink response ancestor"

# The invariant input has the same regular-file, no-symlink custody boundary.
linked_invariants="$test_root/linked-invariants.json"
ln -s "$fixtures/invariants-pass.json" "$linked_invariants"
linked_invariant_output="$test_root/linked-invariant-output"
emit_run "$linked_invariant_output"
expect_refusal "symlink invariant assessment" \
  ingest_run "$linked_invariant_output" "$fixtures/responses-pass" "$linked_invariants"
[[ -L "$linked_invariants" ]]
assert_unpublished "$linked_invariant_output" "symlink invariant assessment"

# A competing conforming publisher holds an atomic lock and is never raced.
locked_output="$test_root/locked-publication"
emit_run "$locked_output"
mkdir "$locked_output/.publication-lock"
expect_refusal "publication lock contention" \
  ingest_run "$locked_output" "$fixtures/responses-pass" "$fixtures/invariants-pass.json"
rmdir "$locked_output/.publication-lock"
assert_unpublished "$locked_output" "publication lock contention"

# Failed invariant publication is valid but nonaccepting.
fail_responses="$test_root/fail-responses"
mkdir -p "$fail_responses"
cp "$fixtures"/responses-pass/*.json "$fail_responses/"
cp "$fixtures/cc662-l3-fail.json" "$fail_responses/cc662-l3.json"
fail_output="$test_root/fail"
emit_run "$fail_output"
ingest_run "$fail_output" "$fail_responses" "$fixtures/invariants-fail.json"
jq -e '
  .hard_invariant_gate.passed == false and
  .disposition.value == "request_planning"
' "$fail_output/publication/recursive-cell-run.json" >/dev/null

# Unknown is typed but nonaccepting and must carry a next MCA.
unknown_invariants="$test_root/invariants-unknown.json"
jq '
  (.items[] | select(.id == "H05") | .status) = "unknown" |
  (.items[] | select(.id == "H05") | .next_mca) = "Resolve H05 evidence."
' "$fixtures/invariants-pass.json" >"$unknown_invariants"
unknown_output="$test_root/unknown"
emit_run "$unknown_output"
ingest_run "$unknown_output" "$fixtures/responses-pass" "$unknown_invariants"
jq -e '
  .hard_invariant_gate.passed == false and .disposition.value == "hold"
' "$unknown_output/publication/recursive-cell-run.json" >/dev/null

# Refuse a failed invariant without the scoped systemic card.
scoped_output="$test_root/scoped-refusal"
emit_run "$scoped_output"
expect_refusal "failed invariant without scoped card" \
  ingest_run "$scoped_output" "$fixtures/responses-pass" "$fixtures/invariants-fail.json"
assert_unpublished "$scoped_output" "failed invariant without scoped card"

# Refuse incomplete per-level ingestion.
missing_responses="$test_root/missing-responses"
mkdir -p "$missing_responses"
for target in cc662-system cc662-l0 cc662-l1 cc662-l2 cc662-l3; do
  cp "$fixtures/responses-pass/$target.json" "$missing_responses/"
done
missing_output="$test_root/missing"
emit_run "$missing_output"
expect_refusal "missing L4 response" \
  ingest_run "$missing_output" "$missing_responses" "$fixtures/invariants-pass.json"
assert_unpublished "$missing_output" "missing L4 response"

# Refuse valid witness bytes when the assessment response map is stale.
replayed_responses="$test_root/replayed-responses"
mkdir -p "$replayed_responses"
cp "$fixtures"/responses-pass/*.json "$replayed_responses/"
jq '.summary = "Different but otherwise valid witness bytes."' \
  "$fixtures/responses-pass/cc662-l2.json" >"$replayed_responses/cc662-l2.json"
replay_output="$test_root/replay"
emit_run "$replay_output"
expect_refusal "assessment replayed against different responses" \
  ingest_run "$replay_output" "$replayed_responses" "$fixtures/invariants-pass.json"
assert_unpublished "$replay_output" "assessment response digest mismatch"

# Refuse malformed non-empty next_fixes shapes.
for variant in invalid-axis extra-field missing-fix; do
  malformed="$test_root/next-fixes-$variant"
  mkdir -p "$malformed"
  cp "$fixtures"/responses-pass/*.json "$malformed/"
  case "$variant" in
    invalid-axis)
      jq '.next_fixes = [{"axis":"delta","fix":"Bad axis."}]' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
    extra-field)
      jq '.next_fixes = [{"axis":"beta","fix":"Repair.","extra":true}]' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
    missing-fix)
      jq '.next_fixes = [{"axis":"beta"}]' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
  esac
  malformed_output="$test_root/malformed-next-fixes-$variant"
  emit_run "$malformed_output"
  expect_refusal "malformed next_fixes $variant" \
    ingest_run "$malformed_output" "$malformed" "$fixtures/invariants-pass.json"
  assert_unpublished "$malformed_output" "malformed next_fixes $variant"
done

# Refuse malformed optional-card shapes while omission remains valid above.
for variant in secondary-nonlist secondary-primary secondary-invalid secondary-duplicate extra-field; do
  malformed="$test_root/secondary-axes-$variant"
  mkdir -p "$malformed"
  cp "$fixtures"/responses-pass/*.json "$malformed/"
  case "$variant" in
    secondary-nonlist)
      jq '.defect_cards[0].secondary_axes = "gamma"' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
    secondary-primary)
      jq '.defect_cards[0].secondary_axes = ["beta"]' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
    secondary-invalid)
      jq '.defect_cards[0].secondary_axes = ["delta"]' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
    secondary-duplicate)
      jq '.defect_cards[0].secondary_axes = ["gamma","gamma"]' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
    extra-field)
      jq '.defect_cards[0].unexpected = true' \
        "$fixtures/responses-pass/cc662-l1.json" >"$malformed/cc662-l1.json" ;;
  esac
  malformed_output="$test_root/malformed-card-$variant"
  emit_run "$malformed_output"
  expect_refusal "malformed optional defect-card field $variant" \
    ingest_run "$malformed_output" "$malformed" "$fixtures/invariants-pass.json"
  assert_unpublished "$malformed_output" "malformed optional defect-card field $variant"
done

# Refuse a witness whose standard target identity moved.
moved_responses="$test_root/moved-responses"
mkdir -p "$moved_responses"
cp "$fixtures"/responses-pass/*.json "$moved_responses/"
jq '.target = "cc662-l3"' "$fixtures/responses-pass/cc662-l2.json" \
  >"$moved_responses/cc662-l2.json"
moved_output="$test_root/moved"
emit_run "$moved_output"
expect_refusal "moved witness target" \
  ingest_run "$moved_output" "$moved_responses" "$fixtures/invariants-pass.json"
assert_unpublished "$moved_output" "moved witness target"

# Refuse prompt mutation between emit and ingest.
mutated_output="$test_root/mutated"
emit_run "$mutated_output"
printf '\nmutation\n' >>"$mutated_output/emission/prompts/cc662-l2.md"
expect_refusal "mutated emitted prompt" \
  ingest_run "$mutated_output" "$fixtures/responses-pass" "$fixtures/invariants-pass.json"
assert_unpublished "$mutated_output" "mutated prompt"

# A forced mid-target engine failure leaves no canonical or staging output.
mid_output="$test_root/mid-failure"
emit_run "$mid_output"
expect_refusal "mid-ingest target failure" \
  env FAKE_COH_FAIL_TARGET=cc662-l2 PYTHONDONTWRITEBYTECODE=1 \
  "$runner" ingest "${common[@]}" --output "$mid_output" \
  --responses "$fixtures/responses-pass" --invariants "$fixtures/invariants-pass.json"
assert_unpublished "$mid_output" "mid-ingest target failure"

# Refuse unshaped source provenance.
expect_refusal "mutable CM revision" \
  "$runner" emit --coh "$fake_coh" --root "$pinned_root" \
  --output "$test_root/bad-revision" --engine-revision "$revision" \
  --cm-revision "mutable-main" --timestamp "2026-07-19T00:00:00Z"

# A final CUE failure after all six reports also leaves no publication.
schema_failure_output="$test_root/schema-failure"
emit_run "$schema_failure_output"
printf '\n#RecursiveCellRun: {schema: "impossible-r4-test"}\n' \
  >>"$pinned_root/$cm_rel/runner/recursive-cell-run.schema.cue"
expect_refusal "final CUE validation failure" \
  ingest_run "$schema_failure_output" "$fixtures/responses-pass" "$fixtures/invariants-pass.json"
assert_unpublished "$schema_failure_output" "final CUE validation failure"

# Refuse an active runner that does not match the target-root authority copy.
printf '\n# authority mutation\n' >>"$runner"
expect_refusal "active/authority runner mismatch" \
  "$source_runner" emit "${common[@]}" --output "$test_root/runner-mismatch"

echo "PASS recursive-cell State-A runner: typed witnesses, atomic publication, aggregate, gate, provenance, refusals"
