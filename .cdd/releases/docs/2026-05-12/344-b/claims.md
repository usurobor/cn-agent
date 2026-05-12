---
cycle: 344-b
role: alpha
type: claims
---
# Honest-Claim Manifest — Cycle #344-B

## Reproducibility claims

- Claim: `notify.sh` handles all 4 canonical CDD events: `cycle-open`, `beta-verdict`, `cycle-rc`, `cycle-merge`
  Verification: `grep -E "cycle-open|beta-verdict|cycle-rc|cycle-merge" src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/notify.sh` → 4 hits

- Claim: `activation/templates/README.md` Quick Start walkthrough is ≤300 words
  Verification: `wc -w src/packages/cnos.cdd/skills/cdd/activation/templates/README.md` ≤ 300

- Claim: `telegram-notifier/README.md` is ≤300 words
  Verification: `wc -w src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/README.md` ≤ 300

- Claim: All YAML files are syntactically valid
  Verification: `python3 -c "import yaml; [yaml.safe_load(open(f).read()) for f in ['src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/cdd-notify.yml', 'src/packages/cnos.cdd/skills/cdd/activation/templates/github-actions/cdd-artifact-validate.yml', 'src/packages/cnos.cdd/skills/cdd/activation/templates/github-actions/cdd-cycle-on-merge.yml']]"` → no exception

- Claim: `notify.sh` passes shell syntax check
  Verification: `bash -n src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/notify.sh` → exit 0

## Source-of-truth alignment claims

- Claim: Event names in `notify.sh` match `activation/SKILL.md §10.1` exactly
  Verification: `grep -E '"cycle-open"|"beta-verdict"|"cycle-rc"|"cycle-merge"' src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/notify.sh` → 4 hits matching §10.1 event vocabulary table

- Claim: Secret names in templates use `CDD_` prefix per `activation/SKILL.md §11`
  Verification: `grep -E 'CDD_TELEGRAM_BOT_TOKEN|CDD_TELEGRAM_CHAT_ID' src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/notify.sh` → ≥2 hits; no other token variables used

## Wiring claims

- Claim: `cdd-artifact-validate.yml` calls `scripts/validate-release-gate.sh` (which exists in cnos)
  Verification: `test -f scripts/validate-release-gate.sh` → exit 0 (file exists in cnos repo)

- Claim: No hardcoded tokens appear in any template file
  Verification: `grep -rE '[0-9]{9}:[A-Za-z0-9_-]{35}' src/packages/cnos.cdd/skills/cdd/activation/templates/` → 0 hits
