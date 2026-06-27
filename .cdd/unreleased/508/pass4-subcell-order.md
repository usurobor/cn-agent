---
cycle: 508
artifact: pass4-subcell-order.md
role: alpha
ac: AC6
---

# Pass 4A — Sub-Cell Order (AC6)

**Purpose:** Document the recommended 4B→4C→4D→4E physical-move order with rationale, per-cell scope, and explicit reference to the AC4 golden-impact map for sequencing decisions.

**These moves do NOT happen in 4A.** This document is planning-only; all moves are deferred to 4B–4E cells after operator review of this audit.

---

## Sequencing Rationale Summary

The recommended order is **4B → 4C → 4D → 4E**, sequenced from lowest-risk/fewest-dependencies to highest-risk/most-dependencies.

| Pass | Scope | Risk | Blocker |
|---|---|---|---|
| 4B | Low-ref reader surfaces (guides, evidence, lineage, doctrine) | Low | None (confirmed by AC1: low ref counts, no golden-bound paths) |
| 4C | Development process bundle (gamma/cdd non-snapshots, rules, plans, etc.) | Medium | `cn cdd-verify` hardcoded path + snapshot hygiene |
| 4D | Reference bundles (protocol, agent-runtime, CLI, package-system, schemas) | Medium-High | High source/package-doctrine citation count; test-fixture dependencies |
| 4E | Architecture bundles (cognitive-substrate, security, architecture) | Medium | Cross-bundle merge complexity; versioned snapshots in beta/architecture |

**Key constraint from AC4 (golden-impact map):** `docs/gamma/conventions/` is the ONLY golden-bound bundle. It is NOT in 4B–4E as currently scoped — it requires a special pre-move golden-re-render step. The `conventions/` bundle move should be handled as a dedicated sub-task within 4C or as 4C-prefix, NOT in 4B (where "low-risk" is the bar).

---

## 4B — Low-Risk Reader Surfaces

### Scope

| From | To | Ref Count | Notes |
|---|---|---|---|
| `docs/beta/guides/` | `guides/` | 2 | Merges with existing `docs/guides/` reader-intent dir |
| `docs/beta/evidence/` | `evidence/` | 1 | Merges with existing `docs/evidence/` reader-intent dir |
| `docs/beta/lineage/` | `concepts/lineage/` | 8 | Has 1 frozen snapshot (3.14.4/); non-snapshot files only |
| `docs/alpha/doctrine/` | `concepts/doctrine/` | 63 | No versioned snapshots; highest-ref 4B bundle; all are markdown-link or inline-path-citation |

### Rationale

These four bundles have the following properties that make them lowest-risk:
1. **No golden-bound references** — AC4 impact map confirms zero golden-bound paths in any of these bundles.
2. **No source/package-doctrine hardcodes** — `src/` files do not hardcode these paths as code-level defaults (only as doc comments, which can be updated alongside the move).
3. **Low test-fixture exposure** — `test/` and `.github/` files do not check these specific paths.
4. **Modest ref counts** — `beta/guides/` (2), `beta/evidence/` (1), `beta/lineage/` (8) are the lowest-count bundles in the inventory; `alpha/doctrine/` (63) is higher but has no frozen snapshots.

**Caveat for `docs/alpha/doctrine/` (63 refs):** While there are no golden-bound or code-hardcoded paths in doctrine, there are 6 `markdown-link` references using GitHub absolute URLs from other doctrine files (e.g., ethics, coherence, inheritance essays cross-referencing each other). These links will break if the doctrine bundle moves without updating the link targets. The move cell must update these cross-references as part of the physical move.

**`docs/beta/lineage/3.14.4/`** (frozen snapshot): Stays in place when `docs/beta/lineage/` moves. The snapshot directory at `docs/beta/lineage/3.14.4/` is in the do-not-move list (Section 2 of AC5); only the non-snapshot files move.

### Explicitly excluded from 4B

- `docs/gamma/conventions/` — golden-bound per AC4 impact map; must NOT be in 4B.
- `docs/beta/governance/` — high ref count (58), versioned snapshot, heavily cited; not low-risk.

---

## 4C — Development Process Bundle

### Scope

