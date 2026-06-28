# Provider Contract v1

**Version:** 0.1.0
**Status:** Draft
**Doc-Class:** canonical-spec
**Owns:** provider runtime contract, subprocess protocol, lifecycle, health, execute semantics, receipts, policy boundary
**Does-Not-Own:** package distribution, extension manifest schema, command dispatch, orchestrator IR, package index format

---

## 0. Purpose

This document defines the runtime contract between the cnos kernel and executable capability providers. It is the missing concrete layer beneath the runtime-extensions model.

This spec answers:

1. How does the kernel launch a provider?
2. How does the kernel learn what the provider can do?
3. How are typed ops executed?
4. How are health, failure, receipts, and shutdown handled?
5. How does the kernel remain the policy decision point?

---

## 1. Relationship to existing docs

This document refines, but does not replace:

- **RUNTIME-EXTENSIONS.md** — extension architecture, discovery, lifecycle, open op model
- **PACKAGE-SYSTEM.md** — package content classes and install pipeline
- **RUNTIME-CONTRACT.md** — runtime self-model

Current architecture already says:

- extensions are package contents, not standalone package systems
- extensions default to subprocess execution
- interface governs describe, health, execute, shutdown
- the runtime remains the policy decision point
- extension ops enter the open typed-op registry

This document makes that protocol concrete.

---

## 2. Definitions

### 2.1 Provider

A provider is an executable runtime endpoint that implements one or more typed capabilities.

Examples:

- `cnos.net.http`
- `cnos.transport.git`
- `cnos.llm.local`
- `cnos.llm.anthropic`

A provider is **not**:

- a skill
- a command
- an orchestrator
- a package

A provider is a runtime-facing executable surface that may be shipped inside a package.

### 2.2 Extension

An extension is the package-level declaration that a provider exists and what ops it contributes. The extension manifest points to the provider executable and declares:

- provider identity
- backend kind
- ops
- permissions
- engine compatibility

### 2.3 Provider protocol

The provider protocol is the message contract between the cnos kernel and the provider process.

### 2.4 Typed op

A typed op is the common execution envelope that the kernel validates and then dispatches either to:

- a built-in handler, or
- a provider

This contract only governs provider-backed op execution.

---

## 3. Design Goals

The provider contract must be:

- **simple** — easy to implement in Go, Rust, or another language
- **explicit** — no reflection, no hidden registration
- **isolated** — subprocess by default
- **kernel-governed** — kernel controls op shape and permissions
- **traceable** — every provider action can produce receipts and trace events
- **replaceable** — multiple providers can implement distinct capability families without changing the kernel

---

## 4. Transport

### 4.1 Default transport

The default provider transport is:

- subprocess
- stdio
- line-delimited JSON (NDJSON)

Why:

- language-agnostic
- simple to debug
- easy to spawn
- no local socket coordination
- clear failure boundary

### 4.2 Non-goals for v1

Do not support yet:

- in-process plugin loading
- gRPC
- HTTP loopback servers
- remote provider protocols

Those may come later without changing the provider semantics.

---

## 5. Provider lifecycle

The kernel interacts with a provider through this lifecycle:

1. **spawn**
2. **handshake**
3. **describe**
4. **health**
5. **execute**
6. **shutdown**

### 5.1 Spawn

The kernel launches the provider executable from the installed package tree. The executable path comes from the extension manifest.

### 5.2 Handshake

Immediately after spawn, the kernel sends a protocol hello. Provider must answer with protocol compatibility.

### 5.3 Describe

The provider returns:

- identity
- protocol version
- implemented capabilities
- supported op kinds
- implementation metadata

### 5.4 Health

The kernel may ask whether the provider is healthy before first use and later on demand.

### 5.5 Execute

The kernel sends one typed op request with:

- request ID
- validated op envelope
- effective permissions
- limits
- execution context metadata

The provider returns:

- success or error
- typed result payload
- artifacts if any
- timings if available

