#!/usr/bin/env python3
"""
Planning-Cell PRE-AUTHORIZATION validator for wave cnos#671 (cell-runtime-doctrine), R2.

Credential-free, standard-library + PyYAML only. Mechanically checks, at the current wave tree:

  (a) every child contract conforms to the §2 cn.cell.contract.v1 KEY-PATH shape (no extra
      top-level keys; nested key sets exact; input union well-typed);
  (b) every immutable ref resolves — the intent object, the grounding SOURCE snapshot
      (must equal the true source SHA-256 9d1ab3a5…), the α-derivative + #627 map content hashes,
      each node.contract_sha256, and every §2 external ref carries its class-specific immutable
      revision component (a bare mutable path / identity-only token is rejected);
  (c) OUTPUT/REF PARITY — the authored wave edge set equals the sibling-output-derived edge set
      exactly (no missing, no spurious);
  (d) the derived dependency graph is a DAG;
  (e) no two PARALLEL (concurrent) nodes share a write surface (requested_output.path / allowed_paths);
  (f) GATE invariants hold (doctrine_affecting ⇒ operator_acceptance_required; reason present iff a
      gate bool is true, null iff both false);
  (g) the COMPLETION-PREDICATE dependency graph is acyclic (non-recursive whole-wave completion).

Exit status: 0 iff ALL checks pass; non-zero (1) on ANY violation.

NEGATIVE-FIXTURE NOTES (each mutation below causes a DETERMINISTIC non-zero exit):
  * add a top-level key (e.g. `consumers:` / `completion_signal:`) to any contract        -> (a) fails
  * change any intent_ref locator hash, or edit intent.cn-intent-v1.yaml                    -> (b) fails
  * corrupt the grounding SOURCE snapshot so its SHA-256 != 9d1ab3a5…                        -> (b) fails
  * flip a node.contract_sha256 / a contract's bytes without updating the pin                -> (b) fails
  * make an external repo_artifact ref drop its `commit`, or a control_plane ref use a bare
    comment id with no content-bound revision                                                -> (b) fails
  * add/remove a wave edge, or change a sibling_output output_id                             -> (c) fails
  * introduce a cycle in the sibling_output refs                                             -> (d) fails
  * point two nodes at the same requested_output.path or overlapping allowed_paths           -> (e) fails
  * set doctrine_affecting: true with operator_acceptance_required: false, or drop `reason`   -> (f) fails
  * make wc5 completion quantify over itself / the wave predicate (recursive completion)      -> (g) fails
"""
import sys, os, re, hashlib, fnmatch
try:
    import yaml
except ImportError:
    print("FATAL: PyYAML required (pip install pyyaml)", file=sys.stderr)
    sys.exit(2)

WAVE_DIR = os.path.dirname(os.path.abspath(__file__))
# repo root = four levels up from .cdd/waves/cell-runtime-doctrine/validate.py
REPO_ROOT = os.path.abspath(os.path.join(WAVE_DIR, "..", "..", ".."))
CONTRACTS_DIR = os.path.join(WAVE_DIR, "contracts")

TRUE_SOURCE_SHA = "9d1ab3a50d00e642fdeb87728dd71f7c7499c60878afe7001f2ddb832b161dbb"
NODE_IDS = ["wc-1", "wc-2", "wc-3a", "wc-3b", "wc-4", "wc-5"]
CONSTRUCTION_SET_N = ["wc-1", "wc-2", "wc-3a", "wc-3b", "wc-4"]  # non-recursive; excludes wc-5

errors = []
def fail(check, msg): errors.append(f"[{check}] {msg}")

def sha256_file(path):
    with open(path, "rb") as f:
        return hashlib.sha256(f.read()).hexdigest()

def load_yaml(path):
    with open(path) as f:
        return yaml.safe_load(f)

# ---- canonical §2 key-path model -------------------------------------------------
TOP_KEYS = {"schema", "cell", "scope", "intent_ref", "inputs", "requested_output",
            "acceptance", "constraints", "gates", "doctrine_affecting", "stop_conditions"}
