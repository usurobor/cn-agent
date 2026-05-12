# Self-Coherence — Cycle #344-B

## §Gap

**Issue:** #344 (meta-issue) — Cycle B scope: Reference notifier implementation

**Version/mode:** CDD 3.15.0 / substantial change (multi-AC, new template artifacts)

**Gap being closed:** `cdd/activation/SKILL.md §10` defines a transport-agnostic notification
interface (event vocabulary + adapter contract) but no reference implementation exists.
Tenants activating CDD have no copy-pasteable starting point for wiring notifications.

Cycle B closes the gap by delivering:
1. `activation/templates/telegram-notifier/` — reference notify.sh + cdd-notify.yml + README
2. `activation/templates/github-actions/` — cdd-artifact-validate.yml + cdd-cycle-on-merge.yml
3. `activation/templates/README.md` — top-level overview of all templates

---

## §Skills

**Tier 1 — CDD lifecycle:**
- `CDD.md` 3.15.0 — canonical lifecycle, §10 event vocabulary, §10.2 adapter contract
- `alpha/SKILL.md` — α role contract (this session)
- `activation/SKILL.md` — notification interface spec (§10, §10.1, §10.2, §11)

**Tier 2 — Engineering:**
- Shell authoring conventions (bash -n lint, shellcheck where available)
- YAML syntax (python3 yaml.safe_load validation before commit)
- Markdown authoring (word-count for README constraints)

**Tier 3 — Issue-specific:**
- No additional Tier 3 skills loaded; the work is template/docs authoring with
  no Go/TypeScript/specialized-runtime surface.

**Design:** Not required — this is a new template tree with no inter-package
contract dependencies. The scope is purely additive (new directory under
`activation/templates/`); no existing surfaces are modified.

**Plan:** Not required — artifact order is prescribed by the issue ACs (B.AC1–B.AC5);
no non-trivial sequencing beyond that order.

---

## §ACs

### B.AC1 — Telegram notifier reference impl works end-to-end

**Evidence:**
- `activation/templates/telegram-notifier/notify.sh` — single-file shell script, uses
  `curl` to call `https://api.telegram.org/bot${CDD_TELEGRAM_BOT_TOKEN}/sendMessage`
- Uses `$CDD_TELEGRAM_BOT_TOKEN` and `$CDD_TELEGRAM_CHAT_ID` from environment (set via
  GitHub secrets in the workflow); no hardcoded values anywhere in the script
- All 4 events handled: `cycle-open` (L42), `beta-verdict` (L50), `cycle-rc` (L57),
  `cycle-merge` (L63) in the `case "$EVENT" in` dispatch block
- `cdd-notify.yml` wires the workflow: jobs `notify-cycle-open`, `notify-beta-verdict`,
  `notify-cycle-rc`, `notify-cycle-merge` each set `CDD_EVENT` and call `notify.sh`
- Secrets-absent path: lines L32–39 check for empty TOKEN/CHAT_ID, log a warning, exit 0

**Status: MET**

### B.AC2 — Artifact validator catches missing files

**Evidence:**
- `activation/templates/github-actions/cdd-artifact-validate.yml`
- Triggers `on: push: branches: ['cycle/**', 'main']` — both required targets present
- Step "Check validate-release-gate.sh exists" (lines ~26–35) fails with:
  `"ERROR: validate-release-gate.sh not found at scripts/validate-release-gate.sh"`
  and `exit 1` if the script is absent
- Step "Validate CDD artifact presence" calls `bash scripts/validate-release-gate.sh`
  which already produces named-missing-artifact output (verified from the existing script)
- YAML validated via `python3 -c "import yaml; yaml.safe_load(...)"` → YAML OK

**Status: MET**

### B.AC3 — Workflow templates are copy-pasteable

**Evidence:**
- All 3 YAML files (cdd-notify.yml, cdd-artifact-validate.yml, cdd-cycle-on-merge.yml)
  use `${{ secrets.CDD_TELEGRAM_BOT_TOKEN }}` and `${{ secrets.CDD_TELEGRAM_CHAT_ID }}` —
  canonical `CDD_` namespace with no tenant-specific hardcoding
- Each file has inline comments marking what to customize (e.g., `# Customize:` blocks)
- No tenant names, repo names, or hardcoded paths that would need changing beyond
  optional secret name overrides

**Status: MET**

### B.AC4 — README walkthrough

**Evidence:**
- `activation/templates/telegram-notifier/README.md` — 266 words (≤300 limit)
- Steps covered: Step 1 (BotFather bot creation), Step 2 (chat ID retrieval via
  getUpdates), Step 3 (set secrets: CDD_TELEGRAM_BOT_TOKEN + CDD_TELEGRAM_CHAT_ID
  in GitHub Settings), Step 4 (wire workflow files), Step 5 (verify)
- No tokens in the README — only the secret names and API URL pattern

**Status: MET**

### B.AC5 — Notifier honors transport-agnostic event vocabulary

**Evidence:**
- `notify.sh` `case "$EVENT" in` block handles exactly:
  - `cycle-open` — matches `activation/SKILL.md §10.1` row 1
  - `beta-verdict` — matches §10.1 row 2
  - `cycle-rc` — matches §10.1 row 3
  - `cycle-merge` — matches §10.1 row 4
- Unknown event path (line ~72-76): prints warning to stderr and exits 0 with a
  warning log line to stdout — does not fail
- Empty `CDD_EVENT` also exits 0 with warning

**Status: MET**
