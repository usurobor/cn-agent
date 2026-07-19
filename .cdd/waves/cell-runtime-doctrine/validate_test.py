#!/usr/bin/env python3
# wave-revision: R7
"""
Executable adversarial-mutation harness for validate.py (wave cnos#671 R7).

Materializes TWENTY-NINE adversarial mutations, each in its own temp copy of the wave tree, HONESTLY
re-pinning every changed contract_sha256 and the oracle-registry content hash (via resync_all), runs
validate.py on each, and asserts EACH exits non-zero FOR ITS OWN NAMED PREDICATE — while the UNMODIFIED
tree exits 0.

The SIX original negatives (β used the first five to false-pass the R2 validator; the sixth is the
mechanically-verifiable oracle-ownership negative, against the registry):
  1. wrong intent id            -> (b): intent_ref.id != the actual intent object id
  2. nonsense cell.class        -> (a): enum cell.class ∈ {working,planning,cohering}
  3. nonexistent artifact path  -> (b): git cannot resolve <commit>:<path> (DOES-NOT-EXIST.md)
  4. tautological whole-wave    -> (g): predicate dependency graph has a CYCLE
  5. flipped fixture expected   -> (g): fixture computed result != authored `expected`
  6. placeholder oracle operand -> (h): an mv registry entry's checker left as a placeholder (non-concrete)

The FIVE mechanically-verifiable bijection mutations (check (h)):
  7. remove one mv registry entry (acceptance predicate + mv projection row left intact) -> (h): a
     mechanically-verifiable predicate in the projection with NO registry owner entry (MISSING).
  8. duplicate one entry under another owner                                          -> (h): DUPLICATE owner
  9. reclassify one mv entry to evidenced                                             -> (j): registry↔projection classification mismatch
 10. leave one mechanically-verifiable acceptance predicate unowned                   -> (h): MISSING entry
 11. add a registry owner entry absent from the contract's acceptance predicates      -> (h): extra owner entry

The FOUR R6 TOTAL-classification mutations (check (j)):
 12. COORDINATED OMISSION — remove a predicate from the registry + BOTH md surfaces + adjust the mv count,
     LEAVE it in the child contract. checks (h) and (i) STILL PASS (this is exactly the R5-breaking
     mutation); (j) MUST fail: an acceptance predicate absent from the assurance registry (UNCLASSIFIED).
 13. DOUBLE-CLASSIFY one child predicate (a second assurance entry)                   -> (j): DOUBLE-CLASSIFIED
 14. PHANTOM registry entry for a predicate no contract declares                      -> (j): PHANTOM
 15. category-without-required-fields (change a category, omit its required fields)   -> (j): missing required field

The SEVEN R7 cn.wave.v1 ENVELOPE mutations (check (k)) — Blocker-1; β mutated the wave object and the
selective pre-R7 extractor passed them all:
 16. schema: nonsense                 -> (k): wave.schema != cn.wave.v1
 17. stop_conditions: []              -> (k): STOP conditions must be a NON-EMPTY seq
 18. gates: {}                        -> (k): wave.gates MISSING keys
 19. critical_path truncated [wc-2,wc-1] -> (k): critical_path does not terminate at the terminal node
 20. wc-2 node output_path docs/WRONG.md -> (k): node.output_path != contract.requested_output.path (parity)
 21. first edge kind: nonsense        -> (k): edge kind not in {depends_on}
 22. duplicate `schema:` key          -> (k): StrictSafeLoader rejects a DUPLICATE mapping key

The SEVEN R7 EVIDENCE-DERIVED-completion mutations (check (g)) — Blocker-2; β set child_complete to a
constant that ignored output/acceptance/V/receipt and the boolean-only model passed:
 23. child_complete.definition = constant `true` (ignores constituents) -> (g): definition IGNORES a required constituent
 24. flip requested_output_produced (claim child_complete: true)        -> (g): naked boolean not backed by constituents
 25. flip acceptance_all_pass                                           -> (g): naked boolean not backed by constituents
 26. flip class_v_pass                                                  -> (g): naked boolean not backed by constituents
 27. flip receipt_bound                                                 -> (g): naked boolean not backed by constituents
 28. remove a constituent from a record                                 -> (g): record missing constituent
 29. wc-5 completion: flip receipt_bound (claim true)                   -> (g): naked boolean (completion, not just readiness)

Run:  python3 validate_test.py   (exit 0 iff all 29 fail for their own predicate AND clean passes)
"""
import os, sys, re, shutil, tempfile, subprocess, hashlib, copy
import yaml