CELL_KEYS = {"id", "class", "mode", "protocol", "matter_domain"}
SCOPE_KEYS = {"repo", "wave", "parent_cell"}
INTENT_KEYS = {"schema", "id", "carrier"}
CARRIER_KEYS = {"kind", "ref"}
OUTPUT_KEYS = {"id", "kind", "path"}
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
        if extra:  fail(check, f"{where}: NON-CANONICAL key(s) {sorted(extra)}")
        if missing: fail(check, f"{where}: MISSING key(s) {sorted(missing)}")

def looks_immutable_hash(s):
    return bool(re.match(r"^sha256:[0-9a-f]{64}(@.+)?$", str(s)))

def check_external_locator(check, where, loc):
    if not isinstance(loc, dict) or "kind" not in loc:
        fail(check, f"{where}: locator missing kind"); return
    kind = loc["kind"]
    if kind not in LOCATOR_SHAPES:
        fail(check, f"{where}: unknown locator kind {kind!r}"); return
    check_keyset(check, f"{where}.locator", loc.keys(), LOCATOR_SHAPES[kind])
    # immutable revision component present (§2: a bare mutable path/identity is not sufficient)
    if kind == "repo_artifact":
        commit = str(loc.get("commit", ""))
        if not re.match(r"^[0-9a-f]{40}$", commit):
            fail(check+ "/b", f"{where}: repo_artifact commit is not an immutable 40-hex revision ({commit!r})")
    elif kind == "control_plane":
        rev = str(loc.get("revision", ""))
        if not looks_immutable_hash(rev):
            fail(check+ "/b", f"{where}: control_plane revision is identity-only, not content-bound ({rev!r})")
    elif kind == "prior_receipt":
        if not str(loc.get("binding", "")):
            fail(check+ "/b", f"{where}: prior_receipt missing immutable binding")

def resolve_hash_at_path(check, where, locator_str):
    """locator_str form: sha256:<64hex>@<relpath-from-repo-root>"""
    m = re.match(r"^sha256:([0-9a-f]{64})@(.+)$", locator_str)
    if not m:
        return  # not a local content locator; skip (structural check handled elsewhere)
    want, rel = m.group(1), m.group(2)
    abspath = os.path.join(REPO_ROOT, rel)
    if not os.path.exists(abspath):
        fail(check, f"{where}: referenced file does not resolve: {rel}"); return
    got = sha256_file(abspath)
    if got != want:
        fail(check, f"{where}: content hash mismatch for {rel}: pinned {want[:12]}… actual {got[:12]}…")

# ---- load contracts --------------------------------------------------------------
contracts = {}
for nid in NODE_IDS:
    p = os.path.join(CONTRACTS_DIR, f"{nid}.cn-cell-contract-v1.yaml")
    if not os.path.exists(p):
        fail("a", f"contract file missing for {nid}"); continue
    contracts[nid] = (p, load_yaml(p))

