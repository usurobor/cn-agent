<!-- wave-revision: R11 -->
# Wave-boundary pre-authorization evidence — oracle-ownership bijection (cnos#671 R11)

This is the **result/evidence binding** referenced by `wave.cn-wave-v1.yaml`
`gates.wave_authorization.preauthorization_gates[]` and by `oracle-registry.yaml`
`wave_predicates.wave_oracle_ownership_bijection_enforced`. It records the **PASS** of the
materialized wave-boundary validator, **content-bound to wave revision R11**.

The wave authorization is **NOT authorization-ready** unless this validator resolves at its
pinned `validator_sha256` AND its invocation exits 0 with `bijective: true`. Removing or
corrupting `oracle_ownership_bijection.go` breaks the pinned hash and the gate → the wave
returns to the pre-authorization hold (see `wave.stop_conditions.revoked_authorization`).

## Bound artifact

- **validator:** `.cdd/waves/cell-runtime-doctrine/wave-validators/oracle_ownership_bijection.go`
- **validator_sha256:** `271887be7cab33d592e009c5a8c04bfa20ed09a31fe5017b0e19e7065be7c832`
- **positive_fixture:** `fixtures/oracle-ownership.one-checker-each.positive.yaml`
  (`sha256:fca4150023a27c9043d999f29b07cf1592319a63e276c773a29834f7058f78fb`)
- **negative_fixture:** `fixtures/oracle-ownership.double-owned.negative.yaml`
  (`sha256:9943f6e4a58e338f049c1cc470f455da62b6ba18a9ea08a289cc150037b3a9ca`)
- **authorized_revision:** `R11`

## Invocation (credential-free; no modules, no network)

```
go run oracle_ownership_bijection.go ..                 # real six contracts + registry
make -C .cdd/waves/cell-runtime-doctrine/wave-validators all
```

## Required result (PASS)

Real wave over the six `contracts/*.cn-cell-contract-v1.yaml` `acceptance.predicates` +
`oracle-registry.yaml` `assurance:` entries:

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

Fixture behavior (proves the validator discriminates):

- `oracle-ownership.one-checker-each.positive.yaml` → `bijective: true`, **exit 0**.
- `oracle-ownership.double-owned.negative.yaml` → `double-owned (registry dup): 1`, **exit 1** (non-zero).

## What the bijection proves

Every child acceptance predicate `(owner, predicate)` maps to **exactly one** registry
`assurance:` entry and vice-versa (78 ⇄ 78; no missing / phantom / duplicate), and every
`mechanically-verifiable` predicate binds **exactly one** concrete `checker|schema` owner.
This is a **whole-wave, cross-contract** property — owned by no child WC — so it is proven at
the **wave boundary, before any WC dispatches**.