WAVE_DIR = os.path.dirname(os.path.abspath(__file__))
CHILDREN = ["wc-1", "wc-2", "wc-3a", "wc-3b", "wc-4", "wc-5"]

def _git_toplevel(path):
    r = subprocess.run(["git", "-C", path, "rev-parse", "--show-toplevel"], capture_output=True, text=True)
    if r.returncode != 0:
        print("FATAL: could not find git repo root for", path, file=sys.stderr); sys.exit(2)
    return r.stdout.strip()

REPO_ROOT = _git_toplevel(WAVE_DIR)

def sha256_file(path):
    with open(path, "rb") as f:
        return hashlib.sha256(f.read()).hexdigest()

def run_validator(tree_dir):
    env = dict(os.environ)
    env["WAVE_VALIDATOR_REPO_ROOT"] = REPO_ROOT           # git objects resolve against the real repo
    env.pop("WAVE_VALIDATOR_WAVE_DIR", None)              # WAVE_DIR defaults to dirname(validate.py) = tree_dir
    r = subprocess.run([sys.executable, os.path.join(tree_dir, "validate.py")],
                       capture_output=True, text=True, env=env)
    return r.returncode, r.stdout + r.stderr

def copy_tree():
    tmp = tempfile.mkdtemp(prefix="wave671-negfix-")
    dst = os.path.join(tmp, "wave")
    shutil.copytree(WAVE_DIR, dst)
    return tmp, dst

def edit(path, old, new, count=1):
    s = open(path).read()
    assert old in s, f"anchor not found in {os.path.basename(path)}: {old!r}"
    open(path, "w").write(s.replace(old, new, count))

def contract_path(tree, nid):
    return os.path.join(tree, "contracts", f"{nid}.cn-cell-contract-v1.yaml")

def resync_all(tree):
    """HONESTLY re-pin every hash: the oracle-registry content hash into each contract's registry ref,
    then each contract's SHA-256 into wave.cn-wave-v1.yaml. Isolates each mutation to its own check."""
    R = sha256_file(os.path.join(tree, "oracle-registry.yaml"))
    wave_p = os.path.join(tree, "wave.cn-wave-v1.yaml")
    wtext = open(wave_p).read()
    for nid in CHILDREN:
        cf = contract_path(tree, nid)
        ct = open(cf).read()
        ct = re.sub(r"(sha256:)[0-9a-f]{64}(@\.cdd/waves/cell-runtime-doctrine/oracle-registry\.yaml)",
                    r"\g<1>" + R + r"\g<2>", ct)
        open(cf, "w").write(ct)
        newsha = sha256_file(cf)
        wtext = re.sub(
            r'(contract_ref: "contracts/' + re.escape(nid) + r'\.cn-cell-contract-v1\.yaml"\n\s*contract_sha256: ")[0-9a-f]{64}(")',
            r"\g<1>" + newsha + r"\g<2>", wtext)
    open(wave_p, "w").write(wtext)

def load_registry(tree):
    return yaml.safe_load(open(os.path.join(tree, "oracle-registry.yaml")))

def dump_registry(tree, data):
    # safe_dump strips YAML comments; re-inject the `wave-revision:` marker line so a registry mutation
    # stays isolated to its intended check (a stripped marker would incidentally trip check (i)).
    body = yaml.safe_dump(data, sort_keys=False, width=100000, allow_unicode=True)
    with open(os.path.join(tree, "oracle-registry.yaml"), "w") as f:
        f.write("# wave-revision: R7\n" + body)

def find_entry(d, owner, pred):
    return next(e for e in d["assurance"] if e.get("owner") == owner and e.get("predicate") == pred)

def md_add_mv_projection_row(tree, row):
    edit(os.path.join(tree, "acceptance-oracles.md"), "\n\n---\n\n*Rows marked", "\n" + row + "\n\n---\n\n*Rows marked")

def md_add_complete_row(tree, row):
    edit(os.path.join(tree, "acceptance-oracles.md"), "\n\n---\n\n## Registry projection",
         "\n" + row + "\n\n---\n\n## Registry projection")

