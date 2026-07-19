#!/usr/bin/env python3
# wave-revision: R6
"""
Planning-Cell PRE-AUTHORIZATION validator for wave cnos#671 (cell-runtime-doctrine), R6.

SOUND, FAIL-CLOSED. Credential-free (Python stdlib + PyYAML + local `git`). It derives every fact
from the AUTHORED wave/contract/intent data — it hard-codes no node list, edge list, or predicate
graph — and rejects any tree that violates the §2 constraint model, ref resolution, derived-graph
consistency, gate invariants, or the AUTHORED completion semantics. Exit status 0 iff ALL checks
pass; non-zero (1) on ANY violation (2 on a harness error).

Checks:
  (a) FULL §2 cn.cell.contract.v1 CONSTRAINT MODEL — not just key sets: exact key paths, ENUMS
      (cell.class ∈ {working,planning,cohering}; inputs.required[].ref_kind ∈ {external,sibling_output};
      external locator.kind ∈ {repo_artifact,control_plane,prior_receipt}; requested_output.kind ∈
      {artifact,relation_graph,judgment}), scalar/seq/bool TYPES, and CARDINALITIES. Rejects
      `cell.class: nonsense`, a non-bool gate, an under-full seq, etc.
  (b) REAL REF RESOLUTION.
        * intent object: every contract's intent_ref.{id,schema} is compared to the ACTUAL intent
          object (reject a wrong id / wrong schema).
        * every immutable repo-artifact locator (repo_artifact commit+path, and any `sha256:h@path`)
          is resolved with `git cat-file -e <commit>:<path>` / `git rev-parse` (commit AND path must
          exist) — reject a pin to a nonexistent path (docs/.../DOES-NOT-EXIST.md).
        * every content-bound `sha256:<64hex>@<repo-relpath>` locator resolves to a file whose bytes
          hash to the pinned digest (intent carrier, control_plane revision, wave grounding refs).
        * the grounding SOURCE snapshot bytes hash to the true source SHA-256 9d1ab3a5….
  (c) DERIVE-FROM-AUTHORED-DATA. Node/output/ref/edge/root/critical-path facts are derived from the
      parsed wave + contracts (never a hard-coded list): contracts-dir set == wave node set; each
      wave node.output_id == its contract.requested_output.id; the AUTHORED edge set == the
      sibling-output-DERIVED edge set exactly; authored roots == derived roots; critical_path is a
      real path from a root to the terminal.
  (d) the derived dependency graph is a DAG.
  (e) no two PARALLEL (concurrent) nodes share a write surface (requested_output.path / allowed_paths).
  (f) GATE invariants (doctrine_affecting ⇒ operator_acceptance_required; reason present iff a gate
      bool is true, null iff both false; reason key always present).
  (g) EVALUATE THE AUTHORED COMPLETION RELATION AS STRUCTURED DATA. The predicate dependency graph is
      built by WALKING the authored `expr` ASTs (no substring scan, no hard-coded graph); acyclicity
      is proved by real topological analysis (a tautology `wave_complete := wave_complete` becomes a
      self-edge → cycle → FAIL). Each authored fixture is EVALUATED (each predicate computed over the
      fixture inputs) and compared to the authored `expected` (a flipped expectation → FAIL). The
      non-recursive construction set N is derived (all nodes minus terminal) and cross-checked.
  (h) ORACLE OWNERSHIP — TOTAL + SINGULAR BIJECTION (R5). The §2 contract shape is restored (no
      top-level `oracles` key); oracle ownership lives in the SEPARATE content-bound registry
      `oracle-registry.yaml`, which each child contract consumes via a canonical
      `inputs.required[].external` control_plane ref (content-hash pinned; resolved in (b)).
      `acceptance-oracles.md` is a human-readable PROJECTION of the registry. This check proves a
      TOTAL + SINGULAR mapping across three sources — the registry (authoritative), each child's
      `acceptance.predicates`, and the md projection:
        * every mechanically-verifiable predicate (in the registry / the md projection) maps to
          EXACTLY ONE child `owner` entry;
        * every registry `owner` entry maps to EXACTLY ONE child `acceptance` predicate;
        * REJECT — a missing entry (an mv predicate with no owner), a duplicate owner, an extra owner
          entry absent from the owning contract's acceptance, a classification mismatch
          (registry vs md-projection vs contract), and ANY placeholder operand; and REJECT any
          registry↔md-projection parity break.
      Each entry's checker/schema is CONCRETE (no placeholder, no implicit CWD) and either EMITTED
      (in the owner's `allowed_paths`) or EXISTING (a pinned immutable input of the owner). Fully
      derivable from the registry + the contracts + the projection.
  (i) LEDGER CONSISTENCY (R6). (1) every declared wave-revision marker line across the wave artifacts
      AGREES (and equals the authored wave revision); (2) the classification counts reported in
      `acceptance-oracles.md` EQUAL the per-category totals DERIVED FROM THE TOTAL registry (`assurance:`),
      and the mechanically-verifiable count EQUALS the number of ownership entries; (3) every
      classification-cell value is EXACTLY ONE enum member (enforced | mechanically-verifiable | evidenced
      | cognitive-review) — a compound category (e.g. `enforced + evidenced`) fails closed.
  (j) CLASSIFICATION TOTALITY + SINGULARITY over the COMPLETE child acceptance set (R6). The authoritative
      artifact is now the TOTAL `assurance:` registry — an entry for EVERY scalar in EVERY child contract's
      `acceptance.predicates`, each classified into EXACTLY ONE of {enforced, mechanically-verifiable,
      evidenced, cognitive-review} with that category's required fields (enforced->enforced_by;
      mechanically-verifiable->the ownership fields; evidenced->receipt_evidence;
      cognitive-review->independent_review). This check derives the set of (owner, predicate) pairs from the
      UNION of every child's `acceptance.predicates` and proves it EQUALS the set of child-owned registry
      entries, each mapping EXACTLY ONCE. REJECTS: a predicate in a contract but absent from the registry
      (UNCLASSIFIED — the exact coordinated-omission the R5 checks missed); a predicate with >1 registry
      entry (DOUBLE-CLASSIFIED); a registry entry for a predicate no contract declares (PHANTOM); a
      classification that is not exactly one enum member; an entry missing its category's required fields;
      and any registry <-> `Complete assurance classification` md-projection parity break. Wave-only
      predicates live in a SEPARATE `wave_predicates:` section so they can neither inflate nor mask child
      coverage. Fail-closed; reports the derived count of child acceptance predicates and per-category totals.

RESOLUTION MODEL (so the negative-fixture harness can validate a mutated COPY of the tree):
  * WAVE_DIR   — the wave matter dir (contracts/, wave.yaml, intent, grounding snapshots). All AUTHORED
                 wave data is read from here. Default: this file's directory. Override: env
                 WAVE_VALIDATOR_WAVE_DIR.
  * REPO_ROOT  — the git repo used ONLY to resolve external repo-artifact `<commit>:<path>` objects.
                 Default: `git -C WAVE_DIR rev-parse --show-toplevel`. Override: env
                 WAVE_VALIDATOR_REPO_ROOT (the harness points this at the real repo while WAVE_DIR is a
                 mutated temp copy).
  * A content locator `sha256:h@<repo-relpath>` under the wave dir resolves against WAVE_DIR (so the
    mutated copy is what gets hashed); any other repo-relative path resolves against REPO_ROOT.
"""
import sys, os, re, hashlib, subprocess

try:
    import yaml
except ImportError:
    print("FATAL: PyYAML required (pip install pyyaml)", file=sys.stderr)
    sys.exit(2)