### 5.6 Shutdown

The kernel requests a clean shutdown when appropriate. The kernel may also terminate the provider process if:

- timeout
- protocol violation
- repeated health failure
- policy change

---

## 6. Protocol envelope

All protocol messages are JSON objects with:

```json
{
  "id": "uuid-or-monotonic-id",
  "method": "hello|describe|health|execute|shutdown",
  "params": { ... }
}
```

Responses are:

```json
{
  "id": "same-id",
  "ok": true,
  "result": { ... }
}
```

or:

```json
{
  "id": "same-id",
  "ok": false,
  "error": {
    "code": "string",
    "message": "human-readable error",
    "retryable": false
  }
}
```

**Rule:**

- `id` is required for every request/response pair
- providers must never emit unframed freeform output on stdout
- human/debug logs go to stderr
- stdout is reserved for protocol messages only

---

## 7. Handshake

### 7.1 Request

```json
{
  "id": "1",
  "method": "hello",
  "params": {
    "protocol": "cn.ext.v1",
    "kernel_version": "3.40.0"
  }
}
```

### 7.2 Response

```json
{
  "id": "1",
  "ok": true,
  "result": {
    "protocol": "cn.ext.v1",
    "provider_version": "0.1.0",
    "compatible": true
  }
}
```

**Rule:** If the provider:

- does not answer,
- answers malformed JSON,
- or declares an incompatible protocol,

the kernel rejects it and emits `extension.rejected`.

---

## 8. Describe

### 8.1 Request

```json
{
  "id": "2",
  "method": "describe",
  "params": {}
}
```

### 8.2 Response

```json
{
  "id": "2",
  "ok": true,
  "result": {
    "name": "cnos.transport.git",
    "protocol": "cn.ext.v1",
    "version": "0.1.0",
    "capabilities": [
      "transport.git.sync",
      "transport.git.packet"
    ],
    "ops": [
      "git_sync",
      "git_packet_send"
    ]
  }
}
```

**Rule:** The provider description must be consistent with the extension manifest. The runtime compares:

- name
- interface / protocol
- version (if strict version match is enabled)
- capabilities (when present in the manifest)
- ops

If they disagree, the kernel rejects the provider. **Manifest declares. Provider confirms. Kernel decides.**

---

## 9. Health

### 9.1 Request

```json
{
  "id": "3",
  "method": "health",
  "params": {}
}
```

### 9.2 Response

```json
{
  "id": "3",
  "ok": true,
  "result": {
    "healthy": true,
    "details": "ready"
  }
}
```

**Rule:** Health is advisory for scheduling/use. The kernel still owns final policy. Repeated health failure may mark the provider:

- discovered
- but not active

---

## 10. Execute

### 10.1 Kernel responsibilities before execute

The kernel must:

- resolve the op kind to exactly one provider
- validate the common op envelope
- validate effective permissions
- construct the execution limits
- assign a request ID
- record `extension.op.start`

The provider is not allowed to decide its own authority envelope.

### 10.2 Request

```json
{
  "id": "42",
  "method": "execute",
  "params": {
    "op": {
      "kind": "http_get",
      "class": "observe",
      "args": {
        "url": "https://example.com"
      }
    },
    "permissions": {
      "network": true,
      "allowed_domains": ["example.com"],
      "allow_secrets": []
    },
    "limits": {
      "timeout_sec": 30,
      "max_bytes": 65536
    },
    "context": {
      "hub_name": "sigma",
      "trace_id": "..."
    }
  }
}
```

### 10.3 Success response

```json
{
  "id": "42",
  "ok": true,
  "result": {
    "status": "ok",
    "payload": {
      "format": "artifact:text",
      "content": "..."
    },
    "artifacts": [],
    "timings_ms": {
      "total": 184
    }
  }
}
```

### 10.4 Error response

```json
{
  "id": "42",
  "ok": false,
  "error": {
    "code": "domain_not_allowed",
    "message": "requested domain is not permitted",
    "retryable": false
  }
}
```

