<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs] -->

# Self-coherence — Cycle #369

## §Gap

**Issue:** [#369 — Phase 2: Add draft CDD contract, receipt, and boundary decision schemas](https://github.com/usurobor/cnos/issues/369).

**Parent roadmap:** #366 (coherence-cell executability). This cycle is **Phase 2**.

**Predecessors:**

- #367 (Phase 1) — froze the parent-facing validator interface in `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`. Named `ValidationVerdict` vs `BoundaryDecision`, the override discipline, the δ-authoritative firing point, and the receipt-derivation rule. Schema syntax deferred to this cycle.
- #364 (closed) — `COHERENCE-CELL.md` doctrine. Frozen by both #367 and this cycle.

**Mode:** docs + schema. No code. No validator behavior change. `cn-cdd-verify` untouched (AC8).

**Version:** schemas pinned at v1 (`cnos.cdd.contract.v1`, `cnos.cdd.receipt.v1`, `cnos.cdd.boundary_decision.v1`). v1 is the **schema-artifact** version, not a CDD protocol version (AC2).

**Gap closed by this cycle:** #367 named the receptor; the types were prose, not declarative. Without `.cue`, Phase 3's validator has to invent shape mid-cycle. This cycle pins the shape so Phase 3 becomes a structural-reader implementation cycle.

**In-scope surfaces (per issue §Scope):**

- `schemas/cdd/contract.cue`
- `schemas/cdd/receipt.cue`
- `schemas/cdd/boundary_decision.cue`
- `schemas/cdd/README.md` (single-source prose; CUE files reference via header comment)
- `schemas/cdd/fixtures/valid-receipt.yaml`
- `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml`
- `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml`
- `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml`
- `.cdd/unreleased/369/*.md` (cycle evidence)

**Out of scope (per issue §Non-goals):** any change to `cn-cdd-verify/`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `delta/SKILL.md` (must not exist yet), `ROLES.md`, `CDD.md`, `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md` (sole permitted touch is an optional trailing-pointer; this cycle takes the safe path and omits it).

## §Skills

**Tier 1 (CDD lifecycle):**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm (consumed; not edited).
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface; load order, pre-review gate (14 rows), incremental self-coherence discipline, SHA-convention rules, polyglot re-audit rule.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority. Consulted to interpret AC boundaries.

**Tier 2:** none load-bearing for a docs+schema cycle with no shell/test toolchain authoring. (The repo's `cnos.eng/` Tier 2 bundles describe language-specific authoring; CUE-shape authoring is governed by the precedent of `schemas/skill.cue` per the issue's "Skills to load" section — "No new schema skill — the existing `schemas/skill.cue` is the precedent and form authority for CUE-shape patterns in this repo.")

**Tier 3 (issue-named):**

- `src/packages/cnos.core/skills/design/SKILL.md` — applied as authoring constraints (§3.2 single source of truth; §3.5 interfaces belong to consumers; §3.7 separate runtime surfaces; §3.10 prefer package/install cohesion).
- `src/packages/cnos.core/skills/write/SKILL.md` — applied to README prose: decisive over exhaustive; surface-containment statements internal to the document.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — used to interpret AC boundaries; no issue edits this cycle.

**Form authorities consulted (not edited):**

- `schemas/skill.cue` + `schemas/README.md` + `schemas/fixtures/skill-frontmatter/` — schema-layer precedent: CUE owns shape; prose lives in `schemas/{subsystem}/README.md`; fixture corpus lives under `schemas/{subsystem}/fixtures/`.
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — semantic owner of the interface this cycle types. Schemas reference it via header comment per active design constraint.
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — doctrine. Recursive scope-lift framing inherited from §Recursion Equation.

## §ACs

Per-AC evidence. Oracles run on HEAD of `cycle/369` at the implementation SHA `96ddf4f5` (the README commit, last implementation commit before this readiness signal — see [α/SKILL.md §2.6 SHA convention]).

### AC1 — Three CUE schemas exist at canonical paths and compile

| Oracle | Result |
|---|---|
| `test -f schemas/cdd/contract.cue` | PASS |
| `test -f schemas/cdd/receipt.cue` | PASS |
| `test -f schemas/cdd/boundary_decision.cue` | PASS |
| `cue vet -c=false schemas/cdd/contract.cue schemas/cdd/receipt.cue schemas/cdd/boundary_decision.cue` | exit 0 |

All three files at canonical paths under `schemas/cdd/` (mirrors `schemas/skill.cue` precedent — repo-root schema surface, not package-internal). All compile under schema-only mode. Files: [`schemas/cdd/contract.cue`](../../../schemas/cdd/contract.cue), [`schemas/cdd/receipt.cue`](../../../schemas/cdd/receipt.cue), [`schemas/cdd/boundary_decision.cue`](../../../schemas/cdd/boundary_decision.cue).

### AC2 — Schema IDs pinned at v1

| Oracle | Result |
|---|---|
| `rg 'cnos\.cdd\.contract\.v1' schemas/cdd/contract.cue` | 2 hits (header + reference) |
| `rg 'cnos\.cdd\.receipt\.v1' schemas/cdd/receipt.cue` | 1 hit (header) |
| `rg 'cnos\.cdd\.boundary_decision\.v1' schemas/cdd/boundary_decision.cue` | 1 hit (header) |

Each schema declares its versioned ID in the header comment. README states explicitly: "`v1` is the **schema-artifact** version, not a CDD protocol version."

### AC3 — Recursive scope-lift invariant typed in README + schemas

| Oracle | Result |
|---|---|
| `schemas/cdd/README.md` contains `## Scope-Lift Invariant` section | PASS — README line 38 |
| Section names all three projections, framed as projection-under-scope-lift | PASS — `### Projection 1/2/3` headings; closing subsection "Projection-under-scope-lift, not role-renaming" makes the framing explicit |
| `receipt.cue` references `boundary_decision.cue` (same-package reference) | PASS — `receipt.cue:53` references `#ValidationVerdict`; `:60` references `#BoundaryDecision`; `:65` references `#Transmissibility` |
| `receipt.cue` types `transmissibility: #Transmissibility` | PASS — `receipt.cue:65` |
| `boundary_decision.cue` defines `#Transmissibility` as a CUE-typed function of verdict × action | PASS — `#Transmissibility` enum defined at `boundary_decision.cue:86`; the derivation (verdict × action → enum) lives in `receipt.cue` lines 99–118 because both inputs are visible there |

The recursion is **stated** in `schemas/cdd/README.md` §Scope-Lift Invariant (three projections, each in its own subsection with explicit projection-under-scope-lift framing) and **structurally typed** in CUE (`#Transmissibility` enum + the if-chain enforcing the verdict × action → transmissibility table). Authors cannot set arbitrary transmissibility values — the table is computed, not asserted (see AC4).

### AC4 — `#Transmissibility` and `#BoundaryDecision` enforced structurally

| Oracle | Result |
|---|---|
| `rg 'accept\|release\|reject\|repair_dispatch\|override' schemas/cdd/boundary_decision.cue` returns all five values | PASS — `boundary_decision.cue:67` `action: "accept" \| "release" \| "reject" \| "repair_dispatch" \| "override"` |
| `rg '#Override' schemas/cdd/boundary_decision.cue` ≥ 1 hit | PASS — multiple hits; definition at `boundary_decision.cue:43` |
| Fixture `verdict: FAIL, action: accept, transmissibility: accepted` fails `cue vet` | PASS — verified via scratch fixture (`transmissibility: conflicting values "degraded" and "accepted"` + explicit `_|_` at `receipt.cue:116`). AC4 oracle hint "(AC7 covers this)" refers to the same structural rule that `invalid-override-masks-verdict.yaml` exercises in its (FAIL × override) form |
| Fixture `verdict: PASS, action: override` fails `cue vet` | PASS — verified via scratch fixture; `transmissibility` forced to `_|_` at `receipt.cue:110` |
| Fixture `action: override` and no `#Override` block fails `cue vet` | PASS — verified via scratch fixture; all `#Override` required fields surface as `incomplete value` errors |

Additional structural rejection (not in the AC4 oracle but enforced by the schema): an `override` block attached to `action: accept` (or any non-override action) fails vet because `boundary_decision.cue:77` has `override?: _|_` in the `if action != "override"` branch.

### AC5 — `#Receipt` required-field rule

| Oracle | Result |
|---|---|
| `boundary_decision` typed as `#BoundaryDecision` (non-null, required) | PASS — `receipt.cue:60` |
| `protocol_gap_count: int & >=0` required | PASS — `receipt.cue:71` |
| `protocol_gap_refs: [...#ProtocolGapRef]` required | PASS — `receipt.cue:72` |
| `#ProtocolGapRef` defined with `id`, `source`, `ref` | PASS — `receipt.cue:36–40` |
| Seven evidence refs typed `string` (required) | PASS — `evidence_root_ref`, `self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref` all at `receipt.cue:79–85`; `ci_refs?` optional at `:86` |
| Receipt omitting `boundary_decision` fails vet | PASS — covered by AC7 row 2 (`invalid-fail-no-boundary-decision.yaml`) |
| Receipt with `protocol_gap_count: 2, protocol_gap_refs: []` fails vet | PASS — verified via scratch fixture; `protocol_gap_count: conflicting values 0 and 2` (consistency constraint `protocol_gap_count: len(protocol_gap_refs)` at `receipt.cue:74` rejects the drift) |
| Receipt omitting any required evidence ref fails vet | Structural — every field at `receipt.cue:79–85` is unmarked (required); a fixture omitting any one surfaces "incomplete value string" with `-c`. Not separately enumerated as a fixture; covered by the schema requirement and the `cue vet -c` semantic |

### AC6 — One valid receipt fixture passes `cue vet`

Oracle:
```
cue vet -c -d '#Receipt' \
  schemas/cdd/contract.cue \
  schemas/cdd/boundary_decision.cue \
  schemas/cdd/receipt.cue \
  schemas/cdd/fixtures/valid-receipt.yaml
```

Result: exit 0. Fixture demonstrates `verdict: PASS`, `action: accept`, `transmissibility: accepted`, `protocol_gap_count: 0`, `protocol_gap_refs: []`, all seven evidence refs populated. Fixture: [`schemas/cdd/fixtures/valid-receipt.yaml`](../../../schemas/cdd/fixtures/valid-receipt.yaml).

### AC7 — Three doctrine-load-bearing invalid fixtures fail `cue vet`

| Fixture | Exits non-zero? | Failure reason mapped to load-bearing case |
|---|---|---|
| [`invalid-override-masks-verdict.yaml`](../../../schemas/cdd/fixtures/invalid-override-masks-verdict.yaml) | exit 1 | `transmissibility: conflicting values "degraded" and "accepted"` — the AC4 (FAIL × override) → degraded mapping rejects author-asserted `accepted`. Exercises #367 AC6 override-polarity: an override is degraded, not a way to mask a non-PASS verdict |
| [`invalid-fail-no-boundary-decision.yaml`](../../../schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml) | exit 1 | `boundary_decision: unresolved disjunction "accept" \| "release" \| "reject" \| "repair_dispatch" \| "override"` — the missing required `boundary_decision` surfaces through the `boundary_decision.action` reads in receipt.cue's if-chain. Exercises #367 AC3 δ-authoritative-ordering: no δ decision means the receipt is in the closure-emitted-but-not-accepted intermediate state, not a valid parent-facing receipt |
| [`invalid-gamma-preflight-authoritative.yaml`](../../../schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml) | exit 1 | Same error class as row 2 — `boundary_decision: unresolved disjunction …`. The `preflight_only: true` marker the fixture sets passes the open schema (extension keys are permitted) but cannot substitute for the required `boundary_decision`. Exercises #367 AC3 γ-preflight-non-authoritative clause: γ preflight cannot populate the receipt's δ-boundary block |

Failure reasons in this self-coherence map each invalid fixture to its named load-bearing case from #367's AC3/AC6.

### AC8 — `cn-cdd-verify` behavior unchanged

Oracle: `git diff origin/main..HEAD -- src/packages/cnos.cdd/commands/cdd-verify/`.

Result: empty. Zero changes under `cdd-verify/`. Phase 3 wires the validator; Phase 2 (this cycle) types the receptor only.

### AC9 — Surface containment

Oracle: `git diff origin/main..HEAD --stat`.

Result (9 files):

```
 .cdd/unreleased/369/self-coherence.md              |  57 ++++++
 schemas/cdd/README.md                              | 206 +++++++++++++++++++++
 schemas/cdd/boundary_decision.cue                  |  86 +++++++++
 schemas/cdd/contract.cue                           |  64 +++++++
 .../invalid-fail-no-boundary-decision.yaml         |  36 ++++
 .../invalid-gamma-preflight-authoritative.yaml     |  45 +++++
 .../fixtures/invalid-override-masks-verdict.yaml   |  58 ++++++
 schemas/cdd/fixtures/valid-receipt.yaml            |  35 ++++
 schemas/cdd/receipt.cue                            | 130 +++++++++++++
```

All changes land under `schemas/cdd/` (new directory: three `.cue` + README + 4 YAML fixtures = 8 files) or `.cdd/unreleased/369/` (cycle evidence = 1 file at signal time, more to follow as further self-coherence sections land). No prohibited file touched: `cdd-verify/` empty diff (AC8); `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `ROLES.md`, `CDD.md`, `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md` not modified (verified by absence from `--stat`). The optional `RECEIPT-VALIDATION.md` trailing-pointer exception was **not** exercised — this cycle takes the safe path and omits any touch on the doctrine document.

### AC10 — `cue vet` invocation pattern documented

`schemas/cdd/README.md` carries both invocations in multi-line literal form under `## How to run`:

- Schema-only compile (lines 132–138): `cue vet -c=false` with all three CUE files explicitly listed, no placeholders.
- Fixture validation (lines 146–152): `cue vet -c -d '#Receipt'` with all three CUE files plus `schemas/cdd/fixtures/{fixture}.yaml` — `{fixture}` is the only placeholder, marking the truly varying part per the issue's "literal multi-line form, placeholders only for the truly varying part" constraint.

Phase 3 inherits both invocations directly from the README (the active design constraint in the issue body names this as the authoritative documentation surface).