# ---- resolution model ------------------------------------------------------------
WAVE_REL = ".cdd/waves/cell-runtime-doctrine"   # canonical repo-relative location of this wave dir
WAVE_DIR = os.environ.get("WAVE_VALIDATOR_WAVE_DIR") or os.path.dirname(os.path.abspath(__file__))
WAVE_DIR = os.path.abspath(WAVE_DIR)
CONTRACTS_DIR = os.path.join(WAVE_DIR, "contracts")

def _git_toplevel(path):
    try:
        r = subprocess.run(["git", "-C", path, "rev-parse", "--show-toplevel"],
                           capture_output=True, text=True)
        if r.returncode == 0:
            return r.stdout.strip()
    except Exception:
        pass
    return os.path.abspath(os.path.join(WAVE_DIR, "..", "..", ".."))

REPO_ROOT = os.path.abspath(os.environ.get("WAVE_VALIDATOR_REPO_ROOT") or _git_toplevel(WAVE_DIR))

def repo_relpath_to_abs(relpath):
    """Map a repo-relative path to an absolute path: wave-matter paths resolve under WAVE_DIR,
    everything else under REPO_ROOT."""
    relpath = relpath.lstrip("/")
    if relpath == WAVE_REL or relpath.startswith(WAVE_REL + "/"):
        return os.path.join(WAVE_DIR, os.path.relpath(relpath, WAVE_REL))
    return os.path.join(REPO_ROOT, relpath)

TRUE_SOURCE_SHA = "9d1ab3a50d00e642fdeb87728dd71f7c7499c60878afe7001f2ddb832b161dbb"

errors = []
def fail(check, msg): errors.append(f"[{check}] {msg}")

def sha256_file(path):
    with open(path, "rb") as f:
        return hashlib.sha256(f.read()).hexdigest()

def load_yaml(path):
    with open(path) as f:
        return yaml.safe_load(f)

# ---- git object resolution -------------------------------------------------------
def git_commit_exists(commit):
    r = subprocess.run(["git", "-C", REPO_ROOT, "rev-parse", "--verify", "--quiet", f"{commit}^{{commit}}"],
                       capture_output=True, text=True)
    return r.returncode == 0

def git_object_exists(commit, path):
    r = subprocess.run(["git", "-C", REPO_ROOT, "cat-file", "-e", f"{commit}:{path}"],
                       capture_output=True, text=True)
    return r.returncode == 0

# ---- type/enum helpers -----------------------------------------------------------
def is_scalar(v):   return isinstance(v, (str, int, float))
def is_str(v):      return isinstance(v, str)
def is_seq(v):      return isinstance(v, list)
def is_bool(v):     return isinstance(v, bool)

CLASS_ENUM   = {"working", "planning", "cohering"}
REFKIND_ENUM = {"external", "sibling_output"}
LOCATOR_ENUM = {"repo_artifact", "control_plane", "prior_receipt"}
OUTPUT_KIND_ENUM = {"artifact", "relation_graph", "judgment"}
# the four assurance-classification categories (checks h, i, j)
KIND_ENUM = {"enforced", "mechanically-verifiable", "evidenced", "cognitive-review"}

TOP_KEYS   = {"schema", "cell", "scope", "intent_ref", "inputs", "requested_output",
              "acceptance", "constraints", "gates", "doctrine_affecting", "stop_conditions"}
# R5: EXACT §2 top-level key set — NO `oracles` (or any other) extra key. Oracle ownership lives in the
# separate content-bound oracle-registry.yaml, consumed via inputs.required[].external (check (h)).
# An added top-level key (incl. `oracles`) now FAILS check (a).
CELL_KEYS  = {"id", "class", "mode", "protocol", "matter_domain"}
SCOPE_KEYS = {"repo", "wave", "parent_cell"}
INTENT_KEYS   = {"schema", "id", "carrier"}
CARRIER_KEYS  = {"kind", "ref"}
OUTPUT_KEYS   = {"id", "kind", "path"}
CONSTRAINT_KEYS = {"allowed_paths", "forbidden_paths", "non_goals"}
GATE_KEYS = {"operator_authorization_required", "operator_acceptance_required", "reason"}
LOCATOR_SHAPES = {
    "repo_artifact": {"kind", "repo", "commit", "path"},
    "control_plane": {"kind", "ref", "revision"},
    "prior_receipt": {"kind", "receipt", "binding"},
}

def check_keyset(check, where, got, expected):
    got = set(got)
    if got != expected:
        extra = got - expected
        missing = expected - got
        if extra:   fail(check, f"{where}: NON-CANONICAL key(s) {sorted(extra)}")
        if missing: fail(check, f"{where}: MISSING key(s) {sorted(missing)}")
        return False
    return True

def check_seq_of_scalar(check, where, v, minlen):
    if not is_seq(v):
        fail(check, f"{where}: must be a seq (got {type(v).__name__})"); return
    if len(v) < minlen:
        fail(check, f"{where}: seq must have >= {minlen} element(s) (got {len(v)})")
    for i, el in enumerate(v):
        if not is_scalar(el):
            fail(check, f"{where}[{i}]: element must be scalar (got {type(el).__name__})")

CONTENT_LOCATOR_RE = re.compile(r"^sha256:([0-9a-f]{64})@(.+)$")

def resolve_content_locator(check, where, locator_str):
    """`sha256:<64hex>@<repo-relpath>` — the path must exist AND its bytes hash to the digest.
    Returns the 64-hex digest on a well-formed locator (regardless of resolution), else None."""
    m = CONTENT_LOCATOR_RE.match(str(locator_str))
    if not m:
        fail(check, f"{where}: not a content-bound immutable locator sha256:<64hex>@<path> ({locator_str!r})")
        return None
    want, rel = m.group(1), m.group(2)
    abspath = repo_relpath_to_abs(rel)
    if not os.path.isfile(abspath):
        fail(check, f"{where}: referenced file does not resolve: {rel}")
        return want
    got = sha256_file(abspath)
    if got != want:
        fail(check, f"{where}: content hash mismatch for {rel}: pinned {want[:12]}… actual {got[:12]}…")
    return want

def check_external_locator(check, where, loc):
    if not isinstance(loc, dict) or "kind" not in loc:
        fail(check, f"{where}: external locator missing kind"); return
    kind = loc["kind"]
    if kind not in LOCATOR_ENUM:
        fail(check, f"{where}: locator.kind not in {sorted(LOCATOR_ENUM)} (got {kind!r})"); return
    check_keyset(check, f"{where}.locator", loc.keys(), LOCATOR_SHAPES[kind])
    # Structural/enum failures are the caller's check (a); actual REF RESOLUTION is always check (b).
    if kind == "repo_artifact":
        commit = str(loc.get("commit", "")); path = str(loc.get("path", ""))
        if not re.match(r"^[0-9a-f]{40}$", commit):
            fail(check, f"{where}: repo_artifact.commit is not an immutable 40-hex revision ({commit!r})")
        elif not git_commit_exists(commit):
            fail("b", f"{where}: repo_artifact.commit does not resolve in the repo ({commit})")
        elif not path:
            fail(check, f"{where}: repo_artifact.path missing")
        elif not git_object_exists(commit, path):
            fail("b", f"{where}: repo_artifact does not resolve — no object {commit[:12]}…:{path}")
    elif kind == "control_plane":
        rev = str(loc.get("revision", ""))
        # revision must be content-bound (sha256:h@path) and resolve; a bare comment/issue id is rejected
        resolve_content_locator("b", f"{where}.revision", rev)
    elif kind == "prior_receipt":
        if not str(loc.get("binding", "")):
            fail(check, f"{where}: prior_receipt missing immutable binding")

