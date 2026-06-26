# reference/

**Canonical specs, APIs, CLI, schemas — the load-bearing definitions.**

> **Pass 1 overlay.** These bundles currently live under `alpha/`. They are
> referenced from `src/` skills and CI golden files, so their eventual move
> is a package-system change gated on a path-dependency audit — not a
> navigation tweak.

## Protocol (CN)

The protocol feature bundle, in reading order:

- [GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md](../alpha/protocol/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md) — the CN protocol whitepaper (RELEASE v3.0.1).
- [PROTOCOL.md](../alpha/protocol/PROTOCOL.md) — FSM design: state diagrams, transition tables.
- [THREAD-API.md](../alpha/protocol/THREAD-API.md) — agent content API.
- [THREAD-EVENT-MODEL.md](../alpha/protocol/THREAD-EVENT-MODEL.md) — thread event model.
- [MESSAGE-PACKET-TRANSPORT.md](../alpha/protocol/MESSAGE-PACKET-TRANSPORT.md) — fail-closed packet transport (#150).

## Runtime

- [AGENT-RUNTIME.md](../alpha/agent-runtime/AGENT-RUNTIME.md) — CN Shell, typed ops, N-pass orchestration, receipts.
- [PROVIDER-CONTRACT-v1.md](../alpha/agent-runtime/PROVIDER-CONTRACT-v1.md) — the provider contract (LLM and other capability endpoints).
- [HYBRID-LLM-ROUTING.md](../alpha/agent-runtime/HYBRID-LLM-ROUTING.md) — model routing policy.
- [RUNTIME-EXTENSIONS.md](../alpha/runtime-extensions/RUNTIME-EXTENSIONS.md) — capability providers, discovery, isolation.

## CLI

- [CLI.md](../alpha/cli/CLI.md) — command reference.
- [DAEMON.md](../alpha/cli/DAEMON.md) — daemon mode.
- [SETUP-INSTALLER.md](../alpha/cli/SETUP-INSTALLER.md) — install script spec.

## Packages

- [PACKAGE-SYSTEM.md](../alpha/package-system/PACKAGE-SYSTEM.md) — the package model.
- [PACKAGE-AUTHORING.md](../alpha/package-system/PACKAGE-AUTHORING.md) · [PACKAGE-ARTIFACTS.md](../alpha/package-system/PACKAGE-ARTIFACTS.md) · [BUILD-AND-DIST.md](../alpha/package-system/BUILD-AND-DIST.md)

## Composition & schemas

- [CTB](../alpha/ctb/README.md) — triadic agent-composition language (draft).
- [schemas/](../alpha/schemas) — JSON schemas. (Schema *design* docs currently live under [`beta/schema/`](../beta/schema); the two homes consolidate here in a later pass.)
