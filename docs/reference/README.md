# reference/

**Canonical specs, APIs, CLI, schemas — the load-bearing definitions.**

## Protocol (CN)

The protocol feature bundle, in reading order:

- [GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md](protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md) — the CN protocol whitepaper (RELEASE v3.0.1).
- [PROTOCOL.md](protocol/cn/PROTOCOL.md) — FSM design: state diagrams, transition tables.
- [THREAD-API.md](protocol/cn/THREAD-API.md) — agent content API.
- [THREAD-EVENT-MODEL.md](protocol/cn/THREAD-EVENT-MODEL.md) — thread event model.
- [MESSAGE-PACKET-TRANSPORT.md](protocol/cn/MESSAGE-PACKET-TRANSPORT.md) — fail-closed packet transport (#150).

## Runtime

- [AGENT-RUNTIME.md](runtime/AGENT-RUNTIME.md) — CN Shell, typed ops, N-pass orchestration, receipts.
- [PROVIDER-CONTRACT-v1.md](runtime/PROVIDER-CONTRACT-v1.md) — the provider contract (LLM and other capability endpoints).
- [HYBRID-LLM-ROUTING.md](runtime/HYBRID-LLM-ROUTING.md) — model routing policy.
- [RUNTIME-EXTENSIONS.md](runtime/extensions/RUNTIME-EXTENSIONS.md) — capability providers, discovery, isolation.

## CLI

- [CLI.md](cli/CLI.md) — command reference.
- [DAEMON.md](cli/DAEMON.md) — daemon mode.
- [SETUP-INSTALLER.md](cli/SETUP-INSTALLER.md) — install script spec.

## Packages

- [PACKAGE-SYSTEM.md](packages/PACKAGE-SYSTEM.md) — the package model.
- [PACKAGE-AUTHORING.md](packages/PACKAGE-AUTHORING.md) · [PACKAGE-ARTIFACTS.md](packages/PACKAGE-ARTIFACTS.md) · [BUILD-AND-DIST.md](packages/BUILD-AND-DIST.md)

## Composition & schemas

- [CTB](ctb/README.md) — triadic agent-composition language (draft).
- [schemas/](schemas/) — JSON schemas and schema design docs.

## Vocabulary

- [governance/GLOSSARY.md](governance/GLOSSARY.md) — canonical short definitions of cnos terms.