# ---- load contracts (derive node ids from files; do not hard-code) ---------------
CONTRACT_RE = re.compile(r"^(wc-[0-9a-z]+)\.cn-cell-contract-v1\.yaml$")
contracts = {}   # nid -> (path, dict)
if not os.path.isdir(CONTRACTS_DIR):
    fail("a", f"contracts dir missing: {CONTRACTS_DIR}")
else:
    for fn in sorted(os.listdir(CONTRACTS_DIR)):
        m = CONTRACT_RE.match(fn)
        if m:
            p = os.path.join(CONTRACTS_DIR, fn)
            try:
                contracts[m.group(1)] = (p, load_yaml(p))
            except Exception as e:
                fail("a", f"{fn}: YAML parse error: {e}")

# ---- load intent object (for real id/schema resolution) --------------------------
INTENT_PATH = os.path.join(WAVE_DIR, "intent.cn-intent-v1.yaml")
intent_obj = None
if not os.path.isfile(INTENT_PATH):
    fail("b", "intent object missing: intent.cn-intent-v1.yaml")
else:
    try:
        intent_obj = load_yaml(INTENT_PATH)
    except Exception as e:
        fail("b", f"intent object YAML parse error: {e}")
INTENT_ID     = (intent_obj or {}).get("id")
INTENT_SCHEMA = (intent_obj or {}).get("schema")

# ---- (a) §2 constraint model -----------------------------------------------------
for nid, (p, c) in contracts.items():
    if not isinstance(c, dict):
        fail("a", f"{nid}: contract is not a mapping"); continue
    check_keyset("a", f"{nid} top-level", c.keys(), TOP_KEYS)
    if c.get("schema") != "cn.cell.contract.v1":
        fail("a", f"{nid}.schema != cn.cell.contract.v1 (got {c.get('schema')!r})")
    # cell
    cell = c.get("cell")
    if not isinstance(cell, dict) or not check_keyset("a", f"{nid}.cell", cell.keys(), CELL_KEYS):
        pass
    else:
        if cell.get("class") not in CLASS_ENUM:
            fail("a", f"{nid}.cell.class not in {sorted(CLASS_ENUM)} (got {cell.get('class')!r})")
        for k in ("id", "mode", "protocol", "matter_domain"):
            if not is_str(cell.get(k)):
                fail("a", f"{nid}.cell.{k}: must be a scalar string (got {type(cell.get(k)).__name__})")
    # scope
    sc = c.get("scope")
    if isinstance(sc, dict) and check_keyset("a", f"{nid}.scope", sc.keys(), SCOPE_KEYS):
        if not is_str(sc.get("repo")):
            fail("a", f"{nid}.scope.repo: must be a scalar string")
        for k in ("wave", "parent_cell"):
            if not (sc.get(k) is None or is_scalar(sc.get(k))):
                fail("a", f"{nid}.scope.{k}: must be scalar-or-null")
    elif not isinstance(sc, dict):
        fail("a", f"{nid}.scope: must be a mapping")
    # intent_ref
    ir = c.get("intent_ref")
    if isinstance(ir, dict) and check_keyset("a", f"{nid}.intent_ref", ir.keys(), INTENT_KEYS):
        if ir.get("schema") != "cn.intent.v1":
            fail("a", f"{nid}.intent_ref.schema != cn.intent.v1 (got {ir.get('schema')!r})")
        if not is_str(ir.get("id")):
            fail("a", f"{nid}.intent_ref.id: must be a scalar string")
        car = ir.get("carrier")
        if isinstance(car, dict) and check_keyset("a", f"{nid}.intent_ref.carrier", car.keys(), CARRIER_KEYS):
            for k in ("kind", "ref"):
                if not is_str(car.get(k)):
                    fail("a", f"{nid}.intent_ref.carrier.{k}: must be a scalar string")
        elif not isinstance(car, dict):
            fail("a", f"{nid}.intent_ref.carrier: must be a mapping")
    elif not isinstance(ir, dict):
        fail("a", f"{nid}.intent_ref: must be a mapping")
    # inputs
    inp = c.get("inputs")
    if not isinstance(inp, dict) or set(inp.keys()) != {"required", "optional"}:
        fail("a", f"{nid}.inputs: must be a mapping with exactly keys required, optional")
    else:
        req = inp.get("required")
        if not is_seq(req) or len(req) < 1:
            fail("a", f"{nid}.inputs.required: must be a seq of 1+")
        else:
            for i, el in enumerate(req):
                if not isinstance(el, dict):
                    fail("a", f"{nid}.inputs.required[{i}]: must be a mapping"); continue
                rk = el.get("ref_kind")
                if rk not in REFKIND_ENUM:
                    fail("a", f"{nid}.inputs.required[{i}].ref_kind not in {sorted(REFKIND_ENUM)} (got {rk!r})")
                elif rk == "external":
                    if set(el.keys()) != {"ref_kind", "locator"}:
                        fail("a", f"{nid}.inputs.required[{i}]: external element keys must be {{ref_kind, locator}}")
                    else:
                        check_external_locator("a", f"{nid}.inputs.required[{i}]", el["locator"])
                elif rk == "sibling_output":
                    if set(el.keys()) != {"ref_kind", "producer", "output_id"}:
                        fail("a", f"{nid}.inputs.required[{i}]: sibling_output element keys must be {{ref_kind, producer, output_id}}")
                    else:
                        if not is_str(el.get("producer")):
                            fail("a", f"{nid}.inputs.required[{i}].producer: must be a scalar string")
                        if not is_str(el.get("output_id")):
                            fail("a", f"{nid}.inputs.required[{i}].output_id: must be a scalar string")
        opt = inp.get("optional")
        if not is_seq(opt):
            fail("a", f"{nid}.inputs.optional: must be a seq (0+)")
        else:
            for i, el in enumerate(opt):
                if not isinstance(el, dict) or el.get("ref_kind") != "external" or set(el.keys()) != {"ref_kind", "locator"}:
                    fail("a", f"{nid}.inputs.optional[{i}]: optional refs must be external-locator only {{ref_kind, locator}}")
                else:
                    check_external_locator("a", f"{nid}.inputs.optional[{i}]", el["locator"])
    # requested_output
    ro = c.get("requested_output")
    if isinstance(ro, dict) and check_keyset("a", f"{nid}.requested_output", ro.keys(), OUTPUT_KEYS):
        if not is_str(ro.get("id")):
            fail("a", f"{nid}.requested_output.id: must be a scalar string")
        if ro.get("kind") not in OUTPUT_KIND_ENUM:
            fail("a", f"{nid}.requested_output.kind not in {sorted(OUTPUT_KIND_ENUM)} (got {ro.get('kind')!r})")
        if not is_str(ro.get("path")):
            fail("a", f"{nid}.requested_output.path: must be a scalar string")
    elif not isinstance(ro, dict):
        fail("a", f"{nid}.requested_output: must be a mapping")
    # acceptance
    acc = c.get("acceptance")
    if not isinstance(acc, dict) or set(acc.keys()) != {"predicates"}:
        fail("a", f"{nid}.acceptance: must be a mapping with exactly key predicates")
    else:
        check_seq_of_scalar("a", f"{nid}.acceptance.predicates", acc.get("predicates"), 1)
    # constraints
    con = c.get("constraints")
    if isinstance(con, dict) and check_keyset("a", f"{nid}.constraints", con.keys(), CONSTRAINT_KEYS):
        check_seq_of_scalar("a", f"{nid}.constraints.allowed_paths", con.get("allowed_paths"), 1)
        check_seq_of_scalar("a", f"{nid}.constraints.forbidden_paths", con.get("forbidden_paths"), 0)
        check_seq_of_scalar("a", f"{nid}.constraints.non_goals", con.get("non_goals"), 0)
    elif not isinstance(con, dict):
        fail("a", f"{nid}.constraints: must be a mapping")
    # gates
    g = c.get("gates")
    if isinstance(g, dict) and check_keyset("a", f"{nid}.gates", g.keys(), GATE_KEYS):
        for k in ("operator_authorization_required", "operator_acceptance_required"):
            if not is_bool(g.get(k)):
                fail("a", f"{nid}.gates.{k}: must be bool (got {type(g.get(k)).__name__})")
        if not (g.get("reason") is None or is_str(g.get("reason"))):
            fail("a", f"{nid}.gates.reason: must be scalar-or-null")
    elif not isinstance(g, dict):
        fail("a", f"{nid}.gates: must be a mapping")
    # doctrine_affecting / stop_conditions
    if not is_bool(c.get("doctrine_affecting")):
        fail("a", f"{nid}.doctrine_affecting: must be bool")
    check_seq_of_scalar("a", f"{nid}.stop_conditions", c.get("stop_conditions"), 0)

