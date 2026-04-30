# Operator Manual

Day-2 operations for a running cnos agent. Assumes install and setup are done.

For install: [README.md quickstart](README.md).
For setup: [SETUP-INSTALLER.md](docs/alpha/cli/SETUP-INSTALLER.md).

---

## Current operation

A hub created with `cn init` and prepared with `cn setup` becomes an active
coherent agent when a capable model is pointed at it as working context.

Current path:

1. Open Claude UI or Claude Code CLI with the hub directory as context.
2. Ask the model to read the hub identity, skills, memory, config, and threads.
3. Work from the hub. The model is the active reasoning body; the hub is the
   durable identity and memory.

---

## Target runtime (planned)

The following surfaces are target architecture, not shipped today. They are
documented here for design continuity.

### Start / Stop

**Daemon** (long-running, Telegram or peer-only):

```bash
cn agent --daemon              # foreground
systemctl start cn-<name>      # systemd (service name varies by hub)
systemctl stop cn-<name>
```

**Oneshot** (cron, one cycle then exit):

```bash
cn agent                       # single maintenance + drain cycle
```

Crontab entry:
```cron
*/5 * * * * cd /path/to/hub && cn agent >> /var/log/cn.log 2>&1
```

Both modes run the same protocol loop: sync peers, materialize inbox, flush outbox, drain queue. The daemon loops; oneshot exits after one pass.

See [AUTOMATION.md](docs/beta/guides/AUTOMATION.md) for scheduler config, drain limits, and systemd unit setup.

### Configuration

Secrets: `.cn/secrets.env` (auto-loaded, gitignored, chmod 600).

```
ANTHROPIC_KEY=sk-ant-...
TELEGRAM_TOKEN=123456:ABC...
```

Settings: `.cn/config.json`.

```json
{
  "runtime": {
    "model": "claude-sonnet-4-6",
    "max_tokens": 4096,
    "poll_interval": 30,
    "allowed_users": [12345678]
  }
}
```

Full config options: [AUTOMATION.md scheduler settings](docs/beta/guides/AUTOMATION.md).

---

## 1. Observing

### Unified log (`cn logs`)

The unified log is the single entry point for operator observability. One line per event, correlated by message ID.

```bash
cn logs                        # last 50 entries, human-formatted
cn logs --since 2h             # last 2 hours
cn logs --errors               # warnings and errors only
cn logs --msg tg-12345         # trace one message end-to-end
cn logs --json                 # raw JSONL (for scripts)
cn logs --kind invocation.end  # filter by event kind
cn logs -n 100                 # last 100 entries
```

Five event kinds per invocation:

| Kind | Meaning |
|------|---------|
| `message.received` | Message dequeued (or rejected/empty at Warn severity) |
| `invocation.start` | LLM call begins |
| `invocation.end` | LLM call completes (includes pass count, ops, duration) |
| `message.sent` | Response projected to Telegram |
| `error` | LLM failure, N-pass exhaustion, or poll error |

Storage: `logs/unified/YYYYMMDD.jsonl` (schema `cn.ulog.v1`). ~2KB per invocation.

### Status and health

```bash
cn status                      # hub name, versions, package drift
cn doctor                      # 20+ checks: tools, config, packages, extensions
```

`cn doctor` validates: git/curl/patch installed, hub structure intact, package chain consistent (deps.json -> deps.lock.json -> vendor/), runtime contract valid, version coherence, origin remote configured.

### Other log locations

| Log | Path | Use when |
|-----|------|----------|
| Unified log | `logs/unified/YYYYMMDD.jsonl` | First stop. Use `cn logs`. |
| Daemon log | `logs/daemon.log` | Service lifecycle, sync activity |
| Event log | `logs/events/YYYYMMDD.jsonl` | Detailed system telemetry (50+ event types) |
| Input log | `logs/input/tg-<id>.md` | Full prompt sent to LLM |
| Output log | `logs/output/tg-<id>.md` | Agent's raw response |

---

## 2. Maintaining

### Updating the binary

```bash
cn update                      # checks GitHub Releases, downloads if newer
cn --version                   # verify
```

`cn update` detects platform (linux-x64, macos-x64, macos-arm64), downloads atomically, and reconciles packages post-update.

Manual update (if `cn update` fails):

```bash
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
```

After updating, restart the daemon if running: `systemctl restart cn-<name>` (target — when daemon is running).

See [BUILD-RELEASE.md](docs/beta/guides/BUILD-RELEASE.md) for rollback procedure.

### Packages

```bash
cn deps                        # list installed packages
cn deps doctor                 # verify installed matches lockfile
cn deps restore                # reinstall from lockfile
cn deps update                 # re-resolve and update lockfile
```

Desired state: `.cn/deps.json`. Resolved state: `.cn/deps.lock.json`. Installed: `.cn/vendor/packages/`.

See [PACKAGE-SYSTEM.md](docs/alpha/package-system/PACKAGE-SYSTEM.md).

### Peers

```bash
cn peer                        # list peers
cn peer add <name> <url>       # add peer
cn peer sync                   # fetch all peer repos
```

See [HANDSHAKE.md](docs/beta/guides/HANDSHAKE.md) for establishing peer-to-peer coordination.

---

## 3. Releasing

One command:

```bash
scripts/release.sh 3.34.0
```

What it does:
1. **Preflight** — clean tree, on main, synced with origin, tag doesn't exist
2. **Bump** — writes VERSION file
3. **Stamp** — propagates version to cn.json, package manifests
4. **Check** — verifies all version sources agree
5. **Confirm** — shows staged diff, asks before committing
6. **Ship** — commit, tag `v3.34.0`, push main + tag

The release workflow (`.github/workflows/release.yml`) triggers on the tag push and handles everything else: builds 4 platform binaries, generates `checksums.txt`, creates the GitHub release with RELEASE.md notes.

### Before releasing

- Update RELEASE.md with release notes for this version
- Update CHANGELOG.md with the version row
- Merge all PRs for this release to main

### Version source of truth

`VERSION` file → `scripts/stamp-versions.sh` derives cn.json + package manifests → `scripts/check-version-consistency.sh` validates. One source, everything else derived.

---

## 4. Troubleshooting

Start with `cn logs --errors`. Every failure path emits to the unified log.

| Symptom | Command | What to check |
|---------|---------|---------------|
| Agent isn't responding | `cn logs --errors --since 1h` | Poll errors, rejected users, LLM failures |
| Message sent but no reply | `cn logs --msg tg-<id>` | Trace receive -> invoke -> respond |
| Agent responds but does nothing | Check for `(0 ops)` in daemon.log | Normal for conversational replies |
| Daemon keeps restarting | `journalctl -u cn-<name> --since "1h ago"` | Missing secrets, bad permissions |
| Telegram messages not arriving | Check `.cn/secrets.env` for TELEGRAM_TOKEN | One poller per token |
| Unknown peer errors | `cn peer` | Verify peer name and URL |
| Package drift | `cn doctor` then `cn deps restore` | Lockfile vs vendor mismatch |

See [TROUBLESHOOTING.md](docs/beta/guides/TROUBLESHOOTING.md) for detailed diagnostics.

---

## 5. Quick Reference

| Task | Command |
|------|---------|
| View recent activity | `cn logs` |
| Check health | `cn doctor` |
| View hub state | `cn status` |
| Update binary | `cn update` |
| Sync peers | `cn sync` |
| Process one item | `cn agent --process` (planned) |
| Full CLI reference | `cn --help` or [CLI.md](docs/alpha/cli/CLI.md) |