**Rule:** The provider may only execute within the permission envelope supplied by the kernel. It must not:

- widen permissions,
- request undeclared secrets,
- or interpret missing permissions permissively.

---

## 11. Shutdown

### 11.1 Request

```json
{
  "id": "99",
  "method": "shutdown",
  "params": {}
}
```

### 11.2 Response

```json
{
  "id": "99",
  "ok": true,
  "result": {
    "stopped": true
  }
}
```

The kernel may skip graceful shutdown if the process is already unhealthy or timed out.

---

## 12. Policy boundary

The provider may declare:

- what it implements
- what it needs

The kernel decides:

- whether it is enabled
- whether it is active
- which op kinds it owns
- what permissions it gets
- what limits apply
- whether it is invoked at all

This is the central boundary of the model.

---

## 13. Receipts and traceability

The runtime must emit at least:

- `extension.discovered`
- `extension.loaded`
- `extension.rejected`
- `extension.disabled`
- `extension.health.ok`
- `extension.health.error`
- `extension.op.start`
- `extension.op.ok`
- `extension.op.error`

Each execute must be traceable by:

- request ID
- provider name
- op kind
- result status
- timings
- failure code if any

Providers do not own the trace model. They contribute data to it.

---

## 14. Package and install layout

Providers are shipped inside packages.

**Source:**

```text
src/agent/extensions/<provider-name>/
  cn.extension.json       ← cn.extension.v1 manifest with capabilities field
  host/...
  docs/...
  schemas/...
```

**Built package:**

```text
packages/<pkg>/extensions/<provider-name>/
  cn.extension.json
  host/...
  docs/...
  schemas/...
```

**Installed:**

```text
.cn/vendor/packages/<pkg>@<version>/extensions/<provider-name>/
  cn.extension.json
  host/...
  docs/...
  schemas/...
```

There is no separate provider manifest. The `cn.extension.v1` manifest with `capabilities` is the single declaration surface.

The package is the distribution unit. The provider is the executable runtime surface.

---

## 15. Platform-aware payloads

A provider package may carry platform-specific binaries.

The package metadata remains:

- platform-neutral at the logical package layer
- platform-aware at the artifact resolution layer

The provider contract does not require one language or one binary format. It only requires that the installed entrypoint satisfy this protocol.

---

## 16. First-class example: cnos.transport.git

This spec is designed so a package like `cnos.transport.git` can exist coherently.

It may provide:

- **a command form:**
  - `git-cn`
- **a provider form:**
  - `cnos.transport.git`

The command is for operator use. The provider is for runtime use. Both may live in one package. They do not become the same runtime thing.

---

## 17. Doctor requirements

`cn doctor` should validate:

- provider manifest parseability
- protocol compatibility
- entrypoint existence
- executable bit where relevant
- duplicate effective op kinds
- conflicting provider identities
- health status if enabled
- provider/manifest description mismatch

---

## 18. Runtime contract implications

A later runtime contract revision should surface provider state explicitly, for example under `body.providers` or `body.capabilities.providers`, including:

- installed providers
- active providers
- health
- op kinds contributed
- effective permission envelope summary

This document defines the provider contract first; runtime-contract exposure can follow.

---

## 19. Invariants

1. Providers are package contents, not a second package system.
2. Providers are not commands.
3. Providers are not skills.
4. Providers do not define policy; the kernel does.
5. The default transport is subprocess + stdio JSON.
6. The kernel must be able to reject malformed or incompatible providers deterministically.
7. Protocol stdout is machine-only; human/debug output goes to stderr.

---

## 20. Non-goals

This spec does not:

- redefine extension discovery
- redefine package distribution
- define the command registry
- define orchestrator execution
- require in-process loading
- require one implementation language
- solve provider marketplace/indexing

---

## 21. Short form

- **Package** distributes.
- **Extension** declares.
- **Provider** executes.
- **Kernel** decides.