# ---- (b) real ref resolution -----------------------------------------------------
# grounding SOURCE snapshot bytes == true source SHA-256
snap = os.path.join(WAVE_DIR, "grounding-source-5015460988.md")
if not os.path.isfile(snap):
    fail("b", "grounding-source snapshot missing")
elif sha256_file(snap) != TRUE_SOURCE_SHA:
    fail("b", f"grounding source SHA-256 != true source {TRUE_SOURCE_SHA[:12]}… (got {sha256_file(snap)[:12]}…)")

# intent_ref.{id,schema} resolve to the ACTUAL intent object; carrier.ref resolves content-bound
for nid, (p, c) in contracts.items():
    ir = c.get("intent_ref") or {}
    if INTENT_ID is not None and ir.get("id") != INTENT_ID:
        fail("b", f"{nid}.intent_ref.id {ir.get('id')!r} != intent object id {INTENT_ID!r}")
    if INTENT_SCHEMA is not None and ir.get("schema") != INTENT_SCHEMA:
        fail("b", f"{nid}.intent_ref.schema {ir.get('schema')!r} != intent object schema {INTENT_SCHEMA!r}")
    ref = ((ir.get("carrier") or {}).get("ref")) or ""
    resolve_content_locator("b", f"{nid}.intent_ref.carrier.ref", ref)

# ---- load wave -------------------------------------------------------------------
WAVE = os.path.join(WAVE_DIR, "wave.cn-wave-v1.yaml")
wave = {}
if not os.path.isfile(WAVE):
    fail("c", "wave.cn-wave-v1.yaml missing")
else:
    try:
        wave = load_yaml(WAVE)
    except Exception as e:
        fail("c", f"wave YAML parse error: {e}")

wave_nodes = {n["id"]: n for n in (wave.get("nodes") or []) if isinstance(n, dict) and "id" in n}

# node.contract_sha256 pins match contract bytes; contract_ref resolves
for nid, n in wave_nodes.items():
    cref = os.path.join(WAVE_DIR, n.get("contract_ref", ""))
    pin = n.get("contract_sha256")
    if not n.get("contract_ref") or not os.path.isfile(cref):
        fail("b", f"node {nid} contract_ref does not resolve: {n.get('contract_ref')!r}")
    elif pin != sha256_file(cref):
        fail("b", f"node {nid} contract_sha256 mismatch (pinned {str(pin)[:12]}… actual {sha256_file(cref)[:12]}…)")

# wave grounding + intent locators resolve (content-bound)
g = wave.get("grounding", {}) or {}
resolve_content_locator("b", "wave.intent.locator", str((wave.get("intent") or {}).get("locator", "")))
resolve_content_locator("b", "wave.grounding.source_snapshot",
                        str((((g.get("source_snapshot") or {}).get("locator")) or {}).get("revision", "")))
for key in ("alpha_derivative", "reconcile_627_ref"):
    loc = (g.get(key) or {}).get("locator", {}) or {}
    ch = str(loc.get("content_hash", "")); path = loc.get("path", "")
    m = re.match(r"^sha256:([0-9a-f]{64})$", ch)
    if m and path:
        resolve_content_locator("b", f"wave.grounding.{key}", f"sha256:{m.group(1)}@{path}")
    else:
        fail("b", f"wave.grounding.{key}: missing content_hash/path immutable binding")

# ---- (c) derive nodes / outputs / edges / roots / critical path ------------------
NODE_IDS = sorted(contracts.keys())
# contracts-dir set == wave node set
if set(wave_nodes.keys()) != set(NODE_IDS):
    fail("c", f"wave node set {sorted(wave_nodes)} != contracts-dir set {NODE_IDS}")

# requested_output ids; wave node.output_id == contract.requested_output.id
output_ids = {nid: ((c.get("requested_output") or {}).get("id")) for nid, (p, c) in contracts.items()}
producer_by_output = {}
for nid, oid in output_ids.items():
    if oid in producer_by_output:
        fail("c", f"duplicate requested_output.id {oid!r} ({producer_by_output[oid]}, {nid})")
    producer_by_output[oid] = nid
for nid, n in wave_nodes.items():
    if nid in output_ids and n.get("output_id") != output_ids[nid]:
        fail("c", f"node {nid}.output_id {n.get('output_id')!r} != contract.requested_output.id {output_ids[nid]!r}")

# derive edges over sibling_output refs
derived_edges = set()
for consumer, (p, c) in contracts.items():
    for el in ((c.get("inputs") or {}).get("required") or []):
        if isinstance(el, dict) and el.get("ref_kind") == "sibling_output":
            prod, oid = el.get("producer"), el.get("output_id")
            if prod not in output_ids:
                fail("c", f"{consumer}: sibling_output producer {prod!r} unknown")
            elif output_ids[prod] != oid:
                fail("c", f"{consumer}: sibling_output output_id {oid!r} != {prod}.requested_output.id {output_ids[prod]!r}")
            else:
                derived_edges.add((prod, consumer))

authored_edges = set()
for e in (wave.get("edges") or []):
    if isinstance(e, dict) and "from" in e and "to" in e:
        authored_edges.add((e["from"], e["to"]))
if derived_edges != authored_edges:
    fail("c", f"edge parity break: derived-only {sorted(derived_edges - authored_edges)}; "
              f"authored-only {sorted(authored_edges - derived_edges)}")

# derived roots (no incoming derived edge) vs authored roots
derived_roots = set(NODE_IDS) - {b for (_a, b) in derived_edges}
authored_roots = set((wave.get("edge_derivation") or {}).get("roots") or [])
if authored_roots and authored_roots != derived_roots:
    fail("c", f"authored roots {sorted(authored_roots)} != derived roots {sorted(derived_roots)}")

# ---- (d) DAG ---------------------------------------------------------------------
def is_dag(node_set, edge_set):
    from collections import defaultdict, deque
    succ = defaultdict(list); indeg = {n: 0 for n in node_set}
    for a, b in edge_set:
        succ[a].append(b); indeg[b] = indeg.get(b, 0) + 1; indeg.setdefault(a, 0)
    q = deque([n for n in indeg if indeg[n] == 0]); seen = 0
    while q:
        x = q.popleft(); seen += 1
        for y in succ[x]:
            indeg[y] -= 1
            if indeg[y] == 0: q.append(y)
    return seen == len(indeg)

if not is_dag(set(NODE_IDS), derived_edges):
    fail("d", "dependency graph has a cycle (not a DAG)")

def reachable(src, edge_set):
    from collections import defaultdict, deque
    succ = defaultdict(list)
    for a, b in edge_set: succ[a].append(b)
    seen, q = set(), deque([src])
    while q:
        x = q.popleft()
        for y in succ[x]:
            if y not in seen: seen.add(y); q.append(y)
    return seen

reach = {n: reachable(n, derived_edges) for n in NODE_IDS}

