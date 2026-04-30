# cnos — Coherence Network OS

[![Build](https://github.com/usurobor/cnos/actions/workflows/build.yml/badge.svg)](https://github.com/usurobor/cnos/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](./LICENSE)

> Status: cnos is under active construction and transitioning toward v4.
> The Go rewrite is in progress. This README separates what ships today from
> the doctrine, protocol, runtime, and language targets that guide the work.

## What is cnos?

cnos is a recurrent coherence system with Git as its lowest durable substrate.

A coherent system preserves an inspectable relation across cycles. In cnos, that relation is carried by doctrine, repositories, commits, threads, packages, runtime receipts, traces, releases, and agents.

Git gives the durable substrate: identity, memory, history, forks, refs, and commits. CN adds the protocol conventions. cn provides the current CLI and the target runtime body. CTB is the emerging language/checker layer for composing coherent agents.

## Layers

```
Recurrent coherence system
├─ doctrine          first principles, conduct, standing, judgment, inheritance
├─ Git substrate     durable identity, history, refs, forks, commits
├─ CN protocol       repo conventions, threads, messages, signatures
├─ cn CLI today      init, setup, status, doctor, build, update
├─ cn runtime target governed typed ops, receipts, bounded execution
├─ coherent agents   sense, compare, choose, act/learn, review
└─ CTB draft         triadic composition + witnessed close-out language/checker
```

## Coherent agents

A coherent agent preserves the boundary it acts on.

An agent begins at the edge between its model and the observable environment. When another agent enters, that edge becomes shared. The conduct rule is: do not close a local gap by opening a larger shared one.

Coherent agency leaves an inspectable surface. Claims remain attached to evidence. Verdicts name the contract they judge. Requests for repair name the gap they saw. When no move preserves every affected boundary, judgment names the boundary protected, the boundary breached, and the debt not called closure.

The core operations are:

| Operation | Role |
| --- | --- |
| MCP | Form the best current picture of boundaries, evidence, constraints, and available moves. |
| CAP | Make the atomic move: act on the world through MCA or update the model through MCI. |
| CLP | Review the result, inspect α/β/γ, patch the weakest axis, and stop before taste becomes false work. |
| CDD | Apply the same boundary-preserving review discipline to cnos itself. |

The agent's direct I/O is pure text. A narrow agent, or skill, may behave like a pure function. Wider agents compose subagents, preserve witnesses, and return close-outs.

The target runtime routes side effects through cn: Git, network, file I/O, typed operations, receipts, and bounded execution. The shipped Go binary does not yet include the full agent runtime.

## Network model

Agents connect through Git. Each agent has a hub: a repository that stores identity, state, threads, packages, and history. Agents communicate by exchanging refs and thread files. A forge may host the repo, but the forge is not the protocol.

| Concept | Meaning |
| --- | --- |
| Hub | A Git repo that holds an agent's identity, state, config, threads, and installed packages. |
| Peer | Another agent hub listed in repo state. |
| Thread | A markdown work or conversation surface with frontmatter and history. |
| Doctrine | The always-on conduct layer: coherence, standing, judgment, inheritance, and review law. |
| CN | The Git-native protocol conventions for hubs, threads, packages, messages, and signatures. |
| cn | The current CLI and target runtime body. |
| Agent | A boundary-preserving articulation of coherence that senses, compares, acts or learns, reviews, and leaves inspectable evidence surfaces. |
| Skill | A narrow agent, usually callable as a bounded task behavior. |
| CTB | Draft language/checker layer for triadic agent composition and witnessed close-outs. |

## Why Git?

For agents, a repo is durable identity. It can be cloned, forked, audited, and moved.

For humans, a repo makes agent work inspectable. Decisions can become commits. Collaborations can become merges. Failures can leave evidence instead of vanishing behind a platform boundary.

For the network, Git keeps the protocol small. No central server owns the memory. No platform is the canonical surface. The durable object is the repo.

## Current state

cnos is moving from a v3 prototype into a modular Go service for v4.

### What ships today

The shipped Go binary provides the kernel CLI path.

| Command | Purpose |
| --- | --- |
| cn help | Show available commands. |
| cn init [name] | Create a new agent hub. |
| cn setup | Install cognitive packages and write the dependency manifest. |
| cn deps list | List installed package dependencies. |
| cn deps restore | Restore package dependencies. |
| cn deps doctor | Check dependency state. |
| cn status | Show hub state. |
| cn doctor | Check hub health. |
| cn build | Assemble packages from source. Use --check for CI. |
| cn update | Self-update with SHA-256 verification and atomic install. |

Today, you can create a hub, install cognitive packages, inspect hub state, and build packages.

### Not shipped yet

The Go binary does not yet include:

- the full agent runtime;
- the daemon or scheduler;
- CN Shell execution;
- thread queue processing;
- inbox/outbox operations;
- peer sync;
- package command discovery and dispatch;
- Telegram integration;
- CTB v0.2 enforcement;
- ctb-check.

## Install

### Release binary

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
```

### From source

```sh
git clone https://github.com/usurobor/cnos.git
cd cnos/src/go
go build -o cn ./cmd/cn
```

Requirements:

- Unix-like OS;
- Git;
- Go 1.22+ for source builds;
- curl for the install script.

## Quick start

```sh
# Create a hub.
cn init my-agent

# Enter the hub directory printed by cn init.
cd <hub-dir>

# Install cognitive packages.
cn setup

# Check the hub.
cn doctor
cn status
```

cn setup installs cognitive packages into:

```
.cn/vendor/packages/
```

After setup, the hub has package content and dependency state, but the full agent runtime is still planned.

## Roadmap to v4

Each phase should ship a working binary.

| Phase | Target | Status |
| --- | --- | --- |
| 1 | CLI skeleton and modular dispatch | complete |
| 2 | Core commands: help, init, status, doctor | complete |
| 3 | Build commands: deps, build, setup, update | complete |
| 4 | Package command discovery and dispatch | next |
| 5 | Agent runtime: scheduler, CN Shell, threads, peers | planned |

Design references:

- [Go kernel commands](./docs/alpha/agent-runtime/GO-KERNEL-COMMANDS.md)
- [Architectural invariants](./docs/alpha/architecture/INVARIANTS.md)
- [Agent runtime](./docs/alpha/agent-runtime/AGENT-RUNTIME.md)

## Repository map

```
cnos/
├── src/
│   ├── go/               active Go CLI source
│   ├── packages/         cognitive packages and skills
│   └── ocaml/            legacy runtime, retained until Phase 5 deletion
├── docs/                 doctrine, protocol, runtime, architecture, CTB, process
├── dist/                 built packages and binaries
├── scripts/              build and release scripts
└── test/                 tests
```

Package source lives under:

```
src/packages/
```

cn build assembles those packages into:

```
dist/packages/
```

cn setup installs built packages into a hub under:

```
.cn/vendor/packages/
```

## Agent hub shape

A hub created by cn init and prepared by cn setup has this target shape:

```
hub/
├── .cn/
│   ├── config.json
│   ├── deps.json
│   └── vendor/packages/
├── spec/
│   ├── SOUL.md
│   └── USER.md
├── threads/
│   ├── in/
│   ├── mail/
│   ├── doing/
│   ├── archived/
│   ├── adhoc/
│   └── reflections/
└── state/
    ├── peers.md
    └── queue/
```

Secrets belong in uncommitted secret surfaces, not in Git history.

## Further reading

cnos is documented in layers. Start with the thesis, then follow doctrine, protocol, runtime, CTB, or TSC depending on the question.

Full index:

- [docs/README.md](./docs/README.md)

### System frame

- [Thesis](./docs/THESIS.md) — cnos as a recurrent coherence system.
- [Coherence System](./docs/alpha/essays/COHERENCE-SYSTEM.md) — coherence as the primary system principle.
- [Foundations](./docs/alpha/essays/FOUNDATIONS.md) — doctrine and coherence loop.
- [Manifesto](./docs/alpha/essays/MANIFESTO.md) — human+AI commons, auditability, forkability, sovereignty.

### Agent doctrine

- [Doctrine README](./docs/alpha/doctrine/README.md) — reading order and composition of the doctrine essays.
- [Coherence for Agents](./docs/alpha/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md) — shared boundaries and inspectable relation.
- [Ethics for Agents](./docs/alpha/doctrine/ethics-for-agents/ETHICS-FOR-AGENTS.md) — affected standing under asymmetry.
- [Judgment for Agents](./docs/alpha/doctrine/judgment-for-agents/JUDGMENT-FOR-AGENTS.md) — boundary selection under forced loss.
- [Inheritance for Agents](./docs/alpha/doctrine/inheritance-for-agents/INHERITANCE-FOR-AGENTS.md) — named failure modes inherited as contestable constraints.

### CN protocol and Git substrate

- [Whitepaper](./docs/alpha/protocol/WHITEPAPER.md) — CN protocol and Git substrate.
- [Protocol](./docs/alpha/protocol/PROTOCOL.md) — protocol states and message surfaces.
- [Security model](./docs/alpha/security/SECURITY-MODEL.md) — trust, sandboxing, signatures, and capability boundaries.

### Agent architecture and runtime

- [CAA](./docs/alpha/agent-runtime/CAA.md) — coherent agent architecture.
- [Agent Runtime](./docs/alpha/agent-runtime/AGENT-RUNTIME.md) — target runtime body, CN Shell, typed ops, receipts, bounded execution.
- [Architecture](./docs/beta/architecture/ARCHITECTURE.md) — how doctrine, specs, runtime, packages, traces, releases, and agents relate.
- [CDD](./docs/gamma/cdd/CDD.md) — coherence-driven development process.
- [Engineering levels](./docs/gamma/ENGINEERING-LEVELS.md) — L5/L6/L7 rubric.
- [Changelog](./CHANGELOG.md) — release coherence ledger.

### CTB language layer

- [CTB README](./docs/alpha/ctb/README.md) — CTB document map and authority rules.
- [CTB language spec v0.1](./docs/alpha/ctb/LANGUAGE-SPEC.md) — current skill-module baseline.
- [CTB language spec v0.2 draft](./docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md) — draft agent-module and composition target.
- [CTB semantics notes](./docs/alpha/ctb/SEMANTICS-NOTES.md) — triadic carrier and agent-composition rationale.
- [CTB vision](./docs/alpha/ctb/CTB-v4.0.0-VISION.md) — strategic vision and roadmap.

Status: CTB v0.2 and ctb-check are draft targets. The shipped runtime does not yet enforce CTB v0.2 witness or composition obligations.

### TSC upstream foundation

CTB's triadic carrier and witness discipline are grounded in the separate TSC repository:

- [TSC repository](https://github.com/usurobor/tsc)
- [C≡ spec](https://github.com/usurobor/tsc/blob/main/spec/c-equiv.md) — term algebra, tri(·,·,·), equivalence, normal forms, α/β/γ evaluators.
- [TSC Core](https://github.com/usurobor/tsc/blob/main/spec/tsc-core.md) — dimensional coherence scores, aggregate C_Σ, confidence intervals, independence, composition bounds.
- [TSC Operational](https://github.com/usurobor/tsc/blob/main/spec/tsc-oper.md) — witnesses, floors, verification controller, verdict logic, provenance bundles.

## Contributing

Fork, branch, make changes, run tests, and submit a pull request.

```sh
go test ./...
```

Commit style: `type: short description`

Common types: feat fix docs chore refactor test

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## Support

cnos is open source.

- Individuals: [GitHub Sponsors](https://github.com/sponsors/usurobor)
- Organizations: [Sustainability](./docs/beta/SUSTAINABILITY.md)

## License

[Apache License 2.0](./LICENSE)
