# `cn repo status` — mocked human deliverables + receipt parity contract

**Status:** Design surface (pre-implementation). Author: δ/α, cnos#656 (Phase 1
of the cnos#655 `cn repo` lifecycle wave).
**Purpose:** Lock the *human-visible* deliverables of `cn repo status` and the
`.cn/repo.state.json` ledger **before** implementation, mirroring the
`cn repo install` wave's own discipline
([`cn-repo-install-MOCKS.md`](cn-repo-install-MOCKS.md), A6 of
[`cn-repo-lifecycle.md`](cn-repo-lifecycle.md)). The implementing cell's
closeout carries the [Receipt parity contract](#receipt-parity-contract)
below with `missed: 0`.

> **Mock-first discipline.** These are *targets*, drawn before implementation.
> The implementing cell replaces "intended" with "observed" and records any
> delta.

---

## Mock A — `cn repo status` (clean, no drift)

```console
$ cn repo status
→ cn repo status — /home/dev/acme-api
✓ source: stable @ 3.82.6 (github:usurobor/cnos)
✓ deps.json / deps.lock.json / vendor: in sync (3 packages)
    cnos.core  3.82.6
    cnos.cdd   3.82.6
    cnos.cds   3.82.6
✓ dispatch: cds (engine tier) — cnos-cds-dispatch.yml matches ledger
✓ labels: ok (canonical set present, no drift)
✓ orphan vendored packages: none
✓ locally-edited managed files: none
  update available: no (3.82.6 is the latest resolved release)

No drift.
```

## Mock B — `cn repo status` (drifted)

```console
$ cn repo status
→ cn repo status — /home/dev/acme-api
✓ source: stable @ 3.82.4 (github:usurobor/cnos)
⚠ deps.lock.json / vendor: cnos.cds 3.82.3 vendored, lock wants 3.82.4 (run `cn repo install`)
⚠ dispatch: cds (engine tier) — cnos-cds-dispatch.yml differs from ledger
    classification: renderer-moved (matches a fresh re-render of the recorded
    tier — not a local edit; run `cn repo repair` to update the ledger sha)
✗ labels: drifted (missing: status:blocked; unknown: priority:urgent)
⚠ orphan vendored packages: cnos.legacy-foo (not in deps.lock.json)
⚠ locally-edited managed files: .cn/deps.json (differs from ledger AND from
    what install would write — user edit, not renderer drift)
  update available: yes — 3.82.6 (run `cn repo update`)

Drift detected (see ⚠/✗ rows above). Exit code: 1 (--check only).
```

## Mock C — `cn repo status --json`

```console
$ cn repo status --json
```
```json
{
  "schema": "cn.repo.status.v1",
  "source": { "channel": "stable", "release": "3.82.4",
              "index": "github:usurobor/cnos/releases/3.82.4/index.json" },
  "packages": [
    { "name": "cnos.cds", "desired": "3.82.4", "locked": "3.82.4",
      "vendored": "3.82.3", "in_sync": false }
  ],
  "dispatch": { "present": true, "tier": "cds", "id": "cnos-cds-dispatch",
                "path": ".github/workflows/cnos-cds-dispatch.yml",
                "drift": "renderer_moved" },
  "labels": { "status": "drifted", "missing": ["status:blocked"],
              "unknown": ["priority:urgent"] },
  "orphan_packages": ["cnos.legacy-foo"],
  "local_edits": [ { "path": ".cn/deps.json", "classification": "user_edit" } ],
  "update_available": { "available": true, "release": "3.82.6" },
  "ledger": { "present": true, "reconstructed": false },
  "drift": true
}
```

## Mock D — `cn repo status --check` exit codes

```console
$ cn repo status --check ; echo "exit=$?"
... (same output as Mock A or B) ...
exit=0   # Mock A (no drift)
exit=1   # Mock B (drift detected)
```

**No writes.** `cn repo status` (with or without `--json`/`--check`) never
writes `.cn/repo.state.json` or any other file (P3). When the ledger is
absent, `--json` includes a `ledger.reconstructed: true` best-effort candidate
(A3) computed in-memory only.

---

## Human-deliverable invariants

| ID | Invariant |
|---|---|
| A1 | Reports source channel/release/index. |
| A2 | Reports each package's desired vs locked vs vendored version and an `in_sync` verdict per package. |
| A3 | Reports dispatch presence, tier, id, path, and drift classification (`matches_ledger` / `user_edit` / `renderer_moved` / absent). |
| A4 | Reports canonical label status: `ok` / `missing: [...]` / `unknown: [...]` (unknown = present but not in the canonical set — informational, never an error). |
| A5 | Reports orphan vendored packages (materialized but absent from `deps.lock.json`). |
| A6 | Reports locally-edited managed files, distinguishing `user_edit` from ledger-vs-fresh-render drift per the same classification as A3. |
| A7 | Reports whether a newer release is available, without performing a resolution that mutates any file (read-only network check, best-effort; degrades to `unknown` offline). |
| B1 | `--json` emits `schema: cn.repo.status.v1` and is machine-parseable (`encoding/json` round-trips). |
| B2 | `--json` sets top-level `drift: true` iff any of packages/dispatch/labels/orphans/local_edits reports non-clean. |
| C1 | `--check` exits nonzero (1) iff `drift == true`; exits 0 otherwise. Without `--check`, exit code is always 0 regardless of drift (drift is reported, not fatal). |
| D1 | `cn repo status` never writes `.cn/repo.state.json` or any other file, with or without any flag (P3). Verified: `git status --porcelain` empty before and after. |
| D2 | When `.cn/repo.state.json` is absent, status still runs (A3 backfill semantics) and reports `ledger.present: false`, `ledger.reconstructed: true` in `--json`, with dispatch/local-edit classification degraded to best-effort (no ledger-recorded sha to compare against — only `matches_fresh_render` / `unknown`). |
| E1 | Help text (`cn repo status --help` and `cn help`) explicitly distinguishes this command from the existing global `cn status` (binary/hub health) — e.g. "repo-scoped: is THIS repo's CNOS install in sync? (see `cn status` for hub/binary health)". |

---

## Receipt parity contract

The implementing cell (`cell_kind: implementation`) MUST emit this block in
its γ close-out. Any row with `verdict: miss` blocks convergence to
`status:review`.

```yaml
mock_parity:
  source: docs/development/design/cn-repo-status-MOCKS.md
  source_commit: "<sha of this file the cell built against>"
  rows:
    # one row per invariant A1..A7, B1..B2, C1, D1..D2, E1
    - id: A1
      expectation: "reports source channel/release/index"
      observed: "<what the built command actually printed>"
      evidence: "<test name / transcript path / commit>"
      verdict: match | exceed | miss
      how: "<why it matches>"
  summary:
    matched: <n>
    exceeded: <n>
    missed: <n>          # MUST be 0 to converge
    exceed_justified: <bool>
```

**Rules** (mirrors `cn-repo-install-MOCKS.md`'s Receipt parity contract):
- Coverage: every invariant ID above has exactly one row; a missing ID is
  itself a `miss`.
- `exceed` requires `how` to state the extra capability and why it is safe.
- Evidence is durable: a test, a transcript under `.cdd/unreleased/656/`, or a
  commit — not prose.
