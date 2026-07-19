#!/usr/bin/env python3
"""
Executable negative-fixture harness for validate.py (wave cnos#671 R3).

Materializes the SIX adversarial mutations (the FIVE the external β used to false-pass the R2 validator,
plus the R4 oracle-ownership fixture), each in its own temp copy of the wave tree (honestly updating
node.contract_sha256 when a contract is mutated, exactly as β did), runs validate.py on each, and asserts
EACH exits non-zero FOR ITS OWN NAMED PREDICATE — while the UNMODIFIED tree exits 0.

  1. wrong intent id           -> check (b): intent_ref.id != the actual intent object id
  2. nonsense cell.class       -> check (a): enum cell.class ∈ {working,planning,cohering}
  3. nonexistent artifact path -> check (b): git cannot resolve <commit>:<path> (DOES-NOT-EXIST.md)
  4. tautological whole-wave   -> check (g): predicate dependency graph has a CYCLE
  5. flipped fixture expected  -> check (g): fixture computed result != authored `expected`
  6. placeholder oracle operand-> check (h): a mechanically-verifiable oracle's checker left as a
                                  placeholder (unowned/non-concrete) → oracle ownership fails

Run:  python3 validate_test.py         (exit 0 iff all 6 fail for their own predicate AND clean passes)
"""
import os, sys, re, shutil, tempfile, subprocess, hashlib

WAVE_DIR = os.path.dirname(os.path.abspath(__file__))

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

def repin_contract(tree_dir, nid, old_sha):
    """Honestly update node.contract_sha256 in wave.yaml after mutating a contract (as β did)."""
    new_sha = sha256_file(os.path.join(tree_dir, "contracts", f"{nid}.cn-cell-contract-v1.yaml"))
    edit(os.path.join(tree_dir, "wave.cn-wave-v1.yaml"), old_sha, new_sha, count=1)

# ---- the five mutations ----------------------------------------------------------
def mut_wrong_intent_id(tree):
    f = os.path.join(tree, "contracts", "wc-1.cn-cell-contract-v1.yaml")
    old = sha256_file(f)
    edit(f, 'id: "intent-2026-0719-671-cell-runtime-doctrine"',
            'id: "intent-WRONG-does-not-match-the-object"')
    repin_contract(tree, "wc-1", old)

def mut_nonsense_class(tree):
    f = os.path.join(tree, "contracts", "wc-1.cn-cell-contract-v1.yaml")
    old = sha256_file(f)
    edit(f, 'class: "working"', 'class: "nonsense"')
    repin_contract(tree, "wc-1", old)

def mut_nonexistent_path(tree):
    f = os.path.join(tree, "contracts", "wc-1.cn-cell-contract-v1.yaml")
    old = sha256_file(f)
    edit(f, 'path: "docs/architecture/CELL-RUNTIME.md"',
            'path: "docs/architecture/DOES-NOT-EXIST.md"')
    repin_contract(tree, "wc-1", old)

def mut_tautology(tree):
    f = os.path.join(tree, "wave.cn-wave-v1.yaml")
    edit(f,
         '      expr:\n        and:\n          - { all_over_N: "child_complete" }\n          - { pred: "wc5_complete" }',
         '      expr: { pred: "wave_complete" }')   # wave_complete := wave_complete (tautology)

def mut_flipped_expected(tree):
    f = os.path.join(tree, "wave.cn-wave-v1.yaml")
    edit(f, 'expected: { wc5_ready: true, wc5_complete: true, wave_complete: true }',
            'expected: { wc5_ready: true, wc5_complete: true, wave_complete: false }')

def mut_placeholder_oracle(tree):
    # Leave a mechanically-verifiable oracle's checker as a PLACEHOLDER operand (unowned / non-concrete):
    # check (h) must reject it (placeholder operand + not in allowed_paths + command mismatch).
    f = os.path.join(tree, "contracts", "wc-2.cn-cell-contract-v1.yaml")
    old = sha256_file(f)
    edit(f, 'checker: ".cdd/unreleased/wc-2/fixtures/check_v_signature.py"',
            'checker: "<checker-to-be-decided>"')
    repin_contract(tree, "wc-2", old)

CASES = [
    ("wrong-intent-id",       mut_wrong_intent_id,   "b", "intent_ref.id"),
    ("nonsense-class",        mut_nonsense_class,    "a", "cell.class not in"),
    ("nonexistent-artifact",  mut_nonexistent_path,  "b", "DOES-NOT-EXIST.md"),
    ("tautological-wave",     mut_tautology,         "g", "CYCLE"),
    ("flipped-fixture",       mut_flipped_expected,  "g", "authored expected"),
    ("placeholder-oracle",    mut_placeholder_oracle, "h", "placeholder operand"),
]

def main():
    failures = []

    # clean (unmodified copy) must PASS
    tmp, tree = copy_tree()
    try:
        rc, out = run_validator(tree)
        if rc != 0 or "RESULT: PASS" not in out:
            failures.append(f"[clean] expected exit 0 / PASS, got exit {rc}")
        print(f"[clean-tree]            exit={rc}  {'OK (PASS)' if rc == 0 else 'UNEXPECTED'}")
    finally:
        shutil.rmtree(tmp, ignore_errors=True)

    # each mutation must FAIL for its own named check
    for name, mut, check, substr in CASES:
        tmp, tree = copy_tree()
        try:
            mut(tree)
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
        print("NEGATIVE-FIXTURE HARNESS: FAIL")
        for x in failures:
            print("  -", x)
        sys.exit(1)
    print("NEGATIVE-FIXTURE HARNESS: PASS — clean tree exits 0; all 6 adversarial "
          "mutations exit non-zero for their own named predicate.")
    sys.exit(0)

if __name__ == "__main__":
    main()