def md_remove_lines_containing(tree, substr, expect):
    p = os.path.join(tree, "acceptance-oracles.md")
    lines = open(p).read().splitlines(keepends=True)
    kept = [ln for ln in lines if substr not in ln]
    removed = len(lines) - len(kept)
    assert removed == expect, f"expected to remove {expect} md line(s) containing {substr!r}, removed {removed}"
    open(p, "w").write("".join(kept))

KEYPATH = "key_path_parity__template_and_worked_instance_normalize_to_identical_key_path_sets"
CM_TO_V = "cm_to_v_edge_pinned__v_consumes_an_immutable_cm_ref_as_a_named_input_the_frozen_contract_times_receipt_signature_is_extended_to_carry"
SCOPE_GUARD = "scope_guard__no_cm_internals_no_cell_or_wave_fsm_no_migration_authored_here"

# ---- the six original mutations --------------------------------------------------
def mut_wrong_intent_id(tree):
    edit(contract_path(tree, "wc-1"), 'id: "intent-2026-0719-671-cell-runtime-doctrine"',
         'id: "intent-WRONG-does-not-match-the-object"')

def mut_nonsense_class(tree):
    edit(contract_path(tree, "wc-1"), 'class: "working"', 'class: "nonsense"')

def mut_nonexistent_path(tree):
    edit(contract_path(tree, "wc-1"), 'path: "docs/architecture/CELL-RUNTIME.md"',
         'path: "docs/architecture/DOES-NOT-EXIST.md"')

def mut_tautology(tree):
    edit(os.path.join(tree, "wave.cn-wave-v1.yaml"),
         '      expr:\n        and:\n          - { all_over_N: "child_complete" }\n          - { pred: "wc5_complete" }',
         '      expr: { pred: "wave_complete" }')   # wave_complete := wave_complete (tautology)

def mut_flipped_expected(tree):
    edit(os.path.join(tree, "wave.cn-wave-v1.yaml"),
         'expected: { wc5_ready: true, wc5_complete: true, wave_complete: true }',
         'expected: { wc5_ready: true, wc5_complete: true, wave_complete: false }')

def mut_placeholder_oracle(tree):
    # Leave an mv registry entry's checker as a PLACEHOLDER operand (non-concrete): check (h) rejects it.
    d = load_registry(tree)
    for e in d["assurance"]:
        if e.get("checker", "").endswith("check_v_signature.py"):
            e["checker"] = "<checker-to-be-decided>"
            break
    dump_registry(tree, d)

# ---- the five mechanically-verifiable bijection mutations (check (h)) -------------
def mut_remove_entry(tree):
    # remove the wc-2 cm_to_v mv entry from the registry; leave its acceptance predicate + the mv
    # projection row intact -> an mv projection predicate with no registry owner (MISSING entry).
    d = load_registry(tree)
    d["assurance"] = [e for e in d["assurance"] if e.get("predicate") != CM_TO_V]
    dump_registry(tree, d)

def mut_duplicate_owner(tree):
    d = load_registry(tree)
    src = find_entry(d, "wc-1", KEYPATH)
    dup = copy.deepcopy(src)
    dup["owner"] = "wc-2"     # same mv predicate now owned by a second child -> mapping not singular
    d["assurance"].append(dup)
    dump_registry(tree, d)

def mut_reclassify(tree):
    d = load_registry(tree)
    find_entry(d, "wc-1", KEYPATH)["classification"] = "evidenced"   # registry vs md-projection mismatch
    dump_registry(tree, d)

SYN_UNOWNED = "synthetic_unowned_mechanically_verifiable_predicate_present_in_contract_absent_from_registry"
def mut_unowned_predicate(tree):
    # add an mv predicate to wc-1 acceptance + the mv projection, but NOT to the registry.
    edit(contract_path(tree, "wc-1"), "acceptance:\n  predicates:\n",
         "acceptance:\n  predicates:\n    - \"%s\"\n" % SYN_UNOWNED)
    row = ("| wc-1 | `%s` | mechanically-verifiable | emitted | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_unowned.py` | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_unowned.positive.json` | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_unowned.negative.json` |") % SYN_UNOWNED
    md_add_mv_projection_row(tree, row)