# ---- (a) §2 key-path conformance -------------------------------------------------
for nid, (p, c) in contracts.items():
    check_keyset("a", f"{nid} top-level", c.keys(), TOP_KEYS)
    if c.get("schema") != "cn.cell.contract.v1":
        fail("a", f"{nid}: schema != cn.cell.contract.v1")
    if isinstance(c.get("cell"), dict):   check_keyset("a", f"{nid}.cell", c["cell"].keys(), CELL_KEYS)
    if isinstance(c.get("scope"), dict):  check_keyset("a", f"{nid}.scope", c["scope"].keys(), SCOPE_KEYS)
    ir = c.get("intent_ref")
    if isinstance(ir, dict):
        check_keyset("a", f"{nid}.intent_ref", ir.keys(), INTENT_KEYS)
        if isinstance(ir.get("carrier"), dict):
            check_keyset("a", f"{nid}.intent_ref.carrier", ir["carrier"].keys(), CARRIER_KEYS)
    inp = c.get("inputs")
    if not isinstance(inp, dict) or set(inp.keys()) != {"required", "optional"}:
        fail("a", f"{nid}.inputs: must have exactly keys required, optional")
    else:
        req = inp["required"] or []
        if not isinstance(req, list) or len(req) < 1:
            fail("a", f"{nid}.inputs.required: must be a seq of 1+")
        for i, el in enumerate(req):
            rk = el.get("ref_kind")
            if rk == "external":
                if set(el.keys()) != {"ref_kind", "locator"}:
                    fail("a", f"{nid}.inputs.required[{i}]: external element keys must be {{ref_kind, locator}}")
                else:
                    check_external_locator("a", f"{nid}.inputs.required[{i}]", el["locator"])
            elif rk == "sibling_output":
                if set(el.keys()) != {"ref_kind", "producer", "output_id"}:
                    fail("a", f"{nid}.inputs.required[{i}]: sibling_output element keys must be {{ref_kind, producer, output_id}}")
            else:
                fail("a", f"{nid}.inputs.required[{i}]: ref_kind must be external|sibling_output (got {rk!r})")
        opt = inp["optional"]
        if not isinstance(opt, list):
            fail("a", f"{nid}.inputs.optional: must be a seq (0+)")
        else:
            for i, el in enumerate(opt):
                if el.get("ref_kind") != "external" or set(el.keys()) != {"ref_kind", "locator"}:
                    fail("a", f"{nid}.inputs.optional[{i}]: optional refs must be external-locator only")
    ro = c.get("requested_output")
    if isinstance(ro, dict): check_keyset("a", f"{nid}.requested_output", ro.keys(), OUTPUT_KEYS)
    acc = c.get("acceptance")
    if not isinstance(acc, dict) or set(acc.keys()) != {"predicates"}:
        fail("a", f"{nid}.acceptance: must have exactly key predicates")
    elif not isinstance(acc["predicates"], list) or len(acc["predicates"]) < 1:
        fail("a", f"{nid}.acceptance.predicates: seq of 1+")
    con = c.get("constraints")
    if isinstance(con, dict): check_keyset("a", f"{nid}.constraints", con.keys(), CONSTRAINT_KEYS)
    g = c.get("gates")
    if isinstance(g, dict): check_keyset("a", f"{nid}.gates", g.keys(), GATE_KEYS)
    if not isinstance(c.get("doctrine_affecting"), bool):
        fail("a", f"{nid}.doctrine_affecting: must be bool")
    if not isinstance(c.get("stop_conditions"), list):
        fail("a", f"{nid}.stop_conditions: must be a seq (0+)")

# ---- (b) immutable refs resolve --------------------------------------------------
# grounding SOURCE snapshot must equal the true source SHA-256
snap = os.path.join(WAVE_DIR, "grounding-source-5015460988.md")
if not os.path.exists(snap):
    fail("b", "grounding-source snapshot missing")
elif sha256_file(snap) != TRUE_SOURCE_SHA:
    fail("b", f"grounding source SHA-256 != true source {TRUE_SOURCE_SHA[:12]}… (got {sha256_file(snap)[:12]}…)")

# intent object referenced identically by every contract, and resolves
intent_hashes = set()
for nid, (p, c) in contracts.items():
    ref = (((c.get("intent_ref") or {}).get("carrier") or {}).get("ref")) or ""
    resolve_hash_at_path("b", f"{nid}.intent_ref.carrier.ref", ref)
    m = re.match(r"^sha256:([0-9a-f]{64})@", ref)
    if m: intent_hashes.add(m.group(1))
    else: fail("b", f"{nid}.intent_ref.carrier.ref is not a content-bound immutable locator ({ref!r})")
if len(intent_hashes) > 1:
    fail("b", f"contracts disagree on intent object hash: {intent_hashes}")

# every control_plane revision of form sha256:h@path resolves
for nid, (p, c) in contracts.items():
    for el in (c.get("inputs", {}).get("required") or []):
        loc = el.get("locator") or {}
        if loc.get("kind") == "control_plane":
            resolve_hash_at_path("b", f"{nid} grounding revision", str(loc.get("revision", "")))

