#!/usr/bin/env python3
# wave-revision: R5
"""
Executable adversarial-mutation harness for validate.py (wave cnos#671 R5).

Materializes ELEVEN adversarial mutations, each in its own temp copy of the wave tree, HONESTLY
re-pinning every changed contract_sha256 and the oracle-registry content hash (via resync_all), runs
validate.py on each, and asserts EACH exits non-zero FOR ITS OWN NAMED PREDICATE — while the UNMODIFIED
tree exits 0.

The SIX original negatives (β used the first five to false-pass the R2 validator; the sixth is the
oracle-ownership negative, now against the registry):
  1. wrong intent id            -> (b): intent_ref.id != the actual intent object id
  2. nonsense cell.class        -> (a): enum cell.class ∈ {working,planning,cohering}
  3. nonexistent artifact path  -> (b): git cannot resolve <commit>:<path> (DOES-NOT-EXIST.md)
  4. tautological whole-wave    -> (g): predicate dependency graph has a CYCLE
  5. flipped fixture expected   -> (g): fixture computed result != authored `expected`
  6. placeholder oracle operand -> (h): a registry entry's checker left as a placeholder (non-concrete)

The FIVE R5 bijection mutations (Fix 2):
  7. remove one registry entry (acceptance predicate + md projection row left intact) -> (h): a
     mechanically-verifiable predicate in the projection with NO registry owner entry (MISSING). This is
     the exact class of mutation the external β used to break R4 — it MUST now fail.
  8. duplicate one entry under another owner                                          -> (h): DUPLICATE owner
  9. reclassify one entry to evidenced/cognitive                                      -> (h): classification mismatch
 10. leave one mechanically-verifiable acceptance predicate unowned (in a contract +
     md projection, absent from the registry)                                        -> (h): MISSING entry
 11. add a registry owner entry absent from the contract's acceptance predicates      -> (h): extra owner entry

Run:  python3 validate_test.py     (exit 0 iff all 11 fail for their own predicate AND clean passes)
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
    with open(os.path.join(tree, "oracle-registry.yaml"), "w") as f:
        yaml.safe_dump(data, f, sort_keys=False, width=100000)

def md_add_projection_row(tree, row):
    md = os.path.join(tree, "acceptance-oracles.md")
    edit(md, "\n\n---\n\n*Rows marked", "\n" + row + "\n\n---\n\n*Rows marked")

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
    # Leave a registry entry's checker as a PLACEHOLDER operand (non-concrete): check (h) rejects it.
    d = load_registry(tree)
    for e in d["oracles"]:
        if e.get("checker", "").endswith("check_v_signature.py"):
            e["checker"] = "<checker-to-be-decided>"
            break
    dump_registry(tree, d)

# ---- the five R5 bijection mutations ---------------------------------------------
def mut_remove_entry(tree):
    # remove wc-2 cm_to_v_edge_pinned entry from the registry; leave its acceptance predicate + the md
    # projection row intact -> a projection predicate with no registry owner (MISSING entry).
    d = load_registry(tree)
    d["oracles"] = [e for e in d["oracles"]
                    if not e["predicate"].startswith("cm_to_v_edge_pinned__")]
    dump_registry(tree, d)

def mut_duplicate_owner(tree):
    d = load_registry(tree)
    src = next(e for e in d["oracles"]
               if e["predicate"] == "key_path_parity__template_and_worked_instance_normalize_to_identical_key_path_sets")
    dup = copy.deepcopy(src)
    dup["owner"] = "wc-2"     # same predicate now owned by a second child -> mapping not singular
    d["oracles"].append(dup)
    dump_registry(tree, d)

def mut_reclassify(tree):
    d = load_registry(tree)
    d["oracles"][0]["classification"] = "evidenced"   # registry vs mechanically-verifiable mismatch
    dump_registry(tree, d)

SYN_UNOWNED = "synthetic_unowned_mechanically_verifiable_predicate_present_in_contract_absent_from_registry"
def mut_unowned_predicate(tree):
    # add an mv predicate to wc-1 acceptance + the md projection, but NOT to the registry.
    edit(contract_path(tree, "wc-1"), "acceptance:\n  predicates:\n",
         "acceptance:\n  predicates:\n    - \"%s\"\n" % SYN_UNOWNED)
    row = ("| wc-1 | `%s` | mechanically-verifiable | emitted | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_unowned.py` | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_unowned.positive.json` | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_unowned.negative.json` |") % SYN_UNOWNED
    md_add_projection_row(tree, row)

SYN_EXTRA = "synthetic_registry_owner_entry_absent_from_the_contract_acceptance_predicates"
def mut_extra_owner(tree):
    # add a registry owner entry (owner wc-1) whose predicate is NOT in wc-1's acceptance predicates.
    d = load_registry(tree)
    d["oracles"].append({
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
    row = ("| wc-1 | `%s` | mechanically-verifiable | emitted | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_extra.py` | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_extra.positive.json` | "
           "`.cdd/unreleased/wc-1/fixtures/synthetic_extra.negative.json` |") % SYN_EXTRA
    md_add_projection_row(tree, row)

CASES = [
    ("wrong-intent-id",       mut_wrong_intent_id,    "b", "intent_ref.id"),
    ("nonsense-class",        mut_nonsense_class,     "a", "cell.class not in"),
    ("nonexistent-artifact",  mut_nonexistent_path,   "b", "DOES-NOT-EXIST.md"),
    ("tautological-wave",     mut_tautology,          "g", "CYCLE"),
    ("flipped-fixture",       mut_flipped_expected,   "g", "authored expected"),
    ("placeholder-oracle",    mut_placeholder_oracle, "h", "placeholder operand"),
    ("remove-registry-entry", mut_remove_entry,       "h", "has NO registry owner entry"),
    ("duplicate-owner",       mut_duplicate_owner,    "h", "DUPLICATE owner"),
    ("reclassify-entry",      mut_reclassify,         "h", "classification"),
    ("unowned-mv-predicate",  mut_unowned_predicate,  "h", "has NO registry owner entry"),
    ("extra-owner-entry",     mut_extra_owner,        "h", "extra owner entry"),
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
                                f"substr({substr!r})={substr_present}")
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
    print("ADVERSARIAL-MUTATION HARNESS: PASS — clean tree exits 0; all 11 adversarial "
          "mutations exit non-zero for their own named predicate.")
    sys.exit(0)

if __name__ == "__main__":
    main()
