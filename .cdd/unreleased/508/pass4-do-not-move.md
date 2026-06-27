---
cycle: 508
artifact: pass4-do-not-move.md
role: alpha
ac: AC5
---

# Pass 4A — Do-Not-Move List (AC5)

**Purpose:** Enumerate all paths that must NOT be moved, renamed, or deleted during Pass 4B–4E physical moves. Entries are named (not guessed) based on `find` output and AC4 golden-impact analysis.

**Exclusion reasons:**
- **frozen/historical** — `.cdd/` records, version-stamped snapshot directories, PRA/audit docs
- **golden-bound** — paths cited in `*.golden.yml` files (see AC4 impact map)
- **source/package-doctrine hardcoded** — paths hardcoded in `src/` as a default or code constant

---

## Section 1 — `.cdd/` Records (frozen/historical)

All files and directories under `.cdd/` are frozen historical records of the CDD framework. They must not be moved, renamed, or deleted. This includes:

### `.cdd/releases/` (version-stamped release records)
```
.cdd/releases/3.55.0/
.cdd/releases/3.56.0/
.cdd/releases/3.56.1/
.cdd/releases/3.56.2/
.cdd/releases/3.57.0/
.cdd/releases/3.58.0/
.cdd/releases/3.59.0/
.cdd/releases/3.60.0/
.cdd/releases/3.61.0/
.cdd/releases/3.62.0/
.cdd/releases/3.63.0/
.cdd/releases/3.65.0/
.cdd/releases/3.66.0/
.cdd/releases/3.67.0/
.cdd/releases/3.68.0/
.cdd/releases/3.69.0/
.cdd/releases/3.70.0/
.cdd/releases/3.71.0/
.cdd/releases/3.72.0/
.cdd/releases/3.73.0/
.cdd/releases/3.74.0/
.cdd/releases/3.75.0/
.cdd/releases/3.76.0/
.cdd/releases/3.77.0/
.cdd/releases/3.78.0/
.cdd/releases/3.79.0/
.cdd/releases/3.80.0/
.cdd/releases/3.81.0/
.cdd/releases/3.82.0/
.cdd/releases/docs/
```

### `.cdd/unreleased/` (open-cycle records, all cycles)
```
.cdd/unreleased/425/
.cdd/unreleased/426/
.cdd/unreleased/427/
.cdd/unreleased/428/
.cdd/unreleased/429/
.cdd/unreleased/430/
.cdd/unreleased/485/
.cdd/unreleased/486/
.cdd/unreleased/487/
.cdd/unreleased/491/
.cdd/unreleased/496/
.cdd/unreleased/497/
.cdd/unreleased/500/
.cdd/unreleased/506/
.cdd/unreleased/508/  ← this cycle's artifact root
```

### `.cdd/` root config and control files
```
.cdd/CADENCE
.cdd/CDD-VERSION
.cdd/DISPATCH
.cdd/OPERATORS
.cdd/legacy-exceptions.yml
```

### `.cdd/iterations/` (cross-repo iteration records)
```
.cdd/iterations/INDEX.md
.cdd/iterations/cross-repo/ (all subdirectories)
.cdd/iterations/wave-2026-05-12.md
.cdd/iterations/wave-2026-05-19-protocol-hygiene.md
```

### `.cdd/waves/` (wave records)
```
.cdd/waves/2026-05-17-coherence-cell/
.cdd/waves/2026-05-19-protocol-hygiene/
.cdd/waves/breadth-2026-05-12/
.cdd/waves/completeness-2026-05-12/
.cdd/waves/hardening-2026-05-12/
```

### `.cdd/MCAs/` and `.cdd/proposals/`
```
.cdd/MCAs/INDEX.md
.cdd/proposals/gamma-iteration-after-3.61.0.md
```

---

## Section 2 — Version-Stamped Snapshot Directories (frozen/historical)

All directories at path pattern `docs/{triad}/{bundle}/{X.Y.Z}/` are frozen historical snapshots. They must not be moved even when the containing bundle moves — the snapshot directories stay where they are.

**Discovered by:** `find docs -type d -name '[0-9]*.[0-9]*.[0-9]*'`

### `docs/alpha/` snapshots
```
docs/alpha/agent-runtime/3.7.0/
docs/alpha/agent-runtime/3.8.0/
docs/alpha/agent-runtime/3.10.0/
docs/alpha/agent-runtime/3.14.0/
docs/alpha/runtime-extensions/1.0.6/
```

### `docs/beta/` snapshots
```
docs/beta/architecture/3.14.4/
docs/beta/governance/3.14.4/
docs/beta/lineage/3.14.4/
docs/beta/schema/3.14.4/
```