| From | To | Ref Count | Notes |
|---|---|---|---|
| `docs/gamma/cdd/` (non-snapshot) | `development/cdd/` | ~385 total* | Snapshots frozen; only CDD.md, CDD-PACKAGE-AUDIT.md, RATIONALE.md, etc. move |
| `docs/gamma/rules/` (non-snapshot) | `development/rules/` | 10 | Has 1 frozen snapshot (3.14.5/) |
| `docs/gamma/plans/` | `development/plans/` | 35 | No versioned snapshots |
| `docs/gamma/checklists/` | `development/checklists/` | 0 | No refs; safe to move |
| `docs/gamma/smoke/` | `development/smoke/` | 2 | Low ref count |
| `docs/gamma/kata/` | `development/kata/` | 0 | No refs; safe to move |
| `docs/gamma/design/` | `development/` or `development/design/` | 5 | Small bundle; development tier |
| `docs/gamma/conventions/` | (TBD after golden re-render) | 29 | **Golden-bound — requires pre-move golden re-render** |
| `docs/gamma/essays/` (non-snapshot) | defer | 13 | Overlaps `docs/papers/`; operator decision |

*Note: The 385 refs for `docs/gamma/cdd/` include 352 from versioned snapshot files (which stay in place). The actual non-snapshot cdd ref count is 385 - (versioned snapshot refs) ≈ 33 non-frozen refs.

### Rationale

4C comes before 4D and 4E because:
1. **Lower src/ coupling than 4D/4E bundles** — `docs/gamma/cdd/` is the primary 4C bundle; its src/ dependency is a single hardcoded default in `run.go` (manageable within the move cell).
2. **No golden-bound paths in non-conventions bundles** — only `conventions/` is golden-bound; the rest of the gamma bundle can move independently.
3. **Development-process content has lower external link surface** — gamma bundles are primarily linked from other gamma docs and `.cdd/` records; `.cdd/` records are frozen and will not be updated.

### Key constraint from AC4 golden-impact map

**`docs/gamma/conventions/`** is the ONLY bundle in the entire 4A inventory that is golden-bound. Per `pass4-golden-impact-map.md`:
- Both `cnos-cds-dispatch.golden.yml` and `cnos-agent-admin.golden.yml` reference `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`.
- Moving this file triggers a CI failure in `install-wake-golden.yml` unless both prompt sources and golden fixtures are updated first.
- The move sequence must be: update `prompt.md` sources → re-render goldens → commit → then move `AGENT-ACTIVATION-LOG-v0.md`.
- This makes `docs/gamma/conventions/` the **most operationally complex move** in 4C.

### Key constraint from AC1 source/package-doctrine

**`docs/gamma/cdd` hardcoded path** in `src/packages/cnos.cdd/commands/cdd-verify/run.go:59` must be updated simultaneously with (or before) the physical move of `docs/gamma/cdd/` non-snapshot content. The 4C cell for `gamma/cdd/` MUST include a code change to `run.go`, `README.md`, and `test-cn-cdd-verify.sh`.

---

## 4D — Reference Bundles

### Scope

| From | To | Ref Count | Notes |
|---|---|---|---|
| `docs/alpha/protocol/` | `reference/protocol/cn/` | 19 | Protocol reference |
| `docs/alpha/agent-runtime/` (non-snapshot) | `reference/runtime/` | 116 total* | 4 frozen snapshots stay; high src/ citation count |
| `docs/alpha/cli/` | `reference/cli/` | 0 | No refs; safe to move |
| `docs/alpha/package-system/` | `reference/packages/` | 41 | High src/ citation; `pkg.go` comment references |
| `docs/alpha/ctb/` | `reference/ctb/` | 18 | Used in schemas/README.md; `test/cmd/` references |
| `docs/alpha/schemas/` + `docs/beta/schema/` | `reference/schemas/` | 8 + 3 = 11 | Merge required; `build.yml` checks `docs/alpha/schemas/protocol-contract.json` |
| `docs/alpha/runtime-extensions/` (non-snapshot) | `reference/runtime/` | 8 | 1 frozen snapshot (1.0.6/) stays |

*Note: `docs/alpha/agent-runtime/` 116 refs include versioned snapshot refs (frozen). The 4 snapshot directories (3.7.0, 3.8.0, 3.10.0, 3.14.0) stay in place.

### Rationale

4D comes after 4B and 4C because:
1. **Highest source/package-doctrine dependency** — `src/go/` and `src/ocaml/` code files reference `docs/alpha/agent-runtime/`, `docs/alpha/package-system/`, and `docs/alpha/protocol/` as design authorities. These comments must be updated alongside the move, requiring code-change coordination.
2. **Test-fixture dependencies** — `test/cmd/cn_contract_test.ml`, `cn_sandbox_test.ml`, and `cn_traceability_test.ml` reference `docs/alpha/schemas/protocol-contract.json` and `docs/alpha/`. The CI `build.yml` workflow runs a `diff` against `docs/alpha/schemas/protocol-contract.json`. Moving `docs/alpha/schemas/` requires updating the build.yml diff command.
3. **High ref counts** — `docs/alpha/agent-runtime/` (116), `docs/alpha/package-system/` (41) are the highest-ref non-frozen bundles; more citation repointing work.