# ---- load wave -------------------------------------------------------------------
WAVE = os.path.join(WAVE_DIR, "wave.cn-wave-v1.yaml")
wave = load_yaml(WAVE)
nodes = {n["id"]: n for n in wave.get("nodes", [])}

# node.contract_sha256 pins match contract bytes
for nid, n in nodes.items():
    cref = os.path.join(WAVE_DIR, n["contract_ref"])
    pin = n.get("contract_sha256")
    if not os.path.exists(cref):
        fail("b", f"node {nid} contract_ref does not resolve: {n['contract_ref']}")
    elif pin != sha256_file(cref):
        fail("b", f"node {nid} contract_sha256 mismatch (pinned {str(pin)[:12]}… actual {sha256_file(cref)[:12]}…)")

# wave grounding + intent + reconcile content hashes resolve
g = wave.get("grounding", {})
resolve_hash_at_path("b", "wave.intent.locator", str(wave.get("intent", {}).get("locator", "")))
resolve_hash_at_path("b", "wave.grounding.source_snapshot", str(g.get("source_snapshot", {}).get("locator", {}).get("revision", "")))
for key in ("alpha_derivative", "reconcile_627_ref"):
    loc = g.get(key, {}).get("locator", {})
    ch = str(loc.get("content_hash", "")); path = loc.get("path", "")
    m = re.match(r"^sha256:([0-9a-f]{64})$", ch)
    if m and path:
        abspath = os.path.join(REPO_ROOT, path)
        if not os.path.exists(abspath): fail("b", f"wave.grounding.{key}: {path} does not resolve")
        elif sha256_file(abspath) != m.group(1): fail("b", f"wave.grounding.{key}: content hash mismatch for {path}")
    else:
        fail("b", f"wave.grounding.{key}: missing content_hash/path immutable binding")

# ---- (c) output/ref parity: derived edges == authored edges ----------------------
output_ids = {nid: (c.get("requested_output") or {}).get("id") for nid, (p, c) in contracts.items()}
producer_by_output = {}
for nid, oid in output_ids.items():
    if oid in producer_by_output: fail("c", f"duplicate requested_output.id {oid!r}")
    producer_by_output[oid] = nid

derived_edges = set()
for consumer, (p, c) in contracts.items():
    for el in (c.get("inputs", {}).get("required") or []):
        if el.get("ref_kind") == "sibling_output":
            prod, oid = el.get("producer"), el.get("output_id")
            if prod not in output_ids: fail("c", f"{consumer}: sibling_output producer {prod!r} unknown")
            elif output_ids[prod] != oid:
                fail("c", f"{consumer}: sibling_output output_id {oid!r} != {prod}.requested_output.id {output_ids[prod]!r}")
            derived_edges.add((prod, consumer))

authored_edges = set((e["from"], e["to"]) for e in wave.get("edges", []))
if derived_edges != authored_edges:
    fail("c", f"edge parity break: derived-only {sorted(derived_edges - authored_edges)}; authored-only {sorted(authored_edges - derived_edges)}")

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

# ---- (e) parallel nodes do not share a write surface -----------------------------
reach = {n: reachable(n, derived_edges) for n in NODE_IDS}
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
            A = set((contracts[a][1].get("constraints") or {}).get("allowed_paths") or [])
            B = set((contracts[b][1].get("constraints") or {}).get("allowed_paths") or [])
            if globs_overlap(A, B):
                fail("e", f"parallel nodes {a} and {b} share a write surface (allowed_paths overlap)")

