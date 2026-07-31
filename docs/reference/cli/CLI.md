# CN CLI Reference

**Status:** Target CLI surface. The shipped binary is the Go `cn`; its commands are
listed in [README "What ships today"](../../../README.md#what-ships-today), the
authoritative shipped set. Commands here outside that set are marked **(planned)** —
they belong to the agent runtime, which is not yet shipped.
**Author:** usurobor (aka Axiom)
**Contributors:** Sigma

---

## Design Principle

**Agent = brain. cn = body.**

Agent thinks and decides. `cn` senses and executes. The agent never touches git, filesystem, or network directly. Everything goes through `cn`.

## Usage

```
cn <command> [options]
```

## Commands

### Agent Decisions (GTD) — planned

The five GTD operations on threads:

```
cn delete <thread>           Discard thread
cn defer <thread> [until]    Postpone thread
cn delegate <thread> <peer>  Forward thread to peer
cn do <thread>               Claim/start thread → threads/doing/
cn done <thread>             Complete thread → threads/archived/
```

Plus direct communication:

```
cn reply <thread> <message>  Append reply to thread
cn send <peer> <message>     Send message to peer (or self)
```

### Agent Output (Structured) — planned

For agent-to-cn structured responses (used by actor loop):

```
cn out do reply --message <msg>
cn out do send --to <peer> --message <msg>
cn out do surface --desc <desc>
cn out do noop --reason <reason>          Reason must be ≥10 chars, non-trivial
cn out do commit --artifact <path>
cn out defer --reason <reason>
cn out delegate --to <peer>
cn out delete --reason <reason>
```

### Agent Runtime — planned

```
cn agent                     Oneshot scheduler: one cycle then exit
cn agent --process           Single-shot: process one queued item and exit
cn agent --daemon            Daemon scheduler: continuous loop + optional Telegram
cn agent --stdio             Interactive REPL (stdin → LLM → stdout)
```

### Orchestration — planned

```
cn sync                      Fetch inbound + send outbound
cn next                      Get next inbox item (with cadence priority)
cn read <thread>             Display thread content
cn inbox [check]             List inbound branches / materialized threads
cn inbox process             Materialize inbound branches as threads
cn inbox flush               Execute queued agent decisions
cn outbox [check]            List outbound threads
cn outbox flush              Push outbox threads to peers
cn queue [list]              Show task queue
cn queue clear               Empty task queue
cn mca [list]                List managed concern aggregations
cn mca add <description>     Surface MCA for community pickup
```

### Thread Creation — planned

```
cn adhoc <title>             Create adhoc thread in threads/adhoc/
cn daily                     Create/show daily reflection
cn weekly                    Create/show weekly reflection
```

### Dependencies

```
cn deps [list]               List installed packages
cn deps restore              Install from lockfile (deterministic)
cn deps doctor               Verify installed assets match lockfile
cn deps add <pkg>            Add dependency to .cn/deps.json        (planned)
cn deps remove <pkg>         Remove dependency                     (planned)
cn deps update [pkg]         Update lockfile (re-resolve in range)  (planned)
cn deps vendor               Commit vendor tree for airgapped use  (planned)
```

Shipped: `list`, `restore`, `doctor`. The rest are planned.

### Build

```
cn build                     Assemble package tarballs in dist/packages/ from src/packages/ sources
cn build --check             Verify built output matches src/packages/ sources (CI mode)
cn build clean               Remove generated content                              (planned)
```

Shipped: `build`, `--check`. `clean` is planned.

### Observability — planned

```
cn logs                      Human-formatted tail (last 50 entries)
cn logs --since <duration>   Filter by time (e.g. 2h, 30m, 1d)
cn logs --msg <id>           Trace single message by correlation ID
cn logs --errors             Show only warnings and errors
cn logs --json               Raw JSONL output
cn logs --kind <kind>        Filter by event kind
cn logs -n <count>           Limit number of entries
```

Reads from `logs/unified/YYYYMMDD.jsonl`. Alias: `l`.

### Hub Management

```
cn init [name]               Create new hub
cn setup                     Interactive hub setup (config, secrets, optional systemd)
cn status                    Show hub state
cn doctor                    Health check
cn update                    Update cn to latest version
cn commit [message]          Stage + commit                                       (planned)
cn push                      Push to origin                                       (planned)
cn save [message]            Commit + push (shorthand)                            (planned)
cn release [version]         Tag + create GitHub release                          (planned)
cn peer [list]               List peers                                           (planned)
cn peer add <name> <url>     Add peer                                             (planned)
cn peer remove <name>        Remove peer                                          (planned)
cn peer sync                 Fetch all peer repos                                 (planned)
```

Shipped: `init`, `setup`, `status`, `doctor`, `update`. The git and peer commands are planned.

### Aliases

```
i = inbox    o = outbox    s = status    d = doctor
l = logs     in = agent
```

### Flags

```
--help, -h       Show help
--version, -V    Show version
--json           Machine-readable output
--quiet, -q      Minimal output
--dry-run        Show what would happen
--verbose, -v    Detailed output
```

## Command Types (Implementation)

> The three implementation sections below (command types, dispatch, file tree)
> describe the **archived OCaml reference implementation**, not the shipped Go
> binary. The active runtime is Go under `src/go/`; the OCaml tree is off `main`
> — see [OCAML-THREAD-REFERENCE.md](../legacy/OCAML-THREAD-REFERENCE.md).

From `cn_lib.ml` — exhaustive variant, compiler warns on missing cases:

```ocaml
type command =
  | Help | Version | Status | Doctor
  | Init of string option
  | Inbox of Inbox.cmd         (* Check | Process | Flush *)
  | Outbox of Outbox.cmd       (* Check | Flush *)
  | Peer of Peer.cmd           (* List | Add | Remove | Sync *)
  | Queue of Queue.cmd         (* List | Clear *)
  | Mca of Mca.cmd             (* List | Add *)
  | Sync | Next
  | Read of string
  | Reply of string * string
  | Send of string * string
  | Gtd of Gtd.cmd             (* Delete | Defer | Delegate | Do | Done *)
  | Out of Out.gtd             (* Structured agent output *)
  | Commit of string option
  | Push
  | Save of string option
  | Agent of Agent.mode        (* Cron | Process | Daemon | Stdio *)
  | Update | Setup
  | Adhoc of string
  | Daily | Weekly
  | Release of string option   (* Tag + GH release; optional version override *)
  | Deps of Deps.cmd           (* List | Restore | Doctor | Add | Remove | Update | Vendor *)
  | Build of Build.cmd          (* Packages | Check | Clean *)
  | Logs of Logs.cmd            (* Show *)
```

## Dispatch

`cn.ml` is ~100 lines of pure routing. It:

1. Parses flags (`--dry-run`, `--json`, etc.)
2. Parses command (string list → `command option`)
3. Finds hub (`Cn_hub.find_hub_path`)
4. Routes to module (`Cn_mail.inbox_check`, `Cn_runtime.run_oneshot`, etc.)

Commands that work without a hub: `help`, `version`, `init`, `update`.
All others require being inside a hub directory.

`Agent` mode loads config (`Cn_config.load`) before dispatch to the runtime.

## Directory Structure

`cn init` creates:

```
cn-<name>/
 +-- .cn/
 |    +-- config.json
 +-- threads/
 |    +-- in/
 |    +-- mail/
 |    |    +-- inbox/
 |    |    +-- outbox/
 |    |    +-- sent/
 |    +-- doing/
 |    +-- deferred/
 |    +-- archived/
 |    +-- adhoc/
 |    +-- reflections/
 |         +-- daily/
 |         +-- weekly/
 +-- state/
 |    +-- peers.md
 |    +-- runtime.md
 |    +-- queue/
 |    +-- mca/
 |    +-- agent.lock
 |    +-- conversation.json
 |    +-- telegram.offset
 +-- logs/
 |    +-- unified/       (operator log: YYYYMMDD.jsonl, schema cn.ulog.v1)
 |    +-- input/
 |    +-- output/
 +-- spec/
```

See [ARCHITECTURE.md](../../architecture/ARCHITECTURE.md) for full directory layout details.

## Implementation

```
src/
 +-- cli/
 |    +-- cn.ml              CLI dispatch
 +-- lib/
 |    +-- cn_lib.ml          Types, parsing, help text (pure)
 |    +-- cn_json.ml         JSON parser/emitter (pure, zero-dep)
 +-- protocol/
 |    +-- cn_protocol.ml     Typed FSMs
 |    +-- cn_protocol.mli    Interface
 +-- ffi/
 |    +-- cn_ffi.ml          System bindings (Fs, Path, Process, Http)
 +-- transport/
 |    +-- cn_io.ml           Protocol I/O
 |    +-- git.ml             Raw git operations
 |    +-- inbox_lib.ml       Inbox utilities
 +-- cmd/
      +-- cn_runtime.ml      Agent runtime orchestrator (dequeue → LLM → finalize)
      +-- cn_context.ml      Context packer (skills, conversation, capabilities, artifacts)
      +-- cn_llm.ml          Claude API client (curl-backed, no --fail)
      +-- cn_telegram.ml     Telegram Bot API client (send, typing, reactions)
      +-- cn_config.ml       Config loader (env vars + .cn/config.json)
      +-- cn_dotenv.ml       .env file loader (.cn/secrets.env)
      +-- cn_agent.ml        Queue, input/output, op execution
      +-- cn_shell.ml        CN Shell: capability runtime (typed ops, N-pass)
      +-- cn_executor.ml     Op executor (dispatch per kind)
      +-- cn_sandbox.ml      Path sandbox (reject escapes, denylist)
      +-- cn_capabilities.ml Capability discovery (budget, allowlists)
      +-- cn_projection.ml   Reply projection (Telegram routing, dedup)
      +-- cn_orchestrator.ml N-pass orchestration (bounded bind loop)
      +-- cn_mail.ml         Inbox/outbox transport
      +-- cn_gtd.ml          GTD commands + thread creation
      +-- cn_mca.ml          MCA commands
      +-- cn_commands.ml     Peer + git commands
      +-- cn_system.ml       Init, setup, update, status, doctor
      +-- cn_hub.ml          Hub discovery, paths, logging
      +-- cn_ulog.ml         Unified operator log (append-only JSONL writer/reader)
      +-- cn_logs.ml         cn logs CLI (filtering, human formatting)
      +-- cn_fmt.ml          Output formatting, timestamps
```

### Build

```bash
dune build src/cli/cn.exe   # Native OCaml binary
```

Installed to `/usr/local/bin/cn`.

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error |

## Non-Goals

- No GUI (terminal only)
- No network services (git is the transport, Telegram is optional projection)
- No cloud dependencies (works offline except LLM API calls)
- No in-call tool loop (agent produces ops post-call; CN Shell executes them)

## Related

- [AGENT-RUNTIME.md](../runtime/AGENT-RUNTIME.md) — Full runtime spec (CN Shell, typed ops, N-pass orchestration, receipts)
- [SECURITY-MODEL.md](../../architecture/security/SECURITY-MODEL.md) — Security architecture
- [SETUP-INSTALLER.md](SETUP-INSTALLER.md) — Install script specification
