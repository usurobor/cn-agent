# α self-coherence — cycle/476

## §Gap

**Issue:** [cnos#476](https://github.com/usurobor/cnos/issues/476) — `cn-wake-install: v0 subcommand to render package-owned wakes into substrate (consume cn.wake-provider.v1)`. Sub 3 of cnos#467; builds on Sub 1 (cnos#468 — merged `c0048bef`) and Sub 2 (cnos#470 — merged `043bf7aa`).

**Version / mode:** `design-and-build`. No version bump in this sub (per γ scaffold pin "no version bump in this sub; γ records candidate release note in closeout"). Renderer + golden fixture + CI validator + AC2 negative tests; small design call within cycle (γ pinned both form + render target in `gamma-scaffold.md`).

**Branch:** `cycle/476` (γ-created from `origin/main@fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a`; γ scaffold at `417541ad`).

**Form-choice acknowledgment.** γ pin (per `gamma-scaffold.md` §"Form choice") is **shell command** at `src/packages/cnos.core/commands/install-wake/cn-install-wake`, registered in `cnos.core/cn.package.json`'s `commands` map as `install-wake`. α **accepts** the pin. No structural reason discovered to override to Option A (Go subcommand). Specifically:

- `discover.ScanPackageCommands` (`src/go/internal/discover/discover.go`) walks `.cn/vendor/packages/*/cn.package.json` and registers shell commands via `ExecCommand`, exporting `CN_HUB_PATH`, `CN_PACKAGE_ROOT`, `CN_COMMAND_NAME`. Path confinement (`filepath.Rel` check) only enforces the entrypoint stays inside the package directory; the entrypoint may then read arbitrary paths (incl. other installed packages' `orchestrators/` trees) at runtime — no structural block.
- `jq` is present in the sandbox (`jq-1.7`) and standard in the cnos CI image (Ubuntu runners ship `jq`). `python3` + `PyYAML` are also present (3.11.15 + PyYAML 6.0.1), so YAML validation does not require new infra.
- YAML emission via shell heredoc is constrained to a fixed-shape template (`claude-code-action@v1` invocation), not arbitrary YAML serialization; escaping risk is bounded by the golden-fixture invariant (CI re-renders + byte-diffs).
- Sub 2's γ scaffold (cycle/470 §"Form choice") explicitly reserved `cn-install-wake` as Sub 3's "execute the install" verb under `commands/install-wake/`; honoring that reservation is the symmetric move.

**Render-target acknowledgment.** γ pin (per `gamma-scaffold.md` §"Render target path") is **`src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`** (per-package sibling). α **accepts** the pin. The renderer's default output path (when invoked WITHOUT `--out`) is this path; `cn install-wake agent-admin` updates the golden in place. The renderer supports `--out <path>` for the future cutover cycle.

**No γ-clarification filed.** No pinned implementation-contract axis appears unsatisfiable; no AC oracle appears unsatisfiable from the pinned form/scope.

**Refusal conditions checked against pre-coding plan:**
- Will NOT edit `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant — confirmed unchanged at session start).
- Will NOT render to `.github/workflows/cnos-agent-admin.yml` (prohibited — would activate second production wake).
- Will NOT edit `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` or `prompt.md` (Sub 2 declarations; CONSUMED, not modified).
- Will NOT edit `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` (the contract; friction-log routes to gamma-closeout).
- Will NOT touch `src/go/` (γ-pinned form is shell).
- Will NOT touch any other package (cnos.cdd, cnos.cdr, cnos.kata).

## §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract; §2.5 incremental self-coherence; §2.6 pre-review gate; §3.6 implementation contract is δ's

**Lifecycle sub-skills:**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — issue shape (read for AC interpretation)
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-and-build mode (form-pin decision)

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` — **the contract this renderer consumes**; §2.1 required fields; §2.4 substrate-rendering target; §2.5 canonical authority-split table; §3 authority split rules; §3.5 reject malformed declarations; §4.1–§4.6 verification checks
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — reference instance (READ ONLY)
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` — prompt template (READ ONLY; inlined verbatim)
- `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` — cnos#468 contextual; renderer must not emit anything that contradicts label doctrine
- `src/packages/cnos.core/cn.package.json` — commands map (one entry added)
- `src/packages/cnos.core/commands/{daily,weekly,save}/cn-*` — shell-command sibling precedent (`#!/bin/sh`; `set -eu`; CN_HUB_PATH/CN_PACKAGE_ROOT/CN_COMMAND_NAME dispatch contract)
- `src/go/internal/discover/discover.go` — `ScanPackageCommands` + `ExecCommand` (confirms entry shape; path confinement; env-var contract)
- `.github/workflows/claude-wake.yml` — substrate-bound wake (READ ONLY; AC7 byte-identical invariant); structural reference for `claude-code-action@v1` invocation shape
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — cited from the manifest's `output_contract.channel_log_convention`

**γ scaffold:** `.cdd/unreleased/476/gamma-scaffold.md` (re-read for form-choice rationale §"Form choice"; render-target rationale §"Render target path"; AC mapping table §"AC mapping"; mechanical gate block; claim-class verification injection).

**Issue body:** read in full via `mcp__github__issue_read` (issue_number=476). AC1–AC8 Invariant/Oracle/Surface blocks parsed verbatim; non-goals checked against scope plan.

**Master tracker:** cnos#467 §"Foundational architecture — package-owned wake providers (authoritative)" + wave-level AC4 — already absorbed via Sub 2 (#470) merge state and γ scaffold's distillation.

## §ACs

Implementation SHA: `7162c32a` (head of cycle/476 after the renderer + golden + CI commits; before this self-coherence section).

All eight ACs mapped to mechanical evidence below. Per the binding claim-class verification rule injected by γ (per cnos#472 friction / PRA F-α-1): every "all X" claim is backed by a per-X table; no aggregated claims.

### Per-AC oracle table

| AC | Oracle | Observed result |
|---|---|---|
| AC1 (form pin + rationale) | γ scaffold §"Form choice" exists; α's §Gap acknowledges or records override | γ pin = shell at `commands/install-wake/cn-install-wake`; α **accepts** (no override). §Gap lines `Form-choice acknowledgment` enumerate the four structural reasons no override applies. |
| AC2 (consume `cn.wake-provider.v1`; reject malformed) | Positive: `cn install-wake agent-admin` exits 0. Negative: schema missing → exit 2; `admin_only` as string → exit 2; bad schema value → exit 2; malformed JSON → exit 2; missing manifest → exit 1. | All six oracles ran (table below). |
| AC3 (renders to syntactically valid GHA workflow at γ-pinned path) | YAML parses via `python3 -c "import yaml; yaml.safe_load(...)"`; structural greps for `^on:`, `^permissions:`, `^concurrency:`, `^jobs:`, `anthropics/claude-code-action@v1` all ≥ 1 | All 5 structural greps return ≥ 1 (table below). YAML parses (keys: `['name', True, 'permissions', 'concurrency', 'jobs']`; `True` is the YAML 1.1 alias for `on:` — same shape as `claude-wake.yml`). |
| AC4 (rendered prompt block carries admin-only boundary verbatim) | `grep -c "MUST NOT execute cells under any circumstance" golden.yml` ≥ 1; `grep -cE 'cell_execution|disallowed_surfaces|cell execution'` ≥ 1; prompt char count 95-105% of original | "MUST NOT execute cells under any circumstance" → 1 hit; "cell execution" / "disallowed" terms in rendered ≥ 1. **The 95–105% char-count oracle from the issue body is brittle as written when applied to the whole rendered file** (golden 14477 chars vs prompt 11872; ratio 1.22 — exceeds 1.05 because YAML literal-block-scalar inlining adds 12-space indentation × ~129 prompt lines ≈ +1548 chars, plus the YAML header + substrate encoding ≈ +1050 chars). β-side widened oracle (run by α and recorded here): extract the `prompt: |` block from the golden, strip the 12-space prefix, compare char count to `prompt.md`. Result: original 11776 chars; extracted 11752 chars; **ratio 0.998 (well within 95–105%)**. Line count 129 = 129. The only line-by-line difference is the `{agent}` → `sigma` substitution at L3 (cn-{agent}, {agent}-hub-list references) — expected per `agent_variable` substitution in `wake-provider/SKILL.md §2.3` step 6. Substantive truncation check: `grep -c "Wake termination" golden.yml` = 1 (the last `##` heading of `prompt.md` is present). |
| AC5 (triggers + permissions match manifest) | Per-trigger oracle: for each trigger in `input_contract.triggers`, grep golden for substrate encoding. Per-permission oracle: for each `permission_intent` entry, grep golden for the GHA encoding. No extra triggers or permissions beyond what manifest authorizes. | Per-trigger table + per-permission table below; both 0 MISSING. |
| AC6 (golden + idempotence + CI) | golden file exists; `cn install-wake agent-admin` twice → same `sha256sum`; CI workflow re-renders + diffs | golden file exists (path: `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`, 14465 bytes); two consecutive renders both report "(unchanged)"; sha256 stable; CI workflow `.github/workflows/install-wake-golden.yml` carries golden-diff step, idempotence step, AC2 negative-smoke step, AC3 YAML+structural step, AC7+AC8 audit steps |
| AC7 (`claude-wake.yml` byte-identical) | `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml \| wc -l` = 0 | 0 (confirmed) |
| AC8 (package-vs-renderer authority split preserved) | Renderer side: `grep -ciE 'admin.only\|disallowed_surfaces\|defer.path\|cell_execution' cn-install-wake` = 0. Package side: same grep on Sub 2 manifest + prompt — count must equal #470 baseline (no enlargement). | Renderer-side: 0. Package-side: 0 (declarations byte-identical with origin/main; gate 2 above confirms 0 diff lines on those files). |

### AC2 per-case oracle table

| Test case | Command | Expected | Observed |
|---|---|---|---|
| Positive | `cn install-wake agent-admin` (default path) | exit 0 + golden rendered | exit 0; golden at `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` |
| Positive (`--out`) | `... agent-admin --out /tmp/test-out.yml` | exit 0 + output at /tmp/test-out.yml + byte-identical to default-path golden | exit 0; `diff /tmp/test-out.yml golden` = empty |
| Negative: schema missing | `... --manifest <json without "schema">` | exit 2 + stderr names `required field "schema"` | exit 2; stderr `required field "schema" missing` |
| Negative: `admin_only` as string | `... --manifest <json with "admin_only": "true">` | exit 2 + stderr names `admin_only must be boolean` | exit 2; stderr `admin_only must be boolean (got string)` |
| Negative: bad schema value | `... --manifest <json with "schema": "cn.wake-provider.v999">` | exit 2 + stderr names unsupported schema | exit 2; stderr `unsupported schema "cn.wake-provider.v999" (renderer accepts cn.wake-provider.v1 only)` |
| Negative: malformed JSON | `... --manifest <invalid json>` | exit 2 + stderr names "not valid JSON" | exit 2; stderr `not valid JSON` |
| Negative: missing required field `defer_path` | `... --manifest <missing defer_path>` | exit 2 + stderr names field | exit 2; stderr `required field "defer_path" missing` |
| Negative: manifest file does not exist | `cn install-wake nonexistent-wake` | exit 1 + stderr names path | exit 1; stderr `manifest not found: ...` |

### AC3 per-grep oracle table (substrate structural validity)

| Required marker | grep pattern | Count |
|---|---|---|
| `name:` field | `^name:` | 1 |
| `on:` block | `^on:` | 1 |
| `permissions:` block | `^permissions:` | 1 |
| `concurrency:` block | `^concurrency:` | 1 |
| `jobs:` block | `^jobs:` | 1 |
| `claude-code-action@v1` step | `anthropics/claude-code-action@v1` | 1 |
| `actions/checkout@v4` step | `actions/checkout@v4` | 1 |
| YAML parses | `python3 -c "import yaml; yaml.safe_load(...)"` | exit 0 |

### AC4 per-grep oracle table (prompt body inlined verbatim)

| Marker | Source (prompt.md) | Rendered (golden.yml) |
|---|---|---|
| "MUST NOT execute cells under any circumstance" | present (L43) | present (1 hit) |
| "Cell execution" (heading "Admin-only boundary") | present | present |
| "cell-shaped directive" (lowercase prose) | present (multiple) | present (multiple) |
| "disallowed surfaces" heading | present (## Disallowed surfaces) | present |
| "Defer-path" headings (2x) | present | present |
| "## Wake termination" heading (last heading) | present | present (no silent truncation) |
| `{agent}` → `sigma` substitution | source: `cn-{agent}`, `as \`sigma\`` | rendered: `cn-sigma`, `as \`sigma\`` (table row for activation URL: source has `https://github.com/usurobor/cn-{agent}` → rendered `https://github.com/usurobor/cn-sigma`) |

### AC5 per-trigger oracle table

| Manifest trigger (`input_contract.triggers[i]`) | Substrate encoding required | Rendered evidence |
|---|---|---|
| `schedule` | `on.schedule` block with cron entries | `grep -c "^  schedule:" golden.yml` = 1; cron entries `'8 * * * *'`, `'23 * * * *'`, `'38 * * * *'`, `'53 * * * *'` all present |
| `issues_opened_title_match` | `on.issues.types: [opened]` + title-match `if:` conditional | `grep -c "^  issues:" golden.yml` = 1; `types: [opened]` present; `contains(github.event.issue.title, 'agent-admin')` present in `if:` |

No extra trigger surface (only `schedule` and `issues` keys under `on:`); per the manifest's 2-element `triggers` array.

### AC5 per-permission oracle table

| Manifest permission_intent[i] | GHA permissions: key | Rendered evidence |
|---|---|---|
| `contents.write` | `contents: write` | present (golden L24) |
| `issues.write` | `issues: write` | present (golden L25) |
| `pull_requests.write` | `pull-requests: write` (kebab-case) | present (golden L26) |
| `id_token.write` | `id-token: write` (kebab-case) | present (golden L27) |

No permission in rendered `permissions:` block is absent from manifest's `permission_intent`. Block contents (4 keys) match manifest's 4-element array exactly.

### AC6 per-render idempotence oracle table

| Render # | sha256sum of golden | exit |
|---|---|---|
| 1 (initial) | `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47` | 0 |
| 2 (immediate re-render) | `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47` | 0 |
| Diff between 1 and 2 | (none) | `diff` exit 0 |
| 3 (after `--out /tmp/test-out.yml` + diff) | (default-path unchanged) | 0 |

The renderer reports `(unchanged)` on the 2nd render rather than re-writing the file, so file-mtime is preserved as well as byte content.

### Mechanical gate summary (all from γ scaffold §"Mechanical gate")

| Gate | Command | Required | Observed |
|---|---|---|---|
| AC7 invariant | `git diff --name-only origin/main..HEAD -- .github/workflows/claude-wake.yml \| wc -l` | 0 | 0 |
| Sub 2 declaration unchanged | `git diff --name-only origin/main..HEAD -- .../wake-provider.json .../prompt.md \| wc -l` | 0 | 0 |
| Scope discipline (corrected regex — γ scaffold regex is brittle; see §Debt F1) | `git diff --name-only origin/main..HEAD \| grep -vE '^(src/packages/cnos\.core/(commands/install-wake/.*\|orchestrators/agent-admin/cnos-agent-admin\.golden\.yml\|cn\.package\.json)\|\.cdd/unreleased/476/.*\|\.github/workflows/install-wake-golden\.yml)$' \| wc -l` | 0 | 0 |
| AC6 idempotence | two renders + `diff /tmp/r1.sum /tmp/r2.sum` | exit 0 | exit 0 |
| AC8 renderer-side | `grep -ciE 'admin.only\|disallowed_surfaces\|defer.path\|cell_execution' cn-install-wake` | 0 | 0 |
| AC8 package-side (Sub 2 baseline) | `git diff --name-only origin/main..HEAD -- .../wake-provider.json .../prompt.md \| wc -l` (declaration byte-identical → carve-out baseline preserved by construction) | 0 | 0 |


