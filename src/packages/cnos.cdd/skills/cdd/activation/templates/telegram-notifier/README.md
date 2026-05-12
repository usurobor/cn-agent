# CDD Telegram Notifier — Setup Walkthrough

This template wires CDD cycle events to a Telegram chat using a Bot API token.
Four events are covered: `cycle-open`, `beta-verdict`, `cycle-rc`, `cycle-merge`.

---

## Step 1 — Create a bot via BotFather

1. Open Telegram and start a chat with **@BotFather**.
2. Send `/newbot`, follow prompts to name your bot.
3. BotFather replies with an API token in the form `123456:ABC-DEF...`.
   Copy it — you will set it as a secret; do not save it anywhere else.

## Step 2 — Get your chat ID

1. Add the bot to the Telegram group or channel that should receive notifications
   (or start a direct chat with it).
2. Send any message in the chat.
3. Fetch `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates` in a browser.
4. Find `"chat":{"id": ...}` in the response. That integer is your chat ID.

## Step 3 — Set repository secrets

In your GitHub repository: **Settings → Secrets and variables → Actions → New repository secret**

| Secret name | Value |
|---|---|
| `CDD_TELEGRAM_BOT_TOKEN` | The token from Step 1 |
| `CDD_TELEGRAM_CHAT_ID` | The chat ID from Step 2 |

Never commit either value to version control.

## Step 4 — Wire the workflow

1. Copy `notify.sh` to `.github/cdd/notify.sh` in your repository and make it executable:
   `chmod +x .github/cdd/notify.sh`
2. Copy `cdd-notify.yml` to `.github/workflows/cdd-notify.yml`.
3. Commit and push to `main`.

## Step 5 — Verify

Push to a `cycle/*` branch. A Telegram message should appear within seconds.
If not, check the Actions run log — `notify.sh` logs the delivery attempt to stdout.

---

*Adapter contract reference: `cdd/activation/SKILL.md §10.2`*
