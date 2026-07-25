<!-- wave-revision: R15 -->
# Wave-boundary pre-authorization evidence — oracle-ownership bijection (cnos#671 R15)

This is the **result/evidence binding** referenced by `wave.cn-wave-v1.yaml`
`gates.wave_authorization.preauthorization_gates[]` and by `oracle-registry.yaml`
`wave_predicates.wave_oracle_ownership_bijection_enforced`. It records the **PASS** of the
materialized wave-boundary validator, **content-bound to wave revision R15**. (The validator and its
78⇄78 result are UNCHANGED since R12 — R13 re-pins the evidence to the current revision; the R12
soundness repair below is retained as history.)

The wave authorization is **NOT authorization-ready** unless this validator resolves at its
pinned `validator_sha256` AND its invocation exits 0 with `bijective: true`. Removing or
corrupting `oracle_ownership_bijection.go` breaks the pinned hash and the gate → the wave
returns to the pre-authorization hold (see `wave.stop_conditions.revoked_authorization`).

## R12 soundness repair (external-β R11 BLOCKER)

The R11 validator hand-parsed YAML by indentation/prefix and silently dropped CUE-valid
flow-style predicate lists (`predicates: ["x"]`), producing a demonstrated complete-wave
**false PASS** (an unregistered child predicate reached authorization-ready; R11 reported
"66 child predicates, 66 registry entries, bijective: true"). **R12 removes all ad-hoc YAML
parsing:** every input is normalized through `cue export --out json` (the same CUE that
structurally vets the wave) and decoded via `encoding/json`, so flow and block serialization
are indistinguishable to the program. Owners are derived from **semantic** identity (contract
`cell.id`, cross-checked against the wave node ids), never filename text. Any parse/normalization
loss, or an empty semantic owner/predicate, **fails closed (exit 2)**. Requires `cue` on PATH
(already required by the wave for structural `cue vet`); `os/exec` + `encoding/json` are stdlib,
so no `go.mod` is needed.

## Bound artifact

- **validator:** `.cdd/waves/cell-runtime-doctrine/wave-validators/oracle_ownership_bijection.go`
- **validator_sha256:** `ccf3080289f6c32f81239b7237cbcf7a984111069a8a8ffd97ada7a3e8815bf7`
- **positive_fixture:** `fixtures/oracle-ownership.one-checker-each.positive.yaml`
  (`sha256:485131df76ca9772dcc733bb65a8b55d258314e63096402f52524b495acc17ec`)
- **negative_fixture:** `fixtures/oracle-ownership.double-owned.negative.yaml`
  (`sha256:6569a056ac408cb49633d045d0000528417bf41131efe73b43cf27a046807c86`)
- **flow-style positive fixture:** `fixtures/oracle-ownership.flow-style.positive.yaml`
  (`sha256:4e88dde9ac4641f28787fc7cfa319f3ac7925f015c26bc00c88ea0ebdf2a436c`)
- **authorized_revision:** `R15`

## Invocation (credential-free; no modules, no network; requires `cue` + `go` on PATH)

```
go run oracle_ownership_bijection.go ..                 # real six contracts + registry (wave-dir mode)
make -C .cdd/waves/cell-runtime-doctrine/wave-validators all
```

`make all` builds the validator once and asserts the real-wave bijection PLUS the exact
per-fixture exit code (`go run` collapses every non-zero to 1; the built binary preserves them).

## Required result (PASS)

Real wave over the six `contracts/*.cn-cell-contract-v1.yaml` `acceptance.predicates` +
`oracle-registry.yaml` `assurance:` entries (mode: wave-dir; owners are the contracts' `cell.id`):

```
child acceptance predicates:  78
registry assurance entries:   78
mechanically-verifiable:      30
missing   (child -> registry): 0
phantom   (registry -> child): 0
child duplicates:              0
double-owned (registry dup):   0
mech missing/≠1 checker owner: 0
bijective: true
RESULT: PASS — oracle-ownership / classification bijection holds (exit 0)
```

## Fixture suite (exact exit codes — 0 PASS · 1 bijection/ownership defect · 2 fail-closed)

Positives (**exit 0**):

- `oracle-ownership.one-checker-each.positive.yaml` — block-style legal bijection.
- `oracle-ownership.flow-style.positive.yaml` — the SAME bijection in CUE-valid **flow** style;
  proves `cue export` normalizes flow ≡ block and does not break a positive.
- `wavedir.positive/` — minimal **wave-dir** mode: `cell.id ∈ wave nodes`, bijective.

Bijection-defect negatives (**exit 1**):

- `oracle-ownership.double-owned.negative.yaml` — a predicate owned twice.
- `oracle-ownership.flow-missing-registry.negative.yaml` — **flow-style** child predicate with no
  registry entry (the fixture-mode analog of the R11 false pass; R11 dropped it, R12 reports missing).
- `oracle-ownership.flow-mismatched-owner.negative.yaml` — flow-style registry `(owner,predicate)`
  matching no child (phantom + missing).
- `wavedir-flow-unregistered.negative/` — **wave-dir** dir-mode analog of the complete-wave mutation:
  a flow-style unregistered contract predicate → missing.

Fail-closed negatives (**exit 2**):

- `oracle-ownership.empty-owner.negative.yaml` — an empty semantic owner must not collapse to a match.
- `oracle-ownership.cellid-not-a-node.negative.yaml` — a child owner absent from `wave_nodes`
  (the single-file analog of `cell.id` ∉ wave nodes).

## Complete-wave mutation (the exact R11 false PASS, now rejected)

Reproducing external-β's demonstration on the real wave — replace WC-1's `acceptance.predicates` with the
CUE-valid flow form `predicates: ["wc1_actual_but_unregistered"]` and remove all 12 WC-1 registry
assurance entries — R11 reported `66 child predicates, 66 registry entries, bijective: true` (exit 0). R12:

```
child acceptance predicates:  67        # the flow predicate is now SEEN (was silently dropped)
registry assurance entries:   66
missing   (child -> registry): 1        # wc-1 :: wc1_actual_but_unregistered
bijective: false
RESULT: FAIL — oracle-ownership / classification bijection broken (exit 1)
```

## What the bijection proves

Every child acceptance predicate `(owner, predicate)` maps to **exactly one** registry
`assurance:` entry and vice-versa (78 ⇄ 78; no missing / phantom / duplicate), and every
`mechanically-verifiable` predicate binds **exactly one** concrete `checker|schema` owner.
This is a **whole-wave, cross-contract** property — owned by no child WC — so it is proven at
the **wave boundary, before any WC dispatches**. Inputs are normalized by `cue export`, so the
serialization style of the matter cannot hide a predicate from the gate.
