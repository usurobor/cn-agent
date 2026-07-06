# `cn repo install` — mocked human deliverables + receipt parity contract

> **Command name (operator ruling).** The installer is **`cn repo install`** —
> a `repo` noun-group that matches the noun-group command resolver and leaves
> room for future `cn repo status` / `cn repo doctor` / `cn repo uninstall`. It
> coexists with the `cn install <target>` group (`cn install wake` #549,
> `cn install cnos.core` #493). An interim draft used `cn install repo`; the
> operator verdict pinned `cn repo install`.

**Status:** Design surface (pre-implementation). Author: kappa/operator review.
**Purpose:** Lock the *human-visible* deliverables of the CDS repo-installer wave
**before** any code is written, and define the mechanism by which the
implementing cell must **prove in its receipt** that what shipped matches or
exceeds these mocks — and *how*.

This file is the source of truth an implementation cell (mode:
`design-and-build`) cites. Its ACs map 1:1 to the mocks below; its closeout
carries the [Receipt parity contract](#receipt-parity-contract).

> **Mock-first discipline.** These are *targets*, drawn before implementation.
> Nothing here ships yet. Where a mock shows output, that is the intended
> experience, not observed behavior. The implementing cell replaces "intended"
> with "observed" and records any delta.

---

## Mock A — `cn repo install --dry-run` (base mode)

Intended terminal experience, run from the root of an arbitrary repo with only
`cn` on PATH:

```console
$ cn repo install --dry-run
→ cnos repo install (dry-run) — no files will be written
✓ Git repository root: /home/dev/acme-api
✓ Resolved cnos release: v3.82.0
  Would fetch package index + tarballs → <cache>/index.json (3 packages)
  Would write .cn/deps.json:
      cnos.core  3.82.0
      cnos.cdd   3.82.0
      cnos.cds   3.82.0
  Would run: (internal) deps lock   → .cn/deps.lock.json
  Would run: (internal) deps restore → .cn/vendor/packages/ (3 packages)
  Would ensure .gitignore contains: .cn/vendor/
  Dispatch: none (base install only — no .github/workflows/ changes)

Planned committed diff (3 files):
  A  .cn/deps.json
  A  .cn/deps.lock.json
  M  .gitignore

Run without --dry-run to apply.
```

**Flags** (base mode): `--release latest|<tag>` resolves the CNOS release;
`--index <path-or-url>` overrides the index for tests / local / offline
development; `--packages a,b,c` overrides the default set; `--dispatch none`
(default) / `cds`; `--dry-run`.

**Human-deliverable invariants A1–A5**

| ID | Invariant |
|---|---|
| A1 | Fails with a clear error if not at a git repo root (does not silently walk up or scaffold). |
| A2 | Prints the exact release tag it will use; `--release latest` resolves to a concrete tag in the output. |
| A3 | Lists the precise planned committed diff — and it is exactly `.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore`. |
| A4 | States dispatch mode explicitly and, in base mode, states no `.github/workflows/` change. |
| A5 | Writes nothing (verified: `git status --porcelain` empty after `--dry-run`). |

---

## Mock B — base-install pull request (the committed diff)

Intended PR contents after `cn repo install` (no dispatch). Vendored packages
are **not** in the diff — they rehydrate from `cn.lock`.

```diff
+++ b/.cn/deps.json
+{
+  "schema": "cn.deps.v1",
+  "profile": "cds",
+  "packages": [
+    { "name": "cnos.core", "version": "3.82.0" },
+    { "name": "cnos.cdd",  "version": "3.82.0" },
+    { "name": "cnos.cds",  "version": "3.82.0" }
+  ]
+}

+++ b/.cn/deps.lock.json  (schema cn.lock.v2 — SHA-256-pinned entries)
+++ b/.gitignore          (+ .cn/vendor/)
```

**Invariants B1–B3**

| ID | Invariant |
|---|---|
| B1 | Diff is exactly three files (`.cn/deps.json`, `.cn/deps.lock.json`, `.gitignore`); no `.github/workflows/` file present. |
| B2 | `.cn/deps.lock.json` validates against `cn.lock.v2` and pins every package in `deps.json` by SHA-256. Versions are **exact** (resolved from the index), not semver ranges. |
| B3 | **Idempotent**: a second `cn repo install` produces no further diff (`git status --porcelain` empty), with no formatting churn. |
| B4 | **Dispatch guard**: until the renderer is generalized (#609), `cn repo install --dispatch cds` fails with an explicit "not implemented until renderer generalization lands" error and writes **no** partial `.github/workflows/cnos-cds-dispatch.yml`. |

---

## Mock C — `cn repo install --dispatch cds` for a **non-sigma** agent

The load-bearing mock: dispatch rendered for agent `acme` with the caller's own
secret. This is the proof the sigma binding is gone.

```console
$ cn repo install --dispatch cds \
    --agent acme \
    --workflow-pat-secret ACME_WORKFLOW_PAT \
    --bot-name "acme-bot" \
    --bot-id 12345678
→ cnos repo install — dispatch: cds (agent: acme)
✓ base install complete (.cn/deps.json, cn.lock)
✓ rendered .github/workflows/cnos-cds-dispatch.yml
  identity:  acme  (bot acme-bot / 12345678)
  pat secret: ACME_WORKFLOW_PAT
⚠ dispatch grants scheduled write access on merge. Review before merging.
  This changes .github/workflows/ — the installing token needs `workflow` scope.
```

Intended rendered header + wiring (the sections that carry identity):

```yaml
concurrency:
  group: cds-dispatch-acme          # was: cds-dispatch-sigma
  cancel-in-progress: false
...
      - uses: actions/checkout@v4
        with:
          token: ${{ secrets.ACME_WORKFLOW_PAT }}   # was: SIGMA_WORKFLOW_PAT
      - uses: anthropics/claude-code-action@v1
        with:
          claude_code_oauth_token: ${{ secrets.CLAUDE_CODE_OAUTH_TOKEN }}
          github_token: ${{ secrets.ACME_WORKFLOW_PAT }}
          bot_name: "acme-bot"        # was: sigma@cnos.cn-sigma.cnos
          bot_id: "12345678"          # was: 41898282
```

**Invariants C1–C6**

| ID | Invariant |
|---|---|
| C1 | Base install runs first; dispatch is layered on top. |
| C2 | Missing `--workflow-pat-secret` / identity → **fails early** with a message naming what's required (no partial render). |
| C3 | Rendered file writes exactly to `.github/workflows/cnos-cds-dispatch.yml`. |
| C4 | **Zero sigma leak** — `grep -iE 'sigma|SIGMA_WORKFLOW_PAT|41898282' cnos-cds-dispatch.yml` returns nothing when agent ≠ sigma. (Reuse the AC7/AC8 leak-audit pattern already in `install-wake-golden.yml`.) |
| C5 | `--agent sigma` still reproduces the current sigma-bound output byte-for-byte (backward compat; the existing golden fixture still passes). |
| C6 | Dispatch never pushes to `main`; PR-only, and the CLI states the `workflow`-scope PAT requirement. |

---

## Mock D — GitHub UI (no-terminal) path

`workflow_dispatch` inputs on `.github/workflows/cnos-install.yml`, delegating to
the same `cn repo install` command (not reimplemented in YAML):

```yaml
inputs:
  release:              { default: latest }
  install_dispatch:     { default: false }        # true → --dispatch cds
  workflow_pat_secret:  { default: CNOS_WORKFLOW_PAT }
  agent:                { default: cnos }
```

**Invariants D1–D4**

| ID | Invariant |
|---|---|
| D1 | Base run (dispatch off) opens a PR with the Mock-B diff, using the default `GITHUB_TOKEN`. |
| D2 | Dispatch run either opens a PR containing the Mock-C workflow **or** fails with a clear message that a `workflow`-scoped PAT is required — never a half-applied state. |
| D3 | The job body invokes `cn repo install …` — install logic is not duplicated in shell. |
| D4 | Never pushes to `main`. |

---

## Mock E — tenant-portable dispatch wake (dogfooding: cnos#606 C4, C5)

The first external tenant (`tsc`, cnos#606) found the rendered wake does
`cd src/go && go build ./cmd/cn` — valid only inside the cnos repo — and that
`cn init` drops a full agent-hub (`spec/SOUL.md`, `agent/`, `threads/`,
`state/`) into the tree. A tenant wants **labels + wakes + `.cn/` runtime**, not
a home hub, and its wake must acquire `cn` the way a tenant can.

Intended rendered dispatch wake for a tenant (identity elided; see Mock C):

```yaml
      - name: Install cn (tenant-portable — no src/go build)
        run: curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh
             | BIN_DIR="$HOME/.local/bin" sh
      # ... claude-code-action work phase ...
      - name: Mechanical checkpoint + PR finalizer
        run: cn cell finalize --base-sha "${CN_WAKE_BASE_SHA:-}"   # installed cn, NOT go build
```

**Invariants E1–E4**

| ID | Invariant |
|---|---|
| E1 | `cn repo install` (base) writes **no** agent-hub scaffold — no `spec/SOUL.md`, `agent/`, `threads/`, `state/`; only `.cn/` + `cn.lock` + `.gitignore` (C4). |
| E2 | A rendered tenant dispatch wake contains **no** `cd src/go` / `go build ./cmd/cn`; it acquires `cn` via `install.sh` or a pinned release (C5). |
| E3 | The rendered wake's finalizer/engine steps invoke the installed `cn` (e.g. `cn cell finalize`), runnable in a repo with no `src/go`. |
| E4 | `--agent sigma` inside the cnos repo still reproduces the current `go build` self-wake byte-for-byte (backward compat; the golden is not broken by tenant mode). |

---

## Mock F — CLI ergonomics (dogfooding: cnos#606 C3)

The tenant hit: `cn --version` → "Unknown command"; `cn init --help` **scaffolded
a stray `cn---help/` directory**; `cn help` lists `issues-fsm` / `cell-finalize`
but the real invocation is `cn issues fsm` / `cn cell finalize`.

```console
$ cn --version
cn 3.82.1 (…)
$ cn repo install --help
cn repo install — make this repository CDS-ready …
$ cn init --help        # a flag is never a hub name
✗ unknown flag --help for `cn init` (did you mean `cn init <name>`?)
```

**Invariants F1–F4**

| ID | Invariant |
|---|---|
| F1 | `cn --version` / `-V` prints the version (exit 0), not "Unknown command". |
| F2 | `--help` / `-h` works on `cn` and every command incl. `cn repo install`. |
| F3 | `cn init --help` (any unrecognized `--flag`) **refuses** with a clear error and scaffolds nothing — no stray `cn---help/`. |
| F4 | `cn help` display names match real invocation (`cn issues fsm`, `cn cell finalize`), or list both forms — no name/invocation mismatch. |

---

## Mock G — PAT-free mechanical FSM engine + secrets runbook (dogfooding: cnos#606 C6, C7)

`cn issues fsm evaluate --apply` is mechanical — it needs no agent token. The
tenant had no fallback: every wake required `SIGMA_WORKFLOW_PAT`, and there was
no doc naming the secrets or a hub-less label path.

**Invariants G1–G4**

| ID | Invariant |
|---|---|
| G1 | `cn repo install --dispatch cds` can render a **mechanical FSM-engine** wake that runs `cn issues fsm evaluate --apply` on the default `GITHUB_TOKEN` — no agent PAT required. |
| G2 | The agent (claude-code-action) wake tier remains the only one that needs a `workflow`-scope PAT + `CLAUDE_CODE_OAUTH_TOKEN`; the two tiers are separable. |
| G3 | Dispatch install ensures the canonical FSM labels **without a full hub** (hub-less `cn labels sync` / the cnos#493 mechanism) (C7). |
| G4 | Installing dispatch emits/points to a tenant secrets runbook naming every required secret per tier. |

---

## Receipt parity contract

The implementing cell (`cell_kind: implementation`) MUST emit this block in its
γ close-out. It is how the cell **confirms match-or-exceed, and how**. It is a
hard gate: any row with `verdict: miss` blocks convergence to `status:review`.

```yaml
mock_parity:
  source: docs/development/design/cn-repo-install-MOCKS.md
  source_commit: "<sha of this file the cell built against>"
  rows:
    # one row per invariant A1..A5, B1..B4, C1..C6, D1..D4, E1..E4, F1..F4, G1..G4
    - id: A3
      expectation: "dry-run lists exactly .cn/deps.json, .cn/deps.lock.json, .gitignore"
      observed: "<what the built command actually printed>"
      evidence: "<test name / transcript path / commit>"
      verdict: match | exceed | miss
      how: "<why it matches; if exceed, the additional capability and why it is safe;
             if miss, the gap + planned disposition>"
  summary:
    matched: <n>
    exceeded: <n>
    missed: <n>          # MUST be 0 to converge
    exceed_justified: <bool>   # every 'exceed' row has a non-empty, safe `how`
```

**Rules**
- **Coverage**: every invariant ID in Mocks A–G has exactly one row. A missing ID
  is itself a `miss`. (Mocks E–G are the cnos#606 dogfooding surface — tenant
  portability, ergonomics, PAT-free engine — and are not optional.)
- **`exceed` is not a free pass**: an `exceed` verdict requires `how` to state the
  extra capability *and* why it does not violate a non-goal (e.g. exceeding into
  autonomous-write territory in base mode is a `miss`, not an `exceed`).
- **Evidence is durable**: `evidence` points at a test, a captured transcript
  under `.cdd/unreleased/<N>/`, or a commit — not prose.
- β rejects the cell if any row's `observed`/`how` is unfalsifiable against the
  cited evidence.

---

## Non-goals (inherited by the wave)

Hosted live package registry · direct pushes to `main` · silent autonomous-write
install · sigma-only default · agent-in-the-loop for base install · manual
terminal for the GitHub UI path.