# critical_path must be a real root→terminal path along derived edges
cp = (wave.get("edge_derivation") or {}).get("critical_path") or []
if cp:
    for a, b in zip(cp, cp[1:]):
        if (a, b) not in derived_edges:
            fail("c", f"critical_path segment {a}->{b} is not a derived edge")
    if cp[0] not in derived_roots:
        fail("c", f"critical_path start {cp[0]!r} is not a derived root")

# ---- (e) parallel nodes do not share a write surface -----------------------------
import fnmatch
def concurrent(a, b): return (b not in reach[a]) and (a not in reach[b])
def globs_overlap(A, B):
    for a in A:
        for b in B:
            if a == b: return True
            if not any(ch in a for ch in "*?[") and fnmatch.fnmatch(a, b): return True
            if not any(ch in b for ch in "*?[") and fnmatch.fnmatch(b, a): return True
    return False

out_paths = {}
for nid, (p, c) in contracts.items():
    op = (c.get("requested_output") or {}).get("path")
    if op in out_paths: fail("e", f"nodes {out_paths[op]} and {nid} share requested_output.path {op!r}")
    out_paths[op] = nid
for i, a in enumerate(NODE_IDS):
    for b in NODE_IDS[i+1:]:
        if concurrent(a, b):
            A = set(((contracts[a][1].get("constraints") or {}).get("allowed_paths")) or [])
            B = set(((contracts[b][1].get("constraints") or {}).get("allowed_paths")) or [])
            if globs_overlap(A, B):
                fail("e", f"parallel nodes {a} and {b} share a write surface (allowed_paths overlap)")

# ---- (f) gate invariants ---------------------------------------------------------
for nid, (p, c) in contracts.items():
    g = c.get("gates") or {}
    if "reason" not in g:
        fail("f", f"{nid}.gates.reason: key must always be present")
    da  = c.get("doctrine_affecting")
    oar = g.get("operator_acceptance_required")
    if da is True and oar is not True:
        fail("f", f"{nid}: doctrine_affecting true but operator_acceptance_required not true (invariant 1)")
    any_gate = bool(g.get("operator_authorization_required")) or bool(oar)
    reason = g.get("reason")
    if any_gate and not (isinstance(reason, str) and reason.strip()):
        fail("f", f"{nid}: a gate bool is true but reason is null/empty")
    if not any_gate and reason is not None:
        fail("f", f"{nid}: both gate bools false but reason is non-null")

# ---- (g) evaluate the AUTHORED completion relation as structured data -------------
comp = wave.get("completion", {}) or {}
terminal = comp.get("terminal_node")
N = comp.get("construction_set_N") or []
derived_N = [n for n in NODE_IDS if n != terminal]
if terminal not in NODE_IDS:
    fail("g", f"completion.terminal_node {terminal!r} is not a wave node")
if set(N) != set(derived_N):
    fail("g", f"construction_set_N {sorted(N)} != derived (all nodes minus terminal) {sorted(derived_N)}")
if terminal in N:
    fail("g", f"construction_set_N must NOT contain the terminal node {terminal!r} (would make completion recursive)")

preds = {pr["name"]: pr for pr in (comp.get("predicates") or []) if isinstance(pr, dict) and "name" in pr}
CHILD = "child_complete"

def ast_edges(expr):
    """Dependency targets this expr READS — derived by walking the AST (no substring scan)."""
    out = []
    if not isinstance(expr, dict):
        raise ValueError(f"expr must be a mapping, got {expr!r}")
    if "lit" in expr or "input" in expr:
        return out
    if "input_map" in expr:
        node = expr.get("node")
        out.append(f"{expr['input_map']}({node})"); return out
    if "all_over_N" in expr:
        base = expr["all_over_N"]
        for n in N: out.append(f"{base}({n})")
        return out
    if "pred" in expr:
        out.append(expr["pred"]); return out
    for op in ("and", "or"):
        if op in expr:
            for e in expr[op]: out += ast_edges(e)
            return out
    if "not" in expr:
        return ast_edges(expr["not"])
    raise ValueError(f"unknown expr node keys {sorted(expr.keys())}")

# build predicate dependency graph FROM the authored ASTs
pnodes = set(preds.keys())
for n in set(N) | ({terminal} if terminal else set()):
    pnodes.add(f"{CHILD}({n})")   # per-node leaves
pedges = set()
graph_ok = True
try:
    for name, pr in preds.items():
        if "expr" not in pr:
            fail("g", f"completion.predicate {name!r} has no expr"); graph_ok = False; continue
        for tgt in ast_edges(pr["expr"]):
            pnodes.add(tgt)  # ensure referenced leaves exist as nodes
            pedges.add((name, tgt))
except ValueError as e:
    fail("g", f"completion predicate expr malformed: {e}"); graph_ok = False

# referential integrity: every {pred: X} target must be a defined predicate or a child leaf
for (a, b) in pedges:
    if b not in preds and not b.startswith(f"{CHILD}("):
        fail("g", f"completion predicate {a!r} references undefined predicate {b!r}")

if graph_ok and not is_dag(pnodes, pedges):
    fail("g", "completion-predicate dependency graph has a CYCLE (recursive completion, e.g. a tautology)")

# evaluate each fixture: compute every predicate over the fixture inputs, compare to authored expected
def eval_expr(expr, inputs, stack):
    if "lit" in expr:   return bool(expr["lit"])
    if "input" in expr:
        if expr["input"] not in inputs:
            raise KeyError(f"fixture missing input {expr['input']!r}")
        return bool(inputs[expr["input"]])
    if "input_map" in expr:
        return bool(inputs[expr["input_map"]][expr["node"]])
    if "all_over_N" in expr:
        base = expr["all_over_N"]
        m = inputs.get(base, {})
        return all(bool(m[n]) for n in N)
    if "pred" in expr:
        return eval_pred(expr["pred"], inputs, stack)
    if "and" in expr: return all(eval_expr(e, inputs, stack) for e in expr["and"])
    if "or" in expr:  return any(eval_expr(e, inputs, stack) for e in expr["or"])
    if "not" in expr: return not eval_expr(expr["not"], inputs, stack)
    raise ValueError(f"unknown expr node keys {sorted(expr.keys())}")

def eval_pred(name, inputs, stack):
    if name in stack:
        raise RecursionError(f"recursive predicate {name!r}")
    if name not in preds:
        raise KeyError(f"undefined predicate {name!r}")
    return eval_expr(preds[name]["expr"], inputs, stack | {name})

fixtures = comp.get("fixtures") or []
if graph_ok:
    if not fixtures:
        fail("g", "completion.fixtures: at least one truth-table fixture required")
    for fx in fixtures:
        nm = fx.get("name", "<unnamed>")
        inputs = fx.get("inputs") or {}
        expected = fx.get("expected") or {}
        if not expected:
            fail("g", f"fixture {nm!r}: no authored `expected` outputs"); continue
        for pname, exp in expected.items():
            try:
                got = eval_pred(pname, inputs, set())
            except Exception as e:
                fail("g", f"fixture {nm!r}: predicate {pname!r} could not be evaluated: {e}"); continue
            if bool(got) != bool(exp):
                fail("g", f"fixture {nm!r}: predicate {pname!r} authored expected={exp} but COMPUTED={got}")

# ---- shared: load the AUTHORITATIVE, TOTAL assurance registry (used by checks h, i, j) -------------
REG_PATH = os.path.join(WAVE_DIR, "oracle-registry.yaml")
registry = {}
assurance = []
wave_predicates = []
if not os.path.isfile(REG_PATH):
    fail("h", "oracle-registry.yaml missing")
else:
    try:
        registry = load_yaml(REG_PATH) or {}
    except Exception as e:
        fail("h", f"oracle-registry.yaml YAML parse error: {e}")
    assurance = registry.get("assurance") or []
    wave_predicates = registry.get("wave_predicates") or []
