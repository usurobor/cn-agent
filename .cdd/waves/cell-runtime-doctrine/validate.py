#!/usr/bin/env python3
"""
Planning-Cell PRE-AUTHORIZATION validator for wave cnos#671 (cell-runtime-doctrine), R3.

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

TOP_KEYS   = {"schema", "cell", "scope", "intent_ref", "inputs", "requested_output",
              "acceptance", "constraints", "gates", "doctrine_affecting", "stop_conditions"}
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

# ---- report ----------------------------------------------------------------------
print("=" * 78)
print("Pre-authorization validator — wave cnos#671 (cell-runtime-doctrine) R3 — SOUND")
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
if errors:
    print(f"RESULT: FAIL ({len(errors)} violation(s))")
    sys.exit(1)
print("RESULT: PASS — all seven checks green at this wave tree.")
sys.exit(0)
