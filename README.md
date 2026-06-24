# cnos — Coherence Network OS

[![Build](https://github.com/usurobor/cnos/actions/workflows/build.yml/badge.svg)](https://github.com/usurobor/cnos/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](./LICENSE)

cnos gives AI agents a Git-native home.

A cnos hub is a repo an agent can wake into, read, write, remember through, and use to coordinate with humans and other agents. Git stores the durable state. `cn` prepares and checks the hub. A model such as Claude or GPT acts as the temporary reasoning body.

The goal is simple: agents should not just do tasks. They should leave receipts, evidence, decisions, and memory that make the whole system more coherent over time.

## Why cnos exists

AI agents are useful, but agent work often fails in predictable ways.

| Common failure | cnos answer |
|---|---|
| It sounds done, but is not done. | Work closes with receipts that can be checked against a contract. |
| It claims evidence it did not produce. | Receipts point to files, commits, logs, and other durable evidence. |
| It reviews itself. | Production and review are separated in the protocol. |
| It forgets why a decision was made. | Threads, close-outs, and commits leave a durable trail in Git. |
| It fixes one local thing and breaks the larger system. | Work is measured against the coherence of the whole project. |
| It creates follow-up work without grounding. | Follow-up issues should come from receipts, findings, and measured gaps. |
| It learns nothing from repeated failure. | Repeated failures become protocol or package improvements. |

cnos turns agent work from one-off task execution into a traceable path of decreasing incoherence.

## Try it

Install `cn`:

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
```

Create a hub:

```sh
cn init my-agent
cd my-agent
cn setup
cn doctor
```

Then activate the hub in a model body.

## Activate an agent

Activation means a model reads a hub and becomes the active reasoning body for that hub.

The canonical procedure lives at [`cnos.core/skills/agent/activate/SKILL.md`](./src/packages/cnos.core/skills/agent/activate/SKILL.md). It loads, in order:

```text
Kernel → CA skills → Persona → Operator → hub state → identity confirmation
```

### Activate in hosted chat

Use Claude, ChatGPT, or another hosted model with fetch/paste access. Give it the hub URL and ask it to follow the activation skill.

```text
Activate as this hub: {hub-url}
First read the cnos activation skill:
https://raw.githubusercontent.com/usurobor/cnos/main/src/packages/cnos.core/skills/agent/activate/SKILL.md
Then load the hub identity, operator contract, skills, memory, and current state.
Confirm who you are before doing any work.
```

If the model cannot fetch URLs, paste the activation skill and hub files manually.

### Activate in CLI / workspace

Use Claude Code, Codex CLI, or another local coding agent with filesystem access. From inside the hub directory:

```sh
claude -p "Read this hub and activate from its identity, operator contract, skills, memory, and current state. Confirm who you are before doing any work."
```

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

## How the project is organized

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

`cn build` produces tarballs under `dist/packages/`; `cn setup` installs them into a hub at `.cn/vendor/packages/`.

## Read more

| Question | Link |
| --- | --- |
| What is the big idea? | [docs/THESIS.md](./docs/THESIS.md) |
| How does the Git-native protocol work? | [CN whitepaper](./docs/alpha/protocol/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md) |
| What is the target runtime? | [Agent runtime](./docs/alpha/agent-runtime/AGENT-RUNTIME.md) |
| How do agents activate? | [Activate skill](./src/packages/cnos.core/skills/agent/activate/SKILL.md) |
| What is the recursive cell model? | [Cell of Cells](./docs/gamma/essays/CELL-OF-CELLS.md) |
| Why "decreasing incoherence"? | [Decreasing Incoherence](./docs/gamma/essays/DECREASING-INCOHERENCE.md) |
| What is CDD? | [CDD kernel](./src/packages/cnos.cdd/skills/cdd/CDD.md) |
| What is CDS? | [CDS package](./src/packages/cnos.cds/README.md) |
| What is CDR? | [CDR package](./src/packages/cnos.cdr/README.md) |
| What is handoff? | [Handoff package](./src/packages/cnos.handoff/README.md) |

Full docs index: [`docs/README.md`](./docs/README.md).

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
