# git-cn Package Design

**Version:** 0.1.0
**Status:** Draft
**Issue:** #218
**Parent:** POLYGLOT-PACKAGES-AND-PROVIDERS.md, PROVIDER-CONTRACT-v1.md

## Decision

`git-cn` belongs in `cnos.transport.git`, as both a package command and a provider-backed extension. The command is for humans. The provider identity is for the runtime. The package is the shared distribution unit.

## Rationale

The docs now distinguish commands from providers and explicitly say both can live in the same package without becoming the same runtime thing. The polyglot packages/provider design uses `git-cn` as the first-class example of exactly this pattern:

- **command form** for operator use
- **provider form** for runtime transport use

## Two surfaces

### Command surface

For direct operator use, under the package's `commands/` content class.

- `git-cn` = operator-facing executable command
- Invoked as `cn git-cn` (or however command naming settles)
- One-shot execution, exact-dispatch

### Extension/provider surface

Under `extensions/`, declaring `cnos.transport.git` with capabilities such as `transport.git.sync` and `transport.git.packet`.

- `cnos.transport.git` = runtime-facing provider identity
- Subprocess protocol (cn.ext.v1)
- Long-lived or on-demand, invoked by the runtime not the operator

## Kernel vs package/provider boundary

The kernel keeps:

- protocol semantics
- envelope rules
- receipts
- routing invariants
- dedupe / rejection policy

The package/provider owns:

- concrete transport implementation
- Git pack/delta operations
- Sync mechanics

This split is already spelled out in POLYGLOT-PACKAGES-AND-PROVIDERS.md §6.

## Package structure

```text
packages/cnos.transport.git/
  cn.package.json
  commands/
    git-cn/
      cn-git-cn              # or git-cn, depending on final command naming
  extensions/
    cnos.transport.git/
      cn.extension.json
      bin/
        git-cn-linux-x64
        git-cn-linux-arm64
        git-cn-macos-x64
        git-cn-macos-arm64
      schemas/
      docs/
```

This matches the package-system model, which already treats commands and extensions as first-class content classes inside one package substrate.

## Extension manifest

```json
{
  "schema": "cn.extension.v1",
  "name": "cnos.transport.git",
  "version": "0.1.0",
  "interface": "cn.ext.v1",
  "kind": "capability-provider",
  "backend": {
    "kind": "subprocess",
    "command": ["bin/git-cn-linux-x64"]
  },
  "capabilities": [
    "transport.git.sync",
    "transport.git.packet"
  ],
  "ops": [
    {
      "kind": "git_sync",
      "class": "effect",
      "request_schema": "schemas/git_sync.json",
      "response_format": "artifact:json"
    },
    {
      "kind": "git_packet_send",
      "class": "effect",
      "request_schema": "schemas/git_packet_send.json",
      "response_format": "artifact:json"
    }
  ],
  "permissions": {
    "git": true,
    "default_read_only": false,
    "allow_secrets": []
  },
  "engines": {
    "cnos": ">=3.40.0 <4.0.0"
  }
}
```

## Implementation language

Rust. The FSMs that model Git's object graph and state transitions need:

- algebraic types (`enum` + exhaustive `match`)
- memory-sensitive object graph traversal
- CPU-bound pack/delta operations
- correctness over velocity

Go would fight the FSM modeling (no sum types, no exhaustive matching). Rust is the natural production target for the same reason Go is right for the kernel — match the language to the workload.

## Short form

- **Package:** `cnos.transport.git`
- **Command:** `git-cn` (operator-facing)
- **Provider:** `cnos.transport.git` (runtime-facing)
- **Language:** Rust
- **Protocol:** cn.ext.v1 (subprocess + stdio JSON)
- **Kernel owns:** semantics. **Package owns:** implementation.
