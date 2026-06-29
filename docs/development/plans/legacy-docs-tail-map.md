# Legacy docs tail map (H4A audit)

**Purpose:** classify everything still living under `docs/alpha/`, `docs/beta/`,
`docs/gamma/`, `docs/essays/` after the 4A–4E reader-intent migration, so the
H4B/H4C/H5 cleanup is a known, bounded move — not another 4D-style surprise.

**Audit (this commit):** `docs/{alpha,beta,gamma,essays}` =
**41 real-active** · **13 frozen-historical** (date-based PRAs) ·
**192 frozen snapshots** (`X.Y.Z/`) · **149 redirect stubs** (`# Moved`) ·
**1 golden-bound**. Only the **41 real-active** + the **1 golden-bound** need
action; everything else stays in place (frozen / stub).

## Rule

No active content should live under `alpha/`, `beta/`, `gamma/`, `essays/`. Only
frozen history (snapshots, PRAs), redirect stubs, or explicitly golden-bound
material may remain. Active links must point to real homes, not stubs.

## H4B/H4C worklist — real active content (40 files) + proposed homes

### `docs/gamma/kata/` → `docs/development/kata/` (23 files) — active_deferred
Kata packets, templates, and run records:
`A1-open-op-registry-conflict.md`, `A2-extension-registry-engine-compatibility.md`,
`B1-runtime-contract-v2-parity-review.md`, `C1-capability-growth-boundary.md`,
`C2-browser-capability-and-ecosystem-boundary.md`,
`TEMPLATE-KATA-PACKET.md`, `TEMPLATE-RUN-RECORD.md`, `TEMPLATE-SCORE-SHEET.md`,
and `runs/` (15: A1/B1/C1 × {cold,selected}-001 + score sheets + `*-DELTA.md`).
`docs/development/README.md` already points at `kata/`.

### `docs/beta/governance/` → `docs/reference/governance/` (4) — active_deferred
`DOCUMENTATION-SYSTEM.md`, `GLOSSARY.md`, `NAMING.md`, `README.md`.

### `docs/gamma/design/` → `docs/development/design/` (2) — active_deferred
`README.md`, `ccnf-o-track-a1-survey.md`.

### `docs/gamma/` top-level engineering docs → `docs/development/` (3) — active_deferred
`ENGINEERING-LANE-CLARITY.md`, `ENGINEERING-LEVELS.md` → `docs/development/`;
`KATA-EVALUATION.md` → `docs/development/kata/` (with the kata).

### `docs/gamma/smoke/` → `docs/evidence/smoke/` (1) — active_deferred
`cds-dispatch-smoke-20260623.md`.

### `docs/alpha/` actives (5) — needs destination decision
| File | Proposed home | Note |
|------|---------------|------|
| `DESIGN-CONSTRAINTS.md` | `docs/reference/` or `docs/development/rules/` | actively linked from `docs/architecture/README.md` + OCaml legacy doc |
| `HUB-PLACEMENT-MODELS.md` | `docs/reference/` or `docs/architecture/` | linked from `docs/architecture/README.md` |
| `architecture/INVARIANTS.md` | `docs/architecture/` or `docs/development/rules/` | linked from `docs/architecture/README.md` |
| `design/WRITER-PACKAGE.md` | `docs/reference/packages/` or `docs/architecture/design/` | |
| `vision/AGENT-NETWORK.md` | `docs/concepts/vision/` or `docs/concepts/` | |

### `docs/beta/` actives (3) — needs destination decision
| File | Proposed home | Note |
|------|---------------|------|
| `EXTENSION-REGISTRY.md` | `docs/reference/runtime/` | |
| `SUSTAINABILITY.md` | root `SUSTAINABILITY.md` or `docs/reference/governance/` | linked from `README.md` Support section |
| `architecture/README.md` | **retire** (it is the 4E redirect index, `# Architecture — Moved`) | retire once `beta/architecture/` has no inbound links beyond the frozen `3.14.4/` |

## H5 — golden-bound (separate high-risk cell)

`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` → `docs/reference/conventions/`.
**Do not move casually.** Proven golden-bound (referenced by the wake goldens
`cnos-cds-dispatch.golden.yml` / `cnos-agent-admin.golden.yml`). Requires its own
cell: wake-golden re-render + golden review + CI proof + explicit receipt.

## Stay in place (no action)

- **frozen snapshots** (192): `X.Y.Z/` version dirs under any legacy bundle — immutable.
- **frozen-historical PRAs** (13): `docs/gamma/cdd/docs/<date>/…/POST-RELEASE-ASSESSMENT*.md` — date-based release evidence, immutable.
- **redirect stubs** (149): `# Moved` files at old paths — retire later (H4 tail) only once no active inbound links resolve through them.
- **`docs/essays/`**: redirect stub(s) only (e.g. `agent-first.md`) — retire candidate once inbound links clear.

## Sequencing note

H4B/H4C are pure active-link-preserving moves (stub at old path, repoint active
links) — same mechanics as 4B–4E, now de-risked by this map. H5 is separate
(golden). Retiring stubs (149) and `essays/` is the final tail, after a grace
period. `gamma/conventions/` stays until H5.