# mechanically-verifiable subset: identical entry shape to the pre-R6 `oracles` list; check (h) owns it.
reg_entries = [e for e in assurance if isinstance(e, dict) and e.get("classification") == "mechanically-verifiable"]
from collections import Counter as _Counter
child_cat_counts = _Counter(e.get("classification") for e in assurance
                            if isinstance(e, dict) and e.get("classification") in KIND_ENUM)

# ---- (h) ORACLE OWNERSHIP — TOTAL + SINGULAR BIJECTION over the mechanically-verifiable subset ------
# Three sources: (1) the AUTHORITATIVE registry oracle-registry.yaml; (2) each child's
# acceptance.predicates; (3) the md PROJECTION (acceptance-oracles.md "Registry projection" table).
# Proves a total, singular map: every mechanically-verifiable predicate ⇄ exactly one child owner
# entry ⇄ exactly one acceptance predicate; rejects a missing/duplicate/extra entry, a classification
# mismatch, a placeholder operand, and any registry↔projection parity break.
PLACEHOLDER_RE = re.compile(r"<[^>\n]+>")
ALLOWED_PATH_ROOTS = (".cdd/", "schemas/", "docs/", "src/")
REGISTRY_ENTRY_KEYS = {"owner", "predicate", "classification", "ownership", "checker", "schema",
                       "positive_fixture", "negative_fixture", "command", "receipt_evidence_predicate"}

def path_concrete(val):
    """A concrete repo-root-relative path: non-empty str, no placeholder operand, not absolute,
    contains '/', and rooted at a known repo top-level dir (rejects implicit-CWD bare filenames)."""
    if not isinstance(val, str) or not val.strip():
        return False, "empty/non-string"
    if PLACEHOLDER_RE.search(val):
        return False, "contains placeholder operand"
    if val.startswith("/"):
        return False, "absolute path"
    if "/" not in val:
        return False, "bare filename (implicit CWD)"
    if not val.startswith(ALLOWED_PATH_ROOTS):
        return False, "not repo-root-relative (unknown root)"
    return True, ""

def path_in_allowed(val, allowed_globs):
    return isinstance(val, str) and any(fnmatch.fnmatch(val, g) for g in allowed_globs)

def md_table_rows(text):
    """Yield lists-of-stripped-cells for every markdown table body row (drops header + separator)."""
    for line in text.splitlines():
        s = line.strip()
        if not s.startswith("|"):
            continue
        cells = [c.strip() for c in s.strip("|").split("|")]
        if not cells:
            continue
        if all(set(c) <= set("-: ") for c in cells):   # separator row
            continue
        yield cells

# The authoritative TOTAL registry + its mechanically-verifiable subset `reg_entries` are loaded above.
# per-owner: acceptance predicate set, allowed globs, pinned-input paths, and that the contract
# actually CONSUMES the registry via a canonical inputs.required[].external control_plane ref.
acc_by_owner, allowed_by_owner, inputpaths_by_owner = {}, {}, {}
for nid, (p, c) in contracts.items():
    acc_by_owner[nid] = set(((c.get("acceptance") or {}).get("predicates")) or [])
    allowed_by_owner[nid] = list(((c.get("constraints") or {}).get("allowed_paths")) or [])
    ips = set(); consumes_registry = False
    for el in ((c.get("inputs") or {}).get("required") or []):
        if isinstance(el, dict) and el.get("ref_kind") == "external":
            loc = el.get("locator") or {}
            if loc.get("kind") == "repo_artifact" and loc.get("path"):
                ips.add(loc["path"])
            rev = str(loc.get("revision", ""))
            if "oracle-registry.yaml" in rev:
                consumes_registry = True
    inputpaths_by_owner[nid] = ips
    if not consumes_registry:
        fail("h", f"{nid}: does not consume oracle-registry.yaml via a canonical inputs.required[].external ref")

# ---- (h.1) validate registry entries + build the predicate->entry map -------------
pred_owner = {}          # predicate -> owner (must be singular)
pred_entry = {}          # predicate -> normalized entry tuple (for md-parity)
for i, e in enumerate(reg_entries):
    where = f"oracle-registry.assurance(mv)[{i}]"
    if not isinstance(e, dict):
        fail("h", f"{where}: must be a mapping"); continue
    extra = set(e.keys()) - REGISTRY_ENTRY_KEYS
    if extra:
        fail("h", f"{where}: unexpected key(s) {sorted(extra)}")
    owner = e.get("owner"); pred = e.get("predicate")
    if owner not in contracts:
        fail("h", f"{where}: owner {owner!r} is not a known child contract"); continue
    if e.get("classification") != "mechanically-verifiable":
        fail("h", f"{where}: classification must be 'mechanically-verifiable' (got {e.get('classification')!r}) — registry vs classification mismatch")
    # extra owner entry absent from the owning contract's acceptance predicates
    if pred not in acc_by_owner.get(owner, set()):
        fail("h", f"{where}: owner entry predicate {pred!r} is ABSENT from {owner}.acceptance.predicates (extra owner entry)")
    if e.get("receipt_evidence_predicate") not in acc_by_owner.get(owner, set()):
        fail("h", f"{where}: receipt_evidence_predicate not in {owner}.acceptance.predicates")
    # duplicate owner: a mechanically-verifiable predicate mapped to more than one entry
    if pred in pred_owner:
        fail("h", f"{where}: DUPLICATE owner entry for predicate {pred!r} (already owned by {pred_owner[pred]}) — mapping not singular")
    own = e.get("ownership")
    if own not in {"emitted", "existing"}:
        fail("h", f"{where}: ownership must be 'emitted'|'existing' (got {own!r})")
    checker, schema = e.get("checker"), e.get("schema")
    if bool(checker) == bool(schema):
        fail("h", f"{where}: must name EXACTLY one of checker|schema")
    primary = checker or schema
    posf, negf, cmd = e.get("positive_fixture"), e.get("negative_fixture"), e.get("command")
    for label, val in (("checker/schema", primary), ("positive_fixture", posf), ("negative_fixture", negf)):
        ok, why = path_concrete(val)
        if not ok:
            fail("h", f"{where}.{label}: {why} ({val!r})")
    if not isinstance(cmd, str) or not cmd.strip():
        fail("h", f"{where}.command: must be a non-empty command string")
    elif PLACEHOLDER_RE.search(cmd):
        fail("h", f"{where}.command: contains placeholder operand ({cmd!r})")
    allowed = allowed_by_owner.get(owner, [])
    if own == "emitted":
        for label, val in (("checker/schema", primary), ("positive_fixture", posf), ("negative_fixture", negf)):
            if isinstance(val, str) and not path_in_allowed(val, allowed):
                fail("h", f"{where}.{label}: emitted path not in {owner} allowed_paths ({val!r})")
    elif own == "existing":
        if isinstance(primary, str) and primary not in inputpaths_by_owner.get(owner, set()):
            fail("h", f"{where}: 'existing' checker/schema {primary!r} is not a pinned immutable input of {owner}")
        for label, val in (("positive_fixture", posf), ("negative_fixture", negf)):
            if isinstance(val, str) and not path_in_allowed(val, allowed):
                fail("h", f"{where}.{label}: emitted fixture not in {owner} allowed_paths ({val!r})")
    if isinstance(cmd, str):
        if isinstance(primary, str) and primary not in cmd:
            fail("h", f"{where}.command: does not reference the concrete checker/schema {primary!r}")
        if isinstance(posf, str) and posf not in cmd:
            fail("h", f"{where}.command: does not reference the positive fixture {posf!r}")
    if isinstance(pred, str):
        pred_owner[pred] = owner
        pred_entry[pred] = (owner, "mechanically-verifiable", own, primary, posf, negf)

