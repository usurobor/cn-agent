# CDD Activation Templates

Copy-pasteable templates for wiring CDD into a tenant repository.
No edits needed beyond replacing secret names if you diverged from the `CDD_` namespace.

---

## Quick Start — Telegram Notifier

Get CDD cycle notifications in Telegram in four steps:

**1. Create a bot.** Message **@BotFather** on Telegram, send `/newbot`, follow prompts.
BotFather returns an API token (`123456:ABC-DEF...`). Keep it secret.

**2. Get your chat ID.** Add the bot to the target group, send any message, then fetch
`https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`. The `"chat":{"id": ...}` integer
is your chat ID.

**3. Set repository secrets.**
In GitHub: **Settings → Secrets and variables → Actions → New repository secret**

| Secret name | Value |
|---|---|
| `CDD_TELEGRAM_BOT_TOKEN` | Token from step 1 |
| `CDD_TELEGRAM_CHAT_ID` | Chat ID from step 2 |

**4. Wire the workflow.**
Copy `telegram-notifier/notify.sh` → `.github/cdd/notify.sh` and
`telegram-notifier/cdd-notify.yml` → `.github/workflows/`. Commit and push.

Push to a `cycle/*` branch to verify — a message should arrive within seconds.
See `telegram-notifier/README.md` for the full reference walkthrough.

---

## telegram-notifier/

CDD notification adapter (`cdd/activation/SKILL.md §10.2`) for Telegram.

| File | Purpose |
|---|---|
| `notify.sh` | Shell script: posts to Telegram Bot API for all 4 CDD events |
| `cdd-notify.yml` | GitHub Actions workflow: triggers notify.sh on the right events |
| `README.md` | Detailed bot registration + secret setup walkthrough |

## github-actions/

CI workflow templates for CDD artifact governance.

| File | Purpose |
|---|---|
| `cdd-artifact-validate.yml` | Validates `.cdd/` structure on every push to `cycle/**` and `main`; requires `scripts/validate-release-gate.sh` in the tenant repository |
| `cdd-cycle-on-merge.yml` | Runs project test suite + emits cycle-merge notification on merge to main |

Copy both files to `.github/workflows/` and replace the test command placeholder
in `cdd-cycle-on-merge.yml` with your project's test runner.

---

*Governing spec: `cdd/activation/SKILL.md §9–§10`*
