---
cycle: 508
artifact: pass4-move-map.md
role: alpha
ac: AC3
---

# Pass 4A — Bundle-by-Bundle Move Map (AC3)

**Scope:** Maps every directory under `docs/alpha/`, `docs/beta/`, `docs/gamma/` to its reader-intent destination (or explicit "stays/defer" note). Per-bundle reference counts drawn from `pass4-path-inventory.txt` (AC1 total: 663).

**Important:** These are PROPOSED destinations for 4B–4E physical moves only. No files are moved in 4A.

---

## Bundle Discovery

Run used to enumerate real directories:
```
find docs/alpha docs/beta docs/gamma -maxdepth 1 -type d
```

Results: 32 bundle directories (14 under alpha, 6 under beta, 9 under gamma) plus 3 triad roots.

---

## docs/alpha/ Bundles (14 sub-bundles)

| Bundle | Ref Count (AC1) | Proposed Reader-Intent Destination | Pass | Notes |
|---|---|---|---|---|
| `docs/alpha/agent-runtime/` | 116 | `reference/runtime/` | 4D | High-ref count; has 4 version-stamped snapshots (frozen); non-snapshot content moves |
| `docs/alpha/architecture/` | 5 | `architecture/` | 4E | Merges into architecture tier |
| `docs/alpha/cli/` | 0 | `reference/cli/` | 4D | No refs in AC1 scope; safe to move when ready |
| `docs/alpha/cognitive-substrate/` | 14 | `architecture/cognitive-substrate/` | 4E | Architecture-tier content |
| `docs/alpha/ctb/` | 18 | `reference/ctb/` | 4D | Protocol reference material |
| `docs/alpha/design/` | 26 | stays/defer | — | Design docs with high inline citation count; needs separate triage. Not in 4B–4E scope per current issue. |
| `docs/alpha/doctrine/` | 63 | `concepts/doctrine/` | 4B | Low-risk reader surface; reader-intent destination is concepts tier |
| `docs/alpha/essays/` | 1 | `papers/` or `concepts/` | defer | 1 reference only; overlap with existing `docs/papers/` — needs operator decision |
| `docs/alpha/package-system/` | 41 | `reference/packages/` | 4D | Reference material; high src/ citation count |
| `docs/alpha/protocol/` | 19 | `reference/protocol/cn/` | 4D | Protocol reference |
| `docs/alpha/runtime-extensions/` | 8 | `reference/runtime/` (alongside agent-runtime) | 4D | 1 version-stamped snapshot (1.0.6/); non-snapshot content moves |
| `docs/alpha/schemas/` | 8 | `reference/schemas/` | 4D | Merges with `docs/beta/schema/` into unified schemas ref |
| `docs/alpha/security/` | 1 | `architecture/security/` | 4E | Architecture-tier content |
| `docs/alpha/vision/` | 3 | stays/defer | — | 3 references only; purpose unclear without deeper content review; operator decision needed |

**Alpha root-level files** (not in a sub-bundle): ~20 refs to `docs/alpha/DESIGN-CONSTRAINTS.md`, `docs/alpha/HUB-PLACEMENT-MODELS.md`, and similar. These need explicit triage — defer to operator.

---

## docs/beta/ Bundles (6 sub-bundles)

| Bundle | Ref Count (AC1) | Proposed Reader-Intent Destination | Pass | Notes |
|---|---|---|---|---|
| `docs/beta/architecture/` | 12 | `architecture/` | 4E | Has 1 version-stamped snapshot (3.14.4/); non-snapshot content moves; merges with alpha/architecture |
| `docs/beta/evidence/` | 1 | `evidence/` | 4B | Low-risk; 1 reference |
| `docs/beta/governance/` | 58 | stays/defer | — | High ref count (58); includes versioned snapshot 3.14.4/ with 25 refs; non-snapshot content (GLOSSARY.md, DOCUMENTATION-SYSTEM.md) is heavily cited. Needs dedicated sub-cell. Not assigned to 4B–4E without further operator review. |
| `docs/beta/guides/` | 2 | `guides/` | 4B | Low-risk; 2 references |
| `docs/beta/lineage/` | 8 | `concepts/lineage/` | 4B | Low-risk reader surface; has 1 version-stamped snapshot (3.14.4/) |
| `docs/beta/schema/` | 3 | `reference/schemas/` | 4D | Merges with `docs/alpha/schemas/` |