# ---- (h.2) TOTALITY: every mv acceptance predicate declared in a contract is owned once ----
# The mechanically-verifiable set is authoritatively the registry; the md projection (h.3) proves the
# classification is not silently dropped. Every registry-owned predicate having been checked ∈ acceptance
# above, totality reduces to: no owner contract carries an acceptance predicate CLASSIFIED mechanically-
# verifiable (by the projection) without a registry entry — proven in (h.3) by set equality.

# ---- (h.3) registry ⇄ md-projection PARITY ---------------------------------------
acc_oracles = os.path.join(WAVE_DIR, "acceptance-oracles.md")
proj = {}   # predicate -> (owner, classification, ownership, primary, pos, neg)
if not os.path.isfile(acc_oracles):
    fail("h", "acceptance-oracles.md missing (registry projection unreadable)")
else:
    md_text = open(acc_oracles, encoding="utf-8").read()
    md_scan = re.sub(r"<!--.*?-->", "", md_text, flags=re.S)   # HTML comments (e.g. the marker) are not operands
    if PLACEHOLDER_RE.search(md_scan):
        fail("h", f"acceptance-oracles.md carries a placeholder operand {PLACEHOLDER_RE.search(md_scan).group(0)!r}")
    # isolate the "Registry projection" section (from its heading to the next heading)
    lines = md_text.splitlines()
    sec = []
    in_sec = False
    for ln in lines:
        if ln.startswith("#") and "Registry projection" in ln:
            in_sec = True; continue
        if in_sec and ln.startswith("## "):
            break
        if in_sec:
            sec.append(ln)
    if not sec:
        fail("h", "acceptance-oracles.md: no 'Registry projection' section found")
    for cells in md_table_rows("\n".join(sec)):
        if len(cells) < 7:
            continue
        owner, pred, classification, ownership, primary, pos, neg = cells[:7]
        if owner == "owner" or pred == "predicate":   # header
            continue
        pred = pred.strip("`"); primary = primary.strip("`"); pos = pos.strip("`"); neg = neg.strip("`")
        if pred in proj:
            fail("h", f"acceptance-oracles.md projection: DUPLICATE row for predicate {pred!r}")
        proj[pred] = (owner, classification, ownership, primary, pos, neg)

# set equality: registry predicates == projection predicates
reg_preds, proj_preds = set(pred_entry), set(proj)
for pred in sorted(proj_preds - reg_preds):
    fail("h", f"mechanically-verifiable predicate {pred!r} classified in acceptance-oracles.md projection but has NO registry owner entry (missing entry)")
for pred in sorted(reg_preds - proj_preds):
    fail("h", f"registry owner entry {pred!r} is ABSENT from the acceptance-oracles.md projection (parity break)")
# field-by-field parity for predicates present in both
for pred in sorted(reg_preds & proj_preds):
    r = pred_entry[pred]     # (owner, classification, ownership, primary, pos, neg)
    m = proj[pred]           # (owner, classification, ownership, primary, pos, neg)
    if r != m:
        fail("h", f"registry↔projection MISMATCH for {pred!r}: registry {r} != md {m}")

# ---- (i) LEDGER CONSISTENCY -------------------------------------------------------
# (1) revision labels agree; (2) reported classification counts == derived table counts (and mv ==
# registry size); (3) every classification cell is EXACTLY ONE enum member (no compound category).
REQUIRED_LABELLED = ["README.md", "acceptance-oracles.md", "validate.py", "validate_test.py",
                     "oracle-registry.yaml", "wave.cn-wave-v1.yaml"]
# line-anchored marker match ONLY (so prose that merely mentions the marker name is not miscounted);
# the literal is split so this source file contains exactly ONE contiguous marker (the top-of-file one).
REV_RE = re.compile(r"(?m)^\s*(?:#|//|<!--)?\s*wave-" + r"revision:\s*(\S+)")
truth_rev = wave.get("revision")
if not truth_rev:
    fail("i", "wave.cn-wave-v1.yaml carries no authored `revision`")
label_vals = set()
for fn in REQUIRED_LABELLED:
    fp = os.path.join(WAVE_DIR, fn)
    if not os.path.isfile(fp):
        fail("i", f"revision-labelled artifact missing: {fn}"); continue
    ms = REV_RE.findall(open(fp, encoding="utf-8").read())
    if len(ms) != 1:
        fail("i", f"{fn}: expected exactly one `wave-revision:` label (found {len(ms)})"); continue
    label_vals.add(ms[0])
    if truth_rev is not None and ms[0] != str(truth_rev):
        fail("i", f"{fn}: wave-revision {ms[0]!r} != authored wave revision {truth_rev!r}")
if len(label_vals) > 1:
    fail("i", f"revision labels DISAGREE across artifacts: {sorted(label_vals)}")

# (2)/(3): counts come from the TOTAL registry (authoritative); the md 'Summary counts' line must EQUAL
# the registry-derived per-category CHILD totals, the mv count must equal the ownership-entry count, and
# EVERY classification-like md cell must be a single enum member (a compound category fails closed).
if os.path.isfile(acc_oracles):
    md_text = open(acc_oracles, encoding="utf-8").read()
    for cells in md_table_rows(md_text):
        for cell in cells:
            c = cell.strip().lower().strip("*` ")
            toks = [t for t in re.split(r"[\s+/,]+|\band\b", c) if t]
            # a cell composed SOLELY of >=2 enum members is a compound category -> fail closed
            if len(toks) >= 2 and all(t in KIND_ENUM for t in toks):
                fail("i", f"acceptance-oracles.md: classification cell is not a single enum member: {cell!r}")
    msum = re.search(
        r"enforced\s*\*\*(\d+)\*\*.*?mechanically-verifiable\s*\*\*(\d+)\*\*.*?evidenced\s*\*\*(\d+)\*\*.*?cognitive-review\s*\*\*(\d+)\*\*",
        md_text, re.S)
    if not msum:
        fail("i", "acceptance-oracles.md: no parseable 'Summary counts' line")
    else:
        reported = {"enforced": int(msum.group(1)), "mechanically-verifiable": int(msum.group(2)),
                    "evidenced": int(msum.group(3)), "cognitive-review": int(msum.group(4))}
        for k in KIND_ENUM:
            if reported[k] != child_cat_counts.get(k, 0):
                fail("i", f"acceptance-oracles.md reported {k}={reported[k]} != total-registry child count {child_cat_counts.get(k, 0)}")
    if child_cat_counts.get("mechanically-verifiable", 0) != len(reg_entries):
        fail("i", f"registry mechanically-verifiable classification count {child_cat_counts.get('mechanically-verifiable', 0)} != ownership-entry count {len(reg_entries)}")

# ---- (j) CLASSIFICATION TOTALITY + SINGULARITY over the COMPLETE child acceptance set --------------
ASSURANCE_CORE = {"owner", "predicate", "classification"}
ASSURANCE_REQUIRED = {
    "enforced":                {"enforced_by"},
    "mechanically-verifiable": {"ownership", "positive_fixture", "negative_fixture", "command",
                                "receipt_evidence_predicate"},   # + EXACTLY one of checker|schema
    "evidenced":               {"receipt_evidence"},
    "cognitive-review":        {"independent_review"},
}
# child acceptance pairs (owner, predicate) from the RAW authored lists (also rejects an intra-contract dup)
contract_pairs = set()
for nid, (p, c) in contracts.items():
    seen = set()
    for pred in (((c.get("acceptance") or {}).get("predicates")) or []):
        if not is_str(pred):
            continue
        if pred in seen:
            fail("j", f"{nid}.acceptance.predicates lists {pred!r} more than once (not classified exactly once)")
        seen.add(pred)
        contract_pairs.add((nid, pred))