### `docs/gamma/cdd/` snapshots (67 directories)
```
docs/gamma/cdd/3.13.0/
docs/gamma/cdd/3.14.1/
docs/gamma/cdd/3.14.4/
docs/gamma/cdd/3.14.5/
docs/gamma/cdd/3.14.6/
docs/gamma/cdd/3.14.7/
docs/gamma/cdd/3.15.0/
docs/gamma/cdd/3.15.2/
docs/gamma/cdd/3.16.0/
docs/gamma/cdd/3.16.1/
docs/gamma/cdd/3.16.2/
docs/gamma/cdd/3.17.0/
docs/gamma/cdd/3.18.0/
docs/gamma/cdd/3.19.0/
docs/gamma/cdd/3.20.0/
docs/gamma/cdd/3.22.0/
docs/gamma/cdd/3.24.0/
docs/gamma/cdd/3.25.0/
docs/gamma/cdd/3.26.0/
docs/gamma/cdd/3.28.0/
docs/gamma/cdd/3.29.0/
docs/gamma/cdd/3.31.0/
docs/gamma/cdd/3.32.0/
docs/gamma/cdd/3.33.0/
docs/gamma/cdd/3.34.0/
docs/gamma/cdd/3.35.0/
docs/gamma/cdd/3.36.0/
docs/gamma/cdd/3.37.0/
docs/gamma/cdd/3.38.0/
docs/gamma/cdd/3.39.0/
docs/gamma/cdd/3.40.0/
docs/gamma/cdd/3.41.0/
docs/gamma/cdd/3.42.0/
docs/gamma/cdd/3.43.0/
docs/gamma/cdd/3.44.0/
docs/gamma/cdd/3.45.0/
docs/gamma/cdd/3.46.0/
docs/gamma/cdd/3.47.0/
docs/gamma/cdd/3.48.0/
docs/gamma/cdd/3.49.0/
docs/gamma/cdd/3.50.0/
docs/gamma/cdd/3.51.0/
docs/gamma/cdd/3.52.0/
docs/gamma/cdd/3.53.0/
docs/gamma/cdd/3.54.0/
docs/gamma/cdd/3.55.0/
docs/gamma/cdd/3.56.0/
docs/gamma/cdd/3.56.1/
docs/gamma/cdd/3.56.2/
docs/gamma/cdd/3.57.0/
docs/gamma/cdd/3.58.0/
docs/gamma/cdd/3.59.0/
docs/gamma/cdd/3.60.0/
docs/gamma/cdd/3.61.0/
docs/gamma/cdd/3.62.0/
docs/gamma/cdd/3.63.0/
docs/gamma/cdd/3.78.0/
docs/gamma/cdd/3.79.0/
docs/gamma/cdd/3.81.0/
```

### `docs/gamma/essays/` snapshots
```
docs/gamma/essays/3.14.5/
```

### `docs/gamma/rules/` snapshots
```
docs/gamma/rules/3.14.5/
```

---

## Section 3 — Golden-Bound Paths (AC4 impact map)

The following path is cited by both known `*.golden.yml` files and must not move until those golden files are updated and re-rendered:

```
docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md
```

**Specifically blocked by:**
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (line 225)
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (lines 82, 168)

**Unblock sequence (required before moving `docs/gamma/conventions/`):**
1. Update `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` with new path
2. Update `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` with new path
3. Re-render: `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`
4. Re-render: `cn install-wake agent-admin --out .github/workflows/cnos-agent-admin.yml`
5. Update both `*.golden.yml` files: `cn install-wake cds-dispatch` + `cn install-wake agent-admin`
6. Verify `install-wake-golden.yml` CI passes

Until steps 1–6 are complete, the entire `docs/gamma/conventions/` bundle is on the do-not-move list.

---

## Section 4 — Source/Package-Doctrine Hardcoded Path

```
docs/gamma/cdd
```

**Hardcoded in:** `src/packages/cnos.cdd/commands/cdd-verify/run.go:59`
```go
a.Bundle = "docs/gamma/cdd" // default per bash predecessor
```

Moving `docs/gamma/cdd/` requires a code change to `run.go` (and corresponding update to `README.md` + `test-cn-cdd-verify.sh`) before or simultaneously with the physical move. This is a non-trivial dependency: the `cn cdd-verify` CLI command will break if `docs/gamma/cdd/` is relocated without updating this default.

---

## Section 5 — Historical PRA and Audit Docs

The following non-versioned docs are historical audit records that should not be moved without careful operator review (their triad paths are cited extensively by other docs):

- `docs/gamma/cdd/CDD-PACKAGE-AUDIT.md` — 27 references in AC1 inventory; a historical audit document
- `docs/alpha/design/WRITER-PACKAGE.md` — 26 references; design history document

These are NOT versioned-snapshot paths and thus not in Section 2, but their high reference count and audit-record nature make them candidates for careful triage before any move.

---

## Summary

| Category | Count | Status |
|---|---|---|
| `.cdd/` root files | 5 | DO NOT MOVE |
| `.cdd/releases/` directories | 29 | DO NOT MOVE |
| `.cdd/unreleased/` directories | 15+ | DO NOT MOVE |
| `.cdd/iterations/` records | many | DO NOT MOVE |
| `.cdd/waves/` records | 5+ | DO NOT MOVE |
| `docs/alpha/agent-runtime/` snapshots | 4 dirs | DO NOT MOVE |
| `docs/alpha/runtime-extensions/` snapshots | 1 dir | DO NOT MOVE |
| `docs/beta/` snapshots | 4 dirs | DO NOT MOVE |
| `docs/gamma/cdd/` snapshots | 59 dirs | DO NOT MOVE |
| `docs/gamma/essays/` snapshots | 1 dir | DO NOT MOVE |
| `docs/gamma/rules/` snapshots | 1 dir | DO NOT MOVE |
| `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` | 1 file | DO NOT MOVE (golden-bound) |
| `docs/gamma/cdd` (path as code constant) | 1 path | DO NOT MOVE without code change |