**Key constraint: `docs/alpha/schemas/protocol-contract.json`**
The file `docs/alpha/schemas/protocol-contract.json` is explicitly checked by:
- `.github/workflows/build.yml:223`: `diff docs/alpha/schemas/protocol-contract.json test/cmd/protocol-contract.json`
- `test/cmd/cn_contract_test.ml:3`: "Verifies that code constants match docs/alpha/schemas/protocol-contract.json"
- `test/cmd/cn_traceability_test.ml:5`: references the same file

Moving `docs/alpha/schemas/` requires simultaneous updates to `build.yml` and the OCaml test files.

---

## 4E — Architecture Bundles

### Scope

| From | To | Ref Count | Notes |
|---|---|---|---|
| `docs/alpha/cognitive-substrate/` | `architecture/cognitive-substrate/` | 14 | Merges into `docs/architecture/` |
| `docs/alpha/security/` | `architecture/security/` | 1 | Low ref count |
| `docs/beta/architecture/` (non-snapshot) | `architecture/` | 12 | 1 frozen snapshot (3.14.4/) stays; merges with `docs/alpha/architecture/` (5 refs) |
| `docs/alpha/architecture/` | `architecture/` | 5 | Merges with beta/architecture |

### Rationale

4E comes last because:
1. **Merge complexity** — `alpha/architecture/` and `beta/architecture/` must merge into a unified `docs/architecture/` tier; this requires content deduplication and may require operator decisions on naming conflicts.
2. **Lowest urgency** — architecture content has lower operational coupling than reference bundles (fewer hardcoded path dependencies) but more semantic complexity.
3. **Versioned snapshot in beta/architecture** — `docs/beta/architecture/3.14.4/` stays frozen; only the non-snapshot ARCHITECTURE.md and similar files move.
4. **Existing `docs/architecture/` directory** — The reader-intent portal already has a `docs/architecture/` dir; the 4E move is a merge, not a new directory creation.

---

## Golden-Impact Map Reference Summary

Per `pass4-golden-impact-map.md` (AC4):

| Bundle | Golden-Bound? | Move Pass | Sequencing Impact |
|---|---|---|---|
| `docs/gamma/conventions/` | **YES** — both golden.yml files | 4C | Requires pre-move golden re-render; most complex single-file move |
| All other bundles | NO | 4B/4C/4D/4E | No golden re-render required for the move itself |

The only riskiest-because-golden-bound bundle is `docs/gamma/conventions/`. It is sequenced in 4C (not 4B) because of this golden-bound complexity, consistent with the AC4 impact map's finding. All 4B bundles have been verified against the AC4 impact map to confirm zero golden-bound dependencies.

---

## Constraints Summary

| Constraint | Affected Bundle | Blocking Pass | Resolution |
|---|---|---|---|
| Golden-bound path (`AGENT-ACTIVATION-LOG-v0.md`) | `docs/gamma/conventions/` | 4C | Pre-move: update both prompt.md + re-render both goldens |
| Code-hardcoded path (`docs/gamma/cdd` in `run.go`) | `docs/gamma/cdd/` | 4C | Simultaneous code change in `cdd-verify/run.go` |
| CI path-check (`docs/alpha/schemas/protocol-contract.json` in `build.yml`) | `docs/alpha/schemas/` | 4D | Simultaneous update to `build.yml` + test files |
| Test-fixture paths (`docs/alpha/` in OCaml tests) | `docs/alpha/schemas/`, `docs/alpha/` generally | 4D | Update `test/cmd/*.ml` files alongside move |
| Frozen snapshot directories (70+ dirs) | Multiple bundles | All | Snapshots STAY; only non-snapshot files move in each pass |

---

## AC6 Certification

- The ordering explicitly cites `pass4-golden-impact-map.md` (AC4): ✅
- The golden-bound bundle (`docs/gamma/conventions/`) is NOT placed in 4B: ✅ (placed in 4C with explicit pre-move re-render gate)
- Riskiest bundles (`gamma/cdd/`, `alpha/agent-runtime/`, `alpha/schemas/`) are in 4C/4D, not 4B: ✅
- All 4B bundles verified against AC4 impact map as zero-golden-bound: ✅
