---
cycle: 508
artifact: pass4-golden-impact-map.md
role: alpha
ac: AC4
---

# Pass 4A — CI/Golden Impact Map (AC4)

**Purpose:** Identify every `docs/(alpha|beta|gamma)/` path cited in known `*.golden.yml` files, map each to its bundle, and state the re-render/guard implication if that bundle moves.

**Golden files discovered in `pass4-golden-inventory.txt`:**

```
./.github/workflows/install-wake-golden.yml   ← CI workflow (not a golden fixture)
./.cdd/unreleased/508/pass4-golden-inventory.txt  ← self-reference (this run)
./src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
./src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

The `install-wake-golden.yml` is a CI workflow that *verifies* the golden fixtures match re-renders — it does not itself contain triad-doc path references in its own body (confirmed by inspection).

---

## Golden File 1: `cnos-cds-dispatch.golden.yml`

**Full path:** `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`

**Rendered from:** `orchestrators/cds-dispatch/prompt.md` + `orchestrators/cds-dispatch/wake-provider.json`

**Triad paths cited (line 225):**

| Path Cited | Bundle | Reader-Intent Destination |
|---|---|---|
| `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` | `docs/gamma/conventions/` | stays/defer (golden-bound — see below) |

**Context:** The path appears in an embedded wake prompt as a cross-reference citation:
```
- **Channel log convention (read-only for dispatch):** [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) — cited to document the admin/dispatch writer-locality split; the dispatch wake does NOT write channel entries.
```

**Re-render/guard implication:**
- The `.golden.yml` is byte-checked by `.github/workflows/install-wake-golden.yml` on every push to `main`/`cycle/*`.
- If `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` moves, the relative path `../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` in the rendered golden breaks.
- The prompt source at `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` MUST be updated to point to the new location.
- After updating the prompt, `cn install-wake cds-dispatch` must be re-run to regenerate the golden.
- The CI `install-wake-golden.yml` workflow will fail if the golden is not re-rendered.
- **Gate:** `docs/gamma/conventions/` bundle must NOT move until both golden files are updated and re-rendered.

---

## Golden File 2: `cnos-agent-admin.golden.yml`

**Full path:** `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`

**Rendered from:** `orchestrators/agent-admin/prompt.md` + `orchestrators/agent-admin/wake-provider.json`

**Triad paths cited (2 occurrences):**

| Line | Path Cited | Bundle | Reader-Intent Destination |
|---|---|---|---|
| 82 | `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` | `docs/gamma/conventions/` | stays/defer (golden-bound) |
| 168 | `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` | `docs/gamma/conventions/` | stays/defer (golden-bound) |

**Context line 82:** Inside status-reporting instruction:
```
per the entry format in [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) §5.
```

**Context line 168:** Inside cross-references section:
```
- **Channel log convention:** [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) — entry format, cursor mechanics, class taxonomy.
```

**Re-render/guard implication:**
- Both occurrences use the same absolute relative path (`../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`).
- If `AGENT-ACTIVATION-LOG-v0.md` moves, BOTH occurrences break.
- The prompt source at `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` MUST be updated.
- After updating, `cn install-wake agent-admin` must be re-run to regenerate the golden.
- **Gate:** Same as golden file 1 — `docs/gamma/conventions/` must not move without this re-render.

---

## Additional Golden Files

From `pass4-golden-inventory.txt`, the only additional `*.golden.yml` files in the repo are the 2 named above. The `install-wake-golden.yml` is a CI workflow, not a golden fixture. No additional golden fixtures with triad path references were found.

**Verification command run:**
```bash
grep -n 'docs/(alpha|beta|gamma)/' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
# → 1 match (line 225)

grep -n 'docs/(alpha|beta|gamma)/' src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
# → 2 matches (lines 82, 168)
```

---

## Impact Summary

| Golden File | Paths Cited | Bundle(s) Affected | Re-render Required? |
|---|---|---|---|
| `cnos-cds-dispatch.golden.yml` | 1 (`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`) | `docs/gamma/conventions/` | Yes — `cn install-wake cds-dispatch` |
| `cnos-agent-admin.golden.yml` | 2 (`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` × 2) | `docs/gamma/conventions/` | Yes — `cn install-wake agent-admin` |

**Only one bundle is golden-bound: `docs/gamma/conventions/`**

This bundle must be treated as a special case in the 4B–4E sequencing. It cannot move until:
1. `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` is updated
2. `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` is updated
3. Both golden files are re-rendered and committed
4. The `install-wake-golden.yml` CI check passes

All other bundles are free from golden-bound constraints (though they may have other constraints per `pass4-do-not-move.md` and `pass4-move-map.md`).