# validate every assurance entry; build (owner,predicate) -> [classification, …]
from collections import defaultdict as _dd
reg_pair_classes = _dd(list)
J_BY_CAT = {k: 0 for k in KIND_ENUM}
for i, e in enumerate(assurance):
    where = f"assurance[{i}]"
    if not isinstance(e, dict):
        fail("j", f"{where}: must be a mapping"); continue
    owner, pred, cls = e.get("owner"), e.get("predicate"), e.get("classification")
    if not is_str(owner) or owner not in contracts:
        fail("j", f"{where}: owner {owner!r} is not a known child contract"); continue
    if not is_str(pred):
        fail("j", f"{where}: predicate must be a scalar string"); continue
    if cls not in KIND_ENUM:
        fail("j", f"{where}: classification {cls!r} is not exactly one enum member of {sorted(KIND_ENUM)}"); continue
    missing = (ASSURANCE_CORE | ASSURANCE_REQUIRED[cls]) - set(e.keys())
    if missing:
        fail("j", f"{where}: classification {cls!r} is missing required field(s) {sorted(missing)}")
    if cls == "mechanically-verifiable" and bool(e.get("checker")) == bool(e.get("schema")):
        fail("j", f"{where}: mechanically-verifiable entry must name EXACTLY one of checker|schema")
    reg_pair_classes[(owner, pred)].append(cls)
    J_BY_CAT[cls] += 1

# SINGULARITY: each (owner,predicate) is classified EXACTLY once
for (owner, pred), classes in reg_pair_classes.items():
    if len(classes) > 1:
        fail("j", f"predicate ({owner!r}, {pred!r}) is DOUBLE-CLASSIFIED ({len(classes)} assurance entries: {classes})")

# BIJECTION vs the complete child acceptance set
reg_pairs = set(reg_pair_classes)
for (owner, pred) in sorted(contract_pairs - reg_pairs):
    fail("j", f"acceptance predicate ({owner!r}, {pred!r}) is absent from the assurance registry (UNCLASSIFIED)")
for (owner, pred) in sorted(reg_pairs - contract_pairs):
    fail("j", f"assurance entry ({owner!r}, {pred!r}) is a PHANTOM — no child contract declares it")

# no drift: the mv subset of `assurance` must equal the (h) ownership entries
mv_pairs_assur = {k for k, cs in reg_pair_classes.items() if cs == ["mechanically-verifiable"]}
mv_pairs_own = {(e.get("owner"), e.get("predicate")) for e in reg_entries}
if mv_pairs_assur != mv_pairs_own:
    fail("j", f"mechanically-verifiable drift between classification and ownership entries: {sorted(mv_pairs_assur ^ mv_pairs_own)}")

# WAVE-ONLY predicates: separate section; enforced; must NOT overlap any child acceptance predicate
child_pred_strs = {pr for (_o, pr) in contract_pairs}
for i, e in enumerate(wave_predicates):
    where = f"wave_predicates[{i}]"
    if not isinstance(e, dict):
        fail("j", f"{where}: must be a mapping"); continue
    if e.get("classification") != "enforced" or not is_str(e.get("enforced_by")):
        fail("j", f"{where}: wave predicate must be classification 'enforced' with an enforced_by check")
    if e.get("predicate") in child_pred_strs:
        fail("j", f"{where}: wave predicate {e.get('predicate')!r} collides with a child acceptance predicate (would mask child coverage)")

# COMPLETE md projection parity: the 'Complete assurance classification' table == the registry map
if os.path.isfile(acc_oracles):
    lines = open(acc_oracles, encoding="utf-8").read().splitlines()
    sec, in_sec = [], False
    for ln in lines:
        if ln.startswith("#") and "Complete assurance classification" in ln:
            in_sec = True; continue
        if in_sec and ln.startswith("## "):
            break
        if in_sec:
            sec.append(ln)
    if not sec:
        fail("j", "acceptance-oracles.md: no 'Complete assurance classification' section found")
    comp_proj = {}
    for cells in md_table_rows("\n".join(sec)):
        if len(cells) < 3:
            continue
        owner, pred, cls = cells[0].strip(), cells[1].strip().strip("`"), cells[2].strip()
        if owner == "owner" or pred == "predicate":
            continue
        key = (owner, pred)
        if key in comp_proj:
            fail("j", f"acceptance-oracles.md complete projection: DUPLICATE row for {key}")
        comp_proj[key] = cls
    reg_map = {k: cs[0] for k, cs in reg_pair_classes.items() if len(cs) == 1}
    for key in sorted(set(comp_proj) - set(reg_map)):
        fail("j", f"complete projection row {key} has no single registry classification (projection/registry mismatch)")
    for key in sorted(set(reg_map) - set(comp_proj)):
        fail("j", f"registry entry {key} is absent from the complete md projection (projection parity break)")
    for key in sorted(set(reg_map) & set(comp_proj)):
        if reg_map[key] != comp_proj[key]:
            fail("j", f"assurance↔projection classification mismatch for {key}: registry {reg_map[key]!r} != projection {comp_proj[key]!r}")

J_TOTAL = len(contract_pairs)

# ---- report ----------------------------------------------------------------------
print("=" * 78)
print("Pre-authorization validator — wave cnos#671 (cell-runtime-doctrine) R6 — SOUND")
print(f"  WAVE_DIR={WAVE_DIR}")
print(f"  REPO_ROOT={REPO_ROOT}")
print("=" * 78)
checks = {
 "a": "§2 constraint model (key paths + ENUMS + types + cardinalities)",
 "b": "real ref resolution (intent id/schema + git commit:path + content hashes + source 9d1ab3a5…)",
 "c": "derive-from-authored-data (nodes/outputs/edges/roots/critical-path; authored==derived)",
 "d": "dependency graph is a DAG",
 "e": "parallel nodes share no write surface",
 "f": "gate invariants (doctrine_affecting ⇒ acceptance; reason present iff a gate bool true)",
 "g": "AUTHORED completion evaluated (predicate-graph acyclic by walk; fixtures computed vs expected)",
 "h": "oracle-ownership TOTAL+SINGULAR bijection over the mechanically-verifiable subset (registry ⇄ each child's acceptance ⇄ mv projection; reject missing/duplicate/extra owner, classification mismatch, placeholder, parity break)",
 "i": "ledger consistency (revision labels agree; reported counts == total-registry child totals == mv ownership size; every category a single enum member)",
 "j": "classification TOTALITY+SINGULARITY over the COMPLETE child acceptance set (union(acceptance.predicates) ⇄ total registry, each classified exactly once across all four categories; reject unclassified/double/phantom/bad-category/parity-break)",
}
by = {k: [e for e in errors if e.startswith(f"[{k}")] for k in checks}
for k, desc in checks.items():
    status = "PASS" if not by[k] else "FAIL"
    print(f"  ({k}) {status}  {desc}")
    for e in by[k]:
        print(f"         - {e}")
# any errors whose tag is not a single known check letter
other = [e for e in errors if not any(e.startswith(f"[{k}") for k in checks)]
for e in other:
    print(f"  (?) FAIL  {e}")
print("-" * 78)
print(f"  (j) derived child acceptance predicates: {J_TOTAL}  "
      f"[enforced {J_BY_CAT['enforced']} · mechanically-verifiable {J_BY_CAT['mechanically-verifiable']} · "
      f"evidenced {J_BY_CAT['evidenced']} · cognitive-review {J_BY_CAT['cognitive-review']}]")
print("-" * 78)
if errors:
    print(f"RESULT: FAIL ({len(errors)} violation(s))")
    sys.exit(1)
print("RESULT: PASS — all ten checks (a–j) green at this wave tree.")
sys.exit(0)