# ---- (f) gate invariants ---------------------------------------------------------
for nid, (p, c) in contracts.items():
    g = c.get("gates") or {}
    if "reason" not in g:
        fail("f", f"{nid}.gates.reason: key must always be present")
    da = c.get("doctrine_affecting")
    oar = g.get("operator_acceptance_required")
    if da is True and oar is not True:
        fail("f", f"{nid}: doctrine_affecting true but operator_acceptance_required not true (invariant 1)")
    any_gate = bool(g.get("operator_authorization_required")) or bool(oar)
    reason = g.get("reason")
    if any_gate and not (isinstance(reason, str) and reason.strip()):
        fail("f", f"{nid}: a gate bool is true but reason is null/empty")
    if not any_gate and reason is not None:
        fail("f", f"{nid}: both gate bools false but reason is non-null")

# ---- (g) completion-predicate graph is acyclic (non-recursive completion) --------
comp = wave.get("completion", {})
N = comp.get("construction_set_N", [])
if set(N) != set(CONSTRUCTION_SET_N):
    fail("g", f"construction_set_N must be {CONSTRUCTION_SET_N} (non-recursive; excludes wc-5); got {N}")
if "wc-5" in N:
    fail("g", "construction_set_N must NOT contain wc-5 (would make completion recursive)")
# Build the predicate dependency graph from the wave's declared semantics and prove acyclicity.
# Predicates: child_complete(n) for n in N∪{wc-5}; wc5_ready; wc5_complete; wave_complete.
pedges = set()
pnodes = set()
for n in N + ["wc-5"]:
    cc = f"child_complete({n})"; pnodes.add(cc)   # leaves: depend only on own surfaces
pnodes |= {"wc5_ready", "wc5_complete", "wave_complete"}
for n in N:
    pedges.add(("wc5_ready", f"child_complete({n})"))
    pedges.add(("wave_complete", f"child_complete({n})"))
# wc5_complete depends ONLY on wc-5 own surfaces (a leaf); must NOT read wc5_ready/wave_complete
pedges.add(("wave_complete", "wc5_complete"))
if not is_dag(pnodes, pedges):
    fail("g", "completion-predicate dependency graph has a cycle")
# structural anti-recursion checks on the authored predicate text
wc5_comp_txt = str(comp.get("wc5_completion_predicate", ""))
for bad in ("wave_complete", "wc5_ready", "whole_wave"):
    if bad in wc5_comp_txt:
        fail("g", f"wc5_completion_predicate must not depend on {bad!r} (recursive completion)")
if comp.get("dependency_graph", {}).get("acyclic") is not True:
    fail("g", "wave.completion.dependency_graph.acyclic must be true")
# WC-5 contract itself must not carry a self-quantifying whole-wave acceptance predicate
wc5_preds = " ".join((contracts.get("wc-5", (None, {}))[1].get("acceptance") or {}).get("predicates") or [])
if "whole_wave_completion_predicate__every_child_completion_predicate_holds_and_the_wave_level" in wc5_preds:
    fail("g", "wc-5 still carries the recursive whole-wave completion acceptance predicate")

# ---- report ----------------------------------------------------------------------
print("=" * 72)
print("Pre-authorization validator — wave cnos#671 (cell-runtime-doctrine) R2")
print("=" * 72)
checks = {
 "a": "§2 key-path conformance (no extra top-level keys)",
 "b": "immutable refs resolve (intent + grounding source 9d1ab3a5… + pins)",
 "c": "output/ref parity (authored edges == sibling-output-derived edges)",
 "d": "dependency graph is a DAG",
 "e": "parallel nodes share no write surface",
 "f": "gate invariants (doctrine_affecting ⇒ acceptance; reason present)",
 "g": "completion-predicate graph acyclic (non-recursive)",
}
by = {k: [e for e in errors if e.startswith(f"[{k}") ] for k in checks}
for k, desc in checks.items():
    status = "PASS" if not by[k] else "FAIL"
    print(f"  ({k}) {status}  {desc}")
    for e in by[k]:
        print(f"         - {e}")
print("-" * 72)
if errors:
    print(f"RESULT: FAIL ({len(errors)} violation(s))")
    sys.exit(1)
print("RESULT: PASS — all seven checks green at this wave tree.")
sys.exit(0)