SYN_EXTRA = "synthetic_registry_owner_entry_absent_from_the_contract_acceptance_predicates"
def mut_extra_owner(tree):
    # add a registry owner entry (owner wc-1) whose predicate is NOT in wc-1's acceptance predicates.
    d = load_registry(tree)
    d["assurance"].append({
        "owner": "wc-1",
        "predicate": SYN_EXTRA,
        "classification": "mechanically-verifiable",
        "ownership": "emitted",
        "checker": ".cdd/unreleased/wc-1/fixtures/synthetic_extra.py",
        "positive_fixture": ".cdd/unreleased/wc-1/fixtures/synthetic_extra.positive.json",
        "negative_fixture": ".cdd/unreleased/wc-1/fixtures/synthetic_extra.negative.json",
        "command": "python3 .cdd/unreleased/wc-1/fixtures/synthetic_extra.py .cdd/unreleased/wc-1/fixtures/synthetic_extra.positive.json",
        "receipt_evidence_predicate": "emits_machine_readable_acceptance_fixtures_into_receipt__contract_template_yaml__worked_instance_yaml__input_union_fixtures__each_validated_by_the_named_acceptance_oracle_command",
    })
    dump_registry(tree, d)
    # keep registry<->projection parity so the ONLY (h) failure is the acceptance-absence
    mv_row = ("| wc-1 | `%s` | mechanically-verifiable | emitted | "
              "`.cdd/unreleased/wc-1/fixtures/synthetic_extra.py` | "
              "`.cdd/unreleased/wc-1/fixtures/synthetic_extra.positive.json` | "
              "`.cdd/unreleased/wc-1/fixtures/synthetic_extra.negative.json` |") % SYN_EXTRA
    md_add_mv_projection_row(tree, mv_row)
    md_add_complete_row(tree, "| wc-1 | `%s` | mechanically-verifiable |" % SYN_EXTRA)

# ---- the four R6 TOTAL-classification mutations (check (j)) -----------------------
def mut_coordinated_omission(tree):
    # THE finding: remove the wc-2 cm_to_v predicate from the registry AND both md surfaces (mv projection
    # + complete projection) AND decrement the mv summary count — but LEAVE it in wc-2's acceptance.
    # (h) and (i) stay green (this is exactly the R5-breaking omission); (j) catches the unclassified
    # obligation. The predicate stays in contracts/wc-2 acceptance untouched.
    d = load_registry(tree)
    d["assurance"] = [e for e in d["assurance"] if e.get("predicate") != CM_TO_V]
    dump_registry(tree, d)
    md_remove_lines_containing(tree, CM_TO_V, expect=2)   # the mv projection row + the complete projection row
    edit(os.path.join(tree, "acceptance-oracles.md"),
         "mechanically-verifiable **30**", "mechanically-verifiable **29**")

def mut_double_classify(tree):
    # a SECOND assurance entry for the wc-1 key_path predicate (a different category) -> not singular.
    d = load_registry(tree)
    d["assurance"].append({
        "owner": "wc-1", "predicate": KEYPATH, "classification": "cognitive-review",
        "independent_review": "synthetic second classification for the same predicate",
    })
    dump_registry(tree, d)

SYN_PHANTOM = "synthetic_phantom_registry_predicate_no_child_contract_declares_it"
def mut_phantom(tree):
    d = load_registry(tree)
    d["assurance"].append({
        "owner": "wc-1", "predicate": SYN_PHANTOM, "classification": "cognitive-review",
        "independent_review": "synthetic phantom entry",
    })
    dump_registry(tree, d)
    md_add_complete_row(tree, "| wc-1 | `%s` | cognitive-review |" % SYN_PHANTOM)  # keep parity; only phantom fires

def mut_category_without_fields(tree):
    # change the wc-1 scope_guard entry from cognitive-review to evidenced WITHOUT its required
    # receipt_evidence field -> (j) rejects the category-without-required-fields.
    d = load_registry(tree)
    e = find_entry(d, "wc-1", SCOPE_GUARD)
    e["classification"] = "evidenced"
    e.pop("receipt_evidence", None)
    dump_registry(tree, d)

# ---- the seven R7 cn.wave.v1 ENVELOPE mutations (check (k)) -----------------------
def wave_path(tree):
    return os.path.join(tree, "wave.cn-wave-v1.yaml")

def mut_wave_schema_nonsense(tree):
    edit(wave_path(tree), 'schema: cn.wave.v1\nwave: "671"', 'schema: nonsense\nwave: "671"')

