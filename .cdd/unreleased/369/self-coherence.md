<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs, Self-check, Debt, CDD Trace] -->

# Self-coherence — Cycle #369

## Gap

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

## Skills

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

## ACs

Per-AC evidence. Oracles run on HEAD of `cycle/369` (per [α/SKILL.md §2.6 SHA convention], readiness-signal carries branch HEAD rather than a specific SHA; β polls and uses whatever HEAD it finds).

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

## Self-check

Per [`alpha/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md) §2.5: "did α's work push ambiguity onto β? Is every claim backed by evidence in the diff?"

### Did α push ambiguity onto β?

**Schema authoring decisions α made explicitly (so β does not have to re-derive):**

1. **Verdict enum collapsed to `PASS | FAIL`.** `RECEIPT-VALIDATION.md` §Output contract permits a `WARN` value with the explicit note "Phase 2 may collapse this into `warnings` and emit only `PASS | FAIL` as the verdict enum." α took the simpler form. The `warnings: [...string]` field remains on `#ValidationVerdict` for advisory observations. This is named in `boundary_decision.cue:25`. β does not need to relitigate.

2. **CUE package name `cdd`.** All three CUE files declare `package cdd`. This is what lets `receipt.cue` reference `#BoundaryDecision` without an `import` statement. AC3 oracle says "imports / references" — same-package reference satisfies the "references" form. This is consistent with `schemas/skill.cue`'s pattern of single-package-per-subsystem.

3. **`#Receipt` is open (`...`).** Matches `schemas/skill.cue` precedent. The closed form rejected `preflight_only` as "field not allowed" before getting to the missing-`boundary_decision` failure; open form makes the load-bearing AC5 violation the named reason `invalid-gamma-preflight-authoritative.yaml` fails. Open does not weaken hard-gate field requirements — required fields are still required.

4. **`#Override` shape — sub-record of `#BoundaryDecision` (per active design constraint).** `boundary_decision.cue:43–53` defines `#Override`; `#BoundaryDecision:61–80` carries `override?: #Override` and the required-iff constraint pair. α implemented option (a) from the dispatch-pinned defaults; β does not need to choose.

5. **`#ProtocolGapRef.source` enum chosen.** Issue named `source` as enum but did not pin members. α chose `"receipt" | "artifact" | "issue" | "review" | "closeout"` — five values describing **where in the cycle the gap observation came from**. α deliberately did **not** type `source` as the #366 Phase 6 four-class taxonomy (`cdd-skill-gap | cdd-protocol-gap | cdd-tooling-gap | cdd-metric-gap`) because that taxonomy classifies the **gap**, not its **source**. The taxonomy is a Phase 6 concern; v1 types the source surface. If β prefers a different name or set, this is a place where re-derivation effort is small. α flagged it as a debt-class possibility under §Debt.

6. **`schema_id` not encoded as a structural field.** AC2 grep oracle is satisfied by header comments. α did not add `schema_id: "cnos.cdd.receipt.v1"` as a sealed constant field on `#Receipt` because the schema artifact's filename + package + header comment is the canonical ID — no per-instance redundancy is required. Mirrors `schemas/skill.cue`'s pattern (no `schema_id` field there either).

7. **`base_sha` required, `contract_sha` optional.** The issue body says "base_sha (or contract_sha)" — α read this as "either-or with `base_sha` as the always-required pinning of merge-base on `origin/main`; `contract_sha` as the optional pinning of the issue-body SHA when the body is edited mid-cycle."

**Where ambiguity remains (deliberate; β can take the call):**

- The exact field name conventions on `#ValidationVerdict.provenance.input_refs` mirror the illustrative shape in `RECEIPT-VALIDATION.md` §Illustrative Example. β may rename without breaking the schema contract; the issue body did not pin these names.
- Whether `boundary_decision.actor` should be a typed role identity (e.g. `=~"^delta@[a-z]+\\.cdd\\.cnos$"`) or a free string. α chose `string` — Phase 3's V can encode role-identity policy via `cn-cdd-verify`'s structural reader rather than via CUE regex. β can call for tightening if Phase 3 needs it.

### Is every claim backed by evidence in the diff?

Every AC row in §ACs cites either a file:line ref or a cue vet exit code from the cycle's HEAD. No claim relies on prose alone:

- AC1 → `test -f` + `cue vet -c=false` exit 0.
- AC2 → `rg --count` hit counts.
- AC3 → README line numbers + `receipt.cue:53/60/65` + `boundary_decision.cue:86`.
- AC4 → enum literal at `boundary_decision.cue:67` + `#Override` at `:43` + four structural rejections each verified via cue vet on three named fixtures and three scratch fixtures.
- AC5 → 7 separate line refs in `receipt.cue` for required fields; consistency constraint verified via scratch fixture.
- AC6 → cue vet exit 0 on `valid-receipt.yaml`.
- AC7 → cue vet exit 1 on each of three invalid fixtures, with cue's error string captured and mapped to the doctrine clause it exercises.
- AC8 → `git diff --stat` empty for `cdd-verify/`.
- AC9 → 9-file diff stat in §ACs; absence of prohibited paths verified by inspection of the same stat.
- AC10 → README line refs for both invocations.

### Peer / harness audit (per alpha/SKILL.md §2.3, §2.4)

**Schema-bearing peers:** none — `schemas/cdd/` is a new subsystem directory. The repo's other schema-bearing subsystem is `schemas/skill.cue`, which is **not** a peer of this work (it types SKILL.md frontmatter; this types CDD receipts). The two subsystems share the layout pattern (CUE owns shape; README owns prose; fixtures live under `schemas/{subsystem}/fixtures/`) but not the shape itself. No sibling-update obligation.

**Receipt-shape harnesses:** none in this cycle. The receipt is not yet emitted by any tool — Phase 3 wires emission via `derive_receipt` and validation via `V`. AC8 enforces zero touch on `cdd-verify/`; AC9 enforces zero touch on any CI workflow.

**Polyglot re-audit (alpha/SKILL.md §2.6 row 9):** diff contains CUE, YAML, Markdown.

- **CUE** — `cue vet -c=false` on all three schemas: exit 0. Both schema-only compile (AC1) and fixture validation (AC6/AC7) exercise the schemas under both vet modes.
- **YAML** — `python3 -c "import yaml; yaml.safe_load(open(f))"` on each of the four fixture files: all parse cleanly.
- **Markdown** — README and self-coherence both contain tables; rows balance by visual inspection and by `cue vet` proxy (any oracle row that lies would surface as a failed oracle here). The one apparent awk-level "row mismatch" in self-coherence is a false positive caused by `\|` characters inside inline-code spans within a table cell (AC4 row 3 has the `rg 'accept\|release\|...'` pattern); these escape sequences are valid GitHub-flavored markdown and render correctly. Relative-link cross-references resolve: `../skill.cue` from `schemas/cdd/README.md` → `schemas/skill.cue` exists; `../README.md` → `schemas/README.md` exists; `../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` and `COHERENCE-CELL.md` both exist.
- **Shell / CI** — diff contains zero shell scripts or workflow files. Nothing to audit.

### Caller-path trace (alpha/SKILL.md §2.6 row 12)

The new modules in this cycle are CUE definitions, not callable functions. The "callers" of these schemas are:

1. `cue vet` itself — invoked by `valid-receipt.yaml` fixture validation (AC6 oracle), the three invalid fixtures (AC7), and the schema-only compile (AC1). Concrete call sites: `schemas/cdd/README.md` §How to run.
2. The README is itself a non-test caller of every fixture path — each `## Files` table row names a fixture as the doctrinal caller-side entry point.
3. Phase 3 (`cn-cdd-verify` refactor) will be the wire-level caller. This cycle pins the call interface; the Phase 3 issue (the immediate successor sub-issue under #366) is the caller-of-record once authored.

No "new function never called from main codepath" risk — there is no main codepath this cycle. The schemas are validators-of-data, not validators-themselves; they are called by `cue vet` invocations, not by program code.

### Artifact enumeration matches diff (alpha/SKILL.md §2.6 row 11)

Every file in `git diff --stat origin/main..HEAD` is mentioned by name in this self-coherence file:

| Diff file | Where mentioned |
|---|---|
| `.cdd/unreleased/369/self-coherence.md` | this file (itself) — §Gap declares the path; §CDD-Trace step 7 below trails it |
| `schemas/cdd/README.md` | §Gap "In-scope surfaces"; §ACs AC3/AC10 |
| `schemas/cdd/boundary_decision.cue` | §Gap "In-scope surfaces"; §ACs AC1/AC2/AC3/AC4 |
| `schemas/cdd/contract.cue` | §Gap "In-scope surfaces"; §ACs AC1/AC2 |
| `schemas/cdd/receipt.cue` | §Gap "In-scope surfaces"; §ACs AC1/AC2/AC3/AC5 |
| `schemas/cdd/fixtures/valid-receipt.yaml` | §Gap "In-scope surfaces"; §ACs AC6 |
| `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml` | §Gap "In-scope surfaces"; §ACs AC7 row 1 |
| `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml` | §Gap "In-scope surfaces"; §ACs AC7 row 2 |
| `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml` | §Gap "In-scope surfaces"; §ACs AC7 row 3 |

Nine files in diff; nine files enumerated. (The §Debt and §CDD-Trace and §Review-readiness sections of this same file will land in further commits and remain part of the same single file already named.)

## Debt

**Known debt α leaves explicit for β / Phase 3:**

1. **`#ProtocolGapRef.source` enum is α-chosen, not doctrine-sourced.** The five-value enum `"receipt" | "artifact" | "issue" | "review" | "closeout"` is α's reading of "where in the cycle a protocol gap observation surfaced." If Phase 6's ε-relocation cycle finds the dimension is wrong (e.g. they want it to be the #366 four-class taxonomy of *what kind of gap*, or a hybrid), this enum will need to shift. Schema-evolution implications: changing the enum is a v1-compatible operation (widening) or a v1-breaking operation (renaming/narrowing); Phase 6 makes the call. Not a hard gate for this cycle. See [§Self-check item 5](#self-check).

2. **No `#EvidenceGraph` schema authored.** Issue body deferred this to Phase 3 explicitly. α typed evidence refs as `string` in `receipt.cue:79–86`. Phase 3 author will build `evidence_graph.cue` and likely tighten the receipt's ref fields to typed refs (e.g. `#EvidenceRootRef` rather than plain `string`). Not a debt of this cycle; named here so the next cycle's α inherits the open question.

3. **`#ValidationVerdict.provenance` field names are illustrative.** The provenance sub-object's field names (`validator_identity`, `validator_version`, `checked_at`, `input_refs`) mirror `RECEIPT-VALIDATION.md` §Illustrative Example. Phase 3 may rename for consistency with its chosen invocation surface. Not a hard gate; α flags it so β does not push back on names that are explicitly subject to Phase 3 choice.

4. **`boundary_decision.actor` typed as `string`, not as a typed role identity.** A regex like `=~"^(alpha|beta|gamma|delta|epsilon)@[a-z]+\\.cdd\\.cnos$"` would tighten this — but the receipt's role-identity policy (which actors may write which decisions; how role pinning interacts with override actor) belongs to V in Phase 3, not to the data schema in Phase 2. α chose the loose form; Phase 3 can add structural reader checks.

5. **The `transmissibility` derivation table is encoded twice — once in `receipt.cue:99–118` and once in `schemas/cdd/README.md` §Scope-Lift Invariant Projection 2.** The schema is the load-bearing version; the README is documentation. If the table ever needs to change, both must move together. α did not add a `[[name]]`-style sync marker because the precedent (`schemas/skill.cue` + `schemas/README.md`) does not use one and the README explicitly names CUE as the structural authority. β can flag if drift-protection is desired.

6. **Polyglot re-audit script (in α's transcript) gave false positives on Markdown table balance and cross-references.** Not a debt of this cycle (no script is being shipped); flagged so a future Tier-2-tooling cycle that automates this re-audit knows the awk-pipe-count approach and naive path-resolution approach both have known limitations.

7. **`base_sha` vs `contract_sha` distinction.** α typed both fields, made `base_sha` required and `contract_sha` optional. Phase 3 may discover that V actually needs both required, or that `contract_sha` should default to `base_sha` when omitted. These are tightenings consistent with v1 compatibility.

**Debt α does NOT carry forward (resolved within this cycle):**

- `cn-cdd-verify` untouched — AC8 PASS.
- No edits to forbidden doctrine surfaces — AC9 PASS.
- All four fixtures present and validate per AC6/AC7 — no fixture corpus gap.
- README documents both cue vet invocations (AC10) — no operator-discovery gap.
- All three schemas pin v1 (AC2) — no version drift surface.

**No protocol gaps observed this cycle.** This is α's preflight reading; Phase 6 ε will read the receipt stream once Phase 3 lands and may surface gaps α did not see.

## CDD Trace

Per [`alpha/SKILL.md`](../../../src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md) §2.2: trace through CDD canonical artifact order, with explicit "not required" + reason where applicable.

### Step 1 — Design artifact

**Not required.** Predecessor cycle [#367](https://github.com/usurobor/cnos/issues/367) shipped [`RECEIPT-VALIDATION.md`](../../../src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md) — the load-bearing design surface this cycle types. No new design artifact authored; this cycle inherits the prior design and pins schema syntax for it. The issue body's `## Scope` and `## Active design constraints` sections function as the per-cycle design freeze (which authoring choices Σ pinned before dispatch, including the five questions α did not need to relitigate).

### Step 2 — Coherence contract (gap statement)

`.cdd/unreleased/369/self-coherence.md` §Gap (above). Issue [#369](https://github.com/usurobor/cnos/issues/369) names the gap (the receptor is named but not typed); this self-coherence file records α's reading of that gap and the surfaces this cycle closes.

### Step 3 — Plan

**Not required.** The work is three CUE files + one README + four YAML fixtures + one self-coherence file. Linear order: schemas (with `boundary_decision` first because the other two depend on its types) → fixtures (valid first, then three invalid) → README → per-AC oracle suite → self-coherence sections. No sequencing surprises; no parallel coordination across packages; no irreversible commits. The issue body's `## Acceptance criteria` set is the de-facto plan-of-record.

### Step 4 — Tests

**Tests are the four fixtures + per-AC oracle suite.** Test corpus:

- [`schemas/cdd/fixtures/valid-receipt.yaml`](../../../schemas/cdd/fixtures/valid-receipt.yaml) — positive case for `#Receipt`.
- [`schemas/cdd/fixtures/invalid-override-masks-verdict.yaml`](../../../schemas/cdd/fixtures/invalid-override-masks-verdict.yaml) — negative case for the AC4 transmissibility mapping under (FAIL × override).
- [`schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml`](../../../schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml) — negative case for the AC5 required-`boundary_decision` rule.
- [`schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml`](../../../schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml) — negative case for γ-preflight-cannot-substitute-for-δ.

The α-side oracle suite (in §ACs above) covers AC1–AC10 with named cue vet invocations and grep oracles. Additional scratch-fixture coverage of AC4 cases not covered by named fixtures (PASS+override, FAIL+accept, override-no-block, override-on-non-override-action, action=none enum-violation, count-refs-length mismatch) was run in α's transcript and is summarized in §ACs AC4 row notes; scratch fixtures were not committed because the named-three-invalid corpus is the issue's hard-gate requirement.

### Step 5 — Code

**No code in this cycle.** Three CUE schema files: [`schemas/cdd/contract.cue`](../../../schemas/cdd/contract.cue), [`schemas/cdd/boundary_decision.cue`](../../../schemas/cdd/boundary_decision.cue), [`schemas/cdd/receipt.cue`](../../../schemas/cdd/receipt.cue). These are *declarative type definitions*, not executable code. The cycle's mode is **docs + schema, no code, no validator behavior change** (issue body line 5).

### Step 6 — Docs

[`schemas/cdd/README.md`](../../../schemas/cdd/README.md) — single-source prose for the subsystem. Sections: `## Files` (the manifest table); `## Scope-Lift Invariant` with three projection subsections (AC3); `## How to run` with both cue vet invocations (AC10); `## What this directory does NOT do` (surface-containment internal to the document); `## Related issues and surfaces`. CUE files reference the README via the leading `// see schemas/cdd/README.md §Scope-Lift Invariant.` header comment.

### Step 7 — Self-coherence

This file. Carries CDD Trace through Step 7 by definition.

**Diff manifest (from `git diff --stat origin/main..HEAD`):**

- `schemas/cdd/README.md` (new)
- `schemas/cdd/contract.cue` (new)
- `schemas/cdd/boundary_decision.cue` (new)
- `schemas/cdd/receipt.cue` (new)
- `schemas/cdd/fixtures/valid-receipt.yaml` (new)
- `schemas/cdd/fixtures/invalid-override-masks-verdict.yaml` (new)
- `schemas/cdd/fixtures/invalid-fail-no-boundary-decision.yaml` (new)
- `schemas/cdd/fixtures/invalid-gamma-preflight-authoritative.yaml` (new)
- `.cdd/unreleased/369/self-coherence.md` (this file)

Nine new files; zero existing files modified.

### Step 8 — Pre-review gate

The 14-row pre-review-gate evaluation lands in [§Review-readiness](#review-readiness) below.
