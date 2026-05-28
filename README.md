# cnos — Coherence Network OS

[![Build](https://github.com/usurobor/cnos/actions/workflows/build.yml/badge.svg)](https://github.com/usurobor/cnos/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](./LICENSE)

**cnos gives AI agents a Git-native home.**

An agent hub is a Git repository that can hold identity, skills, memory, tasks, messages, and evidence. A model can enter that hub, read its state, do bounded work, write results back, and leave a trail that humans and other agents can inspect.

The short version:

```text
Git repo       durable agent home
cn             CLI that creates and manages hubs
packages       skills and protocols an agent can load
threads        messages, tasks, and work history
receipts       evidence that work actually happened
```

The goal is simple: make agent work less ephemeral, less hand-wavy, and easier to review.

## Why this exists

AI agents fail in familiar ways. cnos pairs each failure with a Git-backed answer:

```text
they sound done but are not done           →  work closes with a receipt against a contract
they claim evidence they did not produce   →  evidence lives in files and commits
they forget why a decision was made        →  close-outs leave a trail in the repo
they generate follow-up work ungrounded    →  next moves cite the receipt they came from
they learn nothing from repeated failure   →  repeated gaps surface in the receipt stream
```

The repo is the durable object. The model is the temporary reasoning body.

## Try it

### 1. Install `cn`

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
```

Or build from source:

```sh
git clone https://github.com/usurobor/cnos.git
cd cnos/src/go
go build -o cn ./cmd/cn
```

Requirements:

```text
Unix-like OS
Git
Go 1.22+ if building from source
curl if using install.sh
```

### 2. Create a hub

```sh
cn init my-agent
cd <hub-dir>
cn setup
cn doctor
cn status
```

`cn setup` installs cognitive packages into:

```text
.cn/vendor/packages/
```

At this point you have a prepared hub. The full autonomous `cn agent` runtime is still planned, but the hub can already be activated by a capable model.

## Activate an agent

Activation means a model reads a hub and becomes the active reasoning body for that hub.

The canonical procedure lives at [`cnos.core/skills/agent/activate/SKILL.md`](./src/packages/cnos.core/skills/agent/activate/SKILL.md). It loads, in order:

```text
Kernel → CA skills → Persona → Operator → hub state → identity confirmation
```

There are two common activation paths.

### Option A — hosted chat

Use this for Claude, ChatGPT, or another hosted model with WebFetch or pasted context.

Paste:

```text
Activate as https://github.com/usurobor/cn-sigma.

First read:
https://raw.githubusercontent.com/usurobor/cnos/main/src/packages/cnos.core/skills/agent/activate/SKILL.md

Then follow that skill against the cn-sigma hub.
Load Kernel, CA skills, Persona, Operator, hub state, and confirm identity before doing any other work.
```

If the model cannot fetch URLs, paste the activation skill and hub files manually.

### Option B — CLI / workspace

Use this for Claude Code CLI, Codex CLI, or another agent body with shell and Git access.

```sh
git clone https://github.com/usurobor/cnos.git
git clone https://github.com/usurobor/cn-sigma.git

cd cn-sigma

claude -p "Read ../cnos/src/packages/cnos.core/skills/agent/activate/SKILL.md, then activate against this hub. Load Kernel, CA skills, Persona, Operator, hub state, and confirm identity before doing any other work."
```

The CLI path is best when available because the model can read local files directly.

## What ships today

The Go `cn` binary currently supports the hub and package-management path:

| Command | Purpose |
| --- | --- |
| `cn help` | Show available commands. |
| `cn init [name]` | Create a new agent hub. |
| `cn setup` | Install cognitive packages and write dependency state. |
| `cn deps list` | List installed package dependencies. |
| `cn deps restore` | Restore package dependencies. |
| `cn deps doctor` | Check dependency state. |
| `cn status` | Show hub state. |
| `cn doctor` | Check hub health. |
| `cn build` | Assemble packages from source. Use `--check` for CI. |
| `cn update` | Self-update with SHA-256 verification and atomic install. |

The shipped Go binary does **not** yet include the full agent runtime.

Planned runtime work includes:

```text
cn agent
daemon / scheduler
thread queue processing
inbox / outbox sync
peer sync
CN Shell execution
package command discovery and dispatch
CTB enforcement
```

## How cnos is organized

```text
cnos/
├── src/go/        Go CLI source
├── src/packages/  cognitive packages and skills
├── src/ocaml/     legacy runtime, retained during the Go rewrite
├── docs/          doctrine, protocol, runtime, architecture, essays
├── dist/          built packages and binaries
├── scripts/       build and release scripts
└── test/          tests
```

The current package set:

| Package | What it owns |
| --- | --- |
| `cnos.core` | Core skills and foundational operations. |
| `cnos.cdd` | The generic recursive cell protocol. |
| `cnos.cds` | Software work: code, tests, branches, CI, releases. |
| `cnos.cdr` | Research work: claims, data, methods, reports, review. |
| `cnos.handoff` | Messages and receipts between agents, repos, roles, and packages. |
| `cnos.eng` | Engineering skills for implementation and testing. |
| `cnos.kata` | Runtime verification kata. |
| `cnos.cdd.kata` | CDD-method verification kata. |

`cn build` assembles packages into:

```text
dist/packages/
```

`cn setup` installs them into a hub under:

```text
.cn/vendor/packages/
```

## Core idea

cnos treats work as a small loop:

```text
receive a message
do bounded work
review it
write a receipt
validate the result
use the result to choose the next move
```

In project docs this is called a **coherence cell**.

You do not need to understand the full theory to try cnos. The practical idea is:

```text
make every important piece of agent work leave inspectable evidence
```

## GitHub-hosted agents

A future cnos hub can run without a long-lived server:

```text
GitHub repo      durable agent home
GitHub Actions   wake compute
cn sync          peer/message transport
cn agent         one-shot reasoning loop
model API/CLI    temporary reasoning body
```

The agent wakes, pulls peers, runs the `cn` loop, writes messages or receipts, pushes its repo, and sleeps.

This is planned work. The current Go binary prepares the hub and package state; the full `cn sync` / `cn agent` wake loop is still part of the runtime roadmap.

## Read more

Start here:

| Question | Link |
| --- | --- |
| What is the big idea? | [docs/THESIS.md](./docs/THESIS.md) |
| How does the Git-native protocol work? | [CN whitepaper](./docs/alpha/protocol/WHITEPAPER.md) |
| What is the target runtime? | [Agent runtime](./docs/alpha/agent-runtime/AGENT-RUNTIME.md) |
| How do agents activate? | [Activate skill](./src/packages/cnos.core/skills/agent/activate/SKILL.md) |
| What is the recursive cell model? | [Cell of Cells](./docs/gamma/essays/CELL-OF-CELLS.md) |
| Why "decreasing incoherence"? | [Decreasing Incoherence](./docs/gamma/essays/DECREASING-INCOHERENCE.md) |
| What is CDD? | [CDD kernel](./src/packages/cnos.cdd/skills/cdd/CDD.md) |
| What is CDS? | [CDS package](./src/packages/cnos.cds/README.md) |
| What is CDR? | [CDR package](./src/packages/cnos.cdr/README.md) |
| What is handoff? | [Handoff package](./src/packages/cnos.handoff/README.md) |

Full docs index:

```text
docs/README.md
```

## Contributing

```sh
go test ./...
```

Fork, branch, make a focused change, run tests, and submit a pull request.

Common commit types:

```text
feat
fix
docs
chore
refactor
test
```

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## Support

- Individuals: [GitHub Sponsors](https://github.com/sponsors/usurobor)
- Organizations: [Sustainability](./docs/beta/SUSTAINABILITY.md)

## License

[Apache License 2.0](./LICENSE)