def mut_wave_dup_schema_key(tree):
    # a SECOND `schema:` mapping key — yaml.safe_load would silently collapse it; StrictSafeLoader rejects.
    edit(wave_path(tree), 'schema: cn.wave.v1\nwave: "671"',
         'schema: cn.wave.v1\nschema: cn.wave.v1\nwave: "671"')

def mut_wave_empty_stop(tree):
    p = wave_path(tree); s = open(p).read()
    start = s.index("stop_conditions:")
    end = s.index("# Non-recursive whole-wave completion")
    open(p, "w").write(s[:start] + "stop_conditions: []\n\n" + s[end:])

def mut_wave_empty_gates(tree):
    p = wave_path(tree); s = open(p).read()
    start = s.index("gates:\n  wave_authorization:")
    end = s.index("# Wave-level STOP conditions")
    open(p, "w").write(s[:start] + "gates: {}\n\n" + s[end:])

def mut_wave_critpath_truncated(tree):
    edit(wave_path(tree), 'critical_path: [ "wc-2", "wc-1", "wc-3b", "wc-5" ]',
         'critical_path: [ "wc-2", "wc-1" ]')

def mut_wave_node_output_path(tree):
    edit(wave_path(tree),
         'output_id: "coherence-measurement-contract"\n    output_path: "docs/architecture/COHERENCE-MEASUREMENT.md"',
         'output_id: "coherence-measurement-contract"\n    output_path: "docs/WRONG.md"')

def mut_wave_edge_kind_nonsense(tree):
    edit(wave_path(tree), '{ from: "wc-2",  to: "wc-1",  kind: "depends_on" }',
         '{ from: "wc-2",  to: "wc-1",  kind: "nonsense" }')

# ---- the seven R7 EVIDENCE-DERIVED-completion mutations (check (g)) ---------------
D_WC1 = ("{ requested_output_produced: true, acceptance_all_pass: true, class_v_pass: true, "
         "receipt_bound: true, evidence_bound: true }  # all-complete wc-1")
D_WC5 = ("{ requested_output_produced: true, acceptance_all_pass: true, class_v_pass: true, "
         "receipt_bound: true, evidence_bound: true }  # all-complete wc-5")

def mut_cc_definition_constant_true(tree):
    # the exact Codex mutation: child_complete becomes constant `true`, ignoring output/acceptance/V/receipt.
    edit(wave_path(tree),
         ('    definition:\n      all:\n'
          '        - { record: "requested_output_produced" }\n'
          '        - { record: "acceptance_all_pass" }\n'
          '        - { record: "class_v_pass" }\n'
          '        - { record: "receipt_bound" }\n'
          '        - { record: "evidence_bound" }   # content-bound β/γ evidence present where required'),
         '    definition:\n      all:\n        - { lit: true }')

def _flip_D_wc1(tree, field):
    edit(wave_path(tree), D_WC1, D_WC1.replace(f"{field}: true", f"{field}: false", 1))

def mut_flip_output(tree):     _flip_D_wc1(tree, "requested_output_produced")
def mut_flip_acceptance(tree): _flip_D_wc1(tree, "acceptance_all_pass")
def mut_flip_classv(tree):     _flip_D_wc1(tree, "class_v_pass")
def mut_flip_receipt(tree):    _flip_D_wc1(tree, "receipt_bound")

def mut_remove_constituent(tree):
    # drop the requested_output_produced constituent from wc-1's all-complete record (still valid YAML).
    edit(wave_path(tree), D_WC1, D_WC1.replace("requested_output_produced: true, ", "", 1))

def mut_wc5_completion_naked(tree):
    # flip wc-5's receipt_bound false while its claimed child_complete stays true -> completion (not just
    # readiness) fails: the naked boolean is not backed by wc-5's own constituents.
    edit(wave_path(tree), D_WC5, D_WC5.replace("receipt_bound: true", "receipt_bound: false", 1))

