# KATA-A2

## Extension Registry Engine Compatibility

**Kata ID:** A2-extension-registry-engine-compatibility
**Family:** Local implementation + invariant proof
**Difficulty:** Medium
**Failure Family:** A structural invalid state is accepted because validation is too shallow.
**Expected Governing Skills (hidden):** coding, testing

---

## Public Task Statement

You are extending the cnos runtime extension registry. The design rule is:

> If an extension manifest is incompatible with the current cnos version, it must be rejected deterministically.
> Incompatible manifests must not remain silently active.

Current bug:

- the registry builder validates duplicate names and duplicate op kinds
- but it does not reject manifests whose `engines.cnos` constraint is incompatible with the current runtime version

Your task:

1. propose the smallest coherent fix
2. describe the tests needed
3. if you include code or pseudocode, keep it minimal and implementation-oriented

Do not redesign the whole extension system. Focus on deterministic compatibility rejection.

---

## Public Context

### Existing behavior

The registry builder currently does this:

```ocaml
type ext_manifest = {
  name : string;
  version : string;
  interface : string;
  ops : ext_op list;
  engines_cnos : string option;
}

and ext_op = {
  kind : string;
  class_ : [ `Observe | `Effect ];
}

type registry = {
  by_name : (string, ext_manifest) Hashtbl.t;
  by_op_kind : (string, string) Hashtbl.t;
  rejected : (string * string) list ref; (* ext name, reason *)
}

let build_registry ~current_cnos_version manifests =
  let by_name = Hashtbl.create 16 in
  let by_op_kind = Hashtbl.create 16 in
  let rejected = ref [] in
  List.iter (fun m ->
    (* duplicate name rejection already exists *)
    if Hashtbl.mem by_name m.name then
      rejected := (m.name, "duplicate_name") :: !rejected
    else begin
      Hashtbl.add by_name m.name m;
      List.iter (fun op ->
        if Hashtbl.mem by_op_kind op.kind then
          rejected := (m.name, "duplicate_op_kind") :: !rejected
        else
          Hashtbl.add by_op_kind op.kind m.name
      ) m.ops
    end
  ) manifests;
  { by_name; by_op_kind; rejected }
```

### Desired rule

If `engines_cnos` is present and does not match `current_cnos_version`:

- the manifest is rejected
- none of its ops are added to the active registry
- the rejection reason is explicit
- behavior is deterministic regardless of manifest order

### Example manifests

Manifest A:

```json
{
  "name": "cnos.net.http",
  "version": "1.0.0",
  "engines": { "cnos": ">=3.12.0 <4.0.0" },
  "ops": [{ "kind": "http_get", "class": "observe" }]
}
```

Manifest B:

```json
{
  "name": "org.example.future",
  "version": "1.0.0",
  "engines": { "cnos": ">=4.0.0 <5.0.0" },
  "ops": [{ "kind": "future_op", "class": "observe" }]
}
```

Current runtime version: 3.13.0

---

## Public Deliverable

Return:

1. Invariant
2. Minimal fix
3. Tests
4. Known edge cases

---

## Hidden Evaluator Rubric

### α

Strong answer should:

- state the invariant explicitly:
  - only compatible manifests become active
  - incompatible manifests are rejected before contributing ops
- keep rejection deterministic
- keep the fix local to registry validation/build
- avoid redesigning runtime dispatch or install lifecycle

Weak answer:

- logs a warning but still installs the extension
- rejects too late, after ops are already registered
- makes compatibility a best-effort advisory check

### β

Strong answer should:

- align with runtime-extensions design:
  - compatibility is part of registry discovery/build
  - doctor/traceability can later expose the rejected state
- keep built-ins separate
- avoid collapsing registry validation into runtime dispatch or install/update logic

Weak answer:

- proposes "just hide it from Runtime Contract" but keep it active internally
- mixes registry compatibility with host health or enable/disable policy

### γ

Strong answer should:

- choose a small coherent patch
- mention rejection reason shape
- suggest strong tests
- avoid overclaiming full ecosystem behavior

Weak answer:

- redesigns the whole registry
- leaves "what happens to ops from rejected manifests" ambiguous

### Hidden Checks / Tests

Expected good answer should cover:

**Positive cases:**

- compatible manifest is accepted
- incompatible manifest is rejected
- rejected manifest contributes no ops

**Negative cases:**

- order of manifests does not change compatibility outcome
- incompatible manifest with unique name still rejects
- incompatible manifest with duplicate op kind does not create ambiguous partial state

**Stronger proof signals:**

Bonus if the answer suggests:

- property-like order-independence testing
- clear rejection reasons (`incompatible_engine`)
- testing both with and without `engines_cnos`

### Transfer Target

The evaluator should look for the same pattern transfer as A1:

- invariant-first reasoning
- deterministic rejection
- stronger proof family than a few happy-path examples
