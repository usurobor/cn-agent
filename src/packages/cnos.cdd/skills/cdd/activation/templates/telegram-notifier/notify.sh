#!/usr/bin/env bash
# notify.sh — CDD Telegram notification adapter (reference implementation)
#
# Adapter contract: cdd/activation/SKILL.md §10.2
# Event vocabulary: cdd/activation/SKILL.md §10.1
#
# Required environment variables:
#   CDD_TELEGRAM_BOT_TOKEN  — Telegram Bot API token (GitHub secret)
#   CDD_TELEGRAM_CHAT_ID    — Telegram chat or channel ID (GitHub secret)
#   CDD_EVENT               — one of: cycle-open | beta-verdict | cycle-rc | cycle-merge
#
# Optional payload variables (used to construct message text):
#   CDD_CYCLE_N   — cycle number (issue number)
#   CDD_VERDICT   — APPROVED or REQUEST_CHANGES (for beta-verdict event)
#   CDD_ROUND     — review round number (for beta-verdict and cycle-rc events)
#   CDD_SHA       — commit SHA (for cycle-rc and cycle-merge events)
#   CDD_TITLE     — issue/cycle title (for cycle-open event)
#
# Exit codes:
#   0 — successful delivery, or graceful skip with warning
#   1 — delivery failed (curl error or Telegram API error)
#
# Usage (from GitHub Actions workflow step):
#   bash scripts/notify.sh
#   or copy to any path and call directly.
#
# To test locally (after setting env vars):
#   CDD_EVENT=cycle-open CDD_CYCLE_N=42 CDD_TITLE="My cycle" bash notify.sh

set -euo pipefail

# ── Secret presence check ─────────────────────────────────────────────────────

TOKEN="${CDD_TELEGRAM_BOT_TOKEN:-}"
CHAT_ID="${CDD_TELEGRAM_CHAT_ID:-}"

if [[ -z "$TOKEN" || -z "$CHAT_ID" ]]; then
  echo "WARNING: CDD_TELEGRAM_BOT_TOKEN or CDD_TELEGRAM_CHAT_ID not set — skipping notification" >&2
  # Graceful failure per B.AC1 (unknown-token case) and §10.2 contract row 4 (log attempt)
  echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) notify.sh: secrets absent — event=${CDD_EVENT:-<unset>} destination=telegram skipped"
  exit 0
fi

# ── Event dispatch ────────────────────────────────────────────────────────────

EVENT="${CDD_EVENT:-}"
CYCLE_N="${CDD_CYCLE_N:-?}"
VERDICT="${CDD_VERDICT:-}"
ROUND="${CDD_ROUND:-?}"
SHA="${CDD_SHA:-}"
TITLE="${CDD_TITLE:-}"

# Construct the short SHA for display (first 8 chars if available)
SHORT_SHA=""
if [[ -n "$SHA" ]]; then
  SHORT_SHA="${SHA:0:8}"
fi

case "$EVENT" in
  cycle-open)
    # Payload minimum: cycle number, issue title, branch name
    BRANCH_LABEL="cycle/${CYCLE_N}"
    MESSAGE="*CDD cycle-open* — cycle \`#${CYCLE_N}\`
Branch: \`${BRANCH_LABEL}\`
Title: ${TITLE:-<no title>}"
    ;;

  beta-verdict)
    # Payload minimum: cycle number, verdict, round number
    MESSAGE="*CDD beta-verdict* — cycle \`#${CYCLE_N}\` round ${ROUND}
Verdict: *${VERDICT:-<unknown>}*"
    ;;

  cycle-rc)
    # Payload minimum: cycle number, RC round number, fix-round commit SHA
    MESSAGE="*CDD cycle-rc* — cycle \`#${CYCLE_N}\` round ${ROUND}
Fix-round commit: \`${SHORT_SHA:-<unknown>}\`"
    ;;

  cycle-merge)
    # Payload minimum: cycle number, merge commit SHA
    MESSAGE="*CDD cycle-merge* — cycle \`#${CYCLE_N}\`
Merge commit: \`${SHORT_SHA:-<unknown>}\`"
    ;;

  "")
    echo "WARNING: CDD_EVENT not set — skipping notification" >&2
    echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) notify.sh: CDD_EVENT unset — skipped"
    exit 0
    ;;

  *)
    # Unknown events: warn but do not fail (B.AC5 requirement)
    echo "WARNING: unknown CDD_EVENT='${EVENT}' — not a canonical event; skipping notification" >&2
    echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) notify.sh: unknown event='${EVENT}' destination=telegram skipped"
    exit 0
    ;;
esac

# ── Delivery ──────────────────────────────────────────────────────────────────

TELEGRAM_API_URL="https://api.telegram.org/bot${TOKEN}/sendMessage"

# Log the delivery attempt before sending (§10.2 contract row 4)
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) notify.sh: delivering event='${EVENT}' cycle=${CYCLE_N} destination=telegram"

RESPONSE=$(curl \
  --silent \
  --show-error \
  --fail \
  --max-time 15 \
  --retry 2 \
  --retry-delay 3 \
  -X POST "$TELEGRAM_API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"chat_id\": \"${CHAT_ID}\",
    \"text\": $(printf '%s' "$MESSAGE" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'),
    \"parse_mode\": \"Markdown\"
  }" 2>&1) || {
    echo "ERROR: curl failed delivering event='${EVENT}' to Telegram — ${RESPONSE}" >&2
    exit 1
  }

# Check Telegram API ok field
OK=$(echo "$RESPONSE" | python3 -c 'import json,sys; d=json.load(sys.stdin); print(d.get("ok",""))' 2>/dev/null || echo "")
if [[ "$OK" != "True" ]]; then
  echo "ERROR: Telegram API returned non-ok response for event='${EVENT}': ${RESPONSE}" >&2
  exit 1
fi

echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) notify.sh: delivered event='${EVENT}' cycle=${CYCLE_N} destination=telegram ok"
exit 0