CASES = [
    ("wrong-intent-id",       mut_wrong_intent_id,        "b", "intent_ref.id"),
    ("nonsense-class",        mut_nonsense_class,          "a", "cell.class not in"),
    ("nonexistent-artifact",  mut_nonexistent_path,        "b", "DOES-NOT-EXIST.md"),
    ("tautological-wave",     mut_tautology,               "g", "CYCLE"),
    ("flipped-fixture",       mut_flipped_expected,        "g", "authored expected"),
    ("placeholder-oracle",    mut_placeholder_oracle,      "h", "placeholder operand"),
    ("remove-registry-entry", mut_remove_entry,            "h", "has NO registry owner entry"),
    ("duplicate-owner",       mut_duplicate_owner,         "h", "DUPLICATE owner"),
    ("reclassify-entry",      mut_reclassify,              "j", "classification mismatch"),
    ("unowned-mv-predicate",  mut_unowned_predicate,       "h", "has NO registry owner entry"),
    ("extra-owner-entry",     mut_extra_owner,             "h", "extra owner entry"),
    ("coordinated-omission",  mut_coordinated_omission,    "j", "absent from the assurance registry"),
    ("double-classify",       mut_double_classify,         "j", "DOUBLE-CLASSIFIED"),
    ("phantom-entry",         mut_phantom,                 "j", "PHANTOM"),
    ("category-no-fields",    mut_category_without_fields, "j", "missing required field"),
    # --- R7 cn.wave.v1 envelope (check (k)) ---
    ("wave-schema-nonsense",  mut_wave_schema_nonsense,    "k", "wave.schema != cn.wave.v1"),
    ("wave-empty-stop",       mut_wave_empty_stop,         "k", "STOP conditions must be a NON-EMPTY"),
    ("wave-empty-gates",      mut_wave_empty_gates,        "k", "wave.gates: MISSING"),
    ("wave-critpath-trunc",   mut_wave_critpath_truncated, "k", "does not terminate at the terminal node"),
    ("wave-node-outputpath",  mut_wave_node_output_path,   "k", "docs/WRONG.md"),
    ("wave-edge-kind",        mut_wave_edge_kind_nonsense, "k", "kind not in ['depends_on']"),
    ("wave-dup-schema-key",   mut_wave_dup_schema_key,     "k", "DUPLICATE mapping key"),
    # --- R7 evidence-derived completion (check (g)) ---
    ("cc-def-constant-true",  mut_cc_definition_constant_true, "g", "IGNORES required constituent"),
    ("cc-flip-output",        mut_flip_output,             "g", "naked boolean"),
    ("cc-flip-acceptance",    mut_flip_acceptance,         "g", "naked boolean"),
    ("cc-flip-classv",        mut_flip_classv,             "g", "naked boolean"),
    ("cc-flip-receipt",       mut_flip_receipt,            "g", "naked boolean"),
    ("cc-remove-constituent", mut_remove_constituent,      "g", "missing constituent"),
    ("cc-wc5-naked",          mut_wc5_completion_naked,    "g", "naked boolean"),
]

def main():
    failures = []

    # clean (unmodified copy) must PASS
    tmp, tree = copy_tree()
    try:
        rc, out = run_validator(tree)
        if rc != 0 or "RESULT: PASS" not in out:
            failures.append(f"[clean] expected exit 0 / PASS, got exit {rc}\n{out}")
        print(f"[clean-tree]              exit={rc}  {'OK (PASS)' if rc == 0 else 'UNEXPECTED'}")
    finally:
        shutil.rmtree(tmp, ignore_errors=True)

    # each mutation must FAIL for its own named check
    for name, mut, check, substr in CASES:
        tmp, tree = copy_tree()
        try:
            mut(tree)
            resync_all(tree)          # honestly re-pin any changed contract_sha256 / registry hash
            rc, out = run_validator(tree)
            own_check_failed = f"({check}) FAIL" in out
            substr_present = substr in out
            ok = (rc != 0) and own_check_failed and substr_present
            if not ok:
                failures.append(f"[{name}] exit={rc} own_check({check})_failed={own_check_failed} "
                                f"substr({substr!r})={substr_present}\n{out}")
            print(f"[{name:22}] exit={rc}  check({check}) FAIL={own_check_failed}  "
                  f"substr={substr_present}  {'OK' if ok else 'UNEXPECTED'}")
        finally:
            shutil.rmtree(tmp, ignore_errors=True)

    print("-" * 70)
    if failures:
        print("ADVERSARIAL-MUTATION HARNESS: FAIL")
        for x in failures:
            print("  -", x)
        sys.exit(1)
    print("ADVERSARIAL-MUTATION HARNESS: PASS — clean tree exits 0; all 29 adversarial "
          "mutations exit non-zero for their own named predicate.")
    sys.exit(0)

if __name__ == "__main__":
    main()