**Beta root-level files**: ~3 refs to `docs/beta/EXTENSION-REGISTRY.md`, `docs/beta/SUSTAINABILITY.md`. Defer to operator.

---

## docs/gamma/ Bundles (9 sub-bundles)

| Bundle | Ref Count (AC1) | Proposed Reader-Intent Destination | Pass | Notes |
|---|---|---|---|---|
| `docs/gamma/cdd/` | 385 | `development/cdd/` | 4C | **Highest ref count in inventory.** 65+ version-stamped snapshots (all frozen/historical — must NOT move); non-snapshot content (CDD.md, CDD-PACKAGE-AUDIT.md, RATIONALE.md) may move; `cn cdd-verify` tool hardcodes `docs/gamma/cdd` as default bundle path — **source/package-doctrine dependency blocks bare move without code change**. |
| `docs/gamma/checklists/` | 0 | `development/checklists/` | 4C | 0 refs in AC1 scope; safe to move |
| `docs/gamma/conventions/` | 29 | stays/defer | — | 29 refs; 3 are **generated/golden-bound** (both golden.yml files cite `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`); golden files must be re-rendered before this bundle moves. Do NOT move in 4B–4E without golden re-render. |
| `docs/gamma/design/` | 5 | `development/` (or `development/design/`) | 4C | Low ref count; development tier |
| `docs/gamma/essays/` | 13 | `papers/` or `concepts/` | defer | 13 refs; has 1 version-stamped snapshot (3.14.5/); overlaps with existing `docs/papers/` — operator decision needed |
| `docs/gamma/kata/` | 0 | `development/kata/` | 4C | 0 refs in AC1 scope; safe to move |
| `docs/gamma/plans/` | 35 | `development/plans/` | 4C | 35 refs; non-versioned; development tier |
| `docs/gamma/rules/` | 10 | `development/rules/` | 4C | 10 refs; has 1 version-stamped snapshot (3.14.5/); non-snapshot content moves |
| `docs/gamma/smoke/` | 2 | `development/smoke/` | 4C | 2 refs; development tier |

**Gamma root-level files**: ~19 refs to `docs/gamma/ENGINEERING-LANE-CLARITY.md`, `docs/gamma/ENGINEERING-LEVELS.md`, `docs/gamma/KATA-EVALUATION.md`. Defer to operator.

---

## Ref-Count Summary by Bundle

| Rank | Bundle | Refs | Class (dominant) |
|---|---|---|---|
| 1 | `docs/gamma/cdd/` | 385 | frozen/historical (versioned snapshots dominate) |
| 2 | `docs/alpha/agent-runtime/` | 116 | mixed (versioned snapshots + inline-path-citation) |
| 3 | `docs/beta/governance/` | 58 | frozen/historical + inline-path-citation |
| 4 | `docs/alpha/doctrine/` | 63 | inline-path-citation + markdown-link |
| 5 | `docs/alpha/package-system/` | 41 | inline-path-citation + source/package-doctrine |
| 6 | `docs/gamma/plans/` | 35 | inline-path-citation |
| 7 | `docs/alpha/design/` | 26 | inline-path-citation |
| 8 | `docs/gamma/conventions/` | 29 | **generated/golden-bound** (critical) + test-fixture |
| 9 | `docs/alpha/protocol/` | 19 | inline-path-citation + source/package-doctrine |
| 10 | `docs/alpha/ctb/` | 18 | source/package-doctrine + markdown-link |

---

## Golden-Bound Bundle Alert

**`docs/gamma/conventions/`** — DO NOT MOVE without golden re-render.

Both known `*.golden.yml` files reference `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`. Moving this file requires:
1. Updating the prompt source (`prompt.md`) for both `cds-dispatch` and `agent-admin` wake-providers
2. Re-rendering both golden files via `cn install-wake`
3. Updating the corresponding `.github/workflows/` files

This is a Pass 4C or later task, not 4B.

**`docs/gamma/cdd/`** — DO NOT MOVE `docs/gamma/cdd` path without updating `src/packages/cnos.cdd/commands/cdd-verify/run.go:59` (hardcoded default `"docs/gamma/cdd"`).

---

## Coverage Check

All 32 sub-bundle directories found by `find docs/alpha docs/beta docs/gamma -maxdepth 1 -type d` are mapped above (32 entries: 14 alpha + 6 beta + 9 gamma + 3 roots). No bundle is missing from the map. ✅
